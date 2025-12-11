# 推荐系统测试脚本
# 解决 sonic 库的链接器问题

Write-Host "运行推荐系统测试..." -ForegroundColor Green

# 设置环境变量
$env:GOSUMDB = "sum.golang.org"

# 运行测试（使用 -checklinkname=0 解决 sonic 库的链接问题）
go test -v -ldflags="-checklinkname=0"

Write-Host "`n测试完成！" -ForegroundColor Green
