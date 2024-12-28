package user_queue

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"user queue",
		fx.Provide(NewGreeter),
	)
)
