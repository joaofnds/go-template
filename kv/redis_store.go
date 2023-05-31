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

func (store *RedisStore) Set(ctx context.Context, key, value string) error {
	return store.client.Set(ctx, key, value, 0).Err()
}

func (store *RedisStore) Get(ctx context.Context, key string) (string, error) {
	cmd := store.client.Get(ctx, key)
	if cmd.Err() != nil {
		return "", ErrNotFound
	}

	return cmd.Val(), nil
}

func (store *RedisStore) Del(ctx context.Context, key string) error {
	return store.client.Del(ctx, key).Err()
}
