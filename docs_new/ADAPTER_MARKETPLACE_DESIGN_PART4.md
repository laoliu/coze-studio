# MetaWorkflow V2.0 适配器市场设计 - 第四部分

## 第三部分：适配器市场平台

## 3.1 市场架构

### 3.1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户层                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  开发者控制台 │  │   用户商店   │  │  管理后台    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API网关层                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 认证授权  │  │ 限流控制  │  │ 日志审计  │  │ API版本  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      业务服务层                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  适配器管理  │  │   发布审核   │  │   搜索推荐   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  用户管理    │  │   评分评论   │  │   统计分析   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      数据存储层                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  PostgreSQL │  │   Redis     │  │   S3存储    │             │
│  │  (元数据)    │  │  (缓存)     │  │  (文件)     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐                               │
│  │ Elasticsearch│  │   ClickHouse│                               │
│  │  (全文搜索)  │  │  (日志分析)  │                               │
│  └─────────────┘  └─────────────┘                               │
└─────────────────────────────────────────────────────────────────┘
```

### 3.1.2 数据模型设计

```python
# src/marketplace/models.py

from sqlalchemy import Column, Integer, String, Text, DateTime, Boolean, Float, JSON, ForeignKey, Enum
from sqlalchemy.orm import relationship
from datetime import datetime
import enum

from src.db.database import Base


class AdapterPublishStatus(enum.Enum):
    """发布状态"""
    DRAFT = "draft"              # 草稿
    SUBMITTED = "submitted"      # 已提交
    REVIEWING = "reviewing"      # 审核中
    APPROVED = "approved"        # 已批准
    REJECTED = "rejected"        # 已拒绝
    PUBLISHED = "published"      # 已发布
    DEPRECATED = "deprecated"    # 已废弃


class MarketplaceAdapter(Base):
    """市场适配器"""
    __tablename__ = "marketplace_adapters"
    
    # 基本信息
    id = Column(Integer, primary_key=True)
    adapter_id = Column(String(100), unique=True, nullable=False)  # 唯一标识
    name = Column(String(100), nullable=False)
    display_name = Column(String(200), nullable=False)
    description = Column(Text, nullable=False)
    tagline = Column(String(200))  # 一句话描述
    
    # 作者信息
    author_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    author = relationship("User", back_populates="adapters")
    organization = Column(String(200))
    
    # 版本信息
    current_version = Column(String(20), nullable=False)
    latest_version_id = Column(Integer, ForeignKey("adapter_versions.id"))
    
    # 分类和标签
    category = Column(String(50), nullable=False)  # 主分类：k12/higher_ed/enterprise
    subcategory = Column(String(50))  # 子分类：chemistry/physics/math
    tags = Column(JSON)  # 标签列表
    
    # 能力声明
    supported_domains = Column(JSON)  # 支持的领域
    supported_activity_types = Column(JSON)  # 支持的活动类型
    supported_output_formats = Column(JSON)  # 支持的输出格式
    
    # 依赖
    platform_version_min = Column(String(20))  # 最低平台版本
    platform_version_max = Column(String(20))  # 最高平台版本
    dependencies = Column(JSON)  # 依赖的其他适配器
    
    # 发布状态
    status = Column(Enum(AdapterPublishStatus), default=AdapterPublishStatus.DRAFT)
    published_at = Column(DateTime)
    deprecated_at = Column(DateTime)
    
    # 统计数据
    download_count = Column(Integer, default=0)
    install_count = Column(Integer, default=0)
    active_install_count = Column(Integer, default=0)  # 活跃安装数
    rating_average = Column(Float, default=0.0)
    rating_count = Column(Integer, default=0)
    
    # 质量指标
    quality_score = Column(Float)  # 综合质量分
    test_coverage = Column(Float)  # 测试覆盖率
    security_scan_passed = Column(Boolean, default=False)
    performance_score = Column(Float)  # 性能评分
    
    # 资源链接
    homepage_url = Column(String(500))
    documentation_url = Column(String(500))
    repository_url = Column(String(500))
    support_url = Column(String(500))
    
    # 媒体资源
    icon_url = Column(String(500))
    screenshots = Column(JSON)  # 截图URLs
    demo_video_url = Column(String(500))
    
    # 许可证
    license = Column(String(50))  # MIT/Apache-2.0/GPL-3.0等
    is_open_source = Column(Boolean, default=True)
    
    # 定价（为未来商业化预留）
    is_free = Column(Boolean, default=True)
    pricing_model = Column(String(50))  # free/freemium/paid
    price = Column(Float)  # 月费
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # 关系
    versions = relationship("AdapterVersion", back_populates="adapter", cascade="all, delete-orphan")
    reviews = relationship("AdapterReview", back_populates="adapter", cascade="all, delete-orphan")


