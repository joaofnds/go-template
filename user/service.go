package user

import (
	"context"

	"go.uber.org/zap"
)

type UserService struct {
	repo   *UserRepository
	logger *zap.Logger
}

func NewUserService(repo *UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo, logger}
}

func (service *UserService) CreateUser(name string) (User, error) {
	user := User{name}
	err := service.repo.CreateUser(context.Background(), user)

	if err != nil {
		service.logger.Error("failed to create user", zap.Error(err))
	}

	return user, err
}

func (service *UserService) DeleteAll() error {
	err := service.repo.DeleteAll(context.Background())

	if err != nil {
		service.logger.Error("failed to delete all", zap.Error(err))
	}

	return err
}

func (service *UserService) List() ([]User, error) {
	return service.repo.All(context.Background())
}

func (service *UserService) FindByName(name string) (User, bool) {
	user, err := service.repo.FindByName(context.Background(), name)

	if err != nil {
		service.logger.Error("failed to find user by name", zap.Error(err))
		return user, false
	}

	return user, true
}

func (service *UserService) Remove(user User) error {
	err := service.repo.Delete(context.Background(), user)

	if err != nil {
		service.logger.Error("failed to remove user", zap.Error(err), zap.String("name", user.Name))
	}

	return err
}
