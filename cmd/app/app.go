package main

import (
	"app/adapter/authz_casbin"
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
	"app/config"
	"app/internal/appcontext"
	"app/kv"
	"app/user/user_module"

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
		authz_casbin.Module,

		event.Module,
		queue.ClientModule,
		http.Module,
		postgres.Module,
		redis.Module,

		user_module.Module,
		kv.Module,
	).Run()
}
