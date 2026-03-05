# White Label Integration Modes

> White Label is built on top of Quest + Community. Projects must first set up their Community (tasks, points, leaderboard, etc.) before configuring White Label integration.

---

## Prerequisites

Before using any White Label integration mode, the project must have:
- An active **Community** with configured modules (tasks, points, leaderboard, milestones, etc.)
- Optionally, active **Quest** campaigns

White Label does not replace Community — it **extends** Community by embedding it into the project's own domain and UX.

---

## Four Integration Modes

### Mode 1: Domain Mode

**Code required:** None
**Operator:** Marketing team (no engineering needed)

**How it works:**
- Maps the project's entire Community page to the project's custom domain (e.g., `rewards.yourproject.xyz`)
- The Community page — with all its tasks, points, leaderboard, etc. — becomes accessible via the project's URL
- DNS configuration (CNAME record) is the only technical step

**Pros:**
- Zero code — marketing team can handle everything
- Fastest setup (minutes)
- Page can be used standalone or linked from the project's main site

**Cons:**
- Users still log in via TaskOn's authentication (not the project's login system)
- Limited CSS customization — primarily uses TaskOn's default styles with color/branding configuration
- The page looks and feels like a TaskOn page under the project's domain

**Best for:** Projects that want quick deployment and don't need deep UX integration.

---

### Mode 2: Embed Mode

**Code required:** Minimal (iframe/script tag integration)
**Operator:** Marketing team + minor engineering support

**How it works:**
- Embeds the entire Community page into the project's website as a second-level module (e.g., a "Rewards" tab within the project's main site)
- Users log in once via the project's own authentication — no separate TaskOn login required (SSO integration)
- The embedded page appears within the project's site navigation

**Pros:**
- Single sign-on — users only log in via the project's system
- Appears as a native section of the project's website
- Minimal code integration (embed script)

**Cons:**
- Limited CSS customization — primarily uses TaskOn's styles with color/branding overrides
- The embedded module's look & feel may not perfectly match the project's design system

**Best for:** Projects that want seamless user authentication and basic integration without heavy development.

---

### Mode 3: Widget + Page Builder Mode

**Code required:** Low (paste widget embed code into site)
**Operator:** Marketing team (widget placement) + minimal engineering (embed code)

**How it works:**

#### Widgets
- Individual Community modules are packaged as **embeddable widgets**
- Each widget corresponds to a specific Community feature that has already been created
- Projects can embed individual widgets into **any page, at any position** within their website
- Widgets inherit the project's CSS as closely as possible, maximizing visual consistency
- No separate login required — uses the project's authentication

**Available widget types** (each requires a corresponding Community instance):
| Widget | Source |
|--------|--------|
| Task List | Community → Tasks |
| TaskChain | Community → TaskChain |
| DayChain | Community → DayChain |
| Sector | Community → Sectors |
| Quest | Quest campaigns |
| Leaderboard | Community → Leaderboard |
| Milestone | Community → Milestones |
| User Center | Community → User profile/assets |
| Trading Race (CEX) | Community → Trading competitions |
| Rewards Shop | Community → Benefits Shop |

**Important:** Creating a widget requires that the corresponding Community feature instance already exists. A widget is essentially a **display wrapper + configuration layer** around an existing Community module.

#### Page Builder
- Composes multiple widgets into a **full page layout** (like a user engagement center)
- Pages created in Page Builder are themselves widgets — they can be embedded into the project's site
- Useful for creating:
  - **Permanent pages**: A dedicated "Rewards Hub" or "User Center" as a second-level section
  - **Temporary pages**: Landing pages for ad campaigns, seasonal events, or time-limited promotions
- **Dependency**: Page Builder requires existing widgets — it arranges and lays out widgets that have already been created

**Creation dependency chain:**
```
Community features (tasks, points, leaderboard...)
    ↓ must exist first
Widgets (display wrappers for each feature)
    ↓ must exist first
Page Builder pages (layout compositions of widgets)
    ↓ can then be
Embedded into project's website via embed code
```

**Pros:**
- Maximum CSS consistency with the project's website
- Granular placement — embed individual widgets anywhere
- Page Builder enables rapid creation of campaign/event landing pages
- No separate login — uses project's authentication
- Marketing team can create and iterate pages without engineering

**Cons:**
- Requires pasting embed code into the project's site (minimal engineering)
- Must create Community features → Widgets → Pages in sequence
- More setup steps than Domain or Embed modes

**Best for:** Projects that need pixel-perfect brand consistency and flexible, composable growth experiences.

---

### Mode 4: SDK Integration

**Code required:** Full development effort
**Operator:** Engineering team

**How it works:**
- Full programmatic access to TaskOn's growth engine via SDK/API
- Projects build their own UI on top of TaskOn's backend services
- Complete control over UX, data flow, and business logic
- Headless architecture — TaskOn provides the engine, the project provides the interface

**Pros:**
- 100% custom UX — no constraints on design or interaction patterns
- Deepest integration with the project's tech stack
- Full data access and control
- Maximum flexibility for unique use cases

**Cons:**
- Requires engineering development time
- Longer implementation timeline
- Project is responsible for UI maintenance and updates

**Best for:** Projects with dedicated engineering teams that need completely custom experiences.

---

## Mode Comparison

| Aspect | Domain | Embed | Widget + Page Builder | SDK |
|--------|--------|-------|----------------------|-----|
| **Code needed** | None | Minimal | Low | Full |
| **Setup time** | Minutes | Hours | Hours–Days | Days–Weeks |
| **CSS match** | Low (TaskOn styles) | Low (TaskOn styles) | High (project CSS) | Full (custom) |
| **Login** | TaskOn auth | Project auth (SSO) | Project auth | Project auth |
| **Placement** | Standalone page | Embedded section | Any page, any position | Custom |
| **Page creation** | N/A | N/A | Page Builder | Custom dev |
| **Temporary pages** | No | No | Yes (landing pages, events) | Custom dev |
| **Operator** | Marketing | Marketing + Eng | Marketing + Eng | Engineering |
| **Best for** | Quick launch | Basic integration | Brand-consistent | Full custom |

---

## Recommended Progression

Most projects follow this adoption path:

```
Domain Mode (Day 1)
  → "Let's see if this works for our users"
  → Zero effort, validate the concept

Embed Mode (Week 1)
  → "Users shouldn't need a separate login"
  → Add SSO, embed as site section

Widget + Page Builder (Week 2-4)
  → "We want this to look like our product"
  → Decompose into widgets, customize CSS, build pages

SDK Integration (Month 2+)
  → "We need completely custom flows"
  → Build unique UX on top of TaskOn's engine
```

---

## Key Principle: Build Bottom-Up

The creation flow is always bottom-up:

```
1. Set up Community (tasks, points, leaderboard, milestones, etc.)
2. Create Quest campaigns (if needed)
3. Package features as Widgets (for Mode 3)
4. Compose Widgets into Pages via Page Builder (for Mode 3)
5. Embed into project's website
```

You cannot skip steps. Widgets need Community features. Page Builder needs Widgets. This ensures data integrity and prevents broken experiences.
