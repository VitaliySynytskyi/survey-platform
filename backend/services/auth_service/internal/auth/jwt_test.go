package auth

import (
	"testing"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/google/uuid"
)

func TestTokenGeneration(t *testing.T) {
	// Create a token manager with test settings
	secretKey := "test-secret-key"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 24 * time.Hour
	tokenManager := NewTokenManager(secretKey, accessExpiry, refreshExpiry)

	// Create a test user
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  domain.RoleUser,
	}

	// Generate tokens
	tokenResponse, err := tokenManager.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Check if tokens are not empty
	if tokenResponse.AccessToken == "" {
		t.Error("Access token is empty")
	}
	if tokenResponse.RefreshToken == "" {
		t.Error("Refresh token is empty")
	}
	if tokenResponse.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}

	// Validate the access token
	claims, err := tokenManager.ValidateToken(tokenResponse.AccessToken, AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate access token: %v", err)
	}

	// Check claims
	if claims.UserID != user.ID.String() {
		t.Errorf("UserID mismatch: got %s, want %s", claims.UserID, user.ID.String())
	}
	if claims.Email != user.Email {
		t.Errorf("Email mismatch: got %s, want %s", claims.Email, user.Email)
	}
	if claims.Role != string(user.Role) {
		t.Errorf("Role mismatch: got %s, want %s", claims.Role, string(user.Role))
	}
}

func TestTokenValidation(t *testing.T) {
	// Create a token manager with test settings
	secretKey := "test-secret-key"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 24 * time.Hour
	tokenManager := NewTokenManager(secretKey, accessExpiry, refreshExpiry)

	// Create a test user
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  domain.RoleUser,
	}

	// Generate tokens
	tokenResponse, err := tokenManager.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test valid access token
	_, err = tokenManager.ValidateToken(tokenResponse.AccessToken, AccessToken)
	if err != nil {
		t.Errorf("Valid access token should be accepted: %v", err)
	}

	// Test valid refresh token
	_, err = tokenManager.ValidateToken(tokenResponse.RefreshToken, RefreshToken)
	if err != nil {
		t.Errorf("Valid refresh token should be accepted: %v", err)
	}

	// Test invalid token type
	_, err = tokenManager.ValidateToken(tokenResponse.AccessToken, RefreshToken)
	if err == nil {
		t.Error("Access token should not be accepted as refresh token")
	}

	// Test invalid token
	_, err = tokenManager.ValidateToken("invalid-token", AccessToken)
	if err == nil {
		t.Error("Invalid token should be rejected")
	}
}

func TestRefreshTokens(t *testing.T) {
	// Create a token manager with test settings
	secretKey := "test-secret-key"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 24 * time.Hour
	tokenManager := NewTokenManager(secretKey, accessExpiry, refreshExpiry)

	// Create a test user
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  domain.RoleUser,
	}

	// Generate tokens
	tokenResponse, err := tokenManager.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Refresh tokens
	newTokenResponse, err := tokenManager.RefreshTokens(tokenResponse.RefreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh tokens: %v", err)
	}

	// Check if new tokens are not empty
	if newTokenResponse.AccessToken == "" {
		t.Error("New access token is empty")
	}
	if newTokenResponse.RefreshToken == "" {
		t.Error("New refresh token is empty")
	}

	// Check if new tokens are different from old ones
	if newTokenResponse.AccessToken == tokenResponse.AccessToken {
		t.Error("New access token should be different from old one")
	}
	if newTokenResponse.RefreshToken == tokenResponse.RefreshToken {
		t.Error("New refresh token should be different from old one")
	}
}
