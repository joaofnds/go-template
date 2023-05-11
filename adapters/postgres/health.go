package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type HealthChecker struct {
	client *sqlx.DB
}

func NewHealthChecker(client *sqlx.DB) HealthChecker {
	return HealthChecker{client}
}

func (c HealthChecker) CheckHealth(ctx context.Context) error {
	return c.client.PingContext(ctx)
}
