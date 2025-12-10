# MetaWorkflow V2.0 适配器市场设计 - 第五部分

## 第四部分：搜索、推荐与用户体验

## 4.1 搜索系统

### 4.1.1 Elasticsearch索引设计

```python
# src/marketplace/search/indexing.py

from elasticsearch import Elasticsearch
from typing import Dict, List, Any

# Elasticsearch索引映射
ADAPTER_INDEX_MAPPING = {
    "mappings": {
        "properties": {
            # 基本信息
            "adapter_id": {"type": "keyword"},
            "name": {"type": "keyword"},
            "display_name": {
                "type": "text",
                "fields": {
                    "keyword": {"type": "keyword"},
                    "suggest": {"type": "completion"}
                },
                "analyzer": "standard"
            },
            "description": {
                "type": "text",
                "analyzer": "standard"
            },
            "tagline": {
                "type": "text",
                "analyzer": "standard"
            },
            
            # 分类和标签
            "category": {"type": "keyword"},
            "subcategory": {"type": "keyword"},
            "tags": {"type": "keyword"},
            
            # 能力
            "supported_domains": {"type": "keyword"},
            "supported_activity_types": {"type": "keyword"},
            "supported_output_formats": {"type": "keyword"},
            
            # 作者
            "author_name": {
                "type": "text",
                "fields": {"keyword": {"type": "keyword"}}
            },
            "organization": {"type": "keyword"},
            
            # 质量指标（用于排序）
            "quality_score": {"type": "float"},
            "rating_average": {"type": "float"},
            "rating_count": {"type": "integer"},
            "download_count": {"type": "integer"},
            "install_count": {"type": "integer"},
            "active_install_count": {"type": "integer"},
            
            # 状态
            "status": {"type": "keyword"},
            "is_free": {"type": "boolean"},
            "is_open_source": {"type": "boolean"},
            
            # 时间
            "published_at": {"type": "date"},
            "updated_at": {"type": "date"},
            
            # 版本信息
            "current_version": {"type": "keyword"},
            "platform_version_min": {"type": "keyword"},
            
            # 全文搜索字段（组合多个字段）
            "all_text": {
                "type": "text",
                "analyzer": "standard"
            }
        }
    },
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 2,
        "analysis": {
            "analyzer": {
                "standard": {
                    "type": "standard",
                    "stopwords": "_english_"
                }
            }
        }
    }
}


class AdapterSearchService:
    """适配器搜索服务"""
    
    def __init__(self, es_client: Elasticsearch):
        self.es = es_client
        self.index_name = "marketplace_adapters"
    
    async def index_adapter(self, adapter: MarketplaceAdapter) -> None:
        """
        索引适配器
        
        Args:
            adapter: 适配器对象
        """
        doc = {
            "adapter_id": adapter.adapter_id,
            "name": adapter.name,
            "display_name": adapter.display_name,
            "description": adapter.description,
            "tagline": adapter.tagline,
            "category": adapter.category,
            "subcategory": adapter.subcategory,
            "tags": adapter.tags or [],
            "supported_domains": adapter.supported_domains or [],
            "supported_activity_types": adapter.supported_activity_types or [],
            "supported_output_formats": adapter.supported_output_formats or [],
            "author_name": adapter.author.display_name or adapter.author.username,
            "organization": adapter.organization,
            "quality_score": adapter.quality_score or 0.0,
            "rating_average": adapter.rating_average,
            "rating_count": adapter.rating_count,
            "download_count": adapter.download_count,
            "install_count": adapter.install_count,
            "active_install_count": adapter.active_install_count,
            "status": adapter.status.value,
            "is_free": adapter.is_free,
            "is_open_source": adapter.is_open_source,
            "published_at": adapter.published_at.isoformat() if adapter.published_at else None,
            "updated_at": adapter.updated_at.isoformat(),
            "current_version": adapter.current_version,
            "platform_version_min": adapter.platform_version_min,
            # 组合搜索字段
            "all_text": f"{adapter.display_name} {adapter.description} {adapter.tagline} {' '.join(adapter.tags or [])}"
        }
        
        await self.es.index(
            index=self.index_name,
            id=adapter.adapter_id,
            document=doc
        )
    
    async def search(
        self,
        query: str = None,
        category: str = None,
        tags: List[str] = None,
        domains: List[str] = None,
        activity_types: List[str] = None,
        min_rating: float = None,
        is_free: bool = None,
        sort_by: str = "relevance",
        page: int = 1,
        page_size: int = 20
    ) -> Dict[str, Any]:
        """
        搜索适配器
        
        Args:
            query: 搜索关键词
            category: 分类过滤
            tags: 标签过滤
            domains: 领域过滤
            activity_types: 活动类型过滤
            min_rating: 最低评分
            is_free: 是否免费
            sort_by: 排序方式 (relevance/rating/downloads/newest)
            page: 页码
            page_size: 每页数量
            
        Returns:
            Dict: 搜索结果
        """
        # 构建查询
        must_conditions = []
        filter_conditions = [
            {"term": {"status": "published"}}  # 只搜索已发布的
        ]
        
        # 关键词搜索
        if query:
            must_conditions.append({
                "multi_match": {
                    "query": query,
                    "fields": [
                        "display_name^3",      # display_name权重最高
                        "tagline^2",           # tagline次之
                        "description",
                        "tags",
                        "all_text"
                    ],
                    "type": "best_fields",
                    "fuzziness": "AUTO"        # 允许模糊匹配
                }
            })
        
        # 分类过滤
        if category:
            filter_conditions.append({"term": {"category": category}})
        
        # 标签过滤
        if tags:
            filter_conditions.append({"terms": {"tags": tags}})
        
        # 领域过滤
        if domains:
            filter_conditions.append({"terms": {"supported_domains": domains}})
        
        # 活动类型过滤
        if activity_types:
            filter_conditions.append({"terms": {"supported_activity_types": activity_types}})
        
        # 评分过滤
        if min_rating is not None:
            filter_conditions.append({
                "range": {"rating_average": {"gte": min_rating}}
            })
        
        # 免费过滤
        if is_free is not None:
            filter_conditions.append({"term": {"is_free": is_free}})
        
        # 构建完整查询
        search_query = {
            "bool": {
                "must": must_conditions if must_conditions else [{"match_all": {}}],
                "filter": filter_conditions
            }
        }
        
        # 排序
        sort_options = self._get_sort_options(sort_by)
        
        # 执行搜索
        from_index = (page - 1) * page_size
        
        response = await self.es.search(
            index=self.index_name,
            query=search_query,
            sort=sort_options,
            from_=from_index,
            size=page_size,
            track_total_hits=True
        )
        
        # 处理结果
        hits = response['hits']
        total = hits['total']['value']
        
        adapters = []
        for hit in hits['hits']:
            adapter_data = hit['_source']
            adapter_data['_score'] = hit['_score']
            adapters.append(adapter_data)
        
        return {
            "total": total,
            "page": page,
            "page_size": page_size,
            "total_pages": (total + page_size - 1) // page_size,
            "adapters": adapters
        }
    
    def _get_sort_options(self, sort_by: str) -> List[Dict]:
        """获取排序选项"""
        sort_map = {
            "relevance": [
                {"_score": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"download_count": {"order": "desc"}}
            ],
            "rating": [
                {"rating_average": {"order": "desc"}},
                {"rating_count": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "downloads": [
                {"download_count": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "newest": [
                {"published_at": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ],
            "popularity": [
                {"active_install_count": {"order": "desc"}},
                {"rating_average": {"order": "desc"}},
                {"_score": {"order": "desc"}}
            ]
        }
        
        return sort_map.get(sort_by, sort_map["relevance"])
    
    async def suggest(self, prefix: str, size: int = 10) -> List[str]:
        """
        自动补全建议
        
        Args:
            prefix: 输入前缀
            size: 返回数量
            
        Returns:
            List[str]: 建议列表
        """
        response = await self.es.search(
            index=self.index_name,
            suggest={
                "adapter_suggest": {
                    "prefix": prefix,
                    "completion": {
                        "field": "display_name.suggest",
                        "size": size,
                        "skip_duplicates": True
                    }
                }
            }
        )
        
        suggestions = []
        for option in response['suggest']['adapter_suggest'][0]['options']:
            suggestions.append(option['text'])
        
        return suggestions


# 全局搜索服务实例
search_service = AdapterSearchService(es_client)
```

