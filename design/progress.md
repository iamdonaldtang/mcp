# TaskOn Website & B-End Product Design — Progress Tracker

> Last updated: 2026-03-05
> Design file: `design/pencil-new.pen` (Pencil MCP tools only, never Read/Grep)
> **Active TODO**: `design/todo.md` — remaining design tasks
>
> ### Completion Summary
> - **P1-P14**: 66+9 pages designed (M01-M14 + B01-B48 + B31a-B31h + C01-C06 = 76 page codes)
> - **Design Review**: 90 issues → Fixed 68 / Deferred 19 / Remaining 5 (see `todo.md`)
> - **P0 Community Redesign**: Empty State (3 retention strategies) + Wizard (3-step) + Checklist (activation funnel)
> - **P0b Getting Started Checklists**: Community (5-step) + WL (6-step) — activation-driven, not feature-driven
> - **P0c WL Sidebar Sub-menu**: 8 WL pages with expandable Overview/Widgets/Pages sub-items
> - **P0d Frontend Handoff Audit**: Page codes corrected, button routing complete, Dev Kit/Promo Kit/Auto-detection specs added
> - **Canvas Layout**: Organized grid, documented in `design/layout.md`
> - **Frontend Requirements**: `docs/website_frontend_requirements.md` v4.1 — 76 page codes, complete button routing

---

## Project Overview

TaskOn is a Web3 Growth Engine platform. We are designing the complete website (marketing pages) and B-end product interface (after-login dashboard & product hubs).

**Two independent systems:**
- **C-end App**: Users complete tasks, earn crypto (accessed via "Launch App")
- **B-end Dashboard**: Projects manage campaigns, grow users (accessed via "Login")
- These are completely separate systems with independent URLs, navigation, and user models.

---

## Key Strategic Decisions (Confirmed)

| Decision | Result |
|----------|--------|
| Homepage focus | C-end user acquisition, B-end routing via "For Projects" |
| Products | 4 independent products: Quest, Community, White Label, Boost |
| Boost positioning | Independent 4th product — CPA/CPS pay-per-result model |
| B-end dashboard theme | Light Theme (professional SaaS style) |
| Design depth | Complete flows including multi-step creation wizards |
| C-end / B-end boundary | Completely independent systems |
| Header (Brand Homepage) | Ultra-minimal: Logo + [For Projects] + [Launch App] |
| Footer structure | 3 columns: For Users / For Projects / Company |
| Pricing architecture | Quest+Community share Platform Pricing (tab switch); WL has separate White Label Pricing page; Boost = CPA (no pricing page) |
| Billing cycles | Monthly + Annual (save 20%) |
| Pricing tiers | 2 per product: Free / Pro ($79/mo annual) for Quest & Community; Standard ($499/mo) / Pro ($999/mo) for WL |
| WL positioning | WL = Quest + Community + exclusive features (Page Builder, widgets, custom domain, deep customization, user segmentation, retention analytics). Low-code SDK, first integration < 7 days |

## Product Hierarchy

```
Boost (零门槛)  →  Quest (自运营)  →  Community (深度留存)  →  White Label (全栈拥有)
"Give me results"  "I'll run campaigns"  "I'll build community"  "I own everything"
Effort: Zero        Low                   Medium                  High
Control: Low        Medium                High                    Full
```

**Brand Colors:** Quest=#5D7EF1(blue), Community=#48BB78(green), White-Label=#9B7EE0(purple), Boost=#ED8936(orange)

---

## Design Scope — 38 Pages Total

### Before Login (Marketing Site) — 8 pages

| # | Page | ID in .pen | Status | Notes |
|---|------|-----------|--------|-------|
| 1 | Brand Homepage | `QszRH` | DONE | Dark theme, C-end focused. Hero floating cards, leaderboard, value props, B-end routing section, social proof + supported chains, final CTA. Header: Logo + [For Projects] + [Launch App]. Footer: For Users / For Projects / Company. |
| 2 | Projects Landing Page | `Lz2vL` | DONE ✅ | B-end front door. Light theme. Hero: "Stop Renting Traffic" + trust badges. Sections: 4 Growth Engines (Boost/Quest/Community/WL cards), Build vs Buy comparison (red/green), Go Live in 3 Steps, Social Proof (22M+/2000+/$50M+/95%), Final CTA. New section added: `jcmer` (How It Works). |
| 3 | Quest Product Page | `gXQur` | DONE ✅ | Self-serve acquisition engine. Blue accent. Hero: "10,000 Real Users. 10 Minutes. Zero Code." Features: Launch in 10 min, 95% Real Participation, Full Attribution, 500+ Templates. Final CTA cross-sells to Community. |
| 4 | Community Product Page | `GyyL4` | DONE ✅ | Deep retention engine. Green accent. Hero: "80% Retention at 90 Days. One-Click Community System." Copy updated to outcome-focused. Cross-sells to White Label. |
| 5 | White Label Product Page | `cbBdG` | DONE ✅ | Full-stack growth engine. Purple accent. Hero: "Your Domain. Your Brand. Your Growth Engine." WL = Quest + Community + exclusive (Page Builder, embed widgets, deep customization, user segmentation, retention analytics, custom domain). "Choose Your Integration" section: 3 mode cards (Domain Mode / Widget + Page Builder ★ RECOMMENDED / Full SDK) + Embed Mode note + Widget+PB workflow strip (① Community → ② Widgets → ③ Pages → ④ Embed). Cross-sells from Quest/Community. |
| 6 | Boost Product Page | `Lym65` | DONE ✅ | Pay-per-result engine. Orange accent. Hero: "Guaranteed Results. Pay Only for Real Users." How It Works (3 steps), Key Features (CPA/Anti-Sybil/Tracking/Network), Campaign Types dark section (4 types with pricing), Final CTA cross-sells to Quest. |
| 7 | Platform Pricing | `CXtOH` | DONE ✅ | Quest+Community pricing. Light theme. Billing toggle (Monthly/Annual, save 20%). Product tabs (Quest/Community). 2 tiers: Free ($0) / Pro ($79/mo annual). Feature lists. WL upsell banner. 4 FAQ items. |
| 8 | White Label Pricing | `EDoSn` | DONE ✅ | WL-specific pricing. Purple accent. Hero: "White Label Pricing". Back link to Platform Pricing. 2 tiers: Standard ($499/mo) / Pro ($999/mo, RECOMMENDED). All Quest+Community features included. Low-code SDK emphasis. CTA section purple. |

