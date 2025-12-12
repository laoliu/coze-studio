/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * 节点推荐 Hook
 */

import { useState, useCallback, useEffect, useRef } from 'react';
import type {
  RecommendedNode,
  RecommendationRequest,
} from '../../typing/recommendation';
import { recommendationService } from '../../../services/recommendation-service';

interface UseNodeRecommendationParams {
  /** 源节点ID */
  sourceNodeId?: string;
  /** 源节点类型 */
  sourceNodeType?: string;
  /** 源节点名称 */
  sourceNodeName?: string;
  /** 工作流ID */
  workflowId?: string;
  /** 是否启用推荐 */
  enabled?: boolean;
  /** 推荐数量限制 */
  limit?: number;
  /** 输出格式 */
  outputFormat?: string;
}

interface UseNodeRecommendationReturn {
  /** 推荐列表 */
  recommendations: RecommendedNode[];
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: Error | null;
  /** 请求ID（用于反馈） */
  requestId: string;
  /** 重新获取推荐 */
  refetch: () => Promise<void>;
  /** 记录用户选择 */
  recordSelection: (nodeType: string) => Promise<void>;
  /** 记录用户忽略 */
  recordDismiss: () => Promise<void>;
}

/**
 * 节点推荐 Hook
 * 
 * @example
 * const { recommendations, loading, recordSelection } = useNodeRecommendation({
 *   sourceNodeId: 'node-123',
 *   sourceNodeType: 'LLM',
 *   workflowId: 'wf-456',
 *   enabled: true,
 * });
 */
export function useNodeRecommendation(
  params: UseNodeRecommendationParams
): UseNodeRecommendationReturn {
  const [recommendations, setRecommendations] = useState<RecommendedNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [requestId, setRequestId] = useState<string>('');

  // 使用ref避免重复请求
  const abortControllerRef = useRef<AbortController | null>(null);

  /**
   * 获取推荐
   */
  const fetchRecommendations = useCallback(async () => {
    // 检查必要参数
    if (!params.enabled || !params.sourceNodeId || !params.sourceNodeType || !params.workflowId) {
      setRecommendations([]);
      return;
    }

    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    abortControllerRef.current = new AbortController();
    setLoading(true);
    setError(null);

    try {
      const requestParams: Omit<RecommendationRequest, 'includeReason'> = {
        workflowId: params.workflowId,
        sourceNodeId: params.sourceNodeId,
        sourceNodeType: params.sourceNodeType,
        sourceNodeName: params.sourceNodeName,
        outputFormat: params.outputFormat || 'json',
        limit: params.limit || 10,
      };

      const response = await recommendationService.getRecommendations(requestParams);

      setRecommendations(response.recommendations);
      setRequestId(response.metadata.requestId);
      setError(null);
    } catch (err) {
      // 忽略取消的请求
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      
      console.error('[NodeRecommendation] Failed to fetch recommendations:', err);
      setError(err as Error);
      setRecommendations([]);
    } finally {
      setLoading(false);
      abortControllerRef.current = null;
    }
  }, [
    params.enabled,
    params.sourceNodeId,
    params.sourceNodeType,
    params.sourceNodeName,
    params.workflowId,
    params.outputFormat,
    params.limit,
  ]);

  /**
   * 记录用户选择的节点
   */
  const recordSelection = useCallback(
    async (nodeType: string) => {
      if (!requestId) {
        console.warn('[NodeRecommendation] No requestId available for feedback');
        return;
      }

      try {
        await recommendationService.recordFeedback({
          requestId,
          selectedNode: nodeType,
          userAction: 'selected',
        });
        console.log('[NodeRecommendation] Feedback recorded:', nodeType);
      } catch (err) {
        console.error('[NodeRecommendation] Failed to record feedback:', err);
        // 不抛出错误，避免影响用户体验
      }
    },
    [requestId]
  );

  /**
   * 记录用户忽略推荐
   */
  const recordDismiss = useCallback(async () => {
    if (!requestId) {
      return;
    }

    try {
      await recommendationService.recordFeedback({
        requestId,
        userAction: 'dismissed',
      });
      console.log('[NodeRecommendation] Dismiss recorded');
    } catch (err) {
      console.error('[NodeRecommendation] Failed to record dismiss:', err);
    }
  }, [requestId]);

  // 自动获取推荐
  useEffect(() => {
    fetchRecommendations();

    // 清理函数
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [fetchRecommendations]);

  return {
    recommendations,
    loading,
    error,
    requestId,
    refetch: fetchRecommendations,
    recordSelection,
    recordDismiss,
  };
}
