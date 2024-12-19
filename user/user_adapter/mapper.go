package user_adapter

import (
	"app/user"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UserMongoMapper struct{}

func NewUserMongoMapper() *UserMongoMapper {
	return &UserMongoMapper{}
}

func (mapper *UserMongoMapper) FromDocument(doc bson.M) (user.User, error) {
	id, ok := doc["_id"].(string)
	if !ok {
		return user.User{}, fmt.Errorf("invalid id: %v", doc["_id"])
	}

	email, ok := doc["email"].(string)
	if !ok {
		return user.User{}, fmt.Errorf("invalid email: %v", doc["email"])
	}

	createdAtString, ok := doc["created_at"].(string)
	if !ok {
		return user.User{}, fmt.Errorf("invalid created_at: %v", doc["created_at"])
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtString)
	if err != nil {
		return user.User{}, fmt.Errorf("invalid created_at: %v", doc["created_at"])
	}

	return user.User{
		ID:        id,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}

func (mapper *UserMongoMapper) ToDocument(user user.User) bson.M {
	return bson.M{
		"_id":        user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt.UTC().Format(time.RFC3339),
	}
}
