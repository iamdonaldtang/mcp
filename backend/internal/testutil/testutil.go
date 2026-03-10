// Package testutil provides shared test setup helpers for handler and middleware tests.
package testutil

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/taskon/backend/internal/model"
)

// SetupTestDB creates an in-memory SQLite database and auto-migrates all models.
// StringArray fields (PostgreSQL text[]) are stored as TEXT in SQLite; handlers
// that only read/write the column as a whole work fine, but array-specific SQL won't.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// AutoMigrate all models. SQLite ignores PostgreSQL-specific type hints
	// like "type:uuid" and "type:text[]", which is acceptable for unit tests.
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}
