# Service Discovery з Consul у Мікросервісній Платформі для Опитувань

Цей документ описує інтеграцію Service Discovery за допомогою Consul у нашій мікросервісній платформі для проведення соціологічних опитувань.

## Загальний огляд

У нашій мікросервісній архітектурі ми використовуємо Consul для виявлення сервісів (Service Discovery). Це дозволяє мікросервісам динамічно знаходити один одного без жорстко прописаних адрес, що спрощує масштабування та підвищує стійкість системи.

## Реалізовані компоненти:

1. **Health Check ендпоінти** у кожному мікросервісі:
   - `/health` ендпоінт перевіряє підключення до баз даних (PostgreSQL або MongoDB) та інших критичних залежностей
   - Повертає стандартизований JSON-відповідь зі статусом сервісу та його залежностей

2. **Реєстрація сервісів у Consul**:
   - Кожен мікросервіс реєструється в Consul при запуску
   - Генерується унікальний ID для кожного екземпляра сервісу
   - Сервіс дереєструється з Consul при коректному завершенні роботи

3. **API Gateway взаємодія з Consul**:
   - API Gateway використовує Consul для динамічного виявлення сервісів
   - Circuit Breaker на рівні API Gateway обробляє ситуації, коли сервіси недоступні

4. **Tracing через Correlation ID**:
   - HTTP заголовок `X-Correlation-ID` використовується для відстеження запитів через мікросервіси
   - Якщо заголовок не передано, створюється новий унікальний ID
   - Кожен мікросервіс передає цей ID далі при викликах інших сервісів
   - Полегшує відстеження та діагностику проблем у розподіленому середовищі

## Технічна реалізація

### Консул

Ми використовуємо офіційний Docker образ Consul:

```yaml
consul:
  image: hashicorp/consul:1.15
  container_name: consul
  ports:
    - "8500:8500"
    - "8600:8600/udp"
  environment:
    - CONSUL_BIND_INTERFACE=eth0
  command: "agent -dev -client=0.0.0.0"
```

### Реєстрація сервісу в Consul

Кожен мікросервіс використовує спільний код із пакета `backend/pkg/consul/client.go` для реєстрації в Consul:

```go
// Реєстрація сервісу в Consul
serviceID := fmt.Sprintf("%s-%s", cfg.ServiceName, uuid.New().String())
err = consulClient.RegisterService(
    serviceID,
    cfg.ServiceName,
    cfg.Server.Host,
    cfg.Server.Port,
    []string{},
    fmt.Sprintf("http://%s:%s/health", cfg.Server.Host, cfg.Server.Port),
)
```

### Health Check ендпоінт

Кожен мікросервіс має ендпоінт для перевірки стану:

```go
// HealthCheckHandler обробляє запити до /health ендпоінту
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    status := HealthStatus{
        Status:      "healthy",
        ServiceName: h.serviceName,
    }

    // Перевірка підключення до бази даних
    if err := h.checkDatabase(r.Context()); err != nil {
        status.Status = "unhealthy"
        status.Database = "disconnected"
        status.Details = err.Error()
        w.WriteHeader(http.StatusServiceUnavailable)
    } else {
        status.Database = "connected"
        w.WriteHeader(http.StatusOK)
    }

    // Відповідь у форматі JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}
```

### Correlation ID for Distributed Tracing

У пакеті `backend/pkg/tracing/correlation.go` реалізовано middleware для роботи з Correlation ID:

```go
// Middleware додає обробку Correlation ID до HTTP запитів
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Отримати Correlation ID з заголовка або згенерувати новий
        correlationID := r.Header.Get(CorrelationIDHeader)
        if correlationID == "" {
            correlationID = uuid.New().String()
        }

        // Додати Correlation ID до заголовка відповіді
        w.Header().Set(CorrelationIDHeader, correlationID)

        // Встановити Correlation ID в контексті запиту
        ctx := context.WithValue(r.Context(), contextKey, correlationID)
        r = r.WithContext(ctx)

        // Викликати наступний обробник
        next.ServeHTTP(w, r)
    })
}
```

## Тестування Service Discovery

Для перевірки коректної роботи Service Discovery ми розробили інтеграційні тести, які перевіряють:

1. Наскрізні сценарії через API Gateway
2. Реакцію системи на недоступність окремих сервісів
3. Відновлення після повернення сервісів у робочий стан

## Подальші кроки

1. **Розширене тестування відмовостійкості**:
   - Додаткові тести для Circuit Breaker
   - Симуляція різних сценаріїв відмов сервісів та мережі

2. **Розширення Distributed Tracing**:
   - Інтеграція з OpenTelemetry/Jaeger для повноцінного трасування
   - Збір та аналіз трас запитів через усі сервіси

3. **Покращення моніторингу**:
   - Додавання Prometheus для збору метрик
   - Налаштування графіків та алертів у Grafana

4. **Автоматизація масштабування**:
   - Інтеграція з Kubernetes для автоматичного масштабування
   - Використання Consul для балансування навантаження між екземплярами сервісів 