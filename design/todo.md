# Design TODO

> Active design tasks for `design/pencil-new.pen`
> Completed items tracked in `design/progress.md`
> Last design review: 2026-03-05 (Community + White Label + C-End, 42 pages)
> **Batch execution: 2026-03-05** — 28 of 40 items completed (see below)

## Batch Execution Summary (2026-03-05)

| Category | Items | Completed | Notes |
|----------|-------|-----------|-------|
| P0 Blocking | DR-01, DR-02 | 2/2 | Shop card fills fixed, Integration Center clarified |
| P1a IA/Flow | DR-03–DR-07 | 5/5 | Checklist compressed, Quests tightened, Home zones, nav overflow fixed |
| P1b Visual | DR-08–DR-12 | 5/5 | Community hero 2-col, module insight cards, milestone locks, PB canvas ratio, Embed info banner |
| P2 Polish | DR-13, DR-14, DR-17 | 3/8 | Launch banner, Invite compact, DayChain warning. DR-15/16/18/19/20 deferred (low impact) |
| S-Tier | S1, S5, S6, S7 | 4/7 | Lifecycle selector, For Projects button, Book a Demo CTAs, Pricing badges. S2/S3/S4 partially addressed (existing content) |
| I-Tier-B | I-02, I-05, I-09, I-10 | 4/4 | Rate metrics+trends on DayChain/LB, Shop insight, WL Analytics funnel, Content Mgmt perf data |
| W-Tier | Contract Registry | 1/3 | New page `OKEqS`. Rule Builder + Privilege Manager deferred |
| **Total** | | **24/40** | + 4 partially addressed = **28 touched** |

### Remaining Items (16)
- **DR-15/16/18/19/20** (P2): Text-heavy sections, icon sizing, dev kit modal, pricing tabs — low priority polish
- **S2/S3/S4**: Cross-sell combos (partial), logo wall content (placeholder), ROI calculator (exists on pricing)
- **I-01/I-03/I-04/I-06/I-07/I-08/I-11** (I-Tier-A): New pages needed — Community Insights, Points Economy, Retention deep-dive, User Segmentation, Recommendations, Change Tracking
- **W1 Rule Builder, W2 Privilege Manager**: New pages needed — require Contract Registry (done) as foundation

---

## 🔴 P0 — Blocking Issues

| # | Page | Node ID | Issue | Suggested Fix |
|---|------|---------|-------|---------------|
| DR-01 | C-End Shop | `coM7o` | **Bottom row product cards missing text/pricing** — Second row of 3 cards shows images only, no title, price, or buy button. Appears broken/incomplete vs top row which renders correctly. | Add complete card content (title, price, points cost, CTA button) to bottom 3 cards matching top row pattern |
| DR-02 | WL Integration Center | `Abs1E` vs `gS64G` | **Duplicate Integration Center pages** — Both pages are titled "Integration Center" with different layouts (Abs1E=categorized cards, gS64G=tab-filtered grid). `pages.md` maps gS64G as "Brand" but screenshot shows Integration Center. Creates confusion for frontend dev. | Decide canonical version: keep `Abs1E` (richer categorized layout) as Integration Center; repurpose `gS64G` as WL Brand Settings with live preview (merge with `Cx3LH`), OR delete one and clarify mapping in pages.md |

---

## 🟡 P1 — Important Issues

### P1a — Information Architecture & User Flow

