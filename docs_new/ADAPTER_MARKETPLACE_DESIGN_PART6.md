# MetaWorkflow V2.0 适配器市场设计 - 第六部分（最终部分）

## 第五部分：统计分析与开发者工具

## 5.1 统计分析系统

### 5.1.1 ClickHouse数据仓库设计

```sql
-- 适配器使用事件表
CREATE TABLE adapter_usage_events (
    event_id UUID,
    event_time DateTime,
    event_type String,  -- install/uninstall/activate/deactivate/execute
    
    -- 适配器信息
    adapter_id String,
    adapter_name String,
    adapter_version String,
    
    -- 用户信息
    user_id UInt64,
    organization String,
    
    -- 上下文信息
    domain String,
    activity_type String,
    workflow_id String,
    
    -- 性能指标
    execution_time_ms UInt32,
    memory_usage_mb UInt32,
    
    -- 结果
    success Boolean,
    error_message String,
    
    -- 设备信息
    platform String,
    platform_version String,
    
    -- 时间分区
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, event_time);


-- 适配器下载事件表
CREATE TABLE adapter_download_events (
    download_id UUID,
    download_time DateTime,
    
    adapter_id String,
    adapter_version String,
    
    user_id UInt64,
    
    -- 来源
    source String,  -- marketplace/api/cli
    referrer String,
    
    -- 地理位置
    country String,
    region String,
    
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, download_time);


-- 适配器评分事件表
CREATE TABLE adapter_rating_events (
    rating_id UInt64,
    rating_time DateTime,
    
    adapter_id String,
    user_id UInt64,
    
    rating UInt8,  -- 1-5
    previous_rating UInt8,  -- 如果是更新评分
    
    date Date
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (adapter_id, rating_time);
```

### 5.1.2 统计分析服务

