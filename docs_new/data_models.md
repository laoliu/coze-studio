# 通用领域内容生成·元工作流平台 - 数据模型设计 ⭐ V2.0升级

**文档版本**: V2.0 ⭐ 全面升级  
**日期**: 2025-12-09  
**作者**: Meta-Workflow数据团队

**V2.0重大更新**:
- 🔥 **架构升级**：从K12垂直工具升级到通用多领域平台
- 🔥 **新增领域适配器模型**（8个核心类）：支持K12/美术史/畅销书/职业培训等
- 🔥 **新增AI能力模型**（20+个类）：AI目标生成器、内容发现引擎、叙述引擎
- 🔥 **新增多源内容模型**（10个类）：支持Wikipedia、WikiArt、YouTube等10+数据源
- 🔥 **课程规划数据模型**（5个核心模型）
- 🔥 **模板本地化数据模型**（4个核心模型）
- 🔥 **组件扩展数据模型**（3个核心模型）
- 📊 **新增PostgreSQL表结构**：40+个新表

**V1.0保留内容**:
- ✅ 核心数据模型（WorkflowContext等）
- ✅ 工作流配置模型
- ✅ 执行状态模型
- ✅ 内容输出模型

---

## 目录

1. [核心数据模型](#1-核心数据模型)
2. [工作流配置模型](#2-工作流配置模型)
3. [执行状态模型](#3-执行状态模型)
4. [内容输出模型](#4-内容输出模型)
5. [课程规划数据模型](#5-课程规划数据模型)
6. [模板本地化数据模型](#6-模板本地化数据模型)
7. [组件扩展数据模型](#7-组件扩展数据模型)
8. [**V2.0新增：领域适配器数据模型**](#8-v20新增领域适配器数据模型) ⭐ 新增
9. [**V2.0新增：AI目标生成器数据模型**](#9-v20新增ai目标生成器数据模型) ⭐ 新增
10. [**V2.0新增：AI内容发现引擎数据模型**](#10-v20新增ai内容发现引擎数据模型) ⭐ 新增
11. [**V2.0新增：故事化叙述引擎数据模型**](#11-v20新增故事化叙述引擎数据模型) ⭐ 新增
12. [**V2.0新增：多源内容检索数据模型**](#12-v20新增多源内容检索数据模型) ⭐ 新增
13. [数据库Schema](#13-数据库schema)

---

## 1. 核心数据模型

### 1.1 WorkflowContext (工作流上下文)

所有节点共享的上下文对象，贯穿整个执行流程。

```python
from dataclasses import dataclass, field
from typing import Dict, Any, Optional, List
from datetime import datetime
from enum import Enum

class Subject(str, Enum):
    """学科枚举"""
    MATH = "数学"
    CHINESE = "语文"
    ENGLISH = "英语"
    SCIENCE = "科学"
    HISTORY = "历史"
    GEOGRAPHY = "地理"
    PHYSICS = "物理"
    CHEMISTRY = "化学"
    BIOLOGY = "生物"
    ART = "美术"
    MUSIC = "音乐"
    PE = "体育"

class Region(str, Enum):
    """地区枚举"""
    CN = "中国"
    US = "美国"
    UK = "英国"
    SG = "新加坡"
    JP = "日本"
    KR = "韩国"

class PedagogyModel(str, Enum):
    """教学法枚举"""
    INQUIRY = "探究式"
    LECTURE = "讲授式"
    FLIPPED = "翻转课堂"
    PROJECT = "项目式"
    GAMIFICATION = "游戏化"

@dataclass
class WorkflowContext:
    """工作流上下文"""
    
    # === 基础属性 ===
    instance_id: str  # 工作流实例ID
    subject: Subject  # 学科
    grade: int  # 年级 (1-12)
    region: Region  # 地区
    language: str = "zh-CN"  # 语言
    
    # === 教学属性 ===
    topics: List[str] = field(default_factory=list)  # 知识点列表
    lesson_type: str = "新授课"  # 课型
    duration: int = 45  # 时长(分钟)，微课建议3-15分钟，常规课45-90分钟
    textbook_version: Optional[str] = None  # 教材版本
    
    # === 教学法属性 ===
    pedagogy_model: PedagogyModel = PedagogyModel.LECTURE
    interaction_level: str = "medium"  # low/medium/high
    bloom_level: str = "application"  # remember/understand/apply/analyze/evaluate/create
    
    # === 受众属性 ===
    age_range: str = ""  # 年龄段
    interests: List[str] = field(default_factory=list)  # 兴趣标签
    common_errors: List[str] = field(default_factory=list)  # 常见错误
    
    # === 执行状态 ===
    node_outputs: Dict[str, Any] = field(default_factory=dict)  # 各节点的输出
    metadata: Dict[str, Any] = field(default_factory=dict)  # 元数据
    created_at: datetime = field(default_factory=datetime.utcnow)
    
    def set_node_output(self, node_id: str, output: Any):
        """设置节点输出"""
        self.node_outputs[node_id] = output
    
    def get_node_output(self, node_id: str) -> Optional[Any]:
        """获取节点输出"""
        return self.node_outputs.get(node_id)
    
    def grade_to_age(self) -> int:
        """年级转年龄"""
        return self.grade + 6  # 简化映射
```

### 1.2 PersonaCard (用户画像)

```python
@dataclass
class CognitiveDimension:
    """认知维度"""
    abstract_thinking: str  # concrete/transitioning/formal
    attention_span_minutes: int  # 注意力持续时间
    prerequisite_mastery: int  # 前置知识掌握度 (0-100)
    learning_pace: str  # slow/medium/fast

@dataclass
class EmotionalDimension:
    """情感维度"""
    subject_interest: str  # low/medium/high
    motivation_type: str  # extrinsic/intrinsic/mixed
    frustration_tolerance: str  # low/medium/high
    confidence_level: str  # low/medium/high

@dataclass
class BehavioralDimension:
    """行为维度"""
    learning_styles: List[str]  # visual/auditory/kinesthetic
    interaction_preference: str  # individual/pair/group
    common_error_patterns: List[str]
    preferred_activities: List[str]

@dataclass
class PersonaCard:
    """用户画像卡片"""
    persona_id: str
    target_group: str  # 目标群体描述
    
    cognitive: CognitiveDimension
    emotional: EmotionalDimension
    behavioral: BehavioralDimension
    
    # 推荐策略
    recommended_narrative_style: str  # 推荐叙事风格
    recommended_pacing: str  # 推荐节奏
    recommended_scaffolding_level: str  # 推荐脚手架级别
    
    key_considerations: List[str]  # 关键注意事项
    
    def to_dict(self) -> Dict[str, Any]:
        """转为字典"""
        return {
            "persona_id": self.persona_id,
            "target_group": self.target_group,
            "cognitive": {
                "abstract_thinking": self.cognitive.abstract_thinking,
                "attention_span_minutes": self.cognitive.attention_span_minutes,
                "prerequisite_mastery": self.cognitive.prerequisite_mastery,
                "learning_pace": self.cognitive.learning_pace
            },
            "emotional": {
                "subject_interest": self.emotional.subject_interest,
                "motivation_type": self.emotional.motivation_type,
                "frustration_tolerance": self.emotional.frustration_tolerance,
                "confidence_level": self.emotional.confidence_level
            },
            "behavioral": {
                "learning_styles": self.behavioral.learning_styles,
                "interaction_preference": self.behavioral.interaction_preference,
                "common_error_patterns": self.behavioral.common_error_patterns,
                "preferred_activities": self.behavioral.preferred_activities
            },
            "recommendations": {
                "narrative_style": self.recommended_narrative_style,
                "pacing": self.recommended_pacing,
                "scaffolding_level": self.recommended_scaffolding_level
            },
            "key_considerations": self.key_considerations
        }
```

### 1.3 LessonObjectives (教学目标)

```python
@dataclass
class Objective:
    """单个教学目标"""
    dimension: str  # knowledge/skill/attitude (知识/技能/情感)
    bloom_level: str  # 布鲁姆层级
    description: str  # 目标描述
    assessment_criteria: str  # 评估标准

@dataclass
class LessonObjectives:
    """教学目标集合"""
    knowledge_objectives: List[Objective]  # 知识目标
    skill_objectives: List[Objective]  # 能力目标
    attitude_objectives: List[Objective]  # 情感目标
    
    core_competencies: List[str]  # 核心素养
    curriculum_alignment: Dict[str, str]  # 课标对齐
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "knowledge": [
                {"level": obj.bloom_level, "description": obj.description, "assessment": obj.assessment_criteria}
                for obj in self.knowledge_objectives
            ],
            "skill": [
                {"level": obj.bloom_level, "description": obj.description, "assessment": obj.assessment_criteria}
                for obj in self.skill_objectives
            ],
            "attitude": [
                {"level": obj.bloom_level, "description": obj.description, "assessment": obj.assessment_criteria}
                for obj in self.attitude_objectives
            ],
            "core_competencies": self.core_competencies,
            "curriculum_alignment": self.curriculum_alignment
        }
```

### 1.4 LessonScript (课程脚本)

```python
@dataclass
class InteractionActivity:
    """互动活动"""
    type: str  # question/discussion/experiment/game/exercise
    description: str
    duration: int  # 分钟
    materials: List[str]  # 所需材料
    expected_outcome: str  # 预期结果

@dataclass
class LessonSection:
    """课程环节"""
    section_id: str
    section_type: str  # introduction/exploration/explanation/elaboration/evaluation
    title: str
    duration: int  # 分钟
    
    narrative: str  # 叙事内容
    teacher_script: str  # 教师话术
    student_activities: List[str]  # 学生活动
    
    interactions: List[InteractionActivity]  # 互动活动
    media_resources: List[str]  # 媒体资源URL
    
    key_points: List[str]  # 关键知识点
    potential_difficulties: List[str]  # 可能的困难点
    scaffolding_strategies: List[str]  # 脚手架策略

@dataclass
class LessonScript:
    """完整课程脚本"""
    lesson_id: str
    title: str
    subtitle: Optional[str]
    
    metadata: Dict[str, Any]  # 元数据 (学科/年级/时长等)
    
    story_background: Optional[str]  # 故事化背景
    
    sections: List[LessonSection]  # 各个环节
    
    materials_list: List[str]  # 完整材料清单
    homework: Optional[str]  # 课后作业
    extension_activities: List[str]  # 拓展活动
    
    qa_pairs: List[Dict[str, str]]  # 预设问答对
    
    references: List[Dict[str, str]]  # 参考资料
    
    def total_duration(self) -> int:
        """计算总时长"""
        return sum(section.duration for section in self.sections)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "lesson_id": self.lesson_id,
            "title": self.title,
            "subtitle": self.subtitle,
            "metadata": self.metadata,
            "story_background": self.story_background,
            "sections": [
                {
                    "section_id": s.section_id,
                    "type": s.section_type,
                    "title": s.title,
                    "duration": s.duration,
                    "narrative": s.narrative,
                    "teacher_script": s.teacher_script,
                    "student_activities": s.student_activities,
                    "interactions": [
                        {
                            "type": ia.type,
                            "description": ia.description,
                            "duration": ia.duration,
                            "materials": ia.materials,
                            "expected_outcome": ia.expected_outcome
                        }
                        for ia in s.interactions
                    ],
                    "media_resources": s.media_resources,
                    "key_points": s.key_points,
                    "potential_difficulties": s.potential_difficulties,
                    "scaffolding_strategies": s.scaffolding_strategies
                }
                for s in self.sections
            ],
            "materials_list": self.materials_list,
            "homework": self.homework,
            "extension_activities": self.extension_activities,
            "qa_pairs": self.qa_pairs,
            "references": self.references,
            "total_duration": self.total_duration()
        }
```

---

## 2. 工作流配置模型

### 2.1 WorkflowDefinition (工作流定义)

```python
@dataclass
class NodeCondition:
    """节点执行条件"""
    field: str  # 条件字段 (如 "subject", "grade")
    operator: str  # 运算符 (==, !=, in, not in, >, <)
    value: Any  # 比较值

@dataclass
class NodeConfig:
    """节点配置"""
    node_id: str
    node_type: str  # 节点类型 (PersonaNode, ContentNode等)
    depends_on: List[str] = field(default_factory=list)  # 依赖的节点ID
    
    condition: Optional[NodeCondition] = None  # 执行条件
    
    config: Dict[str, Any] = field(default_factory=dict)  # 节点特定配置
    
    retry_policy: Dict[str, Any] = field(default_factory=dict)  # 重试策略
    timeout: int = 60  # 超时时间(秒)

@dataclass
class WorkflowDefinition:
    """工作流定义"""
    workflow_id: str
    name: str
    version: str
    description: str
    
    # 触发条件
    trigger_conditions: Dict[str, Any]
    
    # 全局配置
    global_config: Dict[str, Any]
    
    # 节点列表
    nodes: List[NodeConfig]
    
    # 输出配置
    output_config: Dict[str, Any]
    
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "workflow_id": self.workflow_id,
            "name": self.name,
            "version": self.version,
            "description": self.description,
            "trigger_conditions": self.trigger_conditions,
            "global_config": self.global_config,
            "nodes": [
                {
                    "node_id": n.node_id,
                    "node_type": n.node_type,
                    "depends_on": n.depends_on,
                    "condition": {
                        "field": n.condition.field,
                        "operator": n.condition.operator,
                        "value": n.condition.value
                    } if n.condition else None,
                    "config": n.config,
                    "retry_policy": n.retry_policy,
                    "timeout": n.timeout
                }
                for n in self.nodes
            ],
            "output_config": self.output_config
        }
```

### 2.2 Rule (规则定义)

```python
@dataclass
class RuleAction:
    """规则动作"""
    action_type: str  # insert_node/set_config/bind_resource/set_template
    parameters: Dict[str, Any]

@dataclass
class Rule:
    """规则"""
    rule_id: str
    name: str
    priority: int  # 优先级 (数字越大优先级越高)
    
    conditions: Dict[str, Any]  # 条件表达式
    
    actions: List[RuleAction]  # 触发的动作
    
    enabled: bool = True
    
    def matches(self, context: WorkflowContext) -> bool:
        """检查是否匹配"""
        for key, value in self.conditions.items():
            context_value = getattr(context, key, None)
            
            if context_value is None:
                return False
            
            if isinstance(value, list):
                if context_value not in value:
                    return False
            else:
                if context_value != value:
                    return False
        
        return True
```

---

## 2.2 工作流模板模型 ⭐ 新增

### 2.2.1 WorkflowTemplate (工作流模板)

```python
@dataclass
class WorkflowTemplate:
    """工作流模板 - 从元工作流实例生成的可复用模板"""
    
    template_id: str  # 唯一标识
    template_name: str  # 模板名称（支持唯一性查找）
    version: str  # 版本号 (如 "1.0", "1.1", "2.0")
    base_template_id: Optional[str] = None  # 父模板ID（用于版本链）
    
    # 元数据
    metadata: Dict[str, Any] = field(default_factory=dict)
    # metadata 结构示例:
    # {
    #     "subject": "数学",
    #     "grade_range": [1, 3],  # 适用年级范围
    #     "region": "中国",
    #     "created_by": "user_id",
    #     "created_at": "2025-01-01T00:00:00",
    #     "updated_at": "2025-01-15T10:30:00",
    #     "description": "小学低年级数学课件生成模板",
    #     "source_instance_id": "instance_123",  # 源工作流实例
    #     "tags": ["数学", "低年级", "基础"],
    #     "is_public": True,  # 是否公开到模板市场
    #     "fork_count": 5,  # 被复制次数
    #     "popularity_score": 8.5  # 流行度评分
    # }
    
    # 工作流结构
    nodes: List[Dict[str, Any]]  # 节点配置列表
    edges: List[Dict[str, Any]]  # 边（依赖关系）列表
    
    # 参数定义
    parameters: Dict[str, Any] = field(default_factory=dict)
    # parameters 结构示例:
    # {
    #     "required": ["topics", "duration"],  # 必填参数
    #     "optional": ["textbook_version", "enable_story"],  # 可选参数
    #     "defaults": {"duration": 45, "enable_story": True},  # 默认值
    #     "constraints": {
    #         "duration": {"min": 20, "max": 90},
    #         "grade": {"min": 1, "max": 12}
    #     }
    # }
    
    # 执行统计
    execution_history: Dict[str, Any] = field(default_factory=dict)
    # execution_history 结构示例:
    # {
    #     "total_runs": 152,
    #     "success_rate": 0.96,
    #     "avg_duration": 120.5,  # 秒
    #     "avg_cost": 0.35,  # 美元
    #     "avg_quality_score": 8.2,
    #     "last_run_at": "2025-01-20T14:30:00"
    # }
    
    # 用户反馈
    ratings: Dict[str, Any] = field(default_factory=dict)
    # ratings 结构示例:
    # {
    #     "avg_score": 4.5,
    #     "total_ratings": 42,
    #     "rating_distribution": {"5": 25, "4": 12, "3": 3, "2": 1, "1": 1}
    # }
    
    # 变更记录
    changelog: List[Dict[str, Any]] = field(default_factory=list)
    # changelog 结构示例:
    # [
    #     {
    #         "version": "1.1",
    #         "date": "2025-01-10T09:00:00",
    #         "author": "user_456",
    #         "changes": "优化了故事设计节点的提示词",
    #         "modifications": {
    #             "nodes": [{"action": "update", "node_id": "story_design", ...}]
    #         }
    #     }
    # ]
    
    # 优化建议
    optimization_suggestions: List[Dict[str, Any]] = field(default_factory=list)
    # 从执行数据分析生成的优化建议
    
    def to_dict(self) -> Dict[str, Any]:
        """转换为字典"""
        return {
            "template_id": self.template_id,
            "template_name": self.template_name,
            "version": self.version,
            "base_template_id": self.base_template_id,
            "metadata": self.metadata,
            "nodes": self.nodes,
            "edges": self.edges,
            "parameters": self.parameters,
            "execution_history": self.execution_history,
            "ratings": self.ratings,
            "changelog": self.changelog,
            "optimization_suggestions": self.optimization_suggestions
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'WorkflowTemplate':
        """从字典创建实例"""
        return cls(**data)
```

### 2.2.2 TemplateExecutionRecord (模板执行记录)

```python
@dataclass
class TemplateExecutionRecord:
    """模板执行记录 - 用于追踪模板使用和优化"""
    
    record_id: str
    instance_id: str  # 工作流实例ID
    template_id: str
    template_version: str
    
    executed_at: datetime
    executed_by: str  # 用户ID
    
    # 执行轨迹
    execution_trace: Dict[str, Any]
    # execution_trace 结构示例:
    # {
    #     "nodes": [
    #         {
    #             "node_id": "persona_analysis",
    #             "node_type": "llm_generation",
    #             "started_at": "2025-01-20T10:00:00",
    #             "completed_at": "2025-01-20T10:00:15",
    #             "duration": 15.2,
    #             "status": "success",
    #             "retry_count": 0,
    #             "input_size": 1024,
    #             "output_size": 2048,
    #             "model_used": "gpt-4",
    #             "tokens": 500,
    #             "cost": 0.01
    #         },
    #         # ... 更多节点
    #     ],
    #     "total_duration": 125.8,
    #     "success": True,
    #     "error_nodes": []
    # }
    
    # 资源消耗
    resource_usage: Dict[str, Any]
    # resource_usage 结构示例:
    # {
    #     "total_tokens": 15000,
    #     "total_cost": 0.35,
    #     "models_used": ["gpt-4", "deepseek-v2", "bert-base-chinese"],
    #     "cache_hits": 3,
    #     "cache_misses": 7
    # }
    
    # 质量指标
    quality_metrics: Dict[str, Any]
    # quality_metrics 结构示例:
    # {
    #     "overall_score": 8.5,
    #     "content_quality": 9.0,
    #     "logic_consistency": 8.0,
    #     "factual_accuracy": 9.5,
    #     "compliance_score": 10.0,
    #     "user_rating": 4.5  # 用户评分（如果有）
    # }
    
    # 用户反馈
    user_feedback: Optional[Dict[str, Any]] = None
    # user_feedback 结构示例:
    # {
    #     "rating": 5,
    #     "comments": "生成的课件非常好，故事设计很吸引人",
    #     "issues": ["个别地方需要微调"],
    #     "submitted_at": "2025-01-20T15:00:00"
    # }
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "record_id": self.record_id,
            "instance_id": self.instance_id,
            "template_id": self.template_id,
            "template_version": self.template_version,
            "executed_at": self.executed_at.isoformat(),
            "executed_by": self.executed_by,
            "execution_trace": self.execution_trace,
            "resource_usage": self.resource_usage,
            "quality_metrics": self.quality_metrics,
            "user_feedback": self.user_feedback
        }
```

### 2.2.3 OptimizationSuggestion (优化建议)

```python
@dataclass
class OptimizationSuggestion:
    """优化建议 - 从执行数据分析生成"""
    
    suggestion_id: str
    template_id: str
    
    # 建议类型
    suggestion_type: str  # "performance" | "quality" | "cost" | "structure"
    severity: str  # "high" | "medium" | "low"
    
    # 问题描述
    issue: str
    affected_nodes: List[str] = field(default_factory=list)
    
    # 建议内容
    suggestion: str
    expected_improvement: Dict[str, Any] = field(default_factory=dict)
    # expected_improvement 示例:
    # {
    #     "performance": {"duration_reduction": "30%"},
    #     "cost": {"cost_reduction": "15%"},
    #     "quality": {"score_improvement": "+0.5"}
    # }
    
    # 支持数据
    supporting_data: Dict[str, Any] = field(default_factory=dict)
    
    # 状态
    status: str = "pending"  # "pending" | "applied" | "rejected" | "testing"
    created_at: datetime = field(default_factory=datetime.utcnow)
    
    # 应用记录
    applied_at: Optional[datetime] = None
    applied_by: Optional[str] = None
    result: Optional[Dict[str, Any]] = None
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "suggestion_id": self.suggestion_id,
            "template_id": self.template_id,
            "suggestion_type": self.suggestion_type,
            "severity": self.severity,
            "issue": self.issue,
            "affected_nodes": self.affected_nodes,
            "suggestion": self.suggestion,
            "expected_improvement": self.expected_improvement,
            "supporting_data": self.supporting_data,
            "status": self.status,
            "created_at": self.created_at.isoformat(),
            "applied_at": self.applied_at.isoformat() if self.applied_at else None,
            "applied_by": self.applied_by,
            "result": self.result
        }
```

---

## 3. 执行状态模型

### 3.1 WorkflowInstance (工作流实例)

```python
class InstanceStatus(str, Enum):
    """实例状态"""
    PENDING = "pending"  # 待执行
    RUNNING = "running"  # 执行中
    COMPLETED = "completed"  # 已完成
    FAILED = "failed"  # 失败
    CANCELLED = "cancelled"  # 已取消

@dataclass
class NodeExecution:
    """节点执行记录"""
    node_id: str
    status: str  # pending/running/completed/failed
    
    started_at: Optional[datetime] = None
    completed_at: Optional[datetime] = None
    
    input_data: Dict[str, Any] = field(default_factory=dict)
    output_data: Optional[Any] = None
    
    error: Optional[str] = None
    retry_count: int = 0
    
    metrics: Dict[str, Any] = field(default_factory=dict)  # 性能指标

@dataclass
class WorkflowInstance:
    """工作流实例"""
    instance_id: str
    workflow_id: str
    
    status: InstanceStatus
    
    context: WorkflowContext
    
    node_executions: Dict[str, NodeExecution] = field(default_factory=dict)
    
    created_at: datetime = field(default_factory=datetime.utcnow)
    started_at: Optional[datetime] = None
    completed_at: Optional[datetime] = None
    
    error_log: List[str] = field(default_factory=list)
    
    metrics: Dict[str, Any] = field(default_factory=dict)
    
    def get_progress(self) -> float:
        """获取执行进度 (0-1)"""
        if not self.node_executions:
            return 0.0
        
        completed = sum(
            1 for ne in self.node_executions.values()
            if ne.status == "completed"
        )
        total = len(self.node_executions)
        
        return completed / total if total > 0 else 0.0
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "instance_id": self.instance_id,
            "workflow_id": self.workflow_id,
            "status": self.status.value,
            "progress": self.get_progress(),
            "created_at": self.created_at.isoformat(),
            "started_at": self.started_at.isoformat() if self.started_at else None,
            "completed_at": self.completed_at.isoformat() if self.completed_at else None,
            "node_executions": {
                node_id: {
                    "status": ne.status,
                    "started_at": ne.started_at.isoformat() if ne.started_at else None,
                    "completed_at": ne.completed_at.isoformat() if ne.completed_at else None,
                    "retry_count": ne.retry_count,
                    "error": ne.error
                }
                for node_id, ne in self.node_executions.items()
            },
            "metrics": self.metrics
        }
```

---

## 4. 内容输出模型

### 4.1 GeneratedContent (生成内容)

```python
@dataclass
class ContentMetadata:
    """内容元数据"""
    subject: str
    grade: int
    region: str
    topics: List[str]
    
    workflow_id: str
    instance_id: str
    
    generated_at: datetime
    generator_version: str
    
    quality_score: Optional[float] = None  # 质量评分
    compliance_status: str = "pending"  # pending/passed/failed
    
    models_used: List[str] = field(default_factory=list)  # 使用的模型
    total_tokens: int = 0  # 总token消耗
    total_cost: float = 0.0  # 总成本

@dataclass
class GeneratedContent:
    """生成的内容"""
    content_id: str
    content_type: str  # lesson_plan/script/exercise/multimedia
    
    metadata: ContentMetadata
    
    data: Dict[str, Any]  # 实际内容数据
    
    # 审核信息
    review_status: str = "pending"  # pending/approved/rejected
    reviewer_feedback: Optional[str] = None
    
    # 使用统计
    view_count: int = 0
    usage_count: int = 0
    avg_rating: Optional[float] = None
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "content_id": self.content_id,
            "content_type": self.content_type,
            "metadata": {
                "subject": self.metadata.subject,
                "grade": self.metadata.grade,
                "region": self.metadata.region,
                "topics": self.metadata.topics,
                "workflow_id": self.metadata.workflow_id,
                "instance_id": self.metadata.instance_id,
                "generated_at": self.metadata.generated_at.isoformat(),
                "quality_score": self.metadata.quality_score,
                "compliance_status": self.metadata.compliance_status,
                "models_used": self.metadata.models_used,
                "total_tokens": self.metadata.total_tokens,
                "total_cost": self.metadata.total_cost
            },
            "data": self.data,
            "review_status": self.review_status,
            "usage_stats": {
                "view_count": self.view_count,
                "usage_count": self.usage_count,
                "avg_rating": self.avg_rating
            }
        }
```

---

## 5. 课程规划数据模型 ⭐新增

### 5.1 Curriculum (课程体系)

整个学期或学年的课程规划。

```python
from typing import List, Dict, Optional, Any
from dataclasses import dataclass, field
from datetime import datetime, date
from enum import Enum

class CurriculumLevel(str, Enum):
    """课程级别"""
    SEMESTER = "学期"  # 一个学期
    ACADEMIC_YEAR = "学年"  # 一整学年
    TEXTBOOK = "教材"  # 整本教材

class PlanningStatus(str, Enum):
    """规划状态"""
    DRAFT = "草稿"
    PLANNING = "规划中"
    REVIEW = "审核中"
    APPROVED = "已批准"
    IN_PROGRESS = "进行中"
    COMPLETED = "已完成"
    ARCHIVED = "已归档"

@dataclass
class Curriculum:
    """课程体系数据模型"""
    
    # === 基础信息 ===
    curriculum_id: str  # 课程ID
    name: str  # 课程名称，例如："九年级化学上册"
    subject: str  # 学科
    grade: int  # 年级 (1-12)
    region: str  # 地区 (CN/US/UK等)
    
    level: CurriculumLevel = CurriculumLevel.SEMESTER
    textbook_version: Optional[str] = None  # 教材版本，例如："人教版"
    
    # === 时间规划 ===
    start_date: date  # 开始日期
    end_date: date  # 结束日期
    total_weeks: int  # 总周数
    lessons_per_week: int  # 每周课时数
    total_lessons: int  # 总课时数
    
    # === 课程结构 ===
    course_units: List[str] = field(default_factory=list)  # 单元ID列表
    # 单元之间的依赖关系
    unit_dependencies: Dict[str, List[str]] = field(default_factory=dict)
    # 例如: {"unit_3": ["unit_1", "unit_2"]}  # unit_3依赖unit_1和unit_2
    
    # === 故事化设计 ===
    story_universe: Optional[str] = None  # 故事宇宙ID（如果采用故事化教学）
    narrative_structure: Optional[str] = None  # 叙事结构类型
    # 例如: "hero_journey" | "mystery_solving" | "exploration"
    
    # === 知识体系 ===
    knowledge_graph_id: Optional[str] = None  # 知识图谱ID
    core_competencies: List[str] = field(default_factory=list)  # 核心素养目标
    curriculum_alignment: Dict[str, str] = field(default_factory=dict)  # 课标对齐
    
    # === 元数据 ===
    status: PlanningStatus = PlanningStatus.DRAFT
    created_by: str = ""  # 创建者
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    
    # === 统计数据 ===
    progress: Dict[str, Any] = field(default_factory=dict)
    # progress 示例:
    # {
    #     "completed_lessons": 15,
    #     "total_lessons": 40,
    #     "completion_rate": 0.375,
    #     "current_unit": "unit_2",
    #     "on_schedule": True
    # }
    
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def get_total_duration_weeks(self) -> int:
        """计算总周数"""
        return (self.end_date - self.start_date).days // 7
    
    def get_completion_rate(self) -> float:
        """获取完成率"""
        return self.progress.get('completion_rate', 0.0)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "curriculum_id": self.curriculum_id,
            "name": self.name,
            "subject": self.subject,
            "grade": self.grade,
            "region": self.region,
            "level": self.level,
            "textbook_version": self.textbook_version,
            "start_date": self.start_date.isoformat(),
            "end_date": self.end_date.isoformat(),
            "total_weeks": self.total_weeks,
            "lessons_per_week": self.lessons_per_week,
            "total_lessons": self.total_lessons,
            "course_units": self.course_units,
            "unit_dependencies": self.unit_dependencies,
            "story_universe": self.story_universe,
            "narrative_structure": self.narrative_structure,
            "knowledge_graph_id": self.knowledge_graph_id,
            "core_competencies": self.core_competencies,
            "curriculum_alignment": self.curriculum_alignment,
            "status": self.status,
            "created_by": self.created_by,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "progress": self.progress,
            "metadata": self.metadata
        }
```

### 5.2 CourseUnit (课程单元)

课程的一个单元或章节。

```python
@dataclass
class CourseUnit:
    """课程单元数据模型"""
    
    # === 基础信息 ===
    unit_id: str  # 单元ID
    curriculum_id: str  # 所属课程ID
    
    sequence_number: int  # 单元序号（第几单元）
    title: str  # 单元标题
    subtitle: Optional[str] = None  # 副标题
    
    # === 教学内容 ===
    topics: List[str] = field(default_factory=list)  # 主题/知识点列表
    learning_objectives: List[str] = field(default_factory=list)  # 学习目标
    
    # 核心问题（驱动问题）
    driving_questions: List[str] = field(default_factory=list)
    
    # === 课时规划 ===
    total_lessons: int = 0  # 本单元课时数
    lesson_metas: List[str] = field(default_factory=list)  # 课时元信息ID列表
    
    # 课时分配策略
    lesson_distribution: Dict[str, int] = field(default_factory=dict)
    # 例如: {
    #     "introduction": 1,      # 导入课
    #     "new_content": 3,       # 新授课
    #     "practice": 2,          # 练习课
    #     "review": 1,            # 复习课
    #     "assessment": 1         # 评估课
    # }
    
    # === 时间安排 ===
    estimated_weeks: float = 0.0  # 预计周数
    start_week: Optional[int] = None  # 开始周
    end_week: Optional[int] = None  # 结束周
    
    # === 故事化元素 ===
    story_arc_id: Optional[str] = None  # 故事弧ID
    narrative_theme: Optional[str] = None  # 叙事主题
    
    # === 依赖关系 ===
    prerequisite_units: List[str] = field(default_factory=list)  # 前置单元
    prerequisite_knowledge: List[str] = field(default_factory=list)  # 前置知识点
    
    # === 资源 ===
    resources: List[Dict[str, str]] = field(default_factory=list)
    # resources 示例:
    # [
    #     {"type": "video", "url": "...", "title": "..."},
    #     {"type": "document", "url": "...", "title": "..."}
    # ]
    
    # === 元数据 ===
    status: str = "planned"  # planned | in_progress | completed
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "unit_id": self.unit_id,
            "curriculum_id": self.curriculum_id,
            "sequence_number": self.sequence_number,
            "title": self.title,
            "subtitle": self.subtitle,
            "topics": self.topics,
            "learning_objectives": self.learning_objectives,
            "driving_questions": self.driving_questions,
            "total_lessons": self.total_lessons,
            "lesson_metas": self.lesson_metas,
            "lesson_distribution": self.lesson_distribution,
            "estimated_weeks": self.estimated_weeks,
            "start_week": self.start_week,
            "end_week": self.end_week,
            "story_arc_id": self.story_arc_id,
            "narrative_theme": self.narrative_theme,
            "prerequisite_units": self.prerequisite_units,
            "prerequisite_knowledge": self.prerequisite_knowledge,
            "resources": self.resources,
            "status": self.status,
            "metadata": self.metadata
        }
```

### 5.3 LessonMeta (课时元信息)

单个课时的规划信息（元数据）。

```python
@dataclass
class LessonMeta:
    """课时元信息 - 课程规划阶段的抽象描述"""
    
    # === 基础信息 ===
    lesson_meta_id: str  # 课时元信息ID
    unit_id: str  # 所属单元ID
    curriculum_id: str  # 所属课程ID
    
    sequence_number: int  # 课时序号（本单元第几课）
    global_sequence: int  # 全局序号（整个学期第几课）
    
    # === 课时定位 ===
    title: str  # 课时标题
    lesson_type: str  # 课型
    # 例如: "新授课" | "练习课" | "实验课" | "复习课" | "评估课"
    
    # === 教学内容 ===
    topics: List[str] = field(default_factory=list)  # 本课知识点
    objectives: List[str] = field(default_factory=list)  # 教学目标
    key_points: List[str] = field(default_factory=list)  # 重点
    difficult_points: List[str] = field(default_factory=list)  # 难点
    
    # === 时间规划 ===
    duration: int = 45  # 时长（分钟）
    scheduled_week: Optional[int] = None  # 计划周次
    scheduled_date: Optional[date] = None  # 计划日期
    
    # === 故事化元素 ===
    story_episode_id: Optional[str] = None  # 故事片段ID
    narrative_hook: Optional[str] = None  # 叙事钩子
    
    # === 依赖关系 ===
    prerequisite_lessons: List[str] = field(default_factory=list)  # 前置课时
    
    # === 生成配置 ===
    workflow_template_id: Optional[str] = None  # 工作流模板ID
    generation_config: Dict[str, Any] = field(default_factory=dict)
    # generation_config 示例:
    # {
    #     "pedagogy_model": "inquiry",
    #     "interaction_level": "high",
    #     "enable_story": True,
    #     "media_preferences": ["video", "interactive"]
    # }
    
    # === 生成状态 ===
    generated: bool = False  # 是否已生成完整内容
    lesson_script_id: Optional[str] = None  # 生成的课程脚本ID
    generation_quality_score: Optional[float] = None  # 生成质量分数
    
    # === 元数据 ===
    status: str = "planned"  # planned | generating | generated | reviewed | approved
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "lesson_meta_id": self.lesson_meta_id,
            "unit_id": self.unit_id,
            "curriculum_id": self.curriculum_id,
            "sequence_number": self.sequence_number,
            "global_sequence": self.global_sequence,
            "title": self.title,
            "lesson_type": self.lesson_type,
            "topics": self.topics,
            "objectives": self.objectives,
            "key_points": self.key_points,
            "difficult_points": self.difficult_points,
            "duration": self.duration,
            "scheduled_week": self.scheduled_week,
            "scheduled_date": self.scheduled_date.isoformat() if self.scheduled_date else None,
            "story_episode_id": self.story_episode_id,
            "narrative_hook": self.narrative_hook,
            "prerequisite_lessons": self.prerequisite_lessons,
            "workflow_template_id": self.workflow_template_id,
            "generation_config": self.generation_config,
            "generated": self.generated,
            "lesson_script_id": self.lesson_script_id,
            "generation_quality_score": self.generation_quality_score,
            "status": self.status,
            "metadata": self.metadata
        }
```

### 5.4 StoryUniverse (故事宇宙)

整个课程的故事化设计框架。

```python
@dataclass
class Character:
    """故事角色"""
    character_id: str
    name: str  # 角色名
    role: str  # 角色类型：protagonist/mentor/antagonist/sidekick
    description: str  # 角色描述
    personality: List[str]  # 性格特征
    avatar_url: Optional[str] = None  # 头像URL

@dataclass
class StoryArc:
    """故事弧 - 一个单元的故事线"""
    arc_id: str
    unit_id: str  # 对应的课程单元ID
    
    arc_title: str  # 故事弧标题
    arc_type: str  # 故事类型：adventure/mystery/challenge/discovery
    
    # 故事结构
    setup: str  # 起：故事背景和开端
    conflict: str  # 承：冲突和挑战
    climax: str  # 转：高潮
    resolution: str  # 合：解决和总结
    
    # 涉及的角色
    characters: List[str] = field(default_factory=list)  # 角色ID列表
    
    # 故事片段（对应各个课时）
    episodes: List[str] = field(default_factory=list)  # 片段ID列表

@dataclass
class StoryEpisode:
    """故事片段 - 一个课时的故事内容"""
    episode_id: str
    arc_id: str  # 所属故事弧
    lesson_meta_id: str  # 对应的课时元信息ID
    
    episode_number: int  # 片段序号
    title: str  # 片段标题
    
    # 故事内容
    hook: str  # 开头钩子（吸引注意）
    development: str  # 发展（引入知识点）
    connection: str  # 连接（与下一课时的联系）
    
    # 知识点融入
    knowledge_integration: Dict[str, str] = field(default_factory=dict)
    # 例如: {"化学反应": "通过角色做实验来展示"}

@dataclass
class StoryUniverse:
    """故事宇宙数据模型"""
    
    # === 基础信息 ===
    universe_id: str
    curriculum_id: str  # 所属课程ID
    
    name: str  # 故事宇宙名称
    description: str  # 描述
    
    # === 世界观设定 ===
    world_setting: str  # 世界设定
    # 例如: "未来科技世界" | "魔法学院" | "历史穿越" | "侦探事务所"
    
    theme: str  # 主题
    tone: str  # 基调：humorous/serious/adventurous/mysterious
    
    # === 角色体系 ===
    characters: List[Character] = field(default_factory=list)
    main_character_id: str = ""  # 主角ID
    
    # === 故事结构 ===
    story_arcs: List[StoryArc] = field(default_factory=list)
    episodes: List[StoryEpisode] = field(default_factory=list)
    
    # === 跨单元连续性 ===
    continuity_elements: List[str] = field(default_factory=list)
    # 例如: ["持续的目标", "重复出现的挑战", "累积的成就系统"]
    
    # === 可视化资源 ===
    visual_assets: Dict[str, str] = field(default_factory=dict)
    # 例如: {
    #     "world_map": "url...",
    #     "character_designs": "url...",
    #     "location_images": "url..."
    # }
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def get_arc_by_unit(self, unit_id: str) -> Optional[StoryArc]:
        """根据单元ID获取故事弧"""
        for arc in self.story_arcs:
            if arc.unit_id == unit_id:
                return arc
        return None
    
    def get_episode_by_lesson(self, lesson_meta_id: str) -> Optional[StoryEpisode]:
        """根据课时ID获取故事片段"""
        for episode in self.episodes:
            if episode.lesson_meta_id == lesson_meta_id:
                return episode
        return None
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "universe_id": self.universe_id,
            "curriculum_id": self.curriculum_id,
            "name": self.name,
            "description": self.description,
            "world_setting": self.world_setting,
            "theme": self.theme,
            "tone": self.tone,
            "characters": [
                {
                    "character_id": c.character_id,
                    "name": c.name,
                    "role": c.role,
                    "description": c.description,
                    "personality": c.personality,
                    "avatar_url": c.avatar_url
                }
                for c in self.characters
            ],
            "main_character_id": self.main_character_id,
            "story_arcs": [
                {
                    "arc_id": arc.arc_id,
                    "unit_id": arc.unit_id,
                    "arc_title": arc.arc_title,
                    "arc_type": arc.arc_type,
                    "setup": arc.setup,
                    "conflict": arc.conflict,
                    "climax": arc.climax,
                    "resolution": arc.resolution,
                    "characters": arc.characters,
                    "episodes": arc.episodes
                }
                for arc in self.story_arcs
            ],
            "episodes": [
                {
                    "episode_id": ep.episode_id,
                    "arc_id": ep.arc_id,
                    "lesson_meta_id": ep.lesson_meta_id,
                    "episode_number": ep.episode_number,
                    "title": ep.title,
                    "hook": ep.hook,
                    "development": ep.development,
                    "connection": ep.connection,
                    "knowledge_integration": ep.knowledge_integration
                }
                for ep in self.episodes
            ],
            "continuity_elements": self.continuity_elements,
            "visual_assets": self.visual_assets,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "metadata": self.metadata
        }
```

### 5.5 KnowledgeGraph (知识图谱)

课程的知识点依赖关系图谱。

```python
@dataclass
class KnowledgeNode:
    """知识节点"""
    node_id: str
    name: str  # 知识点名称
    description: str  # 描述
    
    # 分类
    category: str  # concept/skill/fact/procedure
    difficulty_level: int  # 难度等级 1-5
    bloom_level: str  # 布鲁姆层级
    
    # 关联课时
    related_lessons: List[str] = field(default_factory=list)  # LessonMeta ID列表
    
    # 估计学习时间
    estimated_learning_time: int = 0  # 分钟

@dataclass
class KnowledgeEdge:
    """知识边 - 知识点之间的关系"""
    edge_id: str
    source_node_id: str  # 源知识点
    target_node_id: str  # 目标知识点
    
    # 关系类型
    relationship_type: str
    # 例如: "prerequisite" | "supports" | "extends" | "contradicts" | "parallels"
    
    # 关系强度
    strength: float = 1.0  # 0.0-1.0
    
    # 说明
    description: Optional[str] = None

@dataclass
class KnowledgeGraph:
    """知识图谱数据模型"""
    
    graph_id: str
    curriculum_id: str
    
    name: str
    description: str
    
    # 图谱结构
    nodes: List[KnowledgeNode] = field(default_factory=list)
    edges: List[KnowledgeEdge] = field(default_factory=list)
    
    # 分层信息（用于拓扑排序）
    levels: Dict[str, int] = field(default_factory=dict)
    # 例如: {"node_1": 1, "node_2": 1, "node_3": 2, ...}
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def get_prerequisites(self, node_id: str) -> List[str]:
        """获取某知识点的所有前置知识点"""
        prereqs = []
        for edge in self.edges:
            if edge.target_node_id == node_id and edge.relationship_type == "prerequisite":
                prereqs.append(edge.source_node_id)
        return prereqs
    
    def get_learning_path(self, target_node_id: str) -> List[str]:
        """获取学习路径（从基础到目标知识点）"""
        # 简化实现：基于拓扑排序
        path = []
        visited = set()
        
        def dfs(node_id: str):
            if node_id in visited:
                return
            visited.add(node_id)
            
            # 先访问所有前置
            for prereq in self.get_prerequisites(node_id):
                dfs(prereq)
            
            path.append(node_id)
        
        dfs(target_node_id)
        return path
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "graph_id": self.graph_id,
            "curriculum_id": self.curriculum_id,
            "name": self.name,
            "description": self.description,
            "nodes": [
                {
                    "node_id": n.node_id,
                    "name": n.name,
                    "description": n.description,
                    "category": n.category,
                    "difficulty_level": n.difficulty_level,
                    "bloom_level": n.bloom_level,
                    "related_lessons": n.related_lessons,
                    "estimated_learning_time": n.estimated_learning_time
                }
                for n in self.nodes
            ],
            "edges": [
                {
                    "edge_id": e.edge_id,
                    "source_node_id": e.source_node_id,
                    "target_node_id": e.target_node_id,
                    "relationship_type": e.relationship_type,
                    "strength": e.strength,
                    "description": e.description
                }
                for e in self.edges
            ],
            "levels": self.levels,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "metadata": self.metadata
        }
```

---

## 6. 模板本地化数据模型 ⭐新增

### 6.1 LocalizationConfig (本地化配置)

模板本地化的配置和任务信息。

```python
from typing import List, Dict, Optional, Any
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum

class LocalizationStatus(str, Enum):
    """本地化状态"""
    PENDING = "待开始"
    MAPPING = "标准映射中"
    TRANSLATING = "翻译中"
    ADAPTING = "文化适配中"
    REVIEWING = "审核中"
    TESTING = "测试中"
    COMPLETED = "已完成"
    FAILED = "失败"

@dataclass
class LocalizationConfig:
    """本地化配置数据模型"""
    
    # === 基础信息 ===
    localization_id: str
    job_name: str  # 任务名称
    
    # === 源和目标 ===
    source_template_id: str  # 源模板ID
    source_region: str  # 源地区（例如: "CN"）
    source_language: str  # 源语言（例如: "zh-CN"）
    
    target_region: str  # 目标地区（例如: "US"）
    target_language: str  # 目标语言（例如: "en-US"）
    
    # === 本地化策略 ===
    localization_strategy: Dict[str, Any] = field(default_factory=dict)
    # localization_strategy 示例:
    # {
    #     "translation_model": "gpt-4",
    #     "preserve_structure": True,
    #     "adaptation_level": "high",  # low/medium/high
    #     "compliance_check": True,
    #     "human_review_required": True
    # }
    
    # === 适配范围 ===
    adaptation_scopes: List[str] = field(default_factory=list)
    # 例如: ["standards", "culture", "pedagogy", "examples", "measurement_units"]
    
    # === 标准映射 ===
    standard_mappings: List[str] = field(default_factory=list)  # StandardMapping ID列表
    
    # === 文化适配记录 ===
    cultural_adaptations: List[Dict[str, Any]] = field(default_factory=list)
    # cultural_adaptations 示例:
    # [
    #     {
    #         "source_element": "春节",
    #         "target_element": "Christmas",
    #         "category": "cultural_reference",
    #         "reason": "更贴近目标文化背景"
    #     }
    # ]
    
    # === 教学法调整 ===
    pedagogy_adjustments: List[Dict[str, Any]] = field(default_factory=list)
    # pedagogy_adjustments 示例:
    # [
    #     {
    #         "aspect": "classroom_structure",
    #         "source_value": "teacher-centered",
    #         "target_value": "student-centered",
    #         "reason": "符合美国教学习惯"
    #     }
    # ]
    
    # === 合规性检查 ===
    compliance_issues: List[str] = field(default_factory=list)  # ComplianceIssue ID列表
    
    # === 质量控制 ===
    quality_report_id: Optional[str] = None
    automation_rate: float = 0.0  # 自动化率 (0.0-1.0)
    human_review_points: List[str] = field(default_factory=list)  # 需要人工审核的点
    
    # === 执行状态 ===
    status: LocalizationStatus = LocalizationStatus.PENDING
    current_phase: Optional[str] = None  # 当前阶段
    progress: float = 0.0  # 进度 (0.0-1.0)
    
    # === 输出 ===
    target_template_id: Optional[str] = None  # 生成的目标模板ID
    
    # === 元数据 ===
    created_by: str = ""
    created_at: datetime = field(default_factory=datetime.utcnow)
    completed_at: Optional[datetime] = None
    
    estimated_duration: Optional[int] = None  # 预计耗时（小时）
    actual_duration: Optional[int] = None  # 实际耗时（小时）
    
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "localization_id": self.localization_id,
            "job_name": self.job_name,
            "source_template_id": self.source_template_id,
            "source_region": self.source_region,
            "source_language": self.source_language,
            "target_region": self.target_region,
            "target_language": self.target_language,
            "localization_strategy": self.localization_strategy,
            "adaptation_scopes": self.adaptation_scopes,
            "standard_mappings": self.standard_mappings,
            "cultural_adaptations": self.cultural_adaptations,
            "pedagogy_adjustments": self.pedagogy_adjustments,
            "compliance_issues": self.compliance_issues,
            "quality_report_id": self.quality_report_id,
            "automation_rate": self.automation_rate,
            "human_review_points": self.human_review_points,
            "status": self.status,
            "current_phase": self.current_phase,
            "progress": self.progress,
            "target_template_id": self.target_template_id,
            "created_by": self.created_by,
            "created_at": self.created_at.isoformat(),
            "completed_at": self.completed_at.isoformat() if self.completed_at else None,
            "estimated_duration": self.estimated_duration,
            "actual_duration": self.actual_duration,
            "metadata": self.metadata
        }
```

### 6.2 StandardMapping (课标映射)

源地区课标到目标地区课标的映射关系。

```python
@dataclass
class StandardMapping:
    """课标映射数据模型"""
    
    mapping_id: str
    localization_id: str  # 所属本地化任务ID
    
    # === 源课标 ===
    source_standard_id: str  # 源课标ID
    source_standard_name: str  # 源课标名称
    source_description: str  # 源课标描述
    source_region: str  # 源地区
    
    # === 目标课标 ===
    target_standard_id: str  # 目标课标ID
    target_standard_name: str  # 目标课标名称
    target_description: str  # 目标课标描述
    target_region: str  # 目标地区
    
    # === 映射关系 ===
    mapping_type: str  # exact/partial/approximate/no_match
    similarity_score: float = 0.0  # 相似度分数 (0.0-1.0)
    
    # === 差异说明 ===
    differences: List[str] = field(default_factory=list)
    # 例如: ["源课标更强调计算，目标课标更强调概念理解"]
    
    # === 调整建议 ===
    adaptation_suggestions: List[str] = field(default_factory=list)
    # 例如: ["增加概念解释环节", "减少机械计算练习"]
    
    # === 元数据 ===
    confidence: float = 0.0  # 映射置信度 (0.0-1.0)
    verified_by_human: bool = False  # 是否经过人工验证
    created_at: datetime = field(default_factory=datetime.utcnow)
    
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "mapping_id": self.mapping_id,
            "localization_id": self.localization_id,
            "source_standard_id": self.source_standard_id,
            "source_standard_name": self.source_standard_name,
            "source_description": self.source_description,
            "source_region": self.source_region,
            "target_standard_id": self.target_standard_id,
            "target_standard_name": self.target_standard_name,
            "target_description": self.target_description,
            "target_region": self.target_region,
            "mapping_type": self.mapping_type,
            "similarity_score": self.similarity_score,
            "differences": self.differences,
            "adaptation_suggestions": self.adaptation_suggestions,
            "confidence": self.confidence,
            "verified_by_human": self.verified_by_human,
            "created_at": self.created_at.isoformat(),
            "metadata": self.metadata
        }
```

### 6.3 ComplianceIssue (合规问题)

本地化过程中发现的合规性问题。

```python
@dataclass
class ComplianceIssue:
    """合规性问题数据模型"""
    
    issue_id: str
    localization_id: str  # 所属本地化任务ID
    
    # === 问题分类 ===
    category: str  # legal/cultural/educational/safety/privacy
    severity: str  # critical/high/medium/low
    
    # === 问题描述 ===
    title: str  # 问题标题
    description: str  # 详细描述
    
    # 问题位置
    location: Dict[str, Any] = field(default_factory=dict)
    # location 示例:
    # {
    #     "node_id": "script_generator",
    #     "section": "introduction",
    #     "line_number": 45,
    #     "original_text": "..."
    # }
    
    # === 违规内容 ===
    flagged_content: str  # 被标记的内容
    violated_rule: str  # 违反的规则
    rule_reference: Optional[str] = None  # 规则引用链接
    
    # === 建议修改 ===
    suggested_fix: Optional[str] = None  # 建议的修复方案
    alternative_content: Optional[str] = None  # 替代内容
    
    # === 处理状态 ===
    status: str = "open"  # open/fixed/wont_fix/false_positive
    resolution: Optional[str] = None  # 解决方案说明
    resolved_by: Optional[str] = None  # 解决者
    resolved_at: Optional[datetime] = None  # 解决时间
    
    # === 元数据 ===
    detected_at: datetime = field(default_factory=datetime.utcnow)
    detector: str = ""  # 检测器名称（auto/human）
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "issue_id": self.issue_id,
            "localization_id": self.localization_id,
            "category": self.category,
            "severity": self.severity,
            "title": self.title,
            "description": self.description,
            "location": self.location,
            "flagged_content": self.flagged_content,
            "violated_rule": self.violated_rule,
            "rule_reference": self.rule_reference,
            "suggested_fix": self.suggested_fix,
            "alternative_content": self.alternative_content,
            "status": self.status,
            "resolution": self.resolution,
            "resolved_by": self.resolved_by,
            "resolved_at": self.resolved_at.isoformat() if self.resolved_at else None,
            "detected_at": self.detected_at.isoformat(),
            "detector": self.detector,
            "metadata": self.metadata
        }
```

### 6.4 QualityReport (质量报告)

本地化质量评估报告。

```python
@dataclass
class QualityMetric:
    """质量指标"""
    metric_name: str
    score: float  # 0.0-10.0
    weight: float  # 权重
    details: Optional[str] = None

@dataclass
class QualityReport:
    """质量报告数据模型"""
    
    report_id: str
    localization_id: str
    
    # === 总体评分 ===
    overall_score: float = 0.0  # 0.0-10.0
    grade: str = ""  # A/B/C/D/F
    
    # === 分项评分 ===
    metrics: List[QualityMetric] = field(default_factory=list)
    # metrics 示例:
    # [
    #     QualityMetric("translation_accuracy", 9.2, 0.25),
    #     QualityMetric("cultural_appropriateness", 8.5, 0.20),
    #     QualityMetric("pedagogy_alignment", 8.8, 0.20),
    #     QualityMetric("standard_compliance", 9.5, 0.20),
    #     QualityMetric("consistency", 9.0, 0.15)
    # ]
    
    # === 问题统计 ===
    issue_summary: Dict[str, int] = field(default_factory=dict)
    # issue_summary 示例:
    # {
    #     "critical": 0,
    #     "high": 2,
    #     "medium": 5,
    #     "low": 8
    # }
    
    # === 优点与改进 ===
    strengths: List[str] = field(default_factory=list)
    weaknesses: List[str] = field(default_factory=list)
    recommendations: List[str] = field(default_factory=list)
    
    # === 人工审核 ===
    human_review_required: bool = False
    human_review_points: List[str] = field(default_factory=list)
    
    # === 测试结果 ===
    test_results: Optional[Dict[str, Any]] = None
    # test_results 示例:
    # {
    #     "pilot_users": 10,
    #     "avg_rating": 4.5,
    #     "completion_rate": 0.95,
    #     "feedback": [...]
    # }
    
    # === 元数据 ===
    generated_at: datetime = field(default_factory=datetime.utcnow)
    generated_by: str = ""  # auto/human
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def calculate_overall_score(self):
        """计算加权总分"""
        if not self.metrics:
            return 0.0
        
        total_score = sum(m.score * m.weight for m in self.metrics)
        total_weight = sum(m.weight for m in self.metrics)
        
        self.overall_score = total_score / total_weight if total_weight > 0 else 0.0
        
        # 评级
        if self.overall_score >= 9.0:
            self.grade = "A"
        elif self.overall_score >= 8.0:
            self.grade = "B"
        elif self.overall_score >= 7.0:
            self.grade = "C"
        elif self.overall_score >= 6.0:
            self.grade = "D"
        else:
            self.grade = "F"
        
        return self.overall_score
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "report_id": self.report_id,
            "localization_id": self.localization_id,
            "overall_score": self.overall_score,
            "grade": self.grade,
            "metrics": [
                {
                    "metric_name": m.metric_name,
                    "score": m.score,
                    "weight": m.weight,
                    "details": m.details
                }
                for m in self.metrics
            ],
            "issue_summary": self.issue_summary,
            "strengths": self.strengths,
            "weaknesses": self.weaknesses,
            "recommendations": self.recommendations,
            "human_review_required": self.human_review_required,
            "human_review_points": self.human_review_points,
            "test_results": self.test_results,
            "generated_at": self.generated_at.isoformat(),
            "generated_by": self.generated_by,
            "metadata": self.metadata
        }
```

---

## 7. 组件扩展数据模型 ⭐新增

### 7.1 NodeMetadata (节点元数据)

工作流节点的元数据定义。

```python
from typing import List, Dict, Optional, Any
from dataclasses import dataclass, field
from datetime import datetime

@dataclass
class NodeMetadata:
    """节点元数据"""
    
    # === 基础信息 ===
    node_id: str  # 节点唯一标识
    node_name: str  # 节点名称
    version: str  # 版本号 (语义化版本，例如: "1.2.0")
    
    # === 分类 ===
    category: str  # 分类
    # 例如: "content_generation" | "quality_assurance" | "interaction_design" | "assessment"
    
    subcategory: Optional[str] = None  # 子分类
    tags: List[str] = field(default_factory=list)  # 标签
    
    # === 描述 ===
    description: str = ""  # 功能描述
    long_description: Optional[str] = None  # 详细描述
    
    # === 适用范围 ===
    applicable_subjects: List[str] = field(default_factory=list)  # 适用学科
    applicable_grades: str = ""  # 适用年级，例如: "1-12" 或 "7-9"
    applicable_lesson_types: List[str] = field(default_factory=list)  # 适用课型
    
    # === 作者信息 ===
    author: str = ""  # 作者
    author_email: Optional[str] = None  # 作者邮箱
    organization: Optional[str] = None  # 组织/公司
    license: str = "MIT"  # 许可证
    
    # === 依赖项 ===
    dependencies: List[str] = field(default_factory=list)  # 依赖的包或服务
    # 例如: ["llm_client >= 2.0.0", "unity_renderer >= 3.5.0"]
    
    # === 性能指标 ===
    estimated_duration: int = 10  # 预计执行时间（秒）
    memory_limit: int = 512  # 内存限制（MB）
    cpu_limit: float = 1.0  # CPU限制（核数）
    
    # === 输入输出Schema ===
    input_schema: Dict[str, Any] = field(default_factory=dict)
    output_schema: Dict[str, Any] = field(default_factory=dict)
    
    # === 配置Schema ===
    config_schema: Dict[str, Any] = field(default_factory=dict)
    
    # === 文档 ===
    documentation_url: Optional[str] = None  # 文档链接
    example_url: Optional[str] = None  # 示例链接
    
    # === 版本信息 ===
    changelog: List[Dict[str, str]] = field(default_factory=list)
    # changelog 示例:
    # [
    #     {"version": "1.2.0", "date": "2025-12-09", "changes": "新增降级策略"},
    #     {"version": "1.1.0", "date": "2025-11-01", "changes": "性能优化"}
    # ]
    
    # === 统计数据 ===
    download_count: int = 0  # 下载次数
    usage_count: int = 0  # 使用次数
    avg_rating: float = 0.0  # 平均评分
    
    # === 认证 ===
    certification_level: Optional[str] = None  # bronze/silver/gold/platinum
    verified: bool = False  # 是否经过官方验证
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.utcnow)
    updated_at: datetime = field(default_factory=datetime.utcnow)
    deprecated: bool = False  # 是否已弃用
    deprecation_message: Optional[str] = None
    
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "node_id": self.node_id,
            "node_name": self.node_name,
            "version": self.version,
            "category": self.category,
            "subcategory": self.subcategory,
            "tags": self.tags,
            "description": self.description,
            "long_description": self.long_description,
            "applicable_subjects": self.applicable_subjects,
            "applicable_grades": self.applicable_grades,
            "applicable_lesson_types": self.applicable_lesson_types,
            "author": self.author,
            "author_email": self.author_email,
            "organization": self.organization,
            "license": self.license,
            "dependencies": self.dependencies,
            "estimated_duration": self.estimated_duration,
            "memory_limit": self.memory_limit,
            "cpu_limit": self.cpu_limit,
            "input_schema": self.input_schema,
            "output_schema": self.output_schema,
            "config_schema": self.config_schema,
            "documentation_url": self.documentation_url,
            "example_url": self.example_url,
            "changelog": self.changelog,
            "download_count": self.download_count,
            "usage_count": self.usage_count,
            "avg_rating": self.avg_rating,
            "certification_level": self.certification_level,
            "verified": self.verified,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "deprecated": self.deprecated,
            "deprecation_message": self.deprecation_message,
            "metadata": self.metadata
        }
```

### 7.2 ComponentPackage (组件包)

组件的打包信息。

```python
@dataclass
class ComponentPackage:
    """组件包数据模型"""
    
    package_id: str  # 包ID
    node_id: str  # 节点ID
    version: str  # 版本号
    
    # === 包信息 ===
    package_name: str  # 包名称
    display_name: str  # 显示名称
    
    # === 文件信息 ===
    files: Dict[str, str] = field(default_factory=dict)
    # files 示例:
    # {
    #     "main": "nodes/virtual_lab_node.py",
    #     "tests": "tests/test_virtual_lab_node.py",
    #     "docs": "docs/virtual_lab_node.md",
    #     "examples": "examples/virtual_lab_examples.yaml"
    # }
    
    # === 包大小 ===
    package_size: int = 0  # 字节
    
    # === 下载信息 ===
    download_url: Optional[str] = None
    checksum: Optional[str] = None  # MD5/SHA256校验和
    
    # === 定价 ===
    pricing_model: str = "free"  # free/freemium/paid
    free_quota: int = 0  # 免费调用次数
    price_per_call: float = 0.0  # 单次调用价格
    
    # === 元数据 ===
    published_at: datetime = field(default_factory=datetime.utcnow)
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "package_id": self.package_id,
            "node_id": self.node_id,
            "version": self.version,
            "package_name": self.package_name,
            "display_name": self.display_name,
            "files": self.files,
            "package_size": self.package_size,
            "download_url": self.download_url,
            "checksum": self.checksum,
            "pricing_model": self.pricing_model,
            "free_quota": self.free_quota,
            "price_per_call": self.price_per_call,
            "published_at": self.published_at.isoformat(),
            "metadata": self.metadata
        }
```

### 7.3 ComponentUsageRecord (组件使用记录)

组件的使用统计记录。

```python
@dataclass
class ComponentUsageRecord:
    """组件使用记录"""
    
    record_id: str
    node_id: str  # 节点ID
    version: str  # 使用的版本
    
    # === 使用信息 ===
    user_id: str  # 使用者ID
    workflow_instance_id: str  # 工作流实例ID
    
    # === 执行信息 ===
    executed_at: datetime  # 执行时间
    duration: float = 0.0  # 执行时长（秒）
    memory_used: float = 0.0  # 内存使用（MB）
    cpu_used: float = 0.0  # CPU使用（核·秒）
    
    # === 结果 ===
    status: str = "success"  # success/failure/timeout/error
    error_message: Optional[str] = None
    
    # === 降级信息 ===
    fallback_activated: bool = False
    fallback_reason: Optional[str] = None
    
    # === 质量评分 ===
    output_quality_score: Optional[float] = None  # 输出质量分数
    
    # === 用户反馈 ===
    user_rating: Optional[int] = None  # 1-5星
    user_comment: Optional[str] = None
    
    # === 计费 ===
    billable: bool = False  # 是否计费
    cost: float = 0.0  # 费用
    
    # === 元数据 ===
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "record_id": self.record_id,
            "node_id": self.node_id,
            "version": self.version,
            "user_id": self.user_id,
            "workflow_instance_id": self.workflow_instance_id,
            "executed_at": self.executed_at.isoformat(),
            "duration": self.duration,
            "memory_used": self.memory_used,
            "cpu_used": self.cpu_used,
            "status": self.status,
            "error_message": self.error_message,
            "fallback_activated": self.fallback_activated,
            "fallback_reason": self.fallback_reason,
            "output_quality_score": self.output_quality_score,
            "user_rating": self.user_rating,
            "user_comment": self.user_comment,
            "billable": self.billable,
            "cost": self.cost,
            "metadata": self.metadata
        }
```

---

## 8. V2.0新增：领域适配器数据模型 ⭐

### 8.1 概述

领域适配器是V2.0的核心创新，实现"孵化器-鸡"分离，支持多领域扩展。

**核心数据模型**：
1. `DomainType`（枚举）- 领域类型
2. `AdapterStatus`（枚举）- 适配器状态
3. `AdapterMetadata` - 适配器元数据
4. `DomainInput` - 领域特定输入
5. `DomainOutput` - 领域特定输出
6. `AdapterConfig` - 适配器配置
7. `AdapterRegistry` - 注册表（运行时数据结构）
8. `RoutingResult` - 路由结果

---

### 8.2 枚举类型

```python
from enum import Enum

class DomainType(Enum):
    """领域类型枚举"""
    K12 = "k12"  # K12教育
    ART_HISTORY = "art_history"  # 美术史
    BESTSELLER_INTERPRETATION = "bestseller_interpretation"  # 畅销书解读
    VOCATIONAL_TRAINING = "vocational_training"  # 职业培训
    FINANCE = "finance"  # 金融（预留）
    HEALTH = "health"  # 健康（预留）
    PROGRAMMING = "programming"  # 编程（预留）
    
    @property
    def display_name(self) -> str:
        """显示名称"""
        names = {
            "k12": "K12教育",
            "art_history": "美术史",
            "bestseller_interpretation": "畅销书解读",
            "vocational_training": "职业技能培训",
            "finance": "金融理财",
            "health": "健康养生",
            "programming": "编程开发"
        }
        return names.get(self.value, self.value)


class AdapterStatus(Enum):
    """适配器状态枚举"""
    INSTALLED = "installed"  # 已安装但未激活
    ACTIVE = "active"  # 已激活（可用）
    INACTIVE = "inactive"  # 已停用
    ERROR = "error"  # 错误状态
    UPDATING = "updating"  # 更新中
    DEPRECATED = "deprecated"  # 已弃用
```

---

### 8.3 核心数据类

```python
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Any
from datetime import datetime

@dataclass
class AdapterMetadata:
    """
    适配器元数据
    
    存储适配器的基本信息、依赖、配置等
    """
    # === 基本信息 ===
    adapter_id: str  # 唯一标识，如 "k12-adapter-v2"
    name: str  # 显示名称
    domain: DomainType  # 领域类型
    version: str  # 语义化版本号 (e.g., "2.0.1")
    author: str  # 作者
    author_email: Optional[str] = None
    description: str = ""  # 描述
    
    # === 依赖要求 ===
    requires_python: str = ">=3.9"  # Python版本要求
    dependencies: List[str] = field(default_factory=list)  # Python包依赖
    # 示例: ["numpy>=1.20", "pandas>=1.3"]
    
    # === 配置 ===
    config_schema: Dict[str, Any] = field(default_factory=dict)  # 配置JSON Schema
    default_config: Dict[str, Any] = field(default_factory=dict)  # 默认配置
    supported_languages: List[str] = field(default_factory=lambda: ["zh"])  # 支持的语言
    
    # === 状态 ===
    status: AdapterStatus = AdapterStatus.INSTALLED
    enabled: bool = True
    install_date: datetime = field(default_factory=datetime.now)
    last_updated: Optional[datetime] = None
    
    # === 元信息 ===
    homepage: Optional[str] = None  # 主页URL
    repository: Optional[str] = None  # 代码仓库URL
    documentation: Optional[str] = None  # 文档URL
    license: str = "MIT"  # 开源协议
    tags: List[str] = field(default_factory=list)  # 标签
    
    # === 统计信息 ===
    usage_count: int = 0  # 使用次数
    average_rating: float = 0.0  # 平均评分
    total_ratings: int = 0  # 评分总数
    
    def to_dict(self) -> Dict[str, Any]:
        """转换为字典"""
        return {
            "adapter_id": self.adapter_id,
            "name": self.name,
            "domain": self.domain.value,
            "version": self.version,
            "author": self.author,
            "author_email": self.author_email,
            "description": self.description,
            "requires_python": self.requires_python,
            "dependencies": self.dependencies,
            "config_schema": self.config_schema,
            "default_config": self.default_config,
            "supported_languages": self.supported_languages,
            "status": self.status.value,
            "enabled": self.enabled,
            "install_date": self.install_date.isoformat(),
            "last_updated": self.last_updated.isoformat() if self.last_updated else None,
            "homepage": self.homepage,
            "repository": self.repository,
            "documentation": self.documentation,
            "license": self.license,
            "tags": self.tags,
            "usage_count": self.usage_count,
            "average_rating": self.average_rating,
            "total_ratings": self.total_ratings
        }


@dataclass
class DomainInput:
    """
    领域特定输入
    
    由适配器的parse_input()方法生成
    """
    domain: DomainType  # 领域类型
    raw_input: str  # 原始用户输入
    parsed_data: Dict[str, Any]  # 解析后的结构化数据
    
    # K12示例: {"grade": 3, "subject": "数学", "topic": "除法"}
    # 美术史示例: {"artist": "莫奈", "artwork": "印象·日出", "movement": "印象派"}
    # 畅销书示例: {"book_title": "原则", "author": "Ray Dalio", "chapters": [1, 2, 3]}
    # 职业培训示例: {"skill": "Excel", "tool": "数据透视表", "level": "初级"}
    
    confidence: float = 1.0  # 解析置信度 (0-1)
    metadata: Dict[str, Any] = field(default_factory=dict)  # 额外元数据
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "domain": self.domain.value,
            "raw_input": self.raw_input,
            "parsed_data": self.parsed_data,
            "confidence": self.confidence,
            "metadata": self.metadata,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class DomainOutput:
    """
    领域特定输出
    
    由适配器的apply_domain_template()方法生成
    """
    domain: DomainType  # 领域类型
    content: Any  # 具体内容（格式由适配器定义）
    # K12: 教案dict
    # 美术史: 艺术鉴赏指南dict
    # 畅销书: 书籍解读dict
    # 职业培训: 培训课程dict
    
    format: str = "markdown"  # 输出格式 (markdown/html/pdf/json)
    template_version: str = "1.0"  # 模板版本
    
    # === 质量指标 ===
    quality_score: float = 0.0  # 整体质量分数 (0-1)
    completeness: float = 0.0  # 完整性 (0-1)
    accuracy: float = 0.0  # 准确性 (0-1)
    
    # === 元数据 ===
    metadata: Dict[str, Any] = field(default_factory=dict)
    generated_at: datetime = field(default_factory=datetime.now)
    generation_time: float = 0.0  # 生成耗时（秒）
    
    # === 资源引用 ===
    assets: List[str] = field(default_factory=list)  # 关联资源（图片、视频等）
    references: List[str] = field(default_factory=list)  # 参考资料链接
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "domain": self.domain.value,
            "content": self.content,
            "format": self.format,
            "template_version": self.template_version,
            "quality_score": self.quality_score,
            "completeness": self.completeness,
            "accuracy": self.accuracy,
            "metadata": self.metadata,
            "generated_at": self.generated_at.isoformat(),
            "generation_time": self.generation_time,
            "assets": self.assets,
            "references": self.references
        }


@dataclass
class AdapterConfig:
    """
    适配器配置
    
    存储适配器的运行时配置
    """
    adapter_id: str  # 适配器ID
    settings: Dict[str, Any] = field(default_factory=dict)  # 配置项
    # 示例:
    # {
    #     "max_objectives": 10,
    #     "preferred_narrative_strategy": "artist_biography",
    #     "content_sources": ["wikiart", "wikipedia"]
    # }
    
    enabled: bool = True  # 是否启用
    priority: int = 1  # 优先级（1最高，用于同领域多适配器）
    
    # === 性能配置 ===
    timeout: int = 30  # 超时时间（秒）
    max_retries: int = 3  # 最大重试次数
    cache_enabled: bool = True  # 是否启用缓存
    cache_ttl: int = 3600  # 缓存TTL（秒）
    
    # === 降级配置 ===
    fallback_enabled: bool = True  # 是否启用降级
    fallback_adapter_id: Optional[str] = None  # 降级适配器ID
    
    # === 监控配置 ===
    metrics_enabled: bool = True  # 是否启用指标收集
    logging_level: str = "INFO"  # 日志级别
    
    created_at: datetime = field(default_factory=datetime.now)
    updated_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "adapter_id": self.adapter_id,
            "settings": self.settings,
            "enabled": self.enabled,
            "priority": self.priority,
            "timeout": self.timeout,
            "max_retries": self.max_retries,
            "cache_enabled": self.cache_enabled,
            "cache_ttl": self.cache_ttl,
            "fallback_enabled": self.fallback_enabled,
            "fallback_adapter_id": self.fallback_adapter_id,
            "metrics_enabled": self.metrics_enabled,
            "logging_level": self.logging_level,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat()
        }


@dataclass
class RoutingResult:
    """
    路由结果
    
    记录适配器路由的决策过程
    """
    selected_adapter_id: str  # 选中的适配器ID
    domain: DomainType  # 领域
    
    # === 路由策略 ===
    strategy: str  # 路由策略 (explicit/keyword/llm/history/fallback)
    confidence: float  # 置信度 (0-1)
    
    # === 候选适配器 ===
    candidates: List[Dict[str, Any]] = field(default_factory=list)
    # 示例:
    # [
    #     {"adapter_id": "k12-adapter", "score": 0.9, "reason": "keyword match"},
    #     {"adapter_id": "art-history-adapter", "score": 0.3, "reason": "low match"}
    # ]
    
    # === 决策过程 ===
    decision_tree: Dict[str, Any] = field(default_factory=dict)
    # 记录每个策略的评分和选择原因
    
    # === 性能指标 ===
    routing_time: float = 0.0  # 路由耗时（秒）
    
    # === 元数据 ===
    timestamp: datetime = field(default_factory=datetime.now)
    user_id: Optional[str] = None
    request_id: Optional[str] = None
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "selected_adapter_id": self.selected_adapter_id,
            "domain": self.domain.value,
            "strategy": self.strategy,
            "confidence": self.confidence,
            "candidates": self.candidates,
            "decision_tree": self.decision_tree,
            "routing_time": self.routing_time,
            "timestamp": self.timestamp.isoformat(),
            "user_id": self.user_id,
            "request_id": self.request_id
        }
```

---

### 8.4 运行时数据结构

```python
@dataclass
class AdapterRegistryState:
    """
    适配器注册表状态（运行时）
    
    不持久化到数据库，仅存在于内存中
    """
    adapters: Dict[str, 'BaseDomainAdapter'] = field(default_factory=dict)
    # key: adapter_id, value: 适配器实例
    
    metadata: Dict[str, AdapterMetadata] = field(default_factory=dict)
    # key: adapter_id, value: 元数据
    
    domain_to_adapter: Dict[DomainType, str] = field(default_factory=dict)
    # key: domain, value: adapter_id (默认适配器)
    
    configs: Dict[str, AdapterConfig] = field(default_factory=dict)
    # key: adapter_id, value: 配置
    
    health_status: Dict[str, bool] = field(default_factory=dict)
    # key: adapter_id, value: 健康状态
    
    last_health_check: Dict[str, datetime] = field(default_factory=dict)
    # key: adapter_id, value: 上次健康检查时间
```

---

### 8.5 数据库表设计

```sql
-- 适配器元数据表
CREATE TABLE adapter_metadata (
    adapter_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    domain VARCHAR(50) NOT NULL,  -- k12/art_history/bestseller/vocational
    version VARCHAR(20) NOT NULL,
    
    -- 作者信息
    author VARCHAR(100) NOT NULL,
    author_email VARCHAR(200),
    description TEXT,
    
    -- 依赖
    requires_python VARCHAR(50) DEFAULT '>=3.9',
    dependencies JSONB DEFAULT '[]'::jsonb,  -- ["numpy>=1.20", ...]
    
    -- 配置
    config_schema JSONB DEFAULT '{}'::jsonb,
    default_config JSONB DEFAULT '{}'::jsonb,
    supported_languages JSONB DEFAULT '["zh"]'::jsonb,
    
    -- 状态
    status VARCHAR(20) DEFAULT 'installed',  -- installed/active/inactive/error/updating
    enabled BOOLEAN DEFAULT true,
    install_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP,
    
    -- 元信息
    homepage VARCHAR(500),
    repository VARCHAR(500),
    documentation VARCHAR(500),
    license VARCHAR(50) DEFAULT 'MIT',
    tags JSONB DEFAULT '[]'::jsonb,
    
    -- 统计
    usage_count INTEGER DEFAULT 0,
    average_rating FLOAT DEFAULT 0.0,
    total_ratings INTEGER DEFAULT 0,
    
    CONSTRAINT unique_adapter_version UNIQUE (adapter_id, version),
    CONSTRAINT check_domain CHECK (domain IN ('k12', 'art_history', 'bestseller_interpretation', 'vocational_training', 'finance', 'health', 'programming')),
    CONSTRAINT check_status CHECK (status IN ('installed', 'active', 'inactive', 'error', 'updating', 'deprecated'))
);

CREATE INDEX idx_adapter_domain ON adapter_metadata(domain);
CREATE INDEX idx_adapter_status ON adapter_metadata(status);
CREATE INDEX idx_adapter_enabled ON adapter_metadata(enabled);


-- 适配器配置表
CREATE TABLE adapter_configs (
    config_id SERIAL PRIMARY KEY,
    adapter_id VARCHAR(100) NOT NULL,
    
    -- 配置内容
    settings JSONB DEFAULT '{}'::jsonb,
    enabled BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 1,
    
    -- 性能配置
    timeout INTEGER DEFAULT 30,
    max_retries INTEGER DEFAULT 3,
    cache_enabled BOOLEAN DEFAULT true,
    cache_ttl INTEGER DEFAULT 3600,
    
    -- 降级配置
    fallback_enabled BOOLEAN DEFAULT true,
    fallback_adapter_id VARCHAR(100),
    
    -- 监控配置
    metrics_enabled BOOLEAN DEFAULT true,
    logging_level VARCHAR(10) DEFAULT 'INFO',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_config_adapter FOREIGN KEY (adapter_id) 
        REFERENCES adapter_metadata(adapter_id) ON DELETE CASCADE,
    CONSTRAINT check_logging_level CHECK (logging_level IN ('DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL'))
);

CREATE INDEX idx_config_adapter ON adapter_configs(adapter_id);


-- 适配器路由记录表
CREATE TABLE adapter_routing_records (
    record_id SERIAL PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL,
    
    -- 选中的适配器
    selected_adapter_id VARCHAR(100) NOT NULL,
    domain VARCHAR(50) NOT NULL,
    
    -- 路由策略
    strategy VARCHAR(50) NOT NULL,  -- explicit/keyword/llm/history/fallback
    confidence FLOAT NOT NULL,
    
    -- 候选适配器
    candidates JSONB DEFAULT '[]'::jsonb,
    decision_tree JSONB DEFAULT '{}'::jsonb,
    
    -- 性能
    routing_time FLOAT DEFAULT 0.0,
    
    -- 元数据
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(100),
    raw_input TEXT,
    
    CONSTRAINT fk_routing_adapter FOREIGN KEY (selected_adapter_id) 
        REFERENCES adapter_metadata(adapter_id) ON DELETE CASCADE
);

CREATE INDEX idx_routing_adapter ON adapter_routing_records(selected_adapter_id);
CREATE INDEX idx_routing_domain ON adapter_routing_records(domain);
CREATE INDEX idx_routing_strategy ON adapter_routing_records(strategy);
CREATE INDEX idx_routing_timestamp ON adapter_routing_records(timestamp);
CREATE INDEX idx_routing_user ON adapter_routing_records(user_id);


-- 适配器使用统计表
CREATE TABLE adapter_usage_stats (
    stat_id SERIAL PRIMARY KEY,
    adapter_id VARCHAR(100) NOT NULL,
    
    -- 时间维度
    stat_date DATE NOT NULL,
    stat_hour INTEGER,  -- 0-23，NULL表示日统计
    
    -- 使用统计
    total_requests INTEGER DEFAULT 0,
    successful_requests INTEGER DEFAULT 0,
    failed_requests INTEGER DEFAULT 0,
    timeout_requests INTEGER DEFAULT 0,
    
    -- 性能统计
    avg_response_time FLOAT DEFAULT 0.0,
    p50_response_time FLOAT DEFAULT 0.0,
    p95_response_time FLOAT DEFAULT 0.0,
    p99_response_time FLOAT DEFAULT 0.0,
    
    -- 质量统计
    avg_quality_score FLOAT DEFAULT 0.0,
    avg_user_rating FLOAT DEFAULT 0.0,
    total_ratings INTEGER DEFAULT 0,
    
    -- 降级统计
    fallback_count INTEGER DEFAULT 0,
    fallback_rate FLOAT DEFAULT 0.0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_stats_adapter FOREIGN KEY (adapter_id) 
        REFERENCES adapter_metadata(adapter_id) ON DELETE CASCADE,
    CONSTRAINT unique_adapter_stat UNIQUE (adapter_id, stat_date, stat_hour)
);

CREATE INDEX idx_stats_adapter_date ON adapter_usage_stats(adapter_id, stat_date);
CREATE INDEX idx_stats_date ON adapter_usage_stats(stat_date);
```

---

## 9. V2.0新增：AI目标生成器数据模型 ⭐

### 9.1 概述

AI目标生成器将模糊的用户需求转化为结构化的学习目标。

**核心数据模型**：
1. `UserRequest` - 用户原始请求
2. `AudienceProfile` - 受众画像
3. `LearningObjective` - 学习目标
4. `ObjectiveValidationResult` - 目标验证结果
5. `InteractiveFeedback` - 交互式反馈

---

### 9.2 核心数据类

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from datetime import datetime
from enum import Enum

class BloomLevel(Enum):
    """布鲁姆认知层级"""
    REMEMBER = 1  # 记忆
    UNDERSTAND = 2  # 理解
    APPLY = 3  # 应用
    ANALYZE = 4  # 分析
    EVALUATE = 5  # 评价
    CREATE = 6  # 创造
    
    @property
    def display_name(self) -> str:
        names = {
            1: "记忆",
            2: "理解",
            3: "应用",
            4: "分析",
            5: "评价",
            6: "创造"
        }
        return names[self.value]


@dataclass
class UserRequest:
    """
    用户原始请求
    
    从API接收的原始输入
    """
    request_id: str  # 请求ID
    raw_input: str  # 原始输入文本
    domain: Optional[DomainType] = None  # 用户指定的领域（可选）
    
    # === 用户信息 ===
    user_id: Optional[str] = None
    session_id: Optional[str] = None
    
    # === 上下文信息 ===
    language: str = "zh"  # 语言
    region: Optional[str] = None  # 地区
    
    # === 偏好设置 ===
    preferences: Dict[str, Any] = field(default_factory=dict)
    # 示例: {
    #   "max_objectives": 8, 
    #   "difficulty_level": "medium",
    #   "duration": 10,  # 课程时长（分钟），微课建议3-15分钟
    #   "enable_story": True
    # }
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.now)
    ip_address: Optional[str] = None
    user_agent: Optional[str] = None
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "request_id": self.request_id,
            "raw_input": self.raw_input,
            "domain": self.domain.value if self.domain else None,
            "user_id": self.user_id,
            "session_id": self.session_id,
            "language": self.language,
            "region": self.region,
            "preferences": self.preferences,
            "created_at": self.created_at.isoformat(),
            "ip_address": self.ip_address,
            "user_agent": self.user_agent
        }


@dataclass
class AudienceProfile:
    """
    受众画像
    
    AI分析或用户提供的受众特征
    """
    profile_id: str  # 画像ID
    
    # === 基本属性 ===
    age_range: str  # 年龄段 (e.g., "8-10", "18-25", "40+")
    education_level: str  # 教育水平 (e.g., "小学", "初中", "本科", "研究生")
    role: str  # 角色 (e.g., "student", "teacher", "parent", "professional")
    
    # === 兴趣和背景 ===
    interests: List[str] = field(default_factory=list)  # 兴趣标签
    prior_knowledge: List[str] = field(default_factory=list)  # 已有知识
    learning_goals: List[str] = field(default_factory=list)  # 学习目标
    
    # === 学习特征 ===
    learning_style: str = "visual"  # visual/auditory/kinesthetic/reading
    motivation_level: str = "medium"  # low/medium/high
    attention_span: int = 30  # 注意力时长（分钟）
    
    # === 挑战和需求 ===
    common_errors: List[str] = field(default_factory=list)  # 常见错误
    pain_points: List[str] = field(default_factory=list)  # 痛点
    special_needs: List[str] = field(default_factory=list)  # 特殊需求
    
    # === 元数据 ===
    confidence: float = 1.0  # AI推断置信度 (0-1)
    source: str = "user_provided"  # user_provided/ai_inferred
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "profile_id": self.profile_id,
            "age_range": self.age_range,
            "education_level": self.education_level,
            "role": self.role,
            "interests": self.interests,
            "prior_knowledge": self.prior_knowledge,
            "learning_goals": self.learning_goals,
            "learning_style": self.learning_style,
            "motivation_level": self.motivation_level,
            "attention_span": self.attention_span,
            "common_errors": self.common_errors,
            "pain_points": self.pain_points,
            "special_needs": self.special_needs,
            "confidence": self.confidence,
            "source": self.source,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class LearningObjective:
    """
    学习目标
    
    AI生成的结构化学习目标
    """
    objective_id: str  # 目标ID
    description: str  # 目标描述
    
    # === 分类 ===
    bloom_level: int  # 布鲁姆层级 (1-6)
    objective_type: str  # knowledge/skill/attitude
    domain: DomainType  # 领域
    
    # === 层级关系 ===
    parent_id: Optional[str] = None  # 父目标ID（用于目标树）
    order: int = 0  # 排序
    
    # === 评估标准 ===
    success_criteria: List[str] = field(default_factory=list)  # 成功标准
    assessment_methods: List[str] = field(default_factory=list)  # 评估方法
    
    # === 时间估算 ===
    estimated_time: int = 0  # 预估时长（分钟），根据课程类型：微课3-15分钟，常规课45-90分钟
    difficulty: str = "medium"  # easy/medium/hard
    
    # === 元数据 ===
    metadata: Dict[str, Any] = field(default_factory=dict)
    # 领域特定信息，如 K12: {"grade": 3, "subject": "数学"}
    
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "objective_id": self.objective_id,
            "description": self.description,
            "bloom_level": self.bloom_level,
            "objective_type": self.objective_type,
            "domain": self.domain.value,
            "parent_id": self.parent_id,
            "order": self.order,
            "success_criteria": self.success_criteria,
            "assessment_methods": self.assessment_methods,
            "estimated_time": self.estimated_time,
            "difficulty": self.difficulty,
            "metadata": self.metadata,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class ObjectiveValidationResult:
    """
    目标验证结果
    
    对生成的学习目标进行质量检查的结果
    """
    validation_id: str  # 验证ID
    objective_id: str  # 关联的目标ID
    
    # === 验证维度 ===
    is_measurable: bool = False  # 是否可衡量
    is_achievable: bool = False  # 是否可达成
    is_specific: bool = False  # 是否具体
    is_relevant: bool = False  # 是否相关
    is_time_bound: bool = False  # 是否有时间限制
    
    # === 综合评分 ===
    overall_score: float = 0.0  # 综合评分 (0-1)
    
    # === 问题和建议 ===
    issues: List[str] = field(default_factory=list)  # 发现的问题
    suggestions: List[str] = field(default_factory=list)  # 改进建议
    
    # === 元数据 ===
    validated_at: datetime = field(default_factory=datetime.now)
    validator: str = "ai"  # ai/human
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "validation_id": self.validation_id,
            "objective_id": self.objective_id,
            "is_measurable": self.is_measurable,
            "is_achievable": self.is_achievable,
            "is_specific": self.is_specific,
            "is_relevant": self.is_relevant,
            "is_time_bound": self.is_time_bound,
            "overall_score": self.overall_score,
            "issues": self.issues,
            "suggestions": self.suggestions,
            "validated_at": self.validated_at.isoformat(),
            "validator": self.validator
        }


@dataclass
class InteractiveFeedback:
    """
    交互式反馈
    
    用户对生成的目标的反馈，用于迭代改进
    """
    feedback_id: str  # 反馈ID
    objective_id: str  # 关联的目标ID
    
    # === 反馈类型 ===
    feedback_type: str  # like/dislike/modify/delete/add
    
    # === 反馈内容 ===
    comment: Optional[str] = None  # 文本评论
    modification: Optional[Dict[str, Any]] = None  # 修改建议
    # 示例: {"description": "新的描述", "bloom_level": 3}
    
    # === 元数据 ===
    user_id: Optional[str] = None
    created_at: datetime = field(default_factory=datetime.now)
    applied: bool = False  # 是否已应用
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "feedback_id": self.feedback_id,
            "objective_id": self.objective_id,
            "feedback_type": self.feedback_type,
            "comment": self.comment,
            "modification": self.modification,
            "user_id": self.user_id,
            "created_at": self.created_at.isoformat(),
            "applied": self.applied
        }
```

---

### 9.3 数据库表设计

```sql
-- 用户请求表
CREATE TABLE user_requests (
    request_id VARCHAR(100) PRIMARY KEY,
    raw_input TEXT NOT NULL,
    domain VARCHAR(50),
    
    -- 用户信息
    user_id VARCHAR(100),
    session_id VARCHAR(100),
    
    -- 上下文
    language VARCHAR(10) DEFAULT 'zh',
    region VARCHAR(50),
    
    -- 偏好
    preferences JSONB DEFAULT '{}'::jsonb,
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(50),
    user_agent TEXT
);

CREATE INDEX idx_requests_user ON user_requests(user_id);
CREATE INDEX idx_requests_session ON user_requests(session_id);
CREATE INDEX idx_requests_created ON user_requests(created_at);


-- 受众画像表
CREATE TABLE audience_profiles (
    profile_id VARCHAR(100) PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL,
    
    -- 基本属性
    age_range VARCHAR(20) NOT NULL,
    education_level VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,
    
    -- 兴趣和背景
    interests JSONB DEFAULT '[]'::jsonb,
    prior_knowledge JSONB DEFAULT '[]'::jsonb,
    learning_goals JSONB DEFAULT '[]'::jsonb,
    
    -- 学习特征
    learning_style VARCHAR(20) DEFAULT 'visual',
    motivation_level VARCHAR(10) DEFAULT 'medium',
    attention_span INTEGER DEFAULT 30,
    
    -- 挑战和需求
    common_errors JSONB DEFAULT '[]'::jsonb,
    pain_points JSONB DEFAULT '[]'::jsonb,
    special_needs JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    confidence FLOAT DEFAULT 1.0,
    source VARCHAR(20) DEFAULT 'user_provided',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_profile_request FOREIGN KEY (request_id) 
        REFERENCES user_requests(request_id) ON DELETE CASCADE
);

CREATE INDEX idx_profiles_request ON audience_profiles(request_id);


-- 学习目标表
CREATE TABLE learning_objectives (
    objective_id VARCHAR(100) PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    
    -- 分类
    bloom_level INTEGER NOT NULL CHECK (bloom_level >= 1 AND bloom_level <= 6),
    objective_type VARCHAR(20) NOT NULL,
    domain VARCHAR(50) NOT NULL,
    
    -- 层级关系
    parent_id VARCHAR(100),
    order_index INTEGER DEFAULT 0,
    
    -- 评估标准
    success_criteria JSONB DEFAULT '[]'::jsonb,
    assessment_methods JSONB DEFAULT '[]'::jsonb,
    
    -- 时间估算
    estimated_time INTEGER DEFAULT 0,
    difficulty VARCHAR(10) DEFAULT 'medium',
    
    -- 元数据
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_objective_request FOREIGN KEY (request_id) 
        REFERENCES user_requests(request_id) ON DELETE CASCADE,
    CONSTRAINT fk_objective_parent FOREIGN KEY (parent_id) 
        REFERENCES learning_objectives(objective_id) ON DELETE CASCADE
);

CREATE INDEX idx_objectives_request ON learning_objectives(request_id);
CREATE INDEX idx_objectives_parent ON learning_objectives(parent_id);
CREATE INDEX idx_objectives_domain ON learning_objectives(domain);
CREATE INDEX idx_objectives_bloom ON learning_objectives(bloom_level);


-- 目标验证结果表
CREATE TABLE objective_validation_results (
    validation_id VARCHAR(100) PRIMARY KEY,
    objective_id VARCHAR(100) NOT NULL,
    
    -- 验证维度
    is_measurable BOOLEAN DEFAULT false,
    is_achievable BOOLEAN DEFAULT false,
    is_specific BOOLEAN DEFAULT false,
    is_relevant BOOLEAN DEFAULT false,
    is_time_bound BOOLEAN DEFAULT false,
    
    -- 综合评分
    overall_score FLOAT DEFAULT 0.0,
    
    -- 问题和建议
    issues JSONB DEFAULT '[]'::jsonb,
    suggestions JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    validated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validator VARCHAR(20) DEFAULT 'ai',
    
    CONSTRAINT fk_validation_objective FOREIGN KEY (objective_id) 
        REFERENCES learning_objectives(objective_id) ON DELETE CASCADE
);

CREATE INDEX idx_validations_objective ON objective_validation_results(objective_id);
CREATE INDEX idx_validations_score ON objective_validation_results(overall_score);


-- 交互式反馈表
CREATE TABLE interactive_feedback (
    feedback_id VARCHAR(100) PRIMARY KEY,
    objective_id VARCHAR(100) NOT NULL,
    
    -- 反馈类型
    feedback_type VARCHAR(20) NOT NULL,
    
    -- 反馈内容
    comment TEXT,
    modification JSONB,
    
    -- 元数据
    user_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    applied BOOLEAN DEFAULT false,
    
    CONSTRAINT fk_feedback_objective FOREIGN KEY (objective_id) 
        REFERENCES learning_objectives(objective_id) ON DELETE CASCADE,
    CONSTRAINT check_feedback_type CHECK (feedback_type IN ('like', 'dislike', 'modify', 'delete', 'add'))
);

CREATE INDEX idx_feedback_objective ON interactive_feedback(objective_id);
CREATE INDEX idx_feedback_user ON interactive_feedback(user_id);
CREATE INDEX idx_feedback_applied ON interactive_feedback(applied);
```

---

## 10. V2.0新增：AI内容发现引擎数据模型 ⭐

### 10.1 概述

AI内容发现引擎通过向量检索和结构化提取，从多个内容源中发现相关内容。

**核心数据模型**：
1. `ContentSource` - 内容源配置
2. `SearchQuery` - 搜索查询
3. `ContentItem` - 内容条目
4. `StructuredContent` - 结构化内容
5. `RetrievalResult` - 检索结果
6. `QualityScore` - 质量评分

---

### 10.2 枚举类型

```python
from enum import Enum

class SourceType(Enum):
    """内容源类型"""
    WIKIPEDIA = "wikipedia"
    WIKIART = "wikiart"
    ARXIV = "arxiv"
    YOUTUBE = "youtube"
    EDUCATIONAL_DB = "educational_db"
    CUSTOM_API = "custom_api"
    LOCAL_FILES = "local_files"


class ContentType(Enum):
    """内容类型"""
    TEXT = "text"
    IMAGE = "image"
    VIDEO = "video"
    AUDIO = "audio"
    PDF = "pdf"
    STRUCTURED_DATA = "structured_data"


class QualityLevel(Enum):
    """质量等级"""
    EXCELLENT = "excellent"  # 90-100
    GOOD = "good"  # 70-89
    ACCEPTABLE = "acceptable"  # 50-69
    POOR = "poor"  # 0-49
```

---

### 10.3 核心数据类

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from datetime import datetime

@dataclass
class ContentSource:
    """
    内容源配置
    
    定义如何从特定来源获取内容
    """
    source_id: str  # 源ID
    name: str  # 源名称
    source_type: SourceType  # 源类型
    
    # === 连接配置 ===
    base_url: str  # 基础URL
    api_key: Optional[str] = None  # API密钥
    auth_type: str = "none"  # none/api_key/oauth/basic
    
    # === 能力配置 ===
    supported_content_types: List[ContentType] = field(default_factory=list)
    max_results_per_query: int = 100
    rate_limit: int = 10  # 每秒请求数
    
    # === 解析配置 ===
    parser_class: str = "DefaultParser"  # 解析器类名
    extraction_config: Dict[str, Any] = field(default_factory=dict)
    # 示例: {"selectors": {"title": "h1.title", "content": "div.content"}}
    
    # === 质量配置 ===
    default_quality_weight: float = 1.0  # 默认质量权重
    reliability_score: float = 0.8  # 可靠性评分 (0-1)
    
    # === 状态 ===
    enabled: bool = True
    last_successful_fetch: Optional[datetime] = None
    error_count: int = 0
    
    # === 元数据 ===
    metadata: Dict[str, Any] = field(default_factory=dict)
    created_at: datetime = field(default_factory=datetime.now)
    updated_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "source_id": self.source_id,
            "name": self.name,
            "source_type": self.source_type.value,
            "base_url": self.base_url,
            "api_key": "***" if self.api_key else None,
            "auth_type": self.auth_type,
            "supported_content_types": [ct.value for ct in self.supported_content_types],
            "max_results_per_query": self.max_results_per_query,
            "rate_limit": self.rate_limit,
            "parser_class": self.parser_class,
            "extraction_config": self.extraction_config,
            "default_quality_weight": self.default_quality_weight,
            "reliability_score": self.reliability_score,
            "enabled": self.enabled,
            "last_successful_fetch": self.last_successful_fetch.isoformat() if self.last_successful_fetch else None,
            "error_count": self.error_count,
            "metadata": self.metadata,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat()
        }


@dataclass
class SearchQuery:
    """
    搜索查询
    
    封装搜索参数和向量嵌入
    """
    query_id: str  # 查询ID
    text: str  # 查询文本
    
    # === 向量嵌入 ===
    embedding: Optional[List[float]] = None  # 查询向量 (1536维)
    embedding_model: str = "text-embedding-3-small"
    
    # === 过滤条件 ===
    filters: Dict[str, Any] = field(default_factory=dict)
    # 示例: {"content_type": "image", "min_quality": 0.7, "language": "zh"}
    
    # === 搜索参数 ===
    top_k: int = 10  # 返回结果数
    similarity_threshold: float = 0.5  # 相似度阈值
    
    # === 目标源 ===
    target_sources: List[str] = field(default_factory=list)  # 源ID列表（空表示所有）
    
    # === 上下文 ===
    domain: Optional[DomainType] = None
    objective_id: Optional[str] = None  # 关联的学习目标ID
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "query_id": self.query_id,
            "text": self.text,
            "embedding": self.embedding,
            "embedding_model": self.embedding_model,
            "filters": self.filters,
            "top_k": self.top_k,
            "similarity_threshold": self.similarity_threshold,
            "target_sources": self.target_sources,
            "domain": self.domain.value if self.domain else None,
            "objective_id": self.objective_id,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class ContentItem:
    """
    内容条目
    
    从内容源获取的单个内容项
    """
    item_id: str  # 条目ID
    source_id: str  # 来源ID
    
    # === 基本信息 ===
    title: str  # 标题
    content: str  # 内容文本
    content_type: ContentType  # 内容类型
    
    # === 内容详情 ===
    url: Optional[str] = None  # 原始URL
    author: Optional[str] = None  # 作者
    publish_date: Optional[datetime] = None  # 发布日期
    language: str = "zh"  # 语言
    
    # === 多媒体资源 ===
    images: List[str] = field(default_factory=list)  # 图片URL列表
    videos: List[str] = field(default_factory=list)  # 视频URL列表
    attachments: List[str] = field(default_factory=list)  # 附件URL列表
    
    # === 向量嵌入 ===
    embedding: Optional[List[float]] = None  # 内容向量
    embedding_model: str = "text-embedding-3-small"
    
    # === 标签和分类 ===
    tags: List[str] = field(default_factory=list)  # 标签
    categories: List[str] = field(default_factory=list)  # 分类
    
    # === 元数据 ===
    metadata: Dict[str, Any] = field(default_factory=dict)
    # 领域特定元数据，如美术史: {"artist": "莫奈", "movement": "印象派"}
    
    fetched_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "item_id": self.item_id,
            "source_id": self.source_id,
            "title": self.title,
            "content": self.content,
            "content_type": self.content_type.value,
            "url": self.url,
            "author": self.author,
            "publish_date": self.publish_date.isoformat() if self.publish_date else None,
            "language": self.language,
            "images": self.images,
            "videos": self.videos,
            "attachments": self.attachments,
            "embedding": self.embedding,
            "embedding_model": self.embedding_model,
            "tags": self.tags,
            "categories": self.categories,
            "metadata": self.metadata,
            "fetched_at": self.fetched_at.isoformat()
        }


@dataclass
class StructuredContent:
    """
    结构化内容
    
    从ContentItem中提取的结构化数据
    """
    structured_id: str  # 结构化内容ID
    item_id: str  # 关联的ContentItem ID
    
    # === 结构化字段 ===
    structured_data: Dict[str, Any]  # 结构化数据
    # K12示例: {"concept": "除法", "examples": [...], "exercises": [...]}
    # 美术史示例: {"artwork": {...}, "artist_bio": {...}, "analysis": {...}}
    # 畅销书示例: {"summary": "...", "key_points": [...], "quotes": [...]}
    
    # === 提取信息 ===
    extraction_method: str = "llm"  # llm/regex/xpath/custom
    extractor_model: Optional[str] = None  # 如 "gpt-4"
    extraction_confidence: float = 0.0  # 提取置信度 (0-1)
    
    # === 验证结果 ===
    is_validated: bool = False
    validation_errors: List[str] = field(default_factory=list)
    
    # === 元数据 ===
    extracted_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "structured_id": self.structured_id,
            "item_id": self.item_id,
            "structured_data": self.structured_data,
            "extraction_method": self.extraction_method,
            "extractor_model": self.extractor_model,
            "extraction_confidence": self.extraction_confidence,
            "is_validated": self.is_validated,
            "validation_errors": self.validation_errors,
            "extracted_at": self.extracted_at.isoformat()
        }


@dataclass
class QualityScore:
    """
    质量评分
    
    对ContentItem的质量进行多维度评估
    """
    score_id: str  # 评分ID
    item_id: str  # 关联的ContentItem ID
    
    # === 质量维度 ===
    relevance: float = 0.0  # 相关性 (0-1)
    accuracy: float = 0.0  # 准确性 (0-1)
    completeness: float = 0.0  # 完整性 (0-1)
    readability: float = 0.0  # 可读性 (0-1)
    authority: float = 0.0  # 权威性 (0-1)
    freshness: float = 0.0  # 时效性 (0-1)
    
    # === 综合评分 ===
    overall_score: float = 0.0  # 综合评分 (0-1)
    quality_level: QualityLevel = QualityLevel.ACCEPTABLE
    
    # === 评分详情 ===
    scoring_method: str = "llm"  # llm/rule_based/hybrid
    scorer_model: Optional[str] = None  # 如 "gpt-4"
    
    # === 问题标记 ===
    issues: List[str] = field(default_factory=list)
    # 示例: ["outdated_information", "citation_needed", "bias_detected"]
    
    # === 元数据 ===
    scored_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "score_id": self.score_id,
            "item_id": self.item_id,
            "relevance": self.relevance,
            "accuracy": self.accuracy,
            "completeness": self.completeness,
            "readability": self.readability,
            "authority": self.authority,
            "freshness": self.freshness,
            "overall_score": self.overall_score,
            "quality_level": self.quality_level.value,
            "scoring_method": self.scoring_method,
            "scorer_model": self.scorer_model,
            "issues": self.issues,
            "scored_at": self.scored_at.isoformat()
        }


@dataclass
class RetrievalResult:
    """
    检索结果
    
    单次搜索的完整结果，包含多个内容项和排序信息
    """
    result_id: str  # 结果ID
    query_id: str  # 关联的查询ID
    
    # === 结果列表 ===
    items: List[ContentItem] = field(default_factory=list)
    
    # === 排序信息 ===
    item_scores: Dict[str, float] = field(default_factory=dict)
    # key: item_id, value: 相似度/相关性评分
    
    ranking_method: str = "vector_similarity"  # vector_similarity/hybrid/reranker
    
    # === 统计信息 ===
    total_matches: int = 0  # 总匹配数（可能大于返回数）
    search_time: float = 0.0  # 搜索耗时（秒）
    
    # === 源分布 ===
    source_distribution: Dict[str, int] = field(default_factory=dict)
    # key: source_id, value: 该源的结果数
    
    # === 元数据 ===
    retrieved_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "result_id": self.result_id,
            "query_id": self.query_id,
            "items": [item.to_dict() for item in self.items],
            "item_scores": self.item_scores,
            "ranking_method": self.ranking_method,
            "total_matches": self.total_matches,
            "search_time": self.search_time,
            "source_distribution": self.source_distribution,
            "retrieved_at": self.retrieved_at.isoformat()
        }
```

---

### 10.4 数据库表设计

```sql
-- 内容源配置表
CREATE TABLE content_sources (
    source_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    
    -- 连接配置
    base_url VARCHAR(500) NOT NULL,
    api_key VARCHAR(500),  -- 加密存储
    auth_type VARCHAR(20) DEFAULT 'none',
    
    -- 能力配置
    supported_content_types JSONB DEFAULT '[]'::jsonb,
    max_results_per_query INTEGER DEFAULT 100,
    rate_limit INTEGER DEFAULT 10,
    
    -- 解析配置
    parser_class VARCHAR(100) DEFAULT 'DefaultParser',
    extraction_config JSONB DEFAULT '{}'::jsonb,
    
    -- 质量配置
    default_quality_weight FLOAT DEFAULT 1.0,
    reliability_score FLOAT DEFAULT 0.8,
    
    -- 状态
    enabled BOOLEAN DEFAULT true,
    last_successful_fetch TIMESTAMP,
    error_count INTEGER DEFAULT 0,
    
    -- 元数据
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_source_type CHECK (source_type IN ('wikipedia', 'wikiart', 'arxiv', 'youtube', 'educational_db', 'custom_api', 'local_files')),
    CONSTRAINT check_auth_type CHECK (auth_type IN ('none', 'api_key', 'oauth', 'basic'))
);

CREATE INDEX idx_sources_type ON content_sources(source_type);
CREATE INDEX idx_sources_enabled ON content_sources(enabled);


-- 搜索查询表
CREATE TABLE search_queries (
    query_id VARCHAR(100) PRIMARY KEY,
    text TEXT NOT NULL,
    
    -- 向量嵌入（存储在Milvus，这里只记录引用）
    embedding_model VARCHAR(100) DEFAULT 'text-embedding-3-small',
    
    -- 过滤条件
    filters JSONB DEFAULT '{}'::jsonb,
    
    -- 搜索参数
    top_k INTEGER DEFAULT 10,
    similarity_threshold FLOAT DEFAULT 0.5,
    
    -- 目标源
    target_sources JSONB DEFAULT '[]'::jsonb,
    
    -- 上下文
    domain VARCHAR(50),
    objective_id VARCHAR(100),
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_query_objective FOREIGN KEY (objective_id) 
        REFERENCES learning_objectives(objective_id) ON DELETE SET NULL
);

CREATE INDEX idx_queries_domain ON search_queries(domain);
CREATE INDEX idx_queries_objective ON search_queries(objective_id);
CREATE INDEX idx_queries_created ON search_queries(created_at);


-- 内容条目表
CREATE TABLE content_items (
    item_id VARCHAR(100) PRIMARY KEY,
    source_id VARCHAR(100) NOT NULL,
    
    -- 基本信息
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    
    -- 内容详情
    url VARCHAR(1000),
    author VARCHAR(200),
    publish_date TIMESTAMP,
    language VARCHAR(10) DEFAULT 'zh',
    
    -- 多媒体资源
    images JSONB DEFAULT '[]'::jsonb,
    videos JSONB DEFAULT '[]'::jsonb,
    attachments JSONB DEFAULT '[]'::jsonb,
    
    -- 向量嵌入（存储在Milvus）
    embedding_model VARCHAR(100) DEFAULT 'text-embedding-3-small',
    
    -- 标签和分类
    tags JSONB DEFAULT '[]'::jsonb,
    categories JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    metadata JSONB DEFAULT '{}'::jsonb,
    fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_item_source FOREIGN KEY (source_id) 
        REFERENCES content_sources(source_id) ON DELETE CASCADE,
    CONSTRAINT check_content_type CHECK (content_type IN ('text', 'image', 'video', 'audio', 'pdf', 'structured_data'))
);

CREATE INDEX idx_items_source ON content_items(source_id);
CREATE INDEX idx_items_type ON content_items(content_type);
CREATE INDEX idx_items_language ON content_items(language);
CREATE INDEX idx_items_fetched ON content_items(fetched_at);
CREATE INDEX idx_items_tags ON content_items USING GIN(tags);


-- 结构化内容表
CREATE TABLE structured_contents (
    structured_id VARCHAR(100) PRIMARY KEY,
    item_id VARCHAR(100) NOT NULL,
    
    -- 结构化字段
    structured_data JSONB NOT NULL,
    
    -- 提取信息
    extraction_method VARCHAR(50) DEFAULT 'llm',
    extractor_model VARCHAR(100),
    extraction_confidence FLOAT DEFAULT 0.0,
    
    -- 验证结果
    is_validated BOOLEAN DEFAULT false,
    validation_errors JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    extracted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_structured_item FOREIGN KEY (item_id) 
        REFERENCES content_items(item_id) ON DELETE CASCADE
);

CREATE INDEX idx_structured_item ON structured_contents(item_id);
CREATE INDEX idx_structured_validated ON structured_contents(is_validated);


-- 质量评分表
CREATE TABLE quality_scores (
    score_id VARCHAR(100) PRIMARY KEY,
    item_id VARCHAR(100) NOT NULL,
    
    -- 质量维度
    relevance FLOAT DEFAULT 0.0,
    accuracy FLOAT DEFAULT 0.0,
    completeness FLOAT DEFAULT 0.0,
    readability FLOAT DEFAULT 0.0,
    authority FLOAT DEFAULT 0.0,
    freshness FLOAT DEFAULT 0.0,
    
    -- 综合评分
    overall_score FLOAT DEFAULT 0.0,
    quality_level VARCHAR(20) DEFAULT 'acceptable',
    
    -- 评分详情
    scoring_method VARCHAR(50) DEFAULT 'llm',
    scorer_model VARCHAR(100),
    
    -- 问题标记
    issues JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    scored_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_score_item FOREIGN KEY (item_id) 
        REFERENCES content_items(item_id) ON DELETE CASCADE,
    CONSTRAINT check_quality_level CHECK (quality_level IN ('excellent', 'good', 'acceptable', 'poor'))
);

CREATE INDEX idx_scores_item ON quality_scores(item_id);
CREATE INDEX idx_scores_overall ON quality_scores(overall_score);
CREATE INDEX idx_scores_level ON quality_scores(quality_level);


-- 检索结果表
CREATE TABLE retrieval_results (
    result_id VARCHAR(100) PRIMARY KEY,
    query_id VARCHAR(100) NOT NULL,
    
    -- 结果列表（item_id数组）
    item_ids JSONB NOT NULL,
    
    -- 排序信息
    item_scores JSONB DEFAULT '{}'::jsonb,  -- {item_id: score}
    ranking_method VARCHAR(50) DEFAULT 'vector_similarity',
    
    -- 统计信息
    total_matches INTEGER DEFAULT 0,
    search_time FLOAT DEFAULT 0.0,
    
    -- 源分布
    source_distribution JSONB DEFAULT '{}'::jsonb,
    
    -- 元数据
    retrieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_result_query FOREIGN KEY (query_id) 
        REFERENCES search_queries(query_id) ON DELETE CASCADE
);

CREATE INDEX idx_results_query ON retrieval_results(query_id);
CREATE INDEX idx_results_retrieved ON retrieval_results(retrieved_at);
```

---

## 11. V2.0新增：故事化叙述引擎数据模型 ⭐

### 11.1 概述

故事化叙述引擎将枯燥内容转化为引人入胜的故事化叙述，提升学习体验。

**核心数据模型**：
1. `NarrativeStrategy` - 叙述策略
2. `StoryElement` - 故事元素
3. `NarrativeSection` - 叙述段落
4. `NarrativeContent` - 完整叙述内容
5. `QualityIssue` - 质量问题

---

### 11.2 枚举类型

```python
from enum import Enum

class NarrativeStyle(Enum):
    """叙述风格"""
    STORY_TELLING = "story_telling"  # 故事讲述
    HISTORICAL_NARRATIVE = "historical_narrative"  # 历史叙述
    CASE_STUDY = "case_study"  # 案例研究
    DIALOGUE = "dialogue"  # 对话式
    DOCUMENTARY = "documentary"  # 纪录片式
    ADVENTURE = "adventure"  # 冒险式
    BIOGRAPHY = "biography"  # 传记式


class StoryArc(Enum):
    """故事弧线"""
    HEROES_JOURNEY = "heroes_journey"  # 英雄之旅
    THREE_ACT = "three_act"  # 三幕式
    FIVE_ACT = "five_act"  # 五幕式
    PROBLEM_SOLUTION = "problem_solution"  # 问题-解决
    CHRONOLOGICAL = "chronological"  # 时间顺序
    FLASHBACK = "flashback"  # 倒叙


class EmotionalTone(Enum):
    """情感基调"""
    INSPIRING = "inspiring"  # 激励
    CURIOUS = "curious"  # 好奇
    MYSTERIOUS = "mysterious"  # 神秘
    HUMOROUS = "humorous"  # 幽默
    DRAMATIC = "dramatic"  # 戏剧化
    CALM = "calm"  # 平静
    EXCITING = "exciting"  # 兴奋
```

---

### 11.3 核心数据类

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from datetime import datetime

@dataclass
class NarrativeStrategy:
    """
    叙述策略
    
    定义如何将内容转化为故事化叙述
    """
    strategy_id: str  # 策略ID
    name: str  # 策略名称
    domain: DomainType  # 适用领域
    
    # === 风格设置 ===
    narrative_style: NarrativeStyle  # 叙述风格
    story_arc: StoryArc  # 故事弧线
    emotional_tone: EmotionalTone  # 情感基调
    
    # === 目标受众 ===
    target_age_range: str  # 目标年龄段
    reading_level: str  # 阅读水平
    
    # === 叙述元素 ===
    use_characters: bool = True  # 是否使用人物
    use_dialogue: bool = True  # 是否使用对话
    use_metaphors: bool = True  # 是否使用隐喻
    use_suspense: bool = False  # 是否制造悬念
    
    # === 结构配置 ===
    max_sections: int = 10  # 最大段落数
    section_length_range: tuple = (200, 500)  # 段落长度范围（字数）
    
    # === 提示词模板 ===
    prompt_template: str = ""  # LLM提示词模板
    # 示例: "将以下内容改写为{narrative_style}风格的故事，目标受众是{target_age_range}..."
    
    # === 质量标准 ===
    min_engagement_score: float = 0.7  # 最低吸引力评分
    min_coherence_score: float = 0.8  # 最低连贯性评分
    
    # === 元数据 ===
    created_by: str = "system"
    created_at: datetime = field(default_factory=datetime.now)
    usage_count: int = 0
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "strategy_id": self.strategy_id,
            "name": self.name,
            "domain": self.domain.value,
            "narrative_style": self.narrative_style.value,
            "story_arc": self.story_arc.value,
            "emotional_tone": self.emotional_tone.value,
            "target_age_range": self.target_age_range,
            "reading_level": self.reading_level,
            "use_characters": self.use_characters,
            "use_dialogue": self.use_dialogue,
            "use_metaphors": self.use_metaphors,
            "use_suspense": self.use_suspense,
            "max_sections": self.max_sections,
            "section_length_range": list(self.section_length_range),
            "prompt_template": self.prompt_template,
            "min_engagement_score": self.min_engagement_score,
            "min_coherence_score": self.min_coherence_score,
            "created_by": self.created_by,
            "created_at": self.created_at.isoformat(),
            "usage_count": self.usage_count
        }


@dataclass
class StoryElement:
    """
    故事元素
    
    叙述中的关键元素（人物、场景、冲突等）
    """
    element_id: str  # 元素ID
    element_type: str  # character/setting/conflict/resolution/theme
    
    # === 元素内容 ===
    name: str  # 元素名称
    description: str  # 描述
    
    # === 人物特定字段（element_type=character） ===
    character_traits: List[str] = field(default_factory=list)  # 性格特点
    character_role: Optional[str] = None  # protagonist/antagonist/mentor/sidekick
    
    # === 场景特定字段（element_type=setting） ===
    time_period: Optional[str] = None  # 时间
    location: Optional[str] = None  # 地点
    atmosphere: Optional[str] = None  # 氛围
    
    # === 冲突特定字段（element_type=conflict） ===
    conflict_type: Optional[str] = None  # internal/external/social
    stakes: Optional[str] = None  # 赌注/后果
    
    # === 关联 ===
    related_elements: List[str] = field(default_factory=list)  # 相关元素ID
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "element_id": self.element_id,
            "element_type": self.element_type,
            "name": self.name,
            "description": self.description,
            "character_traits": self.character_traits,
            "character_role": self.character_role,
            "time_period": self.time_period,
            "location": self.location,
            "atmosphere": self.atmosphere,
            "conflict_type": self.conflict_type,
            "stakes": self.stakes,
            "related_elements": self.related_elements,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class NarrativeSection:
    """
    叙述段落
    
    完整叙述内容中的一个段落/章节
    """
    section_id: str  # 段落ID
    narrative_id: str  # 所属叙述内容ID
    
    # === 位置信息 ===
    order: int  # 顺序
    act: Optional[int] = None  # 幕（用于多幕结构）
    
    # === 内容 ===
    title: str = ""  # 段落标题
    content: str = ""  # 段落内容
    word_count: int = 0  # 字数
    
    # === 故事功能 ===
    narrative_function: str = "development"  # exposition/rising_action/climax/falling_action/resolution/development
    
    # === 涉及元素 ===
    story_elements: List[str] = field(default_factory=list)  # 涉及的故事元素ID
    
    # === 转场 ===
    transition_in: Optional[str] = None  # 入场转场
    transition_out: Optional[str] = None  # 出场转场
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "section_id": self.section_id,
            "narrative_id": self.narrative_id,
            "order": self.order,
            "act": self.act,
            "title": self.title,
            "content": self.content,
            "word_count": self.word_count,
            "narrative_function": self.narrative_function,
            "story_elements": self.story_elements,
            "transition_in": self.transition_in,
            "transition_out": self.transition_out,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class NarrativeContent:
    """
    完整叙述内容
    
    基于某个学习目标生成的完整故事化叙述
    """
    narrative_id: str  # 叙述内容ID
    objective_id: str  # 关联的学习目标ID
    strategy_id: str  # 使用的策略ID
    
    # === 基本信息 ===
    title: str  # 标题
    subtitle: Optional[str] = None  # 副标题
    
    # === 内容结构 ===
    sections: List[NarrativeSection] = field(default_factory=list)  # 段落列表
    story_elements: List[StoryElement] = field(default_factory=list)  # 故事元素
    
    # === 统计信息 ===
    total_word_count: int = 0  # 总字数
    estimated_reading_time: int = 0  # 预估阅读时间（分钟）
    
    # === 质量评估 ===
    engagement_score: float = 0.0  # 吸引力评分 (0-1)
    coherence_score: float = 0.0  # 连贯性评分 (0-1)
    educational_value: float = 0.0  # 教育价值评分 (0-1)
    age_appropriateness: float = 0.0  # 年龄适宜性 (0-1)
    
    # === 多媒体资源 ===
    images: List[str] = field(default_factory=list)  # 配图URL
    audio_narration: Optional[str] = None  # 音频朗读URL
    
    # === 元数据 ===
    generated_at: datetime = field(default_factory=datetime.now)
    generation_time: float = 0.0  # 生成耗时（秒）
    llm_model: str = "gpt-4"
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "narrative_id": self.narrative_id,
            "objective_id": self.objective_id,
            "strategy_id": self.strategy_id,
            "title": self.title,
            "subtitle": self.subtitle,
            "sections": [s.to_dict() for s in self.sections],
            "story_elements": [e.to_dict() for e in self.story_elements],
            "total_word_count": self.total_word_count,
            "estimated_reading_time": self.estimated_reading_time,
            "engagement_score": self.engagement_score,
            "coherence_score": self.coherence_score,
            "educational_value": self.educational_value,
            "age_appropriateness": self.age_appropriateness,
            "images": self.images,
            "audio_narration": self.audio_narration,
            "generated_at": self.generated_at.isoformat(),
            "generation_time": self.generation_time,
            "llm_model": self.llm_model
        }


@dataclass
class QualityIssue:
    """
    质量问题
    
    在叙述内容中发现的质量问题
    """
    issue_id: str  # 问题ID
    narrative_id: str  # 关联的叙述内容ID
    section_id: Optional[str] = None  # 关联的段落ID（可选）
    
    # === 问题类型 ===
    issue_type: str  # coherence/factual_error/inappropriate_content/readability/engagement
    severity: str = "medium"  # low/medium/high/critical
    
    # === 问题详情 ===
    description: str = ""  # 问题描述
    affected_text: Optional[str] = None  # 受影响的文本片段
    position: Optional[int] = None  # 文本位置（字符索引）
    
    # === 修复建议 ===
    suggestion: Optional[str] = None  # 修复建议
    auto_fixable: bool = False  # 是否可自动修复
    
    # === 状态 ===
    status: str = "open"  # open/fixed/ignored
    fixed_at: Optional[datetime] = None
    
    # === 元数据 ===
    detected_by: str = "ai"  # ai/human/rule_based
    detected_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "issue_id": self.issue_id,
            "narrative_id": self.narrative_id,
            "section_id": self.section_id,
            "issue_type": self.issue_type,
            "severity": self.severity,
            "description": self.description,
            "affected_text": self.affected_text,
            "position": self.position,
            "suggestion": self.suggestion,
            "auto_fixable": self.auto_fixable,
            "status": self.status,
            "fixed_at": self.fixed_at.isoformat() if self.fixed_at else None,
            "detected_by": self.detected_by,
            "detected_at": self.detected_at.isoformat()
        }
```

---

### 11.4 数据库表设计

```sql
-- 叙述策略表
CREATE TABLE narrative_strategies (
    strategy_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    domain VARCHAR(50) NOT NULL,
    
    -- 风格设置
    narrative_style VARCHAR(50) NOT NULL,
    story_arc VARCHAR(50) NOT NULL,
    emotional_tone VARCHAR(50) NOT NULL,
    
    -- 目标受众
    target_age_range VARCHAR(20) NOT NULL,
    reading_level VARCHAR(50) NOT NULL,
    
    -- 叙述元素
    use_characters BOOLEAN DEFAULT true,
    use_dialogue BOOLEAN DEFAULT true,
    use_metaphors BOOLEAN DEFAULT true,
    use_suspense BOOLEAN DEFAULT false,
    
    -- 结构配置
    max_sections INTEGER DEFAULT 10,
    section_length_min INTEGER DEFAULT 200,
    section_length_max INTEGER DEFAULT 500,
    
    -- 提示词模板
    prompt_template TEXT,
    
    -- 质量标准
    min_engagement_score FLOAT DEFAULT 0.7,
    min_coherence_score FLOAT DEFAULT 0.8,
    
    -- 元数据
    created_by VARCHAR(100) DEFAULT 'system',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usage_count INTEGER DEFAULT 0,
    
    CONSTRAINT check_narrative_style CHECK (narrative_style IN ('story_telling', 'historical_narrative', 'case_study', 'dialogue', 'documentary', 'adventure', 'biography')),
    CONSTRAINT check_story_arc CHECK (story_arc IN ('heroes_journey', 'three_act', 'five_act', 'problem_solution', 'chronological', 'flashback')),
    CONSTRAINT check_emotional_tone CHECK (emotional_tone IN ('inspiring', 'curious', 'mysterious', 'humorous', 'dramatic', 'calm', 'exciting'))
);

CREATE INDEX idx_strategies_domain ON narrative_strategies(domain);
CREATE INDEX idx_strategies_style ON narrative_strategies(narrative_style);


-- 故事元素表
CREATE TABLE story_elements (
    element_id VARCHAR(100) PRIMARY KEY,
    narrative_id VARCHAR(100) NOT NULL,
    element_type VARCHAR(50) NOT NULL,
    
    -- 元素内容
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    
    -- 人物特定字段
    character_traits JSONB DEFAULT '[]'::jsonb,
    character_role VARCHAR(50),
    
    -- 场景特定字段
    time_period VARCHAR(100),
    location VARCHAR(200),
    atmosphere VARCHAR(100),
    
    -- 冲突特定字段
    conflict_type VARCHAR(50),
    stakes TEXT,
    
    -- 关联
    related_elements JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_element_type CHECK (element_type IN ('character', 'setting', 'conflict', 'resolution', 'theme'))
);

CREATE INDEX idx_elements_narrative ON story_elements(narrative_id);
CREATE INDEX idx_elements_type ON story_elements(element_type);


-- 叙述段落表
CREATE TABLE narrative_sections (
    section_id VARCHAR(100) PRIMARY KEY,
    narrative_id VARCHAR(100) NOT NULL,
    
    -- 位置信息
    order_index INTEGER NOT NULL,
    act INTEGER,
    
    -- 内容
    title VARCHAR(300),
    content TEXT NOT NULL,
    word_count INTEGER DEFAULT 0,
    
    -- 故事功能
    narrative_function VARCHAR(50) DEFAULT 'development',
    
    -- 涉及元素
    story_elements JSONB DEFAULT '[]'::jsonb,
    
    -- 转场
    transition_in VARCHAR(200),
    transition_out VARCHAR(200),
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_narrative_function CHECK (narrative_function IN ('exposition', 'rising_action', 'climax', 'falling_action', 'resolution', 'development'))
);

CREATE INDEX idx_sections_narrative ON narrative_sections(narrative_id);
CREATE INDEX idx_sections_order ON narrative_sections(narrative_id, order_index);


-- 叙述内容表
CREATE TABLE narrative_contents (
    narrative_id VARCHAR(100) PRIMARY KEY,
    objective_id VARCHAR(100) NOT NULL,
    strategy_id VARCHAR(100) NOT NULL,
    
    -- 基本信息
    title VARCHAR(300) NOT NULL,
    subtitle VARCHAR(300),
    
    -- 统计信息
    total_word_count INTEGER DEFAULT 0,
    estimated_reading_time INTEGER DEFAULT 0,
    
    -- 质量评估
    engagement_score FLOAT DEFAULT 0.0,
    coherence_score FLOAT DEFAULT 0.0,
    educational_value FLOAT DEFAULT 0.0,
    age_appropriateness FLOAT DEFAULT 0.0,
    
    -- 多媒体资源
    images JSONB DEFAULT '[]'::jsonb,
    audio_narration VARCHAR(500),
    
    -- 元数据
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    generation_time FLOAT DEFAULT 0.0,
    llm_model VARCHAR(50) DEFAULT 'gpt-4',
    
    CONSTRAINT fk_narrative_objective FOREIGN KEY (objective_id) 
        REFERENCES learning_objectives(objective_id) ON DELETE CASCADE,
    CONSTRAINT fk_narrative_strategy FOREIGN KEY (strategy_id) 
        REFERENCES narrative_strategies(strategy_id) ON DELETE SET NULL
);

CREATE INDEX idx_narratives_objective ON narrative_contents(objective_id);
CREATE INDEX idx_narratives_strategy ON narrative_contents(strategy_id);
CREATE INDEX idx_narratives_engagement ON narrative_contents(engagement_score);
CREATE INDEX idx_narratives_generated ON narrative_contents(generated_at);

-- 添加外键约束（延后创建）
ALTER TABLE story_elements 
    ADD CONSTRAINT fk_element_narrative FOREIGN KEY (narrative_id) 
    REFERENCES narrative_contents(narrative_id) ON DELETE CASCADE;

ALTER TABLE narrative_sections 
    ADD CONSTRAINT fk_section_narrative FOREIGN KEY (narrative_id) 
    REFERENCES narrative_contents(narrative_id) ON DELETE CASCADE;


-- 质量问题表
CREATE TABLE quality_issues (
    issue_id VARCHAR(100) PRIMARY KEY,
    narrative_id VARCHAR(100) NOT NULL,
    section_id VARCHAR(100),
    
    -- 问题类型
    issue_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) DEFAULT 'medium',
    
    -- 问题详情
    description TEXT NOT NULL,
    affected_text TEXT,
    position INTEGER,
    
    -- 修复建议
    suggestion TEXT,
    auto_fixable BOOLEAN DEFAULT false,
    
    -- 状态
    status VARCHAR(20) DEFAULT 'open',
    fixed_at TIMESTAMP,
    
    -- 元数据
    detected_by VARCHAR(20) DEFAULT 'ai',
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_issue_narrative FOREIGN KEY (narrative_id) 
        REFERENCES narrative_contents(narrative_id) ON DELETE CASCADE,
    CONSTRAINT fk_issue_section FOREIGN KEY (section_id) 
        REFERENCES narrative_sections(section_id) ON DELETE CASCADE,
    CONSTRAINT check_issue_type CHECK (issue_type IN ('coherence', 'factual_error', 'inappropriate_content', 'readability', 'engagement')),
    CONSTRAINT check_severity CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT check_status CHECK (status IN ('open', 'fixed', 'ignored'))
);

CREATE INDEX idx_issues_narrative ON quality_issues(narrative_id);
CREATE INDEX idx_issues_section ON quality_issues(section_id);
CREATE INDEX idx_issues_status ON quality_issues(status);
CREATE INDEX idx_issues_severity ON quality_issues(severity);
```

---

## 12. V2.0新增：多源内容检索数据模型 ⭐

### 12.1 概述

多源内容检索整合Milvus向量数据库、Redis缓存、MinIO对象存储，实现高效的内容检索和存储。

**核心数据模型**：
1. `VectorSearchQuery` - 向量搜索查询
2. `VectorSearchResult` - 向量搜索结果
3. `CacheEntry` - 缓存条目
4. `ObjectStorageMetadata` - 对象存储元数据
5. `MultiSourceQuery` - 多源查询

---

### 12.2 枚举类型

```python
from enum import Enum

class StorageBackend(Enum):
    """存储后端类型"""
    MILVUS = "milvus"  # 向量数据库
    REDIS = "redis"  # 缓存
    MINIO = "minio"  # 对象存储
    POSTGRESQL = "postgresql"  # 关系数据库


class CacheStrategy(Enum):
    """缓存策略"""
    LRU = "lru"  # 最近最少使用
    LFU = "lfu"  # 最不经常使用
    TTL = "ttl"  # 基于时间
    ADAPTIVE = "adaptive"  # 自适应


class IndexType(Enum):
    """向量索引类型（Milvus）"""
    FLAT = "FLAT"  # 暴力搜索
    IVF_FLAT = "IVF_FLAT"  # 倒排文件
    IVF_SQ8 = "IVF_SQ8"  # 倒排文件+标量量化
    IVF_PQ = "IVF_PQ"  # 倒排文件+乘积量化
    HNSW = "HNSW"  # 层次可导航小世界图
    ANNOY = "ANNOY"  # Spotify的ANNOY算法
```

---

### 12.3 核心数据类

```python
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Any
from datetime import datetime

@dataclass
class VectorSearchQuery:
    """
    向量搜索查询（Milvus）
    
    用于在Milvus中执行向量相似度搜索
    """
    query_id: str  # 查询ID
    collection_name: str  # 集合名称（如 "content_embeddings"）
    
    # === 查询向量 ===
    query_vector: List[float]  # 查询向量（1536维）
    
    # === 搜索参数 ===
    top_k: int = 10  # 返回结果数
    metric_type: str = "L2"  # L2/IP/COSINE
    search_params: Dict[str, Any] = field(default_factory=dict)
    # 示例: {"nprobe": 10, "ef": 64}  # 针对不同索引类型的参数
    
    # === 过滤条件 ===
    filter_expression: Optional[str] = None
    # 示例: "domain == 'k12' and quality_score > 0.7"
    
    # === 输出字段 ===
    output_fields: List[str] = field(default_factory=list)
    # 示例: ["item_id", "title", "quality_score"]
    
    # === 元数据 ===
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "query_id": self.query_id,
            "collection_name": self.collection_name,
            "query_vector": self.query_vector,
            "top_k": self.top_k,
            "metric_type": self.metric_type,
            "search_params": self.search_params,
            "filter_expression": self.filter_expression,
            "output_fields": self.output_fields,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class VectorSearchResult:
    """
    向量搜索结果（Milvus）
    
    Milvus返回的搜索结果
    """
    result_id: str  # 结果ID
    query_id: str  # 关联的查询ID
    
    # === 结果列表 ===
    hits: List[Dict[str, Any]] = field(default_factory=list)
    # 示例:
    # [
    #     {"id": "item_123", "distance": 0.23, "entity": {"title": "...", "quality_score": 0.85}},
    #     {"id": "item_456", "distance": 0.31, "entity": {...}}
    # ]
    
    # === 统计信息 ===
    total_hits: int = 0  # 命中数
    search_time: float = 0.0  # 搜索耗时（毫秒）
    
    # === 元数据 ===
    searched_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "result_id": self.result_id,
            "query_id": self.query_id,
            "hits": self.hits,
            "total_hits": self.total_hits,
            "search_time": self.search_time,
            "searched_at": self.searched_at.isoformat()
        }


@dataclass
class CacheEntry:
    """
    缓存条目（Redis）
    
    存储在Redis中的缓存项
    """
    cache_key: str  # 缓存键
    cache_value: Any  # 缓存值（序列化为JSON）
    
    # === 缓存策略 ===
    ttl: int = 3600  # 生存时间（秒），0表示永久
    strategy: CacheStrategy = CacheStrategy.TTL
    
    # === 访问统计 ===
    hit_count: int = 0  # 命中次数
    last_accessed: datetime = field(default_factory=datetime.now)
    
    # === 元数据 ===
    data_type: str = "json"  # json/string/binary
    compressed: bool = False  # 是否压缩
    size_bytes: int = 0  # 大小（字节）
    
    created_at: datetime = field(default_factory=datetime.now)
    expires_at: Optional[datetime] = None  # 过期时间
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "cache_key": self.cache_key,
            "cache_value": self.cache_value,
            "ttl": self.ttl,
            "strategy": self.strategy.value,
            "hit_count": self.hit_count,
            "last_accessed": self.last_accessed.isoformat(),
            "data_type": self.data_type,
            "compressed": self.compressed,
            "size_bytes": self.size_bytes,
            "created_at": self.created_at.isoformat(),
            "expires_at": self.expires_at.isoformat() if self.expires_at else None
        }


@dataclass
class ObjectStorageMetadata:
    """
    对象存储元数据（MinIO）
    
    MinIO中存储对象的元数据
    """
    object_id: str  # 对象ID
    bucket_name: str  # 桶名称
    object_key: str  # 对象键（路径）
    
    # === 对象信息 ===
    content_type: str  # MIME类型（如 "image/png", "application/pdf"）
    size_bytes: int  # 大小（字节）
    etag: str  # ETag（MD5哈希）
    
    # === 存储配置 ===
    storage_class: str = "STANDARD"  # STANDARD/REDUCED_REDUNDANCY
    encryption: bool = False  # 是否加密
    
    # === 访问控制 ===
    is_public: bool = False  # 是否公开访问
    signed_url: Optional[str] = None  # 签名URL（临时访问）
    signed_url_expires: Optional[datetime] = None  # 签名URL过期时间
    
    # === 关联信息 ===
    related_entity_type: Optional[str] = None  # 关联实体类型（如 "content_item"）
    related_entity_id: Optional[str] = None  # 关联实体ID
    
    # === 元数据 ===
    custom_metadata: Dict[str, str] = field(default_factory=dict)
    # 示例: {"domain": "k12", "quality": "high"}
    
    uploaded_at: datetime = field(default_factory=datetime.now)
    last_modified: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "object_id": self.object_id,
            "bucket_name": self.bucket_name,
            "object_key": self.object_key,
            "content_type": self.content_type,
            "size_bytes": self.size_bytes,
            "etag": self.etag,
            "storage_class": self.storage_class,
            "encryption": self.encryption,
            "is_public": self.is_public,
            "signed_url": self.signed_url,
            "signed_url_expires": self.signed_url_expires.isoformat() if self.signed_url_expires else None,
            "related_entity_type": self.related_entity_type,
            "related_entity_id": self.related_entity_id,
            "custom_metadata": self.custom_metadata,
            "uploaded_at": self.uploaded_at.isoformat(),
            "last_modified": self.last_modified.isoformat()
        }


@dataclass
class MultiSourceQuery:
    """
    多源查询
    
    统一封装来自多个存储后端的查询
    """
    query_id: str  # 查询ID
    query_text: str  # 查询文本
    
    # === 查询目标 ===
    target_backends: List[StorageBackend] = field(default_factory=list)
    # 示例: [StorageBackend.MILVUS, StorageBackend.POSTGRESQL]
    
    # === 查询策略 ===
    strategy: str = "parallel"  # parallel/sequential/cascade
    # parallel: 并行查询所有后端
    # sequential: 顺序查询，直到找到足够结果
    # cascade: 级联查询，先查缓存，miss后查向量库，再查数据库
    
    # === 结果聚合 ===
    merge_method: str = "union"  # union/intersect/rank_fusion
    max_results: int = 20  # 最大结果数
    
    # === 缓存控制 ===
    use_cache: bool = True
    cache_ttl: int = 3600  # 缓存TTL（秒）
    
    # === 超时控制 ===
    timeout: float = 5.0  # 超时时间（秒）
    
    # === 元数据 ===
    domain: Optional[DomainType] = None
    user_id: Optional[str] = None
    created_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "query_id": self.query_id,
            "query_text": self.query_text,
            "target_backends": [b.value for b in self.target_backends],
            "strategy": self.strategy,
            "merge_method": self.merge_method,
            "max_results": self.max_results,
            "use_cache": self.use_cache,
            "cache_ttl": self.cache_ttl,
            "timeout": self.timeout,
            "domain": self.domain.value if self.domain else None,
            "user_id": self.user_id,
            "created_at": self.created_at.isoformat()
        }


@dataclass
class MultiSourceResult:
    """
    多源查询结果
    
    聚合多个存储后端的查询结果
    """
    result_id: str  # 结果ID
    query_id: str  # 关联的查询ID
    
    # === 结果列表 ===
    items: List[Dict[str, Any]] = field(default_factory=list)
    # 统一格式的结果项，包含来源标记
    
    # === 来源分布 ===
    backend_distribution: Dict[str, int] = field(default_factory=dict)
    # key: backend名称, value: 该后端贡献的结果数
    
    # === 性能统计 ===
    total_time: float = 0.0  # 总耗时（秒）
    backend_times: Dict[str, float] = field(default_factory=dict)
    # key: backend名称, value: 该后端的查询耗时
    
    # === 缓存状态 ===
    cache_hit: bool = False
    cache_key: Optional[str] = None
    
    # === 元数据 ===
    retrieved_at: datetime = field(default_factory=datetime.now)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "result_id": self.result_id,
            "query_id": self.query_id,
            "items": self.items,
            "backend_distribution": self.backend_distribution,
            "total_time": self.total_time,
            "backend_times": self.backend_times,
            "cache_hit": self.cache_hit,
            "cache_key": self.cache_key,
            "retrieved_at": self.retrieved_at.isoformat()
        }
