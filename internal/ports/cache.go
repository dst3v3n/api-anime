// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// cache.go define CachePort, la interfaz que debe implementar cualquier sistema de caché.
// Esto permite cambiar la implementación de caché (Valkey, Redis, memoria, etc.) sin
// afectar la lógica de negocio de la aplicación.
package ports

import (
	"context"
	"errors"
)

var ErrCacheMiss = errors.New("cache: key not found")

// CachePort define el contrato que debe cumplir cualquier implementación de caché.
// Proporciona operaciones básicas de almacenamiento, recuperación, eliminación y verificación de existencia.
type CachePort interface {
	// Exists verifica si una clave existe en el caché.
	Exists(ctx context.Context, key string) (bool, error)

	// Get recupera un valor del caché por su clave y lo deserializa en el destino proporcionado.
	Get(ctx context.Context, key string, dest interface{}) error

	// Set almacena un valor en el caché con una clave especificada.
	Set(ctx context.Context, key string, value interface{}) error

	// Delete elimina una clave del caché.
	Delete(ctx context.Context, key string) error
}
