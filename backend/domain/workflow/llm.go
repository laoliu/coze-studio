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

package workflow

import (
	"context"

	"github.com/cloudwego/eino/components/model"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config/modelmgr"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/llm"
)

// GetLLMClient 获取 LLM 客户端（公开函数）
func GetLLMClient(ctx context.Context) (model.BaseChatModel, *modelmgr.Model, error) {
	return llm.GetLLMClient(ctx)
}

// GetLLMConfigurationStatus 获取 LLM 配置状态（公开函数）
func GetLLMConfigurationStatus(ctx context.Context) map[string]interface{} {
	return llm.GetConfigurationStatus(ctx)
}
