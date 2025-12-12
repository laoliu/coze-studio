/*
 * Copyright 2025 Coze Studio. All rights reserved.
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/recommendation"
)

// 推荐请求（简化版）
type RecommendRequest struct {
	WorkflowID     string                 `json:"workflow_id"`
	SourceNodeID   string                 `json:"source_node_id"`
	SourceNodeType string                 `json:"source_node_type"`
	SourceOutputs  map[string]interface{} `json:"source_outputs,omitempty"`
	Limit          int                    `json:"limit,omitempty"`
}

// 推荐响应（简化版）
type RecommendResponse struct {
	Success         bool                 `json:"success"`
	RequestID       string               `json:"request_id"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Error           string               `json:"error,omitempty"`
}

type RecommendationItem struct {
	Type     string  `json:"type"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason"`
	Category string  `json:"category"`
}

var engine *recommendation.Engine

// 初始化推荐引擎
func init() {
	var err error
	engine, err = recommendation.NewEngine(nil)
	if err != nil {
		log.Fatalf("Failed to initialize recommendation engine: %v", err)
	}
	log.Println("✅ Recommendation engine initialized successfully")
}

// 推荐接口
func handleRecommend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RecommendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body: "+err.Error())
		return
	}

	// 验证必填字段
	if req.WorkflowID == "" || req.SourceNodeID == "" || req.SourceNodeType == "" {
		sendError(w, "Missing required fields")
		return
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = 5
	}

	// 转换为推荐引擎请求
	engineReq := &recommendation.RecommendRequest{
		WorkflowID:     req.WorkflowID,
		SourceNodeID:   req.SourceNodeID,
		SourceNodeType: entity.NodeType(req.SourceNodeType),
		SourceOutputs:  convertOutputs(req.SourceOutputs),
		Limit:          req.Limit,
	}

	// 调用推荐引擎
	resp, err := engine.Recommend(context.Background(), engineReq)
	if err != nil {
		sendError(w, "Recommendation failed: "+err.Error())
		return
	}

	// 转换响应
	items := make([]RecommendationItem, len(resp.Recommendations))
	for i, rec := range resp.Recommendations {
		items[i] = RecommendationItem{
			Type:     string(rec.NodeType),
			Score:    rec.Score,
			Reason:   rec.Reason,
			Category: rec.Category,
		}
	}

	response := RecommendResponse{
		Success:         true,
		RequestID:       resp.Meta.RecommendationID,
		Recommendations: items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 反馈接口
func handleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var feedback struct {
		RecommendationID string `json:"recommendation_id"`
		SelectedType     string `json:"selected_type"`
		UserAction       string `json:"user_action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		sendError(w, "Invalid request body: "+err.Error())
		return
	}

	// 记录反馈
	engineFeedback := &recommendation.RecommendationFeedback{
		RecommendationID: feedback.RecommendationID,
		SelectedType:     entity.NodeType(feedback.SelectedType),
		IsUseful:         feedback.UserAction == "accepted",
	}

	err := engine.RecordFeedback(context.Background(), engineFeedback)
	if err != nil {
		sendError(w, "Failed to record feedback: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Feedback recorded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 健康检查
func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "node-recommendation",
		"version": "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 辅助函数
func convertOutputs(outputs map[string]interface{}) map[string]*vo.TypeInfo {
	if outputs == nil {
		return nil
	}

	result := make(map[string]*vo.TypeInfo)
	for key, val := range outputs {
		if m, ok := val.(map[string]interface{}); ok {
			typeInfo := &vo.TypeInfo{}
			if t, ok := m["type"].(string); ok {
				typeInfo.Type = vo.DataType(t)
			}
			if desc, ok := m["desc"].(string); ok {
				typeInfo.Desc = desc
			}
			result[key] = typeInfo
		}
	}
	return result
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(RecommendResponse{
		Success: false,
		Error:   message,
	})
}

// CORS 中间件
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 允许所有来源（测试用）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// 注册路由（带 CORS 支持）
	http.HandleFunc("/api/recommend", corsMiddleware(handleRecommend))
	http.HandleFunc("/api/feedback", corsMiddleware(handleFeedback))
	http.HandleFunc("/health", corsMiddleware(handleHealth))

	// 启动服务器
	port := "8080"
	fmt.Printf("\n╔══════════════════════════════════════════╗\n")
	fmt.Printf("║   节点推荐 API 测试服务器               ║\n")
	fmt.Printf("╚══════════════════════════════════════════╝\n\n")
	fmt.Printf("🚀 Server starting on http://localhost:%s\n\n", port)
	fmt.Printf("📌 Available endpoints:\n")
	fmt.Printf("   POST /api/recommend  - 获取节点推荐\n")
	fmt.Printf("   POST /api/feedback   - 提交用户反馈\n")
	fmt.Printf("   GET  /health         - 健康检查\n\n")
	fmt.Printf("🧪 Test with:\n")
	fmt.Printf("   curl -X POST http://localhost:%s/api/recommend \\\n", port)
	fmt.Printf("     -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("     -d '{\"workflow_id\":\"test\",\"source_node_id\":\"1\",\"source_node_type\":\"LLM\",\"limit\":5}'\n\n")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
