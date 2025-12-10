# 基础领域模型实现说明

**创建时间**: 2025-12-10  
**对应文档**: `/docs_new/ARCHITECTURE_CORE_CONCEPTS.md`, `/docs_new/COZE_INTEGRATION_PLAN.md`

## 📋 概述

本文档说明根据 `docs_new` 中的架构设计，在现有 Coze Studio 工程基础上创建的基础领域模型。

## 🏗️ 已创建的模块

### 1. Adapter 模块（领域适配器）

**位置**: `backend/domain/adapter/`

#### 核心文件

- **`entity/adapter.go`**: 适配器实体定义
  - `Adapter`: 适配器数据模型
  - `AdapterType`: 适配器类型枚举（K12、美术史、职业培训等）
  - `AdapterStatus`: 适配器状态（草稿、激活、停用、归档）
  - `ActivityType`: 活动类型（概念理解、实验探究、问题解决等）
  - `AdapterCapabilities`: 能力声明
  - `AdapterConfig`: 配置信息
  - `AdapterDependencies`: 依赖项

- **`interface.go`**: 适配器接口定义
  - `DomainAdapter`: 所有领域适配器必须实现的接口
    - `ParseRequest()`: 解析用户请求
    - `ValidateContext()`: 验证请求上下文
    - `GenerateLearningObjectives()`: 生成学习目标
    - `DiscoverContent()`: 发现内容
    - `CustomizeWorkflow()`: 定制工作流
    - `FormatOutput()`: 格式化输出
    - `ValidateOutput()`: 验证输出质量
  - 相关数据结构：
    - `UserInput`: 用户输入
    - `RequestContext`: 请求上下文
    - `LearningObjective`: 学习目标
    - `Content`: 内容
    - `WorkflowConfig`: 工作流配置
    - `QualityReport`: 质量报告

- **`service/adapter_manager.go`**: 适配器管理器
  - `AdapterRegistry`: 适配器注册表
    - 支持适配器的注册、注销、查询
    - 支持按类型、能力查找适配器
  - `AdapterManager`: 适配器管理器
    - 提供适配器生命周期管理
    - 封装适配器调用接口
    - 支持自动选择适配器

#### 设计要点

- **插件化架构**: 每个领域适配器都是独立的插件，实现统一接口
- **动态注册**: 支持运行时动态注册和发现适配器
- **领域隔离**: 不同领域的规则和逻辑通过适配器封装
- **可扩展性**: 新增领域只需实现 `DomainAdapter` 接口

---

### 2. Capability 模块（AI能力）

**位置**: `backend/domain/capability/`

#### 核心文件

- **`entity/capability.go`**: AI能力实体定义
  - `Capability`: AI能力数据模型
  - `CapabilityType`: 能力类型枚举
    - 意图识别 (intent_recognition)
    - 学习目标生成 (objective_generation)
    - 内容发现 (content_discovery)
    - 故事化叙述 (narrative_generation)
    - 质量评估 (quality_assessment)
    - 多模态生成 (multimodal_generation)
    - 画像建模 (persona_modeling)
    - 知识映射 (knowledge_mapping)
  - `CapabilityStatus`: 能力状态

- **`interface.go`**: AI能力接口定义
  - `AICapability`: 所有AI能力的基础接口
    - `Execute()`: 执行能力
    - `Validate()`: 验证输入
    - `GetConfig()`: 获取配置
  - 专用能力接口：
    - `IntentRecognizer`: 意图识别
    - `ObjectiveGenerator`: 学习目标生成
    - `ContentDiscovery`: 内容发现
    - `NarrativeGenerator`: 故事化叙述
    - `QualityAssessor`: 质量评估
    - `MultimodalGenerator`: 多模态生成

- **`service/capability_manager.go`**: AI能力管理器
  - `CapabilityRegistry`: 能力注册表
  - `CapabilityManager`: 能力管理器
    - 支持能力的注册、查询、调用
    - 提供类型化的能力调用方法

#### 设计要点

- **能力抽象**: AI能力作为独立的、可复用的模块
- **领域无关**: AI能力不依赖特定领域，可被多个适配器共享
- **统一接口**: 所有能力实现统一的执行接口
- **类型安全**: 提供类型化的专用接口

---

### 3. Component 模块（组件）

**位置**: `backend/domain/component/`

#### 核心文件

- **`entity/component.go`**: 组件实体定义
  - `Component`: 组件数据模型
  - `ComponentType`: 组件类型（工具、插件、小部件、集成）
  - `ComponentCategory`: 组件分类（可视化、交互、分析、导出、MCP）
  - `ComponentStatus`: 组件状态
  - `ComponentInterface`: 组件接口定义

