// backend/tests/contracts/survey_response_message_test.go
package contracts

import (
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/assert"
)

// Структура для відповіді на опитування
type SurveyResponseMessage struct {
	ResponseID  string               `json:"response_id"`
	SurveyID    string               `json:"survey_id"`
	UserID      string               `json:"user_id"`
	SubmittedAt time.Time            `json:"submitted_at"`
	Answers     []QuestionAnswerPair `json:"answers"`
}

type QuestionAnswerPair struct {
	QuestionID string      `json:"question_id"`
	Type       string      `json:"type"`
	Answer     interface{} `json:"answer"`
}

func TestSurveyResponseMessagePact(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "survey_taking_service",
		Provider: "response_processor_service",
	}

	// Налаштування Message Pact
	messagePact := dsl.MessagePact{
		Consumer: "survey_taking_service",
		Provider: "response_processor_service",
	}

	// Визначення прикладу повідомлення для контракту
	surveyResponseExample := SurveyResponseMessage{
		ResponseID:  "resp-123e4567-e89b-12d3-a456-426614174000",
		SurveyID:    "123e4567-e89b-12d3-a456-426614174000",
		UserID:      "user-123e4567-e89b-12d3-a456-426614174001",
		SubmittedAt: time.Now().UTC(),
		Answers: []QuestionAnswerPair{
			{
				QuestionID: "q1",
				Type:       "single-choice",
				Answer:     "opt1",
			},
			{
				QuestionID: "q2",
				Type:       "multiple-choice",
				Answer:     []string{"opt1", "opt3"},
			},
			{
				QuestionID: "q3",
				Type:       "open-text",
				Answer:     "This is a text answer",
			},
		},
	}

	// Створення контракту для повідомлення
	messageDescriptor := pact.AddMessage()
	messageDescriptor.
		Given("Відповідь на опитування готова до обробки").
		ExpectsToReceive("Повідомлення з відповіддю на опитування").
		WithContent(types.JSONMatcher{
			"response_id":  dsl.Term("resp-123e4567-e89b-12d3-a456-426614174000", `^resp-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
			"survey_id":    dsl.Term("123e4567-e89b-12d3-a456-426614174000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
			"user_id":      dsl.Term("user-123e4567-e89b-12d3-a456-426614174001", `^user-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
			"submitted_at": dsl.Term("2023-04-01T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?$`),
			"answers": dsl.EachLike(map[string]interface{}{
				"question_id": dsl.Term("q1", `^q\d+$`),
				"type":        dsl.Term("single-choice", "(single-choice|multiple-choice|open-text|scale|matrix-single|matrix-multiple)"),
				"answer":      dsl.Like("opt1"), // Note: This is a simplification as answer can be different types
			}, 1),
		}).
		AsType(&surveyResponseExample)

	// Перевірка відповідності прикладу контракту
	err := pact.VerifyMessageConsumer(&surveyResponseExample, nil)
	assert.NoError(t, err)
}
