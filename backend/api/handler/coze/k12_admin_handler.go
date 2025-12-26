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

package coze

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/adapter/k12/cache"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/k12/llm"
	"github.com/coze-dev/coze-studio/backend/domain/adapter/models"
	"github.com/coze-dev/coze-studio/backend/infra/database"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// K12AdminHandler K12知识图谱管理API处理器
type K12AdminHandler struct {
	db           *gorm.DB
	dao          database.K12DAO
	llmService   llm.LLMService
	cacheManager *cache.CacheManager
	logger       hlog.FullLogger
}

// NewK12AdminHandler 创建K12管理API处理器
func NewK12AdminHandler(db *gorm.DB) (*K12AdminHandler, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	// 创建DAO
	dao := database.NewK12DAO(db)

	// 创建LLM服务
	llmService := llm.NewCozeLLMService("", "", "", 30*time.Second, hlog.DefaultLogger())

	// 创建缓存配置
	cacheConfig := &cache.CacheConfig{
		EnableMemory: true,
		EnableRedis:  false, // 默认不启用Redis
		EnableMySQL:  true,
		DefaultTTL:   24 * time.Hour,
	}

	// 创建缓存管理器
	cacheManager, err := cache.NewCacheManager(cacheConfig, db, hlog.DefaultLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to create cache manager: %w", err)
	}

	return &K12AdminHandler{
		db:           db,
		dao:          dao,
		llmService:   llmService,
		cacheManager: cacheManager,
		logger:       hlog.DefaultLogger(),
	}, nil
}

// ==================== API 1: 生成知识图谱 ====================

// GenerateKnowledgeGraphRequest 生成知识图谱请求
type GenerateKnowledgeGraphRequest struct {
	Subject string `json:"subject" binding:"required"` // 学科: math, physics, chemistry, biology
	Grade   string `json:"grade" binding:"required"`   // 年级: grade_1 ~ grade_12
	Force   bool   `json:"force"`                      // 强制重新生成（忽略缓存）
}

// GenerateKnowledgeGraph 生成知识图谱
// POST /api/adapter/k12/admin/knowledge-graph/generate
func (h *K12AdminHandler) GenerateKnowledgeGraph(ctx context.Context, c *app.RequestContext) {
	var req GenerateKnowledgeGraphRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}

	// 1. 如果force=true，先删除缓存
	if req.Force {
		cacheKey := h.cacheManager.GenerateKnowledgeGraphKey(req.Subject, req.Grade)
		h.cacheManager.Delete(ctx, cacheKey)
		h.logger.Infof("[K12Admin] Force regenerate, cache cleared: %s", cacheKey)
	}

	// 2. 尝试从缓存获取
	cacheKey := h.cacheManager.GenerateKnowledgeGraphKey(req.Subject, req.Grade)
	if cached, ok := h.cacheManager.Get(ctx, cacheKey); ok && !req.Force {
		h.logger.Infof("[K12Admin] Cache hit for knowledge graph: %s/%s", req.Subject, req.Grade)
		c.JSON(consts.StatusOK, utils.H{
			"code": 0,
			"msg":  "success",
			"data": utils.H{
				"cached":  true,
				"content": cached,
			},
		})
		return
	}

	// 3. 调用LLM生成知识图谱
	h.logger.Infof("[K12Admin] Generating knowledge graph via LLM: %s/%s", req.Subject, req.Grade)
	startTime := time.Now()

	llmReq := &llm.KnowledgeGraphRequest{
		Subject: models.Subject(req.Subject),
		Grade:   models.GradeLevel(req.Grade),
	}

	result, err := h.llmService.GenerateKnowledgeGraph(ctx, llmReq)
	if err != nil {
		logs.Errorf("Failed to generate knowledge graph: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to generate knowledge graph: " + err.Error(),
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Infof("[K12Admin] Knowledge graph generated in %v", duration)

	// 4. 保存到缓存
	if err := h.cacheManager.Set(ctx, cacheKey, result, 24*time.Hour); err != nil {
		h.logger.Warnf("[K12Admin] Failed to cache result: %v", err)
	}

	// 5. 保存到数据库（可选）
	// TODO: 解析result并保存每个知识点到数据库

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"cached":   false,
			"content":  result,
			"duration": duration.Milliseconds(),
		},
	})
}

