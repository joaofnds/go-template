package user

import "go.uber.org/fx"

var Module = fx.Module(
	"user",
	fx.Provide(NewPromHabitInstrumentation),
	fx.Provide(NewUserRepository),
	fx.Provide(NewUserService),
)
