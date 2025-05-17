package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the server
type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	JWTSecretKey     string
	Consul           ConsulConfig
	ServiceName      string
	SurveyServiceURL string
}

// ConsulConfig holds all Consul configuration
type ConsulConfig struct {
	Address string
	Enabled bool
}

// ServerConfig holds all server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds all database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8081"), // Different port from auth_service
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "survey_user"),
			Password: getEnv("DB_PASSWORD", "survey_password"),
			DBName:   getEnv("DB_NAME", "survey_platform"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWTSecretKey: getEnv("JWT_SECRET_KEY", "your_jwt_secret_key"),
		Consul: ConsulConfig{
			Address: getEnv("CONSUL_ADDR", "consul:8500"),
			Enabled: getBoolEnv("CONSUL_ENABLED", true),
		},
		ServiceName:      getEnv("SERVICE_NAME", "user_service"),
		SurveyServiceURL: getEnv("SURVEY_SERVICE_URL", "http://survey_service:8082"),
	}
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an integer from an environment variable
func getIntEnv(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Warning: Invalid integer for %s. Using default value.", key)
		return defaultValue
	}

	return val
}

// Helper function to get a boolean from an environment variable
func getBoolEnv(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		log.Printf("Warning: Invalid boolean for %s. Using default value.", key)
		return defaultValue
	}

	return val
}
