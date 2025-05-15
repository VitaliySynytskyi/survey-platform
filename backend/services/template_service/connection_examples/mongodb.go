package connection_examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig holds the connection details
type MongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// ConnectToMongo establishes a connection to MongoDB
func ConnectToMongo(config MongoConfig) (*mongo.Client, *mongo.Database, error) {
	// Create connection URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		config.User, config.Password, config.Host, config.Port)

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Get a handle to the specified database
	db := client.Database(config.DBName)

	return client, db, nil
}

// Example usage:
func MongoExample() {
	config := MongoConfig{
		Host:     "localhost", // or "mongodb" when using Docker network
		Port:     "27017",
		User:     "mongo_admin",
		Password: "mongo_password",
		DBName:   "survey_data",
	}

	client, db, err := ConnectToMongo(config)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		// Disconnect when done
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Here you would use the database handle to access collections
	// Example:
	// collection := db.Collection("surveys")
	// result, err := collection.InsertOne(context.Background(), bson.M{"name": "Sample Survey"})

	_ = db // Placeholder to avoid unused variable warning
}
