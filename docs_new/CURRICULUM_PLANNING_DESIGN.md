# 教材级课程规划·元工作流设计方案

**文档版本**: V1.0  
**日期**: 2025-12-09  
**场景**: 整本教材的课程体系规划

---

## 📋 需求概述

### 背景

当前系统专注于**单节课**的内容生成（粒度：45分钟），而实际教学中需要从**整本教材**层面进行系统规划。

### 目标

设计一个**两层元工作流架构**：

```
┌─────────────────────────────────────────────────────────┐
│           上层：教材级元工作流 (Curriculum Planning)      │
│   输入：9年级·化学·上册                                   │
│   输出：课程规划、主题划分、故事线、课型序列              │
└─────────────────────────────────────────────────────────┘
                          ↓ 驱动
┌─────────────────────────────────────────────────────────┐
│           下层：单课级元工作流 (Lesson Generation)        │
│   输入：主题1-第1课·新授课                                │
│   输出：具体的课件内容                                    │
└─────────────────────────────────────────────────────────┘
```

---

## 🎯 核心功能需求

### 1. 教材分析与主题划分

**输入**：
- 教材信息：年级、学科、册别、出版社
- 课程标准文档
- 教学日历（总课时、周课时）

**输出**：
- **主题列表**（章节、单元）
  - 主题名称
  - 知识点列表
  - 课时分配
  - 难度等级
  - 前置知识依赖

**示例**：
```yaml
# 9年级化学上册·人教版
curriculum:
  subject: "化学"
  grade: 9
  semester: "上册"
  publisher: "人教版"
  total_periods: 72  # 总课时
  
  units:
    - unit_id: "unit_01"
      name: "走进化学世界"
      periods: 8
      topics:
        - "物质的变化和性质"
        - "化学实验基本方法"
        - "走进化学实验室"
      difficulty: "基础"
      prerequisites: []
    
    - unit_id: "unit_02"
      name: "我们周围的空气"
      periods: 10
      topics:
        - "空气的组成"
        - "氧气的性质"
        - "制取氧气"
      difficulty: "中等"
      prerequisites: ["unit_01"]
```

### 2. 课时序列规划

**功能**：
- 自动生成整学期的课程序列
- 合理安排课型（新授课、复习课、实验课、练习课、测验）
- 考虑学习曲线和遗忘曲线

**课型分类**：
```python
class LessonType(str, Enum):
    NEW = "新授课"           # 讲授新知识
    REVIEW = "复习课"        # 巩固旧知识
    EXPERIMENT = "实验课"    # 动手实验
    PRACTICE = "练习课"      # 习题训练
    INTERACTIVE = "互动课"   # 讨论、游戏
    ASSESSMENT = "测验课"    # 阶段测试
```

**排课规则**：
```yaml
scheduling_rules:
  # 新授课后的巩固策略
  - after: "NEW"
    follow_within: 2  # 2课时内
    follow_type: "PRACTICE"
    probability: 0.8
  
  # 单元结束前的复习
  - before: "UNIT_END"
    precede_within: 3
    precede_type: "REVIEW"
    required: true
  
  # 实验课的位置
  - type: "EXPERIMENT"
    position: "after_theory"  # 理论讲解后
    gap: 1  # 间隔1课时
```

### 3. 故事世界观与主线设计

**场景**：
- 为整本教材设计统一的故事世界观
- 每个单元是故事的一个篇章
- 每节课是故事的一个情节

**设计要素**：

```yaml
story_universe:
  theme: "化学侦探社"
  setting: "未来都市·元素市"
  
  main_characters:
    - name: "小凯"
      role: "主角·化学侦探"
      personality: "好奇、严谨"
    
    - name: "元素精灵"
      role: "导师·神秘生物"
      personality: "睿智、幽默"
  
  story_arc:
    # 主线：揭开元素市的秘密
    main_plot:
      act1: "神秘事件（第1-3单元）"
      act2: "追寻真相（第4-6单元）"
      act3: "终极对决（第7-8单元）"
    
    # 支线：每个单元的小故事
    sub_plots:
      - unit: "unit_01"
        title: "实验室的奇怪现象"
        conflict: "物质变化之谜"
        resolution: "理解化学变化本质"
      
      - unit: "unit_02"
        title: "消失的氧气"
        conflict: "空气成分被破坏"
        resolution: "掌握氧气制取方法"
```

