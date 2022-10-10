package test

import (
	webFiber "web/http/fiber"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NopHTTPInstrumentation = fx.Decorate(newNopHTTPInstrumentation)

type nopHTTPInstrumentation struct{}

type PromHabitInstrumentation struct{}

func newNopHTTPInstrumentation() webFiber.HTTPInstrumentation {
	return &nopHTTPInstrumentation{}
}

func (i *nopHTTPInstrumentation) Middleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
