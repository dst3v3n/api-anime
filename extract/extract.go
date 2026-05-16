// Package extract proporciona funciones de conveniencia para extraer URLs de video
// desde páginas embebidas de reproductores de video.
// Este paquete actúa como una interfaz simplificada que expone la funcionalidad
// de extracción de video sin requerir conocimiento de detalles internos de implementación.
package extract

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/adapters/videoextractor"
)

// ExtractURL extrae la URL directa de reproducción de video desde una página embebida.
// Es una función de conveniencia que crea una instancia del extractor de StreamTape
// y ejecuta la extracción de URL. Útil para casos simples donde no se requiere
// reutilización del extractor.
// Parámetros:
//   - ctx: contexto para control de ciclo de vida, timeout y cancelación
//   - url: URL del reproductor embebido de donde se extraerá el video (ej: streamtape.com/e/xyz/)
//
// Retorna:
//   - urlResponse: URL directa del archivo de video (tipo string)
//   - err: error si falla la navegación, extracción o timeout
//
// Nota: Esta función requiere Chrome/Chromium instalado en el sistema
//
//	para la automatización del navegador (Chromedp).
func ExtractURL(service string, ctx context.Context, url string) (urlResponse string, err error) {
	switch service {
	case "streamtape":
		return videoextractor.NewSteamTape().ExtractVideoURL(ctx, url)
	case "streamwish":
		return videoextractor.NewStreamWish().ExtractVideoURL(ctx, url)
	default:
		return
	}
}