```python
# src/marketplace/analytics/stats.py

from clickhouse_driver import Client
from datetime import datetime, timedelta
from typing import Dict, List, Any

class AdapterAnalyticsService:
    """适配器统计分析服务"""
    
    def __init__(self, clickhouse_client: Client):
        self.ch = clickhouse_client
    
    async def get_adapter_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict[str, Any]:
        """
        获取适配器统计数据
        
        Args:
            adapter_id: 适配器ID
            start_date: 开始日期
            end_date: 结束日期
            
        Returns:
            Dict: 统计数据
        """
        stats = {
            "overview": await self._get_overview_stats(adapter_id, start_date, end_date),
            "downloads": await self._get_download_stats(adapter_id, start_date, end_date),
            "usage": await self._get_usage_stats(adapter_id, start_date, end_date),
            "performance": await self._get_performance_stats(adapter_id, start_date, end_date),
            "rating_trend": await self._get_rating_trend(adapter_id, start_date, end_date),
            "user_demographics": await self._get_user_demographics(adapter_id, start_date, end_date)
        }
        
        return stats
    
    async def _get_overview_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """概览统计"""
        query = """
        SELECT
            countIf(event_type = 'install') as total_installs,
            countIf(event_type = 'uninstall') as total_uninstalls,
            countIf(event_type = 'execute') as total_executions,
            countIf(event_type = 'execute' AND success = 1) as successful_executions,
            countIf(event_type = 'execute' AND success = 0) as failed_executions,
            uniq(user_id) as unique_users,
            avg(execution_time_ms) as avg_execution_time,
            quantile(0.95)(execution_time_ms) as p95_execution_time
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        row = result[0]
        
        return {
            "total_installs": row[0],
            "total_uninstalls": row[1],
            "net_installs": row[0] - row[1],
            "total_executions": row[2],
            "successful_executions": row[3],
            "failed_executions": row[4],
            "success_rate": row[3] / row[2] if row[2] > 0 else 0,
            "unique_users": row[5],
            "avg_execution_time_ms": round(row[6], 2),
            "p95_execution_time_ms": round(row[7], 2)
        }
    
    async def _get_download_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """下载统计（按天）"""
        query = """
        SELECT
            date,
            count() as downloads,
            uniq(user_id) as unique_users
        FROM adapter_download_events
        WHERE adapter_id = %(adapter_id)s
          AND download_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY date
        ORDER BY date
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "daily_data": [
                {
                    "date": row[0].isoformat(),
                    "downloads": row[1],
                    "unique_users": row[2]
                }
                for row in result
            ]
        }
    
    async def _get_usage_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """使用统计"""
        # 按领域统计
        domain_query = """
        SELECT
            domain,
            count() as executions,
            uniq(user_id) as users
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY domain
        ORDER BY executions DESC
        LIMIT 10
        """
        
        domain_result = self.ch.execute(
            domain_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        # 按活动类型统计
        activity_query = """
        SELECT
            activity_type,
            count() as executions
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY activity_type
        ORDER BY executions DESC
        """
        
        activity_result = self.ch.execute(
            activity_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "by_domain": [
                {"domain": row[0], "executions": row[1], "users": row[2]}
                for row in domain_result
            ],
            "by_activity_type": [
                {"activity_type": row[0], "executions": row[1]}
                for row in activity_result
            ]
        }
    
    async def _get_performance_stats(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """性能统计"""
        query = """
        SELECT
            toStartOfHour(event_time) as hour,
            avg(execution_time_ms) as avg_time,
            quantile(0.50)(execution_time_ms) as p50_time,
            quantile(0.95)(execution_time_ms) as p95_time,
            quantile(0.99)(execution_time_ms) as p99_time,
            avg(memory_usage_mb) as avg_memory
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY hour
        ORDER BY hour
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "hourly_performance": [
                {
                    "hour": row[0].isoformat(),
                    "avg_time_ms": round(row[1], 2),
                    "p50_time_ms": round(row[2], 2),
                    "p95_time_ms": round(row[3], 2),
                    "p99_time_ms": round(row[4], 2),
                    "avg_memory_mb": round(row[5], 2)
                }
                for row in result
            ]
        }
    
    async def _get_rating_trend(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """评分趋势"""
        query = """
        SELECT
            toStartOfDay(rating_time) as day,
            avg(rating) as avg_rating,
            count() as rating_count
        FROM adapter_rating_events
        WHERE adapter_id = %(adapter_id)s
          AND rating_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY day
        ORDER BY day
        """
        
        result = self.ch.execute(
            query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "daily_ratings": [
                {
                    "date": row[0].isoformat(),
                    "avg_rating": round(row[1], 2),
                    "rating_count": row[2]
                }
                for row in result
            ]
        }
    
    async def _get_user_demographics(
        self,
        adapter_id: str,
        start_date: datetime,
        end_date: datetime
    ) -> Dict:
        """用户画像"""
        # 按组织统计
        org_query = """
        SELECT
            organization,
            uniq(user_id) as user_count,
            count() as execution_count
        FROM adapter_usage_events
        WHERE adapter_id = %(adapter_id)s
          AND event_type = 'execute'
          AND event_time BETWEEN %(start_date)s AND %(end_date)s
        GROUP BY organization
        ORDER BY execution_count DESC
        LIMIT 10
        """
        
        org_result = self.ch.execute(
            org_query,
            {
                'adapter_id': adapter_id,
                'start_date': start_date,
                'end_date': end_date
            }
        )
        
        return {
            "top_organizations": [
                {
                    "organization": row[0],
                    "user_count": row[1],
                    "execution_count": row[2]
                }
                for row in org_result
            ]
        }


# 全局分析服务实例
analytics_service = AdapterAnalyticsService(clickhouse_client)
```

### 5.1.3 统计API

