package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// SimpleMockHealthCheck простий мок для тестування HealthStatus
func TestHealthResponse(t *testing.T) {
	// Тестування відповіді зі статусом "healthy"
	t.Run("healthy response", func(t *testing.T) {
		// Створюємо запит
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Створюємо ResponseRecorder для запису відповіді
		rr := httptest.NewRecorder()

		// Функція-обробник, яка повертає "healthy" відповідь
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status := HealthStatus{
				Status:      "healthy",
				MongoDB:     "connected",
				ServiceName: "test-service",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(status)
		})

		// Викликаємо обробник
		handler.ServeHTTP(rr, req)

		// Перевіряємо статус-код
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("обробник повернув неправильний статус-код: отримано %v, очікувалося %v",
				status, http.StatusOK)
		}

		// Перевіряємо тіло відповіді
		var response HealthStatus
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("не вдалося розпарсити відповідь: %v", err)
		}

		if response.Status != "healthy" {
			t.Errorf("обробник повернув неправильний статус: отримано %v, очікувалося %v",
				response.Status, "healthy")
		}

		if response.MongoDB != "connected" {
			t.Errorf("обробник повернув неправильний статус MongoDB: отримано %v, очікувалося %v",
				response.MongoDB, "connected")
		}
	})

	// Тестування відповіді зі статусом "unhealthy"
	t.Run("unhealthy response", func(t *testing.T) {
		// Створюємо запит
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Створюємо ResponseRecorder для запису відповіді
		rr := httptest.NewRecorder()

		// Функція-обробник, яка повертає "unhealthy" відповідь
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status := HealthStatus{
				Status:      "unhealthy",
				MongoDB:     "disconnected",
				Details:     "failed to connect to MongoDB",
				ServiceName: "test-service",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(status)
		})

		// Викликаємо обробник
		handler.ServeHTTP(rr, req)

		// Перевіряємо статус-код
		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("обробник повернув неправильний статус-код: отримано %v, очікувалося %v",
				status, http.StatusServiceUnavailable)
		}

		// Перевіряємо тіло відповіді
		var response HealthStatus
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("не вдалося розпарсити відповідь: %v", err)
		}

		if response.Status != "unhealthy" {
			t.Errorf("обробник повернув неправильний статус: отримано %v, очікувалося %v",
				response.Status, "unhealthy")
		}

		if response.MongoDB != "disconnected" {
			t.Errorf("обробник повернув неправильний статус MongoDB: отримано %v, очікувалося %v",
				response.MongoDB, "disconnected")
		}
	})
}
