package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCommunityFlow_FullLifecycle(t *testing.T) {
	r, _ := setupTestApp(t)

	// 1. Register + login
	token := registerAndLogin(t, r)

	// 2. Overview should be empty
	w := doRequest(r, http.MethodGet, "/api/v1/community/overview", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["state"] != "empty" {
		t.Fatalf("expected state 'empty', got %v", data["state"])
	}

	// 3. Create community
	w = doRequest(r, http.MethodPost, "/api/v1/community", map[string]string{
		"name":        "Test Community",
		"description": "test desc",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	// 4. Overview should now have the community name
	w = doRequest(r, http.MethodGet, "/api/v1/community/overview", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if data["name"] != "Test Community" {
		t.Fatalf("expected community name 'Test Community', got %v", data["name"])
	}

	// 5. Create sector
	w = doRequest(r, http.MethodPost, "/api/v1/community/sectors", map[string]string{
		"name": "General",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	resp = parseResponse(t, w)
	sectorData := resp["data"].(map[string]interface{})
	sectorID := sectorData["id"].(string)
	if sectorID == "" {
		t.Fatal("sector creation did not return an id")
	}

	// 6. Create task
	w = doRequest(r, http.MethodPost, "/api/v1/community/tasks", map[string]interface{}{
		"sector_id": sectorID,
		"name":      "Follow on Twitter",
		"type":      "social",
		"points":    50,
	}, token)
	assertStatus(t, w, http.StatusCreated)

	resp = parseResponse(t, w)
	taskData := resp["data"].(map[string]interface{})
	taskID := taskData["id"].(string)
	if taskID == "" {
		t.Fatal("task creation did not return an id")
	}

	// 7. List tasks → should have 1
	w = doRequest(r, http.MethodGet, "/api/v1/community/tasks", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 task, got %d", len(items))
	}

	// 8. Update task status to active
	w = doRequest(r, http.MethodPut, fmt.Sprintf("/api/v1/community/tasks/%s", taskID), map[string]string{
		"status": "active",
	}, token)
	assertStatus(t, w, http.StatusOK)

	// 9. Delete task
	w = doRequest(r, http.MethodDelete, fmt.Sprintf("/api/v1/community/tasks/%s", taskID), nil, token)
	assertStatus(t, w, http.StatusOK)

	// 10. List tasks → should have 0
	w = doRequest(r, http.MethodGet, "/api/v1/community/tasks", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	items = data["items"].([]interface{})
	if len(items) != 0 {
		t.Fatalf("expected 0 tasks after delete, got %d", len(items))
	}
}
