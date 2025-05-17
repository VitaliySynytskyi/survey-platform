package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	consul "github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/config"
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

	// Initialize Consul client if enabled
	var consulClient *consul.Client
	if cfg.Consul.Enabled {
		consulClient, err = consul.NewClient(cfg.Consul.Address)
		if err != nil {
			log.Printf("Warning: Failed to create Consul client: %v", err)
		} else {
			// Generate a unique ID for this service instance
			serviceID := fmt.Sprintf("%s-%s", cfg.ServiceName, uuid.New().String())

			// Register service with Consul
			// Use the service name as the address since we're in a Docker container
			// This will be the DNS name in the Docker network
			serviceAddress := cfg.ServiceName
			portNumber, _ := strconv.Atoi(cfg.Server.Port)
			healthCheckURL := fmt.Sprintf("http://%s:%s/health", serviceAddress, cfg.Server.Port)

			err = consulClient.RegisterService(
				serviceID,
				cfg.ServiceName,
				serviceAddress,
				portNumber,
				[]string{},
				healthCheckURL,
			)

			if err != nil {
				log.Printf("Warning: Failed to register service with Consul: %v", err)
			} else {
				log.Printf("Service registered with Consul: %s", serviceID)

				// Set up deregistration on shutdown
				defer func() {
					if err := consulClient.DeregisterService(serviceID); err != nil {
						log.Printf("Warning: Failed to deregister service from Consul: %v", err)
					} else {
						log.Printf("Service deregistered from Consul: %s", serviceID)
					}
				}()
			}
		}
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for signal to gracefully shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")

	// Stop the server gracefully
	server.Stop()
}
