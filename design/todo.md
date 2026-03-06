# Design TODO

> Active design tasks for `design/pencil-new.pen`
> Completed items tracked in `design/progress.md`
> Last review: 2026-03-06 — Onboarding flow audit (Community + White Label, 13 pages)
> Previous batch: 2026-03-05/06 — 40 items completed (see Archive at bottom)

---

## OB — Onboarding Flow Restructuring (Active)

> **Root cause**: Both Community and White Label wizards only make structural decisions (name, toggle modules, choose path) but skip content creation, tool configuration, and user-side preview. Users "launch" empty shells.
> **Principle**: Wizard creates a Minimum Viable Community/Deployment with template content → Guided state helps flesh out and customize → Active state is data-driven monitoring.
> **Dependency**: Community OB must be resolved first — WL depends on a functioning Community.

### Diagnosis Summary

| Product | What Wizard Does | What Wizard Skips | Result |
|---------|-----------------|-------------------|--------|
| Community | Name + toggle modules | Content creation (tasks, rewards, points rules), C-end preview | Empty community with enabled-but-unfilled modules |
| White Label | Choose integration mode + domain DNS + brand colors | Widget configuration, page building, integration testing, end-to-end preview | Branding applied to nothing deployable; Dev Kit has no content |

**Community gaps**: Wizard Step 2 toggles modules ON but creates zero content for any of them. Step 3 "Review & Launch" reviews metadata only (name + which systems are on), has no C-end preview, and launches an empty community.

**White Label gaps**: (1) Empty State has 15 clickable targets but no single entry CTA. (2) 3 deployment paths on Empty != 4 modes in Wizard (naming/count mismatch). (3) Wizard Step 2 only has Domain Configuration — recommended path "Widget + Page Builder" has no Step 2 page. (4) 6 tool sub-pages are standalone islands with no guided sequence. (5) Checklist puts "Send Dev Kit" before widget/page configuration.

**Cross-product**: Same decision appears on multiple pages with different options each time (WL deployment choice x3, domain config x2, brand customization x2).

---

### OB-C: Community Onboarding Restructuring

> Flow: Empty `zzZ8D` → Wizard (4 steps) → Guided `S1EIA` → Active `vFRHi` → Deep `TQR51`
> Wizard changes from 3 steps to 4 steps. Step 3 (Quick Setup) and Step 4 (Preview & Publish) are new.

#### OB-C1: Redesign Empty State `zzZ8D`

**Current**: 3 strategy cards (decorative) + "Create Your First Community" button + "Or configure individual modules →"
**Problem**: Strategy cards are informational only — selection is not carried into Wizard. Secondary CTA "configure individual modules" has no valid target.

**New design**:
- 3 strategy cards become **selectable** (radio card behavior):
  - Click a card → border highlight + checkmark → card expands to show:
    - "This strategy includes:" list of auto-enabled modules
    - "Expected outcome:" one-line result metric
    - Mini C-end screenshot thumbnail (what users will see)
  - Only one card selected at a time
- Primary CTA (below cards): **"Create Community with This Strategy →"**
  - Appears/activates after a strategy is selected
  - Carries strategy choice into Wizard (pre-selects modules in Step 2)
- Secondary CTA: **"Or start from scratch →"** (enters Wizard with no pre-selection)
- **Remove** "Or configure individual modules →" (no valid destination)
- Keep "How It Works" 4-step strip and Resources section

**Strategy → Module mapping** (pre-selection in Wizard Step 2):
| Strategy | Auto-enabled modules |
|----------|---------------------|
| Activate New Users | Sectors & Tasks, Points & Level, TaskChain |
| Drive Daily Engagement | Sectors & Tasks, Points & Level, DayChain, Leaderboard |
| Maximize Retention | Sectors & Tasks, Points & Level, DayChain, Milestones, Benefits Shop |

#### OB-C2: Update Wizard Step 1 — Customize `Gzpeu`

**Current**: Name, description, brand color, live preview. Stepper: Customize → Configure Modules → Review & Launch
**Changes**:
- Update stepper to 4 steps: **Customize → Modules → Quick Setup → Preview & Publish**
- Show selected strategy name above the form: "Strategy: Activate New Users" (with "Change" link)
- Keep rest of page as-is (name, description, brand color, live preview work well)

#### OB-C3: Redesign Wizard Step 2 — Modules `8NeyG`

**Current**: 4 system cards with toggles + sub-module checkboxes. Right panel shows count summary.
**Problem**: Modules are abstract toggles — project owner doesn't understand what each does or what enabling it commits them to.

**New design**:
- Selected strategy shown at top: "Recommended for: [strategy name]" with strategy-based modules pre-checked
- Each module gets enhanced information:
  - Current: checkbox + name + one-line description
  - **Add**: C-end mini preview thumbnail (40x30px, what the module looks like to users)
  - **Add**: Effect metric: "Users who complete TaskChains have 2.3x higher retention"
  - **Add**: Expandable "What this does" panel (2-3 sentences, shows before/after user behavior)
