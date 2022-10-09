package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) Repository {
	return &MongoRepository{client.Database("template").Collection("users")}
}

func (repo *MongoRepository) CreateUser(ctx context.Context, user User) error {
	_, err := repo.collection.InsertOne(ctx, user)
	return err
}

func (repo *MongoRepository) FindByName(ctx context.Context, name string) (User, error) {
	var user User
	result := repo.collection.FindOne(context.Background(), bson.M{"name": name})
	err := result.Decode(&user)

	return user, err
}

func (repo *MongoRepository) Delete(ctx context.Context, user User) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"name": user.Name})
	return err
}

func (repo *MongoRepository) DeleteAll(ctx context.Context) error {
	_, err := repo.collection.DeleteMany(ctx, bson.M{})
	return err
}

func (repo *MongoRepository) All(ctx context.Context) ([]User, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
