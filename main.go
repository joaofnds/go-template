package main

import (
	"app/adapters/health"
	"app/adapters/http"
	"app/adapters/logger"
	"app/adapters/metrics"
	"app/adapters/mongo"
	"app/adapters/redis"
	"app/config"
	"app/kv"
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
		mongo.Module,
		redis.Module,
		kv.Module,
	).Run()
}
