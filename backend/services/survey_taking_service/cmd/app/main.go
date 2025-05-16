package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	consul "github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/client"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/rabbitmq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Set up MongoDB client
	mongoClient, err := connectToMongoDB(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer disconnectFromMongoDB(mongoClient)

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

	// Add CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:80, http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register API routes
	handler.RegisterRoutes(router)

	// Add health check endpoint
	healthHandler := api.NewHealthHandler(cfg.ServiceName, mongoClient, producer.GetConnection())
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Initialize Consul client if enabled
	var consulClient *consul.Client
	var serviceID string
	if cfg.Consul.Enabled {
		consulClient, err = consul.NewClient(cfg.Consul.Address)
		if err != nil {
			log.Printf("Warning: Failed to create Consul client: %v", err)
		} else {
			// Generate a unique ID for this service instance
			serviceID = fmt.Sprintf("%s-%s", cfg.ServiceName, uuid.New().String())

			// Register service with Consul
			serviceAddress := cfg.Server.Host
			servicePort, _ := strconv.Atoi(cfg.Server.Port)
			healthCheckURL := fmt.Sprintf("http://%s:%s/health", serviceAddress, cfg.Server.Port)

			err = consulClient.RegisterService(
				serviceID,
				cfg.ServiceName,
				serviceAddress,
				servicePort,
				[]string{},
				healthCheckURL,
			)

			if err != nil {
				log.Printf("Warning: Failed to register service with Consul: %v", err)
			} else {
				log.Printf("Service registered with Consul: %s", serviceID)
			}
		}
	}

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

	// Deregister service from Consul if it was registered
	if cfg.Consul.Enabled && consulClient != nil && serviceID != "" {
		if err := consulClient.DeregisterService(serviceID); err != nil {
			log.Printf("Warning: Failed to deregister service from Consul: %v", err)
		} else {
			log.Printf("Service deregistered from Consul: %s", serviceID)
		}
	}

	// Create a deadline context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// connectToMongoDB establishes a connection to MongoDB
func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// disconnectFromMongoDB closes the connection to MongoDB
func disconnectFromMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}
