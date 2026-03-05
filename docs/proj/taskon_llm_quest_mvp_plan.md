# TaskOn LLM Quest MVP 实施规划

> 生成时间：2026-02-03
> 目标：实现对话式Quest创建功能的MVP

---

## 一、需求确认总结

### 1.1 项目目标
将TaskOn从传统SaaS表单界面升级为**对话式配置引擎**，让项目方从业人员通过与LLM聊天即可完成Quest创建。

### 1.2 技术选型

| 组件 | 选择 | 说明 |
|------|------|------|
| **后端语言** | Go | 用户指定 |
| **LLM** | Claude API | 支持Tool Use，中文好 |
| **协议** | MCP (Model Context Protocol) | JSON-RPC 2.0 |
| **前端** | React + TypeScript | 与现有规范一致 |

### 1.3 MVP功能范围

| 功能 | 优先级 | 说明 |
|------|--------|------|
| 对话式Quest创建 | P0 | 通过多轮对话收集需求 |
| 任务模板配置 | P0 | 支持OffChain(社交)+OnChain(链上)任务 |
| 奖励设置 | P0 | Token/Points/NFT/Whitelist |
| 配置预览与发布 | P0 | JSON配置生成、预览、发布 |

### 1.4 现有API分析

**核心端点：**
- `POST /v1/setCampaign` - 创建/修改Campaign
- `POST /v1/getCampaignInfo` - 获取Campaign详情
- `POST /v1/getSetCampaignParams` - 获取创建参数
- `POST /v1/getCampaignList` - 查询Campaign列表

**核心Schema：**
- `SetCampaignParams` - 创建Campaign的完整参数
- `TaskParamItem` - 任务配置
- `WinnerRewards` - 奖励配置
- `EligibilityParamItem` - 前置条件配置

**支持的任务模板(TaskTemplateId)：**
```
// 社交类 (OffChain)
FollowTwitter, RetweetTwitter, LikeATweet, PostTweet, QuoteTweetAndTag
JoinDiscord, JoinTelegram, Youtube, ReplyTweet

// 链上类 (OnChain)
TokenBalance, NftHolder, ContractInteractive, SwapVolume
ChaingeSwap, ChaingeBridge, SwapDexContractInteractive

// 互动类
PowTask, QuizChoose, MultiQuizzes, LearnAndQuiz
DailyConnect, Invite, BindEmail
```

---

## 二、整体架构设计

