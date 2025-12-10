# 文档总览 - 通用领域内容生成·元工作流平台 V2.0

**版本**: V2.0  
**更新日期**: 2025-12-09  
**状态**: 设计阶段 → 开发准备中

本目录包含元工作流平台的**完整V2.0设计文档**，涵盖多领域支持、AI能力增强、模板系统等核心特性。

---

## 📚 文档结构

```
metaworkflow/
├── README.md                          # 项目总览（V2.0）
├── docs/                              # 设计文档
│   ├── README.md                      # 📖 本文档（文档导航）
│   ├── requirements.md                # 需求分析文档 ⭐
│   ├── architecture.md                # 架构设计文档 ⭐
│   ├── data_models.md                 # 数据模型设计
│   ├── api_design.md                  # API接口设计（V2.0 78+端点）
│   ├── roadmap.md                     # 开发路线图 ⭐ V2.0（6人团队，1.5个月）
│   ├── DEPLOYMENT_GUIDE.md            # 部署指南（Docker + K8s）⭐ 新增
│   ├── UPDATE_SUMMARY.md              # V2.0升级总结 ⭐ 新增
│   ├── USE_CASES.md                   # 使用场景文档
│   ├── TEMPLATE_LOCALIZATION_GUIDE.md # 模板本地化指南
│   └── COMPONENT_EXTENSION_GUIDE.md   # 组件扩展指南
├── config/                            # 配置示例
│   ├── workflows/                     # 工作流配置
│   │   └── math_lesson_basic_cn.yaml # 数学课程工作流
│   ├── templates/                     # 工作流模板 ⭐ 新增
│   │   ├── math_lesson_template_v1.yaml
│   │   └── README.md
│   ├── rules/                         # 规则配置
│   │   └── global_rules.yaml         # 全局规则库
│   └── prompts/                       # Prompt模板
│       └── math_content_generation_cn.yaml
├── examples/                          # 示例配置 ⭐ 新增
│   ├── curriculum_chemistry_9_cn.yaml # 化学课标示例
│   ├── curriculum_chemistry_9_us.yaml # 美国化学课标
│   └── README.md
├── src/                               # 源代码（开发中）
└── tests/                             # 测试代码（待开发）
```

---

## 🎯 核心文档导读

### 1. [需求分析文档](requirements.md) 📋

**阅读时长**: 30分钟  
**适合人群**: 所有项目参与者  
**版本**: V2.0 ⭐ 已更新

**主要内容**:
- ✅ 项目背景与核心问题（从K12到多领域平台）
- ✅ V2.0核心特性（AI目标生成 + 内容发现 + 叙述引擎）
- ✅ 业务需求分析（多领域、多受众、多场景）
- ✅ 功能需求详解（模板系统、智能推荐、本地化适配）
- ✅ 非功能需求（性能、安全、可扩展性）
- ✅ 用户画像与核心用例

**核心亮点**:
- **多领域支持**: K12教育 + 艺术史 + 畅销书解读 + 职业培训
- **AI能力升级**: 从模糊需求到结构化目标，从预设知识库到实时检索
- **故事化叙述**: 7种叙述策略，让内容更吸引人
- **模板系统**: 预编排工作流模板，降低使用门槛
- **课程时长参数化**: 支持3-15分钟微课和45-90分钟常规课 ⭐ 最新更新

**关键词**: 多领域、AI驱动、实时检索、故事化、模板化

---

### 2. [架构设计文档](architecture.md) 🏗️

**阅读时长**: 50分钟  
**适合人群**: 技术团队、架构师  
**版本**: V2.0 ⭐ 已更新

**主要内容**:
- ✅ V2.0总体架构设计（7层架构）
- ✅ **4大核心模块**（Domain Adapter、AI Objectives、Content Discovery、Narrative Engine）⭐
- ✅ 工作流引擎与规则引擎
- ✅ **向量检索系统**（Milvus集成）⭐ 新增
- ✅ **工作流模板管理器**（智能推荐、自动构建）⭐ 新增
- ✅ **多源内容检索引擎**（Wikipedia、WikiArt、YouTube等）⭐ 新增
- ✅ 技术选型与理由（PostgreSQL + Milvus + Redis + MinIO）
- ✅ 数据流设计
- ✅ 部署架构（Docker + Kubernetes）

