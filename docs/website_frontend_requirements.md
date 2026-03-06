# TaskOn 官网前端开发需求文档 v4.2

> 本文档定义了 TaskOn 官网（marketing pages）、B-end 产品界面（after-login dashboard & product hubs）、及 C-end 用户界面（White Label 社区前台）的完整页面编码、按钮→页面路由映射、动态数据规格、及静态内容边界。
> 供前端/后端工程师对照实施。
>
> v4.2 更新（OB Restructuring）：
> - 新增 B55=Community Wizard Step 4 Preview & Publish (`7mVsZ`), B56=WL Wizard Step 4 Preview (`WsH2y`)
> - B35 重定义：Review & Launch → Quick Setup（模板内容自动生成+内联编辑）
> - B19 重命名：WL Embed Options → WL Deployment Settings
> - B37 更新：Choose integration mode → Choose deployment path (Embed/Domain/SDK)
> - B38 更新：Complete Setup → Next: Preview（不再是最终步，路由至 B56）
> - B09/B10/B11/B14/B15/B17 描述更新（OB 重构后的新页面内容）
> - Community Wizard: 3→4 steps (B13→B34→B35→B55), WL Wizard: 3→4 steps (B37→B17→B38→B56)
> - 共 88 页面编码（M01-M14 + B01-B56 + B31a-B31i + C01-C09）。B48 (Dev Kit Page) 待设计。
>
> v4.1 更新：
> - 修正 Community Wizard 页面编码：B13=Step1(`Gzpeu`), B34=Step2(`8NeyG`), B35=Step3(`qknQZ`)。旧节点 `cNwNP`/`Sq7A2` 已删除。
> - 新增 B10/B15 Getting Started Checklist 按钮路由（Promo Kit、Dev Kit、Auto-detection、Social Share）
> - 新增 §六 功能组件规格：Promo Kit Generator、Dev Kit Page (B48)、Auto-detection Systems
> - 新增 WL Sidebar 子菜单路由（Overview/Widgets/Pages）
> - 新增 B31a-B31h（8 个 Community 模块管理页面）

---

## 一、页面编码总表

### Marketing Pages（登录前，M01-M14）

