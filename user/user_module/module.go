package user_module

import (
	"app/internal/mill"
	"app/user"
	"app/user/user_adapter"
	"app/user/user_queue"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"user",

		fx.Provide(user_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *user_adapter.PostgresRepository) user.Repository { return repo }, fx.Private),
		fx.Provide(user.NewEmitter, fx.Private),

		fx.Provide(user.NewUserService),
	)

	ListenerModule = fx.Module(
		"user listener",

		user_queue.Module,

		fx.Provide(user_adapter.NewPromProbe, fx.Private),

		fx.Invoke(func(
			processor *cqrs.EventProcessor,
			greeter *user_queue.Greeter,
			probe *user_adapter.PromProbe,
		) error {
			return mill.RegisterEventHandlers(processor, greeter, probe)
		}),
	)
)
