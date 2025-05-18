package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/survey-app/survey-service/internal/models"
)

// PostgresRepository implements the Repository interface
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Transaction management
func (r *PostgresRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *PostgresRepository) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (r *PostgresRepository) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

// CreateSurvey creates a new survey in the database
func (r *PostgresRepository) CreateSurvey(ctx context.Context, survey *models.Survey) (int, error) {
	query := `
		INSERT INTO surveys (creator_id, title, description, is_active, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query,
		survey.CreatorID,
		survey.Title,
		survey.Description,
		survey.IsActive,
		survey.StartDate,
		survey.EndDate,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	survey.ID = id
	survey.CreatedAt = createdAt
	survey.UpdatedAt = updatedAt

	return id, nil
}

// CreateSurveyTx creates a new survey in the database using a transaction
func (r *PostgresRepository) CreateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) (int, error) {
	query := `
		INSERT INTO surveys (creator_id, title, description, is_active, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time

	err := tx.QueryRow(ctx, query,
		survey.CreatorID,
		survey.Title,
		survey.Description,
		survey.IsActive,
		survey.StartDate,
		survey.EndDate,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	survey.ID = id
	survey.CreatedAt = createdAt
	survey.UpdatedAt = updatedAt

	return id, nil
}

// GetSurvey retrieves a survey by its ID
func (r *PostgresRepository) GetSurvey(ctx context.Context, id int) (*models.Survey, error) {
	query := `
		SELECT id, creator_id, title, description, is_active, start_date, end_date, created_at, updated_at
		FROM surveys
		WHERE id = $1
	`

	var survey models.Survey
	var startDate, endDate *time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&survey.ID,
		&survey.CreatorID,
		&survey.Title,
		&survey.Description,
		&survey.IsActive,
		&startDate,
		&endDate,
		&survey.CreatedAt,
		&survey.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}

	// Set optional dates if provided
	if startDate != nil {
		survey.StartDate = *startDate
	}
	if endDate != nil {
		survey.EndDate = *endDate
	}

	// Get questions
	questions, err := r.GetQuestionsBySurveyID(ctx, id)
	if err != nil {
		return nil, err
	}
	survey.Questions = questions

	return &survey, nil
}

// GetSurveys retrieves all surveys for a creator
func (r *PostgresRepository) GetSurveys(ctx context.Context, creatorID int) ([]*models.Survey, error) {
	query := `
		SELECT id, creator_id, title, description, is_active, start_date, end_date, created_at, updated_at
		FROM surveys
		WHERE creator_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []*models.Survey

	for rows.Next() {
		var survey models.Survey
		var startDate, endDate *time.Time

		err := rows.Scan(
			&survey.ID,
			&survey.CreatorID,
			&survey.Title,
			&survey.Description,
			&survey.IsActive,
			&startDate,
			&endDate,
			&survey.CreatedAt,
			&survey.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Set optional dates if provided
		if startDate != nil {
			survey.StartDate = *startDate
		}
		if endDate != nil {
			survey.EndDate = *endDate
		}

		surveys = append(surveys, &survey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return surveys, nil
}

// UpdateSurvey updates a survey in the database
func (r *PostgresRepository) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	query := `
		UPDATE surveys
		SET title = $1, description = $2, is_active = $3, start_date = $4, end_date = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query,
		survey.Title,
		survey.Description,
		survey.IsActive,
		survey.StartDate,
		survey.EndDate,
		survey.ID,
	).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("survey not found")
		}
		return err
	}

	survey.UpdatedAt = updatedAt
	return nil
}

// UpdateSurveyTx updates a survey in the database using a transaction
func (r *PostgresRepository) UpdateSurveyTx(ctx context.Context, tx pgx.Tx, survey *models.Survey) error {
	query := `
		UPDATE surveys
		SET title = $1, description = $2, is_active = $3, start_date = $4, end_date = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`
	var updatedAt time.Time
	err := tx.QueryRow(ctx, query,
		survey.Title,
		survey.Description,
		survey.IsActive,
		survey.StartDate,
		survey.EndDate,
		survey.ID,
	).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("survey not found during tx update")
		}
		return err
	}
	survey.UpdatedAt = updatedAt
	return nil
}