```

---

### 12.4 Milvus集合Schema

```python
# Milvus集合定义（Python代码）
from pymilvus import CollectionSchema, FieldSchema, DataType

def create_content_embeddings_collection():
    """
    创建内容嵌入集合
    
    存储所有内容的向量嵌入
    """
    fields = [
        FieldSchema(name="id", dtype=DataType.VARCHAR, is_primary=True, max_length=100),
        FieldSchema(name="item_id", dtype=DataType.VARCHAR, max_length=100),
        FieldSchema(name="source_id", dtype=DataType.VARCHAR, max_length=100),
        FieldSchema(name="domain", dtype=DataType.VARCHAR, max_length=50),
        FieldSchema(name="content_type", dtype=DataType.VARCHAR, max_length=50),
        FieldSchema(name="title", dtype=DataType.VARCHAR, max_length=500),
        FieldSchema(name="quality_score", dtype=DataType.FLOAT),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536),  # OpenAI ada-002
        FieldSchema(name="created_at", dtype=DataType.INT64),  # Unix timestamp
    ]
    
    schema = CollectionSchema(
        fields=fields,
        description="Content embeddings for vector similarity search",
        enable_dynamic_field=True  # 允许动态字段
    )
    
    # 索引配置
    index_params = {
        "metric_type": "L2",
        "index_type": "IVF_FLAT",
        "params": {"nlist": 1024}
    }
    
    return schema, index_params


