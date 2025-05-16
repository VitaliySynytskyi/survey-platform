# Інструкція з розгортання мікросервісної платформи для опитувань

## 1. Підготовка до розгортання

### Вимоги до системи
- Docker та Docker Compose версії 1.28.0 або вище
- Мінімум 4GB RAM
- 20GB вільного дискового простору
- Підключення до інтернету для завантаження образів Docker

### Структура файлів для розгортання
- `docker-compose.yml` - основний файл для розгортання з використанням Docker Compose
- `.env` - файл з змінними середовища (потрібно створити перед запуском)

## 2. Створення файлу .env

Перед запуском системи необхідно створити файл `.env` в кореневій директорії проекту з наступним вмістом:

```
# Database Configuration
POSTGRES_USER=survey_user
POSTGRES_PASSWORD=survey_password
POSTGRES_DB=survey_platform
POSTGRES_PORT=5432

# MongoDB Configuration
MONGO_USER=mongo_admin
MONGO_PASSWORD=mongo_password
MONGO_PORT=27017

# RabbitMQ Configuration
RABBITMQ_USER=rabbit_user
RABBITMQ_PASSWORD=rabbit_password
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672

# Service Ports
AUTH_SERVICE_PORT=8080
USER_SERVICE_PORT=8081
SURVEY_SERVICE_PORT=8082
SURVEY_TAKING_SERVICE_PORT=8083
ANALYTICS_SERVICE_PORT=8084
API_GATEWAY_PORT=8000
FRONTEND_PORT=80

# JWT Secret
JWT_SECRET=your_jwt_secret_key_change_this_in_production

# Consul Configuration
CONSUL_PORT=8500

# Environment (development, staging, production)
ENVIRONMENT=development
```

**ВАЖЛИВО**: У виробничому середовищі обов'язково змініть всі паролі та JWT_SECRET на унікальні та надійні значення.

## 3. Запуск системи з Docker Compose

### Перший запуск
```bash
# Завантажити та зібрати всі образи
docker-compose build

# Запустити всі сервіси
docker-compose up -d

# Перевірити статус всіх контейнерів
docker-compose ps
```

### Перевірка роботи системи
1. Відкрийте браузер і перейдіть на `http://localhost` для доступу до фронтенду.
2. API Gateway доступний за адресою `http://localhost:8000`.
3. Консоль управління Consul доступна за адресою `http://localhost:8500`.
4. Інтерфейс RabbitMQ доступний за адресою `http://localhost:15672` (логін та пароль з .env файлу).

### Зупинка системи
```bash
docker-compose down
```

### Зупинка системи з видаленням даних (для повного перезапуску)
```bash
docker-compose down -v
```

## 4. Оновлення системи

Для оновлення системи до нової версії:

```bash
# Отримати останні зміни з репозиторію
git pull

# Перебудувати та перезапустити контейнери
docker-compose build
docker-compose up -d
```

## 5. Масштабування (опціонально)

Docker Compose дозволяє масштабувати окремі сервіси:

```bash
# Наприклад, для запуску 3 екземплярів response_processor_service
docker-compose up -d --scale response_processor_service=3
```

**Примітка**: Не всі сервіси можна масштабувати таким чином через обмеження на порти. Для повноцінного масштабування рекомендується використовувати Kubernetes.

## 6. Моніторинг логів

```bash
# Для перегляду логів усіх сервісів
docker-compose logs -f

# Для перегляду логів конкретного сервісу
docker-compose logs -f service_name
```

## 7. Резервне копіювання даних

Для створення резервних копій даних PostgreSQL:

```bash
docker-compose exec postgres pg_dump -U survey_user survey_platform > backup_$(date +%Y%m%d).sql
```

Для резервного копіювання MongoDB:

```bash
docker-compose exec mongodb mongodump --username mongo_admin --password mongo_password --out /dump
docker cp $(docker-compose ps -q mongodb):/dump ./mongo_backup_$(date +%Y%m%d)
```

## 8. Вирішення проблем

### Перевірка статусу сервісів:
```bash
docker-compose ps
```

### Перезапуск окремого сервісу:
```bash
docker-compose restart service_name
```

### Перевірка логів конкретного сервісу:
```bash
docker-compose logs -f service_name
```

### Перевірка мережевих з'єднань:
```bash
docker network inspect microservices-survey-app3_default
``` 