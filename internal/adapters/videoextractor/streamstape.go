// Package videoextractor implementa adaptadores concretos para extraer URLs de video
// desde diferentes plataformas de reprodución embebidas.
// streamtape.go contiene la implementación específica para StreamStape,
// utilizando automatización de navegador (Chromedp) para acceder a contenido
// protegido por JavaScript y extraer la URL directa del video.
package videoextractor

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dst3v3n/api-anime/internal/ports"
)

// StreamStape es la implementación del puerto VideoExtractor para StreamStape.
// Utiliza Chromedp (automatización de Chrome headless) para navegación y extracción.
type StreamStape struct{}

// NewSteamStape crea una nueva instancia del extractor de videos de StreamTape.
// Retorna una interfaz VideoExtractor para permitir polimorfismo e inyección de dependencias.
func NewSteamStape() ports.VideoExtractor {
	return StreamStape{}
}

// ExtractVideoURL extrae la URL directa de reproducción desde una página embebida de StreamStape.
// Proceso:
//  1. Crea un contexto de navegación con Chromedp (Chrome headless)
//  2. Navega a la URL embebida proporcionada
//  3. Espera a que el elemento <video> sea visible en el DOM
//  4. Aguarda 2 segundos para asegurar que el script del reproductor se haya ejecutado completamente
//  5. Extrae el atributo 'src' del elemento <video> usando JavaScript
//  6. Retorna la URL directa del video o error si falla la navegación/extracción
//
// Parámetros:
//   - ctx: contexto para control de ciclo de vida (cancelación, timeout)
//   - embedURL: URL del reproductor embebido de StreamTape (ej: https://streamtape.com/e/abc123/)
//
// Retorna:
//   - videoSrc: URL directa del archivo de video (string)
//   - err: error si falla la navegación, espera, evaluación JavaScript o timeout
//
// Nota: Este método requiere que Chrome/Chromium esté instalado en el sistema.
//
//	Es un método CPU-intensivo y debe usarse con cuidado para no saturar recursos.
func (t StreamStape) ExtractVideoURL(ctx context.Context, embedURL string, resolution string) (string, error) {
	// Crea un nuevo contexto de navegación con Chromedp
	// defer cancel() libera recursos del navegador cuando se completa la función
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Variable que almacenará la URL extraída del atributo src del elemento video
	var videoSrc string

	// Ejecuta una serie de acciones de navegación:
	// 1. Navigate: Carga la página del reproductor
	// 2. WaitVisible: Espera a que el elemento <video> sea visible (indica que la página cargó)
	// 3. Sleep: Pausa de 2 segundos para ejecutación de scripts del reproductor
	// 4. Evaluate: Ejecuta JavaScript para obtener el src del video
	err := chromedp.Run(ctx,
		chromedp.Navigate(embedURL),
		chromedp.WaitVisible(`video`, chromedp.ByQuery),
		chromedp.Sleep(20*time.Millisecond),

		chromedp.Evaluate(`document.querySelector('video').src`, &videoSrc),
	)

	return videoSrc, err
}
