# Frontend Requirement Audit — 6-Dimension Gap Analysis

**Audit Date**: 2026-03-07
**Scope**: `req_community.md` v2.0 + `req_white_label.md` v2.0
**Dimensions**: Page Lifecycle | Button Conditions | Navigation | Data Changes | B→C Cascade | State Completeness

| Priority | Community | White Label | Total |
|----------|-----------|-------------|-------|
| **P0** | 14 | 66 | **80** |
| **P1** | 128 | 124 | **252** |
| **P2** | 12 | 7 | **19** |
| **Total** | **154** | **197** | **351** |

---

# Part 1: Community Product (`req_community.md`)

---

## B09 Community Hub Empty

### Dimension 1: Page Lifecycle
- [P1] Ambiguous: `GET /api/community/status` on mount specified, but whether strategy card data is fetched from API or hardcoded not stated
- [P1] Missing: Destroy — no spec for whether selected strategy card state persists on return

### Dimension 2: Button Conditions
- (Complete — C-01/C-02/C-05 cover selection, disabled state, route guard)

### Dimension 3: Navigation
- (Complete)

### Dimension 4: Data Changes
- [P2] Missing: Strategy card data source — hardcoded or API? If API, cache/refresh specs missing

### Dimension 5: B→C Cascade
- (N/A — pre-creation page)

### Dimension 6: State Completeness
- (Complete — C-04/C-05/§2.3.4 cover loading/guard/error)

---

## B10 Community Hub Guided

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — C-20 specifies onboarding progress API but doesn't list APIs for ACTIVE MODULES and ADD MORE MODULES sections
- [P1] Missing: Destroy — WebSocket `/ws/community/participants` cleanup on page leave not specified (only reconnect logic in C-15)

### Dimension 2: Button Conditions
- [P0] Missing: "Browse Configuration Templates" link (C-19) — visibility condition unclear; navigating to wizard for an already-created community is confusing
- [P1] Missing: Module cards "Configure" vs "Manage" label — determination logic (module status?) not specified

### Dimension 3: Navigation
- [P1] Ambiguous: C-19 "Browse Configuration Templates →" goes to `B13 /community/create?template=browse` — but user already has a community. Creates second or edits existing?

### Dimension 4: Data Changes
- [P0] Missing: After C-14 (enable module), which data blocks refresh? Does checklist (C-16/C-17) auto-update? Does ADD MORE MODULES section remove the enabled module?
- [P1] Missing: Data source for module cards — part of onboarding progress API or separate `GET /api/community/modules`?

### Dimension 5: B→C Cascade
- [P1] Missing: When module enabled via C-14 "+ Enable" — does C-end tab appear immediately or only after configuration?

### Dimension 6: State Completeness
- [P1] Missing: WebSocket initial connection timeout for "First 10 participants" — C-15 mentions polling fallback after 3 retries but no initial timeout

---

## B11 Community Hub Active

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — API call order (C-26 stats with 60s cache, C-24 per-module stats) not specified as parallel/serial
- [P1] Missing: Route guard for community deleted/paused by another admin
- [P2] Missing: Destroy — no cleanup specs

### Dimension 2: Button Conditions
- [P1] Missing: Checklist Banner (C-22/C-23) — after dismiss, does it reappear if all onboarding steps are later completed?
- [P1] Missing: "Add More Modules" section — disappears when all 9 modules enabled?

### Dimension 3: Navigation
- [P1] Missing: Quick Actions "Create Task" → B31 — auto-open task creation modal or just navigate?

### Dimension 4: Data Changes
- [P1] Missing: After enabling module via "+ Enable" — does Module Performance section auto-refresh?
- [P1] Missing: Stats row trend comparison period — "vs last 7 days" configurable or fixed?

### Dimension 5: B→C Cascade
- (Covered by Appendix B — module enable/disable)

### Dimension 6: State Completeness
- [P1] Missing: Module enabled but zero data (e.g., Leaderboard with 0 participants) — per-module card "no data yet" state

---

## B12 Community Hub Deep

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — multiple data sources (stats, AI insights, integrations, analytics chart) — init order and parallel/serial not specified
- [P1] Missing: Pre-API page state — skeleton for all sections?
- [P2] Missing: Destroy — no cleanup specs

### Dimension 2: Button Conditions
- [P1] Missing: AI Insights "Dismiss All" — visible only when 2+ insights?
- [P1] Missing: Integration cards "Reconnect" button styling for error state

### Dimension 3: Navigation
- (Complete — C-29/C-30/C-31/C-33 cover all forward nav)

### Dimension 4: Data Changes
- [P1] Missing: After dismissing AI insight (C-27) — does section collapse when all dismissed?
- [P1] Missing: WAU chart data source — same as `/api/community/stats` or separate analytics API? Cache duration?
- [P1] Missing: Integrations data source — separate API or bundled?

### Dimension 5: B→C Cascade
- (N/A — read-only dashboard)

### Dimension 6: State Completeness
- [P1] Missing: Integration card "error" state (C-32) — retry/auto-reconnect mechanism or only manual "Fix Connection"?

---

## B13 Wizard Step 1: Customize

### Dimension 1: Page Lifecycle
- [P1] Ambiguous: Init — C-43 `GET /api/community/drafts` + URL `?template={id}` — which takes priority when both exist?

### Dimension 2: Button Conditions
- [P1] Missing: "Cancel" button — logic for determining which Hub page to return to (B09/B10/B11/B12)
- [P1] Missing: "Save Draft" — always enabled or only when form dirty?

### Dimension 3: Navigation
- [P0] Missing: "Cancel"/"Back" target resolution — doc says "→ B09/B10/B11/B12" but no logic specified

### Dimension 4: Data Changes
- (Complete — C-38/C-39/C-40 cover preview and save)

### Dimension 5: B→C Cascade
- (N/A — wizard draft)

### Dimension 6: State Completeness
- [P1] Missing: Draft restore conflicts with URL template param — resolution?
- [P2] Missing: Maximum number of drafts

---

## B34 Wizard Step 2: Configure Modules

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — relies on wizard shared state; what if state lost (browser crash)?
- [P1] Missing: Destroy — does navigating back auto-save module selections?

### Dimension 2: Button Conditions
- (Complete — C-46/C-48 cover required module disabled state)

### Dimension 3: Navigation
- [P1] Missing: Direct URL access to Step 2 without completing Step 1 — redirect behavior?

### Dimension 4: Data Changes
- [P1] Missing: Module toggle — persisted on each change or only on Save Draft/Next?

### Dimension 5: B→C Cascade
- (N/A — wizard draft)

### Dimension 6: State Completeness
- [P1] Missing: Module removed from product (e.g., deprecated) — toggle list dynamic?

---

## B35 Wizard Step 3: Quick Setup

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — auto-generated template contents (starter tasks, level names) — hardcoded or API-driven?
- [P1] Missing: Destroy — back to Step 2, change modules, return to Step 3 — are inline edits preserved or regenerated?

### Dimension 2: Button Conditions
- [P1] Missing: "Save Draft" — does it save inline-edited content to draft API?
- [P1] Missing: "+ Add Task" (C-56) — maximum number of wizard tasks?

### Dimension 3: Navigation
- [P1] Missing: "Edit after setup →" links (C-55) open module mgmt in new tab — but community not published yet. 404?

### Dimension 4: Data Changes
- [P1] Ambiguous: C-52/C-53 inline edits save to wizard local state — also auto-saved to draft API like C-40?

### Dimension 5: B→C Cascade
- (N/A — wizard draft)

### Dimension 6: State Completeness
- [P0] Missing: Step 2 has only required modules (blank template) — what auto-generated content appears for Step 3?

---

## B55 Wizard Step 4: Preview & Publish

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — C-62 `GET /api/community/readiness` called on mount or on Publish click? C-60 slug check — auto or user-triggered?
- [P1] Missing: Destroy — custom slug preserved if user navigates back?

### Dimension 2: Button Conditions
- (Complete — C-63/C-65 cover enabled/disabled and full flow)

### Dimension 3: Navigation
- [P1] Missing: After publish → B10 Guided. But if editing existing community, should go to B11/B12?

### Dimension 4: Data Changes
- [P1] Missing: C-end preview mock data source — wizard state (local) or draft API? Complex mock data (leaderboard entries)?

### Dimension 5: B→C Cascade
- [P0] Missing: After Publish — exact C-end pages created/activated? "Initially minimal" per b2c_operation_mapping.md §3.1 is undefined

