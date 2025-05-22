package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/handlers"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/repository"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/service"
)

func main() {
	log.Println("Starting response-service...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from compose or K8s")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize MongoDB repository
	mongoRepo, err := repository.NewMongoRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB repository: %v", err)
	}
	defer func() {
		if err := mongoRepo.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize service
	responseService := service.NewResponseService(mongoRepo, cfg.SurveyServiceURL)

	// Initialize handler
	responseHandler := handlers.NewResponseHandler(responseService)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "response-service",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/responses", responseHandler.SubmitResponse)
		// This new route is added to fetch responses for a specific survey
		// It needs to be distinct from POST /responses and likely grouped under surveys logically
		// e.g. /api/v1/surveys/:surveyId/responses
		// For now, as per SurveyResponses.vue, the call is to /api/v1/surveys/:id/responses,
		// so this route should be registered here if response-service handles it directly.
		api.GET("/surveys/:surveyId/responses", responseHandler.GetSurveyResponsesHandler)

		// Route for exporting survey responses as CSV
		api.GET("/surveys/:surveyId/responses/export", responseHandler.ExportSurveyResponsesCSV)

		// Route for survey analytics
		api.GET("/surveys/:surveyId/analytics", responseHandler.GetSurveyAnalytics)

		// Remove old placeholder routes for /responses if they were here
	}

	// Start server
	port := cfg.ServerPort // Use port from config
	if port == "" {
		port = "8083" // Default port if not in config or .env
	}

	log.Printf("Starting response-service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
