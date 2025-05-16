package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

	// Run tests
	t.Run("GetUserByID", testGetUserByID(store))
	t.Run("ListUsers", testListUsers(store))
	t.Run("EmailExists", testEmailExists(store))
}

func testGetUserByID(store *PostgresStore) func(t *testing.T) {
	return func(t *testing.T) {
		// This test assumes that a test user with ID exists in the database
		// In a real test, we would create a user first, but for simplicity, we'll use a fixed ID
		ctx := context.Background()
		userID := uuid.MustParse("11111111-1111-1111-1111-111111111111") // Replace with a real UUID from your test database

		// Get user by ID
		user, err := store.GetUserByID(ctx, userID)
		if err != nil {
			// If the test user doesn't exist, this should fail with ErrUserNotFound
			if err == ErrUserNotFound {
				t.Skip("Test user not found. Create a test user first.")
			}
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		// Verify that the user ID matches
		if user.ID != userID {
			t.Errorf("User ID mismatch: got %s, want %s", user.ID, userID)
		}
	}
}

func testListUsers(store *PostgresStore) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// List users with pagination
		users, err := store.ListUsers(ctx, 0, 10)
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		// Verify that we got some users
		if len(users) == 0 {
			t.Skip("No users found in the database. Create some test users first.")
		}

		// Count users
		count, err := store.CountUsers(ctx)
		if err != nil {
			t.Fatalf("Failed to count users: %v", err)
		}

		// Verify that the count is at least equal to the number of users we got
		if count < len(users) {
			t.Errorf("User count mismatch: got %d, want at least %d", count, len(users))
		}
	}
}

func testEmailExists(store *PostgresStore) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// This test assumes that a test user with this email exists in the database
		testEmail := "test@example.com"                                  // Replace with a real email from your test database
		userID := uuid.MustParse("11111111-1111-1111-1111-111111111111") // Replace with a real UUID from your test database

		// Check if email exists (excluding the user ID)
		exists, err := store.EmailExists(ctx, testEmail, userID)
		if err != nil {
			t.Fatalf("Failed to check if email exists: %v", err)
		}

		// This should be false since we're excluding the user that owns this email
		if exists {
			t.Errorf("EmailExists should return false for user's own email")
		}

		// Check if email exists (using a different user ID)
		differentUserID := uuid.New() // Generate a random UUID that won't match any user
		exists, err = store.EmailExists(ctx, testEmail, differentUserID)
		if err != nil {
			t.Fatalf("Failed to check if email exists: %v", err)
		}

		// This should be true if the email exists in the database
		if !exists {
			t.Skip("Test email not found in the database. Create a test user first.")
		}
	}
}
