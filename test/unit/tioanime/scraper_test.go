package tioanime

import (
	"context"
	"testing"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/tioanime"
)

func TestParseAnimeSearch(t *testing.T) {
	testCases := []struct {
		name        string
		animeSearch string
		page        uint
		wantError   bool
		description string
	}{
		{
			name:        "búsqueda anime 'naruto' página 1",
			animeSearch: "naruto",
			page:        1,
			wantError:   false,
			description: "debe buscar correctamente el anime 'naruto' en la página 1 sin errores",
		},
		{
			name:        "búsqueda anime 'one piece' página 0",
			animeSearch: "one piece",
			page:        0,
			wantError:   false,
			description: "debe buscar correctamente el anime 'one piece' con página 0 sin errores",
		},
		{
			name:        "búsqueda de todos los animes con paginación",
			animeSearch: "",
			page:        0,
			wantError:   false,
			description: "debe listar todos los animes cuando la búsqueda está vacía",
		},
		{
			name:        "búsqueda anime 'una pieza' página 0",
			animeSearch: "una pieza",
			page:        0,
			wantError:   true,
			description: "debe manejar el error al buscar un anime que no existe",
		},
		{
			name:        "búsqueda anime 'dragon bill' en la página 0",
			animeSearch: "dragon bill",
			page:        0,
			wantError:   true,
			description: "debe manejar el error al buscar un anime que no existe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := tioanime.NewClient()
			ctx := context.Background()

			result, err := parser.SearchAnime(ctx, tc.animeSearch, tc.page)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.Animes) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.description)
				}
			}
		})
	}
}

func TestParseAnimeInfo(t *testing.T) {
	testCases := []struct {
		name        string
		ID          string
		wantError   bool
		description string
	}{
		{
			name:        "obtener información anime 'one-piece-tv'",
			ID:          "one-piece-tv",
			wantError:   false,
			description: "debe obtener correctamente la información del anime 'one-piece-tv'",
		},
		{
			name:        "obtener información anime con ID inválido 'one piece'",
			ID:          "one piece",
			wantError:   true,
			description: "debe manejar el error cuando el ID es inválido",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := tioanime.NewClient()
			ctx := context.Background()

			result, err := parser.AnimeInfo(ctx, tc.ID)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.ID) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.description)
				}
			}
		})
	}
}

func TestParseEpisodeStreamingProvider(t *testing.T) {
	testCases := []struct {
		name        string
		ID          string
		episode     uint
		wantError   bool
		description string
	}{
		{
			name:        "proveedores de streaming 'one-piece-tv' episodio 1000",
			ID:          "one-piece-tv",
			episode:     1000,
			wantError:   false,
			description: "debe obtener los proveedores de streaming del episodio 1000 de 'one-piece-tv'",
		},
		{
			name:        "proveedores de streaming 'one-piece-tv' episodio 2000",
			ID:          "one-piece-tv",
			episode:     2000,
			wantError:   true,
			description: "debe manejar el error cuando el episodio no existe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := tioanime.NewClient()
			ctx := context.Background()

			result, err := parser.EpisodeStreamingProvider(ctx, tc.ID, tc.episode)
			t.Log(result)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.ID) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.description)
				}
			}
		})
	}
}
