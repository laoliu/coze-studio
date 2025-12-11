// Copyright 2025 Coze Studio. All rights reserved.

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AdapterConfig 适配器配置
type AdapterConfig struct {
	LLM      LLMConfig      `yaml:"llm"`
	Database DatabaseConfig `yaml:"database"`
	Content  ContentConfig  `yaml:"content"`
	Debug    DebugConfig    `yaml:"debug"`
}

// LLMConfig LLM 配置
type LLMConfig struct {
	Provider    string  `yaml:"provider"`    // openai, ark, claude, gemini
	Model       string  `yaml:"model"`
	APIKey      string  `yaml:"api_key"`
	BaseURL     string  `yaml:"base_url"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	MaxConnections     int    `yaml:"max_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}

// ContentConfig 内容配置
type ContentConfig struct {
	SearchLimit         int     `yaml:"search_limit"`
	RelevanceThreshold  float64 `yaml:"relevance_threshold"`
}

// DebugConfig 调试配置
type DebugConfig struct {
	EnableMockLLM   bool `yaml:"enable_mock_llm"`
	EnableMockDB    bool `yaml:"enable_mock_db"`
	LogPrompts      bool `yaml:"log_prompts"`
	LogResponses    bool `yaml:"log_responses"`
}

// LoadConfig 加载配置
func LoadConfig(path string) (*AdapterConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AdapterConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 从环境变量覆盖敏感信息
	if apiKey := os.Getenv("LLM_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}

	return &config, nil
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *AdapterConfig {
	config := &AdapterConfig{
		LLM: LLMConfig{
			Provider:    getEnv("LLM_PROVIDER", "openai"),
			Model:       getEnv("LLM_MODEL", "gpt-4"),
			APIKey:      getEnv("LLM_API_KEY", ""),
			BaseURL:     getEnv("LLM_BASE_URL", ""),
			Temperature: getEnvFloat("LLM_TEMPERATURE", 0.7),
			MaxTokens:   getEnvInt("LLM_MAX_TOKENS", 2000),
		},
		Database: DatabaseConfig{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnvInt("DB_PORT", 3306),
			User:               getEnv("DB_USER", "coze"),
			Password:           getEnv("DB_PASSWORD", "coze123"),
			Database:           getEnv("DB_NAME", "opencoze"),
			MaxConnections:     getEnvInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 10),
		},
		Content: ContentConfig{
			SearchLimit:        getEnvInt("CONTENT_SEARCH_LIMIT", 10),
			RelevanceThreshold: getEnvFloat("CONTENT_RELEVANCE_THRESHOLD", 0.7),
		},
		Debug: DebugConfig{
			EnableMockLLM:  getEnvBool("DEBUG_ENABLE_MOCK_LLM", false),
			EnableMockDB:   getEnvBool("DEBUG_ENABLE_MOCK_DB", false),
			LogPrompts:     getEnvBool("DEBUG_LOG_PROMPTS", false),
			LogResponses:   getEnvBool("DEBUG_LOG_RESPONSES", false),
		},
	}

	return config
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		var floatValue float64
		fmt.Sscanf(value, "%f", &floatValue)
		return floatValue
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
