# 数据库Schema创建完成总结

## 概述

已成功完成 MetaWorkflow V2.0 的数据库Schema设计和创建，包括21张核心表、5个迁移脚本、完整的索引策略和示例数据。

**完成时间**: 2025-12-10  
**数据库**: PostgreSQL 13+  
**文档参考**: docs_new/DATABASE_DESIGN.md

## 创建的文件清单

### SQL迁移脚本 (5个)

| 文件名 | 说明 | 行数 | 创建内容 |
|--------|------|------|---------|
| **00_init_functions.sql** | 初始化函数 | ~100 | 4个数据库函数 |
| **01_users_and_permissions.sql** | 用户权限 | ~280 | 5张表 + 5角色 + 19权限 |
| **02_core_business.sql** | 核心业务 | ~250 | 4张核心业务表 |
| **03_adapter_marketplace.sql** | 适配器市场 | ~450 | 7张市场相关表 + 7分类 |
| **04_audit_and_logs.sql** | 审计日志 | ~280 | 5张日志监控表 |

### 辅助文件 (4个)

| 文件名 | 说明 | 行数 |
|--------|------|------|
| **migrate.sh** | 迁移主脚本 | ~250 | Bash脚本，支持 up/down/reset |
| **README.md** | 使用文档 | ~450 | 完整使用指南 |
| **seeds/sample_data.sql** | 示例数据 | ~350 | 测试用示例数据 |
| **DATABASE_SCHEMA_COMPLETION_SUMMARY.md** | 本文件 | ~400 | 完成总结 |

## 数据库表结构

### 📊 表统计

**总计**: 21张表  
**预估总行数**: ~138M 行（含日志数据）  
**预估存储**: 20GB+ （初始）

### 表分类明细

#### 1. 用户与权限 (5张表)

| 表名 | 说明 | 关键字段 | 索引数 |
|------|------|---------|--------|
| `users` | 用户表 | username, email, is_developer | 4 |
| `roles` | 角色表 | name, level | 2 |
| `permissions` | 权限表 | name, resource, action | 3 |
| `role_permissions` | 角色权限关联 | role_id, permission_id | 2 |
| `user_roles` | 用户角色关联 | user_id, role_id | 3 |

**初始化数据**:
- ✅ 5个预定义角色（user, developer, moderator, admin, superadmin）
- ✅ 19个基础权限（workflow, adapter, user, system相关）

#### 2. 核心业务 (4张表)

| 表名 | 说明 | 关键字段 | 索引数 |
|------|------|---------|--------|
| `workflows` | 工作流记录 | domain, status, outputs | 11 |
| `node_executions` | 节点执行 | workflow_id, node_type, status | 7 |
| `learning_objectives` | 学习目标 | workflow_id, bloom_level | 5 |
| `teaching_activities` | 教学活动 | workflow_id, activity_type | 6 |

**支持的领域**:
- k12_education（K12教育）
- art_history（艺术史）
- bestseller（畅销书）
- vocational（职业培训）
- higher_education（高等教育）
- professional_training（专业培训）

**支持的Bloom层级**:
- remember（记忆）
- understand（理解）
- apply（应用）
- analyze（分析）
- evaluate（评价）
- create（创造）

#### 3. 适配器市场 (7张表)

| 表名 | 说明 | 关键字段 | 索引数 |
|------|------|---------|--------|
| `adapter_categories` | 分类 | name, parent_id, path | 4 |
| `adapters` | 适配器 | name, status, manifest | 10 |
| `adapter_versions` | 版本 | adapter_id, version | 5 |
| `adapter_dependencies` | 依赖 | adapter_id, depends_on_id | 2 |
| `adapter_reviews` | 评论 | adapter_id, user_id, rating | 6 |
| `adapter_stats` | 统计 | adapter_id, stat_date | 3 |
| `user_subscriptions` | 订阅 | user_id, adapter_id | 3 |

**初始化数据**:
- ✅ 7个分类（教育、K12、高等教育、职业、内容、艺术、文学）

**特殊功能**:
- ✅ 全文搜索索引（GIN + tsvector）
- ✅ 自动评分统计触发器
- ✅ 层级分类支持

#### 4. 审计与日志 (5张表)

| 表名 | 说明 | 关键字段 | 索引数 |
|------|------|---------|--------|
| `audit_logs` | 审计日志 | user_id, resource_type, action | 7 |
| `api_logs` | API日志 | method, path, status_code | 10 |
| `error_logs` | 错误日志 | error_type, severity | 6 |
| `system_events` | 系统事件 | event_type, event_category | 7 |
| `performance_metrics` | 性能指标 | metric_name, metric_type | 4 |

**日志级别**:
- debug, info, warning, error, critical

## 数据库函数和触发器

