package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/client"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/model"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/rabbitmq"
)

// Handler represents the API handler
type Handler struct {
	surveyClient client.SurveyClient
	producer     *rabbitmq.Producer
}

// NewHandler creates a new API handler
func NewHandler(surveyClient client.SurveyClient, producer *rabbitmq.Producer) *Handler {
	return &Handler{
		surveyClient: surveyClient,
		producer:     producer,
	}
}

// GetPublicSurvey handles GET /surveys/{surveyId}/public
func (h *Handler) GetPublicSurvey(w http.ResponseWriter, r *http.Request) {
	// Get surveyId from URL
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	// Validate surveyId
	if surveyID == "" {
		http.Error(w, "Missing survey ID", http.StatusBadRequest)
		return
	}

	// Get survey from the survey service or database
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	survey, err := h.surveyClient.GetSurvey(ctx, surveyID)
	if err != nil {
		log.Printf("Error getting survey: %v", err)
		http.Error(w, "Failed to get survey", http.StatusInternalServerError)
		return
	}

	// Respond with the public view of the survey
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(survey)
}

// SubmitSurveyResponse handles POST /surveys/{surveyId}/responses
func (h *Handler) SubmitSurveyResponse(w http.ResponseWriter, r *http.Request) {
	// Get surveyId from URL
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	// Validate surveyId
	if surveyID == "" {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var response model.SurveyResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate survey response
	if len(response.Answers) == 0 {
		http.Error(w, "No answers provided", http.StatusBadRequest)
		return
	}

	// Set survey ID and submission time
	response.SurveyID = surveyID
	response.SubmittedAt = time.Now().UTC().Format(time.RFC3339)

	// Get respondent info from auth token if available
	// (simplified here, actual implementation would check auth token)
	// response.RespondentID = extractUserIDFromToken(r)

	// If no authenticated user, ensure anonymous ID is set
	if response.RespondentID == "" && response.AnonymousID == "" {
		http.Error(w, "Either respondent ID or anonymous ID must be provided", http.StatusBadRequest)
		return
	}

	// Publish to RabbitMQ
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.producer.PublishResponse(ctx, response)
	if err != nil {
		log.Printf("Error publishing to RabbitMQ: %v", err)
		http.Error(w, "Failed to process response", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Survey response accepted for processing",
	})
}

// RegisterRoutes registers API routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/surveys/{surveyId}/public", h.GetPublicSurvey).Methods("GET")
	router.HandleFunc("/surveys/{surveyId}/responses", h.SubmitSurveyResponse).Methods("POST")
}
