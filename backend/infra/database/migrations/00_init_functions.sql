-- =============================================================================
-- MetaWorkflow V2.0 - 初始化函数和触发器
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建通用的数据库函数和触发器
-- =============================================================================

-- 删除已存在的函数（如果有）
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
DROP FUNCTION IF EXISTS update_adapter_rating_stats CASCADE;

-- -----------------------------------------------------------------------------
-- 函数: update_updated_at_column
-- 说明: 自动更新表的 updated_at 字段
-- 用法: 在表的 BEFORE UPDATE 触发器中调用
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS '自动更新updated_at字段的触发器函数';

-- -----------------------------------------------------------------------------
-- 函数: update_adapter_rating_stats
-- 说明: 更新适配器的评分统计信息
-- 触发条件: adapter_reviews 表的 INSERT/UPDATE/DELETE
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_adapter_rating_stats()
RETURNS TRIGGER AS $$
DECLARE
    v_adapter_id UUID;
    v_avg_rating DECIMAL(3, 2);
    v_review_count INTEGER;
BEGIN
    -- 确定要更新的 adapter_id
    IF TG_OP = 'DELETE' THEN
        v_adapter_id := OLD.adapter_id;
    ELSE
        v_adapter_id := NEW.adapter_id;
    END IF;
    
    -- 计算平均评分和评论数
    SELECT 
        COALESCE(AVG(rating), 0.00),
        COUNT(*)
    INTO 
        v_avg_rating,
        v_review_count
    FROM adapter_reviews
    WHERE adapter_id = v_adapter_id
        AND status = 'published';
    
    -- 更新适配器表的统计信息
    UPDATE adapters
    SET 
        average_rating = v_avg_rating,
        review_count = v_review_count,
        updated_at = NOW()
    WHERE id = v_adapter_id;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_adapter_rating_stats() IS '更新适配器评分统计的触发器函数';

-- -----------------------------------------------------------------------------
-- 函数: generate_random_uuid
-- 说明: 为旧版本PostgreSQL提供UUID生成支持
-- 注意: PostgreSQL 13+ 已内置 gen_random_uuid()
-- -----------------------------------------------------------------------------
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- 如果使用 PostgreSQL < 13，取消注释上面这行

-- -----------------------------------------------------------------------------
-- 函数: validate_semantic_version
-- 说明: 验证语义化版本号格式 (例如: 1.2.3, 2.0.0-beta.1)
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION validate_semantic_version(version TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN version ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMENT ON FUNCTION validate_semantic_version(TEXT) IS '验证语义化版本号格式';

-- -----------------------------------------------------------------------------
-- 函数: jsonb_merge
-- 说明: 深度合并两个JSONB对象
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION jsonb_merge(current_data JSONB, new_data JSONB)
RETURNS JSONB AS $$
BEGIN
    RETURN current_data || new_data;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMENT ON FUNCTION jsonb_merge(JSONB, JSONB) IS '合并两个JSONB对象';

-- 打印初始化完成信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 初始化函数创建完成';
END $$;
