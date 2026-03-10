package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/pkg/response"
	"gorm.io/gorm"
)

type CCommunityHandler struct {
	DB *gorm.DB
}

// Home returns aggregate home page data for C-end
func (h *CCommunityHandler) Home(c *gin.Context) {
	communityID := c.GetString("community_id")
	if communityID == "" {
		communityID = c.Query("community_id")
	}

	var community model.Community
	if err := h.DB.First(&community, "id = ?", communityID).Error; err != nil {
		response.NotFound(c, "community not found")
		return
	}

	// Pulse stats
	var totalMembers, thisWeek, liveActive, tasksDone int64
	h.DB.Model(&model.CommunityMember{}).Where("community_id = ?", communityID).Count(&totalMembers)
	h.DB.Model(&model.CommunityMember{}).Where("community_id = ? AND joined_at > NOW() - INTERVAL '7 days'", communityID).Count(&thisWeek)
	h.DB.Model(&model.CommunityMember{}).Where("community_id = ? AND updated_at > NOW() - INTERVAL '15 minutes'", communityID).Count(&liveActive)
	h.DB.Model(&model.TaskCompletion{}).Where("community_id = ?", communityID).Count(&tasksDone)

	// Tab visibility
	var hasQuests, hasLeaderboard, hasSprint, hasMilestone, hasShop bool
	h.DB.Model(&model.Task{}).Where("community_id = ? AND status = ?", communityID, "active").
		Select("COUNT(*) > 0").Scan(&hasQuests)
	h.DB.Model(&model.LeaderboardConfig{}).Where("community_id = ? AND is_enabled = ?", communityID, true).
		Select("COUNT(*) > 0").Scan(&hasLeaderboard)
	h.DB.Model(&model.LBSprint{}).Where("community_id = ? AND (status = ? OR status = ?)", communityID, "active", "completed").
		Select("COUNT(*) > 0").Scan(&hasSprint)
	h.DB.Model(&model.Milestone{}).Where("community_id = ? AND status = ?", communityID, "active").
		Select("COUNT(*) > 0").Scan(&hasMilestone)
	h.DB.Model(&model.ShopItem{}).Where("community_id = ? AND status = ?", communityID, "active").
		Select("COUNT(*) > 0").Scan(&hasShop)

	response.OK(c, gin.H{
		"community": gin.H{
			"name":        community.Name,
			"brand_color": community.BrandColor,
			"logo":        community.Logo,
			"status":      community.Status,
		},
		"pulse": gin.H{
			"total_members": totalMembers,
			"this_week":     thisWeek,
			"live_active":   liveActive,
			"tasks_done":    tasksDone,
		},
		"tabVisibility": gin.H{
			"home":        true,
			"quests":      hasQuests,
			"leaderboard": hasLeaderboard,
			"lbSprint":    hasSprint,
			"milestone":   hasMilestone,
			"shop":        hasShop,
		},
	})
}

// Tasks returns sectors with their tasks for C-end display
func (h *CCommunityHandler) Tasks(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var sectors []model.Sector
	h.DB.Where("community_id = ? AND status = ?", communityID, "active").
		Order("sort_order ASC").
		Preload("Tasks", "status = ?", "active").
		Find(&sectors)

	// If user is logged in, annotate task statuses with their completions
	if walletAddr != "" {
		var completions []model.TaskCompletion
		h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).Find(&completions)

		completionMap := make(map[string]string)
		for _, tc := range completions {
			completionMap[tc.TaskID] = tc.Status
		}

		for i := range sectors {
			for j := range sectors[i].Tasks {
				if status, ok := completionMap[sectors[i].Tasks[j].ID]; ok {
					sectors[i].Tasks[j].Status = status
				}
			}
		}
	}

	response.OK(c, sectors)
}

// Announcements returns active announcements
func (h *CCommunityHandler) Announcements(c *gin.Context) {
	communityID := c.GetString("community_id")

	var announcements []model.Announcement
	h.DB.Where("community_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)",
		communityID, true, time.Now()).
		Order("created_at DESC").
		Find(&announcements)

	response.OK(c, announcements)
}

