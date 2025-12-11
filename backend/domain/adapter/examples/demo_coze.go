// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DemoCozeK12Adapter 演示使用 Coze 原生 modelbuilder 的 K12 适配器
// 这个版本直接使用 Coze 的 LLM 基础设施，适合在 Coze 环境中运行
func DemoCozeK12Adapter() {
	fmt.Println("=== Coze Integrated K12 Adapter Demo ===\n")

	ctx := context.Background()

	// 1. 从环境变量获取 Model ID
	modelID, err := getModelIDFromEnv()
	if err != nil {
		fmt.Printf("❌ Failed to get model ID: %v\n", err)
		fmt.Println("\n提示：请设置环境变量 COZE_MODEL_ID")
		fmt.Println("  export COZE_MODEL_ID=12345")
		fmt.Println("\n或在 Coze Admin 中配置 Model，然后获取 Model ID")
		return
	}
	fmt.Printf("✓ Model ID: %d\n", modelID)

	// 2. 初始化内容仓储
	contentRepo, err := initializeCozeContentRepository()
	if err != nil {
		fmt.Printf("❌ Failed to initialize content repository: %v\n", err)
		return
	}
	fmt.Println("✓ Content repository initialized")

	// 3. 创建 Coze 集成版 K12 适配器
	adapter, err := NewCozeK12AdapterByModelID(ctx, modelID, contentRepo)
	if err != nil {
		fmt.Printf("❌ Failed to create adapter: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("  1. Model ID 不存在")
		fmt.Println("  2. Coze 数据库未启动")
		fmt.Println("  3. Model 配置不完整")
		return
	}
	fmt.Println("✓ Coze K12 Adapter created\n")

	// 4. 运行适配器演示
	runCozeAdapterDemo(adapter)
}

// getModelIDFromEnv 从环境变量获取 Model ID
func getModelIDFromEnv() (int64, error) {
	modelIDStr := os.Getenv("COZE_MODEL_ID")
	if modelIDStr == "" {
		// 尝试从配置文件读取
		return getModelIDFromConfig()
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid COZE_MODEL_ID: %v", err)
	}

	return modelID, nil
}

// getModelIDFromConfig 从配置文件获取默认 Model ID
func getModelIDFromConfig() (int64, error) {
	// TODO: 从 Coze 配置文件读取默认 Model
	// 这里暂时返回错误，提示用户设置环境变量
	return 0, fmt.Errorf("COZE_MODEL_ID not set")
}

// initializeCozeContentRepository 初始化连接到 Coze 数据库的内容仓储
func initializeCozeContentRepository() (repository.ContentRepository, error) {
	// 检查是否使用 Mock 模式
	if os.Getenv("DEBUG_ENABLE_MOCK_DB") == "true" {
		fmt.Println("  Using Mock Content Repository (Debug Mode)")
		return setupMockContentRepository(), nil
	}

	// 连接 Coze 数据库
	dsn := getCozeDBDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Coze database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	return repository.NewGormContentRepository(db), nil
}

// getCozeDBDSN 获取 Coze 数据库连接字符串
func getCozeDBDSN() string {
	// 从环境变量读取，如果没有则使用默认值
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "coze")
	password := getEnvOrDefault("DB_PASSWORD", "coze123")
	database := getEnvOrDefault("DB_NAME", "opencoze")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
}

