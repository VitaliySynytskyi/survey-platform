package middleware

import (
	"net/http"
	"strings"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
)

// AuthMiddleware is a middleware that checks for a valid JWT token
type AuthMiddleware struct {
	tokenManager *auth.TokenManager
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(tokenManager *auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

// Authenticate is a middleware that checks for a valid JWT access token
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

		// Validate the token
		claims, err := m.tokenManager.ValidateToken(tokenString, auth.AccessToken)
		if err != nil {
			if err == auth.ErrExpiredToken {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}

		// Extract user from claims
		user, err := m.tokenManager.ExtractUserFromClaims(claims)
		if err != nil {
			http.Error(w, "Failed to extract user from token", http.StatusInternalServerError)
			return
		}

		// Add user to request context
		ctx := r.Context()
		ctx = domain.ContextWithUser(ctx, user)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole is a middleware that checks if the authenticated user has the required role
func (m *AuthMiddleware) RequireRole(role domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (already set by Authenticate middleware)
			user, ok := domain.UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has the required role
			if user.Role != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
