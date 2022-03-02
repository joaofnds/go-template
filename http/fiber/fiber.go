package fiber

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
		fx.Invoke(HookFiber),
	)
)

func NewFiber() *fiber.App {
	return fiber.New()
}

func HookFiber(lc fx.Lifecycle, app *fiber.App, config config.AppConfig) {
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
