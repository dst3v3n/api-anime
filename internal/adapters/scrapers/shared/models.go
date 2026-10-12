package shared

import (
	"net/http"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter    *rate.Limiter
	httpClient *http.Client
}

type URL struct {
	BaseURL      string
	SearchURL    string
	AnimeInfoURL string
	EpisodeURL   string
}
