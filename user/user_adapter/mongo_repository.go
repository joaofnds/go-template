package user_adapter

import (
	"app/user"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{client.Database("template").Collection("users")}
}

func (repository *MongoRepository) CreateUser(ctx context.Context, newUser user.User) error {
	_, err := repository.collection.InsertOne(ctx, newUser)
	return translateErr(err)
}

func (repository *MongoRepository) FindByName(ctx context.Context, name string) (user.User, error) {
	var userFound user.User
	result := repository.collection.FindOne(ctx, bson.M{"name": name})
	return userFound, translateErr(result.Decode(&userFound))
}

func (repository *MongoRepository) Delete(ctx context.Context, userToDelete user.User) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"name": userToDelete.Name})
	return translateErr(err)
}

func (repository *MongoRepository) DeleteAll(ctx context.Context) error {
	_, err := repository.collection.DeleteMany(ctx, bson.M{})
	return translateErr(err)
}

func (repository *MongoRepository) All(ctx context.Context) ([]user.User, error) {
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, translateErr(err)
	}
	defer cursor.Close(ctx)

	var users []user.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, translateErr(err)
	}

	return users, nil
}

func translateErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, mongo.ErrNoDocuments):
		return user.ErrNotFound
	default:
		return user.ErrRepository
	}
}
