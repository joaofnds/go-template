package main

import (
	"app/adapter/casbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/adapter/validation"
	"app/adapter/watermill"
	"app/config"
	"app/internal/appcontext"
	"app/user/user_module"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		appcontext.Module,
		config.Module,
		logger.Module,
		validation.Module,
		watermill.Module,

		casbin.Module,
		queue.WorkerModule,
		postgres.Module,

		user_module.WorkerModule,
	).Run()
}
