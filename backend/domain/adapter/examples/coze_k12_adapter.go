// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	crossworkflow "github.com/coze-dev/coze-studio/backend/crossdomain/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/prompts"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/repository"
)

// CozeK12Adapter 使用 Coze 原生 modelbuilder 的 K12 适配器
// 这是生产环境推荐的实现，直接使用 Coze 的 LLM 基础设施
type CozeK12Adapter struct {
	chatModel   modelbuilder.ToolCallingChatModel
	contentRepo repository.ContentRepository
	promptMgr   *prompts.PromptManager
}

// NewCozeK12AdapterByModelID 通过 Coze Model ID 创建适配器
// modelID: Coze 数据库中的 model_id (在 Admin -> Models 中配置)
func NewCozeK12AdapterByModelID(ctx context.Context, modelID int64, contentRepo repository.ContentRepository) (*CozeK12Adapter, error) {
	// 使用 Coze 的 modelbuilder 从数据库加载 Model 配置
	chatModel, _, err := modelbuilder.BuildModelByID(ctx, modelID, nil)
	if err != nil {
		return nil, fmt.Errorf("build model by id failed: %w", err)
	}

	return &CozeK12Adapter{
		chatModel:   chatModel,
		contentRepo: contentRepo,
		promptMgr:   prompts.NewPromptManager(),
	}, nil
}

// NewCozeK12AdapterBySettings 通过 Bot Settings 创建适配器
// modelInfo: 从 Bot 配置中获取的 Model 信息
func NewCozeK12AdapterBySettings(ctx context.Context, modelInfo *bot_common.ModelInfo, contentRepo repository.ContentRepository) (*CozeK12Adapter, error) {
	// 使用 Coze 的 modelbuilder 从 ModelInfo 构建
	chatModel, _, err := modelbuilder.BuildModelBySettings(ctx, modelInfo)
	if err != nil {
		return nil, fmt.Errorf("build model from settings failed: %w", err)
	}

	return &CozeK12Adapter{
		chatModel:   chatModel,
		contentRepo: contentRepo,
		promptMgr:   prompts.NewPromptManager(),
	}, nil
}

// GetInfo 返回适配器信息
func (a *CozeK12Adapter) GetInfo() entity.AdapterInfo {
	return entity.AdapterInfo{
		AdapterID:   "coze_k12",
		Name:        "coze_k12",
		DisplayName: "Coze K12 Adapter",
		Version:     "1.0.0",
		Description: "K12 教育适配器 (使用 Coze ModelBuilder)",
		Type:        entity.AdapterTypeK12,
		Capabilities: &entity.AdapterCapabilities{
			Domains:   []string{"K12", "STEM", "语言学习", "数学"},
			Grades:    []string{"1-12"},
			Languages: []string{"zh-CN", "en-US"},
		},
	}
}

// ParseRequest 解析用户请求
func (a *CozeK12Adapter) ParseRequest(input *entity.UserInput) (*entity.AdapterContext, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	// 基础验证
	if input.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	if input.Domain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	// 解析年级
	grade := parseGrade(input.Grade)

	ctx := &AdapterContext{
		Context: context.Background(),
		Input:   input,
		Domain:  input.Domain,
		Grade:   grade,
		Metadata: map[string]any{
			"duration":      input.Duration,
			"activity_type": input.ActivityType,
			"language":      input.Language,
		},
	}

	return ctx, nil
}

