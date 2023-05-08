package http

import (
	"app/adapters/health"
	"app/kv"
	"app/user"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(
		app *fiber.App,
		healthController *health.Controller,
		userController *user.Controller,
		kvController *kv.Controller,
	) {
		healthController.Register(app)
		userController.Register(app)
		kvController.Register(app)
	}),
)
