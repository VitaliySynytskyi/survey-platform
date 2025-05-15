package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/response_processor_service/internal/model"
)

// ResponseRepository handles operations on survey responses in MongoDB
type ResponseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewResponseRepository creates a new response repository
func NewResponseRepository(cfg config.MongoDBConfig) (*ResponseRepository, error) {
	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get collection
	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	return &ResponseRepository{
		client:     client,
		collection: collection,
	}, nil
}

// SaveResponse saves a survey response to MongoDB
func (r *ResponseRepository) SaveResponse(ctx context.Context, response *model.SurveyResponse) error {
	// Insert response
	_, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return fmt.Errorf("failed to insert response: %w", err)
	}

	log.Printf("Saved response for survey %s", response.SurveyID.Hex())
	return nil
}

// ProcessRabbitMQMessage processes a message from RabbitMQ and saves it to MongoDB
func (r *ResponseRepository) ProcessRabbitMQMessage(ctx context.Context, message *model.RabbitMQMessage) error {
	// Convert message to survey response
	surveyID, err := primitive.ObjectIDFromHex(message.SurveyID)
	if err != nil {
		return fmt.Errorf("invalid survey ID: %w", err)
	}

	// Parse submitted time
	submittedAt, err := time.Parse(time.RFC3339, message.SubmittedAt)
	if err != nil {
		// Use current time if parsing fails
		submittedAt = time.Now().UTC()
	}

	// Create survey response
	response := &model.SurveyResponse{
		ID:           primitive.NewObjectID(),
		SurveyID:     surveyID,
		RespondentID: message.RespondentID,
		AnonymousID:  message.AnonymousID,
		Answers:      message.Answers,
		SubmittedAt:  submittedAt,
	}

	// Save to MongoDB
	return r.SaveResponse(ctx, response)
}

// Close closes the MongoDB connection
func (r *ResponseRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
