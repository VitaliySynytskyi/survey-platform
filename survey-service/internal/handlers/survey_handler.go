package handlers

import (
	"log"
	"net/http"
	"strconv"

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

// CreateSurvey handles POST /api/v1/surveys request
func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req models.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get creator ID from the JWT token
	creatorID := 1 // Temporary hardcoded value

	survey := &models.Survey{
		CreatorID:   creatorID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	// Pass the requested questions (req.Questions) to the service method
	id, err := h.surveyService.CreateSurvey(c.Request.Context(), survey, req.Questions)
	if err != nil {
		log.Printf("Error creating survey with questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create survey"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Survey created successfully"})
}

// GetSurveys handles GET /api/v1/surveys request
func (h *SurveyHandler) GetSurveys(c *gin.Context) {
	// TODO: Get creator ID from the JWT token
	creatorID := 1 // Temporary hardcoded value

	surveys, err := h.surveyService.GetSurveys(c.Request.Context(), creatorID)
	if err != nil {
		log.Printf("Error getting surveys: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve surveys"})
		return
	}

	if surveys == nil {
		c.JSON(http.StatusOK, []models.Survey{}) // Return empty array instead of null
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

	survey, err := h.surveyService.GetSurvey(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error getting survey: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve survey"})
		return
	}

	if survey == nil {
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

	// Get the existing survey details (primarily for checking existence and creator ID)
	existingSurvey, err := h.surveyService.GetSurvey(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error getting survey for update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve survey for update"})
		return
	}
	if existingSurvey == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}

	// TODO: Validate creator ID from JWT matches existingSurvey.CreatorID

	// Prepare the survey model for update (basic fields)
	surveyToUpdate := &models.Survey{
		ID:          id,
		CreatorID:   existingSurvey.CreatorID, // Keep original creator
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		// Questions will be handled by the service layer
	}

	// The service layer will handle the logic for updating/creating questions and options.
	// Pass the survey data and the requested questions to the service.
	if err := h.surveyService.UpdateSurveyWithQuestions(c.Request.Context(), surveyToUpdate, req.Questions); err != nil {
		log.Printf("Error updating survey with questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update survey"})
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

	if err := h.surveyService.DeleteSurvey(c.Request.Context(), id); err != nil {
		log.Printf("Error deleting survey: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete survey"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey deleted successfully"})
}

// AddQuestion handles POST /api/v1/surveys/:id/questions request
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

	// Set survey ID from URL parameter
	req.SurveyID = surveyID

	questionID, err := h.surveyService.AddQuestion(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Error adding question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add question"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": questionID, "message": "Question added successfully"})
}