```python
# src/marketplace/api/analytics.py

from fastapi import APIRouter, Depends, HTTPException, Query
from datetime import datetime, timedelta

router = APIRouter(prefix="/api/v1/marketplace/analytics", tags=["analytics"])


@router.get("/adapters/{adapter_id}/stats")
async def get_adapter_statistics(
    adapter_id: str,
    start_date: Optional[datetime] = Query(None),
    end_date: Optional[datetime] = Query(None),
    current_user: User = Depends(get_current_user)
):
    """
    获取适配器统计数据
    
    仅适配器作者或管理员可访问
    """
    # 验证权限
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=403,
            detail="Only adapter author or admin can view statistics"
        )
    
    # 默认时间范围：最近30天
    if not end_date:
        end_date = datetime.utcnow()
    if not start_date:
        start_date = end_date - timedelta(days=30)
    
    from src.marketplace.analytics.stats import analytics_service
    
    stats = await analytics_service.get_adapter_stats(
        adapter_id=adapter_id,
        start_date=start_date,
        end_date=end_date
    )
    
    return stats


@router.get("/marketplace/overview")
async def get_marketplace_overview(
    current_user: User = Depends(get_admin_user)
):
    """
    获取市场整体概览
    
    仅管理员可访问
    """
    # 统计各类数据
    total_adapters = await db.query(MarketplaceAdapter).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).count()
    
    total_developers = await db.query(User).filter(
        User.is_developer == True
    ).count()
    
    total_downloads = await db.query(
        func.sum(MarketplaceAdapter.download_count)
    ).scalar()
    
    total_installs = await db.query(
        func.sum(MarketplaceAdapter.install_count)
    ).scalar()
    
    avg_rating = await db.query(
        func.avg(MarketplaceAdapter.rating_average)
    ).filter(
        MarketplaceAdapter.rating_count > 0
    ).scalar()
    
    # 分类统计
    category_stats = await db.query(
        MarketplaceAdapter.category,
        func.count(MarketplaceAdapter.id).label('count')
    ).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.category
    ).all()
    
    return {
        "overview": {
            "total_adapters": total_adapters,
            "total_developers": total_developers,
            "total_downloads": total_downloads or 0,
            "total_installs": total_installs or 0,
            "avg_rating": round(avg_rating, 2) if avg_rating else 0.0
        },
        "by_category": [
            {"category": cat, "count": count}
            for cat, count in category_stats
        ]
    }
```

---

## 5.2 开发者控制台

### 5.2.1 开发者仪表板API

```python
# src/marketplace/api/developer.py

from fastapi import APIRouter, Depends, HTTPException
from typing import List

router = APIRouter(prefix="/api/v1/marketplace/developer", tags=["developer"])


@router.get("/dashboard")
async def get_developer_dashboard(
    current_user: User = Depends(get_developer_user)
):
    """
    获取开发者仪表板数据
    
    包含：
    - 适配器列表
    - 总体统计
    - 最近评论
    - 待处理任务
    """
    # 获取开发者的所有适配器
    adapters = await db.query(MarketplaceAdapter).filter(
        MarketplaceAdapter.author_id == current_user.id
    ).all()
    
    adapter_summaries = []
    total_downloads = 0
    total_installs = 0
    total_ratings = 0
    
    for adapter in adapters:
        total_downloads += adapter.download_count
        total_installs += adapter.install_count
        total_ratings += adapter.rating_count
        
        adapter_summaries.append({
            "adapter_id": adapter.adapter_id,
            "name": adapter.display_name,
            "status": adapter.status.value,
            "version": adapter.current_version,
            "downloads": adapter.download_count,
            "installs": adapter.active_install_count,
            "rating": adapter.rating_average,
            "rating_count": adapter.rating_count,
            "published_at": adapter.published_at.isoformat() if adapter.published_at else None,
            "updated_at": adapter.updated_at.isoformat()
        })
    
    # 获取最近的评论
    recent_reviews = await db.query(AdapterReview).join(
        MarketplaceAdapter
    ).filter(
        MarketplaceAdapter.author_id == current_user.id,
        AdapterReview.is_visible == True
    ).order_by(
        desc(AdapterReview.created_at)
    ).limit(10).all()
    
    # 待审核的版本
    pending_versions = await db.query(AdapterVersion).join(
        MarketplaceAdapter
    ).filter(
        MarketplaceAdapter.author_id == current_user.id,
        AdapterVersion.review_status.in_([
            AdapterPublishStatus.SUBMITTED,
            AdapterPublishStatus.REVIEWING
        ])
    ).all()
    
    return {
        "summary": {
            "total_adapters": len(adapters),
            "published_adapters": len([a for a in adapters if a.status == AdapterPublishStatus.PUBLISHED]),
            "total_downloads": total_downloads,
            "total_installs": total_installs,
            "total_ratings": total_ratings,
            "avg_rating": sum(a.rating_average for a in adapters) / len(adapters) if adapters else 0
        },
        "adapters": adapter_summaries,
        "recent_reviews": [
            {
                "adapter_name": review.adapter.display_name,
                "rating": review.rating,
                "title": review.title,
                "user": review.user.display_name or review.user.username,
                "created_at": review.created_at.isoformat()
            }
            for review in recent_reviews
        ],
        "pending_reviews": [
            {
                "adapter_name": version.adapter.display_name,
                "version": version.version,
                "status": version.review_status.value,
                "submitted_at": version.created_at.isoformat()
            }
            for version in pending_versions
        ]
    }


@router.post("/adapters/{adapter_id}/reply-review/{review_id}")
async def reply_to_review(
    adapter_id: str,
    review_id: int,
    reply_text: str = Body(..., max_length=2000),
    current_user: User = Depends(get_developer_user)
):
    """
    回复用户评论
    
    只有适配器作者可以回复
    """
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id:
        raise HTTPException(
            status_code=403,
            detail="Only adapter author can reply to reviews"
        )
    
    review = await db.query(AdapterReview).get(review_id)
    
    if not review or review.adapter_id != adapter.id:
        raise HTTPException(status_code=404, detail="Review not found")
    
    # 添加回复
    review.author_reply = reply_text
    review.author_replied_at = datetime.utcnow()
    
    await db.commit()
    
    # 通知用户
    await notify_user_of_reply(review.user_id, adapter.display_name)
    
    return {"message": "Reply added successfully"}
```

