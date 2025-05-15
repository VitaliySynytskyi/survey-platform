package store

import (
	"context"
	"testing"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
)

// TestPostgresStore is an integration test that requires a real PostgreSQL database
// These tests should be skipped during normal unit testing and run only in CI/CD pipeline
// or when explicitly requested
func TestPostgresStore(t *testing.T) {
	// Skip this test by default
	t.Skip("Skipping PostgreSQL integration test. Run with -tags=integration to enable.")

	// Connection string for test database
	connString := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"

	// Create store
	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Ensure schema
	err = store.EnsureSchema(context.Background())
	if err != nil {
		t.Fatalf("Failed to ensure schema: %v", err)
	}

	// Run tests
	t.Run("CreateAndFindUser", testCreateAndFindUser(store))
	t.Run("EmailExists", testEmailExists(store))
}

func testCreateAndFindUser(store *PostgresStore) func(t *testing.T) {
	return func(t *testing.T) {
		// Create a test user
		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		user, err := domain.NewUser(email, password, domain.RoleUser)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Save user to database
		err = store.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to save user: %v", err)
		}

		// Find user by email
		foundUser, err := store.FindUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("Failed to find user by email: %v", err)
		}
		if foundUser == nil {
			t.Fatal("User not found by email")
		}
		if foundUser.Email != email {
			t.Errorf("Email mismatch: got %s, want %s", foundUser.Email, email)
		}

		// Find user by ID
		foundUserById, err := store.FindUserByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to find user by ID: %v", err)
		}
		if foundUserById == nil {
			t.Fatal("User not found by ID")
		}
		if foundUserById.ID != user.ID {
			t.Errorf("ID mismatch: got %s, want %s", foundUserById.ID, user.ID)
		}
	}
}

func testEmailExists(store *PostgresStore) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// Should find existing email
		existingEmail := "test@example.com"
		exists, err := store.EmailExists(ctx, existingEmail)
		if err != nil {
			t.Fatalf("Failed to check if email exists: %v", err)
		}
		if !exists {
			t.Errorf("Email %s should exist", existingEmail)
		}

		// Should not find non-existing email
		nonExistingEmail := "nonexisting@example.com"
		exists, err = store.EmailExists(ctx, nonExistingEmail)
		if err != nil {
			t.Fatalf("Failed to check if email exists: %v", err)
		}
		if exists {
			t.Errorf("Email %s should not exist", nonExistingEmail)
		}
	}
}
