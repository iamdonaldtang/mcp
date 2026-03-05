# TaskOn LLM + MCP + Widget 架构升级方案

## 🎯 项目目标

将TaskOn从传统SaaS界面升级为**对话式配置引擎**，实现：
1. 用户通过与LLM对话快速明确需求
2. LLM引导用户完成运营认知提升
3. 需求确认后输出标准化JSON配置
4. Widget根据JSON配置动态渲染组件
5. 后端API改造为MCP服务，实现标准化数据交互

---

## 🏗️ 整体架构设计

```
┌─────────────────────────────────────────────────────────────────────┐
│                        用户界面层 (React)                            │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │ Chat UI      │  │ Widget       │  │ Preview Panel            │   │
│  │ (对话界面)    │  │ Renderer     │  │ (配置预览)                │   │
│  │              │  │ (组件渲染器)  │  │                          │   │
│  └──────┬───────┘  └──────┬───────┘  └────────────┬─────────────┘   │
│         │                 │                       │                  │
│         ▼                 ▼                       ▼                  │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │              JSON Configuration Store                        │    │
│  │              (配置状态管理 - Zustand/Redux)                   │    │
│  └──────────────────────────┬──────────────────────────────────┘    │
└─────────────────────────────┼───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      LLM 对话层 (Orchestrator)                       │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │ Conversation │  │ Intent       │  │ JSON Schema              │   │
│  │ Manager      │  │ Classifier   │  │ Generator                │   │
│  │ (对话管理)    │  │ (意图识别)    │  │ (配置生成器)              │   │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘   │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │ Knowledge Base (运营知识库)                                    │   │
│  │ - 游戏化组件最佳实践                                           │   │
│  │ - 行业案例模板                                                 │   │
│  │ - 指标体系与ROI框架                                            │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      MCP 服务层 (Protocol Layer)                     │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    MCP Gateway / Router                      │    │
│  │                    (JSON-RPC 2.0 over HTTP)                  │    │
│  └─────────────────────────────┬───────────────────────────────┘    │
│                                │                                     │
│  ┌─────────────┬───────────────┼───────────────┬─────────────┐      │
│  │             │               │               │             │      │
│  ▼             ▼               ▼               ▼             ▼      │
│ ┌────────┐ ┌────────┐   ┌────────────┐  ┌────────┐  ┌────────────┐ │
│ │Quest   │ │Points  │   │Leaderboard │  │Rewards │  │Anti-Sybil  │ │
│ │MCP     │ │MCP     │   │MCP Server  │  │MCP     │  │MCP Server  │ │
│ │Server  │ │Server  │   │            │  │Server  │  │            │ │
│ └────────┘ └────────┘   └────────────┘  └────────┘  └────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       数据持久层 (Database)                          │
│  PostgreSQL / MongoDB / Redis                                        │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📊 核心数据流

### Flow 1: 对话式配置生成

```
用户输入 → LLM理解意图 → 多轮对话澄清需求 → 生成配置JSON → Widget渲染预览 → 用户确认 → MCP持久化
```

### Flow 2: 配置执行与数据同步

```
Widget交互 → 触发MCP Tool调用 → 后端处理 → 返回结果 → 更新Widget状态
```

---

## 🤖 LLM 对话层详细设计

### 1. 对话管理器 (Conversation Manager)

```typescript
interface ConversationState {
  sessionId: string;
  currentPhase: 'discovery' | 'clarification' | 'configuration' | 'confirmation';
  collectedRequirements: PartialCampaignConfig;
  missingFields: string[];
  conversationHistory: Message[];
}

