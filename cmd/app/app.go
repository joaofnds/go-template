package main

import (
	"app/adapter/health"
	"app/adapter/http"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/adapter/redis"
	user "app/user/module"

	"app/config"
	"app/kv"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,
		metrics.Module,
		health.Module,

		queue.ClientModule,
		http.Module,
		postgres.Module,
		redis.Module,

		user.Module,
		kv.Module,
	).Run()
}