---

## 5.3 开发者SDK和CLI工具

### 5.3.1 Python SDK

```python
# metaworkflow_sdk/client.py

"""
MetaWorkflow 适配器市场 SDK

安装：
    pip install metaworkflow-sdk

使用示例：
    from metaworkflow_sdk import MarketplaceClient
    
    client = MarketplaceClient(api_key="your_api_key")
    
    # 发布适配器
    client.publish_adapter(
        name="my-chemistry-adapter",
        version="1.0.0",
        package_path="./dist/my-adapter.zip"
    )
"""

import requests
from typing import Optional, Dict, Any
from pathlib import Path

class MarketplaceClient:
    """市场API客户端"""
    
    def __init__(
        self,
        api_key: str,
        base_url: str = "https://marketplace.metaworkflow.ai/api/v1"
    ):
        self.api_key = api_key
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {api_key}",
            "User-Agent": "MetaWorkflow-SDK/1.0"
        })
    
    def create_adapter(
        self,
        name: str,
        display_name: str,
        description: str,
        category: str,
        **kwargs
    ) -> Dict[str, Any]:
        """
        创建适配器（草稿）
        
        Args:
            name: 适配器名称（唯一标识）
            display_name: 显示名称
            description: 描述
            category: 分类
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 创建结果
        """
        payload = {
            "name": name,
            "display_name": display_name,
            "description": description,
            "category": category,
            **kwargs
        }
        
        response = self.session.post(
            f"{self.base_url}/marketplace/publish/adapters",
            json=payload
        )
        response.raise_for_status()
        
        return response.json()
    
    def upload_version(
        self,
        adapter_id: str,
        version: str,
        package_path: str,
        changelog: str,
        **kwargs
    ) -> Dict[str, Any]:
        """
        上传适配器版本
        
        Args:
            adapter_id: 适配器ID
            version: 版本号
            package_path: 包文件路径
            changelog: 变更日志
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 上传结果
        """
        package_file = Path(package_path)
        
        if not package_file.exists():
            raise FileNotFoundError(f"Package file not found: {package_path}")
        
        # 上传版本信息
        version_data = {
            "version": version,
            "changelog": changelog,
            **kwargs
        }
        
        # 上传文件
        with open(package_file, 'rb') as f:
            files = {'package_file': (package_file.name, f, 'application/zip')}
            
            response = self.session.post(
                f"{self.base_url}/marketplace/publish/adapters/{adapter_id}/versions",
                data=version_data,
                files=files
            )
        
        response.raise_for_status()
        
        return response.json()
    
    def submit_for_review(self, adapter_id: str) -> Dict[str, Any]:
        """
        提交审核
        
        Args:
            adapter_id: 适配器ID
        
        Returns:
            Dict: 提交结果
        """
        response = self.session.post(
            f"{self.base_url}/marketplace/publish/adapters/{adapter_id}/submit"
        )
        response.raise_for_status()
        
        return response.json()
    
    def publish_adapter(
        self,
        name: str,
        display_name: str,
        description: str,
        category: str,
        version: str,
        package_path: str,
        changelog: str,
        auto_submit: bool = True,
        **kwargs
    ) -> Dict[str, Any]:
        """
        一键发布适配器
        
        组合了创建、上传、提交审核三个步骤
        
        Args:
            name: 适配器名称
            display_name: 显示名称
            description: 描述
            category: 分类
            version: 版本号
            package_path: 包文件路径
            changelog: 变更日志
            auto_submit: 是否自动提交审核
            **kwargs: 其他可选参数
        
        Returns:
            Dict: 发布结果
        """
        # 1. 创建适配器
        print(f"Creating adapter: {name}")
        create_result = self.create_adapter(
            name=name,
            display_name=display_name,
            description=description,
            category=category,
            **kwargs
        )
        
        adapter_id = create_result['adapter_id']
        print(f"✓ Adapter created: {adapter_id}")
        
        # 2. 上传版本
        print(f"Uploading version: {version}")
        upload_result = self.upload_version(
            adapter_id=adapter_id,
            version=version,
            package_path=package_path,
            changelog=changelog
        )
        print(f"✓ Version uploaded: {upload_result['version']}")
        
        # 3. 提交审核
        if auto_submit:
            print("Submitting for review...")
            submit_result = self.submit_for_review(adapter_id)
            print(f"✓ Submitted for review: {submit_result['status']}")
            
            return {
                "adapter_id": adapter_id,
                "version": version,
                "status": submit_result['status'],
                "message": "Adapter published and submitted for review"
            }
        else:
            return {
                "adapter_id": adapter_id,
                "version": version,
                "status": "draft",
                "message": "Adapter created. Use submit_for_review() to submit."
            }
    
    def get_adapter_stats(
        self,
        adapter_id: str,
        start_date: Optional[str] = None,
        end_date: Optional[str] = None
    ) -> Dict[str, Any]:
        """
        获取适配器统计数据
        
        Args:
            adapter_id: 适配器ID
            start_date: 开始日期 (ISO格式)
            end_date: 结束日期 (ISO格式)
        
        Returns:
            Dict: 统计数据
        """
        params = {}
        if start_date:
            params['start_date'] = start_date
        if end_date:
            params['end_date'] = end_date
        
        response = self.session.get(
            f"{self.base_url}/marketplace/analytics/adapters/{adapter_id}/stats",
            params=params
        )
        response.raise_for_status()
        
        return response.json()
```

