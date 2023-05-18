package kv

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"kv",
	fx.Provide(NewRedisStore),
	fx.Provide(func(redisStore *RedisStore) Store { return redisStore }),
	fx.Provide(NewController),
	fx.Invoke(func(app *fiber.App, controller *Controller) { controller.Register(app) }),
)

type Store interface {
	Get(string) (string, error)
	Set(string, string) error
	Del(string) error
}
