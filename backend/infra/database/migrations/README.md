# 数据库迁移 (Database Migrations)

本目录包含 MetaWorkflow V2.0 的所有数据库迁移脚本。

## 📋 目录结构

```
migrations/
├── README.md                       # 本文件
├── migrate.sh                      # 主迁移脚本
├── 00_init_functions.sql          # 初始化函数和触发器
├── 01_users_and_permissions.sql   # 用户与权限表
├── 02_core_business.sql           # 核心业务表
├── 03_adapter_marketplace.sql     # 适配器市场表
├── 04_audit_and_logs.sql          # 审计与日志表
└── seeds/                          # 测试数据（可选）
    └── sample_data.sql
```

## 🗄️ 数据库架构

### 表清单

| 序号 | 表名 | 说明 | 行数估计 |
|------|------|------|---------|
| **用户与权限** ||||
| 1 | users | 用户表 | 10K |
| 2 | roles | 角色表 | 10 |
| 3 | permissions | 权限表 | 100 |
| 4 | role_permissions | 角色权限关联 | 50 |
| 5 | user_roles | 用户角色关联 | 10K |
| **核心业务** ||||
| 6 | workflows | 工作流记录 | 100K |
| 7 | node_executions | 节点执行记录 | 1M |
| 8 | learning_objectives | 学习目标 | 500K |
| 9 | teaching_activities | 教学活动 | 100K |
| **适配器市场** ||||
| 10 | adapter_categories | 适配器分类 | 100 |
| 11 | adapters | 适配器 | 1K |
| 12 | adapter_versions | 适配器版本 | 5K |
| 13 | adapter_dependencies | 依赖关系 | 2K |
| 14 | adapter_reviews | 评论 | 10K |
| 15 | adapter_stats | 统计数据 | 1K |
| 16 | user_subscriptions | 用户订阅 | 5K |
| **审计与日志** ||||
| 17 | audit_logs | 审计日志 | 10M |
| 18 | api_logs | API日志 | 100M |
| 19 | error_logs | 错误日志 | 1M |
| 20 | system_events | 系统事件 | 5M |
| 21 | performance_metrics | 性能指标 | 10M |

### 数据库要求

- **PostgreSQL**: 13+ （推荐 14+）
- **扩展**: 
  - `uuid-ossp` 或 PostgreSQL 13+ 内置 UUID 支持
  - `pg_trgm` (可选，用于模糊搜索)
- **存储**: 建议至少 20GB 初始空间
- **字符集**: UTF-8

## 🚀 快速开始

### 1. 环境准备

确保已安装 PostgreSQL 13+ 并启动服务：

```bash
# 检查 PostgreSQL 版本
psql --version

# 启动 PostgreSQL 服务
sudo service postgresql start

# 或使用 Docker
docker run --name metaworkflow-db \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:14
```

### 2. 配置数据库连接

设置环境变量：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=metaworkflow
export DB_USER=postgres
export DB_PASSWORD=postgres
```

或创建 `.env` 文件：

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=metaworkflow
DB_USER=postgres
DB_PASSWORD=postgres
```

### 3. 执行迁移

```bash
# 运行所有迁移（向上）
./migrate.sh up

# 回滚所有迁移（向下）
./migrate.sh down

# 重置数据库（删除并重建）
./migrate.sh reset
```

### 4. 验证迁移

```bash
# 连接到数据库
psql -h localhost -U postgres -d metaworkflow

# 查看所有表
\dt

# 查看表结构
\d users
\d workflows
\d adapters

# 查看索引
\di

# 退出
\q
```

## 📝 迁移脚本说明

### 00_init_functions.sql

创建通用的数据库函数和触发器：

- `update_updated_at_column()` - 自动更新 `updated_at` 字段
- `update_adapter_rating_stats()` - 更新适配器评分统计
- `validate_semantic_version()` - 验证语义化版本号
- `jsonb_merge()` - 合并 JSONB 对象

### 01_users_and_permissions.sql

创建用户和权限相关表：

- `users` - 用户基本信息、状态、偏好设置
- `roles` - 角色定义（user, developer, moderator, admin, superadmin）
- `permissions` - 权限定义（workflow, adapter, user, system 等）
- `role_permissions` - 角色权限映射
- `user_roles` - 用户角色映射

**初始化数据**:
- 5 个预定义角色
- 基础权限集（19 个权限）

### 02_core_business.sql

创建核心业务表：

- `workflows` - 工作流执行记录，支持多领域（K12、艺术史、畅销书等）
- `node_executions` - 节点级别的执行记录
- `learning_objectives` - 学习目标（支持 Bloom 分类法）
- `teaching_activities` - 教学活动（多种活动类型）

**特性**:
- JSONB 字段存储灵活数据
- GIN 索引支持 JSON 查询
- 时间分区支持（可选）

### 03_adapter_marketplace.sql

创建适配器市场表：

