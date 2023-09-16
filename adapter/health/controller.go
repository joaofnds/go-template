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
