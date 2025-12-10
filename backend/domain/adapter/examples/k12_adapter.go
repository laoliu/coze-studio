// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/adapter"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
)

// K12Adapter K12教育领域适配器
// 实现 DomainAdapter 接口
type K12Adapter struct {
	info   *entity.AdapterInfo
	config *K12Config
}

// K12Config K12适配器配置
type K12Config struct {
	SupportedSubjects []string          // 支持的科目
	GradeMapping      map[string]string // 年级映射
	CurriculumStandard string           // 课程标准
}

// NewK12Adapter 创建K12适配器
func NewK12Adapter() *K12Adapter {
	return &K12Adapter{
		info: &entity.AdapterInfo{
			AdapterID:   "k12_education",
			Name:        "k12_education",
			DisplayName: "K12教育",
			Version:     "1.0.0",
			Description: "K12教育领域适配器，支持化学、物理、数学、生物等学科",
			Type:        entity.AdapterTypeK12,
			Capabilities: &entity.AdapterCapabilities{
				ActivityTypes: []entity.ActivityType{
					entity.ActivityTypeConcept,
					entity.ActivityTypeExperiment,
					entity.ActivityTypeProblem,
					entity.ActivityTypeProject,
				},
				Domains:   []string{"chemistry", "physics", "mathematics", "biology"},
				Features:  []string{"rag", "multimodal", "interactive"},
				Grades:    []string{"grade_7", "grade_8", "grade_9", "grade_10", "grade_11", "grade_12"},
				Languages: []string{"zh-CN", "en-US"},
			},
		},
		config: &K12Config{
			SupportedSubjects: []string{"化学", "物理", "数学", "生物"},
			GradeMapping: map[string]string{
				"初一": "grade_7",
				"初二": "grade_8",
				"初三": "grade_9",
				"高一": "grade_10",
				"高二": "grade_11",
				"高三": "grade_12",
			},
			CurriculumStandard: "中国教育部课程标准2022版",
		},
	}
}

// GetInfo 获取适配器信息
func (a *K12Adapter) GetInfo() *entity.AdapterInfo {
	return a.info
}

// ParseRequest 解析用户请求
func (a *K12Adapter) ParseRequest(input *entity.UserInput) (*entity.RequestContext, error) {
	// 1. 验证输入
	if input.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	
	// 2. 标准化年级
	grade := input.Grade
	if mappedGrade, ok := a.config.GradeMapping[input.Grade]; ok {
		grade = mappedGrade
	}
	
	// 3. 构建请求上下文
	ctx := &entity.RequestContext{
		Topic:        input.Topic,
		Domain:       input.Domain,
		Grade:        grade,
		Duration:     input.Duration,
		ActivityType: input.ActivityType,
		Language:     input.Language,
		DomainSpecific: map[string]any{
			"curriculum_standard": a.config.CurriculumStandard,
			"subject":            input.Domain,
		},
	}
	
	// 4. 设置默认值
	if ctx.Language == "" {
		ctx.Language = "zh-CN"
	}
	if ctx.Duration == 0 {
		ctx.Duration = 45 // 默认45分钟（一节课）
	}
	
	// 5. 添加学习者画像
	if input.Grade != "" {
		ctx.LearnerProfile = &entity.LearnerProfile{
			Grade:          grade,
			KnowledgeLevel: a.inferKnowledgeLevel(grade),
			LearningStyle:  "mixed", // 默认混合式
		}
	}
	
	return ctx, nil
}

// ValidateContext 验证请求上下文
func (a *K12Adapter) ValidateContext(ctx *entity.RequestContext) error {
	// 1. 验证必需字段
	if ctx.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	
	// 2. 验证科目
	validSubject := false
	for _, subject := range a.config.SupportedSubjects {
		if ctx.Domain == subject {
			validSubject = true
			break
		}
	}
	if !validSubject {
		return fmt.Errorf("unsupported subject: %s", ctx.Domain)
	}
	
	// 3. 验证年级
	if ctx.Grade != "" {
		validGrade := false
		for _, grade := range a.info.Capabilities.Grades {
			if ctx.Grade == grade {
				validGrade = true
				break
			}
		}
		if !validGrade {
			return fmt.Errorf("unsupported grade: %s", ctx.Grade)
		}
	}
	
	// 4. 验证活动类型
	if ctx.ActivityType != "" {
		validActivity := false
		for _, at := range a.info.Capabilities.ActivityTypes {
			if ctx.ActivityType == at {
				validActivity = true
				break
			}
		}
		if !validActivity {
			return fmt.Errorf("unsupported activity type: %s", ctx.ActivityType)
		}
	}
	
	// 5. 验证时长
	if ctx.Duration < 10 || ctx.Duration > 120 {
		return fmt.Errorf("duration must be between 10 and 120 minutes")
	}
	
	return nil
}