// ==================== API 2: 查询知识点 ====================

// QueryKnowledgePointsRequest 查询知识点请求
type QueryKnowledgePointsRequest struct {
	Subject     string `form:"subject"`      // 学科筛选
	Grade       string `form:"grade"`        // 年级筛选
	SearchQuery string `form:"q"`            // 搜索关键词
	Page        int    `form:"page"`         // 页码，从1开始
	PageSize    int    `form:"page_size"`    // 每页数量
}

// QueryKnowledgePoints 查询知识点列表
// GET /api/adapter/k12/admin/knowledge-points
func (h *K12AdminHandler) QueryKnowledgePoints(ctx context.Context, c *app.RequestContext) {
	var req QueryKnowledgePointsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大页面大小
	}

	// 从数据库查询
	var knowledgePoints []*database.KnowledgePointEntity
	var total int64

	query := h.db.Model(&database.KnowledgePointEntity{})

	// 应用筛选条件
	if req.Subject != "" {
		query = query.Where("subject = ?", req.Subject)
	}
	if req.Grade != "" {
		query = query.Where("grade = ?", req.Grade)
	}
	if req.SearchQuery != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+req.SearchQuery+"%", "%"+req.SearchQuery+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		logs.Errorf("Failed to count knowledge points: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to query knowledge points: " + err.Error(),
		})
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Order("created_at DESC").
		Find(&knowledgePoints).Error; err != nil {
		logs.Errorf("Failed to query knowledge points: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to query knowledge points: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"items":      knowledgePoints,
			"total":      total,
			"page":       req.Page,
			"page_size":  req.PageSize,
			"total_pages": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
		},
	})
}

// ==================== API 3: 刷新缓存 ====================

// RefreshCacheRequest 刷新缓存请求
type RefreshCacheRequest struct {
	CacheType string `json:"cache_type" binding:"required"` // knowledge_graph, knowledge_point, etc.
	Key       string `json:"key" binding:"required"`        // 缓存键值
}

// RefreshCache 刷新指定缓存
// POST /api/adapter/k12/admin/cache/refresh
func (h *K12AdminHandler) RefreshCache(ctx context.Context, c *app.RequestContext) {
	var req RefreshCacheRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}

	// 构建完整缓存键
	var fullKey string
	switch req.CacheType {
	case "knowledge_graph":
		fullKey = fmt.Sprintf("k12:kg:%s", req.Key)
	case "knowledge_point":
		fullKey = fmt.Sprintf("k12:kp:%s", req.Key)
	case "learning_resource":
		fullKey = fmt.Sprintf("k12:res:%s", req.Key)
	case "exercise":
		fullKey = fmt.Sprintf("k12:ex:%s", req.Key)
	default:
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  fmt.Sprintf("Unknown cache type: %s", req.CacheType),
		})
		return
	}

	// 删除缓存
	h.cacheManager.Delete(ctx, fullKey)

	h.logger.Infof("[K12Admin] Cache refreshed: %s", fullKey)

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "Cache refreshed successfully",
		"data": utils.H{
			"cache_key": fullKey,
		},
	})
}

// ==================== API 4: 清空缓存 ====================

// ClearCacheRequest 清空缓存请求
type ClearCacheRequest struct {
	CacheLevel string `json:"cache_level"` // memory, redis, mysql, all
}

// ClearCache 清空所有缓存
// DELETE /api/adapter/k12/admin/cache/clear
func (h *K12AdminHandler) ClearCache(ctx context.Context, c *app.RequestContext) {
	var req ClearCacheRequest
	if err := c.BindJSON(&req); err != nil {
		// 如果没有请求体，默认清空所有缓存
		req.CacheLevel = "all"
	}

	if req.CacheLevel == "" {
		req.CacheLevel = "all"
	}

	h.logger.Infof("[K12Admin] Clearing cache level: %s", req.CacheLevel)

	// 清空缓存
	if err := h.cacheManager.Clear(ctx); err != nil {
		logs.Errorf("Failed to clear cache: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to clear cache: " + err.Error(),
		})
		return
	}

	// 重置统计
	h.cacheManager.ResetStats()

	h.logger.Infof("[K12Admin] All caches cleared")

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "All caches cleared successfully",
		"data": utils.H{
			"cleared_at": time.Now(),
		},
	})
}

