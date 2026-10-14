package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GroupRepository struct {
	coll *mongo.Collection
}

func NewGroupRepository(db *mongo.Database) *GroupRepository {
	return &GroupRepository{coll: db.Collection("groups")}
}

func (r *GroupRepository) ExistsByID(ctx context.Context, id bson.ObjectID) (bool, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{"_id": id}, options.Count().SetLimit(1))
	return n > 0, err
}