interface Message {
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: number;
  uiComponents?: UIComponent[];  // 可选的UI组件
  configDelta?: Partial<CampaignConfig>;  // 配置增量
}
```

### 2. 意图分类与需求提取

```typescript
// LLM System Prompt 核心指令
const SYSTEM_PROMPT = `
你是TaskOn的增长顾问AI，帮助Web3项目方配置营销活动。

## 你的职责：
1. 理解用户的增长目标（获客/激活/留存/转化）
2. 推荐合适的游戏化组件组合
3. 引导用户完善配置细节
4. 生成标准化的配置JSON

## 对话流程：
Phase 1 - Discovery（发现需求）
- 询问项目类型（DEX/Perps/Lending/L2等）
- 了解当前生命周期阶段（冷启动/Pre-TGE/Post-listing等）
- 明确核心目标指标（UAW/TVL/Volume/Retention）

Phase 2 - Clarification（澄清细节）
- 确认目标用户画像
- 讨论预算和时间周期
- 推荐组件组合并解释原因

Phase 3 - Configuration（生成配置）
- 逐项确认配置参数
- 生成完整的JSON配置
- 提供配置说明

Phase 4 - Confirmation（确认发布）
- 展示配置预览
- 确认发布或修改

## 可用组件：
${JSON.stringify(AVAILABLE_WIDGETS, null, 2)}

## JSON Schema定义：
${JSON.stringify(CAMPAIGN_CONFIG_SCHEMA, null, 2)}

当用户需求明确时，使用以下格式输出配置：
\`\`\`json-config
{配置JSON}
\`\`\`
`;
```

### 3. Structured Output 配置

```typescript
// 使用 Claude API 的 tool_use 功能生成结构化配置
const tools = [
  {
    name: "generate_campaign_config",
    description: "生成完整的活动配置JSON",
    input_schema: {
      type: "object",
      properties: {
        campaign_type: {
          type: "string",
          enum: ["quest", "trading_race", "points_season", "referral"]
        },
        basic_info: {
          type: "object",
          properties: {
            name: { type: "string" },
            description: { type: "string" },
            start_time: { type: "string", format: "date-time" },
            end_time: { type: "string", format: "date-time" }
          },
          required: ["name", "start_time", "end_time"]
        },
        tasks: {
          type: "array",
          items: { "$ref": "#/definitions/TaskConfig" }
        },
        rewards: {
          type: "object",
          "$ref": "#/definitions/RewardConfig"
        },
        eligibility: {
          type: "object",
          "$ref": "#/definitions/EligibilityConfig"
        },
        widgets: {
          type: "array",
          items: { "$ref": "#/definitions/WidgetConfig" }
        }
      },
      required: ["campaign_type", "basic_info", "tasks", "rewards"]
    }
  },
  {
    name: "update_config_field",
    description: "更新配置的特定字段",
    input_schema: {
      type: "object",
      properties: {
        path: { type: "string", description: "JSON path, e.g., 'rewards.token.amount'" },
        value: { type: "any" }
      },
      required: ["path", "value"]
    }
  },
  {
    name: "suggest_widget_combination",
    description: "根据目标推荐组件组合",
    input_schema: {
      type: "object",
      properties: {
        goal: { type: "string", enum: ["acquisition", "activation", "retention", "monetization"] },
        project_type: { type: "string" },
        budget_range: { type: "string" }
      }
    }
  }
];
```

---

## 📦 JSON Schema 设计

### 核心配置结构

```typescript
// /schemas/campaign-config.schema.json

interface CampaignConfig {
  // 元数据
  $schema: string;
  version: string;
  
  // 基础信息
  campaign: {
    id?: string;
    type: 'quest' | 'trading_race' | 'points_season' | 'referral' | 'event';
    name: string;
    description?: string;
    cover_image?: string;
    visibility: 'public' | 'private';
    period: {
      start: string;  // ISO 8601
      end: string;
    };
  };
  
  // 任务配置
  tasks: TaskConfig[];
  
  // 奖励配置
  rewards: RewardConfig;
  
  // 准入门槛
  eligibility?: EligibilityConfig;
  
  // 游戏化组件
  widgets: WidgetConfig[];
  
  // 反女巫配置
  anti_sybil?: AntiSybilConfig;
  
  // 分析追踪
  analytics?: AnalyticsConfig;
}

interface TaskConfig {
  id: string;
  type: 'offchain' | 'onchain' | 'api' | 'poh';
  category: string;  // e.g., 'x', 'discord', 'swap', 'stake'
  template: string;  // 具体任务模板ID
  params: Record<string, any>;  // 模板参数
  points?: number;
  required: boolean;
  recurrence?: 'once' | 'daily' | 'weekly' | 'monthly';
  validation_timeframe?: {
    start?: string;
    end?: string;
  };
}

