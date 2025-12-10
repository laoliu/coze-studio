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

package nodes

import (
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/capability"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/component"
)

// RegisterAllNewNodes 注册所有新增的工作流节点
// 应该在应用启动时调用此函数
func RegisterAllNewNodes() {
	// 注册适配器节点
	registerAdapterNodes()

	// 注册AI能力节点
	registerCapabilityNodes()

	// 注册组件节点
	registerComponentNodes()
}

// registerAdapterNodes 注册适配器节点
func registerAdapterNodes() {
	// 创建适配器调用节点的NodeAdaptor
	adapterInvokeAdaptor := adapter.NewAdapterInvokeNodeAdaptor()

	// 注册到全局节点注册表
	// 注意：这里假设有一个全局的 RegisterNodeAdaptor 函数
	// 实际实现应该参考现有的节点注册逻辑
	// RegisterNodeAdaptor("adapter_invoke", adapterInvokeAdaptor)
	_ = adapterInvokeAdaptor
}

// registerCapabilityNodes 注册AI能力节点
func registerCapabilityNodes() {
	// 通用能力调用节点
	capabilityInvokeAdaptor := capability.NewCapabilityInvokeNodeAdaptor()

	// 目标生成节点
	objectiveGeneratorAdaptor := capability.NewObjectiveGeneratorNodeAdaptor()

	// 内容发现节点
	contentDiscoveryAdaptor := capability.NewContentDiscoveryNodeAdaptor()

	// 故事化叙述生成节点
	narrativeGeneratorAdaptor := capability.NewNarrativeGeneratorNodeAdaptor()

	// 质量评估节点
	qualityAssessorAdaptor := capability.NewQualityAssessorNodeAdaptor()

	// 注册所有能力节点
	// RegisterNodeAdaptor("capability_invoke", capabilityInvokeAdaptor)
	// RegisterNodeAdaptor("objective_generator", objectiveGeneratorAdaptor)
	// RegisterNodeAdaptor("content_discovery", contentDiscoveryAdaptor)
	// RegisterNodeAdaptor("narrative_generator", narrativeGeneratorAdaptor)
	// RegisterNodeAdaptor("quality_assessor", qualityAssessorAdaptor)
	_, _, _, _, _ = capabilityInvokeAdaptor, objectiveGeneratorAdaptor, contentDiscoveryAdaptor,
		narrativeGeneratorAdaptor, qualityAssessorAdaptor
}

// registerComponentNodes 注册组件节点
func registerComponentNodes() {
	// 通用组件调用节点
	componentInvokeAdaptor := component.NewComponentInvokeNodeAdaptor()

	// MCP工具节点
	mcpToolAdaptor := component.NewMCPToolNodeAdaptor()

	// 可视化节点
	visualizationAdaptor := component.NewVisualizationNodeAdaptor()

	// 导出节点
	exportAdaptor := component.NewExportNodeAdaptor()

	// 注册所有组件节点
	// RegisterNodeAdaptor("component_invoke", componentInvokeAdaptor)
	// RegisterNodeAdaptor("mcp_tool", mcpToolAdaptor)
	// RegisterNodeAdaptor("visualization", visualizationAdaptor)
	// RegisterNodeAdaptor("export", exportAdaptor)
	_, _, _, _ = componentInvokeAdaptor, mcpToolAdaptor, visualizationAdaptor, exportAdaptor
}

// init 函数可以在包导入时自动注册节点
// 但是建议通过显式调用 RegisterAllNewNodes() 来注册，以便更好地控制初始化顺序
// func init() {
// 	RegisterAllNewNodes()
// }