### 5.3.2 CLI工具

```python
# metaworkflow_cli/main.py

"""
MetaWorkflow CLI 工具

安装：
    pip install metaworkflow-cli

命令：
    mw-cli init              # 初始化适配器项目
    mw-cli build             # 构建适配器包
    mw-cli publish           # 发布到市场
    mw-cli stats <adapter>   # 查看统计数据
"""

import click
import yaml
from pathlib import Path
from metaworkflow_sdk import MarketplaceClient

@click.group()
def cli():
    """MetaWorkflow 命令行工具"""
    pass


@cli.command()
@click.argument('name')
@click.option('--template', default='basic', help='项目模板')
def init(name: str, template: str):
    """初始化新的适配器项目"""
    click.echo(f"Creating adapter project: {name}")
    
    project_dir = Path(name)
    project_dir.mkdir(exist_ok=True)
    
    # 创建目录结构
    (project_dir / "src").mkdir(exist_ok=True)
    (project_dir / "tests").mkdir(exist_ok=True)
    (project_dir / "docs").mkdir(exist_ok=True)
    
    # 创建manifest.yaml
    manifest = {
        "name": name,
        "display_name": name.replace('-', ' ').title(),
        "version": "0.1.0",
        "description": f"Description for {name}",
        "author": {
            "name": "Your Name",
            "email": "your.email@example.com"
        },
        "compatibility": {
            "platform_version": ">=2.0.0"
        }
    }
    
    with open(project_dir / "manifest.yaml", 'w') as f:
        yaml.dump(manifest, f, default_flow_style=False)
    
    # 创建adapter.py模板
    adapter_template = '''from metaworkflow.adapters import DomainAdapter

class MyAdapter(DomainAdapter):
    """My custom adapter"""
    
    async def parse_request(self, request):
        # TODO: Implement parsing logic
        pass
    
    async def discover_content(self, topic, context):
        # TODO: Implement content discovery
        pass
'''
    
    with open(project_dir / "src" / "adapter.py", 'w') as f:
        f.write(adapter_template)
    
    click.echo(f"✓ Project created at {project_dir}")
    click.echo("\nNext steps:")
    click.echo(f"  cd {name}")
    click.echo("  # Edit src/adapter.py and manifest.yaml")
    click.echo("  mw-cli build")


@cli.command()
@click.option('--output', '-o', default='dist', help='输出目录')
def build(output: str):
    """构建适配器包"""
    import zipfile
    import shutil
    
    click.echo("Building adapter package...")
    
    # 读取manifest
    if not Path('manifest.yaml').exists():
        click.echo("Error: manifest.yaml not found", err=True)
        return
    
    with open('manifest.yaml', 'r') as f:
        manifest = yaml.safe_load(f)
    
    name = manifest['name']
    version = manifest['version']
    
    # 创建输出目录
    output_dir = Path(output)
    output_dir.mkdir(exist_ok=True)
    
    # 创建zip包
    package_path = output_dir / f"{name}-{version}.zip"
    
    with zipfile.ZipFile(package_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
        # 添加manifest
        zipf.write('manifest.yaml')
        
        # 添加源代码
        for py_file in Path('src').rglob('*.py'):
            zipf.write(py_file, py_file)
        
        # 添加文档
        if Path('README.md').exists():
            zipf.write('README.md')
    
    click.echo(f"✓ Package built: {package_path}")
    click.echo(f"  Size: {package_path.stat().st_size / 1024:.1f} KB")


@cli.command()
@click.option('--api-key', envvar='METAWORKFLOW_API_KEY', help='API密钥')
@click.option('--auto-submit/--no-auto-submit', default=True, help='自动提交审核')
def publish(api_key: str, auto_submit: bool):
    """发布适配器到市场"""
    if not api_key:
        click.echo("Error: API key required. Set METAWORKFLOW_API_KEY or use --api-key", err=True)
        return
    
    # 读取manifest
    with open('manifest.yaml', 'r') as f:
        manifest = yaml.safe_load(f)
    
    name = manifest['name']
    version = manifest['version']
    
    # 查找包文件
    package_path = Path('dist') / f"{name}-{version}.zip"
    
    if not package_path.exists():
        click.echo(f"Error: Package not found: {package_path}", err=True)
        click.echo("Run 'mw-cli build' first", err=True)
        return
    
    click.echo(f"Publishing {name} v{version}...")
    
    # 使用SDK发布
    client = MarketplaceClient(api_key=api_key)
    
    try:
        result = client.publish_adapter(
            name=name,
            display_name=manifest['display_name'],
            description=manifest['description'],
            category=manifest.get('category', 'other'),
            version=version,
            package_path=str(package_path),
            changelog=manifest.get('changelog', 'Initial release'),
            auto_submit=auto_submit
        )
        
        click.echo(f"✓ {result['message']}")
        click.echo(f"  Adapter ID: {result['adapter_id']}")
        click.echo(f"  Status: {result['status']}")
        
    except Exception as e:
        click.echo(f"Error: {str(e)}", err=True)


@cli.command()
@click.argument('adapter_id')
@click.option('--api-key', envvar='METAWORKFLOW_API_KEY', help='API密钥')
@click.option('--days', default=30, help='统计天数')
def stats(adapter_id: str, api_key: str, days: int):
    """查看适配器统计数据"""
    if not api_key:
        click.echo("Error: API key required", err=True)
        return
    
    from datetime import datetime, timedelta
    
    client = MarketplaceClient(api_key=api_key)
    
    end_date = datetime.utcnow()
    start_date = end_date - timedelta(days=days)
    
    stats_data = client.get_adapter_stats(
        adapter_id=adapter_id,
        start_date=start_date.isoformat(),
        end_date=end_date.isoformat()
    )
    
    overview = stats_data['overview']
    
    click.echo(f"\nStatistics for {adapter_id} (Last {days} days)")
    click.echo("=" * 50)
    click.echo(f"Downloads:        {overview['total_installs']}")
    click.echo(f"Active Installs:  {overview['net_installs']}")
    click.echo(f"Executions:       {overview['total_executions']}")
    click.echo(f"Success Rate:     {overview['success_rate']*100:.1f}%")
    click.echo(f"Unique Users:     {overview['unique_users']}")
    click.echo(f"Avg Exec Time:    {overview['avg_execution_time_ms']:.0f}ms")
    click.echo(f"P95 Exec Time:    {overview['p95_execution_time_ms']:.0f}ms")


if __name__ == '__main__':
    cli()
```

