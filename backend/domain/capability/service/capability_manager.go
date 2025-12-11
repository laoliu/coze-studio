// Copyright 2025 Coze Studio. All rights reserved.

package service

import (
	"context"
	"errors"
	"sync"

	"github.com/coze-dev/coze-studio/backend/domain/capability"
)

var (
	ErrCapabilityNotFound      = errors.New("capability not found")
	ErrCapabilityAlreadyExists = errors.New("capability already exists")
	ErrInvalidCapability       = errors.New("invalid capability")
)

// CapabilityRegistry AI能力注册表
type CapabilityRegistry struct {
	capabilities map[string]capability.AICapability
	metadata     map[string]*capability.CapabilityInfo
	mu           sync.RWMutex
}

// NewCapabilityRegistry 创建能力注册表
func NewCapabilityRegistry() *CapabilityRegistry {
	return &CapabilityRegistry{
		capabilities: make(map[string]capability.AICapability),
		metadata:     make(map[string]*capability.CapabilityInfo),
	}
}

// Register 注册能力
func (r *CapabilityRegistry) Register(cap capability.AICapability) error {
	if cap == nil {
		return ErrInvalidCapability
	}

	info := cap.GetCapabilityInfo()
	if info == nil || info.CapabilityID == "" {
		return ErrInvalidCapability
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.capabilities[info.CapabilityID]; exists {
		return ErrCapabilityAlreadyExists
	}

	r.capabilities[info.CapabilityID] = cap
	r.metadata[info.CapabilityID] = info

	return nil
}

// Unregister 注销能力
func (r *CapabilityRegistry) Unregister(capabilityID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.capabilities[capabilityID]; !exists {
		return ErrCapabilityNotFound
	}

	delete(r.capabilities, capabilityID)
	delete(r.metadata, capabilityID)

	return nil
}

// Get 获取能力
func (r *CapabilityRegistry) Get(capabilityID string) (capability.AICapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cap, exists := r.capabilities[capabilityID]
	if !exists {
		return nil, ErrCapabilityNotFound
	}

	return cap, nil
}

// GetInfo 获取能力信息
func (r *CapabilityRegistry) GetInfo(capabilityID string) (*capability.CapabilityInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.metadata[capabilityID]
	if !exists {
		return nil, ErrCapabilityNotFound
	}

	return info, nil
}

// List 列出所有已注册的能力
func (r *CapabilityRegistry) List() []*capability.CapabilityInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*capability.CapabilityInfo, 0, len(r.metadata))
	for _, info := range r.metadata {
		infos = append(infos, info)
	}

	return infos
}

// FindByType 按类型查找能力
func (r *CapabilityRegistry) FindByType(capType string) []*capability.CapabilityInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []*capability.CapabilityInfo
	for _, info := range r.metadata {
		if info.Type == capType {
			infos = append(infos, info)
		}
	}

	return infos
}

// CapabilityManager AI能力管理器
type CapabilityManager struct {
	registry *CapabilityRegistry
}

// NewCapabilityManager 创建能力管理器
func NewCapabilityManager() *CapabilityManager {
	return &CapabilityManager{
		registry: NewCapabilityRegistry(),
	}
}

// RegisterCapability 注册能力
func (m *CapabilityManager) RegisterCapability(cap capability.AICapability) error {
	return m.registry.Register(cap)
}

// GetCapability 获取能力
func (m *CapabilityManager) GetCapability(capabilityID string) (capability.AICapability, error) {
	return m.registry.Get(capabilityID)
}

// ListCapabilities 列出所有能力
func (m *CapabilityManager) ListCapabilities() []*capability.CapabilityInfo {
	return m.registry.List()
}

// Execute 执行能力
func (m *CapabilityManager) Execute(ctx context.Context, capabilityID string, input *capability.CapabilityInput) (*capability.CapabilityOutput, error) {
	cap, err := m.registry.Get(capabilityID)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if err := cap.Validate(input); err != nil {
		return nil, err
	}

	return cap.Execute(ctx, input)
}

// ExecuteIntentRecognition 执行意图识别
func (m *CapabilityManager) ExecuteIntentRecognition(ctx context.Context, input string) (*capability.Intent, error) {
	// 查找意图识别能力
	infos := m.registry.FindByType("intent_recognition")
	if len(infos) == 0 {
		return nil, ErrCapabilityNotFound
	}

	cap, err := m.registry.Get(infos[0].CapabilityID)
	if err != nil {
		return nil, err
	}

	recognizer, ok := cap.(capability.IntentRecognizer)
	if !ok {
		return nil, errors.New("capability does not implement IntentRecognizer")
	}

	return recognizer.RecognizeIntent(ctx, input)
}

// ExecuteObjectiveGeneration 执行学习目标生成
func (m *CapabilityManager) ExecuteObjectiveGeneration(ctx context.Context, params *capability.ObjectiveParams) ([]*capability.Objective, error) {
	// 查找目标生成能力
	infos := m.registry.FindByType("objective_generation")
	if len(infos) == 0 {
		return nil, ErrCapabilityNotFound
	}

	cap, err := m.registry.Get(infos[0].CapabilityID)
	if err != nil {
		return nil, err
	}

	generator, ok := cap.(capability.ObjectiveGenerator)
	if !ok {
		return nil, errors.New("capability does not implement ObjectiveGenerator")
	}

	return generator.GenerateObjectives(ctx, params)
}

// ExecuteContentDiscovery 执行内容发现
func (m *CapabilityManager) ExecuteContentDiscovery(ctx context.Context, params *capability.DiscoveryParams) ([]*capability.Content, error) {
	// 查找内容发现能力
	infos := m.registry.FindByType("content_discovery")
	if len(infos) == 0 {
		return nil, ErrCapabilityNotFound
	}

	cap, err := m.registry.Get(infos[0].CapabilityID)
	if err != nil {
		return nil, err
	}

	discovery, ok := cap.(capability.ContentDiscovery)
	if !ok {
		return nil, errors.New("capability does not implement ContentDiscovery")
	}

	return discovery.DiscoverContent(ctx, params)
}

// ExecuteNarrativeGeneration 执行故事化叙述生成
func (m *CapabilityManager) ExecuteNarrativeGeneration(ctx context.Context, params *capability.NarrativeParams) (*capability.Narrative, error) {
	// 查找叙述生成能力
	infos := m.registry.FindByType("narrative_generation")
	if len(infos) == 0 {
		return nil, ErrCapabilityNotFound
	}

	cap, err := m.registry.Get(infos[0].CapabilityID)
	if err != nil {
		return nil, err
	}

	generator, ok := cap.(capability.NarrativeGenerator)
	if !ok {
		return nil, errors.New("capability does not implement NarrativeGenerator")
	}

	return generator.GenerateNarrative(ctx, params)
}

// ExecuteQualityAssessment 执行质量评估
func (m *CapabilityManager) ExecuteQualityAssessment(ctx context.Context, content interface{}) (*capability.QualityReport, error) {
	// 查找质量评估能力
	infos := m.registry.FindByType("quality_assessment")
	if len(infos) == 0 {
		return nil, ErrCapabilityNotFound
	}

	cap, err := m.registry.Get(infos[0].CapabilityID)
	if err != nil {
		return nil, err
	}

	assessor, ok := cap.(capability.QualityAssessor)
	if !ok {
		return nil, errors.New("capability does not implement QualityAssessor")
	}

	return assessor.AssessQuality(ctx, content)
}
