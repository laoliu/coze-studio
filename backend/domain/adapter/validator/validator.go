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

package validator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
)

var (
	ErrInvalidAdapterID  = errors.New("invalid adapter ID")
	ErrInvalidVersion    = errors.New("invalid version")
	ErrMissingMetadata   = errors.New("missing required metadata")
	ErrInvalidSchema     = errors.New("invalid schema")
	ErrSecurityViolation = errors.New("security violation detected")
)

// AdapterValidator 适配器验证器
type AdapterValidator struct {
	// 可以注入安全策略、Schema 验证器等
}

// NewAdapterValidator 创建验证器
func NewAdapterValidator() *AdapterValidator {
	return &AdapterValidator{}
}

// ValidateAdapter 验证适配器
func (v *AdapterValidator) ValidateAdapter(ctx context.Context, adp adapter.SimplifiedDomainAdapter) error {
	// 1. 获取元数据
	metadata := adp.GetMetadata()
	if metadata == nil {
		return ErrMissingMetadata
	}

	// 2. 验证 ID 格式
	if err := v.validateAdapterID(metadata.AdapterID); err != nil {
		return err
	}

	// 3. 验证版本格式
	if err := v.validateVersion(metadata.Version); err != nil {
		return err
	}

	// 4. 验证必填字段
	if metadata.Name == "" {
		return fmt.Errorf("%w: name is required", ErrMissingMetadata)
	}
	if metadata.DisplayName == "" {
		return fmt.Errorf("%w: display_name is required", ErrMissingMetadata)
	}
	if metadata.Description == "" {
		return fmt.Errorf("%w: description is required", ErrMissingMetadata)
	}

	// 5. 验证 Schema
	if err := v.validateSchemas(metadata); err != nil {
		return err
	}

	// 6. 安全检查
	if err := v.securityCheck(metadata); err != nil {
		return err
	}

	return nil
}

// validateAdapterID 验证适配器 ID 格式
// 格式: domain.adapter_name (例如: k12.math_problem_generator)
func (v *AdapterValidator) validateAdapterID(adapterID string) error {
	if adapterID == "" {
		return ErrInvalidAdapterID
	}

	// 检查格式：字母、数字、点、下划线、连字符
	pattern := `^[a-z0-9][a-z0-9._-]*[a-z0-9]$`
	matched, err := regexp.MatchString(pattern, adapterID)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("%w: must match pattern %s", ErrInvalidAdapterID, pattern)
	}

	return nil
}

// validateVersion 验证版本格式
// 支持语义化版本: 1.0.0, 1.2.3-beta.1
func (v *AdapterValidator) validateVersion(version string) error {
	if version == "" {
		return ErrInvalidVersion
	}

	// 语义化版本正则
	pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	matched, err := regexp.MatchString(pattern, version)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("%w: must be a valid semantic version (e.g., 1.0.0)", ErrInvalidVersion)
	}

	return nil
}

// validateSchemas 验证输入输出 Schema
func (v *AdapterValidator) validateSchemas(metadata *adapter.AdapterMetadata) error {
	// 输入 Schema 必须存在
	if metadata.InputSchema == nil || len(metadata.InputSchema) == 0 {
		return fmt.Errorf("%w: input schema is required", ErrInvalidSchema)
	}

	// 输出 Schema 必须存在
	if metadata.OutputSchema == nil || len(metadata.OutputSchema) == 0 {
		return fmt.Errorf("%w: output schema is required", ErrInvalidSchema)
	}

	// TODO: 可以添加更详细的 JSON Schema 验证
	// 使用 github.com/xeipuuv/gojsonschema 等库

	return nil
}

// securityCheck 安全检查
func (v *AdapterValidator) securityCheck(metadata *adapter.AdapterMetadata) error {
	// 1. 检查是否包含敏感关键词
	sensitiveKeywords := []string{
		"password", "secret", "token", "key", "credential",
		"private", "confidential",
	}

	description := metadata.Description
	for _, keyword := range sensitiveKeywords {
		if containsCaseInsensitive(description, keyword) {
			// 警告：不阻止，但记录日志
			// TODO: 添加日志记录
		}
	}

	// 2. 检查域名限制（如果配置了白名单）
	// TODO: 实现域名白名单检查

	// 3. 检查是否尝试访问系统资源
	// TODO: 实现资源访问检查

	return nil
}

// ValidateInput 验证适配器输入
func (v *AdapterValidator) ValidateInput(ctx context.Context, input *adapter.AdapterInput) error {
	if input == nil {
		return errors.New("input is nil")
	}

	// 1. 验证必填字段
	if input.Domain == "" {
		return errors.New("domain is required")
	}
	if input.ActivityType == "" {
		return errors.New("activity_type is required")
	}
	if input.Parameters == nil {
		return errors.New("parameters is required")
	}

	// 2. 验证参数格式
	// TODO: 根据 Schema 验证参数

	return nil
}

// containsCaseInsensitive 不区分大小写的字符串包含检查
func containsCaseInsensitive(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
