#!/bin/bash
# =============================================================================
# MetaWorkflow V2.0 - 数据库迁移主脚本
# =============================================================================
# 创建时间: 2025-12-10
# 说明: 按顺序执行所有数据库迁移脚本
# 使用方法: ./migrate.sh [up|down|reset]
# =============================================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置（从环境变量或配置文件读取）
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-metaworkflow}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

# 迁移文件目录
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 迁移文件列表（按执行顺序）
MIGRATION_FILES=(
    "00_init_functions.sql"
    "01_users_and_permissions.sql"
    "02_core_business.sql"
    "03_adapter_marketplace.sql"
    "04_audit_and_logs.sql"
)

# 打印信息函数
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

# 执行SQL文件
execute_sql_file() {
    local sql_file=$1
    local file_path="${MIGRATIONS_DIR}/${sql_file}"
    
    if [ ! -f "$file_path" ]; then
        print_error "Migration file not found: $sql_file"
        return 1
    fi
    
    print_info "Executing: $sql_file"
    
    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -f "$file_path" \
        -v ON_ERROR_STOP=1 \
        --quiet
    
    if [ $? -eq 0 ]; then
        print_success "✓ $sql_file executed successfully"
        return 0
    else
        print_error "✗ Failed to execute $sql_file"
        return 1
    fi
}

# 检查数据库连接
check_database_connection() {
    print_info "Checking database connection..."
    
    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "postgres" \
        -c "SELECT 1;" \
        > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        print_success "Database connection successful"
        return 0
    else
        print_error "Cannot connect to database"
        print_error "Host: $DB_HOST:$DB_PORT, User: $DB_USER, Database: postgres"
        return 1
    fi
}

# 创建数据库（如果不存在）
create_database() {
    print_info "Creating database if not exists: $DB_NAME"
    
    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "postgres" \
        -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" \
        | grep -q 1
    
    if [ $? -eq 0 ]; then
        print_warning "Database $DB_NAME already exists"
    else
        PGPASSWORD=$DB_PASSWORD psql \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "postgres" \
            -c "CREATE DATABASE $DB_NAME WITH ENCODING='UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8';" \
            > /dev/null 2>&1
        
        if [ $? -eq 0 ]; then
            print_success "Database $DB_NAME created successfully"
        else
            print_error "Failed to create database $DB_NAME"
            return 1
        fi
    fi
}

# 删除数据库
drop_database() {
    print_warning "Dropping database: $DB_NAME"
    read -p "Are you sure you want to drop the database? (yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        print_info "Operation cancelled"
        return 0
    fi
    
    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "postgres" \
        -c "DROP DATABASE IF EXISTS $DB_NAME;" \
        > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        print_success "Database $DB_NAME dropped successfully"
        return 0
    else
        print_error "Failed to drop database $DB_NAME"
        return 1
    fi
}

# 向上迁移（创建所有表）
migrate_up() {
    print_info "=========================================="
    print_info "Starting database migration (UP)"
    print_info "=========================================="
    
    # 检查连接
    check_database_connection || exit 1
    
    # 创建数据库
    create_database || exit 1
    
    # 执行迁移文件
    for migration_file in "${MIGRATION_FILES[@]}"; do
        execute_sql_file "$migration_file" || exit 1
    done
    
    print_success "=========================================="
    print_success "All migrations completed successfully!"
    print_success "=========================================="
}

# 向下迁移（删除所有表）
migrate_down() {
    print_warning "=========================================="
    print_warning "Starting database rollback (DOWN)"
    print_warning "=========================================="
    
    read -p "Are you sure you want to drop all tables? (yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        print_info "Operation cancelled"
        return 0
    fi
    
    # 这里可以添加具体的 DROP TABLE 语句
    # 或者直接删除整个数据库
    drop_database
    
    print_success "=========================================="
    print_success "Rollback completed"
    print_success "=========================================="
}

# 重置数据库（删除并重建）
migrate_reset() {
    print_warning "=========================================="
    print_warning "Resetting database"
    print_warning "=========================================="
    
    migrate_down
    migrate_up
}

# 显示帮助信息
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up       Run all migrations (default)"
    echo "  down     Rollback all migrations"
    echo "  reset    Drop and recreate database"
    echo "  help     Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST      Database host (default: localhost)"
    echo "  DB_PORT      Database port (default: 5432)"
    echo "  DB_NAME      Database name (default: metaworkflow)"
    echo "  DB_USER      Database user (default: postgres)"
    echo "  DB_PASSWORD  Database password (default: postgres)"
    echo ""
    echo "Example:"
    echo "  DB_NAME=my_db DB_USER=myuser DB_PASSWORD=mypass $0 up"
}

# 主函数
main() {
    local command="${1:-up}"
    
    case "$command" in
        up)
            migrate_up
            ;;
        down)
            migrate_down
            ;;
        reset)
            migrate_reset
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
