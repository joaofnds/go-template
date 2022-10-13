package user

import (
	"net/http"
	"strings"
	"web/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewController(service *user.Service, logger *zap.Logger) *Controller {
	return &Controller{service, logger}
}

type Controller struct {
	service *user.Service
	logger  *zap.Logger
}

func (c *Controller) Register(app *fiber.App) {
	app.Get("/users", c.List)
	app.Post("/users", c.Create)
	app.Delete("/users/:name", c.Delete)
}

func (c *Controller) List(ctx *fiber.Ctx) error {
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

func (c *Controller) Create(ctx *fiber.Ctx) error {
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

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	err := c.service.Remove(user.User{Name: name})
	if err != nil {
		c.logger.Error("failed to remove user", zap.String("name", name))
	}

	c.logger.Info("user removed", zap.String("name", name))

	return ctx.SendStatus(http.StatusOK)
}
