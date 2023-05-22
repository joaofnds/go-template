package kv

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (s *RedisStore) Set(ctx context.Context, key, value string) error {
	return s.client.Set(ctx, key, value, 0).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	cmd := s.client.Get(ctx, key)
	if cmd.Err() != nil {
		return "", ErrNotFound
	}

	return cmd.Val(), nil
}

func (s *RedisStore) Del(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
