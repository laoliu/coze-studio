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
	"fmt"
)

// GenerateKnowledgeGraphKey 生成知识图谱缓存键（公开方法）
func (cm *CacheManager) GenerateKnowledgeGraphKey(subject, grade string) string {
	return cm.generateKnowledgeGraphKey(subject, grade)
}

// GenerateKnowledgePointKey 生成知识点缓存键（公开方法）
func (cm *CacheManager) GenerateKnowledgePointKey(knowledgeID string) string {
	return cm.generateKnowledgePointKey(knowledgeID)
}

// GenerateResourceKey 生成学习资源缓存键（公开方法）
func (cm *CacheManager) GenerateResourceKey(knowledgeID, resourceType string) string {
	return cm.generateResourceKey(knowledgeID, resourceType)
}

// GenerateExerciseKey 生成习题缓存键（公开方法）
func (cm *CacheManager) GenerateExerciseKey(knowledgeID, difficulty string) string {
	return cm.generateExerciseKey(knowledgeID, difficulty)
}

// ParseCacheKey 解析缓存键，返回类型和参数
func ParseCacheKey(key string) (cacheType string, params map[string]string, err error) {
	// key格式: "k12:type:param1:param2..."
	if len(key) < 7 || key[:4] != "k12:" {
		return "", nil, fmt.Errorf("invalid cache key format: %s", key)
	}

	params = make(map[string]string)

	// 根据前缀判断类型
	switch {
	case len(key) > 7 && key[:7] == "k12:kg:":
		// knowledge_graph: k12:kg:subject:grade
		cacheType = "knowledge_graph"
		// 这里可以进一步解析subject和grade

	case len(key) > 7 && key[:7] == "k12:kp:":
		// knowledge_point: k12:kp:knowledge_id
		cacheType = "knowledge_point"
		params["knowledge_id"] = key[7:]

	case len(key) > 8 && key[:8] == "k12:res:":
		// resource: k12:res:knowledge_id:resource_type
		cacheType = "learning_resource"
		// 解析参数

	case len(key) > 7 && key[:7] == "k12:ex:":
		// exercise: k12:ex:knowledge_id:difficulty
		cacheType = "exercise"
		// 解析参数

	default:
		return "", nil, fmt.Errorf("unknown cache key prefix: %s", key)
	}

	return cacheType, params, nil
}
