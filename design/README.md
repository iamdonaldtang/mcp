# TaskOn Onboarding v2 设计项目

## 🎯 项目目标
重新设计TaskOn的项目侧Onboarding流程和核心产品介绍。
采用Role-based设计，让CMO/Community/CTO/CEO快速找到对他们的价值。

## 🛠️ 工具与工作流
- **设计**：Pencil (VSCode插件)
- **AI助手**：Claude (VSCode或浏览器)
- **文档**：Markdown
- **版本控制**：Git

## 📁 文件结构
design/
├─ AI_CONTEXT.md              # 项目背景（Claude会读）
├─ AI_PROMPTS.md              # 10个设计反馈提示词
├─ DESIGN_NOTES.md            # 设计决策日志
├─ EXPORT_FOR_CLAUDE.txt      # 一次性导入Claude的完整上下文
├─ onboarding/
│  ├─ 01_landing_page.pen     # Landing Page设计（Pencil文件）
│  ├─ 02_role_selector.pen    # Role Selector设计
│  ├─ 03_cmo_path.pen         # CMO路径设计
│  ├─ ...
│  └─ ...
├─ components/
│  └─ components_library.pen  # UI组件库
├─ exports/
│  └─ (PNG导出文件)
└─ README.md                   # 本文件

## 🚀 快速开始

### 1. 环境准备（5分钟）
```bash
# 确保安装了必要的工具
# VSCode + Pencil插件 + Claude for VSCode
```

### 2. 阅读项目背景（10分钟）
```bash
# 打开并阅读这些文件
cat design/AI_CONTEXT.md
cat design/EXPORT_FOR_CLAUDE.txt
```

### 3. 启动Claude助手（2分钟）
```bash
# 在Claude Chat中粘贴 EXPORT_FOR_CLAUDE.txt 的内容
# 这样Claude就有完整的项目上下文了
```

### 4. 开始设计（立即）
```bash
# 打开Pencil，开始设计Landing Page
# 每30分钟截图给Claude反馈一次
```

## 📊 设计进度

- [ ] Week 1
  - [ ] Landing Page完成
  - [ ] Role Selector完成
  - [ ] CMO Path完成
  
- [ ] Week 2
  - [ ] Community Manager Path完成
  - [ ] CTO Path完成
  - [ ] CEO Path完成
  - [ ] Handoff文档完成

- [ ] Week 3-4
  - [ ] 前端实现
  - [ ] 测试
  - [ ] 上线

## 🤖 与Claude合作

### 快速工作流
1. 在Pencil中设计
2. 截图 (Cmd+Shift+4 或 Win+Shift+S)
3. 复制一个提示词 (来自 AI_PROMPTS.md)
4. 粘贴提示词 + 截图到Claude Chat
5. 获得具体反馈
6. 调整设计
7. 重复 1-6

### 提示词库位置
快速审查、角色适配、文案优化等提示词在：
→ design/AI_PROMPTS.md

### 记录决策
每个重要的设计决策记录在：
→ design/DESIGN_NOTES.md

## 📱 核心概念速记

**4个用户角色**：CMO / Community / CTO / CEO
**各自的快赢**：
- CMO: 30分钟启动Quest
- Community: 一键配置Points系统
- CTO: 无代码嵌入组件
- CEO: ROI计算器

**设计原则**：
✨ 简洁 > 复杂
✨ 故事 > 功能列表
✨ 快赢 > 全功能展示
✨ 证明 > 承诺

## 🎨 设计系统

**颜色**：蓝(#5D7EF1) 绿(#48BB78) 紫(#9B7EE0) 橙(#ED8936)
**字体**：Inter
**标题大小**：Hero 48px, Section 32px
**Body**：16px Regular
**间距单位**：8px

## 📞 常见问题

**Q: 如何开始与Claude合作？**  
A: 看 [Pencil_Claude_Integration_Guide.md](../Pencil_Claude_Integration_Guide.md)

**Q: 我应该多久找Claude反馈一次？**  
A: 每完成一个小功能（约30分钟）就截图问一下。

**Q: 设计完了怎么交给前端？**  
A: 写Handoff.md并使用AI_PROMPTS.md中的"前端开发准备"提示词。

## 🔗 相关文档

- [TaskOn增长引擎手册](../../taskon_growth_engine_playbook_en.md)
- [Pencil快速上手指南](../Pencil_Quick_Start_Guide.md)
- [执行检查清单](../Design_Execution_Checklist.md)
- [Pencil+Claude集成指南](../Pencil_Claude_Integration_Guide.md)

---

**准备好了？打开Pencil，开始设计吧！🚀**

EOF

echo "✅ README.md 创建完成"

# 7. 最后信息
echo ""
echo "🎉 所有文件初始化完成！"
echo ""
echo "接下来的步骤："
echo "1️⃣  打开 VSCode"
echo "2️⃣  阅读 design/AI_CONTEXT.md"
echo "3️⃣  打开 Claude Chat，粘贴 design/EXPORT_FOR_CLAUDE.txt 的内容"
echo "4️⃣  在Pencil中创建Landing Page设计"
echo "5️⃣  截图给Claude反馈"
echo ""
echo "开始设计吧！💪"