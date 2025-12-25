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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/validator"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/wrapper"
	"github.com/coze-dev/coze-studio/backend/domain/component"
	compEntity "github.com/coze-dev/coze-studio/backend/domain/component/entity"
	compService "github.com/coze-dev/coze-studio/backend/domain/component/service"
	pluginEntity "github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
	pluginModel "github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
)

// ErrValidationFailed 验证失败错误
var ErrValidationFailed = errors.New("adapter validation failed")

// AdapterService 适配器服务
// 基于 coze-studio 现有的 Plugin 和 Component 系统实现
type AdapterService struct {
	pluginRepo   repository.PluginRepository   // 复用 Plugin 仓储
	componentMgr *compService.ComponentManager // 复用 Component 管理器
	validator    *validator.AdapterValidator   // 适配器验证器
	db           *gorm.DB
	idGen        idgen.IDGenerator
}

// NewAdapterService 创建适配器服务
func NewAdapterService(
	pluginRepo repository.PluginRepository,
	componentMgr *compService.ComponentManager,
	validator *validator.AdapterValidator,
	db *gorm.DB,
	idGen idgen.IDGenerator,
) *AdapterService {
	return &AdapterService{
		pluginRepo:   pluginRepo,
		componentMgr: componentMgr,
		validator:    validator,
		db:           db,
		idGen:        idGen,
	}
}

// RegisterAdapterRequest 注册适配器请求
type RegisterAdapterRequest struct {
	// 基本信息
	AdapterID   string `json:"adapter_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Version     string `json:"version" binding:"required"`

	// 分类
	Category string   `json:"category" binding:"required"` // domain, visualization, etc.
	Tags     []string `json:"tags"`

	// 能力声明
	SupportedDomains       []string `json:"supported_domains"`
	SupportedActivityTypes []string `json:"supported_activity_types"`

	// 接口定义
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`

	// 适配器配置
	AdapterConfig *entity.AdapterConfig `json:"adapter_config"`

	// 实现（可以是代码或远程服务端点）
	Implementation adapter.SimplifiedDomainAdapter `json:"-"`                    // 不序列化
	ServerURL      string                          `json:"server_url,omitempty"` // 作者信息
	AuthorID       int64                           `json:"author_id"`            // Set by handler from context, not from request body
	AuthorName     string                          `json:"author_name"`

	// 元数据
	Homepage      string `json:"homepage"`
	Repository    string `json:"repository"`
	Documentation string `json:"documentation"`
	License       string `json:"license"`
	Icon          string `json:"icon"`
}

