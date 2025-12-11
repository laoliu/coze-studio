#!/bin/bash
# 一键设置 K12 Workflows 脚本
# 使用方法: ./setup_k12_workflows.sh

set -e  # 遇到错误立即退出

echo "════════════════════════════════════════════════════════"
echo "      K12 Workflow 创建脚本"
echo "════════════════════════════════════════════════════════"
echo ""

# 配置
MYSQL_CONTAINER="coze-mysql"
MYSQL_USER="root"
MYSQL_PASS="root"
DATABASE="opencoze"
API_ENDPOINT="http://localhost:8888/api/workflow_api/create"
SPACE_ID="1"
START_ID=1000000

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# 检查 Docker 容器是否运行
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "步骤 1: 检查 Coze 环境"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if ! docker ps | grep -q $MYSQL_CONTAINER; then
    print_error "MySQL 容器未运行"
    echo "请先启动 Coze 服务:"
    echo "  cd /home/liunix/work/coze-studio"
    echo "  make web"
    exit 1
fi
print_success "MySQL 容器运行中"

# 检查 API 是否可访问
if ! curl -s -f "$API_ENDPOINT" > /dev/null 2>&1; then
    if ! curl -s http://localhost:8888 > /dev/null 2>&1; then
        print_warning "Coze API 可能未启动 (端口 8888)"
        echo "请确保 coze-web 服务正在运行"
        echo ""
    fi
fi

# 步骤 2: 设置 AUTO_INCREMENT
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "步骤 2: 设置数据库 AUTO_INCREMENT = $START_ID"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查表是否存在
if ! docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE -e "SHOW TABLES LIKE 'workflow_meta';" 2>/dev/null | grep -q workflow_meta; then
    print_error "workflow_meta 表不存在"
    exit 1
fi

# 检查当前最大 ID
CURRENT_MAX=$(docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE -sN -e "SELECT COALESCE(MAX(id), 0) FROM workflow_meta;" 2>/dev/null)
echo "当前最大 Workflow ID: $CURRENT_MAX"

if [ "$CURRENT_MAX" -ge "$START_ID" ]; then
    print_warning "表中已有 ID >= $START_ID 的记录"
    echo "无法设置 AUTO_INCREMENT = $START_ID"
    echo ""
    echo "选项:"
    echo "  1. 使用当前最大 ID + 1 作为起始 (推荐)"
    echo "  2. 删除现有记录后重新设置"
    echo "  3. 退出"
    echo ""
    read -p "请选择 [1/2/3]: " choice
    
    case $choice in
        1)
            START_ID=$((CURRENT_MAX + 1))
            echo "将使用 ID $START_ID 开始"
            ;;
        2)
            print_warning "这将删除 ID >= $START_ID 的所有 Workflow 记录"
            read -p "确定继续? [y/N]: " confirm
            if [ "$confirm" != "y" ]; then
                echo "已取消"
                exit 0
            fi
            docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE -e "DELETE FROM workflow_meta WHERE id >= $START_ID;" 2>/dev/null
            print_success "已删除旧记录"
            ;;
        3|*)
            echo "已取消"
            exit 0
            ;;
    esac
fi

# 设置 AUTO_INCREMENT
docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE <<EOF 2>/dev/null
ALTER TABLE workflow_meta AUTO_INCREMENT = $START_ID;
EOF

print_success "AUTO_INCREMENT 已设置为 $START_ID"

# 步骤 3: 创建 Workflows
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "步骤 3: 创建 K12 Workflows"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 定义要创建的 Workflows
declare -a WORKFLOWS=(
    "K12 通用教学|适用于 K12 通用教学活动的工作流，支持多种教学场景"
    "K12 概念学习|专门用于概念理解的教学活动，包含导入、讲解、练习、总结等环节"
    "K12 练习|用于知识巩固和技能训练的练习活动"
    "K12 探索学习|引导学生自主探索和发现的探究式学习活动"
)

