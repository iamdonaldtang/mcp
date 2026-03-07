# Community 产品 B端前端开发需求文档

> 版本: v2.1 | 日期: 2026-03-07
> v2.1 变更: 基于 6 维度审计 (todolist.md) 补充全部 154 项缺失规格 — 页面生命周期 (Init/Destroy)、按钮可见/禁用条件、导航逻辑、数据变更级联、B→C 端影响、状态完整性; 新增 D21 Task Editor Modal; 解决 B32/D20 触发矛盾
> v2.0 变更: 新增 §2.3 全局交互规范、每页操作详情表、完整 Modal 字段规格、级联效果、边界条件（共补充 153 项缺失内容）
> 基于设计稿 `design/pencil-new.pen` + 现有文档 `website_frontend_requirements.md` v4.2
> 供前端/后端工程师对照实施

---

## 目录

1. [模块概述](#1-模块概述)
2. [全局架构](#2-全局架构)
3. [Hub 页面（4 状态）](#3-hub-页面4-状态)
4. [创建向导（4 步）](#4-创建向导4-步)
5. [模块管理页（9 + 2）](#5-模块管理页9--2)
6. [运营辅助页](#6-运营辅助页)
7. [Modal 弹窗](#7-modal-弹窗)
8. [侧栏架构](#8-侧栏架构)
9. [API 接口汇总](#9-api-接口汇总)
10. [状态路由策略](#10-状态路由策略)

---

## 1. 模块概述

### 1.1 产品定位

Community 是 TaskOn 的用户留存产品，核心价值是"让用户留下来"——通过沉没成本、损失厌恶、目标梯度、禀赋效应等心理机制，为零切换成本的 Web3 用户创造离开成本。

### 1.2 页面编码总览

| 分类 | 页面编码 | 数量 |
|------|---------|------|
| Hub（4 状态） | B09, B10, B11, B12 | 4 |
| 创建向导（4 步） | B13, B34, B35, B55 | 4 |
| 模块管理 | B31, B31a-B31i, B49, B50 | 12 |
| 运营辅助 | B32, B33, B54, B61 | 4 |
| Modal | D01-D11, D16-D21 | 17 |
| **合计** | | **41** |

### 1.3 品牌色

- 产品主色: `#48BB78`（绿色）
- 页面背景: `#0A0F1A`（深色主题）
- 卡片背景: `#111B27`
- 边框: `#1E293B`
- 主文本: `#F1F5F9`
- 次级文本: `#94A3B8`
- 强调文本: `#CBD5E1`

---

## 2. 全局架构

### 2.1 页面布局

所有 Community B端页面共享统一布局：

```
┌─────────────────────────────────────────────┐
│ Sidebar (240px)  │  TopBar (56px, full width) │
│                  ├────────────────────────────│
│  Logo            │  Content Area              │
│  ────            │  (fluid, padding 32-48px)  │
│  Home            │                            │
│  ─ PRODUCTS ─    │                            │
│  Quest           │                            │
│  Community ★     │                            │
│    ▸ Overview    │                            │
│    ▸ Sectors     │                            │
│    ▸ ...modules  │                            │
│  White Label     │                            │
│  Boost           │                            │
│  ────            │                            │
│  Analytics       │                            │
│  Settings        │                            │
└─────────────────────────────────────────────┘
```

### 2.2 主题规格

| 属性 | 值 |
|------|-----|
| 页面背景 | `#0A0F1A` |
| Sidebar 背景 | `#111B27`，右边框 `#1E293B` 1px |
| TopBar 背景 | `#111B27`，下边框 `#1E293B` 1px |
| 卡片背景 | `#111B27`，圆角 12px，边框 `#1E293B` 1px |
| 字体 | Inter 全局 |
| 页面标题 | 24px bold `#F1F5F9` |
| Section 标签 | 12px 600 `#94A3B8` letterSpacing 1px 大写 |
| Body 文本 | 14px normal `#94A3B8` |
| 按钮主色 | `#48BB78` fill, `#FFFFFF` text |
| 按钮次级 | `#48BB78` stroke 1px, `#48BB78` text |
| 状态徽章 | Active: `#0A2E1A`/`#16A34A`; Draft: `#1F1A08`/`#D97706`; Paused: `#2D1515`/`#DC2626` |
| 图标系统 | Material Symbols Rounded |

### 2.3 全局交互规范 (v2.0 新增)

#### 2.3.1 Toast 通知

| 类型 | 背景色 | 图标 | 持续时间 | 示例 |
|------|--------|------|---------|------|
| 成功 | `#0A2E1A` border `#16A34A` | `check_circle` | 3s auto-dismiss | "Draft saved" / "Task published" |
| 错误 | `#2D1515` border `#DC2626` | `error` | 5s (需手动关闭) | "Failed to save, please try again" |
| 警告 | `#1F1A08` border `#D97706` | `warning` | 4s auto-dismiss | "You've used 85% of your quota" |
| 信息 | `#0F1A2E` border `#3B82F6` | `info` | 3s auto-dismiss | "Link copied to clipboard" |

**规格**: 固定在视口右上角, 距顶 72px (TopBar下方), 右侧 24px, 宽度 360px, 圆角 8px, padding 12px 16px, z-index 9999。多条 toast 从上往下堆叠 (gap 8px), 最多同时显示 3 条。

#### 2.3.2 Loading 状态

| 场景 | 表现 |
|------|------|
| 页面首次加载 | Skeleton 骨架屏: 卡片区域用 `#1E293B` 圆角矩形占位, 高度匹配真实内容, 脉冲动画 1.5s ease-in-out infinite |
| 表格数据加载 | 3-5 行灰色占位行 (高度 48px), 列宽匹配真实列 |
| 按钮操作中 | 按钮内文字替换为 spinner (16px), 按钮 disabled, 宽度不变 |
| 内联操作 | 操作元素旁出现 16px spinner, 原内容保留但 opacity 0.5 |
| Modal 内提交 | 提交按钮 spinner + disabled; modal 不可关闭; 背景 overlay 不可点击 |

#### 2.3.3 空状态

| 场景 | 规格 |
|------|------|
| 表格无数据 | 居中显示: 48px 图标 (`inbox` material icon, `#475569`) + 标题 16px `#94A3B8` "No [items] yet" + 副文 14px `#64748B` + Primary CTA "+ Create First [Item]" |
| 搜索无结果 | 图标 `search_off` + "No results for '[query]'" + "Try different keywords or clear filters" + "Clear Filters" 链接 |
| 模块未启用 | 图标 (模块图标, dimmed) + "Enable [module] to get started" + "+ Enable [Module]" CTA |

#### 2.3.4 错误状态

| 场景 | 规格 |
|------|------|
| API 500 | 页面级: 图标 `cloud_off` 64px + "Something went wrong" + "Our team has been notified. Please try again." + "Retry" 按钮 |
| API 404 | 图标 `find_in_page` + "Page not found" + "Back to Community" 链接 |
| 网络断开 | 顶部固定 Banner (红色 `#2D1515`): "You're offline. Changes will sync when reconnected." + 不自动消失 |
| 权限不足 | 图标 `lock` + "You don't have access to this feature" + "Upgrade Plan →" 链接 |

#### 2.3.5 确认 Dialog

**通用规格**: Modal 居中, 宽度 420px, 圆角 12px, 背景 `#111B27`, border `#1E293B`.

| 元素 | 规格 |
|------|------|
| 标题 | 18px bold `#F1F5F9` |
| 正文 | 14px `#94A3B8`, max 3 行 |
| 取消按钮 | 左侧, `#94A3B8` text, hover `#CBD5E1` |
| 确认按钮 | 右侧, 危险操作用红色 `#DC2626` fill, 安全操作用绿色 `#48BB78` fill |
| 背景遮罩 | `#000000` 50% opacity, 点击遮罩 = 点击取消 |

**危险操作确认文案模板**: "Delete '[name]'? This action cannot be undone."

#### 2.3.6 拖拽排序

| 属性 | 规格 |
|------|------|
| 拖拽手柄 | `drag_indicator` icon, 16px, `#475569`, hover `#94A3B8`, cursor `grab` |
| 拖拽中 | 被拖元素: 半透明 (opacity 0.8) + 阴影 `0 4px 12px rgba(0,0,0,0.3)`; 目标位置: 2px `#48BB78` 虚线指示器 |
| 放下 | 200ms ease transition 到新位置; 立即触发 PUT reorder API; 失败则回滚并 toast error |
| 触摸支持 | 长按 300ms 激活拖拽 (移动端) |

#### 2.3.7 乐观更新 (Optimistic Update)

状态切换类操作 (toggle, status change) 一律使用乐观更新:
1. 用户操作 → 立即更新 UI
2. 发送 API 请求
3. 成功 → 无额外反馈
4. 失败 → 回滚 UI + toast error "Operation failed. Please try again."

---

## 3. Hub 页面（4 状态）

同一 URL `/community`，根据用户数据动态切换状态。

### 3.1 状态切换逻辑

| 条件 | 显示页面 |
|------|---------|
| 未创建 Community（0 modules） | B09 Empty |
| Onboarding 未完成 | B10 Guided |
| 1-3 个活跃模块 | B11 Active |
| 4+ 模块，高使用率 | B12 Deep |

---

### 3.2 B09 — Community Hub Empty

**设计稿**: Node `zzZ8D` | 尺寸: 1440×1400px

#### 页面概述
- **用途**: 新用户首次进入 Community 的欢迎页
- **入口**: 侧栏 "Community" | Dashboard 目标卡片
- **用户角色**: 未创建 Community 的项目方
- **核心目标**: 引导用户选择留存策略 → 进入创建向导

#### 页面结构

```
Content Area (padding: 48px 64px, gap: 40px)
├── Welcome Section
│   ├── Icon (community icon, 64×64, green bg #0A1F1A)
│   ├── Title: "Welcome to Community" (24px bold)
│   └── Subtitle (16px #94A3B8)
├── "RETENTION STRATEGIES" Section Label
├── Strategy Cards Row (3 cards, gap: 24px)
│   ├── Card 1: "Activate New Users" (绿色选中态)
│   ├── Card 2: "Drive Daily Engagement" (默认态)
│   └── Card 3: "Maximize Retention" (默认态)
├── "HOW IT WORKS" Section Label
├── Engine Strip (4-step flow: Quest → Activate → Engage → Retain)
├── CTA Block
│   ├── Primary: "Create Community with This Strategy" (#48BB78)
│   └── Secondary: "Or start from scratch →" (#94A3B8 link)
├── Divider (1px #1E293B)
└── Resources Row (3 cards)
    ├── "Video Tutorial"
    ├── "Retention Playbook"
    └── "Learn More" → M04 (ext, new tab)
```

#### Strategy Card 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | `activate` / `engagement` / `retention` |
| title | string | 策略名称 |
| description | string | 3 行描述 |
| icon | string | Material icon name |
| iconColor | string | `#48BB78` / `#ED8936` / `#9B7EE0` |
| iconBg | string | `#0A1F1A` / `#1F1508` / `#1A1033` |
| modules | string[] | 预配置模块列表 |
| metric | string | 关键指标名 |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 并行调用: `GET /api/community/status` (路由守卫, 见 C-05)
- **策略卡片数据来源**: 硬编码前端常量 (非 API) — 3 种策略的 title/description/modules/metric 数据写在前端代码中; 无 API 请求; 无缓存/刷新需求
- 加载态: 页面框架 (sidebar + header) 立即渲染; 内容区 skeleton (3 卡片占位 + engine strip 占位); status API 返回后渲染完整内容
- 路由守卫: 已有 Community → 重定向至 B10/B11/B12 (§10.1)

**页面离开 (Destroy)**:
- 选中策略卡片状态: **不持久化** — 仅存在于组件本地 state; 用户返回时默认预选第一张
- 无 WebSocket/timer 需清理

#### 交互逻辑

1. **策略卡片选择**: 单选互斥，选中态显示绿色边框 2px `#48BB78` + 展开 includes 详情
2. **CTA 按钮**: 携带选中策略 ID → `B13 /community/create?template={strategy_id}`
3. **"Or start from scratch"**: → `B13 /community/create?template=blank`
4. **Engine Strip**: 纯展示，4 个步骤框 (Quest/Activate/Engage/Retain) 用箭头连接

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Create Community with This Strategy" | → B13 `/community/create?template={id}` |
| "Or start from scratch →" | → B13 `/community/create?template=blank` |
| "Video Tutorial" | → (ext) 视频 |
| "Retention Playbook" | → (ext) 帮助中心 |
| "Learn More" → | → M04 (ext, new tab) |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-01 | 策略卡片选择 | click card | — (本地状态) | 选中卡片: 绿色 border 2px `#48BB78` + 展开 includes; 其他卡片取消选中; **默认预选第一张**; 再次点击同一卡片**不取消选中** (始终有一个选中) | — | — |
| C-02 | CTA "Create Community" | click | — | 携带 `?template={id}` 跳转 B13; **无策略选中时 CTA disabled** (opacity 0.5 + cursor not-allowed + tooltip "Select a strategy to continue") | — | — |
| C-03 | Engine Strip 步骤 hover | hover | — | 显示 tooltip: "Step 1: Quest — Acquire users" / "Step 2: Activate — Complete onboarding" / "Step 3: Engage — Drive daily habits" / "Step 4: Retain — Create leaving costs" | — | — |
| C-04 | 页面加载 | page mount | `GET /api/community/status` | Skeleton: 3 strategy cards + engine strip + resource cards 占位 (见 §2.3.2) | — | 500 → 全局错误页 |
| C-05 | 已有 Community 路由守卫 | route enter | `GET /api/community/status` | **若已有 Community → 302 重定向到 B10/B11/B12** (基于 §10.1 状态判断); 此页仅 `status=none` 时可访问 | — | — |
| C-06 | Resource 卡片点击 | click card | — | "Video Tutorial" → 固定外部 URL (YouTube/Vimeo); "Retention Playbook" → 固定外部 URL (help center); URL 从环境变量/CMS 配置取, 非硬编码 | — | — |

---

### 3.3 B10 — Community Hub Guided

**设计稿**: Node `S1EIA` | 尺寸: 1440×1650px

#### 页面概述
- **用途**: 创建后的引导工作区，帮助完成 Onboarding
- **入口**: 发布 Community 后自动跳转
- **用户角色**: 已创建但 Onboarding 未完成的项目方
- **核心目标**: 按 Checklist 完成设置 → 分享 → 获得首批用户

#### 页面结构

```
Content Area (padding: 32px 48px, gap: 32px)
├── Header
│   ├── Title: "Getting Started" (24px bold)
│   ├── Subtitle: "Complete these steps to get the most out of your community"
│   └── Badge: "Setting Up" (amber #1F1A08 bg)
├── Checklist Card (#111B27, rounded 12px, padding 24px)
│   ├── Progress Header: "Getting Started" + "3 of 9 complete"
│   ├── Progress Bar (green #48BB78)
│   ├── ── COMPLETED BY WIZARD ──
│   │   ├── ✅ Community created with strategy
│   │   ├── ✅ 3 starter tasks live
│   │   └── ✅ Points & Levels configured
│   ├── ── ENRICH YOUR COMMUNITY ──
│   │   ├── ○ Add more tasks (expandable, with hint)
│   │   ├── ○ Set up your Benefits Shop (expandable)
│   │   └── ○ Customize DayChain rewards (expandable)
│   ├── ── GO LIVE ──
│   │   ├── ○ Preview your community as a user
│   │   ├── ○ Share with your community (expanded)
│   │   │   ├── Share Link + "Copy" button
│   │   │   ├── Social buttons: Twitter / Discord / Telegram
│   │   │   └── "Generate Promo Kit" button → D19
│   │   └── ○ First 10 participants (auto-detect)
├── "ACTIVE MODULES" Label
├── Module Cards Row (3 cards: Tasks / Points / Leaderboard)
│   └── Each: icon + name + description + "Configure" / "Manage" button
├── "ADD MORE MODULES" Label
├── Add Module Row (4 items: TaskChain / DayChain / Milestones / Benefits Shop)
│   └── Each: icon + name + brief + "+ Enable" link
├── "Browse Configuration Templates →" Link
├── Divider
├── "RESOURCES" Label
└── Resources Row (3 cards: Community Playbook / Points Strategy / Learn More)
```

#### Checklist 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 唯一标识 |
| section | enum | `wizard` / `enrich` / `community_go_live` |
| label | string | 步骤描述 |
| status | enum | `completed` / `in_progress` / `pending` |
| expandable | boolean | 是否可展开详情 |
| hint | string? | 引导提示文案 |
| autoDetect | boolean | 是否 WebSocket 自动检测完成 |
| action | object? | `{ type: 'link' \| 'api' \| 'auto', target: string }` |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 并行调用:
  - `GET /api/community/onboarding/progress` — Checklist 进度 (C-20)
  - `GET /api/community/modules` — 模块列表和状态 (用于 ACTIVE MODULES + ADD MORE MODULES)
  - WebSocket `/ws/community/participants` — "First 10 participants" 实时计数
- 加载态: Checklist 区域 skeleton (5 行占位); Module 卡片 skeleton (3 + 4 卡片占位)
- 路由守卫: onboarding 已完成 → 重定向至 B11/B12

**页面离开 (Destroy)**:
- WebSocket `/ws/community/participants`: **断开连接** (调用 `ws.close()`)
- Checklist 展开状态: 不持久化 (下次进入默认状态)
- 无表单脏数据需处理

**模块卡片按钮标签逻辑**:
- ACTIVE MODULES 区的卡片: 模块 status=`active` 且有配置数据 (≥1 条实例) → **"Manage"**; 模块 status=`active` 但无配置数据 (0 条实例) → **"Configure"**
- ADD MORE MODULES 区: 所有未启用模块显示 **"+ Enable"**

#### 交互逻辑

1. **Checklist 展开/折叠**: 每个 enrich/go_live 项可点击展开显示详情和操作按钮
2. **Progress Bar**: 自动计算 `completedCount / totalCount`，绿色填充
3. **Share 步骤**: 展开后显示链接复制、社交分享按钮、Promo Kit 生成器
4. **First 10 participants**: WebSocket 自动检测，实时更新计数; **初始连接超时**: 5s 内未建立连接 → 降级为 polling (`GET /api/community/stats` 每 30s)
5. **Module 卡片**: 已激活模块显示绿色边框 + "Manage"/"Configure" 按钮；未激活显示 "+ Enable"
6. **Auto-complete items**: wizard 阶段步骤在创建完成后自动标记 ✅

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| 模块卡片 "Configure" / "Manage" | → B31a-B31i（对应模块管理页）|
| "+ Enable" 模块 | (API) 启用模块 → 刷新页面 |
| "Browse Configuration Templates →" | → B09 (策略浏览模式, 非创建向导) |
| "Preview your community" | → B33 (Preview Mode) |
| Share "Copy" | (action) 复制链接 |
| Social buttons | (ext) 预填分享链接 |
| "Generate Promo Kit" | → D19 Modal |
| Resources links | → (ext) 帮助中心 / M04 |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-07 | Checklist 展开/折叠 | click **整行** (非仅箭头) | — | 箭头 rotate 180°; 内容区 slide-down 200ms; **同时只展开一项** (accordion 模式, 展开新项自动关闭旧项) | — | — |
| C-08 | "Add more tasks" 展开内容 | expand | — | 显示: 当前任务数摘要 ("3 tasks configured") + 已有任务列表 (name + points, readonly) + "Go to Sectors & Tasks →" 按钮 → B31 | — | — |
| C-09 | "Set up Benefits Shop" 展开内容 | expand | — | 显示: 状态 ("Not configured yet") + 简述 "Add rewards your community members can redeem with points" + "Open Benefits Shop →" 按钮 → B31g | — | — |
| C-10 | "Customize DayChain rewards" 展开内容 | expand | — | 显示: 当前 DayChain 配置摘要 (base reward + milestones) + "Configure DayChain →" 按钮 → B31c | — | — |
| C-11 | "Preview" 展开内容 | expand | — | 显示: "See your community as a participant" + 内嵌缩略图 (200px 高度, C端首页截图) + "Open Full Preview →" 按钮 → B33 (新标签) | — | — |
| C-12 | Share Twitter | click | — | 打开 Twitter intent URL: `https://twitter.com/intent/tweet?text={encodedText}&url={communityUrl}`. 预填文案模板: "Join our community on @TaskOnXYZ! Complete tasks, earn points, and level up. {url} #Web3 #TaskOn" | — | — |
| C-13 | Share Discord/Telegram | click | — | **Discord**: 复制预填文案到剪贴板 + toast "Message copied! Paste in Discord" (不打开客户端). **Telegram**: 打开 `https://t.me/share/url?url={url}&text={text}` | — | — |
| C-14 | "+ Enable" 模块 | click | `PUT /api/community/modules/{type}/enable` | 按钮 → spinner → 成功: 卡片从 ADD MORE MODULES 移到 ACTIVE MODULES 区 (slide 动画 300ms) + toast "Module enabled" + 按钮变 "Configure". **页面刷新范围**: ① ACTIVE MODULES section 重新渲染 (新增卡片); ② ADD MORE MODULES section 移除该模块; ③ Checklist (C-16/C-17) 自动更新 — 若新模块关联 checklist 步骤 (如 DayChain → "Customize DayChain"), 步骤出现 + progress bar 分母增加; ④ Sidebar 子菜单刷新 (新模块项出现). **C端影响**: 模块启用后 C端 tab **不会立即出现** — 需先在模块管理页 (B31a-B31i) 配置并发布至少 1 个实例, 才按 §10.3 规则显示 C端 tab | Checklist 自动更新; Sidebar 刷新; C端需配置后才显示 | toast "Failed to enable module" |
| C-15 | "First 10 participants" 自动计数 | WebSocket `/ws/community/participants` | — | 实时计数: "0/10" → "5/10" → "10/10"; 数字变化时 scale 动画 (1.0 → 1.2 → 1.0, 300ms); 10/10 时自动标记 ✅ + confetti 动画. **WebSocket 断连**: 每 5s 重试, 最多 3 次; 3 次失败后显示 "Unable to track live — Refresh to update" | 步骤自动标记完成 | 降级为 polling GET /api/community/stats 每 30s |
| C-16 | 步骤自动 vs 手动完成 | — | — | **自动检测**: ✅ Community created (wizard 完成即标记) / ✅ Starter tasks live (wizard 完成即标记) / ✅ Points configured (wizard 完成即标记) / ○ First 10 participants (WebSocket 自动). **手动完成**: ○ Add more tasks (检测 task count > 3) / ○ Set up Shop (检测 shop items > 0) / ○ Customize DayChain (检测 DayChain configured) / ○ Preview (检测 preview visited flag) / ○ Share (检测 share click event) | — | — |
| C-17 | Progress bar 更新 | on step status change | — | 公式: `completedCount / totalCount * 100%`; 绿色 `#48BB78` 填充; 更新触发: 每次步骤状态从 pending → completed 时, bar 动画 300ms ease | — | — |
| C-18 | 全部步骤完成 | auto | — | Progress = 100% → Banner "Congratulations! Your community is fully set up." (green bg) + "Go to Dashboard →" 按钮; **不自动切换页面** (用户手动导航或刷新后由 §10.1 路由到 B11) | 下次进入 /community 时路由到 B11 Active | — |
| C-19 | "Browse Configuration Templates →" | click | — | → B09 `/community?view=strategies` (返回 Empty Hub 的策略选择视图, 以浏览模式展示 3 种策略模板); **注意**: 不跳转 B13 Wizard (用户已有 Community, 不应重新进入创建流程); 此链接的目的是让用户了解其他策略的模块配置组合, 以便在当前 Community 中手动启用对应模块; **可见条件**: 仅当用户当前启用模块数 < 9 时显示 (已全部启用则无需浏览) | — | — |
| C-20 | 已完成步骤数据来源 | page mount | `GET /api/community/onboarding/progress` | 数据从 **API** 取 (非 localStorage), 返回 `{ steps: [{ id, status, completedAt }] }`; 首次加载显示 skeleton checklist (5 行占位) | — | 500 → 重试按钮 |

---

### 3.4 B11 — Community Hub Active

**设计稿**: Node `vFRHi` | 尺寸: 1440×1350px

#### 页面概述
- **用途**: 社区运营中的日常管理视图
- **入口**: 侧栏 "Community"
- **用户角色**: 1-3 个活跃模块的项目方
- **核心目标**: 监控模块表现 → 快速操作

#### 页面结构

```
Content Area (padding: 32px 48px, gap: 32px)
├── Header
│   ├── Title: "My Community" (24px bold)
│   ├── Subtitle: "1,247 members · 3 modules active · Strategy: Activate New Users"
│   └── Badge: "Active" (#0A2E1A bg, #16A34A text)
├── Quick Stats Row (4 cards, gap: 16px)
│   ├── "Total Members" — 1,247 ↑12%
│   ├── "Active This Week" — 342 ↑6%
│   ├── "Points Distributed" — 24,850 ↑16%
│   └── "Tasks Completed" — 8,931 ↑23%
├── Checklist Banner (amber, collapsed: "4/5 remaining... Get your first 10 participants!")
├── "MODULE PERFORMANCE" Label
├── Module Performance Cards Row (3 cards)
│   ├── Tasks: completions/month, unique, trend + "Manage" button
│   ├── Points: earned this week, avg/user + "Manage" button
│   └── Leaderboard: active participants + "View" button
├── "ADD MORE MODULES" Label
├── Add Module Row (4: TaskChain / DayChain / Milestones / Benefits Shop)
├── "QUICK ACTIONS" Label
├── Quick Actions Row (3 buttons)
│   ├── "Create Task" → B31
│   ├── "Add Reward" → B31g (Benefits Shop)
│   └── "View Analytics" → B54 (Community Insights)
├── Divider
├── "RESOURCES" Label
└── Resources Row (2 cards)
```

#### Quick Stats 数据模型

| 字段 | 类型 | 来源 |
|------|------|------|
| totalMembers | number | `/api/community/stats` |
| activeThisWeek | number | `/api/community/stats` |
| pointsDistributed | number | `/api/community/stats` |
| tasksCompleted | number | `/api/community/stats` |
| *Trend (each)* | number (%) | 同 API，与上一周期对比 |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- **并行**调用: `GET /api/community/stats` (C-26, 60s cache) + `GET /api/community/modules/{type}/stats` per active module (C-24) + `GET /api/community/onboarding/progress` (Checklist Banner)
- 加载态: Stats Row 4 skeleton 卡片; Module Performance skeleton 卡片; Checklist Banner skeleton
- **路由守卫**: `GET /api/community/status` — 若 Community 已被其他管理员删除/暂停 → 显示全局错误页 "This community has been paused or deleted" + "Contact Admin" 链接

**页面离开 (Destroy)**:
- 无 WebSocket/timer 需清理
- 筛选/展开状态不持久化

#### 交互逻辑

1. **Stats 卡片**: 纯展示，绿色/红色趋势箭头; **趋势对比周期固定为 "vs last 7 days"** (不可配置); 后端返回 `trendPercent` + `trendPeriod` 字段
2. **Checklist Banner**: 可折叠，显示剩余步骤数，点击展开跳转 B10 Guided 完整 Checklist; **dismiss 后不再显示** (即使 onboarding 步骤后来全部完成也不重新出现, 因为 dismiss 是持久化的 C-23); onboarding 全部完成时 banner 自动消失 (无需 dismiss)
3. **Module Cards**: 每个显示关键指标 + 趋势 + "Manage" CTA; **模块无数据态** (0 completions/0 participants): 卡片显示 "0" + 副文 "No activity yet" (灰色斜体), 趋势箭头不显示
4. **Quick Actions**: 3 个快捷按钮直达常用操作; "Create Task" → B31 (仅导航, **不自动打开 D21 modal**); 到达 B31 后用户手动点击 "+ New Task"
5. **"ADD MORE MODULES" section**: 当用户已启用全部 9 个可选模块时, **整个 section 隐藏** (不显示空的 ADD MORE MODULES 标签)

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Module "Manage" / "View" | → B31a-B31i |
| "+ Enable" 模块 | (API) 启用 → 同 B10 C-14 逻辑: Module Performance 区新增卡片 + ADD MORE MODULES 移除 |
| "Create Task" | → B31 Sectors & Tasks (仅导航) |
| "Add Reward" | → B31g Benefits Shop |
| "View Analytics" | → B54 Community Insights |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-21 | Stats 卡片趋势 hover | hover 趋势箭头 | — | Tooltip: "vs last 7 days" (或 "vs last 30 days" 取决于 stats 周期); 格式: "+12% vs last 7 days" | — | — |
| C-22 | Checklist Banner 点击 | click banner body | — | **导航到 B10 Guided** (整页跳转, 非 in-place accordion); banner 显示: "4/5 remaining · Get your first 10 participants!" | — | — |
| C-23 | Checklist Banner "×" 关闭 | click × | `PUT /api/community/onboarding/dismiss` | Banner slide-up 消失; **持久关闭** (API 存储 dismissed flag); 关闭后本 session 及后续 session 不再显示; Onboarding 全部完成后 banner 自动消失不需要关闭 | — | 失败 → 静默 (降级为 session-only 隐藏) |
| C-24 | Module Performance 卡片指标 | render | `GET /api/community/modules/{type}/stats` | 每种模块显示不同指标: **Tasks**: Completions/month + Unique completers + Trend. **Points**: Points earned this week + Avg pts/user + Distribution. **Leaderboard**: Active participants + Top 3 concentration + Avg position change. **TaskChain**: Active chains + Completion rate + Drop-off step. **DayChain**: Active streak rate + Avg streak days + Day 7 pass-through. **Badges**: Badges earned (week) + Unique holders + Earn rate. **LB Sprint**: Active participants + Tasks completed + Time remaining. **Milestones**: Completions + Claim rate + Next milestone. **Shop**: Redemptions (week) + Points spent + Items sold out. **Lucky Wheel**: Spins today + Win rate + Points consumed | — | 指标加载失败 → 卡片显示 "—" |
| C-25 | Module Performance 卡片点击 | click card body (非按钮区) | — | **整个卡片可点击**, 跳转对应模块管理页 (B31a-B31i); cursor pointer; hover 效果: 背景 `#161F2E` | — | — |
| C-26 | Quick Stats 刷新 | — | `GET /api/community/stats` | **无手动刷新按钮**; 页面加载时获取一次; 60s cache; 切换 tab 回来时若超过 60s 则重新获取; 无 WebSocket 推送 (成本不值得) | — | — |

---

### 3.5 B12 — Community Hub Deep

**设计稿**: Node `TQR51` | 尺寸: 1440×1400px

#### 页面概述
- **用途**: 高级管理视图，展示所有模块 + 分析 + 集成
- **入口**: 侧栏 "Community"（4+ 模块时自动切换）
- **用户角色**: 深度使用的项目方（4+ 模块）

#### 页面结构

```
Content Area (padding: 32px 48px, gap: 32px)
├── Header + "Active" Badge
├── Quick Stats Row (4 cards)
├── "ACTIVE MODULES" Label
├── Module Cards Row 1 (3 cards: Tasks / Points / Leaderboard)
├── Module Cards Row 2 (2 cards: TaskChain / DayChain)
├── AI Insights & Suggestions Card (blue border #3B82F6)
│   ├── "DayChain streaks drop 34% at Day 7" — warning
│   ├── "Benefits Shop items consumed for 22% of users" — suggestion
│   └── "Leaderboard Sprint drives 2.3x engagement" — insight
├── "INTEGRATIONS" Label
├── Integrations Row (4 cards: Twitter/Discord/Telegram + "All Integrations →")
│   └── Connected: green border; Available: default border
├── "ENGAGEMENT OVERVIEW" Label
└── Analytics Row (2 cards)
    ├── Weekly Active Users chart (bar chart)
    └── Retention Metrics (7-day, 30-day, D1/MAU, link to full analytics)
```

#### AI Insights 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 唯一标识 |
| type | enum | `warning` / `suggestion` / `insight` |
| title | string | 简述 |
| description | string | 详细建议 |
| bgColor | string | warning=#1F0D0D, suggestion=#1F1A08, insight=#0F1A2E |
| actionUrl | string? | 可选跳转链接 |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- **并行**调用 (4 组):
  - `GET /api/community/stats` — Quick Stats (60s cache)
  - `GET /api/community/insights/ai` — AI Insights (300s cache, C-28)
  - `GET /api/community/integrations` — Integration 状态列表 (60s cache)
  - `GET /api/community/insights/overview` — WAU chart + Retention 数据 (60s cache); **注意**: WAU chart 使用 insights overview API (非 stats API), 因为需要 weekly 分组数据
- 加载态: **每个 section 独立 skeleton** — Stats Row: 4 卡片 skeleton; Module Cards: 3-5 卡片 skeleton; AI Insights: 3 行 skeleton; Integrations: 4 卡片 skeleton; Analytics: 2 chart skeleton
- 任一 API 失败: 对应 section 显示 §2.3.4 错误态 + "Retry" 按钮; 其他 section 正常渲染

**页面离开 (Destroy)**:
- 无 WebSocket/timer 需清理
- AI Insights dismiss 状态已通过 API 持久化 (C-27)

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Module cards | → B31a-B31i |
| "All Integrations →" | → B61 Integration Center |
| Integration "Configure" | → B61 对应集成配置 |
| Retention "Full Analytics →" | → B54 Community Insights |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-27 | AI Insights 单条 dismiss | click × on insight card | `PUT /api/community/insights/ai/{id}/dismiss` | 卡片 fade-out 200ms; **持久 dismiss** (API 记录, 该条洞察不再出现); "Dismiss All" 按钮在 insights 卡片右上角: **仅当 ≥2 条 insight 时显示** (1 条时仅显示单条 ×); 批量 dismiss 当前所有 → 确认 dialog "Dismiss all insights? New insights will appear as they're generated.". **全部 dismiss 后**: AI Insights section **整体折叠** — 显示灰色占位: "No active insights. New insights will appear after the next analysis cycle (daily)." (不移除 section, 保留空间位置) | — | 失败 → session-only 隐藏 |
| C-28 | AI Insights 加载 | page mount | `GET /api/community/insights/ai` | 页面加载取最新; **缓存 5 分钟**; 无手动刷新按钮 (AI 每 24h 生成新批次); 新的 insight 出现时卡片有微妙 pulse 动画 | — | 失败 → insights 区块不显示 (不影响页面其他部分) |
| C-29 | AI Insight action link | click "View DayChain Cliff →" | — | 跳转目标: actionUrl 字段指定的页面 + 锚点. 示例: "View DayChain Cliff →" → B31c `/community/modules/daychain#streak-distribution`; "Review Shop Items →" → B31g; "Create LB Sprint →" → D05 | — | — |
| C-30 | WAU chart hover | hover bar | — | Tooltip 格式: "Week of Mar 1\nActive Users: 342\n↑12% vs previous week"; **点击 bar**: 跳转 B54 并传 URL 参数 `?dateRange={weekStart},{weekEnd}` 预选对应周 | — | — |
| C-31 | Retention "Full Analytics →" | click | — | → B54 `/community/insights?tab=retention` (带 tab 参数预选 retention tab) | — | — |
| C-32 | Integration error 状态 | render | — | 错误集成卡片: 红色边框 `#DC2626` 1px + `error` icon 红色 + 错误描述 (如 "Token expired"); "Fix Connection" 按钮 → 跳转 B61 对应集成配置 + 自动发起 re-auth 流程; **仅手动修复**, 无自动重连/重试机制 (避免 API 滥用); "Reconnect" 按钮样式: 红色 outlined, 14px | — | — |
| C-33 | Integration "Configure" | click | — | → B61 `/community/integrations#{integration_id}` (带锚点, B61 页面加载后 scrollIntoView 到对应集成卡片 + 短暂高亮动画 500ms) | — | — |

---

## 4. 创建向导（4 步）

### 4.1 向导流程

```
B13 Customize → B34 Modules → B35 Quick Setup → B55 Preview & Publish
```

### 4.2 全局向导模板

```
┌──────────────────────────────────────────────┐
│ Sidebar │ TopBar: "Create Community" + "Save Draft" │
│         ├─────────────────────────────────────│
│         │ Stepper: ①Customize ②Modules ③Quick Setup ④Preview │
│         ├─────────────────────────────────────│
│         │ Step Content (2-column layout)       │
│         │  Left: form/config      Right: preview │
│         ├─────────────────────────────────────│
│         │ Action Bar: [Back] [Next: xxx] (green) │
└──────────────────────────────────────────────┘
```

**Stepper 组件规格**:
- 4 步圆点 + 连线，完成态绿色 ✓，当前态绿色实心，未来态灰色空心
- 步骤名称 13px Inter `#94A3B8` (inactive) / `#48BB78` (active/done)

---

### 4.3 B13 — Step 1: Customize

**设计稿**: Node `Gzpeu` | 编码: B13

#### 页面概述
- **URL**: `/community/create` (step 1)
- **功能**: 设置社区基础信息（名称、描述、品牌色）

#### 左侧表单

| 字段 | 类型 | 必填 | 验证规则 |
|------|------|------|---------|
| Community Name | text input | ✅ | 3-50 字符 |
| Description | textarea | ✅ | 10-500 字符 |
| Brand Color | color picker (8 预设 + 自定义) | ✅ | 有效 hex |

**预设颜色**: `#48BB78`(绿), `#5D7EF1`(蓝), `#9B7EE0`(紫), `#ED8936`(橙), `#ECC94B`(黄), `#F56565`(红), `#38B2AC`(青), `#63B3ED`(浅蓝)

#### 右侧预览

实时预览卡片，展示：
- 社区名称 + 品牌色背景
- 统计占位 (0 Members / 0 Tasks / 0 Levels)
- "Active Modules" 占位

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Save Draft" | (API) `POST /api/community/drafts` |
| "Cancel" / "Back" | → B09/B10/B11/B12 (返回 Hub) |
| "Next: Modules" | → B34 (验证表单后) |

#### 字段验证规则 (v2.0 新增)

| 字段 | 类型 | 必填 | 验证规则 | 错误文案 | UI 反馈 |
|------|------|------|---------|---------|---------|
| Community Name | text | ✅ | 3-50 字符; trim 后判断; 不允许纯空格; 允许中英文/数字/特殊字符 | 空: "Community name is required" / 过短: "Name must be at least 3 characters" / 过长: "Name must be 50 characters or less" | 实时字符计数: 输入框右下角 `{current}/50`, 临近上限 (>40) 变 amber, 超限变红; 错误时 border 变红 `#DC2626` |
| Description | textarea | ✅ | 10-500 字符; **plain text only** (不支持富文本/Markdown); 允许换行 (textarea); trim 后判断 | 空: "Description is required" / 过短: "Description must be at least 10 characters" / 过长: "Description must be 500 characters or less" | 实时字符计数同上, `{current}/500`; 4 行高度, 可拖拽调整高度 |
| Brand Color | color picker | ✅ | 有效 hex 格式: `#` + 6 位 hex 字符 (0-9, A-F, a-f) | "Please enter a valid hex color (e.g., #48BB78)" | 预设色块: 选中显示白色圆环 ring 2px + scale 1.1; "Custom" 选项展开 hex 输入框; 输入实时验证 + 颜色预览圆形色块即时更新; 非法值: 输入框红色 border + 预览色块不变 |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-34/35 | 字段输入 | typing | — | 实时字符计数更新; 见验证规则表 | — | — |
| C-36/37 | Brand Color 选择/自定义 | click 预设 / 输入 hex | — | 选中预设: 白色 ring + 右侧预览更新. 自定义 hex: 输入框实时验证 (onChange) + 预览色块更新 (有效值时) + 颜色选中为 "Custom" | — | — |
| C-38 | 右侧 Preview 更新 | any form change | — | **所有 3 个字段** (name, description, color) 变化均触发预览实时刷新; debounce 200ms; 预览内社区名称、品牌色背景、描述文本同步更新 | — | — |
| C-39 | "Save Draft" | click | `POST /api/community/drafts` | 按钮 → spinner → 成功: toast "Draft saved" → 失败: toast "Failed to save, please try again" | — | toast error |
| C-40 | 自动保存 | 每 30s (form dirty) | `POST /api/community/drafts` | **静默保存 (无 toast)**; 保存中 topbar "Save Draft" 按钮显示小 spinner; 保存成功后 spinner 消失; 30s 内无修改不触发 | — | 静默失败 (不打断用户) |
| C-41 | "Next: Modules" | click | `POST /api/community/drafts` (先保存) | ① 前端验证 → ② 失败: scroll to 第一个错误字段 + focus + shake 动画; ③ 通过: 保存 draft + 跳转 B34 | — | 保存失败 → toast error, 不跳转 |
| C-42 | 浏览器刷新/关闭 | beforeunload | — | 若 form dirty (有未保存修改): 浏览器原生确认框 "You have unsaved changes. Leave?" | — | — |
| C-43 | 恢复 Draft | page mount | `GET /api/community/drafts` | 若有 draft: 自动填充表单 (name, description, color) + topbar 显示 "Draft resumed" (绿色小标签, 3s fade); 无确认弹窗, 直接恢复. **Draft 与 URL template 冲突解决**: ① 有 draft + 有 `?template=` URL 参数 → **draft 优先** (已有进度更重要), 忽略 URL template; topbar 额外显示 "Resuming your previous draft" ② 无 draft + 有 `?template=` → 按 template 预选策略 ③ 无 draft + 无 template → 空白表单. **"Cancel"/"Back" 目标**: 通过 `document.referrer` 或 `sessionStorage` 记录来源页; 优先返回来源页; 无来源信息时 → 按 §10.1 路由到当前 Hub 状态 (B09/B10/B11/B12). **"Save Draft" 按钮**: **始终可用** (不依赖 form dirty) — 点击即保存当前表单状态; 内容未修改时 Save 仍发送 PUT (幂等, 无副作用) | — | 恢复失败 → 空表单 |

---

### 4.4 B34 — Step 2: Configure Modules

**设计稿**: Node `8NeyG` | 编码: B34

#### 页面概述
- **URL**: `/community/create` (step 2)
- **功能**: 基于策略预选模块，允许调整

#### 左侧：模块配置面板

分 4 个系统区块（对应 4-System Architecture）:

**1. Task Engine**
| 模块 | 默认状态 | 标签 |
|------|---------|------|
| Sectors & Tasks | ✅ 开启 | `Required` (不可关闭) |
| TaskChain | 根据策略 | — |
| DayChain | 根据策略 | — |

**2. Points & Recognition**
| 模块 | 默认状态 | 标签 |
|------|---------|------|
| Points & Level | ✅ 开启 | `Required` |
| Badges | 根据策略 | — |
| Leaderboard | 根据策略 | — |

**3. Incentive Campaigns**
| 模块 | 默认状态 |
|------|---------|
| LB Sprint | 根据策略 |
| Milestone | 根据策略 |
| Lucky Wheel | 关闭 |

**4. Rewards Economy**
| 模块 | 默认状态 |
|------|---------|
| Benefits Shop | 根据策略 |

#### 策略 → 模块预选映射

| 策略 | 预选模块 |
|------|---------|
| Activate New Users | Sectors & Tasks ★, Points & Level ★, TaskChain, Leaderboard |
| Drive Daily Engagement | Sectors & Tasks ★, Points & Level ★, DayChain, Leaderboard, LB Sprint |
| Maximize Retention | Sectors & Tasks ★, Points & Level ★, DayChain, Milestones, Benefits Shop, Lucky Wheel |
| Blank | Sectors & Tasks ★, Points & Level ★ only |

#### 右侧：Summary 面板

- "Your community will include" — 列出已启用模块
- "ESTIMATED POINTS EARNED" — 基于模块类型估算积分产出

#### 交互逻辑

1. **Toggle 开关**: 每个非 Required 模块可切换启用/禁用
2. **模块展开**: 点击模块行可展开查看 C-end 预览效果描述
3. **Required 标签**: 绿色 badge，toggle 不可操作
4. **Summary 实时更新**: 切换模块后右侧面板即时刷新

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Save Draft" | (API) |
| "Back" | → B13 |
| "Next: Quick Setup" | → B35 |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-44 | 模块 Toggle 动画 | click toggle | — | 200ms ease transition; ON: green fill `#48BB78` + 圆形滑块 slide right; OFF: gray fill `#475569` + slide left | 右侧 Summary 面板实时更新 | — |
| C-45 | 模块行展开 | click 行 (非 toggle 区域) | — | 展开区域 slide-down 200ms: 显示 C-end 效果描述 (2-3 行文案) + 效果截图缩略图 (160×100px); **再次点击收起**; accordion 模式: 展开新项关闭旧项 | — | — |
| C-46 | Required 模块 toggle | hover (disabled) | — | Toggle: `opacity 0.5` + `cursor: not-allowed`; hover 显示 tooltip: "This module is required and cannot be disabled" | — | — |
| C-47 | Toggle 模块 → Summary 更新 | toggle change | — | Summary 模块列表即时增减 (slide animation); "Estimated Points Earned" 重算: 基础分 = Sectors&Tasks (100) + Points&Level (50); 每额外模块 +30-80 根据类型; 显示 "~{total} XP / user / week (estimated)" | — | — |
| C-48 | "Save Draft" / "Back" | click | `POST /api/community/drafts` | "Save Draft" 同 B13. "Back" → 跳回 B13; **wizard 共享状态** (React context / Redux), B13 表单数据保留不丢失; Back **自动保存模块选择到 draft** (静默 POST, 同 C-40 auto-save 模式). **State 丢失恢复 (browser crash)**: wizard 共享状态丢失时, 从 `GET /api/community/drafts` 恢复 — draft 包含 `{step1: {name,desc,color}, step2: {modules: [...]}, step3: {...}}`; 无 draft → 重定向到 Step 1. **直接 URL 访问 Step 2**: 若无 wizard state 且无 draft → 重定向到 `/community/create` (Step 1); 有 draft → 恢复 draft 并显示 Step 2. **模块 toggle 持久化**: 每次 toggle 变化立即写入 wizard state (内存), 30s auto-save 同步到 draft API. **模块列表**: 由前端硬编码 (9 个固定模块), 非 API 动态; 未来新增模块需前端更新 | — | — |
| C-49 | 策略来源显示 | render | — | 页面顶部: "Based on: **Activate New Users** strategy" (绿色标签) + "You can toggle any module on or off" (副文); 若 `template=blank` 则显示 "Starting from scratch — only required modules are enabled" | — | — |
| C-50 | 全部可选模块关闭 | toggle off all optional | — | **允许** (只剩 Required 模块); 右侧 Summary 显示: "Minimum configuration — only required modules enabled. Consider enabling more for better retention." (amber 文案) | — | — |

---

### 4.5 B35 — Step 3: Quick Setup

**设计稿**: Node `qknQZ` | 编码: B35

#### 页面概述
- **URL**: `/community/create` (step 3)
- **功能**: 基于已选模块自动生成模板内容，支持内联编辑

#### 左侧：Auto-generated Content

根据 Step 2 启用的模块，显示 expandable 卡片：

**Sectors & Tasks 卡片**
- 提示: "We've prepared a 'Getting Started' sector with 3 starter tasks."
- 预填任务列表:
  1. Follow @Project on Twitter — 50 XP
  2. Join Discord — 50 XP
  3. KYC Verification — 100 XP
- "Edit tasks after setup →" 链接

**Points & Level 卡片**
- 预填: Point name = XP, Daily cap = 500
- Level progression: Newcomer → Active → Contributor → Expert → Legend
- "Customize levels after setup →" 链接

**DayChain 卡片** (if enabled)
- 预填: Daily check-in reward = 10 XP
- Day 7 milestone bonus = 2x, Day 30 bonus = 5x
- "Adjust streak rewards after setup →" 链接

#### 右侧：Launch Checklist

- ✅ Sector with 3 tasks
- ✅ Points system: XP, 5 levels
- ✅ Daily check-in with streak bonuses
- "What we'll set up for you:" section
- "Things to customize after launch" section

#### 交互逻辑

1. **Expandable Cards**: 每个模块卡片可展开/折叠
2. **Inline Edit**: 任务名称、点数值可直接点击编辑（不需要跳转）
3. **"Edit after setup" links**: 提示用户发布后可进一步自定义

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Save Draft" | (API) |
| "Back" | → B34 |
| "Next: Preview & Publish" | → B55 |

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-51 | 模块卡片默认状态 | render | — | **默认全部展开** (若 ≤ 3 个模块); 4+ 模块时默认折叠, 显示摘要: 模块图标 + 名称 + 项目数量 ("3 tasks configured"); 折叠状态显示单行摘要. **自动生成内容来源**: 硬编码前端模板 (非 API 驱动); 每种策略有对应的 starter content 模板; 内容在首次进入 Step 3 时生成到 wizard state. **返回 Step 2 修改后**: 若新增模块 → 新模块卡片用默认模板内容追加; 若移除模块 → 对应卡片删除, **已编辑内容丢失** (无警告, 因为 wizard state 整体更新); 未修改的模块卡片内联编辑**保留** | — | — |
| C-52 | 任务名称内联编辑 | click 任务名称文本 | — | 文本 → contenteditable; 蓝色 border 2px `#3B82F6`; Enter 保存 (blur 也保存); Esc 取消 (恢复原值); 编辑时选中全部文本 | 右侧 Launch Checklist 不变 (仅 wizard 本地状态更新) | — |
| C-53 | XP 值内联编辑 | click XP 数字 | — | 数字 → number input (宽度 80px); min=1, max=10000; Enter/blur 保存到 wizard 本地状态; Esc 取消 | — | — |
| C-54 | 内联编辑验证 | blur/enter | — | 任务名称: 不允许空白 (trim 后 length=0 → 恢复原值 + shake 动画); XP 值: 必须为正整数 (非数字/≤0/浮点 → 恢复原值 + red flash). **内联编辑 auto-save**: 内联编辑的值保存到 wizard state (内存); 受 C-40 30s auto-save 覆盖 → 30s 后自动同步到 draft API; "Save Draft" 按钮: 手动保存所有 Step 3 内联编辑内容到 draft (包括修改的 task name/XP/新增 task) | — | — |
| C-55 | "Edit after setup →" 链接 | click | — | 文本链接, 紫色 `#9B7EE0`, hover underline; 点击: **注意 — Community 尚未发布时, 模块管理页不可访问 (无 Community ID)**; 因此: wizard 未发布状态下, 链接显示但点击后 → toast "Available after publishing. You can customize everything in detail once your community is live." (不导航); 已发布的 Community 编辑模式: 正常新标签打开对应页面 | — | — |
| C-56 | "+ Add Task" | click | — | 列表末尾追加空白任务行: 名称 contenteditable (自动聚焦, placeholder "New task name...") + XP 默认值 50; Enter 确认; Esc 移除空行; **最大数量限制**: wizard 阶段最多 10 个 tasks (含预填的 3 个); 达到限制: "+ Add Task" disabled + tooltip "Maximum 10 tasks in wizard. Add more after publishing." | — | — |
| C-57 | 动态卡片列表 | render | — | 仅显示 Step 2 中启用的模块对应卡片; 禁用的模块无对应卡片; **列表由 wizard 状态驱动**, 返回 Step 2 修改模块后 Step 3 即时更新 | — | — |
| C-57b | Blank 模板 Step 3 内容 | render (template=blank) | — | 当 Step 2 仅有 Required 模块 (Sectors & Tasks + Points & Level) 时, Step 3 显示: **Sectors & Tasks 卡片**: 同策略模板 — 预填 "Getting Started" sector + 3 starter tasks (Follow Twitter 50XP, Join Discord 50XP, KYC Verification 100XP); **Points & Level 卡片**: 同策略模板 — Point name=XP, 5 levels (Newcomer 0/Active 100/Contributor 500/Expert 2000/Legend 10000); **右侧 Launch Checklist**: ✅ Sector with 3 tasks + ✅ Points system: XP, 5 levels; **无额外模块卡片**. 内容均为**硬编码模板** (非 API 驱动), 存储在 wizard 本地状态; 所有预填内容支持 C-52/C-53 内联编辑 | — | — |

---

### 4.6 B55 — Step 4: Preview & Publish

**设计稿**: Node `7mVsZ` | 编码: B55

#### 页面概述
- **URL**: `/community/create` (step 4)
- **功能**: C端预览 + 发布准备检查 + 一键发布

#### 左侧：C-End Preview Mock

嵌入式 C端预览，展示用户视角：
- Community 名称 + Logo
- Tab 导航: Home / Quests / Leaderboard
- 用户状态栏: "GaLius_3042 · 0 XP · Level 1"
- 任务卡片列表（来自 Step 3 配置）
- "Powered by TaskOn" footer

**支持 Tab 切换**: Home / Quests / Leaderboard (mock data)

#### 右侧：Launch Readiness

**Readiness Checklist**:
| 项 | 状态 | 说明 |
|----|------|------|
| Community name & branding set | ✅ | 自动检测 |
| 3 modules enabled & configured | ✅ | 自动检测 |
| 3 starter tasks with XP rewards | ✅ | 自动检测 |
| 5 levels with progression path | ✅ | 自动检测 |

**Community URL**:
- 显示 `share.taskon.xyz/community/{slug}`
- "Copy" 按钮

**After Publishing 提示卡**:
- "Community will go live immediately"
- "Clone URL to share with users"
- "Add more tasks & rewards anytime"
- "Track engagement in Analytics"

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Save Draft" | (API) `POST /api/community/drafts` |
| "Back" | → B35 |
| "Publish Community" | (API) `POST /api/community/publish` → 触发 D20 → 成功后跳转 B10 |

**重要**: "Publish Community" 按钮点击后先触发 **D20 Publish Readiness Check** Modal（检查订阅状态 + Twitter 授权），全部通过后才执行发布。**发布后跳转目标**: 首次创建 → B10 Guided (onboarding 开始); 编辑已有 Community (若未来支持 re-publish) → 按 §10.1 路由到对应 Hub 状态 (B11/B12)。

#### 操作详情 (v2.0 新增)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-58 | C-end preview 内 tab 点击 | click Home/Quests/Leaderboard | — | **在 iframe/mock 内切换**, 不触发 B端路由; preview 组件内部管理自己的 tab state; 切换时 content 区域 fade transition 200ms | — | — |
| C-59 | Desktop/Mobile toggle | click toggle | — | **Desktop**: preview 宽度 100%, 无边框. **Mobile**: preview 宽度 375px, 居中显示, 圆角边框模拟手机 (border 2px `#1E293B`, border-radius 24px, padding-top 40px 模拟 notch) | — | — |
| C-60 | Community URL slug 生成 | auto (on name input) | `GET /api/community/check-slug?slug={slug}` | 基于 Community Name 自动生成: lowercase + 空格→hyphen + 移除特殊字符; 显示 "share.taskon.xyz/community/{slug}"; **重复检测**: debounce 500ms 查 API, 重复时自动加数字后缀 (-1, -2...); slug 可手动编辑. **Slug 持久化**: 自定义 slug 保存到 wizard state + draft API; **Back 后保留**: 返回 Step 3 再回来, slug 从 wizard state 恢复 (不重新生成). **Preview 数据来源**: C-end mock 数据来自 wizard state (本地): 任务列表 = Step 3 配置的 tasks; 用户身份 = 固定 mock (同 B33 C-137); Leaderboard/Shop = mock 占位数据; 不调用 C端 API | — | — |
| C-61 | "Copy" URL | click | — | `navigator.clipboard.writeText(url)` → 成功: toast "Copied!" + 按钮文字变 "Copied ✓" 2s 后恢复. 失败 fallback: 选中文本 + toast "Press Ctrl+C to copy" | — | clipboard API 失败 → fallback |
| C-62 | Readiness Checklist 自动检测 | **page mount** (自动执行, 非 Publish 点击时) | `GET /api/community/readiness` (或从 wizard state 前端计算) | 每项检测条件: ① "Community name & branding" → name + color 非空 (**自动通过**, wizard 必填). ② "Modules enabled & configured" → activeModules ≥ 1. ③ "Starter tasks with XP" → tasks.length ≥ 1 && all have points > 0. ④ "Levels configured" → levels.length ≥ 2. 全通过: 每项绿色 ✅ + "Publish" 按钮高亮. **Readiness 检测来源**: 优先从 wizard state (内存) 计算 (无需 API); API 仅用于验证服务端状态 (如 slug 可用性). **Readiness 有问题时**: 用户仍可查看 C-end preview (预览不依赖 readiness); 仅 "Publish" 按钮 disabled; 未通过项显示红色 ❌ + 修复提示 (如 "Go back to Step 3 to add tasks") | — | — |
| C-63 | "Publish Community" 完整流程 | click | — | ① 前端 readiness 全通过? → ②否: 按钮 disabled, 显示提示. ②是: 打开 D20 Modal → ③ D20 检测 (订阅+Twitter) → ④全通过: 自动关闭 D20 → `POST /api/community/publish` → ⑤ 按钮 loading spinner → ⑥ 成功: toast "Community published!" + 跳转 B10 Guided. **C端发布后初始状态**: 发布后 C端立即可访问, 包含: Home (显示配置的 sectors/tasks)、Quests tab (wizard 配置的 starter tasks)、用户状态栏 (积分 0 + Level 1); 其他 tab (Leaderboard/Sprint/Milestone/Shop) 按 §10.3 规则, 仅在对应模块有有效数据时出现; 若 wizard 仅配置了 Required 模块 (Sectors&Tasks + Points&Level), C端仅显示 Home + Quests 两个 tab | B10 Guided 自动出现, onboarding 开始; C端 URL 立即可访问 | D20 未通过 → 见 D20 规格; Publish API 失败 → 见 C-64 |
| C-64 | Publish API 失败 | — | — | 关闭 D20 → toast "Publish failed: {error_message}" (红色, 5s) → 停留在 B55; "Publish" 按钮恢复可用; 用户可重试 | — | — |
| C-65 | Readiness 全通过状态 | auto | — | 所有 4 项 ✅: "Publish Community" 按钮从 disabled (opacity 0.5) 变为 active (绿色 `#48BB78` fill, cursor pointer, hover darken) | — | — |

---

## 5. 模块管理页（9 + 2）

### 5.1 通用管理页模板

所有模块管理页共享统一结构（基于 Quest Management `XvXEQ` 模板复制）：

```
Content Area
├── Breadcrumb: "Home > Community > {Module Name}"
├── Header Row
│   ├── Title: "{Module Name}" (24px bold)
│   ├── Subtitle: "{count} items · {active_count} active"
│   └── Primary Button: "+ Create {Item}" (green)
├── Stats Row (4 cards)
│   └── Each: label (12px #94A3B8) + value (24px bold) + trend
├── Insight Banner (amber/green/red, contextual tip)
├── Filter Tabs: All | Active | Completed | Draft
├── Search Bar + Sort Dropdown
├── Data Table
│   ├── Columns: Name | Status | Key Metric | Date | Actions
│   └── Each row: click → edit modal or detail page
└── Pagination: "Showing X of Y items" + Previous / Next
```

#### 通用表格操作规格 (v2.0 新增, C-66 ~ C-79)

以下操作适用于所有使用表格模板的模块管理页 (B31a-B31i, B49, B50):

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| C-66 | 表格行 hover | hover row | — | 背景色: `#161F2E`; Actions 列按钮 **始终显示** (非 hover 才显示, 避免触屏不友好) | — |
| C-67 | 表格行点击 | click row (非 Actions 区) | — | 打开对应 edit modal (D01-D11); Actions 列按钮有独立点击区域, 不受行点击影响; row click zone: 除 Actions 列外的所有列; cursor pointer | — |
| C-68 | "Duplicate" 操作 | click Actions → Duplicate | `POST /api/community/modules/{type}/instances/{id}/duplicate` | ⋮ 菜单: Duplicate / Delete; 点击 Duplicate → 新 draft 出现在列表**顶部** + 高亮 flash 动画 1s + toast "{Name} duplicated as draft" | toast "Failed to duplicate" |
| C-69 | "Delete" 操作 | click Actions → Delete | `DELETE /api/community/modules/{type}/instances/{id}` | 确认 dialog: "Delete '{name}'? This cannot be undone." → 确认 → 行 fade-out 200ms → toast "{Name} deleted" → pagination 更新 | toast "Failed to delete" |
| C-70 | Status badge 切换 | click badge / toggle | `PUT /api/community/modules/{type}/instances/{id}` body `{status}` | **乐观更新**: 立即切换 badge 颜色 (Active↔Paused); API 失败 → 回滚 badge + toast error | 回滚 + toast |
| C-71 | Search bar | input (debounce 300ms) | `GET ...?search={query}` | 输入出现时显示 × 清空按钮; 无结果: 空状态 (§2.3.3 搜索无结果样式); search 保留在 URL query string `?q={query}` | — |
| C-72 | Filter tabs | click tab | `GET ...?status={filter}` | Active tab: 绿色下划线 2px `#48BB78` + 文字 `#F1F5F9`; **URL 参数同步**: `?status=active`; 刷新后保持筛选; "All" tab 显示总数 badge | — |
| C-73 | Sort dropdown | click dropdown | `GET ...?sort={field}&order={asc\|desc}` | 选项: Date Created (默认, desc) / Date Modified / Name (asc) / Status / 关键指标 (模块特定); 选中项显示 checkmark `✓`; 同步 URL `?sort=name&order=asc` | — |
| C-74 | Pagination | click Prev/Next / page number | `GET ...?page={n}&limit=20` | 默认每页 20 条; 显示 "Showing {start}-{end} of {total} items"; 首页 Prev disabled, 末页 Next disabled; 页码按钮: 1 2 3 ... 10 (省略中间) | — |
| C-75 | Insight Banner dismiss | click × | — | Banner slide-up 消失; **session-only** (刷新后重新出现); banner 含 action link 时: link 为文本链接, 绿色, 点击跳转对应页面 | — |
| C-76 | 列表为空 (0 条) | render | — | 空状态: 居中图标 + "No {items} yet" + 副文 + Primary CTA "+ Create First {Item}" (绿色按钮) → 对应 create modal | — |
| C-77 | 数据加载中 | page mount | — | 表格区域: 3-5 行 skeleton (灰色占位行, 高 48px, 脉冲动画); Stats Row: 4 个 skeleton 卡片 | — |
| C-78 | API 加载失败 | error | — | 表格区域替换为: 图标 `cloud_off` + "Failed to load data" + "Retry" 按钮; Retry → 重新调用 GET API | — |
| C-79 | Bulk 操作 | checkbox | `PUT /api/.../bulk` body `{ids, action}` | 行前 checkbox: 未选 → 空框, 选中 → 绿色 ✓; 表头全选 checkbox; 选中 ≥1 行后顶部出现 bulk action bar (蓝色 bg `#0F1A2E`): 显示 "{n} selected" + "Activate" / "Pause" / "Delete" 按钮; Delete 需确认 dialog | 部分失败 → toast "3 of 5 items updated, 2 failed" |

### 5.2 侧栏子菜单结构

进入模块管理后，Community 侧栏展开为带 5 个 section header 的子菜单：

```
▾ Community (expanded)
  ─ Overview          → B10/B11/B12
  TASKS
    Sectors & Tasks   → B31
    TaskChain         → B31b
    DayChain          → B31c
  POINTS & RECOGNITION
    Points & Level    → B31a
    Badges            → B31i
    Leaderboard       → B31d
  CAMPAIGNS
    LB Sprint         → B31e
    Milestones        → B31f
  REWARDS
    Benefits Shop     → B31g
    Lucky Wheel       → B31h
  SETTINGS
    Access Rules      → B49
    Homepage Editor   → B50
```

**侧栏样式规格**:
- Section headers: 非可点击, fontSize 10, fontWeight 600, fill `#94A3B8`, letterSpacing 1, 大写
- Active item: 绿色 bg `#ECFDF5` → 暗色 `#0A1F1A` + text `#48BB78` + fontWeight 600
- Inactive item: no fill + text `#94A3B8`
- Sub-item: padding `[8,12,8,40]`, fontSize 13, icons 16×16

---

### 5.3 B31 — Sectors & Tasks

**设计稿**: Node `Wug7d` | URL: `/community/sectors`

#### 页面概述
- **功能**: 管理任务分区和具体任务，支持拖拽排序

#### 独特结构（非标准表格模板）

```
Content Area
├── Header: "Sectors & Tasks" + "+ New Sector" + "+ New Task"
├── Stats Row: Total Sectors (4) / Active Tasks (12) / Draft Tasks (3) / Completions (8,931)
├── Sector Groups (expandable)
│   ├── "Getting Started" sector (Active badge)
│   │   ├── Task 1: "Follow @Project on Twitter" [Required] — 50pts — 1,967 — ☑
│   │   ├── Task 2: "Join Discord Server" [Social] — 30pts — 882 — ☑
│   │   ├── Task 3: "Complete KYC Verification" [Verification] — 100pts — 356 — ☑
│   │   └── Task 4: "Bridge Assets to Base" [On-chain] — 100pts — — ☑
│   ├── "Daily Engagement" sector (Active badge)
│   │   ├── Daily Check-in [Recurring] — 10pts — 3,621 — ☑
│   │   ├── Share Daily Market Take [Custom] — 25pts — 1,241 — ☑
│   │   └── Refer a Friend [Referral] — 150pts — 264 — ☑
│   └── "Advanced Trading" sector (Draft badge, Hidden)
│       └── ...
├── Tip Banner: "Drag sectors to reorder..."
```

#### Task 行数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 任务 ID |
| name | string | 任务名称 |
| type | enum | `social` / `onchain` / `verification` / `custom` / `recurring` / `referral` |
| status | enum | `active` / `draft` / `hidden` / `expired` |
| points | number | XP 奖励值 |
| completions | number | 完成次数 |
| sectorId | string | 所属分区 ID |
| order | number | 排序权重 |
| isRequired | boolean | 是否必做 |
| deadline | ISO8601? | 截止时间 |
| maxClaims | number? | 最大领取次数 |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/sectors?include=tasks` — 返回 sector 列表 + 每个 sector 下的 tasks (嵌套结构)
- 加载态: Stats Row 4 skeleton + sector 区域 3-5 个 skeleton 块
- URL 参数支持: `?sector={sectorId}` → 自动展开并 scrollIntoView 到指定 sector; `?task={taskId}` → 展开所在 sector + 高亮该 task 行 (500ms 黄色 flash)
- Stats Row: `GET /api/community/stats` (复用 Hub API, 只取 sectors/tasks 相关字段)

**页面离开 (Destroy)**:
- **内联编辑 (contenteditable) 未保存**: 若有未 blur 的 contenteditable 字段 (C-86 双击编辑中), 自动触发 blur → 保存 → 然后导航
- 拖拽中途离开: 取消拖拽, 不发送 reorder API
- Sector 展开/折叠状态: **通过 localStorage 持久化** (key: `community_sectors_expanded`, value: `[sectorId1, sectorId2]`); 下次进入恢复展开状态

#### 交互逻辑

1. **拖拽排序**: Sector 和 Task 都支持拖拽手柄排序 (drag handle icon)
2. **Sector 折叠/展开**: 点击 Sector header 可折叠其下任务列表
3. **任务行操作**: 每行有编辑✏️、复制、删除、显隐切换
4. **状态切换**: Task status toggle 直接在行内切换 active/draft/hidden; **"Publish Task" 按钮可见条件**: 仅 `status=draft` 的 task 行显示 "Publish" 按钮 (绿色小按钮); `active`/`hidden`/`archived` 状态不显示; **不支持批量 publish** (逐条操作)
5. **Sector 显隐**: Sector header 有 visibility toggle; **Sector 内所有 Task 为 Draft 时**: C端自动隐藏该 sector (无需手动 hide sector); 由 §10.3 规则控制 — "Quests tab: 至少 1 个已发布任务"

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "+ New Sector" | (action) 新增 sector 输入框 |
| "+ New Task" | (action) 新增 task 表单/modal |
| Task row edit ✏️ | (action/modal) 编辑任务详情 |
| "Publish Task" | (API) → 触发 D20 Readiness Check |

#### 操作详情 (v2.0 新增, C-80 ~ C-90)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-80 | Sector 拖拽排序 | drag handle (左侧 `drag_indicator`) | `PUT /api/community/sectors/reorder` body `{order: [id1,id2,...]}` | 拖拽手柄: 左侧 16px icon; 拖拽中: 整个 sector 块高亮 border `#48BB78` + 半透明; 目标位置: 2px 绿色虚线指示器; 松开 → 200ms transition | — | 失败 → 回滚排序 + toast |
| C-81 | Task 同 Sector 内排序 | drag handle | `PUT /api/community/tasks/reorder` body `{sectorId, order: [id1,...]}` | 同 §2.3.6 拖拽规格; 仅在同一 sector 内上下拖拽 | — | 失败 → 回滚 |
| C-82 | Task 跨 Sector 拖拽 | drag to different sector | `PUT /api/community/tasks/{id}` body `{sectorId: newId, order: n}` | 拖拽到另一 Sector 区域: 目标 Sector 展开 + 虚线指示器出现; 松开 → task 从原 sector 消失, 出现在目标 sector 对应位置 | — | 失败 → 回滚 |
| C-83 | Sector header ⋮ 菜单 | click ⋮ | — | 菜单项: **Edit Name** (inline edit: header 文字 → contenteditable, Enter 保存 `PUT /api/community/sectors/{id}`, Esc 取消) / **Hide** (sector 整体 hidden, tasks 不影响其 status, `PUT ...` body `{visible:false}`) / **Delete** (需确认, 见 C-84) | Hide → C端对应 sector 不可见 | — |
| C-84 | Sector 删除非空 | click Delete (non-empty) | — | **阻止删除**: 提示 dialog "This sector contains {n} tasks. Move or delete all tasks first before deleting this sector." (OK 按钮关闭 dialog); **空 sector**: 确认 dialog → `DELETE /api/community/sectors/{id}` → 成功 → sector 消失 + toast | — | — |
| C-85 | Task 行 Actions | click icons | — | 4 个图标按钮 (16px, gap 8px): ✏️ 编辑 (→ **D21 Task Editor Modal**, 预填当前 task 数据) / ⊕ 复制 (→ `POST duplicate`, 新 draft 出现在同 sector 下方) / 👁️ 显隐切换 (`PUT {visible:false}` → task `status` 变 `hidden`; 乐观更新, hidden 态: 行 opacity 0.5; C端即时隐藏该任务) / 🗑️ 删除 (确认 dialog → `DELETE`) | — | — |
| C-86 | Task 积分值内联编辑 | **双击** points cell | `PUT /api/community/tasks/{id}` body `{points}` | 双击 → number input (宽 80px, min 1, max 10000); Enter/blur → 保存 PUT; Esc → 取消; 验证: 正整数, 范围 1-10000 | — | PUT 失败 → 恢复原值 + toast |
| C-87 | "+ New Sector" | click | — | 底部追加空白 sector: 名称输入框 (autofocus, placeholder "Sector name..."); Enter → `POST /api/community/sectors` body `{name}` → 成功: 输入框变为 sector header; Esc → 取消移除空行 | — | POST 失败 → toast + 保留输入框 |
| C-88 | "+ New Task" | click (全局按钮, header 区) | — | 打开 **D21 Task Editor Modal** (见 §7.3.17); **非每个 sector 下独立按钮**; modal 内 sector 字段默认选中当前展开的 sector (若有) | — | — |
| C-89 | "Publish Task" | click (仅 Draft 状态任务) | — | → D20 Readiness Check → 通过 → `PUT /api/community/tasks/{id}` body `{status: 'active'}` → badge 变 Active 绿色 (动画: badge scale 1.0→1.2→1.0, 300ms) + toast "Task published". **Stats 更新**: Stats Row "Active Tasks" 计数 +1 (前端乐观更新); "Draft Tasks" 计数 -1. **拖拽排序后 Stats**: 拖拽排序 (C-80/C-81/C-82) **不影响 Stats Row** (排序不改变数据量, 无需刷新) | C端即时出现该任务 | D20 未通过 → 见 D20 规格 |
| C-89b | Task Active→Hidden (👁️ toggle) C端影响 | click eye icon | — | 乐观更新: 行 opacity 0.5 + badge 变 Hidden (灰色); **C端影响即时生效** — 下一次 C端 API 请求时该任务不再返回; 已有 session 的 C端用户在页面刷新时看到变化; **不影响** 已完成该任务的用户积分/记录 | C端 Quests 列表即时移除该任务 | — |
| C-90 | 删除有完成记录的任务 | click Delete on completed task | — | **软删除 (archive)**: 确认 dialog "This task has {n} completions. It will be archived (not permanently deleted). Users who completed it will keep their points." → 确认 → `PUT {status: 'archived'}` → 行从默认视图消失 (可在 filter tab "Archived" 中查看). **C端影响**: 任务立即从 C端 Quests 列表中移除 (API 实时, 非 CDN 缓存); 已完成用户保留积分和完成记录. **注意**: 此处 "archive" 等同于 b2c_operation_mapping.md §4.2.8 所述的 "permanently removed from C-end" — 即 C端完全不可见, 但 B端数据保留用于分析和审计 | 已获得积分不回收; C端 Quests 列表即时更新 | — |

---

### 5.4 B31a — Points & Level

**设计稿**: Node `zCfKQ` | URL: `/community/modules/points`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Points Issued | 1,284,500 |
| Active Earners | 3,847 |
| Avg Points / User | 334 |
| Level-Up Events | 1,205 |

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Level Name | 等级名称（Bronze/Silver/Gold...） |
| Status | Active / Archived |
| Threshold | 积分阈值 |
| Members | 当前该等级人数 |
| Date | 创建/修改日期 |

#### Primary Action
- "+ Add Level" → D01 Points & Level Editor Modal (`8DhXJ`)

#### Insight Banner
- "Level Manipulation Alert: 3 Th... (安全提示)"

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 并行调用: `GET /api/community/point-types` (积分类型列表) + `GET /api/community/modules/points/levels` (等级表) + `GET /api/community/stats` (Stats Row 数据)
- 加载态: Stats Row skeleton + Point Types section skeleton (2 卡片) + Level table skeleton (5 行)

**页面离开 (Destroy)**: 无需清理

**空数据态**:
- 0 个 levels: Stats Row "Level-Up Events" 显示 "0"; Level 表格显示 §2.3.3 空状态 + "+ Add Level" CTA
- "+ Add Level" 按钮: **始终可用** (不依赖 point type 选择 — M-03 在 D01 内选); 但若 0 个 point types → D01 打开时 M-03 下拉为空 + 提示 "Create a point type first" + "Add Point Type" 快捷链接 (执行 C-91)

#### 操作详情 (v2.0 新增, C-91 ~ C-95)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-91 | 积分类型管理 | page top "Point Types" section | `GET/POST/PUT /api/community/point-types` | 页面顶部独立 section (在 Stats Row 下方): 标签 "POINT TYPES" + 当前类型列表 (卡片: 名称 + 符号 + 类型色); "+ Add Type" 按钮 → inline 展开: Name 输入 (必填, 唯一) + Symbol 输入 (1-5 字符, 如 "XP", "GEM") + 颜色选择 → Save. **最大数量**: 10 种积分类型; 达到限制: "+ Add Type" disabled + tooltip "Maximum 10 point types". **跨页面感知**: 新增积分类型后, B31d (Leaderboard) 和 B31e (LB Sprint) 的 Point Type 下拉数据源自动刷新 (D04 M-18 / D05 M-23 的下拉数据来自 `GET /api/community/point-types`); **无需手动刷新**, 因为 modal 打开时实时拉取. **C端影响**: 新积分类型创建后 **不立即出现在 C端** — 需配置 tasks 来 earn 该积分类型 + 配置 Leaderboard 后才在 C端 User Status Bar 显示 | 新积分类型可用于 Leaderboard/Sprint/Task 配置 | — |
| C-92 | 积分获取规则 | — | — | **不单独配置页面**; 积分获取与 Task 绑定 (每个 Task 在 B31 中指定 points 值和 point type); 此页仅管理 Level 阈值和等级体系; 页面可显示 "Points are earned through tasks. Configure task rewards in Sectors & Tasks." 提示 | — | — |
| C-93 | 新增等级 threshold 验证 | D01 modal 保存时 | — | Threshold 值必须 > 前一等级 threshold 且 < 后一等级 threshold; 保存时验证整体顺序; 违反: D01 内 threshold 字段红色 border + "Threshold must be between {prev} and {next}" | — | — |
| C-94 | 修改等级 threshold (有用户) | D01 modal 保存时 | — | 若该等级有 members > 0 且 threshold 提高 → 警告 dialog: "Raising this threshold may cause {n} members to be demoted to the previous level. Continue?" → 确认 → PUT → 用户等级重新计算. **C端降级通知机制**: 服务端 batch 重算 (5 分钟内完成); 被降级用户 C端效果: ① 等级 badge 实时更新 (WebSocket 推送或下次 API 请求时); ② C端 User Status Bar 显示新等级; ③ **无主动通知** (不弹 toast/push) — 用户仅在查看个人资料或等级页时看到变化 (避免负面体验); ④ 降级后若用户不再满足 Level-gated Shop 商品门控, 商品变为 locked 状态. **B端影响**: B31a Stats Row "Level-Up Events" 不计降级 (仅统计升级); 表格 "Members" 列即时更新反映重新计算结果 | 受影响用户等级变化 + C端等级显示更新 (5 分钟内); B31a Members 列刷新 | — |
| C-95 | 删除等级 (有成员) | click Delete (有 members) | — | 确认 dialog: "Delete '{levelName}'? {n} members at this level will be reassigned to '{previousLevelName}'." → 确认 → DELETE → 成员自动降级到下一低等级; 若删除最低等级: 成员变为 "No Level" 状态. **"No Level" C端显示**: 用户状态栏 Level badge 显示 "—" (灰色 dash); 等级进度条隐藏; tooltip "No level assigned" | 用户等级重新分配; C端等级 badge 更新 (5 分钟内) | — |

---

### 5.5 B31b — TaskChain

**设计稿**: Node `lpdtp` | URL: `/community/modules/taskchain`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Chains | 2 |
| Active Chains | 1 |
| Completions | 1,847 |
| Avg Completion Rate | 72.3% |

#### 独特组件：Chain Step Funnel
- 水平漏斗图，展示每步的完成率和流失率
- 颜色: 绿色(高) → 黄色(中) → 红色(低)
- 标注: "Step 2→3: 6% drop off"

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Chain Name | 链名称 |
| Status | Active / Draft |
| Steps/Completions/Rate | 步数/完成数/完成率 |
| Date | 创建日期 |

#### Primary Action
- "+ Create Chain" → D02 TaskChain Editor (`bZiB5`)
- "Activate Chain" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 并行调用: `GET /api/community/modules/taskchain/instances` (chains 列表 + stats) — 单一 API 返回列表 + 每条 chain 的 funnel 数据; Stats Row 从同一 response 聚合
- Funnel Chart: **显示当前选中/active chain** 的步骤漏斗; 若多条 chain, 默认显示第一条 Active chain; 表格行点击切换 funnel 展示的 chain; 无 Active chain 时 funnel 显示 "Activate a chain to see funnel data"

**页面离开 (Destroy)**: 无需清理

**"Activate Chain" 按钮可见条件**:
- 仅 `status=draft` 的 chain 行显示 "Activate" (绿色小按钮)
- `status=paused` 的 chain: 显示 **"Reactivate"** 按钮 (同样触发 D20)
- `status=active` 的 chain: 显示 "Pause" toggle (C-97)

**After Activate Stats**: 激活成功后 Stats Row "Active Chains" +1 (乐观更新); Funnel Chart 切换到新激活 chain 的数据 (此时为空漏斗, 所有步骤 0%)

**Chain 引用已删除 Task**: 若 chain step 引用的 task 已被删除/archived, 表格行该 chain 显示 amber warning icon + tooltip "Step {N} references a deleted task '{name}'. Edit the chain to fix."; 该 chain 仍可保持 Active (C端显示已完成步骤, 但 deleted task 步骤标记为 "Task no longer available — auto-completed")

#### 操作详情 (v2.0 新增, C-96 ~ C-98)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-96 | Funnel 图 hover | hover step bar | — | Tooltip: "Step {N}: {stepName}\n{count} completed ({percent}% of step {N-1})\n↓{dropoff}% drop off" | — | — |
| C-97 | Chain Active → Pause | click status toggle | `PUT /api/community/modules/taskchain/instances/{id}` body `{status:'paused'}` | 确认 dialog: "Pause this chain? Users currently in progress will have their progress **frozen** (preserved but paused). They can resume when the chain is reactivated." → 确认 → 乐观更新 badge. **C端 UX**: ① TaskChain 入口仍可见 (不隐藏), 但显示 "Chain Paused" 状态横幅 (amber bg `#1F1A08`, 文案: "This chain is temporarily paused. Your progress is saved and will resume when it restarts."); ② 已完成步骤保持 ✅ 状态; ③ 下一步骤按钮变为 disabled (灰色 + "Paused"); ④ 新用户看到 chain 但无法开始 (按钮: "Coming Soon"); ⑤ Reactivate 后所有用户从暂停位置继续 | 正在进行的用户进度冻结 (不丢失); C端 TaskChain 显示 Paused 状态 | — |
| C-98 | "Activate Chain" | click | — | → D20 → 通过 → `PUT {status:'active'}` → Chain 开始对**新用户**生效; 已有进度的用户保持当前进度 (不重置); toast "Chain activated" | C端出现 TaskChain 入口 | — |

---

### 5.6 B31c — DayChain

**设计稿**: Node `fLLVb` | URL: `/community/modules/daychain`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Active Streak Rate | 68% |
| Completion Rate | 72% |
| Day 7 Pass-through | 58% |
| Avg Streak Days | 8.3 |

#### 独特组件：Streak Distribution Chart
- 柱状图，X 轴为连续天数 (1-30)，Y 轴为用户数
- 标注 "Day 7 cliff: 34% drop off"，红色高亮

#### Insight Banner
- "Streak Bottleneck — 43% of users break at Day 7. Consider adding Day 7 bonus..."

#### Primary Action
- "+ Create Chain" → D03 DayChain Config (`nEAUB`)
- "Activate Chain" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/daychain/instances` — 返回 chains 列表 + stats + streak distribution chart 数据 (chart 数据作为 response 子字段, 非单独 API)
- Streak Distribution Chart: 显示当前选中 Active DayChain 的分布; 默认第一条; 点击表格行切换
- **3+ DayChains**: 标准分页适用 (同 §5.1); 默认每页 20 条; 3 条不触发分页

**页面离开 (Destroy)**: 无需清理

**D03 Modal 保存后**: B31c streak distribution chart **自动刷新** (关闭 D03 后重新 GET API); Stats Row 也刷新

**DayChain 配置变更 C端传播**:
- 修改 Base Reward / Milestone Bonus: **即时生效** — C端用户下次 check-in 时获得新奖励值
- 修改 Daily Task: **即时生效** — C端显示新的每日任务; **不影响已有 streak** (用户当天已完成旧任务的算已完成)
- 修改 Grace Period: **即时生效** — 影响断链判定逻辑; 已在 grace period 中的用户适用新规则
- 总结: 所有配置变更即时生效 (非"下一天")

#### 操作详情 (v2.0 新增, C-99 ~ C-102)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-99 | Streak 分布图 hover | hover bar | — | Tooltip: "{count} users have maintained a {n}-day streak" | — | — |
| C-100 | "Day 7 cliff" 红色区域点击 | click highlighted zone | — | 打开 D03 DayChain Config Modal + **自动聚焦** Day 7 Milestone Bonus 字段 (scroll to + highlight border pulse 500ms) | — | — |
| C-101 | Grace Period / Streak Freeze | D03 modal 内 | — | **Grace Period**: 字段在 D03 内 (见 M-16); 0-24 hours; 0=无容忍. **Streak Freeze**: 当前版本**不支持**道具机制; 预留字段但 UI disabled + tooltip "Coming soon" | — | — |
| C-102 | 多个 DayChain 支持 | — | — | 当前设计**支持多个 DayChain** (表格列表模式); 每个 DayChain 可绑定不同的每日任务; 无数量上限但建议 ≤ 3 (Insight Banner: "Having too many DayChains may confuse users") | — | — |

---

### 5.7 B31d — Leaderboard

**设计稿**: Node `Emmab` | URL: `/community/modules/leaderboard`

#### 关键区别
- **Leaderboard ≠ LB Sprint**: 此页管理周期性排行榜（周/月/全部），**无额外激励**
- 基于自定义积分类型（EXP/GEM 等）

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Participation Rate | 74% |
| Weekly Active Contestants | 1,420 |
| Top 3 Concentration | 34% |
| Avg Position Change | ±4.2 |

#### Filter Additions
- 标准 All/Active/Archive 之外，增加 "Point Type" 下拉筛选

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Leaderboard Name | 名称 |
| Status | Active / Archived |
| Period | Weekly/Monthly/All Time |
| Participants | 参与人数 |
| Date | 创建日期 |

#### Primary Action
- "+ Create Leaderboard" → D04 Leaderboard Config (`j8UnD`)

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/leaderboard/instances` (list + stats), 并行 `GET /api/community/point-types` (C-103 filter dropdown)
- Stats 与 instances 来自同一 API response

**页面离开 (Destroy)**: 无需清理

**Archive vs Delete 关系**:
- **Archive**: 软操作 — Leaderboard 变为 Archived 状态, C端隐藏, 历史数据保留; 表格中可见 (filter "Archived" tab)
- **Delete**: 硬删除 — 从 ⋮ 菜单 "Delete" (确认 dialog); **仅 Archived 状态可 Delete** (Active leaderboard 只能先 Archive); Delete 后永久移除
- **Unarchive**: **不支持** — Archive 是单向操作; 若需恢复, 需重新创建 (相同配置)
- **After archive Stats**: "Participation Rate" 重新计算时 **排除 Archived leaderboards** (仅统计 Active); Stats Row 乐观更新

#### 操作详情 (v2.0 新增, C-103 ~ C-105)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-103 | "Point Type" 下拉 | click filter dropdown | `GET /api/community/point-types` | 数据来源: B31a 已配置的积分类型列表; 下拉选项: "All" + 各积分类型名 (EXP/GEM/自定义); 选中后过滤表格 | — | — |
| C-104 | Archive 操作 | click Actions → Archive | `PUT /api/community/modules/leaderboard/instances/{id}` body `{status:'archived'}` | 确认 dialog: "Archive this leaderboard? It will be hidden from the community page but historical data is preserved." → 确认 → badge 变 Archived (灰色) | C端 Leaderboard tab 中该条目隐藏; 历史数据保留可查 | — |
| C-105 | 多 Leaderboard 同积分类型 | create | — | **允许**: 同一积分类型可创建多个 Leaderboard (不同 Period); 例如 EXP Weekly + EXP Monthly; **最大数量**: 每种积分类型 ≤ 5 个 Leaderboard; 超出: "+ Create" 按钮 disabled + tooltip "Maximum 5 leaderboards per point type" | — | — |

---

### 5.8 B31e — LB Sprint

**设计稿**: Node `FO9JR` | URL: `/community/modules/lb-sprint`

#### 关键区别
- **LB Sprint**: 限时排行榜竞赛，有开始/结束日期，附带**非积分类激励**（NFT/Token/WL Spot）

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total LB Sprints | 2 |
| Active Participants | 3,421 |
| Tasks Completed | 8,932 |
| NFTs + Rewards Given | 12 NFTs |

#### Data Table 列

| 列名 | 说明 |
|------|------|
| LB Sprint Name | 名称 |
| Status | Active / Completed / Draft |
| Reward | 奖品类型（$50 USDT / NFT...） |
| Duration | 起止日期 |
| Participants | 参与人数 |

#### Primary Action
- "+ Create LB Sprint" → D05 LB Sprint Editor (`NnzO9`)
- "Launch Sprint" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/lb-sprint/instances` (list + stats)
- **Active sprint countdown**: 使用 `setInterval(1000)` 客户端倒计时; End Date 从 API 获取 (UTC timestamp); 客户端计算 remaining time; **时钟漂移处理**: 每 60s 与 `GET /api/server-time` 同步校准 (或从 response header `Date` 取服务端时间); 漂移 > 5s 时静默校准

**页面离开 (Destroy)**: 清理 countdown `setInterval` timer

**按钮条件补充**:
- **"Launch Sprint"**: 仅 `status=draft` 或 `status=scheduled` (提前启动) 时可见; Scheduled sprint "Launch Now" 按钮 → 确认 "Start this sprint immediately instead of the scheduled time?" → 触发 D20
- **"End Early"**: 仅 `status=active` 时可见; **reward 正在分发时 disabled** (C-108 Auto 模式分发中, 按钮 disabled + tooltip "Rewards are being distributed...")
- **"View Results" → "Retry"**: 仅管理员角色可见 (非 viewer); Retry 针对单个 Failed tier 重新执行分发
- **Concurrent sprints**: **允许多个同时 Active** 的 Sprint; 用户可参与多个; Stats 汇总所有 Active sprints 数据
- **Sprint 结束 C端行为**: Sprint 结束后, C端 LB Sprint tab 显示 "Ended" 最终排名 (保留 7 天); 7 天后从 C端移除; 仅 Active/刚结束(7天内) 的 Sprint 在 C端可见
- **0 participants at end**: reward 分发跳过 (无需发放); status 变 Completed; View Results 显示空排名表 + "No participants"

**End Early 后 Stats**: Sprint 立即变 Completed; Stats "Active Participants" 归零 (该 sprint); 总计 "Tasks Completed" 保留已有数据

**Auto distribute 中间态**: 分发开始 → 表格行显示 spinner + "Distributing..." 文字 (替代 status badge); 分发完成 → badge 变 "Completed" + 若有 Failed tier → amber warning icon

#### 操作详情 (v2.0 新增, C-106 ~ C-109)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-106 | "End Early" | click (Active Sprint only) | `PUT /api/community/modules/lb-sprint/instances/{id}/end` | 确认 dialog: "End this Sprint now? Winners will be calculated based on current rankings. Rewards will be distributed according to your reward settings." → 确认 → status 变 Completed → 奖励分发 (C-108) | C端 Sprint 页面显示最终排名 | — |
| C-107 | "View Results" (Completed Sprint) | click | — | 弹窗/面板显示: 最终排名表 (rank / user / points / reward) + 奖励发放状态 (Sent ✅ / Pending ⏳ / Failed ❌); Failed 项有 "Retry" 按钮 | — | — |
| C-108 | 奖励发放方式 | auto / manual | `POST /api/community/modules/lb-sprint/instances/{id}/distribute` | 由 D05 Modal 中 M-27 配置: **Auto**: Sprint 结束后自动发放 (服务端 cron, 结束后 5 分钟内); **Manual**: 管理员点击 "Distribute Rewards" 按钮触发. **Token/NFT 预充值流程**: ① D05 保存含 Token/NFT 奖励时, Sprint 详情显示 "Deposit Required" 警告卡片 (amber bg): "Deposit rewards to escrow before launching." + "Deposit Now →" 按钮; ② 点击 → 连接钱包 → 合约 approve + deposit 交互; ③ B31e 表格行内显示充值状态: "Deposit: ✅ Sufficient" (绿色) / "Deposit: ⚠️ Insufficient" (amber, 显示差额); ④ 充值不足时 "Launch Sprint" 按钮 disabled + tooltip "Deposit required rewards before launching"; ⑤ 发放时余额不足 → tier 标记 Failed + "Top Up" 按钮 | — | 发放失败 → 该 tier 标记 Failed + toast + "Top Up" 按钮 |
| C-109 | Sprint 定时开始 | D05 设置 Start Date > today | — | Start Date 设置为未来日期 → status = `Scheduled` (蓝色 badge); 到达 start time 后服务端自动激活 → status 变 Active; 表格中 Scheduled Sprint 显示倒计时 "Starts in 2d 5h" | — | — |

---

### 5.9 B31f — Milestones

**设计稿**: Node `WFdZQ` | URL: `/community/modules/milestone`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Sets | 2 |
| Completions | 847 |
| Rewards Claimed | 623 |
| Claim Rate | 73.6% |

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Milestone Name | 名称 |
| Status | Active / Draft / Not Started |
| Tiers | 阶梯数 |
| Completions | 完成人数 |
| Date | 创建日期 |

#### Primary Action
- "+ Create Milestone" → D06 Milestone Editor (`gtOam`)
- "Activate Milestone" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/milestone/instances` (list + stats + per-milestone tier data)

**页面离开 (Destroy)**: 无需清理

**"Activate Milestone" 按钮位置**: 表格行内 Actions 列, 仅 `status=draft` 时显示 "Activate" 按钮 (触发 D20)

**Milestone threshold 修改 C端影响** (C-112 补充):
- 已达标但未领取用户: **保持 "claimable" 状态** (C端显示绿色 "Claim" 按钮)
- 未达标用户: 根据新 threshold 重新判断
- Tier detail panel (C-110): D06 保存后 **自动刷新** — 展开的 tier 详情即时反映新 threshold/reward

**All tiers claimed by all eligible users**: Milestone 不自动变为 "Completed" — 仍保持 Active (新用户仍可达到阈值); 若管理员希望结束, 手动 Archive (⋮ 菜单)

#### 操作详情 (v2.0 新增, C-110 ~ C-112)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-110 | 展开 Milestone 行 | click row | — | 行下方展开 Tier 详情面板: 每 Tier 一行 (Tier 1: 100pts → Badge "Explorer"; Tier 2: 500pts → Shop Item "Exclusive NFT"; ...); 显示每 Tier 的 completions 和 claim rate | — | — |
| C-111 | Claim Rate 数字点击 | click claim rate % | — | 打开 D18 Segment Detail Panel, 预填 filter: 该 Milestone 的已领取/未领取用户分群; 面板显示 wallet address + claim status + claim date | — | — |
| C-112 | 修改已激活 Milestone threshold | D06 modal 编辑 | — | **允许修改**, 但弹出警告 dialog: "Changing the threshold may affect {n} users who have already reached the current threshold but haven't claimed. They will retain their eligibility." → 确认 → PUT; 已达标但未领取的用户保持资格 | — | — |

---

### 5.10 B31g — Benefits Shop

**设计稿**: Node `7yPWx` | URL: `/community/modules/shop`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Items | 5 |
| Total Redemptions | 1,432 |
| Points Spent | 284,500 |
| Items Sold Out | 2 |

#### Filter Additions
- 标准 filter 外增加 "Price Range" 和 "Search categories..." 下拉

#### Insight Banner
- "Best sellers expire in 3 days. NFT Badge sold to 64%... Replenish before depletion."

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Item Name | 商品名称 |
| Category | 分类 |
| Price | 积分价格 |
| Redemptions/Stock | 兑换数/库存 |
| Status | Active / Sold Out / Draft |

#### Primary Action
- "+ Add Item" → D07 Shop Item Editor (`b1JOT`)
- "Publish Item" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/shop/instances` (list + stats); "Price Range" 和 "Search categories" 筛选: **服务端 API 参数** (`?priceMin=&priceMax=&category=`), 非客户端过滤

**页面离开 (Destroy)**: 无需清理

**按钮条件补充**:
- **"Replenish"**: 可见条件 — `status=sold_out` **或** `stock < 10 且 stock > 0` (low stock 也显示 Replenish); Sold Out: amber 按钮; Low stock: 灰色 outlined 按钮
- **"Publish Item"**: 仅 `status=draft` 行显示 (触发 D20)
- **Replenish 后 Stats**: "Items Sold Out" 计数 -1 (若从 Sold Out 恢复); Insight Banner 如果提示 "X items sold out" → 文案更新 (或 banner 消失 if X=0)
- **Badge gate 引用已删除 badge**: 商品仍可购买 (gate **自动降级为 "All users"**); 表格行显示 amber warning: "Badge gate removed: referenced badge was deleted" + tooltip 详情; 管理员可在 D07 中重新设置 gate

#### 操作详情 (v2.0 新增, C-113 ~ C-116)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-113 | Sold Out "Replenish" | click Replenish (Sold Out 行) | `PUT /api/community/modules/shop/instances/{id}` body `{stock}` | 行内显示 "Replenish" 按钮 (amber) → 点击 → inline number input (min 1, placeholder "Enter quantity") + "Confirm" 按钮 → PUT → stock 更新 → status 恢复 Active → toast "Stock replenished: +{n} items" | C端商品重新可用 | — |
| C-114 | 可用性门控 (Level/Badge-gated) | D07 modal 配置 | — | 配置在 D07 M-38 字段; 表格中 gated 商品显示小图标: 🔒 + gate 类型 (如 "Lv.5+", "Badge: Explorer"); C端: 不满足条件用户看到商品但购买按钮 disabled + 提示 "Reach Level 5 to unlock" | — | — |
| C-115 | 商品图片 | render / click | — | 表格中: 40×40 缩略图 (圆角 6px, object-fit cover); 点击缩略图 → 弹出大图预览 (Modal, max 600×400, 背景 overlay) + "×" 关闭 | — | — |
| C-116 | 库存为 0 自动处理 | stock reaches 0 | — | **自动暂停销售**: status 变 "Sold Out" (红色 badge); C端显示 "Sold Out" 标签; **不发通知给管理员** (通过 Insight Banner 提示: "2 items are sold out. Replenish stock to keep your shop active.") | C端商品显示 Sold Out 状态 | — |

---

### 5.11 B31h — Lucky Wheel

**设计稿**: Node `sme5a` | URL: `/community/modules/wheel`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Spins | 12,847 |
| Winners | 3,241 |
| Prizes Awarded | $8,400 |
| Avg Spins / Day | 428 |

#### Insight Banner
- "Day 6-7 decline matches spin fatigue. Reduce frequency or add surprise jackpot."

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Wheel Name | 名称 |
| Status | Active / Draft |
| Total Spins / Prize Pool | 总转数/奖池 |
| Win Rate / Date | 中奖率/日期 |

#### Primary Action
- "+ Create Wheel" → D08 Lucky Wheel Config (`k2gwC`)
- "Activate Wheel" → 触发 D20

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/wheel/instances` (list + stats)

**页面离开 (Destroy)**: 无需清理

**"Activate Wheel" 条件**: ① `status=draft`; ② 至少 1 个 non-"Nothing" prize (C-120); ③ 概率总和 = 100% (C-117); ④ Token/NFT prizes: 需预充值 (同 B31e C-108 预充值流程, 但规模通常更小)

**After activation Stats**: "Total Spins" 初始为 0; Stats Row 更新 "Active Wheels" +1 (若有此字段)

**Prize modified while wheel active**: 修改通过 D08 编辑 → 保存 → **C端下次 spin 时生效** (非当前 spin session); 用户看到的 wheel 视觉在下次页面加载时更新

**Token/NFT prizes 预充值**: 同 B31e 预充值流程; 差异: wheel 的预充值需考虑总 spin 次数 × 中奖概率 × prize value; 预充值不足时 wheel 自动暂停 (status → paused) + B端 Insight Banner 警告 "Lucky Wheel paused: prize pool depleted"

#### 操作详情 (v2.0 新增, C-117 ~ C-120)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-117 | 奖品概率配置 | D08 modal 内 | — | 每个奖品: name + type + value + 概率 % input (number, 0-100, step 0.1); 底部实时显示: "Total: {sum}%" — 绿色 (=100%) / 红色 (≠100%); ≠100% 时 Save 按钮 disabled + 红色提示 "Prize probabilities must total 100%" | — | — |
| C-118 | 转盘预览 | D08 modal 内 | — | Modal 右侧: 可视化转盘 (CSS/SVG 圆形分割); 每个奖品占比例扇区 + 颜色 + 名称; 修改概率时实时更新扇区大小; 仅视觉预览, 不可交互旋转 | — | — |
| C-119 | Spin Cost 编辑 | D08 modal 内 | — | 在 D08 Modal 内 (M-44); **非内联编辑** (不在表格行内); number input, min 0 (0=免费), 无上限; 显示积分类型符号 (如 "50 XP per spin") | — | — |
| C-120 | 全部 "Nothing" 奖品 | D08 modal save | — | **不允许激活**: 若所有奖品类型均为 "Nothing", Save 允许 (draft), 但 "Activate" 按钮 → 弹出提示 "Add at least one prize (non-Nothing) before activating the wheel." | — | — |

---

### 5.12 B31i — Badges

**设计稿**: Node `BJLsz` | URL: `/community/modules/badges`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Badges | 12 |
| Badges Earned | 847 |
| Unique Holders | 312 |
| Earn Rate | 72.4% |

#### 独特筛选
- Badge Category: Achievement / Engagement / Special

#### Insight Banner
- "12 badges designed — 8 earned by +52% of users. Consider prioritizing less-earned badges."

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Badge Name | 名称 + 图标预览 |
| Category | Achievement / Engagement / Special |
| Condition | 获取条件简述 |
| Earned | 获得人数 |
| Status | Active / Draft / Archived |

#### Primary Action
- "+ Create Badge" → D09 Badge Editor (`YbFvp`)

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/modules/badges/instances` (list + stats)

**页面离开 (Destroy)**: 无需清理

**Auto-trigger badge**: 系统服务端自动评估 (cron 或 event-driven); **无 B端手动检查触发按钮** — 管理员不能手动触发 auto-evaluation; Auto-trigger 条件达成时自动发放, 通常在 task 完成/积分变化/streak 达标等事件后 1-5 分钟内

**Manual badge grant (C-123) 后**: 表格 "Earned" 列 **即时 +1** (乐观更新); **C端通知**: 用户在 C端 User Center 收到 in-app notification (小红点 badge): "You've earned the '{badgeName}' badge!" + badge 图标; **无 push notification** (仅 in-app)

**Auto-trigger badge C端通知**: 同 manual — in-app notification

**Badge 删除级联**:
- Shop (D07 M-38 gate): 商品 gate 自动降级为 "All users" (同 B31g badge gate deleted)
- Milestone (D06 M-31 reward): 对应 tier 的 badge reward 变为 "Badge deleted" 红色文案; 用户仍可达到 threshold 但无 badge 发放; 管理员需编辑 milestone 替换 reward
- TaskChain (D02 M-10 reward): Completion Reward badge 变为 null; 用户完成 chain 不获得 badge (仅获得积分)
- **删除前确认 dialog**: "This badge is referenced by {n} milestones and {m} shop items. Deleting it will remove these references. Continue?"

**Auto-condition "Complete X tasks" — tasks deleted**: 若条件引用的 sector 中所有 tasks 被删除, 条件变为 **不可达成**; B端 badge 行显示 amber warning: "Auto-condition may be unattainable: referenced tasks archived"; 管理员需编辑 badge 修改条件

#### 操作详情 (v2.0 新增, C-121 ~ C-123)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-121 | "Holders" 数字点击 | click number in Earned column | — | 打开 D18 Segment Detail Panel, 预填: 该 Badge 的持有者列表; 面板显示 wallet address + earned date + badge icon | — | — |
| C-122 | Badge 图标 | render / D09 modal | — | 表格中: 32×32 图标预览 (圆角 4px); D09 Modal (M-49): 支持上传图片 (PNG/SVG, max 1MB, 建议 128×128) **或** 选择预设图标 (50+ Material Symbols icons 网格, 每行 8 个, 可搜索) | — | 上传: 格式错误 → "Only PNG or SVG files up to 1MB"; 尺寸不够 → 警告 (非阻止) |
| C-123 | 手动发放 Badge | D09 modal 内 (Manual only 模式) | `POST /api/community/modules/badges/instances/{id}/grant` body `{walletAddress}` | M-54: wallet address 输入框 (0x 验证); "Grant Badge" 按钮 → POST → toast "Badge granted to {address}" → Holders 数 +1 | C端用户即时收到 Badge 通知 | 地址无效 → "Invalid wallet address"; 已持有 → "User already holds this badge" |

---

### 5.13 B49 — Access Rules

**设计稿**: Node `g1CNC` | URL: `/community/settings/access-rules`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Rules | 4 |
| Active Rules | 3 |
| Users Validated | 1,240 |
| Access Denied | 8,912 |

#### Filter Tabs
- All / Active / Paused / Deactivated

#### Data Table 列

| 列名 | 说明 |
|------|------|
| Rule Name | 规则名称 |
| Type | Token Gate / NFT Hold / Level / Invite Only |
| Status | Active / Paused |
| Condition | 条件简述 |
| Date | 创建日期 |

#### Row Actions
- Toggle enable/disable: (API) `PUT /api/community/settings/access-rules/:id`
- Click row → D10 Access Rule Editor (`8HgbJ`)

#### Primary Action
- "+ Create Rule" → D10

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/settings/access-rules` (list + stats)

**页面离开 (Destroy)**: 无需清理

**规则优先级排序后 (C-124)**: "Preview Rule" (C-125) **不自动刷新** — 用户需手动重新 "Test" 以看到新顺序的效果

**多规则交互逻辑**: 规则间为 **AND 关系** — 所有 Active 规则按优先级顺序依次评估, 所有规则都必须通过才允许访问; C-125 Preview 清楚显示每条规则的 pass/fail 状态

#### 操作详情 (v2.0 新增, C-124 ~ C-126)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-124 | 规则优先级排序 | drag handle | `PUT /api/community/settings/access-rules/reorder` body `{order: [id1,...]}` | 拖拽手柄调整规则执行顺序 (从上到下依次评估); 同 §2.3.6 拖拽规格; 排序即时保存 | 规则评估顺序改变影响 C端用户准入判断 | 失败 → 回滚 |
| C-125 | "Preview Rule" | click "Preview" button (表头区) | `POST /api/community/settings/access-rules/preview` body `{walletAddress}` | 输入 wallet address → "Test" 按钮 → 发送 POST → 显示结果面板: 每条规则的 pass/fail 状态 (✅/❌) + 最终结论 "Access Granted" (绿色) 或 "Access Denied: {rule_name}" (红色) | — | 地址无效 → "Invalid address" |
| C-126 | Token Gate 规则依赖集成 | D10 modal 保存时 | — | 选择 Rule Type = "Token Gate" → 检查 B61 Integration Center 是否已配置链上验证 (Multi-Chain 或 On-Chain Verification); **未配置**: Token Gate 选项可选但 Save 后规则 status = "Inactive" + 表格行显示 amber warning icon + tooltip "Requires blockchain integration. Configure in Integration Center →" (链接 → B61). **B61 连接后行为**: Token Gate 规则**不会自动激活** — 管理员需手动在 B49 中将规则 toggle 为 Active; 连接后 amber warning 消失, 规则变为可激活状态 | — | — |
| C-126b | 新规则对已登录 C端用户的影响 | 规则变为 Active 时 | — | **仅影响新访问**: 已在 C端 session 中的用户**不会被踢出**; 规则在用户下次页面刷新或新 session 开始时生效; Access Rules 执行逻辑: 用户每次访问 Community URL 时, 服务端按优先级顺序 (C-124) 依次评估所有 Active 规则, **规则间为 AND 关系** (所有 Active 规则都必须通过); 任一规则 fail → 显示 M-60 配置的 Denial Message + 对应解决方案链接 | C端: 不满足条件的新访客被拒绝进入; 已在线用户下次加载时受影响 | — |

---

### 5.14 B50 — Homepage Editor

**设计稿**: Node `5Wm6B` | URL: `/community/settings/homepage`

#### Stats Row
| 指标 | 示例值 |
|------|--------|
| Total Sections | 6 |
| Visible | 4 |
| Page Views (24h) | 3,847 |
| Avg Session | 2m 34s |

#### Filter Tabs
- All / Visible / Hidden

#### Data Table 列 (拖拽排序)

| 列名 | 说明 |
|------|------|
| Section Name | 分区名称 |
| Type | Banner / Widget / Text / Custom |
| Visibility | Visible ✓ / Hidden |
| Last Modified | 最后修改时间 |

#### 特殊交互
- **拖拽排序**: 行支持拖拽手柄重新排序
- **Visibility Toggle**: 每行可切换显隐
- **Preview Button**: Header 副按钮 → B33 Preview Mode
- Click row → D11 Homepage Section Editor (`rDDZo`)

#### Primary Action
- "+ Add Section" → D11

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/settings/homepage` (sections list + stats)

**页面离开 (Destroy)**: 无需清理

**"即时生效" 与 "Publish" 步骤说明**: Homepage Editor 的排序/显隐 **不需要单独 "Publish" 步骤** — 所有操作 (C-127 排序, C-128 显隐) 通过独立 PUT API 即时生效; C端在下次页面加载时反映变化; 无"批量发布"按钮, 无 Draft 中间态; 这是设计决策 (homepage 配置属于低风险运营操作)

**D11 section 引用已归档 widget**: 若 section 类型为 "Leaderboard Widget" 且引用的 Leaderboard 已 Archived → section 自动 **不可见** (C端不显示); B端表格行显示 amber warning: "Referenced leaderboard archived" + "Fix" 链接 (→ D11 编辑, 重新选择实例)

#### 操作详情 (v2.0 新增, C-127 ~ C-130)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-127 | 拖拽行排序 | drag handle | `PUT /api/community/settings/homepage/reorder` body `{order: [id1,...]}` | 同 §2.3.6; 排序**即时生效** (无单独发布步骤); 保存成功 → 无 toast (静默); C端首页下次加载时反映新排序 | C端首页 section 顺序更新 | 失败 → 回滚 + toast |
| C-128 | Visibility toggle | click toggle | `PUT /api/community/settings/homepage/sections/{id}` body `{visible}` | 乐观更新: toggle ON → section 可见; toggle OFF → section 隐藏; **即时生效** (C端下次加载反映); 无 toast | C端对应 section 出现/消失 | 失败 → 回滚 |
| C-129 | "Preview" 按钮 | click (header 副按钮) | — | → B33 Preview Mode, **新标签打开** (target=_blank); 不离开当前 Homepage Editor 页面 | — | — |
| C-130 | Section 类型及配置 | D11 modal 内 | — | 允许的类型: **Banner** (image upload + click URL + alt text) / **Quest Widget** (选择 sector 或 "All Tasks") / **Leaderboard Widget** (选择 Leaderboard 实例) / **Points Widget** (积分余额 + 等级进度) / **Text** (rich text: bold/italic/link/heading) / **Custom HTML** (code editor with basic syntax highlighting; 自动 XSS 过滤: 移除 script/on* attributes) | — | — |

---

## 6. 运营辅助页

### 6.1 B32 — Content Management

**设计稿**: Node `lhR14` | URL: `/community/content`

#### 页面概述
- **功能**: 管理 C端首页的公告、推荐位、模块状态

#### 页面结构

```
Content Area
├── Header: "Content Management" + "Preview C-End" + "+ New Announcement"
├── ANNOUNCEMENTS Section
│   ├── "Announcement Carousel" (Active, 3 items)
│   │   ├── Item 1: "Token Launch — Double Points!" [Pinned]
│   │   ├── Item 2: "Community AMA with Founders"
│   │   └── Item 3: "Token Launch Party" [Scheduled]
├── FEATURED SLOTS Section
│   ├── "Featured Grid (2×3)" (Active)
│   │   ├── Slot 1: "Token Launch Quest" ✓
│   │   ├── Slot 2: "Weekly Sprint" ✓
│   │   ├── Slot 3: "Milestone Rewards" ✓
│   │   ├── Slot 4: "Referral Program" ✓
│   │   ├── Slot 5: "+ Add Featured"
│   │   └── Slot 6: "+ Add Featured"
├── MODULE STATUS OVERVIEW Section
│   └── 6 cards: Points System / Leaderboard / TaskChain / DayChain / Benefits Shop / Lucky Wheel
│       └── Each: status indicator (Active/Not Configured) + "Configure" link
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Preview C-End" | → B33 |
| "+ New Announcement" | → D16 Announcement Editor (`6TLjE`) |
| Announcement actions | (API) pin/unpin, edit→D16, delete |
| "+ Add Featured" | → D17 Featured Slot Editor (`DVVpL`) |
| Module "Configure" | → B31a-B31h |
| "Publish" (content) | → 触发 D20 |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- **并行**调用:
  - `GET /api/community/content/announcements` — 公告列表
  - `GET /api/community/content/featured` — Featured slots
  - `GET /api/community/modules` — Module Status Overview
- 加载态: 每个 section 独立 skeleton

**页面离开 (Destroy)**: 无需清理

**Announcement 编辑/删除 UI**: 每条 announcement 行右侧 ⋮ 菜单 (3 项: Edit → 打开 D16 预填 / Delete → 确认 dialog / Pin/Unpin toggle); **不使用内联按钮** (节省空间)

**Pin 后 carousel 排序**: Pinned 公告始终排在 carousel 列表**第一位**; 同时只能有 1 个 pinned (M-74); 新 pin → 旧 pin 自动 unpin; carousel 列表排序: pinned → 按发布时间倒序

**Scheduled announcement 自动发布**: 服务端 cron 到达 schedule time → status 变 Active; B端 **不主动通知管理员** (不弹 toast/push); 下次刷新 B32 页面时自动看到 status 变化; **无实时 WebSocket 推送** (低优先级)

**Featured slot 已满 (6/6)**: "+ Add Featured" 格子不再显示; 标题区副文: "All 6 featured slots are filled. Remove a slot to add a new one."

**Featured slot 引用已结束的 Sprint/Quest**: 自动显示 "Ended" 状态标签 (灰色); **不自动移除** — 管理员手动替换; C端显示 "Ended" 标签但仍可点击查看历史

**Module Status Overview — module 引用已删除资源**: 模块卡片正常显示 status; 资源级别的 broken reference 不影响 module status 显示; module status 仅反映 Active/Not Configured

#### 操作详情 (v2.0 新增, C-131 ~ C-136)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-131 | Announcement Pin/Unpin | click pin icon | `PUT /api/community/content/announcements/{id}` body `{pinned}` | Pin icon toggle: 📌 (pinned, amber fill) ↔ 📌 (unpinned, gray outline); Pinned 公告始终排在 carousel 列表顶部; 乐观更新 | C端 carousel 顶部固定该公告 | — |
| C-132 | Announcement "Scheduled" | render / D16 配置 | — | 显示发布时间: "Scheduled for Mar 15, 2026 10:00 UTC" (灰色副文); 时间到达后服务端定时任务自动发布 → status 变 Active → carousel 中出现; B端列表中 Scheduled 状态: 蓝色 badge + 时间显示 | — | — |
| C-133 | "+ Add Featured" | click empty slot | — | 弹出 D17 Featured Slot Editor: 选择内容类型 (Quest / LB Sprint / Milestone / External URL) → 选择具体实例 (从已发布列表) → Save → 格子填充缩略图 + 标题 | — | — |
| C-134 | Featured Slot 已填充 | render + hover | — | 显示: 缩略图 + 标题; hover: 显示 overlay 操作: "×" 移除 (确认: "Remove from featured?" → DELETE → 格子恢复 "+ Add Featured") + "✏️" 编辑 (→ D17 预填数据). **移除后 C端影响**: **即时生效** (C端下次 API 请求时格子消失); Featured Grid 自动收缩 (非留空格子) | C端 Featured Grid 即时更新 | — |
| C-135 | Module Status "Configure" | click | — | 每个模块卡片的 "Configure" 链接: Points→B31a, Leaderboard→B31d, TaskChain→B31b, DayChain→B31c, Benefits Shop→B31g, Lucky Wheel→B31h; 已配置: 绿色指示灯 + "Manage"; 未配置: 灰色指示灯 + "Configure" | — | — |
| C-136 | Announcement "Publish" | click (Draft announcement) | — | 单条公告发布 **不走 D20** (公告属于内容运营, 非产品功能发布); 直接 `PUT {status:'active'}` → 乐观更新 badge → toast "Announcement published" | C端 carousel 出现该公告 | — |

---

### 6.2 B33 — Preview Mode

**设计稿**: Node `2UiNC` | URL: `/community/preview`

#### 页面概述
- **功能**: 嵌入 C端完整预览，模拟用户视角

#### 页面结构

```
Full Page
├── Preview Banner (non-dismissible, amber)
│   ├── "⚠️ Preview Mode — This is how participants see your community."
│   ├── Desktop / Mobile toggle
│   └── [Exit Preview] button
├── Embedded C-End Frame
│   ├── C-End Header (dark #0F172A)
│   ├── Nav Tabs: Home / Quests / Leaderboard / Sprint / Milestones / Shop
│   ├── C-End Content (mock data)
│   └── Footer: "Powered by TaskOn"
```

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 数据来源: **B端 Preview API** (`GET /api/community/preview`) — 返回 C端页面数据的 preview 版本 (使用真实 B端配置 + mock 用户数据); 非 C端 API (避免产生真实数据)
- URL 参数: `?from=B32` / `?from=B50` / `?from=B55` — 用于 Exit Preview 返回; `?tab=home` / `?tab=quests` — 可选, 预选 C端 tab
- 加载态: 全页 skeleton (C端 header + content 占位)

**页面离开 (Destroy)**: 若使用 iframe 嵌入 C端 → 清理 iframe src; 无表单需保存

**Exit Preview 路径修复**: 使用 **URL 参数 `?from=`** 判断来源, **不使用 `document.referrer`** (隐私策略可能导致 referrer 为空); 无 `?from` 参数时默认返回 B32 Content Management

**Mobile preview (375px)**: 使用 **CSS transform scale** 缩放 (非真实 375px viewport); 比例 = 375/实际容器宽度; **不模拟触摸交互** (仅视觉预览); 点击交互同 desktop (鼠标事件, 非 touch event)

**Preview 无启用模块**: C端显示 Home tab (仅 Welcome 信息 + 品牌色) + 空 Quests (§2.3.3 空状态); 无其他 tab; Preview 仍可正常渲染

#### Mock Data 策略

| 数据类型 | 数据来源 |
|---------|---------|
| 用户身份 | Mock: "Preview User", Lv.5, 1,250 pts |
| 任务状态 | 混合: completed / available / locked |
| 排行榜 | Mock: 10 条目 |
| 积分余额 | Mock: 1,250 |
| 连续签到 | Mock: 7 天 |
| 商品列表 | 真实 B端配置 |
| 公告 | 真实 B端配置 |

#### 交互逻辑

1. **Desktop/Mobile Toggle**: 切换预览视口宽度 (desktop=100%, mobile=375px centered)
2. **Tab 导航**: 可在预览内切换 Home/Quests/Leaderboard 等
3. **Exit Preview**: → B32 Content Management

#### 操作详情 (v2.0 新增, C-137 ~ C-140)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-137 | Mock 用户身份 | render | — | 固定 mock: 用户名 "Preview User", Level 5, 积分 1,250, 7-day streak; 不可修改; 头像: 默认灰色占位头像 | — | — |
| C-138 | Preview 内链接/按钮拦截 | click any interactive element | — | **拦截所有点击**: 不执行真实 API 操作; 仅模拟 UI 状态变化 (如点击 "Claim" → 按钮变 "Claimed ✓" 但不实际发放); 链接点击: 不导航, 显示 toast "Links are disabled in preview mode" | — | — |
| C-139 | "Exit Preview" | click button | — | 返回**来源页面**: 从 B32 进入 → 返回 B32; 从 B50 进入 → 返回 B50; 从 B55 进入 → 返回 B55; 通过 `document.referrer` 或 URL 参数 `?from=B32` 判断 | — | — |
| C-140 | Preview Banner | render | — | **non-dismissible**: 无 × 关闭按钮; 始终固定在页面顶部; amber 背景 `#1F1A08` border `#D97706`; 高度 48px; 内容: ⚠️ icon + "Preview Mode" + toggle + "Exit Preview" 按钮 | — | — |

---

### 6.3 B54 — Community Insights

**设计稿**: Node `olPfE` | URL: `/community/insights`

#### 页面概述
- **功能**: 跨模块分析 + 经济健康度 + 用户分群 + 留存

#### 页面结构

```
Content Area
├── Header: "Community Insights" + Date Picker + Module Filter + "Export CSV" + "Export PDF"
├── Key Metrics Bar (4 stats, highlighted)
│   ├── D30 Retention Rate: 42.7% (green)
│   ├── Daily Active Users: 1,847
│   ├── Points Burn Rate: 24.5k/day
│   └── Economy Balance: "Healthy"
├── Points Economy Chart — Earn vs Burn Trend
│   └── Dual-axis bar chart (6 months), with Earned/Burned/Net lines
├── Bottom Row (2 panels)
│   ├── User Segments (pie chart + table)
│   │   ├── Power Users: 127 (14%)
│   │   ├── Active: 402 (44%)
│   │   ├── At Risk: 193 (21%)
│   │   └── Dormant: 1,221 (41%)
│   └── Retention by Module (horizontal bars)
│       ├── Tasks: 1,470 (48.4%)
│       ├── Points: 68 (24.2%)
│       ├── Leaderboard: 1,314 (17.1%)
│       └── DayChain: 577 (5.8%)
```

#### API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/community/insights/overview` | GET | 60s |
| `/api/community/insights/economy` | GET | 300s |
| `/api/community/insights/segments` | GET | 300s |
| `/api/community/insights/retention` | GET | 300s |
| `/api/community/insights/export` | GET | N/A |

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- **全部并行**调用 5 个 API (不同 cache 时长):
  - `/api/community/insights/overview` (60s) — Key Metrics + WAU chart
  - `/api/community/insights/economy` (300s) — Points Economy Chart
  - `/api/community/insights/segments` (300s) — User Segments pie chart + table
  - `/api/community/insights/retention` (300s) — Retention by Module bars
  - (optional) `/api/community/modules` (60s) — Module Filter dropdown 数据
- **每个 section 独立 skeleton + 独立加载**: overview skeleton = 4-stat bar + chart 占位; economy skeleton = dual-axis chart 占位; segments = pie chart + 4 行表格; retention = 4 horizontal bars
- 任一 API 失败: 对应 section 显示 §2.3.4 错误态; 其他 section 正常

**页面离开 (Destroy)**: 无需清理

**Export 按钮**: "Export CSV" / "Export PDF" — **0 数据时 disabled** (opacity 0.5 + tooltip "No data to export for selected date range")

**Module Filter (C-142) 补充**: 选中模块无数据时, 图表区域显示 "No data for selected modules in this date range" (非空图表, 也非隐藏图表); Filter 选择**不同步 URL** (仅 Date Picker 同步 URL `?from=&to=`); Filter 选择变化 → **客户端过滤** (前端从已加载数据中过滤, 非 re-fetch API)

**Date 变化 (C-141)**: 所有 API **重新请求** (忽略 cache, 因为 date range 变了 → cache key 不同); 切换 date range = 新的 cache key

**Retention 条形图点击 (C-147)**: 跳转对应模块管理页 **不传日期参数** (模块管理页无 date range 概念)

**已禁用模块在 Retention chart**: 若模块在统计期间内曾 Active 但现在 Inactive → **仍显示在 chart 中** (附带 "(disabled)" 灰色标签); 保留历史数据完整性

#### 操作详情 (v2.0 新增, C-141 ~ C-147)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-141 | Date Picker | click | — | 预设快捷选项: "Last 7 days" / "Last 30 days" / "Last 90 days"; 自定义范围: calendar 双日期选择 (点击 start → 点击 end, 蓝色高亮选中范围); 选中后所有 chart 和 stats 重新加载 (带 `?from=&to=` 参数) | — | — |
| C-142 | Module Filter | click | — | 多选 dropdown: "All Modules" (默认) + 各模块名 (Tasks / Points / Leaderboard / DayChain / TaskChain / Badges / LB Sprint / Milestones / Shop / Lucky Wheel); checkbox 多选; 选中后图表数据过滤; 显示选中数: "3 modules selected" | — | — |
| C-143 | Economy Chart hover | hover bar | — | Tooltip 格式: "{Month}\nEarned: {earned} {pointType}\nBurned: {burned} {pointType}\nNet: +{net}" (绿色为正, 红色为负); 三行数据 + 竖线对齐 | — | — |
| C-144 | User Segment 卡片点击 | click segment slice / row | — | 打开 D18 Segment Detail Panel: 预填该分群 (Power / Active / At Risk / Dormant); 面板显示用户列表 + 搜索 + 导出 | — | — |
| C-145 | "Export CSV" | click | `GET /api/community/insights/export?format=csv` | **异步导出**: 点击 → 按钮 spinner + "Generating..." → 后台生成 → 完成: 自动触发浏览器下载 + toast "Report downloaded"; 大数据量时 (>10k rows): 显示进度百分比 | — | 生成失败 → toast "Export failed, please try again" |
| C-146 | "Export PDF" | click | `GET /api/community/insights/export?format=pdf` | 同 CSV, PDF 格式含: 图表截图 (echarts/chart.js 的 toDataURL) + 数据表格 + 日期范围标注 | — | 同上 |
| C-147 | Retention by Module 条形图 | hover / click | — | Hover: tooltip 显示精确数值 "{module}: {count} users ({percent}% retained)"; 点击条形 → 跳转对应模块管理页 (Tasks→B31, Points→B31a, etc.) | — | — |

---

### 6.4 B61 — Community Integration Center

**设计稿**: Node `ZL5K5` | URL: `/community/integrations`

#### 页面概述
- **功能**: Community 专属集成管理（不含 WL 专属的 SDK/Iframe/PB）
- **入口**: B12 Deep Hub → "All Integrations →"

#### 页面结构

```
Content Area
├── Breadcrumb: "Back to Community"
├── Header: "Integration Center"
├── Status Bar: "2 of 9 integrations active"
├── ── Social & Community ──
│   ├── Twitter / X [Connected] — "Configure"
│   ├── Discord [Connected] — "Configure"
│   └── Telegram [Available] — "Connect"
├── ── Blockchain & Wallet ──
│   ├── Multi-Chain [Available] — "Connect"
│   ├── Wallet Connect [Available] — "Connect"
│   └── On-Chain Verification [Available] — "Connect"
├── ── Analytics & Data ──
│   ├── Google Analytics [Available] — "Connect"
│   ├── Webhooks [Available] — "Connect"
│   └── Data Export [Available] — "Connect"
```

#### Integration Card 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 集成 ID |
| name | string | 名称 |
| category | enum | `social` / `blockchain` / `analytics` |
| status | enum | `connected` / `available` / `error` |
| description | string | 简述 |
| icon | string | 图标标识 |
| configUrl | string | 配置页 URL |

#### 交互逻辑

- **Connected**: 绿色边框 + "Configure" 按钮（绿色）
- **Available**: 默认边框 + "Connect" 按钮（紫色）
- **All buttons**: 点击 → 对应集成配置页/流程

#### 页面生命周期 (v2.0 补充)

**页面初始化 (Init)**:
- 调用: `GET /api/community/integrations` — 返回 9 个集成的状态列表 (connected/available/error)
- 加载态: 9 个 card skeleton (3 per section)
- **URL anchor 支持** (C-33 sends `#{integration_id}`): 页面加载完成后 → `document.getElementById(integration_id).scrollIntoView({behavior:'smooth'})` + 目标卡片边框短暂高亮 (绿色 glow 500ms)

**页面离开 (Destroy)**: 若有 OAuth popup 仍在打开 → 关闭 popup window (`popup.close()`)

**Breadcrumb "Back to Community"**: → 按 §10.1 路由到当前 Hub 状态 (B10/B11/B12)

**"Disconnect" 按钮**: 在 C-152 Configure 展开面板内; **始终可见** (无需二次确认弹窗前的隐藏); 点击后显示确认 dialog; 确认后 → `DELETE /api/community/integrations/{id}` → 卡片变回 "Available" 状态

**集成与订阅层级**: **无层级限制** — 所有集成对所有付费用户可用; 免费试用用户也可配置 (但 Publish 被 D20 拦截)

**After connect/disconnect**: Status Bar "2 of 9 integrations active" **即时更新** (前端乐观更新计数)

**Blockchain 集成 → B49 关系**: 连接 Blockchain 后 → B49 Token Gate / NFT Hold 规则类型变为**可用** (不自动创建规则, 不自动激活, 见 C-126)

**Twitter disconnect → Publish 影响**: 断开 Twitter 后 → 所有 D20 的 Twitter 检查项变为 ❌; 但可通过 "Publish anyway" (M-92) 跳过 Twitter 检查 (仅当订阅已通过时); 断开不影响已发布内容 (已发布的 tasks/modules 继续运行)

#### 操作详情 (v2.0 新增, C-148 ~ C-153)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 级联效果 | 错误处理 |
|---|------|---------|---------|------------|---------|---------|
| C-148 | Twitter "Connect" | click | OAuth 2.0 | 流程: 点击 → popup window (600×700) → Twitter OAuth 授权页 → 用户 "Authorize" → 回调 URL `/api/integrations/twitter/callback` → popup 关闭 → 主页面状态变 "Connected" (绿色 border + ✓) + toast "Twitter connected". **D20 联动**: 连接成功后, D20 Readiness Check 的 Twitter 项将通过; D20 每次打开时**重新检测** (非缓存回调) — 调用 `GET /api/integrations/status` 检查 Twitter 连接是否有效; 5 分钟缓存 (M-93) | D20 Readiness Check 的 Twitter 项通过 | popup 被阻止 → toast "Please allow popups"; 授权失败 → toast "Twitter authorization failed" |
| C-149 | Discord "Connect" | click | OAuth 2.0 + Bot invite | 流程: 点击 → popup → Discord OAuth → 授权后跳转 Bot 邀请页 → 用户选择 server → 确认 → 回调 → Connected; 需要: 读取消息/管理角色权限 | — | 同上 |
| C-150 | Telegram "Connect" | click | — | 流程: 展开 inline 配置面板: 输入 Bot Token (从 @BotFather 获取) + "Verify" 按钮 → `POST /api/integrations/telegram/verify` body `{token}` → 验证 bot 存在且可用 → Connected | — | Token 无效 → "Invalid bot token. Make sure you copied the full token from @BotFather" |
| C-151 | Blockchain "Connect" | click | — | 展开配置面板: 选择链 (dropdown: Ethereum / BSC / Polygon / Base / Arbitrum / Optimism / Avalanche) → 使用默认 RPC (推荐, TaskOn 提供) 或 自定义 RPC endpoint (URL input + "Test Connection" 按钮) → `POST /api/integrations/blockchain` → Connected | 启用 Token Gate 功能 (B49 Access Rules) | RPC unreachable → "Unable to connect to RPC endpoint. Check the URL and try again." |
| C-152 | "Configure" 已连接集成 | click | — | 展开 inline 配置面板: 显示当前配置摘要 (如 "@ProjectName connected on Mar 1") + 修改选项 (如切换 Twitter 账号) + "Disconnect" 按钮 (红色 text, 确认 dialog: "Disconnect {integration}? Features depending on this integration may stop working.") | — | — |
| C-153 | 集成 error 状态 | render | — | 错误卡片: 红色边框 `#DC2626` 1px + `error` icon (红色) + 错误原因文案 (如 "Token expired" / "RPC unreachable" / "Bot removed from server") + "Reconnect" 按钮 (重新发起 OAuth/验证流程); **仅手动 "Reconnect"**, 无自动重连 | — | — |
| C-154 | Wallet Connect "Connect" | click | — | 展开配置面板: ① 选择 Wallet Connect 版本 (v2 推荐); ② 输入 WalletConnect Project ID (从 cloud.walletconnect.com 获取); ③ "Test Connection" → `POST /api/integrations/wallet-connect/verify` → 成功: Connected. **用途**: 用户通过钱包登录 C端 community | — | Project ID 无效 → "Invalid Project ID. Get one from cloud.walletconnect.com" |
| C-155 | On-Chain Verification "Connect" | click | — | 展开配置面板: ① 选择链 (复用 C-151 的链列表); ② 选择验证类型: Token Balance / NFT Ownership / Contract Interaction; ③ 输入合约地址 (0x 验证); ④ "Test" → `POST /api/integrations/onchain-verify/test` → 成功: Connected. **用途**: 用于 B49 Access Rules 的 Token Gate / NFT Hold 规则验证 | 启用 B49 Token Gate / NFT Hold 规则类型 | 合约地址无效 → "Contract not found on selected chain" |
| C-156 | Google Analytics "Connect" | click | — | 展开配置面板: ① 输入 GA4 Measurement ID (格式 `G-XXXXXXXXXX`); ② 可选: 输入 API Secret (用于服务端事件); ③ "Verify" → `POST /api/integrations/ga/verify` → 验证格式有效 → Connected. **注意**: 仅格式验证, 不实际验证 GA 账户所有权; 实际追踪需前端嵌入 gtag.js | — | 格式错误 → "Invalid Measurement ID format. Expected G-XXXXXXXXXX" |
| C-157 | Webhooks "Connect" | click | — | 展开配置面板: ① 输入 Webhook URL (HTTPS required); ② 选择触发事件 (多选 checkbox): User Joined / Task Completed / Level Up / Points Earned / Reward Claimed; ③ 可选: 添加 Secret Header (用于验签); ④ "Send Test Event" → `POST /api/integrations/webhooks/test` → 目标 URL 返回 2xx → Connected + 显示 "Test event sent successfully"; **Connected 后**: 卡片显示已配置的事件列表 + 最近 5 次 delivery 状态 (✅ 200 / ❌ 500 / ⏳ Timeout) | — | URL unreachable → "Webhook endpoint did not respond. Check URL and try again." |
| C-158 | Data Export "Connect" | click | — | 展开配置面板: ① 选择导出目标: Manual Download (默认, 无需配置) / S3 Bucket / Google Cloud Storage; ② S3/GCS: 输入 bucket name + access key + secret; ③ "Test Upload" → `POST /api/integrations/data-export/test` → 上传小测试文件 → 成功: Connected. **Manual Download**: 始终可用 (无需 "Connect"), 仅 S3/GCS 需要配置; 卡片对 Manual Download 显示 "Available" 而非 "Connect" | — | 权限不足 → "Cannot write to bucket. Check IAM permissions." |

---

## 7. Modal 弹窗

### 7.1 Community 相关 Modal 清单

| 编码 | 名称 | Node ID | 触发页面 |
|------|------|---------|---------|
| D01 | Points & Level Editor | `8DhXJ` | B31a |
| D02 | TaskChain Editor | `bZiB5` | B31b |
| D03 | DayChain Config | `nEAUB` | B31c |
| D04 | Leaderboard Config | `j8UnD` | B31d |
| D05 | LB Sprint Editor | `NnzO9` | B31e |
| D06 | Milestone Editor | `gtOam` | B31f |
| D07 | Shop Item Editor | `b1JOT` | B31g |
| D08 | Lucky Wheel Config | `k2gwC` | B31h |
| D09 | Badge Editor | `YbFvp` | B31i |
| D10 | Access Rule Editor | `8HgbJ` | B49 |
| D11 | Homepage Section Editor | `rDDZo` | B50 |
| D16 | Announcement Editor | `6TLjE` | B32 |
| D17 | Featured Slot Editor | `DVVpL` | B32 |
| D18 | Segment Detail Panel | `4FPLn` | B54 |
| D19 | Promo Kit Generator | `2qNbJ` | B10/B15 |
| D20 | Publish Readiness Check | `fY99y` | ALL Publish buttons |
| D21 | Task Editor (v2.0 新增) | — | B31 |

### 7.2 D20 — Publish Readiness Check (通用)

**设计稿**: Node `fY99y`

#### 触发时机
所有 "Publish" / "Go Live" / "Launch" / "Activate" 按钮点击时。

#### Community 触发页面清单

| 页面 | 触发按钮 |
|------|---------|
| B55 Wizard Step 4 | "Publish Community" |
| B31 Sectors & Tasks | "Publish Task" |
| B31b TaskChain | "Activate Chain" |
| B31c DayChain | "Activate Chain" |
| B31e LB Sprint | "Launch Sprint" |
| B31f Milestones | "Activate Milestone" |
| B31g Benefits Shop | "Publish Item" |
| B31h Lucky Wheel | "Activate Wheel" |
| B32 Content Mgmt | ~~"Publish" (公告/内容)~~ — **已移除**: 公告发布不走 D20 (见 C-136 说明) |

**D20 触发范围说明 (v2.0 勘误)**: B32 Content Management 页面的公告 (Announcement) 发布**不触发 D20** — 公告属于内容运营行为, 非产品功能发布。C-136 为准确行为定义。仅"模块级别的 Publish/Activate/Launch"操作触发 D20, 包括: Task 发布 (B31)、Chain/Wheel/Sprint/Milestone 激活 (B31b-B31h)、Shop Item 发布 (B31g)、Community 首次发布 (B55)。

#### Checklist 项

| # | 检测项 | 通过条件 | 失败处理 |
|---|--------|---------|---------|
| 1 | 订阅状态 | 试用期内 OR 付费有效 | "Upgrade Plan →" 链接 |
| 2 | 官方 Twitter 授权 | 已连接并验证 | "Connect Twitter →" 链接 |

#### 交互逻辑

1. **全部通过**: 自动关闭 Modal → 执行发布 API
2. **部分未通过**: Publish 按钮 disabled，显示未通过项 + 解决链接
3. **检测中**: 显示 loading spinner

#### D20 补充规格 (v2.0 新增, M-91 ~ M-93)

| # | 操作 | 触发方式 | 说明 |
|---|------|---------|------|
| M-91 | 检查超时 | 10s timer | 10s 内未完成检测 → 显示 "Verification timed out." + "Publish anyway?" 链接 (仅订阅已通过时显示) |
| M-92 | "Publish anyway" | click | 跳过 readiness check 强制发布; **仅当订阅检查已通过时允许** (Twitter 未通过可跳过); 点击后直接执行发布 API |
| M-93 | 结果缓存 | auto | 5 分钟内再次触发同一 Publish 按钮 → 跳过重新检查, 直接显示上次结果; 缓存 key: `{entityType}_{entityId}_readiness`; 5 分钟后过期重新检测. **存储位置**: `sessionStorage` (非 localStorage) — tab 级别隔离; **两个 tab 行为**: 各自独立缓存, 不互相影响; 刷新 tab → cache 清空 (sessionStorage 行为) |

**D20 补充规格 (v2.0)**:
- **检测 API**: ① 订阅检查: `GET /api/subscription/status` → `{active: boolean, plan: string, expiresAt: ISO8601}`; ② Twitter 检查: `GET /api/integrations/status?type=twitter` → `{connected: boolean, username: string}`; 两个 API **并行调用**; **单项超时**: 每个 API 独立 5s 超时 (M-91 的 10s 是总超时)
- **M-92 "Publish anyway" 精确条件**: ① 订阅检查通过 (active=true) + Twitter 检查失败 (connected=false 或超时) → "Publish anyway" 可点击; ② 订阅检查失败 → "Publish anyway" **不可用** (必须解决订阅问题); ③ 订阅检查超时 + Twitter 通过 → "Publish anyway" 可点击 (给予通过的信任); ④ 两项都超时 → "Publish anyway" 不可用, 仅显示 "Retry checks" 按钮

---

### 7.3 Modal 完整字段规格 (v2.0 新增)

以下为 Community 所有 16 个 Modal 的完整字段定义。

**通用 Modal 规格**: 居中, 宽度 600px (可调), 圆角 12px, 背景 `#111B27`, border `#1E293B`, 背景遮罩 `#000` 50%, Esc 关闭 (form dirty 时需确认), 标题 18px bold `#F1F5F9`, 底部 action bar: Cancel (left) + Save/Create (right, green `#48BB78`).

#### 7.3.1 D01 — Points & Level Editor

**触发**: B31a "+ Add Level" / 行点击编辑 | **Node**: `8DhXJ` | **宽度**: 520px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-01 | Level Name | text input | ✅ | 1-30 字符, 同 community 内唯一 (不区分大小写) | 如 "Bronze" / "Silver"; 重名: "A level with this name already exists" |
| M-02 | Threshold | number input | ✅ | ≥0 整数; 必须 > 前一等级 threshold 且 < 后一等级 (见 C-93); 第一个等级 threshold 可为 0 | 达到该积分值即升级 |
| M-03 | Point Type | select dropdown | ✅ | 从已配置类型选 | 数据源: `GET /api/community/point-types`; 默认选第一个 (通常 XP) |
| M-04 | Level Badge | icon picker | — | 可选 | 方式: 选择预设图标 (20+ 奖杯/星星/盾牌图标) 或 上传图片 (PNG/SVG, max 512KB, 64×64 建议) |
| M-05 | Level Perks | textarea | — | max 200 字符 | 等级特权描述; 显示在 C端等级详情; placeholder "Describe perks for this level..." |
| M-06 | 保存行为 | — | — | — | 新建: `POST /api/community/modules/points/levels`; 编辑: `PUT .../levels/{id}`; 成功: 关闭 modal + 刷新列表 + toast "Level saved"; 失败: modal 不关闭 + 字段旁显示 API 错误 |

**D01 补充规格 (v2.0)**:
- **Modal 打开时 API**: `GET /api/community/point-types` (M-03 dropdown 数据); 若 0 个 point types → M-03 显示 "No point types. Create one first." + inline "+ Add" link (触发 C-91); modal 内显示 loading spinner 直到 API 返回
- **Esc 关闭 + dirty form**: 适用 §7.3 通用规格 — form dirty 时 Esc → 确认 dialog "Discard changes?"; form clean → 直接关闭
- **"Save"/"Create" disabled 条件**: threshold 值违反顺序 (C-93) 时 Save disabled + M-02 红色 border + 错误提示
- **Create vs Edit 模式区别**: **标题**: Create → "Create Level"; Edit → "Edit Level: {levelName}". **预填**: Edit 模式所有字段预填当前值 (name, threshold, point type, badge, perks). **Delete 按钮**: 仅 Edit 模式显示 (modal 左下角红色文字链接 "Delete Level")
- **新 Level C端影响**: 新建的 Level **即时生效** — 符合 threshold 的用户在下次积分变化或服务端 batch 重算时自动升级 (5 分钟内); 无需额外 "Publish" 步骤

#### 7.3.2 D02 — TaskChain Editor

**触发**: B31b "+ Create Chain" / 行点击编辑 | **Node**: `bZiB5` | **宽度**: 680px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-07 | Chain Name | text input | ✅ | 1-50 字符 | |
| M-08 | Steps | dynamic list (sortable) | ✅ | 至少 2 步, 最多 20 步 | 每步: 序号 + 任务选择 + 拖拽手柄; "+ Add Step" 按钮在列表底部 |
| M-09 | Step 任务选择 | multi-select dropdown | ✅ per step | 每步至少 1 个任务; **同一任务不可跨步骤重复选择** — dropdown 自动过滤已被其他步骤选中的任务 (灰色 disabled + tooltip "Already used in Step {N}") | 数据源: `GET /api/community/tasks?status=active` (B31 已发布任务); 显示 task name + sector; 搜索过滤; **0 条 active tasks 时**: modal 仍可打开, 但任务选择框显示 "No active tasks. Create tasks in Sectors & Tasks first." + 链接 → B31; Save 按钮 disabled |
| M-10 | Completion Reward | 积分数 + badge select | — | 积分 ≥ 0; badge 从 B31i 选 | 全部步骤完成后的额外奖励; 不填则无额外奖励 (各步骤自身积分仍获得) |
| M-11 | 步骤顺序 | drag handle | — | — | 拖拽调整步骤顺序; 同 §2.3.6 规格 |

**保存**: 同 D01 模式; Chain 保存为 Draft; 需手动 "Activate" (→ D20).

**D02 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/tasks?status=active` (M-09 dropdown); **0 active tasks**: modal 可打开, 但显示 info banner "No active tasks available. Create and publish tasks in Sectors & Tasks first." + "Go to Sectors & Tasks →" 链接; M-09 dropdown 为空; Save 按钮 disabled
- **"+ Add Step" disabled**: 达到 20 步时 disabled + tooltip "Maximum 20 steps per chain"
- **"Save" disabled**: < 2 步时 disabled + tooltip "A chain requires at least 2 steps"
- **保存后 B31b refresh**: 关闭 D02 → B31b 表格自动刷新 (重新 GET instances API); 若新建 → 新 draft 出现在列表顶部

#### 7.3.3 D03 — DayChain Config

**触发**: B31c "+ Create Chain" / 行点击编辑 | **Node**: `nEAUB` | **宽度**: 560px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-12 | Chain Name | text input | ✅ | 1-50 字符 | |
| M-13 | Daily Task | task select dropdown | ✅ | 从已发布任务选 | 选择每日需完成的任务 (用户每天完成此任务 = 连续签到) |
| M-14 | Base Reward | number input | ✅ | ≥1 整数 | 每日完成基础积分; 显示积分类型符号 |
| M-15 | Milestone Bonuses | dynamic list: day + multiplier | — | day: ≥2 整数, 不重复; multiplier: ≥1.1, 最大 10x | 预填: Day 7=2x, Day 14=3x, Day 30=5x; 可增删; "+ Add Milestone" |
| M-16 | Grace Period | number input (slider) | — | 0-24, 单位 hours; 默认 0 | 断链容忍时长; 0=无容忍 (严格每日); slider + 数字显示 |

**D03 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/tasks?status=active` (M-13 daily task dropdown); 同 D02
- **编辑 Active DayChain 对用户 streak 影响**: 修改 daily task → **不影响已有 streak 计数**; 用户当天如果已完成旧任务, 当天算已完成; 次日开始按新 task 判定; 修改 milestone bonus / grace period → 即时生效 (见 B31c 页面生命周期)
- **M-15 Milestone Bonus 排序**: 添加 Day 5 (between Day 3 and Day 7) → **自动按 day 值升序排列** (保存时后端排序, 同 D06 M-30); UI 中拖拽排序不可用 (顺序由 day 值决定); 重复 day 值: 保存时报错 "Day {N} milestone already exists"

#### 7.3.4 D04 — Leaderboard Config

**触发**: B31d "+ Create Leaderboard" / 行点击编辑 | **Node**: `j8UnD` | **宽度**: 520px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-17 | Name | text input | ✅ | 1-50 字符 | |
| M-18 | Point Type | select | ✅ | 从已配置类型选 | 基于哪种积分排名 |
| M-19 | Period | radio group | ✅ | Weekly / Monthly / All Time | 重置规则见 M-21 |
| M-20 | Display Top N | number input | — | 默认 100, range 10-1000 | C端展示排行榜条目数 |
| M-21 | 重置时间 | display (readonly) | — | — | 自动根据 Period: Weekly = 周一 00:00 UTC; Monthly = 1日 00:00 UTC; All Time = 不重置 |

**D04 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/point-types` (M-18 dropdown)
- **编辑 Active Leaderboard Period 变更** (如 Weekly → Monthly): 确认 dialog "Changing the period will reset current rankings. Historical data from the previous period will be preserved. Continue?" → 确认 → rankings 归零重新计算; C端即时显示重置后的排名

#### 7.3.5 D05 — LB Sprint Editor

**触发**: B31e "+ Create LB Sprint" / 行点击编辑 | **Node**: `NnzO9` | **宽度**: 640px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-22 | Sprint Name | text input | ✅ | 1-50 字符 | |
| M-23 | Point Type | select | ✅ | 从已配置类型选 | 基于哪种积分排名 |
| M-24 | Start Date | date picker | ✅ | ≥ today (不可选过去日期) | calendar 组件, 时间精确到小时 |
| M-25 | End Date | date picker | ✅ | > Start Date; 最长 90 天 | End Date - Start Date ≤ 90d; 违反: "Sprint duration cannot exceed 90 days" |
| M-26 | Reward Tiers | dynamic list | ✅ | 至少 1 tier | 每 tier: Rank Range (from-to, 如 1-3) + Reward Type (select: Token/NFT/WL Spot/Badge/Points) + Quantity/Value; "+ Add Tier" 按钮; rank ranges 不可重叠 |
| M-27 | Reward 发放方式 | radio | ✅ | Auto / Manual | Auto: 结束后自动发放 (需预充值); Manual: 需管理员点击 "Distribute" (见 C-108) |

**D05 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/point-types` (M-23 dropdown)
- **M-26 Rank overlap 错误**: Tier rank ranges 输入后实时验证; 重叠示例: Tier 1=1-3, Tier 2=2-5 → Tier 2 rank "from" 字段红色 border + 错误文案 "Rank range overlaps with Tier 1 (1-3)"; **Save disabled** 当存在重叠
- **最小 Sprint 持续时长**: **1 小时** (End Date - Start Date ≥ 1h); 违反: "Sprint must be at least 1 hour long"
- **编辑 Scheduled Sprint**: **未开始前所有字段可编辑**, 包括 Start Date、End Date、Reward Tiers; 已开始 (Active) 的 Sprint: 仅 End Date 可提前 (End Early) + Reward Tiers 不可修改 (已锁定)

#### 7.3.6 D06 — Milestone Editor

**触发**: B31f "+ Create Milestone" / 行点击编辑 | **Node**: `gtOam` | **宽度**: 600px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-28 | Milestone Name | text input | ✅ | 1-50 字符 | |
| M-29 | Tiers | dynamic list | ✅ | 至少 1 tier, 最多 10 | 每 tier: Threshold (积分值) + Reward (类型+值); "+ Add Tier" |
| M-30 | Tiers 顺序 | auto | — | — | 按积分阈值升序自动排序 (不可手动排序); 保存时后端排序 |
| M-31 | 奖励类型 per tier | select | ✅ | Badge / Shop Item / Custom text | Badge: 从 B31i 已有 badge 选 (dropdown); Shop Item: 从 B31g 已有商品选; Custom: 文本输入 (max 100 字符) |

**D06 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/modules/badges/instances` (M-31 Badge dropdown) + `GET /api/community/modules/shop/instances` (M-31 Shop Item dropdown)
- **"+ Add Tier" max**: 达到 10 tiers 时 disabled + tooltip "Maximum 10 tiers per milestone"
- **Adding tier to active milestone**: **允许** — 新 tier 添加后, 已超过新 threshold 的用户 **retroactively 获得资格** (可 claim); 用户无需重新达标; 确认 dialog: "Adding this tier will make {n} existing users eligible to claim immediately."
- **M-31 Badge dropdown 为空 (B31i 无 badges)**: dropdown 显示 "No badges created yet" + "Create Badge →" 链接 (→ B31i, 新标签); 选择 Badge 类型时 Save 按钮 disabled 直到选中有效 badge

#### 7.3.7 D07 — Shop Item Editor

**触发**: B31g "+ Add Item" / 行点击编辑 | **Node**: `b1JOT` | **宽度**: 600px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-32 | Item Name | text input | ✅ | 1-60 字符 | |
| M-33 | Description | textarea | — | max 300 字符 | placeholder "Describe what the user gets..." |
| M-34 | Category | select | ✅ | NFT / Token / Merchandise / Experience / Other | 影响 C端分类筛选 |
| M-35 | Image | file upload | — | PNG/JPG/SVG, max 2MB, 建议 400×400 | 拖拽上传 + 点击上传; 预览 thumbnail 120×120 |
| M-36 | Price (Points) | number input | ✅ | ≥1 整数 | 显示积分类型符号 (如 "100 XP") |
| M-37 | Stock | radio + number | ✅ | "Unlimited" radio / "Limited" radio + number input (≥1) | Limited 选中时 number input 出现 |
| M-38 | Availability Gate | select | — | All (默认) / Level (min level N) / Badge (must hold badge X) | Level: 展开 level number input; Badge: 展开 badge select (从 B31i) |
| M-39 | Status | radio | ✅ | Save as Draft / Publish Now | Publish Now → 触发 D20 |

**D07 补充规格 (v2.0)**:
- **Modal 打开 API**: `GET /api/community/modules/points/levels` (M-38 Level gate dropdown) + `GET /api/community/modules/badges/instances` (M-38 Badge gate dropdown)
- **M-39 "Publish Now" 流程**: 选择 "Publish Now" + 点击 "Save" → ① 先保存 item (POST/PUT) → ② 保存成功后自动触发 D20 → ③ D20 通过 → `PUT {status:'active'}` → 关闭 D07 + toast "Item published"; **两步操作** (save + publish), 非一步
- **M-35 Image drag-drop 状态**: **Hover**: 拖拽区 border 变为 `#48BB78` 虚线 + 背景色 `#0A1F1A`; **Invalid file**: border 变红 `#DC2626` + "Only PNG/JPG/SVG files up to 2MB" 提示 (1s 后恢复); **Uploading**: 进度条 (绿色 `#48BB78`, 100% 后消失) + upload icon 旋转; **Success**: 预览 thumbnail 出现 + "×" 删除按钮
- **Active item 价格编辑 C端影响**: 价格变更 **即时生效**; C端下次加载时显示新价格; **无 "cart" 概念** — Community Shop 是积分兑换, 用户点击 "Redeem" 即时扣积分, 无购物车/结算流程

#### 7.3.8 D08 — Lucky Wheel Config

**触发**: B31h "+ Create Wheel" / 行点击编辑 | **Node**: `k2gwC` | **宽度**: 720px (含右侧预览)

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-40 | Wheel Name | text input | ✅ | 1-50 字符 | |
| M-41 | Prizes | dynamic list | ✅ | ≥2 prizes; 总概率 = 100% | 每 prize 一行: Name + Type + Value + Probability%; "+ Add Prize" 按钮 |
| M-42 | Prize 类型 | select per prize | ✅ | Points / NFT / Token / Nothing | Nothing = 未中奖安慰 |
| M-43 | 概率验证 | computed display | — | — | 底部实时: "Total: {sum}%" — 绿色 (=100%) / 红色 (≠100%); ≠100% 时 Save disabled (见 C-117) |
| M-44 | Spin Cost | number input | ✅ | ≥0 整数 (0=免费) | 每次抽奖消耗积分 |
| M-45 | Spin Limit | radio group | ✅ | Once per user / Daily (reset 00:00 UTC) / Unlimited | |
| M-46 | Duration | date range picker | — | start ≤ end; 不填则永久有效 | 可选; 空 = 永久 |

**右侧**: 可视化转盘预览 (见 C-118).

**D08 补充规格 (v2.0)**:
- **转盘预览加载**: modal 打开时转盘预览区显示 skeleton (圆形占位 200×200); prizes 数据加载完成后立即渲染转盘; 编辑 prize 时预览实时更新 (无 loading delay)
- **M-41 "All Nothing" 激活限制自动解除**: 用户修改任一 prize 从 "Nothing" 变为其他类型 → 保存后 (draft) → "Activate" 按钮**自动变为可用** (限制解除); 无需额外操作; 判断逻辑在 Activate 按钮点击时实时检查 (非 cache)

#### 7.3.9 D09 — Badge Editor

**触发**: B31i "+ Create Badge" / 行点击编辑 | **Node**: `YbFvp` | **宽度**: 560px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-47 | Badge Name | text input | ✅ | 1-30 字符, 唯一 | 重名: "A badge with this name already exists" |
| M-48 | Description | textarea | — | max 200 字符 | |
| M-49 | Icon | file upload / preset picker | ✅ | Upload: PNG/SVG, max 1MB, 建议 128×128; Preset: 50+ icons | Tab 切换: "Upload" / "Preset"; Preset: 8×N grid, searchable |
| M-50 | Category | select | ✅ | Achievement / Engagement / Special | |
| M-51 | Earn Condition | radio | ✅ | Auto-trigger / Manual only | Auto: 系统自动发放; Manual: 管理员手动 (见 C-123) |
| M-52 | Auto Condition | select + params | ✅ (if Auto) | — | 选项: "Complete {N} tasks" (number input) / "Reach Level {N}" (level select) / "Maintain {N}-day streak" (number) / "Earn {N} points" (number + point type) |
| M-53 | Is Rare | toggle | — | 默认 off | Rare badge: C端有特殊视觉 (金色 glow + sparkle 动画) |
| M-54 | 手动发放 | wallet input + button | — (Manual only) | 0x + 40 hex chars | 见 C-123 |

**D09 补充规格 (v2.0)**:
- **50+ Preset icons**: **前端打包资源** (非 API 请求) — Material Symbols icons 子集, bundled in app; 分类: Achievement (trophies, stars, medals) / Social (groups, chat) / Activity (trending, lightning) / Custom; grid 内搜索: 按 icon name 过滤 (client-side)
- **M-51 Auto→Manual 切换**: 切换时 **清空 M-52 auto-condition 值** (无确认弹窗); M-52 section **隐藏** (slide-up); 反向 Manual→Auto: M-52 section 出现 (slide-down), 需重新配置条件
- **M-54 "Grant Badge" 可用模式**: **仅 Edit 模式** 显示 (badge 必须已创建); Create 模式下 M-54 section 隐藏; Manual only 模式的 badge 在 Create 后进入 Edit 才能手动发放
- **M-52 条件详解**: "Complete {N} tasks" — N: 正整数 ≥1; 范围: **所有 sectors 的所有 active tasks** (非特定 sector); "Reach Level {N}" — dropdown 选择已配置的 level (非数字输入); "Maintain {N}-day streak" — N: 正整数 ≥1; 任一 DayChain 达标即可; "Earn {N} points" — N: 正整数 ≥1 + Point Type select (从已配置类型选)

#### 7.3.10 D10 — Access Rule Editor

**触发**: B49 "+ Create Rule" / 行点击编辑 | **Node**: `8HgbJ` | **宽度**: 560px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-55 | Rule Name | text input | ✅ | 1-50 字符 | |
| M-56 | Rule Type | select | ✅ | Token Gate / NFT Hold / Level Requirement / Invite Only | 选择后下方显示对应参数面板 |
| M-57 | Token Gate Params | 3 fields | ✅ (if Token Gate) | Contract address (0x+40hex) + Minimum balance (≥0) + Chain (从 B61 已配置链选) | 依赖 B61 链上集成 (见 C-126) |
| M-58 | Level Req Params | number input | ✅ (if Level) | min level ≥1 | "Minimum level to access" |
| M-59 | Invite Only Params | code or file | ✅ (if Invite) | 邀请码 (auto-generate 或 custom 6-12 chars) 或 whitelist CSV upload (wallet addresses, max 10k rows) | |
| M-60 | Denial Message | textarea | — | max 200 字符 | 不满足规则时 C端显示; placeholder "You need [condition] to access this community." |

**D10 补充规格 (v2.0)**:
- **Modal 打开 API**: M-57 Token Gate 需 `GET /api/community/integrations` (获取已配置的链列表); **0 chains (B61 未配置 blockchain)**: Chain dropdown 显示 "No blockchain integrations configured" + "Set Up in Integration Center →" 链接 (→ B61); Token Gate 选项仍可选, 但 Save 后 rule status = Inactive (见 C-126)
- **M-56 Rule Type 切换**: 切换类型时 **清空前一类型的参数** (无确认弹窗); 如 Token Gate → Level Requirement: contract address/chain/balance 清空, level input 出现
- **M-59 Invite Only CSV 验证**: **逐行验证** (非聚合); 上传后显示验证结果表: ✅ Valid / ❌ Invalid per row; 错误类型: "Invalid address format (row 15)" / "Duplicate address (row 23, same as row 7)"; 允许部分导入: "Import {valid} valid addresses, skip {invalid} invalid?" (确认 dialog); **最大 10,000 行** — 超出: "File too large. Maximum 10,000 addresses per upload."

#### 7.3.11 D11 — Homepage Section Editor

**触发**: B50 "+ Add Section" / 行点击编辑 | **Node**: `rDDZo` | **宽度**: 640px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-61 | Section Type | select | ✅ | Banner / Quest Widget / Leaderboard Widget / Points Widget / Text / Custom HTML | 选择后下方内容面板变化 |
| M-62 | Title | text input | ✅ | 1-60 字符 | Section 标题, 显示在 C端 |
| M-63 | Banner 内容 | image upload + URL input | ✅ (if Banner) | Image: PNG/JPG, max 2MB, 建议 1200×300; Link URL: valid URL | 点击 banner 跳转链接 |
| M-64 | Widget 内容 | widget select | ✅ (if Widget type) | 选择具体模块实例 | Quest Widget: 选择 sector; Leaderboard Widget: 选择 leaderboard instance; Points Widget: 无需选择 (自动显示) |
| M-65 | Text 内容 | rich text editor | ✅ (if Text) | max 2000 字符 | 工具栏: Bold / Italic / Link / Heading / List |
| M-66 | Custom HTML 内容 | code editor (textarea) | ✅ (if Custom HTML) | max 5000 字符; XSS 过滤 | 自动过滤: `<script>`, `on*` attributes, `javascript:` URLs; 保存前显示 "HTML will be sanitized for security" |
| M-67 | Visibility | select | ✅ | All users / Logged-in only / Level-gated (展开 min level input) | |

**D11 补充规格 (v2.0)**:
- **Modal 打开 API**: M-64 Widget 类型需数据: Quest Widget → `GET /api/community/sectors` (sector 列表); Leaderboard Widget → `GET /api/community/modules/leaderboard/instances` (leaderboard 列表); Points Widget → 无需 (自动)
- **M-66 Custom HTML "sanitized" 警告**: **保存时显示** (非输入时) — 点击 Save → 执行 sanitize → 若有内容被移除 → 显示警告 dialog: "The following elements were removed for security: {list of removed tags/attributes}. Save sanitized version?" → 确认 → 保存; 无被移除内容 → 直接保存
- **Custom HTML C端生效时机**: 保存后 **即时生效** (同 B50 C-128); 无需额外 B50 visibility ON (D11 保存的 section 默认 visibility = ON, 除非用户在 B50 手动关闭)
- **M-65 Rich Text 暗色主题渲染**: 编辑器工具栏: 暗色背景 `#1A2332`, 按钮 `#94A3B8`; 编辑区: 白色文字 `#F1F5F9` on 暗色背景 `#0D1521`; **Link 颜色**: `#5D7EF1` (蓝色, 非默认浏览器蓝); Heading: `#F1F5F9` bold; C端渲染: 沿用 C端主题色 (暗色 header + 浅色 content, 确保 readability)

#### 7.3.12 D16 — Announcement Editor

**触发**: B32 "+ New Announcement" / Announcement actions edit | **Node**: `6TLjE` | **宽度**: 560px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-68 | Title | text input | ✅ | 1-80 字符 | 字符计数: `{n}/80` |
| M-69 | Content | textarea | ✅ | 1-500 字符 | plain text + 支持 1 个 URL (自动链接化); 字符计数 |
| M-70 | Type | select | ✅ | General / Event / Alert | 影响 C端显示: General=蓝色图标, Event=绿色, Alert=红色 |
| M-71 | Image | file upload | — | PNG/JPG, max 2MB | 拖拽/点击上传; 预览 thumbnail |
| M-72 | CTA Button | text + URL (2 inputs) | — | Button text max 20 chars; URL valid format | 如 "Learn More" + `https://...`; 不填则无按钮 |
| M-73 | Schedule | radio | ✅ | "Publish Now" / "Schedule" | Schedule: 展开 date+time picker (≥ now + 5min) |
| M-74 | Pin to Top | toggle | — | 默认 off | Pinned 公告始终置顶; 同时只能有 1 个 pinned (新 pin → 旧 pin 自动 unpin) |

**D16 补充规格 (v2.0)**:
- **Edit 已发布公告**: 所有字段可编辑 (Title, Content, Type, Image, CTA, Schedule, Pin); 编辑已发布公告 → 保存后**即时生效** C端 (无需重新 Publish)
- **M-73 Schedule 时间流逝**: 若 modal 打开时设置了 Schedule 时间, 但 modal 保持打开超过该时间 → 保存时**前端验证**: "Scheduled time has passed. Publish now or select a new time." → 两个按钮: "Publish Now" / "Reschedule"
- **M-74 Pin 警告**: 启用 Pin toggle 时, 若已有另一条 pinned 公告 → inline 提示 (不阻止, 不弹 dialog): "This will replace the currently pinned announcement: '{currentPinnedTitle}'" (amber 文案, 在 toggle 下方)
- **Scheduled announcement 发布后 B端通知**: **无主动通知** (不弹 B端 toast/push); 管理员下次刷新 B32 时看到 status 变化

#### 7.3.13 D17 — Featured Slot Editor

**触发**: B32 "+ Add Featured" / Featured slot edit | **Node**: `DVVpL` | **宽度**: 480px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-75 | Content Type | select | ✅ | Quest / LB Sprint / Milestone / External URL | 选择后内容面板变化 |
| M-76 | Content Select | select dropdown | ✅ (if ≠ External) | — | 从已发布的对应内容中选; 数据源: GET API; 显示 name + status |
| M-77 | External URL | URL input | ✅ (if External) | valid URL format | |
| M-78 | Custom Title | text input | — | max 60 字符 | 覆盖原标题; 空 = 使用原标题 |
| M-79 | Custom Image | file upload | — | PNG/JPG, max 2MB | 覆盖原缩略图; 空 = 使用原图 |

**D17 补充规格 (v2.0)**:
- **Modal 打开 API**: M-76 Content Select 需 API: Quest → `GET /api/community/tasks?status=active`; LB Sprint → `GET /api/community/modules/lb-sprint/instances?status=active,completed`; Milestone → `GET /api/community/modules/milestone/instances?status=active`; 显示 name + status
- **Featured slot → sprint 结束**: Sprint 结束后 featured slot **不自动移除** — 显示 "Ended" 灰色 badge; C端: slot 仍可见, 显示 "Ended" 标签; 管理员手动替换或移除

#### 7.3.14 D18 — Segment Detail Panel

**触发**: B54 User Segment 点击 / B31f Claim Rate 点击 / B31i Holders 点击 | **Node**: `4FPLn` | **宽度**: 720px (side panel, right slide-in)

**注意**: D18 是只读面板, 无编辑功能。

| # | 字段 | 类型 | 说明 |
|---|------|------|------|
| M-80 | Segment 类型 | display badge | Power (绿色) / Active (蓝色) / At Risk (amber) / Dormant (灰色); 或 Badge/Milestone 名称 |
| M-81 | User 列表 | data table | 列: Wallet Address (truncated 0x1234...abcd) / Last Active (relative time) / Total Points / Modules Used (icon badges) |
| M-82 | Search | text input | 按 wallet address 搜索; debounce 300ms |
| M-83 | Filter | dropdown | 按模块使用情况筛选: All / Tasks only / Points only / ... |
| M-84 | Export CSV | button | 导出当前面板数据; 同 C-145 异步模式; 文件名: `segment_{type}_{date}.csv` |

**Panel 交互**: 右侧 slide-in (宽度 720px, 高度 100vh); 背景 overlay; 点击 overlay 关闭; 顶部 × 关闭按钮.

**D18 补充规格 (v2.0)**:
- **Panel 打开 API**: `GET /api/community/insights/segments/{type}/users?page=1&limit=50` — 返回分页用户列表; **分页规格**: 每页 50 条; 底部 "Load More" 按钮 (非传统分页器, 因为 side panel 空间有限); 或 infinite scroll (scroll 到底部自动加载下一页)
- **Loading skeleton**: 用户表格: 5 行 skeleton (头像圆 + 地址文字占位 + 3 列占位); Panel header (segment badge + count) 从传入参数直接显示, 不等 API
- **M-84 Export CSV disabled**: 0 用户时 disabled + tooltip "No users to export"; 正在导出时 disabled + spinner
- **M-82 Search / M-83 Filter**: **服务端** — 搜索: `?search={query}` (debounce 300ms, 按 wallet address 匹配); Filter: `?module={type}` (服务端过滤); 结果替换当前列表 (重置 pagination)

#### 7.3.15 D19 — Promo Kit Generator

**触发**: B10 "Generate Promo Kit" / B15 同 | **Node**: `2qNbJ` | **宽度**: 680px

| # | 字段 | 类型 | 说明 |
|---|------|------|------|
| M-85 | Target Platform | radio group | Twitter / Discord / Telegram; 切换后文案模板和字数限制变化 |
| M-86 | AI 生成文案 | display + editable textarea | 初始: AI 生成 (`POST /api/promo-kit/generate` body `{platform, communityName}`); 可二次编辑; 字数限制: Twitter 280 / Discord 2000 / Telegram 4096; 实时字符计数. **加载态**: textarea 显示 3 行 skeleton 占位 + "Generating copy..." 提示 (预计 3-8s); **失败态**: 显示 fallback 模板文案 (硬编码通用模板: "Join our community on TaskOn! {url}") + amber 提示 "AI generation failed. You can edit the template above or try again." + "Retry" 按钮 |
| M-87 | Generated Banner | image display | AI 生成品牌化图片 (社区名 + brand color + 数据摘要); "Regenerate" 按钮 → 重新生成 (POST API); **加载态**: 图片区域显示 skeleton (高度 200px) + spinner + "Generating banner..." (预计 5-15s); **失败态**: 显示默认 TaskOn 品牌 banner (静态资源) + "Failed to generate custom banner. Retry?" 链接; **Regenerate 限制**: 每次打开 modal 最多 regenerate 5 次 (防止 API 滥用); 达到限制: "Regenerate" 按钮 disabled + tooltip "Generation limit reached. Close and reopen to try again." |
| M-88 | Copy Text | button | clipboard API → toast "Text copied!" |
| M-89 | Download Banner | button | 触发 PNG 下载 (filename: `{communityName}_promo_{platform}.png`) |
| M-90 | Share on [Platform] | button | 打开对应分享链接: Twitter intent / Discord (复制) / Telegram share URL; 按钮文案: "Share on Twitter" 等 |

**D19 补充规格 (v2.0)**:
- **M-87 "Regenerate" 覆盖用户编辑**: 若用户已编辑 M-86 textarea 内容, 点击 "Regenerate" → 确认 dialog: "Regenerate will replace your edits with new AI-generated text. Continue?" → 确认 → 重新 POST API → 新文案替换 textarea; 取消 → 保留用户编辑
- **M-85 Platform 切换**: 切换 target platform (如 Twitter → Discord) → **自动重新生成** 文案 (POST API with new platform); 字数限制同步更新 (280 → 2000); **若用户已编辑**: 确认 dialog "Switching platform will regenerate text. Your edits will be lost. Continue?"; Banner 图片**不重新生成** (仅文案变化)
- **加载和失败态**: 已在 M-86/M-87 字段说明中详细规格化 (见上方修改)

#### 7.3.17 D21 — Task Editor (v2.0 新增)

**触发**: B31 "+ New Task" (C-88) / Task 行编辑 ✏️ (C-85) | **宽度**: 600px

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| M-94 | Task Name | text input | ✅ | 1-80 字符; 同 sector 内唯一 | 重名: "A task with this name already exists in this sector" |
| M-95 | Description | textarea | — | max 500 字符 | placeholder "Describe what the user needs to do..." |
| M-96 | Task Type | select | ✅ | `social` / `onchain` / `verification` / `custom` / `recurring` / `referral` | 选择后下方显示类型特定参数 (如 social → Twitter URL 输入; onchain → 合约地址+链选择) |
| M-97 | Sector | select dropdown | ✅ | 从已有 sectors 选 | 数据源: `GET /api/community/sectors`; 新建 modal 默认选中当前展开的 sector |
| M-98 | Points Amount | number input | ✅ | ≥1 整数, max 10000 | 显示积分类型符号 (如 "50 XP") |
| M-99 | Point Type | select | ✅ | 从已配置类型选 | 数据源: 同 M-03; 默认选第一个 |
| M-100 | Verification Method | select | ✅ | Auto (API 验证) / Manual (管理员审核) / Self-report (用户自行标记) | Auto: 系统自动验证完成 (如 Twitter follow API check); Manual: 提交后管理员审批; Self-report: 用户点击即完成 |
| M-101 | Deadline | date picker | — | ≥ today (可选) | 空 = 永不过期; 过期后 status 自动变 `expired` |
| M-102 | Max Claims | number input | — | ≥1 整数 (可选) | 空 = 无限次; 达到上限后 status 自动变 `completed` |
| M-103 | Is Required | toggle | — | 默认 off | Required tasks 在 C端显示 `Required` 标签; 用户必须完成才能参与部分活动 |
| M-104 | Status | radio | ✅ | Save as Draft / Publish Now | Publish Now → 触发 D20 |

**保存行为**: 新建: `POST /api/community/tasks` body `{...fields, sectorId}`; 编辑: `PUT /api/community/tasks/{id}`; 成功: 关闭 modal + B31 页面刷新 (sector 列表重新加载) + toast "Task saved"; 编辑模式: 所有字段预填当前值, 标题显示 "Edit Task" (新建: "Create Task").

---

## 8. 侧栏架构

### 8.1 全局侧栏（未进入模块时）

```
TaskOn (Logo, 20px bold)
────────────────
Home              → /
─ PRODUCTS ─
Quest             → /quest
Community ★       → /community  (active: green bg #0A1F1A)
White Label       → /white-label
Boost             → /boost
────────────────
Analytics         → /analytics
Settings          → /settings
```

### 8.2 Community 展开侧栏（进入模块管理时）

见 §5.2 详细结构。展开动画：chevron `keyboard_arrow_up` 旋转 + 子项 slide-down。

---

## 9. API 接口汇总

### 9.1 Community 核心 API

| Endpoint | Method | 用途 | 页面 | Cache |
|----------|--------|------|------|-------|
| `/api/community/stats` | GET | Hub 统计数据 | B10-B12 | 60s |
| `/api/community/list` | GET | 社区列表 | B10-B12 | 30s |
| `/api/community/drafts` | POST | 保存草稿 | B13,B34,B35,B55 | N/A |
| `/api/community/publish` | POST | 发布社区 | B55 | N/A |
| `/api/community/sectors` | GET/PUT | 分区管理 | B31 | 30s |
| `/api/community/content` | GET | 内容管理 | B32 | 30s |
| `/api/community/modules/:type` | GET/PUT | 模块配置 | B31a-B31i | 60s |
| `/api/community/modules/:type/instances` | GET/POST | 模块实例 CRUD | B31a-B31i | 30s |
| `/api/community/insights/overview` | GET | 分析概览 | B54 | 60s |
| `/api/community/insights/economy` | GET | 经济分析 | B54 | 300s |
| `/api/community/insights/segments` | GET | 用户分群 | B54 | 300s |
| `/api/community/insights/retention` | GET | 留存分析 | B54 | 300s |
| `/api/community/insights/export` | GET | 导出报告 | B54 | N/A |
| `/api/community/integrations` | GET | 集成列表 | B61 | 60s |
| `/api/community/settings/access-rules` | GET/PUT | 访问规则 | B49 | 60s |
| `/api/community/settings/homepage` | GET/PUT | 首页编辑器 | B50 | 60s |
| `/api/promo-kit/generate` | POST | Promo Kit 生成 | B10 | N/A |

#### v2.0 新增 API

| Endpoint | Method | 用途 | 页面 | Cache |
|----------|--------|------|------|-------|
| `/api/community/status` | GET | 社区存在状态 (路由守卫) | B09 | 30s |
| `/api/community/onboarding/progress` | GET | Checklist 进度 | B10 | 30s |
| `/api/community/onboarding/dismiss` | PUT | 关闭 Checklist Banner | B11 | N/A |
| `/api/community/insights/ai` | GET | AI 洞察列表 | B12 | 300s |
| `/api/community/insights/ai/:id/dismiss` | PUT | 关闭单条洞察 | B12 | N/A |
| `/api/community/drafts` | GET | 获取草稿 (恢复) | B13 | N/A |
| `/api/community/readiness` | GET | 发布准备检查 | B55 | N/A |
| `/api/community/check-slug` | GET | URL slug 唯一性检查 | B55 | N/A |
| `/api/community/point-types` | GET/POST/PUT | 积分类型管理 | B31a | 60s |
| `/api/community/sectors/reorder` | PUT | Sector 排序 | B31 | N/A |
| `/api/community/tasks/reorder` | PUT | Task 排序 | B31 | N/A |
| `/api/community/tasks/:id` | PUT/DELETE | Task CRUD | B31 | N/A |
| `/api/community/modules/:type/instances/:id/duplicate` | POST | 复制实例 | ALL modules | N/A |
| `/api/community/modules/lb-sprint/instances/:id/end` | PUT | 提前结束 Sprint | B31e | N/A |
| `/api/community/modules/lb-sprint/instances/:id/distribute` | POST | 分发奖励 | B31e | N/A |
| `/api/community/modules/badges/instances/:id/grant` | POST | 手动发放 Badge | B31i | N/A |
| `/api/community/settings/access-rules/reorder` | PUT | 规则排序 | B49 | N/A |
| `/api/community/settings/access-rules/preview` | POST | 模拟规则判断 | B49 | N/A |
| `/api/community/settings/homepage/reorder` | PUT | 首页 section 排序 | B50 | N/A |
| `/api/integrations/twitter/callback` | GET | Twitter OAuth 回调 | B61 | N/A |
| `/api/integrations/telegram/verify` | POST | Telegram Bot 验证 | B61 | N/A |
| `/api/integrations/blockchain` | POST | 链上集成配置 | B61 | N/A |
| `/api/community/tasks` | POST | 新建 Task | B31 (D21) | N/A |
| `/api/community/preview` | GET | C端预览数据 | B33 | 30s |
| `/api/subscription/status` | GET | 订阅状态 | D20 | 30s |
| `/api/integrations/status` | GET | 集成状态查询 | D20/B61 | 30s |
| `/api/integrations/wallet-connect/verify` | POST | WalletConnect 验证 | B61 | N/A |
| `/api/integrations/onchain-verify/test` | POST | On-Chain 验证测试 | B61 | N/A |
| `/api/integrations/ga/verify` | POST | GA4 验证 | B61 | N/A |
| `/api/integrations/webhooks/test` | POST | Webhook 测试 | B61 | N/A |
| `/api/integrations/data-export/test` | POST | 数据导出测试 | B61 | N/A |
| `/api/community/insights/segments/{type}/users` | GET | 分群用户列表 | D18 | 60s |
| `/api/server-time` | GET | 服务端时间 (时钟同步) | B31e | N/A |

### 9.2 WebSocket 端点

| Endpoint | 用途 | 页面 |
|----------|------|------|
| `/ws/community/participants` | 首批用户实时计数 | B10 Checklist Step 5 |
| `/ws/community/stats` | 实时统计更新 | B11/B12 |

---

## 10. 状态路由策略

### 10.1 Hub 状态判断

```typescript
function getCommunityHubState(community: Community): PageCode {
  if (!community) return 'B09'; // Empty
  if (!community.onboardingComplete) return 'B10'; // Guided
  if (community.activeModules <= 3) return 'B11'; // Active
  return 'B12'; // Deep
}
```

### 10.2 Entity Lifecycle

```
Community: (none) → Draft → Active → Paused → Deleted
Module:    Not Configured → Active → Inactive → Not Configured
Task:      Draft → Active → Completed/Expired/Inactive → Deleted
Sector:    Active → Hidden → Deleted
```

### 10.3 C-End Tab Visibility Rules

C端 Tab 由 B端模块激活状态动态控制：

| C-End Tab | 需要 B-End 模块 | 显示条件 |
|-----------|---------------|---------|
| Home | (always) | 始终显示 |
| Quests | Active tasks | 至少 1 个已发布任务 |
| Leaderboard | Leaderboard module | 已启用且至少 1 条记录 |
| LB Sprint | LB Sprint module | 有活跃或刚结束的 Sprint |
| Milestone | Milestones module | 已启用且至少 1 个阈值 |
| Shop | Benefits Shop | 已启用且至少 1 个商品 |

---

## 附录 A: 设计稿 Node ID 索引

| 页面 | Node ID |
|------|---------|
| B09 Community Hub (Empty) | `zzZ8D` |
| B10 Community Hub (Guided) | `S1EIA` |
| B11 Community Hub (Active) | `vFRHi` |
| B12 Community Hub (Deep) | `TQR51` |
| B13 Community Wizard Step 1 | `Gzpeu` |
| B34 Community Wizard Step 2 | `8NeyG` |
| B35 Community Wizard Step 3 | `qknQZ` |
| B55 Community Wizard Step 4 | `7mVsZ` |
| B31 Sectors & Tasks | `Wug7d` |
| B31a Points & Level | `zCfKQ` |
| B31b TaskChain | `lpdtp` |
| B31c DayChain | `fLLVb` |
| B31d Leaderboard | `Emmab` |
| B31e LB Sprint | `FO9JR` |
| B31f Milestones | `WFdZQ` |
| B31g Benefits Shop | `7yPWx` |
| B31h Lucky Wheel | `sme5a` |
| B31i Badges | `BJLsz` |
| B49 Access Rules | `g1CNC` |
| B50 Homepage Editor | `5Wm6B` |
| B32 Content Management | `lhR14` |
| B33 Preview Mode | `2UiNC` |
| B54 Community Insights | `olPfE` |
| B61 Integration Center | `ZL5K5` |

---

## 附录 B: 跨模块级联效果 (v2.0 新增)

模块之间存在依赖和联动关系。以下为完整的级联效果矩阵。

### B.1 模块启用/禁用级联

| 操作 | 直接影响 | C端影响 | 注意事项 |
|------|---------|---------|---------|
| 启用 Sectors & Tasks | — (Required, 始终启用) | Quests tab 出现 | — |
| 启用 Points & Level | — (Required, 始终启用) | 用户状态栏显示积分+等级 | — |
| 启用 TaskChain | B31b 管理页可用; 需要 Sectors & Tasks 中有 active tasks | C端任务页出现 TaskChain 入口 | 若无 active tasks, TaskChain 无法配置步骤 |
| 启用 DayChain | B31c 管理页可用; 需要至少 1 个 active task 作为 daily task | C端出现连续签到组件 | 同上 |
| 启用 Leaderboard | B31d 管理页可用; 需要至少 1 种积分类型 (B31a) | C端 Leaderboard tab 出现 | — |
| 启用 LB Sprint | B31e 管理页可用; 需要至少 1 种积分类型 | C端 LB Sprint tab 出现 | — |
| 启用 Milestones | B31f 管理页可用 | C端 Milestones tab 出现 (有 active milestone 时) | — |
| 启用 Benefits Shop | B31g 管理页可用; 需要积分系统 (用积分兑换) | C端 Shop tab 出现 | — |
| 启用 Lucky Wheel | B31h 管理页可用; 需要积分系统 (spin cost) | C端出现抽奖入口 | — |
| 启用 Badges | B31i 管理页可用 | C端用户资料页显示 Badge 收集 | Badge 可用于 Milestone 奖励和 Shop 门控 |
| **禁用模块** | 管理页不可访问 (sidebar 项 dimmed); **已有数据保留** (不删除); status 变 Inactive | 对应 C端 tab/入口消失; **用户已获得的积分/badge/进度保留** | 禁用后重新启用: 数据恢复, 无需重新配置 |

### B.2 数据依赖级联

| 源模块 | 依赖模块 | 级联效果 |
|--------|---------|---------|
| Points & Level | → TaskChain, DayChain, Leaderboard, LB Sprint, Milestones, Benefits Shop, Lucky Wheel | 所有模块的积分来源; 删除积分类型前检查是否被引用 |
| Badges (B31i) | → Milestones (D06 奖励), Shop (D07 门控), TaskChain (D02 奖励) | 删除 badge 前检查引用; 已引用 badge 删除需确认 "This badge is used in X milestones and Y shop items" |
| Sectors & Tasks (B31) | → TaskChain (D02 步骤), DayChain (D03 daily task) | 删除/禁用被引用 task → 对应 Chain 显示警告 "Task '[name]' is no longer available" |
| B61 Integration (链上) | → B49 Access Rules (Token Gate) | 断开链上集成 → Token Gate 规则变 Inactive + 警告 |
| B61 Integration (Twitter) | → D20 Readiness Check | 断开 Twitter → D20 Twitter 检查失败; 影响所有 Publish/Activate 操作 |

### B.3 B端操作 → C端即时影响

| B端操作 | C端影响 | 生效时间 |
|--------|---------|---------|
| 发布新 Task | C端 Quests 页面出现新任务卡片 | 即时 (WebSocket 推送或 C端下次 API 请求) |
| 修改 Task 积分值 | C端任务卡片显示新积分值 | 即时 |
| 删除/隐藏 Task | C端任务卡片消失; 已完成用户保留积分 | 即时 |
| 发布公告 | C端首页 Carousel 出现 | 即时 |
| 修改 Level threshold | 用户等级可能变化 (升/降级); C端等级显示更新 | 5 分钟内 (服务端 batch 重算) |
| Shop 商品 Sold Out | C端商品显示 "Sold Out" 状态 | 即时 |
| 修改品牌色 | C端 Community 页面主题色更新 | 下次 C端页面加载时 |
| Homepage section 排序/显隐 | C端首页 section 顺序/可见性变化 | 下次 C端页面加载时 |

---

## 附录 C: 边界条件与异常处理 (v2.0 新增)

### C.1 并发操作

| 场景 | 处理方式 |
|------|---------|
| 两个管理员同时编辑同一 Task | **Last write wins**: 后保存的覆盖先保存的; 无冲突检测 (SaaS 标准做法); 考虑未来版本增加 `updatedAt` 乐观锁 |
| 两个管理员同时拖拽排序同一列表 | Last write wins; 后端按最后收到的 order 数组存储 |
| 管理员编辑 Modal 期间, 另一管理员删除了该项 | 保存时 API 404 → Modal 显示 "This item has been deleted by another user." + 关闭按钮 → 刷新列表 |

### C.2 权限与订阅限制

| 场景 | 处理方式 |
|------|---------|
| 试用期到期 (14天) | 所有 Publish/Activate 被 D20 拦截 → "Upgrade Plan →"; 已发布内容继续运行 (不下线); 新操作被阻止 |
| 免费计划模块上限 | Community 免费: 3 个模块; 尝试启用第 4 个 → toast "Upgrade to unlock more modules" + 模块 toggle 回弹 |
| 降级计划 | 已启用模块不强制禁用; 但无法启用新模块或创建新实例直到数量回到计划限额内 |
| 非管理员角色 | 只读访问: 所有编辑按钮 disabled + tooltip "Only admins can edit"; 导航仍可用 |

### C.3 数据边界

| 场景 | 限制 | 超限处理 |
|------|------|---------|
| 每个 Community 最大 Sector 数 | 20 | "+ New Sector" disabled + tooltip "Maximum 20 sectors reached" |
| 每个 Sector 最大 Task 数 | 100 | "+ New Task" disabled when sector has 100 tasks |
| 每个 Community 最大 Task 总数 | 500 | 全局 "+ New Task" disabled |
| Leaderboard 每种积分类型最大数 | 5 | 见 C-105 |
| Shop 商品最大数 | 50 | "+ Add Item" disabled |
| DayChain 最大数 | 10 | "+ Create Chain" disabled |
| TaskChain 最大步骤数 | 20 per chain | "+ Add Step" disabled in D02 |
| Badge 最大数 | 100 | "+ Create Badge" disabled |
| Access Rule 最大数 | 20 | "+ Create Rule" disabled |
| Homepage Section 最大数 | 15 | "+ Add Section" disabled |
| Announcement 最大活跃数 | 10 (同时显示) | "+ New Announcement" → toast "Archive some announcements first" |

### C.4 网络与恢复

| 场景 | 处理方式 |
|------|---------|
| API 请求超时 (10s) | 取消请求 + toast "Request timed out. Please try again." + UI 恢复 (乐观更新回滚) |
| 文件上传失败 (Image/CSV) | 进度条变红 + "Upload failed: {reason}" + "Retry" 按钮; 原有图片不变 |
| Draft 数量限制 | 每个 Community 最多 **1 个** active draft (非多 draft 系统); 新建 draft 覆盖旧 draft; 无 draft 列表页面 |
| WebSocket 断连 | 自动重连: 1s → 2s → 4s → 8s → 16s (exponential backoff); 5 次失败后停止 + 页面顶部 banner "Live updates paused. Refresh to reconnect." |
| 浏览器 localStorage 满 | 降级: 不缓存 filter/sort 状态到 localStorage; 功能正常但刷新后重置筛选 |

### C.5 数据为空的特殊处理

| 页面/组件 | 空数据场景 | 显示 |
|----------|----------|------|
| B11 Module Performance | 模块启用但无数据 (0 completions) | 卡片显示 "0" + "No activity yet. Create your first task to get started." |
| B12 AI Insights | 不足 7 天数据, AI 无法生成洞察 | Insights 区块: "AI insights will appear after 7 days of activity." (蓝色 info 卡片) |
| B54 Economy Chart | 无积分收支数据 | 空图表 + "No points data yet. Points activity will appear here after users start earning." |
| B31 任务列表 | 所有 sector 内无任务 | 每个空 sector 内: "No tasks in this sector. Click '+' to add one." |
| D18 Segment Panel | 该分群 0 用户 | 空表格 + "No users in this segment yet." |
