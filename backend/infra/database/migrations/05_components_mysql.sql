-- =============================================================================
-- Simplified Adapter Design - Components Table Migration (MySQL Version)
-- =============================================================================
-- 创建时间: 2025-12-26
-- 说明: 创建 components 表以支持简化的适配器设计 (MySQL 兼容版本)
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 表: components - 组件表（统一管理插件、适配器等组件）
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS components (
    -- 主键
    component_id VARCHAR(255) PRIMARY KEY,

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
    config JSON DEFAULT NULL,
    tags JSON DEFAULT NULL,

    -- 元数据
    icon VARCHAR(500),
    homepage VARCHAR(500),
    repository VARCHAR(500),
    documentation VARCHAR(500),
    license VARCHAR(100) DEFAULT 'MIT',

    -- 统计
    rating DECIMAL(3, 2) DEFAULT 0.0,
    install_count INT DEFAULT 0,
    usage_count BIGINT DEFAULT 0,

    -- 标志
    is_official BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,

    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    published_at TIMESTAMP NULL,
    deprecated_at TIMESTAMP NULL,

    -- 约束
    CONSTRAINT valid_component_type CHECK (type IN ('plugin', 'adapter', 'builtin', 'external')),
    CONSTRAINT valid_component_status CHECK (status IN ('active', 'inactive', 'deprecated', 'draft')),
    CONSTRAINT valid_rating CHECK (rating >= 0.0 AND rating <= 5.0),

    -- 索引
    INDEX idx_components_type (type),
    INDEX idx_components_category (category),
    INDEX idx_components_status (status),
    INDEX idx_components_author_id (author_id),
    INDEX idx_components_rating (rating DESC),
    INDEX idx_components_install_count (install_count DESC),
    INDEX idx_components_usage_count (usage_count DESC),
    INDEX idx_components_created_at (created_at DESC),
    INDEX idx_components_is_official (is_official),
    INDEX idx_components_is_featured (is_featured),
    INDEX idx_components_type_category_status (type, category, status),
    INDEX idx_components_type_rating (type, rating DESC),

    -- 全文搜索索引
    FULLTEXT INDEX idx_components_search (name, display_name, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 表: component_installations - 组件安装记录表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS component_installations (
    -- 主键
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    -- 关联
    component_id VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,

    -- 安装信息
    installed_version VARCHAR(50) NOT NULL,
    installation_status VARCHAR(20) DEFAULT 'active',  -- active, uninstalled, disabled

    -- 配置（用户自定义）
    user_config JSON DEFAULT NULL,

    -- 使用统计
    usage_count BIGINT DEFAULT 0,
    last_used_at TIMESTAMP NULL,

    -- 时间戳
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    uninstalled_at TIMESTAMP NULL,

    -- 约束
    CONSTRAINT valid_installation_status CHECK (installation_status IN ('active', 'uninstalled', 'disabled')),
    CONSTRAINT unique_user_component UNIQUE (component_id, user_id),

    -- 外键
    CONSTRAINT fk_component_installations_component FOREIGN KEY (component_id)
        REFERENCES components(component_id) ON DELETE CASCADE,

    -- 索引
    INDEX idx_component_installations_component_id (component_id),
    INDEX idx_component_installations_user_id (user_id),
    INDEX idx_component_installations_status (installation_status),
    INDEX idx_component_installations_installed_at (installed_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试数据（K12 教育适配器）
INSERT INTO components (
    component_id,
    name,
    display_name,
    description,
    version,
    type,
    category,
    status,
    author_id,
    author_name,
    config,
    tags,
    is_official,
    is_public
) VALUES (
    'k12_education',
    'k12_education',
    'K-12 Education Adapter',
    'Comprehensive K-12 education content generation adapter supporting multiple subjects and grade levels',
    '1.0.0',
    'adapter',
    'domain',
    'active',
    1,
    'System',
    '{"adapter_config": {"max_tokens": 2000, "temperature": 0.7}}',
    '["education", "k12", "teaching", "curriculum"]',
    TRUE,
    TRUE
);
