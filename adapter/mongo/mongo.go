package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module = fx.Module("mongo", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewClient),
		fx.Provide(NewHealthChecker),
	)
	Invokes = fx.Options(
		fx.Invoke(HookConnection),
	)
)

func NewClient(ctx context.Context, config Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.URI)
	return mongo.Connect(ctx, clientOptions)
}

func HookConnection(lifecycle fx.Lifecycle, client *mongo.Client, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := client.Ping(ctx, nil)
			if err != nil {
				logger.Error("failed to ping mongo", zap.Error(err))
				return err
			}

			logger.Info("successfully pinged mongo")

			err = EnsureIndexes(ctx, client)
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

func EnsureIndexes(ctx context.Context, client *mongo.Client) error {
	indexView := client.Database("template").Collection("users").Indexes()
	_, err := indexView.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"email": 1}})
	return err
}
