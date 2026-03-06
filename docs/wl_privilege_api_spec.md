# White Label Privilege & Activity Rules API Specification

> Version: 1.0 Draft
> Date: 2026-03-05
> Status: Design phase — for frontend + backend alignment

---

## Overview

Two WL-exclusive systems that require project-side integration:

1. **Activity Rules** — TaskOn monitors project's on-chain contract events, auto-awards points to users
2. **Privilege System** — TaskOn manages who qualifies for project-native benefits (fee discounts, gas rebates, yield boosts), project applies them

### Responsibility Split

| Concern | TaskOn manages | Project manages |
|---------|---------------|-----------------|
| **Activity Rules** | Contract registry, event→points mapping, anti-sybil, points distribution | Deploying contracts, emitting standard events |
| **Privileges** | Privilege definitions, qualification rules (level/badge/points), user eligibility tracking, analytics | Querying user eligibility, applying benefits in their own logic, reporting usage |

---

## Part 1: Activity Rules API

### 1.1 Contract Registry

Project registers their smart contracts in WL Admin. TaskOn indexes events from these contracts.

#### Register Contract

```
POST /api/v1/wl/{projectId}/contracts
```

Request:
```json
{
  "chainId": 1,
  "address": "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
  "name": "DEX Router",
  "abiSource": "auto"  // "auto" = fetch from block explorer, "manual" = provide ABI
}
```

Response:
```json
{
  "id": "contract_abc123",
  "chainId": 1,
  "address": "0x7a25...488D",
  "name": "DEX Router",
  "status": "verified",
  "detectedEvents": [
    { "name": "Swap", "signature": "Swap(address,uint256,uint256,uint256,uint256,address)", "params": ["sender", "amount0In", "amount1In", "amount0Out", "amount1Out", "to"] },
    { "name": "AddLiquidity", "signature": "...", "params": ["..."] }
  ],
  "detectedFunctions": [
    { "name": "swapExactTokensForTokens", "params": ["amountIn", "amountOutMin", "path", "to", "deadline"] }
  ]
}
```

#### List Contracts

```
GET /api/v1/wl/{projectId}/contracts
```

#### Delete Contract

```
DELETE /api/v1/wl/{projectId}/contracts/{contractId}
```

### 1.2 Activity Rules (Event → Points Mapping)

#### Create Rule

```
POST /api/v1/wl/{projectId}/activity-rules
```

Request:
```json
{
  "name": "Reward Every Swap",
  "contractId": "contract_abc123",
  "trigger": {
    "event": "Swap",
    "conditions": [
      { "param": "amount0In", "operator": "gte", "value": "100000000" }  // >= $100 (USDC 6 decimals)
    ]
  },
  "reward": {
    "points": 50,
    "pointType": "default",       // or custom point type ID
    "multipliers": [
      { "type": "first_daily", "value": 2.0, "description": "First action of the day" },
      { "type": "streak", "value": 1.5, "minDays": 7, "description": "7+ day streak bonus" }
    ]
  },
  "limits": {
    "perUserPerDay": 5,
    "perUserPerWeek": 20,
    "totalDailyBudget": 10000       // max total points distributed per day
  },
  "antiSybil": {
    "minWalletAgeDays": 7,
    "minUniqueInteractions": 3,
    "excludeContracts": false       // if true, only EOA addresses qualify
  },
  "status": "active"
}
```

Response:
```json
{
  "id": "rule_xyz789",
  "name": "Reward Every Swap",
  "status": "active",
  "createdAt": "2026-03-05T10:00:00Z",
  "stats": {
    "triggeredToday": 0,
    "pointsDistributedToday": 0
  }
}
```

#### List Rules

```
GET /api/v1/wl/{projectId}/activity-rules
GET /api/v1/wl/{projectId}/activity-rules/{ruleId}/stats?period=7d
```

#### Preset Templates

```
GET /api/v1/wl/activity-rule-templates
```

