# Design TODO

> Active design tasks and planned changes for `design/pencil-new.pen`

---

## ✅ P0 — Community Empty State Redesign (B09 `zzZ8D`) — DONE

### Problem
Current 3 templates (Points & Levels / Leaderboards & Rewards / Benefits Shop) are **module-level** — they give project owners individual parts but no strategy for HOW to retain users after Quest acquisition.

### New Direction: Retention Strategy Templates

Replace module-level templates with **3 lifecycle-stage strategies**, each being a complete retention loop:

| | Strategy 1 | Strategy 2 | Strategy 3 |
|---|---|---|---|
| **Name** | Activate New Users | Drive Daily Engagement | Maximize Retention |
| **Subtitle** | Convert Quest participants into active members | Build habits that bring users back every day | Create investment that makes leaving costly |
| **AARRR Stage** | Activation | Retention | Deep Retention |
| **Modules** | TaskChain + Points + Levels | DayChain + Leaderboard + Sprint | Benefits Shop + Milestones + Lucky Wheel |
| **Psychology** | Goal gradient + instant feedback | Loss aversion (streaks) + social comparison | Sunk cost (spend points) + variable reinforcement (lucky wheel) |
| **Key Metric** | Activation Rate | Daily Active Users | 30-Day Retention |
| **Economy Role** | **Earn** — complete tasks, earn points, level up | **Compete** — daily check-in, rank, sprint | **Spend** — redeem, milestone rewards, lucky draw |

### User Journey Flow

```
Quest Acquisition → 🎯 Activate → 🔥 Daily Engage → 💎 Deep Retain
                     (TaskChain)   (DayChain)         (Benefits Shop)
                     (Points)      (Leaderboard)      (Milestones)
                     (Levels)      (Sprint)           (Lucky Wheel)
```

### Selection → Onboarding Flow

```
Select Strategy on Empty State (zzZ8D)
    ↓
Community Wizard (Gzpeu) — pre-fills selected strategy's modules
    ↓
Guided Workspace (S1EIA) — shows activated modules + setup checklist
    ↓
Module Management Pages — configure each module in detail
    ↓
Back to Guided Workspace — prompt: "Ready to add next strategy?"
```

### Sectors & Tasks Position
Sectors & Tasks is **infrastructure** — required by ALL strategies. Auto-configured in Wizard Step 1 (not part of any template card).

### Design Changes Required

1. **Redesign template cards section** on `zzZ8D`:
   - Replace 3 module cards with 3 strategy cards
   - Each card: icon + strategy name + subtitle + included modules (as chips/tags) + key metric
   - Visual: lifecycle progression arrow connecting the 3 cards (Activate → Engage → Retain)

2. **Update highlight strip** below cards:
   - Change from individual module names to: "Or configure individual modules: Sectors & Tasks, Points & Level, TaskChain, DayChain, Leaderboard, Sprint, Milestone, Benefits Shop, Lucky Wheel"
   - Add "Start from Scratch" secondary link

3. **Update CTA**:
   - Primary: "Create Your First Community" (opens Wizard with selected strategy pre-filled)
   - If no strategy selected: Wizard opens with blank slate

4. **Optionally add lifecycle diagram**:
   - Small visual showing Quest → Activate → Engage → Retain flow
   - Reinforces the strategic narrative

### Affected Pages — ALL DONE
- `zzZ8D` (B09) — ✅ Redesigned: 3 retention strategy cards + lifecycle flow
- `Gzpeu` (B13) — ✅ Redesigned: Step 1 Customize (3-step stepper, was 4)
- `8NeyG` (B13b) — ✅ Redesigned: Step 2 Configure Modules (9 toggles, strategy pre-fill)
- `qknQZ` (B13c) — ✅ Redesigned: Step 3 Review & Launch (summary + launch button)

### New Wizard Flow (3 steps, was 4)
```
Step 1: Customize (Gzpeu) — Community name, description, brand color
Step 2: Configure Modules (8NeyG) — 9 module toggles, pre-filled by strategy
Step 3: Review & Launch (qknQZ) — Summary + "Launch Community" button
```

---

## 🟡 P1 — Unified Pricing Page Polish

Page `HO2Ny` is functional but may need:
- Tab switching states for Quest/Community/White Label
- Mobile responsive considerations
- Final copy review

---

## 🟢 P2 — Frontend Handoff Prep

- Verify all button routing annotations are complete and accurate
- Cross-check `docs/website_frontend_requirements.md` with actual page designs
- Ensure all page codes (M01-M14, B01-B47, C01-C06) are documented
