-- =============================================================================
-- MetaWorkflow V2.0 - 示例测试数据
-- =============================================================================
-- 创建时间: 2025-12-10
-- 说明: 创建示例数据用于开发和测试
-- 使用方法: psql -d metaworkflow -f sample_data.sql
-- =============================================================================

-- 注意：此文件仅用于开发和测试环境，请勿在生产环境使用

BEGIN;

-- -----------------------------------------------------------------------------
-- 1. 创建测试用户
-- -----------------------------------------------------------------------------
INSERT INTO users (id, username, email, password_hash, full_name, is_active, is_verified, is_developer) VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin', 'admin@example.com', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5QXZX1J2Z3W8e', 'Admin User', true, true, true),
    ('00000000-0000-0000-0000-000000000002', 'alice', 'alice@example.com', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5QXZX1J2Z3W8e', 'Alice Smith', true, true, true),
    ('00000000-0000-0000-0000-000000000003', 'bob', 'bob@example.com', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5QXZX1J2Z3W8e', 'Bob Johnson', true, true, false),
    ('00000000-0000-0000-0000-000000000004', 'charlie', 'charlie@example.com', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5QXZX1J2Z3W8e', 'Charlie Brown', true, true, false)
ON CONFLICT (id) DO NOTHING;

-- 为用户分配角色
INSERT INTO user_roles (user_id, role_id) 
SELECT '00000000-0000-0000-0000-000000000001', id FROM roles WHERE name = 'superadmin'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id) 
SELECT '00000000-0000-0000-0000-000000000002', id FROM roles WHERE name = 'developer'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id) 
SELECT '00000000-0000-0000-0000-000000000003', id FROM roles WHERE name = 'user'
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------------------------------
-- 2. 创建示例适配器
-- -----------------------------------------------------------------------------
INSERT INTO adapters (
    id, 
    name, 
    display_name, 
    description, 
    author_id, 
    domain, 
    status, 
    is_official, 
    latest_version,
    manifest
) VALUES
    (
        '10000000-0000-0000-0000-000000000001',
        'k12-chemistry-adapter',
        'K12化学适配器',
        '专为K12化学教育设计的领域适配器，支持概念理解、实验探索等多种教学活动',
        '00000000-0000-0000-0000-000000000002',
        'k12_education',
        'published',
        true,
        '1.0.0',
        '{
            "name": "k12-chemistry-adapter",
            "version": "1.0.0",
            "domain": "k12_education",
            "author": "Alice Smith",
            "methods": ["parse_request", "generate_objectives", "discover_content", "customize_workflow", "format_output", "validate_output"]
        }'::jsonb
    ),
    (
        '10000000-0000-0000-0000-000000000002',
        'art-history-adapter',
        '艺术史适配器',
        '艺术史教学适配器，支持艺术作品分析、时代背景讲解等',
        '00000000-0000-0000-0000-000000000002',
        'art_history',
        'published',
        true,
        '0.9.0',
        '{
            "name": "art-history-adapter",
            "version": "0.9.0",
            "domain": "art_history",
            "author": "Alice Smith"
        }'::jsonb
    )
ON CONFLICT (id) DO NOTHING;

-- 创建适配器版本
INSERT INTO adapter_versions (
    adapter_id,
    version,
    package_url,
    status,
    is_stable,
    release_notes
) VALUES
    (
        '10000000-0000-0000-0000-000000000001',
        '1.0.0',
        'https://example.com/packages/k12-chemistry-adapter-1.0.0.zip',
        'published',
        true,
        '首个正式版本，支持完整的K12化学教学功能'
    ),
    (
        '10000000-0000-0000-0000-000000000002',
        '0.9.0',
        'https://example.com/packages/art-history-adapter-0.9.0.zip',
        'published',
        true,
        'Beta版本，核心功能已完成'
    )
ON CONFLICT (adapter_id, version) DO NOTHING;

