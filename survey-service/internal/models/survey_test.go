package models

import (
	"testing"
	"time"
)

func TestSurvey(t *testing.T) {
	// Test creation of Survey struct
	now := time.Now()
	survey := Survey{
		ID:          1,
		CreatorID:   100,
		Title:       "Test Survey",
		Description: "This is a test survey",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if survey.ID != 1 {
		t.Errorf("Expected survey ID to be 1, got %d", survey.ID)
	}
	if survey.CreatorID != 100 {
		t.Errorf("Expected survey CreatorID to be 100, got %d", survey.CreatorID)
	}
	if survey.Title != "Test Survey" {
		t.Errorf("Expected survey Title to be 'Test Survey', got %s", survey.Title)
	}
	if !survey.IsActive {
		t.Error("Expected survey to be active")
	}
}

func TestQuestion(t *testing.T) {
	// Test creation of Question struct
	now := time.Now()
	question := Question{
		ID:        1,
		SurveyID:  1,
		Text:      "What is your name?",
		Type:      "text",
		Required:  true,
		OrderNum:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if question.ID != 1 {
		t.Errorf("Expected question ID to be 1, got %d", question.ID)
	}
	if question.Text != "What is your name?" {
		t.Errorf("Expected question Text to be 'What is your name?', got %s", question.Text)
	}
	if question.Type != "text" {
		t.Errorf("Expected question Type to be 'text', got %s", question.Type)
	}
	if !question.Required {
		t.Error("Expected question to be required")
	}
}

func TestQuestionOption(t *testing.T) {
	// Test creation of QuestionOption struct
	now := time.Now()
	option := QuestionOption{
		ID:         1,
		QuestionID: 1,
		Text:       "Option 1",
		OrderNum:   1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if option.ID != 1 {
		t.Errorf("Expected option ID to be 1, got %d", option.ID)
	}
	if option.QuestionID != 1 {
		t.Errorf("Expected option QuestionID to be 1, got %d", option.QuestionID)
	}
	if option.Text != "Option 1" {
		t.Errorf("Expected option Text to be 'Option 1', got %s", option.Text)
	}
}

func TestSurveyWithQuestions(t *testing.T) {
	// Test Survey with nested Questions and Options
	now := time.Now()

	option1 := &QuestionOption{
		ID:         1,
		QuestionID: 1,
		Text:       "Yes",
		OrderNum:   1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	option2 := &QuestionOption{
		ID:         2,
		QuestionID: 1,
		Text:       "No",
		OrderNum:   2,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	question := &Question{
		ID:        1,
		SurveyID:  1,
		Text:      "Do you agree?",
		Type:      "single_choice",
		Required:  true,
		OrderNum:  1,
		Options:   []*QuestionOption{option1, option2},
		CreatedAt: now,
		UpdatedAt: now,
	}

	survey := Survey{
		ID:          1,
		CreatorID:   100,
		Title:       "Survey with Questions",
		Description: "This is a test survey with questions",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
		Questions:   []*Question{question},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if len(survey.Questions) != 1 {
		t.Fatalf("Expected survey to have 1 question, got %d", len(survey.Questions))
	}

	if survey.Questions[0].ID != 1 {
		t.Errorf("Expected question ID to be 1, got %d", survey.Questions[0].ID)
	}

	if len(survey.Questions[0].Options) != 2 {
		t.Fatalf("Expected question to have 2 options, got %d", len(survey.Questions[0].Options))
	}

	if survey.Questions[0].Options[0].Text != "Yes" {
		t.Errorf("Expected first option text to be 'Yes', got %s", survey.Questions[0].Options[0].Text)
	}
}
