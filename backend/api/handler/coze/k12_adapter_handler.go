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

package coze

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/k12"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// K12AdapterHandler K12教育适配器处理器
type K12AdapterHandler struct {
	service *k12.K12AdapterService
}

// NewK12AdapterHandler 创建K12适配器处理器
func NewK12AdapterHandler() *K12AdapterHandler {
	return &K12AdapterHandler{
		service: k12.NewK12AdapterService(),
	}
}

// AnalyzeInput 分析用户输入
// POST /api/adapter/k12/analyze
func (h *K12AdapterHandler) AnalyzeInput(ctx context.Context, c *app.RequestContext) {
	var req models.AnalysisRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.AnalyzeRequest(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Analysis failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GeneratePlanRequest 生成学习计划请求
type GeneratePlanRequest struct {
	UserInput            string   `json:"user_input"`
	MasteredKnowledgeIDs []string `json:"mastered_knowledge_ids"`
}

// GenerateLearningPlan 生成学习计划
// POST /api/adapter/k12/generate-plan
func (h *K12AdapterHandler) GenerateLearningPlan(ctx context.Context, c *app.RequestContext) {
	var req GeneratePlanRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 如果未提供已掌握知识点，使用空列表
	if req.MasteredKnowledgeIDs == nil {
		req.MasteredKnowledgeIDs = []string{}
	}

	plan, err := h.service.GenerateLearningPlan(req.UserInput, req.MasteredKnowledgeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate learning plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    plan,
	})
}

// GetKnowledgeDetails 获取知识点详情
// GET /api/adapter/k12/knowledge/:id
func (h *K12AdapterHandler) GetKnowledgeDetails(ctx context.Context, c *app.RequestContext) {
	knowledgeID := c.Param("id")
	if knowledgeID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Knowledge ID is required",
		})
		return
	}

	details, err := h.service.GetKnowledgeDetails(knowledgeID)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Knowledge point not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    details,
	})
}

// GenerateExercisesRequest 生成习题集请求
type GenerateExercisesRequest struct {
	KnowledgeID string                 `json:"knowledge_id"`
	Difficulty  models.DifficultyLevel `json:"difficulty"`
	Count       int                    `json:"count"`
}

// GenerateExercises 生成习题集
// POST /api/adapter/k12/exercises
func (h *K12AdapterHandler) GenerateExercises(ctx context.Context, c *app.RequestContext) {
	var req GenerateExercisesRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 默认值
	if req.Count <= 0 {
		req.Count = 10
	}

	exercises, err := h.service.GenerateExerciseSet(req.KnowledgeID, req.Difficulty, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate exercises",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    exercises,
		"count":   len(exercises),
	})
}
