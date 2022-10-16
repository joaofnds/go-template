package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"redis",
	fx.Provide(NewClient),
	fx.Invoke(HookRedis),
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func HookRedis(lifecycle fx.Lifecycle, redis *redis.Client) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return redis.Ping(ctx).Err()
		},
	})
}