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

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/domain/component"
	compService "github.com/coze-dev/coze-studio/backend/domain/component/service"
)

func TestAdapterService_ExecuteAdapter(t *testing.T) {
	tests := []struct {
		name          string
		request       *ExecuteAdapterRequest
		setupMock     func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager)
		wantError     bool
		wantErrorType error
		checkResponse func(t *testing.T, resp *ExecuteAdapterResponse)
	}{
		{
			name: "execute adapter successfully",
			request: &ExecuteAdapterRequest{
				AdapterID: "test.adapter",
				Parameters: map[string]interface{}{
					"topic": "math",
					"grade": "5",
				},
				Context: map[string]interface{}{
					"domain":        "k12",
					"activity_type": "concept",
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 注册 mock 组件并返回成功的结果
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "test.adapter",
					executeResult: &component.ComponentOutput{
						ComponentID:   "test.adapter",
						Result:        map[string]interface{}{"status": "success", "content": "test"},
						Metadata:      map[string]interface{}{"adapter_id": "test.adapter"},
						ExecutionTime: 100,
						Success:       true,
					},
				})

				// Mock 数据库更新 (goroutine，不需要验证)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ExecuteAdapterResponse) {
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
				assert.Equal(t, "success", resp.Result["status"])
				assert.Equal(t, "test", resp.Result["content"])
				assert.Equal(t, int64(100), resp.ExecutionTime)
				assert.Empty(t, resp.Error)
			},
		},
		{
			name: "execute adapter not found",
			request: &ExecuteAdapterRequest{
				AdapterID: "nonexistent.adapter",
				Parameters: map[string]interface{}{
					"test": "value",
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 不注册组件，让 GetComponent 返回未找到错误
			},
			wantError:     true,
			wantErrorType: ErrAdapterNotFound,
		},
		{
			name: "execute adapter execution error",
			request: &ExecuteAdapterRequest{
				AdapterID: "test.adapter.fail",
				Parameters: map[string]interface{}{
					"invalid": "params",
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 注册会返回错误的 mock 组件
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID:  "test.adapter.fail",
					executeError: errors.New("execution failed: invalid parameters"),
				})
			},
			wantError: true,
			checkResponse: func(t *testing.T, resp *ExecuteAdapterResponse) {
				// 即使有错误，也应该返回响应对象
				assert.NotNil(t, resp)
				assert.False(t, resp.Success)
				assert.Contains(t, resp.Error, "execution failed")
				assert.GreaterOrEqual(t, resp.ExecutionTime, int64(0))
			},
		},
		{
			name: "execute adapter with empty parameters",
			request: &ExecuteAdapterRequest{
				AdapterID:  "test.adapter.empty",
				Parameters: map[string]interface{}{},
				Context:    map[string]interface{}{},
			},
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "test.adapter.empty",
					executeResult: &component.ComponentOutput{
						ComponentID:   "test.adapter.empty",
						Result:        map[string]interface{}{"status": "empty_params"},
						ExecutionTime: 10,
						Success:       true,
					},
				})
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ExecuteAdapterResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "empty_params", resp.Result["status"])
			},
		},
		{
			name: "execute adapter with complex result",
			request: &ExecuteAdapterRequest{
				AdapterID: "test.adapter.complex",
				Parameters: map[string]interface{}{
					"complexity": "high",
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "test.adapter.complex",
					executeResult: &component.ComponentOutput{
						ComponentID: "test.adapter.complex",
						Result: map[string]interface{}{
							"content": map[string]interface{}{
								"title":    "Complex Result",
								"sections": []interface{}{"intro", "body", "conclusion"},
								"metadata": map[string]interface{}{
									"word_count": 1000,
									"difficulty": "advanced",
								},
							},
						},
						Metadata: map[string]interface{}{
							"processing_steps": 5,
							"ai_calls":         3,
						},
						ExecutionTime: 500,
						Success:       true,
					},
				})
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ExecuteAdapterResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Result["content"])

				content := resp.Result["content"].(map[string]interface{})
				assert.Equal(t, "Complex Result", content["title"])

				assert.Equal(t, int64(500), resp.ExecutionTime)
				assert.Equal(t, 5, resp.Metadata["processing_steps"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mock, _, compMgr := setupTestService(t)

			if tt.setupMock != nil {
				tt.setupMock(mock, compMgr)
			}

			resp, err := service.ExecuteAdapter(context.Background(), tt.request)

			if tt.wantError {
				assert.Error(t, err)
				if tt.wantErrorType != nil {
					assert.ErrorIs(t, err, tt.wantErrorType)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, resp)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}

			// 给 goroutine 一点时间完成
			time.Sleep(10 * time.Millisecond)
		})
	}
}
