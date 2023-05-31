package postgres

import (
	"context"
	"database/sql"
)

type HealthChecker struct {
	client *sql.DB
}

func NewHealthChecker(client *sql.DB) HealthChecker {
	return HealthChecker{client}
}

func (healthChecker HealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.PingContext(ctx)
}
