package config

import (
	"fmt"
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
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	AuthSource string
	URI        string
}

// ConsulConfig holds the configuration for Consul
type ConsulConfig struct {
	Address string
	Enabled bool
}

// Load loads the configuration from environment variables
func Load() *Config {
	mongoHost := getEnv("MONGO_HOST", "localhost")
	mongoPort := getEnv("MONGO_PORT", "27017")
	mongoUser := getEnv("MONGO_USER", "")
	mongoPassword := getEnv("MONGO_PASSWORD", "")
	mongoDatabase := getEnv("MONGO_DB", "survey_platform")
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
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8083"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:       getEnv("RABBITMQ_HOST", "localhost"),
			Port:       getEnv("RABBITMQ_PORT", "5672"),
			Username:   getEnv("RABBITMQ_USER", "guest"),
			Password:   getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:   getEnv("RABBITMQ_EXCHANGE", "survey_responses"),
			Queue:      getEnv("RABBITMQ_QUEUE", "survey_responses_queue"),
			RoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "survey.response"),
		},
		MongoDB: MongoDBConfig{
			Host:       mongoHost,
			Port:       mongoPort,
			User:       mongoUser,
			Password:   mongoPassword,
			Database:   mongoDatabase,
			AuthSource: mongoAuthSource,
			URI:        mongoURI,
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