def create_objective_embeddings_collection():
    """
    创建学习目标嵌入集合
    
    存储学习目标的向量嵌入
    """
    fields = [
        FieldSchema(name="id", dtype=DataType.VARCHAR, is_primary=True, max_length=100),
        FieldSchema(name="objective_id", dtype=DataType.VARCHAR, max_length=100),
        FieldSchema(name="domain", dtype=DataType.VARCHAR, max_length=50),
        FieldSchema(name="bloom_level", dtype=DataType.INT8),
        FieldSchema(name="description", dtype=DataType.VARCHAR, max_length=1000),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536),
        FieldSchema(name="created_at", dtype=DataType.INT64),
    ]
    
    schema = CollectionSchema(
        fields=fields,
        description="Learning objective embeddings"
    )
    
    index_params = {
        "metric_type": "COSINE",
        "index_type": "HNSW",
        "params": {"M": 16, "efConstruction": 200}
    }
    
    return schema, index_params
```

---

### 12.5 Redis数据结构

```python
# Redis键命名规范和数据结构

# 1. 查询结果缓存
# Key: cache:query:{query_hash}
# Type: String (JSON)
# TTL: 3600s
# Value: {"items": [...], "total": 10, "cached_at": "2025-12-09T..."}

# 2. 内容项缓存
# Key: cache:item:{item_id}
# Type: String (JSON)
# TTL: 7200s
# Value: ContentItem.to_dict()

