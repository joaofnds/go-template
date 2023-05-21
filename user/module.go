package user

import (
	"app/user/queue"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	queue.Module,

	fx.Provide(NewUserService),

	// fx.Provide(NewMongoRepository),
	// fx.Provide(func(repo *MongoRepository) Repository { return repo }),

	fx.Provide(NewPostgresRepository),
	fx.Provide(func(repo *PostgresRepository) Repository { return repo }),

	fx.Provide(NewPromProbe),
	fx.Provide(func(probe *PromProbe) Probe { return probe }),

	fx.Provide(NewController),
)
