package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type HealthChecker struct {
	client *redis.Client
}

func NewHealthChecker(client *redis.Client) HealthChecker {
	return HealthChecker{client}
}

func (c HealthChecker) CheckHealth(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
