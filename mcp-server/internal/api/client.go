package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// TaskOnClient TaskOn API客户端
type TaskOnClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTaskOnClient 创建TaskOn API客户端
func NewTaskOnClient(baseURL string) *TaskOnClient {
	return &TaskOnClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetCampaignParams API要求的Campaign参数格式
type SetCampaignParams struct {
	ID                            int               `json:"id"`
	Name                          string            `json:"name"`
	Desc                          string            `json:"desc"`
	CampaignStart                 int64             `json:"campaign_start"`
	CampaignEnd                   int64             `json:"campaign_end"`
	CampaignImage                 string            `json:"campaign_image"`
	CampaignType                  string            `json:"campaign_type"`
	IsDraft                       bool              `json:"is_draft"`
	IsPrivate                     bool              `json:"is_private"`
	EligibilityExpress            string            `json:"eligibility_express"`
	WinnerRewards                 []WinnerRewards   `json:"winner_rewards"`
	QualifierRewards              []RewardInfo      `json:"qualifier_rewards"`
	Eligs                         []EligibilityItem `json:"eligs"`
	Tasks                         []TaskParamItem   `json:"tasks"`
	MinFinishedOptionalTaskNum    int               `json:"min_finished_optional_task_num"`
	MinFinishedOptionalTaskPoints int               `json:"min_finished_optional_task_points"`
	GoogleRecaptcha               bool              `json:"google_recaptcha"`
	IsPreview                     bool              `json:"is_preview"`
}

// TaskParamItem 任务参数项
type TaskParamItem struct {
	TemplateID      string     `json:"template_id"`
	TemplateName    string     `json:"template_name"`
	Params          string     `json:"params"` // JSON字符串
	Points          PointsInfo `json:"points"`
	IsOptional      bool       `json:"is_optional"`
	IsTemplateTask  bool       `json:"is_template_task"`
	IsEcosystemTask bool       `json:"is_ecosystem_task"`
	ClassType       string     `json:"class_type"`
	IsHold          bool       `json:"is_hold"`
	Platform        string     `json:"platform"`
	CustomName      string     `json:"custom_name"`
	Recurrence      string     `json:"recurrence"`
}

// PointsInfo 积分信息
type PointsInfo struct {
	Amount int `json:"amount"`
}

// WinnerRewards 获奖者奖励配置
type WinnerRewards struct {
	WinnerDrawType                string              `json:"winner_draw_type"`
	AutomaticallyWinnerDrawType   string              `json:"automatically_winner_draw_type"`
	WinnerLayerRewards            []WinnerLayerReward `json:"winner_layer_rewards"`
}

// WinnerLayerReward 奖励层级
type WinnerLayerReward struct {
	WinnerLayer WinnerLayer  `json:"winner_layer"`
	MaxWinners  int          `json:"max_winners"`
	Rewards     []RewardInfo `json:"rewards"`
}

// WinnerLayer 获奖层信息
type WinnerLayer struct {
	MaxWinners int `json:"max_winners"`
}

// RewardInfo 奖励信息
type RewardInfo struct {
	RewardType   string `json:"reward_type"`
	Chain        string `json:"chain,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	Amount       string `json:"amount,omitempty"`
	TokenAddress string `json:"token_address,omitempty"`
	PointsName   string `json:"points_name,omitempty"`
}

// EligibilityItem 资格条件项
type EligibilityItem struct {
	TemplateID string `json:"template_id"`
	Params     string `json:"params"` // JSON字符串
}

// CreateCampaign 创建Campaign
func (c *TaskOnClient) CreateCampaign(ctx context.Context, params *SetCampaignParams, authToken string) (*CampaignResponse, error) {
	url := fmt.Sprintf("%s/v1/setCampaign", c.baseURL)

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("序列化参数失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	log.Debug().Str("url", url).Msg("调用TaskOn API创建Campaign")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(respBody))
	}

	var result CampaignResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// CampaignResponse API响应
type CampaignResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    *CampaignData   `json:"data"`
}

// CampaignData Campaign数据
type CampaignData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	ShareLink string `json:"share_link"`
}

// GetCampaignInfo 获取Campaign信息
func (c *TaskOnClient) GetCampaignInfo(ctx context.Context, campaignID int, authToken string) (*CampaignResponse, error) {
	url := fmt.Sprintf("%s/v1/getCampaignInfo?id=%d", c.baseURL, campaignID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result CampaignResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetTaskTemplates 获取任务模板列表
func (c *TaskOnClient) GetTaskTemplates(ctx context.Context) ([]TaskTemplate, error) {
	// 返回预定义的任务模板列表
	return GetSupportedTaskTemplates(), nil
}

// TaskTemplate 任务模板
type TaskTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	ClassType   string   `json:"class_type"`
	Platform    string   `json:"platform"`
	Params      []string `json:"required_params"`
}

// GetSupportedTaskTemplates 获取支持的任务模板
func GetSupportedTaskTemplates() []TaskTemplate {
	return []TaskTemplate{
		// 社交任务 - Twitter
		{ID: "FollowTwitter", Name: "关注Twitter", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"twitter_handle"}},
		{ID: "RetweetTwitter", Name: "转发推文", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"twitter_link"}},
		{ID: "LikeATweet", Name: "点赞推文", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"twitter_link"}},
		{ID: "PostTweet", Name: "发布推文", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"hash_tag"}},
		{ID: "QuoteTweetAndHashTag", Name: "引用推文", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"twitter_link", "hash_tag"}},
		{ID: "ReplyTweet", Name: "回复推文", Category: "social", ClassType: "OffChain", Platform: "Twitter", Params: []string{"twitter_link"}},

		// 社交任务 - Discord/Telegram
		{ID: "JoinDiscord", Name: "加入Discord", Category: "social", ClassType: "OffChain", Platform: "Discord", Params: []string{"discord_server_url"}},
		{ID: "JoinTelegram", Name: "加入Telegram", Category: "social", ClassType: "OffChain", Platform: "Telegram", Params: []string{"tg_group_link"}},

		// 链上任务
		{ID: "TokenBalance", Name: "持有Token", Category: "onchain", ClassType: "OnChain", Platform: "Wallet", Params: []string{"network", "token_contract_addr", "min_balance"}},
		{ID: "NftHolder", Name: "持有NFT", Category: "onchain", ClassType: "OnChain", Platform: "Wallet", Params: []string{"network", "contract_address"}},
		{ID: "SwapVolume", Name: "Swap交易量", Category: "onchain", ClassType: "OnChain", Platform: "Wallet", Params: []string{"chain", "min_volume"}},
		{ID: "ContractInteractive", Name: "合约交互", Category: "onchain", ClassType: "OnChain", Platform: "Wallet", Params: []string{"chain", "contract_address"}},

		// 互动任务
		{ID: "DailyConnect", Name: "每日签到", Category: "engagement", ClassType: "OffChain", Platform: "", Params: []string{}},
		{ID: "Invite", Name: "邀请好友", Category: "engagement", ClassType: "OffChain", Platform: "", Params: []string{"invite_num"}},
		{ID: "QuizChoose", Name: "问答测验", Category: "engagement", ClassType: "OffChain", Platform: "", Params: []string{"question", "options"}},
		{ID: "PowTask", Name: "提交证明", Category: "engagement", ClassType: "OffChain", Platform: "", Params: []string{"pow_type", "desc"}},
		{ID: "BindEmail", Name: "绑定邮箱", Category: "engagement", ClassType: "OffChain", Platform: "", Params: []string{}},
	}
}
