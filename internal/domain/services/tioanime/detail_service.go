package tioanime

import (
	"context"
	"fmt"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	apperrors "github.com/dst3v3n/api-anime/internal/domain/errors"
	"github.com/dst3v3n/api-anime/internal/domain/services/shared"
	"github.com/dst3v3n/api-anime/internal/ports"
)

const serviceName = "tioanime"

type detailService struct {
	scraper     ports.ScraperPort
	cache       ports.CachePort
	enableCache bool
}

func (d *detailService) animeInfo(ctx context.Context, idAnime string) (dto.AnimeInfo, error) {
	const op = "detailService.animeInfo"

	if idAnime == "" {
		return dto.AnimeInfo{}, apperrors.E(serviceName, op, apperrors.KindInvalidInput, apperrors.ErrEmptyAnimeID)
	}

	id := shared.NormalizeID(idAnime)
	cacheKey := fmt.Sprintf("anime-info-%s", id)

	info, err := shared.Cached(ctx, d.cache, d.enableCache, cacheKey, func(ctx context.Context) (dto.AnimeInfo, error) {
		return d.scraper.AnimeInfo(ctx, id)
	})
	if err != nil {
		return dto.AnimeInfo{}, apperrors.Wrap(serviceName, op, err)
	}
	return info, nil
}

func (d *detailService) episodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error) {
	const op = "detailService.episodeStreamingProvider"

	if idAnime == "" {
		return dto.EpisodeStreamingProvider{}, apperrors.E(serviceName, op, apperrors.KindInvalidInput, apperrors.ErrEmptyAnimeID)
	}
	if episode == 0 {
		return dto.EpisodeStreamingProvider{}, apperrors.E(serviceName, op, apperrors.KindInvalidInput, apperrors.ErrInvalidEpisode)
	}

	id := shared.NormalizeID(idAnime)
	cacheKey := fmt.Sprintf("episode-streaming-%s-%d", id, episode)

	res, err := shared.Cached(ctx, d.cache, d.enableCache, cacheKey, func(ctx context.Context) (dto.EpisodeStreamingProvider, error) {
		return d.scraper.EpisodeStreamingProvider(ctx, id, episode)
	})
	if err != nil {
		return dto.EpisodeStreamingProvider{}, apperrors.Wrap(serviceName, op, err)
	}
	return res, nil
}