### 4.1.2 搜索API

```python
# src/marketplace/api/search.py

from fastapi import APIRouter, Query
from typing import Optional, List

router = APIRouter(prefix="/api/v1/marketplace/search", tags=["search"])


@router.get("/adapters")
async def search_adapters(
    q: Optional[str] = Query(None, description="搜索关键词"),
    category: Optional[str] = Query(None, description="分类"),
    tags: Optional[List[str]] = Query(None, description="标签"),
    domains: Optional[List[str]] = Query(None, description="支持的领域"),
    activity_types: Optional[List[str]] = Query(None, description="活动类型"),
    min_rating: Optional[float] = Query(None, ge=0, le=5, description="最低评分"),
    is_free: Optional[bool] = Query(None, description="是否免费"),
    sort_by: str = Query("relevance", description="排序方式"),
    page: int = Query(1, ge=1, description="页码"),
    page_size: int = Query(20, ge=1, le=100, description="每页数量")
):
    """
    搜索适配器
    
    支持的排序方式：
    - relevance: 相关性（默认）
    - rating: 评分
    - downloads: 下载量
    - newest: 最新
    - popularity: 活跃安装量
    """
    from src.marketplace.search.indexing import search_service
    
    results = await search_service.search(
        query=q,
        category=category,
        tags=tags,
        domains=domains,
        activity_types=activity_types,
        min_rating=min_rating,
        is_free=is_free,
        sort_by=sort_by,
        page=page,
        page_size=page_size
    )
    
    return results


@router.get("/suggest")
async def get_suggestions(
    q: str = Query(..., min_length=1, description="输入前缀"),
    size: int = Query(10, ge=1, le=20, description="建议数量")
):
    """获取搜索自动补全建议"""
    from src.marketplace.search.indexing import search_service
    
    suggestions = await search_service.suggest(q, size)
    
    return {"suggestions": suggestions}


@router.get("/filters")
async def get_available_filters():
    """
    获取可用的过滤选项
    
    返回所有可用的分类、标签、领域等选项，
    用于前端构建过滤器UI。
    """
    # 从数据库聚合统计
    categories = await db.query(
        MarketplaceAdapter.category,
        func.count(MarketplaceAdapter.id).label('count')
    ).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.category
    ).all()
    
    # 统计所有标签
    all_tags = await db.query(MarketplaceAdapter.tags).filter(
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).all()
    
    tag_counts = {}
    for (tags,) in all_tags:
        if tags:
            for tag in tags:
                tag_counts[tag] = tag_counts.get(tag, 0) + 1
    
    return {
        "categories": [
            {"name": cat, "count": count}
            for cat, count in categories
        ],
        "tags": [
            {"name": tag, "count": count}
            for tag, count in sorted(tag_counts.items(), key=lambda x: x[1], reverse=True)
        ],
        "sort_options": [
            {"value": "relevance", "label": "相关性"},
            {"value": "rating", "label": "评分"},
            {"value": "downloads", "label": "下载量"},
            {"value": "newest", "label": "最新"},
            {"value": "popularity", "label": "热度"}
        ]
    }
```

