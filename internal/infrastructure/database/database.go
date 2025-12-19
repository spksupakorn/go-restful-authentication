package database

import (
	"context"
	"fmt"
	"time"

	"github.com/spksupakorn/go-restful-authentication/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDB represents the MongoDB connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	logger   *zap.Logger
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *config.Config, logger *zap.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.MongoDB.Timeout)*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Successfully connected to MongoDB",
		zap.String("database", cfg.MongoDB.Database),
	)

	return &MongoDB{
		Client:   client,
		Database: client.Database(cfg.MongoDB.Database),
		logger:   logger,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		m.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return err
	}
	m.logger.Info("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