```
┌─────────────────────────────────────────────────────────────────────┐
│                        用户界面层 (React)                            │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │ Chat UI      │  │ Config       │  │ Preview Panel            │   │
│  │ (对话界面)    │  │ Editor       │  │ (配置预览)                │   │
│  │              │  │ (可选编辑)    │  │                          │   │
│  └──────┬───────┘  └──────┬───────┘  └────────────┬─────────────┘   │
│         │                 │                       │                  │
│         ▼                 ▼                       ▼                  │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │              Quest Config State (Zustand)                    │    │
│  │              JSON配置状态管理                                 │    │
│  └──────────────────────────┬──────────────────────────────────┘    │
└─────────────────────────────┼───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      LLM 对话层 (Go Service)                         │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │ Conversation │  │ Claude API   │  │ JSON Config              │   │
│  │ Manager      │  │ Client       │  │ Generator                │   │
│  │ (对话管理)    │  │ (Tool Use)   │  │ (配置生成器)              │   │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘   │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │ Knowledge Base (知识库)                                        │   │
│  │ - 任务模板库 (TaskTemplateId映射)                              │   │
│  │ - 奖励类型说明                                                 │   │
│  │ - 最佳实践案例                                                 │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      MCP 服务层 (Go)                                 │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    Quest MCP Server                          │    │
│  │                    (JSON-RPC 2.0 over HTTP)                  │    │
│  └─────────────────────────────┬───────────────────────────────┘    │
│                                │                                     │
│  Tools:                        │                                     │
│  ├── create_quest             │  ← 创建Quest配置                    │
│  ├── update_quest             │  ← 更新Quest配置                    │
│  ├── get_task_templates       │  ← 获取任务模板列表                  │
│  ├── validate_config          │  ← 验证配置完整性                    │
│  ├── preview_quest            │  ← 预览Quest                        │
│  └── publish_quest            │  ← 发布Quest                        │
│                                                                      │
│  Resources:                                                          │
│  ├── quest://templates        │  ← 任务模板资源                      │
│  ├── quest://rewards          │  ← 奖励类型资源                      │
│  └── quest://eligibility      │  ← 前置条件资源                      │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      现有后端API                                     │
│  POST /v1/setCampaign  →  创建/修改Campaign                          │
│  POST /v1/getCampaignInfo  →  获取详情                               │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 三、JSON Schema 设计

### 3.1 Quest配置Schema (LLM输出格式)

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "QuestConfig",
  "description": "LLM生成的Quest配置，需转换为SetCampaignParams格式",
  "type": "object",
  "properties": {
    "basic_info": {
      "type": "object",
      "properties": {
        "name": { "type": "string", "description": "活动名称" },
        "description": { "type": "string", "description": "活动描述" },
        "cover_image": { "type": "string", "description": "封面图URL" },
        "start_time": { "type": "string", "format": "date-time" },
        "end_time": { "type": "string", "format": "date-time" },
        "is_private": { "type": "boolean", "default": false }
      },
      "required": ["name", "start_time", "end_time"]
    },
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "template_id": {
            "type": "string",
            "enum": ["FollowTwitter", "RetweetTwitter", "LikeATweet", "PostTweet",
                     "JoinDiscord", "JoinTelegram", "TokenBalance", "NftHolder",
                     "ContractInteractive", "SwapVolume", "QuizChoose", "DailyConnect",
                     "Invite", "PowTask", "BindEmail"]
          },
          "custom_name": { "type": "string", "description": "自定义任务名称" },
          "params": { "type": "object", "description": "任务参数" },
          "points": { "type": "integer", "description": "完成获得积分" },
          "is_optional": { "type": "boolean", "default": false },
          "recurrence": {
            "type": "string",
            "enum": ["once", "daily", "weekly"],
            "default": "once"
          }
        },
        "required": ["template_id", "params"]
      }
    },
    "rewards": {
      "type": "object",
      "properties": {
        "distribution_method": {
          "type": "string",
          "enum": ["fcfs", "lucky_draw", "ranking", "open_to_all"],
          "description": "分发方式：先到先得/抽奖/排名/全员"
        },
        "layers": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "max_winners": { "type": "integer" },
              "rewards": {
                "type": "array",
                "items": {
                  "type": "object",
                  "properties": {
                    "type": {
                      "type": "string",
                      "enum": ["token", "nft", "points", "whitelist", "discord_role"]
                    },
                    "amount": { "type": "string" },
                    "token_symbol": { "type": "string" },
                    "chain": { "type": "string" }
                  },
                  "required": ["type"]
                }
              }
            }
          }
        }
      },
      "required": ["distribution_method"]
    },
    "eligibility": {
      "type": "object",
      "properties": {
        "express": {
          "type": "string",
          "enum": ["and", "or"],
          "default": "and"
        },
        "conditions": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "template_id": {
                "type": "string",
                "enum": ["TokenBalance", "NftHolder", "DiscordRole", "Whitelist",
                         "Poh", "BABTHolder", "XAccountVerification"]
              },
              "params": { "type": "object" }
            }
          }
        }
      }
    }
  },
  "required": ["basic_info", "tasks"]
}
```

### 3.2 任务模板参数Schema (部分示例)

