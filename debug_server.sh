#!/bin/bash
# =============================================================================
# Coze Server 调试助手脚本
# =============================================================================
# 用法: ./debug_server.sh [command]
# 命令:
#   middleware  - 只启动中间件服务
#   stop        - 停止 coze-server 容器
#   logs        - 查看 coze-server 日志
#   attach      - 进入 coze-server 容器
#   status      - 查看所有服务状态
#   clean       - 停止所有服务
#   help        - 显示帮助信息
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目路径
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_DIR="${PROJECT_ROOT}/docker"
COMPOSE_FILE="${DOCKER_DIR}/docker-compose.yml"
DEBUG_COMPOSE_FILE="${DOCKER_DIR}/docker-compose-debug.yml"
ENV_FILE="${DOCKER_DIR}/.env"
DEBUG_ENV_FILE="${DOCKER_DIR}/.env.debug"

# 打印信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# 启动中间件服务
start_middleware() {
    print_header "启动中间件服务"
    
    # 检查配置文件
    if [ ! -f "$DEBUG_ENV_FILE" ]; then
        print_warning "配置文件 $DEBUG_ENV_FILE 不存在，从示例创建..."
        cp "${DEBUG_ENV_FILE}.example" "$DEBUG_ENV_FILE"
    fi
    
    print_info "启动 MySQL, Redis, Elasticsearch, MinIO, Etcd, Milvus, NSQ..."
    
    docker compose -f "$DEBUG_COMPOSE_FILE" \
        --env-file "$DEBUG_ENV_FILE" \
        --profile middleware up -d --wait
    
    print_success "中间件服务启动完成！"
    echo ""
    print_info "服务端口映射："
    echo "  - MySQL:         localhost:3306"
    echo "  - Redis:         localhost:6379"
    echo "  - Elasticsearch: localhost:9200"
    echo "  - MinIO:         localhost:9000 (API), localhost:9001 (Console)"
    echo "  - Milvus:        localhost:19530 (gRPC), localhost:9091 (HTTP)"
    echo "  - NSQ:           localhost:4150 (TCP), localhost:4151 (HTTP)"
    echo ""
    print_info "现在可以在本地启动 coze-server 进行调试了："
    echo "  cd backend && go run main.go"
    echo ""
    echo "或者使用 Make 命令："
    echo "  make server"
}

# 停止 coze-server 容器
stop_server() {
    print_header "停止 coze-server 容器"
    
    if docker ps -a --format '{{.Names}}' | grep -q '^coze-server$'; then
        print_info "停止并删除 coze-server 容器..."
        docker stop coze-server 2>/dev/null || true
        docker rm coze-server 2>/dev/null || true
        print_success "coze-server 容器已停止并删除"
    else
        print_warning "coze-server 容器不存在或未运行"
    fi
}

# 查看日志
show_logs() {
    print_header "查看 coze-server 日志"
    
    if docker ps --format '{{.Names}}' | grep -q '^coze-server$'; then
        print_info "实时查看日志（按 Ctrl+C 退出）..."
        docker logs -f coze-server
    else
        print_error "coze-server 容器未运行"
        exit 1
    fi
}

# 进入容器
attach_container() {
    print_header "进入 coze-server 容器"
    
    if docker ps --format '{{.Names}}' | grep -q '^coze-server$'; then
        print_info "进入容器 shell..."
        docker exec -it coze-server /bin/bash
    else
        print_error "coze-server 容器未运行"
        exit 1
    fi
}

# 查看状态
show_status() {
    print_header "服务状态"
    
    echo ""
    print_info "Docker 容器状态："
    docker ps -a --filter "name=coze-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    echo ""
    print_info "健康检查状态："
    
    containers=("coze-mysql" "coze-redis" "coze-elasticsearch" "coze-minio" "coze-milvus" "coze-server")
    
    for container in "${containers[@]}"; do
        if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
            health=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "unknown")
            if [ "$health" == "healthy" ]; then
                echo -e "  ${GREEN}✓${NC} $container: healthy"
            elif [ "$health" == "starting" ]; then
                echo -e "  ${YELLOW}⏳${NC} $container: starting"
            elif [ "$health" == "unhealthy" ]; then
                echo -e "  ${RED}✗${NC} $container: unhealthy"
            else
                echo -e "  ${BLUE}ℹ${NC} $container: no health check"
            fi
        else
            echo -e "  ${RED}✗${NC} $container: not running"
        fi
    done
    
    echo ""
    print_info "端口监听状态："
    
    ports=(3306 6379 9200 9000 19530 8888)
    port_names=("MySQL" "Redis" "Elasticsearch" "MinIO" "Milvus" "Coze-Server")
    
    for i in "${!ports[@]}"; do
        port="${ports[$i]}"
        name="${port_names[$i]}"
        if lsof -i ":$port" > /dev/null 2>&1 || netstat -tlnp 2>/dev/null | grep -q ":$port "; then
            echo -e "  ${GREEN}✓${NC} $name (port $port): listening"
        else
            echo -e "  ${RED}✗${NC} $name (port $port): not listening"
        fi
    done
}

