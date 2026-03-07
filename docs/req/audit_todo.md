# 需求文档缺失内容审计清单

> 审计日期: 2026-03-06
> 审计对象: `req_community.md` + `req_white_label.md`
> 审计方法: 逐页对照设计稿 Node ID / 已有文档结构 / 标准 SaaS 交互规范，识别缺失内容
> 优先级: P0=阻塞开发 / P1=影响体验 / P2=完善度

---

## 一、缺失内容分类说明

每个页面缺失内容可归为以下类别：

| 类别 | 含义 |
|------|------|
| **OPS** | 区块级操作规格（每个按钮/输入/开关点击后发生什么）|
| **VALID** | 表单字段验证规则（必填/格式/范围/跨字段依赖）|
| **STATE** | UI 状态规格（loading/empty/error/disabled）|
| **MODAL** | Modal 完整字段定义与交互（现有文档只有名称无内容）|
| **CASCADE** | 跨页面/跨模块联动（操作 A 影响页面 B 的表现）|
| **EDGE** | 边界条件（数据为空时、权限不足、网络失败）|

---

## 二、req_community.md 缺失清单

### Hub 页面

#### B09 — Community Hub Empty

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-01 | OPS | 策略卡片是否有默认选中项（第一张预选）；已选中再次点击是否取消选中 | P0 |
| C-02 | OPS | CTA 按钮在无策略选中时是否 disabled；disabled 状态文案 | P0 |
| C-03 | OPS | Engine Strip 4 步（Quest→Activate→Engage→Retain）悬停时是否有 tooltip | P2 |
| C-04 | STATE | 页面数据加载中的 skeleton 规格 | P1 |
| C-05 | EDGE | 若项目已有 Community（导航后退场景），此页是否显示？路由守卫逻辑 | P0 |
| C-06 | OPS | Resource 卡片点击：Video Tutorial / Playbook 具体跳转 URL（是否固定或从 CMS 取）| P2 |

#### B10 — Community Hub Guided

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-07 | OPS | Checklist 项展开/折叠：点击整行还是箭头图标触发；是否多项同时展开 | P0 |
| C-08 | OPS | "Add more tasks" 展开后内容：显示什么（任务列表摘要+按钮）？按钮点击目标 | P0 |
| C-09 | OPS | "Set up Benefits Shop" 展开后内容 | P0 |
| C-10 | OPS | "Customize DayChain rewards" 展开后内容 | P0 |
| C-11 | OPS | "Preview" 展开后内容：直接跳转 B33 还是展开内嵌预览 | P0 |
| C-12 | OPS | Share 步骤 Twitter 分享：预填文案格式（含 @项目名、链接、hashtag 的模板） | P1 |
| C-13 | OPS | Share 步骤 Discord/Telegram：复制链接 vs 打开 Discord 客户端 | P1 |
| C-14 | OPS | "+ Enable" 模块：点击后 loading 状态、成功 toast 文案、失败回滚 | P0 |
| C-15 | STATE | "First 10 participants" 步骤：WebSocket 断连重试逻辑；0/10 → 5/10 → 10/10 计数动画 | P1 |
| C-16 | OPS | Checklist 步骤自动标记完成：哪些步骤由系统自动检测（vs 用户手动点击）| P0 |
| C-17 | OPS | Progress bar：公式（completed/total）；进度更新触发时机 | P0 |
| C-18 | CASCADE | 完成 Checklist 所有步骤后页面状态变化：自动切换到 B11 Active 还是需要刷新 | P0 |
| C-19 | STATE | "Browse Configuration Templates →" 点击跳转 B13 时携带什么参数 | P1 |
| C-20 | EDGE | Checklist 进入 B10 时，已完成步骤数据从哪里取（API 还是本地 localStorage）| P0 |

#### B11 — Community Hub Active

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-21 | OPS | Stats 卡片趋势箭头 hover tooltip 文案（"vs last 7 days"）及格式 | P2 |
| C-22 | OPS | Checklist Banner：点击行为（展开 in-place accordion 还是导航到 B10）| P0 |
| C-23 | OPS | Checklist Banner "×" 关闭：session-only 还是持久关闭（存 localStorage 还是 API）| P1 |
| C-24 | OPS | Module Performance 卡片：每种模块类型显示哪 2-3 个关键指标（各模块指标不同）| P0 |
| C-25 | OPS | Module Performance 卡片点击 card body（非按钮）是否也跳转模块页 | P1 |
| C-26 | EDGE | Quick Stats 数据刷新策略：手动刷新按钮、自动轮询间隔、WebSocket 推送 | P1 |

#### B12 — Community Hub Deep

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-27 | OPS | AI Insights 卡片单条 dismiss：session-only 还是持久；"Dismiss All" 操作 | P1 |
| C-28 | OPS | AI Insights 刷新：页面加载取最新；缓存时间；是否有手动刷新 | P1 |
| C-29 | OPS | AI Insights action link（如 "View DayChain Cliff →"）目标页面及锚点 | P1 |
| C-30 | OPS | Analytics chart (Weekly Active Users) hover tooltip 格式；点击 bar 跳转 B54 并选中对应日期 | P1 |
| C-31 | OPS | Retention Metrics "Full Analytics →" 跳转 B54 时是否带 filter 参数（如 retention tab 预选）| P1 |
| C-32 | OPS | Integration 卡片 error 状态（红色边框）："Fix Connection" 按钮行为 | P0 |
| C-33 | CASCADE | 在 B12 点击 "Configure" 某集成 → 跳转 B61 并自动滚动到对应集成项 | P1 |

---

### 创建向导

#### B13 — Step 1: Customize

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-34 | VALID | Community Name：实时字符计数显示位置（输入框右下角）；3-50 字符违反时红色 border + 错误文案 | P0 |
| C-35 | VALID | Description：实时字符计数；10-500 字符；允许换行？富文本还是 plain text | P0 |
| C-36 | OPS | Brand Color 预设色块：选中后显示白色圆环 ring；"Custom" 选项展开后 hex 输入框格式（# + 6位）| P0 |
| C-37 | OPS | Brand Color 自定义：hex 输入实时验证 + 颜色预览更新；非法值红框 | P0 |
| C-38 | OPS | 右侧 Preview 实时更新：哪些字段变化触发预览刷新（name/color/description？）| P0 |
| C-39 | OPS | "Save Draft" 成功：toast "Draft saved"；失败：toast "Failed to save, please try again" | P0 |
| C-40 | OPS | 自动保存：每 30s 若表单 dirty 则静默保存（无 toast）；保存中 topbar save 按钮 loading | P1 |
| C-41 | OPS | "Next: Modules" 点击：先执行前端验证 → 失败则 scroll to 第一个错误字段并聚焦；通过则保存 draft + 跳转 B34 | P0 |
| C-42 | EDGE | 浏览器刷新/关闭：提示 "You have unsaved changes. Leave?" 确认框 | P1 |
| C-43 | STATE | 若已有 draft，进入 B13 时是否恢复 draft 数据；提示 "Resume your draft?" | P1 |