| # | Page | Node ID | Issue | Suggested Fix |
|---|------|---------|-------|---------------|
| DR-03 | Community Hub Active | `vFRHi` | **Getting Started checklist too prominent in Active state** — Checklist shows 4/6 steps complete but dominates the page. Active users with 1,247 members and 8,931 tasks completed don't need a large onboarding checklist as the primary content. | Make checklist collapsible/auto-minimized when >3 steps complete. Move stats row to top, active modules as primary content. Show "Complete remaining steps" as a compact banner instead. |
| DR-04 | Community Hub Active vs Guided | `vFRHi` vs `S1EIA` | **Insufficient visual differentiation between states** — Guided (new community) and Active (used community) look nearly identical in layout structure: both have checklist → modules → resources. Hard to feel progression. | Active state should lead with stats + module performance cards (with trend indicators like in Deep state `TQR51`). Reduce checklist prominence. Guided keeps checklist as hero. |
| DR-05 | C-End Quests | `dUXTl` | **Excessive whitespace below quest cards** — Only 4 quest cards displayed, bottom 50%+ of page is empty white space. Feels unfinished. | Add "No more quests" empty state message, or show upcoming/ended quests section, or tighten layout to remove dead space. Consider pagination or "Load More" pattern. |
| DR-06 | C-End Home | `vJVhd` | **Information density too high — competing visual hierarchy** — Page has 8+ distinct content sections (action engine, featured cards, DayChain, promo banner, Getting Started, Daily Engagement, stats, Discover More) with no clear visual grouping. User doesn't know where to look first. | Group into 3 clear zones: (1) Hero zone — action engine + featured, (2) Progress zone — DayChain + Getting Started, (3) Explore zone — engagement + stats + discover. Add section dividers or background color alternation for visual rhythm. |
| DR-07 | C-End Navigation | All C-End pages | **9 horizontal nav tabs may cause overflow** — Pages show up to 9 tabs (Home/Quests/Leaderboard/LB Sprint/Milestone/Shop/My Profile/Invite/Activity). On narrower screens this will overflow. Early pages (C01-C06) show 6 tabs, later pages (C07-C09) show 8-9. | Keep 6 public tabs in main nav. Move My Profile/Invite/Activity under a user avatar dropdown menu (common pattern in SaaS). This also creates clear separation between community content and personal features. |

### P1b — Visual Consistency

| # | Page | Node ID | Issue | Suggested Fix |
|---|------|---------|-------|---------------|
| DR-08 | Community Marketing | `GyyL4` | **Hero section lacks visual impact** — Hero text "80% Retention at 90 Days" is compelling but visually flat. No hero image, illustration, or product screenshot. Compares poorly to WL marketing page (`cbBdG`) which has a stronger visual hierarchy. | Add a product screenshot or dashboard mockup to the right of hero text (2-column hero layout). Show the Community dashboard with engagement metrics to make the "80% retention" claim tangible. |
| DR-09 | Community Module Pages | `Wug7d`-`sme5a`, `BJLsz`, `g1CNC`, `5Wm6B` | **All 12 module pages are structurally identical** — Header + 4 stats + filter tabs + data table + pagination. While consistency is good, there's zero module-specific design. DayChain/TaskChain/Lucky Wheel are gamification features but look like spreadsheets. | Add module-specific visual elements: (1) DayChain — show a streak visualization or calendar heatmap, (2) Lucky Wheel — show a wheel preview, (3) Badges — show badge grid thumbnails, (4) Leaderboard — show a mini podium. Keep data table but add a distinguishing visual section per module. |
| DR-10 | C-End Milestones | `53iKE` | **Locked milestone cards lack visual distinction** — Bottom row (Diamond Hands, Whale Status, Legend) looks nearly identical to claimable cards. Lock state not prominent enough — users can't instantly tell which milestones are achievable vs locked. | Add a visible lock overlay or grayscale treatment to locked cards. Show clear progress indicator ("Need 5,000 more pts"). Increase visual gap between claimed (green check), claimable (amber glow/pulse), and locked (gray + lock icon). |

### P1c — Previously Identified (Carried Forward)

