#!/bin/bash

# Docker 容器代理配置快速修复脚本

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}   Docker 容器代理快速配置${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# 获取 Windows 主机 IP
WIN_HOST_IP=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}')
echo -e "${GREEN}✓ Windows 主机 IP: $WIN_HOST_IP${NC}"

# 询问代理端口
read -p "请输入代理端口 [默认: 7890]: " PROXY_PORT
PROXY_PORT=${PROXY_PORT:-7890}
echo -e "${GREEN}✓ 代理端口: $PROXY_PORT${NC}"
echo ""

# 测试 WSL 代理是否可用
echo -e "${YELLOW}🔍 测试 WSL 代理连接...${NC}"
if timeout 5 curl -s -x "http://${WIN_HOST_IP}:${PROXY_PORT}" https://www.google.com > /dev/null 2>&1; then
    echo -e "${GREEN}✅ WSL 代理连接正常${NC}"
else
    echo -e "${RED}❌ WSL 代理连接失败${NC}"
    echo "请确保 Windows 代理工具正在运行并开启了'允许局域网连接'"
    exit 1
fi
echo ""

# 方法 1: 使用环境变量（临时，重启容器后失效）
echo -e "${BLUE}方法 1: 临时设置容器代理（立即生效，重启后失效）${NC}"
echo ""

PROXY_URL="http://${WIN_HOST_IP}:${PROXY_PORT}"

echo "设置环境变量到容器..."

# 为运行中的容器设置代理（需要重启应用进程）
docker exec coze-server sh -c "
export HTTP_PROXY='$PROXY_URL'
export HTTPS_PROXY='$PROXY_URL'
export http_proxy='$PROXY_URL'
export https_proxy='$PROXY_URL'
export NO_PROXY='localhost,127.0.0.1,mysql,redis,nats,minio'
echo '✅ 代理环境变量已设置（当前会话）'
env | grep -i proxy
" 2>/dev/null || true

echo ""
echo -e "${YELLOW}⚠️  注意: 方法 1 只在当前会话生效，容器重启后失效${NC}"
echo ""

# 方法 2: 修改 docker-compose.yml（永久）
echo -e "${BLUE}方法 2: 修改 docker-compose.yml（永久生效，推荐）${NC}"
echo ""

COMPOSE_FILE="/home/liunix/work/coze-studio/docker/docker-compose.yml"

if [ -f "$COMPOSE_FILE" ]; then
    # 备份
    BACKUP_FILE="${COMPOSE_FILE}.backup.proxy.$(date +%Y%m%d_%H%M%S)"
    cp "$COMPOSE_FILE" "$BACKUP_FILE"
    echo -e "${GREEN}✓ 已备份到: $BACKUP_FILE${NC}"
    
    # 检查是否已有代理配置
    if grep -q "HTTP_PROXY.*http://" "$COMPOSE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  检测到已有代理配置，跳过自动添加${NC}"
        echo "请手动编辑文件检查配置"
    else
        echo ""
        echo -e "${YELLOW}请在 coze-server 服务的 environment 部分添加:${NC}"
        echo ""
        echo -e "${GREEN}      - HTTP_PROXY=http://${WIN_HOST_IP}:${PROXY_PORT}${NC}"
        echo -e "${GREEN}      - HTTPS_PROXY=http://${WIN_HOST_IP}:${PROXY_PORT}${NC}"
        echo -e "${GREEN}      - http_proxy=http://${WIN_HOST_IP}:${PROXY_PORT}${NC}"
        echo -e "${GREEN}      - https_proxy=http://${WIN_HOST_IP}:${PROXY_PORT}${NC}"
        echo -e "${GREEN}      - NO_PROXY=localhost,127.0.0.1,mysql,redis,nats,minio${NC}"
        echo ""
        
        read -p "是否尝试自动添加? (y/n) [n]: " AUTO_ADD
        
        if [ "$AUTO_ADD" = "y" ]; then
            # 尝试在 coze-server 的 environment 部分添加
            # 查找 coze-server 服务的 environment 部分
            if grep -A 20 "coze-server:" "$COMPOSE_FILE" | grep -q "environment:"; then
                echo "正在添加代理配置..."
                
                # 使用 sed 在 environment 部分后添加
                sed -i "/coze-server:/,/environment:/{
                    /environment:/a\\
      - HTTP_PROXY=http://${WIN_HOST_IP}:${PROXY_PORT}\\
      - HTTPS_PROXY=http://${WIN_HOST_IP}:${PROXY_PORT}\\
      - http_proxy=http://${WIN_HOST_IP}:${PROXY_PORT}\\
      - https_proxy=http://${WIN_HOST_IP}:${PROXY_PORT}\\
      - NO_PROXY=localhost,127.0.0.1,mysql,redis,nats,minio
                }" "$COMPOSE_FILE" 2>/dev/null || {
                    echo -e "${RED}自动添加失败，请手动编辑${NC}"
                }
                
                echo -e "${GREEN}✓ 已尝试添加配置，请手动验证${NC}"
            else
                echo -e "${YELLOW}未找到 environment 部分，请手动添加${NC}"
            fi
        fi
    fi
    
    echo ""
    echo -e "${YELLOW}编辑文件:${NC}"
    echo "  nano $COMPOSE_FILE"
    echo ""
    
    read -p "配置完成后按 Enter 继续..."
