package health

import (
	"context"
	"web/health"

	"go.uber.org/fx"
)

var UnhealthyHealthService = fx.Decorate(NewUnhealthyHealthService)

func NewUnhealthyHealthService() health.HealthChecker {
	return &unhealthyHealthService{}
}

type unhealthyHealthService struct{}

func (c *unhealthyHealthService) CheckHealth(_ context.Context) health.HealthCheck {
	return health.HealthCheck{
		Mongo: health.Status{Status: health.StatusDown},
		Redis: health.Status{Status: health.StatusDown},
	}
}
