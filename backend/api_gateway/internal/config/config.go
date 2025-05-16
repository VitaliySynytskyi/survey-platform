package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the API Gateway configuration
type Config struct {
	Server         ServerConfig
	Services       ServicesConfig
	CORS           CORSConfig
	RateLimiting   RateLimitingConfig
	CircuitBreaker CircuitBreakerConfig
	Consul         ConsulConfig
}

// ServerConfig holds the HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// ConsulConfig holds the configuration for Consul service discovery
type ConsulConfig struct {
	Enabled  bool
	Address  string
	UseForSD bool // Whether to use Consul for service discovery
}

// ServicesConfig holds the configuration for different services
type ServicesConfig struct {
	Auth              ServiceConfig
	User              ServiceConfig
	Survey            ServiceConfig
	SurveyTaking      ServiceConfig
	ResponseProcessor ServiceConfig
	Analytics         ServiceConfig
}

// ServiceConfig represents a microservice configuration
type ServiceConfig struct {
	URL        string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

// CORSConfig holds the CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// RateLimitingConfig holds rate limiting configuration
type RateLimitingConfig struct {
	Enabled   bool
	Limit     int
	Burst     int
	TimeFrame time.Duration
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled       bool
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	TripThreshold float64
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 15*time.Second),
		},
		Consul: ConsulConfig{
			Enabled:  getEnvBool("CONSUL_ENABLED", true),
			Address:  getEnv("CONSUL_ADDR", "consul:8500"),
			UseForSD: getEnvBool("USE_CONSUL_FOR_SD", true),
		},
		Services: ServicesConfig{
			Auth: ServiceConfig{
				URL:        getEnv("AUTH_SERVICE_URL", "http://auth_service:8081"),
				Timeout:    getEnvDuration("AUTH_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("AUTH_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("AUTH_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
			User: ServiceConfig{
				URL:        getEnv("USER_SERVICE_URL", "http://user_service:8082"),
				Timeout:    getEnvDuration("USER_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("USER_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("USER_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
			Survey: ServiceConfig{
				URL:        getEnv("SURVEY_SERVICE_URL", "http://survey_service:8083"),
				Timeout:    getEnvDuration("SURVEY_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("SURVEY_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("SURVEY_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
			SurveyTaking: ServiceConfig{
				URL:        getEnv("SURVEY_TAKING_SERVICE_URL", "http://survey_taking_service:8084"),
				Timeout:    getEnvDuration("SURVEY_TAKING_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("SURVEY_TAKING_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("SURVEY_TAKING_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
			ResponseProcessor: ServiceConfig{
				URL:        getEnv("RESPONSE_PROCESSOR_SERVICE_URL", "http://response_processor_service:8085"),
				Timeout:    getEnvDuration("RESPONSE_PROCESSOR_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("RESPONSE_PROCESSOR_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("RESPONSE_PROCESSOR_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
			Analytics: ServiceConfig{
				URL:        getEnv("ANALYTICS_SERVICE_URL", "http://analytics_service:8086"),
				Timeout:    getEnvDuration("ANALYTICS_SERVICE_TIMEOUT", 5*time.Second),
				MaxRetries: getEnvInt("ANALYTICS_SERVICE_MAX_RETRIES", 3),
				RetryDelay: getEnvDuration("ANALYTICS_SERVICE_RETRY_DELAY", 100*time.Millisecond),
			},
		},
		CORS: CORSConfig{
			AllowedOrigins:   getEnvStringSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   getEnvStringSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvStringSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
			ExposedHeaders:   getEnvStringSlice("CORS_EXPOSED_HEADERS", []string{"Link"}),
			AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvInt("CORS_MAX_AGE", 300),
		},
		RateLimiting: RateLimitingConfig{
			Enabled:   getEnvBool("RATE_LIMITING_ENABLED", false),
			Limit:     getEnvInt("RATE_LIMITING_LIMIT", 100),
			Burst:     getEnvInt("RATE_LIMITING_BURST", 150),
			TimeFrame: getEnvDuration("RATE_LIMITING_TIMEFRAME", time.Minute),
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:       getEnvBool("CIRCUIT_BREAKER_ENABLED", true),
			MaxRequests:   uint32(getEnvInt("CIRCUIT_BREAKER_MAX_REQUESTS", 5)),
			Interval:      getEnvDuration("CIRCUIT_BREAKER_INTERVAL", 30*time.Second),
			Timeout:       getEnvDuration("CIRCUIT_BREAKER_TIMEOUT", 60*time.Second),
			TripThreshold: getEnvFloat("CIRCUIT_BREAKER_TRIP_THRESHOLD", 0.5),
		},
	}
}

// Helper functions to get environment variables with default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}