- `adapter_categories` - 分类体系（支持层级结构）
- `adapters` - 适配器基本信息、状态、统计
- `adapter_versions` - 版本管理（语义化版本）
- `adapter_dependencies` - 依赖关系图
- `adapter_reviews` - 用户评论和评分
- `adapter_stats` - 每日统计数据
- `user_subscriptions` - 用户订阅通知

**初始化数据**:
- 7 个基础分类（教育、内容创作等）

**特性**:
- 全文搜索索引
- 自动更新评分统计触发器
- 依赖循环检测约束

### 04_audit_and_logs.sql

创建审计和日志表：

- `audit_logs` - 审计日志（所有重要操作）
- `api_logs` - API 访问日志
- `error_logs` - 错误日志（分级）
- `system_events` - 系统事件流
- `performance_metrics` - 性能指标

**特性**:
- 高写入性能优化
- 时间分区支持（建议用于日志表）
- 慢查询和错误专用索引

## 🔧 高级用法

### 手动执行单个迁移

```bash
psql -h localhost -U postgres -d metaworkflow \
  -f 01_users_and_permissions.sql
```

### 只创建特定表

编辑 SQL 文件，注释掉不需要的表定义。

### 加载测试数据

```bash
# 如果有测试数据文件
psql -h localhost -U postgres -d metaworkflow \
  -f seeds/sample_data.sql
```

### 备份数据库

```bash
# 备份整个数据库
pg_dump -h localhost -U postgres -d metaworkflow \
  -F c -f metaworkflow_backup.dump

# 备份 Schema 定义
pg_dump -h localhost -U postgres -d metaworkflow \
  --schema-only -f metaworkflow_schema.sql

# 备份数据
pg_dump -h localhost -U postgres -d metaworkflow \
  --data-only -f metaworkflow_data.sql
```

### 恢复数据库

```bash
# 从备份恢复
pg_restore -h localhost -U postgres -d metaworkflow \
  -c metaworkflow_backup.dump
```

## 📊 性能优化

### 索引策略

1. **主键索引**: 所有表自动创建 UUID 主键索引
2. **外键索引**: 所有外键字段创建索引
3. **GIN 索引**: JSON/数组字段使用 GIN 索引
4. **复合索引**: 常见查询模式创建复合索引
5. **部分索引**: 条件过滤创建部分索引

### 分区建议

对于大表（如日志表），建议启用分区：

```sql
-- 示例：workflows 表按月分区
CREATE TABLE workflows_2025_12 PARTITION OF workflows
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE workflows_2026_01 PARTITION OF workflows
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
```

### 维护任务

```bash
# 定期执行（建议每周）

# 1. 重建索引
psql -d metaworkflow -c "REINDEX DATABASE metaworkflow;"

# 2. 更新统计信息
psql -d metaworkflow -c "ANALYZE;"

# 3. 清理过期数据（日志表）
psql -d metaworkflow -c "
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '90 days';
"

# 4. 检查表膨胀
psql -d metaworkflow -c "
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
"
```

## 🔍 常见问题

### Q1: 迁移失败怎么办？

**A**: 检查错误信息，常见原因：
- 数据库连接配置错误
- PostgreSQL 版本过低
- 权限不足
- 表已存在（使用 `IF NOT EXISTS` 可避免）

### Q2: 如何修改现有表结构？

**A**: 创建新的迁移文件：

```sql
-- 05_alter_users_add_phone.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_phone ON users(phone);
```

### Q3: 如何处理数据迁移冲突？

**A**: 使用事务和条件语句：

```sql
BEGIN;

-- 检查是否已存在
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        INSERT INTO users (username, email, password_hash) 
        VALUES ('admin', 'admin@example.com', 'hashed_password');
    END IF;
END $$;

COMMIT;
```

### Q4: 生产环境如何安全迁移？

**A**: 推荐步骤：
1. 在测试环境完整测试
2. 备份生产数据库
3. 选择低流量时段
4. 使用事务包裹迁移
5. 准备回滚方案
6. 监控性能指标

## 📚 相关文档

- [DATABASE_DESIGN.md](../../../docs_new/DATABASE_DESIGN.md) - 完整数据库设计文档
- [架构文档](../../../docs_new/architecture.md) - 系统架构设计
- [API 设计](../../../docs_new/api_design.md) - API 接口设计

## 🤝 贡献指南

添加新迁移时：

1. 创建新文件，命名格式：`XX_description.sql`
2. 添加完整的注释和文档
3. 使用 `IF NOT EXISTS` 避免重复创建
4. 添加必要的索引和约束
5. 更新 `migrate.sh` 中的 `MIGRATION_FILES` 数组
6. 更新本 README 文档

## 📞 支持

如有问题，请：

1. 查看错误日志
2. 阅读相关文档
3. 搜索已知问题
4. 提交 Issue

---

**版本**: 2.0  
**最后更新**: 2025-12-10  
**维护者**: Coze Studio Team
