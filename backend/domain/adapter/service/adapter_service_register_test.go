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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/validator"
	"github.com/coze-dev/coze-studio/backend/domain/component"
	compService "github.com/coze-dev/coze-studio/backend/domain/component/service"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
)

// setupTestService 创建测试服务
func setupTestService(t *testing.T) (*AdapterService, sqlmock.Sqlmock, *MockPluginRepository, *compService.ComponentManager) {
	// 创建 mock 数据库
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	// 创建 mock 依赖
	mockPluginRepo := &MockPluginRepository{}
	componentMgr := compService.NewComponentManager() // 使用真实的 ComponentManager
	mockIDGen := &MockIDGenerator{}
	val := validator.NewAdapterValidator()

	service := NewAdapterService(
		mockPluginRepo,
		componentMgr,
		val,
		gormDB,
		mockIDGen,
	)

	return service, mock, mockPluginRepo, componentMgr
}

// registerMockComponent 注册一个 mock 组件到 ComponentManager
// 如果 executor 为 nil，则注册一个默认的 mock executor
func registerMockComponent(compMgr *compService.ComponentManager, executor component.ComponentExecutor) {
	if executor == nil {
		executor = &mockComponentExecutor{}
	}
	// 忽略注册错误，因为在测试中组件可能已经存在
	_ = compMgr.RegisterComponent(executor)
}

