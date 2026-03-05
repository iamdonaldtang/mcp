# TaskOn 官网前端开发需求文档 v4.0

> 本文档定义了 TaskOn 官网（marketing pages）、B-end 产品界面（after-login dashboard & product hubs）、及 C-end 用户界面（White Label 社区前台）的完整页面编码、按钮→页面路由映射、动态数据规格、及静态内容边界。
> 供前端/后端工程师对照实施。
>
> v4.0 更新：新增 B31-B47（17 个 B-end 页面）、M09-M14（6 个 Marketing 页面）、C01-C06（6 个 C-end 页面），共 67 页面编码。完成全量按钮→页面路由映射，清除已设计页面的 TODO 标记。仅剩 T01-T04（Quest Campaign Detail/Edit/Step3/Step4）为 SKIPPED。

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
| **M07** | Platform Pricing | `CXtOH` | `/pricing` | Light |
| **M08** | White Label Pricing | `EDoSn` | `/pricing/enterprise` | Light, purple |
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
| **B09** | Community Hub (Empty) | `zzZ8D` | `/community` | Templates + highlight strip |
| **B10** | Community Hub (Guided) | `S1EIA` | `/community` | Onboarding checklist + modules |
| **B11** | Community Hub (Active) | `vFRHi` | `/community` | Quick stats + checklist + modules |
| **B12** | Community Hub (Deep) | `TQR51` | `/community` | 6 modules + engagement chart |
| **B13** | Community Wizard Step 2 | `Gzpeu` | `/community/create` (step 2) | Customize: name, color, preview |
| **B14** | White Label Hub (Empty) | `Ir6Tq` | `/white-label` | 3 goals + 6 toolkit cards |
| **B15** | White Label Hub (Active) | `BnkYW` | `/white-label` | Stats + 6 toolkit (4 configured) |
| **B16** | White Label Hub (Mgmt) | `UPAfV` | `/white-label` | All tools + deployments + analytics |
| **B17** | WL Wizard Step 2 | `CXzmy` | `/white-label/setup` (step 2) | Domain + SSL + embed code |
| **B18** | WL Domain Setup | `5bmH9` | `/white-label/domain` | 3-step domain config + status |
| **B19** | WL Embed Options | `RgCVQ` | `/white-label/embed` | 3 embed mode cards + comparison |
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
| **B32** | Content Management | `lhR14` | `/community/content` | Announcements + featured + modules |
| **B33** | Preview Mode | `2UiNC` | `/community/preview` | Embedded C-end preview |
| **B34** | Community Wizard Step 1 | `cNwNP` | `/community/create` (step 1) | Choose community type |
| **B35** | Community Wizard Step 3 | `Sq7A2` | `/community/create` (step 3) | Set modules |
| **B36** | Community Module Detail | `usBsM` | `/community/modules/:type` | Module config (Points example) |
| **B37** | WL Wizard Step 1 | `NNwid` | `/white-label/setup` (step 1) | Choose integration mode |
| **B38** | WL Wizard Step 3 | `5nCtO` | `/white-label/setup` (step 3) | Customize branding |
| **B39** | Boost Campaign Detail | `Sq4jV` | `/boost/:id` | Campaign metrics + management |
| **B40** | WL Brand Settings | `Cx3LH` | `/white-label/brand` | Logo, colors, typography |
| **B41** | WL SDK & API | `lQxT5` | `/white-label/sdk` | API keys, SDK docs, endpoints |
| **B42** | WL Iframe Embed | `ByGS0` | `/white-label/embed/iframe` | Iframe embed flow |
| **B43** | WL Page Analytics | `69HPh` | `/white-label/pages/:id/analytics` | Page performance metrics |
| **B44** | WL Integration Config | `gS64G` | `/white-label/integrations/:type` | Integration setup (template) |
| **B45** | Analytics Dashboard | `fLxTr` | `/analytics` | Cross-product analytics |
| **B46** | Settings | `ESrVt` | `/settings` | Account settings |
| **B47** | Settings / Profile | `Nh7xq` | `/settings/profile` | Profile management |

