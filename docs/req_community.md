# Community 产品 B端前端开发需求文档

> 版本: v1.0 | 日期: 2026-03-06
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
| Modal | D01-D11, D16-D20 | 16 |
| **合计** | | **40** |

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

#### 交互逻辑

1. **Checklist 展开/折叠**: 每个 enrich/go_live 项可点击展开显示详情和操作按钮
2. **Progress Bar**: 自动计算 `completedCount / totalCount`，绿色填充
3. **Share 步骤**: 展开后显示链接复制、社交分享按钮、Promo Kit 生成器
4. **First 10 participants**: WebSocket 自动检测，实时更新计数
5. **Module 卡片**: 已激活模块显示绿色边框 + "Manage" 按钮；未激活显示 "+ Enable"
6. **Auto-complete items**: wizard 阶段步骤在创建完成后自动标记 ✅

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| 模块卡片 "Configure" / "Manage" | → B31a-B31i（对应模块管理页）|
| "+ Enable" 模块 | (API) 启用模块 → 刷新页面 |
| "Browse Configuration Templates →" | → B13 |
| "Preview your community" | → B33 (Preview Mode) |
| Share "Copy" | (action) 复制链接 |
| Social buttons | (ext) 预填分享链接 |
| "Generate Promo Kit" | → D19 Modal |
| Resources links | → (ext) 帮助中心 / M04 |

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

#### 交互逻辑

1. **Stats 卡片**: 纯展示，绿色/红色趋势箭头
2. **Checklist Banner**: 可折叠，显示剩余步骤数，点击展开跳转 B10 Guided 完整 Checklist
3. **Module Cards**: 每个显示关键指标 + 趋势 + "Manage" CTA
4. **Quick Actions**: 3 个快捷按钮直达常用操作

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Module "Manage" / "View" | → B31a-B31i |
| "+ Enable" 模块 | (API) 启用 → 刷新 |
| "Create Task" | → B31 Sectors & Tasks |
| "Add Reward" | → B31g Benefits Shop |
| "View Analytics" | → B54 Community Insights |

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

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Module cards | → B31a-B31i |
| "All Integrations →" | → B61 Integration Center |
| Integration "Configure" | → B61 对应集成配置 |
| Retention "Full Analytics →" | → B54 Community Insights |

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

**重要**: "Publish Community" 按钮点击后先触发 **D20 Publish Readiness Check** Modal（检查订阅状态 + Twitter 授权），全部通过后才执行发布。

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

#### 交互逻辑

1. **拖拽排序**: Sector 和 Task 都支持拖拽手柄排序 (drag handle icon)
2. **Sector 折叠/展开**: 点击 Sector header 可折叠其下任务列表
3. **任务行操作**: 每行有编辑✏️、复制、删除、显隐切换
4. **状态切换**: Task status toggle 直接在行内切换 active/draft/hidden
5. **Sector 显隐**: Sector header 有 visibility toggle

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "+ New Sector" | (action) 新增 sector 输入框 |
| "+ New Task" | (action) 新增 task 表单/modal |
| Task row edit ✏️ | (action/modal) 编辑任务详情 |
| "Publish Task" | (API) → 触发 D20 Readiness Check |

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
| B32 Content Mgmt | "Publish" (公告/内容) |

#### Checklist 项

| # | 检测项 | 通过条件 | 失败处理 |
|---|--------|---------|---------|
| 1 | 订阅状态 | 试用期内 OR 付费有效 | "Upgrade Plan →" 链接 |
| 2 | 官方 Twitter 授权 | 已连接并验证 | "Connect Twitter →" 链接 |

#### 交互逻辑

1. **全部通过**: 自动关闭 Modal → 执行发布 API
2. **部分未通过**: Publish 按钮 disabled，显示未通过项 + 解决链接
3. **检测中**: 显示 loading spinner

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
