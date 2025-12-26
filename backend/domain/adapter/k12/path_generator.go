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

package k12

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// PathGenerator 学习路径生成器
type PathGenerator struct {
}

// NewPathGenerator 创建学习路径生成器
func NewPathGenerator() *PathGenerator {
	return &PathGenerator{}
}

// GenerateLearningPath 生成学习路径
// targetKnowledgeID: 目标知识点ID
// masteredKnowledgeIDs: 已掌握的知识点ID列表
func (pg *PathGenerator) GenerateLearningPath(targetKnowledgeID string, masteredKnowledgeIDs []string) (*models.LearningPath, error) {
	// 获取目标知识点
	targetKnowledge := GetKnowledgePoint(targetKnowledgeID)
	if targetKnowledge == nil {
		return nil, fmt.Errorf("未找到知识点: %s", targetKnowledgeID)
	}

	// 构建已掌握知识点的映射
	masteredMap := make(map[string]bool)
	for _, id := range masteredKnowledgeIDs {
		masteredMap[id] = true
	}

	// 查找需要学习的前置知识
	requiredKnowledge := pg.findRequiredKnowledge(targetKnowledgeID, masteredMap)

	// 生成学习步骤
	var steps []models.LearningStep
	stepNumber := 1

	// 1. 添加前置知识的学习步骤
	for _, knowledgeID := range requiredKnowledge {
		knowledge := GetKnowledgePoint(knowledgeID)
		if knowledge == nil {
			continue
		}

		step := models.LearningStep{
			Order:         stepNumber,
			Knowledge:     *knowledge,
			Objectives:    []string{fmt.Sprintf("掌握前置知识：%s", knowledge.Name)},
			Resources:     GetResourcesByKnowledge(knowledgeID, ""),
			Exercises:     GetExercisesByKnowledge(knowledgeID, models.DifficultyBasic),
			EstimatedTime: pg.estimateStudyTime(knowledge.Difficulty),
		}
		steps = append(steps, step)
		stepNumber++
	}

	// 2. 添加目标知识点的学习步骤
	// 基础阶段
	basicStep := models.LearningStep{
		Order:         stepNumber,
		Knowledge:     *targetKnowledge,
		Objectives:    []string{fmt.Sprintf("掌握%s的基础概念和原理", targetKnowledge.Name)},
		Resources:     GetResourcesByKnowledge(targetKnowledgeID, models.DifficultyBasic),
		Exercises:     GetExercisesByKnowledge(targetKnowledgeID, models.DifficultyBasic),
		EstimatedTime: 30,
	}
	steps = append(steps, basicStep)
	stepNumber++

	// 提高阶段
	improveStep := models.LearningStep{
		Order:         stepNumber,
		Knowledge:     *targetKnowledge,
		Objectives:    []string{fmt.Sprintf("提高%s的应用能力", targetKnowledge.Name)},
		Resources:     GetResourcesByKnowledge(targetKnowledgeID, models.DifficultyImprove),
		Exercises:     GetExercisesByKnowledge(targetKnowledgeID, models.DifficultyImprove),
		EstimatedTime: 40,
	}
	steps = append(steps, improveStep)
	stepNumber++

	// 拓展阶段（可选）
	advancedStep := models.LearningStep{
		Order:         stepNumber,
		Knowledge:     *targetKnowledge,
		Objectives:    []string{fmt.Sprintf("拓展%s的综合应用", targetKnowledge.Name)},
		Resources:     GetResourcesByKnowledge(targetKnowledgeID, models.DifficultyAdvanced),
		Exercises:     GetExercisesByKnowledge(targetKnowledgeID, models.DifficultyAdvanced),
		EstimatedTime: 50,
	}
	steps = append(steps, advancedStep)
	stepNumber++

	// 3. 添加后续知识点推荐（可选）
	nextPoints := GetNextPoints(targetKnowledgeID)
	for _, nextKnowledge := range nextPoints {
		step := models.LearningStep{
			Order:         stepNumber,
			Knowledge:     *nextKnowledge,
			Objectives:    []string{fmt.Sprintf("推荐学习：%s", nextKnowledge.Name)},
			Resources:     GetResourcesByKnowledge(nextKnowledge.ID, ""),
			Exercises:     []models.Exercise{},
			EstimatedTime: pg.estimateStudyTime(nextKnowledge.Difficulty),
		}
		steps = append(steps, step)
		stepNumber++
	}

	// 计算总学习时间
	totalMinutes := 0
	for _, step := range steps {
		totalMinutes += step.EstimatedTime
	}

	path := &models.LearningPath{
		Steps: steps,
	}

	return path, nil
}

// findRequiredKnowledge 查找需要学习的前置知识（深度优先遍历）
func (pg *PathGenerator) findRequiredKnowledge(knowledgeID string, masteredMap map[string]bool) []string {
	var required []string
	visited := make(map[string]bool)

	var dfs func(string)
	dfs = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true

		// 如果已掌握，则不需要学习，但需要继续检查其前置知识
		if masteredMap[id] {
			return
		}

		knowledge := GetKnowledgePoint(id)
		if knowledge == nil {
			return
		}

		// 先递归处理前置知识
		for _, preID := range knowledge.Prerequisites {
			dfs(preID)
		}

		// 如果不是目标知识点本身，且未掌握，则需要学习
		if id != knowledgeID && !masteredMap[id] {
			required = append(required, id)
		}
	}

	knowledge := GetKnowledgePoint(knowledgeID)
	if knowledge != nil {
		for _, preID := range knowledge.Prerequisites {
			dfs(preID)
		}
	}

	return required
}

// estimateStudyTime 估算学习时间（分钟）
func (pg *PathGenerator) estimateStudyTime(difficulty models.DifficultyLevel) int {
	switch difficulty {
	case models.DifficultyBasic:
		return 20
	case models.DifficultyImprove:
		return 35
	case models.DifficultyAdvanced:
		return 50
	default:
		return 30
	}
}

// GenerateReviewPlan 生成复习计划
func (pg *PathGenerator) GenerateReviewPlan(knowledgeIDs []string) (*models.LearningPath, error) {
	var steps []models.LearningStep
	stepNumber := 1

	for _, knowledgeID := range knowledgeIDs {
		knowledge := GetKnowledgePoint(knowledgeID)
		if knowledge == nil {
			continue
		}

		// 复习步骤：主要是做题
		step := models.LearningStep{
			Order:         stepNumber,
			Knowledge:     *knowledge,
			Objectives:    []string{fmt.Sprintf("复习并巩固 %s", knowledge.Name)},
			Resources:     []models.LearningResource{},
			Exercises:     GetExercisesByKnowledge(knowledgeID, ""),
			EstimatedTime: 15,
		}
		steps = append(steps, step)
		stepNumber++
	}

	path := &models.LearningPath{
		Steps: steps,
	}

	return path, nil
}
