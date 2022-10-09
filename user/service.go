package user

import (
	"context"
)

type Instrumentation interface {
	FailedToCreateUser(error)
	FailedToDeleteAll(error)
	FailedToFindByName(error)
	FailedToRemoveUser(error, User)
	UserCreated()
}

type Repository interface {
	CreateUser(context.Context, User) error
	All(context.Context) ([]User, error)
	FindByName(context.Context, string) (User, error)
	Delete(context.Context, User) error
	DeleteAll(context.Context) error
}

type UserService struct {
	repo            Repository
	instrumentation Instrumentation
}

func NewUserService(repo Repository, instrumentation Instrumentation) *UserService {
	return &UserService{repo, instrumentation}
}

func (service *UserService) CreateUser(name string) (User, error) {
	user := User{name}

	err := service.repo.CreateUser(context.Background(), user)
	if err != nil {
		service.instrumentation.FailedToCreateUser(err)
	}
	service.instrumentation.UserCreated()

	return user, err
}

func (service *UserService) DeleteAll() error {
	err := service.repo.DeleteAll(context.Background())

	if err != nil {
		service.instrumentation.FailedToDeleteAll(err)
	}

	return err
}

func (service *UserService) List() ([]User, error) {
	return service.repo.All(context.Background())
}

func (service *UserService) FindByName(name string) (User, error) {
	user, err := service.repo.FindByName(context.Background(), name)
	if err != nil {
		service.instrumentation.FailedToFindByName(err)
	}

	return user, err
}

func (service *UserService) Remove(user User) error {
	err := service.repo.Delete(context.Background(), user)

	if err != nil {
		service.instrumentation.FailedToRemoveUser(err, user)
	}

	return err
}
