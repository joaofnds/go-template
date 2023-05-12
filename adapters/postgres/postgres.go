package postgres

import (
	"context"
	"database/sql"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"postgres",
	fx.Provide(NewGORMDB),
	fx.Provide(NewSQLDB),
	fx.Provide(NewHealthChecker),
	fx.Invoke(HookConnection),
)

func NewGORMDB(postgresConfig Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(postgresConfig.Addr))
	if err != nil {
		logger.Error("failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func NewSQLDB(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}

func HookConnection(lifecycle fx.Lifecycle, sqlDB *sql.DB, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := sqlDB.PingContext(ctx); err != nil {
				logger.Error("failed to ping db", zap.Error(err))
				return err
			}
			logger.Info("successfully pinged db")

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := sqlDB.Close(); err != nil {
				logger.Error("failed to close db connection", zap.Error(err))
				return err
			}

			return nil
		},
	})
}
