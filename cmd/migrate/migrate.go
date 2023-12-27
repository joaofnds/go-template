package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	"app/adapter/logger"
	"app/adapter/postgres"
	"app/config"

	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`usage:
  go run cmd/migrate/migrate.go up
  go run cmd/migrate/migrate.go down
  go run cmd/migrate/migrate.go down-to 20170506082527
  go run cmd/migrate/migrate.go status
  go run cmd/migrate/migrate.go redo
	go run cmd/migrate/migrate.go create`)
		os.Exit(1)
	}

	app := fx.New(
		logger.NopLoggerProvider,
		config.Module,
		postgres.Module,
		fx.Invoke(func(ctx context.Context, db *sql.DB, config postgres.Config) error {
			goose.SetBaseFS(migrations)
			action, args := os.Args[1], os.Args[2:]
			err := goose.RunContext(ctx, action, db, dir(action), args...)
			if err != nil {
				fmt.Println(strings.ReplaceAll(err.Error(), `\n`, "\n"))
			}
			return err
		}),
	)
	defer func() { _ = app.Stop(context.Background()) }()
	_ = app.Start(context.Background())
}

func dir(action string) string {
	if action == "create" {
		return "cmd/migrate/migrations"
	}

	return "migrations"
}
