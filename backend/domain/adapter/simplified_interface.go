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

package adapter

import (
	"context"
)

// AdapterMetadata 适配器元数据
type AdapterMetadata struct {
	AdapterID              string                 `json:"adapter_id"`
	Name                   string                 `json:"name"`
	DisplayName            string                 `json:"display_name"`
	Version                string                 `json:"version"`
	Description            string                 `json:"description"`
	SupportedDomains       []string               `json:"supported_domains"`
	SupportedActivityTypes []string               `json:"supported_activity_types"`
	InputSchema            map[string]interface{} `json:"input_schema"`
	OutputSchema           map[string]interface{} `json:"output_schema"`
}

// AdapterInput 适配器输入
type AdapterInput struct {
	Domain       string                 `json:"domain"`
	ActivityType string                 `json:"activity_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	Context      map[string]interface{} `json:"context"`
}

// AdapterOutput 适配器输出
type AdapterOutput struct {
	Result   map[string]interface{} `json:"result"`
	Metadata map[string]interface{} `json:"metadata"`
}

// SimplifiedDomainAdapter 简化的领域适配器接口
// 这是新设计的简化接口，用于适配器市场
type SimplifiedDomainAdapter interface {
	// GetMetadata 获取适配器元数据
	GetMetadata() *AdapterMetadata

	// Validate 验证输入
	Validate(ctx context.Context, input *AdapterInput) error

	// Execute 执行适配器逻辑
	Execute(ctx context.Context, input *AdapterInput) (*AdapterOutput, error)
}
