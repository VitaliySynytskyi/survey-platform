package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/model"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.AnalyticsService
}

// NewHandler creates a new HTTP handler
func NewHandler(service *service.AnalyticsService) *Handler {
	return &Handler{
		service: service,
	}
}

// GetSurveyResults handles the request to get survey results
func (h *Handler) GetSurveyResults(w http.ResponseWriter, r *http.Request) {
	// Get survey ID from URL
	surveyID := chi.URLParam(r, "surveyId")
	if surveyID == "" {
		http.Error(w, "Survey ID is required", http.StatusBadRequest)
		return
	}

	// Get user ID and role from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin := auth.IsAdmin(r.Context())

	// Get survey results
	results, err := h.service.GetSurveyResults(r.Context(), surveyID, userID, isAdmin)
	if err != nil {
		handleError(w, err)
		return
	}

	// Return response
	respondJSON(w, http.StatusOK, results)
}

// GetIndividualResponses handles the request to get individual responses
func (h *Handler) GetIndividualResponses(w http.ResponseWriter, r *http.Request) {
	// Get survey ID from URL
	surveyID := chi.URLParam(r, "surveyId")
	if surveyID == "" {
		http.Error(w, "Survey ID is required", http.StatusBadRequest)
		return
	}

	// Get user ID and role from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin := auth.IsAdmin(r.Context())

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	// Parse date filters if provided
	var startDate, endDate *time.Time
	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		if sd, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &sd
		}
	}
	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		if ed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &ed
		}
	}

	// Create filter
	filter := model.IndividualResponsesFilter{
		SurveyID:  surveyID,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	// Get individual responses
	responses, err := h.service.GetIndividualResponses(r.Context(), filter, userID, isAdmin)
	if err != nil {
		handleError(w, err)
		return
	}

	// Return response
	respondJSON(w, http.StatusOK, responses)
}

// Helper functions

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// handleError handles errors
func handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*model.AppError); ok {
		respondJSON(w, appErr.Code, map[string]string{
			"error": appErr.Message,
			"type":  appErr.Type,
		})
		return
	}

	// Default to 500 error
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
