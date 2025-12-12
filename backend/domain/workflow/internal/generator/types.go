package generator

import (
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
)

// WorkflowGenerationRequest 工作流生成请求
type WorkflowGenerationRequest struct {
	UserRequirement  string        `json:"user_requirement"`   // 用户需求描述
	RequirementType  string        `json:"requirement_type"`   // 需求类型：text_generation, data_processing, api_integration, etc.
	Constraints      *Constraints  `json:"constraints"`        // 约束条件
	UserID           string        `json:"user_id"`
	Language         string        `json:"language"`           // 用户语言偏好，默认 zh-CN
}

// Constraints 约束条件
type Constraints struct {
	MaxNodes         int               `json:"max_nodes"`          // 最大节点数量，默认 10
	AllowedNodeTypes []entity.NodeType `json:"allowed_node_types"` // 允许的节点类型
	Budget           *BudgetLimit      `json:"budget"`             // 成本预算
	Performance      *PerformanceReq   `json:"performance"`        // 性能要求
}

// BudgetLimit 成本预算限制
type BudgetLimit struct {
	MaxTokens int     `json:"max_tokens"` // 最大 token 数
	MaxCost   float64 `json:"max_cost"`   // 最大成本（美元）
}

// PerformanceReq 性能要求
type PerformanceReq struct {
	MaxExecutionTime int `json:"max_execution_time"` // 最大执行时间（秒）
}

// WorkflowGenerationResponse 工作流生成响应
type WorkflowGenerationResponse struct {
	WorkflowID       string            `json:"workflow_id"`
	Nodes            []*InternalNode   `json:"nodes"`              // 生成的节点列表
	Edges            []*InternalEdge   `json:"edges"`              // 节点连接关系
	Explanation      string            `json:"explanation"`        // 整体流程说明
	NodeExplanations map[string]string `json:"node_explanations"`  // 每个节点的说明
	EstimatedCost    *CostEstimate     `json:"estimated_cost"`     // 预估成本
	Confidence       float64           `json:"confidence"`         // 方案置信度 0-1
}

// CostEstimate 成本估算
type CostEstimate struct {
	Total      float64            `json:"total"`       // 总成本
	Breakdown  map[string]float64 `json:"breakdown"`   // 成本明细
	TokenCount int                `json:"token_count"` // 预估 token 数
}

// IntentAnalysisResult 意图分析结果
type IntentAnalysisResult struct {
	Intent              *Intent             `json:"intent"`
	Inputs              []*IOParam          `json:"inputs"`
	Outputs             []*IOParam          `json:"outputs"`
	KeySteps            []string            `json:"key_steps"`
	Constraints         *IntentConstraints  `json:"constraints"`
	SuggestedNodes      []*SuggestedNode    `json:"suggested_nodes"`
}

// Intent 意图信息
type Intent struct {
	CoreObjective string  `json:"core_objective"` // 核心目标
	WorkflowType  string  `json:"workflow_type"`  // 工作流类型
	Complexity    string  `json:"complexity"`     // 复杂度：simple, medium, complex
	Confidence    float64 `json:"confidence"`     // 置信度
}

// IOParam 输入/输出参数
type IOParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// IntentConstraints 意图约束
type IntentConstraints struct {
	Quality     string `json:"quality"`     // 质量要求：low, medium, high
	Cost        string `json:"cost"`        // 成本要求：low, medium, high
	Performance string `json:"performance"` // 性能要求：low, medium, high
}

// SuggestedNode 建议的节点
type SuggestedNode struct {
	Type   entity.NodeType `json:"type"`
	Reason string          `json:"reason"`
}

// RawWorkflow LLM 生成的原始工作流结构
type RawWorkflow struct {
	WorkflowName        string     `json:"workflow_name"`
	Description         string     `json:"description"`
	Nodes               []*RawNode `json:"nodes"`
	Edges               []*RawEdge `json:"edges"`
	OverallExplanation  string     `json:"overall_explanation"`
	Confidence          float64    `json:"confidence"`
}

// RawNode 原始节点定义
type RawNode struct {
	ID          string                 `json:"id"`
	Type        entity.NodeType        `json:"type"`
	Name        string                 `json:"name"`
	Position    *Position              `json:"position"`
	Config      map[string]interface{} `json:"config"`
	Explanation string                 `json:"explanation"`
}

// Position 节点位置（用于布局）
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// RawEdge 原始边定义
type RawEdge struct {
	From      string `json:"from"`
	To        string `json:"to"`
	OutputKey string `json:"output_key"`
	InputKey  string `json:"input_key"`
}

// ConfiguredWorkflow 配置后的工作流（使用简化的内部表示）
type ConfiguredWorkflow struct {
	Nodes            []*InternalNode
	Edges            []*InternalEdge
	Explanations     map[string]string
	ConfidenceScore  float64
}

// InternalNode 生成器内部使用的简化节点表示
type InternalNode struct {
	ID          string
	Type        string
	Name        string
	Position    *Position
	Config      map[string]interface{}
}

// InternalEdge 生成器内部使用的简化边表示
type InternalEdge struct {
	From      string
	To        string
	OutputKey string
	InputKey  string
}

// WorkflowExplanations 工作流说明
type WorkflowExplanations struct {
	Overall string            `json:"overall"`
	PerNode map[string]string `json:"per_node"`
}

// WorkflowTemplate 工作流模板
type WorkflowTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	UseCases    []string               `json:"use_cases"`
	Pattern     *TemplatePattern       `json:"pattern"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TemplatePattern 模板模式
type TemplatePattern struct {
	Nodes []*TemplateNode `json:"nodes"`
}

// TemplateNode 模板节点
type TemplateNode struct {
	Type     entity.NodeType        `json:"type"`
	Name     string                 `json:"name"`
	Function string                 `json:"function"`
	Config   map[string]interface{} `json:"config"`
}

// ValidationError 验证错误
type ValidationError struct {
	NodeID   string `json:"node_id"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // error, warning, info
}
