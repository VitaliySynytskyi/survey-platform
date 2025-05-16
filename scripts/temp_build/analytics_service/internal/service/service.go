package service

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/model"
)

// AnalyticsService handles business logic for analytics
type AnalyticsService struct {
	repo *db.AnalyticsRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo *db.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// GetSurveyResults retrieves the aggregated results for a survey
func (s *AnalyticsService) GetSurveyResults(ctx context.Context, surveyID string, userID string, isAdmin bool) (*model.SurveyResults, error) {
	// First get the survey to check ownership
	survey, err := s.repo.GetSurveyById(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	// Check if the user has access to this survey
	if !isAdmin && survey.OwnerID != userID {
		return nil, model.NewAccessDeniedError("you don't have access to this survey")
	}

	// Get survey results
	return s.repo.GetSurveyResults(ctx, surveyID)
}

// GetIndividualResponses retrieves individual responses for a survey
func (s *AnalyticsService) GetIndividualResponses(ctx context.Context, filter model.IndividualResponsesFilter, userID string, isAdmin bool) (*model.IndividualResponsesResult, error) {
	// First get the survey to check ownership
	survey, err := s.repo.GetSurveyById(ctx, filter.SurveyID)
	if err != nil {
		return nil, err
	}

	// Check if the user has access to this survey
	if !isAdmin && survey.OwnerID != userID {
		return nil, model.NewAccessDeniedError("you don't have access to this survey")
	}

	// Get individual responses
	return s.repo.GetIndividualResponses(ctx, filter)
}

// Close closes the service
func (s *AnalyticsService) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
