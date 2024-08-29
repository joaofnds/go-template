package user

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/internal/clock"
	"app/internal/id"
	"context"
)

type Service struct {
	id      id.Generator
	clock   clock.Clock
	repo    Repository
	emitter Emitter
}

func NewUserService(
	id *uuid.Generator,
	clock *time.Clock,
	repo Repository,
	emitter Emitter,
) *Service {
	return &Service{
		id:      id,
		clock:   clock,
		repo:    repo,
		emitter: emitter,
	}
}

func (service *Service) CreateUser(ctx context.Context, name string) (User, error) {
	user := User{
		ID:        service.id.NewID(),
		Name:      name,
		CreatedAt: service.clock.Now(),
	}

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
