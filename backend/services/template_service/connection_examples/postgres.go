package connection_examples

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresConfig holds the connection details
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectToPostgres establishes a connection to PostgreSQL
func ConnectToPostgres(config PostgresConfig) (*sql.DB, error) {
	// Connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check if connection is working
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// Example usage:
func PostgresExample() {
	config := PostgresConfig{
		Host:     "localhost", // or "postgres" when using Docker network
		Port:     "5432",
		User:     "survey_user",
		Password: "survey_password",
		DBName:   "survey_platform",
		SSLMode:  "disable", // or "require" in production
	}

	db, err := ConnectToPostgres(config)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Here you would use the connection to execute queries
	// Example:
	// rows, err := db.Query("SELECT id, name FROM users")
}
