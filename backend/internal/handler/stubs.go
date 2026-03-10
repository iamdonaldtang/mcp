package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"gorm.io/gorm"
)

// StubHandler provides mock responses for endpoints not yet fully implemented.
// This allows the frontend to render all pages without silent 404 failures.
type StubHandler struct {
	DB *gorm.DB
}

// ─── Dashboard & Settings ───

func (h *StubHandler) Dashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	hasCommunity := h.DB.Where("user_id = ?", userID).First(&community).Error == nil

	var wl model.WhiteLabel
	hasWL := h.DB.Where("user_id = ?", userID).First(&wl).Error == nil

	state := "new"
	if hasCommunity {
		state = "active"
	}

	response.OK(c, gin.H{
		"state":             state,
		"has_community":     hasCommunity,
		"has_whitelabel":    hasWL,
		"community_name":    community.Name,
		"community_status":  community.Status,
		"wl_status":         wl.Status,
		"total_members":     0,
		"active_members":    0,
		"total_tasks":       0,
		"tasks_completed":   0,
		"quick_actions":     []gin.H{},
		"recent_activity":   []gin.H{},
	})
}

func (h *StubHandler) Settings(c *gin.Context) {
	response.OK(c, gin.H{
		"organization": gin.H{
			"name":   "TaskOn Project",
			"plan":   "pro",
			"email":  "admin@taskon.xyz",
		},
		"team_members":  []gin.H{},
		"api_keys":      []gin.H{},
		"notifications": gin.H{"email": true, "slack": false},
	})
}

func (h *StubHandler) CreateAPIKey(c *gin.Context) {
	response.Created(c, gin.H{
		"id":         uuid.New().String(),
		"name":       "New API Key",
		"key":        "tk_" + uuid.New().String()[:16],
		"created_at": time.Now(),
	})
}

func (h *StubHandler) InviteTeamMember(c *gin.Context) {
	response.OK(c, gin.H{"message": "invitation sent"})
}

func (h *StubHandler) SubscriptionStatus(c *gin.Context) {
	response.OK(c, gin.H{
		"plan":         "pro",
		"status":       "active",
		"next_billing": time.Now().AddDate(0, 1, 0),
		"features":     []string{"community", "whitelabel", "analytics"},
	})
}

// ─── Community Modules ───

func (h *StubHandler) CommunityModules(c *gin.Context) {
	response.OK(c, gin.H{
		"modules": []gin.H{
			{"id": "tasks", "name": "Sectors & Tasks", "enabled": true, "status": "active"},
			{"id": "points", "name": "Points & Level", "enabled": true, "status": "active"},
			{"id": "daychain", "name": "DayChain", "enabled": true, "status": "active"},
			{"id": "taskchain", "name": "TaskChain", "enabled": false, "status": "inactive"},
			{"id": "leaderboard", "name": "Leaderboard", "enabled": true, "status": "active"},
			{"id": "lb_sprint", "name": "LB Sprint", "enabled": false, "status": "inactive"},
			{"id": "milestone", "name": "Milestone", "enabled": false, "status": "inactive"},
			{"id": "shop", "name": "Benefits Shop", "enabled": false, "status": "inactive"},
			{"id": "wheel", "name": "Lucky Wheel", "enabled": false, "status": "inactive"},
			{"id": "badges", "name": "Badges", "enabled": false, "status": "inactive"},
		},
	})
}

func (h *StubHandler) CommunityModuleStatus(c *gin.Context) {
	response.OK(c, gin.H{
		"tasks":       "active",
		"points":      "active",
		"daychain":    "active",
		"taskchain":   "inactive",
		"leaderboard": "active",
		"lb_sprint":   "inactive",
		"milestone":   "inactive",
		"shop":        "inactive",
		"wheel":       "inactive",
		"badges":      "inactive",
	})
}

func (h *StubHandler) ModuleListOrStats(c *gin.Context) {
	response.OK(c, gin.H{
		"items": []gin.H{},
		"total": 0,
		"stats": gin.H{"total": 0, "active": 0, "draft": 0},
	})
}

func (h *StubHandler) ModuleCreate(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req["id"] = uuid.New().String()
	req["status"] = "draft"
	req["created_at"] = time.Now()
	response.Created(c, req)
}

