#!/bin/bash

# 适配器测试快速脚本
# 用于快速运行所有适配器相关测试

set -e

cd "$(dirname "$0")"

echo "======================================"
echo "  适配器模块测试"
echo "======================================"
echo ""

echo "📋 测试范围:"
echo "  - Validator: 验证器测试"
echo "  - Wrapper: 包装器测试"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 运行 Validator 测试
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 测试 1: Validator 模块"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if go test ./domain/adapter/validator/... -v -cover; then
    echo -e "${GREEN}✅ Validator 测试通过${NC}"
    echo ""
else
    echo -e "${RED}❌ Validator 测试失败${NC}"
    exit 1
fi

# 2. 运行 Wrapper 测试
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 测试 2: Wrapper 模块"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if go test ./domain/adapter/wrapper/... -v -cover; then
    echo -e "${GREEN}✅ Wrapper 测试通过${NC}"
    echo ""
else
    echo -e "${RED}❌ Wrapper 测试失败${NC}"
    exit 1
fi

# 3. 运行所有测试并生成覆盖率报告
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 生成覆盖率报告"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 创建临时目录
mkdir -p coverage

# 生成 Validator 覆盖率
echo "生成 Validator 覆盖率报告..."
go test ./domain/adapter/validator/... -coverprofile=coverage/validator_coverage.out -covermode=atomic
VALIDATOR_COVERAGE=$(go tool cover -func=coverage/validator_coverage.out | grep total | awk '{print $3}')

# 生成 Wrapper 覆盖率
echo "生成 Wrapper 覆盖率报告..."
go test ./domain/adapter/wrapper/... -coverprofile=coverage/wrapper_coverage.out -covermode=atomic
WRAPPER_COVERAGE=$(go tool cover -func=coverage/wrapper_coverage.out | grep total | awk '{print $3}')

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 测试结果总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "${GREEN}✅ 所有测试通过！${NC}"
echo ""
echo "覆盖率统计:"
echo "  • Validator: ${VALIDATOR_COVERAGE}"
echo "  • Wrapper:   ${WRAPPER_COVERAGE}"
echo ""
echo "详细报告:"
echo "  • Validator: coverage/validator_coverage.out"
echo "  • Wrapper:   coverage/wrapper_coverage.out"
echo ""
echo "查看 HTML 报告:"
echo "  go tool cover -html=coverage/validator_coverage.out"
echo "  go tool cover -html=coverage/wrapper_coverage.out"
echo ""
echo -e "${GREEN}🎉 测试完成！${NC}"