else
    echo -e "${RED}✗ 未找到 docker-compose.yml${NC}"
fi

echo ""
echo -e "${BLUE}重启 Docker 服务以应用配置...${NC}"
echo ""

read -p "是否立即重启 coze-server 容器? (y/n) [y]: " RESTART
RESTART=${RESTART:-y}

if [ "$RESTART" = "y" ]; then
    cd /home/liunix/work/coze-studio
    
    echo "停止容器..."
    docker-compose -f docker/docker-compose.yml stop coze-server
    
    echo "启动容器..."
    docker-compose -f docker/docker-compose.yml up -d coze-server
    
    echo "等待容器启动..."
    sleep 5
    
    echo ""
    echo -e "${GREEN}✓ 容器已重启${NC}"
    
    # 验证代理配置
    echo ""
    echo -e "${YELLOW}🔍 验证容器内代理配置...${NC}"
    docker exec coze-server env | grep -i proxy || echo "未找到代理环境变量"
    
    echo ""
    echo -e "${YELLOW}🔍 测试容器内网络连接...${NC}"
    docker exec coze-server sh -c "curl -s -m 5 https://www.google.com > /dev/null && echo '✅ Google.com 可访问' || echo '❌ Google.com 不可访问'"
fi

echo ""
echo -e "${BLUE}================================================${NC}"
echo -e "${GREEN}   配置完成！${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

echo "验证步骤："
echo ""
echo "1. 检查容器代理环境变量:"
echo "   docker exec coze-server env | grep -i proxy"
echo ""
echo "2. 在容器内测试 Google:"
echo "   docker exec coze-server curl -s https://www.google.com"
echo ""
echo "3. 测试 Google API:"
echo "   docker exec coze-server curl -s https://www.googleapis.com"
echo ""
echo "4. 在 Coze 中测试 Google Search Tool"
echo ""

# 创建快速测试脚本
cat > test_docker_proxy.sh << 'EOF'
#!/bin/bash

echo "🔍 测试 Docker 容器代理..."
echo ""

echo "1. 检查代理环境变量:"
docker exec coze-server env | grep -i proxy
echo ""

echo "2. 测试 Google.com:"
docker exec coze-server sh -c "curl -s -m 5 https://www.google.com > /dev/null && echo '✅ 成功' || echo '❌ 失败'"
echo ""

echo "3. 测试 Google APIs:"
docker exec coze-server sh -c "curl -s -m 5 https://www.googleapis.com > /dev/null && echo '✅ 成功' || echo '❌ 失败'"
echo ""

echo "4. 测试 Google Search API:"
if [ -n "$GOOGLE_API_KEY" ] && [ -n "$GOOGLE_CX" ]; then
    docker exec coze-server sh -c "curl -s -m 10 'https://www.googleapis.com/customsearch/v1?key=$GOOGLE_API_KEY&cx=$GOOGLE_CX&q=test' && echo '✅ Google Search API 可用' || echo '❌ Google Search API 不可用'"
else
    echo "⚠️  未设置 GOOGLE_API_KEY 和 GOOGLE_CX，跳过"
fi
EOF

chmod +x test_docker_proxy.sh

echo -e "${GREEN}✓ 已创建测试脚本: test_docker_proxy.sh${NC}"
echo ""
