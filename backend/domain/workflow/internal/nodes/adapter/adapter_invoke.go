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

package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	adapterDomain "github.com/coze-dev/coze-studio/backend/domain/adapter"
	adapterEntity "github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	adapterService "github.com/coze-dev/coze-studio/backend/domain/adapter/service"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
)

const (
	// NodeTypeAdapterInvoke 适配器调用节点
	NodeTypeAdapterInvoke entity.NodeType = "adapter_invoke"
)

func init() {
	// 注册适配器调用节点的适配器
	nodes.RegisterNodeAdaptor(NodeTypeAdapterInvoke, func() nodes.NodeAdaptor {
		return &AdapterInvokeAdaptor{}
	})
}

// AdapterInvokeNode 适配器调用节点
// 用于在工作流中调用领域适配器的各种方法
type AdapterInvokeNode struct {
	adapterManager *adapterService.AdapterManager
	adapterID      string
	method         string // 调用的方法：parse_request, generate_objectives, discover_content等
}

// NewAdapterInvokeNode 创建适配器调用节点
func NewAdapterInvokeNode(adapterManager *adapterService.AdapterManager, adapterID, method string) *AdapterInvokeNode {
	return &AdapterInvokeNode{
		adapterManager: adapterManager,
		adapterID:      adapterID,
		method:         method,
	}
}

// Invoke 执行适配器方法调用
func (n *AdapterInvokeNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	switch n.method {
	case "parse_request":
		return n.parseRequest(ctx, input)
	case "generate_objectives":
		return n.generateObjectives(ctx, input)
	case "discover_content":
		return n.discoverContent(ctx, input)
	case "customize_workflow":
		return n.customizeWorkflow(ctx, input)
	case "format_output":
		return n.formatOutput(ctx, input)
	case "validate_output":
		return n.validateOutput(ctx, input)
	default:
		return nil, fmt.Errorf("unsupported adapter method: %s", n.method)
	}
}

// parseRequest 解析用户请求
func (n *AdapterInvokeNode) parseRequest(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 从输入中提取用户输入
	userInput := &adapterEntity.UserInput{}
	if err := mapToStruct(input, userInput); err != nil {
		return nil, fmt.Errorf("failed to parse user input: %w", err)
	}

	// 调用适配器管理器解析请求
	reqCtx, err := n.adapterManager.ParseRequest(ctx, n.adapterID, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	// 将结果转换为map返回
	return structToMap(reqCtx)
}

// generateObjectives 生成学习目标
func (n *AdapterInvokeNode) generateObjectives(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 从输入中提取请求上下文
	reqCtx := &adapterEntity.RequestContext{}
	if err := mapToStruct(input, reqCtx); err != nil {
		return nil, fmt.Errorf("failed to parse request context: %w", err)
	}

	// 调用适配器生成学习目标
	objectives, err := n.adapterManager.GenerateLearningObjectives(ctx, n.adapterID, reqCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate objectives: %w", err)
	}

	return map[string]any{
		"objectives": objectives,
	}, nil
}

// discoverContent 发现内容
func (n *AdapterInvokeNode) discoverContent(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取请求上下文
	reqCtx := &adapterEntity.RequestContext{}
	if reqCtxData, ok := input["request_context"]; ok {
		if err := mapToStruct(reqCtxData, reqCtx); err != nil {
			return nil, fmt.Errorf("failed to parse request context: %w", err)
		}
	}

	// 提取学习目标
	var objectives []*adapterEntity.LearningObjective
	if objData, ok := input["objectives"]; ok {
		if err := mapToStruct(objData, &objectives); err != nil {
			return nil, fmt.Errorf("failed to parse objectives: %w", err)
		}
	}

	// 调用适配器发现内容
	contents, err := n.adapterManager.DiscoverContent(ctx, n.adapterID, reqCtx, objectives)
	if err != nil {
		return nil, fmt.Errorf("failed to discover content: %w", err)
	}

	return map[string]any{
		"contents": contents,
	}, nil
}

// customizeWorkflow 定制工作流
func (n *AdapterInvokeNode) customizeWorkflow(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取请求上下文
	reqCtx := &adapterEntity.RequestContext{}
	if err := mapToStruct(input, reqCtx); err != nil {
		return nil, fmt.Errorf("failed to parse request context: %w", err)
	}

	// 调用适配器定制工作流
	workflowConfig, err := n.adapterManager.CustomizeWorkflow(ctx, n.adapterID, reqCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to customize workflow: %w", err)
	}

	return structToMap(workflowConfig)
}

// formatOutput 格式化输出
func (n *AdapterInvokeNode) formatOutput(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取工作流结果
	result := &adapterEntity.WorkflowResult{}
	if err := mapToStruct(input, result); err != nil {
		return nil, fmt.Errorf("failed to parse workflow result: %w", err)
	}

	// 调用适配器格式化输出
	output, err := n.adapterManager.FormatOutput(ctx, n.adapterID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to format output: %w", err)
	}

	return map[string]any{
		"formatted_output": output,
	}, nil
}

// validateOutput 验证输出
func (n *AdapterInvokeNode) validateOutput(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取输出数据
	var output interface{}
	if outputData, ok := input["output"]; ok {
		output = outputData
	} else {
		output = input
	}

	// 调用适配器验证输出
	qualityReport, err := n.adapterManager.ValidateOutput(ctx, n.adapterID, output)
	if err != nil {
		return nil, fmt.Errorf("failed to validate output: %w", err)
	}

	return structToMap(qualityReport)
}

// AdapterInvokeAdaptor 适配器调用节点的适配器
type AdapterInvokeAdaptor struct{}

// Adapt 将前端节点转换为后端节点模式
func (a *AdapterInvokeAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	// 从节点配置中提取参数
	adapterID, ok := n.Config["adapter_id"].(string)
	if !ok || adapterID == "" {
		return nil, fmt.Errorf("adapter_id is required")
	}

	method, ok := n.Config["method"].(string)
	if !ok || method == "" {
		return nil, fmt.Errorf("method is required")
	}

	// 创建适配器管理器（实际使用时应该从依赖注入获取）
	adapterManager := adapterService.NewAdapterManager()

	// 创建节点实例
	node := NewAdapterInvokeNode(adapterManager, adapterID, method)

	// 创建节点模式
	nodeSchema := &schema.NodeSchema{
		Node: node,
		// 可以添加更多配置
	}

	return nodeSchema, nil
}

// ToCallbackInput 转换为回调输入格式（用于UI显示）
func (n *AdapterInvokeNode) ToCallbackInput(ctx context.Context, in map[string]any) (*nodes.StructuredCallbackInput, error) {
	return &nodes.StructuredCallbackInput{
		Type: "adapter_invoke",
		Data: in,
	}, nil
}

// ToCallbackOutput 转换为回调输出格式（用于UI显示）
func (n *AdapterInvokeNode) ToCallbackOutput(ctx context.Context, out map[string]any) (*nodes.StructuredCallbackOutput, error) {
	return &nodes.StructuredCallbackOutput{
		Type: "adapter_invoke",
		Data: out,
	}, nil
}

// 工具函数：map转struct
func mapToStruct(m interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

// 工具函数：struct转map
func structToMap(s interface{}) (map[string]any, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// 确保实现了必要的接口
var _ nodes.InvokableNode = (*AdapterInvokeNode)(nil)
var _ nodes.CallbackInputConverted = (*AdapterInvokeNode)(nil)
var _ nodes.CallbackOutputConverted = (*AdapterInvokeNode)(nil)
var _ nodes.NodeAdaptor = (*AdapterInvokeAdaptor)(nil)
