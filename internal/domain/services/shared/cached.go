package shared

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dst3v3n/api-anime/internal/adapters/cache"
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/ports"
	"github.com/valkey-io/valkey-go"
)

func NewCache(cfg config.CacheConfig) (ports.CachePort, error) {
	if !cfg.CacheEnabled {
		return cache.NewNoop(), nil
	}

	addr := fmt.Sprintf("%s:%d", cfg.CacheHost, cfg.CachePort)
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
		SelectDB:    cfg.CacheDB,
	})
	if err != nil {
		return nil, fmt.Errorf("conectar a valkey: %w", err)
	}

	return cache.NewValkeyCache(client), nil
}

func Cached[T any](
	ctx context.Context,
	cache ports.CachePort,
	enabled bool,
	key string,
	load func(context.Context) (T, error),
) (T, error) {
	if enabled {
		var hit T
		err := cache.Get(ctx, key, &hit)
		switch {
		case err == nil:
			return hit, nil
		case !errors.Is(err, ports.ErrCacheMiss):
			slog.WarnContext(ctx, "cache get falló", "key", key, "err", err)
		}
	}

	value, err := load(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	if enabled {
		if err := cache.Set(ctx, key, value); err != nil {
			slog.WarnContext(ctx, "cache set falló", "key", key, "err", err)
		}
	}

	return value, nil
}