func (h *StubHandler) ModuleUpdate(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req["id"] = c.Param("id")
	req["updated_at"] = time.Now()
	response.OK(c, req)
}

// ─── Points & Level ───

func (h *StubHandler) PointTypes(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.OK(c, []gin.H{})
		return
	}
	var pts []model.PointType
	h.DB.Where("community_id = ?", community.ID).Find(&pts)
	response.OK(c, pts)
}

func (h *StubHandler) PointLevels(c *gin.Context) {
	response.OK(c, []gin.H{
		{"level": 1, "name": "Newcomer", "min_points": 0},
		{"level": 2, "name": "Explorer", "min_points": 100},
		{"level": 3, "name": "Contributor", "min_points": 500},
		{"level": 4, "name": "Champion", "min_points": 2000},
		{"level": 5, "name": "Legend", "min_points": 10000},
	})
}

func (h *StubHandler) PointStats(c *gin.Context) {
	response.OK(c, gin.H{
		"total_types":    1,
		"total_issued":   0,
		"total_redeemed": 0,
	})
}

// ─── DayChain Module ───

func (h *StubHandler) DayChainConfig(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.OK(c, gin.H{"enabled": false, "target_days": 30, "daily_points": 10})
		return
	}
	var dc model.DayChainConfig
	if err := h.DB.Where("community_id = ?", community.ID).First(&dc).Error; err != nil {
		response.OK(c, gin.H{"enabled": false, "target_days": 30, "daily_points": 10})
		return
	}
	response.OK(c, dc)
}

