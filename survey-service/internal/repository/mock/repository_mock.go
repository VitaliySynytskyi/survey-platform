package mock

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/survey-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MockTx is a mock transaction that implements pgx.Tx interface
type MockTx struct{}

// Begin implements pgx.Tx
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

// BeginFunc implements pgx.Tx
func (m *MockTx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	return f(m)
}

// Commit implements pgx.Tx
func (m *MockTx) Commit(ctx context.Context) error {
	return nil
}

// Rollback implements pgx.Tx
func (m *MockTx) Rollback(ctx context.Context) error {
	return nil
}

// CopyFrom implements pgx.Tx
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

// SendBatch implements pgx.Tx
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

// LargeObjects implements pgx.Tx
func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

// Prepare implements pgx.Tx
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

// Exec implements pgx.Tx
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

// Query implements pgx.Tx
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

// QueryRow implements pgx.Tx
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return nil
}

// Conn implements pgx.Tx
func (m *MockTx) Conn() *pgx.Conn {
	return nil
}

// MockRepository is a mock implementation of the SurveyRepositoryInterface
type MockRepository struct {
	// Mock behavior flags and return values
	Surveys     map[int]*models.Survey
	Questions   map[int]*models.Question
	Options     map[int]*models.QuestionOption
	SurveyCount int
	ErrorMock   error
}

// NewMockRepository creates a new instance of the mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Surveys:   make(map[int]*models.Survey),
		Questions: make(map[int]*models.Question),
		Options:   make(map[int]*models.QuestionOption),
	}
}

// BeginTx mocks beginning a transaction
func (m *MockRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	if m.ErrorMock != nil {
		return nil, m.ErrorMock
	}
	return &MockTx{}, nil
}

// CommitTx mocks committing a transaction
func (m *MockRepository) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return m.ErrorMock
}

// RollbackTx mocks rolling back a transaction
func (m *MockRepository) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return m.ErrorMock
}

// CreateSurvey mocks creating a survey
func (m *MockRepository) CreateSurvey(ctx context.Context, survey *models.Survey) (int, error) {
	if m.ErrorMock != nil {
		return 0, m.ErrorMock
	}

	survey.ID = len(m.Surveys) + 1
	m.Surveys[survey.ID] = survey
	return survey.ID, nil
}

// GetSurvey mocks retrieving a survey
func (m *MockRepository) GetSurvey(ctx context.Context, id int) (*models.Survey, error) {
	if m.ErrorMock != nil {
		return nil, m.ErrorMock
	}

	survey, exists := m.Surveys[id]
	if !exists {
		return nil, nil
	}
	return survey, nil
}

// ListSurveysByCreatorID mocks listing surveys by creator ID
func (m *MockRepository) ListSurveysByCreatorID(ctx context.Context, creatorID int, offset, limit int) ([]*models.Survey, int, error) {
	if m.ErrorMock != nil {
		return nil, 0, m.ErrorMock
	}

	var surveys []*models.Survey
	for _, survey := range m.Surveys {
		if survey.CreatorID == creatorID {
			surveys = append(surveys, survey)
		}
	}

	// Apply pagination
	totalCount := len(surveys)
	if offset < len(surveys) {
		end := offset + limit
		if end > len(surveys) {
			end = len(surveys)
		}
		surveys = surveys[offset:end]
	} else {
		surveys = []*models.Survey{}
	}

	return surveys, totalCount, nil
}

// ListAllSurveys mocks listing all surveys
func (m *MockRepository) ListAllSurveys(ctx context.Context, isUserAdmin bool, offset, limit int) ([]*models.Survey, int, error) {
	if m.ErrorMock != nil {
		return nil, 0, m.ErrorMock
	}

	var surveys []*models.Survey
	for _, survey := range m.Surveys {
		// Include all surveys for admins, but only active surveys for non-admin users
		if isUserAdmin || survey.IsActive {
			surveys = append(surveys, survey)
		}
	}

	// Apply pagination
	totalCount := len(surveys)
	if offset < len(surveys) {
		end := offset + limit
		if end > len(surveys) {
			end = len(surveys)
		}
		surveys = surveys[offset:end]
	} else {
		surveys = []*models.Survey{}
	}

	return surveys, totalCount, nil
}

// UpdateSurvey mocks updating a survey
func (m *MockRepository) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	if _, exists := m.Surveys[survey.ID]; !exists {
		return nil
	}

	m.Surveys[survey.ID] = survey
	return nil
}

// DeleteSurvey mocks deleting a survey
func (m *MockRepository) DeleteSurvey(ctx context.Context, id int) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	delete(m.Surveys, id)
	return nil
}

// UpdateSurveyStatus mocks updating a survey's status
func (m *MockRepository) UpdateSurveyStatus(ctx context.Context, id int, isActive bool) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	survey, exists := m.Surveys[id]
	if !exists {
		return nil
	}

	survey.IsActive = isActive
	return nil
}

