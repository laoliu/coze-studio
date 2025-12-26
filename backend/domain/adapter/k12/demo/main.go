package main
/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/k12"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

func main() {
	fmt.Println("=== K12教育适配器演示 ===\n")

	// 创建K12适配器服务
	service := k12.NewK12AdapterService()

	// 示例用户输入
	userInput := "我想学习初中数学的二次函数"
	fmt.Printf("用户输入: %s\n\n", userInput)

	// 假设用户已经掌握的知识点
	masteredKnowledgeIDs := []string{
		"real_numbers",           // 实数
		"equations",              // 一元一次方程
		"linear_equations_systems", // 二元一次方程组
	}

	// 1. 分析用户输入
	fmt.Println("【步骤1：分析用户输入】")
	analysisResult, err := service.AnalyzeRequest(&models.AnalysisRequest{
		UserInput: userInput,
	})
	if err != nil {
		log.Fatalf("分析失败: %v", err)
	}

	fmt.Printf("✅ 分析结果：\n")
	fmt.Printf("   学科: %s\n", analysisResult.Subject)
	fmt.Printf("   年级: %s\n", analysisResult.Grade)
	fmt.Printf("   知识点: %s (%s)\n", analysisResult.KnowledgeName, analysisResult.KnowledgeID)
	fmt.Printf("   前置知识: %v\n", analysisResult.Prerequisites)
	fmt.Printf("   置信度: %.2f\n\n", analysisResult.Confidence)

	// 2. 生成学习计划
	fmt.Println("【步骤2：生成个性化学习计划】")
	learningPlan, err := service.GenerateLearningPlan(userInput, masteredKnowledgeIDs)
	if err != nil {
		log.Fatalf("生成学习计划失败: %v", err)
	}

	fmt.Printf("✅ 学习计划生成成功\n\n")

	// 3. 输出学习计划详情
	fmt.Println("【学习计划详情】")
	fmt.Printf("学生水平: %s\n", learningPlan.StudentLevel)
	fmt.Printf("目标知识点: %s\n", learningPlan.TargetKnowledge.Name)
	fmt.Printf("总学习时间: %d 分钟\n\n", learningPlan.TotalTime)

	// 4. 输出前置知识
	if len(learningPlan.Prerequisites) > 0 {
		fmt.Println("📚 前置知识：")
		for i, kp := range learningPlan.Prerequisites {
			status := "❌ 未掌握"
			for _, masteredID := range masteredKnowledgeIDs {
				if kp.ID == masteredID {
					status = "✅ 已掌握"
					break
				}
			}
			fmt.Printf("   %d. %s (%s) - %s\n", i+1, kp.Name, kp.Grade, status)
		}
		fmt.Println()
	}

	// 5. 输出学习路径
	fmt.Println("🎯 学习路径：")
	for _, step := range learningPlan.MainPath.Steps {
		fmt.Printf("\n步骤 %d: %s\n", step.Order, step.Knowledge.Name)
		fmt.Printf("  难度: %s\n", step.Knowledge.Difficulty)
		fmt.Printf("  预计时间: %d 分钟\n", step.EstimatedTime)

		if len(step.Objectives) > 0 {
			fmt.Println("  学习目标:")
			for _, obj := range step.Objectives {
				fmt.Printf("    • %s\n", obj)
			}
		}

		if len(step.Resources) > 0 {
			fmt.Println("  推荐资源:")
			for _, res := range step.Resources {
				fmt.Printf("    • [%s] %s (%d分钟)\n", res.Type, res.Title, res.Duration)
			}
		}

		if len(step.Exercises) > 0 {
			fmt.Printf("  配套习题: %d 道\n", len(step.Exercises))
		}
	}
	fmt.Println()

	// 6. 输出拓展内容
	if len(learningPlan.Extension) > 0 {
		fmt.Println("🚀 拓展学习：")
		for i, kp := range learningPlan.Extension {
			fmt.Printf("   %d. %s (%s - %s)\n", i+1, kp.Name, kp.Grade, kp.Difficulty)
		}
		fmt.Println()
	}

	// 7. 输出完整的JSON（可用于API响应）
	fmt.Println("【完整JSON输出】")
	jsonData, err := json.MarshalIndent(learningPlan, "", "  ")
	if err != nil {
		log.Fatalf("JSON序列化失败: %v", err)
	}
	fmt.Println(string(jsonData))

	fmt.Println("\n=== 演示完成 ===")
}
