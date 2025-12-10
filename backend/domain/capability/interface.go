// Copyright 2025 Coze Studio. All rights reserved.

package capability

import (
	"context"
)

// AICapability AI能力接口
// 所有AI能力实现必须实现此接口
type AICapability interface {
	// GetCapabilityInfo 获取能力信息
	GetCapabilityInfo() *CapabilityInfo
	
	// Execute 执行能力
	Execute(ctx context.Context, input *CapabilityInput) (*CapabilityOutput, error)
	
	// Validate 验证输入
	Validate(input *CapabilityInput) error
	
	// GetConfig 获取配置
	GetConfig() map[string]interface{}
}

// CapabilityInfo 能力信息
type CapabilityInfo struct {
	CapabilityID string            `json:"capability_id"`
	Name         string            `json:"name"`
	DisplayName  string            `json:"display_name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Type         string            `json:"type"`
	Provider     string            `json:"provider"`
	InputSchema  map[string]interface{} `json:"input_schema"`  // 输入模式
	OutputSchema map[string]interface{} `json:"output_schema"` // 输出模式
}

// CapabilityInput 能力输入
type CapabilityInput struct {
	CapabilityID string         `json:"capability_id"`
	Parameters   map[string]any `json:"parameters"` // 输入参数
	Context      map[string]any `json:"context"`    // 上下文信息
}

// CapabilityOutput 能力输出
type CapabilityOutput struct {
	CapabilityID  string         `json:"capability_id"`
	Result        map[string]any `json:"result"`         // 输出结果
	Metadata      map[string]any `json:"metadata"`       // 元数据
	ExecutionTime int64          `json:"execution_time"` // 执行时间(ms)
	Success       bool           `json:"success"`        // 是否成功
	Error         string         `json:"error,omitempty"` // 错误信息
}

// IntentRecognizer 意图识别能力
type IntentRecognizer interface {
	AICapability
	
	// RecognizeIntent 识别用户意图
	RecognizeIntent(ctx context.Context, input string) (*Intent, error)
}

// Intent 意图
type Intent struct {
	IntentType   string         `json:"intent_type"`   // 意图类型
	Confidence   float64        `json:"confidence"`    // 置信度
	Entities     map[string]any `json:"entities"`      // 实体
	Domain       string         `json:"domain"`        // 领域
	ActivityType string         `json:"activity_type"` // 活动类型
}

// ObjectiveGenerator 学习目标生成能力
type ObjectiveGenerator interface {
	AICapability
	
	// GenerateObjectives 生成学习目标
	GenerateObjectives(ctx context.Context, params *ObjectiveParams) ([]*Objective, error)
}

// ObjectiveParams 目标生成参数
type ObjectiveParams struct {
	Topic    string `json:"topic"`    // 主题
	Grade    string `json:"grade"`    // 年级
	Domain   string `json:"domain"`   // 领域
	Duration int    `json:"duration"` // 时长
}

// Objective 学习目标
type Objective struct {
	ID          string  `json:"id"`
	Objective   string  `json:"objective"`    // 目标描述
	Level       string  `json:"level"`        // 布鲁姆级别
	Category    string  `json:"category"`     // 类别
	IsMeasurable bool   `json:"is_measurable"` // 可测量性
	Assessment  string  `json:"assessment"`   // 评估方式
}

// ContentDiscovery 内容发现能力
type ContentDiscovery interface {
	AICapability
	
	// DiscoverContent 发现内容
	DiscoverContent(ctx context.Context, params *DiscoveryParams) ([]*Content, error)
}

// DiscoveryParams 内容发现参数
type DiscoveryParams struct {
	Query       string   `json:"query"`       // 查询
	Domain      string   `json:"domain"`      // 领域
	Sources     []string `json:"sources"`     // 内容源
	MaxResults  int      `json:"max_results"` // 最大结果数
}

// Content 内容
type Content struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`        // 类型
	Title       string         `json:"title"`       // 标题
	Description string         `json:"description"` // 描述
	Source      string         `json:"source"`      // 来源
	URL         string         `json:"url"`         // URL
	Metadata    map[string]any `json:"metadata"`    // 元数据
	Relevance   float64        `json:"relevance"`   // 相关性
}

// NarrativeGenerator 故事化叙述生成能力
type NarrativeGenerator interface {
	AICapability
	
	// GenerateNarrative 生成叙述
	GenerateNarrative(ctx context.Context, params *NarrativeParams) (*Narrative, error)
}

// NarrativeParams 叙述生成参数
type NarrativeParams struct {
	Topic     string         `json:"topic"`     // 主题
	Content   []*Content     `json:"content"`   // 内容
	Style     string         `json:"style"`     // 风格
	Metadata  map[string]any `json:"metadata"`  // 元数据
}

// Narrative 叙述
type Narrative struct {
	Title       string   `json:"title"`       // 标题
	Introduction string  `json:"introduction"` // 引言
	Body        []string `json:"body"`        // 正文
	Conclusion  string   `json:"conclusion"`  // 结论
	Style       string   `json:"style"`       // 风格
}

// QualityAssessor 质量评估能力
type QualityAssessor interface {
	AICapability
	
	// AssessQuality 评估质量
	AssessQuality(ctx context.Context, content interface{}) (*QualityReport, error)
}

// QualityReport 质量报告
type QualityReport struct {
	OverallScore float64              `json:"overall_score"` // 总分
	Dimensions   map[string]float64   `json:"dimensions"`    // 维度评分
	Issues       []*QualityIssue      `json:"issues"`        // 问题
	Suggestions  []string             `json:"suggestions"`   // 建议
	IsApproved   bool                 `json:"is_approved"`   // 是否通过
}

// QualityIssue 质量问题
type QualityIssue struct {
	Severity    string `json:"severity"`    // 严重程度
	Category    string `json:"category"`    // 类别
	Description string `json:"description"` // 描述
	Location    string `json:"location"`    // 位置
	Suggestion  string `json:"suggestion"`  // 建议
}

// MultimodalGenerator 多模态生成能力
type MultimodalGenerator interface {
	AICapability
	
	// GenerateMultimodal 生成多模态内容
	GenerateMultimodal(ctx context.Context, params *MultimodalParams) (*MultimodalContent, error)
}

// MultimodalParams 多模态生成参数
type MultimodalParams struct {
	Type        string         `json:"type"`        // 类型: image, video, audio
	Description string         `json:"description"` // 描述
	Style       string         `json:"style"`       // 风格
	Parameters  map[string]any `json:"parameters"`  // 参数
}

// MultimodalContent 多模态内容
type MultimodalContent struct {
	Type     string         `json:"type"`     // 类型
	URL      string         `json:"url"`      // URL
	Metadata map[string]any `json:"metadata"` // 元数据
}
