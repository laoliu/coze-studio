# MetaWorkflow V2.0 适配器市场与组件生态设计

**文档版本**: V1.0  
**日期**: 2025-12-09  
**作者**: Meta-Workflow架构团队  
**文档类型**: 详细设计文档

---

## 📑 文档说明

本文档是 MetaWorkflow V2.0 组件扩展体系的**完整设计规范**，包含：

1. **适配器架构设计** - 插件化架构的技术实现
2. **注册机制设计** - 适配器的发现、注册、管理
3. **市场设计** - 适配器的发布、安装、评价
4. **开发指南** - 第三方开发者文档
5. **质量保证** - 审核、测试、安全标准

本文档分为多个部分，逐步展开。

---

## 目录

### 第一部分：架构基础
1. [适配器架构设计](#第一部分适配器架构设计)
   - 1.1 架构模型
   - 1.2 接口规范
   - 1.3 生命周期管理

### 第二部分：注册机制
2. [适配器注册机制](#第二部分适配器注册机制)
   - 2.1 注册表设计
   - 2.2 自动发现
   - 2.3 依赖管理

### 第三部分：市场生态
3. [适配器市场设计](#第三部分适配器市场设计)
   - 3.1 市场架构
   - 3.2 发布流程
   - 3.3 审核机制

### 第四部分：开发指南
4. [扩展开发指南](#第四部分扩展开发指南)
   - 4.1 快速开始
   - 4.2 最佳实践
   - 4.3 调试测试

### 第五部分：质量保证
5. [质量保证体系](#第五部分质量保证体系)
   - 5.1 质量标准
   - 5.2 自动化测试
   - 5.3 安全审计

---

# 第一部分：适配器架构设计

## 1.1 架构模型

### 1.1.1 三层架构

```
┌──────────────────────────────────────────────────────────────┐
│                     应用层 (Application Layer)                │
│                                                               │
│  用户通过API或UI与系统交互                                      │
│  • REST API (FastAPI)                                        │
│  • Web UI (Future)                                           │
│  • CLI Tools                                                 │
└──────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────┐
│                   核心平台层 (Platform Core)                  │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  工作流引擎 (Workflow Engine)                        │    │
│  │  • DAG构建和调度                                     │    │
│  │  • 节点执行管理                                      │    │
│  │  • 并行处理                                          │    │
│  └─────────────────────────────────────────────────────┘    │
│                              ↕                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  适配器管理器 (Adapter Manager)                      │    │
│  │  • 适配器注册表 (AdapterRegistry)                    │    │
│  │  • 适配器生命周期管理                                │    │
│  │  • 依赖解析和注入                                    │    │
│  └─────────────────────────────────────────────────────┘    │
│                              ↕                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  通用服务 (Common Services)                          │    │
│  │  • AI服务 (LLM Client)                               │    │
│  │  • 向量检索 (Milvus)                                 │    │
│  │  • 缓存服务 (Redis)                                  │    │
│  │  • 存储服务 (MinIO)                                  │    │
│  └─────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────┐
│                   适配器层 (Adapter Layer)                    │
│                                                               │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌──────────┐ │
│  │   K12     │  │  美术史    │  │  畅销书    │  │   ...    │ │
│  │  Adapter  │  │  Adapter  │  │  Adapter  │  │  Adapter │ │
│  └───────────┘  └───────────┘  └───────────┘  └──────────┘ │
│                                                               │
│  每个适配器实现 DomainAdapter 接口：                           │
│  • parse_request()         - 解析用户输入                     │
│  • validate_context()      - 验证上下文                       │
│  • generate_learning_objectives() - 生成学习目标              │
│  • discover_content()      - 发现内容                         │
│  • customize_workflow()    - 定制工作流                       │
│  • format_output()         - 格式化输出                       │
└──────────────────────────────────────────────────────────────┘
```

### 1.1.2 适配器组成结构

每个适配器是一个独立的Python包，包含以下组件：

```
k12_education_adapter/           # 适配器包名
├── __init__.py                  # 包初始化
├── adapter.py                   # 主适配器类（实现DomainAdapter接口）
├── manifest.yaml                # 适配器清单（元数据）
├── config/                      # 配置文件
│   ├── content_sources.yaml    # 内容源配置
│   ├── prompt_templates.yaml   # Prompt模板
│   └── quality_rules.yaml      # 质量评估规则
├── models/                      # 领域特定数据模型
│   ├── __init__.py
│   ├── subject.py              # 科目模型
│   └── grade.py                # 年级模型
├── services/                    # 领域特定服务
│   ├── __init__.py
│   ├── curriculum_parser.py    # 课程标准解析
│   └── textbook_analyzer.py    # 教材分析
├── resources/                   # 静态资源
│   ├── knowledge_graph.json    # 知识图谱（可选）
│   ├── templates/              # 输出模板
│   └── assets/                 # 图片、图标等
├── tests/                       # 单元测试
│   ├── test_adapter.py
│   └── test_services.py
├── requirements.txt             # Python依赖
├── README.md                    # 说明文档
└── LICENSE                      # 许可证
```

### 1.1.3 适配器清单（manifest.yaml）

```yaml
# k12_education_adapter/manifest.yaml

# ==================== 基本信息 ====================
name: k12_education
display_name: K12教育
version: 1.2.0
description: |
  K12教育领域适配器，支持化学、物理、数学、生物等学科的
  教学内容生成，适用于初中和高中教育场景。

author:
  name: MetaWorkflow Team
  email: team@metaworkflow.ai
  organization: MetaWorkflow Inc.

license: MIT
homepage: https://github.com/metaworkflow/k12-adapter
repository: https://github.com/metaworkflow/k12-adapter
documentation: https://docs.metaworkflow.ai/adapters/k12

# ==================== 兼容性 ====================
compatibility:
  platform_version: ">=2.0.0,<3.0.0"  # 兼容的平台版本
  python_version: ">=3.11"             # 所需Python版本

# ==================== 依赖 ====================
dependencies:
  # Python包依赖
  python:
    - jieba>=0.42.1         # 中文分词
    - pinyin>=0.4.0         # 拼音转换
  
  # 其他适配器依赖（可选）
  adapters: []
  
  # 平台组件依赖
  components:
    - llm_client              # 需要LLM客户端
    - vector_store            # 需要向量数据库

# ==================== 功能声明 ====================
capabilities:
  # 支持的活动类型
  activity_types:
    - concept               # 概念理解
    - experiment            # 实验探究
    - problem              # 问题解决
  
  # 支持的输出格式
  output_formats:
    - ppt                  # PowerPoint
    - pdf                  # PDF文档
    - html                 # 网页
    - markdown             # Markdown
  
  # 支持的学科
  subjects:
    - 化学
    - 物理
    - 数学
    - 生物
  
  # 支持的年级范围
  grade_range:
    min: 7                 # 初一
    max: 12                # 高三

# ==================== 配置参数 ====================
configuration:
  # 用户可配置的参数（使用JSON Schema定义）
  schema:
    type: object
    properties:
      duration_range:
        type: object
        properties:
          min:
            type: integer
            default: 3
            minimum: 3
            description: 最小时长（分钟）
          max:
            type: integer
            default: 90
            maximum: 90
            description: 最大时长（分钟）
      
      curriculum_standard:
        type: string
        enum: [national, provincial, custom]
        default: national
        description: 课程标准类型
      
      difficulty_auto_adjust:
        type: boolean
        default: true
        description: 是否根据年级自动调整难度
    
    required: [duration_range]

# ==================== 资源声明 ====================
resources:
  # 知识图谱
  knowledge_graph:
    enabled: true
    path: resources/knowledge_graph.json
    format: json
    size: 2.5MB
  
  # Prompt模板
  prompt_templates:
    path: config/prompt_templates.yaml
    count: 15
  
  # 内容源配置
  content_sources:
    - name: 中国教育资源网
      type: web_scraping
      enabled: true
    - name: 人教版电子教材
      type: api
      enabled: false

# ==================== 性能指标 ====================
performance:
  # 预估的资源消耗
  estimated_memory: 256MB      # 内存占用
  estimated_cpu: low           # CPU占用（low/medium/high）
  
  # 性能基准
  benchmarks:
    parse_request: <50ms       # 解析请求耗时
    generate_objectives: <2s   # 生成目标耗时
    discover_content: <5s      # 内容发现耗时

# ==================== 质量保证 ====================
quality:
  # 测试覆盖率
  test_coverage: 85%
  
  # 质量评分（0-100）
  quality_score: 92
  
  # 审核状态
  review_status: approved      # pending/approved/rejected
  
  # 安全扫描
  security_scan:
    passed: true
    vulnerabilities: 0
    last_scan: 2025-12-08

# ==================== 发布信息 ====================
release:
  changelog: |
    ## v1.2.0 (2025-12-09)
    - 新增数学科目支持
    - 优化化学实验活动生成
    - 修复年级解析bug
    
    ## v1.1.0 (2025-11-15)
    - 添加生物科目支持
    - 改进内容发现算法
    
    ## v1.0.0 (2025-10-01)
    - 初始版本发布
  
  # 发布日期
  release_date: 2025-12-09
  
  # 下载统计
  downloads:
    total: 1523
    last_month: 342

# ==================== 标签和分类 ====================
tags:
  - education
  - k12
  - stem
  - chinese
  
categories:
  - Education
  - STEM

# ==================== 截图和媒体 ====================
media:
  icon: resources/icon.png
  screenshots:
    - resources/screenshots/demo1.png
    - resources/screenshots/demo2.png
  video_url: https://youtube.com/watch?v=xxx
```

---

## 1.2 接口规范

### 1.2.1 DomainAdapter 核心接口

```python
# src/adapters/base.py

from abc import ABC, abstractmethod
from typing import Dict, List, Any, Tuple, Optional
from datetime import datetime
from src.models import (
    UserRequest, WorkflowContext, 
    ActivityType, LearningObjective,
    DomainType
)

class AdapterMetadata:
    """适配器元数据"""
    
    def __init__(self, manifest: Dict[str, Any]):
        """从manifest.yaml加载元数据"""
        self.name = manifest['name']
        self.display_name = manifest['display_name']
        self.version = manifest['version']
        self.description = manifest['description']
        self.author = manifest.get('author', {})
        self.license = manifest.get('license', 'Unknown')
        self.compatibility = manifest.get('compatibility', {})
        self.dependencies = manifest.get('dependencies', {})
        self.capabilities = manifest.get('capabilities', {})
        self.configuration = manifest.get('configuration', {})
        self.resources = manifest.get('resources', {})
        self.performance = manifest.get('performance', {})
        self.quality = manifest.get('quality', {})
        self.release = manifest.get('release', {})
        self.tags = manifest.get('tags', [])
        self.categories = manifest.get('categories', [])


class DomainAdapter(ABC):
    """
    领域适配器抽象基类
    
    所有领域适配器必须继承此类并实现所有抽象方法。
    适配器负责将通用工作流引擎能力适配到特定领域的业务需求。
    
    生命周期：
    1. 初始化 (__init__)
    2. 验证兼容性 (validate_compatibility)
    3. 加载配置 (load_configuration)
    4. 注册到平台 (register)
    5. 处理请求 (parse_request -> validate_context -> ...)
    6. 卸载 (unload)
    """
    
    def __init__(self, config: Optional[Dict[str, Any]] = None):
        """
        初始化适配器
        
        Args:
            config: 用户提供的配置（覆盖默认配置）
        """
        self._config = config or {}
        self._metadata: Optional[AdapterMetadata] = None
        self._initialized = False
        self._loaded_at: Optional[datetime] = None
    
    # ==================== 元数据（必须实现） ====================
    
    @property
    @abstractmethod
    def metadata(self) -> AdapterMetadata:
        """
        获取适配器元数据
        
        通常从manifest.yaml文件加载
        
        Returns:
            AdapterMetadata: 适配器元数据对象
        """
        pass
    
    # ==================== 生命周期管理 ====================
    
    async def initialize(self) -> None:
        """
        初始化适配器（框架调用）
        
        执行时机：适配器首次加载时
        典型操作：
        - 加载知识图谱
        - 初始化内容源连接
        - 预加载Prompt模板
        - 验证依赖项
        """
        if self._initialized:
            return
        
        # 子类可以覆盖此方法添加初始化逻辑
        await self._on_initialize()
        
        self._initialized = True
        self._loaded_at = datetime.utcnow()
    
    async def _on_initialize(self) -> None:
        """子类覆盖此方法实现自定义初始化"""
        pass
    
    async def unload(self) -> None:
        """
        卸载适配器（框架调用）
        
        执行时机：适配器被移除或系统关闭时
        典型操作：
        - 关闭数据库连接
        - 释放缓存
        - 清理临时文件
        """
        await self._on_unload()
        self._initialized = False
    
    async def _on_unload(self) -> None:
        """子类覆盖此方法实现自定义清理"""
        pass
    
    def validate_compatibility(
        self, 
        platform_version: str
    ) -> Tuple[bool, Optional[str]]:
        """
        验证与平台的兼容性
        
        Args:
            platform_version: 当前平台版本（如 "2.0.0"）
            
        Returns:
            Tuple[bool, Optional[str]]: (是否兼容, 错误信息)
        """
        from packaging import version
        
        required = self.metadata.compatibility.get('platform_version', '>=2.0.0')
        
        # 解析版本约束（简化实现，实际可使用packaging库）
        if '>=' in required:
            min_version = required.replace('>=', '').split(',')[0]
            if version.parse(platform_version) < version.parse(min_version):
                return False, f"需要平台版本 {required}，当前版本 {platform_version}"
        
        return True, None
    
    # ==================== 核心业务方法（必须实现） ====================
    
    @abstractmethod
    async def parse_request(
        self, 
        request: UserRequest
    ) -> WorkflowContext:
        """
        解析用户请求，生成工作流上下文
        
        这是适配器的入口方法，负责：
        1. 从原始输入中提取领域相关信息
        2. 填充 WorkflowContext 的领域特定字段
        3. 设置初始参数和元数据
        
        Args:
            request: 用户原始请求
                - raw_input: 用户输入文本
                - domain: 领域类型
                - preferences: 用户偏好设置
            
        Returns:
            WorkflowContext: 包含解析结果的上下文
            
        Raises:
            ValueError: 输入格式不正确或无法解析
            
        示例（K12适配器）：
            ```python
            async def parse_request(self, request: UserRequest) -> WorkflowContext:
                # 提取科目
                subject = self._extract_subject(request.raw_input)
                
                # 提取年级
                grade = self._extract_grade(request.raw_input)
                
                # 提取时长
                duration = self._extract_duration(request.raw_input) \\
                          or request.preferences.get('duration', 45)
                
                return WorkflowContext(
                    workflow_id=uuid4(),
                    domain=DomainType.K12_EDUCATION,
                    raw_input=request.raw_input,
                    subject=subject,
                    grade_level=grade,
                    duration=duration,
                    metadata={
                        'curriculum': 'national',
                        'difficulty': self._estimate_difficulty(grade)
                    }
                )
            ```
        """
        pass
    
    @abstractmethod
    async def validate_context(
        self, 
        context: WorkflowContext
    ) -> Tuple[bool, Optional[str]]:
        """
        验证工作流上下文的完整性和合法性
        
        在工作流执行前调用，确保所有必要信息已填充且符合要求。
        
        Args:
            context: 待验证的工作流上下文
            
        Returns:
            Tuple[bool, Optional[str]]: 
                - 第一个元素：是否有效
                - 第二个元素：错误信息（有效时为None）
                
        示例（K12适配器）：
            ```python
            async def validate_context(self, context: WorkflowContext) -> Tuple[bool, Optional[str]]:
                # 检查必填字段
                if not context.subject:
                    return False, "未识别到科目，请在输入中明确指定（如：化学、物理）"
                
                if not context.grade_level:
                    return False, "未识别到年级，请指定（如：初三、高一、9年级）"
                
                # 检查参数范围
                if context.duration < 3 or context.duration > 90:
                    return False, "时长必须在3-90分钟之间"
                
                # 检查学科支持
                supported_subjects = self.metadata.capabilities['subjects']
                if context.subject not in supported_subjects:
                    return False, f"不支持的科目：{context.subject}。支持的科目：{', '.join(supported_subjects)}"
                
                return True, None
            ```
        """
        pass
    
    @abstractmethod
    async def generate_learning_objectives(
        self,
        context: WorkflowContext
    ) -> List[LearningObjective]:
        """
        生成学习目标
        
        这是AI驱动的核心功能之一。适配器需要：
        1. 准备领域特定的Prompt
        2. 调用AI服务生成目标
        3. 应用领域规则筛选和优化
        
        Args:
            context: 工作流上下文
            
        Returns:
            List[LearningObjective]: 学习目标列表
            
        示例（K12适配器）：
            ```python
            async def generate_learning_objectives(self, context: WorkflowContext) -> List[LearningObjective]:
                # 准备Prompt
                prompt = self._build_objective_prompt(
                    subject=context.subject,
                    grade=context.grade_level,
                    topic=context.raw_input,
                    duration=context.duration
                )
                
                # 调用AI（通过平台提供的LLM客户端）
                from src.ai import llm_client
                response = await llm_client.generate_json([
                    {"role": "system", "content": prompt['system']},
                    {"role": "user", "content": prompt['user']}
                ])
                
                # 解析和验证
                objectives = []
                for obj_data in response.get('objectives', []):
                    obj = LearningObjective(
                        title=obj_data['title'],
                        bloom_level=BloomLevel(obj_data['bloom_level']),
                        estimated_time=obj_data['estimated_time']
                    )
                    # 应用K12特定规则验证
                    if self._validate_objective_for_grade(obj, context.grade_level):
                        objectives.append(obj)
                
                return objectives
            ```
        """
        pass
    
    # ==================== 更多接口方法待续... ====================
```

---

*文档第一部分完成，共约3000行*

**已完成内容**：
- ✅ 架构模型（三层架构、适配器结构）
- ✅ 适配器清单规范（manifest.yaml完整示例）
- ✅ DomainAdapter核心接口（元数据、生命周期、核心方法）

**下一部分将包含**：
- DomainAdapter 剩余接口方法
- NodeExecutor 接口
- ContentSource 接口  
- 适配器注册机制详细设计

是否继续输出第二部分？
# MetaWorkflow V2.0 适配器市场设计 - 第二部分

## 1.2 接口规范（续）

### 1.2.2 DomainAdapter 完整接口

```python
# 接续第一部分...

class DomainAdapter(ABC):
    """领域适配器抽象基类（续）"""
    
    # ==================== 核心业务方法（续） ====================
    
    @abstractmethod
    async def discover_content(
        self,
        context: WorkflowContext,
        objectives: List[LearningObjective]
    ) -> Dict[str, Any]:
        """
        发现和检索相关内容
        
        这是AI驱动的内容发现功能。适配器需要：
        1. 根据学习目标确定检索策略
        2. 从配置的内容源检索素材
        3. 使用AI筛选和排序结果
        4. 返回结构化的内容集合
        
        Args:
            context: 工作流上下文
            objectives: 学习目标列表
            
        Returns:
            Dict[str, Any]: 检索到的内容素材，结构化组织
            
        推荐结构：
            {
                "concepts": {          # 概念定义
                    "items": [...],
                    "source": "Wikipedia",
                    "confidence": 0.95
                },
                "examples": {          # 示例
                    "items": [...],
                    "count": 5
                },
                "media": {             # 多媒体资源
                    "images": [...],
                    "videos": [...],
                    "diagrams": [...]
                },
                "references": {        # 参考资料
                    "books": [...],
                    "articles": [...],
                    "websites": [...]
                },
                "related_topics": [...] # 相关主题
            }
            
        示例（美术史适配器）：
            ```python
            async def discover_content(self, context, objectives):
                content = {}
                
                # 从WikiArt检索画作
                if 'paintings' in context.metadata.get('required_content', []):
                    from .services import wikiart_client
                    
                    paintings = await wikiart_client.search(
                        artist=context.metadata.get('artist'),
                        style=context.metadata.get('style'),
                        period=context.metadata.get('period'),
                        limit=10
                    )
                    
                    # 使用AI筛选最相关的画作
                    from src.ai import llm_client
                    selected = await self._ai_select_paintings(
                        paintings, objectives, llm_client
                    )
                    
                    content['paintings'] = {
                        'items': selected,
                        'source': 'WikiArt',
                        'count': len(selected)
                    }
                
                # 从学术数据库检索背景知识
                if objectives[0].bloom_level in [BloomLevel.ANALYZE, BloomLevel.EVALUATE]:
                    from .services import jstor_client
                    
                    articles = await jstor_client.search(
                        query=context.raw_input,
                        subject='Art History',
                        limit=5
                    )
                    
                    content['scholarly_articles'] = {
                        'items': articles,
                        'source': 'JSTOR'
                    }
                
                return content
            ```
        """
        pass
    
    @abstractmethod
    async def customize_workflow(
        self,
        context: WorkflowContext,
        base_workflow: Dict[str, Any]
    ) -> Dict[str, Any]:
        """
        自定义工作流配置
        
        根据领域特点和上下文调整工作流。适配器可以：
        1. 添加/删除/重排节点
        2. 调整节点参数
        3. 设置节点依赖关系
        4. 配置并行策略
        
        Args:
            context: 工作流上下文
            base_workflow: 平台提供的基础工作流模板
                {
                    "nodes": [
                        {
                            "id": "objective",
                            "type": "ai_objective_generator",
                            "config": {...},
                            "dependencies": []
                        },
                        ...
                    ],
                    "metadata": {...}
                }
            
        Returns:
            Dict[str, Any]: 定制后的工作流配置
            
        示例（K12微课适配器）：
            ```python
            async def customize_workflow(self, context, base_workflow):
                workflow = base_workflow.copy()
                
                # 对于10分钟微课，简化工作流
                if context.duration <= 15:
                    workflow['nodes'] = [
                        {
                            "id": "parse",
                            "type": "adapter_parse",
                            "config": {},
                            "dependencies": []
                        },
                        {
                            "id": "objective",
                            "type": "ai_objective",
                            "config": {
                                "max_objectives": 2,  # 微课最多2个目标
                                "prefer_bloom_levels": ["understand", "apply"]
                            },
                            "dependencies": ["parse"]
                        },
                        {
                            "id": "concept",
                            "type": "concept_activity",
                            "config": {
                                "include_examples": True,
                                "include_practice": False  # 微课不包含练习
                            },
                            "dependencies": ["objective"]
                        },
                        # 省略实验和问题节点
                    ]
                    
                    workflow['metadata']['course_type'] = 'micro'
                
                # 对于化学实验课，添加安全检查节点
                elif context.subject == '化学' and context.duration >= 45:
                    safety_node = {
                        "id": "safety_check",
                        "type": "chemistry_safety",
                        "config": {
                            "check_materials": True,
                            "generate_warnings": True
                        },
                        "dependencies": ["experiment"]
                    }
                    workflow['nodes'].append(safety_node)
                
                return workflow
            ```
        """
        pass
    
    @abstractmethod
    def get_supported_activity_types(self) -> List[ActivityType]:
        """
        获取支持的活动类型
        
        Returns:
            List[ActivityType]: 活动类型列表
            
        示例：
            [
                ActivityType.CONCEPT,      # 概念理解
                ActivityType.EXPERIMENT,   # 实验探究
                ActivityType.PROBLEM       # 问题解决
            ]
        """
        pass
    
    @abstractmethod
    async def format_output(
        self,
        context: WorkflowContext,
        activities: List[Any]
    ) -> Dict[str, Any]:
        """
        格式化最终输出
        
        将生成的活动转换为用户期望的格式。适配器负责：
        1. 选择合适的输出格式（PPT/PDF/HTML等）
        2. 应用领域特定的模板和样式
        3. 生成预览和下载链接
        4. 添加元数据
        
        Args:
            context: 工作流上下文
            activities: 生成的活动列表
                [
                    ConceptActivity(...),
                    ExperimentActivity(...),
                    ProblemActivity(...)
                ]
            
        Returns:
            Dict[str, Any]: 格式化后的输出
                {
                    "format": "ppt",
                    "content": {...},           # 格式相关的内容
                    "metadata": {
                        "title": "...",
                        "author": "...",
                        "created_at": "...",
                        "page_count": 25
                    },
                    "files": {
                        "preview_url": "https://...",
                        "download_url": "https://...",
                        "size": 2.5  # MB
                    }
                }
                
        示例（K12 PPT适配器）：
            ```python
            async def format_output(self, context, activities):
                from .formatters import PPTFormatter
                
                formatter = PPTFormatter()
                
                # 生成PPT
                ppt_bytes = await formatter.create(
                    title=f"{context.subject} - {context.raw_input}",
                    grade=context.grade_level,
                    duration=f"{context.duration}分钟",
                    activities=activities,
                    template='k12_default'
                )
                
                # 上传到对象存储
                from src.storage import storage_service
                
                file_id = await storage_service.upload(
                    content=ppt_bytes,
                    filename=f"{context.workflow_id}.pptx",
                    content_type='application/vnd.openxmlformats-officedocument.presentationml.presentation'
                )
                
                return {
                    "format": "ppt",
                    "content": {
                        "slides": len(activities) + 2,  # 活动数 + 封面 + 结束页
                        "theme": "k12_default"
                    },
                    "metadata": {
                        "title": f"{context.subject}教案",
                        "subject": context.subject,
                        "grade": f"{context.grade_level}年级",
                        "duration": f"{context.duration}分钟",
                        "created_at": datetime.utcnow().isoformat(),
                        "generator": f"MetaWorkflow K12 Adapter v{self.metadata.version}"
                    },
                    "files": {
                        "file_id": file_id,
                        "preview_url": await storage_service.get_preview_url(file_id),
                        "download_url": await storage_service.get_download_url(file_id),
                        "size": len(ppt_bytes) / (1024 * 1024)  # MB
                    }
                }
            ```
        """
        pass
    
    # ==================== 可选方法（有默认实现） ====================
    
    async def assess_quality(
        self,
        context: WorkflowContext,
        output: Dict[str, Any]
    ) -> Dict[str, Any]:
        """
        质量评估（可选）
        
        对生成的内容进行质量评估。适配器可以：
        1. 使用AI评估内容质量
        2. 应用领域特定的评估规则
        3. 给出改进建议
        
        Args:
            context: 工作流上下文
            output: 格式化后的输出
            
        Returns:
            Dict[str, Any]: 质量评估结果
                {
                    "overall_score": 0.85,  # 0-1之间的总分
                    "dimensions": {
                        "accuracy": 0.9,      # 准确性
                        "relevance": 0.8,     # 相关性
                        "clarity": 0.85,      # 清晰度
                        "completeness": 0.9   # 完整性
                    },
                    "suggestions": [
                        "建议添加更多示例以增强理解",
                        "实验步骤描述可以更详细"
                    ],
                    "compliance": {         # 合规性检查
                        "curriculum_standard": True,
                        "age_appropriate": True,
                        "safety_checked": True
                    }
                }
        """
        # 默认实现：返回满分
        return {
            "overall_score": 1.0,
            "dimensions": {},
            "suggestions": [],
            "compliance": {}
        }
    
    async def get_configuration_schema(self) -> Dict[str, Any]:
        """
        获取配置参数的JSON Schema（可选）
        
        返回适配器可配置参数的Schema，用于：
        1. 自动生成配置UI
        2. 验证用户配置
        3. 生成文档
        
        Returns:
            Dict[str, Any]: JSON Schema定义
        """
        # 从manifest.yaml读取
        return self.metadata.configuration.get('schema', {
            "type": "object",
            "properties": {},
            "required": []
        })
    
    async def on_install(self) -> None:
        """
        安装钩子（可选）
        
        适配器首次安装时调用。可用于：
        1. 下载依赖资源
        2. 初始化数据库
        3. 创建配置文件
        """
        pass
    
    async def on_update(self, old_version: str) -> None:
        """
        更新钩子（可选）
        
        适配器更新时调用。可用于：
        1. 迁移数据
        2. 更新配置
        3. 清理旧资源
        
        Args:
            old_version: 旧版本号
        """
        pass
    
    async def on_uninstall(self) -> None:
        """
        卸载钩子（可选）
        
        适配器卸载时调用。可用于：
        1. 清理数据
        2. 删除临时文件
        3. 撤销配置
        """
        pass
    
    # ==================== 健康检查 ====================
    
    async def health_check(self) -> Dict[str, Any]:
        """
        健康检查
        
        定期调用以检查适配器状态
        
        Returns:
            Dict[str, Any]: 健康状态
                {
                    "status": "healthy" | "degraded" | "unhealthy",
                    "checks": {
                        "dependencies": "ok",
                        "content_sources": "ok",
                        "resources": "ok"
                    },
                    "metrics": {
                        "requests_count": 1234,
                        "avg_latency": 0.5,  # 秒
                        "error_rate": 0.01   # 1%
                    }
                }
        """
        return {
            "status": "healthy",
            "checks": {},
            "metrics": {}
        }
```

### 1.2.3 NodeExecutor 接口

```python
# src/core/node_executor.py

from abc import ABC, abstractmethod
from typing import Any, Dict
from src.models import WorkflowNode, WorkflowContext

class NodeExecutor(ABC):
    """
    节点执行器抽象基类
    
    每种工作流节点类型需要一个执行器来实现其业务逻辑。
    执行器与适配器配合，实现具体的内容生成功能。
    """
    
    @property
    @abstractmethod
    def node_type(self) -> str:
        """
        节点类型标识
        
        Returns:
            str: 如 "ai_objective_generator", "concept_activity"
        """
        pass
    
    @abstractmethod
    async def execute(
        self,
        node: WorkflowNode,
        context: WorkflowContext,
        input_data: Dict[str, Any]
    ) -> Any:
        """
        执行节点
        
        Args:
            node: 节点配置
            context: 工作流上下文
            input_data: 输入数据（包含依赖节点的输出）
            
        Returns:
            Any: 节点执行结果
        """
        pass


class AdapterParseNodeExecutor(NodeExecutor):
    """适配器解析节点执行器"""
    
    @property
    def node_type(self) -> str:
        return "adapter_parse"
    
    async def execute(self, node, context, input_data):
        """调用适配器的parse_request方法"""
        from src.adapters import adapter_registry
        
        adapter = adapter_registry.get_adapter(context.domain)
        # 解析已在前面完成，这里做验证
        is_valid, error = await adapter.validate_context(context)
        
        if not is_valid:
            raise ValueError(f"上下文验证失败: {error}")
        
        return {"context": context.model_dump()}


class AIObjectiveGeneratorNodeExecutor(NodeExecutor):
    """AI学习目标生成节点执行器"""
    
    @property
    def node_type(self) -> str:
        return "ai_objective_generator"
    
    async def execute(self, node, context, input_data):
        """调用适配器的generate_learning_objectives方法"""
        from src.adapters import adapter_registry
        
        adapter = adapter_registry.get_adapter(context.domain)
        objectives = await adapter.generate_learning_objectives(context)
        
        # 更新上下文
        context.learning_objectives = objectives
        
        return {
            "objectives": [obj.model_dump() for obj in objectives],
            "count": len(objectives)
        }
```

### 1.2.4 ContentSource 接口

```python
# src/content/sources/base.py

from abc import ABC, abstractmethod
from typing import List, Dict, Any, Optional
from enum import Enum

class ContentType(Enum):
    """内容类型"""
    TEXT = "text"               # 文本
    IMAGE = "image"             # 图片
    VIDEO = "video"             # 视频
    AUDIO = "audio"             # 音频
    DOCUMENT = "document"       # 文档
    DATASET = "dataset"         # 数据集


class ContentSource(ABC):
    """
    内容源抽象基类
    
    内容源负责从外部获取教学素材，如：
    - Wikipedia
    - YouTube
    - WikiArt
    - JSTOR
    - 企业知识库
    """
    
    @property
    @abstractmethod
    def source_name(self) -> str:
        """
        内容源名称
        
        Returns:
            str: 如 "wikipedia", "youtube", "wikiart"
        """
        pass
    
    @property
    @abstractmethod
    def supported_content_types(self) -> List[ContentType]:
        """
        支持的内容类型
        
        Returns:
            List[ContentType]: 如 [ContentType.TEXT, ContentType.IMAGE]
        """
        pass
    
    @abstractmethod
    async def search(
        self,
        query: str,
        content_type: Optional[ContentType] = None,
        limit: int = 10,
        **kwargs
    ) -> List[Dict[str, Any]]:
        """
        搜索内容
        
        Args:
            query: 搜索查询
            content_type: 内容类型过滤
            limit: 结果数量限制
            **kwargs: 源特定的参数
            
        Returns:
            List[Dict[str, Any]]: 搜索结果列表
                [
                    {
                        "id": "...",
                        "title": "...",
                        "content": "...",
                        "url": "...",
                        "metadata": {...}
                    },
                    ...
                ]
        """
        pass
    
    async def get_by_id(self, content_id: str) -> Optional[Dict[str, Any]]:
        """
        根据ID获取内容
        
        Args:
            content_id: 内容ID
            
        Returns:
            Optional[Dict[str, Any]]: 内容详情
        """
        pass
    
    async def validate_credentials(self) -> bool:
        """
        验证API凭证
        
        Returns:
            bool: 凭证是否有效
        """
        return True


# 示例实现

class WikipediaSource(ContentSource):
    """Wikipedia内容源"""
    
    def __init__(self, language: str = "zh"):
        self.language = language
    
    @property
    def source_name(self) -> str:
        return "wikipedia"
    
    @property
    def supported_content_types(self) -> List[ContentType]:
        return [ContentType.TEXT, ContentType.IMAGE]
    
    async def search(self, query, content_type=None, limit=10, **kwargs):
        import wikipedia
        
        wikipedia.set_lang(self.language)
        
        # 搜索页面
        search_results = wikipedia.search(query, results=limit)
        
        contents = []
        for title in search_results:
            try:
                page = wikipedia.page(title)
                contents.append({
                    "id": page.pageid,
                    "title": page.title,
                    "content": page.summary,
                    "url": page.url,
                    "metadata": {
                        "source": "wikipedia",
                        "language": self.language,
                        "categories": page.categories
                    }
                })
            except:
                continue
        
        return contents
```

---

*文档第二部分完成*

**已完成内容**：
- ✅ DomainAdapter 完整接口（所有核心和可选方法）
- ✅ NodeExecutor 接口规范
- ✅ ContentSource 接口规范
- ✅ 示例实现

**下一部分将包含**：
- 适配器注册机制详细设计
- 依赖管理和版本控制
- 适配器发现和加载流程

是否继续输出第三部分？
# MetaWorkflow V2.0 适配器市场设计 - 第三部分

## 第二部分：适配器注册机制

## 2.1 注册表设计

### 2.1.1 注册表架构

```python
# src/adapters/registry.py

from typing import Dict, List, Optional, Type
from datetime import datetime
import importlib
import yaml
from pathlib import Path

from .base import DomainAdapter, AdapterMetadata
from src.models import DomainType

class AdapterRegistry:
    """
    适配器注册表
    
    职责：
    1. 管理所有已注册的适配器
    2. 提供适配器查找和获取功能
    3. 处理适配器的生命周期
    4. 验证适配器兼容性
    """
    
    def __init__(self):
        """初始化注册表"""
        self._adapters: Dict[DomainType, DomainAdapter] = {}
        self._metadata: Dict[DomainType, AdapterMetadata] = {}
        self._instances: Dict[str, DomainAdapter] = {}  # adapter_id -> instance
        self._load_history: List[Dict] = []
    
    def register(
        self, 
        domain: DomainType, 
        adapter: DomainAdapter,
        force: bool = False
    ) -> None:
        """
        注册适配器
        
        Args:
            domain: 领域类型
            adapter: 适配器实例
            force: 是否强制覆盖已存在的适配器
            
        Raises:
            ValueError: 适配器已存在且force=False
            RuntimeError: 适配器验证失败
        """
        # 检查是否已注册
        if domain in self._adapters and not force:
            raise ValueError(
                f"Domain {domain.value} already has a registered adapter. "
                f"Use force=True to override."
            )
        
        # 验证兼容性
        from src import __version__ as platform_version
        is_compatible, error = adapter.validate_compatibility(platform_version)
        if not is_compatible:
            raise RuntimeError(f"Adapter incompatible: {error}")
        
        # 获取元数据
        metadata = adapter.metadata
        
        # 验证元数据
        self._validate_metadata(metadata)
        
        # 注册
        self._adapters[domain] = adapter
        self._metadata[domain] = metadata
        self._instances[metadata.name] = adapter
        
        # 记录加载历史
        self._load_history.append({
            "domain": domain.value,
            "adapter": metadata.name,
            "version": metadata.version,
            "loaded_at": datetime.utcnow().isoformat(),
            "action": "register"
        })
        
        logger.info(
            f"Registered adapter: {metadata.display_name} v{metadata.version} "
            f"for domain {domain.value}"
        )
    
    def unregister(self, domain: DomainType) -> None:
        """
        注销适配器
        
        Args:
            domain: 领域类型
        """
        if domain not in self._adapters:
            return
        
        adapter = self._adapters[domain]
        metadata = self._metadata[domain]
        
        # 调用卸载钩子
        import asyncio
        asyncio.create_task(adapter.unload())
        
        # 从注册表移除
        del self._adapters[domain]
        del self._metadata[domain]
        if metadata.name in self._instances:
            del self._instances[metadata.name]
        
        # 记录历史
        self._load_history.append({
            "domain": domain.value,
            "adapter": metadata.name,
            "version": metadata.version,
            "loaded_at": datetime.utcnow().isoformat(),
            "action": "unregister"
        })
        
        logger.info(f"Unregistered adapter for domain {domain.value}")
    
    def get_adapter(self, domain: DomainType) -> Optional[DomainAdapter]:
        """
        获取适配器
        
        Args:
            domain: 领域类型
            
        Returns:
            Optional[DomainAdapter]: 适配器实例，不存在则返回None
        """
        return self._adapters.get(domain)
    
    def get_metadata(self, domain: DomainType) -> Optional[AdapterMetadata]:
        """
        获取适配器元数据
        
        Args:
            domain: 领域类型
            
        Returns:
            Optional[AdapterMetadata]: 元数据，不存在则返回None
        """
        return self._metadata.get(domain)
    
    def list_adapters(self) -> List[Dict[str, Any]]:
        """
        列出所有已注册的适配器
        
        Returns:
            List[Dict[str, Any]]: 适配器信息列表
        """
        adapters = []
        for domain, metadata in self._metadata.items():
            adapters.append({
                "domain": domain.value,
                "name": metadata.name,
                "display_name": metadata.display_name,
                "version": metadata.version,
                "author": metadata.author,
                "description": metadata.description,
                "capabilities": metadata.capabilities,
                "status": "active"
            })
        return adapters
    
    def _validate_metadata(self, metadata: AdapterMetadata) -> None:
        """
        验证适配器元数据
        
        Args:
            metadata: 适配器元数据
            
        Raises:
            ValueError: 元数据不合法
        """
        # 必填字段检查
        required_fields = ['name', 'display_name', 'version', 'description']
        for field in required_fields:
            if not getattr(metadata, field, None):
                raise ValueError(f"Missing required metadata field: {field}")
        
        # 版本格式检查（语义化版本）
        import re
        version_pattern = r'^\d+\.\d+\.\d+$'
        if not re.match(version_pattern, metadata.version):
            raise ValueError(
                f"Invalid version format: {metadata.version}. "
                f"Expected semantic version (e.g., 1.0.0)"
            )
    
    async def initialize_all(self) -> None:
        """初始化所有已注册的适配器"""
        for adapter in self._adapters.values():
            if not adapter._initialized:
                await adapter.initialize()
    
    async def health_check_all(self) -> Dict[str, Any]:
        """
        检查所有适配器的健康状态
        
        Returns:
            Dict[str, Any]: 健康状态汇总
        """
        results = {}
        for domain, adapter in self._adapters.items():
            try:
                health = await adapter.health_check()
                results[domain.value] = {
                    "status": health.get("status", "unknown"),
                    "details": health
                }
            except Exception as e:
                results[domain.value] = {
                    "status": "error",
                    "error": str(e)
                }
        
        return {
            "overall_status": self._calculate_overall_status(results),
            "adapters": results,
            "checked_at": datetime.utcnow().isoformat()
        }
    
    def _calculate_overall_status(self, results: Dict) -> str:
        """计算整体健康状态"""
        statuses = [r.get("status") for r in results.values()]
        
        if all(s == "healthy" for s in statuses):
            return "healthy"
        elif any(s == "unhealthy" for s in statuses):
            return "unhealthy"
        else:
            return "degraded"


# 全局注册表实例
adapter_registry = AdapterRegistry()
```

### 2.1.2 自动发现机制

```python
# src/adapters/discovery.py

from pathlib import Path
from typing import List, Dict, Any
import importlib.util
import sys
import yaml

class AdapterDiscovery:
    """
    适配器自动发现
    
    支持多种发现方式：
    1. 从指定目录扫描
    2. 从Python包导入
    3. 从适配器市场下载
    """
    
    def __init__(self, search_paths: List[str] = None):
        """
        初始化发现器
        
        Args:
            search_paths: 适配器搜索路径列表
        """
        self.search_paths = search_paths or [
            "src/adapters",
            "~/.metaworkflow/adapters",
            "/opt/metaworkflow/adapters"
        ]
    
    def discover_all(self) -> List[Dict[str, Any]]:
        """
        发现所有可用的适配器
        
        Returns:
            List[Dict[str, Any]]: 发现的适配器信息列表
        """
        discovered = []
        
        for search_path in self.search_paths:
            path = Path(search_path).expanduser()
            if not path.exists():
                continue
            
            # 扫描目录
            for adapter_dir in path.iterdir():
                if not adapter_dir.is_dir():
                    continue
                
                # 检查是否包含manifest.yaml
                manifest_file = adapter_dir / "manifest.yaml"
                if not manifest_file.exists():
                    continue
                
                # 读取manifest
                with open(manifest_file, 'r', encoding='utf-8') as f:
                    manifest = yaml.safe_load(f)
                
                discovered.append({
                    "path": str(adapter_dir),
                    "manifest": manifest,
                    "source": "local"
                })
        
        return discovered
    
    def load_adapter_from_path(
        self, 
        adapter_path: str
    ) -> DomainAdapter:
        """
        从路径加载适配器
        
        Args:
            adapter_path: 适配器目录路径
            
        Returns:
            DomainAdapter: 适配器实例
            
        Raises:
            ImportError: 无法加载适配器
        """
        adapter_path = Path(adapter_path)
        
        # 读取manifest
        manifest_file = adapter_path / "manifest.yaml"
        with open(manifest_file, 'r', encoding='utf-8') as f:
            manifest = yaml.safe_load(f)
        
        # 加载adapter.py模块
        adapter_module_file = adapter_path / "adapter.py"
        if not adapter_module_file.exists():
            raise ImportError(f"adapter.py not found in {adapter_path}")
        
        # 动态导入
        spec = importlib.util.spec_from_file_location(
            manifest['name'],
            adapter_module_file
        )
        module = importlib.util.module_from_spec(spec)
        sys.modules[manifest['name']] = module
        spec.loader.exec_module(module)
        
        # 查找适配器类（约定：模块中只有一个DomainAdapter子类）
        adapter_class = None
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            if (isinstance(attr, type) and 
                issubclass(attr, DomainAdapter) and 
                attr is not DomainAdapter):
                adapter_class = attr
                break
        
        if adapter_class is None:
            raise ImportError(
                f"No DomainAdapter subclass found in {adapter_module_file}"
            )
        
        # 实例化
        adapter_instance = adapter_class()
        
        return adapter_instance
    
    def auto_register_all(self) -> Dict[str, Any]:
        """
        自动发现并注册所有适配器
        
        Returns:
            Dict[str, Any]: 注册结果统计
        """
        from .registry import adapter_registry
        
        discovered = self.discover_all()
        
        stats = {
            "discovered": len(discovered),
            "registered": 0,
            "failed": 0,
            "skipped": 0,
            "details": []
        }
        
        for item in discovered:
            try:
                manifest = item['manifest']
                
                # 检查是否已注册
                domain_name = manifest['name']
                domain_type = DomainType(domain_name)
                
                if adapter_registry.get_adapter(domain_type):
                    stats['skipped'] += 1
                    stats['details'].append({
                        "adapter": manifest['display_name'],
                        "status": "skipped",
                        "reason": "Already registered"
                    })
                    continue
                
                # 加载适配器
                adapter = self.load_adapter_from_path(item['path'])
                
                # 注册
                adapter_registry.register(domain_type, adapter)
                
                stats['registered'] += 1
                stats['details'].append({
                    "adapter": manifest['display_name'],
                    "version": manifest['version'],
                    "status": "registered"
                })
                
            except Exception as e:
                stats['failed'] += 1
                stats['details'].append({
                    "adapter": item.get('manifest', {}).get('display_name', 'Unknown'),
                    "status": "failed",
                    "error": str(e)
                })
        
        return stats


# 全局发现器实例
adapter_discovery = AdapterDiscovery()
```

### 2.1.3 依赖管理

```python
# src/adapters/dependency.py

from typing import List, Dict, Set, Any
from packaging import version, specifiers
import networkx as nx

class DependencyResolver:
    """
    适配器依赖解析器
    
    功能：
    1. 解析适配器依赖关系
    2. 检测循环依赖
    3. 计算安装顺序
    4. 验证版本兼容性
    """
    
    def __init__(self):
        self.dependency_graph = nx.DiGraph()
    
    def add_adapter(
        self, 
        adapter_name: str, 
        dependencies: Dict[str, str]
    ) -> None:
        """
        添加适配器到依赖图
        
        Args:
            adapter_name: 适配器名称
            dependencies: 依赖字典 {adapter_name: version_spec}
        """
        # 添加节点
        if adapter_name not in self.dependency_graph:
            self.dependency_graph.add_node(adapter_name)
        
        # 添加边
        for dep_name, version_spec in dependencies.items():
            self.dependency_graph.add_edge(
                adapter_name, 
                dep_name,
                version_spec=version_spec
            )
    
    def check_circular_dependencies(self) -> List[List[str]]:
        """
        检查循环依赖
        
        Returns:
            List[List[str]]: 循环依赖链列表，空列表表示无循环
        """
        try:
            cycles = list(nx.simple_cycles(self.dependency_graph))
            return cycles
        except:
            return []
    
    def get_install_order(
        self, 
        adapter_names: List[str]
    ) -> List[str]:
        """
        计算安装顺序（拓扑排序）
        
        Args:
            adapter_names: 要安装的适配器列表
            
        Returns:
            List[str]: 按依赖顺序排列的适配器名称列表
            
        Raises:
            ValueError: 存在循环依赖
        """
        # 检查循环依赖
        cycles = self.check_circular_dependencies()
        if cycles:
            raise ValueError(f"Circular dependencies detected: {cycles}")
        
        # 获取子图
        subgraph = self.dependency_graph.subgraph(adapter_names)
        
        # 拓扑排序
        try:
            install_order = list(nx.topological_sort(subgraph))
            # 反转（依赖项在前）
            return list(reversed(install_order))
        except nx.NetworkXError as e:
            raise ValueError(f"Cannot determine install order: {e}")
    
    def verify_compatibility(
        self,
        adapter_name: str,
        adapter_version: str,
        installed_adapters: Dict[str, str]
    ) -> Tuple[bool, List[str]]:
        """
        验证版本兼容性
        
        Args:
            adapter_name: 适配器名称
            adapter_version: 适配器版本
            installed_adapters: 已安装的适配器 {name: version}
            
        Returns:
            Tuple[bool, List[str]]: (是否兼容, 不兼容原因列表)
        """
        issues = []
        
        # 获取依赖
        if adapter_name not in self.dependency_graph:
            return True, []
        
        dependencies = self.dependency_graph[adapter_name]
        
        for dep_name in dependencies:
            edge_data = self.dependency_graph[adapter_name][dep_name]
            required_version_spec = edge_data.get('version_spec', '*')
            
            # 检查是否已安装
            if dep_name not in installed_adapters:
                issues.append(
                    f"Missing dependency: {dep_name} {required_version_spec}"
                )
                continue
            
            # 检查版本是否兼容
            installed_version = installed_adapters[dep_name]
            
            try:
                spec = specifiers.SpecifierSet(required_version_spec)
                if not spec.contains(installed_version):
                    issues.append(
                        f"Version mismatch: {dep_name} requires {required_version_spec}, "
                        f"but {installed_version} is installed"
                    )
            except Exception as e:
                issues.append(
                    f"Invalid version specifier for {dep_name}: {e}"
                )
        
        return len(issues) == 0, issues


# 全局依赖解析器
dependency_resolver = DependencyResolver()
```

---

## 2.2 适配器生命周期

### 2.2.1 生命周期状态机

```
┌─────────────┐
│  DISCOVERED │  发现：在文件系统或市场中被发现
└──────┬──────┘
       │
       │ 下载/验证
       ▼
┌─────────────┐
│  DOWNLOADED │  已下载：文件已下载到本地
└──────┬──────┘
       │
       │ 安装依赖
       ▼
┌─────────────┐
│  INSTALLED  │  已安装：依赖已安装，但未加载
└──────┬──────┘
       │
       │ 加载模块
       ▼
┌─────────────┐
│   LOADED    │  已加载：模块已导入，但未初始化
└──────┬──────┘
       │
       │ 初始化
       ▼
┌─────────────┐
│   ACTIVE    │  激活：正在使用中
└──────┬──────┘
       │
       │ 停用
       ▼
┌─────────────┐
│  INACTIVE   │  停用：已停止，但未卸载
└──────┬──────┘
       │
       │ 卸载
       ▼
┌─────────────┐
│ UNINSTALLED │  已卸载：已从系统移除
└─────────────┘
```

### 2.2.2 生命周期管理器

```python
# src/adapters/lifecycle.py

from enum import Enum
from typing import Dict, Optional
from datetime import datetime

class AdapterStatus(Enum):
    """适配器状态"""
    DISCOVERED = "discovered"
    DOWNLOADED = "downloaded"
    INSTALLED = "installed"
    LOADED = "loaded"
    ACTIVE = "active"
    INACTIVE = "inactive"
    UNINSTALLED = "uninstalled"
    ERROR = "error"


class AdapterLifecycleManager:
    """适配器生命周期管理器"""
    
    def __init__(self):
        self._status: Dict[str, AdapterStatus] = {}
        self._history: Dict[str, List[Dict]] = {}
    
    async def install(
        self, 
        adapter_name: str,
        source: str
    ) -> bool:
        """
        安装适配器
        
        步骤：
        1. 下载（如果是远程）
        2. 验证签名和完整性
        3. 安装依赖
        4. 运行安装钩子
        5. 更新状态
        
        Args:
            adapter_name: 适配器名称
            source: 来源（local/marketplace/github）
            
        Returns:
            bool: 是否成功
        """
        try:
            self._update_status(adapter_name, AdapterStatus.DISCOVERED)
            
            # 1. 下载
            if source != "local":
                await self._download_adapter(adapter_name, source)
            self._update_status(adapter_name, AdapterStatus.DOWNLOADED)
            
            # 2. 验证
            if not await self._verify_adapter(adapter_name):
                raise RuntimeError("Adapter verification failed")
            
            # 3. 安装依赖
            await self._install_dependencies(adapter_name)
            self._update_status(adapter_name, AdapterStatus.INSTALLED)
            
            # 4. 运行安装钩子
            adapter = self._load_adapter_module(adapter_name)
            await adapter.on_install()
            
            return True
            
        except Exception as e:
            self._update_status(adapter_name, AdapterStatus.ERROR)
            self._record_error(adapter_name, str(e))
            return False
    
    async def activate(self, adapter_name: str) -> bool:
        """激活适配器"""
        try:
            # 加载
            adapter = self._load_adapter_module(adapter_name)
            self._update_status(adapter_name, AdapterStatus.LOADED)
            
            # 初始化
            await adapter.initialize()
            self._update_status(adapter_name, AdapterStatus.ACTIVE)
            
            # 注册
            from .registry import adapter_registry
            domain = self._get_adapter_domain(adapter_name)
            adapter_registry.register(domain, adapter)
            
            return True
            
        except Exception as e:
            self._update_status(adapter_name, AdapterStatus.ERROR)
            return False
    
    async def deactivate(self, adapter_name: str) -> bool:
        """停用适配器"""
        try:
            from .registry import adapter_registry
            
            domain = self._get_adapter_domain(adapter_name)
            adapter_registry.unregister(domain)
            
            self._update_status(adapter_name, AdapterStatus.INACTIVE)
            return True
            
        except Exception as e:
            return False
    
    async def uninstall(self, adapter_name: str) -> bool:
        """卸载适配器"""
        try:
            # 先停用
            await self.deactivate(adapter_name)
            
            # 运行卸载钩子
            adapter = self._load_adapter_module(adapter_name)
            await adapter.on_uninstall()
            
            # 删除文件
            self._remove_adapter_files(adapter_name)
            
            self._update_status(adapter_name, AdapterStatus.UNINSTALLED)
            return True
            
        except Exception as e:
            return False
    
    def get_status(self, adapter_name: str) -> Optional[AdapterStatus]:
        """获取适配器状态"""
        return self._status.get(adapter_name)
    
    def _update_status(
        self, 
        adapter_name: str, 
        status: AdapterStatus
    ) -> None:
        """更新状态"""
        self._status[adapter_name] = status
        
        # 记录历史
        if adapter_name not in self._history:
            self._history[adapter_name] = []
        
        self._history[adapter_name].append({
            "status": status.value,
            "timestamp": datetime.utcnow().isoformat()
        })


# 全局生命周期管理器
lifecycle_manager = AdapterLifecycleManager()
```

---

*文档第三部分完成*

**已完成内容**：
- ✅ 注册表设计（AdapterRegistry）
- ✅ 自动发现机制（AdapterDiscovery）
- ✅ 依赖管理（DependencyResolver）
- ✅ 生命周期管理（AdapterLifecycleManager）

**下一部分将包含**：
- 适配器市场架构设计
- 发布和审核流程
- 用户安装和评价系统

是否继续输出第四部分（适配器市场设计）？
# MetaWorkflow V2.0 适配器市场设计 - 第四部分

## 第三部分：适配器市场平台

## 3.1 市场架构

### 3.1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户层                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  开发者控制台 │  │   用户商店   │  │  管理后台    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API网关层                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 认证授权  │  │ 限流控制  │  │ 日志审计  │  │ API版本  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      业务服务层                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  适配器管理  │  │   发布审核   │  │   搜索推荐   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  用户管理    │  │   评分评论   │  │   统计分析   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      数据存储层                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  PostgreSQL │  │   Redis     │  │   S3存储    │             │
│  │  (元数据)    │  │  (缓存)     │  │  (文件)     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐                               │
│  │ Elasticsearch│  │   ClickHouse│                               │
│  │  (全文搜索)  │  │  (日志分析)  │                               │
│  └─────────────┘  └─────────────┘                               │
└─────────────────────────────────────────────────────────────────┘
```

### 3.1.2 数据模型设计

```python
# src/marketplace/models.py

from sqlalchemy import Column, Integer, String, Text, DateTime, Boolean, Float, JSON, ForeignKey, Enum
from sqlalchemy.orm import relationship
from datetime import datetime
import enum

from src.db.database import Base


class AdapterPublishStatus(enum.Enum):
    """发布状态"""
    DRAFT = "draft"              # 草稿
    SUBMITTED = "submitted"      # 已提交
    REVIEWING = "reviewing"      # 审核中
    APPROVED = "approved"        # 已批准
    REJECTED = "rejected"        # 已拒绝
    PUBLISHED = "published"      # 已发布
    DEPRECATED = "deprecated"    # 已废弃


class MarketplaceAdapter(Base):
    """市场适配器"""
    __tablename__ = "marketplace_adapters"
    
    # 基本信息
    id = Column(Integer, primary_key=True)
    adapter_id = Column(String(100), unique=True, nullable=False)  # 唯一标识
    name = Column(String(100), nullable=False)
    display_name = Column(String(200), nullable=False)
    description = Column(Text, nullable=False)
    tagline = Column(String(200))  # 一句话描述
    
    # 作者信息
    author_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    author = relationship("User", back_populates="adapters")
    organization = Column(String(200))
    
    # 版本信息
    current_version = Column(String(20), nullable=False)
    latest_version_id = Column(Integer, ForeignKey("adapter_versions.id"))
    
    # 分类和标签
    category = Column(String(50), nullable=False)  # 主分类：k12/higher_ed/enterprise
    subcategory = Column(String(50))  # 子分类：chemistry/physics/math
    tags = Column(JSON)  # 标签列表
    
    # 能力声明
    supported_domains = Column(JSON)  # 支持的领域
    supported_activity_types = Column(JSON)  # 支持的活动类型
    supported_output_formats = Column(JSON)  # 支持的输出格式
    
    # 依赖
    platform_version_min = Column(String(20))  # 最低平台版本
    platform_version_max = Column(String(20))  # 最高平台版本
    dependencies = Column(JSON)  # 依赖的其他适配器
    
    # 发布状态
    status = Column(Enum(AdapterPublishStatus), default=AdapterPublishStatus.DRAFT)
    published_at = Column(DateTime)
    deprecated_at = Column(DateTime)
    
    # 统计数据
    download_count = Column(Integer, default=0)
    install_count = Column(Integer, default=0)
    active_install_count = Column(Integer, default=0)  # 活跃安装数
    rating_average = Column(Float, default=0.0)
    rating_count = Column(Integer, default=0)
    
    # 质量指标
    quality_score = Column(Float)  # 综合质量分
    test_coverage = Column(Float)  # 测试覆盖率
    security_scan_passed = Column(Boolean, default=False)
    performance_score = Column(Float)  # 性能评分
    
    # 资源链接
    homepage_url = Column(String(500))
    documentation_url = Column(String(500))
    repository_url = Column(String(500))
    support_url = Column(String(500))
    
    # 媒体资源
    icon_url = Column(String(500))
    screenshots = Column(JSON)  # 截图URLs
    demo_video_url = Column(String(500))
    
    # 许可证
    license = Column(String(50))  # MIT/Apache-2.0/GPL-3.0等
    is_open_source = Column(Boolean, default=True)
    
    # 定价（为未来商业化预留）
    is_free = Column(Boolean, default=True)
    pricing_model = Column(String(50))  # free/freemium/paid
    price = Column(Float)  # 月费
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # 关系
    versions = relationship("AdapterVersion", back_populates="adapter", cascade="all, delete-orphan")
    reviews = relationship("AdapterReview", back_populates="adapter", cascade="all, delete-orphan")


class AdapterVersion(Base):
    """适配器版本"""
    __tablename__ = "adapter_versions"
    
    id = Column(Integer, primary_key=True)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    
    # 版本信息
    version = Column(String(20), nullable=False)
    version_name = Column(String(100))  # 版本代号（可选）
    
    # 变更信息
    changelog = Column(Text)
    release_notes = Column(Text)
    
    # 文件信息
    package_url = Column(String(500), nullable=False)  # 包文件URL
    package_size = Column(Integer)  # 字节
    package_checksum = Column(String(64))  # SHA-256
    
    # manifest.yaml内容
    manifest = Column(JSON, nullable=False)
    
    # 兼容性
    platform_version_min = Column(String(20))
    platform_version_max = Column(String(20))
    breaking_changes = Column(Boolean, default=False)
    
    # 审核信息
    review_status = Column(Enum(AdapterPublishStatus), default=AdapterPublishStatus.DRAFT)
    reviewed_by = Column(Integer, ForeignKey("users.id"))
    reviewed_at = Column(DateTime)
    review_comments = Column(Text)
    
    # 质量检测
    test_results = Column(JSON)  # 测试结果
    security_scan_results = Column(JSON)  # 安全扫描结果
    performance_benchmarks = Column(JSON)  # 性能基准测试
    
    # 统计
    download_count = Column(Integer, default=0)
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    published_at = Column(DateTime)
    
    # 关系
    adapter = relationship("MarketplaceAdapter", back_populates="versions")


class AdapterReview(Base):
    """适配器评论"""
    __tablename__ = "adapter_reviews"
    
    id = Column(Integer, primary_key=True)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    
    # 评分（1-5星）
    rating = Column(Integer, nullable=False)
    
    # 评论内容
    title = Column(String(200))
    content = Column(Text)
    
    # 使用场景
    use_case = Column(String(100))  # 使用场景分类
    adapter_version = Column(String(20))  # 评论时使用的版本
    
    # 点赞统计
    helpful_count = Column(Integer, default=0)
    not_helpful_count = Column(Integer, default=0)
    
    # 审核
    is_verified_user = Column(Boolean, default=False)  # 是否认证用户
    is_moderated = Column(Boolean, default=False)  # 是否经过审核
    is_visible = Column(Boolean, default=True)
    
    # 作者回复
    author_reply = Column(Text)
    author_replied_at = Column(DateTime)
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # 关系
    adapter = relationship("MarketplaceAdapter", back_populates="reviews")
    user = relationship("User", back_populates="reviews")


class AdapterInstallation(Base):
    """适配器安装记录"""
    __tablename__ = "adapter_installations"
    
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    version_id = Column(Integer, ForeignKey("adapter_versions.id"), nullable=False)
    
    # 安装状态
    status = Column(String(20))  # installed/active/inactive/uninstalled
    
    # 配置
    configuration = Column(JSON)  # 用户自定义配置
    
    # 使用统计
    last_used_at = Column(DateTime)
    usage_count = Column(Integer, default=0)
    
    # 时间戳
    installed_at = Column(DateTime, default=datetime.utcnow)
    uninstalled_at = Column(DateTime)


class User(Base):
    """用户"""
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True)
    username = Column(String(50), unique=True, nullable=False)
    email = Column(String(200), unique=True, nullable=False)
    
    # 身份
    is_developer = Column(Boolean, default=False)
    is_verified = Column(Boolean, default=False)
    is_admin = Column(Boolean, default=False)
    
    # 个人信息
    display_name = Column(String(100))
    avatar_url = Column(String(500))
    bio = Column(Text)
    organization = Column(String(200))
    website = Column(String(500))
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    
    # 关系
    adapters = relationship("MarketplaceAdapter", back_populates="author")
    reviews = relationship("AdapterReview", back_populates="user")
```

---

## 3.2 发布流程

### 3.2.1 开发者发布流程

```
┌──────────────┐
│ 1. 创建适配器 │
│   开发本地代码  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 2. 本地测试   │
│  运行单元测试  │
│  质量检查     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 3. 打包       │
│  生成manifest │
│  创建.zip包   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 4. 提交审核   │
│  上传到市场   │
│  填写描述信息  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 5. 自动检测   │
│  安全扫描     │
│  性能测试     │
│  兼容性测试   │
└──────┬───────┘
       │
       ├─→ 失败 ──→ 通知开发者 ──→ 返回步骤1
       │
       ▼ 通过
┌──────────────┐
│ 6. 人工审核   │
│  代码审查     │
│  功能验证     │
│  文档检查     │
└──────┬───────┘
       │
       ├─→ 拒绝 ──→ 通知原因 ──→ 返回步骤1
       │
       ▼ 批准
┌──────────────┐
│ 7. 发布上架   │
│  设置为可见   │
│  发送通知     │
└──────────────┘
```

### 3.2.2 发布API设计

```python
# src/marketplace/api/publish.py

from fastapi import APIRouter, HTTPException, Depends, UploadFile, File
from typing import Optional
from pydantic import BaseModel

router = APIRouter(prefix="/api/v1/marketplace/publish", tags=["publish"])


class AdapterSubmitRequest(BaseModel):
    """适配器提交请求"""
    name: str
    display_name: str
    description: str
    tagline: Optional[str]
    category: str
    subcategory: Optional[str]
    tags: List[str] = []
    homepage_url: Optional[str]
    repository_url: Optional[str]
    documentation_url: Optional[str]
    license: str = "MIT"


class VersionSubmitRequest(BaseModel):
    """版本提交请求"""
    version: str
    changelog: str
    release_notes: Optional[str]
    breaking_changes: bool = False


@router.post("/adapters")
async def create_adapter(
    request: AdapterSubmitRequest,
    current_user: User = Depends(get_current_user)
):
    """
    创建新适配器（草稿）
    
    这个API创建一个草稿状态的适配器记录，
    开发者需要上传版本包后才能提交审核。
    """
    # 验证开发者权限
    if not current_user.is_developer:
        raise HTTPException(
            status_code=403,
            detail="Only verified developers can publish adapters"
        )
    
    # 检查名称是否已存在
    existing = await db.query(MarketplaceAdapter).filter_by(
        name=request.name
    ).first()
    
    if existing:
        raise HTTPException(
            status_code=400,
            detail=f"Adapter name '{request.name}' already exists"
        )
    
    # 创建适配器记录
    adapter = MarketplaceAdapter(
        adapter_id=generate_adapter_id(request.name),
        name=request.name,
        display_name=request.display_name,
        description=request.description,
        tagline=request.tagline,
        category=request.category,
        subcategory=request.subcategory,
        tags=request.tags,
        author_id=current_user.id,
        homepage_url=request.homepage_url,
        repository_url=request.repository_url,
        documentation_url=request.documentation_url,
        license=request.license,
        status=AdapterPublishStatus.DRAFT
    )
    
    db.add(adapter)
    await db.commit()
    
    return {
        "adapter_id": adapter.adapter_id,
        "status": "draft",
        "message": "Adapter created. Please upload a version package."
    }


@router.post("/adapters/{adapter_id}/versions")
async def upload_version(
    adapter_id: str,
    version_info: VersionSubmitRequest,
    package_file: UploadFile = File(...),
    current_user: User = Depends(get_current_user)
):
    """
    上传适配器版本包
    
    步骤：
    1. 验证权限
    2. 验证包格式
    3. 提取manifest.yaml
    4. 上传到对象存储
    5. 创建版本记录
    """
    # 查找适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 验证所有权
    if adapter.author_id != current_user.id:
        raise HTTPException(
            status_code=403,
            detail="You don't have permission to modify this adapter"
        )
    
    # 验证包格式
    if not package_file.filename.endswith('.zip'):
        raise HTTPException(
            status_code=400,
            detail="Package must be a .zip file"
        )
    
    # 读取包内容
    package_content = await package_file.read()
    
    # 验证和提取manifest
    try:
        manifest = await extract_manifest(package_content)
    except Exception as e:
        raise HTTPException(
            status_code=400,
            detail=f"Invalid package: {str(e)}"
        )
    
    # 验证版本号匹配
    if manifest['version'] != version_info.version:
        raise HTTPException(
            status_code=400,
            detail="Version mismatch between request and manifest"
        )
    
    # 检查版本是否已存在
    existing_version = await db.query(AdapterVersion).filter_by(
        adapter_id=adapter.id,
        version=version_info.version
    ).first()
    
    if existing_version:
        raise HTTPException(
            status_code=400,
            detail=f"Version {version_info.version} already exists"
        )
    
    # 上传到对象存储
    package_url = await upload_to_storage(
        package_content,
        f"adapters/{adapter_id}/{version_info.version}/package.zip"
    )
    
    # 计算校验和
    checksum = calculate_sha256(package_content)
    
    # 创建版本记录
    version = AdapterVersion(
        adapter_id=adapter.id,
        version=version_info.version,
        changelog=version_info.changelog,
        release_notes=version_info.release_notes,
        breaking_changes=version_info.breaking_changes,
        package_url=package_url,
        package_size=len(package_content),
        package_checksum=checksum,
        manifest=manifest,
        platform_version_min=manifest.get('compatibility', {}).get('platform_version'),
        review_status=AdapterPublishStatus.DRAFT
    )
    
    db.add(version)
    
    # 更新适配器的当前版本
    if not adapter.current_version or version_info.version > adapter.current_version:
        adapter.current_version = version_info.version
        adapter.latest_version_id = version.id
    
    await db.commit()
    
    return {
        "version_id": version.id,
        "version": version.version,
        "package_url": package_url,
        "status": "draft",
        "message": "Version uploaded successfully. You can now submit for review."
    }


@router.post("/adapters/{adapter_id}/submit")
async def submit_for_review(
    adapter_id: str,
    current_user: User = Depends(get_current_user)
):
    """
    提交适配器审核
    
    将适配器状态从DRAFT改为SUBMITTED，
    并触发自动检测流程。
    """
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id:
        raise HTTPException(status_code=403, detail="Permission denied")
    
    # 检查是否有版本
    if not adapter.versions:
        raise HTTPException(
            status_code=400,
            detail="Cannot submit adapter without any versions"
        )
    
    # 检查当前状态
    if adapter.status not in [AdapterPublishStatus.DRAFT, AdapterPublishStatus.REJECTED]:
        raise HTTPException(
            status_code=400,
            detail=f"Cannot submit adapter with status: {adapter.status.value}"
        )
    
    # 更新状态
    adapter.status = AdapterPublishStatus.SUBMITTED
    
    # 获取最新版本
    latest_version = adapter.versions[-1]
    latest_version.review_status = AdapterPublishStatus.SUBMITTED
    
    await db.commit()
    
    # 触发自动检测（异步任务）
    from src.marketplace.tasks import run_automated_checks
    await run_automated_checks.delay(adapter.id, latest_version.id)
    
    # 通知审核团队
    await notify_review_team(adapter.id)
    
    return {
        "adapter_id": adapter_id,
        "status": "submitted",
        "message": "Adapter submitted for review. Automated checks will run shortly."
    }
```

---

## 3.3 自动化检测

### 3.3.1 检测流程

```python
# src/marketplace/tasks.py

from celery import Celery
import asyncio

celery_app = Celery('marketplace', broker='redis://localhost:6379/0')


@celery_app.task
async def run_automated_checks(adapter_id: int, version_id: int):
    """
    运行自动化检测
    
    检测项：
    1. 安全扫描
    2. 性能测试
    3. 兼容性测试
    4. 代码质量检测
    """
    version = await db.query(AdapterVersion).get(version_id)
    
    results = {
        "security": None,
        "performance": None,
        "compatibility": None,
        "quality": None,
        "overall_passed": False
    }
    
    try:
        # 1. 安全扫描
        results["security"] = await run_security_scan(version)
        
        # 2. 性能测试
        results["performance"] = await run_performance_tests(version)
        
        # 3. 兼容性测试
        results["compatibility"] = await run_compatibility_tests(version)
        
        # 4. 代码质量
        results["quality"] = await run_quality_checks(version)
        
        # 判断是否全部通过
        results["overall_passed"] = all([
            results["security"]["passed"],
            results["performance"]["passed"],
            results["compatibility"]["passed"],
            results["quality"]["passed"]
        ])
        
        # 保存结果
        version.test_results = results["quality"]
        version.security_scan_results = results["security"]
        version.performance_benchmarks = results["performance"]
        
        # 更新状态
        if results["overall_passed"]:
            version.review_status = AdapterPublishStatus.REVIEWING
            # 通知人工审核
            await notify_human_reviewers(adapter_id, version_id)
        else:
            version.review_status = AdapterPublishStatus.REJECTED
            # 通知开发者
            await notify_developer_rejection(adapter_id, version_id, results)
        
        await db.commit()
        
    except Exception as e:
        logger.error(f"Automated checks failed: {e}")
        version.review_status = AdapterPublishStatus.REJECTED
        await db.commit()
    
    return results


async def run_security_scan(version: AdapterVersion) -> Dict:
    """
    安全扫描
    
    检测：
    - 恶意代码
    - 已知漏洞
    - 敏感信息泄露
    - 不安全的依赖
    """
    # 下载包
    package_path = await download_package(version.package_url)
    
    results = {
        "passed": True,
        "issues": [],
        "severity_counts": {"critical": 0, "high": 0, "medium": 0, "low": 0}
    }
    
    # 1. 病毒扫描
    virus_scan = await run_clamav_scan(package_path)
    if not virus_scan["clean"]:
        results["passed"] = False
        results["issues"].append({
            "type": "malware",
            "severity": "critical",
            "description": "Malware detected"
        })
    
    # 2. 依赖漏洞扫描（使用safety或snyk）
    dependency_scan = await run_safety_check(package_path)
    for vuln in dependency_scan["vulnerabilities"]:
        results["issues"].append({
            "type": "vulnerability",
            "severity": vuln["severity"],
            "package": vuln["package"],
            "description": vuln["description"],
            "cve": vuln.get("cve")
        })
        results["severity_counts"][vuln["severity"]] += 1
    
    # 3. 代码静态分析（使用bandit）
    code_scan = await run_bandit_scan(package_path)
    for issue in code_scan["issues"]:
        if issue["severity"] in ["high", "critical"]:
            results["passed"] = False
        results["issues"].append({
            "type": "code_issue",
            "severity": issue["severity"],
            "description": issue["description"],
            "file": issue["file"],
            "line": issue["line"]
        })
    
    # 4. 检查敏感信息
    secrets_scan = await run_trufflehog_scan(package_path)
    if secrets_scan["found"]:
        results["passed"] = False
        results["issues"].append({
            "type": "secrets",
            "severity": "critical",
            "description": "Hardcoded secrets detected"
        })
    
    # Critical或High级别的问题会导致不通过
    if results["severity_counts"]["critical"] > 0 or results["severity_counts"]["high"] > 0:
        results["passed"] = False
    
    return results


async def run_performance_tests(version: AdapterVersion) -> Dict:
    """
    性能测试
    
    测试：
    - 内存占用
    - CPU使用率
    - 响应时间
    - 吞吐量
    """
    results = {
        "passed": True,
        "metrics": {},
        "benchmarks": []
    }
    
    # 加载适配器
    adapter = await load_adapter_for_testing(version)
    
    # 1. 内存测试
    memory_usage = await measure_memory_usage(adapter)
    results["metrics"]["memory_mb"] = memory_usage
    
    # 内存超过500MB视为不通过
    if memory_usage > 500:
        results["passed"] = False
    
    # 2. 响应时间测试
    response_times = []
    for i in range(10):
        start = time.time()
        await adapter.parse_request({
            "activity_type": "concept_understanding",
            "subject": "chemistry",
            "topic": "periodic_table"
        })
        response_times.append(time.time() - start)
    
    avg_response_time = sum(response_times) / len(response_times)
    results["metrics"]["avg_response_time_ms"] = avg_response_time * 1000
    
    # 平均响应时间超过5秒视为不通过
    if avg_response_time > 5.0:
        results["passed"] = False
    
    # 3. 并发测试
    concurrent_results = await run_concurrent_test(adapter, concurrency=10)
    results["metrics"]["concurrent_throughput"] = concurrent_results["throughput"]
    
    return results
```

---

*文档第四部分完成*

**已完成内容**：
- ✅ 市场整体架构（4层架构）
- ✅ 完整的数据模型（适配器、版本、评论、安装、用户）
- ✅ 发布流程设计（7步审核流程）
- ✅ 发布API实现（创建、上传、提交）
- ✅ 自动化检测（安全扫描、性能测试）

**下一部分将包含**：
- 搜索和推荐系统
- 评分评论机制
- 使用统计分析
- 开发者工具和SDK

是否继续输出第五部分（用户体验和分析系统）？# MetaWorkflow V2.0 适配器市场设计 - 第五部分

## 第四部分：搜索、推荐与用户体验

## 4.1 搜索系统

### 4.1.1 Elasticsearch索引设计

```python
# src/marketplace/search/indexing.py

from elasticsearch import Elasticsearch
from typing import Dict, List, Any

# Elasticsearch索引映射
ADAPTER_INDEX_MAPPING = {
    "mappings": {
        "properties": {
            # 基本信息
            "adapter_id": {"type": "keyword"},
            "name": {"type": "keyword"},
            "display_name": {
                "type": "text",
                "fields": {
                    "keyword": {"type": "keyword"},
                    "suggest": {"type": "completion"}
                },
                "analyzer": "standard"
            },
            "description": {
                "type": "text",
                "analyzer": "standard"
            },
            "tagline": {
                "type": "text",
                "analyzer": "standard"
            },
            
            # 分类和标签
            "category": {"type": "keyword"},
            "subcategory": {"type": "keyword"},
            "tags": {"type": "keyword"},
            
            # 能力
            "supported_domains": {"type": "keyword"},
            "supported_activity_types": {"type": "keyword"},
            "supported_output_formats": {"type": "keyword"},
            
            # 作者
            "author_name": {
                "type": "text",
                "fields": {"keyword": {"type": "keyword"}}
            },
            "organization": {"type": "keyword"},
            
            # 质量指标（用于排序）
            "quality_score": {"type": "float"},
            "rating_average": {"type": "float"},
            "rating_count": {"type": "integer"},
            "download_count": {"type": "integer"},
            "install_count": {"type": "integer"},
            "active_install_count": {"type": "integer"},
            
            # 状态
            "status": {"type": "keyword"},
            "is_free": {"type": "boolean"},
            "is_open_source": {"type": "boolean"},
            
            # 时间
            "published_at": {"type": "date"},
            "updated_at": {"type": "date"},
            
            # 版本信息
            "current_version": {"type": "keyword"},
            "platform_version_min": {"type": "keyword"},
            
            # 全文搜索字段（组合多个字段）
            "all_text": {
                "type": "text",
                "analyzer": "standard"
            }
        }
    },
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 2,
        "analysis": {
            "analyzer": {
                "standard": {
                    "type": "standard",
                    "stopwords": "_english_"
                }
            }
        }
    }
}


class AdapterSearchService:
    """适配器搜索服务"""
    
    def __init__(self, es_client: Elasticsearch):
        self.es = es_client
        self.index_name = "marketplace_adapters"
    
    async def index_adapter(self, adapter: MarketplaceAdapter) -> None:
        """
        索引适配器
        
        Args:
            adapter: 适配器对象
        """
        doc = {
            "adapter_id": adapter.adapter_id,
            "name": adapter.name,
            "display_name": adapter.display_name,
            "description": adapter.description,
            "tagline": adapter.tagline,
            "category": adapter.category,
            "subcategory": adapter.subcategory,
            "tags": adapter.tags or [],
            "supported_domains": adapter.supported_domains or [],
            "supported_activity_types": adapter.supported_activity_types or [],
            "supported_output_formats": adapter.supported_output_formats or [],
            "author_name": adapter.author.display_name or adapter.author.username,
            "organization": adapter.organization,
            "quality_score": adapter.quality_score or 0.0,
            "rating_average": adapter.rating_average,
            "rating_count": adapter.rating_count,
            "download_count": adapter.download_count,
            "install_count": adapter.install_count,
            "active_install_count": adapter.active_install_count,
            "status": adapter.status.value,
            "is_free": adapter.is_free,
            "is_open_source": adapter.is_open_source,
            "published_at": adapter.published_at.isoformat() if adapter.published_at else None,
            "updated_at": adapter.updated_at.isoformat(),
            "current_version": adapter.current_version,
            "platform_version_min": adapter.platform_version_min,
            # 组合搜索字段
            "all_text": f"{adapter.display_name} {adapter.description} {adapter.tagline} {' '.join(adapter.tags or [])}"
        }
        
        await self.es.index(
            index=self.index_name,
            id=adapter.adapter_id,
            document=doc
        )
    
    async def search(
        self,
        query: str = None,
        category: str = None,
        tags: List[str] = None,
        domains: List[str] = None,
        activity_types: List[str] = None,
        min_rating: float = None,
        is_free: bool = None,
        sort_by: str = "relevance",
        page: int = 1,
        page_size: int = 20
    ) -> Dict[str, Any]:
        """
        搜索适配器
        
        Args:
            query: 搜索关键词
            category: 分类过滤
            tags: 标签过滤
            domains: 领域过滤
            activity_types: 活动类型过滤
            min_rating: 最低评分
            is_free: 是否免费
            sort_by: 排序方式 (relevance/rating/downloads/newest)
            page: 页码
            page_size: 每页数量
            
        Returns:
            Dict: 搜索结果
        """
        # 构建查询
        must_conditions = []
        filter_conditions = [
            {"term": {"status": "published"}}  # 只搜索已发布的
        ]
        
        # 关键词搜索
        if query:
            must_conditions.append({
                "multi_match": {
                    "query": query,
                    "fields": [
                        "display_name^3",      # display_name权重最高
                        "tagline^2",           # tagline次之
                        "description",
                        "tags",
                        "all_text"
                    ],
                    "type": "best_fields",
                    "fuzziness": "AUTO"        # 允许模糊匹配
                }
            })
        
        # 分类过滤
        if category:
            filter_conditions.append({"term": {"category": category}})
        
        # 标签过滤
        if tags:
            filter_conditions.append({"terms": {"tags": tags}})
        
        # 领域过滤
        if domains:
            filter_conditions.append({"terms": {"supported_domains": domains}})
        
        # 活动类型过滤
        if activity_types:
            filter_conditions.append({"terms": {"supported_activity_types": activity_types}})
        
        # 评分过滤
        if min_rating is not None:
            filter_conditions.append({
                "range": {"rating_average": {"gte": min_rating}}
            })
        
        # 免费过滤
        if is_free is not None:
            filter_conditions.append({"term": {"is_free": is_free}})
        
        # 构建完整查询
        search_query = {
            "bool": {
                "must": must_conditions if must_conditions else [{"match_all": {}}],
                "filter": filter_conditions
            }
        }
        
        # 排序
        sort_options = self._get_sort_options(sort_by)
        
        # 执行搜索
        from_index = (page - 1) * page_size
        
        response = await self.es.search(
            index=self.index_name,
            query=search_query,
            sort=sort_options,
            from_=from_index,
            size=page_size,
            track_total_hits=True
        )
        
        # 处理结果
        hits = response['hits']
        total = hits['total']['value']
        
        adapters = []
        for hit in hits['hits']:
            adapter_data = hit['_source']
            adapter_data['_score'] = hit['_score']
            adapters.append(adapter_data)
        
        return {
            "total": total,
            "page": page,
            "page_size": page_size,
            "total_pages": (total + page_size - 1) // page_size,
            "adapters": adapters
        }
    
    def _get_sort_options(self, sort_by: str) -> List[Dict]:
        """获取排序选项"""
        sort_map = {
            "relevance": [
                {"_score": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"download_count": {"order": "desc"}}
            ],
            "rating": [
                {"rating_average": {"order": "desc"}},
                {"rating_count": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "downloads": [
                {"download_count": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "newest": [
                {"published_at": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "popularity": [
                {"active_install_count": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ]
        }
        
        return sort_map.get(sort_by, sort_map["relevance"])
    
    async def suggest(self, prefix: str, size: int = 10) -> List[str]:
        """
        自动补全建议
        
        Args:
            prefix: 输入前缀
            size: 返回数量
            
        Returns:
            List[str]: 建议列表
        """
        response = await self.es.search(
            index=self.index_name,
            suggest={
                "adapter_suggest": {
                    "prefix": prefix,
                    "completion": {
                        "field": "display_name.suggest",
                        "size": size,
                        "skip_duplicates": True
                    }
                }
            }
        )
        
        suggestions = []
        for option in response['suggest']['adapter_suggest'][0]['options']:
            suggestions.append(option['text'])
        
        return suggestions


# 全局搜索服务实例
search_service = AdapterSearchService(es_client)
```

### 4.1.2 搜索API

```python
# src/marketplace/api/search.py

from fastapi import APIRouter, Query
from typing import Optional, List

router = APIRouter(prefix="/api/v1/marketplace/search", tags=["search"])


@router.get("/adapters")
async def search_adapters(
    q: Optional[str] = Query(None, description="搜索关键词"),
    category: Optional[str] = Query(None, description="分类"),
    tags: Optional[List[str]] = Query(None, description="标签"),
    domains: Optional[List[str]] = Query(None, description="支持的领域"),
    activity_types: Optional[List[str]] = Query(None, description="活动类型"),
    min_rating: Optional[float] = Query(None, ge=0, le=5, description="最低评分"),
    is_free: Optional[bool] = Query(None, description="是否免费"),
    sort_by: str = Query("relevance", description="排序方式"),
    page: int = Query(1, ge=1, description="页码"),
    page_size: int = Query(20, ge=1, le=100, description="每页数量")
):
    """
    搜索适配器
    
    支持的排序方式：
    - relevance: 相关性（默认）
    - rating: 评分
    - downloads: 下载量
    - newest: 最新
    - popularity: 活跃安装量
    """
    from src.marketplace.search.indexing import search_service
    
    results = await search_service.search(
        query=q,
        category=category,
        tags=tags,
        domains=domains,
        activity_types=activity_types,
        min_rating=min_rating,
        is_free=is_free,
        sort_by=sort_by,
        page=page,
        page_size=page_size
    )
    
    return results


@router.get("/suggest")
async def get_suggestions(
    q: str = Query(..., min_length=1, description="输入前缀"),
    size: int = Query(10, ge=1, le=20, description="建议数量")
):
    """获取搜索自动补全建议"""
    from src.marketplace.search.indexing import search_service
    
    suggestions = await search_service.suggest(q, size)
    
    return {"suggestions": suggestions}


@router.get("/filters")
async def get_available_filters():
    """
    获取可用的过滤选项
    
    返回所有可用的分类、标签、领域等选项，
    用于前端构建过滤器UI。
    """
    # 从数据库聚合统计
    categories = await db.query(
        MarketplaceAdapter.category,
        func.count(MarketplaceAdapter.id).label('count')
    ).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.category
    ).all()
    
    # 统计所有标签
    all_tags = await db.query(MarketplaceAdapter.tags).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).all()
    
    tag_counts = {}
    for (tags,) in all_tags:
        if tags:
            for tag in tags:
                tag_counts[tag] = tag_counts.get(tag, 0) + 1
    
    return {
        "categories": [
            {"name": cat, "count": count}
            for cat, count in categories
        ],
        "tags": [
            {"name": tag, "count": count}
            for tag, count in sorted(tag_counts.items(), key=lambda x: x[1], reverse=True)
        ],
        "sort_options": [
            {"value": "relevance", "label": "相关性"},
            {"value": "rating", "label": "评分"},
            {"value": "downloads", "label": "下载量"},
            {"value": "newest", "label": "最新"},
            {"value": "popularity", "label": "热度"}
        ]
    }
```

---

## 4.2 推荐系统

### 4.2.1 推荐算法

```python
# src/marketplace/recommendations.py

from typing import List, Dict, Any
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

class AdapterRecommendationEngine:
    """适配器推荐引擎"""
    
    def __init__(self):
        self.adapter_features = {}  # 适配器特征向量缓存
    
    async def get_recommendations(
        self,
        user_id: int,
        context: Dict[str, Any] = None,
        count: int = 10
    ) -> List[Dict[str, Any]]:
        """
        获取个性化推荐
        
        推荐策略：
        1. 基于用户历史安装（协同过滤）
        2. 基于内容相似度
        3. 热门推荐（新用户冷启动）
        4. 上下文推荐（当前正在做什么）
        
        Args:
            user_id: 用户ID
            context: 上下文信息（当前领域、活动类型等）
            count: 推荐数量
            
        Returns:
            List[Dict]: 推荐的适配器列表
        """
        recommendations = []
        
        # 1. 获取用户安装历史
        user_installations = await self._get_user_installations(user_id)
        
        if user_installations:
            # 有安装历史：协同过滤 + 内容相似
            collaborative_recs = await self._collaborative_filtering(user_id, count // 2)
            content_recs = await self._content_based_filtering(user_installations, count // 2)
            
            recommendations.extend(collaborative_recs)
            recommendations.extend(content_recs)
        else:
            # 新用户：热门推荐
            recommendations = await self._popular_recommendations(count)
        
        # 2. 如果有上下文，添加上下文相关推荐
        if context:
            context_recs = await self._context_based_recommendations(context, count // 3)
            recommendations.extend(context_recs)
        
        # 3. 去重和排序
        recommendations = self._deduplicate_and_rank(recommendations, count)
        
        return recommendations
    
    async def _collaborative_filtering(
        self,
        user_id: int,
        count: int
    ) -> List[Dict]:
        """
        协同过滤推荐
        
        逻辑：找到和用户安装了相同适配器的其他用户，
        推荐他们还安装了什么。
        """
        # 获取用户安装的适配器
        user_adapters = await db.query(
            AdapterInstallation.adapter_id
        ).filter(
            AdapterInstallation.user_id == user_id,
            AdapterInstallation.status == 'active'
        ).all()
        
        user_adapter_ids = [a[0] for a in user_adapters]
        
        # 找到安装了相同适配器的其他用户
        similar_users = await db.query(
            AdapterInstallation.user_id,
            func.count(AdapterInstallation.adapter_id).label('common_count')
        ).filter(
            AdapterInstallation.adapter_id.in_(user_adapter_ids),
            AdapterInstallation.user_id != user_id,
            AdapterInstallation.status == 'active'
        ).group_by(
            AdapterInstallation.user_id
        ).order_by(
            desc('common_count')
        ).limit(50).all()
        
        similar_user_ids = [u[0] for u in similar_users]
        
        # 找到这些用户还安装了什么
        recommendations = await db.query(
            MarketplaceAdapter,
            func.count(AdapterInstallation.user_id).label('install_count')
        ).join(
            AdapterInstallation
        ).filter(
            AdapterInstallation.user_id.in_(similar_user_ids),
            AdapterInstallation.adapter_id.notin_(user_adapter_ids),
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
        ).group_by(
            MarketplaceAdapter.id
        ).order_by(
            desc('install_count')
        ).limit(count).all()
        
        return [
            {
                "adapter": adapter,
                "score": install_count / len(similar_user_ids),
                "reason": "similar_users"
            }
            for adapter, install_count in recommendations
        ]
    
    async def _content_based_filtering(
        self,
        user_installations: List[MarketplaceAdapter],
        count: int
    ) -> List[Dict]:
        """
        基于内容的推荐
        
        逻辑：找到和用户已安装适配器相似的其他适配器。
        相似度基于：分类、标签、支持的领域等。
        """
        # 构建用户偏好特征向量
        user_features = self._build_user_feature_vector(user_installations)
        
        # 获取所有候选适配器
        candidates = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED,
            MarketplaceAdapter.id.notin_([a.id for a in user_installations])
        ).all()
        
        # 计算相似度
        similarities = []
        for candidate in candidates:
            candidate_features = self._build_adapter_feature_vector(candidate)
            similarity = cosine_similarity(
                user_features.reshape(1, -1),
                candidate_features.reshape(1, -1)
            )[0][0]
            
            similarities.append({
                "adapter": candidate,
                "score": similarity,
                "reason": "content_similarity"
            })
        
        # 按相似度排序
        similarities.sort(key=lambda x: x['score'], reverse=True)
        
        return similarities[:count]
    
    async def _popular_recommendations(self, count: int) -> List[Dict]:
        """
        热门推荐
        
        基于：活跃安装量、评分、下载量的加权组合
        """
        adapters = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
        ).all()
        
        # 计算热度分数
        scored_adapters = []
        for adapter in adapters:
            popularity_score = (
                adapter.active_install_count * 0.5 +
                adapter.rating_average * adapter.rating_count * 0.3 +
                adapter.download_count * 0.0001 * 0.2
            )
            
            scored_adapters.append({
                "adapter": adapter,
                "score": popularity_score,
                "reason": "popular"
            })
        
        # 排序
        scored_adapters.sort(key=lambda x: x['score'], reverse=True)
        
        return scored_adapters[:count]
    
    async def _context_based_recommendations(
        self,
        context: Dict[str, Any],
        count: int
    ) -> List[Dict]:
        """
        基于上下文的推荐
        
        Args:
            context: {
                "domain": "chemistry",
                "activity_type": "concept_understanding",
                "subject": "periodic_table"
            }
        """
        filters = []
        
        if "domain" in context:
            filters.append(
                MarketplaceAdapter.supported_domains.contains([context["domain"]])
            )
        
        if "activity_type" in context:
            filters.append(
                MarketplaceAdapter.supported_activity_types.contains([context["activity_type"]])
            )
        
        adapters = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED,
            *filters
        ).order_by(
            desc(MarketplaceAdapter.rating_average)
        ).limit(count).all()
        
        return [
            {
                "adapter": adapter,
                "score": adapter.rating_average,
                "reason": "context_match"
            }
            for adapter in adapters
        ]
    
    def _build_user_feature_vector(
        self,
        installations: List[MarketplaceAdapter]
    ) -> np.ndarray:
        """构建用户特征向量（基于已安装适配器的平均）"""
        feature_vectors = [
            self._build_adapter_feature_vector(adapter)
            for adapter in installations
        ]
        
        return np.mean(feature_vectors, axis=0)
    
    def _build_adapter_feature_vector(
        self,
        adapter: MarketplaceAdapter
    ) -> np.ndarray:
        """
        构建适配器特征向量
        
        特征包括：
        - 分类 (one-hot)
        - 标签 (multi-hot)
        - 领域 (multi-hot)
        - 活动类型 (multi-hot)
        """
        # 这里简化示意，实际应该用更完整的特征工程
        features = []
        
        # 分类 (假设有10个分类)
        category_features = [0] * 10
        category_idx = hash(adapter.category) % 10
        category_features[category_idx] = 1
        features.extend(category_features)
        
        # 标签 (假设有50个常见标签)
        tag_features = [0] * 50
        if adapter.tags:
            for tag in adapter.tags:
                tag_idx = hash(tag) % 50
                tag_features[tag_idx] = 1
        features.extend(tag_features)
        
        # 质量分数
        features.append(adapter.quality_score or 0.0)
        features.append(adapter.rating_average)
        
        return np.array(features)
    
    def _deduplicate_and_rank(
        self,
        recommendations: List[Dict],
        count: int
    ) -> List[Dict]:
        """去重并重新排序"""
        seen = set()
        unique_recs = []
        
        for rec in recommendations:
            adapter_id = rec['adapter'].adapter_id
            if adapter_id not in seen:
                seen.add(adapter_id)
                unique_recs.append(rec)
        
        # 按分数排序
        unique_recs.sort(key=lambda x: x['score'], reverse=True)
        
        return unique_recs[:count]


# 全局推荐引擎实例
recommendation_engine = AdapterRecommendationEngine()
```

### 4.2.2 推荐API

```python
# src/marketplace/api/recommendations.py

from fastapi import APIRouter, Depends

router = APIRouter(prefix="/api/v1/marketplace/recommendations", tags=["recommendations"])


@router.get("/for-you")
async def get_personalized_recommendations(
    count: int = Query(10, ge=1, le=50),
    current_user: User = Depends(get_current_user)
):
    """
    获取个性化推荐
    
    基于用户的安装历史、评分、使用习惯等
    """
    from src.marketplace.recommendations import recommendation_engine
    
    recommendations = await recommendation_engine.get_recommendations(
        user_id=current_user.id,
        count=count
    )
    
    return {
        "recommendations": [
            {
                "adapter_id": rec['adapter'].adapter_id,
                "name": rec['adapter'].display_name,
                "description": rec['adapter'].tagline,
                "rating": rec['adapter'].rating_average,
                "downloads": rec['adapter'].download_count,
                "icon_url": rec['adapter'].icon_url,
                "reason": rec['reason'],
                "score": rec['score']
            }
            for rec in recommendations
        ]
    }


@router.get("/similar/{adapter_id}")
async def get_similar_adapters(
    adapter_id: str,
    count: int = Query(5, ge=1, le=20)
):
    """
    获取相似的适配器
    
    用于"您可能还喜欢"功能
    """
    # 获取目标适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    from src.marketplace.recommendations import recommendation_engine
    
    # 基于内容相似度推荐
    similar = await recommendation_engine._content_based_filtering(
        user_installations=[adapter],
        count=count
    )
    
    return {
        "similar_adapters": [
            {
                "adapter_id": rec['adapter'].adapter_id,
                "name": rec['adapter'].display_name,
                "description": rec['adapter'].tagline,
                "rating": rec['adapter'].rating_average,
                "similarity": rec['score']
            }
            for rec in similar
        ]
    }


@router.get("/trending")
async def get_trending_adapters(
    period: str = Query("week", regex="^(day|week|month)$"),
    count: int = Query(10, ge=1, le=50)
):
    """
    获取趋势适配器
    
    基于最近的下载量增长、评分增长等
    """
    # 计算时间范围
    from datetime import datetime, timedelta
    
    period_map = {
        "day": timedelta(days=1),
        "week": timedelta(weeks=1),
        "month": timedelta(days=30)
    }
    
    since = datetime.utcnow() - period_map[period]
    
    # 统计最近的安装量
    trending = await db.query(
        MarketplaceAdapter,
        func.count(AdapterInstallation.id).label('recent_installs')
    ).join(
        AdapterInstallation
    ).filter(
        AdapterInstallation.installed_at >= since,
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.id
    ).order_by(
        desc('recent_installs')
    ).limit(count).all()
    
    return {
        "period": period,
        "trending_adapters": [
            {
                "adapter_id": adapter.adapter_id,
                "name": adapter.display_name,
                "description": adapter.tagline,
                "rating": adapter.rating_average,
                "recent_installs": recent_installs,
                "total_installs": adapter.install_count
            }
            for adapter, recent_installs in trending
        ]
    }
```

---

## 4.3 评分评论系统

### 4.3.1 评论管理

```python
# src/marketplace/api/reviews.py

from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel, Field

router = APIRouter(prefix="/api/v1/marketplace/reviews", tags=["reviews"])


class ReviewCreateRequest(BaseModel):
    """创建评论请求"""
    rating: int = Field(..., ge=1, le=5, description="评分 1-5星")
    title: Optional[str] = Field(None, max_length=200)
    content: Optional[str] = Field(None, max_length=5000)
    use_case: Optional[str] = Field(None, max_length=100)


@router.post("/adapters/{adapter_id}/reviews")
async def create_review(
    adapter_id: str,
    request: ReviewCreateRequest,
    current_user: User = Depends(get_current_user)
):
    """
    创建评论
    
    用户必须安装过该适配器才能评论
    """
    # 查找适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 检查是否已安装
    installation = await db.query(AdapterInstallation).filter_by(
        user_id=current_user.id,
        adapter_id=adapter.id
    ).first()
    
    if not installation:
        raise HTTPException(
            status_code=403,
            detail="You must install the adapter before reviewing"
        )
    
    # 检查是否已评论
    existing_review = await db.query(AdapterReview).filter_by(
        user_id=current_user.id,
        adapter_id=adapter.id
    ).first()
    
    if existing_review:
        raise HTTPException(
            status_code=400,
            detail="You have already reviewed this adapter"
        )
    
    # 创建评论
    review = AdapterReview(
        adapter_id=adapter.id,
        user_id=current_user.id,
        rating=request.rating,
        title=request.title,
        content=request.content,
        use_case=request.use_case,
        is_verified_user=current_user.is_verified
    )
    
    db.add(review)
    
    # 更新适配器评分统计
    await update_adapter_rating(adapter.id)
    
    await db.commit()
    
    return {"message": "Review created successfully"}


@router.get("/adapters/{adapter_id}/reviews")
async def get_adapter_reviews(
    adapter_id: str,
    sort_by: str = Query("newest", regex="^(newest|helpful|rating_high|rating_low)$"),
    min_rating: Optional[int] = Query(None, ge=1, le=5),
    page: int = Query(1, ge=1),
    page_size: int = Query(20, ge=1, le=100)
):
    """获取适配器的所有评论"""
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 构建查询
    query = db.query(AdapterReview).filter(
        AdapterReview.adapter_id == adapter.id,
        AdapterReview.is_visible == True
    )
    
    # 评分过滤
    if min_rating:
        query = query.filter(AdapterReview.rating >= min_rating)
    
    # 排序
    if sort_by == "newest":
        query = query.order_by(desc(AdapterReview.created_at))
    elif sort_by == "helpful":
        query = query.order_by(desc(AdapterReview.helpful_count))
    elif sort_by == "rating_high":
        query = query.order_by(desc(AdapterReview.rating))
    elif sort_by == "rating_low":
        query = query.order_by(asc(AdapterReview.rating))
    
    # 分页
    total = await query.count()
    reviews = await query.offset((page - 1) * page_size).limit(page_size).all()
    
    return {
        "total": total,
        "page": page,
        "page_size": page_size,
        "reviews": [
            {
                "id": review.id,
                "rating": review.rating,
                "title": review.title,
                "content": review.content,
                "author": {
                    "username": review.user.username,
                    "display_name": review.user.display_name,
                    "is_verified": review.user.is_verified
                },
                "helpful_count": review.helpful_count,
                "created_at": review.created_at.isoformat()
            }
            for review in reviews
        ]
    }
```

---

*文档第五部分完成*

**已完成内容**：
- ✅ Elasticsearch搜索系统（索引设计、多字段搜索、模糊匹配、自动补全）
- ✅ 搜索API（关键词、过滤、排序、分页、过滤器选项）
- ✅ 推荐系统（协同过滤、内容推荐、热门推荐、上下文推荐）
- ✅ 推荐API（个性化推荐、相似推荐、趋势推荐）
- ✅ 评分评论系统（创建评论、查询评论、评分统计）

**下一部分将包含**：
- 统计分析系统
- 开发者控制台
- SDK和工具
- 总结和最佳实践

是否继续输出第六部分（最终部分）？
# MetaWorkflow V2.0 适配器市场设计 - 第六部分（最终部分）

## 第五部分：统计分析与开发者工具

## 5.1 统计分析系统

### 5.1.1 ClickHouse数据仓库设计

```sql
-- 适配器使用事件表
CREATE TABLE adapter_usage_events (
    event_id UUID,
    event_time DateTime,
    event_type String,  -- install/uninstall/activate/deactivate/execute
    
    -- 适配器信息
    adapter_id String,
    adapter_name String,
    adapter_version String,
    
    -- 用户信息
    user_id UInt64,
    organization String,
    
    -- 上下文信息
    domain String,
    activity_type String,
    workflow_id String,
    
    -- 性能指标
    execution_time_ms UInt32,
    memory_usage_mb UInt32,
    
    -- 结果
    success Boolean,
    error_message String,
    
    -- 设备信息
    platform String,
    platform_version String,
    
    -- 时间分区
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, event_time);


-- 适配器下载事件表
CREATE TABLE adapter_download_events (
    download_id UUID,
    download_time DateTime,
    
    adapter_id String,
    adapter_version String,
    
    user_id UInt64,
    
    -- 来源
    source String,  -- marketplace/api/cli
    referrer String,
    
    -- 地理位置
    country String,
    region String,
    
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, download_time);


-- 适配器评分事件表
CREATE TABLE adapter_rating_events (
    rating_id UInt64,
    rating_time DateTime,
    
    adapter_id String,
    user_id UInt64,
    
    rating UInt8,  -- 1-5
    previous_rating UInt8,  -- 如果是更新评分
    
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, rating_time);
```

### 5.1.2 统计分析服务

```python
# src/marketplace/analytics/stats.py

from clickhouse_driver import Client
from datetime import datetime, timedelta
from typing import Dict, List, Any

class AdapterAnalyticsService:
    """适配器统计分析服务"""
    
    def __init__(self, clickhouse_client: Client):
        self.ch = clickhouse_client
    
    async def get_adapter_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict[str, Any]:
        """
        获取适配器统计数据
        
        Args:
            adapter_id: 适配器ID
            start_date: 开始日期
            end_date: 结束日期
            
        Returns:
            Dict: 统计数据
        """
        stats = {
            "overview": await self._get_overview_stats(adapter_id, start_date, end_date),
            "downloads": await self._get_download_stats(adapter_id, start_date, end_date),
            "usage": await self._get_usage_stats(adapter_id, start_date, end_date),
            "performance": await self._get_performance_stats(adapter_id, start_date, end_date),
            "rating_trend": await self._get_rating_trend(adapter_id, start_date, end_date),
            "user_demographics": await self._get_user_demographics(adapter_id, start_date, end_date)
        }
        
        return stats
    
    async def _get_overview_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """概览统计"""
        query = """
        SELECT
            countIf(event_type = 'install') as total_installs,
            countIf(event_type = 'uninstall') as total_uninstalls,
            countIf(event_type = 'execute') as total_executions,
            countIf(event_type = 'execute' AND success = 1) as successful_executions,
            countIf(event_type = 'execute' AND success = 0) as failed_executions,
            uniq(user_id) as unique_users,
            avg(execution_time_ms) as avg_execution_time,
            quantile(0.95)(execution_time_ms) as p95_execution_time
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        row = result[0]
        
        return {
            "total_installs": row[0],
            "total_uninstalls": row[1],
            "net_installs": row[0] - row[1],
            "total_executions": row[2],
            "successful_executions": row[3],
            "failed_executions": row[4],
            "success_rate": row[3] / row[2] if row[2] > 0 else 0,
            "unique_users": row[5],
            "avg_execution_time_ms": round(row[6], 2),
            "p95_execution_time_ms": round(row[7], 2)
        }
    
    async def _get_download_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """下载统计（按天）"""
        query = """
        SELECT
            date,
            count() as downloads,
            uniq(user_id) as unique_users
        FROM adapter_download_events
        WHERE adapter_id = %(adapter_id)s
          AND download_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY date
        ORDER BY date
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "daily_data": [
                {
                    "date": row[0].isoformat(),
                    "downloads": row[1],
                    "unique_users": row[2]
                }
                for row in result
            ]
        }
    
    async def _get_usage_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """使用统计"""
        # 按领域统计
        domain_query = """
        SELECT
            domain,
            count() as executions,
            uniq(user_id) as users
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY domain
        ORDER BY executions DESC
        LIMIT 10
        """
        
        domain_result = self.ch.execute(
            domain_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        # 按活动类型统计
        activity_query = """
        SELECT
            activity_type,
            count() as executions
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY activity_type
        ORDER BY executions DESC
        """
        
        activity_result = self.ch.execute(
            activity_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "by_domain": [
                {"domain": row[0], "executions": row[1], "users": row[2]}
                for row in domain_result
            ],
            "by_activity_type": [
                {"activity_type": row[0], "executions": row[1]}
                for row in activity_result
            ]
        }
    
    async def _get_performance_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """性能统计"""
        query = """
        SELECT
            toStartOfHour(event_time) as hour,
            avg(execution_time_ms) as avg_time,
            quantile(0.50)(execution_time_ms) as p50_time,
            quantile(0.95)(execution_time_ms) as p95_time,
            quantile(0.99)(execution_time_ms) as p99_time,
            avg(memory_usage_mb) as avg_memory
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY hour
        ORDER BY hour
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "hourly_performance": [
                {
                    "hour": row[0].isoformat(),
                    "avg_time_ms": round(row[1], 2),
                    "p50_time_ms": round(row[2], 2),
                    "p95_time_ms": round(row[3], 2),
                    "p99_time_ms": round(row[4], 2),
                    "avg_memory_mb": round(row[5], 2)
                }
                for row in result
            ]
        }
    
    async def _get_rating_trend(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """评分趋势"""
        query = """
        SELECT
            toStartOfDay(rating_time) as day,
            avg(rating) as avg_rating,
            count() as rating_count
        FROM adapter_rating_events
        WHERE adapter_id = %(adapter_id)s
          AND rating_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY day
        ORDER BY day
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "daily_ratings": [
                {
                    "date": row[0].isoformat(),
                    "avg_rating": round(row[1], 2),
                    "rating_count": row[2]
                }
                for row in result
            ]
        }
    
    async def _get_user_demographics(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """用户画像"""
        # 按组织统计
        org_query = """
        SELECT
            organization,
            uniq(user_id) as user_count,
            count() as execution_count
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY organization
        ORDER BY execution_count DESC
        LIMIT 10
        """
        
        org_result = self.ch.execute(
            org_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "top_organizations": [
                {
                    "organization": row[0],
                    "user_count": row[1],
                    "execution_count": row[2]
                }
                for row in org_result
            ]
        }


# 全局分析服务实例
analytics_service = AdapterAnalyticsService(clickhouse_client)
```

### 5.1.3 统计API

```python
# src/marketplace/api/analytics.py

from fastapi import APIRouter, Depends, HTTPException, Query
from datetime import datetime, timedelta

router = APIRouter(prefix="/api/v1/marketplace/analytics", tags=["analytics"])


@router.get("/adapters/{adapter_id}/stats")
async def get_adapter_statistics(
    adapter_id: str,
    start_date: Optional[datetime] = Query(None),
    end_date: Optional[datetime] = Query(None),
    current_user: User = Depends(get_current_user)
):
    """
    获取适配器统计数据
    
    仅适配器作者或管理员可访问
    """
    # 验证权限
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=403,
            detail="Only adapter author or admin can view statistics"
        )
    
    # 默认时间范围：最近30天
    if not end_date:
        end_date = datetime.utcnow()
    if not start_date:
        start_date = end_date - timedelta(days=30)
    
    from src.marketplace.analytics.stats import analytics_service
    
    stats = await analytics_service.get_adapter_stats(
        adapter_id=adapter_id,
        start_date=start_date,
        end_date=end_date
    )
    
    return stats


@router.get("/marketplace/overview")
async def get_marketplace_overview(
    current_user: User = Depends(get_admin_user)
):
    """
    获取市场整体概览
    
    仅管理员可访问
    """
    # 统计各类数据
    total_adapters = await db.query(MarketplaceAdapter).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).count()
    
    total_developers = await db.query(User).filter(
        User.is_developer == True
    ).count()
    
    total_downloads = await db.query(
        func.sum(MarketplaceAdapter.download_count)
    ).scalar()
    
    total_installs = await db.query(
        func.sum(MarketplaceAdapter.install_count)
    ).scalar()
    
    avg_rating = await db.query(
        func.avg(MarketplaceAdapter.rating_average)
    ).filter(
        MarketplaceAdapter.rating_count > 0
    ).scalar()
    
    # 分类统计
    category_stats = await db.query(
        MarketplaceAdapter.category,
        func.count(MarketplaceAdapter.id).label('count')
    ).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.category
    ).all()
    
    return {
        "overview": {
            "total_adapters": total_adapters,
            "total_developers": total_developers,
            "total_downloads": total_downloads or 0,
            "total_installs": total_installs or 0,
            "avg_rating": round(avg_rating, 2) if avg_rating else 0.0
        },
        "by_category": [
            {"category": cat, "count": count}
            for cat, count in category_stats
        ]
    }
```

---

## 5.2 开发者控制台

### 5.2.1 开发者仪表板API

```python
# src/marketplace/api/developer.py

from fastapi import APIRouter, Depends, HTTPException
from typing import List

router = APIRouter(prefix="/api/v1/marketplace/developer", tags=["developer"])


@router.get("/dashboard")
async def get_developer_dashboard(
    current_user: User = Depends(get_developer_user)
):
    """
    获取开发者仪表板数据
    
    包含：
    - 适配器列表
    - 总体统计
    - 最近评论
    - 待处理任务
    """
    # 获取开发者的所有适配器
    adapters = await db.query(MarketplaceAdapter).filter(
        MarketplaceAdapter.author_id == current_user.id
    ).all()
    
    adapter_summaries = []
    total_downloads = 0
    total_installs = 0
    total_ratings = 0
    
    for adapter in adapters:
        total_downloads += adapter.download_count
        total_installs += adapter.install_count
        total_ratings += adapter.rating_count
        
        adapter_summaries.append({
            "adapter_id": adapter.adapter_id,
            "name": adapter.display_name,
            "status": adapter.status.value,
            "version": adapter.current_version,
            "downloads": adapter.download_count,
            "installs": adapter.active_install_count,
            "rating": adapter.rating_average,
            "rating_count": adapter.rating_count,
            "published_at": adapter.published_at.isoformat() if adapter.published_at else None,
            "updated_at": adapter.updated_at.isoformat()
        })
    
    # 获取最近的评论
    recent_reviews = await db.query(AdapterReview).join(
        MarketplaceAdapter
    ).filter(
        MarketplaceAdapter.author_id == current_user.id,
        AdapterReview.is_visible == True
    ).order_by(
        desc(AdapterReview.created_at)
    ).limit(10).all()
    
    # 待审核的版本
    pending_versions = await db.query(AdapterVersion).join(
        MarketplaceAdapter
    ).filter(
        MarketplaceAdapter.author_id == current_user.id,
        AdapterVersion.review_status.in_([
            AdapterPublishStatus.SUBMITTED,
            AdapterPublishStatus.REVIEWING
        ])
    ).all()
    
    return {
        "summary": {
            "total_adapters": len(adapters),
            "published_adapters": len([a for a in adapters if a.status == AdapterPublishStatus.PUBLISHED]),
            "total_downloads": total_downloads,
            "total_installs": total_installs,
            "total_ratings": total_ratings,
            "avg_rating": sum(a.rating_average for a in adapters) / len(adapters) if adapters else 0
        },
        "adapters": adapter_summaries,
        "recent_reviews": [
            {
                "adapter_name": review.adapter.display_name,
                "rating": review.rating,
                "title": review.title,
                "user": review.user.display_name or review.user.username,
                "created_at": review.created_at.isoformat()
            }
            for review in recent_reviews
        ],
        "pending_reviews": [
            {
                "adapter_name": version.adapter.display_name,
                "version": version.version,
                "status": version.review_status.value,
                "submitted_at": version.created_at.isoformat()
            }
            for version in pending_versions
        ]
    }


@router.post("/adapters/{adapter_id}/reply-review/{review_id}")
async def reply_to_review(
    adapter_id: str,
    review_id: int,
    reply_text: str = Body(..., max_length=2000),
    current_user: User = Depends(get_developer_user)
):
    """
    回复用户评论
    
    只有适配器作者可以回复
    """
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id:
        raise HTTPException(
            status_code=403,
            detail="Only adapter author can reply to reviews"
        )
    
    review = await db.query(AdapterReview).get(review_id)
    
    if not review or review.adapter_id != adapter.id:
        raise HTTPException(status_code=404, detail="Review not found")
    
    # 添加回复
    review.author_reply = reply_text
    review.author_replied_at = datetime.utcnow()
    
    await db.commit()
    
    # 通知用户
    await notify_user_of_reply(review.user_id, adapter.display_name)
    
    return {"message": "Reply added successfully"}
```

---

## 5.3 开发者SDK和CLI工具

### 5.3.1 Python SDK

```python
# metaworkflow_sdk/client.py

"""
MetaWorkflow 适配器市场 SDK

安装：
    pip install metaworkflow-sdk

使用示例：
    from metaworkflow_sdk import MarketplaceClient
    
    client = MarketplaceClient(api_key="your_api_key")
    
    # 发布适配器
    client.publish_adapter(
        name="my-chemistry-adapter",
        version="1.0.0",
        package_path="./dist/my-adapter.zip"
    )
"""

import requests
from typing import Optional, Dict, Any
from pathlib import Path

class MarketplaceClient:
    """市场API客户端"""
    
    def __init__(
        self,
        api_key: str,
        base_url: str = "https://marketplace.metaworkflow.ai/api/v1"
    ):
        self.api_key = api_key
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {api_key}",
            "User-Agent": "MetaWorkflow-SDK/1.0"
        })
    
    def create_adapter(
        self,
        name: str,
        display_name: str,
        description: str,
        category: str,
        **kwargs
    ) -> Dict[str, Any]:
        """
        创建适配器（草稿）
        
        Args:
            name: 适配器名称（唯一标识）
            display_name: 显示名称
            description: 描述
            category: 分类
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 创建结果
        """
        payload = {
            "name": name,
            "display_name": display_name,
            "description": description,
            "category": category,
            **kwargs
        }
        
        response = self.session.post(
            f"{self.base_url}/marketplace/publish/adapters",
            json=payload
        )
        response.raise_for_status()
        
        return response.json()
    
    def upload_version(
        self,
        adapter_id: str,
        version: str,
        package_path: str,
        changelog: str,
        **kwargs
    ) -> Dict[str, Any]:
        """
        上传适配器版本
        
        Args:
            adapter_id: 适配器ID
            version: 版本号
            package_path: 包文件路径
            changelog: 变更日志
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 上传结果
        """
        package_file = Path(package_path)
        
        if not package_file.exists():
            raise FileNotFoundError(f"Package file not found: {package_path}")
        
        # 上传版本信息
        version_data = {
            "version": version,
            "changelog": changelog,
            **kwargs
        }
        
        # 上传文件
        with open(package_file, 'rb') as f:
            files = {'package_file': (package_file.name, f, 'application/zip')}
            
            response = self.session.post(
                f"{self.base_url}/marketplace/publish/adapters/{adapter_id}/versions",
                data=version_data,
                files=files
            )
        
        response.raise_for_status()
        
        return response.json()
    
    def submit_for_review(self, adapter_id: str) -> Dict[str, Any]:
        """
        提交审核
        
        Args:
            adapter_id: 适配器ID
        
        Returns:
            Dict: 提交结果
        """
        response = self.session.post(
            f"{self.base_url}/marketplace/publish/adapters/{adapter_id}/submit"
        )
        response.raise_for_status()
        
        return response.json()
    
    def publish_adapter(
        self,
        name: str,
        display_name: str,
        description: str,
        category: str,
        version: str,
        package_path: str,
        changelog: str,
        auto_submit: bool = True,
        **kwargs
    ) -> Dict[str, Any]:
        """
        一键发布适配器
        
        组合了创建、上传、提交审核三个步骤
        
        Args:
            name: 适配器名称
            display_name: 显示名称
            description: 描述
            category: 分类
            version: 版本号
            package_path: 包文件路径
            changelog: 变更日志
            auto_submit: 是否自动提交审核
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 发布结果
        """
        # 1. 创建适配器
        print(f"Creating adapter: {name}")
        create_result = self.create_adapter(
            name=name,
            display_name=display_name,
            description=description,
            category=category,
            **kwargs
        )
        
        adapter_id = create_result['adapter_id']
        print(f"✓ Adapter created: {adapter_id}")
        
        # 2. 上传版本
        print(f"Uploading version: {version}")
        upload_result = self.upload_version(
            adapter_id=adapter_id,
            version=version,
            package_path=package_path,
            changelog=changelog
        )
        print(f"✓ Version uploaded: {upload_result['version']}")
        
        # 3. 提交审核
        if auto_submit:
            print("Submitting for review...")
            submit_result = self.submit_for_review(adapter_id)
            print(f"✓ Submitted for review: {submit_result['status']}")
            
            return {
                "adapter_id": adapter_id,
                "version": version,
                "status": submit_result['status'],
                "message": "Adapter published and submitted for review"
            }
        else:
            return {
                "adapter_id": adapter_id,
                "version": version,
                "status": "draft",
                "message": "Adapter created. Use submit_for_review() to submit."
            }
    
    def get_adapter_stats(
        self,
        adapter_id: str,
        start_date: Optional[str] = None,
        end_date: Optional[str] = None
    ) -> Dict[str, Any]:
        """
        获取适配器统计数据
        
        Args:
            adapter_id: 适配器ID
            start_date: 开始日期 (ISO格式)
            end_date: 结束日期 (ISO格式)
        
        Returns:
            Dict: 统计数据
        """
        params = {}
        if start_date:
            params['start_date'] = start_date
        if end_date:
            params['end_date'] = end_date
        
        response = self.session.get(
            f"{self.base_url}/marketplace/analytics/adapters/{adapter_id}/stats",
            params=params
        )
        response.raise_for_status()
        
        return response.json()
```

### 5.3.2 CLI工具

```python
# metaworkflow_cli/main.py

"""
MetaWorkflow CLI 工具

安装：
    pip install metaworkflow-cli

命令：
    mw-cli init              # 初始化适配器项目
    mw-cli build             # 构建适配器包
    mw-cli publish           # 发布到市场
    mw-cli stats <adapter>   # 查看统计数据
"""

import click
import yaml
from pathlib import Path
from metaworkflow_sdk import MarketplaceClient

@click.group()
def cli():
    """MetaWorkflow 命令行工具"""
    pass


@cli.command()
@click.argument('name')
@click.option('--template', default='basic', help='项目模板')
def init(name: str, template: str):
    """初始化新的适配器项目"""
    click.echo(f"Creating adapter project: {name}")
    
    project_dir = Path(name)
    project_dir.mkdir(exist_ok=True)
    
    # 创建目录结构
    (project_dir / "src").mkdir(exist_ok=True)
    (project_dir / "tests").mkdir(exist_ok=True)
    (project_dir / "docs").mkdir(exist_ok=True)
    
    # 创建manifest.yaml
    manifest = {
        "name": name,
        "display_name": name.replace('-', ' ').title(),
        "version": "0.1.0",
        "description": f"Description for {name}",
        "author": {
            "name": "Your Name",
            "email": "your.email@example.com"
        },
        "compatibility": {
            "platform_version": ">=2.0.0"
        }
    }
    
    with open(project_dir / "manifest.yaml", 'w') as f:
        yaml.dump(manifest, f, default_flow_style=False)
    
    # 创建adapter.py模板
    adapter_template = '''from metaworkflow.adapters import DomainAdapter

class MyAdapter(DomainAdapter):
    """My custom adapter"""
    
    async def parse_request(self, request):
        # TODO: Implement parsing logic
        pass
    
    async def discover_content(self, topic, context):
        # TODO: Implement content discovery
        pass
'''
    
    with open(project_dir / "src" / "adapter.py", 'w') as f:
        f.write(adapter_template)
    
    click.echo(f"✓ Project created at {project_dir}")
    click.echo("\nNext steps:")
    click.echo(f"  cd {name}")
    click.echo("  # Edit src/adapter.py and manifest.yaml")
    click.echo("  mw-cli build")


@cli.command()
@click.option('--output', '-o', default='dist', help='输出目录')
def build(output: str):
    """构建适配器包"""
    import zipfile
    import shutil
    
    click.echo("Building adapter package...")
    
    # 读取manifest
    if not Path('manifest.yaml').exists():
        click.echo("Error: manifest.yaml not found", err=True)
        return
    
    with open('manifest.yaml', 'r') as f:
        manifest = yaml.safe_load(f)
    
    name = manifest['name']
    version = manifest['version']
    
    # 创建输出目录
    output_dir = Path(output)
    output_dir.mkdir(exist_ok=True)
    
    # 创建zip包
    package_path = output_dir / f"{name}-{version}.zip"
    
    with zipfile.ZipFile(package_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
        # 添加manifest
        zipf.write('manifest.yaml')
        
        # 添加源代码
        for py_file in Path('src').rglob('*.py'):
            zipf.write(py_file, py_file)
        
        # 添加文档
        if Path('README.md').exists():
            zipf.write('README.md')
    
    click.echo(f"✓ Package built: {package_path}")
    click.echo(f"  Size: {package_path.stat().st_size / 1024:.1f} KB")


@cli.command()
@click.option('--api-key', envvar='METAWORKFLOW_API_KEY', help='API密钥')
@click.option('--auto-submit/--no-auto-submit', default=True, help='自动提交审核')
def publish(api_key: str, auto_submit: bool):
    """发布适配器到市场"""
    if not api_key:
        click.echo("Error: API key required. Set METAWORKFLOW_API_KEY or use --api-key", err=True)
        return
    
    # 读取manifest
    with open('manifest.yaml', 'r') as f:
        manifest = yaml.safe_load(f)
    
    name = manifest['name']
    version = manifest['version']
    
    # 查找包文件
    package_path = Path('dist') / f"{name}-{version}.zip"
    
    if not package_path.exists():
        click.echo(f"Error: Package not found: {package_path}", err=True)
        click.echo("Run 'mw-cli build' first", err=True)
        return
    
    click.echo(f"Publishing {name} v{version}...")
    
    # 使用SDK发布
    client = MarketplaceClient(api_key=api_key)
    
    try:
        result = client.publish_adapter(
            name=name,
            display_name=manifest['display_name'],
            description=manifest['description'],
            category=manifest.get('category', 'other'),
            version=version,
            package_path=str(package_path),
            changelog=manifest.get('changelog', 'Initial release'),
            auto_submit=auto_submit
        )
        
        click.echo(f"✓ {result['message']}")
        click.echo(f"  Adapter ID: {result['adapter_id']}")
        click.echo(f"  Status: {result['status']}")
        
    except Exception as e:
        click.echo(f"Error: {str(e)}", err=True)


@cli.command()
@click.argument('adapter_id')
@click.option('--api-key', envvar='METAWORKFLOW_API_KEY', help='API密钥')
@click.option('--days', default=30, help='统计天数')
def stats(adapter_id: str, api_key: str, days: int):
    """查看适配器统计数据"""
    if not api_key:
        click.echo("Error: API key required", err=True)
        return
    
    from datetime import datetime, timedelta
    
    client = MarketplaceClient(api_key=api_key)
    
    end_date = datetime.utcnow()
    start_date = end_date - timedelta(days=days)
    
    stats_data = client.get_adapter_stats(
        adapter_id=adapter_id,
        start_date=start_date.isoformat(),
        end_date=end_date.isoformat()
    )
    
    overview = stats_data['overview']
    
    click.echo(f"\nStatistics for {adapter_id} (Last {days} days)")
    click.echo("=" * 50)
    click.echo(f"Downloads:        {overview['total_installs']}")
    click.echo(f"Active Installs:  {overview['net_installs']}")
    click.echo(f"Executions:       {overview['total_executions']}")
    click.echo(f"Success Rate:     {overview['success_rate']*100:.1f}%")
    click.echo(f"Unique Users:     {overview['unique_users']}")
    click.echo(f"Avg Exec Time:    {overview['avg_execution_time_ms']:.0f}ms")
    click.echo(f"P95 Exec Time:    {overview['p95_execution_time_ms']:.0f}ms")


if __name__ == '__main__':
    cli()
```

---

## 5.4 总结与最佳实践

### 5.4.1 架构设计总结

**核心设计原则**：

1. **三层架构（孵化器-鸡-蛋模型）**
   - **平台核心层**：提供通用AI能力（孵化器）
   - **适配器层**：领域专业化（鸡）
   - **内容层**：最终交付物（蛋）

2. **插件化扩展**
   - 标准化接口（DomainAdapter、NodeExecutor、ContentSource）
   - 动态加载机制
   - 依赖管理和版本控制

3. **质量保证**
   - 双重审核（自动检测 + 人工审核）
   - 安全扫描、性能测试、兼容性验证
   - 持续监控和质量评分

4. **开发者友好**
   - 完整的SDK和CLI工具
   - 清晰的文档和示例
   - 快速的发布流程

### 5.4.2 开发最佳实践

**适配器开发**：

```python
# 1. 良好的结构
my-adapter/
├── manifest.yaml          # 元数据声明
├── src/
│   ├── adapter.py        # 主适配器类
│   ├── content_sources/  # 内容源
│   ├── node_executors/   # 节点执行器
│   └── utils/            # 工具函数
├── tests/                 # 测试
├── docs/                  # 文档
└── README.md

# 2. 清晰的接口实现
class MyAdapter(DomainAdapter):
    """
    简洁的文档字符串
    
    说明：
    - 适配器用途
    - 支持的功能
    - 使用示例
    """
    
    async def parse_request(self, request):
        # 详细的注释
        # 清晰的错误处理
        # 合理的日志记录
        pass

# 3. 完善的测试
# tests/test_adapter.py
import pytest

@pytest.mark.asyncio
async def test_parse_request():
    adapter = MyAdapter()
    result = await adapter.parse_request({...})
    assert result is not None

# 4. 版本管理
# 遵循语义化版本（Semantic Versioning）
# MAJOR.MINOR.PATCH
# 1.0.0 -> 1.0.1 (bug fix)
# 1.0.1 -> 1.1.0 (new feature)
# 1.1.0 -> 2.0.0 (breaking change)
```

**性能优化**：

```python
# 1. 异步操作
async def discover_content(self, topic):
    # 并发请求多个数据源
    results = await asyncio.gather(
        self.source1.search(topic),
        self.source2.search(topic),
        self.source3.search(topic)
    )
    return self._merge_results(results)

# 2. 缓存机制
from functools import lru_cache

@lru_cache(maxsize=128)
def _get_static_data(key):
    # 缓存静态数据
    pass

# 3. 批量处理
async def process_batch(self, items):
    # 批量处理，减少API调用
    batch_size = 10
    for i in range(0, len(items), batch_size):
        batch = items[i:i+batch_size]
        await self._process(batch)
```

**安全最佳实践**：

```python
# 1. 输入验证
def parse_request(self, request):
    # 验证输入
    if not isinstance(request, dict):
        raise ValueError("Request must be a dict")
    
    # 清理输入
    topic = request.get('topic', '').strip()
    if not topic:
        raise ValueError("Topic is required")
    
    return topic

# 2. 敏感信息处理
# 永远不要硬编码密钥
# ❌ 错误
API_KEY = "sk-1234567890abcdef"

# ✅ 正确
import os
API_KEY = os.environ.get('EXTERNAL_API_KEY')

# 3. 错误处理
try:
    result = await external_api.call()
except Exception as e:
    # 记录错误，但不暴露敏感信息
    logger.error(f"API call failed: {type(e).__name__}")
    raise AdapterError("Failed to fetch data")
```

### 5.4.3 市场运营建议

**for 平台运营方**：

1. **质量控制**
   - 严格的审核标准
   - 定期安全扫描
   - 用户反馈机制

2. **开发者支持**
   - 详细文档和教程
   - 活跃的开发者社区
   - 技术支持渠道

3. **激励机制**
   - 优质适配器推荐
   - 开发者排行榜
   - 商业化分成（未来）

**for 适配器开发者**：

1. **用户第一**
   - 清晰的文档
   - 快速的响应
   - 持续的更新

2. **质量优先**
   - 充分测试
   - 性能优化
   - 安全防护

3. **社区参与**
   - 及时回复评论
   - 分享最佳实践
   - 贡献示例代码

---

## 6. 文档总结

### 6.1 完整文档结构

本设计文档共**6个部分**，涵盖：

**Part 1 - 架构基础**：
- 三层架构模型
- 适配器包结构
- manifest.yaml规范
- 元数据管理

**Part 2 - 接口规范**：
- DomainAdapter接口（15个方法）
- NodeExecutor接口
- ContentSource接口
- 详细示例

**Part 3 - 注册与生命周期**：
- 注册表设计
- 自动发现机制
- 依赖管理
- 生命周期状态机

**Part 4 - 市场平台**：
- 4层架构
- 完整数据模型
- 发布流程
- 自动化检测

**Part 5 - 用户体验**：
- Elasticsearch搜索
- 推荐系统
- 评分评论
- 统计分析

**Part 6 - 工具与总结**：
- ClickHouse分析
- 开发者仪表板
- Python SDK
- CLI工具
- 最佳实践

### 6.2 关键指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| 适配器审核时间 | < 48小时 | 从提交到审核完成 |
| 安装成功率 | > 99% | 用户安装失败率 < 1% |
| 平台可用性 | 99.9% | 年停机时间 < 8.76小时 |
| 搜索响应时间 | < 200ms | P95响应时间 |
| 推荐准确率 | > 30% | 用户点击率 |
| 开发者满意度 | > 4.0/5.0 | 开发者调查评分 |

### 6.3 实施路线图

**Phase 1 - 基础设施（2周）**：
- ✅ 数据库设计和迁移
- ✅ API框架搭建
- ✅ 认证授权系统

**Phase 2 - 核心功能（4周）**：
- ✅ 适配器注册和管理
- ✅ 发布和审核流程
- ✅ 搜索功能

**Phase 3 - 增强功能（3周）**：
- ✅ 推荐系统
- ✅ 评分评论
- ✅ 统计分析

**Phase 4 - 开发者工具（2周）**：
- ✅ SDK开发
- ✅ CLI工具
- ✅ 文档完善

**Phase 5 - 上线准备（1周）**：
- ⏳ 性能测试
- ⏳ 安全审计
- ⏳ 生产部署

### 6.4 成功标准

**技术指标**：
- [x] 完整的架构设计
- [x] 标准化接口定义
- [x] 自动化质量检测
- [x] 可扩展性设计

**业务指标**：
- [ ] 上线3个月：20+ 适配器
- [ ] 上线6个月：50+ 适配器
- [ ] 上线12个月：100+ 适配器
- [ ] 开发者社区：500+ 成员

**用户体验**：
- [x] 简单的发布流程
- [x] 智能推荐
- [x] 快速搜索
- [x] 详细统计

---

## 结语

MetaWorkflow V2.0 适配器市场设计遵循**开放、标准、质量、生态**的核心理念：

- **开放**：欢迎第三方开发者贡献适配器
- **标准**：统一的接口和规范，保证兼容性
- **质量**：严格的审核和持续的监控
- **生态**：构建健康的开发者和用户生态

通过适配器市场，MetaWorkflow可以快速扩展到新的教育领域，从**K12基础教育**延伸到**高等教育、职业培训、企业培训**等多个场景，真正实现**"孵化器孵化无限可能"**的愿景。

---

**文档版本**: v1.0  
**最后更新**: 2025-12-09  
**作者**: MetaWorkflow架构团队  
**状态**: ✅ 完成

*感谢您阅读完整的适配器市场设计文档！*
