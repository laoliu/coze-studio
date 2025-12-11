# MCP、Tools、组件扩展层关系图谱

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 架构设计说明

---

## 📋 目录

1. [概念定义](#概念定义)
2. [架构层级关系](#架构层级关系)
3. [MCP详解](#mcp详解)
4. [Tools详解](#tools详解)
5. [组件扩展层详解](#组件扩展层详解)
6. [组件市场详解](#组件市场详解)
7. [协作关系](#协作关系)
8. [实际应用场景](#实际应用场景)

---

## 概念定义

### 四个核心概念速览

| 概念 | 全称 | 定位 | 主要用途 | 使用者 |
|------|------|------|---------|--------|
| **MCP** | Model Context Protocol | 标准协议 | AI与外部工具通信的标准接口 | AI模型、开发者 |
| **Tools** | Function Calling Tools | 函数工具 | AI可调用的具体函数 | AI模型（通过MCP） |
| **组件扩展层** | Component Extension Layer | 系统层级 | MetaWorkflow的可插拔扩展机制 | MetaWorkflow开发者 |
| **组件市场** | Component Marketplace | 应用平台 | 组件的发布、分享、交易平台 | 所有用户、开发者 |

### 一句话总结

```
MCP（协议标准）定义了 AI 如何调用 Tools（函数）
组件扩展层（系统机制）支持开发者创建 组件（功能模块）
组件市场（平台）让组件可以分享和复用
Tools 可以包装成 组件，组件也可以暴露为 Tools
```

---

## 架构层级关系

### 完整架构图

```
┌──────────────────────────────────────────────────────────────┐
│                    外部标准与生态                             │
│  ┌────────────────────────────────────────────────────┐      │
│  │          MCP (Model Context Protocol) ⭐           │      │
│  │          - Anthropic主导的开放标准                  │      │
│  │          - 定义AI与工具的通信规范                   │      │
│  │          - JSON-RPC 2.0协议                        │      │
│  ├────────────────────────────────────────────────────┤      │
│  │         OpenAI Function Calling                    │      │
│  │         - OpenAI的函数调用规范                      │      │
│  │         - 类似目的，不同实现                        │      │
│  └────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
                            ↓ (实现标准)
┌──────────────────────────────────────────────────────────────┐
│                   MetaWorkflow系统                            │
│  ┌────────────────────────────────────────────────────┐      │
│  │              应用层 (API)                          │      │
│  └────────────────────────────────────────────────────┘      │
│                            ↓                                  │
│  ┌────────────────────────────────────────────────────┐      │
│  │           领域适配器层                              │      │
│  └────────────────────────────────────────────────────┘      │
│                            ↓                                  │
│  ┌────────────────────────────────────────────────────┐      │
│  │            AI能力层                                 │      │
│  │  ┌──────────────┐  ┌──────────────┐               │      │
│  │  │ LLMClient    │  │ ToolManager  │←── MCP实现    │      │
│  │  │ (模型调用)   │  │ (工具管理)   │               │      │
│  │  └──────────────┘  └──────────────┘               │      │
│  └────────────────────────────────────────────────────┘      │
│                            ↓                                  │
│  ┌────────────────────────────────────────────────────┐      │
│  │          工作流引擎层                               │      │
│  └────────────────────────────────────────────────────┘      │
│                            ↓                                  │
│  ┌────────────────────────────────────────────────────┐      │
│  │          工作流节点层                               │      │
│  │  - 节点可以调用Tools                               │      │
│  │  - 节点可以使用组件                                │      │
│  └────────────────────────────────────────────────────┘      │
│                            ↓                                  │
│  ┌────────────────────────────────────────────────────┐      │
│  │         组件扩展层 ⭐                               │      │
│  │  ┌──────────────────────────────────────────┐     │      │
│  │  │   组件注册中心 (Component Registry)       │     │      │
│  │  ├──────────────────────────────────────────┤     │      │
│  │  │   内置组件库                              │     │      │
│  │  │   - ImageGeneratorComponent              │     │      │
│  │  │   - VirtualLabComponent                  │     │      │
│  │  │   - GamificationComponent                │     │      │
│  │  ├──────────────────────────────────────────┤     │      │
│  │  │   第三方组件                              │     │      │
│  │  │   - WikipediaToolComponent (MCP封装)    │     │      │
│  │  │   - WeatherToolComponent (MCP封装)      │     │      │
│  │  │   - CalculatorToolComponent             │     │      │
│  │  └──────────────────────────────────────────┘     │      │
│  └────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
                            ↓ (发布/安装)
┌──────────────────────────────────────────────────────────────┐
│              组件市场 (Component Marketplace) ⭐              │
│  ┌────────────────────────────────────────────────────┐      │
│  │  组件分类：                                         │      │
│  │  - 内容生成类 (基于MCP的图片/视频/音频生成器)       │      │
│  │  - 质量检查类 (语法检查、事实核查)                  │      │
│  │  - 交互设计类 (游戏化、虚拟实验室)                  │      │
│  │  - 工具集成类 (MCP Tools封装) ⭐                   │      │
│  │  - 数据源类 (Wikipedia、YouTube、Weather API)      │      │
│  └────────────────────────────────────────────────────┘      │
│                                                               │
│  组件格式：                                                   │
│  - Python Package (pip安装)                                  │
│  - MCP Server (独立进程)                                     │
│  - Docker Container (容器化)                                 │
└──────────────────────────────────────────────────────────────┘
```

### 关键关系总结

```
┌─────────────┐
│     MCP     │ ← 协议标准（如何通信）
└─────────────┘
       ↓ (定义接口)
┌─────────────┐
│    Tools    │ ← 具体函数（做什么事）
└─────────────┘
       ↓ (可以封装为)
┌─────────────┐
│    组件     │ ← MetaWorkflow的扩展单元
└─────────────┘
       ↓ (发布到)
┌─────────────┐
│  组件市场   │ ← 分享和复用平台
└─────────────┘
```

---

## MCP详解

### 什么是MCP？

**MCP (Model Context Protocol)** 是由 Anthropic 提出的**开放标准**，用于定义 AI 模型如何与外部工具和数据源进行通信。

```
类比：MCP 就像 USB 标准
- USB 定义了设备如何连接电脑
- MCP 定义了工具如何连接 AI 模型
```

### MCP架构

```
┌──────────────────────────────────────────────────┐
│              MCP 客户端 (Client)                  │
│           (AI应用，如Claude、ChatGPT)             │
└──────────────────────────────────────────────────┘
                    ↓ ↑ (JSON-RPC 2.0)
┌──────────────────────────────────────────────────┐
│              MCP 服务器 (Server)                  │
│           (提供工具和资源的服务)                  │
│  ┌────────────────────────────────────────┐     │
│  │  Tools (工具)                          │     │
│  │  - get_weather(location)              │     │
│  │  - search_wikipedia(query)            │     │
│  │  - run_python_code(code)              │     │
│  ├────────────────────────────────────────┤     │
│  │  Resources (资源)                      │     │
│  │  - 文件系统访问                        │     │
│  │  - 数据库查询                          │     │
│  │  - API调用                            │     │
│  ├────────────────────────────────────────┤     │
│  │  Prompts (提示词模板)                  │     │
│  │  - 预定义的提示词                      │     │
│  └────────────────────────────────────────┘     │
└──────────────────────────────────────────────────┘
```

### MCP规范示例

```json
// MCP服务器暴露工具的格式
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}

// 响应：可用工具列表
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "get_weather",
        "description": "获取指定城市的天气信息",
        "inputSchema": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "城市名称"
            }
          },
          "required": ["location"]
        }
      },
      {
        "name": "search_wikipedia",
        "description": "搜索维基百科",
        "inputSchema": {
          "type": "object",
          "properties": {
            "query": {
              "type": "string",
              "description": "搜索关键词"
            }
          },
          "required": ["query"]
        }
      }
    ]
  }
}
```

### 在MetaWorkflow中的实现

```python
# src/ai/mcp_client.py

from typing import List, Dict, Any
import httpx

class MCPClient:
    """MCP客户端实现"""
    
    def __init__(self, server_url: str):
        self.server_url = server_url
        self.client = httpx.AsyncClient()
    
    async def list_tools(self) -> List[Dict[str, Any]]:
        """
        列出MCP服务器提供的所有工具
        
        Returns:
            工具列表，每个工具包含name、description、inputSchema
        """
        response = await self.client.post(
            f"{self.server_url}/rpc",
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/list",
                "params": {}
            }
        )
        result = response.json()
        return result["result"]["tools"]
    
    async def call_tool(
        self,
        tool_name: str,
        arguments: Dict[str, Any]
    ) -> Any:
        """
        调用MCP工具
        
        Args:
            tool_name: 工具名称
            arguments: 工具参数
        
        Returns:
            工具执行结果
        """
        response = await self.client.post(
            f"{self.server_url}/rpc",
            json={
                "jsonrpc": "2.0",
                "id": 2,
                "method": "tools/call",
                "params": {
                    "name": tool_name,
                    "arguments": arguments
                }
            }
        )
        result = response.json()
        return result["result"]


# 使用示例
mcp_client = MCPClient("http://localhost:8888")

# 列出所有可用工具
tools = await mcp_client.list_tools()
# 输出: [
#   {"name": "get_weather", "description": "...", "inputSchema": {...}},
#   {"name": "search_wikipedia", ...}
# ]

# 调用工具
weather = await mcp_client.call_tool(
    tool_name="get_weather",
    arguments={"location": "北京"}
)
# 输出: {"temperature": 5, "condition": "晴朗"}
```

---

## Tools详解

### 什么是Tools？

**Tools（工具）** 是AI模型可以调用的**具体函数**，用于执行特定任务。

```
分类：
1. 数据获取类：获取天气、搜索Wikipedia、查询数据库
2. 计算执行类：运行代码、数学计算、数据分析
3. 内容生成类：生成图片、生成音频、翻译文本
4. 系统操作类：文件操作、发送邮件、创建日历事件
```

### Tools的两种形式

#### 形式1: OpenAI Function Calling

```python
# OpenAI的函数调用格式

tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "location": {
                        "type": "string",
                        "description": "城市名称，如'北京'、'上海'"
                    },
                    "unit": {
                        "type": "string",
                        "enum": ["celsius", "fahrenheit"],
                        "description": "温度单位"
                    }
                },
                "required": ["location"]
            }
        }
    }
]

# 在LLM调用中使用
response = await openai.chat.completions.create(
    model="gpt-4",
    messages=[{"role": "user", "content": "北京今天天气如何？"}],
    tools=tools,
    tool_choice="auto"
)

# AI决定调用工具
if response.choices[0].message.tool_calls:
    tool_call = response.choices[0].message.tool_calls[0]
    # tool_call.function.name = "get_weather"
    # tool_call.function.arguments = '{"location": "北京"}'
    
    # 执行工具
    result = get_weather(location="北京")
    
    # 将结果返回给AI
    second_response = await openai.chat.completions.create(
        model="gpt-4",
        messages=[
            {"role": "user", "content": "北京今天天气如何？"},
            response.choices[0].message,
            {
                "role": "tool",
                "tool_call_id": tool_call.id,
                "content": str(result)
            }
        ]
    )
```

#### 形式2: MCP Tools

```python
# MCP格式的工具定义（更标准化）

class WeatherMCPServer:
    """天气MCP服务器"""
    
    async def handle_list_tools(self):
        """列出可用工具"""
        return {
            "tools": [
                {
                    "name": "get_weather",
                    "description": "获取指定城市的天气信息",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "location": {
                                "type": "string",
                                "description": "城市名称"
                            }
                        },
                        "required": ["location"]
                    }
                }
            ]
        }
    
    async def handle_call_tool(self, name: str, arguments: dict):
        """执行工具调用"""
        if name == "get_weather":
            return await self._get_weather(arguments["location"])
    
    async def _get_weather(self, location: str):
        """实际的天气查询实现"""
        # 调用天气API
        api_result = await weather_api.get(location)
        return {
            "temperature": api_result.temp,
            "condition": api_result.condition,
            "humidity": api_result.humidity
        }
```

### 在MetaWorkflow中集成Tools

```python
# src/ai/tool_manager.py

class ToolManager:
    """工具管理器（支持MCP和OpenAI Function Calling）"""
    
    def __init__(self):
        self.tools: Dict[str, Tool] = {}
        self.mcp_clients: List[MCPClient] = []
    
    async def register_mcp_server(self, server_url: str):
        """注册MCP服务器"""
        client = MCPClient(server_url)
        
        # 获取服务器提供的工具
        tools = await client.list_tools()
        
        # 注册到本地工具库
        for tool_def in tools:
            self.tools[tool_def["name"]] = MCPTool(
                client=client,
                definition=tool_def
            )
        
        self.mcp_clients.append(client)
    
    def register_local_tool(self, name: str, func: Callable, schema: Dict):
        """注册本地Python函数为工具"""
        self.tools[name] = LocalTool(
            name=name,
            func=func,
            schema=schema
        )
    
    async def call_tool(self, name: str, arguments: Dict) -> Any:
        """调用工具（自动路由到MCP或本地）"""
        if name not in self.tools:
            raise ValueError(f"Tool {name} not found")
        
        tool = self.tools[name]
        return await tool.execute(arguments)
    
    def get_tool_definitions_for_llm(self) -> List[Dict]:
        """
        获取工具定义（用于传给LLM）
        转换为OpenAI Function Calling格式
        """
        definitions = []
        for tool in self.tools.values():
            definitions.append({
                "type": "function",
                "function": {
                    "name": tool.name,
                    "description": tool.description,
                    "parameters": tool.input_schema
                }
            })
        return definitions


# 使用示例
tool_manager = ToolManager()

# 1. 注册MCP服务器的工具
await tool_manager.register_mcp_server("http://localhost:8888")

# 2. 注册本地Python函数
def calculate(expression: str) -> float:
    """计算数学表达式"""
    return eval(expression)

tool_manager.register_local_tool(
    name="calculator",
    func=calculate,
    schema={
        "type": "object",
        "properties": {
            "expression": {
                "type": "string",
                "description": "数学表达式，如'2+3*4'"
            }
        },
        "required": ["expression"]
    }
)

# 3. 在LLM调用中使用
llm_response = await llm_client.chat(
    messages=[{"role": "user", "content": "北京天气如何？2+3等于多少？"}],
    tools=tool_manager.get_tool_definitions_for_llm()
)

# 4. 执行工具调用
if llm_response.tool_calls:
    for tool_call in llm_response.tool_calls:
        result = await tool_manager.call_tool(
            name=tool_call.function.name,
            arguments=json.loads(tool_call.function.arguments)
        )
```

---

## 组件扩展层详解

### 什么是组件扩展层？

**组件扩展层**是MetaWorkflow的**可插拔架构**，允许开发者在不修改核心代码的情况下添加新功能。

### 组件 vs Tools 的区别

| 维度 | Tools | 组件 (Component) |
|------|-------|-----------------|
| **范围** | 单一函数 | 完整功能模块 |
| **复杂度** | 简单（输入→输出） | 复杂（状态管理、配置、UI） |
| **生命周期** | 无状态 | 有完整生命周期（init、execute、cleanup） |
| **调用者** | AI模型 | 工作流节点、用户 |
| **示例** | get_weather() | VirtualLabComponent（包含3D渲染、交互逻辑、状态保存） |

### 组件的三种类型

#### 类型1: 节点型组件（最常见）

```python
# components/virtual_lab_component.py

class VirtualLabComponent(WorkflowNode):
    """虚拟实验室组件（节点型）"""
    
    # 组件元数据
    component_metadata = {
        "name": "虚拟化学实验室",
        "version": "1.2.0",
        "author": "张三",
        "category": "interaction_design",
        "description": "提供3D化学实验模拟",
        "dependencies": ["three.js", "chemistry-engine"],
        "pricing": "免费",
        "license": "MIT"
    }
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__()
        self.config = config
        self.renderer = ThreeJSRenderer()
        self.chemistry_engine = ChemistryEngine()
    
    async def execute(
        self,
        context: WorkflowContext,
        input_data: Dict
    ) -> Dict:
        """
        执行节点
        
        输入: {
            "experiment_type": "燃烧实验",
            "chemicals": ["蜡烛", "氧气"],
            "duration": 300
        }
        
        输出: {
            "lab_url": "https://lab.example.com/session123",
            "interaction_script": {...},
            "expected_results": [...]
        }
        """
        # 1. 生成3D场景
        scene = await self.renderer.create_scene(
            experiment_type=input_data["experiment_type"]
        )
        
        # 2. 添加化学物质
        for chemical in input_data["chemicals"]:
            await scene.add_object(chemical)
        
        # 3. 配置交互逻辑
        interactions = await self._configure_interactions(input_data)
        
        # 4. 保存会话
        session = await self._save_session(scene, interactions)
        
        return {
            "lab_url": session.url,
            "interaction_script": interactions,
            "expected_results": self._predict_results(input_data)
        }
    
    async def _configure_interactions(self, input_data):
        """配置交互逻辑"""
        # 点击、拖拽、化学反应模拟等
        pass
```

#### 类型2: 能力增强型组件

```python
# components/dalle3_image_generator.py

class DALLE3ImageGenerator:
    """DALL-E 3图片生成器（能力增强型组件）"""
    
    component_metadata = {
        "name": "DALL-E 3 图片生成器",
        "version": "2.0.0",
        "category": "capability_extension",
        "extends": "MultimodalGenerator",  # 扩展现有AI能力
        "provider": "OpenAI",
        "pricing": "按次计费"
    }
    
    def __init__(self, api_key: str):
        self.client = openai.AsyncOpenAI(api_key=api_key)
    
    async def generate_image(
        self,
        prompt: str,
        size: str = "1024x1024",
        quality: str = "standard"
    ) -> Image:
        """生成图片"""
        response = await self.client.images.generate(
            model="dall-e-3",
            prompt=prompt,
            size=size,
            quality=quality,
            n=1
        )
        
        return Image(
            url=response.data[0].url,
            revised_prompt=response.data[0].revised_prompt
        )


# 注册到MultimodalGenerator
multimodal_generator.register_provider(
    name="dalle3",
    component=DALLE3ImageGenerator(api_key=os.getenv("OPENAI_API_KEY"))
)
```

#### 类型3: Tool封装型组件

```python
# components/wikipedia_tool_component.py

class WikipediaToolComponent:
    """Wikipedia工具组件（封装MCP Tool）"""
    
    component_metadata = {
        "name": "Wikipedia搜索工具",
        "version": "1.0.0",
        "category": "tool_integration",
        "mcp_compatible": True,
        "data_source": "Wikipedia API"
    }
    
    def __init__(self):
        # 可以内部实现，也可以连接MCP服务器
        self.use_mcp = os.getenv("USE_MCP_WIKIPEDIA", "false") == "true"
        
        if self.use_mcp:
            self.mcp_client = MCPClient("http://localhost:8888")
        else:
            self.wikipedia = wikipediaapi.Wikipedia('en')
    
    async def search(self, query: str, limit: int = 5) -> List[Dict]:
        """搜索Wikipedia"""
        if self.use_mcp:
            # 使用MCP调用
            return await self.mcp_client.call_tool(
                tool_name="search_wikipedia",
                arguments={"query": query, "limit": limit}
            )
        else:
            # 直接调用Wikipedia API
            results = self.wikipedia.search(query, results=limit)
            return [
                {
                    "title": r.title,
                    "summary": r.summary,
                    "url": r.fullurl
                }
                for r in results
            ]
    
    def as_mcp_tool(self) -> Dict:
        """
        将组件暴露为MCP Tool定义
        （可以被其他MCP客户端调用）
        """
        return {
            "name": "search_wikipedia",
            "description": "搜索Wikipedia获取知识",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "query": {
                        "type": "string",
                        "description": "搜索关键词"
                    },
                    "limit": {
                        "type": "integer",
                        "description": "返回结果数量",
                        "default": 5
                    }
                },
                "required": ["query"]
            }
        }
```

---

## 组件市场详解

### 什么是组件市场？

**组件市场 (Component Marketplace)** 是一个**平台**，允许开发者发布、分享、交易组件。

```
类比：
- Chrome Web Store：浏览器扩展市场
- VS Code Marketplace：编辑器插件市场
- Component Marketplace：MetaWorkflow组件市场
```

### 组件市场架构

```
┌──────────────────────────────────────────────────┐
│         组件市场 (Web平台)                        │
├──────────────────────────────────────────────────┤
│                                                   │
│  【浏览和搜索】                                   │
│  - 分类浏览：内容生成/质量检查/交互设计/...       │
│  - 关键词搜索："虚拟实验室"、"图片生成"           │
│  - 标签过滤：#化学 #3D #免费 #MCP兼容             │
│  - 排序：热度/评分/下载量/最新                    │
│                                                   │
│  【组件详情页】                                   │
│  - 基本信息：名称、版本、作者、许可证             │
│  - 功能描述：详细说明、截图、演示视频             │
│  - 技术信息：依赖项、API文档、兼容性              │
│  - 用户评价：评分、评论、使用案例                 │
│                                                   │
│  【安装和管理】                                   │
│  - 一键安装：pip install / npm install            │
│  - 版本管理：升级、回退、卸载                     │
│  - 依赖处理：自动安装依赖                         │
│  - 配置向导：首次使用配置                         │
│                                                   │
│  【开发者中心】                                   │
│  - 发布组件：上传代码、填写元数据                 │
│  - 版本更新：发布新版本                           │
│  - 数据统计：下载量、用户反馈                     │
│  - 收益管理：付费组件的收入                       │
│                                                   │
└──────────────────────────────────────────────────┘
```

### 组件市场分类

```
组件市场分类体系
├── 内容生成类 (Content Generation)
│   ├── 图片生成
│   │   ├── DALL-E 3插件 ($) - OpenAI官方
│   │   ├── Midjourney插件 ($) - 第三方
│   │   └── Stable Diffusion插件 (免费) - 开源
│   ├── 视频生成
│   │   ├── D-ID虚拟人 ($)
│   │   └── FFmpeg处理器 (免费)
│   └── 音频生成
│       ├── Azure TTS (免费额度)
│       └── ElevenLabs ($)
│
├── 质量检查类 (Quality Assurance)
│   ├── 语法检查器 (LanguageTool) - 免费
│   ├── 事实核查器 (FactCheck.org API) - 免费
│   └── 抄袭检测器 (Turnitin API) - $
│
├── 交互设计类 (Interaction Design)
│   ├── 虚拟实验室 (VirtualLab) - $
│   ├── 游戏化组件 (Gamification) - 免费
│   └── AR/VR组件 (Unity WebGL) - $
│
├── 工具集成类 (Tool Integration) ⭐ MCP封装
│   ├── Wikipedia搜索 (MCP) - 免费
│   ├── 天气查询 (MCP) - 免费
│   ├── Google搜索 (MCP) - 需API Key
│   ├── YouTube数据 (MCP) - 需API Key
│   └── 代码执行器 (MCP) - 免费
│
└── 数据源类 (Data Sources)
    ├── WikiArt美术数据库 - 免费
    ├── ChemSpider化学数据库 - 免费
    └── Coursera课程数据 - $
```

### 组件包格式

```yaml
# component.yaml - 组件元数据文件

metadata:
  name: "virtual-chemistry-lab"
  display_name: "虚拟化学实验室"
  version: "1.2.0"
  author: "张三"
  author_email: "zhangsan@example.com"
  license: "MIT"
  homepage: "https://github.com/zhangsan/virtual-lab"
  
  description: |
    提供3D化学实验模拟，支持常见中学化学实验，
    包括燃烧、氧化还原、酸碱中和等。
  
  category: "interaction_design"
  tags:
    - 化学
    - 实验
    - 3D
    - 教育
  
  pricing:
    type: "freemium"  # free / freemium / paid
    free_tier:
      max_experiments_per_month: 100
    paid_tier:
      price: 29.99
      currency: "USD"
      billing: "monthly"

dependencies:
  python: ">=3.9"
  packages:
    - "three.js>=0.150.0"
    - "chemistry-engine>=2.0.0"
    - "pydantic>=2.0.0"
  
  mcp_servers:  # 依赖的MCP服务器
    - name: "chemistry-data"
      url: "http://chemdata.example.com"
      optional: true

compatibility:
  metaworkflow: ">=1.0.0"
  node_types:
    - "experiment_node"
    - "simulation_node"

installation:
  type: "pip"  # pip / npm / docker
  command: "pip install virtual-chemistry-lab"
  
  post_install:
    - "Download 3D models from CDN"
    - "Initialize chemistry database"

configuration:
  required:
    - name: "api_key"
      type: "string"
      description: "ChemData API密钥"
      env_var: "CHEMDATA_API_KEY"
  
  optional:
    - name: "render_quality"
      type: "enum"
      values: ["low", "medium", "high"]
      default: "medium"

screenshots:
  - "https://example.com/screenshots/1.png"
  - "https://example.com/screenshots/2.png"

demo_video: "https://youtube.com/watch?v=xxx"

documentation: "https://docs.example.com/virtual-lab"

support:
  email: "support@example.com"
  issues: "https://github.com/zhangsan/virtual-lab/issues"
  discord: "https://discord.gg/xxx"
```

---

## 协作关系

### 关系图谱

```
┌────────────────────────────────────────────────────────┐
│                     使用场景                           │
├────────────────────────────────────────────────────────┤
│                                                        │
│  【场景1】AI需要调用外部工具                           │
│                                                        │
│  用户："北京今天天气如何？"                             │
│     ↓                                                  │
│  LLM分析：需要获取实时天气数据                          │
│     ↓                                                  │
│  LLM决定：调用 get_weather tool                        │
│     ↓                                                  │
│  ToolManager 执行：                                    │
│     - 如果是MCP Tool → 通过MCP Client调用服务器        │
│     - 如果是本地Tool → 直接执行Python函数              │
│     ↓                                                  │
│  返回结果给LLM → LLM生成回答                           │
│                                                        │
│  ✅ 这里使用：MCP + Tools                             │
│  ❌ 不涉及：组件、组件市场                             │
│                                                        │
├────────────────────────────────────────────────────────┤
│                                                        │
│  【场景2】工作流需要虚拟实验室                         │
│                                                        │
│  K12化学课程工作流执行到 ExperimentNode                │
│     ↓                                                  │
│  检查：是否安装了VirtualLabComponent？                │
│     ↓                                                  │
│  如果已安装：                                          │
│     - 调用组件的execute()方法                         │
│     - 组件内部可能调用MCP Tools获取化学数据            │
│     - 组件生成3D实验室                                │
│     ↓                                                  │
│  返回虚拟实验室URL                                    │
│                                                        │
│  ✅ 这里使用：组件（从组件市场安装）                   │
│  ✅ 组件内部可能使用：MCP Tools                       │
│                                                        │
├────────────────────────────────────────────────────────┤
│                                                        │
│  【场景3】开发者想分享自己的工具                       │
│                                                        │
│  开发者创建：WeatherToolComponent                     │
│     ↓                                                  │
│  选择实现方式：                                        │
│     - 方式A：实现为MCP Server                         │
│     - 方式B：实现为Python包                           │
│     ↓                                                  │
│  打包：创建component.yaml + 代码                      │
│     ↓                                                  │
│  发布到组件市场                                        │
│     ↓                                                  │
│  其他用户：                                            │
│     - 在市场浏览                                       │
│     - 一键安装                                         │
│     - 在工作流中使用                                   │
│                                                        │
│  ✅ 这里使用：组件市场                                │
│  ✅ 组件可以选择：基于MCP或不基于MCP                  │
│                                                        │
└────────────────────────────────────────────────────────┘
```

### 技术栈对应关系

| 层级 | 技术选型 | 示例 |
|------|---------|------|
| **MCP协议层** | JSON-RPC 2.0 | MCP Specification v1.0 |
| **MCP Server实现** | Python/Node.js | FastAPI/Express |
| **MCP Client实现** | httpx/axios | MetaWorkflow的MCPClient |
| **Tools定义** | JSON Schema | OpenAI Function Calling格式 |
| **Tools执行** | Python函数/MCP调用 | async def get_weather() |
| **组件包格式** | Python Package | setup.py + component.yaml |
| **组件注册** | 插件系统 | ComponentRegistry |
| **组件市场后端** | FastAPI + PostgreSQL | REST API |
| **组件市场前端** | React/Vue | Web UI |

### 数据流向图

```
┌──────────────────────────────────────────────────────────┐
│ 1. 用户请求：生成化学课程                                 │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 2. 领域适配器选择工作流模板                               │
│    template: k12_science_workflow                        │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 3. 工作流引擎执行节点                                     │
│    ExperimentNode → 需要虚拟实验室                       │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 4. 检查组件                                               │
│    ComponentRegistry.get("VirtualLabComponent")          │
│    - 已安装 → 使用组件                                   │
│    - 未安装 → 提示用户去组件市场安装                     │
└──────────────────────────────────────────────────────────┘
                    ↓ (假设已安装)
┌──────────────────────────────────────────────────────────┐
│ 5. VirtualLabComponent执行                               │
│    component.execute(context, input_data)                │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 6. 组件内部可能调用MCP Tools                             │
│    - 调用 chemistry-data MCP Server                      │
│    - 获取化学物质的3D模型数据                            │
│    - 获取反应方程式                                      │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 7. 组件生成虚拟实验室                                     │
│    返回: {                                                │
│      "lab_url": "https://lab.example.com/session123",   │
│      "interactions": {...}                               │
│    }                                                     │
└──────────────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────────────┐
│ 8. 工作流继续执行其他节点                                 │
│    最终返回完整课程（包含虚拟实验室链接）                 │
└──────────────────────────────────────────────────────────┘
```

---

## 实际应用场景

### 场景1: 使用MCP Tools获取实时数据

```python
# ========== 场景：AI生成课程时需要最新数据 ==========

# 用户请求
user_input = "讲解2024年诺贝尔化学奖的研究内容"

# 工作流执行
async def generate_nobel_lecture():
    # 1. AI识别需要实时数据
    intent = await intent_recognizer.analyze(user_input)
    # intent.needs_realtime_data = True
    
    # 2. 调用MCP Tool获取诺贝尔奖信息
    nobel_data = await tool_manager.call_tool(
        tool_name="search_wikipedia",  # MCP Tool
        arguments={
            "query": "2024 Nobel Prize in Chemistry",
            "limit": 3
        }
    )
    
    # 3. 将数据传递给内容生成
    content = await narrative_generator.generate(
        topic="诺贝尔化学奖研究",
        reference_data=nobel_data
    )
    
    return content

# 这里的关键：
# ✅ MCP Tool提供最新数据
# ✅ 不需要组件（简单的工具调用）
# ✅ 不涉及组件市场
```

### 场景2: 使用组件增强课程体验

```python
# ========== 场景：使用虚拟实验室组件 ==========

# 用户请求
user_input = "讲解化学燃烧条件，需要虚拟实验"

# 工作流执行
async def generate_chemistry_lesson_with_lab():
    # 1. 领域适配器选择模板
    adapter = K12Adapter()
    request = await adapter.parse_request(user_input)
    template = await adapter.select_workflow_template(request)
    # template = "k12_science_workflow"（包含ExperimentNode）
    
    # 2. 工作流执行到ExperimentNode
    class ExperimentNode(WorkflowNode):
        async def execute(self, context, input_data):
            # 检查是否安装了虚拟实验室组件
            virtual_lab = component_registry.get("VirtualLabComponent")
            
            if virtual_lab:
                # 使用组件生成虚拟实验室
                lab_result = await virtual_lab.execute(
                    context=context,
                    input_data={
                        "experiment_type": "燃烧条件实验",
                        "chemicals": ["蜡烛", "玻璃罩", "氧气"],
                        "duration": 300
                    }
                )
                return lab_result
            else:
                # 组件未安装，返回提示
                return {
                    "error": "VirtualLabComponent not installed",
                    "install_link": "https://marketplace.metaworkflow.com/virtual-lab"
                }
    
    # 3. 虚拟实验室组件内部可能使用MCP
    class VirtualLabComponent:
        async def execute(self, context, input_data):
            # 调用化学数据MCP Server
            chemical_data = await mcp_client.call_tool(
                tool_name="get_chemical_properties",
                arguments={"name": "蜡烛"}
            )
            
            # 使用数据生成3D模型
            scene = await self._create_3d_scene(chemical_data)
            
            return {"lab_url": scene.url}

# 这里的关键：
# ✅ 组件提供复杂功能（3D渲染）
# ✅ 组件从组件市场安装
# ✅ 组件内部可以使用MCP Tools
```

### 场景3: 开发者发布组件到市场

```python
# ========== 场景：开发者创建并发布天气组件 ==========

# 步骤1: 开发组件
class WeatherComponent:
    """天气查询组件"""
    
    component_metadata = {
        "name": "weather-api-component",
        "display_name": "天气查询组件",
        "version": "1.0.0",
        "author": "李四",
        "category": "tool_integration",
        "mcp_compatible": True
    }
    
    def __init__(self, api_key: str):
        self.api_key = api_key
    
    async def get_weather(self, location: str) -> Dict:
        """获取天气"""
        # 实现1：直接调用API
        response = await httpx.get(
            f"https://api.weather.com/v1/current",
            params={"location": location, "key": self.api_key}
        )
        return response.json()
    
    def as_mcp_server(self):
        """
        将组件暴露为MCP Server
        这样其他AI应用也可以使用
        """
        from fastapi import FastAPI
        app = FastAPI()
        
        @app.post("/rpc")
        async def handle_rpc(request: dict):
            if request["method"] == "tools/list":
                return {
                    "tools": [{
                        "name": "get_weather",
                        "description": "获取天气信息",
                        "inputSchema": {
                            "type": "object",
                            "properties": {
                                "location": {"type": "string"}
                            }
                        }
                    }]
                }
            elif request["method"] == "tools/call":
                if request["params"]["name"] == "get_weather":
                    result = await self.get_weather(
                        request["params"]["arguments"]["location"]
                    )
                    return {"result": result}
        
        return app

# 步骤2: 打包组件
# setup.py
from setuptools import setup

setup(
    name="weather-api-component",
    version="1.0.0",
    author="李四",
    packages=["weather_component"],
    install_requires=["httpx>=0.24.0"],
    include_package_data=True,
    package_data={
        "weather_component": ["component.yaml"]
    }
)

# 步骤3: 发布到组件市场
# 通过市场Web界面或CLI工具
$ metaworkflow publish weather-api-component
✅ Component published successfully!
📦 Package: weather-api-component v1.0.0
🔗 Market URL: https://marketplace.metaworkflow.com/weather-api-component

# 步骤4: 其他用户安装
$ metaworkflow install weather-api-component
✅ Component installed successfully!
📝 Please configure API key:
   export WEATHER_API_KEY=your_key_here

# 步骤5: 在工作流中使用
# 方式A：作为组件使用
weather_component = component_registry.get("WeatherComponent")
weather = await weather_component.get_weather("北京")

# 方式B：作为MCP Tool使用
tool_manager.register_mcp_server("http://localhost:8000")
weather = await tool_manager.call_tool("get_weather", {"location": "北京"})
```

### 场景4: 组件市场的完整生态

```
┌────────────────────────────────────────────────────────┐
│              组件市场生态                               │
├────────────────────────────────────────────────────────┤
│                                                        │
│  【官方组件】（MetaWorkflow团队维护）                   │
│  ├─ 基础工具包                                         │
│  │  ├─ Wikipedia Search (MCP) - 免费                  │
│  │  ├─ Calculator - 免费                              │
│  │  └─ Code Executor (MCP) - 免费                     │
│  ├─ 多模态生成                                         │
│  │  ├─ Image Generator (DALL-E 3) - $                │
│  │  └─ TTS Generator (Azure) - 免费额度              │
│  └─ 教学组件                                           │
│      ├─ Quiz Generator - 免费                         │
│      └─ Progress Tracker - 免费                       │
│                                                        │
│  【认证组件】（经过审核的第三方）                       │
│  ├─ 虚拟实验室 (VirtualLab) - $ - ⭐4.8              │
│  ├─ 游戏化引擎 (Gamification) - $ - ⭐4.6            │
│  ├─ AR/VR组件 (ARComponent) - $ - ⭐4.5             │
│  └─ 事实核查器 (FactChecker MCP) - 免费 - ⭐4.7      │
│                                                        │
│  【社区组件】（用户贡献）                               │
│  ├─ Minecraft教育版集成 - 免费 - ⭐4.2               │
│  ├─ Scratch编程集成 - 免费 - ⭐4.0                   │
│  ├─ 中国古诗词数据库 - 免费 - ⭐4.3                   │
│  └─ 化学方程式配平器 - 免费 - ⭐3.9                   │
│                                                        │
│  【MCP生态集成】（外部MCP Servers）                    │
│  ├─ Claude Desktop MCP Servers                        │
│  │  ├─ Filesystem (文件系统访问)                      │
│  │  ├─ GitHub (代码仓库集成)                          │
│  │  └─ Google Drive (云存储)                         │
│  └─ 第三方MCP Servers                                 │
│      ├─ Wolfram Alpha (数学计算)                      │
│      └─ Notion (笔记管理)                             │
│                                                        │
└────────────────────────────────────────────────────────┘
```

---

## 总结

### 关系总结表

| 对比维度 | MCP | Tools | 组件扩展层 | 组件市场 |
|---------|-----|-------|-----------|---------|
| **是什么** | 通信协议 | 函数工具 | 系统架构层 | 应用平台 |
| **定位** | 标准规范 | 具体实现 | 内部机制 | 外部生态 |
| **作用域** | AI生态通用 | AI可调用 | MetaWorkflow专属 | 用户可见 |
| **粒度** | 协议级 | 函数级 | 模块级 | 产品级 |
| **开发者** | Anthropic | 任何人 | MetaWorkflow + 插件开发者 | 社区 |

### 核心理解

#### 1️⃣ MCP是协议标准

```
就像HTTP是Web的协议标准
MCP是AI工具调用的协议标准

不同AI应用（Claude、ChatGPT、MetaWorkflow）
都可以通过MCP调用同一个工具服务器
```

#### 2️⃣ Tools是具体函数

```
get_weather()、search_wikipedia()、run_code()
这些都是Tools（工具函数）

可以通过两种方式暴露：
- OpenAI Function Calling格式
- MCP Server格式
```

#### 3️⃣ 组件扩展层是MetaWorkflow的插件系统

```
允许开发者创建：
- 新的工作流节点（如VirtualLabNode）
- AI能力增强（如DALLE3ImageGenerator）
- Tool封装（如WikipediaToolComponent）

组件可以内部使用MCP Tools
也可以暴露自己为MCP Server
```

#### 4️⃣ 组件市场是分享平台

```
就像Chrome Web Store
开发者可以：
- 发布组件
- 设置价格（免费/付费）
- 获得反馈

用户可以：
- 浏览组件
- 一键安装
- 评价和评论
```

### 实际应用建议

**何时使用MCP？**
- ✅ 想让AI调用外部API
- ✅ 需要实时数据（天气、新闻、股票）
- ✅ 想与其他AI应用共享工具
- ✅ 需要标准化的工具接口

**何时开发组件？**
- ✅ 功能复杂，不是简单函数
- ✅ 有状态管理、配置、UI
- ✅ 想在MetaWorkflow中复用
- ✅ 想分享给其他用户

**何时发布到组件市场？**
- ✅ 组件开发完成且测试通过
- ✅ 有通用价值（不是个人专用）
- ✅ 想获得收益或社区认可
- ✅ 愿意维护和支持用户

### 记忆口诀

```
MCP是协议，定义如何通信
Tools是函数，AI可以调用
组件是模块，扩展系统功能
市场是平台，分享和复用

MCP可封装成组件
组件可暴露为MCP
市场连接开发者和用户
工具赋能AI智能化
```

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow Team  

**相关文档**:
- [核心概念关系图谱](./ARCHITECTURE_CORE_CONCEPTS.md)
- [组件扩展指南](./COMPONENT_EXTENSION_GUIDE.md)
- [AI能力分析报告](./AI_CAPABILITIES_ANALYSIS.md)

**参考资源**:
- [MCP官方文档](https://modelcontextprotocol.io/)
- [OpenAI Function Calling文档](https://platform.openai.com/docs/guides/function-calling)
- [Anthropic Claude API文档](https://docs.anthropic.com/)
