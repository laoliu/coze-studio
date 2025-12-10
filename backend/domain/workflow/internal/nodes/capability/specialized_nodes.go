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
	// NodeTypeObjectiveGenerator 学习目标生成节点
	NodeTypeObjectiveGenerator entity.NodeType = "objective_generator"
	// NodeTypeContentDiscovery 内容发现节点
	NodeTypeContentDiscovery entity.NodeType = "content_discovery"
	// NodeTypeNarrativeGenerator 故事化叙述生成节点
	NodeTypeNarrativeGenerator entity.NodeType = "narrative_generator"
	// NodeTypeQualityAssessor 质量评估节点
	NodeTypeQualityAssessor entity.NodeType = "quality_assessor"
)

func init() {
	// 注册各个专用AI能力节点
	nodes.RegisterNodeAdaptor(NodeTypeObjectiveGenerator, func() nodes.NodeAdaptor {
		return &ObjectiveGeneratorAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeContentDiscovery, func() nodes.NodeAdaptor {
		return &ContentDiscoveryAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeNarrativeGenerator, func() nodes.NodeAdaptor {
		return &NarrativeGeneratorAdaptor{}
	})
	nodes.RegisterNodeAdaptor(NodeTypeQualityAssessor, func() nodes.NodeAdaptor {
		return &QualityAssessorAdaptor{}
	})
}

// ObjectiveGeneratorNode 学习目标生成节点
type ObjectiveGeneratorNode struct {
	capabilityManager *capabilityService.CapabilityManager
}

func NewObjectiveGeneratorNode(capabilityManager *capabilityService.CapabilityManager) *ObjectiveGeneratorNode {
	return &ObjectiveGeneratorNode{
		capabilityManager: capabilityManager,
	}
}

func (n *ObjectiveGeneratorNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 解析参数
	params := &capabilityDomain.ObjectiveParams{}
	if err := mapToStruct(input, params); err != nil {
		return nil, fmt.Errorf("failed to parse objective params: %w", err)
	}

	// 调用能力
	objectives, err := n.capabilityManager.ExecuteObjectiveGeneration(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate objectives: %w", err)
	}

	return map[string]any{
		"objectives": objectives,
	}, nil
}

type ObjectiveGeneratorAdaptor struct{}

func (a *ObjectiveGeneratorAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	capabilityManager := capabilityService.NewCapabilityManager()
	node := NewObjectiveGeneratorNode(capabilityManager)
	return &schema.NodeSchema{Node: node}, nil
}

// ContentDiscoveryNode 内容发现节点
type ContentDiscoveryNode struct {
	capabilityManager *capabilityService.CapabilityManager
}

func NewContentDiscoveryNode(capabilityManager *capabilityService.CapabilityManager) *ContentDiscoveryNode {
	return &ContentDiscoveryNode{
		capabilityManager: capabilityManager,
	}
}

func (n *ContentDiscoveryNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 解析参数
	params := &capabilityDomain.DiscoveryParams{}
	if err := mapToStruct(input, params); err != nil {
		return nil, fmt.Errorf("failed to parse discovery params: %w", err)
	}

	// 调用能力
	contents, err := n.capabilityManager.ExecuteContentDiscovery(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to discover content: %w", err)
	}

	return map[string]any{
		"contents": contents,
	}, nil
}

type ContentDiscoveryAdaptor struct{}

func (a *ContentDiscoveryAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	capabilityManager := capabilityService.NewCapabilityManager()
	node := NewContentDiscoveryNode(capabilityManager)
	return &schema.NodeSchema{Node: node}, nil
}

// NarrativeGeneratorNode 故事化叙述生成节点
type NarrativeGeneratorNode struct {
	capabilityManager *capabilityService.CapabilityManager
}

func NewNarrativeGeneratorNode(capabilityManager *capabilityService.CapabilityManager) *NarrativeGeneratorNode {
	return &NarrativeGeneratorNode{
		capabilityManager: capabilityManager,
	}
}

func (n *NarrativeGeneratorNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 解析参数
	params := &capabilityDomain.NarrativeParams{}
	if err := mapToStruct(input, params); err != nil {
		return nil, fmt.Errorf("failed to parse narrative params: %w", err)
	}

	// 调用能力
	narrative, err := n.capabilityManager.ExecuteNarrativeGeneration(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate narrative: %w", err)
	}

	return structToMap(narrative)
}

type NarrativeGeneratorAdaptor struct{}

func (a *NarrativeGeneratorAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	capabilityManager := capabilityService.NewCapabilityManager()
	node := NewNarrativeGeneratorNode(capabilityManager)
	return &schema.NodeSchema{Node: node}, nil
}

// QualityAssessorNode 质量评估节点
type QualityAssessorNode struct {
	capabilityManager *capabilityService.CapabilityManager
}

func NewQualityAssessorNode(capabilityManager *capabilityService.CapabilityManager) *QualityAssessorNode {
	return &QualityAssessorNode{
		capabilityManager: capabilityManager,
	}
}

func (n *QualityAssessorNode) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 提取要评估的内容
	var content interface{}
	if contentData, ok := input["content"]; ok {
		content = contentData
	} else {
		content = input
	}

	// 调用能力
	qualityReport, err := n.capabilityManager.ExecuteQualityAssessment(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to assess quality: %w", err)
	}

	return structToMap(qualityReport)
}

type QualityAssessorAdaptor struct{}

func (a *QualityAssessorAdaptor) Adapt(ctx context.Context, n *vo.Node, opts ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	capabilityManager := capabilityService.NewCapabilityManager()
	node := NewQualityAssessorNode(capabilityManager)
	return &schema.NodeSchema{Node: node}, nil
}

// 工具函数
func mapToStruct(m interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

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
var _ nodes.InvokableNode = (*ObjectiveGeneratorNode)(nil)
var _ nodes.InvokableNode = (*ContentDiscoveryNode)(nil)
var _ nodes.InvokableNode = (*NarrativeGeneratorNode)(nil)
var _ nodes.InvokableNode = (*QualityAssessorNode)(nil)
var _ nodes.NodeAdaptor = (*ObjectiveGeneratorAdaptor)(nil)
var _ nodes.NodeAdaptor = (*ContentDiscoveryAdaptor)(nil)
var _ nodes.NodeAdaptor = (*NarrativeGeneratorAdaptor)(nil)
var _ nodes.NodeAdaptor = (*QualityAssessorAdaptor)(nil)
