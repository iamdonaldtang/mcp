# TaskOn Widget 组件库架构规范

## 📋 目录

1. [架构总览](#1-架构总览)
2. [组件注册系统](#2-组件注册系统)
3. [核心组件详细规范](#3-核心组件详细规范)
4. [组件通信机制](#4-组件通信机制)
5. [主题系统](#5-主题系统)
6. [MCP集成规范](#6-mcp集成规范)
7. [状态管理](#7-状态管理)
8. [性能优化](#8-性能优化)
9. [测试规范](#9-测试规范)
10. [组件开发指南](#10-组件开发指南)

---

## 1. 架构总览

### 1.1 Widget系统定位

```
┌─────────────────────────────────────────────────────────────────────┐
│                        TaskOn Widget System                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│  │   LLM层     │───▶│  配置JSON   │───▶│  Widget层   │             │
│  │ (生成配置)   │    │  (数据契约)  │    │ (渲染UI)    │             │
│  └─────────────┘    └─────────────┘    └─────────────┘             │
│         │                                     │                      │
│         │                                     │                      │
│         ▼                                     ▼                      │
│  ┌─────────────────────────────────────────────────────┐           │
│  │                    MCP Service Layer                 │           │
│  │     (Quest Server / Points Server / Rewards...)      │           │
│  └─────────────────────────────────────────────────────┘           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 1.2 设计原则

| 原则 | 说明 |
|------|------|
| **配置驱动** | Widget完全由JSON配置驱动，无需编码即可定制 |
| **组合优于继承** | 通过组合基础组件构建复杂功能 |
| **单一职责** | 每个Widget只负责一个明确的功能域 |
| **解耦通信** | 组件间通过事件总线通信，避免直接依赖 |
| **主题隔离** | 样式与逻辑分离，支持完全白标定制 |
| **渐进增强** | 核心功能不依赖JS，增强功能逐步加载 |

### 1.3 技术栈

```typescript
// 核心技术栈
const techStack = {
  framework: 'React 18+',
  language: 'TypeScript 5+',
  styling: 'TailwindCSS + CSS Variables',
  stateManagement: 'Zustand',
  dataFetching: 'TanStack Query + MCP Client',
  validation: 'Zod',
  testing: 'Vitest + Testing Library',
  documentation: 'Storybook',
};
```

### 1.4 目录结构

```
src/widgets/
├── core/                      # 核心基础设施
│   ├── WidgetRegistry.ts      # 组件注册中心
│   ├── WidgetRenderer.tsx     # 动态渲染器
│   ├── WidgetContext.tsx      # 全局上下文
│   ├── EventBus.ts            # 事件总线
│   └── types.ts               # 核心类型定义
│
├── primitives/                # 原子组件
│   ├── Button/
│   ├── Card/
│   ├── Badge/
│   ├── Progress/
│   ├── Avatar/
│   ├── Tooltip/
│   └── ...
│
├── task/                      # 任务类组件
│   ├── TaskList/
│   ├── TaskChain/
│   ├── DayChain/
│   └── TaskCard/
│
├── incentive/                 # 激励类组件
│   ├── Leaderboard/
│   ├── Milestone/
│   ├── LuckyWheel/
│   └── BenefitsShop/
│
├── user/                      # 用户类组件
│   ├── UserCenter/
│   ├── ProfileCard/
│   ├── AssetPanel/
│   └── HistoryList/
│
├── layout/                    # 布局组件
│   ├── PageBuilder/
│   ├── Section/
│   ├── Grid/
│   └── Container/
│
├── hooks/                     # 共享Hooks
│   ├── useMCP.ts
│   ├── useWidget.ts
│   ├── useTheme.ts
│   └── useAuth.ts
│
├── utils/                     # 工具函数
│   ├── config.ts
│   ├── validation.ts
│   └── format.ts
│
└── themes/                    # 主题定义
    ├── default.ts
    ├── dark.ts
    └── types.ts
```

---

## 2. 组件注册系统

### 2.1 Widget注册中心

```typescript
// src/widgets/core/WidgetRegistry.ts

import { z } from 'zod';

/**
 * Widget元数据定义
 */
interface WidgetMetadata {
  /** 唯一标识 */
  type: string;
  /** 显示名称 */
  name: string;
  /** 组件描述 */
  description: string;
  /** 组件分类 */
  category: 'task' | 'incentive' | 'user' | 'layout' | 'primitive';
  /** 组件图标 */
  icon: string;
  /** 版本号 */
  version: string;
  /** 是否支持白标 */
  whitelabel: boolean;
  /** 依赖的MCP服务 */
  mcpDependencies: string[];
  /** 支持的配置项 */
  configSchema: z.ZodSchema;
  /** 默认配置 */
  defaultConfig: Record<string, unknown>;
}

/**
 * Widget注册项
 */
interface WidgetRegistration<P = unknown> {
  metadata: WidgetMetadata;
  component: React.ComponentType<P>;
  /** 配置面板组件（用于可视化编辑） */
  configPanel?: React.ComponentType<{ config: P; onChange: (config: P) => void }>;
  /** 预览组件（用于缩略图展示） */
  preview?: React.ComponentType<{ config: P }>;
}

/**
 * Widget注册中心 - 单例模式
 */
class WidgetRegistry {
  private static instance: WidgetRegistry;
  private widgets: Map<string, WidgetRegistration> = new Map();
  private categories: Map<string, Set<string>> = new Map();

  private constructor() {
    // 初始化分类
    ['task', 'incentive', 'user', 'layout', 'primitive'].forEach(cat => {
      this.categories.set(cat, new Set());
    });
  }

  static getInstance(): WidgetRegistry {
    if (!WidgetRegistry.instance) {
      WidgetRegistry.instance = new WidgetRegistry();
    }
    return WidgetRegistry.instance;
  }

  /**
   * 注册Widget
   */
  register<P>(registration: WidgetRegistration<P>): void {
    const { type, category } = registration.metadata;
    
    if (this.widgets.has(type)) {
      console.warn(`Widget "${type}" already registered, overwriting...`);
    }

    this.widgets.set(type, registration as WidgetRegistration);
    this.categories.get(category)?.add(type);

    console.log(`[WidgetRegistry] Registered: ${type} (${category})`);
  }

  /**
   * 获取Widget
   */
  get(type: string): WidgetRegistration | undefined {
    return this.widgets.get(type);
  }

  /**
   * 获取Widget组件
   */
  getComponent(type: string): React.ComponentType | null {
    return this.widgets.get(type)?.component || null;
  }

  /**
   * 获取分类下所有Widget
   */
  getByCategory(category: string): WidgetRegistration[] {
    const types = this.categories.get(category);
    if (!types) return [];
    return Array.from(types)
      .map(type => this.widgets.get(type))
      .filter(Boolean) as WidgetRegistration[];
  }

  /**
   * 获取所有Widget元数据（用于LLM知识库）
   */
  getAllMetadata(): WidgetMetadata[] {
    return Array.from(this.widgets.values()).map(w => w.metadata);
  }

  /**
   * 验证配置
   */
  validateConfig(type: string, config: unknown): { success: boolean; errors?: z.ZodError } {
    const widget = this.widgets.get(type);
    if (!widget) {
      return { success: false, errors: undefined };
    }

    const result = widget.metadata.configSchema.safeParse(config);
    return {
      success: result.success,
      errors: result.success ? undefined : result.error,
    };
  }

  /**
   * 获取默认配置
   */
  getDefaultConfig(type: string): Record<string, unknown> | null {
    return this.widgets.get(type)?.metadata.defaultConfig || null;
  }
}

export const widgetRegistry = WidgetRegistry.getInstance();
```

### 2.2 Widget注册装饰器

```typescript
// src/widgets/core/decorators.ts

import { widgetRegistry } from './WidgetRegistry';
import type { WidgetMetadata } from './WidgetRegistry';

/**
 * Widget注册装饰器
 * @example
 * @registerWidget({
 *   type: 'task_list',
 *   name: 'Task List',
 *   category: 'task',
 *   ...
 * })
 * export const TaskList: React.FC<TaskListProps> = (props) => { ... }
 */
export function registerWidget(metadata: WidgetMetadata) {
  return function <P>(Component: React.ComponentType<P>) {
    widgetRegistry.register({
      metadata,
      component: Component,
    });
    return Component;
  };
}

/**
 * 批量注册工具
 */
export function registerWidgets(
  widgets: Array<{
    metadata: WidgetMetadata;
    component: React.ComponentType;
    configPanel?: React.ComponentType;
    preview?: React.ComponentType;
  }>
) {
  widgets.forEach(widget => widgetRegistry.register(widget));
}
```

### 2.3 动态Widget渲染器

```typescript
// src/widgets/core/WidgetRenderer.tsx

import React, { Suspense, useMemo } from 'react';
import { widgetRegistry } from './WidgetRegistry';
import { WidgetErrorBoundary } from './WidgetErrorBoundary';
import { WidgetSkeleton } from './WidgetSkeleton';

/**
 * Widget配置接口
 */
export interface WidgetConfig {
  /** Widget类型 */
  type: string;
  /** Widget实例ID */
  id: string;
  /** 显示顺序 */
  position: number;
  /** 配置参数 */
  config: Record<string, unknown>;
  /** 可见性 */
  visible?: boolean;
  /** 样式覆盖 */
  style?: React.CSSProperties;
  /** 类名覆盖 */
  className?: string;
}

interface WidgetRendererProps {
  /** Widget配置数组 */
  widgets: WidgetConfig[];
  /** 布局模式 */
  layout?: 'vertical' | 'grid' | 'tabs';
  /** 全局上下文 */
  context?: Record<string, unknown>;
  /** 空状态渲染 */
  emptyState?: React.ReactNode;
}

/**
 * Widget动态渲染器
 * 根据配置数组动态渲染Widget组件
 */
export const WidgetRenderer: React.FC<WidgetRendererProps> = ({
  widgets,
  layout = 'vertical',
  context = {},
  emptyState,
}) => {
  // 按position排序并过滤不可见的widget
  const sortedWidgets = useMemo(() => {
    return widgets
      .filter(w => w.visible !== false)
      .sort((a, b) => a.position - b.position);
  }, [widgets]);

  if (sortedWidgets.length === 0) {
    return <>{emptyState || <DefaultEmptyState />}</>;
  }

  const layoutClass = {
    vertical: 'flex flex-col gap-6',
    grid: 'grid grid-cols-1 md:grid-cols-2 gap-6',
    tabs: '', // tabs布局单独处理
  }[layout];

  if (layout === 'tabs') {
    return <TabsLayout widgets={sortedWidgets} context={context} />;
  }

  return (
    <div className={layoutClass}>
      {sortedWidgets.map(widgetConfig => (
        <SingleWidgetRenderer
          key={widgetConfig.id}
          widgetConfig={widgetConfig}
          context={context}
        />
      ))}
    </div>
  );
};

/**
 * 单个Widget渲染器
 */
const SingleWidgetRenderer: React.FC<{
  widgetConfig: WidgetConfig;
  context: Record<string, unknown>;
}> = ({ widgetConfig, context }) => {
  const { type, id, config, style, className } = widgetConfig;

  // 获取注册的组件
  const Component = widgetRegistry.getComponent(type);

  if (!Component) {
    return (
      <WidgetNotFound type={type} id={id} />
    );
  }

  // 合并配置和上下文
  const mergedProps = {
    ...config,
    widgetId: id,
    widgetContext: context,
  };

  return (
    <WidgetErrorBoundary
      widgetId={id}
      widgetType={type}
      fallback={<WidgetError type={type} id={id} />}
    >
      <Suspense fallback={<WidgetSkeleton type={type} />}>
        <div style={style} className={className} data-widget-id={id} data-widget-type={type}>
          <Component {...mergedProps} />
        </div>
      </Suspense>
    </WidgetErrorBoundary>
  );
};

/**
 * Widget错误边界
 */
class WidgetErrorBoundary extends React.Component<{
  widgetId: string;
  widgetType: string;
  fallback: React.ReactNode;
  children: React.ReactNode;
}> {
  state = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error(`[Widget Error] ${this.props.widgetType}:${this.props.widgetId}`, error, errorInfo);
    // 可以上报到监控系统
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children;
  }
}

// 辅助组件
const DefaultEmptyState = () => (
  <div className="text-center py-12 text-gray-500">
    <p>No widgets configured</p>
  </div>
);

const WidgetNotFound: React.FC<{ type: string; id: string }> = ({ type, id }) => (
  <div className="p-4 border border-dashed border-red-300 rounded-lg bg-red-50">
    <p className="text-red-600">Widget not found: {type}</p>
    <p className="text-sm text-red-400">ID: {id}</p>
  </div>
);

const WidgetError: React.FC<{ type: string; id: string }> = ({ type, id }) => (
  <div className="p-4 border border-dashed border-yellow-300 rounded-lg bg-yellow-50">
    <p className="text-yellow-600">Widget error: {type}</p>
    <p className="text-sm text-yellow-400">Please try refreshing the page</p>
  </div>
);
```

---

## 3. 核心组件详细规范

### 3.1 TaskList (任务列表)

```typescript
// src/widgets/task/TaskList/types.ts

import { z } from 'zod';

/**
 * 单个任务定义
 */
export const TaskItemSchema = z.object({
  id: z.string(),
  type: z.enum(['offchain', 'onchain', 'api']),
  template: z.string(),
  title: z.string(),
  description: z.string().optional(),
  icon: z.string().optional(),
  points: z.number().default(0),
  required: z.boolean().default(false),
  status: z.enum(['locked', 'available', 'in_progress', 'completed', 'failed']).default('available'),
  params: z.record(z.unknown()).optional(),
  verification: z.object({
    type: z.enum(['auto', 'manual', 'api']),
    endpoint: z.string().optional(),
  }).optional(),
  progress: z.object({
    current: z.number(),
    target: z.number(),
  }).optional(),
  deadline: z.string().datetime().optional(),
  cooldown: z.number().optional(), // 秒
  maxClaims: z.number().optional(),
  claimedCount: z.number().default(0),
});

export type TaskItem = z.infer<typeof TaskItemSchema>;

/**
 * TaskList配置Schema
 */
export const TaskListConfigSchema = z.object({
  // 数据配置
  tasks: z.array(TaskItemSchema).default([]),
  
  // 布局配置
  layout: z.enum(['vertical', 'horizontal', 'grid']).default('vertical'),
  columns: z.number().min(1).max(4).default(1),
  spacing: z.enum(['compact', 'normal', 'loose']).default('normal'),
  
  // 显示配置
  showProgress: z.boolean().default(true),
  showPoints: z.boolean().default(true),
  showStatus: z.boolean().default(true),
  showDescription: z.boolean().default(true),
  showDeadline: z.boolean().default(true),
  
  // 分组配置
  groupBy: z.enum(['none', 'type', 'status', 'category']).default('none'),
  groupLabels: z.record(z.string()).optional(),
  
  // 排序配置
  sortBy: z.enum(['position', 'points', 'status', 'deadline']).default('position'),
  sortOrder: z.enum(['asc', 'desc']).default('asc'),
  
  // 过滤配置
  filterStatus: z.array(z.string()).optional(),
  filterType: z.array(z.string()).optional(),
  
  // 交互配置
  enableClick: z.boolean().default(true),
  enableVerify: z.boolean().default(true),
  autoRefresh: z.boolean().default(false),
  refreshInterval: z.number().default(30000), // ms
  
  // 空状态配置
  emptyState: z.object({
    title: z.string().default('No tasks available'),
    description: z.string().optional(),
    icon: z.string().optional(),
  }).optional(),
  
  // 样式配置
  cardStyle: z.enum(['default', 'minimal', 'elevated']).default('default'),
  accentColor: z.string().optional(),
});

export type TaskListConfig = z.infer<typeof TaskListConfigSchema>;

/**
 * TaskList Props
 */
export interface TaskListProps extends TaskListConfig {
  /** Widget实例ID */
  widgetId: string;
  /** 全局上下文 */
  widgetContext?: Record<string, unknown>;
  /** 任务点击回调 */
  onTaskClick?: (task: TaskItem) => void;
  /** 验证完成回调 */
  onVerifyComplete?: (task: TaskItem, success: boolean) => void;
  /** 自定义任务卡片渲染 */
  renderTask?: (task: TaskItem, defaultRender: React.ReactNode) => React.ReactNode;
}

/**
 * TaskList 事件定义
 */
export interface TaskListEvents {
  'task:click': { taskId: string; task: TaskItem };
  'task:verify:start': { taskId: string };
  'task:verify:success': { taskId: string; points: number };
  'task:verify:failed': { taskId: string; error: string };
  'task:claim': { taskId: string };
  'tasklist:refresh': { widgetId: string };
}
```

```typescript
// src/widgets/task/TaskList/TaskList.tsx

import React, { useMemo, useCallback } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { registerWidget } from '../../core/decorators';
import { useEventBus } from '../../core/EventBus';
import { useMCPClient } from '../../hooks/useMCP';
import { useTheme } from '../../hooks/useTheme';
import { TaskCard } from './TaskCard';
import { TaskListSkeleton } from './TaskListSkeleton';
import { TaskListEmpty } from './TaskListEmpty';
import { TaskListConfigSchema, type TaskListProps, type TaskItem } from './types';

const WIDGET_METADATA = {
  type: 'task_list',
  name: 'Task List',
  description: '展示任务列表，支持多种布局和分组方式',
  category: 'task' as const,
  icon: 'CheckSquare',
  version: '1.0.0',
  whitelabel: true,
  mcpDependencies: ['quest-server'],
  configSchema: TaskListConfigSchema,
  defaultConfig: {
    layout: 'vertical',
    showProgress: true,
    showPoints: true,
    groupBy: 'none',
    sortBy: 'position',
    cardStyle: 'default',
  },
};

export const TaskList: React.FC<TaskListProps> = registerWidget(WIDGET_METADATA)((props) => {
  const {
    widgetId,
    tasks: initialTasks = [],
    layout,
    columns,
    spacing,
    showProgress,
    showPoints,
    showStatus,
    showDescription,
    groupBy,
    sortBy,
    sortOrder,
    filterStatus,
    filterType,
    enableClick,
    enableVerify,
    autoRefresh,
    refreshInterval,
    emptyState,
    cardStyle,
    accentColor,
    onTaskClick,
    onVerifyComplete,
    renderTask,
  } = props;

  const { emit } = useEventBus();
  const { client: mcpClient } = useMCPClient('quest-server');
  const { theme, resolveColor } = useTheme();

  // 获取任务数据
  const {
    data: tasks,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['tasks', widgetId],
    queryFn: async () => {
      if (initialTasks.length > 0) return initialTasks;
      // 从MCP获取任务
      const result = await mcpClient.callTool('get_tasks', { widgetId });
      return result.tasks as TaskItem[];
    },
    refetchInterval: autoRefresh ? refreshInterval : false,
  });

  // 验证任务
  const verifyMutation = useMutation({
    mutationFn: async (taskId: string) => {
      emit('task:verify:start', { taskId });
      const result = await mcpClient.callTool('verify_task', { taskId });
      return result;
    },
    onSuccess: (result, taskId) => {
      if (result.success) {
        emit('task:verify:success', { taskId, points: result.points });
        onVerifyComplete?.(tasks?.find(t => t.id === taskId)!, true);
      } else {
        emit('task:verify:failed', { taskId, error: result.error });
        onVerifyComplete?.(tasks?.find(t => t.id === taskId)!, false);
      }
      refetch();
    },
  });

  // 处理任务点击
  const handleTaskClick = useCallback((task: TaskItem) => {
    if (!enableClick) return;
    emit('task:click', { taskId: task.id, task });
    onTaskClick?.(task);
  }, [enableClick, emit, onTaskClick]);

  // 处理验证
  const handleVerify = useCallback((task: TaskItem) => {
    if (!enableVerify) return;
    verifyMutation.mutate(task.id);
  }, [enableVerify, verifyMutation]);

  // 过滤任务
  const filteredTasks = useMemo(() => {
    if (!tasks) return [];
    return tasks.filter(task => {
      if (filterStatus?.length && !filterStatus.includes(task.status)) return false;
      if (filterType?.length && !filterType.includes(task.type)) return false;
      return true;
    });
  }, [tasks, filterStatus, filterType]);

  // 排序任务
  const sortedTasks = useMemo(() => {
    const sorted = [...filteredTasks].sort((a, b) => {
      let comparison = 0;
      switch (sortBy) {
        case 'points':
          comparison = a.points - b.points;
          break;
        case 'status':
          comparison = a.status.localeCompare(b.status);
          break;
        case 'deadline':
          comparison = (a.deadline || '').localeCompare(b.deadline || '');
          break;
        default:
          comparison = 0;
      }
      return sortOrder === 'desc' ? -comparison : comparison;
    });
    return sorted;
  }, [filteredTasks, sortBy, sortOrder]);

  // 分组任务
  const groupedTasks = useMemo(() => {
    if (groupBy === 'none') return { default: sortedTasks };
    
    return sortedTasks.reduce((groups, task) => {
      const key = task[groupBy as keyof TaskItem] as string || 'other';
      if (!groups[key]) groups[key] = [];
      groups[key].push(task);
      return groups;
    }, {} as Record<string, TaskItem[]>);
  }, [sortedTasks, groupBy]);

  // 布局样式
  const layoutStyles = useMemo(() => {
    const spacingMap = { compact: 'gap-2', normal: 'gap-4', loose: 'gap-6' };
    const baseSpacing = spacingMap[spacing || 'normal'];

    switch (layout) {
      case 'horizontal':
        return `flex flex-row overflow-x-auto ${baseSpacing}`;
      case 'grid':
        return `grid grid-cols-${columns || 1} md:grid-cols-${Math.min(columns || 1, 2)} lg:grid-cols-${columns || 1} ${baseSpacing}`;
      default:
        return `flex flex-col ${baseSpacing}`;
    }
  }, [layout, columns, spacing]);

  // 加载状态
  if (isLoading) {
    return <TaskListSkeleton layout={layout} count={3} />;
  }

  // 空状态
  if (!sortedTasks.length) {
    return <TaskListEmpty {...emptyState} />;
  }

  // 渲染样式变量
  const styleVars = accentColor ? {
    '--widget-accent': resolveColor(accentColor),
  } as React.CSSProperties : undefined;

  return (
    <div 
      className="taskon-task-list"
      style={styleVars}
      data-widget-id={widgetId}
    >
      {Object.entries(groupedTasks).map(([groupKey, groupTasks]) => (
        <div key={groupKey} className="task-group">
          {groupBy !== 'none' && (
            <h3 className="text-sm font-medium text-gray-500 mb-3 uppercase tracking-wider">
              {props.groupLabels?.[groupKey] || groupKey}
            </h3>
          )}
          
          <div className={layoutStyles}>
            {groupTasks.map(task => {
              const cardNode = (
                <TaskCard
                  key={task.id}
                  task={task}
                  variant={cardStyle}
                  showProgress={showProgress}
                  showPoints={showPoints}
                  showStatus={showStatus}
                  showDescription={showDescription}
                  onClick={() => handleTaskClick(task)}
                  onVerify={() => handleVerify(task)}
                  isVerifying={verifyMutation.isPending && verifyMutation.variables === task.id}
                />
              );
              
              return renderTask ? renderTask(task, cardNode) : cardNode;
            })}
          </div>
        </div>
      ))}
    </div>
  );
});
```

### 3.2 Leaderboard (排行榜)

```typescript
// src/widgets/incentive/Leaderboard/types.ts

import { z } from 'zod';

/**
 * 排行榜项
 */
export const LeaderboardEntrySchema = z.object({
  rank: z.number(),
  userId: z.string(),
  address: z.string(),
  displayName: z.string().optional(),
  avatar: z.string().optional(),
  score: z.number(),
  change: z.number().optional(), // 排名变化
  isCurrentUser: z.boolean().default(false),
  metadata: z.record(z.unknown()).optional(),
});

export type LeaderboardEntry = z.infer<typeof LeaderboardEntrySchema>;

/**
 * 奖励层级
 */
export const RewardTierSchema = z.object({
  rankRange: z.tuple([z.number(), z.number()]),
  reward: z.object({
    type: z.enum(['token', 'nft', 'points', 'badge']),
    amount: z.number().optional(),
    name: z.string(),
    icon: z.string().optional(),
  }),
  highlight: z.boolean().default(false),
  label: z.string().optional(),
});

export type RewardTier = z.infer<typeof RewardTierSchema>;

/**
 * Leaderboard配置Schema
 */
export const LeaderboardConfigSchema = z.object({
  // 数据配置
  metric: z.enum(['points', 'volume', 'referrals', 'trades', 'custom']).default('points'),
  metricLabel: z.string().default('Score'),
  metricFormatter: z.string().optional(), // 格式化模板，如 "${value} pts"
  
  // 时间范围
  timeRange: z.enum(['all', '24h', '7d', '30d', 'season', 'custom']).default('all'),
  customStartTime: z.string().datetime().optional(),
  customEndTime: z.string().datetime().optional(),
  
  // 显示配置
  displayCount: z.number().min(3).max(100).default(10),
  showRankChange: z.boolean().default(true),
  showRewardTiers: z.boolean().default(true),
  showCurrentUser: z.boolean().default(true),
  showAvatar: z.boolean().default(true),
  showAddress: z.boolean().default(true),
  addressFormat: z.enum(['full', 'short', 'ens']).default('short'),
  
  // 奖励配置
  rewardTiers: z.array(RewardTierSchema).optional(),
  
  // 交互配置
  enableUserClick: z.boolean().default(true),
  autoRefresh: z.boolean().default(true),
  refreshInterval: z.number().default(60000),
  
  // 动画配置
  animateChanges: z.boolean().default(true),
  highlightDuration: z.number().default(3000),
  
  // 样式配置
  variant: z.enum(['default', 'compact', 'detailed', 'podium']).default('default'),
  accentColor: z.string().optional(),
  highlightTopN: z.number().default(3),
});

export type LeaderboardConfig = z.infer<typeof LeaderboardConfigSchema>;

export interface LeaderboardProps extends LeaderboardConfig {
  widgetId: string;
  widgetContext?: Record<string, unknown>;
  onUserClick?: (entry: LeaderboardEntry) => void;
  onTimeRangeChange?: (range: string) => void;
}
```

```typescript
// src/widgets/incentive/Leaderboard/Leaderboard.tsx

import React, { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion, AnimatePresence } from 'framer-motion';
import { registerWidget } from '../../core/decorators';
import { useMCPClient } from '../../hooks/useMCP';
import { useCurrentUser } from '../../hooks/useAuth';
import { formatAddress, formatNumber } from '../../utils/format';
import { LeaderboardConfigSchema, type LeaderboardProps, type LeaderboardEntry } from './types';

// ... 组件实现

const WIDGET_METADATA = {
  type: 'leaderboard',
  name: 'Leaderboard',
  description: '展示用户排名，支持多种指标和时间范围',
  category: 'incentive' as const,
  icon: 'Trophy',
  version: '1.0.0',
  whitelabel: true,
  mcpDependencies: ['points-server'],
  configSchema: LeaderboardConfigSchema,
  defaultConfig: {
    metric: 'points',
    displayCount: 10,
    showRankChange: true,
    showRewardTiers: true,
    variant: 'default',
  },
};

export const Leaderboard: React.FC<LeaderboardProps> = registerWidget(WIDGET_METADATA)((props) => {
  const {
    widgetId,
    metric,
    metricLabel,
    metricFormatter,
    timeRange,
    displayCount,
    showRankChange,
    showRewardTiers,
    showCurrentUser,
    showAvatar,
    showAddress,
    addressFormat,
    rewardTiers,
    enableUserClick,
    autoRefresh,
    refreshInterval,
    animateChanges,
    variant,
    accentColor,
    highlightTopN,
    onUserClick,
    onTimeRangeChange,
  } = props;

  const [selectedTimeRange, setSelectedTimeRange] = useState(timeRange);
  const { client: mcpClient } = useMCPClient('points-server');
  const { user: currentUser } = useCurrentUser();

  // 获取排行榜数据
  const { data: leaderboard, isLoading } = useQuery({
    queryKey: ['leaderboard', widgetId, metric, selectedTimeRange],
    queryFn: async () => {
      const result = await mcpClient.callTool('get_leaderboard', {
        metric,
        timeRange: selectedTimeRange,
        limit: displayCount,
        includeCurrentUser: showCurrentUser,
        currentUserId: currentUser?.id,
      });
      return result.entries as LeaderboardEntry[];
    },
    refetchInterval: autoRefresh ? refreshInterval : false,
  });

  // 格式化分数
  const formatScore = (score: number) => {
    if (metricFormatter) {
      return metricFormatter.replace('${value}', formatNumber(score));
    }
    return formatNumber(score);
  };

  // 获取奖励层级
  const getRewardForRank = (rank: number) => {
    if (!rewardTiers) return null;
    return rewardTiers.find(tier => 
      rank >= tier.rankRange[0] && rank <= tier.rankRange[1]
    );
  };

  // 时间范围选择器
  const TimeRangeSelector = () => (
    <div className="flex gap-2 mb-4">
      {['24h', '7d', '30d', 'all'].map(range => (
        <button
          key={range}
          className={`px-3 py-1 rounded-full text-sm ${
            selectedTimeRange === range 
              ? 'bg-primary text-white' 
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          }`}
          onClick={() => {
            setSelectedTimeRange(range as any);
            onTimeRangeChange?.(range);
          }}
        >
          {range === 'all' ? 'All Time' : range.toUpperCase()}
        </button>
      ))}
    </div>
  );

  // 领奖台视图 (Top 3)
  const PodiumView = () => {
    if (!leaderboard || leaderboard.length < 3) return null;
    
    const [first, second, third] = leaderboard;
    
    return (
      <div className="flex justify-center items-end gap-4 mb-6 h-40">
        {/* 第二名 */}
        <PodiumItem entry={second} rank={2} height="h-28" />
        {/* 第一名 */}
        <PodiumItem entry={first} rank={1} height="h-36" />
        {/* 第三名 */}
        <PodiumItem entry={third} rank={3} height="h-24" />
      </div>
    );
  };

  const PodiumItem: React.FC<{ entry: LeaderboardEntry; rank: number; height: string }> = ({ 
    entry, rank, height 
  }) => {
    const colors = {
      1: 'bg-yellow-400 text-yellow-900',
      2: 'bg-gray-300 text-gray-700',
      3: 'bg-orange-400 text-orange-900',
    };
    
    return (
      <motion.div
        initial={{ y: 50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: (3 - rank) * 0.2 }}
        className={`flex flex-col items-center ${height} w-24`}
      >
        <div className="relative">
          <img 
            src={entry.avatar || '/default-avatar.png'} 
            className="w-12 h-12 rounded-full border-2 border-white shadow-lg"
            alt={entry.displayName}
          />
          <span className={`absolute -bottom-1 -right-1 w-5 h-5 rounded-full ${colors[rank as 1|2|3]} flex items-center justify-center text-xs font-bold`}>
            {rank}
          </span>
        </div>
        <p className="text-sm font-medium mt-2 truncate w-full text-center">
          {entry.displayName || formatAddress(entry.address, addressFormat)}
        </p>
        <p className="text-xs text-gray-500">{formatScore(entry.score)}</p>
      </motion.div>
    );
  };

  // 列表项
  const LeaderboardRow: React.FC<{ entry: LeaderboardEntry; index: number }> = ({ entry, index }) => {
    const reward = getRewardForRank(entry.rank);
    const isHighlighted = entry.rank <= highlightTopN || entry.isCurrentUser;
    
    return (
      <motion.div
        layout={animateChanges}
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: index * 0.05 }}
        className={`
          flex items-center gap-4 p-3 rounded-lg
          ${entry.isCurrentUser ? 'bg-primary/10 border border-primary' : 'hover:bg-gray-50'}
          ${isHighlighted ? 'font-medium' : ''}
          ${enableUserClick ? 'cursor-pointer' : ''}
        `}
        onClick={() => enableUserClick && onUserClick?.(entry)}
      >
        {/* 排名 */}
        <div className={`
          w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold
          ${entry.rank === 1 ? 'bg-yellow-400 text-yellow-900' : ''}
          ${entry.rank === 2 ? 'bg-gray-300 text-gray-700' : ''}
          ${entry.rank === 3 ? 'bg-orange-400 text-orange-900' : ''}
          ${entry.rank > 3 ? 'bg-gray-100 text-gray-600' : ''}
        `}>
          {entry.rank}
        </div>
        
        {/* 头像 */}
        {showAvatar && (
          <img 
            src={entry.avatar || '/default-avatar.png'}
            className="w-10 h-10 rounded-full"
            alt=""
          />
        )}
        
        {/* 用户信息 */}
        <div className="flex-1 min-w-0">
          <p className="font-medium truncate">
            {entry.displayName || formatAddress(entry.address, addressFormat)}
          </p>
          {showAddress && entry.displayName && (
            <p className="text-xs text-gray-400">
              {formatAddress(entry.address, 'short')}
            </p>
          )}
        </div>
        
        {/* 排名变化 */}
        {showRankChange && entry.change !== undefined && (
          <div className={`text-sm ${entry.change > 0 ? 'text-green-500' : entry.change < 0 ? 'text-red-500' : 'text-gray-400'}`}>
            {entry.change > 0 ? `↑${entry.change}` : entry.change < 0 ? `↓${Math.abs(entry.change)}` : '-'}
          </div>
        )}
        
        {/* 分数 */}
        <div className="text-right">
          <p className="font-semibold">{formatScore(entry.score)}</p>
          {reward && showRewardTiers && (
            <p className="text-xs text-primary">{reward.reward.name}</p>
          )}
        </div>
      </motion.div>
    );
  };

  if (isLoading) {
    return <LeaderboardSkeleton count={displayCount} variant={variant} />;
  }

  return (
    <div className="taskon-leaderboard" data-widget-id={widgetId}>
      {/* 标题栏 */}
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold flex items-center gap-2">
          🏆 {metricLabel} Leaderboard
        </h3>
        <TimeRangeSelector />
      </div>
      
      {/* 领奖台 (仅在 podium 变体) */}
      {variant === 'podium' && <PodiumView />}
      
      {/* 奖励层级说明 */}
      {showRewardTiers && rewardTiers && (
        <div className="flex gap-2 mb-4 flex-wrap">
          {rewardTiers.map((tier, i) => (
            <span 
              key={i}
              className="px-2 py-1 bg-gray-100 rounded text-xs"
            >
              #{tier.rankRange[0]}-{tier.rankRange[1]}: {tier.reward.name}
            </span>
          ))}
        </div>
      )}
      
      {/* 排行榜列表 */}
      <div className="space-y-2">
        <AnimatePresence>
          {leaderboard?.slice(variant === 'podium' ? 3 : 0).map((entry, index) => (
            <LeaderboardRow 
              key={entry.userId} 
              entry={entry} 
              index={index}
            />
          ))}
        </AnimatePresence>
      </div>
      
      {/* 当前用户（如果不在列表中） */}
      {showCurrentUser && currentUser && !leaderboard?.some(e => e.isCurrentUser) && (
        <div className="mt-4 pt-4 border-t border-dashed">
          <p className="text-sm text-gray-500 mb-2">Your Rank</p>
          {/* 显示当前用户排名 */}
        </div>
      )}
    </div>
  );
});
```

### 3.3 DayChain (多日挑战)

```typescript
// src/widgets/task/DayChain/types.ts

import { z } from 'zod';
import { TaskItemSchema } from '../TaskList/types';

/**
 * 单日任务配置
 */
export const DayConfigSchema = z.object({
  day: z.number().min(1),
  tasks: z.array(TaskItemSchema),
  reward: z.object({
    type: z.enum(['points', 'token', 'nft', 'badge']),
    amount: z.number().optional(),
    name: z.string(),
  }).optional(),
  unlockTime: z.string().datetime().optional(), // 绝对解锁时间
  status: z.enum(['locked', 'available', 'completed', 'missed']).default('locked'),
});

export type DayConfig = z.infer<typeof DayConfigSchema>;

/**
 * 连签奖励
 */
export const StreakRewardSchema = z.object({
  days: z.number(),
  reward: z.object({
    type: z.enum(['points', 'token', 'nft', 'badge', 'multiplier']),
    amount: z.number().optional(),
    name: z.string(),
    icon: z.string().optional(),
  }),
});

export type StreakReward = z.infer<typeof StreakRewardSchema>;

/**
 * DayChain配置Schema
 */
export const DayChainConfigSchema = z.object({
  // 基础配置
  name: z.string(),
  description: z.string().optional(),
  totalDays: z.number().min(1).max(365).default(7),
  
  // 时间配置
  startTime: z.string().datetime(),
  resetTime: z.string().default('00:00'), // 每日重置时间 (UTC)
  timezone: z.string().default('UTC'),
  
  // 任务配置
  days: z.array(DayConfigSchema),
  
  // 连签奖励
  streakRewards: z.array(StreakRewardSchema).optional(),
  
  // 规则配置
  allowCatchUp: z.boolean().default(false), // 是否允许补签
  catchUpCost: z.number().optional(), // 补签消耗积分
  maxMissedDays: z.number().default(0), // 最多可错过天数
  
  // 显示配置
  layout: z.enum(['calendar', 'timeline', 'cards']).default('timeline'),
  showStreak: z.boolean().default(true),
  showCountdown: z.boolean().default(true),
  showProgress: z.boolean().default(true),
  
  // 样式
  accentColor: z.string().optional(),
  completedColor: z.string().default('#10B981'),
  missedColor: z.string().default('#EF4444'),
});

export type DayChainConfig = z.infer<typeof DayChainConfigSchema>;

export interface DayChainProps extends DayChainConfig {
  widgetId: string;
  widgetContext?: Record<string, unknown>;
  onDayComplete?: (day: number) => void;
  onStreakReward?: (days: number, reward: StreakReward) => void;
}

/**
 * DayChain状态
 */
export interface DayChainState {
  currentDay: number;
  completedDays: number[];
  missedDays: number[];
  currentStreak: number;
  longestStreak: number;
  nextUnlockTime: string;
}
```

```typescript
// src/widgets/task/DayChain/DayChain.tsx

import React, { useMemo } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { differenceInDays, differenceInSeconds, format, parseISO } from 'date-fns';
import { registerWidget } from '../../core/decorators';
import { useMCPClient } from '../../hooks/useMCP';
import { useCountdown } from '../../hooks/useCountdown';
import { DayChainConfigSchema, type DayChainProps, type DayChainState } from './types';

const WIDGET_METADATA = {
  type: 'day_chain',
  name: 'Day Chain',
  description: '多日连续挑战任务，培养用户习惯',
  category: 'task' as const,
  icon: 'Calendar',
  version: '1.0.0',
  whitelabel: true,
  mcpDependencies: ['quest-server', 'points-server'],
  configSchema: DayChainConfigSchema,
  defaultConfig: {
    totalDays: 7,
    layout: 'timeline',
    showStreak: true,
    showCountdown: true,
    allowCatchUp: false,
  },
};

export const DayChain: React.FC<DayChainProps> = registerWidget(WIDGET_METADATA)((props) => {
  const {
    widgetId,
    name,
    description,
    totalDays,
    days,
    streakRewards,
    allowCatchUp,
    catchUpCost,
    layout,
    showStreak,
    showCountdown,
    showProgress,
    accentColor,
    completedColor,
    missedColor,
    onDayComplete,
    onStreakReward,
  } = props;

  const { client: mcpClient } = useMCPClient('quest-server');

  // 获取DayChain状态
  const { data: state, refetch } = useQuery({
    queryKey: ['daychain', widgetId],
    queryFn: async () => {
      const result = await mcpClient.callTool('get_daychain_state', { widgetId });
      return result as DayChainState;
    },
  });

  // 倒计时
  const countdown = useCountdown(state?.nextUnlockTime);

  // 完成当日任务
  const completeDayMutation = useMutation({
    mutationFn: async (day: number) => {
      const result = await mcpClient.callTool('complete_day', { widgetId, day });
      return result;
    },
    onSuccess: (result, day) => {
      onDayComplete?.(day);
      
      // 检查连签奖励
      const streakReward = streakRewards?.find(r => r.days === result.currentStreak);
      if (streakReward) {
        onStreakReward?.(result.currentStreak, streakReward);
      }
      
      refetch();
    },
  });

  // 补签
  const catchUpMutation = useMutation({
    mutationFn: async (day: number) => {
      const result = await mcpClient.callTool('catch_up_day', { 
        widgetId, 
        day,
        cost: catchUpCost,
      });
      return result;
    },
    onSuccess: () => refetch(),
  });

  // 计算进度
  const progress = useMemo(() => {
    if (!state) return 0;
    return (state.completedDays.length / totalDays) * 100;
  }, [state, totalDays]);

  // 连签进度
  const nextStreakReward = useMemo(() => {
    if (!streakRewards || !state) return null;
    return streakRewards
      .filter(r => r.days > state.currentStreak)
      .sort((a, b) => a.days - b.days)[0];
  }, [streakRewards, state]);

  // Timeline布局
  const TimelineLayout = () => (
    <div className="flex overflow-x-auto gap-4 pb-4">
      {days.map((dayConfig, index) => {
        const dayNumber = index + 1;
        const isCompleted = state?.completedDays.includes(dayNumber);
        const isMissed = state?.missedDays.includes(dayNumber);
        const isCurrent = state?.currentDay === dayNumber;
        const isLocked = dayConfig.status === 'locked';
        
        return (
          <div
            key={dayNumber}
            className={`
              flex-shrink-0 w-24 p-4 rounded-xl text-center transition-all
              ${isCompleted ? 'bg-green-100 border-green-500' : ''}
              ${isMissed ? 'bg-red-50 border-red-300' : ''}
              ${isCurrent ? 'bg-primary/10 border-primary ring-2 ring-primary/30' : ''}
              ${isLocked ? 'bg-gray-50 border-gray-200 opacity-50' : ''}
              ${!isCompleted && !isMissed && !isCurrent && !isLocked ? 'bg-white border-gray-200' : ''}
              border-2
            `}
          >
            {/* Day标签 */}
            <div className="text-xs text-gray-500 mb-1">Day</div>
            <div className="text-2xl font-bold mb-2">{dayNumber}</div>
            
            {/* 状态图标 */}
            {isCompleted && <span className="text-green-500 text-xl">✓</span>}
            {isMissed && <span className="text-red-500 text-xl">✗</span>}
            {isLocked && <span className="text-gray-400 text-xl">🔒</span>}
            
            {/* 当日任务数 */}
            {!isLocked && (
              <div className="text-xs text-gray-500 mt-2">
                {dayConfig.tasks.length} tasks
              </div>
            )}
            
            {/* 当日奖励 */}
            {dayConfig.reward && (
              <div className="text-xs text-primary mt-1">
                +{dayConfig.reward.amount} {dayConfig.reward.name}
              </div>
            )}
            
            {/* 补签按钮 */}
            {allowCatchUp && isMissed && (
              <button
                onClick={() => catchUpMutation.mutate(dayNumber)}
                className="mt-2 text-xs text-blue-500 underline"
                disabled={catchUpMutation.isPending}
              >
                Catch up ({catchUpCost} pts)
              </button>
            )}
          </div>
        );
      })}
    </div>
  );

  // Calendar布局
  const CalendarLayout = () => {
    // 7列日历网格
    const weeks = Math.ceil(totalDays / 7);
    
    return (
      <div className="grid grid-cols-7 gap-2">
        {/* 星期标题 */}
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="text-center text-xs text-gray-500 py-1">
            {day}
          </div>
        ))}
        
        {/* 日期格子 */}
        {Array.from({ length: weeks * 7 }).map((_, index) => {
          const dayNumber = index + 1;
          if (dayNumber > totalDays) return <div key={index} />;
          
          const isCompleted = state?.completedDays.includes(dayNumber);
          const isMissed = state?.missedDays.includes(dayNumber);
          const isCurrent = state?.currentDay === dayNumber;
          
          return (
            <div
              key={index}
              className={`
                aspect-square rounded-lg flex items-center justify-center text-sm font-medium
                ${isCompleted ? 'bg-green-500 text-white' : ''}
                ${isMissed ? 'bg-red-100 text-red-500' : ''}
                ${isCurrent ? 'bg-primary text-white ring-2 ring-primary/30' : ''}
                ${!isCompleted && !isMissed && !isCurrent ? 'bg-gray-100 text-gray-400' : ''}
              `}
            >
              {dayNumber}
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className="taskon-daychain" data-widget-id={widgetId}>
      {/* 标题区 */}
      <div className="mb-6">
        <h3 className="text-xl font-bold">{name}</h3>
        {description && <p className="text-gray-500 mt-1">{description}</p>}
      </div>
      
      {/* 连签统计 */}
      {showStreak && state && (
        <div className="flex gap-6 mb-6 p-4 bg-gradient-to-r from-primary/10 to-transparent rounded-xl">
          <div>
            <div className="text-3xl font-bold text-primary">{state.currentStreak}</div>
            <div className="text-sm text-gray-500">Current Streak</div>
          </div>
          <div>
            <div className="text-3xl font-bold">{state.longestStreak}</div>
            <div className="text-sm text-gray-500">Longest Streak</div>
          </div>
          {nextStreakReward && (
            <div className="ml-auto text-right">
              <div className="text-sm text-gray-500">Next reward at</div>
              <div className="font-medium">
                {nextStreakReward.days} days ({nextStreakReward.days - state.currentStreak} to go)
              </div>
            </div>
          )}
        </div>
      )}
      
      {/* 进度条 */}
      {showProgress && (
        <div className="mb-6">
          <div className="flex justify-between text-sm mb-2">
            <span>Progress</span>
            <span>{state?.completedDays.length || 0} / {totalDays} days</span>
          </div>
          <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div 
              className="h-full bg-primary transition-all duration-500"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>
      )}
      
      {/* 倒计时 */}
      {showCountdown && countdown && (
        <div className="mb-6 p-4 bg-yellow-50 rounded-lg text-center">
          <div className="text-sm text-gray-500 mb-1">Next day unlocks in</div>
          <div className="text-2xl font-mono font-bold">
            {countdown.hours.toString().padStart(2, '0')}:
            {countdown.minutes.toString().padStart(2, '0')}:
            {countdown.seconds.toString().padStart(2, '0')}
          </div>
        </div>
      )}
      
      {/* 日期布局 */}
      {layout === 'timeline' && <TimelineLayout />}
      {layout === 'calendar' && <CalendarLayout />}
      
      {/* 连签奖励预览 */}
      {streakRewards && streakRewards.length > 0 && (
        <div className="mt-6 p-4 bg-gray-50 rounded-lg">
          <h4 className="font-medium mb-3">Streak Rewards</h4>
          <div className="flex flex-wrap gap-3">
            {streakRewards.map((reward, index) => (
              <div 
                key={index}
                className={`
                  px-3 py-2 rounded-lg border
                  ${state && state.currentStreak >= reward.days 
                    ? 'bg-green-100 border-green-300' 
                    : 'bg-white border-gray-200'}
                `}
              >
                <div className="text-sm font-medium">{reward.days} Days</div>
                <div className="text-xs text-gray-500">{reward.reward.name}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
});
```

### 3.4 其他核心组件规范表

| 组件 | 类型 | 核心Props | MCP依赖 | 关键功能 |
|------|------|-----------|---------|----------|
| **TaskChain** | task | tasks[], unlockMode, completionReward | quest-server | 顺序解锁任务链 |
| **Milestone** | incentive | metric, thresholds[], displayStyle | points-server | 阶段性目标追踪 |
| **BenefitsShop** | incentive | items[], categories[], pointsName | rewards-server | 积分兑换商城 |
| **LuckyWheel** | incentive | prizes[], spinCost, dailyLimit | rewards-server | 抽奖转盘 |
| **UserCenter** | user | showAssets, showHistory, showLevel | user-server | 个人中心 |
| **ProfileCard** | user | showBadges, showStats, editable | user-server | 用户资料卡 |
| **AssetPanel** | user | showBalance, showNFTs, showHistory | rewards-server | 资产面板 |

---

## 4. 组件通信机制

### 4.1 事件总线 (EventBus)

```typescript
// src/widgets/core/EventBus.ts

import { createContext, useContext, useEffect, useCallback, useRef } from 'react';
import mitt, { Emitter, EventType, Handler } from 'mitt';

/**
 * 全局事件类型定义
 */
export interface WidgetEvents {
  // 任务事件
  'task:click': { taskId: string; widgetId: string };
  'task:verify:start': { taskId: string };
  'task:verify:success': { taskId: string; points: number };
  'task:verify:failed': { taskId: string; error: string };
  'task:complete': { taskId: string; widgetId: string };
  
  // 积分事件
  'points:earned': { amount: number; source: string; widgetId: string };
  'points:spent': { amount: number; target: string; widgetId: string };
  'points:balance:updated': { balance: number };
  
  // 用户事件
  'user:level:up': { oldLevel: number; newLevel: number };
  'user:badge:earned': { badgeId: string; badgeName: string };
  'user:login': { userId: string };
  'user:logout': void;
  
  // 奖励事件
  'reward:claim:start': { rewardId: string };
  'reward:claim:success': { rewardId: string; type: string };
  'reward:claim:failed': { rewardId: string; error: string };
  
  // 排行榜事件
  'leaderboard:rank:changed': { oldRank: number; newRank: number; userId: string };
  
  // 组件间通信
  'widget:config:changed': { widgetId: string; config: Record<string, unknown> };
  'widget:refresh': { widgetId: string };
  'widget:error': { widgetId: string; error: Error };
  
  // 全局状态
  'app:theme:changed': { theme: string };
  'app:locale:changed': { locale: string };
}

/**
 * 类型安全的事件总线
 */
class TypedEventBus {
  private emitter: Emitter<WidgetEvents>;
  private debugMode: boolean;

  constructor(debug = false) {
    this.emitter = mitt<WidgetEvents>();
    this.debugMode = debug;
  }

  /**
   * 发射事件
   */
  emit<K extends keyof WidgetEvents>(event: K, payload: WidgetEvents[K]): void {
    if (this.debugMode) {
      console.log(`[EventBus] Emit: ${event}`, payload);
    }
    this.emitter.emit(event, payload);
  }

  /**
   * 监听事件
   */
  on<K extends keyof WidgetEvents>(event: K, handler: Handler<WidgetEvents[K]>): void {
    this.emitter.on(event, handler);
  }

  /**
   * 移除监听
   */
  off<K extends keyof WidgetEvents>(event: K, handler: Handler<WidgetEvents[K]>): void {
    this.emitter.off(event, handler);
  }

  /**
   * 监听所有事件
   */
  onAny(handler: (event: keyof WidgetEvents, payload: unknown) => void): void {
    this.emitter.on('*', handler as any);
  }

  /**
   * 清除所有监听
   */
  clear(): void {
    this.emitter.all.clear();
  }
}

// 创建单例
const eventBus = new TypedEventBus(process.env.NODE_ENV === 'development');

// Context
const EventBusContext = createContext(eventBus);

/**
 * EventBus Provider
 */
export const EventBusProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <EventBusContext.Provider value={eventBus}>
      {children}
    </EventBusContext.Provider>
  );
};

/**
 * useEventBus Hook
 */
export function useEventBus() {
  const bus = useContext(EventBusContext);
  
  const emit = useCallback(<K extends keyof WidgetEvents>(
    event: K, 
    payload: WidgetEvents[K]
  ) => {
    bus.emit(event, payload);
  }, [bus]);

  return { emit, bus };
}

/**
 * useEventListener Hook - 自动清理的事件监听
 */
export function useEventListener<K extends keyof WidgetEvents>(
  event: K,
  handler: Handler<WidgetEvents[K]>,
  deps: React.DependencyList = []
) {
  const bus = useContext(EventBusContext);
  const savedHandler = useRef(handler);
  
  useEffect(() => {
    savedHandler.current = handler;
  }, [handler]);

  useEffect(() => {
    const eventHandler: Handler<WidgetEvents[K]> = (payload) => {
      savedHandler.current(payload);
    };
    
    bus.on(event, eventHandler);
    return () => bus.off(event, eventHandler);
  }, [bus, event, ...deps]);
}

export { eventBus };
```

### 4.2 组件间通信示例

```typescript
// 示例：任务完成后更新积分和排行榜

// TaskList组件中
const handleVerifySuccess = (taskId: string, points: number) => {
  emit('task:verify:success', { taskId, points });
  emit('points:earned', { amount: points, source: 'task', widgetId });
};

// PointsDisplay组件中
useEventListener('points:earned', ({ amount }) => {
  // 触发积分动画
  setAnimatingPoints(amount);
  // 更新余额
  refetchBalance();
});

// Leaderboard组件中
useEventListener('points:earned', () => {
  // 可能导致排名变化，刷新排行榜
  refetchLeaderboard();
});

// 全局通知组件中
useEventListener('user:badge:earned', ({ badgeName }) => {
  showToast({
    type: 'success',
    title: 'Badge Earned!',
    message: `You've earned the "${badgeName}" badge!`,
  });
});
```

---

## 5. 主题系统

### 5.1 主题定义

```typescript
// src/widgets/themes/types.ts

/**
 * 颜色方案
 */
export interface ColorScheme {
  // 品牌色
  primary: string;
  primaryHover: string;
  primaryActive: string;
  primaryLight: string;
  
  // 语义色
  success: string;
  warning: string;
  error: string;
  info: string;
  
  // 中性色
  background: string;
  surface: string;
  surfaceHover: string;
  border: string;
  divider: string;
  
  // 文字色
  textPrimary: string;
  textSecondary: string;
  textTertiary: string;
  textDisabled: string;
  textInverse: string;
}

/**
 * 排版系统
 */
export interface Typography {
  fontFamily: {
    sans: string;
    mono: string;
  };
  fontSize: {
    xs: string;
    sm: string;
    base: string;
    lg: string;
    xl: string;
    '2xl': string;
    '3xl': string;
  };
  fontWeight: {
    normal: number;
    medium: number;
    semibold: number;
    bold: number;
  };
  lineHeight: {
    tight: number;
    normal: number;
    relaxed: number;
  };
}

/**
 * 间距系统
 */
export interface Spacing {
  0: string;
  1: string;
  2: string;
  3: string;
  4: string;
  5: string;
  6: string;
  8: string;
  10: string;
  12: string;
  16: string;
  20: string;
  24: string;
}

/**
 * 圆角系统
 */
export interface BorderRadius {
  none: string;
  sm: string;
  md: string;
  lg: string;
  xl: string;
  '2xl': string;
  full: string;
}

/**
 * 阴影系统
 */
export interface Shadows {
  none: string;
  sm: string;
  md: string;
  lg: string;
  xl: string;
}

/**
 * 动画系统
 */
export interface Animation {
  duration: {
    fast: string;
    normal: string;
    slow: string;
  };
  easing: {
    linear: string;
    easeIn: string;
    easeOut: string;
    easeInOut: string;
  };
}

/**
 * 组件特定主题
 */
export interface ComponentThemes {
  button: {
    borderRadius: string;
    paddingX: string;
    paddingY: string;
  };
  card: {
    borderRadius: string;
    padding: string;
    shadow: string;
  };
  badge: {
    borderRadius: string;
    fontSize: string;
  };
  input: {
    borderRadius: string;
    borderWidth: string;
  };
  // ... 其他组件
}

/**
 * 完整主题定义
 */
export interface Theme {
  name: string;
  colors: ColorScheme;
  typography: Typography;
  spacing: Spacing;
  borderRadius: BorderRadius;
  shadows: Shadows;
  animation: Animation;
  components: ComponentThemes;
}
```

### 5.2 默认主题

```typescript
// src/widgets/themes/default.ts

import { Theme } from './types';

export const defaultTheme: Theme = {
  name: 'default',
  
  colors: {
    primary: '#6366F1',
    primaryHover: '#4F46E5',
    primaryActive: '#4338CA',
    primaryLight: '#EEF2FF',
    
    success: '#10B981',
    warning: '#F59E0B',
    error: '#EF4444',
    info: '#3B82F6',
    
    background: '#F9FAFB',
    surface: '#FFFFFF',
    surfaceHover: '#F3F4F6',
    border: '#E5E7EB',
    divider: '#F3F4F6',
    
    textPrimary: '#111827',
    textSecondary: '#6B7280',
    textTertiary: '#9CA3AF',
    textDisabled: '#D1D5DB',
    textInverse: '#FFFFFF',
  },
  
  typography: {
    fontFamily: {
      sans: 'Inter, system-ui, -apple-system, sans-serif',
      mono: 'JetBrains Mono, Consolas, monospace',
    },
    fontSize: {
      xs: '0.75rem',
      sm: '0.875rem',
      base: '1rem',
      lg: '1.125rem',
      xl: '1.25rem',
      '2xl': '1.5rem',
      '3xl': '1.875rem',
    },
    fontWeight: {
      normal: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
    },
    lineHeight: {
      tight: 1.25,
      normal: 1.5,
      relaxed: 1.75,
    },
  },
  
  spacing: {
    0: '0',
    1: '0.25rem',
    2: '0.5rem',
    3: '0.75rem',
    4: '1rem',
    5: '1.25rem',
    6: '1.5rem',
    8: '2rem',
    10: '2.5rem',
    12: '3rem',
    16: '4rem',
    20: '5rem',
    24: '6rem',
  },
  
  borderRadius: {
    none: '0',
    sm: '0.25rem',
    md: '0.375rem',
    lg: '0.5rem',
    xl: '0.75rem',
    '2xl': '1rem',
    full: '9999px',
  },
  
  shadows: {
    none: 'none',
    sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
    md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
    lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
    xl: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
  },
  
  animation: {
    duration: {
      fast: '150ms',
      normal: '300ms',
      slow: '500ms',
    },
    easing: {
      linear: 'linear',
      easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
      easeOut: 'cubic-bezier(0, 0, 0.2, 1)',
      easeInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    },
  },
  
  components: {
    button: {
      borderRadius: '0.5rem',
      paddingX: '1rem',
      paddingY: '0.5rem',
    },
    card: {
      borderRadius: '0.75rem',
      padding: '1.5rem',
      shadow: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
    },
    badge: {
      borderRadius: '9999px',
      fontSize: '0.75rem',
    },
    input: {
      borderRadius: '0.5rem',
      borderWidth: '1px',
    },
  },
};
```

### 5.3 主题Provider和Hook

```typescript
// src/widgets/hooks/useTheme.ts

import React, { createContext, useContext, useState, useCallback, useMemo } from 'react';
import { Theme } from '../themes/types';
import { defaultTheme } from '../themes/default';
import { darkTheme } from '../themes/dark';

interface ThemeContextValue {
  theme: Theme;
  themeName: string;
  setTheme: (name: string) => void;
  customizeTheme: (overrides: DeepPartial<Theme>) => void;
  resolveColor: (color: string) => string;
  cssVariables: Record<string, string>;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

const themes: Record<string, Theme> = {
  default: defaultTheme,
  dark: darkTheme,
};

/**
 * 深度合并对象
 */
function deepMerge<T extends object>(target: T, source: DeepPartial<T>): T {
  const result = { ...target };
  
  for (const key in source) {
    const sourceValue = source[key];
    const targetValue = result[key];
    
    if (sourceValue !== undefined) {
      if (
        typeof sourceValue === 'object' &&
        sourceValue !== null &&
        typeof targetValue === 'object' &&
        targetValue !== null
      ) {
        (result as any)[key] = deepMerge(targetValue as object, sourceValue as object);
      } else {
        (result as any)[key] = sourceValue;
      }
    }
  }
  
  return result;
}

/**
 * 生成CSS变量
 */
function generateCSSVariables(theme: Theme): Record<string, string> {
  const vars: Record<string, string> = {};
  
  // Colors
  Object.entries(theme.colors).forEach(([key, value]) => {
    vars[`--color-${kebabCase(key)}`] = value;
  });
  
  // Typography
  vars['--font-family-sans'] = theme.typography.fontFamily.sans;
  vars['--font-family-mono'] = theme.typography.fontFamily.mono;
  
  Object.entries(theme.typography.fontSize).forEach(([key, value]) => {
    vars[`--font-size-${key}`] = value;
  });
  
  // Spacing
  Object.entries(theme.spacing).forEach(([key, value]) => {
    vars[`--spacing-${key}`] = value;
  });
  
  // Border Radius
  Object.entries(theme.borderRadius).forEach(([key, value]) => {
    vars[`--radius-${key}`] = value;
  });
  
  // Shadows
  Object.entries(theme.shadows).forEach(([key, value]) => {
    vars[`--shadow-${key}`] = value;
  });
  
  return vars;
}

function kebabCase(str: string): string {
  return str.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase();
}

/**
 * Theme Provider
 */
export const ThemeProvider: React.FC<{
  children: React.ReactNode;
  defaultTheme?: string;
  customTheme?: DeepPartial<Theme>;
}> = ({ children, defaultTheme: initialTheme = 'default', customTheme }) => {
  const [themeName, setThemeName] = useState(initialTheme);
  const [customOverrides, setCustomOverrides] = useState<DeepPartial<Theme>>(customTheme || {});
  
  // 计算最终主题
  const theme = useMemo(() => {
    const baseTheme = themes[themeName] || themes.default;
    return deepMerge(baseTheme, customOverrides);
  }, [themeName, customOverrides]);
  
  // 生成CSS变量
  const cssVariables = useMemo(() => generateCSSVariables(theme), [theme]);
  
  // 切换主题
  const setTheme = useCallback((name: string) => {
    if (themes[name]) {
      setThemeName(name);
    }
  }, []);
  
  // 自定义主题
  const customizeTheme = useCallback((overrides: DeepPartial<Theme>) => {
    setCustomOverrides(prev => deepMerge(prev as Theme, overrides) as DeepPartial<Theme>);
  }, []);
  
  // 解析颜色（支持主题变量引用）
  const resolveColor = useCallback((color: string): string => {
    if (color.startsWith('$')) {
      const colorKey = color.slice(1);
      return (theme.colors as any)[colorKey] || color;
    }
    return color;
  }, [theme]);
  
  // 应用CSS变量到root
  React.useEffect(() => {
    const root = document.documentElement;
    Object.entries(cssVariables).forEach(([key, value]) => {
      root.style.setProperty(key, value);
    });
  }, [cssVariables]);
  
  const value = useMemo(() => ({
    theme,
    themeName,
    setTheme,
    customizeTheme,
    resolveColor,
    cssVariables,
  }), [theme, themeName, setTheme, customizeTheme, resolveColor, cssVariables]);
  
  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

/**
 * useTheme Hook
 */
export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
```

---

## 6. MCP集成规范

### 6.1 MCP Client Hook

```typescript
// src/widgets/hooks/useMCP.ts

import { createContext, useContext, useCallback, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

/**
 * MCP工具调用参数
 */
interface MCPToolCallParams {
  name: string;
  arguments: Record<string, unknown>;
}

/**
 * MCP资源读取参数
 */
interface MCPResourceParams {
  uri: string;
}

/**
 * MCP Client接口
 */
interface MCPClient {
  /** 调用MCP工具 */
  callTool: <T = unknown>(name: string, args: Record<string, unknown>) => Promise<T>;
  /** 读取MCP资源 */
  readResource: <T = unknown>(uri: string) => Promise<T>;
  /** 订阅资源变更 */
  subscribe: (uri: string, callback: (data: unknown) => void) => () => void;
  /** 连接状态 */
  isConnected: boolean;
  /** 重连 */
  reconnect: () => Promise<void>;
}

/**
 * MCP Context
 */
interface MCPContextValue {
  clients: Map<string, MCPClient>;
  getClient: (serverName: string) => MCPClient | undefined;
  registerClient: (serverName: string, client: MCPClient) => void;
}

const MCPContext = createContext<MCPContextValue | null>(null);

/**
 * MCP Provider
 */
export const MCPProvider: React.FC<{
  children: React.ReactNode;
  config: {
    servers: Array<{
      name: string;
      url: string;
      authToken?: string;
    }>;
  };
}> = ({ children, config }) => {
  const clientsRef = useRef(new Map<string, MCPClient>());
  
  // 初始化所有MCP客户端
  useEffect(() => {
    config.servers.forEach(async (serverConfig) => {
      const client = await createMCPClient(serverConfig);
      clientsRef.current.set(serverConfig.name, client);
    });
    
    return () => {
      // 清理连接
      clientsRef.current.forEach(client => client.disconnect?.());
      clientsRef.current.clear();
    };
  }, [config]);
  
  const getClient = useCallback((serverName: string) => {
    return clientsRef.current.get(serverName);
  }, []);
  
  const registerClient = useCallback((serverName: string, client: MCPClient) => {
    clientsRef.current.set(serverName, client);
  }, []);
  
  const value = useMemo(() => ({
    clients: clientsRef.current,
    getClient,
    registerClient,
  }), [getClient, registerClient]);
  
  return (
    <MCPContext.Provider value={value}>
      {children}
    </MCPContext.Provider>
  );
};

/**
 * useMCPClient Hook
 */
export function useMCPClient(serverName: string) {
  const context = useContext(MCPContext);
  if (!context) {
    throw new Error('useMCPClient must be used within an MCPProvider');
  }
  
  const client = context.getClient(serverName);
  
  return {
    client,
    isConnected: client?.isConnected ?? false,
  };
}

/**
 * useMCPQuery Hook - 封装MCP资源查询
 */
export function useMCPQuery<T = unknown>(
  serverName: string,
  resourceUri: string,
  options?: {
    enabled?: boolean;
    refetchInterval?: number;
    staleTime?: number;
  }
) {
  const { client } = useMCPClient(serverName);
  
  return useQuery({
    queryKey: ['mcp', serverName, resourceUri],
    queryFn: () => client!.readResource<T>(resourceUri),
    enabled: !!client && (options?.enabled ?? true),
    refetchInterval: options?.refetchInterval,
    staleTime: options?.staleTime,
  });
}

/**
 * useMCPMutation Hook - 封装MCP工具调用
 */
export function useMCPMutation<TData = unknown, TVariables = Record<string, unknown>>(
  serverName: string,
  toolName: string,
  options?: {
    onSuccess?: (data: TData, variables: TVariables) => void;
    onError?: (error: Error, variables: TVariables) => void;
    invalidateQueries?: string[];
  }
) {
  const { client } = useMCPClient(serverName);
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (variables: TVariables) => 
      client!.callTool<TData>(toolName, variables as Record<string, unknown>),
    onSuccess: (data, variables) => {
      options?.onSuccess?.(data, variables);
      
      // 自动失效相关查询
      if (options?.invalidateQueries) {
        options.invalidateQueries.forEach(key => {
          queryClient.invalidateQueries({ queryKey: [key] });
        });
      }
    },
    onError: (error: Error, variables) => {
      options?.onError?.(error, variables);
    },
  });
}

/**
 * useMCPSubscription Hook - 订阅MCP资源变更
 */
export function useMCPSubscription<T = unknown>(
  serverName: string,
  resourceUri: string,
  callback: (data: T) => void
) {
  const { client } = useMCPClient(serverName);
  
  useEffect(() => {
    if (!client) return;
    
    const unsubscribe = client.subscribe(resourceUri, callback as (data: unknown) => void);
    return () => unsubscribe();
  }, [client, resourceUri, callback]);
}
```

### 6.2 MCP Client实现

```typescript
// src/widgets/core/MCPClient.ts

import { EventEmitter } from 'eventemitter3';

interface MCPClientConfig {
  name: string;
  url: string;
  authToken?: string;
  reconnectAttempts?: number;
  reconnectDelay?: number;
}

interface MCPMessage {
  jsonrpc: '2.0';
  id?: string | number;
  method?: string;
  params?: unknown;
  result?: unknown;
  error?: {
    code: number;
    message: string;
    data?: unknown;
  };
}

/**
 * MCP WebSocket Client
 */
export class MCPWebSocketClient extends EventEmitter {
  private ws: WebSocket | null = null;
  private config: MCPClientConfig;
  private requestId = 0;
  private pendingRequests = new Map<string | number, {
    resolve: (value: unknown) => void;
    reject: (error: Error) => void;
    timeout: NodeJS.Timeout;
  }>();
  private subscriptions = new Map<string, Set<(data: unknown) => void>>();
  private reconnectAttempts = 0;
  private _isConnected = false;

  constructor(config: MCPClientConfig) {
    super();
    this.config = {
      reconnectAttempts: 5,
      reconnectDelay: 1000,
      ...config,
    };
  }

  get isConnected(): boolean {
    return this._isConnected;
  }

  /**
   * 连接到MCP服务器
   */
  async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const url = new URL(this.config.url);
        if (this.config.authToken) {
          url.searchParams.set('token', this.config.authToken);
        }
        
        this.ws = new WebSocket(url.toString());
        
        this.ws.onopen = () => {
          this._isConnected = true;
          this.reconnectAttempts = 0;
          this.emit('connected');
          resolve();
        };
        
        this.ws.onclose = (event) => {
          this._isConnected = false;
          this.emit('disconnected', event);
          this.handleReconnect();
        };
        
        this.ws.onerror = (error) => {
          this.emit('error', error);
          reject(error);
        };
        
        this.ws.onmessage = (event) => {
          this.handleMessage(JSON.parse(event.data));
        };
        
      } catch (error) {
        reject(error);
      }
    });
  }

  /**
   * 处理重连
   */
  private handleReconnect(): void {
    if (this.reconnectAttempts >= (this.config.reconnectAttempts || 5)) {
      this.emit('reconnect_failed');
      return;
    }
    
    this.reconnectAttempts++;
    const delay = (this.config.reconnectDelay || 1000) * Math.pow(2, this.reconnectAttempts - 1);
    
    setTimeout(() => {
      this.emit('reconnecting', { attempt: this.reconnectAttempts });
      this.connect().catch(() => {});
    }, delay);
  }

  /**
   * 处理收到的消息
   */
  private handleMessage(message: MCPMessage): void {
    // 响应消息
    if (message.id !== undefined) {
      const pending = this.pendingRequests.get(message.id);
      if (pending) {
        clearTimeout(pending.timeout);
        this.pendingRequests.delete(message.id);
        
        if (message.error) {
          pending.reject(new Error(message.error.message));
        } else {
          pending.resolve(message.result);
        }
      }
      return;
    }
    
    // 通知消息（资源变更）
    if (message.method === 'notifications/resources/updated') {
      const { uri, contents } = message.params as { uri: string; contents: unknown };
      const callbacks = this.subscriptions.get(uri);
      if (callbacks) {
        callbacks.forEach(cb => cb(contents));
      }
    }
  }

  /**
   * 发送请求
   */
  private async sendRequest<T>(method: string, params?: unknown): Promise<T> {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected');
    }
    
    const id = ++this.requestId;
    
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        this.pendingRequests.delete(id);
        reject(new Error('Request timeout'));
      }, 30000);
      
      this.pendingRequests.set(id, { resolve: resolve as any, reject, timeout });
      
      this.ws!.send(JSON.stringify({
        jsonrpc: '2.0',
        id,
        method,
        params,
      }));
    });
  }

  /**
   * 调用MCP工具
   */
  async callTool<T = unknown>(name: string, args: Record<string, unknown>): Promise<T> {
    return this.sendRequest<T>('tools/call', { name, arguments: args });
  }

  /**
   * 读取MCP资源
   */
  async readResource<T = unknown>(uri: string): Promise<T> {
    return this.sendRequest<T>('resources/read', { uri });
  }

  /**
   * 订阅资源变更
   */
  subscribe(uri: string, callback: (data: unknown) => void): () => void {
    if (!this.subscriptions.has(uri)) {
      this.subscriptions.set(uri, new Set());
      // 发送订阅请求
      this.sendRequest('resources/subscribe', { uri }).catch(console.error);
    }
    
    this.subscriptions.get(uri)!.add(callback);
    
    // 返回取消订阅函数
    return () => {
      const callbacks = this.subscriptions.get(uri);
      if (callbacks) {
        callbacks.delete(callback);
        if (callbacks.size === 0) {
          this.subscriptions.delete(uri);
          // 发送取消订阅请求
          this.sendRequest('resources/unsubscribe', { uri }).catch(console.error);
        }
      }
    };
  }

  /**
   * 断开连接
   */
  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.pendingRequests.clear();
    this.subscriptions.clear();
    this._isConnected = false;
  }

  /**
   * 重连
   */
  async reconnect(): Promise<void> {
    this.disconnect();
    this.reconnectAttempts = 0;
    await this.connect();
  }
}

/**
 * 创建MCP客户端
 */
export async function createMCPClient(config: MCPClientConfig): Promise<MCPWebSocketClient> {
  const client = new MCPWebSocketClient(config);
  await client.connect();
  return client;
}
```

---

## 7. 状态管理

### 7.1 Widget Store

```typescript
// src/widgets/stores/widgetStore.ts

import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { devtools } from 'zustand/middleware';
import type { WidgetConfig } from '../core/WidgetRenderer';

interface WidgetState {
  // Widget配置
  widgets: WidgetConfig[];
  
  // Widget状态缓存
  widgetStates: Record<string, Record<string, unknown>>;
  
  // 全局上下文
  globalContext: Record<string, unknown>;
  
  // Actions
  setWidgets: (widgets: WidgetConfig[]) => void;
  addWidget: (widget: WidgetConfig) => void;
  removeWidget: (widgetId: string) => void;
  updateWidgetConfig: (widgetId: string, config: Partial<WidgetConfig['config']>) => void;
  reorderWidgets: (fromIndex: number, toIndex: number) => void;
  
  // Widget状态管理
  setWidgetState: (widgetId: string, state: Record<string, unknown>) => void;
  updateWidgetState: (widgetId: string, updates: Record<string, unknown>) => void;
  clearWidgetState: (widgetId: string) => void;
  
  // 全局上下文
  setGlobalContext: (context: Record<string, unknown>) => void;
  updateGlobalContext: (updates: Record<string, unknown>) => void;
}

export const useWidgetStore = create<WidgetState>()(
  devtools(
    immer((set) => ({
      widgets: [],
      widgetStates: {},
      globalContext: {},
      
      setWidgets: (widgets) => set((state) => {
        state.widgets = widgets;
      }),
      
      addWidget: (widget) => set((state) => {
        state.widgets.push(widget);
        // 按position排序
        state.widgets.sort((a, b) => a.position - b.position);
      }),
      
      removeWidget: (widgetId) => set((state) => {
        state.widgets = state.widgets.filter(w => w.id !== widgetId);
        delete state.widgetStates[widgetId];
      }),
      
      updateWidgetConfig: (widgetId, config) => set((state) => {
        const widget = state.widgets.find(w => w.id === widgetId);
        if (widget) {
          widget.config = { ...widget.config, ...config };
        }
      }),
      
      reorderWidgets: (fromIndex, toIndex) => set((state) => {
        const [removed] = state.widgets.splice(fromIndex, 1);
        state.widgets.splice(toIndex, 0, removed);
        // 更新position
        state.widgets.forEach((w, i) => {
          w.position = i;
        });
      }),
      
      setWidgetState: (widgetId, widgetState) => set((state) => {
        state.widgetStates[widgetId] = widgetState;
      }),
      
      updateWidgetState: (widgetId, updates) => set((state) => {
        if (!state.widgetStates[widgetId]) {
          state.widgetStates[widgetId] = {};
        }
        Object.assign(state.widgetStates[widgetId], updates);
      }),
      
      clearWidgetState: (widgetId) => set((state) => {
        delete state.widgetStates[widgetId];
      }),
      
      setGlobalContext: (context) => set((state) => {
        state.globalContext = context;
      }),
      
      updateGlobalContext: (updates) => set((state) => {
        Object.assign(state.globalContext, updates);
      }),
    })),
    { name: 'widget-store' }
  )
);
```

### 7.2 User Store

```typescript
// src/widgets/stores/userStore.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: string;
  address: string;
  displayName?: string;
  avatar?: string;
  level: number;
  points: number;
  badges: string[];
  joinedAt: string;
}

interface UserState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // Actions
  setUser: (user: User | null) => void;
  updateUser: (updates: Partial<User>) => void;
  addPoints: (amount: number) => void;
  addBadge: (badgeId: string) => void;
  levelUp: () => void;
  logout: () => void;
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      
      setUser: (user) => set({
        user,
        isAuthenticated: !!user,
      }),
      
      updateUser: (updates) => set((state) => ({
        user: state.user ? { ...state.user, ...updates } : null,
      })),
      
      addPoints: (amount) => set((state) => ({
        user: state.user ? {
          ...state.user,
          points: state.user.points + amount,
        } : null,
      })),
      
      addBadge: (badgeId) => set((state) => ({
        user: state.user ? {
          ...state.user,
          badges: [...new Set([...state.user.badges, badgeId])],
        } : null,
      })),
      
      levelUp: () => set((state) => ({
        user: state.user ? {
          ...state.user,
          level: state.user.level + 1,
        } : null,
      })),
      
      logout: () => set({
        user: null,
        isAuthenticated: false,
      }),
    }),
    {
      name: 'taskon-user',
      partialize: (state) => ({ user: state.user }),
    }
  )
);
```

---

## 8. 性能优化

### 8.1 组件懒加载

```typescript
// src/widgets/core/lazyWidgets.ts

