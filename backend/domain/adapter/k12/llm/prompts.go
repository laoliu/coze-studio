/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llm

// KnowledgeGraphPromptTemplate 知识图谱生成提示词模板
const KnowledgeGraphPromptTemplate = `你是K12教育专家，负责构建{{.Subject}}学科{{.Grade}}年级的知识图谱。

任务：生成该学科该年级的15-20个核心知识点，覆盖完整的知识体系。

每个知识点包含：
- id：英文标识符（小写下划线，如"quadratic_functions"）
- name：中文名称
- description：简短描述（1句话）
- difficulty：难度（basic/improve/advanced）
- keywords：关键词列表（3个）
- prerequisites：前置知识点ID列表
- next_points：后续知识点ID列表

严格按此JSON格式输出（不要添加markdown标记或其他文字）：
{
  "subject": "{{.Subject}}",
  "grade": "{{.Grade}}",
  "knowledge_points": [
    {
      "id": "example_id",
      "name": "示例名称",
      "description": "示例描述",
      "difficulty": "basic",
      "keywords": ["关键词1", "关键词2", "关键词3"],
      "prerequisites": [],
      "next_points": []
    }
  ]
}

请立即输出完整的JSON：`

// ResourceGenerationPromptTemplate 学习资源生成提示词模板
const ResourceGenerationPromptTemplate = `你是一位K12教育资源策划专家，擅长为学生推荐优质学习资源。

任务：为以下知识点生成3-5个优质学习资源推荐。

知识点信息：
- 知识点ID: {{.KnowledgeID}}
- 知识点名称: {{.KnowledgeName}}
- 学科: {{.Subject}}
- 年级: {{.Grade}}
- 描述: {{.Description}}

资源要求：
1. 类型多样性：
   - video: 教学视频（讲解清晰，动画生动）
   - textbook: 教材章节（权威出版社）
   - article: 文章资料（深入浅出）
   - interactive: 互动课件（趣味性强）

2. 难度层次：
   - basic: 基础入门（适合初学者）
   - improve: 提高巩固（适合已有基础）
   - advanced: 拓展深化（适合学有余力）

3. 资源信息：
   - title: 资源标题（准确描述内容）
   - description: 资源描述（说明特色和适用场景）
   - url: 资源链接（使用示例链接，格式为 https://example.com/...）
   - duration: 时长（视频为分钟数，文章为预估阅读分钟数）
   - source: 来源（如"人教版教材"、"中学数学在线"等）

4. 内容质量：
   - 符合中国教学大纲
   - 内容准确可靠
   - 表达清晰易懂
   - 有助于理解核心概念

输出格式（严格遵守JSON格式）：
{
  "knowledge_id": "{{.KnowledgeID}}",
  "resources": [
    {
      "id": "resource_{{.KnowledgeID}}_video_1",
      "type": "video",
      "title": "资源标题",
      "description": "资源描述，说明内容特色",
      "url": "https://example.com/resource",
      "difficulty": "basic",
      "duration": 15,
      "source": "来源名称"
    }
  ]
}

请生成3-5个学习资源JSON（直接输出JSON，不要添加其他说明文字）：`

// ExerciseGenerationPromptTemplate 习题生成提示词模板
const ExerciseGenerationPromptTemplate = `你是一位K12教育习题设计专家，擅长设计具有针对性和代表性的练习题。

任务：为以下知识点生成{{.Count}}道练习题。

知识点信息：
- 知识点ID: {{.KnowledgeID}}
- 知识点名称: {{.KnowledgeName}}
- 学科: {{.Subject}}
- 年级: {{.Grade}}
- 目标难度: {{.Difficulty}}

习题要求：
1. 题型多样：
   - 选择题：提供4个选项（A/B/C/D）
   - 填空题：options为空数组
   - 计算题：提供详细步骤
   - 解答题：注重思维过程

2. 难度分级（{{.Difficulty}}）：
{{if eq .Difficulty "basic"}}   - basic（基础）：考查基本概念和简单计算，正确率应达到80%以上
{{else if eq .Difficulty "improve"}}   - improve（提高）：考查综合应用和问题分析，有一定思维含量
{{else if eq .Difficulty "advanced"}}   - advanced（拓展）：考查深度理解和创新思维，具有挑战性
{{end}}

3. 题目要素：
   - question: 题目内容（表述清晰，条件完整）
   - options: 选项数组（选择题必填，其他题型为空）
   - answer: 正确答案（格式规范）
   - solution: 详细解析（包含解题思路和关键步骤）
   - score: 分值（5分、8分、10分、15分等）

4. 质量标准：
   - 答案准确无误
   - 解析详细清晰
   - 覆盖核心知识点
   - 具有教学价值

输出格式（严格遵守JSON格式）：
{
  "knowledge_id": "{{.KnowledgeID}}",
  "difficulty": "{{.Difficulty}}",
  "exercises": [
    {
      "id": "{{.KnowledgeID}}_ex_{{.Difficulty}}_1",
      "question": "题目内容",
      "options": ["A. 选项1", "B. 选项2", "C. 选项3", "D. 选项4"],
      "answer": "正确答案",
      "solution": "详细解题步骤和思路说明",
      "difficulty": "{{.Difficulty}}",
      "score": 5
    }
  ]
}

请生成{{.Count}}道练习题JSON（直接输出JSON，不要添加其他说明文字）：`

// KnowledgePointPromptTemplate 单个知识点生成提示词模板
const KnowledgePointPromptTemplate = `你是一位K12教育专家，负责生成单个知识点的详细信息。

任务：生成以下知识点的完整信息。

输入信息：
- 知识点ID: {{.KnowledgeID}}
- 知识点名称: {{.KnowledgeName}}
- 学科: {{.Subject}}
- 年级: {{.Grade}}

要求：
1. 提供准确的知识点描述（1-2句话）
2. 给出3-5个核心关键词
3. 识别前置知识点（学习该知识点前需要掌握的内容）
4. 推荐后续知识点（学习该知识点后可以延伸的内容）
5. 判断难度等级（basic/improve/advanced）

输出格式（严格遵守JSON格式）：
{
  "id": "{{.KnowledgeID}}",
  "name": "{{.KnowledgeName}}",
  "subject": "{{.Subject}}",
  "grade": "{{.Grade}}",
  "description": "知识点描述",
  "difficulty": "basic",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "prerequisites": ["前置知识点id1", "前置知识点id2"],
  "next_points": ["后续知识点id1", "后续知识点id2"]
}

请生成该知识点的完整信息JSON（直接输出JSON，不要添加其他说明文字）：`
