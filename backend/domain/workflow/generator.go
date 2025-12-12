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
	"context"

	"github.com/cloudwego/eino/components/model"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/generator"
)

// WorkflowGenerationRequest 工作流生成请求（公开接口）
type WorkflowGenerationRequest = generator.WorkflowGenerationRequest

// WorkflowGenerationResponse 工作流生成响应（公开接口）
type WorkflowGenerationResponse = generator.WorkflowGenerationResponse

// Constraints 约束条件（公开接口）
type Constraints = generator.Constraints

// IntentAnalysisResult 意图分析结果（公开接口）
type IntentAnalysisResult = generator.IntentAnalysisResult

// NodeInfo 节点信息（公开接口）
type NodeInfo = generator.InternalNode

// EdgeInfo 边信息（公开接口）
type EdgeInfo = generator.InternalEdge

// Position 节点位置（公开接口）
type Position = generator.Position

// CostEstimate 成本估算（公开接口）
type CostEstimate = generator.CostEstimate

// WorkflowGenerator 工作流生成器接口（公开）
type WorkflowGenerator interface {
	GenerateWorkflow(ctx context.Context, req *WorkflowGenerationRequest) (*WorkflowGenerationResponse, error)
}

// NewWorkflowGenerator 创建工作流生成器（公开工厂函数）
func NewWorkflowGenerator(llmClient model.BaseChatModel) (WorkflowGenerator, error) {
	return generator.NewWorkflowGenerator(llmClient)
}
