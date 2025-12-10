-- =============================================================================
-- MetaWorkflow V2.0 - 核心业务表
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建工作流、节点执行、学习目标、教学活动相关表
-- 依赖: 00_init_functions.sql, 01_users_and_permissions.sql
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 表: workflows - 工作流记录表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS workflows (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联用户
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- 领域和状态
    domain VARCHAR(50) NOT NULL,  -- k12_education, art_history, bestseller, vocational
    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, running, completed, failed, cancelled
    
    -- 输入信息
    raw_input TEXT NOT NULL,
    subject VARCHAR(100),  -- 学科：化学、物理等
    grade_level INTEGER,   -- 年级：1-12
    duration INTEGER,      -- 课程时长（分钟）
    activity_type VARCHAR(50),  -- 活动类型：concept_understanding, experiment_exploration, etc.
    
    -- 执行信息
    current_node VARCHAR(255),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    execution_time FLOAT,  -- 执行时间（秒）
    
    -- 上下文数据（JSON）
    context_data JSONB DEFAULT '{}',
    outputs JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    
    -- 错误信息
    error_message TEXT,
    error_stack TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    CONSTRAINT valid_domain CHECK (domain IN ('k12_education', 'art_history', 'bestseller', 'vocational', 'higher_education', 'professional_training')),
    CONSTRAINT valid_grade_level CHECK (grade_level BETWEEN 1 AND 12 OR grade_level IS NULL)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_workflows_user_id ON workflows(user_id);
CREATE INDEX IF NOT EXISTS idx_workflows_domain ON workflows(domain);
CREATE INDEX IF NOT EXISTS idx_workflows_status ON workflows(status);
CREATE INDEX IF NOT EXISTS idx_workflows_created_at ON workflows(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_workflows_completed_at ON workflows(completed_at DESC) WHERE completed_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_workflows_subject ON workflows(subject) WHERE subject IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_workflows_activity_type ON workflows(activity_type) WHERE activity_type IS NOT NULL;

-- GIN索引用于JSON查询
CREATE INDEX IF NOT EXISTS idx_workflows_context_data ON workflows USING GIN(context_data);
CREATE INDEX IF NOT EXISTS idx_workflows_outputs ON workflows USING GIN(outputs);
CREATE INDEX IF NOT EXISTS idx_workflows_metadata ON workflows USING GIN(metadata);

-- 触发器
DROP TRIGGER IF EXISTS update_workflows_updated_at ON workflows;
CREATE TRIGGER update_workflows_updated_at
    BEFORE UPDATE ON workflows
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE workflows IS '工作流执行记录表';
COMMENT ON COLUMN workflows.domain IS '领域类型：k12_education（K12教育）, art_history（艺术史）, bestseller（畅销书）, vocational（职业培训）';
COMMENT ON COLUMN workflows.duration IS '课程时长（分钟），微课3-15分钟，常规课45-90分钟';
COMMENT ON COLUMN workflows.context_data IS '工作流上下文数据（JSON格式）';
COMMENT ON COLUMN workflows.outputs IS '工作流输出结果（JSON格式）';

-- -----------------------------------------------------------------------------
-- 表: node_executions - 节点执行记录表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS node_executions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联工作流
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    
    -- 节点信息
    node_id VARCHAR(255) NOT NULL,
    node_type VARCHAR(100) NOT NULL,  -- adapter_invoke, capability_invoke, objective_generation, etc.
    
    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, running, completed, failed, skipped
    
    -- 输入输出
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    
    -- 执行信息
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    execution_time FLOAT,  -- 执行时间（秒）
    
    -- 错误信息
    error_message TEXT,
    error_stack TEXT,
    retry_count INTEGER DEFAULT 0,
    
    -- 元数据
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_node_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'skipped'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_node_executions_workflow_id ON node_executions(workflow_id);
CREATE INDEX IF NOT EXISTS idx_node_executions_node_type ON node_executions(node_type);
CREATE INDEX IF NOT EXISTS idx_node_executions_status ON node_executions(status);
CREATE INDEX IF NOT EXISTS idx_node_executions_created_at ON node_executions(created_at DESC);

-- GIN索引用于JSON查询
CREATE INDEX IF NOT EXISTS idx_node_executions_input_data ON node_executions USING GIN(input_data);
CREATE INDEX IF NOT EXISTS idx_node_executions_output_data ON node_executions USING GIN(output_data);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_node_executions_workflow_created ON node_executions(workflow_id, created_at DESC);

COMMENT ON TABLE node_executions IS '节点执行记录表';
COMMENT ON COLUMN node_executions.node_type IS '节点类型：adapter_invoke, capability_invoke, objective_generator等';
COMMENT ON COLUMN node_executions.retry_count IS '重试次数';

-- -----------------------------------------------------------------------------
-- 表: learning_objectives - 学习目标表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS learning_objectives (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联工作流
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    
    -- 目标内容
    title VARCHAR(500) NOT NULL,
    description TEXT,
    bloom_level VARCHAR(20) NOT NULL,  -- remember, understand, apply, analyze, evaluate, create
    
    -- 时间和顺序
    estimated_time INTEGER NOT NULL,  -- 预估时长（分钟）
    sequence_order INTEGER NOT NULL DEFAULT 0,
    
    -- 关键词和前置知识
    keywords TEXT[] DEFAULT '{}',
    prerequisites TEXT[] DEFAULT '{}',
    
    -- 元数据
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_bloom_level CHECK (bloom_level IN ('remember', 'understand', 'apply', 'analyze', 'evaluate', 'create')),
    CONSTRAINT positive_estimated_time CHECK (estimated_time > 0)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_learning_objectives_workflow_id ON learning_objectives(workflow_id);
CREATE INDEX IF NOT EXISTS idx_learning_objectives_bloom_level ON learning_objectives(bloom_level);
CREATE INDEX IF NOT EXISTS idx_learning_objectives_sequence ON learning_objectives(workflow_id, sequence_order);

-- GIN索引用于数组查询
CREATE INDEX IF NOT EXISTS idx_learning_objectives_keywords ON learning_objectives USING GIN(keywords);
CREATE INDEX IF NOT EXISTS idx_learning_objectives_prerequisites ON learning_objectives USING GIN(prerequisites);

COMMENT ON TABLE learning_objectives IS '学习目标表';
COMMENT ON COLUMN learning_objectives.bloom_level IS 'Bloom分类法层级：remember（记忆）, understand（理解）, apply（应用）, analyze（分析）, evaluate（评价）, create（创造）';

-- -----------------------------------------------------------------------------
-- 表: teaching_activities - 教学活动表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS teaching_activities (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联工作流
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    
    -- 活动类型
    activity_type VARCHAR(50) NOT NULL,  -- concept_understanding, experiment_exploration, problem_solving
    title VARCHAR(500) NOT NULL,
    
    -- 活动内容（不同类型存储不同结构）
    content JSONB NOT NULL,
    
    -- 元数据
    duration INTEGER,  -- 活动时长（分钟）
    difficulty VARCHAR(50),  -- easy, medium, hard
    tags TEXT[] DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_activity_type CHECK (activity_type IN (
        'concept_understanding', 
        'experiment_exploration', 
        'problem_solving',
        'creative_synthesis',
        'collaborative_learning',
        'assessment'
    )),
    CONSTRAINT valid_difficulty CHECK (difficulty IN ('easy', 'medium', 'hard') OR difficulty IS NULL)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_teaching_activities_workflow_id ON teaching_activities(workflow_id);
CREATE INDEX IF NOT EXISTS idx_teaching_activities_type ON teaching_activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_teaching_activities_difficulty ON teaching_activities(difficulty) WHERE difficulty IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_teaching_activities_created_at ON teaching_activities(created_at DESC);

-- GIN索引
CREATE INDEX IF NOT EXISTS idx_teaching_activities_content ON teaching_activities USING GIN(content);
CREATE INDEX IF NOT EXISTS idx_teaching_activities_tags ON teaching_activities USING GIN(tags);

-- 触发器
DROP TRIGGER IF EXISTS update_teaching_activities_updated_at ON teaching_activities;
CREATE TRIGGER update_teaching_activities_updated_at
    BEFORE UPDATE ON teaching_activities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE teaching_activities IS '教学活动表';
COMMENT ON COLUMN teaching_activities.activity_type IS '活动类型：concept_understanding（概念理解）, experiment_exploration（实验探索）, problem_solving（问题解决）等';
COMMENT ON COLUMN teaching_activities.content IS '活动内容（JSON格式，不同类型存储不同结构）';

-- 打印创建完成信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 核心业务表创建完成';
    RAISE NOTICE '   - workflows: 工作流记录表';
    RAISE NOTICE '   - node_executions: 节点执行记录表';
    RAISE NOTICE '   - learning_objectives: 学习目标表';
    RAISE NOTICE '   - teaching_activities: 教学活动表';
END $$;
