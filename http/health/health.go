package health

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Providers = fx.Invoke(HealthHandler)
)

func HealthHandler(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}
