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

	compService "github.com/coze-dev/coze-studio/backend/domain/component/service"
)

func TestAdapterService_InstallAdapter(t *testing.T) {
	tests := []struct {
		name          string
		adapterID     string
		userID        int64
		setupMock     func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager)
		wantError     bool
		wantErrorType error
	}{
		{
			name:      "install adapter successfully",
			adapterID: "test.adapter",
			userID:    123,
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 注册 mock 组件
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "test.adapter",
				})

				// Mock 数据库更新
				mock.ExpectExec(`UPDATE "components"`).
					WithArgs(sqlmock.AnyArg(), "test.adapter").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			name:      "install adapter not found",
			adapterID: "nonexistent.adapter",
			userID:    123,
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 不注册组件，让 GetComponent 返回未找到错误
			},
			wantError:     true,
			wantErrorType: ErrAdapterNotFound,
		},
		{
			name:      "install adapter database error",
			adapterID: "test.adapter",
			userID:    123,
			setupMock: func(mock sqlmock.Sqlmock, compMgr *compService.ComponentManager) {
				// 注册 mock 组件
				registerMockComponent(compMgr, &mockComponentExecutor{
					componentID: "test.adapter",
				})

				// Mock 数据库更新失败
				mock.ExpectExec(`UPDATE "components"`).
					WithArgs(sqlmock.AnyArg(), "test.adapter").
					WillReturnError(errors.New("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mock, _, compMgr := setupTestService(t)

			if tt.setupMock != nil {
				tt.setupMock(mock, compMgr)
			}

			err := service.InstallAdapter(context.Background(), tt.adapterID, tt.userID)

			if tt.wantError {
				assert.Error(t, err)
				if tt.wantErrorType != nil {
					assert.ErrorIs(t, err, tt.wantErrorType)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证所有 mock 预期都被满足
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAdapterService_ListAdapters(t *testing.T) {
	tests := []struct {
		name          string
		request       *ListAdaptersRequest
		setupMock     func(mock sqlmock.Sqlmock)
		wantError     bool
		checkResponse func(t *testing.T, resp *ListAdaptersResponse)
	}{
		{
			name: "list all adapters with pagination",
			request: &ListAdaptersRequest{
				Page:     1,
				PageSize: 20,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(`SELECT count\(\*\) FROM "components"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

				// Mock SELECT query
				rows := sqlmock.NewRows([]string{
					"component_id", "name", "display_name", "description", "version",
					"type", "category", "status", "author_id", "author_name",
					"rating", "install_count", "usage_count",
				}).
					AddRow("test.adapter1", "adapter1", "Test Adapter 1", "Description 1", "1.0.0",
						"adapter", "test", "active", 1, "Author 1", 4.5, 100, 50).
					AddRow("test.adapter2", "adapter2", "Test Adapter 2", "Description 2", "1.0.0",
						"adapter", "test", "active", 2, "Author 2", 4.2, 80, 30).
					AddRow("test.adapter3", "adapter3", "Test Adapter 3", "Description 3", "1.1.0",
						"adapter", "test", "active", 3, "Author 3", 4.8, 150, 75)

				mock.ExpectQuery(`SELECT \* FROM "components"`).
					WillReturnRows(rows)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ListAdaptersResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, 3, resp.Total)
				assert.Equal(t, 1, resp.Page)
				assert.Equal(t, 20, resp.PageSize)
				assert.Len(t, resp.Adapters, 3)

				// 验证第一个适配器
				assert.Equal(t, "test.adapter1", resp.Adapters[0].AdapterID)
				assert.Equal(t, "Test Adapter 1", resp.Adapters[0].DisplayName)
				assert.Equal(t, float64(4.5), resp.Adapters[0].Rating)
				assert.Equal(t, 100, resp.Adapters[0].InstallCount)
			},
		},
		{
			name: "list adapters with category filter",
			request: &ListAdaptersRequest{
				Category: "education",
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query with category filter
				mock.ExpectQuery(`SELECT count\(\*\) FROM "components"`).
					WithArgs("adapter", "active", "education").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				// Mock SELECT query with category filter
				rows := sqlmock.NewRows([]string{
					"component_id", "name", "display_name", "description", "version",
					"type", "category", "status", "author_id", "author_name",
					"rating", "install_count", "usage_count",
				}).
					AddRow("k12.adapter", "k12_adapter", "K12 Adapter", "Education adapter", "1.0.0",
						"adapter", "education", "active", 1, "Test Author",
						4.5, 50, 100)

				mock.ExpectQuery(`SELECT \* FROM "components"`).
					WithArgs("adapter", "active", "education", 10).
					WillReturnRows(rows)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ListAdaptersResponse) {
				assert.Equal(t, 1, resp.Total)
				assert.Len(t, resp.Adapters, 1)
				assert.Equal(t, "education", resp.Adapters[0].Category)
			},
		},
		{
			name: "list adapters with tags filter",
			request: &ListAdaptersRequest{
				Tags:     []string{"math", "k12"},
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(`SELECT count\(\*\) FROM "components"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock SELECT query
				rows := sqlmock.NewRows([]string{
					"component_id", "name", "display_name", "description", "version",
					"type", "category", "status", "author_id", "author_name",
					"rating", "install_count", "usage_count",
				}).
					AddRow("math.adapter", "math_adapter", "Math Adapter", "Math problems", "1.0.0",
						"adapter", "education", "active", 1, "Math Author", 4.7, 150, 80).
					AddRow("k12.math", "k12_math", "K12 Math", "K12 math content", "1.0.0",
						"adapter", "education", "active", 2, "K12 Author", 4.6, 120, 60)

				mock.ExpectQuery(`SELECT \* FROM "components"`).
					WillReturnRows(rows)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ListAdaptersResponse) {
				assert.Equal(t, 2, resp.Total)
				assert.Len(t, resp.Adapters, 2)
			},
		},
		{
			name: "list adapters with search query",
			request: &ListAdaptersRequest{
				Query:    "math",
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query with search
				mock.ExpectQuery(`SELECT count\(\*\) FROM "components"`).
					WithArgs("adapter", "active", "math").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				// Mock SELECT query with ts_rank
				rows := sqlmock.NewRows([]string{
					"component_id", "name", "display_name", "description", "version",
					"type", "category", "status", "author_id", "author_name",
					"rating", "install_count", "usage_count",
				}).
					AddRow("math.adapter", "math_adapter", "Math Adapter", "Math problems generator", "1.0.0",
						"adapter", "education", "active", 1, "Math Author", 4.7, 150, 80)

				mock.ExpectQuery(`SELECT \* FROM "components"`).
					WithArgs("adapter", "active", "math", 10).
					WillReturnRows(rows)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ListAdaptersResponse) {
				assert.Equal(t, 1, resp.Total)
				assert.Len(t, resp.Adapters, 1)
				assert.Contains(t, resp.Adapters[0].DisplayName, "Math")
			},
		},
		{
			name: "list adapters empty result",
			request: &ListAdaptersRequest{
				Category: "nonexistent",
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(`SELECT count\(\*\) FROM "components"`).
					WithArgs("adapter", "active", "nonexistent").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock SELECT query (even though count is 0, GORM still executes it)
				rows := sqlmock.NewRows([]string{
					"component_id", "name", "display_name", "description", "version",
					"type", "category", "status", "author_id", "author_name",
					"rating", "install_count", "usage_count",
				})
				mock.ExpectQuery(`SELECT \* FROM "components"`).
					WithArgs("adapter", "active", "nonexistent", 10).
					WillReturnRows(rows)
			},
			wantError: false,
			checkResponse: func(t *testing.T, resp *ListAdaptersResponse) {
				assert.Equal(t, 0, resp.Total)
				assert.Len(t, resp.Adapters, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mock, _, _ := setupTestService(t)

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			resp, err := service.ListAdapters(context.Background(), tt.request)

			if tt.wantError {
				assert.Error(t, err)
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

