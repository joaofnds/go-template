package user

import (
	"net/http"
	"strings"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Providers = fx.Options(
		fx.Invoke(UserListHandler),
		fx.Invoke(UserCreateHandler),
	)
)

func UserListHandler(app *fiber.App, userService *user.UserService) {
	app.Get("/users", func(c *fiber.Ctx) error {
		var out strings.Builder

		for _, u := range userService.List() {
			out.WriteString(u.Name)
		}

		return c.SendString(out.String())
	})
}

func UserCreateHandler(app *fiber.App, userService *user.UserService) {
	app.Post("/users", func(c *fiber.Ctx) error {
		var user user.User

		err := c.BodyParser(&user)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		userService.CreateUser(user.Name)

		return c.SendStatus(http.StatusCreated)
	})
}
