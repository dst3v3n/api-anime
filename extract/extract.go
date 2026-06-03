// Package extract proporciona una fachada de conveniencia para extraer URLs de video
// desde páginas embebidas de diferentes plataformas de reproducción.
// Este paquete actúa como una interfaz simplificada que expone la funcionalidad
// de extracción de video sin requerir conocimiento de los detalles internos de implementación
// de cada adaptador.
package extract

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/adapters/videoextractor"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

// ExtractURL extrae la URL directa de reproducción de video desde una página embebida
// delegando la tarea al extractor correspondiente según el servicio solicitado.
//
// Parámetros:
//   - service: El nombre del proveedor de video (ej: "streamstape", "streamwish").
//   - ctx: Contexto para el control del ciclo de vida, timeouts y cancelación de la solicitud.
//   - url: URL del reproductor embebido de donde se extraerá el video.
//   - resolution: Resolución deseada para el video (ej: "720", "1080"), sujeta a disponibilidad del proveedor.
//
// Retorna:
//   - urlResponse: lista de URLs directas del video con sus resoluciones.
//   - err: Error si el servicio no está soportado o si falla la navegación, extracción o timeout.
//
// Nota:
//
//	Dependiendo del servicio, esta función puede requerir que Chrome/Chromium esté
//	instalado en el sistema para la automatización del navegador mediante Chromedp.
func ExtractURL(service string, ctx context.Context, url string) (urlResponse []dto.VideoURL, err error) {
	switch service {
	case "streamstape":
		return videoextractor.NewSteamStape().ExtractVideoURL(ctx, url)
	case "streamwish":
		return videoextractor.NewStreamWish().ExtractVideoURL(ctx, url)
	default:
		return
	}
}