# 3. 用户会话
# Key: session:{session_id}
# Type: Hash
# TTL: 1800s
# Fields: {user_id, domain, language, preferences}

# 4. 速率限制
# Key: ratelimit:{user_id}:{endpoint}
# Type: String (counter)
# TTL: 60s
# Value: request count

# 5. 分布式锁
# Key: lock:{resource_id}
# Type: String
# TTL: 30s
# Value: lock_token

# 6. 热门内容排行
# Key: trending:{domain}:{period}
# Type: Sorted Set
# Score: view_count
# Members: item_id

# 7. 向量搜索结果缓存
# Key: cache:vector:{embedding_hash}
# Type: String (JSON)
# TTL: 3600s
# Value: VectorSearchResult.to_dict()
```

---

### 12.6 MinIO桶结构

```yaml
# MinIO桶（Bucket）组织结构

buckets:
  # 1. 内容资源桶
  content-assets:
    description: "存储内容相关的多媒体资源"
    versioning: enabled
    encryption: AES256
    structure:
      - images/
          - k12/
          - art_history/
          - bestseller/
      - videos/
      - audio/
      - documents/
          - pdf/
          - docx/
    
  # 2. 用户上传桶
  user-uploads:
    description: "用户上传的文件"
    versioning: enabled
    lifecycle:
      - rule: delete_after_30_days
        prefix: temp/
    structure:
      - {user_id}/
          - avatars/
          - attachments/
    
  # 3. 生成内容桶
  generated-content:
    description: "AI生成的内容和叙述"
    versioning: enabled
    structure:
      - narratives/
          - {domain}/
              - {narrative_id}.json
              - {narrative_id}_assets/
      - lesson_plans/
      - worksheets/
    
  # 4. 备份桶
  backups:
    description: "数据库备份和导出"
    versioning: enabled
    lifecycle:
      - rule: transition_to_glacier_after_90_days
    structure:
      - database/
          - {date}/
      - exports/
          - {export_id}/
