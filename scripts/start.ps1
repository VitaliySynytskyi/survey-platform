# PowerShell скрипт для запуску проекту
Write-Host "Building and starting services..." -ForegroundColor Green
docker-compose up -d

Write-Host "Services started!" -ForegroundColor Green 