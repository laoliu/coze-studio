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

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/infra/database"
)

// ==================== 缓存配置 ====================

// CacheConfig 缓存配置
type CacheConfig struct {
	EnableMemory bool          // 启用内存缓存
	EnableRedis  bool          // 启用Redis缓存
	EnableMySQL  bool          // 启用MySQL缓存
	DefaultTTL   time.Duration // 默认TTL
	RedisAddr    string        // Redis地址
	RedisDB      int           // Redis数据库编号
}

// DefaultCacheConfig 返回默认配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		EnableMemory: true,
		EnableRedis:  false, // 默认不启用Redis（可选）
		EnableMySQL:  true,
		DefaultTTL:   24 * time.Hour, // 默认24小时
		RedisAddr:    "localhost:6379",
		RedisDB:      0,
	}
}

// ==================== 缓存项 ====================

// CacheEntry 缓存项
type CacheEntry struct {
	Key        string    `json:"key"`         // 缓存键
	Value      string    `json:"value"`       // 缓存值（JSON）
	ExpireAt   time.Time `json:"expire_at"`   // 过期时间
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	AccessedAt time.Time `json:"accessed_at"` // 最后访问时间
	HitCount   int64     `json:"hit_count"`   // 命中次数
}

// IsExpired 检查是否过期
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpireAt)
}

// ==================== 缓存统计 ====================

// CacheStats 缓存统计
type CacheStats struct {
	MemoryHits   int64 `json:"memory_hits"`   // 内存缓存命中
	MemoryMisses int64 `json:"memory_misses"` // 内存缓存未命中
	RedisHits    int64 `json:"redis_hits"`    // Redis缓存命中
	RedisMisses  int64 `json:"redis_misses"`  // Redis缓存未命中
	MySQLHits    int64 `json:"mysql_hits"`    // MySQL缓存命中
	MySQLMisses  int64 `json:"mysql_misses"`  // MySQL缓存未命中
	TotalHits    int64 `json:"total_hits"`    // 总命中
	TotalMisses  int64 `json:"total_misses"`  // 总未命中
	HitRate      float64 `json:"hit_rate"`    // 命中率
}

// ==================== 缓存管理器 ====================