**故事连贯性**：
- 每节课的故事要承接上一节
- 重要发现在关键节点揭示
- 角色成长与知识掌握同步

### 4. 知识图谱与依赖关系

**构建教材级知识图谱**：

```
                 物质的变化
                      ↓
        ┌─────────────┴─────────────┐
   化学变化                     物理变化
        ↓                             ↓
   氧化反应 ─────→ 燃烧 ────→ 灭火原理
        ↓
   金属氧化 ─────→ 金属性质 ──→ 金属活动性
```

**依赖管理**：
```python
knowledge_graph = {
    "燃烧条件": {
        "prerequisites": ["氧气的性质", "化学变化"],
        "enables": ["灭火原理", "燃料的利用"],
        "difficulty_level": 2
    },
    "金属活动性": {
        "prerequisites": ["金属的性质", "置换反应"],
        "enables": ["金属冶炼", "金属防护"],
        "difficulty_level": 3
    }
}
```

### 5. 个性化适配

**基于学情调整**：
```python
class CurriculumAdapter:
    """课程规划适配器"""
    
    def adapt_to_class(self, base_plan, class_profile):
        """根据班级情况调整课程计划"""
        
        # 学习能力分析
        if class_profile.avg_ability < 0.6:
            # 降低难度，增加练习课
            base_plan.reduce_difficulty()
            base_plan.add_practice_lessons(ratio=0.3)
        
        # 学习风格偏好
        if class_profile.learning_style == "visual":
            # 增加实验课和视频
            base_plan.increase_experiments()
            base_plan.add_multimedia_content()
        
        # 进度调整
        if class_profile.is_behind_schedule():
            # 合并非核心主题
            base_plan.merge_optional_topics()
            base_plan.optimize_time_allocation()
        
        return base_plan
```

---

## 🏗️ 两层架构设计

### 架构图

```
┌───────────────────────────────────────────────────────────────┐
│                     应用层 (Application)                       │
│  教材管理后台 │ 课程规划看板 │ 教学日历 │ 进度跟踪             │
└───────────────────────────────────────────────────────────────┘
                              ↓
┌───────────────────────────────────────────────────────────────┐
│              上层元引擎 (Curriculum Meta-Engine)               │
│                                                               │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐       │
│  │ 教材分析器   │ → │ 主题划分器   │ → │ 课时规划器   │       │
│  └─────────────┘   └─────────────┘   └─────────────┘       │
│         ↓                 ↓                  ↓               │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐       │
│  │ 知识图谱     │   │ 故事设计器   │   │ 课型编排器   │       │
│  │ 构建器       │   │              │   │              │       │
│  └─────────────┘   └─────────────┘   └─────────────┘       │
│                                                               │
│  输出：CurriculumPlan (课程规划)                              │
└───────────────────────────────────────────────────────────────┘
                              ↓
              ┌───────────────┴───────────────┐
              ↓                               ↓
┌─────────────────────────┐     ┌─────────────────────────┐
│  LessonPlan #1          │     │  LessonPlan #72         │
│  (第1课时)               │ ... │  (第72课时)              │
└─────────────────────────┘     └─────────────────────────┘
              ↓                               ↓
┌───────────────────────────────────────────────────────────────┐
│              下层元引擎 (Lesson Meta-Engine)                   │
│                                                               │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐       │
│  │ 画像分析     │ → │ 教学设计     │ → │ 内容生成     │       │
│  └─────────────┘   └─────────────┘   └─────────────┘       │
│                                                               │
│  输出：LessonScript (课件脚本)                                │
└───────────────────────────────────────────────────────────────┘
```

### 数据流转

