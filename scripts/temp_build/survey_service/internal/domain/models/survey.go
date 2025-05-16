package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuestionType визначає типи питань в опитуванні
type QuestionType string

const (
	SingleChoice   QuestionType = "single-choice"
	MultipleChoice QuestionType = "multiple-choice"
	OpenText       QuestionType = "open-text"
	Scale          QuestionType = "scale"
	MatrixSingle   QuestionType = "matrix-single"
	MatrixMultiple QuestionType = "matrix-multiple"
)

// Option представляє варіант відповіді для питань з вибором
type Option struct {
	Value string `json:"value" bson:"value"`
	Text  string `json:"text" bson:"text"`
}

// ScaleSettings представляє налаштування для питань типу scale
type ScaleSettings struct {
	Min      int    `json:"min" bson:"min"`
	Max      int    `json:"max" bson:"max"`
	MinLabel string `json:"min_label" bson:"min_label"`
	MaxLabel string `json:"max_label" bson:"max_label"`
}

// DisplayLogic представляє логіку умовного відображення питання
type DisplayLogic struct {
	DependsOnQuestionID  string `json:"depends_on_question_id" bson:"depends_on_question_id"`
	DependsOnAnswerValue string `json:"depends_on_answer_value" bson:"depends_on_answer_value"`
}

// Question представляє структуру питання в опитуванні
type Question struct {
	ID            string         `json:"question_id" bson:"question_id"`
	Text          string         `json:"text" bson:"text"`
	Type          QuestionType   `json:"type" bson:"type"`
	IsRequired    bool           `json:"is_required" bson:"is_required"`
	Options       []Option       `json:"options,omitempty" bson:"options,omitempty"`
	ScaleSettings *ScaleSettings `json:"scale_settings,omitempty" bson:"scale_settings,omitempty"`
	MatrixRows    []string       `json:"matrix_rows,omitempty" bson:"matrix_rows,omitempty"`
	MatrixColumns []string       `json:"matrix_columns,omitempty" bson:"matrix_columns,omitempty"`
	DisplayLogic  *DisplayLogic  `json:"display_logic,omitempty" bson:"display_logic,omitempty"`
}

// Survey представляє структуру опитування
type Survey struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	OwnerID     string             `json:"owner_id" bson:"owner_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Questions   []Question         `json:"questions" bson:"questions"`
}

// CreateSurveyRequest представляє запит на створення опитування
type CreateSurveyRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions" binding:"required"`
}

// UpdateSurveyRequest представляє запит на оновлення опитування
type UpdateSurveyRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

// SurveyResponse представляє відповідь з даними опитування
type SurveyResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	OwnerID     string     `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Questions   []Question `json:"questions"`
}

// SurveyListResponse представляє відповідь зі списком опитувань
type SurveyListResponse struct {
	Surveys    []SurveyResponse `json:"surveys"`
	TotalCount int64            `json:"total_count"`
	Page       int64            `json:"page"`
	PerPage    int64            `json:"per_page"`
}

// ToResponse конвертує модель Survey в SurveyResponse
func (s *Survey) ToResponse() SurveyResponse {
	return SurveyResponse{
		ID:          s.ID.Hex(),
		Title:       s.Title,
		Description: s.Description,
		OwnerID:     s.OwnerID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Questions:   s.Questions,
	}
}
