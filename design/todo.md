# Design TODO

> Active design tasks for `design/pencil-new.pen`
> Previous work archived in `design/progress-archive.md`
> Last updated: 2026-03-06

---

## Phase: Publish Readiness Check — COMPLETE

> **Objective**: 所有 Community 和 White Label 中面向用户发布的内容，在 Publish 之前必须通过 2 项公共检测。
> **Status**: 设计完成 — D20 `fY99y` (Row 29, x:1800)

### 检测项

| # | 检测项 | 目的 | 通过条件 | 失败处理 |
|---|--------|------|----------|----------|
| 1 | **订阅状态** | 确保项目有权使用平台 | 试用期内 OR 付费有效期内 | 跳转 Pricing / 续费页面 |
| 2 | **官方推特授权** | 防止项目冒充他人名义发布活动骗用户 | 已连接并验证项目官方 Twitter 账户 | 跳转 Integration / 授权流程 |

### UX 方案

- **触发时机**: 用户点击 "Publish" / "Go Live" / "Launch" 按钮时
- **展示形式**: 弹窗 Modal（新建 D20），2 行 checklist，每行显示检测结果（✓ 绿色通过 / ✗ 红色未通过）
- **通过**: 2 项都通过 → 自动关闭 Modal，执行发布
- **未通过**: 阻止发布，显示未通过项 + 解决方案链接（"Upgrade Plan →" / "Connect Twitter →"）
- **设计参考**: 类似 WL Wizard Step 4 的 readiness checklist，但作为独立 Modal

### 调用清单（Publish Touchpoints）

#### Community 产品

| 页面 | Node ID | 触发按钮 | 场景 |
|------|---------|----------|------|
| Community Wizard Step 4 — Preview & Publish | `7mVsZ` | "Publish Community" | 首次发布社区 |
| Community Hub — Guided (checklist) | `S1EIA` | "Share" / Promo Kit | 首次推广 |
| Sectors & Tasks Mgmt | `Wug7d` | "Publish Task" | 发布新任务 |
| TaskChain Mgmt | `lpdtp` | "Activate Chain" | 激活 TaskChain |
| DayChain Mgmt | `fLLVb` | "Activate Chain" | 激活 DayChain |
| LB Sprint Mgmt | `FO9JR` | "Launch Sprint" | 启动竞赛 |
| Milestone Mgmt | `WFdZQ` | "Activate Milestone" | 激活里程碑 |
| Benefits Shop Mgmt | `7yPWx` | "Publish Item" | 上架商品 |
| Lucky Wheel Mgmt | `sme5a` | "Activate Wheel" | 激活转盘 |
| Content Management | `lhR14` | "Publish" | 发布公告/内容 |

#### White Label 产品

| 页面 | Node ID | 触发按钮 | 场景 |
|------|---------|----------|------|
| WL Wizard Step 4 — Preview & Publish | `WsH2y` | "Publish WL" | 首次发布 WL |
| WL Domain Setup | `5bmH9` | "Edit Domain Settings" | 域名变更后重新生效 |
| WL Page Builder — Editor | `sGDcq` | "Publish Page" | 发布构建的页面 |
| WL Widget Library — Active | `S432k` | "Deploy Widget" | 部署新 Widget |
| WL Hub — Active (checklist) | `BnkYW` | "Announce" / Promo Kit | 首次推广 |

#### Quest 产品

| 页面 | Node ID | 触发按钮 | 场景 |
|------|---------|----------|------|
| Quest Wizard Step 2 | `hBgxY` | "Create Campaign" | 创建 Quest 活动 |
| Quest Hub — Management | `XvXEQ` | "Publish" (draft→active) | 将草稿活动发布 |

#### Boost 产品

| 页面 | Node ID | 触发按钮 | 场景 |
|------|---------|----------|------|
| Boost Wizard Step 4 | `fZpcQ` | "Launch Campaign" | 启动 Boost 活动 |
| Boost Hub — Management | `8gT3V` | "Activate" | 激活暂停的活动 |

### 设计任务

- [x] 新建 Modal D20 `fY99y` — 2 项 checklist（订阅状态 ✓ + 推特授权 ✗）+ 警告条 + Publish 按钮 disabled
- [x] 展示 partial fail 状态（最具信息量的状态，all-pass/all-fail 为此状态的简单变体）
- [x] Canvas label + routing annotation 已添加
- [x] pages.md / layout.md 已更新
- [ ] 在 frontend 文档中标注所有 Publish 按钮 "→ D20" 路由（待 frontend 开发时补充）

