// Copyright 2025 Coze Studio. All rights reserved.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

// CreateK12Workflows 创建 K12 教学所需的 Workflow
// 这个脚本会在 Coze 数据库中创建预定义的 Workflow
func main() {
	fmt.Println("=== Creating K12 Workflows ===\n")

	// 从环境变量获取配置
	spaceID := getEnvInt64("SPACE_ID", 1)
	userID := getEnvInt64("USER_ID", 1)

	fmt.Printf("Space ID: %d\n", spaceID)
	fmt.Printf("User ID: %d\n\n", userID)

	ctx := context.Background()

	// 定义要创建的 Workflows
	workflows := []struct {
		Name        string
		Description string
		DesiredID   int64 // 期望的 ID（注意：实际 ID 由数据库自动生成）
	}{
		{
			Name:        "K12 通用教学",
			Description: "适用于 K12 通用教学活动的工作流，支持多种教学场景",
			DesiredID:   1000000,
		},
		{
			Name:        "K12 概念学习",
			Description: "专门用于概念理解的教学活动，包含导入、讲解、练习、总结等环节",
			DesiredID:   1000001,
		},
		{
			Name:        "K12 练习",
			Description: "用于知识巩固和技能训练的练习活动",
			DesiredID:   1000002,
		},
		{
			Name:        "K12 探索学习",
			Description: "引导学生自主探索和发现的探究式学习活动",
			DesiredID:   1000003,
		},
	}

	// 创建 Workflows
	fmt.Println("⚠️  注意：")
	fmt.Println("   Workflow ID 由数据库自动生成（AUTO_INCREMENT）")
	fmt.Println("   无法直接指定 ID = 1000000")
	fmt.Println("   需要使用以下方法之一：\n")
	fmt.Println("   方法 1: 手动调整数据库 AUTO_INCREMENT 起始值")
	fmt.Println("   方法 2: 创建后记录实际 ID，在代码中使用实际 ID")
	fmt.Println("   方法 3: 使用 Workflow 名称而不是 ID 来查找\n")

	fmt.Println("创建建议：")
	fmt.Println("─────────────────────────────────────────────\n")

	for i, wf := range workflows {
		fmt.Printf("%d. %s\n", i+1, wf.Name)
		fmt.Printf("   描述: %s\n", wf.Description)
		fmt.Printf("   期望 ID: %d (仅供参考)\n", wf.DesiredID)

		// 生成创建命令
		fmt.Printf("   \n   创建命令示例:\n")
		fmt.Printf("   curl -X POST http://localhost:8888/api/workflow_api/create \\\n")
		fmt.Printf("     -H 'Content-Type: application/json' \\\n")
		fmt.Printf("     -d '{\n")
		fmt.Printf("       \"name\": \"%s\",\n", wf.Name)
		fmt.Printf("       \"desc\": \"%s\",\n", wf.Description)
		fmt.Printf("       \"icon_uri\": \"workflow/k12\",\n")
		fmt.Printf("       \"space_id\": \"%d\"\n", spaceID)
		fmt.Printf("     }'\n\n")
	}

	fmt.Println("─────────────────────────────────────────────\n")
	fmt.Println("📝 下一步操作：")
	fmt.Println()
	fmt.Println("1. 使用 Admin UI 创建（推荐）")
	fmt.Println("   访问: http://localhost:8888/admin/workflows")
	fmt.Println("   点击 'Create Workflow' 按钮")
	fmt.Println("   填写名称和描述，保存后查看生成的 ID")
	fmt.Println()
	fmt.Println("2. 使用 API 创建")
	fmt.Println("   复制上面的 curl 命令执行")
	fmt.Println("   或使用下面的代码示例")
	fmt.Println()
	fmt.Println("3. 创建后记录实际 ID")
	fmt.Println("   将实际生成的 ID 更新到 adapter 代码中")
	fmt.Println()

	// 如果想要通过代码创建（需要初始化 Coze 环境）
	createViaCode := os.Getenv("CREATE_VIA_CODE")
	if createViaCode == "true" {
		fmt.Println("\n🔧 通过代码创建 Workflows...")
		fmt.Println("⚠️  需要先初始化 Coze 配置\n")

		// 注意：这里需要初始化 workflow service
		// 实际使用时需要完整的 Coze 环境初始化
		svc := workflow.DefaultSVC()
		if svc == nil {
			log.Fatal("Workflow service not initialized. Please run 'make web' first.")
		}

		for _, wf := range workflows {
			id, err := svc.Create(ctx, &vo.MetaCreate{
				Name:        wf.Name,
				Desc:        wf.Description,
				IconURI:     "workflow/k12",
				SpaceID:     spaceID,
				CreatorID:   userID,
				ContentType: 0, // Workflow
			})
			if err != nil {
				log.Printf("❌ Failed to create '%s': %v\n", wf.Name, err)
				continue
			}

			fmt.Printf("✓ Created '%s' with ID: %d\n", wf.Name, id)
			if id != wf.DesiredID {
				fmt.Printf("  ⚠️  Actual ID (%d) differs from desired ID (%d)\n", id, wf.DesiredID)
				fmt.Printf("  📝 Update your code to use ID: %d\n", id)
			}
		}
	}

	fmt.Println("\n=== Script Completed ===")
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}