- **`interface.go`**: 组件执行器接口
  - `ComponentExecutor`: 组件执行器基础接口
  - 专用组件接口：
    - `MCPTool`: MCP工具（Model Context Protocol）
    - `VisualizationComponent`: 可视化组件
    - `ExportComponent`: 导出组件
    - `InteractionComponent`: 交互组件

- **`service/component_manager.go`**: 组件管理器
  - `ComponentRegistry`: 组件注册表
  - `ComponentManager`: 组件管理器
    - 支持组件的注册、查询、执行
    - 提供类型化的组件调用方法

#### 设计要点

- **组件市场**: 支持第三方组件的发布和安装
- **MCP集成**: 原生支持 Model Context Protocol 工具
- **类型多样**: 支持多种类型的组件（可视化、导出等）
- **独立运行**: 组件可独立部署和运行

---

## 🔄 架构关系

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (API/UI)                           │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                 领域适配器层 (Adapter)                       │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐               │
│  │  K12      │  │  美术史    │  │  职业培训  │               │
│  │  Adapter  │  │  Adapter  │  │  Adapter  │               │
│  └───────────┘  └───────────┘  └───────────┘               │
│                                                              │
│  职责：领域规则、请求解析、工作流定制                         │
└─────────────────────────────────────────────────────────────┘
                          ↓ 调用
┌─────────────────────────────────────────────────────────────┐
│                  AI能力层 (Capability)                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │ 意图识别 │ │ 目标生成 │ │ 内容发现 │ │ 质量评估 │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
│                                                              │
│  职责：提供通用AI算法能力                                     │
└─────────────────────────────────────────────────────────────┘
                          ↓ 使用
┌─────────────────────────────────────────────────────────────┐
│                  组件层 (Component)                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │ MCP工具  │ │ 可视化   │ │ 导出     │ │ 交互     │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
│                                                              │
│  职责：提供可复用的功能组件                                   │
└─────────────────────────────────────────────────────────────┘
```

## 📊 数据模型关系

- **Adapter** ← 依赖 → **Capability**: 适配器调用AI能力
- **Adapter** ← 依赖 → **Component**: 适配器使用组件
- **Capability** ← 使用 → **Component**: AI能力可能使用组件
- **Workflow** ← 编排 → **Capability**: 工作流编排AI能力

## 🔧 使用示例

### 注册适配器

```go
// 创建适配器管理器
adapterMgr := service.NewAdapterManager()

// 注册K12适配器（需要实现DomainAdapter接口）
k12Adapter := NewK12Adapter()
err := adapterMgr.RegisterAdapter(k12Adapter)

// 使用适配器
ctx := context.Background()
reqCtx, err := adapterMgr.ParseRequest(ctx, "k12_education", userInput)
objectives, err := adapterMgr.GenerateLearningObjectives(ctx, "k12_education", reqCtx)
```

### 注册AI能力

```go
// 创建能力管理器
capMgr := service.NewCapabilityManager()

// 注册目标生成能力（需要实现ObjectiveGenerator接口）
objGen := NewObjectiveGeneratorCapability()
err := capMgr.RegisterCapability(objGen)

// 使用能力
objectives, err := capMgr.ExecuteObjectiveGeneration(ctx, params)
```

### 注册组件

```go
// 创建组件管理器
compMgr := service.NewComponentManager()

// 注册MCP工具（需要实现MCPTool接口）
mcpTool := NewVirtualLabTool()
err := compMgr.RegisterComponent(mcpTool)

// 使用组件
result, err := compMgr.InvokeMCPTool(ctx, "virtual_lab", params)
```

## 🚀 下一步

1. **实现具体的适配器**: 创建 K12Adapter、ArtHistoryAdapter 等具体实现
2. **实现AI能力**: 基于LLM实现各个AI能力
3. **创建数据库迁移**: 根据实体定义创建数据库表
4. **实现Repository层**: 连接数据库持久化
5. **创建API接口**: 提供REST API供前端调用
6. **编写单元测试**: 确保模块正确性

## 📚 参考文档

- `/docs_new/ARCHITECTURE_CORE_CONCEPTS.md`: 核心概念
- `/docs_new/COZE_INTEGRATION_PLAN.md`: Coze集成方案
- `/docs_new/ADAPTER_MARKETPLACE_DESIGN.md`: 适配器市场设计
- `/docs_new/AI_CAPABILITIES_ANALYSIS.md`: AI能力分析
- `/docs_new/DATABASE_DESIGN.md`: 数据库设计
