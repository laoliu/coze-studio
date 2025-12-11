// Copyright 2025 Coze Studio. All rights reserved.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	// 导入实际的 adapter 包（需要在 examples 中导出）
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/prompts"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	crossworkflow "github.com/coze-dev/coze-studio/backend/crossdomain/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"strconv"
)

// DemoWorkflow 演示使用 Coze 原生 Workflow 的完整流程
func main() {
	fmt.Println("=== Coze Native Workflow Integration Demo ===\n")

	// 1. 初始化 Coze 配置
	fmt.Println("🔧 Initializing Coze configuration...")
	if err := InitCozeConfig(); err != nil {
		fmt.Printf("❌ Failed to initialize Coze config: %v\n", err)
		fmt.Println("⚠️  提示：请确保 Coze 环境已启动（make web）")
		return
	}
	fmt.Println("✓ Coze configuration initialized")
	fmt.Println()

	// 2. 设置 Content Repository
	contentRepo, err := setupContentRepository()
	if err != nil {
		fmt.Printf("❌ Setup content repository failed: %v\n", err)
		return
	}
	fmt.Println("✓ Content repository initialized")
	fmt.Println()

	// 3. 创建适配器（使用 Coze Model ID）
	fmt.Println("🤖 Creating Coze K12 Adapter...")
	// 这里使用一个示例 Model ID，实际使用时需要替换为真实的 ID
	// 可以在 Coze Admin -> Models 中查看和配置
	modelID := int64(1) // 示例 Model ID

	adapter, err := examples.NewCozeK12AdapterByModelID(context.Background(), modelID, contentRepo)
	if err != nil {
		fmt.Printf("❌ Create adapter failed: %v\n", err)
		fmt.Println("⚠️  提示：请确保 Model ID 在 Coze 数据库中存在")
		return
	}
	fmt.Println("✓ Adapter created successfully")
	fmt.Println()

	// 4. 运行演示
	runWorkflowDemo(adapter)
}

// setupContentRepository 设置内容仓库
func setupContentRepository() (repository.ContentRepository, error) {
	dsn := getCozeDBDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

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
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "coze")
	password := getEnvOrDefault("DB_PASSWORD", "coze123")
	database := getEnvOrDefault("DB_NAME", "opencoze")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
}

// runWorkflowDemo 运行 Workflow 集成演示
func runWorkflowDemo(adapter *examples.CozeK12Adapter) {
	// 1. 显示适配器信息
	info := adapter.GetInfo()
	fmt.Println("📋 Adapter Info:")
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

	fmt.Println("📝 User Input:")
	printJSON(userInput)
	fmt.Println()

	// 3. 解析请求
	reqCtx, err := adapter.ParseRequest(userInput)
	if err != nil {
		fmt.Printf("❌ Parse request failed: %v\n", err)
		return
	}
	fmt.Println("✓ Request parsed")

	// 4. 生成学习目标（使用 Coze LLM）
	fmt.Println("\n🤖 Generating learning objectives with Coze ModelBuilder...")
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
			"user_id":       int64(1),
			"connector_id":  int64(1),
			"connector_uid": "demo_user",
		},
	}

	// 7. 执行 Coze 原生 Workflow（核心部分）
	fmt.Println("\n🚀 Executing Coze Native Workflow...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  重要提示：")
	fmt.Println("   此功能需要在 Coze 系统中预先创建对应的 Workflow")
	fmt.Println("   可以通过以下方式创建：")
	fmt.Println("   1. Admin UI: http://localhost:8888/admin/workflows")
	fmt.Println("   2. API: 使用 domain/workflow 的 Create 方法")
	fmt.Println()
	fmt.Println("   需要创建的 Workflow ID：")
	fmt.Println("   - 1000000: K12 通用教学 Workflow")
	fmt.Println("   - 1000001: K12 概念学习 Workflow (当前使用)")
	fmt.Println("   - 1000002: K12 练习 Workflow")
	fmt.Println("   - 1000003: K12 探索学习 Workflow")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	workflowResult, err := adapter.ExecuteWorkflow(adapterCtx)
	if err != nil {
		fmt.Printf("⚠️  Workflow execution failed: %v\n", err)
		fmt.Println()
		fmt.Println("💡 这是预期的行为，因为演示环境中可能没有创建对应的 Workflow")
		fmt.Println("   在生产环境中使用前，请先通过 Admin UI 或 API 创建 Workflow")
		fmt.Println()
		fmt.Println("📝 继续使用 mock 数据演示后续流程...")
		fmt.Println()
		
		// 使用 mock 数据继续演示
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
		fmt.Println("🎉 Workflow executed successfully!")
		fmt.Println("📊 Workflow Result:")
		printJSON(workflowResult)
		fmt.Println()
	}

	// 8. 格式化输出（使用 Coze LLM 生成叙事）
	fmt.Println("🤖 Formatting output with Coze ModelBuilder...")
	output, err := adapter.FormatOutput(workflowResult)
	if err != nil {
		fmt.Printf("❌ Format output failed: %v\n", err)
		return
	}
	fmt.Println("✓ Output formatted")
	printJSON(output)
	fmt.Println()

	// 9. 验证质量（使用 Coze LLM）
	fmt.Println("🤖 Validating output quality with Coze ModelBuilder...")
	qualityReport, err := adapter.ValidateOutput(output)
	if err != nil {
		fmt.Printf("❌ Validate output failed: %v\n", err)
		return
	}
	fmt.Println("✓ Quality validated")
	printJSON(qualityReport)
	fmt.Println()

	if qualityReport.IsApproved {
		fmt.Println("✅ Output approved - Quality meets standards!")
	} else {
		fmt.Println("⚠️  Output needs improvement")
		fmt.Println("📝 Suggestions:")
		for i, suggestion := range qualityReport.Suggestions {
			fmt.Printf("   %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ Demo Completed Successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("📚 总结：")
	fmt.Println("   ✓ LLM: 使用 Coze ModelBuilder (bizpkg/llm/modelbuilder)")
	fmt.Println("   ✓ Workflow: 使用 Coze 原生 Workflow (crossdomain/workflow)")
	fmt.Println("   ✓ Database: 使用 Coze 数据库")
	fmt.Println()
	fmt.Println("🎯 完全基于 Coze 原生机制，无独立实现！")
}

// printJSON 以格式化 JSON 打印
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
