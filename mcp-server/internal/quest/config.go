package quest

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/taskon/mcp-server/internal/api"
)

// QuestConfig LLM生成的Quest配置格式
type QuestConfig struct {
	BasicInfo   BasicInfo          `json:"basic_info"`
	Tasks       []TaskConfig       `json:"tasks"`
	Rewards     *RewardConfig      `json:"rewards,omitempty"`
	Eligibility *EligibilityConfig `json:"eligibility,omitempty"`
}

// BasicInfo 基础信息
type BasicInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartTime   string `json:"start_time"` // ISO 8601格式
	EndTime     string `json:"end_time"`   // ISO 8601格式
	CoverImage  string `json:"cover_image,omitempty"`
	IsPrivate   bool   `json:"is_private"`
}

// TaskConfig 任务配置
type TaskConfig struct {
	TemplateID string                 `json:"template_id"`
	CustomName string                 `json:"custom_name,omitempty"`
	Params     map[string]interface{} `json:"params"`
	Points     int                    `json:"points"`
	IsOptional bool                   `json:"is_optional"`
	Recurrence string                 `json:"recurrence"` // once, daily, weekly
}

// RewardConfig 奖励配置
type RewardConfig struct {
	DistributionMethod string        `json:"distribution_method"` // fcfs, lucky_draw, ranking, open_to_all
	Layers             []RewardLayer `json:"layers"`
}

// RewardLayer 奖励层级
type RewardLayer struct {
	MaxWinners int          `json:"max_winners"`
	RankRange  *RankRange   `json:"rank_range,omitempty"`
	Rewards    []RewardItem `json:"rewards"`
}

// RankRange 排名范围
type RankRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// RewardItem 奖励项
type RewardItem struct {
	Type         string `json:"type"` // token, nft, points, whitelist
	Amount       string `json:"amount,omitempty"`
	TokenSymbol  string `json:"token_symbol,omitempty"`
	TokenAddress string `json:"token_address,omitempty"`
	Chain        string `json:"chain,omitempty"`
	PointsName   string `json:"points_name,omitempty"`
}

// EligibilityConfig 资格配置
type EligibilityConfig struct {
	Express    string                 `json:"express"` // and, or
	Conditions []EligibilityCondition `json:"conditions"`
}

// EligibilityCondition 资格条件
type EligibilityCondition struct {
	TemplateID string                 `json:"template_id"`
	Params     map[string]interface{} `json:"params"`
}

// ConvertToAPIParams 转换为API参数格式
func ConvertToAPIParams(config *QuestConfig) (*api.SetCampaignParams, error) {
	params := &api.SetCampaignParams{
		ID:           0, // 0表示创建新Campaign
		Name:         config.BasicInfo.Name,
		Desc:         config.BasicInfo.Description,
		CampaignImage: config.BasicInfo.CoverImage,
		CampaignType: "Quest",
		IsDraft:      false,
		IsPrivate:    config.BasicInfo.IsPrivate,
	}

	// 转换时间
	startTime, err := time.Parse(time.RFC3339, config.BasicInfo.StartTime)
	if err != nil {
		return nil, fmt.Errorf("解析开始时间失败: %w", err)
	}
	params.CampaignStart = startTime.Unix()

	endTime, err := time.Parse(time.RFC3339, config.BasicInfo.EndTime)
	if err != nil {
		return nil, fmt.Errorf("解析结束时间失败: %w", err)
	}
	params.CampaignEnd = endTime.Unix()

	// 转换任务
	for _, task := range config.Tasks {
		taskItem, err := convertTask(task)
		if err != nil {
			return nil, fmt.Errorf("转换任务失败: %w", err)
		}
		params.Tasks = append(params.Tasks, taskItem)
	}

	// 转换奖励
	if config.Rewards != nil {
		params.WinnerRewards = convertRewards(config.Rewards)
	}

	// 转换资格条件
	if config.Eligibility != nil {
		params.EligibilityExpress = config.Eligibility.Express
		params.Eligs = convertEligibility(config.Eligibility.Conditions)
	}

	return params, nil
}

