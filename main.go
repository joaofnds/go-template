package main

import (
	"app/adapters/health"
	"app/adapters/http"
	"app/adapters/logger"
	"app/adapters/metrics"
	"app/adapters/mongo"
	"app/adapters/queue"
	"app/adapters/redis"

	"app/config"
	"app/kv"
	"app/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,
		metrics.Module,
		health.Module,

		queue.Module,
		http.Module,
		mongo.Module,
		redis.Module,

		user.Module,
		kv.Module,
	).Run()
}
