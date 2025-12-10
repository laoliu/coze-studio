-- =============================================================================
-- MetaWorkflow V2.0 - 用户与权限表
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建用户、角色、权限相关表
-- 依赖: 00_init_functions.sql
-- =============================================================================

-- 删除已存在的表（开发环境，生产环境需谨慎）
-- DROP TABLE IF EXISTS user_roles CASCADE;
-- DROP TABLE IF EXISTS role_permissions CASCADE;
-- DROP TABLE IF EXISTS user_subscriptions CASCADE;
-- DROP TABLE IF EXISTS permissions CASCADE;
-- DROP TABLE IF EXISTS roles CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

-- -----------------------------------------------------------------------------
-- 表: users - 用户表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
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
    
    -- 约束
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- 触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.is_developer IS '是否为开发者（可发布适配器）';
COMMENT ON COLUMN users.preferences IS '用户偏好设置（JSON格式）';

-- -----------------------------------------------------------------------------
-- 表: roles - 角色表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS roles (
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
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_level ON roles(level DESC);

-- 触发器
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.level IS '权限级别，数值越大权限越高';

-- 初始化角色数据
INSERT INTO roles (name, display_name, level, description) VALUES
    ('user', 'User', 0, 'Regular user with basic permissions'),
    ('developer', 'Developer', 10, 'Can publish and manage adapters'),
    ('moderator', 'Moderator', 30, 'Can review and moderate content'),
    ('admin', 'Administrator', 50, 'Full system administration'),
    ('superadmin', 'Super Administrator', 100, 'Full system access')
ON CONFLICT (name) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 表: permissions - 权限表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS permissions (
    -- 主键
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 权限信息
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL,  -- workflow, adapter, user, etc.
    action VARCHAR(50) NOT NULL,    -- create, read, update, delete, publish, etc.
    description TEXT,
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions(action);
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);

COMMENT ON TABLE permissions IS '权限表';
COMMENT ON COLUMN permissions.resource IS '资源类型：workflow, adapter, user等';
COMMENT ON COLUMN permissions.action IS '操作类型：create, read, update, delete, publish等';

-- 初始化基础权限
INSERT INTO permissions (name, resource, action, description) VALUES
    -- 工作流权限
    ('workflow.create', 'workflow', 'create', 'Create workflows'),
    ('workflow.read', 'workflow', 'read', 'View workflows'),
    ('workflow.update', 'workflow', 'update', 'Update workflows'),
    ('workflow.delete', 'workflow', 'delete', 'Delete workflows'),
    ('workflow.execute', 'workflow', 'execute', 'Execute workflows'),
    
    -- 适配器权限
    ('adapter.create', 'adapter', 'create', 'Create adapters'),
    ('adapter.read', 'adapter', 'read', 'View adapters'),
    ('adapter.update', 'adapter', 'update', 'Update adapters'),
    ('adapter.delete', 'adapter', 'delete', 'Delete adapters'),
    ('adapter.publish', 'adapter', 'publish', 'Publish adapters'),
    ('adapter.review', 'adapter', 'review', 'Review adapters'),
    
    -- 用户权限
    ('user.create', 'user', 'create', 'Create users'),
    ('user.read', 'user', 'read', 'View users'),
    ('user.update', 'user', 'update', 'Update users'),
    ('user.delete', 'user', 'delete', 'Delete users'),
    ('user.manage_roles', 'user', 'manage_roles', 'Manage user roles'),
    
    -- 系统权限
    ('system.admin', 'system', 'admin', 'System administration'),
    ('system.moderate', 'system', 'moderate', 'Content moderation')
ON CONFLICT (name) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 表: role_permissions - 角色权限关联表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS role_permissions (
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
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

COMMENT ON TABLE role_permissions IS '角色权限关联表';

-- -----------------------------------------------------------------------------
-- 表: user_roles - 用户角色关联表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_roles (
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
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles(expires_at) WHERE expires_at IS NOT NULL;

COMMENT ON TABLE user_roles IS '用户角色关联表';
COMMENT ON COLUMN user_roles.expires_at IS '角色过期时间，NULL表示永久有效';

-- 打印创建完成信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 用户与权限表创建完成';
    RAISE NOTICE '   - users: 用户表';
    RAISE NOTICE '   - roles: 角色表（已初始化5个角色）';
    RAISE NOTICE '   - permissions: 权限表（已初始化基础权限）';
    RAISE NOTICE '   - role_permissions: 角色权限关联表';
    RAISE NOTICE '   - user_roles: 用户角色关联表';
END $$;
