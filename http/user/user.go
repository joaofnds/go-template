package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewUserController),
	fx.Invoke(UserHandler),
)

func UserHandler(app *fiber.App, controller *UserController) {
	app.Get("/users", controller.List)
	app.Post("/users", controller.Create)
	app.Delete("/users/:name", controller.Delete)
}