-- 创建适配器评论
INSERT INTO adapter_reviews (
    adapter_id,
    user_id,
    rating,
    title,
    content,
    version,
    is_verified_purchase
) VALUES
    (
        '10000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000003',
        5,
        '非常好用的化学教学工具',
        '这个适配器大大提升了我的化学课程设计效率，自动生成的实验方案非常专业！',
        '1.0.0',
        true
    ),
    (
        '10000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000004',
        4,
        '功能强大，但需要一点学习成本',
        '功能很全面，但对新手来说可能需要花点时间熟悉。文档写得不错。',
        '1.0.0',
        true
    )
ON CONFLICT (adapter_id, user_id) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 3. 创建示例工作流
-- -----------------------------------------------------------------------------
INSERT INTO workflows (
    id,
    user_id,
    domain,
    status,
    raw_input,
    subject,
    grade_level,
    duration,
    activity_type,
    started_at,
    completed_at,
    execution_time,
    context_data,
    outputs
) VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000003',
        'k12_education',
        'completed',
        '我想为七年级学生设计一节关于氧化还原反应的化学课',
        '化学',
        7,
        45,
        'concept_understanding',
        NOW() - INTERVAL '2 hours',
        NOW() - INTERVAL '1 hour 58 minutes',
        120.5,
        '{
            "grade": "七年级",
            "subject": "化学",
            "topic": "氧化还原反应",
            "student_level": "中等"
        }'::jsonb,
        '{
            "objectives_count": 3,
            "activities_count": 2,
            "success": true
        }'::jsonb
    ),
    (
        '20000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000003',
        'k12_education',
        'running',
        '帮我创建一个关于光合作用的生物课程',
        '生物',
        8,
        45,
        'experiment_exploration',
        NOW() - INTERVAL '5 minutes',
        NULL,
        NULL,
        '{
            "grade": "八年级",
            "subject": "生物",
            "topic": "光合作用"
        }'::jsonb,
        '{}'::jsonb
    )
ON CONFLICT (id) DO NOTHING;

-- 创建节点执行记录
INSERT INTO node_executions (
    workflow_id,
    node_id,
    node_type,
    status,
    started_at,
    completed_at,
    execution_time,
    input_data,
    output_data
) VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        'parse_request',
        'adapter_invoke',
        'completed',
        NOW() - INTERVAL '2 hours',
        NOW() - INTERVAL '1 hour 59 minutes 50 seconds',
        10.2,
        '{"user_input": "我想为七年级学生设计一节关于氧化还原反应的化学课"}'::jsonb,
        '{"grade": "七年级", "subject": "化学", "topic": "氧化还原反应"}'::jsonb
    ),
    (
        '20000000-0000-0000-0000-000000000001',
        'generate_objectives',
        'objective_generator',
        'completed',
        NOW() - INTERVAL '1 hour 59 minutes 50 seconds',
        NOW() - INTERVAL '1 hour 59 minutes 20 seconds',
        30.5,
        '{"grade": "七年级", "subject": "化学", "topic": "氧化还原反应"}'::jsonb,
        '{"objectives": [{"title": "理解氧化还原的概念", "bloom_level": "understand"}]}'::jsonb
    )
ON CONFLICT DO NOTHING;

-- 创建学习目标
INSERT INTO learning_objectives (
    workflow_id,
    title,
    description,
    bloom_level,
    estimated_time,
    sequence_order,
    keywords
) VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        '理解氧化还原的基本概念',
        '学生能够准确描述氧化和还原的定义，并识别简单的氧化还原反应',
        'understand',
        15,
        1,
        ARRAY['氧化', '还原', '电子转移', '化学反应']
    ),
    (
        '20000000-0000-0000-0000-000000000001',
        '应用氧化还原概念分析实际案例',
        '学生能够运用氧化还原知识分析日常生活中的化学现象',
        'apply',
        20,
        2,
        ARRAY['应用', '案例分析', '实际问题']
    ),
    (
        '20000000-0000-0000-0000-000000000001',
        '设计简单的氧化还原实验',
        '学生能够设计并执行简单的氧化还原实验，记录实验现象',
        'create',
        10,
        3,
        ARRAY['实验设计', '创造', '科学探究']
    )