#### B34 — Step 2: Configure Modules

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-44 | OPS | 模块 toggle 动画：200ms ease，green fill + slide；toggle 关闭：gray + slide back | P1 |
| C-45 | OPS | 模块行展开（点击非 toggle 区域）：展开显示 C-end 效果描述 + 小图，再次点击收起 | P0 |
| C-46 | OPS | Required 模块 toggle：禁用状态样式（opacity 0.5 + cursor not-allowed）；hover tooltip "This module is required" | P0 |
| C-47 | OPS | 切换模块 toggle → 右侧 Summary 面板：模块列表即时增减；"Estimated Points Earned" 重算规则 | P0 |
| C-48 | OPS | "Save Draft" 同 B13；"Back" 跳回 B13 并保留 B13 表单数据（wizard 共享状态）| P0 |
| C-49 | OPS | 策略来源显示：顶部 "Based on: Activate New Users strategy" + 修改说明 | P1 |
| C-50 | EDGE | 所有可选模块全部关闭（只剩 Required）：是否允许？右侧 Summary 最小配置提示 | P1 |

#### B35 — Step 3: Quick Setup

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-51 | OPS | 模块卡片默认展开还是折叠；若默认折叠则有哪些摘要信息可见 | P0 |
| C-52 | OPS | 内联编辑：任务名称点击 → contenteditable（Enter 保存，Esc 取消）；编辑中蓝色 border | P0 |
| C-53 | OPS | 内联编辑：XP 值点击 → number input（min 1, max 10000）；保存到 wizard 本地状态 | P0 |
| C-54 | VALID | 任务名称：不允许空白；XP 值：必须为正整数 | P0 |
| C-55 | OPS | "Edit after setup →" 链接样式（text link, purple）；点击跳转行为（对应模块管理页？还是新标签？）| P1 |
| C-56 | OPS | "+ Add Task" 内联按钮（如果有）：向列表末尾追加空白任务条目进行编辑 | P1 |
| C-57 | STATE | 如果 Step 2 未选择某模块，Step 3 对应卡片不显示（动态列表）| P0 |

#### B55 — Step 4: Preview & Publish

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-58 | OPS | C-end preview 内部 tab 点击（Home/Quests/Leaderboard）：在 iframe/mock 内切换，不触发 B端路由 | P0 |
| C-59 | OPS | Desktop/Mobile toggle：Mobile 宽度 375px 居中显示；Desktop 100% 宽度 | P0 |
| C-60 | OPS | Community URL slug 生成规则：基于 Community Name 自动生成（lowercase + hyphen）；重复时加数字后缀 | P0 |
| C-61 | OPS | "Copy" URL 按钮：clipboard API → 成功 toast "Copied!"；失败 fallback（选中文本）| P0 |
| C-62 | OPS | Readiness Checklist 每项自动检测：具体 API / 检测条件（哪些是自动通过，哪些需要数据存在）| P0 |
| C-63 | OPS | "Publish Community" 完整流程：① 前端 readiness check → ② 打开 D20 → ③ D20 全通过 → ④ POST publish API → ⑤ loading → ⑥ 成功跳转 B10 | P0 |
| C-64 | EDGE | Publish API 失败：关闭 D20 → toast "Publish failed: [error]" → 停留在 B55 | P0 |
| C-65 | STATE | 所有 readiness 项通过时，"Publish Community" 按钮高亮激活 | P0 |

---

### 模块管理页（通用规格缺失）

#### 通用表格页面共性缺失

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-66 | OPS | 表格行 hover 状态：背景色变化（`#161F2E`）；Actions 列是否始终显示还是 hover 时显示 | P1 |
| C-67 | OPS | 表格行点击（click anywhere on row）→ 打开对应 edit modal；Actions 列独立按钮不受此影响 | P0 |
| C-68 | OPS | "Duplicate" 操作：POST duplicate API → 新 draft 出现在列表顶部 → toast "[Name] duplicated as draft" | P0 |
| C-69 | OPS | "Delete" 操作：确认弹窗（"Delete '[name]'? This cannot be undone."）→ DELETE API → toast "[Name] deleted" → 行消失 | P0 |
| C-70 | OPS | Status badge 行内切换：直接点击 badge 或使用 toggle → PUT API → 乐观更新（先改 UI，失败回滚）| P0 |
| C-71 | OPS | Search bar：debounce 300ms；清空按钮（×）出现条件（有输入时）；无结果空状态 | P0 |
| C-72 | OPS | Filter tabs：URL 参数同步（?status=active）；刷新后保持筛选状态 | P1 |
| C-73 | OPS | Sort dropdown：选项（Date Created/Modified/Name/Status/关键指标）；默认排序；选中状态（checkmark）| P0 |
| C-74 | OPS | Pagination：每页条数（默认 20）；"Showing X-Y of Z items"；首页/末页时 Prev/Next 禁用 | P0 |
| C-75 | OPS | Insight Banner dismiss：× 按钮 → session-only 消失（不持久）；banner 含 action link 时样式 | P1 |
| C-76 | STATE | 列表为空时（0 条）：empty state 插画 + 提示文案 + primary CTA（如 "+ Create First [Item]"）| P0 |
| C-77 | STATE | 数据加载中：表格区域 skeleton（3-5 行灰色占位）| P1 |
| C-78 | STATE | API 加载失败：错误提示 + "Retry" 按钮 | P0 |
| C-79 | OPS | Bulk 操作：行前 checkbox（全选 checkbox 在表头）；选中后顶部出现 bulk action bar（Activate / Pause / Delete）| P1 |

