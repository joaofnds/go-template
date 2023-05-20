package main

import (
	"app/adapters/logger"
	"app/adapters/mongo"
	"app/adapters/queue"

	"app/config"
	"app/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,

		queue.WorkerModule,
		mongo.Module,

		user.Module,
	).Run()
}
