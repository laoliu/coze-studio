package generator

import (
	"bytes"
	"encoding/json"
	"text/template"
)

// PromptTemplate Prompt 模板管理器
type PromptTemplate struct{}

// NewPromptTemplate 创建 Prompt 模板管理器
func NewPromptTemplate() *PromptTemplate {
	return &PromptTemplate{}
}

// intentAnalysisPromptTemplate 意图分析 Prompt 模板
const intentAnalysisPromptTemplate = `你是一个需求分析专家，请分析用户的工作流需求并提取关键信息。

## 用户需求
{{.UserRequirement}}

## 分析任务
请识别以下信息：
1. **核心目标**：用户想要实现什么功能？
2. **输入数据**：工作流需要什么输入？（类型、来源）
3. **输出结果**：期望得到什么输出？（格式、内容）
4. **关键步骤**：需要经过哪些主要处理步骤？
5. **约束条件**：是否有性能、成本、质量要求？
6. **工作流类型**：属于哪种模式？
   - content_generation (内容生成)
   - data_extraction (数据提取)
   - rag_workflow (知识问答)
   - batch_processing (批量处理)
   - api_integration (API 集成)
   - data_transformation (数据转换)
   - decision_automation (自动决策)

## 输出格式（JSON）
{
  "intent": {
    "core_objective": "生成结构化的课程大纲",
    "workflow_type": "content_generation",
    "complexity": "medium",
    "confidence": 0.95
  },
  "inputs": [
    {"name": "course_topic", "type": "string", "description": "课程主题", "required": true}
  ],
  "outputs": [
    {"name": "course_outline", "type": "json", "description": "包含章节的课程大纲", "required": true}
  ],
  "key_steps": [
    "接收课程主题",
    "使用 LLM 生成初始大纲",
    "解析和验证大纲结构",
    "输出最终结果"
  ],
  "constraints": {
    "quality": "high",
    "cost": "medium",
    "performance": "standard"
  },
  "suggested_nodes": [
    {"type": "LLM", "reason": "需要生成创意内容"},
    {"type": "Code", "reason": "需要解析和验证数据"}
  ]
}

请只返回 JSON，不要包含其他文字。
`

