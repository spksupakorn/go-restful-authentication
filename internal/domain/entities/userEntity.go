package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user entity in the domain
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // Never expose password in JSON
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// TableName returns the collection name for MongoDB
func (User) CollectionName() string {
	return "users"
}

// BeforeCreate sets the timestamps before creating
func (u *User) BeforeCreate() {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}
}

// BeforeUpdate updates the UpdatedAt timestamp
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}
