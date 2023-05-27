package user

import (
	"app/user"
	"app/user/adapter"
	"app/user/http"
	"app/user/queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",

	fx.Provide(user.NewUserService),
	fx.Provide(user.NewEventEmitter),

	fx.Provide(adapter.NewPostgresRepository),
	fx.Provide(func(repo *adapter.PostgresRepository) user.Repository { return repo }),

	fx.Provide(adapter.NewPromProbe),

	fx.Provide(http.NewController),
	fx.Provide(queue.NewGreeter),
)
