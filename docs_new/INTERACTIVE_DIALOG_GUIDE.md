# 交互式对话器调用场景指南

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 架构设计说明

---

## 📋 目录

1. [什么是交互式对话器](#什么是交互式对话器)
2. [在系统中的定位](#在系统中的定位)
3. [典型调用场景](#典型调用场景)
4. [与普通工作流的区别](#与普通工作流的区别)
5. [调用时机决策树](#调用时机决策树)
6. [实际案例分析](#实际案例分析)
7. [技术实现方案](#技术实现方案)

---

## 什么是交互式对话器

### 定义

**交互式对话器 (Interactive Dialog)** 是一个AI能力层组件，负责管理**多轮对话**和**动态需求调整**。

```
普通工作流 vs 交互式对话

┌─────────────────────────────────┐
│      普通工作流（一次性生成）     │
├─────────────────────────────────┤
│  用户输入 → 工作流 → 内容输出    │
│      ↓                           │
│    结束（无法修改）               │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│    交互式对话（多轮迭代）         │
├─────────────────────────────────┤
│  用户输入1 → 初步内容             │
│      ↓                           │
│  用户反馈1 → 调整内容             │
│      ↓                           │
│  用户反馈2 → 再次调整             │
│      ↓                           │
│  用户确认 → 最终内容              │
└─────────────────────────────────┘
```

### 核心特点

| 特点 | 普通工作流 | 交互式对话器 |
|------|-----------|-------------|
| **交互次数** | 一次性 | 多轮迭代 |
| **需求明确性** | 要求清晰完整 | 可以逐步澄清 |
| **调整能力** | 需重新生成 | 增量调整 |
| **上下文管理** | 不需要 | 必须记住历史 |
| **适用场景** | 需求明确 | 需求模糊/探索性 |
| **用户负担** | 需要一次性说清楚 | 可以边说边调整 |

---

## 在系统中的定位

### 架构层级

```
┌────────────────────────────────────────────────┐
│              应用层 (API)                       │
│          /api/workflow - 普通工作流              │
│          /api/chat - 交互式对话 ⭐               │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│          领域适配器层 (Domain Adapter)          │
│        解析用户意图，提供领域知识               │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│             AI能力层 (AI Capabilities)          │
│  ┌──────────────────┐  ┌──────────────────┐   │
│  │ 意图识别         │  │ 学习目标生成      │   │
│  └──────────────────┘  └──────────────────┘   │
│  ┌──────────────────┐  ┌──────────────────┐   │
│  │ 内容发现         │  │ 交互式对话器 ⭐   │   │
│  └──────────────────┘  └──────────────────┘   │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│          工作流引擎层 (Workflow Engine)         │
│         普通模式 vs 对话模式                    │
└────────────────────────────────────────────────┘
```

### 与其他AI能力的关系

| AI能力 | 与交互式对话器的关系 |
|--------|-------------------|
| **意图识别** | ✅ **被调用** - 每轮对话都需要识别意图 |
| **学习目标生成** | ✅ **被调用** - 根据对话生成和调整目标 |
| **内容发现** | ✅ **被调用** - 根据对话动态搜索内容 |
| **故事化叙述** | ✅ **被调用** - 根据反馈调整叙述风格 |
| **质量评估** | ✅ **被调用** - 评估每轮生成的质量 |
| **多模态生成** | ❌ **不直接相关** - 对话结束后再调用 |

**关键点**: 交互式对话器是一个**协调者**，它调用其他AI能力来完成多轮对话。

---

## 典型调用场景

### 场景分类

根据实际需求，交互式对话器的调用可以分为以下场景：

```
调用场景分类
├── 1. 需求探索场景（用户不清楚自己要什么）
├── 2. 需求澄清场景（用户表达模糊）
├── 3. 内容迭代场景（对初步内容不满意）
├── 4. 个性化定制场景（需要深度定制）
└── 5. 教学辅导场景（一对一辅导模式）
```

---

### 场景1️⃣: 需求探索场景

**触发条件**: 用户不确定自己要什么，需要系统引导

**典型对话流程**:

```
【用户】: "我想学点化学知识"
    ↓
【系统】: 检测到需求模糊 → ⭐ 启动交互式对话器
    ↓
【对话器】: "请问您是：
           1. 中学生准备考试？
           2. 成人兴趣学习？
           3. 职业技能提升？"
    ↓
【用户】: "成人兴趣学习"
    ↓
【对话器】: "化学领域很广，您更感兴趣：
           1. 日常生活中的化学（食品、清洁剂）
           2. 有机化学（药物、材料）
           3. 无机化学（元素、化合物）"
    ↓
【用户】: "日常生活中的化学"
    ↓
【对话器】: "好的！具体想了解：
           1. 食品添加剂的化学原理
           2. 洗涤剂为什么能去污
           3. 电池的工作原理"
    ↓
【用户】: "食品添加剂"
    ↓
【对话器】: "需求已明确！开始生成内容..."
    ↓
调用普通工作流生成最终内容
```

**关键点**:
- ✅ **逐层收窄范围**: 从大类 → 细分 → 具体主题
- ✅ **提供选项**: 降低用户认知负担
- ✅ **记住上下文**: "成人兴趣学习" + "日常生活化学"

---

### 场景2️⃣: 需求澄清场景

**触发条件**: 用户输入包含歧义或关键信息缺失

**示例1: 歧义澄清**

```python
# ========== 案例：术语歧义 ==========

【用户】: "讲解《日出》这幅画"
    ↓
【系统】: 检测到歧义（多个同名作品）→ ⭐ 启动交互式对话器
    ↓
【对话器】: "找到多个《日出》作品：
           1. 莫奈《日出·印象》（1872年，印象派）
           2. 透纳《日出与海怪》（1845年，浪漫主义）
           3. 弗里德里希《雾海上的日出》（1818年，德国浪漫派）
           
           请问您想了解哪一幅？"
    ↓
【用户】: "莫奈的"
    ↓
【对话器】: "明白了，《日出·印象》。您希望：
           1. 快速浏览（5分钟）
           2. 深度讲解（15分钟）
           3. 艺术史背景（20分钟）"
    ↓
【用户】: "深度讲解"
    ↓
【对话器】: "需求已明确：莫奈《日出·印象》+ 深度讲解15分钟"
    ↓
调用普通工作流生成内容
```

**示例2: 参数补全**

```
【用户】: "讲解化学燃烧"
    ↓
【系统】: 缺少关键参数 → ⭐ 启动交互式对话器
    ↓
【对话器】: "请补充信息：
           1. 目标人群：中学生 / 成人爱好者 / 工程师？
           2. 时长：5分钟 / 10分钟 / 20分钟？
           3. 深度：基础概念 / 详细原理 / 化学方程式？"
    ↓
【用户】: "初三学生，10分钟，需要化学方程式"
    ↓
【对话器】: "收到！为初三学生准备10分钟燃烧讲解，包含方程式"
    ↓
调用普通工作流
```

---

### 场景3️⃣: 内容迭代场景

**触发条件**: 用户对生成的内容不满意，需要调整

**示例: 课件内容调整**

```
【第1轮】普通工作流生成初步内容
    ↓
【用户】: "太学术了，能不能更生活化一点？"
    ↓
【系统】: 检测到调整需求 → ⭐ 切换到交互式对话器
    ↓
【对话器】: 
    - 调用 NarrativeGenerator，style="daily_life"
    - 保留原有学习目标和知识点
    - 调整叙述方式
    ↓
【系统】: 生成调整后的内容（更多生活实例）
    ↓
【用户】: "很好！但是例子能再增加2个吗？"
    ↓
【对话器】:
    - 调用 ContentDiscovery，搜索更多生活案例
    - 保留已有内容
    - 补充2个新例子
    ↓
【系统】: 生成最终版本
    ↓
【用户】: "满意了，导出为PPT"
    ↓
调用 MultimodalGenerator 生成PPT
```

**关键点**:
- ✅ **增量调整**: 不是重新生成，而是在原有基础上修改
- ✅ **保留偏好**: 记住用户的调整历史（"更生活化"）
- ✅ **快速迭代**: 每次调整只改变部分内容

---

### 场景4️⃣: 个性化定制场景

**触发条件**: 用户需要高度定制化的内容

**示例: 企业内训定制**

```
【用户】: "为我们公司财务部做Excel培训"
    ↓
【系统】: 检测到定制需求 → ⭐ 启动交互式对话器
    ↓
【对话器-第1轮】: "了解您的团队：
           1. 人数：___
           2. Excel水平：初级 / 中级 / 高级
           3. 主要使用场景：财务报表 / 数据分析 / 其他"
    ↓
【用户】: "15人，中级水平，主要做财务报表"
    ↓
【对话器-第2轮】: "培训重点（可多选）：
           □ 财务函数（SUM, IF, VLOOKUP）
           □ 数据透视表
           □ 图表制作
           □ 宏和自动化
           □ 财务模板设计"
    ↓
【用户】: "财务函数 + 数据透视表 + 财务模板"
    ↓
【对话器-第3轮】: "培训时长和形式：
           1. 时长：2小时 / 半天 / 全天
           2. 形式：PPT讲解 / 实操演示 / 混合
           3. 是否需要练习数据？"
    ↓
【用户】: "半天，混合形式，需要财务报表练习数据"
    ↓
【对话器-第4轮】: "最后确认：
           ✅ 15人财务部团队
           ✅ 中级Excel水平
           ✅ 培训内容：财务函数 + 透视表 + 模板
           ✅ 半天（4小时）混合形式
           ✅ 提供财务报表练习数据
           
           是否开始生成？"
    ↓
【用户】: "确认"
    ↓
调用普通工作流，传入所有定制参数
```

---

### 场景5️⃣: 教学辅导场景

**触发条件**: 一对一辅导模式，根据学习进度动态调整

**示例: AI化学辅导老师**

```
【学生】: "老师，什么是氧化反应？"
    ↓
【系统】: 检测到教学场景 → ⭐ 启动交互式对话器（辅导模式）
    ↓
【对话器】: 生成基础讲解
    "氧化反应是物质与氧结合的反应。
     比如：铁生锈、木头燃烧都是氧化反应。"
    ↓
【学生】: "那还原反应呢？"
    ↓
【对话器】: 
    - 检测到关联问题
    - 联系上下文（已讲解氧化）
    - 生成对比讲解
    
    "还原反应正好相反，是物质失去氧的反应。
     氧化和还原常常同时发生，比如：
     C + CuO → CO₂ + Cu
     （碳被氧化，氧化铜被还原）"
    ↓
【学生】: "能举个生活中的例子吗？"
    ↓
【对话器】:
    - 检测到需要具象化
    - 调用 ContentDiscovery 搜索生活案例
    
    "好的！比如烤红薯：
     1. 木炭燃烧（氧化）产生热量
     2. 红薯中的淀粉在高温下变化（复杂化学反应）
     
     还有炼铁：
     焦炭夺取铁矿石中的氧（还原），
     铁矿石变成纯铁。"
    ↓
【学生】: "明白了！那考试会怎么考？"
    ↓
【对话器】:
    - 切换到考试辅导模式
    - 生成典型题型
    
    "常见题型：
     1. 判断题：下列反应是否为氧化反应？
     2. 配平题：给出化学方程式，要求配平
     3. 实验题：设计实验验证氧化反应
     
     我给你出一道练习题试试？"
    ↓
【学生】: "好啊"
    ↓
【对话器】: 生成练习题 + 等待答案 + 批改 + 讲解
```

**关键点**:
- ✅ **上下文连贯**: 记住"氧化" → "还原" → "生活案例" → "考试"
- ✅ **教学策略**: 从概念 → 对比 → 具象化 → 应试
- ✅ **动态难度**: 根据学生反应调整深度
- ✅ **苏格拉底式**: 引导学生思考，而非直接给答案

---

## 与普通工作流的区别

### 对比表格

| 维度 | 普通工作流 | 交互式对话器 |
|------|-----------|-------------|
| **输入方式** | 一次性完整输入 | 多轮问答补充 |
| **需求明确度** | 必须清晰明确 | 可以模糊探索 |
| **执行模式** | 线性一次执行 | 循环迭代执行 |
| **调整方式** | 重新执行 | 增量调整 |
| **上下文** | 无需记忆 | 必须记忆历史 |
| **用户体验** | 快速高效 | 友好探索 |
| **适用场景** | 专业用户、明确需求 | 新手用户、探索需求 |
| **技术复杂度** | 简单 | 复杂（状态管理） |
| **API端点** | `/api/workflow/execute` | `/api/chat/message` |
| **数据库** | 不需要会话表 | 需要conversation表 |

### 工作流对比图

```
┌─────────────────────────────────────────────────────────┐
│                   普通工作流模式                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  用户输入："为初三学生讲解化学燃烧条件，10分钟"            │
│      ↓                                                   │
│  领域适配器解析                                          │
│      ↓                                                   │
│  AI能力层处理（意图 → 目标 → 内容 → 叙述）                │
│      ↓                                                   │
│  返回完整内容                                            │
│      ↓                                                   │
│  结束 ✅                                                 │
│                                                          │
│  时间：30秒                                              │
│  交互次数：1次                                            │
└─────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────┐
│                  交互式对话器模式                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  【第1轮】用户："讲化学"                                  │
│      ↓                                                   │
│  系统："中学生还是成人？"                                  │
│      ↓                                                   │
│  【第2轮】用户："中学生"                                  │
│      ↓                                                   │
│  系统："初中还是高中？"                                    │
│      ↓                                                   │
│  【第3轮】用户："初三"                                    │
│      ↓                                                   │
│  系统："想学哪个知识点？"                                  │
│      ↓                                                   │
│  【第4轮】用户："燃烧条件"                                │
│      ↓                                                   │
│  系统："时长：5分钟 / 10分钟 / 15分钟？"                   │
│      ↓                                                   │
│  【第5轮】用户："10分钟"                                  │
│      ↓                                                   │
│  需求收集完成 → 调用普通工作流生成内容                     │
│      ↓                                                   │
│  返回内容                                                │
│      ↓                                                   │
│  【第6轮】用户："太学术了"                                │
│      ↓                                                   │
│  系统：调整叙述风格 → 返回新版本                           │
│      ↓                                                   │
│  【第7轮】用户："满意了"                                  │
│      ↓                                                   │
│  结束 ✅                                                 │
│                                                          │
│  时间：5分钟                                              │
│  交互次数：7轮                                            │
└─────────────────────────────────────────────────────────┘
```

---

## 调用时机决策树

```
接收到用户请求
    ↓
    ┌──────────────────────────────────────┐
    │ 决策1: 需求是否明确完整？             │
    ├──────────────────────────────────────┤
    │ ✅ 是："为初三学生讲解燃烧条件10分钟" │
    │    → 使用普通工作流                   │
    │                                       │
    │ ❌ 否："讲点化学"                     │
    │    → ⭐ 启动交互式对话器（需求探索）  │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策2: 输入是否有歧义？               │
    ├──────────────────────────────────────┤
    │ ✅ 有："讲《日出》" (多个作品同名)    │
    │    → ⭐ 启动交互式对话器（澄清歧义）  │
    │                                       │
    │ ❌ 无："讲莫奈《日出·印象》"          │
    │    → 使用普通工作流                   │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策3: 是否需要多轮迭代？             │
    ├──────────────────────────────────────┤
    │ ✅ 需要：用户对初步结果不满意          │
    │    "太学术了" / "再简单点"            │
    │    → ⭐ 切换到交互式对话器（内容迭代） │
    │                                       │
    │ ❌ 不需要：用户满意初步结果            │
    │    → 使用普通工作流                   │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策4: 是否是教学辅导场景？           │
    ├──────────────────────────────────────┤
    │ ✅ 是："老师，什么是氧化反应？"        │
    │    → ⭐ 启动交互式对话器（辅导模式）  │
    │                                       │
    │ ❌ 否：生成课程内容                   │
    │    → 使用普通工作流                   │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策5: 是否高度定制？                 │
    ├──────────────────────────────────────┤
    │ ✅ 是：企业内训（需收集多个参数）      │
    │    → ⭐ 启动交互式对话器（定制模式）  │
    │                                       │
    │ ❌ 否：标准化内容                     │
    │    → 使用普通工作流                   │
    └──────────────────────────────────────┘
```

### 决策规则总结

| 情况 | 使用模式 | 原因 |
|------|---------|------|
| 需求清晰完整 | 普通工作流 | 效率高，无需多轮 |
| 需求模糊/探索 | ⭐ 交互式对话器 | 逐步收敛需求 |
| 输入有歧义 | ⭐ 交互式对话器 | 需要澄清选择 |
| 缺少关键参数 | ⭐ 交互式对话器 | 补全信息 |
| 对结果不满意 | ⭐ 交互式对话器 | 增量调整 |
| 教学辅导 | ⭐ 交互式对话器 | 动态跟随 |
| 高度定制 | ⭐ 交互式对话器 | 收集多参数 |
| 一次性查询 | 普通工作流 | 直接返回结果 |

---

## 实际案例分析

### 案例1: 新手用户探索学习（完整对话流程）

**背景**: 一位对化学感兴趣的成人用户，没有明确目标

```python
# ========== 完整交互式对话流程 ==========

# 【第1轮】用户启动对话
conversation = await dialog_manager.start_conversation(
    user_id="user_12345",
    initial_message="我想学点化学"
)

# 系统响应：
{
    "conversation_id": "conv_001",
    "message": "很高兴您对化学感兴趣！请问您是：\n1. 中学生准备考试\n2. 成人兴趣学习\n3. 职业技能提升",
    "options": ["中学生", "成人", "职业"],
    "state": "collecting_user_type"
}

# 【第2轮】用户选择
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="成人兴趣学习"
)

# 系统响应：
{
    "message": "明白了！化学领域很广，您更感兴趣：\n1. 日常生活中的化学\n2. 有机化学（药物、材料）\n3. 无机化学（元素、化合物）",
    "options": ["日常生活", "有机化学", "无机化学"],
    "state": "collecting_topic_category",
    "context": {
        "user_type": "成人兴趣学习"
    }
}

# 【第3轮】选择主题
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="日常生活中的化学"
)

# 系统响应：
{
    "message": "太好了！具体想了解：\n1. 食品添加剂的化学原理\n2. 洗涤剂为什么能去污\n3. 电池的工作原理\n4. 其他（请输入）",
    "options": ["食品添加剂", "洗涤剂", "电池", "其他"],
    "state": "collecting_specific_topic",
    "context": {
        "user_type": "成人兴趣学习",
        "category": "日常生活化学"
    }
}

# 【第4轮】确定具体主题
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="食品添加剂"
)

# 系统响应：
{
    "message": "好的！最后确认：\n✅ 成人兴趣学习\n✅ 日常生活化学\n✅ 食品添加剂\n\n时长偏好：5分钟 / 10分钟 / 15分钟？",
    "options": ["5分钟", "10分钟", "15分钟"],
    "state": "collecting_duration",
    "context": {
        "user_type": "成人兴趣学习",
        "category": "日常生活化学",
        "topic": "食品添加剂"
    }
}

# 【第5轮】确定时长
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="10分钟"
)

# 系统响应：
{
    "message": "需求收集完成！开始为您生成内容...",
    "state": "generating_content",
    "context": {
        "user_type": "成人兴趣学习",
        "category": "日常生活化学",
        "topic": "食品添加剂",
        "duration": 10
    }
}

# 【第6轮】生成内容（调用普通工作流）
content = await workflow_engine.execute(
    WorkflowContext(
        domain="chemistry",
        topic="食品添加剂化学原理",
        audience="adult_hobby",
        duration=10,
        style="daily_life"
    )
)

# 返回内容
response = {
    "message": "内容已生成！",
    "content": content,
    "state": "content_ready",
    "actions": ["查看内容", "调整内容", "导出PPT", "结束对话"]
}

# 【第7轮】用户查看后反馈
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="有些化学术语看不懂，能更通俗点吗？"
)

# 系统检测到调整需求 → 调用 NarrativeGenerator
adjusted_content = await narrative_generator.adjust_style(
    original_content=content,
    adjustment="更通俗易懂，减少化学术语，多用生活比喻"
)

# 返回调整后的内容
response = {
    "message": "已为您调整为更通俗的版本！",
    "content": adjusted_content,
    "state": "content_adjusted",
    "actions": ["满意，结束", "继续调整", "导出"]
}

# 【第8轮】用户满意
response = await dialog_manager.handle_turn(
    conversation_id="conv_001",
    user_message="很好，导出为PDF"
)

# 调用多模态生成器
pdf = await multimodal_generator.generate_pdf(
    content=adjusted_content,
    template="adult_learning_handout"
)

# 结束对话
await dialog_manager.end_conversation(
    conversation_id="conv_001",
    final_output=pdf
)
```

**总结**:
- ⏱️ **8轮对话**，从"我想学化学" → PDF讲义
- 🎯 **需求收集**: 第1-5轮逐步收窄范围
- 📝 **内容生成**: 第6轮调用普通工作流
- 🔧 **内容调整**: 第7轮根据反馈优化
- 📄 **格式导出**: 第8轮生成最终产出

---

### 案例2: 歧义消除（《日出》作品）

```python
# ========== 歧义检测和消除流程 ==========

# 【第1轮】用户输入
initial_input = "讲解《日出》这幅画"

# 系统分析：
analysis = await intent_recognizer.analyze(initial_input)
# 结果: {
#     "intent": "art_history_lecture",
#     "entities": [{"artwork": "日出"}],
#     "ambiguity": True,  # ⚠️ 检测到歧义
#     "candidates": [
#         {"artist": "莫奈", "title": "日出·印象", "year": 1872},
#         {"artist": "透纳", "title": "日出与海怪", "year": 1845},
#         {"artist": "弗里德里希", "title": "雾海上的日出", "year": 1818}
#     ]
# }

# 因为有歧义 → 启动交互式对话器
conversation = await dialog_manager.start_conversation(
    user_id="user_67890",
    initial_message=initial_input,
    detected_issue="ambiguity",
    candidates=analysis["candidates"]
)

# 系统响应：
{
    "message": "找到多个《日出》主题的作品：\n\n1. 莫奈《日出·印象》（1872年）\n   - 印象派开山之作\n   - 藏于法国马蒙坦美术馆\n\n2. 透纳《日出与海怪》（1845年）\n   - 英国浪漫主义\n   - 藏于泰特美术馆\n\n3. 弗里德里希《雾海上的日出》（1818年）\n   - 德国浪漫主义\n   - 藏于汉堡美术馆\n\n请问您想了解哪一幅？",
    "options": [
        {"id": 1, "label": "莫奈《日出·印象》"},
        {"id": 2, "label": "透纳《日出与海怪》"},
        {"id": 3, "label": "弗里德里希《雾海上的日出》"}
    ],
    "state": "disambiguating"
}

# 【第2轮】用户选择
response = await dialog_manager.handle_turn(
    conversation_id="conv_002",
    user_message="莫奈的"
)

# 歧义已消除 → 继续收集其他参数
{
    "message": "明白了，《日出·印象》。您希望：\n1. 快速浏览（5分钟）\n2. 深度讲解（15分钟）\n3. 艺术史背景（25分钟）",
    "state": "collecting_depth",
    "context": {
        "artwork": "日出·印象",
        "artist": "莫奈",
        "year": 1872
    }
}

# 【第3轮】选择深度
response = await dialog_manager.handle_turn(
    conversation_id="conv_002",
    user_message="深度讲解"
)

# 需求已完整 → 生成内容
content = await workflow_engine.execute(
    WorkflowContext(
        domain="art_history",
        artist="莫奈",
        artwork="日出·印象",
        depth="deep",
        duration=15
    )
)
```

**价值**:
- ✅ **智能检测歧义**: 不是让用户重新输入，而是提供选项
- ✅ **提供详细信息**: 每个选项附带年代、风格、馆藏信息
- ✅ **快速消歧**: 2轮对话即可确定具体作品

---

### 案例3: 教学辅导模式（化学学习助手）

```python
# ========== AI化学辅导老师 ==========

# 【第1轮】学生提问
conversation = await dialog_manager.start_conversation(
    user_id="student_001",
    mode="tutoring",  # ⭐ 辅导模式
    initial_message="老师，什么是氧化反应？"
)

# 系统分析：
# - 模式：教学辅导
# - 知识点：氧化反应（基础概念）
# - 策略：先简单定义 → 举例 → 引导思考

# 系统响应（第1层：定义）
{
    "message": "氧化反应是物质与氧结合的反应。\n\n简单来说：\n• 铁生锈（铁 + 氧气 → 氧化铁）\n• 木头燃烧（木材 + 氧气 → 二氧化碳 + 水）\n\n这些都是氧化反应。你能想到其他例子吗？",
    "teaching_strategy": "concept_with_examples",
    "expected_response": "学生举例或继续提问",
    "state": "teaching_oxidation_concept"
}

# 【第2轮】学生继续提问
response = await dialog_manager.handle_turn(
    conversation_id="conv_003",
    user_message="那还原反应呢？"
)

# 系统检测：
# - 关联问题（从氧化 → 还原）
# - 教学策略：对比学习

# 系统响应（第2层：对比）
{
    "message": "很好的问题！还原反应正好相反。\n\n【对比】\n氧化反应：得到氧\n还原反应：失去氧\n\n有趣的是，它们常常同时发生！\n\n例如：C + CuO → CO₂ + Cu\n• 碳获得氧 → 被氧化\n• 氧化铜失去氧 → 被还原\n\n这叫做「氧化还原反应」。\n\n你能说说哪个物质被氧化，哪个被还原吗？",
    "teaching_strategy": "comparative_learning_with_quiz",
    "quiz_active": True,
    "state": "teaching_redox_comparison"
}

# 【第3轮】学生回答
response = await dialog_manager.handle_turn(
    conversation_id="conv_003",
    user_message="碳被氧化，氧化铜被还原"
)

# 系统检测：答案正确 → 鼓励 + 深化

# 系统响应（第3层：鼓励+深化）
{
    "message": "✅ 完全正确！你理解得很好！\n\n让我们看个生活中的例子：\n\n【炼铁】\n焦炭（C）夺走铁矿石（Fe₂O₃）中的氧\n→ 铁矿石被还原成纯铁\n→ 焦炭被氧化成二氧化碳\n\n这就是钢铁厂的原理！\n\n你还想知道：\n1. 氧化反应的条件是什么？\n2. 如何判断氧化还原反应？\n3. 考试会怎么考这个知识点？",
    "teaching_strategy": "real_life_application",
    "quiz_result": "correct",
    "encouragement": True,
    "state": "deepening_understanding"
}

# 【第4轮】学生关心考试
response = await dialog_manager.handle_turn(
    conversation_id="conv_003",
    user_message="考试会怎么考？"
)

# 系统切换到考试辅导模式

# 系统响应（第4层：应试）
{
    "message": "考试常见题型：\n\n【题型1】判断题\n「下列反应属于氧化反应的是（  ）」\nA. 2H₂O → 2H₂ + O₂\nB. Fe + O₂ → Fe₃O₄\nC. CaCO₃ → CaO + CO₂\n\n【题型2】配平方程式\n「__Fe + __O₂ → __Fe₃O₄」\n\n【题型3】实验设计\n「设计实验证明铁生锈是氧化反应」\n\n我出一道题给你练习？",
    "teaching_strategy": "exam_preparation",
    "practice_available": True,
    "state": "exam_tutoring"
}

# 【第5轮】学生请求练习
response = await dialog_manager.handle_turn(
    conversation_id="conv_003",
    user_message="好啊，出题"
)

# 生成练习题
{
    "message": "【练习题】\n判断下列反应中，哪些是氧化反应：\n1. 2Mg + O₂ → 2MgO\n2. CuO + H₂ → Cu + H₂O\n3. 2H₂O₂ → 2H₂O + O₂↑\n\n请写出你的答案，我会批改并讲解！",
    "quiz_id": "quiz_001",
    "quiz_type": "multiple_choice",
    "state": "quiz_in_progress"
}

# 【第6轮】学生答题
response = await dialog_manager.handle_turn(
    conversation_id="conv_003",
    user_message="1和3是氧化反应"
)

# 批改答案
{
    "message": "【批改结果】\n✅ 第1题：正确！镁与氧气化合，被氧化\n❌ 第3题：需要注意！这是「分解反应」，不是氧化\n🤔 第2题：你漏掉了！H₂得到氧，被氧化\n\n【详细讲解】\n2. CuO + H₂ → Cu + H₂O\n   • H₂ + O → H₂O（氢被氧化）\n   • CuO - O → Cu（铜被还原）\n   这是氧化还原反应！\n\n3. 2H₂O₂ → 2H₂O + O₂\n   这是过氧化氢分解，不是与外界氧气反应\n\n你理解了吗？需要我再解释哪一题？",
    "quiz_result": {
        "q1": "correct",
        "q2": "missed",
        "q3": "wrong"
    },
    "score": "1/3",
    "state": "quiz_review"
}
```

**辅导模式特点**:
- 🎯 **循序渐进**: 定义 → 对比 → 实例 → 应试
- 💬 **引导式提问**: 不直接给答案，而是引导思考
- 🎓 **即时反馈**: 答题后立即批改和讲解
- 📚 **知识关联**: 从氧化 → 还原 → 氧化还原反应
- ✅ **正向鼓励**: "完全正确！你理解得很好！"
- 🔄 **上下文记忆**: 记住学生学过的概念

---

## 技术实现方案

### 核心架构

```python
# ========== 交互式对话器架构 ==========

class InteractiveDialog:
    """交互式对话管理器"""
    
    def __init__(self):
        # 依赖的其他AI能力
        self.intent_recognizer = IntentRecognizer()
        self.objective_generator = ObjectiveGenerator()
        self.content_discovery = ContentDiscoveryEngine()
        self.narrative_generator = NarrativeGenerator()
        self.quality_assessor = QualityAssessor()
        
        # 会话存储
        self.conversation_store = ConversationStore()
        
        # 状态机
        self.state_machine = DialogStateMachine()
    
    async def start_conversation(
        self,
        user_id: str,
        initial_message: str,
        mode: str = "normal"  # normal, tutoring, customization
    ) -> Conversation:
        """
        启动新对话
        
        Args:
            user_id: 用户ID
            initial_message: 用户的第一句话
            mode: 对话模式（normal普通 / tutoring辅导 / customization定制）
        
        Returns:
            Conversation对象（包含conversation_id, 系统响应等）
        """
        # 1. 创建对话记录
        conversation = await self.conversation_store.create(
            user_id=user_id,
            mode=mode,
            created_at=datetime.now()
        )
        
        # 2. 分析用户输入
        intent_analysis = await self.intent_recognizer.analyze(initial_message)
        
        # 3. 决定是否需要交互式对话
        if intent_analysis.is_complete and not intent_analysis.has_ambiguity:
            # 需求完整清晰 → 直接生成内容
            return await self._direct_generation(conversation, intent_analysis)
        
        # 4. 需求不完整或有歧义 → 启动交互流程
        return await self._start_interactive_flow(
            conversation,
            intent_analysis,
            mode
        )
    
    async def handle_turn(
        self,
        conversation_id: str,
        user_message: str
    ) -> DialogResponse:
        """
        处理对话的一轮
        
        Args:
            conversation_id: 对话ID
            user_message: 用户这一轮的输入
        
        Returns:
            系统响应（包含message, options, state等）
        """
        # 1. 加载对话历史
        conversation = await self.conversation_store.get(conversation_id)
        
        # 2. 添加用户消息到历史
        conversation.add_message(
            role="user",
            content=user_message
        )
        
        # 3. 获取当前状态
        current_state = conversation.state
        
        # 4. 根据状态处理
        if current_state == "collecting_parameters":
            return await self._collect_parameters(conversation, user_message)
        
        elif current_state == "disambiguating":
            return await self._handle_disambiguation(conversation, user_message)
        
        elif current_state == "content_ready":
            return await self._handle_content_feedback(conversation, user_message)
        
        elif current_state == "tutoring":
            return await self._handle_tutoring(conversation, user_message)
        
        else:
            return await self._default_handler(conversation, user_message)
    
    async def _collect_parameters(
        self,
        conversation: Conversation,
        user_input: str
    ) -> DialogResponse:
        """收集缺失的参数"""
        
        # 提取参数
        extracted = await self.intent_recognizer.extract_entities(user_input)
        
        # 更新上下文
        conversation.context.update(extracted)
        
        # 检查还缺什么参数
        missing = self._check_missing_parameters(conversation.context)
        
        if missing:
            # 还有缺失 → 继续询问
            next_question = self._generate_question(missing[0])
            return DialogResponse(
                message=next_question,
                options=self._get_options(missing[0]),
                state="collecting_parameters",
                context=conversation.context
            )
        else:
            # 参数收集完成 → 生成内容
            return await self._generate_content(conversation)
    
    async def _handle_disambiguation(
        self,
        conversation: Conversation,
        user_input: str
    ) -> DialogResponse:
        """处理歧义消除"""
        
        # 识别用户选择
        selected = await self._parse_selection(
            user_input,
            conversation.context["candidates"]
        )
        
        # 更新上下文
        conversation.context["selected_item"] = selected
        conversation.state = "collecting_parameters"
        
        # 继续收集其他参数
        return await self._collect_parameters(conversation, "")
    
    async def _handle_content_feedback(
        self,
        conversation: Conversation,
        user_input: str
    ) -> DialogResponse:
        """处理用户对内容的反馈"""
        
        # 分析反馈意图
        feedback_analysis = await self.intent_recognizer.analyze_feedback(
            user_input
        )
        
        if feedback_analysis.is_satisfied:
            # 用户满意 → 询问是否导出
            return DialogResponse(
                message="很高兴您满意！是否需要：\n1. 导出为PPT\n2. 导出为PDF\n3. 结束对话",
                options=["导出PPT", "导出PDF", "结束"],
                state="finalizing"
            )
        
        elif feedback_analysis.has_adjustment_request:
            # 用户要求调整 → 调用NarrativeGenerator
            adjusted_content = await self.narrative_generator.adjust(
                original=conversation.context["content"],
                adjustment=feedback_analysis.adjustment
            )
            
            conversation.context["content"] = adjusted_content
            
            return DialogResponse(
                message="已为您调整内容！",
                content=adjusted_content,
                options=["满意", "继续调整"],
                state="content_ready"
            )
    
    async def _handle_tutoring(
        self,
        conversation: Conversation,
        user_input: str
    ) -> DialogResponse:
        """辅导模式处理"""
        
        # 1. 识别学生问题类型
        question_type = await self._classify_question(user_input)
        
        # 2. 获取教学策略
        strategy = self._get_teaching_strategy(
            question_type,
            conversation.context["knowledge_graph"]
        )
        
        # 3. 生成教学响应
        if strategy == "concept_explanation":
            # 概念讲解
            response = await self._explain_concept(user_input)
        
        elif strategy == "example_driven":
            # 举例教学
            response = await self._teach_by_examples(user_input)
        
        elif strategy == "quiz":
            # 出题检验
            response = await self._generate_quiz(user_input)
        
        elif strategy == "hint":
            # 引导思考（不直接给答案）
            response = await self._provide_hint(user_input)
        
        # 4. 更新知识图谱（记录学生学过的知识点）
        await self._update_knowledge_graph(
            conversation,
            question=user_input,
            response=response
        )
        
        return DialogResponse(
            message=response,
            state="tutoring",
            teaching_strategy=strategy
        )
    
    async def _generate_content(
        self,
        conversation: Conversation
    ) -> DialogResponse:
        """参数收集完成，生成内容"""
        
        # 1. 构建WorkflowContext
        context = WorkflowContext(
            domain=conversation.context["domain"],
            topic=conversation.context["topic"],
            audience=conversation.context.get("audience_type"),
            duration=conversation.context.get("duration"),
            **conversation.context
        )
        
        # 2. 调用普通工作流生成内容
        content = await workflow_engine.execute(context)
        
        # 3. 保存到对话上下文
        conversation.context["content"] = content
        conversation.state = "content_ready"
        
        # 4. 返回内容
        return DialogResponse(
            message="内容已生成！您可以查看内容，或告诉我需要调整的地方。",
            content=content,
            options=["满意", "太学术了", "太简单了", "增加例子", "导出"],
            state="content_ready"
        )


# ========== 对话状态定义 ==========

class DialogState(Enum):
    """对话状态"""
    COLLECTING_PARAMETERS = "collecting_parameters"  # 收集参数
    DISAMBIGUATING = "disambiguating"  # 消除歧义
    GENERATING_CONTENT = "generating_content"  # 生成内容中
    CONTENT_READY = "content_ready"  # 内容已就绪
    ADJUSTING = "adjusting"  # 调整内容中
    TUTORING = "tutoring"  # 辅导模式
    QUIZ_IN_PROGRESS = "quiz_in_progress"  # 做题中
    FINALIZING = "finalizing"  # 最后确认
    ENDED = "ended"  # 对话结束


# ========== 对话数据模型 ==========

class Conversation(BaseModel):
    """对话模型"""
    id: str
    user_id: str
    mode: str  # normal / tutoring / customization
    state: DialogState
    context: Dict[str, Any]  # 上下文信息
    messages: List[Message]  # 对话历史
    created_at: datetime
    updated_at: datetime
    
    def add_message(self, role: str, content: str):
        """添加消息到历史"""
        self.messages.append(
            Message(
                role=role,  # user / assistant
                content=content,
                timestamp=datetime.now()
            )
        )
        self.updated_at = datetime.now()


class Message(BaseModel):
    """消息模型"""
    role: str  # user / assistant
    content: str
    timestamp: datetime


class DialogResponse(BaseModel):
    """对话响应模型"""
    message: str  # 系统消息
    content: Optional[Any] = None  # 生成的内容（如果有）
    options: Optional[List[str]] = None  # 供用户选择的选项
    state: DialogState  # 当前状态
    context: Optional[Dict[str, Any]] = None  # 上下文信息
    teaching_strategy: Optional[str] = None  # 教学策略（辅导模式）
```

---

### 数据库设计

```sql
-- 对话表
CREATE TABLE conversations (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    mode VARCHAR(20) NOT NULL,  -- normal, tutoring, customization
    state VARCHAR(30) NOT NULL,  -- collecting_parameters, content_ready, etc.
    context JSONB,  -- 上下文信息（JSON格式）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_state (state),
    INDEX idx_created_at (created_at)
);

-- 消息表
CREATE TABLE conversation_messages (
    id SERIAL PRIMARY KEY,
    conversation_id VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL,  -- user, assistant
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (conversation_id) REFERENCES conversations(id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_timestamp (timestamp)
);

-- 知识图谱表（辅导模式用，记录学生学过的知识点）
CREATE TABLE student_knowledge_graph (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    knowledge_point VARCHAR(100) NOT NULL,  -- 知识点名称
    mastery_level INT DEFAULT 0,  -- 掌握程度（0-100）
    last_practiced TIMESTAMP,
    practice_count INT DEFAULT 0,
    
    INDEX idx_user_id (user_id),
    INDEX idx_knowledge_point (knowledge_point)
);
```

---

### API设计

```python
# ========== API路由 ==========

from fastapi import APIRouter, Depends
from typing import Optional

router = APIRouter(prefix="/api/chat", tags=["Interactive Dialog"])

@router.post("/conversations")
async def start_conversation(
    request: StartConversationRequest,
    current_user: User = Depends(get_current_user)
) -> ConversationResponse:
    """
    启动新对话
    
    POST /api/chat/conversations
    {
        "initial_message": "我想学点化学",
        "mode": "normal"  # normal / tutoring / customization
    }
    """
    dialog = InteractiveDialog()
    conversation = await dialog.start_conversation(
        user_id=current_user.id,
        initial_message=request.initial_message,
        mode=request.mode
    )
    return ConversationResponse(conversation=conversation)


@router.post("/conversations/{conversation_id}/messages")
async def send_message(
    conversation_id: str,
    request: SendMessageRequest,
    current_user: User = Depends(get_current_user)
) -> DialogResponse:
    """
    发送消息（对话的一轮）
    
    POST /api/chat/conversations/{conversation_id}/messages
    {
        "message": "成人兴趣学习"
    }
    """
    dialog = InteractiveDialog()
    response = await dialog.handle_turn(
        conversation_id=conversation_id,
        user_message=request.message
    )
    return response


@router.get("/conversations/{conversation_id}")
async def get_conversation(
    conversation_id: str,
    current_user: User = Depends(get_current_user)
) -> Conversation:
    """获取对话详情和历史"""
    store = ConversationStore()
    conversation = await store.get(conversation_id)
    return conversation


@router.delete("/conversations/{conversation_id}")
async def end_conversation(
    conversation_id: str,
    current_user: User = Depends(get_current_user)
) -> dict:
    """结束对话"""
    dialog = InteractiveDialog()
    await dialog.end_conversation(conversation_id)
    return {"message": "Conversation ended"}


@router.get("/conversations")
async def list_conversations(
    current_user: User = Depends(get_current_user),
    limit: int = 20,
    offset: int = 0
) -> List[Conversation]:
    """列出用户的对话列表"""
    store = ConversationStore()
    conversations = await store.list_by_user(
        user_id=current_user.id,
        limit=limit,
        offset=offset
    )
    return conversations
```

---

## 总结

### 调用时机总结表

| 场景 | 触发条件 | 对话轮数 | 价值 |
|------|---------|---------|------|
| **需求探索** | 用户输入模糊（"学点化学"） | 4-6轮 | 帮助用户明确需求 |
| **歧义消除** | 输入有多个解释（"日出"） | 2-3轮 | 快速消歧 |
| **参数补全** | 缺少关键信息 | 3-5轮 | 收集完整参数 |
| **内容迭代** | 对结果不满意 | 1-3轮 | 快速调整，无需重新生成 |
| **定制化** | 企业内训等高度定制 | 5-8轮 | 收集详细需求 |
| **教学辅导** | 一对一学习辅导 | 持续对话 | 动态跟随学习进度 |

### 核心原则

1. **按需启动**: 不是所有请求都需要交互式对话
2. **智能判断**: 系统自动检测是否需要启动对话器
3. **上下文管理**: 必须记住对话历史
4. **状态机驱动**: 清晰的状态转换逻辑
5. **增量调整**: 基于已有内容调整，而非重新生成
6. **用户友好**: 提供选项降低认知负担

### 与普通工作流的配合

```
用户请求
    ↓
【决策点】需求是否清晰完整？
    ↓
    ├─ ✅ 是 → 直接调用普通工作流
    │           ↓
    │         返回结果
    │
    └─ ❌ 否 → ⭐ 启动交互式对话器
                ↓
              多轮对话收集需求
                ↓
              需求完整 ✅
                ↓
              调用普通工作流生成内容
                ↓
              返回内容 → 用户反馈？
                ↓
                ├─ 满意 → 结束
                └─ 不满意 → 继续对话调整
```

### 记忆口诀

```
需求模糊用对话，歧义消除也用它
内容调整很方便，辅导教学更贴心
多轮迭代收需求，上下文中藏智慧
不是每次都启动，智能判断最重要
```

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow Team  
**相关文档**: 
- [多模态生成器调用场景指南](./MULTIMODAL_GENERATOR_GUIDE.md)
- [AI能力分析报告](./AI_CAPABILITIES_ANALYSIS.md)
- [适配器使用指南](./ADAPTER_USAGE_GUIDE.md)
