package repository

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/response-service/internal/models"
)

// ResponseRepositoryInterface defines methods for interacting with response data
type ResponseRepositoryInterface interface {
	CreateResponse(ctx context.Context, response *models.Response) error
	GetResponsesBySurveyID(ctx context.Context, surveyID int) ([]*models.Response, error)
	// Add other methods as needed, e.g., GetResponseByID, GetResponsesByUserID, etc.
}
