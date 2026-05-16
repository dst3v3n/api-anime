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
	"github.com/dst3v3n/api-anime/internal/ports"
)

type StreamWish struct{}

func NewStreamWish() ports.VideoExtractor { return StreamWish{} }

func (t StreamWish) ExtractVideoURL(ctx context.Context, embedURL string, resolution string) (string, error) {
	archivos := []string{"master.m3u8", "master.txt"}

	ctx, cancel := chromedp.NewContext(ctx)
	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer func() { timeoutCancel(); cancel() }()

	var finalVideoURL string
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
					var urlEncontrada bool

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

									if altoResolucion == resolution || resolution == "default" || resolution == "" {

										baseIdx := strings.LastIndex(masterRequestURL, "/")
										if baseIdx != -1 {
											baseURL := masterRequestURL[:baseIdx+1]
											finalVideoURL = baseURL + line
										} else {
											finalVideoURL = line
										}

										urlEncontrada = true
										break
									}
								}
							}
							ultimaEtiquetaStream = ""
						}
					}

					if !urlEncontrada {
						errFiltro = fmt.Errorf("la resolución '%s' no está disponible en este servidor", resolution)
					}

					close(encontrado)
				}(e.RequestID, e.Response.URL)
			}
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(embedURL)); err != nil {
		return "", fmt.Errorf("error al automatizar la navegación: %w", err)
	}

	select {
	case <-encontrado:
		if errFiltro != nil {
			return "", errFiltro
		}
		return finalVideoURL, nil
	case <-ctx.Done():
		return "", fmt.Errorf("tiempo límite excedido buscando la resolución '%s'", resolution)
	}
}
