// Package animeflv implementa la capa de servicios de dominio para AnimeFlv.
// Este archivo (animeflv_service.go) contiene el servicio principal que orquesta
// las operaciones de negocio relacionadas con AnimeFlv, delegando a servicios
// especializados (search, recent, detail) para mantener la separación de responsabilidades.
// Actúa como una fachada que simplifica la interacción con múltiples sub-servicios.
package animeflv

import (
	"context"
	"fmt"

	"github.com/dst3v3n/api-anime/internal/adapters/cache"
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
	"github.com/valkey-io/valkey-go"
)

// AnimeflvService es el servicio principal que coordina las operaciones de AnimeFlv.
// Delega las responsabilidades a servicios especializados para búsqueda, contenido
// reciente y detalles de anime, siguiendo el principio de responsabilidad única.
// Integra caché distribuido (Valkey) en todos los sub-servicios para optimizar rendimiento.
type AnimeflvService struct {
	scraper ports.ScraperPort
	search  searchService
	recent  recentService
	detail  detailService
}

// NewAnimeflvService crea una nueva instancia del servicio AnimeFlv.
// Inicializa el scraper, conexión a Valkey para caché distribuido,
// y todos los sub-servicios necesarios para las operaciones.
func NewAnimeflvService() *AnimeflvService {
	scraper := animeflv.NewClient()

	logger := config.Logging()
	config := config.GetConfig()

	initAddress := fmt.Sprintf("redis://%s:%d/%d", config.CacheHost, config.CachePort, config.CacheDB)

	client, err := valkey.NewClient(valkey.MustParseURL(initAddress))
	if err != nil {
		logger.Error().Err(err).Msg("Error al conectar con Valkey")
		panic(err)
	}
	return &AnimeflvService{
		scraper: scraper,
		search: searchService{
			scraper:     scraper,
			cache:       cache.NewValkeyCache(client),
			enableCache: config.CacheEnabled,
		},
		recent: recentService{
			scraper:     scraper,
			cache:       cache.NewValkeyCache(client),
			enableCache: config.CacheEnabled,
		},
		detail: detailService{
			scraper:     scraper,
			cache:       cache.NewValkeyCache(client),
			enableCache: config.CacheEnabled,
		},
	}
}

// SearchAnime busca animes por nombre con paginación.
// Delega la operación al servicio de búsqueda especializado.
func (afs *AnimeflvService) SearchAnime(ctx context.Context, anime string, page uint) (dto.AnimeResponse, error) {
	return afs.search.SearchAnime(ctx, anime, page)
}

// Search obtiene todos los animes disponibles sin filtros de búsqueda.
// Delega la operación al servicio de búsqueda especializado con caché integrado.
func (afs *AnimeflvService) Search(ctx context.Context, page uint) (dto.AnimeResponse, error) {
	return afs.search.Search(ctx, page)
}

// AnimeInfo obtiene información detallada de un anime específico por su ID.
// Delega la operación al servicio de detalles.
func (afs *AnimeflvService) AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfoResponse, error) {
	return afs.detail.AnimeInfo(ctx, idAnime)
}

// Links obtiene los enlaces de reproducción para un episodio específico.
// Delega la operación al servicio de detalles.
func (afs *AnimeflvService) Links(ctx context.Context, idAnime string, episode uint) (dto.LinkResponse, error) {
	return afs.detail.Links(ctx, idAnime, episode)
}

// RecentAnime obtiene la lista de animes recientemente agregados.
// Delega la operación al servicio de contenido reciente.
func (afs *AnimeflvService) RecentAnime(ctx context.Context) ([]dto.AnimeStruct, error) {
	return afs.recent.RecentAnime(ctx)
}

// RecentEpisode obtiene la lista de episodios recientemente publicados.
// Delega la operación al servicio de contenido reciente.
func (afs *AnimeflvService) RecentEpisode(ctx context.Context) ([]dto.EpisodeListResponse, error) {
	return afs.recent.RecentEpisode(ctx)
}
