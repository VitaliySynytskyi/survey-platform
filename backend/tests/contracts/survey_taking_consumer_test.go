// backend/tests/contracts/survey_taking_consumer_test.go
package contracts

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

// Функція для імітації запиту від Survey Taking Service до Survey Service
func getSurveyDetailsForTaking(pact dsl.Pact, surveyID string) error {
	url := fmt.Sprintf("http://localhost:%d/api/v1/surveys/%s", pact.Server.Port, surveyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer mocked_token")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status 200 but got %d", resp.StatusCode)
	}
	return nil
}

func TestSurveyTakingServiceConsumer(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "survey_taking_service",
		Provider: "survey_service",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// Тестовий ID опитування
	testSurveyID := "123e4567-e89b-12d3-a456-426614174000"

	pact.AddInteraction().
		Given("Опитування існує").
		UponReceiving("Запит на отримання деталей опитування для проходження").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term(fmt.Sprintf("/api/v1/surveys/%s", testSurveyID), "/api/v1/surveys/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
			Headers: dsl.MapMatcher{
				"Authorization": dsl.Term("Bearer mocked_token", `^Bearer\s+[a-zA-Z0-9-_.]+$`),
				"Content-Type":  dsl.Term("application/json", "application/json"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json", "application/json(;?.*)"),
			},
			Body: dsl.Like(map[string]interface{}{
				"id":          dsl.Term(testSurveyID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
				"title":       dsl.Like("Test Survey"),
				"description": dsl.Like("Survey Description"),
				"owner_id":    dsl.Term("123e4567-e89b-12d3-a456-426614174001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
				"created_at":  dsl.Term("2023-04-01T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?$`),
				"updated_at":  dsl.Term("2023-04-01T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?$`),
				"status":      dsl.Term("active", "(draft|active|closed)"),
				"questions": dsl.EachLike(map[string]interface{}{
					"question_id": dsl.Term("q1", `^q\d+$`),
					"text":        dsl.Like("What is your favorite color?"),
					"type":        dsl.Term("single-choice", "(single-choice|multiple-choice|open-text|scale|matrix-single|matrix-multiple)"),
					"is_required": dsl.Like(true),
					"order":       dsl.Like(1),
					"options": dsl.EachLike(map[string]interface{}{
						"option_id": dsl.Term("opt1", `^opt\d+$`),
						"text":      dsl.Like("Red"),
						"order":     dsl.Like(1),
					}, 1),
					"logic_rules": dsl.EachLike(map[string]interface{}{
						"operator":           dsl.Term("equals", "(equals|not_equals|contains|greater_than|less_than)"),
						"value":              dsl.Like("opt1"),
						"target_question_id": dsl.Term("q2", `^q\d+$`),
					}, 0),
				}, 1),
			}),
		})

	err := pact.Verify(func() error {
		return getSurveyDetailsForTaking(pact, testSurveyID)
	})
	assert.NoError(t, err)
}
