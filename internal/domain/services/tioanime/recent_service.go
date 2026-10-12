package tioanime

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/domain/services/shared"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type recentService struct {
	scraper     ports.RecentScraperPort
	cache       ports.CachePort
	enableCache bool
}

func (r *recentService) RecentEpisode(ctx context.Context) ([]dto.RecentEpisode, error) {
	cacheKey := "recent-episode"
	return shared.Cached(ctx, r.cache, r.enableCache, cacheKey, func(ctx context.Context) ([]dto.RecentEpisode, error) {
		return r.scraper.RecentEpisode(ctx)
	})
}

func (r *recentService) RecentAnime(ctx context.Context) (dto.Anime, error) {
	cacheKey := "recent-anime"
	return shared.Cached(ctx, r.cache, r.enableCache, cacheKey, func(ctx context.Context) (dto.Anime, error) {
		return r.scraper.RecentAnime(ctx)
	})
}
