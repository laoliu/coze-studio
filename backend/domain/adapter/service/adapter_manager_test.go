// Copyright 2025 Coze Studio. All rights reserved.

package service

import (
	"context"
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/stretchr/testify/assert"
)

// Mock DomainAdapter 用于测试的mock适配器
type MockDomainAdapter struct {
	info   *entity.AdapterInfo
	execFn func(ctx context.Context, input *adapter.AdapterInput) (*adapter.AdapterOutput, error)
}

func (m *MockDomainAdapter) GetInfo() *entity.AdapterInfo {
	return m.info
}

func (m *MockDomainAdapter) ParseRequest(input *entity.UserInput) (*entity.RequestContext, error) {
	return &entity.RequestContext{}, nil
}

func (m *MockDomainAdapter) ValidateContext(ctx *entity.RequestContext) error {
	return nil
}

func (m *MockDomainAdapter) GenerateLearningObjectives(ctx *entity.RequestContext) ([]*entity.LearningObjective, error) {
	return []*entity.LearningObjective{}, nil
}

func (m *MockDomainAdapter) DiscoverContent(ctx *entity.RequestContext, objectives []*entity.LearningObjective) ([]*entity.Content, error) {
	return []*entity.Content{}, nil
}

func (m *MockDomainAdapter) CustomizeWorkflow(ctx *entity.RequestContext) (*entity.WorkflowConfig, error) {
	return &entity.WorkflowConfig{}, nil
}

func (m *MockDomainAdapter) FormatOutput(result *entity.WorkflowResult) (interface{}, error) {
	return map[string]interface{}{"message": "mock success"}, nil
}

func (m *MockDomainAdapter) ValidateOutput(output interface{}) (*entity.QualityReport, error) {
	return &entity.QualityReport{}, nil
}

