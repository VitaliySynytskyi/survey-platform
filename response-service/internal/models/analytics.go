package models

// SurveyAnalyticsResponse represents the overall analytics data for a survey
type SurveyAnalyticsResponse struct {
	SurveyID          int                 `json:"survey_id"`
	SurveyTitle       string              `json:"survey_title"`
	TotalResponses    int                 `json:"total_responses"`
	QuestionAnalytics []QuestionAnalytics `json:"question_analytics"`
}

// QuestionAnalytics represents analytics for a single question
type QuestionAnalytics struct {
	QuestionID     int                `json:"question_id"`
	QuestionText   string             `json:"question_text"`
	QuestionType   string             `json:"question_type"`
	OptionsSummary []OptionSummary    `json:"options_summary,omitempty"`
	TextResponses  []TextResponseData `json:"text_responses,omitempty"`
}

// OptionSummary provides a summary for a single answer option
type OptionSummary struct {
	OptionID   *int    `json:"option_id,omitempty"`
	OptionText string  `json:"option_text"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TextResponseData holds a single text response
type TextResponseData struct {
	Response string `json:"response"`
}