import { lazy, ComponentType } from 'react';

/**
 * Widget懒加载配置
 */
const widgetModules: Record<string, () => Promise<{ default: ComponentType<any> }>> = {
  'task_list': () => import('../task/TaskList/TaskList'),
  'task_chain': () => import('../task/TaskChain/TaskChain'),
  'day_chain': () => import('../task/DayChain/DayChain'),
  'leaderboard': () => import('../incentive/Leaderboard/Leaderboard'),
  'milestone': () => import('../incentive/Milestone/Milestone'),
  'benefits_shop': () => import('../incentive/BenefitsShop/BenefitsShop'),
  'lucky_wheel': () => import('../incentive/LuckyWheel/LuckyWheel'),
  'user_center': () => import('../user/UserCenter/UserCenter'),
};

/**
 * 创建懒加载组件
 */
export function createLazyWidget(type: string) {
  const moduleLoader = widgetModules[type];
  
  if (!moduleLoader) {
    return null;
  }
  
  return lazy(moduleLoader);
}

/**
 * 预加载Widget
 */
export function preloadWidget(type: string): void {
  const moduleLoader = widgetModules[type];
  if (moduleLoader) {
    moduleLoader();
  }
}

/**
 * 预加载多个Widget
 */
export function preloadWidgets(types: string[]): void {
  types.forEach(preloadWidget);
}
```

### 8.2 虚拟化列表

```typescript
// src/widgets/hooks/useVirtualList.ts

