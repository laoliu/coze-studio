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

package workflow

// NodeRecommendationRequest 节点推荐请求
type NodeRecommendationRequest struct {
	// 工作流 ID
	WorkflowID string `json:"workflow_id" binding:"required"`

	// 源节点信息
	SourceNodeID   string `json:"source_node_id" binding:"required"`
	SourceNodeType string `json:"source_node_type" binding:"required"` // 使用 string 避免循环依赖
	SourceNodeName string `json:"source_node_name,omitempty"`

	// 源节点输出信息（使用 map 避免循环依赖）
	SourceOutputs map[string]interface{} `json:"source_outputs,omitempty"`
	OutputFormat  string                 `json:"output_format,omitempty"` // json, text, markdown

	// 工作流上下文（可选）
	WorkflowContext *WorkflowContextInfo `json:"workflow_context,omitempty"`

	// 推荐配置
	Limit         int  `json:"limit,omitempty"`          // 返回推荐数量，默认 10
	IncludeReason bool `json:"include_reason,omitempty"` // 是否包含推荐理由，默认 true

	// 用户信息
	UserID    string `json:"user_id,omitempty"`
	SpaceID   string `json:"space_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// WorkflowContextInfo 工作流上下文信息（API 层）
type WorkflowContextInfo struct {
	WorkflowName      string            `json:"workflow_name,omitempty"`
	WorkflowObjective string            `json:"workflow_objective,omitempty"`
	TotalNodes        int               `json:"total_nodes,omitempty"`
	HasDatabaseNode   bool              `json:"has_database_node,omitempty"`
	HasLLMNode        bool              `json:"has_llm_node,omitempty"`
	HasLoopNode       bool              `json:"has_loop_node,omitempty"`
	InsideLoop        bool              `json:"inside_loop,omitempty"`
	IsChatWorkflow    bool              `json:"is_chat_workflow,omitempty"`
	ExistingNodes     []NodeSummaryInfo `json:"existing_nodes,omitempty"`
}

// NodeSummaryInfo 节点摘要信息（API 层）
type NodeSummaryInfo struct {
	NodeID   string `json:"node_id"`
	NodeType string `json:"node_type"` // 使用 string 避免循环依赖
	NodeName string `json:"node_name,omitempty"`
}

// NodeRecommendationResponse 节点推荐响应
type NodeRecommendationResponse struct {
	// 推荐结果
	Recommendations []RecommendedNodeInfo `json:"recommendations"`

	// 元数据
	Metadata RecommendationMetadata `json:"metadata"`
}

// RecommendedNodeInfo 推荐的节点信息（API 层）
type RecommendedNodeInfo struct {
	NodeType        string                 `json:"node_type"` // 使用 string 避免循环依赖
	DisplayName     string                 `json:"display_name"`
	Score           float64                `json:"score"`
	Reason          string                 `json:"reason,omitempty"`
	Category        string                 `json:"category,omitempty"`
	SuggestedConfig map[string]interface{} `json:"suggested_config,omitempty"`
	StrategySource  string                 `json:"strategy_source,omitempty"` // rule_engine, llm, statistics
}

// RecommendationMetadata 推荐元数据
type RecommendationMetadata struct {
	RequestID       string   `json:"request_id"`
	TotalCandidates int      `json:"total_candidates"`
	ExecutionTimeMs int64    `json:"execution_time_ms"`
	StrategiesUsed  []string `json:"strategies_used,omitempty"`
}

// RecordFeedbackRequest 记录用户反馈请求
type RecordFeedbackRequest struct {
	RequestID    string `json:"request_id" binding:"required"`
	SelectedNode string `json:"selected_node,omitempty"`        // 使用 string 避免循环依赖
	UserAction   string `json:"user_action" binding:"required"` // selected, dismissed, ignored
	UserComment  string `json:"user_comment,omitempty"`
}

// RecordFeedbackResponse 记录反馈响应
type RecordFeedbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
