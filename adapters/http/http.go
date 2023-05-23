package http

import (
	"app/adapters/health"
	"app/kv"
	userhttp "app/user/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberProvider,
	fx.Invoke(func(
		app *fiber.App,
		healthController *health.Controller,
		userController *userhttp.Controller,
		kvController *kv.Controller,
	) {
		healthController.Register(app)
		userController.Register(app)
		kvController.Register(app)
	}),
)
