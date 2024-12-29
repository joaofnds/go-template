package kv_module

import (
	"app/kv"
	"app/kv/kv_adapter"
	"app/kv/kv_http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"kv",
	fx.Provide(kv_adapter.NewRedisStore),
	fx.Provide(func(redisStore *kv_adapter.RedisStore) kv.Store { return redisStore }),
	fx.Provide(kv_http.NewController),
)
