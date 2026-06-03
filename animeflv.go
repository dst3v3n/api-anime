// Package apianime es la API pública de la librería.
// Proporciona una interfaz simplificada para acceder a todos los servicios de scraping de anime.
// Los usuarios externos importan este paquete para usar la funcionalidad de la librería.
// Implementa el patrón Facade para encapsular la complejidad del sistema interno.
package apianime

import (
	"context"

	"github.com/dst3v3n/api-anime/config"
	log "github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/services/animeflv"
	"github.com/dst3v3n/api-anime/types"
)

// AnimeFlv es la fachada principal que expone públicamente todos los servicios de anime.
// Encapsula el servicio interno de dominio y proporciona métodos para búsqueda, información
// detallada, enlaces de reproducción y contenido reciente.
// Todos los métodos delegan al servicio interno, abstrayendo los detalles de implementación
// y permitiendo a los usuarios externos usar una interfaz simple y consistente.
type AnimeFlv struct {
	service *animeflv.AnimeflvService
}

// NewAnimeFlv crea una nueva instancia del servicio público de AnimeFlv.
// Inicializa el servicio interno con todas sus dependencias (scraper, caché distribuido, etc.).
// El servicio interno maneja automáticamente la configuración, conexión a Valkey y toda la lógica de negocio.
// Retorna un pointer a AnimeFlv listo para usar.
// Nota: Si hay error en la inicialización interna, esto puede causar un panic.
func NewAnimeFlv() (*AnimeFlv, error) {
	logger := log.Logging()
	if err := config.InitConfig(); err != nil {
		logger.Debug().Interface("Error", err).Msg("No se pudo inicializar la configuración")
		return nil, err
	}
	return &AnimeFlv{
		service: animeflv.NewAnimeflvService(),
	}, nil
}

// SearchAnime busca animes por nombre con soporte de paginación.
// Delega la operación al servicio interno de búsqueda que maneja:
// - Validación de parámetros (nombre no vacío)
// - Normalización de entrada (minúsculas, reemplazo de espacios)
// - Caché distribuido (si está habilitado)
// - Scraping desde AnimeFlv
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//   - anime: nombre o palabra clave de búsqueda del anime (no diferencia mayúsculas/minúsculas)
//   - page: número de página para paginación (1-based, 0 se convierte a 1)
//
// Retorna: respuesta con lista de animes encontrados y número total de páginas
func (s *AnimeFlv) SearchAnime(ctx context.Context, anime string, page uint) (types.AnimeResponse, error) {
	return s.service.SearchAnime(ctx, anime, page)
}

// Search obtiene todos los animes disponibles sin filtros de búsqueda.
// Delega la operación al servicio interno que maneja:
// - Paginación de resultados
// - Caché distribuido para optimizar consultas frecuentes
// - Scraping desde la página de búsqueda sin parámetros
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//   - page: número de página para paginación (1-based, 0 se convierte a 1)
//
// Retorna: respuesta con lista de todos los animes disponibles y número total de páginas
func (s *AnimeFlv) Search(ctx context.Context, page uint) (types.AnimeResponse, error) {
	return s.service.Search(ctx, page)
}

// AnimeInfo obtiene información detallada y completa de un anime por su ID.
// Delega la operación al servicio interno que maneja:
// - Validación del ID (no vacío)
// - Normalización del ID (minúsculas)
// - Caché distribuido para evitar scraping repetido
// - Scraping de información completa incluyendo episodios, géneros, estado y animes relacionados
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//   - idAnime: identificador único del anime (ej: "one-piece-tv", "naruto", "bleach")
//
// Retorna: respuesta con información completa del anime (sinopsis, géneros, episodios, etc.)
func (s *AnimeFlv) AnimeInfo(ctx context.Context, idAnime string) (types.AnimeInfoResponse, error) {
	return s.service.AnimeInfo(ctx, idAnime)
}

// Links obtiene los enlaces de reproducción para un episodio específico.
// Delega la operación al servicio interno que maneja:
// - Validación de parámetros
// - Caché distribuido para evitar scraping repetido del mismo episodio
// - Scraping de la página de reproducción
// - Extracción de enlaces de múltiples servidores de video
// Retorna información de múltiples servidores de video con URLs y códigos de embed.
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//   - idAnime: identificador único del anime
//   - episode: número del episodio (1-based)
//
// Retorna: respuesta con título, número de episodio y lista de enlaces de diferentes servidores
func (s *AnimeFlv) Links(ctx context.Context, idAnime string, episode uint) (types.LinkResponse, error) {
	return s.service.Links(ctx, idAnime, episode)
}

// RecentAnime obtiene la lista de animes recientemente agregados al sitio.
// Delega la operación al servicio interno que maneja:
// - Scraping de la página principal de AnimeFlv
// - Extracción de los últimos animes añadidos
// - Caché distribuido para evitar scraping repetido
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//
// Retorna: slice de animes recientes con información básica (ID, título, sinopsis, etc.)
func (s *AnimeFlv) RecentAnime(ctx context.Context) ([]types.AnimeStruct, error) {
	return s.service.RecentAnime(ctx)
}

// RecentEpisode obtiene la lista de episodios recientemente publicados.
// Delega la operación al servicio interno que maneja:
// - Scraping de la página principal de AnimeFlv
// - Extracción de los últimos episodios publicados
// - Caché distribuido para evitar scraping repetido
// Parámetros:
//   - ctx: contexto para control de ciclo de vida y timeout
//
// Retorna: slice de episodios recientes con información resumida (anime, número, capítulo, imagen)
func (s *AnimeFlv) RecentEpisode(ctx context.Context) ([]types.EpisodeListResponse, error) {
	return s.service.RecentEpisode(ctx)
}
