package integration

import (
	"net/http"
	"testing"
)

func TestAuthFlow_RegisterLoginProfile(t *testing.T) {
	r, _ := setupTestApp(t)

	// 1. Register
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "auth@example.com",
		"password": "password123",
		"name":     "Auth User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	regToken, ok := data["token"].(string)
	if !ok || regToken == "" {
		t.Fatal("register did not return a token")
	}

	// 2. Login
	w = doRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "auth@example.com",
		"password": "password123",
	}, "")
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	loginToken, ok := data["token"].(string)
	if !ok || loginToken == "" {
		t.Fatal("login did not return a token")
	}

	// 3. Profile with token
	w = doRequest(r, http.MethodGet, "/api/v1/auth/profile", nil, loginToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if data["name"] != "Auth User" {
		t.Fatalf("expected name 'Auth User', got %v", data["name"])
	}
	if data["email"] != "auth@example.com" {
		t.Fatalf("expected email 'auth@example.com', got %v", data["email"])
	}

	// 4. Profile without token → 401
	w = doRequest(r, http.MethodGet, "/api/v1/auth/profile", nil, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

func TestAuthFlow_DuplicateRegister(t *testing.T) {
	r, _ := setupTestApp(t)

	// 1. Register once
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "dup@example.com",
		"password": "password123",
		"name":     "Dup User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	// 2. Register same email again → 400
	w = doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "dup@example.com",
		"password": "password456",
		"name":     "Dup User 2",
	}, "")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAuthFlow_WrongPassword(t *testing.T) {
	r, _ := setupTestApp(t)

	// Register
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":    "wrong@example.com",
		"password": "password123",
		"name":     "Wrong Pass User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	// Login with wrong password → 401
	w = doRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "wrong@example.com",
		"password": "wrongpassword",
	}, "")
	assertStatus(t, w, http.StatusUnauthorized)
}
