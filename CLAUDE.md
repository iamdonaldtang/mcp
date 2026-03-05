# TaskOn Design Project — Claude Instructions

## What This Project Is

TaskOn is a Web3 Growth Engine platform. This repo contains design files, strategy documents, and reference materials for the complete TaskOn website and B-end product interface redesign.

## Active Work

**P1-P14 COMPLETE — 66 pages designed (67 page codes)**
- **Marketing (M01-M14)**: 8 core pages + 6 support pages (Contact, About, Case Studies, 3 Solutions)
- **B-End (B01-B47)**: Dashboard 3 states + 4 Product Hubs (3 states each) + Wizards + WL sub-pages (Domain/Embed/Widget/Page Builder/Integration/Brand/SDK) + Analytics + Settings + Community Operations (Sectors/Content/Preview)
- **C-End (C01-C06)**: White Label Community front-end (Home/Quests/Leaderboard/Sprint/Milestone/Shop)
- **Skipped**: T01-T04 (Quest campaign detail/edit/wizard steps 3-4 — follow existing patterns). T25 (Blog — external link).

**Design phase complete. Next step: frontend development.**
**Full progress tracker**: `design/progress.md` (READ THIS FIRST in every new session)
**Frontend requirements**: `docs/website_frontend_requirements.md` v4.0 — 67 page codes, complete button→page routing, C-End APIs

## Critical Rules

### Design File Access
- The design file is `design/pencil-new.pen`
- **ONLY use Pencil MCP tools** to read/write .pen files (batch_get, batch_design, get_screenshot, etc.)
- **NEVER use Read, Grep, or cat** on .pen files — contents are encrypted
- Always verify changes visually with `get_screenshot` after major edits
- Always run `snapshot_layout(problemsOnly:true)` to check for layout issues
- Max 25 operations per `batch_design` call

### Design Conventions
- Icons: **Material Symbols Rounded** (NOT Lucide — many Lucide icons fail/missing)
- Absolute positioning requires `layout: "none"` on the parent frame
- Marketing pages: Light theme (white background), 1440px width, 120px horizontal padding
- Brand Homepage: Dark theme (#0A0F1A), already complete
- B-end after-login: Light theme (#F8FAFC), sidebar navigation, higher information density
- Typography: Inter font throughout
- Product brand colors: Quest=#5D7EF1, Community=#48BB78, White-Label=#9B7EE0, Boost=#ED8936
- C-end (White Label community): Dark header (#0F172A) + light content (#F8FAFC), amber (#F59E0B) accent, horizontal nav tabs, "Powered by TaskOn" footer

### Strategy Documents (Read When Needed)
- `docs/TaskOn第一性原理与人性洞察准则.md` — First principles and design philosophy
- `docs/taskon_growth_engine_playbook_en.md` — Growth strategy, AARRR model, competitor analysis
- `docs/taskon_services_en.md` — Product details, pricing models, revenue streams
- `docs/TaskOn_Knowledge_Base_Init_V1.md` — Full product specifications
- `docs/taskon_keywords.md` — Target personas, pain points, messaging
- `docs/website_frontend_requirements.md` — v4.0: 67 page codes, complete button routing, C-End APIs, dev priorities
- `docs/b2c_operation_mapping.md` — 153 B-end to C-end operation mappings, entity lifecycles
- `docs/white_label_integration_modes.md` — WL 4 integration modes, dependency chain
- `design/legacy/` — Old UI screenshots (reference only, known to have major UX problems)

### Communication
- User communicates in Chinese; design content is in English
- Always get alignment before starting large design work
- Think from first principles, not feature lists
- B-end design = prove ROI, minimize cognitive/implementation/time costs

## Key Architecture Decisions

1. **4 Products**: Quest (self-serve acquisition) → Community (retention) → White Label (full-stack) → Boost (CPA/CPS pay-per-result)
2. **C-end and B-end are completely independent systems** — separate URLs, nav, user models
3. **B-end Dashboard**: 3 user states (new/active/power) with adaptive UI
4. **Product Hubs**: 3 states each (empty/active/management)
5. **AARRR funnel**: Landing Page → Product Pages → Sign Up → Dashboard → First Campaign → Retention → Upsell
6. **Navigation after login**: Sidebar (Home / Products / Analytics / Settings) + Top bar (breadcrumb + help + profile)

## Product Value Propositions (First Principles)

| Product | Core User Need | Mechanism |
|---------|---------------|-----------|
| Quest | "Get users fast" — one-time acquisition events | Pay → traffic → task completion → user list |
| Community | "Make users stay" — create leaving costs for zero-switching-cost crypto users | Sunk cost (Points), loss aversion (DayChain), goal gradient (TaskChain), endowment (Milestones) |
| White Label | "Own the experience" — full-stack branded growth | Custom domain, SDK, data ownership |
| Boost | "Guaranteed results" — pay only for outcomes | CPA/CPS model, anti-sybil, managed delivery |

## Three-Layer Page Hierarchy (WHY / WHAT / HOW)

Each product has three page layers. Do NOT mix content between layers.

| Layer | Product Page (before login) | Empty State (after login) | Wizard (creation flow) |
|-------|---------------------------|--------------------------|----------------------|
| **Core Q** | WHY — Why choose this? | WHAT — What can I build? | HOW — How do I build it? |
| **Audience** | Evaluator — "convince me" | New user — "teach me" | Builder — "guide me" |
| **Content** | Value narrative (strategic) | Capability framework (tactical) | Step-by-step config (operational) |
| **Success** | Interest → sign up | Understand → click Create | Configure → activate |
| **Language** | "You will get..." | "You can build..." | "Enter / Select..." |
| **Tone** | Inspire desire | Build confidence | Reduce friction |

**Cross-layer bridges:**
- Product Page → Empty State: via sign-up/login (natural conversion)
- Empty State → Product Page: via "Learn More" link in resources (new tab, NOT main nav)
- Empty State → Wizard: via "Create Your First X" CTA (primary conversion path)

**Empty State design rules:**
1. Show framework, don't sell features — templates show "what you can build", not feature marketing
2. Use highlight strips for differentiation — compact one-liners for unique capabilities (TaskChain/DayChain/Smart Task etc.), NOT full feature sections
3. Single primary action — one CTA to start creating, emphasize ease ("one click")
4. Resources as depth exit — "Learn More", tutorials, playbooks for users who want to go deeper