---

## 4.2 推荐系统

### 4.2.1 推荐算法

```python
# src/marketplace/recommendations.py

from typing import List, Dict, Any
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

class AdapterRecommendationEngine:
    """适配器推荐引擎"""
    
    def __init__(self):
        self.adapter_features = {}  # 适配器特征向量缓存
    
    async def get_recommendations(
        self,
        user_id: int,
        context: Dict[str, Any] = None,
        count: int = 10
    ) -> List[Dict[str, Any]]:
        """
        获取个性化推荐
        
        推荐策略：
        1. 基于用户历史安装（协同过滤）
        2. 基于内容相似度
        3. 热门推荐（新用户冷启动）
        4. 上下文推荐（当前正在做什么）
        
        Args:
            user_id: 用户ID
            context: 上下文信息（当前领域、活动类型等）
            count: 推荐数量
            
        Returns:
            List[Dict]: 推荐的适配器列表
        """
        recommendations = []
        
        # 1. 获取用户安装历史
        user_installations = await self._get_user_installations(user_id)
        
        if user_installations:
            # 有安装历史：协同过滤 + 内容相似
            collaborative_recs = await self._collaborative_filtering(user_id, count // 2)
            content_recs = await self._content_based_filtering(user_installations, count // 2)
            
            recommendations.extend(collaborative_recs)
            recommendations.extend(content_recs)
        else:
            # 新用户：热门推荐
            recommendations = await self._popular_recommendations(count)
        
        # 2. 如果有上下文，添加上下文相关推荐
        if context:
            context_recs = await self._context_based_recommendations(context, count // 3)
            recommendations.extend(context_recs)
        
        # 3. 去重和排序
        recommendations = self._deduplicate_and_rank(recommendations, count)
        
        return recommendations
    
    async def _collaborative_filtering(
        self,
        user_id: int,
        count: int
    ) -> List[Dict]:
        """
        协同过滤推荐
        
        逻辑：找到和用户安装了相同适配器的其他用户，
        推荐他们还安装了什么。
        """
        # 获取用户安装的适配器
        user_adapters = await db.query(
            AdapterInstallation.adapter_id
        ).filter(
            AdapterInstallation.user_id == user_id,
            AdapterInstallation.status == 'active'
        ).all()
        
        user_adapter_ids = [a[0] for a in user_adapters]
        
        # 找到安装了相同适配器的其他用户
        similar_users = await db.query(
            AdapterInstallation.user_id,
            func.count(AdapterInstallation.adapter_id).label('common_count')
        ).filter(
            AdapterInstallation.adapter_id.in_(user_adapter_ids),
            AdapterInstallation.user_id != user_id,
            AdapterInstallation.status == 'active'
        ).group_by(
            AdapterInstallation.user_id
        ).order_by(
            desc('common_count')
        ).limit(50).all()
        
        similar_user_ids = [u[0] for u in similar_users]
        
        # 找到这些用户还安装了什么
        recommendations = await db.query(
            MarketplaceAdapter,
            func.count(AdapterInstallation.user_id).label('install_count')
        ).join(
            AdapterInstallation
        ).filter(
            AdapterInstallation.user_id.in_(similar_user_ids),
            AdapterInstallation.adapter_id.notin_(user_adapter_ids),
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
        ).group_by(
            MarketplaceAdapter.id
        ).order_by(
            desc('install_count')
        ).limit(count).all()
        
        return [
            {
                "adapter": adapter,
                "score": install_count / len(similar_user_ids),
                "reason": "similar_users"
            }
            for adapter, install_count in recommendations
        ]
    
    async def _content_based_filtering(
        self,
        user_installations: List[MarketplaceAdapter],
        count: int
    ) -> List[Dict]:
        """
        基于内容的推荐
        
        逻辑：找到和用户已安装适配器相似的其他适配器。
        相似度基于：分类、标签、支持的领域等。
        """
        # 构建用户偏好特征向量
        user_features = self._build_user_feature_vector(user_installations)
        
        # 获取所有候选适配器
        candidates = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED,
            MarketplaceAdapter.id.notin_([a.id for a in user_installations])
        ).all()
        
        # 计算相似度
        similarities = []
        for candidate in candidates:
            candidate_features = self._build_adapter_feature_vector(candidate)
            similarity = cosine_similarity(
                user_features.reshape(1, -1),
                candidate_features.reshape(1, -1)
            )[0][0]
            
            similarities.append({
                "adapter": candidate,
                "score": similarity,
                "reason": "content_similarity"
            })
        
        # 按相似度排序
        similarities.sort(key=lambda x: x['score'], reverse=True)
        
        return similarities[:count]
    
    async def _popular_recommendations(self, count: int) -> List[Dict]:
        """
        热门推荐
        
        基于：活跃安装量、评分、下载量的加权组合
        """
        adapters = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
        ).all()
        
        # 计算热度分数
        scored_adapters = []
        for adapter in adapters:
            popularity_score = (
                adapter.active_install_count * 0.5 +
                adapter.rating_average * adapter.rating_count * 0.3 +
                adapter.download_count * 0.0001 * 0.2
            )
            
            scored_adapters.append({
                "adapter": adapter,
                "score": popularity_score,
                "reason": "popular"
            })
        
        # 排序
        scored_adapters.sort(key=lambda x: x['score'], reverse=True)
        
        return scored_adapters[:count]
    
    async def _context_based_recommendations(
        self,
        context: Dict[str, Any],
        count: int
    ) -> List[Dict]:
        """
        基于上下文的推荐
        
        Args:
            context: {
                "domain": "chemistry",
                "activity_type": "concept_understanding",
                "subject": "periodic_table"
            }
        """
        filters = []
        
        if "domain" in context:
            filters.append(
                MarketplaceAdapter.supported_domains.contains([context["domain"]])
            )
        
        if "activity_type" in context:
            filters.append(
                MarketplaceAdapter.supported_activity_types.contains([context["activity_type"]])
            )
        
        adapters = await db.query(MarketplaceAdapter).filter(
            MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED,
            *filters
        ).order_by(
            desc(MarketplaceAdapter.rating_average)
        ).limit(count).all()
        
        return [
            {
                "adapter": adapter,
                "score": adapter.rating_average,
                "reason": "context_match"
            }
            for adapter in adapters
        ]
    
    def _build_user_feature_vector(
        self,
        installations: List[MarketplaceAdapter]
    ) -> np.ndarray:
        """构建用户特征向量（基于已安装适配器的平均）"""
        feature_vectors = [
            self._build_adapter_feature_vector(adapter)
            for adapter in installations
        ]
        
        return np.mean(feature_vectors, axis=0)
    
    def _build_adapter_feature_vector(
        self,
        adapter: MarketplaceAdapter
    ) -> np.ndarray:
        """
        构建适配器特征向量
        
        特征包括：
        - 分类 (one-hot)
        - 标签 (multi-hot)
        - 领域 (multi-hot)
        - 活动类型 (multi-hot)
        """
        # 这里简化示意，实际应该用更完整的特征工程
        features = []
        
        # 分类 (假设有10个分类)
        category_features = [0] * 10
        category_idx = hash(adapter.category) % 10
        category_features[category_idx] = 1
        features.extend(category_features)
        
        # 标签 (假设有50个常见标签)
        tag_features = [0] * 50
        if adapter.tags:
            for tag in adapter.tags:
                tag_idx = hash(tag) % 50
                tag_features[tag_idx] = 1
        features.extend(tag_features)
        
        # 质量分数
        features.append(adapter.quality_score or 0.0)
        features.append(adapter.rating_average)
        
        return np.array(features)
    
    def _deduplicate_and_rank(
        self,
        recommendations: List[Dict],
        count: int
    ) -> List[Dict]:
        """去重并重新排序"""
        seen = set()
        unique_recs = []
        
        for rec in recommendations:
            adapter_id = rec['adapter'].adapter_id
            if adapter_id not in seen:
                seen.add(adapter_id)
                unique_recs.append(rec)
        
        # 按分数排序
        unique_recs.sort(key=lambda x: x['score'], reverse=True)
        
        return unique_recs[:count]


# 全局推荐引擎实例
recommendation_engine = AdapterRecommendationEngine()
```

