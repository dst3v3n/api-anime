// Package types es el paquete principal que expone la API pública de la librería.
// Define alias de tipos para re-exportar los DTOs internos, permitiendo que los usuarios
// utilicen estos tipos sin necesidad de importar paquetes internos (internal/).
// Esto mantiene la abstracción y facilita cambios internos sin afectar la API pública.
package types

import (
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

// AnimeResponse contiene el resultado de una búsqueda de animes con información de paginación.
// Incluye la lista de animes encontrados y el total de páginas disponibles.
type AnimeResponse = dto.AnimeResponse

// AnimeStruct contiene la información básica de un anime (título, sinopsis, puntuación, etc.).
type AnimeStruct = dto.AnimeStruct

// AnimeInfoResponse contiene información detallada completa de un anime específico.
// Extiende AnimeStruct con géneros, estado, episodios, animes relacionados y próximo episodio.
type AnimeInfoResponse = dto.AnimeInfoResponse

// LinkResponse contiene los enlaces de reproducción/descarga disponibles para un episodio.
// Incluye múltiples servidores de video con sus URLs y códigos de embed.
type LinkResponse = dto.LinkResponse

// VideoURL representa una URL directa de video con su resolución.
type VideoURL = dto.VideoURL

// EpisodeListResponse contiene información resumida de un episodio en un listado.
// Se utiliza para mostrar episodios recientes sin toda la información completa.
type EpisodeListResponse = dto.EpisodeListResponse

type ConfigCache = config.CacheConfig
