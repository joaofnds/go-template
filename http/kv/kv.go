package kv

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewController),
	fx.Invoke(func(app *fiber.App, controller *Controller) {
		controller.Register(app)
	}),
)
