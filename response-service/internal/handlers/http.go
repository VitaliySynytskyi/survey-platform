package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/contextkeys"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/models"
	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/service"
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
		log.Printf("[HANDLER_ERROR] SubmitResponse: Failed to bind JSON: %v. Request Body: %s", err, c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}
	log.Printf("[HANDLER_INFO] SubmitResponse: Received request for SurveyID: %d, UserID from body: %v", req.SurveyID, req.UserID)

	ctx := c.Request.Context()

	xUserIDStr := c.GetHeader("X-User-ID")
	if xUserIDStr != "" {
		uid, convErr := strconv.Atoi(xUserIDStr)
		if convErr == nil {
			log.Printf("[HANDLER_INFO] SubmitResponse: X-User-ID header found: '%s'. Adding to context.", xUserIDStr)
			ctx = context.WithValue(ctx, contextkeys.UserIDKey, uid)
			req.UserID = &uid
		} else {
			log.Printf("[HANDLER_WARN] SubmitResponse: X-User-ID header '%s' is not a valid integer: %v", xUserIDStr, convErr)
		}
	} else {
		log.Printf("[HANDLER_INFO] SubmitResponse: No X-User-ID header found. UserID from request body (if any) will be used: %v", req.UserID)
		if req.UserID != nil {
			ctx = context.WithValue(ctx, contextkeys.UserIDKey, *req.UserID)
		}
	}

	xUserRolesStr := c.GetHeader("X-User-Roles")
	if xUserRolesStr != "" {
		log.Printf("[HANDLER_INFO] SubmitResponse: X-User-Roles header found: '%s'. Parsing and adding to context.", xUserRolesStr)
		roles := parseRolesHeader(xUserRolesStr)
		if len(roles) > 0 {
			ctx = context.WithValue(ctx, contextkeys.UserRolesKey, roles)
			log.Printf("[HANDLER_INFO] SubmitResponse: Parsed roles: %v", roles)
		} else {
			log.Printf("[HANDLER_WARN] SubmitResponse: X-User-Roles header '%s' parsed to empty list.", xUserRolesStr)
		}
	} else {
		log.Printf("[HANDLER_INFO] SubmitResponse: No X-User-Roles header found.")
	}

	if err := h.responseService.SubmitResponse(ctx, &req); err != nil {
		log.Printf("[HANDLER_ERROR] SubmitResponse: Service call failed: %v", err)
		if strings.Contains(err.Error(), "not active") || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to submit response: %s", err.Error())})
		}
		return
	}

	log.Printf("[HANDLER_INFO] SubmitResponse: Successfully submitted response for SurveyID: %d", req.SurveyID)
	c.JSON(http.StatusCreated, gin.H{"message": "Response submitted successfully"})
}

// Helper function to parse roles from header like "[role1 role2]" or "role1,role2"
func parseRolesHeader(headerValue string) []string {
	cleaned := strings.Trim(headerValue, "[]")
	if strings.Contains(cleaned, " ") { // Likely space-separated: "[role1 role2]"
		return strings.Fields(cleaned)
	}
	if strings.Contains(cleaned, ",") { // Likely comma-separated: "role1,role2"
		return strings.Split(cleaned, ",")
	}
	if cleaned != "" { // Single role
		return []string{cleaned}
	}
	return []string{}
}

