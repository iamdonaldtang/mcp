package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taskon/backend/pkg/response"
)

// Claims for B-end JWT
type BEndClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Claims for C-end wallet JWT
type CEndClaims struct {
	WalletAddress string `json:"wallet_address"`
	CommunityID   string `json:"community_id"`
	jwt.RegisteredClaims
}

// BEndAuth middleware for B-end routes
func BEndAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		claims := &BEndClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// CEndAuth middleware for C-end routes
func CEndAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		claims := &CEndClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("wallet_address", claims.WalletAddress)
		c.Set("community_id", claims.CommunityID)
		c.Next()
	}
}

// CEndOptionalAuth allows unauthenticated access but sets user context if token exists
func CEndOptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.Next()
			return
		}

		claims := &CEndClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err == nil && token.Valid {
			c.Set("wallet_address", claims.WalletAddress)
			c.Set("community_id", claims.CommunityID)
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
