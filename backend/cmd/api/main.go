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
			}

			// White Label Hub
			wl := b.Group("/whitelabel")
			{
				wl.GET("/overview", wlH.GetOverview)
				wl.POST("", wlH.Create)
				wl.PUT("/:id", wlH.Update)

				// Widgets
				wl.GET("/widgets", wlH.ListWidgets)
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

				// Domain
				wl.GET("/domain", wlH.GetDomain)
				wl.POST("/domain", wlH.SetupDomain)

				// SDK
				wl.GET("/sdk", wlH.GetSDKConfig)
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
			cAuth.GET("/community/daychain", cCommunityH.DayChainStatus)
			cAuth.POST("/community/daychain", cCommunityH.DayChainCheckIn)
			cAuth.POST("/milestones/:id/claim", cCommunityH.ClaimMilestone)
			cAuth.POST("/shop/redeem", cCommunityH.RedeemShopItem)
			cAuth.GET("/user/activity", cCommunityH.ActivityFeed)
			cAuth.GET("/invite/link", cCommunityH.InviteLink)
			cAuth.GET("/activity/feed", cCommunityH.ActivityFeed)
		}
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
