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
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/eino/components/model"
	"gopkg.in/yaml.v3"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config/modelmgr"
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// WorkflowGeneratorConfig 工作流生成器配置
type WorkflowGeneratorConfig struct {
	Enabled    bool              `yaml:"enabled"`
	ModelID    int64             `yaml:"model_id"`
	Generation GenerationConfig  `yaml:"generation"`
	Cache      CacheConfig       `yaml:"cache"`
	Retry      RetryConfig       `yaml:"retry"`
	Quality    QualityConfig     `yaml:"quality"`
	Layout     LayoutConfig      `yaml:"layout"`
}

// GenerationConfig 生成配置
type GenerationConfig struct {
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	Timeout     int     `yaml:"timeout"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	EnableIntentCache bool `yaml:"enable_intent_cache"`
	CacheTTL          int  `yaml:"cache_ttl"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int `yaml:"max_retries"`
	RetryInterval int `yaml:"retry_interval"`
}

// QualityConfig 质量控制配置
type QualityConfig struct {
	ConfidenceThreshold float64 `yaml:"confidence_threshold"`
	EnableValidation    bool    `yaml:"enable_validation"`
}

// LayoutConfig 布局配置
type LayoutConfig struct {
	EnableAutoLayout bool `yaml:"enable_auto_layout"`
	NodeSpacingX     int  `yaml:"node_spacing_x"`
	NodeSpacingY     int  `yaml:"node_spacing_y"`
	StartX           int  `yaml:"start_x"`
	StartY           int  `yaml:"start_y"`
}

// TemplateMatchingConfig 模板匹配配置
type TemplateMatchingConfig struct {
	MaxResults          int     `yaml:"max_results"`
	SimilarityThreshold float64 `yaml:"similarity_threshold"`
	EnableRecommendation bool   `yaml:"enable_recommendation"`
}

// Config 完整配置
type Config struct {
	WorkflowGenerator WorkflowGeneratorConfig `yaml:"workflow_generator"`
	TemplateMatching  TemplateMatchingConfig  `yaml:"template_matching"`
}

var (
	config     *Config
	configOnce sync.Once
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	var err error
	configOnce.Do(func() {
		data, readErr := os.ReadFile(configPath)
		if readErr != nil {
			err = fmt.Errorf("failed to read config file: %w", readErr)
			return
		}

		// 替换环境变量
		expandedData := os.ExpandEnv(string(data))

		config = &Config{}
		if unmarshalErr := yaml.Unmarshal([]byte(expandedData), config); unmarshalErr != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", unmarshalErr)
			return
		}

		logs.Infof("Workflow generator config loaded successfully, enabled: %v, model_id: %d",
			config.WorkflowGenerator.Enabled, config.WorkflowGenerator.ModelID)
	})

	return config, err
}

// GetConfig 获取配置
func GetConfig() *Config {
	if config == nil {
		// 尝试加载默认配置
		defaultPath := "resources/conf/workflow/llm_config.yaml"
		cfg, err := LoadConfig(defaultPath)
		if err != nil {
			logs.Errorf("Failed to load config from %s: %v", defaultPath, err)
			// 返回默认配置
			return getDefaultConfig()
		}
		return cfg
	}
	return config
}

// getDefaultConfig 返回默认配置
func getDefaultConfig() *Config {
	return &Config{
		WorkflowGenerator: WorkflowGeneratorConfig{
			Enabled: true,
			ModelID: 0, // 0 表示使用 builtin 模型
			Generation: GenerationConfig{
				Temperature: 0.7,
				MaxTokens:   2000,
				Timeout:     45,
			},
			Cache: CacheConfig{
				EnableIntentCache: true,
				CacheTTL:          60,
			},
			Retry: RetryConfig{
				MaxRetries:    3,
				RetryInterval: 2,
			},
			Quality: QualityConfig{
				ConfidenceThreshold: 0.6,
				EnableValidation:    true,
			},
			Layout: LayoutConfig{
				EnableAutoLayout: true,
				NodeSpacingX:     200,
				NodeSpacingY:     100,
				StartX:           100,
				StartY:           100,
			},
		},
		TemplateMatching: TemplateMatchingConfig{
			MaxResults:          3,
			SimilarityThreshold: 0.5,
			EnableRecommendation: true,
		},
	}
}

// GetLLMClient 获取 LLM 客户端（使用 Coze Studio 现有的模型管理系统）
func GetLLMClient(ctx context.Context) (model.BaseChatModel, *modelmgr.Model, error) {
	cfg := GetConfig()

	if !cfg.WorkflowGenerator.Enabled {
		return nil, nil, fmt.Errorf("workflow generator is disabled in configuration")
	}

	var (
		chatModel modelbuilder.BaseChatModel
		modelInfo *modelmgr.Model
		err       error
	)

	// 构建 LLM 参数 - 转换 float64 到 float32
	temperature := float32(cfg.WorkflowGenerator.Generation.Temperature)
	params := &modelbuilder.LLMParams{
		Temperature: &temperature,
		MaxTokens:   cfg.WorkflowGenerator.Generation.MaxTokens,
	}

	// 如果指定了 model_id，使用指定的模型
	if cfg.WorkflowGenerator.ModelID > 0 {
		logs.CtxInfof(ctx, "Using specified model ID: %d", cfg.WorkflowGenerator.ModelID)
		chatModel, modelInfo, err = modelbuilder.BuildModelByID(ctx, cfg.WorkflowGenerator.ModelID, params)
	} else {
		// 否则使用 builtin 模型
		logs.CtxInfof(ctx, "Using builtin model")
		chatModel, configured, err := modelbuilder.GetBuiltinChatModel(ctx, "WORKFLOW_GENERATOR")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get builtin chat model: %w", err)
		}
		if !configured || chatModel == nil {
			return nil, nil, fmt.Errorf("no builtin model configured or chat model is nil")
		}

		// 获取模型信息（builtin 模型可能没有详细的 modelInfo）
		// 这里返回 nil 也是可以的，generator 可以处理
		return chatModel, nil, nil
	}

	if err != nil {
		logs.CtxErrorf(ctx, "Failed to get LLM client: %v", err)
		return nil, nil, fmt.Errorf("failed to build LLM model: %w", err)
	}

	logs.CtxInfof(ctx, "LLM client initialized successfully")
	return chatModel, modelInfo, nil
}

// IsEnabled 检查功能是否启用
func IsEnabled() bool {
	cfg := GetConfig()
	return cfg.WorkflowGenerator.Enabled
}

// GetTimeout 获取超时时间
func GetTimeout() time.Duration {
	cfg := GetConfig()
	return time.Duration(cfg.WorkflowGenerator.Generation.Timeout) * time.Second
}

// GetMaxRetries 获取最大重试次数
func GetMaxRetries() int {
	cfg := GetConfig()
	return cfg.WorkflowGenerator.Retry.MaxRetries
}

// GetConfidenceThreshold 获取置信度阈值
func GetConfidenceThreshold() float64 {
	cfg := GetConfig()
	return cfg.WorkflowGenerator.Quality.ConfidenceThreshold
}

