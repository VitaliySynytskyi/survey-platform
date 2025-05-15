package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
)

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// UserContext is the key for user information in the request context
type UserContext string

const (
	UserIDKey   UserContext = "userID"
	UserRoleKey UserContext = "userRole"
)

// JWTMiddleware creates middleware for JWT authentication
func JWTMiddleware(cfg config.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			// Parse and validate the token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID gets the user ID from the request context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserRole gets the user role from the request context
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

// IsAdmin checks if the user has admin role
func IsAdmin(ctx context.Context) bool {
	role, ok := GetUserRole(ctx)
	return ok && role == "admin"
}

// CheckSurveyOwnership checks if the user is the owner of the survey
func CheckSurveyOwnership(ctx context.Context, survey *struct {
	OwnerID string `json:"owner_id"`
}) bool {
	userID, ok := GetUserID(ctx)
	if !ok {
		return false
	}

	// Admin can access any survey
	if IsAdmin(ctx) {
		return true
	}

	// Check if the user is the owner
	return userID == survey.OwnerID
}
