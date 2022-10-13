package kv

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	client *redis.Client
}

func NewKV(client *redis.Client) *Store {
	return &Store{client}
}

func (s *Store) Set(key, value string) error {
	return s.client.Set(context.Background(), key, value, 0).Err()
}

func (s *Store) Get(key string) (string, error) {
	cmd := s.client.Get(context.Background(), key)
	if cmd.Err() != nil {
		return "", ErrNotFound
	}

	return cmd.Val(), nil
}

func (s *Store) Del(key string) error {
	return s.client.Del(context.Background(), key).Err()
}
