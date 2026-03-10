package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func init() {
	gin.SetMode(gin.TestMode)
}

// --- helpers ---

func generateBEndToken(secret string, claims BEndClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func generateCEndToken(secret string, claims CEndClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func validBEndClaims() BEndClaims {
	return BEndClaims{
		UserID: "user-1",
		Email:  "alice@example.com",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func validCEndClaims() CEndClaims {
	return CEndClaims{
		WalletAddress: "0xABCDEF1234567890",
		CommunityID:   "comm-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

// echoHandler writes 200 and echoes context values so tests can assert on them.
func echoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":        c.GetString("user_id"),
			"email":          c.GetString("email"),
			"role":           c.GetString("role"),
			"wallet_address": c.GetString("wallet_address"),
			"community_id":   c.GetString("community_id"),
		})
	}
}

// --- BEndAuth tests ---

func TestBEndAuth_ValidToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", BEndAuth(testSecret), echoHandler())

	token := generateBEndToken(testSecret, validBEndClaims())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestBEndAuth_MissingToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", BEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestBEndAuth_ExpiredToken(t *testing.T) {
	claims := validBEndClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	token := generateBEndToken(testSecret, claims)

	r := gin.New()
	r.GET("/test", BEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestBEndAuth_InvalidToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", BEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestBEndAuth_WrongSecret(t *testing.T) {
	token := generateBEndToken("other-secret", validBEndClaims())

	r := gin.New()
	r.GET("/test", BEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- CEndAuth tests ---

func TestCEndAuth_ValidToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", CEndAuth(testSecret), echoHandler())

	token := generateCEndToken(testSecret, validCEndClaims())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCEndAuth_MissingToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", CEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCEndAuth_ExpiredToken(t *testing.T) {
	claims := validCEndClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	token := generateCEndToken(testSecret, claims)

	r := gin.New()
	r.GET("/test", CEndAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- CEndOptionalAuth tests ---

func TestCEndOptionalAuth_NoToken_Continues(t *testing.T) {
	r := gin.New()
	r.GET("/test", CEndOptionalAuth(testSecret), echoHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCEndOptionalAuth_ValidToken_SetsContext(t *testing.T) {
	r := gin.New()
	r.GET("/test", CEndOptionalAuth(testSecret), func(c *gin.Context) {
		wallet := c.GetString("wallet_address")
		if wallet == "" {
			t.Error("expected wallet_address in context")
		}
		c.JSON(http.StatusOK, nil)
	})

	token := generateCEndToken(testSecret, validCEndClaims())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCEndOptionalAuth_InvalidToken_ContinuesWithoutContext(t *testing.T) {
	r := gin.New()
	r.GET("/test", CEndOptionalAuth(testSecret), func(c *gin.Context) {
		wallet := c.GetString("wallet_address")
		if wallet != "" {
			t.Error("expected empty wallet_address for invalid token")
		}
		c.JSON(http.StatusOK, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// --- CORS tests ---

func TestCORS_SetsHeaders(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected ACAO=*, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("expected Access-Control-Allow-Methods header")
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Error("expected Access-Control-Allow-Headers header")
	}
}

func TestCORS_OPTIONS_Returns204(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
