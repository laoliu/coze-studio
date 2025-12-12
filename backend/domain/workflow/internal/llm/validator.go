package llm

import (
	"context"
	"fmt"
	
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ValidateConfiguration 验证 LLM 配置是否正确
func ValidateConfiguration(ctx context.Context) error {
	logs.CtxInfof(ctx, "[LLM Validator] Starting LLM configuration validation...")
	
	// 1. 加载配置
	cfg, err := LoadConfig("backend/conf/workflow/llm_config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load LLM config: %w", err)
	}
	
	// 2. 检查是否启用
	if !cfg.WorkflowGenerator.Enabled {
		logs.CtxWarnf(ctx, "[LLM Validator] Workflow generator is disabled in config")
		return fmt.Errorf("workflow generator is disabled, please set 'workflow_generator.enabled: true' in llm_config.yaml")
	}
	
	logs.CtxInfof(ctx, "[LLM Validator] Workflow generator enabled: model_id=%d", cfg.WorkflowGenerator.ModelID)
	
	// 3. 尝试获取 LLM 客户端
	client, modelInfo, err := GetLLMClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize LLM client: %w", err)
	}
	
	if client == nil {
		return fmt.Errorf("LLM client is nil")
	}
	
	logs.CtxInfof(ctx, "[LLM Validator] Successfully initialized LLM client:")
	logs.CtxInfof(ctx, "[LLM Validator]   - Model ID: %d", modelInfo.ID)
	logs.CtxInfof(ctx, "[LLM Validator]   - Model Name: %s", modelInfo.DisplayInfo.Name)
	logs.CtxInfof(ctx, "[LLM Validator]   - Model Type: %s", modelInfo.Type)
	logs.CtxInfof(ctx, "[LLM Validator]   - Provider: %s", modelInfo.Provider)
	
	// 4. 验证生成配置
	genCfg := cfg.WorkflowGenerator.Generation
	if genCfg.Temperature < 0 || genCfg.Temperature > 2 {
		logs.CtxWarnf(ctx, "[LLM Validator] Temperature %.2f is out of recommended range [0, 2]", genCfg.Temperature)
	}
	
	if genCfg.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be positive, got %d", genCfg.MaxTokens)
	}
	
	if genCfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %d", genCfg.Timeout)
	}
	
	logs.CtxInfof(ctx, "[LLM Validator] Generation config:")
	logs.CtxInfof(ctx, "[LLM Validator]   - Temperature: %.2f", genCfg.Temperature)
	logs.CtxInfof(ctx, "[LLM Validator]   - Max Tokens: %d", genCfg.MaxTokens)
	logs.CtxInfof(ctx, "[LLM Validator]   - Timeout: %ds", genCfg.Timeout)
	
	logs.CtxInfof(ctx, "[LLM Validator] ✅ All validations passed!")
	return nil
}

// GetConfigurationStatus 获取配置状态（用于健康检查）
func GetConfigurationStatus(ctx context.Context) map[string]interface{} {
	status := map[string]interface{}{
		"enabled": false,
		"healthy": false,
	}
	
	cfg, err := LoadConfig("backend/conf/workflow/llm_config.yaml")
	if err != nil {
		status["error"] = err.Error()
		return status
	}
	
	status["enabled"] = cfg.WorkflowGenerator.Enabled
	
	if !cfg.WorkflowGenerator.Enabled {
		status["message"] = "Workflow generator is disabled"
		return status
	}
	
	status["model_id"] = cfg.WorkflowGenerator.ModelID
	
	_, modelInfo, err := GetLLMClient(ctx)
	if err != nil {
		status["error"] = err.Error()
		return status
	}
	
	status["healthy"] = true
	status["model_name"] = modelInfo.DisplayInfo.Name
	status["model_type"] = modelInfo.Type
	status["provider"] = modelInfo.Provider
	status["temperature"] = cfg.WorkflowGenerator.Generation.Temperature
	status["max_tokens"] = cfg.WorkflowGenerator.Generation.MaxTokens
	status["timeout"] = cfg.WorkflowGenerator.Generation.Timeout
	
	return status
}