interface WidgetConfig {
  id: string;
  type: 'task_list' | 'task_chain' | 'day_chain' | 'leaderboard' | 
        'milestone' | 'benefits_shop' | 'lucky_wheel' | 'user_center';
  position: number;  // 排序位置
  config: Record<string, any>;  // 组件特定配置
  styling?: {
    theme?: 'light' | 'dark' | 'custom';
    custom_css?: string;
  };
}

interface RewardConfig {
  distribution_method: 'fcfs' | 'lucky_draw' | 'ranking' | 'open_to_all';
  rewards: {
    type: 'token' | 'nft' | 'points' | 'whitelist' | 'discord_role';
    // ... 具体奖励配置
  }[];
  auto_distribute: boolean;
}
```

### 组件配置Schema示例

```typescript
// /schemas/widgets/leaderboard.schema.json
interface LeaderboardWidgetConfig {
  type: 'leaderboard';
  config: {
    metric: 'points' | 'volume' | 'referrals' | 'custom';
    time_range: 'all' | '7d' | '30d' | 'season';
    display_count: number;
    show_user_rank: boolean;
    incentive_tiers?: {
      rank_range: [number, number];
      reward: RewardConfig;
    }[];
  };
}

// /schemas/widgets/task-chain.schema.json
interface TaskChainWidgetConfig {
  type: 'task_chain';
  config: {
    chain_id: string;
    name: string;
    tasks: {
      task_id: string;
      unlock_condition?: 'previous_completed' | 'time_based' | 'none';
      unlock_delay_hours?: number;
    }[];
    completion_reward?: RewardConfig;
  };
}
```

---

## 🔌 MCP 服务层设计

### 1. MCP Gateway 架构

```typescript
// /services/mcp-gateway/index.ts

import { Server } from '@modelcontextprotocol/sdk/server';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio';

// MCP Server 注册表
const MCP_SERVERS = {
  quest: QuestMCPServer,
  points: PointsMCPServer,
  leaderboard: LeaderboardMCPServer,
  rewards: RewardsMCPServer,
  analytics: AnalyticsMCPServer,
  antiSybil: AntiSybilMCPServer,
};

// Gateway 路由
class MCPGateway {
  private servers: Map<string, MCPServer>;
  
  async handleRequest(request: JSONRPCRequest): Promise<JSONRPCResponse> {
    const { method, params } = request;
    const [serverName, toolName] = method.split('.');
    
    const server = this.servers.get(serverName);
    if (!server) {
      return { error: { code: -32601, message: 'Server not found' } };
    }
    
    return server.callTool(toolName, params);
  }
}
```

### 2. Quest MCP Server 示例

```typescript
// /services/mcp-servers/quest-server.ts

import { Server } from '@modelcontextprotocol/sdk/server';

class QuestMCPServer {
  private server: Server;
  
  constructor() {
    this.server = new Server({
      name: 'taskon-quest-server',
      version: '1.0.0',
    }, {
      capabilities: {
        tools: {},
        resources: {},
      },
    });
    
    this.registerTools();
    this.registerResources();
  }
  
  private registerTools() {
    // 创建Quest
    this.server.setRequestHandler('tools/call', async (request) => {
      const { name, arguments: args } = request.params;
      
      switch (name) {
        case 'create_quest':
          return this.createQuest(args as CampaignConfig);
        case 'update_quest':
          return this.updateQuest(args.id, args.config);
        case 'publish_quest':
          return this.publishQuest(args.id);
        case 'get_quest_stats':
          return this.getQuestStats(args.id);
        default:
          throw new Error(`Unknown tool: ${name}`);
      }
    });
  }
  
  private registerResources() {
    // 注册可读资源
    this.server.setRequestHandler('resources/list', async () => ({
      resources: [
        {
          uri: 'quest://templates',
          name: 'Quest Templates',
          description: 'Available quest templates',
          mimeType: 'application/json',
        },
        {
          uri: 'quest://active',
          name: 'Active Quests',
          description: 'Currently active quests',
          mimeType: 'application/json',
        },
      ],
    }));
    
    this.server.setRequestHandler('resources/read', async (request) => {
      const { uri } = request.params;
      // 根据URI返回相应资源
    });
  }
  