```

---

### 12.7 数据库表设计（PostgreSQL）

```sql
-- 对象存储元数据表（MinIO对象的PostgreSQL记录）
CREATE TABLE object_storage_metadata (
    object_id VARCHAR(100) PRIMARY KEY,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    
    -- 对象信息
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    etag VARCHAR(100) NOT NULL,
    
    -- 存储配置
    storage_class VARCHAR(50) DEFAULT 'STANDARD',
    encryption BOOLEAN DEFAULT false,
    
    -- 访问控制
    is_public BOOLEAN DEFAULT false,
    signed_url TEXT,
    signed_url_expires TIMESTAMP,
    
    -- 关联信息
    related_entity_type VARCHAR(50),
    related_entity_id VARCHAR(100),
    
    -- 元数据
    custom_metadata JSONB DEFAULT '{}'::jsonb,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_object_path UNIQUE (bucket_name, object_key)
);

CREATE INDEX idx_storage_bucket ON object_storage_metadata(bucket_name);
CREATE INDEX idx_storage_entity ON object_storage_metadata(related_entity_type, related_entity_id);
CREATE INDEX idx_storage_uploaded ON object_storage_metadata(uploaded_at);


-- 多源查询记录表
CREATE TABLE multi_source_queries (
    query_id VARCHAR(100) PRIMARY KEY,
    query_text TEXT NOT NULL,
    
    -- 查询目标
    target_backends JSONB NOT NULL,  -- ["milvus", "postgresql"]
    
    -- 查询策略
    strategy VARCHAR(20) DEFAULT 'parallel',
    merge_method VARCHAR(20) DEFAULT 'union',
    max_results INTEGER DEFAULT 20,
    
    -- 缓存控制
    use_cache BOOLEAN DEFAULT true,
    cache_ttl INTEGER DEFAULT 3600,
    
    -- 超时控制
    timeout FLOAT DEFAULT 5.0,
    
    -- 元数据
    domain VARCHAR(50),
    user_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_strategy CHECK (strategy IN ('parallel', 'sequential', 'cascade')),
    CONSTRAINT check_merge_method CHECK (merge_method IN ('union', 'intersect', 'rank_fusion'))
);

