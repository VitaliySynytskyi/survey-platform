package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_taking_service/internal/model"
)

// Producer represents a RabbitMQ producer
type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

// NewProducer creates a new RabbitMQ producer
func NewProducer(cfg config.RabbitMQConfig) (*Producer, error) {
	// Create connection string
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		cfg.Exchange, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return &Producer{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}, nil
}

// PublishResponse publishes a survey response to RabbitMQ
func (p *Producer) PublishResponse(ctx context.Context, response model.SurveyResponse) error {
	// Convert response to message
	message := model.RabbitMQMessage{
		SurveyID:     response.SurveyID,
		RespondentID: response.RespondentID,
		AnonymousID:  response.AnonymousID,
		Answers:      response.Answers,
		SubmittedAt:  response.SubmittedAt,
	}

	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message
	err = p.channel.Publish(
		p.config.Exchange,   // exchange
		p.config.RoutingKey, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			// Set message persistence
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Published response for survey %s", response.SurveyID)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

// GetConnection returns the underlying RabbitMQ connection
func (p *Producer) GetConnection() *amqp.Connection {
	return p.conn
}