Returns:
```json
[
  {
    "id": "template_dex_swap",
    "name": "Reward Every Swap",
    "description": "Award points for each swap above a threshold",
    "category": "DEX",
    "requiredEvents": ["Swap"],
    "defaultConfig": { "...": "..." }
  },
  {
    "id": "template_daily_first",
    "name": "Daily First Action",
    "description": "2x points for the first contract interaction each day",
    "category": "Engagement",
    "requiredEvents": ["any"],
    "defaultConfig": { "...": "..." }
  },
  {
    "id": "template_lp_bonus",
    "name": "Liquidity Provider Bonus",
    "description": "Bonus points for adding liquidity",
    "category": "DeFi",
    "requiredEvents": ["AddLiquidity", "Deposit"],
    "defaultConfig": { "...": "..." }
  },
  {
    "id": "template_staking_milestone",
    "name": "Staking Duration Milestone",
    "description": "Bonus at 7d, 30d, 90d staking duration",
    "category": "DeFi",
    "requiredEvents": ["Stake"],
    "defaultConfig": { "...": "..." }
  }
]
```

### 1.3 Activity Events (read-only, for analytics)

```
GET /api/v1/wl/{projectId}/activity-events?ruleId=rule_xyz789&period=7d
```

Response:
```json
{
  "summary": {
    "totalTriggered": 1247,
    "uniqueUsers": 342,
    "pointsDistributed": 62350,
    "sybilBlocked": 23
  },
  "timeSeries": [
    { "date": "2026-03-04", "triggered": 189, "users": 67, "points": 9450 },
    { "date": "2026-03-03", "triggered": 201, "users": 72, "points": 10050 }
  ]
}
```

---

## Part 2: Privilege System API

### 2.1 Privilege Definitions (B-End Admin manages)

#### Create Privilege

```
POST /api/v1/wl/{projectId}/privileges
```

Request:
```json
{
  "name": "Trading Fee Discount",
  "type": "fee_discount",              // fee_discount | gas_rebate | yield_boost | priority_access | custom
  "value": {
    "discountPercent": 30               // type-specific value schema
  },
  "description": "30% trading fee reduction on all pairs",
  "qualification": {
    "mode": "status",                   // "status" (Mode A) | "achievement" (Mode C)
    "criteria": {
      "type": "level",                  // level | points_balance | badge | milestone | custom
      "operator": "gte",
      "value": 3                        // Level >= 3
    }
  },
  "duration": {
    "type": "while_qualified",          // "while_qualified" (auto-revoke if criteria not met) | "permanent" | "fixed_days"
    "days": null
  },
  "limits": {
    "maxHolders": null,                 // null = unlimited
    "totalBudget": null                 // for gas_rebate: max total refund amount
  },
  "status": "active"
}
```

**Mode A (Status-based) examples:**
```json
// Level-gated fee discount — auto-activate, auto-revoke
{ "qualification": { "mode": "status", "criteria": { "type": "level", "operator": "gte", "value": 3 } }, "duration": { "type": "while_qualified" } }

// Points-balance-gated yield boost — must maintain 1000+ pts
{ "qualification": { "mode": "status", "criteria": { "type": "points_balance", "operator": "gte", "value": 1000 } }, "duration": { "type": "while_qualified" } }
```

**Mode C (Achievement-based) examples:**
```json
// Badge unlock — permanent once earned
{ "qualification": { "mode": "achievement", "criteria": { "type": "badge", "value": "diamond_hands" } }, "duration": { "type": "permanent" } }

// Milestone unlock — permanent
{ "qualification": { "mode": "achievement", "criteria": { "type": "milestone", "value": "first_steps" } }, "duration": { "type": "permanent" } }
```

#### List Privileges

```
GET /api/v1/wl/{projectId}/privileges
```

Response includes live stats:
```json
[
  {
    "id": "priv_fee30",
    "name": "Trading Fee Discount",
    "type": "fee_discount",
    "value": { "discountPercent": 30 },
    "qualification": { "mode": "status", "criteria": { "type": "level", "operator": "gte", "value": 3 } },
    "stats": {
      "activeHolders": 247,
      "totalValueDistributed": "$4,720",  // calculated from usage reports
      "estimatedROI": "8.4x"              // if project reports revenue uplift
    }
  }
]
```

### 2.2 User Privilege Query (Project calls this)