---

## Phase: Community Integration Page — COMPLETE

> **Objective**: Community 产品需要独立的 Integration 页面
> **Status**: 全部完成

### 需求

- Community Hub Deep (B12 `TQR51`) 已有 INTEGRATIONS 区域（Twitter/Discord/Telegram + "All Integrations →"）
- Community Integration 页面 → **B61** `ZL5K5`（Row 8 col 4）
- 集成项：Twitter、Discord、Telegram、合约注册、Webhook
- 不需要 WL 专属的：SDK、Iframe、Page Builder、Custom Domain

### 设计任务

- [x] 新建 Community Integration 页面 → **B61** `ZL5K5`
- [x] Community 侧栏不需要 Integration 子项（从 Overview 进入）
- [x] B12 `TQR51` 已有 INTEGRATIONS 区域 + "All Integrations →" → B61（颜色已修正为绿色）
- ~~[ ] 在 Community Hub Management `UPAfV` 的 TOOLKIT 区域增加 "Integrations" 卡片~~ → 原 todo 中 `UPAfV` 为误写（实为 WL Hub），Community 入口已在 B12 INTEGRATIONS 区域

---

## Phase: Missing Target Pages — IDENTIFIED

> **Objective**: Pages referenced by button routes but not yet designed
> **Status**: Catalogued — for frontend dev phase decision
> **Source**: Button routing audit across all 93 page codes

### SKIPPED Pages (confirmed skip — use existing patterns)

| Target | Referenced From | Route | Notes |
|--------|----------------|-------|-------|
| **T01** Quest Campaign Detail | B02, B05, B06 | `/quest/:id` | View campaign detail — follow existing pattern |
| **T02** Quest Campaign Edit | B02, B05 | `/quest/:id/edit` | Edit campaign — follow existing pattern |
| **T03** Quest Wizard Step 3 | B08 | `/quest/create` step 3 | Set Rewards — follow wizard pattern |
| **T04** Quest Wizard Step 4 | (follows T03) | `/quest/create` step 4 | Review & Launch — follow wizard pattern |
| **T05** Boost Campaign Edit | B28 | `/boost/:id/edit` | Edit boost — follow existing pattern |

### FUTURE Pages (not yet needed — mark as v2)

| Target | Referenced From | Route | Notes |
|--------|----------------|-------|-------|
| Case Study Detail | M11 | `/case-studies/:slug` | Individual case study pages |
| C-End Quest Detail | C02 | `/quests/:id` | In-app quest detail (C-end) |
| C-End User Profile | C03 | `/profile/:address` | Public user profile page |
| C-End Sprint Detail | C04 | `/lb-sprint/:id` | Past sprint detail view |
| Quest Campaign Report | B05 | `/quest/:id/report` | Campaign analytics report |

### Design System Gap

| Item | Notes |
|------|-------|
| **D20** Publish Readiness Check modal | Spec in todo.md above — not yet designed |

---

## Phase: Page Code & Routing Annotation Audit — COMPLETE

> **Objective**: Add page codes (M/B/C/D) to all pages, add button routing annotations below each page on canvas
> **Status**: DONE — 2026-03-06

