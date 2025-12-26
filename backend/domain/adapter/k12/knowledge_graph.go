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

// PhysicsKnowledgeGraph 物理知识图谱
var PhysicsKnowledgeGraph = map[string]*models.KnowledgePoint{
	// 初中物理 - 力学基础
	"measurement_motion": {
		ID:            "measurement_motion",
		Name:          "长度和时间的测量",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "长度、时间的测量方法和单位换算",
		Prerequisites: []string{},
		NextPoints:    []string{"motion_description", "speed"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"测量", "刻度尺", "停表", "单位换算"},
	},
	"speed": {
		ID:            "speed",
		Name:          "速度",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "速度的概念、计算和应用",
		Prerequisites: []string{"measurement_motion"},
		NextPoints:    []string{"forces", "motion_and_forces"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"速度", "平均速度", "匀速运动"},
	},
	"forces": {
		ID:            "forces",
		Name:          "力的概念",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "力的概念、三要素、测量和图示",
		Prerequisites: []string{"speed"},
		NextPoints:    []string{"motion_and_forces", "pressure"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"力", "弹力", "重力", "摩擦力"},
	},
	"motion_and_forces": {
		ID:            "motion_and_forces",
		Name:          "运动和力",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "牛顿第一定律、惯性、二力平衡",
		Prerequisites: []string{"speed", "forces"},
		NextPoints:    []string{"pressure", "work_and_energy"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"牛顿第一定律", "惯性", "二力平衡"},
	},
	"pressure": {
		ID:            "pressure",
		Name:          "压强",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "压强的概念、计算、液体压强和大气压强",
		Prerequisites: []string{"forces"},
		NextPoints:    []string{"buoyancy"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"压强", "液体压强", "大气压强"},
	},
	"buoyancy": {
		ID:            "buoyancy",
		Name:          "浮力",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "浮力的概念、阿基米德原理和应用",
		Prerequisites: []string{"pressure"},
		NextPoints:    []string{"simple_machines"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"浮力", "阿基米德原理", "物体浮沉条件"},
	},
	"work_and_energy": {
		ID:            "work_and_energy",
		Name:          "功和机械能",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "功、功率、动能、势能和机械能守恒",
		Prerequisites: []string{"motion_and_forces"},
		NextPoints:    []string{"simple_machines"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"功", "功率", "动能", "势能", "机械能守恒"},
	},
	"simple_machines": {
		ID:            "simple_machines",
		Name:          "简单机械",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "杠杆、滑轮、机械效率",
		Prerequisites: []string{"work_and_energy", "buoyancy"},
		NextPoints:    []string{},
		Difficulty:    models.DifficultyAdvanced,
		Keywords:      []string{"杠杆", "滑轮", "机械效率"},
	},

	// 初中物理 - 电学
	"electricity_basics": {
		ID:            "electricity_basics",
		Name:          "电流和电路",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade9,
		Description:   "电流、电路、串并联电路",
		Prerequisites: []string{},
		NextPoints:    []string{"voltage_resistance"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"电流", "电路", "串联", "并联"},
	},
	"voltage_resistance": {
		ID:            "voltage_resistance",
		Name:          "电压和电阻",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade9,
		Description:   "电压、电阻的概念和测量",
		Prerequisites: []string{"electricity_basics"},
		NextPoints:    []string{"ohms_law"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"电压", "电阻", "电压表", "滑动变阻器"},
	},
	"ohms_law": {
		ID:            "ohms_law",
		Name:          "欧姆定律",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade9,
		Description:   "欧姆定律及其应用",
		Prerequisites: []string{"voltage_resistance"},
		NextPoints:    []string{"electrical_power", "circuit_analysis"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"欧姆定律", "电流计算", "电阻计算"},
	},
	"electrical_power": {
		ID:            "electrical_power",
		Name:          "电功率",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade9,
		Description:   "电功、电功率、焦耳定律",
		Prerequisites: []string{"ohms_law"},
		NextPoints:    []string{"circuit_analysis"},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"电功", "电功率", "焦耳定律", "额定功率"},
	},
	"circuit_analysis": {
		ID:            "circuit_analysis",
		Name:          "电路综合分析",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade9,
		Description:   "复杂电路的分析和计算",
		Prerequisites: []string{"ohms_law", "electrical_power"},
		NextPoints:    []string{},
		Difficulty:    models.DifficultyAdvanced,
		Keywords:      []string{"动态电路", "电路故障", "综合计算"},
	},

	// 初中物理 - 光学和热学
	"light_reflection": {
		ID:            "light_reflection",
		Name:          "光的反射",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "光的直线传播、反射定律、平面镜成像",
		Prerequisites: []string{},
		NextPoints:    []string{"light_refraction"},
		Difficulty:    models.DifficultyBasic,
		Keywords:      []string{"光的反射", "反射定律", "平面镜"},
	},
	"light_refraction": {
		ID:            "light_refraction",
		Name:          "光的折射",
		Subject:       models.SubjectPhysics,
		Grade:         models.Grade8,
		Description:   "光的折射定律、透镜及其应用",
		Prerequisites: []string{"light_reflection"},
		NextPoints:    []string{},
		Difficulty:    models.DifficultyImprove,
		Keywords:      []string{"光的折射", "透镜", "凸透镜成像"},
	},
}

