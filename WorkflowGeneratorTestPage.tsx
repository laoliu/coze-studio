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

/**
 * LLM 工作流生成器测试页面示例
 * Example test page for LLM workflow generator
 *
 * 使用说明：
 * 1. 将此文件内容添加到你的测试路由中
 * 2. 或者在浏览器 Console 中直接测试
 */

import React, { useState } from 'react';
import {
  Button,
  Space,
  Toast,
  Card,
  Descriptions,
  Spin,
} from '@douyinfe/semi-ui';
import { IconMagicWand, IconRefresh } from '@douyinfe/semi-icons';

// 导入工作流生成 API
import {
  generateWorkflow,
  getLLMStatus,
} from '@coze-workflow/playground/services/workflow-generator-api';

export const WorkflowGeneratorTest: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [llmStatus, setLlmStatus] = useState<any>(null);
  const [generatedWorkflow, setGeneratedWorkflow] = useState<any>(null);

  // 测试 LLM 状态
  const testLLMStatus = async () => {
    setLoading(true);
    try {
      const status = await getLLMStatus();
      setLlmStatus(status);
      Toast.success('LLM 状态检查成功');
      console.log('LLM Status:', status);
    } catch (error) {
      Toast.error('LLM 状态检查失败: ' + error.message);
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  // 测试生成工作流
  const testGenerateWorkflow = async (requirement: string) => {
    setLoading(true);
    setGeneratedWorkflow(null);

    try {
      const result = await generateWorkflow({
        space_id: '0', // 替换为实际的空间ID
        user_requirement: requirement,
        language: 'zh-CN',
        requirement_type: 'specific',
        max_nodes: 10,
      });

      setGeneratedWorkflow(result);
      Toast.success('工作流生成成功！');
      console.log('Generated Workflow:', result);
    } catch (error) {
      Toast.error('生成失败: ' + error.message);
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  // 测试用例
  const testCases = [
    '创建一个文本摘要工作流，输入长文本，输出简洁摘要',
    '构建一个智能客服问答系统，从知识库检索信息后生成回答',
    '实现批量文章翻译功能，从数据库读取文章并逐个翻译',
    '设计一个简历解析器，提取姓名、教育、工作经验等信息',
  ];

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      <h1>🤖 LLM 工作流生成器测试</h1>

      {/* 状态检查 */}
      <Card
        title="1️⃣ LLM 服务状态检查"
        style={{ marginBottom: 24 }}
        extra={
          <Button
            icon={<IconRefresh />}
            onClick={testLLMStatus}
            loading={loading}
          >
            检查状态
          </Button>
        }
      >
        {llmStatus ? (
          <Descriptions
            data={[
              {
                key: '健康状态',
                value: llmStatus.healthy ? '✅ 正常' : '❌ 异常',
              },
              {
                key: '启用状态',
                value: llmStatus.enabled ? '✅ 已启用' : '❌ 未启用',
              },
              { key: '模型名称', value: llmStatus.model_name || 'N/A' },
              { key: '模型ID', value: llmStatus.model_id || 'N/A' },
              {
                key: '温度参数',
                value: llmStatus.configuration?.temperature || 'N/A',
              },
              {
                key: '最大Token',
                value: llmStatus.configuration?.max_tokens || 'N/A',
              },
            ]}
          />
        ) : (
          <p style={{ color: '#999' }}>点击 "检查状态" 按钮查看 LLM 服务状态</p>
        )}
      </Card>

      {/* 生成测试 */}
      <Card title="2️⃣ 工作流生成测试" style={{ marginBottom: 24 }}>
        <Space wrap>
          {testCases.map((testCase, index) => (
            <Button
              key={index}
              icon={<IconMagicWand />}
              onClick={() => testGenerateWorkflow(testCase)}
              loading={loading}
              type={index === 0 ? 'primary' : 'secondary'}
            >
              测试用例 {index + 1}
            </Button>
          ))}
        </Space>

        <div style={{ marginTop: 16 }}>
          {testCases.map((testCase, index) => (
            <p key={index} style={{ fontSize: 12, color: '#666' }}>
              用例 {index + 1}: {testCase}
            </p>
          ))}
        </div>
      </Card>

      {/* 生成结果 */}
      {loading && (
        <Card title="⏳ 生成中..." style={{ marginBottom: 24 }}>
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin size="large" />
            <p style={{ marginTop: 16 }}>AI 正在分析需求并生成工作流...</p>
          </div>
        </Card>
      )}

      {generatedWorkflow && !loading && (
        <Card title="3️⃣ 生成结果" style={{ marginBottom: 24 }}>
          <Descriptions
            data={[
              { key: '工作流ID', value: generatedWorkflow.workflow_id },
              { key: '工作流名称', value: generatedWorkflow.workflow_name },
              { key: '节点数量', value: generatedWorkflow.nodes?.length || 0 },
              { key: '连接数量', value: generatedWorkflow.edges?.length || 0 },
              {
                key: '置信度',
                value: (generatedWorkflow.confidence * 100).toFixed(1) + '%',
              },
              { key: '说明', value: generatedWorkflow.explanation },
            ]}
          />

          {/* 节点列表 */}
          <h4 style={{ marginTop: 24, marginBottom: 12 }}>节点列表：</h4>
          <div
            style={{
              maxHeight: 300,
              overflow: 'auto',
              background: '#f7f7f7',
              padding: 12,
              borderRadius: 4,
            }}
          >
            {generatedWorkflow.nodes?.map((node: any, index: number) => (
              <div
                key={index}
                style={{
                  marginBottom: 8,
                  padding: 8,
                  background: 'white',
                  borderRadius: 4,
                }}
              >
                <strong>[{node.type}]</strong> {node.name}
                {node.description && (
                  <div style={{ fontSize: 12, color: '#666', marginTop: 4 }}>
                    {node.description}
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* 完整 JSON */}
          <h4 style={{ marginTop: 24, marginBottom: 12 }}>完整 JSON：</h4>
          <pre
            style={{
              maxHeight: 400,
              overflow: 'auto',
              background: '#f7f7f7',
              padding: 12,
              borderRadius: 4,
              fontSize: 12,
            }}
          >
            {JSON.stringify(generatedWorkflow, null, 2)}
          </pre>
        </Card>
      )}

      {/* Console 测试 */}
      <Card title="4️⃣ 浏览器 Console 测试">
        <p>在浏览器开发者工具 Console 中运行以下代码：</p>
        <pre
          style={{
            background: '#2d2d2d',
            color: '#f8f8f2',
            padding: 16,
            borderRadius: 4,
            fontSize: 12,
            overflow: 'auto',
          }}
        >
          {`// 测试健康检查
fetch('/api/workflow_api/llm_status', {
  method: 'GET',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' }
})
.then(r => r.json())
.then(data => console.log('LLM 状态:', data));

// 测试生成工作流
fetch('/api/workflow_api/generate', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    space_id: '0',
    user_requirement: '创建一个文本摘要工作流',
    language: 'zh-CN'
  })
})
.then(r => r.json())
.then(data => console.log('生成结果:', data));`}
        </pre>
      </Card>
    </div>
  );
};

export default WorkflowGeneratorTest;
