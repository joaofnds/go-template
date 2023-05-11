package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type HealthChecker struct {
	client *mongo.Client
}

func NewHealthChecker(client *mongo.Client) HealthChecker {
	return HealthChecker{client}
}

func (h HealthChecker) CheckHealth(ctx context.Context) error {
	return h.client.Ping(ctx, nil)
}
