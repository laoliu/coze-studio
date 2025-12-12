package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// WorkflowGenerator 工作流自动生成器
type WorkflowGenerator struct {
	llmModel       model.BaseChatModel
	templateRepo   *TemplateRepository
	layoutEngine   *AutoLayoutEngine
	promptTemplate *PromptTemplate
}

// NewWorkflowGenerator 创建工作流生成器
func NewWorkflowGenerator(llmModel model.BaseChatModel) (*WorkflowGenerator, error) {
	templateRepo, err := NewTemplateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create template repository: %w", err)
	}

	return &WorkflowGenerator{
		llmModel:       llmModel,
		templateRepo:   templateRepo,
		layoutEngine:   NewAutoLayoutEngine(),
		promptTemplate: NewPromptTemplate(),
	}, nil
}

// GenerateWorkflow 根据需求生成完整工作流
func (g *WorkflowGenerator) GenerateWorkflow(
	ctx context.Context,
	req *WorkflowGenerationRequest,
) (*WorkflowGenerationResponse, error) {
	startTime := time.Now()
	logs.CtxInfof(ctx, "[WorkflowGenerator] Start generating workflow for requirement: %s (language: %s)", 
		req.UserRequirement, req.Language)

	// 步骤 1：意图分析
	step1Start := time.Now()
	intent, err := g.analyzeIntent(ctx, req)
	if err != nil {
		logs.CtxErrorf(ctx, "[WorkflowGenerator] Intent analysis failed after %v: %v", 
			time.Since(step1Start), err)
		return nil, fmt.Errorf("意图分析失败: %w", err)
	}
	logs.CtxInfof(ctx, "[WorkflowGenerator] Intent analysis completed in %v: type=%s, confidence=%.2f, steps=%d", 
		time.Since(step1Start), intent.Intent.WorkflowType, intent.Intent.Confidence, len(intent.KeySteps))

	// 步骤 2：匹配最佳实践模板
	step2Start := time.Now()
	templates := g.templateRepo.FindSimilarTemplates(intent.Intent.WorkflowType, 3)
	logs.CtxInfof(ctx, "[WorkflowGenerator] Found %d similar templates in %v", 
		len(templates), time.Since(step2Start))

	// 步骤 3：LLM 生成工作流结构
	step3Start := time.Now()
	rawWorkflow, err := g.llmGenerateWorkflow(ctx, req, intent, templates)
	if err != nil {
		logs.CtxErrorf(ctx, "[WorkflowGenerator] LLM generation failed after %v: %v", 
			time.Since(step3Start), err)
		return nil, fmt.Errorf("LLM 生成失败: %w", err)
	}
	logs.CtxInfof(ctx, "[WorkflowGenerator] Generated workflow in %v: name=%s, nodes=%d, edges=%d, confidence=%.2f", 
		time.Since(step3Start), rawWorkflow.WorkflowName, len(rawWorkflow.Nodes), 
		len(rawWorkflow.Edges), rawWorkflow.Confidence)

	// 步骤 4：节点配置自动推断
	step4Start := time.Now()
	configuredWorkflow := g.configureNodes(rawWorkflow, intent)
	logs.CtxInfof(ctx, "[WorkflowGenerator] Node configuration completed in %v", 
		time.Since(step4Start))

	// 步骤 5：自动布局
	step5Start := time.Now()
	layoutWorkflow := g.layoutEngine.AutoLayout(configuredWorkflow)
	logs.CtxInfof(ctx, "[WorkflowGenerator] Auto layout completed in %v", 
		time.Since(step5Start))

	// 步骤 6：生成解释说明
	step6Start := time.Now()
	explanations := g.generateExplanations(layoutWorkflow, rawWorkflow)
	logs.CtxInfof(ctx, "[WorkflowGenerator] Generated %d explanations in %v", 
		len(explanations), time.Since(step6Start))

	// 构造响应
	response := &WorkflowGenerationResponse{
		WorkflowID:       generateWorkflowID(),
		Nodes:            layoutWorkflow.Nodes,
		Edges:            layoutWorkflow.Edges,
		Explanation:      rawWorkflow.OverallExplanation,
		NodeExplanations: explanations,
		Confidence:       rawWorkflow.Confidence,
	}

	totalTime := time.Since(startTime)
	logs.CtxInfof(ctx, "[WorkflowGenerator] Workflow generation completed successfully in %v (intent:%v, templates:%v, llm:%v, config:%v, layout:%v, explain:%v)", 
		totalTime, 
		step2Start.Sub(step1Start), 
		step3Start.Sub(step2Start),
		step4Start.Sub(step3Start),
		step5Start.Sub(step4Start),
		step6Start.Sub(step5Start),
		time.Since(step6Start))

	return response, nil
}