// DayChainStatus returns current DayChain state for a user
func (h *CCommunityHandler) DayChainStatus(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var config model.DayChainConfig
	if err := h.DB.Where("community_id = ?", communityID).First(&config).Error; err != nil {
		response.OK(c, gin.H{"enabled": false})
		return
	}

	if !config.Enabled {
		response.OK(c, gin.H{"enabled": false})
		return
	}

	var member model.CommunityMember
	if err := h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member).Error; err != nil {
		response.NotFound(c, "member not found")
		return
	}

	checkedInToday := false
	if member.LastCheckIn != nil {
		checkedInToday = member.LastCheckIn.Format("2006-01-02") == time.Now().Format("2006-01-02")
	}

	response.OK(c, gin.H{
		"enabled":          true,
		"current_streak":   member.DayStreak,
		"target_days":      config.TargetDays,
		"checked_in_today": checkedInToday,
		"daily_points":     config.DailyPoints,
	})
}

// DayChainCheckIn performs daily check-in
func (h *CCommunityHandler) DayChainCheckIn(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var member model.CommunityMember
	if err := h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member).Error; err != nil {
		response.NotFound(c, "member not found")
		return
	}

	// Check if already checked in today
	if member.LastCheckIn != nil && member.LastCheckIn.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		response.BadRequest(c, "already checked in today")
		return
	}

	var config model.DayChainConfig
	h.DB.Where("community_id = ?", communityID).First(&config)

	now := time.Now()

	// Check if streak continues or resets
	if member.LastCheckIn != nil {
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		if member.LastCheckIn.Format("2006-01-02") == yesterday {
			member.DayStreak++
		} else {
			member.DayStreak = 1 // streak broken
		}
	} else {
		member.DayStreak = 1
	}

	member.LastCheckIn = &now
	member.TotalPoints += config.DailyPoints
	member.LifetimeXP += config.DailyPoints

	h.DB.Save(&member)

	// Log activity
	pts := config.DailyPoints
	h.DB.Create(&model.ActivityLog{
		CommunityID:   communityID,
		WalletAddress: walletAddr,
		Type:          "streak",
		Title:         "Daily Check-In",
		PointsDelta:   &pts,
	})

	response.OK(c, gin.H{
		"streak":       member.DayStreak,
		"points_earned": config.DailyPoints,
		"total_points":  member.TotalPoints,
	})
}

// UserStatus returns the current user's status
func (h *CCommunityHandler) UserStatus(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var member model.CommunityMember
	if err := h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member).Error; err != nil {
		response.NotFound(c, "member not found")
		return
	}

	// Compute rank
	var rank int64
	h.DB.Model(&model.CommunityMember{}).
		Where("community_id = ? AND total_points > ?", communityID, member.TotalPoints).
		Count(&rank)

	response.OK(c, gin.H{
		"wallet_address": member.WalletAddress,
		"level":          member.Level,
		"xp":             member.LifetimeXP,
		"total_points":   member.TotalPoints,
		"day_streak":     member.DayStreak,
		"rank":           rank + 1,
		"joined_at":      member.JoinedAt,
	})
}

// Leaderboard returns rankings
func (h *CCommunityHandler) Leaderboard(c *gin.Context) {
	communityID := c.GetString("community_id")
	// period := c.DefaultQuery("period", "alltime") // TODO: filter by period

	var members []model.CommunityMember
	h.DB.Where("community_id = ?", communityID).
		Order("total_points DESC").
		Limit(100).
		Find(&members)

	rankings := make([]gin.H, len(members))
	for i, m := range members {
		rankings[i] = gin.H{
			"rank":           i + 1,
			"wallet_address": m.WalletAddress,
			"level":          m.Level,
			"points":         m.TotalPoints,
		}
	}

	response.OK(c, gin.H{
		"rankings": rankings,
	})
}