// runCozeAdapterDemo 运行 Coze 适配器演示
func runCozeAdapterDemo(adapter *CozeK12Adapter) {
	// 1. 显示适配器信息
	info := adapter.GetInfo()
	fmt.Println("Adapter Info:")
	printJSON(info)
	fmt.Println()

	// 2. 用户输入
	userInput := &entity.UserInput{
		Topic:        "光合作用",
		Domain:       "生物",
		Grade:        "高一",
		Duration:     45,
		ActivityType: entity.ActivityTypeConcept,
		Language:     "zh-CN",
		ExtraParams: map[string]any{
			"difficulty": "medium",
		},
	}

	fmt.Println("User Input:")
	printJSON(userInput)
	fmt.Println()

	// 3. 解析请求
	reqCtx, err := adapter.ParseRequest(userInput)
	if err != nil {
		fmt.Printf("❌ Parse request failed: %v\n", err)
		return
	}
	fmt.Println("✓ Request parsed")

	// 4. 生成学习目标（调用 Coze LLM）
	fmt.Println("\n🤖 Generating learning objectives with Coze LLM...")
	objectives, err := adapter.GenerateLearningObjectives(reqCtx)
	if err != nil {
		fmt.Printf("❌ Generate objectives failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Generated %d objectives\n", len(objectives))
	printJSON(objectives)
	fmt.Println()

	// 5. 发现内容（查询 Coze 数据库）
	fmt.Println("🔍 Discovering content from Coze database...")
	contents, err := adapter.DiscoverContent(reqCtx, objectives)
	if err != nil {
		fmt.Printf("❌ Discover content failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Discovered %d content items\n", len(contents))
	printJSON(contents)
	fmt.Println()

	// 6. 准备 AdapterContext（包含必要的元数据）
	adapterCtx := &entity.AdapterContext{
		Context: context.Background(),
		Input:   userInput,
		Domain:  userInput.Domain,
		Grade:   10, // 高一
		Metadata: map[string]any{
			"user_id":       int64(1),          // 演示用户 ID
			"connector_id":  int64(1),          // 演示 Connector ID
			"connector_uid": "demo_user",       // 演示用户标识
		},
	}

	// 7. 执行 Coze 原生 Workflow
	fmt.Println("\n🚀 Executing Coze native Workflow...")
	fmt.Println("⚠️  注意：这需要在 Coze 系统中预先创建对应的 Workflow")
	fmt.Println("   可以通过 Admin -> Workflows 界面创建以下 Workflow：")
	fmt.Println("   - ID 1000000: K12 通用教学 Workflow")
	fmt.Println("   - ID 1000001: K12 概念学习 Workflow")
	fmt.Println("   - ID 1000002: K12 练习 Workflow")
	fmt.Println("   - ID 1000003: K12 探索学习 Workflow")
	fmt.Println()

	workflowResult, err := adapter.ExecuteWorkflow(adapterCtx)
	if err != nil {
		fmt.Printf("⚠️  Workflow execution failed: %v\n", err)
		fmt.Println("   这是预期的，因为演示环境中可能没有创建对应的 Workflow")
		fmt.Println("   在生产环境中，需要先创建这些 Workflow")
		fmt.Println()
		
		// 为了演示继续，使用 mock 数据
		fmt.Println("📝 Using mock workflow result for demo continuation...")
		workflowResult = &entity.WorkflowResult{
			WorkflowID: "1000001",
			Status:     "success",
			Output: map[string]any{
				"topic":       userInput.Topic,
				"domain":      userInput.Domain,
				"grade":       "10",
				"title":       "光合作用 - 概念理解",
				"objectives":  objectives,
				"content":     contents,
				"duration":    userInput.Duration,
				"assessment":  "概念测试、实验观察",
				"materials":   []string{"显微镜", "植物叶片", "PPT"},
				"instructions": "1. 导入\n2. 概念讲解\n3. 实验观察\n4. 总结",
			},
			ExecutionTime: 3500,
		}
	} else {
		fmt.Println("✓ Workflow executed successfully")
		printJSON(workflowResult)
		fmt.Println()
	}

	// 8. 格式化输出（使用 Coze LLM 生成叙事）
	fmt.Println("\n🤖 Formatting output with Coze LLM...")
	output, err := adapter.FormatOutput(workflowResult)
	if err != nil {
		fmt.Printf("❌ Format output failed: %v\n", err)
		return
	}
	fmt.Println("✓ Output formatted")
	printJSON(output)
	fmt.Println()

	// 9. 验证质量（使用 Coze LLM）
	fmt.Println("🤖 Validating output quality with Coze LLM...")
	qualityReport, err := adapter.ValidateOutput(output)
	if err != nil {
		fmt.Printf("❌ Validate output failed: %v\n", err)
		return
	}
	fmt.Println("✓ Quality validated")
	printJSON(qualityReport)
	fmt.Println()

	if qualityReport.IsApproved {
		fmt.Println("✅ Output approved!")
	} else {
		fmt.Println("⚠️  Output needs improvement")
	}

	fmt.Println("\n=== Demo Completed ===")
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
