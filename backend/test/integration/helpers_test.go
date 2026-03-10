package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/handler"
	"github.com/taskon/backend/internal/middleware"
)

const testSecret = "test-secret"

// setupTestApp creates an in-memory SQLite DB, auto-migrates all models,
// and wires up the full Gin router with all routes.
func setupTestApp(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	// Register a UUID generation callback for models that use gen_random_uuid()
	db.Callback().Create().Before("gorm:create").Register("set_uuid", func(tx *gorm.DB) {
		if tx.Statement.Schema == nil {
			return
		}
		for _, field := range tx.Statement.Schema.PrimaryFields {
			if field.DBName == "id" {
				val, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
				strVal, _ := val.(string)
				if isZero || strVal == "" {
					_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, uuid.New().String())
				}
			}
		}
	})

	// Create tables with raw SQL to avoid PostgreSQL-specific syntax issues in SQLite
	sqlDB, _ := db.DB()
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, email TEXT UNIQUE NOT NULL, password TEXT NOT NULL, name TEXT NOT NULL, role TEXT DEFAULT 'admin', avatar TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS communities (id TEXT PRIMARY KEY, user_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, status TEXT DEFAULT 'draft', brand_color TEXT DEFAULT '#48BB78', logo TEXT, enabled_modules TEXT DEFAULT '[]', created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS sectors (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, sort_order INTEGER DEFAULT 0, status TEXT DEFAULT 'active', created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS tasks (id TEXT PRIMARY KEY, sector_id TEXT NOT NULL, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, type TEXT NOT NULL, status TEXT DEFAULT 'draft', points INTEGER DEFAULT 0, icon TEXT DEFAULT 'task_alt', icon_color TEXT DEFAULT '#F59E0B', max_completions INTEGER, cooldown_hours INTEGER, requirement TEXT, current_completions INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS point_types (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, symbol TEXT NOT NULL, icon TEXT, color TEXT, is_default INTEGER DEFAULT 0, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS day_chain_configs (id TEXT PRIMARY KEY, community_id TEXT UNIQUE NOT NULL, enabled INTEGER DEFAULT 0, target_days INTEGER DEFAULT 30, daily_points INTEGER DEFAULT 10, catch_up_enabled INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS leaderboard_configs (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, point_type_id TEXT NOT NULL, periods TEXT DEFAULT '[]', is_enabled INTEGER DEFAULT 1, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS lb_sprints (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, status TEXT DEFAULT 'draft', point_type_id TEXT NOT NULL, starts_at DATETIME, ends_at DATETIME, participants INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS milestones (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, requirement TEXT, threshold INTEGER NOT NULL, reward_type TEXT NOT NULL, reward_value TEXT NOT NULL, status TEXT DEFAULT 'active', claims INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS shop_items (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, image TEXT, category TEXT NOT NULL, price INTEGER NOT NULL, stock INTEGER DEFAULT -1, total_redemptions INTEGER DEFAULT 0, is_time_limited INTEGER DEFAULT 0, expires_at DATETIME, status TEXT DEFAULT 'active', created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS white_labels (id TEXT PRIMARY KEY, community_id TEXT UNIQUE NOT NULL, user_id TEXT NOT NULL, status TEXT DEFAULT 'draft', deployment_path TEXT, custom_domain TEXT, brand_logo TEXT, brand_primary_color TEXT DEFAULT '#F59E0B', brand_accent_color TEXT DEFAULT '#F59E0B', brand_favicon TEXT, brand_custom_css TEXT, is_integration_verified INTEGER DEFAULT 0, first_api_ping_at DATETIME, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS page_builder_pages (id TEXT PRIMARY KEY, white_label_id TEXT NOT NULL, name TEXT NOT NULL, slug TEXT NOT NULL, layout TEXT DEFAULT '[]', is_published INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS reward_rules (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, trigger_type TEXT NOT NULL, trigger_config TEXT DEFAULT '{}', reward_type TEXT NOT NULL, reward_value TEXT NOT NULL, is_active INTEGER DEFAULT 1, execution_count INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS community_members (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, wallet_address TEXT NOT NULL, level INTEGER DEFAULT 1, total_points INTEGER DEFAULT 0, lifetime_xp INTEGER DEFAULT 0, day_streak INTEGER DEFAULT 0, last_check_in DATETIME, joined_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS task_completions (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, task_id TEXT NOT NULL, wallet_address TEXT NOT NULL, status TEXT NOT NULL, points_earned INTEGER DEFAULT 0, completed_at DATETIME, claimed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS shop_redemptions (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, shop_item_id TEXT NOT NULL, wallet_address TEXT NOT NULL, points_spent INTEGER NOT NULL, redeemed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS milestone_claims (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, milestone_id TEXT NOT NULL, wallet_address TEXT NOT NULL, claimed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS activity_logs (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, wallet_address TEXT NOT NULL, type TEXT NOT NULL, title TEXT NOT NULL, description TEXT, points_delta INTEGER, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS domain_setups (id TEXT PRIMARY KEY, white_label_id TEXT UNIQUE NOT NULL, domain TEXT NOT NULL, dns_status TEXT DEFAULT 'pending', ssl_status TEXT DEFAULT 'pending', cname_target TEXT NOT NULL, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS announcements (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, title TEXT NOT NULL, description TEXT, image TEXT, link TEXT, is_active INTEGER DEFAULT 1, created_at DATETIME, expires_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS widget_configs (id TEXT PRIMARY KEY, white_label_id TEXT NOT NULL, module_type TEXT NOT NULL, module_name TEXT NOT NULL, is_configured INTEGER DEFAULT 0, is_active INTEGER DEFAULT 0, embed_code TEXT, settings TEXT DEFAULT '{}', created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS privileges (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL, description TEXT, level_required INTEGER DEFAULT 0, point_cost INTEGER, type TEXT NOT NULL, config TEXT DEFAULT '{}', is_active INTEGER DEFAULT 1, claims INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS referrals (id TEXT PRIMARY KEY, community_id TEXT NOT NULL, referrer_address TEXT NOT NULL, referee_address TEXT NOT NULL, referral_code TEXT NOT NULL, status TEXT DEFAULT 'pending', points_earned INTEGER DEFAULT 0, created_at DATETIME, joined_at DATETIME)`,
	}
	for _, sql := range tables {
		if _, err := sqlDB.Exec(sql); err != nil {
			t.Fatalf("failed to create table: %v\nSQL: %s", err, sql)
		}
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	authH := &handler.AuthHandler{DB: db, Secret: testSecret}
	walletH := &handler.WalletHandler{DB: db, Secret: testSecret}
	commH := &handler.CommunityHubHandler{DB: db}
	cCommH := &handler.CCommunityHandler{DB: db}
	wlH := &handler.WLHubHandler{DB: db}

	// B-end public auth routes
	r.POST("/api/v1/auth/register", authH.Register)
	r.POST("/api/v1/auth/login", authH.Login)

	// B-end protected routes
	bAuth := r.Group("/api/v1")
	bAuth.Use(middleware.BEndAuth(testSecret))
	{
		bAuth.GET("/auth/profile", authH.Profile)

		bAuth.GET("/community/overview", commH.GetOverview)
		bAuth.POST("/community", commH.Create)
		bAuth.POST("/community/sectors", commH.CreateSector)
		bAuth.GET("/community/tasks", commH.ListTasks)
		bAuth.POST("/community/tasks", commH.CreateTask)
		bAuth.PUT("/community/tasks/:id", commH.UpdateTask)
		bAuth.DELETE("/community/tasks/:id", commH.DeleteTask)

		bAuth.GET("/whitelabel/overview", wlH.GetOverview)
		bAuth.POST("/whitelabel", wlH.Create)
		bAuth.GET("/whitelabel/pages", wlH.ListPages)
		bAuth.POST("/whitelabel/pages", wlH.CreatePage)
		bAuth.POST("/whitelabel/rules", wlH.CreateRewardRule)
	}

	// C-end public
	r.POST("/api/c/wallet/connect", walletH.Connect)

	// C-end protected
	cAuth := r.Group("/api/c")
	cAuth.Use(middleware.CEndAuth(testSecret))
	{
		cAuth.GET("/home", cCommH.Home)
		cAuth.GET("/tasks", cCommH.Tasks)
		cAuth.POST("/daychain/checkin", cCommH.DayChainCheckIn)
		cAuth.GET("/leaderboard", cCommH.Leaderboard)
		cAuth.GET("/shop", cCommH.ShopItems)
		cAuth.POST("/shop/redeem", cCommH.RedeemShopItem)
		cAuth.GET("/milestones", cCommH.Milestones)
		cAuth.POST("/milestones/:id/claim", cCommH.ClaimMilestone)
	}

	return r, db
}

// doRequest performs an HTTP request against the Gin engine and returns the recorder.
func doRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseResponse unmarshals the standard {code, data, message} envelope.
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response body: %v\nbody: %s", err, w.Body.String())
	}
	return resp
}

// assertStatus checks the HTTP status code.
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Fatalf("expected status %d, got %d\nbody: %s", expected, w.Code, w.Body.String())
	}
}

// registerAndLogin registers a user and returns the JWT token.
func registerAndLogin(t *testing.T, r *gin.Engine) string {
	t.Helper()
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatal("register did not return a token")
	}
	return token
}