class AdapterVersion(Base):
    """适配器版本"""
    __tablename__ = "adapter_versions"
    
    id = Column(Integer, primary_key=True)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    
    # 版本信息
    version = Column(String(20), nullable=False)
    version_name = Column(String(100))  # 版本代号（可选）
    
    # 变更信息
    changelog = Column(Text)
    release_notes = Column(Text)
    
    # 文件信息
    package_url = Column(String(500), nullable=False)  # 包文件URL
    package_size = Column(Integer)  # 字节
    package_checksum = Column(String(64))  # SHA-256
    
    # manifest.yaml内容
    manifest = Column(JSON, nullable=False)
    
    # 兼容性
    platform_version_min = Column(String(20))
    platform_version_max = Column(String(20))
    breaking_changes = Column(Boolean, default=False)
    
    # 审核信息
    review_status = Column(Enum(AdapterPublishStatus), default=AdapterPublishStatus.DRAFT)
    reviewed_by = Column(Integer, ForeignKey("users.id"))
    reviewed_at = Column(DateTime)
    review_comments = Column(Text)
    
    # 质量检测
    test_results = Column(JSON)  # 测试结果
    security_scan_results = Column(JSON)  # 安全扫描结果
    performance_benchmarks = Column(JSON)  # 性能基准测试
    
    # 统计
    download_count = Column(Integer, default=0)
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    published_at = Column(DateTime)
    
    # 关系
    adapter = relationship("MarketplaceAdapter", back_populates="versions")


class AdapterReview(Base):
    """适配器评论"""
    __tablename__ = "adapter_reviews"
    
    id = Column(Integer, primary_key=True)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    
    # 评分（1-5星）
    rating = Column(Integer, nullable=False)
    
    # 评论内容
    title = Column(String(200))
    content = Column(Text)
    
    # 使用场景
    use_case = Column(String(100))  # 使用场景分类
    adapter_version = Column(String(20))  # 评论时使用的版本
    
    # 点赞统计
    helpful_count = Column(Integer, default=0)
    not_helpful_count = Column(Integer, default=0)
    
    # 审核
    is_verified_user = Column(Boolean, default=False)  # 是否认证用户
    is_moderated = Column(Boolean, default=False)  # 是否经过审核
    is_visible = Column(Boolean, default=True)
    
    # 作者回复
    author_reply = Column(Text)
    author_replied_at = Column(DateTime)
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    # 关系
    adapter = relationship("MarketplaceAdapter", back_populates="reviews")
    user = relationship("User", back_populates="reviews")


class AdapterInstallation(Base):
    """适配器安装记录"""
    __tablename__ = "adapter_installations"
    
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    adapter_id = Column(Integer, ForeignKey("marketplace_adapters.id"), nullable=False)
    version_id = Column(Integer, ForeignKey("adapter_versions.id"), nullable=False)
    
    # 安装状态
    status = Column(String(20))  # installed/active/inactive/uninstalled
    
    # 配置
    configuration = Column(JSON)  # 用户自定义配置
    
    # 使用统计
    last_used_at = Column(DateTime)
    usage_count = Column(Integer, default=0)
    
    # 时间戳
    installed_at = Column(DateTime, default=datetime.utcnow)
    uninstalled_at = Column(DateTime)


class User(Base):
    """用户"""
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True)
    username = Column(String(50), unique=True, nullable=False)
    email = Column(String(200), unique=True, nullable=False)
    
    # 身份
    is_developer = Column(Boolean, default=False)
    is_verified = Column(Boolean, default=False)
    is_admin = Column(Boolean, default=False)
    
    # 个人信息
    display_name = Column(String(100))
    avatar_url = Column(String(500))
    bio = Column(Text)
    organization = Column(String(200))
    website = Column(String(500))
    
    # 时间戳
    created_at = Column(DateTime, default=datetime.utcnow)
    
    # 关系
    adapters = relationship("MarketplaceAdapter", back_populates="author")
    reviews = relationship("AdapterReview", back_populates="user")
