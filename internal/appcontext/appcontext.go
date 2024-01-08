package appcontext

import (
	"context"

	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"appcontext",
		fx.Provide(func(lc fx.Lifecycle) context.Context {
			ctx, cancel := context.WithCancel(context.Background())

			lc.Append(fx.Hook{
				OnStop: func(_ context.Context) error {
					cancel()
					return nil
				},
			})

			return ctx
		}),
	)
)