// DeleteSurvey deletes a survey from the database
func (r *PostgresRepository) DeleteSurvey(ctx context.Context, id int) error {
	// This will cascade delete questions and question options due to foreign key constraints
	query := `DELETE FROM surveys WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("survey not found")
	}

	return nil
}

// CreateQuestion creates a new question in the database
func (r *PostgresRepository) CreateQuestion(ctx context.Context, question *models.Question) (int, error) {
	query := `
		INSERT INTO questions (survey_id, text, type, required, order_num)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query,
		question.SurveyID,
		question.Text,
		question.Type,
		question.Required,
		question.OrderNum,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	question.ID = id
	question.CreatedAt = createdAt
	question.UpdatedAt = updatedAt

	return id, nil
}

// CreateQuestionTx creates a new question in the database using a transaction
func (r *PostgresRepository) CreateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) (int, error) {
	query := `
		INSERT INTO questions (survey_id, text, type, required, order_num)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	var id int
	var createdAt, updatedAt time.Time
	err := tx.QueryRow(ctx, query,
		question.SurveyID,
		question.Text,
		question.Type,
		question.Required,
		question.OrderNum,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return 0, err
	}
	question.ID = id
	question.CreatedAt = createdAt
	question.UpdatedAt = updatedAt
	return id, nil
}

// GetQuestionsBySurveyID retrieves all questions for a survey
func (r *PostgresRepository) GetQuestionsBySurveyID(ctx context.Context, surveyID int) ([]*models.Question, error) {
	query := `
		SELECT id, survey_id, text, type, required, order_num, created_at, updated_at
		FROM questions
		WHERE survey_id = $1
		ORDER BY order_num
	`

	rows, err := r.db.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question

	for rows.Next() {
		var question models.Question

		err := rows.Scan(
			&question.ID,
			&question.SurveyID,
			&question.Text,
			&question.Type,
			&question.Required,
			&question.OrderNum,
			&question.CreatedAt,
			&question.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Fetch options for question types that support them
		switch question.Type {
		case "multiple_choice", "checkbox", "dropdown":
			options, err := r.GetQuestionOptionsByQuestionID(ctx, question.ID)
			if err != nil {
				return nil, err // Or log and continue if one question's options failing shouldn't fail all
			}
			question.Options = options
		default:
			question.Options = []*models.QuestionOption{} // Ensure options is an empty slice, not nil, for other types
		}

		questions = append(questions, &question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// GetQuestionsBySurveyIDTx retrieves all questions for a survey using a transaction
func (r *PostgresRepository) GetQuestionsBySurveyIDTx(ctx context.Context, tx pgx.Tx, surveyID int) ([]*models.Question, error) {
	query := `
		SELECT id, survey_id, text, type, required, order_num, created_at, updated_at
		FROM questions
		WHERE survey_id = $1
		ORDER BY order_num
	`
	rows, err := tx.Query(ctx, query, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*models.Question
	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID, &question.SurveyID, &question.Text, &question.Type,
			&question.Required, &question.OrderNum, &question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		switch question.Type {
		case "multiple_choice", "checkbox", "dropdown":
			// Use GetQuestionOptionsByQuestionIDTx for transactional consistency
			options, err := r.GetQuestionOptionsByQuestionIDTx(ctx, tx, question.ID)
			if err != nil {
				return nil, err
			}
			question.Options = options
		default:
			question.Options = []*models.QuestionOption{}
		}
		questions = append(questions, &question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

// UpdateQuestion updates a question in the database
func (r *PostgresRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	query := `
		UPDATE questions
		SET text = $1, type = $2, required = $3, order_num = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query,
		question.Text,
		question.Type,
		question.Required,
		question.OrderNum,
		question.ID,
	).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("question not found")
		}
		return err
	}

	question.UpdatedAt = updatedAt
	return nil
}

