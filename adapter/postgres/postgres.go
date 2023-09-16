package postgres

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	Module = fx.Module("postgres", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewGORMDB),
		fx.Provide(NewSQLDB),
		fx.Provide(NewHealthChecker),
	)
	Invokes = fx.Options(
		fx.Invoke(EnableTracing),
		fx.Invoke(HookConnection),
	)
)

func NewGORMDB(postgresConfig Config, logger *zap.Logger) (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(postgresConfig.Addr),
		&gorm.Config{
			Logger: NewZapLoggerAdapter(logger, gormlogger.Config{
				LogLevel:      gormlogger.Info,
				SlowThreshold: 100 * time.Millisecond,
			}),
			PrepareStmt:              true,
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		},
	)
}

func EnableTracing(db *gorm.DB) error {
	return db.Use(tracing.NewPlugin())
}

func NewSQLDB(orm *gorm.DB) (*sql.DB, error) { return orm.DB() }

func HookConnection(lifecycle fx.Lifecycle, sqlDB *sql.DB, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return sqlDB.PingContext(ctx) },
		OnStop:  func(ctx context.Context) error { return sqlDB.Close() },
	})
}
