package generator

import (
	"embed"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"gopkg.in/yaml.v3"
)

//go:embed templates.yaml
var templatesFS embed.FS

// TemplateRepository 模板仓库
type TemplateRepository struct {
	templates []*WorkflowTemplate
}

// NewTemplateRepository 创建模板仓库
func NewTemplateRepository() (*TemplateRepository, error) {
	repo := &TemplateRepository{}
	if err := repo.loadTemplates(); err != nil {
		return nil, err
	}
	return repo, nil
}

// loadTemplates 加载模板配置
func (r *TemplateRepository) loadTemplates() error {
	data, err := templatesFS.ReadFile("templates.yaml")
	if err != nil {
		return err
	}

	var config struct {
		Templates []*struct {
			ID          string   `yaml:"id"`
			Name        string   `yaml:"name"`
			Description string   `yaml:"description"`
			Tags        []string `yaml:"tags"`
			UseCases    []string `yaml:"use_cases"`
			Pattern     struct {
				Nodes []struct {
					Type     string                 `yaml:"type"`
					Name     string                 `yaml:"name"`
					Function string                 `yaml:"function"`
					Config   map[string]interface{} `yaml:"config"`
				} `yaml:"nodes"`
			} `yaml:"pattern"`
		} `yaml:"templates"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// 转换为内部类型
	r.templates = make([]*WorkflowTemplate, 0, len(config.Templates))
	for _, tmpl := range config.Templates {
		nodes := make([]*TemplateNode, 0, len(tmpl.Pattern.Nodes))
		for _, node := range tmpl.Pattern.Nodes {
			nodes = append(nodes, &TemplateNode{
				Type:     entity.NodeType(node.Type),
				Name:     node.Name,
				Function: node.Function,
				Config:   node.Config,
			})
		}

		r.templates = append(r.templates, &WorkflowTemplate{
			ID:          tmpl.ID,
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Tags:        tmpl.Tags,
			UseCases:    tmpl.UseCases,
			Pattern: &TemplatePattern{
				Nodes: nodes,
			},
		})
	}

	return nil
}

// FindSimilarTemplates 查找相似模板
func (r *TemplateRepository) FindSimilarTemplates(workflowType string, limit int) []*WorkflowTemplate {
	results := make([]*WorkflowTemplate, 0)

	// 精确匹配工作流类型
	for _, tmpl := range r.templates {
		if tmpl.ID == workflowType {
			results = append(results, tmpl)
			if len(results) >= limit {
				return results
			}
		}
	}

	// 标签匹配
	typeKeywords := strings.Split(workflowType, "_")
	for _, tmpl := range r.templates {
		if len(results) >= limit {
			break
		}

		// 检查是否已添加
		found := false
		for _, existing := range results {
			if existing.ID == tmpl.ID {
				found = true
				break
			}
		}
		if found {
			continue
		}

		// 标签匹配度检查
		for _, tag := range tmpl.Tags {
			for _, keyword := range typeKeywords {
				if strings.Contains(tag, keyword) {
					results = append(results, tmpl)
					break
				}
			}
		}
	}

	// 如果结果不足，添加通用模板
	if len(results) < limit {
		for _, tmpl := range r.templates {
			if len(results) >= limit {
				break
			}

			found := false
			for _, existing := range results {
				if existing.ID == tmpl.ID {
					found = true
					break
				}
			}
			if !found {
				results = append(results, tmpl)
			}
		}
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

// GetAllTemplates 获取所有模板
func (r *TemplateRepository) GetAllTemplates() []*WorkflowTemplate {
	return r.templates
}

// GetTemplateByID 根据 ID 获取模板
func (r *TemplateRepository) GetTemplateByID(id string) *WorkflowTemplate {
	for _, tmpl := range r.templates {
		if tmpl.ID == id {
			return tmpl
		}
	}
	return nil
}