#### B31 — Sectors & Tasks（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-80 | OPS | Sector 拖拽排序：拖拽手柄图标位置；拖拽时整个 sector 块高亮 + 目标位置虚线；松开触发 PUT /api/community/sectors/reorder | P0 |
| C-81 | OPS | Task 在同 Sector 内拖拽排序 | P0 |
| C-82 | OPS | Task 跨 Sector 拖拽：拖拽到另一 Sector 下方 → PUT {taskId, newSectorId, order} | P0 |
| C-83 | OPS | Sector header ⋮ 菜单：Edit Name（inline edit）/ Hide（sector 整体 hidden，tasks 不影响）/ Delete（需确认，且 sector 内无任务时才可删）| P0 |
| C-84 | OPS | Sector 删除非空时：提示 "Move or delete all tasks in this sector first" | P0 |
| C-85 | OPS | Task 行 Actions：✏️ 编辑（modal）/ ⊕ 复制（draft）/ 显隐切换（眼睛图标）/ 🗑️删除 | P0 |
| C-86 | OPS | Task 积分值内联编辑：双击 points cell → number input；验证 1-10000；Enter 保存 PUT；Esc 取消 | P1 |
| C-87 | OPS | "+ New Sector"：底部追加空白 sector 名称输入框；Enter 提交 POST；Esc 取消 | P0 |
| C-88 | OPS | "+ New Task"：全局按钮打开 task creation modal（vs 每个 sector 下有独立 "+ Add Task"）| P0 |
| C-89 | OPS | Task "Publish Task" 按钮（仅 Draft 状态任务有）→ D20 → 通过后 PUT status=active | P0 |
| C-90 | EDGE | 删除有完成记录的任务：是否阻止删除？还是软删除（存档）？ | P0 |

#### B31a — Points & Level（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-91 | OPS | 自定义积分类型管理：页面顶部 "Point Types" tab 或 section；+ Add Type（名称 + 符号）| P0 |
| C-92 | OPS | 积分获取规则配置：每种任务类型默认获得多少积分；是否有单独配置页面 | P0 |
| C-93 | VALID | 新增等级时：threshold 必须大于前一等级 threshold；保存时验证整体顺序 | P0 |
| C-94 | OPS | 修改等级 threshold：若有用户已在该等级，降低 threshold 可能导致用户降级 → 警告弹窗 | P0 |
| C-95 | OPS | 删除等级：若有成员在该等级，提示将降级到下一低等级或保持当前状态 | P0 |

#### B31b — TaskChain（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-96 | OPS | Funnel 图 hover：每步 tooltip 显示 "Step N: X completed (Y% of step N-1)" | P1 |
| C-97 | OPS | Chain 状态切换 Active → Pause：正在进行链的用户进度如何处理（保留/冻结/取消）| P0 |
| C-98 | OPS | "Activate Chain" → D20 → 通过 → Chain 开始对新用户生效（已有进度用户不变）| P0 |

#### B31c — DayChain（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-99 | OPS | Streak 分布图 hover tooltip："X users have maintained Y-day streak" | P1 |
| C-100 | OPS | "Day 7 cliff" 红色高亮区域：点击跳转 DayChain 配置 modal（D03）并聚焦 Day 7 奖励字段 | P1 |
| C-101 | OPS | DayChain 配置：Grace Period（断链容忍期）字段；Streak Freeze（道具）是否支持 | P1 |
| C-102 | EDGE | 多个 DayChain 是否支持？（目前设计是否有此能力）| P0 |

#### B31d — Leaderboard（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-103 | OPS | "Point Type" 下拉：数据来源（从 B31a 取已配置的积分类型列表）| P0 |
| C-104 | OPS | Archive 操作：C-end Leaderboard tab 对应条目隐藏；历史数据保留 | P0 |
| C-105 | OPS | 同一积分类型能否创建多个 Leaderboard（不同 Period）？最大数量限制 | P1 |

#### B31e — LB Sprint（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-106 | OPS | "End Early" 按钮：确认弹窗 "End this Sprint now? Winners will be calculated based on current rankings." → 执行 → 奖励自动分发 | P0 |
| C-107 | OPS | Completed Sprint："View Results" 按钮 → 弹窗显示最终排名 + 奖励发放状态 | P0 |
| C-108 | OPS | 奖励发放：自动还是手动；Token/NFT 奖励发放方式 | P0 |
| C-109 | OPS | Sprint 未来开始：scheduling 功能（设置 start date 后 status=Scheduled）| P1 |

#### B31f — Milestones（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-110 | OPS | 展开 Milestone 行：显示各 Tier 详情（Tier 1: 100pts → Badge A；Tier 2: 500pts → Shop Item B）| P0 |
| C-111 | OPS | Claim Rate 数字点击 → D18 Segment Detail Panel 显示该 Milestone 已/未领取用户 | P1 |
| C-112 | OPS | 修改已激活 Milestone 的 threshold：是否允许？影响已达标但未领取的用户 | P0 |

#### B31g — Benefits Shop（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-113 | OPS | Sold Out 商品行："Replenish" 按钮 → inline number input → 输入数量 → PUT update stock → toast | P0 |
| C-114 | OPS | 商品可用性门控：Level-gated（需等级 ≥ N）/ Badge-gated（需持有 badge X）的配置和显示 | P0 |
| C-115 | OPS | 商品图片：表格中显示缩略图；点击预览大图 | P1 |
| C-116 | EDGE | 库存为 0 时是否自动暂停销售；是否发通知给 B端管理员 | P1 |

#### B31h — Lucky Wheel（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-117 | OPS | 奖品概率配置：每个奖品 % 输入；实时计算总和（必须 = 100%）；不等于 100% 时 Save 按钮 disabled + 提示 | P0 |
| C-118 | OPS | 转盘预览：D08 Modal 中是否有可视化转盘展示概率分布 | P1 |
| C-119 | OPS | Spin Cost 编辑：在 Modal 内（非内联）| P0 |
| C-120 | EDGE | 奖品全部为 "Nothing" 时是否允许激活 | P1 |

#### B31i — Badges（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-121 | OPS | Badges 表中 "Holders" 数字：点击 → D18 Segment Detail Panel 显示持有该 Badge 的用户列表 | P1 |
| C-122 | OPS | Badge 图标：表格中显示 32×32 图标预览；Modal 中支持上传图片（PNG/SVG max 1MB）或选预设 | P0 |
| C-123 | OPS | 手动发放 Badge：管理员可从 D09 Modal 内手动给特定用户发放（输入 wallet address）| P1 |

