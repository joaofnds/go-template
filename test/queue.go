package test

import (
	"app/adapters/redis"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var QueueProvider = fx.Options(
	fx.Provide(func(config redis.Config) *asynq.Client {
		return asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr})
	}),
	fx.Provide(func() *asynq.ServeMux { return asynq.NewServeMux() }),
)