### Dimension 6: State Completeness
- [P1] Missing: Readiness check finds issues mid-wizard (e.g., no tasks) — can user see preview but not publish?

---

## B31 Sectors & Tasks

### Dimension 1: Page Lifecycle
- [P0] Missing: Init — no API endpoint for sector+task tree. Is it `GET /api/community/sectors?include=tasks`?
- [P1] Missing: Destroy — unsaved inline edits (contenteditable) — warn on navigate?
- [P1] Missing: URL params like `?sector=xxx` for auto-expand?

### Dimension 2: Button Conditions
- [P0] Missing: "Publish Task" (C-89) — visible only on Draft tasks, or also Paused/Hidden? Batch publish?
- [P1] Missing: Task visibility toggle (C-85 eye icon) — "hidden" = `status=hidden` or `visible=false`? Different lifecycle states per §10.2

### Dimension 3: Navigation
- [P0] Missing: No D-code for task creation modal (C-88) — fields, validation, behavior unspecified
- [P1] Missing: Task edit modal (C-85 pencil) — what modal? No D-code assigned

### Dimension 4: Data Changes
- [P1] Missing: After drag-reorder (C-80/C-81/C-82) — stats row refresh?
- [P1] Missing: After task publish (C-89 post-D20) — task badge animate to Active? Stats "Active Tasks" increment?

### Dimension 5: B→C Cascade
- [P0] Missing: C-90 says "soft delete (archive)" but b2c_operation_mapping.md §4.2.8 says "permanently removed from C-end" — inconsistency
- [P1] Missing: Task status Active→Hidden via eye icon — C-end impact timing (immediate or next load)?

### Dimension 6: State Completeness
- [P1] Missing: Sector with all tasks in Draft — C-end visibility (auto-hidden per b2c_operation_mapping.md §4.1.8)?

---

## B31a Points & Level

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — Stats Row + Point Types (C-91) + Level table — separate APIs or one? Init order?
- [P2] Missing: Destroy — no cleanup needed

### Dimension 2: Button Conditions
- [P1] Missing: "+ Add Type" (C-91) — max point types? Behavior at limit?
- [P1] Missing: "+ Add Level" — disabled when no point type selected?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D01 modal)

### Dimension 4: Data Changes
- [P0] Missing: C-94 (modify threshold with users) — level table refresh? "Members" count update for demotions?
- [P1] Missing: C-91 (add point type) — cross-page awareness for B31d/B31e point type dropdowns?

### Dimension 5: B→C Cascade
- [P0] Missing: C-94 threshold change causing demotions — C-end user notification mechanism? Level badge real-time update?
- [P1] Missing: C-91 new point type — appears in C-end User Status Bar immediately or only after tasks configured?

### Dimension 6: State Completeness
- [P1] Missing: 0 levels configured — Stats Row "Level-Up Events" shows what? Empty level table state?
- [P1] Missing: C-95 delete lowest level — "No Level" status display in C-end?

---

## B31b TaskChain

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — chains list and funnel chart data — same endpoint or separate?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Activate Chain" — visible only for Draft? Reactivating Paused?
- [P1] Missing: Funnel chart — shown for active chain only or all? Multiple chains which one?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D02)

### Dimension 4: Data Changes
- [P1] Missing: After activating chain (C-98) — funnel chart appear/update? Stats "Active Chains" increment?

### Dimension 5: B→C Cascade
- [P0] Missing: C-97 pause chain — "progress frozen" but C-end UI unspecified. b2c_operation_mapping.md §5.4.6 mentions "Chain paused" message but exact UI (banner? overlay?) not in this doc

### Dimension 6: State Completeness
- [P1] Missing: Chain references deleted/archived task — Appendix B.2 mentions warning but display location unspecified

---

## B31c DayChain

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — streak distribution chart data — separate from chains list?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- (Complete — C-100/C-101 cover chart interaction and status)

### Dimension 3: Navigation
- (Complete — C-100 covers D03 with auto-focus)

### Dimension 4: Data Changes
- [P1] Missing: After modifying DayChain config via D03 — streak distribution chart refresh?

### Dimension 5: B→C Cascade
- [P1] Missing: DayChain config changes (reward, milestone, grace period) — C-end propagation timing (immediate or next day)?

### Dimension 6: State Completeness
- [P1] Missing: 3+ DayChains in table — standard pagination applies?

---

## B31d Leaderboard

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — leaderboard instances list + stats — same endpoint or separate?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Archive" vs "Delete" relationship (C-104 + C-69) — can delete after archive? Unarchive?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D04)

### Dimension 4: Data Changes
- [P1] Missing: After archive (C-104) — stats "Participation Rate" recalculate excluding archived?

### Dimension 5: B→C Cascade
- (Complete — C-104 explicitly covers C-end tab hidden, historical data preserved)

### Dimension 6: State Completeness
- (Complete — C-105 covers max limit)

---

## B31e LB Sprint

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — active sprint countdown timer — client-side or server-synced? Clock drift?
- [P2] Missing: Destroy — timer cleanup

### Dimension 2: Button Conditions
- [P0] Missing: "Launch Sprint" — visible only for Draft? Can Scheduled sprints be launched early?
- [P1] Missing: "End Early" (C-106) — disabled during reward distribution?
- [P1] Missing: "View Results" (C-107) — "Retry" button role restrictions?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D05)

### Dimension 4: Data Changes
- [P1] Missing: C-106 End Early — sprint immediately "Completed"? Stats update?
- [P1] Missing: C-108 reward distribution Auto mode — "Distributing..." intermediate state in table?

### Dimension 5: B→C Cascade
- [P0] Missing: C-108 Token/NFT reward pre-charge flow — where does admin deposit tokens? B31e or separate page?
- [P1] Missing: After sprint ends — C-end "LB Sprint" tab shows final standings or disappears? Duration?

### Dimension 6: State Completeness
- [P1] Missing: Concurrent sprints — multiple active simultaneously? Stats aggregation?
- [P1] Missing: Sprint with 0 participants at end — reward distribution behavior?

---

## B31f Milestones

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — milestones list + stats API endpoint
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Activate Milestone" — location? Table row or header?

### Dimension 3: Navigation
- (Complete — C-110/C-111 cover expansion and detail)

### Dimension 4: Data Changes
- [P1] Missing: C-112 modify threshold — Tier detail panel (C-110) immediate reflection?

### Dimension 5: B→C Cascade
- [P1] Missing: C-112 threshold change — users past old threshold but unclaimed — still "claimable"? Doc says "retain eligibility" but C-end visual state unclear

### Dimension 6: State Completeness
- [P1] Missing: All tiers claimed by all users — "Completed" state/badge?

---

## B31g Benefits Shop

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — shop items + stats API. Filters "Price Range"/"Search categories" — client-side or API params?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Replenish" (C-113) — visible for Sold Out only or all Limited stock?
- [P1] Missing: "Publish Item" — visible for Draft only?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D07)

### Dimension 4: Data Changes
- [P1] Missing: C-113 replenish — stats "Items Sold Out" decrement? Insight Banner update?

### Dimension 5: B→C Cascade
- (Complete — C-113/C-116 cover availability and auto-pause)

### Dimension 6: State Completeness
- [P1] Missing: Item with Level/Badge gate (C-114) — referenced badge deleted → gate behavior?

---

## B31h Lucky Wheel

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — wheel instances + stats API
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Activate Wheel" (C-120) — what conditions besides "all Nothing prizes" block activation?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D08)

### Dimension 4: Data Changes
- [P1] Missing: After activation — stats update?

### Dimension 5: B→C Cascade
- [P1] Missing: Prize modified while wheel active — C-end reflects new prizes on next spin or new sessions?

### Dimension 6: State Completeness
- [P1] Missing: Token/NFT prizes — pre-charge requirement like LB Sprint? Fulfillment flow?

---

## B31i Badges

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — badges list + stats API
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Auto-trigger" badge — manual check trigger available?

### Dimension 3: Navigation
- (Complete — C-121 covers Holders → D18, standard module mgmt + D09)

### Dimension 4: Data Changes
- [P1] Missing: C-123 manual badge grant — table "Earned" column immediate increment?

### Dimension 5: B→C Cascade
- [P1] Missing: C-123 manual badge grant — C-end notification format (toast? push? in-app?)
- [P1] Missing: Badge deletion — cascade to Shop gates (C-114) and Milestone rewards (M-31)?

### Dimension 6: State Completeness
- [P1] Missing: Auto-condition "Complete X tasks" — tasks deleted/archived → condition unattainable? Warning?

