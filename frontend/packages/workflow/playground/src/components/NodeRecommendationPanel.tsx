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
 * NodeRecommendationPanel - React 组件集成示例
 * 展示如何在实际 UI 中使用 useNodeRecommendation Hook
 */

import React, { useEffect, useState } from 'react';

import type { RecommendedNode } from '../typing/recommendation';
import { useNodeRecommendation } from '../hooks/use-node-recommendation';

import styles from './NodeRecommendationPanel.module.less';

export interface NodeRecommendationPanelProps {
  /** 工作流 ID */
  workflowId: string;
  /** 当前选中的节点 */
  selectedNode?: {
    id: string;
    type: string;
    outputs?: Record<string, unknown>;
  };
  /** 添加节点回调 */
  onAddNode?: (nodeType: string) => void;
  /** 是否显示面板 */
  visible?: boolean;
}

/**
 * 节点推荐面板组件
 *
 * 功能：
 * 1. 自动加载当前节点的推荐
 * 2. 显示推荐列表（带分数和原因）
 * 3. 用户点击推荐时添加节点并提交反馈
 * 4. 支持重新加载
 * 5. 错误处理和加载状态
 *
 * @example
 * ```tsx
 * <NodeRecommendationPanel
 *   workflowId="my-workflow"
 *   selectedNode={{
 *     id: 'llm-node-1',
 *     type: 'LLMNode'
 *   }}
 *   onAddNode={(type) => console.log('Add node:', type)}
 *   visible={true}
 * />
 * ```
 */
export const NodeRecommendationPanel: React.FC<
  NodeRecommendationPanelProps
> = ({ workflowId, selectedNode, onAddNode, visible = true }) => {
  const {
    recommendations,
    loading,
    error,
    loadRecommendations,
    submitFeedback,
    clearRecommendations,
    reload,
  } = useNodeRecommendation({ workflowId });

  const [selectedRecommendation, setSelectedRecommendation] = useState<
    string | null
  >(null);

  // 当选中节点变化时，加载推荐
  useEffect(() => {
    if (selectedNode && visible) {
      loadRecommendations({
        sourceNodeId: selectedNode.id,
        sourceNodeType: selectedNode.type,
        sourceOutputs: selectedNode.outputs,
        limit: 5,
      });
    } else {
      clearRecommendations();
    }
  }, [
    selectedNode?.id,
    selectedNode?.type,
    visible,
    loadRecommendations,
    clearRecommendations,
    selectedNode,
  ]);

  // 处理推荐点击
  const handleRecommendationClick = async (recommendation: RecommendedNode) => {
    console.log('🎯 点击推荐卡片:', recommendation.nodeType);

    try {
      setSelectedRecommendation(recommendation.nodeType);

      // 先添加节点
      if (onAddNode) {
        console.log('📝 调用 onAddNode 回调...');
        await onAddNode(recommendation.nodeType);
        console.log('✅ onAddNode 回调执行完成');
      } else {
        console.warn('⚠️ onAddNode 回调未定义');
      }

      // 提交反馈
      console.log('📤 提交反馈...');
      await submitFeedback(recommendation.nodeType, 'selected');
      console.log('✅ 反馈提交成功');

      setSelectedRecommendation(null);
    } catch (err) {
      console.error('❌ 处理推荐点击失败:', err);
      setSelectedRecommendation(null);
    }
  };

  // 处理拒绝推荐
  const handleDismiss = async (recommendation: RecommendedNode) => {
    console.log('🚫 点击关闭按钮:', recommendation.nodeType);

    try {
      await submitFeedback(recommendation.nodeType, 'dismissed');
      console.log('✅ 已拒绝推荐:', recommendation.nodeType);
    } catch (err) {
      console.error('❌ 拒绝推荐失败:', err);
    }
  };

  // 处理重新加载
  const handleReload = () => {
    reload();
  };

  if (!visible) {
    return null;
  }

  return (
    <div className={styles['node-recommendation-panel']}>
      {/* 标题栏 */}
      <div className={styles['panel-header']}>
        <h3 className={styles['panel-title']}>
          <span className={styles.icon}>💡</span>
          推荐的下一步节点
        </h3>
        {recommendations.length > 0 && (
          <button
            className={styles['reload-button']}
            onClick={handleReload}
            disabled={loading}
            title="重新加载推荐"
          >
            🔄
          </button>
        )}
      </div>

      {/* 加载状态 */}
      {loading ? (
        <div className={styles['loading-state']}>
          <div className={styles.spinner}></div>
          <p>正在获取推荐...</p>
        </div>
      ) : null}

      {/* 错误状态 */}
      {error ? (
        <div className={styles['error-state']}>
          <p className={styles['error-message']}>❌ {error.message}</p>
          <button className={styles['retry-button']} onClick={handleReload}>
            重试
          </button>
        </div>
      ) : null}

      {/* 推荐列表 */}
      {!loading && !error && recommendations.length > 0 && (
        <div className={styles['recommendations-list']}>
          {recommendations.map((recommendation, index) => (
            <RecommendationCard
              key={`${recommendation.nodeType}-${index}`}
              recommendation={recommendation}
              rank={index + 1}
              selected={selectedRecommendation === recommendation.nodeType}
              onSelect={() => handleRecommendationClick(recommendation)}
              onDismiss={() => handleDismiss(recommendation)}
            />
          ))}
        </div>
      )}

      {/* 空状态 */}
      {!loading && !error && recommendations.length === 0 && selectedNode ? (
        <div className={styles['empty-state']}>
          <p>暂无推荐节点</p>
          <p className={styles['empty-hint']}>
            当前节点类型: {selectedNode.type}
          </p>
        </div>
      ) : null}

      {/* 未选中节点提示 */}
      {!selectedNode && (
        <div className={styles['empty-state']}>
          <p>👈 请先选择一个节点</p>
          <p className={styles['empty-hint']}>
            选中节点后，系统会自动推荐下一步可以添加的节点
          </p>
        </div>
      )}
    </div>
  );
};

