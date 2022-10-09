package kv

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Providers = fx.Options(
	fx.Provide(NewKVController),
	fx.Invoke(KVHandler),
)

func KVHandler(app *fiber.App, controller *KVController) {
	app.Get("/kv/:key", controller.Get)
	app.Post("/kv/:key/:val", controller.Set)
	app.Delete("/kv/:key", controller.Delete)
}
