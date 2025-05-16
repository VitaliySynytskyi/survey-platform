// backend/tests/contracts/survey_api_consumer_test.go
package contracts

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

// Функція для імітації запиту від API Gateway до Survey Service
func getSurveysFromSurveyService(pact dsl.Pact) error {
	// Приклад URL, який API Gateway міг би викликати (з mock-сервера Pact)
	// Важливо: pact.Server.Port буде динамічним
	url := fmt.Sprintf("http://localhost:%d/api/v1/surveys", pact.Server.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer mocked_token") // Використовуй реальний або мокований токен
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status 200 but got %d", resp.StatusCode)
	}
	// Тут можна додати перевірку тіла відповіді, якщо потрібно,
	// але основна перевірка структури виконується Pact
	return nil
}

func TestSurveyServiceConsumer(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "api_gateway",
		Provider: "survey_service",
		Host:     "localhost", // Де Pact буде запускати mock-сервер
	}
	defer pact.Teardown()

	pact.AddInteraction().
		Given("Існують опитування"). // Стан провайдера
		UponReceiving("Запит на отримання всіх опитувань від api_gateway").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/api/v1/surveys", "/api/v1/surveys"), // Шлях, який очікує survey_service
			Headers: dsl.MapMatcher{
				"Authorization": dsl.Term("Bearer mocked_token", `^Bearer\s+[a-zA-Z0-9-_.]+$`),
				"Content-Type":  dsl.Term("application/json", "application/json"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json", "application/json(;?.*)"), // Дозволяємо charset
			},
			Body: dsl.Like(map[string]interface{}{
				"surveys": dsl.EachLike(map[string]interface{}{
					"id":          dsl.Term("123e4567-e89b-12d3-a456-426614174000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"title":       dsl.Like("Test Survey"),
					"description": dsl.Like("Description"),
					"owner_id":    dsl.Term("123e4567-e89b-12d3-a456-426614174001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`), // Змінено на owner_id згідно моделі SurveyResponse
					"created_at":  dsl.Term("2023-04-01T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?$`),
					"updated_at":  dsl.Term("2023-04-01T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?$`), // Додано updated_at
					"questions": dsl.EachLike(map[string]interface{}{ // Додано структуру Question
						"question_id": dsl.Term("q1", `^q\d+$`),
						"text":        dsl.Like("What is your favorite color?"),
						"type":        dsl.Term("single-choice", "(single-choice|multiple-choice|open-text|scale|matrix-single|matrix-multiple)"),
						"is_required": dsl.Like(true),
						// "options" and "scale_settings" can be added here if needed,
						// but dsl.Like will match if they are present or not.
						// For stricter matching, define them explicitly.
					}, 1),
				}, 1),
				"total_count": dsl.Like(10), // Змінено з pagination.total
				"page":        dsl.Like(1),  // Змінено з pagination.page
				"per_page":    dsl.Like(10), // Змінено з pagination.page_size
			}),
		})

	// Виконання тесту
	err := pact.Verify(func() error {
		// Ця функція викликає реальний код консьюмера (api_gateway),
		// який робить запит до mock-сервера Pact.
		// Зараз ми просто імітуємо запит для генерації контракту.
		// У реальному сценарії тут був би виклик функції з api_gateway,
		// яка робить запит до survey_service.
		return getSurveysFromSurveyService(pact)
	})
	assert.NoError(t, err)

	// Опціонально: опублікувати контракт на Pact Broker
	// pact.PublishContracts("YOUR_PACT_BROKER_URL", "YOUR_PACT_BROKER_TOKEN", "YOUR_APP_VERSION")
}
