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

func (controller *Controller) Register(app *fiber.App) {
	app.Post("/users", controller.Create)
	app.Get("/users", controller.List)
	app.Get("/users/:name", controller.Get)
	app.Delete("/users/:name", controller.Delete)
}

func (controller *Controller) List(ctx *fiber.Ctx) error {
	users, err := controller.service.List(ctx.Context())
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(users)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body UserCreateDTO
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
	}

	createdUser, err := controller.service.CreateUser(ctx.Context(), body.Name)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(createdUser)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	userFound, err := controller.service.FindByName(ctx.Context(), name)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(userFound)
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	userFound, err := controller.service.FindByName(ctx.Context(), name)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	err = controller.service.Remove(ctx.Context(), userFound)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}
