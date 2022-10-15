package http

import (
	"web/health"
	webfiber "web/http/fiber"
	"web/kv"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	webfiber.Module,
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
