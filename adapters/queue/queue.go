package queue

import (
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"queue",
	fx.Provide(NewLogger),
	fx.Provide(func(l *AsyncZapLogger) asynq.Logger { return l }),

	fx.Provide(NewClient),
	fx.Invoke(HookClient),

	fx.Provide(NewServeMux),

	fx.Provide(NewServer),
	fx.Invoke(HookServer),
)
