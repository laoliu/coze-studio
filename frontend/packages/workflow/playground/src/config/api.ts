/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * API 配置 - 不同环境的 API 地址配置
 */

// 判断是否为开发环境
const isDevelopment = process.env.NODE_ENV === 'development';

// 判断是否启用推荐服务测试服务器
const useTestServer = process.env.USE_RECOMMENDATION_TEST_SERVER === 'true' || isDevelopment;

/**
 * API 基础路径配置
 */
export const API_CONFIG = {
  /**
   * 节点推荐服务 API 基础路径
   * 
   * 开发环境:
   * - 如果启用测试服务器: /api (代理到 localhost:8080)
   * - 否则: /api/workflow_api/node (代理到主后端服务)
   * 
   * 生产环境: /api/workflow_api/node
   */
  RECOMMENDATION: useTestServer 
    ? '/api'  // 测试服务器路径 (通过代理转发到 localhost:8080)
    : '/api/workflow_api/node',  // 生产服务路径

  /**
   * 其他 API 路径...
   */
};

/**
 * 环境信息（用于调试）
 */
export const ENV_INFO = {
  isDevelopment,
  useTestServer,
  recommendationBaseURL: API_CONFIG.RECOMMENDATION,
};

// 开发环境下打印配置信息
if (isDevelopment) {
  console.log('[API Config]', ENV_INFO);
}
