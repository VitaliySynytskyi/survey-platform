package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
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
		AllowedOrigins:   []string{"*"},
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
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
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
