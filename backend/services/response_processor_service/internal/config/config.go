package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the configuration for the response_processor_service
type Config struct {
	RabbitMQ RabbitMQConfig
	MongoDB  MongoDBConfig
	Server   ServerConfig
}

// RabbitMQConfig holds the configuration for RabbitMQ
type RabbitMQConfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Exchange      string
	Queue         string
	RoutingKey    string
	PrefetchCount int
}

// MongoDBConfig holds the configuration for MongoDB
type MongoDBConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	Collection string
	AuthSource string
	URI        string
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Port string
}

// Load loads the configuration from environment variables
func Load() *Config {
	mongoHost := getEnv("MONGO_HOST", "localhost")
	mongoPort := getEnv("MONGO_PORT", "27017")
	mongoUser := getEnv("MONGO_USER", "")
	mongoPassword := getEnv("MONGO_PASSWORD", "")
	mongoDatabase := getEnv("MONGO_DB", "survey_platform")
	mongoCollection := getEnv("MONGODB_COLLECTION", "responses")
	mongoAuthSource := getEnv("MONGO_AUTH_SOURCE", "admin")

	// Build MongoDB URI
	var mongoURI string
	if mongoUser != "" && mongoPassword != "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
			mongoUser, mongoPassword, mongoHost, mongoPort, mongoDatabase, mongoAuthSource)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%s/%s",
			mongoHost, mongoPort, mongoDatabase)
	}

	// Allow overriding the built URI with a complete URI
	if uri := getEnv("MONGODB_URI", ""); uri != "" {
		mongoURI = uri
	}

	return &Config{
		RabbitMQ: RabbitMQConfig{
			Host:          getEnv("RABBITMQ_HOST", "localhost"),
			Port:          getEnv("RABBITMQ_PORT", "5672"),
			Username:      getEnv("RABBITMQ_USER", "guest"),
			Password:      getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:      getEnv("RABBITMQ_EXCHANGE", "survey_responses"),
			Queue:         getEnv("RABBITMQ_QUEUE", "survey_responses_queue"),
			RoutingKey:    getEnv("RABBITMQ_ROUTING_KEY", "survey.response"),
			PrefetchCount: getEnvAsInt("RABBITMQ_PREFETCH_COUNT", 1),
		},
		MongoDB: MongoDBConfig{
			Host:       mongoHost,
			Port:       mongoPort,
			User:       mongoUser,
			Password:   mongoPassword,
			Database:   mongoDatabase,
			Collection: mongoCollection,
			AuthSource: mongoAuthSource,
			URI:        mongoURI,
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8085"),
		},
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get integer environment variables with fallback
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