// ==================== API 5: 缓存统计 ====================

// GetCacheStats 获取缓存统计信息
// GET /api/adapter/k12/admin/cache/stats
func (h *K12AdminHandler) GetCacheStats(ctx context.Context, c *app.RequestContext) {
	stats := h.cacheManager.GetStats()

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"memory": utils.H{
				"hits":   stats.MemoryHits,
				"misses": stats.MemoryMisses,
				"hit_rate": func() float64 {
					total := stats.MemoryHits + stats.MemoryMisses
					if total == 0 {
						return 0
					}
					return float64(stats.MemoryHits) / float64(total)
				}(),
			},
			"redis": utils.H{
				"hits":   stats.RedisHits,
				"misses": stats.RedisMisses,
				"hit_rate": func() float64 {
					total := stats.RedisHits + stats.RedisMisses
					if total == 0 {
						return 0
					}
					return float64(stats.RedisHits) / float64(total)
				}(),
			},
			"mysql": utils.H{
				"hits":   stats.MySQLHits,
				"misses": stats.MySQLMisses,
				"hit_rate": func() float64 {
					total := stats.MySQLHits + stats.MySQLMisses
					if total == 0 {
						return 0
					}
					return float64(stats.MySQLHits) / float64(total)
				}(),
			},
			"total": utils.H{
				"hits":     stats.TotalHits,
				"misses":   stats.TotalMisses,
				"hit_rate": stats.HitRate,
			},
		},
	})
}

// ==================== API 6: 生成日志 ====================

// GetGenerationLogsRequest 获取生成日志请求
type GetGenerationLogsRequest struct {
	Subject  string `form:"subject"`   // 学科筛选
	Grade    string `form:"grade"`     // 年级筛选
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页数量
	Days     int    `form:"days"`      // 最近N天，默认7天
}

// GetGenerationLogs 获取生成日志
// GET /api/adapter/k12/admin/generation-logs
func (h *K12AdminHandler) GetGenerationLogs(ctx context.Context, c *app.RequestContext) {
	var req GetGenerationLogsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.PageSize > 200 {
		req.PageSize = 200
	}
	if req.Days <= 0 {
		req.Days = 7
	}

	// 从数据库查询生成日志
	var logs []*database.GenerationLogEntity
	var total int64

	query := h.db.Model(&database.GenerationLogEntity{})

	// 应用筛选条件
	if req.Subject != "" {
		query = query.Where("subject = ?", req.Subject)
	}
	if req.Grade != "" {
		query = query.Where("grade = ?", req.Grade)
	}

	// 时间范围筛选
	startDate := time.Now().AddDate(0, 0, -req.Days)
	query = query.Where("created_at >= ?", startDate)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		h.logger.Errorf("Failed to count generation logs: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to query generation logs: " + err.Error(),
		})
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		h.logger.Errorf("Failed to query generation logs: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to query generation logs: " + err.Error(),
		})
		return
	}

	// 计算统计信息
	var stats struct {
		TotalCount      int64
		SuccessCount    int64
		FailureCount    int64
		AvgResponseTime float64
	}

	h.db.Model(&database.GenerationLogEntity{}).
		Where("created_at >= ?", startDate).
		Count(&stats.TotalCount)

	h.db.Model(&database.GenerationLogEntity{}).
		Where("created_at >= ? AND success = ?", startDate, true).
		Count(&stats.SuccessCount)

	stats.FailureCount = stats.TotalCount - stats.SuccessCount

	// 计算平均响应时间（仅当有记录时）
	if stats.TotalCount > 0 {
		var result struct {
			AvgTime float64 `gorm:"column:avg_time"`
		}
		h.db.Model(&database.GenerationLogEntity{}).
			Where("created_at >= ?", startDate).
			Select("AVG(duration_ms) as avg_time").
			Scan(&result)
		stats.AvgResponseTime = result.AvgTime
	}

	// 计算成功率（避免除以零）
	successRate := 0.0
	if stats.TotalCount > 0 {
		successRate = float64(stats.SuccessCount) / float64(stats.TotalCount)
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"logs":        logs,
			"total":       total,
			"page":        req.Page,
			"page_size":   req.PageSize,
			"total_pages": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
			"statistics": utils.H{
				"total_count":       stats.TotalCount,
				"success_count":     stats.SuccessCount,
				"failure_count":     stats.FailureCount,
				"success_rate":      successRate,
				"avg_response_time": stats.AvgResponseTime,
			},
		},
	})
}

