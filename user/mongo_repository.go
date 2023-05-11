package user

import (
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

func (repo *MongoRepository) CreateUser(ctx context.Context, user User) error {
	_, err := repo.collection.InsertOne(ctx, user)
	return translateErr(err)
}

func (repo *MongoRepository) FindByName(ctx context.Context, name string) (User, error) {
	var user User
	result := repo.collection.FindOne(context.Background(), bson.M{"name": name})
	return user, translateErr(result.Decode(&user))
}

func (repo *MongoRepository) Delete(ctx context.Context, user User) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"name": user.Name})
	return translateErr(err)
}

func (repo *MongoRepository) DeleteAll(ctx context.Context) error {
	_, err := repo.collection.DeleteMany(ctx, bson.M{})
	return translateErr(err)
}

func (repo *MongoRepository) All(ctx context.Context) ([]User, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, translateErr(err)
	}
	defer cursor.Close(ctx)

	var users []User
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
		return ErrNotFound
	default:
		return ErrRepository
	}
}
