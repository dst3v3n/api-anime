// Package animeflv contiene tests de integración para los servicios de dominio de AnimeFlv.
// Este archivo (service_test.go) implementa pruebas que verifican la integración completa
// de los servicios (búsqueda, detalles, contenido reciente) con caché distribuido,
// validando comportamientos como normalización de datos, manejo de cache y delegación al scraper.
package animeflv

import (
	"context"
	"testing"

	"github.com/dst3v3n/api-anime/internal/config"
	"github.com/dst3v3n/api-anime/internal/domain/services/animeflv"
)

func TestSearchAnimeService(t *testing.T) {
	testCases := []struct {
		name        string
		anime       string
		page        uint
		wantError   bool
		description string
	}{
		{
			name:        "búsqueda de anime exitosa",
			anime:       "dragon ball",
			page:        2,
			wantError:   false,
			description: "debe buscar correctamente el anime con paginación",
		},
		{
			name:        "búsqueda de anime fallida",
			anime:       "one piece",
			page:        2,
			wantError:   true,
			description: "no debe encontrar resultados debido a error simulado",
		},
		{
			name:        "búsqueda de anime con nombre vacío",
			anime:       "",
			page:        1,
			wantError:   true,
			description: "debe manejar el error al buscar un anime con nombre vacío",
		},
		{
			name:        "búsqueda de anime con página 0",
			anime:       "bleach",
			page:        0,
			wantError:   false,
			description: "debe buscar correctamente el anime y asignar página por defecto",
		},
	}

	for _, tc := range testCases {
		if err := config.InitConfig(); err != nil {
			t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
		}
		serviceAnimeflv := animeflv.NewAnimeflvService()
		ctx := context.Background()

		_, err := serviceAnimeflv.SearchAnime(ctx, tc.anime, tc.page)
		if (err != nil) != tc.wantError {
			t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
		}
	}
}

func TestSearchService(t *testing.T) {
	testCases := []struct {
		name        string
		wantError   bool
		page        uint
		description string
	}{
		{
			name:        "búsqueda de todos los animes exitosa",
			wantError:   false,
			page:        5,
			description: "debe buscar correctamente todos los animes sin errores",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceAnimeflv := animeflv.NewAnimeflvService()
			ctx := context.Background()

			_, err := serviceAnimeflv.Search(ctx, uint(tc.page))
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
			}
		})
	}
}

func TestAnimeInfoService(t *testing.T) {
	testCases := []struct {
		name        string
		animeID     string
		wantError   bool
		description string
	}{
		{
			name:        "información de anime exitosa",
			animeID:     "chainsaw-man-movie-rezehen",
			wantError:   false,
			description: "debe obtener correctamente la información del anime sin errores",
		},
		{
			name:        "información de anime exitosa",
			animeID:     "haikyuu",
			wantError:   false,
			description: "debe obtener correctamente la información del anime sin errores",
		},
		{
			name:        "información de anime fallida",
			animeID:     "unknown-anime-id",
			wantError:   true,
			description: "debe manejar el error al obtener información de un anime inexistente",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceAnimeflv := animeflv.NewAnimeflvService()
			ctx := context.Background()

			_, err := serviceAnimeflv.AnimeInfo(ctx, tc.animeID)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
			}
		})
	}
}

func TestLinksService(t *testing.T) {
	testCases := []struct {
		name        string
		animeID     string
		episode     uint
		wantError   bool
		description string
	}{
		{
			name:        "enlaces de episodio exitosos",
			animeID:     "naruto-shippuden-hd",
			episode:     200,
			wantError:   false,
			description: "debe obtener correctamente los enlaces del episodio sin errores",
		},
		{
			name:        "enlaces de episodio fallidos",
			animeID:     "unknown-anime-id",
			episode:     1,
			wantError:   true,
			description: "debe manejar el error al obtener enlaces de un episodio inexistente",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceAnimeflv := animeflv.NewAnimeflvService()
			ctx := context.Background()

			_, err := serviceAnimeflv.Links(ctx, tc.animeID, tc.episode)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
			}
		})
	}
}

func TestRecentService(t *testing.T) {
	testCases := []struct {
		name        string
		wantError   bool
		description string
	}{
		{
			name:        "obtención de animes y episodios recientes exitosa",
			wantError:   false,
			description: "debe obtener correctamente los animes y episodios recientes sin errores",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serviceAnimeflv := animeflv.NewAnimeflvService()
			ctx := context.Background()

			_, err := serviceAnimeflv.RecentEpisode(ctx)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
			}

			_, err = serviceAnimeflv.RecentAnime(ctx)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
			}
		})
	}
}
