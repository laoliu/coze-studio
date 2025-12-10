-- =============================================================================
-- MetaWorkflow V2.0 - 适配器市场表
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建适配器、版本、依赖、评论、统计等相关表
-- 依赖: 00_init_functions.sql, 01_users_and_permissions.sql
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 表: adapter_categories - 适配器分类表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapter_categories (
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
CREATE INDEX IF NOT EXISTS idx_adapter_categories_parent_id ON adapter_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_adapter_categories_path ON adapter_categories(path);
CREATE INDEX IF NOT EXISTS idx_adapter_categories_sort_order ON adapter_categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_adapter_categories_name ON adapter_categories(name);

-- 触发器
DROP TRIGGER IF EXISTS update_adapter_categories_updated_at ON adapter_categories;
CREATE TRIGGER update_adapter_categories_updated_at
    BEFORE UPDATE ON adapter_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE adapter_categories IS '适配器分类表';
COMMENT ON COLUMN adapter_categories.path IS '层级路径，用于快速查询子分类';

-- 初始化分类数据
INSERT INTO adapter_categories (name, display_name, description, level, path, sort_order) VALUES
    ('education', 'Education', 'Educational domain adapters', 0, '/education', 1),
    ('k12', 'K-12 Education', 'Primary and secondary education', 1, '/education/k12', 1),
    ('higher_education', 'Higher Education', 'University and college education', 1, '/education/higher_education', 2),
    ('vocational', 'Vocational Training', 'Professional skills training', 1, '/education/vocational', 3),
    ('content', 'Content Creation', 'Content generation and creation', 0, '/content', 2),
    ('art', 'Art & Design', 'Art history and design', 1, '/content/art', 1),
    ('literature', 'Literature', 'Books and literature analysis', 1, '/content/literature', 2)
ON CONFLICT (name) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 表: adapters - 适配器表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapters (
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
    category_id UUID REFERENCES adapter_categories(id) ON DELETE SET NULL,
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
    CONSTRAINT valid_rating CHECK (average_rating BETWEEN 0 AND 5),
    CONSTRAINT valid_domain CHECK (domain IN ('k12_education', 'art_history', 'bestseller', 'vocational', 'higher_education', 'professional_training'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_adapters_name ON adapters(name);
CREATE INDEX IF NOT EXISTS idx_adapters_author_id ON adapters(author_id);
CREATE INDEX IF NOT EXISTS idx_adapters_domain ON adapters(domain);
CREATE INDEX IF NOT EXISTS idx_adapters_category_id ON adapters(category_id);
CREATE INDEX IF NOT EXISTS idx_adapters_status ON adapters(status);
CREATE INDEX IF NOT EXISTS idx_adapters_published_at ON adapters(published_at DESC) WHERE published_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_adapters_download_count ON adapters(download_count DESC);
CREATE INDEX IF NOT EXISTS idx_adapters_rating ON adapters(average_rating DESC);
CREATE INDEX IF NOT EXISTS idx_adapters_is_featured ON adapters(is_featured) WHERE is_featured = TRUE;

-- GIN索引
CREATE INDEX IF NOT EXISTS idx_adapters_tags ON adapters USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_adapters_manifest ON adapters USING GIN(manifest);

-- 全文搜索索引
CREATE INDEX IF NOT EXISTS idx_adapters_fts ON adapters USING GIN(
    to_tsvector('english', display_name || ' ' || COALESCE(description, ''))
);

-- 触发器
DROP TRIGGER IF EXISTS update_adapters_updated_at ON adapters;
CREATE TRIGGER update_adapters_updated_at
    BEFORE UPDATE ON adapters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE adapters IS '适配器表';
COMMENT ON COLUMN adapters.name IS '唯一标识符，如：k12-chemistry-adapter';
COMMENT ON COLUMN adapters.manifest IS '完整的manifest.yaml内容（JSON格式）';
COMMENT ON COLUMN adapters.is_official IS '是否官方适配器';
COMMENT ON COLUMN adapters.is_featured IS '是否精选适配器';

-- -----------------------------------------------------------------------------
-- 表: adapter_versions - 适配器版本表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapter_versions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 关联适配器
    adapter_id UUID NOT NULL REFERENCES adapters(id) ON DELETE CASCADE,
    
    -- 版本信息
    version VARCHAR(50) NOT NULL,  -- 语义化版本：1.2.3
    release_notes TEXT,
    changelog TEXT,
    
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
CREATE INDEX IF NOT EXISTS idx_adapter_versions_adapter_id ON adapter_versions(adapter_id);
CREATE INDEX IF NOT EXISTS idx_adapter_versions_version ON adapter_versions(version);
CREATE INDEX IF NOT EXISTS idx_adapter_versions_status ON adapter_versions(status);
CREATE INDEX IF NOT EXISTS idx_adapter_versions_published_at ON adapter_versions(published_at DESC) WHERE published_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_adapter_versions_download_count ON adapter_versions(download_count DESC);

COMMENT ON TABLE adapter_versions IS '适配器版本表';
COMMENT ON COLUMN adapter_versions.version IS '语义化版本号，如：1.2.3';
COMMENT ON COLUMN adapter_versions.checksum IS 'SHA256校验和';

-- -----------------------------------------------------------------------------
-- 表: adapter_dependencies - 适配器依赖关系表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapter_dependencies (
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
CREATE INDEX IF NOT EXISTS idx_adapter_dependencies_adapter_id ON adapter_dependencies(adapter_id);
CREATE INDEX IF NOT EXISTS idx_adapter_dependencies_depends_on_id ON adapter_dependencies(depends_on_id);

COMMENT ON TABLE adapter_dependencies IS '适配器依赖关系表';
COMMENT ON COLUMN adapter_dependencies.version_constraint IS '版本约束，如：>=1.0.0,<2.0.0';

-- -----------------------------------------------------------------------------
-- 表: adapter_reviews - 适配器评论表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapter_reviews (
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
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_adapter_id ON adapter_reviews(adapter_id);
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_user_id ON adapter_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_rating ON adapter_reviews(rating);
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_created_at ON adapter_reviews(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_helpful ON adapter_reviews(helpful_count DESC);
CREATE INDEX IF NOT EXISTS idx_adapter_reviews_status ON adapter_reviews(status);

-- 触发器
DROP TRIGGER IF EXISTS update_adapter_reviews_updated_at ON adapter_reviews;
CREATE TRIGGER update_adapter_reviews_updated_at
    BEFORE UPDATE ON adapter_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 触发器：更新适配器统计
DROP TRIGGER IF EXISTS update_adapter_stats_on_review ON adapter_reviews;
CREATE TRIGGER update_adapter_stats_on_review
    AFTER INSERT OR UPDATE OR DELETE ON adapter_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_adapter_rating_stats();

COMMENT ON TABLE adapter_reviews IS '适配器评论表';
COMMENT ON COLUMN adapter_reviews.is_verified_purchase IS '是否验证使用过该适配器';

-- -----------------------------------------------------------------------------
-- 表: adapter_stats - 适配器统计表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS adapter_stats (
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
CREATE INDEX IF NOT EXISTS idx_adapter_stats_adapter_id ON adapter_stats(adapter_id);
CREATE INDEX IF NOT EXISTS idx_adapter_stats_date ON adapter_stats(stat_date DESC);
CREATE INDEX IF NOT EXISTS idx_adapter_stats_downloads ON adapter_stats(downloads_today DESC);

COMMENT ON TABLE adapter_stats IS '适配器每日统计表';

-- -----------------------------------------------------------------------------
-- 表: user_subscriptions - 用户订阅表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_subscriptions (
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
    CONSTRAINT unique_user_subscription UNIQUE (user_id, adapter_id),
    CONSTRAINT valid_notification_type CHECK (notification_type IN ('all', 'major_updates', 'security_updates', 'none'))
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_adapter_id ON user_subscriptions(adapter_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_is_active ON user_subscriptions(is_active) WHERE is_active = TRUE;

-- 触发器
DROP TRIGGER IF EXISTS update_user_subscriptions_updated_at ON user_subscriptions;
CREATE TRIGGER update_user_subscriptions_updated_at
    BEFORE UPDATE ON user_subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE user_subscriptions IS '用户订阅表（关注适配器更新）';

-- 打印创建完成信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 适配器市场表创建完成';
    RAISE NOTICE '   - adapter_categories: 适配器分类表（已初始化基础分类）';
    RAISE NOTICE '   - adapters: 适配器表';
    RAISE NOTICE '   - adapter_versions: 适配器版本表';
    RAISE NOTICE '   - adapter_dependencies: 适配器依赖关系表';
    RAISE NOTICE '   - adapter_reviews: 适配器评论表';
    RAISE NOTICE '   - adapter_stats: 适配器统计表';
    RAISE NOTICE '   - user_subscriptions: 用户订阅表';
END $$;
