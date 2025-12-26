package models
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

package models

// Subject 学科
type Subject string

const (
	SubjectMath    Subject = "math"    // 数学
	SubjectPhysics Subject = "physics" // 物理
	SubjectChemistry Subject = "chemistry" // 化学
	SubjectBiology Subject = "biology" // 生物
)

// GradeLevel 年级
type GradeLevel string

const (
	Grade7  GradeLevel = "grade_7"  // 初一
	Grade8  GradeLevel = "grade_8"  // 初二
	Grade9  GradeLevel = "grade_9"  // 初三
	Grade10 GradeLevel = "grade_10" // 高一
	Grade11 GradeLevel = "grade_11" // 高二
	Grade12 GradeLevel = "grade_12" // 高三
)

// DifficultyLevel 难度等级
type DifficultyLevel string

const (
	DifficultyBasic    DifficultyLevel = "basic"    // 基础
	DifficultyImprove  DifficultyLevel = "improve"  // 提高
	DifficultyAdvanced DifficultyLevel = "advanced" // 拓展
)

// KnowledgePoint 知识点
type KnowledgePoint struct {
	ID          string          `json:"id"`           // 知识点ID
	Name        string          `json:"name"`         // 知识点名称
	Subject     Subject         `json:"subject"`      // 所属学科
	Grade       GradeLevel      `json:"grade"`        // 年级
	Description string          `json:"description"`  // 描述
	Prerequisites []string      `json:"prerequisites"` // 前置知识点ID列表
	NextPoints    []string      `json:"next_points"`   // 后续知识点ID列表
	Difficulty    DifficultyLevel `json:"difficulty"`  // 难度
	Keywords      []string        `json:"keywords"`    // 关键词
}

// LearningResource 学习资源
type LearningResource struct {
	ID          string          `json:"id"`
	Type        ResourceType    `json:"type"`         // 资源类型
	Title       string          `json:"title"`        // 标题
	Description string          `json:"description"`  // 描述
	URL         string          `json:"url"`          // 链接
	KnowledgeID string          `json:"knowledge_id"` // 关联知识点
	Difficulty  DifficultyLevel `json:"difficulty"`   // 难度
	Duration    int             `json:"duration"`     // 时长（分钟）
	Source      string          `json:"source"`       // 来源
}

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeVideo    ResourceType = "video"    // 视频
	ResourceTypeExercise ResourceType = "exercise" // 习题
	ResourceTypeTextbook ResourceType = "textbook" // 教材
	ResourceTypeDocument ResourceType = "document" // 文档
)

// Exercise 习题
type Exercise struct {
	ID          string          `json:"id"`
	Question    string          `json:"question"`     // 题目
	Options     []string        `json:"options"`      // 选项（选择题）
	Answer      string          `json:"answer"`       // 答案
	Solution    string          `json:"solution"`     // 解析
	KnowledgeID string          `json:"knowledge_id"` // 关联知识点
	Difficulty  DifficultyLevel `json:"difficulty"`   // 难度
	Score       int             `json:"score"`        // 分值
}

// LearningPath 学习路径
type LearningPath struct {
	Steps []LearningStep `json:"steps"` // 学习步骤
}

// LearningStep 学习步骤
type LearningStep struct {
	Order       int              `json:"order"`        // 顺序
	Knowledge   KnowledgePoint   `json:"knowledge"`    // 知识点
	Objectives  []string         `json:"objectives"`   // 学习目标
	Resources   []LearningResource `json:"resources"`  // 推荐资源
	Exercises   []Exercise       `json:"exercises"`    // 配套习题
	EstimatedTime int            `json:"estimated_time"` // 预计学习时间（分钟）
}

// LearningPlan 学习计划
type LearningPlan struct {
	StudentLevel  string         `json:"student_level"`  // 学生水平
	TargetKnowledge KnowledgePoint `json:"target_knowledge"` // 目标知识点
	Prerequisites []KnowledgePoint `json:"prerequisites"` // 前置知识
	MainPath      LearningPath   `json:"main_path"`      // 主学习路径
	Review        []KnowledgePoint `json:"review"`       // 复习内容
	Extension     []KnowledgePoint `json:"extension"`    // 拓展内容
	TotalTime     int            `json:"total_time"`     // 总学习时间（分钟）
}

// AnalysisRequest 分析请求
type AnalysisRequest struct {
	UserInput string `json:"user_input"` // 用户输入
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	Subject       Subject        `json:"subject"`        // 学科
	Grade         GradeLevel     `json:"grade"`          // 年级
	KnowledgeID   string         `json:"knowledge_id"`   // 知识点ID
	KnowledgeName string         `json:"knowledge_name"` // 知识点名称
	Prerequisites []string       `json:"prerequisites"`  // 前置知识
	Confidence    float64        `json:"confidence"`     // 置信度
}
