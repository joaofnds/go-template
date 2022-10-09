package kv

import (
	"errors"
	"net/http"
	"web/kv"

	"github.com/gofiber/fiber/v2"
)

type KVController struct {
	store *kv.KeyValStore
}

func NewKVController(store *kv.KeyValStore) *KVController {
	return &KVController{store}
}

func (c *KVController) Get(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	if key == "" {
		return fiber.NewError(http.StatusBadRequest, "missing key")
	}

	val, err := c.store.Get(key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).Send([]byte(val))
}

func (c *KVController) Set(ctx *fiber.Ctx) error {
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

func (c *KVController) Delete(ctx *fiber.Ctx) error {
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
