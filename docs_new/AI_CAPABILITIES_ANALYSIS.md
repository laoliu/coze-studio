# AI能力层完整分析与规划

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 能力盘点与补全规划

---

## 📋 目录

1. [当前实现状态](#当前实现状态)
2. [规划中的AI能力](#规划中的ai能力)
3. [遗漏能力分析](#遗漏能力分析)
4. [能力补全计划](#能力补全计划)
5. [实现优先级](#实现优先级)

---

## 当前实现状态

### 已实现的AI能力（Cycle 1）

| 能力模块 | 文件位置 | 实现状态 | 功能完整度 |
|---------|---------|---------|-----------|
| **LLM客户端** | `src/ai/llm_client.py` | ✅ 完成 | 90% |
| **提示词管理** | `src/ai/prompt_manager.py` | ✅ 完成 | 85% |
| **学习目标生成** | `src/ai/objective_generator.py` | 🔶 占位符 | 30% |

### 详细实现分析

#### 1. LLM客户端 (LLMClient) ✅

**文件**: `src/ai/llm_client.py` (141行)

**已实现功能**:
```python
class LLMClient:
    ✅ __init__()           # 初始化（支持多模型配置）
    ✅ chat()              # 基础聊天接口
    ✅ generate_json()     # JSON格式响应生成
```

**技术栈**:
- ✅ 使用 LiteLLM（统一多LLM接口）
- ✅ 支持 OpenAI、Claude、Gemini等
- ✅ 异步调用（async/await）
- ✅ 错误处理和日志记录

**优点**:
- ✅ 接口统一，易于扩展
- ✅ 支持多种LLM provider
- ✅ JSON格式自动清理（去除markdown标记）

**待优化**:
- ⚠️ 缺少重试机制（API失败时）
- ⚠️ 缺少速率限制（Rate Limiting）
- ⚠️ 缺少Token使用统计
- ⚠️ 缺少流式输出支持（Streaming）
- ⚠️ 缺少成本追踪

---

#### 2. 提示词管理器 (PromptManager) ✅

**文件**: `src/ai/prompt_manager.py` (200行)

**已实现功能**:
```python
class PromptManager:
    ✅ load_prompt_config()    # 加载YAML配置
    ✅ render_prompt()         # Jinja2模板渲染
    ✅ get_system_prompt()     # 获取system prompt
    ✅ get_user_prompt()       # 获取user prompt
    ✅ build_messages()        # 构建完整消息列表
```

**技术栈**:
- ✅ YAML配置文件
- ✅ Jinja2模板引擎
- ✅ 模板缓存机制

**优点**:
- ✅ 配置与代码分离
- ✅ 模板可复用
- ✅ 支持变量渲染

**待优化**:
- ⚠️ 缺少Prompt版本管理
- ⚠️ 缺少A/B测试支持
- ⚠️ 缺少Prompt效果评估
- ⚠️ 缺少Few-shot示例管理
- ⚠️ 缺少Prompt优化建议

---

#### 3. 学习目标生成器 (ObjectiveGenerator) 🔶

**文件**: `src/ai/objective_generator.py` (111行)

**已实现功能**:
```python
class ObjectiveGenerator:
    ✅ __init__()                           # 初始化
    🔶 generate_objectives()                # 生成学习目标（当前为占位符）
    ✅ _generate_placeholder_objectives()   # 占位符实现
```

**当前状态**:
```python
# TODO: Cycle 2 实现完整的AI生成逻辑
# 当前返回占位符数据

# 占位符实现（Cycle 1）
objectives = self._generate_placeholder_objectives(context)
```

**待实现**:
- ❌ Prompt构建逻辑
- ❌ LLM调用逻辑
- ❌ 响应解析和验证
- ❌ SMART原则检查
- ❌ 布鲁姆分类法应用

---

## 规划中的AI能力

### V2.0架构中定义的AI能力层

根据 `docs/architecture.md` 和 `docs/ADAPTER_VS_CAPABILITIES.md`：

```
┌────────────────────────────────────────────────────────┐
│               AI能力层 (AI Capability Layer)            │
├────────────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐           │
│  │ 意图识别         │  │ 学习目标生成      │           │
│  │ IntentRecognizer │  │ ObjectiveGen     │           │
│  └──────────────────┘  └──────────────────┘           │
│                                                         │
│  ┌──────────────────┐  ┌──────────────────┐           │
│  │ 内容发现         │  │ 故事化叙述       │           │
│  │ ContentDiscovery │  │ NarrativeGen    │           │
│  └──────────────────┘  └──────────────────┘           │
│                                                         │
│  ┌──────────────────┐  ┌──────────────────┐           │
│  │ 质量评估         │  │ 提示词管理       │           │
│  │ QualityAssessor  │  │ PromptManager   │           │
│  └──────────────────┘  └──────────────────┘           │
└────────────────────────────────────────────────────────┘
```

### 规划的6大核心能力

| # | 能力模块 | 状态 | 优先级 | 文件位置（规划） |
|---|---------|------|-------|----------------|
| 1 | **意图识别** | ❌ 未实现 | P0 | `src/ai/intent_recognizer.py` |
| 2 | **学习目标生成** | 🔶 占位符 | P0 | `src/ai/objective_generator.py` |
| 3 | **内容发现** | ❌ 未实现 | P0 | `src/ai/content_discovery.py` |
| 4 | **故事化叙述** | ❌ 未实现 | P1 | `src/ai/narrative_generator.py` |
| 5 | **质量评估** | ❌ 未实现 | P1 | `src/ai/quality_assessor.py` |
| 6 | **提示词管理** | ✅ 完成 | P0 | `src/ai/prompt_manager.py` |

---

## 遗漏能力分析

### 核心遗漏（P0 - 必须实现）

#### 1. ❌ 意图识别器 (Intent Recognizer)

**缺失原因**: Cycle 1阶段未规划

**功能定义**:
```python
class IntentRecognizer:
    """用户意图识别器"""
    
    async def recognize(self, raw_input: str) -> Intent:
        """
        识别用户输入的意图
        
        输入示例:
        - "讲解初三化学燃烧条件"
        - "分析莫奈的《日出·印象》"
        
        输出:
        Intent(
            intent_type="content_generation",
            task="explain",  # explain, analyze, create, assess
            entities={
                "topic": "燃烧条件",
                "subject": "化学",
                "grade": 9
            },
            confidence=0.95
        )
        """
```

**应用场景**:
1. **领域路由辅助**: 帮助领域路由器选择适配器
2. **任务类型识别**: 区分"讲解"、"分析"、"创作"等任务
3. **实体抽取**: 提取主题、学科、年级等关键信息
4. **意图澄清**: 当输入模糊时，生成澄清问题

**技术方案**:
```python
class IntentRecognizer:
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
        self.intent_prompts = {
            "classify": "判断用户意图类型",
            "extract_entities": "提取关键实体",
            "clarify": "生成澄清问题"
        }
    
    async def recognize(self, raw_input: str) -> Intent:
        # 1. 意图分类
        intent_type = await self._classify_intent(raw_input)
        
        # 2. 实体抽取
        entities = await self._extract_entities(raw_input)
        
        # 3. 计算置信度
        confidence = self._calculate_confidence(intent_type, entities)
        
        # 4. 如果置信度低，生成澄清问题
        if confidence < 0.7:
            clarification = await self._generate_clarification(raw_input)
            return Intent(
                intent_type=intent_type,
                entities=entities,
                confidence=confidence,
                needs_clarification=True,
                clarification_question=clarification
            )
        
        return Intent(
            intent_type=intent_type,
            entities=entities,
            confidence=confidence
        )
```

**估计工作量**: 2-3天

---

#### 2. 🔶 学习目标生成器 (Objective Generator) - 完善

**缺失原因**: Cycle 1仅实现占位符

**需要补充的功能**:

```python
class ObjectiveGenerator:
    """学习目标生成器（完整版）"""
    
    async def generate_objectives(
        self,
        context: WorkflowContext
    ) -> List[LearningObjective]:
        """
        完整的AI生成流程
        
        步骤:
        1. 构建few-shot prompt
        2. 调用LLM生成目标
        3. 解析JSON响应
        4. SMART原则验证
        5. 布鲁姆分类法检查
        6. 时间分配合理性验证
        """
        # 1. 构建prompt
        messages = self._build_prompt_with_examples(context)
        
        # 2. 调用LLM
        response = await self.llm.generate_json(messages)
        
        # 3. 解析响应
        objectives = self._parse_objectives(response)
        
        # 4. 验证和修正
        objectives = self._validate_objectives(objectives, context)
        
        # 5. 时间分配
        objectives = self._allocate_time(objectives, context.duration)
        
        return objectives
    
    def _build_prompt_with_examples(self, context):
        """构建包含few-shot示例的prompt"""
        # 根据领域选择相关示例
        examples = self.example_db.get_examples_by_domain(context.domain)
        
        prompt = f"""
你是一个专业的教学设计专家。请根据以下信息生成学习目标：

领域: {context.domain}
主题: {context.subject}
受众: {context.audience}
时长: {context.duration}分钟

要求:
1. 遵循SMART原则（具体、可衡量、可达成、相关、有时限）
2. 遵循布鲁姆分类法（记忆、理解、应用、分析、评价、创造）
3. 难度适合受众认知水平
4. 时间分配合理

示例：
{self._format_examples(examples)}

请生成3-5个学习目标，JSON格式返回。
"""
        return prompt
    
    def _validate_objectives(self, objectives, context):
        """验证目标合理性"""
        validated = []
        
        for obj in objectives:
            # SMART检查
            if not self._is_smart(obj):
                obj = self._make_smart(obj)
            
            # 布鲁姆层级检查
            if not self._is_bloom_appropriate(obj, context.audience):
                obj = self._adjust_bloom_level(obj)
            
            # 时间合理性检查
            if obj.estimated_time > context.duration * 0.5:
                obj.estimated_time = context.duration * 0.4
            
            validated.append(obj)
        
        return validated
```

**估计工作量**: 2天

---

#### 3. ❌ 内容发现引擎 (Content Discovery Engine)

**缺失原因**: Cycle 1未规划，属于V2.0核心能力

**功能定义**:
```python
class ContentDiscoveryEngine:
    """AI驱动的多源内容发现引擎"""
    
    def __init__(self):
        # 内容源管理
        self.sources = {
            "wikipedia": WikipediaRetriever(),
            "wikiart": WikiArtRetriever(),
            "youtube": YouTubeRetriever(),
            "douban": DoubanRetriever(),
            # ... 更多源
        }
        
        # 检索引擎
        self.bm25_retriever = BM25Retriever()
        self.semantic_retriever = SemanticRetriever()
        
        # 质量评估器
        self.quality_assessor = ContentQualityAssessor()
        
        # 结果融合器
        self.rrf_fusion = RRFFusion()
    
    async def discover(
        self,
        objectives: List[LearningObjective],
        sources: List[str],
        filters: Dict[str, Any]
    ) -> StructuredContent:
        """
        发现和整合多源内容
        
        流程:
        1. 查询生成 - 将学习目标转换为搜索查询
        2. 多源并发检索 - 从多个源并发搜索
        3. 质量评估 - 对每个结果评分
        4. 结果融合 - 使用RRF算法融合
        5. 结构化整合 - 组织成结构化内容
        """
```

**子模块设计**:

##### 3.1 查询生成器 (Query Generator)
```python
class QueryGenerator:
    """将学习目标转换为搜索查询"""
    
    async def generate_queries(
        self,
        objectives: List[LearningObjective]
    ) -> List[SearchQuery]:
        """
        为每个学习目标生成多个查询变体
        
        输入:
        - "理解印象派的艺术特点"
        
        输出:
        - "印象派 艺术特点"
        - "印象派绘画风格"
        - "Impressionism characteristics"
        - "莫奈 雷诺阿 印象派"
        """
        queries = []
        
        for obj in objectives:
            # 关键词提取
            keywords = self._extract_keywords(obj.title)
            
            # 生成同义词变体
            synonyms = await self._generate_synonyms(keywords)
            
            # 生成相关实体查询
            entities = await self._extract_entities(obj.title)
            
            # 组合查询
            queries.extend(self._combine_queries(keywords, synonyms, entities))
        
        return queries
```

##### 3.2 内容源检索器 (Content Retriever)
```python
class WikipediaRetriever:
    """Wikipedia内容检索"""
    
    async def search(
        self,
        query: SearchQuery,
        limit: int = 10
    ) -> List[ContentItem]:
        """
        从Wikipedia检索内容
        
        使用API:
        - Wikipedia API (REST)
        - 支持中英文
        """
        # 调用Wikipedia API
        results = await self.api_client.search(query.text, limit=limit)
        
        # 解析结果
        items = []
        for result in results:
            item = ContentItem(
                source="wikipedia",
                title=result.title,
                url=result.url,
                content=result.extract,
                metadata={
                    "page_id": result.pageid,
                    "language": result.lang
                }
            )
            items.append(item)
        
        return items
```

##### 3.3 质量评估器 (Quality Assessor)
```python
class ContentQualityAssessor:
    """内容质量评估"""
    
    async def assess(
        self,
        content_item: ContentItem,
        objective: LearningObjective
    ) -> QualityScore:
        """
        评估内容质量
        
        维度:
        1. 相关性 (Relevance) - 与学习目标的匹配度
        2. 权威性 (Authority) - 来源的可信度
        3. 准确性 (Accuracy) - 内容的准确性
        4. 完整性 (Completeness) - 信息的完整性
        5. 时效性 (Timeliness) - 信息的新鲜度
        """
        scores = {}
        
        # 1. 相关性评分（语义相似度）
        scores['relevance'] = await self._calculate_relevance(
            content_item.content,
            objective.title
        )
        
        # 2. 权威性评分（基于来源）
        scores['authority'] = self._get_source_authority(content_item.source)
        
        # 3. 准确性评分（事实检查）
        scores['accuracy'] = await self._check_accuracy(content_item.content)
        
        # 4. 完整性评分
        scores['completeness'] = self._check_completeness(
            content_item.content,
            objective
        )
        
        # 5. 时效性评分
        scores['timeliness'] = self._check_timeliness(content_item.metadata)
        
        # 加权平均
        overall_score = (
            scores['relevance'] * 0.4 +
            scores['authority'] * 0.2 +
            scores['accuracy'] * 0.2 +
            scores['completeness'] * 0.1 +
            scores['timeliness'] * 0.1
        )
        
        return QualityScore(
            overall=overall_score,
            dimensions=scores
        )
```

##### 3.4 RRF融合器 (RRF Fusion)
```python
class RRFFusion:
    """Reciprocal Rank Fusion - 结果融合算法"""
    
    def fuse(
        self,
        ranked_lists: Dict[str, List[ScoredContent]],
        k: int = 60
    ) -> List[ScoredContent]:
        """
        融合多个排序列表
        
        RRF公式: score(d) = Σ 1 / (k + rank(d))
        
        输入:
        {
            "bm25": [item1, item2, item3],
            "semantic": [item3, item1, item5],
            "authority": [item1, item4, item2]
        }
        
        输出:
        融合后的排序列表
        """
        scores = defaultdict(float)
        
        for source, ranked_list in ranked_lists.items():
            for rank, item in enumerate(ranked_list, start=1):
                scores[item.id] += 1.0 / (k + rank)
        
        # 按融合分数排序
        fused = sorted(
            scores.items(),
            key=lambda x: x[1],
            reverse=True
        )
        
        return fused
```

**估计工作量**: 5-7天

---

### 重要遗漏（P1 - 应该实现）

#### 4. ❌ 故事化叙述引擎 (Narrative Generator)

**缺失原因**: Cycle 1未规划，属于V2.0增值能力

**功能定义**:
```python
class NarrativeGenerator:
    """故事化叙述生成引擎"""
    
    def __init__(self):
        self.strategy_selector = NarrativeStrategySelector()
        self.character_generator = CharacterGenerator()
        self.plot_generator = PlotGenerator()
        self.scene_generator = SceneGenerator()
    
    async def generate_narrative(
        self,
        content: StructuredContent,
        domain: DomainType,
        audience: AudienceProfile
    ) -> NarrativeContent:
        """
        将枯燥内容转化为故事
        
        流程:
        1. 策略选择 - 选择叙述策略（人物驱动、情节驱动等）
        2. 人物设计 - 设计故事主角
        3. 情节构建 - 三幕式结构
        4. 场景设计 - 具体场景描写
        5. LLM生成 - 生成故事文本
        6. 质量评估 - 评估叙述质量
        """
```

**6种故事化策略**:

| 策略 | 适用场景 | 示例 |
|------|---------|------|
| **人物驱动** | 历史人物、科学家 | "莫奈站在勒阿弗尔港口，面对清晨的雾气..." |
| **情节驱动** | 事件、发现 | "1872年的一个清晨，艺术史将被改写..." |
| **情境嵌入** | 抽象概念 | "想象你正在观看一场魔术表演（燃烧反应）..." |
| **问题解决** | 技能学习 | "小明遇到了一个难题：如何用Excel分析销售数据..." |
| **对比冲突** | 对立概念 | "传统绘画 vs 印象派：一场艺术革命..." |
| **时间旅行** | 历史事件 | "让我们回到1872年，见证印象派的诞生..." |

**子模块设计**:

##### 4.1 叙述策略选择器
```python
class NarrativeStrategySelector:
    """选择合适的叙述策略"""
    
    def select(
        self,
        domain: DomainType,
        audience: AudienceProfile,
        content_type: str
    ) -> NarrativeStrategy:
        """
        根据领域、受众、内容类型选择策略
        
        规则:
        - 美术史 + 成人 → 人物驱动或情节驱动
        - K12 + 儿童 → 情境嵌入或问题解决
        - 抽象概念 → 类比故事
        """
        if domain == DomainType.ART_HISTORY:
            return NarrativeStrategy.CHARACTER_DRIVEN
        elif domain == DomainType.K12_EDUCATION and audience.age < 12:
            return NarrativeStrategy.SCENARIO_EMBEDDED
        else:
            return NarrativeStrategy.PLOT_DRIVEN
```

##### 4.2 三幕式结构生成器
```python
class ThreeActStructure:
    """三幕式故事结构"""
    
    async def generate(
        self,
        content: StructuredContent
    ) -> ThreeActStory:
        """
        生成三幕式结构
        
        第一幕 - 铺垫 (Setup):
        - 引入人物/场景
        - 提出问题/目标
        - 建立情境
        
        第二幕 - 冲突 (Confrontation):
        - 深入探索
        - 遇到挑战
        - 分析解决
        
        第三幕 - 解决 (Resolution):
        - 达成目标
        - 总结收获
        - 留下思考
        """
        act1 = await self._generate_setup(content)
        act2 = await self._generate_confrontation(content)
        act3 = await self._generate_resolution(content)
        
        return ThreeActStory(act1=act1, act2=act2, act3=act3)
```

**估计工作量**: 4-5天

---

#### 5. ❌ 质量评估器 (Quality Assessor)

**缺失原因**: Cycle 1未规划

**功能定义**:
```python
class QualityAssessor:
    """内容质量评估器"""
    
    async def assess(
        self,
        generated_content: Any,
        context: WorkflowContext,
        assessment_type: str
    ) -> QualityReport:
        """
        评估生成内容的质量
        
        评估维度:
        1. 准确性 (Accuracy) - 事实正确性
        2. 完整性 (Completeness) - 是否覆盖所有学习目标
        3. 适配性 (Appropriateness) - 是否适合目标受众
        4. 连贯性 (Coherence) - 逻辑是否清晰
        5. 创新性 (Creativity) - 是否有新颖性
        """
```

**子模块设计**:

##### 5.1 准确性检查器
```python
class AccuracyChecker:
    """事实准确性检查"""
    
    async def check(self, content: str) -> AccuracyScore:
        """
        检查内容准确性
        
        方法:
        1. 事实抽取 - 提取可验证的陈述
        2. 外部验证 - 与权威源对比
        3. 一致性检查 - 检查内部矛盾
        """
        # 提取事实陈述
        facts = await self._extract_facts(content)
        
        # 外部验证
        verification_results = []
        for fact in facts:
            result = await self._verify_fact(fact)
            verification_results.append(result)
        
        # 计算准确性分数
        accuracy = sum(r.is_correct for r in verification_results) / len(facts)
        
        return AccuracyScore(
            score=accuracy,
            verified_facts=len(facts),
            incorrect_facts=[r for r in verification_results if not r.is_correct]
        )
```

##### 5.2 完整性检查器
```python
class CompletenessChecker:
    """内容完整性检查"""
    
    def check(
        self,
        content: GeneratedContent,
        objectives: List[LearningObjective]
    ) -> CompletenessScore:
        """
        检查是否覆盖所有学习目标
        
        方法:
        1. 目标映射 - 将内容映射到目标
        2. 覆盖度计算 - 计算每个目标的覆盖度
        3. 缺失检测 - 找出未覆盖的目标
        """
        coverage = {}
        
        for obj in objectives:
            # 检查内容中是否包含该目标的关键概念
            coverage[obj.id] = self._calculate_coverage(content, obj)
        
        overall_completeness = sum(coverage.values()) / len(objectives)
        
        missing_objectives = [
            obj for obj, cov in coverage.items() if cov < 0.5
        ]
        
        return CompletenessScore(
            score=overall_completeness,
            objective_coverage=coverage,
            missing_objectives=missing_objectives
        )
```

**估计工作量**: 3-4天

---

### 次要遗漏（P2 - 可选实现）

#### 6. ❌ 多模态生成器 (Multimodal Generator)

**功能**: 生成图片、音频、视频等多模态内容

**示例**:
```python
class MultimodalGenerator:
    """多模态内容生成"""
    
    async def generate_image(
        self,
        description: str,
        style: str = "realistic"
    ) -> ImageContent:
        """使用DALL-E或Stable Diffusion生成图片"""
    
    async def generate_audio(
        self,
        text: str,
        voice: str = "default"
    ) -> AudioContent:
        """使用TTS生成音频"""
    
    async def generate_video(
        self,
        script: str,
        images: List[ImageContent]
    ) -> VideoContent:
        """组合图片和音频生成视频"""
```

**估计工作量**: 5-7天

---

#### 7. ❌ 个性化推荐器 (Personalization Engine)

**功能**: 根据用户历史和偏好推荐内容

**示例**:
```python
class PersonalizationEngine:
    """个性化推荐引擎"""
    
    async def recommend(
        self,
        user_profile: UserProfile,
        available_content: List[Content]
    ) -> List[Content]:
        """
        推荐个性化内容
        
        考虑因素:
        1. 用户历史学习记录
        2. 用户偏好设置
        3. 学习风格
        4. 难度适配
        """
```

**估计工作量**: 4-5天

---

#### 8. ❌ 交互式对话器 (Interactive Dialog)

**功能**: 支持多轮对话，动态调整内容

**示例**:
```python
class InteractiveDialog:
    """交互式对话管理"""
    
    async def handle_turn(
        self,
        user_input: str,
        conversation_history: List[Message]
    ) -> Response:
        """
        处理对话轮次
        
        功能:
        1. 意图识别
        2. 上下文管理
        3. 响应生成
        4. 状态更新
        """
```

**估计工作量**: 3-4天

---

## 能力补全计划

### Phase 1: 核心能力补全（P0） - 2周

| 任务 | 工作量 | 负责人 | 依赖 |
|------|--------|--------|------|
| 1. 意图识别器开发 | 3天 | AI组 | LLMClient |
| 2. 学习目标生成器完善 | 2天 | AI组 | PromptManager |
| 3. 内容发现引擎基础版 | 7天 | AI组 | - |

**交付物**:
- `src/ai/intent_recognizer.py`
- `src/ai/objective_generator.py` (完整版)
- `src/ai/content_discovery.py`
- 单元测试和集成测试

---

### Phase 2: 增值能力开发（P1） - 2周

| 任务 | 工作量 | 负责人 | 依赖 |
|------|--------|--------|------|
| 4. 故事化叙述引擎 | 5天 | AI组 | ContentDiscovery |
| 5. 质量评估器 | 4天 | AI组 | - |
| 6. LLM客户端优化 | 1天 | AI组 | - |

**交付物**:
- `src/ai/narrative_generator.py`
- `src/ai/quality_assessor.py`
- 优化后的 `src/ai/llm_client.py`

---

### Phase 3: 扩展能力开发（P2） - 可选

| 任务 | 工作量 | 负责人 | 依赖 |
|------|--------|--------|------|
| 7. 多模态生成器 | 7天 | AI组 | 外部API |
| 8. 个性化推荐器 | 5天 | AI组 | 用户数据 |
| 9. 交互式对话器 | 4天 | AI组 | IntentRecognizer |

---

## 实现优先级

### 优先级矩阵

```
高影响  │  P0: 意图识别    │  P0: 内容发现
       │  P0: 目标生成    │  P1: 故事化叙述
       │                 │  P1: 质量评估
───────┼─────────────────┼──────────────────
       │                 │
低影响  │  P2: 个性化     │  P2: 多模态
       │  P2: 交互对话    │
       │                 │
       └─────────────────┴──────────────────
         低复杂度          高复杂度
```

### 推荐实施路径

**第一阶段（2周）- Cycle 2**:
1. ✅ 完善学习目标生成器（2天）
2. ✅ 开发意图识别器（3天）
3. ✅ 开发内容发现引擎基础版（7天）

**第二阶段（2周）- Cycle 3**:
4. ✅ 开发故事化叙述引擎（5天）
5. ✅ 开发质量评估器（4天）
6. ✅ 优化LLM客户端（1天）

**第三阶段（可选）- Cycle 4**:
7. 开发多模态生成器
8. 开发个性化推荐器
9. 开发交互式对话器

---

## 总结

### 当前AI能力完成度

| 类别 | 已完成 | 占位符 | 未实现 | 总计 |
|------|--------|--------|--------|------|
| **核心能力(P0)** | 1 | 1 | 2 | 4 |
| **增值能力(P1)** | 1 | 0 | 2 | 3 |
| **扩展能力(P2)** | 0 | 0 | 3 | 3 |
| **总计** | 2 | 1 | 7 | 10 |

**完成度**: 30% （3/10）

### 关键发现

✅ **已做好的**:
1. LLM客户端架构良好，支持多provider
2. 提示词管理器设计合理，支持模板化
3. 基础架构清晰，易于扩展

⚠️ **需要补充的**:
1. **意图识别器** - 核心缺失，影响领域路由
2. **内容发现引擎** - V2.0核心能力，必须实现
3. **学习目标生成器** - 需要从占位符升级到完整实现

🔮 **建议优化的**:
1. LLM客户端增加重试、限流、统计功能
2. 提示词管理器增加版本管理和A/B测试
3. 增加完善的错误处理和日志记录

### 下一步行动

**立即行动**（本周）:
1. 制定详细的AI能力开发计划
2. 设计意图识别器的API接口
3. 调研内容发现引擎的技术方案

**近期规划**（2周内）:
1. 完成P0级别的3个核心能力
2. 编写完整的单元测试
3. 集成到现有工作流引擎

**中期目标**（1个月内）:
1. 完成P1级别的2个增值能力
2. 进行端到端测试
3. 性能优化和调优

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow AI Team
