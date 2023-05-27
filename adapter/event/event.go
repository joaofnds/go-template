package event

import (
	useradapters "app/user/adapter"
	userqueue "app/user/queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"event",

	fx.Invoke(func(probe *useradapters.PromProbe) { probe.Listen() }),
	fx.Invoke(func(greeter *userqueue.Greeter) { greeter.Listen() }),
)
