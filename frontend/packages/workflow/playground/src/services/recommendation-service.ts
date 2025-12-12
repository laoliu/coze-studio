/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * 节点推荐API服务
 */

import { API_CONFIG } from '../config/api';
import type {
  RecommendationRequest,
  RecommendationResponse,
  FeedbackRequest,
  FeedbackResponse,
} from '../typing/recommendation';

/**
 * 节点推荐服务类
 */
export class RecommendationService {
  private baseURL = API_CONFIG.RECOMMENDATION;

  /**
   * 获取节点推荐
   */
  async getRecommendations(
    params: Omit<RecommendationRequest, 'includeReason'>
  ): Promise<RecommendationResponse> {
    const response = await fetch(`${this.baseURL}/recommend`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        workflow_id: params.workflowId,
        source_node_id: params.sourceNodeId,
        source_node_type: params.sourceNodeType,
        source_node_name: params.sourceNodeName,
        output_format: params.outputFormat,
        limit: params.limit || 10,
        include_reason: true,
        workflow_context: params.workflowContext
          ? {
              workflow_name: params.workflowContext.workflowName,
              workflow_objective: params.workflowContext.workflowObjective,
              total_nodes: params.workflowContext.totalNodes,
              has_database_node: params.workflowContext.hasDatabaseNode,
              has_llm_node: params.workflowContext.hasLLMNode,
              has_loop_node: params.workflowContext.hasLoopNode,
              inside_loop: params.workflowContext.insideLoop,
              is_chat_workflow: params.workflowContext.isChatWorkflow,
            }
          : undefined,
      }),
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`API error: ${response.status} - ${error}`);
    }

    const data = await response.json();
    
    // 兼容测试服务器和生产服务器的响应格式
    // 测试服务器: { success, request_id, recommendations: [{type, score, reason}] }
    // 生产服务器: { recommendations: [{node_type, ...}], metadata: {request_id, ...} }
    const isTestServerResponse = data.success !== undefined;
    
    if (isTestServerResponse) {
      // 测试服务器格式
      return {
        recommendations: (data.recommendations || []).map((rec: any) => ({
          nodeType: rec.type,
          displayName: rec.type, // 测试服务器没有 display_name
          score: rec.score,
          reason: rec.reason,
          category: rec.category,
          metadata: rec.metadata,
        })),
        metadata: {
          requestId: data.request_id,
          totalCandidates: data.recommendations?.length || 0,
        },
      };
    }
    
    // 生产服务器格式
    return {
      recommendations: (data.recommendations || []).map((rec: any) => ({
        nodeType: rec.node_type,
        displayName: rec.display_name,
        score: rec.score,
        reason: rec.reason,
        category: rec.category,
        suggestedConfig: rec.suggested_config,
        strategySource: rec.strategy_source,
      })),
      metadata: {
        requestId: data.metadata?.request_id,
        totalCandidates: data.metadata?.total_candidates,
        executionTimeMs: data.metadata?.execution_time_ms,
        strategiesUsed: data.metadata?.strategies_used,
      },
    };
  }

  /**
   * 记录用户反馈
   */
  async recordFeedback(params: FeedbackRequest): Promise<FeedbackResponse> {
    // 测试服务器使用 /api/feedback，生产服务器使用 /api/workflow_api/node/recommend/feedback
    const feedbackPath = this.baseURL === '/api' 
      ? `${this.baseURL}/feedback`
      : `${this.baseURL}/recommend/feedback`;
      
    const response = await fetch(feedbackPath, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        request_id: params.requestId,
        selected_node: params.selectedNode,
        user_action: params.userAction,
        user_comment: params.userComment,
      }),
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`Feedback API error: ${response.status} - ${error}`);
    }

    return response.json();
  }
}

/**
 * 单例服务实例
 */
export const recommendationService = new RecommendationService();
