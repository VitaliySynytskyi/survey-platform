package api

import (
	"net/http"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api/handlers"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api/handlers/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewRouter створює та налаштовує маршрутизатор API
func NewRouter(cfg *config.Config, surveyHandler *handlers.SurveyHandler, mongoClient *mongo.Client) http.Handler {
	router := mux.NewRouter()

	// Middleware для всіх запитів
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Middleware для автентифікації
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Health check ендпоінт - публічний доступ
	healthHandler := handlers.NewHealthHandler(mongoClient, cfg.ServiceName)
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// API ендпоінти
	api := router.PathPrefix("/api/v1").Subrouter()

	// Захищені ендпоінти (потрібна автентифікація)
	api.Handle("/surveys", authMiddleware(http.HandlerFunc(surveyHandler.Create))).Methods("POST")
	api.Handle("/surveys/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.GetByID))).Methods("GET")
	api.Handle("/surveys/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.Update))).Methods("PUT")
	api.Handle("/surveys/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.Delete))).Methods("DELETE")
	api.Handle("/users/{userId}/surveys", authMiddleware(http.HandlerFunc(surveyHandler.GetUserSurveys))).Methods("GET")

	// Обробник для 404 помилок
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not found"}`))
	})

	return router
}
