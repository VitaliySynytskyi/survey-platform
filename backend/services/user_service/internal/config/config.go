package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the server
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	JWTSecretKey string
}

// ServerConfig holds all server configuration
type ServerConfig struct {
	Port string
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
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "survey_user"),
			Password: getEnv("DB_PASSWORD", "survey_password"),
			DBName:   getEnv("DB_NAME", "survey_platform"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWTSecretKey: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
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