```json
{
  "FollowTwitter": {
    "twitter_handle": { "type": "string", "description": "Twitter账号，如 @TaskOnXyz" }
  },
  "RetweetTwitter": {
    "twitter_link": { "type": "string", "description": "推文链接" }
  },
  "LikeATweet": {
    "twitter_link": { "type": "string", "description": "推文链接" }
  },
  "JoinDiscord": {
    "discord_server_url": { "type": "string", "description": "Discord服务器邀请链接" }
  },
  "JoinTelegram": {
    "tg_group_link": { "type": "string", "description": "Telegram群组链接" }
  },
  "TokenBalance": {
    "network": { "type": "string", "description": "链名称，如 ethereum, bsc" },
    "token_contract_addr": { "type": "string", "description": "代币合约地址" },
    "token_name": { "type": "string", "description": "代币名称" },
    "min_balance": { "type": "string", "description": "最小余额" }
  },
  "SwapVolume": {
    "chain": { "type": "string" },
    "dex_name": { "type": "string" },
    "min_volume": { "type": "string" },
    "token_pair": { "type": "string" }
  },
  "QuizChoose": {
    "question": { "type": "string" },
    "options": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "label": { "type": "string" },
          "text": { "type": "string" },
          "is_answer": { "type": "boolean" }
        }
      }
    }
  }
}
```

---

## 四、MCP Server 设计 (Go实现)

### 4.1 目录结构

```
taskon-mcp-server/
├── cmd/
│   └── server/
│       └── main.go              # 入口
├── internal/
│   ├── mcp/
│   │   ├── server.go            # MCP Server核心
│   │   ├── tools.go             # Tool定义
│   │   ├── resources.go         # Resource定义
│   │   └── handlers.go          # 请求处理器
│   ├── llm/
│   │   ├── claude.go            # Claude API客户端
│   │   ├── conversation.go      # 对话管理
│   │   └── prompts.go           # System Prompt
│   ├── quest/
│   │   ├── config.go            # Quest配置结构
│   │   ├── converter.go         # 配置格式转换
│   │   └── validator.go         # 配置验证
│   ├── api/
│   │   └── taskon_client.go     # 现有API客户端
│   └── knowledge/
│       ├── templates.go         # 任务模板知识库
│       ├── rewards.go           # 奖励类型知识库
│       └── best_practices.go    # 最佳实践
├── pkg/
│   └── jsonrpc/
│       └── protocol.go          # JSON-RPC 2.0协议
├── configs/
│   ├── templates.json           # 任务模板定义
│   └── prompts/
│       └── system_prompt.md     # System Prompt
├── go.mod
└── go.sum
```

### 4.2 MCP Tools 定义

```go
// internal/mcp/tools.go

var QuestTools = []Tool{
    {
        Name:        "create_quest",
        Description: "创建新的Quest活动配置",
        InputSchema: QuestConfigSchema,
    },
    {
        Name:        "update_quest_field",
        Description: "更新Quest配置的特定字段",
        InputSchema: UpdateFieldSchema,
    },
    {
        Name:        "add_task",
        Description: "向Quest添加一个任务",
        InputSchema: TaskConfigSchema,
    },
    {
        Name:        "remove_task",
        Description: "从Quest移除指定任务",
        InputSchema: RemoveTaskSchema,
    },
    {
        Name:        "set_rewards",
        Description: "设置Quest奖励",
        InputSchema: RewardConfigSchema,
    },
    {
        Name:        "validate_config",
        Description: "验证当前Quest配置是否完整有效",
        InputSchema: ValidateSchema,
    },
    {
        Name:        "publish_quest",
        Description: "发布Quest到TaskOn平台",
        InputSchema: PublishSchema,
    },
    {
        Name:        "suggest_tasks",
        Description: "根据项目类型推荐任务组合",
        InputSchema: SuggestTasksSchema,
    },
}
```

### 4.3 配置转换器

```go
// internal/quest/converter.go

// QuestConfig 是LLM生成的配置格式
type QuestConfig struct {
    BasicInfo   BasicInfo     `json:"basic_info"`
    Tasks       []TaskConfig  `json:"tasks"`
    Rewards     RewardConfig  `json:"rewards"`
    Eligibility EligConfig    `json:"eligibility,omitempty"`
}

// ToSetCampaignParams 转换为API格式
func (q *QuestConfig) ToSetCampaignParams() (*SetCampaignParams, error) {
    params := &SetCampaignParams{
        ID:            0, // 0表示创建新Campaign
        Name:          q.BasicInfo.Name,
        Desc:          q.BasicInfo.Description,
        CampaignStart: q.BasicInfo.StartTime.Unix(),
        CampaignEnd:   q.BasicInfo.EndTime.Unix(),
        CampaignImage: q.BasicInfo.CoverImage,
        IsDraft:       false,
        IsPrivate:     q.BasicInfo.IsPrivate,
    }

    // 转换任务列表
    for _, task := range q.Tasks {
        taskItem, err := convertTask(task)
        if err != nil {
            return nil, err
        }
        params.Tasks = append(params.Tasks, taskItem)
    }

    // 转换奖励配置
    params.WinnerRewards = convertRewards(q.Rewards)

    // 转换前置条件
    if len(q.Eligibility.Conditions) > 0 {
        params.EligibilityExpress = q.Eligibility.Express
        params.Eligs = convertEligibility(q.Eligibility.Conditions)
    }

    return params, nil
}
```

