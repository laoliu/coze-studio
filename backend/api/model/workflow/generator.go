/*
 * Copyright 2025 coze-dev Authors
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

import (
	"github.com/coze-dev/coze-studio/backend/api/model/base"
)

// GenerateWorkflowRequest 工作流生成请求
type GenerateWorkflowRequest struct {
	base.Base
	SpaceID          string   `json:"space_id" vd:"$!=''"`         // 空间 ID
	UserRequirement  string   `json:"user_requirement" vd:"$!=''"`// 用户需求描述
	RequirementType  string   `json:"requirement_type"`            // 需求类型（可选）
	Language         string   `json:"language"`                    // 语言偏好，默认 zh-CN
	MaxNodes         int      `json:"max_nodes"`                   // 最大节点数量（可选）
	AllowedNodeTypes []string `json:"allowed_node_types"`          // 允许的节点类型（可选）
}

// GeneratedNodeInfo 节点信息（API 层定义，避免循环依赖）
type GeneratedNodeInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Position    *NodePosition          `json:"position"`
	Config      map[string]interface{} `json:"config"`
}

// NodePosition 节点位置
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// GeneratedEdgeInfo 边信息（API 层定义，避免循环依赖）
type GeneratedEdgeInfo struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// GenerateWorkflowResponse 工作流生成响应
type GenerateWorkflowResponse struct {
	base.BaseResp
	WorkflowID       string                 `json:"workflow_id"`         // 工作流 ID
	WorkflowName     string                 `json:"workflow_name"`       // 工作流名称
	Nodes            []*GeneratedNodeInfo   `json:"nodes"`               // 生成的节点列表
	Edges            []*GeneratedEdgeInfo   `json:"edges"`               // 节点连接关系
	Explanation      string            `json:"explanation"`         // 整体流程说明
	NodeExplanations map[string]string `json:"node_explanations"`   // 每个节点的说明
	Confidence       float64           `json:"confidence"`          // 方案置信度 0-1
	EstimatedCost    *CostEstimate     `json:"estimated_cost"`      // 预估成本（可选）
}

// CostEstimate 成本估算
type CostEstimate struct {
	TotalCost  float64            `json:"total_cost"`  // 总成本
	TokenCount int                `json:"token_count"` // 预估 token 数
	Breakdown  map[string]float64 `json:"breakdown"`   // 成本明细
}
