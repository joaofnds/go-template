package kv

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewController(store Store) *Controller {
	return &Controller{store}
}

type Controller struct {
	store Store
}

func (c *Controller) Register(app *fiber.App) {
	app.Get("/kv/:key", c.Get)
	app.Post("/kv/:key/:val", c.Set)
	app.Delete("/kv/:key", c.Delete)
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	val, err := c.store.Get(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).Send([]byte(val))
}

func (c *Controller) Set(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	val := ctx.Params("val")
	if val == "" {
		return fiber.NewError(http.StatusBadRequest, "missing val")
	}

	err := c.store.Set(key, val)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	err := c.store.Del(key)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
	}

	return ctx.SendStatus(http.StatusOK)
}
