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
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/google/uuid"
)

// Engine 节点推荐引擎
// 整合多种推荐策略（规则引擎、LLM、统计分析）提供智能节点推荐
type Engine struct {
	// 推荐策略
	strategies []Strategy

	// 评分器
	scorer *Scorer

	// 配置
	config *EngineConfig

	// 日志
	logger logs.FormatLogger
}

// EngineConfig 引擎配置
type EngineConfig struct {
	// 是否启用 LLM 推荐
	EnableLLM bool

	// 是否启用统计推荐
	EnableStatistics bool

	// 默认返回推荐数量
	DefaultLimit int

	// 是否启用缓存
	EnableCache bool

	// 缓存过期时间（秒）
	CacheTTL int
}

// DefaultEngineConfig 默认引擎配置
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		EnableLLM:        false, // MVP 阶段暂不启用 LLM
		EnableStatistics: false, // MVP 阶段暂不启用统计
		DefaultLimit:     10,
		EnableCache:      false, // MVP 阶段暂不启用缓存
		CacheTTL:         300,   // 5 分钟
	}
}

// NewEngine 创建推荐引擎实例
func NewEngine(config *EngineConfig) (*Engine, error) {
	if config == nil {
		config = DefaultEngineConfig()
	}

	// 创建规则引擎（核心策略）
	ruleEngine, err := NewRuleEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create rule engine: %w", err)
	}

	strategies := []Strategy{ruleEngine}

	// 如果启用 LLM，添加 LLM 策略
	if config.EnableLLM {
		// TODO: 创建 LLM 策略
		// llmStrategy := NewLLMStrategy(llmConfig)
		// strategies = append(strategies, llmStrategy)
	}

	// 如果启用统计分析，添加统计策略
	if config.EnableStatistics {
		// TODO: 创建统计策略
		// statsStrategy := NewStatisticsStrategy()
		// strategies = append(strategies, statsStrategy)
	}

	// 创建评分器
	scorer := NewScorer(strategies...)

	return &Engine{
		strategies: strategies,
		scorer:     scorer,
		config:     config,
		logger:     logs.DefaultLogger(),
	}, nil
}

// Recommend 执行节点推荐
func (e *Engine) Recommend(ctx context.Context, req *RecommendRequest) (*RecommendResponse, error) {
	startTime := time.Now()

	// 生成推荐 ID
	recommendationID := uuid.New().String()

	// 验证请求
	if err := e.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = e.config.DefaultLimit
	}

	// 收集各策略的推荐结果
	strategyResults := make(map[string][]RecommendedNode)
	var strategies []StrategyResult

	for _, strategy := range e.strategies {
		strategyStartTime := time.Now()

		// 执行策略推荐
		nodes, err := strategy.Recommend(req)
		if err != nil {
			e.logger.Warnf("strategy %s failed: %v", strategy.Name(), err)
			continue
		}

		strategyResults[strategy.Name()] = nodes

		// 记录策略结果
		strategies = append(strategies, StrategyResult{
			Name:          strategy.Name(),
			Weight:        strategy.Weight(),
			Nodes:         nodes,
			ExecutionTime: time.Since(strategyStartTime).Milliseconds(),
		})
	}

	// 合并和排序推荐结果
	merged := e.scorer.MergeAndRank(strategyResults)

	// 过滤不兼容的推荐
	filtered := e.scorer.FilterIncompatible(merged, req)

	// 提升得分
	boosted := e.scorer.BoostScore(filtered, req)

	// 限制返回数量
	if len(boosted) > req.Limit {
		boosted = boosted[:req.Limit]
	}

	// 计算元信息
	totalCandidates := 0
	for _, nodes := range strategyResults {
		totalCandidates += len(nodes)
	}

	var strategiesUsed []string
	for _, s := range strategies {
		strategiesUsed = append(strategiesUsed, s.Name)
	}

	response := &RecommendResponse{
		Recommendations: boosted,
		Strategies:      strategies,
		Meta: RecommendationMeta{
			TotalCandidates:  totalCandidates,
			FilteredCount:    len(merged) - len(boosted),
			RecommendationID: recommendationID,
			ExecutionTimeMs:  time.Since(startTime).Milliseconds(),
			StrategiesUsed:   strategiesUsed,
		},
	}

	// 记录推荐日志
	e.logger.Infof("Recommendation completed: workflow=%s, source_node=%s, recommendations=%d, time=%dms",
		req.WorkflowID, req.SourceNodeID, len(boosted), response.Meta.ExecutionTimeMs)

	return response, nil
}

// validateRequest 验证推荐请求
func (e *Engine) validateRequest(req *RecommendRequest) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}

	if req.SourceNodeID == "" {
		return fmt.Errorf("source_node_id is required")
	}

	if req.SourceNodeType == "" {
		return fmt.Errorf("source_node_type is required")
	}

	return nil
}

// RecordFeedback 记录用户反馈
func (e *Engine) RecordFeedback(ctx context.Context, feedback *RecommendationFeedback) error {
	// TODO: 将反馈存储到数据库
	// 用于后续的统计分析和模型优化

	e.logger.Infof("Feedback recorded: recommendation_id=%s, useful=%v",
		feedback.RecommendationID, feedback.IsUseful)

	return nil
}

// GetStatistics 获取推荐统计信息
func (e *Engine) GetStatistics(ctx context.Context, workflowID string) (*RecommendationStatistics, error) {
	// TODO: 从数据库查询统计信息
	return &RecommendationStatistics{}, nil
}

// RecommendationStatistics 推荐统计信息
type RecommendationStatistics struct {
	TotalRecommendations int              `json:"total_recommendations"`
	AcceptedCount        int              `json:"accepted_count"`
	AcceptanceRate       float64          `json:"acceptance_rate"`
	TopRecommendedNodes  []NodeStatistics `json:"top_recommended_nodes"`
}

// NodeStatistics 节点统计信息
type NodeStatistics struct {
	NodeType       string  `json:"node_type"`
	RecommendCount int     `json:"recommend_count"`
	AcceptCount    int     `json:"accept_count"`
	AcceptanceRate float64 `json:"acceptance_rate"`
}
