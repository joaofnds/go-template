package user_module

import (
	"app/internal/mill"
	"app/user"
	"app/user/user_adapter"
	"app/user/user_http"
	"app/user/user_queue"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"user",

		fx.Provide(user_adapter.NewPostgresRepository, fx.Private),
		fx.Provide(func(repo *user_adapter.PostgresRepository) user.Repository { return repo }, fx.Private),
		fx.Provide(user.NewEmitter, fx.Private),

		fx.Provide(user.NewUserService),
		fx.Provide(user.NewPermissionService),
	)

	QueueWorkerModule = fx.Module(
		"user queue worker",
		fx.Provide(user_queue.NewGreeterWorker, fx.Private),
		fx.Provide(user_queue.NewPermissionsCleanupWorker, fx.Private),

		fx.Invoke(func(
			mux *asynq.ServeMux,
			greeter *user_queue.GreeterWorker,
			permissionsCleanup *user_queue.PermissionsCleanupWorker,
		) {
			greeter.RegisterQueueHandler(mux)
			permissionsCleanup.RegisterQueueHandler(mux)
		}),
	)

	ListenerModule = fx.Module(
		"user listener",

		fx.Provide(user_queue.NewGreeterQueue, fx.Private),
		fx.Provide(user_queue.NewPermissionsCleanupQueue, fx.Private),
		fx.Provide(user_adapter.NewPromProbe, fx.Private),

		fx.Invoke(func(
			processor *cqrs.EventProcessor,
			probe *user_adapter.PromProbe,
			greeter *user_queue.GreeterQueue,
			permissionsCleanup *user_queue.PermissionsCleanupQueue,
		) error {
			return mill.RegisterEventHandlers(
				processor,
				probe,
				greeter,
				permissionsCleanup,
			)
		}),
	)

	HTTPModule = fx.Module(
		"user http",
		ServiceModule,

		fx.Provide(user_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, userController *user_http.Controller) {
			userController.Register(app)
		}),
	)

	AppModule    = fx.Options(ListenerModule, HTTPModule)
	WorkerModule = fx.Options(QueueWorkerModule)
)
