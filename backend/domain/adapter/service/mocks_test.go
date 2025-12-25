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

	pluginModel "github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/component"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
)

var (
	// ErrComponentNotFound 组件未找到错误
	ErrComponentNotFound = errors.New("component not found")
)

// MockPluginRepository 模拟 Plugin 仓储，实现完整的 repository.PluginRepository 接口
type MockPluginRepository struct {
	CreateDraftPluginFunc func(ctx context.Context, info *entity.PluginInfo) (int64, error)
	GetPluginFunc         func(ctx context.Context, pluginID int64) (*entity.PluginInfo, error)
	UpdatePluginFunc      func(ctx context.Context, info *entity.PluginInfo) error
	DeletePluginFunc      func(ctx context.Context, pluginID int64) error
	ListPluginsFunc       func(ctx context.Context, filters map[string]interface{}) ([]*entity.PluginInfo, error)
}

// CreateDraftPlugin 创建草稿插件
func (m *MockPluginRepository) CreateDraftPlugin(ctx context.Context, info *entity.PluginInfo) (int64, error) {
	if m.CreateDraftPluginFunc != nil {
		return m.CreateDraftPluginFunc(ctx, info)
	}
	return 1, nil
}

// CreateDraftPluginWithCode 创建带代码的草稿插件
func (m *MockPluginRepository) CreateDraftPluginWithCode(ctx context.Context, req *repository.CreateDraftPluginWithCodeRequest) (*repository.CreateDraftPluginWithCodeResponse, error) {
	return &repository.CreateDraftPluginWithCodeResponse{
		Plugin: &entity.PluginInfo{},
		Tools:  []*entity.ToolInfo{},
	}, nil
}

// GetDraftPlugin 获取草稿插件
func (m *MockPluginRepository) GetDraftPlugin(ctx context.Context, pluginID int64, opts ...repository.PluginSelectedOptions) (*entity.PluginInfo, bool, error) {
	if m.GetPluginFunc != nil {
		plugin, err := m.GetPluginFunc(ctx, pluginID)
		return plugin, plugin != nil, err
	}
	return &entity.PluginInfo{}, true, nil
}

// MGetDraftPlugins 批量获取草稿插件
func (m *MockPluginRepository) MGetDraftPlugins(ctx context.Context, pluginIDs []int64, opts ...repository.PluginSelectedOptions) ([]*entity.PluginInfo, error) {
	return []*entity.PluginInfo{}, nil
}

// GetAPPAllDraftPlugins 获取应用的所有草稿插件
func (m *MockPluginRepository) GetAPPAllDraftPlugins(ctx context.Context, appID int64, opts ...repository.PluginSelectedOptions) ([]*entity.PluginInfo, error) {
	return []*entity.PluginInfo{}, nil
}

// ListDraftPlugins 列出草稿插件
func (m *MockPluginRepository) ListDraftPlugins(ctx context.Context, req *repository.ListDraftPluginsRequest) (*repository.ListDraftPluginsResponse, error) {
	return &repository.ListDraftPluginsResponse{
		Plugins: []*entity.PluginInfo{},
		Total:   0,
	}, nil
}

// UpdateDraftPlugin 更新草稿插件
func (m *MockPluginRepository) UpdateDraftPlugin(ctx context.Context, info *entity.PluginInfo) error {
	if m.UpdatePluginFunc != nil {
		return m.UpdatePluginFunc(ctx, info)
	}
	return nil
}

// UpdateDraftPluginWithoutURLChanged 更新草稿插件（不改变URL）
func (m *MockPluginRepository) UpdateDraftPluginWithoutURLChanged(ctx context.Context, info *entity.PluginInfo) error {
	return nil
}

// UpdateDraftPluginWithCode 更新草稿插件代码
func (m *MockPluginRepository) UpdateDraftPluginWithCode(ctx context.Context, req *repository.UpdatePluginDraftWithCode) error {
	return nil
}

// DeleteDraftPlugin 删除草稿插件
func (m *MockPluginRepository) DeleteDraftPlugin(ctx context.Context, pluginID int64) error {
	if m.DeletePluginFunc != nil {
		return m.DeletePluginFunc(ctx, pluginID)
	}
	return nil
}

// DeleteAPPAllPlugins 删除应用的所有插件
func (m *MockPluginRepository) DeleteAPPAllPlugins(ctx context.Context, appID int64) ([]int64, error) {
	return []int64{}, nil
}

// UpdateDebugExample 更新调试示例
func (m *MockPluginRepository) UpdateDebugExample(ctx context.Context, pluginID int64, openapiDoc *pluginModel.Openapi3T) error {
	return nil
}

// GetOnlinePlugin 获取在线插件
func (m *MockPluginRepository) GetOnlinePlugin(ctx context.Context, pluginID int64, opts ...repository.PluginSelectedOptions) (*entity.PluginInfo, bool, error) {
	return &entity.PluginInfo{}, true, nil
}

