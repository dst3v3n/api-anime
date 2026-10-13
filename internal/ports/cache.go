package ports

import (
	"context"
	"errors"
)

var ErrCacheMiss = errors.New("cache: key not found")

type CachePort interface {
	Exists(ctx context.Context, key string) (bool, error)

	Get(ctx context.Context, key string, dest interface{}) error

	Set(ctx context.Context, key string, value interface{}) error

	Delete(ctx context.Context, key string) error
}
