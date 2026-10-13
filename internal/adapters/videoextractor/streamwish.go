package videoextractor

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type StreamWish struct{}

func NewStreamWish() ports.VideoExtractor { return StreamWish{} }

func (t StreamWish) ExtractVideoURL(ctx context.Context, embedURL string) ([]dto.VideoURL, error) {
	archivos := []string{"master.m3u8", "master.txt"}

	ctx, cancel := chromedp.NewContext(ctx)
	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer func() { timeoutCancel(); cancel() }()

	var finalVideoURLs []dto.VideoURL
	var errFiltro error

	encontrado := make(chan struct{})

	chromedp.ListenTarget(ctx, func(ev any) {
		if e, ok := ev.(*network.EventResponseReceived); ok && e.Type == network.ResourceTypeXHR && ctx.Err() == nil {
			if slices.Contains(archivos, path.Base(e.Response.URL)) {
				go func(requestID network.RequestID, masterRequestURL string) {
					var body []byte

					err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
						var errCDP error
						body, errCDP = network.GetResponseBody(requestID).Do(ctx)
						return errCDP
					}))
					if err != nil {
						errFiltro = fmt.Errorf("error al leer el búfer de red de Chrome: %w", err)
						close(encontrado)
						return
					}

					scanner := bufio.NewScanner(bytes.NewReader(body))
					var ultimaEtiquetaStream string

					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())

						if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
							ultimaEtiquetaStream = line
							continue
						}

						if ultimaEtiquetaStream != "" && line != "" && !strings.HasPrefix(line, "#") {
							if idx := strings.Index(ultimaEtiquetaStream, "RESOLUTION="); idx != -1 {
								subStr := ultimaEtiquetaStream[idx+11:]
								if comaIdx := strings.Index(subStr, ","); comaIdx != -1 {
									subStr = subStr[:comaIdx]
								}

								partesResolucion := strings.Split(subStr, "x")

								if len(partesResolucion) == 2 {
									altoResolucion := partesResolucion[1]

									var videoURL string
									baseIdx := strings.LastIndex(masterRequestURL, "/")
									if baseIdx != -1 {
										baseURL := masterRequestURL[:baseIdx+1]
										videoURL = baseURL + line
									} else {
										videoURL = line
									}

									finalVideoURLs = append(finalVideoURLs, dto.VideoURL{
										URL:        videoURL,
										Resolution: altoResolucion + "p",
									})
								}
							}
						}
						ultimaEtiquetaStream = ""
					}
					close(encontrado)
				}(e.RequestID, e.Response.URL)
			}
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(embedURL)); err != nil {
		return nil, fmt.Errorf("error al automatizar la navegación: %w", err)
	}

	select {
	case <-encontrado:
		if errFiltro != nil {
			return nil, errFiltro
		}
		return finalVideoURLs, nil
	}
}