```

---

## 3.2 发布流程

### 3.2.1 开发者发布流程

```
┌──────────────┐
│ 1. 创建适配器 │
│   开发本地代码  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 2. 本地测试   │
│  运行单元测试  │
│  质量检查     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 3. 打包       │
│  生成manifest │
│  创建.zip包   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 4. 提交审核   │
│  上传到市场   │
│  填写描述信息  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ 5. 自动检测   │
│  安全扫描     │
│  性能测试     │
│  兼容性测试   │
└──────┬───────┘
       │
       ├─→ 失败 ──→ 通知开发者 ──→ 返回步骤1
       │
       ▼ 通过
┌──────────────┐
│ 6. 人工审核   │
│  代码审查     │
│  功能验证     │
│  文档检查     │
└──────┬───────┘
       │
       ├─→ 拒绝 ──→ 通知原因 ──→ 返回步骤1
       │
       ▼ 批准
┌──────────────┐
│ 7. 发布上架   │
│  设置为可见   │
│  发送通知     │
└──────────────┘
```

### 3.2.2 发布API设计

```python
# src/marketplace/api/publish.py

from fastapi import APIRouter, HTTPException, Depends, UploadFile, File
from typing import Optional
from pydantic import BaseModel

router = APIRouter(prefix="/api/v1/marketplace/publish", tags=["publish"])


class AdapterSubmitRequest(BaseModel):
    """适配器提交请求"""
    name: str
    display_name: str
    description: str
    tagline: Optional[str]
    category: str
    subcategory: Optional[str]
    tags: List[str] = []
    homepage_url: Optional[str]
    repository_url: Optional[str]
    documentation_url: Optional[str]
    license: str = "MIT"


class VersionSubmitRequest(BaseModel):
    """版本提交请求"""
    version: str
    changelog: str
    release_notes: Optional[str]
    breaking_changes: bool = False


@router.post("/adapters")
async def create_adapter(
    request: AdapterSubmitRequest,
    current_user: User = Depends(get_current_user)
):
    """
    创建新适配器（草稿）
    
    这个API创建一个草稿状态的适配器记录，
    开发者需要上传版本包后才能提交审核。
    """
    # 验证开发者权限
    if not current_user.is_developer:
        raise HTTPException(
            status_code=403,
            detail="Only verified developers can publish adapters"
        )
    
    # 检查名称是否已存在
    existing = await db.query(MarketplaceAdapter).filter_by(
        name=request.name
    ).first()
    
    if existing:
        raise HTTPException(
            status_code=400,
            detail=f"Adapter name '{request.name}' already exists"
        )
    
    # 创建适配器记录
    adapter = MarketplaceAdapter(
        adapter_id=generate_adapter_id(request.name),
        name=request.name,
        display_name=request.display_name,
        description=request.description,
        tagline=request.tagline,
        category=request.category,
        subcategory=request.subcategory,
        tags=request.tags,
        author_id=current_user.id,
        homepage_url=request.homepage_url,
        repository_url=request.repository_url,
        documentation_url=request.documentation_url,
        license=request.license,
        status=AdapterPublishStatus.DRAFT
    )
    
    db.add(adapter)
    await db.commit()
    
    return {
        "adapter_id": adapter.adapter_id,
        "status": "draft",
        "message": "Adapter created. Please upload a version package."
    }


