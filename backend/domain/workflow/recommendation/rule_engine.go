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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"gopkg.in/yaml.v3"
)

// RuleEngine 基于规则的推荐引擎
type RuleEngine struct {
	rules  []RecommendationRule
	name   string
	weight float64
}

// NewRuleEngine 创建规则引擎实例
func NewRuleEngine() (*RuleEngine, error) {
	engine := &RuleEngine{
		name:   "RuleEngine",
		weight: 0.4, // 规则引擎权重 40%
	}

	// 加载规则配置
	if err := engine.loadRules(); err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}

	return engine, nil
}

// loadRules 从 YAML 文件加载推荐规则
func (e *RuleEngine) loadRules() error {
	// 获取规则文件路径（支持多种路径）
	possiblePaths := []string{
		"rules.yaml", // 当前目录（测试时）
		filepath.Join("backend", "domain", "workflow", "recommendation", "rules.yaml"), // 从项目根目录
		filepath.Join("..", "rules.yaml"),                                              // 从 examples 子目录
	}

	var data []byte
	var err error

	// 尝试各个可能的路径
	for _, rulesPath := range possiblePaths {
		data, err = os.ReadFile(rulesPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	// 解析 YAML
	var config struct {
		Rules []RecommendationRule `yaml:"recommendation_rules"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse rules YAML: %w", err)
	}

	e.rules = config.Rules
	return nil
}

// Name 返回策略名称
func (e *RuleEngine) Name() string {
	return e.name
}

// Weight 返回策略权重
func (e *RuleEngine) Weight() float64 {
	return e.weight
}

// Recommend 执行规则匹配和推荐
func (e *RuleEngine) Recommend(req *RecommendRequest) ([]RecommendedNode, error) {
	startTime := time.Now()
	var results []RecommendedNode

	// 遍历所有规则，找到匹配的规则
	for _, rule := range e.rules {
		if e.matchCondition(rule.Condition, req) {
			// 将规则中的推荐项转换为标准格式
			for _, rec := range rule.Recommendations {
				results = append(results, RecommendedNode{
					NodeType:        rec.NodeType,
					Score:           rec.Score,
					Reason:          rec.Reason,
					Category:        rec.Category,
					DisplayName:     getNodeDisplayName(rec.NodeType),
					Icon:            getNodeIcon(rec.NodeType),
					Description:     getNodeDescription(rec.NodeType),
					SuggestedConfig: rec.SuggestedConfig,
					Source:          "rule_engine",
				})
			}
		}
	}

	// 按优先级和得分排序
	results = e.sortAndDedup(results)

	// 记录执行时间
	executionTime := time.Since(startTime).Milliseconds()
	_ = executionTime // 可用于监控

	return results, nil
}

// matchCondition 检查规则条件是否匹配
func (e *RuleEngine) matchCondition(cond RuleCondition, req *RecommendRequest) bool {
	ctx := req.WorkflowContext

	// 检查通配符（总是匹配）
	if cond.AlwaysMatch != nil && *cond.AlwaysMatch {
		return true
	}

	// 检查源节点类型（单个）
	if cond.SourceNodeType != nil && *cond.SourceNodeType != req.SourceNodeType {
		return false
	}

	// 检查源节点类型（多个）
	if len(cond.SourceNodeTypes) > 0 {
		matched := false
		for _, t := range cond.SourceNodeTypes {
			if t == req.SourceNodeType {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查输出格式
	if cond.OutputFormat != nil && *cond.OutputFormat != req.OutputFormat {
		return false
	}

	// 检查是否输出数组
	if cond.OutputHasArray != nil {
		hasArray := e.checkHasArrayOutput(req.SourceOutputs)
		if *cond.OutputHasArray != hasArray {
			return false
		}
	}

	if cond.SourceOutputIsArray != nil {
		hasArray := e.checkHasArrayOutput(req.SourceOutputs)
		if *cond.SourceOutputIsArray != hasArray {
			return false
		}
	}

	// 检查工作流上下文条件
	if ctx != nil {
		// 检查是否有数据库节点
		if cond.WorkflowHasDB != nil && *cond.WorkflowHasDB != ctx.HasDatabaseNode {
			return false
		}

		// 检查工作流复杂度
		if cond.WorkflowComplexity != nil && *cond.WorkflowComplexity != ctx.WorkflowComplexity {
			return false
		}

		// 检查是否在循环内部
		if cond.InsideLoop != nil && *cond.InsideLoop != ctx.InsideLoop {
			return false
		}

		// 检查是否是聊天工作流
		if cond.IsChatWorkflow != nil && *cond.IsChatWorkflow != ctx.IsChatWorkflow {
			return false
		}

		// 检查是否有重复模式
		if cond.HasRepeatedPattern != nil && *cond.HasRepeatedPattern != ctx.HasRepeatedPattern {
			return false
		}

		// 检查是否有大数据集
		if cond.HasLargeDataset != nil && *cond.HasLargeDataset != ctx.HasLargeDataset {
			return false
		}
	}

	// 所有条件都匹配
	return true
}

// checkHasArrayOutput 检查输出中是否包含数组类型
func (e *RuleEngine) checkHasArrayOutput(outputs map[string]*vo.TypeInfo) bool {
	if outputs == nil {
		return false
	}

	for _, typeInfo := range outputs {
		if typeInfo != nil && typeInfo.Type == vo.DataTypeArray {
			return true
		}
	}

	return false
}

// sortAndDedup 对推荐结果排序和去重
func (e *RuleEngine) sortAndDedup(results []RecommendedNode) []RecommendedNode {
	if len(results) == 0 {
		return results
	}

	// 去重：同一个节点类型只保留得分最高的
	deduped := make(map[entity.NodeType]RecommendedNode)
	for _, rec := range results {
		if existing, exists := deduped[rec.NodeType]; exists {
			// 保留得分更高的
			if rec.Score > existing.Score {
				deduped[rec.NodeType] = rec
			}
		} else {
			deduped[rec.NodeType] = rec
		}
	}

	// 转换为切片
	var sorted []RecommendedNode
	for _, rec := range deduped {
		sorted = append(sorted, rec)
	}

	// 按得分降序排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Score > sorted[i].Score {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// getNodeDisplayName 获取节点的显示名称
func getNodeDisplayName(nodeType entity.NodeType) string {
	meta := entity.NodeMetaByNodeType(nodeType)
	if meta != nil {
		return meta.Name
	}
	return string(nodeType)
}

// getNodeIcon 获取节点图标
func getNodeIcon(nodeType entity.NodeType) string {
	meta := entity.NodeMetaByNodeType(nodeType)
	if meta != nil {
		return meta.IconURI
	}
	return ""
}

// getNodeDescription 获取节点描述
func getNodeDescription(nodeType entity.NodeType) string {
	meta := entity.NodeMetaByNodeType(nodeType)
	if meta != nil {
		return meta.Desc
	}
	return ""
}
