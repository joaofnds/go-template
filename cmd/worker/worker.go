package main

import (
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/config"
	user "app/user/module"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,

		queue.WorkerModule,
		postgres.Module,

		user.Module,
	).Run()
}
