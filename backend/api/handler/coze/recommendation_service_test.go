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
	"testing"

	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
)

func TestGetNodeRecommendations(t *testing.T) {
	// 初始化推荐引擎
	if err := initRecommendEngine(); err != nil {
		t.Skipf("Failed to initialize recommendation engine: %v", err)
		return
	}

	// 创建测试请求
	req := &workflow.NodeRecommendationRequest{
		WorkflowID:     "test-workflow",
		SourceNodeID:   "node-1",
		SourceNodeType: entity.NodeTypeLLM,
		OutputFormat:   "json",
		Limit:          5,
		IncludeReason:  true,
	}

	// 转换请求
	recReq := convertToRecommendRequest(req)

	// 验证转换
	if recReq.WorkflowID != req.WorkflowID {
		t.Errorf("WorkflowID mismatch: got %s, want %s", recReq.WorkflowID, req.WorkflowID)
	}

	if recReq.SourceNodeType != req.SourceNodeType {
		t.Errorf("SourceNodeType mismatch: got %s, want %s", recReq.SourceNodeType, req.SourceNodeType)
	}

	t.Logf("Request conversion successful: workflow=%s, source=%s, limit=%d",
		recReq.WorkflowID, recReq.SourceNodeType, recReq.Limit)
}

func TestGetNodeDisplayName(t *testing.T) {
	tests := []struct {
		nodeType entity.NodeType
		want     string
	}{
		{entity.NodeTypeLLM, "大模型"},
		{entity.NodeTypeCodeRunner, "代码"},
		{entity.NodeTypePlugin, "插件"},
		{entity.NodeTypeDatabase, "数据库"},
		{entity.NodeTypeIf, "条件判断"},
		{entity.NodeTypeLoop, "循环"},
		{entity.NodeTypeEnd, "结束"},
	}

	for _, tt := range tests {
		t.Run(string(tt.nodeType), func(t *testing.T) {
			got := getNodeDisplayName(tt.nodeType)
			if got != tt.want {
				t.Errorf("getNodeDisplayName(%s) = %s, want %s",
					tt.nodeType, got, tt.want)
			}
		})
	}
}

func TestConvertToAPIResponse(t *testing.T) {
	// 这是一个示例，展示如何使用 API 响应转换函数
	t.Log("API response conversion functions available")
}
