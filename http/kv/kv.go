package kv

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewKVController),
	fx.Invoke(func(app *fiber.App, controller *KVController) {
		controller.Register(app)
	}),
)
