package http

import (
	"app/user"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(service *user.Service, validator *validator.Validate) *Controller {
	return &Controller{service: service, validator: validator}
}

type Controller struct {
	service   *user.Service
	validator *validator.Validate
}

func (c *Controller) Register(app *fiber.App) {
	app.Post("/users", c.Create)
	app.Get("/users", c.List)
	app.Get("/users/:name", c.Get)
	app.Delete("/users/:name", c.Delete)
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	users, err := c.service.List(ctx.Context())
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(users)
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	var dto UserCreateDTO
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := c.validator.Struct(dto); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
	}

	u, err := c.service.CreateUser(ctx.Context(), dto.Name)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(u)
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	u, err := c.service.FindByName(ctx.Context(), name)
	switch {
	case errors.Is(err, user.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	case errors.Is(err, user.ErrRepository):
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(u)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	u, err := c.service.FindByName(ctx.Context(), name)
	switch {
	case errors.Is(err, user.ErrNotFound):
		return ctx.SendStatus(http.StatusNotFound)
	case errors.Is(err, user.ErrRepository):
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if err = c.service.Remove(ctx.Context(), u); err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}