  // Tool implementations
  private async createQuest(config: CampaignConfig): Promise<MCPToolResult> {
    // 验证配置
    const validation = await this.validateConfig(config);
    if (!validation.valid) {
      return {
        content: [{
          type: 'text',
          text: JSON.stringify({ error: validation.errors })
        }],
        isError: true,
      };
    }
    
    // 创建Quest
    const quest = await this.db.quest.create({ data: config });
    
    return {
      content: [{
        type: 'text',
        text: JSON.stringify({ 
          success: true, 
          quest_id: quest.id,
          preview_url: `https://taskon.xyz/quest/${quest.id}/preview`
        })
      }],
    };
  }
}
```

### 3. MCP Tools 完整定义

```typescript
// /services/mcp-servers/tools-definition.ts

const QUEST_TOOLS = [
  {
    name: 'create_quest',
    description: '创建新的Quest活动',
    inputSchema: {
      type: 'object',
      properties: {
        config: { $ref: '#/definitions/CampaignConfig' }
      },
      required: ['config']
    }
  },
  {
    name: 'update_quest',
    description: '更新Quest配置',
    inputSchema: {
      type: 'object',
      properties: {
        id: { type: 'string' },
        config: { $ref: '#/definitions/PartialCampaignConfig' }
      },
      required: ['id', 'config']
    }
  },
  {
    name: 'publish_quest',
    description: '发布Quest',
    inputSchema: {
      type: 'object',
      properties: {
        id: { type: 'string' },
        notify_community: { type: 'boolean', default: true }
      },
      required: ['id']
    }
  },
  {
    name: 'get_quest_stats',
    description: '获取Quest统计数据',
    inputSchema: {
      type: 'object',
      properties: {
        id: { type: 'string' },
        metrics: {
          type: 'array',
          items: {
            type: 'string',
            enum: ['participants', 'completions', 'conversion_rate', 'rewards_claimed']
          }
        }
      },
      required: ['id']
    }
  },
  {
    name: 'validate_task',
    description: '验证用户任务完成状态',
    inputSchema: {
      type: 'object',
      properties: {
        quest_id: { type: 'string' },
        task_id: { type: 'string' },
        user_wallet: { type: 'string' }
      },
      required: ['quest_id', 'task_id', 'user_wallet']
    }
  }
];

const POINTS_TOOLS = [
  {
    name: 'create_points_system',
    description: '创建积分体系',
    inputSchema: { ... }
  },
  {
    name: 'award_points',
    description: '发放积分',
    inputSchema: { ... }
  },
  {
    name: 'get_leaderboard',
    description: '获取排行榜',
    inputSchema: { ... }
  }
];

const REWARDS_TOOLS = [
  {
    name: 'setup_reward_pool',
    description: '设置奖励池',
    inputSchema: { ... }
  },
  {
    name: 'distribute_rewards',
    description: '分发奖励',
    inputSchema: { ... }
  }
];
```

---

## 🎨 Widget 组件层设计

### 1. Widget Registry & Renderer

```typescript
// /components/widgets/WidgetRegistry.ts

interface WidgetDefinition {
  type: string;
  component: React.ComponentType<any>;
  schema: JSONSchema;  // 配置schema
  defaultConfig: any;
}

