# 🎯 调试 demo.go 快速参考

## ⚡ 最快方式（1 分钟）

```bash
# 直接运行（无需任何依赖！）
cd /home/liunix/work/coze-studio
./debug_adapter.sh k12
```

## 📝 三种调试方式

### 1️⃣ 命令行运行（最简单）

```bash
# 方式 A: 使用脚本
./debug_adapter.sh k12          # K12 演示
./debug_adapter.sh registry     # 注册表演示
./debug_adapter.sh all          # 所有演示

# 方式 B: 直接运行
cd backend/domain/adapter/examples
go run . -demo k12
```

### 2️⃣ VSCode 调试（推荐）

1. **打开文件**: `backend/domain/adapter/examples/demo.go`
2. **设置断点**: 点击行号左侧（红点）
   - 第 15 行：函数入口
   - 第 44 行：用户输入
   - 第 56 行：请求解析
3. **按 F5**: 选择 "Debug Adapter Demo (K12)"
4. **开始调试**:
   - `F10` = 单步跳过
   - `F11` = 单步进入
   - `F5` = 继续执行

### 3️⃣ Delve 调试（高级）

```bash
# 使用脚本启动 Delve
./debug_adapter.sh debug k12

# 或手动启动
cd backend/domain/adapter/examples
dlv debug . -- -demo k12

# Delve 命令
(dlv) break demo.go:15    # 设置断点
(dlv) continue            # 开始执行
(dlv) print userInput     # 打印变量
(dlv) next                # 下一步
```

## 🔍 推荐断点位置

| 行号 | 代码 | 用途 |
|------|------|------|
| 15 | `func DemoK12Adapter()` | 函数入口 |
| 21 | `RegisterAdapter(k12Adapter)` | 注册适配器 |
| 44 | `userInput := &entity.UserInput{...}` | 用户输入 |
| 56 | `ParseRequest(...)` | 解析请求 |
| 65 | `GenerateLearningObjectives(...)` | 生成目标 |
| 77 | `DiscoverContent(...)` | 发现内容 |
| 88 | `CustomizeWorkflow(...)` | 定制工作流 |

## 🎨 调试技巧

### 查看变量

```go
// 在断点处可以查看这些变量
userInput          // 用户输入数据
reqCtx             // 请求上下文
objectives         // 学习目标列表
contents           // 内容列表
workflowConfig     // 工作流配置
qualityReport      // 质量报告
```

### 添加临时打印

```go
// 在任何位置添加调试输出
fmt.Printf("DEBUG: userInput = %+v\n", userInput)
printJSON(objectives)  // 已有的 JSON 格式化函数
```

## 🧪 测试

```bash
# 运行测试
./debug_adapter.sh test

# 运行性能测试
./debug_adapter.sh bench

# 或直接运行
cd backend/domain/adapter/examples
go test -v
```

## 📊 其他命令

```bash
./debug_adapter.sh build      # 构建可执行文件
./debug_adapter.sh clean      # 清理构建产物
./debug_adapter.sh stats      # 代码统计
./debug_adapter.sh vscode     # VSCode 调试指南
./debug_adapter.sh profile    # 性能分析
./debug_adapter.sh help       # 帮助信息
```

## ⚠️ 重要提示

> **✅ demo.go 是独立程序**  
> 无需启动数据库、Redis 等任何服务！  
> 所有数据都在内存中，可以直接运行调试。

## 🆘 遇到问题？

```bash
# 1. 安装依赖
cd /home/liunix/work/coze-studio/backend
go mod tidy

# 2. 重新尝试
./debug_adapter.sh k12

# 3. 查看详细文档
cat docs/DEBUG_ADAPTER_DEMO.md
```

## 📚 相关文档

- **详细调试指南**: `docs/DEBUG_ADAPTER_DEMO.md`
- **快速启动**: `docs/DEBUG_QUICK_START.md`
- **完整调试文档**: `docs/DEBUG_COZE_SERVER.md`

---

**💡 提示**: 使用 `./debug_adapter.sh vscode` 查看 VSCode 调试的详细步骤！
