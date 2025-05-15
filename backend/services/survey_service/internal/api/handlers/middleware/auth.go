package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/VitaliySynytskyi/microservices-survey-app/backend/services/survey_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// UserContext зберігає дані користувача в контексті
type UserContext struct {
	UserID string
	Role   string
}

// ClaimsKey ключ для JWT claims у контексті
type ClaimsKey string

const (
	// UserContextKey ключ для даних користувача в контексті
	UserContextKey ClaimsKey = "user_context"
)

// AuthMiddleware middleware для перевірки JWT токена
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Отримання токена з заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Перевірка формату токена
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Парсинг JWT токена
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Перевірка алгоритму підпису
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(cfg.Auth.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Отримання даних з токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Перевіряємо необхідні поля
			userId, ok := claims["sub"].(string)
			if !ok || userId == "" {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			// Отримання ролі користувача (якщо є)
			role, _ := claims["role"].(string)

			// Додавання даних користувача до контексту
			ctx := context.WithValue(r.Context(), UserContextKey, UserContext{
				UserID: userId,
				Role:   role,
			})

			// Виклик наступного обробника з оновленим контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext отримує дані користувача з контексту
func GetUserFromContext(ctx context.Context) (UserContext, bool) {
	userCtx, ok := ctx.Value(UserContextKey).(UserContext)
	return userCtx, ok
}

// RequireRole middleware для перевірки ролі користувача
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if userCtx.Role != role && userCtx.Role != "admin" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
