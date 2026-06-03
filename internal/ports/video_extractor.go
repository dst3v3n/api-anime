// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// video_extractor.go define VideoExtractor, la interfaz que debe implementar
// cualquier servicio de extracción de URLs de video desde páginas embebidas.
// Esto permite cambiar la implementación de extracción de video (StreamTape, Mega, etc.)
// sin afectar la lógica de negocio de la aplicación.
package ports

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

// VideoExtractor define el contrato que debe cumplir cualquier servicio de extracción de videos.
// Proporciona un método para extraer URLs de reproducción directa desde páginas embebidas
// que contienen reproductores de video, resolviendo direcciones indirectas a URLs reales.
type VideoExtractor interface {
	// ExtractVideoURL extrae la URL directa de reproducción de video desde una página embebida.
	// Parámetros:
	//   - ctx: contexto para control de ciclo de vida y timeout
	//   - embedURL: URL del reproductor embebido o página que contiene el video (ej: streamtape.com/e/...)
	// Retorna:
	//   - videoURLs: lista de URLs directas del video con sus resoluciones
	//   - err: error si falla la extracción o navegación de la página
	ExtractVideoURL(ctx context.Context, embedURL string) (videoURLs []dto.VideoURL, err error)
}