**核心亮点**:
- **领域适配器架构**: 插件化设计，支持热插拔领域模块
- **AI目标生成器**: 基于LLM的需求分析与目标生成
- **内容发现引擎**: 向量检索 + 多源聚合 + 质量评分
- **叙述生成引擎**: 7种叙述策略（时间轴、英雄之旅、问题导向等）
- **模板推荐系统**: 基于协同过滤的智能推荐
- **自动化工作流构建**: 根据需求自动组装节点

**关键代码示例**:
```python
# V2.0 领域适配器接口
class DomainAdapter:
    async def parse_request(self, raw_input: str) -> UserRequest:
        """解析用户输入"""
        
    async def validate_context(self, context: WorkflowContext) -> bool:
        """验证上下文"""
        
    async def customize_workflow(self, workflow: Workflow) -> Workflow:
        """定制工作流"""
```

**关键词**: 领域适配器、AI模块、向量检索、模板系统、微服务

---

### 3. [数据模型设计](data_models.md) 📊

**阅读时长**: 30分钟  
**适合人群**: 后端开发、数据库工程师  
**版本**: V2.0 ⭐ 已更新

**主要内容**:
- ✅ 核心数据模型（WorkflowContext、LearningObjective、ContentFragment）
- ✅ **领域适配器模型**（DomainConfig、AdapterMetadata）⭐ 新增
- ✅ **AI模块数据模型**（AnalysisResult、DiscoveryResult、NarrativeResult）⭐ 新增
- ✅ **工作流模板模型**（WorkflowTemplate、TemplateParameter）⭐ 新增
- ✅ 执行状态模型（WorkflowInstance、NodeExecution）
- ✅ **向量数据模型**（EmbeddingVector、SearchResult）⭐ 新增
- ✅ PostgreSQL表结构设计（23+表）
- ✅ Redis数据结构设计
- ✅ **Milvus向量集合设计**（content_embeddings）⭐ 新增

**核心亮点**:
- **领域配置模型**: 支持多领域动态注册
- **AI模块接口**: 统一的AI能力调用接口
- **模板参数化**: 灵活的模板配置与实例化
- **向量索引**: HNSW索引，支持百万级向量检索
- **课程时长参数**: duration字段支持3-15分钟微课 ⭐ 最新更新

**数据库Schema亮点**:
```sql
-- V2.0 新增：领域适配器表
CREATE TABLE domain_adapters (
    id UUID PRIMARY KEY,
    domain_name VARCHAR(50) NOT NULL,
    adapter_class VARCHAR(200) NOT NULL,
    config JSONB,
    enabled BOOLEAN DEFAULT true
);

-- V2.0 新增：向量索引表
CREATE TABLE content_embeddings (
    id BIGSERIAL PRIMARY KEY,
    content_id UUID NOT NULL,
    embedding_vector VECTOR(1536),
    metadata JSONB
);
```

**关键词**: 多领域、模板化、向量化、参数化、可扩展

---

### 4. [API接口设计](api_design.md) 🔌

**阅读时长**: 40分钟  
**适合人群**: 前后端开发、集成工程师  
**版本**: V2.0 ⭐ 已更新

**主要内容**:
- ✅ RESTful API规范（V2.0版本）
- ✅ **领域适配器管理接口**（9个端点）⭐ 新增
- ✅ **AI目标生成接口**（6个端点）⭐ 新增
- ✅ **AI内容发现接口**（7个端点）⭐ 新增
- ✅ **AI叙述生成接口**（8个端点）⭐ 新增
- ✅ **多源内容检索接口**（8个端点）⭐ 新增
- ✅ **工作流模板管理接口**（10个端点）⭐ 新增
- ✅ 工作流实例管理接口
- ✅ 监控统计接口
- ✅ Webhook回调机制

**核心亮点**:
- **78+个API端点**: 覆盖完整业务流程
- **统一的响应格式**: 标准化错误处理
- **详细的请求/响应示例**: 每个接口都有完整示例
- **Python SDK设计**: 面向对象的客户端库
- **批量操作支持**: 批量生成学习目标、批量检索内容

**API示例**:
```bash
# V2.0 新增：AI目标生成
POST /v2/ai/objectives/analyze
{
  "raw_input": "我想讲解莫奈的印象派绘画，面向初中生，时长10分钟",
  "preferences": {
    "duration": 10,
    "max_objectives": 5
  }
}

# 响应
{
  "success": true,
  "data": {
    "domain": "art_history",
    "audience": {
      "age_group": "中学生",
      "knowledge_level": "beginner"
    },
    "learning_objectives": [
      {
        "title": "认识莫奈及其代表作品",
        "bloom_level": "记忆",
        "estimated_time": 3
      }
    ]
  }
}
```

