package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/rabbitmq"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Set up MongoDB repository
	repository, err := db.NewResponseRepository(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to create response repository: %v", err)
	}

	// Set up RabbitMQ consumer
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ, repository)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// Setup HTTP server for health checks
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		dbErr := repository.CheckHealth(r.Context())

		// Check RabbitMQ connection
		rabbitErr := consumer.CheckHealth()

		if dbErr != nil || rabbitErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"checks": map[string]string{
					"database": dbErr.Error(),
					"rabbitmq": rabbitErr.Error(),
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "response_processor_service",
		})
	}).Methods("GET")

	// Start HTTP server in a goroutine
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Initialize Consul client
	var consulClient *consul.Client
	if consulAddr := os.Getenv("CONSUL_ADDR"); consulAddr != "" {
		consulClient, err = consul.NewClient(consulAddr)
		if err != nil {
			log.Printf("Warning: Failed to create Consul client: %v", err)
		} else {
			// Register service with Consul
			serviceID := fmt.Sprintf("response-processor-service-%s", uuid.New().String())
			serviceName := os.Getenv("SERVICE_NAME")
			if serviceName == "" {
				serviceName = "response_processor_service"
			}

			err = consulClient.RegisterService(
				serviceID,
				"response-processor-service",
				serviceName,
				8085, // Response processor service port
				[]string{"v1", "response-processor"},
				fmt.Sprintf("http://%s:8085/health", serviceName),
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

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages
	err = consumer.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Response processor service started")

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down service...")

	// Cancel context to stop consumer
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Give ongoing operations a chance to complete
	time.Sleep(2 * time.Second)

	// Close MongoDB connection
	closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeCancel()
	if err := repository.Close(closeCtx); err != nil {
		log.Printf("Error closing MongoDB connection: %v", err)
	}

	log.Println("Service exited gracefully")
}
