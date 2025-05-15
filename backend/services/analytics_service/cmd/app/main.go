package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repository
	repo, err := db.NewAnalyticsRepository(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	analyticsService := service.NewAnalyticsService(repo)
	defer func() {
		if err := analyticsService.Close(context.Background()); err != nil {
			log.Printf("Error closing service: %v", err)
		}
	}()

	// Create handler
	handler := api.NewHandler(analyticsService)

	// Setup router
	router := api.SetupRouter(handler, cfg.Auth)

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting analytics service on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
