# 工作流自动生成器 (Workflow Generator)

## 概述

这是 Coze Studio 的工作流自动生成功能模块，允许用户通过自然语言描述需求，自动生成完整的工作流结构。

## 目录结构

```
backend/domain/workflow/internal/generator/
├── types.go                    # 数据结构定义
├── prompt_template.go          # LLM Prompt 模板
├── template_repository.go      # 工作流模板仓库
├── templates.yaml              # 最佳实践模板配置
├── auto_layout.go              # 自动布局引擎
├── workflow_generator.go       # 核心生成引擎
└── README.md                   # 本文档
```

## 核心组件

### 1. WorkflowGenerator (workflow_generator.go)

核心生成引擎，协调整个生成流程。

**主要方法**：
- `GenerateWorkflow()` - 根据需求生成完整工作流

**生成流程**：
```
用户需求 
  ↓
意图分析 (analyzeIntent)
  ↓
模板匹配 (FindSimilarTemplates)
  ↓
LLM 生成工作流结构 (llmGenerateWorkflow)
  ↓
节点配置 (configureNodes)
  ↓
自动布局 (AutoLayout)
  ↓
生成说明 (generateExplanations)
  ↓
返回结果
```

### 2. PromptTemplate (prompt_template.go)

管理 LLM Prompt 模板。

**包含两个主要 Prompt**：
- 意图分析 Prompt - 分析用户需求意图
- 工作流生成 Prompt - 生成完整工作流结构

### 3. TemplateRepository (template_repository.go)

工作流模板仓库，加载和管理最佳实践模板。

**支持的模板类型**：
- content_generation - 内容生成
- data_extraction - 数据提取
- rag_workflow - RAG 知识问答
- batch_processing - 批量处理
- api_integration - API 集成
- data_transformation - 数据转换
- decision_automation - 自动决策

### 4. AutoLayoutEngine (auto_layout.go)

自动布局引擎，计算节点位置。

**算法**：
1. 构建有向图
2. 拓扑排序计算层级
3. 层次布局分配位置

### 5. Types (types.go)

定义所有数据结构：
- 请求/响应类型
- 工作流结构
- 节点配置
- 意图分析结果
等

## 使用示例

```go
package main

import (
    "context"
    
    "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/generator"
)

func main() {
    ctx := context.Background()
    
    // 创建 LLM 模型（需要先配置）
    llmModel := createLLMModel()
    
    // 创建生成器
    gen, err := generator.NewWorkflowGenerator(llmModel)
    if err != nil {
        panic(err)
    }
    
    // 生成工作流
    req := &generator.WorkflowGenerationRequest{
        UserRequirement: "创建一个课程大纲生成器，用户输入课程主题，系统生成包含章节的详细大纲",
        Language:        "zh-CN",
        Constraints: &generator.Constraints{
            MaxNodes: 10,
        },
    }
    
    resp, err := gen.GenerateWorkflow(ctx, req)
    if err != nil {
        panic(err)
    }
    
    // 使用生成的工作流
    fmt.Printf("Generated workflow with %d nodes\n", len(resp.Nodes))
    fmt.Printf("Confidence: %.2f\n", resp.Confidence)
}
```

## API 集成

需要在 `backend/api` 中添加 API 路由：

```go
// POST /api/workflow/generate
type GenerateWorkflowRequest struct {
    UserRequirement string       `json:"user_requirement"`
    Constraints     *Constraints `json:"constraints"`
}

type GenerateWorkflowResponse struct {
    WorkflowID       string            `json:"workflow_id"`
    Nodes            []*vo.Node        `json:"nodes"`
    Edges            []*vo.Edge        `json:"edges"`
    Explanation      string            `json:"explanation"`
    NodeExplanations map[string]string `json:"node_explanations"`
    Confidence       float64           `json:"confidence"`
}
```

## 模板配置

编辑 `templates.yaml` 可以添加或修改工作流模板：

```yaml
templates:
  - id: custom_template
    name: "自定义模板"
    description: "描述"
    tags:
      - tag1
      - tag2
    use_cases:
      - "用例1"
      - "用例2"
    pattern:
      nodes:
        - type: Start
          name: "开始"
        - type: LLM
          name: "处理"
        - type: End
          name: "结束"
```

## 配置要求

1. **LLM 模型**：需要配置支持 JSON 输出的 LLM 模型（如 GPT-4）
2. **温度参数**：
   - 意图分析：0.3（低温度，稳定输出）
   - 工作流生成：0.5（中等温度，平衡创新和稳定）

## 性能指标

- **意图分析时间**：< 2s
- **工作流生成时间**：< 5s
- **总体生成时间**：< 10s
- **生成成功率**：目标 ≥ 85%
- **节点准确性**：目标 ≥ 80%

## 错误处理

主要错误类型：
1. **意图分析失败** - 需求描述过于模糊
2. **LLM 调用失败** - 网络或 API 问题
3. **JSON 解析失败** - LLM 输出格式不正确

建议使用重试机制和降级策略。

## 后续优化

1. **验证引擎** - 检查生成的工作流合法性
2. **成本估算** - 计算工作流执行成本
3. **性能优化** - 缓存模板和 Prompt
4. **A/B 测试** - 优化 Prompt 模板
5. **用户反馈** - 基于反馈改进生成质量

## 相关文档

- [NODE_RECOMMENDATION_SYSTEM.md](../../../../../docs_new/NODE_RECOMMENDATION_SYSTEM.md) - 完整技术设计文档
- [WORKFLOW_AUTO_GENERATION_SUMMARY.md](../../../../../docs_new/WORKFLOW_AUTO_GENERATION_SUMMARY.md) - 功能说明
- [WORKFLOW_AUTO_GENERATION_QUICK_GUIDE.md](../../../../../docs_new/WORKFLOW_AUTO_GENERATION_QUICK_GUIDE.md) - 快速指南

## 许可证

Apache License 2.0
