package models

// SurveyDetailsFromService represents the survey structure fetched from survey-service
type SurveyDetailsFromService struct {
	ID          int                   `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description,omitempty"`
	IsActive    bool                  `json:"is_active"`
	Questions   []QuestionFromService `json:"questions,omitempty"`
}

// QuestionFromService represents a question structure fetched from survey-service
type QuestionFromService struct {
	ID      int                         `json:"id"`
	Text    string                      `json:"text"`
	Type    string                      `json:"type"` // 'text', 'single_choice', etc.
	Options []QuestionOptionFromService `json:"options,omitempty"`
}

// QuestionOptionFromService represents a question option fetched from survey-service
type QuestionOptionFromService struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
