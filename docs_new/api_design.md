# 通用领域内容生成·元工作流平台 - API接口设计 ⭐ V2.0

**文档版本**: V2.0 ⭐ **重大升级：多领域支持 + AI能力增强**  
**日期**: 2025-12-09  
**Base URL**: `https://api.metaworkflow.ai`  
**API版本**: `v2`

---

## 🆕 V2.0 变更概览

**核心升级**：
- ✅ **领域适配器API**：支持K12、美术史、畅销书、职业培训等多领域
- ✅ **AI目标生成API**：智能分析用户需求，生成结构化学习目标
- ✅ **AI内容发现API**：向量检索+多源聚合，智能发现相关内容
- ✅ **AI叙述生成API**：将枯燥内容转化为引人入胜的故事化叙述
- ✅ **多源检索API**：统一接口访问Milvus、Redis、MinIO、PostgreSQL

**技术栈增强**：
- 向量数据库：Milvus（支持亿级向量检索）
- 缓存层：Redis（高性能缓存+分布式锁）
- 对象存储：MinIO（多媒体资源存储）
- 监控：Prometheus + Grafana + Sentry

---

## 目录

### 核心API
1. [API概览](#1-api概览)
2. [认证与授权](#2-认证与授权)
3. [工作流管理API](#3-工作流管理api)
4. [工作流模板管理API](#4-工作流模板管理api)
5. [内容生成API](#5-内容生成api)
6. [课程规划API](#6-课程规划api) ⭐V1新增
7. [模板本地化API](#7-模板本地化api) ⭐V1新增
8. [组件管理API](#8-组件管理api) ⭐V1新增

### V2.0新增API ⭐
9. [领域适配器管理API](#9-领域适配器管理api) 🆕
10. [AI目标生成API](#10-ai目标生成api) 🆕
11. [AI内容发现API](#11-ai内容发现api) 🆕
12. [AI叙述生成API](#12-ai叙述生成api) 🆕
13. [多源内容检索API](#13-多源内容检索api) 🆕

### 通用API
14. [配置管理API](#14-配置管理api)
15. [监控与统计API](#15-监控与统计api)
16. [Webhook回调](#16-webhook回调)
17. [错误码参考](#17-错误码参考)
18. [SDK示例](#18-sdk示例)

---

## 1. API概览

### 1.1 设计原则

- **RESTful风格**: 使用标准HTTP方法(GET/POST/PUT/DELETE)
- **资源导向**: URL代表资源,动词通过HTTP方法表达
- **统一响应格式**: 所有API返回统一的JSON格式
- **版本控制**: 通过URL路径进行版本控制(`/v1/`)
- **幂等性**: PUT和DELETE操作保证幂等性

### 1.2 通用响应格式

**成功响应**:
```json
{
    "success": true,
    "data": {
        // 实际数据
    },
    "metadata": {
        "timestamp": "2025-12-09T10:00:00Z",
        "request_id": "req_abc123"
    }
}
```

**错误响应**:
```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid grade value",
        "details": {
            "field": "grade",
            "constraint": "must be between 1 and 12"
        }
    },
    "metadata": {
        "timestamp": "2025-12-09T10:00:00Z",
        "request_id": "req_abc123"
    }
}
```

### 1.3 HTTP状态码

| 状态码 | 含义 | 使用场景 |
|--------|------|---------|
| 200 | OK | 成功获取资源 |
| 201 | Created | 成功创建资源 |
| 202 | Accepted | 请求已接受,异步处理 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 429 | Too Many Requests | 限流 |
| 500 | Internal Server Error | 服务器错误 |

---

## 2. 认证与授权

### 2.1 获取访问令牌

**接口**: `POST /v1/auth/token`

**请求体**:
```json
{
    "grant_type": "password",
    "username": "user@example.com",
    "password": "secret"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "token_type": "Bearer",
        "expires_in": 3600,
        "refresh_token": "refresh_abc123"
    }
}
```

### 2.2 刷新令牌

**接口**: `POST /v1/auth/refresh`

**请求体**:
```json
{
    "refresh_token": "refresh_abc123"
}
```

### 2.3 使用令牌

所有受保护的API请求需在Header中携带:
```
Authorization: Bearer {access_token}
```

---

## 3. 工作流管理API

### 3.1 创建工作流实例

**接口**: `POST /v1/workflows/instances`

**描述**: 创建并执行一个新的工作流实例

**请求体**:
```json
{
    "subject": "数学",
    "grade": 5,
    "region": "CN",
    "topics": ["分数加减法", "通分"],
    "lesson_type": "新授课",
    "duration": 45,
    "pedagogy_model": "探究式",
    "options": {
        "enable_story": true,
        "interaction_level": "high",
        "output_format": "json"
    }
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "instance_id": "wf_20251209_001",
        "status": "pending",
        "estimated_time": 180,
        "created_at": "2025-12-09T10:00:00Z"
    },
    "metadata": {
        "request_id": "req_abc123"
    }
}
```

### 3.2 查询工作流实例状态

**接口**: `GET /v1/workflows/instances/{instance_id}`

**路径参数**:
- `instance_id`: 工作流实例ID

**响应**:
```json
{
    "success": true,
    "data": {
        "instance_id": "wf_20251209_001",
        "workflow_id": "math_lesson_basic",
        "status": "running",
        "progress": 0.6,
        "current_node": "content_generation",
        "completed_nodes": [
            "persona_analysis",
            "objectives_setting",
            "rag_retrieval"
        ],
        "created_at": "2025-12-09T10:00:00Z",
        "started_at": "2025-12-09T10:00:05Z",
        "estimated_completion": "2025-12-09T10:03:00Z"
    }
}
```

### 3.3 获取工作流结果

**接口**: `GET /v1/workflows/instances/{instance_id}/result`

**前置条件**: 工作流状态为`completed`

**响应**:
```json
{
    "success": true,
    "data": {
        "instance_id": "wf_20251209_001",
        "status": "completed",
        "result": {
            "lesson_id": "MATH_G5_FRACTION_001",
            "title": "分数加减法 - 魔法配方大挑战",
            "metadata": {
                "subject": "数学",
                "grade": 5,
                "duration": 45,
                "generated_at": "2025-12-09T10:02:30Z"
            },
            "sections": [
                {
                    "section_id": "intro",
                    "title": "魔法学院的挑战",
                    "duration": 5,
                    "narrative": "...",
                    "teacher_script": "...",
                    "interactions": [...]
                },
                // ... 更多环节
            ],
            "materials_list": ["PPT", "练习册", "魔法卡片"],
            "references": [
                {"source": "人教版数学五年级上册", "page": 42}
            ]
        },
        "metrics": {
            "total_duration": 175,
            "tokens_used": 15000,
            "cost": 0.45,
            "models_used": ["gpt-4", "deepseek-v2"]
        }
    }
}
```

### 3.4 取消工作流实例

**接口**: `POST /v1/workflows/instances/{instance_id}/cancel`

**响应**:
```json
{
    "success": true,
    "data": {
        "instance_id": "wf_20251209_001",
        "status": "cancelled",
        "message": "Workflow cancelled successfully"
    }
}
```

### 3.5 批量创建工作流

**接口**: `POST /v1/workflows/instances/batch`

**请求体**:
```json
{
    "items": [
        {
            "subject": "数学",
            "grade": 5,
            "topics": ["分数加法"]
        },
        {
            "subject": "语文",
            "grade": 5,
            "topics": ["记叙文写作"]
        }
    ],
    "common_config": {
        "region": "CN",
        "duration": 45
    }
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "batch_id": "batch_001",
        "total": 2,
        "instances": [
            {
                "instance_id": "wf_20251209_001",
                "status": "pending"
            },
            {
                "instance_id": "wf_20251209_002",
                "status": "pending"
            }
        ]
    }
}
```

### 3.6 查询工作流列表

**接口**: `GET /v1/workflows/instances`

**查询参数**:
- `status`: 状态过滤(pending/running/completed/failed)
- `subject`: 学科过滤
- `grade`: 年级过滤
- `from_date`: 开始日期
- `to_date`: 结束日期
- `page`: 页码(默认1)
- `page_size`: 每页数量(默认20)

**示例**: `GET /v1/workflows/instances?status=completed&subject=数学&page=1&page_size=20`

**响应**:
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "instance_id": "wf_20251209_001",
                "subject": "数学",
                "grade": 5,
                "status": "completed",
                "created_at": "2025-12-09T10:00:00Z"
            }
        ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 150,
            "total_pages": 8
        }
    }
}
```

---

## 4. 工作流模板管理API ⭐ 新增

### 4.1 创建模板

**接口**: `POST /v1/templates`

**描述**: 从工作流实例创建可复用模板

**请求体**:
```json
{
    "source_instance_id": "inst_abc123",
    "template_name": "小学低年级数学基础课件",
    "description": "适用于1-3年级数学新授课,包含故事化引入和游戏化互动",
    "tags": ["数学", "低年级", "故事化"],
    "is_public": true
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "template_name": "小学低年级数学基础课件",
        "version": "1.0",
        "metadata": {
            "subject": "数学",
            "grade_range": [1, 3],
            "region": "中国",
            "created_by": "user_123",
            "created_at": "2025-01-20T10:00:00Z",
            "description": "适用于1-3年级数学新授课,包含故事化引入和游戏化互动",
            "source_instance_id": "inst_abc123",
            "tags": ["数学", "低年级", "故事化"],
            "is_public": true
        },
        "nodes": [...],
        "edges": [...]
    }
}
```

### 4.2 查找模板

**接口**: `GET /v1/templates/{template_id}`

**路径参数**:
- `template_id`: 模板ID (必填)

**查询参数**:
- `version`: 版本号 (可选, 默认最新版本)

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "template_name": "小学低年级数学基础课件",
        "version": "1.2",
        "base_template_id": null,
        "metadata": {...},
        "nodes": [...],
        "edges": [...],
        "parameters": {
            "required": ["topics", "duration"],
            "optional": ["textbook_version", "enable_story"],
            "defaults": {
                "duration": 45,
                "enable_story": true
            }
        },
        "execution_history": {
            "total_runs": 152,
            "success_rate": 0.96,
            "avg_duration": 120.5,
            "avg_cost": 0.35,
            "avg_quality_score": 8.2,
            "last_run_at": "2025-01-20T14:30:00Z"
        },
        "ratings": {
            "avg_score": 4.5,
            "total_ratings": 42
        }
    }
}
```

### 4.3 按名称查找模板

**接口**: `GET /v1/templates/by-name/{template_name}`

**路径参数**:
- `template_name`: 模板名称 (必填)

**查询参数**:
- `version`: 版本号 (可选, 默认最新版本)

**响应**: 同 4.2

### 4.4 搜索模板

**接口**: `GET /v1/templates`

**查询参数**:
- `subject`: 学科筛选 (可选)
- `grade`: 年级筛选 (可选)
- `region`: 地区筛选 (可选)
- `tags`: 标签筛选,逗号分隔 (可选)
- `is_public`: 是否公开 (可选, true/false)
- `sort_by`: 排序方式 (可选, `popularity`/`rating`/`recent`, 默认`popularity`)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认20, 最大100)

**请求示例**:
```
GET /v1/templates?subject=数学&grade=3&sort_by=rating&page=1&page_size=20
```

**响应**:
```json
{
    "success": true,
    "data": {
        "templates": [
            {
                "template_id": "tpl_xyz789",
                "template_name": "小学低年级数学基础课件",
                "version": "1.2",
                "metadata": {
                    "subject": "数学",
                    "grade_range": [1, 3],
                    "description": "适用于1-3年级数学新授课",
                    "tags": ["数学", "低年级", "故事化"],
                    "popularity_score": 8.5
                },
                "execution_history": {
                    "total_runs": 152,
                    "success_rate": 0.96
                },
                "ratings": {
                    "avg_score": 4.5,
                    "total_ratings": 42
                }
            },
            // ... 更多模板
        ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 156,
            "total_pages": 8
        }
    }
}
```

### 4.5 实例化模板

**接口**: `POST /v1/templates/{template_id}/instantiate`

**路径参数**:
- `template_id`: 模板ID (必填)

**请求体**:
```json
{
    "parameters": {
        "topics": ["分数加法", "分数减法"],
        "duration": 45,
        "textbook_version": "人教版",
        "enable_story": true,
        "interaction_level": "high"
    },
    "execution_mode": "async"  // "async" | "sync"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "instance_id": "inst_new123",
        "template_id": "tpl_xyz789",
        "template_version": "1.2",
        "status": "running",
        "created_at": "2025-01-20T15:00:00Z",
        "estimated_completion": "2025-01-20T15:02:30Z"
    }
}
```

### 4.6 创建模板新版本

**接口**: `POST /v1/templates/{template_id}/versions`

**路径参数**:
- `template_id`: 模板ID (必填)

**请求体**:
```json
{
    "modifications": {
        "nodes": [
            {
                "action": "update",
                "node_id": "story_design",
                "updates": {
                    "config": {
                        "prompt_template": "updated_prompt_v2"
                    }
                }
            },
            {
                "action": "add",
                "node": {
                    "node_id": "interactive_quiz",
                    "node_type": "llm_generation",
                    "config": {...}
                }
            }
        ],
        "edges": [
            {
                "action": "add",
                "edge": {
                    "from": "content_generation",
                    "to": "interactive_quiz"
                }
            }
        ]
    },
    "change_log": "优化故事设计提示词，增加互动测验节点",
    "version_type": "minor"  // "major" | "minor" | "patch"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "new_version": "1.3",
        "parent_version": "1.2",
        "changes": "优化故事设计提示词，增加互动测验节点",
        "created_at": "2025-01-20T16:00:00Z"
    }
}
```

### 4.7 版本对比

**接口**: `GET /v1/templates/{template_id}/versions/compare`

**路径参数**:
- `template_id`: 模板ID (必填)

**查询参数**:
- `version1`: 版本1 (必填)
- `version2`: 版本2 (必填)

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "version1": "1.2",
        "version2": "1.3",
        "differences": {
            "nodes": {
                "added": [
                    {
                        "node_id": "interactive_quiz",
                        "node_type": "llm_generation"
                    }
                ],
                "removed": [],
                "modified": [
                    {
                        "node_id": "story_design",
                        "changes": {
                            "config.prompt_template": {
                                "old": "updated_prompt_v1",
                                "new": "updated_prompt_v2"
                            }
                        }
                    }
                ]
            },
            "edges": {
                "added": [
                    {"from": "content_generation", "to": "interactive_quiz"}
                ],
                "removed": []
            }
        }
    }
}
```

### 4.8 回滚版本

**接口**: `POST /v1/templates/{template_id}/rollback`

**路径参数**:
- `template_id`: 模板ID (必填)

**请求体**:
```json
{
    "target_version": "1.2",
    "reason": "新版本存在性能问题"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "new_version": "1.4",
        "rollback_to": "1.2",
        "created_at": "2025-01-20T17:00:00Z"
    }
}
```

### 4.9 获取优化建议

**接口**: `GET /v1/templates/{template_id}/suggestions`

**路径参数**:
- `template_id`: 模板ID (必填)

**查询参数**:
- `time_range_days`: 分析时间范围(天数) (可选, 默认30)
- `severity`: 严重程度筛选 (可选, `high`/`medium`/`low`)
- `type`: 建议类型筛选 (可选, `performance`/`quality`/`cost`/`structure`)

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "suggestions": [
            {
                "suggestion_id": "sug_001",
                "type": "performance",
                "severity": "medium",
                "issue": "节点 'content_generation' 平均耗时32.5秒，超过预期",
                "affected_nodes": ["content_generation"],
                "suggestion": "建议优化提示词长度或增加缓存",
                "expected_improvement": {
                    "performance": {"duration_reduction": "30%"}
                },
                "supporting_data": {
                    "avg_duration": 32.5,
                    "max_duration": 45.2,
                    "sample_count": 152
                },
                "status": "pending",
                "created_at": "2025-01-20T10:00:00Z"
            },
            {
                "suggestion_id": "sug_002",
                "type": "quality",
                "severity": "high",
                "issue": "内容质量平均分0.68，低于预期",
                "affected_nodes": ["story_design", "content_generation"],
                "suggestion": "建议优化Prompt模板或增加few-shot示例",
                "expected_improvement": {
                    "quality": {"score_improvement": "+0.3"}
                },
                "status": "pending"
            }
        ],
        "summary": {
            "total": 2,
            "by_severity": {
                "high": 1,
                "medium": 1,
                "low": 0
            },
            "by_type": {
                "performance": 1,
                "quality": 1
            }
        }
    }
}
```

### 4.10 应用优化建议

**接口**: `POST /v1/templates/{template_id}/suggestions/{suggestion_id}/apply`

**路径参数**:
- `template_id`: 模板ID (必填)
- `suggestion_id`: 建议ID (必填)

**请求体**:
```json
{
    "auto_create_version": true,  // 是否自动创建新版本
    "change_log": "应用性能优化建议 #sug_001"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "suggestion_id": "sug_001",
        "new_version": "1.5",
        "status": "applied",
        "applied_at": "2025-01-20T18:00:00Z"
    }
}
```

### 4.11 Fork模板

**接口**: `POST /v1/templates/{template_id}/fork`

**路径参数**:
- `template_id`: 模板ID (必填)

**请求体**:
```json
{
    "new_template_name": "我的数学课件模板",
    "description": "基于基础模板定制的个人版本",
    "modifications": {
        // 可选的初始修改
    }
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_new456",
        "template_name": "我的数学课件模板",
        "version": "1.0",
        "base_template_id": "tpl_xyz789",
        "created_at": "2025-01-20T19:00:00Z"
    }
}
```

### 4.12 获取模板版本列表

**接口**: `GET /v1/templates/{template_id}/versions`

**路径参数**:
- `template_id`: 模板ID (必填)

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "versions": [
            {
                "version": "1.5",
                "created_at": "2025-01-20T18:00:00Z",
                "created_by": "user_123",
                "changes": "应用性能优化建议 #sug_001",
                "is_latest": true
            },
            {
                "version": "1.4",
                "created_at": "2025-01-20T17:00:00Z",
                "created_by": "user_123",
                "changes": "Rollback to version 1.2",
                "is_latest": false
            },
            // ... 更多版本
        ]
    }
}
```

### 4.13 删除模板

**接口**: `DELETE /v1/templates/{template_id}`

**路径参数**:
- `template_id`: 模板ID (必填)

**查询参数**:
- `delete_all_versions`: 是否删除所有版本 (可选, 默认false, 仅删除指定版本)
- `version`: 要删除的版本号 (可选, 默认最新版本)

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "tpl_xyz789",
        "deleted_versions": ["1.5"],
        "deleted_at": "2025-01-20T20:00:00Z"
    }
}
```

---

## 5. 内容生成API

### 5.1 直接生成教案(快捷接口)

**接口**: `POST /v1/content/lesson-plan`

**描述**: 一步生成完整教案(内部调用工作流)

**请求体**:
```json
{
    "subject": "数学",
    "grade": 5,
    "topic": "分数加减法",
    "region": "CN",
    "options": {
        "enable_story": true,
        "output_format": "json"
    }
}
```

**响应** (同工作流结果):
```json
{
    "success": true,
    "data": {
        "lesson_id": "...",
        "title": "...",
        // ... 完整教案数据
    }
}
```

### 5.2 生成用户画像

**接口**: `POST /v1/content/persona`

**请求体**:
```json
{
    "grade": 5,
    "subject": "数学",
    "region": "CN",
    "additional_context": "学生对游戏化学习感兴趣"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "persona_id": "GRADE5_MATH_CN_001",
        "target_group": "五年级数学学生(10-11岁)",
        "cognitive": {
            "abstract_thinking": "transitioning",
            "attention_span_minutes": 15,
            "prerequisite_mastery": 75
        },
        "emotional": {
            "subject_interest": "medium",
            "motivation_type": "mixed"
        },
        "behavioral": {
            "learning_styles": ["visual", "kinesthetic"],
            "interaction_preference": "group_discussion"
        },
        "recommendations": {
            "narrative_style": "游戏化/情境化",
            "pacing": "moderate"
        }
    }
}
```

### 5.3 知识库检索

**接口**: `POST /v1/content/rag/search`

**请求体**:
```json
{
    "query": "分数加法的定义和计算方法",
    "knowledge_base": "cn_math_textbook",
    "top_k": 5,
    "filters": {
        "grade": 5,
        "chapter": "分数"
    }
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "results": [
            {
                "doc_id": "doc_123",
                "content": "同分母分数相加减,分母不变,分子相加减...",
                "score": 0.95,
                "metadata": {
                    "source": "人教版数学五年级上册",
                    "page": 42,
                    "chapter": "分数的加法和减法"
                }
            }
        ],
        "total": 5
    }
}
```

### 5.4 导出内容

**接口**: `GET /v1/content/{content_id}/export`

**查询参数**:
- `format`: 导出格式(json/pdf/pptx/docx)

**示例**: `GET /v1/content/MATH_G5_FRACTION_001/export?format=pdf`

**响应**:
- Content-Type: application/pdf
- Content-Disposition: attachment; filename="lesson_plan.pdf"

---

## 6. 课程规划API ⭐新增

课程规划API用于管理整个学期或学年的课程体系，包括单元划分、课时安排、故事宇宙设计等。

### 6.1 创建课程

**接口**: `POST /v1/curriculum/create`

**请求体**:
```json
{
    "name": "九年级化学上册",
    "subject": "化学",
    "grade": 9,
    "region": "CN",
    "textbook_version": "人教版",
    
    "start_date": "2025-09-01",
    "end_date": "2026-01-20",
    "total_weeks": 20,
    "lessons_per_week": 3,
    
    "enable_story_universe": true,
    "narrative_structure": "mystery_solving",
    
    "core_competencies": [
        "科学探究",
        "证据推理",
        "创新意识"
    ]
}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "curriculum_id": "curr_chem_9_2025_fall",
        "name": "九年级化学上册",
        "subject": "化学",
        "grade": 9,
        "region": "CN",
        "status": "DRAFT",
        "total_lessons": 60,
        "created_at": "2025-12-09T10:00:00Z"
    }
}
```

---

### 6.2 自动规划课程单元

**接口**: `POST /v1/curriculum/{curriculum_id}/plan-units`

**描述**: 基于教材内容自动划分单元、分配课时、设计教学序列。

**请求体**:
```json
{
    "textbook_content": "...",  // 教材全文或章节目录
    "planning_strategy": {
        "unit_size": "medium",  // small/medium/large
        "include_review_lessons": true,
        "include_assessment_lessons": true,
        "story_integration": true
    }
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "job_id": "plan_job_abc123",
        "status": "processing",
        "estimated_duration": 120,  // 秒
        "status_url": "/v1/curriculum/planning-jobs/plan_job_abc123"
    }
}
```

---

### 6.3 获取规划结果

**接口**: `GET /v1/curriculum/planning-jobs/{job_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "job_id": "plan_job_abc123",
        "status": "completed",
        "curriculum_id": "curr_chem_9_2025_fall",
        
        "result": {
            "units": [
                {
                    "unit_id": "unit_1",
                    "sequence_number": 1,
                    "title": "走进化学世界",
                    "topics": ["物质的变化", "化学实验基础"],
                    "total_lessons": 6,
                    "lesson_distribution": {
                        "introduction": 1,
                        "new_content": 3,
                        "practice": 1,
                        "review": 1
                    },
                    "estimated_weeks": 2
                },
                {
                    "unit_id": "unit_2",
                    "sequence_number": 2,
                    "title": "我们周围的空气",
                    "topics": ["氧气", "空气的成分"],
                    "total_lessons": 8,
                    "estimated_weeks": 3
                }
                // ... 更多单元
            ],
            "total_units": 5,
            "total_lessons": 60,
            
            "knowledge_graph": {
                "nodes_count": 45,
                "edges_count": 78,
                "max_depth": 4
            },
            
            "story_universe": {
                "universe_id": "story_chem_detective",
                "name": "化学侦探学院",
                "theme": "通过破案学习化学知识",
                "main_characters": [
                    {
                        "character_id": "char_001",
                        "name": "李小智",
                        "role": "protagonist"
                    }
                ]
            }
        }
    }
}
```

---

### 6.4 获取课程详情

**接口**: `GET /v1/curriculum/{curriculum_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "curriculum_id": "curr_chem_9_2025_fall",
        "name": "九年级化学上册",
        "subject": "化学",
        "grade": 9,
        "region": "CN",
        "status": "APPROVED",
        
        "time_info": {
            "start_date": "2025-09-01",
            "end_date": "2026-01-20",
            "total_weeks": 20,
            "lessons_per_week": 3,
            "total_lessons": 60
        },
        
        "structure": {
            "total_units": 5,
            "units": [
                {
                    "unit_id": "unit_1",
                    "title": "走进化学世界",
                    "sequence_number": 1,
                    "total_lessons": 6,
                    "status": "completed"
                }
                // ... 更多单元
            ]
        },
        
        "progress": {
            "completed_lessons": 15,
            "total_lessons": 60,
            "completion_rate": 0.25,
            "current_unit": "unit_2",
            "on_schedule": true
        },
        
        "story_universe_id": "story_chem_detective",
        "knowledge_graph_id": "kg_chem_9",
        
        "created_at": "2025-12-09T10:00:00Z",
        "updated_at": "2025-12-09T15:30:00Z"
    }
}
```

---

### 6.5 获取单元详情

**接口**: `GET /v1/curriculum/units/{unit_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "unit_id": "unit_1",
        "curriculum_id": "curr_chem_9_2025_fall",
        "sequence_number": 1,
        "title": "走进化学世界",
        "subtitle": "化学让生活更美好",
        
        "content": {
            "topics": ["物质的变化", "化学实验基础", "走进化学实验室"],
            "learning_objectives": [
                "认识化学是一门以实验为基础的科学",
                "初步学会药品的取用等基本实验操作"
            ],
            "driving_questions": [
                "化学与我们的生活有什么关系？",
                "如何安全地进行化学实验？"
            ]
        },
        
        "lessons": [
            {
                "lesson_meta_id": "lesson_001",
                "sequence_number": 1,
                "global_sequence": 1,
                "title": "物质的变化和性质",
                "lesson_type": "新授课",
                "duration": 45,
                "scheduled_week": 1,
                "status": "generated",
                "lesson_script_id": "script_001"
            }
            // ... 更多课时
        ],
        
        "story_arc": {
            "arc_id": "arc_001",
            "arc_title": "侦探学院新生入学",
            "arc_type": "adventure",
            "setup": "李小智收到神秘的化学侦探学院录取通知书...",
            "episodes_count": 6
        },
        
        "time_allocation": {
            "estimated_weeks": 2,
            "start_week": 1,
            "end_week": 2
        }
    }
}
```

---

### 6.6 获取课时元信息

**接口**: `GET /v1/curriculum/lessons/{lesson_meta_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "lesson_meta_id": "lesson_001",
        "unit_id": "unit_1",
        "curriculum_id": "curr_chem_9_2025_fall",
        
        "sequence_number": 1,
        "global_sequence": 1,
        
        "basic_info": {
            "title": "物质的变化和性质",
            "lesson_type": "新授课",
            "duration": 45,
            "scheduled_week": 1,
            "scheduled_date": "2025-09-03"
        },
        
        "content": {
            "topics": ["物理变化", "化学变化", "物质性质"],
            "objectives": [
                "能够区分物理变化和化学变化",
                "初步学会观察和描述化学现象"
            ],
            "key_points": ["化学变化的特征"],
            "difficult_points": ["区分物理变化和化学变化"]
        },
        
        "story_element": {
            "story_episode_id": "episode_001",
            "narrative_hook": "李小智的第一个案件：神秘消失的蜡烛"
        },
        
        "generation_config": {
            "workflow_template_id": "template_inquiry_lesson",
            "pedagogy_model": "inquiry",
            "interaction_level": "high",
            "enable_story": true
        },
        
        "generation_status": {
            "generated": true,
            "lesson_script_id": "script_001",
            "generation_quality_score": 8.7,
            "generated_at": "2025-12-09T14:00:00Z"
        },
        
        "status": "approved"
    }
}
```

---

### 6.7 批量生成课时内容

**接口**: `POST /v1/curriculum/{curriculum_id}/generate-lessons`

**描述**: 基于课程规划，批量生成所有或指定课时的完整教学内容。

**请求体**:
```json
{
    "scope": "unit",  // "all" | "unit" | "custom"
    "unit_ids": ["unit_1", "unit_2"],  // 当scope为unit时
    
    "generation_config": {
        "parallel_jobs": 5,  // 并发生成数量
        "auto_review": true,  // 自动质量审核
        "min_quality_score": 7.5
    }
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "batch_job_id": "batch_gen_xyz789",
        "total_lessons": 14,
        "estimated_duration": 1800,  // 秒
        "status_url": "/v1/curriculum/batch-jobs/batch_gen_xyz789"
    }
}
```

---

### 6.8 获取批量生成进度

**接口**: `GET /v1/curriculum/batch-jobs/{batch_job_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "batch_job_id": "batch_gen_xyz789",
        "status": "in_progress",
        "progress": {
            "total": 14,
            "completed": 8,
            "in_progress": 3,
            "pending": 3,
            "failed": 0,
            "completion_rate": 0.57
        },
        
        "results": [
            {
                "lesson_meta_id": "lesson_001",
                "status": "completed",
                "lesson_script_id": "script_001",
                "quality_score": 8.7
            },
            {
                "lesson_meta_id": "lesson_002",
                "status": "in_progress",
                "progress": 0.65
            }
            // ... 更多
        ],
        
        "statistics": {
            "avg_quality_score": 8.5,
            "total_duration": 720,
            "avg_duration_per_lesson": 90
        },
        
        "started_at": "2025-12-09T15:00:00Z",
        "estimated_completion": "2025-12-09T15:30:00Z"
    }
}
```

---

### 6.9 获取故事宇宙

**接口**: `GET /v1/curriculum/{curriculum_id}/story-universe`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "universe_id": "story_chem_detective",
        "curriculum_id": "curr_chem_9_2025_fall",
        "name": "化学侦探学院",
        "description": "一所神秘的学院，学生通过破解化学谜题来学习化学知识",
        
        "world_setting": "现代都市背景，融入科幻元素",
        "theme": "科学探究与逻辑推理",
        "tone": "轻松幽默，充满悬念",
        
        "characters": [
            {
                "character_id": "char_001",
                "name": "李小智",
                "role": "protagonist",
                "description": "好奇心强的初三学生，擅长观察和推理",
                "personality": ["聪明", "勇敢", "有点粗心"],
                "avatar_url": "https://cdn.example.com/avatars/char_001.png"
            },
            {
                "character_id": "char_002",
                "name": "化学博士",
                "role": "mentor",
                "description": "学院的院长，神秘莫测的化学大师",
                "personality": ["智慧", "幽默", "神秘"]
            }
        ],
        
        "story_arcs": [
            {
                "arc_id": "arc_001",
                "unit_id": "unit_1",
                "arc_title": "侦探学院新生入学",
                "arc_type": "adventure",
                "setup": "李小智收到神秘录取通知书，来到化学侦探学院...",
                "conflict": "入学考试是破解一系列化学谜题",
                "climax": "最终谜题挑战：区分物理变化和化学变化",
                "resolution": "成功入学，开启化学探险之旅",
                "episodes_count": 6
            }
            // ... 更多故事弧
        ],
        
        "continuity_elements": [
            "学院积分系统",
            "持续的案件主线",
            "角色成长弧线"
        ],
        
        "visual_assets": {
            "world_map": "https://cdn.example.com/story/world_map.jpg",
            "character_designs": "https://cdn.example.com/story/characters.pdf"
        }
    }
}
```

---

### 6.10 获取知识图谱

**接口**: `GET /v1/curriculum/{curriculum_id}/knowledge-graph`

**查询参数**:
- `format`: `json` | `neo4j` | `graphml` (可选，默认json)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "graph_id": "kg_chem_9",
        "curriculum_id": "curr_chem_9_2025_fall",
        "name": "九年级化学知识图谱",
        
        "statistics": {
            "total_nodes": 45,
            "total_edges": 78,
            "max_depth": 4,
            "avg_degree": 3.5
        },
        
        "nodes": [
            {
                "node_id": "kn_001",
                "name": "物质的变化",
                "category": "concept",
                "difficulty_level": 2,
                "bloom_level": "understand",
                "related_lessons": ["lesson_001"],
                "estimated_learning_time": 45
            },
            {
                "node_id": "kn_002",
                "name": "物理变化",
                "category": "concept",
                "difficulty_level": 1,
                "bloom_level": "remember",
                "related_lessons": ["lesson_001"]
            },
            {
                "node_id": "kn_003",
                "name": "化学变化",
                "category": "concept",
                "difficulty_level": 2,
                "bloom_level": "understand",
                "related_lessons": ["lesson_001", "lesson_002"]
            }
            // ... 更多节点
        ],
        
        "edges": [
            {
                "edge_id": "edge_001",
                "source_node_id": "kn_001",
                "target_node_id": "kn_002",
                "relationship_type": "includes",
                "strength": 1.0,
                "description": "物质的变化包括物理变化"
            },
            {
                "edge_id": "edge_002",
                "source_node_id": "kn_002",
                "target_node_id": "kn_003",
                "relationship_type": "prerequisite",
                "strength": 0.8,
                "description": "理解物理变化是理解化学变化的基础"
            }
            // ... 更多边
        ],
        
        "levels": {
            "kn_001": 1,
            "kn_002": 2,
            "kn_003": 2,
            "kn_004": 3
        }
    }
}
```

---

### 6.11 获取学习路径

**接口**: `GET /v1/curriculum/knowledge-graph/{graph_id}/learning-path`

**查询参数**:
- `target_node_id`: 目标知识点ID (必需)
- `current_mastery`: 当前已掌握的知识点ID数组 (可选)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "target_node_id": "kn_015",
        "target_node_name": "化学方程式配平",
        
        "learning_path": [
            {
                "step": 1,
                "node_id": "kn_001",
                "node_name": "物质的变化",
                "reason": "基础概念",
                "related_lessons": ["lesson_001"]
            },
            {
                "step": 2,
                "node_id": "kn_005",
                "node_name": "化学反应",
                "reason": "理解化学反应是配平的前提",
                "related_lessons": ["lesson_003"]
            },
            {
                "step": 3,
                "node_id": "kn_010",
                "node_name": "质量守恒定律",
                "reason": "配平的理论依据",
                "related_lessons": ["lesson_008"]
            },
            {
                "step": 4,
                "node_id": "kn_015",
                "node_name": "化学方程式配平",
                "reason": "目标知识点",
                "related_lessons": ["lesson_012", "lesson_013"]
            }
        ],
        
        "total_steps": 4,
        "estimated_learning_time": 180,  // 分钟
        "difficulty_progression": [2, 3, 3, 4]
    }
}
```

---

### 6.12 更新课程进度

**接口**: `PUT /v1/curriculum/{curriculum_id}/progress`

**请求体**:
```json
{
    "completed_lessons": ["lesson_001", "lesson_002", "lesson_003"],
    "current_unit": "unit_2",
    "current_week": 3,
    "notes": "学生对化学变化概念掌握良好"
}
```

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "curriculum_id": "curr_chem_9_2025_fall",
        "progress": {
            "completed_lessons": 15,
            "total_lessons": 60,
            "completion_rate": 0.25,
            "current_unit": "unit_2",
            "current_week": 3,
            "on_schedule": true,
            "schedule_variance_days": 0
        },
        "updated_at": "2025-12-09T16:00:00Z"
    }
}
```

---

### 6.13 列出所有课程

**接口**: `GET /v1/curriculum/list`

**查询参数**:
- `subject`: 学科 (可选)
- `grade`: 年级 (可选)
- `region`: 地区 (可选)
- `status`: 状态 (可选)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认20)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "curriculums": [
            {
                "curriculum_id": "curr_chem_9_2025_fall",
                "name": "九年级化学上册",
                "subject": "化学",
                "grade": 9,
                "region": "CN",
                "status": "IN_PROGRESS",
                "progress": {
                    "completion_rate": 0.25
                },
                "created_at": "2025-09-01T00:00:00Z"
            }
            // ... 更多课程
        ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total_items": 45,
            "total_pages": 3
        }
    }
}
```

---

## 7. 模板本地化API ⭐新增

模板本地化API用于将工作流模板适配到不同地区、语言和教育体系。

### 7.1 创建本地化任务

**接口**: `POST /v1/localization/create`

**请求体**:
```json
{
    "job_name": "九年级化学模板 中国→美国",
    "source_template_id": "template_chem_9_cn",
    "source_region": "CN",
    "source_language": "zh-CN",
    "target_region": "US",
    "target_language": "en-US",
    
    "localization_strategy": {
        "translation_model": "gpt-4",
        "preserve_structure": true,
        "adaptation_level": "high",
        "compliance_check": true,
        "human_review_required": true
    },
    
    "adaptation_scopes": [
        "standards",
        "culture",
        "pedagogy",
        "examples",
        "measurement_units"
    ]
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "job_name": "九年级化学模板 中国→美国",
        "status": "PENDING",
        "estimated_duration": 48,  // 小时
        "status_url": "/v1/localization/loc_job_abc123"
    }
}
```

---

### 7.2 获取本地化状态

**接口**: `GET /v1/localization/{localization_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "job_name": "九年级化学模板 中国→美国",
        "status": "ADAPTING",
        "current_phase": "cultural_adaptation",
        "progress": 0.65,
        
        "source": {
            "template_id": "template_chem_9_cn",
            "region": "CN",
            "language": "zh-CN"
        },
        
        "target": {
            "region": "US",
            "language": "en-US",
            "template_id": null  // 完成后生成
        },
        
        "phases_progress": {
            "standard_mapping": "completed",
            "translation": "completed",
            "cultural_adaptation": "in_progress",
            "pedagogy_adjustment": "pending",
            "compliance_check": "pending",
            "quality_review": "pending"
        },
        
        "statistics": {
            "standard_mappings_count": 12,
            "cultural_adaptations_count": 28,
            "compliance_issues_count": 3,
            "automation_rate": 0.68
        },
        
        "created_at": "2025-12-09T10:00:00Z",
        "estimated_completion": "2025-12-11T10:00:00Z"
    }
}
```

---

### 7.3 获取课标映射

**接口**: `GET /v1/localization/{localization_id}/standard-mappings`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "total_mappings": 12,
        
        "mappings": [
            {
                "mapping_id": "map_001",
                "source_standard": {
                    "id": "CN_CHEM_9_01",
                    "name": "认识物质的组成、结构和性质",
                    "description": "初步认识物质的微观构成...",
                    "region": "CN"
                },
                "target_standard": {
                    "id": "US_NGSS_MS_PS1_1",
                    "name": "Develop models to describe atomic composition",
                    "description": "Develop models to describe the atomic composition of simple molecules...",
                    "region": "US"
                },
                "mapping_type": "partial",
                "similarity_score": 0.75,
                "differences": [
                    "源课标更强调宏观现象，目标课标更强调建模能力"
                ],
                "adaptation_suggestions": [
                    "增加模型构建活动",
                    "强调科学实践过程"
                ],
                "confidence": 0.82,
                "verified_by_human": true
            }
            // ... 更多映射
        ]
    }
}
```