#### B49 — Access Rules（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-124 | OPS | 规则优先级排序：拖拽手柄调整规则执行顺序（从上到下依次评估）| P1 |
| C-125 | OPS | "Preview Rule" 功能：输入 wallet address → 模拟该 wallet 经过所有规则判断的结果 | P1 |
| C-126 | CASCADE | Token Gate 规则依赖 B61 Integration Center 中的链上验证集成；若未配置则规则无法激活 | P0 |

#### B50 — Homepage Editor（独特操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-127 | OPS | 拖拽行排序 → PUT /api/community/settings/homepage/reorder；排序即时生效（无单独发布步骤）| P0 |
| C-128 | OPS | Visibility toggle 行内切换 → PUT → C-end 即时反映（section 消失或出现）| P0 |
| C-129 | OPS | "Preview" 按钮 → B33 Preview Mode（打开新标签还是同标签跳转）| P1 |
| C-130 | OPS | 允许的 Section 类型及各类型的配置内容（Banner/Quest Widget/Leaderboard Widget/Text/Custom HTML）| P0 |

---

### 运营辅助页

#### B32 — Content Management（缺失操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-131 | OPS | Announcement "Pinned" 状态：pin/unpin 操作；被 pin 的公告始终显示在 carousel 顶部 | P0 |
| C-132 | OPS | Announcement "Scheduled" 状态：显示发布时间；时间到自动发布（服务端定时任务）| P0 |
| C-133 | OPS | Featured Slot "+ Add Featured"：弹出选择 UI（选择现有 Quest/Sprint/Milestone 或输入外部 URL）| P0 |
| C-134 | OPS | Featured Slot 已填充：显示内容摘要 + "× 移除" + "编辑" 图标 | P0 |
| C-135 | OPS | Module Status Overview "Configure" 链接：分别跳转 B31a-B31h | P0 |
| C-136 | OPS | Announcement "Publish" → D20（公告发布是否也走 D20 检查？）| P1 |

#### B33 — Preview Mode（缺失规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-137 | OPS | Mock 用户身份固定为 "Preview User"；Level 5；积分 1,250 | P1 |
| C-138 | OPS | C-end preview 内的链接/按钮点击：拦截所有点击（不执行真实操作），仅模拟 UI 状态变化 | P0 |
| C-139 | OPS | "Exit Preview" → 返回来源页面（B32 Content Mgmt 或 B50 Homepage Editor）| P0 |
| C-140 | STATE | Preview Banner 不可关闭（non-dismissible）| P0 |

#### B54 — Community Insights（缺失操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-141 | OPS | Date Picker：预设 7天/30天/90天/自定义范围；自定义使用 calendar 双日期选择 | P1 |
| C-142 | OPS | Module Filter：多选 dropdown（可选 All/Tasks/Points/Leaderboard/DayChain 等）| P1 |
| C-143 | OPS | Economy Chart hover tooltip："EXP Earned: 15,840 / Burned: 12,300 / Net: +3,540" | P1 |
| C-144 | OPS | User Segment 卡片点击 → D18 Segment Detail Panel（显示该分群用户列表）| P0 |
| C-145 | OPS | "Export CSV" → 后台生成 → 下载触发（是否异步？进度提示？）| P1 |
| C-146 | OPS | "Export PDF" → 同上，PDF 格式含图表截图 | P1 |
| C-147 | OPS | Retention by Module 条形图：hover 显示精确数值；点击条形 → 跳转对应模块管理页 | P1 |

#### B61 — Community Integration Center（缺失操作）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| C-148 | OPS | Twitter "Connect"：OAuth 2.0 流程（popup window）→ 授权后回调 → status 变 Connected | P0 |
| C-149 | OPS | Discord "Connect"：OAuth 流程；Bot 邀请链接 | P0 |
| C-150 | OPS | Telegram "Connect"：Bot Token 配置流程 | P0 |
| C-151 | OPS | Blockchain "Connect"：选择支持的链（EVM 兼容）→ 输入 RPC endpoint 或使用默认 | P0 |
| C-152 | OPS | "Configure" 已连接集成：展示当前配置 + 修改选项 + "Disconnect" 按钮 | P0 |
| C-153 | OPS | 集成 error 状态：显示错误原因（Token expired / RPC unreachable 等）+ "Reconnect" 按钮 | P0 |

---

### Modal 完整字段规格（全部缺失）

#### D01 — Points & Level Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-01 | Level Name | text input | 必填, 1-30 字符, 唯一 | "Bronze" / "Silver" 等 |
| M-02 | Threshold | number input | 必填, ≥0, 必须 > 前一等级 threshold | 达到该积分值即升级 |
| M-03 | Point Type | select | 必填, 从已配置类型选 | EXP/GEM/自定义 |
| M-04 | Level Badge | icon picker | 可选 | 升级时显示的图标 |
| M-05 | Level Perks | textarea | 可选, max 200 字符 | 等级特权描述文本 |
| M-06 | 保存行为 | — | — | POST(新建)/PUT(编辑)；成功关闭 modal + 刷新列表 |

#### D02 — TaskChain Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-07 | Chain Name | text input | 必填, 1-50 字符 | |
| M-08 | Steps | dynamic list | 至少 2 步 | 每步选择现有任务（多选）|
| M-09 | Step 任务选择 | multi-select | 每步至少 1 个任务 | 从 B31 已发布任务中选 |
| M-10 | Completion Reward | 积分 + badge | 可选 | 全部步骤完成后的奖励 |
| M-11 | 步骤顺序 | drag handle | — | 拖拽调整步骤顺序 |

#### D03 — DayChain Config

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-12 | Chain Name | text input | 必填 | |
| M-13 | Daily Task | task select | 必填 | 选择每日需完成的任务 |
| M-14 | Base Reward | number | 必填, ≥1 | 每日基础积分 |
| M-15 | Milestone Bonuses | list: day + multiplier | 可选 | Day 7=2x, Day 14=3x, Day 30=5x |
| M-16 | Grace Period | number | 0-24, 单位 hours | 断链容忍时长 |

#### D04 — Leaderboard Config

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-17 | Name | text input | 必填 | |
| M-18 | Point Type | select | 必填 | EXP/GEM/自定义 |
| M-19 | Period | radio | 必填 | Weekly / Monthly / All Time |
| M-20 | Display Top N | number | 默认 100 | C端展示排行榜条目数 |
| M-21 | 重置时间 | — | — | Weekly=周一 00:00 UTC; Monthly=1日 00:00 UTC |

