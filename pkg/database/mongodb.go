package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient represents a MongoDB client
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoClient creates a new MongoDB client
func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Get database
	db := client.Database(dbName)

	return &MongoClient{
		Client: client,
		DB:     db,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// Collection returns a MongoDB collection
func (m *MongoClient) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
