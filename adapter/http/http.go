package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
)
