// Copyright 2025 Coze Studio. All rights reserved.

package entity

import (
	"time"
)

// CapabilityType AI能力类型
type CapabilityType string

const (
	CapabilityTypeIntentRecognition   CapabilityType = "intent_recognition"   // 意图识别
	CapabilityTypeObjectiveGeneration CapabilityType = "objective_generation" // 学习目标生成
	CapabilityTypeContentDiscovery    CapabilityType = "content_discovery"    // 内容发现
	CapabilityTypeNarrativeGeneration CapabilityType = "narrative_generation" // 故事化叙述
	CapabilityTypeQualityAssessment   CapabilityType = "quality_assessment"   // 质量评估
	CapabilityTypeMultimodalGeneration CapabilityType = "multimodal_generation" // 多模态生成
	CapabilityTypePersonaModeling     CapabilityType = "persona_modeling"     // 画像建模
	CapabilityTypeKnowledgeMapping    CapabilityType = "knowledge_mapping"    // 知识映射
)

// CapabilityStatus 能力状态
type CapabilityStatus string

const (
	CapabilityStatusEnabled  CapabilityStatus = "enabled"  // 启用
	CapabilityStatusDisabled CapabilityStatus = "disabled" // 禁用
	CapabilityStatusBeta     CapabilityStatus = "beta"     // 测试
)

// Capability AI能力实体
type Capability struct {
	ID             int64            `json:"id" gorm:"primaryKey;autoIncrement"`
	CapabilityID   string           `json:"capability_id" gorm:"uniqueIndex;size:100;not null;comment:能力唯一标识"`
	Name           string           `json:"name" gorm:"size:100;not null;comment:能力名称"`
	DisplayName    string           `json:"display_name" gorm:"size:200;not null;comment:显示名称"`
	Type           CapabilityType   `json:"type" gorm:"size:50;not null;comment:能力类型"`
	Version        string           `json:"version" gorm:"size:50;not null;comment:版本号"`
	Description    string           `json:"description" gorm:"type:text;comment:描述"`
	
	// 状态
	Status         CapabilityStatus `json:"status" gorm:"size:20;default:enabled;comment:状态"`
	IsCore         bool             `json:"is_core" gorm:"default:false;comment:是否核心能力"`
	
	// 提供者信息
	Provider       string           `json:"provider" gorm:"size:100;comment:提供者"`
	Implementation string           `json:"implementation" gorm:"size:100;comment:实现方式"` // llm, plugin, code
	
	// 配置 (JSON存储)
	Config         JSON             `json:"config" gorm:"type:jsonb;comment:配置信息"`
	// {
	//   "model": "gpt-4",
	//   "parameters": {},
	//   "prompts": {}
	// }
	
	// 性能指标
	AvgLatency     int              `json:"avg_latency" gorm:"default:0;comment:平均延迟(ms)"`
	SuccessRate    float64          `json:"success_rate" gorm:"default:0;comment:成功率"`
	
	// 统计信息
	UsageCount     int64            `json:"usage_count" gorm:"default:0;comment:使用次数"`
	
	// 时间戳
	CreatedAt      time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Capability) TableName() string {
	return "capabilities"
}

// JSON 自定义JSON类型（复用adapter包的定义或在此重新定义）
type JSON []byte
