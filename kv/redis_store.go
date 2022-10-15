package kv

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("not found")

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (s *RedisStore) Set(key, value string) error {
	return s.client.Set(context.Background(), key, value, 0).Err()
}

func (s *RedisStore) Get(key string) (string, error) {
	cmd := s.client.Get(context.Background(), key)
	if cmd.Err() != nil {
		return "", ErrNotFound
	}

	return cmd.Val(), nil
}

func (s *RedisStore) Del(key string) error {
	return s.client.Del(context.Background(), key).Err()
}
