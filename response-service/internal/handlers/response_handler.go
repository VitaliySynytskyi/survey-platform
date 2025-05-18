package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// ExportSurveyResponsesCSV handles GET /api/v1/surveys/:surveyId/responses/export
func (h *ResponseHandler) ExportSurveyResponsesCSV(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID format"})
		return
	}

	responses, err := h.responseService.GetSurveyResponses(c.Request.Context(), surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve responses for export"})
		return
	}

	if len(responses) == 0 {
		// Return an empty CSV with headers if no responses
		b := &bytes.Buffer{}
		writer := csv.NewWriter(b)
		header := []string{"ResponseID", "UserID", "SubmittedAt", "QuestionID", "AnswerValue"}
		if err := writer.Write(header); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
			return
		}
		writer.Flush()
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=survey_%s_responses_%s.csv", surveyIDStr, time.Now().Format("20060102150405")))
		c.Header("Content-Type", "text/csv")
		c.Data(http.StatusOK, "text/csv", b.Bytes())
		return
	}

	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)

	header := []string{"ResponseID", "UserID", "SubmittedAt", "QuestionID", "AnswerValue"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	for _, resp := range responses {
		userIDStr := ""
		if resp.UserID != nil {
			userIDStr = strconv.Itoa(*resp.UserID)
		}
		for _, ans := range resp.Answers {
			var answerValueStr string
			switch v := ans.Value.(type) {
			case string:
				answerValueStr = v
			case []interface{}: // For checkboxes
				var vals []string
				for _, item := range v {
					if s, ok := item.(string); ok {
						vals = append(vals, s)
					} else {
						// Handle cases where items in []interface{} might not be strings
						vals = append(vals, fmt.Sprintf("%v", item))
					}
				}
				answerValueStr = strings.Join(vals, ",")
			case []string: // Should ideally not happen if JSON unmarshals to []interface{}
				answerValueStr = strings.Join(v, ",")
			default:
				answerValueStr = fmt.Sprintf("%v", v)
			}

			row := []string{
				resp.ID.Hex(),
				userIDStr,
				resp.SubmittedAt.Format(time.RFC3339Nano), // Using RFC3339Nano for more precision
				strconv.Itoa(ans.QuestionID),
				answerValueStr,
			}
			if err := writer.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
				return
			}
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error flushing CSV writer"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=survey_%s_responses_%s.csv", surveyIDStr, time.Now().Format("20060102150405")))
	c.Header("Content-Type", "text/csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
