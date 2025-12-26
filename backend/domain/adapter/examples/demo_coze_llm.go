// Copyright 2025 Coze Studio. All rights reserved.
// Coze LLM 集成演示 - 展示如何使用 Coze ModelBuilder

package examples

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudwego/eino/schema"

	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
)

// DemoCozeK12Adapter 演示如何使用 Coze ModelBuilder 调用真实 LLM
func DemoCozeK12Adapter() {
	fmt.Println("=== Coze ModelBuilder 集成演示 ===\n")

	ctx := context.Background()

	// 步骤1: 初始化 Coze 配置
	fmt.Println("📋 步骤1: 初始化 Coze 配置")
	fmt.Println("正在加载数据库和存储配置...")
	if err := InitCozeConfig(ctx); err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		fmt.Println("\n💡 提示:")
		fmt.Println("   - 请确保数据库服务已启动")
		fmt.Println("   - 检查 .env 文件中的配置")
		fmt.Println("   - 确保有正确的数据库连接权限")
		return
	}
	fmt.Println("✓ Coze 配置初始化成功\n")

	// 步骤2: 获取 Model ID
	fmt.Println("📋 步骤2: 获取 Model ID")
	modelID, err := getModelIDFromEnv()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Println("\n💡 配置方法:")
		fmt.Println("   export COZE_MODEL_ID=<你的ModelID>")
		fmt.Println("\n如何获取 Model ID:")
		fmt.Println("   1. 登录 Coze Studio 管理后台")
		fmt.Println("   2. 进入 Model 管理页面")
		fmt.Println("   3. 找到你要使用的 Model，复制其 ID")
		return
	}
	fmt.Printf("✓ Model ID: %d\n\n", modelID)

	// 步骤3: 使用 ModelBuilder 构建 ChatModel
	fmt.Println("📋 步骤3: 构建 ChatModel")
	fmt.Println("正在从数据库加载 Model 配置...")
	chatModel, _, err := modelbuilder.BuildModelByID(ctx, modelID, nil)
	if err != nil {
		fmt.Printf("❌ 构建失败: %v\n", err)
		fmt.Println("\n💡 可能的原因:")
		fmt.Println("   - Model ID 不存在")
		fmt.Println("   - Model 配置不完整")
		fmt.Println("   - 数据库连接失败")
		return
	}
	fmt.Printf("✓ ChatModel 构建成功\n")
	fmt.Printf("   Model ID: %d\n\n", modelID)

	// 步骤4: 调用 LLM 生成内容
	fmt.Println("📋 步骤4: 调用 LLM 生成教学内容")
	fmt.Println("主题: 氧化还原反应（高一化学）\n")

	// 示例1: 生成学习目标
	fmt.Println("示例 1️⃣ : 生成学习目标")
	objectives, err := generateLearningObjectives(ctx, chatModel)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}
	fmt.Println("✓ 学习目标:")
	fmt.Println(objectives)
	fmt.Println()

	// 示例2: 推荐学习资源
	fmt.Println("示例 2️⃣ : 推荐学习资源")
	resources, err := recommendLearningResources(ctx, chatModel)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}
	fmt.Println("✓ 学习资源:")
	fmt.Println(resources)
	fmt.Println()

	// 示例3: 生成教学活动
	fmt.Println("示例 3️⃣ : 生成教学活动")
	activities, err := generateTeachingActivities(ctx, chatModel)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}
	fmt.Println("✓ 教学活动:")
	fmt.Println(activities)
	fmt.Println()

	// 完成
	fmt.Println("=== 演示完成 ===\n")
	fmt.Println("🎉 成功展示了以下功能:")
	fmt.Println("   ✓ 使用 Coze ModelBuilder 加载 Model 配置")
	fmt.Println("   ✓ 调用真实 LLM 生成教育内容")
	fmt.Println("   ✓ 展示 K12 教育场景的实际应用")
	fmt.Println()
	fmt.Println("💡 下一步:")
	fmt.Println("   - 可以修改提示词来生成不同的内容")
	fmt.Println("   - 可以尝试不同的主题和学科")
	fmt.Println("   - 可以集成到完整的适配器实现中")
}

// 生成学习目标
func generateLearningObjectives(ctx context.Context, chatModel modelbuilder.ToolCallingChatModel) (string, error) {
	prompt := `请为高一化学课程"氧化还原反应"生成3个清晰的学习目标。

要求:
1. 每个目标要具体、可衡量
2. 符合布鲁姆分类法
3. 适合45分钟课堂

请直接输出目标列表，每行一个。`

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM 调用失败: %w", err)
	}

	return resp.Content, nil
}

// 推荐学习资源
func recommendLearningResources(ctx context.Context, chatModel modelbuilder.ToolCallingChatModel) (string, error) {
	prompt := `请为"氧化还原反应"主题推荐3个优质学习资源。

要求:
1. 包括视频、文章、互动实验等不同类型
2. 说明每个资源的特点和适用场景
3. 优先推荐免费资源

请直接输出资源列表。`

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM 调用失败: %w", err)
	}

	return resp.Content, nil
}

// 生成教学活动
func generateTeachingActivities(ctx context.Context, chatModel modelbuilder.ToolCallingChatModel) (string, error) {
	prompt := `请为"氧化还原反应"设计一个15分钟的课堂活动。

要求:
1. 活动要有趣且有教育意义
2. 包含明确的步骤说明
3. 说明所需材料和预期效果

请直接输出活动方案。`

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM 调用失败: %w", err)
	}

	return resp.Content, nil
}

// getModelIDFromEnv 从环境变量获取 Model ID
func getModelIDFromEnv() (int64, error) {
	modelIDStr := os.Getenv("COZE_MODEL_ID")
	if modelIDStr == "" {
		return 0, fmt.Errorf("环境变量 COZE_MODEL_ID 未设置")
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的 Model ID: %s", modelIDStr)
	}

	return modelID, nil
}
