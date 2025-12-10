# 工作流节点扩展 (Workflow Node Extensions)

本目录包含扩展的工作流节点，用于集成新的领域适配器、AI能力和组件系统。

## 目录结构

```
nodes/
├── adapter/                    # 适配器节点
│   ├── adapter_invoke.go      # 适配器调用节点实现
│   └── config_schema.json     # 节点配置JSON Schema
├── capability/                 # AI能力节点
│   ├── capability_invoke.go   # 通用能力调用节点
│   ├── specialized_nodes.go   # 专用能力节点（目标生成、内容发现等）
│   └── config_schema.json     # 节点配置JSON Schema
├── component/                  # 组件节点
│   ├── component_invoke.go    # 组件调用节点（MCP工具、可视化、导出等）
│   └── config_schema.json     # 节点配置JSON Schema
├── new_node_metas.go          # 新节点的元数据定义
├── register_new_nodes.go      # 节点注册函数
└── integration_test.go        # 集成测试
```

## 节点类型

### 1. 适配器节点 (Adapter Nodes)

**适配器调用节点 (adapter_invoke)**

用于调用领域适配器的各种方法，实现领域特定的逻辑处理。

支持的方法：
- `parse_request`: 解析用户请求，提取领域特定信息
- `generate_objectives`: 生成领域特定的学习目标
- `discover_content`: 发现和筛选相关内容
- `customize_workflow`: 定制工作流配置
- `format_output`: 格式化输出结果
- `validate_output`: 验证输出质量

配置示例：
```json
{
  "adapter_id": "k12_adapter",
  "method": "parse_request",
  "input_mapping": {
    "user_input": "{{workflow.input.query}}"
  },
  "output_mapping": {
    "parsed_data": "{{node.output}}"
  },
  "timeout": 30000
}
```

### 2. AI能力节点 (Capability Nodes)

#### 通用能力调用节点 (capability_invoke)

用于调用任何注册的AI能力。

配置示例：
```json
{
  "capability_id": "intent_recognizer",
  "capability_type": "intent_recognition",
  "input": {
    "text": "我想学习编程"
  },
  "model_config": {
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 2000
  }
}
```

#### 专用能力节点

**目标生成节点 (objective_generator)**

生成SMART学习目标。

配置示例：
```json
{
  "domain": "K12数学",
  "grade": "七年级",
  "subject": "二次函数",
  "student_level": "中等"
}
```

**内容发现节点 (content_discovery)**

从知识库或外部源发现相关内容。

配置示例：
```json
{
  "query": "二次函数的基础知识",
  "sources": ["knowledge_base", "wikipedia"],
  "max_results": 10
}
```

**故事化叙述生成节点 (narrative_generator)**

生成引人入胜的故事化叙述内容。

配置示例：
```json
{
  "topic": "二次函数的应用",
  "style": "科幻故事",
  "length": "中等"
}
```

**质量评估节点 (quality_assessor)**

评估内容质量并生成详细报告。

配置示例：
```json
{
  "content": "待评估的内容...",
  "criteria": {
    "教育性": 0.8,
    "趣味性": 0.7,
    "准确性": 0.9
  }
}
```

### 3. 组件节点 (Component Nodes)

#### 通用组件调用节点 (component_invoke)

用于调用任何注册的组件。

#### MCP工具节点 (mcp_tool)

调用Model Context Protocol工具。

配置示例：
```json
{
  "component_id": "web_search",
  "component_type": "mcp_tool",
  "input": {
    "query": "二次函数应用案例",
    "max_results": 5
  }
}
```

#### 可视化节点 (visualization)

将数据渲染为可视化图表或界面。

配置示例：
```json
{
  "component_id": "chart_renderer",
  "component_type": "visualization",
  "input": {
    "data": [...],
    "chart_type": "line",
    "title": "函数图像"
  }
}
```

#### 导出节点 (export)

将内容导出为指定格式。

配置示例：
```json
{
  "component_id": "pdf_exporter",
  "component_type": "export",
  "input": {
    "content": "要导出的内容...",
    "title": "学习材料"
  },
  "config": {
    "format": "PDF",
    "layout": "A4",
    "include_toc": true
  }
}
```

## 使用方法

### 1. 注册节点

在应用启动时调用注册函数：

```go
import "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"

func init() {
    // 注册所有新节点
    nodes.RegisterAllNewNodes()
}
```

### 2. 在工作流中使用

在工作流定义中使用节点类型键：