@router.post("/adapters/{adapter_id}/versions")
async def upload_version(
    adapter_id: str,
    version_info: VersionSubmitRequest,
    package_file: UploadFile = File(...),
    current_user: User = Depends(get_current_user)
):
    """
    上传适配器版本包
    
    步骤：
    1. 验证权限
    2. 验证包格式
    3. 提取manifest.yaml
    4. 上传到对象存储
    5. 创建版本记录
    """
    # 查找适配器
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    # 验证所有权
    if adapter.author_id != current_user.id:
        raise HTTPException(
            status_code=403,
            detail="You don't have permission to modify this adapter"
        )
    
    # 验证包格式
    if not package_file.filename.endswith('.zip'):
        raise HTTPException(
            status_code=400,
            detail="Package must be a .zip file"
        )
    
    # 读取包内容
    package_content = await package_file.read()
    
    # 验证和提取manifest
    try:
        manifest = await extract_manifest(package_content)
    except Exception as e:
        raise HTTPException(
            status_code=400,
            detail=f"Invalid package: {str(e)}"
        )
    
    # 验证版本号匹配
    if manifest['version'] != version_info.version:
        raise HTTPException(
            status_code=400,
            detail="Version mismatch between request and manifest"
        )
    
    # 检查版本是否已存在
    existing_version = await db.query(AdapterVersion).filter_by(
        adapter_id=adapter.id,
        version=version_info.version
    ).first()
    
    if existing_version:
        raise HTTPException(
            status_code=400,
            detail=f"Version {version_info.version} already exists"
        )
    
    # 上传到对象存储
    package_url = await upload_to_storage(
        package_content,
        f"adapters/{adapter_id}/{version_info.version}/package.zip"
    )
    
    # 计算校验和
    checksum = calculate_sha256(package_content)
    
    # 创建版本记录
    version = AdapterVersion(
        adapter_id=adapter.id,
        version=version_info.version,
        changelog=version_info.changelog,
        release_notes=version_info.release_notes,
        breaking_changes=version_info.breaking_changes,
        package_url=package_url,
        package_size=len(package_content),
        package_checksum=checksum,
        manifest=manifest,
        platform_version_min=manifest.get('compatibility', {}).get('platform_version'),
        review_status=AdapterPublishStatus.DRAFT
    )
    
    db.add(version)
    
    # 更新适配器的当前版本
    if not adapter.current_version or version_info.version > adapter.current_version:
        adapter.current_version = version_info.version
        adapter.latest_version_id = version.id
    
    await db.commit()
    
    return {
        "version_id": version.id,
        "version": version.version,
        "package_url": package_url,
        "status": "draft",
        "message": "Version uploaded successfully. You can now submit for review."
    }


@router.post("/adapters/{adapter_id}/submit")
async def submit_for_review(
    adapter_id: str,
    current_user: User = Depends(get_current_user)
):
    """
    提交适配器审核
    
    将适配器状态从DRAFT改为SUBMITTED，
    并触发自动检测流程。
    """
    adapter = await db.query(MarketplaceAdapter).filter_by(
        adapter_id=adapter_id
    ).first()
    
    if not adapter:
        raise HTTPException(status_code=404, detail="Adapter not found")
    
    if adapter.author_id != current_user.id:
        raise HTTPException(status_code=403, detail="Permission denied")
    
    # 检查是否有版本
    if not adapter.versions:
        raise HTTPException(
            status_code=400,
            detail="Cannot submit adapter without any versions"
        )
    
    # 检查当前状态
    if adapter.status not in [AdapterPublishStatus.DRAFT, AdapterPublishStatus.REJECTED]:
        raise HTTPException(
            status_code=400,
            detail=f"Cannot submit adapter with status: {adapter.status.value}"
        )
    
    # 更新状态
    adapter.status = AdapterPublishStatus.SUBMITTED
    
    # 获取最新版本
    latest_version = adapter.versions[-1]
    latest_version.review_status = AdapterPublishStatus.SUBMITTED
    
    await db.commit()
    
    # 触发自动检测（异步任务）
    from src.marketplace.tasks import run_automated_checks
    await run_automated_checks.delay(adapter.id, latest_version.id)
    
    # 通知审核团队
    await notify_review_team(adapter.id)
    
    return {
        "adapter_id": adapter_id,
        "status": "submitted",
        "message": "Adapter submitted for review. Automated checks will run shortly."
    }
```

---

## 3.3 自动化检测

### 3.3.1 检测流程

```python
# src/marketplace/tasks.py

from celery import Celery
import asyncio

celery_app = Celery('marketplace', broker='redis://localhost:6379/0')


@celery_app.task
async def run_automated_checks(adapter_id: int, version_id: int):
    """
    运行自动化检测
    
    检测项：
    1. 安全扫描
    2. 性能测试
    3. 兼容性测试
    4. 代码质量检测
    """
    version = await db.query(AdapterVersion).get(version_id)
    
    results = {
        "security": None,
        "performance": None,
        "compatibility": None,
        "quality": None,
        "overall_passed": False
    }
    
    try:
        # 1. 安全扫描
        results["security"] = await run_security_scan(version)
        
        # 2. 性能测试
        results["performance"] = await run_performance_tests(version)
        
        # 3. 兼容性测试
        results["compatibility"] = await run_compatibility_tests(version)
        
        # 4. 代码质量
        results["quality"] = await run_quality_checks(version)
        
        # 判断是否全部通过
        results["overall_passed"] = all([
            results["security"]["passed"],
            results["performance"]["passed"],
            results["compatibility"]["passed"],
            results["quality"]["passed"]
        ])
        
        # 保存结果
        version.test_results = results["quality"]
        version.security_scan_results = results["security"]
        version.performance_benchmarks = results["performance"]
        
        # 更新状态
        if results["overall_passed"]:
            version.review_status = AdapterPublishStatus.REVIEWING
            # 通知人工审核
            await notify_human_reviewers(adapter_id, version_id)
        else:
            version.review_status = AdapterPublishStatus.REJECTED
            # 通知开发者
            await notify_developer_rejection(adapter_id, version_id, results)
        
        await db.commit()
        
    except Exception as e:
        logger.error(f"Automated checks failed: {e}")
        version.review_status = AdapterPublishStatus.REJECTED
        await db.commit()
    
    return results


