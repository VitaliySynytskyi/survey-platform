package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a middleware that validates JWT tokens
type AuthMiddleware struct {
	secretKey string
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

// Authenticate validates the JWT token in the Authorization header
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Parse and validate the token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract user ID and role from claims
		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Invalid token: missing subject claim", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Invalid token: missing role claim", http.StatusUnauthorized)
			return
		}

		// Add user ID and role to context
		ctx := r.Context()
		ctx = domain.ContextWithUserClaims(ctx, domain.UserClaims{
			ID:   userID,
			Role: domain.Role(role),
		})

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin ensures the user has admin role
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user role from context
		role, ok := domain.UserRoleFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user has admin role
		if role != domain.RoleAdmin {
			http.Error(w, "Forbidden: admin role required", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireOwnerOrAdmin ensures the user is either the owner of the resource (matches user ID in URL) or an admin
func (m *AuthMiddleware) RequireOwnerOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID and role from context
		userID, userOk := domain.UserIDFromContext(r.Context())
		role, roleOk := domain.UserRoleFromContext(r.Context())

		if !userOk || !roleOk {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If user is admin, allow access
		if role == domain.RoleAdmin {
			next.ServeHTTP(w, r)
			return
		}

		// Get requested user ID from URL path
		requestedUserID := r.PathValue("id") // Go 1.22+ feature, or use a router like chi to extract path params
		if requestedUserID == "" {
			http.Error(w, "User ID not provided", http.StatusBadRequest)
			return
		}

		// Check if user ID matches the requested ID
		if userID != requestedUserID {
			http.Error(w, "Forbidden: you can only access your own profile", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// validateToken validates a JWT token
func (m *AuthMiddleware) validateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token type (should be "access")
		if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "access" {
			return nil, errors.New("invalid token type")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
