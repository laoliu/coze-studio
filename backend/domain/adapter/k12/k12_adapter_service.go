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

// K12AdapterService K12教育适配器服务
type K12AdapterService struct {
	analyzer      *Analyzer
	pathGenerator *PathGenerator
}

// NewK12AdapterService 创建K12适配器服务
func NewK12AdapterService() *K12AdapterService {
	return &K12AdapterService{
		analyzer:      NewAnalyzer(),
		pathGenerator: NewPathGenerator(),
	}
}

// GenerateLearningPlan 生成完整的学习计划
// userInput: 用户输入，例如 "我想学习初中数学的二次函数"
// masteredKnowledgeIDs: 已掌握的知识点ID列表
func (s *K12AdapterService) GenerateLearningPlan(userInput string, masteredKnowledgeIDs []string) (*models.LearningPlan, error) {
	// 1. 分析用户输入
	analysisResult, err := s.analyzer.AnalyzeInput(userInput)
	if err != nil {
		return nil, fmt.Errorf("分析输入失败: %w", err)
	}

	// 2. 获取目标知识点
	targetKnowledge := GetKnowledgePoint(analysisResult.KnowledgeID)
	if targetKnowledge == nil {
		return nil, fmt.Errorf("未找到知识点: %s", analysisResult.KnowledgeID)
	}

	// 3. 生成学习路径
	learningPath, err := s.pathGenerator.GenerateLearningPath(analysisResult.KnowledgeID, masteredKnowledgeIDs)
	if err != nil {
		return nil, fmt.Errorf("生成学习路径失败: %w", err)
	}

	// 4. 获取前置知识并转换为值类型
	prerequisitesPtr := GetPrerequisites(analysisResult.KnowledgeID)
	var prerequisites []models.KnowledgePoint
	for _, kp := range prerequisitesPtr {
		if kp != nil {
			prerequisites = append(prerequisites, *kp)
		}
	}

	// 5. 获取拓展内容并转换为值类型
	nextPointsPtr := GetNextPoints(analysisResult.KnowledgeID)
	var extension []models.KnowledgePoint
	for _, kp := range nextPointsPtr {
		if kp != nil {
			extension = append(extension, *kp)
		}
	}

	// 6. 计算总学习时间
	totalTime := 0
	for _, step := range learningPath.Steps {
		totalTime += step.EstimatedTime
	}

	// 7. 构建学习计划
	plan := &models.LearningPlan{
		StudentLevel:    fmt.Sprintf("%s %s", analysisResult.Grade, analysisResult.Subject),
		TargetKnowledge: *targetKnowledge,
		Prerequisites:   prerequisites,
		MainPath:        *learningPath,
		Review:          s.identifyReviewContent(masteredKnowledgeIDs),
		Extension:       extension,
		TotalTime:       totalTime,
	}

	return plan, nil
}

// AnalyzeRequest 分析请求
func (s *K12AdapterService) AnalyzeRequest(req *models.AnalysisRequest) (*models.AnalysisResult, error) {
	return s.analyzer.AnalyzeInput(req.UserInput)
}

// identifyReviewContent 识别需要复习的内容
func (s *K12AdapterService) identifyReviewContent(masteredKnowledgeIDs []string) []models.KnowledgePoint {
	var review []models.KnowledgePoint

	// 选择最近学习的几个知识点作为复习内容
	maxReview := 3
	count := 0
	for i := len(masteredKnowledgeIDs) - 1; i >= 0 && count < maxReview; i-- {
		kp := GetKnowledgePoint(masteredKnowledgeIDs[i])
		if kp != nil {
			review = append(review, *kp)
			count++
		}
	}

	return review
}

// GetKnowledgeDetails 获取知识点详情
func (s *K12AdapterService) GetKnowledgeDetails(knowledgeID string) (*models.KnowledgePoint, error) {
	kp := GetKnowledgePoint(knowledgeID)
	if kp == nil {
		return nil, fmt.Errorf("未找到知识点: %s", knowledgeID)
	}
	return kp, nil
}

// GenerateExerciseSet 生成练习题集
func (s *K12AdapterService) GenerateExerciseSet(knowledgeID string, difficulty models.DifficultyLevel, count int) ([]models.Exercise, error) {
	exercises := GetExercisesByKnowledge(knowledgeID, difficulty)
	if len(exercises) == 0 {
		return nil, fmt.Errorf("没有找到相关习题")
	}

	// 如果需要的数量超过现有数量，返回全部
	if count >= len(exercises) {
		return exercises, nil
	}

	// 否则返回前count个
	return exercises[:count], nil
}
