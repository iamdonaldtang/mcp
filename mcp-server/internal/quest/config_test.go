package quest

import (
	"encoding/json"
	"testing"
)

func TestConvertToAPIParams(t *testing.T) {
	config := &QuestConfig{
		BasicInfo: BasicInfo{
			Name:        "SwapX Launch Campaign",
			Description: "Join SwapX and complete tasks to win rewards!",
			StartTime:   "2026-02-10T00:00:00Z",
			EndTime:     "2026-02-17T23:59:59Z",
			IsPrivate:   false,
		},
		Tasks: []TaskConfig{
			{
				TemplateID: "FollowTwitter",
				CustomName: "Follow SwapX on Twitter",
				Params: map[string]interface{}{
					"twitter_handle": "@SwapX",
				},
				Points:     100,
				IsOptional: false,
				Recurrence: "once",
			},
			{
				TemplateID: "JoinDiscord",
				Params: map[string]interface{}{
					"discord_server_url": "https://discord.gg/swapx",
				},
				Points:     50,
				IsOptional: false,
				Recurrence: "once",
			},
			{
				TemplateID: "SwapVolume",
				CustomName: "Complete $100 Swap",
				Params: map[string]interface{}{
					"chain":      "base",
					"min_volume": "100",
				},
				Points:     200,
				IsOptional: false,
				Recurrence: "once",
			},
		},
		Rewards: &RewardConfig{
			DistributionMethod: "lucky_draw",
			Layers: []RewardLayer{
				{
					MaxWinners: 100,
					Rewards: []RewardItem{
						{
							Type:        "token",
							Amount:      "10",
							TokenSymbol: "USDC",
							Chain:       "base",
						},
					},
				},
			},
		},
	}

	apiParams, err := ConvertToAPIParams(config)
	if err != nil {
		t.Fatalf("ConvertToAPIParams failed: %v", err)
	}

	// 验证基础信息
	if apiParams.Name != config.BasicInfo.Name {
		t.Errorf("Name mismatch: got %s, want %s", apiParams.Name, config.BasicInfo.Name)
	}
	if apiParams.CampaignType != "Quest" {
		t.Errorf("CampaignType should be Quest, got %s", apiParams.CampaignType)
	}

	// 验证任务数量
	if len(apiParams.Tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(apiParams.Tasks))
	}

	// 验证任务类型推断
	if apiParams.Tasks[0].ClassType != "OffChain" {
		t.Errorf("FollowTwitter should be OffChain, got %s", apiParams.Tasks[0].ClassType)
	}
	if apiParams.Tasks[0].Platform != "Twitter" {
		t.Errorf("FollowTwitter platform should be Twitter, got %s", apiParams.Tasks[0].Platform)
	}
	if apiParams.Tasks[2].ClassType != "OnChain" {
		t.Errorf("SwapVolume should be OnChain, got %s", apiParams.Tasks[2].ClassType)
	}

	// 验证奖励转换
	if len(apiParams.WinnerRewards) != 1 {
		t.Errorf("Expected 1 winner reward config, got %d", len(apiParams.WinnerRewards))
	}
	if apiParams.WinnerRewards[0].WinnerDrawType != "LuckyDraw" {
		t.Errorf("WinnerDrawType should be LuckyDraw, got %s", apiParams.WinnerRewards[0].WinnerDrawType)
	}

	// 输出JSON验证
	jsonBytes, _ := json.MarshalIndent(apiParams, "", "  ")
	t.Logf("Converted API Params:\n%s", string(jsonBytes))
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *QuestConfig
		expectValid bool
		errorCount  int
	}{
		{
			name: "Valid config",
			config: &QuestConfig{
				BasicInfo: BasicInfo{
					Name:      "Test Quest",
					StartTime: "2026-02-10T00:00:00Z",
					EndTime:   "2026-02-17T23:59:59Z",
				},
				Tasks: []TaskConfig{
					{
						TemplateID: "FollowTwitter",
						Params:     map[string]interface{}{"twitter_handle": "@test"},
					},
				},
			},
			expectValid: true,
			errorCount:  0,
		},
		{
			name: "Missing name",
			config: &QuestConfig{
				BasicInfo: BasicInfo{
					StartTime: "2026-02-10T00:00:00Z",
					EndTime:   "2026-02-17T23:59:59Z",
				},
				Tasks: []TaskConfig{
					{
						TemplateID: "FollowTwitter",
						Params:     map[string]interface{}{"twitter_handle": "@test"},
					},
				},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "End time before start time",
			config: &QuestConfig{
				BasicInfo: BasicInfo{
					Name:      "Test Quest",
					StartTime: "2026-02-17T00:00:00Z",
					EndTime:   "2026-02-10T00:00:00Z",
				},
				Tasks: []TaskConfig{
					{
						TemplateID: "FollowTwitter",
						Params:     map[string]interface{}{"twitter_handle": "@test"},
					},
				},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "No tasks",
			config: &QuestConfig{
				BasicInfo: BasicInfo{
					Name:      "Test Quest",
					StartTime: "2026-02-10T00:00:00Z",
					EndTime:   "2026-02-17T23:59:59Z",
				},
				Tasks: []TaskConfig{},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "Invalid twitter handle",
			config: &QuestConfig{
				BasicInfo: BasicInfo{
					Name:      "Test Quest",
					StartTime: "2026-02-10T00:00:00Z",
					EndTime:   "2026-02-17T23:59:59Z",
				},
				Tasks: []TaskConfig{
					{
						TemplateID: "FollowTwitter",
						Params:     map[string]interface{}{"twitter_handle": "noatsign"},
					},
				},
			},
			expectValid: false,
			errorCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateConfig(tt.config)
			if tt.expectValid && len(errors) > 0 {
				t.Errorf("Expected valid config but got errors: %v", errors)
			}
			if !tt.expectValid && len(errors) == 0 {
				t.Error("Expected invalid config but got no errors")
			}
			if len(errors) != tt.errorCount {
				t.Errorf("Expected %d errors, got %d: %v", tt.errorCount, len(errors), errors)
			}
		})
	}
}

func TestGetTaskSuggestions(t *testing.T) {
	// DEX + acquisition
	suggestions := GetTaskSuggestions("dex", "acquisition", "medium")
	if len(suggestions) == 0 {
		t.Error("Expected suggestions for DEX project")
	}

	// 验证包含关键任务
	hasSwap := false
	hasTwitter := false
	for _, s := range suggestions {
		if s.TemplateID == "SwapVolume" {
			hasSwap = true
		}
		if s.TemplateID == "FollowTwitter" {
			hasTwitter = true
		}
	}
	if !hasSwap {
		t.Error("DEX project should have SwapVolume suggestion")
	}
	if !hasTwitter {
		t.Error("All projects should have FollowTwitter suggestion")
	}

	// GameFi + retention
	gamefiSuggestions := GetTaskSuggestions("gamefi", "retention", "low")
	hasDailyConnect := false
	for _, s := range gamefiSuggestions {
		if s.TemplateID == "DailyConnect" {
			hasDailyConnect = true
		}
	}
	if !hasDailyConnect {
		t.Error("GameFi retention should have DailyConnect suggestion")
	}
}

func TestGeneratePreview(t *testing.T) {
	config := &QuestConfig{
		BasicInfo: BasicInfo{
			Name:        "Test Quest",
			Description: "A test quest",
			StartTime:   "2026-02-10T00:00:00Z",
			EndTime:     "2026-02-17T23:59:59Z",
		},
		Tasks: []TaskConfig{
			{
				TemplateID: "FollowTwitter",
				CustomName: "Follow us",
				Points:     100,
			},
			{
				TemplateID: "JoinDiscord",
				Points:     50,
				IsOptional: true,
			},
		},
		Rewards: &RewardConfig{
			DistributionMethod: "lucky_draw",
			Layers: []RewardLayer{
				{
					MaxWinners: 50,
					Rewards: []RewardItem{
						{Type: "token", Amount: "10", TokenSymbol: "USDC"},
					},
				},
			},
		},
	}

	preview := GeneratePreview(config)
	if preview == "" {
		t.Error("Preview should not be empty")
	}

	// 验证包含关键信息
	if !contains(preview, "Test Quest") {
		t.Error("Preview should contain quest name")
	}
	if !contains(preview, "任务数量: 2") {
		t.Error("Preview should show task count")
	}
	if !contains(preview, "Follow us") {
		t.Error("Preview should show custom task name")
	}
	if !contains(preview, "(可选)") {
		t.Error("Preview should mark optional tasks")
	}

	t.Log(preview)
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