CREATE INDEX idx_multi_queries_user ON multi_source_queries(user_id);
CREATE INDEX idx_multi_queries_domain ON multi_source_queries(domain);
CREATE INDEX idx_multi_queries_created ON multi_source_queries(created_at);


-- 多源查询结果表
CREATE TABLE multi_source_results (
    result_id VARCHAR(100) PRIMARY KEY,
    query_id VARCHAR(100) NOT NULL,
    
    -- 结果列表（item_id数组）
    item_ids JSONB NOT NULL,
    
    -- 来源分布
    backend_distribution JSONB DEFAULT '{}'::jsonb,  -- {"milvus": 5, "postgresql": 3}
    
    -- 性能统计
    total_time FLOAT DEFAULT 0.0,
    backend_times JSONB DEFAULT '{}'::jsonb,  -- {"milvus": 0.23, "postgresql": 0.15}
    
    -- 缓存状态
    cache_hit BOOLEAN DEFAULT false,
    cache_key VARCHAR(200),
    
    -- 元数据
    retrieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_multi_result_query FOREIGN KEY (query_id) 
        REFERENCES multi_source_queries(query_id) ON DELETE CASCADE
);

CREATE INDEX idx_multi_results_query ON multi_source_results(query_id);
CREATE INDEX idx_multi_results_cache ON multi_source_results(cache_hit);
CREATE INDEX idx_multi_results_retrieved ON multi_source_results(retrieved_at);
```

---

## 13. 数据库Schema

### 5.1 PostgreSQL表结构

```sql
-- 工作流定义表
CREATE TABLE workflow_definitions (
    workflow_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    version VARCHAR(20) NOT NULL,
    description TEXT,
    
    trigger_conditions JSONB NOT NULL,
    global_config JSONB NOT NULL,
    nodes JSONB NOT NULL,
    output_config JSONB NOT NULL,
    
    enabled BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(name, version)
);

