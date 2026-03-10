package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/handler"
	"github.com/taskon/backend/internal/middleware"
	"github.com/taskon/backend/internal/model"
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
			if field.DBName == "id" && field.DataType == "string" {
				val, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
				if isZero || val == "" {
					_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, uuid.New().String())
				}
			}
		}
	})

	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		t.Fatalf("failed to migrate: %v", err)
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
