package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/repository"
)

// ResponseServiceInterface defines methods for response-related business logic
type ResponseServiceInterface interface {
	SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error
	GetSurveyResponses(ctx context.Context, surveyID int) ([]*models.Response, error)
}

// ResponseService implements ResponseServiceInterface
type ResponseService struct {
	repo             repository.ResponseRepositoryInterface
	surveyServiceURL string
	httpClient       *http.Client
}

// NewResponseService creates a new ResponseService
func NewResponseService(repo repository.ResponseRepositoryInterface, surveyServiceURL string) *ResponseService {
	return &ResponseService{
		repo:             repo,
		surveyServiceURL: surveyServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Add a timeout for HTTP requests
		},
	}
}

// SubmitResponse handles the business logic for submitting a new survey response
func (s *ResponseService) SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error {
	// Check if the survey is active by calling survey-service
	surveyURL := fmt.Sprintf("%s/api/v1/surveys/%d", s.surveyServiceURL, req.SurveyID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", surveyURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to survey-service: %w", err)
	}

	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call survey-service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return errors.New("survey not found")
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("survey-service returned status %d", httpResp.StatusCode)
	}

	var surveyDetails struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&surveyDetails); err != nil {
		return fmt.Errorf("failed to decode response from survey-service: %w", err)
	}

	if !surveyDetails.IsActive {
		return errors.New("survey is not active and cannot accept new responses")
	}

	// Check if question IDs in answers are valid for the survey, etc. - This part can be expanded later

	response := &models.Response{
		SurveyID: req.SurveyID,
		UserID:   req.UserID,
		Answers:  req.Answers,
		// SubmittedAt will be set by the repository
	}

	return s.repo.CreateResponse(ctx, response)
}

// GetSurveyResponses retrieves all responses for a specific survey
func (s *ResponseService) GetSurveyResponses(ctx context.Context, surveyID int) ([]*models.Response, error) {
	// TODO: Add any transformation or additional logic if needed
	return s.repo.GetResponsesBySurveyID(ctx, surveyID)
}
