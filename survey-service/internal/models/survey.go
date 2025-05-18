package models

import (
	"time"
)

// Survey represents a survey in the system
type Survey struct {
	ID          int         `json:"id"`
	CreatorID   int         `json:"creator_id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	IsActive    bool        `json:"is_active"`
	StartDate   time.Time   `json:"start_date,omitempty"`
	EndDate     time.Time   `json:"end_date,omitempty"`
	Questions   []*Question `json:"questions,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Question represents a question in a survey
type Question struct {
	ID        int               `json:"id"`
	SurveyID  int               `json:"survey_id"`
	Text      string            `json:"text"`
	Type      string            `json:"type"` // 'text', 'single_choice', etc.
	Required  bool              `json:"required"`
	OrderNum  int               `json:"order_num"`
	Options   []*QuestionOption `json:"options,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// QuestionOption represents an option for a question
type QuestionOption struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	Text       string    `json:"text"`
	OrderNum   int       `json:"order_num"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateSurveyRequest represents the data needed to create a new survey
type CreateSurveyRequest struct {
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description"`
	IsActive    bool                    `json:"is_active"`
	StartDate   time.Time               `json:"start_date"`
	EndDate     time.Time               `json:"end_date"`
	Questions   []QuestionUpdateRequest `json:"questions,omitempty"`
}

// UpdateSurveyRequest represents the data needed to update a survey
type UpdateSurveyRequest struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	IsActive    bool                    `json:"is_active"`
	StartDate   time.Time               `json:"start_date"`
	EndDate     time.Time               `json:"end_date"`
	Questions   []QuestionUpdateRequest `json:"questions"`
}

// QuestionUpdateRequest represents data for updating/creating a question within a survey update
type QuestionUpdateRequest struct {
	ID       *int     `json:"id,omitempty"`
	Text     string   `json:"text" binding:"required"`
	Type     string   `json:"type" binding:"required"`
	Required bool     `json:"required"`
	OrderNum int      `json:"order_num"`
	Options  []string `json:"options"`
}

// CreateQuestionRequest represents the data needed to create a new question
type CreateQuestionRequest struct {
	SurveyID int                           `json:"survey_id" binding:"required"`
	Text     string                        `json:"text" binding:"required"`
	Type     string                        `json:"type" binding:"required,oneof=text single_choice"`
	Required bool                          `json:"required"`
	OrderNum int                           `json:"order_num"`
	Options  []CreateQuestionOptionRequest `json:"options"`
}

// CreateQuestionOptionRequest represents the data needed to create a new question option
type CreateQuestionOptionRequest struct {
	Text     string `json:"text" binding:"required"`
	OrderNum int    `json:"order_num"`
}
