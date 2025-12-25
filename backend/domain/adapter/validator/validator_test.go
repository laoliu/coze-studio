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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
)

// mockAdapter 模拟适配器
type mockAdapter struct {
	metadata *adapter.AdapterMetadata
}

func (m *mockAdapter) GetMetadata() *adapter.AdapterMetadata {
	return m.metadata
}

func (m *mockAdapter) Validate(ctx context.Context, input *adapter.AdapterInput) error {
	return nil
}

func (m *mockAdapter) Execute(ctx context.Context, input *adapter.AdapterInput) (*adapter.AdapterOutput, error) {
	return &adapter.AdapterOutput{
		Result:   make(map[string]interface{}),
		Metadata: make(map[string]interface{}),
	}, nil
}

func TestAdapterValidator_ValidateAdapter(t *testing.T) {
	validator := NewAdapterValidator()
	ctx := context.Background()

	tests := []struct {
		name      string
		metadata  *adapter.AdapterMetadata
		wantError bool
		errorType error
	}{
		{
			name: "valid adapter",
			metadata: &adapter.AdapterMetadata{
				AdapterID:   "k12.math_problem_generator",
				Name:        "math_problem_generator",
				DisplayName: "K12 Math Problem Generator",
				Description: "Generate math problems for K12 education",
				Version:     "1.0.0",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"topic": map[string]interface{}{"type": "string"},
					},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"problem": map[string]interface{}{"type": "string"},
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid adapter ID - uppercase",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    "K12.MathProblem",
				Name:         "math_problem_generator",
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "1.0.0",
				InputSchema:  map[string]interface{}{"type": "object"},
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: true,
			errorType: ErrInvalidAdapterID,
		},
		{
			name: "invalid adapter ID - starts with dot",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    ".math.problem",
				Name:         "math_problem_generator",
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "1.0.0",
				InputSchema:  map[string]interface{}{"type": "object"},
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: true,
			errorType: ErrInvalidAdapterID,
		},
		{
			name: "invalid version",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    "k12.math_problem_generator",
				Name:         "math_problem_generator",
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "v1.0", // Invalid semantic version
				InputSchema:  map[string]interface{}{"type": "object"},
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: true,
			errorType: ErrInvalidVersion,
		},
		{
			name: "missing name",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    "k12.math_problem_generator",
				Name:         "", // Missing
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "1.0.0",
				InputSchema:  map[string]interface{}{"type": "object"},
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: true,
			errorType: ErrMissingMetadata,
		},
		{
			name: "missing input schema",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    "k12.math_problem_generator",
				Name:         "math_problem_generator",
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "1.0.0",
				InputSchema:  nil, // Missing
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: true,
			errorType: ErrInvalidSchema,
		},
		{
			name: "valid semantic version with prerelease",
			metadata: &adapter.AdapterMetadata{
				AdapterID:    "k12.math_problem_generator",
				Name:         "math_problem_generator",
				DisplayName:  "K12 Math Problem Generator",
				Description:  "Generate math problems",
				Version:      "1.0.0-beta.1",
				InputSchema:  map[string]interface{}{"type": "object"},
				OutputSchema: map[string]interface{}{"type": "object"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdp := &mockAdapter{metadata: tt.metadata}
			err := validator.ValidateAdapter(ctx, mockAdp)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdapterValidator_ValidateInput(t *testing.T) {
	validator := NewAdapterValidator()
	ctx := context.Background()

	tests := []struct {
		name      string
		input     *adapter.AdapterInput
		wantError bool
	}{
		{
			name: "valid input",
			input: &adapter.AdapterInput{
				Domain:       "k12",
				ActivityType: "concept_understanding",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
			},
			wantError: false,
		},
		{
			name:      "nil input",
			input:     nil,
			wantError: true,
		},
		{
			name: "missing domain",
			input: &adapter.AdapterInput{
				Domain:       "",
				ActivityType: "concept_understanding",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
			},
			wantError: true,
		},
		{
			name: "missing activity type",
			input: &adapter.AdapterInput{
				Domain:       "k12",
				ActivityType: "",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
			},
			wantError: true,
		},
		{
			name: "missing parameters",
			input: &adapter.AdapterInput{
				Domain:       "k12",
				ActivityType: "concept_understanding",
				Parameters:   nil,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateInput(ctx, tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdapterValidator_validateAdapterID(t *testing.T) {
	validator := NewAdapterValidator()

	tests := []struct {
		name      string
		adapterID string
		wantError bool
	}{
		{"valid simple", "k12.math", false},
		{"valid with underscores", "k12.math_problem_generator", false},
		{"valid with hyphens", "k12.math-problem-generator", false},
		{"valid with dots", "k12.math.problem.generator", false},
		{"empty string", "", true},
		{"starts with dot", ".k12.math", true},
		{"ends with dot", "k12.math.", true},
		{"uppercase letters", "K12.Math", true},
		{"spaces", "k12 math", true},
		{"special characters", "k12@math", true},
		{"single character", "a", true}, // 单个字符不符合模式（需要至少2个字符）
		{"very long", "k12.math_problem_generator_for_advanced_students", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateAdapterID(tt.adapterID)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdapterValidator_validateVersion(t *testing.T) {
	validator := NewAdapterValidator()

	tests := []struct {
		name      string
		version   string
		wantError bool
	}{
		{"valid 1.0.0", "1.0.0", false},
		{"valid 2.3.4", "2.3.4", false},
		{"valid with prerelease", "1.0.0-alpha", false},
		{"valid with prerelease and metadata", "1.0.0-beta.1+build.123", false},
		{"empty string", "", true},
		{"no patch version", "1.0", true},
		{"with v prefix", "v1.0.0", true},
		{"invalid format", "1.0.x", true},
		{"negative version", "-1.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateVersion(tt.version)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
