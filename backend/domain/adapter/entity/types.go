// Copyright 2025 Coze Studio. All rights reserved.

package entity

// AdapterInfo 适配器信息（轻量级版本，用于接口返回）
type AdapterInfo struct {
AdapterID    string               `json:"adapter_id"`
Name         string               `json:"name"`
DisplayName  string               `json:"display_name"`
Version      string               `json:"version"`
Description  string               `json:"description"`
Type         AdapterType          `json:"type"`
Capabilities *AdapterCapabilities `json:"capabilities"`
}

// UserInput 用户输入
type UserInput struct {
Topic        string         `json:"topic"`         // 主题
Domain       string         `json:"domain"`        // 领域
Grade        string         `json:"grade"`         // 年级
Duration     int            `json:"duration"`      // 时长（分钟）
ActivityType ActivityType   `json:"activity_type"` // 活动类型
Language     string         `json:"language"`      // 语言
ExtraParams  map[string]any `json:"extra_params"`  // 额外参数
}

// RequestContext 请求上下文
type RequestContext struct {
// 基础信息
UserInput      *UserInput        `json:"user_input"`
Topic          string            `json:"topic"`
Domain         string            `json:"domain"`
Grade          string            `json:"grade"`
Duration       int               `json:"duration"`
ActivityType   ActivityType      `json:"activity_type"`
Language       string            `json:"language"`

// 学习者信息
LearnerProfile *LearnerProfile   `json:"learner_profile,omitempty"`

// 领域特定信息
DomainSpecific map[string]any    `json:"domain_specific"`

// 约束条件
Constraints    *Constraints      `json:"constraints,omitempty"`

// 偏好设置
Preferences    map[string]any    `json:"preferences,omitempty"`

// 元数据
Metadata       map[string]any    `json:"metadata,omitempty"`
}

// LearnerProfile 学习者画像
type LearnerProfile struct {
Age            int      `json:"age"`
Grade          string   `json:"grade"`
KnowledgeLevel string   `json:"knowledge_level"` // beginner, intermediate, advanced
Interests      []string `json:"interests"`
LearningStyle  string   `json:"learning_style"`  // visual, auditory, kinesthetic
SpecialNeeds   []string `json:"special_needs,omitempty"`
}

// Constraints 约束条件
type Constraints struct {
MaxDuration     int      `json:"max_duration,omitempty"`      // 最大时长
MinDuration     int      `json:"min_duration,omitempty"`      // 最小时长
RequiredTopics  []string `json:"required_topics,omitempty"`   // 必须包含的主题
ExcludedTopics  []string `json:"excluded_topics,omitempty"`   // 排除的主题
DifficultyLevel string   `json:"difficulty_level,omitempty"`  // 难度级别
SafetyLevel     string   `json:"safety_level,omitempty"`      // 安全级别
}

// LearningObjective 学习目标
type LearningObjective struct {
ID           string `json:"id"`
Objective    string `json:"objective"`                // 目标描述
Level        string `json:"level"`                    // 布鲁姆分类法级别
Category     string `json:"category"`                 // 类别：知识、技能、态度
IsMeasurable bool   `json:"is_measurable"`            // 是否可测量
Assessment   string `json:"assessment,omitempty"`     // 评估方式
}

// Content 内容
type Content struct {
ID          string         `json:"id"`
Type        string         `json:"type"`                  // text, image, video, audio
Title       string         `json:"title"`
Description string         `json:"description"`
Source      string         `json:"source"`                // 来源
URL         string         `json:"url,omitempty"`
Metadata    map[string]any `json:"metadata,omitempty"`
Relevance   float64        `json:"relevance"`             // 相关性评分
}

// WorkflowConfig 工作流配置
type WorkflowConfig struct {
WorkflowID string         `json:"workflow_id"`
Template   string         `json:"template"`              // 工作流模板
Nodes      []*NodeConfig  `json:"nodes"`                 // 节点配置
Variables  map[string]any `json:"variables"`             // 变量
Timeout    int            `json:"timeout"`               // 超时时间（秒）
}

// NodeConfig 节点配置
type NodeConfig struct {
NodeID       string         `json:"node_id"`
NodeType     string         `json:"node_type"`            // LLM, Plugin, Code等
Config       map[string]any `json:"config"`
Dependencies []string       `json:"dependencies"`         // 依赖的节点ID
}

// WorkflowResult 工作流执行结果
type WorkflowResult struct {
WorkflowID    string         `json:"workflow_id"`
Status        string         `json:"status"`               // success, failed, timeout
Output        map[string]any `json:"output"`               // 输出数据
ExecutionTime int64          `json:"execution_time"`       // 执行时间（毫秒）
ErrorMessage  string         `json:"error_message,omitempty"`
}

// QualityReport 质量报告
type QualityReport struct {
OverallScore float64             `json:"overall_score"`      // 总体评分
Dimensions   map[string]float64  `json:"dimensions"`         // 各维度评分
Issues       []*QualityIssue     `json:"issues"`             // 质量问题
Suggestions  []string            `json:"suggestions"`        // 改进建议
IsApproved   bool                `json:"is_approved"`        // 是否通过
}

// QualityIssue 质量问题
type QualityIssue struct {
Severity    string `json:"severity"`                      // critical, major, minor
Category    string `json:"category"`                      // 类别
Description string `json:"description"`                   // 描述
Location    string `json:"location,omitempty"`            // 位置
Suggestion  string `json:"suggestion,omitempty"`          // 建议
}
