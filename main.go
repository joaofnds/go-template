package main

import (
	"web/config"
	"web/http"
	"web/kv"
	"web/logger"
	"web/metrics"
	"web/mongo"
	"web/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.Module,
		config.Module,
		metrics.Module,
		http.Module,
		user.Module,
		mongo.Module,
		kv.Module,
	).Run()
}