**This is the core integration point.** Project's backend/contract queries TaskOn to check a user's active privileges before applying benefits.

#### Query Single User Privileges

```
GET /api/v1/wl/{projectId}/users/{walletAddress}/privileges
```

Headers:
```
Authorization: Bearer {project_api_key}
```

Response:
```json
{
  "address": "0x1234...abcd",
  "privileges": [
    {
      "id": "priv_fee30",
      "type": "fee_discount",
      "name": "Trading Fee Discount",
      "value": { "discountPercent": 30 },
      "qualifiedVia": "level",           // how user qualified
      "qualifiedValue": 5,               // user's current level
      "activeSince": "2026-02-15T00:00:00Z",
      "expiresAt": null,                 // null = while_qualified / permanent
      "status": "active"
    },
    {
      "id": "priv_gas5",
      "type": "gas_rebate",
      "name": "Monthly Gas Rebate",
      "value": { "rebateUSD": 5, "period": "monthly" },
      "qualifiedVia": "badge",
      "qualifiedValue": "diamond_hands",
      "activeSince": "2026-01-20T00:00:00Z",
      "expiresAt": null,
      "status": "active"
    }
  ],
  "summary": {
    "totalActivePrivileges": 2,
    "totalSavedThisMonth": "$65.70"
  }
}
```

#### Batch Query (for multiple users)

```
POST /api/v1/wl/{projectId}/users/privileges/batch
```

Request:
```json
{
  "addresses": ["0x1234...abcd", "0x5678...efgh", "0x9abc...ijkl"],
  "privilegeType": "fee_discount"     // optional filter
}
```

Response:
```json
{
  "results": {
    "0x1234...abcd": { "fee_discount": { "discountPercent": 30, "status": "active" } },
    "0x5678...efgh": { "fee_discount": null },
    "0x9abc...ijkl": { "fee_discount": { "discountPercent": 20, "status": "active" } }
  }
}
```

#### Quick Check (lightweight, for hot-path contract calls)

```
GET /api/v1/wl/{projectId}/check?address=0x1234...abcd&privilege=fee_discount
```

Response (minimal):
```json
{
  "eligible": true,
  "value": { "discountPercent": 30 }
}
```

Latency target: < 50ms. Cached. Project can also cache locally with webhook for invalidation.

### 2.3 Webhooks (Push notifications to project)

Project registers a webhook URL to receive real-time privilege status changes.

#### Register Webhook

```
POST /api/v1/wl/{projectId}/webhooks
```

Request:
```json
{
  "url": "https://api.myproject.com/taskon/webhook",
  "events": ["privilege.activated", "privilege.deactivated", "privilege.expiring"],
  "secret": "whsec_..."
}
```

#### Webhook Payload Examples

**privilege.activated** — User newly qualifies:
```json
{
  "event": "privilege.activated",
  "timestamp": "2026-03-05T14:30:00Z",
  "data": {
    "address": "0x1234...abcd",
    "privilegeId": "priv_fee30",
    "type": "fee_discount",
    "value": { "discountPercent": 30 },
    "reason": "User reached Level 3"
  }
}
```

**privilege.deactivated** — User no longer qualifies (Mode A: level dropped):
```json
{
  "event": "privilege.deactivated",
  "timestamp": "2026-03-05T15:00:00Z",
  "data": {
    "address": "0x1234...abcd",
    "privilegeId": "priv_fee30",
    "type": "fee_discount",
    "reason": "User dropped below Level 3 (current: Level 2)"
  }
}
```

**privilege.expiring** — Privilege will expire soon (for fixed-duration):
```json
{
  "event": "privilege.expiring",
  "timestamp": "2026-03-05T10:00:00Z",
  "data": {
    "address": "0x1234...abcd",
    "privilegeId": "priv_yield_boost",
    "type": "yield_boost",
    "expiresAt": "2026-03-07T00:00:00Z",
    "hoursRemaining": 38
  }
}
```

### 2.4 Usage Reporting (Project reports back)

Project reports when a privilege was actually used, enabling ROI analytics in WL Admin.

```
POST /api/v1/wl/{projectId}/privilege-usage
```

