package health

import (
	"app/adapters/health"
	"context"

	"go.uber.org/fx"
)

var UnhealthyHealthService = fx.Decorate(NewUnhealthyHealthService)

func NewUnhealthyHealthService() health.Checker {
	return &unhealthyHealthService{}
}

type unhealthyHealthService struct{}

func (c *unhealthyHealthService) CheckHealth(_ context.Context) health.Check {
	return health.Check{
		"mongo": health.Status{Status: health.StatusDown},
		"redis": health.Status{Status: health.StatusDown},
	}
}
