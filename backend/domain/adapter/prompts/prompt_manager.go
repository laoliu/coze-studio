// Copyright 2025 Coze Studio. All rights reserved.

package prompts

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed *.md
var promptFiles embed.FS

// PromptTemplate 提示词模板
type PromptTemplate struct {
	name     string
	template *template.Template
}

// PromptManager 提示词管理器
type PromptManager struct {
	templates map[string]*PromptTemplate
}

// NewPromptManager 创建提示词管理器
func NewPromptManager() *PromptManager {
	return &PromptManager{
		templates: make(map[string]*PromptTemplate),
	}
}

// LoadFromFile 从嵌入文件加载提示词
func (pm *PromptManager) LoadFromFile(filename string) error {
	content, err := promptFiles.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read prompt file %s: %w", filename, err)
	}

	return pm.ParsePrompts(string(content))
}

// ParsePrompts 解析提示词文档
// 支持 Markdown 格式，每个三级标题(###)是一个提示词模板
func (pm *PromptManager) ParsePrompts(content string) error {
	lines := strings.Split(content, "\n")
	var currentName string
	var currentPrompt strings.Builder
	inCodeBlock := false

	for _, line := range lines {
		// 检测代码块
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// 代码块结束，保存提示词
				if currentName != "" {
					promptText := strings.TrimSpace(currentPrompt.String())
					if err := pm.RegisterTemplate(currentName, promptText); err != nil {
						return err
					}
					currentName = ""
					currentPrompt.Reset()
				}
				inCodeBlock = false
			} else {
				inCodeBlock = true
			}
			continue
		}

		// 检测新的提示词标题
		if strings.HasPrefix(line, "### ") && !inCodeBlock {
			// 保存之前的提示词
			if currentName != "" && currentPrompt.Len() > 0 {
				promptText := strings.TrimSpace(currentPrompt.String())
				if err := pm.RegisterTemplate(currentName, promptText); err != nil {
					return err
				}
				currentPrompt.Reset()
			}
			currentName = strings.TrimSpace(strings.TrimPrefix(line, "### "))
			continue
		}

		// 收集提示词内容
		if inCodeBlock {
			currentPrompt.WriteString(line)
			currentPrompt.WriteString("\n")
		}
	}

	// 保存最后一个提示词
	if currentName != "" && currentPrompt.Len() > 0 {
		promptText := strings.TrimSpace(currentPrompt.String())
		if err := pm.RegisterTemplate(currentName, promptText); err != nil {
			return err
		}
	}

	return nil
}

// RegisterTemplate 注册提示词模板
func (pm *PromptManager) RegisterTemplate(name, promptText string) error {
	tmpl, err := template.New(name).Parse(promptText)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	pm.templates[name] = &PromptTemplate{
		name:     name,
		template: tmpl,
	}

	return nil
}

// Render 渲染提示词
func (pm *PromptManager) Render(name string, data interface{}) (string, error) {
	tmpl, ok := pm.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	if err := tmpl.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}

	return buf.String(), nil
}

// GetTemplate 获取模板
func (pm *PromptManager) GetTemplate(name string) (*PromptTemplate, bool) {
	tmpl, ok := pm.templates[name]
	return tmpl, ok
}

// ListTemplates 列出所有模板名称
func (pm *PromptManager) ListTemplates() []string {
	names := make([]string, 0, len(pm.templates))
	for name := range pm.templates {
		names = append(names, name)
	}
	return names
}

// K12PromptData K12 提示词数据结构
type K12PromptData struct {
	Topic            string
	Domain           string
	Grade            string
	Duration         int
	ActivityType     string
	Objectives       []ObjectiveData
	Contents         []ContentData
	ObjectiveSummary string
	Content          string
}

// ObjectiveData 目标数据
type ObjectiveData struct {
	Objective  string
	Level      string
	Category   string
	Assessment string
}

// ContentData 内容数据
type ContentData struct {
	Title       string
	Description string
	Type        string
	URL         string
}

// K12 提示词模板名称常量
const (
	// 学习目标生成
	PromptObjectiveConcept    = "概念理解活动"
	PromptObjectiveExperiment = "实验探究活动"

	// 叙事生成
	PromptNarrativeConcept = "概念教学叙事"
	PromptNarrativeProject = "项目式学习叙事"

	// 其他
	PromptContentSearch = "内容搜索查询生成"
	PromptQualityAssess = "质量评估"
)
