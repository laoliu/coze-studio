#!/bin/bash

# 物理学科K12适配器测试脚本

echo "=== 物理学科K12适配器测试 ==="
echo ""

# 设置基础URL
BASE_URL="http://localhost:8888"

echo "【测试1：分析物理输入 - 欧姆定律】"
echo "输入: 我想学习初中物理的欧姆定律"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习初中物理的欧姆定律"
  }' 2>/dev/null | jq '.'

echo ""
echo "【测试2：生成物理学习计划 - 欧姆定律】"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/generate-plan" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习初中物理的欧姆定律",
    "mastered_knowledge_ids": ["electricity_basics", "voltage_resistance"]
  }' 2>/dev/null | jq '.'

echo ""
echo "【测试3：获取物理知识点详情 - 欧姆定律】"
echo ""

curl -X GET "${BASE_URL}/api/adapter/k12/knowledge/ohms_law" 2>/dev/null | jq '.'

echo ""
echo "【测试4：获取物理习题 - 欧姆定律基础题】"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/exercises" \
  -H "Content-Type: application/json" \
  -d '{
    "knowledge_id": "ohms_law",
    "difficulty": "basic",
    "count": 5
  }' 2>/dev/null | jq '.'

echo ""
echo "【测试5：分析物理输入 - 电功率】"
echo "输入: 我想学习九年级物理的电功率"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习九年级物理的电功率"
  }' 2>/dev/null | jq '.'

echo ""
echo "【测试6：获取物理知识点详情 - 运动和力】"
echo ""

curl -X GET "${BASE_URL}/api/adapter/k12/knowledge/motion_and_forces" 2>/dev/null | jq '.'

echo ""
echo "【测试7：生成物理学习计划 - 电功率（完全从零开始）】"
echo ""

curl -X POST "${BASE_URL}/api/adapter/k12/generate-plan" \
  -H "Content-Type: application/json" \
  -d '{
    "user_input": "我想学习初中物理的电功率",
    "mastered_knowledge_ids": []
  }' 2>/dev/null | jq '.'

echo ""
echo "=== 物理学科测试完成 ==="
echo ""
echo "📊 已测试功能："
echo "  ✅ 物理学科识别"
echo "  ✅ 物理知识点分析（欧姆定律、电功率、运动和力）"
echo "  ✅ 物理学习计划生成"
echo "  ✅ 物理知识点详情查询"
echo "  ✅ 物理习题获取"
echo ""
echo "🎯 物理知识图谱包含："
echo "  • 力学：速度、力、运动和力、压强、浮力、功和机械能、简单机械"
echo "  • 电学：电流和电路、电压和电阻、欧姆定律、电功率、电路分析"
echo "  • 光学：光的反射、光的折射"
echo "  • 总计：16个知识点"
