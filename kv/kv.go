package kv

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"kv",
	fx.Provide(NewRedisStore),
	fx.Provide(func(redisStore *RedisStore) Store { return redisStore }),
	fx.Provide(NewController),
)

type Store interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Del(context.Context, string) error
}
