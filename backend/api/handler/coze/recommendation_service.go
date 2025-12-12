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

package coze

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/recommendation"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	// 推荐引擎单例
	recommendEngine *recommendation.Engine
)

// initRecommendEngine 初始化推荐引擎
func initRecommendEngine() error {
	if recommendEngine != nil {
		return nil
	}

	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		return err
	}

	recommendEngine = engine
	logs.Infof("Node recommendation engine initialized successfully")
	return nil
}

// GetNodeRecommendations 获取节点推荐
// @router /api/workflow_api/node/recommend [POST]
func GetNodeRecommendations(ctx context.Context, c *app.RequestContext) {
	startTime := time.Now()

	// 确保引擎已初始化
	if recommendEngine == nil {
		if err := initRecommendEngine(); err != nil {
			logs.Errorf("Failed to initialize recommendation engine: %v", err)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to initialize recommendation engine",
			})
			return
		}
	}

	// 解析请求
	var req workflow.NodeRecommendationRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = 10
	}
	if !req.IncludeReason {
		req.IncludeReason = true
	}

	// 转换为推荐引擎请求
	recReq := convertToRecommendRequest(&req)

	// 执行推荐
	recResp, err := recommendEngine.Recommend(ctx, recReq)
	if err != nil {
		logs.Errorf("Recommendation failed: %v", err)
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Recommendation failed: " + err.Error(),
		})
		return
	}

	// 转换为 API 响应
	resp := convertToAPIResponse(recResp, time.Since(startTime))

	c.JSON(consts.StatusOK, resp)
}

// RecordNodeFeedback 记录用户反馈
// @router /api/workflow_api/node/recommend/feedback [POST]
func RecordNodeFeedback(ctx context.Context, c *app.RequestContext) {
	var req workflow.RecordFeedbackRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 确保引擎已初始化
	if recommendEngine == nil {
		if err := initRecommendEngine(); err != nil {
			logs.Errorf("Failed to initialize recommendation engine: %v", err)
			c.JSON(consts.StatusInternalServerError, workflow.RecordFeedbackResponse{
				Success: false,
				Message: "Engine not initialized",
			})
			return
		}
	}

	// 记录反馈
	feedback := &recommendation.RecommendationFeedback{
		RecommendationID: req.RequestID,
		SelectedType:     entity.NodeType(req.SelectedNode), // string → NodeType
		IsUseful:         req.UserAction == "selected",
		UserComment:      req.UserComment,
	}

	if err := recommendEngine.RecordFeedback(ctx, feedback); err != nil {
		logs.Errorf("Failed to record feedback: %v", err)
		c.JSON(consts.StatusInternalServerError, workflow.RecordFeedbackResponse{
			Success: false,
			Message: "Failed to record feedback",
		})
		return
	}

	c.JSON(consts.StatusOK, workflow.RecordFeedbackResponse{
		Success: true,
		Message: "Feedback recorded successfully",
	})
}