async def run_security_scan(version: AdapterVersion) -> Dict:
    """
    安全扫描
    
    检测：
    - 恶意代码
    - 已知漏洞
    - 敏感信息泄露
    - 不安全的依赖
    """
    # 下载包
    package_path = await download_package(version.package_url)
    
    results = {
        "passed": True,
        "issues": [],
        "severity_counts": {"critical": 0, "high": 0, "medium": 0, "low": 0}
    }
    
    # 1. 病毒扫描
    virus_scan = await run_clamav_scan(package_path)
    if not virus_scan["clean"]:
        results["passed"] = False
        results["issues"].append({
            "type": "malware",
            "severity": "critical",
            "description": "Malware detected"
        })
    
    # 2. 依赖漏洞扫描（使用safety或snyk）
    dependency_scan = await run_safety_check(package_path)
    for vuln in dependency_scan["vulnerabilities"]:
        results["issues"].append({
            "type": "vulnerability",
            "severity": vuln["severity"],
            "package": vuln["package"],
            "description": vuln["description"],
            "cve": vuln.get("cve")
        })
        results["severity_counts"][vuln["severity"]] += 1
    
    # 3. 代码静态分析（使用bandit）
    code_scan = await run_bandit_scan(package_path)
    for issue in code_scan["issues"]:
        if issue["severity"] in ["high", "critical"]:
            results["passed"] = False
        results["issues"].append({
            "type": "code_issue",
            "severity": issue["severity"],
            "description": issue["description"],
            "file": issue["file"],
            "line": issue["line"]
        })
    
    # 4. 检查敏感信息
    secrets_scan = await run_trufflehog_scan(package_path)
    if secrets_scan["found"]:
        results["passed"] = False
        results["issues"].append({
            "type": "secrets",
            "severity": "critical",
            "description": "Hardcoded secrets detected"
        })
    
    # Critical或High级别的问题会导致不通过
    if results["severity_counts"]["critical"] > 0 or results["severity_counts"]["high"] > 0:
        results["passed"] = False
    
    return results


async def run_performance_tests(version: AdapterVersion) -> Dict:
    """
    性能测试
    
    测试：
    - 内存占用
    - CPU使用率
    - 响应时间
    - 吞吐量
    """
    results = {
        "passed": True,
        "metrics": {},
        "benchmarks": []
    }
    
    # 加载适配器
    adapter = await load_adapter_for_testing(version)
    
    # 1. 内存测试
    memory_usage = await measure_memory_usage(adapter)
    results["metrics"]["memory_mb"] = memory_usage
    
    # 内存超过500MB视为不通过
    if memory_usage > 500:
        results["passed"] = False
    
    # 2. 响应时间测试
    response_times = []
    for i in range(10):
        start = time.time()
        await adapter.parse_request({
            "activity_type": "concept_understanding",
            "subject": "chemistry",
            "topic": "periodic_table"
        })
        response_times.append(time.time() - start)
    
    avg_response_time = sum(response_times) / len(response_times)
    results["metrics"]["avg_response_time_ms"] = avg_response_time * 1000
    
    # 平均响应时间超过5秒视为不通过
    if avg_response_time > 5.0:
        results["passed"] = False
    
    # 3. 并发测试
    concurrent_results = await run_concurrent_test(adapter, concurrency=10)
    results["metrics"]["concurrent_throughput"] = concurrent_results["throughput"]
    
    return results
```

---

*文档第四部分完成*

**已完成内容**：
- ✅ 市场整体架构（4层架构）
- ✅ 完整的数据模型（适配器、版本、评论、安装、用户）
- ✅ 发布流程设计（7步审核流程）
- ✅ 发布API实现（创建、上传、提交）
- ✅ 自动化检测（安全扫描、性能测试）

**下一部分将包含**：
- 搜索和推荐系统
- 评分评论机制
- 使用统计分析
- 开发者工具和SDK

是否继续输出第五部分（用户体验和分析系统）？