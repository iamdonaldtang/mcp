package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"gorm.io/gorm"
)

type CommunityHubHandler struct {
	DB *gorm.DB
}

func (h *CommunityHubHandler) GetOverview(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		// No community yet — return empty state indicator
		response.OK(c, gin.H{"state": "empty"})
		return
	}

	// Count members
	var memberCount int64
	h.DB.Model(&model.CommunityMember{}).Where("community_id = ?", community.ID).Count(&memberCount)

	// Count active members (checked in within 7 days)
	var activeMemberCount int64
	h.DB.Model(&model.CommunityMember{}).
		Where("community_id = ? AND updated_at > NOW() - INTERVAL '7 days'", community.ID).
		Count(&activeMemberCount)

	// Count tasks
	var taskCount int64
	h.DB.Model(&model.Task{}).Where("community_id = ?", community.ID).Count(&taskCount)

	response.OK(c, gin.H{
		"id":              community.ID,
		"name":            community.Name,
		"description":     community.Description,
		"status":          community.Status,
		"brand_color":     community.BrandColor,
		"logo":            community.Logo,
		"enabled_modules": community.EnabledModules,
		"total_members":   memberCount,
		"active_members":  activeMemberCount,
		"total_tasks":     taskCount,
		"created_at":      community.CreatedAt,
	})
}

func (h *CommunityHubHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		BrandColor  string   `json:"brand_color"`
		Logo        string   `json:"logo"`
		Modules     []string `json:"modules"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	community := model.Community{
		UserID:         userID,
		Name:           req.Name,
		Description:    req.Description,
		BrandColor:     req.BrandColor,
		Logo:           req.Logo,
		EnabledModules: req.Modules,
		Status:         "draft",
	}

	if err := h.DB.Create(&community).Error; err != nil {
		response.InternalError(c)
		return
	}

	// Create default point type
	defaultPT := model.PointType{
		CommunityID: community.ID,
		Name:        "Community Points",
		Symbol:      "XP",
		Icon:        "stars",
		Color:       "#F59E0B",
		IsDefault:   true,
	}
	h.DB.Create(&defaultPT)

	// Create default DayChain config
	dc := model.DayChainConfig{
		CommunityID: community.ID,
		Enabled:     false,
		TargetDays:  30,
		DailyPoints: 10,
	}
	h.DB.Create(&dc)

	response.Created(c, community)
}

func (h *CommunityHubHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	communityID := c.Param("id")

	var community model.Community
	if err := h.DB.Where("id = ? AND user_id = ?", communityID, userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var req struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		BrandColor  *string  `json:"brand_color"`
		Logo        *string  `json:"logo"`
		Status      *string  `json:"status"`
		Modules     []string `json:"modules"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.BrandColor != nil {
		updates["brand_color"] = *req.BrandColor
	}
	if req.Logo != nil {
		updates["logo"] = *req.Logo
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Modules != nil {
		updates["enabled_modules"] = model.StringArray(req.Modules)
	}

	if err := h.DB.Model(&community).Updates(updates).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, community)
}

// === Sector CRUD ===

func (h *CommunityHubHandler) ListSectors(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var sectors []model.Sector
	h.DB.Where("community_id = ?", community.ID).Order("sort_order ASC").Preload("Tasks").Find(&sectors)
	response.OK(c, sectors)
}

func (h *CommunityHubHandler) CreateSector(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get next sort order
	var maxOrder int
	h.DB.Model(&model.Sector{}).Where("community_id = ?", community.ID).
		Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	sector := model.Sector{
		CommunityID: community.ID,
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   maxOrder + 1,
		Status:      "active",
	}

	if err := h.DB.Create(&sector).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, sector)
}

// === Task CRUD ===

func (h *CommunityHubHandler) ListTasks(c *gin.Context) {
	userID := c.GetString("user_id")
	sectorID := c.Query("sector_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	query := h.DB.Model(&model.Task{}).Where("community_id = ?", community.ID)
	if sectorID != "" {
		query = query.Where("sector_id = ?", sectorID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tasks []model.Task
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&tasks)

	response.Paginated(c, tasks, total, page, pageSize)
}

func (h *CommunityHubHandler) CreateTask(c *gin.Context) {
	userID := c.GetString("user_id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var req struct {
		SectorID       string `json:"sector_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description"`
		Type           string `json:"type" binding:"required"`
		Points         int    `json:"points"`
		Icon           string `json:"icon"`
		IconColor      string `json:"icon_color"`
		MaxCompletions *int   `json:"max_completions"`
		CooldownHours  *int   `json:"cooldown_hours"`
		Requirement    string `json:"requirement"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task := model.Task{
		SectorID:       req.SectorID,
		CommunityID:    community.ID,
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		Points:         req.Points,
		Icon:           req.Icon,
		IconColor:      req.IconColor,
		MaxCompletions: req.MaxCompletions,
		CooldownHours:  req.CooldownHours,
		Requirement:    req.Requirement,
		Status:         "draft",
	}

	if task.Icon == "" {
		task.Icon = "task_alt"
	}
	if task.IconColor == "" {
		task.IconColor = "#F59E0B"
	}

	if err := h.DB.Create(&task).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, task)
}

func (h *CommunityHubHandler) UpdateTask(c *gin.Context) {
	userID := c.GetString("user_id")
	taskID := c.Param("id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	var task model.Task
	if err := h.DB.Where("id = ? AND community_id = ?", taskID, community.ID).First(&task).Error; err != nil {
		response.NotFound(c, "task not found")
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.DB.Save(&task).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, task)
}

func (h *CommunityHubHandler) DeleteTask(c *gin.Context) {
	userID := c.GetString("user_id")
	taskID := c.Param("id")

	var community model.Community
	if err := h.DB.Where("user_id = ?", userID).First(&community).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	if err := h.DB.Where("id = ? AND community_id = ?", taskID, community.ID).Delete(&model.Task{}).Error; err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, nil)
}
