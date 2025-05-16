package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the analytics service
type Config struct {
	Server  Server
	MongoDB MongoDB
	Auth    Auth
}

// Server holds the configuration for the HTTP server
type Server struct {
	Port string
}

// MongoDB holds the configuration for MongoDB
type MongoDB struct {
	URI           string
	Database      string
	ResponsesColl string
	SurveysColl   string
}

// Auth holds the configuration for authentication
type Auth struct {
	JWTSecret      string
	AuthServiceURL string
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		Server: Server{
			Port: getEnv("SERVER_PORT", "8085"),
		},
		MongoDB: MongoDB{
			URI:           getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:      getEnv("MONGODB_DATABASE", "survey_platform"),
			ResponsesColl: getEnv("MONGODB_RESPONSES_COLLECTION", "responses"),
			SurveysColl:   getEnv("MONGODB_SURVEYS_COLLECTION", "surveys"),
		},
		Auth: Auth{
			JWTSecret:      getEnv("JWT_SECRET", "your-default-secret-key"),
			AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
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
