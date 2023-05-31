package health

import (
	"app/adapter/postgres"
	"app/adapter/redis"

	"context"
)

type Service struct {
	postgresHealth postgres.HealthChecker
	redisHealth    redis.HealthChecker
}

func NewHealthService(
	postgresHealth postgres.HealthChecker,
	redisHealth redis.HealthChecker,
) *Service {
	return &Service{postgresHealth: postgresHealth, redisHealth: redisHealth}
}

func (service *Service) CheckHealth(ctx context.Context) Check {
	return Check{
		"postgres": NewStatus(service.postgresHealth.CheckHealth(ctx)),
		"redis":    NewStatus(service.redisHealth.CheckHealth(ctx)),
	}
}
