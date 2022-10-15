package http

import (
	"web/health"
	webfiber "web/http/fiber"
	"web/http/kv"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	webfiber.Module,
	kv.Providers,
	fx.Invoke(func(
		app *fiber.App,
		healthController *health.Controller,
		userController *user.Controller,
	) {
		healthController.Register(app)
		userController.Register(app)
	}),
)