### 4.2.2 推荐API

```python
# src/marketplace/api/recommendations.py

from fastapi import APIRouter, Depends

router = APIRouter(prefix="/api/v1/marketplace/recommendations", tags=["recommendations"])


@router.get("/for-you")
async def get_personalized_recommendations(
    count: int = Query(10, ge=1, le=50),
    current_user: User = Depends(get_current_user)
):
    """
    获取个性化推荐
    
    基于用户的安装历史、评分、使用习惯等
    """
    from src.marketplace.recommendations import recommendation_engine
    
    recommendations = await recommendation_engine.get_recommendations(
        user_id=current_user.id,
        count=count
    )
    
    return {
        "recommendations": [
            {
                "adapter_id": rec['adapter'].adapter_id,
                "name": rec['adapter'].display_name,
                "description": rec['adapter'].tagline,
                "rating": rec['adapter'].rating_average,
                "downloads": rec['adapter'].download_count,
                "icon_url": rec['adapter'].icon_url,
                "reason": rec['reason'],
                "score": rec['score']
            }
            for rec in recommendations
        ]
    }


@router.get("/similar/{adapter_id}")
async def get_similar_adapters(
    adapter_id: str,
    count: int = Query(5, ge=1, le=20)
):
    """
    获取相似的适配器
    
    用于"您可能还喜欢"功能
    """
    # 获取目标适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    from src.marketplace.recommendations import recommendation_engine
    
    # 基于内容相似度推荐
    similar = await recommendation_engine._content_based_filtering(
        user_installations=[adapter],
        count=count
    )
    
    return {
        "similar_adapters": [
            {
                "adapter_id": rec['adapter'].adapter_id,
                "name": rec['adapter'].display_name,
                "description": rec['adapter'].tagline,
                "rating": rec['adapter'].rating_average,
                "similarity": rec['score']
            }
            for rec in similar
        ]
    }


@router.get("/trending")
async def get_trending_adapters(
    period: str = Query("week", regex="^(day|week|month)$"),
    count: int = Query(10, ge=1, le=50)
):
    """
    获取趋势适配器
    
    基于最近的下载量增长、评分增长等
    """
    # 计算时间范围
    from datetime import datetime, timedelta
    
    period_map = {
        "day": timedelta(days=1),
        "week": timedelta(weeks=1),
        "month": timedelta(days=30)
    }
    
    since = datetime.utcnow() - period_map[period]
    
    # 统计最近的安装量
    trending = await db.query(
        MarketplaceAdapter,
        func.count(AdapterInstallation.id).label('recent_installs')
    ).join(
        AdapterInstallation
    ).filter(
        AdapterInstallation.installed_at >= since,
        MarketplaceAdapter.status == AdapterPublishStatus.PUBLISHED
    ).group_by(
        MarketplaceAdapter.id
    ).order_by(
        desc('recent_installs')
    ).limit(count).all()
    
    return {
        "period": period,
        "trending_adapters": [
            {
                "adapter_id": adapter.adapter_id,
                "name": adapter.display_name,
                "description": adapter.tagline,
                "rating": adapter.rating_average,
                "recent_installs": recent_installs,
                "total_installs": adapter.install_count
            }
            for adapter, recent_installs in trending
        ]
    }
```

