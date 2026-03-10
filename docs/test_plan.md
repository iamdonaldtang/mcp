# TaskOn Test Plan — Unit / Integration / E2E

## 1. Project Overview

- **Frontend**: Vue 3 + TypeScript + Tailwind CSS + Pinia + Vue Router
- **Backend**: Go + Gin + GORM + PostgreSQL + JWT
- **Total**: 57 Vue pages + 30+ Go handlers + 30+ DB models
- **All builds pass**: `vue-tsc --noEmit`, `vite build`, `go build`

## 1.1 Implementation Gaps (Must Fix Before Testing)

以下功能在页面实现阶段未完整实现，需在测试前补充：

### Gap-1: D20 Publish Readiness Check 未封装为共享组件
- **问题**: req doc 要求所有 Publish/Launch/Activate 按钮触发 D20 modal（检查订阅状态 + Twitter 授权），当前各页面的 Publish 按钮直接调用 API，未统一封装
- **影响页面**: B13WizardStep4 (Publish Community), B17WizardStep4 (Publish WL), B31b-B31h (Activate module instances), B31g (Publish Shop Item), B50Content (Publish Announcement 除外，不触发 D20)
- **修复方案**: 创建 `components/common/PublishReadinessCheck.vue` 共享组件，接收 `entityType` + `entityId` props，内部调用 `GET /api/v1/readiness?type=X&id=Y`，检查 2 项（订阅状态 + Twitter 授权），全部通过后 emit `confirmed` 事件。各页面的 Publish 按钮改为先打开此组件
- **缓存**: 5 分钟 sessionStorage 缓存，key = `{entityType}_{entityId}_readiness`

### Gap-2: D19 Promo Kit Generator 未完整实现
- **问题**: B10 (Community Guided) 和 B15 (WL Active) 的 "Announce" 检查步骤需要触发 D19 Promo Kit Generator modal，当前仅有占位按钮
- **影响页面**: B10CommunityGuided (go_live section → Announce step), B15WLActive (GO LIVE section → Announce step)
- **修复方案**: 创建 `components/common/PromoKitModal.vue`，包含平台选择 (Twitter/Discord/Telegram)、AI 生成文案（调用 `POST /api/v1/promo-kit/generate`）、生成横幅图、一键分享按钮

### Gap-3: 跨页面 Store 级联刷新缺失
- **问题**: 各页面独立获取数据，缺少 store 级别的级联刷新机制。例如：
  - B31a 创建任务后，B13WizardStep3 的模板数据不会同步
  - B-end 启用/禁用模块后，C-end 的 tab 可见性不会实时更新
  - B31g 发布商品后，C06Shop 的商品列表不会刷新
- **影响**: 多个页面间的数据一致性
- **修复方案**: 在 Pinia community store 中添加 `lastModified` 时间戳，各页面 onMounted 检查是否需要重新获取数据；或使用 store action 统一管理数据刷新事件

### Gap-4: WebSocket 实时更新仅占位
- **问题**: B15 的集成验证 ping 检测 (`/ws/wl/integration-ping`)、B10 的参与者实时计数器 (`/ws/community/participants`) 仅有模拟逻辑
- **影响页面**: B10CommunityGuided (First 10 Participants step), B15WLActive (Integration Verified step, First User Interaction step)
- **修复方案**: 实现 WebSocket 客户端工具类 `utils/websocket.ts`，封装重连逻辑 + 心跳，各页面通过 composable `useWebSocket(url)` 接入

### Gap-5: 级联删除/依赖检查不完整
- **问题**: 部分删除操作的依赖影响检查未完整实现：
  - 删除 Badge → 应检查是否被 Milestone/Shop/TaskChain 引用
  - 删除 Point Type → 应检查是否有 Leaderboard/LB Sprint/Level 依赖
  - 删除 Contract → 应检查是否有 Rule Builder 规则引用（B48 已实现此项）
  - 断开 Integration → 应检查是否有 Access Rule 依赖
- **影响页面**: B31iBadges, B32PointsLevel, B33AccessRules, B61Integration
- **修复方案**: 删除前调用 `GET /api/v1/xxx/:id/dependencies` 检查引用，有引用时在确认弹窗中列出影响项

