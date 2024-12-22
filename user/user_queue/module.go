package user_queue

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"user queue",
		Providers,
		Invokes,
	)

	Providers = fx.Options(
		fx.Provide(NewGreeter),
	)

	Invokes = fx.Options(
		fx.Invoke(func(greeter *Greeter, processor *cqrs.EventProcessor) error {
			return greeter.RegisterEventHandlers(processor)
		}),
	)
)