---

### 7.4 获取文化适配记录

**接口**: `GET /v1/localization/{localization_id}/cultural-adaptations`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "total_adaptations": 28,
        
        "adaptations": [
            {
                "category": "cultural_reference",
                "source_element": "春节放鞭炮",
                "target_element": "4th of July fireworks",
                "context": "在讲解化学反应时的例子",
                "reason": "使用目标文化中熟悉的节日场景",
                "confidence": 0.9
            },
            {
                "category": "measurement_unit",
                "source_element": "摄氏度 (°C)",
                "target_element": "华氏度 (°F)",
                "context": "温度相关实验",
                "reason": "美国主要使用华氏温标",
                "confidence": 1.0
            },
            {
                "category": "example",
                "source_element": "豆浆制作",
                "target_element": "Making lemonade",
                "context": "溶液概念示例",
                "reason": "更贴近美国学生日常经验",
                "confidence": 0.85
            }
            // ... 更多适配
        ],
        
        "categories_summary": {
            "cultural_reference": 12,
            "measurement_unit": 5,
            "example": 8,
            "name": 3
        }
    }
}
```

---

### 7.5 获取教学法调整

**接口**: `GET /v1/localization/{localization_id}/pedagogy-adjustments`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "total_adjustments": 15,
        
        "adjustments": [
            {
                "aspect": "classroom_structure",
                "source_value": "teacher-centered",
                "target_value": "student-centered",
                "reason": "美国课堂更强调学生主体性",
                "impact": "增加小组讨论和探究活动环节"
            },
            {
                "aspect": "assessment_method",
                "source_value": "mainly_written_tests",
                "target_value": "diverse_assessments",
                "reason": "美国教育强调多元评估",
                "impact": "增加项目式评估、口头报告等"
            },
            {
                "aspect": "time_allocation",
                "source_value": {
                    "lecture": 30,
                    "practice": 10,
                    "summary": 5
                },
                "target_value": {
                    "engagement": 10,
                    "exploration": 15,
                    "explanation": 10,
                    "elaboration": 8,
                    "evaluation": 2
                },
                "reason": "采用5E教学模式",
                "impact": "重新设计教学环节序列"
            }
            // ... 更多调整
        ]
    }
}
```

