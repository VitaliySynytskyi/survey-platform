package domain

import (
	"context"
)

type contextKey string

const (
	userContextKey contextKey = "user"
	userIDKey      contextKey = "user_id"
	userRoleKey    contextKey = "user_role"
)

// UserClaims contains the essential user information from JWT claims
type UserClaims struct {
	ID   string
	Role Role
}

// ContextWithUser adds the user to the context
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves the user from the context
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// ContextWithUserClaims adds user claims to the context
func ContextWithUserClaims(ctx context.Context, claims UserClaims) context.Context {
	ctx = context.WithValue(ctx, userIDKey, claims.ID)
	ctx = context.WithValue(ctx, userRoleKey, claims.Role)
	return ctx
}

// UserIDFromContext retrieves the user ID from the context
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

// UserRoleFromContext retrieves the user role from the context
func UserRoleFromContext(ctx context.Context) (Role, bool) {
	role, ok := ctx.Value(userRoleKey).(Role)
	return role, ok
}
