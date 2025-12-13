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

import { useCallback } from 'react';

import { useService } from '@flowgram-adapter/free-layout-editor';
import {
  WorkflowSelectService,
  PlaygroundConfigEntity,
} from '@flowgram-adapter/free-layout-editor';
import { StandardNodeType } from '@coze-workflow/base';
import { IconCozLightbulb } from '@coze-arch/coze-design/icons';
import { Tooltip, IconButton } from '@coze-arch/coze-design';

import { WorkflowEditService } from '@/services';
import { useFloatLayoutService, useGlobalState } from '@/hooks';
import { LayoutPanelKey } from '@/constants';

// 节点类型名称到 StandardNodeType 的映射
// 后端返回的节点类型名称（如 CodeRunner）需要映射到前端的数字 ID
const NODE_TYPE_MAP: Record<string, string> = {
  // 基础节点
  Entry: '1', // StandardNodeType.Start
  Exit: '2', // StandardNodeType.End - 注意：rules.yaml 中使用 End，但后端常量是 Exit
  End: '2', // 别名，兼容 rules.yaml

  // AI 处理
  LLM: '3', // StandardNodeType.LLM

  // 外部服务
  Plugin: '4', // StandardNodeType.Api

  // 数据处理
  CodeRunner: '5', // StandardNodeType.Code

  // 数据检索
  Dataset: '6', // 知识库检索 StandardNodeType.Dataset
  KnowledgeRetriever: '6', // 别名

  // 流程控制
  If: '8', // 条件判断 StandardNodeType.If
  Selector: '8', // 别名

  SubWorkflow: '9', // StandardNodeType.SubWorkflow
  Lambda: '9', // 别名

  Variable: '11', // StandardNodeType.Variable

  // 数据库
  Database: '12', // SQL自定义 StandardNodeType.Database
  DatabaseCustomSQL: '12', // 别名

  Output: '13', // StandardNodeType.Output
  OutputEmitter: '13', // 别名

  Text: '15', // StandardNodeType.Text
  TextProcessor: '15', // 别名

  Question: '18', // StandardNodeType.Question
  QuestionAnswer: '18', // 别名

  Break: '19', // StandardNodeType.Break

  SetVariable: '20', // StandardNodeType.SetVariable
  VariableAssignerWithinLoop: '20', // 别名

  Loop: '21', // StandardNodeType.Loop

  Intent: '22', // StandardNodeType.Intent
  IntentDetector: '22', // 别名

  DatasetWrite: '27', // 知识库写入 StandardNodeType.DatasetWrite
  KnowledgeIndexer: '27', // 别名

  Batch: '28', // StandardNodeType.Batch

  Continue: '29', // StandardNodeType.Continue

  Input: '30', // StandardNodeType.Input
  InputReceiver: '30', // 别名

  // HTTP 请求
  Http: '45', // StandardNodeType.Http
  HTTPRequester: '45', // 别名

  // 数据库 CRUD
  DatabaseUpdate: '42', // StandardNodeType.DatabaseUpdate
  DatabaseQuery: '43', // StandardNodeType.DatabaseQuery
  DatabaseDelete: '44', // StandardNodeType.DatabaseDelete
  DatabaseCreate: '46', // StandardNodeType.DatabaseCreate
  DatabaseInsert: '46', // 别名
};

/**
 * 节点推荐按钮
 * 点击后在右侧打开推荐面板，显示智能推荐的下一个节点
 */
