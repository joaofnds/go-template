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

func (service *Service) CreateUser(name string) (User, error) {
	user := User{name}

	err := service.repo.CreateUser(context.Background(), user)
	if err != nil {
		service.probe.FailedToCreateUser(err)
	}
	service.probe.UserCreated()

	if err := service.greeter.Enqueue(name); err != nil {
		service.probe.FailedToEnqueue(err)
	}

	return user, err
}

func (service *Service) DeleteAll() error {
	err := service.repo.DeleteAll(context.Background())

	if err != nil {
		service.probe.FailedToDeleteAll(err)
	}

	return err
}

func (service *Service) List() ([]User, error) {
	return service.repo.All(context.Background())
}

func (service *Service) FindByName(name string) (User, error) {
	user, err := service.repo.FindByName(context.Background(), name)
	if err != nil {
		service.probe.FailedToFindByName(err)
	}

	return user, err
}

func (service *Service) Remove(user User) error {
	err := service.repo.Delete(context.Background(), user)

	if err != nil {
		service.probe.FailedToRemoveUser(err, user)
	}

	return err
}
