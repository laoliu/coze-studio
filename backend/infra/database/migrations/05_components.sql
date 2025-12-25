-- =============================================================================
-- Simplified Adapter Design - Components Table Migration
-- =============================================================================
-- 创建时间: 2025-01-XX
-- 说明: 创建 components 表以支持简化的适配器设计
-- 依赖: 00_init_functions.sql, 01_users_and_permissions.sql
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 表: components - 组件表（统一管理插件、适配器等组件）
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS components (
    -- 主键
    component_id VARCHAR(255) PRIMARY KEY,  -- 唯一标识，如：k12.math_problem_generator

    -- 基本信息
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    version VARCHAR(50) NOT NULL,

    -- 类型和分类
    type VARCHAR(50) NOT NULL,  -- plugin, adapter, builtin, external
    category VARCHAR(100) NOT NULL,  -- domain, visualization, integration, etc.

    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'active',  -- active, inactive, deprecated, draft

    -- 作者信息
    author_id BIGINT NOT NULL,
    author_name VARCHAR(255),

    -- 配置（JSON）
    config JSONB DEFAULT '{}',
    tags JSONB DEFAULT '[]',

    -- 元数据
    icon VARCHAR(500),
    homepage VARCHAR(500),
    repository VARCHAR(500),
    documentation VARCHAR(500),
    license VARCHAR(100) DEFAULT 'MIT',

    -- 统计
    rating DECIMAL(3, 2) DEFAULT 0.0,
    install_count INTEGER DEFAULT 0,
    usage_count BIGINT DEFAULT 0,

    -- 标志
    is_official BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,

    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    deprecated_at TIMESTAMP WITH TIME ZONE,

    -- 约束
    CONSTRAINT valid_component_type CHECK (type IN ('plugin', 'adapter', 'builtin', 'external')),
    CONSTRAINT valid_component_status CHECK (status IN ('active', 'inactive', 'deprecated', 'draft')),
    CONSTRAINT valid_rating CHECK (rating >= 0.0 AND rating <= 5.0)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_components_type ON components(type);
CREATE INDEX IF NOT EXISTS idx_components_category ON components(category);
CREATE INDEX IF NOT EXISTS idx_components_status ON components(status);
CREATE INDEX IF NOT EXISTS idx_components_author_id ON components(author_id);
CREATE INDEX IF NOT EXISTS idx_components_rating ON components(rating DESC);
CREATE INDEX IF NOT EXISTS idx_components_install_count ON components(install_count DESC);
CREATE INDEX IF NOT EXISTS idx_components_usage_count ON components(usage_count DESC);
CREATE INDEX IF NOT EXISTS idx_components_created_at ON components(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_components_is_official ON components(is_official) WHERE is_official = TRUE;
CREATE INDEX IF NOT EXISTS idx_components_is_featured ON components(is_featured) WHERE is_featured = TRUE;

-- GIN 索引用于 JSON 查询
CREATE INDEX IF NOT EXISTS idx_components_config ON components USING GIN(config);
CREATE INDEX IF NOT EXISTS idx_components_tags ON components USING GIN(tags);

-- 全文搜索索引
CREATE INDEX IF NOT EXISTS idx_components_search ON components
    USING GIN(to_tsvector('english',
        coalesce(name, '') || ' ' ||
        coalesce(display_name, '') || ' ' ||
        coalesce(description, '')
    ));

-- 组合索引用于常见查询
CREATE INDEX IF NOT EXISTS idx_components_type_category_status ON components(type, category, status);
CREATE INDEX IF NOT EXISTS idx_components_type_rating ON components(type, rating DESC) WHERE status = 'active';

-- 触发器
DROP TRIGGER IF EXISTS update_components_updated_at ON components;
CREATE TRIGGER update_components_updated_at
    BEFORE UPDATE ON components
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE components IS '组件表：统一管理插件、适配器等可扩展组件';
COMMENT ON COLUMN components.component_id IS '组件唯一标识，格式：domain.component_name';
COMMENT ON COLUMN components.type IS '组件类型：plugin（插件）、adapter（适配器）、builtin（内置）、external（外部）';
COMMENT ON COLUMN components.category IS '组件分类：domain（领域）、visualization（可视化）、integration（集成）等';
COMMENT ON COLUMN components.config IS '组件配置（JSON格式），包含 adapter_config、plugin_id 等';
COMMENT ON COLUMN components.tags IS '标签数组（JSON格式），用于搜索和分类';

-- -----------------------------------------------------------------------------
-- 表: component_installations - 组件安装记录表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS component_installations (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- 关联
    component_id VARCHAR(255) NOT NULL REFERENCES components(component_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,

    -- 安装信息
    installed_version VARCHAR(50) NOT NULL,
    installation_status VARCHAR(20) DEFAULT 'active',  -- active, uninstalled, disabled

    -- 配置（用户自定义）
    user_config JSONB DEFAULT '{}',

    -- 使用统计
    usage_count BIGINT DEFAULT 0,
    last_used_at TIMESTAMP WITH TIME ZONE,

    -- 时间戳
    installed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    uninstalled_at TIMESTAMP WITH TIME ZONE,

    -- 约束
    CONSTRAINT valid_installation_status CHECK (installation_status IN ('active', 'uninstalled', 'disabled')),
    CONSTRAINT unique_user_component UNIQUE (component_id, user_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_component_installations_component_id ON component_installations(component_id);
CREATE INDEX IF NOT EXISTS idx_component_installations_user_id ON component_installations(user_id);
CREATE INDEX IF NOT EXISTS idx_component_installations_status ON component_installations(installation_status);
CREATE INDEX IF NOT EXISTS idx_component_installations_installed_at ON component_installations(installed_at DESC);

-- GIN 索引用于用户配置查询
CREATE INDEX IF NOT EXISTS idx_component_installations_user_config ON component_installations USING GIN(user_config);

-- 触发器
DROP TRIGGER IF EXISTS update_component_installations_updated_at ON component_installations;
CREATE TRIGGER update_component_installations_updated_at
    BEFORE UPDATE ON component_installations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE component_installations IS '组件安装记录表：记录用户安装的组件';
COMMENT ON COLUMN component_installations.user_config IS '用户自定义配置（JSON格式）';

-- -----------------------------------------------------------------------------
-- 表: component_reviews - 组件评价表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS component_reviews (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- 关联
    component_id VARCHAR(255) NOT NULL REFERENCES components(component_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,

    -- 评价内容
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    content TEXT,

    -- 状态
    status VARCHAR(20) DEFAULT 'published',  -- published, hidden, deleted

    -- 有用性统计
    helpful_count INTEGER DEFAULT 0,
    unhelpful_count INTEGER DEFAULT 0,

    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,

    -- 约束
    CONSTRAINT valid_review_status CHECK (status IN ('published', 'hidden', 'deleted')),
    CONSTRAINT unique_user_component_review UNIQUE (component_id, user_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_component_reviews_component_id ON component_reviews(component_id);
CREATE INDEX IF NOT EXISTS idx_component_reviews_user_id ON component_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_component_reviews_rating ON component_reviews(rating);
CREATE INDEX IF NOT EXISTS idx_component_reviews_status ON component_reviews(status);
CREATE INDEX IF NOT EXISTS idx_component_reviews_created_at ON component_reviews(created_at DESC);

-- 触发器
DROP TRIGGER IF EXISTS update_component_reviews_updated_at ON component_reviews;
CREATE TRIGGER update_component_reviews_updated_at
    BEFORE UPDATE ON component_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 触发器：更新组件评分
CREATE OR REPLACE FUNCTION update_component_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE components
    SET rating = (
        SELECT COALESCE(AVG(rating), 0.0)
        FROM component_reviews
        WHERE component_id = NEW.component_id
          AND status = 'published'
    )
    WHERE component_id = NEW.component_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_component_rating ON component_reviews;
CREATE TRIGGER trigger_update_component_rating
    AFTER INSERT OR UPDATE OR DELETE ON component_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_component_rating();

COMMENT ON TABLE component_reviews IS '组件评价表：用户对组件的评分和评论';

-- =============================================================================
-- 初始数据
-- =============================================================================

-- 插入示例适配器（如果不存在）
INSERT INTO components (component_id, name, display_name, description, version, type, category, author_id, author_name, status, is_official)
VALUES
    ('k12.math_problem_generator', 'math_problem_generator', 'K12数学题目生成器', '根据知识点和难度生成数学问题', '1.0.0', 'adapter', 'domain', 0, 'System', 'active', TRUE),
    ('art.painting_analyzer', 'painting_analyzer', '艺术史绘画分析器', '分析绘画作品的风格、技法和历史背景', '1.0.0', 'adapter', 'domain', 0, 'System', 'active', TRUE),
    ('vocational.skill_assessor', 'skill_assessor', '职业技能评估器', '评估学员的职业技能水平', '1.0.0', 'adapter', 'domain', 0, 'System', 'active', TRUE)
ON CONFLICT (component_id) DO NOTHING;
