package controllers

import (
	"context"
	"fmt"
	"web/config"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	Module = fx.Options(
		fx.Provide(NewFiber),
		fx.Invoke(RegisterFiber),
		fx.Invoke(RootHandler),
	)
)

func NewFiber() *fiber.App {
	return fiber.New()
}

func RegisterFiber(lc fx.Lifecycle, app *fiber.App, config config.AppConfig) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go app.Listen(fmt.Sprintf(":%d", config.Port))
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})
}
