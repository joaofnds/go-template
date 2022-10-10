package health

import (
	"net/http"
	"web/health"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewHealthController),
	fx.Invoke(func(app *fiber.App, controller *HealthController) {
		controller.Register(app)
	}),
)

type HealthController struct {
	service health.HealthChecker
}

func NewHealthController(service health.HealthChecker) *HealthController {
	return &HealthController{service}
}

func (c *HealthController) Register(app *fiber.App) {
	app.Get("/health", c.CheckHealth)
}

func (c *HealthController) CheckHealth(ctx *fiber.Ctx) error {
	status := http.StatusOK

	h := c.service.CheckHealth(ctx.Context())
	if !h.AllUp() {
		status = http.StatusServiceUnavailable
	}

	return ctx.Status(status).JSON(h)
}
