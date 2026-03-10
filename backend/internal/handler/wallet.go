package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taskon/backend/internal/middleware"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"gorm.io/gorm"
)

type WalletHandler struct {
	DB     *gorm.DB
	Secret string
}

type WalletConnectRequest struct {
	Address     string `json:"address" binding:"required"`
	CommunityID string `json:"community_id" binding:"required"`
}

func (h *WalletHandler) Connect(c *gin.Context) {
	var req WalletConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Check community exists
	var community model.Community
	if err := h.DB.First(&community, "id = ?", req.CommunityID).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	// Upsert community member
	var member model.CommunityMember
	err := h.DB.Where("community_id = ? AND wallet_address = ?", req.CommunityID, req.Address).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		member = model.CommunityMember{
			CommunityID:   req.CommunityID,
			WalletAddress: req.Address,
			Level:         1,
			JoinedAt:      time.Now(),
		}
		if err := h.DB.Create(&member).Error; err != nil {
			response.InternalError(c)
			return
		}
	}

	// Generate C-end JWT
	token, err := h.generateCToken(req.Address, req.CommunityID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, gin.H{
		"token":   token,
		"member":  member,
	})
}

func (h *WalletHandler) generateCToken(address, communityID string) (string, error) {
	claims := middleware.CEndClaims{
		WalletAddress: address,
		CommunityID:   communityID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.Secret))
}
