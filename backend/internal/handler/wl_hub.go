package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"gorm.io/gorm"
)

type WLHubHandler struct {
	DB *gorm.DB
}

func (h *WLHubHandler) GetOverview(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.OK(c, gin.H{"state": "empty"})
		return
	}

	response.OK(c, wl)
}

func (h *WLHubHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	// Ensure user has a community first
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.BadRequest(c, "create a community first")
		return
	}

	var req struct {
		DeploymentPath string `json:"deployment_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wl := model.WhiteLabel{
		CommunityID:    community.ID,
		UserID:         userID,
		Status:         "draft",
		DeploymentPath: req.DeploymentPath,
	}

	if err := h.DB.Create(&wl).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, wl)
}

func (h *WLHubHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	wlID := c.Param("id")

	var wl model.WhiteLabel
	if err := h.DB.Where("id = ? AND user_id = ?", wlID, userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	if err := c.ShouldBindJSON(&wl); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.DB.Save(&wl).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, wl)
}

// === Widget CRUD ===

func (h *WLHubHandler) ListWidgets(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var widgets []model.WidgetConfig
	h.DB.Where("white_label_id = ?", wl.ID).Find(&widgets)
	response.OK(c, widgets)
}

func (h *WLHubHandler) UpdateWidget(c *gin.Context) {
	userID := c.GetString("user_id")
	widgetID := c.Param("id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var widget model.WidgetConfig
	if err := h.DB.Where("id = ? AND white_label_id = ?", widgetID, wl.ID).First(&widget).Error; err != nil {
		response.NotFound(c, "widget not found")
		return
	}

	if err := c.ShouldBindJSON(&widget); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	h.DB.Save(&widget)
	response.OK(c, widget)
}

// === Page Builder ===

func (h *WLHubHandler) ListPages(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var pages []model.PageBuilderPage
	h.DB.Where("white_label_id = ?", wl.ID).Find(&pages)
	response.OK(c, pages)
}

func (h *WLHubHandler) CreatePage(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	page := model.PageBuilderPage{
		WhiteLabelID: wl.ID,
		Name:         req.Name,
		Slug:         req.Slug,
		Layout:       "[]",
	}

	if err := h.DB.Create(&page).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, page)
}

func (h *WLHubHandler) UpdatePage(c *gin.Context) {
	userID := c.GetString("user_id")
	pageID := c.Param("id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var page model.PageBuilderPage
	if err := h.DB.Where("id = ? AND white_label_id = ?", pageID, wl.ID).First(&page).Error; err != nil {
		response.NotFound(c, "page not found")
		return
	}

	if err := c.ShouldBindJSON(&page); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	h.DB.Save(&page)
	response.OK(c, page)
}

// === Smart Rewards ===

func (h *WLHubHandler) ListRewardRules(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var rules []model.RewardRule
	h.DB.Where("community_id = ?", community.ID).Find(&rules)
	response.OK(c, rules)
}

func (h *WLHubHandler) CreateRewardRule(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var rule model.RewardRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	rule.CommunityID = community.ID

	if err := h.DB.Create(&rule).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, rule)
}

func (h *WLHubHandler) ListPrivileges(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var privileges []model.Privilege
	h.DB.Where("community_id = ?", community.ID).Find(&privileges)
	response.OK(c, privileges)
}

func (h *WLHubHandler) CreatePrivilege(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var priv model.Privilege
	if err := c.ShouldBindJSON(&priv); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	priv.CommunityID = community.ID

	if err := h.DB.Create(&priv).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, priv)
}

// === Domain Setup ===

func (h *WLHubHandler) GetDomain(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var domain model.DomainSetup
	if err := h.DB.Where("white_label_id = ?", wl.ID).First(&domain).Error; err != nil {
		response.OK(c, gin.H{"configured": false})
		return
	}

	response.OK(c, domain)
}

func (h *WLHubHandler) SetupDomain(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domain := model.DomainSetup{
		WhiteLabelID: wl.ID,
		Domain:       req.Domain,
		DNSStatus:    "pending",
		SSLStatus:    "pending",
		CNAMETarget:  "cname.taskon.xyz",
	}

	// Upsert
	h.DB.Where("white_label_id = ?", wl.ID).Assign(domain).FirstOrCreate(&domain)

	response.OK(c, domain)
}

// === SDK Config ===

func (h *WLHubHandler) GetSDKConfig(c *gin.Context) {
	userID := c.GetString("user_id")

	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}

	var sdk model.SDKConfig
	if err := h.DB.Where("white_label_id = ?", wl.ID).First(&sdk).Error; err != nil {
		response.OK(c, gin.H{"configured": false})
		return
	}

	response.OK(c, sdk)
}