```python
# 1. 上层元工作流执行
curriculum_plan = await curriculum_engine.generate_plan(
    subject="化学",
    grade=9,
    semester="上册",
    publisher="人教版"
)

# curriculum_plan 包含：
# - units: 单元列表
# - lesson_sequence: 72个课时的元数据
# - story_arc: 故事主线和支线
# - knowledge_graph: 知识依赖图

# 2. 遍历课时序列，触发下层元工作流
for lesson_meta in curriculum_plan.lesson_sequence:
    lesson_script = await lesson_engine.generate_content(
        # 从上层继承的上下文
        unit=lesson_meta.unit,
        topic=lesson_meta.topic,
        lesson_type=lesson_meta.lesson_type,
        story_context=lesson_meta.story_segment,
        
        # 下层特定的参数
        duration=45,
        interaction_level="high",
        enable_experiment=lesson_meta.has_experiment
    )
    
    # 保存生成的课件
    await save_lesson(lesson_script)
```

---

## 📊 核心数据模型

### 1. Curriculum (教材/课程)

```python
@dataclass
class Curriculum:
    """教材级课程规划"""
    
    curriculum_id: str
    
    # 基本信息
    subject: str  # 学科
    grade: int  # 年级
    semester: str  # 上册/下册/全一册
    publisher: str  # 出版社
    
    # 课时规划
    total_periods: int  # 总课时
    weeks: int  # 教学周数
    periods_per_week: int  # 周课时
    
    # 单元划分
    units: List[CourseUnit]
    
    # 课时序列
    lesson_sequence: List[LessonMeta]
    
    # 故事设计
    story_universe: Optional[StoryUniverse] = None
    
    # 知识图谱
    knowledge_graph: Dict[str, KnowledgeNode]
    
    # 元数据
    created_at: datetime
    updated_at: datetime
    created_by: str
```

### 2. CourseUnit (单元/章节)

```python
@dataclass
class CourseUnit:
    """课程单元"""
    
    unit_id: str
    unit_number: int  # 第几单元
    name: str
    
    # 知识内容
    topics: List[str]  # 知识点列表
    objectives: List[str]  # 教学目标
    key_concepts: List[str]  # 核心概念
    
    # 课时分配
    total_periods: int
    lesson_types: Dict[str, int]  # {"NEW": 6, "PRACTICE": 2, "REVIEW": 1}
    
    # 难度与依赖
    difficulty_level: int  # 1-5
    prerequisites: List[str]  # 前置单元ID列表
    
    # 故事线
    story_chapter: Optional[StoryChapter] = None
```

### 3. LessonMeta (课时元数据)

```python
@dataclass
class LessonMeta:
    """单个课时的元数据（规划阶段）"""
    
    lesson_id: str
    sequence_number: int  # 第几课时（1-72）
    
    # 归属
    curriculum_id: str
    unit_id: str
    
    # 内容
    topic: str  # 本节课主题
    lesson_type: LessonType  # 课型
    
    # 课时安排
    duration: int  # 分钟
    scheduled_week: int  # 第几周
    scheduled_day: int  # 周几
    
    # 教学要素
    objectives: List[str]
    key_points: List[str]
    difficult_points: List[str]
    
    # 特殊要求
    has_experiment: bool = False
    experiment_type: Optional[str] = None
    
    needs_multimedia: bool = False
    interaction_level: str = "medium"
    
    # 故事信息
    story_segment: Optional[StorySegment] = None
    
    # 状态
    status: str = "planned"  # planned/generated/reviewed/published
    generated_lesson_id: Optional[str] = None  # 生成后的实际课件ID
```

### 4. StoryUniverse (故事宇宙)

