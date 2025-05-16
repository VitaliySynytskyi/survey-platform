package store

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/domain"
	"github.com/google/uuid"
)

// UserStore defines the interface for user storage
type UserStore interface {
	// GetUserByID finds a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetUserByEmail finds a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateUser updates a user in the database
	UpdateUser(ctx context.Context, user *domain.User) error

	// ListUsers returns a list of all users (with optional pagination)
	ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error)

	// CountUsers returns the total number of users
	CountUsers(ctx context.Context) (int, error)

	// EmailExists checks if an email already exists (excluding the given user ID)
	EmailExists(ctx context.Context, email string, excludeID uuid.UUID) (bool, error)

	// Ping checks if the database connection is alive
	Ping() error

	// Close closes the database connection
	Close()
}
