package shared

import (
	"context"
	"fmt"
	"net/http"
	"time"

	domainErrors "github.com/dst3v3n/api-anime/internal/domain/errors"
	"golang.org/x/time/rate"
)

func NewClient(limit rate.Limit, burst int, timeout time.Duration) *Client {
	return &Client{
		limiter: rate.NewLimiter(limit, burst),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) Wait(ctx context.Context) error {
	return c.limiter.Wait(ctx)
}

func (c *Client) DoRequest(ctx context.Context, op, url, service string) (*http.Response, error) {
	if err := c.Wait(ctx); err != nil {
		return nil, domainErrors.E(service, op, kindFromError(err),
			fmt.Errorf("rate limiter cancelado: %w", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, domainErrors.E(service, op, domainErrors.KindInvalidInput,
			fmt.Errorf("error creando petición HTTP: %w", err))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, domainErrors.E(service, op, kindFromError(err),
			fmt.Errorf("error en la petición HTTP: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, domainErrors.E(service, op, kindFromStatus(resp.StatusCode),
			fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode))
	}

	return resp, nil
}
