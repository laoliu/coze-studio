// Copyright 2025 Coze Studio. All rights reserved.

package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ContentRepository 内容仓储接口
type ContentRepository interface {
	// SearchByKeywords 根据关键词搜索内容
	SearchByKeywords(ctx context.Context, keywords []string, limit int) ([]*Content, error)

	// SearchByTopic 根据主题搜索内容
	SearchByTopic(ctx context.Context, topic, domain, grade string, limit int) ([]*Content, error)

	// GetByID 根据 ID 获取内容
	GetByID(ctx context.Context, id string) (*Content, error)

	// GetByIDs 批量获取内容
	GetByIDs(ctx context.Context, ids []string) ([]*Content, error)

	// Create 创建内容
	Create(ctx context.Context, content *Content) error

	// Update 更新内容
	Update(ctx context.Context, content *Content) error
}

// Content 内容实体
type Content struct {
	ID          string  `gorm:"primaryKey;size:100"`
	Type        string  `gorm:"size:50;index"` // text, image, video, audio
	Title       string  `gorm:"size:500;index"`
	Description string  `gorm:"type:text"`
	Source      string  `gorm:"size:200"`
	URL         string  `gorm:"size:1000"`
	Domain      string  `gorm:"size:100;index"`     // 学科领域
	Grade       string  `gorm:"size:50;index"`      // 年级
	Topic       string  `gorm:"size:200;index"`     // 主题
	Keywords    string  `gorm:"type:jsonb"`         // 关键词数组
	Metadata    string  `gorm:"type:jsonb"`         // 元数据
	Relevance   float64 `gorm:"type:decimal(10,2)"` // 相关性评分
	CreatedAt   int64   `gorm:"autoCreateTime"`
	UpdatedAt   int64   `gorm:"autoUpdateTime"`
}

// TableName 表名
func (Content) TableName() string {
	return "contents"
}

// GormContentRepository Gorm 实现的内容仓储
type GormContentRepository struct {
	db *gorm.DB
}

// NewGormContentRepository 创建 Gorm 内容仓储
func NewGormContentRepository(db *gorm.DB) *GormContentRepository {
	return &GormContentRepository{db: db}
}

// SearchByKeywords 根据关键词搜索
func (r *GormContentRepository) SearchByKeywords(ctx context.Context, keywords []string, limit int) ([]*Content, error) {
	var contents []*Content

	query := r.db.WithContext(ctx)

	// 构建 OR 查询
	for i, keyword := range keywords {
		if i == 0 {
			query = query.Where("title LIKE ? OR description LIKE ? OR topic LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		} else {
			query = query.Or("title LIKE ? OR description LIKE ? OR topic LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if err := query.Order("relevance DESC").Limit(limit).Find(&contents).Error; err != nil {
		return nil, fmt.Errorf("failed to search contents: %w", err)
	}

	return contents, nil
}

// SearchByTopic 根据主题搜索
func (r *GormContentRepository) SearchByTopic(ctx context.Context, topic, domain, grade string, limit int) ([]*Content, error) {
	var contents []*Content

	query := r.db.WithContext(ctx)

	if topic != "" {
		query = query.Where("topic LIKE ?", "%"+topic+"%")
	}
	if domain != "" {
		query = query.Where("domain = ?", domain)
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}

	if err := query.Order("relevance DESC").Limit(limit).Find(&contents).Error; err != nil {
		return nil, fmt.Errorf("failed to search contents by topic: %w", err)
	}

	return contents, nil
}

// GetByID 获取单个内容
func (r *GormContentRepository) GetByID(ctx context.Context, id string) (*Content, error) {
	var content Content
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&content).Error; err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	return &content, nil
}

// GetByIDs 批量获取
func (r *GormContentRepository) GetByIDs(ctx context.Context, ids []string) ([]*Content, error) {
	var contents []*Content
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&contents).Error; err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}
	return contents, nil
}

// Create 创建内容
func (r *GormContentRepository) Create(ctx context.Context, content *Content) error {
	if err := r.db.WithContext(ctx).Create(content).Error; err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}
	return nil
}

// Update 更新内容
func (r *GormContentRepository) Update(ctx context.Context, content *Content) error {
	if err := r.db.WithContext(ctx).Save(content).Error; err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}
	return nil
}

// MockContentRepository Mock 内容仓储（用于测试）
type MockContentRepository struct {
	contents map[string]*Content
}

// NewMockContentRepository 创建 Mock 仓储
func NewMockContentRepository() *MockContentRepository {
	return &MockContentRepository{
		contents: make(map[string]*Content),
	}
}

// AddContent 添加测试内容
func (m *MockContentRepository) AddContent(content *Content) {
	m.contents[content.ID] = content
}

// SearchByKeywords Mock 实现
func (m *MockContentRepository) SearchByKeywords(ctx context.Context, keywords []string, limit int) ([]*Content, error) {
	results := make([]*Content, 0)
	for _, content := range m.contents {
		for _, keyword := range keywords {
			if contains(content.Title, keyword) || contains(content.Description, keyword) {
				results = append(results, content)
				break
			}
		}
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

// SearchByTopic Mock 实现
func (m *MockContentRepository) SearchByTopic(ctx context.Context, topic, domain, grade string, limit int) ([]*Content, error) {
	results := make([]*Content, 0)
	for _, content := range m.contents {
		match := true
		if topic != "" && !contains(content.Topic, topic) {
			match = false
		}
		if domain != "" && content.Domain != domain {
			match = false
		}
		if grade != "" && content.Grade != grade {
			match = false
		}
		if match {
			results = append(results, content)
		}
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

// GetByID Mock 实现
func (m *MockContentRepository) GetByID(ctx context.Context, id string) (*Content, error) {
	content, ok := m.contents[id]
	if !ok {
		return nil, fmt.Errorf("content not found: %s", id)
	}
	return content, nil
}

// GetByIDs Mock 实现
func (m *MockContentRepository) GetByIDs(ctx context.Context, ids []string) ([]*Content, error) {
	results := make([]*Content, 0, len(ids))
	for _, id := range ids {
		if content, ok := m.contents[id]; ok {
			results = append(results, content)
		}
	}
	return results, nil
}

// Create Mock 实现
func (m *MockContentRepository) Create(ctx context.Context, content *Content) error {
	m.contents[content.ID] = content
	return nil
}

// Update Mock 实现
func (m *MockContentRepository) Update(ctx context.Context, content *Content) error {
	m.contents[content.ID] = content
	return nil
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && (s == substr || len(s) > len(substr))
}
