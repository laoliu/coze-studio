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

// Package recommendation 节点推荐系统
//
// 提供基于规则、LLM 和统计分析的智能节点推荐功能。
//
// 主要功能：
//   - 基于规则的快速推荐（规则引擎）
//   - 基于 LLM 的智能推荐（上下文理解）
//   - 基于历史数据的推荐（统计分析）
//   - 多策略融合和评分
//
// 使用示例：
//
//	// 创建推荐引擎
//	engine, err := recommendation.NewEngine(nil)
//	if err != nil {
//	    return err
//	}
//
//	// 构建推荐请求
//	req := &recommendation.RecommendRequest{
//	    WorkflowID:     "workflow-123",
//	    SourceNodeID:   "node-456",
//	    SourceNodeType: entity.NodeTypeLLM,
//	    SourceOutputs: map[string]*vo.TypeInfo{
//	        "output": {Type: entity.TypeString},
//	    },
//	    OutputFormat: "json",
//	    WorkflowContext: &recommendation.WorkflowContext{
//	        HasDatabaseNode: true,
//	        InsideLoop:      false,
//	    },
//	    Limit: 5,
//	}
//
//	// 执行推荐
//	resp, err := engine.Recommend(ctx, req)
//	if err != nil {
//	    return err
//	}
//
//	// 处理推荐结果
//	for _, rec := range resp.Recommendations {
//	    fmt.Printf("推荐节点: %s (得分: %.2f)\n", rec.DisplayName, rec.Score)
//	    fmt.Printf("推荐理由: %s\n", rec.Reason)
//	}
package recommendation

// Version 推荐系统版本
const Version = "1.0.0"

// MVP 阶段功能清单：
// ✅ 规则引擎（核心功能）
// ✅ 多规则匹配和评分
// ✅ 结果排序和去重
// ✅ 推荐 API 接口
// ⏳ LLM 智能推荐（待实现）
// ⏳ 统计分析推荐（待实现）
// ⏳ 缓存机制（待实现）
// ⏳ 用户反馈学习（待实现）