import { useCallback, useMemo, useRef, useState } from 'react';

interface VirtualListOptions<T> {
  items: T[];
  itemHeight: number;
  containerHeight: number;
  overscan?: number;
}

interface VirtualListResult<T> {
  virtualItems: Array<{
    index: number;
    item: T;
    style: React.CSSProperties;
  }>;
  totalHeight: number;
  containerRef: React.RefObject<HTMLDivElement>;
  scrollToIndex: (index: number) => void;
}

/**
 * 虚拟列表Hook
 */
export function useVirtualList<T>({
  items,
  itemHeight,
  containerHeight,
  overscan = 5,
}: VirtualListOptions<T>): VirtualListResult<T> {
  const containerRef = useRef<HTMLDivElement>(null);
  const [scrollTop, setScrollTop] = useState(0);
  
  // 计算可见范围
  const { startIndex, endIndex } = useMemo(() => {
    const start = Math.floor(scrollTop / itemHeight);
    const visibleCount = Math.ceil(containerHeight / itemHeight);
    
    return {
      startIndex: Math.max(0, start - overscan),
      endIndex: Math.min(items.length - 1, start + visibleCount + overscan),
    };
  }, [scrollTop, itemHeight, containerHeight, items.length, overscan]);
  
  // 生成虚拟项
  const virtualItems = useMemo(() => {
    const result = [];
    
    for (let i = startIndex; i <= endIndex; i++) {
      result.push({
        index: i,
        item: items[i],
        style: {
          position: 'absolute' as const,
          top: i * itemHeight,
          height: itemHeight,
          width: '100%',
        },
      });
    }
    
    return result;
  }, [items, startIndex, endIndex, itemHeight]);
  
  // 总高度
  const totalHeight = items.length * itemHeight;
  
  // 滚动到指定索引
  const scrollToIndex = useCallback((index: number) => {
    if (containerRef.current) {
      containerRef.current.scrollTop = index * itemHeight;
    }
  }, [itemHeight]);
  
  // 滚动处理
  const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
    setScrollTop(e.currentTarget.scrollTop);
  }, []);
  
  // 绑定滚动事件
  if (containerRef.current) {
    containerRef.current.onscroll = handleScroll as any;
  }
  
  return {
    virtualItems,
    totalHeight,
    containerRef,
    scrollToIndex,
  };
}
```

### 8.3 性能监控

```typescript
// src/widgets/core/PerformanceMonitor.ts

