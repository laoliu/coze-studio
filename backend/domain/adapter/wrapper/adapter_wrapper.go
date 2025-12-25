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

package wrapper

import (
	"context"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/component"
)

// DomainAdapterWrapper 将 SimplifiedDomainAdapter 包装为 ComponentExecutor
// 这是适配器模式的核心实现
type DomainAdapterWrapper struct {
	pluginID int64
	metadata *adapter.AdapterMetadata
	adapter  adapter.SimplifiedDomainAdapter
}

// NewDomainAdapterWrapper 创建适配器包装器
func NewDomainAdapterWrapper(
	pluginID int64,
	metadata *adapter.AdapterMetadata,
	adp adapter.SimplifiedDomainAdapter,
) *DomainAdapterWrapper {
	return &DomainAdapterWrapper{
		pluginID: pluginID,
		metadata: metadata,
		adapter:  adp,
	}
}

// GetMetadata 获取元数据
func (w *DomainAdapterWrapper) GetMetadata() *adapter.AdapterMetadata {
	return w.metadata
}

// GetAdapter 获取适配器实现
func (w *DomainAdapterWrapper) GetAdapter() adapter.SimplifiedDomainAdapter {
	return w.adapter
}

// AdapterToComponentExecutor 将 DomainAdapterWrapper 转换为 ComponentExecutor
func AdapterToComponentExecutor(wrapper *DomainAdapterWrapper) component.ComponentExecutor {
	return &componentExecutorAdapter{
		wrapper: wrapper,
	}
}

// componentExecutorAdapter 实现 ComponentExecutor 接口
type componentExecutorAdapter struct {
	wrapper *DomainAdapterWrapper
}

// GetComponentInfo 获取组件信息
func (a *componentExecutorAdapter) GetComponentInfo() *component.ComponentInfo {
	return &component.ComponentInfo{
		ComponentID:  a.wrapper.metadata.AdapterID,
		Name:         a.wrapper.metadata.Name,
		DisplayName:  a.wrapper.metadata.DisplayName,
		Version:      a.wrapper.metadata.Version,
		Description:  a.wrapper.metadata.Description,
		Type:         "adapter",
		Category:     "domain_adapter", // 适配器统一使用 domain_adapter 分类
		InputSchema:  a.wrapper.metadata.InputSchema,
		OutputSchema: a.wrapper.metadata.OutputSchema,
	}
}

// Execute 执行组件
func (a *componentExecutorAdapter) Execute(ctx context.Context, input *component.ComponentInput) (*component.ComponentOutput, error) {
	startTime := time.Now()

	// 1. 转换输入
	adapterInput := toAdapterInput(input)

	// 2. 执行适配器
	adapterOutput, err := a.wrapper.adapter.Execute(ctx, adapterInput)
	if err != nil {
		return &component.ComponentOutput{
			ComponentID:   a.wrapper.metadata.AdapterID,
			Success:       false,
			Error:         err.Error(),
			ExecutionTime: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 3. 转换输出
	return &component.ComponentOutput{
		ComponentID:   a.wrapper.metadata.AdapterID,
		Result:        adapterOutput.Result,
		Metadata:      adapterOutput.Metadata,
		ExecutionTime: time.Since(startTime).Milliseconds(),
		Success:       true,
	}, nil
}

// Validate 验证输入
func (a *componentExecutorAdapter) Validate(input *component.ComponentInput) error {
	// 转换为适配器输入并验证
	adapterInput := toAdapterInput(input)
	return a.wrapper.adapter.Validate(context.Background(), adapterInput)
}

// GetConfig 获取配置
func (a *componentExecutorAdapter) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"adapter_id": a.wrapper.metadata.AdapterID,
		"version":    a.wrapper.metadata.Version,
		"plugin_id":  a.wrapper.pluginID,
	}
}

// toAdapterInput 转换为适配器输入
func toAdapterInput(compInput *component.ComponentInput) *adapter.AdapterInput {
	return &adapter.AdapterInput{
		Domain:       getStringFromContext(compInput.Context, "domain"),
		ActivityType: getStringFromContext(compInput.Context, "activity_type"),
		Parameters:   compInput.Parameters,
		Context:      compInput.Context,
	}
}

// getStringFromContext 从 context 获取字符串值
func getStringFromContext(ctx map[string]interface{}, key string) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
