// Package cache implementa el adaptador de caché usando Valkey.
// Proporciona una implementación del puerto CachePort utilizando Valkey como motor de caché.
// Maneja serialización/deserialización de objetos Go a JSON para almacenamiento en Valkey.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/ports"
	"github.com/valkey-io/valkey-go"
)

// Valkey es la implementación concreta del puerto CachePort.
// Encapsula un cliente Valkey para acceso al servidor de caché distribuido.
type Valkey struct {
	client valkey.Client
	config config.CacheConfig
}

// NewValkeyCache crea una nueva instancia del adaptador de caché Valkey.
// Toma un cliente Valkey ya inicializado y retorna una instancia que implementa CachePort.
func NewValkeyCache(client valkey.Client) ports.CachePort {
	enviroment := config.GetConfig()
	return &Valkey{
		client: client,
		config: enviroment,
	}
}

// Get recupera un valor del caché por su clave y lo deserializa en el destino proporcionado.
// Utiliza un tiempo de caché de 1 minuto. Si la clave no existe, retorna nil sin error.
// Si el valor existe pero no puede ser deserializado, retorna un error.
func (v *Valkey) Get(ctx context.Context, key string, dest interface{}) error {
	ttl := v.config.CacheTTL
	resp, error := v.client.DoCache(ctx, v.client.B().Get().Key(key).Cache(), time.Duration(ttl)*time.Minute).ToString()

	if error != nil {
		if valkey.IsValkeyNil(error) {
			return fmt.Errorf("key not found in cache")
		}
		return error
	}

	return deserialize(resp, dest)
}

// Set almacena un valor en el caché con una clave especificada.
// Serializa el valor a JSON y lo guarda con una TTL de 15 minutos.
// Retorna error si falla la serialización o la operación de almacenamiento.
func (v *Valkey) Set(ctx context.Context, key string, value interface{}) error {
	r, _ := regexp.Compile(`^\s*[\[{]`)

	if value == nil {
		return fmt.Errorf("cannot cache nil value")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	if !r.Match(data) {
		return fmt.Errorf("Error: valor no es serializable a JSON")
	}

	cmd := v.client.B().Set().Key(key).Value(string(data)).Ex(time.Minute).Build()
	return v.client.Do(ctx, cmd).Error()
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
	result, err := v.client.Do(ctx, v.client.B().Exists().Key(key).Build()).AsInt64()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
