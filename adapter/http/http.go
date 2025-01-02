package http

import (
	"app/kv/kv_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(
		app *fiber.App,
		kvController *kv_http.Controller,
	) {
		kvController.Register(app)
	}),
)
