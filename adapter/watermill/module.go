package watermill

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module = fx.Module(
		"watermill",
		Providers,
		Invokes,
	)

	Providers = fx.Options(
		fx.Provide(newLogger, fx.Private),
		fx.Provide(NewSonicMarshaler, fx.Private),
		fx.Provide(newRouter, fx.Private),
		fx.Provide(newGoChannel, fx.Private),
		fx.Provide(func(m SonicMarshaler) cqrs.CommandEventMarshaler { return m }, fx.Private),
		fx.Provide(func(c *gochannel.GoChannel) message.Publisher { return c }, fx.Private),
		fx.Provide(func(c *gochannel.GoChannel) message.Subscriber { return c }, fx.Private),

		fx.Provide(newCommandBus),
		fx.Provide(newCommandProcessor),
		fx.Provide(newEventBus),
		fx.Provide(newEventProcessor),
	)
	Invokes = fx.Options(
		fx.Invoke(hookRouter),
	)
)

func newGoChannel(logger watermill.LoggerAdapter) *gochannel.GoChannel {
	return gochannel.NewGoChannel(gochannel.Config{}, logger)
}

func newRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	return message.NewRouter(message.RouterConfig{}, logger)
}

func hookRouter(
	ctx context.Context,
	lifecycle fx.Lifecycle,
	router *message.Router,
	logger *zap.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := router.Run(ctx)
				if err != nil {
					logger.Fatal("watermill router exited with an error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error { return router.Close() },
	})
}
