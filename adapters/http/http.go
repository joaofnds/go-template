package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(NewFiber),
	fx.Invoke(HookFiber),
	fx.Provide(NewPromProbe),
	fx.Provide(func(probe *PromProbe) Probe { return probe }),
)