---

### 7.6 获取合规性问题

**接口**: `GET /v1/localization/{localization_id}/compliance-issues`

**查询参数**:
- `severity`: critical | high | medium | low (可选)
- `status`: open | fixed | wont_fix | false_positive (可选)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "total_issues": 3,
        "open_issues": 1,
        
        "issues": [
            {
                "issue_id": "issue_001",
                "category": "safety",
                "severity": "high",
                "title": "实验安全警告不足",
                "description": "美国要求更详细的实验室安全说明和MSDS引用",
                "location": {
                    "node_id": "script_generator",
                    "section": "实验步骤",
                    "line_number": 45
                },
                "flagged_content": "加热试管...",
                "violated_rule": "US Lab Safety Standards",
                "rule_reference": "https://www.osha.gov/...",
                "suggested_fix": "添加护目镜佩戴提示和紧急处理程序",
                "status": "open",
                "detected_at": "2025-12-09T14:30:00Z"
            },
            {
                "issue_id": "issue_002",
                "category": "cultural",
                "severity": "medium",
                "title": "可能引起文化敏感的内容",
                "description": "某些化学物质名称在目标文化中有特殊含义",
                "status": "fixed",
                "resolution": "已替换为中性表述",
                "resolved_at": "2025-12-09T15:00:00Z"
            }
            // ... 更多问题
        ],
        
        "summary_by_category": {
            "legal": 0,
            "cultural": 1,
            "educational": 0,
            "safety": 2,
            "privacy": 0
        },
        
        "summary_by_severity": {
            "critical": 0,
            "high": 1,
            "medium": 1,
            "low": 1
        }
    }
}
```

---

### 7.7 修复合规性问题

**接口**: `PUT /v1/localization/compliance-issues/{issue_id}/resolve`

**请求体**:
```json
{
    "resolution": "添加了详细的实验室安全警告和MSDS链接",
    "status": "fixed",
    "modified_content": "..."
}
```

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "issue_id": "issue_001",
        "status": "fixed",
        "resolved_by": "user_123",
        "resolved_at": "2025-12-09T16:00:00Z"
    }
}
```

