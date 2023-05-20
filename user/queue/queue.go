package queue

import (
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGreeter),
	fx.Invoke(func(greeter *Greeter, mux *asynq.ServeMux) { greeter.Register(mux) }),
)