interface WidgetPerformanceMetrics {
  widgetId: string;
  widgetType: string;
  renderTime: number;
  mountTime: number;
  updateCount: number;
  lastUpdateTime: number;
  memoryUsage?: number;
}

class WidgetPerformanceMonitor {
  private metrics = new Map<string, WidgetPerformanceMetrics>();
  private observers = new Map<string, PerformanceObserver>();

  /**
   * 开始监控Widget
   */
  startMonitoring(widgetId: string, widgetType: string): void {
    if (this.metrics.has(widgetId)) return;

    this.metrics.set(widgetId, {
      widgetId,
      widgetType,
      renderTime: 0,
      mountTime: 0,
      updateCount: 0,
      lastUpdateTime: Date.now(),
    });

    // 使用Performance API监控
    if (typeof PerformanceObserver !== 'undefined') {
      const observer = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          if (entry.name.includes(widgetId)) {
            this.recordMetric(widgetId, 'renderTime', entry.duration);
          }
        }
      });

      observer.observe({ entryTypes: ['measure'] });
      this.observers.set(widgetId, observer);
    }
  }

  /**
   * 记录渲染开始
   */
  markRenderStart(widgetId: string): void {
    performance.mark(`widget-render-start-${widgetId}`);
  }

  /**
   * 记录渲染结束
   */
  markRenderEnd(widgetId: string): void {
    performance.mark(`widget-render-end-${widgetId}`);
    performance.measure(
      `widget-render-${widgetId}`,
      `widget-render-start-${widgetId}`,
      `widget-render-end-${widgetId}`
    );

    const metrics = this.metrics.get(widgetId);
    if (metrics) {
      metrics.updateCount++;
      metrics.lastUpdateTime = Date.now();
    }
  }

  /**
   * 记录指标
   */
  private recordMetric(
    widgetId: string,
    metric: keyof WidgetPerformanceMetrics,
    value: number
  ): void {
    const metrics = this.metrics.get(widgetId);
    if (metrics) {
      (metrics as any)[metric] = value;
    }
  }

  /**
   * 获取Widget性能指标
   */
  getMetrics(widgetId: string): WidgetPerformanceMetrics | undefined {
    return this.metrics.get(widgetId);
  }

  /**
   * 获取所有性能指标
   */
  getAllMetrics(): WidgetPerformanceMetrics[] {
    return Array.from(this.metrics.values());
  }

  /**
   * 获取慢Widget（渲染时间超过阈值）
   */
  getSlowWidgets(threshold = 16): WidgetPerformanceMetrics[] {
    return this.getAllMetrics().filter(m => m.renderTime > threshold);
  }

  /**
   * 停止监控
   */
  stopMonitoring(widgetId: string): void {
    const observer = this.observers.get(widgetId);
    if (observer) {
      observer.disconnect();
      this.observers.delete(widgetId);
    }
    this.metrics.delete(widgetId);
  }

  /**
   * 清除所有监控
   */
  clear(): void {
    this.observers.forEach(observer => observer.disconnect());
    this.observers.clear();
    this.metrics.clear();
  }
}

