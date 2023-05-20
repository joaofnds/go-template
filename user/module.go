package user

import (
	"app/user/queue"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	queue.Module,

	fx.Provide(NewUserService),

	fx.Provide(NewMongoRepository),
	fx.Provide(func(repo *MongoRepository) Repository { return repo }),

	// fx.Provide(NewPostgresRepository),
	// fx.Provide(func(repo *PostgresRepository) Repository { return repo }),
	// fx.Invoke(AutoMigrate),

	fx.Provide(NewPromProbe),
	fx.Provide(func(probe *PromProbe) Probe { return probe }),

	fx.Provide(NewController),
	fx.Invoke(func(app *fiber.App, controller *Controller) { controller.Register(app) }),
)
