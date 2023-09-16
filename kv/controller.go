package kv

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	store Store
}

func NewController(store Store) *Controller {
	return &Controller{store}
}

func (controller *Controller) Register(app *fiber.App) {
	app.Get("/kv/:key", controller.Get)
	app.Post("/kv/:key/:val", controller.Set)
	app.Delete("/kv/:key", controller.Delete)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	val, err := controller.store.Get(ctx.UserContext(), key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).Send([]byte(val))
}

func (controller *Controller) Set(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	val := ctx.Params("val")
	if val == "" {
		return fiber.NewError(http.StatusBadRequest, "missing val")
	}

	err := controller.store.Set(ctx.UserContext(), key, val)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key param")
	}

	if err := controller.store.Del(ctx.UserContext(), key); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
	}

	return ctx.SendStatus(http.StatusOK)
}
