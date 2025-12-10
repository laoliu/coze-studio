# MetaWorkflow V2.0 - 数据库设计文档

**版本**: 2.0  
**日期**: 2025-12-09  
**状态**: 设计完成 ✅

---

## 📋 目录

1. [概述](#概述)
2. [数据库架构](#数据库架构)
3. [核心业务表](#核心业务表)
4. [适配器市场表](#适配器市场表)
5. [用户与权限表](#用户与权限表)
6. [统计与日志表](#统计与日志表)
7. [索引设计](#索引设计)
8. [数据迁移](#数据迁移)
9. [性能优化](#性能优化)
10. [备份策略](#备份策略)

---

## 概述

### 设计原则

```yaml
原则:
  - 规范化: 遵循3NF，减少数据冗余
  - 可扩展: 支持新领域和活动类型
  - 高性能: 合理使用索引和分区
  - 可追溯: 完整的审计日志
  - 灵活性: JSON字段存储动态数据
```

### 技术栈

| 组件 | 技术选型 | 说明 |
|------|---------|------|
| **主数据库** | PostgreSQL 14+ | ACID事务、JSON支持 |
| **缓存** | Redis 6+ | 热数据缓存 |
| **搜索** | Elasticsearch 8+ | 全文搜索、推荐 |
| **分析** | ClickHouse 23+ | 实时统计分析 |
| **向量库** | Milvus 2+ | 语义搜索 |

### 数据库分类

```
metaworkflow (主库)
├── Core Schema         # 核心业务
├── Marketplace Schema  # 适配器市场
├── Auth Schema         # 用户权限
└── Analytics Schema    # 统计分析

metaworkflow_cache (Redis)
├── Session Cache       # 会话缓存
├── Workflow Cache      # 工作流缓存
└── Query Cache         # 查询缓存

metaworkflow_search (Elasticsearch)
├── adapters           # 适配器搜索
└── content            # 内容搜索

metaworkflow_analytics (ClickHouse)
├── events             # 事件流
└── metrics            # 指标聚合
```

---

## 数据库架构

### ER 关系图

```mermaid
erDiagram
    %% 核心业务实体
    USERS ||--o{ WORKFLOWS : creates
    USERS ||--o{ ADAPTERS : develops
    
    WORKFLOWS ||--|{ NODE_EXECUTIONS : contains
    WORKFLOWS ||--|{ LEARNING_OBJECTIVES : has
    WORKFLOWS ||--o| TEACHING_ACTIVITIES : generates
    
    ADAPTERS ||--o{ ADAPTER_VERSIONS : has
    ADAPTERS ||--o{ ADAPTER_REVIEWS : receives
    ADAPTERS }o--o{ ADAPTER_DEPENDENCIES : depends
    
    USERS ||--o{ ADAPTER_REVIEWS : writes
    USERS ||--o{ USER_SUBSCRIPTIONS : subscribes
    
    %% 统计与日志
    WORKFLOWS ||--o{ AUDIT_LOGS : tracks
    ADAPTERS ||--o{ ADAPTER_STATS : measures
```

### 表清单

| 类别 | 表名 | 说明 | 估计行数 |
|------|------|------|---------|
| **核心** | users | 用户表 | 10K |
| **核心** | workflows | 工作流记录 | 100K |
| **核心** | node_executions | 节点执行记录 | 1M |
| **核心** | learning_objectives | 学习目标 | 500K |
| **核心** | teaching_activities | 教学活动 | 100K |
| **市场** | adapters | 适配器 | 1K |
| **市场** | adapter_versions | 版本 | 5K |
| **市场** | adapter_reviews | 评论 | 10K |
| **市场** | adapter_dependencies | 依赖关系 | 2K |
| **市场** | adapter_stats | 统计 | 1K |
| **市场** | adapter_categories | 分类 | 100 |
| **权限** | roles | 角色 | 10 |
| **权限** | permissions | 权限 | 100 |
| **权限** | user_roles | 用户角色 | 10K |
| **订阅** | user_subscriptions | 用户订阅 | 5K |
| **日志** | audit_logs | 审计日志 | 10M |
| **日志** | api_logs | API日志 | 100M |

---

## 核心业务表

### 1. users - 用户表

```sql
CREATE TABLE users (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 基本信息
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url VARCHAR(500),
    
    -- 账户状态
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE NOT NULL,
    is_developer BOOLEAN DEFAULT FALSE NOT NULL,
    
    -- 元数据
    preferences JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_login_at TIMESTAMP WITH TIME ZONE,
    
    -- 索引
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

-- 索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.is_developer IS '是否为开发者（可发布适配器）';
```

### 2. workflows - 工作流记录表

```sql
CREATE TABLE workflows (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联用户
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- 领域和状态
    domain VARCHAR(50) NOT NULL,  -- k12_education, art_history, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, running, completed, failed, cancelled
    
    -- 输入信息
    raw_input TEXT NOT NULL,
    subject VARCHAR(100),  -- 学科：化学、物理等
    grade_level INTEGER,  -- 年级：1-12
    duration INTEGER,  -- 课程时长（分钟）
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
    CONSTRAINT valid_domain CHECK (domain IN ('k12_education', 'art_history', 'bestseller', 'vocational')),
    CONSTRAINT valid_grade_level CHECK (grade_level BETWEEN 1 AND 12 OR grade_level IS NULL)
);

-- 索引
CREATE INDEX idx_workflows_user_id ON workflows(user_id);
CREATE INDEX idx_workflows_domain ON workflows(domain);
CREATE INDEX idx_workflows_status ON workflows(status);
CREATE INDEX idx_workflows_created_at ON workflows(created_at DESC);
CREATE INDEX idx_workflows_completed_at ON workflows(completed_at DESC) WHERE completed_at IS NOT NULL;
CREATE INDEX idx_workflows_subject ON workflows(subject) WHERE subject IS NOT NULL;
CREATE INDEX idx_workflows_activity_type ON workflows(activity_type) WHERE activity_type IS NOT NULL;

-- GIN索引用于JSON查询
CREATE INDEX idx_workflows_context_data ON workflows USING GIN(context_data);
CREATE INDEX idx_workflows_outputs ON workflows USING GIN(outputs);

-- 触发器
CREATE TRIGGER update_workflows_updated_at
    BEFORE UPDATE ON workflows
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 分区策略（按月分区）
-- CREATE TABLE workflows_2025_12 PARTITION OF workflows
--     FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

COMMENT ON TABLE workflows IS '工作流执行记录表';
COMMENT ON COLUMN workflows.duration IS '课程时长（分钟），微课3-15分钟，常规课45-90分钟';
```

### 3. node_executions - 节点执行记录表

```sql
CREATE TABLE node_executions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联工作流
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    
    -- 节点信息
    node_id VARCHAR(255) NOT NULL,
    node_type VARCHAR(100) NOT NULL,  -- objective_generation, content_enrichment, etc.
    
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
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT valid_node_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'skipped'))
);

-- 索引
CREATE INDEX idx_node_executions_workflow_id ON node_executions(workflow_id);
CREATE INDEX idx_node_executions_node_type ON node_executions(node_type);
CREATE INDEX idx_node_executions_status ON node_executions(status);
CREATE INDEX idx_node_executions_created_at ON node_executions(created_at DESC);

-- GIN索引用于JSON查询
CREATE INDEX idx_node_executions_output_data ON node_executions USING GIN(output_data);

-- 复合索引
CREATE INDEX idx_node_executions_workflow_created ON node_executions(workflow_id, created_at DESC);

COMMENT ON TABLE node_executions IS '节点执行记录表';
```

### 4. learning_objectives - 学习目标表

```sql
CREATE TABLE learning_objectives (
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
CREATE INDEX idx_learning_objectives_workflow_id ON learning_objectives(workflow_id);
CREATE INDEX idx_learning_objectives_bloom_level ON learning_objectives(bloom_level);
CREATE INDEX idx_learning_objectives_sequence ON learning_objectives(workflow_id, sequence_order);

-- GIN索引用于数组查询
CREATE INDEX idx_learning_objectives_keywords ON learning_objectives USING GIN(keywords);

COMMENT ON TABLE learning_objectives IS '学习目标表';
```

### 5. teaching_activities - 教学活动表

```sql
CREATE TABLE teaching_activities (
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
        'creative_synthesis'
    ))
);

-- 索引
CREATE INDEX idx_teaching_activities_workflow_id ON teaching_activities(workflow_id);
CREATE INDEX idx_teaching_activities_type ON teaching_activities(activity_type);
CREATE INDEX idx_teaching_activities_created_at ON teaching_activities(created_at DESC);

-- GIN索引
CREATE INDEX idx_teaching_activities_content ON teaching_activities USING GIN(content);
CREATE INDEX idx_teaching_activities_tags ON teaching_activities USING GIN(tags);

-- 触发器
CREATE TRIGGER update_teaching_activities_updated_at
    BEFORE UPDATE ON teaching_activities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE teaching_activities IS '教学活动表';
```

---

## 适配器市场表

### 6. adapters - 适配器表

```sql
CREATE TABLE adapters (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 基本信息
    name VARCHAR(255) NOT NULL UNIQUE,  -- 唯一标识符，如：k12-chemistry-adapter
    display_name VARCHAR(255) NOT NULL,  -- 显示名称
    description TEXT NOT NULL,
    
    -- 开发者信息
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    organization VARCHAR(255),
    
    -- 分类和标签
    domain VARCHAR(50) NOT NULL,  -- k12_education, art_history, etc.
    category VARCHAR(100),  -- science, math, language, etc.
    tags TEXT[] DEFAULT '{}',
    
    -- 仓库信息
    repository_url VARCHAR(500),
    homepage_url VARCHAR(500),
    documentation_url VARCHAR(500),
    
    -- 统计信息
    download_count INTEGER DEFAULT 0,
    star_count INTEGER DEFAULT 0,
    average_rating DECIMAL(3, 2) DEFAULT 0.00,  -- 0.00 - 5.00
    review_count INTEGER DEFAULT 0,
    
    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'draft',  -- draft, pending_review, published, deprecated
    is_official BOOLEAN DEFAULT FALSE,  -- 是否官方适配器
    is_featured BOOLEAN DEFAULT FALSE,  -- 是否精选
    
    -- 最新版本
    latest_version VARCHAR(50),
    
    -- 元数据
    manifest JSONB NOT NULL,  -- manifest.yaml 内容
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    
    -- 约束
    CONSTRAINT valid_adapter_status CHECK (status IN ('draft', 'pending_review', 'published', 'deprecated')),
    CONSTRAINT valid_rating CHECK (average_rating BETWEEN 0 AND 5)
);

-- 索引
CREATE INDEX idx_adapters_name ON adapters(name);
CREATE INDEX idx_adapters_author_id ON adapters(author_id);
CREATE INDEX idx_adapters_domain ON adapters(domain);
CREATE INDEX idx_adapters_status ON adapters(status);
CREATE INDEX idx_adapters_published_at ON adapters(published_at DESC) WHERE published_at IS NOT NULL;
CREATE INDEX idx_adapters_download_count ON adapters(download_count DESC);
CREATE INDEX idx_adapters_rating ON adapters(average_rating DESC);

-- GIN索引
CREATE INDEX idx_adapters_tags ON adapters USING GIN(tags);
CREATE INDEX idx_adapters_manifest ON adapters USING GIN(manifest);

-- 全文搜索索引
CREATE INDEX idx_adapters_fts ON adapters USING GIN(
    to_tsvector('english', display_name || ' ' || COALESCE(description, ''))
);

-- 触发器
CREATE TRIGGER update_adapters_updated_at
    BEFORE UPDATE ON adapters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE adapters IS '适配器表';
COMMENT ON COLUMN adapters.manifest IS '完整的manifest.yaml内容（JSON格式）';
```

### 7. adapter_versions - 适配器版本表

```sql
CREATE TABLE adapter_versions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联适配器
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    
    -- 版本信息
    version VARCHAR(50) NOT NULL,  -- 语义化版本：1.2.3
    release_notes TEXT,
    
    -- 文件信息
    package_url VARCHAR(500) NOT NULL,  -- 下载链接
    package_size BIGINT,  -- 字节
    checksum VARCHAR(64),  -- SHA256
    
    -- 兼容性
    min_platform_version VARCHAR(50),  -- 最低平台版本要求
    max_platform_version VARCHAR(50),
    python_version VARCHAR(50),  -- 如：>=3.9,<4.0
    
    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'draft',  -- draft, published, deprecated
    is_stable BOOLEAN DEFAULT TRUE,
    
    -- 统计
    download_count INTEGER DEFAULT 0,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    
    -- 唯一约束
    CONSTRAINT unique_adapter_version UNIQUE (adapter_id, version),
    CONSTRAINT valid_version_status CHECK (status IN ('draft', 'published', 'deprecated'))
);

-- 索引
CREATE INDEX idx_adapter_versions_adapter_id ON adapter_versions(adapter_id);
CREATE INDEX idx_adapter_versions_version ON adapter_versions(version);
CREATE INDEX idx_adapter_versions_status ON adapter_versions(status);
CREATE INDEX idx_adapter_versions_published_at ON adapter_versions(published_at DESC);

COMMENT ON TABLE adapter_versions IS '适配器版本表';
```

### 8. adapter_dependencies - 适配器依赖关系表

```sql
CREATE TABLE adapter_dependencies (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 依赖关系
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    depends_on_id UUID NOT NULL REFERENCES adapters(id) ON DELETE RESTRICT,
    
    -- 版本约束
    version_constraint VARCHAR(100),  -- 如：>=1.0.0,<2.0.0
    
    -- 依赖类型
    dependency_type VARCHAR(20) NOT NULL DEFAULT 'required',  -- required, optional
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 唯一约束（防止重复依赖）
    CONSTRAINT unique_dependency UNIQUE (adapter_id, depends_on_id),
    CONSTRAINT no_self_dependency CHECK (adapter_id != depends_on_id),
    CONSTRAINT valid_dependency_type CHECK (dependency_type IN ('required', 'optional'))
);

-- 索引
CREATE INDEX idx_adapter_dependencies_adapter_id ON adapter_dependencies(adapter_id);
CREATE INDEX idx_adapter_dependencies_depends_on_id ON adapter_dependencies(depends_on_id);

COMMENT ON TABLE adapter_dependencies IS '适配器依赖关系表';
```

### 9. adapter_reviews - 适配器评论表

```sql
CREATE TABLE adapter_reviews (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- 评分（1-5星）
    rating INTEGER NOT NULL,
    
    -- 评论内容
    title VARCHAR(200),
    content TEXT,
    
    -- 版本
    version VARCHAR(50),  -- 评论针对的版本
    
    -- 有用性投票
    helpful_count INTEGER DEFAULT 0,
    unhelpful_count INTEGER DEFAULT 0,
    
    -- 状态
    is_verified_purchase BOOLEAN DEFAULT FALSE,  -- 是否验证使用过
    status VARCHAR(20) NOT NULL DEFAULT 'published',  -- published, hidden, deleted
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 约束
    CONSTRAINT unique_user_adapter_review UNIQUE (adapter_id, user_id),
    CONSTRAINT valid_rating CHECK (rating BETWEEN 1 AND 5),
    CONSTRAINT valid_review_status CHECK (status IN ('published', 'hidden', 'deleted'))
);

-- 索引
CREATE INDEX idx_adapter_reviews_adapter_id ON adapter_reviews(adapter_id);
CREATE INDEX idx_adapter_reviews_user_id ON adapter_reviews(user_id);
CREATE INDEX idx_adapter_reviews_rating ON adapter_reviews(rating);
CREATE INDEX idx_adapter_reviews_created_at ON adapter_reviews(created_at DESC);
CREATE INDEX idx_adapter_reviews_helpful ON adapter_reviews(helpful_count DESC);

-- 触发器
CREATE TRIGGER update_adapter_reviews_updated_at
    BEFORE UPDATE ON adapter_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 触发器：更新适配器统计
CREATE TRIGGER update_adapter_stats_on_review
    AFTER INSERT OR UPDATE OR DELETE ON adapter_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_adapter_rating_stats();

COMMENT ON TABLE adapter_reviews IS '适配器评论表';
```

### 10. adapter_categories - 适配器分类表

```sql
CREATE TABLE adapter_categories (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 分类信息
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- 层级结构
    parent_id UUID REFERENCES adapter_categories(id) ON DELETE SET NULL,
    level INTEGER DEFAULT 0,
    path VARCHAR(500),  -- 层级路径，如：/science/chemistry
    
    -- 排序和图标
    sort_order INTEGER DEFAULT 0,
    icon VARCHAR(100),
    
    -- 统计
    adapter_count INTEGER DEFAULT 0,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX idx_adapter_categories_parent_id ON adapter_categories(parent_id);
CREATE INDEX idx_adapter_categories_path ON adapter_categories(path);
CREATE INDEX idx_adapter_categories_sort_order ON adapter_categories(sort_order);

COMMENT ON TABLE adapter_categories IS '适配器分类表';
```

### 11. adapter_stats - 适配器统计表

```sql
CREATE TABLE adapter_stats (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联适配器
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    
    -- 日期
    stat_date DATE NOT NULL,
    
    -- 统计指标
    downloads_today INTEGER DEFAULT 0,
    active_users_today INTEGER DEFAULT 0,
    workflows_created_today INTEGER DEFAULT 0,
    
    -- 累计指标
    total_downloads INTEGER DEFAULT 0,
    total_active_users INTEGER DEFAULT 0,
    total_workflows INTEGER DEFAULT 0,
    
    -- 性能指标
    avg_execution_time FLOAT,  -- 平均执行时间（秒）
    avg_success_rate DECIMAL(5, 2),  -- 成功率（百分比）
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 唯一约束
    CONSTRAINT unique_adapter_stat_date UNIQUE (adapter_id, stat_date)
);

-- 索引
CREATE INDEX idx_adapter_stats_adapter_id ON adapter_stats(adapter_id);
CREATE INDEX idx_adapter_stats_date ON adapter_stats(stat_date DESC);
CREATE INDEX idx_adapter_stats_downloads ON adapter_stats(downloads_today DESC);

COMMENT ON TABLE adapter_stats IS '适配器每日统计表';
```

---

## 用户与权限表

### 12. roles - 角色表

```sql
CREATE TABLE roles (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 角色信息
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- 权限级别
    level INTEGER NOT NULL DEFAULT 0,  -- 0=普通用户, 10=开发者, 50=管理员, 100=超级管理员
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_level ON roles(level DESC);

-- 初始数据
INSERT INTO roles (name, display_name, level, description) VALUES
    ('user', 'User', 0, 'Regular user with basic permissions'),
    ('developer', 'Developer', 10, 'Can publish and manage adapters'),
    ('moderator', 'Moderator', 30, 'Can review and moderate content'),
    ('admin', 'Administrator', 50, 'Full system administration'),
    ('superadmin', 'Super Administrator', 100, 'Full system access');

COMMENT ON TABLE roles IS '角色表';
```

### 13. permissions - 权限表

```sql
CREATE TABLE permissions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 权限信息
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL,  -- workflow, adapter, user, etc.
    action VARCHAR(50) NOT NULL,  -- create, read, update, delete, publish, etc.
    description TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);

COMMENT ON TABLE permissions IS '权限表';
```

### 14. role_permissions - 角色权限关联表

```sql
CREATE TABLE role_permissions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 唯一约束
    CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id)
);

-- 索引
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

COMMENT ON TABLE role_permissions IS '角色权限关联表';
```

### 15. user_roles - 用户角色关联表

```sql
CREATE TABLE user_roles (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,  -- 角色过期时间（可选）
    
    -- 唯一约束
    CONSTRAINT unique_user_role UNIQUE (user_id, role_id)
);

-- 索引
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

COMMENT ON TABLE user_roles IS '用户角色关联表';
```

### 16. user_subscriptions - 用户订阅表

```sql
CREATE TABLE user_subscriptions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    
    -- 订阅类型
    notification_type VARCHAR(50) NOT NULL DEFAULT 'all',  -- all, major_updates, security_updates
    
    -- 状态
    is_active BOOLEAN DEFAULT TRUE,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- 唯一约束
    CONSTRAINT unique_user_subscription UNIQUE (user_id, adapter_id)
);

-- 索引
CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX idx_user_subscriptions_adapter_id ON user_subscriptions(adapter_id);

COMMENT ON TABLE user_subscriptions IS '用户订阅表（关注适配器更新）';
```

---

## 统计与日志表

### 17. audit_logs - 审计日志表

```sql
CREATE TABLE audit_logs (
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- 分区策略（按月分区）
-- CREATE TABLE audit_logs_2025_12 PARTITION OF audit_logs
--     FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

COMMENT ON TABLE audit_logs IS '审计日志表';
```

### 18. api_logs - API日志表

```sql
CREATE TABLE api_logs (
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status_code ON api_logs(status_code);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_response_time ON api_logs(response_time DESC);

-- 分区策略（按月分区）
-- CREATE TABLE api_logs_2025_12 PARTITION OF api_logs
--     FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

COMMENT ON TABLE api_logs IS 'API访问日志表';
```

---

## 索引设计

### 索引策略

```yaml
原则:
  - 主键索引: 所有表使用UUID主键，自动创建B-tree索引
  - 外键索引: 所有外键字段创建索引，提升JOIN性能
  - 查询索引: 基于常见查询模式创建复合索引
  - JSON索引: 使用GIN索引支持JSON查询
  - 全文索引: 使用tsvector支持全文搜索
  - 分区索引: 大表按时间分区，每个分区独立索引
```

### 关键索引列表

| 表名 | 索引类型 | 字段 | 说明 |
|------|---------|------|------|
| workflows | B-tree | (user_id, created_at DESC) | 用户工作流列表 |
| workflows | B-tree | (status, created_at DESC) | 状态过滤 |
| workflows | GIN | context_data | JSON查询 |
| node_executions | B-tree | (workflow_id, created_at) | 工作流节点查询 |
| adapters | GIN | tags | 标签搜索 |
| adapters | GIN | tsvector | 全文搜索 |
| adapter_reviews | B-tree | (adapter_id, rating DESC) | 评分排序 |
| audit_logs | B-tree | (resource_type, resource_id) | 资源审计查询 |

### 索引维护

```sql
-- 定期重建索引（REINDEX）
REINDEX TABLE workflows;

-- 分析表统计信息
ANALYZE workflows;

-- 清理无用索引
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY schemaname, tablename;

-- 查看索引大小
SELECT
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;
```

---

## 数据迁移

### Alembic配置

**1. 安装Alembic**

```bash
pip install alembic
alembic init alembic
```

**2. 配置 alembic.ini**

```ini
[alembic]
script_location = alembic
sqlalchemy.url = postgresql://user:password@localhost:5432/metaworkflow

[loggers]
keys = root,sqlalchemy,alembic

[handlers]
keys = console

[formatters]
keys = generic
```

**3. 创建迁移脚本**

```bash
# 创建迁移
alembic revision --autogenerate -m "Initial schema"

# 应用迁移
alembic upgrade head

# 回滚迁移
alembic downgrade -1
```

**4. 示例迁移脚本**

```python
"""Initial schema

Revision ID: 001
Create Date: 2025-12-09
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers
revision = '001'
down_revision = None
branch_labels = None
depends_on = None

def upgrade():
    # 创建更新时间戳函数
    op.execute("""
        CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = NOW();
            RETURN NEW;
        END;
        $$ language 'plpgsql';
    """)
    
    # 创建users表
    op.create_table('users',
        sa.Column('id', postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column('username', sa.String(100), nullable=False, unique=True),
        sa.Column('email', sa.String(255), nullable=False, unique=True),
        # ... 其他字段
    )
    
    # 创建workflows表
    op.create_table('workflows',
        sa.Column('id', postgresql.UUID(as_uuid=True), primary_key=True),
        # ... 其他字段
    )

def downgrade():
    op.drop_table('workflows')
    op.drop_table('users')
    op.execute('DROP FUNCTION IF EXISTS update_updated_at_column()')
```

### 数据迁移脚本

```python
# scripts/migrate_data.py
"""
数据迁移工具
用于从V1.0迁移到V2.0
"""
import asyncio
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

async def migrate_workflows():
    """迁移工作流数据"""
    # 1. 读取V1.0数据
    # 2. 转换数据结构
    # 3. 写入V2.0表
    pass

async def migrate_adapters():
    """迁移适配器数据"""
    pass

if __name__ == "__main__":
    asyncio.run(migrate_workflows())
    asyncio.run(migrate_adapters())
```

---

## 性能优化

### 1. 查询优化

```sql
-- 使用EXPLAIN ANALYZE分析查询
EXPLAIN ANALYZE
SELECT w.*, COUNT(n.id) as node_count
FROM workflows w
LEFT JOIN node_executions n ON w.id = n.workflow_id
WHERE w.user_id = 'xxx' AND w.status = 'completed'
GROUP BY w.id
ORDER BY w.created_at DESC
LIMIT 20;

-- 优化建议：
-- 1. 添加复合索引 (user_id, status, created_at)
-- 2. 使用物化视图缓存统计数据
```

### 2. 连接池配置

```python
# config/database.py
from sqlalchemy.pool import QueuePool

engine = create_engine(
    DATABASE_URL,
    poolclass=QueuePool,
    pool_size=20,          # 连接池大小
    max_overflow=10,       # 最大溢出连接
    pool_timeout=30,       # 获取连接超时
    pool_recycle=3600,     # 连接回收时间（秒）
    echo=False,
    echo_pool=False
)
```

### 3. 分区策略

```sql
-- 按月分区workflows表
CREATE TABLE workflows (
    -- 字段定义...
) PARTITION BY RANGE (created_at);

-- 创建分区
CREATE TABLE workflows_2025_12 PARTITION OF workflows
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE workflows_2026_01 PARTITION OF workflows
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- 自动创建分区的函数
CREATE OR REPLACE FUNCTION create_monthly_partitions(
    table_name TEXT,
    start_date DATE,
    end_date DATE
)
RETURNS VOID AS $$
DECLARE
    current_date DATE := start_date;
    partition_name TEXT;
    start_range DATE;
    end_range DATE;
BEGIN
    WHILE current_date < end_date LOOP
        partition_name := table_name || '_' || to_char(current_date, 'YYYY_MM');
        start_range := current_date;
        end_range := current_date + INTERVAL '1 month';
        
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
            partition_name, table_name, start_range, end_range
        );
        
        current_date := end_range;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 使用函数创建一年的分区
SELECT create_monthly_partitions('workflows', '2025-01-01'::date, '2026-01-01'::date);
```

### 4. 物化视图

```sql
-- 创建适配器统计物化视图
CREATE MATERIALIZED VIEW adapter_stats_summary AS
SELECT 
    a.id,
    a.name,
    a.download_count,
    a.average_rating,
    COUNT(DISTINCT r.id) as review_count,
    COUNT(DISTINCT w.id) as workflow_count,
    AVG(w.execution_time) as avg_execution_time
FROM adapters a
LEFT JOIN adapter_reviews r ON a.id = r.adapter_id
LEFT JOIN workflows w ON a.domain = w.domain
GROUP BY a.id, a.name, a.download_count, a.average_rating;

-- 创建索引
CREATE INDEX idx_adapter_stats_summary_id ON adapter_stats_summary(id);

-- 定期刷新（使用pg_cron或后台任务）
REFRESH MATERIALIZED VIEW CONCURRENTLY adapter_stats_summary;
```

### 5. 缓存策略

```python
# Redis缓存配置
CACHE_CONFIG = {
    # 工作流缓存
    'workflow': {
        'ttl': 3600,  # 1小时
        'prefix': 'workflow:'
    },
    
    # 适配器缓存
    'adapter': {
        'ttl': 7200,  # 2小时
        'prefix': 'adapter:'
    },
    
    # 查询结果缓存
    'query': {
        'ttl': 300,  # 5分钟
        'prefix': 'query:'
    }
}
```

---

## 备份策略

### 1. 全量备份

```bash
#!/bin/bash
# scripts/backup_full.sh

# 配置
BACKUP_DIR="/var/backups/metaworkflow"
DB_NAME="metaworkflow"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行全量备份
pg_dump -Fc -v \
    -h localhost \
    -U postgres \
    -d $DB_NAME \
    -f $BACKUP_DIR/metaworkflow_full_$DATE.dump

# 压缩备份
gzip $BACKUP_DIR/metaworkflow_full_$DATE.dump

# 删除7天前的备份
find $BACKUP_DIR -name "*.dump.gz" -mtime +7 -delete

echo "Backup completed: metaworkflow_full_$DATE.dump.gz"
```

### 2. 增量备份（WAL归档）

```bash
# postgresql.conf配置
wal_level = replica
archive_mode = on
archive_command = 'cp %p /var/lib/postgresql/archive/%f'
max_wal_senders = 3
```

### 3. 恢复脚本

```bash
#!/bin/bash
# scripts/restore.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: ./restore.sh <backup_file>"
    exit 1
fi

# 停止应用
systemctl stop metaworkflow

# 删除现有数据库
dropdb metaworkflow

# 创建新数据库
createdb metaworkflow

# 恢复备份
pg_restore -v \
    -h localhost \
    -U postgres \
    -d metaworkflow \
    $BACKUP_FILE

# 启动应用
systemctl start metaworkflow

echo "Restore completed"
```

### 4. 定时任务

```bash
# crontab -e

# 每天凌晨2点全量备份
0 2 * * * /path/to/scripts/backup_full.sh

# 每小时增量备份
0 * * * * /path/to/scripts/backup_incremental.sh

# 每周日备份到远程存储
0 3 * * 0 /path/to/scripts/backup_to_s3.sh
```

---

## 监控和维护

### 1. 监控指标

```sql
-- 数据库大小
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
ORDER BY pg_database_size(pg_database.datname) DESC;

-- 表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 慢查询
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 20;

-- 锁等待
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocking_locks.pid AS blocking_pid,
    blocked_activity.query AS blocked_query,
    blocking_activity.query AS blocking_query
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;
```

### 2. 定期维护任务

```sql
-- VACUUM清理
VACUUM ANALYZE workflows;

-- 重建索引
REINDEX TABLE workflows;

-- 更新统计信息
ANALYZE;

-- 清理旧日志（保留90天）
DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '90 days';
DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL '30 days';
```

---

## 附录

### A. 完整建表脚本

完整的SQL脚本请参考：
- `sql/schema.sql` - 完整表结构
- `sql/indexes.sql` - 所有索引
- `sql/functions.sql` - 触发器和函数
- `sql/seed.sql` - 初始数据

### B. 性能基准测试

```bash
# 使用pgbench进行性能测试
pgbench -i -s 50 metaworkflow  # 初始化测试数据
pgbench -c 10 -j 2 -t 1000 metaworkflow  # 运行测试
```

### C. 数据字典

完整的数据字典文档：`docs/DATA_DICTIONARY.md`

---

## 总结

### 数据库特点

✅ **完整性**: 18张表覆盖所有业务场景  
✅ **规范化**: 遵循3NF，减少冗余  
✅ **可扩展**: JSON字段支持灵活扩展  
✅ **高性能**: 合理索引和分区策略  
✅ **可追溯**: 完整的审计日志  
✅ **安全性**: 约束和触发器保证数据一致性  

### 下一步

1. ✅ 数据库设计完成
2. ⏳ 创建Alembic迁移脚本
3. ⏳ 编写ORM仓储层
4. ⏳ 实施备份策略
5. ⏳ 配置监控告警

---

**文档版本**: 1.0  
**最后更新**: 2025-12-09  
**维护者**: MetaWorkflow Team