```python
@dataclass
class StoryUniverse:
    """整本教材的故事世界观"""
    
    theme: str  # 故事主题
    setting: str  # 故事背景
    
    # 角色
    characters: List[StoryCharacter]
    
    # 主线剧情
    main_plot: StoryArc
    
    # 支线剧情（每个单元一个）
    sub_plots: Dict[str, StoryArc]  # unit_id -> StoryArc
    
    # 剧情节点（关键转折点）
    plot_points: List[PlotPoint]

@dataclass
class StoryCharacter:
    """故事角色"""
    name: str
    role: str
    personality: str
    abilities: List[str]
    growth_arc: str  # 角色成长线

@dataclass
class StoryArc:
    """故事线"""
    title: str
    acts: List[str]  # 三幕剧结构
    conflict: str  # 冲突
    resolution: str  # 解决

@dataclass
class StorySegment:
    """单节课的故事片段"""
    segment_id: str
    lesson_id: str
    
    # 情节
    plot_summary: str  # 剧情概要
    scene_description: str  # 场景描述
    
    # 承接
    previous_segment: Optional[str] = None
    next_segment: Optional[str] = None
    
    # 教学整合
    knowledge_integration: str  # 如何融入知识点
```

### 5. KnowledgeNode (知识节点)

```python
@dataclass
class KnowledgeNode:
    """知识图谱节点"""
    
    node_id: str
    name: str  # 知识点名称
    
    # 分类
    type: str  # concept/skill/principle/application
    domain: str  # 所属领域
    
    # 难度
    difficulty_level: int  # 1-5
    bloom_level: str  # remember/understand/apply/analyze/evaluate/create
    
    # 依赖关系
    prerequisites: List[str]  # 前置知识节点ID
    enables: List[str]  # 支撑的后续知识节点ID
    
    # 教学建议
    recommended_lesson_type: str
    estimated_time: int  # 分钟
    
    # 关联
    related_units: List[str]
    related_topics: List[str]
```

---

## 🔧 核心算法

### 1. 主题划分算法

```python
class TopicDivider:
    """主题划分器"""
    
    async def divide_topics(
        self, 
        textbook_content: str,
        curriculum_standard: Dict,
        total_periods: int
    ) -> List[CourseUnit]:
        """
        基于教材内容和课程标准，划分单元和主题
        """
        
        # 1. 提取教材目录结构
        toc = await self.extract_table_of_contents(textbook_content)
        
        # 2. 对照课程标准，识别核心知识点
        core_topics = await self.align_with_standard(toc, curriculum_standard)
        
        # 3. 知识点聚类，形成单元
        units = await self.cluster_topics(core_topics)
        
        # 4. 分配课时（基于难度和重要性）
        units = await self.allocate_periods(units, total_periods)
        
        # 5. 分析依赖关系
        units = await self.analyze_dependencies(units)
        
        return units
    
    async def cluster_topics(self, topics: List[str]) -> List[CourseUnit]:
        """
        使用LLM + 向量聚类进行主题分组
        """
        
        # 1. 获取知识点的向量表示
        embeddings = await self.get_embeddings(topics)
        
        # 2. K-means聚类
        clusters = self.kmeans_clustering(embeddings, n_clusters=8)
        
        # 3. 为每个cluster生成单元名称
        units = []
        for cluster_id, topic_indices in clusters.items():
            cluster_topics = [topics[i] for i in topic_indices]
            
            # 用LLM生成单元名称和目标
            unit_info = await self.llm.generate_unit_info(cluster_topics)
            
            units.append(CourseUnit(
                unit_id=f"unit_{cluster_id:02d}",
                unit_number=cluster_id,
                name=unit_info["name"],
                topics=cluster_topics,
                objectives=unit_info["objectives"],
                difficulty_level=self.estimate_difficulty(cluster_topics)
            ))
        
        return units
```

### 2. 课时序列生成算法

