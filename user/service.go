package user

import (
	"app/adapter/time"
	"app/adapter/uuid"
	"app/internal/clock"
	"app/internal/id"
	"context"
)

type Service struct {
	id          id.Generator
	clock       clock.Clock
	repo        Repository
	emitter     Emitter
	permissions *PermissionService
}

func NewUserService(
	id *uuid.Generator,
	clock *time.Clock,
	repo Repository,
	emitter Emitter,
	permissions *PermissionService,
) *Service {
	return &Service{
		id:          id,
		clock:       clock,
		repo:        repo,
		emitter:     emitter,
		permissions: permissions,
	}
}

func (service *Service) CreateUser(ctx context.Context, email string) (User, error) {
	user := User{
		ID:        service.id.NewID(),
		Email:     email,
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),
	}

	err := service.repo.CreateUser(ctx, user)
	if err != nil {
		_ = service.emitter.FailedToCreateUser(ctx, err)
	} else {
		_ = service.emitter.UserCreated(ctx, user)
	}

	return user, service.permissions.GrantNewUserPermission(user)
}

func (service *Service) DeleteAll(ctx context.Context) error {
	err := service.repo.DeleteAll(ctx)

	if err != nil {
		_ = service.emitter.FailedToDeleteAll(ctx, err)
	}

	return err
}

func (service *Service) List(ctx context.Context) ([]User, error) {
	return service.repo.All(ctx)
}

func (service *Service) FindByID(ctx context.Context, id string) (User, error) {
	user, err := service.repo.FindByID(ctx, id)
	if err != nil {
		_ = service.emitter.FailedToFindByID(ctx, err)
	}

	return user, err
}

func (service *Service) FindByEmail(ctx context.Context, email string) (User, error) {
	user, err := service.repo.FindByEmail(ctx, email)
	if err != nil {
		_ = service.emitter.FailedToFindByName(ctx, err)
	}

	return user, err
}

func (service *Service) Remove(ctx context.Context, user User) error {
	err := service.repo.Delete(ctx, user)

	if err != nil {
		_ = service.emitter.FailedToRemoveUser(ctx, err, user)
	} else {
		_ = service.emitter.UserRemoved(ctx, user)
	}

	return err
}
