package db

import (
	"context"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// Repository defines the interface for a response repository
type Repository interface {
	ProcessRabbitMQMessage(ctx context.Context, message *model.RabbitMQMessage) error
	SaveResponse(ctx context.Context, response *model.SurveyResponse) error
	Close(ctx context.Context) error
	CheckHealth(ctx context.Context) error
}
