package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/db"
)

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status       string            `json:"status"`
	ServiceName  string            `json:"service_name"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Details      string            `json:"details,omitempty"`
}

// SetupRouter sets up the router
func SetupRouter(handler *Handler, cfg config.Auth) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:80", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth middleware
	jwtMiddleware := auth.JWTMiddleware(cfg)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", healthCheckHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)

		// Survey results
		r.Get("/surveys/{surveyId}/results", handler.GetSurveyResults)

		// Individual responses
		r.Get("/surveys/{surveyId}/responses", handler.GetIndividualResponses)
	})

	return r
}

// healthCheckHandler handles the health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()
	healthy := true
	var details string

	status := HealthStatus{
		Status:       "healthy",
		ServiceName:  "analytics_service",
		Dependencies: make(map[string]string),
	}

	// Check MongoDB connection
	repo, err := db.NewAnalyticsRepository(cfg.MongoDB)
	if err != nil {
		healthy = false
		details += "MongoDB: " + err.Error() + "; "
		status.Dependencies["mongodb"] = "disconnected"
	} else {
		status.Dependencies["mongodb"] = "connected"
		defer repo.Close(r.Context())
	}

	// Set status based on health check results
	if !healthy {
		status.Status = "unhealthy"
		status.Details = details
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
