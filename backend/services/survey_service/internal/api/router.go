package api

import (
	"log"
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

	// Middleware для логування всіх запитів
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("SURVEY_SERVICE_ROUTER: Incoming request: Method=%s, URL=%s, User-Agent=%s, RemoteAddr=%s", r.Method, r.URL.String(), r.UserAgent(), r.RemoteAddr)
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
	// Замість PathPrefix("/api/v1"), визначаємо маршрути відносно кореня,
	// оскільки API Gateway вже обробляє /api/v1/surveys префікс

	// Захищені ендпоінти (потрібна автентифікація)
	// POST /api/v1/surveys (через шлюз) -> POST / (тут)
	router.Handle("/", authMiddleware(http.HandlerFunc(surveyHandler.Create))).Methods("POST")
	// GET /api/v1/surveys (через шлюз) -> GET / (тут)
	router.Handle("/", authMiddleware(http.HandlerFunc(surveyHandler.GetAllSurveys))).Methods("GET")

	// GET /api/v1/surveys/{surveyId} (через шлюз) -> GET /{surveyId} (тут)
	router.Handle("/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.GetByID))).Methods("GET")
	// PUT /api/v1/surveys/{surveyId} (через шлюз) -> PUT /{surveyId} (тут)
	router.Handle("/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.Update))).Methods("PUT")
	// DELETE /api/v1/surveys/{surveyId} (через шлюз) -> DELETE /{surveyId} (тут)
	router.Handle("/{surveyId}", authMiddleware(http.HandlerFunc(surveyHandler.Delete))).Methods("DELETE")

	// GET /users/{userId}/surveys (прямий виклик від user_service, або через шлюз, якщо такий маршрут буде додано до шлюзу)
	// Цей маршрут залишається без змін, оскільки він специфічний і не конфліктує з загальними CRUD операціями
	router.Handle("/users/{userId}/surveys", authMiddleware(http.HandlerFunc(surveyHandler.GetUserSurveys))).Methods("GET")

	// Публічний ендпоінт для отримання опитування (без автентифікації)
	// GET /api/v1/surveys/{surveyId}/public (через шлюз) -> GET /{surveyId}/public (тут)
	router.HandleFunc("/{surveyId}/public", surveyHandler.GetPublicSurveyByID).Methods("GET")

	// Публічний ендпоінт для отримання результатів опитування (без автентифікації)
	// GET /api/v1/surveys/{surveyId}/results (через шлюз) -> GET /{surveyId}/results (тут)
	router.HandleFunc("/{surveyId}/results", surveyHandler.GetSurveyResults).Methods("GET")

	// Обробник для 404 помилок
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("SURVEY_SERVICE_ROUTER: Path not found (404): Method=%s, URL=%s", r.Method, r.URL.String())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not found in survey_service"}`))
	})

	return router
}
