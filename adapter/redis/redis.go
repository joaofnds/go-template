package redis

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("redis", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewClient),
		fx.Provide(NewHealthChecker),
	)
	Invokes = fx.Options(
		fx.Invoke(EnableTracing),
		fx.Invoke(HookRedis),
	)
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

func EnableTracing(redis *redis.Client) error {
	return redisotel.InstrumentTracing(redis)
}