**关键词**: V2 API、AI能力、批量操作、模板管理、向量检索

---

### 5. [开发路线图](roadmap.md) 🗓️

**阅读时长**: 40分钟  
**适合人群**: 项目经理、技术负责人、投资方  
**版本**: V2.0 ⭐ 完全重写

**主要内容**:
- ✅ **V2.0开发计划**（1.5个月，6位资深工程师）⭐ 全新
- ✅ **4个Sprint规划**（基础设施 → AI模块 → 领域适配 →测试上线）⭐
- ✅ 团队组织（Team A核心平台组 + Team B AI能力组）
- ✅ 详细任务分解（Day-by-Day任务清单）
- ✅ 里程碑与验收标准（4个关键里程碑）
- ✅ 风险管理（技术、进度、质量风险及缓解措施）
- ✅ 交付清单（代码、文档、部署、测试）

**核心亮点**:
- **紧凑但可执行**: 6-7周完成V2.0全部开发
- **明确的团队分工**: 6人分为核心平台组和AI能力组
- **并行开发策略**: 基础设施与AI模块并行推进
- **详细的日级任务**: 每位工程师都有清晰的每日任务
- **完整的质量保证**: 单元测试>85%、集成测试>90%

**时间线**:
```
Sprint 1 (Week 1-2)  → 基础设施 + 核心框架完成 (M1: 12-29)
Sprint 2 (Week 3-4)  → 3大AI模块完成 (M2: 01-12)
Sprint 3 (Week 5)    → 4个领域适配器 + 集成 (M3: 01-19)
Sprint 4 (Week 6-7)  → 测试 + 优化 + 上线 (M4: 01-31) 🚀
```

**团队配置**:
- **工程师A1** (Tech Lead): 架构设计、核心框架、代码审查
- **工程师A2**: 数据层、API开发、数据库设计
- **工程师A3**: 工作流引擎、任务调度、性能优化
- **工程师B1** (AI Lead): AI目标生成、LLM集成、Prompt工程
- **工程师B2**: 内容发现引擎、向量检索、Milvus集成
- **工程师B3**: 叙述生成引擎、领域适配器开发

**关键词**: 6人团队、1.5个月、4个Sprint、并行开发、质量保证

---

### 6. [部署指南](DEPLOYMENT_GUIDE.md) �

**阅读时长**: 35分钟  
**适合人群**: DevOps、运维工程师  
**版本**: V2.0 ⭐ 新增

**主要内容**:
- ✅ 系统架构概述
- ✅ **环境要求**（Python 3.11+、PostgreSQL 15+、Milvus 2.3+等）
- ✅ **本地开发环境部署**（Docker Compose一键启动）⭐
- ✅ **生产环境部署**（Kubernetes + Helm Charts）⭐
- ✅ **监控系统部署**（Prometheus + Grafana + Sentry）⭐
- ✅ 配置管理（环境变量、配置文件）
- ✅ 数据备份与恢复
- ✅ 故障排查指南

**核心亮点**:
- **一键部署**: `docker-compose up -d` 启动完整环境
- **生产级配置**: Kubernetes YAML + Helm Charts
- **完整的监控栈**: 性能监控、错误追踪、日志聚合
- **自动扩缩容**: HPA配置，支持动态伸缩
- **安全加固**: HTTPS、密钥管理、网络隔离

**部署示例**:
```bash
# 本地开发环境
git clone https://github.com/your-org/metaworkflow.git
cd metaworkflow
cp .env.example .env
docker-compose up -d

# 生产环境（Kubernetes）
helm install metaworkflow ./helm/metaworkflow \
  --namespace production \
  --create-namespace \
  --values values-production.yaml
```

**关键词**: Docker、Kubernetes、监控、安全、高可用

---

### 7. [V2.0升级总结](UPDATE_SUMMARY.md) 📝

**阅读时长**: 15分钟  
**适合人群**: 所有项目参与者  
**版本**: V2.0 ⭐ 新增

**主要内容**:
- ✅ V2.0核心变更概述
- ✅ 新增功能清单
- ✅ 架构升级要点
- ✅ API变更说明
- ✅ 数据模型变更
- ✅ 迁移指南（从V1.0到V2.0）

