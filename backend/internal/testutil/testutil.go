// Package testutil provides shared test setup helpers for handler and middleware tests.
package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// SetupTestDB creates an in-memory SQLite database and auto-migrates all models.
// It registers a UUID callback to simulate PostgreSQL's gen_random_uuid() default.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Register UUID generation callback (SQLite has no gen_random_uuid())
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

	// Manually create tables without PostgreSQL-specific defaults
	sqlDB, _ := db.DB()

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY, email TEXT UNIQUE NOT NULL, password TEXT NOT NULL,
			name TEXT NOT NULL, role TEXT DEFAULT 'admin', avatar TEXT,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS communities (
			id TEXT PRIMARY KEY, user_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, status TEXT DEFAULT 'draft', brand_color TEXT DEFAULT '#48BB78',
			logo TEXT, enabled_modules TEXT DEFAULT '[]',
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS sectors (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, sort_order INTEGER DEFAULT 0, status TEXT DEFAULT 'active',
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY, sector_id TEXT NOT NULL, community_id TEXT NOT NULL,
			name TEXT NOT NULL, description TEXT, type TEXT NOT NULL,
			status TEXT DEFAULT 'draft', points INTEGER DEFAULT 0,
			icon TEXT DEFAULT 'task_alt', icon_color TEXT DEFAULT '#F59E0B',
			max_completions INTEGER, cooldown_hours INTEGER, requirement TEXT,
			current_completions INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS point_types (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			symbol TEXT NOT NULL, icon TEXT, color TEXT, is_default INTEGER DEFAULT 0,
			created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS level_configs (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, level INTEGER NOT NULL,
			name TEXT NOT NULL, min_points INTEGER NOT NULL, badge TEXT)`,
		`CREATE TABLE IF NOT EXISTS task_chains (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, status TEXT DEFAULT 'draft',
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS task_chain_steps (
			id TEXT PRIMARY KEY, task_chain_id TEXT NOT NULL, "order" INTEGER NOT NULL,
			task_id TEXT NOT NULL, bonus_points INTEGER DEFAULT 0)`,
		`CREATE TABLE IF NOT EXISTS day_chain_configs (
			id TEXT PRIMARY KEY, community_id TEXT UNIQUE NOT NULL,
			enabled INTEGER DEFAULT 0, target_days INTEGER DEFAULT 30,
			daily_points INTEGER DEFAULT 10, catch_up_enabled INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS leaderboard_configs (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, point_type_id TEXT NOT NULL,
			periods TEXT DEFAULT '[]', is_enabled INTEGER DEFAULT 1, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS lb_sprints (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, status TEXT DEFAULT 'draft', point_type_id TEXT NOT NULL,
			starts_at DATETIME, ends_at DATETIME, participants INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS sprint_rewards (
			id TEXT PRIMARY KEY, sprint_id TEXT NOT NULL, tier INTEGER NOT NULL,
			threshold INTEGER NOT NULL, reward_type TEXT NOT NULL,
			reward_value TEXT NOT NULL, quantity INTEGER)`,
		`CREATE TABLE IF NOT EXISTS milestones (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, requirement TEXT, threshold INTEGER NOT NULL,
			reward_type TEXT NOT NULL, reward_value TEXT NOT NULL,
			status TEXT DEFAULT 'active', claims INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS shop_items (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, image TEXT, category TEXT NOT NULL,
			price INTEGER NOT NULL, stock INTEGER DEFAULT -1,
			total_redemptions INTEGER DEFAULT 0, is_time_limited INTEGER DEFAULT 0,
			expires_at DATETIME, status TEXT DEFAULT 'active',
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS lucky_wheel_configs (
			id TEXT PRIMARY KEY, community_id TEXT UNIQUE NOT NULL,
			enabled INTEGER DEFAULT 0, cost_points INTEGER DEFAULT 0,
			cost_point_type_id TEXT, daily_spin_limit INTEGER DEFAULT 1,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS wheel_slots (
			id TEXT PRIMARY KEY, wheel_id TEXT NOT NULL, label TEXT NOT NULL,
			reward_type TEXT NOT NULL, reward_value TEXT,
			probability REAL NOT NULL, color TEXT)`,
		`CREATE TABLE IF NOT EXISTS badges (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, icon TEXT, criteria TEXT,
			auto_award INTEGER DEFAULT 0, awarded_count INTEGER DEFAULT 0,
			created_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS access_rules (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, type TEXT NOT NULL,
			config TEXT DEFAULT '{}', enabled INTEGER DEFAULT 1, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS white_labels (
			id TEXT PRIMARY KEY, community_id TEXT UNIQUE NOT NULL, user_id TEXT NOT NULL,
			status TEXT DEFAULT 'draft', deployment_path TEXT,
			custom_domain TEXT, brand_logo TEXT,
			brand_primary_color TEXT DEFAULT '#F59E0B',
			brand_accent_color TEXT DEFAULT '#F59E0B',
			brand_favicon TEXT, brand_custom_css TEXT,
			is_integration_verified INTEGER DEFAULT 0, first_api_ping_at DATETIME,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS widget_configs (
			id TEXT PRIMARY KEY, white_label_id TEXT NOT NULL,
			module_type TEXT NOT NULL, module_name TEXT NOT NULL,
			is_configured INTEGER DEFAULT 0, is_active INTEGER DEFAULT 0,
			embed_code TEXT, settings TEXT DEFAULT '{}',
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS page_builder_pages (
			id TEXT PRIMARY KEY, white_label_id TEXT NOT NULL,
			name TEXT NOT NULL, slug TEXT NOT NULL,
			layout TEXT DEFAULT '[]', is_published INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS reward_rules (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, trigger_type TEXT NOT NULL,
			trigger_config TEXT DEFAULT '{}', reward_type TEXT NOT NULL,
			reward_value TEXT NOT NULL, is_active INTEGER DEFAULT 1,
			execution_count INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS privileges (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, level_required INTEGER DEFAULT 0,
			point_cost INTEGER, type TEXT NOT NULL,
			config TEXT DEFAULT '{}', is_active INTEGER DEFAULT 1,
			claims INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS sdk_configs (
			id TEXT PRIMARY KEY, white_label_id TEXT UNIQUE NOT NULL,
			project_id TEXT UNIQUE NOT NULL, api_key TEXT NOT NULL,
			api_secret TEXT NOT NULL, webhook_url TEXT,
			allowed_origins TEXT DEFAULT '[]', sso_type TEXT,
			oauth_provider TEXT, oauth_client_id TEXT,
			oauth_secret TEXT, oauth_redirect TEXT,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS contract_entries (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			name TEXT NOT NULL, address TEXT NOT NULL,
			chain TEXT NOT NULL, abi TEXT,
			is_verified INTEGER DEFAULT 0, added_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS community_members (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			wallet_address TEXT NOT NULL, level INTEGER DEFAULT 1,
			total_points INTEGER DEFAULT 0, lifetime_xp INTEGER DEFAULT 0,
			day_streak INTEGER DEFAULT 0, last_check_in DATETIME,
			joined_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS task_completions (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			task_id TEXT NOT NULL, wallet_address TEXT NOT NULL,
			status TEXT NOT NULL, points_earned INTEGER DEFAULT 0,
			completed_at DATETIME, claimed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS shop_redemptions (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			shop_item_id TEXT NOT NULL, wallet_address TEXT NOT NULL,
			points_spent INTEGER NOT NULL, redeemed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS milestone_claims (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			milestone_id TEXT NOT NULL, wallet_address TEXT NOT NULL,
			claimed_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS activity_logs (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			wallet_address TEXT NOT NULL, type TEXT NOT NULL,
			title TEXT NOT NULL, description TEXT,
			points_delta INTEGER, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS referrals (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			referrer_address TEXT NOT NULL, referee_address TEXT NOT NULL,
			referral_code TEXT NOT NULL, status TEXT DEFAULT 'pending',
			points_earned INTEGER DEFAULT 0,
			created_at DATETIME, joined_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS domain_setups (
			id TEXT PRIMARY KEY, white_label_id TEXT UNIQUE NOT NULL,
			domain TEXT NOT NULL, dns_status TEXT DEFAULT 'pending',
			ssl_status TEXT DEFAULT 'pending', cname_target TEXT NOT NULL,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS announcements (
			id TEXT PRIMARY KEY, community_id TEXT NOT NULL,
			title TEXT NOT NULL, description TEXT,
			image TEXT, link TEXT, is_active INTEGER DEFAULT 1,
			created_at DATETIME, expires_at DATETIME)`,
	}

	for _, sql := range tables {
		if _, err := sqlDB.Exec(sql); err != nil {
			t.Fatalf("failed to create table: %v\nSQL: %s", err, sql)
		}
	}

	return db
}
