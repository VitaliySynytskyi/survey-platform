package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the survey_taking_service
type Config struct {
	Server      ServerConfig
	RabbitMQ    RabbitMQConfig
	MongoDB     MongoDBConfig
	Consul      ConsulConfig
	ServiceName string
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port string
	Host string
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

// ConsulConfig holds the configuration for Consul
type ConsulConfig struct {
	Address string
	Enabled bool
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8084"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
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
		Consul: ConsulConfig{
			Address: getEnv("CONSUL_ADDR", "consul:8500"),
			Enabled: getEnvAsBool("CONSUL_ENABLED", true),
		},
		ServiceName: getEnv("SERVICE_NAME", "survey_taking_service"),
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

// Helper function to get boolean environment variables with fallback
func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
