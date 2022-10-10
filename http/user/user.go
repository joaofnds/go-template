package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewUserController),
	fx.Invoke(func(app *fiber.App, controller *UserController) {
		controller.Register(app)
	}),
)