```python
class LessonScheduler:
    """课时编排器"""
    
    async def generate_lesson_sequence(
        self,
        units: List[CourseUnit],
        total_weeks: int,
        periods_per_week: int
    ) -> List[LessonMeta]:
        """
        生成整学期的课时序列
        """
        
        lesson_sequence = []
        current_week = 1
        current_period = 1
        
        for unit in units:
            # 为每个单元生成课时序列
            unit_lessons = await self._plan_unit_lessons(unit)
            
            for lesson in unit_lessons:
                # 计算上课时间
                week, day = self._calculate_schedule(
                    current_period, 
                    periods_per_week
                )
                
                lesson.sequence_number = current_period
                lesson.scheduled_week = week
                lesson.scheduled_day = day
                
                lesson_sequence.append(lesson)
                current_period += 1
        
        # 插入复习课和测验课
        lesson_sequence = self._insert_review_lessons(lesson_sequence)
        lesson_sequence = self._insert_assessments(lesson_sequence)
        
        return lesson_sequence
    
    async def _plan_unit_lessons(
        self, 
        unit: CourseUnit
    ) -> List[LessonMeta]:
        """
        为单个单元规划课时
        """
        
        lessons = []
        
        # 1. 新授课（讲解核心知识点）
        for topic in unit.topics:
            lessons.append(LessonMeta(
                lesson_id=self._gen_id(),
                unit_id=unit.unit_id,
                topic=topic,
                lesson_type=LessonType.NEW,
                duration=45,
                objectives=self._gen_objectives(topic),
                has_experiment=self._needs_experiment(topic)
            ))
        
        # 2. 练习课（每2-3个新授课后插入）
        for i in range(len(lessons) // 3):
            lessons.insert(
                (i + 1) * 3,
                LessonMeta(
                    lesson_id=self._gen_id(),
                    unit_id=unit.unit_id,
                    topic=f"{unit.name}·练习",
                    lesson_type=LessonType.PRACTICE,
                    duration=45
                )
            )
        
        # 3. 单元复习课
        lessons.append(LessonMeta(
            lesson_id=self._gen_id(),
            unit_id=unit.unit_id,
            topic=f"{unit.name}·复习",
            lesson_type=LessonType.REVIEW,
            duration=45
        ))
        
        return lessons
```

### 3. 故事线生成算法

```python
class StoryDesigner:
    """故事设计器"""
    
    async def design_story_universe(
        self,
        curriculum: Curriculum,
        units: List[CourseUnit]
    ) -> StoryUniverse:
        """
        为整本教材设计故事宇宙
        """
        
        # 1. 生成故事主题和背景
        theme_prompt = f"""
        为{curriculum.grade}年级{curriculum.subject}课程设计一个贯穿全学期的故事主题。
        
        单元列表：
        {[unit.name for unit in units]}
        
        要求：
        - 主题要吸引{curriculum.grade}年级学生
        - 要能自然融入所有单元的知识点
        - 设定要有足够的想象空间
        """
        
        universe_info = await self.llm.generate(theme_prompt)
        
        # 2. 设计主要角色
        characters = await self._design_characters(universe_info)
        
        # 3. 规划主线剧情（三幕剧结构）
        main_plot = await self._design_main_plot(units)
        
        # 4. 为每个单元设计支线剧情
        sub_plots = {}
        for unit in units:
            sub_plots[unit.unit_id] = await self._design_sub_plot(
                unit, 
                main_plot,
                characters
            )
        
        # 5. 生成关键剧情节点
        plot_points = await self._generate_plot_points(
            main_plot, 
            sub_plots
        )
        
        return StoryUniverse(
            theme=universe_info["theme"],
            setting=universe_info["setting"],
            characters=characters,
            main_plot=main_plot,
            sub_plots=sub_plots,
            plot_points=plot_points
        )
    
    async def generate_lesson_story(
        self,
        lesson_meta: LessonMeta,
        story_universe: StoryUniverse,
        previous_segment: Optional[StorySegment]
    ) -> StorySegment:
        """
        为单节课生成故事片段
        """
        
        unit_story = story_universe.sub_plots[lesson_meta.unit_id]
        
        prompt = f"""
        故事背景：{story_universe.setting}
        当前章节：{unit_story.title}
        本节课主题：{lesson_meta.topic}
        
        上一节课剧情：
        {previous_segment.plot_summary if previous_segment else "（本单元开始）"}
        
        请为本节课设计一个故事片段：
        1. 承接上一节课的剧情
        2. 自然引出本节课的知识点：{lesson_meta.key_points}
        3. 设置悬念，为下一节课铺垫
        4. 时长控制在5-8分钟
        """
        
        segment_info = await self.llm.generate(prompt)
        
        return StorySegment(
            segment_id=self._gen_id(),
            lesson_id=lesson_meta.lesson_id,
            plot_summary=segment_info["summary"],
            scene_description=segment_info["scene"],
            previous_segment=previous_segment.segment_id if previous_segment else None,
            knowledge_integration=segment_info["knowledge_integration"]
        )
```

