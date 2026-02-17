package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HealthChecker struct {
	client *mongo.Client
}

func NewHealthChecker(client *mongo.Client) HealthChecker {
	return HealthChecker{client}
}

func (healthChecker HealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx, nil)
}
