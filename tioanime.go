package apianime

import (
	"context"
	"fmt"

	"github.com/dst3v3n/api-anime/config"
	"github.com/dst3v3n/api-anime/internal/domain/services/tioanime"
	"github.com/dst3v3n/api-anime/types"
)

type Tioanime struct {
	service *tioanime.TioanimeService
}

func NewTioanime() (*Tioanime, error) {
	if err := config.InitConfig(); err != nil {
		return nil, fmt.Errorf("inicializar configuración: %w", err)
	}

	svc, err := tioanime.NewTioanimeService()
	if err != nil {
		return nil, fmt.Errorf("inicializar servicio tioanime: %w", err)
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
