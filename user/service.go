package user

import (
	"app/user/queue"

	"context"
)

type Service struct {
	repo    Repository
	greeter *queue.Greeter
	probe   Probe
}

func NewUserService(repo Repository, greeter *queue.Greeter, probe Probe) *Service {
	return &Service{repo: repo, greeter: greeter, probe: probe}
}

func (service *Service) CreateUser(ctx context.Context, name string) (User, error) {
	user := User{name}

	err := service.repo.CreateUser(ctx, user)
	if err != nil {
		service.probe.FailedToCreateUser(err)
	}
	service.probe.UserCreated()

	if err := service.greeter.Enqueue(name); err != nil {
		service.probe.FailedToEnqueue(err)
	}

	return user, err
}

func (service *Service) DeleteAll(ctx context.Context) error {
	err := service.repo.DeleteAll(ctx)

	if err != nil {
		service.probe.FailedToDeleteAll(err)
	}

	return err
}

func (service *Service) List(ctx context.Context) ([]User, error) {
	return service.repo.All(ctx)
}

func (service *Service) FindByName(ctx context.Context, name string) (User, error) {
	user, err := service.repo.FindByName(ctx, name)
	if err != nil {
		service.probe.FailedToFindByName(err)
	}

	return user, err
}

func (service *Service) Remove(ctx context.Context, user User) error {
	err := service.repo.Delete(ctx, user)

	if err != nil {
		service.probe.FailedToRemoveUser(err, user)
	}

	return err
}
