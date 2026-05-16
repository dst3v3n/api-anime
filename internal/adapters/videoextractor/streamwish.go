package videoextractor

import (
	"context"
	"fmt"
	"path"
	"slices"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type StreamWish struct{}

func NewStreamWish() ports.VideoExtractor { return StreamWish{} }

func (t StreamWish) ExtractVideoURL(ctx context.Context, embedURL string) (string, error) {
	archivos := []string{"master.m3u8", "master.txt"}

	ctx, cancel := chromedp.NewContext(ctx)
	ctx, timeoutCancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() { timeoutCancel(); cancel() }()

	var masterURL string
	encontrado := make(chan struct{})

	chromedp.ListenTarget(ctx, func(ev any) {
		if e, ok := ev.(*network.EventResponseReceived); ok && e.Type == network.ResourceTypeXHR && ctx.Err() == nil {
			if slices.Contains(archivos, path.Base(e.Response.URL)) {
				masterURL = e.Response.URL
				close(encontrado)
			}
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(embedURL)); err != nil && masterURL == "" {
		return "", err
	}

	select {
	case <-encontrado:
		return masterURL, nil
	case <-ctx.Done():
		if masterURL != "" {
			return masterURL, nil
		}
		return "", fmt.Errorf("timeout sin encontrar el video")
	}
}
