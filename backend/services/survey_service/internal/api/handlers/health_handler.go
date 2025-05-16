package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// HealthStatus представляє стан здоров'я сервісу та його залежностей
type HealthStatus struct {
	Status      string `json:"status"`
	MongoDB     string `json:"mongodb,omitempty"`
	Details     string `json:"details,omitempty"`
	ServiceName string `json:"service_name"`
}

// HealthHandler обробляє запити до ендпоінту здоров'я
type HealthHandler struct {
	mongoClient *mongo.Client
	serviceName string
}

// NewHealthHandler створює новий обробник здоров'я
func NewHealthHandler(mongoClient *mongo.Client, serviceName string) *HealthHandler {
	return &HealthHandler{
		mongoClient: mongoClient,
		serviceName: serviceName,
	}
}

// HealthCheck обробляє GET запити до /health ендпоінту
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:      "healthy",
		ServiceName: h.serviceName,
	}

	// Перевірка підключення до MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		status.Status = "unhealthy"
		status.MongoDB = "disconnected"
		status.Details = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.MongoDB = "connected"
		w.WriteHeader(http.StatusOK)
	}

	// Повертаємо відповідь у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