// GenerateLearningObjectives 生成学习目标（使用 Coze LLM）
func (a *CozeK12Adapter) GenerateLearningObjectives(ctx *entity.AdapterContext) ([]entity.LearningObjective, error) {
	// 1. 渲染提示词模板
	prompt, err := a.promptMgr.RenderPrompt("generate_objectives", map[string]any{
		"Topic":        ctx.Input.Topic,
		"Domain":       ctx.Input.Domain,
		"Grade":        ctx.Grade,
		"Duration":     ctx.Input.Duration,
		"ActivityType": ctx.Input.ActivityType,
	})
	if err != nil {
		return nil, fmt.Errorf("render prompt failed: %w", err)
	}

	// 2. 构造 Eino 消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一位经验丰富的K12教育专家，擅长设计科学的学习目标。"),
		schema.UserMessage(prompt + "\n\n请以JSON格式返回学习目标列表。"),
	}

	// 3. 调用 Coze 的 ChatModel (通过 Eino)
	resp, err := a.chatModel.Generate(ctx.Context, messages,
		model.WithResponseFormat(&model.ResponseFormat{
			Type: model.ResponseFormatJSONObject,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("llm generate failed: %w", err)
	}

	// 4. 解析 JSON 响应
	var objectives []entity.LearningObjective
	if err := json.Unmarshal([]byte(resp.Content), &objectives); err != nil {
		return nil, fmt.Errorf("parse json failed: %w", err)
	}

	return objectives, nil
}

// DiscoverContent 发现相关内容
func (a *CozeK12Adapter) DiscoverContent(ctx *entity.AdapterContext, objectives []entity.LearningObjective) ([]*entity.ContentItem, error) {
	// 1. 使用 LLM 生成搜索关键词
	keywords, err := a.generateSearchKeywords(ctx, objectives)
	if err != nil {
		return nil, fmt.Errorf("generate keywords failed: %w", err)
	}

	// 2. 从数据库搜索内容
	contents, err := a.contentRepo.SearchByKeywords(ctx.Context, keywords, 10)
	if err != nil {
		return nil, fmt.Errorf("search content failed: %w", err)
	}

	// 3. 如果没有找到内容，使用 LLM 推荐
	if len(contents) == 0 {
		return a.recommendContentByLLM(ctx, objectives)
	}

	// 4. 转换为适配器格式
	return a.convertToContentItems(contents), nil
}

// generateSearchKeywords 使用 LLM 生成搜索关键词
func (a *CozeK12Adapter) generateSearchKeywords(ctx *entity.AdapterContext, objectives []entity.LearningObjective) ([]string, error) {
	prompt, err := a.promptMgr.RenderPrompt("generate_search_keywords", map[string]any{
		"Topic":      ctx.Input.Topic,
		"Domain":     ctx.Input.Domain,
		"Grade":      ctx.Grade,
		"Objectives": objectives,
	})
	if err != nil {
		return nil, err
	}

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := a.chatModel.Generate(ctx.Context, messages,
		model.WithResponseFormat(&model.ResponseFormat{
			Type: model.ResponseFormatJSONObject,
		}),
	)
	if err != nil {
		return nil, err
	}

	var result struct {
		Keywords []string `json:"keywords"`
	}
	if err := json.Unmarshal([]byte(resp.Content), &result); err != nil {
		return nil, err
	}

	return result.Keywords, nil
}

// recommendContentByLLM 当数据库无内容时，使用 LLM 推荐
func (a *CozeK12Adapter) recommendContentByLLM(ctx *entity.AdapterContext, objectives []entity.LearningObjective) ([]*entity.ContentItem, error) {
	prompt, err := a.promptMgr.RenderPrompt("recommend_content", map[string]any{
		"Topic":      ctx.Input.Topic,
		"Domain":     ctx.Input.Domain,
		"Grade":      ctx.Grade,
		"Objectives": objectives,
	})
	if err != nil {
		return nil, err
	}

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := a.chatModel.Generate(ctx.Context, messages,
		model.WithResponseFormat(&model.ResponseFormat{
			Type: model.ResponseFormatJSONObject,
		}),
	)
	if err != nil {
		return nil, err
	}

	var result struct {
		Contents []*entity.ContentItem `json:"contents"`
	}
	if err := json.Unmarshal([]byte(resp.Content), &result); err != nil {
		return nil, err
	}

	return result.Contents, nil
}


// ExecuteWorkflow 使用 Coze 原生 Workflow 系统执行工作流
// 这个方法直接调用 Coze 的 Workflow 引擎，而不是手动构造配置
func (a *CozeK12Adapter) ExecuteWorkflow(ctx *entity.AdapterContext) (*entity.WorkflowResult, error) {
	// 1. 根据活动类型选择工作流 ID
	// 注意：这些 Workflow 需要预先在 Coze 系统中创建
	// 可以通过 Admin -> Workflows 界面创建，或使用 domain/workflow 的 Create 方法创建
	var workflowID int64
	switch ctx.Input.ActivityType {
	case entity.ActivityTypeConcept:
		workflowID = 1000001 // K12 概念学习 Workflow ID（需要在 Coze 中预先创建）
	case entity.ActivityTypePractice:
		workflowID = 1000002 // K12 练习 Workflow ID
	case entity.ActivityTypeExploration:
		workflowID = 1000003 // K12 探索学习 Workflow ID
	default:
		workflowID = 1000000 // K12 通用教学 Workflow ID
	}

	// 2. 从 Metadata 获取执行配置（带默认值）
	var userID, connectorID int64
	var connectorUID string
	
	if ctx.Metadata != nil {
		if uid, ok := ctx.Metadata["user_id"].(int64); ok {
			userID = uid
		}
		if cid, ok := ctx.Metadata["connector_id"].(int64); ok {
			connectorID = cid
		}
		if cuid, ok := ctx.Metadata["connector_uid"].(string); ok {
			connectorUID = cuid
		}
	}
	
	// 如果没有提供，使用默认值
	if userID == 0 {
		userID = 1 // 默认用户
	}
	if connectorID == 0 {
		connectorID = 1 // 默认 Connector
	}
	if connectorUID == "" {
		connectorUID = "system" // 默认系统用户
	}

	// 3. 配置执行参数
	executeConfig := workflowModel.ExecuteConfig{
		ID:           workflowID,
		From:         workflowModel.FromLatestVersion, // 使用最新版本
		Operator:     userID,
		ConnectorID:  connectorID,
		ConnectorUID: connectorUID,
		Mode:         workflowModel.ExecuteModeRelease, // 生产模式
		BizType:      workflowModel.BizTypeWorkflow,
		SyncPattern:  workflowModel.SyncPatternSync, // 同步执行
		TaskType:     workflowModel.TaskTypeForeground,
	}

	// 4. 准备输入数据
	input := map[string]any{
		"topic":         ctx.Input.Topic,
		"domain":        ctx.Input.Domain,
		"grade":         strconv.Itoa(ctx.Grade),
		"duration":      ctx.Input.Duration,
		"activity_type": string(ctx.Input.ActivityType),
		"language":      ctx.Input.Language,
	}

	// 添加额外参数
	for k, v := range ctx.Input.ExtraParams {
		input[k] = v
	}

	// 4. 调用 Coze 原生 Workflow 系统执行
	execution, _, err := crossworkflow.DefaultSVC().SyncExecute(
		ctx.Context,
		executeConfig,
		input,
	)
	if err != nil {
		return nil, fmt.Errorf("workflow execution failed: %w", err)
	}

	// 5. 转换为 Adapter 的结果格式
	result := &entity.WorkflowResult{
		WorkflowID:    strconv.FormatInt(workflowID, 10),
		Status:        string(execution.Status),
		Output:        execution.Output,
		ExecutionTime: execution.ExecutionTime,
	}

	// 检查执行状态
	if execution.Status != "success" && execution.Status != "completed" {
		result.ErrorMessage = fmt.Sprintf("workflow execution status: %s", execution.Status)
	}

	return result, nil
}

// ExecuteWorkflowStream 流式执行工作流（适用于需要实时反馈的场景）
func (a *CozeK12Adapter) ExecuteWorkflowStream(ctx *entity.AdapterContext, callback func(*crossworkflow.WorkflowMessage) error) (*entity.WorkflowResult, error) {
	// 1. 选择 Workflow ID（与同步版本相同）
	var workflowID int64
	switch ctx.Input.ActivityType {
	case entity.ActivityTypeConcept:
		workflowID = 1000001
	case entity.ActivityTypePractice:
		workflowID = 1000002
	case entity.ActivityTypeExploration:
		workflowID = 1000003
	default:
		workflowID = 1000000
	}

	// 2. 从 Metadata 获取执行配置（带默认值）
	var userID, connectorID int64
	var connectorUID string
	
	if ctx.Metadata != nil {
		if uid, ok := ctx.Metadata["user_id"].(int64); ok {
			userID = uid
		}
		if cid, ok := ctx.Metadata["connector_id"].(int64); ok {
			connectorID = cid
		}
		if cuid, ok := ctx.Metadata["connector_uid"].(string); ok {
			connectorUID = cuid
		}
	}
	
	if userID == 0 {
		userID = 1
	}
	if connectorID == 0 {
		connectorID = 1
	}
	if connectorUID == "" {
		connectorUID = "system"
	}

	// 3. 配置执行参数
	executeConfig := workflowModel.ExecuteConfig{
		ID:           workflowID,
		From:         workflowModel.FromLatestVersion,
		Operator:     userID,
		ConnectorID:  connectorID,
		ConnectorUID: connectorUID,
		Mode:         workflowModel.ExecuteModeRelease,
		BizType:      workflowModel.BizTypeWorkflow,
		SyncPattern:  workflowModel.SyncPatternStream, // 流式模式
		TaskType:     workflowModel.TaskTypeForeground,
	}

	// 3. 准备输入
	input := map[string]any{
		"topic":         ctx.Input.Topic,
		"domain":        ctx.Input.Domain,
		"grade":         strconv.Itoa(ctx.Grade),
		"duration":      ctx.Input.Duration,
		"activity_type": string(ctx.Input.ActivityType),
		"language":      ctx.Input.Language,
	}

	// 4. 流式执行
	streamer, err := crossworkflow.DefaultSVC().StreamExecute(
		ctx.Context,
		executeConfig,
		input,
	)
	if err != nil {
		return nil, fmt.Errorf("workflow stream execution failed: %w", err)
	}

	// 5. 处理流式消息
	var finalOutput map[string]any
	var status string
	var executionTime int64

	for {
		msg, err := streamer.Recv()
		if err != nil {
			break
		}

		// 调用回调函数处理每条消息
		if callback != nil {
			if err := callback(msg); err != nil {
				return nil, fmt.Errorf("callback error: %w", err)
			}
		}

		// 收集最终结果
		if msg != nil && msg.Type == "workflow_finished" {
			status = "success"
			if msg.Extra != nil {
				if output, ok := msg.Extra["output"].(map[string]any); ok {
					finalOutput = output
				}
				if execTime, ok := msg.Extra["execution_time"].(int64); ok {
					executionTime = execTime
				}
			}
		}
	}

	// 6. 构造结果
	result := &entity.WorkflowResult{
		WorkflowID:    strconv.FormatInt(workflowID, 10),
		Status:        status,
		Output:        finalOutput,
		ExecutionTime: executionTime,
	}

	return result, nil
}

// CustomizeWorkflow 保留作为兼容方法，但内部调用 ExecuteWorkflow
// Deprecated: 请使用 ExecuteWorkflow，这个方法将在未来版本中移除
func (a *CozeK12Adapter) CustomizeWorkflow(ctx *entity.AdapterContext) (*entity.WorkflowConfig, error) {
	// 为了向后兼容，保留这个方法但标记为废弃
	// 实际上我们应该直接使用 ExecuteWorkflow
	result, err := a.ExecuteWorkflow(ctx)
	if err != nil {
		return nil, err
	}

	// 将执行结果转换为旧的 WorkflowConfig 格式（用于兼容）
	config := &entity.WorkflowConfig{
		WorkflowID: result.WorkflowID,
		Template:   "coze_native_workflow",
		Variables: map[string]any{
			"status":         result.Status,
			"output":         result.Output,
			"execution_time": result.ExecutionTime,
		},
		Nodes:   []*entity.NodeConfig{}, // Coze Workflow 管理节点，这里留空
		Timeout: 300,                     // 默认 5 分钟超时
	}

	return config, nil
}

// FormatOutput 格式化输出（使用 LLM 生成叙事性文本）
func (a *CozeK12Adapter) FormatOutput(workflowResult *entity.WorkflowResult) (*entity.AdapterOutput, error) {
	// 1. 渲染提示词
	prompt, err := a.promptMgr.RenderPrompt("generate_narrative", map[string]any{
		"WorkflowID": workflowResult.WorkflowID,
		"Status":     workflowResult.Status,
		"Output":     workflowResult.Output,
	})
	if err != nil {
		return nil, fmt.Errorf("render prompt failed: %w", err)
	}

	// 2. 调用 LLM 生成叙事
	messages := []*schema.Message{
		schema.SystemMessage("你是一位教学设计专家，擅长将教学活动转化为清晰、吸引人的教案文本。"),
		schema.UserMessage(prompt),
	}

	resp, err := a.chatModel.Generate(context.Background(), messages)
	if err != nil {
		return nil, fmt.Errorf("llm generate failed: %w", err)
	}

	// 3. 构造输出
	output := &entity.AdapterOutput{
		Content:      resp.Content,
		Format:       entity.FormatMarkdown,
		Metadata:     workflowResult.Output,
		Confidence:   0.9,
		ProcessingMS: workflowResult.ExecutionTime,
	}

	return output, nil
}

// ValidateOutput 验证输出质量
func (a *CozeK12Adapter) ValidateOutput(output *entity.AdapterOutput) (*entity.QualityReport, error) {
	// 1. 渲染质量检查提示词
	prompt, err := a.promptMgr.RenderPrompt("validate_quality", map[string]any{
		"Content":  output.Content,
		"Metadata": output.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("render prompt failed: %w", err)
	}

	// 2. 调用 LLM 进行质量评估
	messages := []*schema.Message{
		schema.SystemMessage("你是一位教学质量评估专家，负责评估教学设计的质量。"),
		schema.UserMessage(prompt),
	}

	resp, err := a.chatModel.Generate(context.Background(), messages,
		model.WithResponseFormat(&model.ResponseFormat{
			Type: model.ResponseFormatJSONObject,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("llm generate failed: %w", err)
	}

	// 3. 解析质量报告
	var report entity.QualityReport
	if err := json.Unmarshal([]byte(resp.Content), &report); err != nil {
		return nil, fmt.Errorf("parse json failed: %w", err)
	}

	return &report, nil
}

// 辅助方法

func (a *CozeK12Adapter) convertToContentItems(contents []*repository.Content) []*entity.ContentItem {
	items := make([]*entity.ContentItem, 0, len(contents))
	for _, c := range contents {
		items = append(items, &ContentItem{
			ID:          c.ID,
			Type:        c.Type,
			Title:       c.Title,
			Description: c.Description,
			Source:      c.Source,
			Metadata: map[string]any{
				"domain":    c.Domain,
				"grade":     c.Grade,
				"topic":     c.Topic,
				"relevance": c.Relevance,
			},
		})
	}
	return items
}

// parseGrade 解析年级字符串为数字
func parseGrade(grade string) int {
	gradeMap := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6,
		"7": 7, "8": 8, "9": 9, "10": 10, "11": 11, "12": 12,
		"一年级": 1, "二年级": 2, "三年级": 3, "四年级": 4,
		"五年级": 5, "六年级": 6, "七年级": 7, "八年级": 8,
		"九年级": 9, "高一": 10, "高二": 11, "高三": 12,
	}
	if g, ok := gradeMap[grade]; ok {
		return g
	}
	return 6 // 默认值
}
