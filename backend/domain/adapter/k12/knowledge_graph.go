package k12
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

// MathKnowledgeGraph 数学知识图谱
var MathKnowledgeGraph = map[string]*models.KnowledgePoint{
	// 基础知识
	"real_numbers": {
		ID:          "real_numbers",
		Name:        "实数",
		Subject:     models.SubjectMath,
		Grade:       models.Grade7,
		Description: "实数的概念、分类和运算",
		Prerequisites: []string{},
		NextPoints:  []string{"equations", "inequalities"},
		Difficulty:  models.DifficultyBasic,
		Keywords:    []string{"有理数", "无理数", "实数运算"},
	},
	"equations": {
		ID:          "equations",
		Name:        "一元一次方程",
		Subject:     models.SubjectMath,
		Grade:       models.Grade7,
		Description: "一元一次方程的概念、解法和应用",
		Prerequisites: []string{"real_numbers"},
		NextPoints:  []string{"linear_equations_systems", "quadratic_equations"},
		Difficulty:  models.DifficultyBasic,
		Keywords:    []string{"方程", "解方程", "应用题"},
	},
	"linear_equations_systems": {
		ID:          "linear_equations_systems",
		Name:        "二元一次方程组",
		Subject:     models.SubjectMath,
		Grade:       models.Grade8,
		Description: "二元一次方程组的概念和解法",
		Prerequisites: []string{"equations"},
		NextPoints:  []string{"quadratic_equations"},
		Difficulty:  models.DifficultyImprove,
		Keywords:    []string{"方程组", "代入法", "加减法"},
	},
	// 目标知识点：二次函数
	"quadratic_equations": {
		ID:          "quadratic_equations",
		Name:        "一元二次方程",
		Subject:     models.SubjectMath,
		Grade:       models.Grade9,
		Description: "一元二次方程的概念、解法和应用",
		Prerequisites: []string{"equations", "linear_equations_systems"},
		NextPoints:  []string{"quadratic_functions"},
		Difficulty:  models.DifficultyImprove,
		Keywords:    []string{"二次方程", "求根公式", "判别式"},
	},
	"quadratic_functions": {
		ID:          "quadratic_functions",
		Name:        "二次函数",
		Subject:     models.SubjectMath,
		Grade:       models.Grade9,
		Description: "二次函数的概念、图像、性质和应用",
		Prerequisites: []string{"quadratic_equations", "linear_functions"},
		NextPoints:  []string{"quadratic_functions_advanced", "inverse_proportional_functions"},
		Difficulty:  models.DifficultyImprove,
		Keywords:    []string{"抛物线", "顶点", "对称轴", "最值"},
	},
	"linear_functions": {
		ID:          "linear_functions",
		Name:        "一次函数",
		Subject:     models.SubjectMath,
		Grade:       models.Grade8,
		Description: "一次函数的概念、图像和应用",
		Prerequisites: []string{"equations"},
		NextPoints:  []string{"quadratic_functions"},
		Difficulty:  models.DifficultyBasic,
		Keywords:    []string{"直线", "斜率", "截距"},
	},
	// 拓展知识
	"quadratic_functions_advanced": {
		ID:          "quadratic_functions_advanced",
		Name:        "二次函数综合应用",
		Subject:     models.SubjectMath,
		Grade:       models.Grade9,
		Description: "二次函数与几何、实际问题的综合应用",
		Prerequisites: []string{"quadratic_functions"},
		NextPoints:  []string{},
		Difficulty:  models.DifficultyAdvanced,
		Keywords:    []string{"综合应用", "数形结合", "实际问题"},
	},
	"inverse_proportional_functions": {
		ID:          "inverse_proportional_functions",
		Name:        "反比例函数",
		Subject:     models.SubjectMath,
		Grade:       models.Grade9,
		Description: "反比例函数的概念、图像和性质",
		Prerequisites: []string{"linear_functions"},
		NextPoints:  []string{},
		Difficulty:  models.DifficultyImprove,
		Keywords:    []string{"反比例", "双曲线", "渐近线"},
	},
}

// GetKnowledgePoint 获取知识点
func GetKnowledgePoint(id string) *models.KnowledgePoint {
	return MathKnowledgeGraph[id]
}

// GetPrerequisites 获取前置知识点列表
func GetPrerequisites(id string) []*models.KnowledgePoint {
	kp := MathKnowledgeGraph[id]
	if kp == nil {
		return nil
	}

	var prerequisites []*models.KnowledgePoint
	for _, preID := range kp.Prerequisites {
		if pre := MathKnowledgeGraph[preID]; pre != nil {
			prerequisites = append(prerequisites, pre)
		}
	}
	return prerequisites
}

// GetNextPoints 获取后续知识点列表
func GetNextPoints(id string) []*models.KnowledgePoint {
	kp := MathKnowledgeGraph[id]
	if kp == nil {
		return nil
	}

	var nextPoints []*models.KnowledgePoint
	for _, nextID := range kp.NextPoints {
		if next := MathKnowledgeGraph[nextID]; next != nil {
			nextPoints = append(nextPoints, next)
		}
	}
	return nextPoints
}
