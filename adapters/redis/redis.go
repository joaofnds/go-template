package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"redis",
	fx.Provide(NewClient),
	fx.Provide(NewHealthChecker),
	fx.Invoke(HookRedis),
)

func NewClient(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})
}

func HookRedis(lifecycle fx.Lifecycle, redis *redis.Client) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return redis.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			return redis.Close()
		},
	})
}
