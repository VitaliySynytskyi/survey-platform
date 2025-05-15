package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the server
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
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

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
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
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "survey_user"),
			Password: getEnv("DB_PASSWORD", "survey_password"),
			DBName:   getEnv("DB_NAME", "survey_platform"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTokenExp:  getDurationEnv("JWT_ACCESS_TOKEN_EXP", 15*time.Minute),
			RefreshTokenExp: getDurationEnv("JWT_REFRESH_TOKEN_EXP", 7*24*time.Hour), // 7 days
		},
	}
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get a duration from an environment variable
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	// Try to parse as seconds
	valInt, err := strconv.Atoi(valStr)
	if err == nil {
		return time.Duration(valInt) * time.Second
	}

	// Try to parse as duration string
	valDuration, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Warning: Invalid duration for %s. Using default value.", key)
		return defaultValue
	}

	return valDuration
}