// ==================== 辅助方法 ====================

// GetSubjects 获取学科列表
// GET /api/adapter/k12/admin/subjects
func (h *K12AdminHandler) GetSubjects(ctx context.Context, c *app.RequestContext) {
	// 从数据库查询统计信息
	type SubjectStats struct {
		Subject          string `gorm:"column:subject"`
		KnowledgeCount   int64  `gorm:"column:knowledge_count"`
		ResourceCount    int64  `gorm:"column:resource_count"`
		ExerciseCount    int64  `gorm:"column:exercise_count"`
		LastGeneratedAt  *time.Time `gorm:"column:last_generated_at"`
	}

	var subjectStats []SubjectStats

	// 查询每个学科的统计信息
	err := h.db.Raw(`
		SELECT
			kp.subject,
			COUNT(DISTINCT kp.id) as knowledge_count,
			COUNT(DISTINCT lr.id) as resource_count,
			COUNT(DISTINCT ex.id) as exercise_count,
			MAX(kp.created_at) as last_generated_at
		FROM knowledge_points kp
		LEFT JOIN learning_resources lr ON kp.id = lr.knowledge_id
		LEFT JOIN exercises ex ON kp.id = ex.knowledge_id
		GROUP BY kp.subject
	`).Scan(&subjectStats).Error

	if err != nil {
		h.logger.Errorf("Failed to query subject stats: %v", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code": 500,
			"msg":  "Failed to query subjects: " + err.Error(),
		})
		return
	}

	// 学科元数据
	subjects := []utils.H{
		{"id": "math", "name": "数学", "icon": "📐"},
		{"id": "physics", "name": "物理", "icon": "⚛️"},
		{"id": "chemistry", "name": "化学", "icon": "🧪"},
		{"id": "biology", "name": "生物", "icon": "🧬"},
		{"id": "english", "name": "英语", "icon": "🔤"},
		{"id": "chinese", "name": "语文", "icon": "📖"},
	}

	// 合并统计数据
	for i := range subjects {
		subjectID := subjects[i]["id"].(string)
		for _, stat := range subjectStats {
			if stat.Subject == subjectID {
				subjects[i]["knowledge_count"] = stat.KnowledgeCount
				subjects[i]["resource_count"] = stat.ResourceCount
				subjects[i]["exercise_count"] = stat.ExerciseCount
				subjects[i]["last_generated_at"] = stat.LastGeneratedAt
				break
			}
		}
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": subjects,
	})
}

// GetSystemInfo 获取系统信息
// GET /api/adapter/k12/admin/system-info
func (h *K12AdminHandler) GetSystemInfo(ctx context.Context, c *app.RequestContext) {
	// 数据库统计
	var knowledgeCount, resourceCount, exerciseCount int64
	h.db.Model(&database.KnowledgePointEntity{}).Count(&knowledgeCount)
	h.db.Model(&database.LearningResourceEntity{}).Count(&resourceCount)
	h.db.Model(&database.ExerciseEntity{}).Count(&exerciseCount)

	// 缓存配置
	cacheConfigs, _ := h.dao.ListCacheConfigs(ctx)

	// 模型ID
	modelID := "100004"
	if modelIDStr := c.Query("model_id"); modelIDStr != "" {
		modelID = modelIDStr
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": 0,
		"msg":  "success",
		"data": utils.H{
			"version": "1.0.0",
			"database": utils.H{
				"knowledge_points":   knowledgeCount,
				"learning_resources": resourceCount,
				"exercises":          exerciseCount,
			},
			"cache": utils.H{
				"configs": cacheConfigs,
			},
			"llm": utils.H{
				"model_id":      modelID,
				"provider":      "Coze ModelBuilder",
				"default_timeout": "30s",
			},
		},
	})
}

// Close 关闭handler（清理资源）
func (h *K12AdminHandler) Close() error {
	if h.cacheManager != nil {
		return h.cacheManager.Close()
	}
	return nil
}

// ==================== 工具函数 ====================

// parseIntDefault 解析整数，失败返回默认值
func parseIntDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
