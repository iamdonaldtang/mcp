# White Label 产品 B端前端开发需求文档

> 版本: v2.1 | 日期: 2026-03-07
> v2.1 变更: 基于 6 维度审计补充 197 项缺失: 每页 Init/Destroy 生命周期、按钮条件/disabled/loading、B44 集成配置详情页 (新增)、D19-WL/D20-WL 完整规格、§2.3 全局交互规范、B15→B16 过渡阈值、侧栏 Active 映射、所有 P0/P1/P2 项
> v2.0 变更: 新增每页操作详情表、完整 Modal 字段验证规则、B16/B42/B51 完整页面规格、级联效果、边界条件（共补充 112 项缺失内容）
> 基于设计稿 `design/pencil-new.pen` + 现有文档 `website_frontend_requirements.md` v4.2
> 供前端/后端工程师对照实施

---

## 目录

1. [模块概述](#1-模块概述)
2. [全局架构](#2-全局架构)
3. [Hub 页面（3 状态）](#3-hub-页面3-状态)
4. [创建向导（4 步）](#4-创建向导4-步)
5. [Domain 管理](#5-domain-管理)
6. [Embed & Deployment](#6-embed--deployment)
7. [Widget Library（3 状态）](#7-widget-library3-状态)
8. [Page Builder（3 状态）](#8-page-builder3-状态)
9. [Brand Settings](#9-brand-settings)
10. [SDK & API](#10-sdk--api)
11. [Integration Center](#11-integration-center)
12. [Smart Rewards](#12-smart-rewards)
13. [Page Analytics](#13-page-analytics)
14. [Dev Kit Page](#14-dev-kit-page)
15. [侧栏架构](#15-侧栏架构)
16. [API 接口汇总](#16-api-接口汇总)
17. [状态路由策略](#17-状态路由策略)
18. [D19 Promo Kit Generator — WL](#18-d19-promo-kit-generator--wl-v20-新增)
19. [D20 Publish Readiness Check — WL](#19-d20-publish-readiness-check--wl-v20-新增)
20. [附录 A: D20 Publish Touchpoints](#附录-a-d20-publish-touchpoints-wl)
21. [附录 B: 设计稿 Node ID 索引](#附录-b-设计稿-node-id-索引)
22. [附录 C: 级联效果 (CASCADE)](#附录-c-级联效果-cascade--v20-新增)
23. [附录 D: 边界条件 (EDGE)](#附录-d-边界条件-edge--v20-新增)

---

## 1. 模块概述

### 1.1 产品定位

White Label 是 TaskOn 的全栈品牌化产品，核心价值是"拥有体验"——项目方可以在自有域名、自有 App 内嵌入完整的 Growth Engine，实现自定义域名、SDK 集成、数据所有权。

**WL = Quest + Community + 独占功能**（Page Builder / Widget Library / Custom Domain / SDK / Deep Customization）

### 1.2 四种部署路径

| 路径 | 代码需求 | 操作者 | 特点 |
|------|---------|-------|------|
| Host on Your Domain | 无 | 市场 | 自定义域名，CNAME 配置，最快上线 |
| Embed in Your App ★推荐 | 低 | 市场+工程 | iframe/Widget/Page Builder，SSO 集成 |
| Build with SDK | 完整 | 工程 | API/SDK 全自定义，Headless 架构 |

### 1.3 页面编码总览

| 分类 | 页面编码 | 数量 |
|------|---------|------|
| Hub（3 状态） | B14, B15, B16 | 3 |
| 创建向导（4 步） | B37, B17/B57/B58/B59/B60, B38, B56 | 8 (含 5 个 Step 2 变体) |
| Embed Options | B19v (4 states) | 4 |
| Domain | B18 | 1 |
| Widget Library | B20, B21, B22 | 3 |
| Page Builder | B23, B24, B25 | 3 |
| Brand Settings | B40 | 1 |
| SDK & API | B41 | 1 |
| Iframe Embed | B42 | 1 |
| Integration Center | B26 | 1 |
| Integration Detail | B44 | 1 |
| Page Analytics | B43 | 1 |
| Contract Registry | B51 | 1 |
| Rule Builder | B52 | 1 |
| Privilege Manager | B53 | 1 |
| Dev Kit | B48 | 1 |
| **合计** | | **32** |

### 1.4 品牌色

- 产品主色: `#9B7EE0`（紫色）
- 页面背景: `#0A0F1A`
- 卡片背景: `#111B27`
- 边框: `#1E293B`
- Active item bg: `#1A1033`（紫色暗调）
- 强调边框: `#9B7EE0`

---

## 2. 全局架构

### 2.1 WL 侧栏子菜单

进入 WL 页面后，侧栏 White Label 项展开为带子菜单：

```
▾ White Label (expanded, chevron: keyboard_arrow_up)
  Overview         → B14/B15/B16
  Widgets          → B20/B22
  Pages            → B23/B25
  Smart Rewards    → B52 (Rule Builder) / B53 (Privilege Manager)
```

**样式规格**:
- Active item: purple bg `#1A1033` + text `#9B7EE0` + fontWeight 600
- Inactive: no fill + text `#94A3B8`
- Sub-item: padding `[8,12,8,40]`, fontSize 13, icons 16×16

### 2.2 共享布局

与 Community 相同的 Sidebar(240px) + TopBar(56px) + Content Area 布局。WL 页面的 Sidebar 中 "White Label" 高亮紫色。

### 2.3 全局交互规范 (Global Interaction Norms)

> 本节定义 WL 所有页面通用的交互规范，各页面不再重复说明。参见 `req_community.md` §2.3 中的完整定义，WL 继承并复用相同规范。

#### 2.3.1 Toast 通知
- 成功: 绿色 bg `#0A2E1A` + border `#16A34A` + check_circle icon, 3s 自动消失, 右上角, max 3 stack
- 错误: 红色 bg `#2D1515` + border `#DC2626` + error icon, 需手动关闭 (×), 包含错误原因文案
- 信息: 蓝色 bg `#0F1A2E` + border `#3B82F6` + info icon, 5s 自动消失

#### 2.3.2 Loading 模式
- **Skeleton**: 页面初始化加载, 各区块独立 skeleton (卡片/表格/图表各自有占位); sidebar + header 立即渲染
- **Inline spinner**: 按钮操作 (Save/Publish/Delete 等), 按钮内 spinner 替换文字, disabled 防重复提交
- **Overlay**: 全页面阻塞操作 (如 WL Publish), 半透明遮罩 + 居中 spinner + 文案

#### 2.3.3 空状态
- 插画 + 主文案 (16px bold) + 副文案 (14px gray) + CTA 按钮
- 表格空态: 表头保留 + body 区域居中显示空状态
- 列表空态: 整块替换为空状态卡片

#### 2.3.4 错误处理
- **API 失败 (非表单)**: skeleton → 3s 后显示 "Unable to load. Check your connection." + "Retry" 按钮 (见 §D.3)
- **表单提交失败**: toast error + 表单数据保留 (不清空) + 按钮恢复 enabled
- **路由守卫失败**: 重定向到 Hub (B14/B15/B16) + toast 说明原因
- **404 (无效 ID)**: 显示 "Page not found" + "← Back to White Label" 链接

#### 2.3.5 表单通用
- **脏数据检测**: 离开含未保存表单的页面 → 浏览器 `beforeunload` 确认弹窗 "You have unsaved changes. Leave anyway?"
- **字段验证**: 实时验证 (blur 触发) + 提交前全量验证; 错误: 红色 border + 字段下方红色错误文案
- **自动保存**: 仅 Wizard 草稿支持自动保存 (debounce 3s → `POST /api/white-label/drafts`); 其余页面需手动 Save

#### 2.3.6 拖拽 (Drag & Drop)
- 拖拽手柄: `drag_indicator` icon (⋮⋮)
- 拖拽中: 被拖元素 opacity 0.5 + 蓝色 border 2px `#3B82F6`; 占位符: 4px 蓝色虚线
- 放下: 300ms ease transition
- 触控: 支持 touch drag (mobile 预览模式)

---

## 3. Hub 页面（3 状态）

同一 URL `/white-label`，根据配置状态动态切换。

### 3.1 状态切换逻辑

| 条件 | 显示页面 |
|------|---------|
| 0 工具已配置 | B14 Empty |
| 1-4 工具已配置 | B15 Active |
| 5+ 工具已配置 AND Monthly Impressions ≥ 1000 AND ≥ 1 active deployment | B16 Management |

---

### 3.2 B14 — White Label Hub Empty

**设计稿**: Node `Ir6Tq`

#### 页面概述
- **用途**: 新用户首次进入 WL 的引导页
- **核心目标**: 引导选择部署路径 → 进入向导或工具页

#### 页面结构

```
Content Area
├── Community Ready Banner (绿色, 如有活跃 Community)
│   └── "Community Ready: Your Community has X tasks, Y points..."  + "View Community →"
├── Header
│   ├── Title: "White Label" (24px bold)
│   └── Subtitle: "Your branding, growth toolkit. All tools work together — start anywhere."
├── "CHOOSE YOUR DEPLOYMENT PATH" Label
├── Deployment Path Cards (3 cards)
│   ├── "Host on Your Domain" — zero code, CNAME
│   │   └── Tags: Custom Domain · Brand Settings
│   ├── "Embed in Your App" ★ RECOMMENDED (紫色边框 selected)
│   │   └── Tags: Widget Library · Page Builder · SSO
│   └── "Build with SDK" — full custom
│       └── Tags: API Keys · SDK Docs · Webhooks
├── CTA Block
│   ├── "Set Up White Label" (紫色按钮)
│   └── "Or start from scratch →"
├── "RESOURCES" Label
└── Resources Row (3 cards: SDK Documentation / Setup Walkthrough / Learn More)
```

#### Deployment Path Card 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | enum | `domain` / `embed` / `sdk` |
| title | string | 路径名 |
| description | string | 简述 |
| recommended | boolean | 是否推荐 |
| tags | string[] | 关联工具标签 |
| codeRequired | string | "None" / "Low" / "Full" |

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "View Community →" | → B10/B11/B12 Community Hub |
| Path card click | (action) 选中路径 |
| "Set Up White Label" | → B37 Wizard Step 1 (携带 path) |
| "Or start from scratch" | → B37 |
| "Custom Domain" tag | → B18 Domain Setup |
| "Widget Library" tag | → B20 |
| "Page Builder" tag | → B23 |
| "Brand Settings" tag | → B40 |
| "SDK & API" tag | → B41 |
| "Integration Center" tag | → B26 |
| Resources links | → (ext) 帮助 / M05 |

#### 操作详情 (v2.0 新增, W-01 ~ W-05)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-01 | 路径卡片选择 | click card | — | 选中: 紫色 border 2px `#9B7EE0` + 内容展开(适用场景+技术要求+预计时间); 同时只能选一个 (单选互斥); 默认 "Embed" 预选 | — |
| W-02 | CTA 文案变化 | path selected | — | 按钮文案动态: "Set Up with Embed" / "Set Up with Domain" / "Set Up with SDK"; 无选中时不会出现 (默认有 Embed 预选) | — |
| W-03 | 路径卡片展开详情 | click (selected card) | — | 选中后展开区 slide-down: 适用场景 (2-3 行) + 技术要求 ("None"/"Low code"/"Full stack") + 预计时间 ("5 min"/"30 min"/"2-4 hours") | — |
| W-04 | "Recommended" badge | render | — | **仅 Embed 路径**有 ★ RECOMMENDED 标识 (紫色 badge, 卡片右上角); 其他路径无 | — |
| W-05 | 路由守卫 (已有 WL) | route enter | `GET /api/white-label/status` | 用户已有 WL 配置 → **302 重定向到 B15/B16** (不应看到 Empty 页); 仅 `configuredTools=0` 时可访问 | — |

**页面初始化 (Init)**:
- 并行调用: `GET /api/white-label/status` (路由守卫 + configuredTools 判断), `GET /api/community/status` (Community Ready Banner 数据)
- 加载态: 各区块独立 skeleton, sidebar + header 立即渲染
- API 失败: 显示 §2.3.4 标准错误态 + Retry 按钮
- 路由守卫: `configuredTools > 0` → 302 到 B15 或 B16; `configuredTools === 0` → 显示 B14
- 路由守卫 API 失败: 默认显示 B14 (不拦截), toast 通知网络错误
- 路由守卫 API 缓存: 60s, 页面 focus 时 revalidate
- Community Ready Banner: `community.status === 'active'` 时显示; 数据含 community.name, community.taskCount, community.pointTypes

**页面离开 (Destroy)**:
- 无状态需清理 (纯展示页面)
- 路径卡片选择状态不持久化 (每次进入重置为 Embed 默认预选)

**按钮条件补充**:
- "Set Up White Label": 始终 enabled (D20 在 Wizard Step 4 Publish 时才触发, 非此处); 无订阅检查
- "View Community →" Banner: 仅 Community status=active 时 Banner 整体可见; API source: `GET /api/community/status`
- Tag 按钮 (Custom Domain, Widget Library 等): 始终 enabled, 纯导航快捷入口; 不因 plan 限制 disabled (plan 限制在各子页面内检查)

**侧栏**: WL 子菜单展开, "Overview" 高亮 active

**浏览器后退**: 从 Wizard 返回时应回到 B14 (如当初从 B14 进入)

---

### 3.3 B15 — White Label Hub Active

**设计稿**: Node `BnkYW`

#### 页面概述
- **用途**: Onboarding 阶段的引导工作区，含 Checklist + Toolkit
- **入口**: 发布 WL 后自动跳转

#### 页面结构

```
Content Area (gap: 32px)
├── Header: "White Label" + Stats subtitle + "Active" badge
├── Getting Started Checklist (#111B27 card, gap: 16px)
│   ├── Progress: "5 of 9 complete"
│   ├── ── CREATED BY WIZARD ──
│   │   ├── ✅ Community: Active (Arbitrum Builders, 12 tasks)
│   │   ├── ✅ Deployment path: Embed in Your App
│   │   ├── ✅ Brand & styling configured
│   ├── ── CONFIGURE YOUR TOOLS ──
│   │   ├── ○ Set up first Widget → Widget Library
│   │   ├── ○ Run your first widget → Widget Library (live)
│   │   ├── ○ Set up your Demo widget → Widget Library
│   │   ├── ○ Build your first page (optional) → Page Builder
│   │   ├── ○ Customize deployment page → maintenance page
│   ├── ── DEPLOY & LAUNCH ──
│   │   ├── ○ Integration verified → (auto-detect API ping)
│   │   ├── ○ Preview as a user → Open Preview
│   │   ├── ○ Announce to users → Share + Promo Kit
│   │   ├── ○ Send Dev Kit to developer → Copy link / Email
│   │   └── ○ First user interaction → (auto-detect)
├── Stats Row (4 cards)
│   ├── Domain: app.myproject.io
│   ├── Widgets Created: 5 active
│   ├── Pages Published: 2 live
│   └── Monthly Impressions: 12,450
├── "YOUR TOOLKIT" Label
├── Toolkit Grid (2 rows × 3)
│   ├── Row 1: Custom Domain / Widget Library / Page Builder
│   └── Row 2: Brand Settings / SDK & API / Integration Center
├── "RESOURCES" Label
└── Resources Row (3 cards)
```

#### Checklist 特殊项

**"Send Dev Kit" 步骤 (展开态)**:
- Dev Kit Link: `taskon.xyz/devkit/{project_id}` + "Copy" 按钮
- "Email to Developer" 按钮 → (API) 发送邮件
- Integration status indicator: WebSocket 自动检测首次 API ping

**"Announce" 步骤**:
- 与 Community B10 相同的 Share + Promo Kit 流程

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Checklist → Widget Library | → B20/B22 |
| Checklist → Page Builder | → B23/B25 |
| Checklist → Dev Kit Copy | (action) 复制链接 |
| Checklist → Email Developer | (API) 发送邮件 |
| Checklist → Open Preview | → B33 Preview Mode |
| Checklist → Promo Kit | → D19 Modal |
| Custom Domain card | → B18 |
| Widget Library card | → B22 |
| Page Builder card | → B25 |
| Brand Settings card | → B40 |
| SDK & API card | → B41 |
| Integration Center card | → B26 |

#### 操作详情 (v2.0 新增, W-06 ~ W-10)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-06 | Checklist 步骤展开内容 | click step row (accordion) | — | 每步展开详情: ① "Set up first Widget" → "Open Widget Library →" 按钮 → B20. ② "Run your first widget" → 状态检测 (已有 deployed widget?). ③ "Build your first page" → "Open Page Builder →" → B23. ④ "Integration verified" → WebSocket 状态指示器 (见 W-08). ⑤ "Preview as user" → "Open Preview" → B33. ⑥ "Announce to users" → Share 区 + "Generate Promo Kit" → D19. ⑦ "Send Dev Kit" → 见 W-07. ⑧ "First user interaction" → 见 W-10 | — |
| W-07 | "Send Dev Kit" | expand step | `POST /api/devkit/{project_id}/generate` | 展开显示: Dev Kit URL `taskon.xyz/devkit/{id}` + "Copy" 按钮 (clipboard → toast "Copied!") + "Open in New Tab" 按钮 (target=_blank) + "Email to Developer" 按钮 → 弹出 email input → `POST /api/devkit/{id}/send-email` → toast "Dev Kit link sent" | 生成失败 → toast error |
| W-08 | "Integration verified" 步骤 | auto (WebSocket) | `/ws/wl/integration-ping` | WebSocket 监听首次来自项目域名的 API ping; 状态指示器: 🔴 "Waiting for first API ping..." → 🟢 "Verified! First ping from {domain} at {time}"; 自动标记步骤 ✅ | WebSocket 断连 → 降级为 polling GET /api/wl/integration-status 每 30s |
| W-09 | Toolkit 卡片 "Configure →" | click | — | 6 个卡片各自跳转: Custom Domain → B18, Widget Library → B20/B22, Page Builder → B23/B25, Brand Settings → B40, SDK & API → B41, Integration Center → B26 | — |
| W-10 | "First user interaction" 步骤 | auto (WebSocket) | `/ws/wl/first-interaction` | WebSocket 监听首个来自 C端的用户交互事件 (task completion / widget click / page view); 检测到后自动标记 ✅ + toast "First user interaction detected!" | 同 W-08 降级 |

**页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/white-label/status` — Hub 状态 + configuredTools + stats row 数据 (Domain / Widgets / Pages / Impressions)
  2. `GET /api/white-label/onboarding` — Checklist 步骤完成状态 (9 步, 每步: id / label / completed / completedAt)
  3. WebSocket `/ws/wl/integration-ping` — 实时监听 API ping (Checklist step: Integration verified)
  4. WebSocket `/ws/wl/first-interaction` — 实时监听用户交互 (Checklist step: First interaction)
- 加载态: Header + Stats Row skeleton + Checklist skeleton (9 行占位) + Toolkit Grid skeleton
- Checklist 状态持久化: **服务端** (通过 onboarding API), 非 localStorage; 多设备同步
- Stats Row 数据源: `GET /api/white-label/status` response 中的 `stats` 字段 (domain / widgetsActive / pagesPublished / monthlyImpressions)
- Stats 刷新策略: 页面加载时获取, 无自动轮询; 从子页面返回时 revalidate (stale-while-revalidate)

**页面离开 (Destroy)**:
- WebSocket `/ws/wl/integration-ping`: 断开连接, 清理 listeners
- WebSocket `/ws/wl/first-interaction`: 断开连接, 清理 listeners
- Checklist 展开/折叠状态: 不持久化 (每次进入默认全折叠)

**B15 → B16 过渡条件**:
- 具体阈值: `configuredTools >= 5 AND monthlyImpressions >= 1000 AND deployments.length >= 1`
- `configuredTools` 定义: 已配置工具数 = (Domain configured ? 1 : 0) + (Widgets deployed count, max 1) + (Pages published count, max 1) + (Brand configured ? 1 : 0) + (SDK keys generated ? 1 : 0) + (Integrations connected count, max 1)
- `monthlyImpressions`: 过去 30 天所有 WL 资产的 page views 总和
- `deployments.length >= 1`: 至少有一个活跃部署路径 (Domain/Widget/Page 任一)
- 过渡为自动 (路由守卫判断), 用户无需手动触发
- 过渡后可回退: 如 impressions 降至 < 1000, 下次访问回到 B15

**Checklist 步骤展开行为**: Accordion 单个展开 (点击新步骤时自动折叠其他已展开步骤)

**"Email to Developer" (W-07) 验证补充**:
- Email input: 必填, email 格式验证 (debounce 300ms, blur 时验证)
- 无效格式: input 红色 border + "Please enter a valid email address"
- Send 按钮: disabled 当 email 为空或格式无效; 点击后 spinner + disabled 防重复; 成功: toast "Dev Kit link sent to {email}" + input 清空; 失败: toast error "Failed to send. Please try again."

**Toolkit 卡片 "Configure →" 显示条件**:
- 每张卡片始终可见 (6 张固定); "Configure →" CTA 始终 enabled
- 卡片显示当前状态 badge: "Active" (green) / "Not Set" (gray) / "Draft" (amber); 状态数据从 `GET /api/white-label/status` 的 `tools` 数组获取

**Loading skeleton**: Header (1 行 title + 1 行 stats 占位) + Checklist (9 行 shimmer) + Stats Row (4 卡片 skeleton) + Toolkit (6 卡片 skeleton)

**Checklist 完成 → B16 过渡**: 不自动跳转; 下次访问 `/white-label` 时路由守卫判断 (见上方阈值条件)

**"Open Preview" → B33**: **新标签页** (target=_blank), 保留 B15 页面状态

**B33 WL-specific Preview 说明**:
- WL Preview 与 Community Preview (req_community.md §B33) 使用相同的 B33 Preview Mode 页面
- **差异**: WL Preview 显示品牌化后的 C端视图 (WL brand colors/logo/fonts), 非 TaskOn 默认主题
- URL: `?from=B15&mode=white-label` → B33 加载 WL brand 配置渲染
- Widget preview: 显示已配置的 WL widgets (非 Community 模块直接视图)
- Domain preview: 如 Domain 已配置, 预览 URL 栏显示自定义域名 mock
- Preview 数据源: `GET /api/white-label/preview` (同 B56 W-40)
- Desktop/Mobile toggle: 支持 (1440px / 375px viewport)

**侧栏**: WL 子菜单展开, "Overview" 高亮 active

---

### 3.4 B16 — White Label Hub Management

**设计稿**: Node `UPAfV`

#### 页面概述
- **用途**: 高级管理视图，所有工具 + 部署 + 分析

#### 页面结构

```
Content Area (gap: 32px)
├── Header: "White Label" + stats + "Active" badge
├── Stats Row (4 cards)
│   ├── Primary Domain: app.myproject.io
│   ├── Widgets Created: 8 active
│   ├── Pages Published: 5 published
│   └── Monthly Impressions: 47,200
├── "YOUR TOOLKIT" Label
├── Toolkit Grid (2 rows × 3) — 同 B15
├── "ACTIVE DEPLOYMENTS" Label
├── Deployments Row (3 cards)
│   ├── app.myproject.io [Live] — Domain deployment
│   │   └── "12,450 impressions · CNAME update..." + "Healthy" badge
│   ├── Embedded Widgets [Active] — Widget deployment
│   │   └── "3 widgets embedded · 1,421 impressions..."
│   └── Published Pages [Active] — Page deployment
│       └── "Landing, Rewards, FAQ · 4 pages..."
├── "USAGE ANALYTICS" Label
└── Analytics Row
    ├── Monthly Impressions Chart (bar, 6 months)
    └── Key Metrics Panel (Avg Time, Bounce Rate, Interactions, Conversion)
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Toolkit cards | → B18/B22/B25/B40/B41/B26 |
| Deployment "Domain" | → B18 |
| Deployment "Widgets" | → B22 |
| Deployment "Pages" | → B25 |
| Analytics → Full | → B43 `/white-label/pages/:id/analytics` (项目级汇总视图) |

#### 操作详情 (v2.0 新增, W-11 ~ W-14)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-11 | 完整页面结构 | render | `GET /api/white-label/status` | **B16 是管理视图** (Node `UPAfV`): Header + 4 Stats + Toolkit Grid (同 B15) + Active Deployments (3 卡片: Domain / Widgets / Pages) + Usage Analytics (chart + metrics) | — |
| W-12 | Deployment Stats | render | `GET /api/white-label/deployments` | 4 stats: 已部署版本 ("v2.3, published 2d ago") / 活跃用户 (24h unique visitors) / 总交互数 (widget interactions + page completions) / API 调用数 (24h, from SDK usage) | — |
| W-13 | Feature Management 卡片 | render + click | — | 5 工具卡片 (Widget Library / Page Builder / Smart Rewards / Brand / SDK): 每个显示状态 badge (Active/Not Set/Draft) + 关键指标 (如 "5 widgets active") + "Manage →" 按钮跳转对应页面 | — |
| W-14 | "View All Deployments" | click | — | → 查看历史发布记录: 弹出 side panel (类似 D18), 显示版本号 / 发布时间 / 发布者 / 变更摘要; 数据源: `GET /api/white-label/deployments/history`; side panel 720px, 右侧滑入, 关闭: × 或点击外部; 分页: 20 条/页 | — |

**页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/white-label/status` — Hub 状态 + Stats Row (Domain / Widgets / Pages / Impressions)
  2. `GET /api/white-label/deployments` — Active Deployments 列表 (3 卡片: Domain / Widgets / Pages)
  3. `GET /api/white-label/analytics/summary` — Usage Analytics 图表 + Key Metrics (Avg Time / Bounce Rate / Interactions / Conversion)
- 加载态: 各区块独立 skeleton; Analytics chart 和 Key Metrics 各自独立 skeleton
- API 失败: 每区块独立显示 §2.3.4 错误态 + Retry

**页面离开 (Destroy)**:
- 无状态需清理

**Analytics Chart API**:
- 端点: `GET /api/white-label/analytics/summary?period=6m&granularity=month`
- 返回: `{ monthlyImpressions: [{month, views}], keyMetrics: {avgTimeOnPage, bounceRate, totalInteractions, conversionRate} }`
- Chart: 柱状图, 6 个月, 紫色 `#9B7EE0`, hover tooltip 显示具体数值
- Key Metrics: 4 个数值卡片, 无数据时显示 "—"
- 空数据 (新项目无流量): chart 显示 "No analytics data yet. Deploy widgets or pages to start tracking." + 空柱状图占位

**W-13 Feature Management 卡片状态判定逻辑**:
- **Widget Library**: Deployed widgets > 0 → "Active" (green); Configured but not deployed > 0 → "Draft" (amber); 0 configured → "Not Set" (gray)
- **Page Builder**: Published pages > 0 → "Active" (green); Draft pages > 0 → "Draft" (amber); 0 pages → "Not Set" (gray)
- **Smart Rewards**: Active rules > 0 OR Active privileges > 0 → "Active" (green); Draft rules/privileges > 0 → "Draft" (amber); 0 → "Not Set" (gray)
- **Brand**: logo 或 primary color 自定义过 → "Active" (green); 未修改过 → "Not Set" (gray)
- **SDK**: API key 已生成 → "Active" (green); 未生成 → "Not Set" (gray)
- 数据来源: `GET /api/white-label/status` response 中的 `tools[]` 数组, 每项含 `{ id, status, keyMetric }`

**Deployment 卡片**:
- 每张卡片点击区域: **整张卡片** 可点击 (cursor: pointer), 跳转对应管理页
- "Healthy" badge 判定: 过去 24h 无 error 日志 + uptime > 99% → "Healthy" (green); 有 warning → "Warning" (amber); 有 error → "Error" (red)
- 部署状态应反映 C端实际状态 (来自 deployments API response)

**"View Full Analytics" 按钮**: 跳转 B43 项目级分析视图; 无数据时按钮仍 enabled (B43 会显示空态)

**Breadcrumb**: "White Label" (不可点, 当前页)

**Loading skeleton**: Header + 4 Stats skeleton + Toolkit Grid (6 卡片 skeleton) + Deployments (3 卡片 skeleton) + Analytics (chart skeleton + 4 metrics skeleton)

**Empty analytics chart** (新项目): 显示空柱状图 + overlay 文案 "No analytics data yet"

**侧栏**: WL 子菜单展开, "Overview" 高亮 active

---

## 4. 创建向导（4 步）

### 4.1 向导流程

```
B37 Path → B17/B57-B60 Configure (path-adaptive) → B38 Brand → B56 Preview
```

**Step 2 有 5 个变体**，根据 Step 1 选择的路径显示不同配置页：

| Path | Step 2 页面 | Node ID |
|------|-----------|---------|
| Embed (Widget) | B17 Widget Config | `CXzmy` |
| Domain | B57 Domain Config | `YGODW` |
| Embed (Iframe) | B58 Iframe Config | `Kr5W5` |
| Embed (Page Builder) | B59 PB Config | `XHwzp` |
| SDK | B60 SDK Config | `eNFmU` |

---

### 4.2 B37 — Step 1: Choose Deployment Path

**设计稿**: Node `NNwid`

#### 页面结构

```
Stepper: ①Path ②Configure ③Brand ④Preview
Content Area (centered)
├── Title: "Choose Your Deployment Path"
├── Subtitle: "Pick how you want to deploy. Other modes unlock for your needs."
├── Path Cards (3, horizontal)
│   ├── "Embed in Your App" ★ RECOMMENDED (selected, 紫色边框)
│   │   └── Tags: Widget Library · Page Builder · SSO Auth
│   ├── "Host on Your Domain" — CNAME 配置
│   │   └── Tags: DNS Setup · Full Portal · Zero Code
│   └── "Build with SDK" — 完全自定义
│       └── Tags: API Keys · SDK Docs · Full Control
└── Action Bar: [Cancel] [Next: Configure] (紫色)
```

#### 交互逻辑

1. **单选互斥**: 选中路径高亮紫色边框，其余默认
2. **★ RECOMMENDED 标签**: "Embed in Your App" 默认选中
3. **Tags**: 纯展示，说明该路径涉及的工具

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Cancel" | → B14/B15/B16 (WL Hub) |
| "Next: Configure" | → B17/B57/B58/B59/B60 (根据选择) |

#### 操作详情 (v2.0 新增, W-15 ~ W-19)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-15 | Embed 路径选中后注释 | select Embed | — | 卡片下方: "Next: Choose embed method (Iframe / Widget Library / Page Builder)" (灰色 14px 副文) | — |
| W-16 | Domain 路径选中后注释 | select Domain | — | "Next: Configure your custom domain and CNAME records" | — |
| W-17 | SDK 路径选中后注释 | select SDK | — | "Next: Generate your API keys and integration code" | — |
| W-18 | "Next" 按钮路由 | click Next | — | Embed → B17 (Widget Config, 默认); 在 B17 内可切换到 Iframe/PB variant. Domain → B57. SDK → B60 | — |
| W-19 | 无选择时 Next | — | — | **Next disabled** (opacity 0.5 + cursor not-allowed); 实际不会出现因为 Embed 默认预选; 但若通过代码清空选择, 则 disabled | — |

**页面初始化 (Init)**:
- **静态页面**: 无 API 调用, 路径卡片数据硬编码
- 草稿恢复: 如有 wizard draft (`GET /api/white-label/drafts`), 恢复之前选择的 path; 无 draft → Embed 默认预选
- 加载态: N/A (静态内容)

**页面离开 (Destroy)**:
- 路径选择状态: 通过 wizard 共享状态管理 (React Context / Zustand), 不单独持久化
- "Cancel" → 返回 Hub (B14/B15/B16, 根据 §17.1 状态路由)
- "Cancel" 确认弹窗: 如已有 draft 修改 (选了非默认 path) → "Discard changes?" 确认; 否则直接跳转

**W-18 Embed → B17 内切换机制**:
- 进入 B17 后, 页面顶部有 Embed Method radio (W-20): Widget Library / Iframe / Page Builder
- 切换到 Iframe: **同标签页内路由跳转**到 B58 (`/white-label/wizard/iframe-config`), 非新标签页
- 切换到 Page Builder: 同标签页路由跳转到 B59 (`/white-label/wizard/pb-config`)
- 切换保留 wizard 共享状态 (已选 path = Embed), stepper 仍显示 Step 2 active
- **非 in-page swap**: 是完整的路由切换 (B17 → B58 / B59), URL 变化, 组件卸载/挂载

**浏览器后退**: Step 2 → Step 1, wizard 状态保留 (path 选择仍在)

---

### 4.3 B17 — Step 2: Widget Configure (Embed path)

**设计稿**: Node `CXzmy`

#### 页面结构

```
Stepper: ①Path ②Configure(active) ③Brand ④Preview
Left Panel: Widget Selection
├── Title: "Choose Your Widgets"
├── Subtitle: "Select which Community modules to embed..."
├── Widget Checklist (based on Community modules)
│   ├── ✅ Leaderboard — Ready to embed, metric: 1,245pts, 148 participants
│   ├── ✅ Task List — 24 active tasks, 148 participants
│   ├── ✅ User Center — Points balance, levels/progress
│   ├── ○ DayChain — "Needs Community setup → Set up in Community →"
│   └── ○ Benefits Shop — "Needs Community setup → Set up in Community →"
├── "Selected 3 widgets — you'll configure each one after setup."

Right Panel: Embed Preview
├── "EMBED PREVIEW" label
├── Preview card: yourproject.community
│   ├── Leaderboard widget mock
│   └── Task widget mock
└── Action Bar: [Back] [Next: Brand] (紫色)
```

#### Widget 状态

| 状态 | 显示 | 说明 |
|------|------|------|
| Ready (green check) | ✅ + 指标 | Community 中已配置 |
| Needs Setup (gray) | ○ + "Set up in Community →" | Community 中未启用 |

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Set up in Community →" | → B10/B11 Community Hub |
| "Back" | → B37 |
| "Next: Brand" | → B38 |

#### Step 2 各变体字段规格 (v2.0 新增, W-20 ~ W-36)

**B17 — Widget Config (Embed path)**

| # | 字段 | 类型 | 必填 | 验证 | 说明 |
|---|------|------|------|------|------|
| W-20 | Embed Method | radio | ✅ | Widget Library (default) / Iframe / Page Builder | 选择影响后续: Widget→继续本页选 widgets; Iframe→跳转 B58; PB→跳转 B59 |
| W-21 | Widget Selection | checkbox list | ✅ (≥1) | 至少选 1 个 module | 列表来源: Community 已配置模块 (green ✅); 未配置模块 (gray ○ + "Set up in Community →") |
| W-22 | SSO Method | radio | ✅ | Wallet Auth / OAuth2+JWT / None (preview only) | None = 只能预览无法交互 |
| W-23 | Target Domain | URL input | ✅ | valid domain format (no protocol) | 将嵌入的页面所在域名, 用于 CORS 配置; placeholder "app.yourproject.io" |

**页面初始化 (Init)**:
- API: `GET /api/community/modules/status` — 获取 Community 各模块配置状态 (用于 Widget Checklist)
- 加载态: Widget Checklist 区域 skeleton (6 行占位); 右侧 Embed Preview skeleton
- API 失败: Checklist 区域显示 §2.3.4 错误态 + Retry
- Community 无配置模块 (全部 "Needs Setup"): 所有模块显示 "Set up in Community →"; "Next: Brand" **disabled** (需至少选 1 个 Ready widget)
- 草稿恢复: 从 wizard draft 恢复已选 widget list + SSO method + target domain

**页面离开 (Destroy)**:
- Widget 选择状态保存到 wizard 共享状态 (非 API)
- 未保存的 SSO / Target Domain 输入: 保留在 wizard 状态中 (内存)

**"Set up in Community →" 链接**: **新标签页** (target=_blank), 防止丢失 wizard 状态; 当前 tab 保留 B17 页面

**"Next: Brand" disabled 视觉**:
- 条件: 0 个 Ready widget 被勾选 (W-21 要求 ≥1)
- 样式: opacity 0.5 + cursor not-allowed + tooltip "Select at least 1 widget to continue"

**SSO Method (W-22)**: 位于 Widget Checklist 下方, Target Domain (W-23) 之前; 三个选项垂直排列

**Target Domain (W-23)**: 位于 SSO Method 下方, Action Bar 之前

**Embed Method 切换 (W-20)**: 切换到 Iframe → 路由跳转 B58; 切换到 PB → 路由跳转 B59; wizard 共享状态保留

**右侧 Embed Preview**: 实时更新; widget 勾选/取消 → preview 中对应 widget block 出现/消失; 动画: 300ms fade-in/out

**Loading skeleton**: Widget Checklist (6 行 shimmer) + Embed Preview (卡片 skeleton)

**B57 — Domain Config**

| # | 字段 | 类型 | 必填 | 验证 | 说明 |
|---|------|------|------|------|------|
| W-24 | Custom Domain | text input | ✅ | valid domain format; 不含 protocol/path | "community.yourproject.io" |
| W-25 | CNAME Target | display (readonly) | — | — | TaskOn 提供: `custom.taskon.xyz`; Copy 按钮 |
| W-26 | DNS Provider | select | — | Cloudflare / Route53 / Namecheap / GoDaddy / Other | 选择后显示对应平台配置教程链接 (外部) |
| W-27 | Verify DNS | button | — | — | 手动触发: `POST /api/white-label/verify-domain` → 显示 polling 状态 ("Checking..." spinner → "✓ Verified" green / "✗ Not found yet, retrying..." amber → 30s 后 auto-retry); **DNS 验证不阻塞 "Next: Brand"**: 用户可输入域名后直接点 Next, 稍后在 B18 完成验证 |

**B57 "Next: Brand" 按钮条件**:
- **必填**: Custom Domain (W-24) 非空 + 格式合法
- **不要求**: DNS 验证通过 (用户可稍后在 B18 完成验证)
- Disabled: domain 为空 → opacity 0.5 + tooltip "Enter a custom domain to continue"

**B57 DNS 验证 polling**: 手动触发 (点击 "Verify DNS" 按钮), 非页面加载自动开始; 触发后 30s 间隔自动重试, 最多 10 次; 超时后停止并显示 "DNS propagation may take up to 48 hours. You can verify later in Domain Settings."

**B57 验证超时**: polling 持续上限 5 分钟 (10 次 × 30s); 超时后显示 amber 提示, "Verify DNS" 按钮恢复为可点击状态

**B58 — Iframe Config**

| # | 字段 | 类型 | 必填 | 验证 | 说明 |
|---|------|------|------|------|------|
| W-28 | Iframe URL | display (auto) | — | — | `share.taskon.io/embed/{project_id}`; Copy 按钮 |
| W-29 | SSO Option | radio | — | With SSO / Without SSO (read-only mode) | Without SSO: iframe 内容可浏览但无法交互 |
| W-30 | Allowed Origins | text input | — | 逗号分隔的域名列表; 各域名 valid format | 允许嵌入的域名白名单; 空 = 允许所有 (不推荐) |
| W-31 | Iframe Code | code block (readonly) | — | — | 自动生成 `<iframe src="..." width="100%" height="600" frameborder="0"></iframe>`; Copy 按钮 |

**B58 页面初始化 (Init)**:
- Iframe URL 中的 `project_id`: 从 wizard 共享状态获取; 首次创建时由后端在 draft 创建时分配 (临时 ID, publish 时确认)
- 静态页面, 无需额外 API
- Copy 按钮 clipboard 失败: fallback 为选中文本 + toast "Press Ctrl+C to copy"

**B58 "Next: Brand" 条件**: **始终 enabled** (W-28~W-31 所有字段均 optional); 点击即可继续

**B58 "Back" 目标**: → B17 (非 B37); 因为 B58 从 B17 的 Embed Method 切换进入, 应返回 B17

**B58 加载/错误态**: 无加载态 (静态内容); 无 API 错误态

**B59 — PB Config (has-pages variant)**

| # | 字段 | 类型 | 必填 | 验证 | 说明 |
|---|------|------|------|------|------|
| W-32 | 已有 Pages 列表 | radio list | ✅ (≥1) | — | 选择基于哪个已有 Page 继续 (Node `zW40A`); 显示 page name + 创建时间 |
| W-33 | 无已有 Pages 时 | template cards | ✅ (≥1) | — | 显示 Page Builder 模板选择 (Node `XHwzp`): Rewards Hub / Community Portal / Custom |

**B59 页面初始化 (Init)**:
- API: `GET /api/white-label/pages` — 获取已有 pages 列表 (用于判断 has-pages vs no-pages 分支)
- 分支逻辑: `pages.length > 0` → 渲染 W-32 (已有 Pages 列表, radio 选择); `pages.length === 0` → 渲染 W-33 (模板卡片选择)
- 加载态: skeleton (页面列表或模板卡片区域)
- API 失败: §2.3.4 错误态 + Retry

**B59 "Next: Brand" 条件**:
- has-pages (W-32): 至少选择 1 个 Page (radio, 必选); 未选 → disabled + tooltip "Select a page to continue"
- no-pages (W-33): 至少选择 1 个 template; 未选 → disabled + tooltip "Select a template to continue"

**B59 边界: 无 pages 且无 templates**: 显示空态 "No pages or templates available. Create widgets first in Widget Library." + "Go to Widget Library →" 链接

**B59 "Back" 目标**: → B17 (从 B17 Embed Method 切换进入)

**B59 Loading state**: pages 列表或 template 卡片区域 skeleton

**B60 — SDK Config**

| # | 字段 | 类型 | 必填 | 验证 | 说明 |
|---|------|------|------|------|------|
| W-34 | SDK Mode | radio | ✅ | Full Custom (Headless) / Hybrid (SDK + TaskOn UI fallback) | Full Custom: 全部自定义, 仅 API; Hybrid: SDK 组件 + fallback UI |
| W-35 | API Key | display (auto) | — | — | 自动生成 `pk_live_xxx`; Copy 按钮; "Regenerate" 按钮 (确认 dialog) |
| W-36 | Project ID | display (readonly) | — | — | 唯一标识符; Copy 按钮 |

**B60 页面初始化 (Init)**:
- API Key 自动生成: 首次进入 B60 时调用 `POST /api/white-label/sdk/keys` 自动生成 key pair; 已有 key 时调用 `GET /api/white-label/sdk` 获取现有 key
- 生成时机: 页面 mount 时自动触发 (无需用户点击)
- 加载态: API Key 区域显示 "Generating API key..." + inline spinner; 其余字段可用
- 生成失败: 显示 amber 提示 "Failed to generate API key. Please try again." + "Retry" 按钮; "Next: Brand" 仍然 **enabled** (key 可在 B41 稍后生成)

**B60 Project ID 来源**: 在 wizard draft 创建时由后端分配 (UUID v4); 如已有 WL project → 使用现有 project ID

**B60 "Regenerate" (W-35) 在 Wizard 中的行为**:
- 首次进入 wizard 时 key 刚生成, "Regenerate" 行为等同于重新生成 (旧 key 失效, 新 key 替换)
- 确认弹窗同 W-78 (B41): "⚠️ This will invalidate your current key."
- 实际上首次 wizard 场景极少使用 (key 刚创建), 但保留按钮以保持与 B41 一致

**B60 "Next: Brand" 条件**: **始终 enabled** (SDK Mode 有默认值 "Full Custom"; API Key 生成失败也不阻塞); 点击时将 SDK 配置写入 wizard 共享状态

**B60 Key 生成失败错误态**: amber 提示卡片 + "Retry" 按钮; key 区域显示 "—" 占位

---

### 4.4 B38 — Step 3: Brand Customization

**设计稿**: Node `5nCtO`

#### 页面结构

```
Stepper: ①Path ②Configure ③Brand(active) ④Preview
Left Panel: Brand Form
├── Project Logo: drag/drop upload zone
├── Brand Colors
│   ├── Primary: color picker (default #48E5E1)
│   ├── Secondary: color picker (default #7C3AED)
├── Typography
│   ├── Heading Font: dropdown (Inter, etc.)
│   ├── Body Font: dropdown
│   ├── Letter Spacing: slider
│   └── Font Size: input (16px)
├── Button Style
│   ├── Primary Button preview (filled)
│   └── Secondary Button preview (outline)

Right Panel: Live Preview
├── "LIVE PREVIEW" label + "Light/Dark" toggle
├── Preview mock (path-adaptive)
│   ├── "YourProject" header with logo
│   ├── "Welcome to YourProject" heading
│   ├── Button previews
│   └── Widget appearance preview
├── "Shaded text with the font size at select-n-review"
└── Action Bar: [Back] [Next: Preview] (紫色)
```

#### Brand 数据模型

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| logo | file (image) | ❌ | SVG/PNG, max 2MB |
| primaryColor | hex string | ✅ | 主品牌色 |
| secondaryColor | hex string | ✅ | 次品牌色 |
| headingFont | string | ✅ | 标题字体 |
| bodyFont | string | ✅ | 正文字体 |
| letterSpacing | number | ❌ | 字间距 |
| fontSize | number | ❌ | 基础字号 |
| buttonStyle | enum | ✅ | `filled` / `outline` / `rounded` |

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Back" | → B17/B57-B60 |
| "Next: Preview" | → B56 |

#### 操作详情 (v2.0 新增, W-37 ~ W-39)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-37 | Brand 字段 (同 B40 子集) | form | — | 字段: Logo (upload, 可选) + Primary Color (hex, 必填) + Secondary Color (hex, 必填) + Heading Font (select, 必填) + Body Font (select, 必填) + Button Style (radio: filled/outline/rounded, 必填); 验证同 B40 | — |
| W-38 | 右侧 Live Preview 更新 | any form change | — | Logo/颜色/字体变化即时反映在右侧 preview mock; debounce 200ms; preview 显示: header (logo + project name) + heading (font preview) + button (style preview) + widget mock (颜色主题) | — |
| W-39 | "Skip for now" | click link | — | 灰色文本链接 "Skip for now — you can configure this later in Brand Settings"; 点击 → 跳过 B38, 使用默认品牌设置, 直接到 B56; toast "You can customize branding anytime from Brand Settings" | — |

**B38 页面初始化 (Init)**:
- API: `GET /api/white-label/brand` — 如已有品牌设置 (如 B40 先前配置), 加载并预填; 否则使用默认值
- 加载态: form skeleton + preview skeleton
- 草稿恢复: 从 wizard draft 恢复已输入的 brand 数据

**B38 页面离开 (Destroy)**:
- Brand 数据保存: 未保存的 brand 数据保留在 wizard 共享状态 (内存); 不自动调用 API
- 脏数据: "Back" 按钮点击时如有修改 → "Unsaved changes" 确认弹窗; "Skip" 不提示 (用户明确跳过)

**B38 Brand 数据保存时机**: **点击 "Next: Preview" 时** 将 brand 数据写入 wizard 共享状态; **不立即** 调用 `PUT /api/white-label/brand` (品牌设置在 B56 Publish 时统一提交); "Skip" 跳过时使用 TaskOn 默认 brand 值

**B38 "Back" 目标**: → 之前的 Step 2 页面 (B17/B57/B58/B59/B60); 系统从 wizard 共享状态中获取 `selectedPath` 来决定返回哪个 Step 2 页面

**Logo 上传在 Wizard 中**: 验证同 B40 (W-71/W-72); 上传到临时存储, publish 时持久化
**Color picker 关闭行为**: 点击 popover 外部关闭; Esc 键关闭; 选择即生效 (无需确认按钮)

---

### 4.5 B56 — Step 4: Preview & Publish

**设计稿**: Node `WsH2y`

#### 页面结构

```
Stepper: ①Path ②Configure ③Brand ④Preview(active)
Left Panel: Deployment Preview
├── Title: "Deployment Preview"
├── Subtitle: "This is how embedded widgets will appear on your website."
├── Preview Frame (path-adaptive)
│   ├── yourproject.io mockup
│   ├── "Community" section
│   │   ├── Leaderboard widget
│   │   ├── Tasks widget
│   │   └── User Center widget
│   └── Footer

Right Panel: Launch Readiness
├── "Launch Readiness" heading
├── "Review before publishing:"
├── Checklist:
│   ├── ✅ Community: Publish Builders (12 tasks)
│   ├── ✅ Path: Embed in Your App
│   ├── ✅ 3 widgets selected
│   ├── ✅ Brand & style configured
│   ├── ⚠️ Begin config for your widgets (amber)
│   └── ⚠️ Page Builder not started yet (amber)
├── "After publishing" info card:
│   ├── "Build widgets in Widget Library (B20)"
│   ├── "Send Dev Kit to your developer"
│   ├── "Track performance in Analytics"
└── Action Bar: [Back] [Publish White Label] (紫色)
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Back" | → B38 |
| "Publish White Label" | (API) `POST /api/white-label/publish` → D20 → 成功后跳转 B15 |

**重要**: 发布按钮先触发 D20 Publish Readiness Check。

#### 操作详情 (v2.0 新增, W-40 ~ W-43)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-40 | WL Preview 内容 | render | `GET /api/white-label/preview` | 左侧 Deployment Preview 根据已配置路径显示: **Embed** → yourproject.io mockup + 已选 widget 模块预览; **Domain** → custom domain portal 预览; **SDK** → API 集成代码 + Headless 架构图; 顶部 desktop/mobile toggle 切换视图 | API 失败 → 显示 "Unable to load preview" + Retry |
| W-41 | Readiness Checklist | render | `GET /api/white-label/readiness` | 右侧 checklist 项 (WL 版): ✅ Community active (name + task count) / ✅ Deployment path selected / ✅ Widgets selected (count) / ✅ Brand configured / ⚠️ "Begin config for your widgets" (amber, 无 widget 已配置时) / ⚠️ "Page Builder not started" (amber, 选了 PB 但未创建 page 时); Domain 路径额外: ✅/⚠️ DNS verified; "After publishing" info card 显示后续步骤 | — |
| W-42 | "Publish White Label" | click | `POST /api/white-label/publish` | 点击 → **先弹出 D20** Publish Readiness Check (订阅+Twitter 验证) → D20 通过 → 调用 publish API → button spinner → 成功: toast "White Label published!" + 自动跳转 B15 Hub Active; Dev Kit URL 自动生成 (见 W-43) | API 失败 → toast error "Publish failed: {reason}" + 留在 B56 |
| W-43 | Dev Kit 自动生成 | publish success | `POST /api/devkit/{project_id}/generate` | Publish 成功后自动触发 Dev Kit 生成: URL = `taskon.xyz/devkit/{project_id}`; B15 Checklist "Send Dev Kit" 步骤自动可用; 生成失败不阻塞 publish (后台重试) | 生成失败 → B15 "Send Dev Kit" 步骤显示 "Generating..." + 后台 30s 重试 |

**B56 页面初始化 (Init)**:
- **并行调用**: `GET /api/white-label/preview` (左侧 Deployment Preview) + `GET /api/white-label/readiness` (右侧 Launch Readiness Checklist)
- Preview 可能较慢 (渲染 mock): 左侧独立 skeleton, 右侧 checklist 可先渲染
- Readiness API 在 mount 时调用 (非 Publish 点击时)
- 加载态: 左侧 Preview 区域全面 skeleton (卡片占位 + shimmer); 右侧 Checklist 显示 loading spinner per item
- API 失败: Preview → "Unable to load preview" + Retry; Readiness → "Unable to check readiness" + Retry

**B56 页面离开 (Destroy)**:
- Draft 保留: 离开不删除 wizard draft; 用户可从 B14/B15 重新进入 wizard 继续
- 无自动保存 (Wizard 最后一步, 无新输入)

**B56 "Publish White Label" 按钮条件**:
- ⚠️ Amber 项为**信息性提示** (informational), **不阻塞** Publish 按钮
- Publish 始终 enabled (只要 readiness API 返回成功)
- 点击 Publish → 弹出 D20 Publish Readiness Check (订阅 + Twitter) → D20 通过后执行 publish API
- Readiness API 失败: Publish 按钮 disabled + "Unable to verify readiness. Please retry."

**Desktop / Mobile preview toggle**: 左侧 Preview 区域顶部, 两个 icon 按钮 (desktop_windows / smartphone); 切换 preview 宽度 (1440px mock → 375px mock); 当前选中 icon filled purple, 另一个 outline gray

**B56 Publish 后 C端影响**:
- 已选 widgets 的 embed code 激活 (访问 embed URL 可见)
- Domain 路径: 自定义域名 portal 上线 (如 DNS 已验证)
- SDK 路径: API key 激活, 可开始调用
- Page Builder: 已选 page 的 share URL 激活
- Dev Kit 自动生成 (W-43)

**B56 Publish 后页面跳转**: D20 通过 → publish API 成功 → toast "White Label published!" → 自动跳转 B15 (Hub Active); B15 会 refetch onboarding 数据

**Preview loading state**: 左侧 Preview 全区域 skeleton (200ms shimmer)

**Readiness API 失败错误态**: 右侧 Checklist 区域显示 "Unable to check readiness" + Retry 按钮; Publish 按钮 disabled

---

## 5. Domain 管理

### 5.1 B18 — Domain Setup

**设计稿**: Node `5bmH9` | URL: `/white-label/domain`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "Domain Management"
├── ── ① Custom Domain ──
│   ├── Status: community.yourproject.com
│   ├── Domain input: HTTPS:// + text input
│   └── Help text
├── ── ② DNS Configuration ──
│   ├── DNS records table
│   │   ├── Type: CNAME | Name: ... | Value: custom.taskon.xyz
│   └── "Copy DNS Record" button
├── ── ③ Brand & Portal ──
│   ├── Logo upload
│   ├── Primary Color picker
│   └── Preview link

Right Panel: Domain Status
├── Domain Status checklist:
│   ├── ✅ Domain added
│   ├── ✅ DNS verified
│   ├── ✅ SSL active
│   └── ✅ Brand applied
├── Portal Preview card (mini C-end preview)
├── Domain Info card: "DNS records can take..."
└── "Edit Domain Settings" button → 触发 D20
```

#### DNS 验证

| Endpoint | Method | Protocol |
|----------|--------|----------|
| `/api/white-label/verify-domain` | POST | HTTP 轮询 (每 10s) |

#### 操作详情 (v2.0 新增, W-44 ~ W-48)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-44 | DNS 验证自动轮询 | page load + manual | `POST /api/white-label/verify-domain` | 页面加载后自动开始 polling (每 30s); 状态指示器: "Checking..." (spinner) → "✓ DNS Verified" (green) / "✗ Not found yet, will retry..." (amber); 手动 "Verify Now" 按钮可立即触发一次检查; 成功后停止 polling | 3 次连续失败 → 降低频率到每 60s; 显示 "DNS propagation may take up to 48 hours" |
| W-45 | CNAME 记录展示 | render | — | DNS Configuration 表格: Type=`CNAME` / Host=`community` (或用户自定义子域名) / Value=`custom.taskon.xyz` / TTL=`300`; "Copy DNS Record" 按钮 → 复制整行文本格式 → toast "Copied!"; Provider 教程链接 (Cloudflare/Route53/Namecheap/GoDaddy) | — |
| W-46 | SSL 状态 | auto (after DNS) | — | DNS 验证通过后自动触发 SSL 配置 (Let's Encrypt); 右侧 Domain Status checklist: ✅ Domain added → ✅ DNS verified → ⏳ "SSL: Provisioning..." (spinner) → ✅ "SSL: Active" (green); SSL 通常 2-5 分钟; Portal Preview 在 SSL active 后才显示 HTTPS 锁图标 | SSL 失败 → "SSL provisioning failed. Retrying..." + 自动重试 3 次 |
| W-47 | "Edit Domain Settings" | click | — | 进入编辑模式: Domain input 变为可编辑; **修改域名后点 Save → 触发 D20** (因为域名变更需要重新验证 DNS + SSL); 注意: 是 save 触发 D20, 不是 edit 按钮本身 | — |
| W-48 | 域名冲突 | verify-domain response | `POST /api/white-label/verify-domain` | 如果域名已被其他 TaskOn 项目 claim: API 返回 409 → 显示红色 error: "This domain is already claimed by another project. Please use a different subdomain." + domain input 红色 border; 不会自动重试 | 已 claimed 时 Save 按钮 disabled |

**B18 页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/white-label/domain` — 当前域名配置 (domain / dns_status / ssl_status / brand_applied)
  2. `GET /api/white-label/brand` — Brand & Portal 区域数据 (logo / primary color)
- 加载态: 各 section (Domain / DNS / Brand / Status Panel) 独立 skeleton
- API 失败: 各 section 独立显示 §2.3.4 错误态 + Retry
- DNS 自动 polling: 页面加载后, 如 dns_status !== 'verified', 自动开始 polling (W-44 规则: 30s 间隔)
- 无域名配置 (首次访问): domain input 为空, DNS section 显示占位文案, Status checklist 全部 ○ 未完成

**B18 页面离开 (Destroy)**:
- DNS polling: 停止 polling, 清除 timer
- 未保存的域名修改: 显示 "Unsaved changes" 确认弹窗

**B18 "Save" 按钮 (W-47 编辑模式) 验证**:
- 可见条件: 仅在 edit mode (点击 "Edit Domain Settings" 后进入编辑态)
- Disabled 条件: domain input 为空 OR 格式无效 OR 域名已被 claim (W-48)
- Loading: 点击后 spinner, disabled 防重复
- 验证: valid domain format (no protocol/path, 允许子域名); 格式错误 → red border + "Please enter a valid domain name"
- **Save 触发 D20 的原因**: 域名变更需要重新验证 DNS + SSL, 是生产环境变更, 需走 D20 确认流程
- Save 流程: Save → D20 弹窗 → D20 通过 → `PUT /api/white-label/domain` → 成功: toast + 退出编辑态 + 重新开始 DNS polling

**DNS 验证 → SSL 流程**:
1. DNS verified → 自动触发 SSL provisioning (Let's Encrypt)
2. SSL provisioning: Status checklist 显示 ⏳ "SSL: Provisioning..." (spinner)
3. SSL 通常 2-5 分钟; 页面通过 polling `GET /api/white-label/domain` 检查 ssl_status (30s 间隔)
4. SSL active → checklist ✅ "SSL: Active" + Portal Preview HTTPS 锁图标
5. SSL 失败: "SSL provisioning failed. Retrying..." + 自动重试最多 3 次 (间隔 60s)
6. 3 次重试均失败: 红色提示 "SSL provisioning failed. Please contact support." + Stop retry

**Brand & Portal 区域 (③)**:
- 显示当前 Brand 设置 (来自 `GET /api/white-label/brand`): Logo 缩略图 + Primary Color 色块 + Preview 链接
- **只读展示** (非编辑): 修改 brand 需跳转 B40; "Edit Brand →" 链接 → B40
- Portal Preview: mini C端预览 (200×300 iframe); 仅 SSL active 后显示; SSL pending 时显示占位 "Portal will be available after SSL is active"

**C端影响**: Custom domain verified + SSL active → community 在自定义域名可访问; 旧 `share.taskon.io` URL 仍有效 (不 redirect); 两个 URL 并存

**Breadcrumb**: "← Back to White Label" → B15 或 B16 (基于 §17.1 状态路由)

**首次访问 (无域名)**: domain input 空, 所有 checklist 项 ○ 未完成, SSL section 不显示

**SSL provisioning 失败后 max retries**: 显示 "SSL provisioning failed after 3 attempts. Contact support@taskon.xyz for assistance." (red text)

**侧栏**: WL 子菜单展开, "Overview" 高亮 active (B18 不在子菜单中, 从 Overview 进入)

---

## 6. Embed & Deployment

### 6.1 B19 — Deployment Settings

**设计稿**: Node `RgCVQ` | URL: `/white-label/deploy`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "Deployment Settings"
├── Current Config Banner (紫色提示)
├── "Choose Your Embed Method" title
├── Embed Method Cards (3, horizontal)
│   ├── "Iframe Embed" — fastest, zero code
│   │   └── What You Get: full portal, zero config, auto-updates
│   │   └── "Get Embed Code →"
│   ├── "Widget Library" ★ RECOMMENDED — 最灵活
│   │   └── What You Get: pick widgets, embed anywhere, CSS match
│   │   └── "Browse Widget Library →"
│   └── "Page Builder" — most powerful
│       └── What You Get: drag-drop, templates, multi-page
│       └── "Open Page Builder →"
├── Quick Comparison Table
│   ├── Setup Time: ~5 min / 10-30 min / 20-60 min
│   ├── Code: None / For widget / Full
│   ├── Layout Control: Fixed / Per widget / Full custom
│   └── Best For: Quick launch / Brand match / Template campaigns
├── Tip banner: "Not sure which to pick?..."
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Get Embed Code →" | → B42 Iframe Embed |
| "Browse Widget Library →" | → B20 Widget Library |
| "Open Page Builder →" | → B23 Page Builder |

#### 操作详情 (v2.0 新增, W-49 ~ W-52)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-49 | Neutral 状态 (B19v) | render | — | Node `Rwq2K`: 3 个 Embed Method 卡片 (Iframe / Widget Library ★RECOMMENDED / Page Builder); 每卡片: 标题 + "What You Get" 列表 (3-4 bullet points) + CTA 按钮; 底部 Quick Comparison Table (Setup Time / Code Required / Layout Control / Best For); Tip banner "Not sure which to pick? Start with Widget Library..." | — |
| W-50 | 选择 Iframe | click "Get Embed Code →" | — | → 跳转 B42 Iframe Embed 页面 (`/white-label/embed/iframe`) | — |
| W-51 | 选择 Widget Library | click "Browse Widget Library →" | `GET /api/white-label/widgets` | 检查已有 widget: 0 → 跳转 B20 (Empty); 1+ → 跳转 B22 (Active) | — |
| W-52 | 选择 Page Builder | click "Open Page Builder →" | `GET /api/white-label/pages` | 检查已有 page: 0 → 跳转 B23 (Empty); 1+ → 跳转 B25 (Active) | — |

**B19 页面初始化 (Init)**:
- API: `GET /api/white-label/status` — 获取当前部署状态 (用于 Current Config Banner + 方法高亮)
- 加载态: Banner 区域 skeleton; 3 个 Embed Method 卡片立即渲染 (静态内容)
- 路由守卫: 无特殊限制, 任何 WL 状态均可访问

**B19 页面离开 (Destroy)**: 无状态需清理

**B19 当前部署方法高亮**: 如已配置某方法 (如 Iframe), 对应卡片显示 "Currently Active" badge (绿色, 右上角); 其他卡片无 badge

**B19 W-51/W-52 API loading**: 点击 CTA 按钮后, 按钮 inline spinner; API 返回后跳转; API 失败 → toast error + 按钮恢复

**B19 Loading/Error**: Banner skeleton (1 行); 3 卡片静态渲染; Comparison Table 静态渲染

**侧栏**: WL 子菜单展开, "Overview" 高亮 active (B19 不在子菜单中)

---

### 6.2 B42 — Iframe Embed

**设计稿**: Node `ByGS0` | URL: `/white-label/embed/iframe`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "Iframe Embed" + "Preview" button
├── ── Embed Configuration ──
│   ├── Source URL: input (readonly, auto-generated)
│   ├── Display Mode: dropdown (Full Page / Sidebar / Modal)
│   ├── Width: input (default 100%)
│   ├── Height: input (default 800px)
├── ── Embed Code ──
│   ├── Code block (<iframe src=... />)
│   └── "Copy Code" button
├── ── SSO Configuration ──
│   ├── JWT Provider: dropdown
│   ├── Redirect URL: input
│   └── "Test Connection" button
```

#### 操作详情 (v2.0 新增, W-53 ~ W-57)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-53 | 完整页面结构 | render | `GET /api/white-label/embed/iframe` | Node `ByGS0`: Breadcrumb ("← Back to White Label") + Header ("Iframe Embed" + "Preview" button) + Embed Configuration (Source URL readonly + Display Mode dropdown + Width/Height inputs) + Embed Code (auto-generated `<iframe>` code block + "Copy Code" button) + SSO Configuration (JWT Provider dropdown + Redirect URL input + "Test Connection" button) | — |
| W-54 | Iframe URL | render | — | Source URL: `share.taskon.io/embed/{project_id}` (readonly, auto-generated); "Copy" 按钮 → clipboard → toast "URL Copied!"; URL 在 WL publish 后自动生效 | — |
| W-55 | Iframe Code 生成 | config change | — | 配置变化时自动更新代码: `<iframe src="share.taskon.io/embed/{project_id}" width="{width}" height="{height}" frameborder="0" allow="clipboard-write"></iframe>`; Display Mode 影响默认尺寸: Full Page (100%×800px) / Sidebar (360px×600px) / Modal (480px×640px); "Copy Code" → clipboard → toast | — |
| W-56 | SSO 配置 | form | — | JWT Provider: 选择后显示对应集成代码 snippet; "Wallet Auth" → 显示 `window.taskonEmbed.auth({ type: 'wallet' })` 代码; "OAuth2/JWT" → 显示 JWT token 传递代码 (`postMessage` API); Redirect URL: SSO 回调地址; 均有 Copy 按钮 | — |
| W-57 | "Test Embed" / "Preview" | click "Preview" button | — | 页面底部展开 inline iframe 预览区 (slide-down animation 300ms): 显示实际 iframe 渲染效果; 预览区含 "Close Preview" 按钮; SSO 未配置时预览显示只读模式 (灰色遮罩 + "SSO not configured" 提示) | iframe 加载失败 → 显示 "Unable to load preview. Check if your domain is configured correctly." |

**B42 页面初始化 (Init)**:
- API: `GET /api/white-label/embed/iframe` — 返回: `{ sourceUrl, displayMode, width, height, ssoConfig: { provider, redirectUrl }, embedCode }`
- 加载态: 全页面 skeleton (表单 + 代码区)
- API 失败: §2.3.4 错误态 + Retry
- 无配置 (首次访问): sourceUrl 自动生成; displayMode 默认 "Full Page"; width/height 使用默认值

**B42 页面离开 (Destroy)**:
- 未保存的 SSO 配置: 显示 "Unsaved changes" 确认弹窗 (通过 §2.3.5 脏数据检测)

**B42 Display Mode 变更**:
- **自动保存**: Display Mode dropdown 变更后自动调用 `PUT /api/white-label/embed/iframe` 保存 (debounce 1s)
- 自动保存期间: dropdown 旁显示 "Saving..." 小文字; 成功: "Saved ✓" (2s 后消失); 失败: toast error
- Display Mode 变更同时更新 Width/Height 默认值: Full Page (100% × 800px) / Sidebar (360px × 600px) / Modal (480px × 640px)

**B42 Width / Height 验证**:
- Width: 数字或百分比; 数字: min 200, max 3840 (px); 百分比: "100%" 格式; 无效 → red border + "Enter a valid width (e.g., 800 or 100%)"
- Height: 仅数字 (px); min 300, max 3840; 无效 → red border + "Enter a valid height (300-3840px)"
- 验证时机: blur 时校验; 自动保存: 验证通过后 debounce 1s 自动调 PUT API

**B42 "Test Connection" (SSO section)**:
- API: `POST /api/white-label/embed/iframe/test-sso`
- 按钮 spinner → 测试 SSO 配置: JWT provider 可达 + redirect URL 有效 → 成功: "✓ SSO connection OK" (green, 5s 后消失) / 失败: "✗ SSO test failed: {reason}" (red)
- 仅在 SSO 配置有内容时 enabled; 空 SSO 时 disabled + tooltip "Configure SSO first"

**B42 Breadcrumb**: "← Back to White Label" → B19 Deployment Settings (直接上级页面)

**B42 Loading state**: 全页面 skeleton (表单 + 代码块 + SSO section)

**B42 Iframe preview fail**: Preview 区域显示 fallback 文案 (W-57), 不影响配置保存

**B42 C端影响**: 配置保存后, 已嵌入的 iframe **自动生效** (下次 iframe 加载时读取新配置); 无需客户端重新部署

**侧栏**: WL 子菜单展开, "Widgets" 高亮 active (B42 归属 Embed 方法, 但侧栏无专项; 使用 "Overview" active)

---

## 7. Widget Library（3 状态）

### 7.1 B20 — Widget Library Empty

**设计稿**: Node `2sSsA` | URL: `/white-label/widgets`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to Embed Options"
├── Header: "Widget Library"
├── "Community Modules" label + count badge
├── Module Cards Grid (2 rows × 3)
│   Row 1 (Configured — green border):
│   ├── Leaderboard — stats preview + "Add Widget →" (green)
│   ├── Task List — stats preview + "Add Widget →" (green)
│   └── User Center — preview + "Add Widget →" (green)
│   Row 2 (Configured — green border):
│   ├── Rewards Shop — stats + "Add Widget →" (green)
│   ├── Daily Check-in — preview + "Add Widget →" (green)
│   └── Quest Card — preview + "Add Widget →" (green)
├── ── Not Yet Configured (amber border) ──
│   └── Milestones — "Set Up in Community →" (amber button)
├── Ready Tip: "Newly vs Needs Setup..."
```

**Widget 两种类型**:
- **Configured** (绿色): Community 中已配置 → "Add Widget →" 创建 widget
- **Not Yet Configured** (amber): Community 中未启用 → "Set Up in Community →" 跳转 Community

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "Add Widget →" | → B21 Widget Config |
| "Set Up in Community →" | → B10/B11 Community Hub |
| "Create Your First Widget" CTA | → B21 (无 moduleType 预选) |

**B20 页面初始化 (Init)**:
- API: `GET /api/community/modules/status` — 获取 Community 各模块配置状态
- 路由守卫: `GET /api/white-label/widgets` → widgets.length > 0 → redirect B22 (Active)
- 加载态: Module Cards Grid skeleton (6 卡片占位)
- API 失败: §2.3.4 错误态 + Retry
- Community 无已配置模块 (全部 amber): 无 "Add Widget →" 绿色卡片; "Create Your First Widget" CTA disabled + tooltip "Set up at least 1 Community module first"; 特殊提示: "All modules need setup. Go to Community to configure your first module."

**B20 页面离开 (Destroy)**: 无状态需清理

**B20 "Add Widget →"**: 跳转 B21 + URL param `?moduleType=leaderboard` (预填 moduleType)

**B20 Loading skeleton**: Module Cards Grid (6 卡片 shimmer)

**B20 Breadcrumb**: "← Back to Embed Options" → B19

**侧栏**: WL 子菜单展开, "Widgets" 高亮 active

---

### 7.2 B21 — Widget Config

**设计稿**: Node `n4pJK` | URL: `/white-label/widgets/:id/config`

#### 页面结构 (2-column)

```
Left Panel: Widget Settings
├── Widget Name: text input
├── Theme: Light / Dark toggle
├── Show Top N Users: number input (default 10)
├── Refresh Interval: dropdown (Every 5 minutes)
├── Display Options (checkboxes):
│   ├── ☑ Show point badges
│   ├── ☑ Show rank change indicators
│   └── ☑ Show avatar address
├── "Save & Get Embed Code" button (紫色)

Right Panel: Live Preview + Embed Code
├── Live Preview card (实时更新)
│   └── Leaderboard widget mock (4 entries)
├── Embed Code block
│   ├── <taskon-leaderboard ... /> code
│   └── "Copy" button
├── Tip: "After configuring, paste the embed code..."
```

#### Widget Config 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | Widget 名称 |
| moduleType | enum | leaderboard / tasks / user-center / shop / daychain / quest |
| theme | enum | `light` / `dark` |
| showTopN | number | 显示条目数 |
| refreshInterval | number | 刷新间隔(秒) |
| displayOptions | object | 显示选项 flags |

**B21 页面初始化 (Init)**:
- **Create mode** (`/white-label/widgets/new?moduleType=leaderboard`): 空表单, moduleType 从 URL param 预填; 如无 param → moduleType select 为空, 需用户选择
- **Edit mode** (`/white-label/widgets/:id/config`): `GET /api/white-label/widgets/:id` → 预填所有字段
- 加载态: edit mode 表单 skeleton; create mode 无需 loading (空表单)
- Edit mode 404: toast "Widget not found" + redirect B22

**B21 页面离开 (Destroy)**:
- 未保存配置: "Unsaved changes" 确认弹窗 (§2.3.5 脏数据检测)

**B21 "Save & Get Embed Code" 条件**:
- Enabled 条件: Widget Name (非空, 1-50 chars) + Module Type (已选) + 所有必填字段有值
- Disabled 时: opacity 0.5 + tooltip "Fill in required fields to save"
- Loading: button spinner + disabled 防重复
- API: Create → `POST /api/white-label/widgets`; Edit → `PUT /api/white-label/widgets/:id`
- **成功后**: 留在 B21 (不跳转); 右侧 Embed Code 区域更新显示生成的 embed code; toast "Widget saved"
- **Create mode 成功后**: URL 变为 `/white-label/widgets/:newId/config` (edit mode); 按钮文案变为 "Save Changes"

**B21 Style Config (补充到页面结构)**:
- 位于 Display Options checkboxes 下方, "Save" 按钮上方
- 字段:
  - Primary Color: color picker (默认 inherit from Brand, 可自定义 hex)
  - Border Radius: slider 0-24px (默认 8px)
  - Padding: select 8/12/16/24px (默认 16px)
- 这些字段在 W-58 已列出, 此处明确其在页面结构中的位置

**B21 Module Type 可修改性**: Create mode → 可选; Edit mode → **只读** (已创建的 widget 不可更改 moduleType, 需删除重建); 只读时 select 显示为 text + lock icon

**B21 Live Preview 更新**: 任何表单字段变化 → debounce 200ms → 右侧 Live Preview 更新; 更新字段: name, theme, showTopN, displayOptions, style config

**B21 Breadcrumb**: "← Back to Widget Library" → B20 或 B22 (根据 widgets 数量)

**B21 Save 失败**: toast error + 表单数据保留

**B21 Edit mode loading state**: 表单 skeleton (6 字段占位) + preview skeleton

**侧栏**: WL 子菜单展开, "Widgets" 高亮 active

---

### 7.3 B22 — Widget Library Active

**设计稿**: Node `S432k` | URL: `/white-label/widgets`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to Embed Options"
├── Header: "Widget Library"
├── "My Widgets" label + count badge
├── My Widgets Row (已创建的 widgets)
│   ├── "Main Leaderboard" card [Active]
│   │   ├── Stats: 1,247 views / 342 interactions / 27.4%
│   │   ├── Embed URL: share.myproject.io/widgets/...
│   │   └── Actions: "Configure" / "Analytics" / "Copy" embed
│   └── "Onboarding Tasks" card [Active]
│       └── Similar structure
├── "Community Modules" label
├── Available Module Cards (remaining not-yet-widgetized)
│   └── Each: preview + "Add Widget →"
├── ── Not Yet Configured ──
│   └── Milestones: "Set Up in Community →"
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Widget "Configure" | → B21 |
| Widget "Analytics" | → B43 Page Analytics |
| Widget "Copy" embed | (action) 复制 embed code |
| "Add Widget →" | → B21 (new) |
| "Deploy Widget" | → 触发 D20 |

#### Widget Library 操作详情 (v2.0 新增, W-58 ~ W-62)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-58 | B21 Widget Config 完整字段 | form | `GET /api/white-label/widgets/:id` | 左侧: Widget Name (text, 必填 1-50) + Module Type (select: Leaderboard/Tasks/User Center/Shop/DayChain/Quest, 必填) + Theme (Light/Dark toggle) + Show Top N (number, 默认 10) + Refresh Interval (select: 5/15/30/60 min) + Display Options (checkbox: point badges / rank change / avatar address) + Style Config (Primary Color: inherit from brand / custom hex; Border Radius: 0-24px slider; Padding: 8/12/16/24px select); 右侧: Live Preview (实时更新) | — |
| W-59 | B21 Embed Code 生成 | save widget | `POST /api/white-label/widgets` | "Save & Get Embed Code" → API call → 成功: 右侧 Embed Code 区显示 `<script src="cdn.taskon.io/widget.js"></script><taskon-{module} project="{id}" widget="{widget_id}" theme="{theme}"></taskon-{module}>`; "Copy" 按钮 → clipboard → toast "Embed code copied!"; Tip: "Paste this code into your website's HTML" | 保存失败 → toast error |
| W-60 | B22 Active 操作 | click per widget | various | 每个 widget 卡片操作: **Configure** → B21 (edit mode, 预填数据); **Copy Embed** → clipboard + toast; **Deploy** → 触发 D20 Publish Readiness Check; **Delete** → confirm dialog "Delete widget '{name}'? This will remove it from all embedded locations." → 确认 → `DELETE /api/white-label/widgets/:id` → 列表刷新 + toast "Widget deleted" | 删除被 Rule 引用 → 403 + toast "Cannot delete: widget is used in active deployments" |
| W-61 | Widget 状态 badge | render | — | 每个 widget 卡片右上角状态: **Configured** (绿色 `#16A34A` bg `#0A2E1A`): 已配置未部署; **Deployed** (蓝色 `#3B82F6` bg `#0F1A2E`): 已部署到生产; **Draft** (灰色 `#64748B` bg `#1E293B`): 配置未完成; 状态来源: API response `.status` 字段 | — |
| W-62 | B20 Template cards | click template | — | Empty 状态下 "Community Modules" 网格: 每个 module card = icon + name + metrics (if configured) + CTA; Configured modules (green border): "Add Widget →" → B21 (预填 moduleType); Not Yet Configured (amber border): "Set Up in Community →" → B10/B11; 模板数据来源: `GET /api/community/modules/status` 获取各模块配置状态 | — |

**B22 页面初始化 (Init)**:
- API: `GET /api/white-label/widgets` — 获取所有已创建 widgets 列表 (含 status, stats)
- 路由守卫: widgets.length === 0 → redirect B20 (Empty)
- 加载态: My Widgets 区域 skeleton (卡片占位) + Community Modules skeleton
- API 失败: §2.3.4 错误态 + Retry

**B22 页面离开 (Destroy)**: 无状态需清理

**B22 "Deploy Widget" → D20 → 后续**:
- D20 通过后: `PUT /api/white-label/widgets/:id` + `{ status: 'deployed' }` → widget 状态变 "Deployed" (蓝色 badge)
- Embed code 激活 (C端可通过 embed code 访问该 widget)
- B22 列表自动刷新 (refetch API)
- toast "Widget deployed successfully"

**B22 "Delete" widget 澄清**:
- W-60 中 "Rule 引用 → 403" 指的是: widget 被 Page Builder 的某个 page 引用时, 删除会影响该 page
- 具体逻辑: `DELETE /api/white-label/widgets/:id` → 后端检查:
  - 被 Page 引用 → 仍可删除, 但 confirm dialog 额外警告 "This widget is used in {n} pages. Those pages will show a placeholder."
  - 被 active deployment 引用 → 403 "Cannot delete: widget is currently deployed. Undeploy first."
- 删除成功: 列表刷新; 如最后一个 widget 被删除 → redirect B20 (Empty)

**B22 Widget 卡片交互**:
- 卡片点击区域: **整张卡片** 可点击 → B21 (edit mode); Actions 按钮 (Configure/Analytics/Copy/Deploy/Delete) 阻止冒泡
- Widget deployed 后: 卡片增加 embed code 显示 + "Copy" 按钮

**B22 "Community Modules" section 更新**: 删除 widget 后, 对应 module 重新出现在 "Community Modules" 区域 (Available → "Add Widget →")

**B22 Loading skeleton**: My Widgets (卡片 shimmer) + Community Modules (网格 shimmer)

**B22 Widget list API 失败**: skeleton → 错误态 "Unable to load widgets" + Retry

**侧栏**: WL 子菜单展开, "Widgets" 高亮 active

---

## 8. Page Builder（3 状态）

### 8.1 B23 — Page Builder Empty

**设计稿**: Node `DRYwN` | URL: `/white-label/pages`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to Embed Options"
├── Header: "Page Builder"
├── "How It Works" (3 steps, horizontal cards)
│   ├── ① Create Widgets — Build individual widgets in Widget Library
│   ├── ② Design Your Page — Arrange widgets in a custom layout
│   └── ③ Publish & Embed — Get URL or embed code
├── "Page Templates" label
├── Template Cards (3)
│   ├── "Rewards Hub" — leaderboard + rewards + shop
│   ├── "Community Portal" — tasks + leaderboard + news
│   └── "Custom Page" — start from empty canvas
├── Callout: widgets needed → Widget Library
├── "Create Your First Page" CTA (紫色)
├── Tip: "Pages are fully customizable..."
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| Template cards | → B24 (pre-filled) |
| Widget callout | → B20 Widget Library |
| "Create Your First Page" | → B24 (空白页, 无模板) |

**B23 页面初始化 (Init)**:
- 路由守卫: `GET /api/white-label/pages` → pages.length > 0 → redirect B25 (Active)
- Widget 可用性检查: `GET /api/white-label/widgets?status=configured` → 用于 template cards 的 widget 依赖提示
- 加载态: Template Cards skeleton (3 卡片占位)
- API 失败: §2.3.4 错误态 + Retry

**B23 页面离开 (Destroy)**: 无状态需清理

**B23 Template 卡片**:
- 每个 template 卡片始终 enabled (点击后进入 B24, B24 内检查 widget 可用性)
- 卡片显示 template 使用的 widget 类型; 如无对应 widget 配置 → 卡片底部 amber 提示 "Requires: Leaderboard, Tasks widgets"

**B23 "Create Your First Page"**: 跳转 B24 空白页 (无模板预填)

**B23 Loading state**: Template Cards skeleton

**侧栏**: WL 子菜单展开, "Pages" 高亮 active

---

### 8.2 B24 — Page Builder Editor

**设计稿**: Node `sGDcq` | URL: `/white-label/pages/new`

#### 页面结构 (2-column)

```
Left Panel: Canvas
├── URL bar mockup: https://share.taskon.io/pages/rewards-hub
├── Page Content
│   ├── Widget Block 1: "Leaderboard Widget" (expandable)
│   │   └── Mock leaderboard content
│   ├── Widget Block 2: "Onboarding Widget" (expandable)
│   │   └── Mock task list
│   └── "+ Add Widget Block" button

Right Panel: Page Settings
├── Page Name: text input
├── URL Slug: text input (auto-generated)
├── Theme: Light / Dark toggle
├── "Widgets on Page" list (draggable order)
│   ├── • Leaderboard (↕ drag)
│   └── • Task List (↕ drag)
├── "Available Widgets" list
│   ├── User Center — "+ Add"
│   ├── Rewards Shop — "+ Add"
│   └── Daily Check-in — "+ Add"
├── "Don't see one? Go to Widget Library →"
├── "Publish Page" button (紫色) → D20
├── "Save Draft" button (secondary)
```

#### 交互逻辑

1. **Widget Block 拖拽**: Canvas 内 widget block 可拖拽重排
2. **"+Add Widget Block"**: 从 Available Widgets 添加到 Page
3. **"× Remove"**: 每个 widget block 可移除
4. **Live Preview**: Canvas 实时反映设置变更
5. **Theme Toggle**: 整页配色切换 Light/Dark

**B24 页面初始化 (Init)**:
- **Create mode** (`/white-label/pages/new` 或 `/white-label/pages/new?template=rewards-hub`):
  - 无 template: 空 canvas + 空设置
  - 有 template: `GET /api/white-label/pages/templates/:templateId` → 预填 page name, widgets, layout
  - Available Widgets: `GET /api/white-label/widgets?status=configured` → 右侧 "Available Widgets" 列表
- **Edit mode** (`/white-label/pages/:id/edit`):
  - `GET /api/white-label/pages/:id` → 预填所有字段 (page name, slug, theme, widgets, order)
  - 加载态: Canvas skeleton + Settings skeleton
  - 404: toast "Page not found" + redirect B25
- 加载态: Create with template → canvas skeleton; Edit → full page skeleton

**B24 页面离开 (Destroy)**:
- 未保存的编辑: "Unsaved changes" 确认弹窗 (§2.3.5)
- 无自动保存 (需手动 Save Draft 或 Publish)

**B24 "Publish Page" 条件**:
- Enabled 条件: Page Name 非空 (1-60 chars) + 至少 1 个 widget block + Slug 验证通过 (唯一)
- Disabled 时: opacity 0.5 + tooltip "Add a page name and at least 1 widget"
- 流程: Publish → D20 弹窗 → D20 通过 → `POST /api/white-label/pages` (create) 或 `PUT /api/white-label/pages/:id` (update) with `{ status: 'published' }` → 成功: toast "Page published!" + **redirect B25** (Active pages list)
- Edit mode: 按钮文案 "Update Page" (仍触发 D20)

**B24 "Save Draft" API**:
- Create: `POST /api/white-label/pages` with `{ status: 'draft', ... }` → 返回 `{ id }` → URL 变为 `/white-label/pages/:id/edit`
- Edit: `PUT /api/white-label/pages/:id` with `{ status: 'draft', ... }`
- 成功: toast "Draft saved" + **留在 B24** (不跳转)
- 失败: toast error + 数据保留

**B24 Page Name 验证**:
- 必填: 1-60 chars
- 无效: 空 → red border + "Page name is required"; > 60 chars → "Maximum 60 characters"
- 验证时机: blur 时

**B24 Slug 验证**: 见 W-66 (debounce 500ms 唯一性检查)
- Slug 冲突: red border + "This slug is already taken. Suggestion: {slug}-2"

**B24 "+ Add Widget Block" (W-64) 边界**:
- 0 available + 0 configured widgets: dropdown 内显示 "No widgets available. Go to Widget Library →" + 链接到 B20

**B24 Post-publish C端影响**:
- Page URL `https://share.taskon.io/pages/{slug}` 立即可访问
- Page 出现在 Domain Portal (如有)
- B43 Analytics 开始收集数据

**B24 Canvas-Settings 双向同步 (W-68)**: 实时同步, 无 debounce; Canvas 拖拽更新列表, 列表拖拽更新 Canvas; 同步通过共享状态 (React state) 实现

**B24 退出**: 无 back 按钮; 通过 Breadcrumb "← Back to Pages" → B25 或 B23; 或 sidebar "Pages"

**B24 Editor loading state (edit mode)**: Canvas 全区域 skeleton + 右侧 Settings skeleton

**侧栏**: WL 子菜单展开, "Pages" 高亮 active

---

### 8.3 B25 — Page Builder Active

**设计稿**: Node `J08v5` | URL: `/white-label/pages`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to Embed Options"
├── Header: "Page Builder" + "+ Create New Page" button
├── "My Pages" label + count
├── Page Card: "Rewards Hub" [Published]
│   ├── Stats: 1,247 views / 342 interactions / 27.6%
│   ├── URL: https://widgets/leaderboard-task-us/rewards-hub
│   ├── Embed code display + "Copy" button
│   └── Actions: "Edit Page" / "Analytics"
├── "Quick Start a New Page" label
├── Template Cards Row (4: Rewards Hub / Community Portal / Contest Page / Quest Landing)
```

#### 按钮路由

| 按钮 | 目标 |
|------|------|
| "+ Create New Page" | → B24 |
| Page "Edit Page" | → B24 (edit existing) |
| Page "Analytics" | → B43 |
| Page "Copy" embed | (action) |
| Template cards | → B24 (pre-filled) |

**B25 页面初始化 (Init)**:
- API: `GET /api/white-label/pages` — 获取所有 pages 列表 (含 status, stats summary)
- 路由守卫: pages.length === 0 → redirect B23 (Empty)
- 加载态: My Pages 区域 skeleton (卡片占位) + Template Cards skeleton
- API 失败: §2.3.4 错误态 + Retry

**B25 页面离开 (Destroy)**: 无状态需清理

**B25 Page 卡片交互**:
- 卡片点击: **整张卡片** 可点击 → B24 (edit mode); Action 按钮阻止冒泡
- Page embed code "Copy": 仅 Published 状态的 pages 显示 embed code + Copy 按钮; Draft 页面无 embed code

**B25 "+ Create New Page"**: 跳转 B24 (create mode, 空白); plan 限额检查: 达到限额 → disabled + tooltip "Upgrade plan for more pages" (见 §D.1)

**B25 Page status badges 样式**:
- Published: 绿色 bg `#0A2E1A` + text `#16A34A`
- Draft: amber bg `#1F1A08` + text `#D97706`
- Unpublished: 灰色 bg `#1E293B` + text `#64748B` (曾 published 后 unpublished)

**B25 Stats 刷新 (W-70)**: 从 B24 返回时 stale-while-revalidate (先显示缓存, 后台 refetch 更新)

**B25 Loading skeleton**: My Pages (卡片 shimmer) + Template Cards (4 卡片 shimmer)

**侧栏**: WL 子菜单展开, "Pages" 高亮 active

#### Page Builder 操作详情 (v2.0 新增, W-63 ~ W-70)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-63 | B24 Canvas 拖拽 | drag widget block | — | 拖拽手柄: 左侧 `⋮⋮` (drag_indicator icon); 拖拽中: 被拖 block opacity 0.5 + 蓝色 border 2px `#3B82F6`; 其他 block 之间出现占位符 (4px 蓝色虚线); 放下: 300ms ease transition + 右侧 "Widgets on Page" 列表同步更新排序; 拖拽范围限制在 Canvas 内 | — |
| W-64 | B24 "+ Add Widget Block" | click | `GET /api/white-label/widgets?status=configured` | 按钮 → 展开 dropdown (slide-down 200ms): 列出 "Available Widgets" (来自 B22 已配置的 widget 列表); 每项: widget name + module type icon + "Add" 按钮; 已在 page 上的 widget 显示 "Already added" (灰色, 不可点); 点击 "Add" → widget block 追加到 Canvas 底部 + dropdown 关闭 | 无可用 widget → dropdown 内显示 "No widgets configured. Go to Widget Library →" |
| W-65 | B24 Widget Block "× Remove" | click × icon | — | 每个 widget block 右上角 × 按钮; hover → ×变红; click → confirm popover "Remove this widget from page?" [Cancel] [Remove]; 确认 → block 300ms slide-up 移除 + 右侧 "Widgets on Page" 列表同步更新; **不会删除 widget 本身** (仅从 page 移除) | — |
| W-66 | B24 URL Slug | input | `GET /api/white-label/pages/check-slug?slug={value}` | 基于 Page Name 自动生成 (toLowerCase + replace spaces with hyphens + remove special chars); 可手动编辑; 输入时 debounce 500ms → API 检查唯一性; 冲突 → 红色 border + "This slug is already taken. Suggestion: {slug}-2"; 完整 URL preview: `https://share.taskon.io/pages/{slug}` | — |
| W-67 | B24 Theme Toggle | click Light/Dark | — | Canvas preview 即时切换: Light → 白色 bg + dark text; Dark → #0A0F1A bg + light text; Widget blocks 内部配色跟随; 300ms transition; Toggle 样式: 当前选中项 filled purple, 另一项 outline | — |
| W-68 | B24 Settings Widget 排序 | drag in list | — | 右侧 "Widgets on Page" 列表: 每项有 ⋮⋮ drag handle; 拖拽重排 → **Canvas 同步更新** widget block 顺序 (300ms transition); 列表与 Canvas 双向同步 (Canvas 拖拽也更新列表, 列表拖拽也更新 Canvas) | — |
| W-69 | B25 "Edit Page" | click | `GET /api/white-label/pages/:id` | 加载已有 Page 数据 (page name, slug, theme, widgets, widget order) → 跳转 B24 (edit mode); B24 top bar 显示 "Edit: {page_name}"; 所有字段预填已有数据; "Publish Page" 按钮变为 "Update Page" (仍触发 D20) | 404 → toast "Page not found" + redirect B25 |
| W-70 | B25 Page stats | render | `GET /api/white-label/pages/:id/analytics?summary=true` | 每个 Page 卡片: Page Views (total) / Unique Visitors (30d) / Widget Clicks (30d) / Completion Rate (%); 数据从 B43 Page Analytics 同一 API 取摘要; "Analytics" 按钮 → B43 (完整分析页) | API 失败 → stats 显示 "—" |

---

## 9. Brand Settings

### 9.1 B40 — Brand Settings

**设计稿**: Node `Cx3LH` | URL: `/white-label/brand`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "Brand Settings" + "Save Changes" button
├── ── Logo ──
│   ├── Current: project-logo.svg (preview)
│   ├── Specs: SVG/PNG, min 256×256, max 2MB
│   └── "Change" button
├── ── Brand Colors ──
│   ├── Primary: color swatch + hex input (#48E5E1)
│   ├── Secondary: color swatch + hex (#7C3AED)
│   ├── Accent: color swatch
│   ├── Background: color swatch
│   └── Additional presets
├── ── Typography ──
│   ├── Heading Font: dropdown (Inter)
│   ├── Body Font: dropdown (Inter)
│   ├── Heading Preview: sample text
├── ── Custom CSS ──
│   ├── Code editor (monospace)
│   └── "Preview" link
├── Footer: "Changes will take effect on next deployment..."
```

#### API

| Endpoint | Method | 说明 |
|----------|--------|------|
| `/api/white-label/brand` | GET | 获取当前品牌设置 |
| `/api/white-label/brand` | PUT | 保存品牌设置 |

#### 操作详情 (v2.0 新增, W-71 ~ W-77)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-71 | Logo 上传 | click "Change" | `POST /api/white-label/brand/logo` (multipart) | "Change" → file picker (accept: .svg, .png); 选择文件 → 客户端验证 (格式+尺寸+大小) → 上传 spinner → 预览区域立即更新; Specs: SVG/PNG, min 256×256px, max 2MB; 当前 logo 显示为 64×64 缩略图 + 文件名 | 见 W-72 |
| W-72 | Logo 上传错误 | validation / upload | — | 格式错误: "Unsupported format. Please upload SVG or PNG." (red text); 尺寸过小: "Image too small. Minimum 256×256 pixels required." (red text); 体积过大: "File too large. Maximum 2MB allowed." (red text); 上传失败: toast "Upload failed. Please try again." + 保留原 logo | — |
| W-73 | Color Picker | click color swatch | — | 点击色块 → 展开 color picker popover: hex input (#000000 format) + hue slider (0-360°) + saturation/brightness 面板; 变化即时反映在: ① 色块 preview ② 右侧 preview 区 (如存在); debounce 100ms; 点击 popover 外部关闭; 支持粘贴 hex 值 (自动补 #) | 无效 hex → input 红色 border + 不更新预览 |
| W-74 | Font Dropdown | select | — | 约 20 种 Google Fonts: Inter (默认), Roboto, Open Sans, Lato, Montserrat, Poppins, Nunito, Raleway, Ubuntu, Playfair Display, Merriweather, Source Code Pro, Fira Code, Space Grotesk, DM Sans, Plus Jakarta Sans, Outfit, Manrope, Lexend, Sora; 选择后 → 右侧 Heading/Body Preview 文本字体即时变化 (font-family swap) | — |
| W-75 | Custom CSS Editor | code input | — | Monospace code editor area (min-height 120px); 基础语法高亮: 单色关键字 (properties: purple, values: green, selectors: blue); "Preview" link → 右侧 preview panel 实时应用 CSS (iframe sandbox, 防 XSS); CSS 限制: max 10KB; 仅支持 .taskon-widget 命名空间下的选择器 | CSS 语法错误 → editor 下方 amber warning "CSS may contain syntax errors" |
| W-76 | "Save Changes" | click | `PUT /api/white-label/brand` | button spinner → 成功: toast "Brand settings updated" + button 恢复; **级联效果**: 已部署的 Widget/Page 在**下次用户访问时**自动加载新品牌设置 (CDN 缓存 TTL 5min → 最迟 5 分钟生效); 不需要手动 re-deploy | 保存失败 → toast error "Failed to save brand settings" |
| W-77 | 已部署时品牌变更 | save success | — | **无需重新发布**: Brand 设置为全局配置, 通过 API 动态加载; 已部署 Widget/Page 在下次请求 `GET /api/white-label/brand` 时获取最新值; 页脚提示: "Changes will take effect on next deployment or within 5 minutes for live widgets."; **不会**弹出 re-deploy 确认 | — |

**B40 页面初始化 (Init)**:
- API: `GET /api/white-label/brand` — 加载当前品牌设置 (logo / colors / fonts / css)
- 加载态: 各 section (Logo / Colors / Typography / CSS) 独立 skeleton
- API 失败: §2.3.4 错误态 + Retry
- 首次访问 (未配置): 所有字段使用 TaskOn 默认值 (logo: 无, primaryColor: #48E5E1, secondaryColor: #7C3AED, headingFont: Inter, bodyFont: Inter); 页面提示 "Customize your brand to match your project identity"

**B40 页面离开 (Destroy)**:
- 未保存修改: "Unsaved changes" 确认弹窗 (§2.3.5 脏数据检测)

**B40 "Save Changes" 按钮条件**:
- Disabled 条件: 表单未修改 (clean) → disabled (opacity 0.5 + tooltip "No changes to save")
- Enabled: 表单有修改 (dirty)
- Loading: button spinner + disabled 防重复

**B40 Preview panel (页面结构补充)**:
- 位于页面右侧 (2-column layout: 左 form 60% / 右 preview 40%)
- Preview 显示: header (logo + project name) + heading (font preview) + button (style preview) + widget mock (颜色主题)
- 实时更新: 表单变化 → debounce 200ms → preview 更新
- Custom CSS: preview 通过 iframe sandbox 隔离渲染 CSS (防 XSS); CSS 变化 → preview 即时反映

**B40 权限 (§D.2)**: Admin / Editor 可编辑; Member → 页面只读, "Save Changes" 隐藏, 表单字段显示为 readonly text

**B40 Loading state**: Logo section skeleton + Colors section skeleton + Typography skeleton + CSS editor skeleton + Preview skeleton

**B40 Custom CSS preview panel**: W-75 中的 "Preview" link 点击 → 右侧 preview panel 实时应用 CSS; 如 preview panel 已显示 → CSS 变化即时生效

**B40 Breadcrumb**: "← Back to White Label" → B15 或 B16 (基于 §17.1)

**侧栏**: WL 子菜单展开, "Overview" 高亮 active (B40 不在子菜单中)

---

## 10. SDK & API

### 10.1 B41 — SDK & API

**设计稿**: Node `lQxT5` | URL: `/white-label/sdk`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "SDK & API" + "View Docs" button (ext)
├── ── API Keys ──
│   ├── Production Key: pk_live_... + "Copy" + "Regenerate"
│   ├── Test Key: pk_test_... + "Copy"
│   └── "+ Generate Key" button
├── ── Quick Start ──
│   ├── Code block: npm install + init snippet
│   └── "Copy Snippet" button
├── ── Webhooks ──
│   ├── Endpoint URL: input + "+ Add Endpoint"
│   ├── Active webhooks list
│   │   ├── https://... → task.completed, points.earned — Active
│   │   └── https://... → user.joined — Active
├── ── API Usage & Limits ──
│   ├── Requests: 12,450 / 50,000 this month
│   ├── Rate Limit: 100 req/min
│   └── Uptime: 99.8%
```

#### 交互逻辑

1. **Copy API Key**: 一键复制，显示 toast "Copied!"
2. **Regenerate Key**: 确认弹窗 → (API) 重新生成
3. **Webhook 管理**: 添加/编辑/删除 webhook endpoints
4. **Usage 显示**: 进度条可视化

#### 操作详情 (v2.0 新增, W-78 ~ W-83)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-78 | "Regenerate Key" | click | `POST /api/white-label/sdk/keys/regenerate` | confirm dialog: "⚠️ Regenerate API Key?" body: "This will **immediately invalidate** your current key. All integrations using this key will stop working until updated." buttons: [Cancel] [Regenerate] (red); 确认 → spinner → 成功: 新 key 显示 (pk_live_xxx) + toast "New API key generated. Update your integrations." + 旧 key 立即失效 | 失败 → toast error; 旧 key 保持有效 |
| W-79 | Webhook "+ Add Endpoint" | click | `POST /api/white-label/sdk/webhooks` | 按钮 → inline 展开表单 (slide-down): URL input (placeholder "https://api.yourproject.io/webhook") + Events 多选 checkbox (见 W-81) + [Cancel] [Save Webhook]; Save → API → 成功: 新 webhook 追加到列表 + toast "Webhook added" + 表单收起 | URL 格式无效 → input 红色 border + "Please enter a valid HTTPS URL"; 保存失败 → toast error |
| W-80 | Webhook "Test" | click per webhook | `POST /api/white-label/sdk/webhooks/:id/test` | 发送测试事件 (`{ event: "test.ping", timestamp: "..." }`) → button spinner → 显示结果: "✓ 200 OK (145ms)" (green) 或 "✗ 500 Internal Server Error (2,340ms)" (red); 结果显示 5s 后自动消失 | 超时 (10s) → "✗ Timeout — endpoint did not respond within 10 seconds" |
| W-81 | Webhook Events | checkbox list | — | 可选事件: ☑ `task.completed` / ☑ `points.earned` / ☑ `user.joined` / ☐ `level.up` / ☐ `badge.earned` / ☐ `sprint.ended` / ☐ `milestone.reached` / ☐ `privilege.granted`; 至少选 1 个; 已有 webhook 编辑时保留已选 | — |
| W-82 | API Usage 进度条 | hover | — | 进度条 (高度 8px, 圆角 4px): 绿色 (#16A34A) 填充 + 灰色 (#1E293B) 背景; hover tooltip: "12,450 of 50,000 requests used (24.9%) — Resets on {month_end_date}"; 进度 >80%: 颜色变 amber (#F59E0B); >95%: 颜色变 red (#EF4444) | — |
| W-83 | API 用量警告 | render (>80%) | — | 用量 >80% 时: 页面顶部 amber 警告条 (bg #1F1A08, border #F59E0B): "⚠️ You've used {percent}% of your monthly API quota ({used}/{total} requests). Consider upgrading your plan." + "Upgrade Plan →" 链接 → M07 Pricing page; >95% 时文案变红: "🚨 API quota nearly exhausted. Service may be throttled." | — |

**B41 页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/white-label/sdk` — API keys (production + test) + usage stats + rate limit info
  2. `GET /api/white-label/sdk/webhooks` — 已配置 webhooks 列表
- 加载态: 各 section (API Keys / Quick Start / Webhooks / Usage) 独立 skeleton
- API 失败: 各 section 独立 §2.3.4 错误态 + Retry
- 返回数据: `{ keys: [{ id, type: 'production'|'test', key: 'pk_live_...', createdAt }], usage: { used, total, rateLimit, uptime }, webhooks: [...] }`

**B41 页面离开 (Destroy)**: 无状态需清理 (无表单编辑态)

**B41 "+ Generate Key" vs "Regenerate" 区分**:
- **"+ Generate Key"**: 创建新 key pair (如无 key 或需要额外 test key); API: `POST /api/white-label/sdk/keys` → 返回新 key; 场景: 项目初次使用 SDK, 或需要单独的 test key
- **"Regenerate"**: 替换已有 key (旧 key 立即失效); API: `POST /api/white-label/sdk/keys/regenerate` (W-78); 场景: key 泄露或需要轮换
- 区别: Generate 增加一个 key; Regenerate 替换现有 key

**B41 Webhook 编辑**:
- 每个 webhook 行的 ⋮ 菜单: [Edit] [Test] [Delete]
- **Edit**: 点击 → inline 展开编辑表单 (URL + Events checkboxes, 预填当前值); [Cancel] [Save]; Save → `PUT /api/white-label/sdk/webhooks/:id` → 成功: toast "Webhook updated" + 表单收起; 失败: toast error
- **Delete**: 点击 → confirm dialog "Delete webhook? Events will no longer be sent to this endpoint." → 确认 → `DELETE /api/white-label/sdk/webhooks/:id` → 成功: 列表刷新 + toast "Webhook deleted"

**B41 Key Regeneration C端影响 (CASCADE)**:
- 旧 key **立即失效**: 所有使用旧 key 的已部署 widgets/pages 将返回 401 Unauthorized
- C端无效 key 处理: widget 显示 "Authentication failed. Contact your project admin." (gray placeholder)
- **重要**: Quick Start 代码块中的 key 自动更新为新 key (W-78 后 refetch)

**B41 "View Docs" 目标**: → `https://docs.taskon.xyz/sdk` (外部链接, target=_blank)

**B41 Empty webhooks section**: 无 webhook 时显示: 空态 "No webhooks configured." + "+ Add Endpoint" 按钮

**B41 Empty API keys**: 不会出现 (WL publish 时自动生成); 但防御性: 如 API 返回空 → 显示 "+ Generate Key" CTA

**B41 Page loading state**: 各 section 独立 skeleton

**侧栏**: WL 子菜单展开, "Overview" 高亮 active (B41 不在子菜单中)

---

## 11. Integration Center

### 11.1 B26 — WL Integration Center

**设计稿**: Node `Abs1E` | URL: `/white-label/integrations`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to White Label"
├── Header: "Integration Center"
├── Status: "2 of 12 integrations active"
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
├── ── Developer Tools (WL 独占) ──
│   ├── API Keys [Available] — "Connect"
│   ├── SDK Configuration [Available] — "Connect"
│   └── SSO / OAuth [Available] — "Connect"
```

**与 Community Integration (B61) 的区别**: WL 多了 "Developer Tools" 分类（API Keys / SDK Config / SSO）。

#### 按钮路由

所有 "Configure" / "Connect" → B44 `/white-label/integrations/:type`

#### 操作详情 (v2.0 新增, W-84 ~ W-88)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-84 | Twitter OAuth | click "Connect" | `/api/white-label/integrations/twitter/auth` | popup window (600×700) → Twitter OAuth 2.0 授权页 → 用户点 "Authorize" → 回调到 TaskOn → popup 关闭 → 父页面刷新: Twitter 状态变 "Connected" (green badge) + 显示 @handle; "Configure" 按钮出现 (→ B44 配置详情) | 用户取消授权 → popup 关闭, 状态不变; Token expired → 见 W-88 |
| W-85 | GA4 Connect | click "Connect" | `POST /api/white-label/integrations/analytics` | inline 展开: "Measurement ID" input (placeholder "G-XXXXXXXXXX") + [Cancel] [Connect]; 提交 → API 验证 ID 格式 (G- prefix + 10 chars) → 成功: status 变 "Connected" + toast "Google Analytics connected"; 配置后所有 WL 页面自动注入 GA tracking code | 无效 ID → "Invalid Measurement ID format. Should be G- followed by 10 characters." |
| W-86 | SSO / OAuth 配置 | click "Connect" (Developer Tools) | — | → 跳转 B44 `/white-label/integrations/sso`: SSO 配置专页; 两种模式: **Wallet Auth** (选择支持的钱包: MetaMask/WalletConnect/Coinbase) + **OAuth2/Custom JWT** (配置: Client ID + Client Secret + Redirect URI + JWT Secret); 保存 → SSO 在所有 WL widget/page 中生效 | — |
| W-87 | SDK Configuration "Connect" | click | — | → 跳转 B41 SDK & API 页面 (不是 B44); 这是快捷入口, 因为 SDK 已有独立完整页面 | — |
| W-88 | 集成 error 状态 | render | — | 集成卡片 error 样式: red border `#EF4444` + error icon (warning, red) + 状态文案 (如 "Token expired" / "API key invalid" / "Connection lost") + "Reconnect →" 按钮; "Reconnect" → 重新触发对应 OAuth/验证流程; 常见错误: Twitter token expired (90d) → re-auth; GA ID changed → re-enter; RPC unreachable → check endpoint | — |
| W-88b | "Disconnect" 操作 | click "Disconnect" | `DELETE /api/white-label/integrations/:type` | 仅 "Connected" 状态的集成显示 "Disconnect" 按钮 (在 ⋮ 菜单或 Configure 按钮旁); 点击 → confirm dialog: "Disconnect {integration}? This will revoke access and may affect active features." → 确认 → API call → 成功: 状态变 "Available" (gray) + toast "{integration} disconnected" + "Connect" 按钮重新出现; 失败: toast error | — |

**B26 页面初始化 (Init)**:
- API: `GET /api/white-label/integrations` — 返回 12 个集成项列表, 每项含 `{ type, name, category, status: 'connected'|'available'|'error', config?, errorMessage? }`
- 加载态: 4 个 section (Social / Blockchain / Analytics / Developer Tools) 各自 skeleton (3-4 卡片占位)
- API 失败: §2.3.4 错误态 + Retry
- 所有集成默认 "Available" 状态 (gray badge)

**B26 页面离开 (Destroy)**: OAuth popup 中途关闭处理: 如 OAuth popup 仍打开, 不强制关闭 (popup 独立生命周期); 父页面不持有 popup 引用

**B26 OAuth popup 完成后刷新**: W-84 popup 关闭 → 父页面通过 `window.addEventListener('message')` 或 polling 检测状态变更 → refetch `GET /api/white-label/integrations` → 更新对应卡片状态

**B26 C端影响**: Twitter/Discord 连接 → C端社交分享功能 enabled; Blockchain 连接 → Token Gate 规则可用; GA 连接 → 所有 WL 页面注入 tracking

**B26 Loading skeleton**: 4 sections × 3 卡片 skeleton

**B26 All "Available" default**: 首次进入, 所有 12 项显示 "Available" gray badge + "Connect" 按钮

**侧栏**: WL 子菜单展开, "Overview" 高亮 active (B26 不在子菜单中)

---

### 11.2 B44 — Integration Configuration Detail (v2.0 新增)

**URL**: `/white-label/integrations/:type`

#### 页面用途
单个集成的详细配置页面, 从 B26 Integration Center 的 "Configure" / "Connect" 按钮进入。每种集成类型有不同的配置表单。

#### 页面初始化 (Init)
- API: `GET /api/white-label/integrations/:type` — 返回集成状态、配置项、连接凭证
- 加载态: config form skeleton
- 路由守卫: `:type` 无效 (不在 12 种支持类型中) → 404 页面, "← Back to Integration Center" 链接到 B26
- `:type` 有效但未连接: 显示 Connect 流程 (OAuth / API Key / 内联配置)

#### 页面结构

```
Content Area
├── Breadcrumb: "White Label > Integrations > {Type Name}"
├── Header
│   ├── Integration icon (40×40) + Name (24px bold)
│   ├── Status badge: Connected (green) / Available (gray) / Error (red)
│   └── "Disconnect" button (secondary, 仅 Connected 时可见)
├── ── Connection Section ──
│   ├── (因类型而异, 见下方类型特定配置)
│   └── "Connect" / "Save" button (purple)
├── ── Settings Section ──
│   ├── (因类型而异, 类型特定设置)
│   └── "Save Settings" button
├── ── Status & Activity ──
│   ├── Last synced: timestamp
│   ├── Events received (24h): count
│   └── "Test Connection" button (secondary)
```

#### 操作表

| 操作 | 触发 | API | 成功 | 失败 |
|------|------|-----|------|------|
| Connect (OAuth) | Click "Connect" | OAuth popup → callback `POST /api/wl/integrations/:type/connect` | 状态→Connected, toast "Connected" | toast error, 保留 Available 状态 |
| Connect (API Key) | Submit key form | `POST /api/wl/integrations/:type/connect` | 同上 | 同上 |
| Save Settings | Click "Save" | `PUT /api/wl/integrations/:type` | toast "Settings saved" | toast error, 表单保留 |
| Test Connection | Click "Test" | `POST /api/wl/integrations/:type/test` | 显示 "✓ Connection OK" (green, 5s 消失) | 显示 "✗ {error detail}" (red) |
| Disconnect | Click "Disconnect" | confirm dialog → `DELETE /api/wl/integrations/:type` | 状态→Available, 清空配置, toast "Disconnected" | toast error |

#### 类型特定配置

**Social & Community**:

| 类型 | Connection | Settings |
|------|-----------|----------|
| **Twitter** | OAuth 2.0 popup (600×700) → callback; 成功显示 @handle | Auto-post enabled (toggle); Default hashtags (text input); Post template (textarea) |
| **Discord** | OAuth popup → Bot 授权 → 选择 server | Notification channel (select from server channels); Event types (checkboxes: task_completed, level_up, sprint_end) |
| **Telegram** | Bot token 输入 (from @BotFather) + Chat ID | Notification events (checkboxes); Language (select) |

**Blockchain & Wallet**:

| 类型 | Connection | Settings |
|------|-----------|----------|
| **Multi-Chain** | Chain 多选 (Ethereum/BSC/Polygon/Arbitrum/Optimism/Base/Avalanche) + 确认 | Custom RPC endpoints (optional, per chain); Block confirmation depth (number, 默认 12) |
| **WalletConnect** | WalletConnect Project ID 输入 (from cloud.walletconnect.com) | 支持的 chains (multi-select); Required methods (multi-select) |
| **On-Chain Verification** | Contract address + ABI (复用 D12 模式) | Verification method (signature/token balance); Cache duration (select: 5min/15min/1h) |

**Analytics & Data**:

| 类型 | Connection | Settings |
|------|-----------|----------|
| **Google Analytics** | Measurement ID 输入 (G-XXXXXXXXXX format) | Data sharing (toggle); Enhanced ecommerce (toggle) |
| **Webhooks** | → 跳转 B41 SDK & API 页面 (快捷入口, 非独立 B44) | — |
| **Data Export** | API key 自动生成 | Export format (select: CSV/JSON); Auto-export frequency (select: daily/weekly/manual); Destination (S3 bucket URL, optional) |

**Developer Tools**:

| 类型 | Connection | Settings |
|------|-----------|----------|
| **API Keys** | → 跳转 B41 SDK & API (快捷入口) | — |
| **SDK Configuration** | → 跳转 B41 SDK & API (快捷入口) | — |
| **SSO / OAuth** | 两种模式: Wallet Auth (选择钱包: MetaMask/WalletConnect/Coinbase) / OAuth2+JWT (Client ID + Client Secret + Redirect URI + JWT Secret) | Session duration (select: 1h/24h/7d/30d); Remember me (toggle) |

> **注意**: Webhooks / API Keys / SDK Configuration 三种类型从 B26 点击后直接跳转 B41, 不经过 B44。B44 处理其余 9 种集成类型。

#### 页面离开 (Destroy)
- 脏数据: 未保存的 Settings 修改 → "Unsaved changes" 确认弹窗

#### 侧栏
WL 子菜单展开, 无子项高亮 (B44 不在子菜单中); Breadcrumb: White Label > Integrations > {Type Name}

---

## 12. Smart Rewards

> **Smart Rewards 是 WL 独占功能**，由 Contract Registry（合约注册）+ Activity Rule Builder（自动化触发）+ Privilege Manager（权益分层）组成。
> 核心逻辑：Contract Registry 注册链上合约 → Rule Builder 定义「什么行为获得什么奖励」→ Privilege Manager 定义「达到什么条件享受什么权益」。
> 三者通过 Points / Badge / Level 等中间状态联动 —— Contract 提供链上事件源，Rule 产出积分/徽章，Privilege 消费积分/徽章作为准入门槛。

### 12.0 B51 — Contract Registry (v2.0 新增, W-89 ~ W-100)

**设计稿**: Node `OKEqS` | URL: `/white-label/contracts`

#### 页面概述

- **核心功能**: 注册和管理链上智能合约，为 Rule Builder (B52) 提供链上事件数据源
- **业务价值**: 项目方注册自己的合约后，Rule Builder 可监听合约事件 (swap / LP / staking / NFT mint 等) 并自动触发奖励
- **与 Rule Builder 的关系**: Contract Registry 是「数据层」（注册合约+解析事件），Rule Builder 是「逻辑层」（IF event THEN reward）
- **与 Privilege Manager 的关系**: Token Gate 类型的 Privilege 通过此处注册的合约进行链上余额查询

#### 页面结构

```
Content Area
├── Header
│   ├── Icon: description (purple bg #1A1033)
│   ├── Title: "Contract Registry"
│   ├── Subtitle: "Register your smart contracts to enable on-chain activity tracking and token-gated privileges"
│   └── "+ Register Contract" button (purple) → D12
├── Stats Row (4 cards)
│   ├── Total Contracts: 12
│   ├── Verified: 10 (green)
│   ├── Events Captured (24h): 8,420
│   └── Active Rules (linked): 8 (→ B52)
├── Contracts Table (card style, rounded corners)
│   ├── Header: "Registered Contracts" + Filter tabs (All / Verified / Pending / Error)
│   ├── Columns: Contract Name | Network | Address | Status | Monitored Events | Actions
│   ├── Row: "DEX Router" | Arbitrum | 0x1234...abcd | ✅ Verified | Swap, AddLiquidity | [View Events] [Edit] [···]
│   ├── Row: "Staking Pool" | Arbitrum | 0x5678...efgh | ✅ Verified | Deposit, Withdraw | [View Events] [Edit] [···]
│   ├── Row: "NFT Collection" | Ethereum | 0x9abc...ijkl | ⏳ Pending | Transfer, Mint | [Verify] [Edit] [···]
│   └── Row: "Token Contract" | BSC | 0xdef0...mnop | ❌ Error | — | [Retry] [Edit] [Delete]
├── SUPPORTED NETWORKS (section header)
│   └── Network badges: Ethereum · BSC · Polygon · Arbitrum · Optimism · Base · Avalanche
├── Tip banner: "Contracts are verified automatically after registration. Verification checks that the contract exists on-chain and its ABI is valid."
```

#### Contract 数据模型

```typescript
interface RegisteredContract {
  id: string;
  name: string;                     // 内部标识名 (如 "DEX Router")
  network: SupportedNetwork;        // 所在链
  address: string;                  // 合约地址 (0x + 40 hex chars)
  abi: object[];                    // 合约 ABI (JSON array)
  status: 'verified' | 'pending' | 'error';
  monitoredEvents: string[];        // 从 ABI 中选择的要监听的事件名列表
  verifiedAt?: string;              // 验证通过时间
  errorMessage?: string;            // status=error 时的错误信息

  // 统计 (只读)
  stats: {
    eventsCaptured24h: number;      // 24h 内捕获的事件数
    linkedRules: number;            // 引用此合约的 Rule 数量
    lastEventAt?: string;           // 最近一次事件时间
  };

  createdAt: string;
  updatedAt: string;
}

type SupportedNetwork =
  | 'ethereum' | 'bsc' | 'polygon' | 'arbitrum'
  | 'optimism' | 'base' | 'avalanche';
```

#### D12 — Contract Register Form (Modal)

**宽度**: 640px

**Modal 结构**:
```
Title Bar: "Register Contract" + × close
Body:
├── Contract Name: text input (placeholder: "e.g. DEX Router, Staking Pool")
├── Network: select dropdown (Ethereum / BSC / Polygon / Arbitrum / Optimism / Base / Avalanche)
│   └── 选择后显示: Network icon + name + Chain ID
├── Contract Address: text input (placeholder: "0x...")
│   └── 验证: 0x prefix + 40 hex characters; 输入后自动 checksum 验证
├── ABI: textarea (placeholder: "Paste contract ABI JSON here...")
│   ├── 或 "Upload ABI File" button (accept: .json)
│   └── 解析状态: 粘贴/上传后自动解析 → "✓ ABI parsed: 12 functions, 5 events found" (green)
│       或 "✗ Invalid ABI format" (red)
├── Events to Monitor: checkbox list (从 ABI 解析结果动态生成)
│   ├── ☑ Swap (address indexed sender, uint256 amountIn, uint256 amountOut)
│   ├── ☑ AddLiquidity (address indexed provider, uint256 amount)
│   ├── ☐ RemoveLiquidity (...)
│   └── ☐ Transfer (...)
│   └── 提示: "Select at least 1 event to monitor"
├── ☑ Verify on Save (toggle, default: on)
│   └── "Automatically verify contract exists on-chain after saving"
Footer: [Cancel] [Register Contract] (purple)
```

**字段验证规格**:

| # | 字段 | 类型 | 必填 | 验证规则 | 说明 |
|---|------|------|------|---------|------|
| W-95 | Contract Name | text input | ✅ | 1-60 chars; 项目内唯一 | 内部标识, 不需要与链上名称一致 |
| W-96 | Network | select | ✅ | 从支持列表选择 | Ethereum / BSC / Polygon / Arbitrum / Optimism / Base / Avalanche |
| W-97 | Contract Address | text input | ✅ | `0x` + 40 hex chars; EIP-55 checksum | 粘贴后自动格式化; 重复地址+同一链 → "This contract is already registered" |
| W-98 | ABI | textarea / JSON upload | ✅ | valid JSON array; 至少包含 1 个 event | 粘贴或上传 .json; 解析后显示 function/event 数量; 不可编辑解析结果 |
| W-99 | Events to Monitor | checkbox list | ✅ (≥1) | 从 ABI events 动态生成 | 每个 event: name + indexed params 签名; 选择 ≥1 个 |
| W-100 | Verify on Save | toggle | — | 默认 on | on → 保存后立即 POST `/api/wl/contracts/:id/verify`; off → 保存为 pending, 后续手动验证 |

**交互逻辑**:

1. **ABI 解析联动**: 粘贴或上传 ABI → 客户端 JSON parse → 提取所有 `type: "event"` 条目 → 动态渲染 Events checkbox list; ABI 无效 → checkbox 列表不显示 + 红色错误
2. **Contract Address 重复检查**: 输入后 debounce 500ms → `GET /api/wl/contracts/check?address={}&network={}` → 重复 → 红色提示 "Already registered"
3. **Verify on Save 流程**: 保存 → `POST /api/wl/contracts` → status=pending → 自动触发 `POST /api/wl/contracts/:id/verify` → 验证合约在链上存在 + ABI 匹配 → status 变 verified/error
4. **编辑模式**: 已验证合约可编辑 name 和 monitored events; **不可修改** address/network/ABI (需删除重建); readonly 字段样式: 灰色 bg `#1E293B` + text `#94A3B8` + lock icon
5. **验证时机**: **所有字段在 submit 时统一验证** (非逐字段 blur); 提交按钮 "Register Contract" 点击 → 客户端全量校验 → 有错误: 滚动到第一个错误字段 + 字段红色 border; 通过: 调用 POST API
6. **ABI 无效 JSON 时机**: 粘贴/上传后**立即** client-side JSON.parse 验证 (非 submit 时); 无效 → 红色错误 "Invalid ABI format"; Events checkbox 列表不显示
7. **ABI 重新粘贴**: 重新粘贴新 ABI → Events checkbox list **重置** (清空之前的选择, 重新从新 ABI 解析生成)
8. **Large ABI parsing**: > 100 events 时显示 "Parsing ABI..." inline spinner (300ms); 解析完成后渲染 checkbox list
9. **Modal 动画**: 打开: fade-in 200ms + scale 0.95→1; 关闭: fade-out 150ms; 关闭方式: × 按钮 / 点击 overlay / Esc 键; 脏数据: 有修改时关闭 → "Discard changes?" 确认

#### 操作详情 (W-89 ~ W-94)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-89 | 完整页面结构 | render | `GET /api/wl/contracts` | Node `OKEqS`: Header + 4 Stats + Contracts Table (filter tabs + columns) + Supported Networks badges + Tip banner; Stats 中 "Active Rules" 数字可点击 → B52 Rule Builder | API 失败 → skeleton → error state "Unable to load contracts" + Retry |
| W-90 | Stats Row | render | `GET /api/wl/contracts/stats` | 4 cards: Total Contracts (count) / Verified (count, green text `#16A34A`) / Events Captured 24h (number + sparkline) / Active Rules (count, 可点击 → B52); "Active Rules" 显示引用任何已注册合约的 rule 数量 | — |
| W-91 | Contracts Table | render + filter | `GET /api/wl/contracts?status={filter}` | 列: Contract Name / Network (chain icon + name) / Address (truncated 0x1234...abcd + Copy icon) / Status badge (Verified=green / Pending=amber spinner / Error=red) / Monitored Events (comma-separated event names) / Actions; Filter tabs: All / Verified / Pending / Error; 分页 20/page | — |
| W-92 | "+ Register Contract" | click | — | → 打开 D12 Modal (create mode, 空表单) | — |
| W-93 | Contract verify | auto/manual | `POST /api/wl/contracts/:id/verify` | **自动**: Verify on Save=on 时保存后立即触发; **手动**: Pending/Error 行的 "Verify" / "Retry" 按钮; 验证过程: button spinner → 链上查询合约代码 + ABI event 签名匹配 → 成功: status=verified + green ✅ + toast "Contract verified" → 失败: status=error + errorMessage 显示 | 常见错误: "Contract not found on {network}" / "ABI does not match deployed bytecode" / "Network RPC timeout" |
| W-94 | Row 操作 | click action buttons | various | **View Events** → side panel (720px): 最近 100 条捕获事件列表 (timestamp + event name + params + tx hash link); **Edit** → D12 (edit mode, 预填, address/network/ABI readonly); **Delete** → confirm dialog "Delete contract '{name}'?" → 检查: 如被 Rule 引用 → dialog 额外警告 "This contract is referenced by {n} active rules. Deleting will disable those rules." → 确认 → `DELETE /api/wl/contracts/:id` → 刷新; **被 Rule 引用时仍可删除** (级联 disable rules) | — |

**B51 页面离开 (Destroy)**: View Events side panel: 关闭 panel, 清理事件列表数据

**B51 Table row 点击行为**: **行不可点击** (非跳转); 操作通过行内 action buttons (View Events / Edit / ⋮ → Delete/Verify/Retry); 这与 B52/B53 的 "row click → modal" 模式不同

**B51 Contract 验证后 Stats 更新**: W-93 verify 成功 → Stats Row "Verified" 数值 +1 (optimistic update); "Pending" 或 "Error" 数值 -1; 无需全页 refetch

**B51 D12 save 后 Table 更新**: D12 关闭 → `GET /api/wl/contracts` refetch 整个列表 (非 optimistic)

**B51 "Active Rules" stat 点击**: → B52 Rule Builder (不带筛选参数, 显示全部 rules)

**B51 Filter tabs URL 同步**: 不同步到 URL params (客户端过滤)

**B51 Individual contract verification error**: 行内 Status badge 显示 "Error" (red) + hover tooltip 显示 errorMessage (如 "Contract not found on Arbitrum")

**侧栏**: WL 子菜单展开, SMART REWARDS section → "Contracts" 高亮 active

#### Contract 状态流转

```
              ┌──────────┐
  Register    │ Pending  │   (刚注册，等待验证)
  ──────────► │          │
              └────┬─────┘
                   │ Verify
        ┌──────────┼──────────┐
        ▼          │          ▼
   ┌──────────┐   │    ┌──────────┐
   │ Verified │   │    │  Error   │
   │          │   │    │          │
   └──────────┘   │    └────┬─────┘
                  │         │ Retry
                  └─────────┘
```

- **Pending**: 刚注册，等待链上验证
- **Verified**: 合约在链上确认存在 + ABI 匹配 → 可被 Rule Builder 引用
- **Error**: 验证失败 (合约不存在 / ABI 不匹配 / RPC 超时) → 可 Retry

#### API

| Endpoint | Method | 说明 | 请求体关键字段 |
|----------|--------|------|---------------|
| `/api/wl/contracts` | GET | 合约列表 | query: `?status=verified&page=1&limit=20` |
| `/api/wl/contracts` | POST | 注册合约 | `{ name, network, address, abi, monitoredEvents, verifyOnSave }` |
| `/api/wl/contracts/:id` | GET | 合约详情 | — |
| `/api/wl/contracts/:id` | PUT | 更新合约 (仅 name + monitoredEvents) | `{ name, monitoredEvents }` |
| `/api/wl/contracts/:id` | DELETE | 删除合约 (级联 disable 引用的 rules) | — |
| `/api/wl/contracts/:id/verify` | POST | 触发链上验证 | — |
| `/api/wl/contracts/:id/events` | GET | 查看捕获事件历史 | query: `?page=1&limit=100` |
| `/api/wl/contracts/check` | GET | 重复检查 | query: `?address=0x...&network=ethereum` |
| `/api/wl/contracts/stats` | GET | 汇总统计 | — |

---

### 12.1 B52 — Activity Rule Builder

**设计稿**: Node `4aAo7` | URL: `/white-label/rules`

#### 页面概述

- **核心功能**: 可视化 IF-THEN 规则引擎，将链上/链下用户行为自动映射为奖励
- **业务价值**: 项目方无需开发即可配置「做了 X 就给 Y」的自动奖励逻辑，驱动用户完成高价值行为
- **与 Community 模块关系**: Rule Builder 产出的积分会流入 Community 的 Points & Level 系统，进而影响 Leaderboard 排名、Benefits Shop 兑换能力、Privilege 资格

#### 页面结构

```
Content Area
├── Header
│   ├── Icon: bolt (amber bg #1F1A08)
│   ├── Title: "Activity Rule Builder"
│   ├── Subtitle: "Create if-then rules to automatically reward users for on-chain actions"
│   └── "+ Create Rule" button (purple) → D13
├── Stats Row (4 cards)
│   ├── Active Rules: 8
│   ├── Rules Triggered (24h): 3,247
│   ├── Points Distributed (24h): 18,540 (↑ 12% vs yesterday, purple highlight)
│   └── Unique Users Rewarded: 892
├── Active Rules Table (card style, rounded corners)
│   ├── Header: "Active Rules"
│   ├── Row: "Reward Every Swap" [Active] — IF swap ≥ $10 on DEX Router → THEN +50 pts · Triggered 1,247 times today
│   ├── Row: "Daily First Action 2x" [Active] — IF first daily transaction → THEN +100 pts (2x multiplier) · Triggered 423 times today
│   ├── Row: "LP Bonus" [Draft] — IF add liquidity ≥ $500 → THEN +200 pts · Triggered 89 times today
│   └── Row: "Staking Milestone" [Active] — IF stake ≥ 30 days continuously → THEN +500 pts (one-time) · Triggered 34 times today
├── RULE PRESETS (section header)
│   ├── "Reward Every Swap" — Auto-reward users for each token swap above a threshold. Best for DEX activity.
│   ├── "Daily First Action 2x" — Double points for the first on-chain action each day. Drives daily active usage.
│   └── "LP Bonus + Staking Milestone" — Reward liquidity provision and long-term staking with tiered point bonuses.
├── ANTI-SYBIL CONFIGURATION (section header)
│   ├── Min Wallet Age: 90 days
│   ├── Min Transactions: >10 unique txs
│   ├── Bot Detection: Enabled ✅ (Preview button)
│   └── Note: "Anti-sybil checks are applied before any rule triggers. Wallets that fail checks are silently excluded."
```

#### Rule 数据模型

```typescript
interface ActivityRule {
  id: string;
  name: string;                    // 规则名称，如 "Reward Every Swap"
  status: 'active' | 'draft' | 'paused';

  // IF 部分 — 触发条件
  trigger: {
    event: TriggerEvent;           // 触发事件类型
    condition: string;             // 条件描述，如 "Any task in Sector: Getting Started"
    value?: number;                // 阈值，如 swap ≥ $10 的 10
    chainId?: string;              // 链 ID（链上事件必填）
    contractAddress?: string;      // 合约地址（链上事件必填，关联 B51 合约注册）
  };

  // THEN 部分 — 执行动作
  action: {
    type: 'award_points' | 'grant_badge' | 'upgrade_tier' | 'webhook';
    pointType?: string;            // 积分类型 ID (EXP / GEM / 自定义)
    value?: number;                // 积分数量，如 50 XP
    badgeId?: string;              // 徽章 ID（type=grant_badge 时）
    tierId?: string;               // 权限层级 ID（type=upgrade_tier 时）
    webhookUrl?: string;           // 外部回调（type=webhook 时）
    multiplier?: number;           // 倍率，如 2x
  };

  // 频率控制
  frequency: 'once' | 'daily' | 'unlimited';

  // 反作弊
  antiSybil: AntiSybilConfig;

  // 统计（只读）
  stats: {
    triggeredToday: number;
    triggeredTotal: number;
    pointsDistributed24h: number;
    uniqueUsersRewarded: number;
  };

  createdAt: string;
  updatedAt: string;
}

// 支持的触发事件
type TriggerEvent =
  // 链上事件（需要 B51 合约注册）
  | 'token_swap'          // 代币兑换
  | 'add_liquidity'       // 添加流动性
  | 'remove_liquidity'    // 移除流动性
  | 'token_transfer'      // 代币转账
  | 'nft_mint'            // NFT 铸造
  | 'nft_transfer'        // NFT 转移
  | 'staking_deposit'     // 质押存入
  | 'staking_withdraw'    // 质押提取
  | 'contract_interaction'// 通用合约交互
  // 平台事件（无需合约注册）
  | 'task_completed'      // 完成任务
  | 'daily_login'         // 每日登录
  | 'referral_success'    // 推荐成功
  | 'milestone_reached'   // 达成里程碑
  | 'level_up';           // 升级

interface AntiSybilConfig {
  minWalletAge: number;    // 最小钱包年龄（天）
  minTransactions: number; // 最小交易数
  botDetection: boolean;   // 是否开启机器人检测
}
```

#### IF-THEN 规则引擎逻辑

```
用户行为 → 事件匹配 → Anti-Sybil 检查 → 条件判断 → 频率检查 → 执行动作 → 更新统计
```

**详细流程**:

1. **事件捕获**: 链上事件通过 B51 注册的合约 + 链上监听器捕获；平台事件由 TaskOn 内部系统产生
2. **Anti-Sybil 前置检查**: 所有规则触发前，先检查钱包年龄/交易数/机器人检测。**不通过的钱包静默排除**（不报错，不通知用户）
3. **条件匹配**: 检查事件是否满足 IF 条件（如 swap 金额 ≥ 阈值、在指定 Sector 完成任务等）
4. **频率检查**:
   - `once`: 每个用户只触发一次（全生命周期去重）
   - `daily`: 每个用户每日最多触发一次（UTC 0:00 重置）
   - `unlimited`: 每次满足条件都触发（适合 swap 奖励等高频场景）
5. **执行动作**:
   - `award_points` → 调用 Points 系统 API 增加积分 → **联动 Leaderboard 排名实时更新**
   - `grant_badge` → 调用 Badge 系统发放徽章 → **可能触发 Privilege 的 Achievement-Based 资格**
   - `upgrade_tier` → 直接提升用户 Privilege 层级
   - `webhook` → POST 到项目方指定 URL（用于项目方自有系统联动）

#### 联动关系图

```
┌─────────────────────────────────────────────────────────────┐
│                    Activity Rule Builder                     │
│              IF (trigger) → THEN (action)                   │
└──────────┬──────────────┬───────────────┬───────────────────┘
           │              │               │
    award_points     grant_badge     upgrade_tier
           │              │               │
           ▼              ▼               ▼
  ┌────────────┐  ┌───────────┐  ┌────────────────┐
  │ Points &   │  │  Badges   │  │   Privilege     │
  │ Level      │  │  (B31i)   │  │   Manager (B53) │
  │ (B31a)     │  └─────┬─────┘  └────────┬───────┘
  └──────┬─────┘        │                 │
         │              │                 │
         ▼              ▼                 ▼
  ┌────────────┐  ┌───────────────┐  ┌──────────────┐
  │ Leaderboard│  │ Privilege     │  │ C-End 用户   │
  │ (C03)      │  │ Qualification │  │ 享受权益     │
  │ LB Sprint  │  │ (Mode C)     │  │ (fee discount│
  │ (C04)      │  └───────────────┘  │  gas rebate  │
  └──────┬─────┘                     │  yield boost)│
         │                           └──────────────┘
         ▼
  ┌────────────┐
  │ Benefits   │
  │ Shop (C06) │
  │ 积分兑换   │
  └────────────┘
```

#### Rule Presets 业务说明

| Preset | 适用场景 | 默认配置 | 预期效果 |
|--------|---------|---------|---------|
| **Reward Every Swap** | DEX 项目促交易量 | IF swap ≥ $10 → +50 pts, unlimited | 每笔合格交易即时奖励，高频激励 |
| **Daily First Action 2x** | 所有项目促 DAU | IF first daily tx → +100 pts (2x), daily | 鼓励每日至少一次链上交互 |
| **LP Bonus + Staking Milestone** | DeFi 项目促 TVL | IF LP ≥ $500 → +200 pts, unlimited; IF stake ≥ 30d → +500 pts, once | 短期流动性 + 长期锁仓双重激励 |

#### Anti-Sybil 配置说明

| 检查项 | 默认值 | 说明 | 失败处理 |
|--------|-------|------|---------|
| **Min Wallet Age** | 90 天 | 钱包首次交易距今天数 | 静默排除，不触发规则 |
| **Min Transactions** | >10 笔 | 钱包历史唯一交易数 | 静默排除，不触发规则 |
| **Bot Detection** | 启用 | 机器学习模型检测刷量行为 | 静默排除 + 标记可疑钱包 |

> **重要**: Anti-Sybil 是全局配置，对该项目下所有 Rule 生效。修改后立即对新触发的事件生效，已发放的奖励不追溯。

#### D13 — Activity Rule Editor (Modal)

**设计稿**: Node `IJZ0E` | 宽度 800px

**Modal 结构**:
```
Title Bar: "Create Activity Rule" + × close
Body:
├── Rule Name: text input (placeholder: "e.g. Reward Daily Logins, NFT Holder Bonus")
├── IF Block (amber border #FED7AA, bg #1F1508)
│   ├── Badge: "IF" (orange #EA580C)
│   ├── Label: "When this event happens..."
│   ├── Event dropdown: "Task Completed" (可选: 见 TriggerEvent 类型)
│   └── Condition dropdown: "Any task in Sector: Getting Started"
├── ↓ Arrow
├── THEN Block (green border #BBF7D0, bg #0A1F1A)
│   ├── Badge: "THEN" (green #16A34A)
│   ├── Label: "Perform this action..."
│   ├── Action dropdown: "Award Points"
│   └── Value input: "50" + unit "XP"
├── Frequency: toggle group [Once per user (default) | Daily | Unlimited]
Footer: [Cancel] [Create Rule] (purple)
```

**交互逻辑**:

1. **Event dropdown 选项联动**:
   - 选择链上事件 (swap/LP/staking 等) → Condition 下拉显示已注册合约列表（来自 B51）
   - 如无合约注册 → 显示 "No contracts registered. Register one first →" 链接到 B51
   - 选择平台事件 (task_completed/daily_login 等) → Condition 下拉显示 Sectors/任务类型
2. **Action dropdown 选项联动**:
   - Award Points → 显示 value 数字输入 + point type 选择器 (EXP/GEM/自定义)
   - Grant Badge → 显示 badge 选择器（来自 B31i 已配置的徽章列表）
   - Upgrade Tier → 显示 tier 选择器（来自 B53 已配置的层级列表）
   - Webhook → 显示 URL 输入框
3. **Frequency 互斥选择**: 三个选项 radio 样式，选中为蓝色填充 (#0F1A2E + #3B82F6 文字)
4. **保存**:
   - 创建模式: POST `/api/wl/rules` → 默认 status=draft
   - 编辑模式: PUT `/api/wl/rules/:id`
   - 保存成功 → 关闭 Modal → 刷新 B52 规则列表
5. **从 Preset 进入**: 自动填充所有字段，用户可修改后保存; Modal 标题变为 "Create Rule from Preset"
6. **Modal 动画**: 打开: fade-in 200ms + scale 0.95→1; 关闭: fade-out 150ms; 关闭方式: × / overlay / Esc; 脏数据: 有修改 → "Discard changes?" 确认

**D13 "Create Rule" / "Save Rule" 按钮条件**:
- 必填字段: Rule Name (1-60 chars) + Event (已选) + Action type (已选) + Action value (如 points 数量) + Frequency (已选)
- Disabled: 任何必填字段为空 → opacity 0.5 + tooltip "Fill in all required fields"
- Loading: button spinner + disabled 防重复

**D13 Create 默认状态**: 创建的 rule 默认 `status: 'draft'`; 需在 B52 列表中通过 toggle 手动激活; 不可在 D13 内直接创建为 Active

**D13 Event dropdown 数据源 API**:
- 链上事件: Event 列表硬编码 (见 TriggerEvent 类型); 选择链上事件后 → Condition dropdown 调用 `GET /api/wl/contracts?status=verified` 获取已验证合约列表
- 平台事件: Event 列表硬编码; 选择 task_completed → Condition dropdown 调用 `GET /api/community/sectors` 获取 Sector 列表; 其他平台事件无子条件

**D13 Action dropdown 数据源 API**:
- Award Points → Point type selector: `GET /api/community/point-types` (EXP / GEM / 自定义)
- Grant Badge → Badge selector: `GET /api/community/badges?status=active` → 0 badges: 显示 "No badges configured. Create one in Community → Badges." + 链接
- Upgrade Tier → Tier selector: `GET /api/wl/privileges?status=active` → 0 tiers: 显示 "No privilege tiers configured. Create one in Privilege Manager." + 链接
- Webhook → URL input (placeholder "https://...")

**D13 无合约时链上事件**: 选择链上事件类型 → Condition dropdown 显示 "No contracts registered. Register one first →" (链接到 B51); 用户无法完成条件配置, "Create Rule" 按钮 disabled

**D13 Edit mode**: 可修改所有字段 (包括 Active rule 的 trigger event type); 但修改触发条件的影响: toast 提醒 "Changing trigger will reset the rule's statistics"

**D13 No badges/tiers visual**: 对应 Action 选择后, selector 区域显示 amber 提示 + 跳转链接 (新标签页)

#### 按钮路由

| 按钮 | 位置 | 目标 | 说明 |
|------|------|------|------|
| "+ Create Rule" | Header | → D13 (create) | 空表单 |
| Row click | Rules Table | → D13 (edit) | 预填已有数据 |
| "Use Preset →" | Preset Cards | → D13 (preset) | 预填模板数据 |
| Rule toggle | Table row | (API) `PUT /api/wl/rules/:id` | 切换 active↔paused |
| Filter tabs | Table 上方 | (前端筛选) | All / Active / Draft / Paused |
| "Preview" | Anti-Sybil Bot Detection | (action) | 预览当前被排除的钱包列表 (见下方 Anti-Sybil Preview 说明) |

#### Rule 状态流转

```
                  ┌─────────┐
     Create Rule  │  Draft  │
     ─────────►   │         │
                  └────┬────┘
                       │ Toggle Enable
                       ▼
                  ┌─────────┐    Toggle Disable    ┌─────────┐
                  │ Active  │ ◄──────────────────► │ Paused  │
                  │         │    Toggle Enable      │         │
                  └────┬────┘                      └────┬────┘
                       │ Delete                         │ Delete
                       ▼                                ▼
                  ┌─────────┐
                  │ Deleted │ (soft delete, 可恢复)
                  └─────────┘
```

- **Draft → Active**: 首次启用，开始监听事件并触发
- **Active → Paused**: 暂停触发，已发放奖励不追溯
- **Paused → Active**: 恢复触发，从恢复时刻起生效
- **Delete**: 软删除，30 天内可恢复 (恢复 UI: 暂无独立页面, 通过 API `/api/wl/rules?status=deleted` 查询 + `PUT /api/wl/rules/:id { status: 'draft' }` 恢复; 未来版本可加入回收站 UI)

#### API

| Endpoint | Method | 说明 | 请求体关键字段 |
|----------|--------|------|---------------|
| `/api/wl/rules` | GET | 获取规则列表 | query: `?status=active&page=1&limit=20` |
| `/api/wl/rules` | POST | 创建规则 | `{ name, trigger, action, frequency }` |
| `/api/wl/rules/:id` | GET | 获取单条规则详情 | — |
| `/api/wl/rules/:id` | PUT | 更新规则 | 同 POST body + `{ status }` |
| `/api/wl/rules/:id` | DELETE | 软删除规则 | — |
| `/api/wl/rules/presets` | GET | 获取预设模板列表 | — |
| `/api/wl/rules/stats` | GET | 获取汇总统计 | query: `?period=24h` |
| `/api/wl/anti-sybil` | GET | 获取反作弊配置 | — |
| `/api/wl/anti-sybil` | PUT | 更新反作弊配置 | `{ minWalletAge, minTransactions, botDetection }` |
| `/api/wl/anti-sybil/preview` | GET | 预览被排除钱包 | query: `?page=1&limit=50` |

**B52 页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/wl/rules` — Rules Table (默认 status=all)
  2. `GET /api/wl/rules/stats` — Stats Row (4 cards)
  3. `GET /api/wl/rules/presets` — Rule Presets (3 cards)
  4. `GET /api/wl/anti-sybil` — Anti-Sybil Configuration
- 加载态: 各 section 独立 skeleton (Stats Row / Rules Table / Presets / Anti-Sybil)
- API 失败: 各 section 独立 §2.3.4 错误态 + Retry

**B52 页面离开 (Destroy)**: 无状态需清理

**B52 Rule toggle (active↔paused)**:
- UI: 每行右侧 toggle switch (active=on purple / paused=off gray)
- 交互: 点击 toggle → **无确认弹窗** (即时切换) → `PUT /api/wl/rules/:id { status: 'active'|'paused' }` → 成功: toggle 动画 + toast "Rule {activated|paused}"; Stats Row 即时更新 (Active Rules ±1)
- Draft → Active: 首次启用, toggle 切换 → `PUT { status: 'active' }` → **仅处理新事件** (不处理 pending/历史事件)
- Paused → Active: 恢复, 从恢复时刻起处理新事件

**B52 Anti-Sybil "Preview" UI**:
- 点击 "Preview" 按钮 → side panel (720px, 右侧滑入)
- Panel 标题: "Excluded Wallets Preview"
- 内容: `GET /api/wl/anti-sybil/preview` → 被排除钱包列表 (地址 / 排除原因 / 首次检测时间)
- 分页: 50/page
- 关闭: × 按钮 / 点击外部

**B52 Anti-Sybil 编辑**:
- Min Wallet Age / Min Transactions: **inline 编辑** (点击数值 → input 变为可编辑 → blur 或 Enter 保存 → `PUT /api/wl/anti-sybil`)
- Bot Detection: toggle switch (同 rule toggle 逻辑)
- 保存后即时生效 (新触发的事件), toast "Anti-sybil settings updated"

**B52 Rule row 点击**: 行点击 → D13 Modal (edit mode, 预填该 rule 数据); 但 toggle / action buttons 区域点击不触发行点击 (阻止冒泡)

**B52 D13 save 后 Table 刷新**: D13 关闭 → refetch `GET /api/wl/rules` 更新列表

**B52 Rule toggle 后 Stats 更新**: optimistic update (Active Rules ±1, 不等待 API 返回)

**B52 "Triggered X times today"**: 页面加载时获取, 非实时 (刷新页面更新)

**B52 Loading skeleton**: Stats Row (4 skeleton) + Rules Table (5 行 skeleton) + Presets (3 卡片 skeleton) + Anti-Sybil (1 card skeleton)

**B52 Soft-delete 恢复**: 当前版本无 UI 入口; Filter tabs 不包含 "Deleted"; 后续版本可加 "Deleted" tab

**侧栏**: WL 子菜单展开, SMART REWARDS section → "Rule Builder" 高亮 active

---

### 12.2 B53 — Privilege Manager

**设计稿**: Node `5xwYN` | URL: `/white-label/privileges`

#### 页面概述

- **核心功能**: 定义用户权益层级，根据积分/等级/徽章/持仓等条件自动或手动授予项目原生权益
- **业务价值**: 让项目方能将 TaskOn 积分体系转化为**真实产品权益**（手续费折扣、Gas 返还、收益加成、优先体验等），形成「行为→积分→权益」的完整价值闭环
- **与 Rule Builder 的关系**: Rule Builder 是「生产端」（产出积分/徽章），Privilege Manager 是「消费端」（消耗积分/徽章作为准入门槛）

#### 页面结构

```
Content Area
├── Header
│   ├── Icon: stars (purple bg #1A1033)
│   ├── Title: "Privilege Manager"
│   ├── Subtitle: "Define project-native privileges that reward loyal users with real product benefits"
│   └── "+ Create Privilege" button (purple) → D14
├── Stats Row (4 cards)
│   ├── Active Privileges: 5
│   ├── Active Holders: 1,247
│   ├── Total Value Distributed: $12,840 (↑ 8% this month, purple highlight)
│   └── API Integration: Connected ✅ (green text)
├── Active Privileges Table (card style)
│   ├── Header: "Active Privileges"
│   ├── Row: "Trading Fee Discount" [Active] — 10% fee discount · Status-based: Level ≥ 3 · 487 active holders
│   ├── Row: "Gas Rebate" [Active] — Up to 50% gas refund · Status-based: 1,000+ points balance · 312 active holders
│   ├── Row: "Yield Boost" [Active] — +2% APY bonus · Achievement-based: Diamond Hands badge · 156 active holders
│   └── Row: "Priority Access" [Active] — Early access to new features · Achievement-based: OG Pioneer milestone · 89 holders
├── QUALIFICATION MODES (section header)
│   ├── Mode A: Status-Based — "Auto-grant privileges based on user level or points balance. Privileges activate/expire automatically as status changes."
│   ├── Mode B: Redemption-Based — "Users spend points in Benefits Shop to claim time-limited privilege vouchers. Managed via existing Shop module."
│   └── Mode C: Achievement-Based — "One-time privilege grant upon earning a specific Badge or completing a Milestone. Permanent or time-limited."
├── INTEGRATION STATUS (section header)
│   ├── API Connection: Enabled ✅
│   ├── Webhook URL: https://...
│   ├── Usage Reporting: Active ✅
│   └── Note: "Your project queries the TaskOn API to check user privilege status and applies benefits in your own product logic."
```

#### Privilege Tier 数据模型

```typescript
interface PrivilegeTier {
  id: string;
  name: string;                     // 层级名称，如 "VIP", "Whale", "OG"
  rankOrder: number;                // 排序序号（1 = 最高层级）
  status: 'active' | 'draft' | 'paused';

  // 视觉配置
  icon: string;                     // Material Symbols icon name (如 "diamond")
  color: string;                    // 层级主题色 (#9B7EE0 / #F59E0B / #3B82F6 / #EF4444)

  // 准入条件（三种模式互斥）
  qualification: {
    mode: 'status_based' | 'redemption_based' | 'achievement_based';

    // Mode A: Status-Based
    condition?: 'token_gate' | 'level_req' | 'points_threshold' | 'manual_assignment';
    tokenGate?: {                   // condition=token_gate
      contractAddress: string;      // 代币合约地址（关联 B51）
      minBalance: number;           // 最低持仓量
      chainId: string;
    };
    levelReq?: {                    // condition=level_req
      minLevel: number;             // 最低等级
    };
    pointsThreshold?: {             // condition=points_threshold
      pointType: string;            // 积分类型 ID
      minBalance: number;           // 最低积分余额
    };

    // Mode B: Redemption-Based
    shopItemId?: string;            // 关联 Benefits Shop 商品 ID

    // Mode C: Achievement-Based
    badgeId?: string;               // 需要持有的徽章 ID
    milestoneId?: string;           // 需要达成的里程碑 ID
  };

  // 权益列表（多选）
  privileges: PrivilegeItem[];

  // 时效
  duration: 'permanent' | 'time_limited';
  durationDays?: number;            // time_limited 时的有效天数

  // 统计（只读）
  stats: {
    activeHolders: number;
    totalValueDistributed: number;
    lastGrantedAt: string;
  };

  createdAt: string;
  updatedAt: string;
}

// 权益项（项目方在自己产品中执行）
interface PrivilegeItem {
  type: 'early_access' | 'exclusive_shop' | 'point_multiplier' | 'custom_badge'
      | 'fee_discount' | 'gas_rebate' | 'yield_boost' | 'priority_support' | 'custom';
  label: string;                    // 显示名称
  value?: number;                   // 具体数值（如 10% discount, 2x multiplier）
  enabled: boolean;
}
```

#### 三种 Qualification Mode 详细说明

##### Mode A: Status-Based（最常用）

**逻辑**: 系统根据用户当前状态**自动判定**是否有资格，资格随状态变化**实时生效/失效**。

| Condition | 判定依据 | 数据源 | 实时性 |
|-----------|---------|--------|-------|
| **Token Gate** | 钱包持有指定代币 ≥ N | B51 注册合约 + 链上查询 | 每次 API 查询时实时检查 |
| **Level Req** | 用户等级 ≥ N | Points & Level (B31a) | 等级变更时即时更新 |
| **Points Threshold** | 积分余额 ≥ N | Points & Level (B31a) | 积分变更时即时更新 |
| **Manual Assignment** | 管理员手动添加 | D15 Members Panel | 手动操作即时生效 |

**联动规则**:
- 用户积分从 1,000 降到 999 → 自动**失去** Points Threshold 类权益
- 用户从 Level 3 升到 Level 4 → 自动**获得** Level Req ≥ 4 的权益
- 用户出售代币导致余额低于阈值 → 下次 API 查询时**失去**权益

##### Mode B: Redemption-Based

**逻辑**: 用户在 C-End Benefits Shop (C06) 中花费积分**主动兑换**权益凭证，凭证有时效。

**联动流程**:
```
B53 创建 Privilege Tier (Mode B)
  → 自动在 Benefits Shop (B31g) 创建对应商品
  → C-End Shop (C06) 展示可兑换
  → 用户花费积分兑换
  → 获得时限权益凭证（如 30 天 VIP）
  → 凭证到期自动失效
```

- **与 Benefits Shop 的关系**: Privilege Manager 创建 Mode B 层级时，会自动在 B31g Benefits Shop 中创建关联商品。商品价格/库存在 B31g 管理，权益内容在 B53 管理。
- **到期提醒**: 凭证到期前 3 天，系统通过 Webhook 通知项目方 → 项目方可推送续费提醒

##### Mode C: Achievement-Based

**逻辑**: 用户获得指定 Badge 或达成指定 Milestone 后，**一次性**获得权益（可永久或限时）。

**联动流程**:
```
Rule Builder (B52) 触发 grant_badge → 用户获得 Badge
  → Privilege Manager 检测到 Badge 匹配 → 自动授予权益

Community Milestone (B31f) 达成 → 用户完成 Milestone
  → Privilege Manager 检测到 Milestone 匹配 → 自动授予权益
```

- **不可撤销**: Badge 一旦获得不可回收，因此 Achievement-Based 权益一旦授予也不自动失效（除非设置了 time_limited）
- **适用场景**: OG 用户永久福利、早期贡献者专属权益

#### 权益项说明

| 权益项 | 说明 | 执行方 | 典型值 |
|--------|------|--------|-------|
| **Early access to new quests** | 新活动/功能提前体验 | 项目方产品 | 提前 24-48h |
| **Exclusive shop items** | 专属商品/NFT 可见 | Benefits Shop | 特定商品只对该层级可见 |
| **2x point multiplier** | 积分翻倍 | TaskOn 系统自动执行 | 2x / 3x / 5x |
| **Custom badge** | 专属徽章展示 | TaskOn 系统自动发放 | 身份标识 |
| **Fee discount** | 交易手续费折扣 | 项目方合约/后端 | 10% / 20% / 50% |
| **Gas rebate** | Gas 费返还 | 项目方合约 | Up to 50% |
| **Yield boost** | 收益加成 | 项目方 DeFi 合约 | +2% APY |
| **Priority support** | 优先客服 | 项目方运营 | 专属频道 |

> **重要**: Fee discount / Gas rebate / Yield boost 等链上权益由**项目方自行实现**。TaskOn 通过 API 提供用户权益状态查询，项目方在自己的合约/后端逻辑中检查并执行。

#### 权益查询 API（项目方调用）

```
项目方产品 → GET /api/wl/privileges/check?wallet=0x...
  → 返回该钱包当前拥有的所有权益列表
  → 项目方根据返回结果在自己产品中执行折扣/返还等逻辑
```

这是**项目方集成的核心接口**，Dev Kit (B48) 中会包含调用示例。

#### D14 — Privilege Tier Editor (Modal)

**设计稿**: Node `FypcB` | 宽度 640px

**Modal 结构**:
```
Title Bar: "Create Privilege Tier" + × close
Body:
├── Row: Tier Name (text input, placeholder: "e.g. VIP, Whale, OG") + Rank Order (number input, default: 1)
├── Row: Tier Icon (icon preview, diamond icon in purple circle) + Tier Color (4 color swatches: purple/amber/blue/red, purple selected with white border)
├── Qualification Condition
│   ├── Dropdown: "Token Gate (hold 100+ tokens)" (options: Token Gate / Level Req / Points Threshold / Manual Assignment)
│   └── Hint text: "Token Gate · Level Req · Points Threshold · Manual Assignment"
├── Privileges (checkbox list)
│   ├── ☑ Early access to new quests
│   ├── ☑ Exclusive shop items
│   ├── ☐ 2x point multiplier
│   └── ☐ Custom badge
Footer: [Cancel] [Create Tier] (purple)
```

**交互逻辑**:

1. **Qualification Condition 联动**:
   - 选 Token Gate → 展开子表单: Contract Address (从 B51 列表选择) + Min Balance + Chain
   - 选 Level Req → 展开子表单: Min Level (number input)
   - 选 Points Threshold → 展开子表单: Point Type (dropdown) + Min Balance
   - 选 Manual Assignment → 无子表单，保存后通过 D15 手动添加成员
2. **Rank Order**: 数值越小层级越高（1 = 最高）。如果输入与已有层级冲突，已有层级自动 +1 后移
3. **Privileges 多选**: 至少选一项。勾选 "2x point multiplier" 时展开倍率输入框
4. **Modal 动画**: 打开: fade-in 200ms + scale 0.95→1; 关闭: fade-out 150ms; 关闭: × / overlay / Esc; 脏数据 → "Discard changes?" 确认
5. **"Create Tier" 按钮条件**:
   - 必填: Tier Name (1-40 chars) + Rank Order (1-99) + Qualification Condition (已选 + 子表单必填项) + Privileges (≥1 checked)
   - Disabled: 任何必填为空 → opacity 0.5 + tooltip "Fill in all required fields"
   - Loading: button spinner + disabled
6. **Edit mode — qualification.mode 锁定视觉**:
   - Active tier: Qualification Condition dropdown 显示为 **readonly text** + lock icon (16px, `#94A3B8`)
   - 子表单 (如 Token Gate 的合约/阈值): **可编辑** (仅 mode 锁定, 阈值可调)
   - Tooltip on lock: "Qualification mode cannot be changed for active tiers"
7. **Token Gate 子表单**:
   - Contract Address: select dropdown, 数据源 `GET /api/wl/contracts?status=verified` → 已验证合约列表; 0 contracts → dropdown 显示 "No contracts registered →" 链接到 B51
   - Min Balance: number input (> 0)
   - Chain: 自动从选择的 contract 填入 (readonly)
8. **Level Req 子表单**: Min Level: number input, 数据源 `GET /api/community/levels` → 获取最大 level 作为验证上限
9. **Duration 选择**: 位于 Privileges checkboxes 下方; radio: "Permanent" (默认) / "Time-limited" → 展开 days input (1-365)
10. **Tier Icon selector**: 点击 icon 区域 → 弹出 icon picker popover; 约 20 个 Material Symbols preset (diamond/star/verified/shield/...); 自定义颜色: 4 个 preset swatches (purple/amber/blue/red) + custom hex input
11. **Rank Order 冲突自动解决**: 输入与已有 tier 冲突 → 保存时已有 tiers 自动 +1; toast 通知 "Rank order adjusted for existing tiers"
12. **保存**:
   - 创建: POST `/api/wl/privileges` → 默认 status=draft
   - 编辑: PUT `/api/wl/privileges/:id`
   - 保存成功 → 关闭 Modal → 刷新 B53 列表

#### D15 — Privilege Members Panel (Modal)

**设计稿**: Node `zNH8l` | 宽度 480px

**Modal 结构**:
```
Title Bar: "VIP Tier Members" (动态层级名) + "24 members" + × close
Body:
├── Search: "Search by address..." (wallet address search)
├── Action Bar: [+ Add Member] (purple) + [↑ Bulk Import] (secondary)
├── Member List
│   ├── 0x7a3b...f291 — Joined Jan 5, 2026 — × remove
│   ├── 0x9c1e...a847 — Joined Feb 12, 2026 — × remove
│   ├── 0x2d5f...c103 — Joined Mar 1, 2026 — × remove
│   └── 0x6b8a...e459 — Joined Mar 3, 2026 — × remove
Footer: [↓ Export CSV]
```

**交互逻辑**:

1. **Add Member**: 弹出地址输入框，输入 wallet address → 验证格式 → 添加到列表
2. **Bulk Import**: 上传 CSV 文件（格式: 每行一个 wallet address），批量添加
3. **Remove (×)**: 确认弹窗 "Remove 0x7a3b...f291 from VIP tier?" → 确认后移除
4. **Search**: 按 wallet address 前缀搜索，实时过滤列表
5. **Export CSV**: 导出当前层级所有成员地址 + 加入时间

**适用场景**: 主要用于 Manual Assignment 模式。Status-Based / Achievement-Based 模式下，成员列表为系统自动管理（只读浏览 + 手动移除）。

**D15 页面初始化 (Init)**:
- API: `GET /api/wl/privileges/:id/members?page=1&limit=50` — 分页加载成员列表
- 加载态: Member List 区域 skeleton (5 行占位); title 中 "X members" 显示 spinner
- 空态 (0 members): "No members yet. Add members manually or wait for automatic qualification."
- Page size: 50/page; 分页: 底部 "Load more" 按钮 (非翻页)

**D15 "+ Add Member" 流程**:
- 点击 → inline input 展开 (slide-down): wallet address input (placeholder "0x...") + [Cancel] [Add] buttons
- 验证: `0x` prefix + 40 hex characters (EIP-55 checksum); 无效 → red border + "Invalid wallet address"
- 重复检查: 客户端检查列表中已有 → "This address is already a member"
- Add 按钮 click → `POST /api/wl/privileges/:id/members { walletAddress }` → 成功: 新成员追加到列表顶部 + title "X members" 数字 +1 + toast "Member added" + input 清空 (保持展开, 方便连续添加)
- 失败: toast error + input 保留

**D15 "Bulk Import" 流程**:
- 点击 → file picker (accept: .csv, .txt)
- 格式: 每行一个 wallet address; 支持有/无 header row
- Max file size: 1MB; Max addresses: 10,000 per import
- 上传后: parsing spinner → 显示 preview: "Found {n} addresses. {valid} valid, {invalid} invalid, {duplicate} duplicates."
- Invalid/duplicate 详情: 可展开查看 (前 20 条无效地址 + 原因)
- Confirm → `POST /api/wl/privileges/:id/members { addresses: [...] }` → 进度条 (0-100%) → 完成: toast "{n} members added" + 列表刷新
- 失败: toast error + 显示失败数量

**D15 "× remove" 每行**:
- 点击 → confirm popover (非 modal): "Remove this member?" [Cancel] [Remove]
- Remove → `DELETE /api/wl/privileges/:id/members/:wallet` → 行 slide-up 移除 (300ms) + title "X members" -1
- Loading: × icon 变为 spinner
- 失败: toast error + 行保留

**D15 "Export CSV"**:
- 点击 → `GET /api/wl/privileges/:id/members/export` → 浏览器下载 CSV
- 异步或同步: 同步 (< 10,000 条直接下载); > 10,000 → 后台生成 + email 通知
- Loading: button spinner

**D15 Title 即时更新**: Add/Remove → "X members" 即时 ±1 (optimistic)

**D15 Read-only mode**: Status-Based / Achievement-Based 模式下, "+ Add Member" 和 "Bulk Import" **隐藏**; "× remove" 仍可见 (管理员手动移除); 列表标题追加 "(Auto-managed)"

**D15 C端影响**:
- Manual add → C端权益变化**即时生效** (API 实时查询)
- Remove → C端权益**即时撤回** (下次 API check 返回无权益)
- Push notification: 由项目方通过 webhook 自行实现 (TaskOn 不推送 C端通知)

#### 按钮路由

| 按钮 | 位置 | 目标 | 说明 |
|------|------|------|------|
| "+ Create Privilege" | Header | → D14 (create) | 空表单 |
| Row click | Privileges Table | → D14 (edit) | 预填已有数据 |
| "Manage Members" (隐含) | Table row 操作区 | → D15 Members Panel | 查看/管理该层级成员 |
| "Most common →" | Mode A card | → D14 (preset: Status-Based) | 预选 qualification.mode |
| "Via Benefits Shop →" | Mode B card | → B31g Benefits Shop | 跳转 Shop 管理 |
| "Configure →" | Mode C card | → D14 (preset: Achievement-Based) | 预选 qualification.mode |
| Filter tabs | Table 上方 | (前端筛选) | All / Active / Draft |
| Row toggle | Table 行内 toggle | (API) `PUT /api/wl/privileges/:id` | active↔paused 状态切换 |
| "Manage Members" | Table 行 ⋮ 菜单 | → D15 Members Panel | 查看/管理该层级成员 |
| "Delete" | Table 行 ⋮ 菜单 | → confirm + `DELETE /api/wl/privileges/:id` | 删除层级 |

#### Privilege 状态流转

```
                  ┌─────────┐
  Create Tier     │  Draft  │   (未发布，可编辑所有字段)
  ────────────►   │         │
                  └────┬────┘
                       │ Activate
                       ▼
                  ┌─────────┐    Pause              ┌─────────┐
                  │ Active  │ ◄────────────────────► │ Paused  │
                  │         │    Resume              │         │
                  └────┬────┘                        └────┬────┘
                       │                                  │
                       │ (Active 后 qualification.mode    │
                       │  不可更改，仅可调整阈值/权益)      │
                       └──────────────────────────────────┘
```

- **Draft**: 可自由编辑所有字段
- **Active**: qualification.mode 锁定不可更改（防止已授权用户突然失去资格）。可调整: 阈值数值、权益列表、层级名称/图标/颜色
- **Paused**: 暂停新用户获得资格，已有用户权益保留。恢复后继续生效

**B53 页面初始化 (Init)**:
- 并行调用:
  1. `GET /api/wl/privileges` — Privileges Table (默认 status=all)
  2. `GET /api/wl/privileges/stats` — Stats Row (4 cards)
  3. `GET /api/wl/privileges/integration` — Integration Status section
- 加载态: 各 section 独立 skeleton (Stats / Table / Qualification Modes / Integration)
- API 失败: 各 section 独立 §2.3.4 错误态 + Retry

**B53 页面离开 (Destroy)**: 无状态需清理

**B53 Privilege row toggle (active↔paused)**:
- UI: 每行右侧 toggle switch (同 B52 Rule toggle 样式)
- 交互: 点击 → 无确认弹窗 → `PUT /api/wl/privileges/:id { status: 'active'|'paused' }` → 成功: toggle 动画 + toast "Privilege {activated|paused}"
- Draft → Active: 首次启用; qualification.mode 从此**锁定**不可更改 (见状态流转)
- Active → Paused: 已有 holders 权益保留; 新用户不再获得资格
- Stats 即时更新: Active Privileges ±1 (optimistic)

**B53 "Manage Members" 位置**:
- 每行 ⋮ (more_horiz) 菜单中: [Edit] [Manage Members] [Delete]
- "Manage Members" → 打开 D15 Members Panel (Modal/Side panel)
- Manual Assignment 模式: D15 内可 Add/Remove; Status-Based/Achievement-Based 模式: D15 只读浏览 + 手动移除

**B53 "Delete" 操作**:
- ⋮ 菜单 → "Delete" → confirm dialog: "Delete privilege tier '{name}'? Active holders will retain their privileges until expiration."
- Mode B 关联: "This will also remove the linked Benefits Shop item." (如有)
- 确认 → `DELETE /api/wl/privileges/:id` → 列表刷新 + toast "Privilege tier deleted"
- 已授予的权益: 不追溯撤回; time_limited 到期自然失效; permanent 保留直到项目方手动处理

**B53 D14 save 后 Table 刷新**: D14 关闭 → refetch `GET /api/wl/privileges`

**B53 "Total Value Distributed" 计算**: `GET /api/wl/privileges/stats` 返回, 后端聚合所有层级的 fee discount / gas rebate 等价值; 无法计算时显示 "—"

**B53 "API Connection: Enabled" 判定**: `GET /api/wl/privileges/integration` → 项目方是否配置了 API key (来自 B41) + 是否有 webhook (来自 B41); 全部满足 → "Enabled ✅"; 部分 → "Partial Setup" (amber); 无 → "Not Connected" (gray)

**B53 "Via Benefits Shop →" (Mode B)**: 跳转 B31g Benefits Shop; Shop 模块未配置时: 跳转 B31g (会显示 Shop 的空态或提示先在 Community 启用)

**B53 Loading skeleton**: Stats Row (4 skeleton) + Privileges Table (4 行 skeleton) + Modes (3 卡片 skeleton) + Integration (1 card skeleton)

**B53 Delete privilege tier — not in operations**: 已通过上方 ⋮ 菜单补充

**侧栏**: WL 子菜单展开, SMART REWARDS section → "Privileges" 高亮 active

#### API

| Endpoint | Method | 说明 | 请求体关键字段 |
|----------|--------|------|---------------|
| `/api/wl/privileges` | GET | 获取层级列表 | query: `?status=active&page=1` |
| `/api/wl/privileges` | POST | 创建层级 | `{ name, rankOrder, icon, color, qualification, privileges, duration }` |
| `/api/wl/privileges/:id` | GET | 获取单个层级详情 | — |
| `/api/wl/privileges/:id` | PUT | 更新层级 | 同 POST body + `{ status }` |
| `/api/wl/privileges/:id` | DELETE | 删除层级 | — |
| `/api/wl/privileges/:id/members` | GET | 获取层级成员列表 | query: `?search=0x&page=1&limit=50` |
| `/api/wl/privileges/:id/members` | POST | 添加成员 | `{ walletAddress }` 或 `{ addresses: [...] }` (批量) |
| `/api/wl/privileges/:id/members/:wallet` | DELETE | 移除成员 | — |
| `/api/wl/privileges/:id/members/export` | GET | 导出 CSV | — |
| `/api/wl/privileges/check` | GET | **项目方调用**: 查询用户权益 | query: `?wallet=0x...` → 返回权益列表 |
| `/api/wl/privileges/stats` | GET | 汇总统计 | — |
| `/api/wl/privileges/integration` | GET | 集成状态 | — |

---

## 13. Page Analytics

### 13.1 B43 — Page Analytics

**设计稿**: Node `69HPh` | URL: `/white-label/pages/:id/analytics`

#### 页面结构

```
Content Area
├── Breadcrumb: "← Back to Pages"
├── Header: "Page Analytics" + Date Range Picker + "Export" button
├── Stats Row (4 cards)
│   ├── Page Views: 12,845
│   ├── Unique Visitors: 3,421
│   ├── Widget Clicks: 1,847
│   └── Total Completions: 423
├── "Page Views Over Time" chart (line/area, D/W/M toggle)
├── "Top Pages" ranking table
│   ├── /rewards-hub — 14,812 (47.1%)
│   ├── /community-portal — 5,478 (12.4%)
│   └── ...
├── "Widget Interactions" table
│   ├── Leaderboard Widget — 1,847 clicks, 88% completion
│   ├── Quest List Widget — 1,923 clicks, 45% completion
│   └── Check-In Pop-up Widget — 889 views, 71% conversion
├── "Conversion Funnel" (3-step visual)
│   └── Page Views (12,845) → Interactions (5,768, 44.9%) → Completions (1,847, 32%)
```

#### 操作详情 (v2.0 新增, W-101 ~ W-105)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-101 | Date Range Picker | click | `GET /api/white-label/pages/:id/analytics?from={}&to={}` | 预设: 7d (默认) / 30d / 90d 快捷按钮 + "Custom" → calendar 双日期选择器; 选择后所有图表+Stats+表格数据刷新; 加载中: skeleton overlay on charts + stats 数字闪烁; custom 最大范围 365 天 | 无数据范围 → charts 显示 "No data for this period" |
| W-102 | Chart D/W/M toggle | click | 同上 (追加 `granularity=day|week|month`) | "Page Views Over Time" chart 粒度切换: Day (每日柱状) / Week (每周柱状) / Month (每月柱状); hover tooltip: "{date}: {views} page views ({delta}% vs previous period)"; chart 类型: area chart (filled), 紫色主色 `#9B7EE0` | — |
| W-103 | "Top Pages" 行点击 | click row | — | 点击某页面行 → 下方 "Widget Interactions" 表格筛选为该页面的 widget 数据; 选中行高亮 (bg `#1A1033`); 再次点击取消筛选 (显示全部); 行数据: 页面 slug + views (number) + share (%) + trend arrow (↑↓) | — |
| W-104 | "Export" | click | `GET /api/white-label/pages/:id/analytics/export?format=csv&from={}&to={}` | button spinner → 后台生成 CSV → 浏览器触发下载; CSV 包含: 日期, 页面, page views, unique visitors, widget clicks, completions; 文件名: `wl-analytics-{slug}-{date}.csv` | 生成失败 → toast "Export failed. Please try again." |
| W-105 | Conversion Funnel | hover + click | — | 3-step 可视化漏斗 (左→右递减): Page Views → Interactions → Completions; 每步 hover tooltip: "{count} ({percent}% of total, {dropoff}% drop from previous step)"; 点击步骤 → smooth scroll 到对应详情区 (Views→chart, Interactions→widget table, Completions→stats card) | — |

**B43 页面初始化 (Init)**:
- API: `GET /api/white-label/pages/:id/analytics?from={7d_ago}&to={today}` — 默认 7 天数据
- 路由守卫: `:id` 无效 → API 返回 404 → 显示 "Page not found" + "← Back to Pages" 链接 (→ B25); toast "Page not found"
- 加载态: Stats Row (4 skeleton) + Chart (chart skeleton) + Top Pages (table skeleton) + Widget Interactions (table skeleton) + Funnel (3 step skeleton)
- API 失败: §2.3.4 错误态 + Retry

**B43 页面离开 (Destroy)**: 无状态需清理

**B43 Date range 变更**: 全页 per-section loading (各图表独立 skeleton overlay); 非全局阻塞

**B43 Top Pages row click → Widget Interactions (W-103)**: 客户端筛选 (已有数据中 filter by page slug); 非 API re-query

**B43 Analytics API 失败**: skeleton → 错误态 + Retry

**B43 Unpublished page**: 显示历史数据 (如有); 标题旁 amber badge "Unpublished"

**B43 Breadcrumb**: "← Back to Pages" → B25 (始终); 不 context-dependent

**B43 Date range Custom min range**: 1 day minimum

**B43 D/W/M toggle 选中样式**: 选中: filled purple bg `#1A1033` + text `#9B7EE0`; 未选: outline border `#1E293B` + text `#94A3B8`

**侧栏**: WL 子菜单展开, "Pages" 高亮 active

---

## 14. Dev Kit Page

### 14.1 B48 — Dev Kit Page

**设计稿**: Node `3jDeL`

#### 页面概述
- **URL**: `taskon.xyz/devkit/{project_id}` (独立页面，无需登录)
- **用户**: 项目方的开发人员（由市场人员发送链接）
- **功能**: 集成指南 + 代码 + SSO 配置 + 验证

#### 页面结构

```
Standalone Page (无 Sidebar，独立暗色主题)
├── Header: "TaskOn / Dev Kit" breadcrumb
├── Title: "{Project Name} — Integration Guide"
├── Subtitle: "Everything you need to integrate TaskOn..."
├── ── ① Install Widget SDK ──
│   ├── Package manager tabs: npm / yarn / CDN
│   └── Code block (copy-ready)
├── ── ② Configure SSO ──
│   ├── JWT Provider selector
│   │   ├── Option A: Wallet Authentication
│   │   └── Option B: OAuth2 / Custom JWT
│   └── Code block (provider config)
├── ── ③ Embed Your Widgets ──
│   ├── Description: "Each widget below..."
│   ├── Expandable widget list:
│   │   ├── ▸ Leaderboard Widget (code + preview)
│   │   ├── ▸ Task List Widget
│   │   └── ▸ User Center Widget
├── ── Ready to Verify? ──
│   ├── Info text: "Once deployed, click verify..."
│   ├── "Verify Integration" button (紫色)
│   └── Status: "Waiting for first API ping from your domain..."
├── Footer: "Questions? Docs / API Reference / Support"
```

#### API

| Endpoint | Method | 说明 |
|----------|--------|------|
| `GET /api/devkit/{project_id}` | GET | 获取项目配置 (widgets, SSO, API key) |
| `POST /api/devkit/{project_id}/verify` | POST | 触发集成验证 |

#### 操作详情 (v2.0 新增, W-106 ~ W-112)

| # | 操作 | 触发方式 | API 调用 | 即时 UI 变化 | 错误处理 |
|---|------|---------|---------|------------|---------|
| W-106 | Dev Kit URL 生成 | WL publish | — | URL: `taskon.xyz/devkit/{project_id}`; project_id 在 WL 首次 publish 时生成 (UUID v4); URL 公开可访问, **无需 TaskOn 登录**; Dev Kit 页面独立暗色主题 (bg `#0A0F1A`), 无 sidebar; 市场人员通过 B15 "Send Dev Kit" 获取链接发给开发者 | — |
| W-107 | Package manager tabs | click tab | — | 3 个 tab: **npm** (默认) / **yarn** / **CDN**; 切换显示对应安装命令: npm → `npm install @taskon/widget-sdk`; yarn → `yarn add @taskon/widget-sdk`; CDN → `<script src="https://cdn.taskon.io/sdk/v1/widget.min.js"></script>`; 每个 code block 右上角 "Copy" 按钮 → clipboard → toast "Copied!" | — |
| W-108 | SSO Provider selector | click radio | — | 2 个选项: **Wallet Authentication** (推荐, 零后端) → 显示 wallet SDK 集成代码 (`taskon.auth.connectWallet()`) + 支持钱包列表 (MetaMask/WalletConnect/Coinbase); **OAuth2 / Custom JWT** → 显示 JWT 配置代码 (JWT Secret input + Redirect URI + 代码 snippet `taskon.auth.jwt({ token })`) ; 代码中 project_id 已预填 | — |
| W-109 | Widget 展开 | click widget row | — | "Embed Your Widgets" 区: 每个已配置 widget 为可展开行 (accordion); 展开 (slide-down 200ms): 嵌入代码 (project_id + widget_id 预填, Copy 按钮) + mini preview (widget 缩略图 200×150); 折叠态: widget name + module type icon + "▸"; 展开态: "▾" + 代码 + preview; 同时可多个展开 | — |
| W-110 | "Verify Integration" | click | `POST /api/devkit/{project_id}/verify` | button → spinner ("Verifying...") → API 检查是否收到来自项目域名的 API ping (检查最近 5 分钟): 有 ping → 成功 (见 W-111); 无 ping → "No API ping detected yet. Make sure your widget code is deployed and your domain is correctly configured." (amber text) + 按钮恢复为 "Try Again" | API 超时 (10s) → "Verification timed out. Please try again." |
| W-111 | 验证成功 | verify response | — | 按钮变为 "✓ Integration Verified" (green, `#16A34A`, disabled 不可再点); 下方显示: "Integration verified at {ISO timestamp} from {detected_domain}" (green text); confetti animation (2s, subtle); "Ready to Verify?" section 标题变为 "✅ Integration Complete"; 状态持久化 (刷新页面仍显示已验证) | — |
| W-112 | Dev Kit 过期 | render | `GET /api/devkit/{project_id}` → 404 | 项目删除或 Dev Kit URL 被撤销时: 整页显示友好错误: 大图标 (link_off, 64px, gray) + "This Dev Kit link is no longer valid" (24px) + "The project owner may have revoked access or deleted this project." (16px gray) + "Contact your project administrator for a new link." + TaskOn logo footer | — |

**B48 页面初始化 (Init)**:
- API: `GET /api/devkit/{project_id}` — 获取项目配置 (widgets, SSO, API key, verification status)
- 网络超时: 10s timeout → 显示 "Unable to load. Please check your connection and try again." + "Retry" 按钮
- 加载态: 全页面 skeleton (标题 + 3 个 section + verify 区域)
- 独立布局: 无 sidebar (standalone page); 暗色主题 bg `#0A0F1A`; 自有 header "TaskOn / Dev Kit"

**B48 验证状态持久化 (W-111)**:
- 验证成功后: **服务端存储** `{ verified: true, verifiedAt: ISO_timestamp, verifiedDomain: 'app.project.io' }`
- 刷新页面: `GET /api/devkit/{project_id}` 返回含 `verified: true` → 直接渲染 W-111 成功态 (绿色 ✓ Integration Verified); 无需重新验证
- 项目方更新 widget/SSO 配置后: Dev Kit 页面**自动反映最新配置** (非 snapshot); verified 状态保留 (配置变更不影响验证状态)
- 重置验证: 项目方在 B15 重新生成 Dev Kit → 旧 Dev Kit URL 失效 (W-112); 新 URL 验证状态重置

**B48 Widget accordion**: 同时可多个展开 (无 max 限制); 展开/折叠动画: slide-down/up 200ms

**B48 0 widgets**: Widget section 显示 "No widgets configured yet. The project owner needs to set up widgets in the Widget Library."

**B48 Footer links**: Docs → `https://docs.taskon.xyz` (external); API Reference → `https://docs.taskon.xyz/api` (external); Support → `mailto:support@taskon.xyz`

**B48 Initial loading state**: 全页面 skeleton

---

## 15. 侧栏架构

### 15.1 WL 展开子菜单 (所有 WL 页面通用)

```
▾ White Label (purple highlight)
  Overview          → B14/B15/B16 (根据状态)
  Widgets           → B20/B22 (根据状态)
  Pages             → B23/B25 (根据状态)
  SMART REWARDS (section header, non-clickable)
  ├── Contracts     → B51
  ├── Rule Builder  → B52
  └── Privileges    → B53
```

Active item 样式: purple bg `#1A1033` + text `#9B7EE0` + fontWeight 600

**注意**: Smart Rewards 是 section header (fontSize 10, fontWeight 600, `#94A3B8`, letterSpacing 1)，下方 3 个子项可点击。

### 15.2 侧栏 Active 状态映射 (完整)

| 页面 | 侧栏 Active 子项 | 说明 |
|------|-----------------|------|
| B14 / B15 / B16 (Hub) | **Overview** | Hub 3 个状态均对应 Overview |
| B37 / B17 / B57-B60 / B38 / B56 (Wizard) | **Overview** | 向导从 Overview 发起 |
| B18 (Domain Setup) | **Overview** | Domain 从 Overview 进入, 不在子菜单中 |
| B19 (Deployment Settings) | **Overview** | 同上 |
| B42 (Iframe Embed) | **Overview** | Iframe 从 B19 进入, 归属 Overview |
| B20 / B21 / B22 (Widget Library) | **Widgets** | — |
| B23 / B24 / B25 (Page Builder) | **Pages** | — |
| B40 (Brand Settings) | **Overview** | Brand 从 Overview 进入 |
| B41 (SDK & API) | **Overview** | SDK 从 Overview 进入 |
| B26 (Integration Center) | **Overview** | Integration 从 Overview 进入 |
| B44 (Integration Detail) | **Overview** | 从 B26 进入, 无独立子项 |
| B51 (Contract Registry) | **Contracts** (Smart Rewards) | — |
| B52 (Rule Builder) | **Rule Builder** (Smart Rewards) | — |
| B53 (Privilege Manager) | **Privileges** (Smart Rewards) | — |
| B43 (Page Analytics) | **Pages** | 从 Pages 进入 |
| B48 (Dev Kit) | N/A | 独立页面, 无 sidebar |

---

## 16. API 接口汇总

### 16.1 White Label 核心 API

| Endpoint | Method | 用途 | 页面 | Cache |
|----------|--------|------|------|-------|
| `/api/white-label/status` | GET | Hub 状态 | B15, B16 | 60s |
| `/api/white-label/drafts` | POST | 保存草稿 | B37, B17, B38, B56 | N/A |
| `/api/white-label/publish` | POST | 发布 WL | B56 | N/A |
| `/api/white-label/verify-domain` | POST | DNS 验证 | B17, B18 | 轮询 |
| `/api/white-label/widgets` | GET | Widget 列表 | B20-B22 | 60s |
| `/api/white-label/widgets/:id` | GET/PUT/DELETE | Widget CRUD | B21 | 60s |
| `/api/white-label/embed-code` | GET | Embed 代码 | B17, B21, B42 | 0 |
| `/api/white-label/pages` | GET | Page 列表 | B23-B25 | 60s |
| `/api/white-label/pages/:id` | GET/PUT/DELETE | Page CRUD | B24 | 60s |
| `/api/white-label/pages/:id/analytics` | GET | Page 分析 | B43 | 60s |
| `/api/white-label/brand` | GET/PUT | 品牌设置 | B40 | 60s |
| `/api/white-label/sdk` | GET | SDK 配置 | B41 | 60s |
| `/api/white-label/sdk/keys` | POST/DELETE | API Key 管理 | B41 | N/A |
| `/api/white-label/integrations` | GET | 集成列表 | B26 | 60s |
| `/api/white-label/integrations/:type` | GET/PUT | 集成配置 | B44 | 60s |
| `/api/wl/integrations/:type/connect` | POST | 集成连接 (OAuth callback / API Key) | B44 | N/A |
| `/api/wl/integrations/:type/test` | POST | 测试连接 | B44 | N/A |
| `/api/white-label/analytics/summary` | GET | 项目级分析汇总 | B16 | 60s |
| `/api/white-label/domain` | GET | 当前域名配置 | B18 | 60s |
| `/api/white-label/domain` | PUT | 更新域名 | B18 | N/A |
| `/api/white-label/onboarding` | GET | Checklist 状态 | B15 | 60s |
| `/api/white-label/embed/iframe` | PUT | 保存 Iframe 配置 | B42 | N/A |
| `/api/white-label/embed/iframe/test-sso` | POST | 测试 SSO | B42 | N/A |
| `/api/white-label/readiness/checks` | GET | D20 检查项 | D20 | 300s |
| `/api/white-label/pages/templates/:id` | GET | 页面模板数据 | B24 | 60s |
| `/api/wl/contracts` | GET/POST | 合约注册 | B51 | 60s |
| `/api/wl/contracts/:id/verify` | POST | 合约验证 | B51 | N/A |
| `/api/wl/rules` | GET/POST/PUT | 活动规则 | B52 | 30s |
| `/api/wl/privileges` | GET/POST/PUT | 权限层级 | B53 | 60s |
| `/api/devkit/:project_id` | GET | Dev Kit 配置 | B48 | 60s |
| `/api/devkit/:project_id/verify` | POST | 集成验证 | B48 | N/A |
| `/api/promo-kit/generate` | POST | Promo Kit | B15 | N/A |
| `/api/white-label/readiness` | GET | 发布前就绪检查 | B56 | 0 |
| `/api/white-label/preview` | GET | 部署预览数据 | B56 | 0 |
| `/api/white-label/embed/iframe` | GET | Iframe 配置 | B42 | 60s |
| `/api/white-label/deployments` | GET | 部署列表 | B16 | 60s |
| `/api/white-label/deployments/history` | GET | 发布历史 | B16 | 60s |
| `/api/white-label/sdk/keys/regenerate` | POST | 重新生成 API Key | B41 | N/A |
| `/api/white-label/sdk/webhooks` | GET/POST | Webhook 管理 | B41 | 60s |
| `/api/white-label/sdk/webhooks/:id/test` | POST | 测试 Webhook | B41 | N/A |
| `/api/white-label/brand/logo` | POST | Logo 上传 | B40 | N/A |
| `/api/white-label/pages/check-slug` | GET | URL slug 唯一性检查 | B24 | 0 |
| `/api/white-label/pages/:id/analytics/export` | GET | 分析数据导出 | B43 | N/A |
| `/api/white-label/integrations/twitter/auth` | GET | Twitter OAuth | B26 | N/A |
| `/api/white-label/integrations/analytics` | POST | GA4 集成 | B26 | N/A |
| `/api/wl/contracts/:id` | GET/PUT/DELETE | 合约 CRUD | B51 | 60s |
| `/api/wl/contracts/:id/events` | GET | 合约事件历史 | B51 | 30s |
| `/api/wl/contracts/check` | GET | 合约重复检查 | B51 (D12) | 0 |
| `/api/wl/contracts/stats` | GET | 合约汇总统计 | B51 | 60s |
| `/api/community/modules/status` | GET | 模块配置状态 | B20 | 60s |
| `/api/devkit/:project_id/generate` | POST | 生成 Dev Kit | B56, B15 | N/A |
| `/api/devkit/:project_id/send-email` | POST | 邮件发送 Dev Kit | B15 | N/A |

### 16.2 WebSocket 端点

| Endpoint | 用途 | 页面 |
|----------|------|------|
| `/ws/wl/integration-ping` | 首次 API ping 检测 | B15 Checklist Step 4 |
| `/ws/wl/first-interaction` | 首个用户交互检测 | B15 Checklist Step 6 |

---

## 17. 状态路由策略

### 17.1 Hub 状态判断

```typescript
function getWLHubState(wl: WhiteLabel): PageCode {
  if (!wl || wl.configuredTools === 0) return 'B14'; // Empty
  if (
    wl.configuredTools >= 5 &&
    wl.monthlyImpressions >= 1000 &&
    wl.deployments.length >= 1
  ) return 'B16'; // Management
  return 'B15'; // Active (1-4 tools, or 5+ but low traffic/no deployments)
}
```

### 17.2 Widget Library 状态

| 条件 | 页面 |
|------|------|
| 0 widgets created | B20 Empty |
| 1+ widgets | B22 Active |

### 17.3 Page Builder 状态

| 条件 | 页面 |
|------|------|
| 0 pages created | B23 Empty |
| 1+ pages | B25 Active |

### 17.4 Wizard Step 2 路由

| Step 1 选择 | Step 2 目标 |
|------------|-----------|
| Embed (default) | B17 Widget Config (`CXzmy`) |
| Embed → Iframe variant | B58 Iframe Config (`Kr5W5`) |
| Embed → PB variant | B59 PB Config (`XHwzp`) |
| Domain | B57 Domain Config (`YGODW`) |
| SDK | B60 SDK Config (`eNFmU`) |

---

## 18. D19 Promo Kit Generator — WL (v2.0 新增)

**设计稿**: Node `2qNbJ` | 宽度 640px

> WL 版 Promo Kit Generator 与 Community 版 (见 `req_community.md` D19) 共享相同的 Modal 结构和交互模式, 但内容生成基于 WL 品牌和部署信息。

#### Modal 结构

```
Title Bar: "Promo Kit" + × close
Body:
├── Platform selector: Twitter (default) / Discord / Telegram
├── Banner Preview (AI 生成, WL brand logo + colors)
│   └── "Regenerate Banner" button
├── Social Copy (AI 生成, 可编辑 textarea)
│   ├── Preview text area (pre-filled with AI content)
│   └── "Regenerate Text" link
├── Action Bar:
│   ├── "Share on {Platform}" button (purple)
│   ├── "Copy Text" button (secondary)
│   └── "Download Banner" button (secondary)
```

#### AI Generation

**Init API**: Modal 打开时调用 `POST /api/promo-kit/generate` with `{ type: 'white_label', projectId, platform: 'twitter' }`
- 返回: `{ bannerUrl, socialCopy, hashtags }`
- **生成时间**: 约 3-8 秒
- **Loading state**: Banner 区域显示 shimmer skeleton + "Generating your promo materials..." 文案; 文本区域显示 typing 动画 (逐字出现效果)

**AI Generation 失败**:
- Banner 失败: 显示 fallback banner (纯色 + logo + 项目名); amber 提示 "AI banner generation failed. Using default template." + "Retry" 按钮
- Text 失败: 显示 fallback template text (项目名 + 通用 CTA); amber 提示 + "Retry"
- 双重失败: 全部使用 fallback; 功能仍可用 (Copy/Share/Download 均正常)

**Regenerate**:
- "Regenerate Banner": `POST /api/promo-kit/generate { type: 'banner_only' }` → banner 区域 spinner → 新 banner 替换
- "Regenerate Text": `POST /api/promo-kit/generate { type: 'text_only' }` → textarea spinner → 新文案替换
- Rate limit: 每种类型 10 次/小时; 超限 → toast "Generation limit reached. Try again later." + button disabled
- **用户已编辑时 Regenerate**: 显示确认 "Regenerate will replace your edits. Continue?" [Cancel] [Regenerate]

**Platform switch**:
- 切换 platform → **自动 regenerate** text (不同 platform 的文案风格/长度/格式不同)
- 如用户已编辑过当前 platform 文案 → 切换前确认 "Switching platform will regenerate text. Your edits will be lost." [Stay] [Switch]
- Banner **不重新生成** (跨 platform 通用)

**Action Buttons**:
- "Share on Twitter": 打开 Twitter compose 窗口 (popup), pre-filled URL: `https://twitter.com/intent/tweet?text={encoded_text}&url={project_url}`
- "Copy Text": clipboard → toast "Copied!"
- "Download Banner": 下载 PNG/JPG; 分辨率: 1200×630 (Twitter card size); 文件大小: 约 200-500KB

**WL-specific 差异 (vs Community)**:
- Banner 使用 WL brand colors + logo (非 Community)
- 文案内容: 强调 "custom branded experience" / "your own domain" 等 WL 特色
- 无 logo/brand configured: banner 使用 TaskOn 默认 + amber 提示 "Set up your brand in Brand Settings for customized banners"

---

## 19. D20 Publish Readiness Check — WL (v2.0 新增)

**设计稿**: Node `fY99y` | 宽度 480px

> D20 是所有 Publish/Deploy/Launch 操作前的统一检查 Modal。WL 与 Community 共享相同的 2-item checklist 结构。
> 权威规范参见 `req_community.md` §7.2 / §7.3, 本节补充 WL 特定内容。

#### Modal 结构

```
Title Bar: "Publish Readiness Check" + × close
Body:
├── Checklist:
│   ├── ① Subscription Status
│   │   ├── ✅ "White Label plan active" (green)
│   │   └── OR ❌ "No active WL subscription" (red) + "Subscribe →" link
│   ├── ② Twitter Authorization
│   │   ├── ✅ "@handle authorized" (green)
│   │   └── OR ❌ "Twitter not connected" (red) + "Connect →" button
├── Result:
│   ├── All pass → "Ready to publish!" (green) + [Publish] button enabled
│   └── Any fail → "Fix issues above to continue" (amber) + [Publish Anyway] (if only subscription fails, see below)
Footer: [Cancel] [Publish] (purple, enabled only when all checks pass OR "Publish Anyway" conditions met)
```

#### Init API

- Modal 打开时调用: `GET /api/white-label/readiness/checks` → 返回 `{ subscription: { status, plan, expiresAt }, twitter: { connected, handle } }`
- **并行执行**: 两个 check 项同时进行, 各自显示 inline spinner
- **超时**: 每项 10s timeout; 超时 → 该项显示 amber "Check timed out" + "Retry" 按钮
- **缓存**: 5 分钟缓存; 5 分钟内重新打开 D20 → 直接显示缓存结果 (不重新检查); 缓存存储: sessionStorage; 两个 tab 各自独立

#### Check 项详情

**① Subscription Check**:
- 通过: WL plan active (Basic/Pro) + 未过期
- 失败: 无 WL 订阅 or 已过期 → "Subscribe →" 链接到 M07 Pricing page (新标签页, 不关闭 D20)
- Timeout: "Unable to verify subscription. Try again."

**② Twitter Authorization**:
- 通过: Twitter OAuth 已完成 + token 未过期
- 失败: 未连接 or token expired → "Connect →" 按钮触发 OAuth popup (同 B26 W-84); popup 完成后自动 re-check
- 目的: 确保项目方发布后能通过 Twitter 进行推广

#### "Publish Anyway" 逻辑

- 仅当 **subscription check timeout** (非明确失败) 时显示 "Publish Anyway" 按钮
- Subscription 明确失败 (无订阅): "Publish" disabled, 无 "Publish Anyway"
- Twitter 失败: "Publish" disabled, 必须连接 Twitter (产品决策: Twitter 是核心推广渠道)

#### Publish 按钮

- Enabled: 两项均 ✅
- 点击: `POST` 到调用方指定的 publish API (B56/B18/B24/B22 各自不同)
- Loading: button spinner
- 成功: Modal 自动关闭 → 调用方处理 post-publish 跳转 (toast + redirect)
- 失败: toast error + Modal 保持打开 + button 恢复

#### Modal 关闭

- × 按钮 / 点击 overlay → Modal 关闭, 不 publish (无 "Cancel publish?" 确认)
- Esc 键 → 同上

#### WL D20 触发点 (5 个)

| 页面 | 触发按钮 | Publish API |
|------|---------|------------|
| B56 Wizard Step 4 | "Publish White Label" | `POST /api/white-label/publish` |
| B18 Domain Setup | "Save" (域名修改) | `PUT /api/white-label/domain` |
| B24 Page Builder Editor | "Publish Page" | `POST /api/white-label/pages` |
| B22 Widget Library Active | "Deploy Widget" | `PUT /api/white-label/widgets/:id { status: 'deployed' }` |
| B15 Hub Active | "Announce" (Promo Kit) | N/A (D19 打开, 非 publish) |

> 注: B15 "Announce" 步骤打开 D19 Promo Kit, 不经过 D20。附录 A 中 "Announce / Promo Kit" 行指的是 D19, 非 D20。

---

## 附录 A: D20 Publish Touchpoints (WL)

| 页面 | 触发按钮 |
|------|---------|
| B56 Wizard Step 4 | "Publish White Label" |
| B18 Domain Setup | "Edit Domain Settings" |
| B24 Page Builder Editor | "Publish Page" |
| B22 Widget Library Active | "Deploy Widget" |
| B15 Hub Active (checklist) | "Announce" / Promo Kit |

---

## 附录 B: 设计稿 Node ID 索引

| 页面 | Node ID |
|------|---------|
| B14 WL Hub (Empty) | `Ir6Tq` |
| B15 WL Hub (Active) | `BnkYW` |
| B16 WL Hub (Management) | `UPAfV` |
| B37 Wizard Step 1 Path | `NNwid` |
| B17 Wizard Step 2 Widget | `CXzmy` |
| B57 Wizard Step 2 Domain | `YGODW` |
| B58 Wizard Step 2 Iframe | `Kr5W5` |
| B59 Wizard Step 2 PB | `XHwzp` |
| B60 Wizard Step 2 SDK | `eNFmU` |
| B38 Wizard Step 3 Brand | `5nCtO` |
| B56 Wizard Step 4 Preview | `WsH2y` |
| B18 Domain Setup | `5bmH9` |
| B19 Deployment Settings | `RgCVQ` |
| B42 Iframe Embed | `ByGS0` |
| B20 Widget Library (Empty) | `2sSsA` |
| B21 Widget Config | `n4pJK` |
| B22 Widget Library (Active) | `S432k` |
| B23 Page Builder (Empty) | `DRYwN` |
| B24 Page Builder (Editor) | `sGDcq` |
| B25 Page Builder (Active) | `J08v5` |
| B26 Integration Center | `Abs1E` |
| B44 Integration Detail | N/A (动态路由页, 无独立设计稿 Node) |
| B40 Brand Settings | `Cx3LH` |
| B41 SDK & API | `lQxT5` |
| B43 Page Analytics | `69HPh` |
| B51 Contract Registry | `OKEqS` |
| B52 Rule Builder | `4aAo7` |
| B53 Privilege Manager | `5xwYN` |
| B48 Dev Kit Page | `3jDeL` |
| B19v Embed Options (Neutral) | `Rwq2K` |
| D12 Contract Register Form | `NcxsI` |
| D13 Activity Rule Editor | `IJZ0E` |
| D14 Privilege Tier Editor | `FypcB` |
| D15 Privilege Members Panel | `zNH8l` |
| D19 Promo Kit Generator | `2qNbJ` |
| D20 Publish Readiness Check | `fY99y` |

---

## 附录 C: 级联效果 (CASCADE) — v2.0 新增

### C.1 Brand 变更级联

| 操作 | 影响范围 | 生效时机 | 说明 |
|------|---------|---------|------|
| B40 "Save Changes" | 所有已部署 Widget + Page | 下次用户访问 (CDN TTL 5min) | Logo/Color/Font 变更通过 API 动态加载, 无需 re-deploy |
| B38 Wizard Brand | 初始配置 | 发布时生效 | Wizard 中的 Brand 设置写入 B40 配置 |

### C.2 Community → WL 级联

| Community 操作 | WL 影响 | 说明 |
|---------------|---------|------|
| 启用新模块 (B34) | B20 Widget Library 新增可用 module card (green) | Widget Library 列表来源 = Community 已配置模块 |
| 禁用模块 (B34) | B20 对应 module card 变为 amber "Set Up in Community →"; **已创建的 Widget 不自动删除**, 状态变 "Source module disabled" (gray) | 保护已部署 widget 不丢失 |
| 删除任务/积分类型 (B31) | Rule Builder (B52) 中引用该任务的 Rule 状态变 "Error: Source deleted" | 显示 amber 警告条, 需手动修复 |
| Points 类型变更 (B31a) | Privilege Manager (B53) 中 Points Threshold 条件可能失效 | 自动检查并 flag 不一致 |

### C.3 Contract → Rule → Privilege 级联

| 操作 | 级联 | 说明 |
|------|------|------|
| B51 删除合约 | B52 引用该合约的 Rules 自动 disable (status=error, reason="Contract deleted") | 用户在删除确认时看到警告 |
| B52 Rule 触发 award_points | B31a Points balance 更新 → Leaderboard 实时排名更新 → B53 Points Threshold 类 Privilege 资格实时变化 | 全自动, 无需手动操作 |
| B52 Rule 触发 grant_badge | B31i Badges 列表更新 → B53 Achievement-Based Privilege 自动授予 | 一次性, 不可撤销 |
| B53 创建 Mode B Privilege | B31g Benefits Shop 自动创建关联商品 | 商品价格/库存在 B31g 管理, 权益在 B53 管理 |
| B53 删除 Privilege Tier | 已持有用户权益在到期前保留; Mode B 关联商品自动下架 | 不追溯已发放权益 |

### C.4 Publish 级联

| 操作 | 级联 | 说明 |
|------|------|------|
| B56 "Publish WL" | Dev Kit 自动生成 (W-43); B15 Checklist 自动标记前 3 步为 ✅; Widget/Page/Domain 配置锁定为 v1 snapshot | 首次发布 |
| B24 "Publish Page" (via D20) | B25 Page 列表新增; B43 Analytics 开始收集数据; B16 "Published Pages" stats +1 | — |
| B22 "Deploy Widget" (via D20) | Widget status → "Deployed" (blue); embed code 激活; B16 "Widgets Created" stats +1 | — |

### C.5 Subscription 级联

| 订阅状态 | WL 影响 | 说明 |
|---------|---------|------|
| WL 订阅过期 | 所有已部署 Widget/Page 显示 "Powered by TaskOn" watermark (不可移除); Brand 自定义部分回退为 TaskOn 默认; SDK API 降级为 readonly | 不删除数据, 续费后恢复 |
| WL 订阅降级 (从 Pro→Basic) | 超出 Basic 限额的 Page/Widget 变为 Draft (不自动删除); 需手动选择保留哪些 | 降级前 14 天邮件通知 |

---

## 附录 D: 边界条件 (EDGE) — v2.0 新增

### D.1 数据上限

| 实体 | 限制 | 超限行为 |
|------|------|---------|
| Widgets per project | 20 (Basic) / 50 (Pro) | "+ Add Widget" disabled + tooltip "Upgrade to Pro for more widgets" |
| Pages per project | 5 (Basic) / 20 (Pro) | 同上 |
| Contracts per project | 10 (Basic) / 50 (Pro) | "+ Register Contract" disabled |
| Rules per project | 20 (Basic) / 100 (Pro) | "+ Create Rule" disabled |
| Privilege Tiers | 10 per project | "+ Create Privilege" disabled |
| Webhooks | 5 per project | "+ Add Endpoint" disabled |
| ABI JSON size | max 500KB | 上传/粘贴超限 → "ABI file too large (max 500KB)" |
| Custom CSS | max 10KB | 编辑器字符计数 + 超限禁止保存 |
| Dev Kit concurrent | 1 per project | 重新生成会覆盖旧链接 |

### D.2 权限与角色

| 操作 | 最低角色 | 未授权行为 |
|------|---------|-----------|
| Publish WL (B56) | Admin | 按钮 disabled + "Only admins can publish" tooltip |
| Regenerate API Key (B41) | Admin | 按钮 disabled |
| Register Contract (B51) | Admin / Editor | Member 无权限 |
| Create/Edit Rule (B52) | Admin / Editor | Member 只读浏览 |
| Manage Privilege Members (D15) | Admin | Editor 只读浏览 |
| Brand Settings (B40) | Admin / Editor | Member 无权限 |
| Dev Kit Send Email (B15) | Admin / Editor | Member 无权限 |

### D.3 网络与恢复

| 场景 | 处理 |
|------|------|
| 页面加载 API 失败 | Skeleton → 3s 后显示 "Unable to load. Check your connection." + "Retry" 按钮 |
| 表单提交失败 | Toast error + 表单数据保留 (不清空) + 按钮恢复 enabled |
| WebSocket 断连 (B15) | 自动重连 3 次 (1s/3s/10s); 失败 → 降级为 HTTP polling (30s 间隔) |
| DNS 验证 RPC 超时 | "DNS verification temporarily unavailable. Will retry automatically." + 60s 后重试 |
| Widget iframe 加载失败 | iframe fallback: 显示 "Widget temporarily unavailable" + 项目 logo |
| Dev Kit 页面 404 | 见 W-112: 友好错误页, 不显示 TaskOn 内部错误信息 |

### D.4 并发与冲突

| 场景 | 处理 |
|------|------|
| 多人同时编辑 B40 Brand | 后保存者覆盖 (last-write-wins); 无实时协作; 保存前不检查版本冲突 |
| 同时 Publish B56 | 第一个请求成功; 第二个返回 409 "Already published" → toast + redirect B15 |
| Widget 被删除时 Page 引用 | Page 中该 widget block 显示 "Widget removed" 占位符 (gray box); 不影响 Page 整体 |
| Rule 引用的合约被删除 | Rule status → error (见 CASCADE C.3); Rule 不自动删除, 等待用户修复或手动删除 |
| Privilege Mode B 关联商品被直接从 Shop 删除 | Privilege Tier 显示 amber 警告 "Linked shop item deleted. Users can no longer redeem." |

### D.5 空数据特殊处理

| 页面 | 空数据情况 | 显示 |
|------|-----------|------|
| B43 Page Analytics | 新页面无访问数据 | 所有 stats = 0; chart = empty state "No data yet. Share your page to start tracking."; 不显示 funnel |
| B51 Contract Registry | 无合约注册 | Table 区域: empty illustration + "Register your first smart contract to start tracking on-chain activity" + "+ Register Contract" CTA |
| B52 Rule Builder | 无 Rules 创建 | Active Rules Table: empty + "Create your first rule to automate rewards" + "+ Create Rule" CTA; 仍显示 Presets 和 Anti-Sybil 区域 |
| B53 Privilege Manager | 无 Privileges 创建 | Table: empty + "Create your first privilege tier" + CTA; 仍显示 Qualification Modes 和 Integration Status |
| B22 Widget Active | 所有 widget 被删除 | 回退到 B20 (Empty) 状态 |
| B25 Page Active | 所有 page 被删除 | 回退到 B23 (Empty) 状态 |
