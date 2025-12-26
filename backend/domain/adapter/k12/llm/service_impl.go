/*
 * Copyright 2025 coze-dev Authors
 *
 * Licens	config := &ModelBuilderClientConfig{
		ModelID:   modelID,
		Timeout:   timeout,
		MaxTokens: 8192, // 增加到8192以确保完整响应（简化提示词后应该足够）
	}der the Apache License, Version 2.0 (the "License");
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

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// CozeLLMService Coze LLM服务实现
type CozeLLMService struct {
	modelBuilderClient *ModelBuilderClient  // 使用ModelBuilder客户端
	logger             hlog.FullLogger
}

// NewCozeLLMService 创建Coze LLM服务
func NewCozeLLMService(apiKey, apiURL, modelName string, timeout time.Duration, logger hlog.FullLogger) *CozeLLMService {
	// 使用ModelBuilder客户端，apiKey/apiURL/modelName参数已废弃
	// 现在从环境变量COZE_MODEL_ID读取Model ID
	modelID := int64(100004) // 默认使用100004
	if modelIDStr := os.Getenv("COZE_MODEL_ID"); modelIDStr != "" {
		if id, err := strconv.ParseInt(modelIDStr, 10, 64); err == nil {
			modelID = id
		}
	}

	config := &ModelBuilderClientConfig{
		ModelID:   modelID,
		Timeout:   timeout,
		MaxTokens: 8192,  // 增加到8192以确保能生成完整的15-20个知识点
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if logger == nil {
		logger = hlog.DefaultLogger()
	}

	modelBuilderClient := NewModelBuilderClient(config, logger)

	return &CozeLLMService{
		modelBuilderClient: modelBuilderClient,
		logger:             logger,
	}
}

// GenerateKnowledgeGraph 生成完整知识图谱
func (s *CozeLLMService) GenerateKnowledgeGraph(ctx context.Context, req *KnowledgeGraphRequest) (*KnowledgeGraph, error) {
	startTime := time.Now()

	// 构建提示词
	prompt, err := s.buildPrompt(KnowledgeGraphPromptTemplate, map[string]interface{}{
		"Subject":           req.Subject,
		"Grade":             req.Grade,
		"AdditionalContext": req.AdditionalContext,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// 调用LLM
	response, err := s.callLLM(ctx, prompt)
	if err != nil {
		s.logGeneration("knowledge_graph", string(req.Subject), string(req.Grade), "", prompt, "", false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	// 打印原始响应（用于调试）- 保存到临时文件以查看完整内容
	tmpFile := "/tmp/llm_response_debug.txt"
	os.WriteFile(tmpFile, []byte(response), 0644)
	s.logger.Infof("[LLM Raw Response] length=%d, saved to %s", len(response), tmpFile)
	s.logger.Infof("[LLM Raw Response First 100 chars] %s", response[:minInt(len(response), 100)])

	// 清理markdown代码块（如果存在）
	cleanedResponse := cleanMarkdownCodeBlock(response)
	s.logger.Infof("[LLM Cleaned Response] original_length=%d, cleaned_length=%d", len(response), len(cleanedResponse))
	s.logger.Infof("[LLM Cleaned Response First 100 chars] %s", cleanedResponse[:minInt(len(cleanedResponse), 100)])

	// 解析响应
	var result struct {
		Subject         string                     `json:"subject"`
		Grade           string                     `json:"grade"`
		KnowledgePoints []*models.KnowledgePoint `json:"knowledge_points"`
	}

	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		s.logGeneration("knowledge_graph", string(req.Subject), string(req.Grade), "", prompt, response, false, err.Error(), time.Since(startTime))
		s.logger.Errorf("[JSON Parse Error] cleaned response: %s", cleanedResponse)
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// 设置knowledge points的学科和年级
	for _, kp := range result.KnowledgePoints {
		kp.Subject = req.Subject
		kp.Grade = req.Grade
	}

	kg := &KnowledgeGraph{
		Subject:         req.Subject,
		Grade:           req.Grade,
		KnowledgePoints: result.KnowledgePoints,
		GeneratedAt:     time.Now(),
	}

	s.logGeneration("knowledge_graph", string(req.Subject), string(req.Grade), "", prompt, response, true, "", time.Since(startTime))
	return kg, nil
}

// GenerateKnowledgePoint 生成单个知识点
func (s *CozeLLMService) GenerateKnowledgePoint(ctx context.Context, req *KnowledgePointRequest) (*models.KnowledgePoint, error) {
	startTime := time.Now()

	prompt, err := s.buildPrompt(KnowledgePointPromptTemplate, map[string]interface{}{
		"KnowledgeID":   req.KnowledgeID,
		"KnowledgeName": req.Name,
		"Subject":       req.Subject,
		"Grade":         req.Grade,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	response, err := s.callLLM(ctx, prompt)
	if err != nil {
		s.logGeneration("knowledge_point", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, "", false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	var kp models.KnowledgePoint
	if err := json.Unmarshal([]byte(response), &kp); err != nil {
		s.logGeneration("knowledge_point", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	s.logGeneration("knowledge_point", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, true, "", time.Since(startTime))
	return &kp, nil
}

// GenerateResources 生成学习资源
func (s *CozeLLMService) GenerateResources(ctx context.Context, req *ResourceRequest) ([]*models.Resource, error) {
	startTime := time.Now()

	prompt, err := s.buildPrompt(ResourceGenerationPromptTemplate, map[string]interface{}{
		"KnowledgeID":   req.KnowledgeID,
		"KnowledgeName": req.KnowledgeName,
		"Subject":       req.Subject,
		"Grade":         req.Grade,
		"Description":   req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	response, err := s.callLLM(ctx, prompt)
	if err != nil {
		s.logGeneration("resource", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, "", false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	var result struct {
		KnowledgeID string             `json:"knowledge_id"`
		Resources   []*models.Resource `json:"resources"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		s.logGeneration("resource", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	s.logGeneration("resource", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, true, "", time.Since(startTime))
	return result.Resources, nil
}

// GenerateExercises 生成习题
func (s *CozeLLMService) GenerateExercises(ctx context.Context, req *ExerciseRequest) ([]*models.Exercise, error) {
	startTime := time.Now()

	prompt, err := s.buildPrompt(ExerciseGenerationPromptTemplate, map[string]interface{}{
		"KnowledgeID":   req.KnowledgeID,
		"KnowledgeName": req.KnowledgeName,
		"Subject":       req.Subject,
		"Grade":         req.Grade,
		"Difficulty":    req.Difficulty,
		"Count":         req.Count,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	response, err := s.callLLM(ctx, prompt)
	if err != nil {
		s.logGeneration("exercise", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, "", false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	var result struct {
		KnowledgeID string             `json:"knowledge_id"`
		Difficulty  string             `json:"difficulty"`
		Exercises   []*models.Exercise `json:"exercises"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		s.logGeneration("exercise", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, false, err.Error(), time.Since(startTime))
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	s.logGeneration("exercise", string(req.Subject), string(req.Grade), req.KnowledgeID, prompt, response, true, "", time.Since(startTime))
	return result.Exercises, nil
}

// buildPrompt 构建提示词
func (s *CozeLLMService) buildPrompt(templateStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// callLLM 调用大模型API
func (s *CozeLLMService) callLLM(ctx context.Context, prompt string) (string, error) {
	// 使用ModelBuilder客户端调用真实的Coze API
	systemPrompt := "你是一个专业的K12教育内容生成助手。请根据用户的要求，生成结构化的教育内容，并以JSON格式返回。"

	content, err := s.modelBuilderClient.CallLLM(ctx, prompt, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to call LLM: %w", err)
	}

	return content, nil
}

// logGeneration 记录生成日志
func (s *CozeLLMService) logGeneration(requestType, subject, grade, knowledgeID, prompt, response string, success bool, errorMsg string, duration time.Duration) {
	log := GenerationLog{
		RequestType:    requestType,
		Subject:        models.Subject(subject),
		Grade:          models.GradeLevel(grade),
		KnowledgeID:    knowledgeID,
		PromptTemplate: prompt,
		LLMResponse:    response,
		Success:        success,
		ErrorMessage:   errorMsg,
		DurationMS:     duration.Milliseconds(),
		CreatedAt:      time.Now(),
	}

	// TODO: 保存到数据库
	if success {
		s.logger.Infof("[LLM Generation Success] type=%s, subject=%s, grade=%s, knowledge_id=%s, duration=%dms",
			requestType, subject, grade, knowledgeID, duration.Milliseconds())
	} else {
		s.logger.Errorf("[LLM Generation Failed] type=%s, subject=%s, grade=%s, knowledge_id=%s, error=%s, duration=%dms",
			requestType, subject, grade, knowledgeID, errorMsg, duration.Milliseconds())
	}

	_ = log // 避免未使用变量警告
}

// cleanMarkdownCodeBlock 清理 markdown 代码块标记
// LLM 可能返回 ```json {...} ``` 格式，需要提取出纯 JSON
func cleanMarkdownCodeBlock(s string) string {
	// 移除开头的 ```json 或 ```
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimSpace(s)
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSpace(s)
	}

	// 移除结尾的 ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}

	return s
}
