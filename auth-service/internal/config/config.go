package config

import (
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	DB  DBConfig
	JWT JWTConfig
}

// DBConfig represents the database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// New returns a new Config instance with values from environment variables
func New() *Config {
	// Set default JWT expiration to 24 hours if not specified
	expirationHours := 24
	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if exp, err := strconv.Atoi(expStr); err == nil {
			expirationHours = exp
		}
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "survey_db"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your_jwt_secret_key"),
			ExpirationHours: expirationHours,
		},
	}
}

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
