#!/bin/bash

# Coze Studio Adapter API 集成测试脚本

BASE_URL="http://localhost:8888"
API_PREFIX="/api/adapter"

echo "=========================================="
echo "  Coze Studio Adapter API 集成测试"
echo "=========================================="
echo ""

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4

    echo -e "${YELLOW}测试: ${name}${NC}"
    echo "URL: ${method} ${BASE_URL}${endpoint}"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -d "${data}" \
            "${BASE_URL}${endpoint}")
    fi

    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    body=$(echo "$response" | sed '/HTTP_CODE:/d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ 成功 (HTTP ${http_code})${NC}"
        echo "响应: $(echo $body | jq -C '.' 2>/dev/null || echo $body)"
    else
        echo -e "${RED}✗ 失败 (HTTP ${http_code})${NC}"
        echo "响应: $body"
    fi
    echo ""
}

# 1. 测试获取适配器列表
test_api "获取适配器列表" "GET" "${API_PREFIX}/list"

# 2. 测试获取适配器列表（带分页）
test_api "获取适配器列表（分页）" "GET" "${API_PREFIX}/list?page=1&page_size=10"

# 3. 测试获取适配器详情
test_api "获取适配器详情" "GET" "${API_PREFIX}/k12_education"

# 4. 测试注册适配器
read -r -d '' REGISTER_DATA << 'EOF'
{
  "adapter_id": "test_adapter_001",
  "name": "测试适配器",
  "description": "用于集成测试的适配器",
  "version": "1.0.0",
  "adapter_type": "education",
  "author": {
    "name": "Test User",
    "email": "test@example.com",
    "organization": "Test Org"
  },
  "supported_activity_types": ["concept", "experiment"],
  "supported_domains": ["mathematics", "physics"],
  "supported_grades": ["grade_7", "grade_8"]
}
EOF

test_api "注册新适配器" "POST" "${API_PREFIX}/register" "$REGISTER_DATA"

# 5. 测试执行适配器
read -r -d '' EXECUTE_DATA << 'EOF'
{
  "input": {
    "topic": "氧化还原反应",
    "domain": "化学",
    "grade": "高一",
    "duration": 45,
    "activity_type": "concept"
  },
  "context": {
    "user_id": 12345,
    "language": "zh-CN"
  }
}
EOF

test_api "执行适配器" "POST" "${API_PREFIX}/k12_education/execute" "$EXECUTE_DATA"

# 6. 测试获取已安装的适配器
test_api "获取已安装的适配器" "GET" "${API_PREFIX}/installed"

echo "=========================================="
echo "  测试完成"
echo "=========================================="
