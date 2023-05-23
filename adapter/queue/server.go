package queue

import (
	"app/adapter/redis"
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

func NewServer(config redis.Config, logger *zap.Logger) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Addr},
		asynq.Config{Concurrency: 10, Logger: NewAsyncZapLogger(logger)},
	)
}

func HookServer(lifecycle fx.Lifecycle, server *asynq.Server, mux *asynq.ServeMux) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Run(mux); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown()
			return nil
		},
	})
}
