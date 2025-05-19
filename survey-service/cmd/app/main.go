package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/survey-app/survey-service/internal/config"
	"github.com/survey-app/survey-service/internal/handlers"
	"github.com/survey-app/survey-service/internal/repository"
	"github.com/survey-app/survey-service/internal/service"
)

func main() {
	log.Println("Starting survey-service application...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg := &config.Config{
		DB: config.DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "survey_db"),
		},
		Port: getEnv("PORT", "8082"),
	}

	// Database connection
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Ping database to verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize repository
	repo := repository.NewPostgresRepository(dbPool)

	// Initialize service
	surveyService := service.NewSurveyService(repo)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "survey-service",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Survey routes
		surveys := api.Group("/surveys")
		{
			surveyHandler := handlers.NewSurveyHandler(surveyService)
			surveys.POST("", surveyHandler.CreateSurvey)
			surveys.GET("/me", surveyHandler.GetMySurveys)
			surveys.GET("/all", surveyHandler.GetAllSurveysPublic)
			surveys.GET("/:id", surveyHandler.GetSurvey)
			surveys.PUT("/:id", surveyHandler.UpdateSurvey)
			surveys.DELETE("/:id", surveyHandler.DeleteSurvey)
			surveys.PATCH("/:id/status", surveyHandler.UpdateSurveyStatus)

			// Question routes
			surveys.POST("/:id/questions", surveyHandler.AddQuestion)
		}

		// Question routes
		questions := api.Group("/questions")
		{
			questionHandler := handlers.NewQuestionHandler(surveyService)
			questions.PUT("/:id", questionHandler.UpdateQuestion)
			questions.DELETE("/:id", questionHandler.DeleteQuestion)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Default port
	}

	log.Printf("Starting survey-service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
