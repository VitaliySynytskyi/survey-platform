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
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api/handlers"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/store/mongodb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Підключення до MongoDB
	mongoClient, err := connectToMongoDB(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer disconnectFromMongoDB(mongoClient)

	// Ініціалізація сховища даних
	db := mongoClient.Database(cfg.MongoDB.Database)
	repository := mongodb.NewSurveyRepository(db)

	// Ініціалізація обробників API
	surveyHandler := handlers.NewSurveyHandler(repository, mongoClient, cfg.MongoDB.Database)

	// Ініціалізація маршрутизатора
	router := api.NewRouter(cfg, surveyHandler, mongoClient)

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
			serviceAddress := cfg.ServiceName
			servicePort := cfg.Server.Port
			healthCheckURL := fmt.Sprintf("http://%s:%d/health", serviceAddress, servicePort)

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

	// Налаштування HTTP-сервера
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в окремій горутині
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Налаштування граціозного завершення роботи
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
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

	// Встановлення таймауту для завершення
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Спроба граціозного завершення роботи сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// connectToMongoDB підключення до MongoDB
func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Перевірка підключення
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}

// disconnectFromMongoDB відключення від MongoDB
func disconnectFromMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}
