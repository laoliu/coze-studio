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
	"sort"
)

// Scorer 推荐结果评分和排序器
type Scorer struct {
	strategies []Strategy
}

// NewScorer 创建评分器实例
func NewScorer(strategies ...Strategy) *Scorer {
	return &Scorer{
		strategies: strategies,
	}
}

// MergeAndRank 合并多个策略的推荐结果并排序
func (s *Scorer) MergeAndRank(strategyResults map[string][]RecommendedNode) []RecommendedNode {
	// 存储所有候选节点，key 为节点类型
	candidates := make(map[string]*ScoredNode)

	// 合并各策略的推荐结果
	for strategyName, nodes := range strategyResults {
		weight := s.getStrategyWeight(strategyName)

		for _, node := range nodes {
			key := string(node.NodeType)

			if existing, exists := candidates[key]; exists {
				// 已存在，累加加权得分
				existing.TotalScore += node.Score * weight
				existing.StrategyScores[strategyName] = node.Score
				// 合并推荐理由
				if node.Reason != "" {
					// 检查是否已经存在相同的理由
					found := false
					for _, r := range existing.Reasons {
						if r == node.Reason {
							found = true
							break
						}
					}
					if !found {
						existing.Reasons = append(existing.Reasons, node.Reason)
					}
				}
			} else {
				// 新节点
				candidates[key] = &ScoredNode{
					Node:           node,
					TotalScore:     node.Score * weight,
					StrategyScores: map[string]float64{strategyName: node.Score},
					Reasons:        []string{node.Reason},
				}
			}
		}
	}

	// 转换为切片并计算最终得分
	var results []RecommendedNode
	for _, scored := range candidates {
		// 归一化得分（总权重可能不为 1）
		totalWeight := 0.0
		for strategyName := range scored.StrategyScores {
			totalWeight += s.getStrategyWeight(strategyName)
		}

		if totalWeight > 0 {
			scored.Node.Score = scored.TotalScore / totalWeight
		}

		// 合并推荐理由
		if len(scored.Reasons) > 1 {
			scored.Node.Reason = s.combineReasons(scored.Reasons)
		}

		results = append(results, scored.Node)
	}

	// 按得分降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// ScoredNode 带评分的节点
type ScoredNode struct {
	Node           RecommendedNode
	TotalScore     float64
	StrategyScores map[string]float64
	Reasons        []string
}

// getStrategyWeight 获取策略权重
func (s *Scorer) getStrategyWeight(strategyName string) float64 {
	for _, strategy := range s.strategies {
		if strategy.Name() == strategyName {
			return strategy.Weight()
		}
	}
	return 0.0
}

// combineReasons 合并多个推荐理由
func (s *Scorer) combineReasons(reasons []string) string {
	if len(reasons) == 0 {
		return ""
	}

	if len(reasons) == 1 {
		return reasons[0]
	}

	// 去重
	uniqueReasons := make(map[string]bool)
	var combined []string

	for _, reason := range reasons {
		if reason != "" && !uniqueReasons[reason] {
			uniqueReasons[reason] = true
			combined = append(combined, reason)
		}
	}

	// 合并为一句话（取第一个理由，或者拼接前两个）
	if len(combined) > 0 {
		if len(combined) == 1 {
			return combined[0]
		}
		// 返回最重要的理由
		return combined[0]
	}

	return reasons[0]
}

// FilterIncompatible 过滤不兼容的推荐
func (s *Scorer) FilterIncompatible(recommendations []RecommendedNode, req *RecommendRequest) []RecommendedNode {
	var filtered []RecommendedNode

	for _, rec := range recommendations {
		if s.isCompatible(rec, req) {
			filtered = append(filtered, rec)
		}
	}

	return filtered
}

// isCompatible 检查节点是否与当前上下文兼容
func (s *Scorer) isCompatible(node RecommendedNode, req *RecommendRequest) bool {
	// TODO: 实现更复杂的兼容性检查
	// 例如：
	// 1. 检查数据类型是否兼容
	// 2. 检查是否有必要的依赖（如数据库节点需要数据库配置）
	// 3. 检查工作流的约束条件

	// 目前简单返回 true，后续可扩展
	return true
}

// BoostScore 根据特定规则提升某些推荐的得分
func (s *Scorer) BoostScore(recommendations []RecommendedNode, req *RecommendRequest) []RecommendedNode {
	for i := range recommendations {
		// 提升得分的规则
		boost := 0.0

		// 1. 如果是常用节点，提升得分
		if s.isCommonlyUsed(string(recommendations[i].NodeType)) {
			boost += 0.05
		}

		// 2. 如果匹配用户历史偏好，提升得分
		if s.matchesUserPreference(string(recommendations[i].NodeType), req) {
			boost += 0.1
		}

		// 3. 如果是最佳实践推荐，提升得分
		if s.isBestPractice(string(recommendations[i].NodeType), string(req.SourceNodeType)) {
			boost += 0.08
		}

		// 应用提升（但不超过 1.0）
		recommendations[i].Score += boost
		if recommendations[i].Score > 1.0 {
			recommendations[i].Score = 1.0
		}
	}

	return recommendations
}

// isCommonlyUsed 检查是否是常用节点
func (s *Scorer) isCommonlyUsed(nodeType string) bool {
	// 定义常用节点列表
	commonNodes := map[string]bool{
		"LLM":      true,
		"Code":     true,
		"If":       true,
		"End":      true,
		"Database": true,
	}

	return commonNodes[nodeType]
}

// matchesUserPreference 检查是否匹配用户偏好
func (s *Scorer) matchesUserPreference(nodeType string, req *RecommendRequest) bool {
	// TODO: 从用户历史数据中学习偏好
	// 目前返回 false，后续可接入统计分析模块
	return false
}

// isBestPractice 检查是否是最佳实践组合
func (s *Scorer) isBestPractice(targetType string, sourceType string) bool {
	// 定义常见的最佳实践组合
	bestPractices := map[string][]string{
		"LLM":      {"Code", "Database", "If"},
		"Code":     {"LLM", "If", "Loop"},
		"Dataset":  {"LLM", "Code"},
		"Database": {"Code", "If", "LLM"},
		"Loop":     {"Code", "LLM", "Break"},
	}

	if targets, exists := bestPractices[sourceType]; exists {
		for _, target := range targets {
			if target == targetType {
				return true
			}
		}
	}

	return false
}