---

## 5.4 总结与最佳实践

### 5.4.1 架构设计总结

**核心设计原则**：

1. **三层架构（孵化器-鸡-蛋模型）**
   - **平台核心层**：提供通用AI能力（孵化器）
   - **适配器层**：领域专业化（鸡）
   - **内容层**：最终交付物（蛋）

2. **插件化扩展**
   - 标准化接口（DomainAdapter、NodeExecutor、ContentSource）
   - 动态加载机制
   - 依赖管理和版本控制

3. **质量保证**
   - 双重审核（自动检测 + 人工审核）
   - 安全扫描、性能测试、兼容性验证
   - 持续监控和质量评分

4. **开发者友好**
   - 完整的SDK和CLI工具
   - 清晰的文档和示例
   - 快速的发布流程

### 5.4.2 开发最佳实践

**适配器开发**：

```python
# 1. 良好的结构
my-adapter/
├── manifest.yaml          # 元数据声明
├── src/
│   ├── adapter.py        # 主适配器类
│   ├── content_sources/  # 内容源
│   ├── node_executors/   # 节点执行器
│   └── utils/            # 工具函数
├── tests/                 # 测试
├── docs/                  # 文档
└── README.md

# 2. 清晰的接口实现
class MyAdapter(DomainAdapter):
    """
    简洁的文档字符串
    
    说明：
    - 适配器用途
    - 支持的功能
    - 使用示例
    """
    
    async def parse_request(self, request):
        # 详细的注释
        # 清晰的错误处理
        # 合理的日志记录
        pass

# 3. 完善的测试
# tests/test_adapter.py
import pytest

@pytest.mark.asyncio
async def test_parse_request():
    adapter = MyAdapter()
    result = await adapter.parse_request({...})
    assert result is not None

# 4. 版本管理
# 遵循语义化版本（Semantic Versioning）
# MAJOR.MINOR.PATCH
# 1.0.0 -> 1.0.1 (bug fix)
# 1.0.1 -> 1.1.0 (new feature)
# 1.1.0 -> 2.0.0 (breaking change)
```