---

## 五、LLM System Prompt 设计

### 5.1 核心Prompt结构

```markdown
# TaskOn Quest创建助手

你是TaskOn的Quest创建助手，帮助项目方通过对话快速创建营销活动。

## 你的能力
1. 理解用户的增长目标和需求
2. 推荐合适的任务组合
3. 生成符合规范的Quest配置JSON
4. 回答关于TaskOn功能的问题

## 对话流程

### Phase 1: DISCOVERY (发现需求)
收集基础信息：
- 项目类型（DEX/Perps/Lending/L2/GameFi等）
- 目标（获客/激活/留存/交易量）
- 预算和时间周期

### Phase 2: CONFIGURATION (配置任务)
确认具体配置：
- 活动名称、时间
- 任务列表（社交任务/链上任务）
- 每个任务的积分
- 奖励设置

### Phase 3: CONFIRMATION (确认发布)
- 展示配置摘要
- 确认发布

## 可用任务模板

### 社交任务 (OffChain)
| 模板ID | 名称 | 参数 |
|--------|------|------|
| FollowTwitter | 关注Twitter | twitter_handle: 账号 |
| RetweetTwitter | 转发推文 | twitter_link: 推文链接 |
| LikeATweet | 点赞推文 | twitter_link: 推文链接 |
| PostTweet | 发布推文 | hash_tag: 话题标签 |
| JoinDiscord | 加入Discord | discord_server_url: 服务器链接 |
| JoinTelegram | 加入Telegram | tg_group_link: 群组链接 |

### 链上任务 (OnChain)
| 模板ID | 名称 | 参数 |
|--------|------|------|
| TokenBalance | 持有Token | network, token_contract_addr, min_balance |
| NftHolder | 持有NFT | network, contract_address |
| SwapVolume | Swap交易量 | chain, dex_name, min_volume |
| ContractInteractive | 合约交互 | chain, contract_address |

### 互动任务
| 模板ID | 名称 | 参数 |
|--------|------|------|
| DailyConnect | 每日签到 | 无需参数 |
| Invite | 邀请好友 | invite_num: 邀请人数 |
| QuizChoose | 问答测验 | question, options |
| PowTask | 提交证明 | pow_type: Image/URL/Text |

## 奖励类型
- token: ERC20代币奖励
- nft: NFT奖励
- points: 积分奖励
- whitelist: 白名单
- discord_role: Discord角色

## 分发方式
- fcfs: 先到先得
- lucky_draw: 抽奖
- ranking: 按排名
- open_to_all: 全员获得

## 输出格式

当需求明确后，使用 create_quest 工具生成配置：

```json
{
  "basic_info": {
    "name": "活动名称",
    "description": "活动描述",
    "start_time": "2026-02-10T00:00:00Z",
    "end_time": "2026-02-24T23:59:59Z"
  },
  "tasks": [
    {
      "template_id": "FollowTwitter",
      "params": { "twitter_handle": "@ProjectName" },
      "points": 100
    }
  ],
  "rewards": {
    "distribution_method": "lucky_draw",
    "layers": [
      {
        "max_winners": 100,
        "rewards": [{ "type": "token", "amount": "10", "token_symbol": "USDC" }]
      }
    ]
  }
}
```

## 对话原则
1. 每轮聚焦1-2个问题，不要一次问太多
2. 给建议时说明理由
3. 使用简洁的中文交流
4. 配置要完整可执行
```

