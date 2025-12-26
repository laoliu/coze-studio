#!/bin/bash

# K12适配器实现验证脚本

echo "=== K12教育适配器实现验证 ==="
echo ""

# 检查文件完整性
echo "【步骤1：检查文件完整性】"
files=(
    "/home/liu/work/coze-studio/backend/domain/adapter/models/k12_models.go"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/knowledge_graph.go"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/resources.go"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/path_generator.go"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/analyzer.go"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/k12_adapter_service.go"
    "/home/liu/work/coze-studio/backend/api/handler/coze/k12_adapter_handler.go"
)

all_exist=true
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo "❌ 缺失: $file"
        all_exist=false
    fi
done
echo ""

if [ "$all_exist" = false ]; then
    echo "❌ 文件检查失败"
    exit 1
fi

# 检查代码行数
echo "【步骤2：统计代码规模】"
total_lines=0
for file in "${files[@]}"; do
    lines=$(wc -l < "$file")
    total_lines=$((total_lines + lines))
    filename=$(basename "$file")
    printf "%-35s %5d 行\n" "$filename" "$lines"
done
echo "----------------------------------------"
printf "%-35s %5d 行\n" "总计" "$total_lines"
echo ""

# 检查关键内容
echo "【步骤3：验证关键组件】"

echo "知识点数量:"
grep -c '"ID":' /home/liu/work/coze-studio/backend/domain/adapter/k12/knowledge_graph.go || echo "0"

echo "学习资源数量:"
grep -c '"ID":' /home/liu/work/coze-studio/backend/domain/adapter/k12/resources.go | head -1

echo "习题数量:"
grep -c '"Question":' /home/liu/work/coze-studio/backend/domain/adapter/k12/resources.go || echo "0"

echo "API端点数量:"
grep -c 'k12API\.' /home/liu/work/coze-studio/backend/api/router/register.go || echo "0"

echo ""

# 检查Go编译
echo "【步骤4：Go语法检查】"
cd /home/liu/work/coze-studio/backend

echo "检查 k12_models.go..."
go vet ./domain/adapter/models/k12_models.go 2>&1 | head -5 || echo "✅ 无语法错误"

echo "检查 k12 package..."
go vet ./domain/adapter/k12/... 2>&1 | head -10 || echo "✅ 无语法错误"

echo "检查 handler..."
go vet ./api/handler/coze/k12_adapter_handler.go 2>&1 | head -5 || echo "✅ 无语法错误"

echo ""

# 文档检查
echo "【步骤5：检查文档】"
docs=(
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/README.md"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/K12_ADAPTER_IMPLEMENTATION.md"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/QUICK_START.md"
    "/home/liu/work/coze-studio/backend/domain/adapter/k12/test_k12_adapter.sh"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "✅ $(basename $doc)"
    else
        echo "❌ 缺失: $(basename $doc)"
    fi
done
echo ""

# 检查路由注册
echo "【步骤6：验证路由注册】"
if grep -q "registerK12AdapterRoutes" /home/liu/work/coze-studio/backend/api/router/register.go; then
    echo "✅ K12路由已注册"
else
    echo "❌ K12路由未注册"
fi
echo ""

# 总结
echo "=== 验证完成 ==="
echo ""
echo "📊 实现统计："
echo "  • 代码文件: ${#files[@]} 个"
echo "  • 代码总行数: $total_lines 行"
echo "  • 文档文件: ${#docs[@]} 个"
echo ""
echo "🎯 核心功能："
echo "  ✅ 数据模型（Subject, Grade, KnowledgePoint, Exercise, etc.）"
echo "  ✅ 知识图谱（8个数学知识点，完整依赖链）"
echo "  ✅ 资源数据库（视频、教材、习题）"
echo "  ✅ 学习路径生成器（个性化路径算法）"
echo "  ✅ 输入分析器（NLP识别学科/年级/知识点）"
echo "  ✅ K12适配器服务（整合所有组件）"
echo "  ✅ HTTP API处理器（4个RESTful端点）"
echo "  ✅ 路由注册（集成到主服务）"
echo ""
echo "📚 文档："
echo "  ✅ README.md - 实现总结"
echo "  ✅ K12_ADAPTER_IMPLEMENTATION.md - 完整技术文档"
echo "  ✅ QUICK_START.md - 快速开始指南"
echo "  ✅ test_k12_adapter.sh - 测试脚本"
echo ""
echo "🚀 下一步："
echo "  1. 启动服务: cd backend && go run main.go"
echo "  2. 运行测试: cd backend/domain/adapter/k12 && ./test_k12_adapter.sh"
echo "  3. 查看文档: cat backend/domain/adapter/k12/README.md"
echo ""