---

## B49 Access Rules

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — rules list + stats API
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: Rule toggle enable/disable — C-end users currently inside community when Token Gate enabled — kicked out?

### Dimension 3: Navigation
- (Complete — standard module mgmt + D10)

### Dimension 4: Data Changes
- [P1] Missing: C-124 rule priority reorder — "Preview Rule" (C-125) auto-refresh?

### Dimension 5: B→C Cascade
- [P0] Missing: C-126 Token Gate without blockchain integration — rule "Inactive"; after B61 connects blockchain, does rule auto-activate?
- [P1] Missing: New access rule enabled — existing C-end users who don't meet criteria — locked out or only new visitors?

### Dimension 6: State Completeness
- [P1] Missing: Multiple rules interaction — AND or OR logic? C-125 implies sequential but boolean logic unspecified

---

## B50 Homepage Editor

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — sections list + stats API
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- (Complete — C-127/C-128/C-129)

### Dimension 3: Navigation
- (Complete — C-129 preview in new tab)

### Dimension 4: Data Changes
- [P1] Missing: C-127/C-128 "即时生效" + "C端下次加载时反映" — is there a Publish step? Could cause confusion

### Dimension 5: B→C Cascade
- [P1] Missing: D11 section references widget (e.g., Leaderboard) later archived/deleted — auto-hide or broken state?

### Dimension 6: State Completeness
- (Complete — C-130/C.3 cover types and limits)

---

## B32 Content Management

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — multiple sections (Announcements, Featured, Module Status) — how many API calls? Parallel/serial?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: Announcement edit/delete — in ⋮ menu or inline buttons?
- [P1] Missing: C-136 says announcements DON'T go through D20, but §7.2 trigger table lists B32 with "Publish (公告/内容)" as D20 trigger — **CONTRADICTION**

### Dimension 3: Navigation
- (Complete — C-133/C-134/C-135)

### Dimension 4: Data Changes
- [P1] Missing: C-131 pin announcement — carousel preview reorder?
- [P1] Missing: C-132 scheduled announcement auto-publish — B-end UI reflect transition? Auto-refresh?

### Dimension 5: B→C Cascade
- [P1] Missing: C-134 remove featured slot — C-end immediate or next load?
- [P1] Missing: Featured slot pointing to deleted Quest/Sprint/Milestone — auto-remove or broken?

### Dimension 6: State Completeness
- [P1] Missing: All 6 featured slots filled — "+ Add Featured" UI hidden?
- [P1] Missing: Module Status Overview — module in error state (depends on deleted resource)?

---

## B33 Preview Mode

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — preview data source: C-end API or B-end preview API?
- [P1] Missing: URL params — `?from=B32` only? What about `?tab=` for specific C-end tab?
- [P2] Missing: Destroy — embedded frame cleanup

### Dimension 2: Button Conditions
- (Complete — C-138/C-140)

### Dimension 3: Navigation
- [P1] Missing: C-139 "Exit Preview" uses `document.referrer` — unreliable (privacy). Should rely on URL param `?from=` only

### Dimension 4: Data Changes
- (N/A — read-only)

### Dimension 5: B→C Cascade
- (N/A — doesn't affect production)

### Dimension 6: State Completeness
- [P1] Missing: Mobile preview (375px) — responsive adaptation or scaled-down? Touch interactions?
- [P1] Missing: Preview with no modules enabled — C-end shows what?

---

## B54 Community Insights

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — 5 APIs with different cache durations — all parallel on mount? Per-section loading?
- [P1] Missing: Pre-API state — skeleton for 4-stat bar + economy chart + segments + retention?
- [P2] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: "Export CSV"/"Export PDF" — disabled when 0 data?
- [P1] Missing: Module Filter (C-142) — selected modules with no data — charts empty or filter shows "No data"?

### Dimension 3: Navigation
- [P1] Missing: URL param sync — C-141 has `?from=&to=`, C-142 module filter doesn't sync to URL
- [P1] Missing: C-147 click retention bar → module mgmt — carries current date range as URL params?

### Dimension 4: Data Changes
- [P1] Missing: Date change (C-141) — different cache durations per API (60s/300s) — cache busted on date change?
- [P1] Missing: Module Filter change (C-142) — refetch or client-side filter?

### Dimension 5: B→C Cascade
- (N/A — read-only analytics)

### Dimension 6: State Completeness
- [P1] Missing: Retention chart — recently disabled module still in chart with historical data?

---

## B61 Community Integration Center

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — API endpoint for integration status list? Loading skeleton for 9 cards?
- [P1] Missing: URL anchor support (C-33 sends `#{integration_id}`) — scroll to anchored card on mount
- [P2] Missing: Destroy — OAuth popup cleanup mid-flow

### Dimension 2: Button Conditions
- [P1] Missing: "Disconnect" (C-152) — always visible or after confirmation?
- [P1] Missing: Integration cards — disabled for higher subscription tier?

### Dimension 3: Navigation
- [P1] Missing: "Back to Community" breadcrumb — which Hub page? Use §10.1 state routing?

### Dimension 4: Data Changes
- [P0] Missing: After connecting integration (e.g., Twitter C-148) — D20 Readiness Check should pass. Event/callback mechanism or D20 re-checks each time?
- [P1] Missing: After disconnect (C-152) — "2 of 9 integrations active" immediate update?

### Dimension 5: B→C Cascade
- [P1] Missing: Blockchain integration (C-151) — immediately enables Token Gate rules in B49? Or manual activation needed?
- [P1] Missing: Twitter disconnect — can user still publish new content?

### Dimension 6: State Completeness
- [P1] Missing: 5 integrations (Wallet Connect, On-Chain Verification, Google Analytics, Webhooks, Data Export) have NO connection flow specs — only Twitter/Discord/Telegram/Blockchain documented

---

## D01 Points & Level Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — API for M-03 point type dropdown? Loading inside modal?
- [P1] Missing: Close on Esc — dirty form confirmation (§7.3) but no explicit D01 spec

### Dimension 2: Button Conditions
- [P1] Missing: "Save"/"Create" — disabled when threshold out of order?
- [P1] Missing: Edit vs Create mode — title/visual differences?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete — M-06)

### Dimension 5: B→C Cascade
- [P1] Missing: New level — immediate C-end user level assignments or only after "published"?

### Dimension 6: State Completeness
- [P1] Missing: M-02 threshold — edit mode pre-fill? Create vs edit distinction?

---

## D02 TaskChain Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — `GET /api/community/tasks?status=active` for M-09. 0 active tasks — can modal open?

### Dimension 2: Button Conditions
- [P1] Missing: "+ Add Step" — disabled at max 20 steps?
- [P1] Missing: "Save" — disabled when < 2 steps?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P1] Missing: After save — B31b table auto-refreshed?

### Dimension 5: B→C Cascade
- (N/A — saves as Draft)

### Dimension 6: State Completeness
- [P0] Missing: M-09 — same task in Step 1 and Step 2 allowed? Duplicate task across steps?

---

## D03 DayChain Config (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — `GET /api/community/tasks?status=active` for M-13

### Dimension 2: Button Conditions
- (Complete — M-15/M-16)

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Editing active DayChain — affects users' current streak? Changing daily task mid-streak?

### Dimension 6: State Completeness
- [P1] Missing: M-15 milestone bonuses — add Day 5 between Day 3 and Day 7? Sort behavior?

---

## D04 Leaderboard Config (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — `GET /api/community/point-types` for M-18

### Dimension 2: Button Conditions
- (Complete)

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Editing active leaderboard Period (Weekly→Monthly) — resets rankings?

### Dimension 6: State Completeness
- (Complete)

---

## D05 LB Sprint Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — `GET /api/community/point-types` for M-23

### Dimension 2: Button Conditions
- [P1] Missing: M-26 reward tier rank overlap — error message format and location?
- [P1] Missing: "Save" — disabled on overlapping ranks?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- (N/A — saves as Draft)

### Dimension 6: State Completeness
- [P1] Missing: M-25 — minimum sprint duration? 1 hour allowed?
- [P1] Missing: Editing Scheduled sprint — all fields changeable before start? Including Start Date?

---

## D06 Milestone Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — M-31 references badges (B31i) and shop items (B31g) — API calls?

### Dimension 2: Button Conditions
- [P1] Missing: M-29 "+ Add Tier" — disabled at max 10?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Adding tier to active milestone — retroactive eligibility for users past threshold?

