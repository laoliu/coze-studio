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

package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ExampleUsage 演示DAO的使用方法
func ExampleUsage() {
	// 1. 连接数据库
	dsn := "root:root@tcp(127.0.0.1:3307)/coze_studio?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 2. 创建DAO实例
	dao := NewK12DAO(db)
	ctx := context.Background()

	// ========== 示例1: 创建知识点 ==========
	fmt.Println("=== 示例1: 创建知识点 ===")
	kp := &KnowledgePointEntity{
		ID:          "ohms_law",
		Name:        "欧姆定律",
		Subject:     "physics",
		Grade:       "grade_9",
		Description: "电流、电压和电阻之间的关系",
		Keywords:    JSONStringArray{"欧姆定律", "电流", "电压", "电阻"},
		Prerequisites: JSONStringArray{"electric_current", "voltage"},
		NextPoints:    JSONStringArray{"series_circuit", "parallel_circuit"},
		Difficulty:    "basic",
		Metadata:      JSONMap{"chapter": "电学基础", "order": 1},
		Version:       1,
	}

	// 设置过期时间（7天后）
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	kp.ExpiresAt = &expiresAt

	err = dao.CreateKnowledgePoint(ctx, kp)
	if err != nil {
		fmt.Printf("创建知识点失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功创建知识点: %s\n", kp.Name)
	}

	// ========== 示例2: 查询知识点 ==========
	fmt.Println("\n=== 示例2: 查询知识点 ===")
	retrievedKP, err := dao.GetKnowledgePoint(ctx, "ohms_law")
	if err != nil {
		fmt.Printf("查询知识点失败: %v\n", err)
	} else {
		fmt.Printf("✅ 查到知识点: %s (学科: %s, 年级: %s)\n",
			retrievedKP.Name, retrievedKP.Subject, retrievedKP.Grade)
		fmt.Printf("   关键词: %v\n", retrievedKP.Keywords)
		fmt.Printf("   前置知识: %v\n", retrievedKP.Prerequisites)
	}

	// ========== 示例3: 按条件列出知识点 ==========
	fmt.Println("\n=== 示例3: 按条件列出知识点 ===")
	filter := KnowledgePointFilter{
		Subject: "physics",
		Grade:   "grade_9",
		Limit:   10,
	}
	kps, err := dao.ListKnowledgePoints(ctx, filter)
	if err != nil {
		fmt.Printf("列出知识点失败: %v\n", err)
	} else {
		fmt.Printf("✅ 找到 %d 个物理（初三）知识点\n", len(kps))
		for _, k := range kps {
			fmt.Printf("   - %s (难度: %s)\n", k.Name, k.Difficulty)
		}
	}

	// ========== 示例4: 批量创建学习资源 ==========
	fmt.Println("\n=== 示例4: 批量创建学习资源 ===")
	resources := []*LearningResourceEntity{
		{
			ID:          "ohms_law_video_1",
			KnowledgeID: "ohms_law",
			Type:        "video",
			Title:       "欧姆定律基础讲解",
			Description: "通过动画演示电流、电压和电阻的关系",
			URL:         "https://example.com/videos/ohms_law_basic",
			Difficulty:  "basic",
			Duration:    15,
			Source:      "Khan Academy",
		},
		{
			ID:          "ohms_law_textbook_1",
			KnowledgeID: "ohms_law",
			Type:        "textbook",
			Title:       "人教版物理教材 - 欧姆定律",
			Description: "教材第15章第2节",
			URL:         "https://example.com/textbooks/physics_grade9_ch15",
			Difficulty:  "basic",
			Source:      "人民教育出版社",
		},
	}

	err = dao.BatchCreateLearningResources(ctx, resources)
	if err != nil {
		fmt.Printf("批量创建资源失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功创建 %d 个学习资源\n", len(resources))
	}

	// ========== 示例5: 查询知识点的所有资源 ==========
	fmt.Println("\n=== 示例5: 查询知识点的所有资源 ===")
	kpResources, err := dao.GetResourcesByKnowledgeID(ctx, "ohms_law")
	if err != nil {
		fmt.Printf("查询资源失败: %v\n", err)
	} else {
		fmt.Printf("✅ 欧姆定律有 %d 个学习资源:\n", len(kpResources))
		for _, r := range kpResources {
			fmt.Printf("   - %s (%s, %d分钟)\n", r.Title, r.Type, r.Duration)
		}
	}

	// ========== 示例6: 批量创建习题 ==========
	fmt.Println("\n=== 示例6: 批量创建习题 ===")
	exercises := []*ExerciseEntity{
		{
			ID:          "ohms_law_ex_1",
			KnowledgeID: "ohms_law",
			Question:    "一个电阻为10Ω的电阻器，通过它的电流为2A，求两端的电压。",
			Options:     JSONStringArray{"A. 5V", "B. 10V", "C. 20V", "D. 40V"},
			Answer:      "C",
			Solution:    "根据欧姆定律 U = I × R = 2A × 10Ω = 20V",
			Difficulty:  "basic",
			Score:       5,
		},
		{
			ID:          "ohms_law_ex_2",
			KnowledgeID: "ohms_law",
			Question:    "在恒定电压12V的情况下，如何使电流从2A增加到3A？",
			Options:     JSONStringArray{"A. 增大电阻", "B. 减小电阻", "C. 保持电阻不变", "D. 断开电路"},
			Answer:      "B",
			Solution:    "由 R = U/I，当电压不变时，电流增大需要减小电阻",
			Difficulty:  "improve",
			Score:       10,
		},
	}

	err = dao.BatchCreateExercises(ctx, exercises)
	if err != nil {
		fmt.Printf("批量创建习题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功创建 %d 道习题\n", len(exercises))
	}

	// ========== 示例7: 按难度查询习题 ==========
	fmt.Println("\n=== 示例7: 按难度查询习题 ===")
	exerciseFilter := ExerciseFilter{
		KnowledgeID: "ohms_law",
		Difficulty:  "basic",
	}
	basicExercises, err := dao.ListExercises(ctx, exerciseFilter)
	if err != nil {
		fmt.Printf("查询习题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 找到 %d 道基础难度习题\n", len(basicExercises))
		for _, ex := range basicExercises {
			fmt.Printf("   - %s (分值: %d)\n", ex.Question[:30]+"...", ex.Score)
		}
	}

	// ========== 示例8: 记录生成日志 ==========
	fmt.Println("\n=== 示例8: 记录生成日志 ===")
	log := &GenerationLogEntity{
		RequestType:    "knowledge_point",
		Subject:        "physics",
		Grade:          "grade_9",
		KnowledgeID:    "ohms_law",
		PromptTemplate: "生成物理知识点...",
		LLMResponse:    `{"name":"欧姆定律","description":"..."}`,
		Success:        true,
		DurationMs:     1500,
	}

	err = dao.CreateGenerationLog(ctx, log)
	if err != nil {
		fmt.Printf("创建日志失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功记录生成日志 (ID: %d, 耗时: %dms)\n", log.ID, log.DurationMs)
	}

	// ========== 示例9: 查询最近的生成日志 ==========
	fmt.Println("\n=== 示例9: 查询最近的生成日志 ===")
	logFilter := GenerationLogFilter{
		Subject: "physics",
		Limit:   5,
	}
	logs, err := dao.ListGenerationLogs(ctx, logFilter)
	if err != nil {
		fmt.Printf("查询日志失败: %v\n", err)
	} else {
		fmt.Printf("✅ 找到 %d 条物理相关的生成日志\n", len(logs))
		for _, l := range logs {
			status := "✅ 成功"
			if !l.Success {
				status = "❌ 失败"
			}
			fmt.Printf("   - %s | %s | %dms | %s\n",
				l.RequestType, l.CreatedAt.Format("2006-01-02 15:04:05"),
				l.DurationMs, status)
		}
	}

	// ========== 示例10: 获取缓存配置 ==========
	fmt.Println("\n=== 示例10: 获取缓存配置 ===")
	cacheConfigs, err := dao.ListCacheConfigs(ctx)
	if err != nil {
		fmt.Printf("查询缓存配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 缓存配置:\n")
		for _, c := range cacheConfigs {
			days := c.TTLSeconds / 86400
			fmt.Printf("   - %s: %d天 (%s)\n", c.ItemType, days, c.Description)
		}
	}

	// ========== 示例11: 查询已过期的知识点 ==========
	fmt.Println("\n=== 示例11: 查询已过期的知识点 ===")
	expiredKPs, err := dao.GetExpiredKnowledgePoints(ctx, 10)
	if err != nil {
		fmt.Printf("查询过期知识点失败: %v\n", err)
	} else {
		if len(expiredKPs) == 0 {
			fmt.Println("✅ 没有过期的知识点")
		} else {
			fmt.Printf("⚠️  发现 %d 个过期知识点，需要刷新:\n", len(expiredKPs))
			for _, kp := range expiredKPs {
				fmt.Printf("   - %s (过期时间: %s)\n",
					kp.Name, kp.ExpiresAt.Format("2006-01-02 15:04:05"))
			}
		}
	}

	// ========== 示例12: 统计数据 ==========
	fmt.Println("\n=== 示例12: 统计数据 ===")
	totalKPs, _ := dao.CountKnowledgePoints(ctx, KnowledgePointFilter{})
	physicsKPs, _ := dao.CountKnowledgePoints(ctx, KnowledgePointFilter{Subject: "physics"})
	totalResources, _ := dao.CountResources(ctx, ResourceFilter{})
	totalExercises, _ := dao.CountExercises(ctx, ExerciseFilter{})

	fmt.Printf("✅ 数据统计:\n")
	fmt.Printf("   - 知识点总数: %d (物理: %d)\n", totalKPs, physicsKPs)
	fmt.Printf("   - 资源总数: %d\n", totalResources)
	fmt.Printf("   - 习题总数: %d\n", totalExercises)

	fmt.Println("\n=== 测试完成 ===")
}