### 4. 知识图谱构建算法

```python
class KnowledgeGraphBuilder:
    """知识图谱构建器"""
    
    async def build_knowledge_graph(
        self,
        units: List[CourseUnit]
    ) -> Dict[str, KnowledgeNode]:
        """
        构建教材级知识图谱
        """
        
        graph = {}
        
        # 1. 为每个知识点创建节点
        for unit in units:
            for topic in unit.topics:
                node = await self._create_knowledge_node(topic, unit)
                graph[node.node_id] = node
        
        # 2. 分析知识点之间的依赖关系
        for node_id, node in graph.items():
            # 使用LLM识别前置知识
            prerequisites = await self._identify_prerequisites(
                node, 
                graph
            )
            node.prerequisites = [p.node_id for p in prerequisites]
            
            # 反向建立"支撑"关系
            for prereq_id in node.prerequisites:
                graph[prereq_id].enables.append(node_id)
        
        # 3. 拓扑排序验证（检测循环依赖）
        if self._has_cycle(graph):
            # 使用LLM修复循环依赖
            graph = await self._fix_cycles(graph)
        
        # 4. 计算难度等级（基于依赖深度）
        for node in graph.values():
            node.difficulty_level = self._calculate_difficulty(node, graph)
        
        return graph
    
    def _calculate_difficulty(
        self, 
        node: KnowledgeNode,
        graph: Dict[str, KnowledgeNode]
    ) -> int:
        """
        基于依赖深度计算难度
        """
        
        if not node.prerequisites:
            return 1  # 基础知识
        
        # 递归计算最大依赖深度
        max_depth = 0
        for prereq_id in node.prerequisites:
            prereq_depth = self._calculate_difficulty(
                graph[prereq_id], 
                graph
            )
            max_depth = max(max_depth, prereq_depth)
        
        return min(max_depth + 1, 5)  # 难度上限为5
```

---

## 🔗 与单课级元工作流的集成

### 集成方式

```python
class IntegratedWorkflowEngine:
    """集成的两层工作流引擎"""
    
    def __init__(self):
        self.curriculum_engine = CurriculumMetaEngine()
        self.lesson_engine = LessonMetaEngine()
    
    async def generate_full_curriculum(
        self,
        subject: str,
        grade: int,
        semester: str,
        publisher: str,
        auto_generate_lessons: bool = False
    ) -> Curriculum:
        """
        生成完整课程（上层工作流）
        """
        
        # 1. 执行上层元工作流
        curriculum_plan = await self.curriculum_engine.generate_plan(
            subject=subject,
            grade=grade,
            semester=semester,
            publisher=publisher
        )
        
        # 2. 可选：自动生成所有课时内容
        if auto_generate_lessons:
            await self._generate_all_lessons(curriculum_plan)
        
        return curriculum_plan
    
    async def _generate_all_lessons(
        self, 
        curriculum: Curriculum
    ):
        """
        批量生成所有课时的内容（下层工作流）
        """
        
        previous_story = None
        
        for lesson_meta in curriculum.lesson_sequence:
            # 1. 生成本节课的故事片段
            if curriculum.story_universe:
                story_segment = await self.story_designer.generate_lesson_story(
                    lesson_meta,
                    curriculum.story_universe,
                    previous_story
                )
                lesson_meta.story_segment = story_segment
                previous_story = story_segment
            
            # 2. 执行下层元工作流，生成课件内容
            lesson_script = await self.lesson_engine.generate_content(
                # 从上层继承的参数
                subject=curriculum.subject,
                grade=curriculum.grade,
                topic=lesson_meta.topic,
                lesson_type=lesson_meta.lesson_type,
                
                # 故事上下文
                story_context=lesson_meta.story_segment.plot_summary if lesson_meta.story_segment else None,
                
                # 教学要求
                objectives=lesson_meta.objectives,
                key_points=lesson_meta.key_points,
                difficult_points=lesson_meta.difficult_points,
                
                # 特殊要求
                enable_experiment=lesson_meta.has_experiment,
                interaction_level=lesson_meta.interaction_level,
                
                # 其他参数
                duration=lesson_meta.duration
            )
            
            # 3. 保存课件
            lesson_id = await self.save_lesson(lesson_script)
            lesson_meta.generated_lesson_id = lesson_id
            lesson_meta.status = "generated"
            
            print(f"✅ 已生成第{lesson_meta.sequence_number}课时：{lesson_meta.topic}")
    
    async def generate_single_lesson(
        self,
        curriculum_id: str,
        lesson_number: int,
        **override_params
    ) -> LessonScript:
        """
        单独生成某一课时（支持覆盖参数）
        """
        
        # 1. 加载课程规划
        curriculum = await self.load_curriculum(curriculum_id)
        
        # 2. 获取课时元数据
        lesson_meta = curriculum.lesson_sequence[lesson_number - 1]
        
        # 3. 合并参数
        params = {
            "subject": curriculum.subject,
            "grade": curriculum.grade,
            "topic": lesson_meta.topic,
            "lesson_type": lesson_meta.lesson_type,
            **override_params  # 用户可覆盖
        }
        
        # 4. 执行下层元工作流
        lesson_script = await self.lesson_engine.generate_content(**params)
        
        return lesson_script
```

