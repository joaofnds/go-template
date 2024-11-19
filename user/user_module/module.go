package user_module

import (
	"app/user"
	"app/user/user_adapter"
	"app/user/user_http"
	"app/user/user_queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",

	fx.Provide(user.NewUserService),
	fx.Provide(user.NewEventEmitter),

	fx.Provide(user_adapter.NewPostgresRepository),
	fx.Provide(func(repo *user_adapter.PostgresRepository) user.Repository { return repo }),

	fx.Provide(user_adapter.NewPromProbe),

	fx.Provide(user_http.NewController),
	fx.Provide(user_queue.NewGreeter),
)
