package model

// SurveyPublic represents the public view of a survey for respondents
type SurveyPublic struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Questions   []QuestionPublic `json:"questions"`
	CreatedAt   string           `json:"created_at,omitempty"`
	UpdatedAt   string           `json:"updated_at,omitempty"`
}

// QuestionPublic represents the public view of a question for respondents
type QuestionPublic struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Type        string         `json:"type"` // single-choice, multiple-choice, open-text, scale, matrix
	Required    bool           `json:"required"`
	Order       int            `json:"order"`
	Options     []OptionPublic `json:"options,omitempty"`
	Description string         `json:"description,omitempty"`
}

// OptionPublic represents the public view of an option for respondents
type OptionPublic struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Order int    `json:"order"`
}

// SurveyResponse represents a response submitted for a survey
type SurveyResponse struct {
	SurveyID     string   `json:"survey_id"`
	RespondentID string   `json:"respondent_id,omitempty"`
	AnonymousID  string   `json:"anonymous_id,omitempty"`
	Answers      []Answer `json:"answers"`
	SubmittedAt  string   `json:"submitted_at"`
}

// Answer represents an answer to a survey question
type Answer struct {
	QuestionID string      `json:"question_id"`
	Value      interface{} `json:"value"` // can be string, []string, number, etc. depending on question type
}

// RabbitMQMessage represents a message to be sent to RabbitMQ
type RabbitMQMessage struct {
	SurveyID     string   `json:"survey_id"`
	RespondentID string   `json:"respondent_id,omitempty"`
	AnonymousID  string   `json:"anonymous_id,omitempty"`
	Answers      []Answer `json:"answers"`
	SubmittedAt  string   `json:"submitted_at"`
}
