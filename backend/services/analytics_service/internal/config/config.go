package config

import (
	"fmt"
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
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	ResponsesColl string
	SurveysColl   string
	AuthSource    string
	URI           string
}

// Auth holds the configuration for authentication
type Auth struct {
	JWTSecret      string
	AuthServiceURL string
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
		Server: Server{
			Port: getEnv("SERVER_PORT", "8084"),
		},
		MongoDB: MongoDB{
			Host:          mongoHost,
			Port:          mongoPort,
			User:          mongoUser,
			Password:      mongoPassword,
			Database:      mongoDatabase,
			URI:           mongoURI,
			ResponsesColl: getEnv("MONGODB_RESPONSES_COLLECTION", "responses"),
			SurveysColl:   getEnv("MONGODB_SURVEYS_COLLECTION", "surveys"),
			AuthSource:    mongoAuthSource,
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