### Dimension 6: State Completeness
- [P1] Missing: M-31 badge reward — if no badges in B31i, dropdown shows what? Link to create?

---

## D07 Shop Item Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — M-38 references levels (B31a) and badges (B31i) — API calls?

### Dimension 2: Button Conditions
- [P1] Missing: M-39 "Publish Now" — triggers D20 on Save or two-step?
- [P1] Missing: M-35 image upload — drag-drop states (hover, invalid file, uploading progress)

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Active item price edit — C-end immediate? Users with item in "cart"?

### Dimension 6: State Completeness
- (Complete — M-37/M-38 cover conditional fields)

---

## D08 Lucky Wheel Config (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — wheel preview (C-118) loading state?

### Dimension 2: Button Conditions
- (Complete — C-117/C-120)

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete — C-118 real-time preview)

### Dimension 5: B→C Cascade
- (N/A — saves as Draft)

### Dimension 6: State Completeness
- [P1] Missing: M-41 both prizes "Nothing" — saves as Draft OK per C-120, but does activation restriction auto-lift when fixed?

---

## D09 Badge Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: 50+ preset icons — fetched from API or bundled?

### Dimension 2: Button Conditions
- [P1] Missing: M-51 Auto-trigger→Manual switch — auto-condition section (M-52) clears values?
- [P1] Missing: M-54 "Grant Badge" — available in create mode or edit only?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete — C-123)

### Dimension 5: B→C Cascade
- [P1] Missing: Auto-trigger badge condition met — C-end notification mechanism?

### Dimension 6: State Completeness
- [P1] Missing: M-52 "Complete {N} tasks" — all tasks or specific sector? "Earn {N} points" — which point type?

---

## D10 Access Rule Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — M-57 Token Gate needs chain list from B61. 0 chains?

### Dimension 2: Button Conditions
- [P1] Missing: M-56 switch rule type — previous type's params clear?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P0] Missing: New Active access rule — currently logged-in C-end users who fail criteria — kicked out?

### Dimension 6: State Completeness
- [P1] Missing: M-59 Invite Only CSV — duplicates? Invalid addresses? Error reporting per-row or aggregate?

---

## D11 Homepage Section Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — M-64 widget data from leaderboard/sector instances — API calls?

### Dimension 2: Button Conditions
- [P1] Missing: M-66 Custom HTML "sanitized" warning — on input or Save?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Custom HTML section — sanitized HTML live on C-end immediately or after B50 visibility ON?

### Dimension 6: State Completeness
- [P1] Missing: M-65 rich text — dark theme rendering (link color, heading sizes)?

---

## D16 Announcement Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Edit mode — which fields editable after publishing?

### Dimension 2: Button Conditions
- [P1] Missing: M-73 schedule — time passes while modal open?
- [P1] Missing: M-74 pin — warning about unpinning existing pinned announcement?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Scheduled announcement auto-publish — B-end notification that it went live?

### Dimension 6: State Completeness
- (Complete — M-68 through M-74)

---

## D17 Featured Slot Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open — M-76 needs published quests/sprints/milestones — API calls?

### Dimension 2: Button Conditions
- (Complete)

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (Complete)

### Dimension 5: B→C Cascade
- [P1] Missing: Featured slot → sprint ends — auto-remove or "ended" state?

### Dimension 6: State Completeness
- (Complete)

---

## D18 Segment Detail Panel

### Dimension 1: Page Lifecycle
- [P1] Missing: Panel open — API endpoint for user list? Pagination spec?
- [P1] Missing: Loading skeleton for user table

### Dimension 2: Button Conditions
- [P1] Missing: M-84 Export CSV — disabled at 0 users?

### Dimension 3: Navigation
- (N/A — side panel)

### Dimension 4: Data Changes
- [P1] Missing: M-82/M-83 search/filter — client-side or server-side? Debounce?

### Dimension 5: B→C Cascade
- (N/A — read-only)

### Dimension 6: State Completeness
- (Complete — C.5 covers empty segment)

---

## D19 Promo Kit Generator (Modal)

### Dimension 1: Page Lifecycle
- [P0] Missing: Modal open — `POST /api/promo-kit/generate` for AI text + AI image. Total generation time? Loading state?
- [P1] Missing: AI generation failure — retry? Manual fallback?

### Dimension 2: Button Conditions
- [P1] Missing: M-87 "Regenerate" — rate limit? Infinite regeneration?
- [P1] Missing: M-86 editable textarea — "Regenerate" overwrites user edits? Warning?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- (N/A — ephemeral content)

### Dimension 5: B→C Cascade
- (N/A — external sharing)

### Dimension 6: State Completeness
- [P1] Missing: M-85 platform switch (Twitter→Discord) — auto-regenerate or keep text? Confirmation if user edited?

---

## D20 Publish Readiness Check (Modal) — Community

### Dimension 1: Page Lifecycle
- [P1] Missing: Which API checks subscription + Twitter auth? Timeout per check or overall (M-91 says 10s)?

### Dimension 2: Button Conditions
- (Complete — M-91/M-92/M-93)

### Dimension 3: Navigation
- (N/A — modal with inline resolution links)

### Dimension 4: Data Changes
- [P1] Missing: M-93 5-min cache — stored where? Two tabs behavior?

### Dimension 5: B→C Cascade
- (Cascade from calling page, not D20)

### Dimension 6: State Completeness
- [P0] Missing: §7.2 trigger table says B32 Content Mgmt triggers D20, but C-136 says announcements DON'T — **CONTRADICTION needs resolution**
- [P1] Missing: M-92 "Publish anyway" — subscription times out but Twitter passes — can still "Publish anyway"?

---

# Part 2: White Label Product (`req_white_label.md`)

---

## B14 White Label Hub Empty

### Dimension 1: Page Lifecycle
- [P0] Missing: No explicit init API. W-05 mentions `GET /api/white-label/status` for route guard, but no spec for other APIs (Community status for banner), loading/skeleton, pre-API state
- [P1] Missing: Destroy — not specified

### Dimension 2: Button Conditions
- [P1] Missing: "Set Up White Label" — disabled state if no subscription? D20 at publish only?
- [P1] Missing: "View Community" banner — API source for determining active Community
- [P2] Missing: Tag buttons (Custom Domain, Widget Library) — disabled if unavailable on plan?

### Dimension 3: Navigation
- [P1] Missing: Browser back from wizard should return to B14
- [P1] Missing: Sidebar highlight for B14 specifically

### Dimension 4: Data Changes
- [P1] Missing: Community Ready Banner data source and refresh
- [P2] Missing: Route guard API cache duration

### Dimension 5: B→C Cascade
- (N/A — pre-configuration)

### Dimension 6: State Completeness
- [P1] Missing: Loading state / skeleton
- [P1] Missing: API failure for route guard — show B14 by default or error?
- [P2] Missing: Subscription limit — plan without WL

---

## B15 White Label Hub Active

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API calls not consolidated. W-06 through W-10 reference various APIs/WebSockets, no unified "on load call these in parallel" with per-block loading
- [P1] Missing: Destroy — WebSocket cleanup (`/ws/wl/integration-ping`, `/ws/wl/first-interaction`) not specified
- [P1] Missing: Checklist state persistence — server-side or localStorage? API?

### Dimension 2: Button Conditions
- [P0] Missing: "Email to Developer" (W-07) — field validation (email format, required), disabled/loading state for send button
- [P1] Missing: Toolkit cards "Configure →" — show/hide conditions
- [P1] Missing: Checklist steps — multiple open simultaneously or accordion single?

### Dimension 3: Navigation
- [P1] Missing: "Open Preview" (B33) — same tab or new tab?

### Dimension 4: Data Changes
- [P0] Missing: Stats Row (Domain/Widgets/Pages/Impressions) — data source API not specified
- [P1] Missing: Stats refresh strategy — polling? Manual? WebSocket?
- [P1] Missing: Checklist step completion — "5 of 9 complete" update mechanism?

### Dimension 5: B→C Cascade
- [P1] Missing: "Announce to users" step — WL-specific Promo Kit C-end differences from Community

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton
- [P1] Missing: All checklist complete — auto-transition to B16?
- [P0] Missing: B15→B16 transition criteria — "5+ tools, high traffic" vague. "High traffic" = Monthly Impressions > what?

---

## B16 White Label Hub Management

### Dimension 1: Page Lifecycle
- [P0] Missing: Init APIs — W-11 `GET /api/white-label/status` + W-12 `GET /api/white-label/deployments` — call order, per-section loading, skeleton not specified
- [P1] Missing: Destroy — not specified

