package http

import (
	"app/adapter/health"
	"app/kv/kv_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(
		app *fiber.App,
		healthController *health.Controller,
		kvController *kv_http.Controller,
	) {
		healthController.Register(app)
		kvController.Register(app)
	}),
)