// AllKnowledgeGraphs 所有学科的知识图谱
var AllKnowledgeGraphs = map[models.Subject]map[string]*models.KnowledgePoint{
	models.SubjectMath:    MathKnowledgeGraph,
	models.SubjectPhysics: PhysicsKnowledgeGraph,
}

// GetKnowledgePoint 获取知识点（支持多学科）
func GetKnowledgePoint(id string) *models.KnowledgePoint {
	// 先尝试从数学知识图谱查找
	if kp := MathKnowledgeGraph[id]; kp != nil {
		return kp
	}
	// 再从物理知识图谱查找
	if kp := PhysicsKnowledgeGraph[id]; kp != nil {
		return kp
	}
	return nil
}

// GetKnowledgePointBySubject 根据学科获取知识点
func GetKnowledgePointBySubject(subject models.Subject, id string) *models.KnowledgePoint {
	if graph, ok := AllKnowledgeGraphs[subject]; ok {
		return graph[id]
	}
	return nil
}

// GetPrerequisites 获取前置知识点列表（支持多学科）
func GetPrerequisites(id string) []*models.KnowledgePoint {
	kp := GetKnowledgePoint(id)
	if kp == nil {
		return nil
	}

	var prerequisites []*models.KnowledgePoint
	for _, preID := range kp.Prerequisites {
		if pre := GetKnowledgePoint(preID); pre != nil {
			prerequisites = append(prerequisites, pre)
		}
	}
	return prerequisites
}

// GetNextPoints 获取后续知识点列表（支持多学科）
func GetNextPoints(id string) []*models.KnowledgePoint {
	kp := GetKnowledgePoint(id)
	if kp == nil {
		return nil
	}

	var nextPoints []*models.KnowledgePoint
	for _, nextID := range kp.NextPoints {
		if next := GetKnowledgePoint(nextID); next != nil {
			nextPoints = append(nextPoints, next)
		}
	}
	return nextPoints
}

// ListKnowledgePointsBySubject 列出指定学科的所有知识点
func ListKnowledgePointsBySubject(subject models.Subject) []*models.KnowledgePoint {
	graph, ok := AllKnowledgeGraphs[subject]
	if !ok {
		return nil
	}

	var points []*models.KnowledgePoint
	for _, kp := range graph {
		points = append(points, kp)
	}
	return points
}

// ListKnowledgePointsByGrade 列出指定年级的知识点
func ListKnowledgePointsByGrade(subject models.Subject, grade models.GradeLevel) []*models.KnowledgePoint {
	graph, ok := AllKnowledgeGraphs[subject]
	if !ok {
		return nil
	}

	var points []*models.KnowledgePoint
	for _, kp := range graph {
		if kp.Grade == grade {
			points = append(points, kp)
		}
	}
	return points
}
