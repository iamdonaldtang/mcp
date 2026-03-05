package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/taskon/mcp-server/internal/api"
	"github.com/taskon/mcp-server/internal/llm"
	"github.com/taskon/mcp-server/internal/quest"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应限制
	},
}

// Server MCP服务器
type Server struct {
	taskonClient *api.TaskOnClient
	claudeClient *llm.ClaudeClient
	sessions     map[string]*Session
	sessionMu    sync.RWMutex
	systemPrompt string
}

// Session 用户会话
type Session struct {
	ID           string
	Messages     []llm.Message
	CurrentQuest *quest.QuestConfig
	State        ConversationState
}

// ConversationState 对话状态
type ConversationState string

const (
	StateDiscovery      ConversationState = "discovery"
	StateConfiguration  ConversationState = "configuration"
	StateConfirmation   ConversationState = "confirmation"
)

// NewServer 创建MCP服务器
func NewServer(taskonClient *api.TaskOnClient, claudeClient *llm.ClaudeClient) *Server {
	return &Server{
		taskonClient: taskonClient,
		claudeClient: claudeClient,
		sessions:     make(map[string]*Session),
		systemPrompt: getSystemPrompt(),
	}
}

// HandleWebSocket 处理WebSocket连接
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket升级失败")
		return
	}
	defer conn.Close()

	sessionID := uuid.New().String()
	session := &Session{
		ID:       sessionID,
		Messages: []llm.Message{},
		State:    StateDiscovery,
	}

	s.sessionMu.Lock()
	s.sessions[sessionID] = session
	s.sessionMu.Unlock()

	defer func() {
		s.sessionMu.Lock()
		delete(s.sessions, sessionID)
		s.sessionMu.Unlock()
	}()

	log.Info().Str("session_id", sessionID).Msg("新的WebSocket会话建立")

	// 发送欢迎消息
	welcomeMsg := MCPMessage{
		JSONRPC: "2.0",
		Method:  "session.created",
		Params: map[string]interface{}{
			"session_id": sessionID,
			"message":    "你好！我是TaskOn Quest Creator，可以帮助你创建营销活动。请告诉我你的项目类型和增长目标。",
		},
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Error().Err(err).Msg("发送欢迎消息失败")
		return
	}

	for {
		var msg MCPMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket读取错误")
			}
			break
		}

		response := s.handleMCPMessage(r.Context(), session, &msg)
		if err := conn.WriteJSON(response); err != nil {
			log.Error().Err(err).Msg("发送响应失败")
			break
		}
	}
}

// MCPMessage MCP消息格式
type MCPMessage struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id,omitempty"`
	Method  string                 `json:"method,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Result  interface{}            `json:"result,omitempty"`
	Error   *MCPError              `json:"error,omitempty"`
}

// MCPError MCP错误
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleMCPMessage 处理MCP消息
func (s *Server) handleMCPMessage(ctx context.Context, session *Session, msg *MCPMessage) *MCPMessage {
	switch msg.Method {
	case "chat":
		return s.handleChat(ctx, session, msg)
	case "get_templates":
		return s.handleGetTemplates(msg)
	case "get_session":
		return s.handleGetSession(session, msg)
	default:
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error:   &MCPError{Code: -32601, Message: "Method not found"},
		}
	}
}

// handleChat 处理对话请求
func (s *Server) handleChat(ctx context.Context, session *Session, msg *MCPMessage) *MCPMessage {
	userMessage, ok := msg.Params["message"].(string)
	if !ok || userMessage == "" {
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error:   &MCPError{Code: -32602, Message: "Invalid message parameter"},
		}
	}

	// 添加用户消息到会话
	session.Messages = append(session.Messages, llm.Message{
		Role: "user",
		Content: []llm.ContentBlock{
			{Type: "text", Text: userMessage},
		},
	})

	// 调用Claude API
	response, err := s.claudeClient.Chat(ctx, session.Messages, s.systemPrompt)
	if err != nil {
		log.Error().Err(err).Msg("Claude API调用失败")
		return &MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error:   &MCPError{Code: -32000, Message: fmt.Sprintf("LLM error: %v", err)},
		}
	}

	// 处理工具调用
	toolCalls := llm.ExtractToolCalls(response)
	var toolResults []interface{}

	for _, call := range toolCalls {
		result := s.executeToolCall(ctx, session, &call)
		toolResults = append(toolResults, result)
	}

	// 添加助手响应到会话
	session.Messages = append(session.Messages, llm.Message{
		Role:    "assistant",
		Content: response.Content,
	})

	textResponse := llm.ExtractTextResponse(response)

	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"message":      textResponse,
			"tool_calls":   toolCalls,
			"tool_results": toolResults,
			"state":        session.State,
			"quest":        session.CurrentQuest,
		},
	}
}

// executeToolCall 执行工具调用
func (s *Server) executeToolCall(ctx context.Context, session *Session, call *llm.ToolCall) interface{} {
	log.Debug().Str("tool", call.Name).Msg("执行工具调用")

	switch call.Name {
	case "create_quest":
		return s.toolCreateQuest(session, call.Input)
	case "add_task":
		return s.toolAddTask(session, call.Input)
	case "suggest_tasks":
		return s.toolSuggestTasks(call.Input)
	case "validate_config":
		return s.toolValidateConfig(call.Input)
	case "preview_config":
		return s.toolPreviewConfig(call.Input)
	case "publish_quest":
		return s.toolPublishQuest(ctx, session, call.Input)
	default:
		return map[string]interface{}{
			"error": fmt.Sprintf("Unknown tool: %s", call.Name),
		}
	}
}

// toolCreateQuest 创建Quest配置
func (s *Server) toolCreateQuest(session *Session, input json.RawMessage) interface{} {
	var config quest.QuestConfig
	if err := json.Unmarshal(input, &config); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	session.CurrentQuest = &config
	session.State = StateConfirmation

	return map[string]interface{}{
		"success": true,
		"message": "Quest配置已创建",
		"config":  config,
	}
}

// toolAddTask 添加任务
func (s *Server) toolAddTask(session *Session, input json.RawMessage) interface{} {
	if session.CurrentQuest == nil {
		session.CurrentQuest = &quest.QuestConfig{}
	}

	var task quest.TaskConfig
	if err := json.Unmarshal(input, &task); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	session.CurrentQuest.Tasks = append(session.CurrentQuest.Tasks, task)
	session.State = StateConfiguration

	return map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("已添加任务: %s", task.TemplateID),
		"task_count": len(session.CurrentQuest.Tasks),
	}
}

// toolSuggestTasks 推荐任务
func (s *Server) toolSuggestTasks(input json.RawMessage) interface{} {
	var params struct {
		ProjectType string `json:"project_type"`
		Goal        string `json:"goal"`
		BudgetLevel string `json:"budget_level"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	suggestions := quest.GetTaskSuggestions(params.ProjectType, params.Goal, params.BudgetLevel)
	return map[string]interface{}{
		"suggestions": suggestions,
	}
}

