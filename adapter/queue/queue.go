package queue

import (
	userqueue "app/user/queue"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var (
	ClientModule    = fx.Module("queue-client", ClientProviders, ClientInvokes)
	ClientProviders = fx.Options(
		fx.Provide(NewServeMux),
		fx.Provide(NewClient),
	)
	ClientInvokes = fx.Options(
		fx.Invoke(HookClient),
		fx.Invoke(RegisterQueues),
	)

	WorkerModule    = fx.Module("queue-worker", ClientModule, WorkerProviders, WorkerInvokes)
	WorkerProviders = fx.Options(
		fx.Provide(NewServer),
	)
	WorkerInvokes = fx.Options(
		fx.Invoke(HookServer),
	)
)

func RegisterQueues(
	mux *asynq.ServeMux,
	greeter *userqueue.Greeter,
) {
	greeter.Register(mux)
}
