package user

import (
	"context"
)

type Service struct {
	repo    Repository
	emitter Emitter
}

func NewUserService(repo Repository, emitter Emitter) *Service {
	return &Service{repo: repo, emitter: emitter}
}

func (service *Service) CreateUser(ctx context.Context, name string) (User, error) {
	user := User{name}

	err := service.repo.CreateUser(ctx, user)
	if err != nil {
		service.emitter.FailedToCreateUser(err)
	} else {
		service.emitter.UserCreated(user)
	}

	return user, err
}

func (service *Service) DeleteAll(ctx context.Context) error {
	err := service.repo.DeleteAll(ctx)

	if err != nil {
		service.emitter.FailedToDeleteAll(err)
	}

	return err
}

func (service *Service) List(ctx context.Context) ([]User, error) {
	return service.repo.All(ctx)
}

func (service *Service) FindByName(ctx context.Context, name string) (User, error) {
	user, err := service.repo.FindByName(ctx, name)
	if err != nil {
		service.emitter.FailedToFindByName(err)
	}

	return user, err
}

func (service *Service) Remove(ctx context.Context, user User) error {
	err := service.repo.Delete(ctx, user)

	if err != nil {
		service.emitter.FailedToRemoveUser(err, user)
	}

	return err
}