export const performanceMonitor = new WidgetPerformanceMonitor();
```

---

## 9. 测试规范

### 9.1 组件测试模板

```typescript
// src/widgets/task/TaskList/__tests__/TaskList.test.tsx

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { TaskList } from '../TaskList';
import { ThemeProvider } from '../../../hooks/useTheme';
import { EventBusProvider } from '../../../core/EventBus';
import { MCPProvider } from '../../../hooks/useMCP';

// Mock MCP Client
const mockMCPClient = {
  callTool: vi.fn(),
  readResource: vi.fn(),
  subscribe: vi.fn(() => () => {}),
  isConnected: true,
};

// 测试数据
const mockTasks = [
  {
    id: 'task-1',
    type: 'offchain',
    template: 'x_follow',
    title: 'Follow Twitter',
    points: 100,
    status: 'available',
    required: true,
  },
  {
    id: 'task-2',
    type: 'onchain',
    template: 'swap_token',
    title: 'Complete Swap',
    points: 200,
    status: 'locked',
    required: false,
  },
];

// 测试Wrapper
const TestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return (
    <QueryClientProvider client={queryClient}>
      <MCPProvider config={{ servers: [] }}>
        <ThemeProvider>
          <EventBusProvider>
            {children}
          </EventBusProvider>
        </ThemeProvider>
      </MCPProvider>
    </QueryClientProvider>
  );
};

