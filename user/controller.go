package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewController(service *Service) *Controller {
	return &Controller{service}
}

type Controller struct {
	service *Service
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
	var user User
	err := ctx.BodyParser(&user)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	_, err = c.service.CreateUser(user.Name)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	user, err := c.service.FindByName(name)
	switch {
	case errors.Is(err, ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	case errors.Is(err, ErrRepository):
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if err = c.service.Remove(user); err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}
