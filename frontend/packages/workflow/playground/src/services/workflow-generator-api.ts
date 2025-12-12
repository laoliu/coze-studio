/**
 * 工作流生成器 API 服务
 * Workflow Generator API Service
 */

import { request } from '@coze-arch/bot-http';

/**
 * 生成工作流请求参数
 */
export interface GenerateWorkflowRequest {
  /** 空间 ID */
  space_id: string;
  /** 用户需求描述 */
  user_requirement: string;
  /** 语言 (zh-CN / en-US) */
  language?: string;
  /** 生成约束条件 */
  constraints?: {
    /** 最大节点数 */
    max_nodes?: number;
    /** 允许的节点类型 */
    allowed_node_types?: string[];
    /** 是否需要错误处理 */
    require_error_handling?: boolean;
  };
}

/**
 * 节点位置
 */
export interface NodePosition {
  x: number;
  y: number;
}

/**
 * 工作流节点
 */
export interface WorkflowNode {
  /** 节点 ID */
  id: string;
  /** 节点类型 */
  type: string;
  /** 节点标题 */
  title: string;
  /** 节点描述 */
  description?: string;
  /** 节点配置 */
  config?: Record<string, any>;
  /** 节点位置 */
  position?: NodePosition;
}

/**
 * 工作流连接
 */
export interface WorkflowEdge {
  /** 连接 ID */
  id: string;
  /** 源节点 ID */
  source: string;
  /** 目标节点 ID */
  target: string;
  /** 源节点句柄 */
  source_handle?: string;
  /** 目标节点句柄 */
  target_handle?: string;
}

/**
 * 生成工作流响应
 */
export interface GenerateWorkflowResponse {
  /** 生成的工作流 ID */
  workflow_id: string;
  /** 工作流名称 */
  workflow_name: string;
  /** 节点列表 */
  nodes: WorkflowNode[];
  /** 连接列表 */
  edges: WorkflowEdge[];
  /** 生成说明 */
  explanation: string;
  /** 置信度 (0-1) */
  confidence: number;
}

/**
 * API 响应包装
 */
interface APIResponse<T> {
  code: number;
  msg: string;
  data: T;
}

/**
 * 生成工作流
 * 
 * @param params 生成请求参数
 * @returns 生成的工作流数据
 * 
 * @example
 * ```typescript
 * const workflow = await generateWorkflow({
 *   space_id: '123456',
 *   user_requirement: '我需要一个批量处理图片的工作流',
 *   language: 'zh-CN'
 * });
 * console.log(workflow.nodes); // 生成的节点列表
 * ```
 */
export async function generateWorkflow(
  params: GenerateWorkflowRequest
): Promise<GenerateWorkflowResponse> {
  try {
    const response = await request.post<APIResponse<GenerateWorkflowResponse>>(
      '/api/workflow_api/generate',
      params,
      {
        timeout: 60000, // 60s 超时
      }
    );

    if (response.code !== 0) {
      throw new Error(response.msg || '生成工作流失败');
    }

    return response.data;
  } catch (error: any) {
    console.error('[WorkflowGenerator] Generate failed:', error);
    
    // 提供更友好的错误信息
    if (error.message?.includes('timeout')) {
      throw new Error('生成超时，请稍后重试');
    } else if (error.message?.includes('network')) {
      throw new Error('网络错误，请检查网络连接');
    } else {
      throw new Error(error.message || '生成工作流失败，请重试');
    }
  }
}

/**
 * 验证生成请求参数
 * 
 * @param params 请求参数
 * @returns 验证结果
 */
export function validateGenerateRequest(
  params: GenerateWorkflowRequest
): { valid: boolean; error?: string } {
  if (!params.space_id || params.space_id.trim() === '') {
    return { valid: false, error: '空间 ID 不能为空' };
  }

  if (!params.user_requirement || params.user_requirement.trim() === '') {
    return { valid: false, error: '需求描述不能为空' };
  }

  if (params.user_requirement.length < 10) {
    return { valid: false, error: '需求描述至少需要 10 个字符' };
  }

  if (params.user_requirement.length > 2000) {
    return { valid: false, error: '需求描述不能超过 2000 个字符' };
  }

  return { valid: true };
}

/**
 * 获取示例需求列表
 */
export function getExampleRequirements(language: string = 'zh-CN'): string[] {
  if (language === 'zh-CN') {
    return [
      '我需要一个批量处理图片的工作流，能够压缩图片并添加水印',
      '创建一个内容生成工作流，从知识库中检索信息，然后生成文章',
      '构建一个数据分析工作流，读取 Excel 文件，进行数据清洗和统计',
      '设计一个 API 集成工作流，调用第三方天气 API 并格式化数据',
      '实现一个自动化审批工作流，根据条件判断是否需要人工审核',
      '开发一个多语言翻译工作流，支持文本的批量翻译',
    ];
  } else {
    return [
      'I need a batch image processing workflow that can compress images and add watermarks',
      'Create a content generation workflow that retrieves information from a knowledge base and generates articles',
      'Build a data analysis workflow that reads Excel files, cleans data, and performs statistics',
      'Design an API integration workflow that calls a third-party weather API and formats the data',
      'Implement an automated approval workflow that determines whether manual review is needed based on conditions',
      'Develop a multilingual translation workflow that supports batch text translation',
    ];
  }
}

/**
 * 格式化置信度显示
 */
export function formatConfidence(confidence: number): string {
  const percentage = (confidence * 100).toFixed(1);
  
  if (confidence >= 0.9) {
    return `${percentage}% (高)`;
  } else if (confidence >= 0.7) {
    return `${percentage}% (中)`;
  } else {
    return `${percentage}% (低)`;
  }
}

/**
 * 获取节点类型的中文名称
 */
export function getNodeTypeName(type: string, language: string = 'zh-CN'): string {
  const nameMap: Record<string, Record<string, string>> = {
    'zh-CN': {
      Start: '开始',
      End: '结束',
      LLM: '大语言模型',
      HTTP: 'HTTP 请求',
      Code: '代码执行',
      Condition: '条件判断',
      Loop: '循环',
      Knowledge: '知识库',
      Plugin: '插件',
      DataTransform: '数据转换',
    },
    'en-US': {
      Start: 'Start',
      End: 'End',
      LLM: 'Large Language Model',
      HTTP: 'HTTP Request',
      Code: 'Code Execution',
      Condition: 'Condition',
      Loop: 'Loop',
      Knowledge: 'Knowledge Base',
      Plugin: 'Plugin',
      DataTransform: 'Data Transform',
    },
  };

  return nameMap[language]?.[type] || type;
}
