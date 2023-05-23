package queue

import (
	"app/adapter/redis"
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

func NewClient(config redis.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr})
}

func HookClient(lifecycle fx.Lifecycle, client *asynq.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return client.Close()
		},
	})
}
