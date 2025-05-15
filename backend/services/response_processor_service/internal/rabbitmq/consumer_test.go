package rabbitmq

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/streadway/amqp"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// MockDelivery mocks the amqp.Delivery structure for testing
type MockDelivery struct {
	body         []byte
	ackCalled    bool
	rejectCalled bool
	requeue      bool
}

// Ack implements the Acknowledger interface
func (m *MockDelivery) Ack(multiple bool) error {
	m.ackCalled = true
	return nil
}

// Nack implements the Acknowledger interface
func (m *MockDelivery) Nack(multiple bool, requeue bool) error {
	return nil
}

// Reject implements the Acknowledger interface
func (m *MockDelivery) Reject(requeue bool) error {
	m.rejectCalled = true
	m.requeue = requeue
	return nil
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name    string
		message model.RabbitMQMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid message",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011", // Valid ObjectID
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012", // Valid ObjectID
						Value:      "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "Missing SurveyID",
			message: model.RabbitMQMessage{
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012",
						Value:      "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "survey_id is required",
		},
		{
			name: "Invalid SurveyID format",
			message: model.RabbitMQMessage{
				SurveyID:     "invalid-id",
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012",
						Value:      "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "invalid survey_id format",
		},
		{
			name: "Missing both RespondentID and AnonymousID",
			message: model.RabbitMQMessage{
				SurveyID: "507f1f77bcf86cd799439011",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012",
						Value:      "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "either respondent_id or anonymous_id must be provided",
		},
		{
			name: "Empty Answers",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011",
				RespondentID: "user123",
				Answers:      []model.Answer{},
				SubmittedAt:  "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "at least one answer must be provided",
		},
		{
			name: "Answer missing QuestionID",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011",
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						Value: "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "answer at index 0 is missing question_id",
		},
		{
			name: "Answer with invalid QuestionID format",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011",
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "invalid-id",
						Value:      "answer1",
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "invalid question_id format for answer at index 0",
		},
		{
			name: "Answer missing Value",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011",
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012",
						Value:      nil,
					},
				},
				SubmittedAt: "2023-05-01T12:00:00Z",
			},
			wantErr: true,
			errMsg:  "answer at index 0 is missing value",
		},
		{
			name: "Missing SubmittedAt",
			message: model.RabbitMQMessage{
				SurveyID:     "507f1f77bcf86cd799439011",
				RespondentID: "user123",
				Answers: []model.Answer{
					{
						QuestionID: "507f1f77bcf86cd799439012",
						Value:      "answer1",
					},
				},
			},
			wantErr: true,
			errMsg:  "submitted_at timestamp is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessage(&tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if msg := err.Error(); msg[:len(tt.errMsg)] != tt.errMsg {
					t.Errorf("validateMessage() error message = %v, want to start with %v", msg, tt.errMsg)
				}
			}
		})
	}
}

func TestHandleDelivery(t *testing.T) {
	validMessage := model.RabbitMQMessage{
		SurveyID:     "507f1f77bcf86cd799439011", // Valid ObjectID
		RespondentID: "user123",
		Answers: []model.Answer{
			{
				QuestionID: "507f1f77bcf86cd799439012", // Valid ObjectID
				Value:      "answer1",
			},
		},
		SubmittedAt: "2023-05-01T12:00:00Z",
	}

	validJSON, _ := json.Marshal(validMessage)
	invalidJSON := []byte(`{"invalid": json}`)

	invalidMessage := model.RabbitMQMessage{} // Missing required fields

	invalidMessageJSON, _ := json.Marshal(invalidMessage)

	tests := []struct {
		name             string
		delivery         *MockDelivery
		repository       *db.MockResponseRepository
		wantReject       bool
		wantRequeue      bool
		wantAck          bool
		wantProcessCount int
	}{
		{
			name: "Valid message",
			delivery: &MockDelivery{
				body: validJSON,
			},
			repository: &db.MockResponseRepository{
				ProcessedMessages: []*model.RabbitMQMessage{},
			},
			wantReject:       false,
			wantAck:          true,
			wantProcessCount: 1,
		},
		{
			name: "Invalid JSON",
			delivery: &MockDelivery{
				body: invalidJSON,
			},
			repository: &db.MockResponseRepository{
				ProcessedMessages: []*model.RabbitMQMessage{},
			},
			wantReject:       true,
			wantRequeue:      false, // Don't requeue for invalid JSON
			wantProcessCount: 0,
		},
		{
			name: "Invalid message (validation fails)",
			delivery: &MockDelivery{
				body: invalidMessageJSON,
			},
			repository: &db.MockResponseRepository{
				ProcessedMessages: []*model.RabbitMQMessage{},
			},
			wantReject:       true,
			wantRequeue:      false, // Don't requeue for validation failures
			wantProcessCount: 0,
		},
		{
			name: "Repository processing error",
			delivery: &MockDelivery{
				body: validJSON,
			},
			repository: &db.MockResponseRepository{
				ShouldFail:        true,
				ErrorMsg:          "database error",
				ProcessedMessages: []*model.RabbitMQMessage{},
			},
			wantReject:       true,
			wantRequeue:      true, // Requeue for processing errors
			wantProcessCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := &Consumer{
				repository: tt.repository,
			}

			mockDelivery := amqp.Delivery{
				Body:         tt.delivery.body,
				Acknowledger: tt.delivery,
			}

			consumer.handleDelivery(context.Background(), mockDelivery)

			if tt.delivery.rejectCalled != tt.wantReject {
				t.Errorf("handleDelivery() reject called = %v, want %v", tt.delivery.rejectCalled, tt.wantReject)
			}

			if tt.delivery.rejectCalled && tt.delivery.requeue != tt.wantRequeue {
				t.Errorf("handleDelivery() requeue = %v, want %v", tt.delivery.requeue, tt.wantRequeue)
			}

			if tt.delivery.ackCalled != tt.wantAck {
				t.Errorf("handleDelivery() ack called = %v, want %v", tt.delivery.ackCalled, tt.wantAck)
			}

			if len(tt.repository.ProcessedMessages) != tt.wantProcessCount {
				t.Errorf("handleDelivery() processed messages count = %v, want %v", len(tt.repository.ProcessedMessages), tt.wantProcessCount)
			}
		})
	}
}
