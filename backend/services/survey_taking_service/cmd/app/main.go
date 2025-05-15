package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/client"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/rabbitmq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up survey client
	surveyClient := client.NewHTTPSurveyClient(os.Getenv("SURVEY_SERVICE_URL"))

	// Set up RabbitMQ producer
	producer, err := rabbitmq.NewProducer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	// Set up API handler
	handler := api.NewHandler(surveyClient, producer)

	// Set up router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Set up HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting survey_taking_service on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
