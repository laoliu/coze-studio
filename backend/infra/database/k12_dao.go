package database
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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ==================== 数据库实体定义 ====================

// JSONStringArray 字符串数组的JSON类型
type JSONStringArray []string

// Scan 实现 sql.Scanner 接口
func (j *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONStringArray value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONStringArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

// JSONMap 通用JSON映射类型
type JSONMap map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONMap value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return json.Marshal(j)
}

// KnowledgePointEntity 知识点数据库实体
type KnowledgePointEntity struct {
	ID            string          `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	Name          string          `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Subject       string          `gorm:"column:subject;type:varchar(32);not null;index:idx_subject_grade" json:"subject"`
	Grade         string          `gorm:"column:grade;type:varchar(32);not null;index:idx_subject_grade" json:"grade"`
	Description   string          `gorm:"column:description;type:text" json:"description"`
	Keywords      JSONStringArray `gorm:"column:keywords;type:json" json:"keywords"`
	Prerequisites JSONStringArray `gorm:"column:prerequisites;type:json" json:"prerequisites"`
	NextPoints    JSONStringArray `gorm:"column:next_points;type:json" json:"next_points"`
	Difficulty    string          `gorm:"column:difficulty;type:varchar(32)" json:"difficulty"`
	Metadata      JSONMap         `gorm:"column:metadata;type:json" json:"metadata"`
	Version       int             `gorm:"column:version;default:1" json:"version"`
	CreatedAt     time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	ExpiresAt     *time.Time      `gorm:"column:expires_at;index:idx_expires_at" json:"expires_at"`
}

// TableName 指定表名
func (KnowledgePointEntity) TableName() string {
	return "knowledge_points"
}

// LearningResourceEntity 学习资源数据库实体
type LearningResourceEntity struct {
	ID          string     `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	KnowledgeID string     `gorm:"column:knowledge_id;type:varchar(64);not null;index:idx_knowledge_id" json:"knowledge_id"`
	Type        string     `gorm:"column:type;type:varchar(32);not null;index:idx_type" json:"type"`
	Title       string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	URL         string     `gorm:"column:url;type:text" json:"url"`
	Difficulty  string     `gorm:"column:difficulty;type:varchar(32)" json:"difficulty"`
	Duration    int        `gorm:"column:duration;type:int" json:"duration"`
	Source      string     `gorm:"column:source;type:varchar(255)" json:"source"`
	Metadata    JSONMap    `gorm:"column:metadata;type:json" json:"metadata"`
	Version     int        `gorm:"column:version;default:1" json:"version"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	ExpiresAt   *time.Time `gorm:"column:expires_at;index:idx_expires_at" json:"expires_at"`
}

// TableName 指定表名
func (LearningResourceEntity) TableName() string {
	return "learning_resources"
}

// ExerciseEntity 习题数据库实体
type ExerciseEntity struct {
	ID          string          `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	KnowledgeID string          `gorm:"column:knowledge_id;type:varchar(64);not null;index:idx_knowledge_id" json:"knowledge_id"`
	Question    string          `gorm:"column:question;type:text;not null" json:"question"`
	Options     JSONStringArray `gorm:"column:options;type:json" json:"options"`
	Answer      string          `gorm:"column:answer;type:text;not null" json:"answer"`
	Solution    string          `gorm:"column:solution;type:text" json:"solution"`
	Difficulty  string          `gorm:"column:difficulty;type:varchar(32);index:idx_difficulty" json:"difficulty"`
	Score       int             `gorm:"column:score;type:int" json:"score"`
	Metadata    JSONMap         `gorm:"column:metadata;type:json" json:"metadata"`
	Version     int             `gorm:"column:version;default:1" json:"version"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	ExpiresAt   *time.Time      `gorm:"column:expires_at;index:idx_expires_at" json:"expires_at"`
}

// TableName 指定表名
func (ExerciseEntity) TableName() string {
	return "exercises"
}

