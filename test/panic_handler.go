package test

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var PanicHandler = fx.Invoke(addPanicHandler)

func addPanicHandler(app *fiber.App) {
	app.All("panic", func(c *fiber.Ctx) error {
		panic("panic handler")
	})
}
