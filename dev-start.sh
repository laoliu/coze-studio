#!/bin/bash
#
# Copyright 2025 coze-dev Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


# 开发环境快速启动脚本
# 不使用 Docker 打包，支持热加载

set -e

echo "🚀 启动 Coze Studio 开发环境..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查中间件是否运行
check_middleware() {
    echo -e "${YELLOW}检查中间件服务...${NC}"

    # 检查 MySQL
    if docker ps | grep -q coze-mysql; then
        echo -e "${GREEN}✓ MySQL 已运行${NC}"
    else
        echo -e "${RED}✗ MySQL 未运行${NC}"
        echo "  启动中间件: cd docker && docker-compose -f docker-compose-debug.yml --profile middleware up -d"
        exit 1
    fi

    # 检查 Redis
    if docker ps | grep -q coze-redis; then
        echo -e "${GREEN}✓ Redis 已运行${NC}"
    else
        echo -e "${RED}✗ Redis 未运行${NC}"
        exit 1
    fi

    echo ""
}

# 启动后端
start_backend() {
    echo -e "${YELLOW}启动后端服务 (端口 8890)...${NC}"

    cd backend

    # 检查 go.mod
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}错误: backend/go.mod 不存在${NC}"
        exit 1
    fi

    # 设置环境变量
    export GOPROXY=https://goproxy.cn,direct
    export GO111MODULE=on

    # 启动后端 (热加载可以使用 air 或 fresh)
    echo "  使用 go run 启动后端..."
    go run -gcflags="all=-N -l" -ldflags="-checklinkname=0" main.go &

    BACKEND_PID=$!
    echo "  后端进程 PID: $BACKEND_PID"
    echo $BACKEND_PID > ../backend.pid

    cd ..
    echo -e "${GREEN}✓ 后端启动完成${NC}"
    echo ""
}

# 启动前端
start_frontend() {
    echo -e "${YELLOW}启动前端服务 (端口 8889)...${NC}"

    cd frontend/apps/coze-studio

    # 检查 package.json
    if [ ! -f "package.json" ]; then
        echo -e "${RED}错误: frontend/apps/coze-studio/package.json 不存在${NC}"
        exit 1
    fi

    # 启动前端开发服务器 (热加载)
    echo "  使用 pnpm dev 启动前端..."
    pnpm dev --port 8889 &

    FRONTEND_PID=$!
    echo "  前端进程 PID: $FRONTEND_PID"
    echo $FRONTEND_PID > ../../../frontend.pid

    cd ../../..
    echo -e "${GREEN}✓ 前端启动完成${NC}"
    echo ""
}

# 主函数
main() {
    echo "════════════════════════════════════════"
    echo "  Coze Studio 开发环境"
    echo "════════════════════════════════════════"
    echo ""

    # 检查中间件
    check_middleware

    # 启动服务
    start_backend
    start_frontend

    echo "════════════════════════════════════════"
    echo -e "${GREEN}✓ 所有服务启动完成！${NC}"
    echo ""
    echo "📝 服务地址:"
    echo "  - 前端: http://localhost:8889"
    echo "  - 后端: http://localhost:8890"
    echo ""
    echo "🔍 查看日志:"
    echo "  - 后端: tail -f backend/logs/*.log"
    echo "  - 前端: 查看终端输出"
    echo ""
    echo "🛑 停止服务:"
    echo "  ./dev-stop.sh"
    echo ""
    echo "💡 代码修改会自动热加载，无需重启"
    echo "════════════════════════════════════════"

    # 等待
    wait
}

# 捕获 Ctrl+C
trap 'echo ""; echo "停止所有服务..."; ./dev-stop.sh; exit 0' INT TERM

main
