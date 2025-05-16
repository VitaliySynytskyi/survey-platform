package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Базова URL API Gateway
	apiGatewayURL = "http://localhost:8080"
)

// Структури для авторизації
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

// Структури для опитувань
type Survey struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatorID   string     `json:"creator_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	Questions   []Question `json:"questions"`
}

type Question struct {
	ID         string      `json:"id,omitempty"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	Required   bool        `json:"required"`
	Options    []string    `json:"options,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

type Validation struct {
	Min       *int    `json:"min,omitempty"`
	Max       *int    `json:"max,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`
	MinLength *int    `json:"min_length,omitempty"`
	MaxLength *int    `json:"max_length,omitempty"`
}

// Структури для відповідей на опитування
type SurveyResponse struct {
	ID        string    `json:"id,omitempty"`
	SurveyID  string    `json:"survey_id"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Answers   []Answer  `json:"answers"`
}

type Answer struct {
	QuestionID string      `json:"question_id"`
	Value      interface{} `json:"value"`
}

// Допоміжна функція для виконання HTTP запитів
func makeRequest(t *testing.T, method, url string, body interface{}, token string) ([]byte, int) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return respBody, resp.StatusCode
}

// Тест, що перевіряє наскрізний сценарій: реєстрація -> логін -> створення опитування -> отримання опитування
func TestUserSurveyFlow(t *testing.T) {
	// Пропустимо тест, якщо не вказано запускати тести в інтеграційному режимі
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// Унікальний email для цього тесту, щоб уникнути конфліктів
	userEmail := fmt.Sprintf("test-user-%d@example.com", time.Now().UnixNano())

	// 1. Реєстрація нового користувача
	t.Log("Step 1: Registering a new user")
	user := User{
		Email:    userEmail,
		Password: "Test@123456",
		Name:     "Test User",
	}

	respBody, statusCode := makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/register", user, "")
	require.Equal(t, http.StatusCreated, statusCode)

	var regResponse map[string]interface{}
	err := json.Unmarshal(respBody, &regResponse)
	require.NoError(t, err)

	// Переконуємося, що реєстрація пройшла успішно
	assert.NotNil(t, regResponse["user_id"])

	// 2. Логін користувача
	t.Log("Step 2: Logging in")
	loginPayload := map[string]string{
		"email":    userEmail,
		"password": "Test@123456",
	}

	respBody, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/login", loginPayload, "")
	require.Equal(t, http.StatusOK, statusCode)

	var loginResp LoginResponse
	err = json.Unmarshal(respBody, &loginResp)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.UserID)

	userID := loginResp.UserID
	token := loginResp.AccessToken

	// 3. Створення опитування
	t.Log("Step 3: Creating a survey")
	survey := Survey{
		Title:       "Test Survey",
		Description: "This is a test survey created by the integration test",
		Questions: []Question{
			{
				Title:    "What is your age?",
				Type:     "number",
				Required: true,
				Validation: &Validation{
					Min: intPtr(18),
					Max: intPtr(100),
				},
			},
			{
				Title:    "What is your favorite color?",
				Type:     "single-choice",
				Required: true,
				Options:  []string{"Red", "Green", "Blue", "Yellow"},
			},
		},
	}

	respBody, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/surveys", survey, token)
	require.Equal(t, http.StatusCreated, statusCode)

	var createdSurvey Survey
	err = json.Unmarshal(respBody, &createdSurvey)
	require.NoError(t, err)
	assert.NotEmpty(t, createdSurvey.ID)
	assert.Equal(t, survey.Title, createdSurvey.Title)
	assert.Equal(t, 2, len(createdSurvey.Questions))

	surveyID := createdSurvey.ID

	// 4. Отримання опитування користувача
	t.Log("Step 4: Getting user's surveys")
	respBody, statusCode = makeRequest(t, "GET", fmt.Sprintf("%s/api/v1/users/%s/surveys", apiGatewayURL, userID), nil, token)
	require.Equal(t, http.StatusOK, statusCode)

	var userSurveys map[string]interface{}
	err = json.Unmarshal(respBody, &userSurveys)
	require.NoError(t, err)

	// Переконуємося, що список опитувань містить створене опитування
	surveys, ok := userSurveys["surveys"].([]interface{})
	require.True(t, ok)

	found := false
	for _, s := range surveys {
		surveyMap, ok := s.(map[string]interface{})
		if ok && surveyMap["id"] == surveyID {
			found = true
			break
		}
	}
	assert.True(t, found, "The created survey should be in the user's surveys list")

	// 5. Отримання опитування за ID
	t.Log("Step 5: Getting survey by ID")
	respBody, statusCode = makeRequest(t, "GET", apiGatewayURL+"/api/v1/surveys/"+surveyID, nil, token)
	require.Equal(t, http.StatusOK, statusCode)

	var retrievedSurvey Survey
	err = json.Unmarshal(respBody, &retrievedSurvey)
	require.NoError(t, err)
	assert.Equal(t, surveyID, retrievedSurvey.ID)
	assert.Equal(t, survey.Title, retrievedSurvey.Title)
	assert.Len(t, retrievedSurvey.Questions, 2)
}

// Допоміжна функція для створення вказівника на int
func intPtr(i int) *int {
	return &i
}

// Тест наскрізного сценарію: створення опитування -> проходження опитування -> перевірка результатів
func TestSurveyResponseFlow(t *testing.T) {
	// Пропустимо тест, якщо не вказано запускати тести в інтеграційному режимі
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// Унікальний email для цього тесту, щоб уникнути конфліктів
	userEmail := fmt.Sprintf("test-user-%d@example.com", time.Now().UnixNano())

	// 1. Реєстрація користувача-організатора опитувань
	user := User{
		Email:    userEmail,
		Password: "Test@123456",
		Name:     "Survey Creator",
	}

	_, statusCode := makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/register", user, "")
	require.Equal(t, http.StatusCreated, statusCode)

	// 2. Логін користувача-організатора
	loginPayload := map[string]string{
		"email":    userEmail,
		"password": "Test@123456",
	}

	respBody, statusCode := makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/login", loginPayload, "")
	require.Equal(t, http.StatusOK, statusCode)

	var loginResp LoginResponse
	err := json.Unmarshal(respBody, &loginResp)
	require.NoError(t, err)
	creatorToken := loginResp.AccessToken

	// 3. Створення опитування
	survey := Survey{
		Title:       "Customer Satisfaction Survey",
		Description: "Please provide your feedback on our services",
		Questions: []Question{
			{
				Title:    "How satisfied are you with our service?",
				Type:     "scale",
				Required: true,
				Options:  []string{"1", "2", "3", "4", "5"},
			},
			{
				Title:    "What improvements would you like to see?",
				Type:     "open-text",
				Required: false,
			},
		},
	}

	respBody, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/surveys", survey, creatorToken)
	require.Equal(t, http.StatusCreated, statusCode)

	var createdSurvey Survey
	err = json.Unmarshal(respBody, &createdSurvey)
	require.NoError(t, err)
	surveyID := createdSurvey.ID

	// Отримуємо IDs питань
	questionIDs := []string{createdSurvey.Questions[0].ID, createdSurvey.Questions[1].ID}

	// 4. Реєстрація користувача-респондента
	respondentEmail := fmt.Sprintf("test-respondent-%d@example.com", time.Now().UnixNano())
	respondent := User{
		Email:    respondentEmail,
		Password: "Test@123456",
		Name:     "Survey Respondent",
	}

	_, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/register", respondent, "")
	require.Equal(t, http.StatusCreated, statusCode)

	// 5. Логін користувача-респондента
	loginPayload = map[string]string{
		"email":    respondentEmail,
		"password": "Test@123456",
	}

	respBody, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/auth/login", loginPayload, "")
	require.Equal(t, http.StatusOK, statusCode)

	var respondentLogin LoginResponse
	err = json.Unmarshal(respBody, &respondentLogin)
	require.NoError(t, err)
	respondentToken := respondentLogin.AccessToken
	respondentID := respondentLogin.UserID

	// 6. Відправлення відповідей на опитування
	response := SurveyResponse{
		SurveyID: surveyID,
		UserID:   respondentID,
		Answers: []Answer{
			{
				QuestionID: questionIDs[0],
				Value:      "5", // Максимальна оцінка
			},
			{
				QuestionID: questionIDs[1],
				Value:      "More mobile-friendly interface would be great!",
			},
		},
	}

	respBody, statusCode = makeRequest(t, "POST", apiGatewayURL+"/api/v1/take/"+surveyID+"/responses", response, respondentToken)
	require.Equal(t, http.StatusCreated, statusCode)

	// Почекайте, поки відповідь буде оброблена асинхронно
	time.Sleep(2 * time.Second)

	// 7. Отримання результатів опитування (організатором)
	respBody, statusCode = makeRequest(t, "GET", apiGatewayURL+"/api/v1/analytics/surveys/"+surveyID+"/results", nil, creatorToken)
	require.Equal(t, http.StatusOK, statusCode)

	var results map[string]interface{}
	err = json.Unmarshal(respBody, &results)
	require.NoError(t, err)

	// Переконуємося, що результати містять відповіді
	assert.Equal(t, surveyID, results["survey_id"])
	assert.Equal(t, float64(1), results["total_responses"]) // Повинен бути один респондент
}