// TestAdapterRegistry_Register 测试注册适配器
func TestAdapterRegistry_Register(t *testing.T) {
	tests := []struct {
		name        string
		adapter     adapter.DomainAdapter
		expectedErr error
	}{
		{
			name: "register valid adapter",
			adapter: &MockDomainAdapter{
				info: &entity.AdapterInfo{
					AdapterID:   "test.adapter",
					Name:        "Test Adapter",
					Version:     "1.0.0",
					Description: "Test description",
					Type:        entity.AdapterTypeK12,
				},
			},
			expectedErr: nil,
		},
		{
			name:        "register nil adapter",
			adapter:     nil,
			expectedErr: ErrInvalidAdapter,
		},
		{
			name: "register adapter with nil info",
			adapter: &MockDomainAdapter{
				info: nil,
			},
			expectedErr: ErrInvalidAdapter,
		},
		{
			name: "register adapter with empty ID",
			adapter: &MockDomainAdapter{
				info: &entity.AdapterInfo{
					AdapterID: "",
					Name:      "Empty ID Adapter",
				},
			},
			expectedErr: ErrInvalidAdapter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewAdapterRegistry()
			err := registry.Register(tt.adapter)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// TestAdapterRegistry_RegisterDuplicate 测试重复注册
func TestAdapterRegistry_RegisterDuplicate(t *testing.T) {
	registry := NewAdapterRegistry()

	adapter1 := &MockDomainAdapter{
		info: &entity.AdapterInfo{
			AdapterID: "test.adapter",
			Name:      "Test Adapter",
			Version:   "1.0.0",
		},
	}

	// 第一次注册应该成功
	err := registry.Register(adapter1)
	assert.NoError(t, err)

	// 第二次注册相同ID应该失败
	err = registry.Register(adapter1)
	assert.Equal(t, ErrAdapterAlreadyExists, err)
}

// TestAdapterRegistry_Get 测试获取适配器
func TestAdapterRegistry_Get(t *testing.T) {
	registry := NewAdapterRegistry()

	adapter := &MockDomainAdapter{
		info: &entity.AdapterInfo{
			AdapterID: "test.adapter",
			Name:      "Test Adapter",
			Version:   "1.0.0",
		},
	}

	// 注册适配器
	err := registry.Register(adapter)
	assert.NoError(t, err)

	// 获取存在的适配器
	got, err := registry.Get("test.adapter")
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "test.adapter", got.GetInfo().AdapterID)

	// 获取不存在的适配器
	got, err = registry.Get("nonexistent.adapter")
	assert.Equal(t, ErrAdapterNotFound, err)
	assert.Nil(t, got)
}

// TestAdapterRegistry_GetInfo 测试获取适配器信息
func TestAdapterRegistry_GetInfo(t *testing.T) {
	registry := NewAdapterRegistry()

	adapter := &MockDomainAdapter{
		info: &entity.AdapterInfo{
			AdapterID:   "test.adapter",
			Name:        "Test Adapter",
			Version:     "1.0.0",
			Description: "Test description",
		},
	}

	// 注册适配器
	err := registry.Register(adapter)
	assert.NoError(t, err)

	// 获取存在的适配器信息
	info, err := registry.GetInfo("test.adapter")
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "test.adapter", info.AdapterID)
	assert.Equal(t, "Test Adapter", info.Name)
	assert.Equal(t, "1.0.0", info.Version)

	// 获取不存在的适配器信息
	info, err = registry.GetInfo("nonexistent.adapter")
	assert.Equal(t, ErrAdapterNotFound, err)
	assert.Nil(t, info)
}

// TestAdapterRegistry_Unregister 测试注销适配器
func TestAdapterRegistry_Unregister(t *testing.T) {
	registry := NewAdapterRegistry()

	adapter := &MockDomainAdapter{
		info: &entity.AdapterInfo{
			AdapterID: "test.adapter",
			Name:      "Test Adapter",
			Version:   "1.0.0",
		},
	}

	// 注册适配器
	err := registry.Register(adapter)
	assert.NoError(t, err)

	// 确认适配器存在
	_, err = registry.Get("test.adapter")
	assert.NoError(t, err)

	// 注销适配器
	err = registry.Unregister("test.adapter")
	assert.NoError(t, err)

	// 确认适配器已被删除
	_, err = registry.Get("test.adapter")
	assert.Equal(t, ErrAdapterNotFound, err)

	// 注销不存在的适配器
	err = registry.Unregister("nonexistent.adapter")
	assert.Equal(t, ErrAdapterNotFound, err)
}

// TestAdapterRegistry_List 测试列出所有适配器
func TestAdapterRegistry_List(t *testing.T) {
	registry := NewAdapterRegistry()

	// 空注册表
	list := registry.List()
	assert.Empty(t, list)

	// 注册多个适配器
	adapters := []*MockDomainAdapter{
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter1",
				Name:      "Adapter 1",
				Version:   "1.0.0",
				Type:      entity.AdapterTypeK12,
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter2",
				Name:      "Adapter 2",
				Version:   "2.0.0",
				Type:      entity.AdapterTypeArtHistory,
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter3",
				Name:      "Adapter 3",
				Version:   "1.5.0",
				Type:      entity.AdapterTypeK12,
			},
		},
	}

	for _, adp := range adapters {
		err := registry.Register(adp)
		assert.NoError(t, err)
	}

	// 列出所有适配器
	list = registry.List()
	assert.Len(t, list, 3)

	// 验证所有适配器都在列表中
	ids := make(map[string]bool)
	for _, info := range list {
		ids[info.AdapterID] = true
	}
	assert.True(t, ids["adapter1"])
	assert.True(t, ids["adapter2"])
	assert.True(t, ids["adapter3"])
}

// TestAdapterRegistry_FindByType 测试按类型查找适配器
func TestAdapterRegistry_FindByType(t *testing.T) {
	registry := NewAdapterRegistry()

	// 注册不同类型的适配器
	adapters := []*MockDomainAdapter{
		{
			info: &entity.AdapterInfo{
				AdapterID: "domain1",
				Name:      "Domain Adapter 1",
				Type:      entity.AdapterTypeK12,
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "domain2",
				Name:      "Domain Adapter 2",
				Type:      entity.AdapterTypeK12,
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "viz1",
				Name:      "Visualization Adapter",
				Type:      entity.AdapterTypeArtHistory,
			},
		},
	}

	for _, adp := range adapters {
		err := registry.Register(adp)
		assert.NoError(t, err)
	}

	// 查找 Domain 类型
	domainAdapters := registry.FindByType(entity.AdapterTypeK12)
	assert.Len(t, domainAdapters, 2)
	for _, info := range domainAdapters {
		assert.Equal(t, entity.AdapterTypeK12, info.Type)
	}

	// 查找 Visualization 类型
	vizAdapters := registry.FindByType(entity.AdapterTypeArtHistory)
	assert.Len(t, vizAdapters, 1)
	assert.Equal(t, "viz1", vizAdapters[0].AdapterID)

	// 查找不存在的类型
	analyticsAdapters := registry.FindByType(entity.AdapterTypeVocational)
	assert.Empty(t, analyticsAdapters)
}

