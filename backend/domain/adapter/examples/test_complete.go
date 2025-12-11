// Copyright 2025 Coze Studio. All rights reserved.
// 完整的 ModelBuilder 测试示例 - 包含配置初始化

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	"github.com/coze-dev/coze-studio/backend/infra/mysql"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
)

func main() {
	fmt.Println("=== Coze ModelBuilder 完整测试 ===\n")
	
	ctx := context.Background()
	
	// 步骤 1: 加载环境变量
	fmt.Println("步骤 1/4: 加载环境变量...")
	if err := loadEnvironment(); err != nil {
		log.Printf("警告: %v (将使用系统环境变量)\n", err)
	} else {
		fmt.Println("✓ 环境变量加载成功")
	}
	
	// 步骤 2: 初始化数据库
	fmt.Println("\n步骤 2/4: 初始化数据库连接...")
	db, err := mysql.New()
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v\n", err)
	}
	fmt.Println("✓ 数据库连接成功")
	
	// 步骤 3: 初始化存储
	fmt.Println("\n步骤 3/4: 初始化存储...")
	oss, err := storage.New(ctx)
	if err != nil {
		log.Fatalf("❌ 存储初始化失败: %v\n", err)
	}
	fmt.Println("✓ 存储初始化成功")
	
	// 步骤 4: 初始化配置
	fmt.Println("\n步骤 4/4: 初始化 Coze 配置...")
	if err := config.Init(ctx, db, oss); err != nil {
		log.Fatalf("❌ 配置初始化失败: %v\n", err)
	}
	fmt.Println("✓ Coze 配置初始化成功")
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("配置初始化完成，开始测试 ModelBuilder")
	fmt.Println(strings.Repeat("=", 50) + "\n")
	
	// 获取 Model ID
	modelIDStr := os.Getenv("COZE_MODEL_ID")
	if modelIDStr == "" {
		log.Fatal("❌ 请设置环境变量 COZE_MODEL_ID\n示例: export COZE_MODEL_ID=100002")
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		log.Fatalf("❌ 无效的 Model ID: %v", err)
	}

	fmt.Printf("Model ID: %d\n\n", modelID)
	
	// 使用 ModelBuilder 创建 ChatModel
	fmt.Println("正在使用 ModelBuilder 创建 ChatModel...")
	chatModel, toolsModel, err := modelbuilder.BuildModelByID(ctx, modelID, nil)
	if err != nil {
		log.Fatalf("❌ 创建 Model 失败: %v\n\n可能的原因:\n1. Model ID %d 不存在\n2. Model 配置不正确\n3. 缺少 API Key 配置\n\n提示: 运行 ./test_coze_db.sh 查看可用的 Models\n", err, modelID)
	}

	fmt.Println("\n✅ 成功创建 ChatModel!")
	fmt.Println(strings.Repeat("=", 50))
	
	if chatModel != nil {
		fmt.Println("✓ ChatModel: 可用 (用于对话)")
	}
	if toolsModel != nil {
		fmt.Println("✓ ToolsModel: 可用 (用于工具调用)")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("\n📚 下一步:")
	fmt.Println("  1. 使用此 ChatModel 进行对话")
	fmt.Println("  2. 参考 simple_coze_demo.go 查看完整示例")
	fmt.Println("  3. 访问 http://localhost:8888/admin 使用 Web UI")
}

func loadEnvironment() error {
	// 查找项目根目录的 .env 文件
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 向上查找 .env 文件
	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			return godotenv.Load(envFile)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return fmt.Errorf(".env 文件未找到")
}