func (h *StubHandler) DayChainUpdate(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	var req struct {
		Enabled        *bool `json:"enabled"`
		TargetDays     *int  `json:"target_days"`
		DailyPoints    *int  `json:"daily_points"`
		CatchUpEnabled *bool `json:"catch_up_enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.TargetDays != nil {
		updates["target_days"] = *req.TargetDays
	}
	if req.DailyPoints != nil {
		updates["daily_points"] = *req.DailyPoints
	}
	if req.CatchUpEnabled != nil {
		updates["catch_up_enabled"] = *req.CatchUpEnabled
	}
	h.DB.Model(&model.DayChainConfig{}).Where("community_id = ?", community.ID).Updates(updates)
	var dc model.DayChainConfig
	h.DB.Where("community_id = ?", community.ID).First(&dc)
	response.OK(c, dc)
}

func (h *StubHandler) DayChainStats(c *gin.Context) {
	response.OK(c, gin.H{
		"active_streakers": 0,
		"avg_streak":       0,
		"longest_streak":   0,
		"total_checkins":   0,
	})
}

func (h *StubHandler) DayChainDistribution(c *gin.Context) {
	response.OK(c, []gin.H{
		{"range": "1-3 days", "count": 0},
		{"range": "4-7 days", "count": 0},
		{"range": "8-14 days", "count": 0},
		{"range": "15-30 days", "count": 0},
		{"range": "30+ days", "count": 0},
	})
}

// ─── Community Settings ───

func (h *StubHandler) AccessRules(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.OK(c, []gin.H{})
		return
	}
	var rules []model.AccessRule
	h.DB.Where("community_id = ?", community.ID).Find(&rules)
	response.OK(c, rules)
}

func (h *StubHandler) CreateAccessRule(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	rule := model.AccessRule{
		CommunityID: community.ID,
		Type:        "open",
		Enabled:     true,
	}
	if t, ok := req["type"].(string); ok {
		rule.Type = t
	}
	h.DB.Create(&rule)
	response.Created(c, rule)
}

func (h *StubHandler) AccessRuleStats(c *gin.Context) {
	response.OK(c, gin.H{"total": 0, "active": 0, "token_gate": 0, "nft_gate": 0})
}

func (h *StubHandler) AccessRulePreview(c *gin.Context) {
	response.OK(c, gin.H{"eligible_count": 0, "total_checked": 0})
}

func (h *StubHandler) AccessRuleReorder(c *gin.Context) {
	response.OK(c, gin.H{"message": "reordered"})
}

func (h *StubHandler) HomepageSections(c *gin.Context) {
	response.OK(c, []gin.H{
		{"id": "hero", "type": "hero", "title": "Welcome", "sort_order": 1, "visible": true},
		{"id": "tasks", "type": "task_list", "title": "Active Tasks", "sort_order": 2, "visible": true},
		{"id": "leaderboard", "type": "leaderboard", "title": "Top Members", "sort_order": 3, "visible": true},
	})
}

func (h *StubHandler) CreateHomepageSection(c *gin.Context) {
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	req["id"] = uuid.New().String()
	response.Created(c, req)
}

func (h *StubHandler) HomepageStats(c *gin.Context) {
	response.OK(c, gin.H{"total_sections": 3, "visible_sections": 3})
}

func (h *StubHandler) HomepageWidgetInstances(c *gin.Context) {
	response.OK(c, []gin.H{})
}

func (h *StubHandler) HomepageReorder(c *gin.Context) {
	response.OK(c, gin.H{"message": "reordered"})
}

// ─── Community Operations ───

func (h *StubHandler) CommunityChecklist(c *gin.Context) {
	response.OK(c, gin.H{
		"items": []gin.H{
			{"id": "create_community", "label": "Create your community", "completed": true},
			{"id": "add_tasks", "label": "Add tasks to a sector", "completed": false},
			{"id": "setup_shop", "label": "Set up Benefits Shop", "completed": false},
			{"id": "customize_rewards", "label": "Customize module rewards", "completed": false},
			{"id": "preview", "label": "Preview as user", "completed": false},
			{"id": "share", "label": "Share your community", "completed": false},
			{"id": "first_10", "label": "Get first 10 participants", "completed": false},
		},
		"progress": 14,
	})
}

func (h *StubHandler) OnboardingProgress(c *gin.Context) {
	response.OK(c, gin.H{
		"current_step": 1,
		"total_steps":  7,
		"dismissed":    false,
		"progress":     14,
	})
}

func (h *StubHandler) OnboardingDismiss(c *gin.Context) {
	response.OK(c, gin.H{"dismissed": true})
}

func (h *StubHandler) CommunityPreview(c *gin.Context) {
	response.OK(c, gin.H{
		"url":      "https://community.taskon.xyz/preview",
		"sections": []gin.H{},
	})
}

func (h *StubHandler) CommunityPublish(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	h.DB.Model(&community).Update("status", "active")
	response.OK(c, gin.H{"status": "active", "message": "published"})
}

func (h *StubHandler) CommunityStats(c *gin.Context) {
	response.OK(c, gin.H{
		"total_members": 0, "active_members": 0,
		"total_tasks": 0, "completed_tasks": 0,
		"total_points_issued": 0,
	})
}

func (h *StubHandler) SlugCheck(c *gin.Context) {
	response.OK(c, gin.H{"available": true})
}

func (h *StubHandler) SmartRewardsOverview(c *gin.Context) {
	response.OK(c, gin.H{
		"total_rules":        0,
		"active_rules":       0,
		"total_executions":   0,
		"total_privileges":   0,
	})
}

func (h *StubHandler) ContentAnnouncements(c *gin.Context) {
	response.OK(c, []gin.H{})
}

func (h *StubHandler) CreateContentAnnouncement(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Link        string `json:"link"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	ann := model.Announcement{
		CommunityID: community.ID,
		Title:       req.Title,
		Description: req.Description,
		Image:       req.Image,
		Link:        req.Link,
		IsActive:    true,
	}
	h.DB.Create(&ann)
	response.Created(c, ann)
}

func (h *StubHandler) ContentFeatured(c *gin.Context) {
	response.OK(c, []gin.H{})
}

func (h *StubHandler) CreateContentFeatured(c *gin.Context) {
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	req["id"] = uuid.New().String()
	response.Created(c, req)
}

func (h *StubHandler) CommunityIntegrations(c *gin.Context) {
	response.OK(c, gin.H{
		"social": []gin.H{
			{"id": "twitter", "name": "Twitter/X", "connected": false},
			{"id": "discord", "name": "Discord", "connected": false},
			{"id": "telegram", "name": "Telegram", "connected": false},
		},
		"blockchain": []gin.H{
			{"id": "ethereum", "name": "Ethereum", "connected": false},
			{"id": "polygon", "name": "Polygon", "connected": false},
			{"id": "bnb", "name": "BNB Chain", "connected": false},
		},
		"analytics": []gin.H{
			{"id": "ga4", "name": "Google Analytics", "connected": false},
		},
	})
}

func (h *StubHandler) WizardDraftSave(c *gin.Context) {
	response.OK(c, gin.H{"message": "draft saved"})
}

func (h *StubHandler) WizardDraftLoad(c *gin.Context) {
	response.OK(c, gin.H{"draft": nil})
}

// ─── White Label Extras ───

func (h *StubHandler) WLBrand(c *gin.Context) {
	userID := c.GetString("user_id")
	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.OK(c, gin.H{
			"primary_color": "#F59E0B",
			"accent_color":  "#F59E0B",
			"logo":          "",
			"favicon":       "",
			"custom_css":    "",
		})
		return
	}
	response.OK(c, gin.H{
		"primary_color": wl.BrandPrimaryColor,
		"accent_color":  wl.BrandAccentColor,
		"logo":          wl.BrandLogo,
		"favicon":       wl.BrandFavicon,
		"custom_css":    wl.BrandCustomCSS,
	})
}

