# 核心概念关系图谱

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 架构说明文档

---

## 📋 目录

1. [四大核心概念](#四大核心概念)
2. [架构层级关系](#架构层级关系)
3. [概念详解](#概念详解)
4. [协作关系](#协作关系)
5. [实际案例](#实际案例)
6. [扩展机制](#扩展机制)

---

## 四大核心概念

### 概念速览表

| 概念 | 英文 | 定位 | 层级 | 数量级 | 生命周期 |
|------|------|------|------|--------|---------|
| **领域适配器** | Domain Adapter | 领域规则翻译器 | 最上层 | 10-50个 | 长期稳定 |
| **AI能力** | AI Capability | 通用AI工具 | 中间层 | 10-20个 | 持续增强 |
| **工作流节点** | Workflow Node | 执行单元 | 底层 | 无限 | 动态组合 |
| **组件** | Component | 可复用功能块 | 扩展层 | 100-1000个 | 市场化 |

### 一句话总结

```
领域适配器（翻译官）调用 AI能力（工具箱），
AI能力驱动 工作流节点（执行单元），
组件（插件）扩展 节点和能力的功能。
```

---

## 架构层级关系

### 完整架构图

```
┌────────────────────────────────────────────────────────────┐
│                      应用层 (API)                          │
│                    用户请求入口                             │
└────────────────────────────────────────────────────────────┘
                           ↓
┌────────────────────────────────────────────────────────────┐
│                  领域适配器层 ⭐                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ K12适配器    │  │ 美术史适配器  │  │ 职业培训适配器│    │
│  │ K12Adapter   │  │ ArtAdapter   │  │ VocAdapter   │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                            │
│  职责：解析用户输入、应用领域规则、选择工作流模板         │
└────────────────────────────────────────────────────────────┘
                           ↓
┌────────────────────────────────────────────────────────────┐
│                   AI能力层 ⭐                              │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐            │
│  │意图识别    │ │学习目标生成│ │内容发现    │            │
│  │Intent      │ │Objective   │ │Content     │            │
│  │Recognizer  │ │Generator   │ │Discovery   │            │
│  └────────────┘ └────────────┘ └────────────┘            │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐            │
│  │故事化叙述  │ │质量评估    │ │多模态生成  │            │
│  │Narrative   │ │Quality     │ │Multimodal  │            │
│  │Generator   │ │Assessor    │ │Generator   │            │
│  └────────────┘ └────────────┘ └────────────┘            │
│                                                            │
│  职责：提供通用AI算法能力（不依赖具体领域）               │
└────────────────────────────────────────────────────────────┘
                           ↓
┌────────────────────────────────────────────────────────────┐
│                 工作流引擎层 ⭐                             │
│                  Workflow Engine                           │
│                                                            │
│  职责：编排节点、管理执行顺序、处理依赖关系                │
└────────────────────────────────────────────────────────────┘
                           ↓
┌────────────────────────────────────────────────────────────┐
│                 工作流节点层 ⭐                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │PersonaNode   │  │ObjectiveNode │  │ContentNode   │    │
│  │角色定义节点  │  │目标生成节点  │  │内容生成节点  │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │NarrativeNode │  │AssessmentNode│  │ExportNode    │    │
│  │叙述生成节点  │  │评估节点      │  │导出节点      │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                            │
│  职责：执行具体任务、调用AI能力、产生输出                  │
└────────────────────────────────────────────────────────────┘
                           ↓
┌────────────────────────────────────────────────────────────┐
│                   组件扩展层 ⭐                             │
│  ┌──────────────────────────────────────────────────┐     │
│  │              组件市场 (Marketplace)               │     │
│  ├──────────────────────────────────────────────────┤     │
│  │ 🔧 内容生成类组件                                 │     │
│  │    - 图片生成器 (ImageGenerator)                 │     │
│  │    - 音频生成器 (AudioGenerator)                 │     │
│  │    - 视频生成器 (VideoGenerator)                 │     │
│  │                                                   │     │
│  │ 📊 质量检查类组件                                 │     │
│  │    - 语法检查器 (GrammarChecker)                 │     │
│  │    - 事实核查器 (FactChecker)                    │     │
│  │                                                   │     │
│  │ 🎯 交互设计类组件                                 │     │
│  │    - 游戏化组件 (GamificationNode)               │     │
│  │    - 虚拟实验室 (VirtualLabNode)                 │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  职责：扩展节点功能、提供专业化能力、支持第三方集成        │
└────────────────────────────────────────────────────────────┘
```

### 调用关系图

```
用户请求："为初三学生讲化学燃烧条件"
    ↓
【1. 领域适配器】K12Adapter
    - 识别：学科=化学，年级=初三，主题=燃烧条件
    - 应用规则：中学化学教学标准
    - 选择模板：K12理科课程模板
    ↓
【2. AI能力层】调用多个AI能力
    ├─ IntentRecognizer：确认意图=教学内容生成
    ├─ ObjectiveGenerator：生成学习目标（3-5个）
    ├─ ContentDiscovery：搜索燃烧条件相关知识
    └─ NarrativeGenerator：组织成课程讲解
    ↓
【3. 工作流引擎】WorkflowEngine
    - 加载模板：K12理科课程工作流
    - 编排节点：PersonaNode → ObjectiveNode → ContentNode → NarrativeNode
    - 管理执行：按DAG顺序执行
    ↓
【4. 工作流节点】按顺序执行
    ├─ PersonaNode：定义学习者画像（初三学生）
    ├─ ObjectiveNode：调用 ObjectiveGenerator 生成目标
    ├─ ContentNode：调用 ContentDiscovery 获取内容
    ├─ NarrativeNode：调用 NarrativeGenerator 生成叙述
    └─ AssessmentNode：生成练习题
    ↓
【5. 组件（可选）】扩展功能
    ├─ ImageGenerator组件：生成实验器材图
    ├─ VirtualLabNode组件：生成虚拟实验
    └─ GamificationNode组件：添加游戏化元素
    ↓
最终输出：完整的化学课程内容
```

---

## 概念详解

### 1️⃣ 领域适配器 (Domain Adapter)

**定义**: 特定领域的"翻译官"和"规则管家"

```python
class K12Adapter(DomainAdapter):
    """K12教育领域适配器"""
    
    domain_name = "K12教育"
    
    # 领域知识
    GRADE_MAPPING = {
        "初一": "7年级",
        "初二": "8年级",
        "初三": "9年级",
        # ...
    }
    
    SUBJECT_STANDARDS = {
        "化学": {
            "curriculum_standard": "义务教育化学课程标准(2022版)",
            "required_experiments": [...],
            "core_concepts": [...]
        }
    }
    
    async def parse_request(self, user_input: str, preferences: Dict) -> UserRequest:
        """
        解析用户输入，提取领域特定信息
        
        输入: "为初三学生讲化学燃烧条件"
        输出: {
            "domain": "K12",
            "subject": "化学",
            "grade": 9,
            "topic": "燃烧条件",
            "audience_type": "初中生",
            "teaching_standard": "义务教育化学课程标准"
        }
        """
        # 1. 提取年级
        grade = self._extract_grade(user_input)
        
        # 2. 提取学科
        subject = self._extract_subject(user_input)
        
        # 3. 提取主题
        topic = self._extract_topic(user_input)
        
        # 4. 应用领域规则
        curriculum = self.SUBJECT_STANDARDS[subject]
        
        return UserRequest(
            domain="K12",
            subject=subject,
            grade=grade,
            topic=topic,
            curriculum_standard=curriculum["curriculum_standard"],
            core_concepts=curriculum["core_concepts"]
        )
    
    async def validate_content(self, content: Content) -> ValidationResult:
        """验证内容是否符合K12教学标准"""
        # 检查是否符合课程标准
        # 检查难度是否适合年级
        # 检查是否包含必要实验
        pass
    
    async def select_workflow_template(self, request: UserRequest) -> str:
        """根据领域规则选择合适的工作流模板"""
        if request.subject in ["物理", "化学", "生物"]:
            return "k12_science_workflow"  # 理科模板（包含实验）
        elif request.subject in ["语文", "英语"]:
            return "k12_language_workflow"  # 文科模板
        else:
            return "k12_general_workflow"  # 通用模板
```

**关键特征**:
- ✅ **领域专属**: 一个领域一个适配器（K12、美术史、职业培训等）
- ✅ **规则管理**: 封装领域规则、标准、知识
- ✅ **输入解析**: 从自然语言提取结构化信息
- ✅ **模板选择**: 根据领域特点选择工作流模板
- ✅ **内容验证**: 确保输出符合领域标准

---

### 2️⃣ AI能力 (AI Capability)

**定义**: 通用的AI算法能力，不依赖具体领域

```python
class ObjectiveGenerator:
    """学习目标生成器（AI能力）"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def generate(self, context: Dict[str, Any]) -> List[LearningObjective]:
        """
        生成学习目标
        
        输入: {
            "topic": "燃烧条件",
            "audience_type": "初中生",
            "duration": 10
        }
        
        输出: [
            "理解燃烧的三个条件：可燃物、氧气、达到着火点",
            "能够解释为什么需要同时满足三个条件",
            "掌握灭火的三种方法及其原理"
        ]
        """
        prompt = self._build_prompt(context)
        response = await self.llm.generate(prompt)
        objectives = self._parse_objectives(response)
        
        # 验证目标符合SMART原则
        validated = self._validate_smart(objectives)
        
        return validated
    
    def _validate_smart(self, objectives: List[str]) -> List[LearningObjective]:
        """
        验证学习目标是否符合SMART原则
        S - Specific（具体）
        M - Measurable（可测量）
        A - Achievable（可达成）
        R - Relevant（相关）
        T - Time-bound（有时限）
        """
        # 使用Bloom分类法验证
        # 检查是否使用了可测量的动词（理解、掌握、应用）
        pass
```

**AI能力清单**:

| AI能力 | 功能 | 输入 | 输出 | 领域依赖 |
|--------|------|------|------|---------|
| **IntentRecognizer** | 意图识别 | 自然语言文本 | 意图+实体 | ❌ 无 |
| **ObjectiveGenerator** | 学习目标生成 | 主题+受众 | 学习目标列表 | ❌ 无 |
| **ContentDiscovery** | 内容发现 | 关键词 | 相关内容 | ❌ 无 |
| **NarrativeGenerator** | 故事化叙述 | 知识点 | 叙述文本 | ❌ 无 |
| **QualityAssessor** | 质量评估 | 内容 | 评分+建议 | ❌ 无 |
| **MultimodalGenerator** | 多模态生成 | 描述 | 图/音/视频 | ❌ 无 |
| **InteractiveDialog** | 交互对话 | 用户消息 | 系统响应 | ❌ 无 |

**关键特征**:
- ✅ **通用性**: 不包含领域特定逻辑
- ✅ **可复用**: 所有领域都可以使用
- ✅ **独立性**: 可以单独测试和优化
- ✅ **AI驱动**: 基于大语言模型或专业算法

---

### 3️⃣ 工作流节点 (Workflow Node)

**定义**: 工作流中的执行单元，封装具体任务

```python
class ObjectiveNode(WorkflowNode):
    """学习目标生成节点"""
    
    node_type = "objective_generation"
    
    def __init__(self, objective_generator: ObjectiveGenerator):
        self.objective_generator = objective_generator
    
    async def execute(self, context: WorkflowContext, input_data: Dict) -> Dict:
        """
        执行节点任务
        
        输入（来自上游节点）: {
            "persona": {
                "age_group": "初中生",
                "grade": 9
            },
            "topic": "燃烧条件",
            "duration": 10
        }
        
        输出（传给下游节点）: {
            "objectives": [
                "理解燃烧的三个条件",
                "能够解释为什么需要同时满足三个条件",
                "掌握灭火的三种方法"
            ]
        }
        """
        # 1. 调用AI能力
        objectives = await self.objective_generator.generate(
            topic=input_data["topic"],
            audience_type=input_data["persona"]["age_group"],
            duration=input_data["duration"]
        )
        
        # 2. 记录执行结果
        await self._save_execution_record(context, objectives)
        
        # 3. 返回结果（供下游节点使用）
        return {
            "objectives": objectives
        }
```

**节点分类**:

```
工作流节点分类
├── 输入节点（Input Nodes）
│   ├── PersonaNode - 定义学习者画像
│   ├── RequirementNode - 收集需求信息
│   └── ConstraintNode - 设置约束条件
│
├── 处理节点（Processing Nodes）
│   ├── ObjectiveNode - 生成学习目标
│   ├── ContentNode - 发现和组织内容
│   ├── NarrativeNode - 生成叙述
│   ├── ActivityNode - 设计教学活动
│   └── AssessmentNode - 生成评估题
│
├── 控制节点（Control Nodes）
│   ├── ConditionNode - 条件分支
│   ├── LoopNode - 循环迭代
│   └── ParallelNode - 并行执行
│
└── 输出节点（Output Nodes）
    ├── FormatNode - 格式化输出
    ├── ExportNode - 导出文件
    └── PublishNode - 发布内容
```

**关键特征**:
- ✅ **单一职责**: 每个节点只做一件事
- ✅ **可组合**: 通过DAG组合成复杂流程
- ✅ **标准接口**: 统一的execute()方法
- ✅ **状态管理**: 记录执行状态和结果

---

### 4️⃣ 组件 (Component)

**定义**: 可插拔的功能扩展模块，可以是新节点或能力增强

```python
# ========== 组件示例1：新节点类型 ==========

class VirtualLabNode(WorkflowNode):
    """虚拟实验室节点（组件）"""
    
    node_type = "virtual_lab"
    category = "interaction_design"  # 组件分类
    
    # 组件元数据
    metadata = {
        "name": "虚拟化学实验室",
        "version": "1.0.0",
        "author": "张三",
        "description": "提供交互式化学实验模拟",
        "dependencies": ["three.js", "chemistry-engine"],
        "tags": ["实验", "化学", "3D"],
        "rating": 4.8,
        "downloads": 1523
    }
    
    async def execute(self, context: WorkflowContext, input_data: Dict) -> Dict:
        """
        生成虚拟实验室
        
        输入: {
            "experiment_type": "燃烧实验",
            "equipment": ["烧杯", "蜡烛", "火柴"],
            "safety_level": "low_risk"
        }
        
        输出: {
            "virtual_lab_url": "https://lab.example.com/combustion",
            "interaction_script": {...},
            "safety_tips": [...]
        }
        """
        # 1. 生成3D场景
        scene = await self._generate_3d_scene(input_data)
        
        # 2. 添加交互逻辑
        interactions = await self._add_interactions(input_data)
        
        # 3. 配置安全提示
        safety = await self._configure_safety(input_data)
        
        return {
            "virtual_lab_url": scene.url,
            "interaction_script": interactions,
            "safety_tips": safety
        }


# ========== 组件示例2：能力增强 ==========

class ImageGeneratorComponent:
    """图片生成组件（增强MultimodalGenerator）"""
    
    component_type = "capability_extension"
    extends = "MultimodalGenerator"  # 扩展现有AI能力
    
    metadata = {
        "name": "DALL-E 3 图片生成器",
        "version": "2.0.0",
        "provider": "OpenAI",
        "pricing": "按次计费"
    }
    
    async def generate_image(self, prompt: str, **options) -> Image:
        """增强图片生成能力"""
        # 调用DALL-E 3 API
        image = await openai.images.generate(
            model="dall-e-3",
            prompt=prompt,
            size=options.get("size", "1024x1024"),
            quality=options.get("quality", "standard")
        )
        return image
```

**组件分类**:

| 类别 | 说明 | 示例 |
|------|------|------|
| **内容生成类** | 生成特定类型内容 | 图片生成器、视频生成器、PPT生成器 |
| **质量检查类** | 检查和优化内容 | 语法检查器、事实核查器、抄袭检测 |
| **交互设计类** | 增强交互体验 | 虚拟实验室、游戏化组件、AR组件 |
| **数据集成类** | 对接外部数据源 | Wikipedia连接器、YouTube API |
| **能力增强类** | 扩展AI能力 | 特定模型集成、专业算法 |

**关键特征**:
- ✅ **插件化**: 不修改核心代码即可安装
- ✅ **市场化**: 可发布到组件市场供他人使用
- ✅ **版本管理**: 支持多版本共存
- ✅ **依赖声明**: 明确声明依赖关系
- ✅ **热加载**: 运行时动态加载

---

## 协作关系

### 关系矩阵

| 关系 | 领域适配器 | AI能力 | 工作流节点 | 组件 |
|------|-----------|--------|-----------|------|
| **调用** AI能力 | ✅ 调用 | - | ✅ 调用 | ❌ 不直接调用 |
| **使用** 工作流模板 | ✅ 选择模板 | ❌ 不涉及 | - | ❌ 不涉及 |
| **执行** 工作流节点 | ❌ 不执行 | ❌ 不执行 | - | - |
| **扩展** 现有功能 | ❌ 很少扩展 | ⚠️ 可被扩展 | ⚠️ 可被扩展 | ✅ 主要用途 |
| **独立存在** | ✅ 可独立 | ✅ 可独立 | ❌ 必须在工作流中 | ✅ 可独立发布 |

### 依赖关系图

```
┌────────────────────────────────────────────────┐
│          领域适配器 (Domain Adapter)           │
│          - 不依赖其他三者                      │
│          - 但会选择工作流模板                  │
└────────────────────────────────────────────────┘
            ↓ (选择模板)
┌────────────────────────────────────────────────┐
│            工作流模板 (Workflow Template)       │
│            - 定义节点序列和依赖                 │
└────────────────────────────────────────────────┘
            ↓ (实例化)
┌────────────────────────────────────────────────┐
│          工作流节点 (Workflow Node)            │
│          - 依赖 AI能力                         │
│          - 可选依赖 组件                       │
└────────────────────────────────────────────────┘
            ↓ (调用)
┌─────────────────┐          ┌─────────────────┐
│   AI能力        │          │   组件          │
│ AI Capability   │  ←扩展── │  Component      │
│ - 独立实现      │          │  - 可扩展能力   │
└─────────────────┘          └─────────────────┘
```

### 协作流程详解

#### 场景：生成K12化学课程

```
┌──────────────────────────────────────────────────────────┐
│ 步骤1: 用户输入 → 领域适配器                              │
├──────────────────────────────────────────────────────────┤
│ 输入: "为初三学生讲解化学燃烧条件，10分钟"                │
│   ↓                                                       │
│ K12Adapter.parse_request()                               │
│   - 提取：年级=9，学科=化学，主题=燃烧条件                │
│   - 应用规则：《义务教育化学课程标准》                    │
│   - 选择模板：k12_science_workflow                       │
│   ↓                                                       │
│ 输出: UserRequest {                                      │
│   domain: "K12",                                         │
│   subject: "化学",                                       │
│   grade: 9,                                              │
│   topic: "燃烧条件",                                     │
│   template: "k12_science_workflow"                      │
│ }                                                        │
└──────────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────────┐
│ 步骤2: 工作流引擎加载模板                                 │
├──────────────────────────────────────────────────────────┤
│ WorkflowEngine.load_template("k12_science_workflow")    │
│   ↓                                                       │
│ 工作流DAG:                                                │
│   PersonaNode → ObjectiveNode → ContentNode              │
│       ↓              ↓             ↓                      │
│   NarrativeNode → ExperimentNode → AssessmentNode       │
│       ↓                                                   │
│   FormatNode                                             │
└──────────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────────┐
│ 步骤3: 执行节点（逐个）                                   │
├──────────────────────────────────────────────────────────┤
│ 【节点1】PersonaNode                                      │
│   - 输入：grade=9, subject=化学                          │
│   - 处理：定义学习者画像                                  │
│   - 输出：{"age_group": "初中生", "knowledge_level": 3}  │
│                                                           │
│ 【节点2】ObjectiveNode ⭐ 调用AI能力                      │
│   - 输入：persona + topic                                │
│   - 调用：ObjectiveGenerator.generate()                 │
│   - AI能力处理：                                          │
│     * 分析主题：燃烧条件                                  │
│     * 生成目标：[理解燃烧三要素, 掌握灭火方法, ...]      │
│   - 输出：{"objectives": [目标1, 目标2, 目标3]}          │
│                                                           │
│ 【节点3】ContentNode ⭐ 调用AI能力                        │
│   - 输入：objectives                                     │
│   - 调用：ContentDiscovery.search()                     │
│   - AI能力处理：                                          │
│     * 搜索Wikipedia：燃烧、氧化反应                      │
│     * 搜索教材库：初中化学教材相关章节                    │
│     * 融合排序：RRF算法                                  │
│   - 输出：{"content_chunks": [内容1, 内容2, ...]}       │
│                                                           │
│ 【节点4】NarrativeNode ⭐ 调用AI能力                      │
│   - 输入：objectives + content_chunks                    │
│   - 调用：NarrativeGenerator.generate()                 │
│   - AI能力处理：                                          │
│     * 选择策略：情境嵌入式（生活化场景）                  │
│     * 生成叙述："想象你在野营，生篝火需要什么？..."      │
│   - 输出：{"narrative": "完整讲解文本"}                  │
│                                                           │
│ 【节点5】ExperimentNode ⭐ 可选调用组件                   │
│   - 输入：topic=燃烧条件                                 │
│   - 检查：是否安装了VirtualLabNode组件？                │
│   - 如果已安装：                                          │
│     * 调用组件：VirtualLabNode.execute()                │
│     * 生成：虚拟实验室URL                                │
│   - 如果未安装：                                          │
│     * 跳过或生成文字版实验步骤                            │
│   - 输出：{"experiment": {...}}                         │
│                                                           │
│ 【节点6】AssessmentNode ⭐ 调用AI能力                     │
│   - 输入：objectives + narrative                         │
│   - 调用：ObjectiveGenerator.generate_questions()       │
│   - AI能力处理：生成练习题                               │
│   - 输出：{"questions": [题1, 题2, 题3]}                │
│                                                           │
│ 【节点7】FormatNode                                       │
│   - 输入：所有上游节点的输出                              │
│   - 处理：格式化为最终课程结构                            │
│   - 输出：完整的课程JSON                                 │
└──────────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────────┐
│ 步骤4: 返回结果                                           │
├──────────────────────────────────────────────────────────┤
│ {                                                        │
│   "title": "化学燃烧条件",                               │
│   "objectives": [...],                                   │
│   "content": "想象你在野营...",                          │
│   "experiment": {虚拟实验室},                            │
│   "questions": [...]                                     │
│ }                                                        │
└──────────────────────────────────────────────────────────┘
```

### 关键协作模式

#### 模式1: 领域适配器 → 工作流模板选择

```python
# K12适配器根据学科选择不同模板

if subject in ["物理", "化学", "生物"]:
    template = "k12_science_workflow"  
    # 包含：实验节点、安全提示节点
    
elif subject in ["数学"]:
    template = "k12_math_workflow"
    # 包含：例题节点、练习题节点、可视化节点
    
elif subject in ["语文", "英语"]:
    template = "k12_language_workflow"
    # 包含：朗读节点、词汇节点、写作节点
```

#### 模式2: 工作流节点 → AI能力调用

```python
class ContentNode(WorkflowNode):
    """内容生成节点"""
    
    def __init__(self):
        # 依赖注入AI能力
        self.content_discovery = ContentDiscoveryEngine()
        self.narrative_generator = NarrativeGenerator()
    
    async def execute(self, context, input_data):
        # 1. 调用AI能力1：发现内容
        raw_content = await self.content_discovery.search(
            keywords=input_data["keywords"]
        )
        
        # 2. 调用AI能力2：组织叙述
        narrative = await self.narrative_generator.generate(
            content=raw_content,
            style="educational"
        )
        
        return {"narrative": narrative}
```

#### 模式3: 组件扩展节点功能

```python
# 原始节点
class ExperimentNode(WorkflowNode):
    async def execute(self, context, input_data):
        # 默认：生成文字版实验步骤
        return {"experiment_text": "实验步骤：1. ..."}

# 组件扩展后
class ExperimentNode(WorkflowNode):
    def __init__(self):
        # 检查是否安装了虚拟实验室组件
        self.virtual_lab = self._load_component("VirtualLabNode")
    
    async def execute(self, context, input_data):
        if self.virtual_lab:
            # 使用组件：生成3D虚拟实验室
            return await self.virtual_lab.execute(context, input_data)
        else:
            # 默认：生成文字版
            return {"experiment_text": "实验步骤：1. ..."}
```

#### 模式4: 组件扩展AI能力

```python
# 原始AI能力
class MultimodalGenerator:
    async def generate_image(self, prompt):
        # 默认：使用Stable Diffusion
        return await self.sd_client.generate(prompt)

# 组件扩展后
class MultimodalGenerator:
    def __init__(self):
        # 加载图片生成组件
        self.providers = self._load_image_components()
        # providers = {
        #     "dall-e-3": DALLE3Component,
        #     "midjourney": MidjourneyComponent,
        #     "stable-diffusion": SDComponent
        # }
    
    async def generate_image(self, prompt, provider="stable-diffusion"):
        # 根据配置选择不同的组件
        component = self.providers.get(provider)
        if component:
            return await component.generate(prompt)
        else:
            # 降级到默认
            return await self.default_generator(prompt)
```

---

## 实际案例

### 案例1: K12化学课程（完整流程）

```
用户输入：
"为初三学生讲解化学燃烧条件，需要包含实验演示，10分钟"

┌─────────────────────────────────────────────────┐
│ 1. 领域适配器处理                                │
├─────────────────────────────────────────────────┤
│ K12Adapter:                                     │
│   - 识别：年级=9，学科=化学，主题=燃烧条件       │
│   - 应用规则：                                   │
│     * 课程标准：义务教育化学课程标准(2022)       │
│     * 必修实验：燃烧条件探究实验                 │
│     * 安全要求：需添加安全提示                   │
│   - 选择模板：k12_science_workflow              │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│ 2. 工作流引擎执行                                │
├─────────────────────────────────────────────────┤
│ 工作流DAG（7个节点）：                           │
│                                                  │
│ PersonaNode                                     │
│    ↓                                             │
│ ObjectiveNode ← (调用ObjectiveGenerator)        │
│    ↓                                             │
│ ContentNode ← (调用ContentDiscovery)            │
│    ↓                                             │
│ NarrativeNode ← (调用NarrativeGenerator)        │
│    ↓                                             │
│ ExperimentNode ← (调用VirtualLabNode组件)       │
│    ↓                                             │
│ SafetyNode ← (K12特有节点)                      │
│    ↓                                             │
│ AssessmentNode ← (调用ObjectiveGenerator)       │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│ 3. 各节点执行详情                                │
├─────────────────────────────────────────────────┤
│ PersonaNode:                                    │
│   输出 → {age: 14-15, grade: 9, knowledge: 3}  │
│                                                  │
│ ObjectiveNode (AI能力):                         │
│   输入 ← persona + topic                        │
│   AI生成 → [                                    │
│     "理解燃烧的三个条件",                        │
│     "能够解释灭火原理",                          │
│     "掌握安全用火知识"                           │
│   ]                                              │
│                                                  │
│ ContentNode (AI能力):                           │
│   输入 ← objectives                             │
│   AI搜索 → Wikipedia + 教材库 + 视频库           │
│   融合 → RRF算法排序                            │
│   输出 → 结构化内容块                            │
│                                                  │
│ NarrativeNode (AI能力):                         │
│   输入 ← content                                │
│   AI叙述 → "想象你在野营生篝火..."              │
│   策略 → 情境嵌入式（生活化）                    │
│                                                  │
│ ExperimentNode (组件):                          │
│   检测到安装了VirtualLabNode组件                │
│   调用组件 → 生成3D虚拟化学实验室               │
│   输出 → {                                       │
│     "lab_url": "https://...",                   │
│     "equipment": ["蜡烛", "玻璃罩", "火柴"],     │
│     "steps": [...]                              │
│   }                                              │
│                                                  │
│ SafetyNode (K12特有):                           │
│   根据实验生成安全提示                           │
│   输出 → [                                       │
│     "实验需在老师指导下进行",                    │
│     "注意火源安全",                              │
│     "保持通风"                                   │
│   ]                                              │
│                                                  │
│ AssessmentNode (AI能力):                        │
│   根据objectives生成练习题                       │
│   输出 → [选择题3道, 简答题2道]                  │
└─────────────────────────────────────────────────┘
                    ↓
最终输出：完整的K12化学课程（包含虚拟实验室）
```

### 案例2: 美术史讲解（不同领域）

```
用户输入：
"讲解莫奈的《日出·印象》"

┌─────────────────────────────────────────────────┐
│ 1. 领域适配器处理                                │
├─────────────────────────────────────────────────┤
│ ArtHistoryAdapter:                              │
│   - 识别：艺术家=莫奈，作品=日出·印象            │
│   - 应用规则：                                   │
│     * 流派：印象派                               │
│     * 时代背景：19世纪法国                       │
│     * 必须包含：画作高清图、艺术技法分析         │
│   - 选择模板：art_history_workflow              │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│ 2. 工作流引擎执行                                │
├─────────────────────────────────────────────────┤
│ 工作流DAG（5个节点）：                           │
│                                                  │
│ ArtworkInfoNode                                 │
│    ↓                                             │
│ ContentNode ← (调用ContentDiscovery)            │
│    ↓                                             │
│ NarrativeNode ← (调用NarrativeGenerator)        │
│    ↓                                             │
│ ImageNode ← (调用MultimodalGenerator组件)       │
│    ↓                                             │
│ FormatNode                                      │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│ 3. 关键差异（vs K12）                            │
├─────────────────────────────────────────────────┤
│ ❌ 没有SafetyNode（美术史不需要）                │
│ ❌ 没有ExperimentNode（不涉及实验）              │
│ ✅ 有ImageNode（高清作品图必须）                 │
│ ✅ NarrativeNode使用不同策略（艺术鉴赏风格）     │
└─────────────────────────────────────────────────┘
```

**关键点**: 
- 同样的AI能力（ContentDiscovery, NarrativeGenerator）
- 不同的领域适配器（K12Adapter vs ArtHistoryAdapter）
- 不同的工作流模板（science vs art_history）
- 不同的节点组合（实验 vs 图片）

---

## 扩展机制

### 如何扩展系统？

```
扩展类型决策树
    ↓
需要新的领域规则？
    ├─ 是 → 创建新的领域适配器
    │       例：VocationalAdapter（职业培训）
    │
    └─ 否 → 需要新的AI算法能力？
            ├─ 是 → 扩展AI能力层
            │       例：PersonalizationEngine（个性化推荐）
            │
            └─ 否 → 需要新的执行逻辑？
                    ├─ 是 → 创建新的工作流节点
                    │       例：SimulationNode（模拟仿真）
                    │
                    └─ 否 → 想复用/共享功能？
                            └─ 是 → 创建组件发布到市场
                                    例：AR组件、游戏化组件
```

### 扩展方式对比

| 扩展方式 | 何时使用 | 工作量 | 复用性 | 示例 |
|---------|---------|--------|--------|------|
| **新增领域适配器** | 进入新领域 | 中等 | 低（领域专属） | VocationalAdapter |
| **扩展AI能力** | 需要新的AI算法 | 大 | 高（所有领域可用） | PersonalizationEngine |
| **新增工作流节点** | 新的执行逻辑 | 小 | 中（特定流程使用） | SimulationNode |
| **开发组件** | 可复用功能 | 小-中 | 高（市场化） | VirtualLabNode |

### 实际扩展示例

#### 示例1: 新增职业培训领域

```python
# 步骤1: 创建领域适配器
class VocationalAdapter(DomainAdapter):
    """职业培训领域适配器"""
    
    domain_name = "职业培训"
    
    SKILL_CATEGORIES = {
        "办公软件": ["Excel", "Word", "PowerPoint"],
        "编程语言": ["Python", "JavaScript", "Java"],
        "设计工具": ["Photoshop", "Figma", "AutoCAD"]
    }
    
    async def parse_request(self, user_input):
        # 提取：培训类型、技能、职业等级
        pass
    
    async def select_workflow_template(self, request):
        if request.training_type == "实操":
            return "vocational_hands_on_workflow"  # 包含ScreenRecordingNode
        else:
            return "vocational_theory_workflow"

# 步骤2: 创建专属工作流模板
# config/workflows/vocational_hands_on.yaml
workflow:
  name: 职业培训-实操类
  nodes:
    - id: persona
      type: persona_node
    - id: objective
      type: objective_node
    - id: content
      type: content_node
    - id: demo_video  # ⭐ 专属节点
      type: screen_recording_node
    - id: practice
      type: practice_node
    - id: assessment
      type: assessment_node

# 步骤3: 可选 - 开发专属组件
class ScreenRecordingNode(WorkflowNode):
    """屏幕录制节点（组件）"""
    
    async def execute(self, context, input_data):
        # 生成软件操作演示视频
        return {"demo_video_url": "..."}
```

#### 示例2: 扩展个性化推荐能力

```python
# 步骤1: 创建新的AI能力
class PersonalizationEngine:
    """个性化推荐引擎（新AI能力）"""
    
    async def recommend_content(
        self,
        user_profile: Dict,
        learning_history: List,
        current_topic: str
    ) -> List[Content]:
        """
        根据用户画像和学习历史推荐内容
        
        算法：
        - 协同过滤
        - 内容相似度
        - 知识图谱推理
        """
        # 分析用户偏好
        preferences = self._analyze_preferences(user_profile, learning_history)
        
        # 推荐相似内容
        recommendations = self._collaborative_filtering(current_topic, preferences)
        
        return recommendations

# 步骤2: 在节点中使用新能力
class PersonalizedContentNode(WorkflowNode):
    """个性化内容节点"""
    
    def __init__(self):
        self.content_discovery = ContentDiscoveryEngine()
        self.personalization = PersonalizationEngine()  # ⭐ 使用新能力
    
    async def execute(self, context, input_data):
        # 1. 常规内容发现
        general_content = await self.content_discovery.search(...)
        
        # 2. 个性化过滤和排序
        personalized = await self.personalization.recommend_content(
            user_profile=context.user_profile,
            learning_history=context.learning_history,
            current_topic=input_data["topic"]
        )
        
        return {"content": personalized}
```

---

## 总结

### 记忆口诀

```
领域适配器是翻译官，解析需求选模板
AI能力是工具箱，通用算法可复用
工作流节点是执行兵，一个任务一个兵
组件是插件库，市场化扩展功能

适配器调用不了节点，节点必须调能力
能力独立可测试，组件扩展最灵活
```

### 快速决策表

**我要实现一个新功能，应该在哪一层？**

| 功能描述 | 应该实现为 | 理由 |
|---------|-----------|------|
| 支持新领域（如烹饪） | 领域适配器 | 需要领域特定规则 |
| 新的AI算法（如知识图谱推理） | AI能力 | 通用算法，可复用 |
| 新的执行逻辑（如3D渲染） | 工作流节点 | 特定任务执行 |
| 可选功能（如AR展示） | 组件 | 可插拔、可分享 |

### 核心原则再强调

1. **领域适配器**: 一个领域一个，封装领域知识
2. **AI能力**: 通用算法，不包含领域逻辑
3. **工作流节点**: 单一职责，可组合
4. **组件**: 插件化，市场化

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow Team  

**相关文档**:
- [领域适配器使用指南](./ADAPTER_USAGE_GUIDE.md)
- [AI能力分析报告](./AI_CAPABILITIES_ANALYSIS.md)
- [组件扩展指南](./COMPONENT_EXTENSION_GUIDE.md)
- [多模态生成器调用场景](./MULTIMODAL_GENERATOR_GUIDE.md)
- [交互式对话器调用场景](./INTERACTIVE_DIALOG_GUIDE.md)
