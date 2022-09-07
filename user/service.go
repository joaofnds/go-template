package user

import (
	"context"
)

type UserService struct {
	repo            *UserRepository
	instrumentation Instrumentation
}

func NewUserService(repo *UserRepository, instrumentation Instrumentation) *UserService {
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
