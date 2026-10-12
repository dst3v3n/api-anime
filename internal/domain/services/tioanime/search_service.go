package tioanime

import (
	"context"
	"fmt"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/domain/services/shared"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type searchService struct {
	scraper     ports.ScraperPort
	cache       ports.CachePort
	enableCache bool
}

func (s *searchService) Search(ctx context.Context, anime string, page uint) (dto.Anime, error) {
	cacheKey := fmt.Sprintf("search-%s-%d", anime, page)
	return shared.Cached(ctx, s.cache, s.enableCache, cacheKey, func(ctx context.Context) (dto.Anime, error) {
		return s.scraper.SearchAnime(ctx, anime, page)
	})
}