// CacheManager 3级缓存管理器
type CacheManager struct {
	config      *CacheConfig
	memoryCache *sync.Map              // L1: 内存缓存
	redisClient *redis.Client          // L2: Redis缓存
	db          *gorm.DB               // L3: MySQL缓存
	dao         database.K12DAO        // DAO层（接口）
	stats       *CacheStats            // 统计信息
	statsMutex  sync.RWMutex           // 统计锁
	logger      hlog.FullLogger
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(config *CacheConfig, db *gorm.DB, logger hlog.FullLogger) (*CacheManager, error) {
	if config == nil {
		config = DefaultCacheConfig()
	}

	if logger == nil {
		logger = hlog.DefaultLogger()
	}

	cm := &CacheManager{
		config:      config,
		memoryCache: &sync.Map{},
		db:          db,
		dao:         database.NewK12DAO(db),
		stats:       &CacheStats{},
		logger:      logger,
	}

	// 初始化Redis（如果启用）
	if config.EnableRedis {
		cm.redisClient = redis.NewClient(&redis.Options{
			Addr: config.RedisAddr,
			DB:   config.RedisDB,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := cm.redisClient.Ping(ctx).Err(); err != nil {
			logger.Warnf("[CacheManager] Redis connection failed: %v, Redis cache disabled", err)
			cm.config.EnableRedis = false
		} else {
			logger.Infof("[CacheManager] Redis cache enabled")
		}
	}

	logger.Infof("[CacheManager] Initialized: Memory=%v, Redis=%v, MySQL=%v",
		config.EnableMemory, config.EnableRedis, config.EnableMySQL)

	return cm, nil
}

// ==================== 缓存键生成 ====================

// generateKnowledgeGraphKey 生成知识图谱缓存键
func (cm *CacheManager) generateKnowledgeGraphKey(subject, grade string) string {
	return fmt.Sprintf("k12:kg:%s:%s", subject, grade)
}

// generateKnowledgePointKey 生成知识点缓存键
func (cm *CacheManager) generateKnowledgePointKey(knowledgeID string) string {
	return fmt.Sprintf("k12:kp:%s", knowledgeID)
}

// generateResourceKey 生成学习资源缓存键
func (cm *CacheManager) generateResourceKey(knowledgeID, resourceType string) string {
	return fmt.Sprintf("k12:res:%s:%s", knowledgeID, resourceType)
}

// generateExerciseKey 生成习题缓存键
func (cm *CacheManager) generateExerciseKey(knowledgeID, difficulty string) string {
	return fmt.Sprintf("k12:ex:%s:%s", knowledgeID, difficulty)
}

// ==================== L1: 内存缓存 ====================

// getFromMemory 从内存缓存获取
func (cm *CacheManager) getFromMemory(key string) (*CacheEntry, bool) {
	if !cm.config.EnableMemory {
		return nil, false
	}

	value, ok := cm.memoryCache.Load(key)
	if !ok {
		cm.incrementStat("memory_misses")
		return nil, false
	}

	entry, ok := value.(*CacheEntry)
	if !ok {
		cm.memoryCache.Delete(key)
		cm.incrementStat("memory_misses")
		return nil, false
	}

	// 检查过期
	if entry.IsExpired() {
		cm.memoryCache.Delete(key)
		cm.incrementStat("memory_misses")
		return nil, false
	}

	// 更新访问信息
	entry.AccessedAt = time.Now()
	entry.HitCount++

	cm.incrementStat("memory_hits")
	cm.logger.Debugf("[Cache:Memory] Hit: %s", key)
	return entry, true
}

// setToMemory 设置到内存缓存
func (cm *CacheManager) setToMemory(key string, entry *CacheEntry) {
	if !cm.config.EnableMemory {
		return
	}

	cm.memoryCache.Store(key, entry)
	cm.logger.Debugf("[Cache:Memory] Set: %s", key)
}

// deleteFromMemory 从内存缓存删除
func (cm *CacheManager) deleteFromMemory(key string) {
	if !cm.config.EnableMemory {
		return
	}

	cm.memoryCache.Delete(key)
	cm.logger.Debugf("[Cache:Memory] Delete: %s", key)
}

// ==================== L2: Redis缓存 ====================

// getFromRedis 从Redis缓存获取
func (cm *CacheManager) getFromRedis(ctx context.Context, key string) (*CacheEntry, bool) {
	if !cm.config.EnableRedis || cm.redisClient == nil {
		return nil, false
	}

	data, err := cm.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			cm.logger.Warnf("[Cache:Redis] Get error: %v", err)
		}
		cm.incrementStat("redis_misses")
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		cm.logger.Warnf("[Cache:Redis] Unmarshal error: %v", err)
		cm.incrementStat("redis_misses")
		return nil, false
	}

	// 检查过期
	if entry.IsExpired() {
		cm.redisClient.Del(ctx, key)
		cm.incrementStat("redis_misses")
		return nil, false
	}

	// 更新访问信息
	entry.AccessedAt = time.Now()
	entry.HitCount++

	cm.incrementStat("redis_hits")
	cm.logger.Debugf("[Cache:Redis] Hit: %s", key)

	// 回填到内存缓存
	cm.setToMemory(key, &entry)

	return &entry, true
}

// setToRedis 设置到Redis缓存
func (cm *CacheManager) setToRedis(ctx context.Context, key string, entry *CacheEntry) error {
	if !cm.config.EnableRedis || cm.redisClient == nil {
		return nil
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	ttl := time.Until(entry.ExpireAt)
	if ttl <= 0 {
		ttl = cm.config.DefaultTTL
	}

	if err := cm.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	cm.logger.Debugf("[Cache:Redis] Set: %s, TTL: %v", key, ttl)
	return nil
}

// deleteFromRedis 从Redis缓存删除
func (cm *CacheManager) deleteFromRedis(ctx context.Context, key string) {
	if !cm.config.EnableRedis || cm.redisClient == nil {
		return
	}

	cm.redisClient.Del(ctx, key)
	cm.logger.Debugf("[Cache:Redis] Delete: %s", key)
}

// ==================== L3: MySQL缓存 ====================

// getFromMySQL 从MySQL缓存获取
func (cm *CacheManager) getFromMySQL(ctx context.Context, key string) (*CacheEntry, bool) {
	if !cm.config.EnableMySQL {
		return nil, false
	}

	// 从cache_config表查询TTL配置
	configs, err := cm.dao.ListCacheConfigs(ctx)
	if err != nil {
		cm.logger.Warnf("[Cache:MySQL] Get config error: %v", err)
		cm.incrementStat("mysql_misses")
		return nil, false
	}

	// 根据key类型查找配置
	var cacheType string
	if len(key) > 7 {
		switch key[:7] {
		case "k12:kg:":
			cacheType = "knowledge_graph"
		case "k12:kp:":
			cacheType = "knowledge_point"
		}
	}

	// 查找对应的配置
	var foundConfig *database.CacheConfigEntity
	for i := range configs {
		if configs[i].ItemType == cacheType {
			foundConfig = configs[i]
			break
		}
	}

	if foundConfig != nil {
		cm.logger.Debugf("[Cache:MySQL] Found config for %s: TTL=%ds", cacheType, foundConfig.TTLSeconds)
	}

	// 这里简化处理：根据key类型从不同表查询
	// 实际应该根据key解析出具体的查询参数
	// 这里只是示例框架

	cm.incrementStat("mysql_misses")
	cm.logger.Debugf("[Cache:MySQL] Miss: %s", key)
	return nil, false
}

// setToMySQL 设置到MySQL缓存（实际是更新数据库记录）
func (cm *CacheManager) setToMySQL(ctx context.Context, key string, entry *CacheEntry) error {
	if !cm.config.EnableMySQL {
		return nil
	}

	// MySQL缓存实际是将数据持久化到数据库表中
	// 这里已经通过DAO层操作完成，不需要额外存储
	cm.logger.Debugf("[Cache:MySQL] Set: %s (via DAO)", key)
	return nil
}

// ==================== 公共方法 ====================

// Get 获取缓存（自动尝试3级缓存）
func (cm *CacheManager) Get(ctx context.Context, key string) (interface{}, bool) {
	// L1: 内存缓存
	if entry, ok := cm.getFromMemory(key); ok {
		var result interface{}
		if err := json.Unmarshal([]byte(entry.Value), &result); err == nil {
			return result, true
		}
	}

	// L2: Redis缓存
	if entry, ok := cm.getFromRedis(ctx, key); ok {
		var result interface{}
		if err := json.Unmarshal([]byte(entry.Value), &result); err == nil {
			return result, true
		}
	}

	// L3: MySQL缓存
	if entry, ok := cm.getFromMySQL(ctx, key); ok {
		var result interface{}
		if err := json.Unmarshal([]byte(entry.Value), &result); err == nil {
			// 回填到上层缓存
			cm.setToMemory(key, entry)
			_ = cm.setToRedis(ctx, key, entry)
			return result, true
		}
	}

	return nil, false
}

// Set 设置缓存（写入所有级别）
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if ttl == 0 {
		ttl = cm.config.DefaultTTL
	}

	// 创建缓存项
	entry := &CacheEntry{
		Key:        key,
		Value:      string(data),
		ExpireAt:   time.Now().Add(ttl),
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
		HitCount:   0,
	}

	// 写入所有级别
	cm.setToMemory(key, entry)
	if err := cm.setToRedis(ctx, key, entry); err != nil {
		cm.logger.Warnf("[CacheManager] Redis set error: %v", err)
	}
	if err := cm.setToMySQL(ctx, key, entry); err != nil {
		cm.logger.Warnf("[CacheManager] MySQL set error: %v", err)
	}

	return nil
}

// Delete 删除缓存（所有级别）
func (cm *CacheManager) Delete(ctx context.Context, key string) {
	cm.deleteFromMemory(key)
	cm.deleteFromRedis(ctx, key)
	// MySQL缓存通过DAO层的Delete方法处理

	cm.logger.Infof("[CacheManager] Deleted: %s", key)
}

// Clear 清空所有缓存
func (cm *CacheManager) Clear(ctx context.Context) error {
	// 清空内存缓存
	if cm.config.EnableMemory {
		cm.memoryCache = &sync.Map{}
		cm.logger.Infof("[CacheManager] Memory cache cleared")
	}

	// 清空Redis缓存
	if cm.config.EnableRedis && cm.redisClient != nil {
		if err := cm.redisClient.FlushDB(ctx).Err(); err != nil {
			cm.logger.Warnf("[CacheManager] Redis flush error: %v", err)
		} else {
			cm.logger.Infof("[CacheManager] Redis cache cleared")
		}
	}

	// MySQL缓存需要通过DAO层删除数据
	cm.logger.Infof("[CacheManager] All caches cleared")

	return nil
}

// ==================== 统计方法 ====================

// incrementStat 增加统计计数
func (cm *CacheManager) incrementStat(statName string) {
	cm.statsMutex.Lock()
	defer cm.statsMutex.Unlock()

	switch statName {
	case "memory_hits":
		cm.stats.MemoryHits++
		cm.stats.TotalHits++
	case "memory_misses":
		cm.stats.MemoryMisses++
		cm.stats.TotalMisses++
	case "redis_hits":
		cm.stats.RedisHits++
		cm.stats.TotalHits++
	case "redis_misses":
		cm.stats.RedisMisses++
		cm.stats.TotalMisses++
	case "mysql_hits":
		cm.stats.MySQLHits++
		cm.stats.TotalHits++
	case "mysql_misses":
		cm.stats.MySQLMisses++
		cm.stats.TotalMisses++
	}

	// 更新命中率
	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total > 0 {
		cm.stats.HitRate = float64(cm.stats.TotalHits) / float64(total)
	}
}

// GetStats 获取缓存统计
func (cm *CacheManager) GetStats() *CacheStats {
	cm.statsMutex.RLock()
	defer cm.statsMutex.RUnlock()

	// 返回副本
	statsCopy := *cm.stats
	return &statsCopy
}

// ResetStats 重置统计
func (cm *CacheManager) ResetStats() {
	cm.statsMutex.Lock()
	defer cm.statsMutex.Unlock()

	cm.stats = &CacheStats{}
	cm.logger.Infof("[CacheManager] Stats reset")
}

// ==================== 清理方法 ====================

// Close 关闭缓存管理器
func (cm *CacheManager) Close() error {
	if cm.redisClient != nil {
		if err := cm.redisClient.Close(); err != nil {
			return err
		}
	}

	cm.logger.Infof("[CacheManager] Closed")
	return nil
}
