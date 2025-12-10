// Copyright 2025 Coze Studio. All rights reserved.

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/coze-studio/backend/domain/adapter"
	"github.com/coze-studio/backend/domain/adapter/entity"
)

var (
	ErrAdapterNotFound      = errors.New("adapter not found")
	ErrAdapterAlreadyExists = errors.New("adapter already exists")
	ErrAdapterNotRegistered = errors.New("adapter not registered")
	ErrInvalidAdapter       = errors.New("invalid adapter")
)

// AdapterRegistry 适配器注册表
// 负责适配器的注册、发现和管理
type AdapterRegistry struct {
	adapters map[string]adapter.DomainAdapter
	metadata map[string]*entity.AdapterInfo
	mu       sync.RWMutex
}

// NewAdapterRegistry 创建适配器注册表
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]adapter.DomainAdapter),
		metadata: make(map[string]*entity.AdapterInfo),
	}
}

// Register 注册适配器
func (r *AdapterRegistry) Register(adp adapter.DomainAdapter) error {
	if adp == nil {
		return ErrInvalidAdapter
	}

	info := adp.GetInfo()
	if info == nil || info.AdapterID == "" {
		return ErrInvalidAdapter
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adapters[info.AdapterID]; exists {
		return ErrAdapterAlreadyExists
	}

	r.adapters[info.AdapterID] = adp
	r.metadata[info.AdapterID] = info

	return nil
}

// Unregister 注销适配器
func (r *AdapterRegistry) Unregister(adapterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adapters[adapterID]; !exists {
		return ErrAdapterNotFound
	}

	delete(r.adapters, adapterID)
	delete(r.metadata, adapterID)

	return nil
}

// Get 获取适配器
func (r *AdapterRegistry) Get(adapterID string) (adapter.DomainAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adp, exists := r.adapters[adapterID]
	if !exists {
		return nil, ErrAdapterNotFound
	}

	return adp, nil
}

// GetInfo 获取适配器信息
func (r *AdapterRegistry) GetInfo(adapterID string) (*entity.AdapterInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.metadata[adapterID]
	if !exists {
		return nil, ErrAdapterNotFound
	}

	return info, nil
}

// List 列出所有已注册的适配器
func (r *AdapterRegistry) List() []*entity.AdapterInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*entity.AdapterInfo, 0, len(r.metadata))
	for _, info := range r.metadata {
		infos = append(infos, info)
	}

	return infos
}

// FindByType 按类型查找适配器
func (r *AdapterRegistry) FindByType(adapterType entity.AdapterType) []*entity.AdapterInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []*entity.AdapterInfo
	for _, info := range r.metadata {
		if info.Type == adapterType {
			infos = append(infos, info)
		}
	}

	return infos
}

// FindByCapability 按能力查找适配器
func (r *AdapterRegistry) FindByCapability(activityType entity.ActivityType) []*entity.AdapterInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []*entity.AdapterInfo
	for _, info := range r.metadata {
		if info.Capabilities != nil {
			for _, at := range info.Capabilities.ActivityTypes {
				if at == activityType {
					infos = append(infos, info)
					break
				}
			}
		}
	}

	return infos
}

// AdapterManager 适配器管理器
// 提供适配器的生命周期管理和调用服务
type AdapterManager struct {
	registry *AdapterRegistry
}

// NewAdapterManager 创建适配器管理器
func NewAdapterManager() *AdapterManager {
	return &AdapterManager{
		registry: NewAdapterRegistry(),
	}
}

// RegisterAdapter 注册适配器
func (m *AdapterManager) RegisterAdapter(adp adapter.DomainAdapter) error {
	return m.registry.Register(adp)
}

// GetAdapter 获取适配器
func (m *AdapterManager) GetAdapter(adapterID string) (adapter.DomainAdapter, error) {
	return m.registry.Get(adapterID)
}

// ListAdapters 列出所有适配器
func (m *AdapterManager) ListAdapters() []*entity.AdapterInfo {
	return m.registry.List()
}

// ParseRequest 解析用户请求（自动选择适配器）
func (m *AdapterManager) ParseRequest(ctx context.Context, adapterID string, input *entity.UserInput) (*entity.RequestContext, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	return adp.ParseRequest(input)
}

// GenerateLearningObjectives 生成学习目标
func (m *AdapterManager) GenerateLearningObjectives(ctx context.Context, adapterID string, reqCtx *entity.RequestContext) ([]*entity.LearningObjective, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	// 验证上下文
	if err := adp.ValidateContext(reqCtx); err != nil {
		return nil, fmt.Errorf("context validation failed: %w", err)
	}

	return adp.GenerateLearningObjectives(reqCtx)
}

// DiscoverContent 发现内容
func (m *AdapterManager) DiscoverContent(ctx context.Context, adapterID string, reqCtx *entity.RequestContext, objectives []*entity.LearningObjective) ([]*entity.Content, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	return adp.DiscoverContent(reqCtx, objectives)
}

// CustomizeWorkflow 定制工作流
func (m *AdapterManager) CustomizeWorkflow(ctx context.Context, adapterID string, reqCtx *entity.RequestContext) (*entity.WorkflowConfig, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	return adp.CustomizeWorkflow(reqCtx)
}

// FormatOutput 格式化输出
func (m *AdapterManager) FormatOutput(ctx context.Context, adapterID string, result *entity.WorkflowResult) (interface{}, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	return adp.FormatOutput(result)
}

// ValidateOutput 验证输出
func (m *AdapterManager) ValidateOutput(ctx context.Context, adapterID string, output interface{}) (*entity.QualityReport, error) {
	adp, err := m.registry.Get(adapterID)
	if err != nil {
		return nil, err
	}

	return adp.ValidateOutput(output)
}

// AutoSelectAdapter 自动选择适配器
// 根据用户输入自动选择最合适的适配器
func (m *AdapterManager) AutoSelectAdapter(ctx context.Context, input *entity.UserInput) (string, error) {
	// TODO: 实现智能适配器选择逻辑
	// 1. 根据领域匹配
	// 2. 根据活动类型匹配
	// 3. 使用AI意图识别
	
	// 目前简单实现：根据领域匹配
	infos := m.registry.List()
	for _, info := range infos {
		if info.Capabilities != nil {
			for _, domain := range info.Capabilities.Domains {
				if domain == input.Domain {
					return info.AdapterID, nil
				}
			}
		}
	}

	return "", ErrAdapterNotFound
}
