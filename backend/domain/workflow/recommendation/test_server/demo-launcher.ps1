# 快速启动推荐系统演示

Write-Host "🎯 节点推荐系统 - 演示启动器" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 检查当前目录
$currentDir = Get-Location
Write-Host "📁 当前目录: $currentDir" -ForegroundColor Yellow

# 检查是否在正确的目录
if (-not (Test-Path "main.go")) {
    Write-Host "❌ 错误: 请在 test_server 目录中运行此脚本" -ForegroundColor Red
    Write-Host "   正确路径: H:\work\python\coze-studio\backend\domain\workflow\recommendation\test_server" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ 目录检查通过" -ForegroundColor Green
Write-Host ""

# 检查服务器状态
Write-Host "🔍 检查服务器状态..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri 'http://localhost:8080/health' -Method Get -ErrorAction Stop
    Write-Host "✅ 服务器已运行: $($response.service) v$($response.version)" -ForegroundColor Green
    $serverRunning = $true
} catch {
    Write-Host "⚠️  服务器未运行，需要启动" -ForegroundColor Yellow
    $serverRunning = $false
}

Write-Host ""

# 如果服务器未运行，启动它
if (-not $serverRunning) {
    Write-Host "🚀 正在启动测试服务器..." -ForegroundColor Cyan
    Write-Host "   (这可能需要几秒钟)" -ForegroundColor Gray
    Write-Host ""
    
    # 在后台启动服务器
    $env:GOSUMDB = "sum.golang.org"
    
    Write-Host "⚙️  执行命令: go run -ldflags=`"-checklinkname=0`" main.go" -ForegroundColor Gray
    Write-Host ""
    Write-Host "⏳ 请在新终端窗口中运行:" -ForegroundColor Yellow
    Write-Host "   cd H:\work\python\coze-studio\backend\domain\workflow\recommendation\test_server" -ForegroundColor White
    Write-Host "   `$env:GOSUMDB=`"sum.golang.org`"; go run -ldflags=`"-checklinkname=0`" main.go" -ForegroundColor White
    Write-Host ""
    
    $continue = Read-Host "服务器启动后，按 Enter 继续..."
    
    # 再次检查
    try {
        $response = Invoke-RestMethod -Uri 'http://localhost:8080/health' -Method Get -ErrorAction Stop
        Write-Host "✅ 服务器已启动: $($response.service) v$($response.version)" -ForegroundColor Green
    } catch {
        Write-Host "❌ 服务器仍未运行，请先启动服务器" -ForegroundColor Red
        exit 1
    }
    
    Write-Host ""
}

# 显示演示选项
Write-Host "🎨 选择演示方式:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. 打开 HTML 演示页面 (推荐) ✨" -ForegroundColor White
Write-Host "   - 美观的交互式界面" -ForegroundColor Gray
Write-Host "   - 实时推荐显示" -ForegroundColor Gray
Write-Host "   - 无需前端启动" -ForegroundColor Gray
Write-Host ""
Write-Host "2. 打开浏览器测试页面" -ForegroundColor White
Write-Host "   - 9个测试场景" -ForegroundColor Gray
Write-Host "   - 手动触发测试" -ForegroundColor Gray
Write-Host ""
Write-Host "3. 查看使用指南" -ForegroundColor White
Write-Host "   - 完整的使用说明" -ForegroundColor Gray
Write-Host ""
Write-Host "4. 运行集成测试" -ForegroundColor White
Write-Host "   - 自动化测试脚本" -ForegroundColor Gray
Write-Host ""

$choice = Read-Host "请选择 (1-4)"

switch ($choice) {
    "1" {
        Write-Host ""
        Write-Host "🌐 正在打开演示页面..." -ForegroundColor Cyan
        Start-Process "http://localhost:8080/demo.html"
        Write-Host "✅ 已在浏览器中打开 demo.html" -ForegroundColor Green
        Write-Host ""
        Write-Host "💡 使用提示:" -ForegroundColor Yellow
        Write-Host "   1. 点击左侧的工作流节点" -ForegroundColor White
        Write-Host "   2. 右侧会显示推荐的下一个节点" -ForegroundColor White
        Write-Host "   3. 点击推荐项可以发送反馈" -ForegroundColor White
        Write-Host "   4. 使用底部按钮测试其他功能" -ForegroundColor White
    }
    "2" {
        Write-Host ""
        Write-Host "🌐 正在打开测试页面..." -ForegroundColor Cyan
        $htmlPath = Join-Path $currentDir "test-hook-manual.html"
        Start-Process $htmlPath
        Write-Host "✅ 已打开 test-hook-manual.html" -ForegroundColor Green
    }
    "3" {
        Write-Host ""
        Write-Host "📖 正在打开使用指南..." -ForegroundColor Cyan
        $guidePath = Join-Path (Split-Path $currentDir) "HOW_TO_VIEW_FRONTEND.md"
        if (Test-Path $guidePath) {
            code $guidePath
            Write-Host "✅ 已在 VS Code 中打开指南" -ForegroundColor Green
        } else {
            Write-Host "❌ 找不到指南文件: $guidePath" -ForegroundColor Red
        }
    }
    "4" {
        Write-Host ""
        Write-Host "🧪 正在运行集成测试..." -ForegroundColor Cyan
        .\test-integration.ps1
    }
    default {
        Write-Host ""
        Write-Host "❌ 无效的选择" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📚 更多信息请查看:" -ForegroundColor Cyan
Write-Host "   HOW_TO_VIEW_FRONTEND.md" -ForegroundColor White
Write-Host ""
Write-Host "🎉 享受推荐系统！" -ForegroundColor Green
Write-Host ""