ON CONFLICT DO NOTHING;

-- 创建教学活动
INSERT INTO teaching_activities (
    workflow_id,
    activity_type,
    title,
    content,
    duration,
    difficulty,
    tags
) VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        'concept_understanding',
        '氧化还原概念讲解',
        '{
            "introduction": "通过生活中的例子引入氧化还原概念",
            "main_points": ["氧化定义", "还原定义", "电子转移"],
            "examples": ["铁生锈", "燃烧反应"],
            "summary": "总结氧化还原的本质特征"
        }'::jsonb,
        15,
        'medium',
        ARRAY['概念讲解', '理论知识']
    ),
    (
        '20000000-0000-0000-0000-000000000001',
        'experiment_exploration',
        '铜与硝酸银反应实验',
        '{
            "objective": "观察氧化还原反应的实验现象",
            "materials": ["铜丝", "硝酸银溶液", "烧杯"],
            "procedure": ["将铜丝放入硝酸银溶液", "观察现象", "记录结果"],
            "safety": ["戴护目镜", "通风环境"],
            "expected_results": "铜丝表面析出银，溶液变蓝"
        }'::jsonb,
        20,
        'medium',
        ARRAY['实验', '动手操作', '现象观察']
    )
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------------------------------
-- 4. 创建审计日志示例
-- -----------------------------------------------------------------------------
INSERT INTO audit_logs (
    user_id,
    resource_type,
    resource_id,
    action,
    status,
    changes,
    ip_address
) VALUES
    (
        '00000000-0000-0000-0000-000000000001',
        'adapter',
        '10000000-0000-0000-0000-000000000001',
        'publish',
        'success',
        '{"status": {"from": "draft", "to": "published"}}'::jsonb,
        '192.168.1.100'::inet
    ),
    (
        '00000000-0000-0000-0000-000000000003',
        'workflow',
        '20000000-0000-0000-0000-000000000001',
        'create',
        'success',
        '{"workflow_id": "20000000-0000-0000-0000-000000000001"}'::jsonb,
        '192.168.1.101'::inet
    )
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------------------------------
-- 5. 创建系统事件示例
-- -----------------------------------------------------------------------------
INSERT INTO system_events (
    event_type,
    event_category,
    user_id,
    resource_type,
    resource_id,
    event_data
) VALUES
    (
        'adapter_published',
        'adapter',
        '00000000-0000-0000-0000-000000000002',
        'adapter',
        '10000000-0000-0000-0000-000000000001',
        '{"adapter_name": "k12-chemistry-adapter", "version": "1.0.0"}'::jsonb
    ),
    (
        'workflow_completed',
        'workflow',
        '00000000-0000-0000-0000-000000000003',
        'workflow',
        '20000000-0000-0000-0000-000000000001',
        '{"execution_time": 120.5, "success": true}'::jsonb
    ),
    (
        'user_registered',
        'user',
        '00000000-0000-0000-0000-000000000004',
        'user',
        '00000000-0000-0000-0000-000000000004',
        '{"username": "charlie", "email": "charlie@example.com"}'::jsonb
    )
ON CONFLICT DO NOTHING;

COMMIT;

-- 打印成功信息
DO $$ 
BEGIN 
    RAISE NOTICE '✅ 示例数据创建完成';
    RAISE NOTICE '   - 4 个测试用户（admin, alice, bob, charlie）';
    RAISE NOTICE '   - 2 个示例适配器（K12化学、艺术史）';
    RAISE NOTICE '   - 2 个工作流（1个已完成，1个运行中）';
    RAISE NOTICE '   - 3 个学习目标';
    RAISE NOTICE '   - 2 个教学活动';
    RAISE NOTICE '   - 审计日志和系统事件';
    RAISE NOTICE '';
    RAISE NOTICE '测试账号:';
    RAISE NOTICE '   - admin@example.com (超级管理员)';
    RAISE NOTICE '   - alice@example.com (开发者)';
    RAISE NOTICE '   - bob@example.com (普通用户)';
    RAISE NOTICE '   默认密码: password123';
END $$;
