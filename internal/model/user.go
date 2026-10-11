package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Email       *string       `bson:"email,omitempty"`
	PhoneNumber string        `bson:"phone_number"`
	Password    string        `bson:"password"`
	RoleID      bson.ObjectID `bson:"role_id"`
	GroupID     bson.ObjectID `bson:"group_id"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}