---

### 7.8 获取质量报告

**接口**: `GET /v1/localization/{localization_id}/quality-report`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "report_id": "qr_abc123",
        "localization_id": "loc_job_abc123",
        "overall_score": 8.5,
        "grade": "B",
        
        "metrics": [
            {
                "metric_name": "translation_accuracy",
                "score": 9.2,
                "weight": 0.25,
                "details": "专业术语翻译准确，语言流畅"
            },
            {
                "metric_name": "cultural_appropriateness",
                "score": 8.5,
                "weight": 0.20,
                "details": "文化适配基本到位，少数细节需调整"
            },
            {
                "metric_name": "pedagogy_alignment",
                "score": 8.8,
                "weight": 0.20,
                "details": "教学法调整符合目标地区习惯"
            },
            {
                "metric_name": "standard_compliance",
                "score": 9.5,
                "weight": 0.20,
                "details": "课标映射准确，符合NGSS要求"
            },
            {
                "metric_name": "consistency",
                "score": 7.0,
                "weight": 0.15,
                "details": "术语一致性需要提升"
            }
        ],
        
        "issue_summary": {
            "critical": 0,
            "high": 1,
            "medium": 2,
            "low": 5
        },
        
        "strengths": [
            "课标映射准确全面",
            "文化适配考虑周到",
            "教学法调整合理"
        ],
        
        "weaknesses": [
            "部分术语翻译不一致",
            "少数安全提示不够详细"
        ],
        
        "recommendations": [
            "建立术语库，确保翻译一致性",
            "补充实验室安全说明",
            "增加更多本地化示例"
        ],
        
        "human_review_required": true,
        "human_review_points": [
            "专业术语最终审定",
            "文化敏感内容复核",
            "教学效果预测"
        ],
        
        "generated_at": "2025-12-09T16:30:00Z"
    }
}
```

---

### 7.9 启动人工审核

**接口**: `POST /v1/localization/{localization_id}/request-review`

**请求体**:
```json
{
    "review_type": "expert",  // "expert" | "native_speaker" | "educator"
    "review_scope": "full",  // "full" | "critical_only"
    "deadline": "2025-12-12T00:00:00Z",
    "notes": "请重点审核化学术语翻译和实验安全说明"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "review_id": "review_xyz789",
        "localization_id": "loc_job_abc123",
        "status": "pending_assignment",
        "estimated_completion": "2025-12-11T18:00:00Z"
    }
}
```

---

### 7.10 完成本地化并发布

**接口**: `POST /v1/localization/{localization_id}/publish`

**请求体**:
```json
{
    "target_template_name": "Chemistry Grade 9 - US Version",
    "publish_notes": "基于中国版本本地化，已通过质量审核",
    "tags": ["NGSS", "chemistry", "grade-9", "US"]
}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "localization_id": "loc_job_abc123",
        "status": "COMPLETED",
        "target_template_id": "template_chem_9_us",
        "target_template": {
            "template_id": "template_chem_9_us",
            "name": "Chemistry Grade 9 - US Version",
            "version": "1.0.0",
            "region": "US",
            "language": "en-US",
            "published_at": "2025-12-09T17:00:00Z"
        },
        "localization_summary": {
            "duration_hours": 36,
            "automation_rate": 0.68,
            "quality_score": 8.5,
            "total_adaptations": 45,
            "human_review_hours": 8
        }
    }
}
```

---

### 7.11 列出本地化任务

**接口**: `GET /v1/localization/list`

**查询参数**:
- `source_region`: 源地区 (可选)
- `target_region`: 目标地区 (可选)
- `status`: 状态 (可选)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认20)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "localizations": [
            {
                "localization_id": "loc_job_abc123",
                "job_name": "九年级化学模板 中国→美国",
                "source_region": "CN",
                "target_region": "US",
                "status": "COMPLETED",
                "progress": 1.0,
                "quality_score": 8.5,
                "created_at": "2025-12-09T10:00:00Z",
                "completed_at": "2025-12-10T22:00:00Z"
            }
            // ... 更多任务
        ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total_items": 15,
            "total_pages": 1
        }
    }
}
```

---

## 8. 组件管理API ⭐新增

组件管理API用于管理工作流节点组件，包括注册、安装、更新、使用统计等。

### 8.1 搜索组件

**接口**: `GET /v1/components/search`

**查询参数**:
- `keyword`: 关键词 (可选)
- `category`: 分类 (可选)
- `subject`: 适用学科 (可选)
- `min_rating`: 最低评分 (可选)
- `certified_only`: 仅显示认证组件 (可选，布尔值)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认20)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "components": [
            {
                "node_id": "virtual_lab_node_v1",
                "node_name": "虚拟实验室节点",
                "version": "1.2.0",
                "category": "interaction_design",
                "tags": ["实验", "虚拟", "交互", "3D"],
                "description": "提供化学、物理实验的虚拟仿真能力",
                "author": "张三",
                "organization": "某某科技",
                "certification_level": "gold",
                "verified": true,
                "avg_rating": 4.7,
                "download_count": 1523,
                "pricing_model": "freemium"
            }
            // ... 更多组件
        ],
        "facets": {
            "categories": {
                "content_generation": 145,
                "interaction_design": 67,
                "quality_assurance": 89
            },
            "certification_levels": {
                "platinum": 12,
                "gold": 45,
                "silver": 123
            }
        },
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total_items": 256,
            "total_pages": 13
        }
    }
}
```

---

### 8.2 获取组件详情

**接口**: `GET /v1/components/{node_id}`

**查询参数**:
- `version`: 指定版本 (可选，默认最新版本)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "node_id": "virtual_lab_node_v1",
        "node_name": "虚拟实验室节点",
        "version": "1.2.0",
        "category": "interaction_design",
        "subcategory": "simulation",
        "tags": ["实验", "虚拟", "交互", "3D"],
        
        "description": "提供化学、物理实验的虚拟仿真能力，支持Unity WebGL渲染",
        "long_description": "...",
        
        "applicable_subjects": ["化学", "物理"],
        "applicable_grades": "7-12",
        "applicable_lesson_types": ["实验课", "探究课"],
        
        "author_info": {
            "author": "张三",
            "author_email": "zhangsan@example.com",
            "organization": "某某科技",
            "license": "MIT"
        },
        
        "dependencies": [
            "llm_client >= 2.0.0",
            "unity_renderer >= 3.5.0",
            "model_library >= 1.0.0"
        ],
        
        "performance": {
            "estimated_duration": 18,  // 秒
            "memory_limit": 512,  // MB
            "cpu_limit": 2.0,  // 核
            "success_rate": 0.966
        },
        
        "input_schema": {
            "required": ["experiment_name", "experiment_type"],
            "optional": ["apparatus_list", "steps"]
        },
        
        "output_schema": {
            "guaranteed": ["virtual_lab_url", "metadata"],
            "conditional": {
                "video_url": "当Unity渲染失败时提供"
            }
        },
        
        "documentation_url": "https://docs.example.com/virtual_lab_node",
        "example_url": "https://examples.example.com/virtual_lab",
        
        "changelog": [
            {
                "version": "1.2.0",
                "date": "2025-12-09",
                "changes": "新增降级策略，支持视频fallback"
            },
            {
                "version": "1.1.0",
                "date": "2025-11-01",
                "changes": "性能优化，内存使用减少30%"
            }
        ],
        
        "statistics": {
            "download_count": 1523,
            "usage_count": 5847,
            "avg_rating": 4.7,
            "review_count": 89
        },
        
        "certification": {
            "certification_level": "gold",
            "verified": true,
            "certification_date": "2025-11-15"
        },
        
        "pricing": {
            "pricing_model": "freemium",
            "free_quota": 100,
            "price_per_call": 0.05
        },
        
        "created_at": "2025-10-01T00:00:00Z",
        "updated_at": "2025-12-09T10:00:00Z"
    }
}
```

---

### 8.3 安装组件

**接口**: `POST /v1/components/install`

**请求体**:
```json
{
    "node_id": "virtual_lab_node_v1",
    "version": "1.2.0",  // 可选，默认最新版本
    "workspace_id": "workspace_abc"  // 工作空间ID
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "installation_id": "install_xyz789",
        "node_id": "virtual_lab_node_v1",
        "version": "1.2.0",
        "status": "installing",
        "progress": 0.0,
        "status_url": "/v1/components/installations/install_xyz789"
    }
}
```

---

### 8.4 获取安装状态

**接口**: `GET /v1/components/installations/{installation_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "installation_id": "install_xyz789",
        "node_id": "virtual_lab_node_v1",
        "version": "1.2.0",
        "status": "completed",
        "progress": 1.0,
        
        "steps": [
            {
                "step": "download",
                "status": "completed",
                "duration": 5.2
            },
            {
                "step": "verify",
                "status": "completed",
                "duration": 1.5
            },
            {
                "step": "install_dependencies",
                "status": "completed",
                "duration": 12.8
            },
            {
                "step": "register",
                "status": "completed",
                "duration": 0.3
            }
        ],
        
        "total_duration": 19.8,
        "installed_at": "2025-12-09T17:30:00Z"
    }
}
```

---

### 8.5 卸载组件

**接口**: `DELETE /v1/components/{node_id}`

**查询参数**:
- `workspace_id`: 工作空间ID (必需)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "node_id": "virtual_lab_node_v1",
        "uninstalled_at": "2025-12-09T18:00:00Z"
    }
}
```

---

### 8.6 获取已安装组件列表

**接口**: `GET /v1/components/installed`

**查询参数**:
- `workspace_id`: 工作空间ID (必需)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "workspace_id": "workspace_abc",
        "installed_components": [
            {
                "node_id": "virtual_lab_node_v1",
                "node_name": "虚拟实验室节点",
                "version": "1.2.0",
                "installed_at": "2025-12-09T17:30:00Z",
                "usage_count": 45,
                "last_used_at": "2025-12-09T18:00:00Z"
            }
            // ... 更多组件
        ],
        "total_count": 12
    }
}
```

---

### 8.7 更新组件

**接口**: `PUT /v1/components/{node_id}/update`

**请求体**:
```json
{
    "target_version": "1.3.0",  // 可选，默认最新版本
    "workspace_id": "workspace_abc"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "update_id": "update_qwe456",
        "node_id": "virtual_lab_node_v1",
        "from_version": "1.2.0",
        "to_version": "1.3.0",
        "status": "updating",
        "status_url": "/v1/components/updates/update_qwe456"
    }
}
```

---

### 8.8 获取组件使用统计

**接口**: `GET /v1/components/{node_id}/statistics`

