package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/models"
	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/repository"
)

// UserService handles user operations
type UserService struct {
	repo repository.Repository
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, userID int, req struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}) (*models.User, error) {
	// Get current user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// If email is changing, check if it's already in use
	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)
		if existingUser != nil {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Save changes
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}
