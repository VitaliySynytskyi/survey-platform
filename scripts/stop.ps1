# PowerShell скрипт для зупинки контейнерів
Write-Host "Stopping infrastructure containers..." -ForegroundColor Yellow
docker-compose -f docker-compose-infra.yml down

Write-Host "All containers stopped." -ForegroundColor Green 