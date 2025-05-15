package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the survey_taking_service
type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
	MongoDB  MongoDBConfig
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port string
}

// RabbitMQConfig holds the configuration for RabbitMQ
type RabbitMQConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Exchange   string
	Queue      string
	RoutingKey string
}

// MongoDBConfig holds the configuration for MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8084"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:       getEnv("RABBITMQ_HOST", "localhost"),
			Port:       getEnv("RABBITMQ_PORT", "5672"),
			Username:   getEnv("RABBITMQ_USERNAME", "guest"),
			Password:   getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:   getEnv("RABBITMQ_EXCHANGE", "survey_responses"),
			Queue:      getEnv("RABBITMQ_QUEUE", "survey_responses_queue"),
			RoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "survey.response"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "survey_platform"),
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