// UpdateQuestionTx updates a question in the database using a transaction
func (r *PostgresRepository) UpdateQuestionTx(ctx context.Context, tx pgx.Tx, question *models.Question) error {
	query := `
		UPDATE questions
		SET text = $1, type = $2, required = $3, order_num = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	var updatedAt time.Time
	err := tx.QueryRow(ctx, query,
		question.Text, question.Type, question.Required, question.OrderNum, question.ID,
	).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("question not found during tx update")
		}
		return err
	}
	question.UpdatedAt = updatedAt
	return nil
}

// DeleteQuestion deletes a question from the database
func (r *PostgresRepository) DeleteQuestion(ctx context.Context, id int) error {
	// This will cascade delete question options due to foreign key constraints
	query := `DELETE FROM questions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("question not found")
	}

	return nil
}

// DeleteQuestionTx deletes a question from the database using a transaction
func (r *PostgresRepository) DeleteQuestionTx(ctx context.Context, tx pgx.Tx, id int) error {
	query := `DELETE FROM questions WHERE id = $1`
	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("question not found during tx delete")
	}
	return nil
}

// CreateQuestionOption creates a new question option in the database
func (r *PostgresRepository) CreateQuestionOption(ctx context.Context, option *models.QuestionOption) (int, error) {
	query := `
		INSERT INTO question_options (question_id, text, order_num)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query,
		option.QuestionID,
		option.Text,
		option.OrderNum,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	option.ID = id
	option.CreatedAt = createdAt
	option.UpdatedAt = updatedAt

	return id, nil
}

// CreateQuestionOptionTx creates a new question option using a transaction
func (r *PostgresRepository) CreateQuestionOptionTx(ctx context.Context, tx pgx.Tx, option *models.QuestionOption) (int, error) {
	query := `
		INSERT INTO question_options (question_id, text, order_num)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	var id int
	var createdAt, updatedAt time.Time
	err := tx.QueryRow(ctx, query,
		option.QuestionID, option.Text, option.OrderNum,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return 0, err
	}
	option.ID = id
	option.CreatedAt = createdAt
	option.UpdatedAt = updatedAt
	return id, nil
}

// GetQuestionOptionsByQuestionID retrieves all options for a question
func (r *PostgresRepository) GetQuestionOptionsByQuestionID(ctx context.Context, questionID int) ([]*models.QuestionOption, error) {
	query := `
		SELECT id, question_id, text, order_num, created_at, updated_at
		FROM question_options
		WHERE question_id = $1
		ORDER BY order_num
	`

	rows, err := r.db.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*models.QuestionOption

	for rows.Next() {
		var option models.QuestionOption

		err := rows.Scan(
			&option.ID,
			&option.QuestionID,
			&option.Text,
			&option.OrderNum,
			&option.CreatedAt,
			&option.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		options = append(options, &option)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

// GetQuestionOptionsByQuestionIDTx retrieves all options for a question using a transaction
func (r *PostgresRepository) GetQuestionOptionsByQuestionIDTx(ctx context.Context, tx pgx.Tx, questionID int) ([]*models.QuestionOption, error) {
	query := `
		SELECT id, question_id, text, order_num, created_at, updated_at
		FROM question_options
		WHERE question_id = $1
		ORDER BY order_num
	`
	rows, err := tx.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var options []*models.QuestionOption
	for rows.Next() {
		var option models.QuestionOption
		err := rows.Scan(
			&option.ID, &option.QuestionID, &option.Text, &option.OrderNum, &option.CreatedAt, &option.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		options = append(options, &option)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return options, nil
}

// DeleteQuestionOptions deletes all options for a question
func (r *PostgresRepository) DeleteQuestionOptions(ctx context.Context, questionID int) error {
	query := `DELETE FROM question_options WHERE question_id = $1`

	_, err := r.db.Exec(ctx, query, questionID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteQuestionOptionsTx deletes all options for a question using a transaction
func (r *PostgresRepository) DeleteQuestionOptionsTx(ctx context.Context, tx pgx.Tx, questionID int) error {
	query := `DELETE FROM question_options WHERE question_id = $1`
	_, err := tx.Exec(ctx, query, questionID)
	// We don't check RowsAffected here, as it's okay if a question had no options to delete
	return err
}
