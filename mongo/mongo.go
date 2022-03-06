package mongo

import (
	"context"
	"web/config"
	"web/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Deps      = fx.Options(config.Module, logger.Module)
	Providers = fx.Options(fx.Provide(NewConnection), fx.Provide(NewDB))
)

func NewConnection(config config.AppConfig, logger *zap.Logger) (*mongo.Client, error) {
	ctx := context.Background()

	logger.Info("config", zap.Reflect("config", config))
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	logger.Info("connected to mongo")

	return client, nil
}

func NewDB(c *mongo.Client) *mongo.Database {
	return c.Database("template")
}
