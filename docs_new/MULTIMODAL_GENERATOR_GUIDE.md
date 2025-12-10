# 多模态生成器调用场景指南

**版本**: 1.0  
**日期**: 2025-12-10  
**状态**: 架构设计说明

---

## 📋 目录

1. [什么是多模态生成器](#什么是多模态生成器)
2. [在系统中的定位](#在系统中的定位)
3. [典型调用场景](#典型调用场景)
4. [调用时机决策树](#调用时机决策树)
5. [实际案例分析](#实际案例分析)
6. [技术实现方案](#技术实现方案)

---

## 什么是多模态生成器

### 定义

**多模态生成器 (Multimodal Generator)** 是一个AI能力层组件，负责生成**非文本**的富媒体内容。

```
┌─────────────────────────────────────────────────┐
│       多模态生成器 (Multimodal Generator)        │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌──────────────┐  ┌──────────────┐            │
│  │ 图片生成器    │  │ 音频生成器    │            │
│  │ Image Gen    │  │ Audio Gen    │            │
│  └──────────────┘  └──────────────┘            │
│                                                  │
│  ┌──────────────┐  ┌──────────────┐            │
│  │ 视频生成器    │  │ PPT生成器     │            │
│  │ Video Gen    │  │ PPT Gen      │            │
│  └──────────────┘  └──────────────┘            │
│                                                  │
└─────────────────────────────────────────────────┘
```

### 核心能力

| 能力 | 说明 | 技术方案 |
|------|------|---------|
| **图片生成** | 配图、示意图、插画 | DALL-E 3, Midjourney, Stable Diffusion |
| **音频生成** | 文本转语音(TTS)、背景音乐 | Azure TTS, ElevenLabs, OpenAI TTS |
| **视频生成** | 动画、演示视频、课件视频 | FFmpeg, D-ID, Runway |
| **PPT生成** | 课件、培训材料、演示文稿 | python-pptx, 模板引擎 |
| **PDF生成** | 讲义、手册、报告 | ReportLab, WeasyPrint |

---

## 在系统中的定位

### 架构层级

```
┌────────────────────────────────────────────────┐
│              应用层 (API)                       │
│          用户请求 → 返回完整内容                 │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│          领域适配器层 (Domain Adapter)          │
│        K12适配器、美术史适配器等                 │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│             AI能力层 (AI Capabilities)          │
│  ┌──────────────────┐  ┌──────────────────┐   │
│  │ 学习目标生成      │  │ 内容发现         │   │
│  └──────────────────┘  └──────────────────┘   │
│  ┌──────────────────┐  ┌──────────────────┐   │
│  │ 故事化叙述       │  │ 多模态生成 ⭐     │   │
│  └──────────────────┘  └──────────────────┘   │
└────────────────────────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│          工作流引擎层 (Workflow Engine)         │
│         节点编排 → 调用各种能力                  │
└────────────────────────────────────────────────┘
```

### 与其他能力的关系

| AI能力 | 输出 | 调用多模态生成器？ |
|--------|------|------------------|
| **意图识别** | 结构化意图数据 | ❌ 否 |
| **学习目标生成** | 学习目标列表 | ❌ 否 |
| **内容发现** | 文本内容 | ❌ 否 |
| **故事化叙述** | 文本故事 | ❌ 否 |
| **质量评估** | 评分和建议 | ❌ 否 |
| **多模态生成** | 图片/音频/视频/PPT | - |

**关键点**: 多模态生成器是在**文本内容已经准备好**之后，为了**增强呈现效果**而调用的。

---

## 典型调用场景

### 场景分类

根据实际需求，多模态生成器的调用可以分为以下场景：

```
调用场景分类
├── 1. 内容增强场景（配图、配音）
├── 2. 格式转换场景（导出PPT、PDF）
├── 3. 交互增强场景（动画、视频）
└── 4. 可访问性场景（TTS、字幕）
```

---

### 场景1️⃣: 内容增强场景

**目的**: 为文本内容配图、配音，增强视觉和听觉效果

#### 1.1 配图生成

**调用时机**: 文本内容生成完成后，需要视觉辅助

**示例**:
```python
# ========== 工作流执行流程 ==========

# Step 1: 生成文本内容
text_content = await narrative_generator.generate(context)
# 输出: "莫奈站在勒阿弗尔港口，面对清晨的雾气..."

# Step 2: 为文本配图 ⭐ 调用多模态生成器
if context.preferences.get("include_images", True):
    images = await multimodal_generator.generate_images([
        {
            "prompt": "莫奈在勒阿弗尔港口作画，印象派风格",
            "style": "impressionist_painting",
            "size": "1024x1024"
        },
        {
            "prompt": "《日出·印象》高清作品图",
            "style": "photorealistic",
            "size": "1920x1080"
        }
    ])
    
    # Step 3: 嵌入到内容中
    enhanced_content = content_assembler.embed_images(
        text=text_content,
        images=images
    )
```

**适用场景**:
- ✅ **美术史**: 需要高清作品图、艺术家肖像
- ✅ **K12理科**: 需要实验器材图、分子结构图
- ✅ **职业培训**: 需要软件截图、操作示意图
- ✅ **畅销书解读**: 需要封面图、作者照片

---

#### 1.2 TTS音频生成

**调用时机**: 需要语音朗读、有声内容

**示例**:
```python
# ========== 场景：儿童艺术启蒙 ==========

# Step 1: 生成适合儿童的故事化内容
story = await narrative_generator.generate(
    context=context,
    style="children_friendly"
)
# 输出: "很久很久以前，有一位叫莫奈的画家..."

# Step 2: 生成朗读音频 ⭐ 调用多模态生成器
if context.audience_type == "儿童艺术启蒙":
    audio = await multimodal_generator.generate_audio(
        text=story,
        voice="children_narrator",  # 儿童友好的声音
        speed=0.9,  # 稍慢语速
        emotion="warm"  # 温暖的情感
    )
    
    # Step 3: 合并到内容包
    content_package = {
        "text": story,
        "audio": audio,
        "duration": audio.duration
    }
```

**适用场景**:
- ✅ **儿童教育**: 需要温柔的朗读声音
- ✅ **语言学习**: 需要标准发音示范
- ✅ **盲人可访问**: 需要TTS朗读全文
- ✅ **有声课程**: 需要讲师语音

---

### 场景2️⃣: 格式转换场景

**目的**: 将结构化内容导出为特定格式（PPT、PDF、视频等）

#### 2.1 PPT课件生成

**调用时机**: 教师需要导出为PPT格式，用于课堂演示

**示例**:
```python
# ========== 场景：K12化学课PPT生成 ==========

# Step 1-4: 正常工作流（生成文本内容）
workflow_result = await workflow_engine.execute(context)
# 包含: 学习目标、教学内容、活动设计等

# Step 5: 用户点击"导出为PPT" ⭐ 调用多模态生成器
ppt = await multimodal_generator.generate_ppt(
    content=workflow_result,
    template="k12_chemistry_template",  # K12化学模板
    options={
        "include_title_slide": True,
        "include_objectives_slide": True,
        "include_images": True,
        "include_experiment_photos": True,
        "theme": "modern_education",
        "language": "zh-CN"
    }
)

# 生成的PPT结构:
# Slide 1: 标题页（课题名称、教师信息）
# Slide 2: 学习目标（3-5个目标）
# Slide 3-5: 知识讲解（配图、示意图）
# Slide 6: 实验演示（步骤图、安全提示）
# Slide 7: 练习题（选择题、简答题）
# Slide 8: 课堂小结

# Step 6: 下载PPT文件
download_url = await storage.upload(ppt, "lesson_plan.pptx")
```

**适用场景**:
- ✅ **K12教育**: 教师需要PPT课件上课
- ✅ **企业培训**: 培训师需要PPT材料
- ✅ **会议演讲**: 需要演示文稿
- ✅ **产品演示**: 需要产品介绍PPT

---

#### 2.2 PDF讲义生成

**调用时机**: 需要可打印的文档、学习材料

**示例**:
```python
# ========== 场景：畅销书读书会讲义 ==========

# Step 1-4: 生成读书会内容
content = await workflow_engine.execute(context)
# 包含: 书籍摘要、讨论问题、阅读笔记等

# Step 5: 导出为PDF讲义 ⭐ 调用多模态生成器
pdf = await multimodal_generator.generate_pdf(
    content=content,
    template="reading_club_handout",
    options={
        "include_cover": True,  # 封面（书籍封面图）
        "include_toc": True,  # 目录
        "include_quotes": True,  # 金句摘录
        "include_discussion_questions": True,
        "format": "A4",
        "orientation": "portrait",
        "font": "思源宋体"
    }
)

# 生成的PDF结构:
# Page 1: 封面（书名、作者、读书会主题）
# Page 2: 目录
# Page 3-5: 书籍摘要和核心观点
# Page 6-8: 讨论问题（10-15个问题）
# Page 9: 金句摘录
# Page 10: 阅读笔记模板

# Step 6: 下载PDF
download_url = await storage.upload(pdf, "reading_handout.pdf")
```

**适用场景**:
- ✅ **学习讲义**: 课后复习材料
- ✅ **培训手册**: 操作指南、参考手册
- ✅ **报告文档**: 研究报告、总结报告
- ✅ **电子书**: 系列课程整理成册

---

### 场景3️⃣: 交互增强场景

**目的**: 生成动态、交互式的多媒体内容

#### 3.1 教学动画视频

**调用时机**: 需要动态演示过程（实验、操作步骤等）

**示例**:
```python
# ========== 场景：化学燃烧实验视频 ==========

# Step 1-4: 生成实验步骤和讲解
experiment_content = await workflow_engine.execute(context)

# Step 5: 生成实验演示视频 ⭐ 调用多模态生成器
video = await multimodal_generator.generate_video(
    script=experiment_content.experiment_steps,
    options={
        "style": "animated_diagram",  # 动画示意图风格
        "duration": 180,  # 3分钟
        "include_narration": True,  # 包含旁白
        "include_subtitles": True,  # 包含字幕
        "include_safety_warnings": True,  # 安全提示
        "transitions": "smooth",
        "background_music": "educational_soft"
    }
)

# 视频内容:
# 0:00-0:30: 介绍燃烧的定义和条件
# 0:30-1:30: 动画演示燃烧三要素（可燃物、氧气、温度）
# 1:30-2:30: 实验步骤演示（点燃蜡烛、观察现象）
# 2:30-3:00: 总结和安全提示

# Step 6: 嵌入到课件或单独下载
video_url = await storage.upload(video, "combustion_demo.mp4")
```

**适用场景**:
- ✅ **科学实验**: 化学、物理实验演示
- ✅ **软件操作**: Excel、Photoshop操作教程
- ✅ **艺术创作**: 绘画过程、雕塑技法
- ✅ **运动技能**: 体育动作分解演示

---

#### 3.2 交互式演示

**调用时机**: 需要用户可以交互的内容（点击、拖拽等）

**示例**:
```python
# ========== 场景：虚拟化学实验室 ==========

# Step 1-4: 生成实验内容
lab_content = await workflow_engine.execute(context)

# Step 5: 生成交互式实验 ⭐ 调用多模态生成器
interactive_lab = await multimodal_generator.generate_interactive(
    content=lab_content,
    type="virtual_lab",
    options={
        "equipment": ["烧杯", "试管", "酒精灯", "火柴"],
        "chemicals": ["蜡烛", "氧气瓶"],
        "interactions": [
            "点击火柴点燃",
            "拖拽蜡烛到酒精灯",
            "观察燃烧现象"
        ],
        "feedback": True,  # 操作反馈
        "safety_check": True  # 安全检查（错误操作提示）
    }
)

# 生成:
# - HTML5交互式页面
# - 可嵌入iframe或导出为独立网页
# - 支持鼠标/触摸交互
```

**适用场景**:
- ✅ **虚拟实验**: 危险实验的安全模拟
- ✅ **交互练习**: 拖拽题、连线题
- ✅ **游戏化学习**: 知识问答、闯关游戏
- ✅ **3D模型**: 分子结构、建筑模型

---

### 场景4️⃣: 可访问性场景

**目的**: 提升内容的可访问性（盲人、听障人士等）

#### 4.1 无障碍TTS

**调用时机**: 内容需要支持视障用户

**示例**:
```python
# ========== 场景：视障用户访问美术史内容 ==========

# Step 1-4: 生成美术史讲解
art_content = await workflow_engine.execute(context)

# Step 5: 生成详细的音频描述 ⭐ 调用多模态生成器
accessible_audio = await multimodal_generator.generate_accessible_audio(
    content=art_content,
    options={
        "include_image_descriptions": True,  # 详细描述图片
        "include_navigation": True,  # 章节导航提示
        "voice": "professional_narrator",
        "format": "mp3",
        "chapters": True  # 分章节，便于跳转
    }
)

# 音频内容:
# 00:00 - 章节1：克劳德·莫奈生平简介
# 05:30 - 章节2：印象派的艺术特点
# 12:00 - 图片描述：《日出·印象》作品，画面中可以看到...
# 18:30 - 章节3：莫奈的艺术技法
```

**适用场景**:
- ✅ **视障用户**: 完整的音频版本
- ✅ **驾驶场景**: 开车时收听课程
- ✅ **多任务场景**: 做家务时听课
- ✅ **语言学习**: 听力练习材料

---

## 调用时机决策树

```
用户请求内容生成
    ↓
工作流执行（生成文本内容）
    ↓
内容已生成 ✅
    ↓
    ┌──────────────────────────────────────┐
    │ 决策1: 需要视觉增强吗？               │
    ├──────────────────────────────────────┤
    │ ✅ 美术史、理科、培训 → 调用图片生成   │
    │ ❌ 纯文字内容、文学 → 不调用           │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策2: 需要听觉增强吗？               │
    ├──────────────────────────────────────┤
    │ ✅ 儿童教育、语言学习 → 调用音频生成   │
    │ ❌ 成人阅读、书面材料 → 不调用         │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策3: 需要特定格式导出吗？           │
    ├──────────────────────────────────────┤
    │ ✅ 用户点击"导出PPT" → 调用PPT生成    │
    │ ✅ 用户点击"导出PDF" → 调用PDF生成    │
    │ ❌ 仅在线阅读 → 不调用                │
    └──────────────────────────────────────┘
    ↓
    ┌──────────────────────────────────────┐
    │ 决策4: 需要动态演示吗？               │
    ├──────────────────────────────────────┤
    │ ✅ 实验步骤、操作流程 → 调用视频生成   │
    │ ❌ 静态讲解 → 不调用                  │
    └──────────────────────────────────────┘
    ↓
返回最终内容（文本 + 多模态内容）
```

---

## 实际案例分析

### 案例1: K12化学课件生成

```python
# ========== 完整工作流 ==========

# 1. 用户输入
user_input = "讲解初三化学燃烧条件，10分钟"

# 2. 领域适配器解析
adapter = K12Adapter()
request = await adapter.parse_request(user_input, {})

# 3. AI能力层：生成文本内容
objectives = await objective_generator.generate(request)
content = await content_discovery.search(objectives)
narrative = await narrative_generator.generate(content)

# 4. 多模态生成器：增强内容 ⭐
# 4.1 生成配图（实验器材图）
images = await multimodal_generator.generate_images([
    {"prompt": "化学实验器材：烧杯、试管、酒精灯"},
    {"prompt": "蜡烛燃烧示意图"}
])

# 4.2 生成实验视频（可选）
if request.preferences.get("include_video"):
    video = await multimodal_generator.generate_video(
        script=narrative.experiment_section,
        style="animated_diagram"
    )

# 5. 用户导出为PPT ⭐
if request.output_format == "ppt":
    ppt = await multimodal_generator.generate_ppt(
        content={
            "objectives": objectives,
            "narrative": narrative,
            "images": images,
            "video": video
        },
        template="k12_chemistry"
    )
    
    return ppt  # 返回PPT文件
else:
    # 返回在线版本（HTML）
    return {
        "narrative": narrative,
        "images": images,
        "video": video
    }
```

**调用多模态生成器的时机**:
1. ⏰ **Step 4.1**: 文本内容生成完成后，配图增强
2. ⏰ **Step 4.2**: 根据需求，生成演示视频
3. ⏰ **Step 5**: 用户主动请求导出PPT

---

### 案例2: 美术史讲解（儿童版）

```python
# ========== 完整工作流 ==========

# 1. 用户输入
user_input = "给7岁儿童讲解莫奈的《日出·印象》，15分钟"

# 2. 领域适配器解析
adapter = ArtHistoryAdapter()
request = await adapter.parse_request(user_input, {
    "audience_type": "children"
})

# 3. AI能力层：生成儿童友好的故事
objectives = await objective_generator.generate(request)
content = await content_discovery.search(objectives)
story = await narrative_generator.generate(
    content,
    style="children_friendly"  # 故事化叙述
)
# 输出: "很久很久以前，有一位叫莫奈的画家..."

# 4. 多模态生成器：增强儿童版内容 ⭐
# 4.1 生成高清作品图
artwork_image = await multimodal_generator.generate_image({
    "prompt": "莫奈《日出·印象》高清作品图",
    "size": "1920x1080",
    "quality": "high"
})

# 4.2 生成儿童友好的TTS音频 ⭐ 关键！
audio = await multimodal_generator.generate_audio(
    text=story,
    voice="children_narrator",  # 温柔的儿童讲述者声音
    speed=0.85,  # 稍慢，便于理解
    emotion="warm_storytelling",  # 讲故事的温暖语气
    include_background_music=True,  # 轻柔背景音乐
    music_style="classical_soft"
)

# 4.3 生成动画视频（可选）
if request.preferences.get("include_animation"):
    animation = await multimodal_generator.generate_video(
        images=[artwork_image],
        narration=audio,
        style="children_animation",  # 儿童动画风格
        effects=["zoom_in", "pan", "fade"]  # 画面效果
    )

# 5. 返回完整的儿童版内容包
return {
    "text": story,  # 文字故事
    "images": [artwork_image],  # 高清作品图
    "audio": audio,  # ⭐ 关键：儿童友好的朗读音频
    "video": animation,  # 可选：动画视频
    "format": "interactive_storybook"  # 交互式故事书格式
}
```

**调用多模态生成器的理由**:
- ✅ **儿童受众**: 7岁儿童更适合听故事（而非阅读）
- ✅ **视觉学习**: 高清作品图帮助理解
- ✅ **沉浸体验**: 背景音乐增强故事氛围
- ✅ **可访问性**: 家长可以在睡前播放音频

---

### 案例3: 企业Excel培训（完整多模态）

```python
# ========== 完整工作流 ==========

# 1. 用户输入
user_input = "2小时Excel数据分析培训，包含实操演示"

# 2. 领域适配器解析
adapter = VocationalTrainingAdapter()
request = await adapter.parse_request(user_input, {
    "training_type": "excel_data_analysis",
    "duration": 120,
    "include_demo": True
})

# 3. AI能力层：生成培训内容
objectives = await objective_generator.generate(request)
content = await content_discovery.search(objectives)
training_materials = await narrative_generator.generate(content)

# 4. 多模态生成器：生成完整培训包 ⭐

# 4.1 生成培训PPT ⭐
ppt = await multimodal_generator.generate_ppt(
    content=training_materials,
    template="corporate_training",
    options={
        "sections": [
            "课程介绍",
            "理论讲解",
            "操作演示",
            "案例分析",
            "练习题",
            "总结"
        ],
        "include_screenshots": True,  # Excel截图
        "include_step_by_step": True,  # 操作步骤图
        "theme": "professional_blue"
    }
)

# 4.2 生成操作演示视频 ⭐
demo_videos = []
for section in training_materials.demo_sections:
    video = await multimodal_generator.generate_video(
        script=section.steps,
        type="screen_recording_simulation",  # 屏幕录制模拟
        options={
            "show_mouse_clicks": True,  # 显示鼠标点击
            "show_keyboard_input": True,  # 显示键盘输入
            "include_narration": True,  # 旁白讲解
            "highlight_important_areas": True,  # 高亮重点区域
            "duration": section.duration
        }
    )
    demo_videos.append(video)

# 4.3 生成练习数据文件 ⭐
practice_files = await multimodal_generator.generate_excel_files([
    {
        "name": "练习数据1_销售明细.xlsx",
        "data": training_materials.practice_data_1,
        "sheets": ["销售明细", "产品目录"]
    },
    {
        "name": "练习数据2_财务报表.xlsx",
        "data": training_materials.practice_data_2,
        "sheets": ["收入支出", "利润分析"]
    }
])

# 4.4 生成操作手册PDF ⭐
handbook = await multimodal_generator.generate_pdf(
    content=training_materials,
    template="operation_handbook",
    options={
        "include_screenshots": True,
        "include_shortcuts": True,  # 快捷键列表
        "include_troubleshooting": True,  # 常见问题
        "format": "A4",
        "bookmarks": True  # PDF书签
    }
)

# 5. 返回完整培训包
return {
    "ppt": ppt,  # ⭐ 培训PPT
    "demo_videos": demo_videos,  # ⭐ 操作演示视频
    "practice_files": practice_files,  # ⭐ 练习数据
    "handbook": handbook,  # ⭐ 操作手册PDF
    "estimated_duration": 120  # 2小时
}
```

**调用多模态生成器的价值**:
- ✅ **PPT**: 培训师授课使用
- ✅ **演示视频**: 学员课后复习操作步骤
- ✅ **练习文件**: 实操练习（带真实数据）
- ✅ **操作手册**: 工作中查阅参考

---

## 技术实现方案

### 架构设计

```python
class MultimodalGenerator:
    """多模态生成器（AI能力层）"""
    
    def __init__(self, config: Dict[str, Any]):
        # 图片生成引擎
        self.image_generator = ImageGenerator(
            provider=config.get("image_provider", "dalle3"),
            api_key=config.get("dalle_api_key")
        )
        
        # 音频生成引擎
        self.audio_generator = AudioGenerator(
            provider=config.get("audio_provider", "azure_tts"),
            api_key=config.get("tts_api_key")
        )
        
        # 视频生成引擎
        self.video_generator = VideoGenerator(
            provider=config.get("video_provider", "ffmpeg"),
            temp_dir=config.get("temp_dir", "/tmp")
        )
        
        # PPT生成引擎
        self.ppt_generator = PPTGenerator(
            template_dir=config.get("template_dir", "templates/ppt")
        )
        
        # PDF生成引擎
        self.pdf_generator = PDFGenerator(
            template_dir=config.get("template_dir", "templates/pdf")
        )
    
    async def generate_image(
        self,
        prompt: str,
        **options
    ) -> ImageContent:
        """生成单张图片"""
        return await self.image_generator.generate(prompt, **options)
    
    async def generate_images(
        self,
        prompts: List[Dict[str, Any]]
    ) -> List[ImageContent]:
        """批量生成图片"""
        tasks = [
            self.image_generator.generate(p["prompt"], **p)
            for p in prompts
        ]
        return await asyncio.gather(*tasks)
    
    async def generate_audio(
        self,
        text: str,
        **options
    ) -> AudioContent:
        """生成音频（TTS）"""
        return await self.audio_generator.text_to_speech(text, **options)
    
    async def generate_video(
        self,
        script: Union[str, List[Dict]],
        **options
    ) -> VideoContent:
        """生成视频"""
        return await self.video_generator.generate(script, **options)
    
    async def generate_ppt(
        self,
        content: Dict[str, Any],
        template: str,
        **options
    ) -> PPTContent:
        """生成PPT课件"""
        return await self.ppt_generator.generate(content, template, **options)
    
    async def generate_pdf(
        self,
        content: Dict[str, Any],
        template: str,
        **options
    ) -> PDFContent:
        """生成PDF文档"""
        return await self.pdf_generator.generate(content, template, **options)
```

### 调用示例（完整）

```python
# ========== 在工作流引擎中调用 ==========

class WorkflowEngine:
    def __init__(self):
        self.multimodal = MultimodalGenerator(config)
    
    async def execute(self, context: WorkflowContext):
        # Step 1-4: 生成文本内容（其他AI能力）
        text_content = await self._generate_text_content(context)
        
        # Step 5: 根据需求调用多模态生成器
        result = {
            "text": text_content
        }
        
        # 5.1 是否需要配图？
        if context.preferences.get("include_images", True):
            result["images"] = await self.multimodal.generate_images(
                self._extract_image_prompts(text_content)
            )
        
        # 5.2 是否需要音频？
        if context.preferences.get("include_audio", False):
            result["audio"] = await self.multimodal.generate_audio(
                text=text_content.narrative,
                voice=self._select_voice(context.audience_type)
            )
        
        # 5.3 是否需要导出PPT？
        if context.output_format == "ppt":
            result["ppt"] = await self.multimodal.generate_ppt(
                content=result,
                template=self._select_template(context.domain)
            )
        
        return result
```

---

## 总结

### 调用时机总结

| 场景 | 调用时机 | 优先级 |
|------|---------|--------|
| **配图** | 文本生成后，需要视觉辅助 | 🟡 中 |
| **TTS音频** | 儿童/视障/语言学习场景 | 🟢 高 |
| **PPT导出** | 用户主动请求导出 | 🟢 高 |
| **PDF导出** | 用户主动请求导出 | 🟡 中 |
| **演示视频** | 实验/操作步骤演示 | 🟡 中 |
| **交互内容** | 需要用户交互体验 | 🔵 低 |

### 核心原则

1. **按需调用**: 不是每次都调用，根据用户需求决定
2. **异步生成**: 多模态内容生成较慢，采用异步
3. **可选增强**: 文本内容是核心，多模态是增强
4. **用户驱动**: 很多时候由用户主动触发（导出按钮）

### 记忆口诀

```
文本内容是基础，多模态是锦上添花
配图配音看需求，导出格式看用户
实验操作要视频，儿童教育要音频
不是每次都调用，按需选择最优化
```

---

**文档版本**: 1.0  
**最后更新**: 2025-12-10  
**维护者**: MetaWorkflow Team
