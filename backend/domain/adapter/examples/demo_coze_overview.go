// Copyright 2025 Coze Studio. All rights reserved.
// Coze LLM 演示 - 无需数据库版本（仅用于演示 API 调用流程）

package examples

import (
	"fmt"
)

// DemoCozeK12AdapterNoDB 演示 Coze 集成流程（不实际调用数据库）
// 用于展示如何使用 Coze ModelBuilder，但不需要真实的基础设施
func DemoCozeK12AdapterNoDB() {
	fmt.Println("=== Coze LLM 集成演示（演示版）===\n")

	fmt.Println("⚠️  提示：这是演示版本，展示集成流程但不调用真实服务\n")

	fmt.Println("📋 Real 版本需要的步骤:\n")

	// 步骤1
	fmt.Println("1️⃣  初始化 Coze 配置")
	fmt.Println("   - 加载 .env 环境变量")
	fmt.Println("   - 连接 MySQL 数据库")
	fmt.Println("   - 初始化存储服务（OSS/S3）")
	fmt.Println("   - 加载 Model 配置")
	fmt.Println()

	// 步骤2
	fmt.Println("2️⃣  获取 Model ID")
	fmt.Println("   - 从环境变量读取: COZE_MODEL_ID")
	fmt.Println("   - 示例: export COZE_MODEL_ID=12345")
	fmt.Println()

	// 步骤3
	fmt.Println("3️⃣  使用 ModelBuilder 构建 ChatModel")
	fmt.Println("   - 代码: chatModel, _, err := modelbuilder.BuildModelByID(ctx, modelID, nil)")
	fmt.Println("   - 从数据库加载 Model 配置（Provider、API Key、参数等）")
	fmt.Println("   - 返回可用的 ChatModel 实例")
	fmt.Println()

	// 步骤4
	fmt.Println("4️⃣  调用 LLM 生成内容")
	fmt.Println("   - 构造消息: messages := []*schema.Message{schema.UserMessage(prompt)}")
	fmt.Println("   - 调用 LLM: resp, err := chatModel.Generate(ctx, messages)")
	fmt.Println("   - 解析响应: resp.Content")
	fmt.Println()

	// 示例
	fmt.Println("📝 示例提示词:\n")

	fmt.Println("示例 1 - 生成学习目标:")
	fmt.Println("```")
	fmt.Println("请为高一化学课程\"氧化还原反应\"生成3个清晰的学习目标。")
	fmt.Println("要求：每个目标要具体、可衡量，符合布鲁姆分类法。")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("预期输出:")
	fmt.Println("  1. 理解氧化还原反应的基本概念和原理")
	fmt.Println("  2. 掌握氧化剂和还原剂的判断方法")
	fmt.Println("  3. 能够书写并配平简单的氧化还原反应方程式")
	fmt.Println()

	fmt.Println("示例 2 - 推荐学习资源:")
	fmt.Println("```")
	fmt.Println("请为\"氧化还原反应\"主题推荐3个优质学习资源。")
	fmt.Println("包括视频、文章、互动实验等不同类型。")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("预期输出:")
	fmt.Println("  1. 【视频】氧化还原反应动画演示（化学大师网）")
	fmt.Println("  2. 【文章】氧化还原反应知识点详解（人教版教材配套）")
	fmt.Println("  3. 【互动】在线化学实验室 - 氧化还原模拟器")
	fmt.Println()

	// 完成
	fmt.Println("=== 演示完成 ===\n")

	fmt.Println("🚀 如何运行真实版本:\n")
	fmt.Println("1. 启动数据库:")
	fmt.Println("   cd /home/liu/work/coze-studio")
	fmt.Println("   docker-compose -f docker/docker-compose-local.yml up -d mysql")
	fmt.Println()
	fmt.Println("2. 配置 Model:")
	fmt.Println("   - 在 Coze Admin 中创建 Model")
	fmt.Println("   - 配置 Provider (OpenAI/Claude/etc.)")
	fmt.Println("   - 获取 Model ID")
	fmt.Println()
	fmt.Println("3. 设置环境变量:")
	fmt.Println("   export COZE_MODEL_ID=<your_model_id>")
	fmt.Println()
	fmt.Println("4. 运行真实演示:")
	fmt.Println("   go run cmd/adapter-demo/main.go -demo=coze")
	fmt.Println()
	fmt.Println("💡 当前可用的演示:")
	fmt.Println("   - go run cmd/adapter-demo/main.go -demo=k12       (Mock版本，无需数据库)")
	fmt.Println("   - go run cmd/adapter-demo/main.go -demo=registry  (注册表演示)")
	fmt.Println("   - go run cmd/adapter-demo/main.go -demo=coze-demo (本演示)")
}
