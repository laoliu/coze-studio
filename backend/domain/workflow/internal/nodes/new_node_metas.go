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
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
)

// 新增节点类型的元数据定义
var (
	// AdapterInvokeNodeMeta 适配器调用节点元数据
	AdapterInvokeNodeMeta = &entity.NodeTypeMeta{
		ID:              1001,
		Key:             "adapter_invoke",
		DisplayKey:      "adapter_invoke",
		Name:            "适配器调用",
		Category:        "adapter",
		Color:           "#8B5CF6",
		Desc:            "调用领域适配器的方法，实现领域特定的逻辑处理",
		IconURI:         "/icons/adapter.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Adapter Invoke",
		EnUSDescription: "Invoke domain adapter methods for domain-specific logic",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  false,
		},
	}

	// CapabilityInvokeNodeMeta AI能力调用节点元数据
	CapabilityInvokeNodeMeta = &entity.NodeTypeMeta{
		ID:              1002,
		Key:             "capability_invoke",
		DisplayKey:      "capability_invoke",
		Name:            "AI能力调用",
		Category:        "ai_capability",
		Color:           "#10B981",
		Desc:            "调用AI能力，如意图识别、目标生成、内容发现等",
		IconURI:         "/icons/ai_capability.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "AI Capability Invoke",
		EnUSDescription: "Invoke AI capabilities like intent recognition, objective generation, etc",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 60000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  true,
		},
	}

	// ObjectiveGeneratorNodeMeta 学习目标生成节点元数据
	ObjectiveGeneratorNodeMeta = &entity.NodeTypeMeta{
		ID:              1003,
		Key:             "objective_generator",
		DisplayKey:      "objective_generator",
		Name:            "学习目标生成",
		Category:        "ai_capability",
		Color:           "#10B981",
		Desc:            "根据输入生成SMART学习目标",
		IconURI:         "/icons/objective.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Objective Generator",
		EnUSDescription: "Generate SMART learning objectives based on input",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  true,
		},
	}

	// ContentDiscoveryNodeMeta 内容发现节点元数据
	ContentDiscoveryNodeMeta = &entity.NodeTypeMeta{
		ID:              1004,
		Key:             "content_discovery",
		DisplayKey:      "content_discovery",
		Name:            "内容发现",
		Category:        "ai_capability",
		Color:           "#10B981",
		Desc:            "从知识库或外部源发现相关内容",
		IconURI:         "/icons/content.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Content Discovery",
		EnUSDescription: "Discover relevant content from knowledge base or external sources",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  true,
		},
	}

	// NarrativeGeneratorNodeMeta 故事化叙述生成节点元数据
	NarrativeGeneratorNodeMeta = &entity.NodeTypeMeta{
		ID:              1005,
		Key:             "narrative_generator",
		DisplayKey:      "narrative_generator",
		Name:            "故事化叙述生成",
		Category:        "ai_capability",
		Color:           "#10B981",
		Desc:            "生成引人入胜的故事化叙述内容",
		IconURI:         "/icons/narrative.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Narrative Generator",
		EnUSDescription: "Generate engaging narrative content",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 60000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  true,
		},
	}

	// QualityAssessorNodeMeta 质量评估节点元数据
	QualityAssessorNodeMeta = &entity.NodeTypeMeta{
		ID:              1006,
		Key:             "quality_assessor",
		DisplayKey:      "quality_assessor",
		Name:            "质量评估",
		Category:        "ai_capability",
		Color:           "#10B981",
		Desc:            "评估内容质量，生成详细的质量报告",
		IconURI:         "/icons/quality.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Quality Assessor",
		EnUSDescription: "Assess content quality and generate detailed quality reports",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  true,
		},
	}

	// ComponentInvokeNodeMeta 组件调用节点元数据
	ComponentInvokeNodeMeta = &entity.NodeTypeMeta{
		ID:              1007,
		Key:             "component_invoke",
		DisplayKey:      "component_invoke",
		Name:            "组件调用",
		Category:        "component",
		Color:           "#F59E0B",
		Desc:            "调用可复用的功能组件",
		IconURI:         "/icons/component.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Component Invoke",
		EnUSDescription: "Invoke reusable functional components",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  false,
		},
	}

	// MCPToolNodeMeta MCP工具节点元数据
	MCPToolNodeMeta = &entity.NodeTypeMeta{
		ID:              1008,
		Key:             "mcp_tool",
		DisplayKey:      "mcp_tool",
		Name:            "MCP工具",
		Category:        "component",
		Color:           "#F59E0B",
		Desc:            "调用Model Context Protocol工具",
		IconURI:         "/icons/mcp.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "MCP Tool",
		EnUSDescription: "Invoke Model Context Protocol tools",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  false,
		},
	}

	// VisualizationNodeMeta 可视化节点元数据
	VisualizationNodeMeta = &entity.NodeTypeMeta{
		ID:              1009,
		Key:             "visualization",
		DisplayKey:      "visualization",
		Name:            "可视化",
		Category:        "component",
		Color:           "#F59E0B",
		Desc:            "将数据渲染为可视化图表或界面",
		IconURI:         "/icons/visualization.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Visualization",
		EnUSDescription: "Render data as visualizations or interfaces",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 15000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  false,
		},
	}

	// ExportNodeMeta 导出节点元数据
	ExportNodeMeta = &entity.NodeTypeMeta{
		ID:              1010,
		Key:             "export",
		DisplayKey:      "export",
		Name:            "导出",
		Category:        "component",
		Color:           "#F59E0B",
		Desc:            "将内容导出为指定格式（PDF、PPT、Word等）",
		IconURI:         "/icons/export.svg",
		SupportBatch:    false,
		Disabled:        false,
		EnUSName:        "Export",
		EnUSDescription: "Export content to specified formats (PDF, PPT, Word, etc.)",
		ExecutableMeta: entity.ExecutableMeta{
			IsComposite:      false,
			DefaultTimeoutMS: 30000,
			PreFillZero:      false,
			PostFillNil:      false,
			MayUseChatModel:  false,
		},
	}
)

// NewNodeMetas 新增节点元数据列表
var NewNodeMetas = []*entity.NodeTypeMeta{
	AdapterInvokeNodeMeta,
	CapabilityInvokeNodeMeta,
	ObjectiveGeneratorNodeMeta,
	ContentDiscoveryNodeMeta,
	NarrativeGeneratorNodeMeta,
	QualityAssessorNodeMeta,
	ComponentInvokeNodeMeta,
	MCPToolNodeMeta,
	VisualizationNodeMeta,
	ExportNodeMeta,
}

// NewCategories 新增节点分类
var NewCategories = []entity.Category{
	{
		Key:      "adapter",
		Name:     "领域适配器",
		EnUSName: "Domain Adapter",
	},
	{
		Key:      "ai_capability",
		Name:     "AI能力",
		EnUSName: "AI Capability",
	},
	{
		Key:      "component",
		Name:     "组件",
		EnUSName: "Component",
	},
}

// RegisterNewNodeMetas 注册新节点元数据
// 应该在系统初始化时调用此函数
func RegisterNewNodeMetas() {
	// 这里可以将新节点元数据添加到全局的NodeTypeMetas中
	// 具体实现取决于现有系统的初始化逻辑
}