create_workflow() {
    local name="$1"
    local desc="$2"
    
    echo ""
    echo "创建: $name"
    
    response=$(curl -s -X POST "$API_ENDPOINT" \
      -H 'Content-Type: application/json' \
      -d "{
        \"name\": \"$name\",
        \"desc\": \"$desc\",
        \"icon_uri\": \"workflow/k12\",
        \"space_id\": \"$SPACE_ID\"
      }")
    
    # 检查响应
    if echo "$response" | grep -q '"id"'; then
        workflow_id=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        print_success "已创建 (ID: $workflow_id)"
        return 0
    else
        print_error "创建失败"
        echo "响应: $response"
        return 1
    fi
}

# 创建所有 Workflows
for workflow in "${WORKFLOWS[@]}"; do
    IFS='|' read -r name desc <<< "$workflow"
    create_workflow "$name" "$desc"
done

# 步骤 4: 验证创建结果
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "步骤 4: 验证创建结果"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE <<EOF 2>/dev/null
SELECT 
    id as 'Workflow ID',
    name as '名称',
    LEFT(description, 40) as '描述',
    created_at as '创建时间'
FROM workflow_meta 
WHERE id >= $START_ID
ORDER BY id;
EOF

# 步骤 5: 生成更新代码提示
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "步骤 5: 更新代码"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 获取实际创建的 IDs
ACTUAL_IDS=$(docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASS $DATABASE -sN -e "SELECT id FROM workflow_meta WHERE id >= $START_ID ORDER BY id;" 2>/dev/null)
IDS_ARRAY=($ACTUAL_IDS)

if [ ${#IDS_ARRAY[@]} -eq 4 ]; then
    echo ""
    echo "在 coze_k12_adapter.go 中更新 Workflow ID:"
    echo ""
    echo "func (a *CozeK12Adapter) ExecuteWorkflow(ctx *entity.AdapterContext) (*entity.WorkflowResult, error) {"
    echo "    var workflowID int64"
    echo "    switch ctx.Input.ActivityType {"
    echo "    case entity.ActivityTypeConcept:"
    echo "        workflowID = ${IDS_ARRAY[1]}  // K12 概念学习"
    echo "    case entity.ActivityTypePractice:"
    echo "        workflowID = ${IDS_ARRAY[2]}  // K12 练习"
    echo "    case entity.ActivityTypeExploration:"
    echo "        workflowID = ${IDS_ARRAY[3]}  // K12 探索学习"
    echo "    default:"
    echo "        workflowID = ${IDS_ARRAY[0]}  // K12 通用教学"
    echo "    }"
    echo "    ..."
    echo "}"
    echo ""
    
    # 提供直接可用的更新命令
    echo "快速更新命令:"
    echo "────────────────────────────────────────────────────"
    echo "sed -i 's/workflowID = 1000001/workflowID = ${IDS_ARRAY[1]}/g' \\"
    echo "  /home/liunix/work/coze-studio/backend/domain/adapter/examples/coze_k12_adapter.go"
    echo ""
    echo "sed -i 's/workflowID = 1000002/workflowID = ${IDS_ARRAY[2]}/g' \\"
    echo "  /home/liunix/work/coze-studio/backend/domain/adapter/examples/coze_k12_adapter.go"
    echo ""
    echo "sed -i 's/workflowID = 1000003/workflowID = ${IDS_ARRAY[3]}/g' \\"
    echo "  /home/liunix/work/coze-studio/backend/domain/adapter/examples/coze_k12_adapter.go"
    echo ""
    echo "sed -i 's/workflowID = 1000000/workflowID = ${IDS_ARRAY[0]}/g' \\"
    echo "  /home/liunix/work/coze-studio/backend/domain/adapter/examples/coze_k12_adapter.go"
    echo ""
fi

echo ""
print_success "设置完成！"
echo ""
echo "════════════════════════════════════════════════════════"
echo "  下一步: 运行 demo 测试 Workflow 执行"
echo "════════════════════════════════════════════════════════"
echo ""
echo "cd /home/liunix/work/coze-studio/backend/domain/adapter/examples/workflow_demo"
echo "go run demo_workflow.go coze_init.go"
echo ""
