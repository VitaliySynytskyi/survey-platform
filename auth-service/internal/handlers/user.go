package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/survey-app/auth-service/internal/service"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetCurrentUser retrieves the current user's information
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get user from database
	user, err := h.service.GetUserByID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser updates the current user's information
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Parse request body
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	user, err := h.service.UpdateUser(c.Request.Context(), userID.(int), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
