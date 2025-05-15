package db

import (
	"context"
	"fmt"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// MockResponseRepository is a mock implementation of the ResponseRepository for testing
type MockResponseRepository struct {
	ProcessRabbitMQMessageFunc func(ctx context.Context, message *model.RabbitMQMessage) error
	ProcessedMessages          []*model.RabbitMQMessage
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
