package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taskon/backend/internal/middleware"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	Secret string
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var user model.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Unauthorized(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c)
		return
	}

	token, err := h.generateBToken(user)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"role":   user.Role,
			"avatar": user.Avatar,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Check if email already exists
	var existing model.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		response.BadRequest(c, "email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c)
		return
	}

	user := model.User{
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
		Role:     "admin",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		response.InternalError(c)
		return
	}

	token, err := h.generateBToken(user)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")

	var user model.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.OK(c, gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role,
		"avatar": user.Avatar,
	})
}

func (h *AuthHandler) generateBToken(user model.User) (string, error) {
	claims := middleware.BEndClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.Secret))
}
