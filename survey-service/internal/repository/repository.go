package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/VitaliySynytskyi/survey-platform/survey-service/internal/models"
)

// SurveyRepositoryInterface defines the interface for survey and question database operations
type SurveyRepositoryInterface interface {
	// Transaction management
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) error
	RollbackTx(ctx context.Context, tx pgx.Tx) error

	// Survey operations (non-transactional)
	CreateSurvey(ctx context.Context, survey *models.Survey) (int, error)
	GetSurvey(ctx context.Context, id int) (*models.Survey, error)
	ListSurveysByCreatorID(ctx context.Context, creatorID int, offset, limit int) ([]*models.Survey, int, error)
	ListAllSurveys(ctx context.Context, isUserAdmin bool, offset, limit int) ([]*models.Survey, int, error)
	UpdateSurvey(ctx context.Context, survey *models.Survey) error
	DeleteSurvey(ctx context.Context, id int) error
	UpdateSurveyStatus(ctx context.Context, id int, isActive bool) error

	// Question operations (non-transactional)
	CreateQuestion(ctx context.Context, question *models.Question) (int, error)
	GetQuestionByID(ctx context.Context, id int) (*models.Question, error)
	GetQuestionsBySurveyID(ctx context.Context, surveyID int) ([]*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, id int) error

	// QuestionOption operations (non-transactional)
	CreateQuestionOption(ctx context.Context, option *models.QuestionOption) (int, error)
	GetQuestionOptionsByQuestionID(ctx context.Context, questionID int) ([]*models.QuestionOption, error)
	DeleteQuestionOptions(ctx context.Context, questionID int) error

	// Survey operations (transactional)
	CreateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) (int, error)
	UpdateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) error

	// Question operations (transactional)
	CreateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) (int, error)
	GetQuestionsBySurveyIDTx(ctx context.Context, tx pgx.Tx, surveyID int) ([]*models.Question, error)
	UpdateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) error
	DeleteQuestionTx(ctx context.Context, tx pgx.Tx, id int) error

	// QuestionOption operations (transactional)
	CreateQuestionOptionTx(ctx context.Context, tx pgx.Tx, option *models.QuestionOption) (int, error)
	// GetQuestionOptionsByQuestionIDTx is not strictly needed if GetQuestionsBySurveyIDTx handles it, but can be added for consistency
	DeleteQuestionOptionsTx(ctx context.Context, tx pgx.Tx, questionID int) error
	DeleteQuestionsBySurveyIDTx(ctx context.Context, tx pgx.Tx, surveyID int) error
}
