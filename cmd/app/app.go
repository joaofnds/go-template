package main

import (
	"app/adapter/event"
	"app/adapter/featureflags"
	"app/adapter/health"
	"app/adapter/http"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/adapter/redis"
	"app/adapter/time"
	"app/adapter/tracing"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/authz"
	"app/config"
	"app/internal/appcontext"
	"app/kv"
	user "app/user/module"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		appcontext.Module,
		config.Module,
		logger.Module,
		metrics.Module,
		tracing.Module,
		health.Module,
		validation.Module,
		featureflags.Module,
		uuid.Module,
		time.Module,

		event.Module,
		queue.ClientModule,
		http.Module,
		postgres.Module,
		redis.Module,

		authz.Module,
		user.Module,
		kv.Module,
	).Run()
}
