package user

import "go.uber.org/fx"

var Module = fx.Module(
	"user",
	fx.Provide(NewController),
	fx.Provide(NewUserService),

	fx.Provide(NewMongoRepository),
	fx.Provide(func(repo *MongoRepository) Repository { return repo }),

	fx.Provide(NewPromInstrumentation),
	fx.Provide(func(instr *PromInstrumentation) Instrumentation { return instr }),
)
