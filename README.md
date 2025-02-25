## Features

- [Dependency Injection](cmd/app/app.go#L37) with [Fx](https://github.com/uber-go/fx)
- [Configuration](config/config.go#L50) with [Viper](https://github.com/spf13/viper)
- [Logging](adapter/logger/logger.go#L11) with [Zap](https://github.com/uber-go/zap)
- [Metrics](adapter/metrics/metrics.go#L32) with [Prometheus](https://github.com/prometheus/client_golang)
- [Tracing](adapter/tracing/tracing.go#L21) with [Open Telemetry](https://opentelemetry.io/)
- [Health checks](adapter/health/health_http/controller.go#L19)
- [Feature](adapter/featureflags/featureflags.go#L20) [flags](user/user_http/controller.go#L78) with [Go Feature Flag](https://github.com/thomaspoignant/go-feature-flag)
- [Validation](user/user_http/dto.go#L4) with [Validator](https://github.com/go-playground/validator)
- [Authentication](adapter/casdoor/token_provider.go#38) with [Casdoor](https://github.com/casdoor/casdoor)
- [Authorization](adapter/casbin/permission_manager.go#L20) with [Casbin](https://github.com/casbin/casbin)
- [HTTP](adapter/http/fiber.go#L35) with [Fiber](https://github.com/gofiber/fiber)
- [C](adapter/watermill/command.go)[Q](adapter/watermill/event.go)[R](user/emitter.go#18)[S](user/user_adapter/prom_probe.go#28) with [Watermill](https://github.com/ThreeDotsLabs/watermill)
- [Background](adapter/queue/client.go#L12) [tasks](user/user_queue/greeter_queue.go#L20)/[workers](cmd/worker/worker.go#L17) with [Asynq](https://github.com/hibiken/asynq)
- [Testing](user/service_test.go#L96) with [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega)
- [Migrations](cmd/migrate/migrate.go#L23) with [Goose](https://github.com/pressly/goose)
- [Sto](user/user_adapter/mongo_repository.go)[ra](kv/kv_adapter/redis_store.go)[ge](user/user_adapter/postgres_repository.go) with [Mongo](https://github.com/mongodb/mongo-go-driver), [Redis](https://github.com/redis/go-redis), and [Gorm](https://github.com/go-gorm/gorm) ([Postgres](https://github.com/go-gorm/postgres))
- [Version](.github/workflows/commit.yaml#L25) [management](.releaserc.yaml) with [Semantic Release](https://github.com/semantic-release/semantic-release)
- [Image](Dockerfile) (using [distroless](https://github.com/GoogleContainerTools/distroless)) [publishing to GitHub Container Registry](.github/workflows/build.yaml)

## Setup

```sh
# set config path manually
export CONFIG_PATH="$PWD/config/config.yaml"
# or use mise (recommended)
mise trust

docker compose up -d

# run migrations
go run cmd/migrate/migrate.go up

# run tests
go test ./...

# start server
go run cmd/app/app.go

# start worker
go run cmd/worker/worker.go
```
