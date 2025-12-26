#!/bin/bash

# K12适配器测试脚本

echo "=== K12教育适配器测试 ==="
echo ""

# 设置基础URL
BASE_URL="http://localhost:8080"

echo "【测试1：分析用户输入】"
echo "输入: 我想学习初中数学的二次函数"
echo ""

# 测试分析接口（假设会添加此接口）
curl -X POST "${BASE_URL}/api/adapter/k12/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习初中数学的二次函数"
  }' | jq '.'

echo ""
echo "【测试2：生成学习计划】"
echo ""

# 测试生成学习计划接口
curl -X POST "${BASE_URL}/api/adapter/k12/generate-plan" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习初中数学的二次函数",
    "mastered_knowledge_ids": ["real_numbers", "equations", "linear_equations_systems"]
  }' | jq '.'

echo ""
echo "【测试3：获取知识点详情】"
echo "知识点ID: quadratic_functions"
echo ""

curl -X GET "${BASE_URL}/api/adapter/k12/knowledge/quadratic_functions" | jq '.'

echo ""
echo "【测试4：生成练习题集】"
echo "知识点: quadratic_functions, 难度: basic, 数量: 5"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/exercises" \
  -H "Content-Type: application/json" \
  -d '{
    "knowledge_id": "quadratic_functions",
    "difficulty": "basic",
    "count": 5
  }' | jq '.'

echo ""
echo "=== 测试完成 ==="