// GetSurveyResponsesHandler handles GET requests to /surveys/:surveyId/responses
func (h *ResponseHandler) GetSurveyResponsesHandler(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey_id format"})
		return
	}

	ctx := c.Request.Context()
	xUserIDStr := c.GetHeader("X-User-ID")
	if xUserIDStr != "" {
		if uid, convErr := strconv.Atoi(xUserIDStr); convErr == nil {
			ctx = context.WithValue(ctx, contextkeys.UserIDKey, uid)
		}
	}
	xUserRolesStr := c.GetHeader("X-User-Roles")
	if xUserRolesStr != "" {
		if roles := parseRolesHeader(xUserRolesStr); len(roles) > 0 {
			ctx = context.WithValue(ctx, contextkeys.UserRolesKey, roles)
		}
	}

	responses, err := h.responseService.GetSurveyResponses(ctx, surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get survey responses: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// ExportSurveyResponsesCSV handles GET requests to /surveys/:surveyId/responses/export
func (h *ResponseHandler) ExportSurveyResponsesCSV(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		log.Printf("[HANDLER_ERROR] ExportSurveyResponsesCSV: Invalid survey_id format '%s': %v", surveyIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey_id format"})
		return
	}

	log.Printf("[HANDLER_INFO] ExportSurveyResponsesCSV: Received request for SurveyID: %d", surveyID)
	ctx := c.Request.Context() // Original context

	// Forward Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		log.Printf("[HANDLER_INFO] ExportSurveyResponsesCSV: Forwarding Authorization header to service context.")
		ctx = context.WithValue(ctx, contextkeys.AuthorizationHeaderKey, authHeader)
	} else {
		log.Printf("[HANDLER_WARN] ExportSurveyResponsesCSV: No Authorization header found in incoming request.")
	}

	// Forward X-User-ID if present
	xUserIDStr := c.GetHeader("X-User-ID")
	if xUserIDStr != "" {
		if uid, convErr := strconv.Atoi(xUserIDStr); convErr == nil {
			log.Printf("[HANDLER_INFO] ExportSurveyResponsesCSV: Forwarding X-User-ID: %s to service context.", xUserIDStr)
			ctx = context.WithValue(ctx, contextkeys.UserIDKey, uid)
		} else {
			log.Printf("[HANDLER_WARN] ExportSurveyResponsesCSV: X-User-ID header '%s' is not a valid integer: %v", xUserIDStr, convErr)
		}
	}

	// Forward X-User-Roles if present
	xUserRolesStr := c.GetHeader("X-User-Roles")
	if xUserRolesStr != "" {
		if roles := parseRolesHeader(xUserRolesStr); len(roles) > 0 {
			log.Printf("[HANDLER_INFO] ExportSurveyResponsesCSV: Forwarding X-User-Roles: %v to service context.", roles)
			ctx = context.WithValue(ctx, contextkeys.UserRolesKey, roles)
		}
	}

	csvData, filename, err := h.responseService.ExportSurveyResponsesCSV(ctx, surveyID)
	if err != nil {
		log.Printf("[HANDLER_ERROR] ExportSurveyResponsesCSV: Service call failed for SurveyID %d: %v", surveyID, err)
		// Check for specific error types if needed, e.g., survey not found
		if strings.Contains(err.Error(), "not found") { // Basic check, could be more robust
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export responses: %s", err.Error())})
		}
		return
	}

	log.Printf("[HANDLER_INFO] ExportSurveyResponsesCSV: Successfully generated CSV for SurveyID %d. Filename: %s. Size: %d bytes", surveyID, filename, len(csvData))

	// Set headers for CSV download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(csvData))
}

// GetSurveyAnalytics handles GET requests to /surveys/{survey_id}/analytics
func (h *ResponseHandler) GetSurveyAnalytics(c *gin.Context) {
	surveyIDStr := c.Param("surveyId")
	surveyID, err := strconv.Atoi(surveyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Survey ID format in URL path"})
		return
	}

	ctx := c.Request.Context()
	xUserIDStr := c.GetHeader("X-User-ID")
	if xUserIDStr != "" {
		if uid, convErr := strconv.Atoi(xUserIDStr); convErr == nil {
			ctx = context.WithValue(ctx, contextkeys.UserIDKey, uid)
		}
	}
	xUserRolesStr := c.GetHeader("X-User-Roles")
	if xUserRolesStr != "" {
		if roles := parseRolesHeader(xUserRolesStr); len(roles) > 0 {
			ctx = context.WithValue(ctx, contextkeys.UserRolesKey, roles)
		}
	}

	analytics, err := h.responseService.GetSurveyAnalytics(ctx, surveyID)
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