// Milestones returns milestones with user claim status
func (h *CCommunityHandler) Milestones(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var milestones []model.Milestone
	h.DB.Where("community_id = ? AND status = ?", communityID, "active").
		Order("threshold ASC").
		Find(&milestones)

	// Get user's claims
	var claims []model.MilestoneClaim
	if walletAddr != "" {
		h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).Find(&claims)
	}
	claimedMap := make(map[string]bool)
	for _, cl := range claims {
		claimedMap[cl.MilestoneID] = true
	}

	// Get user's points
	var member model.CommunityMember
	h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member)

	result := make([]gin.H, len(milestones))
	for i, m := range milestones {
		status := "locked"
		if claimedMap[m.ID] {
			status = "earned"
		} else if member.TotalPoints >= m.Threshold {
			status = "claimable"
		}
		result[i] = gin.H{
			"id":          m.ID,
			"name":        m.Name,
			"description": m.Description,
			"requirement": m.Requirement,
			"threshold":   m.Threshold,
			"reward_type": m.RewardType,
			"reward":      m.RewardValue,
			"status":      status,
		}
	}

	response.OK(c, result)
}

// ClaimMilestone claims a milestone reward
func (h *CCommunityHandler) ClaimMilestone(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")
	milestoneID := c.Param("id")

	var milestone model.Milestone
	if err := h.DB.Where("id = ? AND community_id = ?", milestoneID, communityID).First(&milestone).Error; err != nil {
		response.NotFound(c, "milestone not found")
		return
	}

	// Check not already claimed
	var existingClaim model.MilestoneClaim
	if err := h.DB.Where("milestone_id = ? AND wallet_address = ?", milestoneID, walletAddr).First(&existingClaim).Error; err == nil {
		response.BadRequest(c, "already claimed")
		return
	}

	// Check user meets threshold
	var member model.CommunityMember
	if err := h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member).Error; err != nil {
		response.NotFound(c, "member not found")
		return
	}
	if member.TotalPoints < milestone.Threshold {
		response.BadRequest(c, "insufficient points")
		return
	}

	// Create claim
	claim := model.MilestoneClaim{
		CommunityID:   communityID,
		MilestoneID:   milestoneID,
		WalletAddress: walletAddr,
		ClaimedAt:     time.Now(),
	}
	h.DB.Create(&claim)

	// Increment milestone claims
	h.DB.Model(&milestone).UpdateColumn("claims", gorm.Expr("claims + 1"))

	response.OK(c, gin.H{"claimed": true})
}

// ShopItems returns shop items
func (h *CCommunityHandler) ShopItems(c *gin.Context) {
	communityID := c.GetString("community_id")

	var items []model.ShopItem
	h.DB.Where("community_id = ? AND status = ?", communityID, "active").
		Order("created_at DESC").
		Find(&items)

	response.OK(c, items)
}

