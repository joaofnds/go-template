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

func (c HealthChecker) CheckHealth(ctx context.Context) error {
	return c.client.PingContext(ctx)
}
