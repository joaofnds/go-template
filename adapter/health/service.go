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

func (s *Service) CheckHealth(ctx context.Context) Check {
	return Check{
		"postgres": NewStatus(s.postgresHealth.CheckHealth(ctx)),
		"redis":    NewStatus(s.redisHealth.CheckHealth(ctx)),
	}
}
