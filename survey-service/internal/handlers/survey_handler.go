package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/survey-app/survey-service/internal/models"
	"github.com/survey-app/survey-service/internal/service"
)

// SurveyHandler handles HTTP requests related to surveys
type SurveyHandler struct {
	surveyService service.SurveyServiceInterface
}

// NewSurveyHandler creates a new survey handler
func NewSurveyHandler(surveyService service.SurveyServiceInterface) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
	}
}

// Helper function to get user context from Gin context
func getUserContext(c *gin.Context) (context.Context, error) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		log.Println("Error: X-User-ID header is missing")
		return nil, errors.New("X-User-ID header is missing")
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Error parsing X-User-ID header: %v", err)
		return nil, errors.New("invalid X-User-ID header")
	}

	rolesStr := c.GetHeader("X-User-Roles")
	var roles []string
	if rolesStr != "" {
		// Expecting rolesStr to be like "[role1 role2]"
		if len(rolesStr) > 2 && rolesStr[0] == '[' && rolesStr[len(rolesStr)-1] == ']' {
			rolesContent := rolesStr[1 : len(rolesStr)-1]
			if rolesContent != "" { // handle case like "[]"
				roles = strings.Split(rolesContent, " ")
			} else {
				roles = []string{} // empty slice for "[]"
			}
		} else {
			// If not in expected bracketed format, maybe it's a single role or unformatted.
			// This part might need adjustment based on actual header format from API gateway if it changes.
			// For now, let's log a warning and proceed with a single role if not bracketed.
			log.Printf("Warning: X-User-Roles header ('%s') not in expected bracketed format. Treating as single role or space-separated if no brackets.", rolesStr)
			roles = strings.Split(rolesStr, " ") // simple split if no brackets
		}
	} else {
		log.Println("Warning: X-User-Roles header is missing or empty")
		roles = []string{} // No roles, or handle as an error if roles are strictly required
	}

	// Filter out empty strings that might result from multiple spaces
	var cleanRoles []string
	for _, r := range roles {
		if r != "" {
			cleanRoles = append(cleanRoles, r)
		}
	}

	// Use exported keys from service package
	ctx := context.WithValue(c.Request.Context(), service.UserIDKey, userID)
	ctx = context.WithValue(ctx, service.UserRolesKey, cleanRoles)
	return ctx, nil
}

// CreateSurvey handles POST /api/v1/surveys request
func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req models.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// CreatorID will be set by the service layer from the context

	survey := &models.Survey{
		// CreatorID is now handled by the service using context
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	id, err := h.surveyService.CreateSurvey(userCtx, survey, req.Questions)
	if err != nil {
		log.Printf("Error creating survey with questions: %v", err)
		// Check for specific error types if service layer provides them (e.g. forbidden)
		if strings.Contains(err.Error(), "forbidden") { // Basic check
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create survey"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Survey created successfully"})
}

// GetSurveys handles GET /api/v1/surveys request
func (h *SurveyHandler) GetSurveys(c *gin.Context) {
	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// creatorID is no longer passed directly; service uses context

	surveys, err := h.surveyService.GetSurveys(userCtx) // Pass userCtx, creatorID removed
	if err != nil {
		log.Printf("Error getting surveys: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve surveys"})
		return
	}

	if surveys == nil {
		c.JSON(http.StatusOK, []models.Survey{})
		return
	}

	c.JSON(http.StatusOK, surveys)
}

// GetSurvey handles GET /api/v1/surveys/:id request
func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	// For GetSurvey, user context might be needed if non-active surveys have restricted access
	// Assuming for now GetSurvey is public for active surveys, or service handles auth
	userCtx, err := getUserContext(c) // Get context for potential auth in service
	if err != nil {
		// If getting context fails, it implies an issue with required headers for authenticated access.
		// If this endpoint is truly public for some cases, this error might be too strict.
		// For now, let's assume if headers are present, they should be valid.
		// If this endpoint can be hit anonymously, getUserContext needs to be more lenient or not called.
		// Given our API Gateway changes, most survey routes require auth.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication context error: " + err.Error()})
		return
	}

	survey, err := h.surveyService.GetSurvey(userCtx, id) // Pass userCtx
	if err != nil {
		log.Printf("Error getting survey: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		} else if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve survey"})
		}
		return
	}

	// surveyService.GetSurvey should return nil, err if not found (or error for forbidden)
	// The check for survey == nil might be redundant if errors are handled properly
	if survey == nil { // This case should ideally be caught by specific errors from service
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}

	c.JSON(http.StatusOK, survey)
}

// UpdateSurvey handles PUT /api/v1/surveys/:id request
func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	var req models.UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// The service layer will handle fetching existing survey and checking ownership/admin rights
	// CreatorID and original survey details are handled by the service.
	surveyToUpdate := &models.Survey{
		ID: id,
		// CreatorID will be validated by service against context userID or admin role
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.surveyService.UpdateSurveyWithQuestions(userCtx, surveyToUpdate, req.Questions); err != nil {
		log.Printf("Error updating survey with questions: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		} else if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update survey"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey updated successfully"})
}

// DeleteSurvey handles DELETE /api/v1/surveys/:id request
func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.surveyService.DeleteSurvey(userCtx, id); err != nil {
		log.Printf("Error deleting survey: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		} else if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete survey"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey deleted successfully"})
}

// UpdateSurveyStatus handles PATCH /api/v1/surveys/:id/status request
func (h *SurveyHandler) UpdateSurveyStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	var req models.UpdateSurveyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IsActive == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_active field is required"})
		return
	}

	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.surveyService.UpdateSurveyStatus(userCtx, id, *req.IsActive); err != nil {
		log.Printf("Error updating survey status: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		} else if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update survey status"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey status updated successfully"})
}

// AddQuestion handles POST /api/v1/surveys/:id/questions request
// This handler needs to ensure the user adding a question is the survey owner or admin.
// The service.AddQuestion method will need to perform this check.
func (h *SurveyHandler) AddQuestion(c *gin.Context) {
	surveyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	var req models.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCtx, err := getUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.SurveyID = surveyID // SurveyID is from path param

	questionID, err := h.surveyService.AddQuestion(userCtx, &req) // Pass userCtx
	if err != nil {
		log.Printf("Error adding question: %v", err)
		if strings.Contains(err.Error(), "not found") { // e.g. survey not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found or not authorized"})
		} else if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden to add question to this survey"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add question"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": questionID, "message": "Question added successfully"})
}
