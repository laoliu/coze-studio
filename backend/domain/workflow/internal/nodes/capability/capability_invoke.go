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

package capability

import (
	"context"
	"encoding/json"
	"fmt"

	capabilityDomain "github.com/coze-dev/coze-studio/backend/domain/capability"
	capabilityService "github.com/coze-dev/coze-studio/backend/domain/capability/service"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
)

const (
	// NodeTypeCapabilityInvoke AI能力调用节点
	NodeTypeCapabilityInvoke entity.NodeType = "capability_invoke"
)

func init() {
	// 注册AI能力调用节点的适配器
	nodes.RegisterNodeAdaptor(NodeTypeCapabilityInvoke, func() nodes.NodeAdaptor {
		return &CapabilityInvokeAdaptor{}
	})
}

// CapabilityInvokeNode AI能力调用节点
// 用于在工作流中调用各种AI能力
type CapabilityInvokeNode struct {
	capabilityManager *capabilityService.CapabilityManager
	capabilityID      string
}

// NewCapabilityInvokeNode 创建AI能力调用节点
func NewCapabilityInvokeNode(capabilityManager *capabilityService.CapabilityManager, capabilityID string) *CapabilityInvokeNode {
	return &CapabilityInvokeNode{
		capabilityManager: capabilityManager,
		capabilityID:      capabilityID,
	}
}

// Invoke 执行AI能力调用
func (n *CapabilityInvokeNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 构建能力输入
	capInput := &capabilityDomain.CapabilityInput{
		CapabilityID: n.capabilityID,
		Parameters:   input,
		Context:      make(map[string]any),
	}

	// 如果输入中包含context字段，提取它
	if ctxData, ok := input["context"]; ok {
		if ctxMap, ok := ctxData.(map[string]any); ok {
			capInput.Context = ctxMap
		}
	}

	// 如果输入中包含parameters字段，使用它作为参数
	if params, ok := input["parameters"]; ok {
		if paramsMap, ok := params.(map[string]any); ok {
			capInput.Parameters = paramsMap
		}
	}

	// 执行能力
	output, err := n.capabilityManager.Execute(ctx, n.capabilityID, capInput)
	if err != nil {
		return nil, fmt.Errorf("failed to execute capability %s: %w", n.capabilityID, err)
	}

	// 检查是否成功
	if !output.Success {
		return nil, fmt.Errorf("capability execution failed: %s", output.Error)
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

// CapabilityInvokeAdaptor AI能力调用节点的适配器
type CapabilityInvokeAdaptor struct{}

// Adapt 将前端节点转换为后端节点模式
func (a *CapabilityInvokeAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	// 从节点配置中提取能力ID
	capabilityID, ok := n.Config["capability_id"].(string)
	if !ok || capabilityID == "" {
		return nil, fmt.Errorf("capability_id is required")
	}

	// 创建能力管理器（实际使用时应该从依赖注入获取）
	capabilityManager := capabilityService.NewCapabilityManager()

	// 创建节点实例
	node := NewCapabilityInvokeNode(capabilityManager, capabilityID)

	// 创建节点模式
	nodeSchema := &schema.NodeSchema{
		Node: node,
	}

	return nodeSchema, nil
}

// ToCallbackInput 转换为回调输入格式
func (n *CapabilityInvokeNode) ToCallbackInput(ctx context.Context, in map[string]any) (*nodes.StructuredCallbackInput, error) {
	return &nodes.StructuredCallbackInput{
		Type: "capability_invoke",
		Data: map[string]any{
			"capability_id": n.capabilityID,
			"input":         in,
		},
	}, nil
}

// ToCallbackOutput 转换为回调输出格式
func (n *CapabilityInvokeNode) ToCallbackOutput(ctx context.Context, out map[string]any) (*nodes.StructuredCallbackOutput, error) {
	return &nodes.StructuredCallbackOutput{
		Type: "capability_invoke",
		Data: out,
	}, nil
}

// 确保实现了必要的接口
var _ nodes.InvokableNode = (*CapabilityInvokeNode)(nil)
var _ nodes.CallbackInputConverted = (*CapabilityInvokeNode)(nil)
var _ nodes.CallbackOutputConverted = (*CapabilityInvokeNode)(nil)
var _ nodes.NodeAdaptor = (*CapabilityInvokeAdaptor)(nil)
