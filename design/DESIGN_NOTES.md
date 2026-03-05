# 设计决策日志

用途：记录每个重要的设计决策和AI反馈

## 日期格式
[YYYY-MM-DD] - [屏幕名称]

---

## 示例条目

### [2024-01-25] - Landing Page Hero Background

【决策】
使用蓝色渐变背景 (#5D7EF1 → #9B7EE0) for Hero区块

【理由】
- 与TaskOn现有品牌一致
- 蓝色传达信任和技术感
- 渐变增加视觉深度

【考虑过的替代方案】
- 纯蓝色：太平面
- 红色：过于激进，不符合SaaS风格
- 白色+图片：可维护性差

【AI反馈 (Claude)】
"好的选择。建议检查：浅色文字在蓝色背景上是否有足够的对比度（应该≥4.5:1）。"

【最终状态】
✅ 已调整，对比度满足WCAG AA标准

---

## 你的日志从这里开始

### [2026-01-29] - Landing Page 重大重构：业务模式清晰化

【决策】
将Landing Page从"功能展示"重构为"业务模式选择"，清晰展示TaskOn的两个核心商业模式：
1. 白标SaaS（Build on Your Domain）
2. 平台模式（Grow on TaskOn Platform）

【理由】
- **原问题**：之前的设计混淆了功能和商业模式，用户不清楚TaskOn提供什么服务
- **业务洞察**：TaskOn实际上提供两种截然不同的价值：
  - 白标：嵌入式工具，用户在项目方自己网站完成任务（数据归项目方）
  - 平台：托管服务，用户在TaskOn.xyz完成任务（快速获客）
- **用户困惑点**：项目方需要先理解"我的用户在哪里完成任务"才能选择合适的工具

【设计变更】
1. **Hero Section**：
   - 旧："One-Stop Platform for Marketing Growth"
   - 新："Your Complete Web3 Growth Stack"
   - 副标题强调"在你的域名构建 或 在我们平台启动"

2. **Value Props Section**：
   - 旧：4个功能卡片（增长、参与、PMF、品牌）
   - 新：2x2结构
     - 蓝色卡片组：白标路径（Build on Your Domain + Complete Control）
     - 绿色卡片组：平台路径（Grow on TaskOn + Quest/Community）

3. **Role Selector**：
   - 旧：CMO/Community/CTO/CEO（按职位）
   - 新：Growth/Community/Dev/Executive（按需求和推荐路径）
   - 每个角色明确指向一个路径

【考虑过的替代方案】
- 方案A：完全分离两个产品页面（Platform vs White-label）
  - 弃用原因：增加导航复杂度，用户可能不知道从哪里进入
- 方案B：保持原有功能展示，在深层才区分模式
  - 弃用原因：用户前期困惑太多，跳出率高
- 方案C：✅ 当前方案 - 首页即展示两条路径，让用户快速选择

【文档更新】
- ✅ design/EXPORT_FOR_CLAUDE.txt - 添加业务模式说明
- ✅ design/AI_CONTEXT.md - 更新角色和路径映射
- ✅ design/DESIGN_NOTES.md - 记录此次决策（本条）

【下一步】
- [ ] 添加图标到卡片和角色选择器
- [ ] 创建"Platform vs White-label对比表"页面
- [ ] 为每个团队角色设计详细的流程页面
- [ ] 添加真实案例和数据支撑

【最终状态】
✅ 已完成 - 后续优化为决策树结构（见下一条）

---

### [2026-01-29] - 决策树重构：3选1 + 二级选择

【决策】
将Landing Page重构为清晰的决策树结构：
- 第一层：3个主要选择（Build on My Domain / Grow on TaskOn / Use Both）
- 第二层：针对"Grow on TaskOn"的细分（Quest Only / Quest + Community）

【理由】
- **用户反馈**：需要更清晰的"我的用户在哪里完成任务"决策路径
- **Playbook对齐**：完全对应playbook第18节的产品层级（Quest → Community → White Label）
- **决策优先级**：用户首先需要决定"用户在哪"，然后才是"功能多全"

【设计变更】
1. **Choose Your Growth Model Section**：
   - 3个大卡片，清晰展示3个选择
   - 每个卡片包含：Badge + Icon + Title + Description + Features + CTA
   - 颜色区分：蓝色（白标）、绿色（平台）、紫色（混合）

2. **Platform Detail Section**（新增）：
   - 标题："If You Choose TaskOn Platform..."
   - 2个对比卡片：Quest Only vs Quest + Community
   - Quest卡片标注"STARTER"，Community卡片标注"COMPLETE"
   - 清晰列出各自功能差异

3. **保留Team Paths**：
   - 4个团队角色卡片保留，作为辅助导航
   - 每个角色指向推荐路径

【产品层级映射】（基于playbook）
- Quest（STARTER）：获客+激活，托管在TaskOn，22M用户库，基础验证
- Community（COMPLETE）：Quest + Points/Levels + Leaderboards + Milestones + Benefits Shop
- White Label（PREMIUM）：完全嵌入，你的域名，完整归因+ROI分析

【用户决策流程】
1. 我的用户在哪里完成任务？
   - 我的网站 → 白标
   - TaskOn平台 → 进入第2步
   - 两个都要 → 混合模式

2. （如果选择平台）我需要多完整的功能？
   - 只要快速获客 → Quest Only
   - 需要完整留存系统 → Quest + Community

【文档更新】
- ✅ design/EXPORT_FOR_CLAUDE.txt - 更新为决策树结构
- ✅ design/AI_CONTEXT.md - 更新为两层选择
- ✅ design/DESIGN_NOTES.md - 记录此次决策（本条）

【下一步】
- [ ] 为3个主要选择各创建详细页面
- [ ] 创建对比表（Quest vs Community vs White Label）
- [ ] 添加真实案例和定价信息
- [ ] 完善Team Paths与模型选择的关联

【最终状态】
✅ 决策树结构完成 - 清晰的3选1 + 二级选择

---

### [2026-01-29] - 导航栏结构设计

【决策】
设计清晰的导航结构，区分"产品"（自助工具）和"解决方案"（服务交付）：

**Product**（产品下拉）：
- Quest - 快速获客活动
- Community - 完整留存系统
- White Label - 嵌入式解决方案

**Solutions**（解决方案下拉）：
- Growth Delivery (CPA/CPS/BPS) - 增长交付服务
- Custom Development - 定制化开发
- Managed Services - 代运营
- Co-Marketing - 联合营销

【理由】
- **清晰区分**：Product = 用户自己用的工具，Solutions = TaskOn提供的服务
- **完整覆盖**：对应taskon_services_en.md的4类服务
- **符合认知**：B2B SaaS常见的Product vs Solutions区分
- **弱化服务**：首页主推Product，Solutions作为二级入口

【映射关系】
| 导航项 | 对应服务类别 | 目标用户 |
|---|---|---|
| Product → Quest | Subscription (Quest) | 需要快速获客的团队 |
| Product → Community | Subscription (Community) | 需要长期留存的团队 |
| Product → White Label | Subscription (White Label) | 需要完全控制的团队 |
| Solutions → Growth Delivery | CPA/CPS/BPS | 需要专业增长服务的团队 |
| Solutions → Custom Dev | Custom/Managed (Dev) | 需要定制开发的团队 |
| Solutions → Managed Services | Custom/Managed (Ops) | 需要代运营的团队 |
| Solutions → Co-Marketing | Joint Programs | 需要联合营销的团队 |

【下一步】
- [ ] 设计Product和Solutions的下拉菜单详情
- [ ] 为每个Solutions项创建详情页
- [ ] 在首页适当位置添加Solutions的入口提示

【最终状态】
✅ 导航结构设计完成

---