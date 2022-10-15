package main

import (
	"web/config"
	"web/health"
	"web/http"
	"web/kv"
	"web/logger"
	"web/metrics"
	"web/mongo"
	"web/redis"
	"web/user"

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
		redis.Module,
		mongo.Module,
		kv.Module,
	).Run()
}
