package repository

import (
	"context"
	"errors"

	"github.com/shuyou-ai/shuyou-go/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrDuplicateKey = errors.New("duplicate key")

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{coll: db.Collection(model.UserCollectionName)}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.PrepareCreate()

	if _, err := r.coll.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateKey
		}
		return err
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.coll.FindOne(ctx, activeFilter(bson.M{"_id": id})).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.coll.FindOne(ctx, activeFilter(bson.M{"username": username})).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func activeFilter(filter bson.M) bson.M {
	filter["$or"] = []bson.M{
		{"deleted_at": bson.M{"$exists": false}},
		{"deleted_at": nil},
	}
	return filter
}
