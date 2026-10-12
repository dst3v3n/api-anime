// Package cache implementa el adaptador de caché usando Valkey.
// Proporciona una implementación del puerto CachePort utilizando Valkey como motor de caché.
// Maneja serialización/deserialización de objetos Go a JSON para almacenamiento en Valkey.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/ports"
	"github.com/valkey-io/valkey-go"
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

func (v *Valkey) Get(ctx context.Context, key string, dest any) error {
	res := v.client.Do(ctx, v.client.B().JsonGet().Key(key).Build())

	if err := res.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return fmt.Errorf("cache get %q: %w", key, err)
		}
		return fmt.Errorf("cache get %q: %w", key, err)
	}

	if err := res.DecodeJSON(dest); err != nil {
		return fmt.Errorf("cache decode %q: %w", key, err)
	}
	return nil
}

func (v *Valkey) Set(ctx context.Context, key string, value any) error {
	if value == nil {
		return errors.New("cache: cannot store nil value")
	}

	ttl := time.Duration(config.GetConfig().CacheTTL) * time.Minute

	setCmd := v.client.B().JsonSet().
		Key(key).
		Path("$").
		Value(valkey.JSON(value)).
		Build()

	expireCmd := v.client.B().Expire().
		Key(key).
		Seconds(int64(ttl.Seconds())).
		Build()

	results := v.client.DoMulti(
		ctx,
		v.client.B().Multi().Build(),
		setCmd,
		expireCmd,
		v.client.B().Exec().Build(),
	)

	// Errores de transporte o de encolado (por ejemplo, el módulo JSON no existe).
	for _, r := range results {
		if err := r.Error(); err != nil {
			return fmt.Errorf("cache set %q: %w", key, err)
		}
	}

	// EXEC devuelve un arreglo con el resultado de cada comando ejecutado.
	execResults, err := results[len(results)-1].ToArray()
	if err != nil {
		return fmt.Errorf("cache set %q: %w", key, err)
	}
	for _, msg := range execResults {
		if err := msg.Error(); err != nil {
			return fmt.Errorf("cache set %q: %w", key, err)
		}
	}
	return nil
}

// Delete elimina una clave del caché.
// Retorna error si la operación de eliminación falla.
func (v *Valkey) Delete(ctx context.Context, key string) error {
	return v.client.Do(ctx, v.client.B().Del().Key(key).Build()).Error()
}

// Exists verifica si una clave existe en el caché.
// Retorna true si la clave existe, false en caso contrario.
// Retorna error si la operación de verificación falla.
func (v *Valkey) Exists(ctx context.Context, key string) (bool, error) {
	return v.client.Do(ctx, v.client.B().Exists().Key(key).Build()).AsBool()
}
