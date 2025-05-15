package store

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/google/uuid"
)

// UserStore defines the interface for user storage
type UserStore interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *domain.User) error

	// FindUserByEmail finds a user by email
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// FindUserByID finds a user by ID
	FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// EmailExists checks if an email already exists
	EmailExists(ctx context.Context, email string) (bool, error)

	// EnsureSchema ensures that the required database schema exists
	EnsureSchema(ctx context.Context) error

	// Close closes the database connection
	Close()
}
