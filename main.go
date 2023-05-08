package main

import (
	"app/adapters/health"
	"app/adapters/http"
	"app/adapters/logger"
	"app/adapters/metrics"
	"app/config"
	"app/kv"
	"app/storage/mongo"
	"app/storage/redis"
	"app/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.Module,
		config.Module,
		metrics.Module,
		health.Module,
		http.Module,
		user.Module,
		// postgres.Module,
		mongo.Module,
		redis.Module,
		kv.Module,
	).Run()
}