export const RecommendationButton = () => {
  const floatLayoutService = useFloatLayoutService();
  const selectService = useService(WorkflowSelectService);
  const editService = useService<WorkflowEditService>(WorkflowEditService);
  const playgroundConfig = useService(PlaygroundConfigEntity);
  const { workflowId } = useGlobalState();

  // 添加节点的回调函数 - 独立定义，避免闭包问题
  const handleAddNode = useCallback(
    async (nodeType: string) => {
      try {
        // 将推荐的节点类型名称转换为 StandardNodeType
        const standardNodeType = NODE_TYPE_MAP[nodeType];

        if (!standardNodeType) {
          console.error(
            `❌ 未知的节点类型: ${nodeType}，可用类型:`,
            Object.keys(NODE_TYPE_MAP),
          );
          throw new Error(`不支持的节点类型: ${nodeType}`);
        }

        console.log(`🔄 节点类型转换: ${nodeType} → ${standardNodeType}`);

        // 特殊节点类型处理：这些节点需要用户额外选择/配置，不能直接创建
        if (
          standardNodeType === StandardNodeType.SubWorkflow ||
          standardNodeType === StandardNodeType.Api ||
          standardNodeType === StandardNodeType.Imageflow
        ) {
          console.warn(
            `⚠️ 节点类型 ${nodeType} 需要额外配置，当前暂不支持直接从推荐添加`,
          );
          // TODO: 未来可以打开对应的选择模态框
          // 例如: openWorkflowModal() 或 openPluginModal()
          alert(
            `${nodeType} 节点需要选择具体的${
              standardNodeType === StandardNodeType.SubWorkflow
                ? '子工作流'
                : standardNodeType === StandardNodeType.Api
                  ? '插件'
                  : '图像流'
            }，请从左侧节点面板手动添加`,
          );
          return;
        }

        // 获取当前选中的节点（实时获取，不使用闭包）
        const currentSelectedNode = selectService.activatedNode;

        // 获取节点位置
        let position = { clientX: 0, clientY: 0 };
        const OFFSET_Y = 150;
        const HALF = 2;

        if (
          currentSelectedNode?.position?.x !== undefined &&
          currentSelectedNode?.position?.y !== undefined
        ) {
          // 在选中节点下方偏移
          position = {
            clientX: currentSelectedNode.position.x,
            clientY: currentSelectedNode.position.y + OFFSET_Y,
          };
          console.log(
            '📍 在选中节点下方添加，节点位置:',
            currentSelectedNode.position,
          );
        } else {
          // 如果没有选中节点或节点没有位置信息，在画布中心添加
          const { viewport } = playgroundConfig;
          if (viewport) {
            position = {
              clientX: viewport.x + viewport.width / HALF,
              clientY: viewport.y + viewport.height / HALF,
            };
            console.log('📍 在画布中心添加，viewport:', viewport);
          } else {
            // 如果viewport也没有，使用一个默认位置
            const DEFAULT_POS = 200;
            position = {
              clientX: DEFAULT_POS,
              clientY: DEFAULT_POS,
            };
            console.log('📍 使用默认位置');
          }
        }

        console.log(
          `🔧 准备添加节点: ${nodeType} (${standardNodeType})，位置:`,
          position,
        );

        // 调用 editService 添加节点
        console.log('🔨 调用 editService.addNode...');
        console.log('  - 原始类型:', nodeType);
        console.log('  - 转换后类型:', standardNodeType);
        console.log('  - 转换后类型的类型:', typeof standardNodeType);
        console.log('  - 转换后类型的值:', JSON.stringify(standardNodeType));
        console.log('  - position:', position);
        console.log('  - StandardNodeType枚举:', StandardNodeType);
        console.log('  - 映射表:', NODE_TYPE_MAP);

        let addedNode;
        try {
          console.log('⏳ 正在调用 addNode 方法...');
          console.log(`  - 节点类型: ${standardNodeType} (${nodeType})`);
          console.log('  - 位置:', position);

          // 明确将字符串转换为 StandardNodeType
          // 对于普通节点，传递 undefined 作为 nodeJson，系统会使用默认配置
          addedNode = await editService.addNode(
            standardNodeType as StandardNodeType,
            undefined, // nodeJson - 让系统使用默认配置
            position,
            false, // isDrag
          );

          console.log('✅ addNode 调用成功，返回值:', addedNode);
        } catch (addNodeError) {
          console.error('💥 addNode 调用失败:', addNodeError);
          console.error('  - 节点类型:', standardNodeType, `(${nodeType})`);
          console.error('  - 错误详情:', addNodeError);
          throw addNodeError;
        }

        if (addedNode) {
          console.log(`🎉 成功添加节点: ${nodeType}`, addedNode);
        } else {
          console.warn(
            `⚠️ 添加节点返回空: ${nodeType}，可能节点类型不支持或需要额外配置`,
          );
        }
      } catch (error) {
        console.error('❌ 添加节点失败:', error);
        throw error;
      }
    },
    [selectService, editService, playgroundConfig],
  );

  const handleClick = useCallback(() => {
    // 获取当前选中的节点
    const selectedNode = selectService.activatedNode;

    // 打开推荐面板
    floatLayoutService.open(LayoutPanelKey.NodeRecommendation, 'right', {
      workflowId,
      selectedNode: selectedNode
        ? {
            id: selectedNode.id,
            type: selectedNode.type,
            // 可以传递更多节点信息
          }
        : undefined,
      onAddNode: handleAddNode,
    });
  }, [floatLayoutService, selectService, workflowId, handleAddNode]);

  // 推荐功能在只读模式下也可用，因为它只是提供建议，不会修改工作流
  return (
    <Tooltip content="智能推荐下一个节点">
      <IconButton
        icon={<IconCozLightbulb />}
        color="secondary"
        onClick={handleClick}
        data-testid="workflow-recommendation-button"
      />
    </Tooltip>
  );
};
