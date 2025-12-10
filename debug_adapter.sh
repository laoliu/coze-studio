#!/bin/bash
# =============================================================================
# Adapter Demo 快速调试脚本
# =============================================================================
# 用法: ./debug_adapter.sh [command]
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
DEMO_DIR="${PROJECT_ROOT}/backend/cmd/adapter-demo"
BACKEND_DIR="${PROJECT_ROOT}/backend"

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

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        print_warning "Go 命令未找到，正在检查常见安装位置..."
        
        # 检查常见的 Go 安装路径
        possible_paths=(
            "/usr/local/go/bin/go"
            "/usr/bin/go"
            "$HOME/go/bin/go"
            "$HOME/.go/bin/go"
        )
        
        for go_path in "${possible_paths[@]}"; do
            if [ -x "$go_path" ]; then
                print_success "找到 Go: $go_path"
                export PATH="$(dirname $go_path):$PATH"
                break
            fi
        done
        
        if ! command -v go &> /dev/null; then
            print_error "Go 未安装。请先安装 Go:"
            echo "  https://golang.org/dl/"
            echo ""
            echo "或者在 Ubuntu/Debian 上:"
            echo "  sudo apt update"
            echo "  sudo apt install golang-go"
            exit 1
        fi
    fi
    
    go_version=$(go version | awk '{print $3}')
    print_info "Go version: $go_version"
}

# 运行 K12 演示
run_k12() {
    print_header "运行 K12 适配器演示"
    
    cd "$BACKEND_DIR"
    
    print_info "执行: go run ./cmd/adapter-demo -demo k12"
    echo ""
    
    go run ./cmd/adapter-demo -demo k12
    
    echo ""
    print_success "K12 演示完成！"
}

# 运行注册表演示
run_registry() {
    print_header "运行适配器注册表演示"
    
    cd "$BACKEND_DIR"
    
    print_info "执行: go run ./cmd/adapter-demo -demo registry"
    echo ""
    
    go run ./cmd/adapter-demo -demo registry
    
    echo ""
    print_success "注册表演示完成！"
}

# 运行所有演示
run_all() {
    print_header "运行所有演示"
    
    cd "$BACKEND_DIR"
    
    print_info "执行: go run ./cmd/adapter-demo -demo all"
    echo ""
    
    go run ./cmd/adapter-demo -demo all
    
    echo ""
    print_success "所有演示完成！"
}

# 运行测试
run_tests() {
    print_header "运行单元测试"
    
    cd "$BACKEND_DIR/domain/adapter/examples"
    
    print_info "执行: go test -v"
    echo ""
    
    go test -v
    
    echo ""
    print_success "测试完成！"
}

# 运行性能测试
run_bench() {
    print_header "运行性能测试"
    
    cd "$BACKEND_DIR/domain/adapter/examples"
    
    print_info "执行: go test -v -bench=. -benchmem"
    echo ""
    
    go test -v -bench=. -benchmem
    
    echo ""
    print_success "性能测试完成！"
}

# 使用 Delve 调试
debug_delve() {
    print_header "使用 Delve 调试"
    
    # 检查 delve 是否安装
    if ! command -v dlv &> /dev/null; then
        print_warning "Delve 未安装，正在安装..."
        go install github.com/go-delve/delve/cmd/dlv@latest
        print_success "Delve 安装完成"
    fi
    
    cd "$BACKEND_DIR"
    
    demo_type="${1:-k12}"
    
    print_info "启动 Delve 调试器"
    print_info "演示类型: $demo_type"
    echo ""
    print_info "常用命令:"
    echo "  break demo.go:15  - 设置断点"
    echo "  continue          - 继续执行"
    echo "  next              - 单步跳过"
    echo "  step              - 单步进入"
    echo "  print var         - 打印变量"
    echo "  list              - 显示代码"
    echo "  exit              - 退出"
    echo ""
    
    dlv debug ./cmd/adapter-demo -- -demo "$demo_type"
}

# 构建可执行文件
build_demo() {
    print_header "构建演示程序"
    
    cd "$BACKEND_DIR"
    
    print_info "执行: go build -o adapter-demo ./cmd/adapter-demo"
    
    go build -o adapter-demo ./cmd/adapter-demo
    
    if [ -f "adapter-demo" ]; then
        print_success "构建成功！"
        print_info "可执行文件: ${BACKEND_DIR}/adapter-demo"
        echo ""
        print_info "运行方式:"
        echo "  cd ${BACKEND_DIR}"
        echo "  ./adapter-demo -demo k12"
        echo "  ./adapter-demo -demo registry"
        echo "  ./adapter-demo -demo all"
    else
        print_error "构建失败！"
        exit 1
    fi
}

