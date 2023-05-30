package http

import (
	"app/adapter/health"
	"app/kv"
	userhttp "app/user/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(app *fiber.App, healthController *health.Controller) { healthController.Register(app) }),
	fx.Invoke(func(app *fiber.App, userController *userhttp.Controller) { userController.Register(app) }),
	fx.Invoke(func(app *fiber.App, kvController *kv.Controller) { kvController.Register(app) }),
)