// RegisterAdapterResponse 注册适配器响应
type RegisterAdapterResponse struct {
	AdapterID   string `json:"adapter_id"`
	ComponentID string `json:"component_id"`
	PluginID    int64  `json:"plugin_id,omitempty"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

// RegisterAdapter 注册适配器
// 这是核心方法，将适配器注册到系统中
func (s *AdapterService) RegisterAdapter(ctx context.Context, req *RegisterAdapterRequest) (*RegisterAdapterResponse, error) {
	// 1. 验证适配器
	if req.Implementation != nil {
		if err := s.validator.ValidateAdapter(ctx, req.Implementation); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
		}
	}

	// 2. 检查是否已存在
	existing, err := s.componentMgr.GetComponent(req.AdapterID)
	if err == nil && existing != nil {
		return nil, ErrAdapterAlreadyExists
	}

	// 3. 创建 Plugin（如果需要远程调用）
	var pluginID int64
	if req.ServerURL != "" {
		pluginInfo := &pluginEntity.PluginInfo{
			PluginInfo: &pluginModel.PluginInfo{
				ServerURL: &req.ServerURL,
			},
		}

		pluginID, err = s.pluginRepo.CreateDraftPlugin(ctx, pluginInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to create plugin: %w", err)
		}
	}

	// 4. 创建 Component（适配器）
	config := make(map[string]interface{})
	if req.AdapterConfig != nil {
		config["adapter_config"] = req.AdapterConfig
	}
	if pluginID > 0 {
		config["plugin_id"] = pluginID
	}

	component := &compEntity.Component{
		ComponentID:   req.AdapterID,
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Type:          compEntity.ComponentTypeAdapter,
		Category:      compEntity.ComponentCategory(req.Category),
		Version:       req.Version,
		Description:   req.Description,
		AuthorID:      req.AuthorID,
		AuthorName:    req.AuthorName,
		Status:        compEntity.ComponentStatusActive,
		IsOfficial:    false,
		IsPublic:      true,
		Config:        compEntity.JSON(mustMarshal(config)),
		Tags:          compEntity.JSON(mustMarshal(req.Tags)),
		Icon:          req.Icon,
		Homepage:      req.Homepage,
		Repository:    req.Repository,
		Documentation: req.Documentation,
		License:       req.License,
	}

	// 5. 保存到数据库
	if err := s.db.WithContext(ctx).Create(component).Error; err != nil {
		return nil, fmt.Errorf("failed to create component: %w", err)
	}

	// 6. 如果提供了实现，创建包装器并注册到 ComponentRegistry
	if req.Implementation != nil {
		metadata := &adapter.AdapterMetadata{
			AdapterID:              req.AdapterID,
			Name:                   req.Name,
			DisplayName:            req.DisplayName,
			Version:                req.Version,
			Description:            req.Description,
			SupportedDomains:       req.SupportedDomains,
			SupportedActivityTypes: req.SupportedActivityTypes,
			InputSchema:            req.InputSchema,
			OutputSchema:           req.OutputSchema,
		}

		adapterWrapper := wrapper.NewDomainAdapterWrapper(
			pluginID,
			metadata,
			req.Implementation,
		)

		// 转换为 ComponentExecutor
		compExecutor := wrapper.AdapterToComponentExecutor(adapterWrapper)

		if err := s.componentMgr.RegisterComponent(compExecutor); err != nil {
			return nil, fmt.Errorf("failed to register component executor: %w", err)
		}
	}

	return &RegisterAdapterResponse{
		AdapterID:   req.AdapterID,
		ComponentID: req.AdapterID,
		PluginID:    pluginID,
		Status:      "registered",
		Message:     "Adapter registered successfully",
	}, nil
}

// InstallAdapter 安装适配器
func (s *AdapterService) InstallAdapter(ctx context.Context, adapterID string, userID int64) error {
	// 1. 获取适配器
	_, err := s.componentMgr.GetComponent(adapterID)
	if err != nil {
		return ErrAdapterNotFound
	}

	// 2. 记录安装
	// TODO: 创建 component_installations 表记录

	// 3. 更新安装计数
	if err := s.db.WithContext(ctx).
		Model(&compEntity.Component{}).
		Where("component_id = ?", adapterID).
		UpdateColumn("install_count", gorm.Expr("install_count + ?", 1)).
		Error; err != nil {
		return fmt.Errorf("failed to update install count: %w", err)
	}

	return nil
}

// ExecuteAdapterRequest 执行适配器请求
type ExecuteAdapterRequest struct {
	AdapterID  string                 `json:"adapter_id"`            // Set by handler from URL param, not from request body
	Parameters map[string]interface{} `json:"parameters" binding:"required"`
	Context    map[string]interface{} `json:"context"`
}

// ExecuteAdapterResponse 执行适配器响应
type ExecuteAdapterResponse struct {
	Result        map[string]interface{} `json:"result"`
	Metadata      map[string]interface{} `json:"metadata"`
	ExecutionTime int64                  `json:"execution_time"` // 毫秒
	Success       bool                   `json:"success"`
	Error         string                 `json:"error,omitempty"`
}

// ExecuteAdapter 执行适配器
func (s *AdapterService) ExecuteAdapter(ctx context.Context, req *ExecuteAdapterRequest) (*ExecuteAdapterResponse, error) {
	startTime := time.Now()

	// 1. 获取适配器
	comp, err := s.componentMgr.GetComponent(req.AdapterID)
	if err != nil {
		return nil, ErrAdapterNotFound
	}

	// 2. 准备输入
	input := &component.ComponentInput{
		ComponentID: req.AdapterID,
		Parameters:  req.Parameters,
		Context:     req.Context,
	}

	// 3. 执行
	output, err := comp.Execute(ctx, input)
	if err != nil {
		return &ExecuteAdapterResponse{
			Success:       false,
			Error:         err.Error(),
			ExecutionTime: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 4. 更新使用统计
	go func() {
		_ = s.db.Model(&compEntity.Component{}).
			Where("component_id = ?", req.AdapterID).
			UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).
			Error
	}()

	return &ExecuteAdapterResponse{
		Result:        output.Result,
		Metadata:      output.Metadata,
		ExecutionTime: output.ExecutionTime,
		Success:       output.Success,
		Error:         output.Error,
	}, nil
}

// ListAdaptersRequest 列表查询请求
type ListAdaptersRequest struct {
	Category string   `form:"category"`
	Tags     []string `form:"tags"`
	Domain   string   `form:"domain"`
	Query    string   `form:"query"` // 搜索关键词
	Page     int      `form:"page" binding:"min=1"`
	PageSize int      `form:"page_size" binding:"min=1,max=100"`
}

// ListAdaptersResponse 列表查询响应
type ListAdaptersResponse struct {
	Adapters   []*AdapterInfo `json:"adapters"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// AdapterInfo 适配器信息
type AdapterInfo struct {
	AdapterID    string    `json:"adapter_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Description  string    `json:"description"`
	Version      string    `json:"version"`
	Category     string    `json:"category"`
	Tags         []string  `json:"tags"`
	AuthorName   string    `json:"author_name"`
	Icon         string    `json:"icon"`
	Rating       float64   `json:"rating"`
	InstallCount int       `json:"install_count"`
	UsageCount   int64     `json:"usage_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListAdapters 列表查询
func (s *AdapterService) ListAdapters(ctx context.Context, req *ListAdaptersRequest) (*ListAdaptersResponse, error) {
	query := s.db.WithContext(ctx).
		Model(&compEntity.Component{}).
		Where("type = ?", compEntity.ComponentTypeAdapter).
		Where("status = ?", compEntity.ComponentStatusActive)

	// 分类过滤
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 标签过滤（JSON 查询）
	if len(req.Tags) > 0 {
		// PostgreSQL JSONB 查询
		query = query.Where("tags ?& array[?]", req.Tags)
	}

	// 全文搜索
	if req.Query != "" {
		// 使用 PostgreSQL 全文搜索
		query = query.Where(
			"to_tsvector('english', coalesce(name, '') || ' ' || coalesce(display_name, '') || ' ' || coalesce(description, '')) @@ plainto_tsquery('english', ?)",
			req.Query,
		)
	}

	// 总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var components []*compEntity.Component
	offset := (req.Page - 1) * req.PageSize

	baseQuery := query.Offset(offset).Limit(req.PageSize)

	// 如果有搜索，按相关度排序
	if req.Query != "" {
		orderClause := "ts_rank(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(display_name, '') || ' ' || coalesce(description, '')), plainto_tsquery('english', ?)) DESC, rating DESC"
		if err := baseQuery.Order(gorm.Expr(orderClause, req.Query)).Find(&components).Error; err != nil {
			return nil, err
		}
	} else {
		if err := baseQuery.Order("rating DESC, install_count DESC").Find(&components).Error; err != nil {
			return nil, err
		}
	}

	// 转换为 AdapterInfo
	adapters := make([]*AdapterInfo, len(components))
	for i, comp := range components {
		adapters[i] = &AdapterInfo{
			AdapterID:    comp.ComponentID,
			Name:         comp.Name,
			DisplayName:  comp.DisplayName,
			Description:  comp.Description,
			Version:      comp.Version,
			Category:     string(comp.Category),
			AuthorName:   comp.AuthorName,
			Icon:         comp.Icon,
			Rating:       comp.Rating,
			InstallCount: comp.InstallCount,
			UsageCount:   comp.UsageCount,
			CreatedAt:    comp.CreatedAt,
			UpdatedAt:    comp.UpdatedAt,
		}

		// 解析 Tags
		if len(comp.Tags) > 0 {
			_ = json.Unmarshal(comp.Tags, &adapters[i].Tags)
		}
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &ListAdaptersResponse{
		Adapters:   adapters,
		Total:      int(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// mustMarshal 辅助函数：序列化为 JSON，panic on error
func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