### 创建的函数 (4个)

1. **update_updated_at_column()**
   - 功能：自动更新 `updated_at` 字段
   - 使用：所有带 `updated_at` 字段的表

2. **update_adapter_rating_stats()**
   - 功能：自动更新适配器评分统计
   - 触发：`adapter_reviews` 表的 INSERT/UPDATE/DELETE

3. **validate_semantic_version()**
   - 功能：验证语义化版本号格式
   - 用途：适配器版本验证

4. **jsonb_merge()**
   - 功能：合并JSONB对象
   - 用途：JSON数据合并操作

### 创建的触发器 (9个)

| 表名 | 触发器 | 功能 |
|------|--------|------|
| users | update_users_updated_at | 自动更新时间戳 |
| roles | update_roles_updated_at | 自动更新时间戳 |
| workflows | update_workflows_updated_at | 自动更新时间戳 |
| teaching_activities | update_teaching_activities_updated_at | 自动更新时间戳 |
| adapter_categories | update_adapter_categories_updated_at | 自动更新时间戳 |
| adapters | update_adapters_updated_at | 自动更新时间戳 |
| adapter_reviews | update_adapter_reviews_updated_at | 自动更新时间戳 |
| adapter_reviews | update_adapter_stats_on_review | 更新评分统计 |
| user_subscriptions | update_user_subscriptions_updated_at | 自动更新时间戳 |

## 索引策略

### 索引统计

- **B-tree 索引**: ~80个
- **GIN 索引**: ~15个（JSON/数组/全文搜索）
- **部分索引**: ~10个（条件过滤）
- **复合索引**: ~20个（多列查询优化）

### 关键索引

#### 全文搜索索引

```sql
-- 适配器全文搜索
CREATE INDEX idx_adapters_fts ON adapters USING GIN(
    to_tsvector('english', display_name || ' ' || COALESCE(description, ''))
);
```

#### JSON查询索引

```sql
-- 工作流上下文查询
CREATE INDEX idx_workflows_context_data ON workflows USING GIN(context_data);
CREATE INDEX idx_workflows_outputs ON workflows USING GIN(outputs);
```

#### 性能优化索引

```sql
-- 慢查询索引
CREATE INDEX idx_api_logs_slow_queries ON api_logs(response_time DESC, created_at DESC) 
    WHERE response_time > 1000;

-- 错误日志索引
CREATE INDEX idx_api_logs_errors ON api_logs(status_code, created_at DESC) 
    WHERE status_code >= 400;
```

## 使用方法

### 1. 快速开始

```bash
cd /home/liunix/work/coze-studio/backend/infra/database/migrations

# 配置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=metaworkflow
export DB_USER=postgres
export DB_PASSWORD=postgres

# 执行迁移
./migrate.sh up
```

### 2. 加载示例数据

```bash
psql -h localhost -U postgres -d metaworkflow \
  -f seeds/sample_data.sql
```

### 3. 验证安装

```bash
# 连接数据库
psql -h localhost -U postgres -d metaworkflow

# 查看所有表
\dt

# 查看表结构
\d users
\d workflows
\d adapters

# 查看示例数据
SELECT username, email, is_developer FROM users;
SELECT name, display_name, status FROM adapters;
SELECT domain, status, subject FROM workflows;

# 退出
\q
```

## 特性亮点

### 1. 数据完整性

- ✅ 所有外键约束
- ✅ CHECK约束验证数据有效性
- ✅ UNIQUE约束防止重复
- ✅ NOT NULL约束确保必填字段

### 2. 性能优化

- ✅ 合理的索引策略
- ✅ GIN索引支持JSON查询
- ✅ 部分索引优化存储
- ✅ 复合索引优化常见查询

### 3. 可扩展性

- ✅ JSONB字段存储灵活数据
- ✅ 数组类型支持标签等多值字段
- ✅ 时间分区支持（可选）
- ✅ 层级分类支持

### 4. 审计追溯

- ✅ 完整的审计日志
- ✅ API访问日志
- ✅ 错误日志分级
- ✅ 系统事件记录
- ✅ 性能指标监控

### 5. 开发友好

- ✅ 详细的注释说明
- ✅ 完整的示例数据
- ✅ 自动化迁移脚本
- ✅ 友好的README文档

## 示例数据

### 测试账号

| 用户名 | 邮箱 | 角色 | 密码 |
|--------|------|------|------|
| admin | admin@example.com | 超级管理员 | password123 |
| alice | alice@example.com | 开发者 | password123 |
| bob | bob@example.com | 普通用户 | password123 |
| charlie | charlie@example.com | 普通用户 | password123 |

### 示例适配器

1. **K12化学适配器** (k12-chemistry-adapter)
   - 版本：1.0.0
   - 状态：已发布
   - 评分：4.5/5.0
   - 评论：2条

