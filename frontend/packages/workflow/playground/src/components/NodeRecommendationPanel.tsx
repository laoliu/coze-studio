/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * NodeRecommendationPanel - React 组件集成示例
 * 展示如何在实际 UI 中使用 useNodeRecommendation Hook
 */

import React, { useEffect, useState } from 'react';
import { useNodeRecommendation } from '../hooks/use-node-recommendation';
import type { RecommendedNode } from '../typing/recommendation';
import './NodeRecommendationPanel.module.less';

export interface NodeRecommendationPanelProps {
  /** 工作流 ID */
  workflowId: string;
  /** 当前选中的节点 */
  selectedNode?: {
    id: string;
    type: string;
    outputs?: Record<string, any>;
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
export const NodeRecommendationPanel: React.FC<NodeRecommendationPanelProps> = ({
  workflowId,
  selectedNode,
  onAddNode,
  visible = true,
}) => {
  const {
    recommendations,
    loading,
    error,
    loadRecommendations,
    submitFeedback,
    clearRecommendations,
    reload,
  } = useNodeRecommendation({ workflowId });

  const [selectedRecommendation, setSelectedRecommendation] = useState<string | null>(null);

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
  }, [selectedNode?.id, selectedNode?.type, visible]);

  // 处理推荐点击
  const handleRecommendationClick = async (recommendation: RecommendedNode) => {
    try {
      setSelectedRecommendation(recommendation.nodeType);
      
      // 提交反馈
      await submitFeedback(recommendation.nodeType, 'selected');
      
      // 调用添加节点回调
      onAddNode?.(recommendation.nodeType);
      
      // 提示成功
      console.log(`✅ 已添加节点: ${recommendation.nodeType}`);
      
      setSelectedRecommendation(null);
    } catch (err) {
      console.error('Failed to add node:', err);
      setSelectedRecommendation(null);
    }
  };

  // 处理拒绝推荐
  const handleDismiss = async (recommendation: RecommendedNode) => {
    await submitFeedback(recommendation.nodeType, 'dismissed');
    console.log(`❌ 已拒绝推荐: ${recommendation.nodeType}`);
  };

  // 处理重新加载
  const handleReload = () => {
    reload();
  };

  if (!visible) {
    return null;
  }

  return (
    <div className="node-recommendation-panel">
      {/* 标题栏 */}
      <div className="panel-header">
        <h3 className="panel-title">
          <span className="icon">💡</span>
          推荐的下一步节点
        </h3>
        {recommendations.length > 0 && (
          <button 
            className="reload-button"
            onClick={handleReload}
            disabled={loading}
            title="重新加载推荐"
          >
            🔄
          </button>
        )}
      </div>

      {/* 加载状态 */}
      {loading && (
        <div className="loading-state">
          <div className="spinner"></div>
          <p>正在获取推荐...</p>
        </div>
      )}

      {/* 错误状态 */}
      {error && (
        <div className="error-state">
          <p className="error-message">❌ {error.message}</p>
          <button className="retry-button" onClick={handleReload}>
            重试
          </button>
        </div>
      )}

      {/* 推荐列表 */}
      {!loading && !error && recommendations.length > 0 && (
        <div className="recommendations-list">
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
      {!loading && !error && recommendations.length === 0 && selectedNode && (
        <div className="empty-state">
          <p>暂无推荐节点</p>
          <p className="empty-hint">
            当前节点类型: {selectedNode.type}
          </p>
        </div>
      )}

      {/* 未选中节点提示 */}
      {!selectedNode && (
        <div className="empty-state">
          <p>👈 请先选择一个节点</p>
          <p className="empty-hint">
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
  const scorePercent = Math.round(recommendation.score * 100);
  
  // 根据分数确定颜色
  const getScoreColor = (score: number) => {
    if (score >= 0.9) return 'excellent';
    if (score >= 0.8) return 'good';
    if (score >= 0.6) return 'medium';
    return 'low';
  };

  return (
    <div 
      className={`recommendation-card ${selected ? 'selected' : ''}`}
      onClick={onSelect}
    >
      {/* 排名徽章 */}
      <div className="rank-badge">#{rank}</div>

      {/* 节点信息 */}
      <div className="node-info">
        <div className="node-header">
          <span className="node-type">{recommendation.nodeType}</span>
          <span className={`score-badge ${getScoreColor(recommendation.score)}`}>
            {scorePercent}%
          </span>
        </div>

        {/* 推荐原因 */}
        {recommendation.reason && (
          <p className="reason">{recommendation.reason}</p>
        )}

        {/* 分类标签 */}
        {recommendation.metadata?.category && (
          <span className="category-tag">
            {recommendation.metadata.category}
          </span>
        )}
      </div>

      {/* 操作按钮 */}
      <div className="card-actions">
        <button
          className="dismiss-button"
          onClick={(e) => {
            e.stopPropagation();
            onDismiss();
          }}
          title="不感兴趣"
        >
          ✕
        </button>
      </div>

      {/* 加载指示器 */}
      {selected && (
        <div className="card-loading">
          <div className="spinner-small"></div>
        </div>
      )}
    </div>
  );
};

export default NodeRecommendationPanel;
