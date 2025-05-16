package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role type for user roles
type Role string

// Role constants
const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Role         Role      `json:"role"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserProfile represents the data that can be updated by a user
type UserProfile struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UserResponse is the public representation of a user
// We use this to avoid exposing sensitive fields in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// RoleUpdate represents a request to update a user's role
type RoleUpdate struct {
	Role Role `json:"role"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}

// ApplyUpdate applies the profile updates to the user
func (u *User) ApplyUpdate(profile UserProfile) error {
	if profile.Email != "" && profile.Email != u.Email {
		u.Email = profile.Email
	}

	if profile.Password != "" {
		hashedPassword, err := HashPassword(profile.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = hashedPassword
	}

	if profile.FirstName != "" {
		u.FirstName = profile.FirstName
	}

	if profile.LastName != "" {
		u.LastName = profile.LastName
	}

	u.UpdatedAt = time.Now().UTC()
	return nil
}

// SetRole updates the user's role
func (u *User) SetRole(role Role) {
	u.Role = role
	u.UpdatedAt = time.Now().UTC()
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword checks if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Validate validates the profile update
func (p *UserProfile) Validate() map[string]string {
	errors := make(map[string]string)

	if p.Email != "" {
		// Add email validation if needed
		// For example, check if it matches an email pattern
	}

	if p.Password != "" && len(p.Password) < 6 {
		errors["password"] = "Password must be at least 6 characters long"
	}

	return errors
}

// Validate validates the role update
func (r *RoleUpdate) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Role != RoleUser && r.Role != RoleAdmin {
		errors["role"] = "Invalid role"
	}

	return errors
}
