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

package coze

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/service"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/validator"
	compService "github.com/coze-dev/coze-studio/backend/domain/component/service"
	pluginRepo "github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	adapterService *service.AdapterService
)

// InitAdapterService 初始化适配器服务
func InitAdapterService(db *gorm.DB, idGen idgen.IDGenerator) {
	if adapterService == nil {
		// 创建依赖
		componentMgr := compService.NewComponentManager()
		adapterValidator := validator.NewAdapterValidator()
		pluginRepository := pluginRepo.NewPluginRepo(&pluginRepo.PluginRepoComponents{
			IDGen: idGen,
			DB:    db,
		})

		// 创建 AdapterService
		adapterService = service.NewAdapterService(
			pluginRepository,
			componentMgr,
			adapterValidator,
			db,
			idGen,
		)
		logs.Infof("Adapter service initialized successfully")
	}
}

// RegisterAdapter 注册适配器
func RegisterAdapter(ctx context.Context, c *app.RequestContext) {
	var req service.RegisterAdapterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户 ID（如果需要认证）
	// userID := c.GetInt64("user_id")
	// req.AuthorID = userID

	resp, err := adapterService.RegisterAdapter(ctx, &req)
	if err != nil {
		if err == service.ErrAdapterAlreadyExists {
			c.JSON(consts.StatusConflict, utils.H{
				"code": 409,
				"msg":  err.Error(),
			})
			return
		}
		logs.Errorf("Failed to register adapter: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to register adapter: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": resp,
	})
}

// ListAdapters 获取适配器列表
func ListAdapters(ctx context.Context, c *app.RequestContext) {
	var req service.ListAdaptersRequest

	// 解析查询参数
	if category := c.Query("category"); category != "" {
		req.Category = category
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

	resp, err := adapterService.ListAdapters(ctx, &req)
	if err != nil {
		logs.Errorf("Failed to list adapters: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to list adapters: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": resp,
	})
}

// GetAdapterInfo 获取适配器详情
func GetAdapterInfo(ctx context.Context, c *app.RequestContext) {
	adapterID := c.Param("adapter_id")
	if adapterID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Adapter ID is required",
		})
		return
	}

	// 获取适配器列表，然后筛选指定的适配器
	resp, err := adapterService.ListAdapters(ctx, &service.ListAdaptersRequest{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		logs.Errorf("Failed to list adapters: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to get adapter info: " + err.Error(),
		})
		return
	}

	// 查找指定的适配器
	for _, adapter := range resp.Adapters {
		if adapter.AdapterID == adapterID {
			c.JSON(consts.StatusOK, utils.H{
				"code": 0,
				"msg":  "success",
				"data": adapter,
			})
			return
		}
	}

	c.JSON(consts.StatusNotFound, utils.H{
		"code": 404,
		"msg":  "Adapter not found",
	})
}

// ExecuteAdapter 执行适配器
func ExecuteAdapter(ctx context.Context, c *app.RequestContext) {
	adapterID := c.Param("adapter_id")
	if adapterID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Adapter ID is required",
		})
		return
	}

	var req service.ExecuteAdapterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}
	req.AdapterID = adapterID

	resp, err := adapterService.ExecuteAdapter(ctx, &req)
	if err != nil {
		logs.Errorf("Failed to execute adapter: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to execute adapter: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": resp,
	})
}

// InstallAdapter 安装适配器
func InstallAdapter(ctx context.Context, c *app.RequestContext) {
	adapterID := c.Param("adapter_id")
	if adapterID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Adapter ID is required",
		})
		return
	}

	// 从上下文获取用户 ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code": 401,
			"msg":  "Unauthorized",
		})
		return
	}

	err := adapterService.InstallAdapter(ctx, adapterID, userID)
	if err != nil {
		logs.Errorf("Failed to install adapter: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to install adapter: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "Adapter installed successfully",
		"data": utils.H{
			"adapter_id": adapterID,
			"user_id":    userID,
		},
	})
}

// GetInstalledAdapters 获取已安装的适配器列表
func GetInstalledAdapters(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户 ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		// 如果没有用户 ID，返回所有可用的适配器
		userID = 0
	}

	// 获取所有已注册的适配器
	resp, err := adapterService.ListAdapters(ctx, &service.ListAdaptersRequest{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		logs.Errorf("Failed to list adapters: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to get installed adapters: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"user_id":  userID,
			"adapters": resp.Adapters,
			"total":    resp.Total,
		},
	})
}