#### D05 — LB Sprint Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-22 | Sprint Name | text input | 必填 | |
| M-23 | Point Type | select | 必填 | 基于哪种积分排名 |
| M-24 | Start Date | date picker | 必填, ≥ today | |
| M-25 | End Date | date picker | 必填, > Start Date | |
| M-26 | Reward Tiers | dynamic list | 至少 1 tier | 每 tier: 排名范围 + 奖励类型(Token/NFT/WL Spot) + 数量 |
| M-27 | Reward 发放方式 | radio | 必填 | Auto (结束后自动) / Manual (需管理员触发) |

#### D06 — Milestone Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-28 | Milestone Name | text input | 必填 | |
| M-29 | Tiers | dynamic list | 至少 1 tier | 每 tier: 积分阈值 + 奖励(Badge/Shop Item/Custom) |
| M-30 | Tiers 顺序 | — | — | 按积分阈值升序自动排序（不可手动排序）|
| M-31 | 奖励类型 | select per tier | 必填 | Badge（从 B31i 选）/ Shop Item（从 B31g 选）/ Custom text |

#### D07 — Shop Item Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-32 | Item Name | text input | 必填, 1-60 字符 | |
| M-33 | Description | textarea | 可选, max 300 字符 | |
| M-34 | Category | select | 必填 | NFT / Token / Merchandise / Experience / Other |
| M-35 | Image | file upload | 可选, PNG/JPG/SVG, max 2MB | 商品图片 |
| M-36 | Price (Points) | number | 必填, ≥1 | |
| M-37 | Stock | radio + number | 必填 | Unlimited / Limited (输入数量 ≥1) |
| M-38 | Availability Gate | select | 可选 | All / Level (min level N) / Badge (must hold badge X) |
| M-39 | Status | radio | 必填 | Save as Draft / Publish Now (→ D20) |

#### D08 — Lucky Wheel Config

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-40 | Wheel Name | text input | 必填 | |
| M-41 | Prizes | dynamic list | ≥2 prizes, 总概率 = 100% | 每 prize: 名称 + 类型 + 值 + 概率% |
| M-42 | Prize 类型 | select per prize | 必填 | Points / NFT / Token / Nothing (consolation) |
| M-43 | 概率验证 | computed | 实时显示 total% | 总和 ≠ 100% 时 Save 禁用 + 红色提示 |
| M-44 | Spin Cost | number | 必填, ≥0 | 每次抽奖消耗积分（0=免费）|
| M-45 | Spin Limit | radio | 必填 | Once per user / Daily / Unlimited |
| M-46 | Duration | date range | 可选 | 不填则永久有效 |

#### D09 — Badge Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-47 | Badge Name | text input | 必填, 唯一 | |
| M-48 | Description | textarea | 可选 | |
| M-49 | Icon | file upload / preset picker | 必填 | PNG/SVG max 1MB; 或 50+ 预设图标 |
| M-50 | Category | select | 必填 | Achievement / Engagement / Special |
| M-51 | Earn Condition | select | 必填 | Auto-trigger / Manual only |
| M-52 | Auto Condition (if auto) | select + params | 条件必填 | "Complete X tasks" / "Reach Level Y" / "X-day streak" / "Earn X points" |
| M-53 | Is Rare | toggle | 默认 off | Rare badge 在 C端有特殊视觉效果 |
| M-54 | 手动发放 | wallet input | 可选（Manual only 模式）| 发放给特定 wallet address |

#### D10 — Access Rule Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-55 | Rule Name | text input | 必填 | |
| M-56 | Rule Type | select | 必填 | Token Gate / NFT Hold / Level Requirement / Invite Only |
| M-57 | Token Gate Params | contract + min balance + chain | type=Token Gate 时必填 | 关联 B61 已配置链 |
| M-58 | Level Req Params | min level number | type=Level 时必填 | |
| M-59 | Invite Only Params | code or whitelist upload | type=Invite Only 时必填 | |
| M-60 | Denial Message | textarea | 可选, max 200 字符 | 不满足规则时 C端显示的提示 |

#### D11 — Homepage Section Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-61 | Section Type | select | 必填 | Banner / Quest Widget / Leaderboard Widget / Points Widget / Text / Custom HTML |
| M-62 | Title | text input | 必填 | |
| M-63 | 内容（Banner）| image upload + link | type=Banner 时必填 | |
| M-64 | 内容（Widget）| widget select | type=Widget 时必填 | 选择具体模块实例 |
| M-65 | 内容（Text）| rich text editor | type=Text 时必填 | |
| M-66 | 内容（Custom HTML）| code editor | type=Custom 时必填 | 基本 XSS 过滤 |
| M-67 | Visibility | select | 必填 | All users / Logged-in only / Level-gated (min level N) |

#### D16 — Announcement Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-68 | Title | text input | 必填, max 80 字符 | |
| M-69 | Content | textarea | 必填, max 500 字符 | 支持 plain text + 1 URL |
| M-70 | Type | select | 必填 | General / Event / Alert (影响 C端显示样式) |
| M-71 | Image | file upload | 可选 | PNG/JPG max 2MB |
| M-72 | CTA Button | text + URL | 可选 | "Learn More" + 跳转链接 |
| M-73 | Schedule | radio | 必填 | Publish Now / Schedule (日期时间选择) |
| M-74 | Pin to Top | toggle | 默认 off | Pinned 公告始终置顶 |

#### D17 — Featured Slot Editor

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| M-75 | Content Type | select | 必填 | Quest / LB Sprint / Milestone / External URL |
| M-76 | Content Select | select | type≠External 时必填 | 从已发布的对应内容中选择 |
| M-77 | External URL | URL input | type=External 时必填 | |
| M-78 | Custom Title | text input | 可选 | 覆盖原标题 |
| M-79 | Custom Image | file upload | 可选 | 覆盖原图 |

#### D18 — Segment Detail Panel

| # | 字段 | 类型 | 说明 |
|---|------|------|------|
| M-80 | Segment 类型 | display | Power / Active / At Risk / Dormant |
| M-81 | User 列表 | table | wallet address / last active / total points / modules used |
| M-82 | Search | text input | 按 wallet address 搜索 |
| M-83 | Filter | — | 可按模块使用情况筛选 |
| M-84 | Export CSV | button | 导出当前分群用户数据 |

#### D19 — Promo Kit Generator

