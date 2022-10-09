package kv

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("not found")

type KeyValStore struct {
	client *redis.Client
}

func NewKV(client *redis.Client) *KeyValStore {
	return &KeyValStore{client}
}

func (store *KeyValStore) Set(key, value string) error {
	return store.client.Set(context.Background(), key, value, 0).Err()
}

func (store *KeyValStore) Get(key string) (string, error) {
	s := store.client.Get(context.Background(), key)
	if s.Err() != nil {
		return "", ErrNotFound
	}

	return s.Val(), nil
}

func (store *KeyValStore) Del(key string) error {
	return store.client.Del(context.Background(), key).Err()
}
