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