const widgetRegistry: Map<string, WidgetDefinition> = new Map([
  ['task_list', {
    type: 'task_list',
    component: TaskListWidget,
    schema: taskListSchema,
    defaultConfig: { layout: 'vertical', showProgress: true }
  }],
  ['task_chain', {
    type: 'task_chain',
    component: TaskChainWidget,
    schema: taskChainSchema,
    defaultConfig: { direction: 'horizontal', showConnectors: true }
  }],
  ['day_chain', {
    type: 'day_chain',
    component: DayChainWidget,
    schema: dayChainSchema,
    defaultConfig: { streakRewards: true }
  }],
  ['leaderboard', {
    type: 'leaderboard',
    component: LeaderboardWidget,
    schema: leaderboardSchema,
    defaultConfig: { pageSize: 10, showCurrentUser: true }
  }],
  ['milestone', {
    type: 'milestone',
    component: MilestoneWidget,
    schema: milestoneSchema,
    defaultConfig: { style: 'progress_bar' }
  }],
  ['benefits_shop', {
    type: 'benefits_shop',
    component: BenefitsShopWidget,
    schema: benefitsShopSchema,
    defaultConfig: { categories: ['all'] }
  }],
  ['lucky_wheel', {
    type: 'lucky_wheel',
    component: LuckyWheelWidget,
    schema: luckyWheelSchema,
    defaultConfig: { spinCost: 100 }
  }],
  ['user_center', {
    type: 'user_center',
    component: UserCenterWidget,
    schema: userCenterSchema,
    defaultConfig: { showAssets: true, showHistory: true }
  }],
]);
```

### 2. Dynamic Widget Renderer

```tsx
// /components/widgets/WidgetRenderer.tsx

interface WidgetRendererProps {
  config: WidgetConfig[];
  mcpClient: MCPClient;
  theme?: ThemeConfig;
}

const WidgetRenderer: React.FC<WidgetRendererProps> = ({ 
  config, 
  mcpClient,
  theme 
}) => {
  return (
    <ThemeProvider theme={theme}>
      <WidgetContainer>
        {config
          .sort((a, b) => a.position - b.position)
          .map((widgetConfig) => {
            const definition = widgetRegistry.get(widgetConfig.type);
            if (!definition) {
              console.warn(`Unknown widget type: ${widgetConfig.type}`);
              return null;
            }
            
            const WidgetComponent = definition.component;
            
            return (
              <WidgetWrapper 
                key={widgetConfig.id}
                styling={widgetConfig.styling}
              >
                <WidgetComponent
                  config={widgetConfig.config}
                  mcpClient={mcpClient}
                />
              </WidgetWrapper>
            );
          })}
      </WidgetContainer>
    </ThemeProvider>
  );
};
```

### 3. Widget与MCP通信

```tsx
// /components/widgets/LeaderboardWidget.tsx

const LeaderboardWidget: React.FC<LeaderboardWidgetProps> = ({ 
  config, 
  mcpClient 
}) => {
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        // 通过MCP获取数据
        const result = await mcpClient.callTool('points.get_leaderboard', {
          metric: config.metric,
          time_range: config.time_range,
          limit: config.display_count,
        });
        
        setLeaderboard(JSON.parse(result.content[0].text).data);
      } catch (error) {
        console.error('Failed to fetch leaderboard:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchLeaderboard();
    
    // 设置实时更新
    const interval = setInterval(fetchLeaderboard, 30000);
    return () => clearInterval(interval);
  }, [config, mcpClient]);
  
  return (
    <LeaderboardContainer>
      {loading ? (
        <Skeleton rows={config.display_count} />
      ) : (
        <>
          {leaderboard.map((entry, index) => (
            <LeaderboardRow key={entry.userId} rank={index + 1} {...entry} />
          ))}
          {config.show_user_rank && <CurrentUserRank mcpClient={mcpClient} />}
        </>
      )}
    </LeaderboardContainer>
  );
};
```

---

## 🔄 前后端交互协议

### 1. JSON-RPC 2.0 消息格式

```typescript
// Request
interface MCPRequest {
  jsonrpc: "2.0";
  id: string | number;
  method: string;  // e.g., "quest.create_quest"
  params?: any;
}

// Response
interface MCPResponse {
  jsonrpc: "2.0";
  id: string | number;
  result?: MCPToolResult;
  error?: {
    code: number;
    message: string;
    data?: any;
  };
}

// Tool Result
interface MCPToolResult {
  content: {
    type: 'text' | 'image' | 'resource';
    text?: string;
    data?: string;  // base64 for images
    uri?: string;   // for resources
    mimeType?: string;
  }[];
  isError?: boolean;
}
```

### 2. 前端MCP Client

```typescript
// /services/mcp-client.ts