// workflowGenerationPromptTemplate 工作流生成 Prompt 模板
const workflowGenerationPromptTemplate = `你是工作流设计专家，请设计一个完整的工作流方案。

## 用户需求
{{.UserRequirement}}

## 意图分析结果
{{.IntentJSON}}

{{if .Templates}}
## 参考模板
{{range .Templates}}
### {{.Name}}
- 描述: {{.Description}}
- 适用场景: {{range .UseCases}}{{.}}, {{end}}
- 节点序列: {{.NodeSequence}}
{{end}}
{{end}}

## 可用节点类型
1. **Start** - 工作流入口，定义输入参数
   - 用途: 接收外部输入
   - 配置: inputs (参数列表)

2. **LLM** - 调用大语言模型生成内容
   - 用途: 文本生成、理解、转换
   - 配置: model, prompt, temperature, output_format

3. **Code** - 执行 Python 代码处理数据
   - 用途: 数据处理、计算、转换
   - 配置: code (Python 代码)

4. **Plugin** - 调用外部 API 或服务
   - 用途: 集成第三方服务
   - 配置: plugin_id, parameters

5. **Knowledge** - 知识库检索
   - 用途: 从知识库检索相关信息
   - 配置: knowledge_base_id, query, top_k

6. **Database** - 数据库操作
   - 用途: 读写数据库
   - 配置: operation (query/insert/update), table, conditions

7. **If** - 条件判断分支
   - 用途: 根据条件选择不同路径
   - 配置: condition, branches

8. **Loop** - 循环处理数组数据
   - 用途: 批量处理、遍历数组
   - 配置: items_path, loop_variable

9. **End** - 工作流出口
   - 用途: 定义最终输出
   - 配置: outputs (输出字段)

## 设计原则
1. **简洁高效**: 使用最少的节点完成任务
2. **鲁棒性**: 考虑异常处理和边界情况
3. **可维护性**: 节点职责单一，连接清晰
4. **最佳实践**: 遵循常见的工作流模式

## 输出格式（JSON）
{
  "workflow_name": "课程大纲生成器",
  "description": "根据用户输入的课程主题，自动生成结构化的课程大纲",
  "nodes": [
    {
      "id": "node_1",
      "type": "Start",
      "name": "开始",
      "position": {"x": 100, "y": 100},
      "config": {
        "inputs": [
          {"name": "course_topic", "type": "string", "required": true, "description": "课程主题"}
        ]
      },
      "explanation": "接收用户输入的课程主题作为工作流的输入参数"
    },
    {
      "id": "node_2",
      "type": "LLM",
      "name": "生成课程大纲",
      "position": {"x": 100, "y": 250},
      "config": {
        "model": "gpt-4",
        "prompt": "请根据以下课程主题生成详细的课程大纲，包含至少5个章节，每个章节包含标题和要点。\\n\\n课程主题: {{course_topic}}\\n\\n请以 JSON 格式输出，结构如下：\\n{\\\"chapters\\\": [{\\\"title\\\": \\\"章节标题\\\", \\\"points\\\": [\\\"要点1\\\", \\\"要点2\\\"]}]}",
        "output_format": "json",
        "temperature": 0.7
      },
      "explanation": "使用 LLM 根据主题生成结构化的课程大纲"
    },
    {
      "id": "node_3",
      "type": "Code",
      "name": "验证大纲结构",
      "position": {"x": 100, "y": 400},
      "config": {
        "code": "import json\\n\\ndef main(outline_json):\\n    outline = json.loads(outline_json)\\n    chapters = outline.get('chapters', [])\\n    \\n    # 验证章节数量\\n    if len(chapters) < 3:\\n        return {'valid': False, 'error': '章节数量不足'}\\n    \\n    # 验证章节结构\\n    for chapter in chapters:\\n        if 'title' not in chapter or 'points' not in chapter:\\n            return {'valid': False, 'error': '章节结构不完整'}\\n    \\n    return {'valid': True, 'outline': outline}"
      },
      "explanation": "验证生成的大纲是否符合要求"
    },
    {
      "id": "node_4",
      "type": "End",
      "name": "结束",
      "position": {"x": 100, "y": 550},
      "config": {
        "outputs": [
          {"name": "course_outline", "type": "json", "source": "node_3.outline"}
        ]
      },
      "explanation": "输出最终的课程大纲"
    }
  ],
  "edges": [
    {"from": "node_1", "to": "node_2", "output_key": "course_topic", "input_key": "course_topic"},
    {"from": "node_2", "to": "node_3", "output_key": "output", "input_key": "outline_json"},
    {"from": "node_3", "to": "node_4", "output_key": "outline", "input_key": "course_outline"}
  ],
  "overall_explanation": "这个工作流首先接收课程主题，然后使用 LLM 生成初始大纲，接着用 Python 代码验证大纲结构是否符合要求，最后输出验证后的课程大纲。整个流程简洁高效，确保输出质量。",
  "confidence": 0.92
}

请只返回 JSON，不要包含其他文字。确保 JSON 格式正确，所有字符串都正确转义。
`

// BuildIntentAnalysisPrompt 构建意图分析 Prompt
func (pt *PromptTemplate) BuildIntentAnalysisPrompt(requirement string) string {
	tmpl := template.Must(template.New("intent").Parse(intentAnalysisPromptTemplate))
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, map[string]interface{}{
		"UserRequirement": requirement,
	})
	return buf.String()
}

// BuildWorkflowGenerationPrompt 构建工作流生成 Prompt
func (pt *PromptTemplate) BuildWorkflowGenerationPrompt(
	requirement string,
	intent *IntentAnalysisResult,
	templates []*WorkflowTemplate,
) string {
	// 序列化意图分析结果
	intentJSON, _ := json.MarshalIndent(intent, "", "  ")

	// 准备模板数据
	templateData := make([]map[string]interface{}, 0, len(templates))
	for _, tmpl := range templates {
		nodeSeq := make([]string, 0, len(tmpl.Pattern.Nodes))
		for _, node := range tmpl.Pattern.Nodes {
			nodeSeq = append(nodeSeq, string(node.Type))
		}
		templateData = append(templateData, map[string]interface{}{
			"Name":         tmpl.Name,
			"Description":  tmpl.Description,
			"UseCases":     tmpl.UseCases,
			"NodeSequence": nodeSeq,
		})
	}

	tmpl := template.Must(template.New("workflow").Parse(workflowGenerationPromptTemplate))
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, map[string]interface{}{
		"UserRequirement": requirement,
		"IntentJSON":      string(intentJSON),
		"Templates":       templateData,
	})
	return buf.String()
}
