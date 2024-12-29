package http

import (
	"app/adapter/health"
	"app/authn/authn_http"
	"app/kv/kv_http"
	"app/user/user_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(
		app *fiber.App,
		authnController *authn_http.Controller,
		healthController *health.Controller,
		kvController *kv_http.Controller,
		userController *user_http.Controller,
	) {
		authnController.Register(app)
		healthController.Register(app)
		kvController.Register(app)
		userController.Register(app)
	}),
)