/**
 * 推荐卡片组件
 */
interface RecommendationCardProps {
  recommendation: RecommendedNode;
  rank: number;
  selected: boolean;
  onSelect: () => void;
  onDismiss: () => void;
}

const RecommendationCard: React.FC<RecommendationCardProps> = ({
  recommendation,
  rank,
  selected,
  onSelect,
  onDismiss,
}) => {
  const SCORE_PERCENT_MULTIPLIER = 100;
  const scorePercent = Math.round(
    recommendation.score * SCORE_PERCENT_MULTIPLIER,
  );

  // 根据分数确定颜色
  const SCORE_EXCELLENT = 0.9;
  const SCORE_GOOD = 0.8;
  const SCORE_MEDIUM = 0.6;
  const getScoreColor = (score: number) => {
    if (score >= SCORE_EXCELLENT) {
      return 'excellent';
    }
    if (score >= SCORE_GOOD) {
      return 'good';
    }
    if (score >= SCORE_MEDIUM) {
      return 'medium';
    }
    return 'low';
  };

  return (
    <div
      className={`${styles['recommendation-card']} ${selected ? styles.selected : ''}`}
      onClick={e => {
        console.log('🃏 卡片被点击:', recommendation.nodeType);
        onSelect();
      }}
    >
      {/* 排名徽章 */}
      <div className={styles['rank-badge']}>#{rank}</div>

      {/* 节点信息 */}
      <div className={styles['node-info']}>
        <div className={styles['node-header']}>
          <span className={styles['node-type']}>{recommendation.nodeType}</span>
          <span
            className={`${styles['score-badge']} ${styles[getScoreColor(recommendation.score)]}`}
          >
            {scorePercent}%
          </span>
        </div>

        {/* 推荐原因 */}
        {recommendation.reason ? (
          <p className={styles.reason}>{recommendation.reason}</p>
        ) : null}

        {/* 分类标签 */}
        {recommendation.metadata?.category ? (
          <span className={styles['category-tag']}>
            {recommendation.metadata.category}
          </span>
        ) : null}
      </div>

      {/* 操作按钮 */}
      <div className={styles['card-actions']}>
        <button
          className={styles['dismiss-button']}
          onClick={e => {
            console.log('🔴 关闭按钮被点击');
            e.stopPropagation();
            e.preventDefault();
            onDismiss();
          }}
          title="不感兴趣"
        >
          ✕
        </button>
      </div>

      {/* 加载指示器 */}
      {selected ? (
        <div className={styles['card-loading']}>
          <div className={styles['spinner-small']}></div>
        </div>
      ) : null}
    </div>
  );
};

export default NodeRecommendationPanel;
