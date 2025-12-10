# MetaWorkflow V2.0 适配器市场设计 - 第三部分

## 第二部分：适配器注册机制

## 2.1 注册表设计

### 2.1.1 注册表架构

```python
# src/adapters/registry.py

from typing import Dict, List, Optional, Type
from datetime import datetime
import importlib
import yaml
from pathlib import Path

from .base import DomainAdapter, AdapterMetadata
from src.models import DomainType

class AdapterRegistry:
    """
    适配器注册表
    
    职责：
    1. 管理所有已注册的适配器
    2. 提供适配器查找和获取功能
    3. 处理适配器的生命周期
    4. 验证适配器兼容性
    """
    
    def __init__(self):
        """初始化注册表"""
        self._adapters: Dict[DomainType, DomainAdapter] = {}
        self._metadata: Dict[DomainType, AdapterMetadata] = {}
        self._instances: Dict[str, DomainAdapter] = {}  # adapter_id -> instance
        self._load_history: List[Dict] = []
    
    def register(
        self, 
        domain: DomainType, 
        adapter: DomainAdapter,
        force: bool = False
    ) -> None:
        """
        注册适配器
        
        Args:
            domain: 领域类型
            adapter: 适配器实例
            force: 是否强制覆盖已存在的适配器
            
        Raises:
            ValueError: 适配器已存在且force=False
            RuntimeError: 适配器验证失败
        """
        # 检查是否已注册
        if domain in self._adapters and not force:
            raise ValueError(
                f"Domain {domain.value} already has a registered adapter. "
                f"Use force=True to override."
            )
        
        # 验证兼容性
        from src import __version__ as platform_version
        is_compatible, error = adapter.validate_compatibility(platform_version)
        if not is_compatible:
            raise RuntimeError(f"Adapter incompatible: {error}")
        
        # 获取元数据
        metadata = adapter.metadata
        
        # 验证元数据
        self._validate_metadata(metadata)
        
        # 注册
        self._adapters[domain] = adapter
        self._metadata[domain] = metadata
        self._instances[metadata.name] = adapter
        
        # 记录加载历史
        self._load_history.append({
            "domain": domain.value,
            "adapter": metadata.name,
            "version": metadata.version,
            "loaded_at": datetime.utcnow().isoformat(),
            "action": "register"
        })
        
        logger.info(
            f"Registered adapter: {metadata.display_name} v{metadata.version} "
            f"for domain {domain.value}"
        )
    
    def unregister(self, domain: DomainType) -> None:
        """
        注销适配器
        
        Args:
            domain: 领域类型
        """
        if domain not in self._adapters:
            return
        
        adapter = self._adapters[domain]
        metadata = self._metadata[domain]
        
        # 调用卸载钩子
        import asyncio
        asyncio.create_task(adapter.unload())
        
        # 从注册表移除
        del self._adapters[domain]
        del self._metadata[domain]
        if metadata.name in self._instances:
            del self._instances[metadata.name]
        
        # 记录历史
        self._load_history.append({
            "domain": domain.value,
            "adapter": metadata.name,
            "version": metadata.version,
            "loaded_at": datetime.utcnow().isoformat(),
            "action": "unregister"
        })
        
        logger.info(f"Unregistered adapter for domain {domain.value}")
    
    def get_adapter(self, domain: DomainType) -> Optional[DomainAdapter]:
        """
        获取适配器
        
        Args:
            domain: 领域类型
            
        Returns:
            Optional[DomainAdapter]: 适配器实例，不存在则返回None
        """
        return self._adapters.get(domain)
    
    def get_metadata(self, domain: DomainType) -> Optional[AdapterMetadata]:
        """
        获取适配器元数据
        
        Args:
            domain: 领域类型
            
        Returns:
            Optional[AdapterMetadata]: 元数据，不存在则返回None
        """
        return self._metadata.get(domain)
    
    def list_adapters(self) -> List[Dict[str, Any]]:
        """
        列出所有已注册的适配器
        
        Returns:
            List[Dict[str, Any]]: 适配器信息列表
        """
        adapters = []
        for domain, metadata in self._metadata.items():
            adapters.append({
                "domain": domain.value,
                "name": metadata.name,
                "display_name": metadata.display_name,
                "version": metadata.version,
                "author": metadata.author,
                "description": metadata.description,
                "capabilities": metadata.capabilities,
                "status": "active"
            })
        return adapters
    
    def _validate_metadata(self, metadata: AdapterMetadata) -> None:
        """
        验证适配器元数据
        
        Args:
            metadata: 适配器元数据
            
        Raises:
            ValueError: 元数据不合法
        """
        # 必填字段检查
        required_fields = ['name', 'display_name', 'version', 'description']
        for field in required_fields:
            if not getattr(metadata, field, None):
                raise ValueError(f"Missing required metadata field: {field}")
        
        # 版本格式检查（语义化版本）
        import re
        version_pattern = r'^\d+\.\d+\.\d+$'
        if not re.match(version_pattern, metadata.version):
            raise ValueError(
                f"Invalid version format: {metadata.version}. "
                f"Expected semantic version (e.g., 1.0.0)"
            )
    
    async def initialize_all(self) -> None:
        """初始化所有已注册的适配器"""
        for adapter in self._adapters.values():
            if not adapter._initialized:
                await adapter.initialize()
    
    async def health_check_all(self) -> Dict[str, Any]:
        """
        检查所有适配器的健康状态
        
        Returns:
            Dict[str, Any]: 健康状态汇总
        """
        results = {}
        for domain, adapter in self._adapters.items():
            try:
                health = await adapter.health_check()
                results[domain.value] = {
                    "status": health.get("status", "unknown"),
                    "details": health
                }
            except Exception as e:
                results[domain.value] = {
                    "status": "error",
                    "error": str(e)
                }
        
        return {
            "overall_status": self._calculate_overall_status(results),
            "adapters": results,
            "checked_at": datetime.utcnow().isoformat()
        }
    
    def _calculate_overall_status(self, results: Dict) -> str:
        """计算整体健康状态"""
        statuses = [r.get("status") for r in results.values()]
        
        if all(s == "healthy" for s in statuses):
            return "healthy"
        elif any(s == "unhealthy" for s in statuses):
            return "unhealthy"
        else:
            return "degraded"


# 全局注册表实例
adapter_registry = AdapterRegistry()
```

