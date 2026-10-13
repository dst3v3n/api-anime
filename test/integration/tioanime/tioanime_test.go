package tioanime

import (
	"context"
	"testing"

	apianime "github.com/dst3v3n/api-anime"
)

func newTestTioanime(t *testing.T) *apianime.Tioanime {
	t.Helper()

	api, err := apianime.NewTioanime()
	if err != nil {
		t.Fatalf("no se pudo inicializar Tioanime: %v", err)
	}

	return api
}

func TestIntegrationSearch(t *testing.T) {
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
			name:        "búsqueda anime 'dragon bill' página 0",
			animeSearch: "dragon bill",
			page:        0,
			wantError:   true,
			description: "debe manejar el error al buscar un anime que no existe",
		},
	}

	api := newTestTioanime(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := api.Search(ctx, tc.animeSearch, tc.page)
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

func TestIntegrationAnimeInfo(t *testing.T) {
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

	api := newTestTioanime(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := api.AnimeInfo(ctx, tc.ID)
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

func TestIntegrationEpisodeStreamingProvider(t *testing.T) {
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
		{
			name:        "proveedores de streaming con ID inválido",
			ID:          "one piece",
			episode:     1,
			wantError:   true,
			description: "debe manejar el error cuando el ID del anime es inválido",
		},
	}

	api := newTestTioanime(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := api.EpisodeStreamingProvider(ctx, tc.ID, tc.episode)
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

func TestIntegrationRecentAnime(t *testing.T) {
	api := newTestTioanime(t)
	ctx := context.Background()

	result, err := api.RecentAnime(ctx)
	if err != nil {
		t.Fatalf("error inesperado obteniendo animes recientes: %v", err)
	}

	if len(result.Animes) == 0 {
		t.Error("esperaba animes recientes, pero no se encontraron")
	}
}

func TestIntegrationRecentEpisode(t *testing.T) {
	api := newTestTioanime(t)
	ctx := context.Background()

	result, err := api.RecentEpisode(ctx)
	if err != nil {
		t.Fatalf("error inesperado obteniendo episodios recientes: %v", err)
	}

	if len(result) == 0 {
		t.Error("esperaba episodios recientes, pero no se encontraron")
	}
}

func TestIntegrationContextCancelled(t *testing.T) {
	api := newTestTioanime(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := api.Search(ctx, "naruto", 1); err == nil {
		t.Error("esperaba error con el contexto cancelado en Search, pero no hubo error")
	}

	if _, err := api.AnimeInfo(ctx, "one-piece-tv"); err == nil {
		t.Error("esperaba error con el contexto cancelado en AnimeInfo, pero no hubo error")
	}

	if _, err := api.EpisodeStreamingProvider(ctx, "one-piece-tv", 1000); err == nil {
		t.Error("esperaba error con el contexto cancelado en EpisodeStreamingProvider, pero no hubo error")
	}

	if _, err := api.RecentAnime(ctx); err == nil {
		t.Error("esperaba error con el contexto cancelado en RecentAnime, pero no hubo error")
	}

	if _, err := api.RecentEpisode(ctx); err == nil {
		t.Error("esperaba error con el contexto cancelado en RecentEpisode, pero no hubo error")
	}
}
