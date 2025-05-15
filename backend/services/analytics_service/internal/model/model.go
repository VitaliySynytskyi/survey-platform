package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SurveyResponse represents a response to a survey
type SurveyResponse struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SurveyID     primitive.ObjectID `bson:"survey_id" json:"survey_id"`
	RespondentID string             `bson:"respondent_id,omitempty" json:"respondent_id,omitempty"`
	AnonymousID  string             `bson:"anonymous_id,omitempty" json:"anonymous_id,omitempty"`
	Answers      []Answer           `bson:"answers" json:"answers"`
	SubmittedAt  time.Time          `bson:"submitted_at" json:"submitted_at"`
}

// Answer represents an answer to a survey question
type Answer struct {
	QuestionID string      `bson:"question_id" json:"question_id"`
	Value      interface{} `bson:"value" json:"value"` // string, []string, number depending on question type
}

// SurveyQuestion represents a question in a survey
type SurveyQuestion struct {
	ID          string         `bson:"_id" json:"id"`
	Title       string         `bson:"title" json:"title"`
	Type        string         `bson:"type" json:"type"` // single-choice, multiple-choice, open-text, scale, matrix
	Required    bool           `bson:"required" json:"required"`
	Order       int            `bson:"order" json:"order"`
	Options     []SurveyOption `bson:"options,omitempty" json:"options,omitempty"`
	Description string         `bson:"description,omitempty" json:"description,omitempty"`
}

// SurveyOption represents an option for a question
type SurveyOption struct {
	ID    string `bson:"_id" json:"id"`
	Text  string `bson:"text" json:"text"`
	Order int    `bson:"order" json:"order"`
}

// Survey represents a survey with its questions
type Survey struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	OwnerID     string             `bson:"owner_id" json:"owner_id"`
	Questions   []SurveyQuestion   `bson:"questions" json:"questions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// SurveyResults represents aggregated results for a survey
type SurveyResults struct {
	SurveyID         string           `json:"survey_id"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	TotalResponses   int              `json:"total_responses"`
	QuestionResults  []QuestionResult `json:"question_results"`
	CompletionRate   float64          `json:"completion_rate,omitempty"`
	AverageTimeSpent string           `json:"average_time_spent,omitempty"`
}

// QuestionResult represents aggregated results for a question
type QuestionResult struct {
	QuestionID   string            `json:"question_id"`
	Title        string            `json:"title"`
	Type         string            `json:"type"`
	ResponseData interface{}       `json:"response_data"` // Type depends on question type
	TotalAnswers int               `json:"total_answers"`
	Analytics    QuestionAnalytics `json:"analytics,omitempty"`
}

// QuestionAnalytics contains analytics for specific question types
type QuestionAnalytics struct {
	// For single-choice and multiple-choice
	OptionCounts map[string]OptionCount `json:"option_counts,omitempty"`

	// For scale questions
	Average *float64 `json:"average,omitempty"`
	Median  *float64 `json:"median,omitempty"`
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`

	// For open-text (optional)
	WordFrequency map[string]int `json:"word_frequency,omitempty"`
}

// OptionCount represents count for an option
type OptionCount struct {
	OptionID string  `json:"option_id"`
	Text     string  `json:"text"`
	Count    int     `json:"count"`
	Percent  float64 `json:"percent"`
}

// IndividualResponsesFilter represents filter parameters for fetching individual responses
type IndividualResponsesFilter struct {
	SurveyID  string     `json:"survey_id"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// IndividualResponsesResult represents paginated individual responses
type IndividualResponsesResult struct {
	SurveyID     string           `json:"survey_id"`
	TotalCount   int              `json:"total_count"`
	Responses    []SurveyResponse `json:"responses"`
	CurrentPage  int              `json:"current_page"`
	TotalPages   int              `json:"total_pages"`
	ItemsPerPage int              `json:"items_per_page"`
}