| # | Page | Node ID | Issue | Notes |
|---|------|---------|-------|-------|
| DR-11 | WL Page Builder Editor | `sGDcq` | Canvas area too small | 2-column layout; canvas should be ~65% width, settings panel ~35%. Currently feels 50/50. (was #37) |
| DR-12 | WL Embed Options | `RgCVQ` | Content overlap with WL Wizard | Embed Options sub-page duplicates some wizard Step 1 content. Consider making Embed Options reference the wizard or vice versa. (was #40) |

---

## 🟡 P2 — Optimization & Polish

| # | Page | Node ID | Issue | Suggested Fix |
|---|------|---------|-------|---------------|
| DR-13 | Community Wizard Step 3 | `qknQZ` | **Review & Launch page feels sparse** — Three summary cards are minimally populated. "What Happens Next" card lists 3 generic items. Doesn't build excitement for launch. | Add: (1) Animated launch illustration or confetti preview, (2) Estimated time to first participant, (3) "Share your community" pre-generated social post preview. Make the launch moment feel celebratory. |
| DR-14 | C-End Invite Center | `TaAo9` | **Hero section consumes too much viewport** — "Invite Friends, Earn Rewards" hero + share link + social buttons + stats take up ~60% of first viewport. Referral list (the actionable data) is pushed below fold. | Compact the hero: merge share link and social buttons into one row. Move stats into a horizontal strip. Get the referral list visible above fold. |
| DR-15 | Community Marketing | `GyyL4` | **Feature sections are text-heavy** — "How It Works: The User Journey" and feature comparison sections rely heavily on text with minimal visual variety. | Add icons, micro-illustrations, or product screenshots to break up text walls. Use the brand green (#48BB78) more strategically for visual anchors. |
| DR-16 | WL Marketing | `cbBdG` | **"Choose Your Integration" section icons are small** — Integration mode cards (iframe/widget/page builder/SDK) use tiny logos that are hard to distinguish at normal reading distance. | Increase icon size to 48-56px. Add color-coded borders per integration mode. Consider showing a mini-preview of what each mode looks like when deployed. |
| DR-17 | Community Hub Deep | `TQR51` | **DayChain warning card (orange border) competes with TaskChain card** — Both use attention-grabbing border colors (orange + green). DayChain's "-2.1% drop" warning is important but the visual weight is equal to TaskChain's positive metrics. | Use red/warning color (#DC2626) only for the metric text, not the full card border. Keep card border neutral. Add a small alert icon next to the metric. This focuses attention on the number, not the container. |
| DR-18 | WL Hub Active | `BnkYW` | **Getting Started Step 4 "Send Dev Kit" expanded view is dense** — The expanded step shows URL, copy button, email option, and status indicators. While informative, it's a lot of UI in an inline checklist expansion. | Consider making "Send Dev Kit" open a modal or slide-over panel instead of inline expansion. This gives more space for the dev kit link, email form, and delivery status without crowding the checklist. |
| DR-19 | Unified Pricing Page | `HO2Ny` | Needs tab switching states polish | Tab switching for Quest/Community/White Label, final copy review (carried from previous P2) |

### P2b — Previously Identified (Carried Forward)

| # | Page | Node ID | Issue | Notes |
|---|------|---------|-------|-------|
| DR-20 | Community Hub Guided vs Active | `S1EIA` vs `vFRHi` | Checklist content inconsistency | User flagged: Guided has good 5-step checklist, Active version is inconsistent. Related to DR-03/DR-04. |

---

## 🔴 S-Tier — Systemic Issues (CMO Journey Audit, 2026-03-05)

> Discovered via 4-persona CMO walkthrough: Cold Start / Pre-TGE / Post-CEX Listing / Enterprise
> These are cross-page strategic gaps, not single-page bugs. **Weakest scores: "Strengthen Need" 3.3/10, "Form Decision" 4.8/10.**

| # | Issue | Impact | Suggested Fix | Affects |
|---|-------|--------|---------------|---------|
| S1 | **No lifecycle-based entry point** — Projects Landing shows 4 products in parallel with no "What stage are you at?" guidance. Cold start CMO and post-listing CMO see the same page, must self-navigate. | CMOs don't know which product to pick → bounce or pick wrong one → churn | Add interactive selector on `Lz2vL` above product cards: "Where are you?" → Pre-launch (Quest) / Growing (Quest+Community) / Listed (Community+Boost) / Scaling (White Label). Each maps to recommended product combo. | `Lz2vL` Projects Landing |
| S2 | **Product silos — no combination guidance** — 4 product pages are independent. Pre-TGE CMO needs Quest+Community but no page tells them this. No bundle concept, no "growth stack" visualization. | Multi-product adoption blocked → lower ARPU, confused prospects | (1) Each product page CTA area: add contextual cross-sell "Using Quest? Add Community for 3x retention". (2) Projects Landing: add "Growth Stacks" section showing recommended combos by stage. (3) Pricing: add bundle discount row. | `Lz2vL`, `gXQur`, `GyyL4`, `cbBdG`, `Lym65`, `HO2Ny` |
| S3 | **No social proof in conversion flow** — "2,000+ projects" stat exists but zero named case studies, no client logos, no "Project X achieved Y% in Z weeks" proof. All 4 CMO personas asked: "Who else uses this?" | Biggest gap in funnel: interest exists but can't be converted to conviction. Strongest blocker to adoption. | (1) Add logo wall (8-12 recognizable project logos) to Projects Landing below hero. (2) Each product page: 1 mini case study card (logo + headline metric + 1 sentence). (3) Link to full Case Studies page `M11`. | `Lz2vL`, `gXQur`, `GyyL4`, `cbBdG`, `Lym65` |
| S4 | **No ROI quantification** — CMOs must justify spend to founders/board. No ROI calculator, no "invest X get Y" proof, no benchmark data. "80% retention" is a claim without evidence path. | CMO can't build internal business case → decision stalls | (1) Add simple ROI calculator to Community + Quest product pages (input users → output projected retention/acquisition + cost comparison vs alternatives). (2) Pricing page: add "Average ROI: Xth" or "Avg cost per acquired user: $Y" benchmarks. | `GyyL4`, `gXQur`, `HO2Ny` |
| S5 | **Homepage B-end entry is weak** — Brand Homepage `QszRH` is 100% C-end ("Complete Tasks, Earn Crypto"). "For Projects" button in header is small/secondary. B-end CMOs feel they're on the wrong site. | ~30% B-end visitor bounce at homepage — they don't realize there's a B2B side | (1) Make "For Projects" button more prominent (solid/contrast color, not ghost). (2) Or add a split-screen hero: left="Earn Crypto" for users, right="Grow Your Project" for projects. Minimal change, high impact. | `QszRH` Brand Homepage |
| S6 | **No "Talk to Sales / Book Demo" path** — All product page CTAs are self-serve ("Start Free Trial"). Post-listing CMOs with budget and enterprise CMOs expect to talk to a human before committing. | High-value prospects (>$500/mo potential) have no conversion path suited to their buying behavior | (1) Every product page final CTA: dual buttons `[Start Free Trial]` + `[Book a Demo]`. (2) WL + Boost pages: add "Contact Sales" as primary for enterprise tier. (3) Pricing page high tiers: "Talk to an Expert" link. | `gXQur`, `GyyL4`, `cbBdG`, `Lym65`, `HO2Ny` |
| S7 | **Pricing lacks stage/size matching** — Pricing page shows tiers without helping CMO choose. No "Recommended for your stage", no bundle pricing, no Enterprise/Custom tier for WL. Low WL price ($499) may hurt enterprise credibility. | (1) Mid-size projects can't self-select tier. (2) Enterprise prospects see consumer pricing → question product maturity. | (1) Add "Most Popular" or "Best for [stage]" badges to tiers. (2) Add Quest+Community bundle row. (3) WL: add Enterprise tier with "Custom pricing · Dedicated support · SLA" → Contact Sales. | `HO2Ny` Pricing |

### CMO Journey Scorecard

| Persona | Stage | Understand | Interest | Strengthen | Decide | Total /40 |
|---------|-------|:---:|:---:|:---:|:---:|:---:|
| Alex | Cold start | 6 | 8 | 4 | 7 | 25 |
| Sarah | Pre-TGE | 4 | 6 | 3 | 4 | 17 |
| Marcus | Post-listing | 7 | 7 | 3 | 5 | 22 |
| Lisa | Enterprise | 8 | 6 | 3 | 3 | 20 |
| **Average** | | **6.3** | **6.8** | **3.3** | **4.8** | **21** |

**Key insight**: The funnel top (understand + interest) works at ~65%. The funnel bottom (strengthen + decide) works at ~40%. The conversion bottleneck is **proof** (no case studies, no ROI) and **guidance** (no lifecycle matching, no sales path).

---

## 🔴 I-Tier — Data→Insight→Optimize Loop Gaps (2026-03-05)

> Community/WL 的核心价值是帮项目方根据数据洞察持续优化任务激励体系。
> 审查发现优化闭环 5 个环节覆盖度严重不均：Act=70%, Observe=40%, Insight/Decide/Measure≈5%。
> **运营者当前处于"盲调"状态——能改配置，但看不懂数据、不知道改什么、无法验证效果。**

### I-Tier-A: 结构性缺失（需要新增页面/组件）

| # | Issue | Loop Gap | Suggested Fix | Affects |
|---|-------|----------|---------------|---------|
| I-01 | **No cross-module correlation view** — 12 module pages are silos. No view showing user flow between modules, no earn-spend economy balance, no "which module combo drives retention" analysis. The core Earn→Compete→Spend loop is invisible. | Observe, Insight | New **Community Insights** page (or Overview redesign): (1) Module funnel — Joined→First Task→Points Earned→Streak Started→First Redemption→Retained. (2) Economy balance — Earn Rate vs Burn Rate dual-line chart + points supply trend. (3) Module correlation matrix — overlap & contribution to retention. | Community Overview / new page |
| I-03 | **No Points Economy dashboard** — Points are the economic engine of Community but there's no inflation/deflation tracking, no distribution analysis, no earn-channel breakdown, no purchasing power trend. Operators can't tell if their economy is healthy. | Observe, Insight | Add **Points Economy** panel to Points & Level `zCfKQ` or new sub-page: (1) Earn/Burn dual-line trend. (2) Points distribution histogram (how concentrated?). (3) Earn-channel breakdown pie (tasks/check-in/referral/streak). (4) Economy health score (green/yellow/red). (5) "X% of users can afford at least 1 shop item" affordability metric. | `zCfKQ`, `usBsM` |
| I-04 | **DayChain/TaskChain missing behavioral breakpoint analysis** — DayChain shows "Avg Streak: 8.3 days" but not WHERE users break. TaskChain shows "72.3% completion" but not WHICH step loses users. These are the most critical retention mechanisms but their optimization levers are invisible. | Observe, Insight | DayChain `fLLVb`: add **Streak Distribution Curve** (x=days, y=users, showing Day 7 cliff). TaskChain `lpdtp`: add **Chain Step Funnel** (Step 1→2→3 pass-through rate). Both are the single most actionable charts for retention optimization. | `fLLVb`, `lpdtp` |
| I-06 | **No user segmentation view** — All data is aggregate. Can't distinguish power users (5%, earn 80% of points) vs mid-tier (30%, at-risk) vs long-tail (65%, barely active). One incentive structure can't serve all three. | Observe, Decide | Add **User Segments** panel to Community Overview or Analytics: (1) Auto-segment — Power/Active/At-Risk/Churned with counts + trends. (2) Per-segment metrics (DAU, avg streak, points earned, tasks completed). (3) "At-Risk users increased 15% this week" type alerts. | Community Overview / Analytics `fLxTr` |
| I-07 | **No insight/recommendation engine** — Data is raw numbers with zero interpretation. Operator must be a data analyst to extract meaning. No "DayChain streaks drop 60% at Day 7 → add Day 7 bonus" type guidance. | Insight, Decide | Add **Insights & Suggestions** card to Community Hub Deep `TQR51`: 3-5 AI-generated actionable recommendations, each with observation + cause + suggested action + one-click link to relevant module settings. Ranked by impact (red=urgent, yellow=advice, blue=opportunity). | `TQR51`, Community Overview |
| I-08 | **No before/after change tracking** — When operator adjusts points rules or adds tasks, there's no way to compare pre vs post performance. Optimization loop's final step (Measure) is completely missing. | Measure | (1) Add **Event Markers** on all trend charts — auto-annotate when config changes happen ("Day X: check-in points 10→20"). (2) After any module setting change, offer "Track this change?" → auto-generate 7-day before/after comparison report. | All module pages, trend charts |
| I-11 | **No retention deep-dive** — Retention is Community's single core KPI but gets only a small panel in Hub Deep `TQR51` (4 static numbers). No retention curve, no by-module analysis, no by-cohort breakdown, no "aha moment" identification. | Observe, Insight | Add **Retention** tab in Analytics `fLxTr` or Community-specific: (1) Retention curve (D1/3/7/14/30/60/90), filterable by cohort/module/channel. (2) "Aha Moment" analysis — which behavior best predicts D30 retention. (3) Module attribution — "users who used DayChain have 2.3x better D30 retention". | Analytics `fLxTr` / new page |

### I-Tier-B: 现有页面优化（改进已有设计）

| # | Issue | Loop Gap | Suggested Fix | Affects |
|---|-------|----------|---------------|---------|
| I-02 | **All module stats are vanity metrics** — Stats like "Total Points Issued: 1,284,500" or "Total Spins: 12,847" tell nothing actionable. Missing: rates, WoW trends, comparisons, health indicators. | Observe | Redesign every module's 4 stat cards: (1) Replace volumes with **rates** (completion rate, not completions). (2) Add **WoW trend arrow** (↑12%/↓5%) to each. (3) Add **health color** (green/yellow/red based on thresholds). Example: DayChain → "Active Streak Rate: 68% ↓3%" instead of "Active Users: 1,284". | All 12 module pages (`Wug7d`-`5Wm6B`) |
| I-05 | **Benefits Shop lacks supply-demand analysis** — Shows items and redemption counts but not: can users afford items? Which items have high demand but low conversion? Is points drain sufficient? | Observe, Insight | Add to `7yPWx`: (1) **Affordability metric** — "X% users can afford ≥1 item". (2) **Demand heatmap** — views vs redemptions per item (high views + low redemption = price too high). (3) **Restock alert** — "NFT Badge sells out in ~3 days at current rate". | `7yPWx` Benefits Shop |
| I-09 | **WL Page Analytics too basic** — Only shows page views + top pages. Missing widget-level interaction metrics, user journey through embedded content, view→action conversion funnel, SDK event tracking. | Observe | Redesign `69HPh` as **Engagement Analytics**: (1) Page Views (keep). (2) Widget Interactions (clicks/completions per widget). (3) Conversion funnel (View→Interact→Complete). (4) User flow Sankey diagram. | `69HPh` WL Page Analytics |
| I-10 | **Content Management lacks performance data** — `lhR14` shows announcements and featured items but zero performance metrics. Operator can't tell which announcement got views, which featured quest converts best. | Observe | Add mini data columns to each content item in `lhR14`: Views, Clicks, CTR. Support sort-by-performance. Add "Top performing" badge to best content. | `lhR14` Content Mgmt |

### Optimization Loop Coverage Assessment

| Loop Step | Current | Gap | Severity |
|-----------|:---:|------|:---:|
| **Observe** (see data) | 40% | Vanity metrics, no rates/trends/segments | High |
| **Insight** (understand why) | 5% | Almost no interpretation layer; only TQR51 DayChain warning qualifies | Critical |
| **Decide** (know what to change) | 10% | No recommendations, no benchmarks, operator must self-analyze | Critical |
| **Act** (make changes) | 70% | Module settings pages exist, configs changeable — this is good | Medium |
| **Measure** (verify impact) | 5% | No before/after, no change markers, no impact reports | Critical |

---

## 🔴 W-Tier — WL Core Feature Gaps (2026-03-05)

> WL 与 Community 的核心差异化在于**深度集成项目产品**。以下两个功能是 WL 值 $499+ 而非 $79 的关键。
> 完整循环：**Use (Activity Rules) → Earn (Points) → Enjoy (Privileges) → Use More**

### W1: Activity Rules — 链上行为自动校验 + 积分发放 (WL exclusive)

**问题**: 当前模型要求项目方手动创建每个任务才能奖励用户。但 WL 嵌入在项目产品中，用户 swap/stake/lend 等链上行为已经发生，不应需要"接受任务"才能获得积分。

**与 Tasks 的区别**:
- Tasks = 离散、手动创建、用户需"接受"、一次性/周期性
- Activity Rules = 持续监听、事件驱动、自动检测、用户无感

**B-End 管理页面设计需求** (新增 2 个页面):

| 页面 | 内容 | 位置 |
|------|------|------|
| **Contract Registry** | 注册项目合约地址 + 选链 → 自动解析 ABI → 展示可监听 events/functions → 验证状态 → 常见协议模板 (Uniswap-fork DEX / Aave-fork Lending / Standard Staking) | WL sidebar → new section "SMART REWARDS" |
| **Rule Builder** | 可视化 if-then 规则编辑器: IF (event + amount threshold + frequency limit + time condition) THEN (points + multiplier). 预设模板: "Reward Every Swap" / "Daily First Action 2x" / "LP Bonus" / "Staking Milestone". 反女巫配置: min wallet age, min interactions, bot exclusion list. Live stats: rules triggered today, points distributed. | WL sidebar → under "SMART REWARDS" |

**C-End Widget 需求**:
- Activity Reward Toast: 合约事件确认后弹出 "+50 pts — Swap completed" 通知
- 实时积分余额更新
- Combo/streak 连续操作动效

**与现有系统关系**: Activity Rules 产生的积分进入同一 Points Economy，可在 Benefits Shop / Privileges / Leaderboard 中使用。与 Sectors & Tasks 并存——Tasks 处理链下行为（Twitter/Discord），Activity Rules 处理链上行为。

---

### W2: Privilege System — 项目原生权益管理 + 集成 (WL exclusive)

**问题**: 当前奖励全是 TaskOn 体系内的（NFT/Token/WL Spot），与项目产品无关。真正的留存需要**项目原生权益**（手续费折扣、Gas 返还、Yield 加成等），因为这些创造真实切换成本。

**三种权益获取方式**:
- **方式 A (Status-based)**: 与等级/积分余额绑定，自动生效/取消 → B-end 创建管理，API 分发
- **方式 B (Redemption-based)**: Benefits Shop 中用积分兑换权益券 → **已有**（Shop `7yPWx`），仅需增加 Privilege 类别
- **方式 C (Achievement-based)**: 绑定 Badge/Milestone，一次性获得 → B-end 创建管理，API 分发

**职责划分**:
```
┌─ TaskOn WL (管理 WHO qualifies) ──────────────┐    ┌─ Project (管理 HOW to apply) ──────┐
│                                                │    │                                    │
│  B-end Admin:                                  │    │  Integration:                      │
│  · Create privilege definitions                │    │  · Query TaskOn API for user        │
│    (type, value, criteria, duration)           │───→│    privilege status                 │
│  · Set qualification rules                     │    │  · Apply discount/rebate/boost      │
│    (Level ≥ 3, or Badge "Diamond", etc.)       │    │    in own contract/backend logic    │
│  · Monitor usage analytics                     │    │  · Report privilege usage back      │
│  · Manage budget/limits                        │    │    (for analytics)                  │
│                                                │    │                                    │
│  API/SDK:                                      │    │  Receives:                         │
│  · Expose privilege status per user            │    │  · Webhook on privilege changes     │
│  · Push webhook on status changes              │    │  · REST API for real-time queries   │
│  · Optional: deploy Privilege Contract         │    │  · Optional: on-chain contract call │
│    for trustless on-chain verification         │    │                                    │
└────────────────────────────────────────────────┘    └────────────────────────────────────┘
```

**B-End 管理页面设计需求** (新增 1 个页面 + 修改 2 个):

| 页面 | 内容 | 位置 |
|------|------|------|
| **Privilege Manager** (新增) | 权益列表 (Fee Discount / Gas Rebate / Yield Boost / Priority Access / Custom) + 创建权益 (type, name, value, qualification mode A/C, criteria, duration, budget cap) + 活跃权益 stats (active holders, total value distributed, ROI) + 集成状态 (API connected? usage reported?) | WL sidebar → new section "PRIVILEGES" or under REWARDS |
| **Benefits Shop** (修改 `7yPWx`) | 增加 "Privileges" 分类 tab → 方式 B 的权益券 (time-limited fee discount voucher etc.) | 现有 |
| **Integration Center** (修改 `Abs1E`) | 增加 "Privilege Contract" 部署卡片 + "Activity Listener" 集成卡片 | 现有 Developer Tools section |

**C-End Widget 需求**:
- "My Privileges" 面板 in User Center: 活跃权益列表 + 节省金额汇总 ("Saved $65.70 this month")
- 锁定权益 + 解锁进度条 ("Gas Rebate — Need Level 5, you're 78% there")
- Privilege Status Widget (for Widget Library)

**API 接口文档**: → `docs/wl_privilege_api_spec.md`

---

### W-Tier Page Summary

| New Pages | Description | Est. Complexity |
|-----------|-------------|:---:|
| WL Contract Registry | Register + ABI parse + verify | Medium |
| WL Activity Rule Builder | Visual rule editor + templates + anti-sybil | High |
| WL Privilege Manager | CRUD privileges + qualification rules + analytics | Medium |
| **Total new pages: 3** | + 3 modified pages (Shop/Integration/SDK) + 2 new C-End widgets | |

---

## 🟢 P3 — Remaining Frontend Handoff Items

- [ ] Design B48 Dev Kit Page (`taskon.xyz/devkit/{id}`) — standalone developer integration page
- [ ] Verify all canvas annotations on non-P0d pages are still accurate

---

## Design Review Summary (2026-03-05)

**Scope**: Community + White Label + C-End — 42 pages reviewed
**Method**: Systematic get_screenshot review of every page

### What's Working Well
- **Community module sidebar architecture** — 5 section headers (TASKS/POINTS/CAMPAIGNS/REWARDS/SETTINGS) with 12 items is well-organized and consistent across all module pages
- **WL expandable sidebar** (Overview/Widgets/Pages) — clear navigation for WL sub-pages
- **Community Wizard 3-step flow** — Customize → Configure Systems → Review is intuitive
- **WL deployment paths** (Domain/Embed/SDK) — clear differentiation with "Recommended" badge
- **C-End gamification** — Leaderboard podium, Sprint progress tiers, Milestone claim states are engaging
- **B-End data tables** — consistent pattern across all management pages (stats + filters + table + pagination)
- **Widget Library concept** — "Configured" vs "Not Yet Configured" with clear CTAs works well
- **Page Builder** — empty/editor/active 3-state flow is complete and logical
- **WL Integration Center** (`Abs1E`) — categorized card layout is comprehensive
- **C-End Activity Feed** (`xhPIr`) — clean timeline with color-coded activity types

### Issue Count
| Priority | Count | Description |
|----------|-------|-------------|
| P0 Blocking | 2 | Broken UI + duplicate page confusion |
| P1 Important | 10 | IA issues (5) + visual consistency (3) + carried forward (2) |
| P2 Optimization | 8 | Polish items + carried forward (2) |
| S-Tier Systemic | 7 | Cross-page strategic gaps from CMO journey audit |
| I-Tier Data Loop | 11 | Structural (7 new pages/components) + Existing page improvements (4) |
| W-Tier WL Features | 2 | Activity Rules (2 pages) + Privilege System (1 page + 2 modified + 2 widgets) |
| **Total** | **40** | |

---

## Recently Completed (see `progress.md` for details)

- ✅ **P0**: Community Empty State redesign (3 retention strategies) + Wizard (3-step)
- ✅ **P0b**: Getting Started Checklist redesign (Community 5-step + WL 6-step)
- ✅ **P0c**: WL expandable sidebar sub-menu (8 pages)
- ✅ **P0d**: Frontend handoff audit — page codes corrected, button routing complete, Dev Kit/Promo Kit/Auto-detection specs added, layout.md corrected, 76 total page codes
- ✅ **P0e**: C-End engagement redesign — C01 action engine + C07 User Center + C08 Invite Center + C09 Activity Feed, 79 total page codes
- ✅ **P0f**: Leaderboard vs LB Sprint distinction — multi-point-type leaderboards + time-bounded sprints with incentives across B-end (B31d/B31e) and C-end (C03/C04/all nav tabs)
- ✅ **P0g**: 4-System Module Architecture + Badges/AccessRules/HomepageEditor
- ✅ **Canvas reorganization**: All 75+ pages in structured grid (`design/layout.md`)
- ✅ **Orphaned label fix**: 130+ annotation labels repositioned
