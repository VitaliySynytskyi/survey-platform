package mock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/survey-service/internal/models"
)

func TestNewMockRepository(t *testing.T) {
	repo := NewMockRepository()

	if repo.Surveys == nil {
		t.Error("Expected Surveys map to be initialized")
	}

	if repo.Questions == nil {
		t.Error("Expected Questions map to be initialized")
	}

	if repo.Options == nil {
		t.Error("Expected Options map to be initialized")
	}
}

func TestMockRepositorySurveyOperations(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	// Create a survey
	now := time.Now()
	survey := &models.Survey{
		CreatorID:   1,
		Title:       "Test Survey",
		Description: "Test Description",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	id, err := repo.CreateSurvey(ctx, survey)
	if err != nil {
		t.Fatalf("Unexpected error creating survey: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected survey ID to be 1, got %d", id)
	}

	// Get survey
	fetchedSurvey, err := repo.GetSurvey(ctx, id)
	if err != nil {
		t.Fatalf("Unexpected error getting survey: %v", err)
	}
	if fetchedSurvey.Title != "Test Survey" {
		t.Errorf("Expected survey title 'Test Survey', got '%s'", fetchedSurvey.Title)
	}

	// Update survey
	survey.Title = "Updated Survey"
	err = repo.UpdateSurvey(ctx, survey)
	if err != nil {
		t.Fatalf("Unexpected error updating survey: %v", err)
	}

	fetchedSurvey, _ = repo.GetSurvey(ctx, id)
	if fetchedSurvey.Title != "Updated Survey" {
		t.Errorf("Expected updated survey title 'Updated Survey', got '%s'", fetchedSurvey.Title)
	}

	// List surveys
	surveys, count, err := repo.ListSurveysByCreatorID(ctx, 1, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error listing surveys: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
	if len(surveys) != 1 {
		t.Errorf("Expected 1 survey, got %d", len(surveys))
	}

	// Update status
	err = repo.UpdateSurveyStatus(ctx, id, false)
	if err != nil {
		t.Fatalf("Unexpected error updating status: %v", err)
	}

	fetchedSurvey, _ = repo.GetSurvey(ctx, id)
	if fetchedSurvey.IsActive {
		t.Error("Expected survey to be inactive")
	}

	// Test error handling
	testErr := errors.New("test error")
	repo.ErrorMock = testErr

	_, err = repo.CreateSurvey(ctx, survey)
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}

	// Delete survey
	repo.ErrorMock = nil
	err = repo.DeleteSurvey(ctx, id)
	if err != nil {
		t.Fatalf("Unexpected error deleting survey: %v", err)
	}

	fetchedSurvey, _ = repo.GetSurvey(ctx, id)
	if fetchedSurvey != nil {
		t.Error("Expected survey to be deleted")
	}
}

func TestMockRepositoryQuestionOperations(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	// Create a survey first
	survey := &models.Survey{
		CreatorID: 1,
		Title:     "Test Survey",
		IsActive:  true,
	}
	surveyID, _ := repo.CreateSurvey(ctx, survey)

	// Create a question
	now := time.Now()
	question := &models.Question{
		SurveyID:  surveyID,
		Text:      "Test Question",
		Type:      "text",
		Required:  true,
		OrderNum:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := repo.CreateQuestion(ctx, question)
	if err != nil {
		t.Fatalf("Unexpected error creating question: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected question ID to be 1, got %d", id)
	}

	// Get question
	fetchedQuestion, err := repo.GetQuestionByID(ctx, id)
	if err != nil {
		t.Fatalf("Unexpected error getting question: %v", err)
	}
	if fetchedQuestion.Text != "Test Question" {
		t.Errorf("Expected question text 'Test Question', got '%s'", fetchedQuestion.Text)
	}

	// Update question
	question.Text = "Updated Question"
	err = repo.UpdateQuestion(ctx, question)
	if err != nil {
		t.Fatalf("Unexpected error updating question: %v", err)
	}

	fetchedQuestion, _ = repo.GetQuestionByID(ctx, id)
	if fetchedQuestion.Text != "Updated Question" {
		t.Errorf("Expected updated question text 'Updated Question', got '%s'", fetchedQuestion.Text)
	}

	// Get questions by survey ID
	questions, err := repo.GetQuestionsBySurveyID(ctx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting questions by survey ID: %v", err)
	}
	if len(questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(questions))
	}

	// Delete question
	err = repo.DeleteQuestion(ctx, id)
	if err != nil {
		t.Fatalf("Unexpected error deleting question: %v", err)
	}

	fetchedQuestion, _ = repo.GetQuestionByID(ctx, id)
	if fetchedQuestion != nil {
		t.Error("Expected question to be deleted")
	}
}

func TestMockRepositoryQuestionOptionOperations(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	// Create a survey and question first
	survey := &models.Survey{
		CreatorID: 1,
		Title:     "Test Survey",
		IsActive:  true,
	}
	surveyID, _ := repo.CreateSurvey(ctx, survey)

	question := &models.Question{
		SurveyID: surveyID,
		Text:     "Test Question",
		Type:     "single_choice",
		Required: true,
		OrderNum: 1,
	}
	questionID, _ := repo.CreateQuestion(ctx, question)

	// Create a question option
	now := time.Now()
	option := &models.QuestionOption{
		QuestionID: questionID,
		Text:       "Option 1",
		OrderNum:   1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	id, err := repo.CreateQuestionOption(ctx, option)
	if err != nil {
		t.Fatalf("Unexpected error creating option: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected option ID to be 1, got %d", id)
	}

	// Get options by question ID
	options, err := repo.GetQuestionOptionsByQuestionID(ctx, questionID)
	if err != nil {
		t.Fatalf("Unexpected error getting options: %v", err)
	}
	if len(options) != 1 {
		t.Errorf("Expected 1 option, got %d", len(options))
	}
	if options[0].Text != "Option 1" {
		t.Errorf("Expected option text 'Option 1', got '%s'", options[0].Text)
	}

	// Delete options
	err = repo.DeleteQuestionOptions(ctx, questionID)
	if err != nil {
		t.Fatalf("Unexpected error deleting options: %v", err)
	}

	options, _ = repo.GetQuestionOptionsByQuestionID(ctx, questionID)
	if len(options) != 0 {
		t.Errorf("Expected options to be deleted, got %d", len(options))
	}
}

func TestMockRepositoryTransactions(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	// Test transaction operations
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatalf("Unexpected error beginning transaction: %v", err)
	}

	// Create survey in transaction
	now := time.Now()
	survey := &models.Survey{
		CreatorID:   1,
		Title:       "Transaction Survey",
		Description: "Created in transaction",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
	}

	id, err := repo.CreateSurveyTx(ctx, tx, survey)
	if err != nil {
		t.Fatalf("Unexpected error creating survey in transaction: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected survey ID to be 1, got %d", id)
	}

	// Create question in transaction
	question := &models.Question{
		SurveyID: id,
		Text:     "Transaction Question",
		Type:     "text",
		Required: true,
		OrderNum: 1,
	}

	qid, err := repo.CreateQuestionTx(ctx, tx, question)
	if err != nil {
		t.Fatalf("Unexpected error creating question in transaction: %v", err)
	}

	// Create option in transaction
	option := &models.QuestionOption{
		QuestionID: qid,
		Text:       "Transaction Option",
		OrderNum:   1,
	}

	_, err = repo.CreateQuestionOptionTx(ctx, tx, option)
	if err != nil {
		t.Fatalf("Unexpected error creating option in transaction: %v", err)
	}

	// Commit transaction
	err = repo.CommitTx(ctx, tx)
	if err != nil {
		t.Fatalf("Unexpected error committing transaction: %v", err)
	}

	// Verify data was saved
	fetchedSurvey, _ := repo.GetSurvey(ctx, id)
	if fetchedSurvey.Title != "Transaction Survey" {
		t.Errorf("Expected survey title 'Transaction Survey', got '%s'", fetchedSurvey.Title)
	}

	// Test rollback with error
	testErr := errors.New("test error")
	repo.ErrorMock = testErr

	_, err = repo.BeginTx(ctx)
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestDeleteQuestionsBySurveyIDTx(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	// Create a survey
	survey := &models.Survey{
		CreatorID: 1,
		Title:     "Test Survey",
		IsActive:  true,
	}
	surveyID, _ := repo.CreateSurvey(ctx, survey)

	// Create questions for the survey
	question1 := &models.Question{
		SurveyID: surveyID,
		Text:     "Question 1",
		Type:     "text",
	}
	repo.CreateQuestion(ctx, question1)

	question2 := &models.Question{
		SurveyID: surveyID,
		Text:     "Question 2",
		Type:     "text",
	}
	repo.CreateQuestion(ctx, question2)

	// Create questions for another survey
	otherSurvey := &models.Survey{
		CreatorID: 1,
		Title:     "Other Survey",
		IsActive:  true,
	}
	otherSurveyID, _ := repo.CreateSurvey(ctx, otherSurvey)

	otherQuestion := &models.Question{
		SurveyID: otherSurveyID,
		Text:     "Other Question",
		Type:     "text",
	}
	repo.CreateQuestion(ctx, otherQuestion)

	// Delete questions for the first survey
	tx, _ := repo.BeginTx(ctx)
	err := repo.DeleteQuestionsBySurveyIDTx(ctx, tx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error deleting questions: %v", err)
	}
	repo.CommitTx(ctx, tx)

	// Verify questions are deleted for the first survey
	questions, _ := repo.GetQuestionsBySurveyID(ctx, surveyID)
	if len(questions) != 0 {
		t.Errorf("Expected 0 questions for survey %d, got %d", surveyID, len(questions))
	}

	// Verify questions for other survey are not affected
	otherQuestions, _ := repo.GetQuestionsBySurveyID(ctx, otherSurveyID)
	if len(otherQuestions) != 1 {
		t.Errorf("Expected 1 question for survey %d, got %d", otherSurveyID, len(otherQuestions))
	}
}