func (h *StubHandler) WLBrandUpdate(c *gin.Context) {
	userID := c.GetString("user_id")
	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	updates := map[string]interface{}{}
	if v, ok := req["primary_color"]; ok {
		updates["brand_primary_color"] = v
	}
	if v, ok := req["accent_color"]; ok {
		updates["brand_accent_color"] = v
	}
	if v, ok := req["custom_css"]; ok {
		updates["brand_custom_css"] = v
	}
	h.DB.Model(&wl).Updates(updates)
	response.OK(c, gin.H{"message": "brand updated"})
}

func (h *StubHandler) WLBrandLogoUpload(c *gin.Context) {
	response.OK(c, gin.H{"url": "https://placehold.co/200x200/48BB78/ffffff?text=Logo"})
}

func (h *StubHandler) WLContracts(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.OK(c, []gin.H{})
		return
	}
	var contracts []model.ContractEntry
	h.DB.Where("community_id = ?", community.ID).Find(&contracts)
	response.OK(c, contracts)
}

func (h *StubHandler) WLContractCreate(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
		Chain   string `json:"chain" binding:"required"`
		ABI     string `json:"abi"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	entry := model.ContractEntry{
		CommunityID: community.ID,
		Name:        req.Name,
		Address:     req.Address,
		Chain:       req.Chain,
		ABI:         req.ABI,
		AddedAt:     time.Now(),
	}
	h.DB.Create(&entry)
	response.Created(c, entry)
}

func (h *StubHandler) WLContractCheck(c *gin.Context) {
	response.OK(c, gin.H{"valid": true, "verified": false})
}

func (h *StubHandler) WLContractStats(c *gin.Context) {
	response.OK(c, gin.H{"total": 0, "verified": 0})
}

func (h *StubHandler) WLIntegrations(c *gin.Context) {
	response.OK(c, gin.H{
		"integrations": []gin.H{
			{"id": "sso_wallet", "name": "Wallet SSO", "status": "not_configured"},
			{"id": "sso_oauth", "name": "OAuth SSO", "status": "not_configured"},
			{"id": "webhook", "name": "Webhooks", "status": "not_configured"},
		},
	})
}

func (h *StubHandler) WLOnboarding(c *gin.Context) {
	response.OK(c, gin.H{
		"items": []gin.H{
			{"id": "create_wl", "label": "Create White Label", "completed": true},
			{"id": "configure_widgets", "label": "Configure widgets", "completed": false},
			{"id": "build_page", "label": "Build a page", "completed": false},
			{"id": "preview", "label": "Preview your site", "completed": false},
			{"id": "send_devkit", "label": "Send Dev Kit", "completed": false},
			{"id": "integration_verified", "label": "Integration verified", "completed": false},
			{"id": "announce", "label": "Announce launch", "completed": false},
			{"id": "first_interaction", "label": "First user interaction", "completed": false},
		},
		"progress": 12,
	})
}

func (h *StubHandler) WLPublish(c *gin.Context) {
	userID := c.GetString("user_id")
	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}
	h.DB.Model(&wl).Update("status", "active")
	response.OK(c, gin.H{"status": "active", "message": "published"})
}

func (h *StubHandler) WLReadiness(c *gin.Context) {
	response.OK(c, gin.H{
		"ready": false,
		"checks": []gin.H{
			{"id": "subscription", "label": "Active subscription", "passed": true},
			{"id": "twitter_auth", "label": "Twitter authenticated", "passed": false},
		},
	})
}

func (h *StubHandler) WLRules(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.OK(c, []gin.H{})
		return
	}
	var rules []model.RewardRule
	h.DB.Where("community_id = ?", community.ID).Find(&rules)
	response.OK(c, rules)
}

func (h *StubHandler) WLRuleCreate(c *gin.Context) {
	userID := c.GetString("user_id")
	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}
	var req struct {
		Name        string `json:"name" binding:"required"`
		TriggerType string `json:"trigger_type" binding:"required"`
		RewardType  string `json:"reward_type" binding:"required"`
		RewardValue string `json:"reward_value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	rule := model.RewardRule{
		CommunityID: community.ID,
		Name:        req.Name,
		TriggerType: req.TriggerType,
		RewardType:  req.RewardType,
		RewardValue: req.RewardValue,
		IsActive:    true,
	}
	h.DB.Create(&rule)
	response.Created(c, rule)
}

