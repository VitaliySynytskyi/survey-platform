package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/config"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration from environment variables
	cfg := config.Load()

	// Create a new server
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up the database
	if err := server.SetupDatabase(); err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Initialize Consul client
	consulClient, err := consul.NewClient(os.Getenv("CONSUL_ADDR"))
	if err != nil {
		log.Printf("Warning: Failed to create Consul client: %v", err)
	} else {
		// Register service with Consul
		serviceID := fmt.Sprintf("auth-service-%s", uuid.New().String())
		err = consulClient.RegisterService(
			serviceID,
			"auth-service",
			os.Getenv("SERVICE_NAME"),
			8080,
			[]string{"v1", "auth"},
			fmt.Sprintf("http://%s:8080/health", os.Getenv("SERVICE_NAME")),
		)
		if err != nil {
			log.Printf("Warning: Failed to register service with Consul: %v", err)
		} else {
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

	// Wait for interrupt signal
	server.WaitForSignal()

	// Stop the server gracefully
	server.Stop()
}
