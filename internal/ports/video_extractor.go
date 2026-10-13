package ports

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

type VideoExtractor interface {
	ExtractVideoURL(ctx context.Context, embedURL string) (videoURLs []dto.VideoURL, err error)
}