| # | 字段 | 类型 | 说明 |
|---|------|------|------|
| M-85 | Target Platform | radio | Twitter / Discord / Telegram |
| M-86 | AI 生成文案 | display + editable textarea | 可二次编辑；字数限制随平台变化（Twitter 280）|
| M-87 | Generated Banner | image display | AI 生成品牌化图片；"Regenerate" 按钮 |
| M-88 | Copy Text | button | clipboard API |
| M-89 | Download Banner | button | PNG 下载 |
| M-90 | Share on [Platform] | button | 打开对应平台分享链接 |

#### D20 — Publish Readiness Check（补充缺失）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| M-91 | OPS | 检查超时（10s）：显示 "Verification timed out. Publish anyway?" 选项 | P0 |
| M-92 | OPS | "Publish anyway"：跳过 readiness check 强制发布（仅订阅通过时允许）| P1 |
| M-93 | OPS | Readiness check 结果缓存：5分钟内再次触发同一 Publish 按钮时，跳过重新检查 | P1 |

---

## 三、req_white_label.md 缺失清单

### Hub 页面

#### B14 — WL Hub Empty（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-01 | OPS | 路径卡片选择：click → 紫色 border 2px；同时只能选一个 | P0 |
| W-02 | OPS | 路径选择后 CTA 文案变化："Start with Embed / Domain / SDK" | P0 |
| W-03 | OPS | 路径卡片展开：选中后展开详情（适用场景 + 技术要求 + 预计时间）| P1 |
| W-04 | OPS | "Recommended" badge：Embed 路径有 ★ 标识；其他路径无 | P1 |
| W-05 | EDGE | 用户已有 WL 配置时的路由守卫（此页不应可访问）| P0 |

#### B15 — WL Hub Active / Guided（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-06 | OPS | Checklist 每步展开内容：  ①Configure widgets→"Open Widget Library →" ②Build pages→"Open Page Builder →" ③Preview→内嵌 preview ④Send Dev Kit→生成 Dev Kit URL + Copy 按钮 ⑤Integration verified→WebSocket 状态指示器 ⑥Announce→D19 Promo Kit | P0 |
| W-07 | OPS | "Send Dev Kit"：POST /api/devkit/{project_id}/generate → 返回 URL → 显示 taskon.xyz/devkit/abc123 + Copy 按钮 + 也可 "Open in New Tab" | P0 |
| W-08 | OPS | "Integration verified" 步骤：WebSocket /ws/wl/integration-ping 监听第一个来自项目域名的 API ping；自动标记完成 | P0 |
| W-09 | OPS | Toolkit 卡片区（Widget Library / Page Builder / Domain / Brand）："Configure →" 各自跳转对应子页面 | P0 |
| W-10 | OPS | "First interaction" 步骤：WebSocket /ws/wl/first-interaction 监听；自动标记 | P0 |

#### B16 — WL Hub Management（整体缺失）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-11 | OPS | 完整页面结构（设计稿 Node `UPAfV`）未在文档中描述 | P0 |
| W-12 | OPS | Deployment Stats 区块：已部署版本 / 活跃用户 / 总交互数 / API 调用数（24h）| P0 |
| W-13 | OPS | Feature Management 卡片：Widget Library / Page Builder / Smart Rewards / Brand / SDK 各自状态 + 快速入口 | P0 |
| W-14 | OPS | "View All Deployments" → 查看历史发布记录 | P1 |

---

### 创建向导

#### B37 — Step 1: Path Selection（缺失详情）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-15 | OPS | Embed 路径选中后：下方注释 "Next: Choose embed method (Iframe / Widget / Page Builder)" | P0 |
| W-16 | OPS | Domain 路径选中后：注释 "Next: Configure your custom domain and CNAME" | P0 |
| W-17 | OPS | SDK 路径选中后：注释 "Next: Generate your API keys and integration code" | P0 |
| W-18 | OPS | "Next" 按钮根据路径跳转不同 Step 2 变体（B17/B57/B60）| P0 |
| W-19 | EDGE | 无选择时 "Next" disabled | P0 |

#### Step 2 变体（B17/B57/B58/B59/B60）—— 各自缺失字段规格

**B17 — Step 2: Widget Config**

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-20 | Embed Method | radio | 必填 | Widget Library / Iframe / Page Builder（选择影响后续配置）|
| W-21 | Widget Selection | checkbox list | ≥1 | 选择要嵌入的 Community Modules（Leaderboard / Task List / User Center 等）|
| W-22 | SSO Method | radio | 必填 | Wallet Auth / OAuth2+JWT / None (preview only) |
| W-23 | Target Domain | URL input | 必填 | 将嵌入的页面所在域名（用于 CORS 配置）|

**B57 — Step 2: Domain Config**

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-24 | Custom Domain | text input | 必填, valid domain format | "community.yourproject.io" |
| W-25 | CNAME Target | display (read-only) | — | TaskOn 提供的目标 CNAME 记录值 |
| W-26 | DNS Provider | select | 可选 | Cloudflare / Route53 / Namecheap 等（显示对应配置教程链接）|
| W-27 | Verify DNS | button | — | 手动触发 DNS 验证；显示 polling 状态 |

**B58 — Step 2: Iframe Config**

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-28 | Iframe URL | display (auto-generated) | — | share.taskon.io/embed/{project_id} |
| W-29 | SSO Option | radio | 可选 | With SSO / Without SSO (read-only mode) |
| W-30 | Allowed Origins | text input | 可选 | 允许嵌入的域名白名单（逗号分隔）|
| W-31 | Iframe Code | code block | — | 自动生成 <iframe> 代码；Copy 按钮 |

**B59 — Step 2: PB Config (has-pages variant)**

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-32 | 已有 Pages 列表 | radio list | ≥1 选 | 选择基于哪个已有 Page 继续（Node `zW40A`）|
| W-33 | 无已有 Pages 时 | template cards | ≥1 选 | 选择 Page Builder 模板（Node `XHwzp`）|

**B60 — Step 2: SDK Config**

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-34 | SDK Mode | radio | 必填 | Full Custom (Headless) / Hybrid (SDK + TaskOn UI fallback) |
| W-35 | API Key | display (auto-generated) | — | 显示 pk_live_xxx；Copy 按钮 |
| W-36 | Project ID | display | — | 唯一项目标识符 |

#### B38 — Step 3: Brand（缺失字段规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-37 | MODAL | 完整字段列表（与 B40 Brand Settings 相同字段子集）| P0 |
| W-38 | OPS | 右侧实时 Preview 更新（颜色/Logo 变化即时反映）| P0 |
| W-39 | OPS | "Skip for now" 链接（允许跳过，之后从 B40 再配置）| P1 |

