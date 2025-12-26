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

package k12

import (
	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// LearningResourceDB 学习资源数据库
var LearningResourceDB = map[string][]models.LearningResource{
	// 二次函数的学习资源
	"quadratic_functions": {
		// 视频资源
		{
			ID:          "video_qf_basic",
			Type:        models.ResourceTypeVideo,
			Title:       "二次函数基础概念",
			Description: "详细讲解二次函数的定义、一般式、顶点式和交点式",
			URL:         "https://example.com/videos/quadratic-functions-basic",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyBasic,
			Duration:    15,
			Source:      "中学数学在线",
		},
		{
			ID:          "video_qf_graph",
			Type:        models.ResourceTypeVideo,
			Title:       "二次函数的图像与性质",
			Description: "通过动画演示二次函数的图像变化和性质",
			URL:         "https://example.com/videos/quadratic-functions-graph",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyImprove,
			Duration:    20,
			Source:      "数学可视化",
		},
		// 教材资源
		{
			ID:          "textbook_qf",
			Type:        models.ResourceTypeTextbook,
			Title:       "人教版九年级数学上册 - 第22章 二次函数",
			Description: "教育部审定教材，系统讲解二次函数",
			URL:         "https://example.com/textbooks/grade9-math-chapter22",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyBasic,
			Duration:    0,
			Source:      "人民教育出版社",
		},
	},
	// 一元二次方程的学习资源
	"quadratic_equations": {
		{
			ID:          "video_qe_basic",
			Type:        models.ResourceTypeVideo,
			Title:       "一元二次方程的解法",
			Description: "讲解因式分解法、配方法和求根公式",
			URL:         "https://example.com/videos/quadratic-equations",
			KnowledgeID: "quadratic_equations",
			Difficulty:  models.DifficultyBasic,
			Duration:    18,
			Source:      "中学数学在线",
		},
	},
	// 一次函数的学习资源
	"linear_functions": {
		{
			ID:          "video_lf_basic",
			Type:        models.ResourceTypeVideo,
			Title:       "一次函数入门",
			Description: "从实际问题引入一次函数的概念",
			URL:         "https://example.com/videos/linear-functions",
			KnowledgeID: "linear_functions",
			Difficulty:  models.DifficultyBasic,
			Duration:    12,
			Source:      "中学数学在线",
		},
	},
}

// ExerciseDB 习题数据库
var ExerciseDB = map[string][]models.Exercise{
	"quadratic_functions": {
		// 基础题
		{
			ID:          "qf_ex_basic_1",
			Question:    "函数 y = x² - 4x + 3 的顶点坐标是多少？",
			Options:     []string{"A. (2, -1)", "B. (2, 1)", "C. (-2, -1)", "D. (-2, 1)"},
			Answer:      "A",
			Solution:    "将函数配方为顶点式：y = (x-2)² - 1，因此顶点坐标为 (2, -1)",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
		{
			ID:          "qf_ex_basic_2",
			Question:    "二次函数 y = -2x² 的图像开口方向是？",
			Options:     []string{"A. 向上", "B. 向下", "C. 向左", "D. 向右"},
			Answer:      "B",
			Solution:    "a = -2 < 0，所以抛物线开口向下",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyBasic,
			Score:       3,
		},
		// 提高题
		{
			ID:          "qf_ex_improve_1",
			Question:    "已知二次函数 y = ax² + bx + c 的图像过点 (0,3)、(1,0)、(-1,8)，求 a、b、c 的值。",
			Options:     []string{},
			Answer:      "a=3, b=-6, c=3",
			Solution:    "将三个点代入函数解析式，得到三元一次方程组：\nc=3\na+b+c=0\na-b+c=8\n解得 a=3, b=-6, c=3",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyImprove,
			Score:       10,
		},
		// 拓展题
		{
			ID:          "qf_ex_advanced_1",
			Question:    "某商店销售一种商品，成本为每件40元。经调查发现，当售价为50元时，每天可售出100件；售价每提高1元，每天少售出5件。问售价定为多少元时，每天的利润最大？最大利润是多少？",
			Options:     []string{},
			Answer:      "售价60元，最大利润2000元",
			Solution:    "设售价为x元，则销量为100-5(x-50)=350-5x件\n利润 y = (x-40)(350-5x) = -5x² + 550x - 14000\n配方得 y = -5(x-55)² + 1125\n但由于销量不能为负，350-5x≥0，x≤70\n当x=55时取得最大值，但需要验证实际约束条件\n在约束范围内，当x=60时，利润为(60-40)(350-5×60)=20×50=1000元\n...(完整解答)",
			KnowledgeID: "quadratic_functions",
			Difficulty:  models.DifficultyAdvanced,
			Score:       15,
		},
	},
	"quadratic_equations": {
		{
			ID:          "qe_ex_basic_1",
			Question:    "解方程：x² - 5x + 6 = 0",
			Options:     []string{},
			Answer:      "x₁=2, x₂=3",
			Solution:    "因式分解：(x-2)(x-3)=0\n所以 x-2=0 或 x-3=0\n解得 x₁=2, x₂=3",
			KnowledgeID: "quadratic_equations",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
	},
}

// GetResourcesByKnowledge 根据知识点ID获取学习资源
func GetResourcesByKnowledge(knowledgeID string, difficulty models.DifficultyLevel) []models.LearningResource {
	resources := LearningResourceDB[knowledgeID]
	if difficulty == "" {
		return resources
	}

	var filtered []models.LearningResource
	for _, resource := range resources {
		if resource.Difficulty == difficulty {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// GetExercisesByKnowledge 根据知识点ID和难度获取习题
func GetExercisesByKnowledge(knowledgeID string, difficulty models.DifficultyLevel) []models.Exercise {
	exercises := ExerciseDB[knowledgeID]
	if difficulty == "" {
		return exercises
	}

	var filtered []models.Exercise
	for _, exercise := range exercises {
		if exercise.Difficulty == difficulty {
			filtered = append(filtered, exercise)
		}
	}
	return filtered
}
