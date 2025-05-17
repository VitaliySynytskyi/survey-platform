package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/db"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	config     config.RabbitMQConfig
	repository *db.ResponseRepository
	deliveries <-chan amqp.Delivery
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(cfg config.RabbitMQConfig, repository *db.ResponseRepository) (*Consumer, error) {
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

	// Set QoS (prefetch)
	err = channel.Qos(cfg.PrefetchCount, 0, false)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
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

	// Declare queue
	q, err := channel.QueueDeclare(
		cfg.Queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind queue to exchange
	err = channel.QueueBind(
		q.Name,         // queue name
		cfg.RoutingKey, // routing key
		cfg.Exchange,   // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// Declare dead-letter exchange and queue for failed messages
	dlxName := fmt.Sprintf("%s.dlx", cfg.Exchange)
	dlqName := fmt.Sprintf("%s.dlq", cfg.Queue)

	err = channel.ExchangeDeclare(
		dlxName,  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLX: %w", err)
	}

	_, err = channel.QueueDeclare(
		dlqName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	err = channel.QueueBind(
		dlqName, // queue name
		"#",     // routing key (catch all)
		dlxName, // exchange
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind DLQ: %w", err)
	}

	return &Consumer{
		conn:       conn,
		channel:    channel,
		config:     cfg,
		repository: repository,
	}, nil
}

// Start starts consuming messages from the queue
func (c *Consumer) Start(ctx context.Context) error {
	// Start consuming
	deliveries, err := c.channel.Consume(
		c.config.Queue, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume from queue: %w", err)
	}

	c.deliveries = deliveries

	// Process messages
	go c.processMessages(ctx)

	log.Printf("Started consuming from queue %s", c.config.Queue)
	return nil
}

// processMessages processes messages from the queue
func (c *Consumer) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Context done, stopping consumer")
			return
		case delivery, ok := <-c.deliveries:
			if !ok {
				log.Println("Channel closed, stopping consumer")
				return
			}

			// Process delivery
			c.handleDelivery(ctx, delivery)
		}
	}
}

// validateMessage validates the RabbitMQ message structure
func validateMessage(message *model.RabbitMQMessage) error {
	// Validate SurveyID (must be a valid ObjectID)
	if message.SurveyID == "" {
		return fmt.Errorf("survey_id is required")
	}

	// Check if SurveyID is a valid ObjectID
	_, err := primitive.ObjectIDFromHex(message.SurveyID)
	if err != nil {
		return fmt.Errorf("invalid survey_id format: %w", err)
	}

	// Validate that we have either a respondent ID or an anonymous ID
	if message.RespondentID == "" && message.AnonymousID == "" {
		return fmt.Errorf("either respondent_id or anonymous_id must be provided")
	}

	// Validate answers
	if len(message.Answers) == 0 {
		return fmt.Errorf("at least one answer must be provided")
	}

	// Validate each answer has a questionID and value
	for i, answer := range message.Answers {
		if answer.QuestionID == "" {
			return fmt.Errorf("answer at index %d is missing question_id", i)
		}

		// Temporarily commented out to allow non-ObjectID question_ids
		/*
			_, err := primitive.ObjectIDFromHex(answer.QuestionID)
			if err != nil {
				return fmt.Errorf("invalid question_id format for answer at index %d: %w", i, err)
			}
		*/

		if answer.Value == nil {
			return fmt.Errorf("answer at index %d is missing value", i)
		}
	}

	// Validate submitted_at timestamp
	if message.SubmittedAt == "" {
		return fmt.Errorf("submitted_at timestamp is required")
	}

	return nil
}

// handleDelivery processes a single delivery
func (c *Consumer) handleDelivery(ctx context.Context, delivery amqp.Delivery) {
	// Parse message
	var message model.RabbitMQMessage
	err := json.Unmarshal(delivery.Body, &message)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		// Reject and don't requeue - this is a malformed message
		delivery.Reject(false)
		return
	}

	// Validate message
	err = validateMessage(&message)
	if err != nil {
		log.Printf("Message validation failed: %v", err)
		// Reject and don't requeue - send to DLQ
		delivery.Reject(false)
		return
	}

	// Process message
	err = c.repository.ProcessRabbitMQMessage(ctx, &message)
	if err != nil {
		log.Printf("Failed to process message: %v", err)
		// Reject and requeue based on error type
		// For non-recoverable errors (like invalid data), don't requeue
		delivery.Reject(true)
		return
	}

	// Acknowledge successful processing
	err = delivery.Ack(false)
	if err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// CheckHealth verifies the RabbitMQ connection is healthy
func (c *Consumer) CheckHealth() error {
	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed or nil")
	}

	if c.channel == nil {
		return fmt.Errorf("RabbitMQ channel is nil")
	}

	return nil
}
