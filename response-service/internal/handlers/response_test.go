package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service for testing
type MockResponseService struct {
	mock.Mock
}

func (m *MockResponseService) CreateResponse(ctx context.Context, response *models.Response) (string, error) {
	args := m.Called(ctx, response)
	return args.String(0), args.Error(1)
}

func (m *MockResponseService) GetResponses(ctx context.Context, surveyID string) ([]models.Response, error) {
	args := m.Called(ctx, surveyID)
	return args.Get(0).([]models.Response), args.Error(1)
}

func (m *MockResponseService) GetResponseByID(ctx context.Context, id string) (*models.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Response), args.Error(1)
}

func (m *MockResponseService) GetUserResponses(ctx context.Context, userID int) ([]models.Response, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Response), args.Error(1)
}

func (m *MockResponseService) GetSurveyStats(ctx context.Context, surveyID string) (*models.SurveyStats, error) {
	args := m.Called(ctx, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SurveyStats), args.Error(1)
}

func setupTestRouter(mockService *MockResponseService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create handler with mock service
	handler := NewResponseHandler(mockService)

	// Set up routes
	r.POST("/api/v1/responses", handler.CreateResponse)
	r.GET("/api/v1/responses/survey/:survey_id", handler.GetResponses)
	r.GET("/api/v1/responses/:id", handler.GetResponseByID)
	r.GET("/api/v1/responses/stats/:survey_id", handler.GetSurveyStats)

	return r
}

func TestCreateResponse(t *testing.T) {
	mockService := new(MockResponseService)
	router := setupTestRouter(mockService)

	t.Run("Create response successfully", func(t *testing.T) {
		// Test data
		responseID := primitive.NewObjectID().Hex()
		response := models.Response{
			SurveyID: "survey123",
			UserID:   123,
			Answers: []models.Answer{
				{
					QuestionID: "q1",
					Value:      "blue",
				},
				{
					QuestionID: "q2",
					Value:      "4",
				},
			},
		}

		// Set up mock service expectation
		mockService.On("CreateResponse", mock.Anything, mock.AnythingOfType("*models.Response")).Return(responseID, nil).Once()

		// Create request
		jsonData, _ := json.Marshal(response)
		req, _ := http.NewRequest("POST", "/api/v1/responses", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "123") // Mock middleware

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var result map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, responseID, result["id"])

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Missing user ID", func(t *testing.T) {
		// Create request without X-User-ID header
		response := models.Response{
			SurveyID: "survey123",
			Answers:  []models.Answer{},
		}

		jsonData, _ := json.Marshal(response)
		req, _ := http.NewRequest("POST", "/api/v1/responses", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		// No X-User-ID header

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetResponses(t *testing.T) {
	mockService := new(MockResponseService)
	router := setupTestRouter(mockService)

	t.Run("Get responses for survey", func(t *testing.T) {
		// Test data
		surveyID := "survey123"
		responses := []models.Response{
			{
				ID:       primitive.NewObjectID().Hex(),
				SurveyID: surveyID,
				UserID:   123,
				Answers: []models.Answer{
					{
						QuestionID: "q1",
						Value:      "red",
					},
				},
			},
			{
				ID:       primitive.NewObjectID().Hex(),
				SurveyID: surveyID,
				UserID:   456,
				Answers: []models.Answer{
					{
						QuestionID: "q1",
						Value:      "blue",
					},
				},
			},
		}

		// Set up mock service expectation
		mockService.On("GetResponses", mock.Anything, surveyID).Return(responses, nil).Once()

		// Create request
		req, _ := http.NewRequest("GET", "/api/v1/responses/survey/"+surveyID, nil)

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var result []models.Response
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, surveyID, result[0].SurveyID)

		// Verify mock
		mockService.AssertExpectations(t)
	})
}

func TestGetSurveyStats(t *testing.T) {
	mockService := new(MockResponseService)
	router := setupTestRouter(mockService)

	t.Run("Get survey statistics", func(t *testing.T) {
		// Test data
		surveyID := "survey123"
		stats := &models.SurveyStats{
			SurveyID:       surveyID,
			TotalResponses: 2,
			CompletionRate: 100.0,
			QuestionStats: map[string]models.QuestionStat{
				"q1": {
					QuestionID: "q1",
					AnswerDistribution: map[string]int{
						"red":  1,
						"blue": 1,
					},
				},
			},
		}

		// Set up mock service expectation
		mockService.On("GetSurveyStats", mock.Anything, surveyID).Return(stats, nil).Once()

		// Create request
		req, _ := http.NewRequest("GET", "/api/v1/responses/stats/"+surveyID, nil)

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var result models.SurveyStats
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, surveyID, result.SurveyID)
		assert.Equal(t, 2, result.TotalResponses)
		assert.Equal(t, 100.0, result.CompletionRate)
		assert.Contains(t, result.QuestionStats, "q1")
		assert.Equal(t, 2, len(result.QuestionStats["q1"].AnswerDistribution))

		// Verify mock
		mockService.AssertExpectations(t)
	})
}

func TestGetResponseByID(t *testing.T) {
	mockService := new(MockResponseService)
	router := setupTestRouter(mockService)

	t.Run("Get response by ID", func(t *testing.T) {
		// Test data
		responseID := primitive.NewObjectID().Hex()
		response := &models.Response{
			ID:       responseID,
			SurveyID: "survey123",
			UserID:   123,
			Answers: []models.Answer{
				{
					QuestionID: "q1",
					Value:      "green",
				},
			},
		}

		// Set up mock service expectation
		mockService.On("GetResponseByID", mock.Anything, responseID).Return(response, nil).Once()

		// Create request
		req, _ := http.NewRequest("GET", "/api/v1/responses/"+responseID, nil)

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var result models.Response
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, responseID, result.ID)
		assert.Equal(t, "survey123", result.SurveyID)
		assert.Equal(t, 1, len(result.Answers))
		assert.Equal(t, "green", result.Answers[0].Value)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Response not found", func(t *testing.T) {
		// Test data
		responseID := primitive.NewObjectID().Hex()

		// Set up mock service expectation
		mockService.On("GetResponseByID", mock.Anything, responseID).Return(nil, service.ErrResponseNotFound).Once()

		// Create request
		req, _ := http.NewRequest("GET", "/api/v1/responses/"+responseID, nil)

		// Perform request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Verify mock
		mockService.AssertExpectations(t)
	})
}
