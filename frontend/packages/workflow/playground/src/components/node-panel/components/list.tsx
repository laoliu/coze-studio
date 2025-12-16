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

import {
  forwardRef,
  type MouseEvent as MouseEventType,
  type RefObject,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
  useCallback,
} from 'react';

import { throttle } from 'lodash-es';
import classNames from 'classnames';
import {
  Root,
  ScrollAreaCorner,
  ScrollAreaScrollbar,
  ScrollAreaThumb,
  ScrollAreaViewport,
} from '@radix-ui/react-scroll-area';
import { I18n } from '@coze-arch/i18n';
import { IconCozEmpty, IconCozMagnifier } from '@coze-arch/coze-design/icons';
import { EmptyState, Input } from '@coze-arch/coze-design';
import {
  type WorkflowNodeEntity,
  useService,
  type WorkflowLinePortInfo,
} from '@flowgram-adapter/free-layout-editor';

import { type UnionNodeTemplate } from '@/typing';
import { WorkflowPlaygroundContext } from '@/workflow-playground-context';
import type { WorkflowLineService } from '@/services/workflow-line-service';

import { NodePanelContextProvider } from '../hooks/node-panel-context';
import { RecommendationSection } from './recommendation-section';
import { useSearchNode, useTemplateNodeList } from '../hooks';
import { PANEL_WIDTH, THROTTLE_INTERVAL } from '../constant';
import { SearchResultNodeList } from './search-result-node-list';
import {
  FavoritePluginNodeList,
  type FavoritePluginNodeListRefType,
} from './plugin-node/favorite-plugin-node-list';
import { AtomCategoryList } from './atom-category-list';

