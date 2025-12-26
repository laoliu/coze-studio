#!/bin/bash

# K12知识图谱管理API测试脚本
# 用法: ./test_k12_admin_api.sh

BASE_URL="http://localhost:8888"
API_BASE="/api/adapter/k12/admin"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  K12 知识图谱管理API 端到端测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_code="${5:-0}"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${YELLOW}[$TOTAL_TESTS] 测试: $test_name${NC}"
    echo -e "   请求: $method $endpoint"

    if [ -z "$data" ]; then
        response=$(curl -s -X "$method" "$BASE_URL$API_BASE$endpoint")
    else
        response=$(curl -s -X "$method" "$BASE_URL$API_BASE$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    # 检查响应
    code=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('code', -1))" 2>/dev/null)

    if [ "$code" = "$expected_code" ]; then
        echo -e "   ${GREEN}✓ 通过${NC} (code=$code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "$response" | python3 -m json.tool 2>/dev/null | head -20
    else
        echo -e "   ${RED}✗ 失败${NC} (expected code=$expected_code, got=$code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo "$response" | python3 -m json.tool 2>/dev/null
    fi
    echo ""
}

# ========================================
# 测试 1: 系统信息API
# ========================================
test_api "获取系统信息" "GET" "/system-info" "" 0

# ========================================
# 测试 2: 获取学科列表
# ========================================
test_api "获取学科列表" "GET" "/subjects" "" 0

# ========================================
# 测试 3: 缓存统计（初始状态）
# ========================================
test_api "获取缓存统计（初始）" "GET" "/cache/stats" "" 0

# ========================================
# 测试 4: 生成知识图谱（首次，LLM调用）
# ========================================
echo -e "${BLUE}--- 关键测试: 首次生成知识图谱 ---${NC}"
echo -e "   预期: 调用LLM生成，耗时3-5秒"
echo -e "   开始时间: $(date '+%H:%M:%S')"

START_TIME=$(date +%s)
test_api "生成知识图谱（数学-7年级）" "POST" "/knowledge-graph/generate" \
    '{"subject": "math", "grade": "grade_7", "force": false}' 0
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo -e "   ${BLUE}耗时: ${DURATION}秒${NC}"
if [ $DURATION -ge 3 ] && [ $DURATION -le 10 ]; then
    echo -e "   ${GREEN}✓ 耗时符合预期（3-10秒，LLM生成）${NC}"
else
    echo -e "   ${YELLOW}⚠ 耗时异常: ${DURATION}秒${NC}"
fi
echo ""

# ========================================
# 测试 5: 再次生成（缓存命中）
# ========================================
echo -e "${BLUE}--- 关键测试: 缓存命中测试 ---${NC}"
echo -e "   预期: 从缓存读取，<100ms"

START_TIME=$(date +%s%3N)
test_api "再次生成知识图谱（应缓存命中）" "POST" "/knowledge-graph/generate" \
    '{"subject": "math", "grade": "grade_7", "force": false}' 0
END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

echo -e "   ${BLUE}耗时: ${DURATION}ms${NC}"
if [ $DURATION -lt 200 ]; then
    echo -e "   ${GREEN}✓ 缓存命中！响应时间<200ms${NC}"
else
    echo -e "   ${YELLOW}⚠ 响应时间较长: ${DURATION}ms${NC}"
fi
echo ""

# ========================================
# 测试 6: 缓存统计（有数据后）
# ========================================
test_api "获取缓存统计（有数据）" "GET" "/cache/stats" "" 0

# ========================================
# 测试 7: 查询知识点列表
# ========================================
test_api "查询知识点列表（数学）" "GET" "/knowledge-points?subject=math&page=1&page_size=10" "" 0

# ========================================
# 测试 8: 查询生成日志
# ========================================
test_api "查询生成日志（最近7天）" "GET" "/generation-logs?days=7&page=1&page_size=20" "" 0

# ========================================
# 测试 9: 刷新指定缓存
# ========================================
test_api "刷新指定缓存" "POST" "/cache/refresh" \
    '{"cache_type": "knowledge_graph", "key": "math:grade_7"}' 0

# ========================================
# 测试 10: 验证缓存失效
# ========================================
echo -e "${BLUE}--- 验证: 缓存失效后重新生成 ---${NC}"
START_TIME=$(date +%s)
test_api "刷新后再次生成（应重新调用LLM）" "POST" "/knowledge-graph/generate" \
    '{"subject": "math", "grade": "grade_7", "force": false}' 0
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo -e "   ${BLUE}耗时: ${DURATION}秒${NC}"
if [ $DURATION -ge 2 ]; then
    echo -e "   ${GREEN}✓ 缓存已失效，重新生成${NC}"
else
    echo -e "   ${YELLOW}⚠ 可能仍从缓存读取${NC}"
fi
echo ""

# ========================================
# 测试 11: 清空所有缓存
# ========================================
test_api "清空所有缓存" "DELETE" "/cache/clear" \
    '{"cache_level": "all"}' 0

# ========================================
# 测试 12: 强制刷新生成
# ========================================
test_api "强制刷新生成（force=true）" "POST" "/knowledge-graph/generate" \
    '{"subject": "physics", "grade": "grade_8", "force": true}' 0

# ========================================
# 测试结果汇总
# ========================================
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  测试结果汇总${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "总测试数: ${TOTAL_TESTS}"
echo -e "${GREEN}通过: ${PASSED_TESTS}${NC}"
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}失败: ${FAILED_TESTS}${NC}"
else
    echo -e "失败: 0"
fi

SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
echo -e "成功率: ${SUCCESS_RATE}%"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}"
    echo "  ✓✓✓ 所有测试通过！ ✓✓✓"
    echo -e "${NC}"
    exit 0
else
    echo -e "${RED}"
    echo "  ✗✗✗ 部分测试失败 ✗✗✗"
    echo -e "${NC}"
    exit 1
fi
