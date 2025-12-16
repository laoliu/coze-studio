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


# 开发环境状态检查

echo "🔍 Coze Studio 开发环境状态"
echo "════════════════════════════════════════"
echo ""

# 检查前端开发服务器
echo "📱 前端开发服务器:"
if curl -s http://localhost:8889 > /dev/null; then
    echo "  ✅ 运行中 - http://localhost:8889"
else
    echo "  ❌ 未运行"
fi
echo ""

# 检查后端服务
echo "🔧 后端服务:"
if docker ps | grep -q coze-server-debug; then
    echo "  ✅ 容器运行中"
    if curl -s http://localhost:8890/health > /dev/null 2>&1; then
        echo "  ✅ API 可访问 - http://localhost:8890"
    else
        echo "  ⚠️  容器运行但 API 可能未就绪"
    fi
else
    echo "  ❌ 容器未运行"
fi
echo ""

# 检查推荐服务 API
echo "💡 推荐服务 API:"
RECOMMEND_TEST=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8890/api/recommend \
    -H "Content-Type: application/json" \
    -d '{"workflowId":"test","sourceNodeId":"test","sourceNodeType":"LLM"}' 2>/dev/null)

if [ "$RECOMMEND_TEST" = "200" ] || [ "$RECOMMEND_TEST" = "400" ] || [ "$RECOMMEND_TEST" = "500" ]; then
    echo "  ✅ API 响应 (HTTP $RECOMMEND_TEST)"
else
    echo "  ❌ API 无响应 (HTTP $RECOMMEND_TEST)"
fi
echo ""

# 检查中间件
echo "🗄️  中间件服务:"
services=("coze-mysql:MySQL" "coze-redis:Redis" "coze-minio:Minio" "coze-elasticsearch:Elasticsearch")
for service in "${services[@]}"; do
    IFS=':' read -r container name <<< "$service"
    if docker ps | grep -q "$container"; then
        health=$(docker inspect --format='{{.State.Health.Status}}' $container 2>/dev/null || echo "running")
        if [ "$health" = "healthy" ]; then
            echo "  ✅ $name (healthy)"
        else
            echo "  ⚠️  $name (running)"
        fi
    else
        echo "  ❌ $name (not running)"
    fi
done
echo ""

# 检查端口占用
echo "🔌 端口状态:"
ports=(8889 8890 3307 6379 9000)
for port in "${ports[@]}"; do
    if lsof -i:$port > /dev/null 2>&1; then
        process=$(lsof -i:$port -t 2>/dev/null | head -1)
        if [ -n "$process" ]; then
            cmd=$(ps -p $process -o comm= 2>/dev/null || echo "unknown")
            echo "  ✅ $port: $cmd"
        else
            echo "  ✅ $port: 已占用"
        fi
    else
        echo "  ❌ $port: 空闲"
    fi
done
echo ""

echo "════════════════════════════════════════"
echo "💡 提示:"
echo "  - 访问应用: http://localhost:8889"
echo "  - 查看日志: docker logs -f coze-server-debug"
echo "  - 停止开发: Ctrl+C (前端) + docker stop coze-server-debug"
echo "════════════════════════════════════════"