// GenerateLearningObjectives 生成学习目标
func (a *K12Adapter) GenerateLearningObjectives(ctx *entity.RequestContext) ([]*entity.LearningObjective, error) {
	// TODO: 调用AI能力层的ObjectiveGenerator
	// 这里先返回占位符数据
	
	objectives := []*entity.LearningObjective{
		{
			ID:          "obj_1",
			Objective:   fmt.Sprintf("理解%s的基本概念和原理", ctx.Topic),
			Level:       "理解",
			Category:    "知识",
			IsMeasurable: true,
			Assessment:  "概念解释、案例分析",
		},
		{
			ID:          "obj_2",
			Objective:   fmt.Sprintf("能够应用%s解决实际问题", ctx.Topic),
			Level:       "应用",
			Category:    "技能",
			IsMeasurable: true,
			Assessment:  "问题解决、实验操作",
		},
	}
	
	return objectives, nil
}

// DiscoverContent 发现内容
func (a *K12Adapter) DiscoverContent(ctx *entity.RequestContext, objectives []*entity.LearningObjective) ([]*entity.Content, error) {
	// TODO: 调用AI能力层的ContentDiscovery
	// 这里先返回占位符数据
	
	contents := []*entity.Content{
		{
			ID:          "content_1",
			Type:        "text",
			Title:       fmt.Sprintf("%s - 概念介绍", ctx.Topic),
			Description: "基础概念和理论背景",
			Source:      "教材",
			Relevance:   0.95,
		},
		{
			ID:          "content_2",
			Type:        "video",
			Title:       fmt.Sprintf("%s - 实验演示", ctx.Topic),
			Description: "实验步骤和现象观察",
			Source:      "教学视频库",
			URL:         "https://example.com/video/demo",
			Relevance:   0.88,
		},
	}
	
	return contents, nil
}

// CustomizeWorkflow 定制工作流
func (a *K12Adapter) CustomizeWorkflow(ctx *entity.RequestContext) (*entity.WorkflowConfig, error) {
	// 根据活动类型选择不同的工作流模板
	var template string
	var nodes []*entity.NodeConfig
	
	switch ctx.ActivityType {
	case entity.ActivityTypeConcept:
		template = "concept_understanding_workflow"
		nodes = a.buildConceptWorkflowNodes(ctx)
		
	case entity.ActivityTypeExperiment:
		template = "experiment_workflow"
		nodes = a.buildExperimentWorkflowNodes(ctx)
		
	case entity.ActivityTypeProblem:
		template = "problem_solving_workflow"
		nodes = a.buildProblemSolvingWorkflowNodes(ctx)
		
	default:
		template = "default_workflow"
		nodes = a.buildDefaultWorkflowNodes(ctx)
	}
	
	return &entity.WorkflowConfig{
		WorkflowID: fmt.Sprintf("workflow_%s_%s", ctx.Domain, ctx.ActivityType),
		Template:   template,
		Nodes:      nodes,
		Variables: map[string]any{
			"topic":    ctx.Topic,
			"grade":    ctx.Grade,
			"duration": ctx.Duration,
		},
		Timeout: 300, // 5分钟超时
	}, nil
}

// FormatOutput 格式化输出
func (a *K12Adapter) FormatOutput(result *entity.WorkflowResult) (interface{}, error) {
	// 格式化为K12教学活动格式
	output := map[string]interface{}{
		"activity_id":   result.WorkflowID,
		"title":         result.Output["title"],
		"objectives":    result.Output["objectives"],
		"content":       result.Output["content"],
		"narrative":     result.Output["narrative"],
		"assessment":    result.Output["assessment"],
		"duration":      result.Output["duration"],
		"materials":     result.Output["materials"],
		"instructions":  result.Output["instructions"],
	}
	
	return output, nil
}

