# 通用领域内容生成·元工作流平台 - 架构设计文档

**文档版本**: V2.0 ⭐ 重大升级  
**日期**: 2025-12-09  
**作者**: Meta-Workflow架构团队  
**状态**: 设计阶段

> **V2.0变化**：从K12垂直工具架构升级为多领域通用平台架构  
> **核心升级**：新增AI目标生成器、AI内容发现引擎、故事化叙述引擎、领域适配器架构

---

## 目录

1. [架构概述](#1-架构概述)
2. [总体架构设计](#2-总体架构设计)
3. [核心模块设计](#3-核心模块设计)
4. [技术选型](#4-技术选型)
5. [数据流设计](#5-数据流设计)
6. [部署架构](#6-部署架构)
7. [安全架构](#7-安全架构)

---

## 1. 架构概述

### 1.0 V2.0架构理念：孵化器-鸡-蛋模型 ⭐ 新增

**架构隐喻**：

```
┌─────────────────────────────────────────────────────────┐
│                    孵化器（Incubator）                    │
│              元工作流引擎 - 通用基础设施                   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  AI目标生成 + AI内容发现 + 故事化叙述 + 工作流编排 │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                         │
                         │ 孵化/支持
                         ▼
┌─────────────────────────────────────────────────────────┐
│                      鸡（Chickens）                       │
│                   领域适配器插件                          │
│                                                          │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐     │
│  │ K12  │  │美术史│  │畅销书│  │职业  │  │企业  │ ... │
│  │教育  │  │教育  │  │解读  │  │培训  │  │培训  │     │
│  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘     │
└─────────────────────────────────────────────────────────┘
                         │
                         │ 生产/输出
                         ▼
┌─────────────────────────────────────────────────────────┐
│                      蛋（Eggs）                           │
│                   内容产物实例                            │
│                                                          │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐     │
│  │数学  │  │莫奈  │  │原则  │  │Excel │  │领导力│ ... │
│  │教案  │  │讲解稿│  │读书会│  │培训  │  │工作坊│     │
│  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘     │
└─────────────────────────────────────────────────────────┘
```

**核心设计思想**：

1. **孵化器层（通用能力）**：
   - 与具体领域无关的通用AI能力
   - 提供标准化的孵化环境（API、基础设施、工具链）
   - 一次开发，多领域复用

2. **鸡层（领域专业化）**：
   - 每个领域有独特的"基因"（质量标准、内容源、生成策略）
   - 插件化设计，可独立开发和部署
   - 领域越多，平台价值越大（网络效应）

3. **蛋层（内容产物）**：
   - 最终交付给用户的内容
   - 每个"蛋"都带有其"母鸡"的特征
   - 标准化输出格式，易于管理和分发

**架构优势**：

| 传统架构 | V2.0孵化器架构 |
|---------|---------------|
| 每个领域独立开发系统 | 领域共享通用孵化器 |
| 重复造轮子，成本高 | 边际成本趋近于0 |
| 领域间无协同效应 | 领域越多，价值越大 |
| 扩展周期3-6个月 | 扩展周期1-2周 |
| 难以统一质量标准 | 孵化器保证基础质量 |

---

### 1.1 架构原则 ⭐ V2.0扩展

本系统采用**分层架构** + **事件驱动** + **插件化**的设计理念，遵循以下核心原则：

| 原则 | 说明 | 实现方式 |
|------|------|---------|
| **高内聚低耦合** | 模块职责单一，接口清晰 | 领域驱动设计(DDD) |
| **可扩展性优先** ⭐ | 支持水平扩展和功能扩展 | **插件化架构 + 领域适配器** |
| **配置即代码** | 工作流定义使用声明式配置 | YAML/JSON DSL |
| **故障隔离** | 单点故障不影响整体 | 熔断机制 + 降级策略 |
| **可观测性** | 全链路监控和追踪 | 分布式链路追踪 + 指标采集 |
| **领域自治** ⭐ 新增 | 每个领域独立管理质量标准 | 领域适配器接口 |
| **AI驱动** ⭐ 新增 | 从被动配置到主动生成 | AI目标生成 + 内容发现 + 故事化 |
| **多源整合** ⭐ 新增 | 从单一知识库到全网检索 | 混合检索 + 质量评估 |

### 1.2 架构风格 ⭐ V2.0扩展

**分层架构（Layered Architecture）**：
- 应用层（Presentation Layer）
- **领域适配器层（Domain Adapter Layer）** ⭐ 新增
- 业务逻辑层（Business Logic Layer）
  - **AI能力层（AI Capability Layer）** ⭐ 新增
  - 工作流编排层（Orchestration Layer）
- 数据访问层（Data Access Layer）
  - **多源内容层（Multi-Source Content Layer）** ⭐ 新增
- 基础设施层（Infrastructure Layer）

**事件驱动架构（Event-Driven Architecture）**：
- 工作流节点间通过事件通信
- 支持异步处理和解耦
- 便于扩展和监控

**插件化架构（Plugin Architecture）** ⭐ 新增：
- 领域适配器作为插件动态加载
- 支持第三方开发者贡献适配器
- 热插拔，无需重启系统

### 1.3 关键设计决策 ⭐ V2.0扩展

| 决策点 | V1.0选择 | V2.0选择 ⭐ | V2.0理由 |
|--------|---------|-----------|---------|
| **工作流引擎** | 自研 + LangGraph | 同左 | 需要高度定制化的编排能力 |
| **配置存储** | PostgreSQL + Redis | 同左 | 结构化数据 + 高速缓存 |
| **向量数据库** | Milvus / Qdrant | 同左 + **Faiss**（本地） | 支持大规模向量检索 + 本地部署选项 |
| **消息队列** | Redis Streams | 同左 + **Celery**（分布式任务）| 支持长时间内容发现任务 |
| **API网关** | FastAPI | 同左 | Python生态 + 高性能 |
| **监控系统** | Prometheus + Grafana | 同左 + **Sentry**（错误追踪）| 开源 + 可定制 + 错误管理 |
| **内容检索** | N/A | **Haystack** ⭐ 新增 | 统一的RAG框架 |
| **爬虫框架** | N/A | **Scrapy + BeautifulSoup** ⭐ 新增 | 高效的网页数据获取 |
| **多源API管理** | N/A | **自研API Gateway** ⭐ 新增 | 统一管理Wikipedia、YouTube等API |
| **LLM调用** | OpenAI SDK | **LiteLLM** ⭐ 新增 | 支持多LLM provider（OpenAI/Claude/本地模型）|
| **适配器注册** | N/A | **自研Plugin Registry** ⭐ 新增 | 动态加载和管理领域适配器 |
| **缓存策略** | Redis | Redis + **分层缓存** ⭐ 新增 | 会话/用户/全局三级缓存 |

---

## 2. 总体架构设计 ⭐ V2.0扩展

### 2.1 V2.0系统架构图 ⭐ 重大更新

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         应用层 (Application Layer)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ Web管理后台   │  │ API Gateway  │  │ 监控看板      │  │ 适配器市场    │    │
│  │ (React)      │  │ (FastAPI)    │  │ (Grafana)    │  │ (Marketplace)│    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│               领域适配器层 (Domain Adapter Layer) ⭐ 新增                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                   领域适配器注册中心 (Adapter Registry)                │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────┐  │   │
│  │  │ K12教育  │  │ 美术史   │  │ 畅销书   │  │ 职业培训 │  │ 更多  │  │   │
│  │  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │  │  ...  │  │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └───────┘  │   │
│  │          │            │            │            │            │        │   │
│  │          └────────────┴────────────┴────────────┴────────────┘        │   │
│  │                               ▼                                       │   │
│  │                        自动路由器 (Router)                             │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                  AI能力层 (AI Capability Layer) ⭐ 新增                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐       │
│  │  AI目标生成器      │  │  AI内容发现引擎    │  │  故事化叙述引擎    │       │
│  │ (Objective Gen)   │  │ (Content Disc)    │  │ (Narrative Gen)   │       │
│  ├───────────────────┤  ├───────────────────┤  ├───────────────────┤       │
│  │ • 需求解析        │  │ • 多源检索        │  │ • 叙述策略选择     │       │
│  │ • 受众分析        │  │ • 混合检索策略     │  │ • 故事元素生成     │       │
│  │ • 场景适配        │  │ • 质量评估        │  │ • 三幕式结构      │       │
│  │ • 目标生成        │  │ • 内容整合        │  │ • 质量控制        │       │
│  │ • 交互式完善      │  │ • 分层缓存        │  │ • 多样性支持      │       │
│  └───────────────────┘  └───────────────────┘  └───────────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                元工作流引擎层 (Meta-Workflow Engine Layer)                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │              工作流编排引擎 (Orchestration Engine)                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │参数解析器│→ │规则引擎  │→ │DAG构建器 │→ │执行调度器│            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ 状态管理器    │  │ 模型路由器    │  │ 配置管理器    │  │ 插件管理器    │    │
│  │ (State Mgr)  │  │ (LLM Router) │  │ (Config Mgr) │  │ (Plugin Mgr) │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│              多源内容层 (Multi-Source Content Layer) ⭐ 新增                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     内容源管理器 (Content Source Manager)             │   │
│  │                                                                        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Wikipedia   │  │  WikiArt    │  │  YouTube    │  │   豆瓣读书   │ │   │
│  │  │    API      │  │   Crawler   │  │  Data API   │  │   Crawler   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Goodreads   │  │   GitHub    │  │Stack Overflow│ │     MDN     │ │   │
│  │  │  (Archive)  │  │    API      │  │     API     │  │   Docs      │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │  │  混合检索引擎 (Hybrid Retrieval Engine)                          │ │   │
│  │  │  • BM25关键词检索  • 语义向量检索  • 结构化查询  • 图像检索     │ │   │
│  │  └─────────────────────────────────────────────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    数据存储层 (Data Storage Layer)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ PostgreSQL   │  │    Redis     │  │  Milvus/     │  │    MinIO     │    │
│  │ (关系型DB)    │  │  (缓存+队列) │  │  Qdrant      │  │ (对象存储)    │    │
│  │              │  │              │  │ (向量DB)     │  │              │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    外部服务层 (External Services Layer)                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ OpenAI GPT-4 │  │ Claude API   │  │ 本地LLM模型   │  │ 监控告警      │    │
│  │              │  │              │  │ (Ollama)     │  │ (Sentry)     │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 架构分层说明 ⭐ V2.0扩展

#### 2.2.1 应用层（Application Layer）

**V1.0原有组件**：
- Web管理后台：React/Vue前端，用户配置工作流
- API Gateway：FastAPI，对外提供RESTful API
- 监控看板：Grafana，系统运行监控

**V2.0新增组件** ⭐：
- **适配器市场（Adapter Marketplace）**：
  - 浏览和搜索第三方领域适配器
  - 适配器安装、更新、卸载
  - 用户评分和评论系统
  - 适配器版本管理

#### 2.2.2 领域适配器层（Domain Adapter Layer）⭐ 新增

**核心组件**：

**1. 领域适配器注册中心（Adapter Registry）**
- 功能：
  * 动态加载和卸载领域适配器
  * 适配器版本管理（v1.0、v1.1等）
  * 适配器依赖管理（某适配器依赖其他适配器）
  * 适配器健康检查（定期ping检查可用性）
  
- 数据结构：
  ```python
  class AdapterMetadata:
      id: str  # "art-history-adapter-v1.0"
      name: str  # "美术史教育适配器"
      version: str  # "1.0.0"
      author: str  # "官方/第三方开发者"
      domain: DomainType  # ART_HISTORY
      status: AdapterStatus  # ACTIVE / INACTIVE / DEPRECATED
      config_schema: dict  # 配置文件JSON Schema
      quality_metrics: QualityMetrics  # 质量指标（准确率、满意度）
  ```

**2. 自动路由器（Auto Router）**
- 功能：
  * 根据用户请求自动选择合适的领域适配器
  * 示例：请求包含"莫奈"→ 路由到美术史适配器
  * 支持多适配器协作（跨领域内容生成）
  
- 路由策略：
  ```python
  class RoutingStrategy:
      # 策略1：基于关键词匹配
      def route_by_keywords(request: UserRequest) -> DomainAdapter:
          if "莫奈" in request.topic or "印象派" in request.topic:
              return art_history_adapter
          elif "Excel" in request.topic or "数据分析" in request.topic:
              return vocational_training_adapter
          # ...
      
      # 策略2：基于LLM分类
      def route_by_llm_classification(request: UserRequest) -> DomainAdapter:
          prompt = f"判断以下需求属于哪个领域：{request.topic}"
          domain = llm.classify(prompt)
          return registry.get_adapter(domain)
  ```

**3. 已集成适配器（Initial Adapters）**
| 适配器 | 领域 | 状态 | 特色功能 |
|--------|------|------|---------|
| K12EducationAdapter | K12教育 | ✅ 官方 | 课标对齐、教学法选择、知识图谱 |
| ArtHistoryAdapter ⭐ | 美术史 | ✅ 官方 | WikiArt集成、艺术流派分类、作品分析 |
| BestsellerAdapter ⭐ | 畅销书 | ✅ 官方 | 书评整合、讨论提纲生成、金句提取 |
| VocationalAdapter ⭐ | 职业培训 | ✅ 官方 | 技能树规划、项目式学习、就业导向 |
| EnterpriseAdapter ⭐ | 企业培训 | 🚧 计划中 | 案例库、ROI评估、合规性检查 |

#### 2.2.3 AI能力层（AI Capability Layer）⭐ 新增

**核心引擎**：

**1. AI目标生成器（AI Objective Generator）**
- 架构：
  ```
  UserRequest → NLP解析 → 受众推断 → 场景分析 → LLM生成目标 → 结构化输出
  ```
  
- 关键技术：
  * Few-shot Prompting：提供领域示例引导LLM
  * Chain-of-Thought：让LLM解释推理过程
  * JSON Schema输出：保证输出结构化
  * 后处理验证：检查目标合理性（SMART原则）

**2. AI内容发现引擎（AI Content Discovery Engine）**
- 架构：
  ```
  学习目标 → 查询生成 → 多源并发检索 → 质量评估 → 结果融合 → 结构化整合
  ```
  
- 检索流程：
  ```python
  class ContentDiscoveryPipeline:
      def discover(self, objectives: List[LearningObjective]) -> StructuredContent:
          # Step 1: 查询生成
          queries = self.query_generator.generate(objectives)
          
          # Step 2: 多源并发检索
          results = await asyncio.gather(
              self.wikipedia_retriever.search(queries),
              self.wikiart_retriever.search(queries),
              self.youtube_retriever.search(queries),
              # ...
          )
          
          # Step 3: 质量评估
          scored_results = self.quality_assessor.score(results)
          
          # Step 4: 结果融合（RRF算法）
          fused_results = self.rrf_fusion.fuse(scored_results)
          
          # Step 5: 结构化整合
          structured = self.content_structurer.structure(fused_results)
          
          return structured
  ```

**3. 故事化叙述引擎（Narrative Generation Engine）**
- 架构：
  ```
  原始内容 → 策略选择 → 人物设计 → 情节构建 → 场景设计 → LLM生成 → 质量评估
  ```
  
- 生成流程：
  ```python
  class NarrativeGenerator:
      def generate_story(self, content: StructuredContent) -> NarrativeContent:
          # Step 1: 选择叙述策略
          strategy = self.strategy_selector.select(content.domain, content.audience)
          
          # Step 2: 生成故事元素
          character = self.character_generator.generate(strategy)
          plot = self.plot_generator.generate(strategy, content.key_points)
          scene = self.scene_generator.generate(strategy, character)
          
          # Step 3: 三幕式结构
          act1 = self.generate_setup(character, scene)
          act2 = self.generate_confrontation(plot, content.challenges)
          act3 = self.generate_resolution(plot, content.conclusions)
          
          # Step 4: LLM生成文本
          narrative = self.llm.generate_narrative(
              elements=[act1, act2, act3],
              style=strategy.style,
              tone=strategy.tone
          )
          
          # Step 5: 质量评估
          quality_score = self.quality_assessor.assess(narrative)
          if quality_score < threshold:
              narrative = self.refine_narrative(narrative, feedback)
          
          return narrative
  ```

#### 2.2.4 元工作流引擎层（Meta-Workflow Engine Layer）

**V1.0保留组件**：
- 参数解析器：解析用户输入的配置参数
- 规则引擎：根据规则选择工作流模板
- DAG构建器：动态构建工作流有向无环图
- 执行调度器：调度工作流节点执行
- 状态管理器：管理工作流执行状态
- 模型路由器：根据任务选择合适的LLM
- 配置管理器：管理系统和用户配置

**V2.0新增组件** ⭐：
- **插件管理器（Plugin Manager）**：
  * 管理领域适配器的生命周期
  * 热加载和热卸载适配器
  * 适配器沙箱隔离（避免相互影响）
  * 适配器权限控制

#### 2.2.5 多源内容层（Multi-Source Content Layer）⭐ 新增

**核心组件**：

**1. 内容源管理器（Content Source Manager）**
- 管理10+个内容数据源
- 统一的接口抽象（`ContentSource` interface）
- 配额管理和限流（避免超过API限制）
- 降级策略（主源失败→备用源）

**2. 内容源列表**：
| 内容源 | 领域 | 类型 | 访问方式 | 限制 |
|--------|------|------|---------|------|
| Wikipedia | 通用 | API | REST API | 200 req/s |
| WikiArt | 美术史 | 爬虫 | Scrapy | robots.txt |
| YouTube Data API | 职业培训/通用 | API | REST API | 10K units/day |
| 豆瓣读书 | 畅销书 | 爬虫 | Scrapy | 反爬虫 |
| Goodreads | 畅销书 | 归档数据 | 本地JSON | N/A |
| GitHub | 职业培训 | API | REST API | 5K req/hour |
| Stack Overflow | 职业培训 | API | REST API | 300 req/day |
| MDN Web Docs | 职业培训 | 爬虫 | Scrapy | 开放数据 |
| Google Arts & Culture | 美术史 | 爬虫 | Selenium | 反爬虫 |
| 课标文档 | K12 | 本地文件 | 文件系统 | N/A |

**3. 混合检索引擎（Hybrid Retrieval Engine）**
- **BM25关键词检索**：传统的关键词匹配
- **语义向量检索**：使用Embedding进行语义相似度匹配
- **结构化查询**：针对结构化数据（如JSON）的精确查询
- **图像检索**：基于图像相似度的检索（美术作品）
- **RRF融合算法**：Reciprocal Rank Fusion，融合多种检索结果

#### 2.2.6 数据存储层（Data Storage Layer）

**V1.0保留组件**：
- PostgreSQL：关系型数据库（用户、配置、工作流定义）
- Redis：缓存和消息队列
- Milvus/Qdrant：向量数据库（向量检索）

**V2.0新增组件** ⭐：
- **MinIO**：对象存储
  * 存储爬取的HTML/图片/视频
  * 存储生成的内容产物（PPT、PDF）
  * 支持大文件存储和CDN加速

#### 2.2.7 外部服务层（External Services Layer）

**V1.0保留组件**：
- OpenAI GPT-4：主力LLM

**V2.0新增组件** ⭐：
- **Claude API**：Anthropic的LLM，作为GPT-4的补充
- **本地LLM模型（Ollama）**：支持离线部署和成本优化
- **Sentry**：错误追踪和监控
│  │              └─────────────────┘                   │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 状态管理器    │  │ 模型路由器    │  │ 配置管理器    │          │
│  │ (State Mgr)  │  │ (LLM Router) │  │ (Config Mgr) │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
          │                  │                  │
┌─────────┼──────────────────┼──────────────────┼─────────────────┐
│         ▼                  ▼                  ▼                 │
│              组件服务层 (Component Layer)                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 画像分析节点  │  │ 教学设计节点  │  │ 脚本生成节点  │          │
│  │ PersonaNode  │  │ DesignNode   │  │ ScriptNode   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 校验节点      │  │ 媒体生成节点  │  │ 合规审查节点  │          │
│  │ VerifyNode   │  │ MediaNode    │  │ CompliNode   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
          │                  │                  │
┌─────────┼──────────────────┼──────────────────┼─────────────────┐
│         ▼                  ▼                  ▼                 │
│              基础设施层 (Infrastructure Layer)                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ LLM服务池     │  │ RAG检索引擎   │  │ 资源仓库      │          │
│  │ GPT/Claude   │  │ Vector DB    │  │ Asset Store  │          │
│  │ DeepSeek     │  │ Embedding    │  │ (S3/OSS)     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 关系数据库    │  │ 缓存层        │  │ 消息队列      │          │
│  │ PostgreSQL   │  │ Redis        │  │ RabbitMQ     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 架构分层详解

#### 2.2.1 应用层 (Application Layer)

**职责**：用户交互界面和API入口

**核心组件**：

1. **Web管理后台**
   - 工作流可视化配置界面
   - 工作流模板管理界面
   - 模板版本对比与回滚
   - 执行监控看板
   - 配置管理界面
   - 用户权限管理

2. **API Gateway**
   - RESTful API接口
   - 请求认证与鉴权
   - 限流与熔断
   - API文档自动生成

3. **监控看板**
   - 实时性能指标
   - 执行成功率
   - 成本分析
   - 告警管理

#### 2.2.2 元引擎层 (Meta-Engine Layer)

**职责**：工作流的动态构建和执行管理

**核心组件**：

1. **参数解析器 (Parameter Parser)**
   ```python
   class ParameterParser:
       def parse(self, raw_request):
           """
           解析用户请求，提取关键参数
           返回: WorkflowContext对象
           """
           context = WorkflowContext()
           context.subject = self.extract_subject(raw_request)
           context.grade = self.extract_grade(raw_request)
           context.region = self.extract_region(raw_request)
           context.pedagogy = self.infer_pedagogy(context)
           return context
   ```

2. **规则引擎 (Rule Engine)**
   ```python
   class RuleEngine:
       def match_rules(self, context):
           """
           根据上下文匹配配置规则
           返回: NodeList, ResourceBindings
           """
           rules = self.load_rules(context.subject, context.region)
           
           # 决定需要哪些节点
           required_nodes = []
           if context.subject == "数学":
               required_nodes.append("LogicVerificationNode")
           if context.region in ["US", "EU"]:
               required_nodes.append("GDPRComplianceNode")
           
           # 决定资源绑定
           knowledge_base = self.select_knowledge_base(context)
           prompt_template = self.select_prompt_template(context)
           
           return required_nodes, {
               "knowledge_base": knowledge_base,
               "prompt_template": prompt_template
           }
   ```

3. **DAG构建器 (DAG Builder)**
   ```python
   class DAGBuilder:
       def build(self, context, nodes, resources):
           """
           构建有向无环图工作流
           """
           graph = WorkflowGraph()
           
           # 添加基础骨架节点
           graph.add_node("persona_analysis", PersonaNode())
           graph.add_node("objective_setting", ObjectiveNode())
           graph.add_node("content_generation", ContentNode())
           
           # 根据规则插入条件节点
           for node_type in nodes:
               graph.insert_node(node_type, position="auto")
           
           # 绑定资源
           graph.bind_resources(resources)
           
           # 验证DAG合法性
           if not graph.is_acyclic():
               raise CyclicDependencyError()
           
           return graph
   ```

4. **执行调度器 (Scheduler)**
   - 管理工作流实例的生命周期
   - 调度节点执行顺序
   - 处理并发和依赖
   - 失败重试和降级
   - **记录执行轨迹**：保存每次执行的节点列表、参数、结果

5. **状态管理器 (State Manager)**
   - 持久化工作流执行状态
   - 支持断点续传
   - 提供状态查询接口

6. **模型路由器 (LLM Router)**
   - 根据节点类型选择合适的模型
   - 负载均衡
   - 超时和重试
   - 成本优化

7. **工作流模板管理器 (Workflow Template Manager)** ⭐ 新增
   - 模板生成与存储
   - 版本控制与对比
   - 模板查找与实例化
   - 执行记忆与关联

#### 2.2.3 组件服务层 (Component Layer)

**职责**：可插拔的业务节点实现

**节点接口规范**：
```python
from abc import ABC, abstractmethod

class WorkflowNode(ABC):
    """工作流节点基类"""
    
    @abstractmethod
    def execute(self, context: WorkflowContext) -> NodeOutput:
        """
        执行节点逻辑
        
        Args:
            context: 工作流上下文，包含所有前置节点的输出
        
        Returns:
            NodeOutput: 节点执行结果
        """
        pass
    
    @abstractmethod
    def validate_input(self, context: WorkflowContext) -> bool:
        """验证输入是否满足执行条件"""
        pass
    
    def on_error(self, error: Exception) -> ErrorHandlingStrategy:
        """错误处理策略"""
        return ErrorHandlingStrategy.RETRY
```

**核心节点实现**：

1. **画像分析节点 (PersonaNode)**
   ```python
   class PersonaNode(WorkflowNode):
       def execute(self, context):
           # 提取基础信息
           age = context.grade_to_age()
           subject = context.subject
           
           # 调用LLM生成画像
           prompt = self.build_prompt(age, subject, context.region)
           persona_data = self.llm.generate(prompt)
           
           # 结构化解析
           persona = self.parse_persona(persona_data)
           
           return NodeOutput(
               data=persona,
               metadata={"model": "gpt-4", "tokens": 1500}
           )
   ```

2. **教学设计节点 (DesignNode)**
   - 生成教学目标（知识/能力/情感）
   - 设计教学环节（导入/探究/总结/练习）
   - 分配时间和活动

3. **脚本生成节点 (ScriptNode)**
   - 基于RAG检索知识点
   - 应用故事化模板
   - 生成详细教学脚本

4. **校验节点 (VerifyNode)**
   - 事实核查（与知识库比对）
   - 逻辑一致性检查
   - 格式规范检查

5. **合规审查节点 (ComplianceNode)**
   - 敏感词扫描
   - 区域法规检查
   - 年龄适宜性评估

#### 2.2.4 基础设施层 (Infrastructure Layer)

**职责**：提供基础技术能力

**核心组件**：

1. **LLM服务池**
   - 多模型支持（GPT-4, Claude, DeepSeek等）
   - 连接池管理
   - 请求队列
   - 成本追踪

2. **RAG检索引擎**
   - 向量数据库（Milvus/Qdrant）
   - Embedding模型（OpenAI/BGE）
   - 混合检索（向量+关键词）
   - 重排序（Reranker）

3. **资源仓库**
   - 提示词模板库
   - 配置文件库
   - 生成内容归档
   - 媒体资源存储

4. **数据持久化**
   - PostgreSQL：结构化数据
   - Redis：缓存和会话
   - S3/OSS：大文件存储

---

## 3. 核心模块设计 ⭐ V2.0扩展

> **V2.0变化**：新增AI目标生成器、AI内容发现引擎、故事化叙述引擎、领域适配器架构4个核心模块

### 3.0 V2.0新增模块总览 ⭐

V2.0在V1.0的基础上，新增了4个核心模块，支持多领域内容生成：

```
┌──────────────────────────────────────────────────────────────┐
│              V2.0 核心模块架构                                 │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────────────────────────────────────┐     │
│  │  3.0.1 AI目标生成器 (AI Objective Generator)         │     │
│  │  从模糊需求自动推断学习目标                           │     │
│  └─────────────────────────────────────────────────────┘     │
│                           ↓                                   │
│  ┌─────────────────────────────────────────────────────┐     │
│  │  3.0.2 AI内容发现引擎 (AI Content Discovery Engine)  │     │
│  │  从多源实时检索和整合内容                            │     │
│  └─────────────────────────────────────────────────────┘     │
│                           ↓                                   │
│  ┌─────────────────────────────────────────────────────┐     │
│  │  3.0.3 故事化叙述引擎 (Narrative Generation Engine)  │     │
│  │  将知识点转化为引人入胜的故事                         │     │
│  └─────────────────────────────────────────────────────┘     │
│                           ↓                                   │
│  ┌─────────────────────────────────────────────────────┐     │
│  │  3.0.4 领域适配器架构 (Domain Adapter Architecture)  │     │
│  │  插件化的领域专业化机制                              │     │
│  └─────────────────────────────────────────────────────┘     │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

---

### 3.0.1 AI目标生成器架构设计 ⭐ 新增

#### 概述

**功能**：从用户的模糊需求自动生成结构化的学习目标。

**输入示例**：
```
"我想给中学生讲莫奈的《日出·印象》，大概10分钟"
```

**输出示例**：
```yaml
学习目标:
  - 了解印象派的诞生背景（19世纪法国艺术革新）
  - 理解莫奈的色彩运用和笔触技法
  - 感受《日出·印象》所传达的情感与氛围
  - 认识该画作对现代艺术的影响
课程时长: 10分钟  # 由用户输入指定，支持3-15分钟微课
```

---

#### 模块架构

```
┌──────────────────────────────────────────────────────────────┐
│             AI目标生成器 (AI Objective Generator)             │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐                                            │
│  │ 用户请求输入  │                                            │
│  │ (UserRequest)│                                            │
│  └──────┬───────┘                                            │
│         │                                                     │
│         ▼                                                     │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块1: 需求解析器 (Request Parser)    │                    │
│  │ • 提取主题（topic）                   │                    │
│  │ • 识别领域（domain）                  │                    │
│  │ • 解析约束条件（时长、受众等）         │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块2: 受众分析器 (Audience Analyzer) │                    │
│  │ • 推断年龄段                          │                    │
│  │ • 推断认知水平                        │                    │
│  │ • 推断学习动机                        │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块3: 场景适配器 (Context Adapter)   │                    │
│  │ • 根据时长调整目标深度                │                    │
│  │ • 根据场景调整目标类型                │                    │
│  │ • 考虑受众特征调整表述                │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块4: 目标生成器 (Objective Creator) │                    │
│  │ • LLM Prompt工程                      │                    │
│  │ • Few-shot示例引导                    │                    │
│  │ • JSON Schema输出                     │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块5: 后处理与验证 (Validator)       │                    │
│  │ • SMART原则检查                       │                    │
│  │ • 目标数量控制（3-5个）               │                    │
│  │ • 布鲁姆分类学层级验证                │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────────────────────┐                    │
│  │ 模块6: 交互式完善 (Interactive Refine)│                    │
│  │ • 展示给用户确认                      │                    │
│  │ • 支持用户修改调整                    │                    │
│  │ • 记录反馈用于优化                    │                    │
│  └──────────────┬───────────────────────┘                    │
│                 │                                             │
│                 ▼                                             │
│  ┌──────────────────────┐                                    │
│  │ 输出：学习目标列表    │                                    │
│  │ (LearningObjectives) │                                    │
│  └──────────────────────┘                                    │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

---

#### 核心类设计

**1. 数据模型**

```python
from dataclasses import dataclass
from typing import List, Optional
from enum import Enum

class DomainType(Enum):
    """领域类型"""
    K12_EDUCATION = "k12_education"
    ART_HISTORY = "art_history"
    BESTSELLER = "bestseller"
    VOCATIONAL_TRAINING = "vocational_training"
    ENTERPRISE_TRAINING = "enterprise_training"

class AudienceLevel(Enum):
    """受众认知水平"""
    BEGINNER = "beginner"          # 初学者
    INTERMEDIATE = "intermediate"  # 中级
    ADVANCED = "advanced"          # 高级
    EXPERT = "expert"              # 专家

class LearningContext(Enum):
    """学习场景"""
    CLASSROOM = "classroom"        # 课堂教学
    SELF_STUDY = "self_study"      # 自学
    WORKSHOP = "workshop"          # 工作坊
    ONLINE_COURSE = "online_course" # 在线课程
    READING_CLUB = "reading_club"  # 读书会

@dataclass
class UserRequest:
    """用户请求"""
    raw_input: str  # 原始输入："我想给中学生讲莫奈的《日出·印象》，10分钟"
    domain: Optional[DomainType] = None  # AI推断
    topic: Optional[str] = None  # AI提取
    audience: Optional[str] = None  # AI提取
    duration: Optional[int] = None  # 分钟，微课建议3-15分钟，常规课45-90分钟
    context: Optional[LearningContext] = None  # AI推断
    constraints: Optional[dict] = None  # 其他约束

@dataclass
class AudienceProfile:
    """受众画像"""
    age_range: tuple  # (12, 15) 中学生
    cognitive_level: AudienceLevel
    prior_knowledge: str  # "基本了解艺术史"
    motivation: str  # "培养艺术鉴赏能力"
    attention_span: int  # 注意力持续时间（分钟）

@dataclass
class LearningObjective:
    """学习目标"""
    description: str  # 目标描述
    bloom_level: str  # 布鲁姆分类学层级（remember/understand/apply/analyze/evaluate/create）
    priority: int  # 优先级（1最高）
    estimated_time: int  # 预计用时（分钟）
    assessment_method: Optional[str] = None  # 评估方法

@dataclass
class ObjectiveGenerationResult:
    """目标生成结果"""
    objectives: List[LearningObjective]
    audience_profile: AudienceProfile
    confidence_score: float  # 0-1，AI的置信度
    suggestions: List[str]  # 给用户的建议
    alternative_objectives: Optional[List[LearningObjective]] = None  # 备选目标
```

**2. 核心类实现**

```python
class AIObjectiveGenerator:
    """AI目标生成器"""
    
    def __init__(
        self,
        llm_client: LLMClient,
        domain_registry: DomainRegistry,
        config: ObjectiveGeneratorConfig
    ):
        self.llm = llm_client
        self.registry = domain_registry
        self.config = config
        
        # 子模块
        self.request_parser = RequestParser(llm_client)
        self.audience_analyzer = AudienceAnalyzer(llm_client)
        self.context_adapter = ContextAdapter()
        self.objective_creator = ObjectiveCreator(llm_client)
        self.validator = ObjectiveValidator()
        self.refiner = InteractiveRefiner()
    
    async def generate_objectives(
        self,
        user_request: UserRequest
    ) -> ObjectiveGenerationResult:
        """
        主流程：生成学习目标
        
        Args:
            user_request: 用户请求
            
        Returns:
            ObjectiveGenerationResult: 生成结果
        """
        # Step 1: 解析需求
        parsed_request = await self.request_parser.parse(user_request)
        logger.info(f"Parsed request: domain={parsed_request.domain}, topic={parsed_request.topic}")
        
        # Step 2: 受众分析
        audience_profile = await self.audience_analyzer.analyze(parsed_request)
        logger.info(f"Audience: age={audience_profile.age_range}, level={audience_profile.cognitive_level}")
        
        # Step 3: 场景适配
        adapted_context = self.context_adapter.adapt(
            parsed_request,
            audience_profile
        )
        logger.info(f"Context adapted: duration={adapted_context.duration}, depth={adapted_context.depth}")
        
        # Step 4: 生成目标
        raw_objectives = await self.objective_creator.create(
            parsed_request,
            audience_profile,
            adapted_context
        )
        logger.info(f"Generated {len(raw_objectives)} raw objectives")
        
        # Step 5: 验证和优化
        validated_objectives = self.validator.validate(
            raw_objectives,
            audience_profile,
            adapted_context
        )
        logger.info(f"Validated {len(validated_objectives)} objectives")
        
        # Step 6: 生成备选方案
        alternative_objectives = await self._generate_alternatives(
            validated_objectives,
            audience_profile
        )
        
        # 构建结果
        result = ObjectiveGenerationResult(
            objectives=validated_objectives,
            audience_profile=audience_profile,
            confidence_score=self._calculate_confidence(validated_objectives),
            suggestions=self._generate_suggestions(validated_objectives, audience_profile),
            alternative_objectives=alternative_objectives
        )
        
        return result
    
    async def refine_objectives(
        self,
        current_objectives: List[LearningObjective],
        user_feedback: dict
    ) -> List[LearningObjective]:
        """
        交互式完善目标
        
        Args:
            current_objectives: 当前目标
            user_feedback: 用户反馈
            
        Returns:
            改进后的目标列表
        """
        return await self.refiner.refine(current_objectives, user_feedback)
    
    def _calculate_confidence(self, objectives: List[LearningObjective]) -> float:
        """计算置信度分数"""
        # 基于目标数量、SMART合规性、布鲁姆层级分布等计算
        score = 0.0
        
        # 目标数量合理性（3-5个为佳）
        if 3 <= len(objectives) <= 5:
            score += 0.3
        
        # SMART原则合规性
        smart_compliance = sum(1 for obj in objectives if self._is_smart(obj)) / len(objectives)
        score += smart_compliance * 0.4
        
        # 布鲁姆层级分布合理性（涵盖多个层级）
        bloom_levels = set(obj.bloom_level for obj in objectives)
        score += (len(bloom_levels) / 6) * 0.3  # 6个布鲁姆层级
        
        return min(score, 1.0)
    
    def _is_smart(self, objective: LearningObjective) -> bool:
        """检查目标是否符合SMART原则"""
        # Specific, Measurable, Achievable, Relevant, Time-bound
        # 简化实现：检查是否有明确的动词、可测量的结果
        return (
            len(objective.description) > 10 and
            objective.bloom_level is not None and
            objective.estimated_time > 0
        )
    
    def _generate_suggestions(
        self,
        objectives: List[LearningObjective],
        audience: AudienceProfile
    ) -> List[str]:
        """生成给用户的建议"""
        suggestions = []
        
        total_time = sum(obj.estimated_time for obj in objectives)
        if total_time > 60:
            suggestions.append("建议：总时长超过60分钟，考虑拆分为多个模块")
        
        if audience.cognitive_level == AudienceLevel.BEGINNER:
            has_advanced_objectives = any(
                obj.bloom_level in ["analyze", "evaluate", "create"]
                for obj in objectives
            )
            if has_advanced_objectives:
                suggestions.append("注意：初学者可能对高阶目标感到困难，建议增加铺垫")
        
        return suggestions
    
    async def _generate_alternatives(
        self,
        objectives: List[LearningObjective],
        audience: AudienceProfile
    ) -> List[LearningObjective]:
        """生成备选目标"""
        # 为不同难度级别生成备选方案
        prompt = f"""
        当前学习目标：
        {[obj.description for obj in objectives]}
        
        受众：{audience.age_range}岁，认知水平：{audience.cognitive_level}
        
        请生成2个备选方案：
        1. 更简化的版本（降低难度）
        2. 更深入的版本（提高难度）
        """
        
        # 调用LLM生成备选方案
        # （实际实现略）
        return []
```

**3. 子模块实现示例**

```python
class RequestParser:
    """需求解析器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def parse(self, request: UserRequest) -> UserRequest:
        """
        解析用户请求，提取关键信息
        
        使用LLM进行结构化提取
        """
        prompt = f"""
        请从以下用户需求中提取关键信息：
        
        需求：{request.raw_input}
        
        请以JSON格式输出：
        {{
            "domain": "领域（k12_education/art_history/bestseller/vocational_training）",
            "topic": "主题",
            "audience": "目标受众",
            "duration": 时长（分钟）,
            "context": "场景（classroom/self_study/workshop等）",
            "constraints": {{"其他约束": "值"}}
        }}
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            response_format={"type": "json_object"},
            temperature=0.3  # 低温度，保证稳定性
        )
        
        extracted = json.loads(response.content)
        
        # 更新UserRequest
        request.domain = DomainType(extracted["domain"])
        request.topic = extracted["topic"]
        request.audience = extracted["audience"]
        request.duration = extracted.get("duration")
        request.context = LearningContext(extracted["context"]) if extracted.get("context") else None
        request.constraints = extracted.get("constraints", {})
        
        return request


class AudienceAnalyzer:
    """受众分析器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def analyze(self, request: UserRequest) -> AudienceProfile:
        """
        分析目标受众特征
        
        基于受众描述推断年龄、认知水平、先验知识等
        """
        prompt = f"""
        请分析以下受众的特征：
        
        受众描述：{request.audience}
        学习主题：{request.topic}
        领域：{request.domain.value}
        
        请以JSON格式输出：
        {{
            "age_range": [最小年龄, 最大年龄],
            "cognitive_level": "beginner/intermediate/advanced/expert",
            "prior_knowledge": "描述先验知识",
            "motivation": "学习动机",
            "attention_span": 注意力持续时间（分钟）
        }}
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            response_format={"type": "json_object"},
            temperature=0.5
        )
        
        data = json.loads(response.content)
        
        return AudienceProfile(
            age_range=tuple(data["age_range"]),
            cognitive_level=AudienceLevel(data["cognitive_level"]),
            prior_knowledge=data["prior_knowledge"],
            motivation=data["motivation"],
            attention_span=data["attention_span"]
        )


class ObjectiveCreator:
    """目标生成器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
        self.few_shot_examples = self._load_few_shot_examples()
    
    async def create(
        self,
        request: UserRequest,
        audience: AudienceProfile,
        context: dict
    ) -> List[LearningObjective]:
        """
        生成学习目标
        
        使用Few-shot Prompting引导LLM生成高质量目标
        """
        # 根据领域选择Few-shot示例
        examples = self._select_examples(request.domain)
        
        prompt = f"""
        请为以下学习内容设计学习目标：
        
        主题：{request.topic}
        领域：{request.domain.value}
        受众：{audience.age_range}岁，认知水平：{audience.cognitive_level.value}
        时长：{request.duration}分钟
        场景：{request.context.value if request.context else '未指定'}
        
        受众分析：
        - 先验知识：{audience.prior_knowledge}
        - 学习动机：{audience.motivation}
        - 注意力持续时间：{audience.attention_span}分钟
        
        参考示例：
        {self._format_examples(examples)}
        
        要求：
        1. 生成3-5个学习目标
        2. 符合SMART原则（具体、可测量、可达成、相关、有时限）
        3. 覆盖布鲁姆分类学的多个层级
        4. 总用时不超过{request.duration}分钟
        5. 适合{audience.cognitive_level.value}水平的受众
        
        请以JSON格式输出：
        {{
            "objectives": [
                {{
                    "description": "目标描述",
                    "bloom_level": "remember/understand/apply/analyze/evaluate/create",
                    "priority": 1-5,
                    "estimated_time": 分钟数,
                    "assessment_method": "评估方法"
                }},
                ...
            ]
        }}
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            response_format={"type": "json_object"},
            temperature=0.7,  # 中等温度，保持创造性
            max_tokens=1500
        )
        
        data = json.loads(response.content)
        
        objectives = [
            LearningObjective(**obj_data)
            for obj_data in data["objectives"]
        ]
        
        return objectives
    
    def _load_few_shot_examples(self) -> dict:
        """加载Few-shot示例库"""
        return {
            DomainType.ART_HISTORY: [
                {
                    "topic": "梵高《星夜》",
                    "audience": "成人艺术爱好者",
                    "objectives": [
                        {
                            "description": "了解梵高创作《星夜》的背景（圣雷米精神病院时期）",
                            "bloom_level": "remember",
                            "priority": 1,
                            "estimated_time": 5
                        },
                        {
                            "description": "理解后印象派的色彩运用和笔触技法",
                            "bloom_level": "understand",
                            "priority": 2,
                            "estimated_time": 8
                        },
                        {
                            "description": "分析《星夜》中旋转的星空与柏树的象征意义",
                            "bloom_level": "analyze",
                            "priority": 3,
                            "estimated_time": 10
                        }
                    ]
                }
            ],
            DomainType.VOCATIONAL_TRAINING: [
                {
                    "topic": "Excel数据透视表",
                    "audience": "市场部员工",
                    "objectives": [
                        {
                            "description": "掌握数据透视表的创建流程",
                            "bloom_level": "apply",
                            "priority": 1,
                            "estimated_time": 15
                        },
                        {
                            "description": "能够使用数据透视表分析销售数据",
                            "bloom_level": "apply",
                            "priority": 2,
                            "estimated_time": 20
                        }
                    ]
                }
            ]
            # 更多领域示例...
        }
    
    def _select_examples(self, domain: DomainType) -> list:
        """根据领域选择相关示例"""
        return self.few_shot_examples.get(domain, [])
    
    def _format_examples(self, examples: list) -> str:
        """格式化示例为Prompt文本"""
        formatted = []
        for ex in examples:
            formatted.append(f"""
            示例：
            主题：{ex['topic']}
            受众：{ex['audience']}
            学习目标：
            {json.dumps(ex['objectives'], indent=2, ensure_ascii=False)}
            """)
        return "\n".join(formatted)


class ObjectiveValidator:
    """目标验证器"""
    
    def validate(
        self,
        objectives: List[LearningObjective],
        audience: AudienceProfile,
        context: dict
    ) -> List[LearningObjective]:
        """
        验证和优化目标列表
        
        检查：
        1. SMART原则合规性
        2. 布鲁姆层级分布
        3. 时间分配合理性
        4. 目标数量合理性
        """
        validated = []
        
        for obj in objectives:
            # 检查SMART原则
            if not self._check_smart(obj):
                logger.warning(f"Objective不符合SMART原则: {obj.description}")
                obj = self._fix_smart(obj)
            
            # 检查布鲁姆层级合理性
            if not self._check_bloom_level(obj, audience):
                logger.warning(f"Objective的布鲁姆层级不适合受众: {obj.description}")
                obj = self._adjust_bloom_level(obj, audience)
            
            validated.append(obj)
        
        # 检查总时长
        total_time = sum(obj.estimated_time for obj in validated)
        if context.get("duration") and total_time > context["duration"]:
            logger.warning(f"总时长({total_time}min)超过限制({context['duration']}min)")
            validated = self._adjust_time_allocation(validated, context["duration"])
        
        # 检查目标数量
        if len(validated) > 5:
            logger.warning(f"目标数量({len(validated)})过多，保留优先级最高的5个")
            validated = sorted(validated, key=lambda x: x.priority)[:5]
        elif len(validated) < 3:
            logger.warning(f"目标数量({len(validated)})过少")
        
        return validated
    
    def _check_smart(self, obj: LearningObjective) -> bool:
        """检查是否符合SMART原则"""
        # 简化实现
        return (
            len(obj.description) > 10 and
            obj.bloom_level is not None and
            obj.estimated_time > 0 and
            obj.priority > 0
        )
    
    def _fix_smart(self, obj: LearningObjective) -> LearningObjective:
        """修复SMART问题"""
        # 简化实现：添加默认值
        if obj.estimated_time == 0:
            obj.estimated_time = 10
        if obj.priority == 0:
            obj.priority = 3
        return obj
    
    def _check_bloom_level(self, obj: LearningObjective, audience: AudienceProfile) -> bool:
        """检查布鲁姆层级是否适合受众"""
        high_level_blooms = ["analyze", "evaluate", "create"]
        
        if audience.cognitive_level == AudienceLevel.BEGINNER:
            # 初学者不宜太多高阶目标
            return obj.bloom_level not in high_level_blooms
        
        return True
    
    def _adjust_bloom_level(self, obj: LearningObjective, audience: AudienceProfile) -> LearningObjective:
        """调整布鲁姆层级"""
        if audience.cognitive_level == AudienceLevel.BEGINNER:
            # 将高阶目标降级
            level_map = {
                "analyze": "understand",
                "evaluate": "understand",
                "create": "apply"
            }
            obj.bloom_level = level_map.get(obj.bloom_level, obj.bloom_level)
        
        return obj
    
    def _adjust_time_allocation(self, objectives: List[LearningObjective], max_duration: int) -> List[LearningObjective]:
        """调整时间分配"""
        total_time = sum(obj.estimated_time for obj in objectives)
        scale_factor = max_duration / total_time * 0.9  # 留10%缓冲
        
        for obj in objectives:
            obj.estimated_time = int(obj.estimated_time * scale_factor)
        
        return objectives
```

---

#### Prompt工程示例

**目标生成Prompt模板**：

```python
OBJECTIVE_GENERATION_PROMPT = """
你是一位经验丰富的教学设计专家，擅长为不同领域的学习内容设计高质量的学习目标。

# 任务
请为以下学习内容设计学习目标：

## 基本信息
- **主题**：{topic}
- **领域**：{domain}
- **受众**：{audience_description}
- **时长**：{duration}分钟
- **场景**：{learning_context}

## 受众分析
- **年龄范围**：{age_range}
- **认知水平**：{cognitive_level}
- **先验知识**：{prior_knowledge}
- **学习动机**：{motivation}
- **注意力持续时间**：{attention_span}分钟

## 设计要求
1. 生成3-5个学习目标
2. 符合SMART原则（具体、可测量、可达成、相关、有时限）
3. 覆盖布鲁姆分类学的多个层级：
   - Remember（记忆）：回忆事实、概念
   - Understand（理解）：解释思想、概念
   - Apply（应用）：在新情境中使用信息
   - Analyze（分析）：区分部分和整体
   - Evaluate（评价）：根据标准做出判断
   - Create（创造）：产生新的想法或产品
4. 总用时不超过{duration}分钟
5. 适合{cognitive_level}水平的受众

## 参考示例
{few_shot_examples}

## 输出格式
请以JSON格式输出：
```json
{{
  "objectives": [
    {{
      "description": "使用动作动词开头的明确目标描述",
      "bloom_level": "remember/understand/apply/analyze/evaluate/create",
      "priority": 1,
      "estimated_time": 10,
      "assessment_method": "如何评估该目标是否达成"
    }}
  ],
  "rationale": "设计这些目标的理由",
  "notes": "给教师/讲师的额外建议"
}}
```

# 注意事项
- 使用清晰的动作动词（如：识别、解释、应用、分析、评价、创造）
- 避免模糊的表述（如："了解一些"、"大概知道"）
- 考虑受众的先验知识，避免过难或过简单
- 确保目标之间有逻辑递进关系
- 为每个目标分配合理的时间
"""
```

---

#### 性能优化

**1. 缓存策略**

```python
class ObjectiveGeneratorCache:
    """目标生成器缓存"""
    
    def __init__(self, redis_client: Redis):
        self.redis = redis_client
        self.ttl = 7 * 24 * 3600  # 7天
    
    def get_cached_objectives(self, request_hash: str) -> Optional[ObjectiveGenerationResult]:
        """获取缓存的目标"""
        cached = self.redis.get(f"objectives:{request_hash}")
        if cached:
            return ObjectiveGenerationResult(**json.loads(cached))
        return None
    
    def cache_objectives(self, request_hash: str, result: ObjectiveGenerationResult):
        """缓存目标"""
        self.redis.setex(
            f"objectives:{request_hash}",
            self.ttl,
            json.dumps(result.__dict__, default=str)
        )
    
    def calculate_request_hash(self, request: UserRequest) -> str:
        """计算请求哈希"""
        # 基于关键字段计算哈希
        key_fields = f"{request.domain}|{request.topic}|{request.audience}|{request.duration}"
        return hashlib.md5(key_fields.encode()).hexdigest()
```

**2. 批处理优化**

```python
async def batch_generate_objectives(
    requests: List[UserRequest]
) -> List[ObjectiveGenerationResult]:
    """批量生成目标（并发处理）"""
    tasks = [generator.generate_objectives(req) for req in requests]
    results = await asyncio.gather(*tasks, return_exceptions=True)
    
    # 处理异常
    for i, result in enumerate(results):
        if isinstance(result, Exception):
            logger.error(f"Request {i} failed: {result}")
            results[i] = None
    
    return [r for r in results if r is not None]
```

---

#### 监控指标

```python
class ObjectiveGeneratorMetrics:
    """目标生成器监控指标"""
    
    def __init__(self):
        self.generation_count = Counter('objective_generation_total', 'Total objective generations')
        self.generation_duration = Histogram('objective_generation_duration_seconds', 'Duration of objective generation')
        self.confidence_score = Histogram('objective_confidence_score', 'Confidence score of generated objectives')
        self.user_acceptance_rate = Gauge('objective_user_acceptance_rate', 'Rate of user accepting generated objectives')
    
    def record_generation(self, duration: float, confidence: float, accepted: bool):
        """记录一次目标生成"""
        self.generation_count.inc()
        self.generation_duration.observe(duration)
        self.confidence_score.observe(confidence)
        if accepted:
            self.user_acceptance_rate.inc()
```

---

### 3.0.2 AI内容发现引擎架构设计 ⭐ 新增

#### 概述

**功能**：从多个公开数据源实时检索、评估和整合内容，为内容生成提供高质量素材。

**核心挑战**：
- V1.0依赖预设知识库，无法覆盖所有领域
- V2.0支持多领域，需要从全网动态获取内容
- 不同数据源格式、质量、可靠性差异大
- 需要智能筛选和整合海量信息

**解决方案**：
- 多源并发检索（Wikipedia、WikiArt、YouTube等10+数据源）
- 混合检索策略（BM25关键词 + 语义向量 + 结构化查询）
- 质量评估模型（权威性、相关性、时效性等5维度）
- 分层缓存（会话/用户/全局三级）

---

#### 模块架构

```
┌────────────────────────────────────────────────────────────────────┐
│          AI内容发现引擎 (AI Content Discovery Engine)               │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────┐                                              │
│  │ 学习目标输入      │                                              │
│  │ (LearningObjectives)                                            │
│  └─────────┬────────┘                                              │
│            │                                                        │
│            ▼                                                        │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块1: 查询生成器 (Query Generator)                 │           │
│  │ • 从学习目标提取关键词                              │           │
│  │ • 生成多种查询变体（同义词、相关词）                │           │
│  │ • 针对不同数据源优化查询                            │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块2: 多源检索器 (Multi-Source Retriever)          │           │
│  │                                                     │           │
│  │  ┌──────────────────────────────────────────────┐  │           │
│  │  │  并发检索调度器 (Concurrent Scheduler)        │  │           │
│  │  └──────────────────────────────────────────────┘  │           │
│  │           │          │          │          │        │           │
│  │           ▼          ▼          ▼          ▼        │           │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────┐    │           │
│  │  │Wikipedia│  │ WikiArt │  │ YouTube │  │... │    │           │
│  │  │Retriever│  │Retriever│  │Retriever│  │    │    │           │
│  │  └─────────┘  └─────────┘  └─────────┘  └────┘    │           │
│  │       │            │            │           │       │           │
│  │       └────────────┴────────────┴───────────┘       │           │
│  │                      │                              │           │
│  │                      ▼                              │           │
│  │  ┌──────────────────────────────────────────────┐  │           │
│  │  │  结果聚合器 (Result Aggregator)               │  │           │
│  │  └──────────────────────────────────────────────┘  │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块3: 混合检索引擎 (Hybrid Retrieval Engine)       │           │
│  │ • BM25 关键词检索                                   │           │
│  │ • Embedding 语义检索                                │           │
│  │ • 结构化查询（JSON/SQL）                            │           │
│  │ • 图像相似度检索（美术作品）                        │           │
│  │ • RRF融合算法                                       │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块4: 内容质量评估器 (Quality Assessor)            │           │
│  │ • 权威性评分（数据源可信度）                        │           │
│  │ • 相关性评分（与学习目标匹配度）                    │           │
│  │ • 时效性评分（内容新鲜度）                          │           │
│  │ • 完整性评分（信息丰富度）                          │           │
│  │ • 可读性评分（语言质量）                            │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块5: 内容整合器 (Content Integrator)              │           │
│  │ • 去重（相似内容合并）                              │           │
│  │ • 结构化（转换为统一格式）                          │           │
│  │ • 补充（填补信息缺口）                              │           │
│  │ • 标注（来源、可信度、使用建议）                    │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块6: 缓存管理器 (Cache Manager)                   │           │
│  │ • 会话缓存（Session Cache）：1小时                 │           │
│  │ • 用户缓存（User Cache）：1天                      │           │
│  │ • 全局缓存（Global Cache）：7天                    │           │
│  │ • LRU淘汰策略                                       │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌──────────────────────────┐                                      │
│  │ 输出：结构化内容          │                                      │
│  │ (StructuredContent)      │                                      │
│  └──────────────────────────┘                                      │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

---

#### 核心类设计

**1. 数据模型**

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from enum import Enum
from datetime import datetime

class ContentSourceType(Enum):
    """内容源类型"""
    WIKIPEDIA = "wikipedia"
    WIKIART = "wikiart"
    YOUTUBE = "youtube"
    DOUBAN = "douban"
    GOODREADS = "goodreads"
    GITHUB = "github"
    STACKOVERFLOW = "stackoverflow"
    MDN = "mdn"
    GOOGLE_ARTS = "google_arts_culture"
    ACADEMIC_DB = "academic_database"

class RetrievalStrategy(Enum):
    """检索策略"""
    KEYWORD = "keyword"       # BM25关键词检索
    SEMANTIC = "semantic"     # 语义向量检索
    STRUCTURED = "structured" # 结构化查询
    IMAGE = "image"          # 图像检索
    HYBRID = "hybrid"        # 混合检索

@dataclass
class ContentSource:
    """内容源配置"""
    type: ContentSourceType
    name: str
    base_url: str
    api_key: Optional[str] = None
    rate_limit: int = 100  # 每分钟请求数
    timeout: int = 10  # 超时时间（秒）
    retry_count: int = 3
    priority: int = 1  # 优先级（1最高）
    enabled: bool = True

@dataclass
class Query:
    """检索查询"""
    keywords: List[str]  # 关键词列表
    domain: str  # 领域
    language: str = "zh"  # 语言
    filters: Dict[str, Any] = field(default_factory=dict)  # 过滤条件
    limit: int = 10  # 返回结果数量

@dataclass
class ContentItem:
    """内容条目"""
    id: str  # 唯一标识
    source: ContentSourceType  # 来源
    title: str  # 标题
    content: str  # 内容正文
    url: str  # 原始URL
    author: Optional[str] = None
    publish_date: Optional[datetime] = None
    metadata: Dict[str, Any] = field(default_factory=dict)  # 元数据
    
    # 质量评分
    authority_score: float = 0.0  # 权威性 0-1
    relevance_score: float = 0.0  # 相关性 0-1
    timeliness_score: float = 0.0  # 时效性 0-1
    completeness_score: float = 0.0  # 完整性 0-1
    readability_score: float = 0.0  # 可读性 0-1
    
    # 聚合评分
    overall_score: float = 0.0  # 综合评分 0-1

@dataclass
class StructuredContent:
    """结构化内容"""
    topic: str  # 主题
    domain: str  # 领域
    summary: str  # 摘要
    key_points: List[str]  # 核心要点
    details: Dict[str, Any]  # 详细信息
    sources: List[ContentItem]  # 原始来源
    references: List[str]  # 参考链接
    images: List[str] = field(default_factory=list)  # 相关图片
    videos: List[str] = field(default_factory=list)  # 相关视频
    created_at: datetime = field(default_factory=datetime.now)
    confidence_score: float = 0.0  # 置信度

@dataclass
class RetrievalResult:
    """检索结果"""
    query: Query
    items: List[ContentItem]
    total_count: int
    search_time: float  # 搜索耗时（秒）
    from_cache: bool = False
```

**2. 核心类实现**

```python
class AIContentDiscoveryEngine:
    """AI内容发现引擎"""
    
    def __init__(
        self,
        source_manager: ContentSourceManager,
        retrieval_engine: HybridRetrievalEngine,
        quality_assessor: QualityAssessor,
        content_integrator: ContentIntegrator,
        cache_manager: CacheManager,
        config: DiscoveryEngineConfig
    ):
        self.source_manager = source_manager
        self.retrieval_engine = retrieval_engine
        self.quality_assessor = quality_assessor
        self.integrator = content_integrator
        self.cache = cache_manager
        self.config = config
        
        # 子模块
        self.query_generator = QueryGenerator()
    
    async def discover_content(
        self,
        objectives: List[LearningObjective],
        domain: DomainType,
        user_id: Optional[str] = None
    ) -> StructuredContent:
        """
        主流程：发现和整合内容
        
        Args:
            objectives: 学习目标列表
            domain: 领域类型
            user_id: 用户ID（用于缓存）
            
        Returns:
            StructuredContent: 结构化的整合内容
        """
        start_time = time.time()
        
        # Step 1: 生成查询
        queries = self.query_generator.generate_queries(objectives, domain)
        logger.info(f"Generated {len(queries)} queries for {len(objectives)} objectives")
        
        # Step 2: 检查缓存
        cache_key = self._generate_cache_key(queries, domain)
        cached_content = await self.cache.get(cache_key, user_id)
        if cached_content:
            logger.info("Content found in cache")
            return cached_content
        
        # Step 3: 多源并发检索
        retrieval_results = await self._retrieve_from_multiple_sources(
            queries,
            domain
        )
        logger.info(f"Retrieved {sum(len(r.items) for r in retrieval_results)} items from {len(retrieval_results)} sources")
        
        # Step 4: 质量评估
        scored_items = await self.quality_assessor.assess_batch(
            retrieval_results,
            objectives
        )
        logger.info(f"Assessed {len(scored_items)} items")
        
        # Step 5: 过滤低质量内容
        filtered_items = [
            item for item in scored_items
            if item.overall_score >= self.config.min_quality_threshold
        ]
        logger.info(f"Filtered to {len(filtered_items)} high-quality items")
        
        # Step 6: 内容整合
        structured_content = await self.integrator.integrate(
            filtered_items,
            objectives,
            domain
        )
        logger.info(f"Integrated content: {len(structured_content.key_points)} key points")
        
        # Step 7: 缓存结果
        await self.cache.set(cache_key, structured_content, user_id)
        
        elapsed_time = time.time() - start_time
        logger.info(f"Content discovery completed in {elapsed_time:.2f}s")
        
        structured_content.metadata = {
            "discovery_time": elapsed_time,
            "sources_count": len(set(item.source for item in filtered_items)),
            "total_items": len(filtered_items)
        }
        
        return structured_content
    
    async def _retrieve_from_multiple_sources(
        self,
        queries: List[Query],
        domain: DomainType
    ) -> List[RetrievalResult]:
        """
        从多个数据源并发检索
        
        使用asyncio并发调用多个检索器
        """
        # 根据领域选择合适的数据源
        sources = self.source_manager.get_sources_for_domain(domain)
        
        # 为每个数据源创建检索任务
        tasks = []
        for source in sources:
            if not source.enabled:
                continue
            
            retriever = self.source_manager.get_retriever(source.type)
            for query in queries:
                task = self._retrieve_with_timeout(
                    retriever,
                    query,
                    source
                )
                tasks.append(task)
        
        # 并发执行所有检索任务
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # 过滤异常结果
        valid_results = []
        for i, result in enumerate(results):
            if isinstance(result, Exception):
                logger.error(f"Retrieval task {i} failed: {result}")
            elif result is not None:
                valid_results.append(result)
        
        return valid_results
    
    async def _retrieve_with_timeout(
        self,
        retriever: BaseRetriever,
        query: Query,
        source: ContentSource
    ) -> Optional[RetrievalResult]:
        """
        带超时的检索，包含重试机制
        """
        for attempt in range(source.retry_count):
            try:
                result = await asyncio.wait_for(
                    retriever.retrieve(query),
                    timeout=source.timeout
                )
                return result
            except asyncio.TimeoutError:
                logger.warning(f"{source.name} retrieval timeout (attempt {attempt + 1}/{source.retry_count})")
                if attempt == source.retry_count - 1:
                    return None
            except Exception as e:
                logger.error(f"{source.name} retrieval error: {e}")
                return None
    
    def _generate_cache_key(self, queries: List[Query], domain: DomainType) -> str:
        """生成缓存键"""
        # 基于查询关键词和领域生成哈希
        key_parts = [domain.value]
        for q in queries:
            key_parts.extend(sorted(q.keywords))
        
        key_string = "|".join(key_parts)
        return hashlib.md5(key_string.encode()).hexdigest()


class QueryGenerator:
    """查询生成器"""
    
    def __init__(self, llm_client: Optional[LLMClient] = None):
        self.llm = llm_client
    
    def generate_queries(
        self,
        objectives: List[LearningObjective],
        domain: DomainType
    ) -> List[Query]:
        """
        从学习目标生成检索查询
        
        策略：
        1. 提取关键词（NER + 关键短语提取）
        2. 生成同义词和相关词
        3. 针对不同数据源优化查询
        """
        queries = []
        
        for obj in objectives:
            # 提取关键词
            keywords = self._extract_keywords(obj.description)
            
            # 生成同义词
            expanded_keywords = self._expand_keywords(keywords)
            
            # 创建查询
            query = Query(
                keywords=expanded_keywords,
                domain=domain.value,
                filters={"bloom_level": obj.bloom_level}
            )
            queries.append(query)
        
        return queries
    
    def _extract_keywords(self, text: str) -> List[str]:
        """提取关键词"""
        # 简化实现：分词 + 停用词过滤
        # 实际可使用 jieba、spaCy 等工具
        words = text.split()
        stopwords = {"的", "了", "和", "是", "在", "a", "the", "and", "of"}
        keywords = [w for w in words if w not in stopwords and len(w) > 1]
        return keywords[:5]  # 限制关键词数量
    
    def _expand_keywords(self, keywords: List[str]) -> List[str]:
        """扩展关键词（同义词、相关词）"""
        # 简化实现
        # 实际可使用词向量、同义词词典等
        expanded = keywords.copy()
        
        # 示例：为"莫奈"添加相关词
        synonyms_map = {
            "莫奈": ["Claude Monet", "印象派"],
            "数据透视表": ["pivot table", "数据分析"]
        }
        
        for kw in keywords:
            if kw in synonyms_map:
                expanded.extend(synonyms_map[kw])
        
        return list(set(expanded))  # 去重


class HybridRetrievalEngine:
    """混合检索引擎"""
    
    def __init__(
        self,
        bm25_index: BM25Index,
        vector_db: VectorDatabase,
        config: RetrievalConfig
    ):
        self.bm25 = bm25_index
        self.vector_db = vector_db
        self.config = config
    
    async def hybrid_search(
        self,
        query: Query,
        strategy: RetrievalStrategy = RetrievalStrategy.HYBRID
    ) -> List[ContentItem]:
        """
        混合检索
        
        结合多种检索策略，使用RRF融合结果
        """
        results = []
        
        if strategy in [RetrievalStrategy.KEYWORD, RetrievalStrategy.HYBRID]:
            # BM25 关键词检索
            bm25_results = await self._bm25_search(query)
            results.append(("bm25", bm25_results))
        
        if strategy in [RetrievalStrategy.SEMANTIC, RetrievalStrategy.HYBRID]:
            # 语义向量检索
            semantic_results = await self._semantic_search(query)
            results.append(("semantic", semantic_results))
        
        # RRF 融合
        if strategy == RetrievalStrategy.HYBRID and len(results) > 1:
            fused_results = self._reciprocal_rank_fusion(results)
            return fused_results
        elif results:
            return results[0][1]
        else:
            return []
    
    async def _bm25_search(self, query: Query) -> List[ContentItem]:
        """BM25关键词检索"""
        query_text = " ".join(query.keywords)
        scores = self.bm25.get_scores(query_text)
        
        # 获取Top-K结果
        top_k_indices = np.argsort(scores)[-query.limit:][::-1]
        
        results = []
        for idx in top_k_indices:
            if scores[idx] > 0:
                item = self.bm25.get_document(idx)
                item.relevance_score = float(scores[idx])
                results.append(item)
        
        return results
    
    async def _semantic_search(self, query: Query) -> List[ContentItem]:
        """语义向量检索"""
        query_text = " ".join(query.keywords)
        query_embedding = await self._get_embedding(query_text)
        
        # 向量数据库检索
        search_results = await self.vector_db.search(
            vector=query_embedding,
            limit=query.limit,
            filters=query.filters
        )
        
        results = []
        for result in search_results:
            item = ContentItem(**result.metadata)
            item.relevance_score = result.score
            results.append(item)
        
        return results
    
    def _reciprocal_rank_fusion(
        self,
        ranked_lists: List[tuple],
        k: int = 60
    ) -> List[ContentItem]:
        """
        Reciprocal Rank Fusion (RRF) 算法
        
        将多个排序结果融合为一个排序
        
        RRF Score = Σ(1 / (k + rank))
        """
        scores = {}
        
        for strategy, items in ranked_lists:
            for rank, item in enumerate(items, start=1):
                item_id = item.id
                rrf_score = 1.0 / (k + rank)
                
                if item_id in scores:
                    scores[item_id]["score"] += rrf_score
                else:
                    scores[item_id] = {
                        "item": item,
                        "score": rrf_score
                    }
        
        # 按融合分数排序
        sorted_items = sorted(
            scores.values(),
            key=lambda x: x["score"],
            reverse=True
        )
        
        # 更新relevance_score为融合后的分数
        results = []
        for entry in sorted_items:
            item = entry["item"]
            item.relevance_score = entry["score"]
            results.append(item)
        
        return results
    
    async def _get_embedding(self, text: str) -> List[float]:
        """获取文本的向量表示"""
        # 调用Embedding API（OpenAI/BGE-M3等）
        # 简化实现
        return [0.0] * 768  # 768维向量


class QualityAssessor:
    """内容质量评估器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def assess_batch(
        self,
        retrieval_results: List[RetrievalResult],
        objectives: List[LearningObjective]
    ) -> List[ContentItem]:
        """批量评估内容质量"""
        all_items = []
        for result in retrieval_results:
            all_items.extend(result.items)
        
        # 并发评估
        tasks = [self.assess_item(item, objectives) for item in all_items]
        assessed_items = await asyncio.gather(*tasks)
        
        return assessed_items
    
    async def assess_item(
        self,
        item: ContentItem,
        objectives: List[LearningObjective]
    ) -> ContentItem:
        """
        评估单个内容条目的质量
        
        5个维度：
        1. 权威性（Authority）：数据源可信度
        2. 相关性（Relevance）：与学习目标的匹配度
        3. 时效性（Timeliness）：内容的新鲜度
        4. 完整性（Completeness）：信息的丰富度
        5. 可读性（Readability）：语言的易读性
        """
        # 1. 权威性评分
        item.authority_score = self._assess_authority(item.source)
        
        # 2. 相关性评分（使用LLM）
        item.relevance_score = await self._assess_relevance(item, objectives)
        
        # 3. 时效性评分
        item.timeliness_score = self._assess_timeliness(item.publish_date)
        
        # 4. 完整性评分
        item.completeness_score = self._assess_completeness(item.content)
        
        # 5. 可读性评分
        item.readability_score = self._assess_readability(item.content)
        
        # 综合评分（加权平均）
        item.overall_score = (
            item.authority_score * 0.25 +
            item.relevance_score * 0.35 +
            item.timeliness_score * 0.15 +
            item.completeness_score * 0.15 +
            item.readability_score * 0.10
        )
        
        return item
    
    def _assess_authority(self, source: ContentSourceType) -> float:
        """评估权威性"""
        # 基于数据源的可信度评分
        authority_scores = {
            ContentSourceType.WIKIPEDIA: 0.9,
            ContentSourceType.WIKIART: 0.85,
            ContentSourceType.ACADEMIC_DB: 0.95,
            ContentSourceType.YOUTUBE: 0.6,
            ContentSourceType.DOUBAN: 0.7,
            ContentSourceType.GITHUB: 0.8,
            ContentSourceType.STACKOVERFLOW: 0.85,
            ContentSourceType.MDN: 0.9,
        }
        return authority_scores.get(source, 0.5)
    
    async def _assess_relevance(
        self,
        item: ContentItem,
        objectives: List[LearningObjective]
    ) -> float:
        """评估相关性（使用LLM）"""
        objectives_text = "\n".join([obj.description for obj in objectives])
        
        prompt = f"""
        请评估以下内容与学习目标的相关性（0-1分）：
        
        学习目标：
        {objectives_text}
        
        内容标题：{item.title}
        内容摘要：{item.content[:500]}
        
        请只返回一个0-1之间的数字，表示相关性分数。
        """
        
        try:
            response = await self.llm.generate(
                prompt=prompt,
                max_tokens=10,
                temperature=0.3
            )
            score = float(response.content.strip())
            return max(0.0, min(1.0, score))
        except:
            # 降级方案：关键词匹配
            return self._keyword_based_relevance(item, objectives)
    
    def _keyword_based_relevance(
        self,
        item: ContentItem,
        objectives: List[LearningObjective]
    ) -> float:
        """基于关键词的相关性评分（降级方案）"""
        # 提取目标关键词
        obj_keywords = set()
        for obj in objectives:
            words = obj.description.lower().split()
            obj_keywords.update(words)
        
        # 提取内容关键词
        content_keywords = set(item.content.lower().split())
        
        # 计算Jaccard相似度
        if not obj_keywords:
            return 0.5
        
        intersection = obj_keywords & content_keywords
        union = obj_keywords | content_keywords
        
        return len(intersection) / len(union) if union else 0.0
    
    def _assess_timeliness(self, publish_date: Optional[datetime]) -> float:
        """评估时效性"""
        if not publish_date:
            return 0.5  # 未知日期给中等分
        
        days_old = (datetime.now() - publish_date).days
        
        if days_old < 30:
            return 1.0
        elif days_old < 180:
            return 0.8
        elif days_old < 365:
            return 0.6
        elif days_old < 730:
            return 0.4
        else:
            return 0.2
    
    def _assess_completeness(self, content: str) -> float:
        """评估完整性"""
        # 基于内容长度和结构
        length = len(content)
        
        if length < 100:
            return 0.2
        elif length < 500:
            return 0.5
        elif length < 2000:
            return 0.8
        else:
            return 1.0
    
    def _assess_readability(self, content: str) -> float:
        """评估可读性"""
        # 简化实现：基于句子长度和词汇复杂度
        # 实际可使用Flesch-Kincaid等专业指标
        sentences = content.split('。')
        if not sentences:
            return 0.5
        
        avg_sentence_length = sum(len(s) for s in sentences) / len(sentences)
        
        # 理想句子长度：15-30字
        if 15 <= avg_sentence_length <= 30:
            return 1.0
        elif 10 <= avg_sentence_length <= 40:
            return 0.7
        else:
            return 0.4


class ContentIntegrator:
    """内容整合器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def integrate(
        self,
        items: List[ContentItem],
        objectives: List[LearningObjective],
        domain: DomainType
    ) -> StructuredContent:
        """
        整合多个内容条目为结构化内容
        
        步骤：
        1. 去重（相似内容合并）
        2. 提取核心要点
        3. 生成摘要
        4. 补充缺失信息
        5. 标注来源
        """
        # Step 1: 去重
        deduplicated = self._deduplicate(items)
        logger.info(f"After deduplication: {len(deduplicated)} items")
        
        # Step 2: 排序（按overall_score降序）
        sorted_items = sorted(
            deduplicated,
            key=lambda x: x.overall_score,
            reverse=True
        )
        
        # Step 3: 提取核心要点
        key_points = await self._extract_key_points(sorted_items, objectives)
        
        # Step 4: 生成摘要
        summary = await self._generate_summary(sorted_items, objectives)
        
        # Step 5: 提取详细信息
        details = self._extract_details(sorted_items, domain)
        
        # Step 6: 收集引用
        references = [item.url for item in sorted_items[:10]]  # Top 10来源
        
        # 构建结构化内容
        structured = StructuredContent(
            topic=objectives[0].description if objectives else "未知主题",
            domain=domain.value,
            summary=summary,
            key_points=key_points,
            details=details,
            sources=sorted_items[:20],  # 保留Top 20来源
            references=references,
            confidence_score=self._calculate_confidence(sorted_items)
        )
        
        return structured
    
    def _deduplicate(self, items: List[ContentItem]) -> List[ContentItem]:
        """去重：合并相似内容"""
        if not items:
            return []
        
        # 简化实现：基于标题相似度
        unique_items = []
        seen_titles = set()
        
        for item in items:
            # 标题标准化
            normalized_title = item.title.lower().strip()
            
            # 检查是否已存在相似标题
            if normalized_title not in seen_titles:
                unique_items.append(item)
                seen_titles.add(normalized_title)
        
        return unique_items
    
    async def _extract_key_points(
        self,
        items: List[ContentItem],
        objectives: List[LearningObjective]
    ) -> List[str]:
        """提取核心要点"""
        # 使用LLM从多个内容中提取关键信息
        contents = "\n\n---\n\n".join([
            f"【{item.title}】\n{item.content[:1000]}"
            for item in items[:5]  # 使用Top 5内容
        ])
        
        objectives_text = "\n".join([obj.description for obj in objectives])
        
        prompt = f"""
        请从以下内容中提取与学习目标相关的核心要点：
        
        学习目标：
        {objectives_text}
        
        内容：
        {contents}
        
        请提取5-8个核心要点，每个要点一行，用"- "开头。
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=500,
            temperature=0.5
        )
        
        # 解析要点
        key_points = [
            line.strip()[2:]  # 去掉"- "
            for line in response.content.split('\n')
            if line.strip().startswith('- ')
        ]
        
        return key_points
    
    async def _generate_summary(
        self,
        items: List[ContentItem],
        objectives: List[LearningObjective]
    ) -> str:
        """生成摘要"""
        contents = "\n\n".join([
            f"{item.title}: {item.content[:500]}"
            for item in items[:3]
        ])
        
        objectives_text = "\n".join([obj.description for obj in objectives])
        
        prompt = f"""
        请为以下内容生成一个简洁的摘要（3-5句话）：
        
        学习目标：
        {objectives_text}
        
        内容：
        {contents}
        
        摘要应该：
        1. 涵盖学习目标的核心内容
        2. 语言简洁明了
        3. 适合目标受众理解
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=300,
            temperature=0.5
        )
        
        return response.content.strip()
    
    def _extract_details(
        self,
        items: List[ContentItem],
        domain: DomainType
    ) -> Dict[str, Any]:
        """提取详细信息（领域特定）"""
        details = {}
        
        if domain == DomainType.ART_HISTORY:
            # 美术史：提取艺术家、年代、流派等
            details["artists"] = self._extract_artists(items)
            details["periods"] = self._extract_periods(items)
            details["styles"] = self._extract_styles(items)
        
        elif domain == DomainType.VOCATIONAL_TRAINING:
            # 职业培训：提取工具、技能、步骤等
            details["tools"] = self._extract_tools(items)
            details["skills"] = self._extract_skills(items)
        
        # 通用信息
        details["sources_count"] = len(items)
        details["avg_quality_score"] = sum(item.overall_score for item in items) / len(items) if items else 0
        
        return details
    
    def _extract_artists(self, items: List[ContentItem]) -> List[str]:
        """提取艺术家名字（示例）"""
        # 简化实现：使用NER或关键词匹配
        artists = set()
        keywords = ["莫奈", "梵高", "毕加索", "达芬奇"]
        
        for item in items:
            for keyword in keywords:
                if keyword in item.content:
                    artists.add(keyword)
        
        return list(artists)
    
    def _extract_periods(self, items: List[ContentItem]) -> List[str]:
        """提取时期"""
        # 简化实现
        return []
    
    def _extract_styles(self, items: List[ContentItem]) -> List[str]:
        """提取艺术风格"""
        # 简化实现
        return []
    
    def _extract_tools(self, items: List[ContentItem]) -> List[str]:
        """提取工具"""
        # 简化实现
        return []
    
    def _extract_skills(self, items: List[ContentItem]) -> List[str]:
        """提取技能"""
        # 简化实现
        return []
    
    def _calculate_confidence(self, items: List[ContentItem]) -> float:
        """计算整合内容的置信度"""
        if not items:
            return 0.0
        
        # 基于以下因素：
        # 1. 来源数量（越多越好）
        source_score = min(len(items) / 10, 1.0) * 0.3
        
        # 2. 平均质量分数
        avg_quality = sum(item.overall_score for item in items) / len(items)
        quality_score = avg_quality * 0.5
        
        # 3. 来源多样性（不同数据源）
        unique_sources = len(set(item.source for item in items))
        diversity_score = min(unique_sources / 5, 1.0) * 0.2
        
        return source_score + quality_score + diversity_score


class CacheManager:
    """缓存管理器（三级缓存）"""
    
    def __init__(self, redis_client: Redis):
        self.redis = redis_client
        
        # 缓存TTL设置
        self.session_ttl = 3600  # 1小时
        self.user_ttl = 86400  # 1天
        self.global_ttl = 604800  # 7天
    
    async def get(
        self,
        key: str,
        user_id: Optional[str] = None,
        session_id: Optional[str] = None
    ) -> Optional[StructuredContent]:
        """
        获取缓存内容
        
        查询顺序：会话缓存 → 用户缓存 → 全局缓存
        """
        # Level 1: 会话缓存
        if session_id:
            cached = await self._get_from_cache(f"session:{session_id}:{key}")
            if cached:
                logger.info(f"Cache hit: session")
                return cached
        
        # Level 2: 用户缓存
        if user_id:
            cached = await self._get_from_cache(f"user:{user_id}:{key}")
            if cached:
                logger.info(f"Cache hit: user")
                return cached
        
        # Level 3: 全局缓存
        cached = await self._get_from_cache(f"global:{key}")
        if cached:
            logger.info(f"Cache hit: global")
            return cached
        
        logger.info("Cache miss")
        return None
    
    async def set(
        self,
        key: str,
        content: StructuredContent,
        user_id: Optional[str] = None,
        session_id: Optional[str] = None
    ):
        """
        设置缓存内容
        
        同时写入多级缓存
        """
        serialized = json.dumps(content.__dict__, default=str)
        
        # Level 3: 全局缓存（总是写入）
        await self.redis.setex(
            f"global:{key}",
            self.global_ttl,
            serialized
        )
        
        # Level 2: 用户缓存
        if user_id:
            await self.redis.setex(
                f"user:{user_id}:{key}",
                self.user_ttl,
                serialized
            )
        
        # Level 1: 会话缓存
        if session_id:
            await self.redis.setex(
                f"session:{session_id}:{key}",
                self.session_ttl,
                serialized
            )
    
    async def _get_from_cache(self, cache_key: str) -> Optional[StructuredContent]:
        """从Redis获取缓存"""
        cached_data = await self.redis.get(cache_key)
        if cached_data:
            data = json.loads(cached_data)
            return StructuredContent(**data)
        return None
    
    async def invalidate(self, key: str, scope: str = "all"):
        """
        失效缓存
        
        scope: "session" | "user" | "global" | "all"
        """
        if scope == "all":
            pattern = f"*:{key}"
        else:
            pattern = f"{scope}:*:{key}"
        
        keys = await self.redis.keys(pattern)
        if keys:
            await self.redis.delete(*keys)
```

---

#### 数据源适配器示例

```python
class WikipediaRetriever(BaseRetriever):
    """Wikipedia检索器"""
    
    def __init__(self, api_key: Optional[str] = None):
        self.base_url = "https://zh.wikipedia.org/w/api.php"
        self.session = aiohttp.ClientSession()
    
    async def retrieve(self, query: Query) -> RetrievalResult:
        """从Wikipedia检索"""
        start_time = time.time()
        
        params = {
            "action": "query",
            "format": "json",
            "list": "search",
            "srsearch": " ".join(query.keywords),
            "srlimit": query.limit,
            "utf8": 1
        }
        
        async with self.session.get(self.base_url, params=params) as response:
            data = await response.json()
        
        items = []
        for result in data.get("query", {}).get("search", []):
            item = ContentItem(
                id=f"wiki_{result['pageid']}",
                source=ContentSourceType.WIKIPEDIA,
                title=result["title"],
                content=result["snippet"],
                url=f"https://zh.wikipedia.org/wiki/{result['title']}"
            )
            items.append(item)
        
        search_time = time.time() - start_time
        
        return RetrievalResult(
            query=query,
            items=items,
            total_count=len(items),
            search_time=search_time
        )


class WikiArtRetriever(BaseRetriever):
    """WikiArt检索器（爬虫）"""
    
    def __init__(self):
        self.base_url = "https://www.wikiart.org"
    
    async def retrieve(self, query: Query) -> RetrievalResult:
        """从WikiArt爬取艺术作品信息"""
        # 使用Scrapy或BeautifulSoup爬取
        # 这里是简化示例
        items = []
        
        # 搜索艺术家
        search_url = f"{self.base_url}/en/search/{'+'.join(query.keywords)}"
        
        # 爬取逻辑（简化）
        # ...
        
        return RetrievalResult(
            query=query,
            items=items,
            total_count=len(items),
            search_time=0.0
        )
```

---

#### 监控指标

```python
class ContentDiscoveryMetrics:
    """内容发现引擎监控指标"""
    
    def __init__(self):
        self.discovery_count = Counter('content_discovery_total', 'Total content discoveries')
        self.discovery_duration = Histogram('content_discovery_duration_seconds', 'Duration')
        self.sources_used = Counter('content_sources_used', 'Content sources used', ['source'])
        self.cache_hit_rate = Gauge('content_cache_hit_rate', 'Cache hit rate')
        self.quality_score = Histogram('content_quality_score', 'Quality score')
    
    def record_discovery(
        self,
        duration: float,
        sources: List[ContentSourceType],
        cache_hit: bool,
        quality: float
    ):
        """记录一次内容发现"""
        self.discovery_count.inc()
        self.discovery_duration.observe(duration)
        
        for source in sources:
            self.sources_used.labels(source=source.value).inc()
        
        if cache_hit:
            self.cache_hit_rate.inc()
        
        self.quality_score.observe(quality)
```

---

### 3.0.3 故事化叙述引擎架构设计 ⭐ 新增

#### 概述

**功能**：将枯燥的知识点转化为引人入胜的故事化内容，提升学习兴趣和记忆效果。

**核心挑战**：
- V1.0仅支持教案生成，内容单调乏味
- V2.0需要跨领域（美术史、畅销书、职业培训等）生成吸引人的叙述
- 不同领域、受众、场景需要不同的叙述策略
- 平衡趣味性和知识准确性

**解决方案**：
- 领域特定的叙述策略库（美术史用艺术家故事、职业培训用场景案例）
- 多种故事元素（人物、冲突、转折、金句）
- AI驱动的自适应叙述生成（根据受众调整语气和深度）
- 质量控制机制（事实核查、可读性检测、情感分析）

---

#### 模块架构

```
┌────────────────────────────────────────────────────────────────────┐
│      故事化叙述引擎 (Story-based Narrative Engine)                  │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────────────────────────────┐                  │
│  │ 输入：结构化内容 + 学习目标 + 受众画像        │                  │
│  │ (StructuredContent, Objectives, Audience)    │                  │
│  └─────────┬────────────────────────────────────┘                  │
│            │                                                        │
│            ▼                                                        │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块1: 叙述策略选择器 (Narrative Strategy Selector) │           │
│  │                                                     │           │
│  │  ┌───────────────────────────────────────────┐     │           │
│  │  │  策略库 (Strategy Library)                │     │           │
│  │  │  ┌──────────────────────────────────┐    │     │           │
│  │  │  │ 美术史策略 (Art History)        │    │     │           │
│  │  │  │ • 艺术家传记式                   │    │     │           │
│  │  │  │ • 作品创作背景                   │    │     │           │
│  │  │  │ • 艺术流派演进                   │    │     │           │
│  │  │  └──────────────────────────────────┘    │     │           │
│  │  │  ┌──────────────────────────────────┐    │     │           │
│  │  │  │ 畅销书策略 (Bestseller)          │    │     │           │
│  │  │  │ • 作者创作历程                   │    │     │           │
│  │  │  │ • 核心观点案例化                 │    │     │           │
│  │  │  │ • 金句提炼                       │    │     │           │
│  │  │  └──────────────────────────────────┘    │     │           │
│  │  │  ┌──────────────────────────────────┐    │     │           │
│  │  │  │ 职业培训策略 (Vocational)        │    │     │           │
│  │  │  │ • 真实工作场景                   │    │     │           │
│  │  │  │ • 案例驱动讲解                   │    │     │           │
│  │  │  │ • 步骤可视化                     │    │     │           │
│  │  │  └──────────────────────────────────┘    │     │           │
│  │  │  ┌──────────────────────────────────┐    │     │           │
│  │  │  │ K12教育策略 (K12)                │    │     │           │
│  │  │  │ • 生活化类比                     │    │     │           │
│  │  │  │ • 趣味性引入                     │    │     │           │
│  │  │  │ • 渐进式讲解                     │    │     │           │
│  │  │  └──────────────────────────────────┘    │     │           │
│  │  └───────────────────────────────────────────┘     │           │
│  │            │                                        │           │
│  │            ▼                                        │           │
│  │  ┌──────────────────────────────────┐              │           │
│  │  │ 策略匹配算法                      │              │           │
│  │  │ • 领域权重：60%                  │              │           │
│  │  │ • 受众特征：25%                  │              │           │
│  │  │ • 内容类型：15%                  │              │           │
│  │  └──────────────────────────────────┘              │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块2: 故事元素生成器 (Story Element Generator)     │           │
│  │                                                     │           │
│  │  ┌────────────────────┐  ┌────────────────────┐   │           │
│  │  │ 人物挖掘           │  │ 冲突识别           │   │           │
│  │  │ (Character)        │  │ (Conflict)         │   │           │
│  │  │ • 历史人物         │  │ • 观点对立         │   │           │
│  │  │ • 虚构角色         │  │ • 技术挑战         │   │           │
│  │  │ • 受众代入         │  │ • 历史争议         │   │           │
│  │  └────────────────────┘  └────────────────────┘   │           │
│  │                                                     │           │
│  │  ┌────────────────────┐  ┌────────────────────┐   │           │
│  │  │ 转折设计           │  │ 金句提炼           │   │           │
│  │  │ (Plot Twist)       │  │ (Quotable)         │   │           │
│  │  │ • 意外发现         │  │ • 核心观点         │   │           │
│  │  │ • 认知升级         │  │ • 记忆锚点         │   │           │
│  │  │ • 结论反转         │  │ • 情感共鸣         │   │           │
│  │  └────────────────────┘  └────────────────────┘   │           │
│  │                                                     │           │
│  │  ┌────────────────────┐  ┌────────────────────┐   │           │
│  │  │ 场景构建           │  │ 细节丰富           │   │           │
│  │  │ (Scene)            │  │ (Detail)           │   │           │
│  │  │ • 时空设定         │  │ • 五感描写         │   │           │
│  │  │ • 环境渲染         │  │ • 数据支撑         │   │           │
│  │  │ • 氛围营造         │  │ • 案例补充         │   │           │
│  │  └────────────────────┘  └────────────────────┘   │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块3: 叙述转化器 (Narrative Transformer)           │           │
│  │                                                     │           │
│  │  阶段1: 开场设计 (Opening)                          │           │
│  │  • 悬念式开场："你知道莫奈为什么要画同一个池塘40遍吗？" │           │
│  │  • 问题式开场："为什么Excel的数据透视表如此强大？"   │           │
│  │  • 场景式开场："1872年的巴黎，一位画家..."        │           │
│  │                                                     │           │
│  │  阶段2: 主体展开 (Body)                             │           │
│  │  • 故事线穿插知识点                                  │           │
│  │  • 递进式深化理解                                    │           │
│  │  • 案例与理论结合                                    │           │
│  │                                                     │           │
│  │  阶段3: 收尾升华 (Closing)                          │           │
│  │  • 核心观点总结                                      │           │
│  │  • 启发性思考                                        │           │
│  │  • 行动建议                                          │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块4: 质量控制器 (Quality Controller)              │           │
│  │                                                     │           │
│  │  ✓ 事实核查 (Fact Checking)                         │           │
│  │    • 历史事件准确性                                  │           │
│  │    • 数据引用来源                                    │           │
│  │    • 概念定义正确性                                  │           │
│  │                                                     │           │
│  │  ✓ 可读性检测 (Readability)                         │           │
│  │    • Flesch-Kincaid分数                             │           │
│  │    • 句子长度分布                                    │           │
│  │    • 词汇复杂度                                      │           │
│  │                                                     │           │
│  │  ✓ 情感分析 (Sentiment Analysis)                    │           │
│  │    • 情绪曲线波动                                    │           │
│  │    • 吸引力评分                                      │           │
│  │    • 受众匹配度                                      │           │
│  │                                                     │           │
│  │  ✓ 知识覆盖度 (Coverage)                            │           │
│  │    • 学习目标达成率                                  │           │
│  │    • 关键概念出现                                    │           │
│  │    • 深度平衡                                        │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────┐           │
│  │ 模块5: 迭代优化器 (Iterative Optimizer)             │           │
│  │                                                     │           │
│  │  🔄 自动改进循环（最多3轮）                          │           │
│  │    1. 检测质量问题                                   │           │
│  │    2. 生成改进提示                                   │           │
│  │    3. LLM重新生成                                    │           │
│  │    4. 验证改进效果                                   │           │
│  │                                                     │           │
│  │  改进维度：                                          │           │
│  │  • 趣味性不足 → 增加故事元素                        │           │
│  │  • 准确性问题 → 补充来源引用                        │           │
│  │  • 可读性差 → 简化句式                              │           │
│  │  • 覆盖度低 → 补充知识点                            │           │
│  └─────────────────┬───────────────────────────────────┘           │
│                    │                                                │
│                    ▼                                                │
│  ┌──────────────────────────────────────────┐                      │
│  │ 输出：故事化叙述内容                      │                      │
│  │ (NarrativeContent)                       │                      │
│  └──────────────────────────────────────────┘                      │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

---

#### 核心类设计

**1. 数据模型**

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from enum import Enum
from datetime import datetime

class NarrativeStrategy(Enum):
    """叙述策略"""
    ARTIST_BIOGRAPHY = "artist_biography"  # 艺术家传记式
    ARTWORK_BACKSTORY = "artwork_backstory"  # 作品背景故事
    MOVEMENT_EVOLUTION = "movement_evolution"  # 流派演进
    AUTHOR_JOURNEY = "author_journey"  # 作者历程
    CASE_DRIVEN = "case_driven"  # 案例驱动
    STEP_BY_STEP = "step_by_step"  # 步骤可视化
    ANALOGY_BASED = "analogy_based"  # 类比驱动
    PROGRESSIVE = "progressive"  # 渐进式
    SUSPENSE = "suspense"  # 悬念式
    PROBLEM_SOLVING = "problem_solving"  # 问题解决式

class StoryElementType(Enum):
    """故事元素类型"""
    CHARACTER = "character"  # 人物
    CONFLICT = "conflict"  # 冲突
    PLOT_TWIST = "plot_twist"  # 转折
    QUOTE = "quote"  # 金句
    SCENE = "scene"  # 场景
    DETAIL = "detail"  # 细节

@dataclass
class NarrativeStrategyProfile:
    """叙述策略配置"""
    strategy: NarrativeStrategy
    domain: DomainType
    description: str
    target_audience: List[str]  # 适用受众
    
    # 策略参数
    story_elements: List[StoryElementType]  # 使用的故事元素
    tone: str  # 语气（formal/casual/inspiring）
    complexity_level: int = 3  # 复杂度 1-5
    emotion_curve: str = "rising"  # 情绪曲线（flat/rising/波动）
    
    # 质量阈值
    min_readability: float = 0.6
    min_engagement: float = 0.7
    max_fact_errors: int = 0

@dataclass
class StoryElement:
    """故事元素"""
    type: StoryElementType
    content: str  # 元素内容
    source: Optional[str] = None  # 来源
    confidence: float = 1.0  # 置信度
    metadata: Dict[str, Any] = field(default_factory=dict)

@dataclass
class NarrativeSection:
    """叙述章节"""
    title: str
    content: str
    story_elements: List[StoryElement] = field(default_factory=list)
    knowledge_points: List[str] = field(default_factory=list)  # 覆盖的知识点
    word_count: int = 0

@dataclass
class NarrativeContent:
    """故事化叙述内容"""
    # 基本信息
    title: str
    topic: str
    domain: DomainType
    strategy: NarrativeStrategy
    
    # 结构化内容
    opening: NarrativeSection  # 开场
    body: List[NarrativeSection]  # 主体（多个章节）
    closing: NarrativeSection  # 结尾
    
    # 附加信息
    key_quotes: List[str] = field(default_factory=list)  # 金句
    references: List[str] = field(default_factory=list)  # 参考资料
    
    # 质量指标
    readability_score: float = 0.0  # 可读性 0-1
    engagement_score: float = 0.0  # 吸引力 0-1
    accuracy_score: float = 0.0  # 准确性 0-1
    coverage_score: float = 0.0  # 覆盖度 0-1
    overall_quality: float = 0.0  # 综合质量 0-1
    
    # 元数据
    word_count: int = 0
    estimated_reading_time: int = 0  # 预估阅读时间（分钟）
    created_at: datetime = field(default_factory=datetime.now)
    revision_count: int = 0  # 迭代次数

@dataclass
class QualityIssue:
    """质量问题"""
    dimension: str  # 维度（readability/accuracy/engagement/coverage）
    severity: str  # 严重程度（low/medium/high）
    description: str  # 问题描述
    suggestion: str  # 改进建议
    affected_section: Optional[str] = None  # 受影响的章节
```

**2. 核心类实现**

```python
class StoryBasedNarrativeEngine:
    """故事化叙述引擎"""
    
    def __init__(
        self,
        strategy_selector: NarrativeStrategySelector,
        element_generator: StoryElementGenerator,
        transformer: NarrativeTransformer,
        quality_controller: QualityController,
        optimizer: IterativeOptimizer,
        llm_client: LLMClient,
        config: NarrativeEngineConfig
    ):
        self.strategy_selector = strategy_selector
        self.element_generator = element_generator
        self.transformer = transformer
        self.quality_controller = quality_controller
        self.optimizer = optimizer
        self.llm = llm_client
        self.config = config
    
    async def generate_narrative(
        self,
        structured_content: StructuredContent,
        objectives: List[LearningObjective],
        audience: AudienceProfile,
        domain: DomainType
    ) -> NarrativeContent:
        """
        主流程：生成故事化叙述内容
        
        Args:
            structured_content: AI内容发现引擎输出的结构化内容
            objectives: 学习目标
            audience: 受众画像
            domain: 领域类型
            
        Returns:
            NarrativeContent: 故事化叙述内容
        """
        logger.info(f"Starting narrative generation for topic: {structured_content.topic}")
        
        # Step 1: 选择叙述策略
        strategy_profile = await self.strategy_selector.select_strategy(
            domain=domain,
            audience=audience,
            content_type=structured_content.metadata.get("content_type", "general")
        )
        logger.info(f"Selected strategy: {strategy_profile.strategy.value}")
        
        # Step 2: 生成故事元素
        story_elements = await self.element_generator.generate_elements(
            content=structured_content,
            strategy=strategy_profile,
            objectives=objectives
        )
        logger.info(f"Generated {len(story_elements)} story elements")
        
        # Step 3: 转化为叙述内容（初始版本）
        narrative = await self.transformer.transform(
            content=structured_content,
            elements=story_elements,
            strategy=strategy_profile,
            objectives=objectives
        )
        logger.info(f"Generated initial narrative: {narrative.word_count} words")
        
        # Step 4: 质量控制
        quality_issues = await self.quality_controller.evaluate(
            narrative=narrative,
            objectives=objectives,
            strategy=strategy_profile
        )
        logger.info(f"Quality check: {len(quality_issues)} issues found")
        
        # Step 5: 迭代优化（最多3轮）
        if quality_issues and self.config.enable_optimization:
            narrative = await self.optimizer.optimize(
                narrative=narrative,
                issues=quality_issues,
                max_iterations=self.config.max_iterations
            )
            logger.info(f"Optimization completed: {narrative.revision_count} iterations")
        
        # Step 6: 最终质量评分
        narrative = await self._calculate_final_quality(narrative, objectives)
        
        logger.info(f"Narrative generation completed. Quality: {narrative.overall_quality:.2f}")
        
        return narrative
    
    async def _calculate_final_quality(
        self,
        narrative: NarrativeContent,
        objectives: List[LearningObjective]
    ) -> NarrativeContent:
        """计算最终质量分数"""
        narrative.overall_quality = (
            narrative.readability_score * 0.2 +
            narrative.engagement_score * 0.3 +
            narrative.accuracy_score * 0.3 +
            narrative.coverage_score * 0.2
        )
        return narrative


class NarrativeStrategySelector:
    """叙述策略选择器"""
    
    def __init__(self):
        # 预定义策略库
        self.strategy_library = self._build_strategy_library()
    
    def _build_strategy_library(self) -> Dict[DomainType, List[NarrativeStrategyProfile]]:
        """构建策略库"""
        library = {
            DomainType.ART_HISTORY: [
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.ARTIST_BIOGRAPHY,
                    domain=DomainType.ART_HISTORY,
                    description="以艺术家的生平经历为线索讲述艺术作品",
                    target_audience=["art_enthusiast", "student"],
                    story_elements=[
                        StoryElementType.CHARACTER,
                        StoryElementType.SCENE,
                        StoryElementType.PLOT_TWIST,
                        StoryElementType.QUOTE
                    ],
                    tone="inspiring",
                    complexity_level=3,
                    emotion_curve="rising"
                ),
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.ARTWORK_BACKSTORY,
                    domain=DomainType.ART_HISTORY,
                    description="讲述艺术作品创作的背景故事",
                    target_audience=["general_public"],
                    story_elements=[
                        StoryElementType.SCENE,
                        StoryElementType.DETAIL,
                        StoryElementType.CONFLICT
                    ],
                    tone="casual",
                    complexity_level=2
                ),
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.MOVEMENT_EVOLUTION,
                    domain=DomainType.ART_HISTORY,
                    description="讲述艺术流派的演进历程",
                    target_audience=["researcher", "student"],
                    story_elements=[
                        StoryElementType.CHARACTER,
                        StoryElementType.CONFLICT,
                        StoryElementType.PLOT_TWIST
                    ],
                    tone="formal",
                    complexity_level=4
                )
            ],
            DomainType.BESTSELLER_INTERPRETATION: [
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.AUTHOR_JOURNEY,
                    domain=DomainType.BESTSELLER_INTERPRETATION,
                    description="讲述作者的写作历程和思想演变",
                    target_audience=["reader", "writer"],
                    story_elements=[
                        StoryElementType.CHARACTER,
                        StoryElementType.QUOTE,
                        StoryElementType.DETAIL
                    ],
                    tone="inspiring"
                ),
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.CASE_DRIVEN,
                    domain=DomainType.BESTSELLER_INTERPRETATION,
                    description="用书中案例串联核心观点",
                    target_audience=["general_public"],
                    story_elements=[
                        StoryElementType.SCENE,
                        StoryElementType.QUOTE,
                        StoryElementType.DETAIL
                    ],
                    tone="casual"
                )
            ],
            DomainType.VOCATIONAL_TRAINING: [
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.CASE_DRIVEN,
                    domain=DomainType.VOCATIONAL_TRAINING,
                    description="用真实工作场景案例讲解技能",
                    target_audience=["professional", "student"],
                    story_elements=[
                        StoryElementType.SCENE,
                        StoryElementType.DETAIL,
                        StoryElementType.PROBLEM_SOLVING
                    ],
                    tone="casual",
                    complexity_level=3
                ),
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.STEP_BY_STEP,
                    domain=DomainType.VOCATIONAL_TRAINING,
                    description="步骤化可视化讲解",
                    target_audience=["beginner"],
                    story_elements=[
                        StoryElementType.DETAIL,
                        StoryElementType.SCENE
                    ],
                    tone="formal",
                    complexity_level=2
                )
            ],
            DomainType.K12: [
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.ANALOGY_BASED,
                    domain=DomainType.K12,
                    description="用生活化类比讲解抽象概念",
                    target_audience=["student"],
                    story_elements=[
                        StoryElementType.SCENE,
                        StoryElementType.DETAIL
                    ],
                    tone="casual",
                    complexity_level=2
                ),
                NarrativeStrategyProfile(
                    strategy=NarrativeStrategy.SUSPENSE,
                    domain=DomainType.K12,
                    description="用悬念式开场吸引注意力",
                    target_audience=["student"],
                    story_elements=[
                        StoryElementType.PLOT_TWIST,
                        StoryElementType.SCENE
                    ],
                    tone="casual",
                    emotion_curve="波动"
                )
            ]
        }
        return library
    
    async def select_strategy(
        self,
        domain: DomainType,
        audience: AudienceProfile,
        content_type: str = "general"
    ) -> NarrativeStrategyProfile:
        """
        选择最合适的叙述策略
        
        匹配算法：
        - 领域权重：60%
        - 受众特征：25%
        - 内容类型：15%
        """
        candidates = self.strategy_library.get(domain, [])
        
        if not candidates:
            # 降级：使用通用策略
            return self._get_default_strategy(domain)
        
        # 为每个候选策略打分
        scored_strategies = []
        for strategy in candidates:
            score = 0.0
            
            # 1. 领域匹配（60分）
            if strategy.domain == domain:
                score += 60.0
            
            # 2. 受众匹配（25分）
            audience_match = self._calculate_audience_match(
                strategy.target_audience,
                audience
            )
            score += audience_match * 25.0
            
            # 3. 内容类型匹配（15分）
            # 简化实现
            score += 15.0
            
            scored_strategies.append((strategy, score))
        
        # 返回得分最高的策略
        best_strategy = max(scored_strategies, key=lambda x: x[1])[0]
        logger.info(f"Selected strategy: {best_strategy.strategy.value} (score: {max(scored_strategies, key=lambda x: x[1])[1]:.1f})")
        
        return best_strategy
    
    def _calculate_audience_match(
        self,
        target_audience: List[str],
        audience: AudienceProfile
    ) -> float:
        """计算受众匹配度"""
        # 简化实现：检查受众标签是否在目标受众中
        audience_tags = set([audience.role.lower(), audience.education_level.lower()])
        target_tags = set(target_audience)
        
        intersection = audience_tags & target_tags
        if not target_tags:
            return 0.5
        
        return len(intersection) / len(target_tags)
    
    def _get_default_strategy(self, domain: DomainType) -> NarrativeStrategyProfile:
        """获取默认策略"""
        return NarrativeStrategyProfile(
            strategy=NarrativeStrategy.PROGRESSIVE,
            domain=domain,
            description="渐进式讲解（默认策略）",
            target_audience=["general"],
            story_elements=[StoryElementType.SCENE, StoryElementType.DETAIL],
            tone="casual"
        )


class StoryElementGenerator:
    """故事元素生成器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def generate_elements(
        self,
        content: StructuredContent,
        strategy: NarrativeStrategyProfile,
        objectives: List[LearningObjective]
    ) -> List[StoryElement]:
        """
        生成故事元素
        
        根据策略要求的元素类型生成对应内容
        """
        elements = []
        
        for element_type in strategy.story_elements:
            if element_type == StoryElementType.CHARACTER:
                chars = await self._generate_characters(content, strategy)
                elements.extend(chars)
            
            elif element_type == StoryElementType.CONFLICT:
                conflicts = await self._generate_conflicts(content, objectives)
                elements.extend(conflicts)
            
            elif element_type == StoryElementType.PLOT_TWIST:
                twists = await self._generate_plot_twists(content)
                elements.extend(twists)
            
            elif element_type == StoryElementType.QUOTE:
                quotes = await self._generate_quotes(content)
                elements.extend(quotes)
            
            elif element_type == StoryElementType.SCENE:
                scenes = await self._generate_scenes(content, strategy)
                elements.extend(scenes)
            
            elif element_type == StoryElementType.DETAIL:
                details = await self._generate_details(content)
                elements.extend(details)
        
        logger.info(f"Generated {len(elements)} story elements")
        return elements
    
    async def _generate_characters(
        self,
        content: StructuredContent,
        strategy: NarrativeStrategyProfile
    ) -> List[StoryElement]:
        """生成人物元素"""
        prompt = f"""
        基于以下内容，提取或创造适合的人物角色：
        
        主题：{content.topic}
        领域：{content.domain}
        关键信息：{content.summary}
        
        叙述策略：{strategy.description}
        
        请列出2-3个核心人物，每个人物包括：
        1. 姓名/角色
        2. 简短描述（1-2句话）
        3. 与主题的关联
        
        格式：
        - 人物1: [姓名] - [描述]
        - 人物2: [姓名] - [描述]
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=300,
            temperature=0.7
        )
        
        # 解析人物
        characters = []
        for line in response.content.split('\n'):
            if line.strip().startswith('- '):
                char_text = line.strip()[2:]
                characters.append(StoryElement(
                    type=StoryElementType.CHARACTER,
                    content=char_text,
                    confidence=0.8
                ))
        
        return characters
    
    async def _generate_conflicts(
        self,
        content: StructuredContent,
        objectives: List[LearningObjective]
    ) -> List[StoryElement]:
        """生成冲突元素"""
        prompt = f"""
        基于以下学习目标和内容，识别可以用于叙述的冲突点：
        
        学习目标：
        {chr(10).join([obj.description for obj in objectives[:3]])}
        
        内容摘要：{content.summary}
        
        请列出1-2个冲突点，例如：
        - 观点对立（不同学派的争论）
        - 技术挑战（如何解决某个难题）
        - 历史争议（对某个事件的不同解读）
        
        格式：
        - 冲突1: [描述]
        - 冲突2: [描述]
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=200,
            temperature=0.7
        )
        
        conflicts = []
        for line in response.content.split('\n'):
            if line.strip().startswith('- '):
                conflict_text = line.strip()[2:]
                conflicts.append(StoryElement(
                    type=StoryElementType.CONFLICT,
                    content=conflict_text,
                    confidence=0.7
                ))
        
        return conflicts
    
    async def _generate_plot_twists(
        self,
        content: StructuredContent
    ) -> List[StoryElement]:
        """生成转折元素"""
        # 简化实现：从关键要点中提取意外发现
        twists = []
        
        for point in content.key_points[:2]:
            if any(keyword in point for keyword in ["然而", "意外", "发现", "实际上"]):
                twists.append(StoryElement(
                    type=StoryElementType.PLOT_TWIST,
                    content=point,
                    confidence=0.6
                ))
        
        return twists
    
    async def _generate_quotes(
        self,
        content: StructuredContent
    ) -> List[StoryElement]:
        """生成金句元素"""
        prompt = f"""
        基于以下内容，提炼2-3句简洁有力的金句（适合引用和记忆）：
        
        主题：{content.topic}
        核心要点：
        {chr(10).join(content.key_points[:5])}
        
        金句要求：
        - 简洁（10-30字）
        - 易记
        - 体现核心观点
        
        格式：
        1. [金句1]
        2. [金句2]
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=200,
            temperature=0.7
        )
        
        quotes = []
        for line in response.content.split('\n'):
            if line.strip() and line.strip()[0].isdigit():
                quote_text = line.strip().split('. ', 1)[1] if '. ' in line else line.strip()
                quotes.append(StoryElement(
                    type=StoryElementType.QUOTE,
                    content=quote_text,
                    confidence=0.9
                ))
        
        return quotes
    
    async def _generate_scenes(
        self,
        content: StructuredContent,
        strategy: NarrativeStrategyProfile
    ) -> List[StoryElement]:
        """生成场景元素"""
        # 简化实现：基于内容来源创建场景
        scenes = []
        
        if content.sources:
            # 使用第一个来源创建场景
            source = content.sources[0]
            scene_text = f"场景：{source.title}的相关背景"
            scenes.append(StoryElement(
                type=StoryElementType.SCENE,
                content=scene_text,
                source=source.url,
                confidence=0.7
            ))
        
        return scenes
    
    async def _generate_details(
        self,
        content: StructuredContent
    ) -> List[StoryElement]:
        """生成细节元素"""
        # 从metadata中提取有趣的细节
        details = []
        
        if "details" in content.metadata:
            for key, value in content.metadata["details"].items():
                if value:
                    detail_text = f"{key}: {value}"
                    details.append(StoryElement(
                        type=StoryElementType.DETAIL,
                        content=detail_text,
                        confidence=0.8
                    ))
        
        return details[:3]  # 限制数量


class NarrativeTransformer:
    """叙述转化器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def transform(
        self,
        content: StructuredContent,
        elements: List[StoryElement],
        strategy: NarrativeStrategyProfile,
        objectives: List[LearningObjective]
    ) -> NarrativeContent:
        """
        将结构化内容和故事元素转化为完整叙述
        
        三阶段：开场 → 主体 → 结尾
        """
        # 阶段1: 生成开场
        opening = await self._generate_opening(
            content, elements, strategy
        )
        
        # 阶段2: 生成主体（分多个章节）
        body_sections = await self._generate_body(
            content, elements, strategy, objectives
        )
        
        # 阶段3: 生成结尾
        closing = await self._generate_closing(
            content, objectives, strategy
        )
        
        # 提取金句
        key_quotes = [
            e.content for e in elements
            if e.type == StoryElementType.QUOTE
        ]
        
        # 构建完整叙述
        narrative = NarrativeContent(
            title=content.topic,
            topic=content.topic,
            domain=content.domain,
            strategy=strategy.strategy,
            opening=opening,
            body=body_sections,
            closing=closing,
            key_quotes=key_quotes,
            references=content.references,
            word_count=self._calculate_word_count(opening, body_sections, closing)
        )
        
        # 估算阅读时间（平均250字/分钟）
        narrative.estimated_reading_time = max(1, narrative.word_count // 250)
        
        return narrative
    
    async def _generate_opening(
        self,
        content: StructuredContent,
        elements: List[StoryElement],
        strategy: NarrativeStrategyProfile
    ) -> NarrativeSection:
        """生成开场"""
        # 根据策略选择开场方式
        if strategy.emotion_curve == "波动":
            opening_style = "悬念式"
            instruction = "用一个引人入胜的问题或意外事实开场"
        elif strategy.tone == "inspiring":
            opening_style = "场景式"
            instruction = "用生动的场景描写开场，营造氛围"
        else:
            opening_style = "问题式"
            instruction = "用一个引发思考的问题开场"
        
        # 提取相关元素
        characters = [e for e in elements if e.type == StoryElementType.CHARACTER]
        scenes = [e for e in elements if e.type == StoryElementType.SCENE]
        
        prompt = f"""
        为以下内容创作一个吸引人的开场段落（150-250字）：
        
        主题：{content.topic}
        摘要：{content.summary}
        
        开场风格：{opening_style}
        指导：{instruction}
        
        可用元素：
        {chr(10).join([f"- {e.content}" for e in (characters + scenes)[:3]])}
        
        要求：
        1. 立即抓住读者注意力
        2. 自然引出主题
        3. 设定叙述基调
        4. 字数控制在150-250字
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=400,
            temperature=0.8
        )
        
        opening_text = response.content.strip()
        
        return NarrativeSection(
            title="开场",
            content=opening_text,
            word_count=len(opening_text)
        )
    
    async def _generate_body(
        self,
        content: StructuredContent,
        elements: List[StoryElement],
        strategy: NarrativeStrategyProfile,
        objectives: List[LearningObjective]
    ) -> List[NarrativeSection]:
        """生成主体（多个章节）"""
        sections = []
        
        # 根据学习目标数量决定章节数
        num_sections = min(len(objectives), 5)  # 最多5个章节
        
        for i, objective in enumerate(objectives[:num_sections]):
            # 为每个学习目标创建一个章节
            section = await self._generate_body_section(
                objective=objective,
                content=content,
                elements=elements,
                strategy=strategy,
                section_index=i
            )
            sections.append(section)
        
        return sections
    
    async def _generate_body_section(
        self,
        objective: LearningObjective,
        content: StructuredContent,
        elements: List[StoryElement],
        strategy: NarrativeStrategyProfile,
        section_index: int
    ) -> NarrativeSection:
        """生成主体的一个章节"""
        # 提取相关知识点
        relevant_points = [
            point for point in content.key_points
            if any(keyword in point for keyword in objective.description.split())
        ][:2]
        
        # 提取相关元素
        relevant_elements = elements[section_index*2:(section_index+1)*2]
        
        prompt = f"""
        为以下学习目标创作一个章节（300-500字）：
        
        学习目标：{objective.description}
        认知层级：{objective.bloom_level}
        
        核心信息：
        {chr(10).join(relevant_points)}
        
        故事元素：
        {chr(10).join([f"- {e.content}" for e in relevant_elements])}
        
        叙述策略：{strategy.description}
        语气：{strategy.tone}
        
        要求：
        1. 将知识点融入故事叙述
        2. 保持叙述连贯性
        3. 适当使用故事元素增强趣味
        4. 确保知识准确性
        5. 字数300-500字
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=800,
            temperature=0.7
        )
        
        section_content = response.content.strip()
        
        return NarrativeSection(
            title=f"第{section_index+1}部分：{objective.description[:20]}",
            content=section_content,
            knowledge_points=[objective.description],
            word_count=len(section_content)
        )
    
    async def _generate_closing(
        self,
        content: StructuredContent,
        objectives: List[LearningObjective],
        strategy: NarrativeStrategyProfile
    ) -> NarrativeSection:
        """生成结尾"""
        prompt = f"""
        为以下内容创作一个升华性的结尾（150-250字）：
        
        主题：{content.topic}
        核心观点：
        {chr(10).join(content.key_points[:3])}
        
        学习目标：
        {chr(10).join([obj.description for obj in objectives[:3]])}
        
        要求：
        1. 总结核心观点
        2. 提供启发性思考
        3. 给出行动建议（如适用）
        4. 呼应开场
        5. 字数150-250字
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=400,
            temperature=0.7
        )
        
        closing_text = response.content.strip()
        
        return NarrativeSection(
            title="结尾",
            content=closing_text,
            word_count=len(closing_text)
        )
    
    def _calculate_word_count(
        self,
        opening: NarrativeSection,
        body: List[NarrativeSection],
        closing: NarrativeSection
    ) -> int:
        """计算总字数"""
        return (
            opening.word_count +
            sum(s.word_count for s in body) +
            closing.word_count
        )


class QualityController:
    """质量控制器"""
    
    def __init__(self, llm_client: LLMClient):
        self.llm = llm_client
    
    async def evaluate(
        self,
        narrative: NarrativeContent,
        objectives: List[LearningObjective],
        strategy: NarrativeStrategyProfile
    ) -> List[QualityIssue]:
        """
        评估叙述内容质量
        
        4个维度：事实准确性、可读性、吸引力、知识覆盖度
        """
        issues = []
        
        # 1. 事实核查
        accuracy_issues = await self._check_facts(narrative)
        issues.extend(accuracy_issues)
        narrative.accuracy_score = 1.0 - (len(accuracy_issues) * 0.2)
        
        # 2. 可读性检测
        readability_issues = self._check_readability(narrative, strategy)
        issues.extend(readability_issues)
        narrative.readability_score = max(0, 1.0 - len(readability_issues) * 0.15)
        
        # 3. 吸引力评估
        engagement_issues = await self._check_engagement(narrative, strategy)
        issues.extend(engagement_issues)
        narrative.engagement_score = max(0, 1.0 - len(engagement_issues) * 0.15)
        
        # 4. 知识覆盖度
        coverage_issues = self._check_coverage(narrative, objectives)
        issues.extend(coverage_issues)
        narrative.coverage_score = max(0, 1.0 - len(coverage_issues) * 0.2)
        
        return issues
    
    async def _check_facts(self, narrative: NarrativeContent) -> List[QualityIssue]:
        """事实核查"""
        # 使用LLM检查明显的事实错误
        issues = []
        
        # 检查引用来源
        if not narrative.references:
            issues.append(QualityIssue(
                dimension="accuracy",
                severity="medium",
                description="缺少参考资料来源",
                suggestion="添加内容来源引用"
            ))
        
        # TODO: 更复杂的事实核查（对比知识库、时间线一致性等）
        
        return issues
    
    def _check_readability(
        self,
        narrative: NarrativeContent,
        strategy: NarrativeStrategyProfile
    ) -> List[QualityIssue]:
        """可读性检测"""
        issues = []
        
        # 检查句子长度
        all_text = narrative.opening.content + " ".join([s.content for s in narrative.body])
        sentences = all_text.split('。')
        
        avg_sentence_length = sum(len(s) for s in sentences) / len(sentences) if sentences else 0
        
        # 理想句子长度：15-35字
        if avg_sentence_length > 50:
            issues.append(QualityIssue(
                dimension="readability",
                severity="medium",
                description=f"句子平均长度过长（{avg_sentence_length:.0f}字）",
                suggestion="将长句拆分为更短的句子"
            ))
        elif avg_sentence_length < 10:
            issues.append(QualityIssue(
                dimension="readability",
                severity="low",
                description=f"句子平均长度过短（{avg_sentence_length:.0f}字）",
                suggestion="适当增加句子长度，增强表达"
            ))
        
        # 检查段落长度
        for i, section in enumerate(narrative.body):
            if section.word_count > 800:
                issues.append(QualityIssue(
                    dimension="readability",
                    severity="medium",
                    description=f"第{i+1}章节过长（{section.word_count}字）",
                    suggestion="将章节拆分为更小的段落",
                    affected_section=section.title
                ))
        
        return issues
    
    async def _check_engagement(
        self,
        narrative: NarrativeContent,
        strategy: NarrativeStrategyProfile
    ) -> List[QualityIssue]:
        """吸引力评估"""
        issues = []
        
        # 检查开场是否吸引人
        if len(narrative.opening.content) < 100:
            issues.append(QualityIssue(
                dimension="engagement",
                severity="high",
                description="开场过短，难以吸引读者",
                suggestion="扩展开场，增加引人入胜的元素"
            ))
        
        # 检查故事元素使用
        element_count = sum(len(s.story_elements) for s in narrative.body)
        if element_count == 0:
            issues.append(QualityIssue(
                dimension="engagement",
                severity="medium",
                description="缺少故事元素",
                suggestion="增加人物、场景、冲突等故事元素"
            ))
        
        # 检查金句
        if not narrative.key_quotes:
            issues.append(QualityIssue(
                dimension="engagement",
                severity="low",
                description="缺少令人印象深刻的金句",
                suggestion="提炼1-3句核心观点作为金句"
            ))
        
        return issues
    
    def _check_coverage(
        self,
        narrative: NarrativeContent,
        objectives: List[LearningObjective]
    ) -> List[QualityIssue]:
        """知识覆盖度检查"""
        issues = []
        
        # 检查每个学习目标是否被覆盖
        for obj in objectives:
            is_covered = False
            
            # 检查主体章节
            for section in narrative.body:
                if obj.description in section.knowledge_points or \
                   any(keyword in section.content for keyword in obj.description.split()[:3]):
                    is_covered = True
                    break
            
            if not is_covered:
                issues.append(QualityIssue(
                    dimension="coverage",
                    severity="high",
                    description=f"学习目标未覆盖：{obj.description}",
                    suggestion="添加相关内容章节"
                ))
        
        return issues


class IterativeOptimizer:
    """迭代优化器"""
    
    def __init__(
        self,
        transformer: NarrativeTransformer,
        quality_controller: QualityController,
        llm_client: LLMClient
    ):
        self.transformer = transformer
        self.quality_controller = quality_controller
        self.llm = llm_client
    
    async def optimize(
        self,
        narrative: NarrativeContent,
        issues: List[QualityIssue],
        max_iterations: int = 3
    ) -> NarrativeContent:
        """
        迭代优化叙述内容
        
        最多迭代3轮，每轮解决最严重的问题
        """
        current_narrative = narrative
        current_issues = issues
        iteration = 0
        
        while iteration < max_iterations and current_issues:
            iteration += 1
            logger.info(f"Optimization iteration {iteration}: {len(current_issues)} issues to fix")
            
            # 按严重程度排序问题
            sorted_issues = sorted(
                current_issues,
                key=lambda x: {"high": 3, "medium": 2, "low": 1}[x.severity],
                reverse=True
            )
            
            # 处理最严重的问题（最多3个）
            issues_to_fix = sorted_issues[:3]
            
            # 生成改进提示
            improvement_prompt = self._generate_improvement_prompt(
                current_narrative,
                issues_to_fix
            )
            
            # 重新生成受影响的部分
            current_narrative = await self._regenerate_content(
                current_narrative,
                improvement_prompt,
                issues_to_fix
            )
            
            # 重新评估
            # 注意：这里简化了，实际应该调用完整的评估流程
            current_issues = [
                issue for issue in current_issues
                if issue not in issues_to_fix
            ]
            
            current_narrative.revision_count = iteration
            
            logger.info(f"Iteration {iteration} completed. Remaining issues: {len(current_issues)}")
        
        return current_narrative
    
    def _generate_improvement_prompt(
        self,
        narrative: NarrativeContent,
        issues: List[QualityIssue]
    ) -> str:
        """生成改进提示"""
        prompt_parts = ["请改进以下内容，解决这些问题：\n"]
        
        for i, issue in enumerate(issues, 1):
            prompt_parts.append(f"{i}. [{issue.dimension}] {issue.description}")
            prompt_parts.append(f"   建议：{issue.suggestion}\n")
        
        return "\n".join(prompt_parts)
    
    async def _regenerate_content(
        self,
        narrative: NarrativeContent,
        improvement_prompt: str,
        issues: List[QualityIssue]
    ) -> NarrativeContent:
        """重新生成内容"""
        # 简化实现：只重新生成受影响的章节
        for issue in issues:
            if issue.affected_section:
                # 找到对应章节
                for i, section in enumerate(narrative.body):
                    if section.title == issue.affected_section:
                        # 重新生成该章节
                        new_content = await self._regenerate_section(
                            section,
                            improvement_prompt
                        )
                        narrative.body[i].content = new_content
                        narrative.body[i].word_count = len(new_content)
        
        # 重新计算总字数
        narrative.word_count = (
            narrative.opening.word_count +
            sum(s.word_count for s in narrative.body) +
            narrative.closing.word_count
        )
        
        return narrative
    
    async def _regenerate_section(
        self,
        section: NarrativeSection,
        improvement_prompt: str
    ) -> str:
        """重新生成章节"""
        prompt = f"""
        {improvement_prompt}
        
        原始内容：
        {section.content}
        
        请重新生成改进后的内容，保持原有结构和知识点，但解决上述问题。
        """
        
        response = await self.llm.generate(
            prompt=prompt,
            max_tokens=800,
            temperature=0.7
        )
        
        return response.content.strip()
```

---

#### 监控指标

```python
class NarrativeEngineMetrics:
    """叙述引擎监控指标"""
    
    def __init__(self):
        self.generation_count = Counter('narrative_generation_total', 'Total generations')
        self.generation_duration = Histogram('narrative_generation_duration_seconds', 'Duration')
        self.strategy_usage = Counter('narrative_strategy_usage', 'Strategy usage', ['strategy'])
        self.quality_scores = Histogram('narrative_quality_score', 'Quality score', ['dimension'])
        self.revision_count = Histogram('narrative_revision_count', 'Revision iterations')
    
    def record_generation(
        self,
        duration: float,
        strategy: NarrativeStrategy,
        quality: NarrativeContent,
        revision_count: int
    ):
        """记录一次叙述生成"""
        self.generation_count.inc()
        self.generation_duration.observe(duration)
        self.strategy_usage.labels(strategy=strategy.value).inc()
        
        self.quality_scores.labels(dimension='readability').observe(quality.readability_score)
        self.quality_scores.labels(dimension='engagement').observe(quality.engagement_score)
        self.quality_scores.labels(dimension='accuracy').observe(quality.accuracy_score)
        self.quality_scores.labels(dimension='coverage').observe(quality.coverage_score)
        
        self.revision_count.observe(revision_count)
```

---

### 3.0.4 领域适配器架构设计 ⭐ 新增

#### 概述

**功能**：将通用的元工作流引擎适配到具体领域，实现"孵化器-鸡"的解耦。

**核心挑战**：
- V1.0将K12逻辑硬编码在引擎中，扩展新领域需改动核心代码
- V2.0需要支持4+领域（K12、美术史、畅销书、职业培训等）
- 不同领域的数据结构、业务逻辑、内容模板差异巨大
- 需要支持第三方开发者贡献新领域适配器

**解决方案**：
- 插件化架构：适配器作为独立插件，热插拔
- 标准化接口：BaseDomainAdapter定义统一契约
- 自动发现机制：注册中心自动扫描和注册适配器
- 智能路由：根据用户请求自动选择合适的适配器
- 版本管理：支持适配器的安装、更新、卸载

**"孵化器-鸡-蛋"模型中的定位**：
- **孵化器（Incubator）**：元工作流引擎（领域无关）
- **鸡（Chicken）**：领域适配器（本模块）
- **蛋（Egg）**：生成的内容实例

---

#### 模块架构

```
┌────────────────────────────────────────────────────────────────────┐
│           领域适配器层 (Domain Adapter Layer)                       │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────────────────────────────────────────┐      │
│  │  用户请求 (User Request)                                 │      │
│  │  {                                                       │      │
│  │    "domain": "art_history",                             │      │
│  │    "input": "莫奈的印象派作品...",                       │      │
│  │    "user_profile": {...}                                │      │
│  │  }                                                       │      │
│  └─────────────────────┬────────────────────────────────────┘      │
│                        │                                            │
│                        ▼                                            │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │ 模块1: 适配器注册中心 (Adapter Registry)                │       │
│  │                                                         │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │  已注册适配器列表 (Registered Adapters)          │  │       │
│  │  │  ┌────────────────────────────────────────────┐ │  │       │
│  │  │  │ K12Adapter                                 │ │  │       │
│  │  │  │ • domain: "k12"                           │ │  │       │
│  │  │  │ • version: "2.0.1"                        │ │  │       │
│  │  │  │ • status: "active"                        │ │  │       │
│  │  │  └────────────────────────────────────────────┘ │  │       │
│  │  │  ┌────────────────────────────────────────────┐ │  │       │
│  │  │  │ ArtHistoryAdapter                          │ │  │       │
│  │  │  │ • domain: "art_history"                   │ │  │       │
│  │  │  │ • version: "1.0.0"                        │ │  │       │
│  │  │  │ • status: "active"                        │ │  │       │
│  │  │  └────────────────────────────────────────────┘ │  │       │
│  │  │  ┌────────────────────────────────────────────┐ │  │       │
│  │  │  │ BestsellerAdapter                          │ │  │       │
│  │  │  │ • domain: "bestseller"                    │ │  │       │
│  │  │  │ • version: "1.0.0"                        │ │  │       │
│  │  │  │ • status: "active"                        │ │  │       │
│  │  │  └────────────────────────────────────────────┘ │  │       │
│  │  │  ┌────────────────────────────────────────────┐ │  │       │
│  │  │  │ VocationalAdapter                          │ │  │       │
│  │  │  │ • domain: "vocational_training"           │ │  │       │
│  │  │  │ • version: "1.0.0"                        │ │  │       │
│  │  │  │ • status: "active"                        │ │  │       │
│  │  │  └────────────────────────────────────────────┘ │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │                                                         │       │
│  │  功能：                                                  │       │
│  │  • 自动发现和注册适配器                                  │       │
│  │  • 管理适配器生命周期（安装/启用/禁用/卸载）              │       │
│  │  • 版本控制和依赖管理                                    │       │
│  │  • 健康检查和故障隔离                                    │       │
│  └─────────────────┬───────────────────────────────────────┘       │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │ 模块2: 自动路由器 (Auto Router)                         │       │
│  │                                                         │       │
│  │  路由策略：                                              │       │
│  │  1️⃣ 显式指定：request.domain == "art_history"          │       │
│  │     → 直接路由到 ArtHistoryAdapter                      │       │
│  │                                                         │       │
│  │  2️⃣ 智能推断（用户未指定领域）：                        │       │
│  │     • 关键词匹配："莫奈" → art_history (90%)           │       │
│  │     • 内容分析：LLM分析用户输入 → 领域分类              │       │
│  │     • 历史记录：用户过往领域偏好                         │       │
│  │                                                         │       │
│  │  3️⃣ 降级策略：                                          │       │
│  │     • 目标适配器不可用 → 回退到默认适配器                │       │
│  │     • 多候选适配器 → 选择置信度最高的                    │       │
│  └─────────────────┬───────────────────────────────────────┘       │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │ 模块3: 适配器基类 (BaseDomainAdapter)                   │       │
│  │                                                         │       │
│  │  必须实现的接口方法：                                     │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 1. validate_input(raw_input) -> bool             │  │       │
│  │  │    验证输入是否符合领域规范                        │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 2. parse_input(raw_input) -> DomainInput         │  │       │
│  │  │    将原始输入解析为领域特定数据结构                │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 3. enrich_objectives(objectives) -> Objectives   │  │       │
│  │  │    用领域知识增强学习目标                          │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 4. customize_content_sources() -> List[Source]   │  │       │
│  │  │    指定领域特定的内容数据源                        │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 5. select_narrative_strategy() -> Strategy       │  │       │
│  │  │    选择适合领域的叙述策略                          │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 6. apply_domain_template(content) -> Output      │  │       │
│  │  │    应用领域特定的输出模板                          │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  │  ┌──────────────────────────────────────────────────┐  │       │
│  │  │ 7. validate_output(output) -> bool               │  │       │
│  │  │    验证输出是否符合领域标准                        │  │       │
│  │  └──────────────────────────────────────────────────┘  │       │
│  └─────────────────┬───────────────────────────────────────┘       │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │ 模块4: 具体适配器实现 (Concrete Adapters)               │       │
│  │                                                         │       │
│  │  ┌────────────────────┬────────────────────┐           │       │
│  │  │ K12Adapter         │ ArtHistoryAdapter  │           │       │
│  │  ├────────────────────┼────────────────────┤           │       │
│  │  │ • 年级/学科解析    │ • 艺术家/作品识别  │           │       │
│  │  │ • 知识点拆解       │ • 流派/年代判定    │           │       │
│  │  │ • 教学目标生成     │ • 艺术史数据源     │           │       │
│  │  │ • 教案模板         │ • 鉴赏模板         │           │       │
│  │  └────────────────────┴────────────────────┘           │       │
│  │                                                         │       │
│  │  ┌────────────────────┬────────────────────┐           │       │
│  │  │ BestsellerAdapter  │ VocationalAdapter  │           │       │
│  │  ├────────────────────┼────────────────────┤           │       │
│  │  │ • 书籍信息提取     │ • 技能分类         │           │       │
│  │  │ • 章节分析         │ • 工具识别         │           │       │
│  │  │ • 核心观点提炼     │ • 实操步骤生成     │           │       │
│  │  │ • 解读模板         │ • 培训课程模板     │           │       │
│  │  └────────────────────┴────────────────────┘           │       │
│  └─────────────────┬───────────────────────────────────────┘       │
│                    │                                                │
│                    ▼                                                │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │ 模块5: 插件管理器 (Plugin Manager)                      │       │
│  │                                                         │       │
│  │  功能：                                                  │       │
│  │  • 适配器安装：从文件/Git/Marketplace加载               │       │
│  │  • 依赖解析：检查Python包、配置文件依赖                  │       │
│  │  • 沙箱隔离：每个适配器独立命名空间                      │       │
│  │  • 热更新：无需重启更新适配器                            │       │
│  │  • 权限控制：限制适配器访问范围                          │       │
│  │                                                         │       │
│  │  生命周期管理：                                          │       │
│  │  Install → Validate → Activate → (Update) → Deactivate │       │
│  └─────────────────┬───────────────────────────────────────┘       │
│                    │                                                │
│                    ▼                                                │
│  ┌──────────────────────────────────────────┐                      │
│  │ 输出：领域特定的工作流配置                │                      │
│  │ (Domain-Specific Workflow Config)        │                      │
│  └──────────────────────────────────────────┘                      │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

---

#### 核心类设计

**1. 数据模型**

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any, Protocol
from enum import Enum
from abc import ABC, abstractmethod
import importlib
import inspect

class AdapterStatus(Enum):
    """适配器状态"""
    INSTALLED = "installed"  # 已安装但未激活
    ACTIVE = "active"  # 已激活
    INACTIVE = "inactive"  # 已停用
    ERROR = "error"  # 错误状态
    UPDATING = "updating"  # 更新中

class DomainType(Enum):
    """领域类型（扩展版）"""
    K12 = "k12"
    ART_HISTORY = "art_history"
    BESTSELLER_INTERPRETATION = "bestseller_interpretation"
    VOCATIONAL_TRAINING = "vocational_training"
    FINANCE = "finance"  # 未来扩展
    HEALTH = "health"  # 未来扩展
    PROGRAMMING = "programming"  # 未来扩展

@dataclass
class AdapterMetadata:
    """适配器元数据"""
    adapter_id: str  # 唯一标识，如 "k12-adapter-v2"
    name: str  # 显示名称
    domain: DomainType  # 领域类型
    version: str  # 语义化版本号
    author: str  # 作者
    description: str  # 描述
    
    # 依赖
    requires_python: str = ">=3.9"  # Python版本要求
    dependencies: List[str] = field(default_factory=list)  # Python包依赖
    
    # 配置
    config_schema: Dict[str, Any] = field(default_factory=dict)  # 配置模式
    supported_languages: List[str] = field(default_factory=lambda: ["zh"])  # 支持的语言
    
    # 状态
    status: AdapterStatus = AdapterStatus.INSTALLED
    enabled: bool = True
    
    # 元信息
    homepage: Optional[str] = None
    repository: Optional[str] = None
    license: str = "MIT"
    tags: List[str] = field(default_factory=list)

@dataclass
class DomainInput:
    """领域特定输入（由适配器解析）"""
    domain: DomainType
    raw_input: str  # 原始用户输入
    parsed_data: Dict[str, Any]  # 解析后的结构化数据
    metadata: Dict[str, Any] = field(default_factory=dict)

@dataclass
class DomainOutput:
    """领域特定输出（由适配器生成）"""
    domain: DomainType
    content: Any  # 具体内容（格式由适配器定义）
    format: str  # 输出格式（markdown/html/pdf等）
    metadata: Dict[str, Any] = field(default_factory=dict)

@dataclass
class AdapterConfig:
    """适配器配置"""
    adapter_id: str
    settings: Dict[str, Any] = field(default_factory=dict)
    enabled: bool = True
```

**2. 适配器基类**

```python
class BaseDomainAdapter(ABC):
    """
    领域适配器基类
    
    所有领域适配器必须继承此类并实现所有抽象方法
    """
    
    def __init__(self, config: AdapterConfig):
        self.config = config
        self.metadata = self._get_metadata()
    
    @classmethod
    @abstractmethod
    def _get_metadata(cls) -> AdapterMetadata:
        """
        返回适配器元数据
        
        示例：
        return AdapterMetadata(
            adapter_id="k12-adapter",
            name="K12教育适配器",
            domain=DomainType.K12,
            version="2.0.1",
            author="MetaWorkflow Team",
            description="K12教育领域适配器"
        )
        """
        pass
    
    @abstractmethod
    def validate_input(self, raw_input: str) -> bool:
        """
        验证输入是否符合领域规范
        
        Args:
            raw_input: 原始用户输入
            
        Returns:
            bool: 输入是否有效
        
        示例（K12）：
        - 检查是否包含"年级"、"学科"等关键词
        - 验证年级范围（1-12）
        """
        pass
    
    @abstractmethod
    async def parse_input(self, raw_input: str) -> DomainInput:
        """
        将原始输入解析为领域特定数据结构
        
        Args:
            raw_input: 原始用户输入
            
        Returns:
            DomainInput: 解析后的领域输入
        
        示例（K12）：
        "三年级数学除法" → {
            "grade": 3,
            "subject": "数学",
            "topic": "除法",
            "knowledge_points": ["除法概念", "除法运算"]
        }
        """
        pass
    
    @abstractmethod
    async def enrich_objectives(
        self,
        objectives: List[LearningObjective],
        domain_input: DomainInput
    ) -> List[LearningObjective]:
        """
        用领域知识增强学习目标
        
        Args:
            objectives: AI目标生成器生成的初始目标
            domain_input: 领域特定输入
            
        Returns:
            List[LearningObjective]: 增强后的学习目标
        
        示例（美术史）：
        - 添加艺术流派背景
        - 关联相关艺术家作品
        - 补充历史时代信息
        """
        pass
    
    @abstractmethod
    def customize_content_sources(
        self,
        domain_input: DomainInput
    ) -> List[ContentSource]:
        """
        指定领域特定的内容数据源
        
        Args:
            domain_input: 领域特定输入
            
        Returns:
            List[ContentSource]: 内容数据源列表
        
        示例（美术史）：
        - WikiArt（必选）
        - Google Arts & Culture
        - 特定博物馆网站
        
        示例（职业培训-Excel）：
        - Microsoft官方文档
        - Stack Overflow
        - YouTube教程
        """
        pass
    
    @abstractmethod
    def select_narrative_strategy(
        self,
        domain_input: DomainInput,
        audience: AudienceProfile
    ) -> NarrativeStrategy:
        """
        选择适合领域的叙述策略
        
        Args:
            domain_input: 领域特定输入
            audience: 受众画像
            
        Returns:
            NarrativeStrategy: 叙述策略
        
        示例（K12）：
        - 小学低年级 → ANALOGY_BASED（类比驱动）
        - 小学高年级 → SUSPENSE（悬念式）
        - 初中 → PROGRESSIVE（渐进式）
        """
        pass
    
    @abstractmethod
    async def apply_domain_template(
        self,
        narrative_content: NarrativeContent,
        domain_input: DomainInput
    ) -> DomainOutput:
        """
        应用领域特定的输出模板
        
        Args:
            narrative_content: 故事化叙述内容
            domain_input: 领域特定输入
            
        Returns:
            DomainOutput: 领域特定输出
        
        示例（K12）：
        - 转换为教案格式（教学目标、重难点、教学过程）
        - 添加课堂活动建议
        - 生成配套练习题
        
        示例（美术史）：
        - 转换为艺术鉴赏指南
        - 添加作品图片
        - 生成参观建议
        """
        pass
    
    @abstractmethod
    def validate_output(self, output: DomainOutput) -> bool:
        """
        验证输出是否符合领域标准
        
        Args:
            output: 领域特定输出
            
        Returns:
            bool: 输出是否有效
        
        示例（K12）：
        - 检查是否包含必需的教案章节
        - 验证教学目标的完整性
        - 检查难度适配性
        """
        pass
    
    # 可选的钩子方法（子类可以选择性覆盖）
    
    def on_activate(self):
        """适配器被激活时调用"""
        pass
    
    def on_deactivate(self):
        """适配器被停用时调用"""
        pass
    
    def get_config_schema(self) -> Dict[str, Any]:
        """返回配置模式（JSON Schema）"""
        return {}
    
    def health_check(self) -> bool:
        """健康检查"""
        return True


class AdapterRegistry:
    """适配器注册中心"""
    
    def __init__(self):
        self._adapters: Dict[str, BaseDomainAdapter] = {}
        self._metadata: Dict[str, AdapterMetadata] = {}
        self._domain_to_adapter: Dict[DomainType, str] = {}
    
    def register(
        self,
        adapter_class: type[BaseDomainAdapter],
        config: Optional[AdapterConfig] = None
    ) -> str:
        """
        注册适配器
        
        Args:
            adapter_class: 适配器类
            config: 配置（可选）
            
        Returns:
            str: 适配器ID
        """
        # 创建适配器实例
        if config is None:
            config = AdapterConfig(adapter_id=f"{adapter_class.__name__}")
        
        adapter = adapter_class(config)
        metadata = adapter.metadata
        
        # 检查ID是否已存在
        if metadata.adapter_id in self._adapters:
            raise ValueError(f"Adapter {metadata.adapter_id} already registered")
        
        # 注册
        self._adapters[metadata.adapter_id] = adapter
        self._metadata[metadata.adapter_id] = metadata
        self._domain_to_adapter[metadata.domain] = metadata.adapter_id
        
        logger.info(f"Registered adapter: {metadata.name} (ID: {metadata.adapter_id})")
        
        # 调用激活钩子
        if metadata.enabled:
            adapter.on_activate()
        
        return metadata.adapter_id
    
    def unregister(self, adapter_id: str):
        """注销适配器"""
        if adapter_id not in self._adapters:
            raise ValueError(f"Adapter {adapter_id} not found")
        
        adapter = self._adapters[adapter_id]
        
        # 调用停用钩子
        adapter.on_deactivate()
        
        # 移除
        domain = self._metadata[adapter_id].domain
        del self._adapters[adapter_id]
        del self._metadata[adapter_id]
        if domain in self._domain_to_adapter:
            del self._domain_to_adapter[domain]
        
        logger.info(f"Unregistered adapter: {adapter_id}")
    
    def get_adapter(self, adapter_id: str) -> Optional[BaseDomainAdapter]:
        """根据ID获取适配器"""
        return self._adapters.get(adapter_id)
    
    def get_adapter_by_domain(self, domain: DomainType) -> Optional[BaseDomainAdapter]:
        """根据领域获取适配器"""
        adapter_id = self._domain_to_adapter.get(domain)
        if adapter_id:
            return self._adapters.get(adapter_id)
        return None
    
    def list_adapters(
        self,
        domain: Optional[DomainType] = None,
        status: Optional[AdapterStatus] = None
    ) -> List[AdapterMetadata]:
        """列出适配器"""
        results = list(self._metadata.values())
        
        if domain:
            results = [m for m in results if m.domain == domain]
        
        if status:
            results = [m for m in results if m.status == status]
        
        return results
    
    def auto_discover(self, package_name: str = "adapters"):
        """
        自动发现和注册适配器
        
        扫描指定包下的所有适配器类并自动注册
        """
        try:
            module = importlib.import_module(package_name)
            
            # 遍历模块中的所有类
            for name, obj in inspect.getmembers(module, inspect.isclass):
                # 检查是否是BaseDomainAdapter的子类
                if issubclass(obj, BaseDomainAdapter) and obj != BaseDomainAdapter:
                    try:
                        self.register(obj)
                        logger.info(f"Auto-discovered adapter: {name}")
                    except Exception as e:
                        logger.error(f"Failed to register {name}: {e}")
        
        except ImportError:
            logger.warning(f"Package {package_name} not found")


class AdapterRouter:
    """适配器自动路由器"""
    
    def __init__(
        self,
        registry: AdapterRegistry,
        llm_client: Optional[LLMClient] = None
    ):
        self.registry = registry
        self.llm = llm_client
    
    async def route(
        self,
        raw_input: str,
        domain: Optional[DomainType] = None,
        user_history: Optional[List[str]] = None
    ) -> Optional[BaseDomainAdapter]:
        """
        路由到合适的适配器
        
        策略：
        1. 显式指定 → 直接路由
        2. 智能推断 → 关键词匹配 + LLM分析
        3. 历史偏好 → 用户过往领域
        4. 降级策略 → 默认适配器
        
        Args:
            raw_input: 原始用户输入
            domain: 显式指定的领域（可选）
            user_history: 用户历史记录（可选）
            
        Returns:
            BaseDomainAdapter: 适配器实例
        """
        # 策略1: 显式指定
        if domain:
            adapter = self.registry.get_adapter_by_domain(domain)
            if adapter and adapter.validate_input(raw_input):
                logger.info(f"Routed to {domain.value} (explicit)")
                return adapter
        
        # 策略2: 关键词匹配
        keyword_domain = self._match_by_keywords(raw_input)
        if keyword_domain:
            adapter = self.registry.get_adapter_by_domain(keyword_domain)
            if adapter and adapter.validate_input(raw_input):
                logger.info(f"Routed to {keyword_domain.value} (keywords)")
                return adapter
        
        # 策略3: LLM智能分析
        if self.llm:
            llm_domain = await self._infer_by_llm(raw_input)
            if llm_domain:
                adapter = self.registry.get_adapter_by_domain(llm_domain)
                if adapter and adapter.validate_input(raw_input):
                    logger.info(f"Routed to {llm_domain.value} (LLM inference)")
                    return adapter
        
        # 策略4: 用户历史偏好
        if user_history:
            history_domain = self._infer_from_history(user_history)
            if history_domain:
                adapter = self.registry.get_adapter_by_domain(history_domain)
                if adapter:
                    logger.info(f"Routed to {history_domain.value} (user history)")
                    return adapter
        
        # 降级：返回默认适配器（K12）
        logger.warning("No suitable adapter found, falling back to K12")
        return self.registry.get_adapter_by_domain(DomainType.K12)
    
    def _match_by_keywords(self, text: str) -> Optional[DomainType]:
        """基于关键词匹配领域"""
        keywords_map = {
            DomainType.K12: ["年级", "学科", "数学", "语文", "英语", "教学", "学生"],
            DomainType.ART_HISTORY: ["艺术", "画家", "作品", "流派", "莫奈", "梵高", "毕加索", "博物馆"],
            DomainType.BESTSELLER_INTERPRETATION: ["书", "作者", "畅销书", "阅读", "书籍", "小说"],
            DomainType.VOCATIONAL_TRAINING: ["Excel", "PPT", "编程", "技能", "操作", "培训", "软件"]
        }
        
        text_lower = text.lower()
        
        # 计算每个领域的匹配分数
        scores = {}
        for domain, keywords in keywords_map.items():
            score = sum(1 for kw in keywords if kw.lower() in text_lower)
            if score > 0:
                scores[domain] = score
        
        # 返回得分最高的领域
        if scores:
            best_domain = max(scores.items(), key=lambda x: x[1])[0]
            if scores[best_domain] >= 2:  # 至少匹配2个关键词
                return best_domain
        
        return None
    
    async def _infer_by_llm(self, text: str) -> Optional[DomainType]:
        """使用LLM推断领域"""
        prompt = f"""
        请判断以下用户输入属于哪个领域：
        
        用户输入：{text}
        
        候选领域：
        - k12: K12教育（小学、初中、高中）
        - art_history: 美术史
        - bestseller_interpretation: 畅销书解读
        - vocational_training: 职业技能培训
        
        请只返回领域代码（如"k12"），不要其他内容。
        """
        
        try:
            response = await self.llm.generate(
                prompt=prompt,
                max_tokens=20,
                temperature=0.1
            )
            
            domain_str = response.content.strip().lower()
            
            # 尝试匹配DomainType
            for domain in DomainType:
                if domain.value == domain_str:
                    return domain
        
        except Exception as e:
            logger.error(f"LLM inference failed: {e}")
        
        return None
    
    def _infer_from_history(self, history: List[str]) -> Optional[DomainType]:
        """从用户历史推断偏好领域"""
        # 简化实现：统计历史中各领域出现频率
        domain_counts = {}
        
        for entry in history[-10:]:  # 只看最近10条
            domain = self._match_by_keywords(entry)
            if domain:
                domain_counts[domain] = domain_counts.get(domain, 0) + 1
        
        if domain_counts:
            return max(domain_counts.items(), key=lambda x: x[1])[0]
        
        return None


class PluginManager:
    """插件管理器"""
    
    def __init__(self, registry: AdapterRegistry):
        self.registry = registry
        self.plugin_dir = Path("plugins/adapters")
        self.plugin_dir.mkdir(parents=True, exist_ok=True)
    
    async def install(
        self,
        source: str,
        source_type: str = "file"
    ) -> str:
        """
        安装适配器
        
        Args:
            source: 来源（文件路径/Git URL/Marketplace ID）
            source_type: 来源类型（file/git/marketplace）
            
        Returns:
            str: 适配器ID
        """
        if source_type == "file":
            return await self._install_from_file(source)
        elif source_type == "git":
            return await self._install_from_git(source)
        elif source_type == "marketplace":
            return await self._install_from_marketplace(source)
        else:
            raise ValueError(f"Unsupported source type: {source_type}")
    
    async def _install_from_file(self, file_path: str) -> str:
        """从文件安装"""
        # 简化实现
        # 1. 复制文件到plugin_dir
        # 2. 动态导入模块
        # 3. 注册到registry
        
        plugin_path = Path(file_path)
        if not plugin_path.exists():
            raise FileNotFoundError(f"Plugin file not found: {file_path}")
        
        # 复制到插件目录
        dest_path = self.plugin_dir / plugin_path.name
        shutil.copy(plugin_path, dest_path)
        
        # 动态导入
        spec = importlib.util.spec_from_file_location(
            plugin_path.stem,
            dest_path
        )
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        
        # 查找适配器类
        for name, obj in inspect.getmembers(module, inspect.isclass):
            if issubclass(obj, BaseDomainAdapter) and obj != BaseDomainAdapter:
                adapter_id = self.registry.register(obj)
                logger.info(f"Installed adapter from file: {adapter_id}")
                return adapter_id
        
        raise ValueError("No valid adapter class found in file")
    
    async def _install_from_git(self, git_url: str) -> str:
        """从Git仓库安装"""
        # TODO: 实现Git克隆和安装
        raise NotImplementedError()
    
    async def _install_from_marketplace(self, adapter_id: str) -> str:
        """从Marketplace安装"""
        # TODO: 实现Marketplace集成
        raise NotImplementedError()
    
    def uninstall(self, adapter_id: str):
        """卸载适配器"""
        self.registry.unregister(adapter_id)
        logger.info(f"Uninstalled adapter: {adapter_id}")
    
    def update(self, adapter_id: str, version: Optional[str] = None):
        """更新适配器"""
        # TODO: 实现版本更新逻辑
        raise NotImplementedError()
    
    def list_installed(self) -> List[AdapterMetadata]:
        """列出已安装的适配器"""
        return self.registry.list_adapters()
```

---

#### 具体适配器实现示例

**1. K12适配器**

```python
class K12Adapter(BaseDomainAdapter):
    """K12教育领域适配器"""
    
    @classmethod
    def _get_metadata(cls) -> AdapterMetadata:
        return AdapterMetadata(
            adapter_id="k12-adapter",
            name="K12教育适配器",
            domain=DomainType.K12,
            version="2.0.1",
            author="MetaWorkflow Team",
            description="支持小学、初中、高中教育内容生成",
            tags=["education", "k12", "teaching"]
        )
    
    def validate_input(self, raw_input: str) -> bool:
        """验证输入是否包含年级/学科信息"""
        grade_keywords = ["年级", "小学", "初中", "高中", "grade"]
        subject_keywords = ["数学", "语文", "英语", "物理", "化学", "math", "chinese"]
        
        has_grade = any(kw in raw_input for kw in grade_keywords)
        has_subject = any(kw in raw_input for kw in subject_keywords)
        
        return has_grade or has_subject
    
    async def parse_input(self, raw_input: str) -> DomainInput:
        """解析年级、学科、主题"""
        # 简化实现：使用正则或LLM提取
        parsed_data = {
            "grade": self._extract_grade(raw_input),
            "subject": self._extract_subject(raw_input),
            "topic": raw_input  # 简化：直接使用原文作为主题
        }
        
        return DomainInput(
            domain=DomainType.K12,
            raw_input=raw_input,
            parsed_data=parsed_data
        )
    
    def _extract_grade(self, text: str) -> Optional[int]:
        """提取年级"""
        import re
        match = re.search(r'(\d+)年级', text)
        if match:
            return int(match.group(1))
        
        grade_map = {
            "一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6,
            "七": 7, "八": 8, "九": 9,
            "小学": 3, "初中": 7, "高中": 10
        }
        
        for key, value in grade_map.items():
            if key in text:
                return value
        
        return None
    
    def _extract_subject(self, text: str) -> Optional[str]:
        """提取学科"""
        subjects = ["数学", "语文", "英语", "物理", "化学", "生物", "历史", "地理"]
        for subj in subjects:
            if subj in text:
                return subj
        return None
    
    async def enrich_objectives(
        self,
        objectives: List[LearningObjective],
        domain_input: DomainInput
    ) -> List[LearningObjective]:
        """增强学习目标（添加K12特有信息）"""
        grade = domain_input.parsed_data.get("grade")
        subject = domain_input.parsed_data.get("subject")
        
        # 为每个目标添加年级和学科信息
        for obj in objectives:
            obj.metadata["grade"] = grade
            obj.metadata["subject"] = subject
            
            # 根据年级调整难度
            if grade and grade <= 6:
                obj.bloom_level = min(obj.bloom_level, 3)  # 小学最高到"应用"
        
        return objectives
    
    def customize_content_sources(
        self,
        domain_input: DomainInput
    ) -> List[ContentSource]:
        """K12领域的数据源"""
        # K12优先使用教育类网站
        return [
            ContentSource(
                type=ContentSourceType.WIKIPEDIA,
                name="Wikipedia",
                base_url="https://zh.wikipedia.org",
                priority=2
            ),
            # 可以添加K12特定数据源，如教育部网站、教材资源等
        ]
    
    def select_narrative_strategy(
        self,
        domain_input: DomainInput,
        audience: AudienceProfile
    ) -> NarrativeStrategy:
        """根据年级选择叙述策略"""
        grade = domain_input.parsed_data.get("grade", 6)
        
        if grade <= 3:
            return NarrativeStrategy.ANALOGY_BASED  # 低年级：类比
        elif grade <= 6:
            return NarrativeStrategy.SUSPENSE  # 小学高年级：悬念
        elif grade <= 9:
            return NarrativeStrategy.PROGRESSIVE  # 初中：渐进式
        else:
            return NarrativeStrategy.PROBLEM_SOLVING  # 高中：问题解决
    
    async def apply_domain_template(
        self,
        narrative_content: NarrativeContent,
        domain_input: DomainInput
    ) -> DomainOutput:
        """转换为教案格式"""
        # 构建教案结构
        lesson_plan = {
            "title": narrative_content.title,
            "grade": domain_input.parsed_data.get("grade"),
            "subject": domain_input.parsed_data.get("subject"),
            
            # 教学目标
            "objectives": {
                "knowledge": [],  # 知识目标
                "ability": [],  # 能力目标
                "emotion": []  # 情感目标
            },
            
            # 教学重难点
            "key_points": narrative_content.key_points[:3],
            "difficult_points": narrative_content.key_points[3:5] if len(narrative_content.key_points) > 3 else [],
            
            # 教学过程
            "teaching_process": {
                "introduction": narrative_content.opening.content,
                "main_content": [section.content for section in narrative_content.body],
                "summary": narrative_content.closing.content
            },
            
            # 课堂活动
            "activities": self._generate_activities(narrative_content),
            
            # 作业
            "homework": self._generate_homework(narrative_content)
        }
        
        # 转换为Markdown格式
        markdown_output = self._format_as_markdown(lesson_plan)
        
        return DomainOutput(
            domain=DomainType.K12,
            content=lesson_plan,
            format="markdown",
            metadata={
                "template_version": "1.0",
                "word_count": len(markdown_output)
            }
        )
    
    def _generate_activities(self, content: NarrativeContent) -> List[str]:
        """生成课堂活动建议"""
        # 简化实现
        return [
            "小组讨论：分享对核心概念的理解",
            "动手实践：通过例题巩固知识点",
            "角色扮演：模拟应用场景"
        ]
    
    def _generate_homework(self, content: NarrativeContent) -> List[str]:
        """生成作业"""
        return [
            "基础题：复习课堂内容",
            "提高题：应用到新场景",
            "思考题：拓展延伸"
        ]
    
    def _format_as_markdown(self, lesson_plan: Dict) -> str:
        """格式化为Markdown"""
        md = f"""# {lesson_plan['title']}

## 基本信息
- **年级**：{lesson_plan['grade']}年级
- **学科**：{lesson_plan['subject']}

## 教学目标
### 知识目标
{chr(10).join(['- ' + obj for obj in lesson_plan['objectives']['knowledge']])}

## 教学重难点
### 重点
{chr(10).join(['- ' + point for point in lesson_plan['key_points']])}

### 难点
{chr(10).join(['- ' + point for point in lesson_plan['difficult_points']])}

## 教学过程
### 导入
{lesson_plan['teaching_process']['introduction']}

### 主体内容
{chr(10).join(lesson_plan['teaching_process']['main_content'])}

### 总结
{lesson_plan['teaching_process']['summary']}

## 课堂活动
{chr(10).join(['- ' + act for act in lesson_plan['activities']])}

## 作业布置
{chr(10).join(['- ' + hw for hw in lesson_plan['homework']])}
"""
        return md
    
    def validate_output(self, output: DomainOutput) -> bool:
        """验证教案完整性"""
        if output.format != "markdown":
            return False
        
        required_keys = ["title", "grade", "subject", "objectives", "teaching_process"]
        content = output.content
        
        return all(key in content for key in required_keys)


class ArtHistoryAdapter(BaseDomainAdapter):
    """美术史领域适配器"""
    
    @classmethod
    def _get_metadata(cls) -> AdapterMetadata:
        return AdapterMetadata(
            adapter_id="art-history-adapter",
            name="美术史适配器",
            domain=DomainType.ART_HISTORY,
            version="1.0.0",
            author="MetaWorkflow Team",
            description="支持艺术家、艺术作品、艺术流派的鉴赏内容生成",
            tags=["art", "history", "painting"]
        )
    
    def validate_input(self, raw_input: str) -> bool:
        art_keywords = ["艺术", "画家", "作品", "流派", "美术", "绘画", "art", "artist", "painting"]
        return any(kw in raw_input.lower() for kw in art_keywords)
    
    async def parse_input(self, raw_input: str) -> DomainInput:
        """解析艺术家、作品、流派"""
        parsed_data = {
            "artist": self._extract_artist(raw_input),
            "artwork": self._extract_artwork(raw_input),
            "movement": self._extract_movement(raw_input)
        }
        
        return DomainInput(
            domain=DomainType.ART_HISTORY,
            raw_input=raw_input,
            parsed_data=parsed_data
        )
    
    def _extract_artist(self, text: str) -> Optional[str]:
        """提取艺术家名字"""
        famous_artists = ["莫奈", "梵高", "毕加索", "达芬奇", "米开朗琪罗", "Monet", "Van Gogh"]
        for artist in famous_artists:
            if artist in text:
                return artist
        return None
    
    def _extract_artwork(self, text: str) -> Optional[str]:
        """提取作品名"""
        # 简化实现
        return None
    
    def _extract_movement(self, text: str) -> Optional[str]:
        """提取艺术流派"""
        movements = ["印象派", "立体主义", "表现主义", "浪漫主义", "文艺复兴"]
        for mov in movements:
            if mov in text:
                return mov
        return None
    
    async def enrich_objectives(
        self,
        objectives: List[LearningObjective],
        domain_input: DomainInput
    ) -> List[LearningObjective]:
        """添加艺术史背景"""
        artist = domain_input.parsed_data.get("artist")
        movement = domain_input.parsed_data.get("movement")
        
        for obj in objectives:
            obj.metadata["artist"] = artist
            obj.metadata["movement"] = movement
        
        return objectives
    
    def customize_content_sources(
        self,
        domain_input: DomainInput
    ) -> List[ContentSource]:
        """美术史特定数据源"""
        return [
            ContentSource(
                type=ContentSourceType.WIKIART,
                name="WikiArt",
                base_url="https://www.wikiart.org",
                priority=1  # 最高优先级
            ),
            ContentSource(
                type=ContentSourceType.GOOGLE_ARTS,
                name="Google Arts & Culture",
                base_url="https://artsandculture.google.com",
                priority=1
            ),
            ContentSource(
                type=ContentSourceType.WIKIPEDIA,
                name="Wikipedia",
                base_url="https://zh.wikipedia.org",
                priority=2
            )
        ]
    
    def select_narrative_strategy(
        self,
        domain_input: DomainInput,
        audience: AudienceProfile
    ) -> NarrativeStrategy:
        """美术史优先使用传记式或作品背景故事"""
        artist = domain_input.parsed_data.get("artist")
        
        if artist:
            return NarrativeStrategy.ARTIST_BIOGRAPHY
        else:
            return NarrativeStrategy.ARTWORK_BACKSTORY
    
    async def apply_domain_template(
        self,
        narrative_content: NarrativeContent,
        domain_input: DomainInput
    ) -> DomainOutput:
        """转换为艺术鉴赏指南"""
        guide = {
            "title": narrative_content.title,
            "artist": domain_input.parsed_data.get("artist"),
            "movement": domain_input.parsed_data.get("movement"),
            
            "introduction": narrative_content.opening.content,
            "appreciation_points": narrative_content.key_points,
            "detailed_analysis": [s.content for s in narrative_content.body],
            "conclusion": narrative_content.closing.content,
            
            "artworks": [],  # 可以从sources中提取
            "references": narrative_content.references
        }
        
        return DomainOutput(
            domain=DomainType.ART_HISTORY,
            content=guide,
            format="markdown"
        )
    
    def validate_output(self, output: DomainOutput) -> bool:
        return "title" in output.content and "appreciation_points" in output.content
```

---

#### 监控指标

```python
class AdapterMetrics:
    """适配器监控指标"""
    
    def __init__(self):
        self.adapter_usage = Counter('adapter_usage_total', 'Adapter usage', ['adapter_id'])
        self.routing_duration = Histogram('adapter_routing_duration_seconds', 'Routing time')
        self.routing_strategy = Counter('adapter_routing_strategy', 'Routing strategy', ['strategy'])
        self.adapter_errors = Counter('adapter_errors_total', 'Adapter errors', ['adapter_id', 'error_type'])
    
    def record_usage(self, adapter_id: str, duration: float, strategy: str):
        """记录适配器使用"""
        self.adapter_usage.labels(adapter_id=adapter_id).inc()
        self.routing_duration.observe(duration)
        self.routing_strategy.labels(strategy=strategy).inc()
    
    def record_error(self, adapter_id: str, error_type: str):
        """记录错误"""
        self.adapter_errors.labels(adapter_id=adapter_id, error_type=error_type).inc()
```

---

### 3.1 工作流编排引擎设计

#### 3.1.1 工作流定义DSL

使用YAML定义工作流模板：

```yaml
# math_lesson_workflow.yaml
workflow:
  id: "math_lesson_basic"
  name: "数学课程标准工作流"
  version: "1.0"
  
  # 触发条件
  trigger:
    subject: "数学"
    grade_range: [1, 12]
  
  # 全局配置
  config:
    timeout: 180  # 秒
    retry_policy:
      max_retries: 3
      backoff: "exponential"
    
  # 节点定义
  nodes:
    - id: "persona"
      type: "PersonaAnalysisNode"
      config:
        model: "gpt-4"
        temperature: 0.7
      
    - id: "objectives"
      type: "ObjectiveSettingNode"
      depends_on: ["persona"]
      config:
        bloom_level: "application"
      
    - id: "rag_retrieval"
      type: "RAGRetrievalNode"
      depends_on: ["objectives"]
      config:
        knowledge_base: "${region}_math_textbook"
        top_k: 5
      
    - id: "logic_check"
      type: "LogicVerificationNode"
      depends_on: ["content"]
      condition: "subject == '数学'"  # 条件节点
      
    - id: "content"
      type: "ContentGenerationNode"
      depends_on: ["objectives", "rag_retrieval"]
      config:
        template: "narrative_math"
        
    - id: "compliance"
      type: "ComplianceCheckNode"
      depends_on: ["content"]
      config:
        rules: "${region}_compliance"
  
  # 输出定义
  output:
    format: "json"
    schema: "lesson_plan_v1"
```

#### 3.1.2 执行流程

```python
class WorkflowExecutor:
    def __init__(self, workflow_def: dict):
        self.workflow = self.parse_workflow(workflow_def)
        self.state_manager = StateManager()
        self.event_bus = EventBus()
    
    async def execute(self, context: WorkflowContext):
        """执行工作流"""
        
        # 1. 创建工作流实例
        instance_id = self.create_instance(context)
        
        try:
            # 2. 构建执行计划
            execution_plan = self.build_execution_plan()
            
            # 3. 按拓扑顺序执行节点
            for level in execution_plan.levels:
                # 同一层级的节点可以并发执行
                tasks = [
                    self.execute_node(node, context, instance_id)
                    for node in level
                ]
                results = await asyncio.gather(*tasks)
                
                # 更新上下文
                for node, result in zip(level, results):
                    context.set_node_output(node.id, result)
                    self.state_manager.save_progress(instance_id, node.id, result)
            
            # 4. 生成最终输出
            final_output = self.format_output(context)
            
            # 5. 发送完成事件
            self.event_bus.publish("workflow.completed", {
                "instance_id": instance_id,
                "output": final_output
            })
            
            return final_output
            
        except Exception as e:
            # 错误处理
            self.handle_error(instance_id, e)
            raise
    
    async def execute_node(self, node, context, instance_id):
        """执行单个节点"""
        
        # 条件检查
        if node.has_condition() and not node.evaluate_condition(context):
            return None  # 跳过此节点
        
        # 输入验证
        if not node.validate_input(context):
            raise ValidationError(f"Node {node.id} input invalid")
        
        # 执行
        try:
            result = await node.execute(context)
            return result
        except Exception as e:
            # 重试逻辑
            if node.should_retry(e):
                return await self.retry_node(node, context)
            else:
                raise
```

### 3.2 工作流模板管理器设计

#### 3.2.1 模板生成与存储

**需求**：元工作流执行后，自动生成可复用的模板并持久化存储。

```python
class WorkflowTemplateManager:
    """工作流模板管理器"""
    
    def __init__(self, db, cache):
        self.db = db
        self.cache = cache
        self.version_control = VersionControl()
    
    async def create_template_from_instance(
        self, 
        instance: WorkflowInstance,
        template_name: str,
        description: str = ""
    ) -> WorkflowTemplate:
        """从工作流实例创建模板"""
        
        # 1. 提取实例的结构信息
        nodes = self._extract_nodes(instance)
        edges = self._extract_edges(instance)
        parameters = self._extract_parameters(instance)
        
        # 2. 生成模板元数据
        template = WorkflowTemplate(
            template_id=self._generate_template_id(),
            template_name=template_name,
            version="1.0",
            metadata={
                "subject": instance.context.subject,
                "grade_range": self._infer_grade_range(instance),
                "region": instance.context.region,
                "created_by": instance.user_id,
                "created_at": datetime.utcnow(),
                "description": description,
                "source_instance_id": instance.instance_id
            },
            nodes=nodes,
            edges=edges,
            parameters=parameters
        )
        
        # 3. 保存到数据库
        await self.db.templates.insert(template.to_dict())
        
        # 4. 缓存热门模板
        await self.cache.set(
            f"template:{template.template_id}",
            template.to_json(),
            ttl=3600
        )
        
        # 5. 建立实例与模板的关联
        await self.db.instance_template_mapping.insert({
            "instance_id": instance.instance_id,
            "template_id": template.template_id,
            "relationship": "source"
        })
        
        return template
    
    def _extract_nodes(self, instance: WorkflowInstance) -> List[Dict]:
        """提取节点信息"""
        nodes = []
        for node_id, execution in instance.node_executions.items():
            nodes.append({
                "node_id": node_id,
                "node_type": execution.node_type,
                "config": execution.config,
                "position": execution.position,  # 用于可视化
                "description": execution.description
            })
        return nodes
    
    def _extract_edges(self, instance: WorkflowInstance) -> List[Dict]:
        """提取节点依赖关系"""
        edges = []
        for node_id, execution in instance.node_executions.items():
            for dependency in execution.depends_on:
                edges.append({
                    "from": dependency,
                    "to": node_id,
                    "type": "data_dependency"
                })
        return edges
    
    def _extract_parameters(self, instance: WorkflowInstance) -> Dict:
        """提取可参数化的字段"""
        return {
            "required": ["topics", "duration"],
            "optional": ["textbook_version", "enable_story", "interaction_level"],
            "defaults": {
                "duration": 45,
                "enable_story": True
            }
        }
```

#### 3.2.2 模板查找与实例化

```python
class WorkflowTemplateManager:
    
    async def find_template(
        self, 
        template_id: str = None,
        template_name: str = None,
        filters: Dict = None
    ) -> Optional[WorkflowTemplate]:
        """查找模板"""
        
        # 优先从缓存查找
        if template_id:
            cached = await self.cache.get(f"template:{template_id}")
            if cached:
                return WorkflowTemplate.from_json(cached)
        
        # 从数据库查找
        query = {}
        if template_id:
            query["template_id"] = template_id
        elif template_name:
            query["template_name"] = template_name
        
        if filters:
            query.update(filters)
        
        result = await self.db.templates.find_one(query)
        
        if result:
            return WorkflowTemplate.from_dict(result)
        
        return None
    
    async def search_templates(
        self,
        subject: str = None,
        grade: int = None,
        region: str = None,
        tags: List[str] = None,
        sort_by: str = "popularity",  # popularity/rating/recent
        limit: int = 20
    ) -> List[WorkflowTemplate]:
        """搜索模板"""
        
        query = {}
        
        if subject:
            query["metadata.subject"] = subject
        
        if grade:
            query["metadata.grade_range"] = {"$lte": grade, "$gte": grade}
        
        if region:
            query["metadata.region"] = region
        
        if tags:
            query["tags"] = {"$in": tags}
        
        # 排序规则
        sort = {
            "popularity": {"execution_history.total_runs": -1},
            "rating": {"ratings.avg_score": -1},
            "recent": {"metadata.created_at": -1}
        }
        
        results = await self.db.templates.find(query) \
            .sort(sort[sort_by]) \
            .limit(limit) \
            .to_list()
        
        return [WorkflowTemplate.from_dict(r) for r in results]
    
    async def instantiate_template(
        self,
        template_id: str,
        parameters: Dict,
        user_id: str
    ) -> WorkflowInstance:
        """实例化模板"""
        
        # 1. 加载模板
        template = await self.find_template(template_id=template_id)
        if not template:
            raise TemplateNotFoundError(f"Template {template_id} not found")
        
        # 2. 验证参数
        self._validate_parameters(template, parameters)
        
        # 3. 创建工作流上下文
        context = WorkflowContext(
            instance_id=self._generate_instance_id(),
            **parameters,
            template_id=template_id,
            template_version=template.version
        )
        
        # 4. 构建工作流实例
        instance = WorkflowInstance(
            instance_id=context.instance_id,
            workflow_id=template.template_id,
            template_id=template.template_id,
            template_version=template.version,
            context=context,
            status=InstanceStatus.PENDING,
            user_id=user_id
        )
        
        # 5. 复制节点配置
        for node in template.nodes:
            instance.node_executions[node["node_id"]] = NodeExecution(
                node_id=node["node_id"],
                node_type=node["node_type"],
                config=node["config"],
                status="pending"
            )
        
        # 6. 保存实例
        await self.db.instances.insert(instance.to_dict())
        
        # 7. 更新模板使用统计
        await self.db.templates.update_one(
            {"template_id": template_id},
            {"$inc": {"execution_history.total_runs": 1}}
        )
        
        return instance
```

#### 3.2.3 模板版本控制

```python
class VersionControl:
    """模板版本控制"""
    
    async def create_new_version(
        self,
        template_id: str,
        modifications: Dict,
        change_log: str,
        user_id: str
    ) -> WorkflowTemplate:
        """创建新版本"""
        
        # 1. 加载当前版本
        current = await self.template_manager.find_template(template_id)
        
        # 2. 应用修改
        new_template = self._apply_modifications(current, modifications)
        
        # 3. 增加版本号
        new_version = self._increment_version(current.version, modifications)
        new_template.version = new_version
        new_template.base_template_id = current.template_id
        new_template.metadata["created_at"] = datetime.utcnow()
        new_template.metadata["created_by"] = user_id
        
        # 4. 记录变更
        new_template.changelog = current.changelog or []
        new_template.changelog.append({
            "version": new_version,
            "date": datetime.utcnow().isoformat(),
            "changes": change_log,
            "modifications": modifications,
            "author": user_id
        })
        
        # 5. 保存新版本
        await self.db.templates.insert(new_template.to_dict())
        
        # 6. 更新版本链
        await self.db.template_versions.insert({
            "template_id": template_id,
            "version": new_version,
            "parent_version": current.version,
            "created_at": datetime.utcnow()
        })
        
        return new_template
    
    def _apply_modifications(
        self, 
        template: WorkflowTemplate, 
        modifications: Dict
    ) -> WorkflowTemplate:
        """应用修改到模板"""
        
        new_template = deepcopy(template)
        
        # 节点修改
        if "nodes" in modifications:
            for mod in modifications["nodes"]:
                if mod["action"] == "add":
                    new_template.nodes.append(mod["node"])
                elif mod["action"] == "remove":
                    new_template.nodes = [
                        n for n in new_template.nodes 
                        if n["node_id"] != mod["node_id"]
                    ]
                elif mod["action"] == "update":
                    for node in new_template.nodes:
                        if node["node_id"] == mod["node_id"]:
                            node.update(mod["updates"])
                elif mod["action"] == "replace":
                    new_template.nodes = [
                        mod["new_node"] if n["node_id"] == mod["node_id"] else n
                        for n in new_template.nodes
                    ]
        
        # 边(依赖关系)修改
        if "edges" in modifications:
            for mod in modifications["edges"]:
                if mod["action"] == "add":
                    new_template.edges.append(mod["edge"])
                elif mod["action"] == "remove":
                    new_template.edges = [
                        e for e in new_template.edges
                        if not (e["from"] == mod["from"] and e["to"] == mod["to"])
                    ]
        
        # 配置修改
        if "config" in modifications:
            new_template.config.update(modifications["config"])
        
        # 验证修改后的模板
        self._validate_template(new_template)
        
        return new_template
    
    def _increment_version(self, current_version: str, modifications: Dict) -> str:
        """根据修改类型决定版本号增长"""
        
        major, minor, patch = map(int, current_version.split("."))
        
        # 判断修改类型
        has_breaking_change = self._has_breaking_change(modifications)
        has_new_feature = self._has_new_feature(modifications)
        
        if has_breaking_change:
            # 重大变更，主版本号+1
            return f"{major + 1}.0.0"
        elif has_new_feature:
            # 新功能，次版本号+1
            return f"{major}.{minor + 1}.0"
        else:
            # 小修改，补丁版本号+1
            return f"{major}.{minor}.{patch + 1}"
    
    async def compare_versions(
        self,
        template_id: str,
        version1: str,
        version2: str
    ) -> Dict:
        """对比两个版本的差异"""
        
        t1 = await self.find_template(template_id, version=version1)
        t2 = await self.find_template(template_id, version=version2)
        
        return {
            "nodes": {
                "added": [n for n in t2.nodes if n not in t1.nodes],
                "removed": [n for n in t1.nodes if n not in t2.nodes],
                "modified": self._find_modified_nodes(t1.nodes, t2.nodes)
            },
            "edges": {
                "added": [e for e in t2.edges if e not in t1.edges],
                "removed": [e for e in t1.edges if e not in t2.edges]
            },
            "config_diff": self._diff_dicts(t1.config, t2.config)
        }
    
    async def rollback_to_version(
        self,
        template_id: str,
        target_version: str,
        user_id: str
    ) -> WorkflowTemplate:
        """回滚到指定版本"""
        
        # 1. 加载目标版本
        target_template = await self.find_template(template_id, version=target_version)
        
        # 2. 创建新版本（基于目标版本）
        new_template = deepcopy(target_template)
        new_template.version = self._get_next_version(template_id)
        new_template.metadata["created_at"] = datetime.utcnow()
        new_template.metadata["created_by"] = user_id
        
        # 3. 记录回滚操作
        new_template.changelog.append({
            "version": new_template.version,
            "date": datetime.utcnow().isoformat(),
            "changes": f"Rollback to version {target_version}",
            "type": "rollback",
            "author": user_id
        })
        
        # 4. 保存
        await self.db.templates.insert(new_template.to_dict())
        
        return new_template
```

#### 3.2.4 执行记忆与反馈优化

```python
class ExecutionMemory:
    """执行记忆系统"""
    
    async def record_execution(
        self,
        instance: WorkflowInstance,
        execution_result: Dict
    ):
        """记录执行详情"""
        
        execution_record = {
            "instance_id": instance.instance_id,
            "template_id": instance.template_id,
            "template_version": instance.template_version,
            
            "executed_at": datetime.utcnow(),
            
            # 执行轨迹
            "execution_trace": {
                "nodes": [
                    {
                        "node_id": ne.node_id,
                        "node_type": ne.node_type,
                        "started_at": ne.started_at,
                        "completed_at": ne.completed_at,
                        "duration": (ne.completed_at - ne.started_at).total_seconds(),
                        "status": ne.status,
                        "retry_count": ne.retry_count,
                        "input_size": len(str(ne.input_data)),
                        "output_size": len(str(ne.output_data))
                    }
                    for ne in instance.node_executions.values()
                ],
                "total_duration": execution_result["duration"],
                "success": execution_result["success"]
            },
            
            # 资源消耗
            "resource_usage": {
                "tokens": execution_result.get("total_tokens", 0),
                "cost": execution_result.get("total_cost", 0.0),
                "models_used": execution_result.get("models_used", [])
            },
            
            # 质量指标
            "quality_metrics": execution_result.get("quality_metrics", {})
        }
        
        # 保存执行记录
        await self.db.execution_records.insert(execution_record)
        
        # 更新模板统计
        await self._update_template_stats(instance.template_id, execution_result)
    
    async def analyze_and_suggest(
        self,
        template_id: str,
        time_range_days: int = 30
    ) -> List[Dict]:
        """分析执行数据并提供优化建议"""
        
        # 1. 查询最近的执行记录
        records = await self.db.execution_records.find({
            "template_id": template_id,
            "executed_at": {
                "$gte": datetime.utcnow() - timedelta(days=time_range_days)
            }
        }).to_list(length=1000)
        
        if not records:
            return []
        
        suggestions = []
        
        # 2. 性能分析
        performance_issues = self._analyze_performance(records)
        suggestions.extend(performance_issues)
        
        # 3. 质量分析
        quality_issues = self._analyze_quality(records)
        suggestions.extend(quality_issues)
        
        # 4. 成本分析
        cost_issues = self._analyze_cost(records)
        suggestions.extend(cost_issues)
        
        # 5. 按严重程度排序
        suggestions.sort(key=lambda s: {"high": 0, "medium": 1, "low": 2}[s["severity"]])
        
        return suggestions
    
    def _analyze_performance(self, records: List[Dict]) -> List[Dict]:
        """性能分析"""
        suggestions = []
        
        # 统计各节点的平均执行时间
        node_durations = defaultdict(list)
        for record in records:
            for node in record["execution_trace"]["nodes"]:
                node_durations[node["node_id"]].append(node["duration"])
        
        # 识别慢节点
        for node_id, durations in node_durations.items():
            avg_duration = sum(durations) / len(durations)
            if avg_duration > 30:  # 超过30秒
                suggestions.append({
                    "type": "performance",
                    "severity": "medium",
                    "node_id": node_id,
                    "issue": f"节点平均耗时{avg_duration:.1f}秒，超过预期",
                    "suggestion": "建议优化或增加缓存",
                    "data": {
                        "avg_duration": avg_duration,
                        "max_duration": max(durations),
                        "min_duration": min(durations)
                    }
                })
        
        return suggestions
    
    def _analyze_quality(self, records: List[Dict]) -> List[Dict]:
        """质量分析"""
        suggestions = []
        
        # 统计质量指标
        quality_scores = [
            r["quality_metrics"].get("overall_score", 0)
            for r in records
            if "quality_metrics" in r
        ]
        
        if quality_scores:
            avg_score = sum(quality_scores) / len(quality_scores)
            if avg_score < 0.7:  # 低于70分
                suggestions.append({
                    "type": "quality",
                    "severity": "high",
                    "issue": f"内容质量平均分{avg_score:.2f}，低于预期",
                    "suggestion": "建议优化Prompt模板或增加few-shot示例"
                })
        
        return suggestions
    
    def _analyze_cost(self, records: List[Dict]) -> List[Dict]:
        """成本分析"""
        suggestions = []
        
        # 统计成本
        total_cost = sum(r["resource_usage"]["cost"] for r in records)
        avg_cost = total_cost / len(records)
        
        # 识别高成本节点
        node_costs = defaultdict(list)
        for record in records:
            # 这里需要从node级别记录成本
            pass
        
        return suggestions
```

---

### 3.2.1 智能模板推荐匹配系统 ⭐ 新增

#### 系统概述

智能模板推荐匹配系统负责根据教师的需求（学科、年级、课型、关键词等），自动推荐最合适的工作流模板。该系统结合了**规则匹配**、**向量语义检索**和**协同过滤**三种策略，实现精准高效的模板推荐。

**核心目标**：
- 🎯 **精准性**：推荐的模板高度匹配教师需求
- ⚡ **实时性**：毫秒级响应时间
- 🔄 **个性化**：基于教师使用历史优化推荐
- 📊 **可解释**：清晰说明推荐理由

#### 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                    Template Recommendation Engine                │
│                        (模板推荐引擎)                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────┴─────────────────────┐
        │                                             │
        ▼                                             ▼
┌──────────────────┐                      ┌──────────────────────┐
│   Rule-based     │                      │   Semantic Search    │
│   Matcher        │◄────────────────────►│   Engine             │
│  (规则匹配器)     │      Hybrid          │  (语义搜索引擎)       │
└──────────────────┘      Strategy        └──────────────────────┘
        │                                             │
        │                                             │
        ▼                                             ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Collaborative Filtering                        │
│                    (协同过滤模块)                                 │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Ranking & Scoring                              │
│                    (排序与评分模块)                               │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Recommendation Results                         │
│                    (推荐结果)                                     │
└──────────────────────────────────────────────────────────────────┘
```

---

#### 三层匹配策略

##### **1. 规则匹配层（Rule-based Matching）**

**目的**：快速过滤出符合硬性条件的候选模板

**匹配规则**：

```python
class TemplateRuleMatcher:
    """基于规则的模板匹配器"""
    
    def match(self, query: TemplateQuery, templates: List[Template]) -> List[Template]:
        """规则匹配"""
        candidates = []
        
        for template in templates:
            score = 0
            match_reasons = []
            
            # 1. 学科匹配（必须）
            if not self._subject_match(query.subject, template.metadata.subject):
                continue
            score += 20
            match_reasons.append(f"学科匹配: {template.metadata.subject}")
            
            # 2. 年级范围匹配（必须）
            if not self._grade_match(query.grade, template.metadata.grade_range):
                continue
            score += 20
            match_reasons.append(f"年级匹配: {template.metadata.grade_range}")
            
            # 3. 地区匹配（可选，加分项）
            if query.region and query.region == template.metadata.region:
                score += 15
                match_reasons.append(f"地区匹配: {template.metadata.region}")
            
            # 4. 课型匹配（可选，加分项）
            if query.lesson_type and query.lesson_type in template.metadata.lesson_types:
                score += 10
                match_reasons.append(f"课型匹配: {query.lesson_type}")
            
            # 5. 标签匹配（可选，加分项）
            if query.tags:
                tag_overlap = set(query.tags) & set(template.metadata.tags)
                if tag_overlap:
                    score += len(tag_overlap) * 5
                    match_reasons.append(f"标签匹配: {', '.join(tag_overlap)}")
            
            # 6. 语言匹配（可选）
            if query.language and query.language == template.metadata.language:
                score += 10
                match_reasons.append(f"语言匹配: {template.metadata.language}")
            
            candidates.append({
                "template": template,
                "rule_score": score,
                "match_reasons": match_reasons
            })
        
        # 按规则分数排序
        candidates.sort(key=lambda x: x["rule_score"], reverse=True)
        return candidates
    
    def _subject_match(self, query_subject: str, template_subject: str) -> bool:
        """学科匹配逻辑"""
        # 精确匹配
        if query_subject == template_subject:
            return True
        
        # 学科别名映射
        subject_aliases = {
            "数学": ["math", "mathematics", "算数"],
            "语文": ["chinese", "语言", "阅读"],
            "英语": ["english", "外语"],
            # ... 更多映射
        }
        
        for canonical, aliases in subject_aliases.items():
            if query_subject in aliases and template_subject in aliases:
                return True
        
        return False
    
    def _grade_match(self, query_grade: int, template_range: Tuple[int, int]) -> bool:
        """年级范围匹配"""
        min_grade, max_grade = template_range
        return min_grade <= query_grade <= max_grade
```

**示例**：
```python
# 查询条件
query = TemplateQuery(
    subject="数学",
    grade=3,
    region="CN",
    lesson_type="新授课",
    tags=["故事化", "探究式"]
)

# 匹配结果
# Template A: 学科✓ 年级✓ 地区✓ 课型✓ 标签✓✓ → 得分80
# Template B: 学科✓ 年级✓ 地区✗ 课型✓ 标签✓  → 得分55
# Template C: 学科✓ 年级✓ 地区✗ 课型✗ 标签✗  → 得分40
```

---

##### **2. 语义检索层（Semantic Search）**

**目的**：基于关键词、描述的语义相似度，找到概念上相关的模板

**技术方案**：
- **向量数据库**：Milvus / Qdrant
- **Embedding模型**：OpenAI text-embedding-3-small / BGE-M3
- **检索策略**：混合检索（向量 + BM25）

**实现流程**：

```python
class TemplateSemanticSearchEngine:
    """模板语义搜索引擎"""
    
    def __init__(self, vector_db: VectorDatabase, embedding_model: EmbeddingModel):
        self.vector_db = vector_db
        self.embedding_model = embedding_model
    
    def search(self, query: TemplateQuery, top_k: int = 20) -> List[Dict]:
        """语义搜索"""
        # 1. 构建查询文本
        query_text = self._build_query_text(query)
        
        # 2. 生成查询向量
        query_vector = self.embedding_model.encode(query_text)
        
        # 3. 向量检索（ANN - Approximate Nearest Neighbor）
        vector_results = self.vector_db.search(
            collection_name="templates",
            query_vector=query_vector,
            top_k=top_k,
            metric_type="COSINE"  # 余弦相似度
        )
        
        # 4. BM25关键词检索（可选，混合检索）
        if query.keywords:
            bm25_results = self._bm25_search(query.keywords, top_k)
            # 合并结果
            results = self._merge_results(vector_results, bm25_results)
        else:
            results = vector_results
        
        # 5. 计算语义相似度分数
        for result in results:
            result["semantic_score"] = result["distance"]  # 0-1之间
            result["match_type"] = "semantic"
        
        return results
    
    def _build_query_text(self, query: TemplateQuery) -> str:
        """构建查询文本（用于embedding）"""
        parts = []
        
        if query.subject:
            parts.append(f"学科: {query.subject}")
        
        if query.grade:
            parts.append(f"年级: {query.grade}年级")
        
        if query.lesson_type:
            parts.append(f"课型: {query.lesson_type}")
        
        if query.keywords:
            parts.append(f"关键词: {', '.join(query.keywords)}")
        
        if query.description:
            parts.append(f"描述: {query.description}")
        
        if query.tags:
            parts.append(f"标签: {', '.join(query.tags)}")
        
        return " | ".join(parts)
    
    def _bm25_search(self, keywords: List[str], top_k: int) -> List[Dict]:
        """BM25关键词检索"""
        # 使用PostgreSQL的全文搜索或Elasticsearch
        query_sql = """
            SELECT 
                template_id,
                ts_rank(search_vector, plainto_tsquery('chinese', %s)) as rank
            FROM templates
            WHERE search_vector @@ plainto_tsquery('chinese', %s)
            ORDER BY rank DESC
            LIMIT %s
        """
        keyword_str = " ".join(keywords)
        results = self.db.execute(query_sql, (keyword_str, keyword_str, top_k))
        return results
    
    def _merge_results(self, vector_results: List, bm25_results: List) -> List:
        """混合检索结果合并（Reciprocal Rank Fusion）"""
        # RRF算法：1 / (k + rank)，k通常取60
        k = 60
        scores = defaultdict(float)
        
        # 向量检索结果
        for rank, result in enumerate(vector_results, 1):
            template_id = result["template_id"]
            scores[template_id] += 1 / (k + rank)
        
        # BM25检索结果
        for rank, result in enumerate(bm25_results, 1):
            template_id = result["template_id"]
            scores[template_id] += 1 / (k + rank)
        
        # 按分数排序
        merged = sorted(scores.items(), key=lambda x: x[1], reverse=True)
        return [{"template_id": tid, "fusion_score": score} for tid, score in merged]
```

**向量化策略**：

每个模板在入库时，会生成以下文本的embedding：

```python
def generate_template_embedding_text(template: Template) -> str:
    """生成模板的向量化文本"""
    parts = [
        f"模板名称: {template.name}",
        f"学科: {template.metadata.subject}",
        f"年级: {template.metadata.grade_range[0]}-{template.metadata.grade_range[1]}年级",
        f"课型: {', '.join(template.metadata.lesson_types)}",
        f"描述: {template.metadata.description}",
        f"标签: {', '.join(template.metadata.tags)}",
        f"教学目标: {template.metadata.learning_objectives}",
        f"适用场景: {template.metadata.use_cases}",
    ]
    return "\n".join(parts)
```

**示例**：
```python
# 查询："我需要一个适合小学三年级的数学课，讲解分数加法，要有趣味性"
# 
# 检索过程：
# 1. 生成query embedding
# 2. 在Milvus中检索top 20相似模板
# 3. 返回语义相似度分数

# 结果示例：
# Template A: 语义相似度 0.89 (描述中包含"分数"、"趣味"关键词)
# Template B: 语义相似度 0.76 (适合三年级，但课题不同)
# Template C: 语义相似度 0.65 (年级不匹配，但教学方法相似)
```

---

##### **3. 协同过滤层（Collaborative Filtering）**

**目的**：基于用户行为数据，提供个性化推荐

**策略**：
- **基于用户的协同过滤**（User-based CF）：找到相似用户喜欢的模板
- **基于物品的协同过滤**（Item-based CF）：找到与用户历史使用模板相似的模板

**实现**：

```python
class CollaborativeFilteringRecommender:
    """协同过滤推荐器"""
    
    def __init__(self, usage_history: UsageHistoryService):
        self.usage_history = usage_history
        self.user_similarity_cache = {}
        self.item_similarity_cache = {}
    
    def recommend(self, user_id: str, candidate_templates: List[str], top_k: int = 10) -> List[Dict]:
        """协同过滤推荐"""
        # 获取用户历史
        user_history = self.usage_history.get_user_history(user_id, limit=50)
        
        if not user_history:
            # 冷启动：无历史数据，返回热门模板
            return self._popular_templates(candidate_templates, top_k)
        
        # 基于物品的协同过滤（通常效果更好）
        item_cf_scores = self._item_based_cf(user_history, candidate_templates)
        
        # 基于用户的协同过滤
        user_cf_scores = self._user_based_cf(user_id, candidate_templates)
        
        # 混合两种策略
        final_scores = {}
        for template_id in candidate_templates:
            item_score = item_cf_scores.get(template_id, 0)
            user_score = user_cf_scores.get(template_id, 0)
            final_scores[template_id] = 0.7 * item_score + 0.3 * user_score
        
        # 排序并返回
        ranked = sorted(final_scores.items(), key=lambda x: x[1], reverse=True)
        return [
            {
                "template_id": tid,
                "cf_score": score,
                "recommendation_reason": "基于您的使用历史"
            }
            for tid, score in ranked[:top_k]
        ]
    
    def _item_based_cf(self, user_history: List[str], candidates: List[str]) -> Dict[str, float]:
        """基于物品的协同过滤"""
        scores = defaultdict(float)
        
        # 对于用户历史中的每个模板
        for used_template_id in user_history:
            # 找到与该模板相似的模板
            similar_templates = self._get_similar_templates(used_template_id)
            
            for candidate_id in candidates:
                if candidate_id in similar_templates:
                    similarity = similar_templates[candidate_id]
                    scores[candidate_id] += similarity
        
        # 归一化
        if scores:
            max_score = max(scores.values())
            scores = {tid: score / max_score for tid, score in scores.items()}
        
        return scores
    
    def _user_based_cf(self, user_id: str, candidates: List[str]) -> Dict[str, float]:
        """基于用户的协同过滤"""
        scores = defaultdict(float)
        
        # 找到相似用户
        similar_users = self._get_similar_users(user_id)
        
        # 聚合相似用户的偏好
        for similar_user_id, similarity in similar_users.items():
            user_prefs = self.usage_history.get_user_history(similar_user_id)
            for template_id in user_prefs:
                if template_id in candidates:
                    scores[template_id] += similarity
        
        # 归一化
        if scores:
            max_score = max(scores.values())
            scores = {tid: score / max_score for tid, score in scores.items()}
        
        return scores
    
    def _get_similar_templates(self, template_id: str) -> Dict[str, float]:
        """获取相似模板（基于共现矩阵）"""
        if template_id in self.item_similarity_cache:
            return self.item_similarity_cache[template_id]
        
        # 从数据库查询共现数据
        query = """
            SELECT 
                t2.template_id,
                COUNT(*) as co_occurrence,
                COUNT(*) * 1.0 / (
                    SELECT COUNT(*) FROM executions WHERE template_id = %s
                ) as similarity
            FROM executions t1
            JOIN executions t2 ON t1.user_id = t2.user_id
            WHERE t1.template_id = %s AND t2.template_id != %s
            GROUP BY t2.template_id
            ORDER BY similarity DESC
            LIMIT 20
        """
        results = self.db.execute(query, (template_id, template_id, template_id))
        
        similar = {row["template_id"]: row["similarity"] for row in results}
        self.item_similarity_cache[template_id] = similar
        return similar
    
    def _get_similar_users(self, user_id: str) -> Dict[str, float]:
        """获取相似用户（基于余弦相似度）"""
        # 简化实现：基于使用模板的重叠度
        query = """
            WITH user_templates AS (
                SELECT ARRAY_AGG(DISTINCT template_id) as templates
                FROM executions
                WHERE user_id = %s
            ),
            other_users AS (
                SELECT 
                    user_id,
                    ARRAY_AGG(DISTINCT template_id) as templates
                FROM executions
                WHERE user_id != %s
                GROUP BY user_id
            )
            SELECT 
                o.user_id,
                -- Jaccard相似度：交集 / 并集
                array_length(
                    array_intersect(u.templates, o.templates), 1
                ) * 1.0 / 
                array_length(
                    array_union(u.templates, o.templates), 1
                ) as similarity
            FROM user_templates u, other_users o
            WHERE array_length(array_intersect(u.templates, o.templates), 1) > 0
            ORDER BY similarity DESC
            LIMIT 10
        """
        results = self.db.execute(query, (user_id, user_id))
        return {row["user_id"]: row["similarity"] for row in results}
    
    def _popular_templates(self, candidates: List[str], top_k: int) -> List[Dict]:
        """热门模板（冷启动时使用）"""
        query = """
            SELECT 
                template_id,
                COUNT(*) as usage_count,
                AVG(quality_score) as avg_quality
            FROM executions
            WHERE template_id = ANY(%s)
            GROUP BY template_id
            ORDER BY usage_count DESC, avg_quality DESC
            LIMIT %s
        """
        results = self.db.execute(query, (candidates, top_k))
        return [
            {
                "template_id": row["template_id"],
                "cf_score": 0.5,  # 默认分数
                "recommendation_reason": "热门推荐"
            }
            for row in results
        ]
```

---

#### 综合排序与评分

将三层策略的结果进行综合评分和排序：

```python
class TemplateRankingService:
    """模板排序服务"""
    
    def rank(
        self,
        rule_matches: List[Dict],
        semantic_results: List[Dict],
        cf_results: List[Dict],
        template_metadata: Dict[str, Template]
    ) -> List[Dict]:
        """综合排序"""
        # 1. 合并所有候选模板
        all_candidates = self._merge_candidates(
            rule_matches, semantic_results, cf_results
        )
        
        # 2. 计算综合分数
        for candidate in all_candidates:
            template_id = candidate["template_id"]
            template = template_metadata[template_id]
            
            # 规则分数（0-100）
            rule_score = candidate.get("rule_score", 0) / 100.0
            
            # 语义分数（0-1）
            semantic_score = candidate.get("semantic_score", 0)
            
            # 协同过滤分数（0-1）
            cf_score = candidate.get("cf_score", 0)
            
            # 模板质量分数（基于历史执行数据）
            quality_score = self._calculate_quality_score(template)
            
            # 流行度分数
            popularity_score = self._calculate_popularity_score(template)
            
            # 加权综合分数
            final_score = (
                0.35 * rule_score +           # 规则匹配权重最高
                0.25 * semantic_score +        # 语义相似度
                0.20 * cf_score +              # 个性化推荐
                0.15 * quality_score +         # 模板质量
                0.05 * popularity_score        # 流行度
            )
            
            candidate["final_score"] = final_score
            candidate["score_breakdown"] = {
                "rule": rule_score,
                "semantic": semantic_score,
                "collaborative_filtering": cf_score,
                "quality": quality_score,
                "popularity": popularity_score
            }
        
        # 3. 排序
        all_candidates.sort(key=lambda x: x["final_score"], reverse=True)
        
        # 4. 多样性调整（避免推荐过于相似的模板）
        diverse_results = self._diversify(all_candidates, diversity_threshold=0.8)
        
        return diverse_results
    
    def _calculate_quality_score(self, template: Template) -> float:
        """计算模板质量分数"""
        if not template.execution_history:
            return 0.5  # 默认中等质量
        
        # 基于历史执行数据计算
        success_rate = template.execution_history.get("success_rate", 0)
        avg_quality = template.execution_history.get("avg_quality_score", 0) / 10.0
        avg_rating = template.ratings.get("avg_score", 3.0) / 5.0
        
        # 加权平均
        quality_score = 0.4 * success_rate + 0.3 * avg_quality + 0.3 * avg_rating
        return quality_score
    
    def _calculate_popularity_score(self, template: Template) -> float:
        """计算流行度分数"""
        total_runs = template.execution_history.get("total_runs", 0)
        
        # 使用对数变换，避免过度偏向热门模板
        if total_runs == 0:
            return 0.0
        
        # log10(runs + 1) / log10(max_runs + 1)
        max_runs = 10000  # 假设的最大执行次数
        popularity = math.log10(total_runs + 1) / math.log10(max_runs + 1)
        return min(popularity, 1.0)
    
    def _diversify(self, candidates: List[Dict], diversity_threshold: float = 0.8) -> List[Dict]:
        """多样性调整（Maximal Marginal Relevance）"""
        if len(candidates) <= 5:
            return candidates  # 候选较少时不调整
        
        selected = [candidates[0]]  # 选择得分最高的
        remaining = candidates[1:]
        
        while len(selected) < 10 and remaining:
            # 计算每个候选与已选择模板的最大相似度
            for candidate in remaining:
                max_similarity = max(
                    self._template_similarity(candidate, s)
                    for s in selected
                )
                # MMR分数：相关性 - λ * 相似度
                candidate["mmr_score"] = (
                    candidate["final_score"] - 0.3 * max_similarity
                )
            
            # 选择MMR分数最高的
            remaining.sort(key=lambda x: x["mmr_score"], reverse=True)
            selected.append(remaining.pop(0))
        
        return selected
    
    def _template_similarity(self, t1: Dict, t2: Dict) -> float:
        """计算两个模板的相似度"""
        # 简化实现：基于标签重叠度
        tags1 = set(t1.get("tags", []))
        tags2 = set(t2.get("tags", []))
        
        if not tags1 or not tags2:
            return 0.0
        
        jaccard = len(tags1 & tags2) / len(tags1 | tags2)
        return jaccard
```

---

#### 完整推荐流程

```python
class TemplateRecommendationService:
    """模板推荐服务（完整流程）"""
    
    def __init__(self):
        self.rule_matcher = TemplateRuleMatcher()
        self.semantic_search = TemplateSemanticSearchEngine()
        self.collaborative_filter = CollaborativeFilteringRecommender()
        self.ranker = TemplateRankingService()
        self.cache = RedisCache()
    
    async def recommend(
        self,
        query: TemplateQuery,
        user_id: str,
        top_k: int = 10
    ) -> List[Dict]:
        """推荐模板（完整流程）"""
        
        # 1. 检查缓存
        cache_key = self._build_cache_key(query, user_id)
        cached = await self.cache.get(cache_key)
        if cached:
            return cached
        
        # 2. 规则匹配（初步筛选）
        all_templates = await self._fetch_all_templates()
        rule_matches = self.rule_matcher.match(query, all_templates)
        
        # 3. 语义检索（扩展候选集）
        semantic_results = await self.semantic_search.search(query, top_k=50)
        
        # 4. 合并候选集
        candidate_ids = self._merge_candidate_ids(rule_matches, semantic_results)
        
        # 5. 协同过滤（个性化）
        cf_results = self.collaborative_filter.recommend(user_id, candidate_ids)
        
        # 6. 综合排序
        template_metadata = await self._fetch_templates(candidate_ids)
        final_ranking = self.ranker.rank(
            rule_matches,
            semantic_results,
            cf_results,
            template_metadata
        )
        
        # 7. 生成推荐解释
        for item in final_ranking[:top_k]:
            item["explanation"] = self._generate_explanation(item)
        
        # 8. 缓存结果（15分钟）
        await self.cache.set(cache_key, final_ranking[:top_k], ttl=900)
        
        return final_ranking[:top_k]
    
    def _generate_explanation(self, recommendation: Dict) -> str:
        """生成推荐解释"""
        reasons = []
        breakdown = recommendation["score_breakdown"]
        
        if breakdown["rule"] > 0.7:
            reasons.append("精确匹配您的需求（学科、年级、课型）")
        
        if breakdown["semantic"] > 0.7:
            reasons.append("内容高度相关")
        
        if breakdown["collaborative_filtering"] > 0.5:
            reasons.append("基于您的使用历史")
        
        if breakdown["quality"] > 0.8:
            reasons.append("高质量模板（用户评分{:.1f}★）".format(
                recommendation.get("avg_rating", 4.0)
            ))
        
        if breakdown["popularity"] > 0.7:
            reasons.append("热门推荐")
        
        if not reasons:
            reasons.append("系统推荐")
        
        return " | ".join(reasons)
```

---

#### API接口设计

```python
# API端点：POST /v1/templates/recommend

@app.post("/v1/templates/recommend")
async def recommend_templates(request: RecommendRequest, user: User = Depends(get_current_user)):
    """智能推荐模板"""
    
    query = TemplateQuery(
        subject=request.subject,
        grade=request.grade,
        region=request.region,
        lesson_type=request.lesson_type,
        keywords=request.keywords,
        description=request.description,
        tags=request.tags
    )
    
    recommendations = await recommendation_service.recommend(
        query=query,
        user_id=user.id,
        top_k=request.top_k or 10
    )
    
    return {
        "success": True,
        "data": {
            "recommendations": [
                {
                    "template_id": rec["template_id"],
                    "template_name": rec["template_name"],
                    "match_score": rec["final_score"],
                    "explanation": rec["explanation"],
                    "score_breakdown": rec["score_breakdown"],
                    "metadata": {
                        "subject": rec["subject"],
                        "grade_range": rec["grade_range"],
                        "lesson_types": rec["lesson_types"],
                        "tags": rec["tags"]
                    },
                    "statistics": {
                        "total_uses": rec.get("total_runs", 0),
                        "success_rate": rec.get("success_rate", 0),
                        "avg_rating": rec.get("avg_rating", 0)
                    }
                }
                for rec in recommendations
            ],
            "query": {
                "subject": request.subject,
                "grade": request.grade,
                "keywords": request.keywords
            }
        }
    }
```

**请求示例**：
```json
POST /v1/templates/recommend
{
    "subject": "数学",
    "grade": 3,
    "region": "CN",
    "lesson_type": "新授课",
    "keywords": ["分数", "加法"],
    "description": "讲解分数加法的基本概念和计算方法",
    "tags": ["故事化", "互动"],
    "top_k": 10
}
```

**响应示例**：
```json
{
    "success": true,
    "data": {
        "recommendations": [
            {
                "template_id": "tpl_math_fraction_001",
                "template_name": "小学三年级数学-分数加法趣味课堂",
                "match_score": 0.92,
                "explanation": "精确匹配您的需求（学科、年级、课型） | 内容高度相关 | 高质量模板（用户评分4.7★）",
                "score_breakdown": {
                    "rule": 0.95,
                    "semantic": 0.88,
                    "collaborative_filtering": 0.75,
                    "quality": 0.85,
                    "popularity": 0.80
                },
                "metadata": {
                    "subject": "数学",
                    "grade_range": [3, 4],
                    "lesson_types": ["新授课", "练习课"],
                    "tags": ["故事化", "互动", "分数", "基础"]
                },
                "statistics": {
                    "total_uses": 1523,
                    "success_rate": 0.96,
                    "avg_rating": 4.7
                }
            },
            {
                "template_id": "tpl_math_fraction_002",
                "template_name": "分数初探-游戏化教学模板",
                "match_score": 0.85,
                "explanation": "精确匹配您的需求（学科、年级、课型） | 基于您的使用历史",
                "score_breakdown": {
                    "rule": 0.90,
                    "semantic": 0.82,
                    "collaborative_filtering": 0.88,
                    "quality": 0.78,
                    "popularity": 0.65
                },
                "metadata": {
                    "subject": "数学",
                    "grade_range": [2, 4],
                    "lesson_types": ["新授课"],
                    "tags": ["游戏化", "互动", "分数"]
                },
                "statistics": {
                    "total_uses": 892,
                    "success_rate": 0.94,
                    "avg_rating": 4.5
                }
            }
            // ... 更多推荐
        ],
        "query": {
            "subject": "数学",
            "grade": 3,
            "keywords": ["分数", "加法"]
        }
    }
}
```

---

#### 性能优化

**1. 缓存策略**
```python
# 三级缓存
# L1: 用户会话缓存（内存，1小时）
# L2: Redis缓存（15分钟）
# L3: 预计算热门推荐（1天）

class CacheStrategy:
    def get_recommendations(self, query_hash: str):
        # L1
        if query_hash in self.session_cache:
            return self.session_cache[query_hash]
        
        # L2
        redis_result = self.redis.get(f"rec:{query_hash}")
        if redis_result:
            self.session_cache[query_hash] = redis_result
            return redis_result
        
        # L3
        if self._is_popular_query(query_hash):
            precomputed = self.db.get_precomputed(query_hash)
            if precomputed:
                return precomputed
        
        return None
```

**2. 异步处理**
```python
# 耗时操作异步化
async def recommend(self, query, user_id):
    # 并行执行多个策略
    rule_task = asyncio.create_task(self.rule_matcher.match(query))
    semantic_task = asyncio.create_task(self.semantic_search.search(query))
    cf_task = asyncio.create_task(self.collaborative_filter.recommend(user_id))
    
    # 等待所有结果
    rule_results, semantic_results, cf_results = await asyncio.gather(
        rule_task, semantic_task, cf_task
    )
    
    # 合并排序
    return self.ranker.rank(rule_results, semantic_results, cf_results)
```

**3. 索引优化**
```sql
-- 模板表索引
CREATE INDEX idx_template_subject_grade ON templates(subject, grade_range);
CREATE INDEX idx_template_region ON templates(region);
CREATE INDEX idx_template_tags ON templates USING GIN(tags);

-- 全文搜索索引
CREATE INDEX idx_template_fts ON templates USING GIN(to_tsvector('chinese', name || ' ' || description));

-- 向量索引（Milvus中）
-- IVF_FLAT索引，适合中等规模数据集
```

**4. 批量处理**
```python
# 批量推荐（适用于批量生成场景）
async def batch_recommend(self, queries: List[TemplateQuery], user_id: str):
    # 合并相似查询
    clustered = self._cluster_similar_queries(queries)
    
    # 批量检索
    results = []
    for cluster in clustered:
        representative = cluster[0]
        recs = await self.recommend(representative, user_id, top_k=20)
        
        # 为cluster中的每个查询分配推荐
        for query in cluster:
            results.append(self._adjust_for_query(recs, query))
    
    return results
```

---

#### 评估指标

**离线评估**：
- **准确率@K**：推荐的前K个模板中，用户最终选择的比例
- **召回率@K**：用户最终选择的模板是否在推荐的前K个中
- **NDCG@K**：归一化折损累积增益（考虑排序质量）
- **覆盖率**：推荐系统能覆盖多少模板（避免马太效应）

**在线评估**：
- **点击率（CTR）**：用户点击推荐模板的比例
- **转化率**：用户实例化推荐模板的比例
- **用户满意度**：基于使用后评分
- **推荐多样性**：推荐结果的多样性

**A/B测试**：
```python
# 实验配置
experiments = {
    "control": {
        "weights": [0.35, 0.25, 0.20, 0.15, 0.05],  # 原始权重
        "traffic": 0.5
    },
    "variant_a": {
        "weights": [0.40, 0.30, 0.15, 0.10, 0.05],  # 提高规则和语义权重
        "traffic": 0.25
    },
    "variant_b": {
        "weights": [0.30, 0.20, 0.30, 0.15, 0.05],  # 提高协同过滤权重
        "traffic": 0.25
    }
}
```

---

### 3.2.2 自动化工作流构建与优化系统 ⭐ 新增

#### 系统概述

自动化工作流构建系统（Automated Workflow Builder）是一个**AI驱动的工作流设计助手**，能够根据用户的课程需求，从0开始自动设计、构建、测试和优化工作流，而不仅仅是推荐现有模板。

**核心能力**：
- 🤖 **智能组装**：基于组件能力库，自动选择和组装节点
- 🧪 **试运行验证**：执行工作流并生成测试结果
- 📊 **自动评判**：使用评判组件或人工反馈评估结果质量
- 🔄 **迭代优化**：根据评判结果自动调整提示词、增删节点、修改连接
- 💬 **交互式设计**：支持用户与系统对话，逐步完善工作流

**与模板推荐的区别**：

| 维度 | 模板推荐系统 (3.2.1) | 自动化构建系统 (3.2.2) |
|------|---------------------|----------------------|
| **输入** | 需求描述 | 需求描述 |
| **输出** | 现有模板列表 | **新构建的工作流** |
| **创新性** | 低（复用已有） | **高（创造新的）** |
| **适用场景** | 常见课型 | **特殊需求、创新课型** |
| **灵活性** | 受限于模板库 | **完全自定义** |
| **学习成本** | 低 | 中（需理解组件） |

**典型场景**：

```
场景1：全新课型
- 用户需求：设计一个"VR沉浸式化学实验课"
- 模板推荐：❌ 没有匹配的现有模板
- 自动构建：✅ 分析需求 → 选择VR组件 + 实验组件 → 构建新工作流

场景2：特殊教学法
- 用户需求：基于游戏化学习理论设计数学课
- 模板推荐：⚠️ 只能找到部分相关模板
- 自动构建：✅ 根据游戏化理论自动设计多阶段任务流

场景3：跨学科融合
- 用户需求：STEAM课程（科学+艺术融合）
- 模板推荐：❌ 单学科模板无法满足
- 自动构建：✅ 组合多学科组件，构建融合工作流
```

---

#### 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│              Automated Workflow Builder                          │
│              (自动化工作流构建系统)                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                ┌─────────────┴─────────────┐
                │                           │
                ▼                           ▼
    ┌───────────────────────┐   ┌───────────────────────┐
    │  Requirement Analyzer  │   │  Component Registry   │
    │    (需求分析器)        │   │   (组件注册表)        │
    │                        │   │                       │
    │ - 意图识别             │   │ - 组件能力描述        │
    │ - 目标提取             │   │ - 参数schema         │
    │ - 约束识别             │   │ - 使用示例           │
    └───────────────────────┘   └───────────────────────┘
                │                           │
                └─────────────┬─────────────┘
                              ▼
                ┌───────────────────────────┐
                │   Workflow Composer       │
                │   (工作流编排器)          │
                │                           │
                │ - 节点选择                │
                │ - DAG构建                 │
                │ - 依赖分析                │
                └───────────────────────────┘
                              │
                              ▼
                ┌───────────────────────────┐
                │   Trial Executor          │
                │   (试运行执行器)          │
                │                           │
                │ - 执行工作流              │
                │ - 收集输出                │
                └───────────────────────────┘
                              │
                              ▼
                ┌───────────────────────────┐
                │   Quality Evaluator       │
                │   (质量评估器)            │
                │                           │
                │ - 自动评判组件            │
                │ - 人工反馈收集            │
                │ - 问题诊断                │
                └───────────────────────────┘
                              │
                              ▼
                ┌───────────────────────────┐
                │   Workflow Optimizer      │
                │   (工作流优化器)          │
                │                           │
                │ - 提示词调整              │
                │ - 节点增删                │
                │ - 连接修改                │
                └───────────────────────────┘
                              │
                              ▼
                ┌───────────────────────────┐
                │   Iteration Controller    │
                │   (迭代控制器)            │
                │                           │
                │ - 收敛判断                │
                │ - 保存版本                │
                └───────────────────────────┘
```

---

#### 组件注册表（Component Registry）

组件注册表维护所有可用组件的**能力描述**，供工作流编排器查询和选择。

**数据结构**：

```python
@dataclass
class ComponentCapability:
    """组件能力描述"""
    
    component_id: str  # 组件唯一ID，如 "llm_generator_v1"
    component_name: str  # 组件名称，如 "LLM内容生成器"
    category: str  # 组件类别：generator / transformer / evaluator / tool
    
    description: str  # 功能描述（自然语言）
    use_cases: List[str]  # 适用场景列表
    
    # 输入输出规范
    input_schema: Dict[str, Any]  # JSON Schema格式
    output_schema: Dict[str, Any]
    
    # 参数配置
    required_params: List[str]  # 必需参数
    optional_params: Dict[str, Any]  # 可选参数及默认值
    
    # 示例
    examples: List[Dict[str, Any]]  # 使用示例
    
    # 依赖和约束
    dependencies: List[str]  # 依赖的其他组件
    constraints: Dict[str, Any]  # 使用约束（如token限制、成本）
    
    # 元数据
    tags: List[str]  # 标签（用于搜索）
    version: str
    created_at: datetime
```

**示例注册数据**：

```yaml
# config/components/llm_generator.yaml
component_id: llm_generator_v1
component_name: LLM内容生成器
category: generator

description: |
  基于大语言模型生成教学内容，支持多种内容类型（习题、案例、讲义等）

use_cases:
  - 生成数学应用题
  - 生成物理实验案例
  - 生成语文作文范文
  - 生成历史故事材料

input_schema:
  type: object
  required:
    - content_type
    - topic
  properties:
    content_type:
      type: string
      enum: [problem, case, example, explanation]
    topic:
      type: string
      description: 主题或知识点
    difficulty:
      type: string
      enum: [easy, medium, hard]
    constraints:
      type: object
      properties:
        word_count: {type: integer}
        include_solution: {type: boolean}

output_schema:
  type: object
  properties:
    content:
      type: string
      description: 生成的内容
    metadata:
      type: object

required_params:
  - model_name  # 如 "gpt-4"
  - temperature
  - max_tokens

optional_params:
  system_prompt: "你是一名资深教师..."
  few_shot_examples: []

examples:
  - name: 生成初中数学应用题
    input:
      content_type: problem
      topic: 一元二次方程
      difficulty: medium
      constraints:
        word_count: 150
        include_solution: true
    config:
      model_name: gpt-4
      temperature: 0.7
      max_tokens: 500
    expected_output:
      content: "某工厂生产..."

dependencies: []

constraints:
  max_tokens: 4000
  estimated_cost_per_call: 0.05  # USD
  rate_limit: 60  # calls/min

tags:
  - llm
  - generator
  - content_creation
  - math
  - science
  
version: "1.0.0"
created_at: "2024-01-15"
```

**注册表管理类**：

```python
class ComponentRegistry:
    """组件注册表管理"""
    
    def __init__(self, config_dir: Path):
        self.config_dir = config_dir
        self.components: Dict[str, ComponentCapability] = {}
        self._load_all_components()
    
    def _load_all_components(self):
        """加载所有组件配置"""
        for yaml_file in self.config_dir.glob("*.yaml"):
            with open(yaml_file, 'r', encoding='utf-8') as f:
                data = yaml.safe_load(f)
                capability = ComponentCapability(**data)
                self.components[capability.component_id] = capability
    
    def search_by_capability(
        self,
        description: str,
        category: Optional[str] = None,
        tags: Optional[List[str]] = None
    ) -> List[ComponentCapability]:
        """
        基于能力描述搜索组件
        
        使用语义搜索 + 规则过滤
        """
        # 1. 语义搜索：计算描述文本的相似度
        embedding = self._get_embedding(description)
        candidates = []
        
        for comp_id, comp in self.components.items():
            comp_embedding = self._get_embedding(
                f"{comp.component_name} {comp.description}"
            )
            similarity = cosine_similarity(embedding, comp_embedding)
            
            # 2. 规则过滤
            if category and comp.category != category:
                continue
            if tags and not set(tags).intersection(comp.tags):
                continue
            
            candidates.append((similarity, comp))
        
        # 3. 排序返回
        candidates.sort(key=lambda x: x[0], reverse=True)
        return [comp for _, comp in candidates]
    
    def get_component(self, component_id: str) -> Optional[ComponentCapability]:
        """根据ID获取组件"""
        return self.components.get(component_id)
    
    def get_by_category(self, category: str) -> List[ComponentCapability]:
        """获取某类别的所有组件"""
        return [
            comp for comp in self.components.values()
            if comp.category == category
        ]
    
    def validate_workflow(self, workflow: Dict[str, Any]) -> List[str]:
        """
        验证工作流配置的合法性
        
        返回错误列表（空列表表示验证通过）
        """
        errors = []
        
        for node in workflow.get('nodes', []):
            comp_id = node.get('component')
            comp = self.get_component(comp_id)
            
            if not comp:
                errors.append(f"未知组件: {comp_id}")
                continue
            
            # 验证必需参数
            config = node.get('config', {})
            for param in comp.required_params:
                if param not in config:
                    errors.append(
                        f"节点 {node['id']} 缺少必需参数: {param}"
                    )
            
            # 验证输入输出schema（简化示例）
            # TODO: 使用jsonschema库进行完整验证
        
        return errors
```

---

#### 需求分析器（Requirement Analyzer）

需求分析器将用户的**自然语言需求**转化为**结构化的工作流设计约束**。

**输入示例**：

```python
user_requirement = """
我需要设计一个初中数学课，主题是"一元二次方程的应用"。
课程要求：
1. 先用一个真实场景引入（比如商品利润问题）
2. 生成5道应用题，难度要递进
3. 每道题要有详细解析
4. 最后用AI批改学生的作业
5. 必须符合中国教学大纲
"""
```

**输出示例**：

```python
@dataclass
class WorkflowRequirement:
    """工作流需求（结构化）"""
    
    # 基础信息
    subject: str = "数学"
    grade: str = "初中"
    topic: str = "一元二次方程的应用"
    region: str = "中国"
    
    # 教学目标
    learning_objectives: List[str] = field(default_factory=lambda: [
        "理解一元二次方程在实际问题中的应用",
        "掌握建立方程模型的方法"
    ])
    
    # 内容需求
    content_requirements: List[ContentRequirement] = field(default_factory=list)
    
    # 流程需求
    workflow_steps: List[WorkflowStep] = field(default_factory=list)
    
    # 约束条件
    constraints: Dict[str, Any] = field(default_factory=dict)
    
    # 评估标准
    evaluation_criteria: List[str] = field(default_factory=list)

@dataclass
class ContentRequirement:
    """内容需求"""
    content_type: str  # scenario / problem / explanation / evaluation
    quantity: int  # 数量
    attributes: Dict[str, Any]  # 属性要求

@dataclass
class WorkflowStep:
    """流程步骤"""
    step_id: str
    description: str
    required_capabilities: List[str]  # 需要的组件能力
    dependencies: List[str]  # 依赖的前置步骤
```

**需求分析器实现**：

```python
class RequirementAnalyzer:
    """需求分析器"""
    
    def __init__(self, llm_client):
        self.llm_client = llm_client
        self.schema = WorkflowRequirement  # 用于JSON schema生成
    
    async def analyze(self, user_input: str) -> WorkflowRequirement:
        """
        分析用户需求，返回结构化描述
        
        使用LLM进行信息提取和结构化
        """
        # 1. 构建提示词
        prompt = f"""
你是一个教育工作流设计助手。请分析用户的课程需求，提取关键信息。

用户需求：
{user_input}

请以JSON格式输出，包含以下字段：
1. subject: 学科
2. grade: 年级
3. topic: 主题
4. region: 地区
5. learning_objectives: 教学目标（列表）
6. content_requirements: 内容需求（列表，每项包含 content_type, quantity, attributes）
7. workflow_steps: 流程步骤（列表，每项包含 step_id, description, required_capabilities, dependencies）
8. constraints: 约束条件（字典）
9. evaluation_criteria: 评估标准（列表）

示例content_type: scenario, problem, explanation, evaluation
示例required_capabilities: content_generation, quality_check, grading
"""
        
        # 2. 调用LLM
        response = await self.llm_client.complete(
            prompt=prompt,
            response_format={"type": "json_object"},
            temperature=0.3
        )
        
        # 3. 解析JSON
        result = json.loads(response.content)
        
        # 4. 转换为数据类
        requirement = WorkflowRequirement(
            subject=result.get("subject", ""),
            grade=result.get("grade", ""),
            topic=result.get("topic", ""),
            region=result.get("region", "中国"),
            learning_objectives=result.get("learning_objectives", []),
            content_requirements=[
                ContentRequirement(**cr)
                for cr in result.get("content_requirements", [])
            ],
            workflow_steps=[
                WorkflowStep(**ws)
                for ws in result.get("workflow_steps", [])
            ],
            constraints=result.get("constraints", {}),
            evaluation_criteria=result.get("evaluation_criteria", [])
        )
        
        return requirement
    
    async def refine_with_interaction(
        self,
        initial_requirement: WorkflowRequirement,
        conversation_history: List[Dict[str, str]]
    ) -> WorkflowRequirement:
        """
        通过交互式对话逐步完善需求
        
        用户可能一开始描述不清，需要系统提问澄清
        """
        # 检查缺失的关键信息
        missing_fields = self._identify_missing_fields(initial_requirement)
        
        if not missing_fields:
            return initial_requirement
        
        # 生成澄清问题
        questions = self._generate_clarifying_questions(missing_fields)
        
        # 假设这里通过对话获得了用户回答
        # answers = await self._ask_user(questions)
        
        # 更新需求（这里简化，实际需要完整的对话管理）
        # updated_requirement = self._merge_answers(initial_requirement, answers)
        
        return initial_requirement  # 简化返回
    
    def _identify_missing_fields(self, req: WorkflowRequirement) -> List[str]:
        """识别缺失的关键字段"""
        missing = []
        
        if not req.subject:
            missing.append("subject")
        if not req.grade:
            missing.append("grade")
        if not req.topic:
            missing.append("topic")
        if not req.learning_objectives:
            missing.append("learning_objectives")
        if not req.content_requirements:
            missing.append("content_requirements")
        
        return missing
    
    def _generate_clarifying_questions(self, missing_fields: List[str]) -> List[str]:
        """生成澄清问题"""
        questions = []
        
        question_templates = {
            "subject": "请问这是什么学科的课程？（如数学、物理、语文等）",
            "grade": "请问是什么年级的课程？（如小学三年级、初中一年级等）",
            "topic": "请问具体的教学主题是什么？",
            "learning_objectives": "请问这节课的教学目标是什么？学生应该学会什么？",
            "content_requirements": "请问需要包含哪些类型的教学内容？（如案例、习题、讲解等）"
        }
        
        for field in missing_fields:
            if field in question_templates:
                questions.append(question_templates[field])
        
        return questions
```

---

#### 工作流编排器（Workflow Composer）

工作流编排器根据**需求分析结果**和**组件注册表**，自动设计工作流DAG。

**核心算法：基于能力匹配的DAG构建**

```python
class WorkflowComposer:
    """工作流编排器"""
    
    def __init__(
        self,
        registry: ComponentRegistry,
        llm_client,
        max_nodes: int = 20
    ):
        self.registry = registry
        self.llm_client = llm_client
        self.max_nodes = max_nodes
    
    async def compose(
        self,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        根据需求自动构建工作流
        
        Returns:
            完整的工作流配置字典
        """
        # 1. 为每个步骤选择组件
        node_plan = await self._plan_nodes(requirement)
        
        # 2. 构建节点配置
        nodes = await self._build_nodes(node_plan, requirement)
        
        # 3. 分析依赖关系，构建边
        edges = self._build_edges(nodes, requirement.workflow_steps)
        
        # 4. 验证DAG合法性
        self._validate_dag(nodes, edges)
        
        # 5. 组装完整工作流
        workflow = {
            "name": f"{requirement.subject}_{requirement.topic}_workflow",
            "description": f"自动生成：{requirement.subject} - {requirement.topic}",
            "metadata": {
                "subject": requirement.subject,
                "grade": requirement.grade,
                "topic": requirement.topic,
                "auto_generated": True,
                "generated_at": datetime.now().isoformat()
            },
            "nodes": nodes,
            "edges": edges
        }
        
        return workflow
    
    async def _plan_nodes(
        self,
        requirement: WorkflowRequirement
    ) -> List[Dict[str, Any]]:
        """
        规划节点：为每个步骤选择最合适的组件
        
        使用LLM + 组件搜索的混合方法
        """
        node_plan = []
        
        for step in requirement.workflow_steps:
            # 1. 基于required_capabilities搜索候选组件
            candidates = []
            for capability_desc in step.required_capabilities:
                matched = self.registry.search_by_capability(
                    description=capability_desc,
                    tags=None
                )
                candidates.extend(matched[:3])  # 每个能力取top3
            
            # 2. 使用LLM选择最佳组件
            best_component = await self._select_best_component(
                step=step,
                candidates=candidates,
                requirement=requirement
            )
            
            node_plan.append({
                "step": step,
                "component": best_component
            })
        
        return node_plan
    
    async def _select_best_component(
        self,
        step: WorkflowStep,
        candidates: List[ComponentCapability],
        requirement: WorkflowRequirement
    ) -> ComponentCapability:
        """
        使用LLM从候选组件中选择最佳组件
        """
        if not candidates:
            raise ValueError(f"没有找到满足步骤 {step.step_id} 的组件")
        
        if len(candidates) == 1:
            return candidates[0]
        
        # 构建提示词
        candidate_descriptions = "\n".join([
            f"{i+1}. {c.component_name} ({c.component_id}): {c.description}"
            for i, c in enumerate(candidates)
        ])
        
        prompt = f"""
请为以下教学流程步骤选择最合适的组件：

步骤描述：{step.description}
所需能力：{', '.join(step.required_capabilities)}

可选组件：
{candidate_descriptions}

课程背景：
- 学科：{requirement.subject}
- 年级：{requirement.grade}
- 主题：{requirement.topic}

请输出最佳组件的编号（1-{len(candidates)}），只输出数字。
"""
        
        response = await self.llm_client.complete(
            prompt=prompt,
            temperature=0.3,
            max_tokens=10
        )
        
        try:
            choice_idx = int(response.content.strip()) - 1
            return candidates[choice_idx]
        except (ValueError, IndexError):
            # 解析失败，返回第一个候选
            return candidates[0]
    
    async def _build_nodes(
        self,
        node_plan: List[Dict[str, Any]],
        requirement: WorkflowRequirement
    ) -> List[Dict[str, Any]]:
        """
        构建节点配置（包括参数配置）
        """
        nodes = []
        
        for i, plan in enumerate(node_plan):
            step = plan["step"]
            component = plan["component"]
            
            # 生成节点ID
            node_id = f"node_{i+1}_{step.step_id}"
            
            # 生成配置参数
            config = await self._generate_config(
                component=component,
                step=step,
                requirement=requirement
            )
            
            node = {
                "id": node_id,
                "name": step.description,
                "component": component.component_id,
                "config": config,
                "metadata": {
                    "step_id": step.step_id,
                    "auto_generated": True
                }
            }
            
            nodes.append(node)
        
        return nodes
    
    async def _generate_config(
        self,
        component: ComponentCapability,
        step: WorkflowStep,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        为组件生成配置参数（特别是提示词）
        
        使用LLM生成定制化的提示词
        """
        config = {}
        
        # 1. 填充必需参数的默认值
        for param in component.required_params:
            if param in component.optional_params:
                config[param] = component.optional_params[param]
        
        # 2. 如果组件需要提示词，使用LLM生成
        if "prompt" in component.required_params or "system_prompt" in component.required_params:
            prompt = await self._generate_prompt_for_component(
                component=component,
                step=step,
                requirement=requirement
            )
            config["prompt"] = prompt
        
        # 3. 根据需求设置特定参数
        if "difficulty" in component.optional_params:
            # 根据年级推断难度
            if "小学" in requirement.grade:
                config["difficulty"] = "easy"
            elif "初中" in requirement.grade:
                config["difficulty"] = "medium"
            else:
                config["difficulty"] = "hard"
        
        return config
    
    async def _generate_prompt_for_component(
        self,
        component: ComponentCapability,
        step: WorkflowStep,
        requirement: WorkflowRequirement
    ) -> str:
        """
        使用LLM生成组件的提示词
        """
        meta_prompt = f"""
你是一个教学工作流设计专家。请为以下组件生成合适的提示词（prompt）。

组件信息：
- 名称：{component.component_name}
- 功能：{component.description}
- 使用场景：{', '.join(component.use_cases[:3])}

当前步骤：{step.description}

课程信息：
- 学科：{requirement.subject}
- 年级：{requirement.grade}
- 主题：{requirement.topic}
- 教学目标：{', '.join(requirement.learning_objectives)}

请生成一个详细的提示词，使组件能够生成符合要求的高质量内容。
提示词应该包含：角色设定、任务描述、输出要求、质量标准。

直接输出提示词内容，不要额外解释。
"""
        
        response = await self.llm_client.complete(
            prompt=meta_prompt,
            temperature=0.7,
            max_tokens=500
        )
        
        return response.content.strip()
    
    def _build_edges(
        self,
        nodes: List[Dict[str, Any]],
        workflow_steps: List[WorkflowStep]
    ) -> List[Dict[str, Any]]:
        """
        根据步骤依赖关系构建边
        """
        edges = []
        
        # 创建step_id到node_id的映射
        step_to_node = {
            node["metadata"]["step_id"]: node["id"]
            for node in nodes
        }
        
        # 构建边
        for step in workflow_steps:
            target_node_id = step_to_node.get(step.step_id)
            
            if not target_node_id:
                continue
            
            # 为每个依赖创建边
            for dep_step_id in step.dependencies:
                source_node_id = step_to_node.get(dep_step_id)
                
                if source_node_id:
                    edges.append({
                        "from": source_node_id,
                        "to": target_node_id,
                        "type": "sequential"
                    })
        
        return edges
    
    def _validate_dag(
        self,
        nodes: List[Dict[str, Any]],
        edges: List[Dict[str, Any]]
    ):
        """
        验证DAG合法性（无环、连通）
        """
        # 1. 构建邻接表
        graph = {node["id"]: [] for node in nodes}
        for edge in edges:
            graph[edge["from"]].append(edge["to"])
        
        # 2. 检测环（使用DFS）
        visited = set()
        rec_stack = set()
        
        def has_cycle(node_id: str) -> bool:
            visited.add(node_id)
            rec_stack.add(node_id)
            
            for neighbor in graph.get(node_id, []):
                if neighbor not in visited:
                    if has_cycle(neighbor):
                        return True
                elif neighbor in rec_stack:
                    return True
            
            rec_stack.remove(node_id)
            return False
        
        for node in nodes:
            if node["id"] not in visited:
                if has_cycle(node["id"]):
                    raise ValueError("工作流包含环，DAG验证失败")
        
        # 3. 检查连通性（简化：只检查是否有孤立节点）
        all_connected_nodes = set()
        for edge in edges:
            all_connected_nodes.add(edge["from"])
            all_connected_nodes.add(edge["to"])
        
        all_node_ids = {node["id"] for node in nodes}
        isolated = all_node_ids - all_connected_nodes
        
        if len(isolated) > 1:  # 允许有一个起始节点不在边中
            print(f"警告：发现孤立节点 {isolated}")
```

**使用示例**：

```python
# 示例：自动构建工作流
async def example_compose_workflow():
    # 1. 初始化
    registry = ComponentRegistry(Path("config/components"))
    llm_client = OpenAIClient(api_key="...")
    composer = WorkflowComposer(registry, llm_client)
    
    # 2. 定义需求
    requirement = WorkflowRequirement(
        subject="数学",
        grade="初中",
        topic="一元二次方程的应用",
        region="中国",
        learning_objectives=[
            "理解一元二次方程在实际问题中的应用",
            "掌握建立方程模型的方法"
        ],
        workflow_steps=[
            WorkflowStep(
                step_id="step_1",
                description="生成一个真实场景引入案例",
                required_capabilities=["scenario_generation", "real_world_context"],
                dependencies=[]
            ),
            WorkflowStep(
                step_id="step_2",
                description="生成5道应用题，难度递进",
                required_capabilities=["problem_generation", "difficulty_control"],
                dependencies=["step_1"]
            ),
            WorkflowStep(
                step_id="step_3",
                description="为每道题生成详细解析",
                required_capabilities=["solution_generation", "step_by_step_explanation"],
                dependencies=["step_2"]
            ),
            WorkflowStep(
                step_id="step_4",
                description="使用AI批改学生作业",
                required_capabilities=["auto_grading", "feedback_generation"],
                dependencies=["step_3"]
            )
        ]
    )
    
    # 3. 自动构建工作流
    workflow = await composer.compose(requirement)
    
    # 4. 保存工作流
    with open("workflows/auto_generated_math_lesson.yaml", 'w') as f:
        yaml.dump(workflow, f, allow_unicode=True)
    
    print("工作流已自动生成！")
    return workflow
```

---

#### 试运行执行器（Trial Executor）

试运行执行器负责**执行新构建的工作流**，收集输出结果用于评估。

```python
class TrialExecutor:
    """试运行执行器"""
    
    def __init__(self, engine):
        """
        Args:
            engine: WorkflowEngine实例
        """
        self.engine = engine
    
    async def execute_trial(
        self,
        workflow: Dict[str, Any],
        trial_inputs: Optional[Dict[str, Any]] = None
    ) -> TrialResult:
        """
        执行试运行
        
        Args:
            workflow: 工作流配置
            trial_inputs: 试运行输入数据（可选）
        
        Returns:
            试运行结果，包含所有节点输出和执行信息
        """
        # 1. 准备输入数据
        if trial_inputs is None:
            trial_inputs = self._generate_mock_inputs(workflow)
        
        # 2. 创建执行上下文
        context = {
            "mode": "trial",
            "collect_intermediate_outputs": True,  # 收集所有中间结果
            "enable_logging": True
        }
        
        # 3. 执行工作流
        try:
            result = await self.engine.execute(
                workflow=workflow,
                inputs=trial_inputs,
                context=context
            )
            
            # 4. 收集执行信息
            trial_result = TrialResult(
                success=True,
                outputs=result.outputs,
                intermediate_outputs=result.intermediate_outputs,
                node_executions=result.node_executions,
                execution_time=result.execution_time,
                error=None
            )
            
        except Exception as e:
            trial_result = TrialResult(
                success=False,
                outputs={},
                intermediate_outputs={},
                node_executions=[],
                execution_time=0,
                error=str(e)
            )
        
        return trial_result
    
    def _generate_mock_inputs(self, workflow: Dict[str, Any]) -> Dict[str, Any]:
        """
        为试运行生成模拟输入数据
        """
        # 简化实现：返回空输入
        # 实际应该根据workflow的input_schema生成
        return {}


@dataclass
class TrialResult:
    """试运行结果"""
    success: bool
    outputs: Dict[str, Any]
    intermediate_outputs: Dict[str, Any]  # 每个节点的输出
    node_executions: List[Dict[str, Any]]  # 节点执行详情
    execution_time: float
    error: Optional[str] = None
```

---

#### 质量评估器（Quality Evaluator）

质量评估器分析试运行结果，诊断问题并给出优化建议。

**评估策略：自动评判 + 人工反馈**

```python
class QualityEvaluator:
    """质量评估器"""
    
    def __init__(
        self,
        auto_evaluator_component: Optional[str] = None,
        llm_client=None
    ):
        """
        Args:
            auto_evaluator_component: 自动评判组件ID（如 "rubric_evaluator"）
            llm_client: LLM客户端（用于基于LLM的评估）
        """
        self.auto_evaluator = auto_evaluator_component
        self.llm_client = llm_client
    
    async def evaluate(
        self,
        trial_result: TrialResult,
        requirement: WorkflowRequirement,
        evaluation_mode: str = "auto"  # auto / human / hybrid
    ) -> EvaluationReport:
        """
        评估试运行结果
        
        Args:
            trial_result: 试运行结果
            requirement: 原始需求
            evaluation_mode: 评估模式
        
        Returns:
            评估报告
        """
        if not trial_result.success:
            # 执行失败，直接返回错误报告
            return EvaluationReport(
                overall_score=0.0,
                passed=False,
                issues=[
                    Issue(
                        severity="critical",
                        category="execution_error",
                        description=trial_result.error,
                        affected_nodes=[],
                        suggestions=["检查节点配置", "验证组件参数"]
                    )
                ],
                strengths=[],
                suggestions=["修复执行错误后重新运行"]
            )
        
        # 执行成功，进行质量评估
        if evaluation_mode == "auto":
            return await self._auto_evaluate(trial_result, requirement)
        elif evaluation_mode == "human":
            return await self._human_evaluate(trial_result, requirement)
        else:  # hybrid
            auto_report = await self._auto_evaluate(trial_result, requirement)
            human_report = await self._human_evaluate(trial_result, requirement)
            return self._merge_reports(auto_report, human_report)
    
    async def _auto_evaluate(
        self,
        trial_result: TrialResult,
        requirement: WorkflowRequirement
    ) -> EvaluationReport:
        """
        自动评估（使用LLM）
        """
        # 1. 构建评估提示词
        evaluation_prompt = f"""
你是一名资深教育专家。请评估以下教学内容生成结果的质量。

课程需求：
- 学科：{requirement.subject}
- 年级：{requirement.grade}
- 主题：{requirement.topic}
- 教学目标：{', '.join(requirement.learning_objectives)}

生成结果：
{json.dumps(trial_result.outputs, ensure_ascii=False, indent=2)}

评估维度：
1. 内容准确性（是否符合知识点要求）
2. 难度适配性（是否符合年级水平）
3. 逻辑连贯性（内容之间的逻辑关系）
4. 教学有效性（是否能达成教学目标）
5. 语言表达（是否清晰易懂）

请以JSON格式输出评估结果：
{{
  "overall_score": 0.0-1.0,
  "dimension_scores": {{
    "accuracy": 0.0-1.0,
    "difficulty": 0.0-1.0,
    "coherence": 0.0-1.0,
    "effectiveness": 0.0-1.0,
    "language": 0.0-1.0
  }},
  "issues": [
    {{
      "severity": "critical|major|minor",
      "category": "content|structure|language",
      "description": "问题描述",
      "location": "具体位置",
      "suggestion": "改进建议"
    }}
  ],
  "strengths": ["优点1", "优点2"],
  "overall_feedback": "总体评价"
}}
"""
        
        # 2. 调用LLM评估
        response = await self.llm_client.complete(
            prompt=evaluation_prompt,
            response_format={"type": "json_object"},
            temperature=0.3
        )
        
        # 3. 解析评估结果
        eval_data = json.loads(response.content)
        
        # 4. 转换为报告对象
        issues = [
            Issue(
                severity=issue["severity"],
                category=issue["category"],
                description=issue["description"],
                affected_nodes=[issue.get("location", "")],
                suggestions=[issue["suggestion"]]
            )
            for issue in eval_data.get("issues", [])
        ]
        
        report = EvaluationReport(
            overall_score=eval_data["overall_score"],
            dimension_scores=eval_data.get("dimension_scores", {}),
            passed=eval_data["overall_score"] >= 0.7,
            issues=issues,
            strengths=eval_data.get("strengths", []),
            suggestions=self._generate_optimization_suggestions(issues)
        )
        
        return report
    
    async def _human_evaluate(
        self,
        trial_result: TrialResult,
        requirement: WorkflowRequirement
    ) -> EvaluationReport:
        """
        人工评估（交互式）
        
        实际应该提供Web界面或CLI界面让教师评估
        这里简化处理
        """
        print("=== 人工评估 ===")
        print(f"课程：{requirement.subject} - {requirement.topic}")
        print(f"输出结果：\n{json.dumps(trial_result.outputs, ensure_ascii=False, indent=2)}")
        
        # 简化：假设自动通过
        return EvaluationReport(
            overall_score=0.8,
            passed=True,
            issues=[],
            strengths=["人工确认内容质量良好"],
            suggestions=[]
        )
    
    def _merge_reports(
        self,
        auto_report: EvaluationReport,
        human_report: EvaluationReport
    ) -> EvaluationReport:
        """合并自动评估和人工评估结果"""
        # 简化：取平均分
        merged_score = (auto_report.overall_score + human_report.overall_score) / 2
        
        return EvaluationReport(
            overall_score=merged_score,
            passed=merged_score >= 0.7,
            issues=auto_report.issues + human_report.issues,
            strengths=auto_report.strengths + human_report.strengths,
            suggestions=auto_report.suggestions + human_report.suggestions
        )
    
    def _generate_optimization_suggestions(
        self,
        issues: List[Issue]
    ) -> List[str]:
        """
        根据问题列表生成优化建议
        """
        suggestions = []
        
        # 按严重程度分组
        critical_issues = [i for i in issues if i.severity == "critical"]
        major_issues = [i for i in issues if i.severity == "major"]
        
        if critical_issues:
            suggestions.append(f"发现{len(critical_issues)}个严重问题，建议优先修复")
        
        if major_issues:
            suggestions.append(f"发现{len(major_issues)}个重要问题，建议改进")
        
        # 按类别提供建议
        content_issues = [i for i in issues if i.category == "content"]
        if content_issues:
            suggestions.append("内容准确性需要改进，建议调整提示词或更换组件")
        
        structure_issues = [i for i in issues if i.category == "structure"]
        if structure_issues:
            suggestions.append("工作流结构需要优化，建议调整节点连接或增加中间步骤")
        
        return suggestions


@dataclass
class EvaluationReport:
    """评估报告"""
    overall_score: float  # 0.0-1.0
    passed: bool
    issues: List['Issue']
    strengths: List[str]
    suggestions: List[str]
    dimension_scores: Dict[str, float] = field(default_factory=dict)


@dataclass
class Issue:
    """问题描述"""
    severity: str  # critical / major / minor
    category: str  # content / structure / language / performance
    description: str
    affected_nodes: List[str]
    suggestions: List[str]
```

---

#### 工作流优化器（Workflow Optimizer）

工作流优化器根据**评估报告**自动调整工作流配置，进行迭代优化。

**优化策略**：

```python
class WorkflowOptimizer:
    """工作流优化器"""
    
    def __init__(self, llm_client, registry: ComponentRegistry):
        self.llm_client = llm_client
        self.registry = registry
    
    async def optimize(
        self,
        workflow: Dict[str, Any],
        evaluation_report: EvaluationReport,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        根据评估报告优化工作流
        
        优化策略：
        1. 调整提示词（针对内容质量问题）
        2. 修改组件参数（针对输出格式问题）
        3. 更换组件（针对能力不足问题）
        4. 增加节点（针对流程缺失问题）
        5. 删除节点（针对冗余问题）
        6. 调整连接（针对逻辑问题）
        """
        if evaluation_report.passed:
            # 评估通过，无需优化
            return workflow
        
        # 复制工作流（避免修改原始版本）
        optimized_workflow = copy.deepcopy(workflow)
        
        # 按严重程度处理问题
        critical_issues = [i for i in evaluation_report.issues if i.severity == "critical"]
        major_issues = [i for i in evaluation_report.issues if i.severity == "major"]
        
        # 1. 优先处理严重问题
        for issue in critical_issues:
            optimized_workflow = await self._fix_issue(
                workflow=optimized_workflow,
                issue=issue,
                requirement=requirement
            )
        
        # 2. 处理重要问题
        for issue in major_issues:
            optimized_workflow = await self._fix_issue(
                workflow=optimized_workflow,
                issue=issue,
                requirement=requirement
            )
        
        return optimized_workflow
    
    async def _fix_issue(
        self,
        workflow: Dict[str, Any],
        issue: Issue,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        修复单个问题
        """
        if issue.category == "content":
            # 内容问题 → 调整提示词
            return await self._optimize_prompts(workflow, issue, requirement)
        
        elif issue.category == "structure":
            # 结构问题 → 调整节点或连接
            return await self._optimize_structure(workflow, issue, requirement)
        
        elif issue.category == "language":
            # 语言问题 → 调整语言风格参数
            return await self._optimize_language(workflow, issue, requirement)
        
        else:
            # 其他问题，返回原workflow
            return workflow
    
    async def _optimize_prompts(
        self,
        workflow: Dict[str, Any],
        issue: Issue,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        优化提示词
        """
        affected_nodes = issue.affected_nodes
        
        for node in workflow["nodes"]:
            if node["id"] in affected_nodes or not affected_nodes:
                # 找到需要优化的节点
                current_prompt = node["config"].get("prompt", "")
                
                # 使用LLM改进提示词
                improved_prompt = await self._improve_prompt(
                    current_prompt=current_prompt,
                    issue_description=issue.description,
                    suggestions=issue.suggestions,
                    requirement=requirement
                )
                
                node["config"]["prompt"] = improved_prompt
        
        return workflow
    
    async def _improve_prompt(
        self,
        current_prompt: str,
        issue_description: str,
        suggestions: List[str],
        requirement: WorkflowRequirement
    ) -> str:
        """
        使用LLM改进提示词
        """
        meta_prompt = f"""
你是一个提示词优化专家。请改进以下提示词，解决识别出的问题。

当前提示词：
```
{current_prompt}
```

问题描述：
{issue_description}

改进建议：
{'; '.join(suggestions)}

课程背景：
- 学科：{requirement.subject}
- 年级：{requirement.grade}
- 主题：{requirement.topic}

请输出改进后的提示词（直接输出提示词内容，不要额外解释）：
"""
        
        response = await self.llm_client.complete(
            prompt=meta_prompt,
            temperature=0.7,
            max_tokens=600
        )
        
        return response.content.strip()
    
    async def _optimize_structure(
        self,
        workflow: Dict[str, Any],
        issue: Issue,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        优化工作流结构
        
        可能的操作：
        - 增加节点（如需要额外的质量检查步骤）
        - 删除节点（如存在冗余）
        - 调整连接（如逻辑顺序不对）
        """
        # 简化实现：根据建议决定操作
        if "增加" in issue.description or "添加" in issue.description:
            # 需要增加节点
            return await self._add_node(workflow, issue, requirement)
        
        elif "删除" in issue.description or "移除" in issue.description:
            # 需要删除节点
            return self._remove_node(workflow, issue)
        
        else:
            # 其他结构调整
            return workflow
    
    async def _add_node(
        self,
        workflow: Dict[str, Any],
        issue: Issue,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        添加新节点
        """
        # 使用LLM确定需要添加什么类型的节点
        analysis_prompt = f"""
问题描述：{issue.description}
建议：{'; '.join(issue.suggestions)}

请确定需要添加什么类型的组件。
输出格式：组件类别（generator/transformer/evaluator/tool），一个词。
"""
        
        response = await self.llm_client.complete(
            prompt=analysis_prompt,
            temperature=0.3,
            max_tokens=20
        )
        
        category = response.content.strip().lower()
        
        # 搜索该类别的组件
        candidates = self.registry.get_by_category(category)
        if not candidates:
            return workflow
        
        # 选择第一个候选组件
        new_component = candidates[0]
        
        # 生成新节点
        new_node_id = f"node_optimized_{len(workflow['nodes'])+1}"
        new_node = {
            "id": new_node_id,
            "name": f"优化添加：{new_component.component_name}",
            "component": new_component.component_id,
            "config": {},
            "metadata": {
                "added_by_optimization": True
            }
        }
        
        workflow["nodes"].append(new_node)
        
        # TODO: 添加相应的边（连接到合适的位置）
        
        return workflow
    
    def _remove_node(
        self,
        workflow: Dict[str, Any],
        issue: Issue
    ) -> Dict[str, Any]:
        """
        删除节点
        """
        for node_id in issue.affected_nodes:
            # 删除节点
            workflow["nodes"] = [
                n for n in workflow["nodes"]
                if n["id"] != node_id
            ]
            
            # 删除相关的边
            workflow["edges"] = [
                e for e in workflow["edges"]
                if e["from"] != node_id and e["to"] != node_id
            ]
        
        return workflow
    
    async def _optimize_language(
        self,
        workflow: Dict[str, Any],
        issue: Issue,
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        优化语言风格
        """
        # 在提示词中添加语言风格要求
        style_instruction = "\n\n语言要求：表达清晰简洁，符合学生年龄特点。"
        
        for node in workflow["nodes"]:
            if "prompt" in node["config"]:
                current_prompt = node["config"]["prompt"]
                if style_instruction not in current_prompt:
                    node["config"]["prompt"] = current_prompt + style_instruction
        
        return workflow
```

---

#### 迭代控制器（Iteration Controller）

迭代控制器管理整个**构建-试运行-评估-优化**的循环过程。

```python
class IterationController:
    """迭代控制器"""
    
    def __init__(
        self,
        composer: WorkflowComposer,
        trial_executor: TrialExecutor,
        evaluator: QualityEvaluator,
        optimizer: WorkflowOptimizer,
        max_iterations: int = 5
    ):
        self.composer = composer
        self.trial_executor = trial_executor
        self.evaluator = evaluator
        self.optimizer = optimizer
        self.max_iterations = max_iterations
    
    async def build_and_optimize(
        self,
        requirement: WorkflowRequirement,
        evaluation_mode: str = "auto",
        convergence_threshold: float = 0.8
    ) -> IterationResult:
        """
        自动构建并迭代优化工作流
        
        Args:
            requirement: 用户需求
            evaluation_mode: 评估模式 (auto/human/hybrid)
            convergence_threshold: 收敛阈值（评分达到此值即停止迭代）
        
        Returns:
            迭代结果，包含最终工作流和历史记录
        """
        iteration_history = []
        
        # 1. 初始构建
        print("🚀 开始构建初始工作流...")
        current_workflow = await self.composer.compose(requirement)
        
        for iteration in range(self.max_iterations):
            print(f"\n📍 第 {iteration + 1} 轮迭代")
            
            # 2. 试运行
            print("  🧪 执行试运行...")
            trial_result = await self.trial_executor.execute_trial(current_workflow)
            
            # 3. 评估
            print("  📊 评估结果质量...")
            evaluation = await self.evaluator.evaluate(
                trial_result=trial_result,
                requirement=requirement,
                evaluation_mode=evaluation_mode
            )
            
            print(f"  得分: {evaluation.overall_score:.2f}")
            
            # 4. 记录历史
            iteration_record = IterationRecord(
                iteration=iteration + 1,
                workflow=copy.deepcopy(current_workflow),
                trial_result=trial_result,
                evaluation=evaluation
            )
            iteration_history.append(iteration_record)
            
            # 5. 判断是否收敛
            if evaluation.passed and evaluation.overall_score >= convergence_threshold:
                print("  ✅ 质量达标，迭代结束")
                break
            
            if iteration == self.max_iterations - 1:
                print("  ⚠️  达到最大迭代次数")
                break
            
            # 6. 优化
            print("  🔧 根据评估结果优化工作流...")
            current_workflow = await self.optimizer.optimize(
                workflow=current_workflow,
                evaluation_report=evaluation,
                requirement=requirement
            )
        
        # 7. 返回最终结果
        final_record = iteration_history[-1]
        
        return IterationResult(
            final_workflow=final_record.workflow,
            final_score=final_record.evaluation.overall_score,
            total_iterations=len(iteration_history),
            converged=final_record.evaluation.passed,
            history=iteration_history
        )
    
    def get_best_workflow(self, result: IterationResult) -> Dict[str, Any]:
        """
        从迭代历史中选择得分最高的工作流
        """
        best_record = max(
            result.history,
            key=lambda r: r.evaluation.overall_score
        )
        return best_record.workflow
    
    async def interactive_refine(
        self,
        workflow: Dict[str, Any],
        requirement: WorkflowRequirement
    ) -> Dict[str, Any]:
        """
        交互式优化工作流
        
        用户可以手动调整，系统辅助验证
        """
        print("🎨 进入交互式优化模式")
        print("可用命令：")
        print("  1. adjust_prompt <node_id> - 调整节点提示词")
        print("  2. add_node - 添加新节点")
        print("  3. remove_node <node_id> - 删除节点")
        print("  4. test - 运行测试")
        print("  5. done - 完成优化")
        
        current_workflow = copy.deepcopy(workflow)
        
        # 简化实现：这里只是示例框架
        # 实际需要完整的交互界面（CLI或Web）
        
        return current_workflow


@dataclass
class IterationRecord:
    """单轮迭代记录"""
    iteration: int
    workflow: Dict[str, Any]
    trial_result: TrialResult
    evaluation: EvaluationReport


@dataclass
class IterationResult:
    """迭代结果"""
    final_workflow: Dict[str, Any]
    final_score: float
    total_iterations: int
    converged: bool
    history: List[IterationRecord]
```

---

#### 完整使用示例

```python
async def main():
    """完整的自动化工作流构建流程"""
    
    # 1. 初始化所有组件
    registry = ComponentRegistry(Path("config/components"))
    llm_client = OpenAIClient(api_key=os.getenv("OPENAI_API_KEY"))
    
    # 2. 创建各个子系统
    requirement_analyzer = RequirementAnalyzer(llm_client)
    composer = WorkflowComposer(registry, llm_client)
    trial_executor = TrialExecutor(WorkflowEngine())
    evaluator = QualityEvaluator(llm_client=llm_client)
    optimizer = WorkflowOptimizer(llm_client, registry)
    controller = IterationController(
        composer=composer,
        trial_executor=trial_executor,
        evaluator=evaluator,
        optimizer=optimizer,
        max_iterations=5
    )
    
    # 3. 用户输入需求
    user_input = """
    我需要设计一个高中物理课，主题是"牛顿第二定律"。
    要求：
    1. 先通过一个日常生活例子引入（比如推购物车）
    2. 生成3个计算题，难度从简单到复杂
    3. 每道题要有详细的解题步骤
    4. 最后生成一个实验设计任务
    """
    
    # 4. 分析需求
    print("📝 分析用户需求...")
    requirement = await requirement_analyzer.analyze(user_input)
    
    # 5. 自动构建并优化
    print("\n🤖 开始自动构建工作流...")
    result = await controller.build_and_optimize(
        requirement=requirement,
        evaluation_mode="auto",
        convergence_threshold=0.8
    )
    
    # 6. 保存最终工作流
    final_workflow = result.final_workflow
    output_path = Path(f"workflows/auto_{requirement.subject}_{requirement.topic}.yaml")
    
    with open(output_path, 'w', encoding='utf-8') as f:
        yaml.dump(final_workflow, f, allow_unicode=True)
    
    # 7. 输出总结
    print("\n" + "="*50)
    print("✨ 工作流构建完成！")
    print(f"📊 最终得分: {result.final_score:.2f}")
    print(f"🔄 迭代次数: {result.total_iterations}")
    print(f"✅ 是否收敛: {'是' if result.converged else '否'}")
    print(f"💾 已保存至: {output_path}")
    
    # 8. 显示迭代历史
    print("\n📈 迭代历史：")
    for record in result.history:
        print(f"  第{record.iteration}轮: 得分 {record.evaluation.overall_score:.2f}")
    
    return final_workflow


if __name__ == "__main__":
    asyncio.run(main())
```

**输出示例**：

```
📝 分析用户需求...
  识别学科: 物理
  识别年级: 高中
  识别主题: 牛顿第二定律

🤖 开始自动构建工作流...
🚀 开始构建初始工作流...
  选择组件: scenario_generator (生成引入案例)
  选择组件: problem_generator (生成计算题)
  选择组件: solution_generator (生成解题步骤)
  选择组件: experiment_designer (设计实验)

📍 第 1 轮迭代
  🧪 执行试运行...
  📊 评估结果质量...
  得分: 0.65
  🔧 根据评估结果优化工作流...
    优化提示词: node_2 (提高难度梯度控制)

📍 第 2 轮迭代
  🧪 执行试运行...
  📊 评估结果质量...
  得分: 0.78
  🔧 根据评估结果优化工作流...
    优化提示词: node_3 (增强解题步骤详细度)

📍 第 3 轮迭代
  🧪 执行试运行...
  📊 评估结果质量...
  得分: 0.85
  ✅ 质量达标，迭代结束

==================================================
✨ 工作流构建完成！
📊 最终得分: 0.85
🔄 迭代次数: 3
✅ 是否收敛: 是
💾 已保存至: workflows/auto_物理_牛顿第二定律.yaml

📈 迭代历史：
  第1轮: 得分 0.65
  第2轮: 得分 0.78
  第3轮: 得分 0.85
```

---

#### API接口设计

```python
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

router = APIRouter(prefix="/api/v1/workflow/auto-builder")


class BuildRequest(BaseModel):
    """构建请求"""
    user_requirement: str
    evaluation_mode: str = "auto"  # auto / human / hybrid
    max_iterations: int = 5
    convergence_threshold: float = 0.8


class BuildResponse(BaseModel):
    """构建响应"""
    workflow_id: str
    final_score: float
    total_iterations: int
    converged: bool
    workflow: Dict[str, Any]


@router.post("/build", response_model=BuildResponse)
async def build_workflow(request: BuildRequest):
    """
    自动构建工作流
    
    POST /api/v1/workflow/auto-builder/build
    
    Request Body:
    {
      "user_requirement": "设计一个初中数学课，主题是一元二次方程...",
      "evaluation_mode": "auto",
      "max_iterations": 5,
      "convergence_threshold": 0.8
    }
    
    Response:
    {
      "workflow_id": "wf_auto_20240115_001",
      "final_score": 0.85,
      "total_iterations": 3,
      "converged": true,
      "workflow": { ... }
    }
    """
    try:
        # 初始化控制器（实际应该使用依赖注入）
        controller = get_iteration_controller()
        analyzer = get_requirement_analyzer()
        
        # 分析需求
        requirement = await analyzer.analyze(request.user_requirement)
        
        # 构建并优化
        result = await controller.build_and_optimize(
            requirement=requirement,
            evaluation_mode=request.evaluation_mode,
            convergence_threshold=request.convergence_threshold
        )
        
        # 保存工作流
        workflow_id = f"wf_auto_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        save_workflow(workflow_id, result.final_workflow)
        
        return BuildResponse(
            workflow_id=workflow_id,
            final_score=result.final_score,
            total_iterations=result.total_iterations,
            converged=result.converged,
            workflow=result.final_workflow
        )
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/iterations/{workflow_id}")
async def get_iteration_history(workflow_id: str):
    """
    获取迭代历史
    
    GET /api/v1/workflow/auto-builder/iterations/wf_auto_20240115_001
    
    Response:
    {
      "workflow_id": "wf_auto_20240115_001",
      "iterations": [
        {
          "iteration": 1,
          "score": 0.65,
          "issues": [...],
          "optimizations": [...]
        },
        ...
      ]
    }
    """
    # 从数据库或缓存中获取迭代历史
    history = load_iteration_history(workflow_id)
    
    return {
        "workflow_id": workflow_id,
        "iterations": history
    }
```

---

#### 性能优化

**1. 并行试运行**：

```python
async def parallel_trial_runs(
    workflows: List[Dict[str, Any]],
    executor: TrialExecutor
) -> List[TrialResult]:
    """并行执行多个工作流的试运行"""
    tasks = [
        executor.execute_trial(wf)
        for wf in workflows
    ]
    results = await asyncio.gather(*tasks)
    return results
```

**2. 缓存组件能力**：

```python
# 使用Redis缓存组件搜索结果
@redis_cache(ttl=3600)
def search_components(description: str, category: str):
    return registry.search_by_capability(description, category)
```

**3. 增量优化**：

```python
# 只优化评分下降的节点
def incremental_optimize(workflow, prev_eval, curr_eval):
    """只优化出现新问题的节点"""
    new_issues = [
        i for i in curr_eval.issues
        if i not in prev_eval.issues
    ]
    # 只针对新问题优化
```

---

### 3.3 模板本地化适配器设计 ⭐ 新增

#### 3.3.1 整体架构

模板本地化适配器（Template Localizer）负责将现有课程模板快速适配到新的区域市场，实现**"快速复制 + 深度适配"**的平衡。

```
┌─────────────────────────────────────────────────────────────┐
│            Template Localizer（模板本地化适配器）            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Curriculum Standard Mapper                    │   │
│  │        (课程标准映射器)                               │   │
│  │  ├─ 标准知识库 (NGSS/Common Core/新课标)            │   │
│  │  ├─ LLM语义映射                                      │   │
│  │  └─ 置信度评分 → 人工审核                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                           ↓                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Localization Translator                       │   │
│  │        (语言本地化引擎)                               │   │
│  │  ├─ 术语表构建                                        │   │
│  │  ├─ 教学语境翻译 (GPT-4/Claude)                      │   │
│  │  ├─ 术语一致性检查                                   │   │
│  │  └─ 语言风格调整                                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                           ↓                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Cultural Adapter                              │   │
│  │        (文化适配引擎)                                 │   │
│  │  ├─ 文化元素识别 (NER + 规则)                        │   │
│  │  ├─ 故事世界观重构                                   │   │
│  │  ├─ 案例/比喻本地化                                  │   │
│  │  └─ 敏感内容审查                                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                           ↓                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Pedagogy Adapter                              │   │
│  │        (教学法转换器)                                 │   │
│  │  ├─ 课型比例调整                                      │   │
│  │  ├─ 教学流程模板替换                                  │   │
│  │  ├─ 评估方式转换                                      │   │
│  │  └─ 互动密度调整                                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                           ↓                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Compliance Checker                            │   │
│  │        (合规性检查器)                                 │   │
│  │  ├─ COPPA/FERPA/GDPR检查                             │   │
│  │  ├─ ADA无障碍标准                                     │   │
│  │  ├─ 自动修复 (如生成alt文本)                         │   │
│  │  └─ 人工审核工作流                                    │   │
│  └─────────────────────────────────────────────────────┘   │
│                           ↓                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │        Quality Assessor                              │   │
│  │        (质量评估器)                                   │   │
│  │  ├─ 多维度评分 (结构/语言/文化/对齐/合规)            │   │
│  │  ├─ 版本对比                                          │   │
│  │  └─ 改进建议生成                                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

#### 3.3.2 核心组件设计

**A. 课程标准映射器 (Curriculum Standard Mapper)**

```python
class CurriculumStandardMapper:
    """课程标准映射器"""
    
    def __init__(self):
        self.standard_kb = StandardKnowledgeBase()
        self.llm = LLMClient()
        self.vector_db = MilvusClient()
    
    async def map_standards(
        self,
        source_standards: List[str],
        source_system: str,
        target_system: str
    ) -> List[StandardMapping]:
        """
        映射课程标准
        
        策略：
        1. 向量检索相似标准（基于语义）
        2. LLM理解分析差异
        3. 专家知识库补充
        4. 置信度评分
        """
        
        mappings = []
        
        for source_std in source_standards:
            # 1. 向量检索候选目标标准
            candidates = await self._retrieve_candidates(
                source_std,
                target_system
            )
            
            # 2. LLM精细对齐
            mapping = await self._llm_alignment(
                source_std,
                candidates,
                source_system,
                target_system
            )
            
            # 3. 置信度评分
            mapping.confidence_score = self._calculate_confidence(mapping)
            
            # 4. 标记需审核项
            if mapping.confidence_score < 0.7:
                mapping.requires_review = True
                mapping.review_reason = "Low confidence mapping"
            
            mappings.append(mapping)
        
        return mappings
    
    async def _llm_alignment(
        self,
        source_std: str,
        candidates: List[str],
        source_system: str,
        target_system: str
    ) -> StandardMapping:
        """使用LLM进行精细对齐"""
        
        prompt = f"""
        你是教育标准专家。请将以下{source_system}标准映射到{target_system}:
        
        源标准：{source_std}
        
        候选目标标准：
        {chr(10).join(f"{i+1}. {c}" for i, c in enumerate(candidates))}
        
        请分析：
        1. 最匹配的目标标准是哪一个？为什么？
        2. 知识点覆盖有何差异？（新增/缺失）
        3. 认知层级是否匹配？（基于布鲁姆分类法）
        4. 教学要求有何不同？
        
        输出JSON格式：
        {{
            "target_standard": "...",
            "match_reason": "...",
            "knowledge_gaps": ["...", "..."],
            "new_requirements": ["...", "..."],
            "bloom_level_diff": "...",
            "pedagogical_diff": "..."
        }}
        """
        
        result = await self.llm.generate(prompt)
        
        return StandardMapping(
            source=source_std,
            target=result["target_standard"],
            match_reason=result["match_reason"],
            knowledge_gaps=result["knowledge_gaps"],
            new_requirements=result["new_requirements"],
            differences={
                "bloom_level": result["bloom_level_diff"],
                "pedagogy": result["pedagogical_diff"]
            }
        )


**知识库结构**：

```sql
-- standards库
CREATE TABLE curriculum_standards (
    id UUID PRIMARY KEY,
    system VARCHAR(50),  -- 'US_NGSS', 'CN_NewCurriculum_2022', etc.
    code VARCHAR(100),   -- 'MS-PS1-1', '初中化学-2.1', etc.
    grade_level VARCHAR(20),
    subject VARCHAR(50),
    content TEXT,
    embedding VECTOR(1536),  -- 语义向量
    bloom_level VARCHAR(20), -- Remember/Understand/Apply/Analyze/Evaluate/Create
    
    metadata JSONB  -- 扩展信息
);

-- 标准映射表（人工校准）
CREATE TABLE standard_mappings (
    id UUID PRIMARY KEY,
    source_system VARCHAR(50),
    target_system VARCHAR(50),
    source_code VARCHAR(100),
    target_code VARCHAR(100),
    
    mapping_type VARCHAR(20),  -- 'exact', 'partial', 'none'
    confidence FLOAT,
    
    verified BOOLEAN DEFAULT FALSE,
    verified_by VARCHAR(100),
    verified_at TIMESTAMP,
    
    notes TEXT
);
```

---

**B. 语言本地化引擎 (Localization Translator)**

```python
class LocalizationTranslator:
    """语言本地化引擎"""
    
    def __init__(self):
        self.llm = LLMClient(model="gpt-4")
        self.glossary_builder = GlossaryBuilder()
        self.style_analyzer = StyleAnalyzer()
    
    async def translate_template(
        self,
        template: WorkflowTemplate,
        source_lang: str,
        target_lang: str
    ) -> WorkflowTemplate:
        """翻译模板（保持教学语境）"""
        
        # 1. 构建术语表
        glossary = await self.glossary_builder.build(
            template=template,
            source_lang=source_lang,
            target_lang=target_lang,
            domain="chemistry",
            grade="9"
        )
        
        # 2. 提取所有文本
        texts = self._extract_translatable_texts(template)
        
        # 3. 批量翻译
        translations = await self._translate_batch(
            texts=texts,
            source_lang=source_lang,
            target_lang=target_lang,
            glossary=glossary,
            style=self._get_target_style(target_lang)
        )
        
        # 4. 应用翻译
        template = self._apply_translations(template, translations)
        
        # 5. 术语一致性检查
        inconsistencies = self._check_consistency(template, glossary)
        if inconsistencies:
            template = self._fix_inconsistencies(template, inconsistencies)
        
        # 6. 语言风格调整
        template = await self._adjust_style(template, target_lang)
        
        return template
    
    async def _translate_batch(
        self,
        texts: List[str],
        source_lang: str,
        target_lang: str,
        glossary: Dict[str, str],
        style: Dict
    ) -> List[str]:
        """批量翻译（保持教学语境）"""
        
        # 分批处理（每批20条）
        batch_size = 20
        batches = [texts[i:i+batch_size] for i in range(0, len(texts), batch_size)]
        
        all_translations = []
        
        for batch in batches:
            prompt = f"""
            你是教育内容本地化专家。请翻译以下教学内容：
            
            源语言：{source_lang}
            目标语言：{target_lang}
            
            **术语表**（必须严格使用）：
            {json.dumps(glossary, ensure_ascii=False, indent=2)}
            
            **风格要求**：
            - 正式度：{style['formality']} (0=非正式, 1=正式)
            - 互动性：{style['interactivity']} (问句比例)
            - 句长：{style['sentence_length']}
            - 人称：{style['use_you'] and '使用第二人称' or '使用第一人称复数'}
            
            **待翻译内容**：
            {chr(10).join(f'{i+1}. {t}' for i, t in enumerate(batch))}
            
            输出JSON数组格式，保持顺序：
            ["翻译1", "翻译2", ...]
            
            注意：
            1. 保持教学语境，不仅是字面翻译
            2. 术语必须使用术语表中的翻译
            3. 调整语言风格以符合目标文化
            """
            
            result = await self.llm.generate(prompt)
            all_translations.extend(result)
        
        return all_translations
    
    def _get_target_style(self, lang: str) -> Dict:
        """获取目标语言风格配置"""
        
        styles = {
            "en-US": {
                "formality": 0.5,
                "interactivity": 0.8,
                "sentence_length": "medium",
                "use_you": True,
                "question_ratio": 0.2,  # 20%是问句
                "exclamation_ratio": 0.1
            },
            "zh-CN": {
                "formality": 0.7,
                "interactivity": 0.4,
                "sentence_length": "short",
                "use_you": False,
                "question_ratio": 0.1,
                "exclamation_ratio": 0.05
            },
            "es-ES": {
                "formality": 0.6,
                "interactivity": 0.7,
                "sentence_length": "medium-long",
                "use_you": True,
                "question_ratio": 0.15,
                "exclamation_ratio": 0.15
            }
        }
        
        return styles.get(lang, styles["en-US"])


**术语表构建器**：

```python
class GlossaryBuilder:
    """术语表构建器"""
    
    async def build(
        self,
        template: WorkflowTemplate,
        source_lang: str,
        target_lang: str,
        domain: str,
        grade: str
    ) -> Dict[str, str]:
        """构建术语表"""
        
        # 1. 提取专业术语
        terms = self._extract_domain_terms(
            template,
            domain=domain
        )
        
        # 2. 查询术语库
        glossary = {}
        for term in terms:
            translation = await self._lookup_term(
                term,
                source_lang=source_lang,
                target_lang=target_lang,
                domain=domain
            )
            
            if translation:
                glossary[term] = translation
            else:
                # 使用LLM翻译未知术语
                translation = await self._translate_term(
                    term,
                    source_lang,
                    target_lang,
                    domain
                )
                glossary[term] = translation
                
                # 保存到术语库
                await self._save_term(term, translation, domain)
        
        return glossary
    
    def _extract_domain_terms(
        self,
        template: WorkflowTemplate,
        domain: str
    ) -> List[str]:
        """提取专业术语"""
        
        # 使用NER + 领域词典
        text = self._extract_all_text(template)
        
        # 化学领域词典
        chemistry_patterns = [
            r'\b[A-Z][a-z]?\d*\b',  # 化学式 (H2O, CO2)
            r'\b\w+反应\b',          # 反应类型
            r'\b\w+剂\b',            # 催化剂、氧化剂
            # ...
        ]
        
        terms = set()
        for pattern in chemistry_patterns:
            matches = re.findall(pattern, text)
            terms.update(matches)
        
        return list(terms)
```

---

**C. 文化适配引擎 (Cultural Adapter)**

```python
class CulturalAdapter:
    """文化适配引擎"""
    
    def __init__(self):
        self.llm = LLMClient()
        self.cultural_kb = CulturalKnowledgeBase()
    
    async def adapt_cultural_elements(
        self,
        template: WorkflowTemplate,
        source_region: str,
        target_region: str
    ) -> WorkflowTemplate:
        """适配文化元素"""
        
        # 1. 识别文化元素
        cultural_elements = self._identify_cultural_elements(
            template,
            source_region
        )
        
        # 2. 故事世界观重构
        if template.metadata.get('story_universe'):
            template.metadata['story_universe'] = await self._localize_story(
                template.metadata['story_universe'],
                target_region
            )
        
        # 3. 案例和比喻本地化
        for element in cultural_elements:
            if element.type == "case_study":
                localized = await self._generate_localized_case(
                    element,
                    target_region
                )
                template = self._replace_element(template, element, localized)
            
            elif element.type == "metaphor":
                localized = await self._localize_metaphor(
                    element,
                    target_region
                )
                template = self._replace_element(template, element, localized)
        
        # 4. 敏感内容审查
        sensitivity_issues = await self._check_sensitivity(
            template,
            target_region
        )
        
        if sensitivity_issues:
            template = await self._fix_sensitivity_issues(
                template,
                sensitivity_issues
            )
        
        return template
    
    def _identify_cultural_elements(
        self,
        template: WorkflowTemplate,
        source_region: str
    ) -> List[CulturalElement]:
        """识别文化元素"""
        
        elements = []
        
        # 加载区域特征库
        cultural_indicators = self.cultural_kb.get_indicators(source_region)
        
        # 提取所有文本
        text = self._extract_all_text(template)
        
        # 使用NER识别
        entities = self.ner_model.extract(text)
        
        for entity in entities:
            if entity.type in ["PERSON", "LOCATION", "EVENT", "FOOD", "HOLIDAY"]:
                # 检查是否为文化特定元素
                if self._is_cultural_specific(entity, source_region):
                    # 提取上下文
                    context = self._extract_context(text, entity, window=100)
                    
                    elements.append(CulturalElement(
                        type=self._classify_element_type(entity),
                        content=entity.text,
                        context=context,
                        location=entity.span
                    ))
        
        return elements
    
    async def _generate_localized_case(
        self,
        element: CulturalElement,
        target_region: str
    ) -> str:
        """生成本地化案例"""
        
        prompt = f"""
        你是教育内容本地化专家。请为{target_region}地区重新设计以下教学案例：
        
        原始案例（来自{element.source_region}）：
        {element.content}
        
        教学上下文：
        {element.context}
        
        学习目标：
        {element.learning_objective}
        
        要求：
        1. 保持教学目标不变
        2. 使用{target_region}学生熟悉的场景和例子
        3. 确保文化适当性
        4. 保持相似的趣味性和吸引力
        
        输出本地化后的案例（仅内容，不要解释）：
        """
        
        localized = await self.llm.generate(prompt)
        
        return localized


**文化敏感性规则库**：

```yaml
# cultural_sensitivity_rules.yaml

us:
  avoid:
    - type: "stereotype"
      examples:
        - "种族刻板印象"
        - "性别角色固化"
        - "宗教偏见"
    
    - type: "controversial"
      examples:
        - "政治立场"
        - "宗教争议"
        - "进化论vs创造论（需中立呈现）"
  
  emphasize:
    - "多元化角色（种族/性别/能力多样性）"
    - "环保意识"
    - "批判性思维"
    - "LGBTQ+包容性"

cn:
  avoid:
    - "政治敏感话题"
    - "领土争议"
    - "宗教宣传"
    - "暴力血腥内容"
  
  emphasize:
    - "社会主义核心价值观"
    - "爱国主义教育"
    - "传统文化传承"

sg:
  avoid:
    - "种族冲突"
    - "宗教对立"
    - "政治批评"
  
  emphasize:
    - "多元文化和谐"
    - "双语教育"
    - "新加坡精神"
```

---

**D. 教学法转换器 (Pedagogy Adapter)**

```python
class PedagogyAdapter:
    """教学法转换器"""
    
    def __init__(self):
        self.pedagogy_templates = PedagogyTemplateLibrary()
    
    async def adapt_pedagogy(
        self,
        template: WorkflowTemplate,
        target_region: str
    ) -> WorkflowTemplate:
        """适配教学法"""
        
        # 1. 获取目标区域教学法配置
        pedagogy_config = self.pedagogy_templates.get_config(target_region)
        
        # 2. 调整课型比例
        template = self._adjust_lesson_types(template, pedagogy_config)
        
        # 3. 替换教学流程模板
        template = self._replace_teaching_flow(template, pedagogy_config)
        
        # 4. 调整评估方式
        template = self._adapt_assessment(template, pedagogy_config)
        
        # 5. 调整互动密度
        template = self._adjust_interaction(template, pedagogy_config)
        
        return template
    
    def _replace_teaching_flow(
        self,
        template: WorkflowTemplate,
        config: Dict
    ) -> WorkflowTemplate:
        """替换教学流程模板"""
        
        # 获取教学流程模板
        flow_template = self.pedagogy_templates.get_flow(
            config['teaching_model']
        )
        
        # 替换每个课时的教学流程
        for lesson in template.metadata.get('lessons', []):
            # 保留学习目标和内容
            objectives = lesson['objectives']
            content = lesson['content']
            
            # 应用新的教学流程
            lesson['teaching_flow'] = flow_template.apply(
                objectives=objectives,
                content=content,
                lesson_type=lesson['type']
            )
        
        return template


**教学法模板库**：

```python
class PedagogyTemplateLibrary:
    """教学法模板库"""
    
    def get_config(self, region: str) -> Dict:
        """获取区域教学法配置"""
        
        configs = {
            "CN": {
                "teaching_model": "5_phase",
                "lesson_types": {
                    "NEW": 0.60,
                    "PRACTICE": 0.15,
                    "EXPERIMENT": 0.10,
                    "REVIEW": 0.10,
                    "ASSESSMENT": 0.05
                },
                "interaction_density": 0.4,
                "assessment_focus": "summative"
            },
            
            "US": {
                "teaching_model": "5E",
                "lesson_types": {
                    "INQUIRY_LAB": 0.30,
                    "DIRECT_INSTRUCTION": 0.25,
                    "COLLABORATIVE": 0.15,
                    "PROBLEM_SOLVING": 0.15,
                    "PRACTICE": 0.10,
                    "FORMATIVE_ASSESSMENT": 0.05
                },
                "interaction_density": 0.8,
                "assessment_focus": "formative"
            }
        }
        
        return configs.get(region, configs["US"])
    
    def get_flow(self, model_name: str) -> TeachingFlowTemplate:
        """获取教学流程模板"""
        
        if model_name == "5E":
            return FiveEModel()
        elif model_name == "5_phase":
            return FivePhaseModel()
        # ...


class FiveEModel(TeachingFlowTemplate):
    """5E探究模型"""
    
    def apply(self, objectives: List[str], content: str, lesson_type: str) -> Dict:
        """应用5E模型"""
        
        return {
            "phases": [
                {
                    "name": "Engage",
                    "duration": 10,
                    "activity": "Present phenomenon and driving question",
                    "teacher_role": "Facilitator",
                    "student_role": "Observer and questioner"
                },
                {
                    "name": "Explore",
                    "duration": 15,
                    "activity": "Hands-on investigation in groups",
                    "teacher_role": "Guide",
                    "student_role": "Investigator"
                },
                {
                    "name": "Explain",
                    "duration": 10,
                    "activity": "Student-led discussion",
                    "teacher_role": "Facilitator",
                    "student_role": "Presenter and discussant"
                },
                {
                    "name": "Elaborate",
                    "duration": 10,
                    "activity": "Apply to new situations",
                    "teacher_role": "Coach",
                    "student_role": "Problem solver"
                },
                {
                    "name": "Evaluate",
                    "duration": 5,
                    "activity": "Formative assessment + reflection",
                    "teacher_role": "Assessor",
                    "student_role": "Self-reflector"
                }
            ],
            "model_name": "5E Inquiry Model"
        }
```

---

**E. 合规性检查器 (Compliance Checker)**

```python
class ComplianceChecker:
    """合规性检查器"""
    
    def __init__(self):
        self.rules = ComplianceRuleLibrary()
        self.vision_model = VisionModel()  # 用于生成alt文本
    
    async def check_compliance(
        self,
        template: WorkflowTemplate,
        target_region: str
    ) -> ComplianceReport:
        """检查合规性"""
        
        report = ComplianceReport()
        
        # 1. 隐私合规
        if target_region == "US":
            coppa_issues = await self._check_coppa(template)
            report.add_category("coppa", coppa_issues)
            
            ferpa_issues = await self._check_ferpa(template)
            report.add_category("ferpa", ferpa_issues)
        
        elif target_region == "EU":
            gdpr_issues = await self._check_gdpr(template)
            report.add_category("gdpr", gdpr_issues)
        
        # 2. 无障碍标准
        accessibility_issues = await self._check_accessibility(
            template,
            standard=self._get_accessibility_standard(target_region)
        )
        report.add_category("accessibility", accessibility_issues)
        
        # 3. 版权合规
        copyright_issues = await self._check_copyright(template)
        report.add_category("copyright", copyright_issues)
        
        return report
    
    async def _check_accessibility(
        self,
        template: WorkflowTemplate,
        standard: str
    ) -> List[ComplianceIssue]:
        """检查无障碍标准（ADA/WCAG）"""
        
        issues = []
        
        # 检查图片alt文本
        images = self._extract_all_images(template)
        for img in images:
            if not img.get('alt_text'):
                issues.append(ComplianceIssue(
                    severity="high",
                    rule="WCAG_2.1_1.1.1",
                    category="accessibility",
                    description=f"Image {img['url']} missing alt text",
                    location=img['location'],
                    auto_fixable=True,
                    fix_method="generate_alt_text"
                ))
        
        # 检查视频字幕
        videos = self._extract_all_videos(template)
        for video in videos:
            if not video.get('captions'):
                issues.append(ComplianceIssue(
                    severity="high",
                    rule="WCAG_2.1_1.2.2",
                    description=f"Video missing captions",
                    auto_fixable=False  # 需要人工添加
                ))
        
        # 检查颜色对比度
        color_issues = self._check_color_contrast(template)
        issues.extend(color_issues)
        
        return issues
    
    async def auto_fix(
        self,
        template: WorkflowTemplate,
        issues: List[ComplianceIssue]
    ) -> WorkflowTemplate:
        """自动修复问题"""
        
        for issue in issues:
            if not issue.auto_fixable:
                continue
            
            if issue.fix_method == "generate_alt_text":
                # 使用视觉模型生成alt文本
                img_url = self._extract_url(issue.location)
                alt_text = await self.vision_model.describe(
                    img_url,
                    context="educational chemistry content",
                    detail="descriptive"
                )
                
                template = self._add_alt_text(template, img_url, alt_text)
        
        return template
```

---

#### 3.3.3 本地化流程编排

```python
class TemplateLocalizer:
    """模板本地化编排器"""
    
    async def localize_template(
        self,
        source_template_id: str,
        target_region: str,
        target_language: str,
        strategy: str = "adaptive"
    ) -> LocalizationJob:
        """执行本地化流程"""
        
        # 1. 创建本地化任务
        job = LocalizationJob(
            source_template_id=source_template_id,
            target_region=target_region,
            target_language=target_language,
            strategy=strategy,
            status="initializing"
        )
        
        try:
            # 2. 克隆模板
            job.update_status("cloning")
            new_template = await self._clone_template(source_template_id)
            
            # 3. 课程标准对齐
            job.update_status("aligning_standards")
            standard_report = await self.standard_mapper.map_standards(
                new_template,
                target_region
            )
            job.add_report("standard_mapping", standard_report)
            
            # 4. 语言本地化
            job.update_status("translating")
            new_template = await self.translator.translate_template(
                new_template,
                source_lang=new_template.metadata['language'],
                target_lang=target_language
            )
            
            # 5. 文化适配
            job.update_status("adapting_culture")
            new_template = await self.cultural_adapter.adapt_cultural_elements(
                new_template,
                source_region=new_template.metadata['region'],
                target_region=target_region
            )
            
            # 6. 教学法调整
            job.update_status("adapting_pedagogy")
            new_template = await self.pedagogy_adapter.adapt_pedagogy(
                new_template,
                target_region
            )
            
            # 7. 合规性检查
            job.update_status("checking_compliance")
            compliance_report = await self.compliance_checker.check_compliance(
                new_template,
                target_region
            )
            
            # 8. 自动修复合规问题
            if compliance_report.has_auto_fixable_issues():
                new_template = await self.compliance_checker.auto_fix(
                    new_template,
                    compliance_report.get_auto_fixable_issues()
                )
            
            # 9. 质量评估
            job.update_status("assessing_quality")
            quality_report = await self.quality_assessor.assess(
                new_template,
                source_template=await self._load_template(source_template_id)
            )
            job.add_report("quality", quality_report)
            
            # 10. 保存并标记为draft
            new_template.metadata['status'] = 'draft'
            await self._save_template(new_template)
            
            job.update_status("completed")
            job.new_template_id = new_template.template_id
            
        except Exception as e:
            job.update_status("failed")
            job.error = str(e)
            raise
        
        return job
```

---

#### 3.3.4 数据模型

```python
@dataclass
class LocalizationJob:
    """本地化任务"""
    
    job_id: str
    source_template_id: str
    target_region: str
    target_language: str
    strategy: str  # adaptive / translate-only / full-rewrite
    
    status: str  # initializing / cloning / aligning_standards / translating / 
                 # adapting_culture / adapting_pedagogy / checking_compliance /
                 # assessing_quality / completed / failed
    
    new_template_id: Optional[str] = None
    
    # 进度追踪
    progress: float = 0.0  # 0-1
    current_step: str = ""
    
    # 报告
    reports: Dict[str, Any] = field(default_factory=dict)
    
    # 审核
    review_items: List[ReviewItem] = field(default_factory=list)
    
    # 时间
    created_at: datetime = field(default_factory=datetime.now)
    updated_at: datetime = field(default_factory=datetime.now)
    completed_at: Optional[datetime] = None
    
    error: Optional[str] = None


@dataclass
class StandardMapping:
    """课程标准映射"""
    
    source: str
    target: str
    match_reason: str
    
    confidence_score: float
    requires_review: bool = False
    review_reason: Optional[str] = None
    
    knowledge_gaps: List[str] = field(default_factory=list)
    new_requirements: List[str] = field(default_factory=list)
    differences: Dict[str, str] = field(default_factory=dict)


@dataclass
class ComplianceIssue:
    """合规问题"""
    
    severity: str  # high / medium / low
    category: str  # coppa / ferpa / gdpr / accessibility / copyright
    rule: str      # 规则代码
    description: str
    
    location: Optional[str] = None  # 在模板中的位置
    
    auto_fixable: bool = False
    fix_method: Optional[str] = None
    suggested_fix: Optional[str] = None


@dataclass
class QualityReport:
    """质量报告"""
    
    # 各维度评分
    structure_score: float  # 结构完整性
    language_score: float   # 语言质量
    cultural_score: float   # 文化适配度
    alignment_score: float  # 课标对齐度
    compliance_score: float # 合规完整度
    
    overall_score: float
    status: str  # 可以发布 / 需要小幅优化 / 需要重大改进
    
    suggestions: List[str] = field(default_factory=list)
```

---

### 3.4 规则引擎设计

#### 3.3.1 规则定义格式

```yaml
# rules/math_rules.yaml
rules:
  - name: "数学课程必须包含逻辑验证"
    priority: 100
    condition:
      subject: "数学"
    actions:
      - type: "insert_node"
        node: "LogicVerificationNode"
        position: "after:content_generation"
      - type: "set_config"
        key: "enable_formula_check"
        value: true
  
  - name: "小学低年级使用游戏化"
    priority: 90
    condition:
      grade: [1, 2, 3]
    actions:
      - type: "set_template"
        template: "gamification"
      - type: "set_config"
        key: "interaction_density"
        value: "high"
  
  - name: "美国地区合规检查"
    priority: 80
    condition:
      region: "US"
    actions:
      - type: "insert_node"
        node: "COPPAComplianceNode"
      - type: "bind_resource"
        resource_type: "knowledge_base"
        resource_id: "us_common_core_standards"
```

#### 3.3.2 规则引擎实现

```python
class RuleEngine:
    def __init__(self):
        self.rules = []
        self.load_all_rules()
    
    def evaluate(self, context: WorkflowContext) -> RuleResult:
        """评估所有规则并返回执行动作"""
        
        matched_rules = []
        
        for rule in sorted(self.rules, key=lambda r: r.priority, reverse=True):
            if self.match_condition(rule.condition, context):
                matched_rules.append(rule)
        
        # 合并所有匹配规则的动作
        result = RuleResult()
        for rule in matched_rules:
            result.merge(rule.actions)
        
        return result
    
    def match_condition(self, condition: dict, context: WorkflowContext) -> bool:
        """检查条件是否满足"""
        
        for key, value in condition.items():
            context_value = context.get(key)
            
            if isinstance(value, list):
                # 范围匹配
                if context_value not in value:
                    return False
            else:
                # 精确匹配
                if context_value != value:
                    return False
        
        return True
```

### 3.4 提示词工程模块设计

#### 3.4.1 提示词组件化

```python
class PromptBuilder:
    """提示词构建器"""
    
    def __init__(self):
        self.components = {
            "context": [],
            "role": [],
            "instruction": [],
            "constraints": [],
            "examples": [],
            "output_format": []
        }
    
    def add_context(self, context: str):
        """添加上下文"""
        self.components["context"].append(context)
        return self
    
    def add_role(self, role: str):
        """设置角色"""
        self.components["role"].append(role)
        return self
    
    def add_instruction(self, instruction: str):
        """添加指令"""
        self.components["instruction"].append(instruction)
        return self
    
    def add_constraint(self, constraint: str):
        """添加约束"""
        self.components["constraints"].append(constraint)
        return self
    
    def add_example(self, example: str):
        """添加示例"""
        self.components["examples"].append(example)
        return self
    
    def build(self) -> str:
        """构建完整提示词"""
        
        sections = []
        
        if self.components["role"]:
            sections.append("[ROLE]\n" + "\n".join(self.components["role"]))
        
        if self.components["context"]:
            sections.append("[CONTEXT]\n" + "\n".join(self.components["context"]))
        
        if self.components["instruction"]:
            sections.append("[INSTRUCTION]\n" + "\n".join(self.components["instruction"]))
        
        if self.components["constraints"]:
            sections.append("[CONSTRAINTS]\n" + "\n".join(self.components["constraints"]))
        
        if self.components["examples"]:
            sections.append("[EXAMPLES]\n" + "\n".join(self.components["examples"]))
        
        return "\n\n".join(sections)
```

#### 3.4.2 模板库管理

```python
class PromptTemplateManager:
    """提示词模板管理器"""
    
    def __init__(self, template_dir: str):
        self.template_dir = template_dir
        self.templates = {}
        self.load_templates()
    
    def get_template(self, subject: str, node_type: str, region: str = "default"):
        """获取模板"""
        
        # 优先级: 区域定制 > 学科定制 > 默认
        template_key = f"{subject}_{node_type}_{region}"
        if template_key in self.templates:
            return self.templates[template_key]
        
        template_key = f"{subject}_{node_type}"
        if template_key in self.templates:
            return self.templates[template_key]
        
        template_key = f"default_{node_type}"
        return self.templates.get(template_key)
    
    def render_template(self, template: str, variables: dict) -> str:
        """渲染模板"""
        from jinja2 import Template
        
        jinja_template = Template(template)
        return jinja_template.render(**variables)
```

### 3.5 RAG检索模块设计

#### 3.5.1 向量检索流程

```python
class RAGRetriever:
    """RAG检索器"""
    
    def __init__(self, vector_db, embedding_model, reranker=None):
        self.vector_db = vector_db
        self.embedding_model = embedding_model
        self.reranker = reranker
    
    async def retrieve(self, query: str, knowledge_base: str, top_k: int = 5):
        """检索相关文档"""
        
        # 1. 查询向量化
        query_embedding = await self.embedding_model.encode(query)
        
        # 2. 向量检索
        vector_results = await self.vector_db.search(
            collection=knowledge_base,
            query_vector=query_embedding,
            top_k=top_k * 2  # 召回更多候选
        )
        
        # 3. 关键词检索（BM25）
        keyword_results = await self.keyword_search(query, knowledge_base, top_k)
        
        # 4. 混合检索结果
        merged_results = self.merge_results(vector_results, keyword_results)
        
        # 5. 重排序
        if self.reranker:
            ranked_results = await self.reranker.rank(query, merged_results)
        else:
            ranked_results = merged_results
        
        # 6. 返回Top-K
        return ranked_results[:top_k]
    
    def merge_results(self, vector_results, keyword_results):
        """合并检索结果"""
        
        # RRF (Reciprocal Rank Fusion)
        merged = {}
        
        for rank, doc in enumerate(vector_results, 1):
            doc_id = doc["id"]
            merged[doc_id] = merged.get(doc_id, 0) + 1.0 / (60 + rank)
        
        for rank, doc in enumerate(keyword_results, 1):
            doc_id = doc["id"]
            merged[doc_id] = merged.get(doc_id, 0) + 1.0 / (60 + rank)
        
        # 按分数排序
        sorted_docs = sorted(merged.items(), key=lambda x: x[1], reverse=True)
        
        return [{"id": doc_id, "score": score} for doc_id, score in sorted_docs]
```

#### 3.5.2 知识库结构

```python
class KnowledgeBase:
    """知识库管理"""
    
    def __init__(self, name: str):
        self.name = name
        self.metadata = {
            "subject": "",
            "region": "",
            "version": "",
            "last_updated": ""
        }
    
    async def add_documents(self, documents: List[Document]):
        """添加文档"""
        
        for doc in documents:
            # 文本分块
            chunks = self.chunk_document(doc)
            
            # 向量化
            embeddings = await self.embedding_model.encode_batch(chunks)
            
            # 存入向量库
            await self.vector_db.insert(
                collection=self.name,
                documents=chunks,
                embeddings=embeddings,
                metadata=[{
                    "source": doc.source,
                    "page": chunk.page,
                    "type": doc.type
                } for chunk in chunks]
            )
    
    def chunk_document(self, document: Document, chunk_size: int = 500):
        """文档分块"""
        
        # 根据文档类型选择分块策略
        if document.type == "textbook":
            # 按章节分块
            return self.chunk_by_section(document)
        elif document.type == "formula":
            # 保持公式完整性
            return self.chunk_preserve_formula(document)
        else:
            # 滑动窗口分块
            return self.chunk_sliding_window(document, chunk_size, overlap=50)
```

---

## 4. 技术选型

### 4.1 核心技术栈

| 组件 | 技术选型 | 版本 | 选择理由 |
|------|---------|------|---------|
| **后端框架** | FastAPI | 0.104+ | 高性能、类型提示、自动文档 |
| **工作流引擎** | LangGraph + 自研 | - | 灵活的图结构编排 |
| **数据库** | PostgreSQL | 15+ | 事务支持、JSON类型 |
| **缓存** | Redis | 7+ | 高性能、持久化支持 |
| **向量数据库** | Milvus | 2.3+ | 开源、高性能、易扩展 |
| **对象存储** | MinIO ⭐ 新增 | 2024+ | S3兼容、开源、支持图片/视频存储 |
| **消息队列** | Redis Streams | - | 轻量级、与Redis集成 |
| **LLM接口** | LiteLLM ⭐ 新增 | 1.0+ | 统一多provider接口（OpenAI/Claude/本地LLM） |
| **监控** | Prometheus + Grafana | - | 开源标准方案 |
| **错误追踪** | Sentry ⭐ 新增 | - | 实时错误监控、堆栈追踪 |
| **日志** | Loguru | - | 简洁易用 |
| **任务调度** | APScheduler | - | Python原生 |

---

### 4.1.5 V2.0新增技术栈 ⭐

#### **1. 多源内容检索框架**

| 技术 | 版本 | 用途 | 选择理由 |
|------|------|------|---------|
| **Haystack** | 1.25+ | 统一检索框架 | • 支持多种检索器（BM25/Dense/Hybrid）<br>• 预置RAG流水线<br>• 与Milvus无缝集成 |
| **Scrapy** | 2.11+ | 网页爬虫 | • 成熟稳定的爬虫框架<br>• 支持异步并发<br>• 丰富的中间件生态 |
| **BeautifulSoup4** | 4.12+ | HTML解析 | • 简单易用<br>• 容错能力强<br>• 适合小规模爬取 |
| **aiohttp** | 3.9+ | 异步HTTP客户端 | • 高性能异步请求<br>• 支持并发控制<br>• 与asyncio完美配合 |

**使用场景**：
```python
# Haystack：统一检索接口
from haystack import Pipeline
from haystack.nodes import BM25Retriever, DenseRetriever

pipeline = Pipeline()
pipeline.add_node(component=BM25Retriever(...), name="BM25", inputs=["Query"])
pipeline.add_node(component=DenseRetriever(...), name="Dense", inputs=["Query"])

# Scrapy：爬取WikiArt艺术作品
class WikiArtSpider(scrapy.Spider):
    name = "wikiart"
    start_urls = ["https://www.wikiart.org/en/artists"]
    
    def parse(self, response):
        for artist in response.css('.artist-name'):
            yield {"name": artist.css('::text').get()}
```

---

#### **2. 外部API集成**

| API服务 | 用途 | 认证方式 | 配额限制 |
|---------|------|---------|---------|
| **Wikipedia API** | 通用知识检索 | 无需认证 | 200 req/s |
| **WikiArt API** | 艺术作品/艺术家信息 | API Key | 1000 req/day（免费） |
| **YouTube Data API v3** | 教学视频检索 | OAuth 2.0 | 10,000 quota/day |
| **Goodreads API** | 书籍信息（畅销书） | API Key | 已弃用，改用爬虫 |
| **豆瓣API** | 中文书籍信息 | 无需认证 | 非官方API，限流严格 |
| **Stack Overflow API** | 技术问答（编程培训） | OAuth 2.0 | 300 req/day（无认证） |
| **GitHub API** | 代码示例（编程培训） | Personal Token | 5,000 req/hour（认证） |
| **MDN Web Docs** | Web技术文档 | 无需认证 | 无限制（静态页面） |

**集成示例**：
```python
# Wikipedia API集成
import aiohttp

async def search_wikipedia(query: str, language: str = "zh") -> List[Dict]:
    url = f"https://{language}.wikipedia.org/w/api.php"
    params = {
        "action": "query",
        "format": "json",
        "list": "search",
        "srsearch": query,
        "srlimit": 10
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.get(url, params=params) as response:
            data = await response.json()
            return data["query"]["search"]

# YouTube Data API集成
from googleapiclient.discovery import build

youtube = build('youtube', 'v3', developerKey=YOUTUBE_API_KEY)

def search_youtube_videos(query: str, max_results: int = 5):
    request = youtube.search().list(
        part="snippet",
        q=query,
        type="video",
        maxResults=max_results
    )
    return request.execute()
```

---

#### **3. LLM统一管理：LiteLLM**

**为什么选择LiteLLM**：
- ✅ 统一接口支持100+ LLM provider（OpenAI、Claude、DeepSeek、本地LLM等）
- ✅ 自动重试、负载均衡、fallback机制
- ✅ 内置Token计算、成本追踪
- ✅ 支持流式输出、函数调用
- ✅ 兼容OpenAI SDK接口，无需大规模代码修改

**核心优势**：
```python
# V1.0方式：针对不同provider需要不同代码
if provider == "openai":
    response = openai.ChatCompletion.create(...)
elif provider == "claude":
    response = anthropic.messages.create(...)
elif provider == "deepseek":
    response = deepseek.chat.create(...)

# V2.0方式：统一接口
from litellm import completion

response = completion(
    model="gpt-4-turbo",  # 或 "claude-3-opus" 或 "deepseek-chat"
    messages=[{"role": "user", "content": "Hello"}]
)
# 自动路由到对应provider，返回统一格式
```

**高级功能**：
```python
# 1. Fallback机制（主模型失败自动切换备用）
response = completion(
    model="gpt-4-turbo",
    messages=messages,
    fallbacks=["claude-3-opus", "deepseek-chat"]
)

# 2. 负载均衡（多个API Key轮询）
litellm.api_key = [KEY1, KEY2, KEY3]

# 3. 成本追踪
from litellm import completion_cost
cost = completion_cost(completion_response=response)

# 4. Token计算
from litellm import token_counter
tokens = token_counter(model="gpt-4", messages=messages)
```

**配置示例**：
```yaml
# litellm_config.yaml
model_list:
  - model_name: gpt-4-turbo
    litellm_params:
      model: gpt-4-turbo
      api_key: ${OPENAI_API_KEY}
      rpm: 500  # 每分钟请求限制
      
  - model_name: claude-3-opus
    litellm_params:
      model: claude-3-opus-20240229
      api_key: ${ANTHROPIC_API_KEY}
      rpm: 1000
      
  - model_name: deepseek-chat
    litellm_params:
      model: deepseek-chat
      api_base: https://api.deepseek.com
      api_key: ${DEEPSEEK_API_KEY}
      rpm: 200
      
  - model_name: local-llama
    litellm_params:
      model: ollama/llama2
      api_base: http://localhost:11434
```

---

#### **4. 向量数据库：Milvus深度应用**

**V2.0新增应用场景**：

| 场景 | Collection | 向量维度 | 索引类型 | 用途 |
|------|-----------|---------|---------|------|
| **文本语义检索** | `content_embeddings` | 1536 (OpenAI Ada) | IVF_FLAT | 多源内容语义搜索 |
| **图像检索** | `artwork_embeddings` | 512 (CLIP) | HNSW | 美术作品相似度搜索 |
| **学习目标去重** | `objective_embeddings` | 768 (BGE-M3) | IVF_SQ8 | 检测重复学习目标 |
| **叙述策略匹配** | `strategy_embeddings` | 384 (MiniLM) | FLAT | 匹配最佳叙述策略 |

**配置示例**：
```python
from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType

# 连接Milvus
connections.connect(host="localhost", port="19530")

# 创建Collection Schema
fields = [
    FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
    FieldSchema(name="source", dtype=DataType.VARCHAR, max_length=50),
    FieldSchema(name="content", dtype=DataType.VARCHAR, max_length=5000),
    FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536)
]
schema = CollectionSchema(fields, description="Content embeddings")

# 创建Collection
collection = Collection(name="content_embeddings", schema=schema)

# 创建索引（IVF_FLAT）
index_params = {
    "metric_type": "IP",  # 内积（Inner Product）
    "index_type": "IVF_FLAT",
    "params": {"nlist": 1024}
}
collection.create_index(field_name="embedding", index_params=index_params)

# 混合检索（向量 + 标量过滤）
search_params = {"metric_type": "IP", "params": {"nprobe": 10}}
results = collection.search(
    data=[query_embedding],
    anns_field="embedding",
    param=search_params,
    limit=10,
    expr='source == "wikipedia"',  # 标量过滤
    output_fields=["content", "source"]
)
```

**性能优化**：
```python
# 1. 批量插入（提升10倍速度）
entities = [
    [source1, source2, ...],
    [content1, content2, ...],
    [embedding1, embedding2, ...]
]
collection.insert(entities)

# 2. 异步搜索（并发多个查询）
futures = []
for query in queries:
    future = collection.search_async(data=[query], ...)
    futures.append(future)

results = [f.result() for f in futures]

# 3. 内存优化（仅加载部分字段）
collection.load(_resource_groups=["rg1"], _load_fields=["embedding"])
```

---

#### **5. 对象存储：MinIO**

**为什么需要MinIO**：
- V2.0新增图像内容（美术作品、教学图片）
- 视频内容（YouTube教程、案例演示）
- 生成的多媒体文件（带图教案、PPT等）

**MinIO优势**：
- ✅ S3兼容API（可随时迁移到AWS S3）
- ✅ 开源免费、部署简单
- ✅ 高性能（支持多磁盘RAID）
- ✅ 支持版本控制、生命周期管理

**存储结构设计**：
```
minio-bucket/
├── artworks/              # 艺术作品图片
│   ├── monet/
│   │   ├── impression_sunrise.jpg
│   │   └── water_lilies.jpg
│   └── van_gogh/
│       └── starry_night.jpg
│
├── videos/                # 教学视频
│   ├── excel/
│   │   └── pivot_table_tutorial.mp4
│   └── programming/
│       └── python_basics.mp4
│
├── generated/             # 生成的内容
│   ├── lessons/
│   │   └── math_grade3_20231201.pdf
│   └── presentations/
│       └── art_history_monet.pptx
│
└── cache/                 # 缓存的外部资源
    └── youtube_thumbnails/
```

**集成示例**：
```python
from minio import Minio
from minio.error import S3Error

# 初始化MinIO客户端
minio_client = Minio(
    "localhost:9000",
    access_key="minioadmin",
    secret_key="minioadmin",
    secure=False
)

# 上传艺术作品图片
def upload_artwork(artist: str, artwork_name: str, file_path: str):
    bucket_name = "artworks"
    object_name = f"{artist}/{artwork_name}.jpg"
    
    # 确保bucket存在
    if not minio_client.bucket_exists(bucket_name):
        minio_client.make_bucket(bucket_name)
    
    # 上传文件
    minio_client.fput_object(
        bucket_name,
        object_name,
        file_path,
        content_type="image/jpeg"
    )
    
    # 返回URL
    return f"http://localhost:9000/{bucket_name}/{object_name}"

# 生成预签名URL（临时访问链接，7天有效）
def get_presigned_url(bucket: str, object_name: str) -> str:
    from datetime import timedelta
    
    url = minio_client.presigned_get_object(
        bucket,
        object_name,
        expires=timedelta(days=7)
    )
    return url

# 批量下载YouTube视频缩略图
async def cache_youtube_thumbnails(video_ids: List[str]):
    for video_id in video_ids:
        thumbnail_url = f"https://img.youtube.com/vi/{video_id}/hqdefault.jpg"
        
        async with aiohttp.ClientSession() as session:
            async with session.get(thumbnail_url) as response:
                if response.status == 200:
                    data = await response.read()
                    
                    # 上传到MinIO
                    minio_client.put_object(
                        "cache",
                        f"youtube_thumbnails/{video_id}.jpg",
                        io.BytesIO(data),
                        length=len(data),
                        content_type="image/jpeg"
                    )
```

**生命周期管理**：
```python
from minio.lifecycleconfig import LifecycleConfig, Rule, Expiration

# 自动删除30天前的缓存文件
config = LifecycleConfig(
    [
        Rule(
            rule_id="expire-cache",
            status="Enabled",
            expiration=Expiration(days=30),
            rule_filter={"prefix": "cache/"}
        )
    ]
)

minio_client.set_bucket_lifecycle("artworks", config)
```

---

#### **6. 错误监控：Sentry**

**为什么需要Sentry**：
- V2.0复杂度大增（4个领域、3个AI引擎、10+数据源）
- 需要实时监控生产环境错误
- 快速定位问题根因（堆栈追踪、上下文）

**集成方式**：
```python
import sentry_sdk
from sentry_sdk.integrations.fastapi import FastApiIntegration
from sentry_sdk.integrations.redis import RedisIntegration

sentry_sdk.init(
    dsn="https://xxx@xxx.ingest.sentry.io/xxx",
    environment="production",
    integrations=[
        FastApiIntegration(),
        RedisIntegration(),
    ],
    traces_sample_rate=0.1,  # 10%的请求追踪性能
    profiles_sample_rate=0.1,
)

# 自定义错误上下文
with sentry_sdk.configure_scope() as scope:
    scope.set_tag("domain", "art_history")
    scope.set_context("adapter", {
        "adapter_id": "art-history-adapter",
        "version": "1.0.0"
    })
    scope.set_user({
        "id": user_id,
        "email": user_email
    })

# 捕获自定义错误
try:
    adapter.parse_input(raw_input)
except Exception as e:
    sentry_sdk.capture_exception(e)
    raise
```

**告警规则配置**：
```yaml
# Sentry告警规则
alerts:
  - name: "适配器路由失败率过高"
    condition: "error_count > 10 in 5 minutes"
    filter: "error.type = AdapterRoutingError"
    notification: "slack://devops-channel"
    
  - name: "LLM调用超时"
    condition: "avg(duration) > 30s"
    filter: "transaction.name = completion"
    notification: "email://oncall@example.com"
    
  - name: "向量检索失败"
    condition: "error_count > 5 in 1 minute"
    filter: "error.type = MilvusConnectionError"
    notification: "pagerduty://高优先级"
```

---

#### **7. 爬虫反反爬策略**

**挑战**：WikiArt、豆瓣等网站有反爬机制

**解决方案**：

| 技术 | 用途 | 实现 |
|------|------|------|
| **User-Agent轮换** | 模拟不同浏览器 | `fake-useragent`库 |
| **请求间隔** | 避免频率过高 | `asyncio.sleep(random.uniform(1, 3))` |
| **代理池** | 分散IP | `ProxyBroker`或购买代理服务 |
| **Session复用** | 保持Cookie | `aiohttp.ClientSession` |
| **JS渲染** | 爬取动态页面 | `Playwright`（Scrapy-Playwright插件） |

**Scrapy配置**：
```python
# settings.py
ROBOTSTXT_OBEY = False  # 根据需要决定是否遵守robots.txt

# 下载延迟（秒）
DOWNLOAD_DELAY = 2
RANDOMIZE_DOWNLOAD_DELAY = True

# 并发限制
CONCURRENT_REQUESTS = 16
CONCURRENT_REQUESTS_PER_DOMAIN = 4

# User-Agent轮换
DOWNLOADER_MIDDLEWARES = {
    'scrapy.downloadermiddlewares.useragent.UserAgentMiddleware': None,
    'scrapy_fake_useragent.middleware.RandomUserAgentMiddleware': 400,
}

# 自动限速
AUTOTHROTTLE_ENABLED = True
AUTOTHROTTLE_START_DELAY = 1
AUTOTHROTTLE_MAX_DELAY = 10
AUTOTHROTTLE_TARGET_CONCURRENCY = 2.0

# 重试
RETRY_TIMES = 3
RETRY_HTTP_CODES = [500, 502, 503, 504, 408, 429]
```

**Playwright动态爬取**：
```python
from scrapy_playwright.page import PageMethod

class DynamicSpider(scrapy.Spider):
    name = "dynamic"
    
    def start_requests(self):
        yield scrapy.Request(
            url="https://example.com",
            meta={
                "playwright": True,
                "playwright_page_methods": [
                    PageMethod("wait_for_selector", ".content"),
                    PageMethod("screenshot", path="page.png"),
                ],
            }
        )
```

---

#### **8. 技术选型对比表**

**向量数据库对比**（为什么选Milvus）：

| 特性 | Milvus | Pinecone | Weaviate | Qdrant |
|------|--------|----------|----------|--------|
| **开源** | ✅ | ❌ | ✅ | ✅ |
| **性能** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **可扩展性** | 分布式集群 | 托管服务 | 单机/集群 | 单机/集群 |
| **成本** | 免费（自托管） | 按量付费 | 免费/付费 | 免费/付费 |
| **社区** | 活跃（CNCF） | 商业主导 | 中等 | 中等 |
| **中文支持** | ✅ | ✅ | ✅ | ✅ |
| **多索引类型** | IVF/HNSW/FLAT | HNSW | HNSW | HNSW |
| **混合检索** | ✅ | ❌ | ✅ | ✅ |

**LLM管理对比**（为什么选LiteLLM）：

| 特性 | LiteLLM | LangChain | 原生SDK |
|------|---------|-----------|---------|
| **统一接口** | ✅ 100+ models | ⚠️ 需封装 | ❌ 各不相同 |
| **Fallback** | ✅ 内置 | ⚠️ 需自己实现 | ❌ 无 |
| **成本追踪** | ✅ 自动 | ❌ 手动 | ❌ 手动 |
| **负载均衡** | ✅ 支持 | ❌ 无 | ❌ 无 |
| **学习曲线** | 低 | 高 | 中 |
| **性能开销** | 极小 | 中等 | 无 |
| **适用场景** | 多模型切换 | 复杂agent | 单一模型 |

**对象存储对比**（为什么选MinIO）：

| 特性 | MinIO | AWS S3 | 阿里云OSS | 本地文件系统 |
|------|-------|--------|----------|-------------|
| **成本** | 免费 | 按量付费 | 按量付费 | 免费 |
| **部署** | 本地/云端 | 仅云端 | 仅云端 | 本地 |
| **S3兼容** | ✅ | ✅ | 部分 | ❌ |
| **可迁移性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐ |
| **性能** | 高（本地） | 高（网络） | 高（网络） | 最高 |
| **可靠性** | 取决于配置 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 低 |
| **适用场景** | 开发/小规模 | 生产环境 | 生产环境 | 原型开发 |

---

### 4.2 LLM模型选型

| 应用场景 | 主要模型 | 备用模型 | 理由 |
|---------|---------|---------|------|
| **创意生成** | GPT-4-Turbo | Claude-3-Opus | 创造力强 |
| **逻辑推理** | DeepSeek-V2 | GPT-4 | 成本低、数学能力强 |
| **事实核查** | RAG + GPT-3.5 | - | 成本低、准确性高 |
| **合规检查** | 本地BERT | - | 数据安全、低成本 |
| **向量化** | OpenAI Ada-002 | BGE-Large | 质量高 |

### 4.3 开发工具链

```yaml
开发环境:
  Python: "3.11+"
  包管理: "Poetry"
  代码格式: "Black + isort"
  类型检查: "Mypy"
  测试: "Pytest"
  
CI/CD:
  版本控制: "Git"
  CI工具: "GitHub Actions"
  容器: "Docker"
  编排: "Kubernetes"
  
文档:
  API文档: "自动生成(FastAPI)"
  架构文档: "Markdown + Mermaid"
  代码文档: "Docstring + Sphinx"
```

---

## 5. 数据流设计

### 5.1 完整数据流图

```
用户请求
    │
    ▼
┌─────────────────┐
│  API Gateway    │  ← 认证、限流、参数验证
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ 参数解析器       │  → 提取: subject, grade, region, topic
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  规则引擎       │  → 匹配规则 → 确定节点列表
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  DAG构建器      │  → 构建执行图
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  执行调度器     │  → 按拓扑顺序调度节点
└─────────────────┘
    │
    ├──────────────────────────────────────┐
    │                                      │
    ▼                                      ▼
┌─────────────────┐                 ┌─────────────────┐
│ Node 1: Persona │                 │ Node 2: RAG     │
│                 │                 │                 │
│ ┌─────────────┐ │                 │ ┌─────────────┐ │
│ │ Prompt      │ │                 │ │ Query Vec   │ │
│ │ Builder     │ │                 │ │             │ │
│ └──────┬──────┘ │                 │ └──────┬──────┘ │
│        ▼        │                 │        ▼        │
│ ┌─────────────┐ │                 │ ┌─────────────┐ │
│ │ LLM Router  │ │                 │ │ Vector DB   │ │
│ └──────┬──────┘ │                 │ └──────┬──────┘ │
│        ▼        │                 │        ▼        │
│ ┌─────────────┐ │                 │ ┌─────────────┐ │
│ │ GPT-4       │ │                 │ │ Top-K Docs  │ │
│ └──────┬──────┘ │                 │ └──────┬──────┘ │
│        ▼        │                 │        ▼        │
│   PersonaCard   │                 │   KnowledgeList │
└────────┬────────┘                 └────────┬────────┘
         │                                   │
         └───────────────┬───────────────────┘
                         ▼
                 ┌─────────────────┐
                 │ Node 3: Content │
                 │   Generation    │
                 │                 │
                 │ Input:          │
                 │  - PersonaCard  │
                 │  - Knowledge    │
                 │                 │
                 │ Output:         │
                 │  - LessonScript │
                 └────────┬────────┘
                          ▼
                 ┌─────────────────┐
                 │ Node 4: Verify  │
                 └────────┬────────┘
                          ▼
                 ┌─────────────────┐
                 │ State Manager   │  → 保存结果到Redis/DB
                 └────────┬────────┘
                          ▼
                      返回用户
```

### 5.2 数据模型设计

详见独立的 `data_models.md` 文档。

---

## 6. 部署架构

### 6.1 容器化部署

```yaml
# docker-compose.yml
version: '3.8'

services:
  # API网关
  api-gateway:
    image: metaworkflow/api-gateway:latest
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/metaworkflow
      - REDIS_URL=redis://redis:6379/0
    depends_on:
      - postgres
      - redis
  
  # 工作流引擎
  workflow-engine:
    image: metaworkflow/engine:latest
    environment:
      - REDIS_URL=redis://redis:6379/0
      - MILVUS_HOST=milvus
    depends_on:
      - redis
      - milvus
  
  # PostgreSQL数据库
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=metaworkflow
    volumes:
      - postgres-data:/var/lib/postgresql/data
  
  # Redis缓存
  redis:
    image: redis:7-alpine
    volumes:
      - redis-data:/data
  
  # Milvus向量数据库
  milvus:
    image: milvusdb/milvus:v2.3.0
    environment:
      - ETCD_ENDPOINTS=etcd:2379
      - MINIO_ADDRESS=minio:9000
    depends_on:
      - etcd
      - minio
  
  # 监控
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
  
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana

volumes:
  postgres-data:
  redis-data:
  prometheus-data:
  grafana-data:
```

### 6.2 Kubernetes部署架构

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: workflow-engine
spec:
  replicas: 3  # 水平扩展
  selector:
    matchLabels:
      app: workflow-engine
  template:
    metadata:
      labels:
        app: workflow-engine
    spec:
      containers:
      - name: engine
        image: metaworkflow/engine:v1.0
        resources:
          requests:
            cpu: "1"
            memory: "2Gi"
          limits:
            cpu: "2"
            memory: "4Gi"
        env:
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: url
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

### 6.3 高可用架构

```
                    ┌─────────────┐
                    │   CDN       │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Load Balance│ (Nginx/ALB)
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼───┐          ┌───▼───┐          ┌───▼───┐
    │API GW │          │API GW │          │API GW │
    │  Pod1 │          │  Pod2 │          │  Pod3 │
    └───┬───┘          └───┬───┘          └───┬───┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼────┐         ┌───▼────┐        ┌───▼────┐
    │Engine  │         │Engine  │        │Engine  │
    │ Pod1   │         │ Pod2   │        │ Pod3   │
    └───┬────┘         └───┬────┘        └───┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼────┐         ┌───▼────┐        ┌───▼────┐
    │Postgres│         │ Redis  │        │Milvus  │
    │ Master │         │Cluster │        │Cluster │
    │   +    │         │        │        │        │
    │ Replica│         │        │        │        │
    └────────┘         └────────┘        └────────┘
```

---

## 7. 安全架构

### 7.1 认证与授权

```python
# OAuth2 + JWT
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

async def get_current_user(token: str = Depends(oauth2_scheme)):
    """验证JWT令牌"""
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        user_id: str = payload.get("sub")
        if user_id is None:
            raise HTTPException(status_code=401)
        return user_id
    except JWTError:
        raise HTTPException(status_code=401)

# RBAC权限控制
class PermissionChecker:
    def __init__(self, required_permission: str):
        self.required_permission = required_permission
    
    async def __call__(self, user=Depends(get_current_user)):
        user_permissions = await get_user_permissions(user)
        if self.required_permission not in user_permissions:
            raise HTTPException(status_code=403)
        return user
```

### 7.2 数据加密

```yaml
传输层加密:
  - TLS 1.3
  - 证书管理: Let's Encrypt / ACM
  
存储层加密:
  - 数据库: AES-256加密
  - 敏感字段: 单独加密存储
  - API密钥: 环境变量 + Secrets Manager
  
应用层:
  - 密码: bcrypt哈希
  - 个人信息: 自动脱敏
```

### 7.3 安全监控

```python
# 审计日志
class AuditLogger:
    def log_action(self, user_id, action, resource, result):
        log_entry = {
            "timestamp": datetime.utcnow(),
            "user_id": user_id,
            "action": action,
            "resource": resource,
            "result": result,
            "ip_address": get_client_ip()
        }
        self.db.audit_logs.insert(log_entry)

# 异常检测
class SecurityMonitor:
    def check_anomaly(self, user_id, action):
        # 检测异常行为
        recent_actions = self.get_recent_actions(user_id, minutes=5)
        
        if len(recent_actions) > 100:  # 频率异常
            self.alert("High frequency access", user_id)
        
        if action == "delete" and len(recent_actions) > 10:  # 批量删除
            self.alert("Suspicious deletion", user_id)
```

---

## 附录

### A. 系统配置示例

见 `config/` 目录下的配置文件。

### B. API接口文档

见独立的 `api_design.md` 文档。

### C. 性能基准测试

见独立的 `performance_benchmark.md` 文档。

---

**文档结束**
