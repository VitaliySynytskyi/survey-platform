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
	log.Printf("SURVEY_TAKING_HANDLER: GetPublicSurvey called with surveyID: %s. Request URL: %s", surveyID, r.URL.String())

	// Validate surveyId
	if surveyID == "" {
		log.Printf("SURVEY_TAKING_HANDLER: GetPublicSurvey error: Missing survey ID")
		http.Error(w, "Missing survey ID", http.StatusBadRequest)
		return
	}

	// Get survey from the survey service or database
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	log.Printf("SURVEY_TAKING_HANDLER: Calling surveyClient.GetSurvey for surveyID: %s", surveyID)
	survey, err := h.surveyClient.GetSurvey(ctx, surveyID)
	if err != nil {
		log.Printf("SURVEY_TAKING_HANDLER: GetPublicSurvey error: Failed to get survey for ID %s from surveyClient. Error: %v", surveyID, err)
		http.Error(w, "Failed to get survey", http.StatusInternalServerError)
		return
	}

	log.Printf("SURVEY_TAKING_HANDLER: GetPublicSurvey successfully retrieved survey for ID %s. Title: %s", surveyID, survey.Title)
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
	log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse called with surveyID: %s. Request URL: %s", surveyID, r.URL.String())

	// Validate surveyId
	if surveyID == "" {
		log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse error: Invalid survey ID (empty)")
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var response model.SurveyResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse error for surveyID %s: Invalid request body: %v", surveyID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate survey response
	if len(response.Answers) == 0 {
		log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse error for surveyID %s: No answers provided", surveyID)
		http.Error(w, "No answers provided", http.StatusBadRequest)
		return
	}

	// Set survey ID and submission time
	response.SurveyID = surveyID
	response.SubmittedAt = time.Now().UTC().Format(time.RFC3339)

	// If no authenticated user, ensure anonymous ID is set
	if response.RespondentID == "" && response.AnonymousID == "" {
		log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse error for surveyID %s: RespondentID and AnonymousID are both empty", surveyID)
		http.Error(w, "Either respondent ID or anonymous ID must be provided", http.StatusBadRequest)
		return
	}

	log.Printf("SURVEY_TAKING_HANDLER: Survey response for surveyID %s is valid. Attempting to publish to RabbitMQ. AnonymousID: %s, RespondentID: %s, Answers: %d", surveyID, response.AnonymousID, response.RespondentID, len(response.Answers))

	// Publish to RabbitMQ
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.producer.PublishResponse(ctx, response)
	if err != nil {
		log.Printf("SURVEY_TAKING_HANDLER: SubmitSurveyResponse error for surveyID %s: Error publishing to RabbitMQ: %v", surveyID, err)
		http.Error(w, "Failed to process response", http.StatusInternalServerError)
		return
	}

	log.Printf("SURVEY_TAKING_HANDLER: Successfully published response for surveyID %s to RabbitMQ", surveyID)
	// Respond with success
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Survey response accepted for processing",
	})
}

// RegisterRoutes registers API routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/{surveyId}/public", h.GetPublicSurvey).Methods("GET")
	router.HandleFunc("/{surveyId}/responses", h.SubmitSurveyResponse).Methods("POST")
}