### 2.1.2 自动发现机制

```python
# src/adapters/discovery.py

from pathlib import Path
from typing import List, Dict, Any
import importlib.util
import sys
import yaml

class AdapterDiscovery:
    """
    适配器自动发现
    
    支持多种发现方式：
    1. 从指定目录扫描
    2. 从Python包导入
    3. 从适配器市场下载
    """
    
    def __init__(self, search_paths: List[str] = None):
        """
        初始化发现器
        
        Args:
            search_paths: 适配器搜索路径列表
        """
        self.search_paths = search_paths or [
            "src/adapters",
            "~/.metaworkflow/adapters",
            "/opt/metaworkflow/adapters"
        ]
    
    def discover_all(self) -> List[Dict[str, Any]]:
        """
        发现所有可用的适配器
        
        Returns:
            List[Dict[str, Any]]: 发现的适配器信息列表
        """
        discovered = []
        
        for search_path in self.search_paths:
            path = Path(search_path).expanduser()
            if not path.exists():
                continue
            
            # 扫描目录
            for adapter_dir in path.iterdir():
                if not adapter_dir.is_dir():
                    continue
                
                # 检查是否包含manifest.yaml
                manifest_file = adapter_dir / "manifest.yaml"
                if not manifest_file.exists():
                    continue
                
                # 读取manifest
                with open(manifest_file, 'r', encoding='utf-8') as f:
                    manifest = yaml.safe_load(f)
                
                discovered.append({
                    "path": str(adapter_dir),
                    "manifest": manifest,
                    "source": "local"
                })
        
        return discovered
    
    def load_adapter_from_path(
        self, 
        adapter_path: str
    ) -> DomainAdapter:
        """
        从路径加载适配器
        
        Args:
            adapter_path: 适配器目录路径
            
        Returns:
            DomainAdapter: 适配器实例
            
        Raises:
            ImportError: 无法加载适配器
        """
        adapter_path = Path(adapter_path)
        
        # 读取manifest
        manifest_file = adapter_path / "manifest.yaml"
        with open(manifest_file, 'r', encoding='utf-8') as f:
            manifest = yaml.safe_load(f)
        
        # 加载adapter.py模块
        adapter_module_file = adapter_path / "adapter.py"
        if not adapter_module_file.exists():
            raise ImportError(f"adapter.py not found in {adapter_path}")
        
        # 动态导入
        spec = importlib.util.spec_from_file_location(
            manifest['name'],
            adapter_module_file
        )
        module = importlib.util.module_from_spec(spec)
        sys.modules[manifest['name']] = module
        spec.loader.exec_module(module)
        
        # 查找适配器类（约定：模块中只有一个DomainAdapter子类）
        adapter_class = None
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            if (isinstance(attr, type) and 
                issubclass(attr, DomainAdapter) and 
                attr is not DomainAdapter):
                adapter_class = attr
                break
        
        if adapter_class is None:
            raise ImportError(
                f"No DomainAdapter subclass found in {adapter_module_file}"
            )
        
        # 实例化
        adapter_instance = adapter_class()
        
        return adapter_instance
    
    def auto_register_all(self) -> Dict[str, Any]:
        """
        自动发现并注册所有适配器
        
        Returns:
            Dict[str, Any]: 注册结果统计
        """
        from .registry import adapter_registry
        
        discovered = self.discover_all()
        
        stats = {
            "discovered": len(discovered),
            "registered": 0,
            "failed": 0,
            "skipped": 0,
            "details": []
        }
        
        for item in discovered:
            try:
                manifest = item['manifest']
                
                # 检查是否已注册
                domain_name = manifest['name']
                domain_type = DomainType(domain_name)
                
                if adapter_registry.get_adapter(domain_type):
                    stats['skipped'] += 1
                    stats['details'].append({
                        "adapter": manifest['display_name'],
                        "status": "skipped",
                        "reason": "Already registered"
                    })
                    continue
                
                # 加载适配器
                adapter = self.load_adapter_from_path(item['path'])
                
                # 注册
                adapter_registry.register(domain_type, adapter)
                
                stats['registered'] += 1
                stats['details'].append({
                    "adapter": manifest['display_name'],
                    "version": manifest['version'],
                    "status": "registered"
                })
                
            except Exception as e:
                stats['failed'] += 1
                stats['details'].append({
                    "adapter": item.get('manifest', {}).get('display_name', 'Unknown'),
                    "status": "failed",
                    "error": str(e)
                })
        
        return stats


# 全局发现器实例
adapter_discovery = AdapterDiscovery()
```

