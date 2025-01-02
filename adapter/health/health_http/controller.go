package health_http

import (
	"app/adapter/health"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service health.Checker
}

func NewController(service health.Checker) *Controller {
	return &Controller{service}
}

func (controller *Controller) Register(app *fiber.App) {
	app.Get("/health", controller.CheckHealth)
}

func (controller *Controller) CheckHealth(ctx *fiber.Ctx) error {
	check := controller.service.CheckHealth(ctx.UserContext())
	if !check.AllUp() {
		return ctx.Status(http.StatusServiceUnavailable).JSON(check)
	}
	return ctx.JSON(check)
}