### Dimension 2: Button Conditions
- [P1] Missing: Deployment cards — click area (full card or buttons only?)
- [P1] Missing: "View All Deployments" (W-14) — side panel close, pagination
- [P1] Missing: "View Full Analytics" — disabled if no data?

### Dimension 3: Navigation
- [P0] Missing: Analytics route "→ B45 `/analytics`" — **B45 not defined**. Typo for B43? Blocks routing
- [P1] Missing: Breadcrumb for B16

### Dimension 4: Data Changes
- [P0] Missing: Analytics chart — no API endpoint in operations table
- [P1] Missing: Key Metrics (Avg Time, Bounce Rate, etc.) — data source API not specified
- [P1] Missing: Deployment "Healthy" badge — health determination logic?

### Dimension 5: B→C Cascade
- [P1] Missing: Deployment status should reflect live C-end state

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton for chart and stats
- [P1] Missing: Empty analytics chart (new project, no traffic)
- [P0] Missing: W-13 "Feature Management" card status determination (Active/Not Set/Draft) logic per tool

---

## B37 Wizard Step 1: Choose Deployment Path

### Dimension 1: Page Lifecycle
- [P1] Missing: Static page or loads API? Draft restoration for previously selected path?
- [P1] Missing: Destroy — selection persistence on navigate away

### Dimension 2: Button Conditions
- (W-19 covers disabled Next — adequate)
- [P2] Missing: "Cancel" — confirmation dialog?

### Dimension 3: Navigation
- [P0] Missing: W-18 "Embed → B17" with in-B17 switching to Iframe/PB — transition mechanism ambiguous: full page nav to B58/B59 or in-page swap?
- [P1] Missing: Browser back from Step 2 → Step 1 — wizard state preserved?

### Dimension 4: Data Changes
- (Local state — minimal)

### Dimension 5: B→C Cascade
- (N/A — wizard)

### Dimension 6: State Completeness
- [P2] Missing: Loading state (likely N/A for static)

---

## B17 Wizard Step 2: Widget Configure (Embed)

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — widget list from Community modules but no endpoint specified. W-21 says "列表来源: Community 已配置模块" but no API
- [P1] Missing: Community has 0 configured modules — all "Needs Setup". Can user proceed with 0 widgets?

### Dimension 2: Button Conditions
- [P0] Missing: "Next: Brand" disabled visual when 0 widgets selected (W-21 says "至少选 1 个")
- [P1] Missing: SSO Method (W-22) — not shown in page structure
- [P1] Missing: Target Domain (W-23) — not shown in page structure

### Dimension 3: Navigation
- [P0] Missing: "Set up in Community →" — new tab or same tab? Same tab loses wizard state
- [P1] Missing: Embed Method radio switch to Iframe/PB — transition behavior

### Dimension 4: Data Changes
- [P1] Missing: Right panel "Embed Preview" — real-time update as widgets selected? Update mechanism?

### Dimension 5: B→C Cascade
- (N/A — wizard)

### Dimension 6: State Completeness
- [P1] Missing: Error state if Community modules API fails
- [P1] Missing: Loading state for module list

---

## B57 Wizard Step 2: Domain Config

### Dimension 1: Page Lifecycle
- [P1] Missing: DNS verification auto-start on load or manual only?

### Dimension 2: Button Conditions
- [P0] Missing: "Next: Brand" — DNS verification required before proceeding? Or skip and verify later?
- [P1] Missing: "Verify DNS" (W-27) — disabled during polling? Loading state?

### Dimension 3: Navigation
- (Adequate — Back→B37, Next→B38)

### Dimension 4: Data Changes
- [P1] Missing: DNS verification polling — auto or manual trigger?

### Dimension 5: B→C Cascade
- (N/A — wizard)

### Dimension 6: State Completeness
- [P1] Missing: Verification timeout — DNS propagation takes hours? Max polling duration?

---

## B58 Wizard Step 2: Iframe Config

### Dimension 1: Page Lifecycle
- [P1] Missing: Auto-generated Iframe URL — project_id source at wizard stage (not yet published)?

### Dimension 2: Button Conditions
- [P1] Missing: "Next: Brand" — any required fields? W-28/W-29/W-30/W-31 all optional. Proceed immediately?
- [P2] Missing: "Copy" buttons — clipboard failure fallback

### Dimension 3: Navigation
- [P1] Missing: Back target — B37 or B17? Came from B17's embed selector → back should go to B17

### Dimension 4: Data Changes
- (Minimal)

### Dimension 5: B→C Cascade
- (N/A)

### Dimension 6: State Completeness
- [P2] Missing: No explicit loading/error states

---

## B59 Wizard Step 2: PB Config

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — W-32 references pages list from `GET /api/white-label/pages` but not explicit
- [P1] Missing: Branching logic "has pages" (W-32) vs "no pages" (W-33) — conditional render not clear

### Dimension 2: Button Conditions
- [P0] Missing: "Next: Brand" — W-32 requires ≥1 selection; W-33 template — must select template to proceed?

### Dimension 3: Navigation
- (Same Back target concern as B58)

### Dimension 4: Data Changes
- (Minimal)

### Dimension 5: B→C Cascade
- (N/A)

### Dimension 6: State Completeness
- [P1] Missing: Loading state for pages list
- [P1] Missing: No pages AND no templates — edge case

---

## B60 Wizard Step 2: SDK Config

### Dimension 1: Page Lifecycle
- [P0] Missing: API Key auto-generation (W-35 `pk_live_xxx`) — when? On page load? `POST /api/white-label/sdk/keys`? Failure handling?
- [P1] Missing: Project ID (W-36) — source before WL published?

### Dimension 2: Button Conditions
- [P0] Missing: "Regenerate" (W-35) — in wizard context premature. Confirmation dialog? What does regenerate mean for just-created key?
- [P1] Missing: "Next: Brand" — required fields?

### Dimension 3: Navigation
- (Adequate)

### Dimension 4: Data Changes
- (Minimal)

### Dimension 5: B→C Cascade
- (N/A)

### Dimension 6: State Completeness
- [P1] Missing: Key generation failure error state

---

## B38 Wizard Step 3: Brand Customization

### Dimension 1: Page Lifecycle
- [P1] Missing: Init — load existing brand settings (if B40 configured before)? Or fresh?
- [P1] Missing: Destroy — auto-save draft? Unsaved changes warning?

### Dimension 2: Button Conditions
- [P1] Missing: "Back" — preserves B38 data on round-trip?

### Dimension 3: Navigation
- [P1] Missing: "Back" target "→ B17/B57-B60" — path-dependent. How does system know which Step 2? Wizard state?

### Dimension 4: Data Changes
- [P0] Missing: Brand data save timing — on "Next: Preview"? Auto-save? Not specified

### Dimension 5: B→C Cascade
- (N/A — brand applies on publish)

### Dimension 6: State Completeness
- [P1] Missing: Logo upload loading/error in wizard context
- [P1] Missing: Color picker popover close behavior

---

## B56 Wizard Step 4: Preview & Publish

### Dimension 1: Page Lifecycle
- [P0] Missing: Init — W-40 `GET /api/white-label/preview` + W-41 `GET /api/white-label/readiness` — parallel or serial? Preview slow?
- [P1] Missing: Destroy — draft preserved if leave without publishing?

### Dimension 2: Button Conditions
- [P0] Missing: "Publish White Label" — disabled if readiness has amber items? Or informational only?
- [P1] Missing: Desktop/mobile preview toggle — not in page structure, not in operations

### Dimension 3: Navigation
- (W-42 success → B15 — adequate)

### Dimension 4: Data Changes
- [P1] Missing: After publish — other pages notified? B15 forced refetch?

### Dimension 5: B→C Cascade
- [P0] Missing: What C-end pages/URLs go live after WL publish? Dev Kit generation (W-43) referenced but C-end specifics not in this doc

### Dimension 6: State Completeness
- [P1] Missing: Preview loading state (left panel)
- [P1] Missing: Readiness API failure error state

---

## B18 Domain Setup

### Dimension 1: Page Lifecycle
- [P0] Missing: Init APIs not consolidated. W-44 implies DNS auto-polling but no `GET /api/white-label/domain` for current config
- [P1] Missing: Destroy — DNS polling stop on page leave?

### Dimension 2: Button Conditions
- [P0] Missing: "Save" (W-47 edit mode) — validation, loading, enabled conditions for domain input
- [P1] Missing: W-47 domain change triggers D20 — but D20 checks Twitter auth. Why would domain edit need Twitter? Possible design error

