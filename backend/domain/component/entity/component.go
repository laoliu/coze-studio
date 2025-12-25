// Copyright 2025 Coze Studio. All rights reserved.

package entity

import (
	"time"
)

// ComponentType 组件类型
type ComponentType string

const (
	ComponentTypeTool        ComponentType = "tool"         // 工具
	ComponentTypePlugin      ComponentType = "plugin"       // 插件
	ComponentTypeWidget      ComponentType = "widget"       // 小部件
	ComponentTypeIntegration ComponentType = "integration"  // 集成
	ComponentTypeAdapter     ComponentType = "adapter"      // 适配器 (NEW)
)

// ComponentCategory 组件分类
type ComponentCategory string

const (
	ComponentCategoryVisualization ComponentCategory = "visualization" // 可视化
	ComponentCategoryInteraction   ComponentCategory = "interaction"   // 交互
	ComponentCategoryAnalysis      ComponentCategory = "analysis"      // 分析
	ComponentCategoryExport        ComponentCategory = "export"        // 导出
	ComponentCategoryMCP           ComponentCategory = "mcp"           // MCP工具
	ComponentCategoryDomain        ComponentCategory = "domain"        // 领域适配器 (NEW)
)

// ComponentStatus 组件状态
type ComponentStatus string

const (
	ComponentStatusActive    ComponentStatus = "active"    // 激活
	ComponentStatusInactive  ComponentStatus = "inactive"  // 停用
	ComponentStatusDeprecated ComponentStatus = "deprecated" // 废弃
)

// Component 组件实体
type Component struct {
	ID            int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	ComponentID   string            `json:"component_id" gorm:"uniqueIndex;size:100;not null;comment:组件唯一标识"`
	Name          string            `json:"name" gorm:"size:100;not null;comment:组件名称"`
	DisplayName   string            `json:"display_name" gorm:"size:200;not null;comment:显示名称"`
	Type          ComponentType     `json:"type" gorm:"size:50;not null;comment:组件类型"`
	Category      ComponentCategory `json:"category" gorm:"size:50;not null;comment:组件分类"`
	Version       string            `json:"version" gorm:"size:50;not null;comment:版本号"`
	Description   string            `json:"description" gorm:"type:text;comment:描述"`

	// 作者信息
	AuthorID      int64             `json:"author_id" gorm:"index;comment:作者ID"`
	AuthorName    string            `json:"author_name" gorm:"size:100;comment:作者名称"`

	// 状态
	Status        ComponentStatus   `json:"status" gorm:"size:20;default:active;comment:状态"`
	IsOfficial    bool              `json:"is_official" gorm:"default:false;comment:是否官方"`
	IsPublic      bool              `json:"is_public" gorm:"default:false;comment:是否公开"`

	// 配置 (JSON存储)
	Config        JSON              `json:"config" gorm:"type:jsonb;comment:配置信息"`
	// {
	//   "parameters": {},
	//   "api_endpoint": "",
	//   "authentication": {}
	// }

	// 接口定义 (JSON存储)
	Interface     JSON              `json:"interface" gorm:"type:jsonb;comment:接口定义"`
	// {
	//   "inputs": [],
	//   "outputs": [],
	//   "methods": []
	// }

	// 依赖 (JSON存储)
	Dependencies  JSON              `json:"dependencies" gorm:"type:jsonb;comment:依赖项"`

	// 元数据
	Tags          JSON              `json:"tags" gorm:"type:jsonb;comment:标签"`
	Icon          string            `json:"icon" gorm:"size:500;comment:图标URL"`
	Homepage      string            `json:"homepage" gorm:"size:500;comment:主页"`
	Repository    string            `json:"repository" gorm:"size:500;comment:代码仓库"`
	Documentation string            `json:"documentation" gorm:"size:500;comment:文档地址"`
	License       string            `json:"license" gorm:"size:50;comment:许可证"`

	// 统计信息
	InstallCount  int               `json:"install_count" gorm:"default:0;comment:安装次数"`
	UsageCount    int64             `json:"usage_count" gorm:"default:0;comment:使用次数"`
	Rating        float64           `json:"rating" gorm:"default:0;comment:评分"`

	// 时间戳
	CreatedAt     time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	PublishedAt   *time.Time        `json:"published_at,omitempty" gorm:"comment:发布时间"`
}

// TableName 指定表名
func (Component) TableName() string {
	return "components"
}

// JSON 自定义JSON类型
type JSON []byte

// ComponentInterface 组件接口定义
type ComponentInterface struct {
	Inputs  []*Parameter `json:"inputs"`  // 输入参数
	Outputs []*Parameter `json:"outputs"` // 输出参数
	Methods []*Method    `json:"methods"` // 方法
}

// Parameter 参数定义
type Parameter struct {
	Name        string `json:"name"`        // 参数名
	Type        string `json:"type"`        // 类型
	Description string `json:"description"` // 描述
	Required    bool   `json:"required"`    // 是否必需
	Default     any    `json:"default"`     // 默认值
}

// Method 方法定义
type Method struct {
	Name        string       `json:"name"`        // 方法名
	Description string       `json:"description"` // 描述
	Inputs      []*Parameter `json:"inputs"`      // 输入
	Outputs     []*Parameter `json:"outputs"`     // 输出
}
