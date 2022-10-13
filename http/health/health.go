package health

import (
	"net/http"
	"web/health"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewHealthController),
	fx.Invoke(func(app *fiber.App, controller *Controller) {
		controller.Register(app)
	}),
)

type Controller struct {
	service health.Checker
}

func NewHealthController(service health.Checker) *Controller {
	return &Controller{service}
}

func (c *Controller) Register(app *fiber.App) {
	app.Get("/health", c.CheckHealth)
}

func (c *Controller) CheckHealth(ctx *fiber.Ctx) error {
	status := http.StatusOK

	h := c.service.CheckHealth(ctx.Context())
	if !h.AllUp() {
		status = http.StatusServiceUnavailable
	}

	return ctx.Status(status).JSON(h)
}
