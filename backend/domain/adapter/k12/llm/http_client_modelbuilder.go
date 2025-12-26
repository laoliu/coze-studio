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

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
)

// ==================== ModelBuilder客户端配置 ====================

// ModelBuilderClientConfig ModelBuilder客户端配置
type ModelBuilderClientConfig struct {
	ModelID    int64         // Coze Model ID (例如: 100004)
	Timeout    time.Duration // 请求超时时间
	MaxTokens  int           // 最大token数
}

// DefaultModelBuilderClientConfig 返回默认配置
func DefaultModelBuilderClientConfig() *ModelBuilderClientConfig {
	return &ModelBuilderClientConfig{
		ModelID:   100004, // 默认模型ID
		Timeout:   30 * time.Second,
		MaxTokens: 8192,  // 增加到8192确保完整响应
	}
}

// ==================== ModelBuilderClient ====================

// ModelBuilderClient 使用Coze ModelBuilder的HTTP客户端
type ModelBuilderClient struct {
	config    *ModelBuilderClientConfig
	chatModel modelbuilder.ToolCallingChatModel
	logger    hlog.FullLogger
}

// NewModelBuilderClient 创建ModelBuilder客户端
func NewModelBuilderClient(config *ModelBuilderClientConfig, logger hlog.FullLogger) *ModelBuilderClient {
	if config == nil {
		config = DefaultModelBuilderClientConfig()
	}

	if logger == nil {
		logger = hlog.DefaultLogger()
	}

	return &ModelBuilderClient{
		config:    config,
		chatModel: nil, // 延迟初始化
		logger:    logger,
	}
}

// initChatModel 初始化ChatModel（延迟加载）
func (c *ModelBuilderClient) initChatModel(ctx context.Context) error {
	if c.chatModel != nil {
		return nil
	}

	c.logger.Infof("[ModelBuilder] Initializing ChatModel with Model ID: %d", c.config.ModelID)

	chatModel, _, err := modelbuilder.BuildModelByID(ctx, c.config.ModelID, nil)
	if err != nil {
		c.logger.Errorf("[ModelBuilder] Failed to build ChatModel: %v", err)
		return fmt.Errorf("failed to build ChatModel: %w", err)
	}

	c.chatModel = chatModel
	c.logger.Infof("[ModelBuilder] ChatModel initialized successfully")
	return nil
}

// CallLLM 调用LLM生成内容
func (c *ModelBuilderClient) CallLLM(ctx context.Context, prompt string, systemPrompt string) (string, error) {
	// 延迟初始化ChatModel
	if err := c.initChatModel(ctx); err != nil {
		return "", err
	}

	startTime := time.Now()

	// 记录请求信息
	c.logger.Infof("[LLM Request] model_id=%d, prompt_length=%d",
		c.config.ModelID, len(prompt))

	// 构建消息
	messages := []*schema.Message{}
	if systemPrompt != "" {
		messages = append(messages, &schema.Message{
			Role:    schema.System,
			Content: systemPrompt,
		})
	}
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: prompt,
	})

	// 调用ChatModel
	response, err := c.chatModel.Generate(ctx, messages, model.WithMaxTokens(c.config.MaxTokens))
	if err != nil {
		duration := time.Since(startTime)
		c.logger.Errorf("[LLM Error] error=%v, duration=%dms", err, duration.Milliseconds())
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	// 记录成功信息
	duration := time.Since(startTime)
	c.logger.Infof("[LLM Success] response_length=%d, duration=%dms",
		len(response.Content), duration.Milliseconds())

	return response.Content, nil
}

// CallLLMWithJSON 调用LLM生成JSON并自动解析
func (c *ModelBuilderClient) CallLLMWithJSON(ctx context.Context, prompt string, systemPrompt string, result interface{}) error {
	content, err := c.CallLLM(ctx, prompt, systemPrompt)
	if err != nil {
		return err
	}

	// 尝试提取JSON（有些模型会在JSON前后加说明文字）
	content = extractJSONContent(content)

	// 解析JSON
	if err := json.Unmarshal([]byte(content), result); err != nil {
		c.logger.Errorf("[LLM Parse Error] failed to parse JSON: %v, content: %s", err, content[:minInt(len(content), 200)])
		return fmt.Errorf("failed to parse LLM response as JSON: %w", err)
	}

	return nil
}

// CallLLMBatch 批量调用LLM API
func (c *ModelBuilderClient) CallLLMBatch(ctx context.Context, prompts []string, systemPrompt string) ([]string, []error) {
	results := make([]string, len(prompts))
	errors := make([]error, len(prompts))

	// 顺序调用（可以后续优化为并发）
	for i, prompt := range prompts {
		content, err := c.CallLLM(ctx, prompt, systemPrompt)
		results[i] = content
		errors[i] = err
	}

	return results, errors
}

// GetModelInfo 获取模型信息
func (c *ModelBuilderClient) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"model_id":   c.config.ModelID,
		"timeout":    c.config.Timeout.String(),
		"max_tokens": c.config.MaxTokens,
		"type":       "Coze ModelBuilder",
	}
}

// ==================== 辅助函数 ====================

// extractJSONContent 从混合文本中提取JSON
func extractJSONContent(text string) string {
	// 查找第一个{和最后一个}
	start := -1
	end := -1

	for i, ch := range text {
		if ch == '{' && start == -1 {
			start = i
		}
		if ch == '}' {
			end = i
		}
	}

	if start == -1 || end == -1 || start >= end {
		return text
	}

	return text[start : end+1]
}

// minInt 返回两个整数中的最小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
