package main

import (
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/adapter/validation"
	"app/adapter/watermill"
	"app/config"
	"app/internal/appcontext"
	"app/user/user_queue"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		appcontext.Module,
		config.Module,
		logger.Module,
		validation.Module,
		watermill.Module,

		queue.WorkerModule,
		postgres.Module,

		user_queue.Module,
	).Run()
}