# 清理构建产物
clean() {
    print_header "清理构建产物"
    
    cd "$BACKEND_DIR"
    
    print_info "删除构建文件..."
    
    rm -f adapter-demo
    rm -f *.prof
    rm -f *.out
    rm -f debug.log
    
    # 也清理 examples 目录
    cd "$BACKEND_DIR/domain/adapter/examples"
    rm -f *.prof
    rm -f *.out
    
    print_success "清理完成！"
}

# 查看代码统计
stats() {
    print_header "代码统计"
    
    cd "$BACKEND_DIR/domain/adapter/examples"
    
    print_info "文件列表:"
    ls -lh *.go
    
    echo ""
    print_info "代码行数:"
    wc -l *.go
    
    echo ""
    print_info "函数统计:"
    grep -c "^func " *.go || true
}

# VSCode 调试
vscode_debug() {
    print_header "VSCode 调试指南"
    
    echo ""
    print_info "步骤 1: 在 VSCode 中打开项目"
    echo "  code ${PROJECT_ROOT}"
    
    echo ""
    print_info "步骤 2: 打开文件"
    echo "  ${DEMO_DIR}/demo.go"
    
    echo ""
    print_info "步骤 3: 设置断点"
    echo "  在需要调试的行号左侧点击，设置红点断点"
    echo "  推荐位置: 第 15, 44, 56, 71 行"
    
    echo ""
    print_info "步骤 4: 启动调试"
    echo "  方式 1: 按 F5 键"
    echo "  方式 2: 点击左侧调试图标，选择配置"
    
    echo ""
    print_info "可用的调试配置:"
    echo "  • Debug Adapter Demo (K12)"
    echo "  • Debug Adapter Demo (Registry)"
    echo "  • Debug Adapter Demo (All)"
    
    echo ""
    print_info "调试快捷键:"
    echo "  F5           - 继续执行"
    echo "  F10          - 单步跳过"
    echo "  F11          - 单步进入"
    echo "  Shift+F11    - 跳出"
    echo "  Shift+F5     - 停止调试"
    
    echo ""
    print_info "详细文档: docs/DEBUG_ADAPTER_DEMO.md"
}

# 性能分析
profile() {
    print_header "性能分析"
    
    cd "$BACKEND_DIR/domain/adapter/examples"
    
    print_info "生成 CPU profile..."
    go test -cpuprofile=cpu.prof -bench=. > /dev/null 2>&1
    
    if [ -f "cpu.prof" ]; then
        print_success "CPU profile 生成完成: cpu.prof"
        
        echo ""
        print_info "查看性能分析:"
        echo "  go tool pprof cpu.prof"
        echo "  go tool pprof -http=:8080 cpu.prof"
        
        echo ""
        read -p "是否打开浏览器查看？(y/N): " confirm
        
        if [ "$confirm" == "y" ] || [ "$confirm" == "Y" ]; then
            print_info "启动 pprof web 界面..."
            go tool pprof -http=:8080 cpu.prof
        fi
    else
        print_error "性能分析失败！"
        exit 1
    fi
}

# 显示帮助
show_help() {
    cat << EOF
Adapter Demo 调试助手

用法: $0 [command]

命令:
  k12          运行 K12 适配器演示
  registry     运行适配器注册表演示
  all          运行所有演示
  test         运行单元测试
  bench        运行性能测试
  debug        使用 Delve 调试器 [k12|registry|all]
  build        构建可执行文件
  clean        清理构建产物
  stats        查看代码统计
  vscode       显示 VSCode 调试指南
  profile      生成性能分析报告
  help         显示此帮助信息

示例:
  # 运行 K12 演示
  $0 k12
  
  # 运行测试
  $0 test
  
  # 使用 Delve 调试
  $0 debug k12
  
  # VSCode 调试指南
  $0 vscode
  
  # 性能分析
  $0 profile

快速开始:
  1. 运行演示:     $0 k12
  2. 查看调试帮助: $0 vscode
  3. 阅读文档:     docs/DEBUG_ADAPTER_DEMO.md

EOF
}

# 主函数
main() {
    # 只在需要 Go 的命令时才检查
    case "${1:-help}" in
        help|--help|-h|vscode)
            # 这些命令不需要 Go
            ;;
        *)
            check_go
            ;;
    esac
    
    case "${1:-help}" in
        k12)
            run_k12
            ;;
        registry)
            run_registry
            ;;
        all)
            run_all
            ;;
        test)
            run_tests
            ;;
        bench)
            run_bench
            ;;
        debug)
            debug_delve "${2:-k12}"
            ;;
        build)
            build_demo
            ;;
        clean)
            clean
            ;;
        stats)
            stats
            ;;
        vscode)
            vscode_debug
            ;;
        profile)
            profile
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
