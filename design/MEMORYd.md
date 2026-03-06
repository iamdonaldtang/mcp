# TaskOn Project Memory

## Project Identity
- **What**: TaskOn is a Web3 Growth Engine — SaaS platform for Web3 projects to acquire, retain, and grow users
- **Products**: Quest (acquisition), Community (retention), White Label (full-stack), Boost (CPA/CPS pay-per-result)
- **Brand Colors**: Quest=#5D7EF1, Community=#48BB78, White-Label=#9B7EE0, Boost=#ED8936
- **Design File**: `design/pencil-new.pen` — ONLY accessible via Pencil MCP tools (never Read/Grep)

## Current Project: Complete Website + B-End Product Design
- **Status**: ALL DESIGN WORK COMPLETE + DARK THEME APPLIED — 93 page codes, 81 pages, 19 modals. T01-T05/T25 skipped.
- **Dark Theme Restyle**: DONE — all pages converted from light (#F8FAFC) to dark (#0A0F1A) theme. Color mapping in `design/todo.md`.
- **Progress**: `design/progress.md` | `design/progress-archive.md` (full history)
- **Frontend Requirements**: `docs/website_frontend_requirements.md` v4.2 — 90+ page codes, complete button routing, C-End APIs
- **B2C Operation Mapping**: `docs/b2c_operation_mapping.md` — 153 B-end to C-end operations mapped
- **Pending design**: Publish Readiness Check modal (D20), Community Integration page refinement
- **Next step**: frontend development (after pending items)

## Page Node IDs
See `design/pages.md` for complete page list with Node IDs. Key frequently-referenced pages:
- Community: Empty `zzZ8D`, Guided `S1EIA`, Active `vFRHi`, Wizard `Gzpeu`/`8NeyG`/`qknQZ`
- C-End: Home `vJVhd`, User Center `PykHF`, Invite Center `TaAo9`, Activity Feed `xhPIr`
- WL: Active `BnkYW`, Mgmt `UPAfV`, Widget `2sSsA`/`n4pJK`/`S432k`, PageBuilder `DRYwN`/`sGDcq`/`J08v5`, Contract Registry `OKEqS`, Rule Builder `4aAo7`, Privilege Manager `5xwYN`
- Community Insights: `olPfE` (cross-module analytics, economy health, user segments, retention by module)
- Community Integration Center: `ZL5K5` (B61)
- WL Embed Options neutral: `Rwq2K`, PB has-pages: `zW40A`
- Quest Mgmt template: `XvXEQ` (used as copy source for module mgmt pages)

## Key Decisions
- C-end and B-end are **completely independent systems**
- B-end uses **Dark Theme** (#0A0F1A page bg, #111B27 cards) — restyled from original light theme
- Boost is **independent 4th product** (not merged into Quest)
- Design includes **complete creation wizard flows** (not just hub pages)
- Brand Homepage header: ultra-minimal (Logo + [For Projects] + [Launch App])
- Footer: 3 columns (For Users / For Projects / Company)
- **Pricing**: Unified Pricing page M07 (`HO2Ny`) with 3 product tabs (Quest/Community/WL). Boost = CPA (no pricing page). Base monthly prices: Quest=$300, Community=$600, WL=$1,500.
- **Billing cycles (per product, CRITICAL)**: Quest & Community = Quarterly/Semi-Annual/Annual (NO monthly). White Label = Monthly/Quarterly/Semi-Annual (NO annual). Design shows Quest tab + Quarterly selected as default state.
- **WL = Quest + Community + exclusive features**: Page Builder, embed widgets, custom domain, deep customization, user segmentation, retention analytics. Low-code SDK, first integration < 7 days.
- **WL 3 Deployment Paths** (Hub goal cards, redefined): Host on Your Domain (zero code) / Embed in Your App (★ RECOMMENDED, Widget+PB+SSO) / Build with SDK (full custom). Domain Setup: `5bmH9`. Embed Options: `RgCVQ` (3 modes: Iframe/Widget Library/Page Builder). Widget Library: `2sSsA`/`n4pJK`/`S432k` (empty/config/active). Page Builder: `DRYwN`/`sGDcq`/`J08v5` (empty/editor/active).
- **Widget Library concept**: Shows "Community Modules" (not "Available Widgets"). Two types: Configured (green, "Add Widget →") and Not Yet Configured (amber, "Set Up in Community →"). Points Balance renamed to "User Center".
- **WL Wizard Embed flow**: Step 1 Path → Step 1.5 Embed Options (neutral `Rwq2K`) → selected method config. 5 path blocks on canvas (Domain/Iframe/Widget/PB/SDK).
- **Smart Rewards sidebar**: Rule Builder (`4aAo7`) + Privilege Manager (`5xwYN`) consolidated under "Smart Rewards" sub-item in WL sidebar. Added to all 8 WL Hub pages + both Smart Rewards pages.
- **PB Config has-pages state**: `zW40A` — shows existing pages to select when user already has pages built.
- **Community Integration Center**: `ZL5K5` (B61) — independent Integration page for Community product (Social+Blockchain+Analytics, no Developer Tools). Entry from Community Hub Deep `TQR51` INTEGRATIONS section.
- **Publish Readiness Check (D20)**: Pre-publish modal with 2-item checklist (subscription status + Twitter auth). Applies to all Publish/Launch/Activate buttons across Community (10 pages), WL (5 pages), Quest (2 pages), Boost (2 pages). Full touchpoint list in `design/todo.md`.

## Activation & Onboarding Philosophy (Critical)

### Core Insight: Setup ≠ Activation
The "death zone" is between setup complete and first real user engagement. Checklists must bridge this gap — they are **activation funnels**, not feature configuration lists.

### Checklist Design Principle: Configure → Verify → Distribute → Celebrate
- **Configure**: TaskOn helps (wizard pre-fills) — LOW friction
- **Verify**: System auto-detects (integration ping, DNS) — ZERO user effort
- **Distribute**: TaskOn generates promo materials — MINIMAL friction
- **Celebrate**: Real-time counters, milestones — INSTANT feedback

### Community Activation Bottleneck = Distribution
Project owner's users are on Twitter/Discord/Telegram. If they don't announce, community stays empty → churn.
**Solution**: Promo Kit Generator (AI-generated social posts + branded banner) + one-click share buttons.

### White Label Activation Bottleneck = Integration + Distribution
Two bottlenecks: (1) code must actually be deployed (technical), (2) users must know about it.
**Solutions**:
- **Dev Kit Page**: Auto-generated standalone page (`taskon.xyz/devkit/abc123`) with everything a developer needs — no TaskOn account required. Includes: integration code (copy-paste ready, project ID pre-filled), SSO setup (2 options: wallet/OAuth), step-by-step guide (est. 30 min), self-service verification tool. Marketing person sends ONE link to dev, zero back-and-forth.
- **Integration Verification**: Auto-detect first API ping from project's domain. Real-time status updates.
- **Promo Kit**: Same as Community — AI-generated social posts + branded banner.

### Getting Started Checklists (OB Restructuring — replaces previous design)
**Community** (`S1EIA`): Module-level content enrichment. ✅ auto (created+tasks+points) → Add more tasks → Set up Shop → Customize module rewards → Preview as user → Share (Promo Kit) → First 10 participants
**White Label** (`BnkYW`): Path-adaptive tool configuration. ✅ auto (community+path+brand) → Configure widgets → Build page → Preview → Send Dev Kit → Integration verified → Announce → First interaction

### Community Wizard (OB: 4 steps, was 3)
Step 1 Customize (`Gzpeu`): name, description, brand color (minor stepper update)
Step 2 Modules (`8NeyG`): strategy-driven pre-selection, module C-end previews + effect descriptions
Step 3 Quick Setup (`qknQZ` repurposed): **template content auto-generated** per enabled module, inline editable
Step 4 Preview & Publish (**NEW page**): C-end preview (desktop/mobile) + readiness checklist + Publish button

### WL Wizard (OB: 4 steps, restructured)
Step 1 Path (`NNwid`): 3 paths aligned with Empty (was 4 modes). Embed/Domain/SDK.
Step 2 Configure (`CXzmy` repurposed): **path-adaptive** — widgets (Embed) / DNS (Domain) / SDK keys (SDK)
Step 3 Brand (`5nCtO`): keep + improved preview
Step 4 Preview & Publish (**NEW page**): deployment preview + readiness checklist + Publish button

### Strategy → Module Mapping (OB-C1)
- Activate New Users: Sectors & Tasks + Points & Level + TaskChain
- Drive Daily Engagement: Sectors & Tasks + Points & Level + DayChain + Leaderboard
- Maximize Retention: Sectors & Tasks + Points & Level + DayChain + Milestones + Benefits Shop

### WL Expandable Sidebar Sub-menu
8+ WL pages have expandable sidebar: WL parent item with `keyboard_arrow_up` chevron + 4 sub-items (Overview/Widgets/Pages/Smart Rewards). Active item: purple bg #1A1033 + text #9B7EE0 + fontWeight 600. Inactive: no fill + text #94A3B8. Sub-item padding: [8,12,8,40], fontSize 13, icons 16x16.
- Overview active: `BnkYW`, `UPAfV`
- Widgets active: `2sSsA`, `n4pJK`, `S432k`
- Pages active: `DRYwN`, `sGDcq`, `J08v5`
- Smart Rewards active: `4aAo7` (Rule Builder), `5xwYN` (Privilege Manager)

### Community Retention Strategy Framework
Empty State (`zzZ8D`) shows 3 lifecycle-stage strategies (not module-level templates):
1. **Activate New Users** (green): TaskChain + Points + Levels — Goal gradient + instant feedback — Activation Rate
2. **Drive Daily Engagement** (orange): DayChain + Leaderboard + Sprint — Loss aversion + social comparison — DAU
3. **Maximize Retention** (purple): Benefits Shop + Milestones + Lucky Wheel — Sunk cost + variable reinforcement — 30-Day Retention
Economy loop: Earn → Compete → Spend. Sectors & Tasks is infrastructure (auto-configured, not in templates).

### Community 4-System Module Architecture (P0g)
9 flat modules → 4 logical systems + 2 infrastructure settings:
1. **Task Engine**: Sectors & Tasks (core★), TaskChain, DayChain
2. **Points & Recognition**: Points & Level (core★), Leaderboard, Badges (NEW)
3. **Incentive Campaigns**: LB Sprint, Milestone, Lucky Wheel
4. **Rewards Economy**: Benefits Shop
**Infrastructure** (under SETTINGS in sidebar, NOT in wizard): Access Rules, Homepage Editor

Sidebar has 5 section headers (TASKS/POINTS/CAMPAIGNS/REWARDS/SETTINGS) — non-clickable, gray uppercase, fontSize 10, fontWeight 600, fill #94A3B8, letterSpacing 1.

New pages: Badges `BJLsz`, Access Rules `g1CNC`, Homepage Editor `5Wm6B`

## Leaderboard vs LB Sprint (Critical Distinction)
- **Leaderboard** = recurring (weekly/monthly/all-time) ranking display based on custom point types (EXP, GEM, etc.). **No extra incentives** — purely points-based competition.
- **LB Sprint (Leaderboard Sprint)** = time-bounded competition with start/end dates, based on custom point types, **with non-points incentives** (NFT/Token/WL Spot).
- Community supports multiple custom point types — each can have its own leaderboard and LB Sprint.
- B-end: B31d=Leaderboard (multi-point-type), B31e=LB Sprint (point type badges + incentives). All sidebar sub-menus say "LB Sprint" not "Sprint".
- C-end: C03=Leaderboard (point type selector EXP/GEM), C04=LB Sprint (sprints with NFT/Token incentives). All nav tabs say "LB Sprint" not "Sprint".

## Technical Patterns (Pencil MCP)
- Icon system: Material Symbols Rounded (NOT Lucide — many Lucide icons fail)
- **`expand_less` icon doesn't exist** — use `keyboard_arrow_up` instead
- x/y positioning requires `layout: "none"` on parent (flexbox ignores x/y)
- **Text nodes do NOT support width property** — neither numeric nor fill_container. Width updates are silently dropped. Use manual `\n` breaks or fixed-width parent frames.
- **flexWrap does NOT work** — must use explicit row frames for grid layouts
- **M() move requires real node ID** — cannot use binding names from I() as the moved node ID (bindings only work as parent references)
- **R() for replacing node types** — use R() to swap icon_font → text or vice versa inside components (e.g., replacing check icon with number text in stepper dots)
- Radial gradients: use subtle values to avoid white wash
- Operations: I()=Insert, U()=Update, D()=Delete, M()=Move, C()=Copy, R()=Replace, G()=Image
- Max 25 operations per batch_design call
- Always verify with get_screenshot after major changes
- Always check snapshot_layout(problemsOnly:true) for issues
- **Sidebar sub-menu insertion pattern**: After I() into sidebar, use M(realNodeId, sidebarId, index) to position correctly (typically index 7 = after WL/Community item)

## Three-Layer Page Hierarchy (WHY / WHAT / HOW)
Each product has three page layers — do NOT mix content between them:
- **Product Page** (before login) = WHY — strategic value narrative, convince evaluator to sign up
- **Empty State** (after login) = WHAT — tactical capability framework, build confidence to start creating
- **Wizard** (creation flow) = HOW — operational step-by-step config, reduce friction to activate

**Cross-layer bridges:**
- Empty State → Product Page: "Learn More" link in resources (new tab, NOT main nav)
- Empty State → Wizard: "Create Your First X" CTA (primary path)

**Empty State design rules:**
1. Show framework, not feature marketing (templates = "what you can build")
2. Highlight strip for differentiation (compact one-liners, NOT full sections)
3. Single primary CTA emphasizing ease ("one click")
4. Resources as depth exit for users who want more

## Product Value Propositions
| Product | Core Value | Psychology Lever |
|---------|-----------|-----------------|
| Quest | "Get users fast" — one-time acquisition | Instant gratification |
| Community | "Make users stay" — create leaving costs | Sunk cost, loss aversion, goal gradient, endowment |
| White Label | "Own the experience" — branded full-stack | Control/ownership |
| Boost | "Guaranteed results" — pay per outcome | Risk aversion |

## Design Patterns
- **Marketing pages**: 1440px wide, 120px padding, Inter font, Hero 48-52px, Section 36px, Body 16px
- **B-End pages**: Sidebar 240px + fluid content, Inter font, Page title 24px, Section 18px, Body 14px
- **Hub Empty State**: Welcome → Template/goal cards → Highlight strip → Primary CTA → Divider → Resource cards
- **Hub Active State (campaign-based)**: Header + stats row + campaign cards (badge + metrics + progress bar) + Quick Actions
- **Hub Active State (platform-based)**: Header + stats row + onboarding checklist/toolkit cards + resources
- **Hub Management (campaign-based)**: Header + power stats + filter tabs + search + data table + pagination
- **Hub Management (platform-based)**: Header + power stats + all-configured cards + deployments/analytics sections
- **Module Management (Community)**: Expandable sidebar sub-menu with 5 section headers (TASKS/POINTS/CAMPAIGNS/REWARDS/SETTINGS) + 12 items (Overview + 9 modules + Access Rules + Homepage Editor) + header/stats/filters/data table/pagination. Each module has its own management page. Sidebar active item: green bg #ECFDF5 + text #48BB78. Sub-items: 13px Inter, padding [8,12,8,40], icons 16x16.
- **Dashboard progression**: New User (goal cards + quick start) → Active (stats + campaign cards + quick actions + upsell) → Power (stats + growth chart + product breakdown + activity feed)
- **Wizard template**: Top bar (back+title+save) → Stepper → Step content (often 2-column) → Action bar (back+next)
- **Cross-sell pattern**: Each product page's Final CTA links to next product in hierarchy
- **C-End pages**: Dark header (#0F172A) + horizontal nav tabs (amber underline) + user status bar + dark content (#0A0F1A). Amber (#F59E0B) primary accent. Gamified task cards (icon box + points badge + action button). "Powered by TaskOn" footer.
- **C-End gamification states**: Sprint (completed/in-progress/not-started), Milestones (claimed/claimable/locked), Shop (affordable/not-enough/sold-out)
- **Status badges**: Active=#0A2E1A/#16A34A, Draft=#1F1A08/#D97706, Completed=#1E293B/#64748B, Paused=#2D1515/#DC2626

## QA Patterns
- **Nav completeness check**: Audit all nav targets to verify pages exist
- **Button routing audit**: Use `website_frontend_requirements.md` §三 as source of truth

## User Preferences
- Communicates in Chinese, design content in English
- Thinks from first principles, values critical thinking over feature listing
- Wants alignment before execution ("确定之后再开始设计细节")
- Prefers structured progress tracking

## Canvas Layout (updated session 5)
- Full grid: `design/layout.md` (30 rows including modals + Dev Kit)
- Modals D01-D19: Rows 23-29 (y:53250-60200)
- Dev Kit B48: Row 30 (y:62000)
- Pricing tab states: Row 2 cols 4-5 (next to original M07)

## Key Docs
- `docs/website_frontend_requirements.md` — v4.2: 90+ page codes, button routing, C-End APIs, §六 新功能组件规格
- `docs/TaskOn第一性原理与人性洞察准则.md` — Core design principles
- `docs/taskon_growth_engine_playbook_en.md` — Product strategy, AARRR
- `docs/b2c_operation_mapping.md` — 153 B→C operation mappings
- `docs/white_label_integration_modes.md` — WL 4 integration modes
- `design/todo.md` — Completed design tasks log + OB reference spec
- `design/layout.md` — Canvas grid map (all page + modal positions)
- `design/pages.md` — Complete page + modal node ID reference
