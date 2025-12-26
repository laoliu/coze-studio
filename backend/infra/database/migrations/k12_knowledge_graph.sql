-- K12 Knowledge Graph Database Schema
-- 动态知识图谱数据库设计

-- 1. 知识点表
CREATE TABLE IF NOT EXISTS knowledge_points (
    id VARCHAR(64) PRIMARY KEY COMMENT '知识点ID',
    name VARCHAR(255) NOT NULL COMMENT '知识点名称',
    subject VARCHAR(32) NOT NULL COMMENT '学科',
    grade VARCHAR(32) NOT NULL COMMENT '年级',
    description TEXT COMMENT '描述',
    keywords JSON COMMENT '关键词数组',
    prerequisites JSON COMMENT '前置知识点ID数组',
    next_points JSON COMMENT '后续知识点ID数组',
    difficulty VARCHAR(32) COMMENT '难度：basic/improve/advanced',
    metadata JSON COMMENT '扩展元数据',
    version INT DEFAULT 1 COMMENT '版本号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    INDEX idx_subject_grade (subject, grade),
    INDEX idx_created_at (created_at),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识点表';

-- 2. 学习资源表
CREATE TABLE IF NOT EXISTS learning_resources (
    id VARCHAR(64) PRIMARY KEY COMMENT '资源ID',
    knowledge_id VARCHAR(64) NOT NULL COMMENT '关联知识点ID',
    type VARCHAR(32) NOT NULL COMMENT '资源类型：video/textbook/article',
    title VARCHAR(255) NOT NULL COMMENT '资源标题',
    description TEXT COMMENT '资源描述',
    url TEXT COMMENT '资源链接',
    difficulty VARCHAR(32) COMMENT '难度等级',
    duration INT COMMENT '时长（分钟）',
    source VARCHAR(255) COMMENT '来源',
    metadata JSON COMMENT '扩展元数据',
    version INT DEFAULT 1 COMMENT '版本号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    FOREIGN KEY (knowledge_id) REFERENCES knowledge_points(id) ON DELETE CASCADE,
    INDEX idx_knowledge_id (knowledge_id),
    INDEX idx_type (type),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学习资源表';

-- 3. 习题表
CREATE TABLE IF NOT EXISTS exercises (
    id VARCHAR(64) PRIMARY KEY COMMENT '习题ID',
    knowledge_id VARCHAR(64) NOT NULL COMMENT '关联知识点ID',
    question TEXT NOT NULL COMMENT '题目',
    options JSON COMMENT '选项数组',
    answer TEXT NOT NULL COMMENT '答案',
    solution TEXT COMMENT '解析',
    difficulty VARCHAR(32) COMMENT '难度等级',
    score INT COMMENT '分值',
    metadata JSON COMMENT '扩展元数据',
    version INT DEFAULT 1 COMMENT '版本号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    FOREIGN KEY (knowledge_id) REFERENCES knowledge_points(id) ON DELETE CASCADE,
    INDEX idx_knowledge_id (knowledge_id),
    INDEX idx_difficulty (difficulty),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='习题表';

-- 4. 生成日志表
CREATE TABLE IF NOT EXISTS generation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    request_type VARCHAR(64) COMMENT '请求类型：knowledge_point/resource/exercise',
    subject VARCHAR(32) COMMENT '学科',
    grade VARCHAR(32) COMMENT '年级',
    knowledge_id VARCHAR(64) COMMENT '知识点ID',
    prompt_template TEXT COMMENT '提示词模板',
    llm_response TEXT COMMENT 'LLM响应',
    success BOOLEAN COMMENT '是否成功',
    error_message TEXT COMMENT '错误信息',
    duration_ms INT COMMENT '耗时（毫秒）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_request_type (request_type),
    INDEX idx_success (success),
    INDEX idx_created_at (created_at),
    INDEX idx_knowledge_id (knowledge_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='生成日志表';

-- 5. 缓存配置表
CREATE TABLE IF NOT EXISTS cache_config (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    item_type VARCHAR(32) NOT NULL COMMENT '内容类型：knowledge_point/resource/exercise',
    ttl_seconds INT NOT NULL COMMENT 'TTL（秒）',
    description VARCHAR(255) COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_item_type (item_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='缓存配置表';

-- 插入默认缓存配置
INSERT INTO cache_config (item_type, ttl_seconds, description) VALUES
('knowledge_point', 604800, '知识点缓存7天'),
('resource', 259200, '资源缓存3天'),
('exercise', 86400, '习题缓存1天')
ON DUPLICATE KEY UPDATE ttl_seconds=VALUES(ttl_seconds);
