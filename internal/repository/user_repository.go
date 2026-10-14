package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"go-starter/internal/model"
)

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phone string) (*model.User, error) {
	return r.findOne(ctx, bson.M{"phone_number": phone})
}

func (r *UserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return r.findOne(ctx, bson.M{"_id": id})
}

func (r *UserRepository) FindAll(ctx context.Context, groupID *bson.ObjectID) ([]model.User, error) {
	filter := bson.M{}
	if groupID != nil {
		filter["group_id"] = *groupID
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	users := []model.User{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicate
	}
	return err
}

func (r *UserRepository) Update(ctx context.Context, id bson.ObjectID, set, unset bson.M) (*model.User, error) {
	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var user model.User
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&user)
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return nil, ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return nil, ErrDuplicate
	case err != nil:
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserRepository) findOne(ctx context.Context, filter bson.M) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}