**查询参数**:
- `time_range`: 时间范围，如 "7d", "30d", "90d" (可选，默认30d)

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "node_id": "virtual_lab_node_v1",
        "time_range": "30d",
        
        "usage_statistics": {
            "total_executions": 5847,
            "successful_executions": 5650,
            "failed_executions": 197,
            "success_rate": 0.966,
            
            "avg_duration": 18.2,  // 秒
            "p50_duration": 16.5,
            "p95_duration": 28.3,
            "p99_duration": 35.1,
            
            "avg_memory_usage": 387,  // MB
            "peak_memory_usage": 498,
            
            "fallback_activations": 198,
            "fallback_rate": 0.034
        },
        
        "quality_metrics": {
            "avg_output_quality_score": 8.7,
            "avg_user_rating": 4.7,
            "user_ratings_count": 89
        },
        
        "trend": {
            "usage_trend": "increasing",
            "success_rate_trend": "stable",
            "quality_trend": "improving"
        },
        
        "top_users": [
            {
                "user_id": "user_123",
                "usage_count": 234
            }
            // ... 更多
        ]
    }
}
```

---

### 8.9 评价组件

**接口**: `POST /v1/components/{node_id}/reviews`

**请求体**:
```json
{
    "rating": 5,  // 1-5星
    "title": "非常好用的虚拟实验室",
    "comment": "渲染效果很棒，降级策略也很贴心",
    "pros": ["功能强大", "稳定性好", "文档完善"],
    "cons": ["内存占用稍高"]
}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "review_id": "review_aaa111",
        "node_id": "virtual_lab_node_v1",
        "rating": 5,
        "created_at": "2025-12-09T19:00:00Z"
    }
}
```

---

### 8.10 发布组件（开发者）

**接口**: `POST /v1/components/publish`

**描述**: 开发者发布新组件到组件市场。

**请求体** (multipart/form-data):
```
node_metadata: {JSON格式的NodeMetadata}
package_file: {组件包ZIP文件}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "publish_id": "pub_bbb222",
        "node_id": "my_custom_node_v1",
        "version": "1.0.0",
        "status": "under_review",
        "review_url": "/v1/components/publish-reviews/pub_bbb222"
    }
}
```

---

## 9. 领域适配器管理API 🆕

### 9.1 概述

领域适配器是V2.0的核心创新，实现"孵化器-鸡"分离架构。通过适配器，平台可以轻松扩展到新领域（K12、美术史、畅销书、职业培训等）。

**核心功能**：
- 适配器注册与发现
- 领域路由与智能选择
- 配置管理与热更新
- 健康检查与降级

---

### 9.2 列出所有适配器

**接口**: `GET /v2/adapters`

**查询参数**:
- `domain` (可选): 过滤领域 (k12/art_history/bestseller/vocational_training)
- `status` (可选): 过滤状态 (active/inactive/error)
- `enabled` (可选): 是否启用 (true/false)
- `page` (可选): 页码，默认1
- `page_size` (可选): 每页数量，默认20

**响应**:
```json
{
    "success": true,
    "data": {
        "adapters": [
            {
                "adapter_id": "k12-adapter-v2",
                "name": "K12教育适配器",
                "domain": "k12",
                "version": "2.0.1",
                "author": "MetaWorkflow Team",
                "description": "支持中国K12教育体系的内容生成",
                "status": "active",
                "enabled": true,
                "supported_languages": ["zh", "en"],
                "install_date": "2025-12-01T10:00:00Z",
                "last_updated": "2025-12-09T08:30:00Z",
                "usage_count": 15234,
                "average_rating": 4.8,
                "health_status": "healthy"
            },
            {
                "adapter_id": "art-history-adapter-v1",
                "name": "美术史适配器",
                "domain": "art_history",
                "version": "1.0.0",
                "author": "Art Team",
                "description": "艺术鉴赏与美术史教学内容生成",
                "status": "active",
                "enabled": true,
                "supported_languages": ["zh", "en"],
                "install_date": "2025-12-05T14:00:00Z",
                "usage_count": 892,
                "average_rating": 4.5,
                "health_status": "healthy"
            }
        ],
        "pagination": {
            "total": 6,
            "page": 1,
            "page_size": 20,
            "total_pages": 1
        }
    }
}
```

---

### 9.3 获取适配器详情

**接口**: `GET /v2/adapters/{adapter_id}`

**响应**:
```json
{
    "success": true,
    "data": {
        "adapter_id": "k12-adapter-v2",
        "name": "K12教育适配器",
        "domain": "k12",
        "version": "2.0.1",
        "author": "MetaWorkflow Team",
        "author_email": "team@metaworkflow.ai",
        "description": "支持中国K12教育体系（1-12年级）的智能内容生成",
        
        "requires_python": ">=3.9",
        "dependencies": [
            "numpy>=1.20",
            "pandas>=1.3",
            "jieba>=0.42"
        ],
        
        "config_schema": {
            "type": "object",
            "properties": {
                "max_objectives": {
                    "type": "integer",
                    "default": 10,
                    "description": "最大学习目标数"
                },
                "grade_system": {
                    "type": "string",
                    "enum": ["cn", "us"],
                    "default": "cn",
                    "description": "年级体系"
                }
            }
        },
        
        "default_config": {
            "max_objectives": 10,
            "grade_system": "cn",
            "enable_ai_objectives": true
        },
        
        "supported_languages": ["zh", "en"],
        "status": "active",
        "enabled": true,
        "install_date": "2025-12-01T10:00:00Z",
        "last_updated": "2025-12-09T08:30:00Z",
        
        "homepage": "https://docs.metaworkflow.ai/adapters/k12",
        "repository": "https://github.com/metaworkflow/k12-adapter",
        "documentation": "https://docs.metaworkflow.ai/adapters/k12/guide",
        "license": "MIT",
        "tags": ["education", "k12", "china"],
        
        "usage_count": 15234,
        "average_rating": 4.8,
        "total_ratings": 256,
        
        "health_status": "healthy",
        "last_health_check": "2025-12-09T10:00:00Z"
    }
}
```

---

### 9.4 注册新适配器

**接口**: `POST /v2/adapters`

**请求头**:
```
Authorization: Bearer {admin_token}
```

**请求体**:
```json
{
    "adapter_id": "finance-adapter-v1",
    "name": "金融理财适配器",
    "domain": "finance",
    "version": "1.0.0",
    "author": "Finance Team",
    "author_email": "finance@metaworkflow.ai",
    "description": "个人理财与投资教育内容生成",
    
    "requires_python": ">=3.9",
    "dependencies": [
        "pandas>=1.3",
        "numpy>=1.20"
    ],
    
    "config_schema": {
        "type": "object",
        "properties": {
            "risk_level": {
                "type": "string",
                "enum": ["low", "medium", "high"],
                "default": "medium"
            }
        }
    },
    
    "default_config": {
        "risk_level": "medium",
        "enable_stock_examples": false
    },
    
    "supported_languages": ["zh", "en"],
    "homepage": "https://finance.metaworkflow.ai",
    "license": "MIT",
    "tags": ["finance", "investment", "education"]
}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "adapter_id": "finance-adapter-v1",
        "status": "installed",
        "install_date": "2025-12-09T10:30:00Z",
        "activation_url": "/v2/adapters/finance-adapter-v1/activate"
    }
}
```

---

### 9.5 更新适配器配置

**接口**: `PUT /v2/adapters/{adapter_id}/config`

**请求体**:
```json
{
    "settings": {
        "max_objectives": 15,
        "grade_system": "us",
        "enable_ai_objectives": true,
        "custom_prompt_template": "为美国{grade}年级学生设计..."
    },
    "enabled": true,
    "priority": 1,
    "timeout": 30,
    "max_retries": 3,
    "cache_enabled": true,
    "cache_ttl": 7200,
    "fallback_enabled": true,
    "fallback_adapter_id": "k12-adapter-v1"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "config_id": 123,
        "adapter_id": "k12-adapter-v2",
        "updated_at": "2025-12-09T10:35:00Z",
        "restart_required": false
    }
}
```

---

### 9.6 激活/停用适配器

**接口**: `POST /v2/adapters/{adapter_id}/activate`

**请求体**:
```json
{
    "enabled": true
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "adapter_id": "k12-adapter-v2",
        "status": "active",
        "enabled": true,
        "activated_at": "2025-12-09T10:40:00Z"
    }
}
```

---

### 9.7 领域路由（智能选择适配器）

**接口**: `POST /v2/adapters/route`

**请求体**:
```json
{
    "raw_input": "我想学习莫奈的《印象·日出》这幅画",
    "domain": null,  // 不指定领域，让AI自动判断
    "user_id": "user_123",
    "prefer_strategy": "llm"  // explicit/keyword/llm/history
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "selected_adapter_id": "art-history-adapter-v1",
        "domain": "art_history",
        "strategy": "llm",
        "confidence": 0.95,
        "routing_time": 0.23,
        
        "candidates": [
            {
                "adapter_id": "art-history-adapter-v1",
                "score": 0.95,
                "reason": "LLM detected art history keywords: 莫奈, 印象·日出, 画"
            },
            {
                "adapter_id": "k12-adapter-v2",
                "score": 0.15,
                "reason": "Low keyword match"
            }
        ],
        
        "decision_tree": {
            "step1_explicit": "No domain specified",
            "step2_keyword": "Partial match: art_history (0.6), k12 (0.2)",
            "step3_llm": "LLM classified as art_history with confidence 0.95",
            "final_decision": "art_history"
        }
    }
}
```

---

### 9.8 健康检查

**接口**: `GET /v2/adapters/{adapter_id}/health`

**响应**:
```json
{
    "success": true,
    "data": {
        "adapter_id": "k12-adapter-v2",
        "health_status": "healthy",  // healthy/degraded/unhealthy
        "last_check": "2025-12-09T10:50:00Z",
        "checks": {
            "api_reachable": {
                "status": "pass",
                "latency_ms": 23
            },
            "dependencies_loaded": {
                "status": "pass"
            },
            "config_valid": {
                "status": "pass"
            },
            "error_rate": {
                "status": "pass",
                "value": "0.002",
                "threshold": "0.05"
            }
        },
        "uptime_seconds": 691200,
        "last_error": null
    }
}
```

---

### 9.9 删除适配器

**接口**: `DELETE /v2/adapters/{adapter_id}`

**请求头**:
```
Authorization: Bearer {admin_token}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "adapter_id": "finance-adapter-v1",
        "deleted_at": "2025-12-09T11:00:00Z",
        "cleanup_completed": true
    }
}
```

---

## 10. AI目标生成API 🆕

### 10.1 概述

AI目标生成器将用户的模糊需求转化为清晰、结构化的学习目标。它是V2.0的核心AI能力之一。

**核心功能**：
- 智能分析用户输入
- 生成结构化学习目标（符合布鲁姆分类法）
- 受众画像分析
- 目标验证与优化
- 交互式反馈迭代

---

### 10.2 提交需求分析请求

**接口**: `POST /v2/ai/objectives/analyze`

**请求体**:
```json
{
    "raw_input": "我想教三年级学生学习除法，他们已经掌握了乘法",
    "domain": "k12",  // 可选，不指定则自动识别
    "language": "zh",
    "preferences": {
        "max_objectives": 8,
        "difficulty_level": "medium",
        "duration": 15,  // 课程总时长（分钟），微课建议3-15分钟，常规课45-90分钟
        "include_examples": true
    },
    "user_id": "user_123",
    "session_id": "session_abc"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "request_id": "req_abc123",
        "status": "processing",
        "estimated_time": 5,
        "poll_url": "/v2/ai/objectives/requests/req_abc123",
        "webhook_url": null
    }
}
```

---

### 10.3 获取分析结果

**接口**: `GET /v2/ai/objectives/requests/{request_id}`

**响应** (200 OK):
```json
{
    "success": true,
    "data": {
        "request_id": "req_abc123",
        "status": "completed",
        "processing_time": 4.2,
        
        "audience_profile": {
            "profile_id": "prof_xyz789",
            "age_range": "8-10",
            "education_level": "小学三年级",
            "role": "student",
            "interests": ["数学", "游戏", "动手操作"],
            "prior_knowledge": ["乘法", "加法", "减法", "数位概念"],
            "learning_goals": ["理解除法概念", "掌握除法计算"],
            "learning_style": "visual",
            "motivation_level": "medium",
            "attention_span": 20,
            "common_errors": ["除数和被除数混淆", "余数理解困难"],
            "pain_points": ["抽象概念理解", "计算准确性"],
            "confidence": 0.92,
            "source": "ai_inferred"
        },
        
        "objectives": [
            {
                "objective_id": "obj_001",
                "description": "理解除法的基本概念，能将除法与日常生活中的分配场景联系起来",
                "bloom_level": 2,  // 理解层级
                "objective_type": "knowledge",
                "domain": "k12",
                "parent_id": null,
                "order": 1,
                "success_criteria": [
                    "能用自己的话解释除法的含义",
                    "能举出生活中至少3个除法的例子",
                    "能区分除法与减法的不同"
                ],
                "assessment_methods": [
                    "口头回答问题",
                    "画图表示除法过程",
                    "举例说明"
                ],
                "estimated_time": 5,  // 此目标预估5分钟（根据总课程时长自动分配）
                "difficulty": "easy",
                "metadata": {
                    "grade": 3,
                    "subject": "数学",
                    "topic": "除法",
                    "curriculum_standard": "新课标2022"
                }
            },
            {
                "objective_id": "obj_002",
                "description": "掌握表内除法（1-9的除法表）的计算方法",
                "bloom_level": 3,  // 应用层级
                "objective_type": "skill",
                "domain": "k12",
                "parent_id": "obj_001",
                "order": 2,
                "success_criteria": [
                    "能正确计算81以内的除法",
                    "计算准确率达到90%以上",
                    "能在3分钟内完成20道表内除法题"
                ],
                "assessment_methods": [
                    "纸笔测试",
                    "计时练习",
                    "口算比赛"
                ],
                "estimated_time": 7,  // 此目标预估7分钟
                "difficulty": "medium",
                "metadata": {
                    "grade": 3,
                    "subject": "数学",
                    "topic": "表内除法"
                }
            },
            {
                "objective_id": "obj_003",
                "description": "理解余数的概念，能正确处理有余数的除法问题",
                "bloom_level": 4,  // 分析层级
                "objective_type": "skill",
                "domain": "k12",
                "parent_id": "obj_002",
                "order": 3,
                "success_criteria": [
                    "能解释什么是余数",
                    "能判断何时会产生余数",
                    "能正确计算有余数的除法"
                ],
                "assessment_methods": [
                    "应用题解答",
                    "概念解释",
                    "实际操作"
                ],
                "estimated_time": 25,
                "difficulty": "medium"
            }
        ],
        
        "validation_results": [
            {
                "validation_id": "val_001",
                "objective_id": "obj_001",
                "is_measurable": true,
                "is_achievable": true,
                "is_specific": true,
                "is_relevant": true,
                "is_time_bound": true,
                "overall_score": 0.95,
                "issues": [],
                "suggestions": []
            },
            {
                "validation_id": "val_002",
                "objective_id": "obj_002",
                "is_measurable": true,
                "is_achievable": true,
                "is_specific": true,
                "is_relevant": true,
                "is_time_bound": true,
                "overall_score": 0.98,
                "issues": [],
                "suggestions": ["可以添加更具体的时间限制"]
            }
        ],
        
        "statistics": {
            "total_objectives": 3,
            "by_bloom_level": {
                "2": 1,
                "3": 1,
                "4": 1
            },
            "by_type": {
                "knowledge": 1,
                "skill": 2
            },
            "total_estimated_time": 70,
            "average_validation_score": 0.96
        }
    }
}
```

---

### 10.4 交互式反馈（调整目标）

**接口**: `POST /v2/ai/objectives/{objective_id}/feedback`

**请求体**:
```json
{
    "feedback_type": "modify",  // like/dislike/modify/delete/add
    "comment": "这个目标太简单了，可以提高难度",
    "modification": {
        "description": "能分析除法在实际问题中的应用，并设计除法问题",
        "bloom_level": 5,  // 提升到评价层级
        "difficulty": "hard"
    },
    "user_id": "user_123"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "feedback_id": "fbk_123",
        "objective_id": "obj_001",
        "applied": true,
        "updated_objective": {
            "objective_id": "obj_001",
            "description": "能分析除法在实际问题中的应用，并设计除法问题",
            "bloom_level": 5,
            "difficulty": "hard",
            "updated_at": "2025-12-09T11:10:00Z"
        },
        "re_validation_required": true
    }
}
```

---

### 10.5 批量生成目标

**接口**: `POST /v2/ai/objectives/batch`

**请求体**:
```json
{
    "requests": [
        {
            "raw_input": "教三年级除法",
            "domain": "k12",
            "preferences": {"duration": 10}  // 10分钟微课
        },
        {
            "raw_input": "介绍莫奈的印象派作品",
            "domain": "art_history",
            "preferences": {"duration": 15}  // 15分钟艺术鉴赏
        },
        {
            "raw_input": "解读《原则》这本书",
            "domain": "bestseller_interpretation",
            "preferences": {"duration": 12}  // 12分钟书籍解读
        }
    ],
    "user_id": "user_123"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "batch_id": "batch_abc",
        "total_requests": 3,
        "status": "processing",
        "request_ids": [
            "req_001",
            "req_002",
            "req_003"
        ],
        "poll_url": "/v2/ai/objectives/batches/batch_abc"
    }
}
```

---

### 10.6 导出目标（多格式）

**接口**: `GET /v2/ai/objectives/requests/{request_id}/export`

**查询参数**:
- `format`: 导出格式 (json/yaml/markdown/pdf)
- `include_validation`: 是否包含验证结果 (true/false)

**响应** (格式为markdown时):
```markdown
# 学习目标报告

