# PowerShell скрипт для запуску проекту з Docker Compose
# Цей скрипт тимчасово копіює необхідні пакети в Docker контекст

# Створюємо тимчасовий каталог для зберігання модифікованих файлів
Write-Host "Створення тимчасової директорії" -ForegroundColor Yellow
$TempDir = "temp_build"
If(Test-Path $TempDir) {
    Remove-Item -Recurse -Force $TempDir
}
New-Item -ItemType Directory -Path $TempDir | Out-Null

# Функція для копіювання та модифікації сервісу
function Prepare-Service {
    param (
        [string]$ServiceName
    )
    
    Write-Host "Підготовка сервісу: $ServiceName" -ForegroundColor Yellow
    
    # Створюємо директорію для сервісу
    $ServiceDir = Join-Path $TempDir $ServiceName
    New-Item -ItemType Directory -Path $ServiceDir | Out-Null
    
    # Копіюємо вміст сервісу
    Copy-Item -Recurse "backend/services/$ServiceName/*" $ServiceDir
    
    # Копіюємо pkg директорію для локальних пакетів
    $PkgDir = Join-Path $ServiceDir "pkg"
    New-Item -ItemType Directory -Path $PkgDir | Out-Null
    Copy-Item -Recurse "backend/pkg" $PkgDir -ErrorAction SilentlyContinue
    
    # Модифікуємо go.mod файл, щоб видалити 'replace' директиви
    $GoModPath = Join-Path $ServiceDir "go.mod"
    if (Test-Path $GoModPath) {
        $GoModContent = Get-Content $GoModPath -Raw
        # Видаляємо replace директиви
        $GoModContent = $GoModContent -replace "replace.*=>.*\n", ""
        Set-Content -Path $GoModPath -Value $GoModContent
    }
}

# Підготовка всіх сервісів
$Services = @(
    "auth_service",
    "user_service", 
    "survey_service", 
    "survey_taking_service", 
    "response_processor_service", 
    "analytics_service"
)

foreach ($Service in $Services) {
    Prepare-Service -ServiceName $Service
}

# Підготовка API Gateway
Write-Host "Підготовка API Gateway" -ForegroundColor Yellow
$ApiGatewayDir = Join-Path $TempDir "api_gateway"
New-Item -ItemType Directory -Path $ApiGatewayDir | Out-Null
Copy-Item -Recurse "backend/api_gateway/*" $ApiGatewayDir
$PkgDir = Join-Path $ApiGatewayDir "pkg"
New-Item -ItemType Directory -Path $PkgDir | Out-Null
Copy-Item -Recurse "backend/pkg" $PkgDir -ErrorAction SilentlyContinue

# Створюємо тимчасовий docker-compose.yml
Write-Host "Створення тимчасового docker-compose.yml" -ForegroundColor Yellow
$DockerCompose = Get-Content "docker-compose.yml" -Raw
# Модифікуємо шляхи в docker-compose.yml
$Services | ForEach-Object {
    $DockerCompose = $DockerCompose -replace "context: \./backend/services/$_", "context: ./$TempDir/$_"
}
$DockerCompose = $DockerCompose -replace "context: \./backend/api_gateway", "context: ./$TempDir/api_gateway"
Set-Content -Path "$TempDir/docker-compose.yml" -Value $DockerCompose

# Запуск Docker Compose
Write-Host "Запуск Docker Compose" -ForegroundColor Green
cd $TempDir
docker-compose up -d

# Підсумок
Write-Host "Проект запущено! Перейдіть за адресою http://localhost для доступу до фронтенду." -ForegroundColor Green
Write-Host "Для зупинки проекту виконайте: docker-compose down в директорії $TempDir" -ForegroundColor Yellow 