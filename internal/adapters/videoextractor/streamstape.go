package videoextractor

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type StreamStape struct{}

func NewSteamStape() ports.VideoExtractor {
	return StreamStape{}
}

func (t StreamStape) ExtractVideoURL(ctx context.Context, embedURL string) ([]dto.VideoURL, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var videoSrc string

	err := chromedp.Run(
		ctx,
		chromedp.Navigate(embedURL),
		chromedp.WaitVisible(`video`, chromedp.ByQuery),
		chromedp.Sleep(20*time.Millisecond),

		chromedp.Evaluate(`document.querySelector('video').src`, &videoSrc),
	)
	if err != nil {
		return nil, err
	}

	return []dto.VideoURL{
		{
			URL:        videoSrc,
			Resolution: "default",
		},
	}, nil
}
