# B-End to C-End Operation Mapping

> Every B-end action has a C-end consequence. This document exhaustively maps all possible B-end operations to their corresponding C-end page changes.

## Table of Contents

1. [Design Principles](#1-design-principles)
2. [Entity Lifecycle States](#2-entity-lifecycle-states)
3. [Community-Level Operations](#3-community-level-operations)
4. [Sector & Task Operations](#4-sector--task-operations)
5. [Component Operations](#5-component-operations)
   - 5.1 Points System
   - 5.2 Level System
   - 5.3 Leaderboard
   - 5.4 TaskChain
   - 5.5 DayChain
   - 5.6 Milestones
   - 5.7 Benefits Shop
   - 5.8 Lucky Wheel
   - 5.9 Leaderboard Sprint
6. [Announcement & Featured Slots](#6-announcement--featured-slots)
7. [Anti-Sybil & Eligibility](#7-anti-sybil--eligibility)
8. [White Label Specifics](#8-white-label-specifics)
9. [Preview System Spec](#9-preview-system-spec)
10. [C-End Tab Visibility Rules](#10-c-end-tab-visibility-rules)
11. [Edge Cases & Error States](#11-edge-cases--error-states)

---

## 1. Design Principles

### B-end = Management Tool, C-end = Consumption Experience

| Dimension | B-End (Project Owner) | C-End (Participant) |
|-----------|----------------------|---------------------|
| Mental model | "Organize and optimize my community" | "What can I do and earn?" |
| Information density | High — all metrics, all states | Low — only actionable items |
| Organization | By management logic (sectors, modules, status) | By participation logic (motivation, urgency) |
| Time perspective | Lifecycle (past → future) | Present moment (what's available now) |

### Core Rule
> **C-end never shows B-end management artifacts.** Draft tasks, deactivated modules, sector management UI, analytics — none of this leaks to C-end. C-end only shows the *result* of B-end decisions.

---

## 2. Entity Lifecycle States

### 2.1 Community Lifecycle

```
(none) → [Create via Wizard] → Draft → [Activate] → Active → [Deactivate] → Paused → [Reactivate] → Active
                                                                                    → [Delete] → Deleted
```

| State | B-End View | C-End View |
|-------|-----------|-----------|
| **None** | Empty Hub (P24) — "Create Your First Community" | Community page does not exist (404 or "Coming Soon") |
| **Draft** | Wizard flow (P28-P30) — configuring | Community page does not exist |
| **Active** | Management Hub (P25-P27) with live metrics | Full community page visible with all active content |
| **Paused** | Hub shows "Paused" banner, all features read-only | Community page shows "This community is currently paused" message, no new actions |
| **Deleted** | Removed from sidebar, archived in settings | Community page returns 404 |

### 2.2 Module Lifecycle (Points, Leaderboard, etc.)

```
Not Configured → [Enable & Configure] → Active → [Edit] → Active
                                                → [Deactivate] → Inactive → [Reactivate] → Active
                                                                          → [Remove] → Not Configured
```

| State | B-End View | C-End View |
|-------|-----------|-----------|
| **Not Configured** | "Add Module" card in Hub | Tab/section does not appear |
| **Active** | Module card with metrics, click → detail page | Tab/section visible, fully functional |
| **Inactive** | Module card grayed with "Inactive" badge | Tab/section hidden, existing user data preserved but not accessible |

### 2.3 Task Lifecycle

```
Draft → [Publish] → Active → [User completes] → Completed (per user)
                           → [Expire] → Expired
                           → [Deactivate] → Inactive → [Reactivate] → Active
                           → [Delete] → Deleted
```

| State | B-End View | C-End View |
|-------|-----------|-----------|
| **Draft** | Listed with "Draft" tag, editable, not visible to users | Not visible |
| **Active** | Listed with "Active" tag, metrics updating | Visible in sector, completable |
| **Completed** (per user) | Completion count increments | Shows ✅ for that user, points credited |
| **Expired** | Listed with "Expired" tag, archived | Removed from active lists, shows in history if partially completed |
| **Inactive** | Listed with "Inactive" tag, can reactivate | Not visible (disappears from sector) |
| **Deleted** | Removed from list (or in trash) | Not visible, historical completions preserved in user records |

### 2.4 Sector Lifecycle

```
(none) → [Create] → Active → [Edit name/order] → Active
                            → [Hide] → Hidden → [Show] → Active
                            → [Delete] → Deleted
```

| State | B-End View | C-End View |
|-------|-----------|-----------|
| **Active** | Sector folder with tasks, draggable for reorder | Section visible on Home tab with header + task cards |
| **Hidden** | Sector with "Hidden" badge, tasks still manageable | Sector and its tasks not visible |
| **Deleted** | Removed, tasks become "uncategorized" | Sector disappears, tasks reassigned or removed |
| **Empty** (0 active tasks) | Sector shown with "No active tasks" note | Sector automatically hidden (no empty sections shown to users) |

---

## 3. Community-Level Operations

| # | B-End Operation | B-End UI Change | C-End Page Change |
|---|----------------|----------------|-------------------|
| 3.1 | **Create Community** (via Wizard) | Hub transitions from Empty → Guided state | Community page becomes accessible (initially minimal) |
| 3.2 | **Edit Community Name** | Name updates in Hub header | Project name updates in C-end header |
| 3.3 | **Edit Community Description** | Description updates in settings | Meta description updates (SEO), may show in About section |
| 3.4 | **Change Brand Color** | Color preview in settings | Accent color changes across all C-end UI elements |
| 3.5 | **Upload Logo** | Logo preview in settings | Logo updates in C-end header |
| 3.6 | **Activate Community** | Hub shows "Active" status, live metrics begin | C-end page becomes fully accessible |
| 3.7 | **Deactivate (Pause) Community** | Hub shows "Paused" banner, actions disabled | C-end shows paused message, no task completion possible |
| 3.8 | **Reactivate Community** | Hub returns to normal active state | C-end page resumes normal functionality |
| 3.9 | **Delete Community** | Community removed from sidebar, data archived | C-end returns 404, all user progress frozen (historical) |
| 3.10 | **Configure Eligibility** (who can join) | Eligibility settings panel | Non-eligible users see "Requirements not met" gate |
| 3.11 | **Add Community Operator** | Team member appears in operator list | No direct C-end change (operator manages from B-end) |
| 3.12 | **Enable Anti-Sybil** | Toggle in settings | Users must pass verification before participating |
| 3.13 | **Configure Referral Program** | Referral settings panel with link generation | "Invite Friends" button appears in C-end, referral tracking active |

---

## 4. Sector & Task Operations

### 4.1 Sector Operations

| # | B-End Operation | C-End Home Tab Change |
|---|----------------|----------------------|
| 4.1.1 | **Create Sector** (e.g., "Getting Started") | New section header appears on Home tab (at bottom by default) |
| 4.1.2 | **Rename Sector** | Section header text updates |
| 4.1.3 | **Reorder Sectors** (drag up/down) | Sections reorder on Home tab to match B-end sequence |
| 4.1.4 | **Hide Sector** | Section disappears from Home tab; tasks inside also hidden |
| 4.1.5 | **Show Sector** | Section reappears at its configured position |
| 4.1.6 | **Delete Sector** | Section disappears; tasks inside become uncategorized (hidden unless reassigned) |
| 4.1.7 | **Move Task to Different Sector** | Task card moves from one section to another |
| 4.1.8 | **Sector with 0 Active Tasks** | (auto) Section auto-hides — empty sectors never shown to C-end |

### 4.2 Task Operations

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 4.2.1 | **Create Task (Draft)** | No C-end change (drafts invisible) |
| 4.2.2 | **Publish Task** | Task card appears in its sector on Home tab + in Quests tab |
| 4.2.3 | **Edit Task Title/Description** | Task card text updates in real-time |
| 4.2.4 | **Change Task Points Value** | Points display on task card updates; already-earned points NOT retroactively changed |
| 4.2.5 | **Change Task Type** | Verification method changes (e.g., manual → auto); UX may change (upload vs. connect) |
| 4.2.6 | **Deactivate Task** | Task card disappears from all C-end views; users who started but didn't complete see "Task no longer available" |
| 4.2.7 | **Reactivate Task** | Task card reappears in its sector and Quests tab |
| 4.2.8 | **Delete Task** | Task permanently removed from C-end; completed users retain earned points |
| 4.2.9 | **Set Task as Recurring** (daily/weekly) | Task resets for users on configured interval; shows cooldown timer after completion |
| 4.2.10 | **Set Task Deadline** | Countdown timer appears on task card; auto-expires at deadline |
| 4.2.11 | **Set Max Claims** | "X/N claimed" counter appears; task grays out when limit reached |
| 4.2.12 | **Set Task as Required** | "Required" badge on task card; may gate other content |
| 4.2.13 | **Reorder Tasks within Sector** | Task card order changes within that section |
| 4.2.14 | **Add Eligibility Criteria to Task** | Lock icon on task card for non-eligible users with requirement tooltip |
| 4.2.15 | **Manual Approve** (for Proof of Work tasks) | User's pending submission changes from "Under Review" → "Approved ✅" or "Rejected ❌" |
| 4.2.16 | **Batch Publish Multiple Tasks** | Multiple task cards appear simultaneously |
| 4.2.17 | **Batch Deactivate** | Multiple task cards disappear simultaneously |

---

## 5. Component Operations

### 5.1 Points System

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.1.1 | **Enable Points System** | Points balance appears in User Status Bar; points earned on task completion |
| 5.1.2 | **Set Points Name** (e.g., "BTC Points") | All points displays use custom name |
| 5.1.3 | **Set Points Icon** | Custom icon appears next to points balance |
| 5.1.4 | **Configure Earning Rules** | Points earned per task type change |
| 5.1.5 | **Create Multiple Points Types** | Multiple balances shown (e.g., "Gold" + "Silver") |
| 5.1.6 | **Adjust Points Value for Task** | Updated points value shown on task card (future completions only) |
| 5.1.7 | **Import External Points** | Users see their imported balance reflected |
| 5.1.8 | **Deactivate Points System** | Points balance hidden; earning paused; existing balances frozen |
| 5.1.9 | **Export Points Data** | No C-end change (B-end only) |

### 5.2 Level System

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.2.1 | **Enable Level System** | Level badge appears in User Status Bar |
| 5.2.2 | **Configure Level Thresholds** | Level requirements change; some users may level up/down |
| 5.2.3 | **Set Level Names** (e.g., "Bronze", "Silver") | Level name updates in user profile |
| 5.2.4 | **Set Level-up Rewards** | Level-up notification includes reward claim |
| 5.2.5 | **Enable Level Gating** | Level-gated content shows lock icon with "Requires Lv.X" |
| 5.2.6 | **Deactivate Level System** | Level badge hidden; gating removed; progression frozen |

### 5.3 Leaderboard

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.3.1 | **Enable Leaderboard** | "Leaderboard" tab appears in C-end navigation |
| 5.3.2 | **Set Metric** (points/volume/referrals) | Leaderboard ranks by new metric; ranking may reshuffle |
| 5.3.3 | **Set Time Range** (all/24h/7d/30d) | Default time range changes; affects who's on top |
| 5.3.4 | **Set Display Count** (top N) | More or fewer users shown in ranking |
| 5.3.5 | **Configure Reward Tiers** | Tier badges appear next to ranks (e.g., "Top 10: 100 USDT") |
| 5.3.6 | **Enable Podium View** | Top 3 shown in podium layout instead of list |
| 5.3.7 | **Enable Rank Change Arrows** | ↑↓ indicators appear next to each rank |
| 5.3.8 | **Set Address Format** (full/short/ENS) | User address display format changes |
| 5.3.9 | **Set Accent Color** | Leaderboard UI accent color changes |
| 5.3.10 | **Deactivate Leaderboard** | "Leaderboard" tab disappears from C-end |

### 5.4 TaskChain

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.4.1 | **Create TaskChain** | Chain appears on Home tab (in its sector) and Quests tab |
| 5.4.2 | **Add Steps to Chain** | New steps appear (locked until prerequisites met) |
| 5.4.3 | **Remove Step from Chain** | Step disappears; users on that step skip to next |
| 5.4.4 | **Reorder Steps** | Chain sequence changes; user progress may reset for reordered steps |
| 5.4.5 | **Set Completion Reward** | Final reward preview shown at chain end |
| 5.4.6 | **Deactivate TaskChain** | Chain disappears from all views; in-progress users see "Chain paused" |
| 5.4.7 | **Delete TaskChain** | Chain permanently removed; earned step rewards retained |

### 5.5 DayChain (Daily Streaks)

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.5.1 | **Create DayChain** | DayChain section appears on Home tab; calendar/timeline view shown |
| 5.5.2 | **Set Total Days** | Calendar length changes (e.g., 7-day → 30-day) |
| 5.5.3 | **Configure Per-day Tasks** | Each day's task list updates |
| 5.5.4 | **Set Streak Rewards** (at milestones) | Streak milestone markers appear on calendar |
| 5.5.5 | **Enable Catch-up** | "Recover Missed Day" button appears for users who missed a day |
| 5.5.6 | **Set Catch-up Cost** | Cost shown on catch-up button (e.g., "Recover for 50 pts") |
| 5.5.7 | **Set Daily Reset Time** | Countdown timer adjusts to new reset time |
| 5.5.8 | **Change Layout** (calendar/timeline/cards) | Visual presentation changes |
| 5.5.9 | **Deactivate DayChain** | DayChain section disappears; existing streaks frozen |
| 5.5.10 | **Delete DayChain** | Permanently removed; streak records in user history |

### 5.6 Milestones

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.6.1 | **Enable Milestones** | "Milestone" tab appears in C-end navigation |
| 5.6.2 | **Set Metric** (points/volume/etc.) | Progress measured against new metric |
| 5.6.3 | **Add Threshold Tier** | New milestone marker appears in progression |
| 5.6.4 | **Remove Threshold Tier** | Milestone marker removed; users past it retain rewards |
| 5.6.5 | **Edit Tier Reward** | Reward display updates for unclaimed tiers; claimed rewards unchanged |
| 5.6.6 | **Change Display Style** | Visual presentation changes |
| 5.6.7 | **Deactivate Milestones** | "Milestone" tab disappears; earned achievements retained |

### 5.7 Benefits Shop

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.7.1 | **Enable Benefits Shop** | "Shop" tab appears in C-end navigation |
| 5.7.2 | **Add Shop Item** | New item card appears in shop |
| 5.7.3 | **Edit Item** (name/price/image) | Item card updates |
| 5.7.4 | **Set Item Stock** | "X remaining" counter appears; sold-out items grayed |
| 5.7.5 | **Set Item as Time-limited** | Countdown timer on item card |
| 5.7.6 | **Deactivate Item** | Item disappears from shop |
| 5.7.7 | **Reactivate Item** | Item reappears in shop |
| 5.7.8 | **Delete Item** | Permanently removed; past redemptions in user history |
| 5.7.9 | **Create Category** | New category filter/section in shop |
| 5.7.10 | **Reorder Items** | Item display order changes |
| 5.7.11 | **Change Points Name in Shop** | All "X points" labels use new name |
| 5.7.12 | **Deactivate Benefits Shop** | "Shop" tab disappears; existing redemptions honored |

### 5.8 Lucky Wheel

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.8.1 | **Enable Lucky Wheel** | "Lucky Wheel" tab appears in C-end navigation |
| 5.8.2 | **Configure Prize Segments** | Wheel segments update (colors, labels, probabilities) |
| 5.8.3 | **Set Spin Cost** (free/points) | Cost display on spin button changes |
| 5.8.4 | **Set Daily Limit** | "X/N spins remaining" counter appears |
| 5.8.5 | **Edit Prize** | Segment label/value updates |
| 5.8.6 | **Add/Remove Segment** | Wheel gains/loses a segment |
| 5.8.7 | **Deactivate Lucky Wheel** | "Lucky Wheel" tab disappears |

### 5.9 Leaderboard Sprint (Time-limited Competition)

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 5.9.1 | **Create Sprint** | "Sprint" tab appears (or section on Home) with countdown to end |
| 5.9.2 | **Set Sprint Duration** | Countdown timer updates |
| 5.9.3 | **Set Sprint Rewards** (per rank) | Prize pool / rank rewards displayed |
| 5.9.4 | **Set Sprint Metric** | Competition metric displayed (points/volume/tasks) |
| 5.9.5 | **End Sprint Early** | Sprint shows "Ended" state, final rankings frozen |
| 5.9.6 | **Sprint Naturally Ends** | (auto) Final rankings displayed, rewards distributed |
| 5.9.7 | **Delete Sprint** | Sprint removed from C-end; distributed rewards retained |

---

## 6. Announcement & Featured Slots

### 6.1 Announcements

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 6.1.1 | **Create Announcement** | Banner appears in announcement carousel on Home tab |
| 6.1.2 | **Edit Announcement** (text/image/link) | Banner content updates |
| 6.1.3 | **Set Announcement Order** | Carousel order changes |
| 6.1.4 | **Set Announcement as Pinned** | Banner always shows first in carousel |
| 6.1.5 | **Set Announcement Expiry** | Banner auto-removes after expiry |
| 6.1.6 | **Deactivate Announcement** | Banner removed from carousel |
| 6.1.7 | **Delete Announcement** | Permanently removed |

### 6.2 Featured Slots (2 rows x 3 = 6 cards on Home tab)

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 6.2.1 | **Assign Featured Slot** (link to quest/sprint/milestone) | Featured card appears in grid with image + title |
| 6.2.2 | **Change Slot Content** | Card image/title/link updates |
| 6.2.3 | **Remove Slot Content** | Slot becomes empty (hidden, or shows placeholder) |
| 6.2.4 | **Reorder Slots** | Card positions change in 2x3 grid |
| 6.2.5 | **All 6 Slots Empty** | Featured section auto-hides on C-end |

---

## 7. Anti-Sybil & Eligibility

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 7.1 | **Enable PoH (Proof of Humanity)** | Verification gate appears before first task; unverified users see "Verify your identity" |
| 7.2 | **Enable Enhanced PoH** | Stronger verification required (e.g., Worldcoin, BrightID) |
| 7.3 | **Enable Anti-Bot PoW** | PoW challenge appears before task submission |
| 7.4 | **Set Eligibility Criteria** (token balance, NFT, etc.) | Non-eligible users see requirements gate with clear instructions |
| 7.5 | **Create User Segment** | Segment-specific content shown/hidden per user group |
| 7.6 | **Add Address to Blacklist** | Blacklisted user sees "Access restricted" message |
| 7.7 | **Remove from Blacklist** | User regains normal access |

---

## 8. White Label Specifics

### 8.1 Widget Operations → C-End Widget Embed

| # | B-End Operation | C-End (Embedded Widget) Change |
|---|----------------|-------------------------------|
| 8.1.1 | **Create Widget** (e.g., Leaderboard widget) | Widget available for embed; shows in Page Builder library |
| 8.1.2 | **Configure Widget Props** | Widget appearance/data updates on all embeds |
| 8.1.3 | **Generate Embed Code** | No C-end change (B-end only; code for developer to embed) |
| 8.1.4 | **Deactivate Widget** | Embedded widget shows "Widget unavailable" placeholder |
| 8.1.5 | **Delete Widget** | Embedded locations show empty space or error |

### 8.2 Page Builder Operations → C-End Custom Page

| # | B-End Operation | C-End (Custom Page) Change |
|---|----------------|---------------------------|
| 8.2.1 | **Create Page** | New page URL becomes available |
| 8.2.2 | **Add Widget to Page** (drag & drop) | Widget appears on the custom page |
| 8.2.3 | **Remove Widget from Page** | Widget disappears from page |
| 8.2.4 | **Reorder Widgets on Page** | Layout sequence changes |
| 8.2.5 | **Set Page Title/Meta** | Page title and SEO metadata update |
| 8.2.6 | **Set Page Theme/Colors** | Entire page color scheme changes |
| 8.2.7 | **Save as Draft** | No C-end change (draft not published) |
| 8.2.8 | **Publish Page** | Page becomes accessible at its URL |
| 8.2.9 | **Unpublish Page** | Page returns 404 or shows "Under maintenance" |
| 8.2.10 | **Delete Page** | Page returns 404; widgets still exist independently |

### 8.3 Domain & Branding Operations

| # | B-End Operation | C-End Change |
|---|----------------|-------------|
| 8.3.1 | **Configure Custom Domain** | C-end accessible at custom URL (e.g., community.bitcoin.com) |
| 8.3.2 | **Set Brand Logo** | Logo updates across all C-end pages |
| 8.3.3 | **Set Brand Colors** | Theme color changes across all C-end pages |
| 8.3.4 | **Set Favicon** | Browser tab icon updates |
| 8.3.5 | **Remove "Powered by TaskOn"** | Footer attribution removed (WL Pro only) |
| 8.3.6 | **Configure SSO** | Login flow changes to project's own auth |

---

## 9. Preview System Spec

### 9.1 Preview Entry Points

| Entry Point | B-End Location | Preview Content | Preview Mode |
|-------------|---------------|----------------|-------------|
| Community Hub → [Preview] | Hub page top-right | Full C-end community (all tabs) | New window / full-screen overlay |
| Widget Config → Preview Panel | Widget config page right side | Single widget with mock data | Side panel (split view, ~40% width) |
| Page Builder → [Preview] | Page builder toolbar | Full custom page | Toggle mode (Edit ↔ Preview) |
| Task → [Preview] | Task edit page | Single task card as C-end user sees it | Inline preview below form |

### 9.2 Preview Data Strategy

| Data Type | Preview Source |
|-----------|---------------|
| User identity | Mock user ("Preview User", Lv.5, 1,250 pts, 7-day streak) |
| Task states | Mix of completed/available/locked to show all visual states |
| Leaderboard | Mock data (10 entries with realistic names/scores) |
| Points balance | Mock balance (e.g., 1,250) |
| Streaks | Mock streak (e.g., 7 days) |
| Shop items | Real items from B-end configuration |
| Announcements | Real announcements from B-end |

### 9.3 Preview Banner

All preview modes show a non-dismissible top banner:
```
⚠️ Preview Mode — This is how participants see your community. [Exit Preview]
```

---

## 10. C-End Tab Visibility Rules

C-end tabs are **dynamically shown/hidden** based on B-end module activation:

| C-End Tab | Requires B-End Module | Shows When |
|-----------|----------------------|-----------|
| **Home** | (always) | Always visible (core tab) |
| **Quests** | Any active tasks | At least 1 published task exists |
| **Leaderboard** | Leaderboard module | Leaderboard enabled AND at least 1 entry |
| **LB Sprint** | Leaderboard Sprint | Active or recently-ended sprint exists |
| **Milestone** | Milestones module | Milestones enabled AND at least 1 threshold |
| **Shop** | Benefits Shop module | Shop enabled AND at least 1 active item |
| **Lucky Wheel** | Lucky Wheel module | Wheel enabled AND at least 1 prize configured |

**Rule**: If a tab has no content (module enabled but 0 items), it is **hidden** from C-end to avoid empty states.

### Home Tab Content Ordering

The Home tab shows content in this fixed order:
1. Announcement carousel (if any active announcements)
2. Featured grid 2×3 (if any featured slots filled)
3. User's active DayChain (if DayChain enabled, prominently placed)
4. User's active TaskChains (if any in-progress)
5. Sectors with tasks (in B-end configured order)

---

## 11. Edge Cases & Error States

### 11.1 Task Edge Cases

| Scenario | B-End View | C-End View |
|----------|-----------|-----------|
| Task published to sector that is hidden | Warning: "Task in hidden sector" | Task not visible (sector hidden) |
| Task deadline passed | Auto-status: "Expired" | Task removed from active lists |
| Task max claims reached | Shows "Max reached (500/500)" | Task grayed: "Fully claimed" |
| All tasks in sector expired/inactive | Sector shows "0 active tasks" | Sector auto-hides |
| Recurring task in cooldown | Shows completion stats | Shows countdown timer: "Available in 3h 22m" |
| Task requires manual review | Shows review queue count | User sees "Under Review ⏳" after submission |

### 11.2 Component Edge Cases

| Scenario | C-End View |
|----------|-----------|
| Leaderboard enabled but 0 participants | Tab hidden (rule: no empty tabs) |
| Shop enabled but 0 items | Tab hidden |
| Shop item out of stock | Item shows "Sold Out" badge, redemption button disabled |
| Lucky Wheel all spins used today | Wheel visible but spin button disabled: "Come back tomorrow" |
| LB Sprint ended, rewards pending | LB Sprint tab shows final standings + "Rewards distributing..." |
| DayChain streak broken | Calendar shows missed day in red; catch-up button (if enabled) |
| DayChain all days completed | Calendar shows 100% complete; congratulation message |

### 11.3 Community-Level Edge Cases

| Scenario | C-End View |
|----------|-----------|
| Community paused | Message: "This community is taking a break. Check back later." |
| Community deleted | 404 page |
| User blacklisted | Message: "Access to this community is restricted." |
| User doesn't meet eligibility | Gate page showing requirements + progress toward meeting them |
| No modules active (fresh community) | Minimal page with project info + "Coming soon" message |

### 11.4 White Label Edge Cases

| Scenario | C-End View |
|----------|-----------|
| Widget deactivated but still embedded | Placeholder: "This widget is currently unavailable" |
| Custom domain DNS not verified | Falls back to taskon.xyz/project-name URL |
| Page Builder page unpublished | 404 at custom URL |
| SSO configured but auth fails | Fallback to TaskOn default login |

---

## Appendix: Operation Count Summary

| Category | Operation Count |
|----------|----------------|
| Community-level | 13 |
| Sector operations | 8 |
| Task operations | 17 |
| Points System | 9 |
| Level System | 6 |
| Leaderboard | 10 |
| TaskChain | 7 |
| DayChain | 10 |
| Milestones | 7 |
| Benefits Shop | 12 |
| Lucky Wheel | 7 |
| Leaderboard Sprint | 7 |
| Announcements | 7 |
| Featured Slots | 5 |
| Anti-Sybil & Eligibility | 7 |
| White Label Widget | 5 |
| White Label Page Builder | 10 |
| White Label Domain/Branding | 6 |
| **Total** | **~153 operations** |

---

*Last updated: 2026-03-04*
*Document version: 1.0*
