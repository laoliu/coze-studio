/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * 节点推荐 React Hook
 */

import { useState, useCallback, useMemo } from 'react';
import { RecommendationService } from '../services/recommendation-service';
import type {
  RecommendationRequest,
  RecommendedNode,
} from '../typing/recommendation';

const recommendationService = new RecommendationService();

export interface UseNodeRecommendationOptions {
  /** 工作流ID */
  workflowId: string;
  /** 是否自动加载推荐 */
  autoLoad?: boolean;
}

export interface UseNodeRecommendationResult {
  /** 推荐列表 */
  recommendations: RecommendedNode[];
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: Error | null;
  /** 请求ID */
  requestId: string;
  /** 加载推荐 */
  loadRecommendations: (
    params: Omit<RecommendationRequest, 'workflowId' | 'includeReason'>
  ) => Promise<void>;
  /** 提交反馈 */
  submitFeedback: (
    selectedType: string,
    userAction: 'selected' | 'dismissed' | 'ignored'
  ) => Promise<void>;
  /** 清空推荐 */
  clearRecommendations: () => void;
  /** 重新加载 */
  reload: () => Promise<void>;
}

/**
 * 节点推荐 Hook
 *
 * @example
 * ```tsx
 * const { recommendations, loading, loadRecommendations } = useNodeRecommendation({
 *   workflowId: 'my-workflow-id'
 * });
 *
 * // 加载推荐
 * useEffect(() => {
 *   if (selectedNode) {
 *     loadRecommendations({
 *       sourceNodeId: selectedNode.id,
 *       sourceNodeType: selectedNode.type,
 *       limit: 5
 *     });
 *   }
 * }, [selectedNode]);
 *
 * // 渲染推荐
 * <RecommendationSection
 *   recommendations={recommendations}
 *   loading={loading}
 *   onSelect={handleSelect}
 * />
 * ```
 */
export function useNodeRecommendation(
  options: UseNodeRecommendationOptions
): UseNodeRecommendationResult {
  const [recommendations, setRecommendations] = useState<RecommendedNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [requestId, setRequestId] = useState<string>('');
  const [lastParams, setLastParams] = useState<
    Omit<RecommendationRequest, 'workflowId' | 'includeReason'> | null
  >(null);

  /**
   * 加载推荐
   */
  const loadRecommendations = useCallback(
    async (params: Omit<RecommendationRequest, 'workflowId' | 'includeReason'>) => {
      try {
        setLoading(true);
        setError(null);
        setLastParams(params);

        const response = await recommendationService.getRecommendations({
          workflowId: options.workflowId,
          ...params,
        });

        setRecommendations(response.recommendations);
        setRequestId(response.metadata.requestId);
      } catch (err) {
        const error = err as Error;
        setError(error);
        setRecommendations([]);
        console.error('[useNodeRecommendation] Failed to load recommendations:', error);
      } finally {
        setLoading(false);
      }
    },
    [options.workflowId]
  );

  /**
   * 重新加载（使用上次的参数）
   */
  const reload = useCallback(async () => {
    if (lastParams) {
      await loadRecommendations(lastParams);
    }
  }, [lastParams, loadRecommendations]);

  /**
   * 提交反馈
   */
  const submitFeedback = useCallback(
    async (selectedType: string, userAction: 'selected' | 'dismissed' | 'ignored') => {
      if (!requestId) {
        console.warn('[useNodeRecommendation] No request ID, skipping feedback');
        return;
      }

      try {
        await recommendationService.recordFeedback({
          requestId: requestId,
          selectedNode: selectedType,
          userAction,
        });

        console.log('[useNodeRecommendation] Feedback submitted successfully', {
          selectedType,
          userAction,
        });
      } catch (err) {
        console.error('[useNodeRecommendation] Failed to submit feedback:', err);
      }
    },
    [requestId]
  );

  /**
   * 清空推荐
   */
  const clearRecommendations = useCallback(() => {
    setRecommendations([]);
    setRequestId('');
    setError(null);
    setLastParams(null);
  }, []);

  return useMemo(
    () => ({
      recommendations,
      loading,
      error,
      requestId,
      loadRecommendations,
      submitFeedback,
      clearRecommendations,
      reload,
    }),
    [
      recommendations,
      loading,
      error,
      requestId,
      loadRecommendations,
      submitFeedback,
      clearRecommendations,
      reload,
    ]
  );
}
