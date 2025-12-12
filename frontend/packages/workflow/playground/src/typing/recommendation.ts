/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * 节点推荐系统 - 类型定义
 */

/**
 * 推荐的节点信息
 */
export interface RecommendedNode {
  /** 节点类型 */
  nodeType: string;
  /** 显示名称 */
  displayName: string;
  /** 推荐得分 (0-1) */
  score: number;
  /** 推荐理由 */
  reason?: string;
  /** 节点分类 */
  category?: string;
  /** 建议的配置 */
  suggestedConfig?: Record<string, any>;
  /** 策略来源 */
  strategySource?: string;
}

/**
 * 推荐元数据
 */
export interface RecommendationMetadata {
  /** 请求ID */
  requestId: string;
  /** 候选总数 */
  totalCandidates: number;
  /** 执行时间(毫秒) */
  executionTimeMs: number;
  /** 使用的策略 */
  strategiesUsed?: string[];
}

/**
 * 推荐响应
 */
export interface RecommendationResponse {
  /** 推荐结果列表 */
  recommendations: RecommendedNode[];
  /** 元数据 */
  metadata: RecommendationMetadata;
}

/**
 * 推荐请求参数
 */
export interface RecommendationRequest {
  /** 工作流ID */
  workflowId: string;
  /** 源节点ID */
  sourceNodeId: string;
  /** 源节点类型 */
  sourceNodeType: string;
  /** 源节点名称 */
  sourceNodeName?: string;
  /** 输出格式 */
  outputFormat?: string;
  /** 推荐数量限制 */
  limit?: number;
  /** 是否包含推荐理由 */
  includeReason?: boolean;
  /** 工作流上下文 */
  workflowContext?: WorkflowContext;
}

/**
 * 工作流上下文
 */
export interface WorkflowContext {
  /** 工作流名称 */
  workflowName?: string;
  /** 工作流目标 */
  workflowObjective?: string;
  /** 节点总数 */
  totalNodes?: number;
  /** 是否有数据库节点 */
  hasDatabaseNode?: boolean;
  /** 是否有LLM节点 */
  hasLLMNode?: boolean;
  /** 是否有循环节点 */
  hasLoopNode?: boolean;
  /** 是否在循环内部 */
  insideLoop?: boolean;
  /** 是否聊天工作流 */
  isChatWorkflow?: boolean;
}

/**
 * 反馈请求参数
 */
export interface FeedbackRequest {
  /** 请求ID */
  requestId: string;
  /** 选择的节点类型 */
  selectedNode?: string;
  /** 用户操作 */
  userAction: 'selected' | 'dismissed' | 'ignored';
  /** 用户评论 */
  userComment?: string;
}

/**
 * 反馈响应
 */
export interface FeedbackResponse {
  /** 是否成功 */
  success: boolean;
  /** 消息 */
  message?: string;
}
