package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var createGroups = Migration{
	Version: 2,
	Name:    "create_groups",
	Up: func(ctx context.Context, db *mongo.Database) error {
		err := createCollection(ctx, db, "groups", bson.M{
			"bsonType":             "object",
			"required":             []string{"name", "email", "created_at", "updated_at"},
			"additionalProperties": false,
			"properties": bson.M{
				"_id":        bson.M{"bsonType": "objectId"},
				"name":       bson.M{"bsonType": "string", "minLength": 1},
				"email":      bson.M{"bsonType": "string", "pattern": emailPattern},
				"created_at": bson.M{"bsonType": "date"},
				"updated_at": bson.M{"bsonType": "date"},
			},
		})
		if err != nil {
			return err
		}

		_, err = db.Collection("groups").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		return err
	},
	Down: func(ctx context.Context, db *mongo.Database) error {
		return db.Collection("groups").Drop(ctx)
	},
}