#### B56 — Step 4: Preview & Publish（缺失规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-40 | OPS | WL Preview 内容：显示已配置的 Widget / Page / Domain 部署效果预览 | P0 |
| W-41 | OPS | Readiness Checklist 项（WL 版）：Path configured / Brand set / Widget/Page configured / Domain verified（条件性）| P0 |
| W-42 | OPS | "Publish White Label" → D20 → 通过 → POST /api/white-label/publish → 跳转 B15 | P0 |
| W-43 | EDGE | Publish 后 Dev Kit 链接自动生成 | P0 |

---

### 子页面

#### B18 — Domain Setup（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-44 | OPS | DNS 验证自动轮询：每 30s 检查一次；显示 "Checking..." → "✓ Verified" / "✗ Not found yet" | P0 |
| W-45 | OPS | CNAME 记录展示格式：Type=CNAME / Host=@ (or community) / Value=xxx.taskon.io / TTL=300 | P0 |
| W-46 | OPS | SSL 状态：DNS 验证通过后，TaskOn 自动配置 Let's Encrypt；显示 "SSL: Pending" → "SSL: Active" | P1 |
| W-47 | OPS | "Edit Domain Settings" 按钮 → D20（此操作为何需要发布检查？）| P1 |
| W-48 | EDGE | 域名已被其他 TaskOn 项目使用：报错 "This domain is already claimed" | P0 |

#### B19 — Embed Options（缺失完整描述）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-49 | OPS | B19v Neutral 状态（Node `Rwq2K`）：显示 3 种 Embed 方式卡片（Iframe/Widget/Page Builder）供选择 | P0 |
| W-50 | OPS | 选择 Iframe → 跳转 B42 Iframe Embed 页面 | P0 |
| W-51 | OPS | 选择 Widget Library → 跳转 B20/B22（根据是否有已配置 widget）| P0 |
| W-52 | OPS | 选择 Page Builder → 跳转 B23/B25（根据是否有已创建 page）| P0 |

#### B42 — Iframe Embed（整体缺失）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-53 | OPS | 完整页面结构（设计稿 Node `ByGS0`）| P0 |
| W-54 | OPS | Iframe URL 展示：完整 embed URL + Copy 按钮 | P0 |
| W-55 | OPS | Iframe Code 生成：`<iframe src="..." width="100%" height="600" frameborder="0"></iframe>` + Copy | P0 |
| W-56 | OPS | SSO 集成说明：步骤指导（如何传递 JWT token 给 iframe）| P0 |
| W-57 | OPS | "Test Embed" 按钮：在当前页内嵌显示 iframe 预览 | P1 |

#### B20/B21/B22 — Widget Library（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-58 | OPS | B21 Widget Config 完整字段：Module selector（哪个 Community Module）/ Style Config（primary color inherit / border radius / padding）/ Preview 面板（实时渲染）| P0 |
| W-59 | OPS | B21 Embed Code：基于 module + style 自动生成 `<script>` 代码；Copy 按钮 | P0 |
| W-60 | OPS | B22 Active：每个 widget 行的操作：Edit（→ B21）/ Copy Embed / Deploy（→ D20）/ Delete（确认）| P0 |
| W-61 | OPS | B22 Widget 状态："Configured" (绿色, 可部署) / "Deployed" (蓝色, 已在用) / "Draft" (灰色) | P0 |
| W-62 | OPS | B20 Template cards：预置 Widget 模板（Leaderboard / Task List / User Center / Check-in）→ 点击 → B21 预填 | P0 |

#### B23/B24/B25 — Page Builder（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-63 | OPS | B24 Canvas 拖拽：Widget Block 拖拽手柄（左边 ⋮⋮）；拖拽时 block 高亮 + 其他 block 占位符出现 | P0 |
| W-64 | OPS | B24 "+ Add Widget Block"：弹出 dropdown 显示 "Available Widgets"（来自 B22 已配置的 widget 列表）| P0 |
| W-65 | OPS | B24 Widget Block "× Remove"：确认后从 canvas 移除；Settings 面板 "Widgets on Page" 列表同步更新 | P0 |
| W-66 | OPS | B24 URL Slug：基于 Page Name 自动生成（小写+连字符）；可手动编辑；输入时验证唯一性 | P0 |
| W-67 | OPS | B24 Theme Toggle (Light/Dark)：切换后 Canvas preview 即时变化 | P0 |
| W-68 | OPS | B24 Settings 面板 Widget 拖拽排序（"Widgets on Page" 列表）→ Canvas 同步变化 | P0 |
| W-69 | OPS | B25 "Edit Page"：加载已有 Page 数据 → 跳转 B24（edit mode）| P0 |
| W-70 | OPS | B25 Page stats：Page Views / Unique Visitors / Widget Clicks / Total Completions 数据来源 | P1 |

#### B40 — Brand Settings（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-71 | OPS | Logo 上传："Change" → file picker（SVG/PNG, min 256×256, max 2MB）→ 验证尺寸/格式 → 上传 → 立即预览 | P0 |
| W-72 | VALID | Logo 上传失败情形：文件格式不支持、尺寸过小、体积过大 → 各自的错误提示文案 | P0 |
| W-73 | OPS | Color Picker：点击色块 → 展开 color picker（hex input + hue slider + 透明度滑块）| P0 |
| W-74 | OPS | Font Dropdown：列出支持的 Google Fonts 字体（约 20 种）+ Preview 文本实时更新 | P1 |
| W-75 | OPS | Custom CSS Editor：基础语法高亮（单色关键字高亮）；"Preview" → 在 preview 面板实时应用 CSS | P1 |
| W-76 | CASCADE | "Save Changes" 成功：toast "Brand updated" → 已部署的 Widget/Page 下次加载时反映新品牌 | P0 |
| W-77 | EDGE | 已部署时修改品牌：是否需要重新发布（re-deploy）还是立即生效 | P0 |

#### B41 — SDK & API（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-78 | OPS | "Regenerate Key"：确认弹窗 "This will immediately invalidate your current key. All integrations using this key will stop working." → 确认 → POST → 新 key 展示 → 旧 key 立即失效 | P0 |
| W-79 | OPS | Webhook "+ Add Endpoint"：inline 展开表单（URL input + Events 多选 checkbox）→ Save → POST /api/white-label/sdk/webhooks | P0 |
| W-80 | OPS | Webhook "Test"：POST 测试事件到该 endpoint → 显示响应状态码 + 响应时间 | P1 |
| W-81 | OPS | Webhook 可选事件列表：task.completed / points.earned / user.joined / level.up / badge.earned / sprint.ended | P0 |
| W-82 | OPS | API Usage 进度条 hover：tooltip "12,450 of 50,000 requests used (24.9%)" | P1 |
| W-83 | EDGE | API key 即将超出用量（>80%）：amber 警告条 "You've used 85% of your monthly quota. Upgrade plan →" | P1 |