// convertTask 转换单个任务
func convertTask(task TaskConfig) (api.TaskParamItem, error) {
	paramsJSON, err := json.Marshal(task.Params)
	if err != nil {
		return api.TaskParamItem{}, err
	}

	return api.TaskParamItem{
		TemplateID:      task.TemplateID,
		TemplateName:    task.TemplateID,
		Params:          string(paramsJSON),
		Points:          api.PointsInfo{Amount: task.Points},
		IsOptional:      task.IsOptional,
		IsTemplateTask:  false,
		IsEcosystemTask: false,
		ClassType:       inferClassType(task.TemplateID),
		IsHold:          false,
		Platform:        inferPlatform(task.TemplateID),
		CustomName:      task.CustomName,
		Recurrence:      capitalizeFirst(task.Recurrence),
	}, nil
}

// inferClassType 推断任务类型
func inferClassType(templateID string) string {
	switch templateID {
	case "TokenBalance", "NftHolder", "SwapVolume", "ContractInteractive":
		return "OnChain"
	default:
		return "OffChain"
	}
}

// inferPlatform 推断平台
func inferPlatform(templateID string) string {
	switch templateID {
	case "FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet", "QuoteTweetAndHashTag", "ReplyTweet":
		return "Twitter"
	case "JoinDiscord":
		return "Discord"
	case "JoinTelegram":
		return "Telegram"
	case "TokenBalance", "NftHolder", "SwapVolume", "ContractInteractive":
		return "Wallet"
	default:
		return ""
	}
}

