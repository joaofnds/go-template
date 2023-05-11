package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"mongo",
	fx.Provide(NewClient),
	fx.Provide(NewHealthChecker),
	fx.Invoke(HookConnection),
)

func NewClient(config Config, logger *zap.Logger) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Error("failed to connect to mongo", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func HookConnection(lc fx.Lifecycle, client *mongo.Client, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := client.Connect(ctx)
			if err != nil {
				logger.Error("failed to connect to mongo", zap.Error(err))
				return err
			}
			logger.Info("successfully connected to mongo")

			err = client.Ping(ctx, nil)
			if err != nil {
				logger.Error("failed to ping mongo", zap.Error(err))
				return err
			}

			logger.Info("successfully pinged mongo")

			err = EnsureIndexes(client)
			if err != nil {
				logger.Error("failed to create index", zap.Error(err))
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})
}

func EnsureIndexes(client *mongo.Client) error {
	indexView := client.Database("template").Collection("users").Indexes()
	_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{Keys: bson.M{"name": 1}})
	return err
}
