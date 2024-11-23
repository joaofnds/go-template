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
	"app/adapter/validation"
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

	ctx := context.Background()

	db, dbErr := dbInstance(ctx)
	if dbErr != nil {
		fmt.Println(dbErr)
		os.Exit(1)
	}

	goose.SetBaseFS(migrations)

	action, args := os.Args[1], os.Args[2:]
	gooseErr := goose.RunContext(ctx, action, db, dir(action), args...)
	if gooseErr != nil {
		fmt.Println(strings.ReplaceAll(gooseErr.Error(), `\n`, "\n"))
		os.Exit(1)
	}
}

func dbInstance(ctx context.Context) (*sql.DB, error) {
	var db *sql.DB

	app := fx.New(
		logger.NopLoggerProvider,
		validation.Module,
		config.Module,
		postgres.Module,
		fx.Populate(&db),
	)

	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("app.Err: %w", err)
	}

	if err := app.Stop(ctx); err != nil {
		return nil, fmt.Errorf("app.Stop: %w", err)
	}

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return db, nil
}

func dir(action string) string {
	if action == "create" {
		return "cmd/migrate/migrations"
	}

	return "migrations"
}