- Required modules (Sectors & Tasks, Points & Level): keep "Required" badge, non-dismissable
- Right panel summary becomes concrete:
  - Current: "2 systems · 3 modules"
  - **New**: "Your community will include: Task system with 3 starter tasks, Points with 5 levels, Daily check-in with streak bonuses"
  - Dynamic based on checked modules
- Remove strategy selector from this page (strategy was chosen in Empty/Step 1)

#### OB-C4: New Wizard Step 3 — Quick Setup (replaces `qknQZ` Review & Launch)

**Current `qknQZ`**: 3 summary cards (Community Details / Enabled Systems / What Happens Next) + "Ready to Launch!" metrics + "Launch Community" button.
**Problem**: Reviews only metadata. No content created. No preview. Launches empty community.

**Complete redesign — new purpose: auto-generate template content with inline customization**:

Layout: Left content area (70%) + Right summary panel (30%)

**Left — Module Quick Setup cards** (one per enabled module, vertical stack):

Each card shows template content pre-filled, with inline edit capability:

| Module | Template Content | Inline Edit Fields |
|--------|-----------------|-------------------|
| Sectors & Tasks | "Getting Started" sector + 3 tasks (Follow on Twitter, Join Discord, Visit Website) | Task titles editable, can delete/add tasks |
| Points & Level | Point name "XP", daily cap 500, 5 levels (Newcomer/Active/Contributor/Expert/Legend) | Point name, daily cap, level names |
| TaskChain | 1 chain "Welcome Journey" (3 steps: Complete profile → Do first task → Invite a friend) | Step names, linked tasks |
| DayChain | Check-in reward 10 XP/day, Day 7 bonus 2x, Day 30 bonus 5x | Reward values, bonus days |
| Leaderboard | Weekly reset, rank by XP, top 10 displayed | Reset period, point type, display count |
| LB Sprint | (No template — needs manual setup, shown as "Set up after launch" with link) | — |
| Milestone | 3 milestones (First Task 50XP, 7-Day Streak 200XP, Level 3 500XP) | Names, thresholds, rewards |
| Benefits Shop | (No items — "Add your first reward after launch" prompt) | — |
| Lucky Wheel | (No template — shown as "Set up after launch") | — |
| Badges | 3 badges (Early Bird, Streak Master, Top Contributor) with conditions | Badge names, conditions |

- Header: "We've prepared a starter pack for your community"
- Sub-header: "Everything is pre-filled with best practices. Edit anything you'd like, or keep defaults."
- Estimated time: "This takes about 3 minutes"
- Each module card: collapsible, shows template summary when collapsed, edit fields when expanded

**Right — Summary panel**:
- "Your community at launch:"
- Checklist: ✅ 1 sector · ✅ 3 tasks · ✅ 5 levels · ✅ Daily check-in · ⚠ No shop items (add later)
- "Users will be able to: complete tasks, earn XP, level up, maintain streaks"

**Bottom CTA**: "Next: Preview & Publish →"

#### OB-C5: New Wizard Step 4 — Preview & Publish (new page, needs new Node ID)

**Purpose**: Let project owner see the C-end experience before publishing. Bridge the gap between configuration and launch.

**Layout**: Full-width preview area (75%) + Readiness panel (25%)

**Left — C-end Preview**:
- Embedded preview of Community Home page as users will see it
- Tab bar showing available C-end pages: Home / Quests / Leaderboard (based on enabled modules)
- Desktop/Mobile toggle (top-right of preview frame)
- Preview uses actual brand colors from Step 1 + template content from Step 3
- Banner at top of preview: "This is how your community looks to users" + "Content is from your starter pack. You can customize everything after publishing."

**Right — Readiness Checklist** (auto-generated):
- ✅ Community name: "Arbitrum Builders"
- ✅ Brand configured (green + custom color)
- ✅ 3 tasks created
- ✅ Points system: XP, 5 levels
- ✅ DayChain: 10 XP/day
- ⚠ Benefits Shop: No items yet (optional, add after publish)
- ⚠ Lucky Wheel: Not configured (optional)
- Overall: "Ready to publish — 5/7 features configured"

**Bottom actions**:
- Primary: **"Publish Community"** (green button)
  - Publishes in Draft state (accessible via URL but not promoted)
  - Post-click: celebration animation (confetti + "Your community is live!") → auto-redirect to Guided state after 3 seconds
- Secondary: "Save as Draft" (save progress, return later)
- Link: "← Back to Quick Setup" (go back and edit)

#### OB-C6: Redesign Guided State `S1EIA`

