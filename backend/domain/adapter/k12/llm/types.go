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

package llm

import (
	"context"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
)

// LLMService LLM服务接口
type LLMService interface {
	// GenerateKnowledgeGraph 生成完整知识图谱
	GenerateKnowledgeGraph(ctx context.Context, req *KnowledgeGraphRequest) (*KnowledgeGraph, error)

	// GenerateKnowledgePoint 生成单个知识点
	GenerateKnowledgePoint(ctx context.Context, req *KnowledgePointRequest) (*models.KnowledgePoint, error)

	// GenerateResources 生成学习资源
	GenerateResources(ctx context.Context, req *ResourceRequest) ([]*models.Resource, error)

	// GenerateExercises 生成习题
	GenerateExercises(ctx context.Context, req *ExerciseRequest) ([]*models.Exercise, error)
}

// KnowledgeGraphRequest 知识图谱生成请求
type KnowledgeGraphRequest struct {
	Subject           models.Subject
	Grade             models.GradeLevel
	AdditionalContext string
}

// KnowledgeGraph 知识图谱
type KnowledgeGraph struct {
	Subject         models.Subject
	Grade           models.GradeLevel
	KnowledgePoints []*models.KnowledgePoint
	GeneratedAt     time.Time
}

// KnowledgePointRequest 知识点生成请求
type KnowledgePointRequest struct {
	Subject     models.Subject
	Grade       models.GradeLevel
	KnowledgeID string
	Name        string
}

// ResourceRequest 资源生成请求
type ResourceRequest struct {
	KnowledgeID   string
	KnowledgeName string
	Subject       models.Subject
	Grade         models.GradeLevel
	Description   string
	Count         int // 生成资源数量，建议3-5个
}

// ExerciseRequest 习题生成请求
type ExerciseRequest struct {
	KnowledgeID   string
	KnowledgeName string
	Subject       models.Subject
	Grade         models.GradeLevel
	Difficulty    models.Difficulty
	Count         int // 生成习题数量
}

// GenerationLog 生成日志
type GenerationLog struct {
	ID             int64
	RequestType    string // knowledge_graph, knowledge_point, resource, exercise
	Subject        models.Subject
	Grade          models.GradeLevel
	KnowledgeID    string
	PromptTemplate string
	LLMResponse    string
	Success        bool
	ErrorMessage   string
	DurationMS     int64
	CreatedAt      time.Time
}
