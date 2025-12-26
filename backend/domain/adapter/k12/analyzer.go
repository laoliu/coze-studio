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
	"fmt"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// Analyzer 输入分析器
type Analyzer struct {
}

// NewAnalyzer 创建分析器
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeInput 分析用户输入
func (a *Analyzer) AnalyzeInput(userInput string) (*models.AnalysisResult, error) {
	input := strings.ToLower(userInput)

	// 分析学科
	subject := a.detectSubject(input)
	if subject == "" {
		return nil, fmt.Errorf("无法识别学科")
	}

	// 分析年级
	grade := a.detectGrade(input)
	if grade == "" {
		return nil, fmt.Errorf("无法识别年级")
	}

	// 分析知识点
	knowledgeID, knowledgeName, confidence := a.detectKnowledge(input, subject)
	if knowledgeID == "" {
		return nil, fmt.Errorf("无法识别知识点")
	}

	// 获取前置知识
	prerequisites := a.getPrerequisiteIDs(knowledgeID)

	result := &models.AnalysisResult{
		Subject:       subject,
		Grade:         grade,
		KnowledgeID:   knowledgeID,
		KnowledgeName: knowledgeName,
		Prerequisites: prerequisites,
		Confidence:    confidence,
	}

	return result, nil
}

// detectSubject 检测学科
func (a *Analyzer) detectSubject(input string) models.Subject {
	subjectKeywords := map[models.Subject][]string{
		models.SubjectMath:      {"数学", "math"},
		models.SubjectPhysics:   {"物理", "physics"},
		models.SubjectChemistry: {"化学", "chemistry"},
		models.SubjectBiology:   {"生物", "biology"},
	}

	for subject, keywords := range subjectKeywords {
		for _, keyword := range keywords {
			if strings.Contains(input, keyword) {
				return subject
			}
		}
	}

	return ""
}

// detectGrade 检测年级
func (a *Analyzer) detectGrade(input string) models.GradeLevel {
	gradeKeywords := map[models.GradeLevel][]string{
		models.Grade7:  {"初一", "七年级", "grade 7", "7年级"},
		models.Grade8:  {"初二", "八年级", "grade 8", "8年级"},
		models.Grade9:  {"初三", "九年级", "grade 9", "9年级"},
		models.Grade10: {"高一", "十年级", "grade 10", "10年级"},
		models.Grade11: {"高二", "十一年级", "grade 11", "11年级"},
		models.Grade12: {"高三", "十二年级", "grade 12", "12年级"},
	}

	// 首先尝试精确匹配
	for grade, keywords := range gradeKeywords {
		for _, keyword := range keywords {
			if strings.Contains(input, keyword) {
				return grade
			}
		}
	}

	// 如果没有明确年级，根据学段推断
	if strings.Contains(input, "初中") {
		return models.Grade8 // 默认初二
	}
	if strings.Contains(input, "高中") {
		return models.Grade10 // 默认高一
	}

	return ""
}

// detectKnowledge 检测知识点
func (a *Analyzer) detectKnowledge(input string, subject models.Subject) (string, string, float64) {
	// 数学知识点关键词映射
	mathKnowledgeKeywords := map[string][]string{
		"quadratic_functions":              {"二次函数", "抛物线", "quadratic function"},
		"quadratic_equations":              {"一元二次方程", "quadratic equation"},
		"linear_functions":                 {"一次函数", "线性函数", "linear function"},
		"linear_equations_systems":         {"二元一次方程组", "方程组", "system of equations"},
		"equations":                        {"一元一次方程", "一次方程", "linear equation"},
		"real_numbers":                     {"实数", "real number"},
		"quadratic_functions_advanced":     {"二次函数综合", "二次函数应用"},
		"inverse_proportional_functions":   {"反比例函数", "inverse function"},
	}

	// 物理知识点关键词映射
	physicsKnowledgeKeywords := map[string][]string{
		"ohms_law":             {"欧姆定律", "ohm's law", "电阻定律"},
		"electrical_power":     {"电功率", "electrical power", "焦耳定律", "电功"},
		"circuit_analysis":     {"电路分析", "电路计算", "动态电路"},
		"motion_and_forces":    {"牛顿第一定律", "惯性", "二力平衡", "运动和力"},
		"forces":               {"力的概念", "弹力", "重力", "摩擦力"},
		"speed":                {"速度", "匀速运动", "平均速度"},
		"pressure":             {"压强", "液体压强", "大气压"},
		"buoyancy":             {"浮力", "阿基米德原理"},
		"work_and_energy":      {"功", "机械能", "动能", "势能"},
		"simple_machines":      {"杠杆", "滑轮", "机械效率"},
		"electricity_basics":   {"电流", "电路", "串联", "并联"},
		"voltage_resistance":   {"电压", "电阻"},
		"light_reflection":     {"光的反射", "平面镜"},
		"light_refraction":     {"光的折射", "透镜", "凸透镜"},
	}

	var keywordMap map[string][]string
	if subject == models.SubjectMath {
		keywordMap = mathKnowledgeKeywords
	} else if subject == models.SubjectPhysics {
		keywordMap = physicsKnowledgeKeywords
	}

	if keywordMap != nil {
		for knowledgeID, keywords := range keywordMap {
			for _, keyword := range keywords {
				if strings.Contains(input, keyword) {
					knowledge := GetKnowledgePoint(knowledgeID)
					if knowledge != nil {
						// 根据关键词匹配程度计算置信度
						confidence := 0.8
						if strings.Contains(input, keywords[0]) {
							confidence = 0.95
						}
						return knowledgeID, knowledge.Name, confidence
					}
				}
			}
		}
	}

	return "", "", 0.0
}

// getPrerequisiteIDs 获取前置知识ID列表
func (a *Analyzer) getPrerequisiteIDs(knowledgeID string) []string {
	prerequisites := GetPrerequisites(knowledgeID)
	var ids []string
	for _, kp := range prerequisites {
		ids = append(ids, kp.ID)
	}
	return ids
}

// AnalyzeProgress 分析学习进度
// userInput: 用户输入
// masteredKnowledgeIDs: 已掌握的知识点ID列表
func (a *Analyzer) AnalyzeProgress(userInput string, masteredKnowledgeIDs []string) (*models.AnalysisResult, []string, error) {
	// 首先分析用户输入
	result, err := a.AnalyzeInput(userInput)
	if err != nil {
		return nil, nil, err
	}

	// 构建已掌握知识点映射
	masteredMap := make(map[string]bool)
	for _, id := range masteredKnowledgeIDs {
		masteredMap[id] = true
	}

	// 检查前置知识掌握情况
	var missingPrerequisites []string
	for _, preID := range result.Prerequisites {
		if !masteredMap[preID] {
			missingPrerequisites = append(missingPrerequisites, preID)
		}
	}

	return result, missingPrerequisites, nil
}