### Dimension 3: Navigation
- [P1] Missing: Breadcrumb "← Back to White Label" — B14/B15/B16 state-based?

### Dimension 4: Data Changes
- [P0] Missing: After DNS verification success — what updates? Checklist auto-update? SSL provisioning (W-46) — refetch or WebSocket?
- [P1] Missing: Brand & Portal section (③) — in page structure but no operations. Same as B40?

### Dimension 5: B→C Cascade
- [P0] Missing: Custom domain verified + SSL active — community accessible at custom URL immediately? Old URL redirect?

### Dimension 6: State Completeness
- [P1] Missing: No domain configured (first visit) — input empty, checklist unchecked
- [P1] Missing: SSL provisioning failure after max retries — user sees what?

---

## B19 Deployment Settings (Embed Options)

### Dimension 1: Page Lifecycle
- [P1] Missing: Static page or loads deployment status?
- [P1] Missing: Route guard — accessible without completing wizard?

### Dimension 2: Button Conditions
- [P1] Missing: Current deployment method highlight/selected state

### Dimension 3: Navigation
- (Complete — three routes to B42/B20-B22/B23-B25)

### Dimension 4: Data Changes
- [P1] Missing: W-51/W-52 GET APIs for widget/page count — loading state for checks

### Dimension 5: B→C Cascade
- (N/A — navigation page)

### Dimension 6: State Completeness
- [P2] Missing: Loading/error states

---

## B42 Iframe Embed

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — W-53 `GET /api/white-label/embed/iframe` — what data returned? URL only or SSO + display + dimensions?
- [P1] Missing: Destroy — unsaved SSO config warning

### Dimension 2: Button Conditions
- [P0] Missing: "Test Connection" (SSO section) — not in operations. Behavior, API, success/failure undefined
- [P1] Missing: "Copy Code" — clipboard failure fallback

### Dimension 3: Navigation
- [P1] Missing: Breadcrumb "← Back" — B19 or B15/B16?

### Dimension 4: Data Changes
- [P0] Missing: Display Mode dropdown change — auto-save? Explicit save? No save button in structure
- [P0] Missing: Width/Height inputs — validation (min/max, numeric), error states

### Dimension 5: B→C Cascade
- [P1] Missing: Config save — embedded iframe in client sites reflect new dimensions/SSO immediately? Re-deploy needed?

### Dimension 6: State Completeness
- [P1] Missing: Loading state
- [P1] Missing: Iframe preview fail (W-57 covers preview only, not page load failure)

---

## B20 Widget Library Empty

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — W-62 `GET /api/community/modules/status` not explicit. Community doesn't exist — route guard?
- [P1] Missing: Route guard — if 1+ widgets exist, redirect to B22

### Dimension 2: Button Conditions
- [P1] Missing: "Add Widget →" — passes moduleType to B21? B20 section doesn't mention this
- [P1] Missing: "Create Your First Widget" CTA — goes to B21 with no module pre-selected?

### Dimension 3: Navigation
- [P1] Missing: Breadcrumb "← Back to Embed Options" → B19 — inconsistent with sidebar "Widgets" active

### Dimension 4: Data Changes
- (Static display of Community module status)

### Dimension 5: B→C Cascade
- (N/A — empty state)

### Dimension 6: State Completeness
- [P1] Missing: Community 0 configured modules — all amber. Special message?
- [P1] Missing: Loading skeleton

---

## B21 Widget Config

### Dimension 1: Page Lifecycle
- [P0] Missing: Init — edit mode (`/white-label/widgets/:id/config`) API `GET /api/white-label/widgets/:id` referenced in W-58 but no explicit flow. Create mode API?
- [P1] Missing: Destroy — unsaved config warning

### Dimension 2: Button Conditions
- [P0] Missing: "Save & Get Embed Code" — enabled conditions. Which fields required? Module Type changeable if pre-selected from B20?
- [P0] Missing: Style Config fields (Primary Color, Border Radius, Padding) — in operations but NOT in page structure. Missing from layout

### Dimension 3: Navigation
- [P0] Missing: After save success — stay on B21 or redirect? No back button in page structure
- [P1] Missing: No breadcrumb in page structure

### Dimension 4: Data Changes
- [P1] Missing: Live Preview update — debounce timing? Which fields trigger update?

### Dimension 5: B→C Cascade
- [P1] Missing: Edit mode save — auto-updates already deployed widgets?

### Dimension 6: State Completeness
- [P1] Missing: Save failure error state
- [P1] Missing: Edit mode loading state

---

## B22 Widget Library Active

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/widgets` implied not explicit
- [P1] Missing: All widgets deleted → redirect to B20 (§D.5 but not B22's section)

### Dimension 2: Button Conditions
- [P0] Missing: "Deploy Widget" → D20 — post-D20 API call to change status? Not specified
- [P0] Missing: "Delete" widget — W-60 "Rule 引用 → 403" — what "Rule"? Inconsistent with §D.4 re Page references
- [P1] Missing: Widget card click → B21 edit?

### Dimension 3: Navigation
- [P1] Missing: "Analytics" per widget → B43 — but B43 is "Page Analytics". Widget analytics view?

### Dimension 4: Data Changes
- [P1] Missing: After delete — "Community Modules" section update to show module available again?

### Dimension 5: B→C Cascade
- [P0] Missing: Widget deployed — embed code starts working immediately? Appears in Page Builder library?

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton for widget cards
- [P1] Missing: Widget list API failure error state

---

## B23 Page Builder Empty

### Dimension 1: Page Lifecycle
- [P1] Missing: Widget availability check for templates?
- [P1] Missing: Route guard — redirect to B25 if pages exist

### Dimension 2: Button Conditions
- [P1] Missing: Template cards — disabled if no widgets configured?
- [P1] Missing: "Create Your First Page" — B24 blank or template selector?

### Dimension 3: Navigation
- (Adequate)

### Dimension 4: Data Changes
- (N/A — static)

### Dimension 5: B→C Cascade
- (N/A — empty)

### Dimension 6: State Completeness
- [P2] Missing: Loading state

---

## B24 Page Builder Editor

### Dimension 1: Page Lifecycle
- [P0] Missing: Init create mode — template data fetch? No API specified
- [P0] Missing: Init edit mode — W-69 `GET /api/white-label/pages/:id` is under B25's operations, B24 has no explicit init
- [P1] Missing: Destroy — auto-save interval? Unsaved changes warning?

### Dimension 2: Button Conditions
- [P0] Missing: "Publish Page" — enabled conditions. Page Name required? At least 1 widget? Slug validated?
- [P0] Missing: "Save Draft" — API endpoint not specified (POST create or PUT update?)
- [P1] Missing: "+ Add Widget Block" (W-64) — 0 available AND 0 configured widgets?

### Dimension 3: Navigation
- [P0] Missing: After "Publish Page" via D20 — stay on B24 or redirect to B25?
- [P1] Missing: After "Save Draft" — stay on B24?
- [P1] Missing: No back button. Exit only via breadcrumb?

### Dimension 4: Data Changes
- [P1] Missing: Canvas-Settings sync (W-68 "双向同步") — real-time or on save? Debounce?

### Dimension 5: B→C Cascade
- [P0] Missing: After publish — C-end URL live immediately? Slug URL (`share.taskon.io/pages/{slug}`) accessible? Appears in domain portal?

### Dimension 6: State Completeness
- [P0] Missing: Page Name validation — required? Min/max? Characters?
- [P1] Missing: Slug uniqueness check failure error state
- [P1] Missing: Editor loading state in edit mode

---

## B25 Page Builder Active

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/pages` implied not explicit
- [P1] Missing: All pages deleted → redirect to B23 (§D.5 but not B25's section)

### Dimension 2: Button Conditions
- [P1] Missing: Page card — full card clickable or only buttons?
- [P1] Missing: "+ Create New Page" — disabled at plan limit (§D.1 but not B25)?
- [P1] Missing: Page embed code "Copy" — only published pages?

### Dimension 3: Navigation
- (Adequate)

### Dimension 4: Data Changes
- [P1] Missing: W-70 page stats refresh on return from B24?

### Dimension 5: B→C Cascade
- (Covered at high level by §8.2)

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton
- [P1] Missing: Page status badges (Published/Draft/Unpublished) visual styles

---

## B40 Brand Settings

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/brand` in §16 but not in page section init
- [P1] Missing: Destroy — unsaved changes warning with dirty form

### Dimension 2: Button Conditions
- [P0] Missing: "Save Changes" — disabled when form clean? Or always enabled?
- [P1] Missing: Role restriction — §D.2 says Admin/Editor, Member cannot. Enforcement: redirect? Read-only?

### Dimension 3: Navigation
- [P1] Missing: Breadcrumb "← Back to White Label" — B15 or B16?

### Dimension 4: Data Changes
- [P0] Missing: Preview panel — referenced in W-73/W-74/W-75 but NOT in page structure. Layout inconsistency

### Dimension 5: B→C Cascade
- (W-76/W-77 adequate — CDN TTL 5min)

### Dimension 6: State Completeness
- [P1] Missing: Loading state
- [P1] Missing: No brand configured yet (first visit) — defaults?
- [P1] Missing: Custom CSS preview panel behavior — page structure doesn't include it

---

## B41 SDK & API

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/sdk` in §16 but not in page section. Data returned?
- [P1] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P0] Missing: "+ Generate Key" — in page structure but not operations. Create new key pair? Differs from "Regenerate"?
- [P1] Missing: Webhook "Edit" — W-79 covers "+ Add" but not editing existing webhooks
- [P1] Missing: Webhook "Delete" — not specified