| Code | Page | Design ID | URL | Theme |
|------|------|-----------|-----|-------|
| **M01** | Brand Homepage | `QszRH` | `/` | Dark (#0A0F1A) |
| **M02** | Projects Landing | `Lz2vL` | `/for-projects` | Light |
| **M03** | Quest Product | `gXQur` | `/products/quest` | Light, blue |
| **M04** | Community Product | `GyyL4` | `/products/community` | Light, green |
| **M05** | White Label Product | `cbBdG` | `/products/white-label` | Light, purple |
| **M06** | Boost Product | `Lym65` | `/products/boost` | Light, orange |
| **M07** | Unified Pricing | `HO2Ny` | `/pricing` | Light, 3 product tabs (Quest/Community/WL) |
| ~~M08~~ | ~~White Label Pricing~~ | ~~`EDoSn`~~ | — | Deprecated: merged into M07 Unified Pricing |
| **M09** | Contact / Book Demo | `4q01T` | `/contact` | Light |
| **M10** | About Us | `03CDo` | `/about` | Light |
| **M11** | Case Studies | `XIChW` | `/case-studies` | Light |
| **M12** | Solutions: CPA/CPS | `wsqIT` | `/solutions/cpa` | Light, orange |
| **M13** | Solutions: Custom | `A7FaV` | `/solutions/custom` | Light |
| **M14** | Solutions: Joint | `huslr` | `/solutions/joint` | Light |

### B-End Pages（登录后，B01-B47）

> 所有 B-End 页面共享 Sidebar + Top Bar。域名: `dashboard.taskon.xyz`。Theme: Light (#F8FAFC)。

| Code | Page | Design ID | URL Path | Notes |
|------|------|-----------|----------|-------|
| **B01** | Dashboard (New User) | `4SMOO` | `/` | Welcome + 4 goal cards |
| **B02** | Dashboard (Active) | `IDezm` | `/` | Stats + campaigns + quick actions |
| **B03** | Dashboard (Power) | `W93vp` | `/` | Growth chart + product breakdown |
| **B04** | Quest Hub (Empty) | `VkSu7` | `/quest` | Templates + CTA |
| **B05** | Quest Hub (Active) | `4Mb1C` | `/quest` | 2 campaigns + stats |
| **B06** | Quest Hub (Management) | `XvXEQ` | `/quest` | Table + filters + pagination |
| **B07** | Quest Wizard Step 1 | `dn8Eu` | `/quest/create` | Template selection |
| **B08** | Quest Wizard Step 2 | `hBgxY` | `/quest/create` (step 2) | Task configuration |
| **B09** | Community Hub (Empty) | `zzZ8D` | `/community` | Selectable strategy cards + single CTA |
| **B10** | Community Hub (Guided) | `S1EIA` | `/community` | 3-section checklist (Wizard/Enrich/Go Live) + modules |
| **B11** | Community Hub (Active) | `vFRHi` | `/community` | Stats + trend arrows + module performance cards |
| **B12** | Community Hub (Deep) | `TQR51` | `/community` | 6 modules + engagement chart |
| **B13** | Community Wizard Step 1 | `Gzpeu` | `/community/create` (step 1) | Customize: name, description, brand color |
| **B14** | White Label Hub (Empty) | `Ir6Tq` | `/white-label` | Selectable deployment path cards + single CTA |
| **B15** | White Label Hub (Active) | `BnkYW` | `/white-label` | 3-section checklist (Wizard/Configure/Deploy) + toolkit |
| **B16** | White Label Hub (Mgmt) | `UPAfV` | `/white-label` | All tools + deployments + analytics |
| **B17** | WL Wizard Step 2 | `CXzmy` | `/white-label/setup` (step 2) | Path-adaptive: widgets/DNS/SDK config |
| **B18** | WL Domain Setup | `5bmH9` | `/white-label/domain` | 3-step domain config + status |
| **B19** | WL Deployment Settings | `RgCVQ` | `/white-label/deploy` | Current deployment path + add additional methods |
| **B20** | WL Widget Library (Empty) | `2sSsA` | `/white-label/widgets` | 7 community modules (6+1 unconfigured) |
| **B21** | WL Widget Config | `n4pJK` | `/white-label/widgets/:id/config` | 2-col: settings + preview + embed code |
| **B22** | WL Widget Library (Active) | `S432k` | `/white-label/widgets` | My Widgets (2) + community modules |
| **B23** | WL Page Builder (Empty) | `DRYwN` | `/white-label/pages` | How it works + templates |
| **B24** | WL Page Builder (Editor) | `sGDcq` | `/white-label/pages/new` | Canvas + settings panel |
| **B25** | WL Page Builder (Active) | `J08v5` | `/white-label/pages` | My Pages (1) + quick start |
| **B26** | WL Integration Center | `Abs1E` | `/white-label/integrations` | 4 categories × 3 integrations |
| **B27** | Boost Hub (Empty) | `stYvi` | `/boost` | 3 goal cards + CTA |
| **B28** | Boost Hub (Active) | `5C3WP` | `/boost` | 2 campaigns (active + review) |
| **B29** | Boost Hub (Management) | `8gT3V` | `/boost` | Table + filters + pagination |
| **B30a** | Boost Wizard Step 1 | `SDfui` | `/boost/create` | Goal selection |
| **B30b** | Boost Wizard Step 2 | `l9tmF` | `/boost/create` (step 2) | Channel selection |
| **B30c** | Boost Wizard Step 3 | `KMtqR` | `/boost/create` (step 3) | Budget + CPA calculator |
| **B30d** | Boost Wizard Step 4 | `fZpcQ` | `/boost/create` (step 4) | Review & submit |
| **B31** | Sectors & Tasks | `Wug7d` | `/community/sectors` | Drag-handle task management |
| **B31a** | Points & Level Mgmt | `zCfKQ` | `/community/modules/points` | Level config + stats table |
| **B31b** | TaskChain Mgmt | `lpdtp` | `/community/modules/taskchain` | Chain config + stats table |
| **B31c** | DayChain Mgmt | `fLLVb` | `/community/modules/daychain` | Streak config + stats table |
| **B31d** | Leaderboard Mgmt | `Emmab` | `/community/modules/leaderboard` | Multi-point-type (EXP/GEM) recurring leaderboards, weekly/monthly/all-time cycles, no extra incentives |
| **B31e** | LB Sprint Mgmt | `FO9JR` | `/community/modules/lb-sprint` | Time-bounded leaderboard competitions, point-type-based, non-points incentives (NFT/Token/WL), start/end dates |
| **B31f** | Milestone Mgmt | `WFdZQ` | `/community/modules/milestone` | Milestone config + stats table |
| **B31g** | Benefits Shop Mgmt | `7yPWx` | `/community/modules/shop` | Item config + stats table |
| **B31h** | Lucky Wheel Mgmt | `sme5a` | `/community/modules/wheel` | Wheel config + stats table |
| **B32** | Content Management | `lhR14` | `/community/content` | Announcements + featured + modules |
| **B33** | Preview Mode | `2UiNC` | `/community/preview` | Embedded C-end preview |
| **B34** | Community Wizard Step 2 | `8NeyG` | `/community/create` (step 2) | Configure modules (toggles + inline config) |
| **B35** | Community Wizard Step 3 | `qknQZ` | `/community/create` (step 3) | Quick Setup: auto-generated template content + inline edit |
| **B36** | Community Module Detail | `usBsM` | `/community/modules/:type` | Module config (Points example) |
| **B37** | WL Wizard Step 1 | `NNwid` | `/white-label/setup` (step 1) | Choose deployment path (Embed/Domain/SDK) |
| **B38** | WL Wizard Step 3 | `5nCtO` | `/white-label/setup` (step 3) | Brand: logo, colors, typography + live preview |
| **B39** | Boost Campaign Detail | `Sq4jV` | `/boost/:id` | Campaign metrics + management |
| **B40** | WL Brand Settings | `Cx3LH` | `/white-label/brand` | Logo, colors, typography |
| **B41** | WL SDK & API | `lQxT5` | `/white-label/sdk` | API keys, SDK docs, endpoints |
| **B42** | WL Iframe Embed | `ByGS0` | `/white-label/embed/iframe` | Iframe embed flow |
| **B43** | WL Page Analytics | `69HPh` | `/white-label/pages/:id/analytics` | Page performance metrics |
| **B44** | WL Integration Config | `gS64G` | `/white-label/integrations/:type` | Integration setup (template) |
| **B45** | Analytics Dashboard | `fLxTr` | `/analytics` | Cross-product analytics |
| **B46** | Settings | `ESrVt` | `/settings` | Account settings |
| **B47** | Settings / Profile | `Nh7xq` | `/settings/profile` | Profile management |
| **B31i** | Badges | `BJLsz` | `/community/modules/badges` | Badge management: create/edit/archive badges, stats table |
| **B49** | Access Rules | `g1CNC` | `/community/settings/access-rules` | Token-gate & role-based access rules management |
| **B50** | Homepage Editor | `5Wm6B` | `/community/settings/homepage` | Drag-order section management for C-End homepage |
| **B51** | WL Contract Registry | `OKEqS` | `/white-label/contracts` | Smart contract registry for on-chain reward distribution |
| **B52** | WL Activity Rule Builder | `4aAo7` | `/white-label/rules` | Visual rule builder for automated activity-based triggers |
| **B53** | WL Privilege Manager | `5xwYN` | `/white-label/privileges` | User privilege tiers, token-gated access levels |
| **B54** | Community Insights | `olPfE` | `/community/insights` | Cross-module analytics, economy health, user segments, retention |
| **B55** | Community Wizard Step 4 | `7mVsZ` | `/community/create` (step 4) | Preview & Publish: C-end preview + readiness checklist |
| **B56** | WL Wizard Step 4 | `WsH2y` | `/white-label/setup` (step 4) | Preview: deployment preview + readiness checklist |

### C-End Pages（White Label 社区前台，C01-C09）

> C-End 页面为 White Label Community 的用户端。域名由项目自定义（如 `community.bitcoin.com`）。Theme: Dark header (#0F172A) + Light content (#F8FAFC)。Accent: Amber (#F59E0B)。
>
> **Leaderboard vs LB Sprint 概念区分**:
> - **Leaderboard (C03)**: 基于自定义积分类型（EXP/GEM/等）的周期性排行榜展示（周/月/全部），**无额外激励**，纯展示排名。B-End 管理在 B31d。
> - **LB Sprint (C04)**: 基于自定义积分的**限时排行榜竞赛**，有明确的开始/结束日期，附带**非积分类激励**（NFT/Token/WL Spot 等）。B-End 管理在 B31e。

| Code | Page | Design ID | URL Path | Notes |
|------|------|-----------|----------|-------|
| **C01** | Community Home | `vJVhd` | `/` | Action engine: personal card + urgent actions + DayChain + tasks + invite banner + community pulse |
| **C02** | Quest Tab | `dUXTl` | `/quests` | Quest card grid 2×2 + filters |
| **C03** | Leaderboard | `KmdSd` | `/leaderboard` | Multi-point-type (EXP/GEM) recurring leaderboard, period filters (Weekly/Monthly/All Time), no extra incentives |
| **C04** | LB Sprint Tab | `y5fUZ` | `/lb-sprint` | Time-bounded leaderboard competition with start/end dates, point-type-based (EXP/GEM), non-points incentives (NFT/tokens/WL), reward tiers, past sprints |
| **C05** | Milestone Tab | `53iKE` | `/milestones` | Level progress + milestone cards (claimed/claimable/locked) |
| **C06** | Shop Tab | `coM7o` | `/shop` | Rewards shop + 6 items (redeem/not enough/sold out) |
| **C07** | User Center | `nvQuK` | `/profile` | Profile card + stats + achievement badges + activity history + referral stats |
| **C08** | Invite Center | `wSSD5` | `/invite` | Referral link + share buttons + reward tiers (3) + top inviters leaderboard |
| **C09** | Activity Feed | `5Nhkz` | `/activity` | Live community feed + filter chips + trending cards |

---

## 二、全局可复用组件

### 2.1 Marketing Header (B端版)

用于 M02-M14（Projects Landing 及产品/定价/辅助页）。

| 元素 | 链接目标 |
|------|----------|
| Logo | → M01 `/` |
| Products ▼ | 下拉: Quest → M03, Community → M04, White Label → M05 |
| Solutions ▼ | 下拉: CPA → M12, Custom → M13, Joint → M14 |
| Pricing | → M07 `/pricing` |
| Case Studies | → M11 `/case-studies` |
| "Book Demo" | → M09 `/contact` |
| "Launch Campaign" | → `app.taskon.xyz/create` (ext) |

### 2.2 Brand Homepage Header (最简版)

仅用于 M01。

| 元素 | 链接目标 |
|------|----------|
| Logo "TaskOn" | → M01 `/` |
| "For Projects" | → M02 `/for-projects` |
| "Launch App" | → `app.taskon.xyz` (ext) |

### 2.3 Footer (全站统一)

| 元素 | 链接目标 |
|------|----------|
| **For Users** | |
| Explore Tasks | → `app.taskon.xyz/explore` (ext) |
| Leaderboard | → `app.taskon.xyz/leaderboard` (ext) |
| Rewards | → `app.taskon.xyz/rewards` (ext) |
| **For Projects** | |
| Quest Campaigns | → M03 `/products/quest` |
| Community & Loyalty | → M04 `/products/community` |
| White-Label SDK | → M05 `/products/white-label` |
| CPA/CPS Delivery | → M12 `/solutions/cpa` |
| Custom & Managed | → M13 `/solutions/custom` |
| Joint Programs | → M14 `/solutions/joint` |
| Platform Pricing | → M07 `/pricing` |
| **Company** | |
| About Us | → M10 `/about` |
| Case Studies | → M11 `/case-studies` |
| Blog | → `/blog` (ext) |
| Contact | → M09 `/contact` |
| Social Icons | → Twitter/GitHub/Discord/Telegram (ext) |

### 2.4 B-End Sidebar (全局)

| 元素 | 链接目标 |
|------|----------|
| Logo "TaskOn" | → B01 `/` |
| Home | → B01 `/` |
| Quest | → B04/B05/B06 `/quest` |
| Community | → B09/B10/B11/B12 `/community` |
| White Label ▾ | → B14/B15/B16 `/white-label` (expandable sub-menu) |
| ├─ Overview | → B14/B15/B16 `/white-label` |
| ├─ Widgets | → B20/B22 `/white-label/widgets` |
| └─ Pages | → B23/B25 `/white-label/pages` |
| Boost | → B27/B28/B29 `/boost` |
| Analytics | → B45 `/analytics` |
| Settings | → B46 `/settings` |

### 2.5 B-End Top Bar (全局)

| 元素 | 链接目标 |
|------|----------|
| Breadcrumb | Dynamic 路由面包屑 |
| Help | → 帮助中心 (ext) |
| Notification bell | → (modal) 通知面板 |
| User avatar ▼ | Profile → B47 `/settings/profile`, Settings → B46 `/settings`, Logout → 登出 |

### 2.6 C-End Header (全局)

用于 C01-C09。Dark theme (#0F172A)。

| 元素 | 链接目标 |
|------|----------|
| Project Logo + Name | → C01 `/` |
| "Connect Wallet" | (action) 钱包连接弹窗 |
| Nav: Home | → C01 `/` |
| Nav: Quests | → C02 `/quests` |
| Nav: Leaderboard | → C03 `/leaderboard` |
| Nav: LB Sprint | → C04 `/lb-sprint` |
| Nav: Milestone | → C05 `/milestones` |
| Nav: Shop | → C06 `/shop` |
| Nav: Profile | → C07 `/profile` |
| Nav: Invite | → C08 `/invite` |
| Nav: Activity | → C09 `/activity` |
| Notification bell (header) | → (modal) 通知面板 |
| User avatar (header) | → C07 `/profile` |

### 2.7 C-End Footer (全局)

| 元素 | 链接目标 |
|------|----------|
| "Powered by TaskOn" | → `taskon.xyz` (ext) |

---

## 三、完整按钮→页面路由映射表

> 标记说明:
> - `→ B05` = 跳转到编码为 B05 的页面
> - `(ext)` = 外部链接
> - `(modal)` = 弹窗/面板
> - `(action)` = 前端交互（无页面跳转）
> - `(API)` = API 调用
> - **SKIPPED** = 目标页面已确认跳过设计（使用现有系统模式）

---

### M01: Brand Homepage (`/`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "For Projects" | Header | → M02 |
| "Launch App" | Header | → `app.taskon.xyz` (ext) |
| "Launch App" 主CTA | Hero | → `app.taskon.xyz` (ext) |
| "For Projects →" | Hero | → M02 |
| Trust Bar Logos | Hero | 无链接 |
| Task Card "Claim" / "Start" | 浮动卡 | → `app.taskon.xyz/quest/:id` (ext) |
| "View Full Leaderboard →" | Happening Now | → `app.taskon.xyz/leaderboard` (ext) |
| "See How It Works" | For Projects | → M02 |
| "Launch App" | Final CTA | → `app.taskon.xyz` (ext) |
| "For Projects →" | Final CTA | → M02 |
| Footer links | Footer | 见 §2.3 |

### M02: Projects Landing (`/for-projects`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Start Free" | Hero | → `app.taskon.xyz/create` (ext) |
| "Book Demo" | Hero | → M09 `/contact` |
| Boost card "Learn More" | Growth Engines | → M06 |
| Quest card "Learn More" | Growth Engines | → M03 |
| Community card "Learn More" | Growth Engines | → M04 |
| White Label card "Learn More" | Growth Engines | → M05 |
| "Start Free" | Final CTA | → `app.taskon.xyz/create` (ext) |
| "Book Demo" | Final CTA | → M09 `/contact` |
| Footer links | Footer | 见 §2.3 |

### M03: Quest Product (`/products/quest`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Start Free" | Hero | → `app.taskon.xyz/create` (ext) |
| "See Pricing" | Hero | → M07 |
| "Start Acquiring Real Users Today" | Final CTA | → `app.taskon.xyz/create` (ext) |
| "Try Community →" | Final CTA (cross-sell) | → M04 |
| Footer links | Footer | 见 §2.3 |

### M04: Community Product (`/products/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Start Free" | Hero | → `app.taskon.xyz/create` (ext) |
| "See Pricing" | Hero | → M07 |
| "Start Building Retention Today" | Final CTA | → `app.taskon.xyz/create` (ext) |
| "Try White Label →" | Final CTA (cross-sell) | → M05 |
| Footer links | Footer | 见 §2.3 |

### M05: White Label Product (`/products/white-label`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Talk to Sales" | Hero | → M09 `/contact` |
| "See Pricing" | Hero | → M08 |
| Domain Mode card | Integration Modes | → M09 `/contact` |
| Widget + Page Builder ★ card | Integration Modes | → M09 `/contact` |
| Full SDK card | Integration Modes | → M09 `/contact` |
| "Own Your Growth Stack" | Final CTA | → M09 `/contact` |
| Footer links | Footer | 见 §2.3 |

### M06: Boost Product (`/products/boost`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Get Started" | Hero | → M09 `/contact` |
| "Ready to Get Guaranteed Results?" | Final CTA | → M09 `/contact` |
| "Try Quest →" | Final CTA (cross-sell) | → M03 |
| Footer links | Footer | 见 §2.3 |

### M07: Unified Pricing (`/pricing`) — 替代旧 M07+M08

> **产品 Tab 切换**: Quest / Community / White Label（3 个 tab，切换后更新定价卡片、功能列表、账单周期选项）
>
> **账单周期规则（按产品不同）:**
> - **Quest & Community**: Quarterly（季付）/ Semi-Annual（半年付）/ Annual（年付）。**不支持月付。**
> - **White Label**: Monthly（月付）/ Quarterly（季付）/ Semi-Annual（半年付）。**不支持年付。**
>
> **定价基准（月等效价格）:**
> - Quest: $300/mo | Community: $600/mo | White Label: $1,500/mo
>
> **设计状态**: 当前画面展示 Quest tab + Quarterly 选中态。前端需根据 tab 切换动态更新 billing toggle 选项。

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| Quest / Community / White Label tabs | Product tabs | (action) 前端切换产品 + 更新账单周期选项 |
| Billing cycle toggle (3 options) | Billing | (action) 前端切换账单周期 + 更新价格 |
| "Start 7-Day Free Trial" | Pricing card | → `app.taskon.xyz/create?product={quest\|community\|wl}&plan=trial` (ext) |
| "Subscribe & Pay Now" | Pricing card | → `app.taskon.xyz/create?product={quest\|community\|wl}&cycle={q\|sa\|a\|m}` (ext) |
| Growth Stacks cards | Cross-sell section | (info only, no navigation) |
| ROI Calculator section | WL tab only | (shown conditionally when WL tab selected) |
| FAQ expand/collapse | FAQ | (action) |
| "Start Your Free Trial Today" | Final CTA | → `app.taskon.xyz/create` (ext) |
| "Book a Demo" | Final CTA | → M09 `/contact` |
| Footer links | Footer | 见 §2.3 |

### M09: Contact / Book Demo (`/contact`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Submit" | Contact form | (API) `POST /api/contact` |
| Footer links | Footer | 见 §2.3 |

### M10: About Us (`/about`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Get Started" | CTA | → M02 |
| Footer links | Footer | 见 §2.3 |

### M11: Case Studies (`/case-studies`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| Case study cards | Content | → `/case-studies/:slug` (detail, future) |
| Footer links | Footer | 见 §2.3 |

### M12: Solutions CPA/CPS (`/solutions/cpa`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Get Started" | CTA | → M09 `/contact` |
| "Learn More about Boost" | Cross-sell | → M06 |
| Footer links | Footer | 见 §2.3 |

### M13: Solutions Custom (`/solutions/custom`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Talk to Sales" | CTA | → M09 `/contact` |
| Footer links | Footer | 见 §2.3 |

### M14: Solutions Joint (`/solutions/joint`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| "Apply Now" | CTA | → M09 `/contact` |
| Footer links | Footer | 见 §2.3 |

---

### B01: Dashboard — New User (`/`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "→ Quest" card | Goal cards | → B04 |
| "→ Community" card | Goal cards | → B09 |
| "→ White Label" card | Goal cards | → B14 |
| "→ Boost" card | Goal cards | → B27 |
| "Create Your First Quest →" | Quick start | → B07 |
| Watch Tutorial | Resources | → (ext) 视频 |
| Quick Start Guide | Resources | → (ext) 帮助中心 |
| Talk to Sales | Resources | → M09 `/contact` (ext) |

### B02: Dashboard — Active User (`/`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "New Campaign" button | Header | → B07 (default to Quest create) |
| "View All →" | Campaign section | → B05 or B06 (Quest hub) |
| Campaign card "View" | Campaign card | → `/quest/:id` **SKIPPED** (campaign detail) |
| Campaign card "Edit" | Campaign card | → `/quest/:id/edit` **SKIPPED** |
| Campaign card "Pause" | Campaign card | (API) |
| "Create Quest" | Quick Actions | → B07 |
| "View Analytics" | Quick Actions | → B45 `/analytics` |
| "Duplicate Campaign" | Quick Actions | (API) |
| "Export Data" | Quick Actions | (action) 下载 |
| "Learn More" (WL upsell) | Upsell banner | → M05 (ext, new tab) |

### B03: Dashboard — Power User (`/`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "New Campaign" button | Header | → B07 |
| "View All →" | Recent Activity | → B45 `/analytics` |
| Product Breakdown items | Panel | Quest → B05, Community → B11, Boost → B28, WL → B15 |

### B04: Quest Hub — Empty (`/quest`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Template "Use Template" (×3) | Template cards | → B07 `/quest/create?template=social\|trading\|liquidity` |
| "Create Your First Quest" | 主CTA | → B07 |
| "or start from a blank canvas" | Sub-CTA | → B07 `/quest/create?template=blank` |
| Video Tutorial | Resources | → (ext) 视频 |
| Best Practices Guide | Resources | → (ext) 帮助中心 |

### B05: Quest Hub — Active (`/quest`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Quest" button | Header | → B07 |
| Campaign card "View" | Camp 1 (Active) | → `/quest/:id` **SKIPPED** |
| Campaign card "Edit" | Camp 1 (Active) | → `/quest/:id/edit` **SKIPPED** |
| Campaign card "Pause" | Camp 1 (Active) | (API) |
| Campaign card "View Report" | Camp 2 (Completed) | → `/quest/:id/report` **SKIPPED** |
| Campaign card "Duplicate" | Camp 2 (Completed) | (API) → B07 (pre-filled) |

### B06: Quest Hub — Management (`/quest`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Quest" button | Header | → B07 |
| Filter tabs (All/Active/Draft/Completed) | Filter bar | (action) 前端筛选 |
| Search input | Filter bar | (action) 前端搜索 |
| Table row click | Table | → `/quest/:id` **SKIPPED** |
| Pagination | Bottom | (action) |

### B07: Quest Wizard Step 1 (`/quest/create`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Save Draft" | Top bar | (API) `POST /api/quest/drafts` |
| Template cards (6) | Content | (action) 选择模板 |
| Category filters | Content | (action) 筛选 |
| "Cancel" | Bottom | → B04/B05/B06 (返回 Quest Hub) |
| "Next: Configure Tasks" | Bottom | → B08 |

### B08: Quest Wizard Step 2 (`/quest/create`, step 2)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| "Add Task" | Task list | (action) 添加任务 |
| Task items click | Task list | (action) 显示编辑器 |
| "Back" | Bottom | → B07 |
| "Next: Set Rewards" | Bottom | → Quest Step 3 **SKIPPED** |

### B09: Community Hub — Empty (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Template CTAs (×3) | Template cards | → B13 `/community/create?template=points\|leaderboard\|benefits` |
| "Create Your First Community" | 主CTA | → B13 (Step 1) |
| "or start from a blank canvas" | Sub-CTA | → B13 `/community/create?template=blank` |
| Video Tutorial | Resources | → (ext) 视频 |
| Retention Playbook | Resources | → (ext) 帮助中心 |
| "Learn More" | Resources | → M04 (ext, new tab) |

### B10: Community Hub — Guided Workspace (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Checklist Step 1: "Create your community" | Checklist (完成) | (action) 已完成标记 |
| Checklist Step 2: "Configure modules" | Checklist (完成) | → B36 `/community/modules/:type` |
| Checklist Step 3: "Preview & publish" | Checklist (完成) | → B33 (Preview Mode) |
| Checklist Step 4: "Share with your community" | Checklist (展开态) | 见下方详细路由 |
| ├─ Share Link "Copy" | Link bar | (action) 复制链接到剪贴板 |
| ├─ Twitter button | Social row | (ext) 预填 share 链接到 Twitter |
| ├─ Discord button | Social row | (ext) 预填 share 链接到 Discord |
| ├─ Telegram button | Social row | (ext) 预填 share 链接到 Telegram |
| └─ "Generate Promo Kit" | Promo Kit row | (action/API) → Promo Kit Generator modal |
| Checklist Step 5: "First 10 participants" | Checklist (待完成) | (auto) WebSocket 自动检测 |
| Module cards (Points, Leaderboard, Tasks) | Active Modules | → B36 `/community/modules/:type` |
| "Add" module (Milestones, Benefits Shop) | Add More Modules | (action) 启用模块 |
| Template row link | Resources | → B13 |
| Video Tutorial | Resources | → (ext) 视频 |
| Retention Playbook | Resources | → (ext) 帮助中心 |

### B11: Community Hub — Active (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Quick stats cards (4) | Stats row | 无链接 (展示) |
| Checklist (4/5 complete) | Checklist card | (action) |
| Module cards with metrics | Active Modules | → B36 `/community/modules/:type` |
| "Add" module | Add More Modules | (action) |
| Template row link | Resources | → B13 |
| Resources | Resources | → (ext) |

### B12: Community Hub — Deep (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Module cards (6, 2 rows) | Modules | → B36 `/community/modules/:type` |
| Engagement Overview chart | Analytics | 无链接 (展示) |
| Retention metrics | Analytics | → B45 `/analytics` |

### B13: Community Wizard Step 1 (`/community/create`, step 1)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/community/drafts` |
| Community Name input | Form | (action) |
| Description textarea | Form | (action) |
| Brand Color picker | Form | (action) |
| "Cancel" | Bottom | → B09/B10/B11/B12 (返回 Community Hub) |
| "Next: Configure Modules" | Bottom | → B34 (Community Step 2) |

### B14: White Label Hub — Empty/Toolbox (`/white-label`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Launch Portal" goal card | Goal cards | → B18 (Domain Setup) |
| "Embed Widgets" goal card | Goal cards | → B19 (Embed Options) |
| "Build Pages" goal card | Goal cards | → B23 (Page Builder) |
| Custom Domain toolkit card | Toolkit | → B18 |
| Widget Library toolkit card | Toolkit | → B20 |
| Page Builder toolkit card | Toolkit | → B23 |
| Brand Settings toolkit card | Toolkit | → B40 `/white-label/brand` |
| SDK & API toolkit card | Toolkit | → B41 `/white-label/sdk` |
| Integration Center toolkit card | Toolkit | → B26 |
| SDK Documentation | Resources | → (ext) 文档 |
| Setup Walkthrough | Resources | → (ext) 视频 |
| "Learn More" | Resources | → M05 (ext, new tab) |

### B15: White Label Hub — Active (`/white-label`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 (WL sidebar 含 Overview/Widgets/Pages 子菜单) |
| Checklist Step 1: "Create your project" | Checklist (完成) | (action) 已完成标记 |
| Checklist Step 2: "Choose deployment path" | Checklist (完成) | → B19 (Embed Options) |
| Checklist Step 3: "Send Dev Kit to developer" | Checklist (展开态) | 见下方详细路由 |
| ├─ Dev Kit Link "Copy" | Link bar | (action) 复制 `taskon.xyz/devkit/{id}` |
| ├─ "Email to Developer" | Action row | (action/API) 发送 Dev Kit 链接邮件 |
| └─ Integration status | Status indicator | (auto) WebSocket 自动检测首次 API ping |
| Checklist Step 4: "Integration verified" | Checklist (待完成) | (auto) 自动检测 API 调用 |
| Checklist Step 5: "Announce to your users" | Checklist (待完成) | 类似 B10 Step 4 (Share + Promo Kit) |
| Checklist Step 6: "First user interaction" | Checklist (待完成) | (auto) WebSocket 自动检测 |
| Custom Domain card | Toolkit Row 1 | → B18 |
| Widget Library card | Toolkit Row 1 | → B22 |
| Page Builder card | Toolkit Row 1 | → B25 |
| Brand Settings card | Toolkit Row 2 | → B40 `/white-label/brand` |
| SDK & API card | Toolkit Row 2 | → B41 `/white-label/sdk` |
| Integration Center card | Toolkit Row 2 | → B26 |
| "View Docs →" | Resources | → (ext) 文档 |
| "Watch Video →" | Resources | → (ext) 视频 |
| "Visit Product Page →" | Resources | → M05 (ext, new tab) |

### B16: White Label Hub — Management (`/white-label`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Custom Domain card | Toolkit Row 1 | → B18 |
| Widget Library card | Toolkit Row 1 | → B22 |
| Page Builder card | Toolkit Row 1 | → B25 |
| Brand Settings card | Toolkit Row 2 | → B40 `/white-label/brand` |
| SDK & API card | Toolkit Row 2 | → B41 `/white-label/sdk` |
| Integration Center card | Toolkit Row 2 | → B26 |
| Domain deployment card | Active Deployments | → B18 |
| Widgets deployment card | Active Deployments | → B22 |
| Pages deployment card | Active Deployments | → B25 |
| Usage Analytics chart | Usage Analytics | → B45 `/analytics` |

### B17: WL Wizard Step 2 (`/white-label/setup`, step 2)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/white-label/drafts` |
| Domain input | Left panel | → (API) `POST /api/white-label/verify-domain` |
| SSL toggle | Left panel | (action) |
| "Back" | Bottom | → B37 (WL Step 1) |
| "Next: Generate SDK" | Bottom | → B38 (WL Step 3) |

### B18: WL Domain Setup (`/white-label/domain`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| Step 1: Enter Domain | Content | (action) 输入域名 |
| Step 2: Configure DNS | Content | (action) DNS 配置 |
| Step 3: Brand Portal | Content | (action) 品牌设置 |
| Setup Status checklist | Right panel | (action) 展示进度 |
| "Need Help?" card | Bottom | → (ext) 帮助中心 |

### B19: WL Deployment Settings (`/white-label/deploy`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| Current deployment path display | Content | 无链接 (展示) |
| "Add Deployment Method" | Content | (action) 添加额外部署方式 |
| "Iframe Embed" card | Embed modes | → B42 `/white-label/embed/iframe` |
| "Open Widget Library" ★ card | Embed modes | → B20 (Widget Library) |
| "Open Page Builder" card | Embed modes | → B23 (Page Builder) |
| Comparison table | Content | 无链接 (展示) |

### B20: WL Widget Library — Empty (`/white-label/widgets`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Embed Options" | Breadcrumb | → B19 |
| Configured widget "Add Widget →" (×6) | Community Modules grid | → B21 (Widget Config) |
| Milestones "Set Up in Community →" | Not Yet Configured | → B10/B11 (Community Hub) |
| "Create Your First Widget" | 主CTA | → B21 |
| Tip | Bottom | 无链接 |

### B21: WL Widget Config (`/white-label/widgets/:id/config`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Widget Library" | Breadcrumb | → B20/B22 |
| Settings form fields | Left panel | (action) |
| Live Preview | Right panel | (action) 实时预览 |
| "Copy" embed code | Embed Code | (action) 复制到剪贴板 |
| Help note link | Bottom | → (ext) 帮助中心 |

### B22: WL Widget Library — Active (`/white-label/widgets`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Embed Options" | Breadcrumb | → B19 |
| My Widget card "Edit" | My Widgets | → B21 |
| My Widget card "Copy" embed code | My Widgets | (action) 复制 |
| Configured widget "Add Widget →" (×4) | Community Modules | → B21 |
| Milestones "Set Up in Community →" | Not Yet Configured | → B10/B11 |
| Tip | Bottom | 无链接 |

### B23: WL Page Builder — Empty (`/white-label/pages`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Embed Options" | Breadcrumb | → B19 |
| How It Works (3 steps) | Content | 无链接 (展示) |
| Template cards (×3) | Page Templates | → B24 (pre-filled template) |
| Widget callout | Content | → B20 (Widget Library) |
| "Create Your First Page" | 主CTA | → B24 |
| Tip | Bottom | 无链接 |

### B24: WL Page Builder — Editor (`/white-label/pages/new`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Pages" | Breadcrumb | → B23/B25 |
| Canvas: drag widget blocks | Canvas (left) | (action) |
| "+ Add Widget Block" | Canvas (left) | (action) 添加 widget block |
| Page Settings form | Right panel | (action) |
| Theme toggle (Light/Dark) | Right panel | (action) |
| "× Remove" widget | Widgets on Page | (action) |
| "+ Add" available widget | Available Widgets | (action) |
| "Go to Widget Library →" | Right panel | → B20/B22 |
| "Publish Page" button | Right panel | (API) |
| "Save Draft" button | Right panel | (API) |

### B25: WL Page Builder — Active (`/white-label/pages`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Embed Options" | Breadcrumb | → B19 |
| "Create New Page" button | Header | → B24 |
| Page card "Copy" embed code | My Pages | (action) 复制 |
| Page card "Edit Page" | My Pages | → B24 (edit existing) |
| Page card "Analytics" | My Pages | → B43 `/white-label/pages/:id/analytics` |
| Quick Start template cards (×4) | Quick Start | → B24 (pre-filled template) |

### B26: WL Integration Center (`/white-label/integrations`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| Twitter/X "Configure" | Social (Connected) | → B44 `/white-label/integrations/twitter` |
| Discord "Configure" | Social (Connected) | → B44 `/white-label/integrations/discord` |
| Telegram "Connect" | Social (Available) | → B44 `/white-label/integrations/telegram` |
| Multi-Chain "Connect" | Blockchain | → B44 `/white-label/integrations/multichain` |
| Wallet Connect "Connect" | Blockchain | → B44 `/white-label/integrations/wallet` |
| On-Chain Verification "Connect" | Blockchain | → B44 `/white-label/integrations/onchain` |
| Google Analytics "Connect" | Analytics | → B44 `/white-label/integrations/ga` |
| Webhooks "Connect" | Analytics | → B44 `/white-label/integrations/webhooks` |
| Data Export "Connect" | Analytics | → B44 `/white-label/integrations/export` |
| API Keys "Connect" | Developer | → B44 `/white-label/integrations/api-keys` |
| SDK Configuration "Connect" | Developer | → B44 `/white-label/integrations/sdk-config` |
| SSO/OAuth "Connect" | Developer | → B44 `/white-label/integrations/sso` |

### B27: Boost Hub — Empty (`/boost`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Goal Card "More Details" (×3) | Goal cards | → (ext) 帮助中心 |
| "Start a Boost Campaign" | 主CTA | → B30a |
| "or talk to our growth team" | Sub-CTA | → M09 `/contact` (ext) |
| How Boost Works | Resources | → (ext) 帮助中心 |
| CPA Campaign Guide | Resources | → (ext) 帮助中心 |

### B28: Boost Hub — Active (`/boost`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Campaign" button | Header | → B30a |
| Campaign card actions (Active) | Camp 1 | View → B39 `/boost/:id`, Edit → `/boost/:id/edit` **SKIPPED** |
| Campaign card actions (Under Review) | Camp 2 | View → B39 `/boost/:id` |

### B29: Boost Hub — Management (`/boost`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Campaign" button | Header | → B30a |
| Filter tabs (All/Active/Draft/Completed) | Filter bar | (action) 筛选 |
| Search input | Filter bar | (action) |
| Table row click | Table | → B39 `/boost/:id` |
| Pagination | Bottom | (action) |

### B30a: Boost Wizard Step 1 (`/boost/create`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/boost/drafts` |
| Goal selection cards (3) | Content | (action) 选择目标 |
| Targeting options | Content | (action) |
| "Cancel" | Bottom | → B27/B28/B29 |
| "Next: Select Channels" | Bottom | → B30b |

### B30b: Boost Wizard Step 2 (`/boost/create`, step 2)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Channel checkboxes (3) | Left panel | (action) |
| Campaign Summary | Right panel | (action) 实时计算 |
| "Back" | Bottom | → B30a |
| "Next: Set Budget" | Bottom | → B30c |

### B30c: Boost Wizard Step 3 (`/boost/create`, step 3)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| CPA pricing input | Form | (action) |
| Budget cap input | Form | (action) |
| Duration picker | Form | (action) |
| Estimated reach calculator | Right panel | (action) 实时计算 |
| "Back" | Bottom | → B30b |
| "Next: Review & Submit" | Bottom | → B30d |

### B30d: Boost Wizard Step 4 (`/boost/create`, step 4)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Review summary | Content | 无链接 (展示) |
| Terms checkbox | Content | (action) |
| "Back" | Bottom | → B30c |
| "Submit for Review" | Bottom | (API) → B28 (success redirect) |

### B31: Sectors & Tasks (`/community/sectors`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Add Sector" button | Header | (action) 新增 sector |
| "Add Task" button (per sector) | Sector header | (action) 新增 task |
| Task drag handles | Task rows | (action) 拖拽排序 |
| Task status toggle | Task row | (API) 切换 active/draft/hidden |
| Task type badge | Task row | 无链接 (展示) |
| Sector visibility toggle | Sector header | (API) 隐藏/显示 sector |

### B31a-B31h: Community Module Management Pages

> 8 个模块管理页面，共享相同的页面模式：Header + Stats (4 cards) + Filter Tabs + Data Table + Pagination。
> URL 模式: `/community/modules/:type`

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create / Add" button | Header | (action/modal) 创建模块实例 |
| Filter tabs (All/Active/Draft/Archived) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑配置 |
| Pagination | Bottom | (action) |

### B31i: Badges (`/community/modules/badges`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Badge" | Header | (action/modal) 创建新 badge（名称、图标、颜色、获取条件） |
| Filter tabs (All/Active/Draft/Archived) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑 badge 配置 |
| Badge category filter | Table header | (action) 按类别筛选（Achievement/Engagement/Special） |
| Pagination | Bottom | (action) |

### B49: Access Rules (`/community/settings/access-rules`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Create Rule" | Header | (action/modal) 创建新访问规则 |
| Filter tabs (All/Active/Paused/Archived) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑规则（条件类型：Token Gate/NFT/Level/Invite） |
| Rule toggle (enable/disable) | Table row | (API) `PUT /api/community/settings/access-rules/:id` |
| Pagination | Bottom | (action) |

### B50: Homepage Editor (`/community/settings/homepage`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "+ Add Section" | Header | (action/modal) 添加新 section（类型：Banner/Widget/Custom） |
| Filter tabs (All/Visible/Hidden) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑 section（名称、类型、可见性、排序） |
| Drag reorder | Table rows | (action) 拖拽排序 section 顺序 |
| Visibility toggle | Table row | (API) `PUT /api/community/settings/homepage/:id` |
| "Preview" | Header secondary | → B33 (Preview Mode) |
| Pagination | Bottom | (action) |

### B32: Content Management (`/community/content`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Preview C-End" button | Header | → B33 (Preview Mode) |
| "Add Announcement" | Announcements | (action/modal) 创建公告 |
| Announcement pin/unpin | Announcement item | (API) |
| "Add Featured" | Featured Slots | (action/modal) 添加推荐位 |
| Featured slot edit/remove | Featured card | (action/API) |
| Module status "Configure" | Module Status grid | → B36 `/community/modules/:type` |

### B33: Preview Mode (`/community/preview`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| "Exit Preview" | Preview banner | → B32 (返回 Content Management) |
| Desktop/Mobile toggle | Preview banner | (action) 切换预览模式 |
| Embedded C-end preview | Preview frame | 无链接 (展示，仅预览) |

### B34: Community Wizard Step 2 (`/community/create`, step 2)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/community/drafts` |
| Module toggles (9 modules) | Left panel | (action) 启用/禁用模块 |
| Module inline config expand | Left panel | (action) 展开模块设置 |
| Summary preview | Right panel | (action) 实时汇总 |
| "Back" | Bottom | → B13 (Community Step 1) |
| "Next: Quick Setup" | Bottom | → B35 (Community Step 3 — Quick Setup) |

### B35: Community Wizard Step 3 — Quick Setup (`/community/create`, step 3)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Module quick-setup cards (expand/collapse) | Left panel | (action) 展开模块编辑 |
| Inline edit fields (task names, point values, etc.) | Left panel | (action) 编辑模板内容 |
| Summary checklist | Right panel | 无链接 (展示) |
| "Back" | Bottom | → B34 (Community Step 2) |
| "Next: Preview & Publish" | Bottom | → B55 (Community Step 4 — Preview & Publish) |

### B55: Community Wizard Step 4 — Preview & Publish (`/community/create`, step 4)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/community/drafts` |
| C-end preview mock (Home/Quests/Leaderboard tabs) | Preview area | (action) 预览各标签页 |
| Readiness checklist items (✅/⚠) | Right panel | 无链接 (展示) |
| Community URL "Copy" button | Right panel | (action) 复制链接 |
| "After publishing" info card | Right panel | 无链接 (展示) |
| "Back" | Bottom | → B35 (Community Step 3 — Quick Setup) |
| "Publish Community" | Bottom | (API) `POST /api/community/publish` → B10 (success redirect to Guided) |

### B36: Community Module Detail (`/community/modules/:type`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Community" | Breadcrumb | → B10/B11/B12 |
| Module settings form | Content | (action) |
| "Save Changes" | Bottom | (API) |
| "Reset to Default" | Bottom | (action) |

### B37: WL Wizard Step 1 — Choose Path (`/white-label/setup`, step 1)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/white-label/drafts` |
| Deployment path cards (Embed/Domain/SDK) | Content | (action) 选择部署路径 |
| "Cancel" | Bottom | → B14/B15/B16 (返回 WL Hub) |
| "Next: Configure" | Bottom | → B17 (WL Step 2) |

### B38: WL Wizard Step 3 — Brand (`/white-label/setup`, step 3)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Branding customization fields (logo, colors, typography) | Content | (action) |
| Live preview (path-adaptive) | Right panel | (action) 实时预览 |
| "Back" | Bottom | → B17 (WL Step 2) |
| "Next: Preview" | Bottom | → B56 (WL Step 4 — Preview) |

### B56: WL Wizard Step 4 — Preview (`/white-label/setup`, step 4)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/white-label/drafts` |
| Deployment preview mock (path-adaptive: widgets/portal/API) | Preview area | (action) 预览部署效果 |
| Readiness checklist items (✅/⚠) | Right panel | 无链接 (展示) |
| "After publishing" info card | Right panel | 无链接 (展示：Widget Library→B20, Page Builder→B23, Dev Kit→B48, Analytics→B43) |
| "Back" | Bottom | → B38 (WL Step 3 — Brand) |
| "Publish White Label" | Bottom | (API) `POST /api/white-label/publish` → B15 (success redirect to WL Active) |

### B39: Boost Campaign Detail (`/boost/:id`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Boost" | Breadcrumb | → B28/B29 |
| Campaign metrics | Content | 无链接 (展示) |
| "Pause Campaign" | Actions | (API) |
| "Increase Budget" | Actions | (action/modal) |
| "Download Report" | Actions | (action) 下载 |

### B40: WL Brand Settings (`/white-label/brand`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| Logo upload | Form | (action) 上传 |
| Color picker | Form | (action) |
| Typography settings | Form | (action) |
| "Save Changes" | Bottom | (API) |
| "Preview" | Bottom | → B33 (Preview Mode) |

### B41: WL SDK & API (`/white-label/sdk`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| "Copy" API key | API Keys | (action) 复制 |
| "Regenerate" API key | API Keys | (API) |
| SDK code snippets | Code blocks | (action) 复制 |
| "View Full Documentation" | Bottom | → (ext) API 文档 |

### B42: WL Iframe Embed (`/white-label/embed/iframe`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Embed Options" | Breadcrumb | → B19 |
| Iframe config form | Content | (action) |
| "Copy" embed code | Code block | (action) 复制 |
| Live preview | Right panel | (action) 实时预览 |
| "Save Configuration" | Bottom | (API) |

### B43: WL Page Analytics (`/white-label/pages/:id/analytics`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Pages" | Breadcrumb | → B25 |
| Date range picker | Header | (action) |
| "Export" | Header | (action) 下载 |
| Chart interactions | Content | (action) hover/zoom |
| "Edit Page" | Header | → B24 |

### B44: WL Integration Config (`/white-label/integrations/:type`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Integrations" | Breadcrumb | → B26 |
| Config form fields | Content | (action) |
| "Test Connection" | Form | (API) |
| "Save" | Bottom | (API) |
| "Disconnect" | Bottom | (API) |

### B51: WL Contract Registry (`/white-label/contracts`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5（WL 子菜单） |
| "Register Contract" | Header | (action/modal) 添加合约地址（chain + address + ABI） |
| Filter tabs (All/Active/Pending/Expired) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 查看/编辑合约详情 |
| "Verify" button | Table row | (API) `POST /api/wl/contracts/:id/verify` |
| Pagination | Bottom | (action) |

### B52: WL Activity Rule Builder (`/white-label/rules`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5（WL 子菜单） |
| "Create Rule" | Header | (action/modal) 可视化规则编辑器 |
| Filter tabs (All/Active/Draft/Paused) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑规则（触发条件 + 动作 + 频率） |
| Rule toggle (enable/disable) | Table row | (API) `PUT /api/wl/rules/:id` |
| Pagination | Bottom | (action) |

### B53: WL Privilege Manager (`/white-label/privileges`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5（WL 子菜单） |
| "Create Tier" | Header | (action/modal) 创建权限层级 |
| Filter tabs (All/Active/Draft) | Filter bar | (action) 前端筛选 |
| Table row click | Table | (action/modal) 编辑层级（名称、条件、权限列表） |
| "Manage Members" | Table row | (action/modal) 查看/管理层级成员 |
| Pagination | Bottom | (action) |

### B54: Community Insights (`/community/insights`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Date range picker | Header | (action) |
| Module filter tabs | Filter bar | (action) 按模块筛选 |
| "Export Report" | Header | (action) 下载 CSV/PDF |
| Segment cards | User Segments | (action/modal) 查看用户分群详情 |
| Economy health indicators | Economy panel | (info only) |
| Retention curve click | Chart | (action) 查看详细留存数据 |

### B45: Analytics Dashboard (`/analytics`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Date range picker | Header | (action) |
| Product filter tabs | Filter bar | (action) |
| "Export Report" | Header | (action) 下载 |
| Campaign links in table | Table | → `/quest/:id` or B39 |

### B46: Settings (`/settings`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Profile tab | Tabs | → B47 |
| Team tab | Tabs | (action) |
| Billing tab | Tabs | (action) |
| API tab | Tabs | → B41 (redirect) |
| "Save Changes" | Per section | (API) |

### B47: Settings / Profile (`/settings/profile`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Settings" | Breadcrumb | → B46 |
| Avatar upload | Form | (action) |
| Profile form fields | Form | (action) |
| "Save Profile" | Bottom | (API) |
| "Change Password" | Security | (action/modal) |

---

### C01: Community Home (`/`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| "Connect Wallet" | Header | (action) 钱包连接弹窗 |
| Announcement banner | Announcements | → (ext) 公告详情 |
| Featured card click | Featured grid | → (ext) 活动详情 |
| DayChain day click | DayChain strip | (action) 签到 |
| Task "Start" / "Claim" | Task cards | (action/API) 开始/领取 |
| Footer | Footer | 见 §2.7 |

### C02: Quest Tab (`/quests`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Filter pills (All/Available/Completed) | Filter bar | (action) 前端筛选 |
| Quest card "Start Quest" | Quest card | (action/API) 开始 quest |
| Quest card "View Details" | Quest card | → `/quests/:id` (quest detail, in-app) |
| Footer | Footer | 见 §2.7 |

### C03: Leaderboard (`/leaderboard`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Point Type selector (EXP/GEM/etc) | Title row | (action) 切换积分类型 |
| Time filter pills (Weekly/Monthly/All Time) | Filter bar | (action) 切换时段 |
| Podium user click | Podium | → `/profile/:address` (user profile, future) |
| Table row user click | Rankings table | → `/profile/:address` (user profile, future) |
| Footer | Footer | 见 §2.7 |

### C04: LB Sprint Tab (`/lb-sprint`)

> Leaderboard Sprint = 有开始/结束日期的限时排行榜竞赛，基于自定义积分类型（EXP/GEM），带非积分类激励（NFT/Token/WL Spot）。区别于 Leaderboard（C03），后者为周期性积分展示，无额外激励。

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Sprint task "Continue" | In-progress task | (action/API) 继续任务，赚取对应积分 |
| Sprint task "Start" | Not-started task | (action/API) 开始任务，赚取对应积分 |
| Reward tier cards (EXP/incentive) | Reward Tiers | 无链接 (展示: +100 EXP / +300 EXP / OG NFT for top 10) |
| Past LB Sprint card | Past LB Sprints | → `/lb-sprint/:id` (sprint detail, future) |
| Footer | Footer | 见 §2.7 |

### C05: Milestone Tab (`/milestones`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| "Claim Reward" | Claimable milestone card | (API) 领取奖励 |
| Locked milestone card | Locked cards | 无链接 (展示，灰色) |
| Claimed milestone card | Claimed cards | 无链接 (展示，绿色) |
| Footer | Footer | 见 §2.7 |

### C06: Shop Tab (`/shop`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Category filter pills (All/NFTs/Vouchers/Merch/Whitelist) | Filter bar | (action) 筛选 |
| "Redeem" button | Affordable items | (API) 兑换 → 确认弹窗 |
| "Not Enough" button | Unaffordable items | 无链接 (disabled) |
| "Sold Out" button | Sold out items | 无链接 (disabled) |
| Footer | Footer | 见 §2.7 |

### C07: User Center (`/profile`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Profile card stats | Profile card | 展示 (Total Points / Tasks Done / Quests / Best Streak) |
| Achievement badges | Achievements grid | 展示 (earned = colored, locked = gray + lock icon) |
| Activity history items | Recent Activity list | 展示 (color-coded dots + point changes) |
| "Invite More Friends" button | Referral Stats | → C08 `/invite` |
| Footer | Footer | 见 §2.7 |

### C08: Invite Center (`/invite`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| "Copy Link" button | Referral link row | (action) 复制邀请链接到剪贴板 |
| "Twitter" share button | Share row | (action) 打开 Twitter 分享窗口 |
| "Telegram" share button | Share row | (action) 打开 Telegram 分享窗口 |
| "Discord" share button | Share row | (action) 打开 Discord 分享窗口 |
| Reward tier cards | Reward Tiers | 展示 (Starter 3人/Ambassador 10人/Legend 50人) |
| Top Inviters list | Leaderboard | 展示 (rank + avatar + address + count) |
| "You" highlight row | Leaderboard | 展示 (amber bg, user's rank) |
| Footer | Footer | 见 §2.7 |

### C09: Activity Feed (`/activity`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Filter chips (All/Tasks/Rewards/Level Up/Invites) | Filter bar | (action) 筛选活动类型 |
| Activity feed items | Feed list | 展示 (color-coded icons: green=task, amber=level, purple=invite, red=wheel, blue=milestone, pink=shop) |
| "Live" indicator | Feed header | 展示 (green dot + auto-refresh 15s) |
| Trending cards | Trending Now | 展示 (Hot Streak / LB Sprint Rush / New Reward) |
| Footer | Footer | 见 §2.7 |

---

## 四、数据刷新策略

| 策略 | 定义 | 典型缓存 |
|------|------|----------|
| **Static** | 硬编码到前端代码，无需 API | ∞ (部署时更新) |
| **Semi-static** | 低频变动，CDN 缓存 | 1h - 24h |
| **Dynamic** | 用户操作频繁变化 | 0 - 60s TTL |
| **Real-time** | 需 WebSocket 推送 | 实时 |

---

## 五、API 汇总表

### Marketing APIs

| Endpoint | Method | Pages | Cache | Priority |
|----------|--------|-------|-------|----------|
| `/api/stats/weekly-earnings` | GET | M01 Hero | 300s | P0 |
| `/api/stats/hourly-earnings` | GET | M01 Happening Now | 60s | P0 |
| `/api/stats/online` | GET/WS | M01 Hero | 10s | P1 |
| `/api/stats/platform` | GET | M01/M02/M10 | 86400s | P0 |
| `/api/stats/active-campaigns` | GET | M01 Value Cards | 300s | P2 |
| `/api/quests/featured` | GET | M01 Hero 任务卡 | 60s | P0 |
| `/api/claims/latest` | GET | M01 Hero 奖励弹窗 | 30s | P1 |
| `/api/activity/recent` | GET/WS | M01 Live Feed | 15s | P0 |
| `/api/leaderboard/weekly` | GET | M01 Leaderboard | 300s | P0 |
| `/api/partners/featured` | GET | M01 Trust Bar | 86400s | P2 |
| `/api/partners/all` | GET | M02 Logo Wall | 86400s | P2 |
| `/api/testimonials/featured` | GET | M02 | 86400s | P2 |
| `/api/pricing/plans` | GET | M07 (Unified Pricing, replaces M07+M08) | 3600s | P1 |
| `/api/contact` | POST | M09 | N/A | P1 |
| `/api/case-studies` | GET | M11 | 86400s | P2 |
| `/ws/activity` | WS | M01 Live Feed | real-time | P1 |
| `/ws/presence` | WS | M01 在线人数 | real-time | P2 |

### B-End APIs

| Endpoint | Method | Pages | Cache | Priority |
|----------|--------|-------|-------|----------|
| `/api/user/profile` | GET | Top Bar, B01-B03, B47 | session | P0 |
| `/api/notifications/unread-count` | GET/WS | Top Bar | real-time | P1 |
| `/api/quest/stats` | GET | B05, B06 | 60s | P1 |
| `/api/quest/campaigns` | GET | B05, B06 | 30s | P0 |
| `/api/quest/templates` | GET | B07 | 3600s | P0 |
| `/api/quest/task-types` | GET | B08 | 3600s | P1 |
| `/api/quest/drafts` | POST | B07, B08 | N/A | P1 |
| `/api/community/stats` | GET | B10-B12 | 60s | P1 |
| `/api/community/list` | GET | B10-B12 | 30s | P1 |
| `/api/community/drafts` | POST | B13, B34, B35, B55 | N/A | P1 |
| `/api/community/publish` | POST | B55 | N/A | P0 |
| `/api/community/sectors` | GET | B31 | 30s | P1 |
| `/api/community/sectors` | PUT | B31 | N/A | P1 |
| `/api/community/content` | GET | B32 | 30s | P1 |
| `/api/community/modules/:type` | GET | B36, B31a-B31h | 60s | P1 |
| `/api/community/modules/:type` | PUT | B36, B31a-B31h | N/A | P1 |
| `/api/community/modules/:type/instances` | GET | B31a-B31h | 30s | P1 |
| `/api/community/modules/:type/instances` | POST | B31a-B31h | N/A | P1 |
| `/api/promo-kit/generate` | POST | B10, B15 | N/A | P2 |
| `/api/devkit/:project_id` | GET | B48 | 60s | P2 |
| `/api/devkit/:project_id/verify` | POST | B48 | N/A | P2 |
| `/api/white-label/status` | GET | B15, B16 | 60s | P1 |
| `/api/white-label/widgets` | GET | B20-B22 | 60s | P1 |
| `/api/white-label/pages` | GET | B23-B25 | 60s | P1 |
| `/api/white-label/pages/:id/analytics` | GET | B43 | 60s | P1 |
| `/api/white-label/integrations` | GET | B26 | 60s | P1 |
| `/api/white-label/integrations/:type` | GET/PUT | B44 | 60s | P1 |
| `/api/white-label/verify-domain` | POST | B17, B18 | N/A (轮询) | P0 |
| `/api/white-label/embed-code` | GET | B17, B21, B42 | 0 (实时) | P1 |
| `/api/white-label/brand` | GET/PUT | B40 | 60s | P1 |
| `/api/white-label/sdk` | GET | B41 | 60s | P1 |
| `/api/white-label/drafts` | POST | B17, B37, B38, B56 | N/A | P1 |
| `/api/white-label/publish` | POST | B56 | N/A | P0 |
| `/api/boost/stats` | GET | B28, B29 | 60s | P1 |
| `/api/boost/campaigns` | GET | B28, B29 | 30s | P1 |
| `/api/boost/campaigns/:id` | GET | B39 | 30s | P1 |
| `/api/boost/channels` | GET | B30b | 3600s | P1 |
| `/api/boost/cpa-estimate` | GET | B30c | 0 (实时) | P1 |
| `/api/boost/drafts` | POST | B30a-B30d | N/A | P1 |
| `/api/boost/submit` | POST | B30d | N/A | P0 |
| `/api/analytics/overview` | GET | B45 | 60s | P1 |
| `/api/analytics/export` | GET | B45 | N/A | P2 |
| `/api/settings` | GET/PUT | B46 | session | P1 |

### C-End APIs

| Endpoint | Method | Pages | Cache | Priority |
|----------|--------|-------|-------|----------|
| `/api/c/community/home` | GET | C01 | 60s | P0 |
| `/api/c/community/announcements` | GET | C01 | 60s | P1 |
| `/api/c/community/featured` | GET | C01 | 60s | P1 |
| `/api/c/community/daychain` | GET/POST | C01 | session | P0 |
| `/api/c/community/tasks` | GET | C01 | 30s | P0 |
| `/api/c/quests` | GET | C02 | 30s | P0 |
| `/api/c/leaderboard` | GET | C03 | 300s | P1 |
| `/api/c/lb-sprint/current` | GET | C04 | 60s | P0 |
| `/api/c/lb-sprint/history` | GET | C04 | 300s | P1 |
| `/api/c/milestones` | GET | C05 | 60s | P0 |
| `/api/c/milestones/:id/claim` | POST | C05 | N/A | P0 |
| `/api/c/shop/items` | GET | C06 | 60s | P0 |
| `/api/c/shop/redeem` | POST | C06 | N/A | P0 |
| `/api/c/user/status` | GET | C01-C09 User Status Bar | session | P0 |
| `/api/c/user/profile` | GET | C07 | session | P0 |
| `/api/c/user/achievements` | GET | C07 | 300s | P1 |
| `/api/c/user/activity` | GET | C07, C09 | 30s | P0 |
| `/api/c/user/referral-stats` | GET | C07, C08 | 60s | P1 |
| `/api/c/invite/link` | GET | C08 | session | P0 |
| `/api/c/invite/leaderboard` | GET | C08 | 300s | P1 |
| `/api/c/invite/tiers` | GET | C08 | 3600s | P2 |
| `/api/c/activity/feed` | GET | C09 | 15s | P0 |
| `/api/c/activity/trending` | GET | C09 | 300s | P1 |
| `/api/c/wallet/connect` | POST | Header | N/A | P0 |

---

## 六、新增功能组件规格

### 6.1 Promo Kit Generator（推广素材生成器）

> 用于 B10 Community Checklist Step 4、B15 WL Checklist Step 5

| 属性 | 规格 |
|------|------|
| 触发 | 点击 "Generate Promo Kit" 按钮 |
| UI | Modal / 侧拉面板 |
| 输入 | 项目名称、描述、已启用模块、品牌色（从 Community/WL 配置自动读取） |
| AI 输出 | 3 条平台特定社交文案 (Twitter/Discord/Telegram) + 品牌推广 Banner |
| Banner 尺寸 | 1200×675 (Twitter/OG) + 1080×1080 (IG/TG) |
| 操作 | 每条文案: "Copy" + "Share" (直接分享到对应平台) |
| API | `POST /api/promo-kit/generate` → 返回文案数组 + 图片 URL |

### 6.2 Dev Kit Page（开发者集成页）

> 用于 B15 WL Checklist Step 3。独立页面，无需 TaskOn 账号即可访问。

| 属性 | 规格 |
|------|------|
| URL | `taskon.xyz/devkit/{project_id}` |
| 用户 | 项目方的开发人员（由市场人员发送链接） |
| 页面编码 | **B48** |
| 内容 | 集成代码 (copy-paste, project ID 预填) + SSO 设置 (wallet/OAuth) + 分步指南 (预计 30 min) + "Verify My Integration" 按钮 |
| 设计状态 | 概念设计完成（BnkYW checklist 引用），页面待设计 |
| API | `GET /api/devkit/{project_id}` (获取项目配置) + `POST /api/devkit/{project_id}/verify` (验证集成) |

### 6.3 Auto-detection Systems（自动检测系统）

> 用于 B10 Step 5、B15 Steps 4 & 6

| 检测类型 | 触发条件 | 协议 | 用途 |
|----------|----------|------|------|
| Integration Ping | 首次从项目域名发起的 API 调用 | WebSocket | B15 Step 4 (Integration verified) |
| DNS Verification | 自定义域名 DNS 配置生效 | HTTP 轮询 | B18 Domain Setup |
| Participant Counter | 用户加入社区 | WebSocket 实时计数 | B10 Step 5 (First 10 participants)、B15 Step 6 (First user interaction) |

---

## 七、SKIPPED 页面说明

以下页面在设计审查中被**有意跳过**，因为它们使用与已设计页面相同的系统模式：

| # | 页面 | 原因 |
|---|------|------|
| T01 | Quest Campaign Detail (`/quest/:id`) | 复用 Boost Campaign Detail (B39) 相同模式 |
| T02 | Quest Campaign Edit (`/quest/:id/edit`) | 复用 Wizard Step 2 (B08) 编辑模式 |
| T03 | Quest Wizard Step 3: Set Rewards | 复用已有系统奖励设置流程 |
| T04 | Quest Wizard Step 4: Review & Launch | 复用 Boost Step 4 (B30d) 相同模式 |
| T05 | Boost Campaign Edit (`/boost/:id/edit`) | 复用 Boost Wizard 步骤编辑模式（同 Quest Edit T02 模式），从 B28 Campaign card "Edit" 按钮进入 |
| T25 | Blog (`/blog`) | 外部链接，不需要设计 |

> **注意**: 所有其他 TODO 页面已在 P12-P14 + P0/P0b/P0c 阶段完成设计。76 个页面编码（含 B48 Dev Kit Page 待设计）。

---

## 八、开发优先级（更新版）

| Phase | 内容 | 页面编码 | 依赖 |
|-------|------|----------|------|
| **P1** | Brand Homepage 静态 + 全站路由 + Header/Footer | M01 | 设计稿确认 |
| **P2** | 首页动态 APIs | M01 | 后端 API |
| **P3** | Projects Landing + 4 产品页 + Boost (全静态) | M02-M06 | 设计稿确认 |
| **P4** | 统一定价页 + Contact 表单 | M07 (Unified), M09 | 后端 API + CMS |
| **P5** | WebSocket (activity, presence) | M01 | 基础设施 |
| **P6** | 辅助 Marketing 页 (About, Case Studies, Solutions) | M10-M14 | 内容就绪 |
| **P7** | B-End: 登录/认证 + Sidebar/TopBar + Dashboard 3 states | B01-B03 | 认证系统 |
| **P8** | B-End: Quest 全线 (Hub 3 states + Wizard 2 steps) | B04-B08 | Quest API |
| **P9** | B-End: Community 全线 (Hub 4 states + Wizard 4 steps + Module Mgmt ×9 + Content) | B09-B13, B31-B36, B31a-B31h, B55 | Community API (Wizard: B13→B34→B35→B55) |
| **P10** | B-End: WL 全线 (Hub 3 states + Wizard 4 steps + 10 sub-pages) | B14-B26, B37-B44, B56 | WL API + DNS (Wizard: B37→B17→B38→B56) |
| **P11** | B-End: Boost 全线 (Hub 3 states + Wizard 4 steps + Campaign Detail) | B27-B30d, B39 | Boost API |
| **P12** | B-End: Analytics + Settings | B45-B47 | 全部 API |
| **P13** | B-End: Preview Mode | B33 | Community + WL API |
| **P14** | C-End: White Label 社区前台 (6 tabs) | C01-C06 | C-End API + Wallet |
| **P15** | C-End: User engagement pages (Profile/Invite/Activity) | C07-C09 | C-End User + Invite + Activity APIs |

---

## 九、页面状态路由策略

部分 URL 对应多个设计状态（根据用户数据动态选择）:

| URL | 状态判断条件 | 显示页面 |
|-----|-------------|----------|
| `/` (Dashboard) | 0 campaigns | B01 (New User) |
| `/` (Dashboard) | 1-5 campaigns | B02 (Active) |
| `/` (Dashboard) | 6+ campaigns | B03 (Power) |
| `/quest` | 0 campaigns | B04 (Empty) |
| `/quest` | 1-3 campaigns | B05 (Active) |
| `/quest` | 4+ campaigns | B06 (Management) |
| `/community` | 0 modules | B09 (Empty) |
| `/community` | Onboarding incomplete | B10 (Guided) |
| `/community` | 1-3 active modules | B11 (Active) |
| `/community` | 4+ modules, high usage | B12 (Deep) |
| `/white-label` | 0 tools configured | B14 (Empty) |
| `/white-label` | 1-4 tools configured | B15 (Active) |
| `/white-label` | 5+ tools, high traffic | B16 (Management) |
| `/boost` | 0 campaigns | B27 (Empty) |
| `/boost` | 1-3 campaigns | B28 (Active) |
| `/boost` | 4+ campaigns | B29 (Management) |
| `/white-label/widgets` | 0 widgets created | B20 (Empty) |
| `/white-label/widgets` | 1+ widgets created | B22 (Active) |
| `/white-label/pages` | 0 pages created | B23 (Empty) |
| `/white-label/pages` | 1+ pages created | B25 (Active) |
