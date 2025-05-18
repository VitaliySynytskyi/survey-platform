package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/service"
)

// ResponseHandler handles HTTP requests related to survey responses
type ResponseHandler struct {
	responseService service.ResponseServiceInterface
}

// NewResponseHandler creates a new ResponseHandler
func NewResponseHandler(responseService service.ResponseServiceInterface) *ResponseHandler {
	return &ResponseHandler{
		responseService: responseService,
	}
}

// SubmitResponse handles POST /api/v1/responses
func (h *ResponseHandler) SubmitResponse(c *gin.Context) {
	var req models.CreateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get UserID from JWT token if available and req.UserID is nil

	if err := h.responseService.SubmitResponse(c.Request.Context(), &req); err != nil {
		// TODO: Differentiate between bad request (400) and server error (500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit response"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Response submitted successfully"})
}

// GetSurveyResponsesHandler handles GET /api/v1/surveys/:surveyId/responses
func (h *ResponseHandler) GetSurveyResponsesHandler(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID format"})
		return
	}

	responses, err := h.responseService.GetSurveyResponses(c.Request.Context(), surveyID)
	if err != nil {
		// TODO: Differentiate errors, e.g., if survey not found (though this service might not know)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve responses"})
		return
	}

	if responses == nil {
		// Return empty array if no responses, consistent with GetSurveys in survey-service
		c.JSON(http.StatusOK, []*models.Response{})
		return
	}

	c.JSON(http.StatusOK, responses)
}
