package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"go-starter/internal/model"
)

type RoleRepository struct {
	coll *mongo.Collection
}

func NewRoleRepository(db *mongo.Database) *RoleRepository {
	return &RoleRepository{coll: db.Collection("roles")}
}

func (r *RoleRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Role, error) {
	var role model.Role
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByNames(ctx context.Context, names []string) ([]model.Role, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := r.coll.Find(ctx, bson.M{"name": bson.M{"$in": names}}, opts)
	if err != nil {
		return nil, err
	}
	roles := []model.Role{}
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}