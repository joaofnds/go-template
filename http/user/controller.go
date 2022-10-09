package user

import (
	"net/http"
	"strings"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController struct {
	service *user.UserService
	logger  *zap.Logger
}

func NewUserController(service *user.UserService, logger *zap.Logger) *UserController {
	return &UserController{service, logger}
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	var out strings.Builder

	users, err := c.service.List()
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	for _, u := range users {
		out.WriteString(u.Name)
	}

	return ctx.SendString(out.String())
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var user user.User
	err := ctx.BodyParser(&user)
	if err != nil {
		c.logger.Error("failed to parse body", zap.Error(err))
		return ctx.SendStatus(http.StatusBadRequest)
	}

	_, err = c.service.CreateUser(user.Name)
	if err != nil {
		c.logger.Error("failed to create user", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	c.logger.Info("user created", zap.Reflect("user", user))

	return ctx.SendStatus(http.StatusCreated)
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	err := c.service.Remove(user.User{Name: name})
	if err != nil {
		c.logger.Error("failed to remove user", zap.String("name", name))
	}

	c.logger.Info("user removed", zap.String("name", name))

	return ctx.SendStatus(http.StatusOK)
}
