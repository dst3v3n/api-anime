package tioanime

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/shared"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
	"golang.org/x/time/rate"
)

var (
	_ ports.ScraperPort       = (*Client)(nil)
	_ ports.RecentScraperPort = (*Client)(nil)
)

const service = "tioanime"

type Client struct {
	url  shared.URL
	base *shared.Client
}

func NewClient() *Client {
	return &Client{
		url: shared.URL{
			BaseURL:      "https://tioanime.com",
			SearchURL:    "https://tioanime.com/directorio",
			AnimeInfoURL: "https://tioanime.com/anime",
			EpisodeURL:   "https://tioanime.com/ver",
		},
		base: shared.NewClient(rate.Limit(3), 5, 30*time.Second),
	}
}

func (c *Client) SearchAnime(ctx context.Context, anime string, page uint) (dto.Anime, error) {
	const op = "SearchAnime"

	params := map[string]string{
		"q": anime,
		"p": strconv.Itoa(int(page)),
	}
	url := shared.BuildURL(c.url.SearchURL, params)

	resp, err := c.base.DoRequest(ctx, op, url, service)
	if err != nil {
		return dto.Anime{}, err
	}
	defer resp.Body.Close()

	return parserSearch(resp.Body)
}

func (c *Client) AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfo, error) {
	const op = "AnimeInfo"
	url := fmt.Sprintf("%s/%s", c.url.AnimeInfoURL, idAnime)

	resp, err := c.base.DoRequest(ctx, op, url, service)
	if err != nil {
		return dto.AnimeInfo{}, err
	}
	defer resp.Body.Close()

	return parserAnimeInfo(resp.Body, idAnime)
}

func (c *Client) EpisodeStreamingProvider(ctx context.Context, idAnime string, episode uint) (dto.EpisodeStreamingProvider, error) {
	const op = "EpisodeStreamingProvider"

	url := fmt.Sprintf("%s/%s-%d", c.url.EpisodeURL, idAnime, episode)

	resp, err := c.base.DoRequest(ctx, op, url, service)
	if err != nil {
		return dto.EpisodeStreamingProvider{}, err
	}
	defer resp.Body.Close()

	return parserEpisodeStreaming(resp.Body, idAnime, episode)
}

func (c *Client) RecentEpisode(ctx context.Context) ([]dto.RecentEpisode, error) {
	const op = "RecentEpisode"

	resp, err := c.base.DoRequest(ctx, op, c.url.BaseURL, "tioanime")
	if err != nil {
		return []dto.RecentEpisode{}, err
	}
	defer resp.Body.Close()

	return parserRecentEpisode(resp.Body)
}

func (c *Client) RecentAnime(ctx context.Context) (dto.Anime, error) {
	const op = "RecentAnime"

	resp, err := c.base.DoRequest(ctx, op, c.url.BaseURL, "tioanime")
	if err != nil {
		return dto.Anime{}, err
	}
	defer resp.Body.Close()

	return parserSearch(resp.Body)
}
