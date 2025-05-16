package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// HealthStatus represents the health status of the service and its dependencies
type HealthStatus struct {
	Status       string            `json:"status"`
	ServiceName  string            `json:"service_name"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Details      string            `json:"details,omitempty"`
}

// HealthHandler is the handler for health check requests
type HealthHandler struct {
	serviceName  string
	mongoClient  *mongo.Client
	rabbitMQConn *amqp.Connection
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(serviceName string, mongoClient *mongo.Client, rabbitMQConn *amqp.Connection) *HealthHandler {
	return &HealthHandler{
		serviceName:  serviceName,
		mongoClient:  mongoClient,
		rabbitMQConn: rabbitMQConn,
	}
}

// HealthCheck handles GET /health requests
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:       "healthy",
		ServiceName:  h.serviceName,
		Dependencies: make(map[string]string),
	}

	healthy := true
	var details string

	// Check MongoDB connection
	if h.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := h.mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			healthy = false
			details += "MongoDB: " + err.Error() + "; "
			status.Dependencies["mongodb"] = "disconnected"
		} else {
			status.Dependencies["mongodb"] = "connected"
		}
	} else {
		status.Dependencies["mongodb"] = "not_configured"
	}

	// Check RabbitMQ connection
	if h.rabbitMQConn != nil {
		if h.rabbitMQConn.IsClosed() {
			healthy = false
			details += "RabbitMQ: connection closed; "
			status.Dependencies["rabbitmq"] = "disconnected"
		} else {
			status.Dependencies["rabbitmq"] = "connected"
		}
	} else {
		status.Dependencies["rabbitmq"] = "not_configured"
	}

	if !healthy {
		status.Status = "unhealthy"
		status.Details = details
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
