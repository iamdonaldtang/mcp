package integration

import (
	"net/http"
	"testing"
)

func TestWLFlow_Lifecycle(t *testing.T) {
	r, _ := setupTestApp(t)

	// 1. Register + login
	token := registerAndLogin(t, r)

	// 2. Create community first (WL requires it)
	w := doRequest(r, http.MethodPost, "/api/v1/community", map[string]string{
		"name":        "WL Community",
		"description": "community for whitelabel",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	// 3. Create white label
	w = doRequest(r, http.MethodPost, "/api/v1/whitelabel", map[string]string{
		"deployment_path": "embed",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["deployment_path"] != "embed" {
		t.Fatalf("expected deployment_path 'embed', got %v", data["deployment_path"])
	}

	// 4. Overview should return WL data
	w = doRequest(r, http.MethodGet, "/api/v1/whitelabel/overview", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if data["deployment_path"] != "embed" {
		t.Fatalf("expected overview to contain deployment_path 'embed', got %v", data["deployment_path"])
	}

	// 5. Create page
	w = doRequest(r, http.MethodPost, "/api/v1/whitelabel/pages", map[string]string{
		"name": "Home",
		"slug": "home",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	// 6. List pages → should have 1
	w = doRequest(r, http.MethodGet, "/api/v1/whitelabel/pages", nil, token)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	pages := resp["data"].([]interface{})
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}

	// 7. Create reward rule
	w = doRequest(r, http.MethodPost, "/api/v1/whitelabel/rules", map[string]string{
		"name":         "Welcome Bonus",
		"trigger_type": "task_complete",
		"reward_type":  "points",
		"reward_value": "100",
	}, token)
	assertStatus(t, w, http.StatusCreated)
}

func TestWLFlow_NoCommunity(t *testing.T) {
	r, _ := setupTestApp(t)

	// Register a new user (unique email to avoid collision with registerAndLogin)
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "wl_nocomm@example.com",
		"password": "password123",
		"name":     "No Comm User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	token := data["token"].(string)

	// Try to create WL without community → 400
	w = doRequest(r, http.MethodPost, "/api/v1/whitelabel", map[string]string{
		"deployment_path": "embed",
	}, token)
	assertStatus(t, w, http.StatusBadRequest)

	resp = parseResponse(t, w)
	msg := resp["message"].(string)
	if msg != "create a community first" {
		t.Fatalf("expected error message 'create a community first', got '%s'", msg)
	}
}