### Dimension 3: Navigation
- [P1] Missing: "View Docs" target URL

### Dimension 4: Data Changes
- [P0] Missing: After "Regenerate Key" (W-78) — "Quick Start" code block auto-update with new key?
- [P1] Missing: Webhook list refresh after adding — optimistic or refetch?

### Dimension 5: B→C Cascade
- [P0] Missing: Key regeneration (W-78) — deployed widgets/pages using old key break immediately? C-end invalid key handling?

### Dimension 6: State Completeness
- [P1] Missing: Empty webhooks section
- [P1] Missing: Empty API keys (no keys yet) — possible?
- [P1] Missing: Page loading state

---

## B26 WL Integration Center

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/integrations` implied not stated. Data returned?
- [P1] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P0] Missing: "Configure" post-connection — what can be configured? Just disconnect? Or settings?
- [P0] Missing: "Disconnect" — not in operations at all
- [P1] Missing: Blockchain & Wallet connect flows (Multi-Chain, WalletConnect, On-Chain Verification) — not in W-84 to W-88

### Dimension 3: Navigation
- [P0] Missing: All buttons route to "B44 `/white-label/integrations/:type`" — **B44 page never defined**. No structure, no operations. Critical gap

### Dimension 4: Data Changes
- [P1] Missing: After OAuth popup closes (W-84) — parent page refresh (full reload or refetch)?

### Dimension 5: B→C Cascade
- [P1] Missing: Twitter/Discord connected — C-end social features enabled?

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton
- [P1] Missing: All integrations "Available" default state?

---

## B51 Contract Registry

### Dimension 1: Page Lifecycle
- (W-89 specifies init — adequate)
- [P1] Missing: Destroy — side panel "View Events" cleanup

### Dimension 2: Button Conditions
- [P0] Missing: Table row click → D12 edit? Or only "Edit" button?
- [P1] Missing: "View Events" side panel — close mechanism, pagination for 100 events
- [P1] Missing: Filter tabs — URL param sync?

### Dimension 3: Navigation
- [P1] Missing: "Active Rules" stat click → B52 — pre-filter to project's contracts?

### Dimension 4: Data Changes
- [P0] Missing: After contract verified (W-93) — stats row auto-update?
- [P1] Missing: After D12 save — table update mechanism (refetch or optimistic)?

### Dimension 5: B→C Cascade
- [P1] Missing: Contract registration indirect cascade to Rule Builder — cross-referenced in §C.3 but not B51

### Dimension 6: State Completeness
- [P1] Missing: Individual contract verification error — row-level display beyond status badge?

---

## D12 Contract Register Form (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open/close animation (fade/slide), close via ×/outside/Esc

### Dimension 2: Button Conditions
- [P0] Missing: "Register Contract" — disabled until validation? Or validates on click?
- [P1] Missing: Loading state during API call

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P0] Missing: ABI invalid JSON — when validation runs (on paste? blur? submit?)
- [P1] Missing: ABI re-pasted — events checkbox list reset?

### Dimension 5: B→C Cascade
- (N/A)

### Dimension 6: State Completeness
- [P1] Missing: Large ABI parsing loading state
- [P1] Missing: Edit mode — readonly fields (address/network/ABI) visual style

---

## B52 Activity Rule Builder

### Dimension 1: Page Lifecycle
- [P0] Missing: Init APIs — Rules Table + Stats + Presets + Anti-Sybil = multiple APIs. Call order, per-section loading not specified
- [P1] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P0] Missing: Rule toggle (active↔paused) — inline toggle? Visual style? Confirmation required?
- [P0] Missing: Anti-Sybil "Preview" — opens what? Modal? Side panel? API listed but UI unspecified
- [P1] Missing: Anti-Sybil config editing — inline or modal? Edit mechanism for Min Wallet Age, Min Transactions
- [P1] Missing: Rule row click → D13 edit — interaction details

### Dimension 3: Navigation
- (Adequate in button route table)

### Dimension 4: Data Changes
- [P0] Missing: After D13 create/edit rule — table refresh mechanism
- [P0] Missing: After rule toggle — stats row immediate update?
- [P1] Missing: "Triggered X times today" — real-time or page-load only?

### Dimension 5: B→C Cascade
- [P0] Missing: Rule activated (Draft→Active) — processes pending events or only new events from activation?

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton for rules table
- [P1] Missing: Soft-delete recovery UI — "30 天内可恢复" but no restore view

---

## D13 Activity Rule Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open/close animation
- [P1] Missing: Opened from Preset — title change?

### Dimension 2: Button Conditions
- [P0] Missing: "Create Rule"/"Save Rule" — disabled conditions. Required fields?
- [P0] Missing: Creates as Draft requiring separate activation? Or can create as Active?
- [P1] Missing: Save button loading state

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P0] Missing: Event dropdown — API for sectors/tasks list?
- [P0] Missing: Badge selector (grant_badge action) — API for badges?
- [P0] Missing: Tier selector (upgrade_tier action) — API for privilege tiers?

### Dimension 5: B→C Cascade
- (Covered in §C.3)

### Dimension 6: State Completeness
- [P1] Missing: No contracts registered — chain events unavailable visual state
- [P1] Missing: No badges/tiers for action types
- [P1] Missing: Edit mode — can change trigger event type of Active rule?

---

## B53 Privilege Manager

### Dimension 1: Page Lifecycle
- [P0] Missing: Init APIs — `GET /api/wl/privileges` + stats + integration. Call order, per-section loading not specified
- [P1] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P0] Missing: Privilege row toggle (active↔paused) — not described in operations
- [P0] Missing: "Manage Members" — mentioned as "(隐含)" but location unspecified
- [P1] Missing: "Via Benefits Shop →" (Mode B) — Shop module not configured?
- [P1] Missing: Table row click → D14 — interaction details

### Dimension 3: Navigation
- (Adequate in button route table)

### Dimension 4: Data Changes
- [P0] Missing: After D14 save — table refresh mechanism
- [P0] Missing: Integration Status "API Connection: Enabled" — determination logic? Setup from this page?
- [P1] Missing: "Total Value Distributed: $12,840" — calculation API?

### Dimension 5: B→C Cascade
- [P0] Missing: Privilege Tier activated — C-end immediate? Push vs pull propagation?
- [P1] Missing: Mode B auto-creates Shop item — Shop module not enabled? Fail silently? Auto-enable?

### Dimension 6: State Completeness
- [P1] Missing: Loading skeleton
- [P1] Missing: Delete privilege tier — not in operations or button routes. Deletion behavior in §C.3 but no page action

---

## D14 Privilege Tier Editor (Modal)

### Dimension 1: Page Lifecycle
- [P1] Missing: Modal open/close (animation, outside click, Esc)
- [P1] Missing: Opened from "Configure →" preset — pre-selection?

### Dimension 2: Button Conditions
- [P0] Missing: "Create Tier" — disabled conditions. Required fields?
- [P0] Missing: Edit mode — qualification mode locked for Active tiers. Visual treatment?
- [P1] Missing: Rank Order conflict auto-resolution — user notification (toast? inline?)?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P0] Missing: Token Gate condition — Contract Address dropdown from B51. API for contracts? 0 contracts?
- [P1] Missing: Level Req — min level dropdown data source?

### Dimension 5: B→C Cascade
- (Covered in B53/§C.3)

### Dimension 6: State Completeness
- [P1] Missing: Tier Icon selector spec — Material Symbols? Text input? Custom color beyond 4 swatches?
- [P1] Missing: Duration fields — "permanent" vs "time_limited" selection not in modal structure

---

## D15 Privilege Members Panel (Modal)

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/wl/privileges/:id/members` — loading state, page size, empty state

