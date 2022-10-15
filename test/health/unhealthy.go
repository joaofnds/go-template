package health

import (
	"app/health"
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
		Mongo: health.Status{Status: health.StatusDown},
		Redis: health.Status{Status: health.StatusDown},
	}
}
