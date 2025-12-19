package repositories

import (
	"context"

	"github.com/spksupakorn/go-restful-authentication/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data access (Port)
// This abstraction allows us to switch database implementations easily
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// GetAll retrieves all users with pagination
	GetAll(ctx context.Context, skip, limit int64) ([]*entities.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// EmailExists checks if an email already exists
	EmailExists(ctx context.Context, email string) (bool, error)
}