#### B26 — WL Integration Center（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-84 | OPS | Twitter OAuth 流程：popup window → 用户授权 → 回调 → Connected 状态 | P0 |
| W-85 | OPS | Google Analytics "Connect"：输入 GA4 Measurement ID → 验证 → Connected | P0 |
| W-86 | OPS | SSO / OAuth (Developer Tools) "Connect"：进入 SSO 配置页（选择 Wallet Auth 或 OAuth2，配置 JWT secret）| P0 |
| W-87 | OPS | SDK Configuration "Connect"：跳转 B41 SDK & API 页面 | P0 |
| W-88 | OPS | 集成 error 状态：红色 border + error 图标 + "Token expired: Reconnect →" | P0 |

#### B51 — Contract Registry（整体缺失）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-89 | OPS | 完整页面结构（设计稿 Node `OKEqS`）：Header / Stats / Contract Table / 按钮路由 | P0 |
| W-90 | OPS | Stats Row：Total Contracts / Verified / Events Captured (24h) / Active Rules (关联) | P0 |
| W-91 | OPS | Contract Table 列：Contract Name / Network / Address / Status(Verified/Pending/Error) / Events / Actions | P0 |
| W-92 | OPS | "+ Register Contract" → D12 Contract Register Form | P0 |
| W-93 | OPS | Contract verify：自动链上验证（检查合约是否真实存在 + ABI 可读）→ status 变 Verified | P0 |
| W-94 | OPS | Contract row 操作：View Events / Edit / Delete（已被 Rule 引用时禁止删除）| P0 |

#### D12 — Contract Register Form（整体缺失）

| # | 字段 | 类型 | 验证 | 说明 |
|---|------|------|------|------|
| W-95 | Contract Name | text input | 必填 | 内部标识名 |
| W-96 | Network | select | 必填 | Ethereum / BSC / Polygon / Base / Arbitrum / Optimism / Avalanche |
| W-97 | Contract Address | text input | 必填, 0x + 40 hex chars | |
| W-98 | ABI | textarea / JSON upload | 必填 | 粘贴 ABI JSON 或上传文件；自动解析 events |
| W-99 | Events to Monitor | checkbox list (from ABI) | ≥1 | 解析 ABI 后列出可监听的 event 列表供选择 |
| W-100 | Verify on Save | toggle | 默认 on | 保存后立即尝试链上验证 |

#### B43 — Page Analytics（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-101 | OPS | Date Range Picker：预设 7d/30d/90d + 自定义范围 | P1 |
| W-102 | OPS | "Page Views Over Time" chart：D/W/M 切换（Day/Week/Month 粒度）；hover tooltip | P1 |
| W-103 | OPS | "Top Pages" 表：点击页面行 → 筛选下方 Widget Interactions 表为该页面数据 | P1 |
| W-104 | OPS | "Export" 按钮：导出当前视图数据为 CSV | P1 |
| W-105 | OPS | Conversion Funnel：各步骤数字 hover tooltip；点击步骤跳转 Page Builder 对应 widget | P1 |

#### B48 — Dev Kit Page（缺失操作规格）

| # | 类别 | 缺失内容 | 优先级 |
|---|------|---------|--------|
| W-106 | OPS | Dev Kit URL 生成：taskon.xyz/devkit/{project_id}；project_id 在 WL 发布时生成；URL 公开可访问（无需登录）| P0 |
| W-107 | OPS | Package manager tabs（npm/yarn/CDN）：切换显示对应安装命令；code block 一键 Copy | P0 |
| W-108 | OPS | SSO Provider selector：选 "Wallet Auth" 显示钱包 SDK 集成代码；选 "OAuth2/JWT" 显示 JWT 配置代码 | P0 |
| W-109 | OPS | Widget 展开：每个 widget 展开显示嵌入代码（project_id 预填）+ mini preview | P0 |
| W-110 | OPS | "Verify Integration" 按钮：POST /api/devkit/{id}/verify → 显示 loading → 检测是否收到该项目域名的 API ping → 成功/失败状态 | P0 |
| W-111 | OPS | 验证成功：按钮变为 "✓ Verified" (green)；显示 "Integration verified at [time] from [domain]" | P0 |
| W-112 | EDGE | Dev Kit 页面过期（project 删除或 dev kit URL 撤销）：显示 "This link is no longer valid" 友好错误页 | P1 |

---

## 四、汇总统计

| 文档 | 总缺失项 | P0 | P1 | P2 |
|------|---------|----|----|-----|
| req_community.md | 153 项 | 89 | 55 | 9 |
| req_white_label.md | 112 项 | 68 | 38 | 6 |
| **合计** | **265 项** | **157** | **93** | **15** |

---

## 五、建议处理顺序

### 第一批（P0 阻塞开发，优先处理）
1. **所有 Modal 字段规格**（D01-D12, D14-D17）— 开发 Modal 前必须明确字段
2. **所有 Wizard 字段验证规则**（B13/B34/B35/B55 / B37/B17/B57/B58/B59/B60/B38/B56）
3. **B51 Contract Registry 整体**（当前完全缺失）
4. **B42 Iframe Embed 整体**（当前完全缺失）
5. **B16 WL Hub Management 整体**（当前完全缺失）
6. **通用表格页面操作规格**（C-66 ~ C-79 适用于所有模块管理页）
7. **B31 Sectors & Tasks 独特操作**（拖拽/跨 Sector/inline edit/bulk）

### 第二批（P1 影响体验）
1. **Hub 页面所有操作细节**（B09-B12 / B14-B16 的 per-block ops）
2. **运营辅助页操作细节**（B32/B33/B54/B61 / B40/B41/B26/B43/B48）
3. **各模块管理页独特操作**（B31a-B31i 特有业务规则）

### 第三批（P2 完善度）
1. **动画规格**（toggle 动画、modal 开关动画、drag 视觉反馈）
2. **Tooltip 文案**（所有 hover tooltip 完整文案）
3. **外部链接 URL**（帮助文档、视频教程具体地址）
