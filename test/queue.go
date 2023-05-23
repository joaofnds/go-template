package test

import (
	"app/adapter/redis"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Queue = fx.Options(
	fx.Provide(func(config redis.Config) *asynq.Client {
		return asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr})
	}),
	fx.Provide(func() *asynq.ServeMux { return asynq.NewServeMux() }),
)
