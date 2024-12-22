package user_module

import (
	"app/user"
	"app/user/user_adapter"
	"app/user/user_queue"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"user",
		Imports,
		Providers,
		Invokes,
	)

	Imports = fx.Options(
		user_queue.Module,
	)

	Providers = fx.Options(
		fx.Provide(user_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *user_adapter.PostgresRepository) user.Repository { return repo }, fx.Private),
		fx.Provide(user.NewEventEmitter, fx.Private),
		fx.Provide(user_adapter.NewPromProbe, fx.Private),

		fx.Provide(user.NewUserService),
	)

	Invokes = fx.Options(
		fx.Invoke(func(promProbe *user_adapter.PromProbe, processor *cqrs.EventProcessor) error {
			return promProbe.RegisterEventHandlers(processor)
		}),
	)
)
