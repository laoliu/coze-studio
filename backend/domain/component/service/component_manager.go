// Copyright 2025 Coze Studio. All rights reserved.

package service

import (
	"context"
	"errors"
	"sync"

	"github.com/coze-dev/coze-studio/backend/domain/component"
)

var (
	ErrComponentNotFound      = errors.New("component not found")
	ErrComponentAlreadyExists = errors.New("component already exists")
	ErrInvalidComponent       = errors.New("invalid component")
)

// ComponentRegistry 组件注册表
type ComponentRegistry struct {
	components map[string]component.ComponentExecutor
	metadata   map[string]*component.ComponentInfo
	mu         sync.RWMutex
}

// NewComponentRegistry 创建组件注册表
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		components: make(map[string]component.ComponentExecutor),
		metadata:   make(map[string]*component.ComponentInfo),
	}
}

// Register 注册组件
func (r *ComponentRegistry) Register(comp component.ComponentExecutor) error {
	if comp == nil {
		return ErrInvalidComponent
	}

	info := comp.GetComponentInfo()
	if info == nil || info.ComponentID == "" {
		return ErrInvalidComponent
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[info.ComponentID]; exists {
		return ErrComponentAlreadyExists
	}

	r.components[info.ComponentID] = comp
	r.metadata[info.ComponentID] = info

	return nil
}

// Unregister 注销组件
func (r *ComponentRegistry) Unregister(componentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[componentID]; !exists {
		return ErrComponentNotFound
	}

	delete(r.components, componentID)
	delete(r.metadata, componentID)

	return nil
}

// Get 获取组件
func (r *ComponentRegistry) Get(componentID string) (component.ComponentExecutor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comp, exists := r.components[componentID]
	if !exists {
		return nil, ErrComponentNotFound
	}

	return comp, nil
}

// GetInfo 获取组件信息
func (r *ComponentRegistry) GetInfo(componentID string) (*component.ComponentInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.metadata[componentID]
	if !exists {
		return nil, ErrComponentNotFound
	}

	return info, nil
}

// List 列出所有已注册的组件
func (r *ComponentRegistry) List() []*component.ComponentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*component.ComponentInfo, 0, len(r.metadata))
	for _, info := range r.metadata {
		infos = append(infos, info)
	}

	return infos
}

// FindByType 按类型查找组件
func (r *ComponentRegistry) FindByType(compType string) []*component.ComponentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []*component.ComponentInfo
	for _, info := range r.metadata {
		if info.Type == compType {
			infos = append(infos, info)
		}
	}

	return infos
}

// FindByCategory 按分类查找组件
func (r *ComponentRegistry) FindByCategory(category string) []*component.ComponentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []*component.ComponentInfo
	for _, info := range r.metadata {
		if info.Category == category {
			infos = append(infos, info)
		}
	}

	return infos
}

// ComponentManager 组件管理器
type ComponentManager struct {
	registry *ComponentRegistry
}

// NewComponentManager 创建组件管理器
func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		registry: NewComponentRegistry(),
	}
}

// RegisterComponent 注册组件
func (m *ComponentManager) RegisterComponent(comp component.ComponentExecutor) error {
	return m.registry.Register(comp)
}

// GetComponent 获取组件
func (m *ComponentManager) GetComponent(componentID string) (component.ComponentExecutor, error) {
	return m.registry.Get(componentID)
}

// ListComponents 列出所有组件
func (m *ComponentManager) ListComponents() []*component.ComponentInfo {
	return m.registry.List()
}

// Execute 执行组件
func (m *ComponentManager) Execute(ctx context.Context, componentID string, input *component.ComponentInput) (*component.ComponentOutput, error) {
	comp, err := m.registry.Get(componentID)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if err := comp.Validate(input); err != nil {
		return nil, err
	}

	return comp.Execute(ctx, input)
}

// InvokeMCPTool 调用MCP工具
func (m *ComponentManager) InvokeMCPTool(ctx context.Context, toolName string, params map[string]any) (any, error) {
	comp, err := m.registry.Get(toolName)
	if err != nil {
		return nil, err
	}

	mcpTool, ok := comp.(component.MCPTool)
	if !ok {
		return nil, errors.New("component is not an MCP tool")
	}

	return mcpTool.InvokeTool(ctx, params)
}

// RenderVisualization 渲染可视化
func (m *ComponentManager) RenderVisualization(ctx context.Context, componentID string, data interface{}) (string, error) {
	comp, err := m.registry.Get(componentID)
	if err != nil {
		return "", err
	}

	vizComp, ok := comp.(component.VisualizationComponent)
	if !ok {
		return "", errors.New("component is not a visualization component")
	}

	return vizComp.Render(ctx, data)
}

// ExportContent 导出内容
func (m *ComponentManager) ExportContent(ctx context.Context, componentID string, content interface{}, format string) ([]byte, error) {
	comp, err := m.registry.Get(componentID)
	if err != nil {
		return nil, err
	}

	exportComp, ok := comp.(component.ExportComponent)
	if !ok {
		return nil, errors.New("component is not an export component")
	}

	return exportComp.Export(ctx, content, format)
}

// HandleInteraction 处理交互
func (m *ComponentManager) HandleInteraction(ctx context.Context, componentID string, event *component.InteractionEvent) (*component.InteractionResponse, error) {
	comp, err := m.registry.Get(componentID)
	if err != nil {
		return nil, err
	}

	interactionComp, ok := comp.(component.InteractionComponent)
	if !ok {
		return nil, errors.New("component is not an interaction component")
	}

	return interactionComp.HandleInteraction(ctx, event)
}
