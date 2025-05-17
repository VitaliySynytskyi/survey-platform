package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"log"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/store"
)

// Handler encapsulates all HTTP handlers
type Handler struct {
	userStore    store.UserStore
	tokenManager *auth.TokenManager
}

// NewHandler creates a new Handler
func NewHandler(userStore store.UserStore, tokenManager *auth.TokenManager) *Handler {
	return &Handler{
		userStore:    userStore,
		tokenManager: tokenManager,
	}
}

// RegisterHandler handles user registration
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// Normalize email (convert to lowercase)
	req.Email = strings.ToLower(req.Email)

	// Check if email already exists
	exists, err := h.userStore.EmailExists(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]string{
				"email": "Email already taken",
			},
		})
		return
	}

	// Create new user
	user, err := domain.NewUser(req.Email, req.Password, domain.RoleUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Save user to database
	if err := h.userStore.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user.ToResponse(),
	})
}

// LoginHandler handles user login
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// Normalize email (convert to lowercase)
	req.Email = strings.ToLower(req.Email)

	// Find user by email
	user, err := h.userStore.FindUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid email or password",
		})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate token pair
	tokenResponse, err := h.tokenManager.GenerateTokenPair(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

// RefreshHandler handles token refresh
func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// Refresh tokens
	tokenResponse, err := h.tokenManager.RefreshTokens(req.RefreshToken)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Invalid refresh token",
			})
			return
		}
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

// LogoutHandler handles user logout
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RefreshRequest // Assuming logout sends refresh token for invalidation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If body is empty or not JSON, it's a bad request or we can proceed if RT is optional
		// For now, let's assume RT is expected for proper invalidation.
		http.Error(w, "Invalid request body: refresh_token expected for logout", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required for logout", http.StatusBadRequest)
		return
	}

	// Attempt to invalidate the refresh token
	err := h.tokenManager.InvalidateRefreshToken(req.RefreshToken)
	if err != nil {
		// Log the error but don't necessarily send a 500, as logout should ideally always succeed on client-side
		// unless the token was already invalid, which is fine.
		// Depending on the error type from InvalidateRefreshToken, we might adjust.
		// For example, if ErrInvalidToken is returned, it's not a server error.
		log.Printf("Error invalidating refresh token during logout: %v", err)
		// We can choose to return 200/204 even if token was already invalid or not found for invalidation
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

// MeHandler handles getting current user information
func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (set by authentication middleware)
	user, ok := domain.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ToResponse())
}
