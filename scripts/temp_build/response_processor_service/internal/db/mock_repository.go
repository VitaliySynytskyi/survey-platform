package db

import (
	"context"
	"fmt"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// MockResponseRepository is a mock implementation of the ResponseRepository for testing
type MockResponseRepository struct {
	ProcessRabbitMQMessageFunc func(ctx context.Context, message *model.RabbitMQMessage) error
	SaveResponseFunc           func(ctx context.Context, response *model.SurveyResponse) error
	CloseFunc                  func(ctx context.Context) error
	ProcessedMessages          []*model.RabbitMQMessage
	SavedResponses             []*model.SurveyResponse
	ShouldFail                 bool
	ErrorMsg                   string
}

// ProcessRabbitMQMessage processes a RabbitMQ message and stores it
func (m *MockResponseRepository) ProcessRabbitMQMessage(ctx context.Context, message *model.RabbitMQMessage) error {
	if m.ProcessRabbitMQMessageFunc != nil {
		return m.ProcessRabbitMQMessageFunc(ctx, message)
	}

	if m.ShouldFail {
		return fmt.Errorf(m.ErrorMsg)
	}

	m.ProcessedMessages = append(m.ProcessedMessages, message)
	return nil
}

// SaveResponse saves a survey response
func (m *MockResponseRepository) SaveResponse(ctx context.Context, response *model.SurveyResponse) error {
	if m.SaveResponseFunc != nil {
		return m.SaveResponseFunc(ctx, response)
	}

	if m.ShouldFail {
		return fmt.Errorf(m.ErrorMsg)
	}

	m.SavedResponses = append(m.SavedResponses, response)
	return nil
}

// Close closes the repository
func (m *MockResponseRepository) Close(ctx context.Context) error {
	if m.CloseFunc != nil {
		return m.CloseFunc(ctx)
	}
	return nil
}
