package repository

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/models"
)

// Repository defines the interface for database operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error

	// Role operations
	GetUserRoles(ctx context.Context, userID int) ([]string, error)
	AddUserRole(ctx context.Context, userID, roleID int) error
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
}
