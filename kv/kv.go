package kv

import "go.uber.org/fx"

var Module = fx.Module(
	"kv",
	fx.Provide(NewRedisStore),
	fx.Provide(func(redisStore *RedisStore) Store { return redisStore }),
	fx.Provide(NewController),
)
