# 节点推荐系统

## 📖 简介

节点推荐系统为 Coze Studio 工作流编辑器提供智能节点推荐功能，帮助用户快速找到最合适的后续节点。

## 🎯 功能特性

### MVP 阶段（v1.0）

- ✅ **基于规则的推荐**：通过 YAML 配置的规则引擎，提供快速、可控的节点推荐
- ✅ **多维度评分**：综合考虑节点类型、数据流、工作流上下文等因素
- ✅ **智能排序**：基于得分和最佳实践自动排序推荐结果
- ✅ **可解释性**：为每个推荐提供清晰的理由

### 未来规划

- ⏳ **LLM 智能推荐**：集成大模型进行上下文理解和智能推理
- ⏳ **统计分析推荐**：基于历史数据和用户行为优化推荐
- ⏳ **个性化推荐**：学习用户偏好，提供定制化推荐
- ⏳ **实时学习**：根据用户反馈持续优化推荐质量

## 📂 目录结构

```
recommendation/
├── doc.go              # 包文档
├── README.md           # 使用说明
├── types.go            # 类型定义
├── engine.go           # 推荐引擎主入口
├── rule_engine.go      # 规则引擎实现
├── scorer.go           # 评分和排序
├── rules.yaml          # 推荐规则配置
├── prompts.yaml        # LLM Prompt 模板（未来使用）
└── llm_strategy.go     # LLM 策略（待实现）
```

## 🚀 快速开始

### 1. 创建推荐引擎

```go
import "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/recommendation"

// 使用默认配置
engine, err := recommendation.NewEngine(nil)
if err != nil {
    log.Fatal(err)
}
```

### 2. 构建推荐请求

```go
req := &recommendation.RecommendRequest{
    WorkflowID:     "workflow-123",
    SourceNodeID:   "node-456",
    SourceNodeType: entity.NodeTypeLLM,
    SourceOutputs: map[string]*vo.TypeInfo{
        "output": {Type: entity.TypeString},
    },
    OutputFormat: "json",
    WorkflowContext: &recommendation.WorkflowContext{
        HasDatabaseNode: true,
        TotalNodes:      5,
    },
    Limit:         5,
    IncludeReason: true,
}
```

### 3. 执行推荐

```go
resp, err := engine.Recommend(context.Background(), req)
if err != nil {
    return err
}

for _, rec := range resp.Recommendations {
    fmt.Printf("节点: %s (%.0f%%)\n", rec.DisplayName, rec.Score*100)
    fmt.Printf("理由: %s\n", rec.Reason)
    fmt.Printf("分类: %s\n\n", rec.Category)
}
```

## 📝 规则配置

推荐规则定义在 `rules.yaml` 文件中。每条规则包含：

- **条件**（Condition）：定义规则何时匹配
- **推荐**（Recommendations）：符合条件时推荐哪些节点

### 示例规则

```yaml
recommendation_rules:
  - name: llm_json_to_code
    priority: high
    condition:
      source_node_type: LLM
      output_format: json
    recommendations:
      - node_type: Code
        score: 0.9
        reason: "解析 LLM 返回的 JSON 数据"
        category: "数据处理"
```

### 条件字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `source_node_type` | NodeType | 源节点类型 |
| `source_node_types` | []NodeType | 源节点类型列表（多选一）|
| `output_format` | string | 输出格式（json/text/markdown）|
| `output_has_array` | bool | 输出是否包含数组 |
| `workflow_has_database` | bool | 工作流是否有数据库节点 |
| `inside_loop` | bool | 是否在循环内部 |
| `is_chat_workflow` | bool | 是否是聊天工作流 |
| `always_match` | bool | 通配符（总是匹配）|

## 🔧 API 接口

### 推荐节点

**请求：** `POST /api/workflow/recommend/nodes`

```json
{
  "workflow_id": "workflow-123",
  "source_node_id": "node-456",
  "source_node_type": "LLM",
  "output_format": "json",
  "limit": 5,
  "include_reason": true
}
```

**响应：**

```json
{
  "recommendations": [
    {
      "node_type": "Code",
      "display_name": "代码节点",
      "score": 0.9,
      "reason": "解析 LLM 返回的 JSON 数据",
      "category": "数据处理",
      "icon": "icon-code.jpg",
      "source": "rule_engine"
    }
  ],
  "meta": {
    "total_candidates": 10,
    "filtered_count": 5,
    "recommendation_id": "rec-123",
    "execution_time_ms": 45
  }
}
```

### 提交反馈

**请求：** `POST /api/workflow/recommend/feedback`

```json
{
  "workflow_id": "workflow-123",
  "recommendation_id": "rec-123",
  "source_node_id": "node-456",
  "selected_node_type": "Code",
  "is_useful": true
}
```

## 📊 评分机制

推荐得分由多个因素决定：

1. **规则得分**：规则配置中定义的基础得分（0-1）
2. **策略权重**：不同策略的权重（规则 40%、LLM 40%、统计 20%）
3. **提升因子**：
   - 常用节点：+0.05
   - 匹配用户偏好：+0.1
   - 最佳实践：+0.08

最终得分 = (规则得分 × 策略权重) + 提升因子

## 🧪 测试

```bash
# 运行单元测试
go test ./backend/domain/workflow/internal/recommendation/...

# 运行测试并查看覆盖率
go test -cover ./backend/domain/workflow/internal/recommendation/...
```

## 📈 性能指标

| 指标 | 目标值 | 当前值 |
|------|--------|--------|
| 推荐准确率 | > 60% | TBD |
| 响应时间 | < 500ms | ~50ms |
| 覆盖率 | > 80% | TBD |

## 🛠️ 扩展开发

### 添加新规则

1. 编辑 `rules.yaml`
2. 添加规则定义
3. 重启服务（规则会自动加载）

### 实现新策略

1. 实现 `Strategy` 接口：

```go
type MyStrategy struct {
    name   string
    weight float64
}

func (s *MyStrategy) Name() string { return s.name }
func (s *MyStrategy) Weight() float64 { return s.weight }
func (s *MyStrategy) Recommend(req *RecommendRequest) ([]RecommendedNode, error) {
    // 实现推荐逻辑
    return recommendations, nil
}
```

2. 在 `engine.go` 中注册策略

## 📚 参考文档

- [节点推荐系统设计文档](../../../../../docs_new/NODE_RECOMMENDATION_SYSTEM.md)
- [工作流节点文档](../nodes/README.md)
- [API 设计文档](../../../../../docs_new/api_design.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

Apache License 2.0
