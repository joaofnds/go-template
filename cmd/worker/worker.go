package main

import (
	"app/adapters/logger"
	"app/adapters/postgres"
	"app/adapters/queue"
	user "app/user/module"

	"app/config"
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
