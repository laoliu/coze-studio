# 领域适配器使用指南

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 架构说明

---

## 📋 目录

1. [适配器的职责](#适配器的职责)
2. [正确的使用方式](#正确的使用方式)
3. [错误的使用方式](#错误的使用方式)
4. [完整示例](#完整示例)
5. [扩展开发](#扩展开发)

---

## 适配器的职责

### ❓ 你的疑问是对的！

**问题**: 文档中的示例看起来像是每个作品都要写一个适配器？

```python
# ❌ 这样的示例可能产生误解
class ArtHistoryAdapter(DomainAdapter):
    async def parse_request(self, raw_input, preferences):
        return UserRequest(
            domain=DomainType.ART_HISTORY,
            art_movement="印象派",        # 写死了？
            artist="克劳德·莫奈",        # 写死了？
            artwork="日出·印象",          # 写死了？
            audience_type="成人艺术爱好者",
            duration=30
        )
```

**答案**: ❌ **不是！** 这些值应该从 `raw_input` 中**动态提取**，而不是写死！

---

## 正确的使用方式

### 适配器的真正职责

适配器是**解析器 + 提取器**，而不是硬编码的数据容器！

```
用户输入（raw_input）
    ↓
适配器.parse_request()  【提取规则】
    ↓
UserRequest（结构化数据）
```

### 正确的实现方式

#### ✅ 美术史适配器（正确版本）

```python
class ArtHistoryAdapter(DomainAdapter):
    """美术史领域适配器"""
    
    # 艺术流派关键词库（领域知识）
    ART_MOVEMENTS = {
        "印象派": ["印象派", "莫奈", "雷诺阿", "德加", "马奈"],
        "立体派": ["立体派", "毕加索", "布拉克"],
        "表现主义": ["表现主义", "蒙克", "康定斯基"],
        "文艺复兴": ["文艺复兴", "达芬奇", "米开朗基罗", "拉斐尔"],
        # ... 更多流派
    }
    
    # 著名艺术家数据库（领域知识）
    ARTISTS = {
        "莫奈": {
            "full_name": "克劳德·莫奈",
            "movement": "印象派",
            "famous_works": ["日出·印象", "睡莲", "干草堆"],
            "period": "1840-1926"
        },
        "梵高": {
            "full_name": "文森特·梵高",
            "movement": "后印象派",
            "famous_works": ["星空", "向日葵", "自画像"],
            "period": "1853-1890"
        },
        # ... 更多艺术家
    }
    
    async def parse_request(
        self,
        raw_input: str,
        preferences: Dict[str, Any]
    ) -> UserRequest:
        """
        从用户输入中提取美术史相关信息
        
        输入示例:
        - "讲解莫奈的《日出·印象》"
        - "分析印象派的艺术特点"
        - "介绍文艺复兴时期的绘画"
        """
        # 1. 提取艺术流派（动态识别）
        art_movement = self._extract_movement(raw_input)
        
        # 2. 提取艺术家（动态识别）
        artist = self._extract_artist(raw_input)
        
        # 3. 提取作品（动态识别）
        artwork = self._extract_artwork(raw_input, artist)
        
        # 4. 推断受众类型（基于上下文和preferences）
        audience_type = self._infer_audience(raw_input, preferences)
        
        # 5. 提取或使用默认时长
        duration = self._extract_duration(raw_input, preferences)
        
        # 6. 构建结构化请求
        return UserRequest(
            raw_input=raw_input,
            domain=DomainType.ART_HISTORY,
            preferences={
                **preferences,
                "art_movement": art_movement,
                "artist": artist,
                "artwork": artwork,
                "audience_type": audience_type,
                "duration": duration
            }
        )
    
    def _extract_movement(self, text: str) -> Optional[str]:
        """从文本中提取艺术流派"""
        for movement, keywords in self.ART_MOVEMENTS.items():
            if any(kw in text for kw in keywords):
                return movement
        return None
    
    def _extract_artist(self, text: str) -> Optional[str]:
        """从文本中提取艺术家"""
        for artist_key, artist_info in self.ARTISTS.items():
            # 检查简称或全名
            if artist_key in text or artist_info["full_name"] in text:
                return artist_info["full_name"]
        return None
    
    def _extract_artwork(
        self,
        text: str,
        artist: Optional[str]
    ) -> Optional[str]:
        """从文本中提取作品名"""
        # 使用正则提取《作品名》格式
        import re
        match = re.search(r'《(.+?)》', text)
        if match:
            artwork_title = match.group(1)
            
            # 如果已知艺术家，验证作品是否属于该艺术家
            if artist:
                artist_key = self._get_artist_key(artist)
                if artist_key and artwork_title in self.ARTISTS[artist_key]["famous_works"]:
                    return artwork_title
            
            return artwork_title
        
        return None
    
    def _infer_audience(
        self,
        text: str,
        preferences: Dict[str, Any]
    ) -> str:
        """推断受众类型"""
        # 优先使用preferences
        if "audience_type" in preferences:
            return preferences["audience_type"]
        
        # 从文本中推断
        if any(kw in text for kw in ["儿童", "小朋友", "幼儿"]):
            return "儿童艺术启蒙"
        elif any(kw in text for kw in ["学生", "青少年", "中学"]):
            return "青少年艺术教育"
        elif any(kw in text for kw in ["专业", "研究", "学术"]):
            return "专业艺术研究"
        else:
            return "成人艺术爱好者"
    
    def _extract_duration(
        self,
        text: str,
        preferences: Dict[str, Any]
    ) -> int:
        """提取时长"""
        # 优先使用preferences
        if "duration" in preferences:
            return preferences["duration"]
        
        # 从文本中提取
        import re
        if match := re.search(r'(\d+)\s*分钟', text):
            return int(match.group(1))
        
        # 默认值（美术史通常需要更长时间）
        return 30
```

---

## 完整示例

### 示例1：动态解析不同的美术作品

```python
adapter = ArtHistoryAdapter(config={})

# ========== 案例1：莫奈的《日出·印象》 ==========
request1 = await adapter.parse_request(
    raw_input="讲解莫奈的《日出·印象》",
    preferences={"duration": 25}
)
print(request1.preferences)
# 输出:
# {
#     "art_movement": "印象派",
#     "artist": "克劳德·莫奈",
#     "artwork": "日出·印象",
#     "audience_type": "成人艺术爱好者",
#     "duration": 25
# }

# ========== 案例2：梵高的《星空》 ==========
request2 = await adapter.parse_request(
    raw_input="给儿童讲解梵高的《星空》，15分钟",
    preferences={}
)
print(request2.preferences)
# 输出:
# {
#     "art_movement": "后印象派",
#     "artist": "文森特·梵高",
#     "artwork": "星空",
#     "audience_type": "儿童艺术启蒙",  # 从"儿童"关键词推断
#     "duration": 15
# }

# ========== 案例3：文艺复兴时期绘画 ==========
request3 = await adapter.parse_request(
    raw_input="分析文艺复兴时期的绘画特点",
    preferences={"audience_type": "专业艺术研究"}
)
print(request3.preferences)
# 输出:
# {
#     "art_movement": "文艺复兴",
#     "artist": None,  # 未指定具体艺术家
#     "artwork": None,  # 未指定具体作品
#     "audience_type": "专业艺术研究",
#     "duration": 30  # 默认值
# }
```

---

### 示例2：K12适配器的动态解析

```python
adapter = K12Adapter(config={})

# ========== 案例1：初三化学 ==========
request1 = await adapter.parse_request(
    raw_input="讲解初三化学燃烧条件，10分钟",
    preferences={}
)
print(request1.preferences)
# 输出:
# {
#     "subject": "化学",
#     "grade": 9,  # 从"初三"提取
#     "duration": 10
# }

# ========== 案例2：高一物理 ==========
request2 = await adapter.parse_request(
    raw_input="高一物理牛顿第一定律，需要实验演示",
    preferences={"duration": 45}
)
print(request2.preferences)
# 输出:
# {
#     "subject": "物理",
#     "grade": 10,  # 从"高一"提取
#     "duration": 45,
#     "needs_experiment": True  # 可以进一步提取这类需求
# }
```

---

## 错误的使用方式

### ❌ 反例1：硬编码数据

```python
# ❌ 错误！每个作品都要写一个适配器？
class MonetSunriseAdapter(DomainAdapter):
    async def parse_request(self, raw_input, preferences):
        return UserRequest(
            domain=DomainType.ART_HISTORY,
            art_movement="印象派",
            artist="莫奈",
            artwork="日出·印象"
        )

class VanGoghStarryNightAdapter(DomainAdapter):
    async def parse_request(self, raw_input, preferences):
        return UserRequest(
            domain=DomainType.ART_HISTORY,
            art_movement="后印象派",
            artist="梵高",
            artwork="星空"
        )

# ❌ 这样就需要成千上万个适配器！
```

**问题**:
- 🚫 维护噩梦（每个作品一个类）
- 🚫 无法扩展（新作品需要修改代码）
- 🚫 违背适配器设计原则

---

### ❌ 反例2：适配器职责过重

```python
# ❌ 错误！适配器不应该包含业务逻辑
class ArtHistoryAdapter(DomainAdapter):
    async def parse_request(self, raw_input, preferences):
        # ❌ 不应该在适配器中生成内容
        content = await self._generate_artwork_analysis(raw_input)
        
        # ❌ 不应该在适配器中调用外部API
        artwork_image = await self._fetch_from_wikiart(raw_input)
        
        return UserRequest(...)
```

**问题**:
- 🚫 适配器应该只做解析和提取
- 🚫 业务逻辑应该在工作流引擎中
- 🚫 违背单一职责原则

---

## 架构关系图

### 正确的数据流

```
┌─────────────────────────────────────────────────────┐
│  用户输入 (raw_input)                                │
│  "讲解莫奈的《日出·印象》，给成人讲解，30分钟"          │
└─────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────┐
│  领域适配器 (DomainAdapter)                          │
│  ┌───────────────────────────────────────────────┐  │
│  │ parse_request()                              │  │
│  │  - 提取艺术流派: _extract_movement()          │  │
│  │  - 提取艺术家: _extract_artist()              │  │
│  │  - 提取作品: _extract_artwork()               │  │
│  │  - 推断受众: _infer_audience()                │  │
│  │  - 提取时长: _extract_duration()              │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────┐
│  UserRequest (结构化数据)                            │
│  {                                                  │
│    domain: "art_history",                          │
│    preferences: {                                  │
│      art_movement: "印象派",    ← 动态提取           │
│      artist: "克劳德·莫奈",     ← 动态提取           │
│      artwork: "日出·印象",      ← 动态提取           │
│      audience_type: "成人",    ← 动态推断           │
│      duration: 30              ← 动态提取           │
│    }                                               │
│  }                                                 │
└─────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────┐
│  AI能力层 / 工作流引擎                               │
│  - 使用这些结构化数据                                │
│  - 生成学习目标                                     │
│  - 发现内容                                         │
│  - 生成叙述                                         │
└─────────────────────────────────────────────────────┘
```

---

## 扩展开发

### 如何添加新的领域知识

#### 方法1：配置文件驱动（推荐）

```python
class ArtHistoryAdapter(DomainAdapter):
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        
        # 从配置文件加载领域知识
        self.knowledge_base = self._load_knowledge_base(
            config.get("knowledge_base_path", "config/art_history_kb.yaml")
        )
    
    def _load_knowledge_base(self, path: str) -> Dict:
        """从YAML文件加载知识库"""
        import yaml
        with open(path, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
```

**配置文件示例** (`config/art_history_kb.yaml`):
```yaml
art_movements:
  印象派:
    keywords: [印象派, 莫奈, 雷诺阿, 德加, 马奈]
    period: "1860-1890"
    characteristics:
      - 捕捉光影变化
      - 户外写生
      - 短促笔触
  
  立体派:
    keywords: [立体派, 毕加索, 布拉克]
    period: "1907-1920"
    characteristics:
      - 几何形状
      - 多视角
      - 抽象表现

artists:
  莫奈:
    full_name: "克劳德·莫奈"
    movement: 印象派
    famous_works:
      - 日出·印象
      - 睡莲系列
      - 干草堆系列
    period: "1840-1926"
    
  梵高:
    full_name: "文森特·梵高"
    movement: 后印象派
    famous_works:
      - 星空
      - 向日葵
      - 自画像
    period: "1853-1890"
```

**优点**:
- ✅ 知识与代码分离
- ✅ 易于维护和更新
- ✅ 支持多语言（i18n）
- ✅ 可以从数据库加载

---

#### 方法2：数据库驱动（企业级）

```python
class ArtHistoryAdapter(DomainAdapter):
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        
        # 连接知识库数据库
        self.knowledge_repo = ArtHistoryKnowledgeRepository(
            db_session=config.get("db_session")
        )
    
    async def _extract_artist(self, text: str) -> Optional[str]:
        """从数据库中查询艺术家"""
        # 从数据库查询（支持模糊匹配、别名等）
        artists = await self.knowledge_repo.search_artists(text)
        
        if artists:
            return artists[0].full_name
        
        return None
```

**优点**:
- ✅ 支持大规模数据
- ✅ 实时更新（无需重启）
- ✅ 支持复杂查询
- ✅ 支持多租户

---

#### 方法3：AI辅助提取（未来增强）

```python
class ArtHistoryAdapter(DomainAdapter):
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        
        # 使用AI能力层的意图识别器
        from ..ai import intent_recognizer
        self.intent_recognizer = intent_recognizer
    
    async def _extract_artist(self, text: str) -> Optional[str]:
        """使用AI提取艺术家"""
        # 如果关键词匹配失败，使用AI
        if not self._keyword_match(text):
            intent = await self.intent_recognizer.recognize(text)
            return intent.entities.get("artist")
        
        return self._keyword_match(text)
```

**优点**:
- ✅ 处理模糊输入
- ✅ 理解上下文
- ✅ 支持多语言
- ✅ 自动学习新知识

---

## 最佳实践

### ✅ DO（应该做）

1. **动态提取信息**
   ```python
   art_movement = self._extract_movement(raw_input)  # ✅
   ```

2. **使用领域知识库**
   ```python
   SUBJECT_KEYWORDS = {...}  # ✅ 集中管理
   ```

3. **支持多种输入格式**
   ```python
   # 支持: "初三"、"初3"、"9年级"
   grade = self._extract_grade(raw_input)  # ✅
   ```

4. **提供默认值和容错**
   ```python
   duration = self._extract_duration(raw_input, preferences)
   if duration is None:
       duration = 30  # ✅ 默认值
   ```

5. **使用preferences覆盖**
   ```python
   # 优先使用用户明确指定的值
   if "audience_type" in preferences:
       return preferences["audience_type"]  # ✅
   ```

---

### ❌ DON'T（不应该做）

1. **硬编码具体数据**
   ```python
   art_movement = "印象派"  # ❌ 不要写死
   ```

2. **在适配器中调用AI能力层**
   ```python
   # ❌ 适配器不应该生成内容
   content = await content_discovery.search(...)
   ```

3. **在适配器中实现业务逻辑**
   ```python
   # ❌ 不应该在适配器中生成学习目标
   objectives = self._generate_objectives(...)
   ```

4. **为每个实例创建适配器**
   ```python
   # ❌ 不要这样
   class MonetAdapter(DomainAdapter): pass
   class VanGoghAdapter(DomainAdapter): pass
   ```

---

## 总结

### 核心原则

1. **适配器 = 解析器**
   - 从自然语言提取结构化信息
   - 应用领域特定的解析规则
   - 不包含业务逻辑

2. **领域知识 = 可配置的**
   - 关键词库
   - 规则映射
   - 默认值

3. **一个领域 = 一个适配器**
   - ArtHistoryAdapter 处理所有美术史内容
   - K12Adapter 处理所有K12教育内容
   - 不是一个作品一个适配器！

### 记忆口诀

```
适配器是解析器，不是数据容器
从输入中提取，不是写死数据
领域知识外置化，配置文件或数据库
一个领域一适配，千万别搞成千个
```

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow Team
