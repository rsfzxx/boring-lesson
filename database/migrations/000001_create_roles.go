package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var createRoles = Migration{
	Version: 1,
	Name:    "create_roles",
	Up: func(ctx context.Context, db *mongo.Database) error {
		names := []string{"Admin", "Leader", "Member"}

		err := createCollection(ctx, db, "roles", bson.M{
			"bsonType":             "object",
			"required":             []string{"name"},
			"additionalProperties": false,
			"properties": bson.M{
				"_id":  bson.M{"bsonType": "objectId"},
				"name": bson.M{"enum": names},
			},
		})
		if err != nil {
			return err
		}

		roles := db.Collection("roles")
		_, err = roles.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			return err
		}

		for _, name := range names {
			_, err := roles.UpdateOne(ctx,
				bson.M{"name": name},
				bson.M{"$setOnInsert": bson.M{"name": name}},
				options.UpdateOne().SetUpsert(true),
			)
			if err != nil {
				return err
			}
		}
		return nil
	},
	Down: func(ctx context.Context, db *mongo.Database) error {
		return db.Collection("roles").Drop(ctx)
	},
}