**核心亮点**:
- **从垂直到水平**: K12教育工具 → 通用领域平台
- **AI能力突破**: 3大AI模块全新上线
- **模板系统**: 预编排工作流模板，降低使用门槛
- **向量检索**: Milvus支持语义搜索
- **多领域支持**: 4个初始领域（可扩展）

**关键词**: 升级总结、变更说明、迁移指南

---
## 🔧 配置示例导读

### 1. [数学课程工作流配置](../config/workflows/math_lesson_basic_cn.yaml)

**内容**:
- 完整的YAML工作流定义
- 10个节点的详细配置
- 条件节点示例
- 超时与重试策略
- 错误处理与降级

**亮点**:
```yaml
# 条件节点: 仅当年级<=3时执行故事化设计
- id: "story_design"
  condition:
    operator: "OR"
    expressions:
      - field: "context.grade"
        operator: "<="
        value: 3
```

**适用场景**: 小学数学课程生成（V1.0示例）

---

### 2. [小学数学课件模板](../config/templates/math_lesson_template_v1.yaml) ⭐ V2.0新增

**内容**:
- 预编排工作流模板定义
- 模板参数化配置
- 适用范围与推荐条件
- 版本管理与本地化
- 使用示例与最佳实践

**亮点**:
```yaml
metadata:
  template_name: "小学低年级数学基础课件"
  category: "k12_education"
  tags: ["数学", "小学", "低年级", "基础课件"]
  
parameters:
  - name: "topics"
    type: "array"
    required: true
    description: "知识点列表"
  
  - name: "duration"
    type: "integer"
    default: 45
    constraints:
      min: 3
      max: 15
    description: "课程时长（分钟），微课建议3-15分钟"
```

**适用场景**: 快速生成小学低年级数学课件

---

### 3. [全局规则配置](../config/rules/global_rules.yaml)

**内容**:
- 内容质量检查规则
- 安全审查规则
- 时长控制规则
- 难度适配规则

**亮点**:
```yaml
# 时长控制规则（V2.0更新）
- rule_id: "duration_control"
  condition:
    field: "context.duration"
    operator: "between"
    value: [3, 90]  # 支持3-90分钟范围
  action:
    type: "adjust_content_detail"
    params:
      micro_course_max: 15  # 微课最大15分钟
      regular_course_min: 45  # 常规课最小45分钟
```

**适用场景**: 全局质量保障

---

### 4. [Prompt模板示例](../config/prompts/math_content_generation_cn.yaml)

**内容**:
- 数学内容生成Prompt模板
- Few-shot示例
- 变量替换规则
- 多版本管理

**适用场景**: LLM内容生成

---

## 🌍 示例配置导读（V2.0新增）

### 1. [中国初三化学课标](../examples/curriculum_chemistry_9_cn.yaml)

**内容**:
- 完整的中国初三化学课程标准
- 知识点体系（7个主题、42个知识点）
- 能力要求（布鲁姆分类法）
- 教学建议

**适用场景**: K12化学课程生成的课标参考

---

### 2. [美国九年级化学课标](../examples/curriculum_chemistry_9_us.yaml)

**内容**:
- 完整的美国Next Generation Science Standards (NGSS)
- 与中国课标的对比
- 本地化差异（单位、术语、教学法）

**适用场景**: 国际化课程生成的本地化参考

---

## 📖 扩展指南文档

### 1. [模板本地化指南](TEMPLATE_LOCALIZATION_GUIDE.md)

**内容**:
- 工作流模板的本地化适配方法
- 国家/地区差异处理
- 文化敏感内容审查
- 本地化测试方法

**适合人群**: 国际化团队、本地化工程师

---

### 2. [组件扩展指南](COMPONENT_EXTENSION_GUIDE.md)

**内容**:
- 自定义节点开发
- 自定义规则引擎插件
- 自定义内容源接入
- 扩展点清单与最佳实践

**适合人群**: 后端开发、插件开发者

---

### 3. [使用场景文档](USE_CASES.md)

**内容**:
- K12教育场景（数学、语文、物理等）
- 艺术史教育场景（名画讲解、艺术家传记）
- 畅销书解读场景（精读笔记、读书会讲稿）
- 职业培训场景（Excel技能、编程入门）
- 完整的端到端示例

**适合人群**: 产品经理、业务团队、新用户

---

## 🚀 快速上手路径

### 面向业务人员