## 受众画像
- **年龄段**: 8-10岁
- **教育水平**: 小学三年级
- **学习风格**: 视觉型

## 学习目标

### 目标1: 理解除法的基本概念
- **认知层级**: 理解 (Bloom Level 2)
- **类型**: 知识
- **预估时长**: 5分钟（根据课程总时长自动分配）
- **难度**: 简单

**成功标准**:
1. 能用自己的话解释除法的含义
2. 能举出生活中至少3个除法的例子
3. 能区分除法与减法的不同

**评估方法**:
- 口头回答问题
- 画图表示除法过程
- 举例说明

...
```

---

## 11. AI内容发现API 🆕

### 11.1 概述

AI内容发现引擎通过向量检索和结构化提取，从Wikipedia、WikiArt、arXiv等多个内容源中智能发现相关内容。

**核心功能**：
- 向量语义搜索（Milvus）
- 多源内容聚合
- 结构化内容提取
- 质量评分与过滤
- 智能排序与重排序

---

### 11.2 配置内容源

**接口**: `POST /v2/ai/discovery/sources`

**请求头**:
```
Authorization: Bearer {admin_token}
```

**请求体**:
```json
{
    "source_id": "wikiart-api",
    "name": "WikiArt艺术数据库",
    "source_type": "custom_api",
    "base_url": "https://www.wikiart.org/en/api/2",
    "api_key": "your_api_key_here",
    "auth_type": "api_key",
    
    "supported_content_types": ["image", "text"],
    "max_results_per_query": 50,
    "rate_limit": 10,
    
    "parser_class": "WikiArtParser",
    "extraction_config": {
        "image_selector": ".artwork-img",
        "title_selector": "h1.artwork-title",
        "artist_selector": ".artist-name"
    },
    
    "default_quality_weight": 0.9,
    "reliability_score": 0.85,
    "enabled": true
}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "source_id": "wikiart-api",
        "status": "active",
        "created_at": "2025-12-09T11:20:00Z",
        "test_query_url": "/v2/ai/discovery/sources/wikiart-api/test"
    }
}
```

---

### 11.3 向量语义搜索

**接口**: `POST /v2/ai/discovery/search`

**请求体**:
```json
{
    "query": "莫奈的印象派作品，特别是关于光影和色彩的运用",
    "domain": "art_history",
    "filters": {
        "content_type": "image",
        "min_quality": 0.7,
        "language": "zh"
    },
    "top_k": 10,
    "similarity_threshold": 0.5,
    "target_sources": ["wikiart-api", "wikipedia"],  // 空数组表示所有源
    "enable_reranking": true,  // 是否使用重排序模型
    "user_id": "user_123"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "query_id": "query_abc123",
        "results": [
            {
                "item_id": "item_001",
                "source_id": "wikiart-api",
                "title": "印象·日出（Impression, Sunrise）",
                "content": "《印象·日出》是克劳德·莫奈于1872年在勒阿弗尔港口创作的一幅油画...",
                "content_type": "image",
                "url": "https://www.wikiart.org/en/claude-monet/impression-sunrise-1872",
                "author": "Claude Monet",
                "publish_date": "1872-01-01T00:00:00Z",
                "language": "zh",
                
                "images": [
                    "https://uploads.wikiart.org/images/claude-monet/impression-sunrise.jpg!Large.jpg"
                ],
                "tags": ["印象派", "莫奈", "日出", "光影"],
                "categories": ["油画", "风景画", "印象派"],
                
                "metadata": {
                    "artist": "莫奈",
                    "movement": "印象派",
                    "year": 1872,
                    "medium": "油画",
                    "dimensions": "48 × 63 cm"
                },
                
                "similarity_score": 0.92,
                "quality_score": 0.88,
                "relevance_score": 0.95,
                "final_score": 0.92
            },
            {
                "item_id": "item_002",
                "source_id": "wikiart-api",
                "title": "睡莲系列（Water Lilies）",
                "content": "莫奈的睡莲系列是其晚年最著名的作品之一...",
                "similarity_score": 0.87,
                "quality_score": 0.90,
                "relevance_score": 0.85,
                "final_score": 0.87,
                "images": [
                    "https://uploads.wikiart.org/images/claude-monet/water-lilies.jpg"
                ],
                "metadata": {
                    "artist": "莫奈",
                    "series": "睡莲系列",
                    "years": "1897-1926"
                }
            }
        ],
        
        "statistics": {
            "total_matches": 45,
            "returned_count": 10,
            "search_time": 0.28,
            "source_distribution": {
                "wikiart-api": 7,
                "wikipedia": 3
            },
            "avg_quality_score": 0.85,
            "embedding_model": "text-embedding-3-small"
        }
    }
}
```

---

### 11.4 结构化内容提取

**接口**: `POST /v2/ai/discovery/extract`

**请求体**:
```json
{
    "item_id": "item_001",
    "extraction_schema": {
        "type": "object",
        "properties": {
            "artwork": {
                "type": "object",
                "properties": {
                    "title": {"type": "string"},
                    "artist": {"type": "string"},
                    "year": {"type": "integer"},
                    "movement": {"type": "string"},
                    "medium": {"type": "string"},
                    "dimensions": {"type": "string"}
                }
            },
            "analysis": {
                "type": "object",
                "properties": {
                    "color_palette": {"type": "array"},
                    "composition": {"type": "string"},
                    "techniques": {"type": "array"},
                    "themes": {"type": "array"}
                }
            },
            "historical_context": {
                "type": "string"
            }
        }
    },
    "extraction_method": "llm",  // llm/regex/xpath
    "llm_model": "gpt-4"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "structured_id": "struct_123",
        "item_id": "item_001",
        
        "structured_data": {
            "artwork": {
                "title": "印象·日出",
                "artist": "克劳德·莫奈",
                "year": 1872,
                "movement": "印象派",
                "medium": "布面油画",
                "dimensions": "48 × 63 cm"
            },
            "analysis": {
                "color_palette": [
                    "橙红色（太阳）",
                    "深蓝色（水面）",
                    "灰色（雾气）",
                    "黑色（船只轮廓）"
                ],
                "composition": "太阳位于画面中央偏左，通过雾气营造朦胧感，船只的剪影形成前景",
                "techniques": [
                    "快速笔触",
                    "色彩并置",
                    "光影对比",
                    "大气透视"
                ],
                "themes": [
                    "工业化时代的港口",
                    "光线的瞬间印象",
                    "自然与工业的共存"
                ]
            },
            "historical_context": "这幅画创作于1872年的勒阿弗尔港口，正值法国工业革命时期。莫奈通过描绘日出时分的港口景象，捕捉了光线和色彩的瞬间印象，这幅作品后来成为印象派运动的命名由来。"
        },
        
        "extraction_confidence": 0.94,
        "is_validated": true,
        "validation_errors": [],
        "extracted_at": "2025-12-09T11:30:00Z",
        "extraction_time": 3.5
    }
}
```

---

### 11.5 质量评分

**接口**: `POST /v2/ai/discovery/items/{item_id}/score`

**请求体**:
```json
{
    "scoring_method": "llm",  // llm/rule_based/hybrid
    "scorer_model": "gpt-4",
    "criteria": {
        "relevance_weight": 0.3,
        "accuracy_weight": 0.25,
        "completeness_weight": 0.2,
        "authority_weight": 0.15,
        "freshness_weight": 0.1
    }
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "score_id": "score_456",
        "item_id": "item_001",
        
        "dimensions": {
            "relevance": 0.95,
            "accuracy": 0.92,
            "completeness": 0.88,
            "readability": 0.90,
            "authority": 0.85,
            "freshness": 0.70
        },
        
        "overall_score": 0.88,
        "quality_level": "excellent",
        
        "scoring_details": {
            "relevance_reason": "内容高度匹配查询意图，直接讲述莫奈的印象派作品",
            "accuracy_reason": "信息准确，来源可靠（WikiArt官方数据库）",
            "completeness_reason": "包含作品基本信息、艺术分析、历史背景，稍缺技法细节",
            "authority_reason": "来自权威艺术数据库，有引用来源"
        },
        
        "issues": [
            "部分技术细节描述不够深入"
        ],
        
        "scored_at": "2025-12-09T11:35:00Z"
    }
}
```

---

### 11.6 获取内容源列表

**接口**: `GET /v2/ai/discovery/sources`

**查询参数**:
- `source_type` (可选): 过滤类型
- `enabled` (可选): 是否启用
- `content_type` (可选): 支持的内容类型

**响应**:
```json
{
    "success": true,
    "data": {
        "sources": [
            {
                "source_id": "wikiart-api",
                "name": "WikiArt艺术数据库",
                "source_type": "custom_api",
                "enabled": true,
                "supported_content_types": ["image", "text"],
                "reliability_score": 0.85,
                "total_items": 250000,
                "last_successful_fetch": "2025-12-09T10:00:00Z",
                "error_count": 0,
                "health_status": "healthy"
            },
            {
                "source_id": "wikipedia",
                "name": "维基百科",
                "source_type": "wikipedia",
                "enabled": true,
                "supported_content_types": ["text", "image"],
                "reliability_score": 0.95,
                "total_items": 6500000,
                "last_successful_fetch": "2025-12-09T11:00:00Z",
                "error_count": 0,
                "health_status": "healthy"
            }
        ],
        "total_sources": 6,
        "active_sources": 6
    }
}
```

---

### 11.7 缓存管理

**接口**: `DELETE /v2/ai/discovery/cache`

**查询参数**:
- `query_hash` (可选): 特定查询的缓存
- `clear_all` (可选): 清空所有缓存

**响应**:
```json
{
    "success": true,
    "data": {
        "cleared_entries": 125,
        "freed_space_mb": 45.2,
        "cleared_at": "2025-12-09T11:40:00Z"
    }
}
```

---

## 12. AI叙述生成API 🆕

### 12.1 概述

AI叙述生成引擎将枯燥的教学内容转化为引人入胜的故事化叙述，提升学习体验和知识留存率。

**核心功能**：
- 多种叙述策略（故事讲述、历史叙述、案例研究等）
- 故事元素生成（人物、场景、冲突、解决）
- 多幕结构编排
- 质量检测与优化
- 多媒体资源配图

---

### 12.2 列出叙述策略

**接口**: `GET /v2/ai/narrative/strategies`

**查询参数**:
- `domain` (可选): 过滤领域
- `target_age` (可选): 目标年龄段
- `narrative_style` (可选): 叙述风格

**响应**:
```json
{
    "success": true,
    "data": {
        "strategies": [
            {
                "strategy_id": "k12-story-telling",
                "name": "K12故事讲述策略",
                "domain": "k12",
                "narrative_style": "story_telling",
                "story_arc": "three_act",
                "emotional_tone": "inspiring",
                "target_age_range": "8-12",
                "reading_level": "小学3-6年级",
                "use_characters": true,
                "use_dialogue": true,
                "use_metaphors": true,
                "max_sections": 8,
                "usage_count": 1234
            },
            {
                "strategy_id": "art-biography",
                "name": "艺术家传记策略",
                "domain": "art_history",
                "narrative_style": "biography",
                "story_arc": "chronological",
                "emotional_tone": "curious",
                "target_age_range": "15+",
                "reading_level": "高中及以上",
                "use_characters": true,
                "use_dialogue": false,
                "max_sections": 12,
                "usage_count": 567
            }
        ],
        "total": 15
    }
}
```

---

### 12.3 创建叙述内容

**接口**: `POST /v2/ai/narrative/generate`

**请求体**:
```json
{
    "objective_id": "obj_001",
    "strategy_id": "art-biography",
    "content_items": ["item_001", "item_002"],  // 从内容发现获取的item_id
    
    "customization": {
        "title": "莫奈：光影的诗人",
        "emotional_tone": "inspiring",
        "max_word_count": 2000,
        "include_images": true,
        "image_count": 5
    },
    
    "user_id": "user_123"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "narrative_id": "narr_abc123",
        "status": "generating",
        "estimated_time": 15,
        "poll_url": "/v2/ai/narrative/narr_abc123"
    }
}
```

---

### 12.4 获取叙述内容

**接口**: `GET /v2/ai/narrative/{narrative_id}`

**响应**:
```json
{
    "success": true,
    "data": {
        "narrative_id": "narr_abc123",
        "objective_id": "obj_001",
        "strategy_id": "art-biography",
        "status": "completed",
        
        "title": "莫奈：光影的诗人",
        "subtitle": "从印象·日出到睡莲系列的艺术之旅",
        
        "sections": [
            {
                "section_id": "sec_001",
                "order": 1,
                "act": 1,
                "title": "勒阿弗尔的黎明",
                "content": "1872年的一个清晨，勒阿弗尔港口笼罩在薄雾之中。32岁的克劳德·莫奈站在旅馆的窗前，凝视着逐渐升起的太阳。橙红色的光芒穿透雾气，在深蓝色的水面上投下一道颤动的倒影。港口里的船只若隐若现，烟囱冒出的烟雾与晨雾融为一体...",
                "word_count": 245,
                "narrative_function": "exposition",
                "story_elements": ["char_001", "setting_001"],
                "transition_out": "这幅后来被称为《印象·日出》的作品，不仅改变了莫奈的命运，更开创了一个全新的艺术流派。"
            },
            {
                "section_id": "sec_002",
                "order": 2,
                "act": 1,
                "title": "印象派的诞生",
                "content": "1874年，莫奈和他的朋友们——雷诺阿、德加、毕沙罗——决定举办一场独立的画展。当时的巴黎艺术界被官方沙龙所主导，他们的作品因为'不够完整'而屡遭拒绝。\n\n一位名叫路易·勒鲁瓦的评论家在报纸上讽刺地写道：'这些画简直就是印象，连壁纸都比它们完整！'他以莫奈的《印象·日出》为靶子，称这些艺术家为'印象派'...",
                "word_count": 312,
                "narrative_function": "rising_action",
                "story_elements": ["char_002", "conflict_001"]
            },
            {
                "section_id": "sec_003",
                "order": 3,
                "act": 2,
                "title": "光影的革命",
                "content": "莫奈对光线的痴迷近乎疯狂。他会在不同时刻对同一场景进行反复描绘——鲁昂大教堂在清晨、正午、黄昏的样子各不相同。他追求的不是物体本身，而是光线照在物体上的瞬间印象...",
                "word_count": 398,
                "narrative_function": "climax"
            }
        ],
        
        "story_elements": [
            {
                "element_id": "char_001",
                "element_type": "character",
                "name": "克劳德·莫奈",
                "description": "法国印象派画家，光影的捕捉者",
                "character_traits": ["敏锐", "坚持", "创新"],
                "character_role": "protagonist"
            },
            {
                "element_id": "setting_001",
                "element_type": "setting",
                "name": "勒阿弗尔港口",
                "time_period": "1872年",
                "location": "法国勒阿弗尔",
                "atmosphere": "朦胧、宁静、充满希望"
            },
            {
                "element_id": "conflict_001",
                "element_type": "conflict",
                "name": "与传统艺术界的冲突",
                "conflict_type": "external",
                "stakes": "艺术理念的认可与生存"
            }
        ],
        
        "statistics": {
            "total_word_count": 1856,
            "total_sections": 8,
            "estimated_reading_time": 12,
            "acts": 3
        },
        
        "quality_assessment": {
            "engagement_score": 0.88,
            "coherence_score": 0.92,
            "educational_value": 0.85,
            "age_appropriateness": 0.90,
            "overall_quality": 0.89
        },
        
        "images": [
            "https://storage.metaworkflow.ai/narratives/narr_abc123/impression-sunrise.jpg",
            "https://storage.metaworkflow.ai/narratives/narr_abc123/rouen-cathedral.jpg",
            "https://storage.metaworkflow.ai/narratives/narr_abc123/water-lilies.jpg"
        ],
        
        "generated_at": "2025-12-09T11:50:00Z",
        "generation_time": 12.3,
        "llm_model": "gpt-4"
    }
}
```

---

### 12.5 质量检测

**接口**: `POST /v2/ai/narrative/{narrative_id}/check-quality`

**响应**:
```json
{
    "success": true,
    "data": {
        "narrative_id": "narr_abc123",
        "quality_assessment": {
            "engagement_score": 0.88,
            "coherence_score": 0.92,
            "educational_value": 0.85,
            "age_appropriateness": 0.90,
            "overall_quality": 0.89
        },
        
        "issues": [
            {
                "issue_id": "issue_001",
                "section_id": "sec_003",
                "issue_type": "readability",
                "severity": "low",
                "description": "第3段第2句过长，可能影响阅读流畅性",
                "affected_text": "他追求的不是物体本身，而是光线照在物体上的瞬间印象，这种对光影变化的执着观察...",
                "position": 156,
                "suggestion": "建议拆分为两句",
                "auto_fixable": true
            },
            {
                "issue_id": "issue_002",
                "section_id": "sec_005",
                "issue_type": "factual_error",
                "severity": "medium",
                "description": "年份可能有误",
                "affected_text": "1890年，莫奈购买了吉维尼的房产",
                "suggestion": "核实准确年份应为1883年",
                "auto_fixable": false
            }
        ],
        
        "recommendations": [
            "增加更多具体的艺术技法描述",
            "可以添加与读者的互动提问",
            "建议在第二幕增加情感高潮"
        ],
        
        "checked_at": "2025-12-09T12:00:00Z"
    }
}
```

---

### 12.6 修复质量问题

**接口**: `POST /v2/ai/narrative/{narrative_id}/fix-issues`

**请求体**:
```json
{
    "issue_ids": ["issue_001"],  // 只修复可自动修复的问题
    "auto_fix_all": false
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "narrative_id": "narr_abc123",
        "fixed_issues": 1,
        "skipped_issues": 1,
        "updated_sections": ["sec_003"],
        "fixed_at": "2025-12-09T12:05:00Z"
    }
}
```

---

### 12.7 导出叙述内容

**接口**: `GET /v2/ai/narrative/{narrative_id}/export`

**查询参数**:
- `format`: 导出格式 (markdown/html/pdf/audio)
- `include_images`: 是否包含图片 (true/false)
- `include_metadata`: 是否包含元数据 (true/false)

**响应** (格式为markdown时):
```markdown
# 莫奈：光影的诗人
> 从印象·日出到睡莲系列的艺术之旅