### C-End Pages（White Label 社区前台，C01-C06）

> C-End 页面为 White Label Community 的用户端。域名由项目自定义（如 `community.bitcoin.com`）。Theme: Dark header (#0F172A) + Light content (#F8FAFC)。Accent: Amber (#F59E0B)。

| Code | Page | Design ID | URL Path | Notes |
|------|------|-----------|----------|-------|
| **C01** | Community Home | `vJVhd` | `/` | Announcements + featured + DayChain + task sectors |
| **C02** | Quest Tab | `dUXTl` | `/quests` | Quest card grid 2×2 + filters |
| **C03** | Leaderboard | `KmdSd` | `/leaderboard` | Podium top 3 + rankings table |
| **C04** | Sprint Tab | `y5fUZ` | `/sprint` | Current sprint + tasks + reward tiers + past sprints |
| **C05** | Milestone Tab | `53iKE` | `/milestones` | Level progress + milestone cards (claimed/claimable/locked) |
| **C06** | Shop Tab | `coM7o` | `/shop` | Rewards shop + 6 items (redeem/not enough/sold out) |

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
| White Label | → B14/B15/B16 `/white-label` |
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

用于 C01-C06。Dark theme (#0F172A)。

| 元素 | 链接目标 |
|------|----------|
| Project Logo + Name | → C01 `/` |
| "Connect Wallet" | (action) 钱包连接弹窗 |
| Nav: Home | → C01 `/` |
| Nav: Quests | → C02 `/quests` |
| Nav: Leaderboard | → C03 `/leaderboard` |
| Nav: Sprint | → C04 `/sprint` |
| Nav: Milestone | → C05 `/milestones` |
| Nav: Shop | → C06 `/shop` |

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

### M07: Platform Pricing (`/pricing`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1 |
| Monthly/Annual toggle | Billing | (action) 前端切换 |
| Quest/Community tabs | Tabs | (action) 前端切换 |
| Free "Get Started Free" | Pricing card | → `app.taskon.xyz/create` (ext) |
| Pro "Start Free Trial" | Pricing card | → `app.taskon.xyz/create?plan=pro` (ext) |
| "Get White Label Pricing" | WL upsell banner | → M08 |
| FAQ expand/collapse | FAQ | (action) |
| Footer links | Footer | 见 §2.3 |

### M08: White Label Pricing (`/pricing/enterprise`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header links | Header | 见 §2.1（"Talk to Sales" → M09） |
| "← Back to Platform Pricing" | Hero | → M07 |
| Standard "Start 14-Day Free Trial" | Pricing card | → `app.taskon.xyz/create?plan=wl-standard` (ext) |
| Pro "Book a Demo" | Pricing card | → M09 `/contact` |
| "Book a Demo" | CTA section | → M09 `/contact` |
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
| "Create Your First Community" | 主CTA | → B34 (Step 1) |
| "or start from a blank canvas" | Sub-CTA | → B34 `/community/create?template=blank` |
| Video Tutorial | Resources | → (ext) 视频 |
| Retention Playbook | Resources | → (ext) 帮助中心 |
| "Learn More" | Resources | → M04 (ext, new tab) |

### B10: Community Hub — Guided Workspace (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Checklist items (5 steps) | Checklist card | (action) 展开/完成步骤 |
| Module cards (Points, Leaderboard, Tasks) | Active Modules | → B36 `/community/modules/:type` |
| "Add" module (Milestones, Benefits Shop) | Add More Modules | (action) 启用模块 |
| Template row link | Resources | → B34 |
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
| Template row link | Resources | → B34 |
| Resources | Resources | → (ext) |

### B12: Community Hub — Deep (`/community`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
| Module cards (6, 2 rows) | Modules | → B36 `/community/modules/:type` |
| Engagement Overview chart | Analytics | 无链接 (展示) |
| Retention metrics | Analytics | → B45 `/analytics` |

### B13: Community Wizard Step 2 (`/community/create`, step 2)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/community/drafts` |
| Community Name input | Form | (action) |
| Description textarea | Form | (action) |
| Brand Color picker | Form | (action) |
| Live Preview | Right panel | (action) 实时预览 |
| "Back" | Bottom | → B34 (Community Step 1) |
| "Next: Points & Rewards" | Bottom | → B35 (Community Step 3) |

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
| Sidebar / Top bar | 全局 | 见 §2.4, §2.5 |
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

### B19: WL Embed Options (`/white-label/embed`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to White Label" | Breadcrumb | → B14/B15/B16 |
| "Iframe Embed" card | Embed modes | → B42 `/white-label/embed/iframe` |
| "Open Widget Library" ★ card | Embed modes | → B20 (Widget Library) |
| "Open Page Builder" card | Embed modes | → B23 (Page Builder) |
| Comparison table | Content | 无链接 (展示) |
| Tip card | Bottom | 无链接 |

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

### B34: Community Wizard Step 1 (`/community/create`, step 1)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/community/drafts` |
| Community type cards | Content | (action) 选择类型 |
| "Cancel" | Bottom | → B09/B10/B11/B12 (返回 Community Hub) |
| "Next: Customize" | Bottom | → B13 (Community Step 2) |

### B35: Community Wizard Step 3 (`/community/create`, step 3)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Module toggles | Content | (action) 启用/禁用模块 |
| Module config expand | Content | (action) 展开模块配置 |
| "Back" | Bottom | → B13 (Community Step 2) |
| "Launch Community" | Bottom | (API) → B10 (success redirect to Guided Workspace) |

### B36: Community Module Detail (`/community/modules/:type`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "← Back to Community" | Breadcrumb | → B10/B11/B12 |
| Module settings form | Content | (action) |
| "Save Changes" | Bottom | (API) |
| "Reset to Default" | Bottom | (action) |

### B37: WL Wizard Step 1 (`/white-label/setup`, step 1)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) `POST /api/white-label/drafts` |
| Integration mode cards | Content | (action) 选择集成模式 |
| "Cancel" | Bottom | → B14/B15/B16 (返回 WL Hub) |
| "Next: Configure" | Bottom | → B17 (WL Step 2) |

### B38: WL Wizard Step 3 (`/white-label/setup`, step 3)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| "Save Draft" | Top bar | (API) |
| Branding customization fields | Content | (action) |
| "Back" | Bottom | → B17 (WL Step 2) |
| "Complete Setup" | Bottom | (API) → B15 (success redirect to Active Hub) |

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
| Time filter pills (This Week/Month/All Time) | Filter bar | (action) 切换时段 |
| Podium user click | Podium | → `/profile/:address` (user profile, future) |
| Table row user click | Rankings table | → `/profile/:address` (user profile, future) |
| Footer | Footer | 见 §2.7 |

### C04: Sprint Tab (`/sprint`)

| 按钮/CTA | 位置 | 目标 |
|----------|------|------|
| Header / Nav tabs | 全局 | 见 §2.6 |
| Sprint task "Continue" | In-progress task | (action/API) 继续任务 |
| Sprint task "Start" | Not-started task | (action/API) 开始任务 |
| Reward tier card | Reward Tiers | 无链接 (展示) |
| Past sprint card | Past Sprints | → `/sprint/:id` (sprint detail, future) |
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
| `/api/pricing/plans` | GET | M07, M08 | 3600s | P1 |
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
| `/api/community/drafts` | POST | B13, B34, B35 | N/A | P1 |
| `/api/community/sectors` | GET | B31 | 30s | P1 |
| `/api/community/sectors` | PUT | B31 | N/A | P1 |
| `/api/community/content` | GET | B32 | 30s | P1 |
| `/api/community/modules/:type` | GET | B36 | 60s | P1 |
| `/api/community/modules/:type` | PUT | B36 | N/A | P1 |
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
| `/api/white-label/drafts` | POST | B17, B37, B38 | N/A | P1 |
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
| `/api/c/sprint/current` | GET | C04 | 60s | P0 |
| `/api/c/sprint/history` | GET | C04 | 300s | P1 |
| `/api/c/milestones` | GET | C05 | 60s | P0 |
| `/api/c/milestones/:id/claim` | POST | C05 | N/A | P0 |
| `/api/c/shop/items` | GET | C06 | 60s | P0 |
| `/api/c/shop/redeem` | POST | C06 | N/A | P0 |
| `/api/c/user/status` | GET | C01-C06 User Status Bar | session | P0 |
| `/api/c/wallet/connect` | POST | Header | N/A | P0 |

---

## 六、SKIPPED 页面说明

以下页面在设计审查中被**有意跳过**，因为它们使用与已设计页面相同的系统模式：

| # | 页面 | 原因 |
|---|------|------|
| T01 | Quest Campaign Detail (`/quest/:id`) | 复用 Boost Campaign Detail (B39) 相同模式 |
| T02 | Quest Campaign Edit (`/quest/:id/edit`) | 复用 Wizard Step 2 (B08) 编辑模式 |
| T03 | Quest Wizard Step 3: Set Rewards | 复用已有系统奖励设置流程 |
| T04 | Quest Wizard Step 4: Review & Launch | 复用 Boost Step 4 (B30d) 相同模式 |
| T05 | Boost Campaign Edit (`/boost/:id/edit`) | 复用 Boost Wizard 步骤编辑模式（同 Quest Edit T02 模式），从 B28 Campaign card "Edit" 按钮进入 |
| T25 | Blog (`/blog`) | 外部链接，不需要设计 |

> **注意**: 所有其他 TODO 页面已在 P12-P14 阶段完成设计。67 个页面编码全部覆盖。

---

## 七、开发优先级（更新版）

| Phase | 内容 | 页面编码 | 依赖 |
|-------|------|----------|------|
| **P1** | Brand Homepage 静态 + 全站路由 + Header/Footer | M01 | 设计稿确认 |
| **P2** | 首页动态 APIs | M01 | 后端 API |
| **P3** | Projects Landing + 4 产品页 + Boost (全静态) | M02-M06 | 设计稿确认 |
| **P4** | 2 定价页 + Contact 表单 | M07, M08, M09 | 后端 API + CMS |
| **P5** | WebSocket (activity, presence) | M01 | 基础设施 |
| **P6** | 辅助 Marketing 页 (About, Case Studies, Solutions) | M10-M14 | 内容就绪 |
| **P7** | B-End: 登录/认证 + Sidebar/TopBar + Dashboard 3 states | B01-B03 | 认证系统 |
| **P8** | B-End: Quest 全线 (Hub 3 states + Wizard 2 steps) | B04-B08 | Quest API |
| **P9** | B-End: Community 全线 (Hub 4 states + Wizard 3 steps + Module Detail + Sectors + Content) | B09-B13, B31-B36 | Community API |
| **P10** | B-End: WL 全线 (Hub 3 states + Wizard 3 steps + 10 sub-pages) | B14-B26, B37-B44 | WL API + DNS |
| **P11** | B-End: Boost 全线 (Hub 3 states + Wizard 4 steps + Campaign Detail) | B27-B30d, B39 | Boost API |
| **P12** | B-End: Analytics + Settings | B45-B47 | 全部 API |
| **P13** | B-End: Preview Mode | B33 | Community + WL API |
| **P14** | C-End: White Label 社区前台 (6 tabs) | C01-C06 | C-End API + Wallet |

---

## 八、页面状态路由策略

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
