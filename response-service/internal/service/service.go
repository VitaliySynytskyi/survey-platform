package service

import (
	"context"

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
	repo repository.ResponseRepositoryInterface
}

// NewResponseService creates a new ResponseService
func NewResponseService(repo repository.ResponseRepositoryInterface) *ResponseService {
	return &ResponseService{repo: repo}
}

// SubmitResponse handles the business logic for submitting a new survey response
func (s *ResponseService) SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error {
	// TODO: Add any validation logic here if needed
	// For example, check if the survey ID is valid by calling survey-service (requires inter-service communication setup)
	// Check if question IDs in answers are valid for the survey, etc.

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
