package ports

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

type ScraperPort interface {
	SearchAnime(ctx context.Context, anime string, page uint) (dto.Anime, error)
	AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfo, error)
	EpisodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error)
}

type RecentScraperPort interface {
	RecentEpisode(ctx context.Context) ([]dto.RecentEpisode, error)
	RecentAnime(ctx context.Context) (dto.Anime, error)
}