describe('TaskList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockMCPClient.callTool.mockResolvedValue({ tasks: mockTasks });
  });

  it('renders task list correctly', async () => {
    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={mockTasks}
          layout="vertical"
          showProgress={true}
          showPoints={true}
        />
      </TestWrapper>
    );

    await waitFor(() => {
      expect(screen.getByText('Follow Twitter')).toBeInTheDocument();
      expect(screen.getByText('Complete Swap')).toBeInTheDocument();
    });
  });

  it('displays points correctly', async () => {
    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={mockTasks}
          showPoints={true}
        />
      </TestWrapper>
    );

    await waitFor(() => {
      expect(screen.getByText('100')).toBeInTheDocument();
      expect(screen.getByText('200')).toBeInTheDocument();
    });
  });

  it('handles task click', async () => {
    const onTaskClick = vi.fn();

    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={mockTasks}
          enableClick={true}
          onTaskClick={onTaskClick}
        />
      </TestWrapper>
    );

    await waitFor(() => {
      const taskCard = screen.getByText('Follow Twitter').closest('[data-testid="task-card"]');
      if (taskCard) {
        fireEvent.click(taskCard);
        expect(onTaskClick).toHaveBeenCalledWith(mockTasks[0]);
      }
    });
  });

  it('shows empty state when no tasks', () => {
    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={[]}
          emptyState={{
            title: 'No tasks available',
            description: 'Check back later',
          }}
        />
      </TestWrapper>
    );

    expect(screen.getByText('No tasks available')).toBeInTheDocument();
  });

  it('filters tasks by status', async () => {
    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={mockTasks}
          filterStatus={['available']}
        />
      </TestWrapper>
    );

    await waitFor(() => {
      expect(screen.getByText('Follow Twitter')).toBeInTheDocument();
      expect(screen.queryByText('Complete Swap')).not.toBeInTheDocument();
    });
  });

  it('groups tasks by type', async () => {
    render(
      <TestWrapper>
        <TaskList
          widgetId="test-widget"
          tasks={mockTasks}
          groupBy="type"
          groupLabels={{
            offchain: 'Social Tasks',
            onchain: 'On-chain Tasks',
          }}
        />
      </TestWrapper>
    );

    await waitFor(() => {
      expect(screen.getByText('Social Tasks')).toBeInTheDocument();
      expect(screen.getByText('On-chain Tasks')).toBeInTheDocument();
    });
  });
});
```

### 9.2 Storybook配置

```typescript
// src/widgets/task/TaskList/TaskList.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { TaskList } from './TaskList';