func (h *StubHandler) WLAntiSybilUpdate(c *gin.Context) {
	response.OK(c, gin.H{"message": "anti-sybil rules updated"})
}

func (h *StubHandler) WLPrivilegeStats(c *gin.Context) {
	response.OK(c, gin.H{"total": 0, "active": 0, "total_claims": 0})
}

func (h *StubHandler) WLSDKGenerateKey(c *gin.Context) {
	response.Created(c, gin.H{
		"api_key":    "tk_" + uuid.New().String()[:16],
		"api_secret": "tks_" + uuid.New().String(),
	})
}

func (h *StubHandler) WLSDKKeys(c *gin.Context) {
	response.Created(c, gin.H{"message": "key created"})
}

func (h *StubHandler) WLSDKWebhooks(c *gin.Context) {
	response.OK(c, gin.H{"message": "webhook saved"})
}

func (h *StubHandler) WLSSOTest(c *gin.Context) {
	response.OK(c, gin.H{"success": true, "message": "SSO test passed"})
}

func (h *StubHandler) WLDevKitSend(c *gin.Context) {
	response.OK(c, gin.H{"message": "dev kit sent"})
}

func (h *StubHandler) WLDomainVerify(c *gin.Context) {
	response.OK(c, gin.H{"status": "pending", "message": "DNS verification initiated"})
}

func (h *StubHandler) WLWidgetCreate(c *gin.Context) {
	userID := c.GetString("user_id")
	var wl model.WhiteLabel
	if err := h.DB.Where("user_id = ?", userID).First(&wl).Error; err != nil {
		response.NotFound(c, "white label not found")
		return
	}
	var req struct {
		ModuleType string `json:"module_type" binding:"required"`
		ModuleName string `json:"module_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	widget := model.WidgetConfig{
		WhiteLabelID: wl.ID,
		ModuleType:   req.ModuleType,
		ModuleName:   req.ModuleName,
		IsConfigured: false,
		IsActive:     false,
		Settings:     "{}",
	}
	h.DB.Create(&widget)
	response.Created(c, widget)
}

func (h *StubHandler) WLWizardDraft(c *gin.Context) {
	response.OK(c, gin.H{"draft": nil})
}

// ─── C-End Extras ───

func (h *StubHandler) CEndUserProfile(c *gin.Context) {
	wallet := c.GetString("wallet_address")
	response.OK(c, gin.H{
		"wallet_address": wallet,
		"username":       "",
		"avatar":         "",
		"bio":            "",
	})
}

func (h *StubHandler) CEndAchievements(c *gin.Context) {
	response.OK(c, []gin.H{})
}

func (h *StubHandler) CEndReferralStats(c *gin.Context) {
	response.OK(c, gin.H{
		"total_referrals":  0,
		"active_referrals": 0,
		"total_earned":     0,
		"referral_code":    uuid.New().String()[:8],
	})
}