**Current**: "My Community" title + "Active" badge + 5-step checklist (3 auto-done + Share + First 10) + module cards with "Configure" + Resources
**Problem**: Steps 1-3 are structural (already done by wizard), steps 4-5 are promotion. Entire content creation phase missing. "Active" badge conflicts with separate Active state.

**New design**:

**Header changes**:
- Title: **"Getting Started"** (not "My Community" — that's for Active state)
- Badge: **"Setting Up"** (amber/yellow, not "Active" green)
- Sub-text: "Complete these steps to get the most out of your community"

**Checklist — module-level content creation tasks** (dynamic based on enabled modules):

```
Getting Started                                          3 of 9 complete

COMPLETED BY WIZARD
✅ Community created with "Activate New Users" strategy
✅ 3 starter tasks live in "Getting Started" sector
✅ Points & Levels configured (XP, 5 levels)

ENRICH YOUR COMMUNITY
○ Add more tasks to your sectors                         → Go to Sectors
  You have 3 tasks. Communities with 10+ tasks
  see 2x higher engagement.

○ Set up your Benefits Shop                              → Go to Shop
  Give users something to spend points on.
  Quick: add 1 reward in 2 minutes.

○ Customize your DayChain rewards                        → Go to DayChain
  Default: 10 XP/day. Consider higher Day 7 bonus
  to prevent streak dropout.

GO LIVE
○ Preview your community as a user                       → Open Preview
  Walk through the experience your users will have.

○ Share with your community                              → Promo Kit
  Community link · Twitter · Discord · Telegram
  AI-generated social post + branded banner

○ First 10 participants                                  auto-detect
  0/10 — waiting for your first users
```

- "ENRICH" section is **dynamic** — only shows items for enabled modules that need attention
- Each item has: task description + motivation metric + direct link to relevant module management page
- Items auto-complete when conditions are met (e.g., "Add more tasks" completes when task count > 5)

**Module cards section** (below checklist):
- Title: "Your Active Modules" (same as current)
- Cards show: module name + current status + "Configure" or "Manage" button
- Keep "Add More Modules" row
- Keep "Browse Configuration Templates" link

**Resources section**: Keep as-is

#### OB-C7: Redesign Active State `vFRHi`

**Current**: Stats row + Getting Started (collapsed, 4/5) + Active Modules (with data) + Add More + Resources
**Problem**: Still shows checklist. Layout nearly identical to Guided. No clear visual progression.

**New design — data-driven, no onboarding content**:

**Header**: "My Community" + **"Active"** badge (green) — this is now the only page with this badge
**Sub-text**: "Your community is live. Complete setup steps below." → **Remove** (no setup reference)

**Stats row** (keep, enhance):
- Same 4 stats: Total Members / Active This Week / Points Distributed / Tasks Completed
- **Add**: WoW trend arrows (↑12% / ↓3%) to each stat
- **Add**: Sparkline mini-chart below each number

**Module Performance Cards** (replaces "Active Modules"):
- Each enabled module as a card, but now showing **performance data**:
  - Sectors & Tasks: "4,231 completions this week (+16%)" + completion rate + active tasks count
  - Points: "Avg 72 XP earned/user/day" + economy balance indicator
  - DayChain: "Avg streak: 8.3 days" + active streak rate + trend
  - Leaderboard: "Top scorer: 3,450 XP" + weekly board resets + participant count
- Each card: "Manage →" button (goes to module management page)
- Cards sorted by: modules needing attention first (red/amber indicators)

**Remaining onboarding** (minimal):
- IF checklist has incomplete items: single-line amber banner at top: "1 step remaining: Get your first 10 participants → View"
- NO expanded checklist, NO checklist section

**Quick Actions row** (new):
- "Create Task" / "Add Reward" / "View Analytics" — 3 shortcut buttons for common actions

**Resources**: Keep but reduce to 2 cards (remove "Learn More" since user is already active)

#### OB-C8: Stepper Update (all Wizard pages)

- Update stepper on `Gzpeu`, `8NeyG`, `qknQZ` (repurposed), and new Step 4 page
- Old: Customize → Configure Modules → Review & Launch (3 steps)
- New: **Customize → Modules → Quick Setup → Preview & Publish** (4 steps)
- Stepper dot style: same as current (circle + label)

---

### OB-W: White Label Onboarding Restructuring

> Flow: Empty `Ir6Tq` → Wizard (4 steps) → Active/Guided `BnkYW` → Mgmt `UPAfV`
> Prerequisite: Community must be in Active or Guided state (not Empty).
> Wizard restructured to 4 clear steps with path-adaptive content.

#### OB-W1: Redesign Empty State `Ir6Tq`

**Current**: Yellow banner "Community Setup Required" + 3 deployment paths + 6 toolkit cards + Resources
**Problems**: (1) Banner is advisory, not blocking. (2) No single entry CTA. (3) 15+ clickable targets. (4) Toolkit cards premature (tools are post-setup). (5) Path names don't match Wizard modes.

**New design — two states based on Community prerequisite**:

**State A: Community NOT ready** (blocking):
- Full-page prerequisite screen (not banner):
  - Icon: shield/lock icon (purple, 48px)
  - Title: "Set Up Community First"
  - Text: "White Label wraps your Community in your own brand and deploys it to your product. You need an active Community with content before setting up White Label."
  - What you'll need: "At least 1 sector with tasks, Points system configured, and optionally a Benefits Shop"
  - Primary CTA: **"Set Up Community →"** (goes to Community Empty `zzZ8D` or Guided `S1EIA`)
  - Secondary: "Learn more about White Label →" (link to WL marketing page `cbBdG`, new tab)
- No toolkit cards, no deployment paths — nothing else to do until Community is ready

**State B: Community IS ready** (normal Empty):
- Community status bar at top: "Your Community: [name] — 12 tasks, 5 levels, 245 members" (green checkmark)
- Title: "White Label" + description
- 3 deployment paths as **selectable cards** (radio behavior, like Community strategy cards):
  - **"Embed in Your App"** ★ Recommended
    - "Add TaskOn growth features directly into your existing website or app"
    - Includes: Widget Library, Page Builder, Iframe Embed
    - Best for: "Teams with an existing web app"
    - Effort: "< 1 day integration"
  - **"Host on Your Domain"**
    - "Run a full community portal on your own domain (e.g., community.yourproject.io)"
    - Includes: Custom domain, SSL, full portal
    - Best for: "Teams who want a standalone community site"
    - Effort: "30 min DNS setup"
  - **"Build with SDK"**
    - "Full programmatic control — build completely custom experiences"
    - Includes: REST API, Webhooks, SDK
    - Best for: "Teams with dedicated developers"
    - Effort: "1-2 weeks development"
- Primary CTA: **"Set Up White Label →"** (purple button, carries path selection to Wizard)
- Remove 6 toolkit cards (user hasn't set anything up yet — tools are premature)
- Keep Resources section (3 cards: SDK Documentation, Setup Walkthrough, Learn More)

#### OB-W2: Redesign Wizard Step 1 — Choose Path `NNwid`

**Current**: 4 integration modes (Widget+PB / Custom Domain / Iframe Embed / Full SDK). Stepper: Mode → Configure → OK → Customize.
**Problem**: 4 modes != Empty's 3 paths. Stepper has 4 steps but only 3 pages. "OK" step undefined.

**New design**:
- Stepper: **Path → Configure → Brand → Preview** (4 steps, all have pages)
- 3 paths (aligned with Empty State):
  - **"Embed in Your App"** ★ Recommended
    - Sub-options (expandable after selecting this card):
      - Widget Library: "Embed individual components (leaderboard, tasks, points)"
      - Page Builder: "Build full pages combining multiple widgets"
      - Iframe: "Embed entire community in an iframe (simplest)"
    - User can check multiple sub-options
  - **"Host on Your Domain"**
    - Shows: "You'll configure DNS in the next step"
  - **"Build with SDK"**
    - Shows: "You'll get API keys and SDK setup in the next step"
- If path was pre-selected from Empty State, that card is pre-highlighted
- CTA: "Next: Configure →"

#### OB-W3: Redesign Wizard Step 2 — Configure (path-adaptive) `CXzmy`

**Current**: Domain Configuration only (DNS + SSL). Only serves "Host on Your Domain" path.
**Problem**: 3 of 4 paths have no Step 2 content.

**New design — page content changes based on selected path**:

**IF "Embed in Your App"** (recommended path):
- Title: "Choose Your Widgets"
- Left: List of Community modules available as widgets (same list as Widget Library `2sSsA`):
  - Each module: checkbox + name + mini-preview + configuration status from Community
  - Green = ready to embed (module configured in Community)
  - Amber = needs Community setup first (with "Set up in Community →" link)
  - Pre-checked: Leaderboard, Task List, User Center (most common)
- Right: Live preview showing selected widgets in a mock embed context
- Bottom: "Selected: 3 widgets — you'll configure each one after setup"
- CTA: "Next: Brand →"

**IF "Host on Your Domain"**:
- Title: "Domain Configuration" (keep current `CXzmy` content)
- Custom Domain input + SSL Certificate + DNS Status panel
- Portal Preview on right
- CTA: "Next: Brand →"

**IF "Build with SDK"**:
- Title: "SDK Setup"
- Auto-generate: Production API Key + Test API Key (shown with copy buttons)
- Quick Start code snippet (same as `lQxT5` Quick Start section)
- Webhook URL configuration (optional)
- Right: "What's included" checklist (REST API, Webhooks, SDK, SSO)
- CTA: "Next: Brand →"

#### OB-W4: Wizard Step 3 — Brand `5nCtO` (keep, minor adjustments)

**Current**: Logo, Brand Colors, Typography, Button Style, mini Live Preview
**Changes**:
- Keep all form fields as-is
- **Improve Live Preview**: Instead of small card preview, show a larger preview of:
  - For Embed path: a widget (e.g., Leaderboard) with applied branding
  - For Domain path: the portal header + first section with applied branding
  - For SDK path: a code snippet + "Your users will see:" mini preview
- Add note at bottom: "You can refine branding anytime in Brand Settings"
- CTA: "Next: Preview →"

#### OB-W5: New Wizard Step 4 — Preview & Publish (new page, needs new Node ID)

**Purpose**: Show the project owner what their WL deployment will look like. Verify readiness before publishing.

**Layout**: Preview area (70%) + Readiness panel (30%)

**Left — Deployment Preview** (adapts to path):
- Embed path: Mock website frame with embedded widgets visible. Shows "Your website" header + TaskOn widgets embedded below.
- Domain path: Full community portal preview at "community.yourproject.io"
- SDK path: "SDK integration preview not available — see documentation for implementation guide" + link to SDK docs. Show API health check: "API Key active ✅"

**Right — Readiness Checklist**:
- ✅ Community: [name] with X tasks, Y rewards
- ✅ Path: Embed in Your App
- ✅ Widgets: 3 selected (Leaderboard, Task List, User Center)
- ✅ Brand: Logo + colors configured
- ⚠ Widget configuration: pending (configure in Widget Library after publish)
- ⚠ Page Builder: no pages yet (optional, build after publish)
- Overall: "Ready to publish — configure widgets in the next step"

**Bottom actions**:
- Primary: **"Publish White Label"** (purple button)
  - Creates WL instance in Draft state
  - Post-click: success screen ("White Label is ready! Configure your widgets to start embedding.") → redirect to Active state
- Secondary: "Save Draft"

#### OB-W6: Redesign Active State `BnkYW` — Checklist rewrite

**Current checklist** (6 steps):
1. ✅ Create your community
2. ✅ Choose deployment path
3. Send Dev Kit to developer (expanded, step 3)
4. Integration verified
5. Announce to your users
6. First user interaction

**Problems**: (1) Step 1 "Create your community" is Community, not WL. (2) "Send Dev Kit" comes before widget/page configuration. (3) No steps for widget config, page building, or preview. (4) Entire content creation phase missing.

**New checklist** (path-adaptive, showing Embed path as example):

```
Getting Started                                          3 of 9 complete

COMPLETED BY WIZARD
✅ Community ready: Arbitrum Builders (12 tasks, 3 rewards)
✅ Deployment path: Embed in Your App
✅ Brand configured (logo + colors)

CONFIGURE YOUR TOOLS
○ Set up Leaderboard widget                              → Widget Library
  Name it, choose point type, get embed code. 2 min.

○ Set up Task List widget                                → Widget Library
  Shows your Community tasks to users. 2 min.

○ Set up User Center widget                              → Widget Library
  Points balance + level progress for users. 2 min.

○ Build your first page (optional)                       → Page Builder
  Combine widgets into a single embeddable page.

DEPLOY & LAUNCH
○ Preview as a user                                      → Open Preview
  See the full experience with your branding.

○ Send Dev Kit to developer                              → Generate Dev Kit
  One link with embed codes, SSO setup, and docs.
  Everything your developer needs to integrate.

○ Integration verified                                   auto-detect
  Waiting for first API call from your domain.

○ Announce to your users                                 → Promo Kit
  AI-generated social posts + branded banner.

○ First user interaction                                 auto-detect
  Waiting for first widget interaction.
```

- "CONFIGURE YOUR TOOLS" section is **dynamic** based on path:
  - Embed path: widget setup steps (one per selected widget) + optional PB step
  - Domain path: "Verify DNS" + "Customize portal" steps
  - SDK path: "Implement SDK" + "Test API integration" steps
- Steps auto-complete when widget is configured / page is published / etc.
- Dev Kit step moves to DEPLOY section (after configuration, not before)

**Stats row** (below checklist):
- Custom Domain / Widget Status / Pages Published / Monthly Impressions
- Keep as current but values update as user completes steps

**Toolkit section**:
- Rename from "Your Toolkit" to "Your Tools"
- Only show tools relevant to chosen path (not all 6)
- Each card shows configuration status: "3 widgets configured" / "0 pages built" / etc.

**Resources**: Keep as-is

#### OB-W7: Eliminate Duplicate Pages

**Embed Options `RgCVQ`** — duplicates Wizard Step 1 (deployment method choice):
- **Action**: Repurpose as "Deployment Settings" (post-setup management page)
- New content: "Current deployment: Embed in Your App" + ability to add additional deployment methods + comparison reference
- OR: Delete page if Wizard + Active State cover the same ground

**Domain Setup `5bmH9`** — overlaps with Wizard Step 2 (domain path):
- **Action**: Keep as standalone tool page (it's more complete than Wizard Step 2 version)
- Wizard Step 2 domain config becomes a simplified version that links to `5bmH9` for advanced settings
- Remove duplicate DNS form from Wizard — instead embed `5bmH9` content or link to it

**Brand Settings `Cx3LH`** — overlaps with Wizard Step 3:
- **Action**: Keep as the full brand management page. Wizard Step 3 is the "quick setup" version.
- Add note on Wizard: "Advanced options available in Brand Settings"
- No changes needed to `Cx3LH` itself

---

### OB Summary — New/Modified Pages

| # | Page | Node ID | Action | Priority |
|---|------|---------|--------|----------|
| OB-C1 | Community Empty | `zzZ8D` | Redesign: selectable strategy cards + dynamic CTA | P0 |
| OB-C2 | Comm Wizard Step 1 | `Gzpeu` | Update: stepper (3→4), show strategy name | P0 |
| OB-C3 | Comm Wizard Step 2 | `8NeyG` | Redesign: module info + C-end previews + strategy pre-selection | P0 |
| OB-C4 | Comm Wizard Step 3 | `qknQZ` | **Full redesign**: Quick Setup with template content | P0 |
| OB-C5 | Comm Wizard Step 4 | **NEW** | **New page**: Preview & Publish with C-end preview | P0 |
| OB-C6 | Comm Guided | `S1EIA` | Redesign: title/badge + module-level checklist | P0 |
| OB-C7 | Comm Active | `vFRHi` | Redesign: remove checklist, data-driven layout | P1 |
| OB-C8 | Stepper update | all wizard | Update stepper labels on 4 pages | P0 |
| OB-W1 | WL Empty | `Ir6Tq` | Redesign: blocking prereq + selectable paths + single CTA | P0 |
| OB-W2 | WL Wizard Step 1 | `NNwid` | Redesign: 3 paths (not 4) + new stepper | P0 |
| OB-W3 | WL Wizard Step 2 | `CXzmy` | Redesign: path-adaptive content (3 variants) | P0 |
| OB-W4 | WL Wizard Step 3 | `5nCtO` | Update: improved preview | P1 |
| OB-W5 | WL Wizard Step 4 | **NEW** | **New page**: Preview & Publish | P0 |
| OB-W6 | WL Active | `BnkYW` | Redesign: checklist rewrite + path-adaptive steps | P0 |
| OB-W7a | Embed Options | `RgCVQ` | Repurpose or delete (duplicates Wizard Step 1) | P1 |
| OB-W7b | Domain Setup | `5bmH9` | Keep as tool page, dedup from Wizard | P2 |

**Total: 2 new pages + 12 redesigned pages + 2 cleanup items**

---

### Execution Order

| Phase | Tasks | Dependency |
|-------|-------|------------|
| **Phase 1**: Community Wizard | OB-C2, OB-C3, OB-C4, OB-C5, OB-C8 | None (start here) |
| **Phase 2**: Community States | OB-C1, OB-C6, OB-C7 | Phase 1 (need to know wizard output) |
| **Phase 3**: WL Wizard | OB-W2, OB-W3, OB-W5 | Phase 1 (Community must work first) |
| **Phase 4**: WL States | OB-W1, OB-W6, OB-W4 | Phase 3 |
| **Phase 5**: Cleanup | OB-W7a, OB-W7b | Phase 4 |

---

## P3 — Remaining Design Work (Lower Priority)

| # | Page | Code | Status | Notes |
|---|------|------|--------|-------|
| P3-01 | Dev Kit Page | B48 | DONE | `3jDeL` — standalone dev integration page with 3-step guide + verify CTA |
| P3-02 | Unified Pricing — Community tab | M07 | DONE | `AJ3T6` — green theme, $600/mo, Community features |
| P3-03 | Unified Pricing — WL tab | M07 | DONE | `P2TZl` — purple theme, $1,500/mo, WL features + ROI Calculator |
| P3-04 | Pricing tab switching polish | M07 | DONE | Each tab: product-colored buttons/icons, billing cycles, feature lists, CTA copy |
| P3-05 | WL Contract Registry sidebar | B51 | DONE | `OKEqS` missing WL sub-menu. Copy from `4aAo7`. |

---

## P4 — Modal / Dialog Designs (19 modals)

> All B-End management pages have buttons or row clicks that open `(action/modal)` dialogs.
> These modals need concrete UI designs for frontend handoff.
> Code series: **D01-D19**. Each modal is associated with a parent page.
> Create and Edit share the same modal layout (empty vs pre-filled).
> 3 trivial modals excluded (Increase Budget B39, Change Password B47, Notification Panel) — use standard system patterns.

### Community Module Modals (D01-D08)

> Parent pages: B31a-B31h. Each module management page has "Create/Add" button + table row click → modal.

| Code | Modal | Parent | Complexity | Key Fields |
|------|-------|--------|------------|------------|
| **D01** | Points & Level Editor | B31a `zCfKQ` | Medium | Level name, XP threshold, badge icon, perks list, status |
| **D02** | TaskChain Editor | B31b `lpdtp` | High | Chain name, sequential steps (linked tasks), completion reward, reset behavior |
| **D03** | DayChain Config | B31c `fLLVb` | Low-Med | Daily reward amount, bonus day schedule (Day 7/14/30), streak break penalty |
| **D04** | Leaderboard Config | B31d `Emmab` | Medium | Point type selector (EXP/GEM/custom), reset period (weekly/monthly/all-time), display count, visibility |
| **D05** | LB Sprint Editor | B31e `FO9JR` | High | Sprint name, point type, start/end dates, reward tiers (Top 1/Top 10/Top 50 with NFT/Token/WL Spot prizes), status |
| **D06** | Milestone Editor | B31f `WFdZQ` | Medium | Milestone name, condition type (points/tasks/level/streak), threshold value, reward type & amount |
| **D07** | Shop Item Editor | B31g `7yPWx` | Medium | Item name, image upload, description, point cost, stock limit, item type (NFT/Voucher/Merch/WL), redemption method |
| **D08** | Lucky Wheel Config | B31h `sme5a` | High | Segment list (name + probability + prize), visual wheel preview, spin cost, daily spin limit, animation settings |

### Community Settings Modals (D09-D11)

| Code | Modal | Parent | Complexity | Key Fields |
|------|-------|--------|------------|------------|
| **D09** | Badge Editor | B31i `BJLsz` | Medium | Badge name, icon, color, category (Achievement/Engagement/Special), earn condition, description |
| **D10** | Access Rule Editor | B49 `g1CNC` | Medium | Rule name, condition type (Token Gate/NFT Hold/Level Req/Invite Count), chain selector, contract address, threshold, affected areas |
| **D11** | Homepage Section Editor | B50 `5Wm6B` | Medium | Section name, type (Banner/Widget/Custom HTML), content config (depends on type), visibility toggle, sort order |

### WL Advanced Modals (D12-D15)

| Code | Modal | Parent | Complexity | Key Fields |
|------|-------|--------|------------|------------|
| **D12** | Contract Register Form | B51 `OKEqS` | Medium | Chain selector, contract address, ABI upload/paste, contract type (ERC-20/ERC-721/ERC-1155/Custom), verification status |
| **D13** | Activity Rule Editor | B52 `4aAo7` | High | Visual IF→THEN builder: Trigger (event type + conditions) → Action (award points/badge/access) → Frequency (once/daily/unlimited), rule preview |
| **D14** | Privilege Tier Editor | B53 `5xwYN` | Medium | Tier name, rank order, qualification conditions (token gate/level/points/manual), privilege list checkboxes, icon/color |
| **D15** | Privilege Members Panel | B53 `5xwYN` | Low-Med | Member list (address, join date, status), search/filter, manual add/remove, bulk actions |

### Content & Analytics Modals (D16-D18)

| Code | Modal | Parent | Complexity | Key Fields |
|------|-------|--------|------------|------------|
| **D16** | Announcement Editor | B32 `lhR14` | Low-Med | Title, body (rich text), pin toggle, schedule (now/later), target audience (all/tier/segment) |
| **D17** | Featured Slot Editor | B32 `lhR14` | Low | Select item (quest/task/reward), display position (1-4), duration (permanent/timed), image override |
| **D18** | Segment Detail Panel | B54 `olPfE` | Medium | Segment name & criteria, member count, user list (address/level/points/last active), export button, create campaign from segment |

### Cross-Product Modal (D19)

| Code | Modal | Parent | Complexity | Key Fields |
|------|-------|--------|------------|------------|
| **D19** | Promo Kit Generator | B10 `S1EIA` / B15 `BnkYW` | Medium | AI-generated social posts (Twitter/Discord/Telegram), branded banner preview, community/WL link, copy/download buttons, regenerate option |

### Summary

| Group | Codes | Count | Status |
|-------|-------|-------|--------|
| Community Module Modals | D01-D08 | 8 | DONE |
| Community Settings Modals | D09-D11 | 3 | DONE |
| WL Advanced Modals | D12-D15 | 4 | DONE |
| Content & Analytics Modals | D16-D18 | 3 | DONE |
| Cross-Product Modal | D19 | 1 | DONE |
| **Total** | **D01-D19** | **19** | **DONE** |

### Design Approach

**Shared modal shell**: All modals use the same outer frame — 640px or 800px wide overlay, title bar with close X, scrollable body, sticky footer with Cancel + Save/Create buttons. Design the shell once, then populate per-modal content.

**Complexity tiers**:
- **Low** (D03, D15, D17): Simple form fields, 1 column, ≤ 6 fields. Estimate: 1 modal per batch.
- **Medium** (D01, D04, D06, D07, D09-D12, D14, D16, D18, D19): Form + preview or multi-section, 1-2 columns. Estimate: 1 modal per batch.
- **High** (D02, D05, D08, D13): Complex UI — chain step builder, date+tier config, wheel segment editor, visual rule builder. Estimate: dedicated batch each.

**Suggested execution order** (by dependency + reuse):
1. D01 (Points) — establishes modal pattern for all module editors
2. D06 (Milestone), D04 (Leaderboard), D07 (Shop Item) — medium complexity, reuse D01 pattern
3. D05 (LB Sprint), D02 (TaskChain) — high complexity, unique layouts
4. D03 (DayChain), D08 (Wheel), D09 (Badge) — remaining modules
5. D10-D11, D12-D15 — settings + WL modals
6. D16-D19 — content, analytics, promo kit

---

## Archive — Previous Review Items (2026-03-05/06, 40/40 complete)

### Completed Batch Summary

| Category | Items | Status |
|----------|-------|--------|
| P0 Blocking (DR-01/02) | 2 | Done — Shop cards fixed, Integration Center clarified |
| P1a IA/Flow (DR-03–07) | 5 | Done — Checklist, Quests, Home zones, nav overflow |
| P1b Visual (DR-08–12) | 5 | Done — Community hero, module cards, milestones, PB canvas, Embed banner |
| P2 Polish (DR-13–20) | 6+2 deferred | Done — Launch banner, Invite, DayChain, icons, pricing |
| S-Tier Systemic (S1–S7) | 7 | Done — Lifecycle selector, growth stacks, logo wall, ROI, For Projects, Book Demo, Pricing badges |
| I-Tier Data Loop (I-01–I-11) | 11 | Done — Community Insights, Points Economy, streak/funnel charts, segments, AI insights, event markers |
| W-Tier WL Features (W1–W2) | 3 | Done — Contract Registry, Rule Builder, Privilege Manager |

### Items Superseded by OB Restructuring

The following items from the previous review are now **absorbed into** the OB restructuring above:

- **DR-03** (Active checklist too prominent) → absorbed by OB-C7
- **DR-04** (Guided vs Active indistinguishable) → absorbed by OB-C6 + OB-C7
- **DR-12** (Embed Options duplicates Wizard) → absorbed by OB-W7a
- **DR-13** (Wizard Step 3 sparse) → absorbed by OB-C4 (complete replacement)
- **DR-18** (Dev Kit expanded view dense) → absorbed by OB-W6 (Dev Kit repositioned)
- **DR-20** (Checklist inconsistency Guided/Active) → absorbed by OB-C6 + OB-C7

### Items NOT Affected by OB (still valid, completed)

- DR-05 (C-End Quests whitespace), DR-06 (C-End Home zones), DR-07 (C-End nav overflow)
- DR-08 (Community Marketing hero), DR-09 (Module page visual variety), DR-10 (Milestone lock states)
- DR-11 (PB canvas ratio), DR-14–DR-17 (various polish)
- All S-Tier, I-Tier, W-Tier items

---

## Reference — Dependency Chain Diagrams

### Community Content Flow
```
Empty State (choose strategy)
    ↓
Wizard Step 1: Customize (name, brand)
    ↓
Wizard Step 2: Modules (strategy → auto-select, understand each module)
    ↓
Wizard Step 3: Quick Setup (template content auto-generated, inline edit)
    ↓
Wizard Step 4: Preview & Publish (C-end preview + readiness check → Draft)
    ↓
Guided State (module-level content enrichment checklist)
    ↓ (all key steps done)
Active State (data-driven monitoring)
    ↓ (scale + time)
Deep State (AI insights + advanced analytics)
```

### White Label Deployment Flow
```
Prerequisite: Community in Guided or Active state
    ↓
Empty State (choose deployment path)
    ↓
Wizard Step 1: Path (Embed / Domain / SDK — aligned with Empty)
    ↓
Wizard Step 2: Configure (path-adaptive: widgets / DNS / SDK keys)
    ↓
Wizard Step 3: Brand (logo, colors, fonts + improved preview)
    ↓
Wizard Step 4: Preview & Publish (deployment preview + readiness → Draft)
    ↓
Active State (tool configuration checklist → widget setup → page build → Dev Kit → deploy)
    ↓ (integration verified + users active)
Mgmt State (deployments + analytics + toolkit management)
```

### WL Tool Dependency Chain
```
Community Modules (content source)
    ↓ depends on
Widget Library (modules → embeddable components)
    ↓ depends on
Page Builder (components → composed pages)
    ↓ feeds into
Deployment (domain / embed code / SDK)
    ↓ requires
Integration (developer implements in project app)
    ↓ validated by
Verification (auto-detect first API call)
```
