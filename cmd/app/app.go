package main

import (
	"app/adapter/casbin"
	"app/adapter/casdoor"
	"app/adapter/featureflags"
	"app/adapter/health/health_module"
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
	"app/adapter/watermill"
	"app/authn/authn_module"
	"app/authz/authz_http"
	"app/config"
	"app/internal/appcontext"
	"app/kv/kv_module"
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
		validation.Module,
		featureflags.Module,
		uuid.Module,
		time.Module,
		casbin.Module,
		casdoor.Module,

		watermill.Module,
		queue.ClientModule,
		postgres.Module,
		redis.Module,
		http.Module,
		authz_http.Module,

		health_module.HTTPModule,
		authn_module.HTTPModule,
		user_module.AppModule,
		kv_module.Module,
	).Run()
}
