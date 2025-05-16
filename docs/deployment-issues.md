# Проблеми розгортання та їх вирішення

## Виявлені проблеми

### 1. Несумісність версій Go
Dockerfile мікросервісів використовували Go 1.21, але в go.mod вказана вимога Go 1.22+. Це призводило до помилки:
```
go: go.mod requires go >= 1.22 (running go 1.21.13; GOTOOLCHAIN=local)
```

**Вирішення:** 
- Оновили Dockerfile для використання GO 1.22: `FROM golang:1.22-alpine AS builder`

### 2. Проблеми з локальними модулями в Docker контексті
Модулі використовують локальні пакети через директиви `replace` в go.mod файлах:
```
replace github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul => ../../pkg/consul
```

При збірці в Docker виникає помилка:
```
go: github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul@v0.0.0 (replaced by ../../pkg/consul): reading /pkg/consul/go.mod: open /pkg/consul/go.mod: no such file or directory
```

**Причина проблеми:**
При збірці Docker не має доступу до файлів за межами контексту збірки, а відносні шляхи в директивах `replace` вказують на директорії, недоступні в контексті Docker.

## Рекомендації для розгортання

### Підхід 1: Модифікація структури проекту
1. Створити тимчасову директорію для збірки, де кожен сервіс міститиме локальну копію залежних пакетів
2. Модифікувати go.mod файли, видаливши директиви `replace`
3. Використовувати оновлений Docker контекст для збірки

Приклад скрипту для підготовки:
```powershell
# PowerShell скрипт для підготовки тимчасової директорії
$TempDir = "temp_build"
If(Test-Path $TempDir) {
    Remove-Item -Recurse -Force $TempDir
}
New-Item -ItemType Directory -Path $TempDir | Out-Null

# Копіюємо вміст сервісу та його залежності
function Prepare-Service {
    param (
        [string]$ServiceName
    )
    
    $ServiceDir = Join-Path $TempDir $ServiceName
    New-Item -ItemType Directory -Path $ServiceDir | Out-Null
    
    # Копіюємо вміст сервісу
    Copy-Item -Recurse "backend/services/$ServiceName/*" $ServiceDir
    
    # Копіюємо pkg директорію
    $PkgDir = Join-Path $ServiceDir "pkg"
    New-Item -ItemType Directory -Path $PkgDir | Out-Null
    Copy-Item -Recurse "backend/pkg" $PkgDir -ErrorAction SilentlyContinue
    
    # Модифікуємо go.mod
    $GoModPath = Join-Path $ServiceDir "go.mod"
    if (Test-Path $GoModPath) {
        $GoModContent = Get-Content $GoModPath -Raw
        # Видаляємо replace директиви
        $GoModContent = $GoModContent -replace "replace.*=>.*\n", ""
        Set-Content -Path $GoModPath -Value $GoModContent
    }
}
```

### Підхід 2: Використання модуля "vendor"
1. Для кожного сервісу виконати `go mod vendor`, що створить директорію vendor з усіма залежностями
2. Включити цю директорію в Docker контекст
3. Оновити Dockerfile для використання vendor:
   ```dockerfile
   COPY . .
   RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main ./cmd/server/main.go
   ```

### Підхід 3: Поетапне розгортання
1. Запустити інфраструктурні компоненти через docker-compose:
   ```
   docker-compose -f docker-compose-infra.yml up -d
   ```

2. Розробити і запускати кожен мікросервіс локально для розробки/тестування

3. Опублікувати контейнери в реєстрі для продакшн

## Поточний стан розгортання

Інфраструктурні компоненти успішно запущені:
- PostgreSQL доступний на порту 5432
- MongoDB доступний на порту 27017
- RabbitMQ доступний на портах 5672 (AMQP) та 15672 (web UI)
- Consul доступний на порту 8500 (web UI)

## Рекомендовані наступні кроки

1. Впровадити один з описаних підходів для збірки мікросервісів
2. Створити окремі CI/CD pipeline для кожного мікросервісу
3. Налаштувати реєстр контейнерів для зберігання образів
4. Впровадити Kubernetes розгортання для продакшн середовища 