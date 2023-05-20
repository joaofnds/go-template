package queue

import (
	userQueue "app/user/queue"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"queue",
	fx.Provide(NewLogger),
	fx.Provide(func(logger *AsyncZapLogger) asynq.Logger { return logger }),

	fx.Provide(NewClient),
	fx.Invoke(HookClient),

	fx.Provide(NewServeMux),

	fx.Provide(NewServer),
	fx.Invoke(HookServer),
)

func Register(
	mux *asynq.ServeMux,
	greeter *userQueue.Greeter,
) {
	greeter.Register(mux)
}
