package http

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

var (
	FiberModule = fx.Module("fiber", FiberProviders, FiberInvokes)

	FiberProviders = fx.Options(
		fx.Provide(NewFiber),
		fx.Provide(NewPromProbe),
		fx.Provide(func(probe *PromProbe) Probe { return probe }),
	)
	FiberInvokes = fx.Options(
		fx.Invoke(HookFiber),
	)
)

type Probe interface {
	Middleware(*fiber.Ctx) error
}

func NewFiber(config Config, probe Probe) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})
	app.Use(otelfiber.Middleware())
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               config.Limiter.Requests,
		Expiration:        config.Limiter.Expiration,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(probe.Middleware)
	app.Use(cors.New())
	return app
}

func HookFiber(lc fx.Lifecycle, app *fiber.App, config Config) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf(":%d", config.Port)); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})
}
