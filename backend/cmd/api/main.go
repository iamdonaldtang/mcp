package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taskon/backend/internal/config"
	"github.com/taskon/backend/internal/handler"
	"github.com/taskon/backend/internal/middleware"
	"github.com/taskon/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Handlers
	authH := &handler.AuthHandler{DB: db, Secret: cfg.JWTSecret}
	walletH := &handler.WalletHandler{DB: db, Secret: cfg.JWTSecret}
	communityH := &handler.CommunityHubHandler{DB: db}
	wlH := &handler.WLHubHandler{DB: db}
	cCommunityH := &handler.CCommunityHandler{DB: db}
	stubH := &handler.StubHandler{DB: db}

	// Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// === B-End API (/api/v1/) ===
	v1 := r.Group("/api/v1")
	{
		// Auth (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authH.Login)
			auth.POST("/register", authH.Register)
			auth.GET("/profile", middleware.BEndAuth(cfg.JWTSecret), authH.Profile)
		}

		// Protected B-end routes
		b := v1.Group("")
		b.Use(middleware.BEndAuth(cfg.JWTSecret))
		{
			// Dashboard & Settings
			b.GET("/dashboard", stubH.Dashboard)
			b.POST("/dashboard", stubH.Dashboard)
			b.GET("/settings", stubH.Settings)
			b.POST("/settings/api-keys", stubH.CreateAPIKey)
			b.POST("/settings/organization/invite", stubH.InviteTeamMember)
			b.GET("/subscription/status", stubH.SubscriptionStatus)

			// Community Hub
			community := b.Group("/community")
			{
				community.GET("/overview", communityH.GetOverview)
				community.POST("", communityH.Create)
				community.PUT("/:id", communityH.Update)

				// Sectors
				community.GET("/sectors", communityH.ListSectors)
				community.POST("/sectors", communityH.CreateSector)

				// Tasks
				community.GET("/tasks", communityH.ListTasks)
				community.POST("/tasks", communityH.CreateTask)
				community.PUT("/tasks/:id", communityH.UpdateTask)
				community.DELETE("/tasks/:id", communityH.DeleteTask)

				// Modules
				community.GET("/modules", stubH.CommunityModules)
				community.GET("/modules/status", stubH.CommunityModuleStatus)

				// Points & Level
				community.GET("/modules/points/types", stubH.PointTypes)
				community.GET("/modules/points/levels", stubH.PointLevels)
				community.GET("/modules/points/stats", stubH.PointStats)

				// DayChain
				community.GET("/modules/daychain", stubH.DayChainConfig)
				community.POST("/modules/daychain", stubH.DayChainUpdate)
				community.GET("/modules/daychain/stats", stubH.DayChainStats)
				community.GET("/modules/daychain/streak-distribution", stubH.DayChainDistribution)

				// TaskChain
				community.GET("/modules/taskchain", stubH.ModuleListOrStats)
				community.POST("/modules/taskchain", stubH.ModuleCreate)
				community.GET("/modules/taskchain/stats", stubH.ModuleListOrStats)
				community.GET("/modules/taskchain/funnel", stubH.ModuleListOrStats)

				// Leaderboard
				community.GET("/modules/leaderboard", stubH.ModuleListOrStats)
				community.POST("/modules/leaderboard", stubH.ModuleCreate)
				community.GET("/modules/leaderboard/stats", stubH.ModuleListOrStats)

				// LB Sprint
				community.GET("/modules/lb-sprint", stubH.ModuleListOrStats)
				community.POST("/modules/lb-sprint", stubH.ModuleCreate)
				community.GET("/modules/lb-sprint/stats", stubH.ModuleListOrStats)

				// Milestone
				community.GET("/modules/milestone", stubH.ModuleListOrStats)
				community.POST("/modules/milestone", stubH.ModuleCreate)
				community.GET("/modules/milestone/stats", stubH.ModuleListOrStats)

				// Benefits Shop
				community.GET("/modules/shop", stubH.ModuleListOrStats)
				community.GET("/modules/shop/stats", stubH.ModuleListOrStats)

				// Lucky Wheel
				community.GET("/modules/wheel", stubH.ModuleListOrStats)
				community.POST("/modules/wheel", stubH.ModuleCreate)
				community.GET("/modules/wheel/stats", stubH.ModuleListOrStats)

				// Badges
				community.GET("/modules/badges", stubH.ModuleListOrStats)
				community.POST("/modules/badges", stubH.ModuleCreate)
				community.GET("/modules/badges/stats", stubH.ModuleListOrStats)

				// Settings
				community.GET("/settings/access-rules", stubH.AccessRules)
				community.POST("/settings/access-rules", stubH.CreateAccessRule)
				community.GET("/settings/access-rules/stats", stubH.AccessRuleStats)
				community.POST("/settings/access-rules/preview", stubH.AccessRulePreview)
				community.PUT("/settings/access-rules/reorder", stubH.AccessRuleReorder)
				community.GET("/settings/homepage/sections", stubH.HomepageSections)
				community.POST("/settings/homepage/sections", stubH.CreateHomepageSection)
				community.GET("/settings/homepage/stats", stubH.HomepageStats)
				community.GET("/settings/homepage/widget-instances", stubH.HomepageWidgetInstances)
				community.PUT("/settings/homepage/reorder", stubH.HomepageReorder)

				// Operations
				community.GET("/checklist", stubH.CommunityChecklist)
				community.GET("/onboarding/progress", stubH.OnboardingProgress)
				community.PUT("/onboarding/dismiss", stubH.OnboardingDismiss)
				community.GET("/preview", stubH.CommunityPreview)
				community.POST("/publish", stubH.CommunityPublish)
				community.GET("/stats", stubH.CommunityStats)
				community.GET("/slug/check", stubH.SlugCheck)
				community.GET("/smart-rewards/overview", stubH.SmartRewardsOverview)
				community.GET("/content/announcements", stubH.ContentAnnouncements)
				community.POST("/content/announcements", stubH.CreateContentAnnouncement)
				community.GET("/content/featured", stubH.ContentFeatured)
				community.POST("/content/featured", stubH.CreateContentFeatured)
				community.GET("/integrations", stubH.CommunityIntegrations)
				community.POST("/wizard/draft", stubH.WizardDraftSave)
				community.PUT("/wizard/draft", stubH.WizardDraftSave)
				community.GET("/wizard/draft", stubH.WizardDraftLoad)
			}

			// White Label Hub
			wl := b.Group("/whitelabel")
			{
				wl.GET("/overview", wlH.GetOverview)
				wl.POST("", wlH.Create)
				wl.PUT("/:id", wlH.Update)

				// Widgets
				wl.GET("/widgets", wlH.ListWidgets)
				wl.POST("/widgets", stubH.WLWidgetCreate)
				wl.PUT("/widgets/:id", wlH.UpdateWidget)

				// Pages
				wl.GET("/pages", wlH.ListPages)
				wl.POST("/pages", wlH.CreatePage)
				wl.PUT("/pages/:id", wlH.UpdatePage)

				// Smart Rewards
				wl.GET("/reward-rules", wlH.ListRewardRules)
				wl.POST("/reward-rules", wlH.CreateRewardRule)
				wl.GET("/privileges", wlH.ListPrivileges)
				wl.POST("/privileges", wlH.CreatePrivilege)
				wl.GET("/privileges/stats", stubH.WLPrivilegeStats)

				// Rules (frontend uses /rules not /reward-rules on some pages)
				wl.GET("/rules", stubH.WLRules)
				wl.POST("/rules", stubH.WLRuleCreate)
				wl.PUT("/rules/anti-sybil", stubH.WLAntiSybilUpdate)

				// Domain
				wl.GET("/domain", wlH.GetDomain)
				wl.POST("/domain", wlH.SetupDomain)
				wl.PUT("/domain", wlH.SetupDomain)
				wl.POST("/domain/verify", stubH.WLDomainVerify)

				// SDK
				wl.GET("/sdk", wlH.GetSDKConfig)
				wl.POST("/sdk/generate-key", stubH.WLSDKGenerateKey)
				wl.POST("/sdk/keys", stubH.WLSDKKeys)
				wl.POST("/sdk/webhooks", stubH.WLSDKWebhooks)
				wl.POST("/sso/test", stubH.WLSSOTest)

				// Brand
				wl.GET("/brand", stubH.WLBrand)
				wl.PUT("/brand", stubH.WLBrandUpdate)
				wl.POST("/brand/logo", stubH.WLBrandLogoUpload)

				// Contracts
				wl.GET("/contracts", stubH.WLContracts)
				wl.POST("/contracts", stubH.WLContractCreate)
				wl.GET("/contracts/check", stubH.WLContractCheck)
				wl.GET("/contracts/stats", stubH.WLContractStats)

				// Integration & Onboarding
				wl.GET("/integrations", stubH.WLIntegrations)
				wl.GET("/onboarding", stubH.WLOnboarding)
				wl.POST("/publish", stubH.WLPublish)
				wl.GET("/readiness", stubH.WLReadiness)
				wl.POST("/devkit/send", stubH.WLDevKitSend)
				wl.GET("/wizard/draft", stubH.WLWizardDraft)
			}
		}
	}

	// === C-End API (/api/c/) ===
	cAPI := r.Group("/api/c")
	{
		// Public (wallet connect)
		cAPI.POST("/wallet/connect", walletH.Connect)

		// Optional auth (can view without login)
		cPublic := cAPI.Group("")
		cPublic.Use(middleware.CEndOptionalAuth(cfg.JWTSecret))
		{
			cPublic.GET("/community/home", cCommunityH.Home)
			cPublic.GET("/community/announcements", cCommunityH.Announcements)
			cPublic.GET("/community/tasks", cCommunityH.Tasks)
			cPublic.GET("/community/achievements", stubH.CEndAchievements)
			cPublic.GET("/quests", cCommunityH.Quests)
			cPublic.GET("/leaderboard", cCommunityH.Leaderboard)
			cPublic.GET("/lb-sprint/current", cCommunityH.LBSprintCurrent)
			cPublic.GET("/lb-sprint/history", cCommunityH.LBSprintHistory)
			cPublic.GET("/milestones", cCommunityH.Milestones)
			cPublic.GET("/shop/items", cCommunityH.ShopItems)
		}

		// Protected C-end (requires wallet auth)
		cAuth := cAPI.Group("")
		cAuth.Use(middleware.CEndAuth(cfg.JWTSecret))
		{
			cAuth.GET("/user/status", cCommunityH.UserStatus)
			cAuth.GET("/user/profile", stubH.CEndUserProfile)
			cAuth.GET("/user/achievements", stubH.CEndAchievements)
			cAuth.GET("/user/referral-stats", stubH.CEndReferralStats)
			cAuth.GET("/user/activity", cCommunityH.ActivityFeed)
			cAuth.GET("/community/daychain", cCommunityH.DayChainStatus)
			cAuth.POST("/community/daychain", cCommunityH.DayChainCheckIn)
			cAuth.POST("/milestones/:id/claim", cCommunityH.ClaimMilestone)
			cAuth.POST("/shop/redeem", cCommunityH.RedeemShopItem)
			cAuth.GET("/invite/link", cCommunityH.InviteLink)
			cAuth.GET("/activity/feed", cCommunityH.ActivityFeed)
		}
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
