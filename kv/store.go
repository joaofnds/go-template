package kv

import (
	"context"
)

type Store interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Del(context.Context, string) error
}
