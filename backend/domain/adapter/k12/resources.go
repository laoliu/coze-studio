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

	// === 物理学科资源 ===
	// 欧姆定律的学习资源
	"ohms_law": {
		{
			ID:          "video_ohm_basic",
			Type:        models.ResourceTypeVideo,
			Title:       "欧姆定律基础讲解",
			Description: "详细讲解欧姆定律的内容、公式推导和基本应用",
			URL:         "https://example.com/videos/ohms-law-basic",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyBasic,
			Duration:    18,
			Source:      "初中物理在线",
		},
		{
			ID:          "video_ohm_experiment",
			Type:        models.ResourceTypeVideo,
			Title:       "欧姆定律实验探究",
			Description: "通过实验演示欧姆定律，讲解电路图和数据分析",
			URL:         "https://example.com/videos/ohms-law-experiment",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyImprove,
			Duration:    22,
			Source:      "物理实验室",
		},
		{
			ID:          "textbook_ohm",
			Type:        models.ResourceTypeTextbook,
			Title:       "人教版九年级物理全一册 - 第17章 欧姆定律",
			Description: "教育部审定教材，系统讲解欧姆定律",
			URL:         "https://example.com/textbooks/grade9-physics-chapter17",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyBasic,
			Duration:    0,
			Source:      "人民教育出版社",
		},
	},

	// 电功率的学习资源
	"electrical_power": {
		{
			ID:          "video_power_basic",
			Type:        models.ResourceTypeVideo,
			Title:       "电功率基础概念",
			Description: "讲解电功、电功率的概念和计算公式",
			URL:         "https://example.com/videos/electrical-power-basic",
			KnowledgeID: "electrical_power",
			Difficulty:  models.DifficultyBasic,
			Duration:    16,
			Source:      "初中物理在线",
		},
		{
			ID:          "video_joule_law",
			Type:        models.ResourceTypeVideo,
			Title:       "焦耳定律及其应用",
			Description: "详细讲解焦耳定律的内容和实际应用",
			URL:         "https://example.com/videos/joule-law",
			KnowledgeID: "electrical_power",
			Difficulty:  models.DifficultyImprove,
			Duration:    20,
			Source:      "物理实验室",
		},
	},

	// 运动和力的学习资源
	"motion_and_forces": {
		{
			ID:          "video_newton_first",
			Type:        models.ResourceTypeVideo,
			Title:       "牛顿第一定律",
			Description: "讲解牛顿第一定律、惯性的概念和应用",
			URL:         "https://example.com/videos/newton-first-law",
			KnowledgeID: "motion_and_forces",
			Difficulty:  models.DifficultyBasic,
			Duration:    15,
			Source:      "初中物理在线",
		},
		{
			ID:          "video_force_balance",
			Type:        models.ResourceTypeVideo,
			Title:       "二力平衡",
			Description: "讲解二力平衡的条件和应用",
			URL:         "https://example.com/videos/force-balance",
			KnowledgeID: "motion_and_forces",
			Difficulty:  models.DifficultyImprove,
			Duration:    12,
			Source:      "初中物理在线",
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

	// === 物理习题 ===
	"ohms_law": {
		// 基础题
		{
			ID:          "ohm_ex_basic_1",
			Question:    "一段导体两端的电压为6V，通过的电流为0.3A，这段导体的电阻是多少？",
			Options:     []string{"A. 2Ω", "B. 20Ω", "C. 1.8Ω", "D. 18Ω"},
			Answer:      "B",
			Solution:    "根据欧姆定律 R = U/I = 6V / 0.3A = 20Ω",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
		{
			ID:          "ohm_ex_basic_2",
			Question:    "某导体的电阻为10Ω，当通过的电流为2A时，导体两端的电压是多少？",
			Options:     []string{"A. 5V", "B. 12V", "C. 20V", "D. 8V"},
			Answer:      "C",
			Solution:    "根据欧姆定律 U = IR = 2A × 10Ω = 20V",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
		// 提高题
		{
			ID:          "ohm_ex_improve_1",
			Question:    "一个滑动变阻器的最大阻值为50Ω，接入电路的阻值为20Ω，通过的电流为0.5A。若要使电流变为0.2A，电源电压不变，滑动变阻器接入电路的阻值应调为多少？",
			Options:     []string{},
			Answer:      "50Ω",
			Solution:    "电源电压 U = I₁R₁ = 0.5A × 20Ω = 10V\n所需电阻 R₂ = U/I₂ = 10V / 0.2A = 50Ω",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyImprove,
			Score:       10,
		},
		// 拓展题
		{
			ID:          "ohm_ex_advanced_1",
			Question:    "如图所示电路，电源电压恒为12V，R₁=10Ω，滑动变阻器R₂的最大阻值为50Ω。当开关S闭合，滑片P在中点时，电流表示数为0.3A。求：(1)R₂接入电路的阻值；(2)电路消耗的总功率。",
			Options:     []string{},
			Answer:      "(1) 30Ω  (2) 3.6W",
			Solution:    "(1) R₂中点阻值为 50Ω/2 = 25Ω\n总电阻 R总 = U/I = 12V/0.3A = 40Ω\n因为串联，R₁ + R₂ = 40Ω\n所以 R₂ = 40Ω - 10Ω = 30Ω\n(2) P = UI = 12V × 0.3A = 3.6W",
			KnowledgeID: "ohms_law",
			Difficulty:  models.DifficultyAdvanced,
			Score:       15,
		},
	},

	"electrical_power": {
		// 基础题
		{
			ID:          "power_ex_basic_1",
			Question:    "一个灯泡上标有\"220V 100W\"，其额定电压和额定功率分别是多少？",
			Options:     []string{"A. 220V, 100W", "B. 100V, 220W", "C. 110V, 50W", "D. 220V, 50W"},
			Answer:      "A",
			Solution:    "灯泡上的标识直接表示额定电压为220V，额定功率为100W",
			KnowledgeID: "electrical_power",
			Difficulty:  models.DifficultyBasic,
			Score:       3,
		},
		{
			ID:          "power_ex_basic_2",
			Question:    "某用电器的功率为1000W，使用时间为2小时，消耗的电能是多少度？",
			Options:     []string{"A. 0.5度", "B. 1度", "C. 2度", "D. 4度"},
			Answer:      "C",
			Solution:    "W = Pt = 1000W × 2h = 2000Wh = 2kWh = 2度",
			KnowledgeID: "electrical_power",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
		// 提高题
		{
			ID:          "power_ex_improve_1",
			Question:    "一个电热水壶，额定电压为220V，额定功率为1800W。求：(1)正常工作时的电流；(2)正常工作时的电阻。",
			Options:     []string{},
			Answer:      "(1) 8.18A  (2) 26.9Ω",
			Solution:    "(1) I = P/U = 1800W / 220V ≈ 8.18A\n(2) R = U²/P = (220V)² / 1800W ≈ 26.9Ω\n或 R = U/I = 220V / 8.18A ≈ 26.9Ω",
			KnowledgeID: "electrical_power",
			Difficulty:  models.DifficultyImprove,
			Score:       10,
		},
	},

	"motion_and_forces": {
		// 基础题
		{
			ID:          "force_ex_basic_1",
			Question:    "下列现象中，能说明物体具有惯性的是？",
			Options:     []string{"A. 汽车刹车后不能立即停下", "B. 物体在光滑表面上滑动", "C. 投出的篮球继续向前运动", "D. 以上都是"},
			Answer:      "D",
			Solution:    "惯性是物体保持原有运动状态的性质。A选项车继续前进，B选项物体继续滑动，C选项球继续飞行，都是惯性的表现。",
			KnowledgeID: "motion_and_forces",
			Difficulty:  models.DifficultyBasic,
			Score:       5,
		},
		// 提高题
		{
			ID:          "force_ex_improve_1",
			Question:    "一个物体受到两个力的作用，F₁ = 10N向右，F₂ = 10N向左。判断这两个力是否是平衡力，并说明理由。",
			Options:     []string{},
			Answer:      "是平衡力",
			Solution:    "这两个力满足二力平衡的条件：\n1. 大小相等（都是10N）\n2. 方向相反（一个向左，一个向右）\n3. 作用在同一物体上\n4. 作用在同一直线上\n因此这是一对平衡力。",
			KnowledgeID: "motion_and_forces",
			Difficulty:  models.DifficultyImprove,
			Score:       8,
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