**性能优化**：

```python
# 1. 异步操作
async def discover_content(self, topic):
    # 并发请求多个数据源
    results = await asyncio.gather(
        self.source1.search(topic),
        self.source2.search(topic),
        self.source3.search(topic)
    )
    return self._merge_results(results)

# 2. 缓存机制
from functools import lru_cache

@lru_cache(maxsize=128)
def _get_static_data(key):
    # 缓存静态数据
    pass

# 3. 批量处理
async def process_batch(self, items):
    # 批量处理，减少API调用
    batch_size = 10
    for i in range(0, len(items), batch_size):
        batch = items[i:i+batch_size]
        await self._process(batch)
```

**安全最佳实践**：

```python
# 1. 输入验证
def parse_request(self, request):
    # 验证输入
    if not isinstance(request, dict):
        raise ValueError("Request must be a dict")
    
    # 清理输入
    topic = request.get('topic', '').strip()
    if not topic:
        raise ValueError("Topic is required")
    
    return topic

# 2. 敏感信息处理
# 永远不要硬编码密钥
# ❌ 错误
API_KEY = "sk-1234567890abcdef"

# ✅ 正确
import os
API_KEY = os.environ.get('EXTERNAL_API_KEY')

# 3. 错误处理
try:
    result = await external_api.call()
except Exception as e:
    # 记录错误，但不暴露敏感信息
    logger.error(f"API call failed: {type(e).__name__}")
    raise AdapterError("Failed to fetch data")
```

