// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/entity"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/service"
)

// DemoK12Adapter 演示K12适配器的使用
func DemoK12Adapter() {
	fmt.Println("=== K12 Adapter Demo ===\n")
	
	// 1. 创建适配器管理器
	adapterMgr := service.NewAdapterManager()
	
	// 2. 创建并注册K12适配器
	k12Adapter := NewK12Adapter()
	err := adapterMgr.RegisterAdapter(k12Adapter)
	if err != nil {
		fmt.Printf("Failed to register adapter: %v\n", err)
		return
	}
	fmt.Println("✓ K12 Adapter registered successfully")
	
	// 3. 列出所有适配器
	adapters := adapterMgr.ListAdapters()
	fmt.Printf("✓ Total adapters: %d\n\n", len(adapters))
	
	// 4. 构建用户输入
	userInput := &entity.UserInput{
		Topic:        "氧化还原反应",
		Domain:       "化学",
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
	
	// 5. 解析请求
	ctx := context.Background()
	reqCtx, err := adapterMgr.ParseRequest(ctx, "k12_education", userInput)
	if err != nil {
		fmt.Printf("Failed to parse request: %v\n", err)
		return
	}
	fmt.Println("✓ Request parsed successfully")
	fmt.Println("Request Context:")
	printJSON(reqCtx)
	fmt.Println()
	
	// 6. 生成学习目标
	objectives, err := adapterMgr.GenerateLearningObjectives(ctx, "k12_education", reqCtx)
	if err != nil {
		fmt.Printf("Failed to generate objectives: %v\n", err)
		return
	}
	fmt.Printf("✓ Generated %d learning objectives\n", len(objectives))
	fmt.Println("Learning Objectives:")
	printJSON(objectives)
	fmt.Println()
	
	// 7. 发现内容
	contents, err := adapterMgr.DiscoverContent(ctx, "k12_education", reqCtx, objectives)
	if err != nil {
		fmt.Printf("Failed to discover content: %v\n", err)
		return
	}
	fmt.Printf("✓ Discovered %d content items\n", len(contents))
	fmt.Println("Contents:")
	printJSON(contents)
	fmt.Println()
	
	// 8. 定制工作流
	workflowConfig, err := adapterMgr.CustomizeWorkflow(ctx, "k12_education", reqCtx)
	if err != nil {
		fmt.Printf("Failed to customize workflow: %v\n", err)
		return
	}
	fmt.Println("✓ Workflow customized successfully")
	fmt.Println("Workflow Config:")
	printJSON(workflowConfig)
	fmt.Println()
	
	// 9. 模拟工作流执行结果
	workflowResult := &entity.WorkflowResult{
		WorkflowID: workflowConfig.WorkflowID,
		Status:     "success",
		Output: map[string]any{
			"title":       "氧化还原反应 - 概念理解",
			"objectives":  objectives,
			"content":     contents,
			"narrative":   "通过实例讲解氧化还原反应...",
			"assessment":  "概念测试、案例分析",
			"duration":    45,
			"materials":   []string{"PPT", "实验器材"},
			"instructions": "1. 导入新课\n2. 概念讲解\n3. 实例分析\n4. 练习巩固",
		},
		ExecutionTime: 2500,
	}
	
	// 10. 格式化输出
	output, err := adapterMgr.FormatOutput(ctx, "k12_education", workflowResult)
	if err != nil {
		fmt.Printf("Failed to format output: %v\n", err)
		return
	}
	fmt.Println("✓ Output formatted successfully")
	fmt.Println("Formatted Output:")
	printJSON(output)
	fmt.Println()
	
	// 11. 验证输出质量
	qualityReport, err := adapterMgr.ValidateOutput(ctx, "k12_education", output)
	if err != nil {
		fmt.Printf("Failed to validate output: %v\n", err)
		return
	}
	fmt.Println("✓ Output validated successfully")
	fmt.Println("Quality Report:")
	printJSON(qualityReport)
	fmt.Println()
	
	if qualityReport.IsApproved {
		fmt.Println("✓ Output approved for use!")
	} else {
		fmt.Println("✗ Output needs improvement")
	}
	
	fmt.Println("\n=== Demo Completed ===")
}

// printJSON 以格式化JSON打印对象
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// DemoAdapterRegistry 演示适配器注册表的使用
func DemoAdapterRegistry() {
	fmt.Println("=== Adapter Registry Demo ===\n")
	
	// 创建适配器管理器
	adapterMgr := service.NewAdapterManager()
	
	// 注册多个适配器
	k12Adapter := NewK12Adapter()
	adapterMgr.RegisterAdapter(k12Adapter)
	
	// 列出所有适配器
	adapters := adapterMgr.ListAdapters()
	fmt.Printf("Total adapters: %d\n\n", len(adapters))
	
	for _, info := range adapters {
		fmt.Printf("Adapter: %s\n", info.DisplayName)
		fmt.Printf("  ID: %s\n", info.AdapterID)
		fmt.Printf("  Version: %s\n", info.Version)
		fmt.Printf("  Type: %s\n", info.Type)
		fmt.Printf("  Supported Activity Types: %v\n", info.Capabilities.ActivityTypes)
		fmt.Printf("  Supported Domains: %v\n", info.Capabilities.Domains)
		fmt.Printf("  Supported Grades: %v\n\n", info.Capabilities.Grades)
	}
	
	// 测试自动选择适配器
	userInput := &entity.UserInput{
		Domain: "化学",
	}
	
	adapterID, err := adapterMgr.AutoSelectAdapter(context.Background(), userInput)
	if err != nil {
		fmt.Printf("Auto select failed: %v\n", err)
		return
	}
	
	fmt.Printf("Auto selected adapter: %s\n", adapterID)
	
	fmt.Println("\n=== Demo Completed ===")
}
