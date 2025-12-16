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
 * 推荐节点卡片组件
 */

import React, { useCallback } from 'react';
import classNames from 'classnames';
import type { RecommendedNode } from '../../../typing/recommendation';
import styles from './styles.module.less';

export interface RecommendationCardProps {
  /** 推荐节点数据 */
  recommendation: RecommendedNode;
  /** 点击回调 */
  onClick?: () => void;
  /** 是否可拖拽 */
  draggable?: boolean;
  /** 拖拽开始回调 */
  onDragStart?: (event: React.DragEvent) => void;
}

/**
 * 推荐节点卡片
 */
export const RecommendationCard: React.FC<RecommendationCardProps> = ({
  recommendation,
  onClick,
  draggable = true,
  onDragStart,
}) => {
  const scorePercentage = Math.round(recommendation.score * 100);

  const handleClick = useCallback(() => {
    console.log(
      '[RecommendationCard] Card clicked, nodeType:',
      recommendation.nodeType,
    );
    console.log('[RecommendationCard] onClick callback exists:', !!onClick);
    onClick?.();
    console.log('[RecommendationCard] onClick invoked');
  }, [onClick, recommendation.nodeType]);

  const handleDragStart = useCallback(
    (event: React.DragEvent) => {
      onDragStart?.(event);
    },
    [onDragStart]
  );

  // 根据得分确定颜色
  const getScoreColor = () => {
    if (scorePercentage >= 80) return styles.scoreHigh;
    if (scorePercentage >= 60) return styles.scoreMedium;
    return styles.scoreLow;
  };

  return (
    <div
      className={classNames(styles.recommendationCard, {
        [styles.draggable]: draggable,
      })}
      onClick={handleClick}
      draggable={draggable}
      onDragStart={handleDragStart}
      role="button"
      tabIndex={0}
      title={recommendation.reason}
    >
      <div className={styles.header}>
        <div className={styles.title}>
          <span className={styles.icon}>✨</span>
          {recommendation.displayName}
        </div>
        <div className={classNames(styles.score, getScoreColor())}>
          {scorePercentage}%
        </div>
      </div>

      {recommendation.reason && (
        <div className={styles.reason}>{recommendation.reason}</div>
      )}

      {recommendation.category && (
        <div className={styles.footer}>
          <span className={styles.category}>
            {recommendation.category}
          </span>
          {recommendation.strategySource && (
            <span className={styles.source}>
              {recommendation.strategySource === 'rule_engine' && '📋 规则'}
              {recommendation.strategySource === 'llm' && '🤖 AI'}
              {recommendation.strategySource === 'statistics' && '📊 统计'}
            </span>
          )}
        </div>
      )}
    </div>
  );
};
