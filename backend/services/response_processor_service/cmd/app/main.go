package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/rabbitmq"
)

func main() {
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
