package main

import (
	"app/config"
	"app/health"
	"app/http"
	"app/kv"
	"app/logger"
	"app/metrics"
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