```json
{
  "nodes": [
    {
      "id": "node_1",
      "type": "adapter_invoke",
      "config": {
        "adapter_id": "k12_adapter",
        "method": "parse_request"
      }
    },
    {
      "id": "node_2",
      "type": "objective_generator",
      "config": {
        "domain": "{{node_1.output.domain}}",
        "grade": "{{node_1.output.grade}}"
      }
    },
    {
      "id": "node_3",
      "type": "export",
      "config": {
        "component_id": "pdf_exporter",
        "input": {
          "content": "{{node_2.output.objectives}}"
        }
      }
    }
  ]
}
```

### 3. 编程方式创建节点

```go
import (
    "context"
    workflowadapter "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/adapter"
)

ctx := context.Background()
config := map[string]interface{}{
    "adapter_id": "k12_adapter",
    "method":     "parse_request",
}

node, err := workflowadapter.NewAdapterInvokeNode(ctx, config)
if err != nil {
    // 处理错误
}

// 执行节点
input := map[string]interface{}{
    "user_input": "我想学习七年级数学",
}
output, err := node.Invoke(ctx, input)
```

## 完整工作流示例

以下是一个K12教学内容生成的完整工作流示例：

```json
{
  "id": "k12_content_generation_workflow",
  "name": "K12教学内容生成工作流",
  "nodes": [
    {
      "id": "parse_request",
      "type": "adapter_invoke",
      "config": {
        "adapter_id": "k12_adapter",
        "method": "parse_request",
        "input_mapping": {
          "user_input": "{{workflow.input.query}}"
        }
      }
    },
    {
      "id": "generate_objectives",
      "type": "objective_generator",
      "config": {
        "domain": "{{parse_request.output.domain}}",
        "grade": "{{parse_request.output.grade}}",
        "subject": "{{parse_request.output.subject}}"
      }
    },
    {
      "id": "discover_content",
      "type": "content_discovery",
      "config": {
        "query": "{{parse_request.output.subject}}",
        "grade": "{{parse_request.output.grade}}"
      }
    },
    {
      "id": "generate_narrative",
      "type": "narrative_generator",
      "config": {
        "objectives": "{{generate_objectives.output}}",
        "content": "{{discover_content.output}}",
        "style": "{{parse_request.output.narrative_style}}"
      }
    },
    {
      "id": "assess_quality",
      "type": "quality_assessor",
      "config": {
        "content": "{{generate_narrative.output}}",
        "objectives": "{{generate_objectives.output}}"
      }
    },
    {
      "id": "visualize",
      "type": "visualization",
      "config": {
        "component_id": "interactive_visualizer",
        "input": {
          "content": "{{generate_narrative.output}}",
          "type": "interactive_lesson"
        }
      }
    },
    {
      "id": "export",
      "type": "export",
      "config": {
        "component_id": "multi_format_exporter",
        "input": {
          "content": "{{generate_narrative.output}}",
          "visualization": "{{visualize.output}}"
        },
        "config": {
          "formats": ["PDF", "HTML", "PPTX"]
        }
      }
    }
  ],
  "edges": [
    {"from": "parse_request", "to": "generate_objectives"},
    {"from": "parse_request", "to": "discover_content"},
    {"from": "generate_objectives", "to": "generate_narrative"},
    {"from": "discover_content", "to": "generate_narrative"},
    {"from": "generate_narrative", "to": "assess_quality"},
    {"from": "assess_quality", "to": "visualize"},
    {"from": "visualize", "to": "export"}
  ]
}
```

## 测试

运行集成测试：

```bash
cd backend/domain/workflow/internal/nodes
go test -v ./...
```

运行性能测试：

```bash
go test -bench=. -benchmem
```

## 扩展指南

### 添加新的适配器方法

1. 在对应的适配器接口中添加新方法
2. 在 `adapter_invoke.go` 的 switch 语句中添加新的 case
3. 更新 `config_schema.json` 中的 method enum

### 添加新的AI能力类型

1. 在 `backend/domain/capability/interface.go` 中定义新能力接口
2. 在 `specialized_nodes.go` 中创建专用节点
3. 添加对应的 NodeAdaptor
4. 更新元数据定义

### 添加新的组件类型

1. 在 `backend/domain/component/interface.go` 中定义新组件接口
2. 在 `component_invoke.go` 中添加新的节点类型
3. 创建 NodeAdaptor
4. 更新配置 Schema

## 注意事项

1. **节点ID分配**: 新节点的元数据ID从1001开始，避免与现有节点冲突
2. **超时设置**: 根据实际情况设置合理的超时时间，AI能力节点建议60秒
3. **错误处理**: 所有节点都应妥善处理错误并返回有意义的错误信息
4. **并发安全**: 节点实现应该是无状态的，支持并发执行
5. **资源清理**: 确保节点执行完成后正确释放资源

## 相关文档

- [领域适配器设计](../../../adapter/README.md)
- [AI能力系统](../../../capability/README.md)
- [组件系统](../../../component/README.md)
- [工作流引擎](../../README.md)