### Gap-6: 表单自动保存与草稿恢复不一致
- **问题**: B13 Wizard 有 30s 自动保存 + 草稿恢复，但 B17 WL Wizard 部分步骤未实现相同机制
- **影响页面**: B17WizardStep2, B17WizardStep3
- **修复方案**: 统一 wizard 自动保存 composable `useWizardDraft(key, data)`，两个 wizard 复用

---

## 2. Test Stack (Recommended)

### Frontend
- **Unit Tests**: Vitest + @vue/test-utils
- **Component Tests**: Vitest + @testing-library/vue
- **E2E Tests**: Playwright (cross-browser)

### Backend
- **Unit Tests**: Go standard `testing` package
- **Integration Tests**: `httptest` + test DB (SQLite or test PostgreSQL)
- **E2E Tests**: Playwright (shared with frontend)

---

## 3. Frontend Unit Tests

### 3.1 Stores (4 files)
| Store | File | Key Tests |
|-------|------|-----------|
| auth | `stores/auth.ts` | login sets token + profile, logout clears state, fetchProfile on 401 clears token |
| community | `stores/community.ts` | fetchOverview populates state, checklist progress computed correctly |
| whitelabel | `stores/whitelabel.ts` | fetchOverview, checklist state |
| c-end | `stores/c-end.ts` | connectWallet sets token, tabVisibility reactivity, disconnectWallet clears |

### 3.2 API Client (1 file)
| Test | Description |
|------|-------------|
| `api/client.ts` | JWT interceptor adds Authorization header, 401 triggers logout, cApi uses separate token |

### 3.3 Shared Components (7 files)
| Component | Key Tests |
|-----------|-----------|
| StatusBadge | Renders correct color/text for each status (active/draft/completed/paused) |
| DataTable | Renders headers + rows, loading skeleton, empty state |
| Pagination | Page click emits event, disables prev/next at bounds |
| StatsCard | Renders value + label + trend arrow |
| EmptyState | Renders icon + text + CTA button |
| Modal | Opens/closes, emits close on overlay click, renders slots |
| TaskCard (C-end) | 7 status states render correctly (available/in_progress/completed/claimed/locked/expired/cooldown) |

### 3.4 Page Components — Priority P0 (Core Flows)
| Page | Key Tests |
|------|-----------|
| B09 CommunityEmpty | 3 strategy cards render, CTA navigates to wizard |
| B10 CommunityGuided | Checklist sections render, steps expand/collapse |
| B13 WizardStep1 | Form validation (name 3-50, desc 10-500), color picker, auto-save fires |
| B13 WizardStep2 | Strategy selection auto-enables modules, required modules locked |
| B13 WizardStep4 | Readiness checklist, publish button disabled when incomplete |
| B31a Sectors | CRUD: create task modal, status toggle, search debounce, pagination |
| B14 WLEmpty | 3 path cards, prerequisite banner when no community |
| B17 WizardStep2 | Path-adaptive rendering (5 variants) |
| C01 Home | DayChain check-in, task sectors render, quick actions |
| C06 Shop | Item grid states (affordable/not_enough/sold_out), redeem modal |

### 3.5 Page Components — Priority P1 (Module Management)
All B31b-B31i follow same pattern — test one thoroughly, spot-check others:
| Test Pattern | Description |
|-------------|-------------|
| Stats row | 4 StatsCard components render with API data |
| Filter tabs | Tab click updates filter, URL param sync |
| Search | Debounced 300ms, clears with X button |
| Data table | Columns match spec, row actions work |
| Create modal | Opens, validates fields, submits POST |
| Edit modal | Pre-fills data, submits PUT |
| Status toggle | Optimistic update, rollback on error |
| Delete | Confirmation dialog, DELETE request |
| Pagination | 20/page default, page change emits |
| Empty state | Shows when no data |

### 3.6 Page Components — Priority P2 (All Remaining)
- B01 Dashboard: 3 user states render correctly
- B05 Analytics: Date range selector, product cards
- B06 Settings: Tab switching, form saves
- B50 Content: Announcements CRUD, featured slots, module status
- B51 Preview: Desktop/Mobile toggle, click interception
- B60 Insights: Segment expansion, retention bars
- B61 Integration: Card states (available/connected/error), OAuth mock
- All WL pages: B40-B49 with their respective modals

---

## 4. Backend Unit Tests

### 4.1 Handler Tests (Priority Order)

#### P0 — Auth & Core
| Handler | File | Tests |
|---------|------|-------|
| Auth | `handler/auth.go` | Register (hash password, create user), Login (verify password, generate JWT), Profile (return user data), invalid credentials → 401 |
| Wallet | `handler/wallet.go` | Connect (upsert member, generate C-end JWT), invalid wallet → 400 |