---

## 第一幕：勒阿弗尔的黎明

1872年的一个清晨，勒阿弗尔港口笼罩在薄雾之中。32岁的克劳德·莫奈站在旅馆的窗前，凝视着逐渐升起的太阳...

![印象·日出](https://storage.metaworkflow.ai/narratives/narr_abc123/impression-sunrise.jpg)

---

## 第一幕：印象派的诞生

1874年，莫奈和他的朋友们——雷诺阿、德加、毕沙罗——决定举办一场独立的画展...

---

**元数据**
- 字数：1,856
- 阅读时间：约12分钟
- 质量评分：8.9/10
```

---

### 12.8 生成音频朗读

**接口**: `POST /v2/ai/narrative/{narrative_id}/synthesize-audio`

**请求体**:
```json
{
    "voice": "zh-CN-XiaoxiaoNeural",  // Azure TTS语音
    "speed": 1.0,
    "format": "mp3"
}
```

**响应** (202 Accepted):
```json
{
    "success": true,
    "data": {
        "audio_id": "audio_123",
        "status": "synthesizing",
        "estimated_time": 30,
        "poll_url": "/v2/ai/narrative/audio/audio_123"
    }
}
```

---

## 13. 多源内容检索API 🆕

### 13.1 概述

多源内容检索API提供统一的接口，聚合来自Milvus向量数据库、Redis缓存、MinIO对象存储、PostgreSQL关系数据库的查询结果。

**核心功能**：
- 统一查询接口
- 多存储后端聚合
- 智能缓存策略
- 级联查询优化
- 对象存储管理

---

### 13.2 统一查询接口

**接口**: `POST /v2/storage/multi-search`

**请求体**:
```json
{
    "query_text": "莫奈的印象派作品",
    "target_backends": ["milvus", "postgresql"],  // 可选: milvus/redis/postgresql/minio
    "strategy": "cascade",  // parallel/sequential/cascade
    "merge_method": "rank_fusion",  // union/intersect/rank_fusion
    "max_results": 20,
    "use_cache": true,
    "cache_ttl": 3600,
    "timeout": 5.0,
    "domain": "art_history",
    "user_id": "user_123"
}
```

**响应说明**:
- `strategy`:
  - `parallel`: 并行查询所有后端，合并结果
  - `sequential`: 顺序查询，直到找到足够结果
  - `cascade`: 级联查询（Redis缓存 → Milvus向量库 → PostgreSQL）
- `merge_method`:
  - `union`: 简单合并，去重
  - `intersect`: 取交集
  - `rank_fusion`: 重排序融合（如RRF算法）

**响应**:
```json
{
    "success": true,
    "data": {
        "query_id": "query_multi_123",
        "results": [
            {
                "item_id": "item_001",
                "title": "印象·日出",
                "content": "...",
                "source_backend": "milvus",
                "score": 0.95,
                "metadata": {
                    "artist": "莫奈",
                    "year": 1872
                }
            },
            {
                "item_id": "item_002",
                "title": "睡莲系列",
                "source_backend": "postgresql",
                "score": 0.87
            }
        ],
        
        "statistics": {
            "total_results": 18,
            "returned_count": 18,
            "backend_distribution": {
                "milvus": 12,
                "postgresql": 6
            },
            "total_time": 0.45,
            "backend_times": {
                "redis": 0.02,  // cache miss
                "milvus": 0.28,
                "postgresql": 0.15
            },
            "cache_hit": false,
            "strategy_used": "cascade"
        }
    }
}
```

---

### 13.3 Milvus向量搜索

**接口**: `POST /v2/storage/milvus/search`

**请求体**:
```json
{
    "collection_name": "content_embeddings",
    "query_vector": [0.123, 0.456, ...],  // 1536维向量
    "top_k": 10,
    "metric_type": "L2",  // L2/IP/COSINE
    "search_params": {
        "nprobe": 10
    },
    "filter_expression": "domain == 'art_history' and quality_score > 0.7",
    "output_fields": ["item_id", "title", "quality_score"]
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "result_id": "milvus_result_123",
        "hits": [
            {
                "id": "vec_001",
                "distance": 0.23,
                "entity": {
                    "item_id": "item_001",
                    "title": "印象·日出",
                    "quality_score": 0.88
                }
            },
            {
                "id": "vec_002",
                "distance": 0.31,
                "entity": {
                    "item_id": "item_002",
                    "title": "睡莲系列",
                    "quality_score": 0.90
                }
            }
        ],
        "total_hits": 10,
        "search_time": 0.28
    }
}
```

---

### 13.4 Redis缓存操作

**接口**: `POST /v2/storage/redis/set`

**请求体**:
```json
{
    "cache_key": "cache:query:莫奈印象派_hash123",
    "cache_value": {
        "items": [...],
        "total": 18,
        "cached_at": "2025-12-09T12:15:00Z"
    },
    "ttl": 3600,
    "strategy": "ttl"
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "cache_key": "cache:query:莫奈印象派_hash123",
        "set_at": "2025-12-09T12:15:00Z",
        "expires_at": "2025-12-09T13:15:00Z",
        "size_bytes": 15234
    }
}
```

---

**接口**: `GET /v2/storage/redis/get`

**查询参数**:
- `cache_key`: 缓存键

**响应**:
```json
{
    "success": true,
    "data": {
        "cache_key": "cache:query:莫奈印象派_hash123",
        "cache_value": {
            "items": [...],
            "total": 18
        },
        "hit_count": 25,
        "last_accessed": "2025-12-09T12:20:00Z",
        "ttl_remaining": 3455
    }
}
```

---

### 13.5 MinIO对象存储

**接口**: `POST /v2/storage/minio/upload`

**请求头**:
```
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**请求体** (multipart/form-data):
```
file: (binary)
bucket_name: content-assets
object_key: narratives/narr_abc123/cover.jpg
content_type: image/jpeg
is_public: false
related_entity_type: narrative_content
related_entity_id: narr_abc123
custom_metadata: {"quality": "high", "domain": "art_history"}
```

