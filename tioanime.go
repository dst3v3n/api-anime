package apianime

import (
	"context"

	"github.com/dst3v3n/api-anime/config"
	apperrors "github.com/dst3v3n/api-anime/internal/domain/errors"
	"github.com/dst3v3n/api-anime/internal/domain/services/tioanime"
	"github.com/dst3v3n/api-anime/types"
)

const serviceName = "apianime"

type Tioanime struct {
	service *tioanime.TioanimeService
}

func NewTioanime() (*Tioanime, error) {
	const op = "Tioanime"

	if err := config.InitConfig(); err != nil {
		return nil, apperrors.E(serviceName, op+" inicializar configuración", apperrors.KindConfig, err)
	}

	svc, err := tioanime.NewTioanimeService()
	if err != nil {
		return nil, apperrors.Wrap(serviceName, op+" inicializar servicio tioanime", err)
	}

	return &Tioanime{service: svc}, nil
}

func (s *Tioanime) Search(ctx context.Context, anime string, page uint) (types.Anime, error) {
	return s.service.Search(ctx, anime, page)
}

func (s *Tioanime) AnimeInfo(ctx context.Context, idAnime string) (types.AnimeInfo, error) {
	return s.service.AnimeInfo(ctx, idAnime)
}

func (s *Tioanime) EpisodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (types.EpisodeStreamingProvider, error) {
	return s.service.EpisodeStreamingProvider(ctx, idAnime, episode)
}

func (s *Tioanime) RecentAnime(ctx context.Context) (types.Anime, error) {
	return s.service.RecentAnime(ctx)
}

func (s *Tioanime) RecentEpisode(ctx context.Context) ([]types.RecentEpisode, error) {
	return s.service.RecentEpisode(ctx)
}
