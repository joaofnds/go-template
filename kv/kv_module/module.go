package kv_module

import (
	"app/kv"
	"app/kv/kv_adapter"
	"app/kv/kv_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"kv",
		fx.Provide(kv_adapter.NewRedisStore),
		fx.Provide(func(redisStore *kv_adapter.RedisStore) kv.Store { return redisStore }),
	)

	HTTPModule = fx.Module(
		"kv http",
		Module,
		fx.Provide(kv_http.NewController, fx.Private),
		fx.Invoke(func(app *fiber.App, kvController *kv_http.Controller) {
			kvController.Register(app)
		}),
	)
)
