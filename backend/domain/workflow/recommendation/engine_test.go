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

package recommendation

import (
	"context"
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

func TestEngineBasicRecommendation(t *testing.T) {
	// 创建推荐引擎
	engine, err := NewEngine(nil)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// 测试场景 1: LLM 节点输出 JSON 后的推荐
	t.Run("LLM JSON Output", func(t *testing.T) {
		req := &RecommendRequest{
			WorkflowID:     "test-workflow-1",
			SourceNodeID:   "llm-node-1",
			SourceNodeType: entity.NodeTypeLLM,
			SourceNodeName: "内容生成",
			OutputFormat:   "json",
			WorkflowContext: &WorkflowContext{
				TotalNodes:      2,
				HasDatabaseNode: false,
			},
			Limit:         5,
			IncludeReason: true,
		}

		resp, err := engine.Recommend(context.Background(), req)
		if err != nil {
			t.Fatalf("Recommendation failed: %v", err)
		}

		// 验证结果
		if len(resp.Recommendations) == 0 {
			t.Error("Expected recommendations, got none")
		}

		// 打印所有推荐
		for _, rec := range resp.Recommendations {
			t.Logf("Recommendation: %s (%.2f) - %s", rec.NodeType, rec.Score, rec.Reason)
		}

		// 应该推荐 Code 节点
		foundCode := false
		for _, rec := range resp.Recommendations {
			if rec.NodeType == entity.NodeTypeCodeRunner {
				foundCode = true
				t.Logf("Found Code recommendation: score=%.2f, reason=%s",
					rec.Score, rec.Reason)
				break
			}
		}

		if !foundCode {
			t.Error("Expected Code node recommendation for LLM JSON output")
		}
	})

	// 测试场景 2: Code 节点输出数组后的推荐
	t.Run("Code Array Output", func(t *testing.T) {
		req := &RecommendRequest{
			WorkflowID:     "test-workflow-2",
			SourceNodeID:   "code-node-1",
			SourceNodeType: entity.NodeTypeCodeRunner,
			SourceNodeName: "数据处理",
			SourceOutputs: map[string]*vo.TypeInfo{
				"items": {Type: vo.DataTypeArray},
			},
			WorkflowContext: &WorkflowContext{
				TotalNodes:     3,
				HasArrayOutput: true,
			},
			Limit: 5,
		}

		resp, err := engine.Recommend(context.Background(), req)
		if err != nil {
			t.Fatalf("Recommendation failed: %v", err)
		}

		// 应该推荐 Loop 节点
		foundLoop := false
		for _, rec := range resp.Recommendations {
			if rec.NodeType == entity.NodeTypeLoop {
				foundLoop = true
				t.Logf("Found Loop recommendation: score=%.2f", rec.Score)
				break
			}
		}

		if !foundLoop {
			t.Error("Expected Loop node recommendation for array output")
		}
	})

	// 测试场景 3: 开始节点后的推荐
	t.Run("Start Node", func(t *testing.T) {
		req := &RecommendRequest{
			WorkflowID:     "test-workflow-3",
			SourceNodeID:   "start-node",
			SourceNodeType: entity.NodeTypeEntry,
			SourceNodeName: "开始",
			WorkflowContext: &WorkflowContext{
				TotalNodes:     1,
				IsChatWorkflow: false,
			},
			Limit: 5,
		}

		resp, err := engine.Recommend(context.Background(), req)
		if err != nil {
			t.Fatalf("Recommendation failed: %v", err)
		}

		if len(resp.Recommendations) == 0 {
			t.Error("Expected recommendations for start node")
		}

		// 应该包含 LLM、Code 等常用节点
		nodeTypes := make(map[entity.NodeType]bool)
		for _, rec := range resp.Recommendations {
			nodeTypes[rec.NodeType] = true
			t.Logf("Recommendation: %s (%.2f)", rec.DisplayName, rec.Score)
		}

		if !nodeTypes[entity.NodeTypeLLM] && !nodeTypes[entity.NodeTypeCodeRunner] {
			t.Error("Expected common nodes (LLM or Code) for start node")
		}
	})
}

func TestRuleEngine(t *testing.T) {
	engine, err := NewRuleEngine()
	if err != nil {
		t.Fatalf("Failed to create rule engine: %v", err)
	}

	t.Run("Load Rules", func(t *testing.T) {
		if len(engine.rules) == 0 {
			t.Error("No rules loaded")
		}
		t.Logf("Loaded %d rules", len(engine.rules))
	})

	t.Run("Match Condition", func(t *testing.T) {
		req := &RecommendRequest{
			SourceNodeType: entity.NodeTypeLLM,
			OutputFormat:   "json",
			WorkflowContext: &WorkflowContext{
				HasDatabaseNode: false,
			},
		}

		recs, err := engine.Recommend(req)
		if err != nil {
			t.Fatalf("Recommendation failed: %v", err)
		}

		if len(recs) == 0 {
			t.Error("Expected recommendations from rule engine")
		}
	})
}

func TestScorer(t *testing.T) {
	ruleEngine, err := NewRuleEngine()
	if err != nil {
		t.Logf("Warning: Could not create rule engine: %v", err)
		// 创建一个简单的 mock 策略用于测试
		ruleEngine = &RuleEngine{
			name:   "RuleEngine",
			weight: 0.4,
			rules:  []RecommendationRule{},
		}
	}
	scorer := NewScorer(ruleEngine)

	t.Run("Merge Results", func(t *testing.T) {
		strategyResults := map[string][]RecommendedNode{
			"RuleEngine": {
				{NodeType: entity.NodeTypeCodeRunner, Score: 0.9, Reason: "Rule reason"},
				{NodeType: entity.NodeTypeLLM, Score: 0.7, Reason: "Another reason"},
			},
		}

		merged := scorer.MergeAndRank(strategyResults)

		if len(merged) != 2 {
			t.Errorf("Expected 2 merged results, got %d", len(merged))
		}

		// 验证排序
		if merged[0].Score < merged[1].Score {
			t.Error("Results not sorted by score")
		}
	})

	t.Run("Deduplication", func(t *testing.T) {
		strategyResults := map[string][]RecommendedNode{
			"RuleEngine": {
				{NodeType: entity.NodeTypeCodeRunner, Score: 0.8},
				{NodeType: entity.NodeTypeCodeRunner, Score: 0.9},
			},
		}

		merged := scorer.MergeAndRank(strategyResults)

		if len(merged) != 1 {
			t.Errorf("Expected 1 deduplicated result, got %d", len(merged))
		}

		// 应该保留得分更高的
		if merged[0].Score < 0.9 {
			t.Error("Should keep higher score in deduplication")
		}
	})
}

func BenchmarkRecommendation(b *testing.B) {
	engine, err := NewEngine(nil)
	if err != nil {
		b.Fatal(err)
	}

	req := &RecommendRequest{
		WorkflowID:     "bench-workflow",
		SourceNodeID:   "node-1",
		SourceNodeType: entity.NodeTypeLLM,
		OutputFormat:   "json",
		WorkflowContext: &WorkflowContext{
			TotalNodes: 5,
		},
		Limit: 10,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Recommend(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