// analyzeIntent 分析用户需求意图
func (g *WorkflowGenerator) analyzeIntent(
	ctx context.Context,
	req *WorkflowGenerationRequest,
) (*IntentAnalysisResult, error) {
	prompt := g.promptTemplate.BuildIntentAnalysisPrompt(req.UserRequirement)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := g.llmModel.Generate(ctx, messages, model.WithTemperature(0.3))
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 解析 JSON 响应
	content := response.Content
	
	// 清理可能的 markdown 代码块标记
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var intent IntentAnalysisResult
	if err := json.Unmarshal([]byte(content), &intent); err != nil {
		logs.CtxErrorf(ctx, "[WorkflowGenerator] Failed to parse intent JSON: %v, content: %s", err, content)
		return nil, fmt.Errorf("解析意图分析结果失败: %w", err)
	}

	return &intent, nil
}

// llmGenerateWorkflow 使用 LLM 生成工作流结构
func (g *WorkflowGenerator) llmGenerateWorkflow(
	ctx context.Context,
	req *WorkflowGenerationRequest,
	intent *IntentAnalysisResult,
	templates []*WorkflowTemplate,
) (*RawWorkflow, error) {
	prompt := g.promptTemplate.BuildWorkflowGenerationPrompt(
		req.UserRequirement,
		intent,
		templates,
	)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := g.llmModel.Generate(ctx, messages, model.WithTemperature(0.5))
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 解析 JSON 响应
	content := response.Content
	
	// 清理可能的 markdown 代码块标记
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var workflow RawWorkflow
	if err := json.Unmarshal([]byte(content), &workflow); err != nil {
		logs.CtxErrorf(ctx, "[WorkflowGenerator] Failed to parse workflow JSON: %v, content: %s", err, content)
		return nil, fmt.Errorf("解析工作流结构失败: %w", err)
	}

	return &workflow, nil
}

// configureNodes 自动配置节点
func (g *WorkflowGenerator) configureNodes(
	workflow *RawWorkflow,
	intent *IntentAnalysisResult,
) *ConfiguredWorkflow {
	configured := &ConfiguredWorkflow{
		Nodes:        make([]*InternalNode, 0, len(workflow.Nodes)),
		Edges:        make([]*InternalEdge, 0, len(workflow.Edges)),
		Explanations: make(map[string]string),
	}

	// 转换节点
	for _, rawNode := range workflow.Nodes {
		node := &InternalNode{
			ID:       rawNode.ID,
			Type:     string(rawNode.Type),
			Name:     rawNode.Name,
			Position: rawNode.Position,
			Config:   rawNode.Config,
		}

		// 如果没有位置信息，会由布局引擎自动计算
		if node.Position == nil {
			node.Position = &Position{X: 0, Y: 0}
		}

		configured.Nodes = append(configured.Nodes, node)
		configured.Explanations[node.ID] = rawNode.Explanation
	}

	// 转换边
	for _, rawEdge := range workflow.Edges {
		edge := &InternalEdge{
			From:      rawEdge.From,
			To:        rawEdge.To,
			OutputKey: rawEdge.OutputKey,
			InputKey:  rawEdge.InputKey,
		}
		configured.Edges = append(configured.Edges, edge)
	}

	configured.ConfidenceScore = workflow.Confidence

	return configured
}

// generateExplanations 生成工作流说明
func (g *WorkflowGenerator) generateExplanations(
	workflow *ConfiguredWorkflow,
	rawWorkflow *RawWorkflow,
) map[string]string {
	explanations := make(map[string]string)

	// 合并现有的解释
	for nodeID, explanation := range workflow.Explanations {
		explanations[nodeID] = explanation
	}

	return explanations
}

// generateWorkflowID 生成工作流 ID
func generateWorkflowID() string {
	// 简单实现，使用时间戳
	return fmt.Sprintf("workflow_%d", time.Now().UnixNano())
}
