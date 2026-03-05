# TaskOn Quest MCP Server

基于Claude API的智能Quest创建服务，通过对话帮助Web3项目方创建营销活动。

## 功能特性

- **对话式创建Quest**: 通过自然语言对话配置活动
- **智能任务推荐**: 根据项目类型和目标推荐最佳任务组合
- **配置验证**: 自动验证Quest配置的完整性和有效性
- **格式转换**: 自动将LLM格式转换为TaskOn API格式

## 快速开始

### 环境要求

- Go 1.21+
- Claude API Key

### 安装依赖

```bash
cd mcp-server
go mod tidy
```

### 运行服务

```bash
# 开发模式
go run cmd/server/main.go --claude-key YOUR_API_KEY --debug

# 或使用环境变量
export CLAUDE_API_KEY=your_api_key
go run cmd/server/main.go --debug
```

### 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--port` | 8080 | 服务端口 |
| `--claude-key` | - | Claude API Key |
| `--taskon-api` | https://api.taskon.xyz | TaskOn API地址 |
| `--debug` | false | 调试模式 |

## API端点

### WebSocket: `/mcp`

MCP协议端点，支持以下方法：

```json
// 发送消息
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "chat",
  "params": {
    "message": "我想创建一个DEX推广活动"
  }
}

// 获取任务模板
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "get_templates"
}

// 获取会话状态
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "get_session"
}
```

### REST: `POST /api/chat`

简化的REST接口：

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "optional-session-id",
    "message": "我想创建一个NFT项目的Quest活动"
  }'
```

### 健康检查: `GET /health`

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## 项目结构

```
mcp-server/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── internal/
│   ├── api/
│   │   └── client.go        # TaskOn API客户端
│   ├── llm/
│   │   └── claude.go        # Claude API客户端
│   ├── mcp/
│   │   └── server.go        # MCP服务器核心
│   └── quest/
│       ├── config.go        # Quest配置和转换
│       └── config_test.go   # 测试用例
├── go.mod
├── go.sum
└── README.md
```

## 支持的任务模板

### 社交任务 (OffChain)

| 模板ID | 名称 | 必填参数 |
|--------|------|----------|
| FollowTwitter | 关注Twitter | twitter_handle |
| RetweetTwitter | 转发推文 | twitter_link |
| LikeATweet | 点赞推文 | twitter_link |
| PostTweet | 发布推文 | hash_tag |
| JoinDiscord | 加入Discord | discord_server_url |
| JoinTelegram | 加入Telegram | tg_group_link |

### 链上任务 (OnChain)

| 模板ID | 名称 | 必填参数 |
|--------|------|----------|
| TokenBalance | 持有Token | network, token_contract_addr, min_balance |
| NftHolder | 持有NFT | network, contract_address |
| SwapVolume | Swap交易量 | chain, min_volume |
| ContractInteractive | 合约交互 | chain, contract_address |

### 互动任务 (OffChain)

| 模板ID | 名称 | 必填参数 |
|--------|------|----------|
| DailyConnect | 每日签到 | 无 |
| Invite | 邀请好友 | invite_num |
| QuizChoose | 问答测验 | question, options |
| PowTask | 提交证明 | pow_type, desc |

## 对话示例

```
用户: 我想创建一个DEX推广活动，主要目标是增加交易量

助手: 了解！DEX项目的交易量目标，我推荐以下任务组合：
1. 关注Twitter - 建立基础关注
2. 加入Discord - 构建社区
3. 完成Swap交易 - 直接驱动交易量

请告诉我：
- 你的Twitter账号是？
- Discord邀请链接？
- 最低交易量要求（USD）？

用户: Twitter是@SwapX，Discord是discord.gg/swapx，最低100美元

助手: [调用create_quest工具生成配置]
配置已生成！活动摘要：
📋 活动名称: SwapX Trading Campaign
📌 任务: 3个
🎁 奖励: 待配置

是否需要添加奖励配置？
```

## 运行测试

```bash
go test -v ./...
```

## 配置文件

相关配置文件位于项目根目录：

- `configs/prompts/quest_creator_system_prompt.md` - 完整System Prompt
- `configs/prompts/system_prompt_short.txt` - 精简版
- `configs/prompts/claude_tools_definition.json` - Claude Tools定义
- `schemas/quest-config.schema.json` - Quest配置Schema
- `schemas/task-params.schema.json` - 任务参数Schema
- `schemas/config-mapping.md` - 格式转换文档

## License

MIT
