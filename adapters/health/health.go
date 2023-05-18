package health

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewHealthService),
	fx.Provide(func(service *Service) Checker { return service }),
	fx.Provide(NewHealthController),
	fx.Invoke(func(app *fiber.App, controller *Controller) { controller.Register(app) }),
)
