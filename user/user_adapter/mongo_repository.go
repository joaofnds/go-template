package user_adapter

import (
	"app/user"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ user.Repository = &MongoRepository{}

type MongoRepository struct {
	collection *mongo.Collection
	mapper     *UserMongoMapper
}

func NewMongoRepository(client *mongo.Client, mapper *UserMongoMapper) *MongoRepository {
	return &MongoRepository{
		collection: client.Database("template").Collection("users"),
		mapper:     mapper,
	}
}

func (repository *MongoRepository) CreateUser(ctx context.Context, newUser user.User) error {
	_, err := repository.collection.InsertOne(ctx, repository.mapper.ToDocument(newUser))
	return translateErr(err)
}

func (repository *MongoRepository) FindByID(ctx context.Context, id string) (user.User, error) {
	return repository.findBy(ctx, bson.M{"_id": id})
}

func (repository *MongoRepository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	return repository.findBy(ctx, bson.M{"email": email})
}

func (repository *MongoRepository) Delete(ctx context.Context, userToDelete user.User) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"email": userToDelete.Email})
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
	defer func() { _ = cursor.Close(ctx) }()

	var users []user.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, translateErr(err)
	}

	return users, nil
}

func (repository *MongoRepository) findBy(ctx context.Context, params bson.M) (user.User, error) {
	result := repository.collection.FindOne(ctx, params)
	if result.Err() != nil {
		return user.User{}, translateErr(result.Err())
	}

	var doc bson.M
	if err := result.Decode(&doc); err != nil {
		return user.User{}, translateErr(err)
	}

	return repository.mapper.FromDocument(doc)
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
