# ORM仓储层设计文档

**版本**: 2.0  
**日期**: 2025-12-09  
**状态**: 完成 ✅

---

## 📋 目录

1. [概述](#概述)
2. [架构设计](#架构设计)
3. [Repository列表](#repository列表)
4. [使用指南](#使用指南)
5. [最佳实践](#最佳实践)
6. [性能优化](#性能优化)
7. [测试](#测试)

---

## 概述

### 设计模式

采用 **Repository Pattern（仓储模式）** + **Unit of Work Pattern（工作单元模式）**：

```
Controller/Service
    ↓
UnitOfWork
    ↓
Repository Layer
    ↓
ORM Models (SQLAlchemy)
    ↓
Database (PostgreSQL)
```

### 核心优势

✅ **解耦业务逻辑和数据访问**  
✅ **统一事务管理**  
✅ **简化测试（Mock Repository）**  
✅ **代码复用（BaseRepository）**  
✅ **类型安全（泛型支持）**  

---

## 架构设计

### 1. BaseRepository（基础仓储）

```python
class BaseRepository(Generic[T]):
    """通用CRUD基类"""
    
    async def create(self, obj: T) -> T
    async def get(self, id: UUID) -> Optional[T]
    async def get_multi(self, skip, limit, filters) -> List[T]
    async def update(self, id: UUID, values) -> Optional[T]
    async def delete(self, id: UUID) -> bool
    async def count(self, filters) -> int
    async def exists(self, id: UUID) -> bool
```

### 2. 具体Repository

继承BaseRepository并添加特定业务方法：

```python
class WorkflowRepository(BaseRepository[WorkflowRecord]):
    """工作流仓储"""
    
    # 继承基础CRUD
    # + 特定业务方法
    async def create_from_context(...)
    async def update_status(...)
    async def list_by_domain(...)
    async def get_statistics(...)
```

### 3. UnitOfWork（工作单元）

统一管理所有Repository和事务：

```python
class UnitOfWork:
    """工作单元"""
    
    workflows: WorkflowRepository
    nodes: NodeExecutionRepository
    learning_objectives: LearningObjectiveRepository
    adapters: AdapterRepository
    # ... 更多Repository
    
    async def commit()
    async def rollback()
    async def flush()
```

---

## Repository列表

### 核心业务Repository

| Repository | 模型 | 主要功能 |
|------------|------|---------|
| **WorkflowRepository** | WorkflowRecord | 工作流CRUD、状态管理、统计 |
| **NodeExecutionRepository** | NodeExecutionRecord | 节点执行、状态更新 |
| **LearningObjectiveRepository** | LearningObjectiveRecord | 学习目标批量操作 |
| **TeachingActivityRepository** | TeachingActivityRecord | 教学活动管理 |

### 适配器相关Repository

| Repository | 模型 | 主要功能 |
|------------|------|---------|
| **AdapterRepository** | AdapterRecord | 适配器CRUD、搜索、排行 |
| **AdapterVersionRepository** | AdapterVersionRecord | 版本管理、最新版本 |
| **AdapterReviewRepository** | AdapterReviewRecord | 评论管理、评分统计 |
| **AdapterStatsRepository** | AdapterStatsRecord | 每日统计、趋势分析 |

---

## 使用指南

### 1. 基本CRUD操作

```python
from src.db import get_uow

async def example():
    async with get_uow() as uow:
        # 创建
        workflow = await uow.workflows.create_from_context(context)
        
        # 查询
        found = await uow.workflows.get(workflow.id)
        
        # 更新
        await uow.workflows.update_status(
            workflow.id,
            WorkflowStatus.RUNNING
        )
        
        # 删除
        await uow.workflows.delete(workflow.id)
        
        # 自动提交
```

### 2. 批量操作

```python
async with get_uow() as uow:
    # 批量创建学习目标
    objectives = [
        {'title': '目标1', 'bloom_level': 'understand', ...},
        {'title': '目标2', 'bloom_level': 'apply', ...}
    ]
    
    created = await uow.learning_objectives.create_batch(
        workflow_id,
        objectives
    )
    
    print(f"创建了 {len(created)} 个学习目标")
```

### 3. 复杂查询

```python
async with get_uow() as uow:
    # 按领域查询
    workflows = await uow.workflows.list_by_domain(
        DomainType.K12_EDUCATION,
        limit=50,
        offset=0
    )
    
    # 带过滤条件
    workflows = await uow.workflows.get_multi(
        skip=0,
        limit=100,
        filters={
            'domain': 'k12_education',
            'status': 'completed'
        },
        order_by='-created_at'
    )
```

### 4. 关联查询

```python
async with get_uow() as uow:
    # 查询工作流及其关联数据（使用eager loading）
    workflow = await uow.workflows.get_with_relations(workflow_id)
    
    print(f"节点数: {len(workflow.nodes)}")
    print(f"学习目标数: {len(workflow.learning_objectives)}")
```

### 5. 统计聚合

```python
async with get_uow() as uow:
    # 工作流统计
    stats = await uow.workflows.get_statistics(
        domain=DomainType.K12_EDUCATION
    )
    
    print(f"总数: {stats['total']}")
    print(f"完成: {stats['completed']}")
    print(f"平均执行时间: {stats['avg_execution_time']}秒")
    
    # 适配器评分统计
    rating_stats = await uow.adapter_reviews.get_rating_stats(adapter_id)
    print(f"平均评分: {rating_stats['average_rating']}")
    print(f"评分分布: {rating_stats['rating_distribution']}")
```

### 6. 搜索功能

```python
async with get_uow() as uow:
    # 全文搜索适配器
    results = await uow.adapters.search(
        "化学实验",
        domain=DomainType.K12_EDUCATION,
        limit=20
    )
    
    for adapter in results:
        print(f"{adapter.display_name}: {adapter.description}")
```

### 7. 事务控制

#### 自动提交（推荐）

```python
async with get_uow() as uow:
    # 所有操作
    await uow.workflows.create(...)
    await uow.nodes.create_node(...)
    # 退出时自动提交
```

#### 手动控制

```python
async with UnitOfWork() as uow:
    try:
        # 操作1
        workflow = await uow.workflows.create(...)
        await uow.flush()  # 刷新但不提交
        
        # 操作2
        await uow.nodes.create_node(...)
        
        # 手动提交
        await uow.commit()
    except Exception as e:
        # 手动回滚
        await uow.rollback()
        raise
```

### 8. FastAPI集成

```python
from fastapi import APIRouter, Depends
from src.db import UnitOfWork, get_db_session

router = APIRouter()

@router.post("/workflows")
async def create_workflow(
    request: WorkflowRequest,
    session: AsyncSession = Depends(get_db_session)
):
    async with UnitOfWork(session) as uow:
        workflow = await uow.workflows.create_from_context(
            request.to_context()
        )
        return {"id": workflow.id}
```

---

## 最佳实践

### 1. 始终使用UnitOfWork

✅ **好的做法**：
```python
async with get_uow() as uow:
    workflow = await uow.workflows.create(...)
    await uow.learning_objectives.create_batch(...)
    # 一个事务
```

❌ **不好的做法**：
```python
async with get_db() as session:
    repo = WorkflowRepository(session)
    await repo.create(...)
    # 手动管理太复杂
```

### 2. 避免N+1查询

✅ **好的做法**：
```python
# 使用eager loading
workflow = await uow.workflows.get_with_relations(workflow_id)
for node in workflow.nodes:  # 已加载，无额外查询
    print(node.node_type)
```

❌ **不好的做法**：
```python
workflow = await uow.workflows.get(workflow_id)
for node_id in workflow.node_ids:
    node = await uow.nodes.get(node_id)  # N次查询！
```

### 3. 使用批量操作

✅ **好的做法**：
```python
# 批量创建
await uow.learning_objectives.create_batch(
    workflow_id,
    objectives_list
)
```

❌ **不好的做法**：
```python
# 循环单个创建
for obj in objectives_list:
    await uow.learning_objectives.create(obj)  # 慢！
```

### 4. 合理使用flush和commit

```python
async with UnitOfWork() as uow:
    # 操作1
    workflow = await uow.workflows.create(...)
    await uow.flush()  # 生成ID但不提交
    
    # 操作2（需要使用workflow.id）
    await uow.nodes.create_node(workflow.id, ...)
    
    # 最后统一提交
    await uow.commit()
```

### 5. 异常处理

```python
async with get_uow() as uow:
    try:
        # 业务逻辑
        workflow = await uow.workflows.create(...)
        
        # 可能抛出异常的操作
        result = await some_external_service()
        
        # 成功则保存结果
        await uow.workflows.update_outputs(workflow.id, result)
        
    except ExternalServiceError as e:
        # 标记为失败
        await uow.workflows.update_status(
            workflow.id,
            WorkflowStatus.FAILED,
            error_message=str(e)
        )
        raise
```

---

## 性能优化

### 1. 使用索引

```python
# Repository中已经对常用查询字段建立索引
# 参考 sql/schema.sql 中的索引定义
```

### 2. 分页查询

```python
# 使用limit和offset
workflows = await uow.workflows.get_multi(
    skip=page * page_size,
    limit=page_size
)
```

### 3. 选择性查询字段

```python
# 如果只需要部分字段，使用原生SQL
from sqlalchemy import select

result = await session.execute(
    select(WorkflowRecord.id, WorkflowRecord.subject)
    .where(WorkflowRecord.domain == 'k12_education')
)
```

### 4. 缓存频繁访问的数据

```python
from functools import lru_cache

@lru_cache(maxsize=100)
async def get_adapter_cached(adapter_id: str):
    async with get_uow() as uow:
        return await uow.adapters.get(adapter_id)
```

### 5. 使用连接池

```python
# database.py 中已配置
engine = create_async_engine(
    DATABASE_URL,
    pool_size=20,        # 连接池大小
    max_overflow=10,     # 最大溢出
    pool_pre_ping=True   # 连接健康检查
)
```

---

## 测试

### 单元测试

```python
import pytest
from src.db import get_uow

@pytest.mark.asyncio
async def test_create_workflow():
    async with get_uow() as uow:
        context = WorkflowContext(...)
        workflow = await uow.workflows.create_from_context(context)
        
        assert workflow.id == context.workflow_id
        assert workflow.status == 'pending'
```

### Mock Repository

```python
from unittest.mock import AsyncMock

@pytest.fixture
def mock_uow():
    uow = AsyncMock()
    uow.workflows = AsyncMock()
    uow.workflows.create_from_context.return_value = WorkflowRecord(...)
    return uow

async def test_service(mock_uow):
    service = WorkflowService(mock_uow)
    result = await service.create_workflow(...)
    
    mock_uow.workflows.create_from_context.assert_called_once()
```

### 集成测试

```python
@pytest.mark.integration
@pytest.mark.asyncio
async def test_full_workflow():
    async with get_uow() as uow:
        # 创建工作流
        workflow = await uow.workflows.create(...)
        
        # 添加节点
        node = await uow.nodes.create_node(...)
        
        # 添加学习目标
        objectives = await uow.learning_objectives.create_batch(...)
        
        # 验证
        retrieved = await uow.workflows.get_with_relations(workflow.id)
        assert len(retrieved.nodes) == 1
        assert len(retrieved.learning_objectives) > 0
```

### 运行测试

```bash
# 运行所有Repository测试
pytest tests/test_repositories.py -v

# 运行特定测试类
pytest tests/test_repositories.py::TestWorkflowRepository -v

# 运行特定测试方法
pytest tests/test_repositories.py::TestWorkflowRepository::test_create_workflow -v

# 显示详细输出
pytest tests/test_repositories.py -v -s

# 生成覆盖率报告
pytest tests/test_repositories.py --cov=src/db --cov-report=html
```

---

## 总结

### 已实现功能

✅ **通用CRUD基类** (BaseRepository)  
✅ **9个专用Repository**  
✅ **工作单元模式** (UnitOfWork)  
✅ **自动事务管理**  
✅ **批量操作支持**  
✅ **复杂查询（过滤、排序、分页）**  
✅ **关联查询（eager loading）**  
✅ **统计聚合功能**  
✅ **全文搜索**  
✅ **完整的类型提示**  
✅ **7个使用示例**  
✅ **40+个单元测试**  

### 文件清单

| 文件 | 行数 | 说明 |
|------|------|------|
| `src/db/repositories.py` | 400+ | 核心Repository实现 |
| `src/db/adapter_repository.py` | 350+ | 适配器Repository |
| `src/db/unit_of_work.py` | 120+ | 工作单元模式 |
| `examples/repository_examples.py` | 300+ | 使用示例 |
| `tests/test_repositories.py` | 400+ | 单元测试 |

### 下一步

1. ⏳ 添加更多Repository（User、Role、Permission等）
2. ⏳ 实现缓存层（Redis集成）
3. ⏳ 添加软删除支持
4. ⏳ 实现审计日志自动记录
5. ⏳ 性能监控和慢查询分析

---

**文档版本**: 1.0  
**最后更新**: 2025-12-09  
**维护者**: MetaWorkflow Team
