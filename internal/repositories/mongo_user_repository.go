package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/spksupakorn/go-restful-authentication/internal/domain/entities"
	"github.com/spksupakorn/go-restful-authentication/internal/domain/repositories"
	"github.com/spksupakorn/go-restful-authentication/internal/infrastructure/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// mongoUserRepository implements the UserRepository interface using MongoDB
type mongoUserRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewMongoUserRepository creates a new instance of MongoDB user repository
func NewMongoUserRepository(db *database.MongoDB, logger *zap.Logger) repositories.UserRepository {
	collection := db.GetCollection(entities.User{}.CollectionName())

	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logger.Error("Failed to create unique index on email", zap.Error(err))
	}

	return &mongoUserRepository{
		collection: collection,
		logger:     logger,
	}
}

// Create creates a new user
func (r *mongoUserRepository) Create(ctx context.Context, user *entities.User) error {
	user.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		}
		r.logger.Error("Failed to create user", zap.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a user by ID
func (r *mongoUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error("Failed to get user by ID", zap.Error(err), zap.String("id", id.Hex()))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error("Failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetAll retrieves all users with pagination
func (r *mongoUserRepository) GetAll(ctx context.Context, skip, limit int64) ([]*entities.User, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		r.logger.Error("Failed to get all users", zap.Error(err))
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		r.logger.Error("Failed to decode users", zap.Error(err))
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// Update updates a user
func (r *mongoUserRepository) Update(ctx context.Context, user *entities.User) error {
	user.BeforeUpdate()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		}
		r.logger.Error("Failed to update user", zap.Error(err), zap.String("id", user.ID.Hex()))
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *mongoUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete user", zap.Error(err), zap.String("id", id.Hex()))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Count returns the total number of users
func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// EmailExists checks if an email already exists
func (r *mongoUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		r.logger.Error("Failed to check email existence", zap.Error(err), zap.String("email", email))
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return count > 0, nil
}
