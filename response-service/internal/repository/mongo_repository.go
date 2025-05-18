package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/survey-app/response-service/internal/config" // Assuming config path
	"github.com/survey-app/response-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoRepository implements ResponseRepositoryInterface for MongoDB
type MongoRepository struct {
	client     *mongo.Client
	dbName     string
	collection *mongo.Collection
}

// NewMongoRepository creates a new MongoRepository instance and connects to MongoDB
func NewMongoRepository(cfg *config.Config) (*MongoRepository, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB!")

	collection := client.Database(cfg.MongoDBName).Collection(cfg.ResponseCollection)

	return &MongoRepository{
		client:     client,
		dbName:     cfg.MongoDBName,
		collection: collection,
	}, nil
}

// CreateResponse inserts a new response document into MongoDB
func (r *MongoRepository) CreateResponse(ctx context.Context, response *models.Response) error {
	response.SubmittedAt = time.Now() // Ensure submitted time is set
	_, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return fmt.Errorf("failed to insert response: %w", err)
	}
	return nil
}

// GetResponsesBySurveyID retrieves all responses for a given surveyID
func (r *MongoRepository) GetResponsesBySurveyID(ctx context.Context, surveyID int) ([]*models.Response, error) {
	filter := bson.M{"surveyId": surveyID}
	findOptions := options.Find()
	// Add sorting if needed, e.g., by SubmittedAt
	// findOptions.SetSort(bson.D{{"submittedAt", -1}} // -1 for descending

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find responses for surveyID %d: %w", surveyID, err)
	}
	defer cursor.Close(ctx)

	var responses []*models.Response
	if err = cursor.All(ctx, &responses); err != nil {
		return nil, fmt.Errorf("failed to decode responses for surveyID %d: %w", surveyID, err)
	}

	// If no documents found, cursor.All returns an empty slice and no error.
	// MongoDB driver doesn't typically return a "not found" error for Find operations that yield no results.
	return responses, nil
}

// Disconnect closes the MongoDB client connection
func (r *MongoRepository) Disconnect(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}
