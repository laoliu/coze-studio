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

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/recommendation"
)

// 示例 1: 基础推荐
func example1_BasicRecommendation() {
	fmt.Println("=== 示例 1: 基础节点推荐 ===\n")

	// 1. 创建推荐引擎
	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 构建推荐请求
	req := &recommendation.RecommendRequest{
		WorkflowID:     "demo-workflow-001",
		SourceNodeID:   "llm-node-1",
		SourceNodeType: entity.NodeTypeLLM,
		SourceNodeName: "内容生成器",
		OutputFormat:   "json",
		WorkflowContext: &recommendation.WorkflowContext{
			WorkflowName:      "AI 内容生成工作流",
			WorkflowObjective: "根据用户输入生成结构化内容",
			TotalNodes:        2,
			HasDatabaseNode:   false,
		},
		Limit:         5,
		IncludeReason: true,
	}

	// 3. 执行推荐
	resp, err := engine.Recommend(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 输出结果
	fmt.Printf("✨ 为节点 '%s' 找到 %d 个推荐\n\n", req.SourceNodeName, len(resp.Recommendations))

	for i, rec := range resp.Recommendations {
		fmt.Printf("%d. %s (%.0f%% 匹配)\n", i+1, rec.DisplayName, rec.Score*100)
		fmt.Printf("   分类: %s\n", rec.Category)
		fmt.Printf("   理由: %s\n", rec.Reason)
		fmt.Printf("   来源: %s\n", rec.Source)
		fmt.Println()
	}

	fmt.Printf("⚡ 执行时间: %dms\n", resp.Meta.ExecutionTimeMs)
	fmt.Printf("📊 候选总数: %d\n", resp.Meta.TotalCandidates)
	fmt.Println()
}

// 示例 2: 数组输出场景
func example2_ArrayOutputRecommendation() {
	fmt.Println("=== 示例 2: 数组输出推荐 ===\n")

	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		log.Fatal(err)
	}

	req := &recommendation.RecommendRequest{
		WorkflowID:     "demo-workflow-002",
		SourceNodeID:   "code-node-1",
		SourceNodeType: entity.NodeTypeCodeRunner,
		SourceNodeName: "数据解析器",
		SourceOutputs: map[string]*vo.TypeInfo{
			"items": {
				Type: vo.DataTypeArray,
				Desc: "解析后的数据项列表",
			},
			"count": {
				Type: vo.DataTypeNumber,
				Desc: "数据项数量",
			},
		},
		WorkflowContext: &recommendation.WorkflowContext{
			TotalNodes:     3,
			HasArrayOutput: true,
		},
		Limit: 5,
	}

	resp, err := engine.Recommend(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🔄 检测到数组输出，推荐使用循环处理\n\n")

	for i, rec := range resp.Recommendations {
		fmt.Printf("%d. %s (%.0f%%)\n", i+1, rec.DisplayName, rec.Score*100)
		fmt.Printf("   %s\n", rec.Reason)
		fmt.Println()
	}
}

// 示例 3: 工作流起始节点
func example3_StartNodeRecommendation() {
	fmt.Println("=== 示例 3: 工作流起始节点推荐 ===\n")

	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		log.Fatal(err)
	}

	req := &recommendation.RecommendRequest{
		WorkflowID:     "demo-workflow-003",
		SourceNodeID:   "start-node",
		SourceNodeType: entity.NodeTypeEntry,
		SourceNodeName: "开始",
		WorkflowContext: &recommendation.WorkflowContext{
			WorkflowName:      "新建工作流",
			WorkflowObjective: "处理用户查询并返回结果",
			TotalNodes:        1,
			IsChatWorkflow:    false,
		},
		Limit: 5,
	}

	resp, err := engine.Recommend(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🚀 为新工作流推荐第一个节点\n\n")

	for i, rec := range resp.Recommendations {
		fmt.Printf("%d. %s\n", i+1, rec.DisplayName)
		fmt.Printf("   📋 %s\n", rec.Category)
		fmt.Printf("   💡 %s\n", rec.Reason)
		fmt.Println()
	}
}

// 示例 4: 带数据库的场景
func example4_DatabaseWorkflowRecommendation() {
	fmt.Println("=== 示例 4: 数据库工作流推荐 ===\n")

	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		log.Fatal(err)
	}

	req := &recommendation.RecommendRequest{
		WorkflowID:     "demo-workflow-004",
		SourceNodeID:   "llm-node-2",
		SourceNodeType: entity.NodeTypeLLM,
		SourceNodeName: "文章生成器",
		OutputFormat:   "text",
		WorkflowContext: &recommendation.WorkflowContext{
			TotalNodes:      3,
			HasDatabaseNode: true,
		},
		Limit: 5,
	}

	resp, err := engine.Recommend(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("💾 检测到工作流包含数据库节点\n\n")

	for i, rec := range resp.Recommendations {
		fmt.Printf("%d. %s (%.1f)\n", i+1, rec.DisplayName, rec.Score)
		fmt.Printf("   %s\n", rec.Reason)
		fmt.Println()
	}
}

// 示例 5: 循环内部推荐
func example5_LoopInsideRecommendation() {
	fmt.Println("=== 示例 5: 循环内部推荐 ===\n")

	engine, err := recommendation.NewEngine(nil)
	if err != nil {
		log.Fatal(err)
	}

	req := &recommendation.RecommendRequest{
		WorkflowID:     "demo-workflow-005",
		SourceNodeID:   "loop-start",
		SourceNodeType: entity.NodeTypeLoop,
		SourceNodeName: "批量处理循环",
		WorkflowContext: &recommendation.WorkflowContext{
			TotalNodes: 4,
			InsideLoop: true,
		},
		Limit: 5,
	}

	resp, err := engine.Recommend(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🔁 循环内部节点推荐\n\n")

	for i, rec := range resp.Recommendations {
		fmt.Printf("%d. %s\n", i+1, rec.DisplayName)
		fmt.Printf("   %s\n", rec.Reason)
		fmt.Println()
	}
}

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Coze Studio 节点推荐系统 - 使用示例    ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	// 运行所有示例
	example1_BasicRecommendation()
	example2_ArrayOutputRecommendation()
	example3_StartNodeRecommendation()
	example4_DatabaseWorkflowRecommendation()
	example5_LoopInsideRecommendation()

	fmt.Println("✅ 所有示例运行完成！")
}
