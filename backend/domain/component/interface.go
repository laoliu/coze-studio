// Copyright 2025 Coze Studio. All rights reserved.

package component

import (
	"context"
)

// ComponentExecutor 组件执行器接口
type ComponentExecutor interface {
	// GetComponentInfo 获取组件信息
	GetComponentInfo() *ComponentInfo
	
	// Execute 执行组件
	Execute(ctx context.Context, input *ComponentInput) (*ComponentOutput, error)
	
	// Validate 验证输入
	Validate(input *ComponentInput) error
	
	// GetConfig 获取配置
	GetConfig() map[string]interface{}
}

// ComponentInfo 组件信息
type ComponentInfo struct {
	ComponentID  string                 `json:"component_id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Category     string                 `json:"category"`
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`
}

// ComponentInput 组件输入
type ComponentInput struct {
	ComponentID string         `json:"component_id"`
	Parameters  map[string]any `json:"parameters"` // 输入参数
	Context     map[string]any `json:"context"`    // 上下文
}

// ComponentOutput 组件输出
type ComponentOutput struct {
	ComponentID   string         `json:"component_id"`
	Result        map[string]any `json:"result"`         // 输出结果
	Metadata      map[string]any `json:"metadata"`       // 元数据
	ExecutionTime int64          `json:"execution_time"` // 执行时间(ms)
	Success       bool           `json:"success"`        // 是否成功
	Error         string         `json:"error,omitempty"` // 错误信息
}

// MCPTool MCP工具接口
// Model Context Protocol 工具的标准接口
type MCPTool interface {
	ComponentExecutor
	
	// GetToolSchema 获取工具模式（符合MCP规范）
	GetToolSchema() *MCPToolSchema
	
	// InvokeTool 调用工具
	InvokeTool(ctx context.Context, params map[string]any) (any, error)
}

// MCPToolSchema MCP工具模式
type MCPToolSchema struct {
	Name        string                 `json:"name"`        // 工具名称
	Description string                 `json:"description"` // 工具描述
	InputSchema map[string]interface{} `json:"inputSchema"` // 输入模式
	// 符合JSON Schema规范
}

// VisualizationComponent 可视化组件
type VisualizationComponent interface {
	ComponentExecutor
	
	// Render 渲染可视化
	Render(ctx context.Context, data interface{}) (string, error) // 返回HTML或图片URL
}

// ExportComponent 导出组件
type ExportComponent interface {
	ComponentExecutor
	
	// Export 导出内容
	Export(ctx context.Context, content interface{}, format string) ([]byte, error)
}

// InteractionComponent 交互组件
type InteractionComponent interface {
	ComponentExecutor
	
	// HandleInteraction 处理交互
	HandleInteraction(ctx context.Context, event *InteractionEvent) (*InteractionResponse, error)
}

// InteractionEvent 交互事件
type InteractionEvent struct {
	EventType string         `json:"event_type"` // 事件类型
	Data      map[string]any `json:"data"`       // 事件数据
	Timestamp int64          `json:"timestamp"`  // 时间戳
}

// InteractionResponse 交互响应
type InteractionResponse struct {
	ResponseType string         `json:"response_type"` // 响应类型
	Data         map[string]any `json:"data"`          // 响应数据
	NextAction   string         `json:"next_action,omitempty"` // 下一步操作
}