Request:
```json
{
  "address": "0x1234...abcd",
  "privilegeId": "priv_fee30",
  "event": "applied",
  "details": {
    "originalFee": "0.003",
    "discountedFee": "0.0021",
    "savedAmount": "0.0009",
    "savedUSD": 1.85,
    "txHash": "0xabc..."
  },
  "timestamp": "2026-03-05T14:35:00Z"
}
```

This feeds back into:
- B-End Privilege Manager: "Total Value Distributed: $4,720" / "Estimated ROI: 8.4x"
- C-End User Center: "Saved $65.70 this month with your privileges"

### 2.5 On-Chain Privilege Contract (Optional, Mode C)

For projects that want trustless on-chain verification (no API dependency), TaskOn deploys a Privilege Registry contract per project.

```solidity
// Simplified interface — TaskOn deploys and manages
interface ITaskOnPrivileges {
    // Returns true if user has the specified privilege
    function hasPrivilege(address user, bytes32 privilegeType) external view returns (bool);

    // Returns the privilege value (e.g., discount basis points)
    function getPrivilegeValue(address user, bytes32 privilegeType) external view returns (uint256);

    // Returns all active privileges for a user
    function getUserPrivileges(address user) external view returns (Privilege[] memory);

    struct Privilege {
        bytes32 privilegeType;    // e.g., keccak256("fee_discount")
        uint256 value;            // e.g., 3000 = 30% (basis points)
        uint256 activeSince;
        uint256 expiresAt;        // 0 = no expiry
    }
}
```

**Project integration example (in their DEX contract):**

```solidity
import "./ITaskOnPrivileges.sol";

contract ProjectDEX {
    ITaskOnPrivileges public taskonPrivileges;
    uint256 public baseFeeRate = 30; // 0.3% in basis points * 100

    function swap(address tokenIn, address tokenOut, uint256 amountIn) external {
        uint256 feeRate = baseFeeRate;

        // Check TaskOn privilege
        if (taskonPrivileges.hasPrivilege(msg.sender, keccak256("fee_discount"))) {
            uint256 discount = taskonPrivileges.getPrivilegeValue(msg.sender, keccak256("fee_discount"));
            feeRate = baseFeeRate * (10000 - discount) / 10000;
        }

        // Apply discounted fee
        uint256 fee = amountIn * feeRate / 10000;
        // ... rest of swap logic
    }
}
```

TaskOn keeps the Privilege Registry updated via Merkle root updates or direct state writes (depending on gas cost/chain).

---

## Part 3: Integration Flow Summary

### For a DEX project integrating both Activity Rules + Privileges:

```
Step 1: Register contracts in WL Admin
        → TaskOn starts monitoring swap/LP events

Step 2: Configure Activity Rules
        → "Every swap ≥ $100 → +50 pts, daily first = 2x"
        → Users automatically earn points when they swap

Step 3: Create Privileges in WL Admin
        → "Level 3+ → 30% fee discount"
        → "Diamond Hands badge → permanent 10% discount"

Step 4: Integrate Privilege Check (choose one)
        → Option A: REST API (GET /check?address=...&privilege=fee_discount)
        → Option B: Webhook listener (cache locally, update on changes)
        → Option C: On-chain contract (taskonPrivileges.hasPrivilege(...))

Step 5: Report Usage (optional but recommended)
        → POST /privilege-usage when discount is applied
        → Enables ROI dashboard in WL Admin

Step 6: Add C-End Widgets
        → Activity Reward Toast (real-time "+50 pts" on swap)
        → Privilege Status Widget (shows active benefits + savings)
```

### Integration Effort Estimate

| Integration Mode | Effort | Latency | Trust Model |
|-----------------|--------|---------|-------------|
| REST API (Quick Check) | 1-2 hours | ~50ms (cached) | API key auth |
| Webhooks + Local Cache | 3-4 hours | ~0ms (local) | Webhook secret |
| On-Chain Contract | 1-2 days | ~0ms (on-chain read) | Trustless |

---

## Part 4: B-End Admin Page Specifications

### Contract Registry Page

