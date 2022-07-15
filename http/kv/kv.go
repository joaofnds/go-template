package kv

import (
	"errors"
	"net/http"
	"web/kv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Providers = fx.Options(
		fx.Invoke(GetListener),
		fx.Invoke(SetListener),
		fx.Invoke(DelListener),
	)
)

func GetListener(app *fiber.App, kvStore *kv.KeyValStore) {
	app.Get("/kv/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		if key == "" {
			return fiber.NewError(http.StatusBadRequest, "missing key")
		}

		val, err := kvStore.Get(key)
		if err != nil {
			if errors.Is(err, kv.ErrNotFound) {
				return fiber.NewError(http.StatusNotFound, err.Error())
			}
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		return c.Status(http.StatusOK).Send([]byte(val))
	})
}

func SetListener(app *fiber.App, kvStore *kv.KeyValStore) {
	app.Post("/kv/:key/:val", func(c *fiber.Ctx) error {
		key := c.Params("key")
		if key == "" {
			return fiber.NewError(http.StatusBadRequest, "missing key")
		}

		val := c.Params("val")
		if val == "" {
			return fiber.NewError(http.StatusBadRequest, "missing val")
		}

		err := kvStore.Set(key, val)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
		}

		c.Status(http.StatusCreated)
		return nil
	})
}

func DelListener(app *fiber.App, kvStore *kv.KeyValStore) {
	app.Delete("/kv/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		if key == "" {
			return fiber.NewError(http.StatusBadRequest, "missing key")
		}

		err := kvStore.Del(key)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "failed to delete key")
		}

		return c.SendStatus(http.StatusOK)
	})
}
