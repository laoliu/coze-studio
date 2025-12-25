// Copyright 2025 Coze Studio. All rights reserved.

package adapter

import (
	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
)

// DomainAdapter 领域适配器接口
// 所有领域适配器必须实现此接口
type DomainAdapter interface {
	// GetInfo 获取适配器信息
	GetInfo() *entity.AdapterInfo

	// ParseRequest 解析用户请求
	// 将用户输入转换为标准化的请求上下文
	ParseRequest(input *entity.UserInput) (*entity.RequestContext, error)

	// ValidateContext 验证请求上下文
	// 检查请求上下文是否符合领域规则
	ValidateContext(ctx *entity.RequestContext) error

	// GenerateLearningObjectives 生成学习目标
	// 根据上下文生成符合领域标准的学习目标
	GenerateLearningObjectives(ctx *entity.RequestContext) ([]*entity.LearningObjective, error)

	// DiscoverContent 发现内容
	// 从知识库或外部源发现相关内容
	DiscoverContent(ctx *entity.RequestContext, objectives []*entity.LearningObjective) ([]*entity.Content, error)

	// CustomizeWorkflow 定制工作流
	// 根据领域特性定制工作流执行逻辑
	CustomizeWorkflow(ctx *entity.RequestContext) (*entity.WorkflowConfig, error)

	// FormatOutput 格式化输出
	// 将执行结果格式化为领域特定的输出格式
	FormatOutput(result *entity.WorkflowResult) (interface{}, error)

	// ValidateOutput 验证输出质量
	// 根据领域标准验证输出质量
	ValidateOutput(output interface{}) (*entity.QualityReport, error)
}
