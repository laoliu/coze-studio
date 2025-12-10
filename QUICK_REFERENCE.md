# 新架构快速参考指南

## 🎯 核心概念速查

### 三层架构

| 层级 | 名称 | 职责 | 示例 |
|------|------|------|------|
| 🎓 | **领域适配器层** | 领域规则翻译 | K12Adapter, ArtHistoryAdapter |
| 🧠 | **AI能力层** | 通用AI算法 | ObjectiveGenerator, ContentDiscovery |
| 🔧 | **组件层** | 可复用功能 | MCPTool, VisualizationComponent |

## 📦 模块路径

```
backend/domain/
├── adapter/          # 领域适配器
│   ├── entity/       # 实体定义
│   ├── service/      # 服务层（管理器）
│   ├── examples/     # 示例实现
│   └── interface.go  # 接口定义
│
├── capability/       # AI能力
│   ├── entity/
│   ├── service/
│   └── interface.go
│
└── component/        # 组件
    ├── entity/
    ├── service/
    └── interface.go
```

## 🔑 关键接口

### DomainAdapter（领域适配器）

```go
type DomainAdapter interface {
    GetInfo() *AdapterInfo
    ParseRequest(*UserInput) (*RequestContext, error)
    ValidateContext(*RequestContext) error
    GenerateLearningObjectives(*RequestContext) ([]*LearningObjective, error)
    DiscoverContent(*RequestContext, []*LearningObjective) ([]*Content, error)
    CustomizeWorkflow(*RequestContext) (*WorkflowConfig, error)
    FormatOutput(*WorkflowResult) (interface{}, error)
    ValidateOutput(interface{}) (*QualityReport, error)
}
```

### AICapability（AI能力）

```go
type AICapability interface {
    GetCapabilityInfo() *CapabilityInfo
    Execute(context.Context, *CapabilityInput) (*CapabilityOutput, error)
    Validate(*CapabilityInput) error
    GetConfig() map[string]interface{}
}
```

### ComponentExecutor（组件执行器）

```go
type ComponentExecutor interface {
    GetComponentInfo() *ComponentInfo
    Execute(context.Context, *ComponentInput) (*ComponentOutput, error)
    Validate(*ComponentInput) error
    GetConfig() map[string]interface{}
}
```

## 💻 代码示例

### 1. 注册和使用适配器

```go
import (
    "github.com/coze-studio/backend/domain/adapter/service"
    "github.com/coze-studio/backend/domain/adapter/entity"
)

// 创建管理器
adapterMgr := service.NewAdapterManager()

// 注册适配器
k12Adapter := NewK12Adapter()
adapterMgr.RegisterAdapter(k12Adapter)

// 使用适配器
userInput := &entity.UserInput{
    Topic:        "氧化还原反应",
    Domain:       "化学",
    Grade:        "高一",
    ActivityType: entity.ActivityTypeConcept,
}

ctx := context.Background()
reqCtx, _ := adapterMgr.ParseRequest(ctx, "k12_education", userInput)
objectives, _ := adapterMgr.GenerateLearningObjectives(ctx, "k12_education", reqCtx)
```

### 2. 注册和使用AI能力

```go
import "github.com/coze-studio/backend/domain/capability/service"

// 创建管理器
capMgr := service.NewCapabilityManager()

// 注册能力
objGen := NewObjectiveGeneratorCapability()
capMgr.RegisterCapability(objGen)

// 使用能力
params := &capability.ObjectiveParams{
    Topic:    "氧化还原反应",
    Grade:    "grade_10",
    Domain:   "chemistry",
    Duration: 45,
}
objectives, _ := capMgr.ExecuteObjectiveGeneration(ctx, params)
```

### 3. 注册和使用组件

```go
import "github.com/coze-studio/backend/domain/component/service"

// 创建管理器
compMgr := service.NewComponentManager()

// 注册组件
mcpTool := NewVirtualLabTool()
compMgr.RegisterComponent(mcpTool)

// 使用组件
params := map[string]any{
    "experiment": "oxidation_reduction",
    "chemicals":  []string{"KMnO4", "H2O2"},
}
result, _ := compMgr.InvokeMCPTool(ctx, "virtual_lab", params)
```

