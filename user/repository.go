package user

import "context"

type Repository interface {
	CreateUser(context.Context, User) error
	All(context.Context) ([]User, error)
	FindByID(context.Context, string) (User, error)
	FindByEmail(context.Context, string) (User, error)
	Delete(context.Context, User) error
	DeleteAll(context.Context) error
}