// convertToRecommendRequest 转换 API 请求为推荐引擎请求
func convertToRecommendRequest(req *workflow.NodeRecommendationRequest) *recommendation.RecommendRequest {
	recReq := &recommendation.RecommendRequest{
		WorkflowID:     req.WorkflowID,
		SourceNodeID:   req.SourceNodeID,
		SourceNodeType: entity.NodeType(req.SourceNodeType), // string → NodeType 类型转换
		SourceNodeName: req.SourceNodeName,
		SourceOutputs:  convertOutputsToTypeInfo(req.SourceOutputs), // 转换 outputs
		OutputFormat:   req.OutputFormat,
		Limit:          req.Limit,
		IncludeReason:  req.IncludeReason,
		UserID:         req.UserID,
		SpaceID:        req.SpaceID,
		ProjectID:      req.ProjectID,
	}

	// 转换工作流上下文
	if req.WorkflowContext != nil {
		recReq.WorkflowContext = &recommendation.WorkflowContext{
			WorkflowName:      req.WorkflowContext.WorkflowName,
			WorkflowObjective: req.WorkflowContext.WorkflowObjective,
			TotalNodes:        req.WorkflowContext.TotalNodes,
			HasDatabaseNode:   req.WorkflowContext.HasDatabaseNode,
			HasKnowledgeNode:  req.WorkflowContext.HasLLMNode, // HasLLMNode → HasKnowledgeNode
			HasLoopNode:       req.WorkflowContext.HasLoopNode,
			InsideLoop:        req.WorkflowContext.InsideLoop,
			IsChatWorkflow:    req.WorkflowContext.IsChatWorkflow,
		}

		// 转换已存在的节点
		if len(req.WorkflowContext.ExistingNodes) > 0 {
			recReq.WorkflowContext.ExistingNodes = make([]recommendation.NodeSummary, len(req.WorkflowContext.ExistingNodes))
			for i, node := range req.WorkflowContext.ExistingNodes {
				recReq.WorkflowContext.ExistingNodes[i] = recommendation.NodeSummary{
					ID:   node.NodeID,                    // NodeID → ID
					Type: entity.NodeType(node.NodeType), // string → NodeType 类型转换
					Name: node.NodeName,                  // NodeName → Name
				}
			}
		}
	}

	return recReq
}

// convertToAPIResponse 转换推荐引擎响应为 API 响应
func convertToAPIResponse(recResp *recommendation.RecommendResponse, executionTime time.Duration) *workflow.NodeRecommendationResponse {
	resp := &workflow.NodeRecommendationResponse{
		Recommendations: make([]workflow.RecommendedNodeInfo, len(recResp.Recommendations)),
		Metadata: workflow.RecommendationMetadata{
			RequestID:       recResp.Meta.RecommendationID, // Meta.RecommendationID → RequestID
			TotalCandidates: recResp.Meta.TotalCandidates,
			ExecutionTimeMs: executionTime.Milliseconds(),
			StrategiesUsed:  recResp.Meta.StrategiesUsed,
		},
	}

	// 转换推荐结果
	for i, rec := range recResp.Recommendations {
		resp.Recommendations[i] = workflow.RecommendedNodeInfo{
			NodeType:        string(rec.NodeType), // NodeType → string 类型转换
			DisplayName:     getNodeDisplayName(rec.NodeType),
			Score:           rec.Score,
			Reason:          rec.Reason,
			Category:        rec.Category,
			SuggestedConfig: rec.SuggestedConfig,
			StrategySource:  rec.Source, // Source → StrategySource
		}
	}

	return resp
}

// getNodeDisplayName 获取节点显示名称
func getNodeDisplayName(nodeType entity.NodeType) string {
	// TODO: 从节点元数据中获取国际化的显示名称
	// 这里只列出确定存在的几个基础类型
	displayNames := map[entity.NodeType]string{
		entity.NodeTypeLLM:         "大模型",
		entity.NodeTypeCodeRunner:  "代码",
		entity.NodeTypePlugin:      "插件",
		entity.NodeTypeLoop:        "循环",
		entity.NodeTypeSubWorkflow: "子工作流",
	}

	if name, ok := displayNames[nodeType]; ok {
		return name
	}
	// 其他节点类型直接返回字符串值
	return string(nodeType)
}

// convertOutputsToTypeInfo 转换输出格式从 map[string]interface{} 到 map[string]*vo.TypeInfo
// 这是一个简化的转换，实际使用中可能需要更详细的类型解析
func convertOutputsToTypeInfo(outputs map[string]interface{}) map[string]*vo.TypeInfo {
	if outputs == nil {
		return nil
	}

	// 对于简单场景，返回 nil
	// 实际使用中，推荐引擎主要依赖 OutputFormat 字段
	// 如果需要详细的类型信息，前端应该传递更结构化的数据
	return nil
}