**响应** (201 Created):
```json
{
    "success": true,
    "data": {
        "object_id": "obj_minio_123",
        "bucket_name": "content-assets",
        "object_key": "narratives/narr_abc123/cover.jpg",
        "content_type": "image/jpeg",
        "size_bytes": 245678,
        "etag": "5d41402abc4b2a76b9719d911017c592",
        "uploaded_at": "2025-12-09T12:25:00Z",
        "url": "https://storage.metaworkflow.ai/content-assets/narratives/narr_abc123/cover.jpg",
        "signed_url": "https://storage.metaworkflow.ai/content-assets/narratives/narr_abc123/cover.jpg?X-Amz-Expires=3600&...",
        "signed_url_expires": "2025-12-09T13:25:00Z"
    }
}
```

---

**接口**: `GET /v2/storage/minio/download/{object_id}`

**查询参数**:
- `expires_in`: 签名URL有效期（秒），默认3600

**响应**:
```json
{
    "success": true,
    "data": {
        "object_id": "obj_minio_123",
        "signed_url": "https://storage.metaworkflow.ai/content-assets/narratives/narr_abc123/cover.jpg?X-Amz-Expires=3600&...",
        "expires_at": "2025-12-09T13:30:00Z",
        "content_type": "image/jpeg",
        "size_bytes": 245678
    }
}
```

---

**接口**: `DELETE /v2/storage/minio/{object_id}`

**响应**:
```json
{
    "success": true,
    "data": {
        "object_id": "obj_minio_123",
        "deleted_at": "2025-12-09T12:35:00Z"
    }
}
```

---

### 13.6 缓存统计

**接口**: `GET /v2/storage/redis/stats`

**响应**:
```json
{
    "success": true,
    "data": {
        "total_keys": 1523,
        "total_memory_mb": 245.6,
        "hit_rate": 0.78,
        "miss_rate": 0.22,
        "avg_ttl": 2456,
        
        "by_key_type": {
            "cache:query": 856,
            "cache:item": 342,
            "session": 125,
            "ratelimit": 98,
            "lock": 5,
            "trending": 12,
            "cache:vector": 85
        },
        
        "top_accessed_keys": [
            {
                "key": "cache:query:k12数学除法_hash456",
                "hit_count": 145,
                "last_accessed": "2025-12-09T12:30:00Z"
            }
        ],
        
        "evicted_keys_last_hour": 23,
        "expired_keys_last_hour": 156
    }
}
```

---

### 13.7 存储健康检查

**接口**: `GET /v2/storage/health`

**响应**:
```json
{
    "success": true,
    "data": {
        "overall_status": "healthy",
        "backends": {
            "milvus": {
                "status": "healthy",
                "latency_ms": 12,
                "collections": 2,
                "total_vectors": 1250000,
                "last_check": "2025-12-09T12:40:00Z"
            },
            "redis": {
                "status": "healthy",
                "latency_ms": 2,
                "total_keys": 1523,
                "memory_usage_mb": 245.6,
                "hit_rate": 0.78,
                "last_check": "2025-12-09T12:40:00Z"
            },
            "minio": {
                "status": "healthy",
                "latency_ms": 18,
                "total_buckets": 4,
                "total_objects": 8965,
                "total_size_gb": 45.2,
                "last_check": "2025-12-09T12:40:00Z"
            },
            "postgresql": {
                "status": "healthy",
                "latency_ms": 8,
                "total_tables": 35,
                "total_rows": 156789,
                "database_size_gb": 2.3,
                "last_check": "2025-12-09T12:40:00Z"
            }
        }
    }
}
```

---

### 13.8 批量预热缓存

**接口**: `POST /v2/storage/redis/warmup`

**请求体**:
```json
{
    "queries": [
        "K12数学除法",
        "莫奈印象派作品",
        "原则这本书的核心观点"
    ],
    "ttl": 7200
}
```

**响应**:
```json
{
    "success": true,
    "data": {
        "total_queries": 3,
        "cached_queries": 3,
        "failed_queries": 0,
        "total_items_cached": 45,
        "total_time": 2.3,
        "cache_size_mb": 12.5,
        "warmed_up_at": "2025-12-09T12:45:00Z"
    }
}
```

---

## 14. 配置管理API

### 6.1 获取工作流定义

**接口**: `GET /v1/config/workflows/{workflow_id}`

**响应**:
```json
{
    "success": true,
    "data": {
        "workflow_id": "math_lesson_basic",
        "name": "数学课程标准工作流",
        "version": "1.0",
        "nodes": [
            {
                "node_id": "persona",
                "node_type": "PersonaAnalysisNode",
                "config": {...}
            }
        ]
    }
}
```

### 6.2 创建/更新工作流定义

**接口**: `POST /v1/config/workflows` (创建)
**接口**: `PUT /v1/config/workflows/{workflow_id}` (更新)

**请求体**:
```json
{
    "workflow_id": "science_lesson_inquiry",
    "name": "科学课探究式工作流",
    "version": "1.0",
    "trigger_conditions": {
        "subject": "科学"
    },
    "nodes": [...]
}
```

### 6.3 获取规则列表

**接口**: `GET /v1/config/rules`

**查询参数**:
- `subject`: 学科过滤
- `region`: 地区过滤
- `enabled`: 是否启用(true/false)

**响应**:
```json
{
    "success": true,
    "data": {
        "rules": [
            {
                "rule_id": "math_logic_check",
                "name": "数学课程必须包含逻辑验证",
                "priority": 100,
                "conditions": {"subject": "数学"},
                "actions": [...],
                "enabled": true
            }
        ]
    }
}
```

### 6.4 更新规则

**接口**: `PUT /v1/config/rules/{rule_id}`

**请求体**:
```json
{
    "enabled": false,
    "priority": 90
}
```

### 5.5 获取提示词模板

**接口**: `GET /v1/config/prompts`

**查询参数**:
- `subject`: 学科
- `node_type`: 节点类型
- `region`: 地区

**响应**:
```json
{
    "success": true,
    "data": {
        "template_id": "math_content_generation_cn",
        "subject": "数学",
        "node_type": "ContentGenerationNode",
        "region": "CN",
        "template": "[ROLE]\n你是一位小学数学教师...\n[INSTRUCTION]\n...",
        "variables": ["topic", "persona", "objectives"],
        "version": "1.2"
    }
}
```

---

## 10. 监控与统计API

### 10.1 获取系统健康状态

**接口**: `GET /v1/monitoring/health`

**响应**:
```json
{
    "success": true,
    "data": {
        "status": "healthy",
        "components": {
            "database": {
                "status": "up",
                "latency_ms": 5
            },
            "redis": {
                "status": "up",
                "latency_ms": 2
            },
            "vector_db": {
                "status": "up",
                "latency_ms": 15
            },
            "llm_providers": {
                "openai": {
                    "status": "up",
                    "latency_ms": 500
                },
                "anthropic": {
                    "status": "up",
                    "latency_ms": 450
                }
            }
        },
        "timestamp": "2025-12-09T10:00:00Z"
    }
}
```

### 6.2 获取执行统计

**接口**: `GET /v1/monitoring/stats`

**查询参数**:
- `from_date`: 开始日期
- `to_date`: 结束日期
- `granularity`: 粒度(hour/day/week/month)

**响应**:
```json
{
    "success": true,
    "data": {
        "period": {
            "from": "2025-12-01T00:00:00Z",
            "to": "2025-12-09T23:59:59Z"
        },
        "total_instances": 1250,
        "success_rate": 0.95,
        "avg_duration_seconds": 165,
        "total_cost": 125.50,
        "by_subject": {
            "数学": 450,
            "语文": 380,
            "英语": 270,
            "科学": 150
        },
        "by_status": {
            "completed": 1188,
            "failed": 62
        }
    }
}
```

### 6.3 获取成本分析

**接口**: `GET /v1/monitoring/cost`

**查询参数**:
- `from_date`: 开始日期
- `to_date`: 结束日期
- `group_by`: 分组维度(model/subject/region)

**响应**:
```json
{
    "success": true,
    "data": {
        "total_cost": 125.50,
        "total_tokens": 2500000,
        "breakdown": {
            "gpt-4": {
                "cost": 85.30,
                "tokens": 1200000,
                "calls": 850
            },
            "deepseek-v2": {
                "cost": 25.20,
                "tokens": 1000000,
                "calls": 500
            },
            "claude-3": {
                "cost": 15.00,
                "tokens": 300000,
                "calls": 100
            }
        },
        "avg_cost_per_instance": 0.10
    }
}
```

### 6.4 获取性能指标

**接口**: `GET /v1/monitoring/metrics`

**查询参数**:
- `metric_name`: 指标名称(可选,不填返回所有)
- `from_time`: 开始时间
- `to_time`: 结束时间

**响应**:
```json
{
    "success": true,
    "data": {
        "metrics": [
            {
                "name": "workflow_duration",
                "values": [
                    {"timestamp": "2025-12-09T10:00:00Z", "value": 165},
                    {"timestamp": "2025-12-09T10:05:00Z", "value": 172}
                ],
                "unit": "seconds"
            },
            {
                "name": "llm_latency",
                "values": [...],
                "unit": "milliseconds"
            }
        ]
    }
}
```

---

## 11. Webhook回调

### 11.1 配置Webhook

**接口**: `POST /v1/webhooks`

**请求体**:
```json
{
    "url": "https://your-server.com/webhook",
    "events": ["workflow.completed", "workflow.failed"],
    "secret": "your_webhook_secret"
}
```

### 7.2 Webhook事件格式

当工作流完成时,系统会POST到配置的URL:

```json
{
    "event": "workflow.completed",
    "timestamp": "2025-12-09T10:02:30Z",
    "data": {
        "instance_id": "wf_20251209_001",
        "status": "completed",
        "result_url": "https://api.metaworkflow.edu/v1/workflows/instances/wf_20251209_001/result"
    },
    "signature": "sha256=..."
}
```

---

## 12. 错误码参考

| 错误码 | HTTP状态 | 说明 |
|--------|---------|------|
| `INVALID_REQUEST` | 400 | 请求参数无效 |
| `UNAUTHORIZED` | 401 | 未认证 |
| `FORBIDDEN` | 403 | 无权限 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `CONFLICT` | 409 | 资源冲突 |
| `RATE_LIMIT_EXCEEDED` | 429 | 超过限流 |
| `WORKFLOW_TIMEOUT` | 500 | 工作流执行超时 |
| `LLM_ERROR` | 502 | LLM服务错误 |
| `INTERNAL_ERROR` | 500 | 内部服务器错误 |

---

## 13. SDK示例

### Python SDK

```python
from metaworkflow import MetaWorkflowClient

client = MetaWorkflowClient(
    api_key="your_api_key",
    base_url="https://api.metaworkflow.edu"
)

# 创建工作流
instance = client.workflows.create(
    subject="数学",
    grade=5,
    topics=["分数加减法"],
    region="CN"
)

# 等待完成
result = instance.wait_for_completion(timeout=300)

# 获取结果
lesson_plan = result.data
print(lesson_plan["title"])

# 导出PDF
pdf_bytes = client.content.export(
    content_id=lesson_plan["lesson_id"],
    format="pdf"
)

with open("lesson.pdf", "wb") as f:
    f.write(pdf_bytes)
```

---

## 附录

### A. 速率限制

| 端点类型 | 限制 |
|---------|------|
| 认证接口 | 10次/分钟 |
| 工作流创建 | 100次/小时 |
| 查询接口 | 1000次/小时 |
| 配置修改 | 50次/小时 |

### B. 变更日志

| 版本 | 日期 | 变更内容 |
|------|------|---------|
| **v2.0** | **2025-12-09** | **🆕 V2.0重大升级**<br>- 新增领域适配器管理API（9个端点）<br>- 新增AI目标生成API（6个端点）<br>- 新增AI内容发现API（7个端点）<br>- 新增AI叙述生成API（8个端点）<br>- 新增多源内容检索API（8个端点）<br>- 支持Milvus向量检索、Redis缓存、MinIO对象存储<br>- 多领域支持（K12、美术史、畅销书、职业培训等） |
| v1.0 | 2025-12-09 | 初始版本<br>- 课程规划API<br>- 模板本地化API<br>- 组件管理API |

---

**文档结束**

**© 2025 MetaWorkflow Platform - V2.0 Multi-Domain AI Content Generation**