// ValidateOutput 验证输出质量
func (a *K12Adapter) ValidateOutput(output interface{}) (*entity.QualityReport, error) {
	// TODO: 调用AI能力层的QualityAssessor
	// 这里先返回占位符数据
	
	report := &entity.QualityReport{
		OverallScore: 85.0,
		Dimensions: map[string]float64{
			"内容准确性": 90.0,
			"教学适切性": 85.0,
			"目标对齐度": 88.0,
			"语言规范性": 82.0,
		},
		Issues: []*entity.QualityIssue{
			{
				Severity:    "minor",
				Category:    "语言",
				Description: "部分术语表述不够规范",
				Suggestion:  "建议使用标准教材术语",
			},
		},
		Suggestions: []string{
			"可以增加更多实例说明",
			"建议添加互动环节",
		},
		IsApproved: true,
	}
	
	return report, nil
}

// 辅助方法

// inferKnowledgeLevel 根据年级推断知识水平
func (a *K12Adapter) inferKnowledgeLevel(grade string) string {
	switch grade {
	case "grade_7", "grade_8":
		return "beginner"
	case "grade_9", "grade_10":
		return "intermediate"
	case "grade_11", "grade_12":
		return "advanced"
	default:
		return "intermediate"
	}
}

// buildConceptWorkflowNodes 构建概念理解工作流节点
func (a *K12Adapter) buildConceptWorkflowNodes(ctx *entity.RequestContext) []*entity.NodeConfig {
	return []*entity.NodeConfig{
		{
			NodeID:   "persona_node",
			NodeType: "LLM",
			Config: map[string]any{
				"prompt": "生成学习者画像",
			},
		},
		{
			NodeID:   "objective_node",
			NodeType: "Plugin",
			Config: map[string]any{
				"capability": "objective_generation",
			},
			Dependencies: []string{"persona_node"},
		},
		{
			NodeID:   "content_node",
			NodeType: "Plugin",
			Config: map[string]any{
				"capability": "content_discovery",
			},
			Dependencies: []string{"objective_node"},
		},
		{
			NodeID:   "narrative_node",
			NodeType: "LLM",
			Config: map[string]any{
				"capability": "narrative_generation",
			},
			Dependencies: []string{"content_node"},
		},
	}
}

// buildExperimentWorkflowNodes 构建实验探究工作流节点
func (a *K12Adapter) buildExperimentWorkflowNodes(ctx *entity.RequestContext) []*entity.NodeConfig {
	return []*entity.NodeConfig{
		{
			NodeID:   "experiment_design_node",
			NodeType: "Plugin",
			Config: map[string]any{
				"component": "experiment_designer",
			},
		},
		{
			NodeID:   "safety_check_node",
			NodeType: "Code",
			Config: map[string]any{
				"script": "check_safety_requirements",
			},
			Dependencies: []string{"experiment_design_node"},
		},
		{
			NodeID:   "procedure_node",
			NodeType: "LLM",
			Config: map[string]any{
				"prompt": "生成实验步骤",
			},
			Dependencies: []string{"safety_check_node"},
		},
	}
}

// buildProblemSolvingWorkflowNodes 构建问题解决工作流节点
func (a *K12Adapter) buildProblemSolvingWorkflowNodes(ctx *entity.RequestContext) []*entity.NodeConfig {
	return []*entity.NodeConfig{
		{
			NodeID:   "problem_generation_node",
			NodeType: "LLM",
			Config: map[string]any{
				"prompt": "生成问题",
			},
		},
		{
			NodeID:   "solution_node",
			NodeType: "LLM",
			Config: map[string]any{
				"prompt": "生成解决方案",
			},
			Dependencies: []string{"problem_generation_node"},
		},
	}
}

// buildDefaultWorkflowNodes 构建默认工作流节点
func (a *K12Adapter) buildDefaultWorkflowNodes(ctx *entity.RequestContext) []*entity.NodeConfig {
	return []*entity.NodeConfig{
		{
			NodeID:   "default_node",
			NodeType: "LLM",
			Config: map[string]any{
				"prompt": "生成教学活动",
			},
		},
	}
}

// 确保 K12Adapter 实现了 DomainAdapter 接口
var _ adapter.DomainAdapter = (*K12Adapter)(nil)
