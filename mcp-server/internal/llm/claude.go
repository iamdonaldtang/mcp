package llm

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

const (
	ClaudeAPIURL = "https://api.anthropic.com/v1/messages"
	ClaudeModel  = "claude-sonnet-4-20250514"
)

// ClaudeClient Claude API客户端
type ClaudeClient struct {
	apiKey     string
	httpClient *http.Client
	tools      []Tool
}

// NewClaudeClient 创建Claude客户端
func NewClaudeClient(apiKey string) *ClaudeClient {
	client := &ClaudeClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // LLM响应可能较慢
		},
	}
	client.tools = client.buildTools()
	return client
}

// Tool Claude工具定义
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"input_schema"`
}

// InputSchema 输入Schema
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property Schema属性
type Property struct {
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Default     interface{}         `json:"default,omitempty"`
}

// Message 消息
type Message struct {
	Role    string        `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ContentBlock 内容块
type ContentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   string          `json:"content,omitempty"`
}

// ChatRequest Claude API请求
type ChatRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
	Tools     []Tool    `json:"tools,omitempty"`
}

// ChatResponse Claude API响应
type ChatResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        Usage          `json:"usage"`
}

// Usage 使用量统计
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Chat 发送对话请求
func (c *ClaudeClient) Chat(ctx context.Context, messages []Message, systemPrompt string) (*ChatResponse, error) {
	req := ChatRequest{
		Model:     ClaudeModel,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages:  messages,
		Tools:     c.tools,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ClaudeAPIURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	log.Debug().Int("message_count", len(messages)).Msg("调用Claude API")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Claude API返回错误 %d: %s", resp.StatusCode, string(respBody))
	}

	var result ChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	log.Debug().
		Str("stop_reason", result.StopReason).
		Int("input_tokens", result.Usage.InputTokens).
		Int("output_tokens", result.Usage.OutputTokens).
		Msg("Claude API响应")

	return &result, nil
}

