package health

import (
	"context"
)

func NewFakeHealthService(real Checker) *FakeHealthService {
	return &FakeHealthService{
		realChecker:  real,
		unhealthy:    UnhealthyHealthService{},
		useUnhealthy: false,
	}
}

type FakeHealthService struct {
	realChecker  Checker
	unhealthy    Checker
	useUnhealthy bool
}

func (service *FakeHealthService) UseUnhealthy() {
	service.useUnhealthy = true
}

func (service *FakeHealthService) CheckHealth(ctx context.Context) Check {
	if service.useUnhealthy {
		return service.unhealthy.CheckHealth(ctx)
	}

	return service.realChecker.CheckHealth(ctx)
}

type UnhealthyHealthService struct{}

func (service UnhealthyHealthService) CheckHealth(context.Context) Check {
	return Check{
		"postgres": Status{Status: StatusDown},
		"redis":    Status{Status: StatusDown},
	}
}
