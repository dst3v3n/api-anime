package cache

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/ports"
)

type noop struct{}

func NewNoop() ports.CachePort { return noop{} }

func (noop) Get(context.Context, string, any) error       { return ports.ErrCacheMiss }
func (noop) Set(context.Context, string, any) error       { return nil }
func (noop) Delete(context.Context, string) error         { return nil }
func (noop) Exists(context.Context, string) (bool, error) { return false, nil }
