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

import React from 'react';

import { IconCozStarFill } from '@coze-arch/coze-design/icons';

interface WorkflowResultProps {
  result: {
    workflow_name: string;
    explanation: string;
    nodes: Array<{
      id: string;
      type: string;
      name: string;
      description: string;
    }>;
  };
}

export const WorkflowResult: React.FC<WorkflowResultProps> = ({ result }) => {
  return (
    <div
      style={{
        padding: '16px',
        backgroundColor: '#f5f7fa',
        borderRadius: '8px',
        marginTop: '16px',
      }}
    >
      <div
        style={{ display: 'flex', alignItems: 'center', marginBottom: '12px' }}
      >
        <IconCozStarFill
          style={{ color: '#FFB800', marginRight: '8px', fontSize: '18px' }}
        />
        <h3 style={{ margin: 0, fontSize: '16px', fontWeight: 600 }}>
          {result.workflow_name}
        </h3>
      </div>

      <p style={{ color: '#666', marginBottom: '16px', lineHeight: '1.6' }}>
        {result.explanation}
      </p>

      <div>
        <h4 style={{ fontSize: '14px', fontWeight: 600, marginBottom: '12px' }}>
          Workflow Nodes:
        </h4>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
          {result.nodes.map((node, index) => (
            <div
              key={index}
              style={{
                padding: '12px',
                backgroundColor: '#fff',
                borderRadius: '6px',
                border: '1px solid #e5e7eb',
              }}
            >
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  marginBottom: '4px',
                }}
              >
                <span
                  style={{
                    display: 'inline-block',
                    width: '20px',
                    height: '20px',
                    lineHeight: '20px',
                    textAlign: 'center',
                    backgroundColor: '#4f46e5',
                    color: '#fff',
                    borderRadius: '50%',
                    fontSize: '12px',
                    marginRight: '8px',
                  }}
                >
                  {index + 1}
                </span>
                <span style={{ fontWeight: 600, fontSize: '14px' }}>
                  {node.name}
                </span>
                <span
                  style={{
                    marginLeft: '8px',
                    padding: '2px 8px',
                    backgroundColor: '#e0e7ff',
                    color: '#4f46e5',
                    borderRadius: '4px',
                    fontSize: '12px',
                  }}
                >
                  {node.type}
                </span>
              </div>
              <p
                style={{
                  margin: 0,
                  color: '#666',
                  fontSize: '13px',
                  paddingLeft: '28px',
                }}
              >
                {node.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