### API设计

```python
# 1. 创建课程规划
POST /v1/curriculums
{
    "subject": "化学",
    "grade": 9,
    "semester": "上册",
    "publisher": "人教版",
    "options": {
        "enable_story": true,
        "story_theme": "化学侦探社",
        "auto_generate_lessons": false
    }
}

# 返回
{
    "curriculum_id": "curr_chem_9_1",
    "total_units": 8,
    "total_lessons": 72,
    "story_universe": {...},
    "lesson_sequence": [...]
}

# 2. 批量生成课时内容
POST /v1/curriculums/{curriculum_id}/generate_all_lessons
{
    "start_from": 1,
    "end_at": 72,
    "parallel": true,  # 并行生成
    "max_workers": 10
}

# 3. 单独生成某一课时
POST /v1/curriculums/{curriculum_id}/lessons/{lesson_number}/generate
{
    "override_params": {
        "interaction_level": "high",
        "enable_experiment": true
    }
}

# 4. 获取课程进度
GET /v1/curriculums/{curriculum_id}/progress
{
    "total_lessons": 72,
    "generated": 15,
    "in_progress": 3,
    "pending": 54,
    "progress_rate": 0.21
}
```

---

## 📈 优势与价值

### 1. 系统性规划
- ✅ 从宏观到微观，完整覆盖一学期教学
- ✅ 知识点有序组织，符合认知规律
- ✅ 课时分配科学，节奏张弛有度

### 2. 故事连贯性
- ✅ 统一的世界观，增强沉浸感
- ✅ 剧情承接自然，激发持续兴趣
- ✅ 角色成长与知识进阶同步

### 3. 个性化适配
- ✅ 根据班级情况调整难度和节奏
- ✅ 支持不同学习风格的教学策略
- ✅ 灵活应对教学进度变化

### 4. 工作流复用
- ✅ 上层规划一次，下层批量生成
- ✅ 单课级工作流完全复用现有能力
- ✅ 模板化管理，快速迭代

---

## 🎯 实施建议

### Phase 1: 上层元工作流开发 (2-3周)
1. 实现主题划分算法
2. 实现课时编排算法
3. 实现知识图谱构建
4. 开发故事设计器

### Phase 2: 集成与测试 (1-2周)
1. 两层工作流集成
2. 端到端测试（生成完整一学期课程）
3. 性能优化（并行生成）

### Phase 3: 功能增强 (2-3周)
1. 个性化适配算法
2. 进度跟踪与调整
3. 教学日历可视化
4. 协同编辑功能

---

**下一步：是否开始更新相关文档？**