// capitalizeFirst 首字母大写
func capitalizeFirst(s string) string {
	if s == "" {
		return "Once"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// convertRewards 转换奖励配置
func convertRewards(rewards *RewardConfig) []api.WinnerRewards {
	winnerReward := api.WinnerRewards{
		WinnerDrawType:              mapDistributionMethod(rewards.DistributionMethod),
		AutomaticallyWinnerDrawType: "Automatic",
	}

	for _, layer := range rewards.Layers {
		layerReward := api.WinnerLayerReward{
			WinnerLayer: api.WinnerLayer{MaxWinners: layer.MaxWinners},
			MaxWinners:  layer.MaxWinners,
		}

		for _, reward := range layer.Rewards {
			rewardInfo := api.RewardInfo{
				RewardType:   capitalizeFirst(reward.Type),
				Chain:        reward.Chain,
				Symbol:       reward.TokenSymbol,
				Amount:       reward.Amount,
				TokenAddress: reward.TokenAddress,
				PointsName:   reward.PointsName,
			}
			layerReward.Rewards = append(layerReward.Rewards, rewardInfo)
		}

		winnerReward.WinnerLayerRewards = append(winnerReward.WinnerLayerRewards, layerReward)
	}

	return []api.WinnerRewards{winnerReward}
}

// mapDistributionMethod 映射分发方式
func mapDistributionMethod(method string) string {
	switch method {
	case "fcfs":
		return "FCFS"
	case "lucky_draw":
		return "LuckyDraw"
	case "ranking":
		return "Ranking"
	case "open_to_all":
		return "OpenToAll"
	default:
		return "LuckyDraw"
	}
}

// convertEligibility 转换资格条件
func convertEligibility(conditions []EligibilityCondition) []api.EligibilityItem {
	var items []api.EligibilityItem
	for _, cond := range conditions {
		paramsJSON, _ := json.Marshal(cond.Params)
		items = append(items, api.EligibilityItem{
			TemplateID: cond.TemplateID,
			Params:     string(paramsJSON),
		})
	}
	return items
}

// ValidateConfig 验证Quest配置
func ValidateConfig(config *QuestConfig) []string {
	var errors []string

	// 验证基础信息
	if config.BasicInfo.Name == "" {
		errors = append(errors, "活动名称不能为空")
	}
	if config.BasicInfo.StartTime == "" {
		errors = append(errors, "开始时间不能为空")
	}
	if config.BasicInfo.EndTime == "" {
		errors = append(errors, "结束时间不能为空")
	}

	// 验证时间
	if config.BasicInfo.StartTime != "" && config.BasicInfo.EndTime != "" {
		start, err1 := time.Parse(time.RFC3339, config.BasicInfo.StartTime)
		end, err2 := time.Parse(time.RFC3339, config.BasicInfo.EndTime)
		if err1 != nil {
			errors = append(errors, "开始时间格式无效")
		}
		if err2 != nil {
			errors = append(errors, "结束时间格式无效")
		}
		if err1 == nil && err2 == nil && !end.After(start) {
			errors = append(errors, "结束时间必须晚于开始时间")
		}
	}

	// 验证任务
	if len(config.Tasks) == 0 {
		errors = append(errors, "至少需要一个任务")
	}

	for i, task := range config.Tasks {
		taskErrors := validateTaskParams(task)
		for _, e := range taskErrors {
			errors = append(errors, fmt.Sprintf("任务 %d: %s", i+1, e))
		}
	}

	return errors
}

// validateTaskParams 验证任务参数
func validateTaskParams(task TaskConfig) []string {
	var errors []string

	if task.TemplateID == "" {
		errors = append(errors, "缺少template_id")
		return errors
	}

	switch task.TemplateID {
	case "FollowTwitter":
		handle, ok := task.Params["twitter_handle"].(string)
		if !ok || handle == "" {
			errors = append(errors, "twitter_handle不能为空")
		} else if !strings.HasPrefix(handle, "@") {
			errors = append(errors, "twitter_handle必须以@开头")
		}

	case "RetweetTwitter", "LikeATweet":
		link, ok := task.Params["twitter_link"].(string)
		if !ok || link == "" {
			errors = append(errors, "twitter_link不能为空")
		} else if !strings.Contains(link, "twitter.com") && !strings.Contains(link, "x.com") {
			errors = append(errors, "无效的推文链接")
		}

	case "JoinDiscord":
		url, ok := task.Params["discord_server_url"].(string)
		if !ok || url == "" {
			errors = append(errors, "discord_server_url不能为空")
		} else if !strings.Contains(url, "discord") {
			errors = append(errors, "无效的Discord链接")
		}

	case "JoinTelegram":
		link, ok := task.Params["tg_group_link"].(string)
		if !ok || link == "" {
			errors = append(errors, "tg_group_link不能为空")
		} else if !strings.Contains(link, "t.me") {
			errors = append(errors, "无效的Telegram链接")
		}

	case "TokenBalance":
		if _, ok := task.Params["network"]; !ok {
			errors = append(errors, "network不能为空")
		}
		if _, ok := task.Params["token_contract_addr"]; !ok {
			errors = append(errors, "token_contract_addr不能为空")
		}
		if _, ok := task.Params["min_balance"]; !ok {
			errors = append(errors, "min_balance不能为空")
		}

	case "SwapVolume":
		if _, ok := task.Params["chain"]; !ok {
			errors = append(errors, "chain不能为空")
		}
		if _, ok := task.Params["min_volume"]; !ok {
			errors = append(errors, "min_volume不能为空")
		}

	case "ContractInteractive":
		if _, ok := task.Params["chain"]; !ok {
			errors = append(errors, "chain不能为空")
		}
		if _, ok := task.Params["contract_address"]; !ok {
			errors = append(errors, "contract_address不能为空")
		}

	case "Invite":
		if _, ok := task.Params["invite_num"]; !ok {
			errors = append(errors, "invite_num不能为空")
		}

	case "QuizChoose":
		if _, ok := task.Params["question"]; !ok {
			errors = append(errors, "question不能为空")
		}
		if _, ok := task.Params["options"]; !ok {
			errors = append(errors, "options不能为空")
		}
	}

	return errors
}

// GeneratePreview 生成配置预览
func GeneratePreview(config *QuestConfig) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📋 活动名称: %s\n", config.BasicInfo.Name))
	if config.BasicInfo.Description != "" {
		sb.WriteString(fmt.Sprintf("📝 描述: %s\n", config.BasicInfo.Description))
	}
	sb.WriteString(fmt.Sprintf("⏰ 时间: %s ~ %s\n", config.BasicInfo.StartTime, config.BasicInfo.EndTime))
	sb.WriteString(fmt.Sprintf("🔒 私有: %v\n", config.BasicInfo.IsPrivate))

	sb.WriteString(fmt.Sprintf("\n📌 任务数量: %d\n", len(config.Tasks)))
	for i, task := range config.Tasks {
		optional := ""
		if task.IsOptional {
			optional = " (可选)"
		}
		name := task.CustomName
		if name == "" {
			name = task.TemplateID
		}
		sb.WriteString(fmt.Sprintf("  %d. %s - %d积分%s\n", i+1, name, task.Points, optional))
	}

	if config.Rewards != nil {
		sb.WriteString(fmt.Sprintf("\n🎁 奖励分发: %s\n", config.Rewards.DistributionMethod))
		for _, layer := range config.Rewards.Layers {
			sb.WriteString(fmt.Sprintf("  - 最多%d人获奖\n", layer.MaxWinners))
			for _, reward := range layer.Rewards {
				if reward.Type == "token" {
					sb.WriteString(fmt.Sprintf("    • %s %s\n", reward.Amount, reward.TokenSymbol))
				} else if reward.Type == "points" {
					sb.WriteString(fmt.Sprintf("    • %s 积分\n", reward.Amount))
				} else {
					sb.WriteString(fmt.Sprintf("    • %s\n", reward.Type))
				}
			}
		}
	}

	return sb.String()
}

// TaskSuggestion 任务建议
type TaskSuggestion struct {
	TemplateID  string `json:"template_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	Priority    string `json:"priority"` // high, medium, low
}

// GetTaskSuggestions 获取任务建议
func GetTaskSuggestions(projectType, goal, budgetLevel string) []TaskSuggestion {
	var suggestions []TaskSuggestion

	// 基础社交任务 - 适用于所有项目
	suggestions = append(suggestions, TaskSuggestion{
		TemplateID:  "FollowTwitter",
		Name:        "关注Twitter",
		Description: "关注项目官方Twitter账号",
		Reason:      "建立社交媒体关注基础，成本低，转化率高",
		Priority:    "high",
	})

	suggestions = append(suggestions, TaskSuggestion{
		TemplateID:  "JoinDiscord",
		Name:        "加入Discord",
		Description: "加入项目Discord社区",
		Reason:      "构建活跃社区，便于后续运营和公告",
		Priority:    "high",
	})

	// 根据项目类型推荐
	switch projectType {
	case "dex", "perps":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "SwapVolume",
			Name:        "完成交易",
			Description: "在平台完成指定交易量",
			Reason:      "直接驱动交易量，验证真实用户",
			Priority:    "high",
		})
	case "lending":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "ContractInteractive",
			Name:        "存款/借款",
			Description: "与借贷协议交互",
			Reason:      "增加TVL，验证协议使用",
			Priority:    "high",
		})
	case "nft":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "NftHolder",
			Name:        "持有NFT",
			Description: "持有项目NFT",
			Reason:      "验证NFT持有者，建立holder社区",
			Priority:    "high",
		})
	case "gamefi":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "DailyConnect",
			Name:        "每日签到",
			Description: "每天登录签到",
			Reason:      "提高日活和留存率",
			Priority:    "high",
		})
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "Invite",
			Name:        "邀请好友",
			Description: "邀请好友参与",
			Reason:      "病毒式增长，降低获客成本",
			Priority:    "medium",
		})
	}

	// 根据目标推荐
	switch goal {
	case "acquisition":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "RetweetTwitter",
			Name:        "转发推文",
			Description: "转发项目公告推文",
			Reason:      "扩大传播范围，触达更多潜在用户",
			Priority:    "medium",
		})
	case "retention":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "DailyConnect",
			Name:        "每日签到",
			Description: "连续签到获得额外奖励",
			Reason:      "培养用户习惯，提高留存率",
			Priority:    "high",
		})
	case "activation":
		suggestions = append(suggestions, TaskSuggestion{
			TemplateID:  "QuizChoose",
			Name:        "知识问答",
			Description: "回答项目相关问题",
			Reason:      "教育用户，提高产品理解度",
			Priority:    "medium",
		})
	}

	return suggestions
}
