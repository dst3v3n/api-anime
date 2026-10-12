package tioanime

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/tioanime"
	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/domain/services/shared"
)

type TioanimeService struct {
	search searchService
	recent recentService
	detail detailService
}

func NewTioanimeService() (*TioanimeService, error) {
	scraper := tioanime.NewClient()
	config := config.GetConfig()
	enabled := config.CacheEnabled

	c, err := shared.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &TioanimeService{
		search: searchService{
			scraper:     scraper,
			cache:       c,
			enableCache: enabled,
		},
		recent: recentService{
			scraper:     scraper,
			cache:       c,
			enableCache: enabled,
		},
		detail: detailService{
			scraper:     scraper,
			cache:       c,
			enableCache: enabled,
		},
	}, nil
}

func (a *TioanimeService) Search(ctx context.Context, anime string, page uint) (dto.Anime, error) {
	return a.search.Search(ctx, anime, page)
}

func (a *TioanimeService) RecentEpisode(ctx context.Context) ([]dto.RecentEpisode, error) {
	return a.recent.RecentEpisode(ctx)
}

func (a *TioanimeService) RecentAnime(ctx context.Context) (dto.Anime, error) {
	return a.recent.RecentAnime(ctx)
}

func (a *TioanimeService) AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfo, error) {
	return a.detail.animeInfo(ctx, idAnime)
}

func (a *TioanimeService) EpisodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error) {
	return a.detail.episodeStreamingProvdider(ctx, idAnime, episode)
}
