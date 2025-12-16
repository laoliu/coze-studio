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

interface GeneratedNode {
  id: string;
  name: string;
  type: string;
  description: string;
  position?: { x: number; y: number };
  config?: Record<string, unknown>;
}

interface GeneratedEdge {
  id: string;
  source: string;
  target: string;
  type: string;
}

// Frontend edge format (what the canvas expects)
interface FrontendEdge {
  sourceNodeID: string;
  targetNodeID: string;
  sourcePortID?: string;
  targetPortID?: string;
}

const AUTH_FAILED_CODE = 700012006;

interface CreateWorkflowApiResponse {
  code?: number;
  msg?: string;
  data?: {
    workflow_id?: string;
  };
  workflow_id?: string;
}

const extractWorkflowId = (createData: CreateWorkflowApiResponse): string => {
  const workflowId = createData.data?.workflow_id || createData.workflow_id;

  if (!workflowId) {
    console.error('Create response:', createData);
    throw new Error('Failed to get workflow ID from create response');
  }

  return workflowId;
};

const createWorkflowMetadata = async (params: {
  name: string;
  desc: string;
  spaceId: string;
}) => {
  const response = await fetch('/api/workflow_api/create', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: params.name,
      desc: params.desc,
      icon_uri: 'workflow',
      space_id: params.spaceId,
    }),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();

  if (data.code !== 0 && data.code !== undefined) {
    if (data.code === AUTH_FAILED_CODE) {
      throw new Error('Authentication failed. Please login and try again.');
    }
    throw new Error(data.msg || 'Failed to create workflow');
  }

  return data;
};

const saveWorkflowCanvas = async (params: {
  workflowId: string;
  spaceId: string;
  nodes: GeneratedNode[];
  edges: GeneratedEdge[];
}) => {
  // Convert backend edge format to frontend edge format
  const frontendEdges: FrontendEdge[] = params.edges.map(edge => ({
    sourceNodeID: edge.source,
    targetNodeID: edge.target,
  }));

  console.log('[saveWorkflowCanvas] Backend edges:', params.edges);
  console.log('[saveWorkflowCanvas] Frontend edges:', frontendEdges);

  // Create schema object with nodes and edges in frontend format
  const schema = {
    nodes: params.nodes,
    edges: frontendEdges,
  };

  console.log('[saveWorkflowCanvas] Schema to save:', JSON.stringify(schema, null, 2));

  const response = await fetch('/api/workflow_api/save', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      workflow_id: params.workflowId,
      space_id: params.spaceId,
      schema: JSON.stringify(schema),
      submit_commit_id: '', // Empty for new workflow
    }),
  });

  if (!response.ok) {
    console.error('Failed to save workflow canvas, but workflow was created');
    return;
  }

  const data = await response.json();
  if (data.code !== 0 && data.code !== undefined) {
    console.error('Failed to save workflow canvas:', data.msg);
  }
};

export const createWorkflow = async (params: {
  name: string;
  desc: string;
  spaceId: string;
  nodes?: GeneratedNode[];
  edges?: GeneratedEdge[];
}) => {
  // Step 1: Create workflow metadata
  const createData = await createWorkflowMetadata(params);
  const workflowId = extractWorkflowId(createData);

  // Step 2: Save canvas data if nodes are provided
  if (params.nodes && params.nodes.length > 0) {
    await saveWorkflowCanvas({
      workflowId,
      spaceId: params.spaceId,
      nodes: params.nodes,
      edges: params.edges || [],
    });
  }

  return {
    ...createData,
    workflow_id: workflowId,
  };
};
