package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api/handlers/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/domain/models"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/store/mongodb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SurveyResponsesCollectionName = "survey_responses"
)

type StoredAnswer struct {
	QuestionID string      `bson:"question_id"`
	Value      interface{} `bson:"value"`
}

type StoredSurveyResponse struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SurveyID string             `bson:"survey_id"`
	Answers  []StoredAnswer     `bson:"answers"`
}

type SurveyHandler struct {
	repository    mongodb.Repository
	mongoDBClient *mongo.Client
	surveyDBName  string
}

func NewSurveyHandler(repo mongodb.Repository, client *mongo.Client, dbName string) *SurveyHandler {
	return &SurveyHandler{
		repository:    repo,
		mongoDBClient: client,
		surveyDBName:  dbName,
	}
}

func (h *SurveyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var createReq models.CreateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateSurveyRequest(&createReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	survey := models.Survey{
		ID:          primitive.NewObjectID(),
		Title:       createReq.Title,
		Description: createReq.Description,
		OwnerID:     userCtx.UserID,
		Questions:   createReq.Questions,
	}

	if err := h.repository.Create(r.Context(), &survey); err != nil {
		http.Error(w, "Failed to create survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(survey.ToResponse())
}

func (h *SurveyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	survey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	if survey.OwnerID != userCtx.UserID && userCtx.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey.ToResponse())
}

func (h *SurveyHandler) Update(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	surveyID := vars["surveyId"]
	objectID, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	existingSurvey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	if existingSurvey.OwnerID != userCtx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var updateReq models.UpdateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updateReq.Title != "" {
		existingSurvey.Title = updateReq.Title
	}

	if updateReq.Description != "" {
		existingSurvey.Description = updateReq.Description
	}

	if updateReq.Questions != nil {
		if err := validateQuestions(updateReq.Questions); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		existingSurvey.Questions = updateReq.Questions
	}

	existingSurvey.ID = objectID

	if err := h.repository.Update(r.Context(), existingSurvey); err != nil {
		http.Error(w, "Failed to update survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingSurvey.ToResponse())
}

func (h *SurveyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	existingSurvey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	if existingSurvey.OwnerID != userCtx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.repository.Delete(r.Context(), surveyID); err != nil {
		http.Error(w, "Failed to delete survey", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SurveyHandler) GetUserSurveys(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID != userCtx.UserID && userCtx.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	surveys, total, err := h.repository.GetByOwnerID(r.Context(), userID, page, perPage)
	if err != nil {
		http.Error(w, "Failed to get surveys", http.StatusInternalServerError)
		return
	}

	surveysResponse := models.SurveyListResponse{
		Surveys:    make([]models.SurveyResponse, 0, len(surveys)),
		TotalCount: total,
		Page:       page,
		PerPage:    perPage,
	}

	for _, survey := range surveys {
		surveysResponse.Surveys = append(surveysResponse.Surveys, survey.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveysResponse)
}

func (h *SurveyHandler) GetAllSurveys(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	surveys, total, err := h.repository.GetAll(r.Context(), page, perPage)
	if err != nil {
		http.Error(w, "Failed to get all surveys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	surveysResponse := models.SurveyListResponse{
		Surveys:    make([]models.SurveyResponse, 0, len(surveys)),
		TotalCount: total,
		Page:       page,
		PerPage:    perPage,
	}

	for _, survey := range surveys {
		surveysResponse.Surveys = append(surveysResponse.Surveys, survey.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveysResponse)
}

func (h *SurveyHandler) GetPublicSurveyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]
	log.Printf("SURVEY_HANDLER: GetPublicSurveyByID called with surveyID: %s", surveyID)

	if surveyID == "" {
		log.Printf("SURVEY_HANDLER: GetPublicSurveyByID error: Survey ID is required")
		http.Error(w, "Survey ID is required", http.StatusBadRequest)
		return
	}

	survey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNotFound) {
			log.Printf("SURVEY_HANDLER: GetPublicSurveyByID error: Survey not found for ID %s. Repository error: %v", surveyID, err)
			http.Error(w, "Survey not found", http.StatusNotFound)
		} else if errors.Is(err, mongodb.ErrInvalidID) {
			log.Printf("SURVEY_HANDLER: GetPublicSurveyByID error: Invalid survey ID format for ID %s. Repository error: %v", surveyID, err)
			http.Error(w, "Invalid survey ID format", http.StatusBadRequest)
		} else {
			log.Printf("SURVEY_HANDLER: GetPublicSurveyByID error: Failed to retrieve survey for ID %s. Repository error: %v", surveyID, err)
			http.Error(w, "Failed to retrieve survey", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("SURVEY_HANDLER: GetPublicSurveyByID successfully retrieved survey for ID %s", surveyID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey.ToResponse())
}

type SurveyResult struct {
	OptionID   string  `json:"option_id"`
	OptionText string  `json:"option_text"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type QuestionResult struct {
	QuestionID       string              `json:"question_id"`
	QuestionText     string              `json:"question_text"`
	Type             models.QuestionType `json:"type"`
	Results          []SurveyResult      `json:"results,omitempty"`
	OpenEndedAnswers []string            `json:"open_ended_answers,omitempty"`
}

type SurveyResultsResponse struct {
	SurveyID        string           `json:"survey_id"`
	SurveyTitle     string           `json:"survey_title"`
	TotalResponses  int              `json:"total_responses"`
	QuestionResults []QuestionResult `json:"question_results"`
}

func (h *SurveyHandler) GetSurveyResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surveyIDHex := vars["surveyId"]
	log.Printf("SURVEY_HANDLER: GetSurveyResults called with surveyID: %s", surveyIDHex)

	if surveyIDHex == "" {
		log.Printf("SURVEY_HANDLER: GetSurveyResults error: Survey ID is required")
		http.Error(w, "Survey ID is required", http.StatusBadRequest)
		return
	}

	survey, err := h.repository.GetByID(r.Context(), surveyIDHex)
	if err != nil {
		if errors.Is(err, mongodb.ErrNotFound) {
			log.Printf("SURVEY_HANDLER: GetSurveyResults error: Survey not found for ID %s. Repo error: %v", surveyIDHex, err)
			http.Error(w, "Survey not found, cannot retrieve results", http.StatusNotFound)
		} else if errors.Is(err, mongodb.ErrInvalidID) {
			log.Printf("SURVEY_HANDLER: GetSurveyResults error: Invalid survey ID format %s. Repo error: %v", surveyIDHex, err)
			http.Error(w, "Invalid survey ID format for results", http.StatusBadRequest)
		} else {
			log.Printf("SURVEY_HANDLER: GetSurveyResults error: Failed to get survey %s. Repo error: %v", surveyIDHex, err)
			http.Error(w, "Failed to retrieve survey before fetching results", http.StatusInternalServerError)
		}
		return
	}
	log.Printf("SURVEY_HANDLER: GetSurveyResults: Successfully fetched survey structure for ID %s, Title: %s", surveyIDHex, survey.Title)

	responsesCollection := h.mongoDBClient.Database(h.surveyDBName).Collection(SurveyResponsesCollectionName)
	filter := bson.M{"survey_id": surveyIDHex}

	cursor, err := responsesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("SURVEY_HANDLER: GetSurveyResults error: Failed to query survey_responses for surveyID %s: %v", surveyIDHex, err)
		http.Error(w, "Failed to fetch survey responses", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var storedResponses []StoredSurveyResponse
	if err = cursor.All(context.Background(), &storedResponses); err != nil {
		log.Printf("SURVEY_HANDLER: GetSurveyResults error: Failed to decode survey_responses for surveyID %s: %v", surveyIDHex, err)
		http.Error(w, "Failed to decode survey responses", http.StatusInternalServerError)
		return
	}
	log.Printf("SURVEY_HANDLER: GetSurveyResults: Fetched %d responses for surveyID %s", len(storedResponses), surveyIDHex)

	aggregatedResults := SurveyResultsResponse{
		SurveyID:        surveyIDHex,
		SurveyTitle:     survey.Title,
		TotalResponses:  len(storedResponses),
		QuestionResults: make([]QuestionResult, 0, len(survey.Questions)),
	}

	for _, qDefinition := range survey.Questions {
		qr := QuestionResult{
			QuestionID:   qDefinition.ID,
			QuestionText: qDefinition.Text,
			Type:         qDefinition.Type,
		}

		switch qDefinition.Type {
		case models.SingleChoice, models.MultipleChoice:
			optionCounts := make(map[string]int)
			for _, option := range qDefinition.Options {
				optionCounts[option.Value] = 0
			}

			for _, resp := range storedResponses {
				for _, ans := range resp.Answers {
					if ans.QuestionID == qDefinition.ID {
						if qDefinition.Type == models.SingleChoice {
							if optValStr, ok := ans.Value.(string); ok {
								optionCounts[optValStr]++
							}
						} else if qDefinition.Type == models.MultipleChoice {
							if optVals, ok := ans.Value.(primitive.A); ok {
								for _, optValInterface := range optVals {
									if optValStr, okStr := optValInterface.(string); okStr {
										optionCounts[optValStr]++
									}
								}
							} else if optValsStr, ok := ans.Value.([]string); ok {
								for _, optValStr := range optValsStr {
									optionCounts[optValStr]++
								}
							}
						}
					}
				}
			}

			qr.Results = make([]SurveyResult, 0, len(qDefinition.Options))
			for _, option := range qDefinition.Options {
				count := optionCounts[option.Value]
				percentage := 0.0
				if aggregatedResults.TotalResponses > 0 {
					percentage = (float64(count) / float64(aggregatedResults.TotalResponses)) * 100
				}
				qr.Results = append(qr.Results, SurveyResult{
					OptionID:   option.Value,
					OptionText: option.Text,
					Count:      count,
					Percentage: percentage,
				})
			}

		case models.OpenText:
			qr.OpenEndedAnswers = make([]string, 0)
			for _, resp := range storedResponses {
				for _, ans := range resp.Answers {
					if ans.QuestionID == qDefinition.ID {
						if textAnswer, ok := ans.Value.(string); ok {
							qr.OpenEndedAnswers = append(qr.OpenEndedAnswers, textAnswer)
						}
					}
				}
			}
		default:
			log.Printf("SURVEY_HANDLER: GetSurveyResults: Skipping aggregation for unsupported question type %s (ID: %s)", qDefinition.Type, qDefinition.ID)
		}
		aggregatedResults.QuestionResults = append(aggregatedResults.QuestionResults, qr)
	}

	log.Printf("SURVEY_HANDLER: GetSurveyResults successfully aggregated results for survey ID %s", surveyIDHex)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aggregatedResults)
}

func validateSurveyRequest(req *models.CreateSurveyRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}

	if len(req.Questions) == 0 {
		return errors.New("at least one question is required")
	}

	return validateQuestions(req.Questions)
}

func validateQuestions(questions []models.Question) error {
	for i, q := range questions {
		if q.Text == "" {
			return errors.New("question text is required")
		}

		switch q.Type {
		case models.SingleChoice, models.MultipleChoice:
			if len(q.Options) < 2 {
				return errors.New("choice questions must have at least 2 options")
			}
		case models.Scale:
			if q.ScaleSettings == nil {
				return errors.New("scale questions must have scale settings")
			}
			if q.ScaleSettings.Min >= q.ScaleSettings.Max {
				return errors.New("scale min must be less than max")
			}
		case models.MatrixSingle, models.MatrixMultiple:
			if len(q.MatrixRows) == 0 {
				return errors.New("matrix questions must have at least one row")
			}
			if len(q.MatrixColumns) == 0 {
				return errors.New("matrix questions must have at least one column")
			}
		case models.OpenText:
		default:
			return errors.New("invalid question type at index " + strconv.Itoa(i))
		}

		if q.DisplayLogic != nil && q.DisplayLogic.DependsOnQuestionID != "" {
			found := false
			for _, prevQ := range questions {
				if prevQ.ID == q.DisplayLogic.DependsOnQuestionID {
					found = true
					break
				}
			}
			if !found {
				return errors.New("display logic references non-existent question")
			}
		}
	}

	return nil
}