---

## 4.3 评分评论系统

### 4.3.1 评论管理

```python
# src/marketplace/api/reviews.py

from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel, Field

router = APIRouter(prefix="/api/v1/marketplace/reviews", tags=["reviews"])


class ReviewCreateRequest(BaseModel):
    """创建评论请求"""
    rating: int = Field(..., ge=1, le=5, description="评分 1-5星")
    title: Optional[str] = Field(None, max_length=200)
    content: Optional[str] = Field(None, max_length=5000)
    use_case: Optional[str] = Field(None, max_length=100)


@router.post("/adapters/{adapter_id}/reviews")
async def create_review(
    adapter_id: str,
    request: ReviewCreateRequest,
    current_user: User = Depends(get_current_user)
):
    """
    创建评论
    
    用户必须安装过该适配器才能评论
    """
    # 查找适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 检查是否已安装
    installation = await db.query(AdapterInstallation).filter_by(
        user_id=current_user.id,
        adapter_id=adapter.id
    ).first()
    
    if not installation:
        raise HTTPException(
            status_code=403,
            detail="You must install the adapter before reviewing"
        )
    
    # 检查是否已评论
    existing_review = await db.query(AdapterReview).filter_by(
        user_id=current_user.id,
        adapter_id=adapter.id
    ).first()
    
    if existing_review:
        raise HTTPException(
            status_code=400,
            detail="You have already reviewed this adapter"
        )
    
    # 创建评论
    review = AdapterReview(
        adapter_id=adapter.id,
        user_id=current_user.id,
        rating=request.rating,
        title=request.title,
        content=request.content,
        use_case=request.use_case,
        is_verified_user=current_user.is_verified
    )
    
    db.add(review)
    
    # 更新适配器评分统计
    await update_adapter_rating(adapter.id)
    
    await db.commit()
    
    return {"message": "Review created successfully"}


@router.get("/adapters/{adapter_id}/reviews")
async def get_adapter_reviews(
    adapter_id: str,
    sort_by: str = Query("newest", regex="^(newest|helpful|rating_high|rating_low)$"),
    min_rating: Optional[int] = Query(None, ge=1, le=5),
    page: int = Query(1, ge=1),
    page_size: int = Query(20, ge=1, le=100)
):
    """获取适配器的所有评论"""
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 构建查询
    query = db.query(AdapterReview).filter(
        AdapterReview.adapter_id == adapter.id,
        AdapterReview.is_visible == True
    )
    
    # 评分过滤
    if min_rating:
        query = query.filter(AdapterReview.rating >= min_rating)
    
    # 排序
    if sort_by == "newest":
        query = query.order_by(desc(AdapterReview.created_at))
    elif sort_by == "helpful":
        query = query.order_by(desc(AdapterReview.helpful_count))
    elif sort_by == "rating_high":
        query = query.order_by(desc(AdapterReview.rating))
    elif sort_by == "rating_low":
        query = query.order_by(asc(AdapterReview.rating))
    
    # 分页
    total = await query.count()
    reviews = await query.offset((page - 1) * page_size).limit(page_size).all()
    
    return {
        "total": total,
        "page": page,
        "page_size": page_size,
        "reviews": [
            {
                "id": review.id,
                "rating": review.rating,
                "title": review.title,
                "content": review.content,
                "author": {
                    "username": review.user.username,
                    "display_name": review.user.display_name,
                    "is_verified": review.user.is_verified
                },
                "helpful_count": review.helpful_count,
                "created_at": review.created_at.isoformat()
            }
            for review in reviews
        ]
    }
```

---

*文档第五部分完成*

**已完成内容**：
- ✅ Elasticsearch搜索系统（索引设计、多字段搜索、模糊匹配、自动补全）
- ✅ 搜索API（关键词、过滤、排序、分页、过滤器选项）
- ✅ 推荐系统（协同过滤、内容推荐、热门推荐、上下文推荐）
- ✅ 推荐API（个性化推荐、相似推荐、趋势推荐）
- ✅ 评分评论系统（创建评论、查询评论、评分统计）

**下一部分将包含**：
- 统计分析系统
- 开发者控制台
- SDK和工具
- 总结和最佳实践

是否继续输出第六部分（最终部分）？
