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
	fx.Invoke(func(app *fiber.App, healthController *health.Controller) { healthController.Register(app) }),
	fx.Invoke(func(app *fiber.App, userController *user_http.Controller) { userController.Register(app) }),
	fx.Invoke(func(app *fiber.App, kvController *kv_http.Controller) { kvController.Register(app) }),
	fx.Invoke(func(app *fiber.App, authnController *authn_http.Controller) { authnController.Register(app) }),
)
