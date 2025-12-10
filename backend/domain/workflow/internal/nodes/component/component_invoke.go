/*
 * Copyright 2025 Coze Studio. All rights reserved.
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

package component

import (
	"context"
	"fmt"

	componentDomain "github.com/coze-dev/coze-studio/backend/domain/component"
	componentService "github.com/coze-dev/coze-studio/backend/domain/component/service"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
)

const (
	// NodeTypeComponentInvoke 组件调用节点
	NodeTypeComponentInvoke entity.NodeType = "component_invoke"
	// NodeTypeMCPTool MCP工具节点
	NodeTypeMCPTool entity.NodeType = "mcp_tool"
	// NodeTypeVisualization 可视化节点
	NodeTypeVisualization entity.NodeType = "visualization"
	// NodeTypeExport 导出节点
	NodeTypeExport entity.NodeType = "export"
)

func init() {
	// 注册组件相关节点
	nodes.RegisterNodeAdaptor(NodeTypeComponentInvoke, func() nodes.NodeAdaptor {
		return &ComponentInvokeAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeMCPTool, func() nodes.NodeAdaptor {
		return &MCPToolAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeVisualization, func() nodes.NodeAdaptor {
		return &VisualizationAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeExport, func() nodes.NodeAdaptor {
		return &ExportAdaptor{}
	})
}

// ComponentInvokeNode 组件调用节点
type ComponentInvokeNode struct {
	componentManager *componentService.ComponentManager
	componentID      string
}

func NewComponentInvokeNode(componentManager *componentService.ComponentManager, componentID string) *ComponentInvokeNode {
	return &ComponentInvokeNode{
		componentManager: componentManager,
		componentID:      componentID,
	}
}

func (n *ComponentInvokeNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 构建组件输入
	compInput := &componentDomain.ComponentInput{
		ComponentID: n.componentID,
		Parameters:  input,
		Context:     make(map[string]any),
	}

	// 如果输入中包含context字段，提取它
	if ctxData, ok := input["context"]; ok {
		if ctxMap, ok := ctxData.(map[string]any); ok {
			compInput.Context = ctxMap
		}
	}

	// 如果输入中包含parameters字段，使用它作为参数
	if params, ok := input["parameters"]; ok {
		if paramsMap, ok := params.(map[string]any); ok {
			compInput.Parameters = paramsMap
		}
	}

	// 执行组件
	output, err := n.componentManager.Execute(ctx, n.componentID, compInput)
	if err != nil {
		return nil, fmt.Errorf("failed to execute component %s: %w", n.componentID, err)
	}

	// 检查是否成功
	if !output.Success {
		return nil, fmt.Errorf("component execution failed: %s", output.Error)
	}

	// 构建输出
	result := map[string]any{
		"result":         output.Result,
		"metadata":       output.Metadata,
		"execution_time": output.ExecutionTime,
		"success":        output.Success,
	}

	return result, nil
}

type ComponentInvokeAdaptor struct{}

func (a *ComponentInvokeAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	componentID, ok := n.Config["component_id"].(string)
	if !ok || componentID == "" {
		return nil, fmt.Errorf("component_id is required")
	}

	componentManager := componentService.NewComponentManager()
	node := NewComponentInvokeNode(componentManager, componentID)
	return &schema.NodeSchema{Node: node}, nil
}

// MCPToolNode MCP工具节点
type MCPToolNode struct {
	componentManager *componentService.ComponentManager
	toolName         string
}

func NewMCPToolNode(componentManager *componentService.ComponentManager, toolName string) *MCPToolNode {
	return &MCPToolNode{
		componentManager: componentManager,
		toolName:         toolName,
	}
}

func (n *MCPToolNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取工具参数
	var params map[string]any
	if paramsData, ok := input["params"]; ok {
		params = paramsData.(map[string]any)
	} else {
		params = input
	}

	// 调用MCP工具
	result, err := n.componentManager.InvokeMCPTool(ctx, n.toolName, params)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke MCP tool %s: %w", n.toolName, err)
	}

	return map[string]any{
		"result": result,
	}, nil
}

type MCPToolAdaptor struct{}

func (a *MCPToolAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	toolName, ok := n.Config["tool_name"].(string)
	if !ok || toolName == "" {
		return nil, fmt.Errorf("tool_name is required")
	}

	componentManager := componentService.NewComponentManager()
	node := NewMCPToolNode(componentManager, toolName)
	return &schema.NodeSchema{Node: node}, nil
}

// VisualizationNode 可视化节点
type VisualizationNode struct {
	componentManager *componentService.ComponentManager
	componentID      string
}

func NewVisualizationNode(componentManager *componentService.ComponentManager, componentID string) *VisualizationNode {
	return &VisualizationNode{
		componentManager: componentManager,
		componentID:      componentID,
	}
}

func (n *VisualizationNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取要可视化的数据
	var data interface{}
	if dataField, ok := input["data"]; ok {
		data = dataField
	} else {
		data = input
	}

	// 渲染可视化
	htmlOrURL, err := n.componentManager.RenderVisualization(ctx, n.componentID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render visualization: %w", err)
	}

	return map[string]any{
		"visualization": htmlOrURL,
	}, nil
}

type VisualizationAdaptor struct{}

func (a *VisualizationAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	componentID, ok := n.Config["component_id"].(string)
	if !ok || componentID == "" {
		return nil, fmt.Errorf("component_id is required")
	}

	componentManager := componentService.NewComponentManager()
	node := NewVisualizationNode(componentManager, componentID)
	return &schema.NodeSchema{Node: node}, nil
}

// ExportNode 导出节点
type ExportNode struct {
	componentManager *componentService.ComponentManager
	componentID      string
	format           string
}

func NewExportNode(componentManager *componentService.ComponentManager, componentID, format string) *ExportNode {
	return &ExportNode{
		componentManager: componentManager,
		componentID:      componentID,
		format:           format,
	}
}

func (n *ExportNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取要导出的内容
	var content interface{}
	if contentField, ok := input["content"]; ok {
		content = contentField
	} else {
		content = input
	}

	// 导出内容
	data, err := n.componentManager.ExportContent(ctx, n.componentID, content, n.format)
	if err != nil {
		return nil, fmt.Errorf("failed to export content: %w", err)
	}

	return map[string]any{
		"exported_data": data,
		"format":        n.format,
	}, nil
}

type ExportAdaptor struct{}

func (a *ExportAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	componentID, ok := n.Config["component_id"].(string)
	if !ok || componentID == "" {
		return nil, fmt.Errorf("component_id is required")
	}

	format, ok := n.Config["format"].(string)
	if !ok || format == "" {
		format = "pdf" // 默认格式
	}

	componentManager := componentService.NewComponentManager()
	node := NewExportNode(componentManager, componentID, format)
	return &schema.NodeSchema{Node: node}, nil
}

// 确保实现了必要的接口
var _ nodes.InvokableNode = (*ComponentInvokeNode)(nil)
var _ nodes.InvokableNode = (*MCPToolNode)(nil)
var _ nodes.InvokableNode = (*VisualizationNode)(nil)
var _ nodes.InvokableNode = (*ExportNode)(nil)
var _ nodes.NodeAdaptor = (*ComponentInvokeAdaptor)(nil)
var _ nodes.NodeAdaptor = (*MCPToolAdaptor)(nil)
var _ nodes.NodeAdaptor = (*VisualizationAdaptor)(nil)
var _ nodes.NodeAdaptor = (*ExportAdaptor)(nil)
