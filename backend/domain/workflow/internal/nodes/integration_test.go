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

package nodes_test

import (
	"context"
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/examples"
	adaptermanager "github.com/coze-dev/coze-studio/backend/domain/adapter/service"
	workflowadapter "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/adapter"
)

// TestAdapterInvokeNode 测试适配器调用节点
func TestAdapterInvokeNode(t *testing.T) {
	ctx := context.Background()

	// 1. 注册K12适配器
	k12Adapter := examples.NewK12Adapter()
	adapterManager := adaptermanager.NewAdapterManager()
	if err := adapterManager.Register(k12Adapter); err != nil {
		t.Fatalf("Failed to register K12 adapter: %v", err)
	}

	// 2. 创建适配器调用节点
	nodeConfig := map[string]interface{}{
		"adapter_id": "k12_adapter",
		"method":     "parse_request",
		"input_mapping": map[string]interface{}{
			"user_input": "{{workflow.input.query}}",
		},
		"output_mapping": map[string]interface{}{
			"parsed_data": "{{node.output}}",
		},
	}

	node, err := workflowadapter.NewAdapterInvokeNode(ctx, nodeConfig)
	if err != nil {
		t.Fatalf("Failed to create adapter invoke node: %v", err)
	}

	// 3. 准备输入数据
	input := map[string]interface{}{
		"user_input": "我想学习七年级数学的二次函数",
	}

	// 4. 执行节点
	output, err := node.Invoke(ctx, input)
	if err != nil {
		t.Fatalf("Failed to invoke adapter node: %v", err)
	}

	// 5. 验证输出
	t.Logf("Node output: %+v", output)

	// 验证输出包含期望的字段
	if output == nil {
		t.Fatal("Node output is nil")
	}

	// 根据实际的适配器实现来验证
	// 例如：检查是否包含 grade, subject 等字段
}

// TestCapabilityInvokeNode 测试AI能力调用节点
// 这个测试需要实际的AI能力实现和模型配置
// 这里提供一个框架示例
func TestCapabilityInvokeNode(t *testing.T) {
	t.Skip("Skipping capability invoke node test - requires AI model setup")

	ctx := context.Background()

	// 准备节点配置
	nodeConfig := map[string]interface{}{
		"capability_id":   "objective_generator",
		"capability_type": "objective_generation",
		"input": map[string]interface{}{
			"domain":        "K12数学",
			"grade":         "七年级",
			"subject":       "二次函数",
			"student_level": "中等",
		},
		"model_config": map[string]interface{}{
			"model":       "gpt-4",
			"temperature": 0.7,
			"max_tokens":  2000,
		},
	}

	// 创建并执行节点
	// node, err := capability.NewCapabilityInvokeNode(ctx, nodeConfig)
	// ...
	_ = nodeConfig
}

// TestComponentInvokeNode 测试组件调用节点
func TestComponentInvokeNode(t *testing.T) {
	t.Skip("Skipping component invoke node test - requires component setup")

	ctx := context.Background()

	// 准备节点配置
	nodeConfig := map[string]interface{}{
		"component_id":   "export_pdf",
		"component_type": "export",
		"input": map[string]interface{}{
			"content": "这是要导出的内容",
			"title":   "学习材料",
		},
		"config": map[string]interface{}{
			"format": "PDF",
			"layout": "A4",
		},
	}

	// 创建并执行节点
	// node, err := component.NewComponentInvokeNode(ctx, nodeConfig)
	// ...
	_ = nodeConfig
}

// TestWorkflowIntegration 测试完整的工作流集成
// 这个测试展示了如何将多个节点串联起来形成一个完整的工作流
func TestWorkflowIntegration(t *testing.T) {
	t.Skip("Skipping workflow integration test - requires full system setup")

	ctx := context.Background()

	// 工作流场景：K12教学内容生成
	// 1. 适配器节点：解析用户请求（识别年级、科目等）
	// 2. 能力节点：生成学习目标
	// 3. 能力节点：内容发现
	// 4. 能力节点：故事化叙述生成
	// 5. 能力节点：质量评估
	// 6. 组件节点：可视化生成
	// 7. 组件节点：导出为PDF

	// 准备工作流输入
	workflowInput := map[string]interface{}{
		"query": "我想学习七年级数学的二次函数",
	}

	// 执行工作流各步骤
	// ...

	_ = ctx
	_ = workflowInput
}

// BenchmarkAdapterInvokeNode 适配器节点性能测试
func BenchmarkAdapterInvokeNode(b *testing.B) {
	ctx := context.Background()

	// 准备测试环境
	k12Adapter := examples.NewK12Adapter()
	adapterManager := adaptermanager.NewAdapterManager()
	if err := adapterManager.Register(k12Adapter); err != nil {
		b.Fatalf("Failed to register K12 adapter: %v", err)
	}

	nodeConfig := map[string]interface{}{
		"adapter_id": "k12_adapter",
		"method":     "parse_request",
	}

	node, err := workflowadapter.NewAdapterInvokeNode(ctx, nodeConfig)
	if err != nil {
		b.Fatalf("Failed to create adapter invoke node: %v", err)
	}

	input := map[string]interface{}{
		"user_input": "我想学习七年级数学的二次函数",
	}

	// 性能测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := node.Invoke(ctx, input)
		if err != nil {
			b.Fatalf("Failed to invoke node: %v", err)
		}
	}
}
