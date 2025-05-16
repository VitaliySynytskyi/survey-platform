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
	// Check MongoDB connection
	repo, err := db.NewAnalyticsRepository(config.Load().MongoDB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "MongoDB connection failed: " + err.Error(),
		})
		return
	}
	defer repo.Close(r.Context())

	// All checks passed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "analytics_service",
	})
}