#### P0 — Community Hub
| Handler | Tests |
|---------|-------|
| GetOverview | Returns stats, empty state for new user |
| Create | Creates community + default PointType + DayChainConfig |
| ListSectors | Pagination, filters |
| CreateSector / CreateTask | Validation, success |
| UpdateTask / DeleteTask | Status toggle, cascade |

#### P0 — C-End
| Handler | Tests |
|---------|-------|
| Home | Returns pulse stats + tab visibility |
| Tasks | Returns tasks with user completion overlay |
| DayChainCheckIn | Streak logic (new streak, continue, break, grace period) |
| Leaderboard | Returns ranked list with user position |
| ShopItems + Redeem | Points deduction, stock update, insufficient points |
| ClaimMilestone | Claim logic, duplicate claim prevention |

#### P1 — White Label
| Handler | Tests |
|---------|-------|
| WL GetOverview | Returns stats by state |
| WL Create/Update | Validation |
| Widget CRUD | List, Create, Update, Delete |
| Page CRUD | Create with slug validation, Publish |
| RewardRule CRUD | Trigger/action validation |

### 4.2 Middleware Tests
| Middleware | Tests |
|-----------|-------|
| BEndAuth | Valid JWT → sets user in context, expired → 401, missing → 401 |
| CEndAuth | Valid C-end JWT → sets wallet in context |
| CEndOptionalAuth | Missing token → continues (no error), valid → sets wallet |
| CORS | Returns correct headers |

### 4.3 Model Tests
| Test | Description |
|------|-------------|
| AutoMigrate | All 30+ models migrate without error |
| Relationships | Foreign keys, cascade deletes |
| Validations | Required fields, unique constraints |

---

## 5. Integration Tests (Backend)

Use `httptest.NewServer` + test PostgreSQL/SQLite.

### 5.1 Auth Flow
```
POST /api/v1/auth/register → 201
POST /api/v1/auth/login → 200 + token
GET /api/v1/auth/profile (with token) → 200
GET /api/v1/auth/profile (no token) → 401
```

### 5.2 Community CRUD Flow
```
POST /api/v1/community → 201 (creates community + defaults)
GET /api/v1/community/overview → 200
POST /api/v1/community/sectors → 201
POST /api/v1/community/tasks → 201
PUT /api/v1/community/tasks/:id → 200 (status toggle)
DELETE /api/v1/community/tasks/:id → 200
GET /api/v1/community/tasks?page=1&limit=20 → 200 (pagination)
```

### 5.3 C-End User Flow
```
POST /api/c/wallet/connect → 200 + c_token
GET /api/c/home → 200 (pulse + tabs)
GET /api/c/tasks → 200 (with completion overlay)
POST /api/c/daychain/checkin → 200 (streak increment)
POST /api/c/daychain/checkin → 400 (duplicate same day)
GET /api/c/leaderboard → 200
GET /api/c/shop → 200
POST /api/c/shop/:id/redeem → 200 (points deducted)
POST /api/c/shop/:id/redeem → 400 (insufficient points)
```

### 5.4 White Label Flow
```
POST /api/v1/whitelabel → 201
POST /api/v1/whitelabel/widgets → 201
PUT /api/v1/whitelabel/widgets/:id → 200
GET /api/v1/whitelabel/pages → 200
POST /api/v1/whitelabel/rules → 201
```

---

## 6. E2E Tests (Playwright)

### 6.1 Setup
- Dev server: `vite preview` or `vite dev`
- Backend: test instance with seeded data
- Browser: Chromium (primary), Firefox + WebKit (cross-browser)

### 6.2 Critical User Journeys

#### Journey 1: B-End Community Creation (P0)
```
1. Login → Dashboard (B01 new user state)
2. Click "Make Users Stay" → B09 Community Empty
3. Click "Create Your First Community" → B13 Wizard Step 1
4. Fill name + description + brand color → Next
5. Select strategy "Activate New Users" → verify modules auto-enabled → Next
6. Verify auto-generated content → edit a task name → Next
7. Verify preview renders → check readiness → Publish
8. Verify redirect to B10 Guided
```

