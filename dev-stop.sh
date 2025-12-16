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


# 停止开发环境服务

set -e

echo "🛑 停止开发环境服务..."

# 停止后端
if [ -f "backend.pid" ]; then
    BACKEND_PID=$(cat backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "  停止后端 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
    fi
    rm backend.pid
fi

# 停止前端
if [ -f "frontend.pid" ]; then
    FRONTEND_PID=$(cat frontend.pid)
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "  停止前端 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
    fi
    rm frontend.pid
fi

# 额外保险：通过端口查找并杀死进程
echo "  清理残留进程..."
lsof -ti:8890 | xargs kill -9 2>/dev/null || true
lsof -ti:8889 | xargs kill -9 2>/dev/null || true

echo "✓ 所有服务已停止"
