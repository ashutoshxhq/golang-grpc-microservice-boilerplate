package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User eporting model for joke
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Username  string             `json:"username" bson:"username,omitempty"`
	Phone     string             `json:"phone" bson:"phone,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	Role      string             `json:"role" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}
