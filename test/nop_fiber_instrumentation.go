package test

import (
	webFiber "web/http/fiber"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NopHTTPInstrumentation = fx.Decorate(newNopFiberInstrumentation)

type nopHTTPInstrumentation struct{}

func newNopFiberInstrumentation() webFiber.Instrumentation {
	return &nopHTTPInstrumentation{}
}

func (i *nopHTTPInstrumentation) Middleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
