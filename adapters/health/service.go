package health

import (
	"app/adapters/mongo"
	"app/adapters/redis"

	"context"
)

type Service struct {
	mongoHealth mongo.HealthChecker
	redisHealth redis.HealthChecker
}

func NewHealthService(
	mongoHealth mongo.HealthChecker,
	redisHealth redis.HealthChecker,
) *Service {
	return &Service{mongoHealth, redisHealth}
}

func (s *Service) CheckHealth(ctx context.Context) Check {
	return Check{
		"mongo": NewStatus(s.mongoHealth.CheckHealth(ctx)),
		"redis": NewStatus(s.redisHealth.CheckHealth(ctx)),
	}
}
