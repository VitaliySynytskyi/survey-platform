package main

import (
	"log"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/config"
)

func main() {
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

	// Wait for interrupt signal
	server.WaitForSignal()

	// Stop the server gracefully
	server.Stop()
}
