package user

import (
	"app/internal/event"
	"app/user"
	"app/user/adapter"
	"app/user/http"
	"app/user/queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	queue.Module,

	fx.Provide(user.NewUserService),

	//fx.Provide(adapter.NewMongoRepository),
	//fx.Provide(func(repo *adapter.MongoRepository) user.Repository { return repo }),

	fx.Provide(adapter.NewPostgresRepository),
	fx.Provide(func(repo *adapter.PostgresRepository) user.Repository { return repo }),

	fx.Provide(adapter.NewPromProbe),
	fx.Provide(func(probe *adapter.PromProbe) user.Probe { return probe }),

	fx.Provide(http.NewController),

	fx.Invoke(func(greeter *queue.Greeter) {
		event.On(func(event user.UserCreated) {
			_ = greeter.Enqueue(event.User.Name)
		})
	}),
)