class TaskOnMCPClient {
  private ws: WebSocket;
  private pendingRequests: Map<string, { resolve: Function; reject: Function }>;
  private requestId: number = 0;
  
  constructor(endpoint: string) {
    this.ws = new WebSocket(endpoint);
    this.pendingRequests = new Map();
    
    this.ws.onmessage = (event) => {
      const response: MCPResponse = JSON.parse(event.data);
      const pending = this.pendingRequests.get(String(response.id));
      
      if (pending) {
        if (response.error) {
          pending.reject(response.error);
        } else {
          pending.resolve(response.result);
        }
        this.pendingRequests.delete(String(response.id));
      }
    };
  }
  
  async callTool(method: string, params?: any): Promise<MCPToolResult> {
    const id = ++this.requestId;
    
    return new Promise((resolve, reject) => {
      this.pendingRequests.set(String(id), { resolve, reject });
      
      const request: MCPRequest = {
        jsonrpc: "2.0",
        id,
        method,
        params,
      };
      
      this.ws.send(JSON.stringify(request));
      
      // Timeout
      setTimeout(() => {
        if (this.pendingRequests.has(String(id))) {
          this.pendingRequests.delete(String(id));
          reject(new Error('Request timeout'));
        }
      }, 30000);
    });
  }
  
  async getResource(uri: string): Promise<any> {
    return this.callTool('resources/read', { uri });
  }
}
```

---

## 📱 对话式UI设计

### 1. Chat Interface with Dynamic Components

```tsx
// /components/chat/ChatInterface.tsx

const ChatInterface: React.FC = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [currentConfig, setCurrentConfig] = useState<Partial<CampaignConfig>>({});
  const [isConfiguring, setIsConfiguring] = useState(false);
  
  const handleUserMessage = async (content: string) => {
    // 添加用户消息
    setMessages(prev => [...prev, { role: 'user', content }]);
    
    // 调用LLM
    const response = await callLLM({
      messages: [...messages, { role: 'user', content }],
      tools: CAMPAIGN_TOOLS,
      context: { currentConfig },
    });
    
    // 处理LLM响应
    const assistantMessage: ChatMessage = {
      role: 'assistant',
      content: response.textContent,
    };
    
    // 检查是否有tool调用
    if (response.toolCalls) {
      for (const toolCall of response.toolCalls) {
        if (toolCall.name === 'generate_campaign_config') {
          assistantMessage.configUpdate = toolCall.input;
          setCurrentConfig(prev => ({
            ...prev,
            ...toolCall.input,
          }));
        } else if (toolCall.name === 'suggest_widget_combination') {
          assistantMessage.uiComponents = [{
            type: 'widget_suggestion',
            data: toolCall.input,
          }];
        }
      }
    }
    
    // 检查是否有内嵌UI组件
    if (response.uiComponents) {
      assistantMessage.uiComponents = response.uiComponents;
    }
    
    setMessages(prev => [...prev, assistantMessage]);
  };
  
  return (
    <ChatContainer>
      <MessageList>
        {messages.map((msg, idx) => (
          <MessageBubble key={idx} role={msg.role}>
            <MessageContent content={msg.content} />
            {msg.uiComponents && (
              <DynamicComponents components={msg.uiComponents} />
            )}
            {msg.configUpdate && (
              <ConfigPreview config={msg.configUpdate} />
            )}
          </MessageBubble>
        ))}
      </MessageList>
      
      <ChatInput onSend={handleUserMessage} />
      
      {Object.keys(currentConfig).length > 0 && (
        <ConfigSidebar 
          config={currentConfig}
          onPublish={handlePublish}
          onEdit={(path, value) => {
            setCurrentConfig(prev => setDeep(prev, path, value));
          }}
        />
      )}
    </ChatContainer>
  );
};
```

### 2. 动态UI组件渲染

```tsx
// /components/chat/DynamicComponents.tsx