### After Login (B-End Product) — 30 pages + 5 C-End pages

| # | Page | ID | Status | Notes |
|---|------|----|--------|-------|
| 9 | B-End Dashboard (New User) | `4SMOO` | DONE ✅ | Welcome + 4 goal cards + Quest recommended quick start + 3 resource cards |
| 10 | B-End Dashboard (Active User) | `IDezm` | DONE ✅ | Stats (3 campaigns / 2,847 participants / 68.3% / $3,200) + 2 campaign cards + Quick Actions + WL upsell |
| 11 | B-End Dashboard (Power User) | `W93vp` | DONE ✅ | Power stats (12 campaigns / 47,280 participants) + 12-month growth chart + Product Breakdown + Recent Activity feed |
| 12 | Quest Hub (Empty) | `VkSu7` | DONE ✅ | 3 template cards + CTA + 2 resource cards |
| 13 | Quest Hub (Active) | `4Mb1C` | DONE ✅ | Stats + 2 campaign cards (Active + Completed) with progress bars |
| 14 | Quest Hub (Management) | `XvXEQ` | DONE ✅ | Power stats + filter tabs + data table 7×6 + pagination |
| 15 | Quest Wizard (Step 1-2) | `dn8Eu`, `hBgxY` | DONE ✅ | Step 1: Template grid. Step 2: Task config (master-detail) |
| 16 | Community Hub (Empty) | `zzZ8D` | DONE ✅ | 3 template cards + highlight strip + one-click CTA + 3 resources |
| 17 | Community Hub (Guided Workspace) | `S1EIA` | DONE ✅ | Onboarding checklist + active modules + add-more modules + resources |
| 18 | Community Hub (Active) | `vFRHi` | DONE ✅ | Quick stats + 4/5 checklist + module cards with real metrics |
| 19 | Community Hub (Deep) | `TQR51` | DONE ✅ | Power stats + 6 module cards (2 rows) + Engagement Overview (chart + retention) |
| 19b | Community: Points & Level | `zCfKQ` | DONE ✅ | Module management: stats + filter + 3-row table (Bronze/Silver/Gold) + expandable sidebar sub-menu |
| 19c | Community: TaskChain | `lpdtp` | DONE ✅ | Module management: chain list with completions/rate + sidebar sub-menu |
| 19d | Community: DayChain | `fLLVb` | DONE ✅ | Module management: streak chains with active users + sidebar sub-menu |
| 19e | Community: Leaderboard | `Emmab` | DONE ✅ | Module management: leaderboard configs with metrics/period + sidebar sub-menu |
| 19f | Community: Sprint | `FO9JR` | DONE ✅ | Module management: sprint campaigns with duration/participants + sidebar sub-menu |
| 19g | Community: Milestone | `WFdZQ` | DONE ✅ | Module management: milestone sets with completions/claim rate + sidebar sub-menu |
| 19h | Community: Benefits Shop | `7yPWx` | DONE ✅ | Module management: items with cost/stock/redemptions + sidebar sub-menu |
| 19i | Community: Lucky Wheel | `sme5a` | DONE ✅ | Module management: wheel configs with spins/prizes/winners + sidebar sub-menu |
| 20 | Community Wizard | `Gzpeu` | DONE ✅ | Step 2: Customize (name, color, live preview) |
| 21 | White Label Hub (Empty/Toolbox) | `Ir6Tq` | DONE ✅ | Goal-oriented: 3 goal cards + 6 toolkit cards + resources |
| 22 | White Label Hub (Active) | `BnkYW` | DONE ✅ | Stats + 6 toolkit cards with configured status + resources |
| 23 | White Label Hub (Management) | `UPAfV` | DONE ✅ | Power stats + all tools configured + Active Deployments + Usage Analytics |
| 24 | White Label Wizard | `CXzmy` | DONE ✅ | Step 2: Configure (domain, SSL, embed code preview) |
| 24b | WL Domain Setup | `5bmH9` | DONE ✅ | "Host on Your Domain" flow: 3 steps (Enter Domain + Configure DNS + Brand Portal) + Setup Status checklist + Portal Preview mockup + Help card |
| 24c | WL Embed Options | `RgCVQ` | DONE ✅ | "Embed in Your App" sub-page: 3 embed mode cards (Iframe Embed / Widget Library ★ RECOMMENDED / Page Builder) with code snippets + Quick Comparison table (5×4) + tip card |
| 24d | WL Widget Library (Empty) | `2sSsA` | DONE ✅ | "Community Modules" (7): 6 configured widget types (Leaderboard/Task List/User Center/Rewards Shop/Daily Check-in/Quest Card) in 2×3 grid + "Not yet configured" section with Milestones (amber, "Set Up in Community →") + tip |
| 24e | WL Widget Config | `n4pJK` | DONE ✅ | Leaderboard config: 2-column (Widget Settings form + Live Preview + Embed Code with Copy button + help note) |
| 24f | WL Widget Library (Active) | `S432k` | DONE ✅ | My Widgets (2) with active Leaderboard + Task List cards (stats + embed code) + Community Modules (5): 4 configured + Milestones (not configured, amber) |
| 24g | WL Page Builder (Empty) | `DRYwN` | DONE ✅ | "How It Works" 3-step process (Create Widgets → Design Page → Publish & Embed) + 3 Page Templates (Rewards Hub/Community Portal/Custom Page) + widget requirement callout + CTA + tip |
| 24h | WL Page Builder (Editor) | `sGDcq` | DONE ✅ | 2-column: Left = canvas preview (browser URL bar, brand header, Leaderboard + Task List widget blocks, "+ Add Widget Block") / Right = Page Settings (name, slug, theme toggle) + Widgets on Page (2) + Available Widgets (3, with "+ Add") + Publish/Save Draft buttons |
| 24i | WL Page Builder (Active) | `J08v5` | DONE ✅ | "Create New Page" header button + My Pages (1) with Rewards Hub card (Published badge, stats 1,247 views/342 clicks/27.4% CTR, embed code, Edit/Analytics) + Quick Start templates (4 options) |
| 24j | WL Integration Center | `Abs1E` | DONE ✅ | Stats bar (2/12 active) + 4 categories × 3 cards: Social (Twitter✓/Discord✓/Telegram) + Blockchain (Multi-Chain/Wallet Connect/On-Chain) + Analytics (GA/Webhooks/Data Export) + Developer (API Keys/SDK/SSO). Connected = green border + badge; Available = gray + purple Connect |
| 25 | Boost Hub (Empty) | `stYvi` | DONE ✅ | 3 goal cards + CTA + 2 resource cards |
| 26 | Boost Hub (Active) | `5C3WP` | DONE ✅ | Stats + 2 campaign cards (Active + Under Review with info banner) |
| 27 | Boost Hub (Management) | `8gT3V` | DONE ✅ | Power stats + filter tabs + table 7×5 (CPA color-coded) + pagination |
| 28 | Boost Wizard (Step 1) | `SDfui` | DONE ✅ | Set Goal: campaign goal selection + targeting |
| 29 | Boost Wizard (Step 2-4) | `l9tmF`, `KMtqR`, `fZpcQ` | DONE ✅ | Step 2: Channels. Step 3: Budget + CPA calculator. Step 4: Review & Submit |
| 30 | B-End Sectors & Tasks | `Wug7d` | DONE ✅ | Community Sectors & Tasks management with drag handles, type badges, status pills, hidden sector |
| 31 | B-End Content Management | `lhR14` | DONE ✅ | Announcements, Featured Slots 2×3, Module Status Overview 2×3 |
| 32 | B-End Preview Mode | `2UiNC` | DONE ✅ | Preview banner + Desktop/Mobile toggle + embedded C-end view |
| C1 | C-End Community Home | `vJVhd` | DONE ✅ | Dark header, nav tabs, status bar, announcements, featured, DayChain, task sectors |
| C2 | C-End Quest Tab | `dUXTl` | DONE ✅ | Quest card grid 2×2 with filters, 4 quest states |
| C3 | C-End Leaderboard | `KmdSd` | DONE ✅ | Podium top 3, "Your Rank" card, rankings table 5 rows |
| C4 | C-End Sprint Tab | `y5fUZ` | DONE ✅ | Current sprint card, 5 tasks (3 states), reward tiers, past sprints |
| C5 | C-End Milestone Tab | `53iKE` | DONE ✅ | Level progress, 6 milestones (claimed/claimable/locked) |
| C6 | C-End Shop Tab | `coM7o` | DONE ✅ | 6 shop items (redeem/not enough/sold out), category filters |

