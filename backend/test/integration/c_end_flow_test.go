package integration

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/taskon/backend/internal/model"
	"gorm.io/gorm"
)

// setupBEndAndCommunity registers a B-end user, creates a community, and returns
// the Gin engine, DB, B-end JWT token, and community ID.
func setupBEndAndCommunity(t *testing.T) (*gin.Engine, *gorm.DB, string, string) {
	t.Helper()

	r, db := setupTestApp(t)
	bToken := registerAndLogin(t, r)

	w := doRequest(r, http.MethodPost, "/api/v1/community", map[string]string{
		"name":        "C-End Community",
		"description": "community for c-end tests",
	}, bToken)
	assertStatus(t, w, http.StatusCreated)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	communityID := data["id"].(string)
	if communityID == "" {
		t.Fatal("community creation did not return an id")
	}

	return r, db, bToken, communityID
}

func TestCEndFlow_WalletAndCommunity(t *testing.T) {
	r, _, _, communityID := setupBEndAndCommunity(t)

	// Wallet connect
	w := doRequest(r, http.MethodPost, "/api/c/wallet/connect", map[string]string{
		"address":      "0xTEST123",
		"community_id": communityID,
	}, "")
	assertStatus(t, w, http.StatusOK)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	cToken, ok := data["token"].(string)
	if !ok || cToken == "" {
		t.Fatal("wallet/connect did not return a token")
	}

	// Home
	w = doRequest(r, http.MethodGet, "/api/c/home", nil, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if _, hasPulse := data["pulse"]; !hasPulse {
		t.Fatal("home response missing 'pulse'")
	}
	if _, hasTab := data["tabVisibility"]; !hasTab {
		t.Fatal("home response missing 'tabVisibility'")
	}

	// Tasks
	w = doRequest(r, http.MethodGet, "/api/c/tasks", nil, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	if resp["data"] == nil {
		t.Fatal("tasks response data should not be nil")
	}

	// Leaderboard
	w = doRequest(r, http.MethodGet, "/api/c/leaderboard", nil, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if _, hasRankings := data["rankings"]; !hasRankings {
		t.Fatal("leaderboard response missing 'rankings'")
	}
}

func TestCEndFlow_DayChainCheckIn(t *testing.T) {
	r, db, _, communityID := setupBEndAndCommunity(t)

	// Enable DayChain for this community
	db.Model(&model.DayChainConfig{}).Where("community_id = ?", communityID).Update("enabled", true)

	// Wallet connect
	w := doRequest(r, http.MethodPost, "/api/c/wallet/connect", map[string]string{
		"address":      "0xDAYCHAIN",
		"community_id": communityID,
	}, "")
	assertStatus(t, w, http.StatusOK)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	cToken := data["token"].(string)

	// First check-in — streak should be 1
	w = doRequest(r, http.MethodPost, "/api/c/daychain/checkin", nil, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	streak := data["streak"].(float64)
	if streak != 1 {
		t.Fatalf("expected streak 1, got %v", streak)
	}

	// Second check-in same day — should be rejected
	w = doRequest(r, http.MethodPost, "/api/c/daychain/checkin", nil, cToken)
	assertStatus(t, w, http.StatusBadRequest)

	resp = parseResponse(t, w)
	msg := resp["message"].(string)
	if msg != "already checked in today" {
		t.Fatalf("expected 'already checked in today', got '%s'", msg)
	}
}

func TestCEndFlow_ShopRedeem(t *testing.T) {
	r, db, _, communityID := setupBEndAndCommunity(t)

	// Create shop item directly in DB
	shopItem := model.ShopItem{
		CommunityID: communityID,
		Name:        "Test Voucher",
		Category:    "voucher",
		Price:       50,
		Stock:       10,
		Status:      "active",
	}
	if err := db.Create(&shopItem).Error; err != nil {
		t.Fatalf("failed to create shop item: %v", err)
	}
	itemID := shopItem.ID

	// Wallet connect
	walletAddr := "0xSHOPPER"
	w := doRequest(r, http.MethodPost, "/api/c/wallet/connect", map[string]string{
		"address":      walletAddr,
		"community_id": communityID,
	}, "")
	assertStatus(t, w, http.StatusOK)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]interface{})
	cToken := data["token"].(string)

	// Give member 100 points
	db.Model(&model.CommunityMember{}).
		Where("community_id = ? AND wallet_address = ?", communityID, walletAddr).
		Update("total_points", 100)

	// List shop items — should have 1
	w = doRequest(r, http.MethodGet, "/api/c/shop", nil, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	items := resp["data"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 shop item, got %d", len(items))
	}

	// First redeem (100 - 50 = 50 remaining)
	w = doRequest(r, http.MethodPost, "/api/c/shop/redeem", map[string]string{
		"item_id": itemID,
	}, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if data["redeemed"] != true {
		t.Fatal("expected redeemed: true")
	}

	// Second redeem (50 - 50 = 0 remaining)
	w = doRequest(r, http.MethodPost, "/api/c/shop/redeem", map[string]string{
		"item_id": itemID,
	}, cToken)
	assertStatus(t, w, http.StatusOK)

	resp = parseResponse(t, w)
	data = resp["data"].(map[string]interface{})
	if data["redeemed"] != true {
		t.Fatal("expected redeemed: true on second redeem")
	}

	// Third redeem — insufficient points (0 < 50)
	w = doRequest(r, http.MethodPost, "/api/c/shop/redeem", map[string]string{
		"item_id": itemID,
	}, cToken)
	assertStatus(t, w, http.StatusBadRequest)

	resp = parseResponse(t, w)
	msg := resp["message"].(string)
	if msg != "insufficient points" {
		t.Fatalf("expected 'insufficient points', got '%s'", msg)
	}
}
