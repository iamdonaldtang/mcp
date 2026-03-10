package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/taskon/backend/internal/model"
	"github.com/taskon/backend/internal/testutil"
)

func setupCEndRouter(db *gorm.DB) *gin.Engine {
	cH := &CCommunityHandler{DB: db}
	wH := &WalletHandler{DB: db, Secret: testJWTSecret}
	r := gin.New()
	r.POST("/wallet/connect", wH.Connect)
	cEnd := r.Group("/c")
	cEnd.Use(func(c *gin.Context) {
		c.Set("wallet_address", "0xTEST1234")
		c.Set("community_id", c.Query("cid"))
		c.Next()
	})
	cEnd.GET("/home", cH.Home)
	cEnd.GET("/tasks", cH.Tasks)
	cEnd.POST("/daychain/checkin", cH.DayChainCheckIn)
	cEnd.GET("/leaderboard", cH.Leaderboard)
	cEnd.GET("/shop", cH.ShopItems)
	cEnd.POST("/shop/redeem", cH.RedeemShopItem)
	cEnd.GET("/milestones", cH.Milestones)
	cEnd.POST("/milestones/:id/claim", cH.ClaimMilestone)
	return r
}

func createTestMember(t *testing.T, db *gorm.DB, communityID, wallet string, points int) model.CommunityMember {
	t.Helper()
	m := model.CommunityMember{
		ID:            uuid.New().String(),
		CommunityID:   communityID,
		WalletAddress: wallet,
		Level:         1,
		TotalPoints:   points,
		LifetimeXP:    points,
		DayStreak:     0,
		JoinedAt:      time.Now(),
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create test member: %v", err)
	}
	return m
}

func createTestDayChainConfig(t *testing.T, db *gorm.DB, communityID string, enabled bool) model.DayChainConfig {
	t.Helper()
	dc := model.DayChainConfig{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		Enabled:     enabled,
		TargetDays:  30,
		DailyPoints: 10,
	}
	if err := db.Create(&dc).Error; err != nil {
		t.Fatalf("failed to create test day chain config: %v", err)
	}
	return dc
}

func createTestShopItem(t *testing.T, db *gorm.DB, communityID, name string, price, stock int) model.ShopItem {
	t.Helper()
	item := model.ShopItem{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		Name:        name,
		Category:    "voucher",
		Price:       price,
		Stock:       stock,
		Status:      "active",
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("failed to create test shop item: %v", err)
	}
	return item
}

func createTestMilestone(t *testing.T, db *gorm.DB, communityID, name string, threshold int) model.Milestone {
	t.Helper()
	m := model.Milestone{
		ID:          uuid.New().String(),
		CommunityID: communityID,
		Name:        name,
		Threshold:   threshold,
		RewardType:  "points",
		RewardValue: "100",
		Status:      "active",
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create test milestone: %v", err)
	}
	return m
}

// --- Wallet Connect tests ---

func TestWalletConnect_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	r := setupCEndRouter(db)

	w := postJSON(r, "/wallet/connect", WalletConnectRequest{
		Address:     "0xABCDEF1234567890",
		CommunityID: community.ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["token"] == nil || data["token"] == "" {
		t.Error("expected token in response")
	}

	// Verify member was created
	var count int64
	db.Model(&model.CommunityMember{}).Where("wallet_address = ?", "0xABCDEF1234567890").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 member, got %d", count)
	}
}

