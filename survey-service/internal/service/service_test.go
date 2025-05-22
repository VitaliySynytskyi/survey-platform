package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/survey-service/internal/models"
	"github.com/VitaliySynytskyi/survey-platform/survey-service/internal/repository/mock"
)

func setupTestContext(userID int, roles []string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UserRolesKey, roles)
	return ctx
}

func TestCreateSurvey(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create a context with user info
	userID := 1
	roles := []string{"user"}
	ctx := setupTestContext(userID, roles)

	// Create a survey and questions
	now := time.Now()
	survey := &models.Survey{
		Title:       "Test Survey",
		Description: "Test Description",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
	}

	questions := []models.QuestionUpdateRequest{
		{
			Text:     "Question 1",
			Type:     "text",
			Required: true,
			OrderNum: 1,
		},
		{
			Text:     "Question 2",
			Type:     "single_choice",
			Required: false,
			OrderNum: 2,
			Options:  []string{"Option 1", "Option 2"},
		},
	}

	// Create the survey
	surveyID, err := service.CreateSurvey(ctx, survey, questions)
	if err != nil {
		t.Fatalf("Unexpected error creating survey: %v", err)
	}
	if surveyID != 1 {
		t.Errorf("Expected survey ID to be 1, got %d", surveyID)
	}

	// Verify survey was saved to the repository
	savedSurvey, err := mockRepo.GetSurvey(ctx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting survey: %v", err)
	}
	if savedSurvey.Title != "Test Survey" {
		t.Errorf("Expected survey title 'Test Survey', got '%s'", savedSurvey.Title)
	}
	if savedSurvey.CreatorID != userID {
		t.Errorf("Expected creator ID %d, got %d", userID, savedSurvey.CreatorID)
	}

	// Check for questions
	savedQuestions, err := mockRepo.GetQuestionsBySurveyID(ctx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting questions: %v", err)
	}
	if len(savedQuestions) != 2 {
		t.Fatalf("Expected 2 questions, got %d", len(savedQuestions))
	}

	// Verify second question has options
	questionID := savedQuestions[1].ID
	options, err := mockRepo.GetQuestionOptionsByQuestionID(ctx, questionID)
	if err != nil {
		t.Fatalf("Unexpected error getting options: %v", err)
	}
	if len(options) != 2 {
		t.Fatalf("Expected 2 options, got %d", len(options))
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	_, err = service.CreateSurvey(ctx, survey, questions)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestGetSurvey(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create a test survey first
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

	ctx := context.Background() // Use a simple context for initial setup
	surveyID, _ := mockRepo.CreateSurvey(ctx, survey)

	// Create questions and options
	question1 := &models.Question{
		SurveyID:  surveyID,
		Text:      "Question 1",
		Type:      "text",
		Required:  true,
		OrderNum:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// Not using q1ID to avoid linter error
	mockRepo.CreateQuestion(ctx, question1)

	question2 := &models.Question{
		SurveyID:  surveyID,
		Text:      "Question 2",
		Type:      "single_choice",
		Required:  false,
		OrderNum:  2,
		CreatedAt: now,
		UpdatedAt: now,
	}
	q2ID, _ := mockRepo.CreateQuestion(ctx, question2)

	option1 := &models.QuestionOption{
		QuestionID: q2ID,
		Text:       "Option 1",
		OrderNum:   1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	mockRepo.CreateQuestionOption(ctx, option1)

	option2 := &models.QuestionOption{
		QuestionID: q2ID,
		Text:       "Option 2",
		OrderNum:   2,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	mockRepo.CreateQuestionOption(ctx, option2)

	// Now create an authenticated context for testing access
	userCtx := setupTestContext(1, []string{"user"})

	// Get the survey as owner
	fetchedSurvey, err := service.GetSurvey(userCtx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting survey: %v", err)
	}
	if fetchedSurvey.Title != "Test Survey" {
		t.Errorf("Expected survey title 'Test Survey', got '%s'", fetchedSurvey.Title)
	}

	// Get survey as a different user - should be allowed because survey is active
	otherUserCtx := setupTestContext(2, []string{"user"})
	// Using publicSurvey to avoid linter error
	publicSurvey, err := service.GetSurvey(otherUserCtx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting public survey: %v", err)
	}
	// Verify publicSurvey is not nil
	if publicSurvey == nil {
		t.Error("Expected non-nil public survey")
	}

	// Set survey to inactive and try again as a different user
	survey.IsActive = false
	mockRepo.UpdateSurvey(ctx, survey)

	_, err = service.GetSurvey(otherUserCtx, surveyID)
	if err == nil {
		t.Error("Expected error for inactive survey accessed by non-owner, got nil")
	}

	// Try with admin access
	adminCtx := setupTestContext(3, []string{"admin"})
	adminSurvey, err := service.GetSurvey(adminCtx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error getting survey with admin: %v", err)
	}
	if adminSurvey == nil {
		t.Error("Admin should be able to access inactive survey but got nil")
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	_, err = service.GetSurvey(userCtx, surveyID)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestUpdateSurveyWithQuestions(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create a test survey first
	now := time.Now()
	survey := &models.Survey{
		CreatorID:   1,
		Title:       "Original Title",
		Description: "Original Description",
		IsActive:    true,
		StartDate:   now,
		EndDate:     now.Add(24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ctx := context.Background() // Simple context for setup
	surveyID, _ := mockRepo.CreateSurvey(ctx, survey)

	// Create original question
	question := &models.Question{
		SurveyID:  surveyID,
		Text:      "Original Question",
		Type:      "text",
		Required:  true,
		OrderNum:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	questionID, _ := mockRepo.CreateQuestion(ctx, question)

	// Create user context for updates
	userCtx := setupTestContext(1, []string{"user"})

	// Create survey update data
	updateSurvey := &models.Survey{
		ID:          surveyID,
		Title:       "Updated Title",
		Description: "Updated Description",
		IsActive:    false,
		StartDate:   now.Add(1 * time.Hour),
		EndDate:     now.Add(48 * time.Hour),
	}

	updateQuestions := []models.QuestionUpdateRequest{
		{
			ID:       &questionID, // Update existing question
			Text:     "Updated Question",
			Type:     "text",
			Required: false,
			OrderNum: 1,
		},
		{
			Text:     "New Question",
			Type:     "single_choice",
			Required: true,
			OrderNum: 2,
			Options:  []string{"New Option 1", "New Option 2"},
		},
	}

	// Update the survey
	err := service.UpdateSurveyWithQuestions(userCtx, updateSurvey, updateQuestions)
	if err != nil {
		t.Fatalf("Unexpected error updating survey: %v", err)
	}

	// Verify survey was updated - need to get fresh
	mockRepo.ErrorMock = nil
	updatedSurvey, _ := mockRepo.GetSurvey(ctx, surveyID)

	if updatedSurvey.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updatedSurvey.Title)
	}

	if updatedSurvey.IsActive {
		t.Error("Expected survey to be inactive")
	}

	// Test unauthorized update
	otherUserCtx := setupTestContext(999, []string{"user"})
	err = service.UpdateSurveyWithQuestions(otherUserCtx, updateSurvey, updateQuestions)
	if err == nil {
		t.Error("Expected error for unauthorized update but got nil")
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	err = service.UpdateSurveyWithQuestions(userCtx, updateSurvey, updateQuestions)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestDeleteSurvey(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create a test survey
	ctx := context.Background()
	survey := &models.Survey{
		CreatorID: 1,
		Title:     "Test Survey",
		IsActive:  true,
	}
	surveyID, _ := mockRepo.CreateSurvey(ctx, survey)

	// Create a question
	question := &models.Question{
		SurveyID: surveyID,
		Text:     "Question",
		Type:     "text",
	}
	questionID, _ := mockRepo.CreateQuestion(ctx, question)

	// Create options
	option := &models.QuestionOption{
		QuestionID: questionID,
		Text:       "Option",
	}
	mockRepo.CreateQuestionOption(ctx, option)

	// Create user context for delete
	userCtx := setupTestContext(1, []string{"user"})

	// Delete the survey
	err := service.DeleteSurvey(userCtx, surveyID)
	if err != nil {
		t.Fatalf("Unexpected error deleting survey: %v", err)
	}

	// Verify survey is deleted
	deletedSurvey, _ := mockRepo.GetSurvey(ctx, surveyID)
	if deletedSurvey != nil {
		t.Error("Expected survey to be deleted")
	}

	// Test unauthorized delete
	survey2 := &models.Survey{
		CreatorID: 2, // Different user
		Title:     "Another Survey",
		IsActive:  true,
	}
	survey2ID, _ := mockRepo.CreateSurvey(ctx, survey2)

	otherUserCtx := setupTestContext(1, []string{"user"})
	err = service.DeleteSurvey(otherUserCtx, survey2ID)
	if err == nil {
		t.Error("Expected error for unauthorized delete but got nil")
	}

	// Test admin can delete any survey
	adminCtx := setupTestContext(3, []string{"admin"})
	survey3 := &models.Survey{
		CreatorID: 2,
		Title:     "Yet Another Survey",
		IsActive:  true,
	}
	survey3ID, _ := mockRepo.CreateSurvey(ctx, survey3)

	err = service.DeleteSurvey(adminCtx, survey3ID)
	if err != nil {
		t.Fatalf("Unexpected error when admin deleting survey: %v", err)
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	err = service.DeleteSurvey(userCtx, surveyID)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestUpdateSurveyStatus(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create a test survey
	ctx := context.Background()
	survey := &models.Survey{
		CreatorID: 1,
		Title:     "Test Survey",
		IsActive:  false,
	}
	surveyID, _ := mockRepo.CreateSurvey(ctx, survey)

	// Create user context
	userCtx := setupTestContext(1, []string{"user"})

	// Update status to active
	err := service.UpdateSurveyStatus(userCtx, surveyID, true)
	if err != nil {
		t.Fatalf("Unexpected error updating status: %v", err)
	}

	// Verify status was updated
	updatedSurvey, _ := mockRepo.GetSurvey(ctx, surveyID)
	if !updatedSurvey.IsActive {
		t.Error("Expected survey to be active")
	}

	// Test unauthorized update
	otherUserCtx := setupTestContext(999, []string{"user"})
	err = service.UpdateSurveyStatus(otherUserCtx, surveyID, false)
	if err == nil {
		t.Error("Expected error for unauthorized update but got nil")
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	err = service.UpdateSurveyStatus(userCtx, surveyID, true)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestListUserSurveys(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create simple context for setup
	ctx := context.Background()

	// Create test surveys for user 1
	for i := 0; i < 5; i++ {
		survey := &models.Survey{
			CreatorID: 1,
			Title:     "User 1 Survey",
			IsActive:  true,
		}
		mockRepo.CreateSurvey(ctx, survey)
	}

	// Create test surveys for user 2
	for i := 0; i < 3; i++ {
		survey := &models.Survey{
			CreatorID: 2,
			Title:     "User 2 Survey",
			IsActive:  true,
		}
		mockRepo.CreateSurvey(ctx, survey)
	}

	// Create user context
	userCtx := setupTestContext(1, []string{"user"})

	// List surveys for user 1
	surveys, total, err := service.ListUserSurveys(userCtx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error listing surveys: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total count 5, got %d", total)
	}

	if len(surveys) != 5 {
		t.Errorf("Expected 5 surveys, got %d", len(surveys))
	}

	// Test pagination
	surveys, total, err = service.ListUserSurveys(userCtx, 0, 2)
	if err != nil {
		t.Fatalf("Unexpected error listing surveys with pagination: %v", err)
	}

	if total != 5 { // Total should still be 5
		t.Errorf("Expected total count 5, got %d", total)
	}

	if len(surveys) != 2 { // But only 2 returned due to limit
		t.Errorf("Expected 2 surveys, got %d", len(surveys))
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	_, _, err = service.ListUserSurveys(userCtx, 0, 10)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestListAllPublicSurveys(t *testing.T) {
	mockRepo := mock.NewMockRepository()
	service := NewSurveyService(mockRepo)

	// Create simple context for setup
	ctx := context.Background()

	// Create active surveys
	for i := 0; i < 3; i++ {
		survey := &models.Survey{
			CreatorID: 1,
			Title:     "Active Survey",
			IsActive:  true,
		}
		mockRepo.CreateSurvey(ctx, survey)
	}

	// Create inactive surveys
	for i := 0; i < 2; i++ {
		survey := &models.Survey{
			CreatorID: 1,
			Title:     "Inactive Survey",
			IsActive:  false,
		}
		mockRepo.CreateSurvey(ctx, survey)
	}

	// Create user context
	userCtx := setupTestContext(2, []string{"user"}) // Different user

	// List public surveys
	surveys, total, err := service.ListAllPublicSurveys(userCtx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error listing public surveys: %v", err)
	}

	// Using total to avoid linter error
	if total < 3 {
		t.Errorf("Expected total count to be at least 3, got %d", total)
	}

	// Should only return active surveys
	if len(surveys) != 3 {
		t.Errorf("Expected 3 active surveys, got %d", len(surveys))
	}

	// Admin should see all surveys
	adminCtx := setupTestContext(3, []string{"admin"})
	adminSurveys, adminTotal, err := service.ListAllPublicSurveys(adminCtx, 0, 10)
	if err != nil {
		t.Fatalf("Unexpected error listing surveys for admin: %v", err)
	}

	// Using adminTotal to avoid linter error
	if adminTotal < 5 {
		t.Errorf("Expected admin total count to be at least 5, got %d", adminTotal)
	}

	// Admin should see all 5 surveys (active and inactive)
	if len(adminSurveys) != 5 {
		t.Errorf("Expected admin to see all 5 surveys, got %d", len(adminSurveys))
	}

	// Test error handling
	mockRepo.ErrorMock = errors.New("test error")
	_, _, err = service.ListAllPublicSurveys(userCtx, 0, 10)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}