### 2.1.3 依赖管理

```python
# src/adapters/dependency.py

from typing import List, Dict, Set, Any
from packaging import version, specifiers
import networkx as nx

class DependencyResolver:
    """
    适配器依赖解析器
    
    功能：
    1. 解析适配器依赖关系
    2. 检测循环依赖
    3. 计算安装顺序
    4. 验证版本兼容性
    """
    
    def __init__(self):
        self.dependency_graph = nx.DiGraph()
    
    def add_adapter(
        self, 
        adapter_name: str, 
        dependencies: Dict[str, str]
    ) -> None:
        """
        添加适配器到依赖图
        
        Args:
            adapter_name: 适配器名称
            dependencies: 依赖字典 {adapter_name: version_spec}
        """
        # 添加节点
        if adapter_name not in self.dependency_graph:
            self.dependency_graph.add_node(adapter_name)
        
        # 添加边
        for dep_name, version_spec in dependencies.items():
            self.dependency_graph.add_edge(
                adapter_name, 
                dep_name,
                version_spec=version_spec
            )
    
    def check_circular_dependencies(self) -> List[List[str]]:
        """
        检查循环依赖
        
        Returns:
            List[List[str]]: 循环依赖链列表，空列表表示无循环
        """
        try:
            cycles = list(nx.simple_cycles(self.dependency_graph))
            return cycles
        except:
            return []
    
    def get_install_order(
        self, 
        adapter_names: List[str]
    ) -> List[str]:
        """
        计算安装顺序（拓扑排序）
        
        Args:
            adapter_names: 要安装的适配器列表
            
        Returns:
            List[str]: 按依赖顺序排列的适配器名称列表
            
        Raises:
            ValueError: 存在循环依赖
        """
        # 检查循环依赖
        cycles = self.check_circular_dependencies()
        if cycles:
            raise ValueError(f"Circular dependencies detected: {cycles}")
        
        # 获取子图
        subgraph = self.dependency_graph.subgraph(adapter_names)
        
        # 拓扑排序
        try:
            install_order = list(nx.topological_sort(subgraph))
            # 反转（依赖项在前）
            return list(reversed(install_order))
        except nx.NetworkXError as e:
            raise ValueError(f"Cannot determine install order: {e}")
    
    def verify_compatibility(
        self,
        adapter_name: str,
        adapter_version: str,
        installed_adapters: Dict[str, str]
    ) -> Tuple[bool, List[str]]:
        """
        验证版本兼容性
        
        Args:
            adapter_name: 适配器名称
            adapter_version: 适配器版本
            installed_adapters: 已安装的适配器 {name: version}
            
        Returns:
            Tuple[bool, List[str]]: (是否兼容, 不兼容原因列表)
        """
        issues = []
        
        # 获取依赖
        if adapter_name not in self.dependency_graph:
            return True, []
        
        dependencies = self.dependency_graph[adapter_name]
        
        for dep_name in dependencies:
            edge_data = self.dependency_graph[adapter_name][dep_name]
            required_version_spec = edge_data.get('version_spec', '*')
            
            # 检查是否已安装
            if dep_name not in installed_adapters:
                issues.append(
                    f"Missing dependency: {dep_name} {required_version_spec}"
                )
                continue
            
            # 检查版本是否兼容
            installed_version = installed_adapters[dep_name]
            
            try:
                spec = specifiers.SpecifierSet(required_version_spec)
                if not spec.contains(installed_version):
                    issues.append(
                        f"Version mismatch: {dep_name} requires {required_version_spec}, "
                        f"but {installed_version} is installed"
                    )
            except Exception as e:
                issues.append(
                    f"Invalid version specifier for {dep_name}: {e}"
                )
        
        return len(issues) == 0, issues


# 全局依赖解析器
dependency_resolver = DependencyResolver()
```

