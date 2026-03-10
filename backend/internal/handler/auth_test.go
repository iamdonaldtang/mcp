package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/internal/testutil"
)

const testJWTSecret = "test-secret"

func init() {
	gin.SetMode(gin.TestMode)
}

// --- helpers ---

func setupAuthRouter(db *gorm.DB) *gin.Engine {
	h := &AuthHandler{DB: db, Secret: testJWTSecret}
	r := gin.New()
	r.POST("/login", h.Login)
	r.POST("/register", h.Register)
	r.GET("/profile", func(c *gin.Context) {
		c.Set("user_id", c.Query("uid"))
		c.Next()
	}, h.Profile)
	return r
}

func createTestUser(t *testing.T, db *gorm.DB, id, email, password, name string) model.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	user := model.User{
		ID:       id,
		Email:    email,
		Password: string(hash),
		Name:     name,
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func postJSON(r *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v\nbody: %s", err, w.Body.String())
	}
	return resp
}

// --- Login tests ---

func TestLogin_ValidCredentials(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestUser(t, db, "u1", "alice@example.com", "password123", "Alice")
	r := setupAuthRouter(db)

	w := postJSON(r, "/login", LoginRequest{Email: "alice@example.com", Password: "password123"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["token"] == nil || data["token"] == "" {
		t.Error("expected token in response")
	}
	user := data["user"].(map[string]interface{})
	if user["email"] != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %v", user["email"])
	}
}

func TestLogin_InvalidEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestUser(t, db, "u1", "alice@example.com", "password123", "Alice")
	r := setupAuthRouter(db)

	w := postJSON(r, "/login", LoginRequest{Email: "nobody@example.com", Password: "password123"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestUser(t, db, "u1", "alice@example.com", "password123", "Alice")
	r := setupAuthRouter(db)

	w := postJSON(r, "/login", LoginRequest{Email: "alice@example.com", Password: "wrongpass"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupAuthRouter(db)

	// Missing email
	w := postJSON(r, "/login", map[string]string{"password": "password123"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing email, got %d", w.Code)
	}

	// Missing password
	w = postJSON(r, "/login", map[string]string{"email": "alice@example.com"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing password, got %d", w.Code)
	}
}

// --- Register tests ---

func TestRegister_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupAuthRouter(db)

	w := postJSON(r, "/register", RegisterRequest{
		Email:    "bob@example.com",
		Password: "password123",
		Name:     "Bob",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["token"] == nil || data["token"] == "" {
		t.Error("expected token in response")
	}

	// Verify user was persisted
	var count int64
	db.Model(&model.User{}).Where("email = ?", "bob@example.com").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestUser(t, db, "u1", "alice@example.com", "password123", "Alice")
	r := setupAuthRouter(db)

	w := postJSON(r, "/register", RegisterRequest{
		Email:    "alice@example.com",
		Password: "password123",
		Name:     "Alice2",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegister_MissingName(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupAuthRouter(db)

	w := postJSON(r, "/register", map[string]string{
		"email":    "bob@example.com",
		"password": "password123",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupAuthRouter(db)

	w := postJSON(r, "/register", RegisterRequest{
		Email:    "bob@example.com",
		Password: "abc",
		Name:     "Bob",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for short password, got %d", w.Code)
	}
}

// --- Profile tests ---

func TestProfile_ValidUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	createTestUser(t, db, "u1", "alice@example.com", "password123", "Alice")
	r := setupAuthRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/profile?uid=u1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Alice" {
		t.Errorf("expected name Alice, got %v", data["name"])
	}
}

func TestProfile_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupAuthRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/profile?uid=nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
