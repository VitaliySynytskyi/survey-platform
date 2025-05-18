package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/service"
)

// ResponseHandler struct holds the service interface
type ResponseHandler struct {
	responseService service.ResponseServiceInterface
}

// NewResponseHandler creates a new ResponseHandler
func NewResponseHandler(rs service.ResponseServiceInterface) *ResponseHandler {
	return &ResponseHandler{responseService: rs}
}

// SubmitResponse handles POST requests to /responses
func (h *ResponseHandler) SubmitResponse(c *gin.Context) {
	var req models.CreateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}

	if err := h.responseService.SubmitResponse(c.Request.Context(), &req); err != nil {
		if strings.Contains(err.Error(), "not active") || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to submit response: %s", err.Error())})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Response submitted successfully"})
}

// GetSurveyResponsesHandler handles GET requests to /surveys/:surveyId/responses
func (h *ResponseHandler) GetSurveyResponsesHandler(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey_id format"})
		return
	}

	responses, err := h.responseService.GetSurveyResponses(c.Request.Context(), surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get survey responses: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// ExportSurveyResponsesCSV handles GET requests to /surveys/:surveyId/responses/export
// This is referenced in main.go, so its signature should be correct for gin.
func (h *ResponseHandler) ExportSurveyResponsesCSV(c *gin.Context) {
	// surveyIDStr := c.Param("surveyId") // Example of getting surveyId if needed
	// Placeholder logic: In a real implementation, you would fetch data,
	// format it as CSV, and set appropriate headers like Content-Disposition.
	c.String(http.StatusNotImplemented, "CSV export functionality is not fully implemented in this handler yet.")
}

// GetSurveyAnalytics handles GET requests to /surveys/{survey_id}/analytics
func (h *ResponseHandler) GetSurveyAnalytics(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Survey ID format in URL path"})
		return
	}

	analytics, err := h.responseService.GetSurveyAnalytics(c.Request.Context(), surveyID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Analytics not found or survey does not exist: %s", err.Error())})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get survey analytics: %s", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, analytics)
}
