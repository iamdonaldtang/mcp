# C端用户页面前端开发需求文档

> 版本: v1.0 | 日期: 2026-03-06
> 基于设计稿 `design/pencil-new.pen` + 现有文档 `website_frontend_requirements.md` v4.2
> 供前端/后端工程师对照实施

---

## 目录

1. [模块概述](#1-模块概述)
2. [全局架构](#2-全局架构)
3. [C01 Community Home](#3-c01-community-home)
4. [C02 Quest Tab](#4-c02-quest-tab)
5. [C03 Leaderboard](#5-c03-leaderboard)
6. [C04 LB Sprint](#6-c04-lb-sprint)
7. [C05 Milestone](#7-c05-milestone)
8. [C06 Shop](#8-c06-shop)
9. [C07 User Center](#9-c07-user-center)
10. [C08 Invite Center](#10-c08-invite-center)
11. [C09 Activity Feed](#11-c09-activity-feed)
12. [Tab 可见性规则](#12-tab-可见性规则)
13. [API 接口汇总](#13-api-接口汇总)
14. [钱包连接](#14-钱包连接)

---

## 1. 模块概述

### 1.1 产品定位

C端是 White Label Community 的用户前台，面向项目方的终端用户。域名由项目方自定义（如 `community.bitcoin.com`）。C端与B端是**完全独立的系统**——独立 URL、导航、用户模型。

### 1.2 页面编码总览

| 编码 | 页面 | Node ID | URL |
|------|------|---------|-----|
| C01 | Community Home | `vJVhd` | `/` |
| C02 | Quest Tab | `dUXTl` | `/quests` |
| C03 | Leaderboard | `KmdSd` | `/leaderboard` |
| C04 | LB Sprint | `y5fUZ` | `/lb-sprint` |
| C05 | Milestone | `53iKE` | `/milestones` |
| C06 | Shop | `coM7o` | `/shop` |
| C07 | User Center | `PykHF` | `/profile` |
| C08 | Invite Center | `TaAo9` | `/invite` |
| C09 | Activity Feed | `xhPIr` | `/activity` |

### 1.3 主题规格

| 属性 | 值 |
|------|-----|
| 页面背景 | `#0A0F1A` (深色) |
| Header 背景 | `#0F172A` |
| 卡片背景 | `#111B27` |
| 边框 | `#1E293B` |
| Primary 强调色 | `#F59E0B` (amber) |
| 主文本 | `#F1F5F9` |
| 次级文本 | `#94A3B8` |
| 字体 | Inter 全局 |
| 图标系统 | Material Symbols Rounded |
| 按钮主色 | `#F59E0B` fill, `#000000` text |
| 状态绿 | `#48BB78` (completed/claimed) |
| 状态灰 | `#64748B` (locked/disabled) |

---

## 2. 全局架构

### 2.1 页面布局

```
┌───────────────────────────────────────────────┐
│ Header: Logo + Project Name | Nav Tabs | Wallet + User │
├───────────────────────────────────────────────┤
│ User Status Bar (if logged in)                 │
├───────────────────────────────────────────────┤
│ Page Content (full width, padding varies)      │
├───────────────────────────────────────────────┤
│ Footer: "Powered by TaskOn"                    │
└───────────────────────────────────────────────┘
```

### 2.2 Header 组件

```
Header (#0F172A, height ~56px)
├── Left: Project Logo (from B-end brand) + Project Name
├── Center: Nav Tabs (horizontal, amber underline active)
│   ├── Home → /
│   ├── Quests → /quests
│   ├── Leaderboard → /leaderboard
│   ├── LB Sprint → /lb-sprint
│   ├── Milestone → /milestones
│   └── Shop → /shop
├── Right:
│   ├── Notification Bell (icon)
│   ├── "Connect Wallet" button (amber outline) — if not connected
│   └── User Avatar + Address (if connected) → /profile
│       └── Settings gear icon
```

**Nav Tab 样式**:
- Active: amber `#F59E0B` text + bottom underline 2px
- Inactive: `#94A3B8` text

### 2.3 User Status Bar

登录后在 Header 下方显示：

```
User Status Bar (#111B27, ~40px)
├── Avatar + Wallet Address (truncated)
├── Level Badge: "Lv.3" (amber)
├── Points: "1,250 XP"
├── Streak: "🔥 7 day streak"
└── Notification count
```

### 2.4 Footer

```
Footer (#0A0F1A, border-top #1E293B)
├── Left: "TaskOn · Earn rewards by completing Web3 tasks"
├── Center: Help Center | Docs | Community | Blog | Terms | Privacy
└── Right: Social icons (Twitter/Discord/Telegram)
├── "Powered by TaskOn" link → taskon.xyz (ext)
```

---

## 3. C01 — Community Home

**设计稿**: Node `vJVhd` | URL: `/`

### 3.1 页面概述
- **功能**: Action Engine — 用户的个人仪表盘 + 紧急操作 + 任务列表
- **核心目标**: 驱动即时行动（签到、完成任务、查看排名）

### 3.2 页面结构

```
Page Content (gap: 24px)
├── User Card (amber accent)
│   ├── Avatar + Wallet Address
│   ├── Stats: 1,249 XP | Level 3 | Daily Streak: 7
│   └── Quick badges: "Early Member" / "Active Weekly"
├── Announcement Carousel (if any)
│   └── Banner: "Token Launch Event — Double Points!"  + "Learn More" button
├── Daily Streak (DayChain)
│   ├── Title: "Daily Streak — Day 7 of 30"
│   ├── Calendar strip (30 days, completed=green, today=amber, future=gray)
│   ├── "Unlock Streak: 22 days to go, get bonus 5x at Day 30"
│   └── "Continue Streak" button (amber)
├── Quick Actions (horizontal scroll)
│   ├── "Start Quest" (amber icon box)
│   ├── "Lucky Wheel" (amber icon box)
│   └── "Invite to Top 10" (amber icon box)
├── "Getting Started" Sector
│   ├── Task 1: "Follow @Project on Twitter" — +50 XP [Start]
│   ├── Task 2: "Join Discord Server" — +50 XP [Start]
│   └── Task 3: "Complete KYC Verification" — +100 XP [Locked]
├── "Daily Engagement" Sector
│   ├── Task: "Daily Check-In" — +10 XP [Done]
│   ├── Task: "Share Daily Market Take" — +25 XP [Start]
│   └── Task: "Invite Friends, Earn..." — +150 XP [Invite Now]
├── Community Pulse (stats bar)
│   ├── 1,247 Members | 342 This Week | 89 Live Active | 5,680 Tasks Done
├── Discover More Section
│   ├── "Sprint Challenges" card → /lb-sprint
│   └── "Rewards Shop" card → /shop
└── Footer
```

### 3.3 Task Card 组件

```
Task Card (#111B27, rounded 12px, padding 16px)
├── Left: Icon box (colored bg, 40×40, rounded 10px)
│   └── Material icon (specific to task type)
├── Center:
│   ├── Task Name (14px 600 #F1F5F9)
│   ├── Subtitle: "Social · 1,350 completed" (12px #94A3B8)
├── Right:
│   ├── Points badge: "+50 XP" (12px 600 #F59E0B)
│   └── Action button: "Start" / "Done" / "Claim" / "Locked"
```

**Task 状态按钮**:
| 状态 | 按钮文本 | 样式 |
|------|---------|------|
| Available | "Start" | amber fill, dark text |
| In Progress | "Continue" | amber outline |
| Completed (unclaimed) | "Claim" | green fill |
| Claimed | "Done ✓" | gray, disabled |
| Locked | lock icon + requirement | gray, disabled |
| Expired | "Expired" | gray, disabled |
| Cooldown | "Available in 3h 22m" | gray, countdown |

### 3.4 数据模型

#### Home Page Aggregate

| 字段 | 类型 | API |
|------|------|-----|
| announcements | Announcement[] | `/api/c/community/announcements` |
| featured | FeaturedSlot[] | `/api/c/community/featured` |
| dayChain | DayChainState | `/api/c/community/daychain` |
| sectors | Sector[] (with tasks) | `/api/c/community/tasks` |
| userStatus | UserStatus | `/api/c/user/status` |
| communityPulse | PulseStats | `/api/c/community/home` |

#### Task 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 任务 ID |
| name | string | 任务名 |
| type | enum | `social` / `onchain` / `verification` / `custom` / `recurring` / `referral` |
| status | enum | `available` / `in_progress` / `completed` / `claimed` / `locked` / `expired` / `cooldown` |
| points | number | XP 奖励 |
| icon | string | Material icon name |
| iconColor | string | 图标颜色 |
| sectorName | string | 所属分区 |
| completions | number | 总完成次数 |
| requirement | string? | 锁定条件描述 |
| cooldownEnd | ISO8601? | 冷却结束时间 |

### 3.5 按钮路由

| 按钮 | 目标 |
|------|------|
| Announcement "Learn More" | → (ext) 公告详情 |
| Featured card click | → (ext) 活动详情 |
| DayChain day click | (API) `POST /api/c/community/daychain` 签到 |
| Task "Start" / "Claim" | (API) 开始/领取任务 |
| "Sprint Challenges" | → C04 `/lb-sprint` |
| "Rewards Shop" | → C06 `/shop` |

---

## 4. C02 — Quest Tab

**设计稿**: Node `dUXTl` | URL: `/quests`

### 4.1 页面结构

```
Page Content
├── Title: "All Quests" + Filter pills
├── Filter Pills: All(25) | Available(9) | Completed(10)
├── Quest Cards Grid (2×2)
│   ├── Quest 1: "Token Launch Social Quest" [Available]
│   │   ├── Image (banner)
│   │   ├── Description
│   │   ├── "67 tasks completed"
│   │   ├── Progress: "3 tasks · 1,247 participants"
│   │   └── "Start Quest" button (amber)
│   ├── Quest 2: "Refer a Friend" [Available]
│   │   └── Similar structure
│   ├── Quest 3: "Community Onboarding" [Completed] (green badge)
│   │   └── "19 tasks · All completed" + checkmark
│   └── Quest 4: "DeFi Trading..." [Ends in 3d]
│       └── Countdown badge + "Start Quest" button
├── Cross-sell Banners (bottom)
│   ├── "Boost your rank in the Sprint →" → C04
│   └── "Unlock Milestones with your points →" → C05
└── Footer
```

### 4.2 Quest Card 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | Quest ID |
| name | string | 名称 |
| description | string | 描述 |
| image | string | Banner 图片 URL |
| status | enum | `available` / `completed` / `in_progress` / `ended` |
| taskCount | number | 任务总数 |
| completedTasks | number | 已完成任务数 |
| participants | number | 参与人数 |
| endsAt | ISO8601? | 结束时间 |
| points | number | 总可获积分 |

### 4.3 按钮路由

| 按钮 | 目标 |
|------|------|
| Filter pills | (action) 前端筛选 |
| "Start Quest" | (API) 开始 quest |
| Quest card click | → `/quests/:id` (detail, future) |

---

## 5. C03 — Leaderboard

**设计稿**: Node `KmdSd` | URL: `/leaderboard`

### 5.1 关键区分
- **Leaderboard ≠ LB Sprint**: 此页为**周期性积分排行榜**（周/月/全部），**无额外激励**
- 支持多积分类型（EXP/GEM 等）

### 5.2 页面结构

```
Page Content
├── Title: "Leaderboard"
├── Point Type Selector: "Community Points ▼" (dropdown)
├── Time Filter: This Week | This Month | All Time
├── Podium (Top 3)
│   ├── #2: Avatar + 9,450 pts (left, shorter)
│   ├── #1: Avatar + 12,000 pts (center, tallest, amber crown)
│   └── #3: Avatar + 7,580 pts (right, shortest)
├── Your Rank Bar (amber highlight)
│   └── "Lv.3 User · #40 people · 1,250 pts"
├── Rankings Table
│   ├── Rank | User | Level | Points
│   ├── #4: ... Lv.7 · 6,892 pts · ↑4 since last week
│   ├── #5: ... Lv.6 · 5,450 pts
│   └── ...continue
├── Cross-sell Banners
│   ├── "Earn more points in the Sprint →" → C04
│   └── "Redeem rewards in the Shop →" → C06
└── Footer
```

### 5.3 Podium 组件

| 位置 | 高度 | 颜色 | 内容 |
|------|------|------|------|
| #1 (center) | 最高 | amber 背景 `#F59E0B` | 头像 + 积分 + 皇冠 |
| #2 (left) | 中等 | `#CBD5E1` 银 | 头像 + 积分 |
| #3 (right) | 最低 | `#CD7F32` 铜 | 头像 + 积分 |

### 5.4 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/leaderboard` | GET | 300s |
| Params: `pointType`, `period` (weekly/monthly/alltime) | | |

---

## 6. C04 — LB Sprint

**设计稿**: Node `y5fUZ` | URL: `/lb-sprint`

### 6.1 关键区分
- **LB Sprint = 限时排行榜竞赛**，有明确开始/结束日期
- 基于自定义积分类型（EXP/GEM），附带**非积分类激励**（NFT/Token/WL Spot）

### 6.2 页面结构

```
Page Content
├── Sprint Header
│   ├── Title: "Weekly LB Sprint"
│   ├── Countdown badge: "Ends in 5d 18h 27m" (amber)
│   └── Points earned: "500 earned this week" (large number)
├── Progress Bar (current progress)
├── Sprint Tasks (task cards with point badges)
│   ├── "Complete trades on DEX" [Completed] — +120 pts
│   ├── "Provide liquidity to any pool" [Completed] — +100 pts
│   └── "Bridge assets to Base network" [In Progress] — +150 pts
├── Sprint Rankings
│   ├── #1: user... — +160 pts [Challenge]
│   └── #2: user... — +150 pts [Beat]
├── Reward Tiers (visual progress track)
│   ├── Tier 1: $50 USDT (green, completed ✓)
│   ├── Tier 2: Complete at $1,000 (green, completed ✓)
│   └── Tier 3: NFT Badge — "Top 10 Traders" (locked, highest)
├── Past Sprints Section
│   ├── Sprint card 1: "4th completions · 100 pts earned" [Completed]
│   └── Sprint card 2: "250 completions · 432 pts earned" [Expired]
├── Cross-sell Banners
│   ├── "Spend your points in the Shop →" → C06
│   └── "Track your rank on the Leaderboard →" → C03
└── Footer
```

### 6.3 Sprint 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | Sprint ID |
| name | string | 名称 |
| status | enum | `active` / `ended` / `upcoming` |
| startsAt | ISO8601 | 开始时间 |
| endsAt | ISO8601 | 结束时间 |
| pointType | string | 积分类型 (EXP/GEM) |
| userPoints | number | 用户已获积分 |
| totalParticipants | number | 总参与者 |
| tasks | SprintTask[] | 竞赛任务 |
| rewardTiers | RewardTier[] | 奖励阶梯 |
| rankings | RankEntry[] | 排名 (top N) |
| userRank | number | 用户当前排名 |

#### RewardTier

| 字段 | 类型 | 说明 |
|------|------|------|
| threshold | number | 积分阈值 |
| reward | string | 奖品描述 |
| type | enum | `token` / `nft` / `whitelist` / `points` |
| status | enum | `earned` / `claimable` / `locked` |

### 6.4 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/lb-sprint/current` | GET | 60s |
| `/api/c/lb-sprint/history` | GET | 300s |

---

## 7. C05 — Milestone

**设计稿**: Node `53iKE` | URL: `/milestones`

### 7.1 页面结构

```
Page Content
├── Header
│   ├── Title: "Milestones"
│   ├── User Level: "Gold Member · 1,250 / 2,000 pts to next"
│   └── Total Progress Bar (green)
├── ── Earned Milestones (green border) ──
│   ├── Milestone 1: "Earn your first 100 points" [Claimed ✓] — checkmark
│   └── Milestone 2: "Reach 500 points" [Claimed ✓] — "Claimed: Exclusive NFT"
├── ── Claimable (amber border, pulsing) ──
│   └── Milestone 3: "Reach 1,500 total points" [Claim Reward] — amber button
├── ── Locked (gray border) ──
│   ├── Milestone 4: "Diamond Hands" — Requires Lv.7, hold 50k+ tokens...
│   │   └── Reward: "OG NFT + $100 in Points"
│   ├── Milestone 5: "Whale Status" — Requires 50,000 point...
│   │   └── Reward: "Exclusive Merch..."
│   └── Milestone 6: "Legend" — Requires 100k pts...
├── Cross-sell Banners
│   ├── "Join this week's Sprint to earn faster →" → C04
│   └── "Complete Quests for 50+ milestones points →" → C02
└── Footer
```

### 7.2 Milestone Card 状态

| 状态 | 边框 | 按钮 | 视觉 |
|------|------|------|------|
| Claimed | green `#48BB78` | "Claimed ✓" (disabled, green) | Green checkmark overlay |
| Claimable | amber `#F59E0B` | "Claim Reward" (amber, pulsing) | Amber glow |
| Locked | gray `#1E293B` | none (gray, lock icon) | Dimmed content |

### 7.3 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/milestones` | GET | 60s |
| `/api/c/milestones/:id/claim` | POST | N/A |

---

## 8. C06 — Shop

**设计稿**: Node `coM7o` | URL: `/shop`

### 8.1 页面结构

```
Page Content
├── Header
│   ├── Title: "Rewards Shop"
│   └── Balance badge: "1,250 BTC Points Available" (amber)
├── Category Filter: All Items | NFTs | Vouchers | Merch | Whitelist
├── Item Grid (2×3)
│   ├── Item 1: "Limited Community Badge" [Available]
│   │   ├── NFT Image (3 variants)
│   │   ├── "NFT · Limited Edition"
│   │   ├── Price: "500 pts" (amber)
│   │   ├── "Redeem" button (amber fill)
│   │   └── "01/50 remaining"
│   ├── Item 2: "Ultra-Rare Trading Skin" [In Stock]
│   │   └── Price: 2,000 pts → "Redeem" (amber)
│   ├── Item 3: "Priority Access Beta" [Low Stock]
│   │   └── Price: 2,000 pts → "Not Enough" (gray badge) → needs 750 more
│   ├── Item 4: "T-Shirt" [Available]
│   │   └── Price: 1,500 pts → "Not Enough" (gray disabled)
│   ├── Item 5: "Ethereum" [Available]
│   │   └── Price: 950 pts → "Redeem" (amber)
│   └── Item 6: "NFT" [Sold Out]
│       └── Price: 1,200 pts → "Sold Out" (gray disabled)
└── Footer
```

### 8.2 Item Card 状态

| 状态 | 按钮 | 说明 |
|------|------|------|
| Affordable | "Redeem" (amber fill) | 用户积分足够 |
| Not Enough | "Not Enough" (gray) + "needs X more" | 积分不足 |
| Sold Out | "Sold Out" (gray, strikethrough) | 库存为 0 |
| Time-limited | Countdown timer | 限时商品 |

### 8.3 Redeem 交互

1. 点击 "Redeem" → 确认弹窗: "Redeem {item} for {price} pts?"
2. 确认 → (API) `POST /api/c/shop/redeem`
3. 成功 → Toast "Redeemed!" + 余额更新
4. 失败 → Error toast

### 8.4 Item 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 商品 ID |
| name | string | 名称 |
| description | string | 描述 |
| image | string | 图片 URL |
| category | enum | `nft` / `voucher` / `merch` / `whitelist` |
| price | number | 积分价格 |
| stock | number | 剩余库存 (-1=unlimited) |
| totalRedemptions | number | 总兑换次数 |
| isTimeLimited | boolean | 是否限时 |
| expiresAt | ISO8601? | 过期时间 |
| status | enum | `available` / `sold_out` |

### 8.5 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/shop/items` | GET | 60s |
| `/api/c/shop/redeem` | POST | N/A |

---

## 9. C07 — User Center

**设计稿**: Node `PykHF` | URL: `/profile`

### 9.1 页面结构

```
Page Content
├── Profile Card (#111B27, large)
│   ├── Avatar (64×64, rounded) + Wallet Address
│   ├── "Joined Jan 2024 · Level 3 · 1,249 XP Lifetime"
│   └── Level Progress Bar
├── Stats Row (4 cards)
│   ├── Total Points: 1,250
│   ├── Tasks Done: 23
│   ├── Day Streak: 14
│   └── Rank: #42
├── "Achievements" Section
│   ├── Badge Grid (6 slots)
│   │   ├── "7-Day Streak" (earned, colored, green)
│   │   ├── "First Quest" (earned, colored, amber)
│   │   ├── "Addict" (earned, colored, green)
│   │   ├── "Bookmark" (locked, gray + lock icon)
│   │   ├── "Cup" (locked, gray)
│   │   └── "Legend" (locked, gray)
├── "Recent Activity" Section
│   ├── Activity 1: "Completed 'Follow Twitter' task" — +50 pts — 2h ago (green dot)
│   ├── Activity 2: "Reached Level 3" — Level Up! — 1h ago (amber dot)
│   ├── Activity 3: "Referred user 0xf5..." — +150 pts — 5h ago (purple dot)
│   └── Activity 4: "Won Lucky Wheel bonus" — +200 pts — 1 day ago (red dot)
├── "Referral Program" Section
│   ├── Stats: Total Referrals (12) / Bonus Points Earned (600) / Conversion Rate (67%)
│   └── "Invite More Friends" button → C08
└── Footer (full footer with links + social + "Powered by TaskOn")
```

### 9.2 Achievement Badge 组件

| 状态 | 视觉 |
|------|------|
| Earned | 彩色图标 + 名称 + 获取日期 |
| Locked | 灰色图标 + lock overlay + "?" 或条件提示 |

### 9.3 Activity Item 颜色编码

| 类型 | 圆点颜色 |
|------|---------|
| Task completion | green `#48BB78` |
| Level up | amber `#F59E0B` |
| Referral | purple `#9B7EE0` |
| Lucky Wheel | red `#EF4444` |
| Milestone | blue `#3B82F6` |
| Shop redeem | pink `#EC4899` |

### 9.4 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/user/profile` | GET | session |
| `/api/c/user/achievements` | GET | 300s |
| `/api/c/user/activity` | GET | 30s |
| `/api/c/user/referral-stats` | GET | 60s |

---

## 10. C08 — Invite Center

**设计稿**: Node `TaAo9` | URL: `/invite`

### 10.1 页面结构

```
Page Content
├── Hero Banner (#111B27, centered)
│   ├── Icon (share icon, amber)
│   ├── Title: "Invite Friends, Earn Rewards"
│   └── Subtitle: "Share your link and earn 30 points for every friend who signs up"
├── Referral Link Bar
│   ├── URL: https://bitswift.community/a/ref={code}
│   └── "Copy Link" button (amber)
├── Social Share Buttons (horizontal, centered)
│   ├── "Twitter" (blue button)
│   ├── "Discord" (purple button)
│   └── "Telegram" (blue button)
├── Stats Row (4 cards)
│   ├── Total Invites Sent: 24
│   ├── Successful Joins: 12
│   ├── Points Earned: 600
│   └── Conversion Rate: 50%
├── "Your Referrals" Section
│   ├── Table: Wallet / Status / Points Earned / Date
│   ├── 0x3a7b...123 | Joined (green) | +50 | Mar 4, 2026
│   ├── 0xcd9f8...3e67 | Joined (green) | +50 | Mar 1, 2026
│   └── 0x55d2...a963 | Pending (amber) | — | Mar 2, 2026
└── Footer
```

### 10.2 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| referralCode | string | 用户唯一邀请码 |
| referralUrl | string | 完整邀请链接 |
| totalInvites | number | 已发送邀请数 |
| successfulJoins | number | 成功加入数 |
| pointsEarned | number | 邀请获得积分 |
| conversionRate | number (%) | 转化率 |
| referrals | Referral[] | 邀请记录列表 |

#### Referral

| 字段 | 类型 |
|------|------|
| address | string (truncated) |
| status | enum: `joined` / `pending` |
| pointsEarned | number |
| date | ISO8601 |

### 10.3 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/invite/link` | GET | session |
| `/api/c/invite/leaderboard` | GET | 300s |
| `/api/c/invite/tiers` | GET | 3600s |
| `/api/c/user/referral-stats` | GET | 60s |

---

## 11. C09 — Activity Feed

**设计稿**: Node `xhPIr` | URL: `/activity`

### 11.1 页面结构

```
Page Content
├── Filter Chips: All | Tasks | Points | Rewards | Level Up | Invites
├── Live Indicator: "● Live" (green dot, auto-refresh 15s)
├── ── Today ──
│   ├── Feed Item: "Completed 'Follow Twitter' task in Getting Started" — +50 pts — 3 hrs ago
│   │   └── (green dot icon)
│   ├── Feed Item: "Reached Level 5 — unlocked Bronze Badge" — Level Up! — 5 hrs ago
│   │   └── (amber dot icon)
│   └── Feed Item: "Invited 0xcd7fbb...0e47 — 'How about this community?'" — +150 pts — 8 hrs ago
│       └── (purple dot icon)
├── ── Yesterday ──
│   ├── Feed Item: "Won 200 points from Lucky Wheel!" — +200 pts — 1 day ago
│   │   └── (red dot icon)
│   ├── Feed Item: "DayChain streak reaches 14 days — keep it going!" — 14 day streak — 1 day ago
│   │   └── (amber dot icon)
│   └── Feed Item: "Redeemed 'Lv.3 OG Role' from Benefits Shop for 500 pts" — -500 pts — 1 day ago
│       └── (pink dot icon)
├── Trending Now Section (if applicable)
│   ├── "Hot Streak" trending card
│   ├── "LB Sprint Rush" trending card
│   └── "New Reward" trending card
└── Footer
```

### 11.2 Feed Item 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 动态 ID |
| type | enum | `task` / `points` / `reward` / `level_up` / `invite` / `wheel` / `streak` / `shop` |
| title | string | 动态标题 |
| description | string? | 详细描述 |
| pointsDelta | number? | 积分变化 (+/-) |
| timestamp | ISO8601 | 时间 |
| iconColor | string | 颜色编码 (同 C07 activity 颜色) |

### 11.3 实时更新

- **Auto-refresh**: 每 15 秒自动拉取新动态
- **Live indicator**: 绿色圆点 + "Live" 文字
- **New item animation**: 新动态从顶部滑入

### 11.4 API

| Endpoint | Method | Cache |
|----------|--------|-------|
| `/api/c/activity/feed` | GET | 15s |
| `/api/c/activity/trending` | GET | 300s |

---

## 12. Tab 可见性规则

C端 Nav Tab 由 B端模块激活状态**动态控制**：

| C-End Tab | B-End 模块要求 | 显示条件 |
|-----------|---------------|---------|
| **Home** | (always) | 始终显示 |
| **Quests** | Active tasks | 至少 1 个已发布任务 |
| **Leaderboard** | Leaderboard module | 已启用且至少 1 条记录 |
| **LB Sprint** | LB Sprint module | 有活跃或刚结束的 Sprint |
| **Milestone** | Milestones module | 已启用且至少 1 个阈值 |
| **Shop** | Benefits Shop | 已启用且至少 1 个商品 |

**规则**: 模块已启用但内容为空 → Tab **隐藏**（避免空状态）。

**Hidden tabs**: User Center (C07)、Invite Center (C08)、Activity Feed (C09) 不在主导航，通过 Header 右侧用户头像 / 页面内链接访问。

---

## 13. API 接口汇总

### 13.1 C-End 完整 API

| Endpoint | Method | 页面 | Cache | Priority |
|----------|--------|------|-------|----------|
| `/api/c/community/home` | GET | C01 | 60s | P0 |
| `/api/c/community/announcements` | GET | C01 | 60s | P1 |
| `/api/c/community/featured` | GET | C01 | 60s | P1 |
| `/api/c/community/daychain` | GET/POST | C01 | session | P0 |
| `/api/c/community/tasks` | GET | C01 | 30s | P0 |
| `/api/c/quests` | GET | C02 | 30s | P0 |
| `/api/c/leaderboard` | GET | C03 | 300s | P1 |
| `/api/c/lb-sprint/current` | GET | C04 | 60s | P0 |
| `/api/c/lb-sprint/history` | GET | C04 | 300s | P1 |
| `/api/c/milestones` | GET | C05 | 60s | P0 |
| `/api/c/milestones/:id/claim` | POST | C05 | N/A | P0 |
| `/api/c/shop/items` | GET | C06 | 60s | P0 |
| `/api/c/shop/redeem` | POST | C06 | N/A | P0 |
| `/api/c/user/status` | GET | ALL | session | P0 |
| `/api/c/user/profile` | GET | C07 | session | P0 |
| `/api/c/user/achievements` | GET | C07 | 300s | P1 |
| `/api/c/user/activity` | GET | C07, C09 | 30s | P0 |
| `/api/c/user/referral-stats` | GET | C07, C08 | 60s | P1 |
| `/api/c/invite/link` | GET | C08 | session | P0 |
| `/api/c/invite/leaderboard` | GET | C08 | 300s | P1 |
| `/api/c/invite/tiers` | GET | C08 | 3600s | P2 |
| `/api/c/activity/feed` | GET | C09 | 15s | P0 |
| `/api/c/activity/trending` | GET | C09 | 300s | P1 |
| `/api/c/wallet/connect` | POST | Header | N/A | P0 |

---

## 14. 钱包连接

### 14.1 连接流程

1. 用户点击 "Connect Wallet" (amber outline button)
2. 弹出钱包选择 Modal
   - MetaMask / WalletConnect / Coinbase Wallet / OKX
3. 授权签名 → 获取 wallet address
4. API: `POST /api/c/wallet/connect` → 返回 session token
5. Header 更新: 显示 avatar + truncated address + level badge

### 14.2 登录状态

| 状态 | Header 显示 | 功能限制 |
|------|-----------|---------|
| 未连接 | "Connect Wallet" button | 只能查看，不能操作任务/兑换 |
| 已连接 | Avatar + Address + Level | 全部功能可用 |

### 14.3 Session 管理

- JWT token 存储在 localStorage
- Token 过期: 7 天
- 自动续期: 每次 API 调用时检查

---

## 附录 A: 设计稿 Node ID 索引

| 页面 | Node ID |
|------|---------|
| C01 Community Home | `vJVhd` |
| C02 Quest Tab | `dUXTl` |
| C03 Leaderboard | `KmdSd` |
| C04 LB Sprint | `y5fUZ` |
| C05 Milestone | `53iKE` |
| C06 Shop | `coM7o` |
| C07 User Center | `PykHF` |
| C08 Invite Center | `TaAo9` |
| C09 Activity Feed | `xhPIr` |

---

## 附录 B: Home Tab 内容排序

C01 Home 页内容固定排序：

1. Announcement carousel (if any)
2. Featured grid 2×3 (if any)
3. User's active DayChain (if enabled)
4. User's active TaskChains (if any)
5. Sectors with tasks (in B-end configured order)
6. Community Pulse stats
7. Discover More section

---

## 附录 C: Edge Cases

| 场景 | C-End 显示 |
|------|-----------|
| Community paused | "This community is taking a break. Check back later." |
| Community deleted | 404 page |
| User blacklisted | "Access to this community is restricted." |
| User doesn't meet eligibility | Gate page with requirements |
| No modules active | Minimal page + "Coming soon" |
| Task deadline passed | Auto-remove from active lists |
| Task max claims reached | "Fully claimed" (grayed out) |
| Shop item sold out | "Sold Out" badge, button disabled |
| LB Sprint ended, rewards pending | Final standings + "Rewards distributing..." |
| DayChain streak broken | Red highlight + catch-up button (if enabled) |
| Leaderboard 0 participants | Tab hidden |
| Widget deactivated (WL embed) | "This widget is currently unavailable" |
| SSO auth fails (WL) | Fallback to TaskOn default login |
