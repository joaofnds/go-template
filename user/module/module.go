package user

import (
	"app/internal/event"
	"app/user"
	"app/user/adapter"
	"app/user/http"
	"app/user/queue"
	"context"

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

	fx.Provide(func(lifecycle fx.Lifecycle) *event.Event[user.UserCreated] {
		e := event.NewEvent[user.UserCreated](10)
		lifecycle.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				e.Close()
				return nil
			},
		})
		return e
	}),

	fx.Invoke(func(event *event.Event[user.UserCreated], greeter *queue.Greeter) {
		go func() {
			for userCreatedEvent := range event.Listen() {
				_ = greeter.Enqueue(userCreatedEvent.User.Name)
			}
		}()
	}),
)
