package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/model"
)

// SurveyClient defines the interface for interacting with the survey service
type SurveyClient interface {
	GetSurvey(ctx context.Context, surveyID string) (*model.SurveyPublic, error)
}

// HTTPSurveyClient implements SurveyClient using HTTP requests
type HTTPSurveyClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPSurveyClient creates a new HTTP-based survey client
func NewHTTPSurveyClient(baseURL string) *HTTPSurveyClient {
	return &HTTPSurveyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSurvey retrieves a survey by ID from the survey service
func (c *HTTPSurveyClient) GetSurvey(ctx context.Context, surveyID string) (*model.SurveyPublic, error) {
	// Construct request URL
	url := fmt.Sprintf("%s/surveys/%s", c.baseURL, surveyID)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var survey model.SurveyPublic
	if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &survey, nil
}

// MockSurveyClient is a mock implementation of SurveyClient for testing
type MockSurveyClient struct {
	surveys map[string]*model.SurveyPublic
}

// NewMockSurveyClient creates a new mock survey client
func NewMockSurveyClient() *MockSurveyClient {
	return &MockSurveyClient{
		surveys: make(map[string]*model.SurveyPublic),
	}
}

// AddSurvey adds a survey to the mock client
func (c *MockSurveyClient) AddSurvey(survey *model.SurveyPublic) {
	c.surveys[survey.ID] = survey
}

// GetSurvey retrieves a survey by ID from the mock client
func (c *MockSurveyClient) GetSurvey(ctx context.Context, surveyID string) (*model.SurveyPublic, error) {
	survey, ok := c.surveys[surveyID]
	if !ok {
		return nil, fmt.Errorf("survey not found: %s", surveyID)
	}
	return survey, nil
}
