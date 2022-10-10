package user

import "go.uber.org/fx"

var Module = fx.Module(
	"user",
	fx.Provide(NewUserService),

	fx.Provide(NewMongoRepository),
	fx.Provide(func(repo *MongoRepository) Repository { return repo }),

	fx.Provide(NewPromHabitInstrumentation),
	fx.Provide(func(intr *PromHabitInstrumentation) Instrumentation { return intr }),
)
