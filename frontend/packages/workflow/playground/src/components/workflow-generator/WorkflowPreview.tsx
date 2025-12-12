import React from 'react';
import { Card, Tag, Descriptions, Space } from '@douyinfe/semi-ui';
import styles from './WorkflowPreview.module.less';

interface Node {
  id: string;
  type: string;
  title: string;
  description?: string;
  config?: Record<string, any>;
  position?: {
    x: number;
    y: number;
  };
}

interface Edge {
  id: string;
  source: string;
  target: string;
  sourceHandle?: string;
  targetHandle?: string;
}

interface WorkflowPreviewProps {
  nodes: Node[];
  edges: Edge[];
  explanation?: string;
  confidence?: number;
}

const WorkflowPreview: React.FC<WorkflowPreviewProps> = ({
  nodes,
  edges,
  explanation,
  confidence
}) => {
  const getNodeTypeColor = (type: string): string => {
    const colorMap: Record<string, string> = {
      start: 'green',
      end: 'red',
      llm: 'blue',
      http: 'orange',
      code: 'purple',
      condition: 'gold',
      loop: 'cyan',
      data_transform: 'magenta',
      default: 'gray'
    };
    return colorMap[type] || colorMap.default;
  };

  return (
    <div className={styles.workflowPreview}>
      {/* 基本信息 */}
      <Card title="工作流概览" bordered={false} className={styles.section}>
        <Descriptions
          column={2}
          data={[
            { label: '节点数量', value: nodes.length },
            { label: '连接数量', value: edges.length },
            { label: '置信度', value: confidence ? `${(confidence * 100).toFixed(1)}%` : '未知' }
          ]}
        />
        {explanation && (
          <div className={styles.explanation}>
            <h4>生成说明</h4>
            <p>{explanation}</p>
          </div>
        )}
      </Card>

      {/* 节点列表 */}
      <Card title="节点列表" bordered={false} className={styles.section}>
        <Space direction="vertical" size="medium" style={{ width: '100%' }}>
          {nodes.map((node, index) => (
            <Card
              key={node.id}
              size="small"
              bordered
              className={styles.nodeCard}
            >
              <div className={styles.nodeHeader}>
                <Space>
                  <span className={styles.nodeIndex}>#{index + 1}</span>
                  <Tag color={getNodeTypeColor(node.type)}>{node.type}</Tag>
                  <strong>{node.title}</strong>
                </Space>
              </div>
              {node.description && (
                <p className={styles.nodeDescription}>{node.description}</p>
              )}
              {node.config && Object.keys(node.config).length > 0 && (
                <div className={styles.nodeConfig}>
                  <div className={styles.configTitle}>配置:</div>
                  <pre className={styles.configContent}>
                    {JSON.stringify(node.config, null, 2)}
                  </pre>
                </div>
              )}
            </Card>
          ))}
        </Space>
      </Card>

      {/* 连接关系 */}
      <Card title="连接关系" bordered={false} className={styles.section}>
        <div className={styles.edgesList}>
          {edges.map((edge, index) => {
            const sourceNode = nodes.find(n => n.id === edge.source);
            const targetNode = nodes.find(n => n.id === edge.target);
            return (
              <div key={edge.id} className={styles.edgeItem}>
                <span className={styles.edgeIndex}>{index + 1}.</span>
                <span className={styles.edgeSource}>
                  {sourceNode?.title || edge.source}
                </span>
                <span className={styles.edgeArrow}>→</span>
                <span className={styles.edgeTarget}>
                  {targetNode?.title || edge.target}
                </span>
              </div>
            );
          })}
        </div>
      </Card>
    </div>
  );
};

export default WorkflowPreview;
