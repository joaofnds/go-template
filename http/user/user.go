package user

import (
	"net/http"
	"strings"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Providers = fx.Options(
		fx.Invoke(UserListHandler),
		fx.Invoke(UserCreateHandler),
		fx.Invoke(DeleteHandler),
	)
)

func UserListHandler(app *fiber.App, userService *user.UserService) {
	app.Get("/users", func(c *fiber.Ctx) error {
		var out strings.Builder

		users, err := userService.List()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		for _, u := range users {
			out.WriteString(u.Name)
		}

		return c.SendString(out.String())
	})
}

func UserCreateHandler(app *fiber.App, userService *user.UserService, logger *zap.Logger) {
	app.Post("/users", func(c *fiber.Ctx) error {
		var user user.User
		err := c.BodyParser(&user)
		if err != nil {
			logger.Error("failed to parse body", zap.Error(err))
			return c.SendStatus(http.StatusBadRequest)
		}

		_, err = userService.CreateUser(user.Name)
		if err != nil {
			logger.Error("failed to create user", zap.Error(err))
			return c.SendStatus(http.StatusInternalServerError)
		}

		logger.Info("user created", zap.Reflect("user", user))

		return c.SendStatus(http.StatusCreated)
	})
}

func DeleteHandler(app *fiber.App, userService *user.UserService, logger *zap.Logger) {
	app.Delete("/users/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		err := userService.Remove(user.User{Name: name})
		if err != nil {
			logger.Error("failed to remove user", zap.String("name", name))
		}

		logger.Info("user removed", zap.String("name", name))

		return c.SendStatus(http.StatusOK)
	})
}