**Header**: "Smart Contracts" + [+ Register Contract]
**Stats**: Registered Contracts: 3 | Active Listeners: 2 | Events Detected Today: 1,247
**Content**:
- Contract cards: name, address (truncated), chain icon, status badge (Verified/Pending/Error), detected events count
- Each card expands to show: ABI events list, linked Activity Rules count, last event timestamp
**Empty state**: "Register your first contract to start auto-rewarding on-chain actions" + protocol template selector

### Rule Builder Page

**Header**: "Activity Rules" + [+ Create Rule]
**Stats**: Active Rules: 4 | Triggered Today: 892 | Points Distributed Today: 44,600
**Filter tabs**: Active | Paused | All
**Content**:
- Rule cards: name, contract name, event type, points per trigger, today's stats (triggers, points, unique users)
- Expandable: conditions detail, multipliers, limits, anti-sybil settings
**Create/Edit flow**: Visual if-then builder (similar to Zapier trigger-action pattern)
- Step 1: Select contract → Select event
- Step 2: Add conditions (amount ≥ X, frequency limit, etc.)
- Step 3: Define reward (base points + multipliers)
- Step 4: Set limits + anti-sybil
- Live preview: "When a user swaps ≥ $100 on DEX Router → Award 50 pts (100 pts if first swap of day). Max 5 per user per day."

### Privilege Manager Page

**Header**: "Privileges" + [+ Create Privilege]
**Stats**: Active Privileges: 3 | Users with Privileges: 482 | Total Value This Month: $4,720 | Est. ROI: 8.4x
**Content**:
- Privilege cards:
  - Icon + Name + Type badge (Fee Discount / Gas Rebate / Yield Boost / Priority Access)
  - Value: "30% fee reduction"
  - Qualification: "Level 3+" or "Badge: Diamond Hands"
  - Mode: "Status-based (auto)" or "Achievement (permanent)"
  - Active holders count + trend
  - Integration status: "API Connected" / "Contract Deployed" / "Not Integrated" (with setup link)
- Expandable: usage chart (privilege applications per day), top users, budget usage

**Create Privilege flow**:
- Step 1: Type (select from: Fee Discount / Gas Rebate / Yield Boost / Priority Access / Custom)
- Step 2: Value (e.g., discount percent, rebate amount, boost percent)
- Step 3: Qualification mode
  - Status-based: select criteria (Level / Points Balance / Custom) + operator + value
  - Achievement-based: select trigger (Badge / Milestone) + specific item
- Step 4: Duration (While qualified / Permanent / Fixed days)
- Step 5: Limits (max holders, total budget)
- Step 6: Integration instructions (show API endpoint + code snippet + webhook setup)

---

## Part 5: C-End Widget Specifications

### Activity Reward Toast Widget

Embedded in project site via WL Widget Library. Shows real-time point rewards.

**Trigger**: Contract event detected for current user
**Display**: Slide-in toast notification (bottom-right or top-right)
**Content**:
```
+50 pts
Swap $250 USDC → ETH
━━━━━━ Total: 1,250 pts
```
**Duration**: 4 seconds, auto-dismiss. Stack up to 3 toasts.
**Combo**: If multiple actions within 30s, show "Combo x3!" animation.

### My Privileges Panel (in User Center Widget)

**Section in User Center widget showing active and locked privileges.**

Active privileges:
```
[green] 30% Fee Discount     Gold Level
        Applied automatically
        Saved $47.20 this month

[green] +2% Yield Boost      Diamond Badge
        Expires in 12 days
        Extra $18.50 earned
```

Locked privileges (goal gradient):
```
[gray/lock] Gas Rebate         Need Level 5
            [=======---] 78% to unlock

[gray/lock] Priority Access    Need Badge: Whale
            Complete 3 more milestones
```

Summary:
```
You saved $65.70 this month with your privileges.
```

### Privilege Status Widget (standalone, for Widget Library)

A compact card showing privilege summary, embeddable anywhere on project site.
```
Your Benefits        Gold Level
━━━━━━━━━━━━━━━━━━━
30% Fee Discount   [Active]
+2% Yield Boost    [12d left]
━━━━━━━━━━━━━━━━━━━
Saved $65.70 this month
[View All Benefits →]
```