// RedeemShopItem redeems a shop item
func (h *CCommunityHandler) RedeemShopItem(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var req struct {
		ItemID string `json:"item_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var item model.ShopItem
	if err := h.DB.Where("id = ? AND community_id = ?", req.ItemID, communityID).First(&item).Error; err != nil {
		response.NotFound(c, "item not found")
		return
	}

	if item.Stock == 0 {
		response.BadRequest(c, "sold out")
		return
	}

	var member model.CommunityMember
	if err := h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).First(&member).Error; err != nil {
		response.NotFound(c, "member not found")
		return
	}

	if member.TotalPoints < item.Price {
		response.BadRequest(c, "insufficient points")
		return
	}

	// Deduct points
	member.TotalPoints -= item.Price
	h.DB.Save(&member)

	// Update stock
	if item.Stock > 0 {
		h.DB.Model(&item).UpdateColumn("stock", gorm.Expr("stock - 1"))
	}
	h.DB.Model(&item).UpdateColumn("total_redemptions", gorm.Expr("total_redemptions + 1"))

	// Create redemption record
	h.DB.Create(&model.ShopRedemption{
		CommunityID:   communityID,
		ShopItemID:    req.ItemID,
		WalletAddress: walletAddr,
		PointsSpent:   item.Price,
		RedeemedAt:    time.Now(),
	})

	// Log activity
	pts := -item.Price
	h.DB.Create(&model.ActivityLog{
		CommunityID:   communityID,
		WalletAddress: walletAddr,
		Type:          "shop",
		Title:         "Redeemed " + item.Name,
		PointsDelta:   &pts,
	})

	response.OK(c, gin.H{"redeemed": true, "remaining_points": member.TotalPoints})
}

// ActivityFeed returns user's activity
func (h *CCommunityHandler) ActivityFeed(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	var activities []model.ActivityLog
	h.DB.Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).
		Order("created_at DESC").
		Limit(50).
		Find(&activities)

	response.OK(c, activities)
}

// InviteLink returns or generates referral link
func (h *CCommunityHandler) InviteLink(c *gin.Context) {
	communityID := c.GetString("community_id")
	walletAddr := c.GetString("wallet_address")

	// Find existing referral or generate code
	var existing model.Referral
	err := h.DB.Where("community_id = ? AND referrer_address = ?", communityID, walletAddr).First(&existing).Error

	code := ""
	if err == nil {
		code = existing.ReferralCode
	} else {
		code = walletAddr[:8] // simple code generation
	}

	// Count stats
	var totalInvites, successfulJoins int64
	var pointsEarned int64
	h.DB.Model(&model.Referral{}).Where("community_id = ? AND referrer_address = ?", communityID, walletAddr).Count(&totalInvites)
	h.DB.Model(&model.Referral{}).Where("community_id = ? AND referrer_address = ? AND status = ?", communityID, walletAddr, "joined").Count(&successfulJoins)
	h.DB.Model(&model.Referral{}).Where("community_id = ? AND referrer_address = ?", communityID, walletAddr).
		Select("COALESCE(SUM(points_earned), 0)").Scan(&pointsEarned)

	convRate := float64(0)
	if totalInvites > 0 {
		convRate = float64(successfulJoins) / float64(totalInvites) * 100
	}

	response.OK(c, gin.H{
		"referral_code":    code,
		"referral_url":     "https://community.example.com/ref=" + code,
		"total_invites":    totalInvites,
		"successful_joins": successfulJoins,
		"points_earned":    pointsEarned,
		"conversion_rate":  convRate,
	})
}

// LBSprintCurrent returns the current active sprint
func (h *CCommunityHandler) LBSprintCurrent(c *gin.Context) {
	communityID := c.GetString("community_id")

	var sprint model.LBSprint
	if err := h.DB.Where("community_id = ? AND status = ?", communityID, "active").
		Preload("Rewards").
		First(&sprint).Error; err != nil {
		response.OK(c, gin.H{"active": false})
		return
	}

	response.OK(c, gin.H{
		"active": true,
		"sprint": sprint,
	})
}

// LBSprintHistory returns past sprints
func (h *CCommunityHandler) LBSprintHistory(c *gin.Context) {
	communityID := c.GetString("community_id")

	var sprints []model.LBSprint
	h.DB.Where("community_id = ? AND status IN ?", communityID, []string{"completed", "ended"}).
		Order("ends_at DESC").
		Find(&sprints)

	response.OK(c, sprints)
}

// Quests returns available quests (sectors grouped as quests)
func (h *CCommunityHandler) Quests(c *gin.Context) {
	communityID := c.GetString("community_id")

	var sectors []model.Sector
	h.DB.Where("community_id = ? AND status = ?", communityID, "active").
		Preload("Tasks", "status = ?", "active").
		Find(&sectors)

	// Map sectors as "quests" for C-end display
	quests := make([]gin.H, len(sectors))
	for i, s := range sectors {
		quests[i] = gin.H{
			"id":          s.ID,
			"name":        s.Name,
			"description": s.Description,
			"task_count":  len(s.Tasks),
			"status":      "available",
		}
	}

	response.OK(c, quests)
}
