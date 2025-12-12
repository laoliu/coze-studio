import React, { useState } from 'react';
import { Modal, Input, Button, Steps, Spin, Toast } from '@douyinfe/semi-ui';
import { IconCheckCircleStroked } from '@douyinfe/semi-icons';
import WorkflowPreview from './WorkflowPreview';
import { generateWorkflow, getExampleRequirements, validateGenerateRequest } from '../../services/workflow-generator-api';
import styles from './index.module.less';

const { TextArea } = Input;

interface WorkflowGeneratorProps {
  visible: boolean;
  onClose: () => void;
  onGenerate: (workflow: any) => void;
}

/**
 * 工作流智能生成器组件
 * 
 * 功能：
 * 1. 接收用户的自然语言需求描述
 * 2. 调用后端 API 生成完整工作流
 * 3. 展示生成结果预览
 * 4. 支持用户确认或重新生成
 */
export const WorkflowGenerator: React.FC<WorkflowGeneratorProps> = ({
  visible,
  onClose,
  onGenerate,
}) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [requirement, setRequirement] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedWorkflow, setGeneratedWorkflow] = useState<any>(null);

  // 示例需求模板
  const exampleRequirements = [
    '创建一个课程大纲生成器，用户输入课程主题，系统生成包含章节的详细大纲',
    '构建一个智能客服问答系统，从知识库检索相关信息后生成回答',
    '实现批量文章翻译功能，读取数据库中的文章列表并逐个翻译',
    '设计一个简历信息提取工作流，从文本中提取姓名、经历等结构化信息',
  ];

  // 生成工作流
  const handleGenerate = async () => {
    if (!requirement.trim()) {
      Toast.warning('请输入需求描述');
      return;
    }

    setIsGenerating(true);
    setCurrentStep(1);

    try {
      // TODO: 调用实际的 API
      // const response = await generateWorkflowAPI({
      //   user_requirement: requirement,
      //   language: 'zh-CN',
      //   constraints: { max_nodes: 10 },
      // });

      // 模拟 API 调用
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // 临时模拟数据
      const mockWorkflow = {
        workflow_id: 'workflow_mock_123',
        nodes: [],
        edges: [],
        explanation: '这是一个模拟的工作流说明',
        confidence: 0.85,
      };

      setGeneratedWorkflow(mockWorkflow);
      setCurrentStep(2);
      Toast.success('工作流生成成功！');
    } catch (error) {
      Toast.error('生成工作流失败，请重试');
      setCurrentStep(0);
      console.error(error);
    } finally {
      setIsGenerating(false);
    }
  };

  // 确认生成
  const handleConfirm = () => {
    if (generatedWorkflow) {
      onGenerate(generatedWorkflow);
      handleReset();
      onClose();
    }
  };

  // 重新生成
  const handleRegenerate = () => {
    setCurrentStep(0);
    setGeneratedWorkflow(null);
  };

  // 重置状态
  const handleReset = () => {
    setCurrentStep(0);
    setRequirement('');
    setGeneratedWorkflow(null);
    setIsGenerating(false);
  };

  // 关闭弹窗
  const handleClose = () => {
    handleReset();
    onClose();
  };

  return (
    <Modal
      title="🤖 智能工作流生成器"
      visible={visible}
      onCancel={handleClose}
      footer={null}
      style={{ width: 800 }}
      className={styles.workflowGeneratorModal}
    >
      <div className={styles.container}>
        {/* 步骤指示器 */}
        <Steps current={currentStep} className={styles.steps}>
          <Steps.Step title="描述需求" />
          <Steps.Step title="生成工作流" />
          <Steps.Step title="预览确认" />
        </Steps>

        {/* 步骤 1: 输入需求 */}
        {currentStep === 0 && (
          <div className={styles.step1}>
            <div className={styles.inputSection}>
              <label className={styles.label}>
                请详细描述您想要创建的工作流功能和流程：
              </label>
              <TextArea
                rows={6}
                placeholder="例如：创建一个课程大纲生成器，用户输入课程主题，系统生成包含章节的详细大纲..."
                value={requirement}
                onChange={setRequirement}
                className={styles.textarea}
              />
            </div>

            <div className={styles.examplesSection}>
              <div className={styles.examplesTitle}>💡 示例需求：</div>
              {exampleRequirements.map((example, index) => (
                <div
                  key={index}
                  className={styles.exampleItem}
                  onClick={() => setRequirement(example)}
                >
                  {example}
                </div>
              ))}
            </div>

            <div className={styles.actions}>
              <Button onClick={handleClose}>取消</Button>
              <Button
                type="primary"
                onClick={handleGenerate}
                disabled={!requirement.trim()}
              >
                🎯 智能生成工作流
              </Button>
            </div>
          </div>
        )}

        {/* 步骤 2: 生成中 */}
        {currentStep === 1 && (
          <div className={styles.step2}>
            <div className={styles.loading}>
              <Spin size={48} />
              <p className={styles.loadingText}>
                🤖 正在分析需求并生成工作流...
              </p>
              <p className={styles.loadingHint}>
                这可能需要 5-10 秒，请稍候
              </p>
            </div>
          </div>
        )}

        {/* 步骤 3: 预览结果 */}
        {currentStep === 2 && generatedWorkflow && (
          <div className={styles.step3}>
            <div className={styles.preview}>
              <h3>📊 生成的工作流预览</h3>
              <div className={styles.workflowInfo}>
                <p>
                  <strong>置信度：</strong>
                  {(generatedWorkflow.confidence * 100).toFixed(0)}%
                </p>
                <p>
                  <strong>节点数量：</strong>
                  {generatedWorkflow.nodes?.length || 0}
                </p>
                <p>
                  <strong>说明：</strong>
                  {generatedWorkflow.explanation}
                </p>
              </div>

              <div className={styles.notice}>
                ℹ️ 提示：生成的工作流将自动添加到画布中，您可以进一步编辑和调整。
              </div>
            </div>

            <div className={styles.actions}>
              <Button onClick={handleRegenerate}>重新生成</Button>
              <Button type="primary" onClick={handleConfirm} icon={<IconCheckCircleStroked />}>
                确认创建
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default WorkflowGenerator;
