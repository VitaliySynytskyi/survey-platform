package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	ServerPort         string
	MongoDBURI         string
	MongoDBName        string
	ResponseCollection string
	SurveyServiceURL   string
}

// New creates a new Config instance
func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		ServerPort:         getEnv("RESPONSE_SERVICE_PORT", "8083"),
		MongoDBURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:        getEnv("MONGO_DATABASE", "survey_responses_db"),
		ResponseCollection: getEnv("MONGO_RESPONSE_COLLECTION", "responses"),
		SurveyServiceURL:   getEnv("SURVEY_SERVICE_URL", "http://survey-service:8082"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
