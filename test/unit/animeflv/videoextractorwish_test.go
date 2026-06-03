package animeflv

import (
	"context"
	"testing"

	"github.com/dst3v3n/api-anime/extract"
)

func TestVideoExtractor_ExtractVideoLinksWish(t *testing.T) {
	testCases := []struct {
		name        string
		urlEmbed    string
		resolution  string
		description string
		wantError   bool
	}{
		{
			name:        "extraer enlaces de video exitosamente",
			urlEmbed:    "https://streamwish.to/e/ss619zjv2ufo",
			description: "debe extraer correctamente los enlaces de video sin errores",
			wantError:   false,
		},
		{
			name:        "extraer enlaces de video con error",
			urlEmbed:    "https://www.animeflv.net/embed/invalid",
			description: "debe manejar correctamente el error al extraer enlaces de video de una URL inválida",
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			links, err := extract.ExtractURL("streamwish", ctx, tc.urlEmbed)
			if (err != nil) != tc.wantError {
				t.Errorf("Test %s fallido: %s. Error esperado: %v, Error obtenido: %v", tc.name, tc.description, tc.wantError, err)
			}
			if !tc.wantError {
				for _, link := range links {
					t.Logf("Enlace de video extraído: %s (%s)", link.URL, link.Resolution)
				}
			}
		})
	}
}