2. **艺术史适配器** (art-history-adapter)
   - 版本：0.9.0
   - 状态：已发布（Beta）

### 示例工作流

1. **氧化还原反应课程**
   - 状态：已完成
   - 学科：化学
   - 年级：七年级
   - 时长：45分钟
   - 学习目标：3个
   - 教学活动：2个

2. **光合作用课程**
   - 状态：运行中
   - 学科：生物
   - 年级：八年级

## 数据库维护

### 定期任务

#### 1. 重建索引（每周）

```bash
psql -d metaworkflow -c "REINDEX DATABASE metaworkflow;"
```

#### 2. 更新统计信息（每天）

```bash
psql -d metaworkflow -c "ANALYZE;"
```

#### 3. 清理过期日志（每月）

```sql
-- 删除90天前的API日志
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '90 days';

-- 删除180天前的审计日志
DELETE FROM audit_logs 
WHERE created_at < NOW() - INTERVAL '180 days';
```

#### 4. 备份数据库（每天）

```bash
pg_dump -h localhost -U postgres -d metaworkflow \
  -F c -f metaworkflow_backup_$(date +%Y%m%d).dump
```

### 性能监控

#### 查看表大小

```sql
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### 查看慢查询

```sql
SELECT query, calls, mean_exec_time, max_exec_time
FROM pg_stat_statements
WHERE mean_exec_time > 100
ORDER BY mean_exec_time DESC
LIMIT 20;
```

## 与应用层集成

### Go代码示例

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

// 连接数据库
func ConnectDB() (*sql.DB, error) {
    connStr := "host=localhost port=5432 user=postgres password=postgres dbname=metaworkflow sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // 测试连接
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}

// 查询用户
func GetUser(db *sql.DB, username string) (*User, error) {
    var user User
    err := db.QueryRow(`
        SELECT id, username, email, full_name, is_active 
        FROM users 
        WHERE username = $1
    `, username).Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.IsActive)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

## 后续工作建议

### 短期 (1-2周)

1. **ORM集成**
   - [ ] 配置 GORM 或 sqlx
   - [ ] 生成 Go model 结构
   - [ ] 创建 Repository 层

2. **数据访问层**
   - [ ] 实现 CRUD 操作
   - [ ] 添加事务支持
   - [ ] 实现分页查询

3. **测试覆盖**
   - [ ] 数据库单元测试
   - [ ] 集成测试
   - [ ] 性能测试

### 中期 (1-2个月)

1. **性能优化**
   - [ ] 启用查询缓存（Redis）
   - [ ] 实现连接池优化
   - [ ] 添加读写分离

2. **监控告警**
   - [ ] 集成 Prometheus
   - [ ] 配置 Grafana 仪表板
   - [ ] 设置慢查询告警

3. **数据迁移工具**
   - [ ] 实现版本管理
   - [ ] 支持回滚机制
   - [ ] 自动化部署

### 长期 (3-6个月)

1. **分布式支持**
   - [ ] 数据分片策略
   - [ ] 跨区域复制
   - [ ] 高可用架构

2. **数据分析**
   - [ ] 实时数据仓库
   - [ ] BI 报表系统
   - [ ] 用户行为分析

## 常见问题

### Q1: 如何修改表结构？

**A**: 创建新的迁移文件：

```sql
-- 05_alter_table.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_phone ON users(phone);
```

然后更新 `migrate.sh` 中的 `MIGRATION_FILES` 数组。

### Q2: 如何处理大量数据导入？

**A**: 使用 `COPY` 命令或批量插入：

```sql
-- 禁用索引
DROP INDEX idx_workflows_user_id;

-- 批量导入
COPY workflows FROM '/path/to/data.csv' CSV HEADER;

-- 重建索引
CREATE INDEX idx_workflows_user_id ON workflows(user_id);
ANALYZE workflows;
```

### Q3: 如何实现软删除？

**A**: 添加 `deleted_at` 字段：

```sql
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
```

## 总结

✅ **完整性**: 21张表，覆盖所有核心业务需求  
✅ **规范性**: 遵循数据库设计最佳实践  
✅ **性能**: 完善的索引策略和查询优化  
✅ **可维护**: 详细的文档和自动化脚本  
✅ **可扩展**: JSONB、数组、分区支持  
✅ **可追溯**: 完整的审计和日志系统  

数据库Schema已完全就绪，可以支持 MetaWorkflow V2.0 的完整功能！

---

**相关文件**:
- 设计文档: docs_new/DATABASE_DESIGN.md
- 迁移脚本: backend/infra/database/migrations/
- 使用指南: backend/infra/database/migrations/README.md
- 整体文档: backend/domain/DOCS_INDEX.md
