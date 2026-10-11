package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go-starter/internal/model"
)

var ErrNotFound = errors.New("data tidak ditemukan")

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, bson.M{"phone_number": phone}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}