#### Journey 2: B-End Module Management (P0)
```
1. Navigate to Sectors & Tasks (B31a)
2. Create sector → verify appears
3. Create task (fill modal) → verify in table
4. Toggle task status → verify optimistic update
5. Search → verify filter
6. Delete task → confirm → verify removed
```

#### Journey 3: C-End User Flow (P0)
```
1. Connect wallet → verify user status bar
2. Check-in DayChain → verify streak increment
3. Browse tasks → click "Start" → verify status change
4. Visit Leaderboard → verify ranking
5. Visit Shop → redeem item → verify points deducted
6. Visit Activity Feed → verify entries appear
```

#### Journey 4: White Label Setup (P1)
```
1. Navigate to WL Empty (B14)
2. Select "Embed" path → Start wizard
3. Complete 4 wizard steps
4. Verify redirect to B15 Active
5. Navigate to Widget Library → create widget
6. Navigate to Brand Settings → change colors → verify preview
```

#### Journey 5: Cross-Product Navigation (P1)
```
1. Dashboard → Community → WL → Analytics → Settings
2. Verify sidebar highlights correct item
3. Verify breadcrumbs update
4. Community sidebar: expand → navigate modules → verify sub-item highlight
5. WL sidebar: expand → navigate sub-items
```

### 6.3 Visual Regression (P2)
- Screenshot each page at 1440px + 375px
- Compare against baseline on PR

---

## 7. Test File Structure

```
frontend/
├── vitest.config.ts
├── src/
│   ├── __tests__/
│   │   ├── stores/
│   │   │   ├── auth.spec.ts
│   │   │   ├── community.spec.ts
│   │   │   ├── whitelabel.spec.ts
│   │   │   └── c-end.spec.ts
│   │   ├── components/
│   │   │   ├── StatusBadge.spec.ts
│   │   │   ├── DataTable.spec.ts
│   │   │   ├── Pagination.spec.ts
│   │   │   ├── Modal.spec.ts
│   │   │   └── TaskCard.spec.ts
│   │   ├── views/
│   │   │   ├── b-community/
│   │   │   │   ├── B09CommunityEmpty.spec.ts
│   │   │   │   ├── B13WizardStep1.spec.ts
│   │   │   │   ├── B31aSectors.spec.ts
│   │   │   │   └── ... (one per page)
│   │   │   ├── b-whitelabel/
│   │   │   │   ├── B14WLEmpty.spec.ts
│   │   │   │   └── ...
│   │   │   └── c-end/
│   │   │       ├── C01Home.spec.ts
│   │   │       └── ...
│   │   └── api/
│   │       └── client.spec.ts
│   └── ...
├── e2e/
│   ├── playwright.config.ts
│   ├── community-creation.spec.ts
│   ├── module-management.spec.ts
│   ├── c-end-user-flow.spec.ts
│   ├── whitelabel-setup.spec.ts
│   └── navigation.spec.ts

backend/
├── internal/
│   ├── handler/
│   │   ├── auth_test.go
│   │   ├── community_hub_test.go
│   │   ├── c_community_test.go
│   │   ├── wl_hub_test.go
│   │   └── wallet_test.go
│   ├── middleware/
│   │   └── auth_test.go
│   └── model/
│       └── models_test.go
└── test/
    └── integration/
        ├── auth_flow_test.go
        ├── community_flow_test.go
        ├── c_end_flow_test.go
        └── wl_flow_test.go
```

## 8. Implementation Priority

| Phase | Scope | Estimated Files | Priority |
|-------|-------|----------------|----------|
| T1 | Frontend test infra (vitest.config + test utils + mocks) | 3 | P0 |
| T2 | Store unit tests (4 stores) | 4 | P0 |
| T3 | Shared component tests (7 components) | 7 | P0 |
| T4 | Backend handler unit tests (auth + community + c-end) | 5 | P0 |
| T5 | Backend middleware tests | 1 | P0 |
| T6 | Backend integration tests (4 flows) | 4 | P0 |
| T7 | Frontend page tests P0 (10 core pages) | 10 | P0 |
| T8 | E2E critical journeys (5 journeys) | 5 | P0 |
| T9 | Frontend page tests P1 (module mgmt pattern) | 9 | P1 |
| T10 | Frontend page tests P2 (all remaining) | ~30 | P2 |
| T11 | E2E visual regression | 1 config + baselines | P2 |

**Recommended execution order**: T1 → T2 → T4 → T5 → T3 → T6 → T7 → T8 → T9 → T10 → T11
