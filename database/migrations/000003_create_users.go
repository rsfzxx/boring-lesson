package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var createUsers = Migration{
	Version: 3,
	Name:    "create_users",
	Up: func(ctx context.Context, db *mongo.Database) error {
		err := createCollection(ctx, db, "users", bson.M{
			"bsonType": "object",
			"required": []string{
				"name", "phone_number", "password", "role_id", "group_id", "created_at", "updated_at",
			},
			"additionalProperties": false,
			"properties": bson.M{
				"_id":          bson.M{"bsonType": "objectId"},
				"name":         bson.M{"bsonType": "string", "minLength": 1},
				"email":        bson.M{"bsonType": "string", "pattern": emailPattern}, // opsional
				"phone_number": bson.M{"bsonType": "string", "pattern": `^\+?[0-9]{8,15}$`},
				"password":     bson.M{"bsonType": "string", "minLength": 20}, // hash bcrypt
				"role_id":      bson.M{"bsonType": "objectId"},
				"group_id":     bson.M{"bsonType": "objectId"},
				"created_at":   bson.M{"bsonType": "date"},
				"updated_at":   bson.M{"bsonType": "date"},
			},
		})
		if err != nil {
			return err
		}

		_, err = db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "phone_number", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).
					SetPartialFilterExpression(bson.M{"email": bson.M{"$type": "string"}}),
			},
			{Keys: bson.D{{Key: "group_id", Value: 1}}},
			{Keys: bson.D{{Key: "role_id", Value: 1}}},
		})
		return err
	},
	Down: func(ctx context.Context, db *mongo.Database) error {
		return db.Collection("users").Drop(ctx)
	},
}