// GenerationLogEntity 生成日志数据库实体
type GenerationLogEntity struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RequestType    string    `gorm:"column:request_type;type:varchar(64);index:idx_request_type" json:"request_type"`
	Subject        string    `gorm:"column:subject;type:varchar(32)" json:"subject"`
	Grade          string    `gorm:"column:grade;type:varchar(32)" json:"grade"`
	KnowledgeID    string    `gorm:"column:knowledge_id;type:varchar(64);index:idx_knowledge_id" json:"knowledge_id"`
	PromptTemplate string    `gorm:"column:prompt_template;type:text" json:"prompt_template"`
	LLMResponse    string    `gorm:"column:llm_response;type:text" json:"llm_response"`
	Success        bool      `gorm:"column:success;index:idx_success" json:"success"`
	ErrorMessage   string    `gorm:"column:error_message;type:text" json:"error_message"`
	DurationMs     int       `gorm:"column:duration_ms;type:int" json:"duration_ms"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime;index:idx_created_at" json:"created_at"`
}

// TableName 指定表名
func (GenerationLogEntity) TableName() string {
	return "generation_logs"
}

// CacheConfigEntity 缓存配置数据库实体
type CacheConfigEntity struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ItemType    string    `gorm:"column:item_type;type:varchar(32);not null;uniqueIndex:uk_item_type" json:"item_type"`
	TTLSeconds  int       `gorm:"column:ttl_seconds;type:int;not null" json:"ttl_seconds"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (CacheConfigEntity) TableName() string {
	return "cache_config"
}

// ==================== DAO 接口定义 ====================

// K12DAO K12知识图谱数据访问接口
type K12DAO interface {
	// ========== 知识点操作 ==========
	// CreateKnowledgePoint 创建知识点
	CreateKnowledgePoint(ctx context.Context, kp *KnowledgePointEntity) error
	// GetKnowledgePoint 根据ID获取知识点
	GetKnowledgePoint(ctx context.Context, id string) (*KnowledgePointEntity, error)
	// UpdateKnowledgePoint 更新知识点
	UpdateKnowledgePoint(ctx context.Context, kp *KnowledgePointEntity) error
	// DeleteKnowledgePoint 删除知识点
	DeleteKnowledgePoint(ctx context.Context, id string) error
	// ListKnowledgePoints 列出知识点（支持过滤）
	ListKnowledgePoints(ctx context.Context, filter KnowledgePointFilter) ([]*KnowledgePointEntity, error)
	// GetExpiredKnowledgePoints 获取已过期的知识点
	GetExpiredKnowledgePoints(ctx context.Context, limit int) ([]*KnowledgePointEntity, error)
	// BatchCreateKnowledgePoints 批量创建知识点
	BatchCreateKnowledgePoints(ctx context.Context, kps []*KnowledgePointEntity) error

	// ========== 学习资源操作 ==========
	// CreateLearningResource 创建学习资源
	CreateLearningResource(ctx context.Context, resource *LearningResourceEntity) error
	// GetLearningResource 根据ID获取学习资源
	GetLearningResource(ctx context.Context, id string) (*LearningResourceEntity, error)
	// UpdateLearningResource 更新学习资源
	UpdateLearningResource(ctx context.Context, resource *LearningResourceEntity) error
	// DeleteLearningResource 删除学习资源
	DeleteLearningResource(ctx context.Context, id string) error
	// ListLearningResources 列出学习资源（支持过滤）
	ListLearningResources(ctx context.Context, filter ResourceFilter) ([]*LearningResourceEntity, error)
	// GetResourcesByKnowledgeID 根据知识点ID获取资源
	GetResourcesByKnowledgeID(ctx context.Context, knowledgeID string) ([]*LearningResourceEntity, error)
	// BatchCreateLearningResources 批量创建学习资源
	BatchCreateLearningResources(ctx context.Context, resources []*LearningResourceEntity) error

	// ========== 习题操作 ==========
	// CreateExercise 创建习题
	CreateExercise(ctx context.Context, exercise *ExerciseEntity) error
	// GetExercise 根据ID获取习题
	GetExercise(ctx context.Context, id string) (*ExerciseEntity, error)
	// UpdateExercise 更新习题
	UpdateExercise(ctx context.Context, exercise *ExerciseEntity) error
	// DeleteExercise 删除习题
	DeleteExercise(ctx context.Context, id string) error
	// ListExercises 列出习题（支持过滤）
	ListExercises(ctx context.Context, filter ExerciseFilter) ([]*ExerciseEntity, error)
	// GetExercisesByKnowledgeID 根据知识点ID获取习题
	GetExercisesByKnowledgeID(ctx context.Context, knowledgeID string) ([]*ExerciseEntity, error)
	// BatchCreateExercises 批量创建习题
	BatchCreateExercises(ctx context.Context, exercises []*ExerciseEntity) error

	// ========== 生成日志操作 ==========
	// CreateGenerationLog 创建生成日志
	CreateGenerationLog(ctx context.Context, log *GenerationLogEntity) error
	// GetGenerationLog 根据ID获取生成日志
	GetGenerationLog(ctx context.Context, id int64) (*GenerationLogEntity, error)
	// ListGenerationLogs 列出生成日志（支持过滤）
	ListGenerationLogs(ctx context.Context, filter GenerationLogFilter) ([]*GenerationLogEntity, error)

	// ========== 缓存配置操作 ==========
	// GetCacheConfig 根据类型获取缓存配置
	GetCacheConfig(ctx context.Context, itemType string) (*CacheConfigEntity, error)
	// ListCacheConfigs 列出所有缓存配置
	ListCacheConfigs(ctx context.Context) ([]*CacheConfigEntity, error)
	// UpdateCacheConfig 更新缓存配置
	UpdateCacheConfig(ctx context.Context, config *CacheConfigEntity) error

	// ========== 统计操作 ==========
	// CountKnowledgePoints 统计知识点数量
	CountKnowledgePoints(ctx context.Context, filter KnowledgePointFilter) (int64, error)
	// CountResources 统计资源数量
	CountResources(ctx context.Context, filter ResourceFilter) (int64, error)
	// CountExercises 统计习题数量
	CountExercises(ctx context.Context, filter ExerciseFilter) (int64, error)
}

// ==================== 过滤器定义 ====================

// KnowledgePointFilter 知识点过滤器
type KnowledgePointFilter struct {
	Subject    string   // 学科过滤
	Grade      string   // 年级过滤
	Difficulty string   // 难度过滤
	IDs        []string // ID列表过滤
	Expired    *bool    // 是否已过期（nil表示不过滤）
	Limit      int      // 限制数量
	Offset     int      // 偏移量
}

// ResourceFilter 资源过滤器
type ResourceFilter struct {
	KnowledgeID string // 知识点ID过滤
	Type        string // 资源类型过滤
	Difficulty  string // 难度过滤
	Limit       int    // 限制数量
	Offset      int    // 偏移量
}

// ExerciseFilter 习题过滤器
type ExerciseFilter struct {
	KnowledgeID string // 知识点ID过滤
	Difficulty  string // 难度过滤
	Limit       int    // 限制数量
	Offset      int    // 偏移量
}

// GenerationLogFilter 生成日志过滤器
type GenerationLogFilter struct {
	RequestType string     // 请求类型过滤
	Subject     string     // 学科过滤
	Grade       string     // 年级过滤
	KnowledgeID string     // 知识点ID过滤
	Success     *bool      // 是否成功（nil表示不过滤）
	StartTime   *time.Time // 开始时间
	EndTime     *time.Time // 结束时间
	Limit       int        // 限制数量
	Offset      int        // 偏移量
}

// ==================== GORM 实现 ====================

// k12DAOImpl K12 DAO的GORM实现
type k12DAOImpl struct {
	db *gorm.DB
}

// NewK12DAO 创建K12 DAO实例
func NewK12DAO(db *gorm.DB) K12DAO {
	return &k12DAOImpl{db: db}
}

// ========== 知识点操作实现 ==========

func (d *k12DAOImpl) CreateKnowledgePoint(ctx context.Context, kp *KnowledgePointEntity) error {
	return d.db.WithContext(ctx).Create(kp).Error
}

func (d *k12DAOImpl) GetKnowledgePoint(ctx context.Context, id string) (*KnowledgePointEntity, error) {
	var kp KnowledgePointEntity
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&kp).Error
	if err != nil {
		return nil, err
	}
	return &kp, nil
}

func (d *k12DAOImpl) UpdateKnowledgePoint(ctx context.Context, kp *KnowledgePointEntity) error {
	return d.db.WithContext(ctx).Save(kp).Error
}

func (d *k12DAOImpl) DeleteKnowledgePoint(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&KnowledgePointEntity{}).Error
}

func (d *k12DAOImpl) ListKnowledgePoints(ctx context.Context, filter KnowledgePointFilter) ([]*KnowledgePointEntity, error) {
	query := d.db.WithContext(ctx).Model(&KnowledgePointEntity{})

	// 应用过滤条件
	if filter.Subject != "" {
		query = query.Where("subject = ?", filter.Subject)
	}
	if filter.Grade != "" {
		query = query.Where("grade = ?", filter.Grade)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}
	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}
	if filter.Expired != nil {
		now := time.Now()
		if *filter.Expired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", now)
		} else {
			query = query.Where("expires_at IS NULL OR expires_at >= ?", now)
		}
	}

	// 应用分页
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var kps []*KnowledgePointEntity
	err := query.Find(&kps).Error
	return kps, err
}

func (d *k12DAOImpl) GetExpiredKnowledgePoints(ctx context.Context, limit int) ([]*KnowledgePointEntity, error) {
	now := time.Now()
	query := d.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Order("expires_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var kps []*KnowledgePointEntity
	err := query.Find(&kps).Error
	return kps, err
}

func (d *k12DAOImpl) BatchCreateKnowledgePoints(ctx context.Context, kps []*KnowledgePointEntity) error {
	if len(kps) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).CreateInBatches(kps, 100).Error
}

// ========== 学习资源操作实现 ==========

func (d *k12DAOImpl) CreateLearningResource(ctx context.Context, resource *LearningResourceEntity) error {
	return d.db.WithContext(ctx).Create(resource).Error
}

func (d *k12DAOImpl) GetLearningResource(ctx context.Context, id string) (*LearningResourceEntity, error) {
	var resource LearningResourceEntity
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (d *k12DAOImpl) UpdateLearningResource(ctx context.Context, resource *LearningResourceEntity) error {
	return d.db.WithContext(ctx).Save(resource).Error
}

func (d *k12DAOImpl) DeleteLearningResource(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&LearningResourceEntity{}).Error
}

func (d *k12DAOImpl) ListLearningResources(ctx context.Context, filter ResourceFilter) ([]*LearningResourceEntity, error) {
	query := d.db.WithContext(ctx).Model(&LearningResourceEntity{})

	// 应用过滤条件
	if filter.KnowledgeID != "" {
		query = query.Where("knowledge_id = ?", filter.KnowledgeID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	// 应用分页
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var resources []*LearningResourceEntity
	err := query.Find(&resources).Error
	return resources, err
}

func (d *k12DAOImpl) GetResourcesByKnowledgeID(ctx context.Context, knowledgeID string) ([]*LearningResourceEntity, error) {
	var resources []*LearningResourceEntity
	err := d.db.WithContext(ctx).
		Where("knowledge_id = ?", knowledgeID).
		Find(&resources).Error
	return resources, err
}

func (d *k12DAOImpl) BatchCreateLearningResources(ctx context.Context, resources []*LearningResourceEntity) error {
	if len(resources) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).CreateInBatches(resources, 100).Error
}

// ========== 习题操作实现 ==========

func (d *k12DAOImpl) CreateExercise(ctx context.Context, exercise *ExerciseEntity) error {
	return d.db.WithContext(ctx).Create(exercise).Error
}

func (d *k12DAOImpl) GetExercise(ctx context.Context, id string) (*ExerciseEntity, error) {
	var exercise ExerciseEntity
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&exercise).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (d *k12DAOImpl) UpdateExercise(ctx context.Context, exercise *ExerciseEntity) error {
	return d.db.WithContext(ctx).Save(exercise).Error
}

func (d *k12DAOImpl) DeleteExercise(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&ExerciseEntity{}).Error
}

func (d *k12DAOImpl) ListExercises(ctx context.Context, filter ExerciseFilter) ([]*ExerciseEntity, error) {
	query := d.db.WithContext(ctx).Model(&ExerciseEntity{})

	// 应用过滤条件
	if filter.KnowledgeID != "" {
		query = query.Where("knowledge_id = ?", filter.KnowledgeID)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	// 应用分页
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var exercises []*ExerciseEntity
	err := query.Find(&exercises).Error
	return exercises, err
}

func (d *k12DAOImpl) GetExercisesByKnowledgeID(ctx context.Context, knowledgeID string) ([]*ExerciseEntity, error) {
	var exercises []*ExerciseEntity
	err := d.db.WithContext(ctx).
		Where("knowledge_id = ?", knowledgeID).
		Find(&exercises).Error
	return exercises, err
}

func (d *k12DAOImpl) BatchCreateExercises(ctx context.Context, exercises []*ExerciseEntity) error {
	if len(exercises) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).CreateInBatches(exercises, 100).Error
}

// ========== 生成日志操作实现 ==========

func (d *k12DAOImpl) CreateGenerationLog(ctx context.Context, log *GenerationLogEntity) error {
	return d.db.WithContext(ctx).Create(log).Error
}

func (d *k12DAOImpl) GetGenerationLog(ctx context.Context, id int64) (*GenerationLogEntity, error) {
	var log GenerationLogEntity
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (d *k12DAOImpl) ListGenerationLogs(ctx context.Context, filter GenerationLogFilter) ([]*GenerationLogEntity, error) {
	query := d.db.WithContext(ctx).Model(&GenerationLogEntity{})

	// 应用过滤条件
	if filter.RequestType != "" {
		query = query.Where("request_type = ?", filter.RequestType)
	}
	if filter.Subject != "" {
		query = query.Where("subject = ?", filter.Subject)
	}
	if filter.Grade != "" {
		query = query.Where("grade = ?", filter.Grade)
	}
	if filter.KnowledgeID != "" {
		query = query.Where("knowledge_id = ?", filter.KnowledgeID)
	}
	if filter.Success != nil {
		query = query.Where("success = ?", *filter.Success)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// 按创建时间倒序
	query = query.Order("created_at DESC")

	// 应用分页
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var logs []*GenerationLogEntity
	err := query.Find(&logs).Error
	return logs, err
}

// ========== 缓存配置操作实现 ==========

func (d *k12DAOImpl) GetCacheConfig(ctx context.Context, itemType string) (*CacheConfigEntity, error) {
	var config CacheConfigEntity
	err := d.db.WithContext(ctx).Where("item_type = ?", itemType).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (d *k12DAOImpl) ListCacheConfigs(ctx context.Context) ([]*CacheConfigEntity, error) {
	var configs []*CacheConfigEntity
	err := d.db.WithContext(ctx).Find(&configs).Error
	return configs, err
}

func (d *k12DAOImpl) UpdateCacheConfig(ctx context.Context, config *CacheConfigEntity) error {
	return d.db.WithContext(ctx).Save(config).Error
}

// ========== 统计操作实现 ==========

func (d *k12DAOImpl) CountKnowledgePoints(ctx context.Context, filter KnowledgePointFilter) (int64, error) {
	query := d.db.WithContext(ctx).Model(&KnowledgePointEntity{})

	// 应用过滤条件
	if filter.Subject != "" {
		query = query.Where("subject = ?", filter.Subject)
	}
	if filter.Grade != "" {
		query = query.Where("grade = ?", filter.Grade)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}
	if filter.Expired != nil {
		now := time.Now()
		if *filter.Expired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", now)
		} else {
			query = query.Where("expires_at IS NULL OR expires_at >= ?", now)
		}
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (d *k12DAOImpl) CountResources(ctx context.Context, filter ResourceFilter) (int64, error) {
	query := d.db.WithContext(ctx).Model(&LearningResourceEntity{})

	// 应用过滤条件
	if filter.KnowledgeID != "" {
		query = query.Where("knowledge_id = ?", filter.KnowledgeID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (d *k12DAOImpl) CountExercises(ctx context.Context, filter ExerciseFilter) (int64, error) {
	query := d.db.WithContext(ctx).Model(&ExerciseEntity{})

	// 应用过滤条件
	if filter.KnowledgeID != "" {
		query = query.Where("knowledge_id = ?", filter.KnowledgeID)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}
