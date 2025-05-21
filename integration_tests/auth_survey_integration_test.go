package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Integration test between Auth Service and Survey Service
// This test should be run with a test environment where all services are running

func TestAuthAndSurveyIntegration(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Configuration
	baseAuthURL := "http://localhost:8081/api/v1"
	baseSurveyURL := "http://localhost:8082/api/v1"

	// Test user credentials
	username := fmt.Sprintf("testuser%d", time.Now().Unix())
	email := fmt.Sprintf("test%d@example.com", time.Now().Unix())
	password := "Password123!"

	var authToken string
	var userID int

	// Step 1: Register a new user
	t.Run("Register User", func(t *testing.T) {
		registerURL := fmt.Sprintf("%s/auth/register", baseAuthURL)

		requestBody := map[string]interface{}{
			"username":   username,
			"email":      email,
			"password":   password,
			"first_name": "Test",
			"last_name":  "User",
		}

		jsonBody, _ := json.Marshal(requestBody)

		resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		// Save auth token for subsequent requests
		authToken = responseData["token"].(string)
		userData := responseData["user"].(map[string]interface{})
		userID = int(userData["id"].(float64))

		resp.Body.Close()
	})

	// Step 2: Create a new survey
	var surveyID string

	t.Run("Create Survey", func(t *testing.T) {
		createSurveyURL := fmt.Sprintf("%s/surveys", baseSurveyURL)

		// Sample survey with different question types
		requestBody := map[string]interface{}{
			"title":       "Integration Test Survey",
			"description": "This survey was created during integration testing",
			"questions": []map[string]interface{}{
				{
					"text":     "What is your favorite color?",
					"type":     "TEXT",
					"required": true,
				},
				{
					"text":     "Rate your experience from 1 to 5",
					"type":     "RATING",
					"required": true,
					"options": map[string]interface{}{
						"min": 1,
						"max": 5,
					},
				},
				{
					"text":     "Select your preferred programming languages",
					"type":     "MULTIPLE_CHOICE",
					"required": false,
					"options": map[string]interface{}{
						"choices":  []string{"Go", "JavaScript", "Python", "Java", "C#"},
						"multiple": true,
					},
				},
			},
		}

		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", createSurveyURL, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		// Save survey ID for subsequent requests
		surveyID = responseData["id"].(string)

		resp.Body.Close()
	})

	// Step 3: Get the created survey
	t.Run("Get Survey", func(t *testing.T) {
		getSurveyURL := fmt.Sprintf("%s/surveys/%s", baseSurveyURL, surveyID)

		req, _ := http.NewRequest("GET", getSurveyURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		// Verify survey data
		assert.Equal(t, surveyID, responseData["id"])
		assert.Equal(t, "Integration Test Survey", responseData["title"])
		assert.Equal(t, float64(userID), responseData["creator_id"])

		questions := responseData["questions"].([]interface{})
		assert.Equal(t, 3, len(questions))

		resp.Body.Close()
	})

	// Step 4: Update the survey
	t.Run("Update Survey", func(t *testing.T) {
		updateSurveyURL := fmt.Sprintf("%s/surveys/%s", baseSurveyURL, surveyID)

		requestBody := map[string]interface{}{
			"title":       "Updated Integration Test Survey",
			"description": "This survey was updated during integration testing",
		}

		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("PATCH", updateSurveyURL, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp.Body.Close()

		// Verify the update
		getSurveyURL := fmt.Sprintf("%s/surveys/%s", baseSurveyURL, surveyID)
		req, _ = http.NewRequest("GET", getSurveyURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		resp, err = client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		assert.Equal(t, "Updated Integration Test Survey", responseData["title"])

		resp.Body.Close()
	})

	// Step 5: Delete the survey
	t.Run("Delete Survey", func(t *testing.T) {
		deleteSurveyURL := fmt.Sprintf("%s/surveys/%s", baseSurveyURL, surveyID)

		req, _ := http.NewRequest("DELETE", deleteSurveyURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp.Body.Close()

		// Verify the survey was deleted
		getSurveyURL := fmt.Sprintf("%s/surveys/%s", baseSurveyURL, surveyID)
		req, _ = http.NewRequest("GET", getSurveyURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		resp, err = client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		resp.Body.Close()
	})
}
