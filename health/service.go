package health

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Checker interface {
	CheckHealth(ctx context.Context) Check
}

type Service struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
}

func NewHealthService(mongoClient *mongo.Client, redisClient *redis.Client) *Service {
	return &Service{mongoClient, redisClient}
}

func (s *Service) CheckHealth(ctx context.Context) Check {
	return Check{
		Mongo: s.MongoHealth(ctx),
		Redis: s.RedisHealth(ctx),
	}
}

func (s *Service) MongoHealth(ctx context.Context) Status {
	if err := s.mongoClient.Ping(ctx, nil); err != nil {
		return Status{Status: StatusDown}
	}
	return Status{Status: StatusUp}
}

func (s *Service) RedisHealth(ctx context.Context) Status {
	if cmd := s.redisClient.Ping(ctx); cmd.Err() != nil {
		return Status{Status: StatusDown}
	}
	return Status{Status: StatusUp}
}
