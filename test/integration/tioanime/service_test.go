package tioanime

import (
	"context"
	"testing"

	"github.com/dst3v3n/api-anime/internal/domain/services/tioanime"
)

var tioAnimeService, _ = tioanime.NewTioanimeService()

func TestServiceSearch(t *testing.T) {
	testCases := []struct {
		name        string
		animeSearch string
		page        uint
		wantError   bool
	}{
		{
			name:        "Búsqueda anime 'Naruto' página 1",
			animeSearch: "naruto",
			page:        1,
			wantError:   false,
		},
		{
			name:        "Búsqueda anime 'one piece' página 0",
			animeSearch: "one piece",
			page:        0,
			wantError:   false,
		},
		{
			name:        "Búsqueda todos los animes con paginación",
			animeSearch: "",
			page:        0,
			wantError:   false,
		},
		{
			name:        "Busqueda anime 'una pieza' página 0",
			animeSearch: "una pieza",
			page:        0,
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := tioAnimeService.Search(ctx, tc.animeSearch, tc.page)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, wantError = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.Animes) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.name)
					t.Error(result)
				}
			}
		})
	}
}

func TestServiceAnimeInfo(t *testing.T) {
	testCases := []struct {
		name      string
		ID        string
		wantError bool
	}{
		{
			name:      "Obtener información anime 'one piece'",
			ID:        "one-piece-tv",
			wantError: false,
		},
		{
			name:      "Obtener información anime 'one piece'",
			ID:        "one piece",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := tioAnimeService.AnimeInfo(ctx, tc.ID)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, wantError = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.ID) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.name)
					t.Error(result)
				}
			}
		})
	}
}

func TestServiceEpisodeStreamingProvider(t *testing.T) {
	testCases := []struct {
		name      string
		ID        string
		episode   uint
		wantError bool
	}{
		{
			name:      "Obtener proveedores de streaming 'one-piece-tv' episode 1000",
			ID:        "one-piece-tv",
			episode:   1000,
			wantError: false,
		},
		{
			name:      "Obtener proveedores de streaming 'one-piece-tv' episode 2000",
			ID:        "one-piece-tv",
			episode:   2000,
			wantError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := tioAnimeService.EpisodeStreamingProvider(ctx, tc.ID, tc.episode)
			t.Log(result)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, wantError = %v", tc.name, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.ID) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.name)
					t.Error(result)
				}
			}
		})
	}
}
