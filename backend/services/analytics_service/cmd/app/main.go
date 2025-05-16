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

	"github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/service"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

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

	// Initialize Consul client
	var consulClient *consul.Client
	if consulAddr := os.Getenv("CONSUL_ADDR"); consulAddr != "" {
		var err error
		consulClient, err = consul.NewClient(consulAddr)
		if err != nil {
			log.Printf("Warning: Failed to create Consul client: %v", err)
		} else {
			// Register service with Consul
			serviceID := fmt.Sprintf("analytics-service-%s", uuid.New().String())
			serviceName := os.Getenv("SERVICE_NAME")
			if serviceName == "" {
				serviceName = "analytics_service"
			}

			err = consulClient.RegisterService(
				serviceID,
				"analytics-service",
				serviceName,
				8086, // Analytics service port
				[]string{"v1", "analytics"},
				fmt.Sprintf("http://%s:8086/health", serviceName),
			)
			if err != nil {
				log.Printf("Warning: Failed to register service with Consul: %v", err)
			} else {
				log.Println("Successfully registered service with Consul")

				// Set up deregistration on shutdown
				defer func() {
					if err := consulClient.DeregisterService(serviceID); err != nil {
						log.Printf("Warning: Failed to deregister service: %v", err)
					} else {
						log.Println("Successfully deregistered service from Consul")
					}
				}()
			}
		}
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
