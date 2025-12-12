/*
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package recommendation 提供工作流节点推荐功能
// 支持基于规则、LLM 和统计分析的多策略推荐
package recommendation

import (
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

// RecommendRequest 节点推荐请求
type RecommendRequest struct {
	// 工作流 ID
	WorkflowID string `json:"workflow_id"`

	// 源节点信息
	SourceNodeID   string          `json:"source_node_id"`
	SourceNodeType entity.NodeType `json:"source_node_type"`
	SourceNodeName string          `json:"source_node_name"`

	// 源节点输出信息
	SourceOutputs map[string]*vo.TypeInfo `json:"source_outputs"`
	OutputFormat  string                  `json:"output_format,omitempty"` // json, text, markdown

	// 工作流上下文
	WorkflowContext *WorkflowContext `json:"workflow_context"`

	// 推荐配置
	Limit         int  `json:"limit"`          // 返回推荐数量，默认 10
	IncludeReason bool `json:"include_reason"` // 是否包含推荐理由

	// 用户信息（用于个性化推荐）
	UserID    string `json:"user_id,omitempty"`
	SpaceID   string `json:"space_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// WorkflowContext 工作流上下文信息
type WorkflowContext struct {
	// 工作流基本信息
	WorkflowName      string `json:"workflow_name,omitempty"`
	WorkflowObjective string `json:"workflow_objective,omitempty"` // 工作流目标描述

	// 已有节点信息
	ExistingNodes []NodeSummary `json:"existing_nodes"`
	TotalNodes    int           `json:"total_nodes"`

	// 工作流特征
	HasDatabaseNode  bool   `json:"has_database_node"`
	HasKnowledgeNode bool   `json:"has_knowledge_node"`
	HasLoopNode      bool   `json:"has_loop_node"`
	InsideLoop       bool   `json:"inside_loop"`        // 当前是否在循环内部
	WorkflowComplexity string `json:"workflow_complexity"` // low, medium, high

	// 数据流特征
	HasArrayOutput     bool `json:"has_array_output"`
	HasLargeDataset    bool `json:"has_large_dataset"`
	HasRepeatedPattern bool `json:"has_repeated_pattern"`

	// 场景标识
	IsChatWorkflow bool   `json:"is_chat_workflow"`
	ScenarioType   string `json:"scenario_type,omitempty"` // content_generation, data_processing, etc.
}

// NodeSummary 节点摘要信息
type NodeSummary struct {
	ID          string          `json:"id"`
	Type        entity.NodeType `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
}

// RecommendResponse 节点推荐响应
type RecommendResponse struct {
	// 推荐结果列表
	Recommendations []RecommendedNode `json:"recommendations"`

	// 各策略的推荐结果（调试用）
	Strategies []StrategyResult `json:"strategies,omitempty"`

	// 推荐元信息
	Meta RecommendationMeta `json:"meta,omitempty"`
}

// RecommendedNode 推荐的节点
type RecommendedNode struct {
	// 节点基本信息
	NodeType    entity.NodeType `json:"node_type"`
	DisplayName string          `json:"display_name"`
	Category    string          `json:"category"` // AI处理, 数据处理, 流程控制等

	// 推荐得分和理由
	Score  float64 `json:"score"`  // 0-1 之间
	Reason string  `json:"reason"` // 推荐理由

	// UI 展示信息
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`

	// 节点模板（可选，用于快速创建）
	Template *vo.Node `json:"template,omitempty"`

	// 配置建议
	SuggestedConfig map[string]interface{} `json:"suggested_config,omitempty"`

	// 推荐来源
	Source string `json:"source,omitempty"` // rule_engine, llm, statistics
}

// StrategyResult 单个策略的推荐结果
type StrategyResult struct {
	Name         string            `json:"name"`   // 策略名称
	Weight       float64           `json:"weight"` // 策略权重
	Nodes        []RecommendedNode `json:"nodes"`
	ExecutionTime int64            `json:"execution_time_ms,omitempty"` // 执行时间（毫秒）
}

// RecommendationMeta 推荐元信息
type RecommendationMeta struct {
	TotalCandidates   int    `json:"total_candidates"`    // 候选节点总数
	FilteredCount     int    `json:"filtered_count"`      // 过滤掉的数量
	RecommendationID  string `json:"recommendation_id"`   // 推荐 ID（用于反馈）
	ExecutionTimeMs   int64  `json:"execution_time_ms"`   // 总执行时间
	StrategiesUsed    []string `json:"strategies_used"`   // 使用的策略列表
}

// RecommendationFeedback 推荐反馈
type RecommendationFeedback struct {
	WorkflowID       string          `json:"workflow_id"`
	RecommendationID string          `json:"recommendation_id"`
	SourceNodeID     string          `json:"source_node_id"`
	RecommendedType  entity.NodeType `json:"recommended_type"`
	SelectedType     entity.NodeType `json:"selected_type"`
	IsUseful         bool            `json:"is_useful"` // 用户是否采纳推荐
	UserComment      string          `json:"user_comment,omitempty"`
	Timestamp        int64           `json:"timestamp"`
}

// RecommendationRule 推荐规则定义（对应 YAML 配置）
type RecommendationRule struct {
	Name          string                 `yaml:"name"`
	Priority      string                 `yaml:"priority"` // high, medium, low
	Condition     RuleCondition          `yaml:"condition"`
	Recommendations []RuleRecommendation `yaml:"recommendations"`
}

// RuleCondition 规则条件
type RuleCondition struct {
	// 源节点条件
	SourceNodeType  *entity.NodeType  `yaml:"source_node_type,omitempty"`
	SourceNodeTypes []entity.NodeType `yaml:"source_node_types,omitempty"`

	// 输出条件
	OutputFormat      *string `yaml:"output_format,omitempty"`       // json, text, markdown
	OutputHasArray    *bool   `yaml:"output_has_array,omitempty"`
	SourceOutputIsArray *bool `yaml:"source_output_is_array,omitempty"`

	// 工作流上下文条件
	WorkflowHasDB      *bool   `yaml:"workflow_has_database,omitempty"`
	WorkflowComplexity *string `yaml:"workflow_complexity,omitempty"` // low, medium, high
	InsideLoop         *bool   `yaml:"inside_loop,omitempty"`
	IsChatWorkflow     *bool   `yaml:"is_chat_workflow,omitempty"`

	// 操作类型
	OperationType *string `yaml:"operation_type,omitempty"` // query, insert, update, delete

	// 特殊条件
	IsBranchEnd        *bool `yaml:"is_branch_end,omitempty"`
	HasRepeatedPattern *bool `yaml:"has_repeated_pattern,omitempty"`
	HasLargeDataset    *bool `yaml:"has_large_dataset,omitempty"`

	// 通配符
	AlwaysMatch *bool `yaml:"always_match,omitempty"`
}

// RuleRecommendation 规则推荐项
type RuleRecommendation struct {
	NodeType        entity.NodeType        `yaml:"node_type"`
	Score           float64                `yaml:"score"`
	Reason          string                 `yaml:"reason"`
	Category        string                 `yaml:"category"`
	SuggestedConfig map[string]interface{} `yaml:"suggested_config,omitempty"`
}

// NodeCategory 节点分类常量
const (
	CategoryAIProcessing   = "AI 处理"
	CategoryDataProcessing = "数据处理"
	CategoryDataStorage    = "数据存储"
	CategoryDataRetrieval  = "数据检索"
	CategoryFlowControl    = "流程控制"
	CategoryExternal       = "外部服务"
	CategoryOrganization   = "流程组织"
	CategoryPerformance    = "性能优化"
)

// Priority 优先级常量
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// Strategy 推荐策略接口
type Strategy interface {
	// Name 返回策略名称
	Name() string

	// Recommend 执行推荐逻辑
	Recommend(req *RecommendRequest) ([]RecommendedNode, error)

	// Weight 返回策略权重
	Weight() float64
}
