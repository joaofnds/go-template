package http

import (
	"web/health"
	webfiber "web/http/fiber"
	"web/http/kv"
	"web/http/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	webfiber.Module,
	user.Providers,
	kv.Providers,
	fx.Invoke(func(
		app *fiber.App,
		healthController *health.Controller,
	) {
		healthController.Register(app)
	}),
)
