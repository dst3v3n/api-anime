package tioanime

import (
	"context"
	"fmt"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/domain/services/shared"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type detailService struct {
	scraper     ports.ScraperPort
	cache       ports.CachePort
	enableCache bool
}

func (d *detailService) animeInfo(ctx context.Context, idAnime string) (dto.AnimeInfo, error) {
	if idAnime == "" {
	}

	id := shared.NormalizeID(idAnime)
	cacheKey := fmt.Sprintf("anime-info-%s", id)

	return shared.Cached(ctx, d.cache, d.enableCache, cacheKey, func(ctx context.Context) (dto.AnimeInfo, error) {
		return d.scraper.AnimeInfo(ctx, id)
	})
}

func (d *detailService) episodeStreamingProvdider(ctx context.Context, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error) {
	if idAnime == "" {
	}

	id := shared.NormalizeID(idAnime)
	cacheKey := fmt.Sprintf("episode-streaming-%s-%d", id, episode)

	return shared.Cached(ctx, d.cache, d.enableCache, cacheKey, func(ctx context.Context) (dto.EpisodeStreamingProvider, error) {
		return d.scraper.EpisodeStreamingProvider(ctx, id, episode)
	})
}
