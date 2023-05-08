package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisHealthChecker struct {
	client *redis.Client
}

func NewRedisHealthChecker(client *redis.Client) RedisHealthChecker {
	return RedisHealthChecker{client}
}

func (c RedisHealthChecker) CheckHealth(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