---

## Remaining Design Scope (P6-P9)

### P6: Wizard Completion + Hub Strategic Redesign — DONE ✅

**Strategic decisions made during P6:**
- **Quest Wizard**: Steps 1-2 already designed, Steps 3-4 skipped (uses existing system patterns)
- **Community**: One-click initialize (no template selection wizard). Community Hub IS the guide → created Guided Workspace state
- **White Label Hub**: Redesigned from mode-based lock-in to goal-oriented toolbox (like Stripe)
- **Boost Wizard**: Completed all remaining steps (1, 3, 4)

**Boost Wizard** (4-step stepper):
| Step | Content | ID | Status |
|------|---------|-----|--------|
| Step 1: Set Goal | Campaign goal selection (User Acquisition / Trading Volume / TVL) + targeting | `SDfui` | DONE ✅ |
| Step 2: Select Channels | Channel options (Explore / Push / Email), campaign summary | `l9tmF` | DONE ✅ |
| Step 3: Set Budget | CPA pricing, budget cap, duration, estimated reach calculator | `KMtqR` | DONE ✅ |
| Step 4: Review & Submit | Summary, terms acceptance, submit for review | `fZpcQ` | DONE ✅ |

**Community Hub Guided Workspace**: `S1EIA` — Onboarding checklist (5 steps, partially complete), active module cards (Points System, Leaderboard, Tasks), add-more modules (Milestones, Benefits Shop), resources section.

**White Label Hub Redesigned**: `Ir6Tq` — Goal-oriented toolbox: 3 goal cards (Launch Portal / Embed Widgets / Build Pages) with tool tags, 6 toolkit cards (Custom Domain, Widget Library, Page Builder, Brand Settings, SDK & API, Integration Center), resources.

### P7: Hub Active States — DONE ✅

| Hub | ID | Content |
|-----|-----|---------|
| Quest Hub (Active) | `4Mb1C` | Stats row (1,847 participants / 1,203 completions / 65.1% rate / $2,500 rewards) + 2 campaign cards (Active "Token Launch Social Quest" with blue progress bar + Completed "Community Onboarding Quest") |
| Community Hub (Active) | `vFRHi` | Quick stats (1,247 members / 342 active / 24,850 points / 8,931 tasks) + 4/5 checklist + active module cards with real metrics + add-more modules |
| White Label Hub (Active) | `BnkYW` | Stats (Custom Domain active / 3 widgets / 2 pages / 12,450 impressions) + 6 toolkit cards with green configured status (SDK not yet configured) + resources |
| Boost Hub (Active) | `5C3WP` | Stats ($1,500 budget / $423 spent / 846 conversions / $0.50 CPA) + 2 campaign cards (Active "DEX Trading" + Under Review "New User Acquisition Q1" with yellow badge + info banner) |

