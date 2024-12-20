package main

import (
	"app/adapter/casbin"
	"app/adapter/casdoor"
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
	"app/authn/authn_http"
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
		casbin.Module,
		casdoor.Module,

		event.Module,
		queue.ClientModule,
		postgres.Module,
		redis.Module,
		http.Module,

		authn_http.Module,
		user_module.Module,
		kv.Module,
	).Run()
}