func TestWalletConnect_MissingFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCEndRouter(db)

	// Missing address
	w := postJSON(r, "/wallet/connect", map[string]string{
		"community_id": "some-id",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing address, got %d", w.Code)
	}

	// Missing community_id
	w = postJSON(r, "/wallet/connect", map[string]string{
		"address": "0xABC",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing community_id, got %d", w.Code)
	}
}

func TestWalletConnect_InvalidCommunity(t *testing.T) {
	db := testutil.SetupTestDB(t)
	r := setupCEndRouter(db)

	w := postJSON(r, "/wallet/connect", WalletConnectRequest{
		Address:     "0xABCDEF",
		CommunityID: "nonexistent-community-id",
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Home tests ---

func TestCEndHome(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestMember(t, db, community.ID, "0xTEST1234", 100)

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/c/home?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})

	// Check community info
	comm := data["community"].(map[string]interface{})
	if comm["name"] != "Test Community" {
		t.Errorf("expected community name 'Test Community', got %v", comm["name"])
	}

	// Check pulse stats
	pulse := data["pulse"].(map[string]interface{})
	if pulse["total_members"].(float64) != 1 {
		t.Errorf("expected total_members 1, got %v", pulse["total_members"])
	}

	// Check tab visibility
	tabs := data["tabVisibility"].(map[string]interface{})
	if tabs["home"] != true {
		t.Error("expected home tab to be visible")
	}
}

// --- Tasks tests ---

func TestCEndTasks(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	sector := createTestSector(t, db, community.ID, "General", 1)
	createTestTask(t, db, community.ID, sector.ID, "Follow Twitter", "social")
	createTestMember(t, db, community.ID, "0xTEST1234", 0)

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/c/tasks?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected 1 sector, got %d", len(data))
	}
	sectorData := data[0].(map[string]interface{})
	tasks := sectorData["tasks"].([]interface{})
	if len(tasks) != 1 {
		t.Errorf("expected 1 task in sector, got %d", len(tasks))
	}
}

// --- DayChain tests ---

func TestDayChainCheckIn_NewStreak(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestDayChainConfig(t, db, community.ID, true)
	createTestMember(t, db, community.ID, "0xTEST1234", 0)

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/c/daychain/checkin?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["streak"].(float64) != 1 {
		t.Errorf("expected streak 1, got %v", data["streak"])
	}
	if data["points_earned"].(float64) != 10 {
		t.Errorf("expected points_earned 10, got %v", data["points_earned"])
	}
}

func TestDayChainCheckIn_Duplicate(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestDayChainConfig(t, db, community.ID, true)

	// Create member with LastCheckIn set to today
	now := time.Now()
	member := model.CommunityMember{
		ID:            uuid.New().String(),
		CommunityID:   community.ID,
		WalletAddress: "0xTEST1234",
		Level:         1,
		TotalPoints:   10,
		DayStreak:     1,
		LastCheckIn:   &now,
		JoinedAt:      now,
	}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("failed to create member: %v", err)
	}

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/c/daychain/checkin?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for duplicate check-in, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Leaderboard tests ---

func TestLeaderboard(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestMember(t, db, community.ID, "0xAAA", 500)
	createTestMember(t, db, community.ID, "0xBBB", 300)
	createTestMember(t, db, community.ID, "0xCCC", 100)

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/c/leaderboard?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	rankings := data["rankings"].([]interface{})
	if len(rankings) != 3 {
		t.Fatalf("expected 3 rankings, got %d", len(rankings))
	}

	// First place should be 0xAAA with 500 points
	first := rankings[0].(map[string]interface{})
	if first["wallet_address"] != "0xAAA" {
		t.Errorf("expected first place '0xAAA', got %v", first["wallet_address"])
	}
	if first["rank"].(float64) != 1 {
		t.Errorf("expected rank 1, got %v", first["rank"])
	}
}

// --- Shop tests ---

func TestShopItems(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	createTestShopItem(t, db, community.ID, "T-Shirt", 100, 10)
	createTestShopItem(t, db, community.ID, "Sticker", 20, 50)

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/c/shop?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 shop items, got %d", len(data))
	}
}

func TestRedeemShopItem(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	item := createTestShopItem(t, db, community.ID, "T-Shirt", 50, 10)
	createTestMember(t, db, community.ID, "0xTEST1234", 200)

	r := setupCEndRouter(db)

	w := postJSON(r, "/c/shop/redeem?cid="+community.ID, map[string]string{
		"item_id": item.ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["redeemed"] != true {
		t.Error("expected redeemed to be true")
	}
	if data["remaining_points"].(float64) != 150 {
		t.Errorf("expected remaining_points 150, got %v", data["remaining_points"])
	}

	// Verify stock was decremented
	var updatedItem model.ShopItem
	db.First(&updatedItem, "id = ?", item.ID)
	if updatedItem.Stock != 9 {
		t.Errorf("expected stock 9, got %d", updatedItem.Stock)
	}
}

func TestRedeemShopItem_InsufficientPoints(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	item := createTestShopItem(t, db, community.ID, "Expensive Item", 1000, 5)
	createTestMember(t, db, community.ID, "0xTEST1234", 10) // Only 10 points

	r := setupCEndRouter(db)

	w := postJSON(r, "/c/shop/redeem?cid="+community.ID, map[string]string{
		"item_id": item.ID,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for insufficient points, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Milestone tests ---

func TestClaimMilestone(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	milestone := createTestMilestone(t, db, community.ID, "First 100 Points", 100)
	createTestMember(t, db, community.ID, "0xTEST1234", 200) // Meets threshold

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/c/milestones/"+milestone.ID+"/claim?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	if data["claimed"] != true {
		t.Error("expected claimed to be true")
	}

	// Verify claim was persisted
	var claimCount int64
	db.Model(&model.MilestoneClaim{}).Where("milestone_id = ? AND wallet_address = ?", milestone.ID, "0xTEST1234").Count(&claimCount)
	if claimCount != 1 {
		t.Errorf("expected 1 claim, got %d", claimCount)
	}
}

func TestClaimMilestone_AlreadyClaimed(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	milestone := createTestMilestone(t, db, community.ID, "First 100 Points", 100)
	createTestMember(t, db, community.ID, "0xTEST1234", 200)

	// Create existing claim
	db.Create(&model.MilestoneClaim{
		ID:            uuid.New().String(),
		CommunityID:   community.ID,
		MilestoneID:   milestone.ID,
		WalletAddress: "0xTEST1234",
		ClaimedAt:     time.Now(),
	})

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/c/milestones/"+milestone.ID+"/claim?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for already claimed, got %d: %s", w.Code, w.Body.String())
	}
}

func TestClaimMilestone_InsufficientPoints(t *testing.T) {
	db := testutil.SetupTestDB(t)
	community := createTestCommunity(t, db, "test-user-1", "Test Community")
	milestone := createTestMilestone(t, db, community.ID, "First 1000 Points", 1000)
	createTestMember(t, db, community.ID, "0xTEST1234", 50) // Below threshold

	r := setupCEndRouter(db)

	req := httptest.NewRequest(http.MethodPost, "/c/milestones/"+milestone.ID+"/claim?cid="+community.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for insufficient points, got %d: %s", w.Code, w.Body.String())
	}
}
