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

import { Tooltip, IconButton } from '@coze-arch/coze-design';
import { IconCozLightbulb } from '@coze-arch/coze-design/icons';
import { useService } from '@flowgram-adapter/free-layout-editor';
import { WorkflowSelectService } from '@flowgram-adapter/free-layout-editor';

import { useFloatLayoutService, useGlobalState } from '@/hooks';
import { LayoutPanelKey } from '@/constants';

/**
 * 节点推荐按钮
 * 点击后在右侧打开推荐面板，显示智能推荐的下一个节点
 */
export const RecommendationButton = () => {
  const floatLayoutService = useFloatLayoutService();
  const selectService = useService(WorkflowSelectService);
  const { workflowId } = useGlobalState();

  const handleClick = useCallback(() => {
    // 获取当前选中的节点
    const selectedNode = selectService.activatedNode;

    // 打开推荐面板
    floatLayoutService.open(LayoutPanelKey.NodeRecommendation, 'right', {
      workflowId: workflowId,
      selectedNode: selectedNode
        ? {
            id: selectedNode.id,
            type: selectedNode.type,
            // 可以传递更多节点信息
          }
        : undefined,
    });
  }, [floatLayoutService, selectService, workflowId]);

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
