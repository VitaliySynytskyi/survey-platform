package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	log.Printf("SURVEY_TAKING_CLIENT: GetSurvey called. BaseURL: '%s', SurveyID: '%s'", c.baseURL, surveyID)
	// Construct request URL
	url := fmt.Sprintf("%s/%s/public", c.baseURL, surveyID)
	log.Printf("SURVEY_TAKING_CLIENT: Constructed survey service request URL: %s", url)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("SURVEY_TAKING_CLIENT: Error creating request to %s: %v", url, err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	log.Printf("SURVEY_TAKING_CLIENT: Sending GET request to %s", url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("SURVEY_TAKING_CLIENT: Error sending request to %s: %v", url, err)
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("SURVEY_TAKING_CLIENT: Received response from %s. Status: %s, StatusCode: %d", url, resp.Status, resp.StatusCode)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("SURVEY_TAKING_CLIENT: Unexpected status code %d from %s", resp.StatusCode, url)
		// It might be useful to read the body here for more error details, if any
		// e.g., bodyBytes, _ := io.ReadAll(resp.Body); log.Printf("Response body: %s", string(bodyBytes))
		return nil, fmt.Errorf("unexpected status code: %d from survey service", resp.StatusCode) // Simplified error message to user
	}

	// Parse response
	var survey model.SurveyPublic
	if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
		log.Printf("SURVEY_TAKING_CLIENT: Error decoding JSON response from %s: %v", url, err)
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("SURVEY_TAKING_CLIENT: Successfully decoded survey response for ID %s. Title: %s", survey.ID, survey.Title)
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
