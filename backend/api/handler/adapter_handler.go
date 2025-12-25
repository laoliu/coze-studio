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
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/service"
	"github.com/coze-dev/coze-studio/backend/pkg/response"
)

// AdapterHandler 适配器处理器
type AdapterHandler struct {
	adapterService AdapterServiceInterface
}

// AdapterServiceInterface 定义适配器服务接口，便于测试
type AdapterServiceInterface interface {
	RegisterAdapter(ctx context.Context, req *service.RegisterAdapterRequest) (*service.RegisterAdapterResponse, error)
	ListAdapters(ctx context.Context, req *service.ListAdaptersRequest) (*service.ListAdaptersResponse, error)
	InstallAdapter(ctx context.Context, adapterID string, userID int64) error
	ExecuteAdapter(ctx context.Context, req *service.ExecuteAdapterRequest) (*service.ExecuteAdapterResponse, error)
}

// NewAdapterHandler 创建处理器
func NewAdapterHandler(adapterService *service.AdapterService) *AdapterHandler {
	return &AdapterHandler{
		adapterService: adapterService,
	}
}

// RegisterRoutes 注册路由
func (h *AdapterHandler) RegisterRoutes(r *gin.RouterGroup) {
	adapters := r.Group("/adapters")
	{
		adapters.POST("", h.RegisterAdapter)
		adapters.GET("", h.ListAdapters)
		adapters.GET("/:id", h.GetAdapter)
		adapters.POST("/:id/install", h.InstallAdapter)
		adapters.POST("/:id/execute", h.ExecuteAdapter)
	}
}

// RegisterAdapter 注册适配器
// @Summary 注册适配器
// @Description 注册新的领域适配器到系统中
// @Tags adapters
// @Accept json
// @Produce json
// @Param request body service.RegisterAdapterRequest true "注册请求"
// @Success 200 {object} service.RegisterAdapterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/adapters [post]
func (h *AdapterHandler) RegisterAdapter(c *gin.Context) {
	var req service.RegisterAdapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	req.AuthorID = userID.(int64)

	resp, err := h.adapterService.RegisterAdapter(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrAdapterAlreadyExists {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register adapter: "+err.Error())
		return
	}

	response.Success(c, resp)
}

// ListAdapters 获取适配器列表
// @Summary 获取适配器列表
// @Description 分页查询适配器列表，支持按分类、标签、领域过滤
// @Tags adapters
// @Accept json
// @Produce json
// @Param category query string false "分类"
// @Param tags query []string false "标签"
// @Param domain query string false "领域"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} service.ListAdaptersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/adapters [get]
func (h *AdapterHandler) ListAdapters(c *gin.Context) {
	var req service.ListAdaptersRequest

	// 解析查询参数
	if category := c.Query("category"); category != "" {
		req.Category = category
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		req.Tags = tags
	}
	if domain := c.Query("domain"); domain != "" {
		req.Domain = domain
	}

	// 分页参数
	req.Page = 1
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	req.PageSize = 20
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			req.PageSize = ps
		}
	}

	resp, err := h.adapterService.ListAdapters(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list adapters: "+err.Error())
		return
	}

	response.Success(c, resp)
}

// GetAdapter 获取适配器详情
// @Summary 获取适配器详情
// @Description 获取指定适配器的详细信息
// @Tags adapters
// @Accept json
// @Produce json
// @Param id path string true "适配器 ID"
// @Success 200 {object} service.AdapterInfo
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/adapters/{id} [get]
func (h *AdapterHandler) GetAdapter(c *gin.Context) {
	adapterID := c.Param("id")
	if adapterID == "" {
		response.Error(c, http.StatusBadRequest, "Adapter ID is required")
		return
	}

	// TODO: 实现 GetAdapter 方法
	response.Error(c, http.StatusNotImplemented, "Not implemented yet")
}

// InstallAdapter 安装适配器
// @Summary 安装适配器
// @Description 为当前用户安装指定的适配器
// @Tags adapters
// @Accept json
// @Produce json
// @Param id path string true "适配器 ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/adapters/{id}/install [post]
func (h *AdapterHandler) InstallAdapter(c *gin.Context) {
	adapterID := c.Param("id")
	if adapterID == "" {
		response.Error(c, http.StatusBadRequest, "Adapter ID is required")
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.adapterService.InstallAdapter(c.Request.Context(), adapterID, userID.(int64))
	if err != nil {
		if err == service.ErrAdapterNotFound {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to install adapter: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Adapter installed successfully",
	})
}

// ExecuteAdapter 执行适配器
// @Summary 执行适配器
// @Description 执行指定适配器的功能
// @Tags adapters
// @Accept json
// @Produce json
// @Param id path string true "适配器 ID"
// @Param request body service.ExecuteAdapterRequest true "执行请求"
// @Success 200 {object} service.ExecuteAdapterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/adapters/{id}/execute [post]
func (h *AdapterHandler) ExecuteAdapter(c *gin.Context) {
	adapterID := c.Param("id")
	if adapterID == "" {
		response.Error(c, http.StatusBadRequest, "Adapter ID is required")
		return
	}

	var req service.ExecuteAdapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	req.AdapterID = adapterID

	resp, err := h.adapterService.ExecuteAdapter(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrAdapterNotFound {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to execute adapter: "+err.Error())
		return
	}

	response.Success(c, resp)
}
