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

/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * 节点推荐API服务
 */

import type {
  RecommendationRequest,
  RecommendationResponse,
  FeedbackRequest,
  FeedbackResponse,
} from '../typing/recommendation';
import { API_CONFIG } from '../config/api';

/**
 * 节点推荐服务类
 */
export class RecommendationService {
  private baseURL = API_CONFIG.RECOMMENDATION;

  /**
   * 获取节点推荐
   */
  async getRecommendations(
    params: Omit<RecommendationRequest, 'includeReason'>,
  ): Promise<RecommendationResponse> {
    const DEFAULT_LIMIT = 10;
    const response = await fetch(`${this.baseURL}/recommend`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // 发送 cookies
      body: JSON.stringify({
        workflow_id: params.workflowId,
        source_node_id: params.sourceNodeId,
        source_node_type: params.sourceNodeType,
        source_node_name: params.sourceNodeName,
        output_format: params.outputFormat,
        limit: params.limit || DEFAULT_LIMIT,
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
        recommendations: (data.recommendations || []).map((rec: unknown) => {
          const recData = rec as Record<string, unknown>;
          return {
            nodeType: recData.type as string,
            displayName: recData.type as string, // 测试服务器没有 display_name
            score: recData.score as number,
            reason: recData.reason as string,
            category: recData.category as string,
            metadata: recData.metadata as Record<string, unknown>,
          };
        }),
        metadata: {
          requestId: data.request_id,
          totalCandidates: data.recommendations?.length || 0,
        },
      };
    }

    // 生产服务器格式
    return {
      recommendations: (data.recommendations || []).map((rec: unknown) => {
        const recData = rec as Record<string, unknown>;
        return {
          nodeType: recData.node_type as string,
          displayName: recData.display_name as string,
          score: recData.score as number,
          reason: recData.reason as string,
          category: recData.category as string,
          suggestedConfig: recData.suggested_config as Record<string, unknown>,
          strategySource: recData.strategy_source as string,
        };
      }),
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
    // 测试服务器使用 /api/feedback，生产服务器使用 /api/workflow_api/recommend/feedback
    const feedbackPath =
      this.baseURL === '/api'
        ? `${this.baseURL}/feedback`
        : `${this.baseURL}/recommend/feedback`;

    const response = await fetch(feedbackPath, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // 发送 cookies
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