### P8: Hub Management States — DONE ✅

| Hub | ID | Content |
|-----|-----|---------|
| Quest Hub (Management) | `XvXEQ` | Power stats (14,283 participants / 9,847 completions / 68.9% / $18,500) + filter tabs (All/Active/Draft/Completed) + search + data table 7 columns × 6 rows + pagination |
| Community Hub (Deep) | `TQR51` | Power stats (5,240 members / 1,203 active / 187,400 points / 34,560 tasks) + 6 module cards (2 rows, includes TaskChain + DayChain) + Engagement Overview (weekly bar chart + retention metrics) |
| White Label Hub (Management) | `UPAfV` | Power stats (6 tools / 8 widgets / 5 pages / 47,200 impressions) + 6 fully-configured toolkit cards + Active Deployments (Domain/Widgets/Pages with live badges) + Usage Analytics (6-month bar chart + key metrics panel) |
| Boost Hub (Management) | `8gT3V` | Power stats ($12,500 budget / $7,830 spent / 5,247 conversions / $0.67 CPA) + filter tabs + table 7 columns × 5 rows (CPA color-coded: green/amber) + pagination |

### P9: Dashboard States — DONE ✅

| Dashboard State | ID | Content |
|-----------------|-----|---------|
| Dashboard (New User) | `4SMOO` | Welcome + 4 goal cards + recommended Quest quick start + 3 resource cards |
| Dashboard (Active User) | `IDezm` | Welcome back header + "New Campaign" button + 4 stats (3 campaigns / 2,847 participants / 68.3% rate / $3,200 spent) + 2 active campaign cards (Quest + Community with progress bars) + Quick Actions (Create/Analytics/Duplicate/Export) + WL upsell banner |
| Dashboard (Power User) | `W93vp` | Welcome back + 4 power stats (12 campaigns / 47,280 participants / 72.1% rate / $28,500 spent) + 12-month participant growth bar chart + Product Breakdown panel (Quest 28,450 / Community 8,340 / Boost 5,890 / WL 4,600 = 47,280 total) + Recent Activity feed (4 entries, color-coded by product) |

---

## Execution Phases

| Phase | Content | Pages | Status |
|-------|---------|-------|--------|
| **P1** | Projects Landing Page + B-End Dashboard | 2 | DONE ✅ |
| **P2** | Quest full line (Product Page + Hub + Wizard) | 3 | DONE ✅ |
| **P3** | Community full line + White Label full line | 6 | DONE ✅ |
| **P4** | Boost full line (Product Page + Hub + Wizard) | 3 | DONE ✅ |
| **P5** | Global review + cross-navigation verification + frontend requirements + navigation mapping | — | DONE ✅ |
| **P6** | Wizard completion + Hub strategic redesign (Boost wizard 3 steps + Community guided workspace + WL hub redesign) | 6 | DONE ✅ |
| **P7** | Hub Active States: all 4 products (1-3 campaigns) | 4 | DONE ✅ |
| **P8** | Hub Management States: all 4 products (many campaigns) | 4 | DONE ✅ |
| **P9** | Dashboard Active + Power states | 2 | DONE ✅ |
| **P10** | WL Hub goal cards redefined (3 deployment paths) + Domain Setup + Embed Options + Widget Library (3 states) + Widget Config + Page Builder (3 states) | 8 | DONE ✅ |
| **P11** | Page-button routing map + TODO identification | — | DONE ✅ |
| **P12** | TODO pages T05-T24 (skip T01-T04) | 19 | DONE ✅ |
| **P13** | B-End Sectors & Tasks, Content Management, C-End Community Home, Quest Tab, Leaderboard, B-End Preview Mode | 6 | DONE ✅ |
| **P14** | C-End Sprint, Milestone, Shop tabs + frontend requirements v4.0 update | 3 | DONE ✅ |
| **P0** | Community Empty State redesign (3 retention strategies) + Wizard (4→3 steps) | 4 | DONE ✅ |
| **P0b** | Getting Started Checklist redesign (Community 5-step + WL 6-step activation funnels) | 2 | DONE ✅ |
| **P0c** | WL expandable sidebar sub-menu (Overview/Widgets/Pages) on 8 pages | 8 | DONE ✅ |
| **P0d** | Frontend handoff audit: page codes + button routing + specs | 0 (doc) | DONE ✅ |

---

## P12: TODO Pages (T05-T24) — DONE ✅

### High Priority (T05-T10, skipping T01-T04)

| # | Page | ID | Status |
|---|------|----|--------|
| T01 | Quest Campaign Detail | — | SKIPPED (uses existing patterns) |
| T02 | Quest Campaign Edit | — | SKIPPED (uses existing patterns) |
| T03 | Quest Wizard Step 3 | — | SKIPPED (uses existing patterns) |
| T04 | Quest Wizard Step 4 | — | SKIPPED (uses existing patterns) |
| T05 | Community Wizard Step 1: Choose Type | `cNwNP` | DONE ✅ |
| T06 | Community Wizard Step 3: Set Modules | `Sq7A2` | DONE ✅ |
| T07 | Community Module Detail (Points) | `usBsM` | DONE ✅ |
| T08 | WL Wizard Step 1: Choose Mode | `NNwid` | DONE ✅ |
| T09 | WL Wizard Step 3: Customize | `5nCtO` | DONE ✅ |
| T10 | Boost Campaign Detail | `Sq4jV` | DONE ✅ |

