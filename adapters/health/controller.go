package health

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Checker
}

func NewHealthController(service Checker) *Controller {
	return &Controller{service}
}

func (c *Controller) Register(app *fiber.App) {
	app.Get("/health", c.CheckHealth)
}

func (c *Controller) CheckHealth(ctx *fiber.Ctx) error {
	check := c.service.CheckHealth(ctx.Context())
	if !check.AllUp() {
		return ctx.Status(http.StatusServiceUnavailable).JSON(check)
	}
	return ctx.JSON(check)
}
