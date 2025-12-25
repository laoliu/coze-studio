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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/component"
)

// mockDomainAdapter 模拟领域适配器
type mockDomainAdapter struct {
	metadata      *adapter.AdapterMetadata
	validateError error
	executeError  error
	executeOutput *adapter.AdapterOutput
}

func (m *mockDomainAdapter) GetMetadata() *adapter.AdapterMetadata {
	return m.metadata
}

func (m *mockDomainAdapter) Validate(ctx context.Context, input *adapter.AdapterInput) error {
	return m.validateError
}

func (m *mockDomainAdapter) Execute(ctx context.Context, input *adapter.AdapterInput) (*adapter.AdapterOutput, error) {
	if m.executeError != nil {
		return nil, m.executeError
	}
	if m.executeOutput != nil {
		return m.executeOutput, nil
	}
	return &adapter.AdapterOutput{
		Result: map[string]interface{}{
			"status": "success",
			"data":   "mock output",
		},
		Metadata: map[string]interface{}{
			"adapter_id": m.metadata.AdapterID,
		},
	}, nil
}

func TestDomainAdapterWrapper_Creation(t *testing.T) {
	metadata := &adapter.AdapterMetadata{
		AdapterID:   "k12.math_problem_generator",
		Name:        "math_problem_generator",
		DisplayName: "K12 Math Problem Generator",
		Version:     "1.0.0",
		Description: "Generate math problems",
	}

	mockAdapter := &mockDomainAdapter{
		metadata: metadata,
	}

	wrapper := NewDomainAdapterWrapper(123, metadata, mockAdapter)

	assert.NotNil(t, wrapper)
	assert.Equal(t, int64(123), wrapper.pluginID)
	assert.Equal(t, metadata, wrapper.GetMetadata())
	assert.Equal(t, mockAdapter, wrapper.GetAdapter())
}

func TestAdapterToComponentExecutor_GetMethods(t *testing.T) {
	metadata := &adapter.AdapterMetadata{
		AdapterID:   "k12.math_problem_generator",
		Name:        "math_problem_generator",
		DisplayName: "K12 Math Problem Generator",
		Version:     "1.0.0",
		Description: "Generate math problems for K12 students",
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
	}

	mockAdapter := &mockDomainAdapter{
		metadata: metadata,
	}

	wrapper := NewDomainAdapterWrapper(123, metadata, mockAdapter)
	executor := AdapterToComponentExecutor(wrapper)

	// Test GetComponentInfo
	info := executor.GetComponentInfo()
	assert.NotNil(t, info)
	assert.Equal(t, "k12.math_problem_generator", info.ComponentID)
	assert.Equal(t, "adapter", info.Type)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "Generate math problems for K12 students", info.Description)
	assert.NotNil(t, info.InputSchema)
	assert.NotNil(t, info.OutputSchema)

	// Test GetConfig
	config := executor.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "k12.math_problem_generator", config["adapter_id"])
	assert.Equal(t, "1.0.0", config["version"])
	assert.Equal(t, int64(123), config["plugin_id"])
}

func TestAdapterToComponentExecutor_Validate(t *testing.T) {
	metadata := &adapter.AdapterMetadata{
		AdapterID: "k12.math_problem_generator",
		Name:      "math_problem_generator",
		Version:   "1.0.0",
	}

	tests := []struct {
		name          string
		validateError error
		input         *component.ComponentInput
		wantError     bool
	}{
		{
			name:          "valid input",
			validateError: nil,
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
				Context: map[string]interface{}{
					"domain":        "k12",
					"activity_type": "concept_understanding",
				},
			},
			wantError: false,
		},
		{
			name:          "validation error",
			validateError: assert.AnError,
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters:  map[string]interface{}{},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := &mockDomainAdapter{
				metadata:      metadata,
				validateError: tt.validateError,
			}

			wrapper := NewDomainAdapterWrapper(123, metadata, mockAdapter)
			executor := AdapterToComponentExecutor(wrapper)

			err := executor.Validate(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdapterToComponentExecutor_Execute(t *testing.T) {
	metadata := &adapter.AdapterMetadata{
		AdapterID: "k12.math_problem_generator",
		Name:      "math_problem_generator",
		Version:   "1.0.0",
	}

	tests := []struct {
		name          string
		executeError  error
		executeOutput *adapter.AdapterOutput
		input         *component.ComponentInput
		wantSuccess   bool
		wantError     bool
	}{
		{
			name:         "successful execution",
			executeError: nil,
			executeOutput: &adapter.AdapterOutput{
				Result: map[string]interface{}{
					"problem": "Solve: 2x + 3 = 7",
					"answer":  "x = 2",
				},
				Metadata: map[string]interface{}{
					"difficulty": "medium",
				},
			},
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
				Context: map[string]interface{}{
					"domain":        "k12",
					"activity_type": "concept_understanding",
				},
			},
			wantSuccess: true,
			wantError:   false,
		},
		{
			name:         "execution error",
			executeError: assert.AnError,
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters: map[string]interface{}{
					"topic": "invalid_topic",
				},
			},
			wantSuccess: false,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := &mockDomainAdapter{
				metadata:      metadata,
				executeError:  tt.executeError,
				executeOutput: tt.executeOutput,
			}

			wrapper := NewDomainAdapterWrapper(123, metadata, mockAdapter)
			executor := AdapterToComponentExecutor(wrapper)

			output, err := executor.Execute(context.Background(), tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.NotNil(t, output)
				assert.False(t, output.Success)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, output)
				assert.Equal(t, tt.wantSuccess, output.Success)

				if tt.wantSuccess {
					assert.NotNil(t, output.Result)
					assert.NotNil(t, output.Metadata)
					assert.GreaterOrEqual(t, output.ExecutionTime, int64(0))
				}
			}
		})
	}
}

func TestToAdapterInput(t *testing.T) {
	tests := []struct {
		name     string
		input    *component.ComponentInput
		expected *adapter.AdapterInput
	}{
		{
			name: "full context",
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
				Context: map[string]interface{}{
					"domain":        "k12",
					"activity_type": "concept_understanding",
					"extra_data":    "some value",
				},
			},
			expected: &adapter.AdapterInput{
				Domain:       "k12",
				ActivityType: "concept_understanding",
				Parameters: map[string]interface{}{
					"topic": "algebra",
				},
				Context: map[string]interface{}{
					"domain":        "k12",
					"activity_type": "concept_understanding",
					"extra_data":    "some value",
				},
			},
		},
		{
			name: "minimal context",
			input: &component.ComponentInput{
				ComponentID: "k12.math_problem_generator",
				Parameters: map[string]interface{}{
					"topic": "geometry",
				},
				Context: nil,
			},
			expected: &adapter.AdapterInput{
				Domain:       "",
				ActivityType: "",
				Parameters: map[string]interface{}{
					"topic": "geometry",
				},
				Context: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toAdapterInput(tt.input)

			assert.Equal(t, tt.expected.Domain, result.Domain)
			assert.Equal(t, tt.expected.ActivityType, result.ActivityType)
			assert.Equal(t, tt.expected.Parameters, result.Parameters)
			assert.Equal(t, tt.expected.Context, result.Context)
		})
	}
}