### Medium Priority (T11-T18)

| # | Page | ID | Status |
|---|------|----|--------|
| T11 | WL Brand Settings | `Cx3LH` | DONE ✅ |
| T12 | WL SDK & API | `lQxT5` | DONE ✅ |
| T13 | WL Iframe Embed | `ByGS0` | DONE ✅ |
| T14 | WL Page Analytics | `69HPh` | DONE ✅ |
| T15 | WL Integration Config | `gS64G` | DONE ✅ |
| T16 | Analytics Dashboard | `fLxTr` | DONE ✅ |
| T17 | Settings | `ESrVt` | DONE ✅ |
| T18 | Settings / Profile | `Nh7xq` | DONE ✅ |

### Low Priority (T19-T24)

| # | Page | ID | Status |
|---|------|----|--------|
| T19 | Contact / Book Demo | `4q01T` | DONE ✅ |
| T20 | About Us | `03CDo` | DONE ✅ |
| T21 | Case Studies | `XIChW` | DONE ✅ |
| T22 | Solutions: CPA/CPS | `wsqIT` | DONE ✅ |
| T23 | Solutions: Custom | `A7FaV` | DONE ✅ |
| T24 | Solutions: Joint | `huslr` | DONE ✅ |
| T25 | Blog | — | DEFERRED (likely external link) |

---

## TODO: Missing Pages (from Page-Button Map) — COMPLETED

> Full details: `docs/website_frontend_requirements.md` §六
> Page codes reference: M01-M08 (marketing), B01-B30 (B-end)

### High Priority (core user flows)

| # | Page | URL | Refs | Notes |
|---|------|-----|------|-------|
| T01 | Quest Campaign Detail | `/quest/:id` | 5 | From B02, B05, B06 — campaign cards "View" / table row click |
| T02 | Quest Campaign Edit | `/quest/:id/edit` | 2 | From B02, B05 — "Edit" button |
| T03 | Quest Wizard Step 3: Set Rewards | `/quest/create` step 3 | 1 | B08 "Next" button |
| T04 | Quest Wizard Step 4: Review & Launch | `/quest/create` step 4 | 1 | Step 3 "Next" |
| T05 | Community Wizard Step 1: Choose Type | `/community/create` step 1 | 1 | B13 "Back" button |
| T06 | Community Wizard Step 3: Set Modules | `/community/create` step 3 | 1 | B13 "Next" button |
| T07 | Community Module Detail | `/community/modules/:type` | 6 | B10, B11, B12 module cards |
| T08 | WL Wizard Step 1: Choose Mode | `/white-label/setup` step 1 | 1 | B17 "Back" |
| T09 | WL Wizard Step 3: Customize | `/white-label/setup` step 3 | 1 | B17 "Next" |
| T10 | Boost Campaign Detail | `/boost/:id` | 3 | B28, B29 — campaign cards / table rows |

### Medium Priority (WL sub-pages, auxiliary)

| # | Page | URL | Refs | Notes |
|---|------|-----|------|-------|
| T11 | WL Brand Settings | `/white-label/brand` | 3 | B14, B15, B16 toolkit cards |
| T12 | WL SDK & API | `/white-label/sdk` | 3 | B14, B15, B16 toolkit cards |
| T13 | WL Iframe Embed flow | `/white-label/embed/iframe` | 1 | B19 Iframe Embed card |
| T14 | WL Page Analytics | `/white-label/pages/:id/analytics` | 1 | B25 page card "Analytics" |
| T15 | WL Integration Config (×12) | `/white-label/integrations/:type` | 12 | B26 — all Connect/Configure buttons. May be modal or single config template. |
| T16 | Analytics Dashboard | `/analytics` | 4 | Sidebar, B02, B03, B12, B16 |
| T17 | Settings | `/settings` | 2 | Sidebar, Top Bar |
| T18 | Settings / Profile | `/settings/profile` | 1 | Top Bar avatar dropdown |

### Low Priority (marketing support pages)

| # | Page | URL | Refs | Notes |
|---|------|-----|------|-------|
| T19 | Contact / Book Demo | `/contact` | 10 | M02, M05, M06, M08, Header, Footer. High ref count but simple form. |
| T20 | About Us | `/about` | 1 | Footer |
| T21 | Case Studies | `/case-studies` | 2 | Header, Footer |
| T22 | Solutions: CPA/CPS | `/solutions/cpa` | 2 | Header dropdown, Footer |
| T23 | Solutions: Custom | `/solutions/custom` | 2 | Header dropdown, Footer |
| T24 | Solutions: Joint | `/solutions/joint` | 2 | Header dropdown, Footer |
| T25 | Blog | `/blog` | 1 | Footer (may be external) |

### Summary

| Priority | Count | Key items |
|----------|-------|-----------|
| High | 10 | Campaign detail/edit, wizard missing steps, module detail |
| Medium | 8 (+12 integration sub) | WL Brand/SDK, Analytics, Settings, Integration configs |
| Low | 7 | Contact, About, Case Studies, 3 Solutions, Blog |
| **Total** | **25** (+12 integration sub-pages) | |

> **Recommendation:** T01 (Quest Campaign Detail) and T16 (Analytics Dashboard) have the highest reference counts and block the most user flows. Consider designing these first.

---

## P13: B-End Operations & C-End Community Pages — DONE ✅

Based on `docs/b2c_operation_mapping.md` (153 B-end to C-end operations mapped), designed 4 B-end management pages and 2 C-end user-facing pages + 1 preview mode page.

