package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type HealthChecker struct {
	client *redis.Client
}

func NewHealthChecker(client *redis.Client) HealthChecker {
	return HealthChecker{client}
}

func (healthChecker HealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx).Err()
}
