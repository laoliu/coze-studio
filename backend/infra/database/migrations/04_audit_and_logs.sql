-- =============================================================================
-- MetaWorkflow V2.0 - 审计与日志表
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建审计日志、API日志等监控相关表
-- 依赖: 00_init_functions.sql, 01_users_and_permissions.sql
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 表: audit_logs - 审计日志表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS audit_logs (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 用户和资源
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    resource_type VARCHAR(50) NOT NULL,  -- workflow, adapter, user, etc.
    resource_id UUID,
    
    -- 操作信息
    action VARCHAR(50) NOT NULL,  -- create, update, delete, publish, etc.
    status VARCHAR(20) NOT NULL,  -- success, failure
    
    -- 详细信息
    changes JSONB,  -- 变更前后对比
    metadata JSONB DEFAULT '{}',
    
    -- IP和客户端
    ip_address INET,
    user_agent TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_audit_status CHECK (status IN ('success', 'failure'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs(status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- GIN索引用于JSON查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_changes ON audit_logs USING GIN(changes);
CREATE INDEX IF NOT EXISTS idx_audit_logs_metadata ON audit_logs USING GIN(metadata);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_created ON audit_logs(resource_type, resource_id, created_at DESC);

COMMENT ON TABLE audit_logs IS '审计日志表';
COMMENT ON COLUMN audit_logs.resource_type IS '资源类型：workflow, adapter, user等';
COMMENT ON COLUMN audit_logs.action IS '操作类型：create, update, delete, publish等';
COMMENT ON COLUMN audit_logs.changes IS '变更前后对比（JSON格式）';

-- -----------------------------------------------------------------------------
-- 表: api_logs - API日志表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS api_logs (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 请求信息
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    query_params JSONB,
    
    -- 用户信息
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- 响应信息
    status_code INTEGER NOT NULL,
    response_time FLOAT NOT NULL,  -- 毫秒
    
    -- 详细信息
    request_body JSONB,
    response_body JSONB,
    error_message TEXT,
    
    -- 客户端信息
    ip_address INET,
    user_agent TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_http_method CHECK (method IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD')),
    CONSTRAINT valid_status_code CHECK (status_code BETWEEN 100 AND 599)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_api_logs_path ON api_logs(path);
CREATE INDEX IF NOT EXISTS idx_api_logs_method ON api_logs(method);
CREATE INDEX IF NOT EXISTS idx_api_logs_status_code ON api_logs(status_code);
CREATE INDEX IF NOT EXISTS idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_api_logs_response_time ON api_logs(response_time DESC);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_api_logs_method_path ON api_logs(method, path);
CREATE INDEX IF NOT EXISTS idx_api_logs_user_created ON api_logs(user_id, created_at DESC);

-- 慢查询索引（响应时间 > 1000ms）
CREATE INDEX IF NOT EXISTS idx_api_logs_slow_queries ON api_logs(response_time DESC, created_at DESC) 
    WHERE response_time > 1000;

-- 错误日志索引（4xx, 5xx）
CREATE INDEX IF NOT EXISTS idx_api_logs_errors ON api_logs(status_code, created_at DESC) 
    WHERE status_code >= 400;

COMMENT ON TABLE api_logs IS 'API访问日志表';
COMMENT ON COLUMN api_logs.response_time IS '响应时间（毫秒）';

-- -----------------------------------------------------------------------------
-- 表: error_logs - 错误日志表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS error_logs (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 错误信息
    error_type VARCHAR(100) NOT NULL,  -- exception, validation, database, network, etc.
    error_code VARCHAR(50),
    error_message TEXT NOT NULL,
    error_stack TEXT,
    
    -- 上下文
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    workflow_id UUID REFERENCES workflows(id) ON DELETE SET NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    
    -- 严重级别
    severity VARCHAR(20) NOT NULL DEFAULT 'error',  -- debug, info, warning, error, critical
    
    -- 详细信息
    context JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    
    -- 客户端信息
    ip_address INET,
    user_agent TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_severity CHECK (severity IN ('debug', 'info', 'warning', 'error', 'critical'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_error_logs_error_type ON error_logs(error_type);
CREATE INDEX IF NOT EXISTS idx_error_logs_severity ON error_logs(severity);
CREATE INDEX IF NOT EXISTS idx_error_logs_user_id ON error_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_workflow_id ON error_logs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_created_at ON error_logs(created_at DESC);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_error_logs_severity_created ON error_logs(severity, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_error_logs_type_created ON error_logs(error_type, created_at DESC);

-- GIN索引
CREATE INDEX IF NOT EXISTS idx_error_logs_context ON error_logs USING GIN(context);

COMMENT ON TABLE error_logs IS '错误日志表';
COMMENT ON COLUMN error_logs.severity IS '严重级别：debug, info, warning, error, critical';

-- -----------------------------------------------------------------------------
-- 表: system_events - 系统事件表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS system_events (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 事件信息
    event_type VARCHAR(100) NOT NULL,  -- user_registered, adapter_published, workflow_completed, etc.
    event_category VARCHAR(50) NOT NULL,  -- user, adapter, workflow, system
    
    -- 关联资源
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    
    -- 事件数据
    event_data JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_event_category CHECK (event_category IN ('user', 'adapter', 'workflow', 'system', 'security'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_system_events_event_type ON system_events(event_type);
CREATE INDEX IF NOT EXISTS idx_system_events_event_category ON system_events(event_category);
CREATE INDEX IF NOT EXISTS idx_system_events_user_id ON system_events(user_id);
CREATE INDEX IF NOT EXISTS idx_system_events_resource ON system_events(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_system_events_created_at ON system_events(created_at DESC);

-- GIN索引
CREATE INDEX IF NOT EXISTS idx_system_events_event_data ON system_events USING GIN(event_data);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_system_events_category_created ON system_events(event_category, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_events_type_created ON system_events(event_type, created_at DESC);

COMMENT ON TABLE system_events IS '系统事件表';
COMMENT ON COLUMN system_events.event_type IS '事件类型：user_registered, adapter_published, workflow_completed等';
COMMENT ON COLUMN system_events.event_category IS '事件分类：user, adapter, workflow, system, security';

-- -----------------------------------------------------------------------------
-- 表: performance_metrics - 性能指标表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS performance_metrics (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 指标信息
    metric_name VARCHAR(100) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,  -- counter, gauge, histogram, summary
    
    -- 指标值
    value FLOAT NOT NULL,
    unit VARCHAR(50),  -- ms, seconds, count, bytes, percent
    
    -- 标签（用于分组和过滤）
    tags JSONB DEFAULT '{}',
    
    -- 时间戳
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_metric_type CHECK (metric_type IN ('counter', 'gauge', 'histogram', 'summary'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_performance_metrics_name ON performance_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_type ON performance_metrics(metric_type);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_timestamp ON performance_metrics(timestamp DESC);

-- GIN索引
CREATE INDEX IF NOT EXISTS idx_performance_metrics_tags ON performance_metrics USING GIN(tags);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_performance_metrics_name_timestamp ON performance_metrics(metric_name, timestamp DESC);

COMMENT ON TABLE performance_metrics IS '性能指标表';
COMMENT ON COLUMN performance_metrics.metric_type IS '指标类型：counter（计数器）, gauge（仪表）, histogram（直方图）, summary（摘要）';
COMMENT ON COLUMN performance_metrics.tags IS '指标标签（JSON格式），用于分组和过滤';

-- 打印创建完成信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 审计与日志表创建完成';
    RAISE NOTICE '   - audit_logs: 审计日志表';
    RAISE NOTICE '   - api_logs: API日志表';
    RAISE NOTICE '   - error_logs: 错误日志表';
    RAISE NOTICE '   - system_events: 系统事件表';
    RAISE NOTICE '   - performance_metrics: 性能指标表';
END $$;