// toolValidateConfig 验证配置
func (s *Server) toolValidateConfig(input json.RawMessage) interface{} {
	var params struct {
		Config quest.QuestConfig `json:"config"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	errors := quest.ValidateConfig(&params.Config)
	return map[string]interface{}{
		"valid":  len(errors) == 0,
		"errors": errors,
	}
}

// toolPreviewConfig 预览配置
func (s *Server) toolPreviewConfig(input json.RawMessage) interface{} {
	var params struct {
		Config quest.QuestConfig `json:"config"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	preview := quest.GeneratePreview(&params.Config)
	return map[string]interface{}{
		"preview": preview,
	}
}

// toolPublishQuest 发布Quest
func (s *Server) toolPublishQuest(ctx context.Context, session *Session, input json.RawMessage) interface{} {
	var params struct {
		Config          quest.QuestConfig `json:"config"`
		NotifyCommunity bool              `json:"notify_community"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	// 验证配置
	errors := quest.ValidateConfig(&params.Config)
	if len(errors) > 0 {
		return map[string]interface{}{
			"success": false,
			"errors":  errors,
		}
	}

	// 转换为API格式
	apiParams, err := quest.ConvertToAPIParams(&params.Config)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	// 调用TaskOn API (需要用户授权token)
	// 这里返回转换后的配置，实际发布由前端携带用户token完成
	return map[string]interface{}{
		"success":    true,
		"message":    "配置已准备就绪，可以发布",
		"api_params": apiParams,
	}
}

// handleGetTemplates 获取任务模板
func (s *Server) handleGetTemplates(msg *MCPMessage) *MCPMessage {
	templates := api.GetSupportedTaskTemplates()
	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  templates,
	}
}

// handleGetSession 获取会话状态
func (s *Server) handleGetSession(session *Session, msg *MCPMessage) *MCPMessage {
	return &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"session_id": session.ID,
			"state":      session.State,
			"quest":      session.CurrentQuest,
		},
	}
}

// HandleChat REST API处理对话
func (s *Server) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"session_id"`
		Message   string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 获取或创建会话
	s.sessionMu.Lock()
	session, exists := s.sessions[req.SessionID]
	if !exists {
		session = &Session{
			ID:       uuid.New().String(),
			Messages: []llm.Message{},
			State:    StateDiscovery,
		}
		s.sessions[session.ID] = session
	}
	s.sessionMu.Unlock()

	// 构造MCP消息
	mcpMsg := &MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "chat",
		Params:  map[string]interface{}{"message": req.Message},
	}

	// 处理对话
	response := s.handleChat(r.Context(), session, mcpMsg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Result)
}

// getSystemPrompt 获取系统提示词
func getSystemPrompt() string {
	return `你是TaskOn Quest Creator，帮助Web3项目方通过对话创建营销活动。

## 你的能力
1. 理解用户的增长目标
2. 推荐合适的任务组合
3. 生成Quest配置JSON
4. 引导用户完善参数

## 对话流程
1. DISCOVERY: 了解项目类型和目标
2. CONFIGURATION: 配置任务和奖励
3. CONFIRMATION: 确认并发布

## 任务模板
### 社交任务
- FollowTwitter: 关注Twitter (参数: twitter_handle)
- RetweetTwitter: 转发推文 (参数: twitter_link)
- LikeATweet: 点赞推文 (参数: twitter_link)
- JoinDiscord: 加入Discord (参数: discord_server_url)
- JoinTelegram: 加入TG (参数: tg_group_link)

### 链上任务
- TokenBalance: 持有Token (参数: network, token_contract_addr, min_balance)
- SwapVolume: Swap交易量 (参数: chain, dex_name, min_volume)
- ContractInteractive: 合约交互 (参数: chain, contract_address)

### 互动任务
- DailyConnect: 每日签到
- Invite: 邀请好友 (参数: invite_num)
- QuizChoose: 问答 (参数: question, options)

## 奖励类型
- token: 代币奖励
- points: 积分奖励
- nft: NFT奖励
- whitelist: 白名单

## 分发方式
- fcfs: 先到先得
- lucky_draw: 抽奖
- ranking: 排名
- open_to_all: 全员

## 行为准则
- 友好、专业、简洁
- 每轮聚焦1-2个问题
- 给建议时说明理由
- 信息不足时主动询问`
}
