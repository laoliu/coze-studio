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

import React, { useState } from 'react';

import { Modal, Button, Toast } from '@coze-arch/coze-design';

import { WorkflowResult } from './workflow-result';
import { RequirementInput } from './requirement-input';
import { createWorkflow } from './create-workflow';

interface WorkflowGeneratorModalProps {
  visible: boolean;
  spaceId: string;
  onClose: () => void;
}

interface GeneratedWorkflow {
  // Backend API uses PascalCase for status fields
  // eslint-disable-next-line @typescript-eslint/naming-convention -- Backend API convention
  StatusCode: number;
  // eslint-disable-next-line @typescript-eslint/naming-convention -- Backend API convention
  StatusMessage: string;
  workflow_id: string;
  workflow_name: string;
  explanation: string;
  nodes: Array<{
    id: string;
    type: string;
    name: string;
    description: string;
  }>;
  edges: Array<{
    id: string;
    source: string;
    target: string;
    type: string;
  }>;
  confidence: number;
}

export const WorkflowGeneratorModal: React.FC<WorkflowGeneratorModalProps> = ({
  visible,
  spaceId,
  onClose,
}) => {
  const [requirement, setRequirement] = useState('');
  const [loading, setLoading] = useState(false);
  const [generatedWorkflow, setGeneratedWorkflow] =
    useState<GeneratedWorkflow | null>(null);

  const handleGenerate = async () => {
    if (!requirement.trim()) {
      Toast.error('Please enter workflow requirements');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('/api/workflow_api/generate', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_requirement: requirement,
          space_id: spaceId,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const responseText = await response.text();
      const MAX_LOG_LENGTH = 500;
      console.log('Raw response:', responseText.substring(0, MAX_LOG_LENGTH));

      let data;
      try {
        data = JSON.parse(responseText);
      } catch (parseError) {
        console.error('JSON parse error:', parseError);
        console.error('Response length:', responseText.length);
        throw new Error('Invalid response format from server');
      }

      // Backend returns {StatusCode: 0, workflow_name, nodes, edges, ...}
      console.log('Parsed data:', data);
      console.log('StatusCode:', data.StatusCode);
      console.log('Has nodes:', !!data.nodes, 'Count:', data.nodes?.length);

      if (data.StatusCode === 0) {
        console.log('Setting generated workflow:', data);
        setGeneratedWorkflow(data);
        Toast.success(
          'Workflow generated successfully! You can find it in the workflow list.',
        );
        console.log('Workflow state should be updated');
      } else {
        throw new Error(data.StatusMessage || 'Failed to generate workflow');
      }
    } catch (error) {
      console.error('Failed to generate workflow:', error);
      Toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to generate workflow. Please try again.',
      );
    } finally {
      setLoading(false);
    }
  };

  const handleCreateWorkflow = async () => {
    if (!generatedWorkflow) {
      return;
    }

    setLoading(true);
    try {
      const result = await createWorkflow({
        name: generatedWorkflow.workflow_name,
        desc: generatedWorkflow.explanation || 'AI Generated Workflow',
        spaceId,
        nodes: generatedWorkflow.nodes,
        edges: generatedWorkflow.edges,
      });

      console.log('Workflow created successfully:', result);

      Toast.success('Workflow created successfully!');
      handleClose();

      // Reload to show the new workflow in the list
      // The user can then click on it to open the editor
      window.location.reload();
    } catch (error) {
      console.error('Failed to create workflow:', error);
      Toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to create workflow. Please try again.',
      );
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setRequirement('');
    setGeneratedWorkflow(null);
    onClose();
  };

  return (
    <Modal
      visible={visible}
      onCancel={handleClose}
      title="AI Workflow Generator"
      width={700}
      footer={
        <div
          style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}
        >
          <Button onClick={handleClose}>Cancel</Button>
          {!generatedWorkflow ? (
            <Button type="primary" onClick={handleGenerate} loading={loading}>
              Generate
            </Button>
          ) : (
            <Button type="primary" onClick={handleCreateWorkflow}>
              Create Workflow
            </Button>
          )}
        </div>
      }
    >
      <RequirementInput
        requirement={requirement}
        setRequirement={setRequirement}
      />

      {generatedWorkflow ? <WorkflowResult result={generatedWorkflow} /> : null}
    </Modal>
  );
};