const meta: Meta<typeof TaskList> = {
  title: 'Widgets/Task/TaskList',
  component: TaskList,
  tags: ['autodocs'],
  argTypes: {
    layout: {
      control: 'select',
      options: ['vertical', 'horizontal', 'grid'],
    },
    cardStyle: {
      control: 'select',
      options: ['default', 'minimal', 'elevated'],
    },
    groupBy: {
      control: 'select',
      options: ['none', 'type', 'status', 'category'],
    },
  },
  decorators: [
    (Story) => (
      <div className="max-w-2xl mx-auto p-4">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof TaskList>;

// 示例任务数据
const sampleTasks = [
  {
    id: '1',
    type: 'offchain' as const,
    template: 'x_follow',
    title: 'Follow @TaskOnXyz on Twitter',
    description: 'Follow our official Twitter account',
    points: 100,
    status: 'completed' as const,
    required: true,
  },
  {
    id: '2',
    type: 'offchain' as const,
    template: 'discord_join',
    title: 'Join Discord Server',
    description: 'Join our community Discord',
    points: 100,
    status: 'available' as const,
    required: true,
  },
  {
    id: '3',
    type: 'onchain' as const,
    template: 'swap_token',
    title: 'Complete First Swap',
    description: 'Swap at least $10 worth of tokens',
    points: 200,
    status: 'available' as const,
    required: false,
    progress: { current: 5, target: 10 },
  },
  {
    id: '4',
    type: 'onchain' as const,
    template: 'add_liquidity',
    title: 'Add Liquidity',
    description: 'Add liquidity to any pool',
    points: 500,
    status: 'locked' as const,
    required: false,
  },
];

export const Default: Story = {
  args: {
    widgetId: 'story-tasklist',
    tasks: sampleTasks,
    layout: 'vertical',
    showProgress: true,
    showPoints: true,
    showStatus: true,
  },
};

export const HorizontalLayout: Story = {
  args: {
    ...Default.args,
    layout: 'horizontal',
  },
};

export const GridLayout: Story = {
  args: {
    ...Default.args,
    layout: 'grid',
    columns: 2,
  },
};

export const GroupedByType: Story = {
  args: {
    ...Default.args,
    groupBy: 'type',
    groupLabels: {
      offchain: '🌐 Social Tasks',
      onchain: '⛓️ On-chain Tasks',
    },
  },
};

export const MinimalStyle: Story = {
  args: {
    ...Default.args,
    cardStyle: 'minimal',
    showDescription: false,
  },
};

export const Empty: Story = {
  args: {
    widgetId: 'story-tasklist-empty',
    tasks: [],
    emptyState: {
      title: 'No Tasks Available',
      description: 'Check back later for new tasks!',
    },
  },
};

export const CustomAccentColor: Story = {
  args: {
    ...Default.args,
    accentColor: '#10B981',
  },
};
```

---

## 10. 组件开发指南

### 10.1 新Widget开发流程

```markdown
## 开发新Widget的步骤

### 1. 定义类型和Schema
- 在 `types.ts` 中定义配置Schema (使用Zod)
- 定义Props接口
- 定义事件类型

### 2. 创建组件
- 使用 `@registerWidget` 装饰器注册
- 实现组件逻辑
- 使用 `useMCPClient` 获取数据
- 使用 `useEventBus` 发送事件

### 3. 添加样式
- 使用Tailwind CSS
- 支持主题变量
- 确保白标可定制

### 4. 编写测试
- 单元测试 (Vitest)
- 组件测试 (Testing Library)
- Storybook文档

### 5. 注册到Registry
- 在 `index.ts` 中导出
- 添加到懒加载配置

### 6. 更新文档
- 更新README
- 添加使用示例
- 更新LLM知识库
```

### 10.2 Widget开发模板

```typescript
// 新Widget开发模板
// src/widgets/{category}/{WidgetName}/

/**
 * 目录结构:
 * ├── index.ts           # 导出入口
 * ├── types.ts           # 类型定义
 * ├── {WidgetName}.tsx   # 主组件
 * ├── {WidgetName}.stories.tsx  # Storybook
 * ├── __tests__/
 * │   └── {WidgetName}.test.tsx
 * └── components/        # 子组件
 *     └── ...
 */

// types.ts
import { z } from 'zod';

export const MyWidgetConfigSchema = z.object({
  // 定义配置项
});

export type MyWidgetConfig = z.infer<typeof MyWidgetConfigSchema>;

export interface MyWidgetProps extends MyWidgetConfig {
  widgetId: string;
  widgetContext?: Record<string, unknown>;
  // 回调函数
}

// MyWidget.tsx
import React from 'react';
import { registerWidget } from '../../core/decorators';
import { useMCPClient } from '../../hooks/useMCP';
import { useEventBus } from '../../core/EventBus';
import { MyWidgetConfigSchema, type MyWidgetProps } from './types';

const WIDGET_METADATA = {
  type: 'my_widget',
  name: 'My Widget',
  description: '描述',
  category: 'task' as const,
  icon: 'IconName',
  version: '1.0.0',
  whitelabel: true,
  mcpDependencies: ['required-server'],
  configSchema: MyWidgetConfigSchema,
  defaultConfig: {
    // 默认配置
  },
};

export const MyWidget: React.FC<MyWidgetProps> = registerWidget(WIDGET_METADATA)((props) => {
  const { widgetId } = props;
  const { client } = useMCPClient('required-server');
  const { emit } = useEventBus();

  // 组件逻辑...

  return (
    <div className="taskon-my-widget" data-widget-id={widgetId}>
      {/* 组件内容 */}
    </div>
  );
});
```

---

## 📎 附录: Widget配置JSON完整示例

```json
{
  "campaignId": "campaign_001",
  "widgets": [
    {
      "id": "widget_tasklist_001",
      "type": "task_list",
      "position": 1,
      "visible": true,
      "config": {
        "tasks": [],
        "layout": "vertical",
        "showProgress": true,
        "showPoints": true,
        "groupBy": "type",
        "groupLabels": {
          "offchain": "🌐 Social Tasks",
          "onchain": "⛓️ On-chain Tasks"
        },
        "cardStyle": "default",
        "enableVerify": true
      }
    },
    {
      "id": "widget_leaderboard_001",
      "type": "leaderboard",
      "position": 2,
      "visible": true,
      "config": {
        "metric": "points",
        "metricLabel": "Points",
        "timeRange": "all",
        "displayCount": 10,
        "showRankChange": true,
        "showRewardTiers": true,
        "rewardTiers": [
          {
            "rankRange": [1, 1],
            "reward": { "type": "token", "amount": 500, "name": "$500 USDC" },
            "highlight": true
          },
          {
            "rankRange": [2, 10],
            "reward": { "type": "token", "amount": 50, "name": "$50 USDC" }
          }
        ],
        "variant": "default"
      }
    },
    {
      "id": "widget_milestone_001",
      "type": "milestone",
      "position": 3,
      "visible": true,
      "config": {
        "metric": "volume",
        "metricLabel": "Trading Volume",
        "thresholds": [
          { "value": 100, "reward": { "type": "points", "amount": 100, "name": "100 Points" } },
          { "value": 500, "reward": { "type": "points", "amount": 300, "name": "300 Points" } },
          { "value": 1000, "reward": { "type": "badge", "name": "Whale Trader" } }
        ],
        "displayStyle": "progress_bar"
      }
    }
  ]
}
```

---

**文档版本**: 1.0.0  
**最后更新**: 2025-01-22  
**维护者**: TaskOn Engineering Team