// buildTools 构建工具定义
func (c *ClaudeClient) buildTools() []Tool {
	return []Tool{
		{
			Name:        "create_quest",
			Description: "创建完整的Quest活动配置。当基础信息、任务和奖励都已确认后使用此工具。",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"basic_info": {
						Type:        "object",
						Description: "活动基础信息",
						Properties: map[string]Property{
							"name":        {Type: "string", Description: "活动名称"},
							"description": {Type: "string", Description: "活动描述"},
							"start_time":  {Type: "string", Description: "开始时间，ISO 8601格式"},
							"end_time":    {Type: "string", Description: "结束时间，ISO 8601格式"},
							"cover_image": {Type: "string", Description: "封面图URL"},
							"is_private":  {Type: "boolean", Description: "是否私有活动", Default: false},
						},
						Required: []string{"name", "start_time", "end_time"},
					},
					"tasks": {
						Type:        "array",
						Description: "任务列表",
						Items: &Property{
							Type: "object",
							Properties: map[string]Property{
								"template_id": {
									Type:        "string",
									Description: "任务模板ID",
									Enum:        []string{"FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet", "JoinDiscord", "JoinTelegram", "TokenBalance", "NftHolder", "SwapVolume", "ContractInteractive", "DailyConnect", "Invite", "QuizChoose", "PowTask", "BindEmail"},
								},
								"custom_name": {Type: "string", Description: "自定义任务名称"},
								"params":      {Type: "object", Description: "任务参数"},
								"points":      {Type: "integer", Description: "积分", Default: 100},
								"is_optional": {Type: "boolean", Description: "是否可选", Default: false},
								"recurrence":  {Type: "string", Description: "重复周期", Enum: []string{"once", "daily", "weekly"}, Default: "once"},
							},
							Required: []string{"template_id", "params"},
						},
					},
					"rewards": {
						Type:        "object",
						Description: "奖励配置",
						Properties: map[string]Property{
							"distribution_method": {
								Type:        "string",
								Description: "分发方式",
								Enum:        []string{"fcfs", "lucky_draw", "ranking", "open_to_all"},
							},
							"layers": {
								Type:        "array",
								Description: "奖励层级",
								Items: &Property{
									Type: "object",
									Properties: map[string]Property{
										"max_winners": {Type: "integer", Description: "最大获奖人数"},
										"rewards": {
											Type: "array",
											Items: &Property{
												Type: "object",
												Properties: map[string]Property{
													"type":         {Type: "string", Enum: []string{"token", "nft", "points", "whitelist"}},
													"amount":       {Type: "string"},
													"token_symbol": {Type: "string"},
													"chain":        {Type: "string"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Required: []string{"basic_info", "tasks"},
			},
		},
		{
			Name:        "add_task",
			Description: "向当前Quest配置添加一个新任务",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"template_id": {
						Type:        "string",
						Description: "任务模板ID",
						Enum:        []string{"FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet", "JoinDiscord", "JoinTelegram", "TokenBalance", "SwapVolume", "ContractInteractive", "DailyConnect", "Invite", "QuizChoose", "PowTask"},
					},
					"params":      {Type: "object", Description: "任务参数"},
					"points":      {Type: "integer", Description: "积分", Default: 100},
					"is_optional": {Type: "boolean", Description: "是否可选", Default: false},
					"custom_name": {Type: "string", Description: "自定义名称"},
				},
				Required: []string{"template_id", "params"},
			},
		},
		{
			Name:        "suggest_tasks",
			Description: "根据项目类型和目标推荐任务组合",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"project_type": {
						Type:        "string",
						Description: "项目类型",
						Enum:        []string{"dex", "perps", "lending", "l2", "gamefi", "nft", "wallet", "other"},
					},
					"goal": {
						Type:        "string",
						Description: "主要目标",
						Enum:        []string{"acquisition", "activation", "retention", "trading_volume"},
					},
					"budget_level": {
						Type:        "string",
						Description: "预算级别",
						Enum:        []string{"low", "medium", "high"},
					},
				},
				Required: []string{"project_type", "goal"},
			},
		},
		{
			Name:        "validate_config",
			Description: "验证当前Quest配置是否完整有效",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"config": {Type: "object", Description: "要验证的配置"},
				},
				Required: []string{"config"},
			},
		},
		{
			Name:        "preview_config",
			Description: "生成配置的预览摘要",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"config": {Type: "object", Description: "要预览的配置"},
				},
				Required: []string{"config"},
			},
		},
		{
			Name:        "publish_quest",
			Description: "发布Quest到TaskOn平台",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"config":           {Type: "object", Description: "完整的Quest配置"},
					"notify_community": {Type: "boolean", Description: "是否通知社区", Default: false},
				},
				Required: []string{"config"},
			},
		},
	}
}

// ExtractToolCalls 从响应中提取工具调用
func ExtractToolCalls(response *ChatResponse) []ToolCall {
	var calls []ToolCall
	for _, block := range response.Content {
		if block.Type == "tool_use" {
			calls = append(calls, ToolCall{
				ID:    block.ID,
				Name:  block.Name,
				Input: block.Input,
			})
		}
	}
	return calls
}

// ToolCall 工具调用
type ToolCall struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// ExtractTextResponse 从响应中提取文本
func ExtractTextResponse(response *ChatResponse) string {
	for _, block := range response.Content {
		if block.Type == "text" {
			return block.Text
		}
	}
	return ""
}

// BuildToolResultMessage 构建工具结果消息
func BuildToolResultMessage(toolUseID string, result interface{}) Message {
	resultJSON, _ := json.Marshal(result)
	return Message{
		Role: "user",
		Content: []ContentBlock{
			{
				Type:      "tool_result",
				ToolUseID: toolUseID,
				Content:   string(resultJSON),
			},
		},
	}
}