# 清理所有服务
clean_all() {
    print_header "清理所有服务"
    
    read -p "确定要停止所有服务吗？(y/N): " confirm
    
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        print_info "操作已取消"
        exit 0
    fi
    
    print_info "停止 docker-compose.yml 的所有服务..."
    docker compose -f "$COMPOSE_FILE" down 2>/dev/null || true
    
    print_info "停止 docker-compose-debug.yml 的所有服务..."
    docker compose -f "$DEBUG_COMPOSE_FILE" --profile '*' down 2>/dev/null || true
    
    print_success "所有服务已停止"
    
    read -p "是否删除数据卷？(y/N): " confirm_data
    
    if [ "$confirm_data" == "y" ] || [ "$confirm_data" == "Y" ]; then
        print_warning "删除数据目录..."
        rm -rf "${DOCKER_DIR}/data"
        print_success "数据已清理"
    fi
}

# 快速调试模式
debug_mode() {
    print_header "快速调试模式"
    
    print_info "步骤 1/3: 停止现有的 coze-server 容器..."
    stop_server
    
    print_info "步骤 2/3: 启动中间件服务..."
    start_middleware
    
    print_info "步骤 3/3: 准备本地调试环境..."
    
    echo ""
    print_success "调试环境准备完成！"
    echo ""
    print_info "下一步操作："
    echo ""
    echo "方式 1: 使用 Make 命令"
    echo "  cd ${PROJECT_ROOT}"
    echo "  make server"
    echo ""
    echo "方式 2: 直接运行 Go"
    echo "  cd ${PROJECT_ROOT}/backend"
    echo "  export APP_ENV=debug"
    echo "  go run main.go"
    echo ""
    echo "方式 3: 使用 IDE 调试"
    echo "  在 VSCode 或 GoLand 中设置断点"
    echo "  按 F5 或点击 Debug 按钮"
    echo ""
    echo "访问地址："
    echo "  - API: http://localhost:8888"
    echo "  - Debug: http://localhost:8889"
}

# 测试连接
test_connections() {
    print_header "测试中间件连接"
    
    # 测试 MySQL
    print_info "测试 MySQL 连接..."
    if docker exec coze-mysql mysql -ucoze -pcoze123 -e "SELECT 1" > /dev/null 2>&1; then
        print_success "MySQL 连接成功"
    else
        print_error "MySQL 连接失败"
    fi
    
    # 测试 Redis
    print_info "测试 Redis 连接..."
    if docker exec coze-redis redis-cli ping > /dev/null 2>&1; then
        print_success "Redis 连接成功"
    else
        print_error "Redis 连接失败"
    fi
    
    # 测试 Elasticsearch
    print_info "测试 Elasticsearch 连接..."
    if curl -s http://localhost:9200 > /dev/null 2>&1; then
        print_success "Elasticsearch 连接成功"
    else
        print_error "Elasticsearch 连接失败"
    fi
    
    # 测试 MinIO
    print_info "测试 MinIO 连接..."
    if curl -s http://localhost:9000/minio/health/live > /dev/null 2>&1; then
        print_success "MinIO 连接成功"
    else
        print_error "MinIO 连接失败"
    fi
}

# 显示帮助
show_help() {
    cat << EOF
Coze Server 调试助手

用法: $0 [command]

命令:
  middleware    只启动中间件服务（MySQL, Redis, ES等）
  debug         快速调试模式（停止server容器 + 启动中间件）
  stop          停止 coze-server 容器
  logs          查看 coze-server 日志
  attach        进入 coze-server 容器
  status        查看所有服务状态
  test          测试中间件连接
  clean         停止所有服务并清理
  help          显示此帮助信息

示例:
  # 启动调试环境
  $0 debug
  
  # 查看服务状态
  $0 status
  
  # 查看日志
  $0 logs
  
  # 清理环境
  $0 clean

更多信息请查看: docs/DEBUG_COZE_SERVER.md
EOF
}

# 主函数
main() {
    check_docker
    
    case "${1:-help}" in
        middleware)
            start_middleware
            ;;
        debug)
            debug_mode
            ;;
        stop)
            stop_server
            ;;
        logs)
            show_logs
            ;;
        attach)
            attach_container
            ;;
        status)
            show_status
            ;;
        test)
            test_connections
            ;;
        clean)
            clean_all
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
