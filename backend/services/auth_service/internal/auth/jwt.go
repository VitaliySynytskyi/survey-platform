package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents the JWT claims
type Claims struct {
	UserID string    `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenType specifies the token type (access or refresh)
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenManager handles JWT token generation and validation
type TokenManager struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokenPair generates a new access and refresh token pair
func (tm *TokenManager) GenerateTokenPair(user *domain.User) (*domain.TokenResponse, error) {
	// Generate access token
	accessToken, accessExpiry, err := tm.generateToken(user, AccessToken, tm.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := tm.generateToken(user, RefreshToken, tm.refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Calculate expiry in seconds
	expiresIn := int64(time.Until(accessExpiry).Seconds())

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// generateToken creates a new token of the specified type
func (tm *TokenManager) generateToken(user *domain.User, tokenType TokenType, expiry time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiryTime := now.Add(expiry)

	claims := Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "survey-platform-auth-service",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiryTime, nil
}

// ValidateToken validates the provided token and returns the claims
func (tm *TokenManager) ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract and validate the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify token type
	if claims.Type != tokenType {
		return nil, fmt.Errorf("%w: token type mismatch", ErrInvalidToken)
	}

	return claims, nil
}

// ExtractUserFromClaims creates a user from token claims
func (tm *TokenManager) ExtractUserFromClaims(claims *Claims) (*domain.User, error) {
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	// Create a minimal user object with the information from the token
	user := &domain.User{
		ID:    userID,
		Email: claims.Email,
		Role:  domain.Role(claims.Role),
	}

	return user, nil
}

// RefreshTokens validates a refresh token and generates a new token pair
func (tm *TokenManager) RefreshTokens(refreshTokenString string) (*domain.TokenResponse, error) {
	// Validate the refresh token
	claims, err := tm.ValidateToken(refreshTokenString, RefreshToken)
	if err != nil {
		return nil, err
	}

	// Extract user from claims
	user, err := tm.ExtractUserFromClaims(claims)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return tm.GenerateTokenPair(user)
}
