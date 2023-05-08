package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoHealthChecker struct {
	client *mongo.Client
}

func NewMongoHealthChecker(client *mongo.Client) MongoHealthChecker {
	return MongoHealthChecker{client}
}

func (c MongoHealthChecker) CheckHealth(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}