const DynamicComponents: React.FC<{ components: UIComponent[] }> = ({ 
  components 
}) => {
  return (
    <ComponentsContainer>
      {components.map((comp, idx) => {
        switch (comp.type) {
          case 'form':
            return <DynamicForm key={idx} {...comp.data} />;
          case 'buttons':
            return <ButtonGroup key={idx} options={comp.data.options} />;
          case 'widget_suggestion':
            return <WidgetSuggestionCard key={idx} {...comp.data} />;
          case 'config_preview':
            return <ConfigPreviewCard key={idx} config={comp.data} />;
          case 'template_selector':
            return <TemplateSelector key={idx} templates={comp.data.templates} />;
          default:
            return null;
        }
      })}
    </ComponentsContainer>
  );
};
```

---

## 🚀 实施路线图

### Phase 1: 基础架构 (4-6周)

| 周次 | 任务 | 交付物 |
|------|------|--------|
| 1-2 | MCP Gateway搭建 | MCP Server框架、JSON-RPC通信层 |
| 2-3 | 核心MCP Server开发 | Quest/Points/Rewards MCP Servers |
| 3-4 | Widget Registry | 基础Widget组件库 |
| 4-6 | 前端MCP Client | WebSocket通信、状态管理 |

### Phase 2: LLM对话层 (3-4周)

| 周次 | 任务 | 交付物 |
|------|------|--------|
| 1-2 | 对话管理系统 | Conversation Manager、意图分类 |
| 2-3 | 知识库构建 | 运营最佳实践、模板库 |
| 3-4 | Structured Output | JSON Schema定义、验证层 |

### Phase 3: 集成与优化 (3-4周)

| 周次 | 任务 | 交付物 |
|------|------|--------|
| 1-2 | 端到端集成 | 完整对话→配置→渲染流程 |
| 2-3 | 白标支持 | 主题系统、嵌入式组件 |
| 3-4 | 性能优化 | 缓存、预加载、错误处理 |

### Phase 4: 上线与迭代 (2周)

| 周次 | 任务 | 交付物 |
|------|------|--------|
| 1 | Beta测试 | 内部测试、Bug修复 |
| 2 | 灰度发布 | 监控、反馈收集 |

---

## 💡 关键技术决策

### 1. MCP vs REST API

| 维度 | MCP | REST API |
|------|-----|----------|
| LLM集成 | ✅ 原生支持 | ❌ 需要适配 |
| 实时通信 | ✅ 双向 | ⚠️ 需要轮询 |
| 标准化 | ✅ 行业标准 | ⚠️ 自定义 |
| 生态系统 | ✅ 快速扩展 | ❌ 孤岛 |

**建议**: 逐步迁移核心API到MCP，保留REST作为降级方案。

### 2. LLM选择

| 模型 | 优势 | 劣势 |
|------|------|------|
| Claude | Tool use强、中文好 | 成本较高 |
| GPT-4o | 生态成熟 | JSON结构有时不稳定 |
| Gemini | 长上下文 | 国内访问问题 |

**建议**: 主用Claude API，配置本地小模型作为意图分类前置。

### 3. 前端技术栈

```
React 18+ + TypeScript + Zustand + TailwindCSS + shadcn/ui
```

---

## 🔐 安全考虑

1. **MCP认证**: OAuth 2.1 + JWT
2. **配置验证**: JSON Schema双重验证（前端+后端）
3. **权限控制**: 基于角色的MCP Tool访问控制
4. **数据隔离**: 租户级别数据隔离
5. **审计日志**: 所有MCP调用记录

---

## 📈 预期收益

1. **降低学习成本**: 用户通过对话即可配置，无需理解复杂UI
2. **提升配置质量**: LLM推荐最佳实践组合
3. **加快上线速度**: 从小时级降到分钟级
4. **支持未来扩展**: MCP标准化便于接入更多AI能力
5. **增强白标能力**: Widget组件化支持深度定制

---

## 📚 参考资源

- [MCP官方文档](https://modelcontextprotocol.io/)
- [MCP SDK (TypeScript)](https://github.com/modelcontextprotocol/typescript-sdk)
- [Claude API Tool Use](https://docs.anthropic.com/en/docs/tool-use)
- [JSON Schema规范](https://json-schema.org/)
