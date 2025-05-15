package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/domain"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/store"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userStore store.UserStore
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// GetUserHandler handles GET requests to retrieve user information
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	// Parse user ID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find user by ID
	user, err := h.userStore.GetUserByID(r.Context(), parsedID)
	if err != nil {
		if err == store.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Return user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ToResponse())
}

// UpdateUserHandler handles PUT requests to update user information
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	// Parse user ID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var profile domain.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := profile.Validate(); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// Find user by ID
	user, err := h.userStore.GetUserByID(r.Context(), parsedID)
	if err != nil {
		if err == store.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Check if email is being updated and if it's already taken
	if profile.Email != "" && profile.Email != user.Email {
		// Normalize email (convert to lowercase)
		profile.Email = strings.ToLower(profile.Email)

		// Check if email is already taken
		exists, err := h.userStore.EmailExists(r.Context(), profile.Email, user.ID)
		if err != nil {
			http.Error(w, "Failed to check email uniqueness", http.StatusInternalServerError)
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
	}

	// Apply updates to user
	if err := user.ApplyUpdate(profile); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Save updated user to database
	if err := h.userStore.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Return updated user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ToResponse())
}

// ListUsersHandler handles GET requests to retrieve a list of users
func (h *UserHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get users
	users, err := h.userStore.ListUsers(r.Context(), offset, pageSize)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Get total count
	total, err := h.userStore.CountUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to count users", http.StatusInternalServerError)
		return
	}

	// Convert users to response format
	userResponses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	// Create response with pagination
	response := map[string]interface{}{
		"users": userResponses,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}

	// Return users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateRoleHandler handles PUT requests to update a user's role
func (h *UserHandler) UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	// Parse user ID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var roleUpdate domain.RoleUpdate
	if err := json.NewDecoder(r.Body).Decode(&roleUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if errors := roleUpdate.Validate(); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// Find user by ID
	user, err := h.userStore.GetUserByID(r.Context(), parsedID)
	if err != nil {
		if err == store.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Update user role
	user.SetRole(roleUpdate.Role)

	// Save updated user to database
	if err := h.userStore.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Return updated user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ToResponse())
}