// TestAdapterRegistry_FindByCapability 测试按能力查找适配器
func TestAdapterRegistry_FindByCapability(t *testing.T) {
	registry := NewAdapterRegistry()

	// 注册具有不同能力的适配器
	adapters := []*MockDomainAdapter{
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter1",
				Name:      "Adapter 1",
				Capabilities: &entity.AdapterCapabilities{
					ActivityTypes: []entity.ActivityType{
						entity.ActivityTypeConcept,
						entity.ActivityTypeExperiment,
					},
				},
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter2",
				Name:      "Adapter 2",
				Capabilities: &entity.AdapterCapabilities{
					ActivityTypes: []entity.ActivityType{
						entity.ActivityTypeProject,
					},
				},
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID: "adapter3",
				Name:      "Adapter 3",
				Capabilities: &entity.AdapterCapabilities{
					ActivityTypes: []entity.ActivityType{
						entity.ActivityTypeConcept,
					},
				},
			},
		},
		{
			info: &entity.AdapterInfo{
				AdapterID:    "adapter4",
				Name:         "Adapter 4",
				Capabilities: nil, // 无能力声明
			},
		},
	}

	for _, adp := range adapters {
		err := registry.Register(adp)
		assert.NoError(t, err)
	}

	// 查找支持 QuestionGeneration 的适配器
	questionAdapters := registry.FindByCapability(entity.ActivityTypeConcept)
	assert.Len(t, questionAdapters, 2)
	ids := make(map[string]bool)
	for _, info := range questionAdapters {
		ids[info.AdapterID] = true
	}
	assert.True(t, ids["adapter1"])
	assert.True(t, ids["adapter3"])

	// 查找支持 WorkflowCustomization 的适配器
	workflowAdapters := registry.FindByCapability(entity.ActivityTypeProject)
	assert.Len(t, workflowAdapters, 1)
	assert.Equal(t, "adapter2", workflowAdapters[0].AdapterID)

	// 查找支持 ContentDiscovery 的适配器
	contentAdapters := registry.FindByCapability(entity.ActivityTypeExperiment)
	assert.Len(t, contentAdapters, 1)
	assert.Equal(t, "adapter1", contentAdapters[0].AdapterID)
}

// TestAdapterRegistry_ConcurrentAccess 测试并发访问
func TestAdapterRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewAdapterRegistry()

	// 预先注册一些适配器
	for i := 0; i < 5; i++ {
		adapter := &MockDomainAdapter{
			info: &entity.AdapterInfo{
				AdapterID: "adapter" + string(rune('0'+i)),
				Name:      "Adapter " + string(rune('0'+i)),
			},
		}
		err := registry.Register(adapter)
		assert.NoError(t, err)
	}

	// 并发读取
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			// 读取操作
			list := registry.List()
			assert.GreaterOrEqual(t, len(list), 5)

			_, err := registry.Get("adapter0")
			assert.NoError(t, err)

			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestAdapterManager_NewAdapterManager 测试创建适配器管理器
func TestAdapterManager_NewAdapterManager(t *testing.T) {
	manager := NewAdapterManager()
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.registry)
}

// TestAdapterManager_RegisterAndGet 测试通过管理器注册和获取适配器
func TestAdapterManager_RegisterAndGet(t *testing.T) {
	manager := NewAdapterManager()

	adapter := &MockDomainAdapter{
		info: &entity.AdapterInfo{
			AdapterID: "test.adapter",
			Name:      "Test Adapter",
			Version:   "1.0.0",
		},
	}

	// 通过管理器注册
	err := manager.RegisterAdapter(adapter)
	assert.NoError(t, err)

	// 通过管理器获取
	got, err := manager.GetAdapter("test.adapter")
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "test.adapter", got.GetInfo().AdapterID)
}

// TestAdapterManager_ListAdapters 测试通过管理器列出适配器
func TestAdapterManager_ListAdapters(t *testing.T) {
	manager := NewAdapterManager()

	// 注册多个适配器
	for i := 0; i < 3; i++ {
		adapter := &MockDomainAdapter{
			info: &entity.AdapterInfo{
				AdapterID: "adapter" + string(rune('0'+i)),
				Name:      "Adapter " + string(rune('0'+i)),
			},
		}
		err := manager.RegisterAdapter(adapter)
		assert.NoError(t, err)
	}

	// 列出所有适配器
	list := manager.ListAdapters()
	assert.Len(t, list, 3)
}
