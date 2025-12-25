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

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/service"
)

// MockAdapterService 模拟 AdapterService
type MockAdapterService struct {
	RegisterAdapterFunc func(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error)
	ListAdaptersFunc    func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error)
	InstallAdapterFunc  func(ctx context.Context, adapterID string, userID int64) error
	ExecuteAdapterFunc  func(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error)
}

func (m *MockAdapterService) RegisterAdapter(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error) {
	if m.RegisterAdapterFunc != nil {
		return m.RegisterAdapterFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *MockAdapterService) ListAdapters(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
	if m.ListAdaptersFunc != nil {
		return m.ListAdaptersFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *MockAdapterService) InstallAdapter(ctx context.Context, adapterID string, userID int64) error {
	if m.InstallAdapterFunc != nil {
		return m.InstallAdapterFunc(ctx, adapterID, userID)
	}
	return errors.New("not implemented")
}

func (m *MockAdapterService) ExecuteAdapter(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error) {
	if m.ExecuteAdapterFunc != nil {
		return m.ExecuteAdapterFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

// setupTestHandler 创建测试handler并注入mock service
func setupTestHandler(mockService *MockAdapterService) *AdapterHandler {
	return &AdapterHandler{
		adapterService: mockService,
	}
}

// setupTestRouter 创建测试路由
func setupTestRouter(handler *AdapterHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	handler.RegisterRoutes(api)
	return router
}

// createValidAdapterRequest 创建有效的适配器注册请求
func createValidAdapterRequest() map[string]interface{} {
	return map[string]interface{}{
		"adapter_id":   "test.adapter",
		"name":         "test_adapter",
		"display_name": "Test Adapter",
		"version":      "1.0.0",
		"description":  "A test adapter",
		"category":     "test",
	}
}

// TestAdapterHandler_RegisterAdapter 测试注册适配器
func TestAdapterHandler_RegisterAdapter(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAdapterService)
		setUserID      bool
		userID         int64
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "register adapter successfully",
			requestBody: map[string]interface{}{
				"adapter_id":  "test.adapter",
				"name":        "test_adapter",
				"display_name": "Test Adapter",
				"version":     "1.0.0",
				"description": "A test adapter",
				"category":    "test",
			},
			setupMock: func(m *MockAdapterService) {
				m.RegisterAdapterFunc = func(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error) {
					return &service.RegisterAdapterResponse{
						AdapterID:   "test.adapter",
						ComponentID: "test.adapter",
						Status:      "registered",
					}, nil
				}
			},
			setUserID:      true,
			userID:         123,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "test.adapter", data["adapter_id"])
				assert.Equal(t, "registered", data["status"])
			},
		},
		{
			name:           "register adapter without user_id",
			requestBody:    createValidAdapterRequest(),
			setupMock:      func(m *MockAdapterService) {},
			setUserID:      false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Unauthorized")
			},
		},
		{
			name:           "register adapter with invalid JSON",
			requestBody:    "invalid json",
			setupMock:      func(m *MockAdapterService) {},
			setUserID:      true,
			userID:         123,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Invalid request")
			},
		},
		{
			name:        "register adapter already exists",
			requestBody: createValidAdapterRequest(),
			setupMock: func(m *MockAdapterService) {
				m.RegisterAdapterFunc = func(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error) {
					return nil, service.ErrAdapterAlreadyExists
				}
			},
			setUserID:      true,
			userID:         123,
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "already exists")
			},
		},
		{
			name:        "register adapter service error",
			requestBody: createValidAdapterRequest(),
			setupMock: func(m *MockAdapterService) {
				m.RegisterAdapterFunc = func(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error) {
					return nil, errors.New("internal service error")
				}
			},
			setUserID:      true,
			userID:         123,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Failed to register adapter")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock service
			mockService := &MockAdapterService{}
			tt.setupMock(mockService)

			// 创建 handler
			handler := setupTestHandler(mockService)

			// 创建 router 并设置 user_id middleware
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// 添加middleware来设置 user_id
			if tt.setUserID {
				router.Use(func(c *gin.Context) {
					c.Set("user_id", tt.userID)
					c.Next()
				})
			}

			// 注册路由
			api := router.Group("/api")
			handler.RegisterRoutes(api)

			// 准备请求体
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// 创建请求
			req := httptest.NewRequest(http.MethodPost, "/api/adapters", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)


			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestAdapterHandler_ListAdapters 测试获取适配器列表
func TestAdapterHandler_ListAdapters(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(*MockAdapterService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "list adapters with default parameters",
			queryParams: "",
			setupMock: func(m *MockAdapterService) {
				m.ListAdaptersFunc = func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
					assert.Equal(t, 1, req.Page)
					assert.Equal(t, 20, req.PageSize)
					return &service.ListAdaptersResponse{
						Adapters: []*service.AdapterInfo{
							{AdapterID: "adapter1", Name: "Adapter 1"},
							{AdapterID: "adapter2", Name: "Adapter 2"},
						},
						Total: 2,
						Page:  1,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(2), data["total"])
				adapters := data["adapters"].([]interface{})
				assert.Len(t, adapters, 2)
			},
		},
		{
			name:        "list adapters with pagination",
			queryParams: "?page=2&page_size=10",
			setupMock: func(m *MockAdapterService) {
				m.ListAdaptersFunc = func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
					assert.Equal(t, 2, req.Page)
					assert.Equal(t, 10, req.PageSize)
					return &service.ListAdaptersResponse{
						Adapters: []*service.AdapterInfo{},
						Total:    0,
						Page:     2,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(0), data["total"])
			},
		},
		{
			name:        "list adapters with category filter",
			queryParams: "?category=education",
			setupMock: func(m *MockAdapterService) {
				m.ListAdaptersFunc = func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
					assert.Equal(t, "education", req.Category)
					return &service.ListAdaptersResponse{
						Adapters: []*service.AdapterInfo{
							{AdapterID: "edu.adapter", Category: "education"},
						},
						Total: 1,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(1), data["total"])
			},
		},
		{
			name:        "list adapters with tags filter",
			queryParams: "?tags=math&tags=k12",
			setupMock: func(m *MockAdapterService) {
				m.ListAdaptersFunc = func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
					assert.Contains(t, req.Tags, "math")
					assert.Contains(t, req.Tags, "k12")
					return &service.ListAdaptersResponse{
						Adapters: []*service.AdapterInfo{},
						Total:    0,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list adapters service error",
			queryParams: "",
			setupMock: func(m *MockAdapterService) {
				m.ListAdaptersFunc = func(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error) {
					return nil, errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Failed to list adapters")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock service
			mockService := &MockAdapterService{}
			tt.setupMock(mockService)

			// 创建 handler 和 router
			handler := setupTestHandler(mockService)
			// service injected via setupTestHandler
			router := setupTestRouter(handler)

			// 创建请求
			url := "/api/adapters" + tt.queryParams
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestAdapterHandler_InstallAdapter 测试安装适配器
func TestAdapterHandler_InstallAdapter(t *testing.T) {
	tests := []struct {
		name           string
		adapterID      string
		setUserID      bool
		userID         int64
		setupMock      func(*MockAdapterService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "install adapter successfully",
			adapterID: "test.adapter",
			setUserID: true,
			userID:    123,
			setupMock: func(m *MockAdapterService) {
				m.InstallAdapterFunc = func(ctx context.Context, adapterID string, userID int64) error {
					assert.Equal(t, "test.adapter", adapterID)
					assert.Equal(t, int64(123), userID)
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Contains(t, data["message"], "successfully")
			},
		},
		{
			name:           "install adapter without user_id",
			adapterID:      "test.adapter",
			setUserID:      false,
			setupMock:      func(m *MockAdapterService) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Unauthorized")
			},
		},
		{
			name:      "install adapter not found",
			adapterID: "nonexistent.adapter",
			setUserID: true,
			userID:    123,
			setupMock: func(m *MockAdapterService) {
				m.InstallAdapterFunc = func(ctx context.Context, adapterID string, userID int64) error {
					return service.ErrAdapterNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "not found")
			},
		},
		{
			name:      "install adapter service error",
			adapterID: "error.adapter",
			setUserID: true,
			userID:    123,
			setupMock: func(m *MockAdapterService) {
				m.InstallAdapterFunc = func(ctx context.Context, adapterID string, userID int64) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Failed to install adapter")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock service
			mockService := &MockAdapterService{}
			tt.setupMock(mockService)

			// 创建 handler
			handler := setupTestHandler(mockService)

			// 创建 router 并设置 middleware
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// 添加middleware来设置 user_id
			if tt.setUserID {
				router.Use(func(c *gin.Context) {
					c.Set("user_id", tt.userID)
					c.Next()
				})
			}

			// 注册路由
			api := router.Group("/api")
			handler.RegisterRoutes(api)

			// 创建请求
			url := "/api/adapters/" + tt.adapterID + "/install"
			req := httptest.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestAdapterHandler_ExecuteAdapter 测试执行适配器
func TestAdapterHandler_ExecuteAdapter(t *testing.T) {
	tests := []struct {
		name           string
		adapterID      string
		requestBody    interface{}
		setupMock      func(*MockAdapterService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "execute adapter successfully",
			adapterID: "test.adapter",
			requestBody: map[string]interface{}{
				"parameters": map[string]interface{}{
					"input": "test input",
				},
			},
			setupMock: func(m *MockAdapterService) {
				m.ExecuteAdapterFunc = func(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error) {
					assert.Equal(t, "test.adapter", req.AdapterID)
					return &service.ExecuteAdapterResponse{
						Result: map[string]interface{}{
							"output": "test output",
						},
						Success:       true,
						ExecutionTime: 100,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.True(t, data["success"].(bool))
				result := data["result"].(map[string]interface{})
				assert.Equal(t, "test output", result["output"])
			},
		},
		{
			name:      "execute adapter with invalid JSON",
			adapterID: "test.adapter",
			requestBody: "invalid json",
			setupMock:      func(m *MockAdapterService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Invalid request")
			},
		},
		{
			name:      "execute adapter not found",
			adapterID: "nonexistent.adapter",
			requestBody: map[string]interface{}{
				"parameters": map[string]interface{}{},
			},
			setupMock: func(m *MockAdapterService) {
				m.ExecuteAdapterFunc = func(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error) {
					return nil, service.ErrAdapterNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "not found")
			},
		},
		{
			name:      "execute adapter service error",
			adapterID: "error.adapter",
			requestBody: map[string]interface{}{
				"parameters": map[string]interface{}{},
			},
			setupMock: func(m *MockAdapterService) {
				m.ExecuteAdapterFunc = func(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error) {
					return nil, errors.New("execution failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp["message"], "Failed to execute adapter")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock service
			mockService := &MockAdapterService{}
			tt.setupMock(mockService)

			// 创建 handler 和 router
			handler := setupTestHandler(mockService)
			// service injected via setupTestHandler
			router := setupTestRouter(handler)

			// 准备请求体
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// 创建请求
			url := "/api/adapters/" + tt.adapterID + "/execute"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)


			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestAdapterHandler_GetAdapter 测试获取适配器详情
func TestAdapterHandler_GetAdapter(t *testing.T) {
	handler := setupTestHandler(&MockAdapterService{})
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/adapters/test.adapter", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// GetAdapter 还未实现，应该返回 501
	assert.Equal(t, http.StatusNotImplemented, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "Not implemented")
}