**C-End Visual Language** (established in P13):
- Dark header (#0F172A) with project branding + wallet connect
- Horizontal nav tabs (amber underline for active)
- User status bar (points, level, streak, rank)
- Amber (#F59E0B) as primary C-end accent
- Gamified task cards with icon boxes + points badges + action buttons
- "Powered by TaskOn" footer

| # | Page | ID | Status | Notes |
|---|------|----|--------|-------|
| P63 | B-End Sectors & Tasks | `Wug7d` | DONE ✅ | Community sidebar. Stats (4 sectors / 12 active / 3 draft / 8,931 completions). Sector 1 "Getting Started" (4 tasks with drag handles, type badges, status pills). Sector 2 "Daily Engagement" (3 tasks). Sector 3 "Advanced Trading" (hidden/collapsed with warning badge). Tip card. Label: `nQ6R7` |
| P64 | B-End Content Management | `lhR14` | DONE ✅ | Community sidebar. "Preview C-End" button. Announcements section (2 items, pinned/active). Featured Slots 2×3 grid (4 filled + 2 empty). Module Status Overview 2×3 (Points/Leaderboard/TaskChain/DayChain active, Benefits Shop + Lucky Wheel not configured). Label: `GB2p5` |
| P65 | C-End Community Home | `vJVhd` | DONE ✅ | Dark header + nav tabs (Home active). User Status Bar. Announcement banner (dark). Featured Grid (3 dark cards). DayChain strip (6 green + day 7 amber TODAY). Sector 1 "Getting Started" (3 task cards: Follow Twitter +50, Join Discord +30, KYC completed +100). Sector 2 "Daily Engagement" (2 task cards). "Powered by TaskOn" footer. Label: `NCqV3` |
| P66 | C-End Quest Tab | `dUXTl` | DONE ✅ | Same dark header, Quests tab active. "All Quests" + filter pills (All/Available/Completed). 2×2 quest card grid: Token Launch Social (500pts), Refer a Friend (150pts), Community Onboarding (completed +200 earned), Advanced Trading (Ends in 2d). Label: `0ivQd` |
| P67 | C-End Leaderboard | `KmdSd` | DONE ✅ | Same dark header, Leaderboard tab active. Time filters (This Week/This Month/All Time). Podium top 3 (#2 CryptoKing 9,450pts, #1 DeFi_Whale 12,800pts with trophy, #3 NFT_Hunter 7,120pts). "Your Rank" highlight card (#42, 1,250pts). Rankings table 5 rows (#4-8) with avatar/name/level/streak/points. Label: `2Obio` |
| P68 | B-End Preview Mode | `2UiNC` | DONE ✅ | B-end sidebar (Community highlighted) + breadcrumb (Community / Content Management / Preview). Amber preview banner with Desktop/Mobile toggle + Exit Preview. Embedded C-end preview showing full community home: dark header, nav tabs, status bar, announcement, featured cards, DayChain, task cards, footer. Label: `0CPEN` |

---

## P14: C-End Remaining Tabs + Document Update — DONE ✅

Analysis of P13 C-End pages revealed 3 missing tabs (Sprint/Milestone/Shop) shown in navigation but not designed. Also updated `docs/website_frontend_requirements.md` from v3.0 to v4.0.

| # | Page | ID | Status | Notes |
|---|------|----|--------|-------|
| P69 | C-End Sprint Tab | `y5fUZ` | DONE ✅ | Current Sprint #12 "Trading Week" with timer + 5 sprint tasks (3 completed/1 in-progress/1 not started) + 3 Reward Tiers (unlocked/2-more/locked) + 2 Past Sprints. Label: `QWqVV` |
| P70 | C-End Milestone Tab | `53iKE` | DONE ✅ | Level badge (Lv.5, 1,250/2,000 pts) + Row 1: 3 milestones (First Steps claimed + Rising Star claimed + Gold Member claimable) + Row 2: 3 locked milestones (Diamond/Whale/Legend, opacity 0.6). Label: `252RE` |
| P71 | C-End Shop Tab | `coM7o` | DONE ✅ | Balance badge (1,250 BTC Points) + category filter pills + Row 1: 3 items (NFT 800pts Redeem + Voucher 500pts Redeem + Whitelist 2,000pts Not Enough) + Row 2: 3 items (Merch 1,500pts Not Enough + Voucher 300pts Redeem + NFT Sold Out 0/50). Label: `3fN8Z` |

**Document update:** `docs/website_frontend_requirements.md` v3.0 → v4.0:
- Added 29 new page codes: B31-B47 (17 B-end), M09-M14 (6 marketing), C01-C06 (6 C-end)
- Total: 67 page codes (M01-M14 + B01-B47 + C01-C06)
- Cleared all TODO markers (replaced with SKIPPED for T01-T04 intentional skips)
- Added complete button routing for all 29 new pages
- Added C-End API table and C-End global components (§2.6, §2.7)
- Updated dev priorities (§七) with C-End phase

---

## After-Login Navigation Architecture

```
┌─ Sidebar ─────────────────┐  ┌─ Top Bar ──────────────────────────┐
│ [TaskOn Logo]              │  │ Breadcrumb ... [Help] [🔔] [Profile]│
│ ─────────────              │  └────────────────────────────────────┘
│ 🏠 Home                   │
│ ─────────────              │
│ Products                   │
│   📋 Quest                 │
│   👥 Community             │
│   🏷️ White Label           │
│   🚀 Boost                 │
│ ─────────────              │
│ 📊 Analytics               │
│ ⚙️ Settings                │
│   └ Team / Billing / API   │
└────────────────────────────┘
```

---

## B-End User Journey (AARRR)

```
Acquisition: Brand Homepage / Search / Ads → Projects Landing Page → Product Pages
    ↓
Interest: Product Pages deep dive → [Start Free] or [Book Demo]
    ↓
Activation: Sign up → Dashboard (new user) → Goal selection → First campaign wizard
    ↓
Retention: Dashboard (active) → Campaign metrics → "Create another" / "Try Community"
    ↓
Revenue: Free → Paid tier → Higher product (Quest→Community→WL)
    ↓
Referral: Success cases → Invite colleagues → Joint programs
```

---

## Before-Login Page Template (Products)

Each product page follows this structure:
1. Hero: Product headline + core promise + CTA
2. "3 Steps to Start" visual flow
3. Core capabilities with product screenshots/mockups
4. Metrics proof (specific numbers from customers)
5. Pricing for this product tier
6. Case study
7. CTA + cross-sell to next product

---

## After-Login Hub Template (3 States)

**Empty State (0 campaigns):**
- Product explainer (1 sentence)
- Template gallery (3 templates)
- Capability highlight strip (compact one-liners for unique features — confidence builder, NOT feature marketing)
- Single CTA: "Create Your First [X]" (emphasize one-click ease)
- Learning resources: Tutorial + Guide + "Learn More" link to Product Page (new tab)

**Active State (1-3 campaigns):**
- Quick stats bar
- Campaign cards with status + metrics + actions
- "Create Another" CTA
- Optimization tips / upsell hints

**Power State (many campaigns):**
- Stats overview dashboard
- List/table view with filters, search, sort, bulk actions
- Analytics + Data export links

---

## Three-Layer Page Hierarchy (WHY / WHAT / HOW)

Each product has three page layers with distinct purposes. Content should NOT be mixed between layers.

| Layer | Product Page (before login) | Empty State (after login) | Wizard (creation flow) |
|-------|---------------------------|--------------------------|----------------------|
| Core Q | WHY — Why choose this? | WHAT — What can I build? | HOW — How do I build it? |
| Audience | Evaluator — "convince me" | New user — "teach me" | Builder — "guide me" |
| Content | Value narrative (strategic) | Capability framework (tactical) | Step-by-step config (operational) |
| Language | "You will get..." | "You can build..." | "Enter / Select..." |

**Cross-layer bridges:**
- Product Page → Empty State: via sign-up/login
- Empty State → Product Page: via "Learn More" in resources (new tab, NOT main nav)
- Empty State → Wizard: via "Create Your First X" CTA

## Product Value Propositions (First Principles)

| Product | Core Value | Human Psychology Lever |
|---------|-----------|----------------------|
| Quest | "Get users fast" — one-time acquisition | Instant gratification (projects get users immediately) |
| Community | "Make users stay" — create leaving costs | Sunk cost (Points), loss aversion (DayChain), goal gradient (TaskChain), endowment (Milestones) |
| White Label | "Own the experience" — branded full-stack | Control/ownership instinct |
| Boost | "Guaranteed results" — pay per outcome | Risk aversion (pay only for results) |

---

## Design System Notes

### Before Login (Marketing)
- Theme: Light (white background)
- Typography: Inter, Hero 48-52px, Section 36px, Body 16px
- Layout: 1440px, padding 120px horizontal
- Components: existing Component Library node `327GX` in pencil-new.pen

### After Login (Product)
- Theme: Light (white/light-gray #F8FAFC background)
- Typography: Inter, Page title 24px, Section 18px, Body 14px
- Layout: Sidebar 240px + Content area fluid
- Density: Higher than marketing (data-forward, compact)
- Cards: White background, subtle border (#E2E8F0), 8-12px border radius

### Brand Homepage (already done)
- Theme: Dark (#0A0F1A / #0F172A / #111827)
- Typography: Inter, Hero 52px
- Special: Floating card composition, radial gradient CTA

---

## Key Node IDs Reference

### Brand Homepage (QszRH) — DONE
- Header: `Dh5nd` (reusable, dark)
- Hero: `75vRS`, title: `GSEll`, CTAs: `NHSGY`
- Product preview: `8MY9q` (layout: "none", floating cards)
- Happening Now: `5bmf8`, leaderboard: `I5YmY`
- Value Props: `KnWao`, For Projects: `dPgTd`
- Social Proof: `Hr8nJ` (includes SUPPORTED CHAINS `OX8JY`)
- Final CTA: `L53RI`, Footer: `R2x4b`

### Existing Page IDs (from pencil-new.pen)
- Projects Landing Page: `Lz2vL`
- Quest Product: `gXQur`
- Community Product: `GyyL4`
- White Label Product: `cbBdG`
- Analytics Product: `iSggN`
- Anti-Sybil Product: `LPxNK`
- SaaS Pricing: `wZxSn`
- CPA/CPS Solution: `v3cov`
- Custom/Managed: `VaGDH`
- Joint Programs: `zeNZr`
- Hybrid Model: `G1Hw2`
- Platform Pricing: `f22ED`
- Enterprise WL: `HKf5j`
- Case Study: `bk7RG`
- Contact/Demo: `ITvVW`
- About: `dKAdu`
- Platform Pricing (new): `CXtOH`
- Enterprise Pricing WL (new): `EDoSn`
- Component Library: `327GX`
- Nav Dropdowns: `kXg2k`, `43K4T`

---

## Legacy Reference (design/legacy/)

### Before Login (for reference only, not reused)
- `boost.png` — CPA product page, dark theme, 3D illustration
- `community.png` — Step-by-step customization showcase
- `quest.png` — Best structured: trust badges + features + templates
- `white label.png` — 6-feature grid + integration modes

### After Login (for reference only)
- `homeafterlogin.png` — Quest + Community dual cards, recommend section
- `questafterlogin.png` — Campaign list with status tabs + template suggestions
- `communityafterlogin.png` — Stats row + customize settings + preview
- `boostafterlogin.png` — Action card grid with pricing per action
- `WL-page-afterlogin.png` — Page Builder with SDK credentials
- `WL-widget-afterlogin.png` — Widget Builder (nearly identical to page builder)
- `WLafterlogin.png` — Custom Domain mode explainer card

### Legacy Problems (don't repeat):
1. Feature exhibition instead of problem-solving
2. Empty states with no guidance
3. No product cross-linking or progression
4. AARRR funnel completely broken
5. Dark theme reducing B-end trust

---

## P0: Community Redesign + Onboarding Overhaul — DONE ✅

### P0 — Community Empty State + Wizard Redesign
- **`zzZ8D` (B09)**: Replaced 3 module-level templates with 3 lifecycle-stage retention strategies:
  1. Activate New Users (green): TaskChain + Points + Levels
  2. Drive Daily Engagement (orange): DayChain + Leaderboard + Sprint
  3. Maximize Retention (purple): Benefits Shop + Milestones + Lucky Wheel
- Added "HOW IT WORKS" lifecycle flow: Quest → Activate → Engage → Retain
- **`Gzpeu` (B13)**: Wizard Step 1 Customize (3 steps, was 4)
- **`8NeyG` (B13b)**: Wizard Step 2 Configure Modules (9 toggles, strategy pre-fill)
- **`qknQZ` (B13c)**: Wizard Step 3 Review & Launch (summary + launch button)

### P0b — Getting Started Checklist Redesign
Philosophy: **Configure → Verify → Distribute → Celebrate** (activation funnels, not feature config lists)
- **`S1EIA` Community** (5 steps): Create → Configure → Preview → **Share** (link+social+Promo Kit) → First 10 participants
- **`BnkYW` WL** (6 steps): Create → Deploy → **Send Dev Kit** (link+email+auto-detect) → Verified → Announce → First interaction
- Key concepts: Promo Kit Generator (AI social posts + banner), Dev Kit Page (standalone, no account needed)

### P0c — WL Expandable Sidebar Sub-menu
Added to all 8 WL pages: `keyboard_arrow_up` chevron + 3 sub-items (Overview/Widgets/Pages)
- Overview active: `BnkYW`, `UPAfV`
- Widgets active: `2sSsA`, `n4pJK`, `S432k`
- Pages active: `DRYwN`, `sGDcq`, `J08v5`

### P0d — Frontend Handoff Audit

Comprehensive audit and correction of all page codes, button routing, and frontend requirements:

1. **Community Wizard page codes fixed**: B13=Step1(`Gzpeu`), B34=Step2(`8NeyG`), B35=Step3(`qknQZ`). Old nodes `cNwNP`/`Sq7A2` removed.
2. **Frontend requirements v4.1**: Added §六 (Promo Kit Generator, Dev Kit Page B48, Auto-detection Systems), updated B10/B15 checklist routing, WL sidebar sub-menu in §2.4.
3. **Canvas annotations updated**: `hrASq` (S1EIA) and `BHGp8` (BnkYW) reflect new checklist buttons.
4. **layout.md fully corrected**: All B-codes now match `website_frontend_requirements.md` (was using wrong sequential codes since Row 10).
5. **8 Community module management pages coded**: B31a-B31h (Points/TaskChain/DayChain/Leaderboard/Sprint/Milestone/Shop/Wheel).
6. **B48 Dev Kit Page assigned**: Conceptualized, routing added, page design pending.
7. **Total page codes**: 76 (M01-M14 + B01-B48 + B31a-B31h + C01-C06).

---

## Design Review — Complete ✅

> Full issue tracker: `design/review-issues.md`

### Review Summary

Two-round systematic review of all 66 pages completed 2026-03-04.

| Round | Dimensions | Issues |
|-------|-----------|--------|
| Round 1 | UX (first principles), Path (AARRR), State coverage, Layout clarity | 71 (3 P0 + 35 P1 + 33 P2) |
| Round 2 | Layout clipping (snapshot_layout), Visual quality, Navigation completeness, Content consistency | 19 (1 P0 + 8 P1 + 10 P2) |
| **Total** | | **90 (4 P0 + 43 P1 + 43 P2)** |

### 4 P0 Critical Issues

| # | Page | Issue |
|---|------|-------|
| 1 | M07 Platform Pricing `CXtOH` | Free($0) tier 完全缺失 |
| 2 | M07 Platform Pricing `CXtOH` | 缺 Feature Comparison Table |
| 3 | M08 WL Pricing `EDoSn` | 缺 Feature Comparison Table |
| 72 | B11 Community Active `vFRHi` | Resources Row 完全不可见 (fully clipped) |

### Recommended Fix Phases

| Phase | Scope | Est. Issues | Priority |
|-------|-------|-------------|----------|
| **Fix-A** | 4 P0: Pricing Free tier + 2 Feature Tables + Community Active clipping | 4 | Must fix |
| **Fix-B** | 系统性文本裁剪修复 (7 pages, Pencil text width limitation) | 7 | Should fix |
| **Fix-C** | P1 功能性问题 (按 batch 分组: Dashboard/Community/Boost/WL/Auxiliary/C-End/Marketing) | ~32 | Should fix |
| **Fix-D** | P2 优化项 (选择性修复, 部分 Won't Fix) | ~43 | Nice to have |

---

## Reference Documents

| Doc | Path | Key Content |
|-----|------|-------------|
| First Principles | `docs/TaskOn第一性原理与人性洞察准则.md` | Efficiency formula, human nature insights, B-end decision heuristics |
| Growth Playbook | `docs/taskon_growth_engine_playbook_en.md` | Full AARRR framework, lifecycle playbooks, competitor analysis |
| Services | `docs/taskon_services_en.md` | 4 revenue streams, pricing models, product details |
| Knowledge Base | `docs/TaskOn_Knowledge_Base_Init_V1.md` | Product specs, widget system, integration details |
| Keywords | `docs/taskon_keywords.md` | Persona definitions, pain points, competitor comparison, SEO |
| Frontend Requirements | `docs/website_frontend_requirements.md` | Static/dynamic content, API specs, dev priorities |