---

## 六、实施路线图

### Phase 1: 基础架构 (Week 1-2)

| 任务 | 交付物 | 优先级 |
|------|--------|--------|
| 搭建Go项目结构 | 基础代码框架 | P0 |
| 实现MCP Server核心 | JSON-RPC 2.0通信层 | P0 |
| 定义Quest Tools | create_quest等6个Tool | P0 |
| 实现配置转换器 | QuestConfig → SetCampaignParams | P0 |
| 对接现有API | TaskOn API Client | P0 |

### Phase 2: LLM对话层 (Week 2-3)

| 任务 | 交付物 | 优先级 |
|------|--------|--------|
| Claude API集成 | Tool Use调用 | P0 |
| 对话管理器 | 状态机管理 | P0 |
| System Prompt | 完整提示词 | P0 |
| 知识库构建 | 任务模板、最佳实践 | P1 |

### Phase 3: 前端集成 (Week 3-4)

| 任务 | 交付物 | 优先级 |
|------|--------|--------|
| Chat UI组件 | 对话界面 | P0 |
| Config Preview | 配置预览面板 | P0 |
| MCP Client | WebSocket通信 | P0 |
| 状态管理 | Zustand Store | P1 |

### Phase 4: 测试与优化 (Week 4-5)

| 任务 | 交付物 | 优先级 |
|------|--------|--------|
| 端到端测试 | 完整流程验证 | P0 |
| Prompt调优 | 对话质量优化 | P1 |
| 错误处理 | 边界情况处理 | P1 |
| 文档编写 | 使用说明 | P2 |

---

## 七、关键技术决策

### 7.1 为什么选择MCP而不是直接REST API

| 维度 | MCP | 直接REST |
|------|-----|----------|
| LLM集成 | ✅ 原生Tool Use支持 | ❌ 需要额外适配 |
| 标准化 | ✅ 行业标准协议 | ⚠️ 自定义格式 |
| 可扩展性 | ✅ 易于添加新工具 | ⚠️ 需要修改代码 |
| 调试 | ✅ 结构化日志 | ⚠️ 自建 |

### 7.2 配置格式转换策略

```
用户输入 → LLM理解 → QuestConfig (简化格式) → SetCampaignParams (API格式) → 后端API
```

- **QuestConfig**: LLM易于理解和生成的简化格式
- **SetCampaignParams**: 现有API要求的完整格式
- **转换层**: 负责字段映射、默认值填充、格式转换

### 7.3 对话状态管理

```go
type ConversationState struct {
    SessionID     string
    Phase         Phase  // DISCOVERY | CONFIGURATION | CONFIRMATION
    CollectedInfo map[string]interface{}
    CurrentConfig *QuestConfig
    MissingFields []string
    History       []Message
}
```

---

## 八、验收标准

### 8.1 功能验收

- [ ] 用户可以通过对话创建包含3-5个任务的Quest
- [ ] 支持至少5种社交任务模板
- [ ] 支持至少3种链上任务模板
- [ ] 支持Token/Points奖励设置
- [ ] 配置可正确转换并调用现有API
- [ ] 发布后可在TaskOn平台查看

### 8.2 质量验收

- [ ] 平均对话轮数 < 10轮完成配置
- [ ] 配置验证通过率 > 95%
- [ ] API调用成功率 > 99%
- [ ] 响应时间 < 3秒 (LLM调用除外)

---

## 九、风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| LLM生成配置不准确 | 高 | 严格的Schema验证 + Few-shot示例 |
| API格式转换错误 | 高 | 完整的单元测试 + 集成测试 |
| 对话理解偏差 | 中 | 多轮确认机制 + 配置预览 |
| 现有API变更 | 低 | 版本控制 + 兼容层 |

---

## 十、下一步行动

1. **立即开始**: 搭建Go项目结构，实现MCP Server骨架
2. **本周完成**: Quest Tools定义，配置转换器
3. **下周目标**: Claude API集成，对话管理器
4. **两周后**: 前端集成，端到端测试

---

*文档版本: V1.0*
*创建时间: 2026-02-03*
