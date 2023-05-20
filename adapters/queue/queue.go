package queue

import (
	userQueue "app/user/queue"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var sharedProviders = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(func(logger *AsyncZapLogger) asynq.Logger { return logger }),
	fx.Provide(NewServeMux),
	fx.Provide(NewClient),
	fx.Invoke(HookClient),
	fx.Invoke(Register),
)

var ClientModule = fx.Module(
	"queue-client",
	sharedProviders,
)

var WorkerModule = fx.Module(
	"queue-worker",
	sharedProviders,
	fx.Provide(NewServer),
	fx.Invoke(HookServer),
)

func Register(
	mux *asynq.ServeMux,
	greeter *userQueue.Greeter,
) {
	greeter.Register(mux)
}