// MGetOnlinePlugins 批量获取在线插件
func (m *MockPluginRepository) MGetOnlinePlugins(ctx context.Context, pluginIDs []int64, opts ...repository.PluginSelectedOptions) ([]*entity.PluginInfo, error) {
	return []*entity.PluginInfo{}, nil
}

// ListCustomOnlinePlugins 列出自定义在线插件
func (m *MockPluginRepository) ListCustomOnlinePlugins(ctx context.Context, spaceID int64, pageInfo dto.PageInfo) ([]*entity.PluginInfo, int64, error) {
	return []*entity.PluginInfo{}, 0, nil
}

// GetVersionPlugin 获取版本插件
func (m *MockPluginRepository) GetVersionPlugin(ctx context.Context, vPlugin pluginModel.VersionPlugin) (*entity.PluginInfo, bool, error) {
	return &entity.PluginInfo{}, true, nil
}

// MGetVersionPlugins 批量获取版本插件
func (m *MockPluginRepository) MGetVersionPlugins(ctx context.Context, vPlugins []pluginModel.VersionPlugin, opts ...repository.PluginSelectedOptions) ([]*entity.PluginInfo, error) {
	return []*entity.PluginInfo{}, nil
}

// PublishPlugin 发布插件
func (m *MockPluginRepository) PublishPlugin(ctx context.Context, draftPlugin *entity.PluginInfo) error {
	return nil
}

// PublishPlugins 批量发布插件
func (m *MockPluginRepository) PublishPlugins(ctx context.Context, draftPlugins []*entity.PluginInfo) error {
	return nil
}

// CopyPlugin 复制插件
func (m *MockPluginRepository) CopyPlugin(ctx context.Context, req *repository.CopyPluginRequest) (*entity.PluginInfo, []*entity.ToolInfo, error) {
	return req.Plugin, req.Tools, nil
}

// MoveAPPPluginToLibrary 移动应用插件到库
func (m *MockPluginRepository) MoveAPPPluginToLibrary(ctx context.Context, draftPlugin *entity.PluginInfo, draftTools []*entity.ToolInfo) error {
	return nil
}

// MockComponentManager 模拟组件管理器
type MockComponentManager struct {
	RegisterComponentFunc func(executor component.ComponentExecutor) error
	GetComponentFunc      func(componentID string) (component.ComponentExecutor, error)
	ListComponentsFunc    func() []component.ComponentExecutor
	UnregisterFunc        func(componentID string) error
}

func (m *MockComponentManager) RegisterComponent(executor component.ComponentExecutor) error {
	if m.RegisterComponentFunc != nil {
		return m.RegisterComponentFunc(executor)
	}
	return nil
}

func (m *MockComponentManager) GetComponent(componentID string) (component.ComponentExecutor, error) {
	if m.GetComponentFunc != nil {
		return m.GetComponentFunc(componentID)
	}
	return nil, ErrComponentNotFound
}

func (m *MockComponentManager) ListComponents() []component.ComponentExecutor {
	if m.ListComponentsFunc != nil {
		return m.ListComponentsFunc()
	}
	return []component.ComponentExecutor{}
}

func (m *MockComponentManager) Unregister(componentID string) error {
	if m.UnregisterFunc != nil {
		return m.UnregisterFunc(componentID)
	}
	return nil
}

// MockIDGenerator 模拟 ID 生成器
type MockIDGenerator struct {
	GenerateFunc func() int64
	GenIDFunc    func(ctx context.Context) (int64, error)
}

func (m *MockIDGenerator) Generate() int64 {
	if m.GenerateFunc != nil {
		return m.GenerateFunc()
	}
	return 12345
}

func (m *MockIDGenerator) GenID(ctx context.Context) (int64, error) {
	if m.GenIDFunc != nil {
		return m.GenIDFunc(ctx)
	}
	return 12345, nil
}

func (m *MockIDGenerator) GenMultiIDs(ctx context.Context, counts int) ([]int64, error) {
	ids := make([]int64, counts)
	for i := 0; i < counts; i++ {
		ids[i] = int64(12345 + i)
	}
	return ids, nil
}

// mockDomainAdapter 模拟领域适配器实现
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
	if m.validateError != nil {
		return m.validateError
	}
	return nil
}

func (m *mockDomainAdapter) Execute(ctx context.Context, input *adapter.AdapterInput) (*adapter.AdapterOutput, error) {
	if m.executeError != nil {
		return nil, m.executeError
	}
	if m.executeOutput != nil {
		return m.executeOutput, nil
	}
	return &adapter.AdapterOutput{
		Result:   map[string]interface{}{"status": "success"},
		Metadata: map[string]interface{}{"adapter_id": m.metadata.AdapterID},
	}, nil
}
