package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/models"
	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/repository"
)

// AuthService handles authentication operations
type AuthService struct {
	repo            repository.Repository
	jwtSecret       string
	expirationHours int
}

// NewAuthService creates a new AuthService instance
func NewAuthService(repo repository.Repository, jwtSecret string, expirationHours int) *AuthService {
	return &AuthService{
		repo:            repo,
		jwtSecret:       jwtSecret,
		expirationHours: expirationHours,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if username already exists
	existingUser, _ := s.repo.GetUserByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	passwordHash, err := models.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	// Save user to database
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Get 'user' role and assign it to the user
	role, err := s.repo.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("error getting user role: %w. Ensure 'user' role exists in the database.", err)
	}

	// Assign role to user
	err = s.repo.AddUserRole(ctx, userID, role.ID)
	if err != nil {
		return nil, fmt.Errorf("error assigning role to user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Get user with roles
	userWithRoles, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user with roles: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *userWithRoles,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !models.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Parse the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	// Check if it's a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Get user ID from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID := int(userIDFloat)

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	return &models.AuthResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		User:         *user,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check if it's an access token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// generateTokens generates access and refresh tokens for a user
func (s *AuthService) generateTokens(user *models.User) (string, string, error) {
	// Create access token
	accessTokenClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"roles":    user.Roles,
		"type":     "access",
		"exp":      time.Now().Add(time.Hour * time.Duration(s.expirationHours)).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Create refresh token (valid for 30 days)
	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
