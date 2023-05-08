package test

import (
	"app/adapters/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NopHTTPInstrumentation = fx.Decorate(newNopFiberInstrumentation)

type nopHTTPInstrumentation struct{}

func newNopFiberInstrumentation() http.Instrumentation {
	return &nopHTTPInstrumentation{}
}

func (i *nopHTTPInstrumentation) Middleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
