package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/internal/testutil"
)

func setupCommunityRouter(db *gorm.DB) *gin.Engine {
	h := &CommunityHubHandler{DB: db}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	r.GET("/overview", h.GetOverview)
	r.POST("/community", h.Create)
	r.POST("/sectors", h.CreateSector)
	r.GET("/tasks", h.ListTasks)
	r.POST("/tasks", h.CreateTask)
	r.PUT("/tasks/:id", h.UpdateTask)
	r.DELETE("/tasks/:id", h.DeleteTask)
	return r
}

func createTestCommunity(t *testing.T, db *gorm.DB, userID, name string) model.Community {
	t.Helper()
	c := model.Community{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		Description: "Test community",
		Status:      "active",
		BrandColor:  "#48BB78",
	}
	if err := db.Create(&c).Error; err != nil {
		t.Fatalf("failed to create test community: %v", err)
	}
	return c
}

func createTestSector(t *testing.T, db *gorm.DB, communityID, name string, sortOrder int) model.Sector {
	t.Helper()
	s := model.Sector{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		Name:        name,
		SortOrder:   sortOrder,
		Status:      "active",
	}
	if err := db.Create(&s).Error; err != nil {
		t.Fatalf("failed to create test sector: %v", err)
	}
	return s
}

func createTestTask(t *testing.T, db *gorm.DB, communityID, sectorID, name, taskType string) model.Task {
	t.Helper()
	task := model.Task{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		SectorID:    sectorID,
		Name:        name,
		Type:        taskType,
		Status:      "active",
		Points:      10,
		Icon:        "task_alt",
		IconColor:   "#F59E0B",
	}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}
	return task
}

// --- Create tests ---

func TestCommunityCreate_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCommunityRouter(db)

	w := postJSON(r, "/community", map[string]interface{}{
		"name":        "My Community",
		"description": "A test community",
		"brand_color": "#FF0000",
		"modules":     []string{"tasks", "points"},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "My Community" {
		t.Errorf("expected name 'My Community', got %v", data["name"])
	}
	if data["status"] != "draft" {
		t.Errorf("expected status 'draft', got %v", data["status"])
	}

	// Verify default point type was created
	var ptCount int64
	communityID := data["id"].(string)
	db.Model(&model.PointType{}).Where("community_id = ?", communityID).Count(&ptCount)
	if ptCount != 1 {
		t.Errorf("expected 1 default point type, got %d", ptCount)
	}

	// Verify day chain config was created
	var dcCount int64
	db.Model(&model.DayChainConfig{}).Where("community_id = ?", communityID).Count(&dcCount)
	if dcCount != 1 {
		t.Errorf("expected 1 day chain config, got %d", dcCount)
	}
}

func TestCommunityCreate_MissingName(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCommunityRouter(db)

	w := postJSON(r, "/community", map[string]interface{}{
		"description": "No name provided",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- GetOverview tests ---

func TestCommunityGetOverview_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCommunityRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/overview", nil)
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

func TestCommunityGetOverview_WithData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	sector := createTestSector(t, db, community.ID, "General", 1)
	createTestTask(t, db, community.ID, sector.ID, "Task 1", "social")

	r := setupCommunityRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/overview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Test Community" {
		t.Errorf("expected name 'Test Community', got %v", data["name"])
	}
	if data["total_tasks"].(float64) != 1 {
		t.Errorf("expected total_tasks 1, got %v", data["total_tasks"])
	}
}

// --- Sector tests ---

func TestCreateSector_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestCommunity(t, db, "test-user-1", "Test Community")
	r := setupCommunityRouter(db)

	w := postJSON(r, "/sectors", map[string]interface{}{
		"name":        "Getting Started",
		"description": "Onboarding tasks",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Getting Started" {
		t.Errorf("expected name 'Getting Started', got %v", data["name"])
	}
	// Auto sort order should be 1 (first sector)
	if data["sort_order"].(float64) != 1 {
		t.Errorf("expected sort_order 1, got %v", data["sort_order"])
	}
}

func TestCreateSector_NoCommunity(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCommunityRouter(db)

	w := postJSON(r, "/sectors", map[string]interface{}{
		"name": "Orphan Sector",
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Task tests ---

func TestListTasks_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestCommunity(t, db, "test-user-1", "Test Community")
	r := setupCommunityRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	total := data["total"].(float64)
	if total != 0 {
		t.Errorf("expected total 0, got %v", total)
	}
}

func TestCreateTask_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	sector := createTestSector(t, db, community.ID, "General", 1)
	r := setupCommunityRouter(db)

	w := postJSON(r, "/tasks", map[string]interface{}{
		"sector_id": sector.ID,
		"name":      "Follow on Twitter",
		"type":      "social",
		"points":    50,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Follow on Twitter" {
		t.Errorf("expected name 'Follow on Twitter', got %v", data["name"])
	}
	if data["status"] != "draft" {
		t.Errorf("expected status 'draft', got %v", data["status"])
	}
	// Should have default icon
	if data["icon"] != "task_alt" {
		t.Errorf("expected default icon 'task_alt', got %v", data["icon"])
	}
	if data["icon_color"] != "#F59E0B" {
		t.Errorf("expected default icon_color '#F59E0B', got %v", data["icon_color"])
	}
}

func TestUpdateTask(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	sector := createTestSector(t, db, community.ID, "General", 1)
	task := createTestTask(t, db, community.ID, sector.ID, "Old Name", "social")
	r := setupCommunityRouter(db)

	body := map[string]interface{}{
		"name":   "Updated Name",
		"points": 100,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+task.ID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %v", data["name"])
	}
}

func TestDeleteTask(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	sector := createTestSector(t, db, community.ID, "General", 1)
	task := createTestTask(t, db, community.ID, sector.ID, "To Delete", "social")
	r := setupCommunityRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+task.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify task was deleted (soft-deleted)
	var count int64
	db.Model(&model.Task{}).Where("id = ?", task.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected task to be deleted, but found %d", count)
	}

	// Verify task still exists with Unscoped (soft delete)
	var unscopedCount int64
	db.Unscoped().Model(&model.Task{}).Where("id = ?", task.ID).Count(&unscopedCount)
	if unscopedCount != 1 {
		t.Errorf("expected soft-deleted task to exist, found %d", unscopedCount)
	}
}

