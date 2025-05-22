package config

import (
	"testing"
)

func TestConfigCreation(t *testing.T) {
	// Test the creation of a Config struct
	cfg := Config{
		DB: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "surveys",
		},
		Port: "8080",
	}

	// Verify field values
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.DB.Host != "localhost" {
		t.Errorf("Expected DB.Host to be 'localhost', got '%s'", cfg.DB.Host)
	}

	if cfg.DB.Port != "5432" {
		t.Errorf("Expected DB.Port to be '5432', got '%s'", cfg.DB.Port)
	}

	if cfg.DB.User != "postgres" {
		t.Errorf("Expected DB.User to be 'postgres', got '%s'", cfg.DB.User)
	}

	if cfg.DB.Password != "password" {
		t.Errorf("Expected DB.Password to be 'password', got '%s'", cfg.DB.Password)
	}

	if cfg.DB.Name != "surveys" {
		t.Errorf("Expected DB.Name to be 'surveys', got '%s'", cfg.DB.Name)
	}
}

func TestDBConfig(t *testing.T) {
	// Test the creation of a DBConfig struct
	dbConfig := DBConfig{
		Host:     "test-host",
		Port:     "1234",
		User:     "test-user",
		Password: "test-password",
		Name:     "test-db",
	}

	// Verify field values
	if dbConfig.Host != "test-host" {
		t.Errorf("Expected Host to be 'test-host', got '%s'", dbConfig.Host)
	}

	if dbConfig.Port != "1234" {
		t.Errorf("Expected Port to be '1234', got '%s'", dbConfig.Port)
	}

	if dbConfig.User != "test-user" {
		t.Errorf("Expected User to be 'test-user', got '%s'", dbConfig.User)
	}

	if dbConfig.Password != "test-password" {
		t.Errorf("Expected Password to be 'test-password', got '%s'", dbConfig.Password)
	}

	if dbConfig.Name != "test-db" {
		t.Errorf("Expected Name to be 'test-db', got '%s'", dbConfig.Name)
	}
}
