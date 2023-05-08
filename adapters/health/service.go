package health

import (
	"app/adapters/mongo"
	"app/adapters/redis"

	"context"
)

type Service struct {
	mongoHealth mongo.MongoHealthChecker
	redisHealth redis.RedisHealthChecker
}

func NewHealthService(
	mongoHealth mongo.MongoHealthChecker,
	redisHealth redis.RedisHealthChecker,
) *Service {
	return &Service{mongoHealth, redisHealth}
}

func (s *Service) CheckHealth(ctx context.Context) Check {
	return Check{
		"mongo": NewStatus(s.mongoHealth.CheckHealth(ctx)),
		"redis": NewStatus(s.redisHealth.CheckHealth(ctx)),
	}
}