## 📊 实体类型速查

### AdapterType（适配器类型）

```go
const (
    AdapterTypeK12         = "k12_education"
    AdapterTypeArtHistory  = "art_history"
    AdapterTypeVocational  = "vocational_training"
    AdapterTypeCustom      = "custom"
)
```

### ActivityType（活动类型）

```go
const (
    ActivityTypeConcept    = "concept"      // 概念理解
    ActivityTypeExperiment = "experiment"   // 实验探究
    ActivityTypeProblem    = "problem"      // 问题解决
    ActivityTypeProject    = "project"      // 项目制学习
    ActivityTypeDiscussion = "discussion"   // 讨论辩论
    ActivityTypeCreation   = "creation"     // 创作表达
    ActivityTypeReflection = "reflection"   // 反思总结
)
```

### CapabilityType（能力类型）

```go
const (
    CapabilityTypeIntentRecognition   = "intent_recognition"
    CapabilityTypeObjectiveGeneration = "objective_generation"
    CapabilityTypeContentDiscovery    = "content_discovery"
    CapabilityTypeNarrativeGeneration = "narrative_generation"
    CapabilityTypeQualityAssessment   = "quality_assessment"
    CapabilityTypeMultimodalGeneration = "multimodal_generation"
)
```

## 🔄 典型工作流

```
用户输入 (UserInput)
    ↓
适配器解析 (ParseRequest)
    ↓
上下文验证 (ValidateContext)
    ↓
生成学习目标 (GenerateLearningObjectives)
    → AI能力: ObjectiveGenerator
    ↓
发现内容 (DiscoverContent)
    → AI能力: ContentDiscovery
    ↓
定制工作流 (CustomizeWorkflow)
    ↓
执行工作流 (Workflow Engine)
    → 多个AI能力 + 组件
    ↓
格式化输出 (FormatOutput)
    ↓
质量验证 (ValidateOutput)
    → AI能力: QualityAssessor
    ↓
最终结果
```

## 🛠️ 实现新适配器的步骤

1. **定义适配器结构**
```go
type MyAdapter struct {
    info   *entity.AdapterInfo
    config *MyConfig
}
```

2. **实现构造函数**
```go
func NewMyAdapter() *MyAdapter {
    return &MyAdapter{
        info: &entity.AdapterInfo{
            AdapterID: "my_adapter",
            // ... 其他信息
        },
    }
}
```

3. **实现所有接口方法**
```go
func (a *MyAdapter) ParseRequest(input *entity.UserInput) (*entity.RequestContext, error) {
    // 实现逻辑
}
// ... 其他方法
```

4. **确保实现接口**
```go
var _ adapter.DomainAdapter = (*MyAdapter)(nil)
```

5. **注册适配器**
```go
adapterMgr := service.NewAdapterManager()
myAdapter := NewMyAdapter()
adapterMgr.RegisterAdapter(myAdapter)
```

## 📚 相关文档

- 详细实施总结: `IMPLEMENTATION_SUMMARY.md`
- 域模型说明: `backend/domain/README.md`
- 架构设计文档: `docs_new/ARCHITECTURE_CORE_CONCEPTS.md`
- Coze集成方案: `docs_new/COZE_INTEGRATION_PLAN.md`

## ❓ 常见问题

### Q: 如何选择使用哪个适配器？

A: 使用 `AdapterManager.AutoSelectAdapter()` 方法，或根据 `Domain` 字段手动选择。

### Q: AI能力和组件有什么区别？

A: AI能力专注于算法（如LLM调用），组件提供具体功能（如导出、可视化）。

### Q: 如何添加新的活动类型？

A: 在 `entity.ActivityType` 中添加新常量，然后在适配器的 `Capabilities` 中声明支持。

### Q: 适配器之间可以互相调用吗？

A: 可以通过依赖声明，但不推荐。应该通过共享AI能力来实现复用。

---

**快速帮助**: 运行 `go run backend/domain/adapter/examples/*.go` 查看完整示例