---

## 2.2 适配器生命周期

### 2.2.1 生命周期状态机

```
┌─────────────┐
│  DISCOVERED │  发现：在文件系统或市场中被发现
└──────┬──────┘
       │
       │ 下载/验证
       ▼
┌─────────────┐
│  DOWNLOADED │  已下载：文件已下载到本地
└──────┬──────┘
       │
       │ 安装依赖
       ▼
┌─────────────┐
│  INSTALLED  │  已安装：依赖已安装，但未加载
└──────┬──────┘
       │
       │ 加载模块
       ▼
┌─────────────┐
│   LOADED    │  已加载：模块已导入，但未初始化
└──────┬──────┘
       │
       │ 初始化
       ▼
┌─────────────┐
│   ACTIVE    │  激活：正在使用中
└──────┬──────┘
       │
       │ 停用
       ▼
┌─────────────┐
│  INACTIVE   │  停用：已停止，但未卸载
└──────┬──────┘
       │
       │ 卸载
       ▼
┌─────────────┐
│ UNINSTALLED │  已卸载：已从系统移除
└─────────────┘
```

### 2.2.2 生命周期管理器

```python
# src/adapters/lifecycle.py

from enum import Enum
from typing import Dict, Optional
from datetime import datetime

class AdapterStatus(Enum):
    """适配器状态"""
    DISCOVERED = "discovered"
    DOWNLOADED = "downloaded"
    INSTALLED = "installed"
    LOADED = "loaded"
    ACTIVE = "active"
    INACTIVE = "inactive"
    UNINSTALLED = "uninstalled"
    ERROR = "error"


class AdapterLifecycleManager:
    """适配器生命周期管理器"""
    
    def __init__(self):
        self._status: Dict[str, AdapterStatus] = {}
        self._history: Dict[str, List[Dict]] = {}
    
    async def install(
        self, 
        adapter_name: str,
        source: str
    ) -> bool:
        """
        安装适配器
        
        步骤：
        1. 下载（如果是远程）
        2. 验证签名和完整性
        3. 安装依赖
        4. 运行安装钩子
        5. 更新状态
        
        Args:
            adapter_name: 适配器名称
            source: 来源（local/marketplace/github）
            
        Returns:
            bool: 是否成功
        """
        try:
            self._update_status(adapter_name, AdapterStatus.DISCOVERED)
            
            # 1. 下载
            if source != "local":
                await self._download_adapter(adapter_name, source)
            self._update_status(adapter_name, AdapterStatus.DOWNLOADED)
            
            # 2. 验证
            if not await self._verify_adapter(adapter_name):
                raise RuntimeError("Adapter verification failed")
            
            # 3. 安装依赖
            await self._install_dependencies(adapter_name)
            self._update_status(adapter_name, AdapterStatus.INSTALLED)
            
            # 4. 运行安装钩子
            adapter = self._load_adapter_module(adapter_name)
            await adapter.on_install()
            
            return True
            
        except Exception as e:
            self._update_status(adapter_name, AdapterStatus.ERROR)
            self._record_error(adapter_name, str(e))
            return False
    
    async def activate(self, adapter_name: str) -> bool:
        """激活适配器"""
        try:
            # 加载
            adapter = self._load_adapter_module(adapter_name)
            self._update_status(adapter_name, AdapterStatus.LOADED)
            
            # 初始化
            await adapter.initialize()
            self._update_status(adapter_name, AdapterStatus.ACTIVE)
            
            # 注册
            from .registry import adapter_registry
            domain = self._get_adapter_domain(adapter_name)
            adapter_registry.register(domain, adapter)
            
            return True
            
        except Exception as e:
            self._update_status(adapter_name, AdapterStatus.ERROR)
            return False
    
    async def deactivate(self, adapter_name: str) -> bool:
        """停用适配器"""
        try:
            from .registry import adapter_registry
            
            domain = self._get_adapter_domain(adapter_name)
            adapter_registry.unregister(domain)
            
            self._update_status(adapter_name, AdapterStatus.INACTIVE)
            return True
            
        except Exception as e:
            return False
    
    async def uninstall(self, adapter_name: str) -> bool:
        """卸载适配器"""
        try:
            # 先停用
            await self.deactivate(adapter_name)
            
            # 运行卸载钩子
            adapter = self._load_adapter_module(adapter_name)
            await adapter.on_uninstall()
            
            # 删除文件
            self._remove_adapter_files(adapter_name)
            
            self._update_status(adapter_name, AdapterStatus.UNINSTALLED)
            return True
            
        except Exception as e:
            return False
    
    def get_status(self, adapter_name: str) -> Optional[AdapterStatus]:
        """获取适配器状态"""
        return self._status.get(adapter_name)
    
    def _update_status(
        self, 
        adapter_name: str, 
        status: AdapterStatus
    ) -> None:
        """更新状态"""
        self._status[adapter_name] = status
        
        # 记录历史
        if adapter_name not in self._history:
            self._history[adapter_name] = []
        
        self._history[adapter_name].append({
            "status": status.value,
            "timestamp": datetime.utcnow().isoformat()
        })


# 全局生命周期管理器
lifecycle_manager = AdapterLifecycleManager()
```

---

*文档第三部分完成*

**已完成内容**：
- ✅ 注册表设计（AdapterRegistry）
- ✅ 自动发现机制（AdapterDiscovery）
- ✅ 依赖管理（DependencyResolver）
- ✅ 生命周期管理（AdapterLifecycleManager）

**下一部分将包含**：
- 适配器市场架构设计
- 发布和审核流程
- 用户安装和评价系统

是否继续输出第四部分（适配器市场设计）？