-- 创建索引
CREATE INDEX idx_workflow_trigger ON workflow_definitions 
    USING gin(trigger_conditions);

-- 工作流实例表
CREATE TABLE workflow_instances (
    instance_id VARCHAR(100) PRIMARY KEY,
    workflow_id VARCHAR(100) NOT NULL REFERENCES workflow_definitions(workflow_id),
    
    status VARCHAR(20) NOT NULL,
    
    context JSONB NOT NULL,
    node_executions JSONB NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    error_log TEXT[],
    metrics JSONB DEFAULT '{}'
);

-- 创建索引
CREATE INDEX idx_instance_workflow ON workflow_instances(workflow_id);
CREATE INDEX idx_instance_status ON workflow_instances(status);
CREATE INDEX idx_instance_created ON workflow_instances(created_at);

-- 规则表
CREATE TABLE rules (
    rule_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    priority INTEGER NOT NULL,
    
    conditions JSONB NOT NULL,
    actions JSONB NOT NULL,
    
    enabled BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rules_priority ON rules(priority DESC);
CREATE INDEX idx_rules_conditions ON rules USING gin(conditions);

-- 生成内容表
CREATE TABLE generated_contents (
    content_id VARCHAR(100) PRIMARY KEY,
    content_type VARCHAR(50) NOT NULL,
    
    subject VARCHAR(50) NOT NULL,
    grade INTEGER NOT NULL,
    region VARCHAR(20) NOT NULL,
    topics TEXT[],
    
    workflow_id VARCHAR(100) REFERENCES workflow_definitions(workflow_id),
    instance_id VARCHAR(100) REFERENCES workflow_instances(instance_id),
    
    data JSONB NOT NULL,
    metadata JSONB NOT NULL,
    
    review_status VARCHAR(20) DEFAULT 'pending',
    reviewer_feedback TEXT,
    
    view_count INTEGER DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    avg_rating DECIMAL(3,2),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_content_subject_grade ON generated_contents(subject, grade);
CREATE INDEX idx_content_region ON generated_contents(region);
CREATE INDEX idx_content_topics ON generated_contents USING gin(topics);
CREATE INDEX idx_content_created ON generated_contents(created_at);

-- 用户画像表
CREATE TABLE personas (
    persona_id VARCHAR(100) PRIMARY KEY,
    target_group VARCHAR(200) NOT NULL,
    
    cognitive JSONB NOT NULL,
    emotional JSONB NOT NULL,
    behavioral JSONB NOT NULL,
    
    recommendations JSONB NOT NULL,
    key_considerations TEXT[],
    
    usage_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 审计日志表
CREATE TABLE audit_logs (
    log_id BIGSERIAL PRIMARY KEY,
    
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(100),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id VARCHAR(100),
    
    ip_address INET,
    user_agent TEXT,
    
    details JSONB,
    result VARCHAR(20)
);

CREATE INDEX idx_audit_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);

-- 工作流模板表 ⭐ 新增
CREATE TABLE workflow_templates (
    template_id VARCHAR(100) PRIMARY KEY,
    template_name VARCHAR(200) NOT NULL,
    version VARCHAR(20) NOT NULL,
    base_template_id VARCHAR(100) REFERENCES workflow_templates(template_id),
    
    -- 元数据
    metadata JSONB NOT NULL,
    -- 包含: subject, grade_range, region, created_by, created_at, description, 
    --      source_instance_id, tags, is_public, fork_count, popularity_score
    
    -- 工作流结构
    nodes JSONB NOT NULL,
    edges JSONB NOT NULL,
    
    -- 参数定义
    parameters JSONB NOT NULL,
    -- 包含: required, optional, defaults, constraints
    
    -- 执行统计
    execution_history JSONB DEFAULT '{}',
    -- 包含: total_runs, success_rate, avg_duration, avg_cost, 
    --      avg_quality_score, last_run_at
    
    -- 用户反馈
    ratings JSONB DEFAULT '{"avg_score": 0, "total_ratings": 0, "rating_distribution": {}}',
    
    -- 变更记录
    changelog JSONB DEFAULT '[]',
    
    -- 优化建议
    optimization_suggestions JSONB DEFAULT '[]',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(template_name, version)
);

-- 创建索引
CREATE INDEX idx_template_name ON workflow_templates(template_name);
CREATE INDEX idx_template_subject ON workflow_templates 
    USING gin((metadata->'subject'));
CREATE INDEX idx_template_grade ON workflow_templates 
    USING gin((metadata->'grade_range'));
CREATE INDEX idx_template_public ON workflow_templates 
    ((metadata->>'is_public')) 
    WHERE (metadata->>'is_public')::boolean = true;
CREATE INDEX idx_template_popularity ON workflow_templates 
    (((metadata->>'popularity_score')::numeric)) DESC;

-- 模板执行记录表 ⭐ 新增
CREATE TABLE template_execution_records (
    record_id VARCHAR(100) PRIMARY KEY,
    instance_id VARCHAR(100) NOT NULL REFERENCES workflow_instances(instance_id),
    template_id VARCHAR(100) NOT NULL REFERENCES workflow_templates(template_id),
    template_version VARCHAR(20) NOT NULL,
    
    executed_at TIMESTAMP NOT NULL,
    executed_by VARCHAR(100),
    
    -- 执行轨迹
    execution_trace JSONB NOT NULL,
    -- 包含: nodes (每个节点的详细执行信息), total_duration, success, error_nodes
    
    -- 资源消耗
    resource_usage JSONB NOT NULL,
    -- 包含: total_tokens, total_cost, models_used, cache_hits, cache_misses
    
    -- 质量指标
    quality_metrics JSONB NOT NULL,
    -- 包含: overall_score, content_quality, logic_consistency, 
    --      factual_accuracy, compliance_score, user_rating
    
    -- 用户反馈
    user_feedback JSONB,
    -- 包含: rating, comments, issues, submitted_at
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_exec_record_template ON template_execution_records(template_id);
CREATE INDEX idx_exec_record_time ON template_execution_records(executed_at);
CREATE INDEX idx_exec_record_user ON template_execution_records(executed_by);
CREATE INDEX idx_exec_record_quality ON template_execution_records 
    (((quality_metrics->>'overall_score')::numeric)) DESC;

-- 优化建议表 ⭐ 新增
CREATE TABLE optimization_suggestions (
    suggestion_id VARCHAR(100) PRIMARY KEY,
    template_id VARCHAR(100) NOT NULL REFERENCES workflow_templates(template_id),
    
    suggestion_type VARCHAR(50) NOT NULL, -- performance/quality/cost/structure
    severity VARCHAR(20) NOT NULL, -- high/medium/low
    
    issue TEXT NOT NULL,
    affected_nodes JSONB DEFAULT '[]',
    
    suggestion TEXT NOT NULL,
    expected_improvement JSONB DEFAULT '{}',
    supporting_data JSONB DEFAULT '{}',
    
    status VARCHAR(20) DEFAULT 'pending', -- pending/applied/rejected/testing
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    applied_at TIMESTAMP,
    applied_by VARCHAR(100),
    result JSONB
);

-- 创建索引
CREATE INDEX idx_suggestion_template ON optimization_suggestions(template_id);
CREATE INDEX idx_suggestion_status ON optimization_suggestions(status);
CREATE INDEX idx_suggestion_severity ON optimization_suggestions(severity);
CREATE INDEX idx_suggestion_type ON optimization_suggestions(suggestion_type);

-- 模板版本链表 ⭐ 新增
CREATE TABLE template_versions (
    id BIGSERIAL PRIMARY KEY,
    template_id VARCHAR(100) NOT NULL REFERENCES workflow_templates(template_id),
    version VARCHAR(20) NOT NULL,
    parent_version VARCHAR(20),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(template_id, version)
);

CREATE INDEX idx_version_template ON template_versions(template_id);

-- 实例-模板关联表 ⭐ 新增
CREATE TABLE instance_template_mapping (
    id BIGSERIAL PRIMARY KEY,
    instance_id VARCHAR(100) NOT NULL REFERENCES workflow_instances(instance_id),
    template_id VARCHAR(100) NOT NULL REFERENCES workflow_templates(template_id),
    relationship VARCHAR(50) NOT NULL, -- source/derived_from/fork
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(instance_id, template_id)
);

CREATE INDEX idx_itm_instance ON instance_template_mapping(instance_id);
CREATE INDEX idx_itm_template ON instance_template_mapping(template_id);

-- 性能指标表
CREATE TABLE metrics (
    metric_id BIGSERIAL PRIMARY KEY,
    
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(10,2),
    
    labels JSONB,
    
    instance_id VARCHAR(100),
    node_id VARCHAR(100)
);

CREATE INDEX idx_metrics_timestamp ON metrics(timestamp);
CREATE INDEX idx_metrics_name ON metrics(metric_name);
```

### 5.2 Redis数据结构

```python
# 工作流实例状态 (Hash)
# Key: workflow:instance:{instance_id}
{
    "status": "running",
    "current_node": "content_generation",
    "progress": "0.6",
    "created_at": "2025-12-09T10:00:00Z"
}

# 节点输出缓存 (Hash)
# Key: workflow:instance:{instance_id}:outputs
{
    "persona": "{...}",
    "objectives": "{...}",
    "rag_retrieval": "{...}"
}

# 执行队列 (List)
# Key: workflow:queue:pending
["instance_id_1", "instance_id_2", "instance_id_3"]

# 模型路由负载 (Sorted Set)
# Key: llm:load
# Score: 当前QPS
{
    "gpt-4": 45.0,
    "claude-3": 30.0,
    "deepseek-v2": 20.0
}

# 配置缓存 (String with TTL)
# Key: config:workflow:{workflow_id}
# TTL: 3600 seconds
"{...workflow definition...}"

# 模板缓存 ⭐ 新增 (String with TTL)
# Key: template:{template_id}
# TTL: 3600 seconds
"{...workflow template...}"

# 热门模板 ⭐ 新增 (Sorted Set)
# Key: templates:hot
# Score: 流行度评分
{
    "template_id_1": 9.5,
    "template_id_2": 8.8,
    "template_id_3": 7.2
}

# 用户模板使用历史 ⭐ 新增 (List)
# Key: user:{user_id}:templates
["template_id_1", "template_id_2", "template_id_3"]

# 模板执行统计 ⭐ 新增 (Hash)
# Key: template:{template_id}:stats
{
    "total_runs": "152",
    "success_rate": "0.96",
    "avg_duration": "120.5",
    "last_run_at": "2025-01-20T14:30:00"
}

# 分布式锁 (String with TTL)
# Key: lock:instance:{instance_id}
# TTL: 300 seconds
"locked_by_worker_1"
```

---

### 8.3 新增表结构 ⭐

#### 课程规划相关表

```sql
-- ==================== 课程规划相关表 ====================

-- 课程体系表
CREATE TABLE curriculums (
    curriculum_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    grade INTEGER NOT NULL CHECK (grade >= 1 AND grade <= 12),
    region VARCHAR(10) NOT NULL,
    
    level VARCHAR(20) NOT NULL,  -- SEMESTER/ACADEMIC_YEAR/TEXTBOOK
    textbook_version VARCHAR(100),
    
    -- 时间规划
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_weeks INTEGER NOT NULL,
    lessons_per_week INTEGER NOT NULL,
    total_lessons INTEGER NOT NULL,
    
    -- 课程结构
    course_units JSONB DEFAULT '[]'::jsonb,  -- 单元ID数组
    unit_dependencies JSONB DEFAULT '{}'::jsonb,  -- 依赖关系
    
    -- 故事化设计
    story_universe_id VARCHAR(100),
    narrative_structure VARCHAR(50),
    
    -- 知识体系
    knowledge_graph_id VARCHAR(100),
    core_competencies JSONB DEFAULT '[]'::jsonb,
    curriculum_alignment JSONB DEFAULT '{}'::jsonb,
    
    -- 元数据
    status VARCHAR(20) DEFAULT 'DRAFT',
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 统计数据
    progress JSONB DEFAULT '{}'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_story_universe FOREIGN KEY (story_universe_id) 
        REFERENCES story_universes(universe_id) ON DELETE SET NULL,
    CONSTRAINT fk_knowledge_graph FOREIGN KEY (knowledge_graph_id) 
        REFERENCES knowledge_graphs(graph_id) ON DELETE SET NULL
);

CREATE INDEX idx_curriculum_subject_grade ON curriculums(subject, grade);
CREATE INDEX idx_curriculum_region ON curriculums(region);
CREATE INDEX idx_curriculum_status ON curriculums(status);
CREATE INDEX idx_curriculum_created_by ON curriculums(created_by);


-- 课程单元表
CREATE TABLE course_units (
    unit_id VARCHAR(100) PRIMARY KEY,
    curriculum_id VARCHAR(100) NOT NULL,
    
    sequence_number INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    subtitle VARCHAR(200),
    
    -- 教学内容
    topics JSONB DEFAULT '[]'::jsonb,
    learning_objectives JSONB DEFAULT '[]'::jsonb,
    driving_questions JSONB DEFAULT '[]'::jsonb,
    
    -- 课时规划
    total_lessons INTEGER DEFAULT 0,
    lesson_metas JSONB DEFAULT '[]'::jsonb,  -- LessonMeta ID数组
    lesson_distribution JSONB DEFAULT '{}'::jsonb,
    
    -- 时间安排
    estimated_weeks FLOAT DEFAULT 0.0,
    start_week INTEGER,
    end_week INTEGER,
    
    -- 故事化元素
    story_arc_id VARCHAR(100),
    narrative_theme VARCHAR(200),
    
    -- 依赖关系
    prerequisite_units JSONB DEFAULT '[]'::jsonb,
    prerequisite_knowledge JSONB DEFAULT '[]'::jsonb,
    
    -- 资源
    resources JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    status VARCHAR(20) DEFAULT 'planned',
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_curriculum FOREIGN KEY (curriculum_id) 
        REFERENCES curriculums(curriculum_id) ON DELETE CASCADE,
    CONSTRAINT unique_curriculum_sequence UNIQUE (curriculum_id, sequence_number)
);

CREATE INDEX idx_unit_curriculum ON course_units(curriculum_id);
CREATE INDEX idx_unit_status ON course_units(status);


-- 课时元信息表
CREATE TABLE lesson_metas (
    lesson_meta_id VARCHAR(100) PRIMARY KEY,
    unit_id VARCHAR(100) NOT NULL,
    curriculum_id VARCHAR(100) NOT NULL,
    
    sequence_number INTEGER NOT NULL,
    global_sequence INTEGER NOT NULL,
    
    -- 课时定位
    title VARCHAR(200) NOT NULL,
    lesson_type VARCHAR(50) NOT NULL,
    
    -- 教学内容
    topics JSONB DEFAULT '[]'::jsonb,
    objectives JSONB DEFAULT '[]'::jsonb,
    key_points JSONB DEFAULT '[]'::jsonb,
    difficult_points JSONB DEFAULT '[]'::jsonb,
    
    -- 时间规划
    duration INTEGER DEFAULT 45,
    scheduled_week INTEGER,
    scheduled_date DATE,
    
    -- 故事化元素
    story_episode_id VARCHAR(100),
    narrative_hook TEXT,
    
    -- 依赖关系
    prerequisite_lessons JSONB DEFAULT '[]'::jsonb,
    
    -- 生成配置
    workflow_template_id VARCHAR(100),
    generation_config JSONB DEFAULT '{}'::jsonb,
    
    -- 生成状态
    generated BOOLEAN DEFAULT false,
    lesson_script_id VARCHAR(100),
    generation_quality_score FLOAT,
    
    -- 元数据
    status VARCHAR(20) DEFAULT 'planned',
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_unit FOREIGN KEY (unit_id) 
        REFERENCES course_units(unit_id) ON DELETE CASCADE,
    CONSTRAINT fk_lesson_curriculum FOREIGN KEY (curriculum_id) 
        REFERENCES curriculums(curriculum_id) ON DELETE CASCADE,
    CONSTRAINT unique_unit_sequence UNIQUE (unit_id, sequence_number),
    CONSTRAINT unique_curriculum_global_sequence UNIQUE (curriculum_id, global_sequence)
);

CREATE INDEX idx_lesson_meta_unit ON lesson_metas(unit_id);
CREATE INDEX idx_lesson_meta_curriculum ON lesson_metas(curriculum_id);
CREATE INDEX idx_lesson_meta_status ON lesson_metas(status);
CREATE INDEX idx_lesson_meta_generated ON lesson_metas(generated);
CREATE INDEX idx_lesson_meta_scheduled_date ON lesson_metas(scheduled_date);


-- 故事宇宙表
CREATE TABLE story_universes (
    universe_id VARCHAR(100) PRIMARY KEY,
    curriculum_id VARCHAR(100) NOT NULL,
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- 世界观设定
    world_setting TEXT,
    theme VARCHAR(100),
    tone VARCHAR(50),
    
    -- 角色体系
    characters JSONB DEFAULT '[]'::jsonb,
    main_character_id VARCHAR(100),
    
    -- 故事结构
    story_arcs JSONB DEFAULT '[]'::jsonb,
    episodes JSONB DEFAULT '[]'::jsonb,
    
    -- 跨单元连续性
    continuity_elements JSONB DEFAULT '[]'::jsonb,
    
    -- 可视化资源
    visual_assets JSONB DEFAULT '{}'::jsonb,
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_universe_curriculum FOREIGN KEY (curriculum_id) 
        REFERENCES curriculums(curriculum_id) ON DELETE CASCADE,
    CONSTRAINT unique_curriculum_universe UNIQUE (curriculum_id)
);

CREATE INDEX idx_story_universe_curriculum ON story_universes(curriculum_id);


-- 知识图谱表
CREATE TABLE knowledge_graphs (
    graph_id VARCHAR(100) PRIMARY KEY,
    curriculum_id VARCHAR(100) NOT NULL,
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- 图谱结构（存储为JSON）
    nodes JSONB DEFAULT '[]'::jsonb,
    edges JSONB DEFAULT '[]'::jsonb,
    
    -- 分层信息
    levels JSONB DEFAULT '{}'::jsonb,
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_graph_curriculum FOREIGN KEY (curriculum_id) 
        REFERENCES curriculums(curriculum_id) ON DELETE CASCADE,
    CONSTRAINT unique_curriculum_graph UNIQUE (curriculum_id)
);

CREATE INDEX idx_knowledge_graph_curriculum ON knowledge_graphs(curriculum_id);
```

#### 模板本地化相关表

```sql
-- ==================== 模板本地化相关表 ====================

-- 本地化配置表
CREATE TABLE localization_configs (
    localization_id VARCHAR(100) PRIMARY KEY,
    job_name VARCHAR(200) NOT NULL,
    
    -- 源和目标
    source_template_id VARCHAR(100) NOT NULL,
    source_region VARCHAR(10) NOT NULL,
    source_language VARCHAR(10) NOT NULL,
    
    target_region VARCHAR(10) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    
    -- 本地化策略
    localization_strategy JSONB DEFAULT '{}'::jsonb,
    
    -- 适配范围
    adaptation_scopes JSONB DEFAULT '[]'::jsonb,
    
    -- 标准映射
    standard_mappings JSONB DEFAULT '[]'::jsonb,  -- StandardMapping ID数组
    
    -- 文化适配记录
    cultural_adaptations JSONB DEFAULT '[]'::jsonb,
    
    -- 教学法调整
    pedagogy_adjustments JSONB DEFAULT '[]'::jsonb,
    
    -- 合规性检查
    compliance_issues JSONB DEFAULT '[]'::jsonb,  -- ComplianceIssue ID数组
    
    -- 质量控制
    quality_report_id VARCHAR(100),
    automation_rate FLOAT DEFAULT 0.0,
    human_review_points JSONB DEFAULT '[]'::jsonb,
    
    -- 执行状态
    status VARCHAR(20) DEFAULT 'PENDING',
    current_phase VARCHAR(50),
    progress FLOAT DEFAULT 0.0,
    
    -- 输出
    target_template_id VARCHAR(100),
    
    -- 元数据
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    
    estimated_duration INTEGER,  -- 小时
    actual_duration INTEGER,  -- 小时
    
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_source_template FOREIGN KEY (source_template_id) 
        REFERENCES workflow_templates(template_id) ON DELETE RESTRICT,
    CONSTRAINT fk_target_template FOREIGN KEY (target_template_id) 
        REFERENCES workflow_templates(template_id) ON DELETE SET NULL,
    CONSTRAINT fk_quality_report FOREIGN KEY (quality_report_id) 
        REFERENCES quality_reports(report_id) ON DELETE SET NULL
);

CREATE INDEX idx_localization_status ON localization_configs(status);
CREATE INDEX idx_localization_source ON localization_configs(source_template_id);
CREATE INDEX idx_localization_target ON localization_configs(target_template_id);
CREATE INDEX idx_localization_regions ON localization_configs(source_region, target_region);


-- 课标映射表
CREATE TABLE standard_mappings (
    mapping_id VARCHAR(100) PRIMARY KEY,
    localization_id VARCHAR(100) NOT NULL,
    
    -- 源课标
    source_standard_id VARCHAR(100) NOT NULL,
    source_standard_name VARCHAR(200) NOT NULL,
    source_description TEXT,
    source_region VARCHAR(10) NOT NULL,
    
    -- 目标课标
    target_standard_id VARCHAR(100) NOT NULL,
    target_standard_name VARCHAR(200) NOT NULL,
    target_description TEXT,
    target_region VARCHAR(10) NOT NULL,
    
    -- 映射关系
    mapping_type VARCHAR(20) NOT NULL,  -- exact/partial/approximate/no_match
    similarity_score FLOAT DEFAULT 0.0,
    
    -- 差异说明
    differences JSONB DEFAULT '[]'::jsonb,
    
    -- 调整建议
    adaptation_suggestions JSONB DEFAULT '[]'::jsonb,
    
    -- 元数据
    confidence FLOAT DEFAULT 0.0,
    verified_by_human BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_mapping_localization FOREIGN KEY (localization_id) 
        REFERENCES localization_configs(localization_id) ON DELETE CASCADE
);

CREATE INDEX idx_standard_mapping_localization ON standard_mappings(localization_id);
CREATE INDEX idx_standard_mapping_regions ON standard_mappings(source_region, target_region);
CREATE INDEX idx_standard_mapping_type ON standard_mappings(mapping_type);


-- 合规性问题表
CREATE TABLE compliance_issues (
    issue_id VARCHAR(100) PRIMARY KEY,
    localization_id VARCHAR(100) NOT NULL,
    
    -- 问题分类
    category VARCHAR(50) NOT NULL,  -- legal/cultural/educational/safety/privacy
    severity VARCHAR(20) NOT NULL,  -- critical/high/medium/low
    
    -- 问题描述
    title VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- 问题位置
    location JSONB DEFAULT '{}'::jsonb,
    
    -- 违规内容
    flagged_content TEXT,
    violated_rule TEXT NOT NULL,
    rule_reference VARCHAR(500),
    
    -- 建议修改
    suggested_fix TEXT,
    alternative_content TEXT,
    
    -- 处理状态
    status VARCHAR(20) DEFAULT 'open',  -- open/fixed/wont_fix/false_positive
    resolution TEXT,
    resolved_by VARCHAR(100),
    resolved_at TIMESTAMP,
    
    -- 元数据
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    detector VARCHAR(50) NOT NULL,  -- auto/human
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_issue_localization FOREIGN KEY (localization_id) 
        REFERENCES localization_configs(localization_id) ON DELETE CASCADE
);

CREATE INDEX idx_compliance_issue_localization ON compliance_issues(localization_id);
CREATE INDEX idx_compliance_issue_category ON compliance_issues(category);
CREATE INDEX idx_compliance_issue_severity ON compliance_issues(severity);
CREATE INDEX idx_compliance_issue_status ON compliance_issues(status);


-- 质量报告表
CREATE TABLE quality_reports (
    report_id VARCHAR(100) PRIMARY KEY,
    localization_id VARCHAR(100) NOT NULL,
    
    -- 总体评分
    overall_score FLOAT DEFAULT 0.0,
    grade VARCHAR(5) DEFAULT '',
    
    -- 分项评分（存储为JSON数组）
    metrics JSONB DEFAULT '[]'::jsonb,
    
    -- 问题统计
    issue_summary JSONB DEFAULT '{}'::jsonb,
    
    -- 优点与改进
    strengths JSONB DEFAULT '[]'::jsonb,
    weaknesses JSONB DEFAULT '[]'::jsonb,
    recommendations JSONB DEFAULT '[]'::jsonb,
    
    -- 人工审核
    human_review_required BOOLEAN DEFAULT false,
    human_review_points JSONB DEFAULT '[]'::jsonb,
    
    -- 测试结果
    test_results JSONB,
    
    -- 元数据
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    generated_by VARCHAR(50) NOT NULL,  -- auto/human
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_report_localization FOREIGN KEY (localization_id) 
        REFERENCES localization_configs(localization_id) ON DELETE CASCADE
);

CREATE INDEX idx_quality_report_localization ON quality_reports(localization_id);
CREATE INDEX idx_quality_report_grade ON quality_reports(grade);
```

#### 组件扩展相关表

```sql
-- ==================== 组件扩展相关表 ====================

-- 节点元数据表
CREATE TABLE node_metadata (
    node_id VARCHAR(100) PRIMARY KEY,
    node_name VARCHAR(200) NOT NULL,
    version VARCHAR(20) NOT NULL,
    
    -- 分类
    category VARCHAR(50) NOT NULL,
    subcategory VARCHAR(50),
    tags JSONB DEFAULT '[]'::jsonb,
    
    -- 描述
    description TEXT NOT NULL,
    long_description TEXT,
    
    -- 适用范围
    applicable_subjects JSONB DEFAULT '[]'::jsonb,
    applicable_grades VARCHAR(20) DEFAULT '',
    applicable_lesson_types JSONB DEFAULT '[]'::jsonb,
    
    -- 作者信息
    author VARCHAR(100) NOT NULL,
    author_email VARCHAR(100),
    organization VARCHAR(200),
    license VARCHAR(50) DEFAULT 'MIT',
    
    -- 依赖项
    dependencies JSONB DEFAULT '[]'::jsonb,
    
    -- 性能指标
    estimated_duration INTEGER DEFAULT 10,  -- 秒
    memory_limit INTEGER DEFAULT 512,  -- MB
    cpu_limit FLOAT DEFAULT 1.0,  -- 核数
    
    -- 输入输出Schema
    input_schema JSONB DEFAULT '{}'::jsonb,
    output_schema JSONB DEFAULT '{}'::jsonb,
    
    -- 配置Schema
    config_schema JSONB DEFAULT '{}'::jsonb,
    
    -- 文档
    documentation_url VARCHAR(500),
    example_url VARCHAR(500),
    
    -- 版本信息
    changelog JSONB DEFAULT '[]'::jsonb,
    
    -- 统计数据
    download_count INTEGER DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    avg_rating FLOAT DEFAULT 0.0,
    
    -- 认证
    certification_level VARCHAR(20),  -- bronze/silver/gold/platinum
    verified BOOLEAN DEFAULT false,
    
    -- 元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deprecated BOOLEAN DEFAULT false,
    deprecation_message TEXT,
    
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT unique_node_version UNIQUE (node_id, version)
);

CREATE INDEX idx_node_category ON node_metadata(category);
CREATE INDEX idx_node_tags ON node_metadata USING gin(tags);
CREATE INDEX idx_node_subjects ON node_metadata USING gin(applicable_subjects);
CREATE INDEX idx_node_certification ON node_metadata(certification_level);
CREATE INDEX idx_node_verified ON node_metadata(verified);
CREATE INDEX idx_node_deprecated ON node_metadata(deprecated);


-- 组件包表
CREATE TABLE component_packages (
    package_id VARCHAR(100) PRIMARY KEY,
    node_id VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    
    -- 包信息
    package_name VARCHAR(200) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    
    -- 文件信息
    files JSONB DEFAULT '{}'::jsonb,
    
    -- 包大小
    package_size BIGINT DEFAULT 0,  -- 字节
    
    -- 下载信息
    download_url VARCHAR(500),
    checksum VARCHAR(100),  -- MD5/SHA256
    
    -- 定价
    pricing_model VARCHAR(20) DEFAULT 'free',  -- free/freemium/paid
    free_quota INTEGER DEFAULT 0,
    price_per_call FLOAT DEFAULT 0.0,
    
    -- 元数据
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_package_node FOREIGN KEY (node_id) 
        REFERENCES node_metadata(node_id) ON DELETE CASCADE,
    CONSTRAINT unique_package_version UNIQUE (node_id, version)
);

CREATE INDEX idx_component_package_node ON component_packages(node_id);
CREATE INDEX idx_component_package_pricing ON component_packages(pricing_model);


-- 组件使用记录表
CREATE TABLE component_usage_records (
    record_id VARCHAR(100) PRIMARY KEY,
    node_id VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    
    -- 使用信息
    user_id VARCHAR(100) NOT NULL,
    workflow_instance_id VARCHAR(100) NOT NULL,
    
    -- 执行信息
    executed_at TIMESTAMP NOT NULL,
    duration FLOAT DEFAULT 0.0,  -- 秒
    memory_used FLOAT DEFAULT 0.0,  -- MB
    cpu_used FLOAT DEFAULT 0.0,  -- 核·秒
    
    -- 结果
    status VARCHAR(20) DEFAULT 'success',  -- success/failure/timeout/error
    error_message TEXT,
    
    -- 降级信息
    fallback_activated BOOLEAN DEFAULT false,
    fallback_reason VARCHAR(200),
    
    -- 质量评分
    output_quality_score FLOAT,
    
    -- 用户反馈
    user_rating INTEGER CHECK (user_rating >= 1 AND user_rating <= 5),
    user_comment TEXT,
    
    -- 计费
    billable BOOLEAN DEFAULT false,
    cost FLOAT DEFAULT 0.0,
    
    -- 元数据
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT fk_usage_node FOREIGN KEY (node_id) 
        REFERENCES node_metadata(node_id) ON DELETE CASCADE,
    CONSTRAINT fk_usage_workflow FOREIGN KEY (workflow_instance_id) 
        REFERENCES execution_instances(instance_id) ON DELETE CASCADE
);

CREATE INDEX idx_usage_node ON component_usage_records(node_id);
CREATE INDEX idx_usage_user ON component_usage_records(user_id);
CREATE INDEX idx_usage_workflow ON component_usage_records(workflow_instance_id);
CREATE INDEX idx_usage_executed_at ON component_usage_records(executed_at);
CREATE INDEX idx_usage_status ON component_usage_records(status);
CREATE INDEX idx_usage_billable ON component_usage_records(billable);
```

---

## 附录

### A. 数据验证Schema (Pydantic)

```python
from pydantic import BaseModel, Field, validator
from typing import List, Optional
from datetime import datetime

class WorkflowRequestSchema(BaseModel):
    """工作流请求验证"""
    subject: str = Field(..., description="学科")
    grade: int = Field(..., ge=1, le=12, description="年级")
    region: str = Field(..., description="地区")
    topics: List[str] = Field(..., min_items=1, description="知识点列表")
    
    lesson_type: Optional[str] = "新授课"
    duration: Optional[int] = Field(45, ge=15, le=90)
    
    @validator('subject')
    def validate_subject(cls, v):
        allowed = ["数学", "语文", "英语", "科学", "历史", "物理", "化学", "生物"]
        if v not in allowed:
            raise ValueError(f"不支持的学科: {v}")
        return v
```

---

**文档结束**