func TestAdapterService_RegisterAdapter(t *testing.T) {
	tests := []struct {
		name          string
		request       *RegisterAdapterRequest
		setupMock     func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager)
		wantError     bool
		wantErrorType error
		checkResponse func(t *testing.T, resp *RegisterAdapterResponse)
	}{
		{
			name: "register adapter without implementation",
			request: &RegisterAdapterRequest{
				AdapterID:   "test.adapter",
				Name:        "test_adapter",
				DisplayName: "Test Adapter",
				Version:     "1.0.0",
				Description: "A test adapter",
				Category:    "test",
				AuthorID:    1,
				AuthorName:  "Test Author",
				Tags:        []string{"test"},
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"input": map[string]interface{}{"type": "string"},
					},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"output": map[string]interface{}{"type": "string"},
					},
				},
				SupportedDomains:       []string{"test"},
				SupportedActivityTypes: []string{"test"},
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 不注册组件，让 GetComponent 返回未找到错误

				// Mock 数据库 INSERT (GORM's Create uses a query, not exec)
				mock.ExpectQuery(`INSERT INTO "components"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *RegisterAdapterResponse) {
				assert.Equal(t, "test.adapter", resp.AdapterID)
				assert.Equal(t, "test.adapter", resp.ComponentID)
				assert.Equal(t, "registered", resp.Status)
			},
		},
		{
			name: "register adapter with implementation",
			request: &RegisterAdapterRequest{
				AdapterID:   "test.adapter.impl",
				Name:        "test_adapter_impl",
				DisplayName: "Test Adapter with Implementation",
				Version:     "1.0.0",
				Description: "A test adapter with implementation",
				Category:    "test",
				AuthorID:    1,
				AuthorName:  "Test Author",
				Tags:        []string{"test"},
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"input": map[string]interface{}{"type": "string"},
					},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"output": map[string]interface{}{"type": "string"},
					},
				},
				SupportedDomains:       []string{"test"},
				SupportedActivityTypes: []string{"test"},
				Implementation: &mockDomainAdapter{
					metadata: &adapter.AdapterMetadata{
						AdapterID:   "test.adapter.impl",
						Name:        "test_adapter_impl",
						DisplayName: "Test Adapter with Implementation",
						Version:     "1.0.0",
						Description: "A test adapter with implementation",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"input": map[string]interface{}{"type": "string"},
							},
						},
						OutputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"output": map[string]interface{}{"type": "string"},
							},
						},
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 不注册组件，让 GetComponent 返回未找到

				// Mock 数据库 INSERT (GORM's Create uses a query with RETURNING)
				mock.ExpectQuery(`INSERT INTO "components"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *RegisterAdapterResponse) {
				assert.Equal(t, "test.adapter.impl", resp.AdapterID)
				assert.Equal(t, "registered", resp.Status)
			},
		},
		{
			name: "register adapter already exists",
			request: &RegisterAdapterRequest{
				AdapterID:   "existing.adapter",
				Name:        "existing_adapter",
				DisplayName: "Existing Adapter",
				Version:     "1.0.0",
				Description: "An existing adapter",
				Category:    "test",
				AuthorID:    1,
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 注册组件，模拟组件已存在
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "existing.adapter",
				})
			},
			wantError:     true,
			wantErrorType: ErrAdapterAlreadyExists,
		},
		{
			name: "register adapter with server URL creates plugin",
			request: &RegisterAdapterRequest{
				AdapterID:   "remote.adapter",
				Name:        "remote_adapter",
				DisplayName: "Remote Adapter",
				Version:     "1.0.0",
				Description: "A remote adapter",
				Category:    "test",
				AuthorID:    1,
				AuthorName:  "Test Author",
				ServerURL:   "https://example.com/adapter",
				InputSchema: map[string]interface{}{
					"type": "object",
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
				},
				SupportedDomains:       []string{"test"},
				SupportedActivityTypes: []string{"test"},
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 不注册组件

				// Mock CreateDraftPlugin
				pluginRepo.CreateDraftPluginFunc = func(ctx context.Context, info *entity.PluginInfo) (int64, error) {
					return 123, nil
				}

				// Mock 数据库 INSERT (GORM's Create uses a query with RETURNING)
				mock.ExpectQuery(`INSERT INTO "components"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *RegisterAdapterResponse) {
				assert.Equal(t, "remote.adapter", resp.AdapterID)
				assert.Equal(t, int64(123), resp.PluginID)
				assert.Equal(t, "registered", resp.Status)
			},
		},
		{
			name: "register adapter plugin creation fails",
			request: &RegisterAdapterRequest{
				AdapterID:              "remote.adapter.fail",
				Name:                   "remote_adapter_fail",
				DisplayName:            "Remote Adapter Fail",
				Version:                "1.0.0",
				Description:            "A remote adapter that fails",
				Category:               "test",
				AuthorID:               1,
				ServerURL:              "https://example.com/adapter",
				InputSchema:            map[string]interface{}{"type": "object"},
				OutputSchema:           map[string]interface{}{"type": "object"},
				SupportedDomains:       []string{"test"},
				SupportedActivityTypes: []string{"test"},
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 不注册组件

				// Mock CreateDraftPlugin - 失败
				pluginRepo.CreateDraftPluginFunc = func(ctx context.Context, info *entity.PluginInfo) (int64, error) {
					return 0, errors.New("plugin creation failed")
				}
			},
			wantError: true,
		},
		{
			name: "register adapter database error",
			request: &RegisterAdapterRequest{
				AdapterID:              "db.error.adapter",
				Name:                   "db_error_adapter",
				DisplayName:            "DB Error Adapter",
				Version:                "1.0.0",
				Description:            "Database error test",
				Category:               "test",
				AuthorID:               1,
				InputSchema:            map[string]interface{}{"type": "object"},
				OutputSchema:           map[string]interface{}{"type": "object"},
				SupportedDomains:       []string{"test"},
				SupportedActivityTypes: []string{"test"},
			},
			setupMock: func(mock sqlmock.Sqlmock, pluginRepo *MockPluginRepository, compMgr *compService.ComponentManager) {
				// 不注册组件

				// Mock 数据库 INSERT - 失败 (GORM's Create uses a query)
				mock.ExpectQuery(`INSERT INTO "components"`).
					WillReturnError(errors.New("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mock, pluginRepo, compMgr := setupTestService(t)

			if tt.setupMock != nil {
				tt.setupMock(mock, pluginRepo, compMgr)
			}

			resp, err := service.RegisterAdapter(context.Background(), tt.request)

			if tt.wantError {
				assert.Error(t, err)
				if tt.wantErrorType != nil {
					assert.ErrorIs(t, err, tt.wantErrorType)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}

			// 验证所有 mock 预期都被满足
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// mockComponentExecutor 用于测试的 mock ComponentExecutor
type mockComponentExecutor struct {
	componentID   string // 可配置的 ComponentID
	executeResult *component.ComponentOutput
	executeError  error
}

func (m *mockComponentExecutor) GetComponentInfo() *component.ComponentInfo {
	id := m.componentID
	if id == "" {
		id = "test" // 默认值
	}
	return &component.ComponentInfo{
		ComponentID: id,
		Name:        id,
		Version:     "1.0.0",
	}
}

func (m *mockComponentExecutor) Execute(ctx context.Context, input *component.ComponentInput) (*component.ComponentOutput, error) {
	if m.executeError != nil {
		return nil, m.executeError
	}
	if m.executeResult != nil {
		return m.executeResult, nil
	}
	return &component.ComponentOutput{Success: true}, nil
}

func (m *mockComponentExecutor) Validate(input *component.ComponentInput) error {
	return nil
}

func (m *mockComponentExecutor) GetConfig() map[string]interface{} {
	return map[string]interface{}{}
}
