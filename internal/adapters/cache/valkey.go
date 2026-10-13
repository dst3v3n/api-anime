package cache

import (
	"context"
	"errors"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/dst3v3n/api-anime/internal/config"
	apperrors "github.com/dst3v3n/api-anime/internal/domain/errors"
	"github.com/dst3v3n/api-anime/internal/ports"
)

var _ ports.CachePort = (*Valkey)(nil)

var ErrInvalidCacheValue = errors.New("cache: valor no serializable a JSON")

const (
	serviceName = "cache"
	defaultTTL  = 60 * time.Minute
)

type Valkey struct {
	client valkey.Client
	config config.CacheConfig
}

func NewValkeyCache(client valkey.Client) ports.CachePort {
	enviroment := config.GetConfig()
	return &Valkey{
		client: client,
		config: enviroment,
	}
}

func (v *Valkey) ttl() time.Duration {
	ttl := time.Duration(v.config.CacheTTL) * time.Minute
	if ttl <= 0 {
		return defaultTTL
	}
	return ttl
}

// fail reemplaza a fmt.Errorf: mismo contexto (op + clave), pero con Kind.
func fail(op, key string, err error) error {
	return apperrors.E(serviceName, op+" "+key, apperrors.Classify(err), err)
}

func (v *Valkey) Get(ctx context.Context, key string, dest any) error {
	const op = "Valkey.Get"

	res := v.client.Do(ctx, v.client.B().JsonGet().Key(key).Build())

	if err := res.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return ports.ErrCacheMiss
		}
		return fail(op, key, err)
	}

	if err := res.DecodeJSON(dest); err != nil {
		return apperrors.E(serviceName, op+" decode "+key, apperrors.KindExternal, err)
	}

	return nil
}

func (v *Valkey) Set(ctx context.Context, key string, value any) error {
	const op = "Valkey.Set"

	data, err := encodeJSON(value)
	if err != nil {
		return apperrors.E(serviceName, op+" "+key, apperrors.KindInvalidInput, err)
	}

	setCmd := v.client.B().JsonSet().
		Key(key).
		Path("$").
		Value(valkey.BinaryString(data)).
		Build()

	expireCmd := v.client.B().Expire().
		Key(key).
		Seconds(int64(v.ttl().Seconds())).
		Build()

	results := v.client.DoMulti(
		ctx,
		v.client.B().Multi().Build(),
		setCmd,
		expireCmd,
		v.client.B().Exec().Build(),
	)

	for _, r := range results {
		if err := r.Error(); err != nil {
			return fail(op, key, err)
		}
	}

	execResults, err := results[len(results)-1].ToArray()
	if err != nil {
		return fail(op, key, err)
	}
	for _, msg := range execResults {
		if err := msg.Error(); err != nil {
			return fail(op, key, err)
		}
	}

	return nil
}

func (v *Valkey) Delete(ctx context.Context, key string) error {
	const op = "Valkey.Delete"

	if err := v.client.Do(ctx, v.client.B().Del().Key(key).Build()).Error(); err != nil {
		return fail(op, key, err)
	}
	return nil
}

func (v *Valkey) Exists(ctx context.Context, key string) (bool, error) {
	const op = "Valkey.Exists"

	n, err := v.client.Do(ctx, v.client.B().Exists().Key(key).Build()).AsInt64()
	if err != nil {
		return false, fail(op, key, err)
	}
	return n > 0, nil
}
