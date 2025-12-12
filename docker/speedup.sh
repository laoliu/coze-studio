#!/bin/bash
# Docker 启动加速脚本
# 用途：预构建包含 Atlas 的 MySQL 镜像，大幅提升启动速度

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CUSTOM_IMAGE="coze-mysql-atlas:8.4.5"

echo "=========================================="
echo "  Coze Studio Docker 启动加速工具"
echo "=========================================="
echo ""

# 检查是否已存在自定义镜像
if docker images | grep -q "coze-mysql-atlas"; then
    echo "✅ 发现已存在的优化镜像: $CUSTOM_IMAGE"
    read -p "是否重新构建？(y/N): " rebuild
    if [[ ! "$rebuild" =~ ^[Yy]$ ]]; then
        echo "跳过构建，直接启动..."
        cd "$SCRIPT_DIR"
        docker compose -f docker-compose-local.yml up -d
        exit 0
    fi
fi

echo ""
echo "📦 步骤 1: 构建优化的 MySQL 镜像（包含 Atlas CLI）"
echo "这一步只需要执行一次，大约需要 3-5 分钟..."
echo ""

cd "$SCRIPT_DIR"
docker build -f mysql-with-atlas.Dockerfile -t "$CUSTOM_IMAGE" .

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 镜像构建成功！"
    echo ""
else
    echo ""
    echo "❌ 镜像构建失败，请检查错误信息"
    exit 1
fi

echo "📝 步骤 2: 备份原 docker-compose-local.yml"
if [ ! -f docker-compose-local.yml.backup ]; then
    cp docker-compose-local.yml docker-compose-local.yml.backup
    echo "✅ 已创建备份: docker-compose-local.yml.backup"
else
    echo "⚠️  备份已存在，跳过"
fi

echo ""
echo "🔧 步骤 3: 修改 docker-compose-local.yml"
echo "将 MySQL 镜像从 mysql:8.4.5 改为 $CUSTOM_IMAGE"
echo ""

# 使用 sed 替换镜像
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|image: mysql:8.4.5|image: $CUSTOM_IMAGE|g" docker-compose-local.yml
else
    # Linux
    sed -i "s|image: mysql:8.4.5|image: $CUSTOM_IMAGE|g" docker-compose-local.yml
fi

# 注释掉 Atlas 安装代码（因为已经预装了）
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' '/if ! command -v atlas/,/fi/s/^/        # /' docker-compose-local.yml
else
    sed -i '/if ! command -v atlas/,/fi/s/^/        # /' docker-compose-local.yml
fi

echo "✅ 配置文件已更新"
echo ""

echo "🚀 步骤 4: 启动 Docker Compose"
echo ""

docker compose -f docker-compose-local.yml up -d

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "  ✅ 启动成功！"
    echo "=========================================="
    echo ""
    echo "以后启动只需要运行："
    echo "  cd $SCRIPT_DIR"
    echo "  docker compose -f docker-compose-local.yml up -d"
    echo ""
    echo "查看日志："
    echo "  docker compose -f docker-compose-local.yml logs -f"
    echo ""
    echo "停止服务："
    echo "  docker compose -f docker-compose-local.yml down"
    echo ""
    echo "如需恢复原配置："
    echo "  cp docker-compose-local.yml.backup docker-compose-local.yml"
    echo ""
else
    echo ""
    echo "❌ 启动失败，请检查日志"
    echo "查看日志命令: docker compose -f docker-compose-local.yml logs"
fi
