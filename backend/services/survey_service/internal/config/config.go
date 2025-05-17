package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config представляє конфігурацію сервісу
type Config struct {
	Server      ServerConfig
	MongoDB     MongoDBConfig
	Auth        AuthConfig
	Consul      ConsulConfig
	ServiceName string
}

// ServerConfig представляє конфігурацію веб-сервера
type ServerConfig struct {
	Port int
	Host string
}

// MongoDBConfig представляє конфігурацію MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// AuthConfig представляє конфігурацію автентифікації
type AuthConfig struct {
	JWTSecret string
}

// ConsulConfig представляє конфігурацію Consul
type ConsulConfig struct {
	Address string
	Enabled bool
}

// LoadConfig завантажує конфігурацію з середовища
func LoadConfig() (*Config, error) {
	// Завантаження змінних середовища з .env файлу, якщо він існує
	_ = godotenv.Load()

	// Конфігурація сервера
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8082"))
	if err != nil {
		return nil, fmt.Errorf("invalid server port: %w", err)
	}
	host := getEnv("SERVER_HOST", "0.0.0.0")

	// Конфігурація MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnv("MONGO_DATABASE", "survey_service")

	// Конфігурація автентифікації
	jwtSecret := getEnv("JWT_SECRET_KEY", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	// Конфігурація Consul
	consulEnabled, _ := strconv.ParseBool(getEnv("CONSUL_ENABLED", "true"))
	consulAddress := getEnv("CONSUL_ADDR", "consul:8500")
	serviceName := getEnv("SERVICE_NAME", "survey_service")

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
		MongoDB: MongoDBConfig{
			URI:      mongoURI,
			Database: mongoDatabase,
		},
		Auth: AuthConfig{
			JWTSecret: jwtSecret,
		},
		Consul: ConsulConfig{
			Address: consulAddress,
			Enabled: consulEnabled,
		},
		ServiceName: serviceName,
	}, nil
}

// getEnv отримує значення змінної середовища або значення за замовчуванням
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