### 5.4.3 市场运营建议

**for 平台运营方**：

1. **质量控制**
   - 严格的审核标准
   - 定期安全扫描
   - 用户反馈机制

2. **开发者支持**
   - 详细文档和教程
   - 活跃的开发者社区
   - 技术支持渠道

3. **激励机制**
   - 优质适配器推荐
   - 开发者排行榜
   - 商业化分成（未来）

**for 适配器开发者**：

1. **用户第一**
   - 清晰的文档
   - 快速的响应
   - 持续的更新

2. **质量优先**
   - 充分测试
   - 性能优化
   - 安全防护

3. **社区参与**
   - 及时回复评论
   - 分享最佳实践
   - 贡献示例代码

---

## 6. 文档总结

### 6.1 完整文档结构

本设计文档共**6个部分**，涵盖：

**Part 1 - 架构基础**：
- 三层架构模型
- 适配器包结构
- manifest.yaml规范
- 元数据管理

**Part 2 - 接口规范**：
- DomainAdapter接口（15个方法）
- NodeExecutor接口
- ContentSource接口
- 详细示例

**Part 3 - 注册与生命周期**：
- 注册表设计
- 自动发现机制
- 依赖管理
- 生命周期状态机

**Part 4 - 市场平台**：
- 4层架构
- 完整数据模型
- 发布流程
- 自动化检测

**Part 5 - 用户体验**：
- Elasticsearch搜索
- 推荐系统
- 评分评论
- 统计分析

**Part 6 - 工具与总结**：
- ClickHouse分析
- 开发者仪表板
- Python SDK
- CLI工具
- 最佳实践

### 6.2 关键指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| 适配器审核时间 | < 48小时 | 从提交到审核完成 |
| 安装成功率 | > 99% | 用户安装失败率 < 1% |
| 平台可用性 | 99.9% | 年停机时间 < 8.76小时 |
| 搜索响应时间 | < 200ms | P95响应时间 |
| 推荐准确率 | > 30% | 用户点击率 |
| 开发者满意度 | > 4.0/5.0 | 开发者调查评分 |

### 6.3 实施路线图

**Phase 1 - 基础设施（2周）**：
- ✅ 数据库设计和迁移
- ✅ API框架搭建
- ✅ 认证授权系统

**Phase 2 - 核心功能（4周）**：
- ✅ 适配器注册和管理
- ✅ 发布和审核流程
- ✅ 搜索功能

**Phase 3 - 增强功能（3周）**：
- ✅ 推荐系统
- ✅ 评分评论
- ✅ 统计分析

**Phase 4 - 开发者工具（2周）**：
- ✅ SDK开发
- ✅ CLI工具
- ✅ 文档完善

**Phase 5 - 上线准备（1周）**：
- ⏳ 性能测试
- ⏳ 安全审计
- ⏳ 生产部署

### 6.4 成功标准

**技术指标**：
- [x] 完整的架构设计
- [x] 标准化接口定义
- [x] 自动化质量检测
- [x] 可扩展性设计

**业务指标**：
- [ ] 上线3个月：20+ 适配器
- [ ] 上线6个月：50+ 适配器
- [ ] 上线12个月：100+ 适配器
- [ ] 开发者社区：500+ 成员

**用户体验**：
- [x] 简单的发布流程
- [x] 智能推荐
- [x] 快速搜索
- [x] 详细统计

---

## 结语

MetaWorkflow V2.0 适配器市场设计遵循**开放、标准、质量、生态**的核心理念：

- **开放**：欢迎第三方开发者贡献适配器
- **标准**：统一的接口和规范，保证兼容性
- **质量**：严格的审核和持续的监控
- **生态**：构建健康的开发者和用户生态

通过适配器市场，MetaWorkflow可以快速扩展到新的教育领域，从**K12基础教育**延伸到**高等教育、职业培训、企业培训**等多个场景，真正实现**"孵化器孵化无限可能"**的愿景。

---

**文档版本**: v1.0  
**最后更新**: 2025-12-09  
**作者**: MetaWorkflow架构团队  
**状态**: ✅ 完成

*感谢您阅读完整的适配器市场设计文档！*