1. **了解系统能做什么** → 阅读[项目总览](../README.md)（15分钟）
2. **查看实际案例** → 阅读[使用场景文档](USE_CASES.md)（20分钟）
3. **理解功能需求** → 阅读[需求分析文档](requirements.md)（30分钟）
4. **查看配置示例** → 查看[config/](../config/)目录下的YAML文件

**总计**: ~1小时

---

### 面向技术人员

1. **理解系统架构** → 阅读[架构设计文档](architecture.md)（50分钟）
2. **理解数据模型** → 阅读[数据模型设计](data_models.md)（30分钟）
3. **理解API接口** → 阅读[API接口设计](api_design.md)（40分钟）
4. **查看部署方案** → 阅读[部署指南](DEPLOYMENT_GUIDE.md)（35分钟）
5. **了解开发计划** → 阅读[开发路线图](roadmap.md)（40分钟）

**总计**: ~3小时

---

### 面向DevOps工程师

1. **了解技术栈** → 阅读[架构设计文档](architecture.md)中的"技术选型"章节（20分钟）
2. **学习部署方案** → 阅读[部署指南](DEPLOYMENT_GUIDE.md)（35分钟）
3. **配置监控系统** → 查看部署指南中的"监控系统部署"章节（15分钟）
4. **准备开发环境** → 跑通本地Docker Compose环境（30分钟）

**总计**: ~1.5小时

---

### 面向项目经理/投资方

1. **了解产品定位** → 阅读[项目总览](../README.md)（15分钟）
2. **理解V2.0升级** → 阅读[V2.0升级总结](UPDATE_SUMMARY.md)（15分钟）
3. **查看开发计划** → 阅读[开发路线图](roadmap.md)（40分钟）
4. **评估风险与资源** → 重点阅读roadmap中的"风险管理"和"交付清单"章节（10分钟）

**总计**: ~1.5小时

---

## 📊 V2.0关键指标

| 指标 | V1.0 | V2.0 | 提升 |
|------|------|------|------|
| **支持领域** | 1个（K12教育） | 4个（K12 + 艺术史 + 畅销书 + 职业培训） | 4x ⭐ |
| **API端点** | ~25个 | 78+个 | 3x ⭐ |
| **数据表** | 15个 | 23+个 | 1.5x |
| **AI能力模块** | 0个 | 3个（目标生成 + 内容发现 + 叙述生成） | 新增 ⭐ |
| **向量检索** | 不支持 | 支持（Milvus） | 新增 ⭐ |
| **工作流模板** | 不支持 | 支持（智能推荐 + 自动构建） | 新增 ⭐ |
| **课程时长** | 固定15分钟 | 3-90分钟可配置 | 灵活性⬆ ⭐ |
| **开发周期** | 24周（4人） | 7周（6人） | 效率提升3.4x ⭐ |

---

## 🔄 文档更新日志

### 2025-12-09 (V2.0)
- ✅ 更新所有核心文档到V2.0版本
- ✅ 新增DEPLOYMENT_GUIDE.md（部署指南）
- ✅ 新增UPDATE_SUMMARY.md（升级总结）
- ✅ 重写roadmap.md（6人团队、1.5个月计划）
- ✅ 更新README.md（多领域平台定位）
- ✅ 更新docs/README.md（本文档）
- ✅ 课程时长参数化（支持3-15分钟微课）⭐ 最新

### 2025-12-08 (V2.0-alpha)
- ✅ architecture.md添加4大AI模块设计
- ✅ api_design.md添加V2.0 API设计（78+端点）
- ✅ data_models.md添加向量数据模型
- ✅ 新增模板系统设计

### 2025-11-30 (V1.0)
- ✅ 初始版本，K12教育垂直工具

---

## 💡 贡献指南

文档维护原则：
1. **保持同步**: 代码变更必须同步更新文档
2. **示例完整**: 每个功能都要有完整的代码/配置示例
3. **版本标注**: 新增功能标注版本号（如"⭐ V2.0新增"）
4. **分层阅读**: 为不同角色提供不同的阅读路径
5. **持续优化**: 根据反馈不断改进文档结构

---

## 📮 联系方式

- **技术负责人**: [待定]
- **产品负责人**: [待定]
- **文档维护**: GitHub Issues
- **紧急联系**: [待定]

---

**© 2025 MetaWorkflow Platform - 通用领域内容生成·元工作流平台 V2.0**  
**文档版本**: V2.0  
**最后更新**: 2025-12-09

---

