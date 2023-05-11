package http

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NopProbeProvider = fx.Decorate(func() Probe { return NopProbe{} })

type NopProbe struct{}

func (p NopProbe) Middleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
