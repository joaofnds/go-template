package user_http

import (
	"app/user/user_module"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user http",
	user_module.Module,
	user_module.ListenerModule,
	fx.Provide(NewController),
)