### 2. [全局规则配置](../config/rules/global_rules.yaml)

**内容**:
- 20+条业务规则
- 学科规则(数学/语文/美术)
- 年级规则(低年级/高年级)
- 地区规则(中国/美国)
- 教学法规则(探究式/翻转课堂)

**亮点**:
```yaml
# 规则: 数学课程必须包含逻辑验证
- rule_id: "math_logic_verification_required"
  conditions:
    subject: "数学"
  actions:
    - type: "insert_node"
      parameters:
        node_type: "LogicVerificationNode"
```

---

### 3. [Prompt模板示例](../config/prompts/math_content_generation_cn.yaml)

**内容**:
- 完整的Prompt模板结构
- 系统角色定义
- 上下文模板
- 指令模板
- Few-shot示例
- 思维链引导

**亮点**:
```yaml
chain_of_thought_guidance: |
  生成时请按以下思路:
  1. 分析学生特点 → 决定是否使用故事化
  2. 确定教学策略 → 决定采用什么教学法
  3. 设计教学流程 → 决定如何拆解
  4. 设计互动活动 → 决定如何提问
  5. 整合教材内容 → 决定如何引用
```

---

## 📖 推荐阅读路径

### 路径1: 业务人员/产品经理
1. `README.md` - 了解项目概览
2. `docs/requirements.md` - 理解业务需求
3. `docs/roadmap.md` - 掌握实施计划
4. `config/workflows/` - 查看配置示例

**预计时间**: 1.5小时

---

### 路径2: 技术负责人/架构师
1. `README.md` - 快速了解
2. `docs/architecture.md` - 深入架构设计
3. `docs/data_models.md` - 理解数据结构
4. `docs/api_design.md` - 掌握接口规范
5. `docs/roadmap.md` - 制定实施计划

**预计时间**: 2.5小时

---

### 路径3: 后端开发工程师
1. `docs/architecture.md` - 理解系统架构
2. `docs/data_models.md` - 熟悉数据模型
3. `docs/api_design.md` - 掌握API规范
4. `config/workflows/` - 理解工作流配置
5. `config/rules/` - 理解规则引擎

**预计时间**: 2小时

---

### 路径4: AI算法工程师
1. `docs/requirements.md` (重点: 3.2 故事化与模块化编排器)
2. `docs/architecture.md` (重点: 3.3 提示词工程模块, 3.4 RAG检索模块)
3. `config/prompts/` - 研究Prompt模板设计
4. `docs/roadmap.md` (重点: Phase 1 Week 1-2 技术验证)

**预计时间**: 2小时

---

## 🎓 关键概念速查

### 元工作流 (Meta-Workflow)
> 用于生成特定工作流的工作流模板。根据输入参数动态组装节点和配置。

### 三级节点体系
| 层级 | 职责 | 示例 |
|------|------|------|
| L1: 战略层 | 决定内容基调 | 用户画像分析 |
| L2: 设计层 | 生成教学内容 | 脚本生成 |
| L3: 生产层 | 质量保障 | 合规校验 |

### RAG (Retrieval-Augmented Generation)
> 检索增强生成,通过知识库检索来抑制LLM幻觉。

### DAG (Directed Acyclic Graph)
> 有向无环图,用于表示工作流的节点依赖关系。

---

## ❓ 常见问题 FAQ

**Q1: 为什么需要"元工作流"而不是直接用LLM?**  
A: 直接用LLM生成长文本容易产生幻觉、结构松散。元工作流将生成过程拆解为可控、可校验的标准化步骤,确保质量和合规性。

**Q2: 新增一个学科需要多久?**  
A: Phase 2完成后,新增学科只需配置规则和Prompt模板,预计2天内完成。

**Q3: 如何保证内容的准确性?**  
A: 通过RAG检索权威教材、逻辑验证节点、事实核查节点三层保障。

**Q4: 支持哪些国家和地区?**  
A: V1.0支持中国和美国,后续将扩展到新加坡、英国等。

**Q5: 成本如何控制?**  
A: 通过模型路由(创意用GPT-4,逻辑用DeepSeek)、Prompt优化、缓存策略来降低成本。

---

## 📞 获取帮助

- **文档问题**: 提交Issue到GitHub
- **技术讨论**: 加入技术社区(待建立)
- **商务合作**: contact@metaworkflow.edu

---

**最后更新**: 2025-12-09  
**文档版本**: V1.0  
**维护者**: Meta-Workflow文档团队
