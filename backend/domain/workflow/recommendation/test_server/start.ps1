# 启动推荐系统测试服务器（在新窗口中）
$scriptPath = $PSScriptRoot
$serverScript = @"
`$env:GOSUMDB = 'sum.golang.org'
Set-Location '$scriptPath'
Write-Host '🚀 Node Recommendation Test Server' -ForegroundColor Green
Write-Host '📁 Directory: `$(Get-Location)' -ForegroundColor Cyan
Write-Host ''
go run -ldflags='-checklinkname=0' main.go
Write-Host ''
Write-Host 'Press any key to close...' -ForegroundColor Yellow
`$null = `$Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
"@

# 在新的 PowerShell 窗口中运行
Start-Process pwsh -ArgumentList "-NoExit", "-Command", $serverScript

Write-Host "✅ Server starting in new window..." -ForegroundColor Green
Write-Host "🌐 URL: http://localhost:8080" -ForegroundColor Cyan
