// Copyright 2025 Coze Studio. All rights reserved.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 0. 初始化 Coze 配置 (必须在使用 ModelBuilder 之前)
	fmt.Println("初始化 Coze 配置...")
	if err := InitCozeConfig(ctx); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	fmt.Println()

	// 1. 从环境变量获取 Model ID
	modelIDStr := os.Getenv("COZE_MODEL_ID")
	if modelIDStr == "" {
		log.Fatal("请设置环境变量 COZE_MODEL_ID")
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		log.Fatalf("无效的 Model ID: %v", err)
	}

	fmt.Printf("=== Coze 集成示例 ===\n")
	fmt.Printf("Model ID: %d\n\n", modelID)

	// 2. 使用 Coze ModelBuilder 创建 ChatModel
	chatModel, _, err := modelbuilder.BuildModelByID(ctx, modelID, nil)
	if err != nil {
		log.Fatalf("创建 Model 失败: %v", err)
	}

	fmt.Println("✓ 成功创建 ChatModel")

	// 3. 调用 LLM 生成内容
	prompt := "请为小学三年级数学课设计一个关于分数概念的学习目标"

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	fmt.Printf("\n用户提示: %s\n\n", prompt)
	fmt.Println("正在调用 LLM...")

	response, err := chatModel.Generate(ctx, messages, model.WithMaxTokens(500))
	if err != nil {
		log.Fatalf("LLM 调用失败: %v", err)
	}

	fmt.Printf("\nLLM 响应:\n%s\n\n", response.Content)

	// 4. JSON 格式输出示例
	jsonPrompt := `请生成一个 JSON 格式的学习目标，包含以下字段：
- id: 唯一标识
- objective: 目标描述
- level: 布鲁姆分类法级别（记忆/理解/应用/分析/评价/创造）
- category: 类别（知识/技能/态度）`

	messages2 := []*schema.Message{
		{
			Role:    schema.User,
			Content: jsonPrompt,
		},
	}

	fmt.Printf("JSON 格式请求:\n%s\n\n", jsonPrompt)
	fmt.Println("正在调用 LLM...")

	response2, err := chatModel.Generate(ctx, messages2, model.WithMaxTokens(300))
	if err != nil {
		log.Fatalf("LLM 调用失败: %v", err)
	}

	fmt.Printf("\nJSON 响应:\n%s\n\n", response2.Content)

	fmt.Println("=== 示例完成 ===")
}
