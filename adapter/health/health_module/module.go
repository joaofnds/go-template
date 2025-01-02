package health_module

import (
	"app/adapter/health"
	"app/adapter/health/health_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"health",
		fx.Provide(health.NewHealthService),
		fx.Provide(func(service *health.Service) health.Checker { return service }),
	)

	HTTPModule = fx.Module(
		"health http",
		Module,
		fx.Provide(health_http.NewController, fx.Private),
		fx.Invoke(func(app *fiber.App, healthController *health_http.Controller) {
			healthController.Register(app)
		}),
	)
)
