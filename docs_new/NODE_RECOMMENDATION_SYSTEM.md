# Coze Studio 节点推荐系统 - 详细实施方案

**文档版本**: V1.1 🆕  
**日期**: 2025-12-12  
**作者**: Coze Studio 开发团队  
**文档类型**: 技术设计与实施指南

**🆕 V1.1 更新亮点**：新增"从零开始的智能工作流构建"功能，支持通过自然语言描述需求自动生成完整工作流！

---

## 📋 目录

1. [系统概述](#1-系统概述)
2. [核心功能设计](#2-核心功能设计)
3. [技术架构](#3-技术架构)
4. [实施方案](#4-实施方案)
5. [代码示例](#5-代码示例)
6. [部署指南](#6-部署指南)
7. [最佳实践](#7-最佳实践)

---

## 1. 系统概述

### 1.1 需求背景

在 Coze Studio 工作流编辑器中，用户需要智能推荐功能来：

1. **场景化推荐**：根据当前工作流的上下文，推荐下一个最合适的节点
2. **智能补全**：基于用户输入、已有节点类型和数据流，自动推荐后续节点
3. **最佳实践**：推荐常见工作流模式和节点组合
4. **提升效率**：减少用户寻找和选择节点的时间

### 1.2 系统定位

**节点推荐系统 = LLM 智能推理 + 代码规则引擎 + 历史数据分析**

```
┌─────────────────────────────────────────────────────────────┐
│                   节点推荐系统架构                            │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │ LLM推理层  │  │ 规则引擎层  │  │ 数据分析层  │            │
│  │            │  │            │  │            │            │
│  │ • 语义理解 │  │ • 节点规则 │  │ • 使用统计 │            │
│  │ • 意图识别 │  │ • 类型约束 │  │ • 模式挖掘 │            │
│  │ • 上下文   │  │ • 最佳实践 │  │ • 协同过滤 │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│         │                │                │                 │
│         └────────────────┴────────────────┘                 │
│                         ↓                                   │
│              ┌──────────────────────┐                       │
│              │   推荐结果聚合器      │                       │
│              └──────────────────────┘                       │
│                         ↓                                   │
│              ┌──────────────────────┐                       │
│              │   排序 & 过滤        │                       │
│              └──────────────────────┘                       │
└─────────────────────────────────────────────────────────────┘
```

### 1.3 功能特性

| 功能 | 说明 | 实现方式 |
|------|------|----------|
| **智能续写** | 根据最后一个节点推荐下一个节点 | LLM + 规则引擎 |
| **模式识别** | 识别常见工作流模式 | 数据分析 + 规则库 |
| **类型推断** | 基于数据流类型推荐兼容节点 | 类型系统 + 规则引擎 |
| **上下文感知** | 考虑整个工作流的上下文 | LLM 语义分析 |
| **个性化推荐** | 基于用户历史习惯 | 协同过滤 |
| **实时反馈** | 用户选择后动态调整推荐 | 在线学习 |

---

## 2. 核心功能设计

### 2.1 推荐触发场景

#### 场景 1：连线推荐
用户从某个节点拖出连线时，显示推荐节点列表

```go
// 伪代码示例
type NodeRecommendationRequest struct {
    SourceNodeID   string                    // 源节点 ID
    SourceNodeType entity.NodeType           // 源节点类型
    SourceOutputs  map[string]*TypeInfo      // 源节点输出类型
    WorkflowContext *WorkflowContext         // 工作流上下文
    UserID         string                    // 用户 ID（用于个性化）
}

type NodeRecommendationResponse struct {
    Recommendations []RecommendedNode `json:"recommendations"`
    Reasoning       string            `json:"reasoning"`  // 推荐理由
}

type RecommendedNode struct {
    NodeType     entity.NodeType `json:"node_type"`
    DisplayName  string          `json:"display_name"`
    Score        float64         `json:"score"`         // 推荐得分 0-1
    Reason       string          `json:"reason"`        // 推荐原因
    Category     string          `json:"category"`      // 节点分类
    Icon         string          `json:"icon"`
    Template     *vo.Node        `json:"template"`      // 节点模板（可选）
}
```

#### 场景 2：空白画布推荐
用户在空白画布添加第一个节点后的推荐

#### 场景 3：搜索增强推荐
用户搜索节点时，优先显示推荐节点

#### 场景 4：从零开始的智能工作流构建 🆕
**核心价值**：根据用户需求描述，自动生成完整的工作流结构

用户只需输入需求描述，系统将：
1. **理解需求意图**：使用 LLM 分析用户的需求和目标
2. **匹配节点类型**：自动选择最合适的节点组合
3. **建立工作流**：生成完整的节点连接和配置
4. **提供解释**：说明每个节点的作用和整体流程逻辑

**示例交互：**

```
用户输入：
"我想创建一个课程大纲生成器，用户输入课程主题，系统生成包含章节的详细大纲"

系统输出：
┌─────────────────────────────────────────────────────────────┐
│ 🎯 工作流方案：课程大纲智能生成                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  [Start] → [LLM: 生成大纲] → [Code: 解析结构化数据]         │
│              ↓                      ↓                        │
│         [If: 检查质量] → [Loop: 优化各章节] → [End]         │
│                                                              │
└─────────────────────────────────────────────────────────────┘

推荐节点配置：
1. ✅ Start 节点
   - 接收参数：course_topic (string)
   
2. ✅ LLM 节点 - "生成课程大纲"
   - 模型：GPT-4
   - Prompt: "根据主题 {{course_topic}} 生成课程大纲..."
   - 输出格式：JSON
   
3. ✅ Code 节点 - "解析结构化数据"
   - 语言：Python
   - 功能：解析 LLM 输出，提取章节信息
   
4. ✅ If 节点 - "检查章节数量"
   - 条件：章节数 >= 3
   - 分支：满足→继续，不满足→重新生成
   
5. ✅ Loop 节点 - "遍历各章节"
   - 循环变量：chapter
   - 操作：为每章节生成详细内容
   
6. ✅ End 节点
   - 输出：完整的课程大纲 JSON
```

**请求数据结构：**

```go
type WorkflowGenerationRequest struct {
    UserRequirement  string            `json:"user_requirement"`   // 用户需求描述
    RequirementType  string            `json:"requirement_type"`   // 需求类型：text_generation, data_processing, api_integration, etc.
    Constraints      *Constraints      `json:"constraints"`        // 约束条件
    UserID           string            `json:"user_id"`
    Language         string            `json:"language"`           // 用户语言偏好
}

type Constraints struct {
    MaxNodes         int               `json:"max_nodes"`          // 最大节点数量
    AllowedNodeTypes []entity.NodeType `json:"allowed_node_types"` // 允许的节点类型
    Budget           *BudgetLimit      `json:"budget"`             // 成本预算（API 调用等）
    Performance      *PerformanceReq   `json:"performance"`        // 性能要求
}

type WorkflowGenerationResponse struct {
    WorkflowID       string                  `json:"workflow_id"`
    Nodes            []vo.Node               `json:"nodes"`          // 生成的节点列表
    Edges            []vo.Edge               `json:"edges"`          // 节点连接关系
    WorkflowGraph    *WorkflowGraph          `json:"workflow_graph"` // 可视化图形
    Explanation      string                  `json:"explanation"`    // 整体流程说明
    NodeExplanations map[string]string       `json:"node_explanations"` // 每个节点的说明
    EstimatedCost    *CostEstimate           `json:"estimated_cost"` // 预估成本
    Confidence       float64                 `json:"confidence"`     // 方案置信度 0-1
}
```

**LLM Prompt 设计：**

```
你是一个工作流设计专家，擅长将用户需求转化为可执行的工作流。

## 用户需求
{{user_requirement}}

## 可用节点类型
1. Start - 工作流入口，定义输入参数
2. LLM - 调用大语言模型生成内容
3. Code - 执行 Python 代码处理数据
4. Plugin - 调用外部 API 或服务
5. Knowledge - 知识库检索
6. Database - 数据库读写操作
7. If - 条件判断分支
8. Loop - 循环处理数组数据
9. End - 工作流出口

## 设计原则
1. 简洁高效：使用最少的节点完成任务
2. 鲁棒性：考虑异常处理和边界情况
3. 可维护性：节点职责单一，连接清晰
4. 最佳实践：遵循常见的工作流模式

## 输出要求
请设计一个完整的工作流方案，包括：
1. 节点列表（顺序排列）
2. 节点配置建议
3. 节点连接关系
4. 整体流程说明

## 输出格式（JSON）
{
  "workflow_name": "课程大纲生成器",
  "description": "根据用户输入的课程主题，自动生成结构化的课程大纲",
  "nodes": [
    {
      "id": "node_1",
      "type": "Start",
      "name": "开始",
      "position": {"x": 100, "y": 100},
      "config": {
        "inputs": [
          {"name": "course_topic", "type": "string", "required": true, "description": "课程主题"}
        ]
      },
      "explanation": "接收用户输入的课程主题作为工作流的输入参数"
    },
    {
      "id": "node_2",
      "type": "LLM",
      "name": "生成课程大纲",
      "position": {"x": 100, "y": 200},
      "config": {
        "model": "gpt-4",
        "prompt": "请根据以下课程主题生成详细的课程大纲...",
        "output_format": "json",
        "temperature": 0.7
      },
      "explanation": "使用 LLM 根据主题生成结构化的课程大纲"
    },
    ...
  ],
  "edges": [
    {"from": "node_1", "to": "node_2", "output_key": "course_topic", "input_key": "topic"},
    {"from": "node_2", "to": "node_3", "output_key": "outline", "input_key": "data"},
    ...
  ],
  "overall_explanation": "这个工作流首先接收课程主题，然后使用 LLM 生成初始大纲...",
  "confidence": 0.92
}
```

### 2.2 推荐算法设计

#### 2.2.1 多策略混合推荐

```python
# Python 伪代码
class NodeRecommendationEngine:
    def __init__(self):
        self.llm_strategy = LLMRecommendationStrategy()
        self.rule_strategy = RuleBasedStrategy()
        self.stats_strategy = StatisticsBasedStrategy()
        
    def recommend(self, request: NodeRecommendationRequest) -> List[RecommendedNode]:
        # 1. 收集各策略的推荐结果
        llm_results = self.llm_strategy.recommend(request, weight=0.4)
        rule_results = self.rule_strategy.recommend(request, weight=0.4)
        stats_results = self.stats_strategy.recommend(request, weight=0.2)
        
        # 2. 合并和排序
        all_results = self.merge_results([llm_results, rule_results, stats_results])
        
        # 3. 去重和过滤
        filtered_results = self.filter_incompatible(all_results, request)
        
        # 4. 返回 Top-N
        return filtered_results[:10]
```

#### 2.2.2 LLM 推理策略

**Prompt 设计示例：**

```
你是一个工作流设计专家。当前用户正在构建一个工作流，请根据以下信息推荐最合适的下一个节点。

## 当前工作流上下文
- 工作流目标：{{workflow_objective}}
- 已有节点：
  1. 开始节点 (Start)
  2. LLM 节点 - 生成课程大纲
  3. Code 节点 - 解析 JSON 数据 ← [当前节点]

## 节点类型说明
可选的节点类型包括：
1. LLM - 调用大语言模型
2. Code - 执行 Python 代码
3. Plugin - 调用外部 API
4. Database - 数据库操作
5. Knowledge - 知识库检索
6. If - 条件判断
7. Loop - 循环处理
8. End - 结束节点

## 当前节点输出
- output: Array<CourseChapter> (课程章节数组)
- chapter_count: Number (章节数量)

## 任务
请推荐 3-5 个最合适的后续节点，并说明推荐理由。

## 输出格式（JSON）
{
  "recommendations": [
    {
      "node_type": "Loop",
      "reason": "需要遍历每个课程章节，为每章生成详细内容",
      "score": 0.95,
      "suggested_config": {
        "loop_variable": "chapter",
        "items_path": "output"
      }
    },
    ...
  ]
}
```

#### 2.2.3 规则引擎策略

**规则定义示例（YAML）：**

```yaml
# backend/domain/workflow/internal/nodes/recommendation/rules.yaml
recommendation_rules:
  # 规则 1：LLM 节点后常接 Code 节点
  - name: llm_to_code
    condition:
      source_node_type: LLM
      output_format: JSON
    recommendations:
      - node_type: Code
        score: 0.8
        reason: "解析 LLM 返回的 JSON 数据"
      - node_type: If
        score: 0.6
        reason: "根据 LLM 结果进行条件判断"
        
  # 规则 2：数组输出后建议用 Loop
  - name: array_to_loop
    condition:
      source_output_type: Array
    recommendations:
      - node_type: Loop
        score: 0.9
        reason: "遍历数组中的每个元素"
      - node_type: Code
        score: 0.7
        reason: "批量处理数组数据"
        
  # 规则 3：数据处理后建议保存
  - name: data_processing_to_storage
    condition:
      source_node_category: 
        - Code
        - LLM
      workflow_has_database: true
    recommendations:
      - node_type: Database
        score: 0.75
        reason: "将处理结果保存到数据库"
      - node_type: Dataset
        score: 0.7
        reason: "将数据写入知识库"
        
  # 规则 4：检索后建议用 LLM 生成
  - name: retrieval_to_generation
    condition:
      source_node_type: 
        - Knowledge
        - Database
    recommendations:
      - node_type: LLM
        score: 0.85
        reason: "基于检索结果生成内容"
        
  # 规则 5：条件判断后建议用 End
  - name: if_to_end
    condition:
      source_node_type: If
      is_final_branch: true
    recommendations:
      - node_type: End
        score: 0.7
        reason: "结束工作流并返回结果"
```

#### 2.2.4 统计分析策略

**基于历史数据的推荐：**

```go
// 数据结构：节点转移概率矩阵
type NodeTransitionMatrix struct {
    transitions map[string]map[string]float64
    // transitions["LLM"]["Code"] = 0.65 表示 LLM 后接 Code 的概率是 65%
}

// 计算推荐得分
func (m *NodeTransitionMatrix) GetTransitionScore(from, to entity.NodeType) float64 {
    if scores, ok := m.transitions[string(from)]; ok {
        if score, ok := scores[string(to)]; ok {
            return score
        }
    }
    return 0.0
}

// 从历史工作流数据中学习
func (m *NodeTransitionMatrix) LearnFromWorkflows(workflows []*WorkflowHistory) {
    // 统计节点转移次数
    counts := make(map[string]map[string]int)
    
    for _, wf := range workflows {
        for i := 0; i < len(wf.Nodes)-1; i++ {
            from := string(wf.Nodes[i].Type)
            to := string(wf.Nodes[i+1].Type)
            
            if counts[from] == nil {
                counts[from] = make(map[string]int)
            }
            counts[from][to]++
        }
    }
    
    // 转换为概率
    for from, toCounts := range counts {
        total := 0
        for _, count := range toCounts {
            total += count
        }
        
        m.transitions[from] = make(map[string]float64)
        for to, count := range toCounts {
            m.transitions[from][to] = float64(count) / float64(total)
        }
    }
}
```

---

## 3. 技术架构

### 3.1 后端架构（Go）

```
backend/domain/workflow/internal/
├── recommendation/                    # 推荐系统核心
│   ├── engine.go                     # 推荐引擎主入口
│   ├── strategy.go                   # 策略接口定义
│   ├── llm_strategy.go               # LLM 推理策略
│   ├── rule_strategy.go              # 规则引擎策略
│   ├── stats_strategy.go             # 统计分析策略
│   ├── scorer.go                     # 评分和排序
│   ├── filter.go                     # 结果过滤
│   ├── rules.yaml                    # 规则配置
│   └── prompts.yaml                  # LLM Prompt 模板
├── nodes/
│   └── recommendation/               # 节点推荐相关节点
│       ├── recommender_node.go       # 推荐节点（可选）
│       └── config.go
└── analytics/                        # 数据分析
    ├── workflow_analyzer.go          # 工作流分析
    ├── pattern_miner.go              # 模式挖掘
    └── transition_matrix.go          # 转移概率矩阵
```

### 3.2 前端架构（TypeScript）

```
frontend/packages/workflow/playground/src/
├── services/
│   └── node-recommendation-service.ts   # 推荐服务
├── components/
│   ├── recommendation-panel/           # 推荐面板组件
│   │   ├── index.tsx
│   │   ├── recommendation-card.tsx
│   │   └── styles.module.less
│   └── node-panel/
│       └── enhanced-node-panel.tsx     # 增强的节点面板
└── hooks/
    └── use-node-recommendation.ts      # 推荐 Hook
```

### 3.3 API 接口设计

#### 3.3.1 推荐接口

```go
// POST /api/workflow/recommend/nodes
type RecommendNodesRequest struct {
    WorkflowID     string                 `json:"workflow_id"`
    SourceNodeID   string                 `json:"source_node_id"`
    Context        *WorkflowContext       `json:"context"`
    Limit          int                    `json:"limit"`           // 返回数量，默认 10
    IncludeReason  bool                   `json:"include_reason"`  // 是否包含推荐理由
}

type RecommendNodesResponse struct {
    Recommendations []RecommendedNode `json:"recommendations"`
    Strategies      []StrategyResult  `json:"strategies"`  // 各策略的得分（调试用）
}

type StrategyResult struct {
    Name   string  `json:"name"`
    Weight float64 `json:"weight"`
    Nodes  []RecommendedNode `json:"nodes"`
}
```

#### 3.3.2 反馈接口

```go
// POST /api/workflow/recommend/feedback
type RecommendationFeedbackRequest struct {
    WorkflowID       string `json:"workflow_id"`
    RecommendationID string `json:"recommendation_id"`
    SourceNodeID     string `json:"source_node_id"`
    SelectedNodeType string `json:"selected_node_type"`
    IsUseful         bool   `json:"is_useful"`  // 用户是否采纳推荐
}
```

---

## 4. 实施方案

### 4.1 实施阶段

#### 阶段 1：基础推荐（MVP）⭐ 优先

**目标**：实现基于规则的基础推荐功能

**工作项**：
1. ✅ 定义节点推荐规则（YAML 配置）
2. ✅ 实现规则引擎
3. ✅ 实现基础推荐 API
4. ✅ 前端推荐面板 UI

**预计时间**：1-2 周

**技术栈**：
- Go 规则引擎
- TypeScript/React 前端组件
- RESTful API

#### 阶段 2：LLM 增强推荐

**目标**：集成 LLM 进行智能推理

**工作项**：
1. ✅ 设计 LLM Prompt 模板
2. ✅ 实现 LLM 调用服务
3. ✅ 结果解析和验证
4. ✅ 多策略融合算法

**预计时间**：2-3 周

#### 阶段 3：数据驱动优化

**目标**：基于历史数据优化推荐

**工作项**：
1. ✅ 收集工作流使用数据
2. ✅ 构建节点转移矩阵
3. ✅ 实现协同过滤算法
4. ✅ A/B 测试框架

**预计时间**：3-4 周

#### 阶段 4：从零开始智能构建工作流 🆕

**目标**：实现基于需求描述自动生成完整工作流

**工作项**：
1. ✅ 设计工作流生成 Prompt 模板
2. ✅ 实现需求意图分析服务
3. ✅ 构建节点模板库和最佳实践库
4. ✅ 实现工作流自动布局算法
5. ✅ 节点配置自动推断
6. ✅ 工作流验证和优化
7. ✅ 前端可视化预览组件

**预计时间**：4-5 周

**核心功能**：

```go
// 工作流自动生成服务
type WorkflowGeneratorService struct {
    llmService        *llm.LLMService
    templateRepo      *TemplateRepository
    layoutEngine      *AutoLayoutEngine
    validationEngine  *WorkflowValidator
}

// 生成工作流的主要方法
func (s *WorkflowGeneratorService) GenerateWorkflow(
    ctx context.Context, 
    req *WorkflowGenerationRequest,
) (*WorkflowGenerationResponse, error) {
    
    // 1. 意图分析：理解用户需求
    intent, err := s.analyzeIntent(ctx, req.UserRequirement)
    if err != nil {
        return nil, err
    }
    
    // 2. 匹配模板：从最佳实践库中找相似工作流
    templates := s.templateRepo.FindSimilarTemplates(intent)
    
    // 3. LLM 生成：基于需求和模板生成工作流结构
    workflow, err := s.llmGenerateWorkflow(ctx, req, intent, templates)
    if err != nil {
        return nil, err
    }
    
    // 4. 节点配置：自动推断每个节点的配置
    configuredWorkflow := s.configureNodes(workflow, intent)
    
    // 5. 自动布局：计算节点位置，使工作流美观
    layoutWorkflow := s.layoutEngine.AutoLayout(configuredWorkflow)
    
    // 6. 验证优化：检查逻辑错误，优化性能
    validatedWorkflow, err := s.validationEngine.ValidateAndOptimize(layoutWorkflow)
    if err != nil {
        return nil, err
    }
    
    // 7. 生成说明：为每个节点和整体流程生成解释
    explanations := s.generateExplanations(validatedWorkflow, intent)
    
    return &WorkflowGenerationResponse{
        Nodes:            validatedWorkflow.Nodes,
        Edges:            validatedWorkflow.Edges,
        Explanation:      explanations.Overall,
        NodeExplanations: explanations.PerNode,
        Confidence:       validatedWorkflow.ConfidenceScore,
    }, nil
}
```

**最佳实践模板库**：

```yaml
# backend/domain/workflow/internal/templates/best_practices.yaml
workflow_templates:
  - id: content_generation
    name: "内容生成工作流"
    description: "使用 LLM 生成内容，包含质量检查和格式化"
    tags: ["llm", "content", "generation"]
    use_cases:
      - "文章生成"
      - "课程大纲生成"
      - "商品描述生成"
    pattern:
      nodes:
        - type: Start
          outputs: ["user_input"]
        - type: LLM
          name: "内容生成"
          inputs: ["user_input"]
          config:
            temperature: 0.7
            output_format: "text"
        - type: Code
          name: "内容检查"
          function: "validate_content"
        - type: If
          condition: "content_quality > 0.8"
          branches:
            - true: next_node
            - false: regenerate
        - type: End
      
  - id: data_extraction
    name: "数据提取与处理"
    description: "从非结构化数据中提取信息并结构化"
    tags: ["extraction", "nlp", "data_processing"]
    use_cases:
      - "简历信息提取"
      - "文档解析"
      - "表单数据提取"
    pattern:
      nodes:
        - type: Start
        - type: LLM
          name: "信息提取"
          config:
            output_format: "json"
        - type: Code
          name: "数据清洗"
        - type: Database
          operation: "insert"
        - type: End
      
  - id: rag_workflow
    name: "RAG 知识问答"
    description: "检索增强生成，基于知识库回答问题"
    tags: ["rag", "knowledge", "qa"]
    use_cases:
      - "智能客服"
      - "文档问答"
      - "知识库检索"
    pattern:
      nodes:
        - type: Start
        - type: Knowledge
          name: "知识检索"
          operation: "search"
        - type: Code
          name: "结果排序"
        - type: LLM
          name: "答案生成"
          config:
            with_context: true
        - type: End
        
  - id: batch_processing
    name: "批量数据处理"
    description: "循环处理大量数据项"
    tags: ["loop", "batch", "parallel"]
    use_cases:
      - "批量翻译"
      - "批量生成"
      - "数据迁移"
    pattern:
      nodes:
        - type: Start
        - type: Database
          operation: "query"
        - type: Loop
          name: "遍历数据"
          loop_nodes:
            - type: LLM
              name: "处理单项"
            - type: Database
              operation: "update"
        - type: End
```

**意图分析 Prompt**：

```
你是一个需求分析专家，请分析用户的工作流需求并提取关键信息。

## 用户需求
{{user_requirement}}

## 分析任务
请识别以下信息：
1. **核心目标**：用户想要实现什么功能？
2. **输入数据**：工作流需要什么输入？（类型、来源）
3. **输出结果**：期望得到什么输出？（格式、内容）
4. **关键步骤**：需要经过哪些主要处理步骤？
5. **约束条件**：是否有性能、成本、质量要求？
6. **工作流类型**：属于哪种模式？
   - content_generation (内容生成)
   - data_extraction (数据提取)
   - rag_workflow (知识问答)
   - batch_processing (批量处理)
   - api_integration (API 集成)
   - data_transformation (数据转换)
   - decision_automation (自动决策)

## 输出格式（JSON）
{
  "intent": {
    "core_objective": "生成结构化的课程大纲",
    "workflow_type": "content_generation",
    "complexity": "medium",  // simple, medium, complex
    "confidence": 0.95
  },
  "inputs": [
    {"name": "course_topic", "type": "string", "description": "课程主题"}
  ],
  "outputs": [
    {"name": "course_outline", "type": "json", "description": "包含章节的课程大纲"}
  ],
  "key_steps": [
    "接收课程主题",
    "使用 LLM 生成初始大纲",
    "解析和验证大纲结构",
    "优化章节内容",
    "输出最终结果"
  ],
  "constraints": {
    "quality": "high",
    "cost": "medium",
    "performance": "standard"
  },
  "suggested_nodes": [
    {"type": "LLM", "reason": "需要生成创意内容"},
    {"type": "Code", "reason": "需要解析和验证数据"},
    {"type": "If", "reason": "需要质量检查"}
  ]
}
```

**前端交互流程**：

```typescript
// frontend/apps/web/src/features/workflow/components/WorkflowGenerator.tsx

interface WorkflowGeneratorProps {
  onWorkflowGenerated: (workflow: GeneratedWorkflow) => void;
}

export const WorkflowGenerator: React.FC<WorkflowGeneratorProps> = ({
  onWorkflowGenerated
}) => {
  const [requirement, setRequirement] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedWorkflow, setGeneratedWorkflow] = useState<GeneratedWorkflow | null>(null);
  
  const handleGenerate = async () => {
    setIsGenerating(true);
    
    try {
      // 1. 调用工作流生成 API
      const response = await workflowAPI.generateWorkflow({
        user_requirement: requirement,
        language: 'zh-CN'
      });
      
      // 2. 显示预览
      setGeneratedWorkflow(response);
      
      // 3. 用户可以编辑调整
      // ...
      
    } catch (error) {
      message.error('生成工作流失败');
    } finally {
      setIsGenerating(false);
    }
  };
  
  return (
    <div className="workflow-generator">
      <Card title="🤖 智能工作流生成器">
        <TextArea
          rows={4}
          placeholder="请描述您想要创建的工作流，例如：&#10;我想创建一个课程大纲生成器，用户输入课程主题，系统生成包含章节的详细大纲"
          value={requirement}
          onChange={(e) => setRequirement(e.target.value)}
        />
        
        <Button
          type="primary"
          loading={isGenerating}
          onClick={handleGenerate}
          disabled={!requirement}
        >
          🎯 生成工作流
        </Button>
      </Card>
      
      {generatedWorkflow && (
        <WorkflowPreview
          workflow={generatedWorkflow}
          onAccept={() => onWorkflowGenerated(generatedWorkflow)}
          onModify={(modified) => setGeneratedWorkflow(modified)}
        />
      )}
    </div>
  );
};
```

### 4.2 技术实施路径

#### 4.2.1 方案 A：使用现有 LLM 节点（推荐）✅

**优点**：
- ✅ 复用现有基础设施
- ✅ 无需新增节点类型
- ✅ 可以直接使用已配置的 LLM 模型
- ✅ 开发速度快

**实现方式**：

```go
// 1. 在 backend/domain/workflow/internal/recommendation/ 新增推荐引擎

// engine.go
package recommendation

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/llm"
    "github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
)

type RecommendationEngine struct {
    llmConfig     *llm.Config           // 复用 LLM 节点配置
    ruleEngine    *RuleEngine
    statsEngine   *StatisticsEngine
}

func NewRecommendationEngine(llmCfg *llm.Config) *RecommendationEngine {
    return &RecommendationEngine{
        llmConfig:   llmCfg,
        ruleEngine:  NewRuleEngine(),
        statsEngine: NewStatisticsEngine(),
    }
}

func (e *RecommendationEngine) Recommend(ctx context.Context, req *RecommendRequest) (*RecommendResponse, error) {
    // 1. 规则推荐
    ruleResults := e.ruleEngine.Recommend(req)
    
    // 2. LLM 推荐（使用现有 LLM 节点）
    llmResults, err := e.llmRecommend(ctx, req)
    if err != nil {
        // LLM 失败时降级到规则推荐
        return &RecommendResponse{Recommendations: ruleResults}, nil
    }
    
    // 3. 统计推荐
    statsResults := e.statsEngine.Recommend(req)
    
    // 4. 融合结果
    merged := e.mergeAndRank(ruleResults, llmResults, statsResults)
    
    return &RecommendResponse{Recommendations: merged}, nil
}

func (e *RecommendationEngine) llmRecommend(ctx context.Context, req *RecommendRequest) ([]RecommendedNode, error) {
    // 构造 Prompt
    prompt := e.buildPrompt(req)
    
    // 调用 LLM 节点
    llmInput := map[string]any{
        "system_prompt": recommendationSystemPrompt,
        "user_prompt":   prompt,
        "output_format": "json",
    }
    
    // 使用现有的 LLM 节点执行
    llmNode, err := e.llmConfig.Build(ctx, &schema.NodeSchema{
        Type: entity.NodeTypeLLM,
        // ... 配置
    })
    if err != nil {
        return nil, err
    }
    
    output, err := llmNode.Invoke(ctx, llmInput)
    if err != nil {
        return nil, err
    }
    
    // 解析 LLM 输出
    return e.parseL LMOutput(output)
}
```

#### 4.2.2 方案 B：使用代码节点（备选）

**适用场景**：
- 需要复杂的规则逻辑
- 需要访问外部 API
- 需要自定义算法

**实现方式**：

```go
// 使用 Code 节点执行推荐脚本
func (e *RecommendationEngine) codeRecommend(ctx context.Context, req *RecommendRequest) ([]RecommendedNode, error) {
    // Python 推荐脚本
    pythonScript := `
import json
from typing import List, Dict

def recommend_nodes(workflow_context: Dict, source_node: Dict) -> List[Dict]:
    """
    基于上下文推荐节点
    """
    recommendations = []
    
    # 规则 1: 如果源节点是 LLM 且输出 JSON，推荐 Code 节点
    if source_node['type'] == 'LLM' and source_node.get('output_format') == 'json':
        recommendations.append({
            'node_type': 'Code',
            'score': 0.9,
            'reason': '解析 LLM 返回的 JSON 数据'
        })
    
    # 规则 2: 如果源节点输出数组，推荐 Loop 节点
    if 'Array' in str(source_node.get('output_types', [])):
        recommendations.append({
            'node_type': 'Loop',
            'score': 0.85,
            'reason': '遍历数组元素'
        })
    
    # ... 更多规则
    
    return sorted(recommendations, key=lambda x: x['score'], reverse=True)

# 主函数
workflow_ctx = json.loads(input['workflow_context'])
source_node = json.loads(input['source_node'])

result = recommend_nodes(workflow_ctx, source_node)
output = {'recommendations': result}
`
    
    // 构造 Code 节点输入
    codeInput := map[string]any{
        "code": pythonScript,
        "input": map[string]any{
            "workflow_context": req.WorkflowContext,
            "source_node": req.SourceNode,
        },
    }
    
    // 执行 Code 节点
    codeNode := code.NewCodeRunnerNode(/* config */)
    output, err := codeNode.Invoke(ctx, codeInput)
    if err != nil {
        return nil, err
    }
    
    return parseCodeOutput(output)
}
```

---

## 5. 代码示例

### 5.1 规则引擎实现

```go
// backend/domain/workflow/internal/recommendation/rule_engine.go
package recommendation

import (
    "context"
    "gopkg.in/yaml.v3"
    "io/ioutil"
)

type RuleEngine struct {
    rules []RecommendationRule
}

type RecommendationRule struct {
    Name          string                 `yaml:"name"`
    Condition     RuleCondition          `yaml:"condition"`
    Recommendations []RuleRecommendation `yaml:"recommendations"`
}

type RuleCondition struct {
    SourceNodeType    *entity.NodeType   `yaml:"source_node_type"`
    SourceNodeTypes   []entity.NodeType  `yaml:"source_node_types"`
    OutputType        *string            `yaml:"output_type"`
    OutputFormat      *string            `yaml:"output_format"`
    WorkflowHasDB     *bool              `yaml:"workflow_has_database"`
    NodeCategory      []string           `yaml:"source_node_category"`
}

type RuleRecommendation struct {
    NodeType entity.NodeType `yaml:"node_type"`
    Score    float64         `yaml:"score"`
    Reason   string          `yaml:"reason"`
    Config   map[string]any  `yaml:"suggested_config"`
}

func NewRuleEngine() *RuleEngine {
    engine := &RuleEngine{}
    engine.loadRules("rules.yaml")
    return engine
}

func (e *RuleEngine) loadRules(filepath string) error {
    data, err := ioutil.ReadFile(filepath)
    if err != nil {
        return err
    }
    
    var config struct {
        Rules []RecommendationRule `yaml:"recommendation_rules"`
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return err
    }
    
    e.rules = config.Rules
    return nil
}

func (e *RuleEngine) Recommend(req *RecommendRequest) []RecommendedNode {
    var results []RecommendedNode
    
    for _, rule := range e.rules {
        if e.matchCondition(rule.Condition, req) {
            for _, rec := range rule.Recommendations {
                results = append(results, RecommendedNode{
                    NodeType:     rec.NodeType,
                    Score:        rec.Score,
                    Reason:       rec.Reason,
                    DisplayName:  getNodeDisplayName(rec.NodeType),
                    Category:     getNodeCategory(rec.NodeType),
                    Template:     buildNodeTemplate(rec.NodeType, rec.Config),
                })
            }
        }
    }
    
    return results
}

func (e *RuleEngine) matchCondition(cond RuleCondition, req *RecommendRequest) bool {
    // 检查源节点类型
    if cond.SourceNodeType != nil && *cond.SourceNodeType != req.SourceNode.Type {
        return false
    }
    
    if len(cond.SourceNodeTypes) > 0 {
        matched := false
        for _, t := range cond.SourceNodeTypes {
            if t == req.SourceNode.Type {
                matched = true
                break
            }
        }
        if !matched {
            return false
        }
    }
    
    // 检查输出类型
    if cond.OutputType != nil {
        hasType := false
        for _, output := range req.SourceNode.Outputs {
            if output.Type == *cond.OutputType {
                hasType = true
                break
            }
        }
        if !hasType {
            return false
        }
    }
    
    // 检查输出格式
    if cond.OutputFormat != nil {
        if req.SourceNode.OutputFormat != *cond.OutputFormat {
            return false
        }
    }
    
    // 检查工作流是否有数据库
    if cond.WorkflowHasDB != nil {
        hasDB := req.WorkflowContext.HasDatabaseNode()
        if *cond.WorkflowHasDB != hasDB {
            return false
        }
    }
    
    return true
}
```

### 5.2 LLM 策略实现

```go
// backend/domain/workflow/internal/recommendation/llm_strategy.go
package recommendation

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/llm"
)

type LLMStrategy struct {
    llmNode *llm.LLM
    prompts *PromptTemplates
}

type PromptTemplates struct {
    SystemPrompt string
    UserTemplate string
}

func NewLLMStrategy(llmNode *llm.LLM) *LLMStrategy {
    return &LLMStrategy{
        llmNode: llmNode,
        prompts: loadPromptTemplates(),
    }
}

func (s *LLMStrategy) Recommend(ctx context.Context, req *RecommendRequest) ([]RecommendedNode, error) {
    // 1. 构建 Prompt
    userPrompt := s.buildPrompt(req)
    
    // 2. 调用 LLM
    llmInput := map[string]any{
        "system_prompt": s.prompts.SystemPrompt,
        "user_prompt":   userPrompt,
        "output_schema": recommendationOutputSchema,
    }
    
    output, err := s.llmNode.Invoke(ctx, llmInput)
    if err != nil {
        return nil, fmt.Errorf("LLM invocation failed: %w", err)
    }
    
    // 3. 解析输出
    return s.parseOutput(output)
}

func (s *LLMStrategy) buildPrompt(req *RecommendRequest) string {
    // 序列化上下文信息
    contextJSON, _ := json.MarshalIndent(req.WorkflowContext, "", "  ")
    sourceNodeJSON, _ := json.MarshalIndent(req.SourceNode, "", "  ")
    
    // 填充模板
    return fmt.Sprintf(s.prompts.UserTemplate, 
        string(contextJSON),
        string(sourceNodeJSON),
        formatAvailableNodes(),
    )
}

func (s *LLMStrategy) parseOutput(output map[string]any) ([]RecommendedNode, error) {
    // 从 LLM 输出中提取推荐结果
    recommendationsRaw, ok := output["recommendations"]
    if !ok {
        return nil, fmt.Errorf("no recommendations in LLM output")
    }
    
    // 反序列化
    var llmRecs []struct {
        NodeType        string         `json:"node_type"`
        Reason          string         `json:"reason"`
        Score           float64        `json:"score"`
        SuggestedConfig map[string]any `json:"suggested_config"`
    }
    
    data, _ := json.Marshal(recommendationsRaw)
    if err := json.Unmarshal(data, &llmRecs); err != nil {
        return nil, err
    }
    
    // 转换为标准格式
    var results []RecommendedNode
    for _, rec := range llmRecs {
        nodeType := entity.NodeType(rec.NodeType)
        results = append(results, RecommendedNode{
            NodeType:    nodeType,
            Score:       rec.Score,
            Reason:      rec.Reason,
            DisplayName: getNodeDisplayName(nodeType),
            Category:    getNodeCategory(nodeType),
            Template:    buildNodeTemplate(nodeType, rec.SuggestedConfig),
        })
    }
    
    return results, nil
}

// Prompt 模板
const recommendationSystemPrompt = `你是一个工作流设计专家助手。你的任务是分析用户当前的工作流状态，并推荐最合适的后续节点。

## 推荐原则
1. **目标导向**：推荐的节点应该帮助用户实现工作流的最终目标
2. **数据兼容**：推荐的节点应该能够接收上游节点的输出
3. **最佳实践**：优先推荐常见且成熟的工作流模式
4. **简洁高效**：避免推荐冗余的节点

## 可用节点类型
{{AVAILABLE_NODES}}

## 输出要求
- 返回 3-5 个最合适的推荐
- 每个推荐必须包含：node_type, reason, score (0-1)
- 按 score 降序排列
- 输出必须是有效的 JSON 格式`

const recommendationUserTemplate = `
## 工作流上下文
%s

## 当前节点信息
%s

## 可用节点类型
%s

请推荐最合适的 3-5 个后续节点。
`
```

### 5.3 前端推荐面板

```typescript
// frontend/packages/workflow/playground/src/components/recommendation-panel/index.tsx
import React, { useEffect, useState } from 'react';
import { Card, Spin, Empty } from '@arco-design/web-react';
import { IconRobot, IconLightning } from '@arco-design/web-react/icon';
import { useNodeRecommendation } from '../../hooks/use-node-recommendation';
import type { RecommendedNode } from '../../types/recommendation';
import styles from './styles.module.less';

export interface RecommendationPanelProps {
  sourceNodeId: string;
  workflowId: string;
  onSelect: (node: RecommendedNode) => void;
  visible: boolean;
}

export const RecommendationPanel: React.FC<RecommendationPanelProps> = ({
  sourceNodeId,
  workflowId,
  onSelect,
  visible,
}) => {
  const { recommendations, loading, error, fetchRecommendations } = 
    useNodeRecommendation();

  useEffect(() => {
    if (visible && sourceNodeId) {
      fetchRecommendations({
        workflowId,
        sourceNodeId,
        limit: 5,
        includeReason: true,
      });
    }
  }, [visible, sourceNodeId, workflowId]);

  if (!visible) return null;

  return (
    <div className={styles.recommendationPanel}>
      <div className={styles.header}>
        <IconRobot className={styles.icon} />
        <span>智能推荐</span>
      </div>

      <div className={styles.content}>
        {loading && (
          <div className={styles.loading}>
            <Spin />
            <div>正在分析工作流...</div>
          </div>
        )}

        {error && (
          <Empty 
            description={error.message || '推荐失败，请重试'}
            icon={<IconLightning />}
          />
        )}

        {!loading && !error && recommendations.length === 0 && (
          <Empty description="暂无推荐节点" />
        )}

        {!loading && !error && recommendations.length > 0 && (
          <div className={styles.recommendations}>
            {recommendations.map((rec, index) => (
              <RecommendationCard
                key={`${rec.nodeType}-${index}`}
                recommendation={rec}
                onClick={() => onSelect(rec)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

interface RecommendationCardProps {
  recommendation: RecommendedNode;
  onClick: () => void;
}

const RecommendationCard: React.FC<RecommendationCardProps> = ({
  recommendation,
  onClick,
}) => {
  return (
    <Card
      className={styles.recommendationCard}
      hoverable
      onClick={onClick}
    >
      <div className={styles.cardHeader}>
        <img src={recommendation.icon} className={styles.nodeIcon} />
        <div className={styles.nodeName}>{recommendation.displayName}</div>
        <div className={styles.score}>
          {Math.round(recommendation.score * 100)}%
        </div>
      </div>
      <div className={styles.reason}>{recommendation.reason}</div>
      <div className={styles.category}>{recommendation.category}</div>
    </Card>
  );
};
```

### 5.4 推荐 Hook

```typescript
// frontend/packages/workflow/playground/src/hooks/use-node-recommendation.ts
import { useState, useCallback } from 'react';
import { recommendNodes } from '../services/node-recommendation-service';
import type { 
  RecommendNodesRequest, 
  RecommendNodesResponse,
  RecommendedNode 
} from '../types/recommendation';

export interface UseNodeRecommendationReturn {
  recommendations: RecommendedNode[];
  loading: boolean;
  error: Error | null;
  fetchRecommendations: (req: RecommendNodesRequest) => Promise<void>;
  clearRecommendations: () => void;
}

export const useNodeRecommendation = (): UseNodeRecommendationReturn => {
  const [recommendations, setRecommendations] = useState<RecommendedNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchRecommendations = useCallback(async (req: RecommendNodesRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await recommendNodes(req);
      setRecommendations(response.recommendations);
    } catch (err) {
      setError(err as Error);
      setRecommendations([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const clearRecommendations = useCallback(() => {
    setRecommendations([]);
    setError(null);
  }, []);

  return {
    recommendations,
    loading,
    error,
    fetchRecommendations,
    clearRecommendations,
  };
};
```

### 5.5 推荐服务

```typescript
// frontend/packages/workflow/playground/src/services/node-recommendation-service.ts
import type { 
  RecommendNodesRequest, 
  RecommendNodesResponse 
} from '../types/recommendation';

const API_BASE = '/api/workflow';

export async function recommendNodes(
  req: RecommendNodesRequest
): Promise<RecommendNodesResponse> {
  const response = await fetch(`${API_BASE}/recommend/nodes`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(req),
  });

  if (!response.ok) {
    throw new Error(`推荐请求失败: ${response.statusText}`);
  }

  return response.json();
}

export async function submitRecommendationFeedback(
  workflowId: string,
  recommendationId: string,
  selectedNodeType: string,
  isUseful: boolean
): Promise<void> {
  await fetch(`${API_BASE}/recommend/feedback`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      workflowId,
      recommendationId,
      selectedNodeType,
      isUseful,
    }),
  });
}
```

### 5.6 工作流自动生成器实现 🆕

#### 5.6.1 后端生成服务

```go
// backend/domain/workflow/internal/generator/workflow_generator.go
package generator

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
    "github.com/coze-dev/coze-studio/backend/domain/workflow/vo"
)

// WorkflowGenerator 工作流自动生成器
type WorkflowGenerator struct {
    llmService       LLMService
    templateRepo     *TemplateRepository
    layoutEngine     *AutoLayoutEngine
    validationEngine *WorkflowValidator
    promptTemplate   *PromptTemplate
}

// GenerateWorkflow 根据需求生成完整工作流
func (g *WorkflowGenerator) GenerateWorkflow(
    ctx context.Context,
    req *WorkflowGenerationRequest,
) (*WorkflowGenerationResponse, error) {
    
    // 步骤 1：意图分析
    intent, err := g.analyzeIntent(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("意图分析失败: %w", err)
    }
    
    // 步骤 2：匹配最佳实践模板
    templates := g.templateRepo.FindSimilarTemplates(intent.WorkflowType, 3)
    
    // 步骤 3：LLM 生成工作流结构
    rawWorkflow, err := g.llmGenerateWorkflow(ctx, req, intent, templates)
    if err != nil {
        return nil, fmt.Errorf("LLM 生成失败: %w", err)
    }
    
    // 步骤 4：节点配置自动推断
    configuredWorkflow := g.configureNodes(rawWorkflow, intent)
    
    // 步骤 5：自动布局
    layoutWorkflow := g.layoutEngine.AutoLayout(configuredWorkflow)
    
    // 步骤 6：验证和优化
    validatedWorkflow, validationErrors := g.validationEngine.Validate(layoutWorkflow)
    if len(validationErrors) > 0 {
        // 尝试自动修复
        validatedWorkflow = g.autoFixErrors(validatedWorkflow, validationErrors)
    }
    
    // 步骤 7：生成解释说明
    explanations := g.generateExplanations(validatedWorkflow, intent)
    
    return &WorkflowGenerationResponse{
        WorkflowID:       generateWorkflowID(),
        Nodes:            validatedWorkflow.Nodes,
        Edges:            validatedWorkflow.Edges,
        Explanation:      explanations.Overall,
        NodeExplanations: explanations.PerNode,
        Confidence:       intent.Confidence,
    }, nil
}

// analyzeIntent 分析用户需求意图
func (g *WorkflowGenerator) analyzeIntent(
    ctx context.Context,
    req *WorkflowGenerationRequest,
) (*IntentAnalysisResult, error) {
    
    prompt := g.promptTemplate.BuildIntentAnalysisPrompt(req.UserRequirement)
    
    llmResp, err := g.llmService.Call(ctx, &LLMRequest{
        Prompt:       prompt,
        Model:        "gpt-4",
        Temperature:  0.3, // 低温度，保证输出稳定
        OutputFormat: "json",
    })
    if err != nil {
        return nil, err
    }
    
    var intent IntentAnalysisResult
    if err := json.Unmarshal([]byte(llmResp.Content), &intent); err != nil {
        return nil, fmt.Errorf("解析意图分析结果失败: %w", err)
    }
    
    return &intent, nil
}
```

#### 5.6.2 前端生成器组件

```typescript
// frontend/apps/web/src/features/workflow/components/WorkflowGenerator.tsx
import React, { useState } from 'react';
import { Card, Button, Input, message, Steps, Spin } from 'antd';
import { RobotOutlined, CheckCircleOutlined, LoadingOutlined } from '@ant-design/icons';
import type { GeneratedWorkflow } from '../types';
import { WorkflowPreview } from './WorkflowPreview';
import { workflowAPI } from '../services/workflow-api';

const { TextArea } = Input;
const { Step } = Steps;

export const WorkflowGenerator: React.FC = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [requirement, setRequirement] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedWorkflow, setGeneratedWorkflow] = useState<GeneratedWorkflow | null>(null);
  
  // 示例需求
  const exampleRequirements = [
    '创建一个课程大纲生成器，输入课程主题，生成包含章节的详细大纲',
    '构建一个智能客服问答系统，从知识库检索相关信息后生成回答',
    '实现批量文章翻译功能，读取数据库中的文章列表并逐个翻译',
  ];
  
  const handleGenerate = async () => {
    if (!requirement.trim()) {
      message.warning('请输入需求描述');
      return;
    }
    
    setIsGenerating(true);
    setCurrentStep(1);
    
    try {
      const response = await workflowAPI.generateWorkflow({
        user_requirement: requirement,
        language: 'zh-CN',
        constraints: { max_nodes: 10 },
      });
      
      setGeneratedWorkflow(response);
      setCurrentStep(2);
      message.success('工作流生成成功！');
      
    } catch (error) {
      message.error('生成工作流失败，请重试');
      setCurrentStep(0);
    } finally {
      setIsGenerating(false);
    }
  };
  
  return (
    <div className="workflow-generator">
      <Steps current={currentStep} style={{ marginBottom: 24 }}>
        <Step title="描述需求" icon={<RobotOutlined />} />
        <Step title="生成工作流" icon={isGenerating ? <LoadingOutlined /> : null} />
        <Step title="预览确认" icon={<CheckCircleOutlined />} />
      </Steps>
      
      {currentStep === 0 && (
        <Card title="🤖 智能工作流生成器">
          <TextArea
            rows={6}
            placeholder="请详细描述您想要创建的工作流功能和流程..."
            value={requirement}
            onChange={(e) => setRequirement(e.target.value)}
          />
          
          <div style={{ marginTop: 16 }}>
            <p>💡 示例需求：</p>
            {exampleRequirements.map((example, index) => (
              <div key={index} onClick={() => setRequirement(example)} 
                   style={{ cursor: 'pointer', padding: 8, background: '#f5f5f5', margin: '4px 0' }}>
                {example}
              </div>
            ))}
          </div>
          
          <Button type="primary" onClick={handleGenerate} loading={isGenerating} 
                  style={{ marginTop: 16 }} icon={<RobotOutlined />}>
            🎯 智能生成工作流
          </Button>
        </Card>
      )}
      
      {currentStep === 1 && (
        <Card>
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin size="large" />
            <p style={{ marginTop: 16 }}>🤖 正在分析需求并生成工作流...</p>
          </div>
        </Card>
      )}
      
      {currentStep === 2 && generatedWorkflow && (
        <WorkflowPreview workflow={generatedWorkflow} />
      )}
    </div>
  );
};
```

---

## 6. 部署指南

### 6.1 配置文件

**rules.yaml**

```yaml
# backend/domain/workflow/internal/recommendation/rules.yaml
recommendation_rules:
  # ... (见前面的规则定义示例)
```

**prompts.yaml**

```yaml
# backend/domain/workflow/internal/recommendation/prompts.yaml
system_prompt: |
  你是一个工作流设计专家助手...

user_prompt_template: |
  ## 工作流上下文
  {{workflow_context}}
  
  ## 当前节点信息
  {{source_node}}
  
  请推荐最合适的后续节点。
```

### 6.2 环境变量

```bash
# .env
RECOMMENDATION_ENABLED=true
RECOMMENDATION_LLM_MODEL=gpt-4
RECOMMENDATION_MAX_RESULTS=10
RECOMMENDATION_CACHE_TTL=300  # 缓存 5 分钟
```

### 6.3 数据库表（可选）

```sql
-- 用于存储推荐反馈和学习
CREATE TABLE workflow_recommendation_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    workflow_id VARCHAR(64) NOT NULL,
    source_node_id VARCHAR(64) NOT NULL,
    recommended_node_type VARCHAR(64) NOT NULL,
    selected_node_type VARCHAR(64),
    is_useful BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_workflow (workflow_id),
    INDEX idx_source (source_node_id)
);

-- 用于存储节点转移统计
CREATE TABLE workflow_node_transitions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    from_node_type VARCHAR(64) NOT NULL,
    to_node_type VARCHAR(64) NOT NULL,
    count INT DEFAULT 1,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_transition (from_node_type, to_node_type)
);
```

---

## 7. 最佳实践

### 7.1 性能优化

1. **缓存推荐结果**
   ```go
   // 使用 Redis 缓存推荐结果
   cacheKey := fmt.Sprintf("recommendation:%s:%s", workflowID, sourceNodeID)
   if cached, err := redisClient.Get(cacheKey); err == nil {
       return parseCache(cached)
   }
   ```

2. **异步计算**
   ```go
   // LLM 推荐异步执行，不阻塞规则推荐
   go func() {
       llmResults := e.llmRecommend(ctx, req)
       e.cache.Set(cacheKey, llmResults, 5*time.Minute)
   }()
   ```

3. **批量推荐**
   ```go
   // 支持批量查询多个源节点的推荐
   type BatchRecommendRequest struct {
       SourceNodeIDs []string
   }
   ```

### 7.2 监控和调试

1. **推荐质量监控**
   ```go
   // 记录推荐指标
   metrics.RecordRecommendation(metrics.RecommendationMetrics{
       SourceNodeType:    req.SourceNode.Type,
       RecommendedCount:  len(results),
       LLMLatency:        llmDuration,
       RuleLatency:       ruleDuration,
       UserAccepted:      feedback.IsUseful,
   })
   ```

2. **A/B 测试**
   ```go
   // 随机分配用户到不同的推荐策略
   if userID%2 == 0 {
       return e.strategyA.Recommend(req)
   } else {
       return e.strategyB.Recommend(req)
   }
   ```

### 7.3 用户体验优化

1. **渐进式加载**
   - 先显示规则推荐（快速）
   - 后续补充 LLM 推荐（较慢但更智能）

2. **推荐理由可视化**
   - 显示推荐得分
   - 解释推荐原因
   - 高亮匹配的规则

3. **用户反馈循环**
   - 收集用户选择数据
   - 定期重训练模型
   - 优化推荐权重

---

## 8. 总结

### 8.1 技术方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **规则引擎** | 快速、可控、可解释 | 需要人工维护规则 | ⭐⭐⭐⭐⭐ |
| **LLM 推理** | 智能、灵活、上下文感知 | 延迟高、成本高 | ⭐⭐⭐⭐ |
| **统计学习** | 数据驱动、自动优化 | 需要大量数据 | ⭐⭐⭐ |
| **混合方案** | 综合优势、效果最佳 | 复杂度高 | ⭐⭐⭐⭐⭐ |
| **🆕 工作流自动生成** | 零门槛、高效率、端到端 | LLM 成本、需验证 | ⭐⭐⭐⭐⭐ |

### 8.2 推荐实施路径

**MVP 阶段（2 周）**：
1. ✅ 实现基于规则的推荐引擎
2. ✅ 定义 20-30 条核心推荐规则
3. ✅ 开发前端推荐面板
4. ✅ 完成基础 API 接口

**增强阶段（3-4 周）**：
1. ✅ 集成 LLM 智能推理
2. ✅ 实现多策略融合算法
3. ✅ 添加用户反馈机制
4. ✅ 性能优化和缓存

**优化阶段（4-6 周）**：
1. ✅ 收集和分析用户数据
2. ✅ 构建节点转移矩阵
3. ✅ 实现个性化推荐
4. ✅ A/B 测试和持续优化

**🆕 工作流自动生成阶段（4-5 周）**：
1. ✅ 实现需求意图分析服务
2. ✅ 构建工作流模板库和最佳实践库
3. ✅ 开发自动布局和验证引擎
4. ✅ 实现前端生成器和预览组件
5. ✅ 集成到工作流编辑器
6. ✅ 用户体验优化和测试

### 8.3 核心功能对比

| 功能类型 | 使用场景 | 交互方式 | 输出结果 |
|---------|---------|---------|---------|
| **连线推荐** | 已有节点，添加下一个 | 拖拽连线触发 | 推荐节点列表 |
| **🆕 零开始构建** | 空白画布，从需求开始 | 自然语言描述 | 完整工作流 |

**连线推荐 vs 零开始构建**：

```
┌─────────────────────────────────────────────────────────────┐
│ 连线推荐：渐进式构建                                          │
├─────────────────────────────────────────────────────────────┤
│  用户：添加 Start 节点                                        │
│  系统：[推荐] LLM, Code, Plugin...                           │
│  用户：选择 LLM 节点并连接                                    │
│  系统：[推荐] Code, If, End...                               │
│  用户：继续选择...                                            │
│  （多次交互，逐步构建）                                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 🆕 零开始构建：一步到位                                       │
├─────────────────────────────────────────────────────────────┤
│  用户：输入需求描述                                           │
│       "创建课程大纲生成器，输入主题，生成章节"                 │
│  系统：[生成完整工作流]                                       │
│       Start → LLM → Code → If → Loop → End                   │
│       （包含所有节点配置和连接）                              │
│  用户：预览、确认或微调                                       │
│  （一次交互，快速完成）                                       │
└─────────────────────────────────────────────────────────────┘
```

### 8.4 成功指标

**连线推荐指标**：
- **推荐准确率**：用户采纳推荐的比例 > 60%
- **响应时间**：推荐结果返回 < 500ms（规则）/ < 2s（LLM）
- **覆盖率**：80% 的场景都能给出推荐
- **用户满意度**：推荐有用率 > 70%

**🆕 工作流自动生成指标**：
- **生成成功率**：能够成功生成可执行工作流 > 85%
- **生成时间**：从需求到工作流 < 10s
- **节点准确性**：生成的节点类型符合需求 > 80%
- **用户采纳率**：直接使用或微调后使用 > 70%
- **置信度阈值**：系统置信度 > 0.8 时推荐给用户

### 8.5 应用场景示例

#### 场景 1：新手用户快速上手 🆕

```
用户：我是新手，想创建一个工作流但不知道从哪开始
系统：请告诉我您想实现什么功能
用户：我想做一个自动回答客户问题的机器人
系统：[自动生成] 
      Start → Knowledge(搜索知识库) → LLM(生成回答) → End
用户：太好了！直接用这个
```

#### 场景 2：专家用户精细控制

```
用户：我要构建一个复杂的数据处理流程
系统：[提供] 工作流生成器 + 连线推荐
用户：先用生成器创建基础结构
系统：[生成] 基础工作流框架
用户：然后手动添加特殊节点，系统提供推荐
系统：[推荐] 基于上下文的节点建议
```

### 8.6 技术架构优势

**层次化设计**：

```
┌─────────────────────────────────────────────────────────────┐
│                      用户界面层                               │
│  ┌──────────────────┐         ┌──────────────────┐          │
│  │  连线推荐面板     │         │ 🆕 工作流生成器  │          │
│  └──────────────────┘         └──────────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层                               │
│  ┌──────────────────┐         ┌──────────────────┐          │
│  │  推荐引擎         │         │ 🆕 生成引擎      │          │
│  │  - 规则策略       │         │  - 意图分析      │          │
│  │  - LLM策略        │         │  - 模板匹配      │          │
│  │  - 统计策略       │         │  - 自动配置      │          │
│  └──────────────────┘         └──────────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      基础设施层                               │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐            │
│  │  LLM   │  │ 规则库 │  │ 模板库 │  │ 数据库 │            │
│  └────────┘  └────────┘  └────────┘  └────────┘            │
└─────────────────────────────────────────────────────────────┘
```

**复用与扩展**：
- ✅ LLM 服务复用于推荐和生成
- ✅ 规则引擎支持两种场景
- ✅ 模板库持续积累最佳实践
- ✅ 统一的验证和布局引擎

---

## 9. 🆕 工作流自动生成详细说明

### 9.1 核心价值

**问题**：
- ❌ 新手用户不知道如何开始
- ❌ 需要了解所有节点类型才能设计工作流
- ❌ 重复构建相似的工作流模式
- ❌ 节点配置复杂，容易出错

**解决方案**：
- ✅ 只需用自然语言描述需求
- ✅ 系统自动选择合适的节点
- ✅ 自动生成完整的工作流结构
- ✅ 智能配置每个节点参数

### 9.2 技术亮点

1. **意图识别**
   - 使用 LLM 理解用户需求
   - 提取关键信息（输入、输出、步骤）
   - 分类工作流类型（生成、提取、问答等）

2. **模板匹配**
   - 内置常见工作流最佳实践
   - 基于相似度匹配模板
   - 自适应调整模板结构

3. **智能配置**
   - 根据需求自动配置 Prompt
   - 推断合适的模型和参数
   - 设置数据流连接关系

4. **自动布局**
   - 拓扑排序计算节点层级
   - 层次化布局算法
   - 优化交叉边，提高可读性

5. **验证优化**
   - 检查数据流类型兼容性
   - 验证必填参数完整性
   - 发现潜在逻辑错误

### 9.3 使用流程

```
1. 用户输入需求
   ↓
2. 系统分析意图 (LLM)
   ↓
3. 匹配最佳实践模板
   ↓
4. 生成工作流结构 (LLM)
   ↓
5. 自动配置节点参数
   ↓
6. 计算节点布局位置
   ↓
7. 验证并优化工作流
   ↓
8. 展示预览和说明
   ↓
9. 用户确认或调整
   ↓
10. 创建工作流
```

### 9.4 最佳实践模板库

系统内置以下工作流模板：

| 模板名称 | 适用场景 | 节点模式 |
|---------|---------|---------|
| **内容生成** | 文章、大纲、文案 | Start → LLM → Code → If → End |
| **数据提取** | 简历、文档解析 | Start → LLM → Code → Database → End |
| **RAG 问答** | 智能客服、知识问答 | Start → Knowledge → Code → LLM → End |
| **批量处理** | 翻译、生成、转换 | Start → Database → Loop[LLM] → End |
| **API 集成** | 第三方服务调用 | Start → Plugin → Code → If → End |
| **数据管道** | ETL、数据迁移 | Start → Database → Code → Loop → Database → End |

### 9.5 与现有功能协同

```
┌─────────────────────────────────────────────────────────────┐
│                    工作流创建方式                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 🆕 零开始自动生成                                         │
│     "我想创建..."  →  [完整工作流]                           │
│                                                              │
│  2. 模板导入 + 手动调整                                       │
│     选择模板  →  [基础结构]  →  手动编辑                     │
│                                                              │
│  3. 渐进式构建 + 推荐辅助                                     │
│     空画布  →  添加节点  →  [推荐]  →  继续添加              │
│                                                              │
│  三种方式可自由组合使用                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 附录

### A. 完整的规则配置示例

见单独文件：`rules.yaml`

### B. LLM Prompt 完整模板

见单独文件：`prompts.yaml`

### C. API 接口完整文档

见单独文件：`API_SPECIFICATION.md`

### D. 测试用例

见单独文件：`TESTING_GUIDE.md`

---

**文档更新记录**

| 日期 | 版本 | 修改内容 | 作者 |
|------|------|----------|------|
| 2025-12-11 | V1.0 | 初始版本 | Coze Studio Team |
| 2025-12-12 | V1.1 | 🆕 新增"从零开始的智能工作流构建"功能：<br>- 场景4：基于需求自动生成完整工作流<br>- 意图分析和模板匹配<br>- 自动布局和验证引擎<br>- 前端生成器组件<br>- 最佳实践模板库 | Coze Studio Team |

