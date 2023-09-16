## Features

- [Dependency Injection](cmd/app/app.go#L23) with [Fx](https://github.com/uber-go/fx)
- [Configuration](config/config.go#L47) with [Viper](https://github.com/spf13/viper)
- [Logging](adapter/logger/logger.go#L10) with [Zap](https://github.com/uber-go/zap)
- [Metrics](adapter/metrics/metrics.go#L22) with [Prometheus](https://github.com/prometheus/client_golang)
- [Tracing](adapter/tracing/tracing.go#L18) with [Open Telemetry](https://opentelemetry.io/)
- [Health checks](adapter/health/controller.go#L18)
- [Feature](adapter/featureflags/featureflags.go#L20) [flags](user/http/controller.go#L69) with [Go Feature Flag](https://github.com/thomaspoignant/go-feature-flag)
- [Validation](user/http/dto.go#L4) with [Validator](https://github.com/go-playground/validator)
- [HTTP](adapter/http/fiber.go#L34) with [Fiber](https://github.com/gofiber/fiber)
- [Background](adapter/queue/client.go#L12) [tasks](user/queue/greeter.go#L33)/[workers](cmd/worker/worker.go#L14) with [Asynq](https://github.com/hibiken/asynq)
- [Testing](user/service_test.go#L68) with [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega)
- [Migrations](cmd/migrate/migrate.go#L22) with [Goose](https://github.com/pressly/goose)
- [Sto](user/adapter/mongo_repository.go)[ra](kv/redis_store.go)[ge](user/adapter/postgres_repository.go) with [Mongo](https://github.com/mongodb/mongo-go-driver), [Redis](https://github.com/redis/go-redis), and [Gorm](https://github.com/go-gorm/gorm) ([Postgres](https://github.com/go-gorm/postgres))
- [Version](.github/workflows/commit.yaml#L64) [management](.releaserc.yaml) with [Semantic Release](https://github.com/semantic-release/semantic-release)
- [Image](Dockerfile) (using [distroless](https://github.com/GoogleContainerTools/distroless)) [publishing to GitHub Container Registry](.github/workflows/build.yaml)

## Setup

```sh
# set config path manually
export CONFIG_PATH="$PWD/config/config.yaml"
# or use direnv (recommended)
direnv allow .

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
