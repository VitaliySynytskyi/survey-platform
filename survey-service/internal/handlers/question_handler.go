package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/survey-app/survey-service/internal/models"
	"github.com/survey-app/survey-service/internal/service"
)

// QuestionHandler handles HTTP requests related to questions
type QuestionHandler struct {
	surveyService service.SurveyServiceInterface
}

// NewQuestionHandler creates a new question handler
func NewQuestionHandler(surveyService service.SurveyServiceInterface) *QuestionHandler {
	return &QuestionHandler{
		surveyService: surveyService,
	}
}

// UpdateQuestion handles PUT /api/v1/questions/:id request
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var req models.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create question model from request
	question := &models.Question{
		ID:       id,
		SurveyID: req.SurveyID,
		Text:     req.Text,
		Type:     req.Type,
		Required: req.Required,
		OrderNum: req.OrderNum,
	}

	// Convert option requests to option models
	var options []*models.QuestionOption
	for _, optReq := range req.Options {
		options = append(options, &models.QuestionOption{
			QuestionID: id,
			Text:       optReq.Text,
			OrderNum:   optReq.OrderNum,
		})
	}

	if err := h.surveyService.UpdateQuestion(c.Request.Context(), question, options); err != nil {
		log.Printf("Error updating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

// DeleteQuestion handles DELETE /api/v1/questions/:id request
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.surveyService.DeleteQuestion(c.Request.Context(), id); err != nil {
		log.Printf("Error deleting question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