import styles from './styles.module.less';
export type NodeListRefType = FavoritePluginNodeListRefType;
interface NodesContainerProps {
  onSelect: (props: {
    event: MouseEventType<HTMLElement>;
    nodeTemplate: UnionNodeTemplate;
  }) => void;
  enableDrag?: boolean;
  containerNode?: WorkflowNodeEntity;
  adaptiveHeight?: number;
  /**
   * Update the status of the node being added. At this time, clickOutside will not close the node panel to avoid triggering onClose to close before the addition of the node is completed, and the previous connection cannot be destroyed.
   * @param isAdding
   * @returns
   */
  onAddingNode?: (isAdding: boolean) => void;
}
export const NodeList = forwardRef<NodeListRefType, NodesContainerProps>(
  (props, ref) => {
    const {
      onSelect,
      enableDrag = false,
      containerNode,
      adaptiveHeight,
      onAddingNode,
    } = props;

    // 获取工作流上下文（用于推荐）
    const context = useService<WorkflowPlaygroundContext>(
      WorkflowPlaygroundContext,
    );

    const nodeCategoryList = useTemplateNodeList(containerNode);
    const nodeListRef = useRef<HTMLDivElement>();
    const [showBorder, setShowBorder] = useState(false);
    const [input, setInput] = useState('');

    // 推荐相关状态
    const [showRecommendation, setShowRecommendation] = useState(false);
    const [sourceNodeInfo, setSourceNodeInfo] = useState<{
      nodeId?: string;
      nodeType?: string;
      nodeName?: string;
    }>({});

    const {
      showSearchResult,
      searchResult,
      isSearching,
      noSearchResult,
      loadMore,
      keyword,
      handleKeywordChange,
    } = useSearchNode({
      atomNodeCategoryList: nodeCategoryList,
    });

    useEffect(() => {
      if (!nodeListRef.current || showSearchResult) {
        setShowBorder(false);
        return;
      }
      const el = nodeListRef.current;
      const handleScroll = throttle(() => {
        if ((el.scrollTop ?? 0) > 0) {
          setShowBorder(true);
        } else {
          setShowBorder(false);
        }
      }, THROTTLE_INTERVAL);
      setShowBorder((nodeListRef.current?.scrollTop ?? 0) > 0);
      el.addEventListener('scroll', handleScroll);
      return () => el.removeEventListener('scroll', handleScroll);
    }, [showSearchResult]);
    const favoritePluginsRef = useRef<FavoritePluginNodeListRefType>();

    useImperativeHandle(ref, () => ({
      refetch: async () => {
        await favoritePluginsRef.current?.refetch();
      },
    }));

    useEffect(() => {
      const text = input.replaceAll(' ', ''); // Remove spaces
      handleKeywordChange(text);
    }, [input]);

    // 处理推荐节点选择
    const handleSelectRecommendation = useCallback(
      async (nodeType: string) => {
        console.log(
          '[Recommendation] handleSelectRecommendation called with nodeType:',
          nodeType,
        );
        setShowRecommendation(false);

        if (!containerNode) {
          console.warn('[Recommendation] No container node available');
          return;
        }

        console.log(
          '[Recommendation] Container node:',
          containerNode.id,
          containerNode.type,
        );

        // 创建一个简单的 mock 事件
        const mockEvent: MouseEventType = {
          type: 'click',
          clientX: 0,
          clientY: 0,
        };

        // 调用原始的 onSelect 来添加节点
        const nodeTemplate: UnionNodeTemplate = {
          type: nodeType,
        };

        console.log('[Recommendation] Calling onSelect to add node...');
        await onSelect({
          event: mockEvent,
          nodeTemplate,
        });
        console.log('[Recommendation] onSelect completed');

        // 节点添加后，我们需要创建连接
        // 延迟一下确保节点已经被添加到画布
        const CONNECTION_DELAY_MS = 100;
        console.log(
          '[Recommendation] Setting up setTimeout for connection creation',
        );
        setTimeout(() => {
          console.log('[Recommendation] setTimeout callback executing');
          try {
            // 获取 workflow-line-service 来创建连接
            const lineService = context.container.get<WorkflowLineService>(
              'WorkflowLineService',
            );
            console.log('[Recommendation] Got lineService:', !!lineService);

            // 获取刚添加的节点（应该是最新添加的）
            const allNodes = context.entityManager
              ?.getAllEntities()
              .filter((e: WorkflowNodeEntity) => e.type === nodeType);

            console.log(
              '[Recommendation] All nodes of type',
              nodeType,
              ':',
              allNodes?.length,
            );

            if (!allNodes || allNodes.length === 0) {
              console.warn('[Recommendation] Could not find newly added node');
              return;
            }

            // 取最后一个（最新添加的）
            const newNode = allNodes[allNodes.length - 1];
            console.log('[Recommendation] New node:', newNode.id, newNode.type);

            // 获取 containerNode 的输出端口
            const sourceOutputPorts = containerNode.ports.filter(
              p => p.direction === 'output',
            );
            console.log(
              '[Recommendation] Source output ports:',
              sourceOutputPorts.length,
            );

            // 获取新节点的输入端口
            const targetInputPorts = newNode.ports.filter(
              p => p.direction === 'input',
            );
            console.log(
              '[Recommendation] Target input ports:',
              targetInputPorts.length,
            );

            if (
              sourceOutputPorts.length === 0 ||
              targetInputPorts.length === 0
            ) {
              console.warn('[Recommendation] No valid ports for connection');
              return;
            }

            // 创建连接：从 containerNode 的第一个输出端口到新节点的第一个输入端口
            const lineInfo: WorkflowLinePortInfo = {
              from: containerNode.id,
              fromPort: sourceOutputPorts[0].id,
              to: newNode.id,
              toPort: targetInputPorts[0].id,
            };
            console.log('[Recommendation] Creating line with info:', lineInfo);
            lineService?.createLine(lineInfo);

            console.log('[Recommendation] Successfully created connection');
          } catch (error) {
            console.error(
              '[Recommendation] Failed to create connection:',
              error,
            );
          }
        }, CONNECTION_DELAY_MS); // 延迟100ms确保节点已添加
      },
      [onSelect, containerNode, context],
    );

    // 当 containerNode 变化时，更新推荐源节点
    useEffect(() => {
      if (containerNode) {
        setSourceNodeInfo({
          nodeId: containerNode.id,
          nodeType: containerNode.type,
          nodeName: containerNode.data?.nodeMeta?.title,
        });
        setShowRecommendation(true);
      } else {
        setShowRecommendation(false);
      }
    }, [containerNode]);

    return (
      <div
        className={styles['node-panel']}
        style={{ height: adaptiveHeight, width: PANEL_WIDTH }}
        data-flow-editor-selectable="false"
        data-testid="workflow.detail.node-panel"
      >
        <svg style={{ width: 0, height: 0, display: 'block' }}>
          <clipPath
            id="favorite-plugin-clip-path"
            clipPathUnits="objectBoundingBox"
          >
            <path d="M0,1 H1 V0 C1,0.552,0.552,1,0,1"></path>
          </clipPath>
        </svg>
        <div
          className={classNames(
            styles['node-search'],
            showBorder ? styles['node-search-shadow'] : '',
          )}
          data-testid="workflow.detail.node-panel.search"
        >
          <Input
            className={styles['node-search-input']}
            showClear
            placeholder={I18n.t('workflow_250306_01')}
            value={input}
            onClear={() => setInput('')}
            onChange={setInput}
            prefix={<IconCozMagnifier style={{ fontSize: '16px' }} />}
          />
        </div>
        {/* 推荐区域 */}
        {!showSearchResult && showRecommendation && sourceNodeInfo.nodeId && (
          <div className={styles['recommendation-wrapper']}>
            <RecommendationSection
              sourceNodeId={sourceNodeInfo.nodeId}
              sourceNodeType={sourceNodeInfo.nodeType || ''}
              sourceNodeName={sourceNodeInfo.nodeName}
              workflowId={context.workflowEntity?.id || ''}
              onSelectNode={handleSelectRecommendation}
              maxDisplay={3}
              enabled={true}
            />
          </div>
        )}
        {noSearchResult ? (
          <EmptyState
            className="mx-auto mt-[128px] max-w-[270px]"
            icon={<IconCozEmpty />}
            size="large"
            title={I18n.t('workflow_250305_001')}
            description={I18n.t('workflow_250305_002')}
          />
        ) : null}
        <Root className={classNames('node-panel-render', styles['node-list'])}>
          <ScrollAreaViewport
            ref={nodeListRef as RefObject<HTMLDivElement>}
            className={styles.viewport}
          >
            <NodePanelContextProvider
              value={{
                onSelect,
                enableDrag,
                keyword,
                getScrollContainer: () => nodeListRef.current,
                onLoadMore: loadMore,
                onAddingNode,
              }}
            >
              {showSearchResult ? (
                <SearchResultNodeList
                  searchResult={searchResult}
                  loading={isSearching}
                />
              ) : (
                <div
                  className={styles['list-wrapper']}
                  data-testid="workflow.detail.node-panel.list"
                >
                  <AtomCategoryList data={[nodeCategoryList[0]]} />
                  <FavoritePluginNodeList
                    ref={
                      favoritePluginsRef as RefObject<FavoritePluginNodeListRefType>
                    }
                  />
                  <AtomCategoryList data={nodeCategoryList.slice(1)} />
                </div>
              )}
            </NodePanelContextProvider>
          </ScrollAreaViewport>
          <ScrollAreaScrollbar
            className={styles.scrollbar}
            orientation="vertical"
          >
            <ScrollAreaThumb className={styles['scrollbar-thumb']} />
            <ScrollAreaCorner className={styles['scrollbar-corner']} />
          </ScrollAreaScrollbar>
        </Root>
      </div>
    );
  },
);