// CreateQuestion mocks creating a question
func (m *MockRepository) CreateQuestion(ctx context.Context, question *models.Question) (int, error) {
	if m.ErrorMock != nil {
		return 0, m.ErrorMock
	}

	question.ID = len(m.Questions) + 1
	m.Questions[question.ID] = question
	return question.ID, nil
}

// GetQuestionByID mocks retrieving a question by ID
func (m *MockRepository) GetQuestionByID(ctx context.Context, id int) (*models.Question, error) {
	if m.ErrorMock != nil {
		return nil, m.ErrorMock
	}

	question, exists := m.Questions[id]
	if !exists {
		return nil, nil
	}
	return question, nil
}

// GetQuestionsBySurveyID mocks retrieving questions by survey ID
func (m *MockRepository) GetQuestionsBySurveyID(ctx context.Context, surveyID int) ([]*models.Question, error) {
	if m.ErrorMock != nil {
		return nil, m.ErrorMock
	}

	var questions []*models.Question
	for _, question := range m.Questions {
		if question.SurveyID == surveyID {
			questions = append(questions, question)
		}
	}

	return questions, nil
}

// UpdateQuestion mocks updating a question
func (m *MockRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	if _, exists := m.Questions[question.ID]; !exists {
		return nil
	}

	m.Questions[question.ID] = question
	return nil
}

// DeleteQuestion mocks deleting a question
func (m *MockRepository) DeleteQuestion(ctx context.Context, id int) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	delete(m.Questions, id)
	return nil
}

// CreateQuestionOption mocks creating a question option
func (m *MockRepository) CreateQuestionOption(ctx context.Context, option *models.QuestionOption) (int, error) {
	if m.ErrorMock != nil {
		return 0, m.ErrorMock
	}

	option.ID = len(m.Options) + 1
	m.Options[option.ID] = option
	return option.ID, nil
}

// GetQuestionOptionsByQuestionID mocks retrieving options by question ID
func (m *MockRepository) GetQuestionOptionsByQuestionID(ctx context.Context, questionID int) ([]*models.QuestionOption, error) {
	if m.ErrorMock != nil {
		return nil, m.ErrorMock
	}

	var options []*models.QuestionOption
	for _, option := range m.Options {
		if option.QuestionID == questionID {
			options = append(options, option)
		}
	}

	return options, nil
}

// DeleteQuestionOptions mocks deleting options for a question
func (m *MockRepository) DeleteQuestionOptions(ctx context.Context, questionID int) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	for id, option := range m.Options {
		if option.QuestionID == questionID {
			delete(m.Options, id)
		}
	}

	return nil
}

// Transactional operations - these are simplified in the mock to call their non-transactional counterparts

// CreateSurveyTx mocks creating a survey in a transaction
func (m *MockRepository) CreateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) (int, error) {
	return m.CreateSurvey(ctx, survey)
}

// UpdateSurveyTx mocks updating a survey in a transaction
func (m *MockRepository) UpdateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) error {
	return m.UpdateSurvey(ctx, survey)
}

// CreateQuestionTx mocks creating a question in a transaction
func (m *MockRepository) CreateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) (int, error) {
	return m.CreateQuestion(ctx, question)
}

// GetQuestionsBySurveyIDTx mocks retrieving questions by survey ID in a transaction
func (m *MockRepository) GetQuestionsBySurveyIDTx(ctx context.Context, tx pgx.Tx, surveyID int) ([]*models.Question, error) {
	return m.GetQuestionsBySurveyID(ctx, surveyID)
}

// UpdateQuestionTx mocks updating a question in a transaction
func (m *MockRepository) UpdateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) error {
	return m.UpdateQuestion(ctx, question)
}

// DeleteQuestionTx mocks deleting a question in a transaction
func (m *MockRepository) DeleteQuestionTx(ctx context.Context, tx pgx.Tx, id int) error {
	return m.DeleteQuestion(ctx, id)
}

// CreateQuestionOptionTx mocks creating a question option in a transaction
func (m *MockRepository) CreateQuestionOptionTx(ctx context.Context, tx pgx.Tx, option *models.QuestionOption) (int, error) {
	return m.CreateQuestionOption(ctx, option)
}

// DeleteQuestionOptionsTx mocks deleting options for a question in a transaction
func (m *MockRepository) DeleteQuestionOptionsTx(ctx context.Context, tx pgx.Tx, questionID int) error {
	return m.DeleteQuestionOptions(ctx, questionID)
}

// DeleteQuestionsBySurveyIDTx mocks deleting questions for a survey in a transaction
func (m *MockRepository) DeleteQuestionsBySurveyIDTx(ctx context.Context, tx pgx.Tx, surveyID int) error {
	if m.ErrorMock != nil {
		return m.ErrorMock
	}

	for id, question := range m.Questions {
		if question.SurveyID == surveyID {
			delete(m.Questions, id)
		}
	}

	return nil
}
