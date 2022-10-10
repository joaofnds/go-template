package health

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthChecker interface {
	CheckHealth(ctx context.Context) HealthCheck
}

type HealthService struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
}

func NewHealthService(mongoClient *mongo.Client, redisClient *redis.Client) *HealthService {
	return &HealthService{mongoClient, redisClient}
}

func (c *HealthService) CheckHealth(ctx context.Context) HealthCheck {
	return HealthCheck{
		Mongo: c.MongoHealth(ctx),
		Redis: c.RedisHealth(ctx),
	}
}

func (c *HealthService) MongoHealth(ctx context.Context) Status {
	if err := c.mongoClient.Ping(ctx, nil); err != nil {
		return Status{Status: StatusDown}
	}
	return Status{Status: StatusUp}
}

func (c *HealthService) RedisHealth(ctx context.Context) Status {
	if cmd := c.redisClient.Ping(ctx); cmd.Err() != nil {
		return Status{Status: StatusDown}
	}
	return Status{Status: StatusUp}
}
