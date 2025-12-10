// Copyright 2025 Coze Studio. All rights reserved.

package entity

import (
	"encoding/json"
	"time"
)

// AdapterType 适配器类型
type AdapterType string

const (
	AdapterTypeK12         AdapterType = "k12_education"      // K12教育
	AdapterTypeArtHistory  AdapterType = "art_history"        // 美术史
	AdapterTypeVocational  AdapterType = "vocational_training" // 职业培训
	AdapterTypeCustom      AdapterType = "custom"             // 自定义
)

// AdapterStatus 适配器状态
type AdapterStatus string

const (
	AdapterStatusDraft     AdapterStatus = "draft"     // 草稿
	AdapterStatusActive    AdapterStatus = "active"    // 激活
	AdapterStatusInactive  AdapterStatus = "inactive"  // 停用
	AdapterStatusArchived  AdapterStatus = "archived"  // 归档
)

// ActivityType 活动类型
type ActivityType string

const (
	ActivityTypeConcept       ActivityType = "concept"        // 概念理解
	ActivityTypeExperiment    ActivityType = "experiment"     // 实验探究
	ActivityTypeProblem       ActivityType = "problem"        // 问题解决
	ActivityTypeProject       ActivityType = "project"        // 项目制学习
	ActivityTypeDiscussion    ActivityType = "discussion"     // 讨论辩论
	ActivityTypeCreation      ActivityType = "creation"       // 创作表达
	ActivityTypeReflection    ActivityType = "reflection"     // 反思总结
)

// Adapter 领域适配器实体
type Adapter struct {
	ID             int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	AdapterID      string        `json:"adapter_id" gorm:"uniqueIndex;size:100;not null;comment:适配器唯一标识"`
	Name           string        `json:"name" gorm:"size:100;not null;comment:适配器名称"`
	DisplayName    string        `json:"display_name" gorm:"size:200;not null;comment:显示名称"`
	Type           AdapterType   `json:"type" gorm:"size:50;not null;comment:适配器类型"`
	Version        string        `json:"version" gorm:"size:50;not null;comment:版本号"`
	Description    string        `json:"description" gorm:"type:text;comment:描述"`
	
	// 作者信息
	AuthorID       int64         `json:"author_id" gorm:"index;comment:作者ID"`
	AuthorName     string        `json:"author_name" gorm:"size:100;comment:作者名称"`
	Organization   string        `json:"organization" gorm:"size:200;comment:组织名称"`
	
	// 状态
	Status         AdapterStatus `json:"status" gorm:"size:20;default:draft;comment:状态"`
	IsOfficial     bool          `json:"is_official" gorm:"default:false;comment:是否官方"`
	IsPublic       bool          `json:"is_public" gorm:"default:false;comment:是否公开"`
	
	// 能力声明 (JSON存储)
	Capabilities   JSON          `json:"capabilities" gorm:"type:jsonb;comment:能力声明"`
	// {
	//   "activity_types": ["concept", "experiment"],
	//   "domains": ["chemistry", "physics"],
	//   "features": ["rag", "multimodal"]
	// }
	
	// 配置 (JSON存储)
	Config         JSON          `json:"config" gorm:"type:jsonb;comment:配置信息"`
	// {
	//   "prompt_templates": {},
	//   "quality_rules": {},
	//   "content_sources": []
	// }
	
	// 依赖 (JSON存储)
	Dependencies   JSON          `json:"dependencies" gorm:"type:jsonb;comment:依赖项"`
	// {
	//   "adapters": [],
	//   "capabilities": ["llm_client", "vector_store"],
	//   "components": []
	// }
	
	// 元数据
	Manifest       JSON          `json:"manifest" gorm:"type:jsonb;comment:完整清单"`
	Tags           JSON          `json:"tags" gorm:"type:jsonb;comment:标签"`
	Homepage       string        `json:"homepage" gorm:"size:500;comment:主页"`
	Repository     string        `json:"repository" gorm:"size:500;comment:代码仓库"`
	Documentation  string        `json:"documentation" gorm:"size:500;comment:文档地址"`
	License        string        `json:"license" gorm:"size:50;comment:许可证"`
	
	// 统计信息
	InstallCount   int           `json:"install_count" gorm:"default:0;comment:安装次数"`
	UsageCount     int64         `json:"usage_count" gorm:"default:0;comment:使用次数"`
	Rating         float64       `json:"rating" gorm:"default:0;comment:评分"`
	ReviewCount    int           `json:"review_count" gorm:"default:0;comment:评论数"`
	
	// 时间戳
	CreatedAt      time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	PublishedAt    *time.Time    `json:"published_at,omitempty" gorm:"comment:发布时间"`
	ArchivedAt     *time.Time    `json:"archived_at,omitempty" gorm:"comment:归档时间"`
}

// TableName 指定表名
func (Adapter) TableName() string {
	return "adapters"
}

// JSON 自定义JSON类型
type JSON json.RawMessage

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*j = JSON(bytes)
	return nil
}

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = JSON(data)
	return nil
}

// AdapterCapabilities 适配器能力声明
type AdapterCapabilities struct {
	ActivityTypes []ActivityType `json:"activity_types"` // 支持的活动类型
	Domains       []string       `json:"domains"`        // 支持的领域
	Features      []string       `json:"features"`       // 支持的特性
	Grades        []string       `json:"grades"`         // 适用年级（K12专用）
	Languages     []string       `json:"languages"`      // 支持的语言
}

// AdapterConfig 适配器配置
type AdapterConfig struct {
	PromptTemplates map[string]string `json:"prompt_templates"` // 提示词模板
	QualityRules    map[string]any    `json:"quality_rules"`    // 质量规则
	ContentSources  []ContentSource   `json:"content_sources"`  // 内容源
}

// ContentSource 内容源配置
type ContentSource struct {
	Type     string         `json:"type"`     // 类型: knowledge_base, api, file
	Name     string         `json:"name"`     // 名称
	Config   map[string]any `json:"config"`   // 配置
	Priority int            `json:"priority"` // 优先级
}

// AdapterDependencies 适配器依赖
type AdapterDependencies struct {
	Adapters     []string `json:"adapters"`     // 依赖的适配器
	Capabilities []string `json:"capabilities"` // 依赖的能力
	Components   []string `json:"components"`   // 依赖的组件
}
