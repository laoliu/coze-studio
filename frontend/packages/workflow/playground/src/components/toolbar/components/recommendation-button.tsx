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
import { IconCozLightbulb } from '@coze-arch/coze-design/icons';
import { Tooltip, IconButton } from '@coze-arch/coze-design';

import { WorkflowEditService } from '@/services';
import { useFloatLayoutService, useGlobalState } from '@/hooks';
import { LayoutPanelKey } from '@/constants';

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

  const handleClick = useCallback(() => {
    // 获取当前选中的节点
    const selectedNode = selectService.activatedNode;

    // 添加节点的回调函数
    const handleAddNode = async (nodeType: string) => {
      try {
        // 获取选中节点的位置，在其下方添加新节点
        let position = { clientX: 0, clientY: 0 };
        const OFFSET_Y = 150;
        if (selectedNode) {
          const nodePosition = selectedNode.position;
          // 在选中节点下方偏移
          position = {
            clientX: nodePosition.x,
            clientY: nodePosition.y + OFFSET_Y,
          };
        } else {
          // 如果没有选中节点，在画布中心添加
          const { viewport } = playgroundConfig;
          const HALF = 2;
          position = {
            clientX: viewport.x + viewport.width / HALF,
            clientY: viewport.y + viewport.height / HALF,
          };
        }

        // 调用 editService 添加节点
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        await editService.addNode(nodeType as any, undefined, position, false);

        console.log(`✅ 成功添加节点: ${nodeType}`);
      } catch (error) {
        console.error('添加节点失败:', error);
        throw error;
      }
    };

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
  }, [
    floatLayoutService,
    selectService,
    editService,
    playgroundConfig,
    workflowId,
  ]);

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
