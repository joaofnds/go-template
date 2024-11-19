package event

import (
	useradapters "app/user/user_adapter"
	userqueue "app/user/user_queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"event",

	fx.Invoke(func(probe *useradapters.PromProbe) { probe.Listen() }),
	fx.Invoke(func(greeter *userqueue.Greeter) { greeter.Listen() }),
)
