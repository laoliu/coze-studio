// Copyright 2025 Coze Studio. All rights reserved.
// 测试 ModelBuilder 的最小示例 - 带配置初始化

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
)

func main() {
	fmt.Println("=== Coze ModelBuilder 测试 ===")
	
	ctx := context.Background()
	
	// 初始化 Coze 配置
	fmt.Println("正在初始化 Coze 配置...")
	if err := InitCozeConfig(ctx); err != nil {
		log.Fatalf("❌ 初始化配置失败: %v\n", err)
	}
	fmt.Println("✓ 配置初始化成功\n")
	
	// 从环境变量获取 Model ID
	modelIDStr := os.Getenv("COZE_MODEL_ID")
	if modelIDStr == "" {
		log.Fatal("请设置环境变量 COZE_MODEL_ID\n示例: export COZE_MODEL_ID=100002")
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		log.Fatalf("无效的 Model ID: %v", err)
	}

	fmt.Printf("Model ID: %d\n", modelID)
	fmt.Println("正在加载 Model 配置...")
	
	// 使用 Coze ModelBuilder 创建 ChatModel
	chatModel, toolsModel, err := modelbuilder.BuildModelByID(ctx, modelID, nil)
	if err != nil {
		log.Fatalf("❌ 创建 Model 失败: %v\n\n可能的原因:\n1. Model ID 不存在\n2. 数据库连接失败\n3. Model 配置不正确\n", err)
	}

	fmt.Println("\n✅ 成功创建 ChatModel!")
	
	if chatModel != nil {
		fmt.Println("   - ChatModel: 可用")
	}
	if toolsModel != nil {
		fmt.Println("   - ToolsModel: 可用")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("\n下一步: 使用此 ChatModel 进行对话")
	fmt.Println("参考: ./simple_coze_demo.go")
}