### What was done
1. **Frame renames**: All 100+ page frames renamed with codes (e.g., `QszRH` → "M01 — Brand Homepage")
2. **Canvas labels**: Gray (#94A3B8) code labels above every page — updated ~40 old P-codes to M/B/C, added ~50 new labels (modules, WL wizard variants, modals, C-end)
3. **Routing annotations**: Amber (#F59E0B) text below every page — button→target mappings for all 93 page codes + 19 modals
4. **Position fixes**: Repositioned misplaced C-End labels (y:42120→49620/51520), B51-B53 annotations (y:38060→45460/45560), B54 annotation (y:42110→49560)
5. **Missing pages catalogued**: See "Missing Target Pages" section above

---

## Phase: Dark Theme Restyle — COMPLETE

> **Objective**: Apply the production dark theme (from `design/legacy/theme/`) to all 81 pages + 19 modals.
> **Status**: DONE — all 15 batches executed + second-pass color fix + visual QA

### Execution Summary

| Step | What | Status |
|------|------|--------|
| Color mapping | Built light→dark mapping from production screenshots | Done |
| Pass 1 (Batch 1-15) | `replace_all_matching_properties` on all page groups | Done |
| Pass 2 (fix) | Caught 11 additional colors: `#f7fafc`, `#e6f2ff`, `#fff5f5`, `#fef2f2`, `#f3f0ff`, `#e9e0ff`, `#f1f5f9`(fill), `#bfdbfe`, `#93c5fd`(fill), `#fffbeb`, `#fee2e2` | Done |
| Visual QA | Spot-checked 12+ pages across all types + layout checks | Done, 0 issues |

### Pages Restyled (100 nodes total)

- **Components**: 327GX, kXg2k, 43K4T, ClmXH
- **B-End (50+ pages)**: Dashboard, Quest, Community, WL, Boost, Analytics, Settings, Dev Kit
- **C-End (9 pages)**: Community Home through Activity Feed
- **Marketing (14 pages)**: M02-M14 (M01 was already dark)
- **Modals (19)**: D01-D19
- **Skipped**: M01 (already dark), CXtOH/EDoSn/PWQV6 (deprecated)

---

## Color Mapping Reference (Light → Dark)

### Fill Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#f8fafc` | `#0A0F1A` | Page background |
| `#ffffff` | `#111B27` | Card/panel background |
| `#fafafa` | `#0D1520` | Subtle bg variant |
| `#f7fafc` | `#0A0F1A` | Near-white bg variant |
| `#f1f5f9` | `#1E293B` | Light gray bg / status |
| `#e2e8f0` | `#1E293B` | Dividers, separators |
| `#d1d5db` | `#2D3748` | Gray fill elements |
| `#eff6ff` | `#0F1A2E` | Blue tint bg |
| `#dbeafe` | `#1A2744` | Blue bg |
| `#ebf4ff` | `#0F1A2E` | Blue tint bg |
| `#e6f2ff` | `#0F1A2E` | Blue tint bg |
| `#bfdbfe` | `#1A2744` | Light blue |
| `#93c5fd` | `#1E3A5E` | Medium blue (as fill) |
| `#f0fdf4` | `#0A1F1A` | Green tint bg |
| `#ecfdf5` | `#0A1F1A` | Green tint bg |
| `#f0fff4` | `#0A1F1A` | Green tint bg |
| `#dcfce7` | `#0A2E1A` | Green badge/status bg |
| `#faf5ff` | `#1A1033` | Purple tint bg |
| `#f3e8ff` | `#1A1033` | Purple tint bg |
| `#ede9fe` | `#1A1033` | Purple tint bg |
| `#f5f3ff` | `#1A1033` | Purple tint bg |
| `#f3f0ff` | `#1A1033` | Purple tint bg |
| `#e9e0ff` | `#1A1033` | Purple tint bg |
| `#fff5eb` | `#1F1508` | Orange tint bg |
| `#ffedd5` | `#1F1508` | Orange tint bg |
| `#fff7ed` | `#1F1508` | Orange tint bg |
| `#fef3c7` | `#1F1A08` | Amber/yellow tint bg |
| `#fffbeb` | `#1F1A08` | Amber light bg |
| `#fff5f5` | `#1F0D0D` | Red tint bg |
| `#fef2f2` | `#1F0D0D` | Red tint bg |
| `#fee2e2` | `#2D1515` | Red status bg |

### Text Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#1e293b` | `#F1F5F9` | Primary text |
| `#2d3748` | `#E2E8F0` | Strong text |
| `#1a202c` | `#FFFFFF` | Strongest text |
| `#4a5568` | `#CBD5E1` | Medium text |
| `#475569` | `#CBD5E1` | Medium text |
| `#718096` | `#94A3B8` | Muted text |
| `#a0aec0` | `#94A3B8` | Muted text |
| `#000000` | `#FFFFFF` | Black → white |
| `#cbd5e1` | `#475569` | Placeholder/disabled |
| `#e2e8f0` | `#334155` | Decorative text |

### Stroke Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#e2e8f0` | `#1E293B` | Default border |
| `#93c5fd` | `#1E3A5E` | Blue accent border |
| `#c6f6d5` | `#1A3D2A` | Green accent border |

### Unchanged Colors
- Product brands: Quest `#5D7EF1`, Community `#48BB78`, WL `#9B7EE0`, Boost `#ED8936`
- Status: Active `#16A34A`, Draft `#D97706`, semantic blues/purples/oranges
- Social: Discord `#5865F2`, Telegram `#229ED9`, YouTube `#FF0000`, LinkedIn `#0A66C2`