### Dimension 2: Button Conditions
- [P0] Missing: "+ Add Member" — inline input or sub-modal? Wallet validation (0x + 40 hex). Duplicate handling
- [P0] Missing: "Bulk Import" — CSV upload flow. Max file size? Max addresses? Progress indicator?
- [P1] Missing: "× remove" — API endpoint loading state per row
- [P1] Missing: "Export CSV" — async or sync? Loading?

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P1] Missing: Add/remove member — title "24 members" immediate update?

### Dimension 5: B→C Cascade
- [P0] Missing: Manual add member — C-end privilege change immediate? Push notification?
- [P0] Missing: Remove member — privilege revoke immediate? Or after current session?

### Dimension 6: State Completeness
- [P1] Missing: 0 members empty state
- [P1] Missing: Read-only mode for Status-Based/Achievement-Based (auto-managed). "× remove" hidden?

---

## B43 Page Analytics

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/white-label/pages/:id/analytics` — invalid `:id` → 404 handling?
- [P1] Missing: Destroy — no cleanup

### Dimension 2: Button Conditions
- [P1] Missing: Date range "Custom" — min range?
- [P1] Missing: D/W/M toggle (W-102) — selected state visual

### Dimension 3: Navigation
- [P1] Missing: Breadcrumb "← Back to Pages" — always B25 or context-dependent?

### Dimension 4: Data Changes
- [P1] Missing: Date range change — per-section loading or global overlay?
- [P1] Missing: Top Pages row click → Widget Interactions (W-103) — client-side or API re-query?

### Dimension 5: B→C Cascade
- (N/A — read-only)

### Dimension 6: State Completeness
- [P1] Missing: Analytics API failure error state
- [P2] Missing: Unpublished page — show historical data?

---

## B48 Dev Kit Page

### Dimension 1: Page Lifecycle
- [P0] Missing: Init API — `GET /api/devkit/{project_id}` — network timeout handling?
- [P1] Missing: Standalone page (no sidebar) — own layout component?

### Dimension 2: Button Conditions
- [P0] Missing: "Verify Integration" (W-110) — after success, state persistence mechanism (server flag? Cookie?)
- [P1] Missing: Multiple widget accordion — W-109 "同时可多个展开" but no max

### Dimension 3: Navigation
- [P1] Missing: Footer links (Docs/API Reference/Support) — target URLs

### Dimension 4: Data Changes
- [P1] Missing: Project owner updates widget/SSO config after Dev Kit generated — Dev Kit auto-reflects or snapshot?

### Dimension 5: B→C Cascade
- (N/A — developer-facing)

### Dimension 6: State Completeness
- [P1] Missing: Initial loading state
- [P1] Missing: Project with 0 widgets — empty widget section

---

## D19 Promo Kit Generator (Modal) — WL

### Dimension 1: Page Lifecycle
- [P0] Missing: Init — AI generation API endpoint, parameters, loading state

### Dimension 2: Button Conditions
- [P0] Missing: "Regenerate" banner — API, loading, rate limit
- [P1] Missing: "Share on [Platform]" — pre-filled tweet compose URL format?
- [P1] Missing: "Copy Text" — clipboard fallback
- [P1] Missing: "Download Banner" — resolution, file size indication

### Dimension 3: Navigation
- (N/A — modal)

### Dimension 4: Data Changes
- [P1] Missing: Platform switch (Twitter/Discord/Telegram) — auto-regenerate or keep text?

### Dimension 5: B→C Cascade
- (N/A)

### Dimension 6: State Completeness
- [P0] Missing: AI generation failure — error state, retry
- [P1] Missing: No logo/brand configured — banner appearance?

---

## D20 Publish Readiness Check (Modal) — WL

### Dimension 1: Page Lifecycle
- [P0] Missing: Init — API endpoints for checks not specified. `GET /api/white-label/readiness` in §16 but modal spec doesn't reference it
- [P1] Missing: Check execution — parallel/serial? Per-item loading spinners?

### Dimension 2: Button Conditions
- [P0] Missing: "Publish" button — enabled when all checks pass? Or always enabled?
- [P0] Missing: Subscription check fail — "Subscribe →" link to M07? Stay in modal or close?
- [P0] Missing: Twitter check fail — OAuth in popup? Within modal or separate?

### Dimension 3: Navigation
- [P1] Missing: Modal close — clicking ×/outside → "Cancel publish?" confirmation?

### Dimension 4: Data Changes
- [P1] Missing: After D20 pass + publish success — modal auto-close or success state first?

### Dimension 5: B→C Cascade
- (Gate, not action — cascade from publish)

### Dimension 6: State Completeness
- [P0] Missing: D20 has no standalone spec in WL document — referenced by 5 pages but undefined. Must include full spec or cross-reference
- [P1] Missing: Check timeout handling (10s from audit_todo.md M-91)
- [P1] Missing: Cache behavior (5-min from audit_todo.md M-93)

---

# Part 3: Cross-Document Gaps

---

## CROSS-01: Missing Page B44 (Integration Configuration Detail)
- [P0] **CRITICAL**: B26 WL Integration Center routes ALL "Configure"/"Connect" buttons to "B44 `/white-label/integrations/:type`" — **B44 has ZERO specification**. No page structure, no operations, no data model. 12 integration types depend on this page.

## CROSS-02: Missing Page B45 (Analytics)
- [P0] **CRITICAL**: B16 WL Hub Management routes "View Full Analytics" → "B45 `/analytics`" — **B45 not defined**. Likely typo for B43 (Page Analytics), but needs confirmation as B43 is page-specific, not project-wide.

## CROSS-03: B33 Preview Mode (WL context)
- [P1] B15 checklist references "Open Preview → B33" but B33 is only defined in `req_community.md`. WL preview differs (widgets, domain portal, branded pages) — no WL-specific preview spec.

## CROSS-04: D20 Specification Fragmentation
- [P1] D20 is referenced by 17 touchpoints across both docs but has no single authoritative spec. Community doc has partial spec (§7.2/§7.3), WL doc references audit_todo.md M-91/M-92/M-93. Need consolidated D20 spec in one location.

## CROSS-05: Sidebar Active State Gaps
- [P1] WL §15.1 defines sidebar submenu but doesn't specify active item for B42 (Iframe), B19 (Deploy Settings), B18 (Domain Setup), B48 (Dev Kit). These pages don't map to Overview/Widgets/Pages/Smart Rewards.

## CROSS-06: Global Interaction Norms Missing (WL)
- [P1] WL doc has no §2.3 "Global Interaction Norms" section. Toast styles, loading patterns, empty state illustrations, error page templates not centrally defined. Community doc has these but WL doesn't inherit explicitly.

## CROSS-07: B32 D20 Trigger Contradiction
- [P0] Community §7.2 trigger table lists B32 Content Mgmt "Publish (公告/内容)" as D20 trigger, but C-136 explicitly says "单条公告发布不走 D20" — **CONTRADICTION**. Which content types on B32 trigger D20?

---

# Summary: Top 10 Critical Blockers (P0)

| # | Issue | Doc | Impact |
|---|-------|-----|--------|
| 1 | **B44 Integration Config page completely undefined** | WL | 12 integration types blocked |
| 2 | **B45 Analytics page undefined** (likely typo for B43) | WL | B16 routing blocked |
| 3 | **D20 has no standalone spec** in WL doc | WL | 5 publish flows ambiguous |
| 4 | **B15→B16 transition threshold** "high traffic" undefined | WL | Route guard unimplementable |
| 5 | **B31 Task Creation/Edit Modal** has no D-code or spec | Community | Core CRUD blocked |
| 6 | **B31 Task archive vs b2c mapping** inconsistency | Community | C-end behavior unclear |
| 7 | **Multiple pages lack init API sequences** | Both | Devs reverse-engineer from scattered ops |
| 8 | **B32 D20 trigger contradiction** (§7.2 vs C-136) | Community | Publish flow ambiguous |
| 9 | **D13 Rule Editor** — event/badge/tier APIs unspecified | WL | Rule creation form blocked |
| 10 | **D15 Members Panel** — add/import/remove flows underspec'd | WL | Member management blocked |
