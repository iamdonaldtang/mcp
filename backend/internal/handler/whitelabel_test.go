package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/internal/testutil"
)

func setupWLRouter(db *gorm.DB) *gin.Engine {
	h := &WLHubHandler{DB: db}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	r.GET("/wl/overview", h.GetOverview)
	r.POST("/wl", h.Create)
	r.GET("/wl/pages", h.ListPages)
	r.POST("/wl/pages", h.CreatePage)
	r.POST("/wl/rewards", h.CreateRewardRule)
	return r
}

func createTestWhiteLabel(t *testing.T, db *gorm.DB, communityID, userID, path string) model.WhiteLabel {
	t.Helper()
	wl := model.WhiteLabel{
		ID:             uuid.New().String(),
		CommunityID:    communityID,
		UserID:         userID,
		Status:         "draft",
		DeploymentPath: path,
	}
	if err := db.Create(&wl).Error; err != nil {
		t.Fatalf("failed to create test white label: %v", err)
	}
	return wl
}

// --- GetOverview tests ---

func TestWLGetOverview_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupWLRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/wl/overview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["state"] != "empty" {
		t.Errorf("expected state 'empty', got %v", data["state"])
	}
}

func TestWLGetOverview_WithData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestWhiteLabel(t, db, community.ID, "test-user-1", "embed")

	r := setupWLRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/wl/overview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["deployment_path"] != "embed" {
		t.Errorf("expected deployment_path 'embed', got %v", data["deployment_path"])
	}
	if data["status"] != "draft" {
		t.Errorf("expected status 'draft', got %v", data["status"])
	}
}

// --- Create tests ---

func TestWLCreate_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestCommunity(t, db, "test-user-1", "Test Community")
	r := setupWLRouter(db)

	w := postJSON(r, "/wl", map[string]interface{}{
		"deployment_path": "domain",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["deployment_path"] != "domain" {
		t.Errorf("expected deployment_path 'domain', got %v", data["deployment_path"])
	}
	if data["status"] != "draft" {
		t.Errorf("expected status 'draft', got %v", data["status"])
	}

	// Verify persisted
	var count int64
	db.Model(&model.WhiteLabel{}).Where("user_id = ?", "test-user-1").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 white label, got %d", count)
	}
}

func TestWLCreate_NoCommunity(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupWLRouter(db)

	w := postJSON(r, "/wl", map[string]interface{}{
		"deployment_path": "embed",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Page tests ---

func TestWLListPages_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestWhiteLabel(t, db, community.ID, "test-user-1", "embed")

	r := setupWLRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/wl/pages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("expected empty pages list, got %d items", len(data))
	}
}

func TestWLCreatePage_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestWhiteLabel(t, db, community.ID, "test-user-1", "embed")

	r := setupWLRouter(db)

	w := postJSON(r, "/wl/pages", map[string]interface{}{
		"name": "Home Page",
		"slug": "home",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Home Page" {
		t.Errorf("expected name 'Home Page', got %v", data["name"])
	}
	if data["slug"] != "home" {
		t.Errorf("expected slug 'home', got %v", data["slug"])
	}
	if data["layout"] != "[]" {
		t.Errorf("expected default layout '[]', got %v", data["layout"])
	}
}

// --- Reward Rule tests ---

func TestWLCreateRewardRule_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestCommunity(t, db, "test-user-1", "Test Community")
	r := setupWLRouter(db)

	w := postJSON(r, "/wl/rewards", map[string]interface{}{
		"name":         "Task Complete Bonus",
		"trigger_type": "task_complete",
		"reward_type":  "points",
		"reward_value": "50",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Task Complete Bonus" {
		t.Errorf("expected name 'Task Complete Bonus', got %v", data["name"])
	}
	if data["trigger_type"] != "task_complete" {
		t.Errorf("expected trigger_type 'task_complete', got %v", data["trigger_type"])
	}

	// Verify persisted
	var count int64
	db.Model(&model.RewardRule{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 reward rule, got %d", count)
	}
}
