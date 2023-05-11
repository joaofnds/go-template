package postgres

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Module = fx.Module(
	"postgres",
	fx.Provide(NewClient),
	fx.Provide(NewHealthChecker),
	fx.Invoke(HookConnection),
)

func NewClient(postgresConfig Config, logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", postgresConfig.Addr)
	if err != nil {
		logger.Error("failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func HookConnection(lifecycle fx.Lifecycle, db *sqlx.DB, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.PingContext(ctx); err != nil {
				logger.Error("failed to ping db", zap.Error(err))
				return err
			}
			logger.Info("successfully pinged db")

			if err := createTables(ctx, db); err != nil {
				logger.Error("failed to create tables", zap.Error(err))
				return err
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := db.Close(); err != nil {
				logger.Error("failed to close db connection", zap.Error(err))
				return err
			}

			return nil
		},
	})
}

func createTables(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE IF NOT EXISTS users (
      id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
  `)
	return err
}
