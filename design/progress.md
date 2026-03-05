# TaskOn Design — Progress Tracker

> Last updated: 2026-03-05
> Design file: `design/pencil-new.pen` (Pencil MCP tools only)
> Page IDs: `design/pages.md` | TODO: `design/todo.md` | Archive: `design/progress-archive.md`

## Status

**72 pages designed, 82 page codes** (M01-M14 + B01-B50 + B31a-B31i + C01-C09)
- Design Review: 90 issues → Fixed 68 / Deferred 19 / Remaining 5 (see `todo.md`)
- Skipped: T01-T04 (Quest campaign detail/edit/wizard 3-4, follow existing patterns), T25 (Blog, external)

---

## Key Decisions

| Decision | Result |
|----------|--------|
| Products | 4 independent: Quest, Community, White Label, Boost |
| B-end theme | Light (#F8FAFC, professional SaaS) |
| C-end / B-end | Completely independent systems |
| Header | Ultra-minimal: Logo + [For Projects] + [Launch App] |
| Pricing | Quest+Community: Free/$79mo; WL: $499/$999mo; Boost: CPA (no pricing page) |
| WL = Quest + Community + exclusive | Page Builder, widgets, custom domain, SDK |
| Community modules | 4 systems (Task Engine / Points & Recognition / Incentive Campaigns / Rewards Economy) + 2 infrastructure (Access Rules / Homepage Editor) |
| Leaderboard vs LB Sprint | Leaderboard = recurring, no incentives; LB Sprint = time-bounded + NFT/Token incentives |

---

## Execution Phases

| Phase | Content | Pages | Status |
|-------|---------|-------|--------|
| **P1** | Projects Landing + Dashboard | 2 | DONE ✅ |
| **P2** | Quest (Product + Hub + Wizard) | 3 | DONE ✅ |
| **P3** | Community + White Label full line | 6 | DONE ✅ |
| **P4** | Boost full line | 3 | DONE ✅ |
| **P5** | Global review + navigation + frontend requirements | — | DONE ✅ |
| **P6** | Wizard completion + Hub strategic redesign | 6 | DONE ✅ |
| **P7** | Hub Active States (4 products) | 4 | DONE ✅ |
| **P8** | Hub Management States (4 products) | 4 | DONE ✅ |
| **P9** | Dashboard Active + Power | 2 | DONE ✅ |
| **P10** | WL deployment paths + Domain/Embed/Widget/PageBuilder | 8 | DONE ✅ |
| **P11** | Page-button routing map | — | DONE ✅ |
| **P12** | TODO pages T05-T24 | 19 | DONE ✅ |
| **P13** | B-End Sectors/Content/Preview + C-End Home/Quest/Leaderboard | 6 | DONE ✅ |
| **P14** | C-End Sprint/Milestone/Shop + frontend req v4.0 | 3 | DONE ✅ |
| **P0** | Community Empty State + Wizard redesign (3 retention strategies) | 4 | DONE ✅ |
| **P0b** | Getting Started Checklists (Community 5-step + WL 6-step) | 2 | DONE ✅ |
| **P0c** | WL expandable sidebar (Overview/Widgets/Pages) on 8 pages | 8 | DONE ✅ |
| **P0d** | Frontend handoff audit: page codes + button routing + specs | 0 (doc) | DONE ✅ |
| **P0e** | C-End engagement: C01 action engine + C07-C09 new pages | 3+1 | DONE ✅ |
| **P0f** | Leaderboard vs LB Sprint distinction across all pages | 0 (updates) | DONE ✅ |
| **P0g** | 4-System Module Architecture + Badges/AccessRules/HomepageEditor | 3+updates | DONE ✅ |

---

## Page Summary (see `design/pages.md` for Node IDs)

### Marketing (M01-M14) — 14 pages
Brand Homepage (dark) · Projects Landing · Quest/Community/WL/Boost Product Pages · Platform Pricing · WL Pricing · Contact · About · Case Studies · 3 Solutions

### B-End (B01-B50 + B31a-B31i) — 59 page codes
- **Dashboard**: New/Active/Power (3 states)
- **Quest**: Hub 3 states + Wizard 2 steps
- **Community**: Hub 4 states + Wizard 3 steps + 12 module/settings pages (Sectors/Points/TaskChain/DayChain/Leaderboard/LBSprint/Milestone/Shop/Wheel/Badges/AccessRules/HomepageEditor) + Content Mgmt + Preview Mode
- **White Label**: Hub 3 states + Wizard + Domain/Embed + Widget 3 states + PageBuilder 3 states + Integration/Brand/SDK/Analytics
- **Boost**: Hub 3 states + Wizard 4 steps
- **Global**: Analytics Dashboard + Settings + Profile

### C-End (C01-C09) — 9 page codes
Home · Quest · Leaderboard · LB Sprint · Milestone · Shop · User Center · Invite Center · Activity Feed

---

## Reference Docs

| Doc | Path |
|-----|------|
| Frontend Requirements | `docs/website_frontend_requirements.md` v4.1 (82 page codes, button routing, APIs) |
| B2C Operation Mapping | `docs/b2c_operation_mapping.md` (153 operations) |
| First Principles | `docs/TaskOn第一性原理与人性洞察准则.md` |
| Growth Playbook | `docs/taskon_growth_engine_playbook_en.md` |
| Knowledge Base | `docs/TaskOn_Knowledge_Base_Init_V1.md` |
| WL Integration Modes | `docs/white_label_integration_modes.md` |
| Canvas Layout | `design/layout.md` |
| Full History | `design/progress-archive.md` (detailed per-phase execution notes) |
