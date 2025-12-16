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
 * 推荐区域主组件
 */

import React, { useMemo } from 'react';
import { RecommendationCard } from './RecommendationCard';
import { useNodeRecommendation } from '../../hooks/use-node-recommendation';
import styles from './styles.module.less';

export interface RecommendationSectionProps {
  /** 源节点ID */
  sourceNodeId: string;
  /** 源节点类型 */
  sourceNodeType: string;
  /** 源节点名称 */
  sourceNodeName?: string;
  /** 工作流ID */
  workflowId: string;
  /** 选择节点回调 */
  onSelectNode?: (nodeType: string) => void;
  /** 最大显示数量 */
  maxDisplay?: number;
  /** 是否启用 */
  enabled?: boolean;
}

/**
 * 推荐区域组件
 */
export const RecommendationSection: React.FC<RecommendationSectionProps> = ({
  sourceNodeId,
  sourceNodeType,
  sourceNodeName,
  workflowId,
  onSelectNode,
  maxDisplay = 5,
  enabled = true,
}) => {
  const {
    recommendations,
    loading,
    error,
    recordSelection,
    recordDismiss,
  } = useNodeRecommendation({
    sourceNodeId,
    sourceNodeType,
    sourceNodeName,
    workflowId,
    enabled,
    limit: 10,
  });

  // 限制显示数量
  const displayRecommendations = useMemo(() => {
    return recommendations.slice(0, maxDisplay);
  }, [recommendations, maxDisplay]);

  const handleSelectNode = async (nodeType: string) => {
    console.log(
      '[RecommendationSection] handleSelectNode called with nodeType:',
      nodeType,
    );
    console.log(
      '[RecommendationSection] onSelectNode callback exists:',
      !!onSelectNode,
    );
    // 记录用户选择
    await recordSelection(nodeType);
    console.log('[RecommendationSection] recordSelection completed');
    // 触发回调
    onSelectNode?.(nodeType);
    console.log('[RecommendationSection] onSelectNode callback invoked');
  };

  // 加载状态
  if (loading) {
    return (
      <div className={styles.recommendationSection}>
        <div className={styles.header}>
          <h4 className={styles.title}>💡 智能推荐</h4>
        </div>
        <div className={styles.loading}>
          <div className={styles.spinner} />
          <span>正在分析工作流...</span>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className={styles.recommendationSection}>
        <div className={styles.header}>
          <h4 className={styles.title}>💡 智能推荐</h4>
        </div>
        <div className={styles.error}>
          <span className={styles.errorIcon}>⚠️</span>
          <span>推荐功能暂时不可用</span>
        </div>
      </div>
    );
  }

  // 空状态
  if (displayRecommendations.length === 0) {
    return (
      <div className={styles.recommendationSection}>
        <div className={styles.header}>
          <h4 className={styles.title}>💡 智能推荐</h4>
        </div>
        <div className={styles.empty}>
          <span className={styles.emptyIcon}>🤔</span>
          <span>暂无推荐节点</span>
        </div>
      </div>
    );
  }

  // 正常状态
  return (
    <div className={styles.recommendationSection}>
      <div className={styles.header}>
        <h4 className={styles.title}>
          💡 智能推荐
          <span className={styles.badge}>{displayRecommendations.length}</span>
        </h4>
        {recommendations.length > maxDisplay && (
          <span className={styles.moreHint}>
            还有 {recommendations.length - maxDisplay} 个推荐
          </span>
        )}
      </div>

      <div className={styles.list}>
        {displayRecommendations.map((recommendation, index) => (
          <RecommendationCard
            key={`${recommendation.nodeType}-${index}`}
            recommendation={recommendation}
            onClick={() => handleSelectNode(recommendation.nodeType)}
            draggable={false}
          />
        ))}
      </div>

      {displayRecommendations.length > 0 && (
        <div className={styles.footer}>
          <button
            className={styles.dismissButton}
            onClick={recordDismiss}
            type="button"
          >
            不需要推荐
          </button>
        </div>
      )}
    </div>
  );
};
