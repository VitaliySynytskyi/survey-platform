package contextkeys

// ContextKey defines the type for context keys to ensure type safety.
type ContextKey string

const (
	// UserIDKey is the context key for the user's ID.
	UserIDKey ContextKey = "userID"
	// UserRolesKey is the context key for the user's roles.
	UserRolesKey ContextKey = "userRoles"
)
